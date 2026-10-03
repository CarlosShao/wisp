# 票 154 AC#1 —— 触发门的读数生成器（可复算）。G1／G1b／G2／G3／G4 出自票 154；
# G5（OpenScope↔CloseScope 成对普查）由票 158 r1 新增，其形状照 G1–G4，另带正控与负一负两腿。
# 票 161 AC#6 追加（161-r6，09-27，**只加腿、不删腿、不改任何一条既有读数行**）：
#   ① 每一腿开跑前先 want 登记「我今天该不该响」，文末新增【聚合退码】块＝声明与实测不符的腿数；
#     响的腿本身不计进退码（G2、G5-负一负 今天按设计就该响）。
#   ② 两族成对普查各成一腿：G6（OpenTask↔CloseTask）与 G7（Defer↔DisposalScope），
#     形状照 G5＝主尺＋正控＋负一负三发都在；这两族另开反向（方向倒过来也算，见 want 的第三味 rev）。
# 票 171 r1 追加（09-27，**只加声明、加味、加形状；一枚腿没删、一条 pathspec 射程没减**）：
#   ① AC#2＝枚数进判据：新增 want_n <基线枚数>，pair 腿的响由「空／非空」变成「实测枚数 vs 基线枚数」；
#     实测 > 基线 才 BAD（新增的未成对躲不进「非空」里了），实测 < 基线 打 SHR（变好了不恒红，但声明过期可见）。
#     pair() 没拿到 want_n 就直接死（rc=4）——不让下一程新加一枚腿又把 AC#2 那枚洞开回去。
#   ② AC#6＝合方词根认句柄那一形：G5 三腿的合方从写死的 'CloseScope' 换成 'CloseScope|Close\(\)'
#     （票 160 之后 risk.OpenScope 交回 *Scope 句柄，收尾是句柄自带的 .Close()，见 bridge.go:706），
#     并新增 pair() 的第四味 handle＝「开方把句柄丢了（`_ =` 或裸调用）⇒ 这一族在这枚文件里结构上关不掉，
#     合方写什么都不算关」。G6／G7 明写 '-'：它们的开方（OpenTask 无返回值、Defer 返回 bool）不交句柄。
#   两格撞同一枚 pair helper ⇒ 同一程做完（票面 171 AC#6 末行「别拆两程各改一次同一个函数」）。
#   ③ 票 171 r1 的形制一字未动：pair 腿仍然「实测 > 基线才 BAD」、拿不到 want_n 直接死（rc=4）。
# 票 171 r2 追加（09-27，**只加声明、加味、加断言；一枚腿没删、一条 pathspec 射程没减**）：
#   ① AC#7＝枚数进 run() 那一族腿：G1／G1b／G2／G3／G4 从此也登记 want_n（run 腿的枚数＝命中行数），
#     记账交给 pair 腿那同一枚 book_count——不再有"只有 0／非 0"的第二套判法。
#     ⚠ 不选"逐枚相等"那一形，理由与 AC#2 结案记录同一条（相等的形会让改好了也变红＝台账 A320／A321）。
#   ② AC#8①＝腿数断言：腿的名册只写在 LEGS_EXPECT 一处（腿数从它现算，不在第二个地方抄一遍号）；
#     want() 先核腿号在不在册／有没有重复，聚合段再核【名册里每一腿都真的记过账】——
#     整条腿被摘掉（want 与它的 run()/pair() 一起删）从此是一枚 BAD，不是一句没人看的打印。
#   ③ AC#8②＝SHR 那一味仍然**不计进退码**（变好了不算言行不一），但单列一张【基线过期】表，
#     不再只是聚合行里那句行内注；核它的那一把＝probes/161/r6/flip-declaration.sh 新增的那一格。
# 票 171 r3 追加（10-03，**只改分子口径；一枚腿没删、一条 pathspec 射程没减、名册本体逐行原样照打**）：
#   ① AC#1（尺洞②＝尾随注释被当成调用点）＝两族腿的分子一次改齐，别只修一扇门（票面 AC#1 那格
#     09-27 19:1x 追加的射程提醒逐字要求这件事；171-v2 表 §14 第一条登记的正是同一枚洞的另一半）：
#       pair 腿——改前只剔"行首是 //"的整行注释，行尾注释里的词照算进调用行（现量两枚见 call_lines
#         上方那一段：G7 反向在锚 4f2a777b 点名的 3 枚里 2 枚的全部证据就是行尾注释）；
#       run 腿——改前分子＝整张名册的行数，连整行注释都算命中（现量：G2 在同锚的 3 行命中里 2 行是注释）。
#       两族从此共用同一处实现 call_lines()；--diffsets 那味文本喂入也改吃它（它改前自己抄了一遍
#       "只剔整行注释"，那就是"同源拷贝逐枚有归属"那起事故的形状）。
#   ② 分子口径变了 ⇒ 十四腿的 want_n 基线在同一次改动里逐枚复算（票面 AC#8 那条"谁改口径就在同一格
#     把 want_n 核下来，别留 SHR 当默认"的规矩）；ring/quiet 若因新分子换位，也一并现量重登。
#     每一腿那一格从此自己说清两个数：# 名册行数＝（含注释，逐行原样，供下一程 diff）与
#     # 命中行数＝／# 未成对枚数＝（分子），被剔掉的行逐枚打在 #   COMMENT-ONLY 那一格里。
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
# 票 171 AC#2 加味：
#   want_n <枚数>                   —— 紧跟 want 登记【基线枚数】。票 171 AC#2 时它只管 pair 腿
#                                     （枚数＝未成对），票 171 AC#7 起 run 腿也必填（枚数＝命中行数）；
#                                     没登记就记账＝死在起跑线上（rc=4），下一程新加一条腿不可能悄悄退回
#                                     「空／非空」那一形。
# 票 171 AC#8① 加味：
#   LEGS_EXPECT                     —— 腿的名册（唯一一处写这些腿号的地方）；want() 核腿号在不在册，
#                                     聚合段核名册里每一腿都记过账，缺一枚＝一枚 BAD。
# ⚠ 这一枚**不是**「任何一腿响⇒脚本非 0」：本文件里 G2、G5-负一负、G6、G6-负一负、G7、G7-负一负
#   今天按设计就该响（G5 主尺那 1 枚是票 158 登记在册的活形状），把它们计进退码＝把票 161 AC#7
#   刚收掉的恒红形状从另一扇门放回来（同一枚错，台账 A320／A321）。
LEG='(还没 want)'
EXPECT='quiet'
REV=0
COUNT=''
AGG=''
AGG_BAD=0
LEG_COUNT=0
STALE=''
STALE_COUNT=0
LEG_SEEN=''
BOOKED=''
WANT_COUNT=0

