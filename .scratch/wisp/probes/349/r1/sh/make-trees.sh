#!/bin/sh
# 建 349-r1 自己的两棵仓外导出树（⛔ 复用别腿的树、⛔ 盘符根）。
# treeBase = 起手锚 7b633fa…（我第一笔之前）；treeDel = 我最后一笔交件锚（由本脚本参数给）。
# 用法：sh make-trees.sh <del-anchor-or-HEAD>
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
DEL="$1"
[ -z "$DEL" ] && DEL=HEAD
BASE=7b633fa9713b0514806c8b808ed4db91b1485cdc
ROOT=/d/tmp/wisp/20261019/349r1
LOG="$ROOT/logs/00-make-trees.log"
mkdir -p "$ROOT/logs"
{
echo "### base=$BASE del=$DEL"
sh .scratch/wisp/probes/349/r1/sh/build-tree.sh "$BASE" treeBase
sh .scratch/wisp/probes/349/r1/sh/build-tree.sh "$DEL" treeDel
echo "### 8.3 short-name proof (tree root must have an ASCII 8.3 alias)"
cmd //c "dir /X D:\\tmp\\wisp\\20261019" 2>/dev/null | sed 's/^/  /'
echo "### go file counts (want non-zero)"
for t in treeBase treeDel; do
  printf '%s go=%s\n' "$t" "$(find "$ROOT/$t" -name '*.go' -type f | wc -l)"
done
} > "$LOG" 2>&1
echo "wrote=$LOG rc=$?"
cat "$LOG"
