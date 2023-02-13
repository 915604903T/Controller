package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/915604903T/ModelController/handlers"
	"github.com/915604903T/ModelController/helper"

	"github.com/NVIDIA/go-nvml/pkg/nvml"

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

	// start send resource info to center server
	go helper.SendResourceInfo()

	// close nvml and two channel
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()
	defer close(handlers.RenderFinish)
	defer close(handlers.RelocaliseFinish)

	// regard tcp port as env viarable
	tcpPort, _ := strconv.Atoi(os.Getenv("PORT"))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", tcpPort),
		Handler: router,
	}
	log.Println("listen on port: ", tcpPort)
	log.Fatal(s.ListenAndServe())
}
