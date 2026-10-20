#!/bin/sh
# 353-r1 交件名册尺（最后一把）＝三把尺同号＋逐枚提交⇄票面 numstat＋终 blob 集。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "ANCHORED=$(git log --all --grep='^票 353 落地腿 353-r1' --format='%h' | wc -l)"
  echo "UNANCHORED=$(git log --all --grep='票 353 落地腿 353-r1' --format='%h' | wc -l)"
  git log --all --format='%s' > "$L/subjects.txt"
  echo "MEMBER_INDEX1=$(awk 'BEGIN{n=0} index($0,"票 353 落地腿 353-r1")==1{n++} END{print n}' "$L/subjects.txt")"
  echo "--- 逐枚（旧→新）---"
  git log --all --grep='^票 353 落地腿 353-r1' --format='%h %ad %s' --date=format:'%H:%M' | sort -k2 | cut -c1-96
  echo "--- 票面 numstat（本腿那些笔）＝期望空 ---"
  for c in $(git log --all --grep='^票 353 落地腿 353-r1' --format='%h'); do
    git show --numstat --format='' "$c" | grep -E 'issues/353' || true
  done
  echo "(end)"
  echo "--- 终 blob 集（HEAD）---"
  for f in internal/winsec/winsec.go internal/winsec/pinnedleaf_342_windows.go internal/winsec/pinnedleaf_342_unix.go \
           internal/winsec/pinnedleaf_342_fallback.go internal/winsec/nametopup_353_windows_test.go \
           internal/winsec/ancestorwindow_342r2_windows_test.go; do
    echo "  $(git rev-parse HEAD:$f | cut -c1-12) $f CR=$(git cat-file blob HEAD:$f | tr -cd '\r' | wc -c)"
  done
  echo "--- 工作树⛔ 脏（本腿写面）---"
  git status --porcelain -- internal/winsec .scratch/wisp/probes/353
  echo "(end status)"
} 2>&1
