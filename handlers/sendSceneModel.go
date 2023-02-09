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

		// write zip files
		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()
		archiveFiles := getFileList(sceneName)
		for _, fileName := range archiveFiles {
			log.Println("zip ", fileName)
			file, err := os.Open(fileName)
			if err != nil {
				panic(err)
			}
			tmpWriter, err := zipWriter.Create(fileName)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(tmpWriter, file)
			if err != nil {
				panic(err)
			}
			file.Close()
		}
	}
}
