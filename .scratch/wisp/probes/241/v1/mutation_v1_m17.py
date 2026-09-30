#!/usr/bin/env python3
"""M17: cross-check 241-r1's own claimed mutant X3 (LevelFullScale 32768 -> 16384).

r1 evidence sec 5 claims X3 reddens exactly 4 named tests. Overlay only; no tracked
file is written. This is a check of the IMPLEMENT LEG'S ROSTER CLAIM, not of the scale.
"""
import json
import os
import subprocess

REPO = "D:/work/workspace/projects plans/Wisp"
SRC = REPO + "/internal/audio/level.go"
PROBE = REPO + "/.scratch/wisp/probes/241/v1"

with open(os.path.join(REPO, "internal", "audio", "level.go"), "rb") as fh:
    text = fh.read().decode("utf-8")

anchor = "const LevelFullScale = 32768.0"
assert text.count(anchor) == 1, "anchor not unique"
mpath = PROBE + "/mut-M17-r1-claimed-X3.go"
with open(mpath.replace("/", os.sep), "wb") as fh:
    fh.write(text.replace(anchor, "const LevelFullScale = 16384.0").encode("utf-8"))

opath = PROBE + "/ov-M17.json"
with open(opath.replace("/", os.sep), "wb") as fh:
    fh.write(json.dumps({"Replace": {SRC: mpath}}).encode("utf-8"))

env = dict(os.environ)
env["PATH"] = "%s/third_party/sherpa-onnx:%s/build:%s" % (REPO, REPO, env["PATH"])
proc = subprocess.run(["go", "test", "-overlay=" + opath, "-count=1", "-v", "./internal/audio/"],
                      cwd=REPO, env=env, capture_output=True, text=True)
out = proc.stdout + proc.stderr
red = sorted({ln.split()[2].rstrip(":") for ln in out.splitlines() if ln.startswith("--- FAIL: ")})
print("M17-r1-claimed-X3 rc=%d red_count=%d red=%s" % (proc.returncode, len(red), red))
claimed = ["TestLevelFullScaleSquareEndpoints", "TestLevelSineAgainstRootTwo",
           "TestLevelSineToleranceIsNamedAndBounded", "TestLevelOverBoundedChannelFromWavInjector"]
print("r1 CLAIMED 4: %s" % claimed)
print("MATCHES CLAIM: %s" % (red == sorted(claimed)))
with open(PROBE + "/M17.full.txt", "wb") as fh:
    fh.write(out.encode("utf-8"))
