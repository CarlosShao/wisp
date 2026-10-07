#!/usr/bin/env bash
# 275-a1 AC#0: per-file gofumpt exit-code roster over `git ls-files '*.go'` (935 files).
# Read-only w.r.t. the repo: only writes under .scratch/wisp/probes/275/a1/logs/.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
GF="D:/work/base/gopath/bin/gofumpt.exe"
LOGD=".scratch/wisp/probes/275/a1/logs"

git ls-files -z '*.go' > "$LOGD/tracked-go.zlist"
wc -c < "$LOGD/tracked-go.zlist" > "$LOGD/tracked-go.zlist.bytes"
git ls-files '*.go' | wc -l > "$LOGD/tracked-go.count"

: > "$LOGD/gofumpt-exitcodes.tsv"
: > "$LOGD/gofumpt-nonzero-detail.txt"

echo "start $(date)" > "$LOGD/ac0.progress"
i=0
while IFS= read -r -d '' f; do
    i=$((i + 1))
    out=$("$GF" -l -- "$f" 2>&1)
    rc=$?
    printf '%s\t%s\n' "$f" "$rc" >> "$LOGD/gofumpt-exitcodes.tsv"
    if [ "$rc" -ne 0 ]; then
        {
            printf '=== file: %s\n=== rc: %s\n=== gofumpt -l combined output (verbatim):\n' "$f" "$rc"
            printf '%s\n' "$out"
            printf '=== first error line: %s\n\n' "$(printf '%s\n' "$out" | sed -n '1p')"
        } >> "$LOGD/gofumpt-nonzero-detail.txt"
    fi
    if [ $((i % 100)) -eq 0 ]; then echo "i=$i $(date)" > "$LOGD/ac0.progress"; fi
done < "$LOGD/tracked-go.zlist"

echo "loop done i=$i $(date)" > "$LOGD/ac0.progress"

# batch form (the shape attrib.sh:341 actually runs) - exit code taken WITHOUT a pipe
git ls-files -z '*.go' | xargs -0 "$GF" -l > "$LOGD/batch-gofumpt-l-tracked.txt" 2> "$LOGD/batch-gofumpt-l-stderr.txt"
echo "batch_xargs_rc=$?" > "$LOGD/batch-rc.txt"
"$GF" -l -- .scratch/wisp/probes/185/c1/mut/fs_broken.go > "$LOGD/single-broken-file-stdout.txt" 2> "$LOGD/single-broken-file-stderr.txt"
echo "single_broken_rc=$?" >> "$LOGD/batch-rc.txt"
wc -l "$LOGD/gofumpt-exitcodes.tsv" >> "$LOGD/batch-rc.txt"
echo "ALL DONE $(date)" >> "$LOGD/ac0.progress"