# ---- 票 171 AC#8①：腿的名册。枚数（14）不在别处再写一遍，聚合段从这张表现算 ----
# 这张表就是"登记期望值写进尺旁边、与九腿 want 同源"那一格的答案：腿号只在这里出现一次，
# want() 拿它核对，聚合段拿它求缺——所以摘掉一整条腿＝这张表缺一枚＝退码跟着动。
LEGS_EXPECT='G1 G1b G2 G3 G4 G5 G5pos G5neg G6 G6pos G6neg G7 G7pos G7neg'
LEG_EXPECT_COUNT=0
for __le in $LEGS_EXPECT; do LEG_EXPECT_COUNT=$((LEG_EXPECT_COUNT + 1)); done

want() { # $1 腿号 $2 ring|quiet [$3 rev]
	# 票 171 AC#8①：腿号先对名册。重复／打错／不在册一律死在起跑线上（rc=3）——
	# 一张对不上账的名册比没有名册更坏：它会把"缺腿"那枚读数折成假 0。
	case " $LEG_SEEN " in
	*" $1 "*)
		echo "想死在起跑线上：腿 $1 第二次 want——同一条腿记两次账＝腿数与实测从此对不上" >&2
		exit 3
		;;
	esac
	if ! printf '%s\n' $LEGS_EXPECT | grep -qxF "$1"; then
		echo "想死在起跑线上：腿 $1 不在 LEGS_EXPECT 名册里——新加一枚腿要同时在名册登记（票 171 AC#8①）" >&2
		exit 3
	fi
	LEG="$1"
	EXPECT="$2"
	REV=0
	COUNT=''
	WANT_COUNT=$((WANT_COUNT + 1))
	LEG_SEEN="$LEG_SEEN $LEG"
	if [ "${3:-}" = "rev" ]; then REV=1; fi
	if [ "$EXPECT" != ring ] && [ "$EXPECT" != quiet ]; then
		echo "想死在起跑线上：腿 $LEG 的声明 '$EXPECT' 不是 ring／quiet 之一——声明打错＝聚合是枚哑弹" >&2
		exit 3
	fi
}

# want_n（票 171 AC#2，票 171 AC#7 起两族腿都管）：基线枚数。放在 want 之后、腿之前＝「开跑前先登记」这条没破。
# 不校「quiet 却登记了 N>0」那种自相矛盾的声明——book_count 里【空/非空】那一味先它一步响，
# 而且 flip-declaration.sh 正是靠改 want 那一行来验「声明还连着退码」，起跑线上把脚本打死会把
# 它要证的那件事（聚合累计声明与实测不符）换成另一件事（脚本崩在登记处）。
want_n() { # $1 基线枚数
	if [ "$LEG" = '(还没 want)' ]; then
		echo "想死在起跑线上：want_n 排在任何 want 之前——基线没有腿可登记" >&2
		exit 3
	fi
	case "${1:-}" in
	'' | *[!0-9]*)
		echo "想死在起跑线上：腿 $LEG 的基线枚数 '${1:-}' 不是数字——枚数打成散文＝聚合是枚哑弹" >&2
		exit 3
		;;
	esac
	COUNT="$1"
}

