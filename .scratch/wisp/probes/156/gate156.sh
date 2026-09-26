#!/usr/bin/env bash
# gate156.sh <tag> - one package-by-package gate run for ticket 156.
#
# AC#7 of the ticket says: gates run ONE PACKAGE AT A TIME (never a whole-repo
# sweep - that is the orchestrator's move after a batch lands), and a run counts
# only if the === RUN lines exist. Two instrument facts this script exists
# because of, both measured in this repo and repeated in the dispatch:
#   * cmd/wisp's test binary links the sherpa-onnx cgo import library, so a
#     missing/oddly-spelled PATH fails at LOAD time (0xc0000135) with 0 === RUN
#     lines: that is "never ran", not "green" (scripts/wisp-cli-tests.sh).
#   * PATH must be the shell's own path form (/d/work/...), NOT `pwd -W`'s
#     drive-letter form, which the MSYS rewrite mangles into the 0xc0000135 shape.
set -u
tag="$1"
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
# this script lives at .scratch/wisp/probes/156/, four levels under the root
root=$(CDPATH= cd -- "$here/../../../.." && pwd)
cd "$root" || exit 1
export PATH="$root/third_party/sherpa-onnx:$PATH"
out="$here/gate-$tag"
mkdir -p "$out"
echo "gate156: tag=$tag pwd=$(pwd) dll_dir=$root/third_party/sherpa-onnx"
for pkg in ./cmd/wisp/ ./internal/tools/; do
    name=$(echo "$pkg" | tr -d './' )
    echo "=== $pkg ==="
    go test -count=1 -v "$pkg" >"$out/$name.log" 2>&1
    rc=$?
    printf '%-16s rc=%-3s RUN=%-4s PASS=%-4s FAIL=%-4s SKIP=%-4s\n' \
        "$name" "$rc" \
        "$(grep -c '^=== RUN' "$out/$name.log")" \
        "$(grep -c '^--- PASS' "$out/$name.log")" \
        "$(grep -c '^--- FAIL' "$out/$name.log")" \
        "$(grep -c '^--- SKIP' "$out/$name.log")"
done
