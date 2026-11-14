#!/bin/sh
# scripts/slo-nail-census.sh - the OUTER tooth over the two self-locking nails of
# the D32 gate (ticket 266, shape 1-A) and, since ticket 274, over the frontend
# provenance gate of the build line (scripts/build.ps1).
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
#   and, since ticket 274, that the BUILD gate (scripts/build.ps1) carries its
#   frontend provenance gate - the nail 274-r2 (commit 91c90aa1) built into step
#   2 to stop a stubbed npm from shipping last month's page bytes (ticket 274
#   AC#10) with array-to-array rosters (ticket 274 AC#11):
#
#       the build gate has exactly 1 provenance gate shape
#       $preRoster  (the pre-npm snapshot array)   = line 158 of build.ps1
#       $freshNames (the array-to-array diff)      = line 209 of build.ps1
#       the unprovenanced Fail                  = line 212 of build.ps1
#       the npm run build call                     = line 165 of build.ps1
#
#   Deleting any of those lines (or commenting them out) now leaves the
#   declaration behind saying the build gate should carry them - the same
#   shape the D32 declaration gives the ticket-263 delete-nails mutant.
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
# THE BUILD GATE DECLARATION (ticket 274) - lives here, not in build.ps1
# ---------------------------------------------------------------------------
# build.ps1's frontend provenance gate (ticket 274 AC#10/AC#11, commits
# 35633445/91c90aa1/674c1920) is four load-bearing shapes, each pinned by an
# exact token that must occur EXACTLY once, uncommented:
#   bld_snapshot   the pre-npm dist roster array (AC#10's snapshot)
#   bld_diff       the array-to-array roster diff (AC#11's array-in-array-out)
#   bld_redline    the named-cause Fail for unprovenanced page bytes
#   bld_npmcall    the npm invocation the snapshot must precede (the ordering
#                  anchor for the three above, standing in for the gate's
#                  "first sample")
# Ordering: the three gate shapes must sit BEFORE the npm call (snapshot) or
# between snapshot and the roster print that claims success - pinned simply as
# "snapshot and diff and redline all precede the npm call" is wrong for the
# redline (it does not), so the declaration is: snapshot before npm call; diff
# and redline after the npm call and before the 'go build' that embeds dist.
# That last anchor is exact and unique in build.ps1:
bld_build_anchor='& $go.Source build -trimpath -ldflags $ldflags -o (Join-Path $outDir '"'"'wisp.exe'"'"') ./cmd/wisp'
# Tokens are matched as substrings of comment-stripped code lines, each
# expected exactly once in the whole file:
bld_snapshot='$preRoster = @($preDistFiles | Sort-Object FullName | ForEach-Object { $_.FullName.Substring($DistDir.Length + 1) })'
bld_diff='$freshNames = @($postRoster | Where-Object { $preRoster -notcontains $_ })'
bld_redline_token='unprovenanced page bytes - nothing in frontend/dist was produced by this build'
bld_npmcall='& $npm.Source run build'
# The success line's provenance half (its ok sentence); double-quoted shell
# string, so the apostrophe is literal and the sigil is escaped:
bld_ok="provenance: this run's npm added new artifact name(s)"

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
(ticket 266) and over the build gate's frontend provenance gate (ticket 274). It
reads scripts/slo-check.ps1 and checks that the two nails this file declares are
still defined, still carry their own ok sentence, still have exactly one
uncommented column-0 call line, and are still called before the first SLO sample.
It also reads scripts/build.ps1 and checks that the frontend provenance gate
(ticket 274 AC#10/AC#11: pre-npm snapshot, array-to-array roster diff, the
unprovenanced Fail, the npm call, and the success line's provenance sentence)
still exists exactly once and in the declared order around the npm call and the
go build that embeds dist.

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
# and, since ticket 274, one row per build-gate token:
#   BTOK <key> <count> <firstline>
analyze() {  # $1 = the file to read
    tr -d '\015' < "$1" | awk \
        -v names="$nail_names" \
        -v oks="$nail_ok" \
        -v anchor="$sampling_anchor" \
        -v t1="$bld_snapshot" \
        -v t2="$bld_npmcall" \
        -v t3="$bld_diff" \
        -v t4="$bld_redline_token" \
        -v t5="$bld_ok" \
        -v t6="$bld_build_anchor" \
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
            # ticket 274 build-gate tokens: substring counts on comment-stripped
            # code lines, first hit line recorded for the verdict text.
            nt[1] = t1; nt[2] = t2; nt[3] = t3; nt[4] = t4; nt[5] = t5; nt[6] = t6
            for (j = 1; j <= 6; j++) {
                if (index(code, nt[j]) > 0) {
                    bcount[j]++
                    if (bline[j] == 0) { bline[j] = NR }
                }
            }
        }
        END {
            if (nn != nok) { exit 3 }
            for (i = 1; i <= nn; i++) {
                printf "NAIL %s %d %d %d %d %d %d %s\n", nm[i], \
                    defcount[i] + 0, defline[i] + 0, okcount[i] + 0, okline[i] + 0, \
                    callcount[i] + 0, callline[i] + 0, ok[i]
            }
            printf "ANCHOR %d\n", anchorline + 0
            printf "BTOK snapshot %d %d\n", bcount[1] + 0, bline[1] + 0
            printf "BTOK npmcall %d %d\n", bcount[2] + 0, bline[2] + 0
            printf "BTOK diff %d %d\n", bcount[3] + 0, bline[3] + 0
            printf "BTOK redline %d %d\n", bcount[4] + 0, bline[4] + 0
            printf "BTOK okline %d %d\n", bcount[5] + 0, bline[5] + 0
            printf "BTOK buildanchor %d %d\n", bcount[6] + 0, bline[6] + 0
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

    # ---------------------------------------------------------------------------
    # the build-gate census (ticket 274): the same rows file carries the
    # BTOK lines for scripts/build.ps1's frontend provenance gate.
    # ---------------------------------------------------------------------------
    bld_target=$here/build.ps1
    bld_rel=scripts/build.ps1
    if [ ! -f "$bld_target" ] || [ ! -r "$bld_target" ] || [ ! -s "$bld_target" ]; then
        cannot_decide "the build gate $bld_rel is missing, unreadable, or empty - ticket 274's provenance declaration cannot be checked against anything, and that is not a pass"
    fi

    # The build-gate rows come from reading THE BUILD GATE, not the D32 gate:
    # the same reader (comment strip included) runs against build.ps1 alone.
    bld_rowsfile=$(mktemp 2>/dev/null) || cannot_decide "cannot create a temporary file for the build-gate rows"
    if ! analyze "$bld_target" > "$bld_rowsfile"; then
        cannot_decide "the reader (tr | awk) failed on $bld_rel - the build gate could not be turned into a reading"
    fi
    if [ ! -s "$bld_rowsfile" ]; then
        cannot_decide "the reader produced no rows at all for $bld_rel"
    fi

    printf '%s: build gate %s - declared provenance tokens: snapshot / npmcall / diff / redline / okline / buildanchor, each exactly once\n' \
        "$PROG" "$bld_rel"

    bt_snapshot=0; bt_npmcall=0; bt_diff=0; bt_redline=0; bt_okline=0; bt_buildanchor=0
    bl_snapshot=0; bl_npmcall=0; bl_diff=0; bl_redline=0; bl_okline=0; bl_buildanchor=0
    bt_rows=0
    while read -r kind name cnt ln; do
        if [ "$kind" != "BTOK" ]; then continue; fi
        bt_rows=$((bt_rows + 1))
        case $name in
            snapshot)    bt_snapshot=$cnt;    bl_snapshot=$ln ;;
            npmcall)     bt_npmcall=$cnt;     bl_npmcall=$ln ;;
            diff)        bt_diff=$cnt;        bl_diff=$ln ;;
            redline)     bt_redline=$cnt;     bl_redline=$ln ;;
            okline)      bt_okline=$cnt;      bl_okline=$ln ;;
            buildanchor) bt_buildanchor=$cnt; bl_buildanchor=$ln ;;
            *) cannot_decide "the reader emitted an unknown build-gate row [$name] - the rows do not cover the declaration" ;;
        esac
    done < "$bld_rowsfile"
    if [ "$bt_rows" -ne 6 ]; then
        cannot_decide "the reader returned $bt_rows build-gate row(s) for a declaration of 6 - no verdict can be built from them"
    fi

    bld_findings=0
    bld_fail() {  # $1 = token key, $2 = the finding
        bld_findings=$((bld_findings + 1))
        printf '%s: FAIL slo-nail-census-build-%s: %s\n' "$PROG" "$1" "$2"
    }

    # each token exactly once, uncommented
    for pair in "snapshot:$bt_snapshot:$bl_snapshot" "npmcall:$bt_npmcall:$bl_npmcall" \
                "diff:$bt_diff:$bl_diff" "redline:$bt_redline:$bl_redline" \
                "okline:$bt_okline:$bl_okline" "buildanchor:$bt_buildanchor:$bl_buildanchor"; do
        key=${pair%%:*}; rest=${pair#*:}; cnt=${rest%%:*}; ln=${rest#*:}
        if [ "$cnt" -eq 0 ]; then
            bld_fail "$key" "the declared build-gate token [$key] occurs zero times in uncommented code of $bld_rel - the ticket 274 provenance gate shape this file declares has been deleted or commented out (ticket 274 AC#10/AC#11; same failure family as ticket 263's delete-nails mutant)."
        elif [ "$cnt" -gt 1 ]; then
            bld_fail "$key" "the declared build-gate token [$key] occurs $cnt times in $bld_rel (first at $bld_rel:$ln); this file declares exactly one."
        fi
    done

    # ordering: snapshot precedes the npm call; diff, redline and okline sit
    # between the npm call and the go build that embeds dist.
    if [ "$bt_snapshot" -gt 0 ] && [ "$bt_npmcall" -gt 0 ] && [ "$bl_snapshot" -gt "$bl_npmcall" ]; then
        bld_fail ordering "the pre-npm snapshot sits at $bld_rel:$bl_snapshot, AFTER the npm run build call at $bld_rel:$bl_npmcall. A snapshot taken after npm cannot prove provenance (ticket 274 AC#10)."
    fi
    for pair in "diff:$bl_diff" "redline:$bl_redline" "okline:$bl_okline"; do
        key=${pair%%:*}; ln=${pair#*:}
        if [ "$ln" -gt 0 ] && [ "$bl_npmcall" -gt 0 ] && [ "$ln" -lt "$bl_npmcall" ]; then
            bld_fail ordering "the build-gate token [$key] sits at $bld_rel:$ln, BEFORE the npm run build call at $bld_rel:$bl_npmcall - the roster comparison must read dist after npm ran, not before (ticket 274 AC#10)."
        fi
        if [ "$ln" -gt 0 ] && [ "$bl_buildanchor" -gt 0 ] && [ "$ln" -gt "$bl_buildanchor" ]; then
            bld_fail ordering "the build-gate token [$key] sits at $bld_rel:$ln, AFTER the go build call at $bld_rel:$bl_buildanchor - go:embed reads dist during that build, so a provenance verdict printed after it protects nothing (ticket 274 AC#10)."
        fi
    done

    if [ "$bld_findings" -gt 0 ]; then
        printf '%s: FAIL %s build-gate declaration violation(s) - the provenance gate of %s is not the shape this file declares (ticket 274). A green from a build whose provenance gate is gone ships page bytes nobody can attribute to this build.\n' \
            "$PROG" "$bld_findings" "$bld_rel"
        return 1
    fi

    printf '%s: build gate PASS - snapshot:%s npmcall:%s diff:%s redline:%s okline:%s buildanchor:%s in the declared order (ticket 274)\n' \
        "$PROG" "$bl_snapshot" "$bl_npmcall" "$bl_diff" "$bl_redline" "$bl_okline" "$bl_buildanchor"

    if [ "$findings" -gt 0 ]; then
        printf '%s: FAIL %s declaration violation(s) - the D32 gate is not the shape this file declares (ticket 266). A green from a gate whose nails are gone is not a reading: ticket 263 measured an all-green screen with both call lines deleted.\n' \
            "$PROG" "$findings"
        return 1
    fi

    printf '%s: PASS - %s/%s nails declared here are defined once, carry their own ok sentence, are called once at column 0, and run before the first sample at %s:%s\n' \
        "$PROG" "$declared_nail_count" "$declared_nail_count" "$rel" "$anchor"
    printf '%s: NOTE - this tooth pins shape only. That those two ok lines were actually PRINTED is not asserted by any outer step yet (ticket 266 shape 2), and one commit deleting this file together with the CI step that calls it would leave nothing ringing. The build-gate half (ticket 274) pins shape only in the same way.\n' \
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
    # ticket 274: the bench also carries a copy of the build gate, so the
    # build-gate half of the declaration gets its own positive control. Normalised
    # to LF like pristine.ps1 (the reader strips \r itself).
    tr -d '\015' < "$here/build.ps1" > "$BENCH/pristine-build.ps1" \
        || cannot_decide "cannot copy scripts/build.ps1 into the bench"
    bench_build=$BENCH/scripts/build.ps1
    pristine_build=$BENCH/pristine-build.ps1
    cp "$pristine_build" "$bench_build"

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

    # ---- 14. ticket 274: the healthy build gate copy is green, and the run
    #          names the build gate's own PASS line. The D32 half is restored
    #          first, so this control reads two healthy copies.
    cp "$pristine" "$bench_target"
    run_bench "$bench_census"
    expect_rc "the pristine build gate is green" 0
    expect_has "the green run prints the build gate PASS line" "build gate PASS"

    # ---- 15. ticket 274 shot 1: delete the provenance diff line -> red
    sed -e '/^\$freshNames = @(\$postRoster | Where-Object { \$preRoster -notcontains \$_ })$/d' \
        "$pristine_build" > "$bench_build"
    run_bench "$bench_census"
    expect_rc "deleting the roster diff turns the build-gate census red" 1
    expect_has "the red names the missing diff token" "slo-nail-census-build-diff"

    # ---- 16. ticket 274 shot 2: comment out the redline -> red
    sed -e 's/^    Fail "frontend step: npm run build reported success but frontend\/dist holds no artifact name outside the snapshot taken before this run'"'"'s npm /    # Fail "frontend step: npm run build reported success but frontend\/dist holds no artifact name outside the snapshot taken before this run'"'"'s npm /' \
        "$pristine_build" > "$bench_build"
    run_bench "$bench_census"
    expect_rc "commenting out the unprovenanced Fail turns the build-gate census red" 1
    expect_has "the red names the missing redline token" "slo-nail-census-build-redline"

    # ---- 17. ticket 274 shot 3: restore -> green again
    cp "$pristine_build" "$bench_build"
    run_bench "$bench_census"
    expect_rc "restoring the build gate turns the census green again" 0
    expect_absent "the restored run carries no FAIL line" "FAIL slo-nail-census"

    # ---- 18. ticket 274 shot 4: snapshot moved AFTER the npm call -> ordering red
    #          (delete the snapshot line, re-insert it right after the npm call)
    sed -e '/^\$preRoster = @(\$preDistFiles | Sort-Object FullName | ForEach-Object { \$_\.FullName\.Substring(\$DistDir\.Length + 1) })$/d' \
        "$pristine_build" > "$bench_build"
    awk -v ins='$preRoster = @($preDistFiles | Sort-Object FullName | ForEach-Object { $_.FullName.Substring($DistDir.Length + 1) })' '
        { print
          if ($0 == "    & $npm.Source run build") { print ins } }' \
        "$bench_build" > "$bench_build.tmp" && mv "$bench_build.tmp" "$bench_build"
    run_bench "$bench_census"
    expect_rc "a snapshot taken after npm turns the build-gate census red" 1
    expect_has "the ordering red names the snapshot-after-npm cause" "AFTER the npm run build call"

    # ---- 19. ticket 274 shot 5: the ok sentence deleted -> red
    sed -e 's/provenance: this run'"'"'s npm added new artifact name(s)/[provenance sentence removed by the census control]/' \
        "$pristine_build" > "$bench_build"
    run_bench "$bench_census"
    expect_rc "deleting the provenance ok sentence turns the build-gate census red" 1
    expect_has "the red names the missing okline token" "slo-nail-census-build-okline"

    # ---- 20. the whole build gate deleted (shots 15-19 all applied at once):
    #          five tokens missing, five reds named
    grep -v -e '^\$preRoster = @(\$preDistFiles' \
            -e '^\$freshNames = @(\$postRoster' \
            -e 'unprovenanced page bytes' \
            -e '^\s*& \$npm\.Source run build' \
            -e 'provenance: this run' \
        "$pristine_build" > "$bench_build"
    run_bench "$bench_census"
    expect_rc "deleting the whole provenance gate turns the census red" 1
    expect_has "the whole-gate red names the snapshot" "slo-nail-census-build-snapshot"
    expect_has "the whole-gate red names the npmcall" "slo-nail-census-build-npmcall"
    expect_has "the whole-gate red names the diff" "slo-nail-census-build-diff"
    expect_has "the whole-gate red names the redline" "slo-nail-census-build-redline"
    expect_has "the whole-gate red names the okline" "slo-nail-census-build-okline"

    # ---- 21. ticket 274 AC#11-shaped mutation: the array diff re-joined into a
    #          string split (the -join/-split shred that AC#11 closed) - the
    #          diff token itself is gone in that shape, so the gate must ring
    sed -e 's/^\$freshNames = @(\$postRoster | Where-Object { \$preRoster -notcontains \$_ })$/\$freshNames = @((\$postRoster -join '"'"' '"'"') -split '"'"' '"'"' | Where-Object { (\$preRoster -join '"'"' '"'"') -notcontains \$_ })/' \
        "$pristine_build" > "$bench_build"
    run_bench "$bench_census"
    expect_rc "reverting the array diff to the join-split shape turns the census red" 1
    expect_has "the join-split red names the diff token" "slo-nail-census-build-diff"

    # leave the bench in its healthy state: the copies are the only things this
    # control ever wrote
    cp "$pristine" "$bench_target"
    cp "$pristine_build" "$bench_build"

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
