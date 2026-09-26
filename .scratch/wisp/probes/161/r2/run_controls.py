#!/usr/bin/env python3
"""
161-r2 CONTROL DESK — fires all three rulers at three fixtures and prints one
table. Run it before believing ANY "0" in the evidence file: a ruler that cannot
ring on POS cannot certify absence on the real tree.

Fixtures live in fix/*.gofix (NOT .go -- an in-tree .go file would be picked up
by the repo-wide gofumpt gate, which is ticket 161 AC#7's own finding). They are
copied into $TEMP so go/parser sees them as ordinary .go without touching the
repo.

Usage: python run_controls.py
"""
import json
import os
import shutil
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
T = os.environ.get("TEMP", "/tmp").replace("\\", "/") + "/161r2"
FIXCTL = T + "/fixctl"

# (label, fixture source, temp subdir)
CASES = [
    ("POS  (ask-shapes, must RING)", "fix/pos_known_positive.gofix", "pos"),
    ("POS2 (risk-pkg producers, must RING ruler C)", "fix/pos2_risk_pkg_producer.gofix", "pos2"),
    ("NEG  (comments+strings only, must be ZERO)", "fix/neg_comment_and_string.gofix", "neg"),
]


def go_ruler(script, root):
    p = subprocess.run(["go", "run", script, root], capture_output=True,
                       text=True, cwd=T)
    try:
        return json.loads(p.stdout), p.stderr[:200]
    except Exception:
        return None, (p.stdout[:400] + p.stderr[:400])


def main():
    # Each invocation gets its OWN fresh run dir. Nothing is ever deleted (the
    # dispatch's create-only rule), and reusing a fixed dir would silently
    # double-count: the previous run's copy of the same fixture is still there.
    run = "%s/run-%s" % (FIXCTL, time.strftime("%Y%m%d-%H%M%S"))
    print("run dir: %s   (previous run dirs are kept, not deleted)" % run)
    for label, src, sub in CASES:
        root = os.path.join(run, sub)
        os.makedirs(root, exist_ok=True)
        shutil.copyfile(os.path.join(HERE, src),
                        os.path.join(root, os.path.basename(src).replace(".gofix", ".go")))
        print("=" * 74)
        print("CONTROL %s\n  fixture=%s  as=%s/%s.go" % (label, src, root,
                                                        os.path.basename(src).replace(".gofix", "")))
        b, err = go_ruler("ruler_b2.go", root)
        if b is None:
            print("  ruler B : UNPARSEABLE OUTPUT ->", err)
        else:
            print("  ruler B : func_nodes=%d parse_errors=%d" %
                  (b["func_nodes"], b["parse_error_count"]))
            print("            method_uses=%s" % json.dumps(b["method_uses"]))
            print("            value_uses =%s" % json.dumps(b["value_uses"]))
            print("            dispatch_calls_production=%s" %
                  json.dumps(b["dispatch_calls_production"]))
            if b["parse_errors"]:
                print("            PARSE ERRORS:",
                      [s["text"][:100] for s in b["parse_errors"]])
        c, err = go_ruler("ruler_c.go", root)
        if c is None:
            print("  ruler C : UNPARSEABLE OUTPUT ->", err)
        else:
            print("  ruler C : total_rows=%d buckets=%s" %
                  (c["total_rows"], json.dumps(c["by_name_position_bucket"])))
        p = subprocess.run([sys.executable, "ruler_a_enum.py", root,
                            "--ext", ".go", "--fixtures"],
                           capture_output=True, text=True, cwd=HERE)
        tot = {}
        for l in p.stdout.split("\n"):
            if l.startswith("== "):
                w = l.split()
                tot[w[1]] = w[2].split("=")[1] + "/" + w[3].split("=")[1]
        print("  ruler A : name=total/prod  " + json.dumps(tot))
    print("=" * 74)
    print("PASS rule: POS and POS2 must each show >=1 non-zero from EVERY ruler that")
    print("claims to measure that layer; NEG must show zero AND parse_errors=0.")
    print("If NEG shows zero while parse_errors>0, the zero is VACUOUS -- fix the")
    print("fixture, not the ruler.")


if __name__ == "__main__":
    main()
