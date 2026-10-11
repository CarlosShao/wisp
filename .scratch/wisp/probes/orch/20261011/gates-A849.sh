#!/usr/bin/env bash
# A849 precommit gates + hygiene. Batch = ledger + handover + my probes + the 代提 msg. Zero .go.
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
L=$D/logs/orch-A849-gates.txt

echo "=== A849 gates (self-landed rc) ===" > "$L"
echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD) batch=A849" >> "$L"
GOFLAGS= go build ./... > "$D/logs/orch-A849-build.txt" 2>&1; echo "gate_build_rc=$?" >> "$L"
GOFLAGS= go vet ./internal/tools/ > "$D/logs/orch-A849-vet-tools.txt" 2>&1; echo "gate_vet_tools_rc=$?" >> "$L"
sh scripts/d22scan.sh > "$D/logs/orch-A849-d22scan.txt" 2>&1; echo "gate_d22scan_rc=$?" >> "$L"
bash scripts/check-path-length-budget.sh --with-self-test > "$D/logs/orch-A849-pathlen.txt" 2>&1; echo "gate_pathlen_rc=$?" >> "$L"
grep -E 'VERDICT|not in roster|worst full path' "$D/logs/orch-A849-pathlen.txt" | tail -3 >> "$L"
echo "d22scan_rc_line_in_artefact=$(grep -oE 'rc_d22scan=[0-9]+|clean - no D22 ban' "$D/logs/orch-A849-d22scan.txt" | tail -1)" >> "$L"
# land rc into silent-green artefacts (⛔ 0-byte cell)
for f in orch-A849-build.txt orch-A849-vet-tools.txt; do
  [ -s "$D/logs/$f" ] || printf 'rc=0 (silent when green; this line lands the rc)\n' >> "$D/logs/$f"
done

git add -- "$D" docs/reports/pending-and-issues.md docs/reports/HANDOVER.md
{ echo "=== hygiene ==="
  echo "staged_total=$(git -c core.quotePath=false diff --cached --name-only | grep -c '')"
  echo "staged_go=$(git -c core.quotePath=false diff --cached --name-only | grep -c '[.]go$' || true)"
  echo "foreign_staged=$(git -c core.quotePath=false diff --cached --name-only | grep -v '^docs/reports/\|^\.scratch/wisp/probes/orch/20261011/' | grep -c '' || true)"
  echo "inflight_staged=$(git -c core.quotePath=false diff --cached --name-only | grep -cE '^internal/|^cmd/|^tools/|^scripts/|^\.github/|^design/|^\.gitignore$|^go\.(mod|sum)$' || true)"
  z=0; while IFS= read -r -d '' f; do [ -s "$f" ] || { echo "ZERO: $f"; z=$((z+1)); }; done < <(git diff --cached --name-only -z)
  echo "zero_byte_staged_listed_above"
  echo "numstat_ledger=$(git diff --cached --numstat -- docs/reports/pending-and-issues.md | tr '\n' ' ')"
  echo "numstat_handover=$(git diff --cached --numstat -- docs/reports/HANDOVER.md | tr '\n' ' ')"
  echo "prefix_ledger=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -"$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)" docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -"$(git show HEAD:docs/reports/HANDOVER.md | wc -l)" docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "index_lock=$(ls .git/index.lock 2>/dev/null | grep -c '' || true)"
  echo "worktree_dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools | wc -l | tr -d ' ')"
} >> "$L"
cat "$L"
