#!/usr/bin/env python3
"""Ticket 156 rulers (probe runner + mutation ruler).

Inherits the three instrument lessons this repo has already paid for, and the
two pitfalls the 156 dispatch names:

  * "did it even run" is decided ONLY by the anchored `=== RUN` count. A load-time
    0xc0000135 (missing/mangled sherpa DLL path) exits 1 with zero RUN lines, which
    is NOT green: this script emits `LOADED=NO` and returns rc=3 for that case.
  * -overlay is NEVER combined with a -cover* flag (Go silently ignores the
    overlay then, so a mutated file with a syntax error still reports `ok`). Any
    go-test argument containing "cover" is refused here.
  * for every mutation, the overlay's OWN bytes are re-read before the run and
    both directions are checked (new text present, old text absent); a spec whose
    literal hits != 1 is a FATAL, not a silently-skipped cell.
  * red/green comes from anchored `^--- FAIL:` lines only, never from a summary
    line or a -skip in prose.

usage:
  python my156.py probe <base-tree> <out-dir>            # AC#1 shot, both trees
  python my156.py mut   <base-tree> <out-dir> <cell>...  # AC#2/AC#5 teeth
"""
import json
import os
import re
import shlex
import subprocess
import sys

REPO = "D:/work/workspace/projects plans/Wisp"
HERE = os.path.dirname(os.path.abspath(__file__))
GO = "go"
DLL_DIRNAME = os.path.join("third_party", "sherpa-onnx")

SEL156 = "TestSLO156"
SEL144 = "TestSLO144|TestSLO147|TestSLO149|TestSLO152"
SELPROBE = "TestP156SubjectDeathProbe"


def spec(rel, old, new):
    return ("mut", rel, old, new)


# ------------------------------------------------------------------ overlays ---
# The AC#1 probe: a file that does NOT exist in the tree, added by the overlay.
PROBE_ADD = ("add", "cmd/wisp/zz156probe_windows_test.go",
             os.path.join(HERE, "zz156probe_windows_test.go"))

# --- AC#2's举证 / AC#5's teeth: every cell below COMPILES, none is a type error.
# m1: the ingredient this ticket chose, removed. exitStatus keeps its shape and
#     stops asking the OS, which is exactly the shipped口径 before ticket 156
#     (ask our own record). If the new cases do not redden here, the ingredient is
#     not load-bearing and this ticket has shipped a decoration.
OS_ARM_REMOVED = spec("cmd/wisp/slo_windows.go",
                      "\t// (b) the OS, for the window this type actually lives in.\n"
                      "\treturn s.queryProcess()\n",
                      "\t// MUTATION m1: the OS arm is gone; the answer comes from\n"
                      "\t// our own record, i.e. the bytes before ticket 156.\n"
                      "\treturn false, exitCodeUnknown\n")
# m2: the arm stays but can never decide "dead" (WAIT_OBJECT_0 is unreachable).
OS_ARM_NEVER_DEAD = spec("cmd/wisp/slo_windows.go",
                         "\t\tif waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited != windows.WAIT_OBJECT_0 {\n"
                         "\t\t\treturn // WAIT_TIMEOUT: alive. WAIT_FAILED: not judged here.\n"
                         "\t\t}\n",
                         "\t\tif waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited != 1 {\n"
                         "\t\t\treturn // MUTATION m2: nothing is ever judged dead.\n"
                         "\t\t}\n")
# m3: the arm stays but decides "dead" for everything except WAIT_FAILED, i.e. it
#     reads a LIVE subject as gone. Only case 16a can catch this one.
OS_ARM_ALWAYS_DEAD = spec("cmd/wisp/slo_windows.go",
                          "\t\tif waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited != windows.WAIT_OBJECT_0 {\n"
                          "\t\t\treturn // WAIT_TIMEOUT: alive. WAIT_FAILED: not judged here.\n"
                          "\t\t}\n",
                          "\t\tif waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited == windows.WAIT_FAILED {\n"
                          "\t\t\treturn // MUTATION m3: everything else is called dead.\n"
                          "\t\t}\n")
