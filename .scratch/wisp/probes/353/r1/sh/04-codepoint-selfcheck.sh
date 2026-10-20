#!/bin/sh
# 353-r1 起手自查：新仪器件的非 ASCII 码位 ⇄ 既有 winsec .go 的码位先例 ⇄ d22scan 的豁免形状。
# 尺＝LC_ALL=C grep -nP '[^\x00-\x7F]'（命中＝非 ASCII 行数）；⛔ 把这条尺当判据＝只是自查。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "== 本腿新件的非 ASCII 行 =="
LC_ALL=C grep -nP '[^\x00-\x7F]' internal/winsec/nametopup_353_windows_test.go || echo "(none)"
echo "== 既有 winsec .go 的非 ASCII 行数（先例）=="
for f in internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/pinnedleaf_342_windows.go internal/winsec/winsec.go internal/winsec/pinnedleaf_342_unix.go internal/winsec/leafrace_342r2_windows_test.go internal/winsec/winsec_other.go; do
  n=$(LC_ALL=C grep -cP '[^\x00-\x7F]' "$f" || true)
  echo "$n  $f"
done
echo "== 非 ASCII 字符集合（既有件里出现过的高位码位）=="
LC_ALL=C grep -hoP '[^\x00-\x7F]' internal/winsec/*.go | sort -u | while read -r ch; do
  printf '%s' "$ch" | od -An -tx1 | tr -d ' \n'
  printf '  %s\n' "$ch"
done | head -40
echo "== d22scan 的豁免形状 =="
grep -n 'TrimSpace\|HasPrefix.*//\|strip\|comment' tools/d22scan/main.go | head -20
