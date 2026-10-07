#!/bin/sh
# attrib.sh - ticket 161 AC#7 step 2, the attribution instrument, r5 copy.
# Written by 161-r5 as a COPY of .scratch/wisp/probes/161/r4/attrib.sh (that file
# is another program's committed artifact: read-only here). Nothing in r4 was
# edited; every rule below that r4 already had is the same rule, same text.
#
# WHAT r5 ADDED ON TOP OF r4, AND WHY EACH ADDITION IS NOT A WEAKENING
# ---------------------------------------------------------------------------
#   --classify-line <text>  the classifier, drivable from TEXT instead of disk.
#   --self-test             feeds the classifier known lines and requires BOTH
#                           directions: lines that must ring and lines that must
#                           stay silent. This is what lets the physical negative
#                           control stop being load-bearing (see r5/README.md).
#   --tracked-only          runs ruler (A) alone, with the same hollow guard.
#                           This is the shape CI executes.
#   negative-control sample copied to probes/161/r5/negative-control/, so a
#   deliberately-unattributable .go file now sits under a path that IS attributable
#   to a real ticket, and the on-disk injury the r4 program left at
#   probes/999/bad-sample.go is no longer the only proof that (B) can ring.
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
#   rc=1  --self-test and any case read the wrong way     -> the classifier moved
#   rc=2  the instrument could not run (no gofumpt, not in a git repo) -> no verdict
#   rc=2  the (A) ruler was handed ZERO files             -> refuse to call that green
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
# prints the same (today) 9 lines - but no reader can tell whether any one of them
# exists in a CI checkout, which is the exact ambiguity that expired the 票 158
# AC#5③ verdict at 20:49 on 2026-09-26 (measured 20:31 empty, the fixtures that
# made it non-empty were committed 18 minutes later in 506cbae). Delete the ticket
# roster check (`ticket_known`) and (B) becomes a tautology: every bench line
# "attributes" to whatever digits its directory happens to contain, so the working
# tree can never report an injury, which is the 恒真检 this repo has already refused
# twice. Delete the (A) ruler and the instrument loses its one non-negotiable
# reading, the one CI actually runs. Delete `parse_err` and rc=2 from gofumpt -
# which means "a file would not parse", NOT "these files need formatting" - is read
# as the second thing (4 of today's 9 lines are parse errors on purpose). Delete
# the `TRACKED_GO -gt 0` guard and a ruler handed nothing reports "empty", which is
# the mistake this file's own ancestor made once (see r4 header: `git ls-files -Z`
# is not a switch). Delete either direction of `--self-test` and the classifier is
# back to being believed rather than checked.
#
# USAGE
#   sh .scratch/wisp/probes/161/r5/attrib.sh                       # both rulers
#   sh .scratch/wisp/probes/161/r5/attrib.sh --quiet               # verdicts + summary
#   sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only        # what CI runs (form A)
#   sh .scratch/wisp/probes/161/r5/attrib.sh --self-test           # classifier, both directions
#   sh .scratch/wisp/probes/161/r5/attrib.sh --classify-line 'some/path.go'   # one line, rc 0/1
# Exit codes as listed above. UNTRACKED/IGNORED lines never change the exit code on
# their own - being attributable is what keeps them out of it.
#
# DELIBERATELY NOT DONE HERE: no exemption list, no path this script formats, no
# write of any kind outside its own log files. It reads; it does not fix.

set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

QUIET=0
MODE=both
CLASS_INPUT=""

while [ $# -gt 0 ]; do
    case $1 in
    --quiet) QUIET=1 ;;
    --tracked-only) MODE=tracked ;;
    --self-test) MODE=selftest ;;
    --classify-line)
        MODE=classify
        [ $# -ge 2 ] || { echo "attrib.sh: --classify-line needs one argument" >&2; exit 2; }
        CLASS_INPUT=$2
        shift
        ;;
    *) echo "attrib.sh: unknown argument '$1' - refusing to guess what was asked" >&2; exit 2 ;;
    esac
    shift
done

