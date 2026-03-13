package runtime

import (
	"context"
	"os/exec"
)

type ampRuntime struct{}

func init() {
	Register(&ampRuntime{})
}

func (a *ampRuntime) Name() string { return "amp" }

func (a *ampRuntime) RunAgent(ctx context.Context, agent Agent, prompt string) (Result, error) {
	// TODO: implement via amp -x CLI
	return Result{Agent: agent.Name}, nil
}

func (a *ampRuntime) IsAvailable() bool {
	_, err := exec.LookPath("amp")
	return err == nil
}
