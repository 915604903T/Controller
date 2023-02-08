package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func runRender(sceneName string) {
	// save scene files in the file with the same name
	cmd := exec.Command("spaintgui-processVoxel",
		"-f", "collaborative_config.ini",
		"--name", sceneName,
		"-s", sceneName, "-t", "Disk")
	// cmd.Env = append(cmd.Env, "CUDA_VISIBLE_DEVICES=1")
	err := cmd.Run()
	if err != nil {
		log.Fatal("run ", sceneName, " error: ", err)
	}
	renderFinish <- sceneName
}

func MakeReceiveFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sceneName := vars["name"]
		log.Print("receive file and run render request: ", sceneName)
		defer r.Body.Close()

		// Create directory to save images, poses, calib.txt
		os.Mkdir(sceneName, 0644)
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
			fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
			if part.FileName() == "" { // this is FormData
				data, _ := ioutil.ReadAll(part)
				fmt.Printf("FormData=[%s]\n", string(data))
			} else { // This is FileData
				//Filename contains the directory
				dst, _ := os.Create(part.FileName())
				defer dst.Close()
				io.Copy(dst, part)
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("save file success!"))
		go runRender(sceneName)
	}
}
