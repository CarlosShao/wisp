#!/usr/bin/env python3
"""Ticket 156 r2 (AC#5 + AC#4(ii)) mutation driver.

This file adds NO new mutation literals. Every spec it runs is imported from
probes/156/my156.py, the ruler the shipped bytes were measured with, so an r2
reading and an r1 log describe the same edit. What this driver adds is only the
"landing must be printed" self-proof the r2 dispatch demands:

  * before the run it prints, per overlaid file: the sha256 of the tree bytes it
    started from, the sha256 of the file handed to the compiler through
    -overlay, whether the mutation's own new text is in that file, and whether
    the old text is gone. my156.apply_overlay already refuses to emit a reading
    when those two checks fail; here the same facts go to the log instead of
    only to a FATAL, so a green can be read as green-because-it-landed.
  * for an asis cell it prints overlay=NONE plus the tree sha256, which is the
    other half of the same claim (re-installed bytes, not a remembered build).
  * it never passes -cover* with -overlay (my156.run FATALs on that), never uses
    -race, and judges red only from anchored '^--- FAIL:' lines.
  * after every run it re-hashes the tree file and prints it again, so a run
    that silently mutated the working tree cannot be reported as an overlay run.

usage: python my156-r2.py <log-dir> <cell>...
"""
import hashlib
import json
import os
import re
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import my156 as R  # the shipped ruler: specs, case_off, run(), the cell literals

REPO = R.REPO
SEL_SLO_ALL = "TestSLO144|TestSLO147|TestSLO149|TestSLO152|TestSLO156"
SEL_156_SHAPE = "TestSLO149Exited|TestSLO156"   # same selector as r1's mut-156 logs
SEL_152_ONLY = "TestSLO144|TestSLO147|TestSLO149|TestSLO152"

ALL_NEW_CASES_OFF = [R.case_off(R.NEW_TEST_REL, f, tag)
                     for tag, f in sorted(R.NEW_CASES.items())]

G1_RESTORED_BORROWED = R.spec(
    "cmd/wisp/slo_windows.go", "\treturn offsetUnknown\n}", "\treturn 0\n}")
CASE14 = "TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne"

CELLS = {
    # --- AC#4(ii): the reading that lives in the 144 file's case 14 comment,
    # --- re-measured at THIS anchor on THESE bytes (both directions).
    "r2-152g1-case14off": ([G1_RESTORED_BORROWED,
                            R.case_off(R.OLD_TEST_REL, CASE14, "case14")], SEL_152_ONLY),
    "r2-152g1-case14live": ([G1_RESTORED_BORROWED], SEL_152_ONLY),
    # --- AC#5: question 一, three ways, and question 二.
    "r2-asis-same-shape": ([], SEL_156_SHAPE),
    "r2-m1-same-shape": ([R.OS_ARM_REMOVED], SEL_156_SHAPE),
    "r2-m1-old-surface-only": ([R.OS_ARM_REMOVED], SEL_152_ONLY),
    "r2-asis-slo-all": ([], SEL_SLO_ALL),
    "r2-m1-slo-all": ([R.OS_ARM_REMOVED], SEL_SLO_ALL),
    "r2-m1-plus-all-new-cases-off": ([R.OS_ARM_REMOVED] + ALL_NEW_CASES_OFF, SEL_SLO_ALL),
    "r2-reinstall-asis-slo-all": ([], SEL_SLO_ALL),
    # --- AC#5 question 二, at the operator layer: the shipped probe prints the
    # --- give-up sentence verbatim, so the same probe run twice (shipped bytes /
    # --- ingredient removed) is a printed before-after pair, not an inference
    # --- from an assertion that did not fire.
    "r2-probe-on-shipped": ([R.PROBE_ADD], "^" + R.SELPROBE + "$"),
    "r2-probe-on-m1": ([R.OS_ARM_REMOVED, R.PROBE_ADD], "^" + R.SELPROBE + "$"),
    # --- named baselines. my156.py's plain cell name "asis" writes asis.log no
    # --- matter which -run selector is in play, so two baselines that answer two
    # --- different questions land in one file and the earlier one is gone. These
    # --- names carry the selector in them; r2 cites only these.
    "r2-asis-152-surface": ([], SEL_152_ONLY),
    "r2-asis-slo-all-after-ac4": ([], SEL_SLO_ALL),
    "r2-m1-slo-all-after-ac4": ([R.OS_ARM_REMOVED], SEL_SLO_ALL),
    "r2-m1-same-shape-after-ac4": ([R.OS_ARM_REMOVED], SEL_156_SHAPE),
}


def sha(path):
    with open(path, "rb") as fh:
        return hashlib.sha256(fh.read()).hexdigest()[:16]


SHIPPED = ["cmd/wisp/slo_windows.go", "cmd/wisp/slo_report_144_windows_test.go",
           "cmd/wisp/slo_exit_os_156_windows_test.go"]


