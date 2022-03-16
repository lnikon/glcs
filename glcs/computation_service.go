package glcs

type ComputationService struct {
}

func NewComputationService() *ComputationService {
    return &ComputationService{}
}
 
func (*ComputationService) Start(desc *ComputationDescription) error {
    return nil
}

func (*ComputationService) Status(name string) (*ComputationStatus, error) {
    return nil, nil
}

func (*ComputationService) Result(name string) (*ComputationResult, error) {
    return nil, nil
}

func (*ComputationService) Stop(name string) error {
    return nil
}

