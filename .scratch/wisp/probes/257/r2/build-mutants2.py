import json
import pathlib
import sys

ROOT = pathlib.Path("D:/work/workspace/projects plans/Wisp")
SRC = ROOT / "cmd/wisp/firstrun.go"
TSRC = ROOT / "cmd/wisp/firstrun_257_nonpreset_test.go"
OUT = ROOT / ".scratch/wisp/probes/257/r2/mut"
src = SRC.read_text(encoding="utf-8")
tsrc = TSRC.read_text(encoding="utf-8")


def sub(content, needle, repl, name):
    n = content.count(needle)
    if n != 1:
        sys.exit("MUTATION SETUP FAIL %s: needle count %d, want 1: %r" % (name, n, needle))
    return content.replace(needle, repl, 1)


jobs = []

# M7: delete the non-preset condition from the receipt (the clause my new file
# pins). Expect TestTicket257R2AC1ReceiptStatesTheNonPresetCondition to go red.
jobs.append(("M7", SRC, sub(
    src,
    "；非预设名必须自己写 protocol，否则这份文件加载不过",
    "。",
    "M7"), "firstrun.go"))

# M8: prove the behavioural half is not vacuous - point shape 1 at a PRESET name
# and the "a non-preset row with no protocol must be refused" step must scream.
m8_lines = tsrc.split("\n")
hit = 0
for i, ln in enumerate(m8_lines):
    if ln.startswith("\tnoProtocol := base + "):
        m8_lines[i] = ln.replace("t257npGhost", "t257Provider")
        hit += 1
if hit != 1:
    sys.exit("MUTATION SETUP FAIL M8: shape-1 line hits %d, want 1" % hit)
jobs.append(("M8", TSRC, "\n".join(m8_lines), "firstrun_257_nonpreset_test.go"))

for name, original, content, fname in jobs:
    if content == original.read_text(encoding="utf-8"):
        sys.exit("MUTATION SETUP FAIL %s: mutant identical to source" % name)
    d = OUT / name
    d.mkdir(parents=True, exist_ok=True)
    mp = d / fname
    mp.write_text(content, encoding="utf-8")
    (d / "overlay.json").write_text(json.dumps({"Replace": {str(original): str(mp)}}), encoding="utf-8")
    print("built %s -> %s" % (name, mp))
