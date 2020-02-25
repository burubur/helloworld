package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	CommitHash = "000000"
	checkErr   = func(err error) {
		if err != nil && err != http.ErrServerClosed {
			log.Println("something error while writing response", err.Error())
		}
	}
)

const (
	port               = 8080
	defGracefulTimeout = 60 * time.Second
)

type resp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Version string      `json:"version"`
}

func main() {
	log.Printf("started open API server version: %s on PORT: %d", CommitHash, port)
	response := resp{
		Code:    "SUCCESS",
		Message: "Hello world!",
		Version: CommitHash,
	}
	marshaledResp, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Println("something error while marshaling response", err.Error())
	}

	pong := resp{
		Code:    "PONG",
		Message: "pong",
		Version: CommitHash,
	}
	marshaledPong, err := json.MarshalIndent(pong, "", "  ")
	if err != nil {
		log.Println("something error while marshaling response", err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(marshaledResp)
		checkErr(err)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(marshaledPong)
		checkErr(err)
	})

	server := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: mux,
	}

	mainCtx, cancel := context.WithTimeout(context.Background(), defGracefulTimeout)
	defer cancel()

	terminationHook := shutdownHook()
	terminationCH := make(chan struct{})
	terminationHook(terminationCH)
	go func() {
		err = server.ListenAndServe()
		checkErr(err)
	}()

	<-terminationCH
	err = server.Shutdown(mainCtx)
	if err != nil {
		log.Println("something error while shutting down the server", err.Error())
	}
	log.Println(fmt.Sprintf("stopped open API on port: %d", port))
}

func shutdownHook() func(chan struct{}) {
	return func(processingCH chan struct{}) {
		go func(proCH chan struct{}) {
			interruptionSignal := make(chan os.Signal, 1)
			signal.Notify(interruptionSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

			// watch interruption signal
			terminationSignal := <-interruptionSignal
			log.Println(fmt.Sprint("caught interruption signal: ", terminationSignal))

			// send interruption signal to the client through the registered channel
			proCH <- struct{}{}

			// stop relaying incoming signals
			signal.Stop(interruptionSignal)
		}(processingCH)
	}
}
