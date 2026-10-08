#!/usr/bin/env python3
"""111-ciif1: insert step-level `if: ${{ !cancelled() }}` on named steps.

Dry-run by default; `--apply` writes. Adds ONLY `if:` lines, placed immediately
before the step's own `run:`/`uses:` key (the shape this file already uses at
:522 and :784). Refuses to touch a step that already carries an `if:` and
refuses to touch the slo-full job at all (scripts/slo-freshness.sh P1).
"""
import re
import sys
import yaml

PATH = ".github/workflows/ci.yml"
GUARD = "if: ${{ !cancelled() }}"

# (job, step index 1-based): no step-level `if:` AND below at least one step that
# can go red. Setup steps (index 1-2 everywhere) stay unguarded: that is the
# shape this file declares at :469 ("EVERY STEP BELOW THE SETUP ONES CARRIES")
# and :521. `Upload SLO report` stays unguarded: :229-233 registers that as
# deliberate (a failed sample must not produce an artifact).
TARGETS = {
    "lint": [3, 4, 5, 6, 7, 9, 10],
    "test-core": [3, 4, 5, 6],
    "slo-smoke": [3, 4, 5],
    "lint-frontend": [3, 4, 5, 6, 7, 8, 9, 10, 11],
}
FORBIDDEN = "slo-full"
apply = "--apply" in sys.argv
lines = open(PATH, "r", encoding="utf-8", newline="").read().split("\n")

# ---- jobs block: top-level `jobs:` at indent 0 ----
jobs_start = next(i for i, t in enumerate(lines, 1) if re.match(r"^jobs:\s*$", t))
jobs_end = len(lines)
for ln in range(jobs_start + 1, len(lines) + 1):
    if re.match(r"^[a-z][a-z0-9_.-]*:", lines[ln - 1]):
        jobs_end = ln - 1
        break
hdrs = [(ln, re.match(r"^  ([a-z][a-z0-9_-]*):\s*$", lines[ln - 1]).group(1))
        for ln in range(jobs_start + 1, jobs_end + 1)
        if re.match(r"^  ([a-z][a-z0-9_-]*):\s*$", lines[ln - 1])]
job_range = {}
for i, (ln, name) in enumerate(hdrs):
    job_range[name] = (ln, hdrs[i + 1][0] - 1 if i + 1 < len(hdrs) else jobs_end)

# ---- number the steps of each job from the raw lines ----
def steps_of(job):
    lo, hi = job_range[job]
    out = []
    cur = None
    for ln in range(lo, hi + 1):
        text = lines[ln - 1]
        if not text.strip():
            continue
        s = text.lstrip(" ")
        ind = len(text) - len(s)
        if s.startswith("#"):
            continue
        m = re.match(r"^ *-\s+(\S[^:]*):\s*(.*)$", text)
        if m and ind >= 6 and (ind - 2) % 2 == 0:
            step_indent = ind + 2
            if cur:
                out.append(cur)
            cur = {"name_line": ln, "name": f"{m.group(1)}: {m.group(2)}",
                   "has_if": None, "body": None, "indent": step_indent}
            key, lineno = m.group(1), ln
        elif cur is not None and ind >= cur["indent"] and ":" in s:
            key = s.split(":", 1)[0].strip()
            lineno = ln
            if cur["indent"] != ind:      # deeper -> nested map (with:/env:), ignore
                continue
        else:
            continue
        if key == "if" and cur["has_if"] is None:
            cur["has_if"] = lineno
        if key in ("run", "uses") and cur["body"] is None:
            cur["body"] = lineno
    if cur:
        out.append(cur)
    return out

plan = []
for job, idxs in TARGETS.items():
    steps = steps_of(job)
    for idx in idxs:
        st = steps[idx - 1]
        if st["has_if"]:
            print(f"SKIP already-guarded: {job} step {idx}")
            continue
        assert st["body"], f"{job} step {idx}: no run/uses key"
        plan.append({"job": job, "idx": idx, "name_line": st["name_line"],
                     "insert_before": st["body"], "indent": st["indent"],
                     "name": st["name"][:64],
                     "body_txt": lines[st["body"] - 1].strip()[:60]})
plan.sort(key=lambda p: p["insert_before"])
print(f"planned insertions: {len(plan)}")
for p in plan:
    print(f"  {p['job']:<14} step {p['idx']:>2} name@:{p['name_line']:>4} "
          f"insert before :{p['insert_before']:>4} ind={p['indent']:>2} "
          f"| {p['body_txt']} | {p['name']}")

flo, fhi = job_range[FORBIDDEN]
bad = [p for p in plan if flo <= p["insert_before"] <= fhi]
print(f"\n{FORBIDDEN} body = {flo}..{fhi}; planned edits inside = {len(bad)} (must be 0)")
assert not bad, "refusing to break slo-freshness.sh P1"

doc = yaml.safe_load(open(PATH, encoding="utf-8"))
print("jobs=%d steps=%d" % (len(doc["jobs"]),
      sum(len(j.get("steps") or []) for j in doc["jobs"].values())))
if not apply:
    print("DRY RUN (no write)")
    sys.exit(0)

out = list(lines)
for p in reversed(plan):
    out.insert(p["insert_before"] - 1, " " * p["indent"] + GUARD)
open(PATH, "w", encoding="utf-8", newline="").write("\n".join(out))
print(f"APPLIED {len(plan)} insertions")
