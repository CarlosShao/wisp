# Build one mutant file from a tracked source by exact single string replacement,
# plus a go test -overlay JSON so the worktree is never written to.
# Usage:
#   python mutate.py --root <repo> --orig <rel/path> --needle <file> --repl <file> \
#                    --out <rel/path under probes> --overlay <rel/path under probes>
# needle/repl are read from small files so shell quoting can never corrupt them.
import argparse
import json
import os
import sys

ap = argparse.ArgumentParser()
ap.add_argument("--root", required=True)
ap.add_argument("--orig", required=True)
ap.add_argument("--needle", required=True)
ap.add_argument("--repl", required=True)
ap.add_argument("--out", required=True)
ap.add_argument("--overlay", required=True)
a = ap.parse_args()

root = os.path.abspath(a.root)
orig = os.path.join(root, a.orig.replace("\\", "/"))
src = open(orig, "r", encoding="utf-8").read()
needle = open(a.needle, "r", encoding="utf-8").read()
repl = open(a.repl, "r", encoding="utf-8").read()

n = src.count(needle)
if n != 1:
    print("MUT-ABORT: needle occurs %d times, need exactly 1" % n)
    sys.exit(3)

mut = src.replace(needle, repl)
outp = os.path.abspath(a.out)
os.makedirs(os.path.dirname(outp), exist_ok=True)
open(outp, "w", encoding="utf-8", newline="\n").write(mut)

ov = {"Replace": {orig.replace("\\", "/"): outp.replace("\\", "/")}}
ovp = os.path.abspath(a.overlay)
open(ovp, "w", encoding="utf-8", newline="\n").write(json.dumps(ov))

# landing proof, same chain
print("MUT-OK orig=%s" % a.orig)
print("MUT-ORIG-HITS needle_hits_in_orig=%d" % src.count(needle))
print("MUT-MUTANT-HITS needle_hits_in_mutant=%d" % mut.count(needle))
print("MUT-REPL-HITS repl_hits_in_mutant=%d" % mut.count(repl))
print("MUT-LINES orig=%d mutant=%d" % (src.count(chr(10)), mut.count(chr(10))))
print("OVERLAY=%s" % ovp)
