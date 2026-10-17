#!/bin/sh
# 335-r1 AC#2 zero-change ruler, same shape as ticket 331's zero-change-check.sh:
# anchor blob vs candidate file, both filtered of WHOLE-LINE comments (^[[:space:]]*//),
# then compared byte for byte.
# Declared range (honestly): this measures only "did any non-whole-line-comment byte
# move". Any change confined to whole-line comments is invisible to it - which is
# exactly the point of AC#2, and why the run ships with a positive control (one
# non-comment byte moved: the ruler must go red) and a negative control (only a
# whole-line comment moved: the ruler stays green, i.e. it is blind by construction).
#   usage: sh zc335r1.sh <anchor> <candidate-file> [repo-dir]
ANCHOR="$1"; CAND="$2"; REPO="${3:-.}"
FILE=internal/panel/display_scrub_216.go
TMPD="$(mktemp -d "${TMPDIR:-/tmp}/zc335.XXXXXX")"
strip() { grep -Ev '^[[:space:]]*//' "$@"; }
( cd "$REPO" && git show "$ANCHOR:$FILE" ) | strip > "$TMPD/a.strip"
strip < "$CAND" > "$TMPD/b.strip"
SA=$(sha256sum < "$TMPD/a.strip" | cut -d' ' -f1)
SB=$(sha256sum < "$TMPD/b.strip" | cut -d' ' -f1)
echo "anchor=$ANCHOR file=$FILE candidate=$CAND"
echo "anchor-stripped-sha256=$SA"
echo "cand-stripped-sha256=$SB"
diff "$TMPD/a.strip" "$TMPD/b.strip" > "$TMPD/d.diff"
rc=$?
if [ -s "$TMPD/d.diff" ]; then nonempty=yes; else nonempty=no; fi
echo "same_sha=$([ "$SA" = "$SB" ] && echo yes || echo no) diff-nonempty=$nonempty diff-rc=$rc"
echo "TMPD=$TMPD"
exit $rc
