package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var roleFiles = []string{
	"architect",
	"researcher",
	"ux-designer",
	"frontend-engineer",
	"backend-engineer",
	"devops-sre",
	"qa-engineer",
	"code-reviewer",
}

func TestInit_OpenCode_CreatesAllFiles(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	assertFileExists(t, filepath.Join(dir, "maestro.yaml"))
	assertDirExists(t, filepath.Join(dir, ".maestro"))

	// Check all OpenCode agent files (roles + maestro)
	for _, name := range roleFiles {
		assertFileExists(t, filepath.Join(dir, ".opencode", "agent", name+".md"))
	}
	assertFileExists(t, filepath.Join(dir, ".opencode", "agent", "maestro.md"))

	assertFileContains(t, filepath.Join(dir, ".gitignore"), ".maestro/")

	// maestro.yaml should list opencode as default tool
	assertFileContains(t, filepath.Join(dir, "maestro.yaml"), "- opencode")
}

func TestInit_UnsupportedTool(t *testing.T) {
	dir := t.TempDir()

	err := Init(dir, "amp")
	if err == nil {
		t.Fatal("expected error for unsupported tool")
	}
	if !strings.Contains(err.Error(), "unsupported tool") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInit_SkipsExistingFiles(t *testing.T) {
	dir := t.TempDir()

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
	for _, name := range roleFiles {
		assertFileExists(t, filepath.Join(dir, ".opencode", "agent", name+".md"))
	}
}

func TestInit_SkipsWhenAlreadyInitialized(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	configPath := filepath.Join(dir, "maestro.yaml")
	custom := []byte("custom: true\n")
	if err := os.WriteFile(configPath, custom, 0o644); err != nil {
		t.Fatalf("writing custom maestro.yaml: %v", err)
	}

	if err := Init(dir); err != nil {
		t.Fatalf("second Init() error = %v", err)
	}

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

func TestInit_AgentFilesHaveFrontmatterAndRole(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// OpenCode files should have OpenCode-specific frontmatter
	ocArchitect := filepath.Join(dir, ".opencode", "agent", "architect.md")
	assertFileContains(t, ocArchitect, "mode: subagent")
	assertFileContains(t, ocArchitect, "# Architect")
}

func TestUpdate_ReScaffoldsFiles(t *testing.T) {
	dir := t.TempDir()

	// First initialize
	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Then update
	if err := Update(dir, "opencode"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Files should still exist and have content
	for _, name := range roleFiles {
		assertFileExists(t, filepath.Join(dir, ".opencode", "agent", name+".md"))
	}
	assertFileExists(t, filepath.Join(dir, ".opencode", "agent", "maestro.md"))
}

func TestUpdate_NotInitialized(t *testing.T) {
	dir := t.TempDir()

	err := Update(dir, "opencode")
	if err == nil {
		t.Fatal("expected error when not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("unexpected error: %v", err)
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
