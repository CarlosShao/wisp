#!/bin/bash
# 353-v1 R23 = 三发落点互不相同的独立证明（对仓外突变树在盘字节；量的）
R=/d/tmp/wisp/20261020/353r1
V=/d/tmp/wisp/20261020/353v1
echo "###### 1) 三棵突变树 pinnedleaf 的 sha256（彼此应不同＝落点不同）＋与 HEAD/基线比对"
for k in 1 2 3; do
  printf '  mut-final-%s sha=%s\n' "$k" "$(sha256sum $R/mut-final-$k/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-16)"
done
echo "  finalL(还原后) sha=$(sha256sum $R/finalL/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-16)"
echo "  finalB(基线)   sha=$(sha256sum $R/finalB/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-16)"
echo "###### 2) 每棵树相对 finalL 的差异行（＝那一发落的刀）"
for k in 1 2 3; do
  echo "  -- mut-final-$k vs finalL --"
  diff -u0 $R/finalL/internal/winsec/pinnedleaf_342_windows.go $R/mut-final-$k/internal/winsec/pinnedleaf_342_windows.go | grep -E '^[+-][^+-]' | sed -e 's/^\(.\{120\}\).*/\1/' | sed 's/^/     /'
done
echo "###### 3) 三把新锚在各自那棵树命中＝1、在别的树＝0（落点⛔ 同）"
declare -a A=( '	_ = base; _ = now; _ = pinned // MUT-1: removal step removed' '	f, err := os.OpenFile(full, flags, perm) // MUT-2: by-name open (baseline shape)' '	if err != nil && !os.SameFile(pinned, now) { // MUT-3: 乙 identity re-check neutered' )
i=1
for a in "${A[@]}"; do
  line="  MUT-$i anchor hits:"
  for k in 1 2 3; do
    line="$line mut$k=$(grep -cF "$a" $R/mut-final-$k/internal/winsec/pinnedleaf_342_windows.go || true)"
  done
  line="$line finalL=$(grep -cF "$a" $R/finalL/internal/winsec/pinnedleaf_342_windows.go || true)"
  echo "$line"
  i=$((i+1))
done
echo "###### 4) 旧批（笔⑤之前 sh/09）的底片四数，供对照（只报数⛔ 当终量凭据）"
L=$R/logs
for f in mut-1.raw.txt mut-2.raw.txt mut-3.raw.txt; do
  printf '  %-16s RUN=%s PASS=%s FAIL=%s SKIP=%s\n' "$f" "$(grep -c '^=== RUN' $L/$f||true)" "$(grep -c '^--- PASS' $L/$f||true)" "$(grep -c '^--- FAIL' $L/$f||true)" "$(grep -c '^--- SKIP' $L/$f||true)"
done
echo "###### 5) mut-1/2/3 旧批分子行（对照用）"
for f in mut-1.raw.txt mut-2.raw.txt mut-3.raw.txt; do
  echo "  -- $f"
  grep -E '353r1 AC#2 over' "$L/$f" | cut -c1-200 | sed 's/^/     /'
done
