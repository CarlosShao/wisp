#!/usr/bin/env python3
# 33-r10 mutation harness (ticket 33 AC#13 item 3).
# Never edits tracked files: mutates a copy outside the repo and feeds it to the
# compiler with -overlay. Runs the built test exe with CWD = cmd/wisp and the
# sherpa DLLs on PATH (the harness the dispatch names).

import json
import os
import subprocess
import sys

REPO = os.path.abspath(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(
    os.path.abspath(__file__))))))
REPO = "D:/work/workspace/projects plans/Wisp"
OUT = "D:/tmp/wisp33r10"
EVID = os.path.join(REPO, ".scratch", "wisp", "probes", "33", "r10", "mut")
TARGET = os.path.join(REPO, "cmd", "wisp", "panel_host_windows.go")
RUN_DIR = os.path.join(REPO, "cmd", "wisp")
DLL = os.path.join(REPO, "third_party", "sherpa-onnx")
PATTERN = "AC13ColdStartPageOver"

os.makedirs(OUT, exist_ok=True)
os.makedirs(EVID, exist_ok=True)

with open(TARGET, "r", encoding="utf-8") as fh:
    src = fh.read()

SWAP_BODY = """func (m *PanelManager) coldStartPageHandover(ctx context.Context, t0 time.Time) float64 {
	rtMs := m.firstRoundTripLocked(ctx, t0)

	if err := m.serveEntry(); err != nil {
		m.serveNotBuiltNoticeLocked()
	}
	return rtMs
}"""

SWAPPED = """func (m *PanelManager) coldStartPageHandover(ctx context.Context, t0 time.Time) float64 {
	if err := m.serveEntry(); err != nil {
		m.serveNotBuiltNoticeLocked()
	}
	rtMs := m.firstRoundTripLocked(ctx, t0)
	return rtMs
}"""

PROBE_SRC = """	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"></head><body>` +
		`<script>if(window.wispProbeRT)window.wispProbeRT();</script></body></html>`)"""

PROBE_WIDENED = """	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"></head><body><p>loading the panel</p>` +
		`<script>if(window.wispProbeRT)window.wispProbeRT();</script></body></html>`)"""

for needle, label in ((SWAP_BODY, "handover body"), (PROBE_SRC, "probe document")):
    if src.count(needle) != 1:
        sys.stdout.write("LANDING CHECK FAILED: %s occurs %d time(s)\n" % (label, src.count(needle)))
        sys.exit(2)

variants = {}

# pristine: byte-identical replacement, must not change the colour.
variants["pristine"] = src

# swap: AC#13's named mutation, must go red.
variants["swap"] = src.replace(SWAP_BODY, SWAPPED)

# shift: same statements, line numbers moved (content anchor, not line anchor).
shift_block = "\n".join([
    "// shift-probe 33-r10: these lines exist to move every following line number",
    "// by eight. The ruler below must not care, because it reads the last document",
    "// the control was handed, not a position in this file.",
    "", "", "",
]) 
variants["shift"] = src.replace(SWAP_BODY, shift_block + SWAP_BODY)

# widen_swap: swap the order AND give the probe page lookable markup, so the
# "body is not script-only" tooth is defeated. The identity tooth must still bite.
variants["widen_swap"] = src.replace(SWAP_BODY, SWAPPED).replace(PROBE_SRC, PROBE_WIDENED)

env = dict(os.environ)
env["PATH"] = DLL + os.pathsep + env["PATH"]

report = []
for name, text in variants.items():
    mut_path = os.path.join(OUT, "panel_host_windows_%s.go" % name).replace("\\", "/")
    with open(mut_path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(text)
    overlay = {"Replace": {TARGET.replace("\\", "/"): mut_path}}
    ov_path = os.path.join(OUT, "overlay_%s.json" % name).replace("\\", "/")
    with open(ov_path, "w", encoding="utf-8", newline="\n") as fh:
        json.dump(overlay, fh)
    exe = os.path.join(OUT, "wisp33r10_%s.exe" % name).replace("\\", "/")

    # landing proof: the overlay must pair target -> replacement, and the
    # replacement must differ from the target except in the pristine run.
    diff = 0 if text == src else 1
    build = subprocess.run(
        ["go", "test", "-c", "-overlay", ov_path, "-o", exe, "./cmd/wisp"],
        cwd=REPO, env=env, capture_output=True, text=True)
    run = subprocess.run(
        [exe, "-test.run", PATTERN, "-test.v", "-test.count", "1"],
        cwd=RUN_DIR, env=env, capture_output=True, text=True)
    log = []
    log.append("variant=%s" % name)
    log.append("overlay_json=%s" % json.dumps(overlay, ensure_ascii=False))
    log.append("replacement_differs_from_target=%d" % diff)
    if os.path.exists(exe):
        log.append("exe_bytes=%d" % os.path.getsize(exe))
    else:
        log.append("exe_bytes=MISSING")
    log.append("buildrc=%d" % build.returncode)
    if build.stdout.strip() or build.stderr.strip():
        log.append("build_output=%s" % (build.stdout + build.stderr).strip())
    log.append("runrc=%d" % run.returncode)
    log.append("test_output_begin")
    log.append((run.stdout + run.stderr).strip())
    log.append("test_output_end")
    verdict = "RED" if run.returncode != 0 else "GREEN"
    log.append("colour=%s" % verdict)
    log.append("rc=%d" % run.returncode)
    body = "\n".join(log) + "\n"
    with open(os.path.join(EVID, "mut-%s.log" % name), "w", encoding="utf-8", newline="\n") as fh:
        fh.write(body)
    report.append("%s buildrc=%d runrc=%d colour=%s diff=%d" % (
        name, build.returncode, run.returncode, verdict, diff))

summary = "\n".join(report) + "\nrc=0\n"
with open(os.path.join(EVID, "mut-summary.log"), "w", encoding="utf-8", newline="\n") as fh:
    fh.write(summary)
sys.stdout.write(summary)
