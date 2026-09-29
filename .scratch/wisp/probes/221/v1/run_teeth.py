import json
import os
import subprocess

ROOT = r"D:/work/workspace/projects plans/Wisp"
MUT = os.path.join(ROOT, ".scratch", "wisp", "probes", "221", "v1", "mutations")
LOGS = os.path.join(ROOT, ".scratch", "wisp", "probes", "221", "v1", "logs")
TEST_REL = "internal/tools/task_cancel_221_legs_test.go"


def ws(p):
    return p.replace("/", "\\")


KEY = ws(ROOT + "/" + TEST_REL)

for m in ("control", "m13-list-marker-lifted", "m14-cancel-marked-again"):
    val = ws(os.path.join(MUT, "teeth-%s-test.go" % m))
    with open(os.path.join(MUT, "teeth-%s.json" % m), "w", encoding="utf-8") as f:
        f.write(json.dumps({"Replace": {KEY: val}}))
    print("wrote overlay for", m, "->", os.path.exists(os.path.join(MUT, "teeth-%s-test.go" % m)))

env = dict(os.environ)
env["PATH"] = os.path.join(ROOT, "third_party", "sherpa-onnx") + os.pathsep + \
    os.path.join(ROOT, "build") + os.pathsep + env["PATH"]

for m in ("control", "m13-list-marker-lifted", "m14-cancel-marked-again"):
    log = os.path.join(LOGS, "teeth-%s.txt" % m)
    with open(log, "wb") as f:
        rc = subprocess.call(["go", "test", "./internal/tools", "-count=1", "-v",
                              "-run", "Test221DeferredMarkerForCancelLiftedButListStillMarked",
                              "-overlay", os.path.join(MUT, "teeth-%s.json" % m)],
                             cwd=ROOT, stdout=f, stderr=subprocess.STDOUT, env=env)
    body = open(log, encoding="utf-8", errors="replace").read().split("\n")
    reds = [l for l in body if l.strip().startswith("--- FAIL")]
    msgs = [l for l in body if "task_cancel_221_legs_test.go:" in l]
    print("teeth-%s rc=%s red=%d" % (m, rc, len(reds)))
    for l in reds:
        print("   ", l.strip())
    for l in msgs[:2]:
        print("    MSG", l.strip()[:180])

st = subprocess.run(["git", "-C", ROOT, "status", "--porcelain", "--", "internal", "cmd"],
                    capture_output=True, text=True).stdout
print("worktree after teeth:", repr(st))
