package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

func dealRenderFinish(sceneName string) {
	url := CenterServerAddr + "/sys/model/" + sceneName
	var buf *bytes.Buffer
	modelFile := sceneName + "/model"
	_, err := os.Stat(modelFile)
	if os.IsNotExist(err) {
		buf = bytes.NewBuffer([]byte("Failed"))
	} else {
		buf = bytes.NewBuffer([]byte(HostAddr))
	}
	request, err := http.NewRequest("GET", url, buf)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("send request to ", url)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("receive error from center server for render finish request: ", resp_body)
		return
	}
}

func generateGlobalPose(poseFileName string, relocInfo relocaliseInfo) globalPose {
	poseFile, err := os.Open(poseFileName)
	if err != nil {
		panic(err)
	}
	defer poseFile.Close()

	scanner := bufio.NewScanner(poseFile)
	dqPose := [2][4]float64{}
	// 2 lines totally, only read the first line, the second line is world scene
	scanner.Scan()
	line := scanner.Text()
	dqPoseStr := strings.Fields(line)[1]
	poseStrs := strings.Split(dqPoseStr, ",")
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			index := i*2 + j
			poseStr := strings.TrimFunc(poseStrs[index], func(r rune) bool {
				return r != '.' && r != '-' && !unicode.IsNumber(r)
			})
			poseData, _ := strconv.ParseFloat(poseStr, 64)
			dqPose[j][i] = poseData
		}
	}
	pose := NewPoseDq(dqPose)
	globalPose := globalPose{
		relocInfo.Scene1Name,
		relocInfo.Scene1IP,
		relocInfo.Scene2Name,
		relocInfo.Scene2IP,
		pose,
	}
	return globalPose
}

func mergeRelocMesh(globalpose globalPose) {
	// generate pose File and name output file
	scene1, scene2 := globalpose.Scene1Name, globalpose.Scene2Name
	namePre := scene1 + "-" + scene2
	poseFileName := "worldPose.txt"
	mergeFileName := namePre + ".ply"
	mesh1FileName, mesh2FileName := scene1+".ply", scene2+".ply"

	// run merge two scene mesh file
	cmd := exec.Command("python", "mergeMesh.py",
		"--file1", mesh1FileName,
		"--file2", mesh2FileName,
		"--pose", poseFileName,
		"--output", mergeFileName)
	fmt.Println("relocalise cmd args: ", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
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
		log.Println("exec python mergeMesh.py error: ", err)
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
}

func dealRelocaliseFinish(relocInfo relocaliseInfo) {
	scene1, scene2 := relocInfo.Scene1Name, relocInfo.Scene2Name
	poseFileName := "global_poses/" + scene1 + "-" + scene2 + ".txt"

	//if file does not exist
	isExist := true
	_, err := os.Stat(poseFileName)
	if os.IsNotExist(err) {
		log.Print(poseFileName, " not exist")
		isExist = false
	}
	var buf *bytes.Buffer
	//if file exist
	if isExist {
		globalPose := generateGlobalPose(poseFileName, relocInfo)
		mergeRelocMesh(globalPose)
		globalPoseStr, err := json.Marshal(globalPose)
		if err != nil {
			log.Println("marshal globalpose err: ", err)
			panic(err)
		}
		buf = bytes.NewBuffer(globalPoseStr)
	} else {
		buf = bytes.NewBuffer([]byte(scene1 + " " + scene2 + " " + "failed"))
	}
	url := CenterServerAddr + "/sys/relocalise"
	request, err := http.NewRequest("GET", url, buf)
	if err != nil {
		panic(err)
	}
	log.Print("Send request to ", url)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("receive error from relocalise info center: ", resp_body)
		return
	}
}

func dealMergeMeshFinish(meshInfo MeshInfo) {
	meshInfoBytes, err := json.Marshal(meshInfo)
	if err != nil {
		log.Println("marshal meshInfo err: ", err)
		panic(err)
	}
	buf := bytes.NewBuffer(meshInfoBytes)

	url := CenterServerAddr + "/sys/mesh"
	request, err := http.NewRequest("GET", url, buf)
	if err != nil {
		log.Println("Create mesh finish request to center server err: ", err)
		panic(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("receive error from mesh finish center: ", resp_body)
		return
	}
}

/*
func DealSignal() {
	for {
		select {
		case sceneName := <-RenderFinish:
			dealRenderFinish(sceneName)
		case relocInfo := <-RelocaliseFinish:
			dealRelocaliseFinish(relocInfo)
		case meshInfo := <-MergeMeshFinish:
			dealMergeMeshFinish(meshInfo)
		}
	}
}
*/
