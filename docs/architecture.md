# ShellForge Architecture

## Overview

ShellForge is a single Go binary (~7.5MB) that provides governed local AI agent execution. Its core value is **governance** — when drivers like Crush or Claude Code are installed, they provide the agentic loop; ShellForge wraps them with AgentGuard policy enforcement.

## 7-Layer Stack

```
┌─────────────────────────────────────────────┐
│  Layer 7: DefenseClaw (Supply Chain)        │  Cisco AI BoM Scanner
├─────────────────────────────────────────────┤
│  Layer 6: OpenShell (Kernel Sandbox)        │  NVIDIA Landlock/Seccomp
├─────────────────────────────────────────────┤
│  Layer 5: AgentGuard (Governance Kernel)    │  Policy enforcement
├─────────────────────────────────────────────┤
│  Layer 4: Dagu (Orchestration)              │  YAML DAG workflows + web UI
├─────────────────────────────────────────────┤
│  Layer 3: Crush (Execution Engine)          │  Go-native AI coding agent
├─────────────────────────────────────────────┤
│  Layer 2: RTK (Token Compression)           │  Auto-compress I/O
├─────────────────────────────────────────────┤
│  Layer 1: Ollama (Local LLM)                │  Metal GPU on Mac
└─────────────────────────────────────────────┘
```

## Go Project Layout

```
cmd/shellforge/
├── main.go         # CLI entry point (subcommands: run, agent, qa, report, serve, swarm, scan, setup, status)
└── status.go       # Ecosystem health check

internal/
├── action/         # Core types: ActionProposal, ActionResult, RiskLevel, Scope
├── agent/          # Native fallback agentic loop
├── correction/     # Denial tracking, escalation levels, corrective feedback for LLM
├── engine/         # Pluggable engine interface (legacy OpenCode/DeepAgents adapters)
├── governance/     # agentguard.yaml parser + policy engine
├── intent/         # Format-agnostic intent parser — extracts actions from ANY LLM output
├── logger/         # Structured JSON logging
├── normalizer/     # Tool call → Canonical Action Representation (CAR) converter
├── ollama/         # Ollama HTTP client (chat, generate)
├── orchestrator/   # Run lifecycle state machine (IDLE → PLANNING → WORKING → EVALUATING)
├── scheduler/      # Priority-aware inference queue with semaphore concurrency control
└── integration/    # RTK, OpenShell, DefenseClaw, TurboQuant, AgentGuard adapters
```

## Driver Architecture

ShellForge governs any CLI agent driver via AgentGuard hooks. The driver handles the agent loop; ShellForge ensures governance is active and spawns it as a subprocess.

Supported drivers for `shellforge run <driver>`:
- **crush** — Go-native AI coding agent (TUI + headless)
- **claude** — Anthropic Claude Code CLI
- **copilot** — GitHub Copilot CLI
- **codex** — OpenAI Codex CLI
- **gemini** — Google Gemini CLI

Fallback: built-in **native** agentic loop (Ollama + tool calling, no external driver required).

## Governance Flow

```
User Request → shellforge run <driver> / shellforge agent
  → Driver subprocess (or native loop)
    → Tool Call → Governance Check (agentguard.yaml)
      → ALLOW → Execute Tool → Return Result
      → DENY  → Correction Engine → Structured feedback → LLM self-corrects
```

## Governed Multi-Agent Pipeline (Phase 7)

```
Intent Parser  →  Normalizer  →  Governance Check
(any LLM fmt)     (→ CAR)        (allow / deny)
                                      │
                               Correction Engine
                               (escalate + coach)
                                      │
                            Orchestrator State Machine
                            IDLE → PLANNING → WORKING
                                → EVALUATING → COMPLETE
                                      │
                               Scheduler Queue
                               (priority + semaphore)
```

## Data Flow

1. User invokes `shellforge agent` (or run, qa, report)
2. CLI loads `agentguard.yaml` governance policy
3. Intent parser normalizes LLM output to unified Action structs
4. Normalizer converts tool calls to Canonical Action Representations (CARs)
5. Each CAR passes through governance check
6. Allowed actions execute (shell commands wrapped by RTK + OpenShell sandbox)
7. Denied actions trigger the correction engine — structured feedback coaches LLM toward compliant alternatives
8. Orchestrator state machine tracks run phase; scheduler queue manages concurrency
9. Results compressed by RTK, fed back to LLM
10. Loop continues until task complete or budget exhausted

## macOS (Apple Silicon) Support

All layers run on Mac M4:
- Ollama uses Metal for GPU acceleration
- RTK, AgentGuard, Crush are native arm64 binaries
- TurboQuant runs via PyTorch (MPS backend) for KV cache optimization
- OpenShell runs inside Docker/Colima (Linux VM for Landlock/Seccomp)
- DefenseClaw installs via pip or source build