# book_count：两族腿共用的记账（票 171 r2 起 run 腿也走这一枚，「只有 0／非 0」那一套从此不存在了；
# 改前那枚只管 run 腿的 book() 因此整枚删掉——留着它就是同一件事的第二个判法，正是"同源拷贝"的形状）。
# $1＝实测枚数（pair 腿＝未成对枚数；run 腿＝命中行数） $2＝那一发 git grep 的 rc
# $3＝kind：unpaired（默认＝pair 腿）| hits（run 腿，票 171 AC#7）。kind 只决定红句里「新增」后面
#   那两个词，判据算术一枚不多一枚不少 ⇒ AC#2 结案记录的逐字红句（腿=G5 …因=新增未成对）仍然对得上。
# 两味**并集**，一枚腿最多计一进退码（不重复计）：
#   ① 票 161 AC#6① 那一味（空／非空 vs want 的 ring／quiet）——一字不改地继续校，
#      否则 .scratch/wisp/probes/161/r6/flip-declaration.sh 那发「翻声明退码跟着变」会被我换成另一件事；
#   ② 票 171 AC#2／AC#7 那一味（实测枚数 vs want_n 基线）——同一条腿里新增的第二枚从此躲不进「非空」。
# 实测 < 基线 打 SHR 不计退码：改好了不算言行不一（要算的话这把门就为好消息响，
# 下一程就去放宽它——那正是票 161 AC#7／台账 A320、A321 收掉的恒红形状）。
# 但 SHR 从此**单列一张【基线过期】表**（票 171 AC#8②）：改前它只是聚合行里的一句行内注，
# 现量那一发（want_n 8 抬成 9、实测 8）＝退码 0、单列表 0 枚，也就是"只有人眼看得见"。
book_count() {
	local n="$1" rrc="$2" kind="${3:-unpaired}" verdict meas bin='quiet' why=''
	if [ -z "$COUNT" ]; then
		echo "想死在起跑线上：腿 $LEG 没登记 want_n 基线——枚数不进判据＝票 171 AC#2／AC#7 那枚洞还开着" >&2
		exit 4
	fi
	meas="${n}枚"
	[ "$n" -gt 0 ] && bin='ring'
	if [ "$rrc" -gt 1 ]; then
		# 尺自己坏了（git grep 出错，不是零命中）：这一腿的 0 枚没有资格说「都关好了」。
		verdict='BAD '
		meas="git-grep-rc-$rrc"
	elif [ "$bin" != "$EXPECT" ]; then
		# ① 票 161 AC#6① 那一味：响不响＝名册空不空，逐字照改前的判法（flip-declaration 连的是这一味）。
		verdict='BAD '
		why=' 因=空/非空那一味与声明不符（票 161 AC#6①）'
	elif [ "$n" -gt "$COUNT" ]; then
		verdict='BAD '
		if [ "$kind" = hits ]; then
			why=' 因=新增命中（票 171 AC#7：实测 > 基线）'
		else
			why=' 因=新增未成对（票 171 AC#2：实测 > 基线）'
		fi
	elif [ "$n" -lt "$COUNT" ]; then
		# 变好＝不算言行不一，但基线于是过期了：行进聚合表（逐字照改前），同时进文末单列的那张表。
		verdict='SHR '
		why=' 注=读数变好了，基线过期，下一程把 want_n 核下来'
		STALE="$STALE# STALE 腿=$LEG 声明=$EXPECT 基线=${COUNT}枚 实测=${COUNT}枚 差=0枚 处置=把 want_n 核下来；样本一枚没变好而是被摘掉的，先回答谁摘的
"
		STALE_COUNT=$((STALE_COUNT + 1))
	else
		verdict='ok  '
	fi
	[ "$verdict" = 'BAD ' ] && AGG_BAD=$((AGG_BAD + 1))
	LEG_COUNT=$((LEG_COUNT + 1))
	BOOKED="$BOOKED $LEG"
	AGG="$AGG# $verdict 腿=$LEG 声明=$EXPECT 基线=${COUNT}枚 实测=$meas$why
"
}


# discarded_files（票 171 AC#6 的 handle 味）：开方那一行没把返回的句柄接下来。
# 两种拼写：(1) 开方之前整行一个赋值号都没有＝裸调用，返回值落地即弃；
#           (2) 赋值给 `_`＝显式丢（今天生产码里那一枚在册活形状 cmd/wisp/panel_assets.go 就是这一形）。
# 句柄被丢＝这一族的收尾在这枚文件里【结构上不可能发生】，所以合方写什么都不该认它关过。
# 先剔定义行：`func (p *Provenance) OpenScope(...)` 那一行天生没有赋值号，把它当调用点
# 就是拿声明说成读数（G5-负一负那枚射程连 internal/risk 的定义本体一起吃，171-r1 合成树的
# 样本也这么写——两处都会踩；本程实测踩过一次，读数见证据件 §3 的改前读数那一格）。
# 只这一处实现：pair() 用它求名册、也用它打批注（同源拷贝逐枚有归属那起事故的教训）。
discarded_files() { # $1 开方 pattern $2 调用行 -> 丢句柄的文件清单
	printf '%s\n' "$2" | grep -Ew "$1" |
		grep -vE ':[0-9]+:[[:space:]]*func[[:space:]]' |
		grep -E "(^[^=]*)$1|(^|[[:space:];{(])_[[:space:]]*=[^;]*$1" |
		cut -d: -f2 | sort -u || true
}

# ---- 票 171 r3 AC#1：分子只吃代码那一半（两族腿共用这一处实现，别写第二份） ----------------
# 尺洞②的形状（161-v2 裁过、票 171 立案）：成对普查那把尺剔"整行注释"，但**行尾注释里的词它照算**，
# 于是"有人在代码里提了一下开方/关范围"被当成"真调用了"。现量两枚（锚 4f2a777b，读数＝
#   probes/171/r3/logs/gate-pre-at-anchor.txt）：
#   internal/observe/goroutine.go:34  CategoryTemporary GoroutineCategory = "temporary" // inside a DisposalScope
#   internal/proc/shutdown.go:85      ReleaseSpeechSessions func(ctx context.Context) error // step 5 (speech, via DisposalScope)
# ⇒ G7 反向今天点名的 3 枚里有 2 枚的"合方在场"证据就是这种行尾注释（161-gates-accept-r2.md:67/:68 同判）。
# run() 那一族更宽：改前它的分子＝整张名册的行数，连整行注释都不剔（同一枚洞的另半边，
#   票面 AC#1 那格 09-27 19:1x 的追加提醒点名的就是它；171-v2 表 §14 第一条登记同案）。
# 只切"注释那一半"，不切射程：一条 pathspec 没动、一枚腿没删、名册本体照旧逐行原样打出来供 diff。
# 为什么按"空白字符＋//"切而不按"第一个//"切：gofmt 之后行尾注释前必有一个空格，而字符串里的
# "https://x" 前面是冒号不是空格——所以 URL 不会被当成注释切掉（这一味是本程现量选的，不是设想）。
strip_comment() { # stdin 名册行 -> 只留代码那一半（行尾的 // 与 /* 注释切掉）
	sed -E -e 's|[[:space:]]//.*$||' -e 's|[[:space:]]/\*.*$||'
}

call_lines() { # $1 名册（git grep 的逐行读数） $2 词表 $3 匹配味（-Ew 默认＝pair 腿与 diffsets 同味；-E＝run 腿与它自己的 git grep -nE 同味）
	# 三道工序，缺一不可：① 剔整行注释（改前 pair() 已有这一味，逐字沿用）；
	# ② 剔行尾注释（票 171 AC#1 新加的那一味）；③ 剔"切完就不再匹配词表"的行——
	# 少了③，一枚只靠注释进场的行仍然留在调用行里，diffsets 照它求差集，等于没修。
	# ③ 用的那一味必须与这一腿自己捞名册时那一味一致（-E 或 -Ew），否则"命中"与"仍匹配"是两套词法。
	local m="${3:--Ew}"
	printf '%s\n' "$1" |
		grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' |
		strip_comment |
		grep "$m" "$2" || true
}

