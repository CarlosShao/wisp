#!/bin/sh
# 353-r1 交件尺：逐枚 hash-object（工作树⇄入库 blob）＋ CR 只量**入库 blob**（★⛔ 对工作树 grep -c $'\r'＝钝尺）。
# 尺＝git hash-object ⇄ git cat-file blob | tr -cd '\r' | wc -c；射程＝本腿那几枚。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
mkdir -p "$L"
REF="${1:-HEAD}"
{
  echo "== REF=$REF FULL=$(git rev-parse "$REF") =="
  for f in internal/winsec/pinnedleaf_342_windows.go internal/winsec/pinnedleaf_342_unix.go \
           internal/winsec/pinnedleaf_342_fallback.go internal/winsec/winsec.go \
           internal/winsec/nametopup_353_windows_test.go; do
    wt=$(git hash-object "$f")
    ib=$(git rev-parse "$REF:$f" 2>/dev/null || echo MISSING)
    cr=$(git cat-file blob "$REF:$f" 2>/dev/null | tr -cd '\r' | wc -c)
    by=$(git cat-file blob "$REF:$f" 2>/dev/null | wc -c)
    ln=$(git cat-file blob "$REF:$f" 2>/dev/null | wc -l)
    echo "$f WORKTREE=$wt BLOB=$ib CR=$cr BYTES=$by LINES=$ln"
  done
  echo "== 承重锚对照（起手 blob ⇄ 现 blob）=="
  git rev-parse a59f5630:internal/winsec/winsec.go | sed 's/^/WINSEC_AT_BASE_ANCHOR=/'
  git rev-parse a59f5630:internal/winsec/pinnedleaf_342_windows.go | sed 's/^/PINNED_AT_BASE_ANCHOR=/'
  echo "== pristine（仓外基线树字节）对照 =="
  sha256sum /d/tmp/wisp/20261020/353r1/pristine-treeB-pinnedleaf_342_windows.go 2>/dev/null
} > "$L/hashes-$REF.txt" 2>&1
cat "$L/hashes-$REF.txt"
