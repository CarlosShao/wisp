#!/bin/sh
# 342-v2 第 6 笔的尺件＋提交（★先出尺、尺落盘、再提交＝正文声称的枚数⇄实际枚数能逐字对回）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261019/342v2
F=.scratch/wisp/probes/342/v2
R=.scratch/wisp/probes/348/v1
cp -f $B/logs/x05-final-selfcheck.txt $F/logs/ 2>/dev/null
cp -f $B/logs/x05-final-selfcheck.txt $R/logs/ 2>/dev/null
cp -f $F/sh/25-final-selfcheck.sh $R/sh/ 2>/dev/null
cp -f $F/sh/26-commit-batch6.sh $R/sh/ 2>/dev/null
{
  echo "### ls-tree 唯一键尺（数已入库、⛔ git ls-files）"
  echo "LS_TREE_342=$(git ls-tree HEAD --name-only -- .scratch/wisp/issues/ | grep -c 'issues/342-')"
  echo "LS_TREE_348=$(git ls-tree HEAD --name-only -- .scratch/wisp/issues/ | grep -c 'issues/348-')"
  echo "### 本腿名册（提交前的笔数＝5；本笔＝第 6 笔）"
  git log --all --format='%H%x09%s' > $B/subjects-now.tsv
  awk -F'\t' 'index($2,"票 342＋票 348 终裁腿 342-v2")==1{print substr($1,1,8)}' $B/subjects-now.tsv
} | tee $F/logs/x06-ls-tree-and-precount.txt
cp -f $F/logs/x06-ls-tree-and-precount.txt $R/logs/
git add $F/logs/x05-final-selfcheck.txt $F/sh/25-final-selfcheck.sh $R/logs/x05-final-selfcheck.txt $R/sh/25-final-selfcheck.sh $F/logs/x06-ls-tree-and-precount.txt $R/logs/x06-ls-tree-and-precount.txt $F/sh/26-commit-batch6.sh $R/sh/26-commit-batch6.sh
git commit -F 'D:\tmp\wisp\20261019\342v2\cm-06.txt' -- $F/logs/x05-final-selfcheck.txt $F/sh/25-final-selfcheck.sh $R/logs/x05-final-selfcheck.txt $R/sh/25-final-selfcheck.sh $F/logs/x06-ls-tree-and-precount.txt $R/logs/x06-ls-tree-and-precount.txt $F/sh/26-commit-batch6.sh $R/sh/26-commit-batch6.sh
echo "COMMIT_RC=$?"
