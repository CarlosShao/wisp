#!/bin/sh
# 349-v1 建我自己的两棵仓外导出树（⛔ 复用 349-r1/343-v1 的树）。
# 配方＝公共块「台面配方」八件套 + docs/ + 手工三件。
# 用法：sh build-tree.sh <anchor> <dirname>
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ANCHOR="$1"; NAME="$2"
ROOT=/d/tmp/wisp/20261019/349v1
D="$ROOT/$NAME"
cd "$REPO" || exit 90
FULL=$(git rev-parse "$ANCHOR") || exit 91
mkdir -p "$D" || exit 92
{
echo "### build-tree anchor=$ANCHOR full=$FULL dir=$D date=$(date -Iseconds)"
( cd "$REPO" && git archive "$FULL" go.mod go.sum internal cmd tools scripts frontend models docs ) > "$D/.archive-$FULL.tar" || echo "ARCHIVE_RC=NONZERO"
( cd "$D" && tar -xf ".archive-$FULL.tar" ) || echo "TAR_RC=NONZERO"
mkdir -p "$D/build" "$D/third_party/sherpa-onnx" "$D/frontend/dist"
cp -f "$REPO"/build/*.dll "$D/build/" 2>/dev/null || echo "dll_build_copy_rc=nonzero"
cp -f "$REPO"/third_party/sherpa-onnx/*.dll "$D/third_party/sherpa-onnx/" 2>/dev/null || echo "dll_sherpa_copy_rc=nonzero"
cp -rf "$REPO"/frontend/dist/. "$D/frontend/dist/" 2>/dev/null || echo "frontend_dist_copy_rc=nonzero"
printf 'go_files=%s\n' "$(find "$D" -name '*.go' -type f | wc -l)"
printf 'dll_files=%s\n' "$(find "$D" -name '*.dll' -type f | wc -l)"
printf 'docs_PLAN=%s\n' "$([ -f "$D/docs/PLAN.md" ] && echo yes || echo no)"
printf 'frontend_dist_files=%s\n' "$(find "$D/frontend/dist" -type f 2>/dev/null | wc -l)"
printf 'shortname_check=%s\n' "$(cmd -c 'dir /x /ad D:\tmp\wisp\20261019' 2>/dev/null | grep -c "$NAME" || true)"
for f in blacklist provenance syncdirs; do
  printf 'blob_%s_in_tree=%s\n' "$f" "$(cd "$D" && git hash-object internal/risk/$f.go)"
done
echo "### done"
}
