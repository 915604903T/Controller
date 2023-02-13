package handlers

import (
	"os"
)

var ClientId string = "1"
var CUDA_DEVICE string
var CenterServerAddr string = "http://127.0.0.1:23333"

var RenderFinish chan string
var RelocaliseFinish chan string //two scene

type pose [4][2]float64

type globalPose struct {
	Scene1Name string `json:"scene1name"`
	Scene1Pose pose   `json:"scene1pose"`
	Scene2Name string `json:"scene2name"`
	Scene2Pose pose   `json:"scene2pose"`
}

type relocaliseInfo struct {
	Scene1Name string `json:"scene1name"`
	Scene1IP   string `json:"scene1ip"`
	Scene2Name string `json:"scene2name"`
	Scene2IP   string `json:"scene2ip"`
}

func init() {
	RenderFinish = make(chan string)
	RelocaliseFinish = make(chan string)
	CUDA_DEVICE = os.Getenv("CUDA_DEVICE")
}
