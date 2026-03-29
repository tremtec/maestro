# 🎵 Maestro

A Unix-compatible CLI tool that orchestrates squad-based development
using [OpenCode](https://opencode.ai) as the agent runtime.

Maestro is a primary agent that leads your agentic opera — coordinating
specialist sub-agents through a 4-phase workflow (Discovery → Synthesis
→ Build → Quality Gate) to reach an objective with minimal human
interaction.

## Features

- **Squad orchestration** — manages specialist sub-agents (Architect,
  Frontend, Backend, DevOps, QA) through structured phases.
- **Long-running tasks** — designed for complex, multi-phase objectives.
- **Minimal human interaction** — only plan approval is required.
- **Markdown as database** — all state, reports, and artifacts are
  plain markdown files in `.maestro/`.
- **Unix principles** — accepts stdin, returns stdout, composes with
  pipes.
- **CLI first** — built with Cobra, uses OpenCode as the agent runtime.

## Getting Started

```bash
go install github.com/tremtec/maestro@latest
maestro init                          # scaffold for opencode
maestro run "Build a REST API for user management"
```

## Commands

| Command               | Description                                     |
| --------------------- | ----------------------------------------------- |
| `maestro init`        | Set up squad, config, and `.maestro/` state dir |
| `maestro run`         | Run a prompt through the workflow              |
| `maestro update`      | Update agent definitions to latest templates    |
| `maestro upgrade`     | Upgrade maestro CLI to latest version          |
| `maestro completions` | Generate shell completions                      |

## Documentation

- [Maestro Overview](docs/maestro.md)
- [Squad Development Workflow](docs/squad-development-workflow.md)