comment_only_lines() { # $1 名册 $2 词表 [$3 匹配味] -> 被分子剔掉的行，逐枚列出来给人复核（读数不藏在算式里）
	local l code m="${3:--Ew}"
	printf '%s\n' "$1" | while IFS= read -r l; do
		[ -z "$l" ] && continue
		code="$(printf '%s\n' "$l" | grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' | strip_comment || true)"
		if ! printf '%s\n' "$code" | grep -q "$m" "$2"; then printf '%s\n' "$l"; fi
	done || true
}

diffsets() { # $1 开方 pattern $2 合方 pattern $3 名册的调用行（已剔注释）[$4 强制点名的文件清单]
	# 求差集那一段独立成枚函数，理由有两条：pair() 的每一腿吃它；票 161 AC#6②（乙）那发
	# 「文本喂入」样本也吃同一份代码——证算术的尺与产生读数的尺于是不会是两把。
	# stdout＝「$1 出现过、同一文件里 $2 一次都没出现过」的文件清单（换行分隔，无尾随换行）。
	# $4（票 171 AC#6）＝就算 $2 出现过也**照样点名**的文件（丢句柄那一族用）；不传＝行为逐字同改前。
	local o c f out='' force="${4:-}"
	o="$(printf '%s\n' "$3" | grep -Ew "$1" | cut -d: -f2 | sort -u || true)"
	c="$(printf '%s\n' "$3" | grep -Ew "$2" | cut -d: -f2 | sort -u || true)"
	for f in $o; do
		if [ -n "$force" ] && printf '%s\n' "$force" | grep -qxF "$f"; then
			out="$out$f
"
			continue
		fi
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
	# 票 171 AC#7：名册先收进变量再原样打出来（行内容一字不差，只是多次数一遍行），
	# 因为这一族腿的响从此看枚数：命中行数 vs want_n 登记的基线。
	# 票 171 AC#1（本程）：分子从此只数"代码里真有那一味"的行——
	#   改前 hits＝整张名册的行数，连整行注释都算命中（现量：G2 在锚 4f2a777b 的 3 行命中里
	#   cmd/wisp/run.go:942 与 internal/tools/subagent_197.go:332 两行都是 // 开头的整行注释）。
	#   名册本体照旧逐行原样打印（下一程仍按行 diff），只是"行数"这一格换成了剔注释之后的数，
	#   并附一张被剔掉的行——两族腿同一处实现（call_lines），不写第二份。
	local roster rc hits raw_hits dropped
	roster="$(git grep -nE $flags "$pat" "$A" -- $@ 2>/dev/null)"; rc=$?
	if [ -n "$roster" ]; then printf '%s\n' "$roster"; fi
	raw_hits="$(printf '%s' "$roster" | grep -c . || true)"
	hits="$(call_lines "$roster" "$pat" -E | grep -c . || true)"
	dropped="$(comment_only_lines "$roster" "$pat" -E)"
	echo "# 名册行数＝$raw_hits （含注释行，逐行原样＝上面那一份，供下一程 diff；这一格一字未改）"
	echo "# 命中行数＝$hits   （票 171 AC#7 的枚数本体；票 171 AC#1 起分子只吃代码那一半——同一条腿里的第二枚从此躲不进「非空」）"
	if [ -n "$dropped" ]; then
		echo "# 被分子剔掉的行（注释里提到词表，代码里没有）："
		printf '%s\n' "$dropped" | sed 's|^|#   COMMENT-ONLY |'
	else
		echo "# 被分子剔掉的行：（空＝今天这一腿没有一枚靠注释进场）"
	fi
	echo "# 基线枚数＝${COUNT:-（没登记）}   （want_n 登记的那一发；AC#7 起 run 腿也必填，没登记＝rc=4）"
	echo "# rc=$rc   （1＝零命中／0＝非空＝该枚子句已响）"
	# 票 161 AC#6① 那一味（空／非空）＋票 171 AC#7 那一味（枚数 vs 基线）都在 book_count 里；
	# rc>1 仍是「尺自己坏了」，那一味在 book_count 里先它一步。改前的 book() 已并入 book_count。
	book_count "$hits" "$rc" hits
}

