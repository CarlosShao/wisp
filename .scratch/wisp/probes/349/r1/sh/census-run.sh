#!/bin/sh
# 349-r1 AC#0 同族普查（我自己那把三形尺：SA 驼峰形／SB 调用形／SC 声明形，＋BEXCL 噪声类）。
# ★两侧都在【仓外导出树】里跑 ⇒ 全包名册锚死在一个 commit（共享工作树里有 1211 枚未跟踪件，
#   其中 .scratch/wisp/probes/**/**.go 一堆，会污染"全仓定义名册"＝把假阳算成定义）。
# ★正控＝起手侧那四枚名必须命中（fires／taskScopeID／disposalScope／Complete），交件侧必须 0。
# 用法：sh census-run.sh <tree> <label>
set -u
REPO='/d/work/workspace/projects plans/Wisp'
BIN=/d/tmp/wisp/20261019/349r1/census/census349r1.exe
TREE="$1"
LABEL="$2"
OUT="$REPO/.scratch/wisp/probes/349/r1/logs/census"
mkdir -p "$OUT" || exit 91
FULL="$OUT/ac0-$LABEL-full.txt"
PROD="$OUT/ac0-$LABEL-prodonly.txt"
{
echo "### AC#0 census label=$LABEL tree=$TREE"
(cd "$TREE" && "$BIN" -pkg internal/risk -roots internal,cmd,tools)
} > "$FULL" 2> "$OUT/ac0-$LABEL-full.err"
echo "rc_full=$?"
head -1 "$FULL"
grep '^shape=' "$FULL"
wc -l "$FULL" "$OUT/ac0-$LABEL-full.err"
{
echo "### AC#0 census restricted to the three production files of the ticket"
for b in blacklist.go provenance.go syncdirs.go; do
  echo "-- $b --"
  (cd "$TREE" && "$BIN" -pkg internal/risk -roots internal,cmd,tools -file "$b")
done
} > "$PROD" 2> "$OUT/ac0-$LABEL-prod.err"
echo "rc_prod=$?"
grep -nE 'fires|taskScopeID|disposalScope|(^|[^a-zA-Z])Complete\(' "$FULL" > "$OUT/ac0-$LABEL-four-names.txt" || echo "ZERO hits for the four names (delivery side wants this)"
wc -l "$OUT/ac0-$LABEL-four-names.txt" "$PROD"
echo "### counters (prod-only files)"
grep '^shape=' "$PROD"
