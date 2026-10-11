#!/usr/bin/env bash
# 304-r1 gates.sh — every gate lands its own literal rc=N line inside its artefact.
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/304/r1"
L="$D/logs"
GOFU="$(go env GOPATH)/bin/gofumpt.exe"
mkdir -p "$L"
cd "$R" || exit 1
HB="/d/tmp/wisp-304-r1/headblobs"; mkdir -p "$HB"
FILES="internal/llm/adaptertest/childjob_windows.go internal/llm/adaptertest/childjob_other.go internal/llm/adaptertest/mockllm.go internal/llm/openaichat/mockllm_integ_test.go"

echo "=== gates @ $(date '+%F %T %z') HEAD=$(git rev-parse --short HEAD)" >  "$L/g-header.txt"

GOFLAGS= go build ./... > "$L/g-build.txt" 2>&1;                 echo "build_rc=$?" >> "$L/g-build.txt"
GOFLAGS= go vet ./internal/llm/... ./internal/proc/ > "$L/g-vet.txt" 2>&1; echo "vet_rc=$?" >> "$L/g-vet.txt"
sh scripts/d22scan.sh > "$L/g-d22scan.txt" 2>&1;                 echo "d22scan_rc=$?" >> "$L/g-d22scan.txt"
bash scripts/check-path-length-budget.sh --with-self-test > "$L/g-pathlen.txt" 2>&1; echo "pathlen_rc=$?" >> "$L/g-pathlen.txt"
bash scripts/portable-tests.sh --scope=census > "$L/g-census-after.txt" 2>&1; echo "census_after_rc=$?" >> "$L/g-census-after.txt"

# format roster, two rulers side by side, same scope (the four files this leg touched)
: > "$L/g-format-worktree.txt"
for f in $FILES; do "$GOFU" -l "$f" >> "$L/g-format-worktree.txt" 2>&1; done
echo "format_worktree_rc=$?" >> "$L/g-format-worktree.txt"
echo "scope=$FILES ruler=gofumpt -l (worktree)" >> "$L/g-format-worktree.txt"

: > "$L/g-format-headblob.txt"
for f in $FILES; do b=$(basename "$f"); git show "HEAD:$f" > "$HB/$b" 2>> "$L/g-format-headblob.txt"; "$GOFU" -l "$HB/$b" >> "$L/g-format-headblob.txt" 2>&1; done
echo "format_headblob_rc=$?" >> "$L/g-format-headblob.txt"
echo "scope=same four paths, dumped from git show HEAD:<path> into $HB ruler=gofumpt -l (HEAD blob)" >> "$L/g-format-headblob.txt"

{ echo "census_before_totals=$(grep -oE 'packages=[0-9]+ with-zero-compiled-tests=[0-9]+ claimed-by-no-scope=[0-9]+ unclaimed-with-tests=[0-9]+' "$L/g-census-before.txt" | tail -1)"
  echo "census_after_totals=$(grep -oE 'packages=[0-9]+ with-zero-compiled-tests=[0-9]+ claimed-by-no-scope=[0-9]+ unclaimed-with-tests=[0-9]+' "$L/g-census-after.txt" | tail -1)"
  if [ "$(grep -oE 'packages=[0-9]+ with-zero-compiled-tests=[0-9]+ claimed-by-no-scope=[0-9]+ unclaimed-with-tests=[0-9]+' "$L/g-census-before.txt" | tail -1)" = "$(grep -oE 'packages=[0-9]+ with-zero-compiled-tests=[0-9]+ claimed-by-no-scope=[0-9]+ unclaimed-with-tests=[0-9]+' "$L/g-census-after.txt" | tail -1)" ]; then
    echo "census_totals_delta=IDENTICAL"; else echo "census_totals_delta=DIFFERS"; fi
  echo "ticket_AC3_stated_baseline=packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0"
  echo "worktree_vs_headblob_format_delta=$(wc -l < "$L/g-format-worktree.txt") worktree-lines vs $(wc -l < "$L/g-format-headblob.txt") headblob-lines (both include their own header/rc lines)"
} > "$L/g-summary.txt"
grep -hE "_rc=" "$L"/g-*.txt | sed 's/^/  /'
cat "$L/g-summary.txt"
