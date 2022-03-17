package glcs

import (
	"context"
	"encoding/json"

	"github.com/go-kit/log"
)

type ComputationService struct {
	Logger log.Logger
}

func structToString(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "")
	if err != nil {
		return ""
	}
	return string(out)
}

func NewComputationService(logger log.Logger) *ComputationService {
	return &ComputationService{Logger: logger}
}

func (cs *ComputationService) Start(ctx context.Context, desc *ComputationDescription) (*ComputationStatus, error) {
	cs.Logger.Log("StartEndpoint", "Called", "Name", desc.Name, "Algorithm", desc.Algorithm, "VertexCount", desc.VertexCount, "Density", desc.Density)
	return &ComputationStatus{Status: "Starting"}, nil
}

func (*ComputationService) Status(ctx context.Context, name string) (*ComputationStatus, error) {
	return &ComputationStatus{Status: "Undefined"}, nil
}

func (*ComputationService) Result(ctx context.Context, name string) (*ComputationResult, error) {
	return &ComputationResult{Status: "InProgress"}, nil
}

func (*ComputationService) Stop(ctx context.Context, name string) (*ComputationStatus, error) {
	return &ComputationStatus{Status: "Undefined"}, nil
}
