package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type GracefulServer struct {
	Server           *http.Server
	shutdownFinished chan struct{}
}

func (s *GracefulServer) ListenAndServe() (err error) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	err = s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		// expected error after calling Server.Shutdown().
		err = nil
	} else if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
		return
	}

	log.Println("[ListenAndServe] waiting for shutdown finishing...")
	<-s.shutdownFinished
	log.Println("[ListenAndServe] shutdown finished")

	return
}

func (s *GracefulServer) WaitForExitingSignal() {
	var waiter = make(chan os.Signal, 1) // buffered channel
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	// blocks here until there's a signal
	<-waiter

	log.Println("[WaitForExitingSignal] receive shutdown signal")
	for k, v := range TimeCost {
		fmt.Println(k, ": ", v, "ms")
	}

	ctx := context.Background()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println("[WaitForExitingSignal] shutting down: " + err.Error())
	} else {
		log.Println("[WaitForExitingSignal] shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
