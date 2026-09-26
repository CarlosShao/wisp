#!/bin/sh
# attrib.sh - ticket 161 AC#7 step 2, the attribution instrument.
# Written by 161-r4. Reads the two gofumpt rulers, never touches either of them.
#
# WHY THIS EXISTS AT ALL (the sentence AC#7 was really asking for)
# ---------------------------------------------------------------------------
# The criterion this instrument replaces read: "gofumpt -l . tools/d22scan
# tools/mockllm must be EMPTY". That reading is STRUCTURALLY unsatisfiable on a
# working tree of this repo, and no amount of formatting fixes it, because
# gofumpt walks .scratch/** (measured: `git ls-files '.scratch/**/*.go'` = 40
# files) and a forensic bench is BY DEFINITION a directory you fill with
# deliberately-broken samples. So the criterion could produce exactly two kinds of
# answer, neither of them the thing it wanted to measure: "empty", which on this
# tree is only reachable by deleting or formatting somebody else's bad sample
# (both forbidden), and "non-empty", which reads identical to a real regression.
#
# The fix is not an exemption for .scratch, not a gofumpt flag, not a wider
# allowlist, and not "the working-tree reading does not count". It is TWO RULERS,
# each of which can actually produce the reading it names:
#
#   (A) TRACKED TREE   git ls-files '*.go' | xargs gofumpt -l
#       MUST BE EMPTY. This is the form CI executes, because a CI checkout
#       contains tracked bytes and nothing else.
#   (B) WORKING TREE   gofumpt -l . tools/d22scan tools/mockllm
#       EXPECTED NON-EMPTY on a bench tree. Non-empty is fine ONLY if every line
#       can be attributed to a named ticket's intentional bad sample. A line that
#       cannot be attributed is a real regression, and this script exits 1 on it.
#
# HARD EXIT RULES (the reason this is an instrument and not a printout)
#   rc=1  (A) printed anything at all                     -> tracked bytes are dirty; CI red
#   rc=1  (B) printed a line this script cannot attribute -> real injury, chase it
#   rc=1  (B) printed a TRACKED line while (A) was empty  -> the two forms disagree
#   rc=2  the instrument could not run (no gofumpt, not in a git repo) -> no verdict
#   rc=0  (A) empty AND every (B) line attributable
#
# WHAT ATTRIBUTION MEANS HERE, EXACTLY
# A line is attributed iff its repo-relative path is
#   .scratch/wisp/probes/<NN>/<...>               AND
#   .scratch/wisp/issues/<NN>-*.md exists        (the ticket is real)
# The first half alone would be self-fulfilling (any number passes), so the second
# half is what makes "归到票 NN" a claim instead of a guess. Consequences stated
# plainly: an attributed line means "ticket <NN> owns this path"; it does NOT prove
# the bytes were written on purpose, and this script cannot tell an intentional bad
# sample from ticket <NN>'s accidental breakage. That is registered in the evidence
# file's "what this did not measure" section, not papered over.
#
# WHICH INGREDIENT, DELETED, LETS THE ORIGINAL BUG BACK IN (AC#7 asks for this
# sentence inside the instrument, not in prose)
# ---------------------------------------------------------------------------
# Delete the `classify()` call at the two call sites and everything below still
# prints the same 8 (today) lines - but no reader can tell whether any one of them
# exists in a CI checkout, which is the exact ambiguity that expired the 票 158
# AC#5③ verdict at 20:49 on 2026-09-26 (measured 20:31 empty, the fixtures that
# made it non-empty were committed 18 minutes later in 506cbae). Delete the ticket
# roster check (`ticket_known`) and (B) becomes a tautology: every bench line
# "attributes" to whatever digits its directory happens to contain, so the working
# tree can never report an injury, which is the恒真检 this repo has already refused
# twice. Delete the (A) ruler and the instrument loses its one non-negotiable
# reading, the one CI actually runs. Delete `parse_err` and rc=2 from gofumpt -
# which means "a file would not parse", NOT "these files need formatting" - is read
# as the second thing (5 of today's 8 lines are parse errors on purpose).
#
# USAGE
#   sh .scratch/wisp/probes/161/r4/attrib.sh              # run both rulers
#   sh .scratch/wisp/probes/161/r4/attrib.sh --quiet      # verdict lines + summary only
# Exit codes as listed above. UNTRACKED/IGNORED lines never change the exit code on
# their own - being attributable is what keeps them out of it.
#
# DELIBERATELY NOT DONE HERE: no exemption list, no path this script formats, no
# write of any kind outside its own log files. It reads; it does not fix.

