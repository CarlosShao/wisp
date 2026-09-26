# 票 154 AC#1 —— 触发门四枚子句的读数生成器（可复算）。
# 用法：bash .scratch/wisp/probes/154/gate-clauses.sh [锚点]   > 读数文件
# 规矩：只做 grep／不写盘上任何东西；名册变了即门响。
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
A="${1:-$(git rev-parse HEAD)}"
GO='internal/**/*.go'
GO2='cmd/**/*.go'
NT=':_test.go'

echo "# 票 154 AC#1 触发门 —— 锚点 $A"
echo "# 生成时刻 $(date -Iseconds)"
echo "# 下一程复算＝同一条命令再跑一次，逐行 diff 名册；名册变长＝门响。"

run() { # $1 标题 $2 pattern $3 额外 git-grep 旗标 $4.. pathspecs
	local title="$1" pat="$2" flags="$3"; shift 3
	echo
	echo "## $title"
	echo "\$ git grep -nE $flags '$pat' $A -- $*"
	# shellcheck disable=SC2068
	git grep -nE $flags "$pat" "$A" -- $@
	local rc=$?
	echo "# rc=$rc   （1＝零命中／0＝非空＝该枚子句已响）"
}

run "G1 生产码里构造 ToolRequest（排 internal/agent/ 自身）" \
	'ToolRequest\{' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/*'

run "G1b 生产码里任意 .Execute( 调用点（排 loop.go 与 bridge.go 自身）" \
	'\.Execute\(' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'

run "G2 OpenTask/CloseTask 的生产出现点（排 bridge.go 自身；今日应只剩 cmd/wisp/run.go 那一发＋它上面的注释）" \
	'OpenTask|CloseTask' '-w' "$GO" "$GO2" ':!*_test.go' ':!internal/tools/bridge.go'

run "G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法" \
	'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' '' internal/agent

run "G4 生产组合根打开 L0 直通（PassThroughUnclassifiedRisk 的字面赋值）" \
	'PassThroughUnclassifiedRisk:' '' '*.go' ':!*_test.go'

echo
echo "## 附：G1 的负一负（同尺摘掉排测试那一条 ⇒ 今日即非空，证明模式吃得住这个形状，不是结构上产不出读数的装饰）"
echo "\$ git grep -cE 'ToolRequest\\{' $A -- internal cmd 里 *_test.go 的枚数"
git grep -E 'ToolRequest\{' "$A" -- "$GO" "$GO2" | grep -c '_test\.go:'
echo
echo "## 附：今日 production .Execute( 全名册（不过滤）＝G1b 的排除项本体"
git grep -nE '\.Execute\(' "$A" -- "$GO" "$GO2" ':!*_test.go'
