package glcs

type Algorithm string

const (
	Prim    Algorithm = "Prim"
	Kruskal           = "Kruskal"
)

type ComputationStatus string

const (
	Starting   ComputationStatus = "Starting"
	Started                      = "Started"
	InProgress                   = "Started"
	Finished                     = "Finished"
	Failed                       = "Failed"
)
type ComputationResult struct {
    status ComputationStatus
    // TODO: More fields to be added.
}

type ComputationDescription struct {
	Name        string
	Algorithm   Algorithm
	VertexCount uint32
	Density     uint32
}

type ComputationServiceInterface interface {
	Start(desc *ComputationDescription) error
	Status(name string) (*ComputationStatus, error)
	Result(name string) (*ComputationResult, error)
    Stop(name string) error
}
