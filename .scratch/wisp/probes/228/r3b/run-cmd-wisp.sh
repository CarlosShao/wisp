#!/bin/sh
# 228-r3b colour runner. Runs the cmd/wisp package in the OUT-OF-REPO export tree
# only (never the shared worktree). Surface picks the temp root volume:
#   c = default C-volume temp (leave TMP/TEMP/TMPDIR alone)
#   d = TMP/TEMP/TMPDIR pointed at a pre-created D:-volume dir
# GOTMPDIR is always the pre-created D: work dir (a missing one == "creating work dir" red).
# Usage: run-cmd-wisp.sh <c|d> <logfile> [extra go-test args...]
set -u
EXPORT=/d/wisp228r3b
SURFACE="$1"; shift
LOG="$1"; shift
export PATH="$EXPORT/third_party/sherpa-onnx:$EXPORT/build:$PATH"
export GOTMPDIR="$EXPORT/tmp"
if [ "$SURFACE" = "d" ]; then
  mkdir -p "$EXPORT/tdst"
  export TMPDIR="$EXPORT/tdst"
  export TMP="$EXPORT/tdst"
  export TEMP="$EXPORT/tdst"
fi
cd "$EXPORT" || { echo "cd export failed" >"$LOG"; exit 9; }
echo "### surface=$SURFACE pwd=$PWD GOTMPDIR=$GOTMPDIR TMP=${TMP:-<default-c>} $(date '+%F %T %z')" >"$LOG"
# -v full verbosity to a file (do NOT pipe through tail here; we need every line for the four counts)
go test ./cmd/wisp/ -count=1 -timeout 30m -v "$@" >>"$LOG" 2>&1
RC=$?
echo "### go-test-rc=$RC surface=$SURFACE $(date '+%F %T %z')" >>"$LOG"
echo "DONE rc=$RC surface=$SURFACE"
