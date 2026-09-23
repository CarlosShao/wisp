#!/bin/sh
# scripts/slo-freshness.sh - the nail ticket 134 AC#3 asks for: a gate that goes
# RED when the slo-full job stops being triggered, and (since ticket 134 AC#6)
# RED when the job keeps being triggered but stops producing valid samples.
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
# Why probe P3 was added (AC#6, the other half of the same trade). AC#6 lets
# scripts/slo-check.ps1 exit 0 when it refuses to sample a busy machine
# ("NO CONCLUSION (machine-contended)", owner 2026-09-23 Q-36). Colouring that
# run green makes P2 useless on its own: the job is now triggered, completes
# and stays "fresh" every single push even when it produced zero numbers for a
# month. So the aging basis may not be "did a job happen", it has to be "did a
# VALID SAMPLE happen". P3 is that, and it reads only things obtainable from
# the API - never a job's conclusion:
#
#   the record P3 ages on = the newest workflow artifact named `slo-full-report`
#   why that is a valid sample  = ci.yml's Upload step names exactly
#                                 build/slo/slo-report.json, and
#                                 scripts/slo-check.ps1 writes that file only
#                                 on a path that actually sampled (it clears
#                                 the previous run's *.json first, and prints
#                                 "0 state file(s) written, slo-report.json NOT
#                                 written" when it refuses). upload-artifact
#                                 creates no artifact at all when nothing
#                                 matches, so "artifact exists" == "report
#                                 existed" == "numbers were produced".
#   what it does NOT prove     = that the numbers passed. A sample that FAILED
#                                (all_pass=false) exits 1, the Upload step is
#                                never reached, so there is no artifact and the
#                                P3 clock does not move. That is on purpose:
#                                such a run reddens the badge by itself, and
#                                counting it would let a permanently-failing
#                                gate look healthy - the same silent death this
#                                nail exists to catch.
#
# Three probes, all required:
#   P1 static   - .github/workflows/ci.yml still carries the triggers that make
#                 slo-full auto-run (push on main/dev + a schedule), and the
#                 slo-full job still has no `if:` / `continue-on-error`. This
#                 probe needs no network, so it catches "someone edited the
#                 trigger block" on the same commit that did it.
#   P2 dynamic  - the newest slo-full JOB in a `ci` workflow run on main/dev is
#                 younger than SLO_FULL_MAX_AGE_DAYS. This catches the deaths
#                 P1 cannot: the self-hosted runner goes offline and the job
#                 sits in `queued` forever, the schedule gets removed on main
#                 only, or the whole workflow stops being referenced. Since
#                 AC#6 this probe alone is NOT enough (see above) - it is kept
#                 because it is still the only one that sees a queued job.
#   P3 sampling - the newest `slo-full-report` ARTIFACT is younger than
#                 SLO_FULL_SAMPLE_MAX_AGE_DAYS, and at least one exists. This
#                 is the probe that stays sharp when every run is a
#                 machine-contended no-conclusion run.
#
# Failure tokens are stable strings so a log grep can name the cause:
#   slo-full-trigger-missing   /  slo-full-never-triggered  /  slo-full-stale
#   slo-full-sample-never      /  slo-full-sample-stale
#
# Test seams (documented, and NOT set by the workflow below): SLO_FRESH_NOW
# overrides "now", SLO_FULL_LAST_TRIGGER injects a job record instead of
# querying GitHub, SLO_FULL_LAST_SAMPLE injects (or, as the literal `none`,
# fixtures away) a valid-sample record, SLO_FRESH_CI points at a different copy
# of ci.yml. They exist so this pin can be mutation-tested without touching a
# shared file or the API; every one of them shapes the INPUT of a check, none
# of them turns a check off. There is no seam that makes a probe pass: an empty
# SLO_FULL_LAST_TRIGGER still queries the API, `none` for SLO_FULL_LAST_SAMPLE
# forces the reddest branch there is, and a missing token exits 2 rather than
# exiting 0 (D22 mode-6: no skip switch).
#
# Exit codes: 0 = all three probes agree slo-full is alive and sampling,
# 1 = at least one failed, 2 = the pin could not look (no gh, no token, no
# ci.yml, unparseable date).
set -eu

