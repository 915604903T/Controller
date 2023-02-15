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
	buf := bytes.NewBuffer([]byte(HostAddr))
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
		poseFile, err := os.Open(poseFileName)
		if err != nil {
			panic(err)
		}
		defer poseFile.Close()
		scanner := bufio.NewScanner(poseFile)
		poses := [2]pose{}
		scenes := [2]string{}
		scenesIp := [2]string{}
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
		if scenes[0] == relocInfo.Scene1Name {
			scenesIp[0], scenesIp[1] = relocInfo.Scene1IP, relocInfo.Scene2IP
		} else {
			scenesIp[0], scenesIp[1] = relocInfo.Scene2IP, relocInfo.Scene1IP
		}
		globalpose := globalPose{
			scenes[0],
			scenesIp[0],
			poses[0],
			scenes[1],
			scenesIp[1],
			poses[1],
		}
		globalposeStr, err := json.Marshal(globalpose)
		if err != nil {
			panic(err)
		}
		buf = bytes.NewBuffer(globalposeStr)
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
		log.Fatal("receive error from globalpose center: ", resp_body)
		return
	}
}

func DealSignal() {
	for {
		select {
		case sceneName := <-RenderFinish:
			// copyLock.RLock()
			dealRenderFinish(sceneName)
			// copyLock.RUnlock()
		case relocInfo := <-RelocaliseFinish:
			// copyLock.RLock()
			dealRelocaliseFinish(relocInfo)
			// copyLock.RUnlock()
		}
	}
}
