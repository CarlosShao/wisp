#!/bin/sh
# 353-r1 同步尺＝把**已入库那批字节**（HEAD blob）取到仓外 staging，再喂两棵树。
# ★⛔ 把工作树当权威（公共块：完成通知一律当未验证）＝这里一律走 git cat-file blob HEAD:<path>。
set -u
R=/d/tmp/wisp/20261020/353r1
cd '/d/work/workspace/projects plans/Wisp' || exit 9
mkdir -p "$R/incoming/treeL" "$R/incoming/treeB"
echo "HEAD=$(git log -1 --format=%h)"
for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go \
         nametopup_353_windows_test.go ancestorwindow_342r2_windows_test.go; do
  git cat-file blob "HEAD:internal/winsec/$f" > "$R/incoming/treeL/$f"
  echo "  STAGED $f blob=$(git rev-parse "HEAD:internal/winsec/$f" | cut -c1-12) bytes=$(wc -c < "$R/incoming/treeL/$f")"
done
# 基线侧＝treeB2：产码四枚＝⛔ 基线锚 a59f5630 的 blob；测试两枚＝HEAD 的（★两侧同一对测试字节）
for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go; do
  git cat-file blob "a59f5630:internal/winsec/$f" > "$R/base/$f"
  echo "  BASE $f blob=$(git rev-parse "a59f5630:internal/winsec/$f" | cut -c1-12)"
done
