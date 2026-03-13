package runtime

import (
	"context"
	"os/exec"
)

type opencodeRuntime struct{}

func init() {
	Register(&opencodeRuntime{})
}

func (o *opencodeRuntime) Name() string { return "opencode" }

func (o *opencodeRuntime) RunAgent(ctx context.Context, agent Agent, prompt string) (Result, error) {
	// TODO: implement via opencode HTTP API
	return Result{Agent: agent.Name}, nil
}

func (o *opencodeRuntime) IsAvailable() bool {
	_, err := exec.LookPath("opencode")
	return err == nil
}
