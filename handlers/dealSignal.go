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
}

func dealRelocaliseFinish(sceneName string) {
	names := strings.Fields(sceneName)
	scene1, scene2 := names[0], names[1]

	poseFileName := scene1 + "-" + scene2 + ".txt"
	poseFile, err := os.Open(poseFileName)
	if err != nil {
		panic(err)
	}
	defer poseFile.Close()
	scanner := bufio.NewScanner(poseFile)
	poses := [2]pose{}
	// only 2 lines
	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		scenePose := strings.Fields(line)
		poseStrs := strings.Split(scenePose[1], ",")
		tmpPose := pose{}
		for j := 0; j < 4; j++ {
			for k := 0; k < 2; k++ {
				index := j*2 + k
				poseStr := strings.TrimFunc(poseStrs[index], func(r rune) bool {
					return r != '.' && !unicode.IsNumber(r)
				})
				poseData, _ := strconv.ParseFloat(poseStr, 32)
				tmpPose[j][i] = float32(poseData)
			}
		}
		poses[i] = tmpPose
	}
	url := centerServerAddr + "/sys/relocalise"
	globalpose := globalPose{
		scene1,
		poses[0],
		scene2,
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
		case sceneName := <-renderFinish:
			dealRelocaliseFinish(sceneName)
		}
	}
}
