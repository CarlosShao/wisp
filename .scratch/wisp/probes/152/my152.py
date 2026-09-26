#!/usr/bin/env python3
"""Ticket 152 mutation ruler (lives OUTSIDE the repository).

Lessons hard-coded from the two rulers this ticket inherits:
  * ticket 149 第 1 格: two specs against the SAME file must CHAIN, and the
    proof is that the overlay's own bytes are re-read before any run happens.
    This script FATALs (emits no reading) unless, for every spec, the new text
    is present and the old text is absent in the file handed to the compiler.
  * ticket 149 第 5 格 / dispatch pitfall 2: -overlay and -cover* together make
    Go silently IGNORE the overlay. This script never passes a cover flag; it
    also refuses any go-test argument containing "cover".
  * dispatch pitfall 1: "did it even run" is decided ONLY by the === RUN count
    (0xc0000135 also comes from a mangled PATH, not just a missing DLL).
  * ticket 149 第 7 格: red/green is counted from anchored '^--- FAIL:' lines
    only, never from a bare grep (the runner's own summary line is a false hit).
"""
import json
import os
import re
import shlex
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
REPO = "D:/work/workspace/projects plans/Wisp"
PROBES = os.path.join(REPO, ".scratch", "wisp", "probes", "152")
GO = "go"

SEL144 = "TestSLO144|TestSLO147|TestSLO149"  # the ticket 149 family, verbatim


def spec(rel, old, new):
    return (rel, old, new)


# ---------------------------------------------------------------- mutations ---
# AC#3's leg, on the ANCHOR bytes: the fallback arm of contradictionOffset.
FALLBACK_TO_ZERO = spec("cmd/wisp/slo_windows.go",
                        "\treturn inputOffset\n}", "\treturn 0\n}")
FALLBACK_TO_MINUS1 = spec("cmd/wisp/slo_windows.go",
                          "\treturn inputOffset\n}", "\treturn -1\n}")
# AC#2's leg: the exited branch's last-reading half (ticket 149's p7, same text).
EXITED_DROP_SUMMARY = spec("cmd/wisp/slo_windows.go",
                           's.pid, s.exitCode(), last.summary())',
                           's.pid, s.exitCode(), "nothing read")')
CASE13_NOT_A_TEST = spec("cmd/wisp/slo_report_144_windows_test.go",
                         "func TestSLO149ExitedGiveUpSentenceCarriesTheLastReading(t *testing.T) {",
                         "func notATEST_152_case13_neutralised(t *testing.T) {")
# AC#3's leg on the FIXED bytes: put the known-lying value back, keeping every
# other part of the fix. Written by the caller once the fixed text exists.


# AC#3's leg, on the FIXED bytes: put the borrowed position back, and make each
# of the two rendering arms lie. These are the mutations that answer "would this
# ring on the unfixed code" for the shape ticket 152 shipped.
RESTORE_BORROWED_VALUE = spec("cmd/wisp/slo_windows.go",
                              "\treturn offsetUnknown\n}", "\treturn 0\n}")
RENDERER_ALWAYS_PRINTS_POSITION = spec("cmd/wisp/slo_windows.go",
                                       "\tif offset < 0 {\n\t\treturn fmt.Errorf(\"%d bytes contradict a subject report at an offset the decoder did not name: %w\", size, err)",
                                       "\tif false {\n\t\treturn fmt.Errorf(\"%d bytes contradict a subject report at an offset the decoder did not name: %w\", size, err)")
RENDERER_NEVER_PRINTS_POSITION = spec("cmd/wisp/slo_windows.go",
                                      "\tif offset < 0 {\n\t\treturn fmt.Errorf(\"%d bytes contradict a subject report at an offset the decoder did not name: %w\", size, err)",
                                      "\tif true {\n\t\treturn fmt.Errorf(\"%d bytes contradict a subject report at an offset the decoder did not name: %w\", size, err)")
CASE14_NOT_A_TEST = spec("cmd/wisp/slo_report_144_windows_test.go",
                         "func TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne(t *testing.T) {",
                         "func notATEST_152_case14_neutralised(t *testing.T) {")


def apply_muts(base, name, specs):
    out = os.path.join(HERE, "mut", name)
    os.makedirs(out, exist_ok=True)
    texts = {}
    for rel, old, new in specs:
        if rel not in texts:
            src = os.path.join(base, rel.replace("/", os.sep))
            with open(src, "r", encoding="utf-8", newline="") as fh:
                texts[rel] = fh.read()
        hits = texts[rel].count(old)
        if hits != 1:
            sys.exit("FATAL %s: literal hits=%d in %s, want exactly 1" % (name, hits, rel))
        texts[rel] = texts[rel].replace(old, new)
    replace = {}
    for rel in texts:
        dst = os.path.join(out, os.path.basename(rel))
        with open(dst, "w", encoding="utf-8", newline="") as fh:
            fh.write(texts[rel])
        # PROOF, against the very bytes that go into the compiler:
        with open(dst, "r", encoding="utf-8", newline="") as fh:
            back = fh.read()
        for (r, old, new) in specs:
            if r != rel:
                continue
            if new not in back or old in back:
                sys.exit("FATAL %s: overlay proof failed for %s "
                         "(new in=%s / old still in=%s) - refusing to emit a reading"
                         % (name, rel, new in back, old in back))
        replace["%s/%s" % (base.replace(os.sep, "/"), rel)] = dst.replace(os.sep, "/")
    overlay = os.path.join(out, "overlay.json")
    with open(overlay, "w", encoding="utf-8") as fh:
        json.dump({"Replace": replace}, fh)
    return overlay, sorted(replace)


