#!/usr/bin/env python3
"""180-c1: split census.tsv into the three reader classes and write reader-classes.txt."""
import os, csv

REPO = r"D:/work/workspace/projects plans/Wisp"
tsv = os.path.join(REPO, ".scratch", "wisp", "probes", "180", "c1", "census.tsv")
dest = os.path.join(REPO, ".scratch", "wisp", "probes", "180", "c1", "reader-classes.txt")

rows = list(csv.DictReader(open(tsv, encoding="utf-8"), delimiter="\t"))
rows = [r for r in rows if r["field"] != "UNPARSED"]

A = [r for r in rows if int(r["hits_all"]) == 0]
B = [r for r in rows if int(r["hits_all"]) > 0 and int(r["hits_outside_cfg"]) == 0]
AMB = [r for r in rows if int(r["hits_outside_cfg"]) > 0 and int(r["hits_attributed"]) == 0]
C = [r for r in rows if int(r["hits_attributed"]) > 0]
GEN = [r for r in rows if r["collision_flag"]]

out = []
out.append("# 180-c1 AC#4 reader classes (census of default-tagged fields in internal/config/schema.go)")
out.append("# ruler A: grep -rn \"\\.<Field>\\b\" --include=*.go internal/ cmd/ | grep -v _test.go   (the ticket's own ruler)")
out.append("# ruler B: same, then drop internal/config/** hits  (does anything OUTSIDE the schema's own package read it?)")
out.append("# population: 69 default-tagged fields (= 70 lines containing 'default:' minus schema.go:8, a doc comment)")
out.append("")
out.append("## CLASS A - zero hits under ruler A : %d" % len(A))
out.append("##   meaning: no non-test line in internal/ or cmd/ mentions .<Field> at all")
for r in A:
    out.append("  %-24s %-18s default=%-14s toml=%s" % (r["field"], r["struct"], r["default"], r["toml_key"]))
out.append("")
out.append("## CLASS B - ruler A > 0 but ruler B == 0 : %d" % len(B))
out.append("##   meaning: only internal/config itself touches it (manager.go block copies, defaults.go reflection,")
out.append("##            catalog.go). No component outside the schema package reads the value and acts on it.")
for r in B:
    out.append("  %-24s %-18s default=%-14s hits_all=%-3s sample=%s"
               % (r["field"], r["struct"], r["default"], r["hits_all"], r["sample_all"]))
out.append("")
out.append("## CLASS C - ATTRIBUTED reader exists outside internal/config : %d" % len(C))
out.append("##   meaning: a production line outside the schema package names BOTH the owning type AND the field.")
out.append("##   THIS is the only class 180-c1 will call 'has a production reader'.")
for r in C:
    out.append("  %-24s %-18s default=%-14s hits_all=%-4s hits_outside=%-3s hits_attr=%-3s sample=%s%s"
               % (r["field"], r["struct"], r["default"], r["hits_all"], r["hits_outside_cfg"],
                  r["hits_attributed"], r["sample_outside_cfg"],
                  "  [GENERIC-NAME]" if r["collision_flag"] else ""))
out.append("")
out.append("## CLASS AMBIG - hits outside internal/config but ATTRIBUTION UNPROVEN : %d" % len(AMB))
out.append("##   The name appears outside the schema package, never on a line that also names the owning type.")
out.append("##   Two readings are possible and this leg does NOT pick one: (a) pure name collision with another")
out.append("##   struct's field, (b) a real reader that copies the value into a local/other struct first, so the")
out.append("##   section name and field name never share a line. Ruler C cannot tell these apart. NOT counted as")
out.append("##   'has reader', NOT counted as dead key. Each needs a code read by the landing leg.")
for r in AMB:
    out.append("  %-24s %-18s hits_all=%-4s hits_outside=%-3s first=%s"
               % (r["field"], r["struct"], r["hits_all"], r["hits_outside_cfg"], r["sample_outside_cfg"]))
out.append("")
out.append("## totals: A=%d B=%d AMBIG=%d C=%d  sum=%d" % (len(A), len(B), len(AMB), len(C), len(A)+len(B)+len(AMB)+len(C)))
out.append("rc=0")

open(dest, "w", encoding="utf-8", newline="\n").write("\n".join(out) + "\n")
print("wrote %s" % dest)
print("A=%d B=%d C=%d generic=%d total=%d" % (len(A), len(B), len(C), len(GEN), len(rows)))
print("rc=0")
