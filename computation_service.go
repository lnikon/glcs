package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-kit/log"
)

type ComputationLaunchParameterType string

const (
	Replicas    ComputationLaunchParameterType = "n"
	VertexCount ComputationLaunchParameterType = "vertex-count"
	Percentage  ComputationLaunchParameterType = "percentage"
)

type Computation struct {
	description *ComputationDescription
	cmd         *exec.Cmd
	result      bytes.Buffer
}

type ComputationCache struct {
	data []Computation
}

type ComputationService struct {
	Logger          log.Logger
	db              *ComputationServiceDbConnector
	upcxxRunner     string
	pgasGraphRunner string

	description *ComputationDescription
	cmd         *exec.Cmd
	result      bytes.Buffer
}

func NewComputationService(logger log.Logger) (*ComputationService, error) {
	db, err := NewComputationServiceDbConnector()
	if err != nil {
		logger.Log("NewComputationService", "Failed", "Error", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}

	upcxxRunnerBinary, lookErr := exec.LookPath(string(UPCXXRun))
	if lookErr != nil {
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", lookErr)
	}

	pgasGraphRunnerBinary, lookErr := exec.LookPath(string(PGASGraphRun))
	if lookErr != nil {
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", lookErr)
	}

	return &ComputationService{Logger: logger, db: db, upcxxRunner: upcxxRunnerBinary, pgasGraphRunner: pgasGraphRunnerBinary}, nil
}

func (cs *ComputationService) Start(ctx context.Context, desc *ComputationDescription) (*ComputationStatus, error) {
	var err error

	cs.description = desc
	cs.cmd, err = cs.launchComputation(desc)
	if err = cs.db.WriteNewComputationIntoDb(desc); err != nil {
		cs.Logger.Log("Start", "Failed", "Error", err)
		if err := cs.cmd.Process.Kill(); err != nil {
			cs.Logger.Log("ProcessKill", "Failed", "Error", err)
			return nil, err
		}
		return nil, err
	}

	go cs.watchComputation()

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

func constructParameter(name string) string {
	return strings.Join([]string{"--", name}, "")
}

func (cs *ComputationService) constructLunchArguments(desc *ComputationDescription) []string {
	return []string{strings.Join([]string{string(Replicas), fmt.Sprintf("%d", desc.Replicas)}, ""), cs.pgasGraphRunner,
		constructParameter(string(VertexCount)), fmt.Sprintf("%d", desc.VertexCount),
		constructParameter(string(Percentage)), fmt.Sprintf("%d", desc.Density)}

	// TODO: Do we ever gonna need this?
	// constructParameter("export-path"), constructResultFileName(desc.Name)}
}

// Launches computation using provided description and returns. Doesn't waits for it to finish.
func (cs *ComputationService) launchComputation(desc *ComputationDescription) (*exec.Cmd, error) {
	upcxxCmd := exec.Command(cs.upcxxRunner, cs.constructLunchArguments(desc)...)
	upcxxCmd.Stdout = &cs.result
	execErr := upcxxCmd.Start()
	if execErr != nil {
		cs.Logger.Log("launchComputation", Failed, "Error", execErr)
		return nil, fmt.Errorf("Unable to start new UPCXX computation, error=%s", execErr)
	}

	return upcxxCmd, nil
}

// Waits until computation is finished then updates its status in DB.
func (cs *ComputationService) watchComputation() {
	if err := cs.cmd.Wait(); err != nil {
		cs.Logger.Log("watchComputation", Failed, "Name", cs.description.Name)
		if err := cs.db.UpdateComputationStatusInDb(cs.description.Name, Failed, ""); err != nil {
			cs.Logger.Log("UpdateComputationStatusInDb", Failed, "Error", err)
		}
		return
	}

	cs.Logger.Log("Computation", Finished, "Name", cs.description.Name, "Result", cs.result.String())
	if err := cs.db.UpdateComputationStatusInDb(cs.description.Name, Finished, cs.result.String()); err != nil {
		cs.Logger.Log("UpdateComputationStatusInDb", Failed, "Error", err)
	}
}
