package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templates embed.FS

// IsInitialized reports whether the target directory already has a maestro project.
// A project is considered initialized when maestro.yaml, .maestro/, and
// .opencode/agent/ all exist.
func IsInitialized(targetDir string) bool {
	paths := []string{
		filepath.Join(targetDir, "maestro.yaml"),
		filepath.Join(targetDir, ".maestro"),
		filepath.Join(targetDir, ".opencode", "agent"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			return false
		}
	}
	return true
}

// Init scaffolds a new maestro project in the given directory.
// It creates the agent markdown files, maestro.yaml, and the .maestro state directory.
// If the project is already initialized, it returns nil immediately.
func Init(targetDir string) error {
	if IsInitialized(targetDir) {
		fmt.Println("Project already initialized, skipping.")
		return nil
	}

	// Create .opencode/agent/ files and maestro.yaml from embedded templates
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Strip the "templates/" prefix to get the relative output path
		rel := strings.TrimPrefix(path, "templates/")
		if rel == "" {
			return nil
		}

		dest := filepath.Join(targetDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}

		// Skip if file already exists
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("  skip %s (already exists)\n", rel)
			return nil
		}

		data, err := templates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", rel, err)
		}

		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", rel, err)
		}

		fmt.Printf("  create %s\n", rel)
		return nil
	})
	if err != nil {
		return fmt.Errorf("scaffolding templates: %w", err)
	}

	// Rename templates output: agent/ -> .opencode/agent/
	agentSrc := filepath.Join(targetDir, "agent")
	if _, err := os.Stat(agentSrc); err == nil {
		agentDest := filepath.Join(targetDir, ".opencode", "agent")
		if err := os.MkdirAll(filepath.Dir(agentDest), 0o755); err != nil {
			return fmt.Errorf("creating .opencode/: %w", err)
		}
		if _, err := os.Stat(agentDest); os.IsNotExist(err) {
			if err := os.Rename(agentSrc, agentDest); err != nil {
				return fmt.Errorf("moving agent/ to .opencode/agent/: %w", err)
			}
		} else {
			// Agent dir already exists, remove the temp copy
			os.RemoveAll(agentSrc)
		}
	}

	// Create .maestro/ state directory
	maestroDir := filepath.Join(targetDir, ".maestro")
	if err := os.MkdirAll(maestroDir, 0o755); err != nil {
		return fmt.Errorf("creating .maestro/: %w", err)
	}
	fmt.Println("  create .maestro/")

	// Ensure .maestro/ is in .gitignore
	if err := ensureGitignore(targetDir); err != nil {
		return err
	}

	return nil
}

// ensureGitignore adds .maestro/ to .gitignore if not already present.
func ensureGitignore(targetDir string) error {
	gitignorePath := filepath.Join(targetDir, ".gitignore")

	var content []byte
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = data
	}

	if strings.Contains(string(content), ".maestro/") {
		return nil
	}

	entry := ".maestro/\n"
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		entry = "\n" + entry
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening .gitignore: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("writing to .gitignore: %w", err)
	}

	fmt.Println("  update .gitignore")
	return nil
}
