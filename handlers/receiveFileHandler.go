package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/shirou/gopsutil/process"
)

func doRender(sceneName string, cmd *exec.Cmd) {
	/*stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = cmd.Stdout*/
	if err := cmd.Start(); err != nil {
		log.Println("[doRender] start render err:", err)
		return
	}
	go measureRender(sceneName, cmd)
	/*for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	*/
}

func measureRender(sceneName string, cmd *exec.Cmd) {
	pid := cmd.Process.Pid
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Println("[measureRender] process", pid, "not exist!!!")
		return
	}
	totalMemory := 0.0
	totalCpuUsage := 0.0
	cnt := 0
	for ; ; time.Sleep(time.Second * 2) {
		isRunning, err := p.IsRunning()
		if err != nil {
			log.Println("[measureRender] process", pid, "is running return error!!!")
			continue
		}
		if !isRunning {
			break
		}
		cpuPercent, _ := p.CPUPercent()
		memoryInfo, _ := p.MemoryInfo()
		totalCpuUsage += cpuPercent
		totalMemory += float64(memoryInfo.RSS) / 1e6
		cnt++
		// log.Println("[measureRender] ", sceneName, "cpuUsage: ", cpuPercent, "%")
		// log.Println("[measureRender] ", sceneName, "RSS: ", float64(memoryInfo.RSS)/1e6, "MB")
	}
	averageMemory := totalMemory / float64(cnt)
	averageCpu := totalCpuUsage / float64(cnt)
	index := sceneName + "-Render"
	PidResourceLock.Lock()
	MemoryCost[index] = averageMemory
	CpuUsage[index] = averageCpu
	PidResourceLock.Unlock()
}

func runRender(sceneName string) {
	// save scene files in the file with the same name
	var start time.Time
	for ; ; time.Sleep(time.Second * 5) {
		cmd := exec.Command("spaintgui-processVoxel",
			"-f", "collaborative_config.ini",
			"--name", sceneName,
			"-s", sceneName, "-t", "Disk")
		cmd.Env = append(cmd.Env, "CUDA_VISIBLE_DEVICES="+CUDA_DEVICE)
		fmt.Println("cmd args: ", cmd.Args)
		// do render until success
		start = time.Now()
		doRender(sceneName, cmd)
		err := cmd.Wait()
		if err != nil {
			log.Println("exec spaintgui-processVoxel error: ", err)
			continue
		} else {
			break
		}
	}
	duration := time.Since(start)
	log.Println("[runRender] Render", sceneName, "cost", duration, "s!!!!!!!!!!!!!!!!!!!!!")
	TimeCostLock.Lock()
	TimeCost[sceneName+"-Render"] = duration
	TimeCostLock.Unlock()

	// RenderFinish <- sceneName
	dealRenderFinish(sceneName)
}

func MakeReceiveFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		vars := mux.Vars(r)
		sceneName := vars["name"]
		log.Print("[MakeReceiveFileHandler] receive file and run render request: ", sceneName)
		defer r.Body.Close()

		// Create directory to save images, poses, calib.txt
		os.Mkdir(sceneName, 0755)
		// read multiple files
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			// fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
			if part.FileName() == "" { // this is FormData
				data, _ := ioutil.ReadAll(part)
				fmt.Printf("FormData=[%s]\n", string(data))
			} else { // This is FileData
				//Filename not contains the directory
				name := filepath.Join(sceneName, part.FileName())
				dst, _ := os.Create(name)
				var data []byte
				_, err := part.Read(data)
				if err != nil {
					log.Println(name, "part read err: ", err)
					panic(err)
				}
				img, format, _ := image.Decode(bytes.NewReader(data))
				log.Println("this is ", name, "format: ", format)
				originImg := resize.Resize(uint(img.Bounds().Max.X*2), 0, img, resize.NearestNeighbor)
				png.Encode(dst, originImg)
				// io.Copy(dst, part)
				dst.Close()
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("save file success!"))
		duration := time.Since(start)

		log.Println("[MakeReceiveFileHandler] Receive", sceneName, "cost", duration, "s!!!!!!!!!!!!!!!!!!!!!")
		TimeCostLock.Lock()
		TimeCost[sceneName+"-ReceiveUserFile"] = duration
		TimeCostLock.Unlock()

		go runRender(sceneName)
	}
}
