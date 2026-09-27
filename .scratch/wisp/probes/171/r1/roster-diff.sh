#!/bin/sh
# roster-diff.sh - ticket 171 r1: "我有没有减射程" answered as a set difference, not as a sentence.
# It takes two gate-clauses.sh readings (the pre-fix one and the post-fix one, both produced by
# the SAME command on the SAME anchor) and prints, per pair leg, the files that only the pre-fix
# ruler named (退场) and the files only the post-fix ruler named (新增).
# Empty 新增 + 退场 only-what-the-fix-explains = no range was dropped. Anything else = report it.
#
# Usage: sh roster-diff.sh <pre-fix-reading.txt> <post-fix-reading.txt>
# rc=0 the two rosters differ by nothing outside the G5 family; rc=1 they do (look at the print).
set -eu

[ $# -eq 2 ] || { echo "usage: sh roster-diff.sh <pre.txt> <post.txt>" >&2; exit 2; }
pre=$1
post=$2
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

for leg in '## G5 OpenScope' '## G5-正控' '## G5-负一负' '## G6 主尺' '## G6-正控' '## G6-负一负' \
    '## G7 主尺' '## G7-正控' '## G7-负一负'; do
    sec() {
        awk -v want="$1" '
            /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
            ins { print }
        ' "$2"
    }
    name_of=$(printf '%s' "$leg" | sed -e 's/^## //' -e 's/ /_/g')
    sec "$leg" "$pre" | grep -E '^#   UNPAIRED' | sed -e 's/(开方调用点=[0-9]*)*.*$//' -e 's/(合方调用点=[0-9]*)*.*$//' | sort \
        >"$here/logs/roster-$name_of-pre.txt"
    sec "$leg" "$post" | grep -E '^#   UNPAIRED' | sed -e 's/(开方调用点=[0-9]*)*.*$//' -e 's/(合方调用点=[0-9]*)*.*$//' | sort \
        >"$here/logs/roster-$name_of-post.txt"
    echo "== $leg"
    printf '   改前 %s 枚／改后 %s 枚\n' \
        "$(grep -c . "$here/logs/roster-$name_of-pre.txt" || true)" \
        "$(grep -c . "$here/logs/roster-$name_of-post.txt" || true)"
    echo "   只在改前被点名（退场）:"
    comm -23 "$here/logs/roster-$name_of-pre.txt" "$here/logs/roster-$name_of-post.txt" | sed -e 's/^/     - /'
    echo "   只在改后被点名（新增）:"
    comm -13 "$here/logs/roster-$name_of-pre.txt" "$here/logs/roster-$name_of-post.txt" | sed -e 's/^/     + /'
done
