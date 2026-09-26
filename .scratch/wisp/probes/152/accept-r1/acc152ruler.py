#!/usr/bin/env python3
"""R1 acceptance ruler for ticket 152 (independent of probes/152/my152.py).

For every cell: copy the named file out of a snapshot, apply literal
replacements, READ BACK and assert new text present / old text absent
(otherwise FATAL, no reading printed), write an overlay json, run
`go test -overlay ... -v -run <selection>` with NO -cover flag, and report
the four counts (top-level === RUN / --- PASS / --- FAIL / --- SKIP) plus the
red names.

usage: acc152ruler.py <snap>:<cell> [<snap>:<cell> ...]     snap in {post,pre}
"""
import json
import os
import re
import shlex
import subprocess
import sys

BASE = "D:/work/tmp/wisp152-accept-r1"
SNAP = {
    "post": BASE + "/snap-post",          # git archive 97cfc6e  (delivered bytes)
    "pre": BASE + "/snap-pre",            # git archive 10e3585^ (pre-fix bytes)
}
MUTDIR = BASE + "/mut"
OVLDIR = BASE + "/ovl"
LOGDIR = BASE + "/logs"

# TestAcc152Helper is deliberately NOT in this selection: it calls os.Exit(7).
SELECTION = "TestSLO144|TestSLO147|TestSLO149|TestSLO152|TestAcc152Probe"

EXITED_ARM = (
    "\t\tif s.exited() {\n"
    "\t\t\treturn nil, fmt.Errorf(\"wisp slo: subject %d exited (code %d) without writing its report (%s)\",\n"
    "\t\t\t\ts.pid, s.exitCode(), last.summary())\n"
    "\t\t}\n"
)
CASE13_OFF = (
    "func TestSLO149ExitedGiveUpSentenceCarriesTheLastReading(t *testing.T) {",
    "func notATest_152acc_case13_neutralised(t *testing.T) {",
)
CASE14_OFF = (
    "func TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne(t *testing.T) {",
    "func notATest_152acc_case14_neutralised(t *testing.T) {",
)
SLO = "cmd/wisp/slo_windows.go"
TST = "cmd/wisp/slo_report_144_windows_test.go"
PROBE = "cmd/wisp/zz152acc1_windows_test.go"
PROBE_SRC = BASE + "/probe/zz152acc1_windows_test.go"

CELLS = {
    # --- baselines -------------------------------------------------------
    "asis-post": [],
    "asis-pre": [],
    "probe-post": [("@PROBE@",)],
    # --- the exited arm: does anything notice if it does not exist? ------
    "x1-exited-arm-deleted": [(SLO, EXITED_ARM, "")],
    "x2-exited-arm-deleted-case13off": [(SLO, EXITED_ARM, ""), (TST,) + CASE13_OFF],
    "x3-case13-off": [(TST,) + CASE13_OFF],
    "z1-exited-loses-last-summary": [
        (SLO, "s.pid, s.exitCode(), last.summary())", 's.pid, s.exitCode(), "nothing read")')
    ],
    "z2-exited-loses-summary-case13off": [
        (SLO, "s.pid, s.exitCode(), last.summary())", 's.pid, s.exitCode(), "nothing read")'),
        (TST,) + CASE13_OFF,
    ],
    # --- AC#3's fallback leg --------------------------------------------
    "y1-fallback-returns-0": [(SLO, "\treturn offsetUnknown\n}", "\treturn 0\n}")],
    "y2-fallback-0-case14off": [
        (SLO, "\treturn offsetUnknown\n}", "\treturn 0\n}"),
        (TST,) + CASE14_OFF,
    ],
    "y3-case14-off": [(TST,) + CASE14_OFF],
    "y4-renderer-drops-the-naming-arm": [
        (SLO,
         "\tif offset < 0 {\n"
         "\t\treturn fmt.Errorf(\"%d bytes contradict a subject report at an offset the decoder did not name: %w\", size, err)\n"
         "\t}\n", "")
    ],
    "y5-delete-the-fallback-line": [(SLO, "\treturn offsetUnknown\n}", "}")],
    "y6-fallback-returns-1": [(SLO, "\treturn offsetUnknown\n}", "\treturn 1\n}")],
    # --- on the PRE-fix bytes: was the borrowed value witnessed at all? ------
    "p1-pre-fallback-was-a-constant": [(SLO, "\treturn inputOffset\n}", "\treturn 0\n}")],
    "p2-pre-fallback-was-minus-one": [(SLO, "\treturn inputOffset\n}", "\treturn -1\n}")],
    # --- the delivered test file dropped onto the pre-fix code (n1 shape) ----
    "n1-newtest-on-pre-code": [("@MAP@", TST, "post")],
}


