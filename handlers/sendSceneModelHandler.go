package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

func getFileList(directory string) []string {
	fileList := []string{}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		path := filepath.Join(directory, name)
		if !file.IsDir() {
			if strings.Contains(name, "txt") || strings.Contains(directory, "model") {
				fileList = append(fileList, path)
			}
		} else if strings.Contains(path, "model") {
			subDirFileList := getFileList(path)
			fileList = append(fileList, subDirFileList...)
		}
	}
	return fileList
}

func MakeSendSceneModelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sceneName := vars["name"]
		log.Print("request zip file ", sceneName, " file for relocalise")
		defer r.Body.Close()

		// add pose File and model to archieve file
		archieveName := sceneName + ".zip"
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archieveName))

		_, err := os.Stat(archieveName)
		if os.IsNotExist(err) {
			// write zip files
			zipWriter := zip.NewWriter(w)
			defer zipWriter.Close()
			archiveFiles := getFileList(sceneName)
			log.Println("[MakeSendSceneModelHandler] create ", archieveName)
			runtime.LockOSThread()
			for _, fileName := range archiveFiles {
				// log.Println("zip ", fileName)
				tmpWriter, err := zipWriter.Create(fileName)
				if err != nil {
					log.Println("create zip file error: ", err)
					panic(err)
				}
				file, err := os.Open(fileName)
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(tmpWriter, file)
				if err != nil {
					log.Println("write zip file error: ", err)
					panic(err)
				}
				file.Close()
			}
			runtime.UnlockOSThread()
		} else {
			zipFile, err := os.Open(archieveName)
			if err != nil {
				log.Println("open archive file", archieveName, " error: ", err)
				panic(err)
			}
			_, err = io.Copy(w, zipFile)
		}

	}
}
