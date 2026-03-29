package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templates embed.FS

// IsInitialized reports whether the target directory already has a maestro project.
func IsInitialized(targetDir string) bool {
	paths := []string{
		filepath.Join(targetDir, "maestro.yaml"),
		filepath.Join(targetDir, ".maestro"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			return false
		}
	}
	return true
}

// Init scaffolds a new maestro project in the given directory.
// Now only supports "opencode" tool.
func Init(targetDir string, tools ...string) error {
	// Default to opencode only
	if len(tools) == 0 {
		tools = []string{"opencode"}
	}

	// Validate - only opencode is supported
	for _, t := range tools {
		if !ValidateTool(t) {
			return fmt.Errorf("unsupported tool: %q (only 'opencode' is supported)", t)
		}
	}

	if IsInitialized(targetDir) {
		fmt.Println("Project already initialized, skipping.")
		return nil
	}

	// Write maestro.yaml from template (opencode only now)
	if err := writeMaestroYAML(targetDir); err != nil {
		return err
	}

	// Scaffold opencode tool
	if err := scaffoldTool(targetDir, "opencode"); err != nil {
		return fmt.Errorf("scaffolding opencode: %w", err)
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

// Update re-scaffolds agent files from the latest templates.
// This updates the .opencode/ directory with current agent prompts.
func Update(targetDir, tool string) error {
	// Validate tool
	if !ValidateTool(tool) {
		return fmt.Errorf("unsupported tool: %q (only 'opencode' is supported)", tool)
	}

	// Check if project is initialized
	if !IsInitialized(targetDir) {
		return fmt.Errorf("project not initialized. Run 'maestro init' first")
	}

	fmt.Printf("Updating %s agent files...\n", tool)

	// Re-scaffold the tool (force overwrite)
	if err := scaffoldToolForce(targetDir, tool); err != nil {
		return fmt.Errorf("updating %s: %w", tool, err)
	}

	fmt.Println("Update complete.")
	return nil
}

// Remove maestro.yaml, .maestro/, .agents/, .opencode/, etc.
func Drop(targetDir string) error {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("reading target directory: %w", err)
	}

	for _, e := range entries {
		name := e.Name()
		if name == "maestro.yaml" ||
			strings.HasPrefix(name, ".maestro") ||
			strings.HasPrefix(name, ".agents") ||
			strings.HasPrefix(name, ".opencode") {

			// absolute path
			path := filepath.Join(targetDir, name)

			fmt.Printf("Removing %s...\n", path)

			// Remove path (file or directory)
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("removing %s: %w", path, err)
			}
		}
	}

	return nil
}

// writeMaestroYAML reads the template and writes it.
func writeMaestroYAML(targetDir string) error {
	dest := filepath.Join(targetDir, "maestro.yaml")
	if _, err := os.Stat(dest); err == nil {
		fmt.Println("  skip maestro.yaml (already exists)")
		return nil
	}

	data, err := templates.ReadFile("templates/maestro.yaml")
	if err != nil {
		return fmt.Errorf("reading maestro.yaml template: %w", err)
	}

	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("writing maestro.yaml: %w", err)
	}
	fmt.Println("  create maestro.yaml")
	return nil
}

// roles returns the list of role names (derived from templates/roles/).
func roles() ([]string, error) {
	entries, err := templates.ReadDir("templates/roles")
	if err != nil {
		return nil, fmt.Errorf("reading roles: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		names = append(names, name)
	}
	return names, nil
}

// splitFrontmatter splits a markdown file with YAML frontmatter into
// the frontmatter fields (without --- fences) and the body.
func splitFrontmatter(data []byte) (frontmatter, body string) {
	s := string(data)
	if !strings.HasPrefix(s, "---\n") {
		return "", s
	}
	end := strings.Index(s[4:], "\n---")
	if end < 0 {
		return "", s
	}
	return s[4 : 4+end], s[4+end+4:]
}

// scaffoldTool writes agent files for a specific tool backend (skip existing).
func scaffoldTool(targetDir, tool string) error {
	roleNames, err := roles()
	if err != nil {
		return err
	}

	for _, role := range roleNames {
		roleRaw, err := templates.ReadFile("templates/roles/" + role + ".md")
		if err != nil {
			return fmt.Errorf("reading role %s: %w", role, err)
		}

		toolRaw, err := templates.ReadFile("templates/" + tool + "/" + role + ".yaml")
		if err != nil {
			return fmt.Errorf("reading %s wrapper for %s: %w", tool, role, err)
		}

		// Parse role frontmatter (description) and body
		roleFM, roleBody := splitFrontmatter(roleRaw)
		// Parse tool wrapper frontmatter (tool-specific fields)
		toolFM, _ := splitFrontmatter(toolRaw)

		// Merge: role frontmatter + tool frontmatter + body
		var merged strings.Builder
		merged.WriteString("---\n")
		merged.WriteString(strings.TrimSpace(roleFM))
		merged.WriteString("\n")
		if toolFM != "" {
			merged.WriteString(strings.TrimSpace(toolFM))
			merged.WriteString("\n")
		}
		merged.WriteString("---")
		merged.WriteString(roleBody)

		dest := agentOutputPath(targetDir, tool, role)
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("  skip %s (already exists)\n", dest)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", dest, err)
		}

		if err := os.WriteFile(dest, []byte(merged.String()), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}
		fmt.Printf("  create %s\n", dest)
	}

	return nil
}

// scaffoldToolForce writes agent files, overwriting existing ones.
func scaffoldToolForce(targetDir, tool string) error {
	roleNames, err := roles()
	if err != nil {
		return err
	}

	for _, role := range roleNames {
		roleRaw, err := templates.ReadFile("templates/roles/" + role + ".md")
		if err != nil {
			return fmt.Errorf("reading role %s: %w", role, err)
		}

		toolRaw, err := templates.ReadFile("templates/" + tool + "/" + role + ".yaml")
		if err != nil {
			return fmt.Errorf("reading %s wrapper for %s: %w", tool, role, err)
		}

		// Parse role frontmatter (description) and body
		roleFM, roleBody := splitFrontmatter(roleRaw)
		// Parse tool wrapper frontmatter (tool-specific fields)
		toolFM, _ := splitFrontmatter(toolRaw)

		// Merge: role frontmatter + tool frontmatter + body
		var merged strings.Builder
		merged.WriteString("---\n")
		merged.WriteString(strings.TrimSpace(roleFM))
		merged.WriteString("\n")
		if toolFM != "" {
			merged.WriteString(strings.TrimSpace(toolFM))
			merged.WriteString("\n")
		}
		merged.WriteString("---")
		merged.WriteString(roleBody)

		dest := agentOutputPath(targetDir, tool, role)
		exists := false
		if _, err := os.Stat(dest); err == nil {
			exists = true
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", dest, err)
		}

		if err := os.WriteFile(dest, []byte(merged.String()), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}

		if exists {
			fmt.Printf("  update %s\n", dest)
		} else {
			fmt.Printf("  create %s\n", dest)
		}
	}

	return nil
}

// agentOutputPath returns the destination path for an agent file based on the tool.
func agentOutputPath(targetDir, tool, role string) string {
	switch tool {
	case "opencode":
		return filepath.Join(targetDir, ".opencode", "agent", role+".md")
	default:
		return filepath.Join(targetDir, "."+tool, "agent", role+".md")
	}
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

	entry := "\n# Maestro Setup\n.maestro/\n"
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		entry = "\n" + entry
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening .gitignore: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing .gitignore: %v\n", err)
		}
	}()

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("writing to .gitignore: %w", err)
	}

	fmt.Println("  update .gitignore")
	return nil
}
