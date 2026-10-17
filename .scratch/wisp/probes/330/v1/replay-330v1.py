#!/usr/bin/env python3
# ticket 330 verdict leg 330-v1: independent replay of 330-r2's paired G2
# mutations. NOTHING here edits a tracked file: every byte touched lives in the
# out-of-repo export tree D:\tmp\w330v1 (git archive of b9069dfa, the delivery
# tip). Each shot restores its files in a finally block and re-checks the
# restored sha256 against the committed baseline; a failed restore exits 4.
import hashlib
import os
import subprocess
import sys

TREE = r"D:\tmp\w330v1"
LOGDIR = r"D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\330\v1\logs"
GOTMP = r"D:\tmp\w330v1tmp"
DOORBELL = os.path.join(TREE, "internal", "agent", "plan_present_doorbell_165r3_test.go")
WORLD = os.path.join(TREE, "internal", "agent", "plan_present_crowded_world_330r2_test.go")
BASE = {
    DOORBELL: "6542773d3d2edcc6d554a47f091ce2dffe35586a92c1f948f79c361454683fca",
    WORLD: "891fbf5ccc9bdcd23f5075538c2a590e441dc87a100c2634a45d29cbbdd2b59d",
}
RUN = ("^(TestCrowdedWorldStillOffersTheBrake330r2|"
       "TestResidencyTierDecidesTheInjectionArm330r2|"
       "TestPlanPresentCallBooksTheToolCallRow165r3)$")

RES_BRAKE_OLD = "\t\tResident:    true,\n"
RES_BRAKE_NEW = "\t\tResident:    false,\n"
RES_DECOY_OLD = "\t\tResident:    false,\n"
RES_DECOY_NEW = "\t\tResident:    true,\n"
QUERY_CN = 'const crowdedWorldQuery330r2 = "帮我把这个项目收拾一下"'
QUERY_EN = 'const crowdedWorldQuery330r2 = "frobnicate widget quietly"'
PARMS = ("\t\tParameters:  json.RawMessage(`{\"type\":\"object\",\"properties\":{\"root\":"
         "{\"type\":\"string\"}},\"required\":[\"root\"],\"additionalProperties\":false}`),\n")
WORLD_BLOCK_OLD = (
    "\treturn append(planPresentDirectory165r3(), ToolInfo{\n"
    "\t\tName:        longTailName330r2,\n"
    "\t\tDescription: longTailDesc330r2,\n"
    "\t\tResident:    false,\n"
    "\t\tRiskLevel:   RiskL1,\n"
    + PARMS +
    "\t})\n")
WORLD_BLOCK_NEW = (
    "\tvar droppedEntry330r2 = ToolInfo{\n"
    "\t\tName:        longTailName330r2,\n"
    "\t\tDescription: longTailDesc330r2,\n"
    "\t\tResident:    false,\n"
    "\t\tRiskLevel:   RiskL1,\n"
    + PARMS +
    "\t}\n"
    "\t_ = droppedEntry330r2\n"
    "\treturn planPresentDirectory165r3()\n")


