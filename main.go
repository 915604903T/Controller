package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/915604903T/ModelController/handlers"

	"github.com/gorilla/mux"
)

const NameExpression = "-a-zA-Z_0-9."

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/render/scene/{name}", handlers.MakeReceiveFileHandler())
	router.HandleFunc("/relocalise/info", handlers.MakeReceiveRelocInfoHandler())
	router.HandleFunc("/relocalise/scene/{name}", handlers.MakeSendSceneModelHandler())

	// start server but not block
	go handlers.DealSignal()

	// regard tcp port as env viarable
	tcpPort, _ := strconv.Atoi(os.Getenv("tcpport"))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", tcpPort),
		Handler: router,
	}
	log.Fatal(s.ListenAndServe())
}
