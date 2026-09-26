#!/usr/bin/env python3
"""Pretty-print / assert on ruler B's JSON (log_rulerB_raw.json)."""
import collections
import json
import os
import sys

p = os.path.join(os.path.dirname(os.path.abspath(__file__)), "log_rulerB_raw.json")
d = json.load(open(p, encoding="utf-8"))

BUCKETS = ["use-in-cmd", "use-in-internal", "use-in-panel", "use-in-test"]

print("=== 1. enum members DECLARED (definition side, from go/parser) ===")
for s in d["enum_decls"]:
    print("   %-28s %-18s %s:%d" % (s["text"], s["kind"], s["file"], s["line"]))

print("\n=== 2. interface method sets carrying an ask verb ===")
for k, v in sorted(d["iface_sigs"].items()):
    if set(v) & {"PendingWindow", "PendingApproval", "Prompt", "Decide", "Update"}:
        print("   %-16s %s" % (k, sorted(v)))

print("\n=== 3. selector USE sites of the wanted family, bucketed ===")
for tok, rows in sorted(d["use_sites"].items()):
    c = collections.Counter(r["kind"] for r in rows)
    line = "   %-18s total=%-4d " % (tok, len(rows))
    line += "  ".join("%s=%d" % (b, c.get(b, 0)) for b in BUCKETS)
    print(line)

print("\n=== 3b. non-test use sites, spelled out ===")
for tok, rows in sorted(d["use_sites"].items()):
    prod = [r for r in rows if r["kind"] in ("use-in-cmd", "use-in-internal", "use-in-panel")]
    print("   -- %s : %d" % (tok, len(prod)))
    for r in sorted(prod, key=lambda x: (x["file"], x["line"])):
        print("      %s:%d [%s] %s" % (r["file"], r["line"], r["kind"], r["text"]))

print("\n=== 4. call-graph reachability from production seeds ===")
r = d["reach"]
print("   seeds:", ", ".join(r["_seeds"]))
print("   production nodes reachable from seeds:", r["_reachable_total"])
for tok in sorted(r):
    if tok.startswith("_"):
        continue
    v = r[tok] or []
    print("   %-18s reachable=%d  %s" % (tok, len(v), v))

print("\n=== 5. statemachine.Machine.Dispatch declaration sites ===")
for s in d["dispatchers"]:
    print("   %s:%d %s" % (s["file"], s["line"], s["text"]))
