#!/bin/bash
# 353-v1 R02 = 乙本体 vs 基线：入库 blob 导出到仓外 + 逐段内容锚尺（量的）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261020/353v1/blob
mkdir -p "$OUT"
for p in internal/winsec/pinnedleaf_342_windows.go internal/winsec/winsec.go internal/winsec/pinnedleaf_342_unix.go internal/winsec/pinnedleaf_342_fallback.go internal/winsec/nametopup_353_windows_test.go internal/winsec/ancestorwindow_342r2_windows_test.go; do
  b=$(basename "$p")
  git cat-file blob "HEAD:$p" > "$OUT/HEAD-$b.txt" 2>&1
  git cat-file blob "a59f5630:$p" > "$OUT/BASE-$b.txt" 2>/dev/null || echo "(absent at baseline) $p"
done
echo "=== byte/line counts (HEAD) ==="
wc -c -l "$OUT"/HEAD-*.txt | cat
echo "=== diff pinnedleaf_342_windows.go BASE..HEAD (unified 3) ==="
diff -u3 "$OUT/BASE-pinnedleaf_342_windows.go.txt" "$OUT/HEAD-pinnedleaf_342_windows.go.txt" | sed -n '1,200p'
echo "RC_DIFF_PINNED=$?"
