#!/usr/bin/env bash
# 票 157 实现程 · 收口后的终检（一枚 commit 一枚，本文件只读数不结论）
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
COMMITS=$(git log --grep='157 第\|157 收口' --pretty=%h)
echo "### W1. 本程 commit（主题锚定取号，枚数只在这里量）"
echo "$COMMITS" | tr '\n' ' '; echo
printf '枚数 = '; echo "$COMMITS" | wc -l
echo "\$ git log --no-walk --pretty=tformat: --name-only <上面那几枚> | sort -u"
git log --no-walk --pretty=tformat: --name-only $COMMITS | sort -u > /tmp/roster157final.txt
cat /tmp/roster157final.txt
printf '名册去重枚数 = '; grep -c . /tmp/roster157final.txt
echo "--- 越界自检（应 rc=1、零枚）---"
grep -v -E '^(docs/evidence/s1/155-three-unjudged-cells-r1\.md|docs/evidence/s1/153-trace-lies-unguarded-r1\.md|\.scratch/wisp/probes/157/)' /tmp/roster157final.txt; echo "rc=$?"
echo
echo "### W2. 禁改面逐支（尺＝最终名册）"
for pat in '^docs/PLAN\.md' '^docs/specs/' '^internal/' '^cmd/' '^third_party/' 'thresholds\.go' 'golden' 'allowlist\.txt' '^scripts/slo-check\.ps1' '^tools/d22scan/' '^docs/reports/' '^\.scratch/wisp/issues/' '^\.scratch/wisp/dispatches/' '^docs/evidence/s1/(?!153|155)'; do
  printf '%-40s 命中 = ' "$pat"; grep -cP "$pat" /tmp/roster157final.txt
done
echo "--- 名册里 docs/evidence 那一支逐枚（应只 153/155 两枚）---"; grep '^docs/evidence/' /tmp/roster157final.txt
printf '正控：同一把尺打 b23c7f7 的 internal/ = '; git show --name-only --pretty=tformat: b23c7f7 | grep -cE '^internal/'
echo
echo "### W3. 零翻勾 / 零 -done 终检（三枚票面）"
for n in 153 155 157; do
  f=$(ls .scratch/wisp/issues/${n}-*.md 2>/dev/null | head -1)
  printf '%s 未勾 = %s ／ 已勾 = %s ／ 文件名带 -done? ' "$n" "$(grep -c '^- \[ \]' "$f")" "$(grep -c '^- \[x\]' "$f")"
  case "$f" in *-done*) echo "YES";; *) echo "NO";; esac
  printf '   最后一次改动 = '; git log -1 --format='%h %ad %s' --date=format:'%H:%M' -- "$f" | cut -c1-44
done
echo "--- 名册里 issues/ 一支的枚数（应 0＝本程一次没提交票面）---"
grep -c '^\.scratch/wisp/issues/' /tmp/roster157final.txt
echo
echo "### W4. 只追加终检：与基版逐字节 cmp（基版＝本程首枚的父）"
BASE=$(git rev-parse 0c539c3^); echo "基版＝$BASE"
for f in docs/evidence/s1/155-three-unjudged-cells-r1.md docs/evidence/s1/153-trace-lies-unguarded-r1.md; do
  git show "$BASE:$f" > /tmp/b2.txt; n=$(wc -l < /tmp/b2.txt); head -n "$n" "$f" > /tmp/p2.txt
  cmp -s /tmp/b2.txt /tmp/p2.txt && echo "$f 前 $n 行 逐字节相同" || { echo "$f 前 $n 行 不同"; cmp /tmp/b2.txt /tmp/p2.txt; }
done
printf '附录标题枚数 = '; grep -c '^## 附录 B（票 157 代记，只追加）' docs/evidence/s1/155-three-unjudged-cells-r1.md docs/evidence/s1/153-trace-lies-unguarded-r1.md | tr '\n' ' '; echo
echo
echo "### W5. 十笔落点枚数（行首锚定尺）"
printf '155 件 B.1-B.8 = '; grep -cE '^### B\.[1-8] ' docs/evidence/s1/155-three-unjudged-cells-r1.md
printf '155 件〔成立〕= '; grep -cE '^### B\.[1-8] .*〔成立〕' docs/evidence/s1/155-three-unjudged-cells-r1.md
printf '155 件〔不成立〕= '; grep -cE '^### B\.[1-8] .*〔不成立〕' docs/evidence/s1/155-three-unjudged-cells-r1.md
printf '153 件 B.9-B.10 = '; grep -cE '^### B\.(9|10) ' docs/evidence/s1/153-trace-lies-unguarded-r1.md
printf '两枚尺列枚数（155／153）= '; grep -cE '^\*\*我这一笔用的尺是什么' docs/evidence/s1/155-three-unjudged-cells-r1.md; grep -cE '^\*\*我这一笔用的尺是什么' docs/evidence/s1/153-trace-lies-unguarded-r1.md
echo
echo "### W6. 未推自证 + 别程文件未被本程碰（工作树此刻仍脏＝别家的活，本程不动）"
git status -sb | head -1
git status --porcelain -- docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md .scratch/wisp/probes/152/my152.py cmd/wisp/slo_windows.go docs/evidence/s1/154-host-id-never-closed-r1.md
echo "（上面这四枚仍应是别家的未提交态：M／M／M／M —— 本程一枚未 commit、未还原）"
printf '本程名册里出现这四枚任一枚？ '; grep -cE '152-|probes/152|slo_windows|154-host' /tmp/roster157final.txt
