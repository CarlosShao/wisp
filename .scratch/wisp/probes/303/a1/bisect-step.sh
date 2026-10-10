#!/bin/sh
# 303-a1 R4 台件（票 303 AC#1 bisect 的一步）——仓外 clone 里跑最小复现集并分类。
# 判据（票面 AC#1 + 派单"本腿最容易做错的地方"）：
#   BAD(1)   ＝ 逐字 no report "ac14r-0" from the page within 15s 且 ^--- FAIL: <名>
#   GOOD(0)  ＝ ^--- PASS: TestAC14AwaitedBindingReplyReachesThePage
#   INVALID(125) ＝ 其余一切（build failed / no tests to run / 0xc0000135 / 别的红句）
# PATH 铺的是【母仓绝对路径】的两枚目录（clone 里没有 third_party/ 与 build/，被 .gitignore 挡着）：
#   /d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx
#   /d/work/workspace/projects plans/Wisp/build
# 形式＝shell 形式（/d/...），⛔ D:/ 形式：见 scripts/wisp-cli-tests.sh:113-118 的实测（D:/ 那形 0xc0000135）。

CLONE="$HOME/wisp-303-bisect"
MOTHER='/d/work/workspace/projects plans/Wisp'
STEPDIR="$HOME/wisp-303-steps"
LOG="$STEPDIR/index.tsv"
NAME='TestAC14AwaitedBindingReplyReachesThePage'
RED='no report "ac14r-0" from the page within 15s'

mkdir -p "$STEPDIR"
[ -f "$LOG" ] || printf 'step\tsha\tverdict\tgo_rc\treason\tbody\n' > "$LOG"

cd "$CLONE" 2>/dev/null || { echo "303-step: FATAL: cannot cd into $CLONE"; exit 125; }
sha=$(git rev-parse HEAD)
n=$(wc -l < "$LOG")
body="$STEPDIR/step$(printf '%03d' "$n")-$sha.txt"

export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"
go test -count=1 -timeout 300s -v -run "$NAME" ./cmd/wisp/ > "$body" 2>&1
rc=$?

reason=""
verdict=""
code=125
if grep -qF "$RED" "$body" && grep -qE "^--- FAIL: $NAME" "$body"; then
    verdict="BAD"; code=1; reason="ticket red sentence verbatim"
elif grep -qE "^--- PASS: $NAME" "$body"; then
    verdict="GOOD"; code=0; reason="page reply hop reached (--- PASS)"
else
    verdict="INVALID"; code=125
    if grep -q '0xc0000135' "$body"; then reason="STATUS_DLL_NOT_FOUND 0xc0000135 (load-time death, cases never ran)"
    elif grep -q 'build failed' "$body"; then reason="build failed"
    elif grep -q '\[compile error\]' "$body"; then reason="compile error"
    elif grep -q 'no tests to run' "$body"; then reason="no tests to run (case absent at this rev)"
    elif grep -qE "^--- FAIL: $NAME" "$body"; then reason="this case FAILed for ANOTHER reason (not the ticket red sentence)"
    elif grep -qE '^FAIL' "$body"; then reason="package FAILed, this case has no --- line"
    else reason="unclassified (go rc=$rc), no PASS line for this case"
    fi
fi

# 只留判据相关的行进索引，整发原文在 body 里
{
    printf '%s\t%s\t%s\t%s\t%s\t%s\n' "$n" "$sha" "$verdict" "$rc" "$reason" "$body"
} >> "$LOG"

echo "303-step sha=$sha verdict=$verdict go_rc=$rc reason=$reason body=$body"
exit "$code"
