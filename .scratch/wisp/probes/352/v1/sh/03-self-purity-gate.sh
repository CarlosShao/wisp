#!/bin/sh
# 352-v1 笔⑥＝本腿纯度自门尺（★撞到自己违令：笔⑤ 的 commit 正文含 1 枚被禁字面字符）
# 尺＝U+00A7 那对字节 C2 A7（grep -o）；正控＝docs/reports/HANDOVER.md 须 >0＝尺灵
# ★CR 一律对入库 blob 量＝git cat-file blob … | tr -cd '\r' | wc -c（⛔ 对工作树用 grep -c $'\r'）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B='票 352 终裁腿 352-v1'
SEC=$(printf '\xc2\xa7')
echo '== 0) 正控（尺须灵）=='
printf '  HANDOVER.md U00A7=%s (须>0)\n' "$(grep -o "$SEC" docs/reports/HANDOVER.md | grep -c . || true)"
echo '== 1) 本腿 5 笔逐枚：commit 正文枚数 ⇄ 入库 blob CR ⇄ 名=='
for c in 0593dd23 ed2bfc62 9cff42da 9e340ec8 a0bed7c0; do
  body=$(git log -1 --format='%B' $c)
  bcount=$(printf '%s' "$body" | grep -o "$SEC" | grep -c . || true)
  printf '  %-9s body_U00A7=%s  ' "$c" "$bcount"
  for f in $(git show --name-only --format='' $c); do
    printf '%s[CR=%s B=%s] ' "$(basename $f)" "$(git cat-file blob $c:$f | tr -cd '\r' | wc -c)" "$(git cat-file blob $c:$f | wc -c)"
  done
  echo
done
echo '== 2) 交付物本体（判语件 blob）枚数＝'
printf '  verdict blob U00A7=%s CR=%s 行数=%s 字节=%s\n' \
  "$(git cat-file blob a0bed7c0:docs/evidence/s1/352-v1-verdict.md | grep -o "$SEC" | grep -c . || true)" \
  "$(git cat-file blob a0bed7c0:docs/evidence/s1/352-v1-verdict.md | tr -cd '\r' | wc -c)" \
  "$(git cat-file blob a0bed7c0:docs/evidence/s1/352-v1-verdict.md | wc -l)" \
  "$(git cat-file blob a0bed7c0:docs/evidence/s1/352-v1-verdict.md | wc -c)"
echo '== 3) 名册尺（锚定本腿自己那枚前缀；含本笔之后会漂＝逐枚 %s 起头才算）=='
echo "  锚定枚数=$(git log --grep="^$B" --format='%h' | grep -c .)"
git log --format='%H%x09%s' > /d/tmp/wisp/20261020/352v1/scratch-v1subj.txt
awk -F'\t' -v P="$B" 'index($2,P)==1{c++} END{print "  awk index==1 枚数=" c+0}' /d/tmp/wisp/20261020/352v1/scratch-v1subj.txt
echo '== 4) 票 352 框终量（本腿一枚没翻）⇄票面尺=='
T=.scratch/wisp/issues/352-wrappath-errno-half-stringified-so-fs-ernotexist-cannot-punch-through-private-route.md
printf '  未勾=%s 已勾=%s 总框=%s 行数=%s Status行=%s\n' \
  "$(git cat-file blob HEAD:$T | grep -cE '^- \[ \]')" \
  "$(git cat-file blob HEAD:$T | grep -cE '^- \[x\]' || true)" \
  "$(git cat-file blob HEAD:$T | grep -cE '^- \[[ x]\]')" \
  "$(git cat-file blob HEAD:$T | wc -l)" \
  "$(git cat-file blob HEAD:$T | grep -n '^\*\*Status:' | cut -d: -f1)"
echo '== 5) 写面自证（本腿全部笔里 .go 枚数须＝0）=='
for c in $(git log --grep="^$B" --format='%H'); do git show --name-only --format='' $c; done | grep -c '\.go$' || true
