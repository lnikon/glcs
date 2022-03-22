package computation

import (
	"bytes"
	"os/exec"
)

type Computation struct {
	description *ComputationDescription
	cmd         *exec.Cmd
	result      bytes.Buffer
}

type ComputationCache struct {
	data []Computation
}
