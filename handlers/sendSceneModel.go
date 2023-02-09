package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func MakeSendSceneModelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sceneName := vars["name"]
		log.Print("request ", sceneName, " file for relocalise")
		defer r.Body.Close()

		// add pose File and model to archieve file
		archieveName := sceneName + ".zip"
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archieveName))

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		err := filepath.Walk(sceneName, func(path string, info fs.FileInfo, err error) error {
			log.Print("this is name: ", info.Name())
			if !info.IsDir() {
				fileName := info.Name()
				if strings.Contains(fileName, "model") ||
					strings.Contains(fileName, ".txt") ||
					strings.Contains(fileName, ".ini") {
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
			return nil
		})
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
	}
}
