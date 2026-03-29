# Maestro

Maestro is a CLI tool that orchestrates squad-based development using
[OpenCode](https://opencode.ai) as the agent runtime. It acts as a
**primary agent** that leads your agentic opera — coordinating specialist
sub-agents through a structured workflow to reach an objective.

Built for long-running tasks with minimal human interaction. Unix
compatible: accepts stdin as input, returns output to stdout.

## What Maestro Does

- **Orchestrates multi-agent workflows** — coordinates Discovery,
  Synthesis, Build, and Quality Gate phases automatically.
- **Manages specialist sub-agents** — delegates work to Architect,
  Researcher, Frontend/Backend Engineers, QA, and others via OpenCode.
- **Minimizes human interaction** — the only required human touchpoint
  is reviewing and approving the plan before the build phase begins.
- **Tracks state as markdown** — all objectives, phase reports, and
  artifacts live in `.maestro/` as plain files.

See [Squad Development Workflow](squad-development-workflow.md) for the
full 4-phase process.

---

## Getting Started

### Install

```bash
go install github.com/tremtec/maestro@latest
```

### Initialize a project

```bash
maestro init                          # scaffold for opencode
```

This will:

1. Set up a squad of sub-agents for OpenCode: `.opencode/agent/*.md`
2. Create a `maestro.yaml` configuration file.
3. Create the `.maestro/` state directory and add it to `.gitignore`.

### Run an objective

```bash
maestro run "Build a REST API for user management"
```

Maestro starts the 4-phase workflow:

1. **Discovery** — sub-agents research the problem in parallel.
2. **Synthesis** — Maestro consolidates findings and presents a plan.
   Waits for human approval.
3. **Build** — delegates tasks to build agents following the dependency
   graph.
4. **Quality Gate** — QA and code review verify the output. Failures
   loop back to the responsible agent.

### Update agent definitions

```bash
maestro update
```

Updates the `.opencode/agent/*.md` files to the latest templates.

### Upgrade maestro

```bash
maestro upgrade
```

Downloads and installs the latest maestro CLI version.

### Shell completions

```bash
maestro completions
```

Generates shell completions (powered by Cobra).

---

## CLI Commands

| Command               | Description                                      |
| --------------------- | ------------------------------------------------ |
| `maestro init`        | Initialize project (squad, config, state folder) |
| `maestro run`         | Start an objective through the workflow          |
| `maestro update`      | Update agent definitions to latest templates    |
| `maestro upgrade`     | Upgrade maestro CLI to latest version            |
| `maestro completions` | Generate shell completions                       |

---

## State Management

Maestro tracks all state in the `.maestro/` directory (git-ignored).

### Directory structure

```
.maestro/
└── 2026-03-12-user-management-api/
    ├── discovery-report.md
    ├── synthesis-report.md
    ├── build-report.md
    └── quality-gate-report.md
```

Each objective gets a folder named `(date)-(objective-title)/`. Each
phase produces a report markdown file:

| File                     | Phase     | Contents                                 |
| ------------------------ | --------- | ---------------------------------------- |
| `discovery-report.md`    | Discovery | Findings from Architect, Researcher, UX  |
| `synthesis-report.md`    | Synthesis | Consolidated plan with milestones        |
| `build-report.md`        | Build     | Implementation status per agent          |
| `quality-gate-report.md` | Quality   | Test results, review feedback, pass/fail |

### Objective lifecycle

```
initiated → discovery → synthesis → [human approval] → build → quality-gate → done
```

The `.maestro/` folder also tracks:

- List of in-progress objectives.
- Each objective's current phase.
- Assets generated per phase.
- Human review status and decisions.

---

## Configuration

Workflow configuration lives in `maestro.yaml` at the project root. See
[Squad Development Workflow](squad-development-workflow.md) for the full
schema.

---

## Design Principles

- **CLI first** — uses OpenCode as the agent runtime, runs from your
  terminal.
- **Unix compatible** — accepts stdin, returns stdout, composes with
  pipes.
- **Minimal human interaction** — only the plan approval step requires
  input.
- **Markdown as database** — all state is plain text, inspectable with
  standard tools.
- **Long-running tasks** — designed for complex, multi-phase objectives
  that take time.
