#!/usr/bin/env bash
# gate156-r4.sh <tag> - AC#7 per-package gate for ticket 156's FOURTH program (r4).
#
# Deliberately the SAME SHAPE as gate156-r3.sh (which is itself gate156.sh's shape),
# because AC#7 asks for an "after" reading that can be differenced against the
# landed "before" run at gate-r3-pre/: same package list, same flag set, same
# anchored patterns, same roster extraction. Anything that changed between the two
# scripts would be an instrument change, not a code change.
#
# Instrument facts this shape exists because of (all measured in this repo):
#   * cmd/wisp's test binary links the sherpa-onnx cgo import library. A missing or
#     drive-letter-form PATH fails at LOAD time (exit status 0xc0000135) with ZERO
#     '=== RUN' lines - "never ran", not "green". Authority:
#     scripts/wisp-cli-tests.sh:99-113 - the shell's own /d/work/... form is the one
#     the loader accepts, and `pwd -W` is the trap.
#   * red/green comes ONLY from anchored '^--- FAIL:'. This package prints the
#     string FAIL from inside PASSING cases: ticket 128's
#     TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 runs `wisp doctor` with
#     APPDATA unset on purpose, and the doctor's own report lines ('[FAIL] data dir
#     resolvable (dev)', '[FAIL] onnxruntime.dll version', '[FAIL] sherpa-onnx C
#     API', and its 'wisp doctor: FAIL' summary) land in that case's t.Logf output.
#     An unanchored grep for FAIL therefore reports 4 reds where there are 0.
#   * one panicking case can swallow dozens of sibling readings, so four numbers
#     alone cannot prove a run covered the same ground as another: the rosters and
#     their two-way comm are the check.
#   * no -race (this runner can die with 0xc0000374 / rc=1 and zero '--- FAIL',
#     which is neither red nor green), no -cover, no -overlay: plain shipped tree.
set -u
tag="$1"
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../.." && pwd)   # lives at .scratch/wisp/probes/156/
cd "$root" || exit 1
export PATH="$root/third_party/sherpa-onnx:$PATH"
out="$here/gate-r4-$tag"
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

printf '%-16s %-4s %-5s %-6s %-6s %-6s %-6s %-5s %s\n' \
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
    printf '%-16s %-4s %-5s %-6s %-6s %-6s %-6s %-5s %s\n' "$name" "$rc" "$load" "$run" "$pass" "$fail" "$skip" "$verdict"
    grep -E '^--- (PASS|FAIL|SKIP)' "$out/$name.log" | sed -E 's/^--- ([A-Z]+): (\S+).*/\1 \2/' | sort > "$out/$name-verdicts.txt"
    grep '^=== RUN' "$out/$name.log" | sed -E 's/^=== RUN +(\S+).*/\1/' | sort > "$out/$name-runs.txt"
    grep -E '^--- FAIL' "$out/$name.log" | sed -E 's/^--- FAIL: (\S+).*/\1/' | sort > "$out/$name-reds.txt"
    grep -E '^--- SKIP' "$out/$name.log" | sed -E 's/^--- SKIP: (\S+).*/\1/' | sort > "$out/$name-skips.txt"
    echo "rc=$rc pkg=$pkg log=$out/$name.log" >> "$out/meta.txt"
done
echo "end_iso=$(date -Is)" >> "$out/meta.txt"
