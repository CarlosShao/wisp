#!/usr/bin/env bash
# 322-v2 红名册作差尺：本腿独立取 10-16 名册（指名那一个原始 log），⛔ 抄编排者的差集结论。
# 尺＝顶格 ^--- FAIL 的用例名（sed 去耗时、sort -u），作差用 comm，逐枚点名。
# 名册落 .scratch/wisp/probes/322/v2/roster/（临时件只建不删）。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 99
O=".scratch/wisp/probes/orch/20261016/logs/orch-322v1-ac4-cmdwisp-full.txt"
D=".scratch/wisp/probes/322/r2/logs/cwisp-full-D.log"
C=".scratch/wisp/probes/322/r2/logs/cwisp-full-C.log"
T=".scratch/wisp/probes/322/v2/roster"
mkdir -p "$T"
{
echo "### 322-v2 作差尺 $(date '+%F %T %z')"
echo "O(10-16 名册)=$O  sha=$(git hash-object "$O" | cut -c1-12)"
grep '^--- FAIL' "$O" | sed 's/^--- FAIL: //;s/ (.*//' | sort -u > "$T/o.txt"
grep '^--- FAIL' "$D" | sed 's/^--- FAIL: //;s/ (.*//' | sort -u > "$T/d.txt"
grep '^--- FAIL' "$C" | sed 's/^--- FAIL: //;s/ (.*//' | sort -u > "$T/c.txt"
sort -u "$T/d.txt" "$T/c.txt" > "$T/union-dc.txt"
echo "--- 名册枚数: O=$(wc -l < "$T/o.txt") D=$(wc -l < "$T/d.txt") C=$(wc -l < "$T/c.txt") UNION(D,C)=$(wc -l < "$T/union-dc.txt")"
echo "--- O 全集（10-16 六枚）"; cat "$T/o.txt"
echo "--- in-O-not-in-D（被今天盖住／消失的旧红候选）"; comm -23 "$T/o.txt" "$T/d.txt"; echo "枚数=$(comm -23 "$T/o.txt" "$T/d.txt" | wc -l)"
echo "--- UNION(D,C) minus O（今天多出的枚，逐名）"; comm -13 "$T/o.txt" "$T/union-dc.txt"
echo "--- D minus C（随台面摆动的枚）"; comm -13 "$T/c.txt" "$T/d.txt"
echo "--- C minus D"; comm -13 "$T/d.txt" "$T/c.txt"
echo "--- SKIP 名册三台对撞（顶格 ^--- SKIP）"
grep '^--- SKIP' "$O" | sed 's/^--- SKIP: //;s/ (.*//' | sort > "$T/skip-o.txt"
grep '^--- SKIP' "$D" | sed 's/^--- SKIP: //;s/ (.*//' | sort > "$T/skip-d.txt"
grep '^--- SKIP' "$C" | sed 's/^--- SKIP: //;s/ (.*//' | sort > "$T/skip-c.txt"
diff "$T/skip-d.txt" "$T/skip-c.txt" && echo "SKIP D⇄C 逐字相同（枚数=$(wc -l < "$T/skip-d.txt")）"
diff "$T/skip-o.txt" "$T/skip-d.txt" && echo "SKIP O⇄D 逐字相同"
echo "### done $(date '+%F %T %z')"
} | tee .scratch/wisp/probes/322/v2/logs/roster-diff.txt
