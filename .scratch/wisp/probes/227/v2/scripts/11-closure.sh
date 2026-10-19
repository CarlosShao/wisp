#!/bin/bash
# 两把收尾尺：
#  A) 编排者那发的 cmd/wisp 块到底跑了几枚（截断定量）⇄ 我这发 ⇄ 源码里的用例总数
#  B) 本腿两发之间摆动的那枚（full-A 红 ⇄ full-B 绿）安静定向复量 3 发
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ORCH=/d/tmp/wisp/20261019/orch227-b/full.log
MINE=/d/tmp/wisp/20261019/227v2/full-B.log
echo '== A) cmd/wisp 块内 === RUN 枚数 =='
for pair in "orch:$ORCH" "mine:$MINE"; do
  tag=${pair%%:*}; L=${pair#*:}
  start=$(grep -n '^=== RUN   TestAC1AlwaysBranchDoesNotRevertAHandEditedKey$' "$L" | head -1 | cut -d: -f1)
  end=$(grep -n '^FAIL	github.com/CarlosShao/wisp/cmd/wisp' "$L" | head -1 | cut -d: -f1)
  n=$(awk -v a="$start" -v b="$end" 'NR>=a && NR<=b && /^=== RUN/' "$L" | wc -l)
  p=$(awk -v a="$start" -v b="$end" 'NR>=a && NR<=b && /^--- PASS/' "$L" | wc -l)
  f=$(awk -v a="$start" -v b="$end" 'NR>=a && NR<=b && /^--- FAIL/' "$L" | wc -l)
  echo "  $tag block lines $start..$end RUN=$n PASS_top=$p FAIL_top=$f"
  echo "  $tag last RUN in block: $(awk -v a="$start" -v b="$end" 'NR>=a && NR<=b && /^=== RUN/' "$L" | tail -1)"
done
echo '  源码 cmd/wisp 顶层用例总数: '
cd "$REPO"; grep -h '^func Test' cmd/wisp/*_test.go | wc -l
echo '== B) 摆动名安静定向 3 发 =='
for i in 1 2 3; do
  pwsh.exe -NoProfile -File "$REPO/.scratch/wisp/probes/227/v2/scripts/02-run.ps1" \
    -Gotmp 'C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-e' \
    -Log "D:\\tmp\\wisp\\20261019\\227v2\\swing-winsec-$i.log" \
    -Tree 'D:\tmp\wisp\20261019\227v2\deliv-tree-ac5' \
    -Sel './internal/winsec -run=TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1'
  echo "  winsec shot$i: $(grep -c '^--- PASS' "/d/tmp/wisp/20261019/227v2/swing-winsec-$i.log" || true)P / $(grep -c '^--- FAIL' "/d/tmp/wisp/20261019/227v2/swing-winsec-$i.log" || true)F $(grep '^RC=' "/d/tmp/wisp/20261019/227v2/swing-winsec-$i.log")"
done
echo DONE
