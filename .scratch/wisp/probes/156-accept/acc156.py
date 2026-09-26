#!/usr/bin/env python3
"""accept-r1 ruler: imports ticket 156's OWN ruler (my156.py) so no mutation
literal is re-typed here -- every edit below is the exact byte sequence the
implementing programs used. Only the cell names, selectors, log dir and the
print/roster layer are mine.

usage:
  python acc156.py probe                 # AC#1/AC#2: probe on anchor + on m1
  python acc156.py teeth                 # AC#5: forms 甲/乙/丙 + their asis pairs
  python acc156.py ac4                    # AC#4(ii): the reading inside case 14's comment
"""
import os
import re
import sys

sys.path.insert(0, os.path.join(".scratch", "wisp", "probes", "156"))
import my156 as R  # noqa: E402

BASE = R.REPO
OUT = os.path.join(BASE, ".scratch", "wisp", "probes", "156-accept")
PROBE_SEL = "^" + R.SELPROBE + "$"
SEL156_AND_149 = "TestSLO149Exited|TestSLO156"
SEL_ALL_SLO = "TestSLO144|TestSLO147|TestSLO149|TestSLO152|TestSLO156"
SEL_OLD_ONLY = "TestSLO144|TestSLO147|TestSLO149|TestSLO152"


def shot(name, entries, sel):
    logdir = os.path.join(OUT, "logs")
    os.makedirs(logdir, exist_ok=True)
    return R.run(BASE, name, entries, logdir, sel)


def probe():
    # AC#1 改前 shape, measured by THIS run: the (b) arm knocked out only.
    shot("accept-probe-on-ship", [R.PROBE_ADD], PROBE_SEL)
    shot("accept-probe-on-m1", [R.PROBE_ADD, R.OS_ARM_REMOVED], PROBE_SEL)


def teeth():
    all_off = [R.case_off(R.NEW_TEST_REL, f, tag) for tag, f in sorted(R.NEW_CASES.items())]
    # 甲: knock out the single ingredient, ask the new cases + case 13.
    shot("accept-teeth-asis-jia", [], SEL156_AND_149)
    shot("accept-teeth-m1-jia", [R.OS_ARM_REMOVED], SEL156_AND_149)
    # 乙: same knockout PLUS the four new holders renamed away.
    shot("accept-teeth-m1-plus-newcases-off-yi", [R.OS_ARM_REMOVED] + all_off, SEL_ALL_SLO)
    # 丙: knockout, but only the pre-ticket roster is asked.
    shot("accept-teeth-asis-old-surface-bing-asis", [], SEL_OLD_ONLY)
    shot("accept-teeth-m1-old-surface-bing", [R.OS_ARM_REMOVED], SEL_OLD_ONLY)


def ac4():
    sel = SEL_OLD_ONLY
    g1 = R.spec("cmd/wisp/slo_windows.go", "\treturn offsetUnknown\n}", "\treturn 0\n}")
    off14 = R.case_off(R.OLD_TEST_REL,
                       "TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne", "case14")
    shot("accept-ac4-asis-152-surface", [], sel)
    shot("accept-ac4-g1-case14off", [g1, off14], sel)
    shot("accept-ac4-g1-case14live", [g1], sel)


def muts():
    """Single-point knockouts of EVERY changed spot in exitStatus/queryProcess, plus
    removing ONE holder at a time. This is the acceptor's own answer to 'which case
    goes quiet when this one spot is reverted' -- r1 reported these, r2 re-ran none."""
    off15 = R.case_off(R.NEW_TEST_REL, R.NEW_CASES["case15"], "case15")
    off17 = R.case_off(R.NEW_TEST_REL, R.NEW_CASES["case17"], "case17")
    off13 = R.case_off(R.OLD_TEST_REL, R.CASE13, "case13")
    shot("accept-mut-asis", [], SEL156_AND_149)
    shot("accept-mut-m2-never-dead", [R.OS_ARM_NEVER_DEAD], SEL156_AND_149)
    shot("accept-mut-m3-always-dead", [R.OS_ARM_ALWAYS_DEAD], SEL156_AND_149)
    shot("accept-mut-m4-arm-a-deleted", [R.SHORTCIRCUIT_DELETED], SEL156_AND_149)
    shot("accept-mut-m5-code-read-dropped", [R.CODE_READ_DROPPED], SEL156_AND_149)
    # removing exactly one holder: does the same knockout still get noticed?
    shot("accept-mut-m1-off-case15", [R.OS_ARM_REMOVED, off15], SEL156_AND_149)
    shot("accept-mut-m1-off-case17", [R.OS_ARM_REMOVED, off17], SEL156_AND_149)
    # is case 13 the holder of arm (a)?  (m4 with case 13 removed -- r1's m4-plus-case13-off)
    shot("accept-mut-m4-off-case13", [R.SHORTCIRCUIT_DELETED, off13], SEL156_AND_149)


if __name__ == "__main__":
    mode = sys.argv[1]
    print("ACCEPT156|MODE|%s" % mode)
    {"probe": probe, "teeth": teeth, "ac4": ac4, "muts": muts}[mode]()