pair() { # 成对普查：$1 标题 $2 开方 pattern $3 合方 pattern $4 handle|- $5.. pathspecs
	# 与 run() 的区别：run() 的响不响看名册非空，pair() 的响不响看「开了却没在同一文件关过」。
	# 一条 git grep 表达不出成对，所以这一枚是【一把尺 ＋ 两次求差集】；
	# 锚点、pathspec 与逐行可 diff 的名册本体都照 run() 的形状。
	# $4（票 171 AC#6）＝handle 时开方丢句柄的文件强制点名；'-' ＝不启用（G6/G7 的开方不交句柄）。
	# 响的粒度（票 171 AC#2）＝want_n 登记的基线枚数，不再是非空——见 book_count()。
	local title="$1" open="$2" close="$3" mode="$4"; shift 4
	echo
	echo "## $title"
	echo "\$ git grep -nEw '$open|$close' $A -- $*"
	# shellcheck disable=SC2068
	local roster rrc calls dropped='' unpaired="" runpaired="" n rn f oc deleg rncalls meas disc='' generic='' gk
	roster="$(git grep -nEw "$open|$close" "$A" -- $@ 2>/dev/null)"; rrc=$?
	echo "# 逐行读数（名册本体，供下一程 diff；这一份含注释行）："
	if [ -n "$roster" ]; then printf '%s\n' "$roster"; else echo "#   （空）"; fi
	# 成对判据只吃调用点：纯注释行剔掉。
	# 不剔的下场是本程实测到的——它点到了票 158 r1 自己测试文件里那句解释 Mark 机理的注释，
	# 一把会因为散文而响的尺不叫读数（详证据件 §2.4）。名册本体照旧带注释，diff 不受影响。
	# 票 171 AC#1 把这一味补全：改前只剔"行首是 //"的整行，行尾注释里的词照算进调用行，
	# 于是"提了一下开方/关范围"＝"真调用了"（两枚现量见 strip_comment 上方那一段）。
	# 剔掉的行不藏在算式里，逐枚打在下面这一格里，谁想复核都对得上名册原行。
	calls="$(call_lines "$roster" "$open|$close")"
	dropped="$(comment_only_lines "$roster" "$open|$close")"
	echo "# 分子口径（票 171 AC#1）＝剔整行注释＋剔行尾注释之后仍然匹配词表 '$open|$close' 的行。"
	if [ -n "$dropped" ]; then
		echo "# 被分子剔掉的行（注释里提到词表，代码里没有）："
		printf '%s\n' "$dropped" | sed 's|^|#   COMMENT-ONLY |'
	else
		echo "# 被分子剔掉的行：（空＝今天这一族没有一枚靠注释进场）"
	fi
	# 票 171 AC#6：handle 味＝先求「开方把句柄丢了」的文件清单，交给 diffsets 当强制点名那一味。
	if [ "$mode" = handle ]; then
		disc="$(discarded_files "$open" "$calls")"
		if [ -n "$disc" ]; then
			echo "# 开方丢了句柄的文件（$mode 味：合方写什么都不算关）："
			printf '%s\n' "$disc" | sed 's|^|#   DISCARD-HANDLE |'
		fi
		# 泛用 Close() 命中的文件：合方认两形之后，这一味要单独说得出来，别混进「关好了」。
		generic="$(printf '%s\n' "$calls" | grep -Ew 'Close\(\)' | cut -d: -f2 | sort -u || true)"
	fi
	# 票 161 AC#6②：求差集那段交给 diffsets()（下面（乙）文本喂入那发样本吃的是同一份算术）。
	unpaired="$(diffsets "$open" "$close" "$calls" "$disc")"
	n="$(printf '%s' "$unpaired" | grep -c . || true)"
	echo "# 成对名册基准（开过、却没在同一文件关过的文件；空＝这把尺没响）："
	if [ -n "$unpaired" ]; then
		while IFS= read -r f; do
			[ -z "$f" ] && continue
			oc="$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | grep -cxF "$f" || true)"
			# 委托关闭的形状：同一文件里出现另一族的收尾动词 ⇒ 人来判「真漏」还是「委托」
			deleg="$(git grep -lEw 'CloseTask|Defer' "$A" -- "$f" 2>/dev/null || true)"
			echo "#   UNPAIRED $f (开方调用点=${oc:-0})${deleg:+  # 同文件另见收尾动词，判「真漏／委托」}"
			if [ "$mode" = handle ] && printf '%s\n' "$disc" | grep -qxF "$f"; then
				echo "#     note 丢句柄那一味点的名：开方没接住返回值，这一族的收尾结构上不可能"
			fi
		done <<< "$unpaired"
	else
		echo "#   （空）"
	fi
	# 票 171 AC#6：合方认下句柄那一形之后，「这一文件的关是谁关的」必须说出来给人判——
	# 泛用 Close() 是本仓到处都是的动词（现量：非测试生产码 42 枚文件里有 Close()，命令＝
	#   git grep -lEw 'Close\(\)' <锚> -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l，
	#   票 171 r2 在锚 f1b99a70 重量仍是 42；改前这句写的 40＝票面 AC#8 那条「顺带一字」的账），
	#   一把只按文件求差集的尺分辨不出那只 Close() 是不是这族的句柄，所以它别想安静地混过去。
	if [ "$mode" = handle ] && [ -n "$generic" ]; then
		gk="$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | sort -u |
			while IFS= read -r f; do [ -z "$f" ] && continue
				printf '%s\n' "$generic" | grep -qxF "$f" && printf '%s\n' "$f"
			done | sort -u)"
		if [ -n "$gk" ]; then
			echo "# 合方只由泛用 Close() 命中的开方文件（这一味＝人来判那只 Close() 是不是这族句柄）："
			printf '%s\n' "$gk" | sed 's|^|#   CLOSE-BY-GENERIC-CLOSE |'
		fi
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
	echo "# 基线枚数＝${COUNT:-（没登记）}   （票 171 AC#2：响＝实测 > 基线；want_n 登记的那一发就是这一行的右边）"
	echo "# git grep rc=$rrc（1＝本射程里这一族一枚都没有——那是「射程里没这东西」，不是「都关好了」；两种 0 枚响含义不同）"
	echo "# rc=$n   （本枚子句的响＝未成对枚数；0＝不响）"
	# 票 161 AC#6①＋票 171 AC#2：pair 腿的实测交给 book_count——它拿 want_n 的基线比枚数，
	# 没登记基线就直接死（rc=4），所以这一形不可能被下一程新加的腿悄悄退回「空／非空」。
	book_count "$n" "$rrc"
}

