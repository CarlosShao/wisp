#!/bin/sh
# 303-a1 定向两发台件（票 303 AC#1 的 (a)/(b) 凭据）
# 用法：sh targeted-pair.sh <tag> <rev>
#   在 <rev> 上跑【两枚】用例（AC#14 最小复现集＋AC#13 dist 陷阱枚），整发原文落 $HOME/wisp-303-steps/pair-<tag>-<rev>.txt
#   尺的判据句只认逐字：no report "ac14r-0" from the page within 15s
# PATH 铺的是【母仓绝对路径】的两枚目录（clone 里没有 third_party/ 与 build/）：
#   /d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx
#   /d/work/workspace/projects plans/Wisp/build

CLONE="$HOME/wisp-303-bisect"
MOTHER='/d/work/workspace/projects plans/Wisp'
STEPDIR="$HOME/wisp-303-steps"
tag="$1"
rev="$2"
[ -n "$tag" ] && [ -n "$rev" ] || { echo "usage: $0 <tag> <rev>"; exit 2; }

cd "$CLONE" 2>/dev/null || { echo "FATAL: cannot cd into $CLONE"; exit 2; }
git checkout -q "$rev" || { echo "FATAL: checkout $rev failed"; exit 2; }
sha=$(git rev-parse HEAD)
body="$STEPDIR/pair-$tag-$sha.txt"

export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"
go test -count=1 -timeout 420s -v \
    -run 'TestAC14AwaitedBindingReplyReachesThePage|TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe' \
    ./cmd/wisp/ > "$body" 2>&1
rc=$?

echo "PAIR tag=$tag rev=$sha go_rc=$rc body=$body"
echo "-- verdict lines --"
grep -E '^(--- (PASS|FAIL|SKIP)|ok |FAIL|PASS)' "$body"
echo "-- 判据句 --"
grep -nF 'no report "ac14r-0" from the page within 15s' "$body"
grep -nF 'the embed resolves no entry' "$body"
echo "-- 页面自己的话（绿发才有）--"
grep -nF "page's own words" "$body"
exit "$rc"
