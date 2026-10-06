#!/usr/bin/env bash
# 236-v1 mutation runner: overlay only, never edits the shared worktree.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
LOGROOT=".scratch/wisp/probes/236/v1/logs/mut"
mkdir -p "$LOGROOT"
SUM="$LOGROOT/summary.txt"
: > "$SUM"
for name in "$@"; do
  ov="D:/tmp/wisp236v1/overlay/$name.json"
  log="$LOGROOT/$name.txt"
  go test ./internal/tools/ -count=1 -v -overlay "$ov" > "$log" 2>&1
  rc=$?
  p=$(grep -c '^--- PASS' "$log"); f=$(grep -c '^--- FAIL' "$log"); s=$(grep -c '^--- SKIP' "$log")
  pan=$(grep -c '^panic' "$log")
  {
    echo "=== $name rc=$rc PASS=$p FAIL=$f SKIP=$s panic-lines=$pan"
    grep '^--- FAIL' "$log" | sed 's/ (.*//'
    grep -E '^(panic|fatal error|FAIL\s+github)' "$log" | head -6
    grep -E '^ok\s+github|^FAIL\s+github|^# ' "$log" | head -4
  } >> "$SUM"
  cp "$ov" ".scratch/wisp/probes/236/v1/logs/overlay/$name.json" 2>/dev/null
  echo "done $name rc=$rc PASS=$p FAIL=$f"
done
