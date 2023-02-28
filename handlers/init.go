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
	MemoryCost = make(map[string]float64)
	CpuUsage = make(map[string]float64)

	CUDA_DEVICE = os.Getenv("CUDA_DEVICE")

	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Printf("Unable to initialize NVML: %v\n", nvml.ErrorString(ret))
		HasGPU = false
	}
}
