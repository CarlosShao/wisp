#!/usr/bin/env bash
# ticket 158 acceptance r1 - independent pair censuses for the two families that
# r1 §5-5 / r2 §F-7 registered as "not censused" residuals.
# Same ruler shape as .scratch/wisp/probes/154/gate-clauses.sh's pair():
#   one git grep -nEw at a PINNED anchor -> strip comment lines -> set difference
#   of files that open vs files that close. Read-only: prints to stdout only.
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
A="${1:-c7f638c}"
GO='internal/**/*.go'
GO2='cmd/**/*.go'

pair() { # $1 title $2 open $3 close $4.. pathspecs (extra excludes last)
	local title="$1" open="$2" close="$3"; shift 3
	local roster calls opens closes f unpaired="" n oc
	roster="$(git grep -nEw "$open|$close" "$A" -- $@)"; 
	calls="$(printf '%s\n' "$roster" | grep -vE ':([0-9]+):[[:space:]]*(//|/\*|\*)' || true)"
	opens="$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | sort -u || true)"
	closes="$(printf '%s\n' "$calls" | grep -Ew "$close" | cut -d: -f2 | sort -u || true)"
	echo
	echo "## $title"
	echo "\$ git grep -nEw '$open|$close' $A -- $*   (roster lines: $(printf '%s\n' "$roster" | grep -c . ), call lines: $(printf '%s\n' "$calls" | grep -c . ))"
	echo "# open call sites per file:"
	printf '%s\n' "$opens" | while IFS= read -r f; do [ -z "$f" ] && continue
		printf '#   O %-58s %s\n' "$f" "$(printf '%s\n' "$calls" | grep -Ew "$open" | cut -d: -f2 | grep -cxF "$f")"; done
	echo "# close call sites per file:"
	printf '%s\n' "$closes" | while IFS= read -r f; do [ -z "$f" ] && continue
		printf '#   C %-58s %s\n' "$f" "$(printf '%s\n' "$calls" | grep -Ew "$close" | cut -d: -f2 | grep -cxF "$f")"; done
	for f in $opens; do
		printf '%s\n' "$closes" | grep -qxF "$f" || unpaired="$unpaired$f
"
	done
	n="$(printf '%s' "$unpaired" | grep -c . || true)"
	echo "# UNPAIRED (opened in a file that never closes it):"
	if [ -n "$unpaired" ]; then printf '%s\n' "$unpaired" | sed '/^$/d' | sed 's/^/#   /'; else echo "#   (none)"; fi
	echo "# 未成对枚数＝$n"
}

echo "# anchor $A ($(git log -1 --format='%h %ad %s' --date=format:'%H:%M:%S' "$A"))"
echo
echo "## 0. token discovery: which open/close-ish verbs exist in the taint-scope area"
git grep -hoEw '(Open|Close|Begin|End|Defer|Disposal|Forget|Release)[A-Za-z0-9_]*' "$A" -- internal/risk internal/tools | sort | uniq -c | sort -rn

pair "P-1 G5 的既有射程（OpenScope<->CloseScope），本轮同尺复算" \
	'OpenScope' 'CloseScope' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*'
pair "P-2 票面 §5-5 第一族：OpenTask<->CloseTask（桥的 scope 账）" \
	'OpenTask' 'CloseTask' "$GO" "$GO2" ':!*_test.go'
pair "P-3 票面 §5-5 第二族：Defer<->DisposalScope（收尾委托那一族，排定义本体包）" \
	'Defer' 'DisposalScope' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*'
pair "P-4 反向：DisposalScope<->Defer（把上一族的开合方向倒过来再量一遍）" \
	'DisposalScope' 'Defer' "$GO" "$GO2" ':!*_test.go' ':!internal/risk/*'