set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

QUIET=0
if [ "${1:-}" = "--quiet" ]; then QUIET=1; fi

# ---- locate the repo root (walk up to .git; never hardcode a depth) ---------
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "attrib.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d
[ -f "$root/go.mod" ] || { echo "attrib.sh: no go.mod at derived repo root $root - refusing" >&2; exit 2; }
cd "$root"

GOPATH_BIN=$(go env GOPATH)
GOFUMPT="$GOPATH_BIN/bin/gofumpt"
[ -f "$GOPATH_BIN/bin/gofumpt.exe" ] && GOFUMPT="$GOPATH_BIN/bin/gofumpt.exe"
[ -x "$GOFUMPT" ] || { echo "attrib.sh: no gofumpt at $GOFUMPT - refusing to answer with a ruler that did not run" >&2; exit 2; }

ver() { "$GOFUMPT" --version 2>&1; }

# ---- helpers ----------------------------------------------------------------
# norm <path-as-the-tool-printed-it> -> forward slashes, no :line:col[: msg] tail
norm() {
    printf '%s\n' "$1" | tr '\\' '/' \
        | sed -e 's/:[0-9][0-9]*:[0-9][0-9]*: .*//' -e 's/:[0-9][0-9]*:[0-9][0-9]*$//'
}

# parse_err <raw line> <normalised path> -> 1 if gofumpt printed a diagnostic for
# that line instead of a bare file name (i.e. the file would not parse at all).
# Compared against norm()'s output rather than pattern-matched, because the
# diagnostic text ("expected ')', found '{'") contains no digit after its colon
# space, which is what a naive `*: *[0-9]*` test silently gets wrong.
parse_err() {
    if [ "$(printf '%s\n' "$1" | tr '\\' '/')" = "$2" ]; then printf '0'; else printf '1'; fi
}

# classify <repo-relative path> -> tracked | untracked | ignored
# This is the (A)/(B) split made per-line. Deleting it is what re-opens AC#7.
classify() {
    if git ls-files --error-unmatch -- "$1" >/dev/null 2>&1; then
        printf 'tracked'
        return
    fi
    if git check-ignore -q -- "$1" 2>/dev/null; then
        printf 'ignored'
    else
        printf 'untracked'
    fi
}

