#!/bin/bash
# 353-v1 R01 = 锚⇄blob⇄逐笔名册尺（量的；射程=git log/show 元数据 + 入库 blob 哈希）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "=== HEAD ==="
git log -1 --format='%H %ad %s' --date=iso
echo "=== anchor a59f5630 ==="
git log -1 --format='%H %ad %s' --date=iso a59f5630
echo "=== baseline blobs at a59f5630 (census/r1 claim db077e30/f3c1b68e/63714a83) ==="
git rev-parse a59f5630:internal/winsec/pinnedleaf_342_windows.go
git rev-parse a59f5630:internal/winsec/winsec.go
git rev-parse a59f5630:internal/winsec/winsec_windows.go
echo "=== HEAD blobs ==="
for p in pinnedleaf_342_windows.go winsec.go winsec_windows.go nametopup_353_windows_test.go ancestorwindow_342r2_windows_test.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go; do
  printf '%s ' "$p"; git rev-parse "HEAD:internal/winsec/$p" 2>&1
done
echo "=== diff --stat a59f5630..HEAD restricted to .go ==="
git diff --name-only a59f5630 HEAD -- '*.go' | cat
echo "=== per-commit name-only for the 6 353-r1 commits ==="
for c in 3a509e50 fa395b33 110ed00b ffc3c62a ad2704d0 0b75e40b; do
  echo "-- $c"
  git show --name-only --format='%h %ad %s' --date=iso "$c" | sed -n '1,80p'
done
