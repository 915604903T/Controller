package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func requestAndSaveFile(fileName, clientAddr string) {
	name := strings.Split(fileName, ".")
	url := clientAddr + "/filemesh/" + name[0]
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("request ", fileName, " err: ", err)
		panic(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("receive response ", fileName, " err: ", err)
		panic(err)
	}
	defer resp.Body.Close()
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Println("create ", fileName, " err: ", err)
		panic(err)
	}
	io.Copy(file, resp.Body)
}

func doMergeMesh(mergeMeshInfo MergeMeshInfo) {
	// if mesh file does not exist, request file from another client
	client1, client2 := mergeMeshInfo.Mesh1.Client, mergeMeshInfo.Mesh2.Client
	mesh1FileName, mesh2FileName := mergeMeshInfo.Mesh1.FileName, mergeMeshInfo.Mesh2.FileName
	if HostAddr != client1 {
		requestAndSaveFile(mesh1FileName, client1)
	} else if HostAddr != client2 {
		requestAndSaveFile(mesh2FileName, client2)
	}

	// process pose and file name
	names := []string{}
	newSceneMap := make(map[string]bool)
	for name := range mergeMeshInfo.Mesh1.Scenes {
		names = append(names, name)
		newSceneMap[name] = true
	}
	for name := range mergeMeshInfo.Mesh2.Scenes {
		names = append(names, name)
		newSceneMap[name] = true
	}
	namePre := strings.Join(names, "-")
	poseFileName := namePre + ".txt"
	writeTmpPoseFile(poseFileName, mergeMeshInfo.PoseMatrix)
	mergeFileName := namePre + ".ply"

	// exec merge mesh program
	cmd := exec.Command("python", "mergeMesh.py",
		"--file1", mesh1FileName,
		"--file2", mesh2FileName,
		"--pose", poseFileName,
		"--output", mergeFileName)
	fmt.Println("relocalise cmd args: ", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
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
		log.Println("exec spaintgui-relocalise error: ", err)
	}

	// Remove unnecessary files
	err = os.Remove(mesh1FileName)
	if err != nil {
		log.Println("remove ", mesh1FileName, " err: ", err)
		panic(err)
	}
	err = os.Remove(mesh2FileName)
	if err != nil {
		log.Println("remove ", mesh2FileName, " err: ", err)
		panic(err)
	}
	err = os.Remove(poseFileName)
	if err != nil {
		log.Println("remove ", poseFileName, " err: ", err)
		panic(err)
	}
	newMeshInfo := MeshInfo{
		Scenes:     newSceneMap,
		WorldScene: mergeMeshInfo.Mesh1.WorldScene,
		FileName:   mergeFileName,
		Client:     HostAddr,
	}

	// send finish signal to MergeMeshFinish
	// MergeMeshFinish <- newMeshInfo
	dealMergeMeshFinish(newMeshInfo)
}

func MakeMergeMeshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("receive merge mesh request")
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		mergeMeshInfo := MergeMeshInfo{}
		err := json.Unmarshal(body, &mergeMeshInfo)
		if err != nil {
			log.Fatal("error de-serializing merge mesh request body: ", body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

		go doMergeMesh(mergeMeshInfo)
	}
}
