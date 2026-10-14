#!/bin/sh
# 票 178 AC#1 读数器：把 `git grep -nEw 'OpenTask|CloseTask' -- internal/ cmd/` 的逐行读数
# 按【文件】分成四堆：① 开且关／② 只开（真缺陷）／③ 只关（关的是别处开的＝票 178 那一形）／
# ④ 注释里提到（该文件里一枚真调用点都没有）。
# 用法：sh .scratch/wisp/probes/178/r1/ac1-census.sh [锚点]  > 读数文件
# 口径：调用行＝剔整行注释（^[[:space:]]*//、/\*、\*）＋切行尾注释（[[:space:]]// 起）之后仍
#   匹配 \bOpenTask\b／\bCloseTask\b 的行——与 .scratch/wisp/probes/154/gate-clauses.sh 的
#   call_lines()/strip_comment()（票 171 AC#1 的分子口径）同一味，逐字一致。
# 本件不做判据、不进退码：它只把射程量成可复制的名册（票 178 AC#1 那句「先量成可复制的名册」）。
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
A="${1:-$(git rev-parse HEAD)}"

strip() { sed -E -e 's|[[:space:]]//.*$||' -e 's|[[:space:]]/\*.*$||'; }
line_is_comment() { printf '%s\n' "$1" | grep -qE '^[[:space:]]*(//|/\*|\*)'; }

ROSTER="$(git grep -nEw 'OpenTask|CloseTask' "$A" -- internal/ cmd/ 2>/dev/null || true)"
RAW_N="$(printf '%s' "$ROSTER" | grep -c . || true)"
WT="$(git grep -nEw 'OpenTask|CloseTask' -- internal/ cmd/ 2>/dev/null || true)"
AN="$(printf '%s\n' "$ROSTER" | sed -e 's|^[0-9a-f][0-9a-f]*:||')"

echo "# 票 178 AC#1 四堆名册 —— 锚点 $A / 生成 $(date -Iseconds)"
echo "# 射程＝git grep -nEw 'OpenTask|CloseTask' <锚> -- internal/ cmd/（含 *_test.go、含 bridge.go 定义本体）"
if [ "$WT" = "$AN" ]; then
	echo "# 工作树 vs 锚点（去 sha 前缀后逐行比）：相同——今天 internal/ cmd/ 这两枚词的工作树读数＝锚点读数"
else
	echo "# 工作树 vs 锚点（去 sha 前缀后逐行比）：不同！以锚点那一份为准，差集如下："
	printf '%s\n' "$WT" | grep -vxF "$AN" | sed -e 's|^|#   WT-ONLY     |' || true
	printf '%s\n' "$AN" | grep -vxF "$WT" | sed -e 's|^|#   ANCHOR-ONLY |' || true
fi
echo "# 名册行数＝$RAW_N（含注释行，逐行原样）"
echo "#"

# ---- 逐行分类，落成一张 <path>|<both|open|close|comment>|<行号>|<原文> 的中间表 ------------
tmp=$(mktemp "${TMPDIR:-/tmp}/wisp178r1-ac1.XXXXXX")
printf '%s\n' "$ROSTER" | while IFS= read -r l; do
	[ -z "$l" ] && continue
	path=$(printf '%s' "$l" | cut -d: -f2)
	lnum=$(printf '%s' "$l" | cut -d: -f3)
	text=$(printf '%s' "$l" | cut -d: -f4-)
	if line_is_comment "$text"; then
		printf '%s|comment|%s|%s\n' "$path" "$lnum" "$text" >>"$tmp"
		continue
	fi
	code=$(printf '%s\n' "$text" | strip)
	o=0; c=0
	printf '%s\n' "$code" | grep -qEw 'OpenTask' && o=1
	printf '%s\n' "$code" | grep -qEw 'CloseTask' && c=1
	case "$o$c" in
	11) printf '%s|both|%s|%s\n' "$path" "$lnum" "$text" >>"$tmp" ;;
	10) printf '%s|open|%s|%s\n' "$path" "$lnum" "$text" >>"$tmp" ;;
	01) printf '%s|close|%s|%s\n' "$path" "$lnum" "$text" >>"$tmp" ;;
	00) printf '%s|comment|%s|%s\n' "$path" "$lnum" "$text" >>"$tmp" ;;
	esac
done
TMP="$tmp"

counts() { # <path> -> 开方/合方各有几枚调用行
	f=$1
	o=$(grep -F "$f|" "$TMP" | grep -cE '\|(open|both)\|' || true)
	c=$(grep -F "$f|" "$TMP" | grep -cE '\|(close|both)\|' || true)
}
pile() { # <堆号+名> <both|open|close|comment> <一行说明>
	mode="$2"
	files=$(cut -d'|' -f1 "$TMP" | sort -u | while IFS= read -r f; do
		[ -z "$f" ] && continue
		o=$(grep -F "$f|" "$TMP" | grep -cE '\|(open|both)\|' || true)
		c=$(grep -F "$f|" "$TMP" | grep -cE '\|(close|both)\|' || true)
		case "$mode" in
		both) [ "$o" -gt 0 ] && [ "$c" -gt 0 ] && printf '%s\n' "$f" ;;
		open) [ "$o" -gt 0 ] && [ "$c" -eq 0 ] && printf '%s\n' "$f" ;;
		close) [ "$o" -eq 0 ] && [ "$c" -gt 0 ] && printf '%s\n' "$f" ;;
		comment) [ "$o" -eq 0 ] && [ "$c" -eq 0 ] && printf '%s\n' "$f" ;;
		esac
	done)
	n=$(printf '%s' "$files" | grep -c . || true)
	echo "#"
	echo "# $1 ＝${n} 枚文件（$3）"
	for f in $files; do
		counts "$f"
		if [ "$mode" = comment ]; then
			echo "#   $f （只有注释行）"
		else
			echo "#   $f (开方调用点=$o 合方调用点=$c)"
		fi
		grep -F "$f|" "$TMP" | grep -E '\|(open|close|both)\|' | sed -e 's|^|#       |' || true
	done
}
pile "① 开且关（同文件既有开方又有合方调用点）" both "不是缺陷；改前那把尺也不点它"
pile "② 只开（真缺陷堆）" open "开了却没在同一文件关过；本票 AC#2 造的就是这一堆"
pile "③ 只关（关的是别处开的）" close "票 178 立案的那一形；逐枚列出供人判，票 178 判它不是缺陷"
pile "④ 注释里提到（只有注释行匹配的文件）" comment "剔注释后一枚真调用点都没有"
echo "#"
echo "# 注释里提到的【行】（逐行；与 gate-clauses.sh 各腿的 COMMENT-ONLY 格同源）＝$(grep -c '|comment|' "$TMP" || true) 行："
grep '|comment|' "$TMP" | sed -e 's|^|#   COMMENT-ONLY |' || true
unp=$(cut -d'|' -f1 "$TMP" | sort -u | while IFS= read -r f; do
	[ -z "$f" ] && continue
	o=$(grep -F "$f|" "$TMP" | grep -cE '\|(open|both)\|' || true)
	c=$(grep -F "$f|" "$TMP" | grep -cE '\|(close|both)\|' || true)
	if { [ "$o" -gt 0 ] && [ "$c" -eq 0 ]; } || { [ "$o" -eq 0 ] && [ "$c" -gt 0 ]; }; then printf 'x\n'; fi
done | grep -c . || true)
echo "# ②＋③ 合起来＝未成对文件 $unp 枚——改前 G6neg 把这两堆混成的那一枚『未成对枚数』（判据只该数②）"
