#!/bin/sh
# scripts/slo-freshness.sh - the nail ticket 134 AC#3 asks for: a gate that goes
# RED when the slo-full job stops being triggered.
#
# Why this file exists at all. Ticket 134 had to decouple slo-full from the
# review fleet that shares this laptop, and the tempting shape (option A: keep
# only `workflow_dispatch`) would have traded "a sample may be contaminated"
# (reversible: mark it noisy) for "D32's two hard thresholds no longer produce
# an automatic verdict" (irreversible: nobody notices the day the gate dies).
# That second disease has a name and a ticket number in this repo - ticket 85
# (lint-tools-never-produced-a-verdict), ticket 71 (gates-must-self-report), and
# the line in ci.yml that says "A skippable job is a job that will one day be
# skipped". So the decoupling was accepted ONLY together with this nail: if
# slo-full goes N days without a single trigger record, something here fails
# with a non-zero exit code and says so out loud.
#
# Two probes, both required:
#   P1 static   - .github/workflows/ci.yml still carries the triggers that make
#                 slo-full auto-run (push on main/dev + a schedule), and the
#                 slo-full job still has no `if:` / `continue-on-error`. This
#                 probe needs no network, so it catches "someone edited the
#                 trigger block" on the same commit that did it.
#   P2 dynamic  - the newest slo-full JOB in a `ci` workflow run on main/dev is
#                 younger than SLO_FULL_MAX_AGE_DAYS. This catches the quieter
#                 deaths: the self-hosted runner goes offline and the job sits
#                 in `queued` forever, the schedule gets removed on main only,
#                 or the whole workflow stops being referenced.
#
# Failure tokens are stable strings so a log grep can name the cause:
#   slo-full-trigger-missing / slo-full-never-triggered / slo-full-stale
#
# Test seams (documented, and NOT set by the workflow below): SLO_FRESH_NOW
# overrides "now", SLO_FULL_LAST_TRIGGER injects a job record instead of
# querying GitHub, SLO_FRESH_CI points at a different copy of ci.yml. They exist
# so this pin can be mutation-tested without touching a shared file or the API.
# There is no seam that turns the checks off: an empty SLO_FULL_LAST_TRIGGER
# still queries the API, and a missing token exits 2 rather than exiting 0
# (D22 mode-6: no skip switch).
#
# Exit codes: 0 = both probes agree slo-full is alive, 1 = at least one failed,
# 2 = the pin could not look (no gh, no token, no ci.yml, unparseable date).
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)

max_age_days=${SLO_FULL_MAX_AGE_DAYS:-3}
ci_file=${SLO_FRESH_CI:-$root/.github/workflows/ci.yml}
record=${SLO_FULL_LAST_TRIGGER:-}
now_override=${SLO_FRESH_NOW:-}

case $max_age_days in
'' | *[!0-9]*)
    echo "slo-freshness: SLO_FULL_MAX_AGE_DAYS must be a whole number of days, got '$max_age_days'" >&2
    exit 2
    ;;
esac

fail() {
    echo "slo-freshness: FAIL $1" >&2
    FAILURES=$((FAILURES + 1))
}

FAILURES=0

# ---------------------------------------------------------------- P1 static --
# The trigger block is checked where it actually lives (top-level `on:`), which
# is also what makes it cheap: if a future edit moves slo-full to a manual-only
# trigger, the two lines below are exactly what disappear.
if [ ! -f "$ci_file" ]; then
    echo "slo-freshness: no workflow file at $ci_file - cannot check the trigger block" >&2
    exit 2
fi

# `on:` is parsed as the YAML boolean true by some readers, so match both.
if ! grep -qE '^[[:space:]]*(on:|true:)[[:space:]]*$' "$ci_file"; then
    fail "slo-full-trigger-missing: $ci_file has no top-level 'on:' block at all"
fi
if ! grep -qE '^[[:space:]]*push:[[:space:]]*$' "$ci_file"; then
    fail "slo-full-trigger-missing: $ci_file has no 'push:' trigger; slo-full can only run when a human presses something"
fi
if ! grep -qE '^[[:space:]]*branches:[[:space:]]*\[[^]]*dev[^]]*\]|^[[:space:]]*-[[:space:]]*(main|dev)[[:space:]]*$' "$ci_file"; then
    fail "slo-full-trigger-missing: $ci_file pushes are not scoped to main/dev, so an automatic run is not guaranteed"
fi
if ! grep -qE "^[[:space:]]*(-[[:space:]]*)?cron:[[:space:]]*['\"]?[0-9]" "$ci_file"; then
    fail "slo-full-trigger-missing: $ci_file has no 'cron:' schedule entry (ticket 134 shape B); nobody has to press a button any more"
fi

# The job body must stay un-skippable. Anchored to the slo-full job so a
# legitimate `if: ${{ !cancelled() }}` in ANOTHER job cannot mask its removal
# here (measured: ci.yml carries 6 such lines outside slo-full, so a file-wide
# grep would pass no matter what happens to this job).
job_body=$(sed -n '/^  slo-full:/,/^  [a-z-]*:$/p' "$ci_file")
if [ -z "$job_body" ]; then
    fail "slo-full-trigger-missing: no 'slo-full:' job found in $ci_file - the merge gate it guards no longer exists"
