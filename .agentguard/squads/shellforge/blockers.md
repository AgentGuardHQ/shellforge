# ShellForge Squad — Blockers

**Updated:** 2026-03-29T20:00Z
**Reported by:** EM run 5 (claude-code:opus:shellforge:em)

---

## P0 — Critical Blockers (2)

### 1. All 3 PRs Awaiting Human Review — BLOCKING SQUAD PROGRESS
**Description:** All 3 open PRs are passing CI (5/5 checks each) but blocked on `REVIEW_REQUIRED`. GitHub branch protection prevents the EM (authored as jpleva91) from self-approving.
**PRs blocked:**
- **#83** — `fix(p0): close governance fail-open vulnerabilities` — closes #58, #59, #62, #67, #69, #75
- **#84** — `fix(docs): update stale Crush comments in cmdEvaluate (#74)` — closes #74
- **#85** — `chore(squad): EM state update — run 4` — squad ops housekeeping

**Action Required:** @jpleva91 or a collaborator must review and approve PRs #83, #84, #85.
**Priority:** Review #83 first — it carries all P0/P1 governance security fixes.

### 2. PR Budget AT LIMIT (3/3) — No New Fix PRs Possible
**Description:** Squad has reached the max of 3 open PRs. No new work can be opened until at least one PR merges.
**Impact:** P2 bugs (#65 scheduler silent error, #66 flattenParams dead code, #52 cmdScan glob broken, #53 README stale) remain queued but cannot be addressed.
**Unblocked by:** Merging any of #83, #84, or #85.

---

## P1 — Remaining Work (queued, no new PRs until budget frees)

### #68 — Zero test coverage across all packages
**Severity:** High — governance runtime with no tests is unshipable
**Impact:** Can't validate fix correctness, no regression protection. Blocks dogfood credibility.
**Assignee:** qa-agent
**URL:** https://github.com/AgentGuardHQ/shellforge/issues/68

### #63 — classifyShellRisk prefix matching too broad
**Severity:** High — false read-only classification on commands starting with `cat`/`ls`/`echo`
**Assignee:** qa-agent
**URL:** https://github.com/AgentGuardHQ/shellforge/issues/63

---

## P2 — Unassigned (queued, blocked by PR budget)

| # | Issue | Notes |
|---|-------|-------|
| #65 | scheduler.go silent os.WriteFile error | Silent failure on job persistence |
| #66 | flattenParams dead code | Logic bug, result overwritten before use |
| #52 | filepath.Glob ** never matches Go files | cmdScan broken for entire scan feature |
| #53 | README stale ./shellforge commands | Docs rot |

---

## Resolved (pending merge of PR #83)

- **#58** — bounded-execution wildcard policy blocked all run_shell → fix in PR #83
- **#62** — cmdEvaluate fail-open on JSON unmarshal → fix in PR #83
- **#75** — govern-shell.sh printf injection → fix in PR #83
- **#67** — govern-shell.sh fragile sed output parsing → fix in PR #83
- **#69** — rm policy only blocked -rf/-fr, not plain rm → fix in PR #83
- **#59** — misleading `# Mode: monitor` comment with `mode: enforce` → fix in PR #83
- **#74** — stale crush references in cmdEvaluate → fix in PR #84

---

## Status Summary

| Item | Status |
|------|--------|
| PR #83 (P0 fixes) | CI ✅ 5/5 — REVIEW BLOCKED |
| PR #84 (P1 docs) | CI ✅ 5/5 — REVIEW BLOCKED |
| PR #85 (EM state) | CI ✅ 5/5 — REVIEW BLOCKED |
| PR budget | 3/3 AT LIMIT |
| Dogfood (#76) | BLOCKED on #83 merge |
| QA-agent (#63, #68) | Active |
| New fix PRs | BLOCKED until budget frees |
| Retry loops | None |
| Blast radius | Low |
