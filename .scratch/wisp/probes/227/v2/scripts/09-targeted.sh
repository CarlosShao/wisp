#!/bin/bash
# 227-v2 定向复量台（串行、⛔ 并发 go test）：
#  A) mklink /J 混拼法探针＝坐实 A1053 那两枚"短路径仍红"的根因是不是 GOTMPDIR 用正斜杠拼法
#  B) 8.3 别名父段树：票 252 那一族的 1 枚残留红（Alignment 取的是仓根的【父段】短名）换名复量
#  C) 争用族 TestResolvePerCallBudget 定向 3 发（A426/A432 在册假红形）
set -u
REPO='/d/work/workspace/projects plans/Wisp'
B='/d/tmp/wisp/20261019/227v2/deliv-tree-ac5'
OUT='/d/tmp/wisp/20261019/227v2/targeted'
mkdir -p "$OUT"
GT='C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-probe2'
mkdir -p "$GT"

echo '== A) mklink 混拼法探针（backslash link ⇄ slash target，正是 A1053 final-attr 红句那形）=='
for combo in "bs bs" "bs sl" "sl bs" "sl sl"; do
  set -- $combo
  mkdir -p "$GT/$1-$2-tgt"
  LNK="$GT/$1-$2-link"; TGT="$GT/$1-$2-tgt"
  [ "$1" = sl ] && LNK="C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-probe2/$1-$2-link"
  [ "$2" = sl ] && TGT="C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-probe2/$1-$2-tgt"
  o=$(cmd //c "mklink /J $LNK $TGT" 2>&1); rc=$?
  echo "  linkform=$1 targetform=$2 rc=$rc out=$o"
done

echo '== B) 别名父段树：cp -r 到父段会出短名别名的落点，只复量票 252 那一族 =='
P2='/d/tmp/wisp/20261019/227v2-secondbench'
mkdir -p "$P2"
if [ ! -d "$P2/deliv-tree-ac5" ]; then cp -rp "$B" "$P2/deliv-tree-ac5"; fi
echo "  second tree root files=$(find "$P2/deliv-tree-ac5/internal/tools" -name '*_test.go' | wc -l)"
pwsh.exe -NoProfile -File "$REPO/.scratch/wisp/probes/227/v2/scripts/02-run.ps1" \
  -Gotmp 'C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-c' \
  -Log 'D:\tmp\wisp\20261019\227v2\targeted-252.log' \
  -Tree 'D:\tmp\wisp\20261019\227v2-secondbench\deliv-tree-ac5' \
  -Sel './internal/tools -run=TestTicket252'
bash "$REPO/.scratch/wisp/probes/227/v2/scripts/03-four-numbers.sh" "/d/tmp/wisp/20261019/227v2/targeted-252.log"

echo '== C) 争用族定向 3 发 =='
for i in 1 2 3; do
  pwsh.exe -NoProfile -File "$REPO/.scratch/wisp/probes/227/v2/scripts/02-run.ps1" \
    -Gotmp 'C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-d' \
    -Log "D:\\tmp\\wisp\\20261019\\227v2\\targeted-budget-$i.log" \
    -Tree 'D:\tmp\wisp\20261019\227v2\deliv-tree-ac5' \
    -Sel './internal/risk -run=TestResolvePerCallBudget'
  echo "  shot$i: $(grep -c '^--- PASS' "/d/tmp/wisp/20261019/227v2/targeted-budget-$i.log" || true) PASS / $(grep -c '^--- FAIL' "/d/tmp/wisp/20261019/227v2/targeted-budget-$i.log" || true) FAIL / $(grep '^RC=' "/d/tmp/wisp/20261019/227v2/targeted-budget-$i.log")"
done
echo TARGETED_DONE
