#!/bin/sh
# 建我自己的两棵仓外导出树：treeBase = 起手锚 29d258f9、treeDel = 交件锚 31d4a476。
# 八件套 + docs/ + 手工三件（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
LOG=/d/tmp/wisp/20261019/343v1/logs
ROOT=/d/tmp/wisp/20261019/343v1
mkdir -p $ROOT/logs

build_tree () { # $1=anchor $2=dirname
  local d=$ROOT/$2
  mkdir -p "$d"
  ( cd "$REPO" && git archive "$1" go.mod go.sum internal cmd tools scripts frontend models docs ) > "$d/.archive-$1.tar"
  ( cd "$d" && tar -xf ".archive-$1.tar" ) || return 1
  # 手工三件（未跟踪，⛔ 在 archive 里）
  mkdir -p "$d/build" "$d/third_party/sherpa-onnx" "$d/frontend/dist"
  cp -f "$REPO"/build/*.dll "$d/build/" 2>/dev/null
  cp -f "$REPO"/third_party/sherpa-onnx/*.dll "$d/third_party/sherpa-onnx/" 2>/dev/null
  cp -rf "$REPO"/frontend/dist/. "$d/frontend/dist/" 2>/dev/null
  return 0
}

{
echo "== build treeBase (29d258f9) =="
build_tree 29d258f9 treeBase && echo BUILT_OK_treeBase || echo BUILT_FAIL_treeBase
echo "== build treeDel (31d4a476) =="
build_tree 31d4a476 treeDel && echo BUILT_OK_treeDel || echo BUILT_FAIL_treeDel
echo "== non-empty proofs =="
for t in treeBase treeDel; do
  printf '%s go_files=%s\n' "$t" "$(find $ROOT/$t -name '*.go' -type f | wc -l)"
  printf '%s dlls=%s\n' "$t" "$(find $ROOT/$t -name '*.dll' -type f | wc -l)"
  printf '%s docs_PLAN=%s\n' "$t" "$([ -f $ROOT/$t/docs/PLAN.md ] && echo yes || echo no)"
  printf '%s blacklist_blob=' "$t"
  ( cd "$ROOT/$t" && git hash-object internal/risk/blacklist.go )
  printf '%s frontend_dist=%s\n' "$t" "$(find $ROOT/$t/frontend/dist -type f 2>/dev/null | wc -l)"
done
echo "== expected: treeBase blob = c1799cd2… ; treeDel blob = 62f2dddb… =="
echo "== cd form test (D:/... forward slash in git-bash) =="
( cd "$ROOT/treeDel" && pwd && ls -d internal/risk ) && echo CD_FORM_OK || echo CD_FORM_BAD
( cd "D:/tmp/wisp/20261019/343v1/treeDel" && pwd ) && echo CD_DRIVEFORM_OK || echo CD_DRIVEFORM_BAD
} > $LOG/g06-treebuild.log 2>&1
echo DONE
