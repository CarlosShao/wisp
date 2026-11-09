# 2026-11-08 票411 修腿：旧射程（bash 工作树展开清单，xargs 分批绕过溢出）vs 新射程（git glob）等价性核
cd "$(git rev-parse --show-toplevel)" || exit 9
A=$(git rev-parse HEAD)
GO='internal/**/*.go'
GO2='cmd/**/*.go'
oldlist=$(mktemp "${TMPDIR:-/tmp}/wisp-411-oldscope.XXXXXX")
for w in $GO $GO2; do printf '%s\n' "$w"; done > "$oldlist"
echo "old scope files: $(wc -l < "$oldlist")"
# 旧射程命中＝分批 git grep（每批 200 枚文件 pathspec＋全量排除项），并集
oldgrep() { # $1flags $2pat $3...排除
	local flags="$1" pat="$2"; shift 2
	local out="" batch
	local batchfile; batchfile=$(mktemp "${TMPDIR:-/tmp}/wisp-411-batch.XXXXXX")
	split -l 200 "$oldlist" "$batchfile."
	for batch in "$batchfile".*; do
		git grep -nE $flags "$pat" "$A" -- $(cat "$batch") "$@" 2>/dev/null
	done
	rm -f "$batchfile".*
}
check() { # $1名 $2flags $3pat $4...排除
	local name="$1" flags="$2" pat="$3"; shift 3
	local newc oldc
	newc=$(git grep -nE $flags "$pat" "$A" -- 'internal/**/*.go' 'cmd/**/*.go' "$@" | sort -u | wc -l)
	oldc=$(oldgrep "$flags" "$pat" "$@" | sort -u | wc -l)
	local st=同; [ "$newc" = "$oldc" ] || st='DIFF'
	echo "$st $name: new=$newc old=$oldc"
}
check G1 '' 'ToolRequest\{' ':!*_test.go' ':!internal/agent/*'
check G1b '' '\.Execute\(' ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'
check G2 '-w' 'OpenTask|CloseTask' ':!*_test.go' ':!internal/tools/bridge.go'
check G5 '-Ew' 'OpenScope|CloseScope|Close\(\)' ':!*_test.go' ':!internal/risk/*'
check G6 '-Ew' 'OpenTask|CloseTask' ':!*_test.go' ':!internal/tools/bridge.go'
check G7 '-Ew' 'Defer(Named)?|(New)?DisposalScope' ':!*_test.go' ':!internal/plugin/disposal.go'
# 开方文件集合（决定 pair 腿未成对枚数的正向名册）
oldfiles() { local pat="$1"; shift; oldgrep -E "$pat" "$@" | cut -d: -f2 | sort -u; }
for spec in "OpenTask :!*_test.go" "Defer(Named)? :!*_test.go"; do
	pat=${spec%% *}
	newf=$(git grep -Ew "$pat" "$A" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | cut -d: -f2 | sort -u | md5sum | cut -d" " -f1)
	oldf=$(oldfiles "$pat" ':!*_test.go' | md5sum | cut -d" " -f1)
	newn=$(git grep -lEw "$pat" "$A" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l)
	oldn=$(oldfiles "$pat" ':!*_test.go' | wc -l)
	[ "$newf" = "$oldf" ] && echo "开方文件集合[$pat] 同 (n=$newn)" || echo "DIFF 开方文件集合[$pat]: new=$newn old=$oldn"
done
rm -f "$oldlist"
