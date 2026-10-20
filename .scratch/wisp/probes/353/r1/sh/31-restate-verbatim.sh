#!/bin/sh
# 353-r1 旧文字重述尺＝改前（★基线锚 a59f5630 的入库 blob）⇄改后（★HEAD 的入库 blob）逐字取回，
# 尺＝git cat-file blob + awk 按内容锚取段（⛔ 绝对行号当锚＝行号会漂；这里取段是为了"逐字"，报数仍带尺名）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
B="$L/restate-before-after.txt"
{
  echo "########## 重述①＝pinnedleaf_342_windows.go 那句 'cannot delete a foreign name either' ##########"
  echo "--- 改前（a59f5630 blob，含行号）---"
  git cat-file blob a59f5630:internal/winsec/pinnedleaf_342_windows.go | grep -n 'foreign name either' -B4 -A3
  echo "--- 改后（HEAD blob）---"
  git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go | grep -n 'fd-relative, so it cannot reach' -A12
  echo
  echo "########## 重述②＝ancestorwindow_342r2_windows_test.go 那枚既有钉（红句 ⇄ RESIDUAL 分支）##########"
  echo "--- 改前：钉本体（a59f5630 blob）---"
  git cat-file blob a59f5630:internal/winsec/ancestorwindow_342r2_windows_test.go | grep -n 'half-made file survived' -B2 -A3
  echo "--- 改前：RESIDUAL 分支（a59f5630 blob）---"
  git cat-file blob a59f5630:internal/winsec/ancestorwindow_342r2_windows_test.go | grep -n 'RESIDUAL OBSERVED' -B2 -A7
  echo "--- 改后：钉本体（HEAD blob）---"
  git cat-file blob HEAD:internal/winsec/ancestorwindow_342r2_windows_test.go | grep -n 'survived the refusal carrying a body' -B10 -A4
  echo "--- 改后：两条分支（HEAD blob）---"
  git cat-file blob HEAD:internal/winsec/ancestorwindow_342r2_windows_test.go | grep -n 'AC#1-one RED: %d of %d entered worlds\|GUARD HELD' -A9
} > "$B" 2>&1
cat "$B"
