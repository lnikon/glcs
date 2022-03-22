package main

import (
    "context"
    "testing"
)

func TestStart(t *testing.T) {
    srv, _ := setup()

    desc := ComputationDescription{
        Name: "TestComputation-1",
        Algorithm: Kruskal,
        VertexCount: 1024,
        Density: 95,
    }

    err := srv.Start(&desc)
    if err != nil {
        t.Errorf("Error: %s", err)
    }
}

func TestStatus(t *testing.T) {
    srv, _ := setup()

    name := "TestComputation-1"

    status, err := srv.Status(name)
    if err != nil {
        t.Errorf("Error: %s", err)
    }

    if status == nil {
        t.Errorf("Error: status is nil")
    }
}

func TestResult(t *testing.T) {
    srv, _ := setup()

    name := "TestComputation-1"

    result, err := srv.Result(name)
    if err != nil {
        t.Errorf("Error: %s", err)
    }

    if result == nil {
        t.Errorf("Error: status is nil")
    }
}

func TestStop(t *testing.T) {
    srv, _ := setup()

    name := "TestComputation-1"

    err := srv.Stop(name)
    if err != nil {
        t.Errorf("Error: %s", err)
    }
}

func setup() (srv ComputationServiceInterface, ctx context.Context) {
    return NewComputationService(), context.Background()
}
