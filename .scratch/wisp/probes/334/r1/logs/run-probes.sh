#!/usr/bin/env bash
# 334-r1 AC#0：把只读探针拷进仓外导出树跑两把台面，读 formsOf/anchorForms/Classify/Gate 实值。
# ⛔ 共享工作树、⛔ 跟踪字节。台面一＝GOTMPDIR 指 D 卷（候选落在主目录链外）；
# 台面二＝CI 那把台面（GOTMPDIR 按 %LOCALAPPDATA% 组、TMP/TEMP 一字不动），
# 逐字照 scripts/portable-tests.sh:152-159 那个形，先 mkdir -p。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT="${1:?导出树}"
LOGS="$REPO/.scratch/wisp/probes/334/r1/logs"
export GOPATH='D:\work\base\gopath'
echo "date=$(date '+%F %T %z') root=$ROOT"
echo "probe_src_hash=$(git -C "$REPO" hash-object "$LOGS/probe-334r1_test.go.txt")"
cp "$LOGS/probe-334r1_test.go.txt" "$ROOT/internal/risk/zz_probe_334r1_test.go"
echo "probe_installed=$ROOT/internal/risk/zz_probe_334r1_test.go"
(cd "$ROOT" && go vet ./internal/risk/ >/dev/null 2>"$LOGS/ac0-03-probe-vet.log"; echo "vet_rc=$?" >>"$LOGS/ac0-03-probe-vet.log")
tail -3 "$LOGS/ac0-03-probe-vet.log"

# ---- 台面一：D 卷 GOTMPDIR ----
D_GOTMP='D:\W334R1Gotmp\probe'
mkdir -p /d/W334R1Gotmp/probe || echo "mkdir_gotmp_D_rc=NONZERO"
echo "### surface D: GOTMPDIR=$D_GOTMP TMP=$TMP(不动)"
( cd "$ROOT" && GOTMPDIR="$D_GOTMP" go test ./internal/risk/ -count=1 -timeout 30m -v \
    -run 'TestProbe334' ) >"$LOGS/ac0-04-probe-surface-D.log" 2>&1
echo "rc=$? run=$(grep -c '^=== RUN' "$LOGS/ac0-04-probe-surface-D.log") pass_top=$(grep -c '^--- PASS' "$LOGS/ac0-04-probe-surface-D.log") fail_top=$(grep -c '^--- FAIL' "$LOGS/ac0-04-probe-surface-D.log") skip=$(grep -c '^--- SKIP' "$LOGS/ac0-04-probe-surface-D.log")" | tee -a "$LOGS/ac0-04-probe-surface-D.log"

# ---- 台面二：CI 形（照 portable-tests.sh:152-159）----
if [ -n "${LOCALAPPDATA:-}" ]; then
    gotmp_base="$LOCALAPPDATA\\wisp-test-gotmp"
    gotmp_run="$gotmp_base\\portable-$(date +%Y%m%d-%H%M%S)-$$"
    mkdir -p "$(cygpath -u "${gotmp_run//\\//}")" && echo "mkdir_ci_gotmp_rc=0" || echo "mkdir_ci_gotmp_rc=NONZERO"
    export GOTMPDIR="$gotmp_run"
    echo "### surface CI-shape: GOTMPDIR=$GOTMPDIR TMP=$TMP(不动) TEMP=$TEMP(不动)"
    ( cd "$ROOT" && go test ./internal/risk/ -count=1 -timeout 30m -v \
        -run 'TestProbe334' ) >"$LOGS/ac0-05-probe-surface-CIshape.log" 2>&1
    echo "rc=$? run=$(grep -c '^=== RUN' "$LOGS/ac0-05-probe-surface-CIshape.log") pass_top=$(grep -c '^--- PASS' "$LOGS/ac0-05-probe-surface-CIshape.log") fail_top=$(grep -c '^--- FAIL' "$LOGS/ac0-05-probe-surface-CIshape.log") skip=$(grep -c '^--- SKIP' "$LOGS/ac0-05-probe-surface-CIshape.log")" | tee -a "$LOGS/ac0-05-probe-surface-CIshape.log"
fi
echo "shared_worktree_risk_status=$(git -C "$REPO" status --porcelain -- internal/risk | wc -l) (应为 0＝本腿没在共享树里写过)"
echo "done=$(date '+%F %T %z')"
