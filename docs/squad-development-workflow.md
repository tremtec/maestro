# Squad Development Workflow

Maestro orchestrates a squad of AI agents through a structured 4-phase
workflow. It acts as a **Product Owner (PO)** — receiving objectives,
coordinating discovery, synthesizing learnings, and driving delivery.

Maestro never writes code. It only orchestrates.

---

## Squad Roles

The squad is organized into three functional groups aligned with the
workflow phases.

### Discovery Team

| Role            | Responsibility                                                        |
| --------------- | --------------------------------------------------------------------- |
| **Architect**   | Technical feasibility, system design, tradeoffs, component boundaries |
| **Researcher**  | Documentation lookup, dependency analysis, prior art, API references  |
| **UX Designer** | User flows, accessibility, UI patterns, interaction design            |

### Build Team

| Role                  | Responsibility                                           |
| --------------------- | -------------------------------------------------------- |
| **Frontend Engineer** | UI components, client-side logic, state management       |
| **Backend Engineer**  | API routes, database schemas, server logic, integrations |
| **DevOps / SRE**      | Infrastructure, CI/CD pipelines, deployment, monitoring  |

### Quality Team

| Role              | Responsibility                                                |
| ----------------- | ------------------------------------------------------------- |
| **QA Engineer**   | Automated tests, linting, integration checks                  |
| **Code Reviewer** | API validation, standards compliance, hallucination detection |

---

## Workflow Phases

### Phase 1 — Discovery

Maestro receives an objective and fans out to the Discovery team in
parallel. Each agent investigates the problem from their perspective.

```text
Objective
  │
  ├──▶ Architect    → architecture proposal, component diagram
  ├──▶ Researcher   → relevant docs, API references, dependencies
  └──▶ UX Designer  → user flows, UI patterns, accessibility notes
```

**Inputs:** User objective (plain text).

**Outputs:** Each agent produces a structured findings report. These are
collected by Maestro for the next phase.

### Phase 2 — Synthesis

Maestro consolidates all discovery findings into a single summary and
generates a task plan with agent assignments.

1. **Summarize learnings** — merge architecture proposals, research
   findings, and UX recommendations into a coherent brief.
2. **Generate plan** — break the objective into ordered milestones,
   each assigned to a Build team agent.
3. **Present for approval** — the user reviews the plan and can adjust
   scope, priorities, or assignments before proceeding.

**Inputs:** Discovery findings from Phase 1.

**Outputs:** An approved execution plan with milestones and agent
assignments.

### Phase 3 — Build

Maestro delegates tasks to the Build team following a dependency graph
(DAG). Independent tasks run in parallel; dependent tasks wait for
their prerequisites.

```text
Plan
  │
  ├──▶ Backend Engineer    (no deps)
  ├──▶ DevOps / SRE        (no deps)
  └──▶ Frontend Engineer   (depends_on: backend)
```

Each agent works within an isolated session managed by OpenCode. Maestro
monitors progress via the event stream and coordinates handoffs between
agents.

**Inputs:** Approved plan from Phase 2.

**Outputs:** Implemented code changes across the codebase.

### Phase 4 — Quality Gate

The Quality team verifies all changes before delivery.

1. **Automated checks** — QA Engineer runs test and lint commands.
2. **Code review** — Code Reviewer validates API usage against
   documentation, checks for hallucinated or deprecated methods.
3. **Feedback loop** — failures are routed back to the responsible
   Build agent with corrective instructions. The cycle repeats until
   all checks pass or `max_retries` is reached.

```text
Changes
  │
  ├──▶ QA Engineer     → go test, go vet
  └──▶ Code Reviewer   → API validation, standards check
         │
         ├── ✅ Pass → Deliver to user
         └── ❌ Fail → Route back to Build agent → retry
```

**Inputs:** Code changes from Phase 3.

**Outputs:** Verified, tested, deliverable code.

---

## Workflow Configuration

The workflow is defined in `maestro.yaml` at the project root.

```yaml
name: my-project
description: Project description

# Supported tools — agent runtimes to scaffold and use
tools:
  - opencode

# Squad definition — maps roles to agent configurations
agents:
  architect:
    role: architect
    phase: discovery

  researcher:
    role: researcher
    phase: discovery

  ux-designer:
    role: ux-designer
    phase: discovery

  frontend-engineer:
    role: frontend-engineer
    phase: build

  backend-engineer:
    role: backend-engineer
    phase: build

  devops-sre:
    role: devops-sre
    phase: build

  qa-engineer:
    role: qa-engineer
    phase: quality

  code-reviewer:
    role: code-reviewer
    phase: quality

# Workflow phases
workflow:
  discovery:
    parallel: true

  synthesis:
    approval: true # pause for user approval before build

  build:
    parallel: true
    tasks:
      - agent: backend-engineer
        depends_on: []
      - agent: devops-sre
        depends_on: []
      - agent: frontend-engineer
        depends_on: [backend-engineer]

  quality_gate:
    commands:
      - go test ./...
      - go vet ./...
    on_failure: fix_and_retry
    max_retries: 3
```

### Agent Definition

Agent definitions are composed from shared role files and tool-specific
wrappers. The role content (instructions, responsibilities, output
format) is shared across all tools — only the frontmatter and file
structure differ.

**OpenCode** — agents are markdown files in `.opencode/agent/`:

```markdown
---
description: Brief description of the agent's expertise
mode: subagent
tools:
  write: true
  edit: true
---

Role instructions...
```

---

## Execution Model

Maestro drives the workflow through a **Runtime** interface that
abstracts the underlying AI tool. OpenCode provides the execution
backend:

```text
┌──────────────┐     Runtime interface    ┌──────────────────┐
│  Maestro CLI │ ◀──────────────────────▶ │  OpenCode Runtime│
│  (Go)        │   RunAgent(prompt) →     │  (HTTP API)      │
└──────────────┘   ← Result              └──────────────────┘
```

### OpenCode Runtime

Uses the OpenCode HTTP API:

1. **Start server** — `opencode serve` exposes the OpenAPI endpoint.
2. **Create sessions** — each agent task runs in an isolated session.
3. **Send prompts** — `POST /session/:id/message` with the task
   description and context from previous phases.
4. **Monitor events** — `GET /event` (SSE) streams real-time progress.
5. **Collect results** — `GET /session/:id/message` retrieves agent
   output for synthesis and review.

Or use the CLI directly:

```bash
opencode "Your prompt here"
```

---

## State Management

Maestro tracks all objective state in the `.maestro/` directory
(git-ignored). Each objective gets a timestamped folder with a phase
report for each completed step.

```
.maestro/
└── 2026-03-12-user-management-api/
    ├── discovery-report.md
    ├── synthesis-report.md
    ├── build-report.md
    └── quality-gate-report.md
```

**Objective lifecycle:**

```
initiated → discovery → synthesis → [human approval] → build → quality-gate → done
```

---

## Key Principles

- **Maestro never writes code.** It plans, delegates, and verifies.
- **Discovery before building.** Every objective starts with research
  and architecture review.
- **Human in the loop.** The user approves the plan before build
  begins.
- **Fail fast, retry smart.** Quality gate failures route back to the
  responsible agent with specific corrective instructions.
- **Parallel by default.** Independent tasks run concurrently to
  minimize wall-clock time.
- **DAG-based ordering.** Dependencies between tasks are explicit,
  ensuring correct sequencing.
