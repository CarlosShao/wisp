import json
import os
import subprocess
import sys

ROOT = r"D:/work/workspace/projects plans/Wisp"
OUT = os.path.join(ROOT, ".scratch", "wisp", "probes", "221", "v1", "mutations")

t = subprocess.run(["git", "-C", ROOT, "cat-file", "blob", "HEAD:internal/tools/task.go"],
                   capture_output=True, check=True).stdout.decode("utf-8")

LIST_LINE = "//\ttask.list    -   DEFERRED with five fields, PLAN.7 :1531\n"
# rebuild the exact header line from HEAD (anchor on its stable prefix)
lines = [l for l in t.split("\n") if l.startswith("//\ttask.list")]
if len(lines) != 1:
    sys.exit("task.list header line not unique: %r" % (lines,))
LIST_LINE = lines[0] + "\n"

CANCELLINE_REMOVED = "\t\t{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},\n"


def emit(name, src):
    p = os.path.join(OUT, name + ".go")
    with open(p, "w", encoding="utf-8", newline="\n") as f:
        f.write(src)
    with open(os.path.join(OUT, name + ".json"), "w", encoding="utf-8") as f:
        json.dump({"Replace": {
            (ROOT + "/internal/tools/task.go").replace("/", "\\"): p.replace("/", "\\")
        }}, f)
    print(name)


# m13: lift task.list's DEFERRED marker WITHOUT wiring it (the shape AC#4 forbids)
emit("m13-list-marker-lifted", t.replace(LIST_LINE, ""))

# m14: put task.cancel's DEFERRED marker back while it stays registered
emit("m14-cancel-marked-again",
     t.replace(LIST_LINE, LIST_LINE + "//\ttask.cancel  -   DEFERRED with five fields, PLAN.7 :1531\n"))
