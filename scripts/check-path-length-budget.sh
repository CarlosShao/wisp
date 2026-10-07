#!/bin/sh
# scripts/check-path-length-budget.sh - the tracked-path-length gate (ticket 262).
#
# WHAT THIS MEASURES AND WHY IT EXISTS
#   Nothing in this repository looked at the length of a TRACKED PATH before this
#   script. Two ticket file names grew past the Windows path limit and
#   `actions/checkout` on the self-hosted runner died with
#   `##[error]error: unable to create file .scratch/wisp/issues/256-...` - which
#   takes the whole `slo-full` job with it, and `slo-full` is the only place the
#   two D32 resource numbers are ever evaluated. The failure surfaced as a red job
#   whose step name says nothing about path length. Ticket 262 exists because the
#   next such name would otherwise arrive just as silently.
#   Census: .scratch/wisp/probes/262/a1/census.md ; every number this file cites
#   was re-run by 262-r1 in .scratch/wisp/probes/262/r1/census.md section 1.
#
# WHO CALLS THIS, AND WHEN (ticket 262 AC#1 - two call sites, both named)
#   1. CI, shape "丁" (the half that runs whether or not anyone remembers to):
#      file .github/workflows/ci.yml, job `lint`, step named
#      "Tracked path-length budget (ticket 262)" - the step declaration is at
#      .github/workflows/ci.yml:136 and its command line, which is the call site
#      proper, is .github/workflows/ci.yml:166.
#      Command: sh scripts/check-path-length-budget.sh --with-self-test
#      Reached on every push to main or dev, on every pull_request, and on the
#      nightly schedule defined in that same file.
#   2. Local pre-push, shape "丙" (the half that fires before anything reaches CI):
#      The orchestrator - the only party holding push rights
#      (.scratch/wisp/issues/README.md rules 1 and 2) - runs the same command
#      immediately before each `git push`.
#      This is a DOCUMENTED invocation rather than an installed hook, on purpose:
#      this repository has no `.githooks/` directory and `git config --get
#      core.hooksPath` exits 1 (measured 2026-10-04), and installing a hook is a
#      git-config change that the ticket-262 unfreeze record
#      (docs/reports/pending-and-issues.md A585) does not grant - A585 authorises
#      ADDING ONE STEP to .github/workflows/ci.yml and nothing else.
#
# WHY THE CI STEP IS NOT PLACED ON THE SELF-HOSTED RUNNER (the dead-gate shape)
#   The obvious placement was "put the ruler where the wall is", i.e. hang it
#   after `actions/checkout` in `slo-full` on the Windows self-hosted runner. That
#   is a door that can only open when it is not needed: a step after checkout
#   executes only if checkout succeeded, and the one path shape that makes
#   checkout fail is exactly the path this ruler exists to catch. The orchestrator
#   ruled it out in docs/reports/pending-and-issues.md A585, retracting the
#   sentence it had itself written into ticket 262 section "what to build".
#   So the runner-side width is reproduced by a NAMED CONSTANT (WORST_PREFIX)
#   instead of by standing on that machine. Stated consequence: this step runs on
#   ubuntu-latest, whose checkout directory is /home/runner/work/wisp/wisp (27
#   characters), and the hosted Windows jobs use D:\a\wisp\wisp (14 characters).
#   Neither machine has the 43-character working directory that broke; what the
#   hosted side can NOT see is the width of the self-hosted one, which is why the
#   check below never uses the host it runs on as its denominator.
#
# THE CHECK - THE THRESHOLD IS THE HAT, NOT THE WALL
#   red iff   relative_length + WORST_PREFIX > FULL_PATH_BUDGET
#   which by construction (re-derived and asserted at every run) is
#           relative_length > HAT_RELATIVE_LIMIT
#   which for a file in .scratch/wisp/issues/ is
#           file-name length > HAT_NAME_LIMIT = 100        = README rule 9
#   A second hat is applied repo-wide, not only to tickets: a path whose LAST
#   COMPONENT exceeds HAT_NAME_LIMIT is an offender wherever it lives, because
#   the relative-length hat alone is looser than rule 9 for a shallow directory
#   (a 115-character name at the repository root is also a 115-character relative
#   path and would slip past it). Measured on the current tree the two hats select
#   the SAME 57 paths, so the second hat costs no extra roster entries.
#
#   HAT_NAME_LIMIT is the ONLY threshold knob in this file. The wall is
#   deliberately not the threshold: the wall sits inside the relative-length open
#   interval (206, 217] (re-derived from the run that died, see the 262-r1
#   evidence file section 1), and a gate placed there only rings on the day
#   checkout already fails. DEBT_CAP and the wall band are printed as HEADROOM
#   READINGS and can never turn this gate green or red.
#   Raising HAT_NAME_LIMIT to 180, 206 or 217 to silence a red run is the exact
#   edit ticket 262 lists under "forbidden actions"; it is also what the
#   counter-shape control in the 262-r1 evidence file demonstrates against.
#
# THE UNIT - CHARACTERS, AND WHY THAT IS A MEASURED CHOICE, NOT A GUESS
#   Rule 9 counts CHARACTERS and Windows counts UTF-16 code units, so for a CJK
#   name the character count is the physically right ruler. But awk's length()
#   counts BYTES unless the awk implementation is multibyte-aware AND the locale
#   is UTF-8 - and CI's ubuntu-latest image ships mawk as the default awk, which
#   is not. This repository has 28 tracked paths with non-ASCII bytes (the 262-r1
#   evidence file counts them), so the question is not rhetorical: in UTF-8 one
#   Chinese character is 3 bytes, and a byte-mode ruler reads a CJK name up to
#   three times too long.
#   The script therefore PROBES its own awk at run time and prints the answer:
#     unit=characters - ruling on characters. The byte view of the same
#         denominator is computed in the same run and the two offender sets are
#         compared, so "the two units agree here" is a reading of this tree,
#         re-verified on every run, never a cited assumption.
#     unit=bytes - ruling on bytes, which is stricter for CJK. If any offender in
#         that mode contains non-ASCII bytes it is a POSSIBLE FALSE RED this awk
#         cannot resolve, and the script refuses (exit 2) instead of reporting a
#         verdict it did not compute. It does not quietly pass, and it does not
#         quietly fail the tree: the fix it asks for is a multibyte awk.
#     anything else - refuse (exit 2).
#   On the current tree all 57 offenders are pure ASCII, so both units give the
#   same set; that equality is printed every run.
#
# HOW THE DENOMINATOR IS READ, AND THE RECORD-SPLIT GUARD
#   `git ls-files` is read in its raw (core.quotepath=false) form, newline
#   separated, rather than with NUL separators - because RS="\0" is not portable
#   across the awks this gate has to run under (gawk here, mawk in CI). That is
#   only safe if no tracked path contains a newline, and that is guarded rather
#   than assumed: with core.quotepath left at its default, git escapes such a
#   path onto ONE line, so `git ls-files | wc -l` and
#   `git -c core.quotepath=false ls-files | wc -l` differ by exactly the number of
#   embedded newlines. If the two counts disagree the script refuses (exit 2).
#
# SHAPE, AND WHY IT IS NOT A COPY OF scripts/d22scan.sh
#   scripts/d22scan.sh is `set -eu` with the seeded-violation positive control as
#   a PRECONDITION: if the control is red the real scan never runs. This gate
#   keeps that ordering and that abort (see --with-self-test below), but must not
#   use `set -e` inside the scan, because AC#2 requires the denominator to be
#   printed on a green run and the red path has to name EVERY offender, not stop
#   at the first. So the scan collects, prints, and decides its exit code once:
#     0 = green;  1 = red (budget violated, or the roster no longer matches the
#     tree);  2 = the gate refuses to answer (unit undecidable, empty
#     denominator, record split, or its own constants drifted apart).
#   No skippable step, no `|| true`, no continue-on-error (D22 run-away mode 6).
#
# THE ROSTER AND ITS GUARDS (ticket 262 AC#4)
#   Exemptions are DATA inside this file, not comments: each roster line assigns
#   one tracked path as a key and its reason as the value (see the awk program
#   below). The guards are
#   A  an over-budget tracked path that is not in the roster               -> red
#   B  a roster entry that is not an over-budget tracked path any more    -> red
#      (the drift pair: A catches a tree that grew past the roster, B catches a
#       roster that was left behind by a rename - on disk both look like green)
#   C  count(roster entries) != count(over-budget paths on the tree)      -> red
#      (deliberately redundant with A and B: this is the one-line arithmetic the
#       ticket asks for, and it also catches a duplicated roster key)
#   D  a roster entry whose reason is empty or whitespace                 -> red
#      Delete any one reason and this gate's behaviour changes. That is AC#4; a
#      reason that only a human reads is not a reason.
set -u

