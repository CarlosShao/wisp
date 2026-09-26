#!/usr/bin/env bash
# 票 157 实现程 · 交件形状尺：commit 名册（--no-walk）＋禁改面逐支＋票面框枚数（锚定尺）
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9

echo "### V1. 本程 commit 名册（按主题锚定枚枚取号，再 --no-walk 取文件名册）"
COMMITS=$(git log --grep='157 第' --pretty=%h)
echo "\$ git log --grep='157 第' --pretty=%h"
echo "$COMMITS"
printf '本程 commit 枚数 = '; echo "$COMMITS" | wc -l
echo "\$ git log --no-walk --pretty=tformat: --name-only $COMMITS | sort -u"
git log --no-walk --pretty=tformat: --name-only $COMMITS | sort -u > /tmp/roster157.txt
cat /tmp/roster157.txt
printf '名册去重枚数 = '; grep -c . /tmp/roster157.txt
echo "--- 对照：不带 --no-walk 会吞整仓历史（本仓今天的坑①，本程自己复到）---"
printf '不带走 --no-walk 的同参数尺 = '; git log --pretty=tformat: --name-only $COMMITS | sort -u | wc -l
echo "--- 越界自检：名册减掉本程声明的三枚写面 ---"
grep -v -E '^(docs/evidence/s1/155-three-unjudged-cells-r1\.md|docs/evidence/s1/153-trace-lies-unguarded-r1\.md|\.scratch/wisp/probes/157/)' /tmp/roster157.txt
echo "越界 rc=$?（上面对空＝零枚）"
echo
echo "### V2. 禁改面逐支（尺＝名册，不＝工作树；frontend/design 不算进宣称）"
for pat in '^docs/PLAN\.md' '^docs/specs/' '^internal/' '^cmd/' '^third_party/' 'thresholds\.go' 'golden' 'allowlist\.txt' '^scripts/slo-check\.ps1' '^tools/d22scan/' '^docs/reports/pending-and-issues\.md' '^docs/reports/HANDOVER\.md' '^docs/reports/injection-timeline\.md' '^\.scratch/wisp/issues/' '^docs/evidence/s1/152-' '^docs/evidence/s1/154-'; do
  printf '%-42s 命中 = ' "$pat"; grep -cE "$pat" /tmp/roster157.txt
done
echo "--- 名册里 docs/evidence/ 那一支到底有哪几枚（应只有本票两枚写面）---"
grep '^docs/evidence/' /tmp/roster157.txt
echo "--- 正控：同一把尺打在交付那枚 commit 上应命中 internal/ ---"
git show --name-only --pretty=tformat: b23c7f7 | grep -cE '^internal/'
echo
echo "### V3. 票面框枚数（锚定尺；工单 :39 写 4 枚、派单 :48 写 5 枚，本程现量）"
T157=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
printf '157 票面未勾框（行首锚定） = '; grep -c '^- \[ \]' $T157
printf '157 票面 AC 形状的框 = '; grep -c '^- \[ \] \*\*AC#' $T157
printf '157 票面已勾框 = '; grep -c '^- \[x\]' $T157
printf '157 票面 -done 字样 = '; grep -c '\-done' $T157
echo "\$ grep -n '4 枚框\|5 枚框' 157 票面"
grep -n '4 枚框\|5 枚框' $T157
echo "--- 对照三枚票面此刻的未勾框枚数（本程不勾，只量）---"
for f in .scratch/wisp/issues/153-*.md .scratch/wisp/issues/155-*.md $T157; do
  printf '%-24s 未勾 = ' "$(basename $f | cut -c1-3)"; grep -c '^- \[ \]' "$f"
done
echo "--- 票面最后一次改动归谁（三枚都该不是本程的）---"
for f in .scratch/wisp/issues/153-*.md .scratch/wisp/issues/155-*.md $T157; do
  printf '%-6s ' "$(basename $f | cut -c1-3)"; git log -1 --format='%h %ad %s' --date=format:'%H:%M' -- "$f" | cut -c1-46
done
echo
echo "### V4. 两枚证据件此刻的行数（只追加自证：旧节一字未改另尺见 V5）"
wc -l docs/evidence/s1/155-three-unjudged-cells-r1.md docs/evidence/s1/153-trace-lies-unguarded-r1.md
echo "--- 附录那一节各文件里几枚标题（应各 1 枚）---"
grep -c '^## 附录 B（票 157 代记，只追加）' docs/evidence/s1/155-three-unjudged-cells-r1.md docs/evidence/s1/153-trace-lies-unguarded-r1.md
echo
echo "### V5. '既有各节一字未改'的那把尺：本程首枚 commit 之前的版本 vs 此刻，截到附录之前"
BASE=$(git rev-parse 0c539c3^)
for f in docs/evidence/s1/155-three-unjudged-cells-r1.md docs/evidence/s1/153-trace-lies-unguarded-r1.md; do
  git show "$BASE:$f" > /tmp/base_$(basename $f)
  n=$(wc -l < /tmp/base_$(basename $f))
  head -n "$n" "$f" > /tmp/prefix_$(basename $f)
  if cmp -s /tmp/base_$(basename $f) /tmp/prefix_$(basename $f); then
    echo "$f 前 $n 行＝与 $BASE 逐字节相同（cmp 静默）"
  else
    echo "$f 前 $n 行 **不同**："; cmp /tmp/base_$(basename $f) /tmp/prefix_$(basename $f)
  fi
done