def build_overlay(snap_key, cell, muts):
    """muts: list of (rel, old, new) or the marker ("@PROBE@",)."""
    snap = SNAP[snap_key]
    files = {}
    originals = {}
    extra = {}
    for m in muts:
        if m[0] == "@PROBE@" and len(m) == 1:
            extra[snap + "/" + PROBE] = PROBE_SRC
            continue
        if m[0] == "@MAP@":
            extra[snap + "/" + m[1]] = SNAP[m[2]] + "/" + m[1]
            continue
        rel, old, new = m
        p = os.path.join(snap, rel.replace("/", os.sep))
        if rel not in originals:
            with open(p, "r", encoding="utf-8", newline="") as fh:
                originals[rel] = fh.read()
        cur = files.get(rel, originals[rel])
        n = cur.count(old)
        if n != 1:
            return None, None, "FATAL %s: pattern for %s matched %d times (want exactly 1)" % (cell, rel, n)
        files[rel] = cur.replace(old, new, 1)
    for m in muts:
        if len(m) != 3 or m[0].startswith("@"):
            continue
        rel, old, new = m
        txt = files[rel]
        if old in txt:
            return None, None, "FATAL %s: old text still present in %s" % (cell, rel)
        if new and new not in txt:
            return None, None, "FATAL %s: new text missing from %s" % (cell, rel)
    replace = dict(extra)
    for rel, txt in files.items():
        out = os.path.join(MUTDIR, snap_key + "--" + cell, rel.replace("/", os.sep))
        os.makedirs(os.path.dirname(out), exist_ok=True)
        with open(out, "w", encoding="utf-8", newline="") as fh:
            fh.write(txt)
        replace[(snap + "/" + rel).replace("\\", "/")] = out.replace("\\", "/")
    ovl = os.path.join(OVLDIR, "%s--%s.json" % (snap_key, cell))
    os.makedirs(OVLDIR, exist_ok=True)
    with open(ovl, "w", encoding="utf-8") as fh:
        json.dump({"Replace": {k.replace("\\", "/"): v for k, v in replace.items()}}, fh)
    return ovl, len(files) + len(extra), None


def go_test(snap_key, cell, ovl, nfiles, run_sel, log_name, selection_default=True):
    env = dict(os.environ)
    # Windows-native PATH form (drive letter + backslashes, `;` separated): the MSYS
    # form only works when bash rewrites it, and this ruler launches natively.
    dlls = (SNAP[snap_key] + "/third_party/sherpa-onnx").replace("/", "\\")
    env["PATH"] = dlls + ";" + env["PATH"]
    cmd = ["go", "test", "-count=1", "-v"]
    if ovl:
        cmd += ["-overlay", ovl]
    cmd += ["-run", run_sel, "./cmd/wisp/"]
    assert not any("cover" in c for c in cmd), "the overlay+cover trap is barred from this ruler"
    os.makedirs(LOGDIR, exist_ok=True)
    log = os.path.join(LOGDIR, log_name + ".log")
    with open(log, "w", encoding="utf-8", newline="\n") as fh:
        fh.write("$ " + " ".join(shlex.quote(c) for c in cmd) + "\n")
        fh.write("cwd=" + SNAP[snap_key] + "\n")
        fh.write("overlay_files=%d\n" % (nfiles or 0))
        if ovl:
            fh.write("overlay=" + open(ovl, encoding="utf-8").read() + "\n")
        p = subprocess.run(cmd, cwd=SNAP[snap_key], env=env, stdout=subprocess.PIPE,
                           stderr=subprocess.STDOUT, text=True, encoding="utf-8", errors="replace")
        fh.write(p.stdout)
    text = open(log, encoding="utf-8").read()
    lines = text.splitlines()
    top_run = [l for l in lines if re.match(r"^=== RUN   \S+$", l)]
    tp = [l for l in lines if re.match(r"^--- PASS:", l)]
    tf = [l for l in lines if re.match(r"^--- FAIL:", l)]
    ts = [l for l in lines if re.match(r"^--- SKIP:", l)]
    summary = dict(cell=log_name, rc=p.returncode, RUN=len(top_run),
                   DID_NOT_RUN=(len(top_run) == 0), PASS=len(tp),
                   FAIL=len(tf), SKIP=len(ts),
                   red=sorted(re.sub(r"^--- FAIL: (\S+).*", r"\1", l) for l in tf))
    with open(os.path.join(LOGDIR, "summary.txt"), "a", encoding="utf-8") as fh:
        fh.write(json.dumps(summary, ensure_ascii=False) + "\n")
    return summary


def main():
    for key in sys.argv[1:]:
        snap_key, cell = key.split(":", 1)
        if cell not in CELLS:
            print(json.dumps(dict(cell=key, fatal="unknown cell")))
            continue
        ovl, n, err = build_overlay(snap_key, cell, CELLS[cell])
        if err:
            print(json.dumps(dict(cell=key, fatal=err)))
            continue
        print(json.dumps(go_test(snap_key, cell, ovl, n, SELECTION, key.replace(":", "__")), ensure_ascii=False))


if __name__ == "__main__":
    main()
