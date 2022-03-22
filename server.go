package main

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

// Creates new HTTP server to handle incoming UPCXX computation requests
func NewHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	loggingMiddleware := LoggingMiddleware(logger)

	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.Use(loggingMiddleware)

	r.Methods("POST").Path("/computation/start").Handler(httptransport.NewServer(
		endpoints.StartEndpoint,
		decodeStartRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/computation/status").Handler(httptransport.NewServer(
		endpoints.StatusEndpoint,
		decodeStatusRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/computation/result").Handler(httptransport.NewServer(
		endpoints.ResultEndpoint,
		decodeResultRequest,
		encodeResponse,
	))

	r.Methods("POST").Path("/computation/stop").Handler(httptransport.NewServer(
		endpoints.StopEndpoint,
		decodeStopRequest,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
