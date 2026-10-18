#!/bin/bash
# 227-v1: final numeric corrections inside the verdict (roster 6→7 commits, union 90→96,
# prefix 14→15, foreign in-flight list). Each replacement is anchored on a unique token and
# verified by a before/after grep; ⛔ rewriting any committed judgement, only its readings.
cd '/d/work/workspace/projects plans/Wisp' || exit 2
v=docs/evidence/s1/227-v1-verdict.md
out=.scratch/wisp/probes/227/v1/logs/31-verdict-numeric-corrections.txt
{
  echo "== before (hits of each stale token) =="
  grep -o '同前缀终量 \*\*14 笔\*\*＝本腿 6＋外腿 8' "$v" | sed 's/^/  /'
  grep -o '本腿 6 笔并集 \*\*90 枚路径\*\*' "$v" | sed 's/^/  /'
  grep -o '终量＝\*\*14 笔＝本腿 6＋外腿 8\*\*' "$v" | sed 's/^/  /'
  grep -o '八枚入库＋一枚在飞 `08-answer-iv-sentences.txt`' "$v" | sed 's/^/  /'
} > "$out" 2>&1
sed -i 's/同前缀终量 \*\*14 笔\*\*＝本腿 6＋外腿 8/同前缀终量 **15 笔**＝本腿 7＋外腿 8/' "$v"
sed -i 's/本腿 6 笔并集 \*\*90 枚路径\*\*⊆三处写面/本腿 7 笔并集 **96 枚路径**⊆三处写面/' "$v"
sed -i 's/终量＝\*\*14 笔＝本腿 6＋外腿 8\*\*/终量＝**15 笔＝本腿 7＋外腿 8**/' "$v"
sed -i 's/八枚入库＋一枚在飞 `08-answer-iv-sentences.txt`/八枚已自行入库＋三枚在飞未提交（`08-*`／`09-probe-dump-source`／`10-gates-four`，逐枚 mtime 见 `logs\/30-*.log`）/' "$v"
sed -i 's|★零相撞自证＝`git log -- docs/evidence/s1/227-v1-verdict.md` 与 `git log -- 票227` 各只含本腿那一笔（外腿 8 笔⛔ 碰这两枚、⛔ 碰任何 `.go`）|★零相撞自证（`logs/29-collision-check.txt`）＝外腿 8 枚路径 ⇄ 本腿 96 枚路径**交集＝0**，且"本腿每个路径的最后提交者⛔ 本腿七笔之一"＝0 枚＝本腿⛔ 卷走别腿字节|' "$v"
{
  echo
  echo "== after (each new token must hit exactly once; stale tokens must be 0) =="
  for pat in '同前缀终量 \*\*15 笔\*\*＝本腿 7＋外腿 8' '本腿 7 笔并集 \*\*96 枚路径\*\*' '终量＝\*\*15 笔＝本腿 7＋外腿 8\*\*' '八枚已自行入库＋三枚在飞未提交' '交集＝0'; do
    printf '  [%s] hits=%s\n' "$pat" "$(grep -c "$pat" "$v" || true)"
  done
  for pat in '14 笔' '90 枚路径' '一枚在飞'; do
    printf '  stale [%s] hits=%s (want 0)\n' "$pat" "$(grep -c "$pat" "$v" || true)"
  done
  echo
  echo "== size after edits: lines=$(wc -l < "$v") (cap 130) bytes=$(wc -c < "$v") =="
  echo "== slot ruler re-run: $(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' "$v" || true) (want 0) | posctrl $(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' docs/reports/HANDOVER.md || true) (want >0)"
  echo "== CR delta introduced by sed (want 0): $(tr -dc '\r' < "$v" | wc -c) CR bytes total (pristine file had 0) =="
} >> "$out" 2>&1
echo "written"
