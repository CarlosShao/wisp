#!/bin/sh
# drive.sh - ticket 136 leg 136s probe driver for scripts/slo-freshness.sh.
#
# Drives the freshness pin through its OWN seams (SLO_FRESH_NOW,
# SLO_FULL_LAST_TRIGGER, SLO_FULL_LAST_SAMPLE, SLO_FULL_ARTIFACT_PAGES) plus a
# `gh` stub placed first on PATH - the same instrument the ticket-134
# acceptance used for its K-shapes. No network, no token, no real collection:
# nothing here runs `wisp slo` or scripts/slo-check.ps1.
#
# usage: sh drive.sh <slo-freshness.sh path> <out file> [world ...]
#   with no world names all worlds run. One call = one tree: run it once on the
#   HEAD snapshot (tree-old) and once on the fixed copy (tree-new).
#   Fixtures/stub/logs live under ${WS:-/d/tmp/wisp136s} (created, never deleted).
set -u

dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
ws=${WS:-/d/tmp/wisp136s}
script=$1
outfile=$2
shift 2
only="$*"

tab=$(printf '\t')
nowv='2026-09-23T09:00:00Z'
trig_fresh='2026-09-23T08:30:00Z|completed|success|https://example.invalid/run/1'
trig_future='2027-01-01T00:00:00Z|completed|success|https://example.invalid/run/2'
samp_fresh="2026-09-23T08:45:00Z${tab}35831084009${tab}10736979707"
samp_stale="2026-09-13T08:45:00Z${tab}35831084009${tab}10736979707"
samp_future="2027-01-01T00:00:00Z${tab}35831084009${tab}10736979707"
all_worlds='control spread cap-zero forge-branch forge-workflow forge-event missing-ts missing-source partial-green partial-notfound future-sample future-job stale-sample k1-no-token k2-no-repo k3-runs-fail k4-artifacts-fail k4b-total-fail k5-empty k6-empty-stale k7-run-fail'

stubdir=$ws/bin; mkdir -p "$stubdir"
cat > "$stubdir/gh" <<'STUB'
#!/bin/sh
# gh stub for the ticket-136s probes: serves only the endpoints
# scripts/slo-freshness.sh uses, out of $STUB_WORLD, and applies the caller's
# --jq filter with real jq when present, else ghjq.py. Unknown endpoint or
# filter => exit 1, so the pin lands in "cannot look" instead of a fake answer.
endpoint=''
filter=''
while [ $# -gt 0 ]; do
    case $1 in
    api) shift ;;
    --jq) filter=$2; shift 2 ;;
    *) endpoint=$1; shift ;;
    esac
done
fail=''
for cls in ${STUB_FAIL:-}; do
    case $cls in
    runs-list) case $endpoint in */actions/workflows/*/runs*) fail=$cls ;; esac ;;
    artifacts) case $endpoint in */actions/artifacts*) fail=$cls ;; esac ;;
    total) case $endpoint in */actions/artifacts*per_page=1\&*) fail=$cls ;; esac ;;
    run) case $endpoint in */actions/runs/*) case $endpoint in */jobs*) ;; *) fail=$cls ;; esac ;; esac ;;
    esac
done
if [ -n "$fail" ]; then
    echo "stub gh: injected failure '$fail' for $endpoint" >&2
    exit 1
fi
case $endpoint in
*/actions/artifacts*)
    case $endpoint in
    *per_page=1\&*) file=$STUB_WORLD/total.json ;;
    *page=*) page=${endpoint##*page=}; page=${page%%&*}; file=$STUB_WORLD/artifacts-page$page.json ;;
    *) echo "stub gh: unhandled artifacts query: $endpoint" >&2; exit 1 ;;
    esac ;;
