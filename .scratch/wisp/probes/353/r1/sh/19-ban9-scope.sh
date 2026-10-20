#!/bin/sh
# 353-r1 ban#9 射程尺＝d22scan 的 phantom-citation 那一道到底扫不扫 _test.go；＋列出我这几枚里所有
# 以 internal/ 起头的注释路径 token（＝会进 ban9 的 Stat 判定的形状）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "== ban9 调用点：walk 的建档范围 =="
  grep -n 'ban9\|walkComments\|_test\.go' tools/d22scan/main.go | sed -n '1,40p'
  echo
  echo "== ban9 那个 for 之前有没有 test 豁免 =="
  awk 'NR>=770 && NR<=805' tools/d22scan/main.go
  echo
  echo "== 我这四枚件里以 internal/ 起头的 token（含行号）=="
  grep -n 'internal/syscall' internal/winsec/pinnedleaf_342_windows.go internal/winsec/pinnedleaf_342_unix.go \
    internal/winsec/pinnedleaf_342_fallback.go internal/winsec/winsec.go internal/winsec/nametopup_353_windows_test.go \
    internal/winsec/ancestorwindow_342r2_windows_test.go 2>/dev/null || echo "(none)"
  echo
  echo "== 正控：仓内既有件里有没有同款拼法（若⛔ 存＝说明这是本腿新引入的形状）=="
  git grep -n 'internal/syscall' -- 'internal/**/*.go' 'cmd/**/*.go' ':!internal/winsec/nametopup_353*' | head -10 || echo "(none besides mine)"
} > "$L/ban9-scope.txt" 2>&1
cat "$L/ban9-scope.txt"
