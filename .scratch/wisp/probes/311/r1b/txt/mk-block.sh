#!/usr/bin/env bash
# mk-block.sh - 311-r1b replication instrument for ticket 311 AC#1 shape B.
# The guard block is COPIED, never re-typed: it is extracted from a source file by
# content anchors and inserted into a target file before a content-anchored line,
# with exactly one substitution (the entry's own echo prefix). That is what makes
# "four entries, one shape" checkable by a line diff instead of by reading.
#
# usage:
#   sh mk-block.sh extract SRCFILE OUTBLOCK
#   sh mk-block.sh insert  SRCFILE TARGETFILE PREFIX ANCHOR_REGEX
#   sh mk-block.sh strip   TARGETFILE            (positive control: block removed only)
#   sh mk-block.sh has     TARGETFILE            (1 = block present, 0 = absent)
set -e
# NOTE: the block header contains a literal "(" - as an awk *dynamic* regex that is
# an unmatched group, so it must arrive escaped. The first two trips of this
# instrument hit exactly that (first: extract produced 0 lines and insert added
# nothing but a blank line - reverted by content, see logs/13; second: the -v layer
# itself un-escaped the backslash). So the markers below are matched as FIXED
# STRINGS with index()/equality, not as regexes: nothing left to escape.
# Strings, not regexes (see NOTE above). WISP_MARK_START lets the same instrument cut
# out the ticket-310 block as well, so "four entries, one shape" is measurable for
# both blocks with ONE ruler (dispatch section 9 asks exactly that question).
START_MARK=${WISP_MARK_START:-'# ---- TICKET 311 AC#1 (shape B'}
END_MARK='fi'
mode=$1

case $mode in
has)
    awk -v sm="$START_MARK" 'index($0, sm) == 1 { found = 1 } END { print (found ? 1 : 0) }' "$2"
    ;;
extract)
    src=$2
    out=$3
    # The block is: the "TICKET 311 AC#1 (shape B" comment header through the first
    # column-0 `fi` after it (every fi inside the block is indented, so that one is
    # the block's own end), with trailing blank lines dropped.
    awk -v sm="$START_MARK" -v em="$END_MARK" '
        !inb && index($0, sm) == 1 { inb = 1 }
        inb { print }
        inb && $0 == em { inb = 0 }
    ' "$src" | awk '{ buf[NR] = $0 } END { n = NR; while (n > 0 && buf[n] == "") n--; for (i = 1; i <= n; i++) print buf[i] }' > "$out"
    echo "extract_lines=$(wc -l < "$out")"
    ;;
insert)
    src=$2
    tgt=$3
    prefix=$4
    anchor=$5
    tmp=$(mktemp)
    blk=$(mktemp)
    sh "$0" extract "$src" "$blk" > /dev/null
    blk_ct=$(wc -l < "$blk")
    if [ "$blk_ct" -lt 40 ]; then
        echo "REFUSED: extracted block is only $blk_ct lines from $src - refusing to insert an empty/partial block"
        exit 97
    fi
    if [ "$(sh "$0" has "$tgt")" = "1" ]; then
        echo "REFUSED: $tgt already contains the ticket 311 block (double insert would make two guards)"
        exit 98
    fi
    if ! grep -qE "$anchor" "$tgt"; then
        echo "REFUSED: content anchor [$anchor] not found in $tgt - placement unverified"
        exit 99
    fi
    sed "s|winsec-tests.sh:|${prefix}:|g" "$blk" > "$blk.p"
    blk_lines=$(wc -l < "$blk.p")
    prefix_lines=$(grep -c "${prefix}: ticket 311" "$blk.p" || true)
    awk -v anchor="$anchor" -v blkfile="$blk.p" '
        $0 ~ anchor && !ins {
            while ((getline line < blkfile) > 0) print line
            print ""
            ins = 1
        }
        { print }
        END { if (!ins) { print "ERROR-BLOCK-NOT-INSERTED" > "/dev/stderr"; exit 3 } }
    ' "$tgt" > "$tmp"
    mv "$tmp" "$tgt"
    rm -f "$blk" "$blk.p"
    echo "insert_target=$tgt entry_prefix=$prefix anchor=$anchor block_lines=$blk_lines prefixed_echo_lines=$prefix_lines"
    ;;
strip)
    tgt=$2
    tmp=$(mktemp)
    awk -v sm="$START_MARK" -v em="$END_MARK" '
        !skipping && index($0, sm) == 1 { skipping = 1; removed++; next }
        skipping {
            removed++
            if ($0 == em) { skipping = 0; eatblank = 1 }
            next
        }
        eatblank == 1 && $0 == "" { eatblank = 0; next }
        { print }
        END { printf("strip_removed_lines=%d\n", removed) > "/dev/stderr" }
    ' "$tgt" > "$tmp" 2>/tmp/mk-block-strip.txt
    if ! grep -q 'strip_removed_lines=[1-9]' /tmp/mk-block-strip.txt; then
        echo "REFUSED: strip found no block in $tgt - a silent no-op would fake a positive control"
        rm -f "$tmp"
        exit 94
    fi
    mv "$tmp" "$tgt"
    cat /tmp/mk-block-strip.txt
    ;;
*)
    echo "unknown mode: $mode"
    exit 2
    ;;
esac
