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
// tools specifies which agent runtimes to scaffold (e.g., ["opencode", "amp"]).
func Init(targetDir string, tools ...string) error {
	if len(tools) == 0 {
		tools = []string{"opencode"}
	}

	for _, t := range tools {
		if !ValidateTool(t) {
			return fmt.Errorf("unsupported tool: %q (supported: %v)", t, SupportedTools())
		}
	}

	if IsInitialized(targetDir) {
		fmt.Println("Project already initialized, skipping.")
		return nil
	}

	// Write maestro.yaml from template (with tools list injected)
	if err := writeMaestroYAML(targetDir, tools); err != nil {
		return err
	}

	// Scaffold each requested tool
	for _, tool := range tools {
		if err := scaffoldTool(targetDir, tool); err != nil {
			return fmt.Errorf("scaffolding %s: %w", tool, err)
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

// writeMaestroYAML reads the template and replaces the tools list.
func writeMaestroYAML(targetDir string, tools []string) error {
	dest := filepath.Join(targetDir, "maestro.yaml")
	if _, err := os.Stat(dest); err == nil {
		fmt.Println("  skip maestro.yaml (already exists)")
		return nil
	}

	data, err := templates.ReadFile("templates/maestro.yaml")
	if err != nil {
		return fmt.Errorf("reading maestro.yaml template: %w", err)
	}

	// Replace the default tools list with the requested ones
	var toolLines strings.Builder
	for _, t := range tools {
		toolLines.WriteString("  - " + t + "\n")
	}
	content := strings.Replace(
		string(data),
		"tools:\n  - opencode\n",
		"tools:\n"+toolLines.String(),
		1,
	)

	if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
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

// scaffoldTool writes agent files for a specific tool backend.
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

// agentOutputPath returns the destination path for an agent file based on the tool.
func agentOutputPath(targetDir, tool, role string) string {
	switch tool {
	case "opencode":
		return filepath.Join(targetDir, ".opencode", "agent", role+".md")
	case "amp":
		if role == "maestro" {
			return filepath.Join(targetDir, "AGENTS.md")
		}
		return filepath.Join(targetDir, ".agents", "skills", role, "SKILL.md")
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
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("writing to .gitignore: %w", err)
	}

	fmt.Println("  update .gitignore")
	return nil
}
