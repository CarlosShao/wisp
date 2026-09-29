import json
import os
import subprocess
import sys

REPO = os.path.abspath(os.getcwd())
d = os.path.join(REPO, ".scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated")
os.makedirs(d, exist_ok=True)
rel = "internal/tools/subagent_222_test.go"
t = subprocess.run(["git", "cat-file", "blob", "HEAD:" + rel], cwd=REPO,
                   capture_output=True, check=True).stdout.decode("utf-8")

old = "\t\t\t\t\trelease: h.release,\n"
new = "\t\t\t\t\trelease: h.release,\n\t\t\t\t\tlate:     200 * time.Millisecond,\n"
assert t.count(old) == 1
t = t.replace(old, new, 1)

old = "\tstarted chan<- struct{}\n\trelease <-chan struct{}\n"
new = "\tstarted chan<- struct{}\n\trelease <-chan struct{}\n\tlate      time.Duration\n"
assert t.count(old) == 1
t = t.replace(old, new, 1)

old = "\t\t\tcase <-c.release:\n\t\t\tcase <-ctx.Done():\n\t\t\t\treturn ctx.Err()\n"
new = "\t\t\tcase <-c.release:\n\t\t\t\tif c.late > 0 {\n\t\t\t\t\ttime.Sleep(c.late)\n\t\t\t\t}\n\t\t\tcase <-ctx.Done():\n\t\t\t\treturn ctx.Err()\n"
assert t.count(old) == 1
t = t.replace(old, new, 1)

old = "h222PreFixBudget = 3 * time.Second"
new = "h222PreFixBudget = 50 * time.Millisecond"
assert t.count(old) == 1
t = t.replace(old, new, 1)

p = os.path.join(d, os.path.basename(rel)).replace("\\", "/")
with open(p, "w", encoding="utf-8", newline="\n") as f:
    f.write(t)
oj = os.path.join(d, "overlay.json")
with open(oj, "w", encoding="utf-8", newline="\n") as f:
    json.dump({"Replace": {(REPO + "/" + rel).replace("\\", "/"): p}}, f, indent=1)
print("OVERLAY=" + oj.replace("\\", "/"))