# ---- 票 161 AC#6②（乙）文本喂入模式：只证「求差集那段算术」不撒谎 -------------------
# 名册从 stdin 进，两个方向的未成对清单出；这一枚不吃 git、不吃 $A、不吃 pathspec，
# 所以它证的不是（甲）合成树那发读数（那才是主证），它证的只是 diffsets() 本身。
# 用法：DIFFSETS_OPEN=OpenTask DIFFSETS_CLOSE=CloseTask sh gate-clauses.sh --diffsets < 名册.txt
# 名册行形状与 pair() 的逐行读数一致（path:line:text）；文件头那三行 # 是锚点标记，正常现象。
if [ "${1:-}" = "--diffsets" ]; then
	ds_roster="$(cat)"
	ds_open="${DIFFSETS_OPEN:-OpenTask}"
	ds_close="${DIFFSETS_CLOSE:-CloseTask}"
	# 票 171 AC#1：这一味也换成 call_lines——文本喂入模式与 pair() 于是吃同一份剔注释的算式，
	# 改前它自己抄了一遍"只剔整行注释"，那正是"同源拷贝逐枚有归属"那起事故的形状。
	ds_calls="$(call_lines "$ds_roster" "$ds_open|$ds_close")"
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
# 票 171 AC#7：run 腿（G1／G1b／G2／G3／G4）的 want_n 也是同一次现量登记的【命中行数】，命令逐枚是
#   git grep -nE <该腿的 pattern> <锚> -- <该腿的 pathspec 原样> | grep -c .
# 本程在锚 f1b99a70 现量：G1=0／G1b=0／G2=2／G3=0／G4=0（G2 那 2 行＝cmd/wisp/run.go 的调用点与它上面的注释，
# run 腿的枚数吃整张名册，注释行也算——这与 pair 腿"只吃调用点"是两回事，票面 AC#1 那枚尺洞不归本程）。
want G1 quiet
want_n 0
run "G1 生产码里构造 ToolRequest（排 internal/agent/ 自身）" \
	'ToolRequest\{' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/*'

want G1b quiet
want_n 0
run "G1b 生产码里任意 .Execute( 调用点（排 loop.go 与 bridge.go 自身）" \
	'\.Execute\(' '' "$GO" "$GO2" ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'

# G2 今天【该响】：名册非空是它的工作——cmd/wisp/run.go 那一发 CloseTask 就在册。
want G2 ring
# 票 171 r3 AC#1：分子换骨（run 腿从此只数代码那一半），基线在同一格现量核下来，不留 SHR 当默认。
#   改前登记的 2 枚＝名册行数（cmd/wisp/run.go:960 那一发调用点＋:942 那行整行注释）；
#   现量现口径：名册行数＝3、命中行数＝1（锚 4f2a777b，命令＝
#   sh .scratch/wisp/probes/154/gate-clauses.sh 4f2a777b | 取 G2 那一格的「# 命中行数＝」）。
#   那 3 行名册逐行原样照打，被剔掉的 2 行逐枚打在 COMMENT-ONLY 格里（读数不藏在算式里）。
want_n 1
run "G2 OpenTask/CloseTask 的生产出现点（排 bridge.go 自身；今日应只剩 cmd/wisp/run.go 那一发＋它上面的注释）" \
	'OpenTask|CloseTask' '-w' "$GO" "$GO2" ':!*_test.go' ':!internal/tools/bridge.go'

want G3 quiet
want_n 0
run "G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法" \
	'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' '' internal/agent

want G4 quiet
want_n 0
run "G4 生产组合根打开 L0 直通（PassThroughUnclassifiedRisk 的字面赋值）" \
	'PassThroughUnclassifiedRisk:' '' '*.go' ':!*_test.go'

# G5（票 158 AC#2 新增）射程＝OpenScope↔CloseScope 的成对普查。
# G1–G4 盯的是「宿主派发工具却没关环路那一枚 id」那一族（过桥那条腿），
# 而「不过桥、直接对 risk 层开一枚 scope」这一形今天有活样板：
#   cmd/wisp/panel_assets.go:243  _ = prov.OpenScope(taintSourceScopeID)   且该文件里 Close()/CloseScope 零枚。
# 排掉 internal/risk/* 的理由：OpenScope/CloseScope 的**定义**在那儿（票 158 地界明令只盘调用者），
# 把定义算成调用点会让名册每改一行注释就响一次，那是噪音不是读数。
# G5 主尺今天【该响】，且响 1 枚：cmd/wisp/panel_assets.go:243 开了 C25 scope 却没在同一文件关。
# 这一枚是票 158 登记在册的活形状、不是本票的靶子（票面 AC#6 那条红线），本程一字节生产码都没动。
# 票 171 AC#6 改的两味（射程一条没减）：合方 'CloseScope' → 'CloseScope|Close\(\)'，
#   因为票 160 之后 CloseScope 这个词根在生产码里【已经一枚都不剩】（现量：git grep -n CloseScope <锚>
#   -- internal cmd 排测试 ⇒ rc=1 零命中；唯一的收尾是 internal/tools/bridge.go:706 的 scope.Close()），
#   只写 CloseScope 的尺于是把「用句柄正常关掉」的文件一律读成漏；另一味＝handle（丢句柄强制点名）。
want G5 ring
want_n 1
pair "G5 OpenScope↔CloseScope 成对普查：生产码里开了 C25 scope 却没在同一文件关过的名册（排 internal/risk 的定义本体、排 *_test.go）" \
	'OpenScope' 'CloseScope|Close\(\)' handle "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*'

echo
echo "## G5 正控（形状照 G1–G4 的「正控先跑」：同一把尺打在已知「开又合」的样本上，必须不响）"
echo "# 样本＝internal/tools/bridge.go（:648 b.scopes[taskID] = b.prov.OpenScope(taskID) 与"
echo "#   :706 closeErr = scope.Close() 同文件成对——票 160 之后收尾是句柄自带的那枚 Close，"
echo "#   行号＝09-27 本程在锚上现量（git grep -nE 'OpenScope|scope\\.Close\\(\\)' <锚> -- internal/tools/bridge.go），"
echo "#   改前那句注释写的 :642/:691 已经因为这次落地位移过，出处＝票 171 票面 AC#6 与台账 A329）"
echo "# 写法上的两条自纠（都是本程实测撞出来的，不是设想的）："
echo "#   (1) pathspec 必须与主尺同形再逐条排除。第一版图省事写成 'internal/tools/**/*.go'，"
echo "#       实测 git grep rc=1、名册空、未成对 0——那是一枚**根本产不出读数的装饰腿**，"
echo "#       它「不响」不能算正控（票 154 的 AC#3 就是拿「零区别」否掉装饰腿的）。"
echo "#       本机现量：git grep -lEw OpenScope <锚> -- internal/tools/**/*.go => rc=1 零命中；"
echo "#                          同尺 -- internal/tools/*.go  => rc=0 命中 bridge.go。"
echo "#   (2) 必须排 *_test.go。第一版没排，正控自己被本程新测试里那句解释 Mark 机理的**注释**"
echo "#       点中——那是假阳性。成对判据从此只吃调用点（见 pair() 里剔注释那一行）。"
want G5pos quiet
want_n 0
pair "G5-正控 同尺＝主尺射程，只多排掉主尺今天点名的那一枚 panel_assets.go ⇒ 剩下的开方全是 balanced 样本" \
	'OpenScope' 'CloseScope|Close\(\)' handle "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*' ':!cmd/wisp/panel_assets.go'

