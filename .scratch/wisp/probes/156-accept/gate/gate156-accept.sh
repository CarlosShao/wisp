#!/usr/bin/env bash
# accept-r1's OWN AC#7 gate. Same shape as gate156-r4.sh / gate156-r3.sh / gate156.sh
# (same package list, same flags, same anchored patterns, same roster extraction),
# because an acceptance run that changed the instrument could not be differenced
# against the landed ones. Output goes to probes/156-accept/gate/accept-post/.
set -u
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../../.." && pwd)   # lives at .scratch/wisp/probes/156-accept/gate/
cd "$root" || exit 1
export PATH="$root/third_party/sherpa-onnx:$PATH"
out="$here/accept-post"
mkdir -p "$out"

{
    echo "tag=accept-r1"
    echo "requested_anchor=7ca1130"
    echo "worktree_head=$(git rev-parse HEAD)"
    echo "branch=$(git branch --show-current)"
    echo "start_iso=$(date -Is)"
    echo "go_version=$(go version)"
    echo "pwd=$root"
    echo "dll_dir shell form: $root/third_party/sherpa-onnx -> ${PATH%%:*}"
    echo "status_cmd_wisp_internal=[$(git status --porcelain -- cmd/wisp internal)]"
    echo "diff_against_anchor_7ca1130=[$(git diff --name-only 7ca1130 -- cmd/wisp internal/tools go.mod go.sum)]"
    for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go \
             cmd/wisp/slo_exit_os_156_windows_test.go; do
        echo "tree_sha256 ${f}=$(git hash-object "$f" | cut -c1-16)"
        echo "anchor_sha256 ${f}=$(git rev-parse 7ca1130:$f | cut -c1-16)"
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
    run=$(grep -a -c '^=== RUN' "$out/$name.log")
    pass=$(grep -a -c '^--- PASS' "$out/$name.log")
    fail=$(grep -a -c '^--- FAIL' "$out/$name.log")
    skip=$(grep -a -c '^--- SKIP' "$out/$name.log")
    load=$(grep -a -c '0xc0000135' "$out/$name.log")
    verdict=$(grep -aE '^(ok|FAIL)[ \t]' "$out/$name.log" | tail -1 | tr -s ' \t' ' ')
    printf '%-16s %-4s %-5s %-6s %-6s %-6s %-6s %-5s %s\n' "$name" "$rc" "$load" "$run" "$pass" "$fail" "$skip" "$verdict"
    grep -aE '^--- (PASS|FAIL|SKIP)' "$out/$name.log" | sed -E 's/^--- ([A-Z]+): (\S+).*/\1 \2/' | sort > "$out/$name-verdicts.txt"
    grep -a '^=== RUN' "$out/$name.log" | sed -E 's/^=== RUN +(\S+).*/\1/' | sort > "$out/$name-runs.txt"
    grep -aE '^--- FAIL' "$out/$name.log" | sed -E 's/^--- FAIL: (\S+).*/\1/' | sort > "$out/$name-reds.txt"
    grep -aE '^--- SKIP' "$out/$name.log" | sed -E 's/^--- SKIP: (\S+).*/\1/' | sort > "$out/$name-skips.txt"
    echo "rc=$rc pkg=$pkg log=$out/$name.log" >> "$out/meta.txt"
done
echo "end_iso=$(date -Is)" >> "$out/meta.txt"
