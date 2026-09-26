#!/usr/bin/env bash
# 157 AC#5 recheck — probe 09: append-only proof of my OWN evidence file against my first commit 4720e53.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=docs/evidence/s1/157-ac5-recheck-r1.md
BASE=4720e53
echo "### Y-A  pure-append proof: the first 239 lines must be byte-identical to the base version"
echo "  \$ git cat-file blob $BASE:$OUT | wc -l  => $(git cat-file blob "$BASE:$OUT" | wc -l)"
echo "  \$ wc -l $OUT (now)                     => $(wc -l < "$OUT")"
git cat-file blob "$BASE:$OUT" | sed -n '1,239p' > /tmp/157ac5-base-head.txt
sed -n '1,239p' "$OUT" > /tmp/157ac5-now-head.txt
cmp /tmp/157ac5-base-head.txt /tmp/157ac5-now-head.txt && echo "  cmp(前 239 行, 两版) => SILENT (逐字节相同 = 纯追加)" || echo "  cmp => DIFFERS (那不是纯追加)"
echo "  \$ git diff --numstat $BASE -- $OUT      => $(git diff --numstat "$BASE" -- "$OUT" | tr '\t' ' ')"
echo "  \$ git diff --numstat (whole worktree, must be only my two families):"
git diff --numstat | sed 's/^/     /'
echo "### Y-B  index right before this commit"
echo "  \$ git diff --cached --name-only (before add) =>"; git diff --cached --name-only | sed 's/^/     /'
