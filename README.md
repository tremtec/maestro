# 🎵 Maestro

A Unix-compatible CLI tool that orchestrates squad-based development
using AI agentic tools like [OpenCode](https://opencode.ai) and
[Amp](https://ampcode.com).

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
- **CLI first** — built with Cobra, supports OpenCode and Amp as agent
  runtimes.

## Getting Started

```bash
go install github.com/tremtec/maestro@latest
maestro init                          # defaults to opencode
maestro init --tool opencode,amp      # scaffold for both tools
maestro "Build a REST API for user management"
```

## Commands

| Command               | Description                                     |
| --------------------- | ----------------------------------------------- |
| `maestro init`        | Set up squad, config, and `.maestro/` state dir |
| `maestro <prompt>`    | Run an objective through the 4-phase workflow   |
| `maestro completions` | Generate shell completions                      |

## Documentation

- [Maestro Overview](docs/maestro.md)
- [Squad Development Workflow](docs/squad-development-workflow.md)
