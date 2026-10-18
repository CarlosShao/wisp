#!/bin/bash
# 339-r1 交件自证尺①：逐枚提交号＋git show --name-only，全集⊆三处写面、.go 只在 internal/risk/** 测试面。
# 带一枚正控＝同尺打在别腿提交上必须点名>0（证明尺会开火）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
OUT=".scratch/wisp/probes/339/r1/logs/selfcheck-ruler1.log"
PAT='^(internal/risk/anchor_miss_reason_(334|339)[a-z0-9_]*_test\.go|\.scratch/wisp/probes/339/r1/|\.scratch/wisp/issues/339-)'
# ⚠ 尺的自我更正：第一稿这里写 `(334|339)_.*_test\.go`，把真身 anchor_miss_reason_334_test.go
# 判成越界（334 后面没有 "_.*_" 那一段）＝尺坏、⛔ 读成红；改正形后重打，正控仍会开火。
: > "$OUT"
ok=1
for c in 7c279453 7437ef88 149b265a 2bf64afa 5d1ae702; do
  full=$(git rev-parse --short=10 "$c")
  echo "== $full $(git log -1 --format=%s $c | cut -c1-60) ==" >> "$OUT"
  git show --name-only --format="" "$c" | sed '/^$/d' >> "$OUT"
  bad=$(git show --name-only --format="" "$c" | sed '/^$/d' | grep -vcE "$PAT" || true)
  gonon=$(git show --name-only --format="" "$c" | sed '/^$/d' | grep '\.go$' | grep -cvE '^internal/risk/[a-z0-9_]+_test\.go$' || true)
  echo "  -> out-of-face=$bad  bad-go=$gonon" >> "$OUT"
  { [ "$bad" -eq 0 ] && [ "$gonon" -eq 0 ]; } || ok=0
done
echo "ALL-5-COMMITS-FACES=$([ $ok -eq 1 ] && echo PASS || echo FAIL)" >> "$OUT"
allpaths=$( { for c in 7c279453 7437ef88 149b265a 2bf64afa 5d1ae702; do git show --name-only --format="" "$c" | sed '/^$/d'; done; } | sort -u | wc -l)
echo "UNIQUE-PATHS-ACROSS-5-COMMITS=$allpaths" >> "$OUT"
echo "=== 正控（尺必须会开火）＝同尺打 227-v1 的提交 7d9cf71a ==="
probe=$(git show --name-only --format="" 7d9cf71a | sed '/^$/d' | grep -vcE "$PAT" || true)
echo "POSITIVE-CONTROL other-leg out-of-face=$probe (must be > 0)" >> "$OUT"
cat "$OUT"
