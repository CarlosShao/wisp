#!/bin/sh
# scripts/ticket-lint.sh - the ticket-pool linter, v1 (report only).
#
# WHAT THIS IS
#   A read-only census of .scratch/wisp/issues/ that answers, in one run, the
#   seven questions the harness-reform plan asks of the ticket pool
#   (docs/audit/2026-10-25-harness-reform-plan.md section 3.2, items 1-7).
#   It touches no ticket, writes nothing inside the repository, and decides
#   nothing: every finding is printed for a human, nothing is repaired.
#
# V1 NEVER BLOCKS A BUILD - A DESIGN REQUIREMENT, NOT AN OVERSIGHT
#   The plan states the shape: "v1 reports only, v2 goes into CI". So the scan
#   exits 0 whenever it RAN TO COMPLETION, whether or not it found problems;
#   findings live in the output, never in the exit code. Exit 1 is reserved for
#   "--self-test found a control that did not behave", exit 2 for "this gate
#   refuses to answer" (no work tree, no pool, no threshold line, empty
#   denominator). Wiring this into .github/workflows/ci.yml is a v2 action and is
#   deliberately NOT done here - if this file starts failing a build, that is
#   the bug, not the achievement.
#
# WHY A SEPARATE RULER FROM scripts/check-path-length-budget.sh (item 5)
#   Item 5 shares its THRESHOLD with that gate and must not restate it: the hat is
#   read out of scripts/check-path-length-budget.sh at run time, from the one
#   line that defines it (HAT_NAME_LIMIT=NNN). This script never writes that file
#   and never carries a second copy of the number. Measured 2026-10-26:
#   HAT_NAME_LIMIT=100, README rule 9. The 262 gate already turns an over-budget
#   TRACKED PATH red, so item 5 is not "another red" - it is the POOL view of the
#   same hat, annotated with which offenders the 262 gate grandfathers by name
#   (its roster), so "new ticket name over the hat" becomes a reading instead of
#   a guess.
#
# THE STATUS-LINE PROBLEM - THE WHOLE REASON THIS SCRIPT IS CAREFUL
#   Ticket faces do not carry one status format. Measured 2026-10-26 over 365
#   tickets, four head shapes plus one continuation shape are in use:
#       A   **Status:** done                             (Sep tickets, 174 files)
#       B   - Status: **DONE（...）**                    (Oct tickets, 137 files)
#       C   - **Status**：**待派**（...）                 (17 files)
#       D   Status: OPEN（...）                           (4 files)
#       E   -> **accepted-done**（...）                   (continuation line)
#   A ruler that knows only shape B misses every Sep ticket; one that knows only
#   A misses every Oct ticket; one that reads only the FIRST status line misses
#   the tickets whose done marker lives on an appended E line. Ticket 101 is the
#   canonical example: its head says "open", and the appended line says
#   "-> **accepted-done**", which is why item 2 asks "does ANY status line carry
#   a done word" and not "does the first one". Measured 2026-10-26: five tickets
#   (95, 97, 99, 101, 102) owe their done marker to an appended line alone, so a
#   first-line-only ruler would print five false reds.
#   Shape E is deliberately narrow, because "a line that starts with an arrow" is
#   NOT by itself a status: this pool also contains body prose such as
#   "  -> within cap or alert" and "-> exit 1". A line counts as shape E only
#   when it starts at column 0, follows a recognised head, and its payload is ONE
#   bolded token (-> **accepted-done**（...）). Every appended status in the pool
#   is written that way; nothing else is admitted.
#   KNOWN LIMITATION (v1): a fifth shape, a blockquote-prefixed
#   "> **Status ... **：", is NOT recognised. Measured 2026-10-26: no ticket
#   relies on it alone and none of the 20 no-status tickets carries one, so the
#   omission costs no verdict today; a v2 ruler should take it in.
#
# WHY NON-ASCII IS KEPT OUT OF BRACKET EXPRESSIONS
#   The awk that runs here is byte-oriented unless the locale is UTF-8; this
#   script PROBES that at run time and prints the answer (the same probe
#   scripts/check-path-length-budget.sh uses). In byte mode a bracket expression
#   like [（(] is a set of BYTES, so its first alternative can match the leading
#   byte of a three-byte character and split it. Measured here 2026-10-26: that
#   produced mojibake status values and a wrong reading. Every non-ASCII literal
#   in this file is therefore used as a COMPLETE SEQUENCE in an alternation
#   ( :|： ), never inside a bracket expression.
#
# WHAT IS READ, AND FROM WHERE (all read-only)
#   <root>/.scratch/wisp/issues/*.md            the pool (README.md excluded)
#   <root>/.scratch/wisp/voided/*.md            withdrawn drafts, for item 4
#   <root>/scripts/check-path-length-budget.sh  HAT_NAME_LIMIT and its roster
#   the pool files themselves                   every status line, every
#                                               line-start checkbox, every full
#                                               date, every Blocked-by field
#
# THE SEVEN ITEMS
#   1 (red)    a -done ticket that still has a line-start unchecked box
#   2 (red)    a -done ticket with no done word on ANY status line
#   3a (red)   a live ticket with no status line at all
#   3b (red)   a live ticket whose status value is outside the vocabulary
#   3c (note)  a value legible only because the v1 alias table normalised it
#   4 (yellow) two files on one ticket number, and the missing-number roster
#              reconciled against voided/
#   5 (red)    a ticket name longer than the hat (262-roster offenders noted)
#   6 (yellow) ageing: in-progress over 72h, review over 14d,
#              blocked-on-owner over 21d
#   7 (stat)   done / live / void / waiting-on-owner / ready-right-now counts
#
#   The vocabulary (item 3b) is the canonical set from the plan's schema v2
#   (ready, in-progress, review, blocked, blocked-on-owner, done, void, reserved)
#   plus a documented alias table for the spellings the pool actually uses. The
#   alias table exists so item 3b reports DEFECTS (values a scheduler could not
#   act on) instead of re-reporting the known migration debt - and every alias
#   hit is still printed as 3c, so that debt stays visible and countable.
#
# UNITS AND DATES
#   Name length is measured in characters when this awk is multibyte-aware and in
#   bytes otherwise; the unit is printed on every run. Dates come from the LAST
#   full YYYY-MM-DD date on the ticket face, and ageing is a calendar-day
#   difference against today, so "72h" is read as "more than 3 days". Day
#   granularity is deliberate: a ticket face carries dates, not timestamps, and
#   inventing clock precision would be a lie.
#
# EXIT CODES
#   0  the scan ran to completion (with or without findings) - the v1 contract
#   0  --self-test: every control behaved
#   1  --self-test: a control did not behave (the ruler has no teeth)
#   2  refused to answer (no work tree, no pool, no threshold line, or a zero
#      denominator)
#
# USAGE
#   sh scripts/ticket-lint.sh                    scan the pool (report only)
#   sh scripts/ticket-lint.sh --self-test        run the built-in controls only
#   sh scripts/ticket-lint.sh --with-self-test   controls first, then the scan
#   sh scripts/ticket-lint.sh --help