else
    for banned in 'if:' 'continue-on-error:'; do
        if printf '%s\n' "$job_body" | grep -qE "^[[:space:]]+${banned%:}:"; then
            fail "slo-full-trigger-missing: the slo-full job now carries '$banned' (D22 mode-6: a skippable job is a job that will one day be skipped)"
        fi
    done
    if ! printf '%s\n' "$job_body" | grep -qE 'slo-check\.ps1[[:space:]]+-Subset[[:space:]]+full'; then
        fail "slo-full-trigger-missing: the slo-full job no longer runs scripts/slo-check.ps1 -Subset full"
    fi
fi

# ---------------------------------------------------------------- P2 dynamic --
to_epoch() {
    # GNU date reads the RFC3339 stamps GitHub emits ("2026-09-23T02:11:10Z").
    # A parse failure is exit 2, never a silent 0: "cannot tell" is not "fresh".
    if ! out=$(date -u -d "$1" +%s 2>/dev/null); then
        echo "slo-freshness: cannot parse timestamp '$1' - refusing to guess its age" >&2
        exit 2
    fi
    printf '%s\n' "$out"
}

if [ -z "$record" ]; then
    repo=${GITHUB_REPOSITORY:-}
    if [ -z "$repo" ]; then
        echo "slo-freshness: no SLO_FULL_LAST_TRIGGER and no GITHUB_REPOSITORY - cannot ask GitHub anything" >&2
        exit 2
    fi
    if [ -z "${GH_TOKEN:-}${GITHUB_TOKEN:-}" ]; then
        echo "slo-freshness: no GH_TOKEN/GITHUB_TOKEN - the freshness probe cannot look (this is not a pass)" >&2
        exit 2
    fi
    if ! command -v gh >/dev/null 2>&1; then
        echo 'slo-freshness: gh is not on PATH, and this pin has no offline fallback' >&2
        exit 2
    fi
    # Newest first, only the events that are supposed to be automatic, only the
    # branches that carry the merge gate. The first run that HAS a slo-full job
    # wins; jobs of runs still being created are handled by the age rule.
    if ! run_ids=$(gh api "repos/$repo/actions/workflows/ci.yml/runs?per_page=25" \
        --jq '.workflow_runs[]
              | select(.event=="push" or .event=="schedule" or .event=="workflow_dispatch")
              | select(.head_branch=="main" or .head_branch=="dev")
              | .id' 2>&1); then
        echo 'slo-freshness: the runs query failed, GitHub said:' >&2
        printf '%s\n' "$run_ids" >&2
        exit 2
    fi
    for rid in $run_ids; do
        line=$(gh api "repos/$repo/actions/runs/$rid/jobs?per_page=100" \
            --jq '.jobs[]
                  | select(.name=="slo-full")
                  | [(.created_at // ""), .status, (.conclusion // "none"), .html_url]
                  | join("|")' </dev/null || true)
        if [ -n "$line" ]; then
            record=$line
            break
        fi
    done
fi

if [ -z "$record" ]; then
    fail 'slo-full-never-triggered: no slo-full job found in the 25 newest ci runs on main/dev'
else
    IFS='|' read -r created status conclusion url <<EOF2
$record
EOF2
    if [ -z "$created" ]; then
        fail "slo-full-never-triggered: the record has no created_at: $record"
    else
        if [ -n "$now_override" ]; then
            now=$(to_epoch "$now_override")
        else
            now=$(date -u +%s)
        fi
        then=$(to_epoch "$created")
        age=$((now - then))
        if [ "$age" -lt 0 ]; then age=0; fi
        age_days=$((age / 86400))
        echo "slo-freshness: newest slo-full job record: created_at=$created status=$status conclusion=$conclusion"
        echo "slo-freshness:   job url: $url"
        echo "slo-freshness:   age: $age_days day(s) ($age s); threshold: $max_age_days day(s)"
        if [ "$age_days" -gt "$max_age_days" ]; then
            fail "slo-full-stale: no slo-full trigger record for $age_days day(s) (> $max_age_days). Either nobody pushed and the schedule stopped firing, or the self-hosted runner has been deaf. D32 has had no automatic verdict since $created."
        elif [ "$status" != completed ] && [ $((age / 3600)) -gt 12 ]; then
            fail "slo-full-stale: the newest slo-full job is still '$status' after $((age / 3600)) hour(s) - a queued job that never starts is a silent death too."
        fi
    fi
fi

if [ "$FAILURES" -ne 0 ]; then
    echo "slo-freshness: $FAILURES probe failure(s) - ticket 134 AC#3 nail is red" >&2
    exit 1
fi
echo 'slo-freshness: OK - slo-full still has automatic triggers and a trigger record inside the window'
exit 0