echo
echo "## G5 负一负（同尺摘掉「排测试＋排定义」那两条 ⇒ 今日应非空，证明这把尺吃得住这个形状、不是结构上产不出读数的装饰）"
# G5-负一负【按设计就该响】：这一腿的作用就是证明这把尺吃得住这个形状。
# 枚数＝本程现量：改前 9 枚 → 改后 8 枚，退场的只有 internal/tools/bridge.go 那一枚
#   （它在 :648 接住句柄、在 :706 用句柄关掉——认得句柄那一形之后这一枚按设计该退场）；
#   其余 8 枚逐枚还在被点名，命令与两读见 docs/evidence/s1/171-pair-legs-r1.md §4。
want G5neg ring
want_n 8
pair "G5-负一负 同尺不过滤任何文件" \
	'OpenScope' 'CloseScope|Close\(\)' handle "$GO" "$GO2"

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
# ⚠ 票 171 r1 现量补一句：那几处行号里 bridge.go 的三枚已经被票 160 的落地挪过——
#   本程在锚上重量＝OpenTask 的调用点 :648、定义本体 :639（OpenTask）与 :689（CloseTask）、
#   收尾 :706。命令：git grep -nE 'func \(b \*Bridge\) (Open|Close)Task|OpenScope|scope\.Close\(\)' <锚> -- internal/tools/bridge.go
#   （台账 A329 那两行说的就是这种位移；这里不改 161 写下的字，只把现量并排放在旁边。）
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
want_n 1
pair "G6 主尺 OpenTask↔CloseTask 成对普查（排 *_test.go、排定义本体 internal/tools/bridge.go；两向都算）" \
	'OpenTask' 'CloseTask' - "$GO" "$GO2" ':!*_test.go' ':!internal/tools/bridge.go'

echo
echo "## G6 正控（同一把尺打在已知「开了又关」的样本上，两向都必须不响）"
echo "# 样本＝internal/tools/bridge.go 自己：:559 的 b.OpenTask 与 :684 的 CloseTask 在同一枚文件里。"
echo "# 这条射程是真吃得住的（本程现量 git grep rc=0、剔注释后 3 行调用点），不是「射程里根本没这东西」"
echo "#   那种空 0——那种 0 不算正控（票 154 AC#3 就是拿「零区别」否掉装饰腿的，G5 的自纠 (1) 同一条）。"
want G6pos quiet rev
want_n 0
pair "G6-正控 同尺只放已知「开又合」的那一枚 internal/tools/bridge.go ⇒ 两向都必须 0 枚" \
	'OpenTask' 'CloseTask' - 'internal/tools/bridge.go' ':!*_test.go'

echo
echo "## G6 负一负（同尺摘掉「排测试＋排定义」两条 ⇒ 今日应非空，证明这把尺吃得住这个形状）"
echo "# 今日真实枚数：正向 0 枚／反向 1 枚＝cmd/wisp/run.go。测试文件里 OpenTask 与 CloseTask 同文件成对"
echo "#   （internal/tools/bridge_scope_open_ticket158_test.go），所以摘掉过滤也点不出正向未成对——"
echo "#   本程不把它谎称成点得出；正向那一向的证＝（甲）合成树。"
want G6neg ring rev
want_n 1
pair "G6-负一负 同尺不过滤任何文件" \
	'OpenTask' 'CloseTask' - "$GO" "$GO2"

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
# 票 171 r3 AC#1：分子换骨之后这一腿的未成对从 3 枚掉到 1 枚，基线在同一格现量核下来（不留 SHR 当默认）。
#   退场的两枚＝internal/observe/goroutine.go 与 internal/proc/shutdown.go——它们"合方在场"的全部证据
#   是行尾注释（goroutine.go:34 / shutdown.go:85，两行逐字打在下一格的 COMMENT-ONLY 里，
#   与 161-gates-accept-r2.md:67/:68 那两行裁定同一数）。
#   留下的那一枚＝internal/memory/retention.go（合方调用点=2：:98 的签名＋:100 的字符串字面量），
#   161-v2 判它"形状使然，且是判据词表缺一味（Go 不并成方）"——本程一字节未动它，它是正控：
#   AC#1 修的是"注释算不算调用点"，不是"谁该被点名"。名册本体 19 行逐行原样照打（r3-roster-diff check 1）。
want_n 1
pair "G7 主尺 Defer↔DisposalScope 成对普查（排 *_test.go、排定义本体 internal/plugin/disposal.go；两向都算）" \
	'Defer(Named)?' '(New)?DisposalScope' - "$GO" "$GO2" ':!*_test.go' ':!internal/plugin/disposal.go'

echo
echo "## G7 正控（同一把尺打在已知「收尾与 scope 同文件」的样本上，两向都必须不响）"
echo "# 样本＝internal/plugin/disposal.go 自己：Defer／DeferNamed 与 DisposalScope 同文件 ⇒ 两向 0 枚，"
echo "#   且名册非空（本程现量 git grep rc=0），不是「射程里没这东西」那种空 0。"
want G7pos quiet rev
want_n 0
pair "G7-正控 同尺只放定义本体 internal/plugin/disposal.go 这一枚已知「两向都有」的样本" \
	'Defer(Named)?' '(New)?DisposalScope' - 'internal/plugin/disposal.go' ':!*_test.go'