set -u

PROG=${0##*/}

usage() {
    cat <<'EOF'
usage: sh scripts/ticket-lint.sh [--self-test | --with-self-test]

  (no option)        census the ticket pool; report only, always exits 0
  --self-test        run the built-in controls on a throwaway pool
  --with-self-test   controls first, then the scan
  -h, --help         this text

exit: 0 ran to completion (findings do not change this), 1 a self-test control
      failed, 2 refused to answer.
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

ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$ROOT" ]; then
    printf '%s: not inside a git work tree - refusing to guess a pool\n' "$PROG" >&2
    exit 2
fi

POOL=$ROOT/.scratch/wisp/issues
VOIDED=$ROOT/.scratch/wisp/voided
LENGTH_GATE=$ROOT/scripts/check-path-length-budget.sh

if [ ! -f "$LENGTH_GATE" ]; then
    printf '%s: REFUSE - %s is missing, and item 5 shares its threshold with that gate (there is exactly one definition of the hat and this script must not carry a second one)\n' \
        "$PROG" "$LENGTH_GATE" >&2
    exit 2
fi

# The one definition of the hat lives in the 262 gate. Read it, never restate it.
HAT_NAME_LIMIT=$(sed -n 's/^HAT_NAME_LIMIT=\([0-9][0-9]*\).*$/\1/p' "$LENGTH_GATE" | head -n 1)
case "$HAT_NAME_LIMIT" in
    ''|*[!0-9]*)
        printf '%s: REFUSE - could not read HAT_NAME_LIMIT out of %s. This gate does not own that number; if the 262 gate renamed it, teach this gate the new name rather than giving it a copy of the value\n' \
            "$PROG" "$LENGTH_GATE" >&2
        exit 2 ;;
esac

# Unit probe, same in spirit as the 262 gate's: never assume whether awk counts
# characters or bytes. "a" + U+25FF + "b" is 5 bytes and 3 characters.
SAMPLE=$(printf 'a\342\226\277b')
UNIT=$(awk -v s="$SAMPLE" 'BEGIN{print length(s)}' 2>/dev/null)
if [ "$UNIT" = "3" ]; then
    UNITNAME=characters
elif [ "$UNIT" = "5" ]; then
    UNITNAME=bytes
else
    UNITNAME=undetermined
fi

TODAY=$(date +%Y-%m-%d)

# ---------------------------------------------------------------------------
# the census: one awk pass over the pool, one TSV record per ticket
# ---------------------------------------------------------------------------
# Emitted fields:
#   C <num> <suffix> <unchecked> <checked> <lastdate> <nstatus> <nhead>
#     <active> <doneflag> <blockedflag> <name> <formats>
census() {
    awk '
    function trim(s) { sub(/^[ \t]+/, "", s); sub(/[ \t]+$/, "", s); return s }
    # A status VALUE is the short token at the head of the line: everything before
    # the first bracket, comma, full stop or space. Prose after the token is
    # dropped and emphasis markers go, so **done** and done become one word.
    function cut(v,  t) {
        gsub(/\*/, "", v); gsub(/`/, "", v); gsub(/~/, "", v)
        gsub(/in progress/, "in-progress", v)
        gsub(/In progress/, "in-progress", v)
        v = trim(v); t = v
        sub(/（.*$/, "", t); sub(/\(.*$/, "", t)
        sub(/，.*$/, "", t); sub(/,.*$/, "", t)
        sub(/。.*$/, "", t); sub(/、.*$/, "", t)
        sub(/；.*$/, "", t); sub(/;.*$/, "", t)
        sub(/[ \t].*$/, "", t)
        return trim(t)
    }
    function flush(  base, num) {
        if (fname == "") return
        base = fname; sub(/.*\//, "", base)
        if (base == "README.md") { fname = ""; return }
        num = base; sub(/-.*$/, "", num)
        sfx = "live"
        if (base ~ /-done\.md$/) sfx = "done"
        else if (base ~ /-void\.md$/) sfx = "void"
        if (active == "") active = "-"
        printf "C\t%s\t%s\t%d\t%d\t%s\t%d\t%d\t%s\t%d\t%d\t%s\t%s\n", \
            num, sfx, un + 0, ck + 0, (lastd == "" ? "-" : lastd), nstat + 0, nhead + 0, \
            active, doneflag + 0, blk + 0, base, (fmts == "" ? "-" : fmts)
    }
    FNR == 1 {
        flush()
        fname = FILENAME; un = 0; ck = 0; lastd = ""; nstat = 0; nhead = 0
        active = ""; doneflag = 0; blk = 0; fmts = ""; sfx = "live"
    }
    {
        line = $0
        if (line ~ /^[ \t]*-[ \t]*\[ \]/) un++
        else if (line ~ /^[ \t]*-[ \t]*\[[xX]\]/) ck++
        if (match(line, /20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]/)) lastd = substr(line, RSTART, 10)
        # the Blocked-by field, both spellings; None / - / n/a count as empty
        if (line ~ /^[ \t]*(-[ \t]*)?\*\*[Bb]locked[ -]?[Bb]y[ \t]*:?\*\*[ \t]*:?/) {
            bv = line
            sub(/^[^:]*:[ \t]*/, "", bv)
            gsub(/\*/, "", bv); bv = trim(bv)
            sub(/[ \t].*$/, "", bv)
            bv = tolower(bv)
            if (bv != "" && bv != "none" && bv != "-" && bv != "n/a" && bv != "na" && bv != "(none)") blk = 1
        }
        v = ""; fmt = ""
        if (line ~ /^[ \t]*\*\*Status:\*\*/) {
            fmt = "A"; v = line; sub(/^[ \t]*\*\*Status:\*\*[ \t]*/, "", v)
        } else if (line ~ /^[ \t]*-[ \t]+Status[ \t]*(:|：)/) {
            fmt = "B"; v = line; sub(/^[ \t]*-[ \t]+Status[ \t]*(:|：)[ \t]*/, "", v)
        } else if (line ~ /^[ \t]*-[ \t]*\*\*Status\*\*[ \t]*(:|：)/) {
            fmt = "C"; v = line; sub(/^[ \t]*-[ \t]*\*\*Status\*\*[ \t]*(:|：)[ \t]*/, "", v)
        } else if (line ~ /^[ \t]*Status[ \t]*(:|：)/) {
            fmt = "D"; v = line; sub(/^[ \t]*Status[ \t]*(:|：)[ \t]*/, "", v)
        } else if (nhead > 0 && line ~ /^→[ \t]*\*\*[^* \t]+\*\*/) {
            # shape E: an appended transition. It must be a top-level arrow whose
            # payload is ONE bolded token. That is how every real appended status
            # in this pool is written, and it is what keeps body prose such as
            # "  -> within cap or alert" or "-> exit 1" out of the status block.
            fmt = "E"; v = line; sub(/^→[ \t]*/, "", v)
        }
        if (v == "") next
        nstat++
        if (fmt != "E") nhead++
        if (index(fmts, fmt) == 0) fmts = fmts fmt
        val = cut(v)
        if (val == "") val = "(empty)"
        active = val
        if (doneflag == 0 && val != "") {
            kind = tolower(val)
            if (index(kind, "done") > 0 || index(val, "结案") > 0) doneflag = 1
        }
    }
    END { flush() }' "$POOL"/*.md
}

# ---------------------------------------------------------------------------
# the report: one awk pass over roster + voided + census
# ---------------------------------------------------------------------------
report() {
    awk -v PROG="$PROG" -v today="$TODAY" -v hat="$HAT_NAME_LIMIT" -v unit="$UNITNAME" '
    function trim(s) { sub(/^[ \t]+/, "", s); sub(/[ \t]+$/, "", s); return s }
    # days since a fixed epoch: ageing needs no mktime and no date arithmetic
    function days(y, m, d,  yy, mm) {
        yy = (m <= 2) ? y - 1 : y
        mm = (m <= 2) ? m + 9 : m - 3
        return 365*yy + int(yy/4) - int(yy/100) + int(yy/400) + int((153*mm + 2)/5) + d
    }
    function dayof(s) { return days(substr(s,1,4)+0, substr(s,6,2)+0, substr(s,9,2)+0) }
    # canonical vocabulary + the v1 alias table. An alias maps a spelling the pool
    # actually uses onto the value a scheduler can act on; anything unmapped is a
    # DEFECT (3b), not debt.
    function canon(v,  t) {
        t = tolower(v)
        if (t == "ready" || t == "ready-for-agent" || t == "open" || t == "filed") return "ready"
        if (v == "已立" || v == "待派" || v == "立而不派" || v == "未开工" || v == "在派" || v == "在派队列里") return "ready"
        if (t == "in-progress" || t == "inprogress" || t == "active") return "in-progress"
        if (v == "进行中" || v == "在做") return "in-progress"
        if (t == "review" || t == "in-review" || t == "ready-for-review") return "review"
        if (t == "rejected-needs-fix" || t == "needs-fix") return "in-progress"
        if (t == "blocked") return "blocked"
        if (t == "blocked-on-owner") return "blocked-on-owner"
        if (t == "done" || t == "accepted-done" || t == "implemented-agent-done" || t == "done-as-refutation") return "done"
        if (v == "已结案" || v == "结案") return "done"
        if (t == "void" || t == "withdrawn" || t == "withdraw") return "void"
        if (t == "reserved") return "reserved"
        return ""
    }
    function agecut(c) {
        if (c == "in-progress") return 3
        if (c == "review") return 14
        if (c == "blocked-on-owner") return 21
        return 0
    }
    $1 == "R" { roster[$2] = 1; next }
    $1 == "V" { voided[$2] = 1; next }
    $1 != "C" { next }
    {
        num = $2; sfx = $3; un = $4; ck = $5; lastd = $6; nstat = $7; nhead = $8
        active = $9; doneflag = $10; blk = $11; name = $12
        tickets++; sfxcount[sfx]++
        # the gap walk below counts with integers, so the seen-set is keyed by the
        # NUMBER and not by the zero-padded spelling ("01" and 1 are one ticket)
        numseen[num + 0] = 1
        if (num + 0 > maxnum) maxnum = num + 0
        if (minnum == 0 || num + 0 < minnum) minnum = num + 0
        if (num in numseen2) dup[num] = numseen2[num] " + " name
        else numseen2[num] = name
        # -- item 1 ----------------------------------------------------------
        if (sfx == "done" && un > 0) {
            r1++
            printf "%s: R1 RED - -done ticket still has %d line-start unchecked box(es): %s  %s\n", PROG, un, num, name
        }
        # -- item 2 ----------------------------------------------------------
        if (sfx == "done" && !doneflag) {
            r2++
            printf "%s: R2 RED - -done ticket with no done word on any of its %d status line(s): %s  %s\n", PROG, nstat, num, name
        }
        # -- item 3 ----------------------------------------------------------
        c = ""
        if (sfx == "live") {
            if (nstat == 0) {
                r3a++
                printf "%s: R3a RED - live ticket with no status line at all: %s  %s\n", PROG, num, name
            } else {
                c = canon(active)
                if (c == "") {
                    r3b++
                    printf "%s: R3b RED - status value outside the vocabulary: %s  value=\"%s\"  %s\n", PROG, num, active, name
                } else if (tolower(active) != c) {
                    r3c++
                    printf "%s: R3c NOTE - legacy spelling normalised to a vocabulary value: %s  \"%s\" -> \"%s\"  %s\n", PROG, num, active, c, name
                }
            }
        }
        # -- item 5 ----------------------------------------------------------
        if (length(name) > hat + 0) {
            r5++
            tag = (name in roster) ? "grandfathered by the 262 roster" : "NOT in the 262 roster - a new offence"
            printf "%s: R5 RED - ticket name longer than the hat (%d): %d %s  (%s)  %s\n", PROG, hat, length(name), unit, tag, name
            if (unit == "bytes" && name ~ /[^ -~]/) {
                printf "%s: R5 NOTE - this name carries non-ASCII bytes while the ruling unit is bytes, so its character length may be smaller: %s\n", PROG, name
            }
        }
        # -- item 6 ----------------------------------------------------------
        ac = agecut(c)
        if (ac > 0 && lastd != "-") {
            age = dayof(today) - dayof(lastd)
            if (age > ac) {
                r6++
                printf "%s: R6 YELLOW - %s older than %d day(s): %s  last activity %s (%d days)  %s\n", PROG, c, ac, num, lastd, age, name
            }
        }
        # -- item 7 ----------------------------------------------------------
        if (sfx == "live") {
            live++
            if (c == "blocked-on-owner") owner++
            if (c == "blocked") blocked++
            if (c == "ready" && !blk) readynow++
            if (c == "ready" && blk) readyblocked++
            if (blk) withblocker++
        }
    }
    END {
        # -- item 4 ----------------------------------------------------------
        for (k in dup) {
            r4dup++
            printf "%s: R4 YELLOW - two files share one ticket number: %s -> %s\n", PROG, k, dup[k]
        }
        for (i = minnum; i <= maxnum; i++) {
            if (i in numseen) continue
            pad = (i < 10) ? "0" i : "" i
            r4gap++
            hit = "-"
            for (v in voided) { split(v, p, "-"); if (p[1] == pad) hit = v }
            if (hit != "-") {
                printf "%s: R4 YELLOW - ticket number %s is absent from the pool and IS explained by voided/: %s\n", PROG, pad, hit
            } else {
                printf "%s: R4 YELLOW - ticket number %s is absent from the pool and voided/ does not explain it\n", PROG, pad
            }
        }
        printf "%s: R7 STAT done=%d live=%d void=%d blocked-on-owner=%d blocked=%d ready-right-now=%d ready-but-blocked=%d\n", \
            PROG, sfxcount["done"] + 0, live + 0, sfxcount["void"] + 0, owner + 0, blocked + 0, readynow + 0, readyblocked + 0
        printf "%s: R7 NOTE ready-right-now = live tickets whose status is ready and whose Blocked-by field is absent or empty; %d of %d live tickets carry a non-empty Blocked-by field, so that count is only as complete as the field is\n", \
            PROG, withblocker + 0, live + 0
        printf "%s: SUMMARY R1=%d R2=%d R3a=%d R3b=%d R3c=%d R4dup=%d R4gap=%d R5=%d R6=%d tickets=%d hat=%d unit=%s\n", \
            PROG, r1 + 0, r2 + 0, r3a + 0, r3b + 0, r3c + 0, r4dup + 0, r4gap + 0, r5 + 0, r6 + 0, tickets + 0, hat, unit
    }' "$1"
}

# ---------------------------------------------------------------------------
# the scan
# ---------------------------------------------------------------------------
run_scan() {
    if [ ! -d "$POOL" ]; then
        printf '%s: REFUSE - the pool directory %s does not exist\n' "$PROG" "$POOL" >&2
        return 2
    fi
    n_pool=0
    for f in "$POOL"/*.md; do
        [ -f "$f" ] && { n_pool=1; break; }
    done
    if [ "$n_pool" -eq 0 ]; then
        printf '%s: REFUSE - the pool %s holds no *.md ticket. A gate whose denominator is empty cannot report\n' "$PROG" "$POOL" >&2
        return 2
    fi
    WORK=$(mktemp -d 2>/dev/null)
    if [ -z "$WORK" ] || [ ! -d "$WORK" ]; then
        printf '%s: REFUSE - cannot create a working directory (mktemp returned "%s")\n' "$PROG" "$WORK" >&2
        return 2
    fi
    ALL=$WORK/all.tsv
    : > "$ALL"
    # the 262 gate roster, as data
    sed -n 's|^[[:space:]]*R\["\(\.scratch/wisp/issues/[^"]*\)"\].*|\1|p' "$LENGTH_GATE" \
        | awk '{ n = $0; sub(/.*\//, "", n); print "R\t" n }' >> "$ALL"
    # the withdrawn drafts, for the item-4 reconciliation
    if [ -d "$VOIDED" ]; then
        for v in "$VOIDED"/*.md; do
            [ -f "$v" ] || continue
            printf 'V\t%s\n' "${v##*/}" >> "$ALL"
        done
    fi
    census >> "$ALL"
    printf '%s: repo=%s\n' "$PROG" "$ROOT"
    printf '%s: pool=%s\n' "$PROG" "$POOL"
    printf '%s: work=%s (kept on disk; this project never deletes temp artifacts)\n' "$PROG" "$WORK"
    printf '%s: unit=%s (probe: length of the 5-byte sample "a U+25FF b" = %s)\n' "$PROG" "$UNITNAME" "$UNIT"
    printf '%s: hat=%d read from scripts/check-path-length-budget.sh (HAT_NAME_LIMIT, README rule 9 - this script carries no copy of it)\n' \
        "$PROG" "$HAT_NAME_LIMIT"
    printf '%s: today=%s (ageing is a calendar-day difference)\n' "$PROG" "$TODAY"
    report "$ALL" | awk '
        # Regroup the findings by rule in the plan order (1, 2, 3a, 3b, 3c, 4, 5,
        # 6, then the item-7 statistics), keeping the pool order inside each
        # group, so a human reads the report the way the plan is written.
        function rank(k) {
            if (k == "R1")  return 1
            if (k == "R2")  return 2
            if (k == "R3a") return 3
            if (k == "R3b") return 4
            if (k == "R3c") return 5
            if (k == "R4")  return 6
            if (k == "R5")  return 7
            if (k == "R6")  return 8
            if (k == "R7")  return 9
            return 10
        }
        {
            if (match($0, /: R[0-9][a-z]? /)) k = substr($0, RSTART + 2, RLENGTH - 3)
            else k = "ZZ"
            if (!(k in seen)) { seen[k] = 1; key[++n] = k }
            buf[k] = buf[k] $0 "\n"
        }
        END {
            for (i = 2; i <= n; i++) {
                t = key[i]; j = i - 1
                while (j >= 1 && rank(key[j]) > rank(t)) { key[j + 1] = key[j]; j-- }
                key[j + 1] = t
            }
            for (i = 1; i <= n; i++) {
                if (i > 1) printf "\n"
                printf "%s", buf[key[i]]
            }
        }'
    printf '%s: VERDICT reported, not enforced - v1 exits 0 by design; the findings above are for a human and this gate does not fail a build\n' "$PROG"
    return 0
}

# ---------------------------------------------------------------------------
# the positive control: a throwaway pool with one planted violation per item
# ---------------------------------------------------------------------------
# Nothing here touches the real pool: the bench is a fresh mktemp directory with
# its own git work tree, its own .scratch/wisp/issues, and a copy of this script
# plus a copy of the 262 gate (so the threshold is read through the same code
# path). Every planted violation must be named; every healthy sample must stay
# silent; and the scan must still exit 0 with violations on the table.
selftest_fail=0
selftest_n=0

st_ok() {
    selftest_n=$((selftest_n + 1))
    printf '%s: control %d ok - %s\n' "$PROG" "$selftest_n" "$1"
}
st_bad() {
    selftest_n=$((selftest_n + 1))
    selftest_fail=$((selftest_fail + 1))
    printf '%s: CONTROL FAILED - %s\n' "$PROG" "$1"
    printf '%s:   (bench kept at %s)\n' "$PROG" "$BENCH"
}
st_has() {  # label file pattern
    if grep -qF -- "$3" "$2"; then st_ok "$1"; else st_bad "$1: expected the output to carry [$3]"; fi
}
st_absent() {  # label file pattern
    if grep -qF -- "$3" "$2"; then st_bad "$1: the output carries [$3] but should not"; else st_ok "$1"; fi
}
st_reds() {  # label file expected-count
    got=$(grep -c -F ' RED - ' "$2")
    if [ "$got" -eq "$3" ]; then st_ok "$1"; else st_bad "$1: expected $3 RED line(s), got $got"; fi
}

# a ticket face with a given status block: $1 = path, $2 = status lines, $3 = rest
selftest_ticket() {
    {
        printf '# %s\n\n' "${1##*/}"
        printf '%s\n' "$2"
        printf '\n%s\n' "$3"
    } > "$1"
}

run_selftest() {
    BENCH=$(mktemp -d 2>/dev/null)
    if [ -z "$BENCH" ] || [ ! -d "$BENCH" ]; then
        printf '%s: REFUSE - the positive control cannot build a throwaway pool (mktemp returned "%s")\n' "$PROG" "$BENCH" >&2
        return 2
    fi
    printf '%s: positive control bench=%s (kept on disk; this project never deletes temp artifacts)\n' "$PROG" "$BENCH"
    if ! git init -q "$BENCH" 2>/dev/null; then
        printf '%s: REFUSE - git init failed in %s\n' "$PROG" "$BENCH" >&2
        return 2
    fi
    mkdir -p "$BENCH/.scratch/wisp/issues" "$BENCH/.scratch/wisp/voided" "$BENCH/scripts" || return 2
    cp "$LENGTH_GATE" "$BENCH/scripts/check-path-length-budget.sh" || return 2
    SELF_DIR=$(cd "$(dirname "$0")" && pwd)
    SELF=$SELF_DIR/$(basename "$0")
    cp "$SELF" "$BENCH/scripts/ticket-lint.sh" || return 2
    B=$BENCH/.scratch/wisp/issues

    # ---- healthy samples: every one of these must stay silent ---------------
    selftest_ticket "$B/01-clean-done-done.md" '**Status:** done' \
        '**Blocked by:** None

- [x] AC#1 checked

Last update: 2026-10-20'
    selftest_ticket "$B/02-appended-done-done.md" '**Status:** open
→ **ready-for-review**（appended line）
→ **accepted-done**（appended line: the done marker is NOT on the first line）' \
        '- [x] AC#1 checked

Last update: 2026-10-21'
    selftest_ticket "$B/03-clean-ready.md" '- Status: ready-for-agent' \
        '**Blocked by:** None

- [ ] AC#1 still open, and that is legal on a live ticket

Last update: 2026-10-25'
    selftest_ticket "$B/04-clean-review-recent.md" '- Status: review' \
        '**Blocked by:** 01-clean-done

Last update: 2026-10-25'

    rc=0
    ( cd "$BENCH" && sh scripts/ticket-lint.sh ) > "$BENCH/out-clean.txt" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ]; then st_ok "the bench baseline exits 0"
    else st_bad "the bench baseline should exit 0, got rc=$rc"; fi
    st_reds "the healthy pool has no RED line at all" "$BENCH/out-clean.txt" 0
    st_has "the healthy pool reports all-zero findings" "$BENCH/out-clean.txt" \
        'SUMMARY R1=0 R2=0 R3a=0 R3b=0 R3c=1 R4dup=0 R4gap=0 R5=0 R6=0'
    st_absent "the done ticket whose only done marker sits on an appended line is NOT reported by item 2" \
        "$BENCH/out-clean.txt" 'R2 RED'

    # ---- planted violations: one per item ----------------------------------
    selftest_ticket "$B/05-open-box-done.md" '- Status: done' \
        '- [ ] AC#2 never checked

Last update: 2026-10-20'
    selftest_ticket "$B/06-dup-left.md" '- Status: in-progress' 'Last update: 2026-10-25'
    selftest_ticket "$B/06-dup-right.md" '- Status: in-progress' 'Last update: 2026-10-25'
    selftest_ticket "$B/07-no-done-word-done.md" '**Status:** 待派' \
        '- [x] AC#1 checked

Last update: 2026-10-20'
    selftest_ticket "$B/08-no-status-live.md" 'No status line anywhere on this face.' 'Last update: 2026-10-25'
    selftest_ticket "$B/09-bad-value-live.md" '- Status: whatever-the-agent-felt-like（not a vocabulary value）' \
        'Last update: 2026-10-25'
    selftest_ticket "$B/10-aging-review.md" '- Status: review' 'Last update: 2018-01-01'
    selftest_ticket "$B/11-blocked-owner.md" '- Status: blocked-on-owner' 'Last update: 2018-01-01'
    selftest_ticket "$B/13-gap-anchor.md" '- Status: in-progress' 'Last update: 2026-10-25'
    selftest_ticket "$B/15-gap-anchor.md" '- Status: in-progress' 'Last update: 2026-10-25'
    selftest_ticket "$B/16-blocked.md" '- Status: blocked' 'Last update: 2026-10-25'
    printf 'withdrawn draft for number 12\n' > "$BENCH/.scratch/wisp/voided/12-first-draft-void.md"

    # a name over the hat, built from the hat itself so the sample cannot drift
    LONG=$(awk -v n="$HAT_NAME_LIMIT" 'BEGIN{ s = "17-"; while (length(s) < n + 2) s = s "x"; print s "-live.md" }')
    selftest_ticket "$B/$LONG" '- Status: blocked' 'Last update: 2026-10-25'

    rc=0
    ( cd "$BENCH" && sh scripts/ticket-lint.sh ) > "$BENCH/out-bad.txt" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ]; then st_ok "the scan still exits 0 with planted violations on the table (v1 reports, never blocks)"
    else st_bad "the scan must exit 0 even with findings, got rc=$rc"; fi
    st_has "item 1 names the -done ticket that still has a line-start unchecked box" "$BENCH/out-bad.txt" \
        'R1 RED - -done ticket still has 1 line-start unchecked box(es): 05'
    st_has "item 2 names the -done ticket with no done word on any status line" "$BENCH/out-bad.txt" \
        'R2 RED - -done ticket with no done word on any of its 1 status line(s): 07'
    st_absent "item 2 does not name 02-appended-done-done.md" "$BENCH/out-bad.txt" '02-appended-done'
    st_has "item 3a names the live ticket with no status line" "$BENCH/out-bad.txt" \
        'R3a RED - live ticket with no status line at all: 08'
    st_has "item 3b names the out-of-vocabulary value verbatim" "$BENCH/out-bad.txt" \
        'R3b RED - status value outside the vocabulary: 09  value="whatever-the-agent-felt-like"'
    st_has "item 3c reports the legacy spelling as a note, not a red" "$BENCH/out-bad.txt" \
        'R3c NOTE - legacy spelling normalised to a vocabulary value: 03  "ready-for-agent" -> "ready"'
    st_has "item 4 names the two files that share one ticket number" "$BENCH/out-bad.txt" \
        'R4 YELLOW - two files share one ticket number: 06'
    st_has "item 4 reconciles gap 12 against voided/" "$BENCH/out-bad.txt" \
        'ticket number 12 is absent from the pool and IS explained by voided/: 12-first-draft-void.md'
    st_has "item 4 reports gap 14 as unexplained" "$BENCH/out-bad.txt" \
        'ticket number 14 is absent from the pool and voided/ does not explain it'
    st_has "item 5 names the over-hat name and says it is not grandfathered" "$BENCH/out-bad.txt" \
        "R5 RED - ticket name longer than the hat ($HAT_NAME_LIMIT): ${#LONG} $UNITNAME  (NOT in the 262 roster - a new offence)"
    st_has "item 6 ages the review ticket" "$BENCH/out-bad.txt" \
        'R6 YELLOW - review older than 14 day(s): 10'
    st_has "item 6 ages the blocked-on-owner ticket on its own threshold" "$BENCH/out-bad.txt" \
        'R6 YELLOW - blocked-on-owner older than 21 day(s): 11'
    st_has "item 7 counts every bucket from the same records" "$BENCH/out-bad.txt" \
        'R7 STAT done=4 live=12 void=0 blocked-on-owner=1 blocked=2 ready-right-now=1 ready-but-blocked=0'

    # ---- teeth: make a healthy sample sick and watch the ruler notice ------
    sed '0,/accepted-done/ s/accepted-done/ready-for-review/' "$B/02-appended-done-done.md" > "$BENCH/t" 
    mv "$BENCH/t" "$B/02-appended-done-done.md"
    rc=0
    ( cd "$BENCH" && sh scripts/ticket-lint.sh ) > "$BENCH/out-mut1.txt" 2>&1
    rc=$?
    st_has "item 2 has teeth: removing the appended done word turns that ticket red" "$BENCH/out-mut1.txt" \
        'R2 RED - -done ticket with no done word on any of its 3 status line(s): 02'
    if [ "$rc" -eq 0 ]; then st_ok "the mutated scan still exits 0"
    else st_bad "the mutated scan must still exit 0, got rc=$rc"; fi

    selftest_ticket "$B/03-clean-ready.md" '- Status: today-i-feel-lucky' 'Last update: 2026-10-25'
    rc=0
    ( cd "$BENCH" && sh scripts/ticket-lint.sh ) > "$BENCH/out-mut2.txt" 2>&1
    rc=$?
    st_has "item 3b has teeth: a second bad value is caught and the note for the old one disappears" "$BENCH/out-mut2.txt" \
        'R3b RED - status value outside the vocabulary: 03  value="today-i-feel-lucky"'
    st_absent "item 3c notice is not stuck on: the normalised spelling is gone once the value is bad" \
        "$BENCH/out-mut2.txt" 'R3c NOTE'

    printf '%s: positive control %s - %d control(s) ran, %d failed\n' \
        "$PROG" "$([ "$selftest_fail" -eq 0 ] && echo PASSED || echo FAILED)" "$selftest_n" "$selftest_fail"
    [ "$selftest_fail" -eq 0 ] && return 0
    return 1
}

rc=0
if [ "$MODE" = selftest ] || [ "$MODE" = both ]; then
    if ! run_selftest; then
        printf '%s: positive control RED - the scan is not run, because a ruler that cannot go red does not get to print a reading\n' "$PROG" >&2
        exit 1
    fi
    if [ "$MODE" = selftest ]; then
        exit 0
    fi
fi

run_scan
SCAN_RC=$?
exit $SCAN_RC
