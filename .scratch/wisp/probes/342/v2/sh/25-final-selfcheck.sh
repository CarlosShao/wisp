#!/bin/sh
# 342-v2 交件自证三把尺（盘上尺，⛔ "命令没报错"）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo '### 尺一＝本腿名册（主题起头锚定行首＝票 342＋票 348 终裁腿 342-v2）'
git log --all --format='%H%x09%s' > /d/tmp/wisp/20261019/342v2/subjects-now.tsv
awk -F'\t' 'index($2,"票 342＋票 348 终裁腿 342-v2")==1{print substr($1,1,8)}' /d/tmp/wisp/20261019/342v2/subjects-now.tsv > /d/tmp/wisp/20261019/342v2/my-roster.txt
echo "MY_COMMITS=$(wc -l < /d/tmp/wisp/20261019/342v2/my-roster.txt)"
while read -r h; do
  echo "-- $h files=$(git show --name-only --format= $h | grep -c . ) go_bytes=$(git show --name-only --format= $h | grep -c '\.go$' || true)"
  git show --name-only --format= $h | grep -vE '^(\.scratch/wisp/probes/34(2|8)/v[12]/|docs/evidence/s1/34(2-v2|8-v1)-verdict\.md$|\.scratch/wisp/issues/34(2|8)-)' | sed 's/^/   OUT_OF_FACE /'
done < /d/tmp/wisp/20261019/342v2/my-roster.txt
echo '### 尺二＝两枚判语件绝对路径＋字节＋行数'
for F in docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md; do
  echo "$PWD/$F  bytes=$(wc -c < "$F") lines=$(wc -l < "$F") CTA7=$(grep -cP '\xc2\xa7' "$F" || true)"
done
echo '### 写面件数／0 字节／.out'
echo "FILES_342_v2=$(find .scratch/wisp/probes/342/v2 -type f | wc -l) FILES_348_v1=$(find .scratch/wisp/probes/348/v1 -type f | wc -l)"
echo "ZERO_BYTE=$(find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 -type f -size 0 | wc -l) DOT_OUT=$(find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 -name '*.out' | wc -l)"
echo '### 尺三＝两票框终量（须＝起手：342 x=1 open=6 ⇄ 348 x=0 open=7）'
for T in .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md; do
  echo "$(basename $T | cut -c1-6) x=$(grep -c '^- \[x\]' "$T") open=$(grep -c '^- \[ \]' "$T") ls_tree_count=$(git ls-tree HEAD --name-only -- .scratch/wisp/issues/ | grep -c "$(basename $T | cut -c1-4)-")"
done
echo '### 跟踪件是否仍有未提交改动（须空）'
git status --porcelain -uno
echo '### HEAD'
git log -1 --format='%h %s' | cut -c1-70
