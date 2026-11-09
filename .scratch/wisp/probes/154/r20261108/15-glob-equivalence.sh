# 2026-11-09 票411：旧射程（bash 两层展开）精确等价形 :(glob) 两层 vs 新射程 :(glob) 递归（只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
A=$(git rev-parse HEAD)
cmp2() { # $1名 $2flags $3pat
	local name="$1" flags="$2" pat="$3"
	local oldn newn
	oldn=$(git grep $flags "$pat" "$A" -- ':(glob)internal/*/*.go' ':(glob)cmd/*/*.go' ':!*_test.go' | wc -l)
	newn=$(git grep $flags "$pat" "$A" -- ':(glob)internal/**/*.go' ':(glob)cmd/**/*.go' ':!*_test.go' | wc -l)
	local st=同; [ "$oldn" = "$newn" ] || st=DIFF
	echo "$st $name: old两层=$oldn new递归=$newn"
}
cmp2 G1  '-E' 'ToolRequest\{'
cmp2 G1b '-E' '\.Execute\('
cmp2 G2  '-Ew' 'OpenTask|CloseTask'
cmp2 G5  '-Ew' 'OpenScope|CloseScope|Close\(\)'
cmp2 G6  '-Ew' 'OpenTask|CloseTask'
cmp2 G7  '-Ew' 'Defer(Named)?|(New)?DisposalScope'
# 深层文件（第 3 层及以下＝旧展开盲区）逐词表命中
echo "== 盲区文件（:  (glob)internal/*/*/*.go 等 3 层以上）上的命中 =="
git grep -nE 'ToolRequest\{|\.Execute\(|OpenTask|CloseTask|OpenScope|CloseScope|Defer|DisposalScope' "$A" -- ':(glob)internal/*/*/**' ':(glob)cmd/*/*/**' ':!*_test.go' | cut -d: -f2 | sort -u | head -20
echo "盲区命中文件枚数=$(git grep -cE 'ToolRequest\{|\.Execute\(|OpenTask|CloseTask|OpenScope|CloseScope|Defer|DisposalScope' "$A" -- ':(glob)internal/*/*/**' ':(glob)cmd/*/*/**' ':!*_test.go' | wc -l)"
