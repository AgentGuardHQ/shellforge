<div align="center">

# 🔥 ShellForge

**Governed local AI agents — a single Go binary, eight integrations, zero cloud.**

[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![GitHub Pages](https://img.shields.io/badge/🌐_Live_Site-agentguardhq.github.io/shellforge-ff6b2b?style=for-the-badge)](https://agentguardhq.github.io/shellforge)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![AgentGuard](https://img.shields.io/badge/Governed_by-AgentGuard-green?style=for-the-badge)](https://github.com/AgentGuardHQ/agentguard)

*Run autonomous AI agents on your machine with policy enforcement on every tool call. No cloud. No API keys. No data leaves your laptop.*

[🌐 Website](https://agentguardhq.github.io/shellforge) · [📖 Docs](docs/architecture.md) · [🗺️ Roadmap](docs/roadmap.md) · [🛡️ AgentGuard](https://github.com/AgentGuardHQ/agentguard)

<img src="https://github.com/user-attachments/assets/a94a8a5e-dfeb-4771-a6ab-465d3c2f01f0" alt="ShellForge — Local Governed Agent Swarm" width="700">

</div>

---

## Quick Start

### Install via Homebrew (macOS / Linux)

```bash
brew tap AgentGuardHQ/tap
brew install shellforge

shellforge setup                 # pulls Ollama + model + verifies stack
shellforge status                # verify 8/8 integrations ✓
shellforge qa                    # run the QA agent
```

### Install from source

```bash
git clone https://github.com/AgentGuardHQ/shellforge.git
cd shellforge

bash scripts/setup.sh --all     # install full 8-layer ecosystem
go build -o shellforge ./cmd/shellforge/
./shellforge status
```

**Requirements:** macOS (Apple Silicon/Intel) or Linux · ~1.3 GB RAM (1.7B model)

---

## What Is ShellForge?

ShellForge is a **governed AI agent runtime** — a single Go binary that orchestrates local LLM inference through [Ollama](https://ollama.com) and wraps every tool call with [AgentGuard](https://github.com/AgentGuardHQ/agentguard) policy enforcement.

**The core insight:** ShellForge's value is **governance**, not the agent loop. When [OpenCode](https://github.com/opencode-ai/opencode) or [DeepAgents](https://github.com/langchain-ai/deepagents) are installed, they provide the agentic loop and tools — ShellForge wraps them with AgentGuard policy enforcement. The native agent loop (`internal/agent/loop.go`) is a fallback for when no external framework is available.

---

## The 8-Layer Ecosystem

Eight open-source integrations. One governed runtime.

| # | Layer | Project | What It Does |
|---|-------|---------|--------------|
| 1 | 🦙 **Infer** | [Ollama](https://ollama.com) | Local LLM inference (Metal GPU on Mac) |
| 2 | ⚡ **Optimize** | [RTK](https://github.com/rtk-ai/rtk) | Token compression — auto-wraps shell output (70–90% reduction) |
| 3 | 🧠 **Quantize** | [TurboQuant](https://github.com/google-research/turboquant) (Google) | KV cache optimization — run 14B models on 8 GB Macs |
| 4 | 🛡️ **Govern** | [AgentGuard](https://github.com/AgentGuardHQ/agentguard) | Governance kernel — enforce/monitor policy on every action |
| 5 | 💻 **Code** | [OpenCode](https://github.com/opencode-ai/opencode) | AI coding framework (Go CLI, native tools) |
| 6 | 🤖 **Orchestrate** | [DeepAgents](https://github.com/langchain-ai/deepagents) (LangChain) | Multi-agent orchestration (autonomous task decomposition) |
| 7 | 🔒 **Sandbox** | [OpenShell](https://github.com/NVIDIA/OpenShell) (NVIDIA) | Kernel sandbox — Landlock + Seccomp BPF (Docker on macOS) |
| 8 | 🐾 **Scan** | [DefenseClaw](https://github.com/cisco-ai-defense/defenseclaw) (Cisco) | Supply chain scanner — AI Bill of Materials generation |

Check integration health at any time:

```bash
./shellforge status
# ✓ Ollama        running (qwen3:1.7b loaded)
# ✓ RTK           v0.4.2
# ✓ TurboQuant    configured
# ✓ AgentGuard    enforce mode (5 rules)
# ✓ OpenCode      v0.1.0
# ✓ DeepAgents    connected
# ✓ OpenShell     Docker sandbox active
# ✓ DefenseClaw   scanner ready
# Status: 8/8 integrations healthy
```

---

## CLI Commands

| Command | Description |
|---------|-------------|
| `./shellforge status` | Show ecosystem health (all 8 integrations) |
| `./shellforge qa` | Run the QA agent (test gap analysis) |
| `./shellforge report` | Run the report agent (markdown summary) |
| `./shellforge agent` | Run a custom agent with a prompt |
| `./shellforge scan` | Run security scan via DefenseClaw |
| `./shellforge setup` | Interactive setup wizard |
| `./shellforge version` | Show version |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Your Prompt                                                │
│  "Analyze this repo for test gaps"                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
                  ┌────────▼────────┐
                  │  ⚡ RTK          │  Strip 70–90% of terminal
                  │  Token Compress  │  noise before the LLM sees it
                  └────────┬────────┘
                           │
                  ┌────────▼────────┐
                  │  🧠 TurboQuant   │  6× KV cache compression
                  │  Quantization    │  14B models on 8 GB RAM
                  └────────┬────────┘
                           │
                  ┌────────▼────────┐
                  │  🦙 Ollama       │  Local inference
                  │  qwen3 · mistral │  any GGUF model
                  └────────┬────────┘
                           │
              ┌────────────▼────────────────┐
              │  🔥 ShellForge (Go binary)   │
              │                              │
              │  ┌─ OpenCode adapter ──────┐ │
              │  │  (preferred agent loop)  │ │  Pluggable engine
              │  ├─ DeepAgents adapter ────┤ │  interface picks the
              │  │  (multi-agent planning)  │ │  best available
              │  ├─ Native loop (fallback) ┤ │  framework at runtime
              │  │  internal/agent/loop.go  │ │
              │  └─────────────────────────┘ │
              │                              │
              │  Tools: read_file │ write_file│
              │  run_shell │ list_files       │
              │  search_files                 │
              └────────────┬────────────────-┘
                           │ tool call
           ════════════════╪════════════════
           ║  🛡️ AgentGuard Governance     ║
           ║  agentguard.yaml              ║
           ║  enforce / monitor            ║
           ║  allow · deny · audit         ║
           ║  every. single. action.       ║
           ════════════════╪════════════════
                           │ approved action
                  ┌────────┼────────┐
                  ▼        ▼        ▼
             ┌────────┐ ┌─────┐ ┌──────────┐
             │ Files  │ │ Git │ │  Shell   │
             │        │ │     │ │ (RTK)    │
             └────────┘ └─────┘ └──────────┘
                  🌍 Your Environment
                  (sandboxed by OpenShell)
```

---

## Governance

ShellForge's core value is **governance**. The AgentGuard engine (`internal/governance/engine.go`) parses `agentguard.yaml` and intercepts every tool call before execution.

```yaml
# agentguard.yaml — policy-as-code for every agent action
mode: enforce  # enforce | monitor

policies:
  - name: no-force-push
    action: deny
    pattern: "git push --force"

  - name: no-destructive-rm
    action: deny
    pattern: "rm -rf"

  - name: file-write-bounds
    action: deny
    description: "Agents can only write within the project directory"

  - name: bounded-execution
    action: deny
    description: "5-minute timeout per agent run"

  - name: no-secret-access
    action: deny
    pattern: "*.env|*id_rsa|*id_ed25519"
```

- **`enforce` mode** — violations are blocked and logged
- **`monitor` mode** — violations are logged but not blocked (use while tuning policies)

---

## Project Structure

```
shellforge/
├── cmd/shellforge/
│   ├── main.go                    # CLI entry (setup, qa, report, agent, scan, status, version)
│   └── status.go                  # Ecosystem health check (all 8 integrations)
├── internal/
│   ├── governance/engine.go       # Parses agentguard.yaml, enforce/monitor mode
│   ├── ollama/client.go           # Ollama HTTP client
│   ├── agent/loop.go              # Multi-turn agentic loop with tool calls (native fallback)
│   ├── tools/
│   │   ├── tools.go               # 5 tools: read_file, write_file, run_shell, list_files, search_files
│   │   └── rtk_shell.go           # RTK-wrapped shell execution
│   ├── logger/logger.go           # Structured JSON logging
│   ├── engine/
│   │   ├── engine.go              # Pluggable engine interface
│   │   ├── opencode.go            # OpenCode subprocess adapter
│   │   └── deepagents.go          # DeepAgents subprocess adapter
│   └── integration/
│       ├── rtk.go                 # RTK token compression
│       ├── openshell.go           # NVIDIA OpenShell sandbox
│       ├── defenseclaw.go         # Cisco DefenseClaw scanner
│       ├── turboquant.go          # Google TurboQuant quantization
│       └── agentguard.go          # AgentGuard kernel integration
├── scripts/
│   ├── setup.sh                   # Interactive installer (--all, --minimal)
│   ├── run-agent.sh
│   ├── run-qa-agent.sh
│   └── run-report-agent.sh
├── agentguard.yaml                # Governance policy (5 rules, enforce mode)
├── go.mod                         # github.com/AgentGuardHQ/shellforge
└── go.sum
```

---

## Build & Development

```bash
# Build
go build -o shellforge ./cmd/shellforge/

# Run directly
go run ./cmd/shellforge/ status
go run ./cmd/shellforge/ qa

# Test
go test ./...
```

### Model Options

| Model | Params | RAM | Best For |
|-------|--------|-----|----------|
| `qwen3:1.7b` | 1.7B | ~1.2 GB | Fast tasks, prototyping |
| `qwen3:4b` | 4B | ~3 GB | Balanced reasoning |
| `mistral:7b` | 7B | ~5 GB | Complex analysis |

---

## Cron Automation

```cron
# crontab -e
*/10 * * * * /path/to/shellforge/scripts/run-qa-agent.sh
*/30 * * * * /path/to/shellforge/scripts/run-report-agent.sh
```

All scripts are idempotent and timeout-safe.

---

## macOS (Apple Silicon / M4) Support

ShellForge runs natively on macOS with Apple Silicon (M1–M4). Notes:

- **Ollama** uses Metal GPU acceleration automatically — no CUDA needed
- **TurboQuant** KV cache compression makes 14B models fit in 8 GB unified memory
- **OpenShell** requires Docker (via [Colima](https://github.com/abiosoft/colima) or Docker Desktop) since Landlock/Seccomp are Linux-only kernel features

```bash
# macOS: install Colima for OpenShell sandbox support
brew install colima docker
colima start
```

---

## The AgentGuard Platform

ShellForge is part of the **AgentGuard** ecosystem:

| Project | What It Does |
|---------|--------------|
| [**AgentGuard**](https://github.com/AgentGuardHQ/agentguard) | Governance gateway — policy enforcement for Claude Code, Codex, Copilot, Gemini, OpenCode, DeepAgents |
| [**AgentGuard Cloud**](https://github.com/AgentGuardHQ/agentguard-cloud) | SaaS dashboard — observability, session replay, swarm org chart |
| **ShellForge** ← you are here | Governed local agent runtime — Go binary + Ollama + 8 integrations |

---

## Contributing

ShellForge is open source and actively developed. We welcome:

- New integration adapters (`internal/integration/`)
- Engine adapters for additional frameworks (`internal/engine/`)
- Governance policy templates
- Tool implementations (`internal/tools/`)
- Documentation improvements

```bash
# Fork, branch, build, test, PR
git checkout -b feat/my-feature
go build ./cmd/shellforge/
go test ./...
```

See [docs/roadmap.md](docs/roadmap.md) for what's planned.

---

<div align="center">

**[🌐 Website](https://agentguardhq.github.io/shellforge)** · **[⭐ Star on GitHub](https://github.com/AgentGuardHQ/shellforge)** · **[🛡️ AgentGuard](https://github.com/AgentGuardHQ/agentguard)**

*Built with 🔥 by humans and agents*

MIT License

</div>
