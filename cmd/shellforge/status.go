package main

import (
"fmt"
"os/exec"
"strings"

"github.com/AgentGuardHQ/shellforge/internal/governance"
"github.com/AgentGuardHQ/shellforge/internal/integration"
"github.com/AgentGuardHQ/shellforge/internal/ollama"
)

func cmdStatusFull() {
fmt.Printf("ShellForge %s — Full Ecosystem Status\n", version)
fmt.Println(strings.Repeat("─", 50))

// ── Model Layer ──
fmt.Println("\n🦙 Model Layer")
if ollama.IsRunning() {
models, _ := ollama.ListModels()
fmt.Printf("  ✓ Ollama: running (%d models)\n", len(models))
for _, m := range models {
tag := ""
if m == ollama.Model {
tag = " ← active"
}
fmt.Printf("    • %s%s\n", m, tag)
}
} else {
fmt.Println("  ✗ Ollama: not running (start: ollama serve)")
}

// ── Token Compression ──
fmt.Println("\n⚡ Token Compression")
rtk := integration.NewRTK()
if rtk.Available() {
fmt.Printf("  ✓ RTK: %s (60-90%% token savings on shell output)\n", rtk.Version())
} else {
fmt.Println("  ○ RTK: not installed (brew install rtk-ai/tap/rtk)")
}

// ── Memory Optimization ──
fmt.Println("\n🧠 Memory Optimization")
tq := integration.NewTurboQuant()
if tq.Available() {
fmt.Println("  ✓ TurboQuant: installed (3-bit KV cache, 6x compression)")
est := tq.EstimateMemory(1.7, 4096)
fmt.Printf("    qwen3:1.7b → %.1fGB standard, %.1fGB with TQ (%.0f%% savings)\n",
est.TotalStandard, est.TotalTQ, est.SavingsPercent)
} else {
fmt.Println("  ○ TurboQuant: not installed (pip install turboquant-pytorch)")
}

// ── Governance ──
fmt.Println("\n🛡️  Governance")
configPath := findGovernanceConfig()
if configPath != "" {
eng, err := governance.NewEngine(configPath)
if err != nil {
fmt.Printf("  ✗ Policy: %s\n", err)
} else {
fmt.Printf("  ✓ Policy: mode=%s, %d rules (%s)\n", eng.Mode, len(eng.Policies), configPath)
for _, p := range eng.Policies {
fmt.Printf("    • %s [%s] %s\n", p.Name, p.Action, p.Description)
}
}
} else {
fmt.Println("  ✗ Policy: no agentguard.yaml found")
}

agk := integration.NewAgentGuardKernel()
if agk.Available() {
fmt.Printf("  ✓ Kernel: %s (full evaluation — blast radius, personas, invariants)\n", agk.Version())
} else {
fmt.Println("  ○ Kernel: built-in YAML evaluator (install kernel for full evaluation)")
}

// ── Agent Engines ──
fmt.Println("\n🤖 Agent Engines")
fmt.Println("  ✓ native: built-in Ollama loop with tool use")
if _, err := exec.LookPath("opencode"); err == nil {
fmt.Println("  ✓ opencode: detected (Go-native AI coding agent)")
} else {
fmt.Println("  ○ opencode: not installed (npm i -g opencode-ai)")
}
checkNodeModule("deepagents", "deepagents", "multi-step planning via LangGraph")

// ── Security ──
fmt.Println("\n🔒 Security")
openshell := integration.NewOpenShell()
if openshell.Available() {
fmt.Println("  ✓ OpenShell: NVIDIA kernel sandbox (Landlock + Seccomp)")
} else {
fmt.Println("  ○ OpenShell: not installed (https://github.com/NVIDIA/OpenShell)")
}

dc := integration.NewDefenseClaw()
if dc.Available() {
fmt.Println("  ✓ DefenseClaw: Cisco supply chain scanner")
} else {
fmt.Println("  ○ DefenseClaw: not installed (https://github.com/cisco/defenseclaw)")
}

// ── Summary ──
fmt.Println("\n" + strings.Repeat("─", 50))
available := countAvailable(rtk.Available(), tq.Available(), agk.Available(),
openshell.Available(), dc.Available())
fmt.Printf("Stack: %d/8 integrations active\n", 2+available)
}

func checkNodeModule(pkg, name, desc string) {
cmd := exec.Command("node", "-e", fmt.Sprintf("require('%s')", pkg))
if cmd.Run() == nil {
fmt.Printf("  ✓ %s: installed (%s)\n", name, desc)
} else {
fmt.Printf("  ○ %s: not installed (npm i %s)\n", name, pkg)
}
}

func countAvailable(flags ...bool) int {
n := 0
for _, f := range flags {
if f {
n++
}
}
return n
}
