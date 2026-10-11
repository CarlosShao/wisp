#!/usr/bin/env bash
# orchestrator re-run of 309-v2's load-bearing claims (verdict leg, non-implementer)
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A851-309v2-verify.txt"
: > "$L"
cd "$R" || exit 1

{ echo "=== A) commit reality + scope (its report: 3974ab4d, 83 files, +8924/-0, go=0, forbidden=0)"
  echo "cat-file=$(git cat-file -t 3974ab4d 2>&1) head=$(git rev-parse --short HEAD)"
  git -c core.quotePath=false show --name-only --format= 3974ab4d > "$D/logs/orch-A851-309v2-paths.txt"
  echo "files=$(grep -c '' "$D/logs/orch-A851-309v2-paths.txt")"
  echo "numstat=$(git diff --numstat 3974ab4d^ 3974ab4d | awk '{a+=$1; d+=$2} END{print "+"a" -"d}')"
  echo "go_paths=$(grep -c '[.]go$' "$D/logs/orch-A851-309v2-paths.txt" || true)"
  echo "forbidden_design=$(grep -c '^design/' "$D/logs/orch-A851-309v2-paths.txt" || true) forbidden_gitignore=$(grep -c '^\.gitignore$' "$D/logs/orch-A851-309v2-paths.txt" || true) forbidden_gomod=$(grep -c '^go\.mod$' "$D/logs/orch-A851-309v2-paths.txt" || true)"
  echo "inflight_faces=$(grep -cE '^internal/|^cmd/|^tools/|^scripts/|^\.github/' "$D/logs/orch-A851-309v2-paths.txt" || true)"
  echo "foreign_outside_allowed=$(grep -vE '^docs/evidence/s1/|^\.scratch/wisp/probes/309/v2/|^\.scratch/wisp/issues/' "$D/logs/orch-A851-309v2-paths.txt" | grep -c '' || true)"
  echo "--- foreign list (empty above = 0):"; grep -vE '^docs/evidence/s1/|^\.scratch/wisp/probes/309/v2/|^\.scratch/wisp/issues/' "$D/logs/orch-A851-309v2-paths.txt"
  echo "ticket309_numstat=$(git diff --numstat 3974ab4d^ 3974ab4d -- .scratch/wisp/issues/ | tr '\n' ' ')"
  echo "evidence_files=$(grep -c '^docs/evidence/s1/' "$D/logs/orch-A851-309v2-paths.txt" || true)"
  echo "staged_now=$(git diff --cached --name-only | grep -c '' || true) dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools scripts | wc -l | tr -d ' ')"
} >> "$L"

{ echo "=== B) ticket 309 state (its report: untouched 3/3)"
  T=$(git -c core.quotePath=false ls-files -- .scratch/wisp/issues | grep -E '/309-' | head -1)
  echo "path=$T bytes=$(wc -c <"$T")"
  echo "boxes_done=$(grep -cE '^[[:space:]]*- \[x\]' "$T") boxes_open=$(grep -cE '^[[:space:]]*- \[ \]' "$T")"
  echo "--- AC#3 line and AC#1/AC#2 lines (first 300 chars each):"
  grep -nE '^- \[[ x]\] \*\*`?AC#(1|2|3)`?' "$T" | cut -c1-300
  echo "--- ticket tail bullet added by it:"
  tail -2 "$T" | cut -c1-200
} >> "$L"

{ echo "=== C) the over-claim sentence it says is falsified (bridge.go:955-970)"
  B=internal/panel/bridge.go
  [ -f "$B" ] || B=$(git ls-files '*bridge.go' | head -1)
  echo "file=$B"
  sed -n '952,970p' "$B" | cat -n
} >> "$L"

{ echo "=== D) the nail: where it lives, its ruler literal, and its denominator constant"
  echo "--- find the test:"
  git -c core.quotePath=false grep -rn 'Test309AssessorForEmptyTaskIDDropsTaintDetector' -- internal | cut -c1-120
  echo "--- find 293 / denominator constants in that file:"
  F=$(git -c core.quotePath=false grep -rl 'Test309AssessorForEmptyTaskIDDropsTaintDetector' -- internal | head -1)
  echo "F=$F"
  [ -n "$F" ] && grep -nE '293|ToolRequest|new\(|wantTotal|composite|&agent\.' "$F" | cut -c1-150 | head -30
} >> "$L"

{ echo "=== E) does the instrument's pattern set cover new(agent.ToolRequest)?"
  echo "--- repo-wide spellings of a ToolRequest creation (three shapes):"
  echo "literal_braces=$(git -c core.quotePath=false grep -rE 'agent\.ToolRequest\{' -- internal cmd tools | grep -c '' || true)"
  echo "amp_form=$(git -c core.quotePath=false grep -rE '&agent\.ToolRequest\{' -- internal cmd tools | grep -c '' || true)"
  echo "new_form=$(git -c core.quotePath=false grep -rE 'new\(agent\.ToolRequest\)' -- internal cmd tools | grep -c '' || true)"
  echo "--- named new() sites:"
  git -c core.quotePath=false grep -rnE 'new\(agent\.ToolRequest\)' -- internal cmd tools | cut -c1-130
} >> "$L"

{ echo "=== F) run the nail on today's tree (green expected; ⛔ I do NOT flip on its word alone)"
  F=$(git -c core.quotePath=false grep -rl 'Test309AssessorForEmptyTaskIDDropsTaintDetector' -- internal | head -1)
  P=$(dirname "$F")
  echo "pkg=./$P/"
  GOFLAGS= go test ./"$P"/ -run '309' -count=1 2>&1 | tail -3
  echo "run_rc=$?"
} >> "$L"

{ echo "=== G) its zero-byte gate artefacts (rule: 0-byte gate file = cell not delivered)"
  for f in gate-build gate-vet gate-d22scan gate-pathlen gate2-d22scan gate2-pathlen r4-grep-rulers s9-behaviour rigA rigB rigC rigD rigE; do
    p=".scratch/wisp/probes/309/v2/logs/$f"
    for ext in log txt; do
      [ -f "$p.$ext" ] && printf "%-22s %s bytes=%s\n" "$f.$ext" "$(grep -c '' "$p.$ext")" "$(wc -c <"$p.$ext")"
    done
  done
  echo "zero_byte=$(find .scratch/wisp/probes/309/v2 -type f -size 0 | grep -c '' || true)"
  echo "--- rc lines self-landed?"; grep -rhoE 'rc=[0-9]+' .scratch/wisp/probes/309/v2/logs/*.log 2>/dev/null | sort | uniq -c | tr '\n' ' '
} >> "$L"

cat "$L"
