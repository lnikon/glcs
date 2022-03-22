package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lnikon/glcs/computation"
)

type startRequest struct {
	Algorithm   computation.Algorithm `json:"algorithm"`
	Name        string                `json:"name"`
	VertexCount uint32                `json:"vertexCount"`
	Density     uint32                `json:"density"`
	Replicas    uint32                `json:"replicas"`
}

type startResponse struct {
	Status computation.ComputationStatus `json:"status"`
}

type statusRequest struct {
	Name string `json:"name"`
}

type statusResponse struct {
	Status computation.ComputationStatus `json:"status"`
}

type resultRequest struct {
	Name string `json:"name"`
}

type resultResponse struct {
	Result computation.ComputationResult `json:"result"`
}

type stopRequest struct {
	Name string `json:"name"`
}

type stopResponse struct {
	Status computation.ComputationStatus `json:"status"`
}

func decodeStartRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req startRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req statusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeResultRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req resultRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeStopRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req stopRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
