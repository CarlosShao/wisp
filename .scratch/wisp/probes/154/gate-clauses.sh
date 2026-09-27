# 票 154 AC#1 —— 触发门的读数生成器（可复算）。G1／G1b／G2／G3／G4 出自票 154；
# G5（OpenScope↔CloseScope 成对普查）由票 158 r1 新增，其形状照 G1–G4，另带正控与负一负两腿。
# 票 161 AC#6 追加（161-r6，09-27，**只加腿、不删腿、不改任何一条既有读数行**）：
#   ① 每一腿开跑前先 want 登记「我今天该不该响」，文末新增【聚合退码】块＝声明与实测不符的腿数；
#     响的腿本身不计进退码（G2、G5-负一负 今天按设计就该响）。
#   ② 两族成对普查各成一腿：G6（OpenTask↔CloseTask）与 G7（Defer↔DisposalScope），
#     形状照 G5＝主尺＋正控＋负一负三发都在；这两族另开反向（方向倒过来也算，见 want 的第三味 rev）。
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

# ---- 票 161 AC#6①：每一腿自己声明，聚合只累计【声明与实测不符】的腿 ----------------
# 形状（判据本体＝票面 AC#6 下面 09-27 09:1x 那个 > 块）：
#   want <腿号> <ring|quiet> [rev]  —— 开跑前登记这一腿今天该不该响；rev＝成对普查两向都算。
#   ring＝该响（run()：名册非空／pair()：未成对≥1 枚）／quiet＝不该响（零命中／0 枚）。
# ⚠ 这一枚**不是**「任何一腿响⇒脚本非 0」：本文件里 G2、G5-负一负、G6、G6-负一负、G7、G7-负一负
#   今天按设计就该响（G5 主尺那 1 枚是票 158 登记在册的活形状），把它们计进退码＝把票 161 AC#7
#   刚收掉的恒红形状从另一扇门放回来（同一枚错，台账 A320／A321）。
LEG='(还没 want)'
EXPECT='quiet'
REV=0
AGG=''
AGG_BAD=0
LEG_COUNT=0

want() { # $1 腿号 $2 ring|quiet [$3 rev]
	LEG="$1"
	EXPECT="$2"
	REV=0
	if [ "${3:-}" = "rev" ]; then REV=1; fi
	if [ "$EXPECT" != ring ] && [ "$EXPECT" != quiet ]; then
		echo "想死在起跑线上：腿 $LEG 的声明 '$EXPECT' 不是 ring／quiet 之一——声明打错＝聚合是枚哑弹" >&2
		exit 3
	fi
}

book() { # $1 这一腿的实测（ring|quiet|git-grep-rc-N）；与 want 的声明比对，只把【不符】计进退码
	local verdict
	if [ "$EXPECT" = "$1" ]; then
		verdict='ok  '
	else
		verdict='BAD '
		AGG_BAD=$((AGG_BAD + 1))
	fi
	LEG_COUNT=$((LEG_COUNT + 1))
	AGG="$AGG# $verdict 腿=$LEG 声明=$EXPECT 实测=$1
"
}

