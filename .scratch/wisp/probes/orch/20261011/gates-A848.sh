#!/usr/bin/env bash
# A848 precommit gates + hygiene. This batch touches ZERO .go (docs + one ticket rename),
# so the Go-surface readings describe the tree, not my batch -- stated, not hidden.
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
L=$D/logs/orch-A848-gates.txt
T308=.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume-done.md
mkdir -p "$D/logs"

{ echo "=== gate rc lines (self-landed) ==="
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD) batch=A848"
} > "$L"

# build
GOFLAGS= go build ./... > "$D/logs/orch-A848-build.txt" 2>&1; echo "gate_build_rc=$?" >> "$L"
# vet risk + tools (tools because 309-r2's code just landed there)
GOFLAGS= go vet ./internal/risk/ > "$D/logs/orch-A848-vet-risk.txt" 2>&1; echo "gate_vet_risk_rc=$?" >> "$L"
GOFLAGS= go vet ./internal/tools/ > "$D/logs/orch-A848-vet-tools.txt" 2>&1; echo "gate_vet_tools_rc=$?" >> "$L"
# d22scan (whole repo, live worktree)
sh scripts/d22scan.sh > "$D/logs/orch-A848-d22scan.txt" 2>&1; echo "gate_d22scan_rc=$?" >> "$L"
echo "d22scan_findings_in_internal_risk=$(grep -c 'internal/risk' "$D/logs/orch-A848-d22scan.txt")" >> "$L"
echo "d22scan_findings_in_internal_tools=$(grep -c 'internal/tools' "$D/logs/orch-A848-d22scan.txt")" >> "$L"
echo "d22scan_total_finding_lines=$(grep -cE 'finding|violation|BAN|ban #' "$D/logs/orch-A848-d22scan.txt")" >> "$L"
# path length gate
bash scripts/check-path-length-budget.sh --with-self-test > "$D/logs/orch-A848-pathlen.txt" 2>&1; echo "gate_pathlen_rc=$?" >> "$L"
grep -E 'GREEN|RED|over_hat|over-hat|tracked=|worst' "$D/logs/orch-A848-pathlen.txt" | tail -6 >> "$L"

# ---- staging + hygiene (quotePath off, -z iteration; ⛔ && chains after grep -c) ----
# ⛔ `git add .scratch/wisp/issues/` —— 那会把别的腿正在写的票尾 bullet 一起暂存。
git add -- "$D" docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T308" >> /dev/null 2>&1
git -c core.quotePath=false diff --cached --name-only > "$D/logs/orch-A848-staged.txt"
n_all=$(grep -c '' "$D/logs/orch-A848-staged.txt")
n_go=$(grep -c '[.]go$' "$D/logs/orch-A848-staged.txt" || true)
n_out=$(grep -c '[.]out$' "$D/logs/orch-A848-staged.txt" || true)
foreign=$(grep -v '^docs/reports/\|^\.scratch/wisp/issues/\|^\.scratch/wisp/probes/orch/20261011/' "$D/logs/orch-A848-staged.txt" | grep -c '' || true)
inflight=$(grep -cE '^internal/|^cmd/|^tools/|^scripts/|^\.github/' "$D/logs/orch-A848-staged.txt" || true)
zero=$(git diff --cached --name-only -z | tr '\0' '\n' | while IFS= read -r f; do [ -s "$f" ] || echo "$f"; done | grep -c '' || true)
{ echo "=== hygiene ==="
  echo "staged_total=$n_all staged_go=$n_go staged_dot_out=$n_out foreign_staged=$foreign inflight_staged=$inflight zero_byte_staged=$zero"
  echo "numstat_ledger=$(git diff --cached --numstat -- docs/reports/pending-and-issues.md | tr '\n' ' ')"
  echo "numstat_handover=$(git diff --cached --numstat -- docs/reports/HANDOVER.md | tr '\n' ' ')"
  echo "numstat_t308=$(git diff --cached --numstat -M -- "$T308" | tr '\n' ' ')"
  echo "rename_t308=$(git diff --cached --name-status -M -- .scratch/wisp/issues/ | tr '\n' ' ')"
  echo "prefix_ledger=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -"$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)" docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -"$(git show HEAD:docs/reports/HANDOVER.md | wc -l)" docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "box_308=$(grep -cE '^[[:space:]]*- \[x\]' "$T308")/$(grep -cE '^[[:space:]]*- \[ \]' "$T308")"
  echo "status_done_lines_308=$(grep -c '^\*\*Status:\*\* \*\*done\*\*' "$T308")"
  echo "worktree_dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools | wc -l | tr -d ' ')"
  echo "ticket309_or_304_touched_by_me=$(git -c core.quotePath=false diff --cached --name-only | grep -cE 'issues/(309|304)' || true)"
} >> "$L"
cat "$L"
