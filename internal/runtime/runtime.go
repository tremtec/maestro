package runtime

import (
	"context"
	"fmt"
)

// Agent represents an agent configuration from maestro.yaml.
type Agent struct {
	Name  string
	Role  string
	Phase string
}

// Result holds the output from running an agent task.
type Result struct {
	Agent  string
	Output string
	Err    error
}

// Runtime defines the interface for an AI agentic tool backend.
type Runtime interface {
	// Name returns the tool identifier (e.g., "opencode", "amp").
	Name() string

	// RunAgent sends a prompt to an agent and returns the result.
	RunAgent(ctx context.Context, agent Agent, prompt string) (Result, error)

	// IsAvailable reports whether the tool CLI is installed and accessible.
	IsAvailable() bool
}

var registry = map[string]Runtime{}

// Register adds a runtime to the global registry.
func Register(r Runtime) {
	registry[r.Name()] = r
}

// Get returns a registered runtime by name.
func Get(name string) (Runtime, error) {
	r, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown runtime: %q (available: %v)", name, Available())
	}
	return r, nil
}

// Available returns the names of all registered runtimes.
func Available() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
