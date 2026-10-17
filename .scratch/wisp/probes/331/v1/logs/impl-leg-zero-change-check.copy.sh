#!/bin/sh
# 331-r1 AC#2 零改动尺：锚 blob ⇄ 候选文件，滤掉【整行注释】(^[[:space:]]*//) 后比字节。
# 射程（如实声明）：只测"非整行注释的字节有没有动"；整行注释的任何改动对本尺不可见。
# 用法: sh zero-change-check.sh <anchor> <candidate-file> [repo-dir]
ANCHOR="$1"; CAND="$2"; REPO="${3:-.}"
TMPD="$(mktemp -d "${TMPDIR:-/tmp}/zcc.XXXXXX")"
strip() { grep -Ev '^[[:space:]]*//' "$@"; }
( cd "$REPO" && git show "$ANCHOR:internal/observe/errors.go" ) | strip > "$TMPD/a.strip"
strip < "$CAND" > "$TMPD/b.strip"
SA=$(sha256sum < "$TMPD/a.strip" | cut -d' ' -f1)
SB=$(sha256sum < "$TMPD/b.strip" | cut -d' ' -f1)
echo "anchor-stripped-sha256=$SA"
echo "cand-stripped-sha256=$SB"
diff "$TMPD/a.strip" "$TMPD/b.strip" > "$TMPD/d.diff"
rc=$?
if [ -s "$TMPD/d.diff" ]; then nonempty=yes; else nonempty=no; fi
echo "diff-nonempty=$nonempty diff-rc=$rc"
echo "TMPD=$TMPD"
exit $rc
