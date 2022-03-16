package glcs

import (
    "context"
    "net/http"
    "os"
    stdlog "log"

    httptransport "github.com/go-kit/kit/transport/http"
    log "github.com/go-kit/log"

    "github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
    var logger log.Logger
    logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
    stdlog.SetOutput(log.NewStdlibAdapter(logger))
    logger = log.With(logger, "ts", log.DefaultTimestampUTC, "log", log.DefaultCaller)

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
