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
	url := centerServerAddr + "/sys/model/" + sceneName
	buf := bytes.NewBuffer([]byte("OK"))
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

func dealRelocaliseFinish(sceneName string) {
	names := strings.Fields(sceneName)
	scene1, scene2 := names[0], names[1]

	poseFileName := "global_poses/" + scene1 + "-" + scene2 + ".txt"
	//if file does not exist
	_, err := os.Stat(poseFileName)
	if os.IsNotExist(err) {
		log.Print(poseFileName, " not exist")
		return
	}
	//if file exist
	poseFile, err := os.Open(poseFileName)
	if err != nil {
		panic(err)
	}
	defer poseFile.Close()
	scanner := bufio.NewScanner(poseFile)
	poses := [2]pose{}
	scenes := [2]string{}
	// only 2 lines
	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		scenePose := strings.Fields(line)
		scenes[i] = scenePose[0]
		poseStrs := strings.Split(scenePose[1], ",")
		tmpPose := pose{}
		for j := 0; j < 4; j++ {
			for k := 0; k < 2; k++ {
				index := j*2 + k
				poseStr := strings.TrimFunc(poseStrs[index], func(r rune) bool {
					return r != '.' && r != '-' && !unicode.IsNumber(r)
				})
				poseData, _ := strconv.ParseFloat(poseStr, 64)
				tmpPose[j][k] = poseData
			}
		}
		poses[i] = tmpPose
	}
	url := centerServerAddr + "/sys/relocalise"
	globalpose := globalPose{
		scenes[0],
		poses[0],
		scenes[1],
		poses[1],
	}
	globalposeStr, err := json.Marshal(globalpose)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer([]byte(globalposeStr))
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
		log.Fatal("receive error from globalpose center: ", resp_body)
		return
	}
}

func DealSignal() {
	for {
		select {
		case sceneName := <-renderFinish:
			dealRenderFinish(sceneName)
		case sceneName := <-relocaliseFinish:
			dealRelocaliseFinish(sceneName)
		}
	}
}
