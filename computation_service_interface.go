package main

import "context"

type Algorithm string

const (
	Prim    Algorithm = "Prim"
	Kruskal           = "Kruskal"
)

type ComputationStatusValue string

const (
	Undefined                    ComputationStatusValue = "Undefined"
	NotStarted                                          = "NotStarted"
	Starting                                            = "Starting"
	Started                                             = "Started"
	FailedToUpdateDb                                    = "FailedToUpdateDb"
	AnotherComputationInProgress                        = "AnotherCommunicationInProcess"
	InProgress                                          = "Started"
	Finished                                            = "Finished"
	Failed                                              = "Failed"
)

type ComputationStatus struct {
	Status ComputationStatusValue
}

type ComputationResult struct {
	Status ComputationStatusValue `json:"status"`
	Result string
	// TODO: More fields to be added.
}

type ComputationDescription struct {
	Name        string    `json:"name"`
	Algorithm   Algorithm `json:"algorithm"`
	VertexCount uint32    `json:"vertexCount"`
	Density     uint32    `json:"density"`
	Replicas    uint32    `json:"replicas"`
}

type ComputationServiceInterface interface {
	Start(ctx context.Context, desc *ComputationDescription) (*ComputationStatus, error)
	Status(ctx context.Context, name string) (*ComputationStatus, error)
	Result(ctx context.Context, name string) (*ComputationResult, error)
	Stop(ctx context.Context, name string) (*ComputationStatus, error)
}
