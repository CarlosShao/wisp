#!/bin/sh
# 212-r3: comment-side phantom-citation survey over tools/** (APPROXIMATION, not the
# instrument's own ruler - ban #9 does not scan tools/ today, and this helper cannot
# reproduce repoPathRe's tokenisation either; it is the粗分母 for the owner's pending
# "widening to tools/**" decision, per 212-v2 §8 #4).
#
# Usage: sh tools-cited.sh <dir-with-.go-files>   -> prints tokens cited in comment
# lines whose spelled-out form does not exist at the repo root.
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../../.." && pwd)
dir=${1:-tools}

cd "$root"
grep -rhE '^[[:space:]]*//' "$dir" --include=*.go \
    | grep -oE '(docs|\.scratch|internal|cmd|tools|scripts)/[A-Za-z0-9_./-]*[A-Za-z0-9_-]' \
    | sort -u \
    | while read -r tok; do
        case "$tok" in
            *'...'*|*'…'*|*'*'*) continue ;;   # class 3, exempt by the region rule
        esac
        if [ ! -e "$tok" ]; then
            printf '%s\n' "$tok"
        fi
      done
