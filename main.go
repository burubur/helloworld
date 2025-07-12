package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
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

var (
    commitHash = "000000"
    checkErr   = func(err error) {
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            slog.Error("something error while writing response", slog.Any("error", err))
        }
    }
)

func main() {
    var logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    var customSlog = slog.New(logHandler)
    slog.SetDefault(customSlog)

    slog.Debug("started open API server version: %s on PORT: %d", commitHash, port)
    response := resp{
        Code:    "SUCCESS",
        Message: "Hello world!",
        Version: commitHash,
    }
    marshaledResp, err := json.MarshalIndent(response, "", "  ")
    if err != nil {
        slog.Error("something error while marshaling response", slog.Any("error", err))
    }

    pong := resp{
        Code:    "PONG",
        Message: "pong",
        Version: commitHash,
    }
    marshaledPong, err := json.MarshalIndent(pong, "", "  ")
    if err != nil {
        slog.Error("something error while marshaling response", slog.Any("error", err))
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
        slog.Error("something error while shutting down the server", slog.Any("error", err))
    }
    slog.Debug(fmt.Sprintf("stopped open API on port: %d", port))
}

func shutdownHook() func(chan struct{}) {
    return func(processingCH chan struct{}) {
        go func(procCH chan struct{}) {
            interruptionSignal := make(chan os.Signal, 1)
            signal.Notify(interruptionSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

            // watch interruption signal
            terminationSignal := <-interruptionSignal
            slog.Debug("caught interruption signal", "signal", terminationSignal)

            // send interruption signal to the client through the registered channel
            procCH <- struct{}{}

            // stop relaying incoming signals
            signal.Stop(interruptionSignal)
        }(processingCH)
    }
}
