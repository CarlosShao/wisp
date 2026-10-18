#!/usr/bin/env bash
# 336-r1 交件侧末读数（具名尺逐枚；★本件所在那一笔【不计入】下面那枚 count＝自指必然滞后，收件侧现跑为准）。
# 派单 AC#6＝我的探针面⛔ 那个分节号字符 ⇒ 本件⛔ 出现该字面字符，尺的 needle 用字节构造（SEC）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
T=".scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md"
R=".scratch/wisp/probes/336/r1/report.md"
W="internal/risk/anchor_fallback_teeth_336_windows_test.go"
O="internal/risk/anchor_fallback_teeth_336_other_test.go"
PREFIX='^票 336 落地腿 336-r1'
SEC=$(printf '\302\247')
echo "### 336-r1 末读数 $(date '+%F %T %z')  head=$(git log -1 --format='%H')"
echo "  SEC_bytes=$(printf '%s' "$SEC" | wc -c)（现量＝2 字节：U+00A7 的 UTF-8 形是 C2 A7；★上一版我在此写死一枚＝3 的期望＝假读数，本版以现量为准，同时证明 needle 真构造出来了、本件⛔ 出现那枚字面字符）"

echo "## 1 分节号逐枚位置（尺＝grep -n，⛔ 只看总数）"
for f in "$R" "$W" "$O" .scratch/wisp/probes/336/r1/logs/*.sh .scratch/wisp/probes/336/r1/logs/*.log; do
  n=$(grep -c "$SEC" "$f" 2>/dev/null || true)
  if [ "$n" != "0" ]; then
    echo "  $f hits=$n"
    grep -n "$SEC" "$f" | sed 's/^/      /' | cut -c1-140
  fi
done
echo "  票面 hits=$(grep -c "$SEC" "$T" || true)（逐枚行号见下：区分编排者 AC#6 判据原文 ⇄ 我的 bullet）"
grep -n "$SEC" "$T" | cut -c1-40 | sed 's/^/      /'
echo "  我的 bullet（第 27 行）hits=$(sed -n '27p' "$T" | grep -c "$SEC" || true)"

echo "## 2 占位符尺（正控先行）"
echo "  正控＝同尺指真含词件：$(git grep -c -E '占位|placeholder|TODO' HEAD -- '.scratch/wisp/probes/334/v1/logs/83-incident-flipcount.log' '.scratch/wisp/probes/334/v1/logs/95-final-selfcheck.log' | tr '\n' ' ')"
echo "  我的件：report.md=$(grep -cE '占位|placeholder|TODO|FIXME' "$R" || true) 牙=$(grep -cE '占位|placeholder|TODO|FIXME' "$W" || true) 绊线=$(grep -cE '占位|placeholder|TODO|FIXME' "$O" || true)"

echo "## 3 提交名册（严格前缀尺；本件那一笔不计入＝自指滞后）"
echo "  count=$(git log --all --grep="$PREFIX" --oneline | wc -l)"
git log --all --grep="$PREFIX" --pretty='%h %ad %s' --date=format:'%H:%M' | cut -c1-96 | sed 's/^/  /'

echo "## 4 写面越界（逐笔 name-only ⊆ 三处）"
git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do
  bad=$(git show --name-only --pretty=format: "$h" \
    | grep -v '^internal/risk/[^/]*_test\.go$' \
    | grep -v '^\.scratch/wisp/probes/336/r1/' \
    | grep -v '^\.scratch/wisp/issues/336-' | grep -c . || true)
  echo "  ${h:0:8} out_of_scope=$bad"
done
echo "  .go 非测试命中总数=$(git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do git show --name-only --pretty=format: "$h"; done | grep '\.go$' | grep -vc '_test\.go$' || true)"
echo "  禁区件命中=$(git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do git show --name-only --pretty=format: "$h"; done | grep -cE 'tokens_fourway_test|l2_grant_boundary|ticket90_persist|thresholds\.go|golden|allowlist\.txt|^go\.(mod|sum)$|ci\.yml|internal/panel/|cmd/wisp/|^frontend/|^design/' || true)"

echo "## 5 blacklist.go 逐枚提交 blob（期望全＝c1799cd2…）"
git log --all --grep="$PREFIX" --pretty='%H' | while read -r h; do echo "  ${h:0:8} $(git rev-parse "$h:internal/risk/blacklist.go")"; done
echo "  工作树＝$(git hash-object internal/risk/blacklist.go)"

echo "## 6 票面（⛔ 翻框）"
echo "  未勾=$(grep -c '^- \[ \]' "$T" || true) 已勾=$(grep -c '^- \[x\]' "$T" || true) 我的 bullet 枚数=$(grep -c '^- \[2026-10-18 12:1x' "$T" || true) 总行数=$(wc -l < "$T")"
echo "  append-only 尺（我这一条 ⇄ HEAD 前那 26 行逐字）＝$(git show "bb247168^:$T" | cmp - <(head -26 "$T") && echo 'IDENTICAL rc=0' || echo DIFFERENT)"

echo "## 7 形状尺（两枚 .go）"
for f in "$W" "$O"; do
  echo "  $f lines=$(wc -l < "$f") bytes=$(wc -c < "$f") cr=$(tr -dc '\r' < "$f" | wc -c) nonascii=$(grep -oP '[^\x00-\x7F]' "$f" | wc -l) exported_non_test=$(grep -E '^(func|type|var|const) [A-Z]' "$f" | grep -vcE '^func Test' || true) t_Skip=$(grep -c 't\.Skip' "$f" || true) tag=$(head -1 "$f")"
done
echo "  report.md lines=$(wc -l < "$R") bytes=$(wc -c < "$R") cr=$(tr -dc '\r' < "$R" | wc -c)"
echo "  logs 枚数=$(ls .scratch/wisp/probes/336/r1/logs | wc -l) 非 0 字节=$(find .scratch/wisp/probes/336/r1/logs -type f -size +0c | wc -l) 0 字节=$(find .scratch/wisp/probes/336/r1/logs -type f -size 0 | wc -l)"
echo "  0 字节件逐枚（⛔ 入库＝尺坏 artefact，临时件只建不删）:"; find .scratch/wisp/probes/336/r1/logs -type f -size 0 | sed 's/^/      /'

echo "## 8 台面字节对账（交件锚那棵 ⇄ 跟踪件）"
echo "  teeth_in_tree_final=$(git hash-object /d/tmp/wisp/20261018/336r1/tree-final/internal/risk/anchor_fallback_teeth_336_windows_test.go)"
echo "  teeth_tracked=$(git rev-parse HEAD:internal/risk/anchor_fallback_teeth_336_windows_test.go)"
echo "  blacklist_tree=$(git hash-object /d/tmp/wisp/20261018/336r1/tree-final/internal/risk/blacklist.go)"
echo "final_readings_rc=0"
