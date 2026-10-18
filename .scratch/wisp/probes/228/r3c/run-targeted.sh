#!/bin/sh
# 228-r3c stage 2: colour runner for the TARGETED -run reads only.
# Runs only in the OUT-OF-REPO export tree; never the shared worktree.
# Surface picks the temp root volume:
#   c = default C-volume temp (TMP/TEMP/TMPDIR untouched)
#   d = TMP/TEMP/TMPDIR pointed at the pre-created D: dir (.../tdst)
# GOTMPDIR is always the pre-created D: dir (.../gotmp).
# Usage: run-targeted.sh <c|d> <logfile-name-only> <one -run pattern>
set -u
ROOT=/d/tmp/wisp/20261018/228r3c
EXPORT=$ROOT/tree
REPO='/d/work/workspace/projects plans/Wisp'
HASHSH="$REPO/.scratch/wisp/probes/228/r3c/hash.sh"
SURFACE="$1"; shift
LOGNAME="$1"; shift
RUN="$1"
LOG=$ROOT/logs/$LOGNAME
mkdir -p "$ROOT/logs"
export PATH="$EXPORT/third_party/sherpa-onnx:$EXPORT/build:$PATH"
export GOTMPDIR="$ROOT/gotmp"
if [ "$SURFACE" = "d" ]; then
  mkdir -p "$ROOT/tdst"
  export TMPDIR="$ROOT/tdst"
  export TMP="$ROOT/tdst"
  export TEMP="$ROOT/tdst"
fi
cd "$EXPORT" || { echo "cd export failed" >"$LOG"; exit 9; }
{
  echo "### leg=228-r3c surface=$SURFACE pwd=$PWD GOTMPDIR=$GOTMPDIR TMP=${TMP:-<default-c>}"
  echo "### run-pattern=$RUN date=$(date '+%F %T %z')"
  echo "### shortname-ruler(cygpath -d tree)=$(cygpath -d "$EXPORT")"
  echo "### hash-of-mutated-file-BEFORE-this-run"
  sh "$HASHSH" "$EXPORT/internal/proc/boot_windows.go"
} >"$LOG" 2>&1
go test ./cmd/wisp/ -count=1 -timeout 30m -v -run "$RUN" >>"$LOG" 2>&1
RC=$?
{
  echo "### go-test-rc=$RC surface=$SURFACE $(date '+%F %T %z')"
  echo "### hash-of-mutated-file-AFTER-this-run"
  sh "$HASHSH" "$EXPORT/internal/proc/boot_windows.go"
} >>"$LOG" 2>&1
# 四数尺：先读进变量再落笔——直接 `grep "$LOG" >>"$LOG"` 会撞上
# "input file is also the output"（GNU grep 拒）＝台面故障，⛔ 计入判据。
RUN=$(grep -c '^=== RUN' "$LOG" || true)
PASSTOP=$(grep -c '^--- PASS' "$LOG" || true)
FAILTOP=$(grep -c '^--- FAIL' "$LOG" || true)
SKIP=$(grep -c '^--- SKIP' "$LOG" || true)
COLOUR=$(grep -E '^--- (PASS|FAIL|SKIP)' "$LOG" || true)
NAMES=$(grep -E '^(=== RUN|--- (PASS|FAIL|SKIP))' "$LOG" || true)
{
  echo "### four-counts (尺=grep -c on this log；RUN=^=== RUN 含子测试 / PASS_TOP=^--- PASS 只顶层 / FAIL_TOP=^--- FAIL / SKIP=^--- SKIP；四数不混用不相减)"
  printf 'RUN=%s\nPASS_TOP=%s\nFAIL_TOP=%s\nSKIP=%s\n' "$RUN" "$PASSTOP" "$FAILTOP" "$SKIP"
  echo "### per-test-colour"
  printf '%s\n' "$COLOUR"
  echo "### per-test-names"
  printf '%s\n' "$NAMES"
} >>"$LOG" 2>&1
echo "DONE rc=$RC surface=$SURFACE log=$LOG"
