#!/usr/bin/env bash
# Ticket 70 control helper: extract the two ticket-314 END ANCHOR blocks BY MARKER
# (not by line number, so the "before" tree and the fixed tree can be compared
# even though the fix shifted every line below the roster). Prints the block's
# line range in the file it read, then the sha256 of the block's bytes.
set -u
repo_root=$1
before=$2
after=$3

block_range() {
    f=$1
    start_marker=$2
    end_prefix=$3
    s=$(grep -nF -- "$start_marker" "$f" | head -1 | cut -d: -f1)
    if [ -z "$s" ]; then
        echo "MISSING START: $start_marker in $f" >&2
        return 1
    fi
    e=$(awk -v n="$s" -v m="$end_prefix" 'NR > n && index($0, m) == 1 { print NR - 1; exit }' "$f")
    if [ -z "$e" ]; then
        echo "MISSING END: $end_prefix in $f" >&2
        return 1
    fi
    printf '%s %s\n' "$s" "$e"
}

hash_block() {
    f=$1
    s=$2
    e=$3
    sed -n "${s},${e}p" "$f" | sha256sum | cut -d" " -f1
}

m1="# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR, computed BEFORE the run."
m2="# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR verdict ----------------"

# block 1's end prefix is the capture line that follows it; block 2's is the
# cleanup line that follows it.
for pair in "before:$before" "after:$after"; do
    label=${pair%%:*}
    f=${pair#*:}
    r=$(block_range "$f" "$m1" 'capture=$(mktemp')
    set -- $r
    printf '%s block1 %s..%s sha256=%s\n' "$label" "$1" "$2" "$(hash_block "$f" "$1" "$2")"
    r=$(block_range "$f" "$m2" 'rm -f "$capture" "$resolved"')
    set -- $r
    printf '%s block2 %s..%s sha256=%s\n' "$label" "$1" "$2" "$(hash_block "$f" "$1" "$2")"
done

# the repo root is only used to make the caller's cwd explicit; nothing here
# writes or checks anything against it.
: "$repo_root"
