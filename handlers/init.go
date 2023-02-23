package handlers

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func init() {
	// RenderFinish = make(chan string)
	// RelocaliseFinish = make(chan relocaliseInfo)
	// MergeMeshFinish = make(chan MeshInfo)
	requestFile = make(map[string]*sync.Mutex)

	TimeCost = make(map[string]time.Duration)

	CUDA_DEVICE = os.Getenv("CUDA_DEVICE")

	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
}
