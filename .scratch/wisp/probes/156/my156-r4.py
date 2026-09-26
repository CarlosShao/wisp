#!/usr/bin/env python3
"""Ticket 156 r4 (AC#4(ii)) mutation driver.

This程 is the FOURTH program on ticket 156. r3 died mid-flight (model-service
connection interruption) after delivering the gofumpt fix and the "before" gate,
so AC#4(ii) is r4's cell. AC#4(ii) says "re-measure at YOUR OWN anchor before you
change it", so this file exists to produce r4's own readings.

Instrument discipline inherited, not re-invented:

  * This driver adds NO new mutation literals. Every spec it runs is imported from
    probes/156/my156.py - the ruler the shipped bytes were measured with - so an
    r4 reading, an r2 reading and an r1 log describe the SAME edit. The cell below
    is the same G1-restored + case-14-renamed pair r2 ran, re-typed here on
    purpose: my156-r2.py is another program's deliverable and is not executed by
    this程.
  * probes/152/my152.py is another cell's unfinished deliverable (it is ` M` on
    disk): not executed, not restored, not committed here.
  * -overlay is never combined with -cover* (my156.run FATALs on that), no -race,
    and red comes ONLY from anchored '^--- FAIL:' lines.
  * before each run it prints, per overlaid file: the sha256 of the tree bytes it
    started from, the sha256 of the file handed to the compiler, whether the
    mutation's own new text is in that file and whether the old text is gone - so
    a green can be read as green-because-it-landed, not green-because-the-overlay
    was silently ignored.
  * after every run it re-prints the tree sha256 and `git status -- cmd/wisp` so a
    run that silently edited the working tree cannot be reported as an overlay run.
  * the log names carry the selector in them. my156.py's bare cell name "asis"
    writes asis.log whatever -run was in play, which is how r2 ended up with one
    file holding two different baselines; r4 never writes that name.

usage: python my156-r4.py <log-dir> <cell>...
"""
import hashlib
import json
import os
import re
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import my156 as R  # the shipped ruler: specs, case_off, run(), apply_overlay

REPO = R.REPO
# r1/r2's "票 152 surface" selector: the four pre-ticket-156 SLO cases, WITHOUT
# this ticket's new TestSLO156 cases. The comment AC#4(ii) fixes promises the
# reading of THIS selector, so this is the only selector r4's AC#4(ii) cells use.
SEL_152_ONLY = "TestSLO144|TestSLO147|TestSLO149|TestSLO152"

# G1 "put the borrowed value back": the arm that returns offsetUnknown instead of a
# number is reverted to the number. Verbatim the spec my156.py ships.
G1_RESTORED_BORROWED = R.spec(
    "cmd/wisp/slo_windows.go", "\treturn offsetUnknown\n}", "\treturn 0\n}")
CASE14 = "TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne"

CELLS = {
    # the reading that belongs in the sentence at :739 - borrowed value restored
    # AND case 14 neutralised, which is exactly what that sentence describes.
    "r4-152g1-case14off": ([G1_RESTORED_BORROWED,
                            R.case_off(R.OLD_TEST_REL, CASE14, "case14")], SEL_152_ONLY),
    # the same edit with case 14 live: the pair that makes the sentence mean
    # something. Expected red, and the red name expected to be case 14 itself.
    "r4-152g1-case14live": ([G1_RESTORED_BORROWED], SEL_152_ONLY),
    # no mutation at all, same selector, same bytes: the baseline the stale
    # "21/14/0" actually is, if the错因 r2 registered is real.
    "r4-asis-152-surface": ([], SEL_152_ONLY),
}

SHIPPED = ["cmd/wisp/slo_windows.go", "cmd/wisp/slo_report_144_windows_test.go",
           "cmd/wisp/slo_exit_os_156_windows_test.go"]


def sha(path):
    with open(path, "rb") as fh:
        return hashlib.sha256(fh.read()).hexdigest()[:16]


def git(*args):
    return subprocess.run(["git"] + list(args), cwd=REPO, capture_output=True,
                          text=True, encoding="utf-8").stdout.strip()


