#!/bin/bash
# 353-v1 R13 = 去牙三发的正控⇄平台冻结逐字⇄还原凭据（量的）
L=/d/tmp/wisp/20261020/353r1/logs
R=/d/tmp/wisp/20261020/353r1
echo "###### mut-1 正控逐字（腿报 name=present）"
grep -E 'AC#2 CONTROL' "$L/final-mut-1.raw.txt" | sed -e 's/^\(.\{330\}\).*/\1/'
echo "###### mut-2 腿报逐字红句（平台把名冻住）"
grep -nE 'being used by another process|did not ENTER|froze' "$L/final-mut-2.raw.txt" | sed -e 's/^\(.\{300\}\).*/\1/' | sed -n '1,10p'
echo "###### mut-3 分子逐字"
grep -E 'over 3 shots' "$L/final-mut-3.raw.txt" | sed -e 's/^\(.\{260\}\).*/\1/'
echo "###### 日志里有没有腿自己的还原记录"
grep -rl 'RESTORE_SHA256' "$L" 2>/dev/null | sed 's/^/   /'
grep -rh 'RESTORE_SHA256' "$L" 2>/dev/null | sed 's/^/   /'
echo "###### 交件侧三枚关键件在盘哈希（＝还原后终态，我自己量）"
for t in finalB finalL; do
  printf '%-8s pinnedleaf %s\n' "$t" "$(sha256sum $R/$t/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-16)"
done
printf 'HEADblob pinnedleaf %s\n' "$(cd '/d/work/workspace/projects plans/Wisp' && git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go | sha256sum | cut -d' ' -f1 | cut -c1-16)"
printf 'BASEblob pinnedleaf %s\n' "$(cd '/d/work/workspace/projects plans/Wisp' && git cat-file blob a59f5630:internal/winsec/pinnedleaf_342_windows.go | sha256sum | cut -d' ' -f1 | cut -c1-16)"
echo "###### mut-final-* 三棵突变树现在的 pinnedleaf 哈希（只建不删，留证）"
for k in 1 2 3; do
  printf 'mut-final-%s %s\n' "$k" "$(sha256sum $R/mut-final-$k/internal/winsec/pinnedleaf_342_windows.go 2>/dev/null | cut -d' ' -f1 | cut -c1-16)"
done
