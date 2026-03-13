package scaffold

import (
	"sort"
	"testing"
)

func TestValidateTool_Supported(t *testing.T) {
	for _, name := range []string{"opencode", "amp"} {
		if !ValidateTool(name) {
			t.Errorf("ValidateTool(%q) = false, want true", name)
		}
	}
}

func TestValidateTool_Unsupported(t *testing.T) {
	for _, name := range []string{"cursor", "cloudcode", "", "OPENCODE"} {
		if ValidateTool(name) {
			t.Errorf("ValidateTool(%q) = true, want false", name)
		}
	}
}

func TestSupportedTools(t *testing.T) {
	got := SupportedTools()
	sort.Strings(got)

	want := []string{"amp", "opencode"}
	if len(got) != len(want) {
		t.Fatalf("SupportedTools() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SupportedTools()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
