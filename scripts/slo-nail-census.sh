#!/bin/sh
# scripts/slo-nail-census.sh - the OUTER tooth over the two self-locking nails of
# the D32 gate (ticket 266, shape 1-A).
#
# WHY THIS FILE EXISTS AT ALL
#   scripts/slo-check.ps1 carries two nails that inspect THAT SAME FILE:
#     Test-ScriptShapeNail     the word half - the file must contain no direct
#                              read of the automatic exit-code variable (ticket 263)
#     Test-ExitCodeInstrument  the capability half - a real program is started,
#                              waited for, and its real exit code read (ticket 263)
#   A self-locking nail is bypassed by deleting the nail, and that is not a quirk,
#   it is the definition of the instrument. Ticket 263's verdict measured exactly
#   that mutant (mutant name delete-nails): with both call lines removed the gate
#   exited 0, printed an all-green screen, and not one self-check line appeared.
#   Nothing elsewhere in the repository knew those nails were supposed to be
#   there, so the repository had two teeth and no tooth over the teeth. The answer
#   is not a third self-check inside slo-check.ps1 - that would be the same shape
#   one more time - it is a declaration that lives in ANOTHER file and says what
#   the gate must contain:
#
#       the D32 gate has exactly 2 self-locking nails
#       Test-ScriptShapeNail      carrying "shape nail ok"
#       Test-ExitCodeInstrument   carrying "instrument self-check ok"
#
#   The two assignments below ARE that declaration. They live in this file; the
#   inspected file never mentions them. Deleting a call from slo-check.ps1 now
#   leaves a declaration behind that says it should exist.
#
# WHAT ONE RUN CHECKS (and what each finding means)
#   1. declaration drift   declared_nail_count must equal the number of rows in
#                          the two tables below. Editing one without the other is
#                          a red, not a silent redefinition of the gate.
#   2. definition          exactly one column-0 "function <Name> {" per nail.
#   3. ok sentence         the nail's own success sentence appears inside ITS OWN
#                          body (from its definition to the next top-level
#                          function). A sentence that survives only in a comment,
#                          or drifted into a neighbour's body, is a red.
#   4. call site           exactly one column-0, uncommented, standalone call
#                          line. Commented-out and block-commented copies do not
#                          count (blocks and line comments are stripped, and a
#                          line comment can never reach column 0 anyway); a second
#                          copy is a red because two runs of a self-check is not
#                          what this file declares.
#   5. ordering            each call line is after its own definition, and both
#                          call lines precede the first sampling call - which is
#                          the property the gate's own comment claims ("Both
#                          nails run before anything is measured", ticket 263).
#   Every finding prints the file and the line numbers actually read, so a red
#   names its evidence instead of a rule.
#
# WHAT THIS FILE CANNOT SEE - written down here instead of discovered later
#   (a) It pins SHAPE, not execution. A call line that exists but is unreachable
#       (for example moved inside a function nobody calls) still counts. The
#       runtime half - "those two ok lines were actually printed by that run" -
#       needs an outer step that reads the gate's stdout or its report; that is
#       shape 2 of ticket 266 and it does not exist today.
#   (b) The whole-file bypass is open and must not be sold as closed: one commit
#       that deletes this file AND the CI step that calls it leaves nothing that
#       rings. The only other outer tooth, probe P3 of scripts/slo-freshness.sh,
#       ages on the newest uploaded slo-full-report ARTIFACT; it is blind to the
#       presence of a nail and fires only after days without a conclusion. This
#       is therefore a tooth against a one-line edit, not a door that closes.
#   (c) The reader is textual. A here-string body carrying a column-0 call line
#       would be counted (none exists in the target today); a call assembled at
#       run time is missed; and the quote tracking below is single-line except
#       for the two constructs that really span lines in this target (block
#       comments and here-strings), which get their own state.
#
# EXIT CODES AND READINGS
#   0  every declaration above holds. The PASS line says so in words.
#   1  a declaration is violated. The finding carries one of
#      slo-nail-census-declaration-drift /
#      slo-nail-census-nail-definition-missing /
#      slo-nail-census-nail-definition-duplicated /
#      slo-nail-census-nail-ok-sentence-missing /
#      slo-nail-census-nail-call-missing /
#      slo-nail-census-nail-call-duplicated /
#      slo-nail-census-ordering
#   2  CANNOT DECIDE - the target is missing, empty, unreadable, or the sampling
#      anchor this file orders against is gone. No verdict is printed, the token
#      slo-nail-census-cannot-decide names the reason, and the exit code says
#      "not a pass" (the honesty tier of ticket 266; same precedent as
#      Fail-InstrumentBroken in ticket 263 and the validity precheck of 134).
#
# USAGE
#   sh scripts/slo-nail-census.sh                   read the real target
#   sh scripts/slo-nail-census.sh --with-self-test  positive control, then read
#   sh scripts/slo-nail-census.sh --self-test       positive control only
#
# The positive control never touches the real tree: it builds a throwaway bench,
# copies THIS script and the target into it, and mutates only those copies - so
# "delete the nails and watch the outer tooth ring" is a reading that can be
# taken on any machine without editing the gate. The target path is derived from
# this script's own location (the scripts/d22scan.sh rule: no cwd dependence and
# no second copy of the path), which is also what lets the bench work with no
# test seam at all: no environment variable can point this reader at another
# file, because that would be a one-line way to make it green.
#
# No step here is skippable, no exit code is swallowed, no threshold, golden or
# Set-StrictMode is touched, and this file adds no rule to scripts/** beyond the
# one tooth it is (ticket 266 scope, AC#3).
set -eu

