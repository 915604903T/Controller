package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func SendResourceInfo() {
	for {
		// copyLock.RLock()

		resourceInfo := ResourceInfo{}
		// get gpu info
		cudaDevice, _ := strconv.Atoi(CUDA_DEVICE)
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
		resourceInfo.MemoryFree = memory.Available

		// send reqeust to center server
		resourceInfoStr, err := json.Marshal(resourceInfo)
		if err != nil {
			panic(err)
		}
		log.Println("this is resourceInfo: ", string(resourceInfoStr))
		buf := bytes.NewBuffer(resourceInfoStr)
		url := CenterServerAddr + "/sys/client/" + ClientId
		request, err := http.NewRequest("GET", url, buf)
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			// panic(err)
			log.Fatal("send request ", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			resp_body, _ := ioutil.ReadAll(resp.Body)
			log.Fatal("receive error from globalpose center: ", resp_body)
			return
		}

		// copyLock.RUnlock()

		// every 10 second trigger once
		time.Sleep(time.Second * 5)
	}
}
