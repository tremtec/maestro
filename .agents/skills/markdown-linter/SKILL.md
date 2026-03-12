---
name: markdown-linter
description: "Formats and lints markdown (.md) files. Use whenever modifying or creating documentation, agent definitions, or skill files."
---

# Markdown Linter Skill

Ensures all markdown documentation meets structural and formatting standards.

## Triggers

Run this skill whenever you:

- Create or update `.md` files (like `README.md`, docs, or agent definitions).
- Create or update agent files in `.opencode/agent/`.
- Create or update skill files in `.agents/skills/`.

## Standard Workflow

After creating or modifying a markdown file:

1. **Verify structure**:
   - Headings use proper hierarchy (no skipping levels).
   - Fenced code blocks specify a language (e.g., ` ```bash `, ` ```go `, ` ```yaml `).
   - Blank lines surround headings, code blocks, and lists.
   - No trailing whitespace.

2. **Check links**:
   - Internal links point to existing files.
   - Relative paths are correct from the file's location.

3. **Verify tables**:
   - Column alignment is consistent.
   - Header separators use correct syntax (`| --- |`).

## Common Issues

- **Fenced code blocks**: Always specify a language.
- **Blank lines**: Ensure empty lines before and after code blocks and headings.
- **Consistent list markers**: Use `-` for unordered lists throughout.
