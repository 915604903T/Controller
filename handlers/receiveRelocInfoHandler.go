package handlers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func unzipFile(archiveName string) {
	// unzip receive file
	archive, err := zip.OpenReader(archiveName + ".zip")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := f.Name
		// log.Println("[receiveRelocInfo]: this is filePath: ", filePath)
		if f.FileInfo().IsDir() {
			log.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			log.Println("io copy error: ", err)
			panic(err)
		}
		dstFile.Close()
		fileInArchive.Close()
	}
}

func requestZipSceneFile(scene, clientAddr string) {
	// Make request to receive scene zip file
	url := clientAddr + "/relocalise/scene/" + scene
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("send request to client to request zip file: ", url)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("receive error from relocalise: ", resp_body)
		return
	}

	// Copy zip from resp body to zip file and unzip it
	archiveName := scene + ".zip"
	scene2ZipFile, err := os.Create(archiveName)
	defer scene2ZipFile.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	runtime.LockOSThread()
	io.Copy(scene2ZipFile, resp.Body)
	unzipFile(scene)
	runtime.UnlockOSThread()

	// Remove redundant file
	_, err = os.Stat(archiveName)
	if os.IsNotExist(err) {
		return
	} else {
		err = os.Remove(archiveName)
		if err != nil {
			panic(err)
		}
	}

}
func getFileAndRelocalise(relocInfo relocaliseInfo) {
	// request for zip file
	scene1, scene2 := relocInfo.Scene1Name, relocInfo.Scene2Name
	var lock *sync.Mutex
	requestFileLock.RLock()
	l, ok := requestFile[scene2]
	if ok {
		lock = l
	}
	requestFileLock.RUnlock()
	if !ok {
		requestFileLock.Lock()
		requestFile[scene2] = &sync.Mutex{}
		lock = requestFile[scene2]
		requestFileLock.Unlock()
	}
	// if scene2 file not exist, request target client to send zip files
	var start time.Time
	var d1 time.Duration
	lock.Lock()
	_, err := os.Stat(scene2)
	if os.IsNotExist(err) {
		start = time.Now()
		requestZipSceneFile(scene2, relocInfo.Scene2IP)
		d1 = time.Since(start)
	}
	lock.Unlock()

	if os.IsNotExist(err) {
		log.Println("[getFileAndRelocalise] Request zip file", scene2, "cost", d1, "ms!!!!!!!!!!!!!!!!!!!!!")
		TimeCostLock.Lock()
		TimeCost[scene2+"-RequestZipFile"] = d1
		TimeCostLock.Unlock()
	}

	// run relocalise process
	cmd := exec.Command("spaintgui-relocalise",
		"-f", "collaborative_config.ini",
		"--scene1", scene1,
		"--scene2", scene2,
		"-s", scene1, "-t", "Disk",
		"-s", scene2, "-t", "Disk",
		"--saveMeshOnExit")
	cmd.Env = append(cmd.Env, "CUDA_VISIBLE_DEVICES="+CUDA_DEVICE)
	pathEnv := os.Environ()
	cmd.Env = append(cmd.Env, pathEnv...)
	fmt.Println("relocalise cmd args: ", cmd.Args)
	/*
		stdout, err := cmd.StdoutPipe()
		defer stdout.Close()
		if err != nil {
			panic(err)
		}
		cmd.Stderr = cmd.Stdout
	*/
	start = time.Now()
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	/*
		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				break
			}
		}
	*/
	if err = cmd.Wait(); err != nil {
		log.Println("exec spaintgui-relocalise error: ", err)

		d2 := time.Since(start)
		log.Println("[getFileAndRelocalise] run relocalise unsuccessfully", scene1, scene2, "cost", d2, "ms!!!!!!!!!!!!!!!!!!!!!")
		TimeCostLock.Lock()
		TimeCost[scene1+"-"+scene2+"-ReocaliseUnsuccessful"] = d2
		TimeCostLock.Unlock()

		dealFailedReloclise()
		return
	}

	// send finish signal to RelocaliseFinish
	// RelocaliseFinish <- relocInfo
	d2 := time.Since(start)
	log.Println("[getFileAndRelocalise] run relocalise", scene1, scene2, "cost", d2, "ms!!!!!!!!!!!!!!!!!!!!!")
	TimeCostLock.Lock()
	TimeCost[scene1+"-"+scene2+"-Reocalise"] = d2
	TimeCostLock.Unlock()
	dealRelocaliseFinish(relocInfo)
}

func MakeReceiveRelocInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("receive request from server for relocalise")
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		relocInfo := relocaliseInfo{}
		err := json.Unmarshal(body, &relocInfo)
		if err != nil {
			log.Fatal("error de-serializing request body: ", body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Receive relocalise info: ", relocInfo.Scene1Name, relocInfo.Scene2Name)
		w.WriteHeader(http.StatusOK)
		// run relocalise
		go getFileAndRelocalise(relocInfo)
	}
}
