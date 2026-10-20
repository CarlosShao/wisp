#!/bin/bash
# 编排者放腿前／收件后的自跑尺（v3，2026-10-20 20:3x）
# 今天靠这族尺撞到／被撞到的东西，逐条写进注释＝⛔ 让下一颗脑袋重新踩：
#   A1068② 防双派＝锚定行首"枚数为真、集合为空"（旧笔主题形「321-r1 落地腿交件（票 321，…）」"票 NN"不在行首）⇒★③判据件维度不可省
#   A1069③ 仓根冻结件被重定向截成 0 字节 ⇒ ★④每次必逐枚列名（只报计数＝第一回就是这么漏的）＋加仓根 mtime 监视
#   A1066③/A1069② 钝尺三枚：CR 对工作树量（autocrlf）／池扫正则抽空／锚跨行注释⇒"0 命中"是换行伪证 ⇒ ★⑤拼流那把与⑥正控先行
#   本轮新坐实：internal/agent ⇄ internal/models 都 import internal/winsec ⇒ ⑦依赖面尺（写腿在 winsec 落笔＝这两包的色都不干净）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
T="${1:?票号}"; P="${2:?腿码前缀，如 票 321 落地腿 321-r1}"; shift 2
echo "== 票 $T ／ 前缀「$P」 =="
echo "① 锚定行首（只证同前缀没在飞）＝$(git log --all --grep="^$P" --oneline | wc -l)"
echo "② 非锚定子串（会吃正文引用）＝$(git log --all --grep="$P" --oneline | wc -l)"
LEG=$(printf '%s' "$P" | grep -oE '[0-9]+-[a-z0-9]+$')
echo "②′ 腿码任意位（旧笔换形的兜底）＝$(git log --all --grep="$LEG" --oneline | wc -l)"
echo "③ ★判据件维度（这票有没有活过）——不可省；不传路径＝⛔ 不许据此放腿"
if [ "$#" -gt 0 ]; then
  for f in "$@"; do
    printf "   %s HEAD有=%s 历史笔=%s\n" "$f" "$(git cat-file -e "HEAD:$f" 2>/dev/null && echo Y || echo N)" "$(git log --all --format='%h' -- "$f" | wc -l)"
    git log --all --format='     %h %ci %s' -- "$f" | head -3 | cut -c1-110
  done
else echo "   ⚠ 空＝这把没跑＝先 grep 票面找出它点名的落件再跑"; fi
echo "④ tracked 脏＝逐枚列名＋mtime（⛔ 只报计数＝go.mod 那起漏在第一回）"
git status --porcelain | grep -v '^??' | sed 's/^/   /' | cat
git status --porcelain | grep -v '^??' | sed 's/^ *..//' | while read -r f; do [ -f "$f" ] && echo "     mtime $f = $(stat -c %y "$f" | cut -c1-19)"; done
echo "⑤ 仓根冻结件 mtime 监视（重定向事故探针）"
for f in go.mod go.sum .gitignore .github/workflows/ci.yml; do
  [ -f "$f" ] && printf "   %s 字节=%s mtime=%s 同HEAD=%s\n" "$f" "$(wc -c < "$f")" "$(stat -c %y "$f" | cut -c12-16)" "$([ "$(git hash-object "$f")" = "$(git rev-parse "HEAD:$f")" ] && echo Y || echo N)"
done
echo "⑥ ★跨行诚实尺（剥行首注释拼流再 grep；⛔ 单枚行内 grep 当"已删"凭据）"
[ -n "${ANCHOR:-}" ] && for rev in HEAD; do
  n=$(git cat-file blob "$rev:${ANCHORFILE:-}" 2>/dev/null | awk '{sub(/^[ \t]*\/\/ ?/,""); printf "%s ", $0}' | grep -o "$ANCHOR" | wc -l)
  echo "   $rev 拼流命中=$n（行内＝$(git grep -c "$ANCHOR" "$rev" -- "${ANCHORFILE:-}" 2>/dev/null | cut -d: -f3)）"
done || echo "   （设 ANCHOR＝… ANCHORFILE＝… 再跑）"
echo "⑦ 依赖面尺（该包 import 了哪些在飞写腿的包＝⛔ 那些色不可信）"
if [ -n "${DEP:-}" ]; then
  for d in winsec models tools; do
    n=$(go list -deps "$DEP" 2>/dev/null | grep -c "wisp/internal/$d")
    echo "   $DEP => internal/$d 依赖=$n"
  done
else echo "   （设 DEP=./internal/agent/ 再跑）"; fi
echo "⑧ 票面框尺＋Status（⛔ 写死预期值，交件时现量＝A1067⑤）"
TT=$(ls .scratch/wisp/issues/${T}-*.md 2>/dev/null | head -1)
[ -n "$TT" ] && echo "   open=$(grep -cE '^- \[ \] \*\*AC#' "$TT" || true) done=$(grep -cE '^- \[x\] \*\*AC#' "$TT" || true) $(grep -m1 'Status' "$TT" | tr -d '\r' | cut -c1-40)"
echo "⑨ CPU 三发（任一 >70 ⇒ ⛔ 放）"
bad=0
for i in 1 2 3; do
  v=$(powershell -NoProfile -Command "(Get-Counter '\Processor(_Total)\% Processor Time' -SampleInterval 1).CounterSamples.CookedValue" 2>/dev/null | tr -dc '0-9.' | cut -d. -f1)
  echo "   $i: ${v:-NA}"; [ -n "$v" ] && [ "$v" -gt 70 ] && bad=$((bad+1))
done
echo "   超闸枚数=$bad ⇒ $([ "$bad" -eq 0 ] && echo '过闸，可放' || echo '⛔ 放')"
