package glcs

import (
	"context"
	"encoding/json"
	"net/http"
)

type startRequest struct {
	algorithm   Algorithm
	name        string
	vertexCount uint32
	density     uint32
}

type startResponse struct {
	status ComputationStatus
}

type statusRequest struct {
	name string
}

type statusResponse struct {
	status ComputationStatus
}

type resultRequest struct {
	name string
}

type resultResponse struct {
	result ComputationResult
}

type stopRequest struct {
	name string
}

type stopResponse struct {
	status ComputationStatus
}

func decodeStartRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req startRequest
	return req, nil
}

func decodeStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req statusRequest
	return req, nil
}

func decodeResultRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req resultRequest
	return req, nil
}

func decodeStopRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req stopRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