# Spell the empty prefix assignment as `CDPATH=''` rather than `CDPATH=`.
# Identical shell semantics (an empty CPATH for the `cd` that follows), but the
# linter (shellcheck 0.11.0) reads the bare `CDPATH= cd` spelling as a typo
# (SC1007) and the "Shell lint for the pin" step in
# .github/workflows/slo-fresh.yml is a hard gate - leaving it would have made
# this nail's own workflow red the first time it ran on a runner that has a
# linter. Measured 2026-09-23, see
# docs/evidence/s1/134-ac6-contended-no-conclusion.md section 3.4.
here=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH='' cd -- "$here/.." && pwd)

max_age_days=${SLO_FULL_MAX_AGE_DAYS:-3}
sample_max_age_days=${SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3}
sample_pages=${SLO_FULL_ARTIFACT_PAGES:-3}
ci_file=${SLO_FRESH_CI:-$root/.github/workflows/ci.yml}
record=${SLO_FULL_LAST_TRIGGER:-}
sample_record=${SLO_FULL_LAST_SAMPLE:-}
now_override=${SLO_FRESH_NOW:-}

case $max_age_days in
'' | *[!0-9]*)
    echo "slo-freshness: SLO_FULL_MAX_AGE_DAYS must be a whole number of days, got '$max_age_days'" >&2
    exit 2
    ;;
esac
case $sample_max_age_days in
'' | *[!0-9]*)
    echo "slo-freshness: SLO_FULL_SAMPLE_MAX_AGE_DAYS must be a whole number of days, got '$sample_max_age_days'" >&2
    exit 2
    ;;
esac
case $sample_pages in
'' | *[!0-9]*)
    echo "slo-freshness: SLO_FULL_ARTIFACT_PAGES must be a whole number of pages, got '$sample_pages'" >&2
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
# here. Re-measured 2026-09-23 at ticket 134 AC#6 (the "6" this comment used to
# claim was stale): ci.yml carries 9 indented `if:` lines, 8 of them
# `${{ !cancelled() }}` and one `always()`, at :178 :215 :291 :351 :389 :399
# :421 :457 :473 - ALL of them outside the slo-full job (whose key is at :537),
# so a file-wide grep would pass no matter what happens to this job.
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

# ------------------------------------------------------------- probes P2/P3 --
to_epoch() {
    # GNU date reads the RFC3339 stamps GitHub emits ("2026-09-23T02:11:10Z").
    # A parse failure is exit 2, never a silent 0: "cannot tell" is not "fresh".
    if ! out=$(date -u -d "$1" +%s 2>/dev/null); then
        echo "slo-freshness: cannot parse timestamp '$1' - refusing to guess its age" >&2
        exit 2
    fi
    printf '%s\n' "$out"
}

now_epoch() {
    # One clock for both dynamic probes, so P2 and P3 can never disagree about
    # what "now" is (SLO_FRESH_NOW is a test seam for both of them at once).
    if [ -n "$now_override" ]; then
        to_epoch "$now_override"
    else
        date -u +%s
    fi
}

can_look() {
    # Shared precondition of every probe that has to ask GitHub something.
    # "Cannot look" exits 2 - it is never a pass and never a fail(1), because a
    # pin that silently skips when it has no token is the disease this file
    # exists to refuse. Called once per probe, before any of its queries.
    repo=${GITHUB_REPOSITORY:-}
    if [ -z "$repo" ]; then
        echo "slo-freshness: no injection seam and no GITHUB_REPOSITORY - cannot ask GitHub anything" >&2
        exit 2
    fi
    if [ -z "${GH_TOKEN:-}${GITHUB_TOKEN:-}" ]; then
        echo "slo-freshness: no GH_TOKEN/GITHUB_TOKEN - the freshness probes cannot look (this is not a pass)" >&2
        exit 2
    fi
    if ! command -v gh >/dev/null 2>&1; then
        echo 'slo-freshness: gh is not on PATH, and this pin has no offline fallback' >&2
        exit 2
    fi
}

