#!/bin/bash
# 227-v1 AC#5 partial, BASE side: same ./internal/... scope on feac8f77 export tree,
# so "zero new reds from the delivery bytes" becomes a measured comparison (逐名红名册).
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR="$root/gotmp"; unset TMPDIR
cd "$root/base" || exit 2
go test ./internal/... -count=1 -timeout 30m -v > "$out/16-ac5-internal-base-full.txt" 2>&1
rc=$?
RUN=$(grep -c '^=== RUN' "$out/16-ac5-internal-base-full.txt" || true)
PASSTOP=$(grep -c '^--- PASS' "$out/16-ac5-internal-base-full.txt" || true)
FAILTOP=$(grep -c '^--- FAIL' "$out/16-ac5-internal-base-full.txt" || true)
SKIP=$(grep -c '^--- SKIP' "$out/16-ac5-internal-base-full.txt" || true)
{
  echo "tree=base(feac8f77) scope=./internal/... (cmd/wisp half NOT run = lane)"
  echo "rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP"
  echo "---- named FAIL roster (base) ----"
  grep '^--- FAIL' "$out/16-ac5-internal-base-full.txt" | sed 's/ (.*//' | sort > "$out/16-base-failnames.txt"
  cat "$out/16-base-failnames.txt"
} > "$out/16-ac5-internal-base.txt" 2>&1
grep '^--- FAIL' "$out/12-ac5-internal-all-full.txt" | sed 's/ (.*//' | sort > "$out/16-deliv-failnames.txt"
{
  echo "== 逐名差集 (deliv minus base) = 本票字节带进来的新红? want EMPTY =="
  comm -13 "$out/16-base-failnames.txt" "$out/16-deliv-failnames.txt"
  echo "== 差集枚数 = $(comm -13 "$out/16-base-failnames.txt" "$out/16-deliv-failnames.txt" | grep -c . || true) (want 0)"
  echo "== 反向差集 (base minus deliv) = 交付件治掉的既有红 (info only) =="
  comm -23 "$out/16-base-failnames.txt" "$out/16-deliv-failnames.txt"
  echo "== 两枚名册枚数: base=$(grep -c . "$out/16-base-failnames.txt" || true) deliv=$(grep -c . "$out/16-deliv-failnames.txt" || true) =="
} >> "$out/16-ac5-internal-base.txt" 2>&1
echo "rc=$rc"
