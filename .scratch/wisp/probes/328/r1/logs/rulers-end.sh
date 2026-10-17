#!/usr/bin/env bash
# 328-r1 交件侧尺：零读数具名 + 起手侧对拉 + 引号形射程核对。
# 每把尺都带：尺文本 / 射程 / 量的 or 算的 / 正控。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
A=648eda04          # 本腿起手锚（现认）
B=f705b082          # 本腿交件侧锚（仪器③）
echo "date=$(date '+%F %T %z')  anchors: start=$A  end=$B"
echo
echo "== 尺1 起（必答②那格）：frontend 里 blockedOnApproval 命中文件数（量的） =="
echo "cmd: git grep -l 'blockedOnApproval' $B -- frontend | wc -l"
git grep -l 'blockedOnApproval' $B -- frontend | wc -l
echo "正控（同尺族）: git grep -l 'correlationId' $B -- frontend | wc -l"
git grep -l 'correlationId' $B -- frontend | wc -l
echo "正控 2（大小写不敏感同族，防我读错名）: git grep -il 'blockedonapproval' $B -- frontend | wc -l"
git grep -il 'blockedonapproval' $B -- frontend | wc -l
echo
echo "== 尺2 页面侧接收器自陈的行号（现量，⛔ 抄普查件） =="
git grep -n 'wispOnPanelSnapshot' $B -- cmd/wisp/panel_snapshot_deliver_329.go | cut -c1-160
echo
echo "== 尺3 生产包那句自陈 =="
git grep -n 'blockedOnApproval is never asserted on a production packet' $B -- cmd/wisp | cut -c1-160
echo
echo "== 尺4 名册那把尺今天的直呼者（量的，交件侧） =="
git grep -n 'splitClashSuffix' $B -- internal cmd | cut -c1-120
echo "正控（同尺族应命中）: git grep -n 'indexWaitingKey' $B -- internal | wc -l"
git grep -n 'indexWaitingKey' $B -- internal | wc -l
echo
echo "== 尺5 队列撞名改写有没有任何用例真产过它（量的，⛔ 抄普查件） =="
git grep -nE '"[^"]*#[0-9]+"' $B -- internal/agent/approval | wc -l
echo "正控（同尺族，全仓）: git grep -nE '\"[A-Za-z0-9_.-]+#[0-9]+\"' $B -- internal cmd | wc -l"
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' $B -- internal cmd | sed 's/^/   /'
echo
echo "== 尺6 cmd/wisp 那包卡片 corr 里带 # 的用例枚数（静态界，本腿⛔ 跑该包） =="
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' $B -- cmd/wisp | wc -l
echo "正控（同尺族在 internal/panel）: git grep -nE '\"[A-Za-z0-9_.-]+#[0-9]+\"' $B -- internal/panel | wc -l"
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' $B -- internal/panel | wc -l
echo
echo "== 尺7 internal/tools 里断言 BlockedOnApproval 的用例枚数（本腿没跑那包 ⇒ 静态界） =="
git grep -n 'BlockedOnApproval' $B -- internal/tools | wc -l
echo
echo "== 尺8 本腿四枚提交只碰写面（逐枚 name-only） =="
for c in b2d55342 1696d01b ffebb6d1 f705b082; do echo "-- $c"; git show --name-only --format='   %h %ad %s' $c | sed 's/^/   /'; done
echo
echo "== 尺9 callCorr 铸造点仍只一枚（交件侧复跑，⛔ 抄派单号） =="
git grep -n 'callCorr' $B -- internal cmd | grep -v '_test.go' | cut -c1-120
echo "done rc=0"
