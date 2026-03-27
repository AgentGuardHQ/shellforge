# ShellForge Architecture

## Overview

ShellForge is a single Go binary (~7.5MB) that provides governed local AI agent execution. Its core value is **governance** — when frameworks like OpenCode or DeepAgents are installed, they provide the agentic loop; ShellForge wraps them with AgentGuard policy enforcement.

## 8-Layer Stack

```
┌─────────────────────────────────────────────┐
│  Layer 8: OpenShell (Kernel Sandbox)        │  NVIDIA Landlock/Seccomp
├─────────────────────────────────────────────┤
│  Layer 7: DefenseClaw (Supply Chain)        │  Cisco AI BoM Scanner
├─────────────────────────────────────────────┤
│  Layer 6: DeepAgents (Multi-Agent)          │  LangChain orchestration
├─────────────────────────────────────────────┤
│  Layer 5: OpenCode (AI Coding)              │  Go CLI, native tools
├─────────────────────────────────────────────┤
│  Layer 4: AgentGuard (Governance Kernel)    │  Policy enforcement
├─────────────────────────────────────────────┤
│  Layer 3: TurboQuant (Quantization)         │  KV cache optimization
├─────────────────────────────────────────────┤
│  Layer 2: RTK (Token Compression)           │  Auto-compress I/O
├─────────────────────────────────────────────┤
│  Layer 1: Ollama (Local LLM)                │  Metal GPU on Mac
└─────────────────────────────────────────────┘
```

## Go Project Layout

```
cmd/shellforge/
├── main.go         # CLI entry point (cobra-style subcommands)
└── status.go       # Ecosystem health check

internal/
├── governance/     # agentguard.yaml parser + policy engine
├── ollama/         # Ollama HTTP client (chat, generate)
├── agent/          # Native fallback agentic loop
├── tools/          # 5 tool implementations + RTK wrapper
├── engine/         # Pluggable engine interface (OpenCode, DeepAgents)
├── logger/         # Structured JSON logging
└── integration/    # RTK, OpenShell, DefenseClaw, TurboQuant, AgentGuard
```

## Engine Architecture

ShellForge uses a pluggable engine system:

1. **OpenCode** (preferred) — subprocess, `--non-interactive` mode, governance-wrapped
2. **DeepAgents** — subprocess, Node.js/Python SDK, governance-wrapped
3. **Native** (fallback) — built-in multi-turn loop with Ollama + tool calling

The engine selection is automatic based on what's installed.

## Governance Flow

```
User Request → Engine (OpenCode/DeepAgents/Native)
  → Tool Call → Governance Check (agentguard.yaml)
    → ALLOW → Execute Tool → Return Result
    → DENY  → Log Violation → Block Execution
```

## Data Flow

1. User invokes `./shellforge qa` (or agent, report, scan)
2. CLI loads `agentguard.yaml` governance policy
3. Detects available engine (OpenCode > DeepAgents > Native)
4. Engine sends prompt to Ollama (via RTK for token compression)
5. LLM responds with tool calls
6. Each tool call passes through governance check
7. Allowed tools execute (shell commands wrapped by RTK + OpenShell sandbox)
8. Results compressed by RTK, fed back to LLM
9. Loop continues until task complete or budget exhausted

## macOS (Apple Silicon) Support

All 8 layers run on Mac M4:
- Ollama uses Metal for GPU acceleration
- RTK, AgentGuard, OpenCode are native arm64 binaries
- TurboQuant runs via PyTorch (MPS backend)
- OpenShell runs inside Docker/Colima (Linux VM for Landlock)
- DefenseClaw installs via pip or source build

## Release Pipeline

ShellForge uses [goreleaser](https://goreleaser.com) for reproducible cross-platform builds.

```
.goreleaser.yaml
├── Builds: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
├── Archives: tar.gz with checksums
└── Homebrew: AgentGuardHQ/homebrew-tap (Formula/shellforge.rb)
```

**To cut a release:**
```bash
git tag v0.2.0
git push origin v0.2.0
# GitHub Actions runs goreleaser, publishes to GitHub Releases + Homebrew tap
```

**ldflags:** `-s -w -X main.version={{.Version}}` — strips debug symbols, injects version.
