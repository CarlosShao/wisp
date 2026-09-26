#!/usr/bin/env bash
# gate156-r3.sh <tag> - AC#7 per-package gate for ticket 156's THIRD program (r3).
#
# Deliberately the same SHAPE as gate156.sh (r1's instrument) so an r3 reading and
# an r1 reading are about the same thing: one `go test -count=1 -v` per package,
# never a whole-repo sweep, four numbers counted from anchored patterns, plus the
# two rosters AC#7 asks for and their two-way `comm`.
#
# Instrument facts this script exists because of (all measured in this repo, none
# invented here):
#   * cmd/wisp's test binary links the sherpa-onnx cgo import library. A missing
#     or drive-letter-form PATH fails at LOAD time (exit status 0xc0000135) with
#     ZERO '=== RUN' lines - that is "never ran", not "green"
#     (scripts/wisp-cli-tests.sh:101-109 is the authority, and it says the
#     shell's own /d/work/... form is the one the loader accepts).
#   * a single panicking case can swallow dozens of sibling readings, so the four
#     numbers alone cannot prove a run covered the same ground: the rosters and
#     their差集 are the check (this ticket's dispatch, gate section).
#   * red/green is read ONLY from anchored '^--- FAIL:'. t.Logf lines also carry a
#     'file:line:' prefix, so an unanchored 'FAIL' grep can be a log line.
#   * no -race (this runner can die with 0xc0000374 / rc=1 and zero '--- FAIL',
#     which is neither red nor green), no -cover, no -overlay: this is the plain
#     shipped tree.
set -u
tag="$1"
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../.." && pwd)   # lives at .scratch/wisp/probes/156/
cd "$root" || exit 1
export PATH="$root/third_party/sherpa-onnx:$PATH"
out="$here/gate-r3-$tag"
mkdir -p "$out"

{
    echo "tag=$tag"
    echo "anchor_head=$(git rev-parse HEAD)"
    echo "anchor_short=$(git rev-parse --short HEAD)"
    echo "branch=$(git branch --show-current)"
    echo "start_iso=$(date -Is)"
    echo "go_version=$(go version)"
    echo "pwd=$root"
    echo "dll_dir=/d/work/... shell form: $root/third_party/sherpa-onnx -> ${PATH%%:*}"
    echo "status_cmd_wisp_internal=[$(git status --porcelain -- cmd/wisp internal)]"
    for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go \
             cmd/wisp/slo_exit_os_156_windows_test.go internal/tools; do
        echo "tree_sha256 ${f}=$(git hash-object "$f" | cut -c1-16)"
    done
    echo "runner_listener_procs=$(tasklist 2>/dev/null | grep -c 'Runner.Listener.exe')"
    echo "runner_worker_procs=$(tasklist 2>/dev/null | grep -c 'Runner.Worker.exe')"
} > "$out/meta.txt" 2>&1

printf '%-16s %-4s %-5s %-6s %-6s %-6s %-5s %s\n' \
    pkg rc 0xc0000135 RUN PASS FAIL SKIP verdict
for pkg in ./cmd/wisp/ ./internal/tools/; do
    name=$(echo "$pkg" | tr -d './')
    go test -count=1 -v "$pkg" > "$out/$name.log" 2>&1
    rc=$?
    run=$(grep -c '^=== RUN' "$out/$name.log")
    pass=$(grep -c '^--- PASS' "$out/$name.log")
    fail=$(grep -c '^--- FAIL' "$out/$name.log")
    skip=$(grep -c '^--- SKIP' "$out/$name.log")
    load=$(grep -c '0xc0000135' "$out/$name.log")
    verdict=$(grep -E '^(ok|FAIL)\s' "$out/$name.log" | tail -1 | tr -s ' \t' ' ')
    printf '%-16s %-4s %-5s %-6s %-6s %-6s %-5s %s\n' "$name" "$rc" "$load" "$run" "$pass" "$fail" "$skip" "$verdict"
    # rosters: top-level verdict names, and all RUN names (subtest RUN lines are
    # NOT indented by this Go version - measured, gate-pre/internaltools.log).
    grep -E '^--- (PASS|FAIL|SKIP)' "$out/$name.log" | sed -E 's/^--- ([A-Z]+): (\S+).*/\1 \2/' | sort > "$out/$name-verdicts.txt"
    grep '^=== RUN' "$out/$name.log" | sed -E 's/^=== RUN +(\S+).*/\1/' | sort > "$out/$name-runs.txt"
    grep -E '^--- FAIL' "$out/$name.log" | sed -E 's/^--- FAIL: (\S+).*/\1/' | sort > "$out/$name-reds.txt"
    grep -E '^--- SKIP' "$out/$name.log" | sed -E 's/^--- SKIP: (\S+).*/\1/' | sort > "$out/$name-skips.txt"
    echo "rc=$rc pkg=$pkg log=$out/$name.log"
done
echo "end_iso=$(date -Is)" >> "$out/meta.txt"
