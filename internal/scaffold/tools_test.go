package scaffold

import (
	"testing"
)

func TestValidateTool_Supported(t *testing.T) {
	if !ValidateTool("opencode") {
		t.Errorf("ValidateTool(%q) = false, want true", "opencode")
	}
}

func TestValidateTool_Unsupported(t *testing.T) {
	for _, name := range []string{"cursor", "cloudcode", "", "OPENCODE", "amp"} {
		if ValidateTool(name) {
			t.Errorf("ValidateTool(%q) = true, want false", name)
		}
	}
}

func TestSupportedTools(t *testing.T) {
	got := SupportedTools()

	want := []string{"opencode"}
	if len(got) != len(want) {
		t.Fatalf("SupportedTools() = %v, want %v", got, want)
	}
	if got[0] != want[0] {
		t.Errorf("SupportedTools()[0] = %q, want %q", got[0], want[0])
	}
}
