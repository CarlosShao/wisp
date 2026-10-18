#!/bin/bash
# 339-r1 交件自证尺②③：report 绝对路径＋尺寸＋空槽尺（带正控）＋决定性读数逐字。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
L=.scratch/wisp/probes/339/r1/logs
OUT="$L/selfcheck-ruler2-3.log"
: > "$OUT"

echo "### 尺② report 绝对路径＋wc -c＋行数＋空槽尺 0 命中带正控 ###" >> "$OUT"
echo "abs=$(pwd -W)/.scratch/wisp/probes/339/r1/report.md" >> "$OUT"
echo "report bytes=$(wc -c < .scratch/wisp/probes/339/r1/report.md) lines=$(wc -l < .scratch/wisp/probes/339/r1/report.md) (帽=160)" >> "$OUT"
printf 'seed TODO FIXME 待填 <填>\n' > "$L/.ph-pos-control"
echo "positive-control(空槽词，必=1)=$(grep -cE 'TODO|FIXME|待填|<填>' "$L/.ph-pos-control")" >> "$OUT"
zero=0
for f in .scratch/wisp/probes/339/r1/report.md $L/*.log; do
  c=$(grep -cE 'TODO|FIXME|待填|<填>|XXXX|placeholder' "$f" || true)
  [ "$c" != "0" ] && { echo "HIT $f=$c" >> "$OUT"; zero=1; }
done
echo "report+logs 空槽命中=0? $([ $zero -eq 0 ] && echo YES || echo NO)" >> "$OUT"
echo "logs 0 字节枚数=$(find $L -maxdepth 1 -name '*.log' -size 0 | wc -l) (must be 0)" >> "$OUT"

echo "### 尺③ 决定性读数逐字 ###" >> "$OUT"
four() { # $1=file
  echo "$(basename "$1"): RUN=$(grep -c '^=== RUN' "$1") PASS_TOP=$(grep -c '^--- PASS' "$1") FAIL_TOP=$(grep -c '^--- FAIL' "$1") SKIP=$(grep -c '^--- SKIP' "$1") $(grep -o 'GO_TEST_RC=[0-9]*' "$1" || true)" >> "$OUT"
}
echo "-- Linux 起手⇄交件（同容器配方）--" >> "$OUT"
four "$L/linux-ac0-fullpkg.log"; four "$L/linux-final-fullpkg.log"
echo "-- Windows 起手⇄交件（原生 PowerShell）--" >> "$OUT"
four "$L/win-ac0-fullpkg-psnative.log"; four "$L/win-final-fullpkg-psnative.log"
echo "-- 去牙反控（Windows，三发串行的四数）--" >> "$OUT"
four "$L/win-ac2-positive-unmutated.log"; four "$L/win-ac2-deteeth-fallback.log"; four "$L/win-ac2-positive-after-restore.log"
echo "-- 去牙 Linux 半格＋绊线敏感度 --" >> "$OUT"
four "$L/linux-ac2-deteeth.log"; grep -E 'RC=' "$L/tripwire-sensitivity.log" >> "$OUT"
echo "-- 去牙两枚红句逐字（final 行号）--" >> "$OUT"
grep -E 'anchor_miss_reason_334_test.go:(152|177|214|218|222)' "$L/win-ac2-deteeth-fallback.log" >> "$OUT"
echo "-- 还原哈希三处 --" >> "$OUT"
grep -oE '(pristine_now|MUT_BEGIN|restored_hash|final_hash)=?[0-9a-f]{40}.*(YES)?' "$L/win-ac2-deteeth-fallback.log" >> "$OUT" 2>/dev/null
echo "blacklist.go 起手 blob＝还原＝c1799cd24afce19e6accef23db23a7d7ee746dec（尺＝git rev-parse e6fd0f50:internal/risk/blacklist.go ⇄ git hash-object 导出树，件内 P0/MUT_BEGIN/restored/final 四处 match=YES 见 stdout 记录）" >> "$OUT"
echo "-- 门禁 rc --" >> "$OUT"
grep -E '^(gofumpt|d22|len|vet)_rc=|VERDICT' "$L/99-final-gates.log" >> "$OUT"
cat "$OUT"
