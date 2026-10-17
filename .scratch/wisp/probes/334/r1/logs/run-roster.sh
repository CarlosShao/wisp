#!/usr/bin/env bash
# 334-r1 成对读数尺（AC#3 名册 ⇄ AC#5 整包）。两发只差树字节（改前锚 ⇄ 本腿提交），
# 台面一律＝CI 那把形：GOTMPDIR 按 %LOCALAPPDATA% 组（逐字照 scripts/portable-tests.sh:152-159），
# TMP／TEMP 一字不动；mkdir -p 先建（没建＝"creating work dir" 红，长得像判据响）。
# 用法：run-roster.sh <导出树> <标签 before|after> <scope roster|full>
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT="${1:?导出树}"; TAG="${2:?标签}"; SCOPE="${3:?roster|full}"
LOGS="$REPO/.scratch/wisp/probes/334/r1/logs"
export GOPATH='D:\work\base\gopath'
ROSTER='TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock'
if [ "$SCOPE" = roster ]; then RUNARG="$ROSTER"; else RUNARG='.'; fi
out="$LOGS/ac3-$TAG-$SCOPE.log"
err="$LOGS/ac3-$TAG-$SCOPE.err.log"

gotmp_base="$LOCALAPPDATA\\wisp-test-gotmp"
gotmp_run="$gotmp_base\\portable-334r1-$TAG-$(date +%H%M%S)-$$"
win=$(cygpath -u "${gotmp_run//\\//}")
mkdir -p "$win"; echo "mkdir_gotmp_rc=$?" | tee -a "$out"
export GOTMPDIR="$gotmp_run"

{
    echo "date=$(date '+%F %T %z') tag=$TAG scope=$SCOPE tree=$ROOT"
    echo "surface: GOTMPDIR=$GOTMPDIR TMP=$TMP TEMP=$TEMP (TMP/TEMP 未改)"
    echo "risk_blob_tree=$(git -C "$REPO" hash-object "$ROOT/internal/risk/blacklist.go")"
    echo "risk_blob_anchor03f7e7a6=$(git -C "$REPO" rev-parse 03f7e7a6:internal/risk/blacklist.go)"
    echo "test_cmd=go test ./internal/risk/ -count=1 -timeout 30m -v -run '$RUNARG'"
} >>"$out"
( cd "$ROOT" && go test ./internal/risk/ -count=1 -timeout 30m -v -run "$RUNARG" ) >"$out.body" 2>"$err"
rc=$?
cat "$out.body" >>"$out"
echo "rc=$rc" >>"$out"
echo "--- 决定性 stderr 头 8 行 ---" >>"$out"; head -8 "$err" >>"$out"
RUN=$(grep -c '^=== RUN' "$out" || true)
PASS_TOP=$(grep -c '^--- PASS' "$out" || true)
FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)
SKIP=$(grep -c '^--- SKIP' "$out" || true)
echo "FOURCOUNT RUN=$RUN PASS_TOP=$PASS_TOP FAIL_TOP=$FAIL_TOP SKIP=$SKIP rc=$rc" | tee -a "$out"
grep '^--- FAIL\|^    --- FAIL' "$out" >"$out.redlist" || true
echo "red_names=$(wc -l <"$out.redlist")" >>"$out"
