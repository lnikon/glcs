package main

import (
    "context"
    "flag"
    "fmt"
    stdlog "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/lnikon/glcs"
)

func main() {
    var (
        httpAddr = flag.String("http", ":8080", "http listen address")
    )

    flag.Parse()
    ctx := context.Background()

    srv := glcs.NewComputationService()
    errChan := make(chan error)

    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errChan <- fmt.Errorf("%s", <-c)
    }()

    endpoints := glcs.Endpoints{
        StartEndpoint: glcs.MakeStartEndpoint(*srv),
        StatusEndpoint: glcs.MakeStatusEndpoint(*srv),
        ResultEndpoint: glcs.MakeResultEndpoint(*srv),
        StopEndpoint: glcs.MakeStopEndpoint(*srv),
    }

    go func() {
        stdlog.Println("glcs is listening on port:", *httpAddr)
        handler := glcs.NewHTTPServer(ctx, endpoints)
        errChan <- http.ListenAndServe(*httpAddr, handler)
    }()

    stdlog.Fatalln(<-errChan)
}

