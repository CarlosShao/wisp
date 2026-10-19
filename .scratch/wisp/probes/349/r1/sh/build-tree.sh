#!/bin/sh
# 349-r1 建我自己的两棵仓外导出树（⛔ 复用别腿的树）。
# 配方（公共块「台面配方」八件套 + docs/ + 手工三件）：
#   git archive <锚> go.mod go.sum internal cmd tools scripts frontend models  ＋ docs
#   build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist（都⛔ 跟踪）
# 树落 D: 卷（短名别名）；GOTMPDIR 指 D: 且先 mkdir -p。
# 用法：sh build-tree.sh <anchor> <dirname>
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ANCHOR="$1"
NAME="$2"
ROOT=/d/tmp/wisp/20261019/349r1
D="$ROOT/$NAME"
cd "$REPO" || exit 90
FULL=$(git rev-parse "$ANCHOR") || exit 91
mkdir -p "$D" || exit 92
{
echo "### build-tree anchor=$ANCHOR full=$FULL dir=$D"
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
printf 'risk_blacklist_hash='; ( cd "$D" && git hash-object internal/risk/blacklist.go )
printf 'risk_provenance_hash='; ( cd "$D" && git hash-object internal/risk/provenance.go )
printf 'risk_syncdirs_hash='; ( cd "$D" && git hash-object internal/risk/syncdirs.go )
# 临时件只建不删：中间 .tar 留着，不 rm。
echo "### done"
}