# ---- locate the repo root (walk up to .git; never hardcode a depth) ---------
# ATTRIB_ROOT is the only r5 addition to the root derivation, and it exists for
# exactly one reason: cell 3 of 161-r5 has to run ruler (A) inside an OUT-OF-REPO
# copy of the tracked set (proof that TRACKED-DIRTY can ring without anybody
# dirtying a shared tree). It is not an exemption either - the copy has its own
# .git and its own index, so the same git questions get asked there.
d=${ATTRIB_ROOT:-$here}
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "attrib.sh: no .git above $d - refusing to answer with a ruler that did not run" >&2; exit 2; }
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

# ---- the classifier, as one place that decides ------------------------------
# classify_line <raw line gofumpt printed>
#   CL_VERDICT : "TRACKED<CI-RED>" | "ticket-<NN>" | "UNATTRIBUTABLE"
#   CL_DETAIL  : why, including the git class and gofumpt's own parse_err flag
#   CL_RC      : 0 = attributed to a named ticket (this line is EXPECTED on a
#                bench tree, it does not ring)
#                1 = this line is NOT attributable, or it is tracked while the
#                    tracked ruler said empty -> REAL INJURY, the caller exits 1
# r4 decided this inline in the (B) loop; r5 factored it out so that the SAME
# function answers a line fed as text (--classify-line, --self-test) and the lines
# read off disk. It is one rule with two feeds, not two rules.
classify_line() {
    cl_path=$(norm "$1")
    cl_cls=$(classify "$cl_path")
    cl_ticket=$(ticket_of "$cl_path")
    cl_known=$(ticket_known "$cl_ticket")
    cl_pe=$(parse_err "$1" "$cl_path")
    CL_PATH=$cl_path
    CL_PARSE_ERR=$cl_pe
    CL_CLASS=$cl_cls
    if [ "$cl_cls" = tracked ]; then
        CL_VERDICT="TRACKED<CI-RED>"
        CL_DETAIL="$cl_path (parse_err=$cl_pe) <== TRACKED BYTES ARE DIRTY"
        CL_RC=1
        return
    fi
    if [ "$cl_known" = 1 ]; then
        CL_VERDICT="ticket-$cl_ticket"
        CL_DETAIL="$cl_path ($cl_cls, parse_err=$cl_pe)"
        CL_RC=0
        return
    fi
    CL_VERDICT="UNATTRIBUTABLE"
    if [ -z "$cl_ticket" ]; then
        CL_DETAIL="$cl_path (no probes/<NN>/ path segment, class=$cl_cls, parse_err=$cl_pe) <== REAL INJURY"
    else
        CL_DETAIL="$cl_path (path claims ticket $cl_ticket but .scratch/wisp/issues/$cl_ticket-*.md does not exist, class=$cl_cls, parse_err=$cl_pe) <== REAL INJURY"
    fi
    CL_RC=1
}

# ---- mode: classify one line of text, then hard-exit ------------------------
if [ "$MODE" = classify ]; then
    classify_line "$CLASS_INPUT"
    printf '%s\t%s\n' "$CL_VERDICT" "$CL_DETAIL"
    exit "$CL_RC"
fi

