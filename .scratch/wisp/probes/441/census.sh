#!/bin/sh
# Ticket 441 AC#4: the census ruler for self-checks shaped like
#   "the raced instrument never reached the window, therefore RED".
#
# It answers ONE question per line: which arm of a window instrument can print a red
# sentence when the window simply did not open on this run? The classification below
# reads the guard line that stands immediately above each self-check, because the
# sentence alone cannot tell a "reached-count == 0" arm (which contention can starve)
# from a fixture guard ("could not build the world", deterministic) or a racer-broken
# arm ("flipErrs > 0", which is a real defect, not a scheduling outcome).
#
# Re-run from anywhere:  sh .scratch/wisp/probes/441/census.sh
# It is read-only: git grep plus awk over the file on disk, no writes, no go test.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../.." && pwd)
cd "$root"

printf 'census root: %s\n' "$root"
printf 'ruler: git grep -n -E %s -- internal  (only *_test.go rows are kept)\n' \
    "'RED instrument|instrument is dead|RED INSTRUMENT'"

hits=$(git grep -n -E 'RED instrument|instrument is dead|RED INSTRUMENT' -- internal \
    | grep '_test\.go:' || true)
if [ -z "$hits" ]; then
    echo "census: 0 rows - the ruler is blind, refusing to report an empty roster" >&2
    exit 1
fi

printf 'rows: %s\n' "$(printf '%s\n' "$hits" | wc -l | tr -d ' ')"
printf 'files: %s\n' "$(printf '%s\n' "$hits" | cut -d: -f1 | sort -u | wc -l | tr -d ' ')"
printf 'NOTE: the arm column is a heuristic over the guard line; the per-site verdicts are in\n'
printf '      docs/evidence/s1/441-instrument-abstention.md section 4, read by hand.\n'
printf '\n%-58s %-28s %-24s %s\n' 'file:line' 'build tag' 'arm (heuristic)' 'guard immediately above'
printf '%s\n' '------------------------------------------------------------------------------------------------------------------------------------'

printf '%s\n' "$hits" | while IFS= read -r hit; do
    file=${hit%%:*}
    rest=${hit#*:}
    line=${rest%%:*}
    tag=$(sed -n '1p' "$file")
    prev=$(awk -v n="$line" 'NR < n && NF { g = $0 } END { print g }' "$file")
    case "$prev" in
        *'flipErrs'*|*'flipErr'*|*'FlipErrs'*)
            arm='racer-broken (real defect)'
            ;;
        *'== 0'*|*'!= '*|*'< 0'*|*'<= '*|*'> 0'*|*'>= '*)
            arm='reached-count (CAN abstain-red)'
            ;;
        *)
            arm='fixture/third shape (read it)'
            ;;
    esac
    short=$(printf '%s' "$prev" | sed 's/^[[:space:]]*//' | cut -c1-72)
    printf '%-58s %-28s %-24s %s\n' "$file:$line" "$tag" "$arm" "$short"
done