diffsets() { # $1 开方 pattern $2 合方 pattern $3 名册的调用行（已剔注释）
	# 求差集那一段独立成枚函数，理由有两条：pair() 的每一腿吃它；票 161 AC#6②（乙）那发
	# 「文本喂入」样本也吃同一份代码——证算术的尺与产生读数的尺于是不会是两把。
	# stdout＝「$1 出现过、同一文件里 $2 一次都没出现过」的文件清单（换行分隔，无尾随换行）。
	local o c f out=''
	o="$(printf '%s\n' "$3" | grep -Ew "$1" | cut -d: -f2 | sort -u || true)"
	c="$(printf '%s\n' "$3" | grep -Ew "$2" | cut -d: -f2 | sort -u || true)"
	for f in $o; do
		printf '%s\n' "$c" | grep -qxF "$f" || out="$out$f
"
	done
	printf '%s' "$out"
}

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
	# 票 161 AC#6①：run() 这一形的响＝名册非空（git grep rc=0）；rc>1 是尺自己坏了，也算不符。
	local meas='quiet'
	[ "$rc" -eq 0 ] && meas='ring'
	if [ "$rc" -gt 1 ]; then meas="git-grep-rc-$rc"; fi
	book "$meas"
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
	local roster rrc calls unpaired="" runpaired="" n rn f oc deleg rncalls meas
	roster="$(git grep -nEw "$open|$close" "$A" -- $@ 2>/dev/null)"; rrc=$?
	echo "# 逐行读数（名册本体，供下一程 diff；这一份含注释行）："
	if [ -n "$roster" ]; then printf '%s\n' "$roster"; else echo "#   （空）"; fi
	# 成对判据只吃调用点：纯注释行剔掉。
	# 不剔的下场是本程实测到的——它点到了票 158 r1 自己测试文件里那句解释 Mark 机理的注释，
	# 一把会因为散文而响的尺不叫读数（详证据件 §2.4）。名册本体照旧带注释，diff 不受影响。
	calls="$(printf '%s\n' "$roster" | grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' || true)"
	# 票 161 AC#6②：求差集那段交给 diffsets()（下面（乙）文本喂入那发样本吃的是同一份算术）。
	unpaired="$(diffsets "$open" "$close" "$calls")"
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
	if [ "$REV" = 1 ]; then
		# 票 161 AC#6②「方向倒过来也算」：合方出现过、同一文件里开方一次都没有，也算一枚未成对。
		# 这一味不是装饰，G6 主尺今天响的那 1 枚只有反向点得到：
		#   cmd/wisp/run.go:563 调 rt.bridge.CloseTask(taskID)，而 OpenTask 唯一的生产调用点在
		#   internal/tools/bridge.go:559——那枚文件是这一族的定义本体、按票 158 的地界必须排掉，
		#   于是正向在那条射程里【结构上就是 0 枚】。这正是要把反向单独立一味的理由。
		runpaired="$(diffsets "$close" "$open" "$calls")"
		rn="$(printf '%s' "$runpaired" | grep -c . || true)"
		echo "# 反向名册（$close 出现过、同一文件却没见过 $open 的文件；这一味＝「关过却没开过」）："
		if [ -n "$runpaired" ]; then
			while IFS= read -r f; do
				[ -z "$f" ] && continue
				rncalls="$(printf '%s\n' "$calls" | grep -Ew "$close" | cut -d: -f2 | grep -cxF "$f" || true)"
				echo "#   UNPAIRED-REV $f (合方调用点=${rncalls:-0})  # 同文件没见过开方：人来判「真漏／委托」"
			done <<< "$runpaired"
		else
			echo "#   （空）"
		fi
		echo "# 反向未成对枚数＝$rn   （0＝反向没响）"
		n=$((n + rn))
	fi
	echo "# 未成对枚数＝$n   （0＝这把尺没响／>0＝逐枚点名如上）"
	echo "# git grep rc=$rrc（1＝本射程里这一族一枚都没有——那是「射程里没这东西」，不是「都关好了」；两种 0 枚响含义不同）"
	echo "# rc=$n   （本枚子句的响＝未成对枚数；0＝不响）"
	# 票 161 AC#6①：pair() 这一形的响＝未成对枚数≥1（两向之和，若这一腿开了 rev）。
	meas='quiet'
	[ "$n" -gt 0 ] && meas='ring'
	if [ "$rrc" -gt 1 ]; then meas="git-grep-rc-$rrc"; fi
	book "$meas"
}

# ---- 票 161 AC#6②（乙）文本喂入模式：只证「求差集那段算术」不撒谎 -------------------
# 名册从 stdin 进，两个方向的未成对清单出；这一枚不吃 git、不吃 $A、不吃 pathspec，
# 所以它证的不是（甲）合成树那发读数（那才是主证），它证的只是 diffsets() 本身。
# 用法：DIFFSETS_OPEN=OpenTask DIFFSETS_CLOSE=CloseTask sh gate-clauses.sh --diffsets < 名册.txt
# 名册行形状与 pair() 的逐行读数一致（path:line:text）；文件头那三行 # 是锚点标记，正常现象。
if [ "${1:-}" = "--diffsets" ]; then
	ds_roster="$(cat)"
	ds_calls="$(printf '%s\n' "$ds_roster" | grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' || true)"
	ds_open="${DIFFSETS_OPEN:-OpenTask}"
	ds_close="${DIFFSETS_CLOSE:-CloseTask}"
	ds_fwd="$(diffsets "$ds_open" "$ds_close" "$ds_calls")"
	ds_rev="$(diffsets "$ds_close" "$ds_open" "$ds_calls")"
	echo "# diffsets 模式：开方=$ds_open 合方=$ds_close（名册 $(printf '%s' "$ds_roster" | grep -c . || true) 行 → 剔注释后调用行 $(printf '%s' "$ds_calls" | grep -c . || true) 行）"
	echo "# 正向（$ds_open 却没 $ds_close）："
	if [ -n "$ds_fwd" ]; then printf '%s\n' "$ds_fwd" | sed 's|^|#   |'; else echo "#   （空）"; fi
	echo "# 反向（$ds_close 却没 $ds_open）："
	if [ -n "$ds_rev" ]; then printf '%s\n' "$ds_rev" | sed 's|^|#   |'; else echo "#   （空）"; fi
	exit 0
