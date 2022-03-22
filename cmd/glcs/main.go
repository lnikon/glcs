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

	"github.com/lnikon/glcs/common"
	"github.com/lnikon/glcs/computation"
	"github.com/lnikon/glcs/server"

	log "github.com/go-kit/log"
)

func setupLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "log", log.DefaultCaller)
	return logger
}

func checkEnvironment(logger log.Logger) error {
	if err := common.CheckUpcxxEnvVars(); err != nil {
		logger.Log("CheckUpcxxEnvVars", "Failed", "Error", err)
		return err
	}

	if err := common.CheckDbEnvVars(); err != nil {
		logger.Log("CheckDbEnvVars", "Failed", "Error", err)
		return err
	}

	if err := common.CheckUpcxxBinaries(); err != nil {
		logger.Log("CheckUpcxxBinaries", "Failed", "Error", err)
		return err
	}

	return nil
}

func setupEndpoints(srv *computation.ComputationService) server.Endpoints {
	return server.Endpoints{
		StartEndpoint:  server.MakeStartEndpoint(srv),
		StatusEndpoint: server.MakeStatusEndpoint(srv),
		ResultEndpoint: server.MakeResultEndpoint(srv),
		StopEndpoint:   server.MakeStopEndpoint(srv),
	}
}

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()

	// Loger which will be shared through the project
	logger := setupLogger()

	// Several env vars and binaries should be correctly set and available for UPCXX
	if err := checkEnvironment(logger); err != nil {
		logger.Log("EnvironmentCheck", "Failed")
		return
	}

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	srv, err := computation.NewComputationService(logger)
	if err != nil {
		logger.Log("NewComputationService", computation.Failed, "Error", err)
		return
	}
	endpoints := setupEndpoints(srv)

	ctx := context.Background()
	go func() {
		logger.Log("listening", *httpAddr)
		handler := server.NewHTTPServer(ctx, endpoints, logger)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	stdlog.Fatalln(<-errChan)
}
