# 111-r5 AC#5 gate: d22scan on TWO pristine snapshots (HEAD^ = before my ci.yml change,
# HEAD = after), so "各 scope 不降" is measured against a real baseline instead of memory.
# Snapshots are built with `git archive | tar -x` under /tmp - never a worktree inside the repo.
set -u
log=.scratch/wisp/probes/111/r5/logs/d22scan-two-snapshots.txt
: > "$log"
before=/tmp/wisp-t111r5-before
after=/tmp/wisp-t111r5-after
mkdir -p "$before" "$after"

{
    echo "== date: $(date '+%Y-%m-%d %H:%M %z')"
    echo "== snapshot heads: before=$(git rev-parse --short HEAD^) after=$(git rev-parse --short HEAD)"
    git archive HEAD^ | tar -x -C "$before" && echo "archive HEAD^ rc=0"
    git archive HEAD | tar -x -C "$after" && echo "archive HEAD rc=0"
} >> "$log" 2>&1

( cd "$before" && sh scripts/d22scan.sh > d22scan-run.log 2>&1; echo "d22scan rc=$?" >> d22scan-run.log )
( cd "$after"  && sh scripts/d22scan.sh > d22scan-run.log 2>&1; echo "d22scan rc=$?" >> d22scan-run.log )

{
    echo
    echo "== rc of each run (the AC#5 requirement is rc=0 on a pristine snapshot) =="
    echo "before: $(grep -h 'd22scan rc=' "$before/d22scan-run.log" | tail -1)"
    echo "after:  $(grep -h 'd22scan rc=' "$after/d22scan-run.log" | tail -1)"
    echo
    echo "== per-scope lines, before vs after (any line here is a scope and its count) =="
    grep -hoE 'ban ?#[0-9]+.*' "$before/d22scan-run.log" | head -40 | sed 's/^/BEFORE  /'
    grep -hoE 'ban ?#[0-9]+.*' "$after/d22scan-run.log" | head -40 | sed 's/^/AFTER   /'
} >> "$log" 2>&1

grep -hoE 'ban ?#[0-9]+.*' "$before/d22scan-run.log" > /tmp/wisp-t111r5-scopes-before.txt 2>/dev/null
grep -hoE 'ban ?#[0-9]+.*' "$after/d22scan-run.log" > /tmp/wisp-t111r5-scopes-after.txt 2>/dev/null

{
    echo
    echo "== scope-line comparison =="
    echo "before_lines=$(wc -l < /tmp/wisp-t111r5-scopes-before.txt)"
    echo "after_lines=$(wc -l < /tmp/wisp-t111r5-scopes-after.txt)"
    if diff -u /tmp/wisp-t111r5-scopes-before.txt /tmp/wisp-t111r5-scopes-after.txt > /tmp/wisp-t111r5-scopes.diff; then
        echo "SCOPES_IDENTICAL=yes"
    else
        echo "SCOPES_DIFFER (diff in /tmp/wisp-t111r5-scopes.diff):"
        head -20 /tmp/wisp-t111r5-scopes.diff
    fi
} >> "$log" 2>&1

echo "log_lines=$(wc -l < "$log")"