def sha(path):
    with open(path, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()


def read(path):
    with open(path, "r", encoding="utf-8", newline="") as f:
        return f.read()


def write(path, text):
    with open(path, "w", encoding="utf-8", newline="") as f:
        f.write(text)


def run_tests(logname, full=False):
    args = ["go", "test", "./internal/agent/", "-count=1", "-timeout", "30m", "-v"]
    if not full:
        args += ["-run", RUN]
    env = dict(os.environ, GOTMPDIR=GOTMP)
    r = subprocess.run(args, cwd=TREE, env=env, capture_output=True, text=True, encoding="utf-8")
    out = (r.stdout or "") + "\n[stderr]\n" + (r.stderr or "")
    with open(os.path.join(LOGDIR, logname), "w", encoding="utf-8", newline="\n") as f:
        f.write("### rc=%d\n%s" % (r.returncode, out))
    keep = []
    for ln in out.splitlines():
        if (ln.startswith("=== RUN") or ln.startswith("--- PASS") or ln.startswith("--- FAIL")
                or ln.startswith("--- SKIP") or "_test.go:" in ln or ln in ("FAIL", "PASS")
                or ln.startswith("ok  \t") or "build failed" in ln or "cannot use" in ln
                or ln.startswith("imports") or "declared and not used" in ln):
            keep.append(ln)
    return r.returncode, keep


def counts(logname):
    with open(os.path.join(LOGDIR, logname), encoding="utf-8") as f:
        t = f.read()
    def gc(p):
        return sum(1 for ln in t.splitlines() if ln.startswith(p))
    return ("RUN=%d PASS_TOP=%d FAIL_TOP=%d SKIP=%d PANIC=%d"
            % (gc("=== RUN"), gc("--- PASS"), gc("--- FAIL"), gc("--- SKIP"), gc("panic")))


def forward_full():
    rc, _ = run_tests("v1-forward-fullpackage.log", full=True)
    print("=== forward FULL package rc=%d  %s" % (rc, counts("v1-forward-fullpackage.log")))


def forward_control():
    rc, keep = run_tests("v1-forward-named3.log")
    print("=== forward named-3-cases rc=%d  %s" % (rc, counts("v1-forward-named3.log")))
    for ln in keep:
        print("    | %s" % ln)


def shot(tag, edits, logname):
    pristine = {}
    mutated = {}
    for path, old, new in edits:
        pristine.setdefault(path, read(path))
        cur = mutated.get(path, pristine[path])
        if cur.count(old) != 1:
            print("ABORT %s: needle occurs %d times in %s" % (tag, cur.count(old), path))
            return 3
        mutated[path] = cur.replace(old, new, 1)
    print("=== shot=%s" % tag)
    try:
        for path, text in mutated.items():
            write(path, text)
            print("    mutated %s sha=%s" % (path, sha(path)))
        rc, keep = run_tests(logname)
        print("    rc=%d  %s  raw_log=%s" % (rc, counts(logname), os.path.join(LOGDIR, logname)))
        for ln in keep:
            print("    | %s" % ln)
    finally:
        for path, text in pristine.items():
            write(path, text)
            same = sha(path) == BASE[path]
            print("    restored %s sha=%s matches_committed_baseline=%s" % (path, sha(path), same))
            if not same:
                print("    !!! RESTORE FAILED - bytes differ from the committed baseline")
                sys.exit(4)
    return 0


def main():
    os.makedirs(LOGDIR, exist_ok=True)
    os.makedirs(GOTMP, exist_ok=True)
    for p in (DOORBELL, WORLD):
        if sha(p) != BASE[p]:
            print("ABORT: export-tree %s sha=%s != committed baseline %s" % (p, sha(p), BASE[p]))
            return 5
    print("tree preflight ok: both target files match the committed baseline at b9069dfa")
    forward_control()
    rc = 0
    rc |= shot("M1_replay-of-G2_mirror-brake-Resident-true-to-false",
               [(DOORBELL, RES_BRAKE_OLD, RES_BRAKE_NEW)], "v1-m1-g2-replay.log")
    rc |= shot("M2_defang-added-non-resident-candidate-Resident-false-to-true",
               [(WORLD, RES_DECOY_OLD, RES_DECOY_NEW)], "v1-m2-defang.log")
    rc |= shot("M3_control-query-with-no-common-term-with-the-decoy",
               [(WORLD, QUERY_CN, QUERY_EN)], "v1-m3-control.log")
    rc |= shot("M4_only-the-added-candidate-deleted-no-tier-mutation",
               [(WORLD, WORLD_BLOCK_OLD, WORLD_BLOCK_NEW)], "v1-m4-deleted-only.log")
    rc |= shot("M5_mandatory-question-added-candidate-deleted-AND-G2-replay",
               [(WORLD, WORLD_BLOCK_OLD, WORLD_BLOCK_NEW),
                (DOORBELL, RES_BRAKE_OLD, RES_BRAKE_NEW)], "v1-m5-deleted-plus-g2-replay.log")
    print("all shots done, rc=%d" % rc)
    return rc


if __name__ == "__main__":
    sys.exit(main())
