#!/bin/bash
# 编排者放腿前／收件后的自跑尺（v4，2026-10-20 22:3x）
# 今天靠这族尺撞到／被撞到的东西逐条写进来＝⛔ 让下一颗脑袋重新踩：
#   ③ A1068 防双派＝锚定行首"枚数为真、集合为空"（旧笔主题形「321-r1 落地腿交件（票 321，…）」⇒"票 NN"不在行首）⇒判据件维度不可省
#   ④ A1069③ 仓根冻结件被重定向截成 0 字节 ⇒ 脏尺逐枚列名＋仓根 mtime 监视
#   ⑥ A1069②/A1070② CR 钝尺＝数字节⛔ 用 grep -c（它数行）、正形＝tr -cd / od，并必带已知含量正控
#   ⑦ A1069⑤ 依赖面尺＝internal/agent⇄cmd/wisp 直 import winsec；internal/models 零条直 import 但 go list -deps＝1（经 observe/statemachine 传递）
#         ⇒ 两把尺不同义（直接 import 面 ⇄ 传递闭包面）＝A1071③ 我这枚问题语义混了；措辞写"依赖闭包含 winsec"、⛔ 写"import"
#   ⑧ A1071⑥ ★票面前提新鲜度＝票面 §现量 里那句可判定断言今天还对不对（139-r1 停手击中我：五把防双派尺全过、票面"零条日志"却是过期读数）
#   ⑨ A1072① ★凡账/交件里写"已放 X 腿"，落笔前必须已有 subagent 返回 id 或盘上 -send 启动动作（＝"先存档再发"那把尺，今天咬了我自己）
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
echo "⑤ 仓根冻结件 mtime／字节监视（重定向事故探针；★字节数用 wc -c、⛔ grep -c）"
for f in go.mod go.sum .gitignore .github/workflows/ci.yml; do
  [ -f "$f" ] && printf "   %s 字节=%s mtime=%s 同HEAD=%s\n" "$f" "$(wc -c < "$f")" "$(stat -c %y "$f" | cut -c12-16)" "$([ "$(git hash-object "$f")" = "$(git rev-parse "HEAD:$f")" ] && echo Y || echo N)"
done
echo "⑥ 跨行诚实尺（剥行首注释拼流再 grep；⛔ 单枚行内 grep 当"已删"凭据）"
if [ -n "${ANCHOR:-}" ]; then
  for rev in HEAD; do
    n=$(git cat-file blob "$rev:${ANCHORFILE:-}" 2>/dev/null | awk '{sub(/^[ \t]*\/\/ ?/,""); printf "%s ", $0}' | grep -o "$ANCHOR" | wc -l)
    echo "   $rev 拼流命中=$n（行内＝$(git grep -c "$ANCHOR" "$rev" -- "${ANCHORFILE:-}" 2>/dev/null | cut -d: -f3)）"
  done
else echo "   （设 ANCHOR＝… ANCHORFILE＝… 再跑）"; fi
echo "⑦ 依赖面尺（直接 import ⇄ 传递闭包两把分开报；⛔ 混说"import"）"
if [ -n "${DEP:-}" ]; then
  for d in winsec; do
    echo "   $d 直 import=$(go list -f '{{join .Imports "\n"}}' "$DEP" 2>/dev/null | grep -c "internal/$d") 传递闭包=$(go list -deps "$DEP" 2>/dev/null | grep -c "internal/$d") 含测试闭包=$(go list -deps -test "$DEP" 2>/dev/null | grep -c "internal/$d")"
  done
else echo "   （设 DEP=./internal/models/ 再跑）"; fi
echo "⑧ ★票面前提新鲜度（把票面 §现量／立票理由里可判定的断言逐枚对 HEAD blob 复量；跑色类⛔ 在此判＝标『待有窗口的腿』）"
TT=$(ls .scratch/wisp/issues/${T}-*.md 2>/dev/null | head -1)
if [ -n "${ASSERT:-}" ]; then
  for a in "${ASSERT[@]}"; do
    echo "   断言「$a」⇄HEAD 现量命中=$(git grep -cE "$a" HEAD -- './**/*.go' 2>/dev/null | awk -F: '{s+=$3}END{print s+0}') 枚（>0 ⇔ 断言若写"零条/无"＝过期）"
  done
else echo "   ⚠ 没设 ASSERT＝这把没跑＝★⛔ 据此放腿（139-r1 那口就是这么咬的）"; fi
echo "⑨ 账号存在性尺（凡引用 A##/Q##/R## 先验在不在；编排者今天栽过空号）"
for acc in "${ACCS[@]:-}"; do [ -n "$acc" ] && echo "   $acc = $(grep -c "^## $acc" docs/reports/pending-and-issues.md)"; done
echo "⑩ 票面框尺（⛔ 写死预期值＝现量并报尺；我那句旧预期会被我自己的翻勾改掉＝A1067⑤）"
[ -n "$TT" ] && echo "   open=$(grep -cE '^- \[ \] \*\*AC#' "$TT" || true) done=$(grep -cE '^- \[x\] \*\*AC#' "$TT" || true) $(grep -m1 'Status' "$TT" | tr -d '\r' | cut -c1-40)"
echo "⑪ CPU 三发（任一 >70 ⇒ ⛔ 放）"
bad=0
for i in 1 2 3; do
  v=$(powershell -NoProfile -Command "(Get-Counter '\Processor(_Total)\% Processor Time' -SampleInterval 1).CounterSamples.CookedValue" 2>/dev/null | tr -dc '0-9.' | cut -d. -f1)
  echo "   $i: ${v:-NA}"; [ -n "$v" ] && [ "$v" -gt 70 ] && bad=$((bad+1))
done
echo "   超闸枚数=$bad ⇒ $([ "$bad" -eq 0 ] && echo '过闸，可放' || echo '⛔ 放')"
echo "⑫ 放腿后自门（＝⑨ 的姊妹）：写完"已放 X 腿"之前，这里必须有 id／或盘上有 -send 存档笔"
echo "   （手工核对：subagent 返回的 id 或 probes/orch/<日期>/txt/*-send.txt）"
