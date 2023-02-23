package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
)

func doRender(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = cmd.Stdout
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
}
func runRender(sceneName string) {
	// save scene files in the file with the same name
	cmd := exec.Command("spaintgui-processVoxel",
		"-f", "collaborative_config.ini",
		"--name", sceneName,
		"-s", sceneName, "-t", "Disk")
	cmd.Env = append(cmd.Env, "CUDA_VISIBLE_DEVICES="+CUDA_DEVICE)
	fmt.Println("cmd args: ", cmd.Args)
	// do render until success
	for {
		doRender(cmd)
		err := cmd.Wait()
		if err != nil {
			log.Println("exec spaintgui-processVoxel error: ", err)
			continue
		} else {
			break
		}
	}

	// RenderFinish <- sceneName
	dealRenderFinish(sceneName)
}

func MakeReceiveFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sceneName := vars["name"]
		log.Print("receive file and run render request: ", sceneName)
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
				dst, _ := os.Create(filepath.Join(sceneName, part.FileName()))
				io.Copy(dst, part)
				dst.Close()
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("save file success!"))
		go runRender(sceneName)
	}
}

