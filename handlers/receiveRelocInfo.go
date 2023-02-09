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
		log.Println("[receiveRelocInfo]: this is filePath: ", filePath)
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
			panic(err)
		}
		dstFile.Close()
		fileInArchive.Close()
	}
}
func getFileAndRelocalise(relocInfo relocaliseInfo) {
	// request for zip file
	scene1, scene2 := relocInfo.Scene1Name, relocInfo.Scene2Name
	targetIp := relocInfo.Scene2IP
	url := targetIp + "/relocalise/scene/" + scene2
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
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
	scene2ZipFile, err := os.Create(scene2 + ".zip")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer scene2ZipFile.Close()
	io.Copy(scene2ZipFile, resp.Body)

	unzipFile(scene2)

	cmd := exec.Command("spaintgui-relocalise",
		"-f", "collaborative_config.ini",
		"--scene1", scene1,
		"--scene2", scene2,
		"-s", scene1, "-t", "Disk",
		"-s", scene2, "-t", "Disk")
	cmd.Env = append(cmd.Env, "CUDA_VISIBLE_DEVICES=2")
	fmt.Println("relocalise cmd args: ", cmd.Args)

    stdout, err := cmd.StdoutPipe()
    if err!=nil {
        panic(err)
    }
	cmd.Stderr = cmd.Stdout
    if err = cmd.Start(); err != nil {
        panic(err)
    }
    for {
        tmp := make([]byte, 1024)
        _, err := stdout.Read(tmp)
        fmt.Print(string(tmp))
        if err != nil {
            break
        }
    }
    if err = cmd.Wait(); err != nil {
        panic(err)
    }
	// err = cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }
	relocaliseFinish <- scene1 + " " + scene2
}
func MakeReceiveRelocInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		relocInfo := relocaliseInfo{}
		err := json.Unmarshal(body, &relocInfo)
		if err != nil {
			log.Fatal("error de-serializing request body: ", body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Print("Receive relocalise info: ", relocInfo.Scene1Name, relocInfo.Scene2Name)
		w.WriteHeader(http.StatusOK)
		// run relocalise
		go getFileAndRelocalise(relocInfo)
	}
}
