package handlers

var centerServerAddr string = "127.0.0.1:23333"
var renderFinish chan string
var relocaliseFinish chan string //two scene

type pose [4][2]float32

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
	renderFinish = make(chan string)
	relocaliseFinish = make(chan string)
}
