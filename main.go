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

	log "github.com/go-kit/log"

	"github.com/lnikon/glcs"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)

	flag.Parse()
	ctx := context.Background()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "log", log.DefaultCaller)

	srv := glcs.NewComputationService(logger)
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := glcs.Endpoints{
		StartEndpoint:  glcs.MakeStartEndpoint(srv),
		StatusEndpoint: glcs.MakeStatusEndpoint(srv),
		ResultEndpoint: glcs.MakeResultEndpoint(srv),
		StopEndpoint:   glcs.MakeStopEndpoint(srv),
	}

	go func() {
		logger.Log("listening", *httpAddr)
		handler := glcs.NewHTTPServer(ctx, endpoints, logger)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	stdlog.Fatalln(<-errChan)
}
