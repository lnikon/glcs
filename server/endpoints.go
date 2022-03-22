package server

import (
	"context"

	"github.com/lnikon/glcs/computation"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	StartEndpoint  endpoint.Endpoint
	StatusEndpoint endpoint.Endpoint
	ResultEndpoint endpoint.Endpoint
	StopEndpoint   endpoint.Endpoint
}

func MakeStartEndpoint(srv *computation.ComputationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(startRequest)

		desc := computation.ComputationDescription{
			Name:        r.Name,
			Algorithm:   r.Algorithm,
			VertexCount: r.VertexCount,
			Density:     r.Density,
			Replicas:    r.Replicas,
		}

		status, err := srv.Start(ctx, &desc)
		if err != nil {
			return nil, err
		}

		return startResponse{*status}, nil
	}
}

func MakeStatusEndpoint(srv *computation.ComputationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(statusRequest)

		status, err := srv.Status(ctx, r.Name)
		if err != nil {
			return nil, err
		}

		return statusResponse{*status}, nil
	}
}

func MakeResultEndpoint(srv *computation.ComputationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(resultRequest)

		result, err := srv.Result(ctx, r.Name)
		if err != nil {
			srv.Logger.Log("ResultEndpoint", "Failed", "Error", err)
			return nil, err
		}

		return resultResponse{*result}, nil
	}
}

func MakeStopEndpoint(srv *computation.ComputationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(stopRequest)

		resp, err := srv.Stop(ctx, r.Name)
		if err != nil {
			return nil, err
		}

		return stopResponse{*resp}, nil
	}
}

// func (e Endpoints) Start(ctx context.Context, desc *ComputationDescription) (*ComputationStatus, error) {
// 	req := startRequest{
// 		Algorithm:   desc.Algorithm,
// 		Name:        desc.Name,
// 		VertexCount: desc.VertexCount,
// 		density:     desc.Density,
// 	}
//
// 	resp, err := e.StartEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	startResp := resp.(startResponse)
// 	return &startResp.status, nil
// }
//
// func (e Endpoints) Status(ctx context.Context, name string) (*ComputationStatus, error) {
// 	req := statusRequest{
// 		name: name,
// 	}
//
// 	resp, err := e.StatusEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// statusResp := statusResponse{}
// 	statusResp := resp.(statusResponse)
// 	return &statusResp.status, nil
// }
//
// func (e Endpoints) Result(ctx context.Context, name string) (*ComputationResult, error) {
// 	req := resultRequest{
// 		name: name,
// 	}
//
// 	resp, err := e.ResultEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	resultResp := resp.(resultResponse)
// 	return &resultResp.result, nil
// }
//
// func (e Endpoints) Stop(ctx context.Context, name string) (*ComputationStatus, error) {
// 	req := stopRequest{
// 		name: name,
// 	}
//
// 	resp, err := e.StopEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	stopResp := resp.(stopResponse)
// 	return &stopResp.status, nil
// }
