#!/usr/bin/env bash
# 336-r1 交件侧三把尺（收件只认盘上三把尺：提交真在／件真在／自己复跑）。
set -u
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 1
TICKET=".scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md"
PREFIX='^票 336 落地腿 336-r1'

echo "### 336-r1 交件侧三把尺（$(date '+%F %T %z')）"
echo "head_now=$(git log -1 --format='%H %ad' --date=iso-strict)"

echo "## 尺①提交真在（严格前缀名册＋逐笔 name-only 枚数）"
git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do
  echo "  ${h:0:8} files=$(git show --name-only --pretty=format: "$h" | grep -c . || true)"
done
echo "  count=$(git log --all --grep="$PREFIX" --pretty='%h' | wc -l)"

echo "## 尺②件真在：写面越界普查（逐笔 ⊆ internal/risk/*_test.go ∪ probes/336/r1 ∪ 票 336）"
git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do
  bad=$(git show --name-only --pretty=format: "$h" \
    | grep -v '^internal/risk/[^/]*_test\.go$' \
    | grep -v '^\.scratch/wisp/probes/336/r1/' \
    | grep -v '^\.scratch/wisp/issues/336-' | grep -c . || true)
  echo "  ${h:0:8} out_of_scope=$bad"
done
echo "  .go 非测试命中总数=$(git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do git show --name-only --pretty=format: "$h"; done | grep '\.go$' | grep -vc '_test\.go$' || true)"

echo "## 尺③自己复跑：blacklist.go 逐枚提交 blob（期望全＝c1799cd2…＝起手同串）"
git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do
  echo "  ${h:0:8} $(git rev-parse "$h:internal/risk/blacklist.go")"
done
echo "  工作树 internal/cmd 脏行=$(git status --porcelain --untracked-files=no -- internal cmd | wc -l)（期望 0）"
echo "  禁区件逐笔命中=$(git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do git show --name-only --pretty=format: "$h"; done | grep -cE 'tokens_fourway_test|l2_grant_boundary|ticket90_persist|thresholds\.go|golden|allowlist\.txt|^go\.(mod|sum)$|ci\.yml|internal/panel/|cmd/wisp/|frontend/|design/' || true)"

echo "## 票面七格（交件侧现跑；⛔ 翻框）"
echo "  未勾=$(grep -c '^- \[ \]' "$TICKET" || true) 已勾=$(grep -c '^- \[x\]' "$TICKET" || true)"
echo "  Status 行=$(grep -m1 '^- Status:' "$TICKET" || true)"

echo "## 我⛔ 入库的 0 字节尺坏 artefact（临时件只建不删⇒留盘上未跟踪、⛔ 当凭据）"
for f in ac4-full-before-dshape-a ac4-full-before-dshape-b ac4-full-after-dshape-a ac4-full-after-dshape-b; do
  p=".scratch/wisp/probes/336/r1/logs/$f.rednames"
  [ -f "$p" ] && echo "  $f.rednames bytes=$(wc -c < "$p")"
done
p=".scratch/wisp/probes/336/r1/logs/teeth-01-pristine-dshape.log.redlist"
[ -f "$p" ] && echo "  teeth-01….log.redlist bytes=$(wc -c < "$p")"
echo "  成因具名＝引号内 glob 不展开＝尺坏；真红名册＝*.log.rednames 逐枚非 0"

echo "## 证据件尺寸（绝对路径＋wc -c＋行数）"
echo "  $(pwd)/.scratch/wisp/probes/336/r1/report.md bytes=$(wc -c < .scratch/wisp/probes/336/r1/report.md) lines=$(wc -l < .scratch/wisp/probes/336/r1/report.md)"
for f in anchor_fallback_teeth_336_windows_test.go anchor_fallback_teeth_336_other_test.go; do
  echo "  $(pwd)/internal/risk/$f bytes=$(wc -c < "internal/risk/$f") lines=$(wc -l < "internal/risk/$f")"
done
echo "  logs 枚数=$(ls .scratch/wisp/probes/336/r1/logs | wc -l) 非 0 字节枚数=$(find .scratch/wisp/probes/336/r1/logs -type f -size +0c | wc -l)"
echo "selfcheck_rc=0"