echo
echo "## G7 负一负（同尺不过滤任何文件 ⇒ 今日应非空，证明这把尺吃得住这个形状）"
echo "# 今日真实枚数：正向 1 枚＝internal/risk/provenance_test.go（出现 Defer 却没在同一文件出现"
echo "#   DisposalScope 这个词）／反向 3 枚＝主尺那三枚 ⇒ 合计 4 枚。"
want G7neg ring rev
# 票 171 r3 AC#1：同上一条腿——不过滤任何文件的这一腿也从 4 枚掉到 2 枚，退场的就是那两枚注释撑起来的
#   （goroutine.go、shutdown.go），留下的两枚是真代码（retention.go 的签名＋字符串、retention_test.go 的
#   测试本体）。现量命令与逐枚位移＝sh .scratch/wisp/probes/171/r3/r3-roster-diff.sh（两把尺打同一枚锚
#   4f2a777b，名册 71 行逐字节相同＝射程未动，位移两枚都有 COMMENT-ONLY 的解释）。
want_n 2
pair "G7-负一负 同尺不过滤任何文件" \
	'Defer(Named)?' '(New)?DisposalScope' - "$GO" "$GO2"

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
echo "# 每腿判定（ok＝说到做到／BAD＝言行不一／SHR＝读数比基线好＝基线过期但不算言行不一；这一份同样逐行可 diff）："
echo "# 票 171 AC#2／AC#7：pair 腿与 run 腿现在都是同一套两列——基线＝want_n 登记的枚数、"
echo "#   实测＝这一腿自己数出来的枚数（pair 腿＝未成对，run 腿＝命中行数）。"
echo "#   退码累计的是【声明与实测不符】，而「不符」现在有三味：空/非空（票 161 AC#6①）、"
echo "#   实测>基线（票 171 AC#2／AC#7）与【腿不在场】（票 171 AC#8①）。"
# ---- 票 171 AC#8①：腿数断言。改前这一格只有「腿数＝14」一句打印，摘掉一整条腿它打印 13 然后退 0 ----
# 名册（LEGS_EXPECT）里每一腿都必须真的记过账；缺一枚就是一枚 BAD，计进退码。
# 反向也校：want 登记过、却没有对应的 run()／pair()（＝只删腿不删声明）同样是一枚 BAD。
MISSING_COUNT=0
for __e in $LEGS_EXPECT; do
	case " $BOOKED " in
	*" $__e "*) ;;
	*)
		MISSING_COUNT=$((MISSING_COUNT + 1))
		AGG="$AGG# BAD  腿=$__e 声明=在册应跑 实测=腿不在场 因=整条腿没跑（票 171 AC#8①：腿数断言）
"
		AGG_BAD=$((AGG_BAD + 1))
		;;
	esac
done
PHANTOM_COUNT=0
for __s in $LEG_SEEN; do
	case " $BOOKED " in
	*" $__s "*) ;;
	*)
		PHANTOM_COUNT=$((PHANTOM_COUNT + 1))
		AGG="$AGG# BAD  腿=$__s 声明=want 登记过 实测=只声明没跑腿 因=腿被摘掉而声明留着（票 171 AC#8①）
"
		AGG_BAD=$((AGG_BAD + 1))
		;;
	esac
done
if [ -n "$AGG" ]; then printf '%s' "$AGG"; fi
echo "# 腿数＝$LEG_COUNT 声明与实测不符＝$AGG_BAD"
echo "# 腿数断言：名册=$LEG_EXPECT_COUNT 声明=$WANT_COUNT 记账=$LEG_COUNT 缺腿=$MISSING_COUNT 空头声明=$PHANTOM_COUNT"
if [ "$LEG_COUNT" = "$LEG_EXPECT_COUNT" ] && [ "$WANT_COUNT" = "$LEG_EXPECT_COUNT" ] &&
	[ "$MISSING_COUNT" = 0 ] && [ "$PHANTOM_COUNT" = 0 ]; then
	echo "# 腿数断言＝相符（名册上每一腿都记了账）"
else
	echo "# 腿数断言＝不符（缺腿=$MISSING_COUNT 空头声明=$PHANTOM_COUNT）——已逐枚计入上面那张表的 BAD"
fi
# ---- 票 171 AC#8②：SHR（实测 < 基线）不计进退码，但单列一张读得到的【基线过期】表 ----
echo
echo "## 基线过期（票 171 AC#8②）：实测低于基线的腿，逐枚列成一张表"
echo "# 这一味**不**计进退码：变好了不算言行不一（算了就是把门换成「为好消息响」的恒红形状，台账 A320／A321）。"
echo "# 改前它只是聚合行里那句「注=读数变好了」，现量＝退码 0、单列表 0 枚（只有人眼看得见）。"
echo "# 表头：腿 | 声明 | 基线（want_n） | 实测 | 差 | 处置"
if [ -n "$STALE" ]; then printf '%s' "$STALE"; else echo "#   （空＝今天没有任何一条腿的读数低于基线）"; fi
echo "# 基线过期枚数＝$STALE_COUNT   （核它的那一把＝.scratch/wisp/probes/161/r6/flip-declaration.sh 新增的那一格）"
echo "# 聚合退码＝$AGG_BAD   （0＝每一腿的响与不响都和自己登记的声明一致，且名册上每一腿都在场；>0＝有腿言行不一）"
echo "# 翻一枚声明让退码跟着变那发证据不在这里，在 .scratch/wisp/probes/161/r6/flip-declaration.sh"
if [ "$AGG_BAD" -gt 125 ]; then
	echo "# 不符腿数大于 125，退码截在 125（够红了，别拿溢出当读数）"
	exit 125
fi
exit "$AGG_BAD"
