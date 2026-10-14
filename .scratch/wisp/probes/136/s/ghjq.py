#!/usr/bin/env python
# ghjq.py - the `--jq` half of the ticket-136s `gh` stub (this host has no jq).
#
# The stub feeds it one fixture JSON and the exact filter text
# scripts/slo-freshness.sh passed to `gh api --jq`; this file implements only
# the filters that pin uses, and an UNKNOWN filter is an error (exit 1) - the
# pin then takes its "cannot look" path (exit 2) instead of being served a
# made-up answer, which is the fail-closed direction this whole ticket is about.
# Whitespace is normalised before matching so the filter text may be re-flowed
# in the script without silently disabling the stub.
#
# Covered filters (as of ticket 136s):
#   .total_count
#   .artifacts[] | select(.name=="slo-full-report") | select(.expired==false)
#     | [created_at, workflow_run.id, id] | join("\t")                (3 cols)
#     | [created_at, workflow_run.id, id, workflow_run.head_branch]   (4 cols)
#   [.path, .event, .head_branch] | join("\t")                        (one run)
#   .workflow_runs[] | select(event in push|schedule|workflow_dispatch)
#     | select(head_branch in main|dev) | .id                 (P2 runs list)
#   .jobs[] | select(.name=="slo-full") | [created_at, status,
#     conclusion, html_url] | join("|")                       (P2 job lookup)
import json
import re
import sys

def norm(f):
    return re.sub(r"\s+", " ", f).strip()

NAME = 'slo-full-report'
F_TOTAL = '.total_count'
F_ART3 = ('.artifacts[] | select(.name=="slo-full-report") | select(.expired==false) '
          '| [(.created_at // ""), ((.workflow_run.id // 0)|tostring), ((.id // 0)|tostring)] '
          '| join("\\t")')
F_ART4 = ('.artifacts[] | select(.name=="slo-full-report") | select(.expired==false) '
          '| [(.created_at // "-"), ((.workflow_run.id // 0)|tostring), ((.id // 0)|tostring), '
          '(.workflow_run.head_branch // "")] | join("\\t")')
F_RUN = '[(.path // ""), (.event // ""), (.head_branch // "")] | join("\\t")'
F_RUNS = ('.workflow_runs[] | select(.event=="push" or .event=="schedule" or .event=="workflow_dispatch") '
          '| select(.head_branch=="main" or .head_branch=="dev") | .id')
F_JOBS = ('.jobs[] | select(.name=="slo-full") '
          '| [(.created_at // ""), .status, (.conclusion // "none"), .html_url] | join("|")')

def s(v):
    return v if isinstance(v, str) else ("" if v is None else str(v))

def jqalt(v, alt):
    # jq's `//`: null (and false, and an absent key) falls through to alt; an
    # empty string does NOT.
    return alt if v is None or v is False else v

def main():
    if len(sys.argv) != 3:
        sys.stderr.write("ghjq.py: usage: ghjq.py <fixture.json> <filter>\n")
        return 1
    try:
        # Windows text mode would translate \n to \r\n, and the \r would ride
        # into the LAST field of a TSV row (where the pin's branch predicate
        # reads it). Output stays LF-only, like real jq.
        sys.stdout.reconfigure(newline="\n")
    except AttributeError:
        pass
    with open(sys.argv[1], "r", encoding="utf-8") as fh:
        data = json.load(fh)
    f = norm(sys.argv[2])
    out = []
    if f == norm(F_TOTAL):
        tc = data.get("total_count")
        out.append(str(tc) if isinstance(tc, int) else "null")
    elif f in (norm(F_ART3), norm(F_ART4)):
        four = (f == norm(F_ART4))
        for a in data.get("artifacts", []):
            if a.get("name") != NAME or a.get("expired") is not False:
                continue
            wr = a.get("workflow_run") or {}
            ts = s(jqalt(a.get("created_at"), "-")) if four else s(a.get("created_at"))
            row = [ts, str(wr.get("id") or 0), str(a.get("id") or 0)]
            if four:
                row.append(s(wr.get("head_branch")))
            out.append("\t".join(row))
    elif f == norm(F_RUN):
        out.append("\t".join([s(data.get("path")), s(data.get("event")), s(data.get("head_branch"))]))
    elif f == norm(F_RUNS):
        for r in data.get("workflow_runs", []):
            if r.get("event") not in ("push", "schedule", "workflow_dispatch"):
                continue
            if r.get("head_branch") not in ("main", "dev"):
                continue
            out.append(str(r.get("id")))
    elif f == norm(F_JOBS):
        for j in data.get("jobs", []):
            if j.get("name") != "slo-full":
                continue
            out.append("|".join([s(j.get("created_at")), s(j.get("status")),
                                 s(j.get("conclusion")) if j.get("conclusion") is not None else "none",
                                 s(j.get("html_url"))]))
    else:
        sys.stderr.write("ghjq.py: unsupported filter (refusing to guess): %s\n" % f)
        return 1
    for line in out:
        print(line)
    return 0

if __name__ == "__main__":
    sys.exit(main())
