#!/bin/sh
# 333-v1 colour runs: SAME command, only the temp-root volume moves. Four runs:
# pre-anchor(P) x {C,D}  and  delivery(A) x {C,D}.  Tree = out-of-repo export stems.
# RUN counts ^=== RUN (incl subtests) / PASS_TOP ^--- PASS (top level only) /
# FAIL_TOP ^--- FAIL / SKIP ^--- SKIP ; -count=1 -timeout 30m, never default timeout.
set -u
LOGD="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/333/v1/logs"
GOFLAGS=
export GOFLAGS
run_one() { # $1=label $2=tree $3=surface(C|D)
  L="$1"; T="$2"; S="$3"
  OUT="$LOGD/$L.log"
  if [ "$S" = "C" ]; then
    export TMP="C:/Users/swq/AppData/Local/Temp" TEMP="C:/Users/swq/AppData/Local/Temp"
    unset GOTMPDIR
  else
    mkdir -p "D:/wisp-333v1-tmp/tmp" "D:/wisp-333v1-tmp/gotmp"
    export TMP="D:/wisp-333v1-tmp/tmp" TEMP="D:/wisp-333v1-tmp/tmp" GOTMPDIR="D:/wisp-333v1-tmp/gotmp"
  fi
  cd "$T"
  { echo "# label=$L tree=$T surface=$S"; echo "# cmd=go test -v ./internal/panel/ -count=1 -timeout 30m"; echo "# TMP=$TMP GOTMPDIR=${GOTMPDIR:-unset}"; } > "$OUT"
  go test -v ./internal/panel/ -count=1 -timeout 30m >> "$OUT" 2>&1
  RC=$?
  echo "rc=$RC" >> "$OUT"
  RUN=$(grep -c '^=== RUN' "$OUT" || true); PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true); SK=$(grep -c '^--- SKIP' "$OUT" || true)
  echo "FOUR $L: RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK rc=$RC"
  grep '^--- FAIL' "$OUT" | sed 's/^--- FAIL *//' | awk '{print $1}' | sort > "$LOGD/$L.failnames.txt"
  sha256sum "$OUT" | cut -d' ' -f1 > "$LOGD/$L.log.sha256"
  cd - >/dev/null
}
run_one col-pre-c  /d/W333V2-hLzM9f C
run_one col-pre-d  /d/W333V2-hLzM9f D
run_one col-del-c  /d/W333V1-SuqCO8 C
run_one col-del-d  /d/W333V1-SuqCO8 D
