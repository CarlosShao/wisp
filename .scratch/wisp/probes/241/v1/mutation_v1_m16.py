#!/usr/bin/env python3
"""M16: is TestLevelFrameIsOneSeamFrame an alive nail or a decoration?

None of the 15 level.go mutants could reach it, because it asserts constants that
live in audio.go. That is not evidence of decoration, so mutate the carrier itself:
FrameSamples 512 -> 511 (which also moves FrameBytes and FrameDuration). Injected
with -overlay only; no tracked file is written.
"""
import json
import os
import subprocess

REPO = "D:/work/workspace/projects plans/Wisp"
SRC = REPO + "/internal/audio/audio.go"
PROBE = REPO + "/.scratch/wisp/probes/241/v1"

with open(os.path.join(REPO, "internal", "audio", "audio.go"), "rb") as fh:
    text = fh.read().decode("utf-8")

anchor = "const FrameSamples = 512"
assert text.count(anchor) == 1, "anchor not unique"
mutated = text.replace(anchor, "const FrameSamples = 511")
mpath = PROBE + "/mut-M16-audio-framesamples-511.go"
with open(mpath.replace("/", os.sep), "wb") as fh:
    fh.write(mutated.encode("utf-8"))

opath = PROBE + "/ov-M16.json"
with open(opath.replace("/", os.sep), "wb") as fh:
    fh.write(json.dumps({"Replace": {SRC: mpath}}).encode("utf-8"))

env = dict(os.environ)
env["PATH"] = "%s/third_party/sherpa-onnx:%s/build:%s" % (REPO, REPO, env["PATH"])
proc = subprocess.run(["go", "test", "-overlay=" + opath, "-count=1", "-v",
                       "-run", "TestLevel", "./internal/audio/"],
                      cwd=REPO, env=env, capture_output=True, text=True)
out = proc.stdout + proc.stderr
red = sorted({ln.split()[2].rstrip(":") for ln in out.splitlines()
               if ln.startswith("--- FAIL: ")})
print("M16-audio-framesamples-511 rc=%d red=%s" % (proc.returncode, red))
target = "TestLevelFrameIsOneSeamFrame"
print("TARGET-ALIVE: %s -> %s" % (target, "YES (it goes red)" if target in red else "NO (survived)"))
with open(PROBE + "/M16.full.txt", "wb") as fh:
    fh.write(out.encode("utf-8"))