PROG=slo-nail-census
MODE=run

# ---------------------------------------------------------------------------
# THE DECLARATION - the part that lives here and not in the inspected file
# ---------------------------------------------------------------------------
declared_nail_count=2
nail_names='Test-ScriptShapeNail|Test-ExitCodeInstrument'
nail_ok='shape nail ok|instrument self-check ok'

# The first sampling call of the gate: what the nails must precede. Anchored on
# the one program the gate actually samples, never on a line number.
sampling_anchor='-FilePath $WispExe'

# ---------------------------------------------------------------------------
# where the inspected file is - derived from this script's own location
# ---------------------------------------------------------------------------
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
target=$here/slo-check.ps1
root=$here
if [ -f "$here/../go.mod" ]; then
    root=$(CDPATH= cd -- "$here/.." && pwd)
fi
case $target in
    "$root"/*) rel=${target#"$root/"} ;;
    *)         rel=$target ;;
esac

usage() {
    cat <<'EOF'
scripts/slo-nail-census.sh - the outer tooth over the D32 gate's self-locking nails
(ticket 266). It reads scripts/slo-check.ps1 and checks that the two nails this
file declares are still defined, still carry their own ok sentence, still have
exactly one uncommented column-0 call line, and are still called before the first
SLO sample.

  sh scripts/slo-nail-census.sh                   read the real target
  sh scripts/slo-nail-census.sh --with-self-test  positive control first, then read
  sh scripts/slo-nail-census.sh --self-test       positive control only

Exit 0 = every declaration holds; 1 = a declaration is violated; 2 = CANNOT
DECIDE (missing/empty/unreadable target, or the sampling anchor is gone) - which
is not a pass and prints no verdict.
EOF
}

while [ $# -gt 0 ]; do
    case $1 in
        --with-self-test|--self-test)
            if [ "$MODE" = "both" ]; then
                printf '%s: --with-self-test and --self-test are exclusive, refusing to guess\n' "$PROG" >&2
                exit 2
            fi
            if [ "$1" = "--with-self-test" ]; then MODE=both; else MODE=selftest; fi
            ;;
        -h|--help) usage; exit 0 ;;
        *)
            printf '%s: unknown argument %s (try --help)\n' "$PROG" "$1" >&2
            exit 2
            ;;
    esac
    shift
done

cannot_decide() {  # $1 = the named reason
    printf '%s: CANNOT DECIDE slo-nail-census-cannot-decide: %s\n' "$PROG" "$1" >&2
    printf '%s: CANNOT DECIDE: no verdict about the D32 gate nails was produced, and this is NOT a pass - no reading may be taken from this run (ticket 266 honesty tier).\n' "$PROG" >&2
    exit 2
}

findings=0
fail_finding() {  # $1 = the full finding, token first
    findings=$((findings + 1))
    printf '%s: FAIL %s\n' "$PROG" "$1"
}

# ---------------------------------------------------------------------------
# the reader: one line at a time, comments stripped, quotes respected
# ---------------------------------------------------------------------------
# Emits one row per declared nail plus the anchor line actually found:
#   NAIL <name> <defcount> <defline> <okcount> <okline> <callcount> <callline> <oktext>
#   ANCHOR <line>            (0 when the anchor does not occur at all)
analyze() {  # $1 = the file to read
    tr -d '\015' < "$1" | awk \
        -v names="$nail_names" \
        -v oks="$nail_ok" \
        -v anchor="$sampling_anchor" \
        -v SQ="'" '
        # Strip the two PowerShell comment forms and keep string contents, so a
        # sentence that survives only inside a comment cannot count as a nail.
        function strip(line,   out, q, n, i, ch) {
            out = ""; q = ""; n = length(line); i = 1
            while (i <= n) {
                ch = substr(line, i, 1)
                if (inblock) {
                    if (ch == "#" && substr(line, i + 1, 1) == ">") { inblock = 0; i += 2; continue }
                    i++; continue
                }
                if (q != "") {
                    if (ch == q) {
                        if (substr(line, i + 1, 1) == q) { out = out ch ch; i += 2; continue }
                        q = ""; out = out ch; i++; continue
                    }
                    if (q == "\"" && ch == "`") {
                        out = out ch
                        if (i < n) { out = out substr(line, i + 1, 1) }
                        i += 2; continue
                    }
                    out = out ch; i++; continue
                }
                if (ch == SQ || ch == "\"") { q = ch; out = out ch; i++; continue }
                if (ch == "<" && substr(line, i + 1, 1) == "#") { inblock = 1; i += 2; continue }
                if (ch == "#") { break }
                out = out ch; i++
            }
            return out
        }
        BEGIN {
            nn = split(names, nm, "|")
            nok = split(oks, ok, "|")
            if (nn != nok) { print "ROWCOUNT " nn " " nok; exit 3 }
            inblock = 0; hs = ""; curfn = ""; anchorline = 0
        }
        {
            line = $0
            if (hs != "") {
                # Inside a here-string: no comments, and the closer is the only
                # thing that ends it.
                code = line
                if (line ~ ("^[[:space:]]*" hs "@[[:space:]]*$")) { hs = "" }
            } else {
                code = strip(line)
                tail = code
                sub(/[[:space:]]+$/, "", tail)
                opener = substr(tail, length(tail) - 1, 2)
                if (opener == "@" SQ) { hs = SQ }
                else if (opener == "@\"") { hs = "\"" }
            }
            # The top-level function this line belongs to, so an ok sentence must
            # sit in ITS body and not somewhere else in the file.
            if (code ~ /^function[[:space:]]+[A-Za-z0-9_-]+/) {
                f = code
                sub(/^function[[:space:]]+/, "", f)
                sub(/[^A-Za-z0-9_-].*$/, "", f)
                curfn = f
            }
            for (i = 1; i <= nn; i++) {
                if (code ~ ("^function[[:space:]]+" nm[i] "[[:space:]]*\\{")) {
                    defcount[i]++
                    if (defline[i] == 0) { defline[i] = NR }
                }
                if (curfn == nm[i] && NR > defline[i] && index(code, ok[i]) > 0) {
                    okcount[i]++
                    if (okline[i] == 0) { okline[i] = NR }
                }
                if (code ~ ("^" nm[i] "[[:space:]]*$")) {
                    callcount[i]++
                    if (callline[i] == 0) { callline[i] = NR }
                }
            }
            if (anchorline == 0 && index(code, anchor) > 0) { anchorline = NR }
        }
        END {
            if (nn != nok) { exit 3 }
            for (i = 1; i <= nn; i++) {
                printf "NAIL %s %d %d %d %d %d %d %s\n", nm[i], \
                    defcount[i] + 0, defline[i] + 0, okcount[i] + 0, okline[i] + 0, \
                    callcount[i] + 0, callline[i] + 0, ok[i]
            }
            printf "ANCHOR %d\n", anchorline + 0
        }'
}

# ---------------------------------------------------------------------------
# the census: the real read of the real target
# ---------------------------------------------------------------------------
census() {
    findings=0

    if [ ! -e "$target" ]; then
        cannot_decide "the inspected file $rel does not exist - the D32 gate it belongs to is not on disk, so nothing can be said about its nails"
    fi
    if [ ! -f "$target" ]; then
        cannot_decide "the inspected path $rel is not a regular file"
    fi
    if [ ! -r "$target" ]; then
        cannot_decide "the inspected file $rel is not readable by this user"
    fi
    if [ ! -s "$target" ]; then
        cannot_decide "the inspected file $rel is empty - an empty gate cannot carry a nail, and an empty read is not a silent pass"
    fi

    # 1. declaration drift. Both tables and the constant must say the same thing.
    rows_names=$(printf '%s\n' "$nail_names" | awk -F'|' '{ print NF }')
    rows_ok=$(printf '%s\n' "$nail_ok" | awk -F'|' '{ print NF }')
    if [ "$rows_names" -ne "$declared_nail_count" ] || [ "$rows_ok" -ne "$declared_nail_count" ]; then
        printf '%s: FAIL slo-nail-census-declaration-drift: declared_nail_count=%s but the name table holds %s row(s) and the ok table %s row(s). The declaration in this file no longer says one thing, so it cannot count anything.\n' \
            "$PROG" "$declared_nail_count" "$rows_names" "$rows_ok"
        printf '%s: no nail was examined in this run; the numbers above are not a reading of the gate.\n' "$PROG"
        return 1
    fi

    rowsfile=$(mktemp 2>/dev/null) || cannot_decide "cannot create a temporary file for the analysis rows"
    if ! analyze "$target" > "$rowsfile"; then
        cannot_decide "the reader (tr | awk) failed on $rel - the file could not be turned into a reading"
    fi
    if [ ! -s "$rowsfile" ]; then
        cannot_decide "the reader produced no rows at all for $rel"
    fi

    anchor=0
    while read -r kind a b c d e f g h; do
        if [ "$kind" = "ANCHOR" ]; then anchor=$a; fi
    done < "$rowsfile"

    if [ "$anchor" -eq 0 ]; then
        cannot_decide "the sampling anchor (a line carrying [$sampling_anchor]) is gone from $rel. This file orders the nails against the gate's first sample and cannot do that any more. If the gate legitimately changed its sampling entry point, update the declaration in this file - a two-file diff somebody reviews - instead of widening this reader silently."
    fi

    printf '%s: file=%s (absolute %s)\n' "$PROG" "$rel" "$target"
    printf '%s: declared nails=%s - the declaration lives in this file, not in %s\n' \
        "$PROG" "$declared_nail_count" "$rel"
    printf '%s: sampling anchor: first line carrying [%s] is %s:%s\n' \
        "$PROG" "$sampling_anchor" "$rel" "$anchor"

    i=0
    while read -r kind name defc defl okc okl callc calll oktext; do
        if [ "$kind" != "NAIL" ]; then continue; fi
        i=$((i + 1))
        printf '%s: nail %s/%s %s - definition %s:%s, ok sentence %s:%s, call %s:%s\n' \
            "$PROG" "$i" "$declared_nail_count" "$name" \
            "$rel" "$defl" "$rel" "$okl" "$rel" "$calll"

        # 2. exactly one definition
        if [ "$defc" -eq 0 ]; then
            fail_finding "slo-nail-census-nail-definition-missing: the declared nail '$name' has no column-0 'function $name {' in $rel. This file declares that the D32 gate carries it (ticket 266 AC#1); a gate without it can go all-green while measuring nothing (ticket 263 mutant delete-nails)."
        elif [ "$defc" -gt 1 ]; then
            fail_finding "slo-nail-census-nail-definition-duplicated: the declared nail '$name' is defined $defc times in $rel (first at $rel:$defl); this file declares exactly one definition."
        fi

        # 3. the nail's own ok sentence, inside its own body
        if [ "$okc" -eq 0 ]; then
            fail_finding "slo-nail-census-nail-ok-sentence-missing: the declared nail '$name' carries no [$oktext] sentence inside its own body in $rel. A nail whose success line was deleted still passes itself and stops being visible on the gate's screen, which is the silent half of ticket 263."
        fi

        # 4. exactly one uncommented column-0 call line
        if [ "$callc" -eq 0 ]; then
            fail_finding "slo-nail-census-nail-call-missing: the declared nail '$name' has no column-0, uncommented, standalone call line in $rel (its definition sits at $rel:$defl). Nothing runs it - this is the ticket 263 delete-nails mutant, and the reason this tooth exists (ticket 266 AC#1/AC#2)."
        elif [ "$callc" -gt 1 ]; then
            fail_finding "slo-nail-census-nail-call-duplicated: the declared nail '$name' has $callc column-0 call lines in $rel (first at $rel:$calll); this file declares exactly one."
        fi

        # 5. ordering: after its own definition, before the first sample
        if [ "$callc" -gt 0 ] && [ "$defl" -gt 0 ] && [ "$calll" -lt "$defl" ]; then
            fail_finding "slo-nail-census-ordering: the call to '$name' sits at $rel:$calll, BEFORE its own definition at $rel:$defl. PowerShell resolves functions in script order, so that line is not a call."
        fi
        if [ "$callc" -gt 0 ] && [ "$calll" -ge "$anchor" ]; then
            fail_finding "slo-nail-census-ordering: the call to '$name' sits at $rel:$calll, not before the first sampling call at $rel:$anchor. The gate claims both nails run before anything is measured (ticket 263); below the sampler, even a printed break arrives after numbers were already claimed."
        fi
    done < "$rowsfile"

    if [ "$i" -ne "$declared_nail_count" ]; then
        cannot_decide "the reader returned $i nail row(s) for a declaration of $declared_nail_count - the rows do not cover the declaration, so no verdict can be built from them"
    fi

    if [ "$findings" -gt 0 ]; then
        printf '%s: FAIL %s declaration violation(s) - the D32 gate is not the shape this file declares (ticket 266). A green from a gate whose nails are gone is not a reading: ticket 263 measured an all-green screen with both call lines deleted.\n' \
            "$PROG" "$findings"
        return 1
    fi

    printf '%s: PASS - %s/%s nails declared here are defined once, carry their own ok sentence, are called once at column 0, and run before the first sample at %s:%s\n' \
        "$PROG" "$declared_nail_count" "$declared_nail_count" "$rel" "$anchor"
    printf '%s: NOTE - this tooth pins shape only. That those two ok lines were actually PRINTED is not asserted by any outer step yet (ticket 266 shape 2), and one commit deleting this file together with the CI step that calls it would leave nothing ringing.\n' \
        "$PROG"
    return 0
}

# ---------------------------------------------------------------------------
# the positive control: mutate a COPY of the gate and watch the tooth ring
# ---------------------------------------------------------------------------
BENCH=''
selftest_n=0
selftest_bad=0
OUT=''
RC=0

st_ok() {
    selftest_n=$((selftest_n + 1))
    printf '%s: control %d ok - %s\n' "$PROG" "$selftest_n" "$1"
}
st_bad() {
    selftest_n=$((selftest_n + 1))
    selftest_bad=$((selftest_bad + 1))
    printf '%s: CONTROL FAILED - %s\n' "$PROG" "$1"
    printf '%s:   (bench kept at %s)\n' "$PROG" "$BENCH"
}
run_bench() {  # $1 = the census script to run; sets OUT and RC
    OUT=$BENCH/out-$((selftest_n + 1)).txt
    RC=0
    ( sh "$1" ) > "$OUT" 2>&1 || RC=$?
}
expect_rc() {  # $1 = label, $2 = wanted code
    if [ "$RC" -eq "$2" ]; then
        st_ok "$1 (rc=$RC)"
    else
        st_bad "$1 - expected rc=$2, got rc=$RC; output at $OUT"
    fi
}
expect_has() {  # $1 = label, $2 = text that must be present
    if grep -qF -- "$2" "$OUT"; then
        st_ok "$1 - output carries [$2]"
    else
        st_bad "$1 - expected the output to carry [$2]; output at $OUT"
    fi
}
expect_absent() {  # $1 = label, $2 = text that must NOT be present
    if grep -qF -- "$2" "$OUT"; then
        st_bad "$1 - output carries [$2] but must not; output at $OUT"
    else
        st_ok "$1 - output does not carry [$2]"
    fi
}

run_selftest() {
    BENCH=$(mktemp -d 2>/dev/null) || cannot_decide "cannot create the positive-control bench (mktemp failed)"
    if [ ! -d "$BENCH" ]; then
        cannot_decide "mktemp returned '$BENCH', which is not a directory"
    fi
    printf '%s: positive control bench=%s (kept on disk; this project never deletes temp artifacts)\n' \
        "$PROG" "$BENCH"
    mkdir -p "$BENCH/scripts" || cannot_decide "cannot create $BENCH/scripts"
    self=$here/$(basename "$0")
    cp "$self" "$BENCH/scripts/slo-nail-census.sh" || cannot_decide "cannot copy this script into the bench"
    # ticket 266-r2: normalise the mutation source to LF, whatever the checkout
    # holds. CI signs the target out with CRLF (eol=crlf), and the sed variants
    # below anchor at end of line, which never matches a trailing \r, so every
    # deliberate breakage would silently miss and the self-test would stay green
    # where it must ring. $pristine feeds only the variants; the census reader
    # strips \r itself, and control 27 (CRLF copy) builds its own CRLF from here.
    tr -d '\015' < "$target" > "$BENCH/pristine.ps1" || cannot_decide "cannot copy $rel into the bench"
    bench_target=$BENCH/scripts/slo-check.ps1
    bench_census=$BENCH/scripts/slo-nail-census.sh
    pristine=$BENCH/pristine.ps1
    cp "$pristine" "$bench_target"

    # ---- 1. the healthy copy must be green, and must name what it read
    run_bench "$bench_census"
    expect_rc "the pristine copy is green" 0
    expect_has "the green run names nail 1/2" "nail 1/2 Test-ScriptShapeNail"
    expect_has "the green run names nail 2/2" "nail 2/2 Test-ExitCodeInstrument"
    expect_has "the green run prints its PASS line" "PASS - 2/2 nails"

    # ---- 2. ticket 266 AC#2 shot 1: delete both call lines -> it must ring
    sed -e '/^Test-ScriptShapeNail$/d' -e '/^Test-ExitCodeInstrument$/d' "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "deleting both nail calls turns the census red" 1
    expect_has "the red names the call token" "slo-nail-census-nail-call-missing"
    expect_has "the red names the shape nail and its file" "the declared nail 'Test-ScriptShapeNail' has no column-0, uncommented, standalone call line in slo-check.ps1"
    expect_has "the red names the capability nail" "the declared nail 'Test-ExitCodeInstrument' has no column-0, uncommented, standalone call line"

    # ---- 3. ticket 266 AC#2 shot 2: restore -> green again
    cp "$pristine" "$bench_target"
    run_bench "$bench_census"
    expect_rc "restoring the two call lines turns the census green again" 0
    expect_absent "the restored run carries no FAIL line" "FAIL slo-nail-census"

    # ---- 4. one call deleted: the red must accuse one nail, not both
    sed -e '/^Test-ScriptShapeNail$/d' "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "deleting ONE call turns the census red" 1
    expect_has "the red names the shape nail's missing call" "the declared nail 'Test-ScriptShapeNail' has no column-0"
    expect_absent "the red does not accuse the capability nail" "the declared nail 'Test-ExitCodeInstrument' has no column-0"

    # ---- 5. both calls commented out with line comments
    sed -e 's/^Test-ScriptShapeNail$/# Test-ScriptShapeNail/' \
        -e 's/^Test-ExitCodeInstrument$/# Test-ExitCodeInstrument/' \
        "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "commenting both calls out turns the census red" 1
    expect_has "the commented-out calls are named as missing" "slo-nail-census-nail-call-missing"

    # ---- 6. both calls inside a block comment (the reader's block state, not
    #         the column-0 anchor, is what has to catch this one)
    awk 'BEGIN { open = 0 }
        { if (!open && $0 == "Test-ScriptShapeNail") { print "<#"; open = 1 }
          print
          if (open && $0 == "Test-ExitCodeInstrument") { print "#>"; open = 0 } }' \
        "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "a block comment around both calls turns the census red" 1
    expect_has "the block-commented calls are named as missing" "slo-nail-census-nail-call-missing"

    # ---- 7. a definition deleted (the call line survives, the nail does not)
    sed -e '/^function Test-ScriptShapeNail {$/d' "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "deleting a nail's definition turns the census red" 1
    expect_has "the red names the missing definition token" "slo-nail-census-nail-definition-missing"

    # ---- 8. a nail's own ok sentence deleted
    sed -e 's/instrument self-check ok/instrument self-check [sentence removed by the census control]/' \
        "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "deleting a nail's ok sentence turns the census red" 1
    expect_has "the red names the missing ok sentence token" "slo-nail-census-nail-ok-sentence-missing"

    # ---- 9. the target missing altogether: CANNOT DECIDE, never green
    mv "$bench_target" "$BENCH/target-moved-away.ps1"
    run_bench "$bench_census"
    expect_rc "a missing target is CANNOT DECIDE (exit 2), never green" 2
    expect_has "the refusal carries its token" "slo-nail-census-cannot-decide"
    expect_has "the refusal says it is not a pass" "NOT a pass"

    # ---- 10. the target empty: same honesty tier
    : > "$bench_target"
    run_bench "$bench_census"
    expect_rc "an empty target is CANNOT DECIDE (exit 2)" 2
    expect_has "the empty-file refusal carries its token" "slo-nail-census-cannot-decide"

    # ---- 11. a CRLF copy stays green: the \r strip is load-bearing, because
    #          git's attributes say eol=crlf for the target while the work tree
    #          currently holds LF (measured 2026-10-15, ticket 266-r1 evidence).
    awk '{ printf "%s\r\n", $0 }' "$pristine" > "$bench_target"
    run_bench "$bench_census"
    expect_rc "a CRLF copy of the target stays green" 0

    # ---- 12. the declaration count edited alone: drift is a red, not a
    #          silent redefinition
    cp "$bench_census" "$BENCH/scripts/census-drift.sh"
    sed -e 's/^declared_nail_count=2$/declared_nail_count=3/' \
        "$BENCH/scripts/census-drift.sh" > "$BENCH/scripts/census-drift.tmp"
    mv "$BENCH/scripts/census-drift.tmp" "$BENCH/scripts/census-drift.sh"
    run_bench "$BENCH/scripts/census-drift.sh"
    expect_rc "a count that no longer matches its tables turns the census red" 1
    expect_has "the red names the drift token" "slo-nail-census-declaration-drift"

    # ---- 13. the two calls moved below the sampler: the ordering claim has teeth
    sed -e '/^Test-ScriptShapeNail$/d' -e '/^Test-ExitCodeInstrument$/d' "$pristine" > "$bench_target"
    printf 'Test-ScriptShapeNail\nTest-ExitCodeInstrument\n' >> "$bench_target"
    run_bench "$bench_census"
    expect_rc "calls moved below the first sample turn the census red" 1
    expect_has "the red names the ordering token" "slo-nail-census-ordering"
    expect_has "the red says why the position matters" "not before the first sampling call"

    # leave the bench in its healthy state: the copies are the only things this
    # control ever wrote
    cp "$pristine" "$bench_target"

    printf '%s: positive control %s - %d control(s) ran, %d failed\n' \
        "$PROG" "$([ "$selftest_bad" -eq 0 ] && echo PASSED || echo FAILED)" \
        "$selftest_n" "$selftest_bad"
    if [ "$selftest_bad" -eq 0 ]; then return 0; fi
    return 1
}

# ---------------------------------------------------------------------------
# main
# ---------------------------------------------------------------------------
if [ "$MODE" = "selftest" ] || [ "$MODE" = "both" ]; then
    if ! run_selftest; then
        printf '%s: positive control RED - the census is not run, because a ruler that cannot go red does not get to print a reading\n' \
            "$PROG" >&2
        exit 1
    fi
    if [ "$MODE" = "selftest" ]; then
        exit 0
    fi
fi

rc=0
census || rc=$?
exit "$rc"
