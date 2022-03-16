package glcs

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	StartEndpoint  endpoint.Endpoint
	StatusEndpoint endpoint.Endpoint
	ResultEndpoint endpoint.Endpoint
	StopEndpoint   endpoint.Endpoint
}

func MakeStartEndpoint(srv ComputationService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        r := request.(startRequest)

        desc := ComputationDescription{
            Name: r.name,
            Algorithm: r.algorithm,
            VertexCount: r.vertexCount,
            Density: r.density,
        }

        err := srv.Start(&desc)
        if err != nil {
            return nil, err
        }

        return startResponse{Starting}, nil
    }
}

func MakeStatusEndpoint(srv ComputationService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        r := request.(statusRequest)

        status, err := srv.Status(r.name)
        if err != nil {
            return nil, err
        }

        return statusResponse{*status}, nil
    }
}

func MakeResultEndpoint(srv ComputationService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        r := request.(resultRequest)

        result, err := srv.Result(r.name)
        if err != nil {
            return nil, err
        }

        return resultResponse{*result}, nil
    }
}
 
func MakeStopEndpoint(srv ComputationService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        r := request.(stopRequest)

        err := srv.Stop(r.name)
        if err != nil {
            return nil, err
        }

        return stopResponse{}, nil
    }
}

