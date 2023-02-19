package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func generateGlobalPose(poseFileName string, relocInfo relocaliseInfo) []byte {
	poseFile, err := os.Open(poseFileName)
	if err != nil {
		panic(err)
	}
	defer poseFile.Close()

	scanner := bufio.NewScanner(poseFile)
	dqPose := [2][4]float64{}
	// 2 lines totally, only read the first line, the second line is world scene
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
	globalpose := globalPose{
		relocInfo.Scene1Name,
		relocInfo.Scene1IP,
		relocInfo.Scene2Name,
		relocInfo.Scene2IP,
		pose,
	}
	globalPoseStr, err := json.Marshal(globalpose)
	if err != nil {
		log.Println("marshal globalpose err: ", err)
		panic(err)
	}
	return globalPoseStr
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
		buf = bytes.NewBuffer(globalPose)
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
