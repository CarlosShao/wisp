#!/usr/bin/env bash
# 编排者收件尺 v5b：三处可疑（115 票面谁翻的 / .go 变化枚数 / agent 三件套真名）。
cd "/d/work/workspace/projects plans/Wisp" || exit 9

P115=$(ls .scratch/wisp/issues/115-*.md | head -1)
echo "=== ① 票 115 的票面：谁翻的框 / 谁改的名 ==="
echo "path=$P115"
echo "-- 全历史（含改名跟踪）："
git log --format='  %h %ct %s' --follow -- "$P115" | cut -c1-120 | head -12
echo "-- HEAD 里这枚件的勾态：未勾=$(git show HEAD:"$P115" | grep -c '^- \[ \]') 已勾=$(git show HEAD:"$P115" | grep -c '^- \[x\]')"
echo "-- 工作树勾态：未勾=$(grep -c '^- \[ \]' "$P115") 已勾=$(grep -c '^- \[x\]' "$P115")"
echo "-- 工作树⇄blob 内容差（字节）：$(git show HEAD:"$P115" | cmp -s - "$P115" && echo SAME || echo DIFF)"
echo "-- blob CR=$(git cat-file blob HEAD:"$P115" | tr -cd '\r' | wc -c) 工作树 CR=$(tr -cd '\r' < "$P115" | wc -c)"
echo "-- AC#1/AC#2 那两行在 HEAD 的形："
git show HEAD:"$P115" | grep -nE '^\- \[[ x]\] AC#[0-9]' | cut -c1-100

echo
echo "=== ② 115-v1 那 6 笔里有没有动票面 ==="
git show --name-only --format='C=%h %s' 9f180275 2>/dev/null | grep -E '^(C=|.*issues/)' | cut -c1-110

echo
echo "=== ③ ec826c63 之后的 .go 变化：拆到「139-v1 当时」⇄「现在」 ==="
echo "-- ec826c63..a19d5bf2（139-v1 声称 5 枚）："
git diff --name-only ec826c63..a19d5bf2 -- '*.go'
echo -n "   枚数="; git diff --name-only ec826c63..a19d5bf2 -- '*.go' | wc -l
echo "-- a19d5bf2..HEAD 之后的 .go："
git diff --name-only a19d5bf2..HEAD -- '*.go'
echo -n "   枚数="; git diff --name-only a19d5bf2..HEAD -- '*.go' | wc -l
echo "-- ★a19d5bf2..HEAD 里哪几笔动了 .go（⛔ 用区间计数当枚数）："
git log --format='%h %s' a19d5bf2..HEAD --name-only -- '*.go' 2>/dev/null | grep -B20 '\.go$' | grep -E '^[0-9a-f]{8} |\.go$' | cut -c1-110 | head -30

echo
echo "=== ④ 139-v1 判语件里那"三件套"到底是哪三枚 ==="
grep -nE 'compress\.go|loop\.go|internal/agent/[a-z_]+\.go' docs/evidence/s1/139-v1-verdict.md | cut -c1-140 | head -20

echo
echo "=== ⑤ agent 目录里现在存在的、判语提到的那枚第三件 ==="
ls internal/agent/*.go | grep -v '_test' | cut -c1-60
echo "-- 各文件 HEAD blob 短哈希："
for f in internal/agent/compress.go internal/agent/loop.go internal/agent/types.go internal/agent/messages.go; do
  h=$(git rev-parse --short=8 HEAD:$f 2>/dev/null); [ -n "$h" ] && echo "   $h $f"
done

echo
echo "=== ⑥ 'history compressed' 那 3 处命中逐枚看（139-v1 说 1）==="
grep -rn 'history compressed' --include='*.go' internal | cut -c1-150

echo
echo "=== ⑦ 353-v1 名册正形 + 是否已交齐（我用的前缀是钝尺）==="
git log --format='%h %ct %s' | grep '^.*票 353 终裁腿 353-v1' | cut -c1-100
echo "-- 锚定/内存同号："
echo -n "   锚定="; git log --format='%s' | grep -c '^票 353 终裁腿 353-v1'
echo -n "   内存="; git log --format='%H%x09%s' | awk -F'\t' 'index($2,"票 353 终裁腿 353-v1")==1{n++} END{print n+0}'

echo
echo "=== ⑧ ★有效锚自查（我在上一轮引用过一枚不存在的哈希）==="
for h in a18f801d ec826c63 c673f67b a19d5bf2 9a0ebab3 35fffa29 3a509e50 ecffb93b fa395b33; do
  r=$(git rev-parse --verify -q "$h^{commit}" 2>/dev/null)
  echo "   $h -> ${r:-INVALID}"
done
