#!/bin/bash
# 321-r1 colour runner.
#   usage = run.sh <treeDir> <runRegex> <logName> [posix|win]
# Runs go test in an OUT-OF-REPO export tree with the sherpa-onnx + build DLL
# dirs at the FRONT of PATH. The dispatch gives the literal POSIX recipe
# PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" (mode=posix, default).
# Mode=win hands the child a Windows-form PATH through cmd.exe instead, because
# `env PATH='D:\...;...'` under MSYS bash cannot even locate `go` (that is how
# the first attempt died with rc=127: env: 'go': No such file or directory).
# Why the PATH at all: without those DLL dirs the test binary dies at load with
# 0xc0000135 and prints ZERO `--- FAIL` - four numbers 0/0/0/0 that look like a
# verdict and are台面故障, never a reading (public block, 颜色尺).
# `cygpath -wp` on the whole PATH is NOT used in the shared working tree (gcc
# disappears => CGO off => sherpa excluded by build constraints).
set -u
TREE="$1"; PAT="$2"; LOG="$3"; MODE="${4:-posix}"
ROOT="/d/tmp/wisp/20261020/321r1"
mkdir -p "$ROOT/gotmp" "$ROOT/logs"
cd "$TREE" || exit 9
GT=$(cygpath -w "$ROOT/gotmp")
L="$ROOT/logs/$LOG"
echo "TREE=$TREE PAT=$PAT MODE=$MODE" | tee "$L.meta"
echo "GOTMPDIR=$GT" | tee -a "$L.meta"
if [ "$MODE" = win ]; then
  P1=$(cygpath -w "$TREE/third_party/sherpa-onnx")
  P2=$(cygpath -w "$TREE/build")
  echo "PATH_HEAD=$P1;$P2" | tee -a "$L.meta"
  {
    tasklist //FI "IMAGENAME eq wisp.exe" 2>&1 | head -4
    cmd.exe /c "set GOTMPDIR=$GT&& set PATH=$P1;$P2;%PATH%&& go test -count=1 -timeout 30m -run $PAT ./cmd/wisp/ -v"
  } > "$L" 2>&1
else
  echo "PATH_HEAD=$TREE/third_party/sherpa-onnx:$TREE/build" | tee -a "$L.meta"
  {
    tasklist //FI "IMAGENAME eq wisp.exe" 2>&1 | head -4
    GOTMPDIR="$GT" PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH" \
      go test -count=1 -timeout 30m -run "$PAT" ./cmd/wisp/ -v
  } > "$L" 2>&1
fi
RC=$?
echo "RC=$RC"
echo "RUN=$(grep -c '^=== RUN' "$L" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$L" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$L" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$L" || true)"
echo "--top names--"
grep '^--- \(PASS\|FAIL\|SKIP\)' "$L" | sed 's/ (.*//' | sort
echo "--bench-failure markers (non-empty = 台面故障, never a reading)--"
grep -n '0xc0000135\|build failed\|setup failed\|no test files\|cannot find module\|creating work dir\|No such file' "$L" | head -5
exit $RC
