package scaffold

// supportedTools lists all known tool backends.
var supportedTools = map[string]bool{
	"opencode": true,
	"amp":      true,
}

// ValidateTool reports whether a tool name is supported.
func ValidateTool(name string) bool {
	return supportedTools[name]
}

// SupportedTools returns the list of valid tool names.
func SupportedTools() []string {
	out := make([]string, 0, len(supportedTools))
	for name := range supportedTools {
		out = append(out, name)
	}
	return out
}