def shipped_lines(name):
    """Every cell prints the bytes it started from, whether or not it overlays:
    an asis reading is only worth what it is worth on a named version."""
    return ["P156R2|TREE|cell=%s|head=%s|%s" % (
        name, git("rev-parse", "--short", "HEAD"),
        "|".join("%s=%s" % (os.path.basename(s), sha(os.path.join(REPO, s)))
                 for s in SHIPPED))]


def git(*args):
    return subprocess.run(["git"] + list(args), cwd=REPO, capture_output=True,
                          text=True, encoding="utf-8").stdout.strip()


def landing_lines(base, name, entries):
    """Print what the compiler is about to be handed - or refuse."""
    out = []
    rels = [] if not entries else sorted({e[1] for e in entries if e[0] == "mut"})
    out.append("P156R2|LANDING|cell=%s|overlay=%s" % (
        name, "NONE(asis: tree bytes as shipped)" if not entries else "yes"))
    for rel in rels:
        tree = os.path.join(base, rel.replace("/", os.sep))
        out.append("P156R2|LANDING|cell=%s|tree_file=%s|tree_sha256=%s" % (name, rel, sha(tree)))
    if entries:
        overlay, replaced = R.apply_overlay(base, name, entries)
        out.append("P156R2|LANDING|cell=%s|overlay_json=%s" % (name, overlay))
        mapping = json.load(open(overlay, encoding="utf-8"))["Replace"]
        for key in replaced:
            rel = key.split("cmd/wisp/")[-1]
            dst = mapping[key]  # the exact path the compiler is told to read
            body = open(dst, encoding="utf-8", newline="").read()
            mine = [e for e in entries if e[0] == "mut" and e[1].endswith(rel)]
            newin = ["new%d=%s" % (i + 1, "yes" if e[3] in body else "NO-REFUSE")
                     for i, e in enumerate(mine)]
            oldin = ["old%d=%s" % (i + 1, "STILL-THERE-REFUSE" if e[2] in body else "gone")
                     for i, e in enumerate(mine)]
            out.append("P156R2|LANDING|cell=%s|compiled_file=cmd/wisp/%s|src_read_by_compiler=%s"
                       "|entries=%d|tmp_sha256=%s|%s|%s" % (
                           name, rel, dst, len(mine), sha(dst), "|".join(newin), "|".join(oldin)))
    return out


def roster(path):
    txt = open(path, encoding="utf-8", newline="").read()
    return sorted(re.findall(r"^--- (PASS|FAIL|SKIP): (\S+)", txt, re.M))


def main():
    if len(sys.argv) < 3:
        sys.exit(__doc__)
    log_dir = sys.argv[1]
    cells = sys.argv[2:]
    os.makedirs(log_dir, exist_ok=True)
    base = REPO
    head = git("rev-parse", "--short", "HEAD")
    print("P156R2|ANCHOR|head=%s|tree=%s|status_cmd_wisp=[%s]" % (
        head, base, git("status", "--porcelain", "--", "cmd/wisp")))
    summary = []
    for name in cells:
        if name not in CELLS:
            sys.exit("unknown cell %s (known: %s)" % (name, ",".join(sorted(CELLS))))
        entries, sel = CELLS[name]
        for ln in shipped_lines(name) + landing_lines(base, name, entries):
            print(ln)
        r = R.run(base, name, entries, log_dir, sel)
        lp = r["log"]
        after = [git("status", "--porcelain", "--", "cmd/wisp")]
        txt = open(lp, encoding="utf-8", newline="").read()
        counts = {
            "RUN": len(re.findall(r"^=== RUN", txt, re.M)),
            "PASS": len(re.findall(r"^--- PASS", txt, re.M)),
            "FAIL": len(re.findall(r"^--- FAIL", txt, re.M)),
            "SKIP": len(re.findall(r"^--- SKIP", txt, re.M)),
        }
        with open(os.path.join(log_dir, name + ".roster.txt"), "w", encoding="utf-8",
                  newline="") as fh:
            for kind, t in roster(lp):
                fh.write("%s %s\n" % (kind, t))
        fails = re.findall(r"^--- FAIL: (\S+)", txt, re.M)
        print("P156R2|RESULT|cell=%s|rc=%d|RUN=%d|PASS=%d|FAIL=%d|SKIP=%d|red=%s"
              "|tree_dirty_after=%s|log=%s" % (
                  name, r["rc"], counts["RUN"], counts["PASS"], counts["FAIL"],
                  counts["SKIP"], ",".join(fails) or "-", after[0] or "clean", lp))
        summary.append("%s RUN=%d PASS=%d FAIL=%d SKIP=%d red=%s" % (
            name, counts["RUN"], counts["PASS"], counts["FAIL"], counts["SKIP"],
            ",".join(fails) or "-"))
    with open(os.path.join(log_dir, "r2-summary.txt"), "a", encoding="utf-8",
              newline="") as fh:
        fh.write("head=%s\n" % head + "".join(s + "\n" for s in summary))
    print("P156R2|DONE|" + ("; ".join(summary)))


if __name__ == "__main__":
    main()