# m4: the OTHER ingredient - our own record after a reaping. Deleting it is the
#     edit that keeps every new verdict here true and still breaks production
#     (os/exec will not lend the handle back after Wait), so it must have a holder.
SHORTCIRCUIT_DELETED = spec("cmd/wisp/slo_windows.go",
                            "\tif ps := s.cmd.ProcessState; ps != nil {\n\t\treturn ps.Exited(), ps.ExitCode()\n\t}\n",
                            "\t// MUTATION m4: the reaped-child arm is gone.\n")
# m5: the liveness question stays, the SECOND syscall is dropped: gone=true with a
#     code nobody asked the OS for.
CODE_READ_DROPPED = spec("cmd/wisp/slo_windows.go",
                         "\t\tif cerr := windows.GetExitCodeProcess(windows.Handle(h), &osCode); cerr != nil {\n\t\t\treturn\n\t\t}\n",
                         "\t\tosCode = 1 // MUTATION m5: GetExitCodeProcess never asked.\n")


def case_off(rel, fname, tag):
    """Neutralise ONE case by renaming its function - code untouched. This is the
    escape hatch AC#5 asks to be measured rather than asserted."""
    return spec(rel,
                "func %s(t *testing.T) {" % fname,
                "func notATEST_156_%s_neutralised(t *testing.T) {" % tag)


NEW_CASES = {
    "case15": "TestSLO156ExitedAsksTheOSForAChildNobodyReaped",
    "case16": "TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees",
    "case17": "TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn",
    "case18": "TestSLO156WaitReadyNamesTheDeadSubjectToo",
}
NEW_TEST_REL = "cmd/wisp/slo_exit_os_156_windows_test.go"
OLD_TEST_REL = "cmd/wisp/slo_report_144_windows_test.go"
CASE13 = "TestSLO149ExitedGiveUpSentenceCarriesTheLastReading"


def apply_overlay(base, name, entries):
    out = os.path.join(os.environ.get("TMP", "/tmp"), "wisp156-ruler", name)
    os.makedirs(out, exist_ok=True)
    texts, adds = {}, {}
    replace = {}
    for ent in entries:
        if ent[0] == "add":
            _, rel, src = ent
            with open(src, "r", encoding="utf-8", newline="") as fh:
                body = fh.read()
            if "package main" not in body:
                sys.exit("FATAL %s: added file %s is not package main" % (name, rel))
            replace["%s/%s" % (base.replace(os.sep, "/"), rel)] = src.replace(os.sep, "/")
            continue
        _, rel, old, new = ent
        if rel not in texts:
            with open(os.path.join(base, rel.replace("/", os.sep)), "r", encoding="utf-8", newline="") as fh:
                texts[rel] = fh.read()
        hits = texts[rel].count(old)
        if hits != 1:
            sys.exit("FATAL %s: literal hits=%d in %s, want exactly 1" % (name, hits, rel))
        texts[rel] = texts[rel].replace(old, new)
        adds.setdefault(rel, []).append((old, new))
    for rel, body in texts.items():
        dst = os.path.join(out, os.path.basename(rel))
        with open(dst, "w", encoding="utf-8", newline="") as fh:
            fh.write(body)
        with open(dst, "r", encoding="utf-8", newline="") as fh:  # PROOF step
            back = fh.read()
        for old, new in adds[rel]:
            if new not in back or old in back:
                sys.exit("FATAL %s: overlay proof failed for %s - refusing to emit a reading" % (name, rel))
        replace["%s/%s" % (base.replace(os.sep, "/"), rel)] = dst.replace(os.sep, "/")
    overlay = os.path.join(out, "overlay.json")
    with open(overlay, "w", encoding="utf-8") as fh:
        json.dump({"Replace": replace}, fh)
    return overlay, sorted(replace)