def shipped_lines(name):
    """Every cell prints the bytes it started from, whether or not it overlays:
    an asis reading is only worth what it is worth on a named version."""
    return ["P156R4|TREE|cell=%s|head=%s|%s" % (
        name, git("rev-parse", "--short", "HEAD"),
        "|".join("%s=%s" % (os.path.basename(s), sha(os.path.join(REPO, s)))
                 for s in SHIPPED))]


def landing_lines(base, name, entries):
    """Print what the compiler is about to be handed - or refuse."""
    out = []
    rels = [] if not entries else sorted({e[1] for e in entries if e[0] == "mut"})
    out.append("P156R4|LANDING|cell=%s|overlay=%s" % (
        name, "NONE(asis: tree bytes as shipped)" if not entries else "yes"))
    for rel in rels:
        tree = os.path.join(base, rel.replace("/", os.sep))
        out.append("P156R4|LANDING|cell=%s|tree_file=%s|tree_sha256=%s" % (name, rel, sha(tree)))
    if entries:
        overlay, replaced = R.apply_overlay(base, name, entries)
        out.append("P156R4|LANDING|cell=%s|overlay_json=%s" % (name, overlay))
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
            out.append("P156R4|LANDING|cell=%s|compiled_file=cmd/wisp/%s|src_read_by_compiler=%s"
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
    print("P156R4|ANCHOR|head=%s|branch=%s|tree=%s|status_cmd_wisp=[%s]" % (
        head, git("branch", "--show-current"), base, git("status", "--porcelain", "--", "cmd/wisp")))
    # contention check: `wisp slo`-shaped runs live on this machine's own
    # self-hosted runner, so a same-machine CPU fight is a reading, not noise.
    tl = subprocess.run(["tasklist"], capture_output=True, text=True).stdout or ""
    print("P156R4|CONTENTION|runner_listener_procs=%d|runner_worker_procs=%d" % (
        tl.count("Runner.Listener.exe"), tl.count("Runner.Worker.exe")))
    summary = []
    for name in cells:
        if name not in CELLS:
            sys.exit("unknown cell %s (known: %s)" % (name, ",".join(sorted(CELLS))))
        entries, sel = CELLS[name]
        for ln in shipped_lines(name) + landing_lines(base, name, entries):
            print(ln)
        r = R.run(base, name, entries, log_dir, sel)
        lp = r["log"]
        after = git("status", "--porcelain", "--", "cmd/wisp")
        txt = open(lp, encoding="utf-8", newline="").read()
        counts = {k: len(re.findall(p, txt, re.M)) for k, p in (
            ("RUN", "^=== RUN"), ("PASS", "^--- PASS"),
            ("FAIL", "^--- FAIL"), ("SKIP", "^--- SKIP"))}
        with open(os.path.join(log_dir, name + ".roster.txt"), "w", encoding="utf-8",
                  newline="") as fh:
            for kind, t in roster(lp):
                fh.write("%s %s\n" % (kind, t))
        fails = re.findall(r"^--- FAIL: (\S+)", txt, re.M)
        print("P156R4|RESULT|cell=%s|rc=%d|RUN=%d|PASS=%d|FAIL=%d|SKIP=%d|red=%s"
              "|tree_dirty_after=[%s]|log=%s" % (
                  name, r["rc"], counts["RUN"], counts["PASS"], counts["FAIL"],
                  counts["SKIP"], ",".join(fails) or "-", after or "clean", lp))
        for s in SHIPPED:
            print("P156R4|AFTER|cell=%s|%s=%s" % (name, os.path.basename(s),
                                                  sha(os.path.join(REPO, s))))
        summary.append("%s RUN=%d PASS=%d FAIL=%d SKIP=%d red=%s" % (
            name, counts["RUN"], counts["PASS"], counts["FAIL"], counts["SKIP"],
            ",".join(fails) or "-"))
    with open(os.path.join(log_dir, "r4-summary.txt"), "a", encoding="utf-8",
              newline="") as fh:
        fh.write("head=%s\n" % head + "".join(s + "\n" for s in summary))
    print("P156R4|DONE|" + ("; ".join(summary)))


if __name__ == "__main__":
    main()