# ticket_of <repo-relative path> -> the <NN> segment, empty if the path is not a bench path
ticket_of() {
    case $1 in
    .scratch/wisp/probes/[0-9]*/*)
        t=${1#.scratch/wisp/probes/}
        printf '%s' "${t%%/*}"
        ;;
    *) printf '' ;;
    esac
}

# ticket_known <NN> -> 1 if a real ticket file carries that number, else 0
ticket_known() {
    [ -n "$1" ] || { printf '0'; return; }
    if ls "$root"/.scratch/wisp/issues/"$1"-*.md >/dev/null 2>&1; then
        printf '1'
    else
        printf '0'
    fi
}

TRACKED_DIRTY=0
UNATTRIBUTABLE=0
DISAGREE=0
B_LINES=0
B_FILES=0
B_TRACKED=0
B_TICKETS=""

# report <verdict-token> <line> - one attributed reading line
report() {
    B_LINES=$((B_LINES + 1))
    [ "$QUIET" = 1 ] || printf '   %s\n' "$1	$2"
}

# ---- ruler A: tracked tree, must be empty -----------------------------------
# The tracked-file count is read and printed FIRST, because an empty (A) is only
# evidence if the ruler was handed files to look at. Measured on this tree: 529
# tracked .go files, 40 of them under .scratch. A ruler that saw zero files and
# printed "empty" would be the恒真检 this repo has already refused twice, so it is
# a hard stop here rather than a footnote (that mistake was made writing this very
# script: `git ls-files -Z` is not a switch, the pipeline still exited 0, and the
# first reading this file ever printed was vacuously green).
TRACKED_GO=$(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ') || TRACKED_GO=0
[ "$TRACKED_GO" -gt 0 ] || { echo "attrib.sh: (A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing" >&2; exit 2; }
[ "$QUIET" = 1 ] || echo "== (A) tracked tree: git ls-files -z '*.go' | xargs -0 $GOFUMPT -l   (tracked .go files handed to it: $TRACKED_GO)"
A_RC=0
A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?
if [ "$A_RC" -gt 1 ]; then
    echo "attrib.sh: (A) gofumpt exited $A_RC on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler" >&2
    exit 2
fi
A_LINES=0
A_FILES=0
if [ -n "$A_OUT" ]; then
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        p=$(norm "$line")
        TRACKED_DIRTY=$((TRACKED_DIRTY + 1))
        who=$(git log -1 --format='%h %ad %an' --date=format:'%m-%d %H:%M' -- "$p" 2>/dev/null || echo '-')
        [ "$QUIET" = 1 ] || printf '   TRACKED-DIRTY	%s	(last commit touching it: %s)	<== CI IS RED ON THIS\n' "$p" "$who"
    done <<EOF
$A_OUT
EOF
    A_LINES=$(printf '%s\n' "$A_OUT" | grep -c . || true)
    A_FILES=$(printf '%s\n' "$A_OUT" | sed -e 's/:[0-9][0-9]*:[0-9][0-9]*\(: \|$\).*//' | tr '\\' '/' | sort -u | grep -c . || true)
fi

# ---- ruler B: working tree, expected non-empty, every line attributable ------
[ "$QUIET" = 1 ] || echo "== (B) working tree: $GOFUMPT -l . tools/d22scan tools/mockllm"
B_RC=0
B_OUT=$("$GOFUMPT" -l . tools/d22scan tools/mockllm 2>&1) || B_RC=$?
if [ -n "$B_OUT" ]; then
    seen_files=""
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        p=$(norm "$line")
        cls=$(classify "$p")
        t=$(ticket_of "$p")
        known=$(ticket_known "$t")
        pe=$(parse_err "$line" "$p")
        case $seen_files in
        *"|$p|"*) ;;
        *) seen_files="$seen_files|$p|"; B_FILES=$((B_FILES + 1)) ;;
        esac
        if [ "$cls" = tracked ]; then
            B_TRACKED=$((B_TRACKED + 1))
            if [ "$A_LINES" = 0 ]; then DISAGREE=1; fi
            report "TRACKED<CI-RED>" "$p (parse_err=$pe)"
        elif [ "$known" = 1 ]; then
            B_TICKETS="$B_TICKETS $t"
            report "ticket-$t" "$p ($cls, parse_err=$pe)"
        else
            UNATTRIBUTABLE=$((UNATTRIBUTABLE + 1))
            if [ -z "$t" ]; then
                report "UNATTRIBUTABLE" "$p (no probes/<NN>/ path segment, class=$cls, parse_err=$pe) <== REAL INJURY"
            else
                report "UNATTRIBUTABLE" "$p (path claims ticket $t but .scratch/wisp/issues/$t-*.md does not exist, class=$cls, parse_err=$pe) <== REAL INJURY"
            fi
        fi
    done <<EOF
$B_OUT
EOF
fi
B_UNIQ_TICKETS=$(printf '%s\n' $B_TICKETS | sort -u | tr '\n' ' ' | sed -e 's/ $//')

# ---- summary + hard exit -----------------------------------------------------
echo
echo "attrib.sh: gofumpt $(ver)"
echo "attrib.sh: (A) tracked-tree  lines=$A_LINES files=$A_FILES  rule: MUST be empty"
echo "attrib.sh: (B) working-tree  lines=$B_LINES files=$B_FILES  tracked=$B_TRACKED attributed_tickets=[$B_UNIQ_TICKETS] unattributable=$UNATTRIBUTABLE"
echo "attrib.sh: (B) gofumpt's own exit status was $(printf '%s' "${B_RC:-?}") - note: rc=2 from gofumpt means 'a file would not parse', not 'these files need formatting'; per-line parse_err flags above are the reading."
echo "attrib.sh: attribution rule = .scratch/wisp/probes/<NN>/** AND a real .scratch/wisp/issues/<NN>-*.md"

RC=0
[ "$TRACKED_DIRTY" -gt 0 ] && RC=1
[ "$UNATTRIBUTABLE" -gt 0 ] && RC=1
[ "$DISAGREE" = 1 ] && RC=1
if [ "$RC" = 0 ]; then
    echo "attrib.sh: GREEN - (A) empty and all $B_LINES working-tree lines attributed to a named ticket"
else
    echo "attrib.sh: RED - tracked-dirty=$TRACKED_DIRTY unattributable=$UNATTRIBUTABLE forms-disagree=$DISAGREE (see the lines above for which)"
fi
exit "$RC"