*/actions/runs/*/jobs*) file=$STUB_WORLD/jobs.json ;;
*/actions/runs/*) rid=${endpoint##*/actions/runs/}; rid=${rid%%/*}; file=$STUB_WORLD/run-$rid.json ;;
*/actions/workflows/*/runs*) file=$STUB_WORLD/runs.json ;;
*) echo "stub gh: unhandled endpoint: $endpoint" >&2; exit 1 ;;
esac
[ -f "$file" ] || { echo "stub gh: no fixture $file for $endpoint" >&2; exit 1; }
if [ -z "$filter" ]; then
    cat "$file"
elif command -v jq >/dev/null 2>&1; then
    jq -r "$filter" "$file"
elif [ -n "${STUB_PY:-}" ]; then
    "$STUB_PY" "$STUB_JQ" "$file" "$filter"
else
    echo 'stub gh: no jq on PATH and no STUB_PY - cannot apply --jq' >&2
    exit 1
fi
STUB
chmod +x "$stubdir/gh"
path_in="$stubdir:/usr/bin:/bin"
pick_py() {
    # A working interpreter, not just a name on PATH (Windows ships a Store
    # alias for python3 that exits 9009 and prints "Python was not found").
    for cand in python3 python py; do
        p=$(command -v "$cand" 2>/dev/null) || continue
        "$p" -c 'pass' >/dev/null 2>&1 && { printf '%s' "$p"; return 0; }
    done
    return 1
}
stub_py=$(pick_py || true)

emit_page() { # $1 total, $2 artifact specs "id,created,run,branch ; ..."
    printf '{"total_count":%s,"artifacts":[' "$1"
    sep=''
    printf '%s\n' "$2" | tr ';' '\n' | while IFS=, read -r id created run branch; do
        [ -n "$id" ] || continue
        wr=''
        [ -z "$run" ] || wr=',"workflow_run":{"id":'"$run"',"head_branch":"'"$branch"'"}'
        ts=''
        [ -z "$created" ] || ts=',"created_at":"'"$created"'"'
        printf '%s{"id":%s,"name":"slo-full-report","expired":false%s%s}' "$sep" "$id" "$ts" "$wr"
        sep=','
    done
    printf ']}'
}

setup_world() { # sets W_* and builds $wdir fixtures
    w=$1
    W_TOTAL=188; W_ART1=''; W_ART2=''; W_ART3=''; W_RUNS=''
    W_FAIL=''; W_PAGES=3; W_TRIG=$trig_fresh; W_SAMP=''
    W_REPO=stub/repo; W_TOKEN=stub-token; W_NOW=$nowv
    case $w in
    control) W_ART1="80,2026-09-23T08:45:00Z,700,dev;81,2026-09-23T04:00:00Z,701,main"
             W_RUNS="700,.github/workflows/ci.yml,push,dev;701,.github/workflows/ci.yml,schedule,main" ;;
    spread) W_ART1="82,2026-09-01T04:00:00Z,800,dev"
             W_ART2="83,2026-09-23T08:45:00Z,801,main"
             W_RUNS="800,.github/workflows/ci.yml,push,dev;801,.github/workflows/ci.yml,schedule,main" ;;
    cap-zero) W_PAGES=0; W_ART1="84,2026-09-23T08:45:00Z,802,dev"
             W_RUNS="802,.github/workflows/ci.yml,push,dev" ;;
    forge-branch) W_ART1="90,2026-09-23T08:45:00Z,701,feature/evil"
             W_RUNS="701,.github/workflows/ci.yml,push,dev" ;;
    forge-workflow) W_ART1="91,2026-09-23T08:45:00Z,702,dev"
             W_RUNS="702,.github/workflows/chore.yml,push,dev" ;;
    forge-event) W_ART1="92,2026-09-23T08:45:00Z,703,dev"
             W_RUNS="703,.github/workflows/ci.yml,pull_request,dev" ;;
    missing-ts) W_ART1="93,,704,dev"
             W_RUNS="704,.github/workflows/ci.yml,push,dev" ;;
    missing-source) W_ART1="94,2026-09-23T08:45:00Z,," ;;
    partial-green) W_TOTAL=350; W_ART1="95,2026-09-23T08:45:00Z,706,dev"
             W_RUNS="706,.github/workflows/ci.yml,push,dev" ;;
    partial-notfound) W_TOTAL=350 ;;
    future-sample) W_SAMP=$samp_future ;;
    future-job) W_TRIG=$trig_future; W_SAMP=$samp_fresh ;;
    stale-sample) W_SAMP=$samp_stale ;;
    k1-no-token) W_TRIG=''; W_SAMP=''; W_TOKEN='' ;;
    k2-no-repo) W_TRIG=''; W_SAMP=''; W_REPO='' ;;
    k3-runs-fail) W_TRIG=''; W_FAIL=runs-list; W_SAMP=$samp_fresh ;;
    k4-artifacts-fail) W_FAIL=artifacts ;;
    k4b-total-fail) W_FAIL=total; W_ART1="97,2026-09-23T08:45:00Z,708,dev"
             W_RUNS="708,.github/workflows/ci.yml,push,dev" ;;
    k5-empty) W_TOTAL=0 ;;
    k6-empty-stale) W_TOTAL=0; W_TRIG='' ;;
    k7-run-fail) W_FAIL=run; W_ART1="96,2026-09-23T08:45:00Z,707,dev"
             W_RUNS="707,.github/workflows/ci.yml,push,dev" ;;
    *) echo "drive.sh: unknown world: $w" >&2; exit 2 ;;
    esac
    wdir=$ws/worlds/$w
    mkdir -p "$wdir"
    printf '{"total_count":%s,"artifacts":[]}' "$W_TOTAL" > "$wdir/total.json"
    emit_page "$W_TOTAL" "$W_ART1" > "$wdir/artifacts-page1.json"
    emit_page "$W_TOTAL" "$W_ART2" > "$wdir/artifacts-page2.json"
    emit_page "$W_TOTAL" "$W_ART3" > "$wdir/artifacts-page3.json"
    printf '{"total_count":0,"workflow_runs":[]}' > "$wdir/runs.json"
    printf '{"jobs":[]}' > "$wdir/jobs.json"
    printf '%s\n' "$W_RUNS" | tr ';' '\n' | while IFS=, read -r rid rpath revent rbranch; do
        [ -n "$rid" ] || continue
        printf '{"id":%s,"path":"%s","event":"%s","head_branch":"%s"}' "$rid" "$rpath" "$revent" "$rbranch" > "$wdir/run-$rid.json"
    done
}

run_world() {
    w=$1
    setup_world "$w"
    mkdir -p "$ws/out"
    logf=$ws/out/$w.out
    env -i \
        PATH="$path_in" \
        HOME="${HOME:-/tmp}" \
        STUB_WORLD="$wdir" \
        STUB_JQ="$dir/ghjq.py" \
        STUB_PY="$stub_py" \
        STUB_FAIL="$W_FAIL" \
        GITHUB_REPOSITORY="$W_REPO" \
        GH_TOKEN="$W_TOKEN" \
        SLO_FRESH_NOW="$W_NOW" \
        SLO_FULL_LAST_TRIGGER="$W_TRIG" \
        SLO_FULL_LAST_SAMPLE="$W_SAMP" \
        SLO_FULL_ARTIFACT_PAGES="$W_PAGES" \
        sh "$script" >"$logf" 2>&1
    rc=$?
    {
        printf '== %s rc=%s\n' "$w" "$rc"
        sed 's/^/   | /' "$logf"
        printf '\n'
    } >> "$outfile"
}

: > "$outfile"
mkdir -p "$ws/out"
for w in ${only:-$all_worlds}; do
    run_world "$w"
done
echo "--- rc summary ($script) ---"
grep '^== ' "$outfile"
