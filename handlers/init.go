package handlers

import (
	"log"
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

var centerServerAddr string = "http://127.0.0.1:23333"
var renderFinish chan string
var relocaliseFinish chan string //two scene

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
	renderFinish = make(chan string)
	relocaliseFinish = make(chan string)
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
		}

		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at index %d: %v", i, nvml.ErrorString(ret))
		}

		fmt.Printf("%v\n", uuid)
	}
}
