# 2026-11-08 票411 修腿：旧射程（bash 工作树展开清单）vs 新射程（git glob）等价性核（只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
A=$(git rev-parse HEAD)
GO='internal/**/*.go'
GO2='cmd/**/*.go'
# 旧射程清单＝改前 bash 未引号展开的真实形状（1084 枚工作树文件）＋字面排除项
oldlist=$(mktemp "${TMPDIR:-/tmp}/wisp-411-oldscope.XXXXXX")
for w in $GO $GO2; do printf '%s\n' "$w"; done > "$oldlist"
echo "old scope files: $(wc -l < "$oldlist")"
# 腿定义：名|flags|pattern|额外排除（与脚本调用点逐字一致）
check() { # $1名 $2flags $3pat $4...排除
	local name="$1" flags="$2" pat="$3"; shift 3
	local newc oldc
	newc=$(git grep -nE $flags "$pat" "$A" -- 'internal/**/*.go' 'cmd/**/*.go' "$@" | wc -l)
	{ cat "$oldlist"; printf '%s\n' "$@"; } > "$oldlist.full"
	oldc=$(git grep -nE $flags --pathspec-from-file="$oldlist.full" "$pat" "$A" | wc -l)
	local st=同; [ "$newc" = "$oldc" ] || st='≠≠≠'
	echo "$st $name: new=$newc old=$oldc (排除: $*)"
}
check G1 '' 'ToolRequest\{' ':!*_test.go' ':!internal/agent/*'
check G1b '' '\.Execute\(' ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'
check G2 '-w' 'OpenTask|CloseTask' ':!*_test.go' ':!internal/tools/bridge.go'
check G5 '-Ew' 'OpenScope|CloseScope|Close\(\)' ':!*_test.go' ':!internal/risk/*'
check G6 '-Ew' 'OpenTask|CloseTask' ':!*_test.go' ':!internal/tools/bridge.go'
check G7 '-Ew' 'Defer(Named)?|(New)?DisposalScope' ':!*_test.go' ':!internal/plugin/disposal.go'
# 未成对（文件差集）层面的核：G6 与 G7 的开方文件集合
for pat in 'OpenTask' 'Defer(Named)?'; do
	newf=$(git grep -Ew "$pat" "$A" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | cut -d: -f2 | sort -u | md5sum)
	oldf=$(git grep -Ew --pathspec-from-file="$oldlist" "$pat" "$A" | cut -d: -f2 | sort -u | md5sum)
	[ "$newf" = "$oldf" ] && echo "开方文件集合[$pat] 同 ($newf)" || echo "≠≠≠ 开方文件集合[$pat] 不同 new=$newf old=$oldf"
done
rm -f "$oldlist" "$oldlist.full"