PROG=${0##*/}
ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$ROOT" ]; then
    printf '%s: not inside a git work tree - refusing to guess a denominator\n' "$PROG" >&2
    exit 2
fi

# ---- constants ---------------------------------------------------------------
# The only threshold knob. Source: .scratch/wisp/issues/README.md rule 9, a
# ticket file name of at most 100 characters including the NN- prefix and the
# .md suffix (README.md:57). Turning it is a contract-adjacent edit needing the
# same kind of named A585-style record that this script's CI step needed, not a
# fix for a red run.
HAT_NAME_LIMIT=100

# Worst-case checkout directory in this repository's CI, plus one separator.
# Source: ticket 262 section 1 item 1, re-measured by 262-r1 from the live runner
# log of run 37158259050: `Working directory is
# 'E:\work\base\actions-runner\_work\wisp\wisp'` = 43 characters, and joining a
# repository-relative path onto it adds the separator, giving 44.
WORST_PREFIX=44

# Length of the ticket directory the hat is converted through. Asserted against
# the real string at run time below, so a directory rename cannot silently move
# the threshold.
ISSUES_PREFIX_LEN=21

# Not thresholds. Printed as headroom readings only, see the header.
DEBT_CAP=180
WALL_LOW=206
WALL_HIGH=217

# ---- the roster --------------------------------------------------------------
# The exemption list lives inside the awk program below (between the
# ROSTER-AWK-BEGIN and ROSTER-AWK-END markers), one line per over-budget tracked
# path, with its reason as data. All 57 entries were first added to the
# repository BEFORE the rule-9 commit (newest add 2026-10-03 16:38:28 +0800, rule
# 9 landed 2026-10-04 08:55:30 +0800 in fd269de1), which is why they share the
# same "filed before the hat" sentence instead of carrying 57 separate stories;
# that is the category reason ticket 262 AC#4 allows, and guard D is what stops
# it from decaying into decoration.
#
# OLD NAMES ARE TRACEABLE (ticket 262 AC#7)
#   Eleven ticket paths were shortened with `git mv` on 2026-10-04 so that
#   `actions/checkout` could pass again: the two that actually broke the run
#   (commit fd269de1) and the nine that were next over the old debt line (commit
#   46079fcc). None of the eleven names below is a tracked path any more - they are
#   recorded here as the former spelling of the file that now sits beside them, so
#   a reader holding an old citation can find its object. Do not "fix" the old
#   names out of .scratch/wisp/probes/** or docs/evidence/s1/**: those files are
#   readings taken at a moment, and rewriting a reading is how a record stops
#   being evidence. Verify the pairs with:
#     git show --name-status -M --format= fd269de1
#     git show --name-status -M --format= 46079fcc
#   broke run 37158259050, former -> now:
#     .scratch/wisp/issues/256-resident-leg-cannot-read-those-risk-config-keys-because-approval-new-is-built-before-the-session-ledger-with-grants-nil-so-move-it-into-assembleruntime-deferred-until-confirming-is-measured.md
#       -> .scratch/wisp/issues/256-resident-gate-built-before-session-grants.md
#     .scratch/wisp/issues/257-on-a-clean-machine-all-seven-panel-settings-fields-are-unwritable-because-the-default-provider-registry-is-nil-while-the-write-side-requires-an-existing-row-and-no-roster-field-can-create-one.md
#       -> .scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md
#       -> (second hop, added 2026-10-07 by the orchestrator under A674 after acceptance
#           leg 262-v1 named this pair as the one stale pointer in this block: commit
#           6c96a425 later appended `-done`, so the path above is NOT a tracked path any
#           more. The pair above stays as written - it records what 46079fcc produced, and
#           rewriting a recorded reading is what this block exists to prevent. Current
#           tracked name, verified by `git ls-files .scratch/wisp/issues | grep /257-`:
#           .scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes-done.md)
#   nine more past the old debt line, former -> now (all in 46079fcc):
#     .scratch/wisp/issues/149-the-corrupt-leg-of-the-offset-field-has-zero-teeth-the-new-comment-claims-a-sufficient-condition-with-counterexamples-and-the-exited-branch-drops-the-summary-done.md
#       -> .scratch/wisp/issues/149-corrupt-offset-leg-zero-teeth-done.md
#     .scratch/wisp/issues/154-the-close-only-covers-the-loop-task-id-so-host-supplied-task-ids-still-have-no-owner-and-concurrent-tasks-have-zero-readings-reserved-with-a-trigger-gate-done.md
#       -> .scratch/wisp/issues/154-close-covers-loop-task-id-only-done.md
#     .scratch/wisp/issues/176-there-is-no-spawner-runasync-has-zero-production-call-sites-and-taskroster-record-has-zero-writers-so-164-ac4-has-no-object-and-175s-unstamped-path-waits-for-one-writer.md
#       -> .scratch/wisp/issues/176-no-spawner-zero-production-call-sites.md
#     .scratch/wisp/issues/177-stamping-external-content-and-letting-the-model-re-read-a-spilled-artifact-are-mutually-exclusive-today-because-the-stub-embeds-the-host-minted-path-and-r4-scans-path-arguments.md
#       -> .scratch/wisp/issues/177-stamp-vs-model-reread-spill-exclusive.md
#     .scratch/wisp/issues/183-the-host-minted-pointer-exemption-does-not-hold-on-the-real-cli-so-r4-refuses-the-models-own-reread-and-the-spilled-output-still-cannot-be-read-back-done.md
#       -> .scratch/wisp/issues/183-host-pointer-exemption-fails-real-cli-done.md
#     .scratch/wisp/issues/253-panel-inbound-has-three-ruler-holes-reachability-nail-only-recognizes-postmessage-config-methods-fall-outside-both-rosters-and-the-legacy-envelope-resolves-null-with-zero-audit.md
#       -> .scratch/wisp/issues/253-panel-inbound-three-ruler-holes.md
#     .scratch/wisp/issues/254-the-new-fifth-tier-winsec-has-no-caller-in-ci-so-the-fresh-guard-c-never-fires-unless-winsec-tests-sh-switches-to-scope-or-explicit-paths-are-proven-audited-done.md
#       -> .scratch/wisp/issues/254-fifth-tier-winsec-no-ci-caller-done.md
#     .scratch/wisp/issues/259-the-grant-binding-layer-is-an-identity-in-production-code-because-spend-is-fed-the-items-own-stored-digest-while-panelapi-and-panelitem-have-capability-side-ruler-holes.md
#       -> .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
#     .scratch/wisp/issues/260-the-configured-cancel-hotkey-never-registers-on-any-path-because-the-idle-pass-skips-the-slot-by-design-and-the-borrow-pass-hardcodes-bare-esc-while-a-lost-borrow-is-zero-symptom.md
#       -> .scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md

usage() {
    cat <<'EOF'
usage: sh scripts/check-path-length-budget.sh [--self-test | --with-self-test]

  (no option)        scan the tracked paths of the current repository
  --self-test        plant an over-budget tracked path in a throwaway repository,
                     assert this gate goes red and names it, then green again
  --with-self-test   control first, then the scan - what both call sites run

exit: 0 green, 1 red (budget violated or roster drifted), 2 the gate refused
EOF
}

MODE=scan
case "${1:-}" in
    --self-test)      MODE=selftest ;;
    --with-self-test) MODE=both ;;
    -h|--help)        usage; exit 0 ;;
    "")               MODE=scan ;;
    *) printf '%s: unknown option: %s\n' "$PROG" "$1" >&2; usage >&2; exit 2 ;;
