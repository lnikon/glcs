package computation

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/lnikon/glcs/common"

	"github.com/go-kit/log"
)

type ComputationLaunchParameterType string

const (
	Replicas    ComputationLaunchParameterType = "n"
	VertexCount ComputationLaunchParameterType = "vertex-count"
	Percentage  ComputationLaunchParameterType = "percentage"
)

type ComputationService struct {
	upcxxRunner     string
	pgasGraphRunner string

	Logger log.Logger
	db     *ComputationServiceDbConnector
	cache  *Computations
}

func NewComputationService(logger log.Logger) (*ComputationService, error) {
	db, err := NewComputationServiceDbConnector()
	if err != nil {
		logger.Log("NewComputationService", "Failed", "Error", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}

	upcxxRunnerBinary, lookErr := exec.LookPath(string(common.UPCXXRun))
	if lookErr != nil {
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", lookErr)
	}

	pgasGraphRunnerBinary, lookErr := exec.LookPath(string(common.PGASGraphRun))
	if lookErr != nil {
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", lookErr)
	}

	return &ComputationService{
		upcxxRunner:     upcxxRunnerBinary,
		pgasGraphRunner: pgasGraphRunnerBinary,
		Logger:          logger,
		db:              db,
		cache:           &Computations{}}, nil
}

func (cs *ComputationService) Start(ctx context.Context, desc *ComputationDescription) (*ComputationStatus, error) {
	computation, err := cs.launchComputation(desc)
	if err != nil {
		cs.Logger.Log("Start", "Failed", "Error", err)
		return nil, err
	}

	// Store in cache
	cs.cache.Append(computation)

	return &ComputationStatus{Status: Starting}, nil
}

func (*ComputationService) Status(ctx context.Context, name string) (*ComputationStatus, error) {
	return &ComputationStatus{Status: Undefined}, nil
}

func (cs *ComputationService) Result(ctx context.Context, name string) (*ComputationResult, error) {
	result, err := cs.db.ReadComputationFromDb(name)
	if err != nil {
		return nil, err
	}

	return &ComputationResult{Status: Finished, Result: result}, nil
}

func (*ComputationService) Stop(ctx context.Context, name string) (*ComputationStatus, error) {
	return &ComputationStatus{Status: "Undefined"}, nil
}

// Asyncly launch computation.
func (cs *ComputationService) launchComputation(desc *ComputationDescription) (*Computation, error) {
	result := bytes.Buffer{}
	upcxxCmd := exec.Command(cs.upcxxRunner, cs.constructLunchArguments(desc)...)
	fmt.Println("args", cs.constructLunchArguments(desc))
	upcxxCmd.Stdout = &result
	err := upcxxCmd.Start()
	if err != nil {
		cs.Logger.Log("launchComputation", Failed, "Error", err)
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", err)
	}

	computation := NewComputation(desc, upcxxCmd, &result)
	if err := cs.db.WriteNewComputationIntoDb(computation); err != nil {
		cs.Logger.Log("WriteNewComputationIntoDb", "Failed", "Error", err)
		if err := computation.Kill(); err != nil {
			cs.Logger.Log("ProcessKill", "Failed", "Error", err)
			return nil, err
		}
		return nil, err
	}

	go cs.watchComputation(computation)

	return computation, nil
}

// Waits until computation is finished then updates its status in DB.
func (cs *ComputationService) watchComputation(computation *Computation) {
	if err := computation.Wait(); err != nil {
		cs.Logger.Log("watchComputation", Failed, "Name", computation.Description().Name)
		if err := cs.db.UpdateComputationStatusInDb(computation.Description().Name, Failed, ""); err != nil {
			cs.Logger.Log("UpdateComputationStatusInDb", Failed, "Error", err)
		}
		return
	}

	cs.Logger.Log("Computation", Finished, "Name", computation.Description().Name, "Result", computation.Result().String())
	if err := cs.db.UpdateComputationStatusInDb(computation.Description().Name, Finished, computation.Result().String()); err != nil {
		cs.Logger.Log("UpdateComputationStatusInDb", Failed, "Error", err)
	}
}