if [ -z "$record" ]; then
    can_look
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

now=$(now_epoch)

if [ -z "$record" ]; then
    fail 'slo-full-never-triggered: no slo-full job found in the 25 newest ci runs on main/dev'
else
    IFS='|' read -r created status conclusion url <<EOF2
$record
EOF2
    if [ -z "$created" ]; then
        fail "slo-full-never-triggered: the record has no created_at: $record"
    else
        then=$(to_epoch "$created")
        age=$((now - then))
        if [ "$age" -lt 0 ]; then age=0; fi
        age_days=$((age / 86400))
        echo "slo-freshness: P2 newest slo-full job record: created_at=$created status=$status conclusion=$conclusion"
        echo "slo-freshness:   job url: $url"
        echo "slo-freshness:   age: $age_days day(s) ($age s); threshold: $max_age_days day(s)"
        if [ "$age_days" -gt "$max_age_days" ]; then
            fail "slo-full-stale: no slo-full trigger record for $age_days day(s) (> $max_age_days). Either nobody pushed and the schedule stopped firing, or the self-hosted runner has been deaf. D32 has had no automatic verdict since $created."
        elif [ "$status" != completed ] && [ $((age / 3600)) -gt 12 ]; then
            fail "slo-full-stale: the newest slo-full job is still '$status' after $((age / 3600)) hour(s) - a queued job that never starts is a silent death too."
        fi
    fi
fi

# ------------------------------------------------------------------- P3 ------
# Ticket 134 AC#6: age on the last VALID SAMPLE, not on the last job. See the
# header of this file for why P2 alone stopped being enough the moment a
# machine-contended run learned to exit 0.
#
# Records are TSV: created_at <TAB> workflow-run-id <TAB> artifact-id.
#
# ONE source of truth for the name. It is a plain assignment, not a knob: an
# env-overridable name would let whoever sets it to `slo-smoke-report` blind
# P3 from the other side (measured - see the AC#6 evidence, mutation M6: with
# that one string swapped, P3 reads a hosted smoke artifact uploaded 55 minutes
# after the self-hosted gate last produced anything and calls it a valid D32
# sample). The label in the messages interpolates the same variable, so the
# filter and what the pin claims to have read cannot drift apart.
sample_artifact=slo-full-report
sample_records=''
sample_scan_note=''
if [ "$sample_record" = none ]; then
    # Fixture branch: the world in which the artifact scan came back empty,
    # i.e. runs keep happening but not one of them produced a report. Only a
    # test seam can reach it with `none`; the real scan reaches it by finding
    # nothing, which is the same code path from here down.
    sample_scan_note="SLO_FULL_LAST_SAMPLE=none (fixture: scan returned no $sample_artifact artifact)"
elif [ -n "$sample_record" ]; then
    sample_records=$sample_record
    sample_scan_note='SLO_FULL_LAST_SAMPLE injected (test seam, 1 record)'
