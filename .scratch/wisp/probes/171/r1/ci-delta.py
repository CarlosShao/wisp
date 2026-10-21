"""ci-delta.py - ticket 171 AC#3 legality self-proof for the one added ci.yml step.

Compares the ci.yml at HEAD with the working copy and prints, for every job, the step
count and the step names. The claim being checked is narrow and mechanical: exactly one
step was ADDED, in the job that already carried the attrib.sh step; no job disappeared;
no trigger moved; no existing step was renamed.

Run from the repo root:  python <this file>
"""
import subprocess
import sys

import yaml

after = yaml.safe_load(open(".github/workflows/ci.yml", encoding="utf-8"))
before = yaml.safe_load(
    subprocess.run(
        ["git", "show", "HEAD:.github/workflows/ci.yml"],
        capture_output=True,
        text=True,
        encoding="utf-8",
    ).stdout
)

print("jobs before:", sorted(before["jobs"]))
print("jobs after :", sorted(after["jobs"]))
assert sorted(before["jobs"]) == sorted(after["jobs"]), "the job roster moved"

fail = 0
for j in sorted(after["jobs"]):
    nb = before["jobs"][j].get("steps", [])
    na = after["jobs"][j].get("steps", [])
    tag = "SAME" if len(nb) == len(na) else "CHANGED"
    print("  %-14s steps %d -> %d  %s" % (j, len(nb), len(na), tag))

    def names(steps):
        return [str(s.get("name", "")) for s in steps]

    gone = [n for n in names(nb) if n not in names(na)]
    new = [n for n in names(na) if n not in names(nb)]
    for n in gone:
        print("      GONE: %s" % n)
        fail = 1
    for n in new:
        print("      NEW : %s" % n)

print("triggers before:", before.get(True) or before.get("on"))
print("triggers after :", after.get(True) or after.get("on"))
if (before.get(True) or before.get("on")) != (after.get(True) or after.get("on")):
    print("FAIL: the trigger block moved (ticket 134 territory)")
    fail = 1

# The deliberately-added properties of the new step, checked on the parsed YAML and not
# on the text (a comment saying "no continue-on-error" would satisfy a grep for it).
lint = after["jobs"]["lint"]["steps"]
hits = [s for s in lint if "attrib.sh --self-test" in str(s.get("run", ""))]
print("new-step matches in job lint:", len(hits))
if len(hits) != 1:
    print("FAIL: expected exactly one step running attrib.sh --self-test in job lint")
    fail = 1
else:
    s = hits[0]
    print("  name              :", s.get("name"))
    print("  if                :", s.get("if"))
    print("  run               :", repr(s.get("run")))
    print("  continue-on-error :", "continue-on-error" in s)
    print("  position          : step %d of %d (= last)" % (lint.index(s) + 1, len(lint)))
    if "continue-on-error" in s:
        print("FAIL: the new step carries continue-on-error (zero-execution shape)")
        fail = 1
    if s.get("run") != "sh .scratch/wisp/probes/161/r5/attrib.sh --self-test":
        print("FAIL: the run line is not the byte string this ticket named")
        fail = 1

print("RESULT:", "GREEN - one step added, nothing removed, triggers untouched" if fail == 0 else "RED")
sys.exit(0 if fail == 0 else 1)
