#!/usr/bin/env bash
# 322-v2 estate⇄整包两台包含尺（派单③／必答(v)：estate 三档只有 D 台 ⇒ 本腿用名册尺核
# "缺 C 台那一维"是否已被整包两台覆盖；⛔ 推理代替读数）。
# 尺＝定向发顶格判定名（去耗时、sort -u）⊂ 整包两台各自的顶格判定名，comm 逐枚。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 99
L=".scratch/wisp/probes/322/r2/logs"
R3B=".scratch/wisp/probes/228/r3b/logs"
O16=".scratch/wisp/probes/orch/20261016/logs/orch-322v1-ac4-cmdwisp-full.txt"
T=".scratch/wisp/probes/322/v2/roster"
mkdir -p "$T"
{
echo "### 322-v2 estate 包含尺 $(date '+%F %T %z')"
grep -E '^--- (PASS|FAIL|SKIP):' "$L/cwisp-estate-resident.log"    | sed 's/^--- [A-Z]*: //;s/ (.*//' | sort -u > "$T/est-res.txt"
grep -E '^--- (PASS|FAIL|SKIP):' "$L/cwisp-estate-deathpoints.log" | sed 's/^--- [A-Z]*: //;s/ (.*//' | sort -u > "$T/est-dp.txt"
grep -E '^--- ' "$L/cwisp-full-D.log" | sed 's/^--- [A-Z]*: //;s/ (.*//' | sort -u > "$T/verd-d.txt"
grep -E '^--- ' "$L/cwisp-full-C.log" | sed 's/^--- [A-Z]*: //;s/ (.*//' | sort -u > "$T/verd-c.txt"
echo "-- estate-resident 判定名枚数=$(wc -l < "$T/est-res.txt") not-in-fullD=$(comm -23 "$T/est-res.txt" "$T/verd-d.txt" | wc -l) not-in-fullC=$(comm -23 "$T/est-res.txt" "$T/verd-c.txt" | wc -l)"
comm -23 "$T/est-res.txt" "$T/verd-d.txt"
comm -23 "$T/est-res.txt" "$T/verd-c.txt"
echo "-- deathpoints 两枚逐名在两台的存在尺："
for t in TestPanelThreadIsSTAAndExitsCleanly TestAC13BringUpSurvivesAReusedThreadQuit; do
  printf "%-50s D=%s C=%s\n" "$t" "$(grep -c "^--- .*: $t " "$L/cwisp-full-D.log"||true)" "$(grep -c "^--- .*: $t " "$L/cwisp-full-C.log"||true)"
done
echo "-- resident 族顶格 RUN/判定/FAIL 尺（『从不执行』旧口径的收窄凭据）："
for f in "$O16" "$L/cwisp-full-D.log" "$L/cwisp-full-C.log"; do
  printf "%-34s topRUNres=%s topVERDres=%s resFAIL=%s\n" "$(basename "$f")" \
    "$(grep '^=== RUN   ' "$f" | grep 'Resident' | grep -vc '/' || true)" \
    "$(grep -cE '^--- [A-Z]*:.*Resident' "$f"||true)" \
    "$(grep -cE '^--- FAIL:.*Resident' "$f"||true)"
done
echo "-- 228/r3b 两台（旧口径出处的死亡台面）同尺："
for f in "$R3B/surface-c.log" "$R3B/surface-d.log"; do
  printf "%-16s topRUNres=%s topVERDres=%s lastRUN=%s\n" "$(basename "$f")" \
    "$(grep '^=== RUN   ' "$f" | grep 'Resident' | grep -vc '/' || true)" \
    "$(grep -cE '^--- [A-Z]*:.*Resident' "$f"||true)" \
    "$(grep '^=== RUN   ' "$f" | tail -1 | cut -c12-64)"
done
echo "### done $(date '+%F %T %z')"
} | tee .scratch/wisp/probes/322/v2/logs/estate-dC-rulers.txt