esac
if [ $# -gt 1 ]; then
    printf '%s: too many arguments - this gate takes at most one mode option\n' "$PROG" >&2
    exit 2
fi

# ---- unit probe ---------------------------------------------------------------
SAMPLE=$(printf 'a\342\226\277b')   # "a" + U+25FF + "b" = 5 bytes, 3 characters
UNIT=$(awk -v s="$SAMPLE" 'BEGIN{print length(s)}' 2>/dev/null)
if [ "$UNIT" = "3" ]; then
    UNITNAME=characters
elif [ "$UNIT" = "5" ]; then
    UNITNAME=bytes
else
    UNITNAME=undetermined
fi

# Derived, never hand-set: see the consistency asserts below.
HAT_RELATIVE=$((HAT_NAME_LIMIT + ISSUES_PREFIX_LEN))
FULL_PATH_BUDGET=$((HAT_RELATIVE + WORST_PREFIX))

# ---- consistency asserts (the constants may not drift apart) -----------------
CONSISTENCY=$(awk -v issues="$ISSUES_PREFIX_LEN" -v namecap="$HAT_NAME_LIMIT" \
    -v hatrel="$HAT_RELATIVE" -v worst="$WORST_PREFIX" -v budget="$FULL_PATH_BUDGET" '
    BEGIN {
        if (length(".scratch/wisp/issues/") != issues)
            printf "drift: ISSUES_PREFIX_LEN=%s but .scratch/wisp/issues/ measures %d characters\n", issues, length(".scratch/wisp/issues/")
        if (hatrel != namecap + issues)
            printf "drift: HAT_RELATIVE=%s is not HAT_NAME_LIMIT + ISSUES_PREFIX_LEN\n", hatrel
        if (budget != hatrel + worst)
            printf "drift: FULL_PATH_BUDGET=%s is not HAT_RELATIVE + WORST_PREFIX\n", budget
    }')
if [ -n "$CONSISTENCY" ]; then
    printf '%s: REFUSE - its own constants no longer derive from each other:\n%s\n' "$PROG" "$CONSISTENCY" >&2
    exit 2
fi

# ---- denominator integrity guard --------------------------------------------
# Two reads of the same list: git escapes an embedded newline onto one line, so a
# disagreement between these two counts is an embedded newline, and a
# newline-separated read would then be inventing records.
Q_COUNT=$(git -C "$ROOT" -c core.quotepath=true ls-files | wc -l | tr -d ' ')
RAW_COUNT=$(git -C "$ROOT" -c core.quotepath=false ls-files | wc -l | tr -d ' ')
if [ "$Q_COUNT" != "$RAW_COUNT" ]; then
    printf '%s: REFUSE - the tracked list read two ways disagrees (%s quoted records vs %s raw records): a tracked path contains a newline and the newline-separated read is unsafe\n' \
        "$PROG" "$Q_COUNT" "$RAW_COUNT" >&2
    exit 2
fi

# ---- the scan ----------------------------------------------------------------
run_scan() {
    if [ "$UNITNAME" = "undetermined" ]; then
        printf '%s: REFUSE - this awk reports the 5-byte probe as "%s", which is neither 3 (characters) nor 5 (bytes); the unit is not knowable here\n' "$PROG" "$UNIT" >&2
        return 2
    fi

    printf '%s: repo=%s\n' "$PROG" "$ROOT"
    printf '%s: unit=%s (probe: length of the 5-byte sample "a U+25FF b" = %s; LANG=%s LC_ALL=%s)\n' \
        "$PROG" "$UNITNAME" "$UNIT" "${LANG:-unset}" "${LC_ALL:-unset}"
    printf '%s: hat: rule 9 name cap=%s + issues dir prefix=%s -> relative hat=%s ; worst checkout prefix=%s -> full-path budget=%s\n' \
        "$PROG" "$HAT_NAME_LIMIT" "$ISSUES_PREFIX_LEN" "$HAT_RELATIVE" "$WORST_PREFIX" "$FULL_PATH_BUDGET"
    printf '%s: check: relative_length + %s > %s, or a last name component > %s\n' \
        "$PROG" "$WORST_PREFIX" "$FULL_PATH_BUDGET" "$HAT_NAME_LIMIT"
    printf '%s: denominator read: %s tracked paths\n' "$PROG" "$RAW_COUNT"

    git -C "$ROOT" -c core.quotepath=false ls-files | awk \
        -v PROG="$PROG" -v hatrel="$HAT_RELATIVE" -v budget="$FULL_PATH_BUDGET" \
        -v worst="$WORST_PREFIX" -v namecap="$HAT_NAME_LIMIT" -v debtcap="$DEBT_CAP" \
        -v walllow="$WALL_LOW" -v wallhigh="$WALL_HIGH" -v unitname="$UNITNAME" '
function isort(a, n,   i, j, t) {
    for (i = 2; i <= n; i++) {
        t = a[i]; j = i - 1
        while (j >= 1 && a[j] > t) { a[j + 1] = a[j]; j-- }
        a[j + 1] = t
    }
}
BEGIN {
    tracked = 0; over = 0; rostered = 0; viol = 0; nst = 0; nnr = 0; nna = 0
    maxlen = 0; maxpath = ""; band = 0; debt = 0; wallb = 0
#ROSTER-AWK-BEGIN
    R[".scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/158-the-gate-stands-outside-the-shape-that-is-live-today-openscope-closescope-has-no-clause-and-the-bridge-scope-open-is-guarded-only-in-another-package.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/148-clonedecision-only-guarded-one-door-the-sibling-replay-outlet-still-hands-out-a-shallow-value-copy-and-what-the-panel-should-receive-is-unwritten.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/179-loop-deciderisk-treats-a-declared-l0-as-unclassified-so-every-readonly-builtin-is-refused-on-the-real-cli-and-no-criterion-holds-that-sentence.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/152-the-real-subject-death-shape-has-never-been-measured-only-faked-to-be-reachable-and-the-offset-fallback-is-a-leg-that-changes-no-reading-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/151-bridge-opentask-sits-on-the-dispatch-hot-path-while-closetask-has-zero-callers-repo-wide-so-the-per-task-taint-scope-is-only-ever-opened-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/184-a-decision-column-of-five-frozen-values-cannot-tell-an-auto-allowed-l0-from-a-human-approval-and-grant-id-has-zero-production-writers-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/171-three-ruler-holes-ticket-161-cannot-carry-aggregate-rs-selftest-not-in-ci-trailing-comments-counted-as-calls-and-count-blind-pair-legs.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/258-the-schema-calls-hotkey-hot-tier-but-the-resident-leg-builds-the-ball-from-default-hotkeys-and-the-only-reloader-caller-is-balldebug.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/155-ticket-153-leaves-three-cells-unjudged-because-my-attack-point-grid-and-the-ticket-ac-grid-are-not-the-same-coordinate-system-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/147-collectreports-named-red-says-document-stopped-at-offset-n-but-the-unwritten-branch-always-carries-0-and-no-assertion-pins-it-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/175-task-output-returns-external-content-with-no-c25-provenance-mark-because-marking-is-name-gated-and-the-name-is-not-on-the-roster.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/182-the-task-monitor-rail-is-outside-ticket-145s-fourteen-row-table-so-nobody-has-counted-which-of-its-stacks-have-a-go-side-source.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/197-the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/181-nothing-in-go-ever-reads-git-so-the-panel-cannot-show-the-branch-and-the-two-mutating-actions-would-need-c17-whitelist-methods.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/195-config-set-needs-a-per-section-write-foundation-before-the-panel-can-touch-it-and-the-only-writer-in-the-repo-is-permmode.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/150-four-deferred-markers-have-no-row-in-the-spec-12-registry-so-the-1-to-1-rule-is-violated-today-and-no-instrument-scans-it.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/137-winsecs-12-posix-rejection-legs-go-green-in-the-symlink-shape-for-the-hosts-own-varlink-not-their-own-planted-link-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-because-the-fs-read-result-becomes-an-undeclared-c25-mark.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/196-two-task-state-vocabularies-already-coexist-and-the-schema-has-no-check-so-the-d43-names-are-not-the-ones-in-use.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/251-winsec-pin-is-a-pin-nobody-compares-the-winsec-scope-runs-through-explicit-paths-so-guard-c-never-sees-it-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/235-the-pool-nail-comment-reason-was-falsified-by-ticket-222-and-only-one-of-three-cases-sees-the-real-bridge-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/156-q-57-landed-as-option-b-exited-asks-the-os-instead-of-our-own-record-and-the-152-leftover-comment-actions-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/178-pair-census-ruler-g6-negative-leg-cannot-express-closing-a-scope-opened-elsewhere-and-rings-every-honest-test.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/dispatches/2026-09-28-114x-accept-179-v1-is-any-of-the-three-new-criteria-tautological-and-who-owns-the-decision-column.md"] = "hist: dispatch note filed before README rule 9 (100-char hat, added 2026-10-04); rule 9 governs ticket names only"
    R[".scratch/wisp/issues/191-the-panel-embedded-terminal-puts-an-execution-surface-on-the-panel-channel-so-it-needs-a-named-ruling-first.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/127-the-resident-gui-leg-installs-its-log-listener-with-zero-nails-deleting-20-lines-keeps-146-tests-green-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/144-wisp-slo-collectreport-reads-a-still-being-written-file-as-a-corrupt-one-one-empty-read-kills-the-run-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/173-the-audit-line-prints-in-allowlist-scope-true-when-nothing-was-judged-blank-path-both-meanings-one-field.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/dispatches/2026-09-28-135x-impl-183-r2-grow-the-two-missing-teeth-and-re-take-the-gates-after-your-own-last-commit.md"] = "hist: dispatch note filed before README rule 9 (100-char hat, added 2026-10-04); rule 9 governs ticket names only"
    R[".scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/140-job-level-wisp-env-test-makes-the-ask-the-os-hardening-legs-never-get-consulted-on-the-windows-runner.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/128-resolvedatadir-falls-back-to-current-directory-when-appdata-is-missing-moving-config-dpapi-and-memory.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/243-a-closed-ticket-number-still-owns-work-in-25-comments-which-is-the-shape-that-fooled-three-legs-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/129-sametree-also-drops-the-absoluteness-segment-so-a-drive-relative-spelling-passes-the-seam-guard-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/169-a-single-decorative-glyph-in-frontend-makes-the-repo-wide-gate-red-and-kills-every-ci-step-after-it.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/135-the-only-teeth-in-ticket-128s-judgement-have-no-self-proof-leg-two-mutations-stacked-go-16-16-green.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/dispatches/2026-09-28-142x-readonly-185-c1-who-should-declare-the-read-back-path-and-what-does-each-route-cost.md"] = "hist: dispatch note filed before README rule 9 (100-char hat, added 2026-10-04); rule 9 governs ticket names only"
    R[".scratch/wisp/issues/246-the-process-that-runs-tasks-is-not-the-process-holding-the-ball-so-no-veto-channel-has-an-executor.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/193-lightweight-must-no-longer-constrain-frontend-animation-and-visuals-make-the-per-panel-gates-named.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/134-slofull-autostarts-on-every-push-on-this-laptop-but-making-it-manual-only-would-kill-the-gate-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/172-no-ruler-links-the-frozen-tier-in-plan-md-to-the-tier-the-gate-actually-gives-a-zero-tests-go-red.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/139-context-compression-leaves-no-trace-zero-log-calls-in-compress-go-zero-readers-of-res-compression.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/126-winsec-tree-attribution-strips-the-volume-so-a-seal-on-one-drive-is-reported-against-another-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/143-wisp-panel-assets-l2-cannot-produce-an-r4-card-the-taint-detector-is-not-wired-on-that-path-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/192-the-browser-and-artifact-preview-stack-is-unclaimed-and-sits-on-the-undefined-web-search-ruling.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); still open, renaming it is outside this gate"
    R[".scratch/wisp/issues/142-the-two-git-spawns-are-not-atomic-so-a-tracked-directory-can-be-skipped-during-the-window-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
    R[".scratch/wisp/issues/141-ban8-emoji-scan-range-is-narrower-than-plan-md-states-41-go-files-carry-u2190-u25ff-today-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"
#ROSTER-AWK-END
}
{
    tracked++
    len = length($0)
    n = split($0, comp, "/")
    namelen = length(comp[n])
    if (len > maxlen) { maxlen = len; maxpath = $0 }
    if (len > hatrel && len <= debtcap) band++
    if (len > debtcap) debt++
    if (len > walllow && len <= wallhigh) wallb++
    if (len > hatrel || namelen > namecap) {
        over++
        LENV[$0] = len
        NAMEV[$0] = namelen
        if ($0 in R) {
            rostered++
            if (R[$0] ~ /^[ \t]*$/) { nnr++; NOD[nnr] = $0 }
        } else {
            viol++; VD[viol] = $0
        }
        if ($0 ~ /[^ -~]/) { nna++; AMD[nna] = $0 }
    }
}
END {
    sizeR = 0
    for (p in R) { sizeR++; if (!(p in LENV)) { nst++; SD[nst] = p } }

    printf "%s: denominator: tracked paths=%d  over-budget=%d  covered by roster=%d  not in roster=%d\n", \
        PROG, tracked, over, rostered, viol
    printf "%s: longest=%d chars relative (%s)\n", PROG, maxlen, maxpath
    printf "%s: worst full path on the self-hosted runner=%d chars (hat budget %d, wall open interval (%d,%d])\n", \
        PROG, maxlen + worst, budget, walllow, wallhigh
    printf "%s: bands: over the hat=%d  of which in the 122..%d middle=%d  past the old debt line(%d)=%d  in the wall interval=%d  roster entries=%d\n", \
        PROG, over, debtcap, band, debtcap, debt, wallb, sizeR

    if (tracked == 0) {
        printf "%s: REFUSE - zero tracked paths. A gate whose denominator is empty cannot report green\n", PROG
        exit 2
    }

    if (viol > 0) {
        isort(VD, viol)
        for (i = 1; i <= viol; i++)
            printf "%s: RED - over budget and NOT in the roster: %s  (relative %d chars, name %d chars, full path on the self-hosted runner %d chars, budget %d)\n", \
                PROG, VD[i], LENV[VD[i]], NAMEV[VD[i]], LENV[VD[i]] + worst, budget
    }
    if (nst > 0) {
        isort(SD, nst)
        for (i = 1; i <= nst; i++)
            printf "%s: RED - roster entry that is not an over-budget tracked path any more (the roster has drifted away from the tree): %s\n", \
                PROG, SD[i]
    }
    if (nnr > 0) {
        isort(NOD, nnr)
        for (i = 1; i <= nnr; i++)
            printf "%s: RED - roster entry carries an EMPTY reason; the reason is data, not a comment: %s\n", \
                PROG, NOD[i]
    }
    if (sizeR != over)
        printf "%s: RED - count guard: the roster holds %d entries, the tree has %d over-budget tracked paths\n", \
            PROG, sizeR, over

    if (nna > 0) {
        isort(AMD, nna)
        for (i = 1; i <= nna; i++)
            printf "%s: NOTE - this offender contains non-ASCII bytes, so its byte length and its character length differ: %s\n", \
                PROG, AMD[i]
        if (unitname == "bytes") {
            printf "%s: REFUSE - unit=bytes and %d offender(s) above carry non-ASCII bytes. Byte length over-reads a CJK name by up to 3 times and this awk cannot recover the character count. The fix is a multibyte awk (gawk) under a UTF-8 locale, not a higher hat.\n", \
                PROG, nna
            exit 2
        }
    }

    if (viol > 0 || nst > 0 || nnr > 0 || sizeR != over) {
        printf "%s: VERDICT RED\n", PROG
        exit 1
    }
    printf "%s: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree\n", PROG
    exit 0
}'
    return $?
}

# Offender lists under each unit, for the cross-check printed with the verdict.
# Same newline-separated read as the scan (the record-split guard above is what
# makes that legal), run once per unit so the comparison is two readings.
byte_view_offenders() {
    git -C "$ROOT" -c core.quotepath=false ls-files | LC_ALL=C awk \
        -v hatrel="$HAT_RELATIVE" -v namecap="$HAT_NAME_LIMIT" '
        { n = split($0, c, "/"); if (length($0) > hatrel || length(c[n]) > namecap) print length($0), $0 }' | LC_ALL=C sort -k2
}
char_view_offenders() {
    git -C "$ROOT" -c core.quotepath=false ls-files | awk \
        -v hatrel="$HAT_RELATIVE" -v namecap="$HAT_NAME_LIMIT" '
        { n = split($0, c, "/"); if (length($0) > hatrel || length(c[n]) > namecap) print length($0), $0 }' | LC_ALL=C sort -k2
}

# ---- positive control --------------------------------------------------------
run_selftest() {
    BENCH=$(mktemp -d 2>/dev/null)
    if [ -z "$BENCH" ] || [ ! -d "$BENCH" ]; then
        printf '%s: REFUSE - the positive control cannot build a throwaway repository (mktemp returned "%s")\n' "$PROG" "$BENCH" >&2
        return 2
    fi
    printf '%s: positive control bench=%s (kept on disk; this project never deletes temp artifacts)\n' "$PROG" "$BENCH"

    PAD=$(printf '%40s' '' | tr ' ' 'x')
    SHORT_OK="ok-short-name.md"
    PLANTED=".scratch/wisp/issues/seeded-$PAD-$PAD-over-the-hat.md"

    if ! git init -q "$BENCH" 2>/dev/null; then
        printf '%s: REFUSE - git init failed in %s\n' "$PROG" "$BENCH" >&2
        return 2
    fi
    mkdir -p "$BENCH/.scratch/wisp/issues"
    printf 'seeded by scripts/check-path-length-budget.sh --self-test\n' > "$BENCH/$SHORT_OK"
    printf 'seeded by scripts/check-path-length-budget.sh --self-test\n' > "$BENCH/$PLANTED"

    # A copy of this gate with the roster block deleted. The bench repository
    # holds none of the real rostered paths, so an honest run against it starts
    # from an empty exemption list; deleting only the block (not rewriting the
    # check) is what keeps this control pointed at the shipped logic. If the
    # markers were ever removed, this sed would produce an unchanged copy, the
    # bench baseline would then be red on 57 stale entries and control 1/3 would
    # fail out loud instead of passing vacuously.
    ROSTER_MARK_BEGIN='#ROSTER-AWK-BEGIN'
    ROSTER_MARK_END='#ROSTER-AWK-END'
    sed "/^${ROSTER_MARK_BEGIN}\$/,/^${ROSTER_MARK_END}\$/d" "$0" > "$BENCH/gate.sh" 2>/dev/null
    if [ ! -s "$BENCH/gate.sh" ] || grep -qE '^    R\["' "$BENCH/gate.sh"; then
        printf '%s: REFUSE - the roster-less copy of itself is missing or still carries roster lines\n' "$PROG" >&2
        return 2
    fi

    git -C "$BENCH" add -- "$SHORT_OK"
    rc=0
    ( cd "$BENCH" && sh gate.sh ) > "$BENCH/run1.txt" 2>&1 || rc=$?
    if [ "$rc" -ne 0 ]; then
        printf '%s: CONTROL FAILED - the bench baseline (one short tracked path, empty roster) should be GREEN, got rc=%s. The gate may be stuck red\n' "$PROG" "$rc" >&2
        cat "$BENCH/run1.txt" >&2
        return 1
    fi
    printf '%s: control 1/3 ok - bench baseline green, so this gate is not stuck red\n' "$PROG"

    git -C "$BENCH" add -- "$PLANTED"
    rc=0
    ( cd "$BENCH" && sh gate.sh ) > "$BENCH/run2.txt" 2>&1 || rc=$?
    if [ "$rc" -ne 1 ]; then
        printf '%s: CONTROL FAILED - planting a tracked path %d characters long should be RED (rc=1), got rc=%s. The gate may be stuck green\n' "$PROG" ${#PLANTED} "$rc" >&2
        cat "$BENCH/run2.txt" >&2
        return 1
    fi
    if ! grep -qF -- "$PLANTED" "$BENCH/run2.txt"; then
        printf '%s: CONTROL FAILED - the gate went red but never NAMED the planted path %s\n' "$PROG" "$PLANTED" >&2
        cat "$BENCH/run2.txt" >&2
        return 1
    fi
    printf '%s: control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:\n' "$PROG"
    sed -n 's/.*\(RED - over budget.*\)/\1/p' "$BENCH/run2.txt"

    git -C "$BENCH" rm --cached --quiet -- "$PLANTED" 2>/dev/null
    rc=0
    ( cd "$BENCH" && sh gate.sh ) > "$BENCH/run3.txt" 2>&1 || rc=$?
    if [ "$rc" -ne 0 ]; then
        printf '%s: CONTROL FAILED - after un-planting, the gate should be GREEN again, got rc=%s\n' "$PROG" "$rc" >&2
        cat "$BENCH/run3.txt" >&2
        return 1
    fi
    printf '%s: control 3/3 ok - green again once the planted path is gone\n' "$PROG"
    printf '%s: positive control PASSED\n' "$PROG"
    return 0
}

rc=0
if [ "$MODE" = selftest ] || [ "$MODE" = both ]; then
    if ! run_selftest; then
        printf '%s: positive control RED - the scan below is not run, because a gate that cannot go red does not get to print green\n' "$PROG" >&2
        exit 1
    fi
    if [ "$MODE" = selftest ]; then
        exit 0
    fi
fi

run_scan
SCAN_RC=$?

if [ "$UNITNAME" = "characters" ]; then
    CTL=$(mktemp -d 2>/dev/null)
    byte_view_offenders > "$CTL/byte.txt"
    char_view_offenders > "$CTL/char.txt"
    if cmp -s "$CTL/byte.txt" "$CTL/char.txt"; then
        printf '%s: unit cross-check: the byte view and the character view select the SAME %d paths on this tree, recomputed just now in both units\n' \
            "$PROG" "$(grep -c . "$CTL/char.txt" || true)"
    else
        printf '%s: unit cross-check: the two units DISAGREE on this tree.\n' "$PROG"
        printf '%s: over-budget only in BYTES (a CJK name a byte-mode awk would over-read):\n' "$PROG"
        LC_ALL=C comm -23 "$CTL/byte.txt" "$CTL/char.txt"
        printf '%s: over-budget only in CHARACTERS:\n' "$PROG"
        LC_ALL=C comm -13 "$CTL/byte.txt" "$CTL/char.txt"
        printf '%s: the ruling stays on characters (rule 9 counts characters); the lines above are a reading of this tree, not a failure of this gate\n' "$PROG"
    fi
fi

exit $SCAN_RC
