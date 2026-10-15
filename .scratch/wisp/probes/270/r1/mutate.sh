#!/bin/sh
# mutate.sh - ticket 270 AC#1's mutation harness (form 甲).
#
# Why physical rewrite and not `go test -overlay`: overlay replaces only the bytes
# the COMPILER sees. 321-r1 measured that every ruler which reads its subject at
# RUNTIME (os.ReadFile / parser.ParseFile) is blind to it, so an overlay mutant of
# a source shape reads as a false green. This leg mutates the file on disk and
# puts it back by sha256; the .bak is kept, never deleted (temporary files are
# created, not removed).
#
#   sh mutate.sh save       - baseline sha256 + .bak of the subject
#   sh mutate.sh mu-fix     - delete the assembly guard branch (post-甲 tree)
#   sh mutate.sh mu-nofix   - restore HEAD's production file, THEN delete the
#                             same branch (the pre-甲 tree, i.e. no call-site check)
#   sh mutate.sh restore    - put .bak back and verify sha256 equals baseline
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../../.." && pwd)
subj="$root/internal/tools/subagent_197.go"
bak="$here/subagent_197.go.bak"
headfile="$here/subagent_197.go.head"
shafile="$here/sha256.txt"

if [ ! -f "$root/go.mod" ]; then
    echo "mutate.sh: derived repo root has no go.mod: $root" >&2
    exit 2
fi

record() {
    sha256sum "$subj" | sed 's# .*##' > "$shafile"
    echo "subject sha256: $(cat "$shafile")"
}

guard_strip() {
    python "$here/guard_strip.py" "$subj"
    echo "guard branch removed from: $subj"
}

case "${1:-}" in
    save)
        cp "$subj" "$bak"
        echo "backup written: $bak"
        record
        ;;
    mu-fix)
        cp "$subj" "$bak"
        record
        guard_strip
        sha256sum "$subj" | sed 's# .*##'
        ;;
    mu-nofix)
        cp "$subj" "$bak"
        record
        git -C "$root" show HEAD:internal/tools/subagent_197.go > "$headfile"
        cp "$headfile" "$subj"
        guard_strip
        sha256sum "$subj" | sed 's# .*##'
        ;;
    restore)
        cp "$bak" "$subj"
        echo "restored from: $bak"
        now=$(sha256sum "$subj" | sed 's# .*##')
        want=$(cat "$shafile")
        if [ "$now" != "$want" ]; then
            echo "mutate.sh: restore did not land on the baseline sha - now=$now want=$want" >&2
            exit 3
        fi
        echo "restore verified: $now"
        ;;
    *)
        echo "usage: sh mutate.sh {save|mu-fix|mu-nofix|restore}" >&2
        exit 2
        ;;
esac
