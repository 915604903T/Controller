package handlers

import (
	"log"
	"os"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

var ClientId string = "1"
var CUDA_DEVICE string
var CenterServerAddr string = "http://127.0.0.1:23333"
var HostAddr string

var RenderFinish chan string
var RelocaliseFinish chan relocaliseInfo //two scene

// var copyLock sync.RWMutex

type pose [4][2]float64

type globalPose struct {
	Scene1Name string `json:"scene1name"`
	Scene1Ip   string `json:"scene1ip"`
	Scene1Pose pose   `json:"scene1pose"`
	Scene2Name string `json:"scene2name"`
	Scene2Ip   string `json:"scene2ip"`
	Scene2Pose pose   `json:"scene2pose"`
}

type relocaliseInfo struct {
	Scene1Name string `json:"scene1name"`
	Scene1IP   string `json:"scene1ip"`
	Scene2Name string `json:"scene2name"`
	Scene2IP   string `json:"scene2ip"`
}

type ResourceInfo struct {
	GPUMemoryFree uint64    `json:"gpumemoryfree"`
	MemoryFree    uint64    `json:"memoryfree"`
	CpuUsage      []float64 `json:"cpuusage"`
}

func init() {
	RenderFinish = make(chan string)
	RelocaliseFinish = make(chan relocaliseInfo)
	CUDA_DEVICE = os.Getenv("CUDA_DEVICE")

	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
}
