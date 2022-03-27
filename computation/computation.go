package computation

import (
	"bytes"
	"os/exec"
	"sync"
)

type Computation struct {
	description *ComputationDescription
	cmd         *exec.Cmd

	result *bytes.Buffer
	mutex  sync.Mutex
}

type Computations struct {
	data  []*Computation
	mutex sync.Mutex
}

func NewComputation(desc *ComputationDescription, cmd *exec.Cmd, result *bytes.Buffer) *Computation {
	return &Computation{description: desc, cmd: cmd, result: result}
}

// Async start upcxx process
func (c *Computation) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.cmd.Start()
}

// Block wait on undelying upcxx process
func (c *Computation) Wait() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, err := c.cmd.Process.Wait()
	return err
}

// Kill undelying upcxx process
func (c *Computation) Kill() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.cmd.Process.Kill()
}

// Returns copy of computation result
func (c *Computation) Result() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.result.String()
}

// Returns copy of computation description
func (c *Computation) Description() ComputationDescription {
	return *c.description
}

// Append to computations collection
func (cs *Computations) Append(c *Computation) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cs.data = append(cs.data, c)
}

func (cs *Computations) Find(name string) *Computation {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for _, computation := range cs.data {
		if computation.Description().Name == name {
			return computation
		}
	}

	return nil
}

// TODO: Should the collection support delteion?
// func (cs *Computations) Remove(name string) error {
// }
