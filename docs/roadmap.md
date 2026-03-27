# Roadmap

## Phase 1 — Foundation ✅
- [x] Ollama integration with low-context wrapper
- [x] 3 simple agents (QA, report, prototype)
- [x] AgentGuard governance policy (monitor mode)
- [x] Script-based execution with cron support
- [x] Memory optimization placeholder

## Phase 2 — Hardening ✅
- [x] Go rewrite — single static binary (~7.5MB), zero Node.js dependencies
- [x] Switch agentguard.yaml to `enforce` mode
- [x] AgentGuard CLI hooks integrated into governance engine
- [x] Token budget tracking per agent per day
- [x] Output quality scoring (simple heuristics)
- [x] Error recovery and retry logic

## Phase 3 — Framework Integration ✅
- [x] **OpenCode** — Go CLI AI coding framework
  - Pluggable engine interface (`internal/engine/`)
  - `--non-interactive` subprocess mode, governance-wrapped
  - Tool-use governance via AgentGuard policy engine
- [x] **DeepAgents** — multi-agent orchestration (LangChain-based)
  - Subprocess engine adapter (`internal/engine/`)
  - Agent decomposition: break goals into sub-tasks
  - Governance-wrapped tool calls

## Phase 4 — Memory & Context ✅
- [x] **RTK v0.31.0** — token compression integrated
  - Auto-wraps shell output and LLM I/O
  - Reduces context window usage by ~40%
- [x] **TurboQuant** — model quantization + KV cache optimization
  - PyTorch MPS backend on Apple Silicon
  - Integrated via `internal/integration/`
- [x] Prompt caching for repeated patterns

## Phase 5 — Security ✅
- [x] **NVIDIA OpenShell** sandbox integration
  - Landlock + Seccomp isolation per agent run
  - Docker/Colima on Mac for Linux kernel features
  - Integrated via `internal/integration/`
- [x] **Cisco DefenseClaw** scanning
  - AI Bill of Materials (BoM) scanner
  - Scan agent skills/plugins pre-install
  - Integrated via `internal/integration/`

## Phase 6 — Scale 🔄 In Progress
- [x] Interactive setup CLI (`shellforge setup`)
- [x] Ecosystem health check (`shellforge status`)
- [ ] Binary releases (goreleaser, GitHub Releases)
- [ ] Multi-model routing (qwen for fast, mistral for quality)
- [ ] Cross-platform support (Linux arm64, Windows)
- [ ] Cloud telemetry integration (AgentGuard Cloud)
- [ ] Dashboard for local swarm observability
