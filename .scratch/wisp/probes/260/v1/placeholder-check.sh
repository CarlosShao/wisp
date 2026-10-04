#!/bin/sh
# 260-v1 placeholder self-check.
#
# The patterns live HERE, not in verdict.md, so that the verdict file itself can
# be measured at zero hits (a table that names the markers it forbids would
# match its own ruler).
#
# Ruler = the two shapes the orchestrator named for a delivered 交件, plus two
# same-family markers, plus an empty-cell sweep for anything left unwritten.
set -u
cd "$(dirname "$0")" || exit 2
FILE=verdict.md
if [ ! -f "$FILE" ]; then
    echo "placeholder-check.sh: verdict.md not found next to this script" >&2
    exit 2
fi

echo "subject: $FILE  ($(wc -l < "$FILE" | tr -d ' ') lines, $(wc -c < "$FILE" | tr -d ' ') bytes)"
total=0
for pat in '（待填）' '（填写中）' '待填' '填写中' 'TODO' 'TBD' '- *\[ \] *$'; do
    n=$(grep -c -- "$pat" "$FILE" 2>/dev/null || true)
    n=${n:-0}
    printf 'pattern %-14s hits=%s\n' "$pat" "$n"
    total=$((total + n))
done

# empty-cell sweep: a markdown table row whose last cell is blank
blanks=$(grep -cE '\|[[:space:]]*\|[[:space:]]*$' "$FILE" || true)
blanks=${blanks:-0}
printf 'trailing-empty-cell hits=%s\n' "$blanks"

echo "VERDICT: total pattern hits=$total, empty trailing cells=$blanks"
if [ "$total" -eq 0 ] && [ "$blanks" -eq 0 ]; then
    echo "placeholder-check.sh: CLEAN - no placeholder marker and no unwritten cell in the delivered table"
    exit 0
fi
echo "placeholder-check.sh: NOT CLEAN - the table still carries a marker or an unwritten cell" >&2
exit 1
