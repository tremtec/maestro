---
name: go-quality-gate
description: "Runs Go build, test, and vet checks to ensure code quality. Use whenever creating or modifying Go source files."
---

# Go Quality Gate Skill

Ensures all Go code meets project quality standards before committing.

## Triggers

Run this skill whenever you:

- Create new Go source files
- Modify existing Go code
- Add new dependencies
- Refactor packages or interfaces

## Standard Workflow

After creating or modifying Go code:

1. **Build**:

   ```bash
   make build
   ```

2. **Run tests**:

   ```bash
   make test
   ```

3. **Run vet** (catches suspicious constructs):

   ```bash
   go vet ./...
   ```

4. **Check results**:
   - ✅ All pass: Code is ready.
   - ❌ Failures: Fix issues before proceeding.

## Common Issues

- **Unused imports**: Remove them or use `_` for side-effect imports.
- **Unused variables**: Remove or use them.
- **Error handling**: Always check returned errors; don't use `_` to ignore them silently.
- **Printf format mismatches**: `go vet` catches format string issues.

## Integration

Always run this skill:

- After implementing new commands or packages
- Before committing code
- After adding or updating dependencies in `go.mod`
