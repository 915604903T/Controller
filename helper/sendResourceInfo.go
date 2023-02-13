package helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/915604903T/ModelController/handlers"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func sendResourceInfo() {
	for {
		resourceInfo := ResourceInfo{}

		// get gpu info
		cudaDevice, _ := strconv.Atoi(handlers.CUDA_DEVICE)
		device, ret := nvml.DeviceGetHandleByIndex(cudaDevice)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device index %d: %v", cudaDevice, nvml.ErrorString(ret))
		}
		memoryInfo, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device memory info %d: %v", cudaDevice, nvml.ErrorString(ret))
		}
		resourceInfo.GPUMemoryFree = memoryInfo.Free

		// get cpu info
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			panic(err)
		}
		resourceInfo.CpuUsage = cpuPercent

		// get memory info
		memory, err := mem.VirtualMemory()
		if err != nil {
			panic(err)
		}
		resourceInfo.MemoryFree = memory.Free

		// send reqeust to center server
		resourceInfoStr, err := json.Marshal(resourceInfo)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewBuffer(resourceInfoStr)
		url := handlers.CenterServerAddr + "/sys/client/" + handlers.ClientId
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

		// every 10 second trigger once
		time.Sleep(time.Second * 10)
	}
}
