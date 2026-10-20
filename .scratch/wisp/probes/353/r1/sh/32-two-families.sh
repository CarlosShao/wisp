#!/bin/sh
# 353-r1 分清两枚仪器（★⛔ 把两族分子混着报＝票 348 AC#3 那条"两列⛔ 合并"的同形错）：
#   A＝ancestorunlink_348_windows_test.go 的 form=symlink-dir/leaf-symlink 行（truncated/deleted＝342⇄348 的号）
#   B＝leafrace_342r2_windows_test.go 的 AC#1-leaf DAMAGE 行（witness truncated/gone＝那枚 racer 自己的号）
# 两侧各 4 发逐字抽出。
set -u
L=/d/tmp/wisp/20261020/353r1/logs
echo "===== A 族（348/342 的分子行，⛔ 与 B 族混）====="
for t in finalB finalL; do
  for s in 1 2 3 4; do
    f="$L/swing-$t-$s.raw.txt"; [ -f "$f" ] || continue
    echo "-- $t shot=$s --"
    grep -oE 'form=(symlink-dir|leaf-symlink) attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$f" | sed 's/^/     /'
  done
done
echo "===== B 族（leafrace_342r2 自己的 DAMAGE 行）====="
for t in finalB finalL; do
  for s in 1 2 3 4; do
    f="$L/swing-$t-$s.raw.txt"; [ -f "$f" ] || continue
    echo "-- $t shot=$s --"
    grep -oE 'AC#1-leaf DAMAGE : witness truncated=[0-9]+ gone=[0-9]+ unreadable=[0-9]+' "$f" | sed 's/^/     /'
    grep -oE 'AC#1-leaf MEASURED.*' "$f" | cut -c1-150 | sed 's/^/     /'
  done
done
echo "===== B 族那枚用例的判据行（它凭什么红）====="
grep -n -A6 'func TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2' '/d/work/workspace/projects plans/Wisp/internal/winsec/leafrace_342r2_windows_test.go' | sed -n '1,2p'
sed -n '325,350p' '/d/work/workspace/projects plans/Wisp/internal/winsec/leafrace_342r2_windows_test.go'