else
    can_look
    # Measured 2026-09-23 on this repo, two facts that both bite:
    #   * the listing is ordered by ARTIFACT ID descending, not by created_at:
    #     page 1 started 05:39:19 / 05:25:01 / 04:53:50 / 05:39:05 / 04:44:08,
    #     and its last three rows were 10:12:45 / 10:00:10 / 10:13:42 (the
    #     SMALLEST id carrying the NEWEST stamp). So "take the first row" is
    #     wrong and "sort the text" is not enough either - keep the max.
    #   * pages past the end answer with an empty artifacts array, not an
    #     error (per_page=100 -> page 1 = 100 rows, page 2 = 70, pages 3-5 = 0),
    #     so the loop below is bounded and cheap even when the cap overshoots.
    # Cap: SLO_FULL_ARTIFACT_PAGES * 100. If a future repo ever holds more
    # artifacts than that the scan is PARTIAL - and a partial max can only be
    # older than the true newest, so the error direction is RED, never green.
    sample_page=1
    while [ "$sample_page" -le "$sample_pages" ]; do
        if ! page_records=$(gh api \
            "repos/$repo/actions/artifacts?per_page=100&page=$sample_page" \
            --jq '.artifacts[]
                  | select(.name=="'"$sample_artifact"'")
                  | select(.expired==false)
                  | [(.created_at // ""), ((.workflow_run.id // 0)|tostring), ((.id // 0)|tostring)]
                  | join("\t")' 2>&1 </dev/null); then
            echo "slo-freshness: P3 the artifacts query (page $sample_page) failed, GitHub said:" >&2
            printf '%s\n' "$page_records" >&2
            exit 2
        fi
        if [ -n "$page_records" ]; then
            if [ -n "$sample_records" ]; then
                sample_records="$sample_records
$page_records"
            else
                sample_records=$page_records
            fi
        fi
        sample_page=$((sample_page + 1))
    done
    sample_scan_note="scanned $sample_pages page(s) x 100 of the artifact listing"
    if [ -z "$sample_records" ]; then
        echo "slo-freshness: P3 $sample_scan_note - no $sample_artifact artifact found"
    fi
fi

sample_count=0
best_ts=''
best_run=''
best_art=''
best_epoch=-1
if [ -n "$sample_records" ]; then
    tab=$(printf '\t')
    while IFS=$tab read -r ts art_run art_id; do
        [ -n "$ts" ] || continue
        sample_count=$((sample_count + 1))
        e=$(to_epoch "$ts")
        if [ "$e" -gt "$best_epoch" ]; then
            best_epoch=$e
            best_ts=$ts
            best_run=$art_run
            best_art=$art_id
        fi
    done <<EOF3
$sample_records
EOF3
fi

if [ "$sample_count" -eq 0 ]; then
    fail "slo-full-sample-never: no uploaded $sample_artifact artifact at all ($sample_scan_note). slo-full is being triggered but has produced no valid sample that this pin can see - which since AC#6 is NOT the same question as 'is the job green'. D32 (Sleeping CPU <=0.5%, private RSS <=25MB) is unverified."
else
    sample_age=$((now - best_epoch))
    if [ "$sample_age" -lt 0 ]; then sample_age=0; fi
    sample_age_days=$((sample_age / 86400))
    echo "slo-freshness: P3 newest VALID SAMPLE: artifact $sample_artifact created_at=$best_ts artifact_id=$best_art run=$best_run"
    echo "slo-freshness:   $sample_scan_note; $sample_count valid-sample record(s) considered"
    echo "slo-freshness:   age: $sample_age_days day(s) ($sample_age s); threshold: $sample_max_age_days day(s)"
    if [ "$sample_age_days" -gt "$sample_max_age_days" ]; then
        fail "slo-full-sample-stale: the last VALID SLO SAMPLE is $sample_age_days day(s) old (> $sample_max_age_days). Jobs may well be green - since ticket 134 AC#6 a machine-contended run exits 0 without a report - so a green job is not the question. The question is when numbers last existed, and that was $best_ts (run $best_run). D32 has been unverified since then."
    fi
fi

if [ "$FAILURES" -ne 0 ]; then
    echo "slo-freshness: $FAILURES probe failure(s) - ticket 134 AC#3/AC#6 nail is red" >&2
    exit 1
fi
echo 'slo-freshness: OK - slo-full still has automatic triggers, a trigger record inside the window, and a VALID SAMPLE inside the window'
exit 0
