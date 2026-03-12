---
name: update-agent-md
description: "Updates AGENTS.md with new technical guidelines and architectural patterns. Use whenever establishing new conventions for the Maestro Go CLI project."
---

# Update AGENTS.md Skill

Keeps AGENTS.md in sync with evolving technical standards and architectural decisions.

## Purpose

AGENTS.md is the single source of truth for project conventions. When you establish new patterns, this skill ensures they're documented.

## When to Update AGENTS.md

Update when:

- **Establishing new Go patterns** (e.g., error handling, package layout)
- **Adding new CLI commands or subcommands**
- **Defining conventions for agent definitions** (`.opencode/agent/` format)
- **Changing workflow phases or state management**
- **Adding new Cobra command patterns**

DO NOT add:

- Feature-specific implementation details
- Temporary decisions or experiments
- Objective-specific state (belongs in `.maestro/`)

## Update Workflow

1. **Identify the change**:
   - What new pattern or convention exists?
   - Is it project-wide or specific to one command?
   - Will future development benefit from knowing this?

2. **Find the right section** in AGENTS.md:
   - **Build & Test Commands**: `go build`, `go test`, `make` targets
   - **Architecture & Structure**: Package layout, module boundaries
   - **Code Style & Conventions**: Naming, error handling, imports
   - **Workflow & Agents**: Agent definition format, phase conventions
   - Create new section if needed (but be selective)

3. **Write concise guidelines**:
   - Start with principle, not implementation
   - Use bullet points
   - Example: "Return errors, don't panic — use `fmt.Errorf` with `%w` for wrapping"

4. **Keep it general**:
   - Focus on the pattern, not the specific feature
   - Other contributors should understand the principle

## Checklist

Before updating AGENTS.md:

- [ ] Is this a reusable pattern (not feature-specific)?
- [ ] Will future development benefit from knowing this?
- [ ] Is it general enough for the project scope?
- [ ] Does it fit an existing section or need a new one?
