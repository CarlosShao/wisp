#!/bin/bash
# 221-v1 mutation runner: one test binary at a time, overlay only, never the worktree.
ROOT="D:/work/workspace/projects plans/Wisp"
cd "$ROOT" || exit 1
MUT="$ROOT/.scratch/wisp/probes/221/v1/mutations"
LOGS="$ROOT/.scratch/wisp/probes/221/v1/logs"
RUNSEL='Test221|Test197|TestEveryRegisteredToolIsClassifiedForMarking175r2'
export PATH="$ROOT/third_party/sherpa-onnx:$ROOT/build:$PATH"

run_one() {
  local name="$1" overlay="$2"
  local log="$LOGS/mut-$name.txt"
  if [ -z "$overlay" ]; then
    go test ./internal/tools -run "$RUNSEL" -v -count=1 > "$log" 2>&1
  else
    go test ./internal/tools -run "$RUNSEL" -v -count=1 -overlay "$overlay" > "$log" 2>&1
  fi
  local rc=$?
  local reds
  reds=$(grep -cE '^[[:space:]]*--- FAIL' "$log")
  local total
  total=$(grep -cE '^[[:space:]]*--- (PASS|FAIL)' "$log")
  echo "=== $name rc=$rc names_seen=$total red=$reds"
  grep -E '^[[:space:]]*--- FAIL' "$log" | sed 's/^/    RED /'
}

run_one asis ""
for j in "$MUT"/*.json; do
  n=$(basename "$j" .json)
  # only run if the corresponding go blob exists
  [ -f "${j%.json}.go" ] || continue
  if ! tasklist //FI "IMAGENAME eq go.exe" | grep -q go.exe; then
    run_one "$n" "$j"
  else
    echo "=== $n SKIPPED: another go.exe is running, not contending"
  fi
done
echo "=== worktree check (must be empty) ==="
git status --porcelain -- internal cmd