fi

# ---- 以下每一腿的 want＝票 161 AC#6① 的声明，逐枚都是 09-27 在锚上现量登记的今天真值 ----
want G1 quiet
run "G1 生产码里构造 ToolRequest（排 internal/agent/ 自身）" \
	'ToolRequest\{' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/*'

want G1b quiet
run "G1b 生产码里任意 .Execute( 调用点（排 loop.go 与 bridge.go 自身）" \
	'\.Execute\(' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'

# G2 今天【该响】：名册非空是它的工作——cmd/wisp/run.go 那一发 CloseTask 就在册。
want G2 ring
run "G2 OpenTask/CloseTask 的生产出现点（排 bridge.go 自身；今日应只剩 cmd/wisp/run.go 那一发＋它上面的注释）" \
	'OpenTask|CloseTask' '-w' "$GO" "$GO2" ':!*_test.go' ':!internal/tools/bridge.go'

want G3 quiet
run "G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法" \
	'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' '' internal/agent

want G4 quiet
run "G4 生产组合根打开 L0 直通（PassThroughUnclassifiedRisk 的字面赋值）" \
	'PassThroughUnclassifiedRisk:' '' '*.go' ':!*_test.go'

# G5（票 158 AC#2 新增）射程＝OpenScope↔CloseScope 的成对普查。
# G1–G4 盯的是「宿主派发工具却没关环路那一枚 id」那一族（过桥那条腿），
# 而「不过桥、直接对 risk 层开一枚 scope」这一形今天有活样板：
#   cmd/wisp/panel_assets.go:232  prov.OpenScope(taintSourceScopeID)   且 cmd/wisp 下 CloseScope 零枚。
# 排掉 internal/risk/* 的理由：OpenScope/CloseScope 的**定义**在那儿（票 158 地界明令只盘调用者），
# 把定义算成调用点会让名册每改一行注释就响一次，那是噪音不是读数。
# G5 主尺今天【该响】，且响 1 枚：cmd/wisp/panel_assets.go:232 开了 C25 scope 却没在同一文件关。
# 这一枚是票 158 登记在册的活形状、不是本票的靶子（票面 AC#6 那条红线），本程一字节生产码都没动。
want G5 ring
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
want G5pos quiet
pair "G5-正控 同尺＝主尺射程，只多排掉主尺今天点名的那一枚 panel_assets.go ⇒ 剩下的开方全是 balanced 样本" \
	'OpenScope' 'CloseScope' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*' ':!cmd/wisp/panel_assets.go'

echo
echo "## G5 负一负（同尺摘掉「排测试＋排定义」那两条 ⇒ 今日应非空，证明这把尺吃得住这个形状、不是结构上产不出读数的装饰）"
# G5-负一负【按设计就该响】（今天 5 枚）：这一腿的作用就是证明这把尺吃得住这个形状。
want G5neg ring
pair "G5-负一负 同尺不过滤任何文件" \
	'OpenScope' 'CloseScope' "$GO" "$GO2"

echo
echo "## 附：G1 的负一负（同尺摘掉排测试那一条 ⇒ 今日即非空，证明模式吃得住这个形状，不是结构上产不出读数的装饰）"
echo "\$ git grep -cE 'ToolRequest\\{' $A -- internal cmd 里 *_test.go 的枚数"
git grep -E 'ToolRequest\{' "$A" -- "$GO" "$GO2" | grep -c '_test\.go:'
echo
echo "## 附：今日 production .Execute( 全名册（不过滤）＝G1b 的排除项本体"
git grep -nE '\.Execute\(' "$A" -- "$GO" "$GO2" ':!*_test.go'

# =====================================================================================
# 票 161 AC#6② —— 两族成对普查各成一腿（161-r6 新增，09-27）。
# 位置一律在文末、既有录用法之后：上面那些名册是要被票 154／158 的证据件按行 diff 的，
# 新腿插在中间会让行号整体后移，插在文末则 diff 只剩纯新增。
# 形状照 G5＝主尺＋正控＋负一负三发都在；want 的第三味 rev＝「方向倒过来也算」。
# 下面注释里那几处行号（bridge.go:559／:633／:684、disposal.go:129／:186／:191／:253、
# run.go:563）是本程 git grep 现量的，命令与全份读数见 docs/evidence/s1/161-gate-aggregate-r6.md。
# =====================================================================================

echo
echo "## G6 主尺（票 161 AC#6②第一族）OpenTask↔CloseTask：派发了环路却没在同一文件关掉那一枚 id"
echo "# 与 G5 的分工：G5 盘「不过桥、直接对 risk 层开 scope」那一形，G6 盘「过桥的 task 生命周期」那一形。"
echo "# 排 internal/tools/bridge.go 的理由与 G5 排 internal/risk 同一条：OpenTask／CloseTask 的定义本体"
echo "#   在那儿（:633 与 :684），把定义算成调用点＝名册每改一行注释就响一次，那是噪音不是读数。"
echo "# 今日真实枚数（09-27 现量，两向分开登记）：正向 0 枚／反向 1 枚＝cmd/wisp/run.go。"
echo "#   正向那 0 枚**不是**「都关好了」：开方唯一的生产调用点（bridge.go:559 b.OpenTask(dec.TaskID)）"
echo "#   正好落在被排掉的定义本体文件里，这条射程的结构就没有开方可数——所以这一腿的响只由反向撑着，"
echo "#   「新造一发开了却没关的它必须响」那件事交给（甲）合成树那发主证，不靠这条射程的空 0。"
want G6 ring rev
pair "G6 主尺 OpenTask↔CloseTask 成对普查（排 *_test.go、排定义本体 internal/tools/bridge.go；两向都算）" \
	'OpenTask' 'CloseTask' "$GO" "$GO2" ':!*_test.go' ':!internal/tools/bridge.go'

echo
echo "## G6 正控（同一把尺打在已知「开了又关」的样本上，两向都必须不响）"
echo "# 样本＝internal/tools/bridge.go 自己：:559 的 b.OpenTask 与 :684 的 CloseTask 在同一枚文件里。"
echo "# 这条射程是真吃得住的（本程现量 git grep rc=0、剔注释后 3 行调用点），不是「射程里根本没这东西」"
echo "#   那种空 0——那种 0 不算正控（票 154 AC#3 就是拿「零区别」否掉装饰腿的，G5 的自纠 (1) 同一条）。"
want G6pos quiet rev
pair "G6-正控 同尺只放已知「开又合」的那一枚 internal/tools/bridge.go ⇒ 两向都必须 0 枚" \
	'OpenTask' 'CloseTask' 'internal/tools/bridge.go' ':!*_test.go'

echo
echo "## G6 负一负（同尺摘掉「排测试＋排定义」两条 ⇒ 今日应非空，证明这把尺吃得住这个形状）"
echo "# 今日真实枚数：正向 0 枚／反向 1 枚＝cmd/wisp/run.go。测试文件里 OpenTask 与 CloseTask 同文件成对"
echo "#   （internal/tools/bridge_scope_open_ticket158_test.go），所以摘掉过滤也点不出正向未成对——"
echo "#   本程不把它谎称成点得出；正向那一向的证＝（甲）合成树。"
want G6neg ring rev
pair "G6-负一负 同尺不过滤任何文件" \
	'OpenTask' 'CloseTask' "$GO" "$GO2"

echo
echo "## G7 主尺（票 161 AC#6②第二族）Defer↔DisposalScope：挂了收尾却没绑到声明的 scope／反向＝拿着 scope 却没挂收尾"
echo "# 定义本体＝internal/plugin/disposal.go（type DisposalScope :129、Defer :186、DeferNamed :191、"
echo "#   Dispose :253），与 G5 排 internal/risk、G6 排 bridge.go 同一条理由，排掉它。"
echo "# 开方写 'Defer(Named)?' 而不是 'Defer'：本程现量 printf 'x.DeferNamed(y)' | grep -Ew 'Defer' ⇒ rc=1，"
echo "#   -w 把 DeferNamed 当另一个词，光写 Defer 会漏掉一半的收尾登记点（换成 Defer(Named)? 后 rc=0）。"
echo "# 合方同理写 '(New)?DisposalScope'：grep -Ew 'DisposalScope' 对 'plugin.NewDisposalScope(...)' 也是 rc=1"
echo "#   ——同一个词被前缀粘住，只排 DisposalScope 会把「用 NewDisposalScope 造了 scope、又 Defer 到它上面」的"
echo "#   文件误判成未成对。今天这一改不动任何读数：本程现量 git grep -nEw NewDisposalScope <锚> -- internal cmd"
echo "#   排测试排定义本体 ⇒ rc=1，生产码里一枚构造点都还没有，所以合方名册与改前逐字相同。"
echo "# 今日真实枚数：正向 0 枚（Defer／DeferNamed 的生产调用点全在被排掉的定义本体里，测试之外零枚）／"
echo "#   反向 3 枚＝internal/memory/retention.go、internal/observe/goroutine.go、internal/proc/shutdown.go。"
echo "#   反向这 3 枚里 retention.go:98 是真拿 scope 的签名，另两枚各有一味是「代码行尾随注释」被算进"
echo "#   调用行的已知粗糙处（pair() 只剔以 // 开头的整行）——逐枚人来判真漏／委托，本程不替它判。"
want G7 ring rev
pair "G7 主尺 Defer↔DisposalScope 成对普查（排 *_test.go、排定义本体 internal/plugin/disposal.go；两向都算）" \
	'Defer(Named)?' '(New)?DisposalScope' "$GO" "$GO2" ':!*_test.go' ':!internal/plugin/disposal.go'

echo
echo "## G7 正控（同一把尺打在已知「收尾与 scope 同文件」的样本上，两向都必须不响）"
echo "# 样本＝internal/plugin/disposal.go 自己：Defer／DeferNamed 与 DisposalScope 同文件 ⇒ 两向 0 枚，"
echo "#   且名册非空（本程现量 git grep rc=0），不是「射程里没这东西」那种空 0。"
want G7pos quiet rev
pair "G7-正控 同尺只放定义本体 internal/plugin/disposal.go 这一枚已知「两向都有」的样本" \
	'Defer(Named)?' '(New)?DisposalScope' 'internal/plugin/disposal.go' ':!*_test.go'

echo
echo "## G7 负一负（同尺不过滤任何文件 ⇒ 今日应非空，证明这把尺吃得住这个形状）"
echo "# 今日真实枚数：正向 1 枚＝internal/risk/provenance_test.go（出现 Defer 却没在同一文件出现"
echo "#   DisposalScope 这个词）／反向 3 枚＝主尺那三枚 ⇒ 合计 4 枚。"
want G7neg ring rev
pair "G7-负一负 同尺不过滤任何文件" \
	'Defer(Named)?' '(New)?DisposalScope' "$GO" "$GO2"

# =====================================================================================
# 票 161 AC#6① —— 聚合退码。放在文件最后一格，是因为前面每一腿都得保持「逐行可 diff」原样。
# 它要修的正是票 158／票 161 那句「门点响了，可它还是退 0」：改之前本文件的退码来自最后那发
# 附录用法的 git grep（今天 rc=0），哪几腿在响都跟退码无关＝聚合是一枚恒绿。
# =====================================================================================
echo
echo "## 聚合退码（票 161 AC#6①）：只累计【声明与实测不符】的腿"
echo "# 判据**不是**「任何一腿响⇒脚本非 0」——见文件头 want 那一段与本票票面 09-27 09:1x 的 > 块："
echo "#   G2／G5-负一负／G6／G6-负一负／G7／G7-负一负 今天按设计就该响，把它们计进退码"
echo "#   ＝把票 161 AC#7 刚收掉的恒红形状从另一扇门放回来（同一枚错，台账 A320／A321）。"
echo "# 每腿判定（ok＝说到做到／BAD＝言行不一；这一份同样逐行可 diff）："
if [ -n "$AGG" ]; then printf '%s' "$AGG"; fi
echo "# 腿数＝$LEG_COUNT 声明与实测不符＝$AGG_BAD"
echo "# 聚合退码＝$AGG_BAD   （0＝每一腿的响与不响都和自己登记的声明一致；>0＝有腿言行不一）"
echo "# 翻一枚声明让退码跟着变那发证据不在这里，在 .scratch/wisp/probes/161/r6/flip-declaration.sh"
if [ "$AGG_BAD" -gt 125 ]; then
	echo "# 不符腿数大于 125，退码截在 125（够红了，别拿溢出当读数）"
	exit 125
fi
exit "$AGG_BAD"
