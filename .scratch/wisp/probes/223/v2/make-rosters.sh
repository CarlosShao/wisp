#!/bin/sh
# 223-v2：逐名红册／绿册抽取＋与 223-v1 那发的名集合对拉（原始两发全量件都在盘上）
cd "D:/work/workspace/projects plans/Wisp" || exit 9
V1=.scratch/wisp/probes/223/v1/full-v.txt
V2=.scratch/wisp/probes/223/v2/full-v.txt
OUT=.scratch/wisp/probes/223/v2
grep -E "^[[:space:]]*--- (PASS|FAIL|SKIP): " "$V2" | sed -E 's/^[[:space:]]*//' | sort > $OUT/v2-all-result-lines.txt
grep -E "^--- PASS: " "$V2" | sed -E 's/^--- PASS: //; s/ \(.*\)$//' | sort > $OUT/v2-green-roster-top.txt
grep -E "^    +--- PASS: " "$V2" | sed -E 's/^ +--- PASS: //; s/ \(.*\)$//' | sort > $OUT/v2-green-roster-sub.txt
grep -E "^--- (FAIL|SKIP): " "$V2" | sed -E 's/^--- (FAIL|SKIP): //; s/ \(.*\)$//' | sort > $OUT/v2-redskip-roster-top.txt
grep -E "^--- (PASS|FAIL|SKIP): " "$V1" | sed -E 's/^--- [A-Z]+: //; s/ \(.*\)$//' | sort -u > $OUT/v1-name-roster-top.txt
grep -E "^--- (PASS|FAIL|SKIP): " "$V2" | sed -E 's/^--- [A-Z]+: //; s/ \(.*\)$//' | sort -u > $OUT/v2-name-roster-top.txt
echo "### 只在 v2 出现的顶层结果名（＝v1 那发没有的）"
comm -13 $OUT/v1-name-roster-top.txt $OUT/v2-name-roster-top.txt
echo "### 只在 v1 出现的顶层结果名（＝v2 这发掉的）"
comm -23 $OUT/v1-name-roster-top.txt $OUT/v2-name-roster-top.txt
echo "### 计数"
wc -l $OUT/v2-green-roster-top.txt $OUT/v2-green-roster-sub.txt $OUT/v2-redskip-roster-top.txt $OUT/v1-name-roster-top.txt $OUT/v2-name-roster-top.txt
