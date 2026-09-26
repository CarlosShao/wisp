# 票 154 AC#1 —— 触发门的读数生成器（可复算）。G1／G1b／G2／G3／G4 出自票 154；
# G5（OpenScope↔CloseScope 成对普查）由票 158 r1 新增，其形状照 G1–G4，另带正控与负一负两腿。
# 用法：bash .scratch/wisp/probes/154/gate-clauses.sh [锚点]   > 读数文件
# 规矩：只做 grep／不写仓里任何东西；名册变了即门响。
#   （G5 的 pair() 为求差集用了一次 bash herestring，临时件走 $TMPDIR、不落仓；
#     读数本体仍然全部由 git grep 在 $A 那棵树上产生，与工作树无关。）
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

pair() { # 成对普查（票 158 G5 用）：$1 标题 $2 开方 pattern $3 合方 pattern $4.. pathspecs
	# 与 run() 的区别：run() 的响不响看名册非空，pair() 的响不响看「开了却没在同一文件关过」。
	# 一条 git grep 表达不出成对，所以这一枚是【一把尺 ＋ 两次求差集】；
	# 锚点、pathspec 与逐行可 diff 的名册本体都照 run() 的形状。
	local title="$1" open="$2" close="$3"; shift 3
	echo
	echo "## $title"
	echo "\$ git grep -nEw '$open|$close' $A -- $*"
	# shellcheck disable=SC2068
	local roster rrc calls opens closes unpaired="" n f oc deleg
	roster="$(git grep -nEw "$open|$close" "$A" -- $@ 2>/dev/null)"; rrc=$?
	echo "# 逐行读数（名册本体，供下一程 diff；这一份含注释行）："
	if [ -n "$roster" ]; then printf '%s\n' "$roster"; else echo "#   （空）"; fi
	# 成对判据只吃调用点：纯注释行剔掉。
	# 不剔的下场是本程实测到的——它点到了票 158 r1 自己测试文件里那句解释 Mark 机理的注释，
	# 一把会因为散文而响的尺不叫读数（详证据件 §2.4）。名册本体照旧带注释，diff 不受影响。
	calls="$(printf '%s\n' "$roster" | grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' || true)"
	opens="$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | sort -u || true)"
	closes="$(printf '%s\n' "$calls" | grep -Ew "$close" | cut -d: -f2 | sort -u || true)"
	for f in $opens; do
		printf '%s\n' "$closes" | grep -qxF "$f" || unpaired="$unpaired$f
"
	done
	n="$(printf '%s' "$unpaired" | grep -c . || true)"
	echo "# 成对名册基准（开过、却没在同一文件关过的文件；空＝这把尺没响）："
	if [ -n "$unpaired" ]; then
		while IFS= read -r f; do
			[ -z "$f" ] && continue
			oc="$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | grep -cxF "$f" || true)"
			# 委托关闭的形状：同一文件里出现另一族的收尾动词 ⇒ 人来判「真漏」还是「委托」
			deleg="$(git grep -lEw 'CloseTask|Defer' "$A" -- "$f" 2>/dev/null || true)"
			echo "#   UNPAIRED $f (开方调用点=${oc:-0})${deleg:+  # 同文件另见收尾动词，判「真漏／委托」}"
		done <<< "$unpaired"
	else
		echo "#   （空）"
	fi
	echo "# 未成对枚数＝$n   （0＝这把尺没响／>0＝逐枚点名如上）"
	echo "# git grep rc=$rrc（1＝本射程里这一族一枚都没有——那是「射程里没这东西」，不是「都关好了」；两种 0 枚响含义不同）"
	echo "# rc=$n   （本枚子句的响＝未成对枚数；0＝不响）"
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

# G5（票 158 AC#2 新增）射程＝OpenScope↔CloseScope 的成对普查。
# G1–G4 盯的是「宿主派发工具却没关环路那一枚 id」那一族（过桥那条腿），
# 而「不过桥、直接对 risk 层开一枚 scope」这一形今天有活样板：
#   cmd/wisp/panel_assets.go:232  prov.OpenScope(taintSourceScopeID)   且 cmd/wisp 下 CloseScope 零枚。
# 排掉 internal/risk/* 的理由：OpenScope/CloseScope 的**定义**在那儿（票 158 地界明令只盘调用者），
# 把定义算成调用点会让名册每改一行注释就响一次，那是噪音不是读数。
pair "G5 OpenScope↔CloseScope 成对普查：生产码里开了 C25 scope 却没在同一文件关过的名册（排 internal/risk 的定义本体、排 *_test.go）" \
	'OpenScope' 'CloseScope' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*'

echo
echo "## G5 正控（形状照 G1–G4 的「正控先跑」：同一把尺打在已知「开又合」的样本上，必须不响）"
echo "# 样本＝internal/tools/bridge.go（:642 prov.OpenScope 与 :691 prov.CloseScope 同文件成对）"
echo "# 写法上的两条自纠（都是本程实测撞出来的，不是设想的）："
echo "#   (1) pathspec 必须与主尺同形再逐条排除。第一版图省事写成 'internal/tools/**/*.go'，"
echo "#       实测 git grep rc=1、名册空、未成对 0——那是一枚**根本产不出读数的装饰腿**，"
echo "#       它「不响」不能算正控（票 154 的 AC#3 就是拿「零区别」否掉装饰腿的）。"
echo "#       本机现量：git grep -lEw OpenScope <锚> -- internal/tools/**/*.go => rc=1 零命中；"
echo "#                          同尺 -- internal/tools/*.go  => rc=0 命中 bridge.go。"
echo "#   (2) 必须排 *_test.go。第一版没排，正控自己被本程新测试里那句解释 Mark 机理的**注释**"
echo "#       点中——那是假阳性。成对判据从此只吃调用点（见 pair() 里剔注释那一行）。"
pair "G5-正控 同尺＝主尺射程，只多排掉主尺今天点名的那一枚 panel_assets.go ⇒ 剩下的开方全是 balanced 样本" \
	'OpenScope' 'CloseScope' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*' ':!cmd/wisp/panel_assets.go'

echo
echo "## G5 负一负（同尺摘掉「排测试＋排定义」那两条 ⇒ 今日应非空，证明这把尺吃得住这个形状、不是结构上产不出读数的装饰）"
pair "G5-负一负 同尺不过滤任何文件" \
	'OpenScope' 'CloseScope' "$GO" "$GO2"

echo
echo "## 附：G1 的负一负（同尺摘掉排测试那一条 ⇒ 今日即非空，证明模式吃得住这个形状，不是结构上产不出读数的装饰）"
echo "\$ git grep -cE 'ToolRequest\\{' $A -- internal cmd 里 *_test.go 的枚数"
git grep -E 'ToolRequest\{' "$A" -- "$GO" "$GO2" | grep -c '_test\.go:'
echo
echo "## 附：今日 production .Execute( 全名册（不过滤）＝G1b 的排除项本体"
git grep -nE '\.Execute\(' "$A" -- "$GO" "$GO2" ':!*_test.go'
