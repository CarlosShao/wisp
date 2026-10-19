#!/bin/sh
# 342-v2 票尾那两枚被禁码位尺（★逐枚分栏：既有字节 ⇄ 我新增字节）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
T342='.scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md'
T348='.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md'
for T in "$T342" "$T348"; do
  HEADHITS=$(git cat-file blob "HEAD:$T" | grep -cP '\xc2\xa7' || true)
  WTHITS=$(grep -cP '\xc2\xa7' "$T" || true)
  MYHITS=$(git diff -- "$T" | grep '^+[^+]' | grep -cP '\xc2\xa7' || true)
  echo "$T HEAD_HITS=$HEADHITS WORKTREE_HITS=$WTHITS MY_ADDED_LINES_HITS=$MYHITS"
done
echo "MY_BULLETS_HITS: 342=$(grep -cP '\xc2\xa7' /d/tmp/wisp/20261019/342v2/bullet-342.txt || true) 348=$(grep -cP '\xc2\xa7' /d/tmp/wisp/20261019/342v2/bullet-348.txt || true)"
echo "VERDICTS_HITS: 342=$(grep -cP '\xc2\xa7' docs/evidence/s1/342-v2-verdict.md || true) 348=$(grep -cP '\xc2\xa7' docs/evidence/s1/348-v1-verdict.md || true)"
echo "POSITIVE_CONTROL HANDOVER=$(grep -cP '\xc2\xa7' docs/reports/HANDOVER.md || true)"
