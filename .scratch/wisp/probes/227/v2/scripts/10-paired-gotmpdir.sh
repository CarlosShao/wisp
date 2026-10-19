#!/bin/bash
# 成对实验＝两枚未定性的根因尺：同一棵树、同一 -run、只差 GOTMPDIR 那一枚根的【拼法】。
# 甲＝反斜杠形（出厂 portable-tests.sh 的 $LOCALAPPDATA\\wisp-test-gotmp 拼法）
# 乙＝正斜杠形（A1053 两发红句里 C:/Users/... 那形）
# 判据＝乙红⇄甲绿 ⇒ "短路径下仍红"与卷无关、与拼法有关＝台损；两形都红 ⇒ 另找根因。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
TREE='D:\tmp\wisp\20261019\227v2\deliv-tree-ac5'
RUNS='TestDirOfBaseOfCutOnThePlatformSeparator|Test20r2JunctionRefusalIsNeverAskable'
for form in backslash slash; do
  case $form in
    backslash) GT='C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-pf-bs' ;;
    slash)     GT='C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-pf-sl' ;;
  esac
  pwsh.exe -NoProfile -File "$REPO/.scratch/wisp/probes/227/v2/scripts/02-run.ps1" \
    -Gotmp "$GT" -Log "D:\\tmp\\wisp\\20261019\\227v2\\paired-$form.log" -Tree "$TREE" \
    -Sel "./internal/tools -run=$RUNS"
  echo "== $form (GOTMPDIR=$GT) =="
  bash "$REPO/.scratch/wisp/probes/227/v2/scripts/03-four-numbers.sh" "/d/tmp/wisp/20261019/227v2/paired-$form.log"
  grep -E 'refusal_20r2_windows_test.go:[0-9]+|pathshape_portable_test.go:[0-9]+' "/d/tmp/wisp/20261019/227v2/paired-$form.log" | head -4
done
echo PAIRED_DONE
