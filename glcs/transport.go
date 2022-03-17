package glcs

import (
	"context"
	"encoding/json"
	"net/http"
)

type startRequest struct {
	Algorithm   Algorithm `json:"algorithm"`
	Name        string    `json:"name"`
	VertexCount uint32    `json:"vertexCount"`
	Density     uint32    `json:"density"`
}

type startResponse struct {
	Status ComputationStatus `json:"status"`
}

type statusRequest struct {
	Name string `json:"name"`
}

type statusResponse struct {
	Status ComputationStatus `json:"status"`
}

type resultRequest struct {
	Name string `json:"name"`
}

type resultResponse struct {
	Result ComputationResult `json:"result"`
}

type stopRequest struct {
	Name string `json:"name"`
}

type stopResponse struct {
	Status ComputationStatus `json:"status"`
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
