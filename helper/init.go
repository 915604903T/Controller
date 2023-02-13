package helper

import (
	"log"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type ResourceInfo struct {
	GPUMemoryFree uint64    `json:"gpumemoryfree"`
	MemoryFree    uint64    `json:"memoryfree"`
	CpuUsage      []float64 `json:"cpuusage"`
}

func init() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
}
