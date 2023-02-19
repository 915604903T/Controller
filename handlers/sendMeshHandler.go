package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func MakeSendMeshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["name"]
		log.Print("request ply file ", fileName, " file for relocalise")
		defer r.Body.Close()

		// add pose File and model to archieve file
		plyFileName := fileName + ".ply"
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", plyFileName))

		plyFile, err := os.Open(plyFileName)
		if err != nil {
			log.Println("open archive file", plyFileName, " error: ", err)
			panic(err)
		}
		// will be inturrupt then panic????
		_, err = io.Copy(w, plyFile)

	}
}
