#!/usr/bin/env bash
# zero156.sh <sha> [<sha>...] - ticket 156 AC#6's contract-axis census, run ONE
# COMMIT AT A TIME.
#
# The尺 is deliberately not a range diff: "git diff A..B -- <roster>" answers a
# question the ticket does not ask (it lets one commit's hit hide in another's
# clean list, and with two sibling程s in the same working tree it can attribute
# someone else's file to this ticket). Each sha below is read from its OWN commit
# object, so every number printed is one commit's number.
#
# A hit count of 0 is the claim. The roster is copied verbatim from the ticket's
# AC#6 list; go.mod/go.sum are NOT in that list and are named separately, because
# AC#2's "no new dependency" line needs them and the ticket's roster does not cover
# them.
set -u
if [ $# -lt 1 ]; then
    echo "usage: zero156.sh <sha>..." >&2
    exit 2
fi

roster=(
    '^docs/PLAN\.md$'
    '^docs/specs/'
    '^internal/risk/'
    '^internal/panel/'
    '^internal/agent/approval/'
    '^internal/observe/'
    'thresholds\.go$'
    'golden'
    '^tools/d22scan/allowlist\.txt$'
    '^scripts/slo-check\.ps1$'
    '^tools/d22scan/'
    '^frontend/'
    '^design/'
    '^docs/reports/pending-and-issues\.md$'
    '^docs/reports/HANDOVER\.md$'
)
extra=(
    '^go\.mod$'
    '^go\.sum$'
)

for sha in "$@"; do
    if ! git rev-parse --verify -q "$sha^{commit}" >/dev/null; then
        echo "zero156: $sha is not a commit in this repo - refusing to print a reading" >&2
        exit 2
    fi
    files=$(git show --name-only --format='' "$sha" | sed '/^$/d')
    total=$(printf '%s\n' "$files" | grep -c . )
    echo "=== $sha $(git log -1 --format='%s' "$sha" | cut -c1-70)"
    echo "--- touched paths ($total):"
    printf '%s\n' "$files" | sed 's/^/      /'
    i=0
    for pat in "${roster[@]}"; do
        i=$((i + 1))
        hits=$(printf '%s\n' "$files" | grep -cE -- "$pat" || true)
        printf '    AC6[%02d] %-40s hits=%s\n' "$i" "$pat" "$hits"
    done
    i=0
    for pat in "${extra[@]}"; do
        i=$((i + 1))
        hits=$(printf '%s\n' "$files" | grep -cE -- "$pat" || true)
        printf '    named-separately[%d] %-27s hits=%s\n' "$i" "$pat" "$hits"
    done
done
