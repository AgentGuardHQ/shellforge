# ShellForge Squad — Blockers

**Updated:** 2026-03-29T18:00Z
**Reported by:** EM run (claude-code:opus:shellforge:em)

---

## P0 — Active Blockers (0)

All 3 P0 governance security bugs are fixed in PR #83 (pending CI + merge).

See PR: https://github.com/AgentGuardHQ/shellforge/pull/83

---

## P1 — Remaining Work

### #68 — Zero test coverage across all packages
**Severity:** High — governance runtime with no tests is unshipable
**Impact:** Can't validate fix correctness, no regression protection. Blocks dogfood credibility.
**Assignee:** qa-agent
**URL:** https://github.com/AgentGuardHQ/shellforge/issues/68

### #63 — classifyShellRisk prefix matching too broad
**Severity:** High — false read-only classification on commands starting with `cat`/`ls`/`echo`
**Assignee:** qa-agent
**URL:** https://github.com/AgentGuardHQ/shellforge/issues/63

### #74 — Stale crush references in main.go
**Severity:** Low-medium — cosmetic but misleading; crush→goose migration was v0.6
**URL:** https://github.com/AgentGuardHQ/shellforge/issues/74

---

## Resolved This Run

- **#58** — bounded-execution wildcard policy matched every run_shell → `engine.go` fix merged in PR #83
- **#62** — cmdEvaluate fail-open on JSON unmarshal → fail-closed fix in PR #83
- **#75** — govern-shell.sh printf injection → jq --arg fix in PR #83
- **#67** — govern-shell.sh fragile sed output parsing → jq fix in PR #83
- **#69** — rm policy only blocked -rf/-fr, not plain rm → policy broadened in PR #83
- **#59** — misleading `# Mode: monitor` comment with `mode: enforce` → fixed in PR #83

---

## Notes

- PR budget: 1/3 open — capacity for 2 more fix PRs
- No retry loops or blast radius concerns
- Dogfood run (#76) unblocked once PR #83 merges
- Test coverage (#68) is now the most pressing remaining gap — no regression safety net