def run(base, name, entries, log_dir, run_pat, sel_extra=()):
    if any("cover" in a for a in sel_extra):
        sys.exit("FATAL %s: -overlay must never meet -cover* (dispatch pitfall 2)" % name)
    overlay, replaced = (None, []) if not entries else apply_overlay(base, name, entries)
    cmd = [GO, "test", "-count=1", "-v", "-run", run_pat, "./cmd/wisp/"]
    if overlay:
        cmd[3:3] = ["-overlay", overlay]
    cmd += list(sel_extra)
    env = dict(os.environ)
    env["PATH"] = os.path.join(base, DLL_DIRNAME) + os.pathsep + env["PATH"]
    proc = subprocess.run(cmd, cwd=base, env=env, capture_output=True, text=True)
    out = proc.stdout + proc.stderr
    os.makedirs(log_dir, exist_ok=True)
    logpath = os.path.join(log_dir, name + ".log")
    with open(logpath, "w", encoding="utf-8", newline="") as fh:
        fh.write("$ %s\n$ tree=%s  overlaid=%s\n%s" % (shellish(cmd), base, replaced or "-", out))
    ran = len(re.findall(r"^=== RUN", out, re.M))
    passed = len(re.findall(r"^--- PASS", out, re.M))
    fails = re.findall(r"^--- FAIL: (\S+)", out, re.M)
    skips = len(re.findall(r"^--- SKIP", out, re.M))
    print("%-34s rc=%-3d RUN=%-4d PASS=%-4d FAIL=%-4d SKIP=%-3d %s%s" % (
        name, proc.returncode, ran, passed, len(fails), skips,
        "LOADED=YES" if ran else "LOADED=NO(dll/PATH) ",
        (" red=" + ",".join(fails)) if fails else ""))
    return {"name": name, "rc": proc.returncode, "ran": ran, "fails": fails, "log": logpath}


def shellish(cmd):
    """Quote every argv element so the logged command line is paste-safe (the
    -run pattern holds a bare '|', which a shell reads as a pipeline)."""
    return " ".join(shlex.quote(a) for a in cmd)


if __name__ == "__main__":
    if len(sys.argv) < 4:
        sys.exit(__doc__)
    mode, base, log_dir = sys.argv[1], sys.argv[2], sys.argv[3]
    if mode == "probe":
        r = run(base, "probe-156", [PROBE_ADD], log_dir, "^" + SELPROBE + "$")
        sys.exit(0 if r["ran"] and not r["fails"] else 1)
    if mode == "mut":
        args = sys.argv[4:]
        sel = SEL156
        cells = []
        for a in args:
            if a.startswith("--sel="):
                sel = a[len("--sel="):]
            else:
                cells.append(a)
        all_off = [case_off(NEW_TEST_REL, f, tag) for tag, f in sorted(NEW_CASES.items())]
        M = {
            "asis": [],
            "m1-os-arm-removed": [OS_ARM_REMOVED],
            "m2-os-arm-never-dead": [OS_ARM_NEVER_DEAD],
            "m3-os-arm-always-dead": [OS_ARM_ALWAYS_DEAD],
            "m4-shortcircuit-deleted": [SHORTCIRCUIT_DELETED],
            "m5-code-read-dropped": [CODE_READ_DROPPED],
            # the escape hatch, measured: does anything still notice once the
            # holders of this ticket's verdicts are removed one at a time?
            "m1-plus-all-new-cases-off": [OS_ARM_REMOVED] + all_off,
            "case15-off": [case_off(NEW_TEST_REL, NEW_CASES["case15"], "case15")],
            "case16-off": [case_off(NEW_TEST_REL, NEW_CASES["case16"], "case16")],
            "case17-off": [case_off(NEW_TEST_REL, NEW_CASES["case17"], "case17")],
            "case18-off": [case_off(NEW_TEST_REL, NEW_CASES["case18"], "case18")],
            "m4-plus-case13-off": [SHORTCIRCUIT_DELETED,
                                   case_off(OLD_TEST_REL, CASE13, "case13")],
            # --- AC#4(ii): re-measure the reading that sits in the 144 file's
            # --- case 14 comment, at THIS anchor, with THIS ticket's bytes under
            # --- it. Same two specs ticket 152's ruler used, re-typed here on
            # --- purpose: probes/152/my152.py is another cell's unfinished
            # --- deliverable and is not executed, not restored and not committed
            # --- by this程. Run with --sel='TestSLO144|TestSLO147|TestSLO149|TestSLO152'.
            "152-g1-restored-borrowed-value-case14off": [
                spec("cmd/wisp/slo_windows.go", "\treturn offsetUnknown\n}", "\treturn 0\n}"),
                case_off(OLD_TEST_REL,
                         "TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne", "case14")],
        }
        for cell in cells:
            if cell not in M:
                sys.exit("unknown cell %s (known: %s)" % (cell, ",".join(sorted(M))))
            run(base, cell, M[cell], log_dir, sel)
        sys.exit(0)
    sys.exit("unknown mode %s" % mode)
