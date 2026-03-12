package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var agentFiles = []string{
	"architect.md",
	"researcher.md",
	"ux-designer.md",
	"frontend-engineer.md",
	"backend-engineer.md",
	"devops-sre.md",
	"qa-engineer.md",
	"code-reviewer.md",
}

func TestInit_CreatesAllFiles(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Check maestro.yaml
	assertFileExists(t, filepath.Join(dir, "maestro.yaml"))

	// Check .maestro/ state directory
	assertDirExists(t, filepath.Join(dir, ".maestro"))

	// Check all agent files
	for _, name := range agentFiles {
		assertFileExists(t, filepath.Join(dir, ".opencode", "agent", name))
	}

	// Check .gitignore contains .maestro/
	assertFileContains(t, filepath.Join(dir, ".gitignore"), ".maestro/")
}

func TestInit_SkipsExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// Pre-create maestro.yaml with custom content
	configPath := filepath.Join(dir, "maestro.yaml")
	original := []byte("custom: config\n")
	if err := os.WriteFile(configPath, original, 0o644); err != nil {
		t.Fatalf("writing pre-existing maestro.yaml: %v", err)
	}

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// maestro.yaml should NOT be overwritten
	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("reading maestro.yaml: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("maestro.yaml was overwritten; got %q, want %q", got, original)
	}

	// Other files should still be created
	for _, name := range agentFiles {
		assertFileExists(t, filepath.Join(dir, ".opencode", "agent", name))
	}
}

func TestInit_SkipsWhenAlreadyInitialized(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	// Modify maestro.yaml to detect if it gets overwritten
	configPath := filepath.Join(dir, "maestro.yaml")
	custom := []byte("custom: true\n")
	if err := os.WriteFile(configPath, custom, 0o644); err != nil {
		t.Fatalf("writing custom maestro.yaml: %v", err)
	}

	// Running again should skip entirely
	if err := Init(dir); err != nil {
		t.Fatalf("second Init() error = %v", err)
	}

	// maestro.yaml should keep the custom content (not re-scaffolded)
	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("reading maestro.yaml: %v", err)
	}
	if string(got) != string(custom) {
		t.Errorf("maestro.yaml was modified on re-init; got %q, want %q", got, custom)
	}
}

func TestIsInitialized(t *testing.T) {
	dir := t.TempDir()

	if IsInitialized(dir) {
		t.Fatal("empty dir should not be initialized")
	}

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if !IsInitialized(dir) {
		t.Fatal("dir should be initialized after Init()")
	}
}

func TestInit_AppendsToExistingGitignore(t *testing.T) {
	dir := t.TempDir()

	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatalf("writing .gitignore: %v", err)
	}

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	got, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}

	if !strings.Contains(string(got), "node_modules/") {
		t.Error(".gitignore lost existing content")
	}
	if !strings.Contains(string(got), ".maestro/") {
		t.Error(".gitignore missing .maestro/ entry")
	}
}

func TestInit_GitignoreNotDuplicated(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}
	if err := Init(dir); err != nil {
		t.Fatalf("second Init() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}

	count := strings.Count(string(got), ".maestro/")
	if count != 1 {
		t.Errorf(".maestro/ appears %d times in .gitignore, want 1", count)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("expected file %s to exist: %v", path, err)
		return
	}
	if info.IsDir() {
		t.Errorf("expected %s to be a file, got directory", path)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("expected directory %s to exist: %v", path, err)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory, got file", path)
	}
}

func assertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("reading %s: %v", path, err)
		return
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("file %s does not contain %q", path, substr)
	}
}
