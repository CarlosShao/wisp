#!/bin/bash
# 353-v1 R03 = winsec.go / unix / fallback 三枚件 BASE(a59f5630)->HEAD 逐字 diff（量的）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261020/353v1/blob
for b in winsec.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go; do
  git cat-file blob "a59f5630:internal/winsec/$b" > "$OUT/BASE-$b.txt" 2>/dev/null
  echo "########## diff BASE..HEAD : $b ##########"
  diff -u3 "$OUT/BASE-$b.txt" "$OUT/HEAD-$b.txt"
  echo "RC_DIFF_$b=$?"
done