def shellish(cmd):
    """Quote every argv element so the logged line is paste-safe.

    A logged 'command original' that cannot be pasted is not a reading: the
    -run pattern contains a bare '|', which a shell would read as a pipeline.
    """
    return " ".join(shlex.quote(a) for a in cmd)


def run(base, name, specs, log_dir, run_pat=SEL144, extra=()):
    if any("cover" in a for a in extra):
        sys.exit("FATAL %s: -overlay must never be combined with -cover* (pitfall 2)" % name)
    if specs:
        overlay, replaced = apply_muts(base, name, specs)
    else:
        overlay, replaced = None, []
    cmd = [GO, "test", "-count=1", "-v", "-run", run_pat, "./cmd/wisp/"]
    if overlay:
        cmd[3:3] = ["-overlay", overlay]
    cmd += list(extra)
    env = dict(os.environ)
    dll = os.path.join(base, "third_party", "sherpa-onnx")
    env["PATH"] = dll + os.pathsep + env["PATH"]
    proc = subprocess.run(cmd, cwd=base, env=env, capture_output=True, text=True)
    out = proc.stdout + proc.stderr
    os.makedirs(log_dir, exist_ok=True)
    logpath = os.path.join(log_dir, name + ".log")
    with open(logpath, "w", encoding="utf-8", newline="") as fh:
        fh.write("$ %s\n$ tree=%s  mutated_files=%s\n%s" % (shellish(cmd), base, replaced or "-", out))
    ran = len(re.findall(r"^=== RUN", out, re.M))
    passed = len(re.findall(r"^--- PASS", out, re.M))
    fails = re.findall(r"^--- FAIL: (\S+)", out, re.M)
    skips = len(re.findall(r"^--- SKIP", out, re.M))
    red_lines = [ln for ln in out.splitlines() if ln.lstrip().startswith("Error Trace") or
                 ln.lstrip().startswith("Error:") or re.match(r"^\s+\S+_test\.go:\d+:", ln)]
    print("%-30s rc=%d RUN=%d PASS=%d FAIL=%d SKIP=%d reds=%s" %
          (name, proc.returncode, ran, passed, len(fails), skips, ",".join(fails) or "-"))
    return {"name": name, "rc": proc.returncode, "ran": ran, "passed": passed,
            "fails": fails, "skips": skips, "log": logpath, "red": red_lines}


if __name__ == "__main__":
    if len(sys.argv) < 3:
        sys.exit("usage: my152.py <base-tree> <log-dir> [cell...] [--sel=PATTERN]")
    base = sys.argv[1]
    log_dir = sys.argv[2]
    sel = SEL144
    cells = []
    for a in sys.argv[3:]:
        if a.startswith("--sel="):
            sel = a[len("--sel="):]
        else:
            cells.append(a)
    M = {
        "asis": [],
        # --- anchor bytes (ticket 149's legs, re-measured here) ---
        "f1-fallback-to-zero": [FALLBACK_TO_ZERO],
        "f2-fallback-to-minus1": [FALLBACK_TO_MINUS1],
        "e1-exited-drop-summary": [EXITED_DROP_SUMMARY],
        "e2-exited-drop-summary-case13off": [EXITED_DROP_SUMMARY, CASE13_NOT_A_TEST],
        "e3-case13-not-a-test-only": [CASE13_NOT_A_TEST],
        # --- fixed bytes (ticket 152's own arms) ---
        "g1-restored-borrowed-value": [RESTORE_BORROWED_VALUE],
        "g1-restored-borrowed-value-case14off": [RESTORE_BORROWED_VALUE, CASE14_NOT_A_TEST],
        "g2-case14-not-a-test-only": [CASE14_NOT_A_TEST],
        "g3-renderer-always-prints-position": [RENDERER_ALWAYS_PRINTS_POSITION],
        "g4-renderer-never-prints-position": [RENDERER_NEVER_PRINTS_POSITION],
    }
    for cell in cells:
        if cell not in M:
            sys.exit("unknown cell %s (known: %s)" % (cell, ",".join(sorted(M))))
        run(base, cell, M[cell], log_dir, run_pat=sel)
