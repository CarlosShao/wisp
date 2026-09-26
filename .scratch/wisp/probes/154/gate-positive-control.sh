# 票 154 AC#1 —— 触发门三枚子句的正控（在仓库外的合成树上跑，不碰本仓）
#
# 规矩：临时件只建不删；本脚本自己只用 mkdir/commit，不用 rm。
# 每枚子句给两发：本仓真实锚点（应为 0 命中或固定名册）＋合成控制树（应非空＝门会响）。
set -u
CTL=/d/tmp/wisp154-gatectl
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO"
ANCHOR=$(git rev-parse HEAD)
echo "ANCHOR=$ANCHOR"

echo
echo "### G1 生产码里构造 agent.ToolRequest（排除 internal/agent/ 自身）"
echo "--- 本仓（真读，应为空）---"
git grep -nE "ToolRequest\{" "$ANCHOR" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/agent/*'
echo "rc=$? (1＝零命中)"
echo "--- 合成树（正控，应非空）---"
cd "$CTL" && git grep -nE "ToolRequest\{" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/agent/*'

echo
echo "### G2 OpenTask/CloseTask 的生产调用者名册（排除 bridge.go 自身）"
echo "--- 本仓（真读，今日只有 run.go:563 那一发 CloseTask＋它的注释）---"
cd "$REPO" && git grep -nwE "OpenTask|CloseTask" "$ANCHOR" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/tools/bridge.go'
echo "--- 合成树（正控：第二枚所有者出现时应变两行）---"
cd "$CTL" && git grep -nwE "OpenTask|CloseTask" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/tools/bridge.go'

echo
echo "### G3 Q-56 那一支落地（Loop 上出现收 task id 的导出方法）"
echo "--- 本仓（真读，应为空）---"
cd "$REPO" && git grep -nE "^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID" "$ANCHOR" -- internal/agent
echo "rc=$? (1＝零命中)"
echo "--- 合成树（正控，应非空）---"
cd "$CTL" && git grep -nE "^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID" -- internal/agent

echo
echo "### 今日名册基准（三枚子句在本仓的完整输出，下一程按此对拍）"
cd "$REPO"
echo "-- G1 --"; git grep -cE "ToolRequest\{" "$ANCHOR" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/agent/*' || echo "(无输出＝零命中)"
echo "-- G1 摘掉 !_test 之后（证明模式本身吃得住这个形状）--"
git grep -nE "ToolRequest\{" "$ANCHOR" -- 'internal/**/*.go' 'cmd/**/*.go' | grep -c "_test.go"
echo "-- 今日 production ToolRequest 构造点（含 internal/agent）--"
git grep -nE "ToolRequest\{" "$ANCHOR" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go'
