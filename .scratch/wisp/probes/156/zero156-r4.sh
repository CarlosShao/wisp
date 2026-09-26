#!/usr/bin/env bash
# zero156-r4.sh - ticket 156 AC#6's contract-axis census, run ONE COMMIT AT A TIME.
#
# Same per-commit principle as r1's landed zero156.sh (which this extends, not
# replaces): a range diff is NOT the ruler, because it lets one commit's hit hide in
# another's clean list and, with sibling programs in one shared working tree, can
# attribute someone else's file to this ticket. Every line printed here is one
# commit's own number, read from that commit object.
#
# What this version adds, because AC#6's dispatch names three instrument traps:
#   1. FALSE ZERO from `git ls-tree HEAD <dir>` without -r -> every inventory read
#      in this script's self-test is `git ls-tree -r`. (Measured: HEAD internal/risk
#      = 1 entry vs -r = 37. The non-recursive form lists the directory itself and
#      makes a live axis look clean.)
#   2. A rename records R in ONE commit, and only counts if BOTH names are examined
#      -> `git show --name-status -M` is used and the OLD and NEW path of every
#      R/C entry are both matched against the roster, so a contract-axis file moved
#      out of the axis cannot pass as a non-hit.
#   3. An empty `--diff-filter=D` does NOT mean no deleted LINES -> for every hit the
#      --numstat deletion column is printed, so "on the axis" and "how much was
#      removed" are separate columns.
#
# AC#6's roster is copied verbatim from the ticket face (:49-51). go.mod/go.sum are
# NOT on that roster; they are named separately because AC#2's "no new dependency"
# line needs them.
#
# usage:
#   zero156-r4.sh selftest            # prove trap #1 is live, count the axis
#   zero156-r4.sh <rev-list-arg>...   # e.g. 9835d81~1..HEAD  (per commit, oldest first)
set -u
root=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$root/../../../.." || exit 1
tmp="$root/zero156-r4-work"
mkdir -p "$tmp"

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

if [ "${1:-}" = "selftest" ]; then
    echo "AC6 selftest - the axis is NOT empty, so a zero on it has to be earned"
    printf '    %-42s nonrecursive=%-4s recursive=%s\n' \
        'internal/risk (trap #1 demo)' \
        "$(git ls-tree HEAD internal/risk | wc -l)" \
        "$(git ls-tree -r HEAD internal/risk | wc -l)"
    total=0
    for p in docs/PLAN.md docs/specs internal/risk internal/panel internal/agent/approval \
             internal/observe scripts/slo-check.ps1 tools/d22scan frontend design \
             docs/reports/pending-and-issues.md docs/reports/HANDOVER.md; do
        n=$(git ls-tree -r HEAD --name-only -- "$p" | grep -c .)
        total=$((total + n))
        printf '    %-42s files=%s\n' "$p" "$n"
    done
    for pat in 'thresholds\.go$' 'golden' 'allowlist\.txt$'; do
        n=$(git ls-tree -r HEAD --name-only | grep -cE -- "$pat")
        printf '    %-42s files=%s (may overlap the dirs above)\n' "/$pat/" "$n"
    done
    echo "    axis inventory (dirs above, summed, overlaps not removed) = $total"
    echo "    positive control: grep of 'golden' over the whole HEAD tree finds the"
    echo "    files the roster means, so a hits=0 line below is a measured zero."
    exit 0
fi

[ $# -ge 1 ] || { echo "usage: $0 selftest | <rev-list-arg>..." >&2; exit 2; }

sha_count=0
axis_hit_commits=0
roster_hit_total=0
for sha in $(git rev-list --reverse "$@"); do
    sha_count=$((sha_count + 1))
    short=$(git log -1 --format='%h' "$sha")
    subj=$(git log -1 --format='%s' "$sha" | cut -c1-64)
    git show --name-status -M --format='' "$sha" | awk 'NF' > "$tmp/$short.statused.txt"
    git show --numstat    -M --format='' "$sha" | awk 'NF' > "$tmp/$short.numstat.txt"

    # every path this commit touched, INCLUDING both names of an R entry
    awk '{ if ($1 ~ /^R/) { print $2; print $3 } else if ($1 ~ /^[ACDMTUX]/) { print $2 } }' \
        "$tmp/$short.statused.txt" | sort -u > "$tmp/$short.paths.txt"

    hits=0
    line=""
    i=0
    for pat in "${roster[@]}"; do
        i=$((i + 1))
        n=$(grep -cE -- "$pat" "$tmp/$short.paths.txt" || true)
        if [ "$n" -gt 0 ]; then
            hits=$((hits + n))
            line="$line AC6[$(printf '%02d' "$i")]=$n"
        fi
    done
    exline=""
    i=0
    for pat in "${extra[@]}"; do
        i=$((i + 1))
        n=$(grep -cE -- "$pat" "$tmp/$short.paths.txt" || true)
        [ "$n" -gt 0 ] && exline="$exline go-dep[$i]=$n"
    done
    # deleted lines on this commit (any path) - trap #3
    delsum=$(awk '{ s += ($1 == "-" ? 0 : $2) } END { print s + 0 }' "$tmp/$short.numstat.txt")

    if [ "$hits" -gt 0 ]; then
        axis_hit_commits=$((axis_hit_commits + 1))
        roster_hit_total=$((roster_hit_total + hits))
        printf '%-9s %-3s AC6-HIT%s  files=%-3s deletedlines=%-5s | %s\n' \
            "$short" "$sha_count" "$line" "$(wc -l < "$tmp/$short.paths.txt")" "$delsum" "$subj"
        for pat in "${roster[@]}"; do
            grep -HE -- "$pat" "$tmp/$short.paths.txt" | sed "s|^|            axis| " 
        done | sed "s|$tmp/$short.paths.txt:||" | sort -u
    else
        printf '%-9s %-3s AC6-zero%s%-14s files=%-3s deletedlines=%-5s | %s\n' \
            "$short" "$sha_count" "$exline" "" "$(wc -l < "$tmp/$short.paths.txt")" "$delsum" "$subj"
    fi
done
echo "=== AC#6 census: $sha_count commits read one at a time; $axis_hit_commits of them touch the contract axis; $roster_hit_total axis-path hits total ==="