# ---- mode: self-test the classifier on synthetic text -----------------------
# WHY BOTH DIRECTIONS ARE IN THE SAME TEST
# A check that rings on whatever you feed it is worth exactly as much as a check
# that never rings (this repo has refused that shape three times, most loudly at
# r4's own discovery that probes/161/r3/gate-order.sh prints `clean on the
# committed set` and rc=0 while all three of its rulers misfire). So every
# expectation below is paired: a line that MUST be called an injury, and a line
# that MUST be called a ticket's own bench artifact, decided by the same
# classify_line(). The "silent" cases are not vacuous either - each one names the
# verdict token it must produce, so a classifier that answered "ticket-999" for
# the 999 line would fail just as loudly as the reverse.
#
# WHAT THIS REPLACES ON DISK
# The r4 program proved "the ruler can ring" by leaving a real wound at
# .scratch/wisp/probes/999/bad-sample.go (unattributable by design, and by this
# repo's create-only rule nobody below the owner may delete it). That made form
# (B) red on every run, forever - a permanently red instrument is as useless as a
# permanently green one. The ring now comes from text, so it is reproducible on a
# fresh clone, and the physical sample it used to require has been copied to
# r5/negative-control/, a path that IS attributable.
if [ "$MODE" = selftest ]; then
    ST_FAIL=0
    ST_RUN=0
    # st_case <expect: ring|silent> <expected verdict token> <raw line>
    st_case() {
        want=$1
        want_verdict=$2
        line=$3
        ST_RUN=$((ST_RUN + 1))
        classify_line "$line"
        got=$CL_RC
        ok=1
        [ "$got" = "$want" ] || ok=0
        [ "$CL_VERDICT" = "$want_verdict" ] || ok=0
        if [ "$ok" = 1 ]; then
            printf '   PASS  want rc=%s verdict=%-17s got rc=%s  %s\n' "$want" "$want_verdict" "$got" "$CL_DETAIL"
        else
            ST_FAIL=$((ST_FAIL + 1))
            printf '   FAIL  want rc=%s verdict=%-17s got rc=%s verdict=%s\n' "$want" "$want_verdict" "$got" "$CL_VERDICT"
            printf '         line: %s\n         why : %s\n' "$line" "$CL_DETAIL"
        fi
    }
    echo "== attrib.sh --self-test: the classifier, fed as text ($(ver))"
    # --- must RING (this is the reading the disk sample used to be the only proof of)
    st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/999/bad-sample.go'
    st_case 1 UNATTRIBUTABLE '.scratch\wisp\probes\999\bad-sample.go:19:2: expected declaration, found '\'
    st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'
    st_case 1 UNATTRIBUTABLE 'internal/tools/not_a_bench_path.go'
    # A tracked path is the form-(A) reading: it must ring here too, because this
    # mode answers "would this line be a silent bench artifact?" and tracked bytes
    # are the opposite of that. tools/d22scan/main.go is a real tracked file, so
    # this case asks the git index a question instead of a fabricated path.
    st_case 1 'TRACKED<CI-RED>' 'tools/d22scan/main.go'
    # --- must stay SILENT (attributable to a real, named ticket)
    st_case 0 ticket-161 '.scratch/wisp/probes/161/r5/negative-control/bad-sample.go'
    st_case 0 ticket-161 '.scratch\wisp\probes\161\r5\negative-control\bad-sample.go:19:2: expected declaration, found '\'
    # A second, DIFFERENT real ticket, so that "attributable" cannot be quietly
    # hardcoded to this program's own number. Every line in this mode is text -
    # that is the whole point of the mode - and what this case checks is the roster
    # half of the rule (does .scratch/wisp/issues/169-*.md exist), not whether some
    # file sits at the path. The 777 ring case above is its pair: a bench path whose
    # number has no ticket file must ring even though it has the right shape.
    st_case 0 ticket-169 '.scratch/wisp/probes/169/a-bench-sample.go'
    # --- AC#6 (275-r2): a behavioural positive control for the rc=2 roster print ---
    # WHY THE TEXT CASES ABOVE CANNOT COVER THIS. Every case so far feeds the
    # classifier a LINE OF TEXT, so not one of them can notice that ruler_a's
    # `A_RC -gt 1` branch stopped printing rows: 275-v1 commented that print block
    # out and read `--self-test cases=8 failures=0` GREEN while stdout went from 23
    # rows to 1 (verdict.md 2(i) and its logs/selftest-neutered.txt). "It prints" was
    # never the claim; "somebody checks that it prints" is. So this case does not
    # trust the printout - it runs the real rc=2 path and asserts on its stdout.
    # WHAT IT RUNS. A throwaway git tree under mktemp (the shape
    # tracked-dirty-proof.sh:46-52 already uses, and it asserts the path is outside
    # the repo rather than trusting the name), holding exactly two tracked .go
    # samples: one that does not parse, and one that parses but is not gofumpt-clean
    # (the two-space composite-literal indent copied from tracked-dirty-proof.sh:113,
    # a shape that repo already measured as TRACKED-DIRTY). go.mod and .gitattributes
    # travel with it: go.mod because this script's own root check requires it (:134),
    # .gitattributes so the tree normalises *.go the way a CI checkout does
    # (same reason as tracked-dirty-proof.sh:67-73). The tree is then handed to THIS
    # script via ATTRIB_ROOT (:127, the instrument's own out-of-repo mechanism) with
    # --tracked-only, i.e. the mode CI runs.
    # WHAT MUST HOLD. rc=2 (the guarded branch really executed) AND at least one
    # 'A-ROSTER' row on that run's STDOUT naming the unformatted sample. Neutralise
    # the print block and the row count reads 0, this case FAILS loudly with the
    # captured output, and the self-test exits 1.
    # WHAT IT DOES NOT PROVE (so nobody over-reads it). Not that the roster is
    # COMPLETE over the real 935-file denominator - that is 275-v1's directed
    # ordering experiment, not this case; not that the rows are well-formatted; not
    # anything about the rc=1 path (this tree is built so that gofumpt exits >=2). It
    # is sensitive to the print block being deleted, redirected to stderr, silenced,
    # or the A-ROSTER token renamed; it is blind to the guard, the exit code and the
    # denominator, all of which other cases and this file's own guards carry.
    # NOTE ON WRITES AND ON A MISSING gofumpt: this is the first mode in the file
    # that writes anything at all - only inside that throwaway tree, outside the
    # repo, and it is left ON DISK when it ends (create-only, as
    # tracked-dirty-proof.sh:24-26 does); its path is printed either way. If gofumpt
    # cannot be found this mode is never reached: the startup check at :140 exits 2
    # with a named reason, which is a loud no and not a silent skip.
    ST_RUN=$((ST_RUN + 1))
    CTL_BROKEN=rosterctl/zz_roster_ctl_broken.go
    CTL_UNFORMATTED=rosterctl/zz_roster_ctl_unformatted.go
    ctl_dir=""
    ctl_note=""
    ctl_rc=0
    ctl_rows=0
    if ctl_dir=$(mktemp -d "${TMPDIR:-/tmp}/wisp-275-r2-rosterctl.XXXXXX" 2>/dev/null); then
        case $ctl_dir/ in
        "$root"/*) ctl_note="the throwaway tree $ctl_dir landed INSIDE the repo $root - that shape is banned here, refusing to run it" ;;
        esac
    else
        ctl_dir=""
        ctl_note="mktemp -d failed - the throwaway tracked tree could not be built, so the roster print could not be run at all"
    fi
    if [ -z "$ctl_note" ] && { [ ! -f "$root/go.mod" ] || [ ! -f "$root/.gitattributes" ]; }; then
        ctl_note="$root is missing go.mod or .gitattributes - the throwaway tree cannot be made to look like a checkout"
    fi
    if [ -z "$ctl_note" ]; then
        if ! { mkdir -p "$ctl_dir/rosterctl" \
            && cp -- "$root/go.mod" "$root/.gitattributes" "$ctl_dir"/ \
            && printf 'package rosterctl\n\nfunc rosterCtlBroken( {\n    return 1\n}\n' > "$ctl_dir/$CTL_BROKEN" \
            && printf 'package rosterctl\n\nvar rosterCtlDirty = map[string]int{\n  "two space indent is not gofumpt": 1,\n}\n' > "$ctl_dir/$CTL_UNFORMATTED"; }; then
            ctl_note="could not populate the throwaway tree $ctl_dir (go.mod, .gitattributes or the two samples)"
        fi
    fi
    if [ -z "$ctl_note" ]; then
        if ! ( cd "$ctl_dir" && git init -q . && git add -- go.mod .gitattributes "$CTL_BROKEN" "$CTL_UNFORMATTED" ) >/dev/null 2>&1; then
            ctl_note="git init or git add failed in $ctl_dir - the two samples are not tracked there, so ruler (A) would be handed nothing"
        fi
    fi
    if [ -z "$ctl_note" ]; then
        ATTRIB_ROOT=$ctl_dir sh "$here/$(basename -- "$0")" --tracked-only \
            > "$ctl_dir/selftest-roster-stdout.txt" 2> "$ctl_dir/selftest-roster-stderr.txt" \
            || ctl_rc=$?
        ctl_rows=$(tr -d '\r' < "$ctl_dir/selftest-roster-stdout.txt" | tr '\\' '/' | grep -c "A-ROSTER.*$CTL_UNFORMATTED" || true)
        if [ "$ctl_rc" != 2 ]; then
            ctl_note="the throwaway tree exited $ctl_rc, not 2 - the A_RC -gt 1 branch this case exists to check never ran (roster rows naming the sample: $ctl_rows)"
        elif [ "$ctl_rows" -lt 1 ]; then
            ctl_note="rc was 2 but stdout carries no A-ROSTER row naming $CTL_UNFORMATTED - the print block did not print, which is exactly the red-with-zero-names shape 275 was filed for"
        fi
    fi
    if [ -z "$ctl_note" ]; then
        printf '   PASS  want rc=2+roster-names-the-sample  got rc=%s rows=%s  throwaway tracked tree left at: %s\n' "$ctl_rc" "$ctl_rows" "$ctl_dir"
    else
        ST_FAIL=$((ST_FAIL + 1))
        printf '   FAIL  want rc=2+roster-names-the-sample  got rc=%s rows=%s\n' "$ctl_rc" "$ctl_rows"
        printf '         why : %s\n' "$ctl_note"
        if [ -n "$ctl_dir" ]; then
            printf '         tree: %s (its stdout, then its stderr:)\n' "$ctl_dir"
            tail -5 "$ctl_dir/selftest-roster-stdout.txt" 2>/dev/null | sed 's/^/         out: /' || true
            tail -3 "$ctl_dir/selftest-roster-stderr.txt" 2>/dev/null | sed 's/^/         err: /' || true
        fi
    fi
    echo "attrib.sh: --self-test cases=$ST_RUN failures=$ST_FAIL"
    if [ "$ST_FAIL" = 0 ]; then
        echo "attrib.sh: SELF-TEST GREEN - the classifier rings on a line it cannot attribute and stays silent on a line a real ticket owns"
        exit 0
    fi
    echo "attrib.sh: SELF-TEST RED - the classifier no longer reads the way AC#7 was decided; do not 'fix' the expectations, fix the classifier"
    exit 1
fi

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

# ruler_a - form (A): tracked bytes only, must be empty.
# The tracked-file count is read and printed FIRST, because an empty (A) is only
# evidence if the ruler was handed files to look at. A ruler that saw zero files and
# printed "empty" would be the 恒真检 this repo has already refused twice, so it is
# a hard stop here rather than a footnote (that mistake was made writing the r4
# version of this file: `git ls-files -Z` is not a switch, the pipeline still exited
# 0, and the first reading it ever printed was vacuously green).
ruler_a() {
    TRACKED_GO=$(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ') || TRACKED_GO=0
    [ "$TRACKED_GO" -gt 0 ] || { echo "attrib.sh: (A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing" >&2; exit 2; }
    [ "$QUIET" = 1 ] || echo "== (A) tracked tree: git ls-files -z '*.go' | xargs -0 $GOFUMPT -l   (tracked .go files handed to it: $TRACKED_GO)"
    A_RC=0
    A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?
    if [ "$A_RC" -gt 1 ]; then
        echo "attrib.sh: (A) gofumpt exited $A_RC on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler" >&2
        # 275-r1 / A663 (shape 丁): print the per-file roster gofumpt DID emit (A_OUT)
        # to stdout BEFORE exiting, so this step is red WITH names instead of red with
        # zero names (the rc=2 path used to exit here, before the loop below ever ran).
        # Nothing about the verdict moves: exit code stays 2, the >1 guard stays, the
        # denominator stays every tracked *.go, .scratch/** is not excluded, and no
        # sample is touched. gofumpt -l lists every file that needs formatting on stdout
        # and prints the parse-error diagnostic on stderr even while exiting non-zero, so
        # A_OUT already carries one row per flagged file plus that diagnostic; this only
        # makes those already-computed rows visible. Raw as gofumpt printed it (repo-
        # relative path; backslashes on Windows, forward slashes in a CI checkout).
        if [ -n "$A_OUT" ]; then
            echo "== (A) roster: every row of A_OUT printed raw, and A_OUT was captured with 2>&1, so these rows are gofumpt's stdout list AND its stderr parse-error diagnostics folded into one stream (count them as rows, not as files); the whole capture is here, it is not a prefix cut off at the unparseable file. exit code unchanged, denominator unchanged:"
            while IFS= read -r roster_line; do
                [ -n "$roster_line" ] || continue
                printf '   A-ROSTER	%s\n' "$roster_line"
            done <<EOF
$A_OUT
EOF
        fi
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
}

ruler_b() {
    [ "$QUIET" = 1 ] || echo "== (B) working tree: $GOFUMPT -l . tools/d22scan tools/mockllm"
    B_RC=0
    B_OUT=$("$GOFUMPT" -l . tools/d22scan tools/mockllm 2>&1) || B_RC=$?
    if [ -n "$B_OUT" ]; then
        seen_files=""
        while IFS= read -r line; do
            [ -n "$line" ] || continue
            classify_line "$line"
            p=$CL_PATH
            case $seen_files in
            *"|$p|"*) ;;
            *) seen_files="$seen_files|$p|"; B_FILES=$((B_FILES + 1)) ;;
            esac
            if [ "$CL_VERDICT" = "TRACKED<CI-RED>" ]; then
                B_TRACKED=$((B_TRACKED + 1))
                if [ "$A_LINES" = 0 ]; then DISAGREE=1; fi
                report "TRACKED<CI-RED>" "$p (parse_err=$CL_PARSE_ERR)"
            elif [ "$CL_RC" = 0 ]; then
                B_TICKETS="$B_TICKETS ${CL_VERDICT#ticket-}"
                report "$CL_VERDICT" "$p ($CL_CLASS, parse_err=$CL_PARSE_ERR)"
            else
                UNATTRIBUTABLE=$((UNATTRIBUTABLE + 1))
                report "UNATTRIBUTABLE" "$CL_DETAIL"
            fi
        done <<EOF
$B_OUT
EOF
    fi
    B_UNIQ_TICKETS=$(printf '%s\n' $B_TICKETS | sort -u | tr '\n' ' ' | sed -e 's/ $//')
}

summary() {
    echo
    echo "attrib.sh: gofumpt $(ver)"
    echo "attrib.sh: (A) tracked-tree  lines=$A_LINES files=$A_FILES  rule: MUST be empty"
}

# ---- mode: tracked-only - what CI runs -------------------------------------
# Same ruler, same hollow guard, same exit code; form (B) is not consulted,
# because (B) is about a bench working tree and a CI checkout has no bench state
# to attribute. That asymmetry is the whole reason AC#7 was re-read as two forms.
if [ "$MODE" = tracked ]; then
    ruler_a
    summary
    echo "attrib.sh: (B) not run in --tracked-only mode - a checkout has no bench working tree to attribute"
    echo "attrib.sh: attribution rule = .scratch/wisp/probes/<NN>/** AND a real .scratch/wisp/issues/<NN>-*.md"
    RC=0
    [ "$TRACKED_DIRTY" -gt 0 ] && RC=1
    if [ "$RC" = 0 ]; then
        echo "attrib.sh: GREEN (tracked-only) - (A) empty over $TRACKED_GO tracked .go files"
    else
        echo "attrib.sh: RED (tracked-only) - tracked-dirty=$TRACKED_DIRTY files_dirty=$A_FILES. The CI gofmt step and this step read the same denominator; this one additionally refuses to answer when handed 0 files."
    fi
    exit "$RC"
fi

# ---- both rulers -----------------------------------------------------------
ruler_a
ruler_b

# ---- summary + hard exit -----------------------------------------------------
B_UNIQ_TICKETS=$(printf '%s\n' $B_TICKETS | sort -u | tr '\n' ' ' | sed -e 's/ $//')
summary
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
