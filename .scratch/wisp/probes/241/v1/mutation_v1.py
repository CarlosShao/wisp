#!/usr/bin/env python3
"""241-v1 mutation harness (rev 2).

Adjudicator leg for ticket 241. Mutates internal/audio/level.go WITHOUT ever
writing to it: every mutant is a file under .scratch, injected through
`go test -overlay`. The tracked source is never modified, so there is nothing
to restore and no chance of zeroing a tracked file.

rev 2 fixes two harness bugs found on the first run:
  (a) M1's original anchor deleted the only use of loop var `v` -> the mutant
      did not compile, which is a harness bug and must never be reported as a
      red test. Now `v * 0` keeps `v` live while pinning the numerator to zero.
  (b) os.path.join emitted backslashes into the overlay JSON value. Values are
      normalised to forward slashes now.
  (c) A build failure is now detected and aborts the harness loudly.
  (d) C0 is a semantically neutral control mutant (comment text only) that MUST
      stay green; if it goes red the harness is not measuring what it claims.

Usage: python mutation_v1.py [MUTANT_NAME ...]     (no args = run all)
"""
import json
import os
import subprocess
import sys

REPO = "D:/work/workspace/projects plans/Wisp"
SRC = "internal/audio/level.go"
ABS_SRC = REPO + "/" + SRC
PROBE = REPO + "/.scratch/wisp/probes/241/v1"

# (name, old, new)
MUTANTS = [
    ("C0-comment-only-control",
     "// THE SCALE. level = sqrt(mean(sample^2)) / LevelFullScale over one C8 seam",
     "// THE SCALE. level = sqrt(mean(sample^2)) / LevelFullScale over one C9 seam"),
    ("M1-numerator-zero",
     "sumSquares += v * v", "sumSquares += v * 0"),
    ("M2-return-constant-half",
     "return rms / LevelFullScale", "return math.Sqrt(0.25) + rms*0"),
    ("M3-tolerance-wide-1e9",
     "const SineLevelTolerance = 5e-5", "const SineLevelTolerance = 1e9"),
    ("M4-tolerance-narrow-1e-9",
     "const SineLevelTolerance = 5e-5", "const SineLevelTolerance = 1e-9"),
    ("M5-tolerance-at-bound-1e-3",
     "const SineLevelTolerance = 5e-5", "const SineLevelTolerance = 1e-3"),
    ("M6-denominator-32767",
     "const LevelFullScale = 32768.0", "const LevelFullScale = 32767.0"),
    ("M7-negative-scale",
     "return rms / LevelFullScale", "return -rms / LevelFullScale"),
    ("M8-gain-two-x",
     "return rms / LevelFullScale", "return 2 * (rms / LevelFullScale)"),
    ("M9-compress-sqrt",
     "return rms / LevelFullScale", "return math.Sqrt(rms / LevelFullScale)"),
    ("M10-framelength-too-lax",
     "if len(frame) != FrameBytes {", "if len(frame) < FrameBytes {"),
    ("M11-mean-divisor-off-by-one",
     "float64(sumSquares) / float64(len(samples))",
     "float64(sumSquares) / float64(len(samples)-1)"),
    ("M12-decode-byte-order-swapped",
     "samples[i] = int16(uint16(frame[2*i]) | uint16(frame[2*i+1])<<8)",
     "samples[i] = int16(uint16(frame[2*i+1]) | uint16(frame[2*i])<<8)"),
    # M13 is the only mutant that gives TestLevelSameInputSameBits any teeth: it
    # keeps the function pure-looking but makes the result depend on call
    # history. Multi-anchor because a package-level accumulator is needed.
    ("M13-hidden-state-accumulator",
     [("// LevelFullScale is the denominator",
       "var v1Accum int64\n\n// LevelFullScale is the denominator"),
      ("\tvar sumSquares int64\n\tfor _, s := range samples {\n\t\tv := int64(s)\n\t\tsumSquares += v * v\n\t}",
       "\tfor _, s := range samples {\n\t\tv := int64(s)\n\t\tv1Accum += v * v\n\t}\n\tsumSquares := v1Accum")]),
    # M14/M15 answer the boundary question M5 opened: the tolerance nail admits
    # values up to EXACTLY 1e-3, so does the scale still have teeth at that
    # admitted ceiling? Both mutants swap the RMS detector for a mean-of-abs
    # detector (the classic plausible-alternative that shares every square-wave
    # endpoint with RMS but differs by 2/sqrt(pi) vs 1/sqrt(2) on a sine).
    ("M14-mean-abs-detector",
     [("\tvar sumSquares int64\n\tfor _, s := range samples {\n\t\tv := int64(s)\n\t\tsumSquares += v * v\n\t}\n\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\treturn rms / LevelFullScale",
       "\tvar sumAbs int64\n\tfor _, s := range samples {\n\t\tsumAbs += int64(math.Abs(float64(s)))\n\t}\n\treturn (float64(sumAbs) / float64(len(samples))) / LevelFullScale")]),
    ("M15-mean-abs-at-admitted-ceiling",
     [("\tvar sumSquares int64\n\tfor _, s := range samples {\n\t\tv := int64(s)\n\t\tsumSquares += v * v\n\t}\n\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\treturn rms / LevelFullScale",
       "\tvar sumAbs int64\n\tfor _, s := range samples {\n\t\tsumAbs += int64(math.Abs(float64(s)))\n\t}\n\treturn (float64(sumAbs) / float64(len(samples))) / LevelFullScale"),
      ("const SineLevelTolerance = 5e-5", "const SineLevelTolerance = 1e-3")]),
]


def read_source():
    with open(os.path.join(REPO, SRC.replace("/", os.sep)), "rb") as fh:
        return fh.read()


def build(name, old, new):
    text = read_source().decode("utf-8")
    pairs = old if isinstance(old, list) else [(old, new)]
    for o, n in pairs:
        hits = text.count(o)
        if hits != 1:
            raise SystemExit("anchor %r found %d times in %s - harness bug, not a result"
                             % (o, hits, SRC))
    for o, n in pairs:
        text = text.replace(o, n)
    mpath = PROBE + "/mut-" + name + ".go"
    with open(mpath.replace("/", os.sep), "wb") as fh:
        fh.write(text.encode("utf-8"))
    return mpath


def write_overlay(mapping, tag):
    opath = PROBE + "/ov-" + tag + ".json"
    with open(opath.replace("/", os.sep), "wb") as fh:
        fh.write(json.dumps({"Replace": mapping}).encode("utf-8"))
    return opath


def run(overlay):
    env = dict(os.environ)
    env["PATH"] = "%s/third_party/sherpa-onnx:%s/build:%s" % (REPO, REPO, env["PATH"])
    cmd = ["go", "test", "-overlay=" + overlay, "-count=1", "-v", "./internal/audio/"]
    proc = subprocess.run(cmd, cwd=REPO, env=env, capture_output=True, text=True)
    out = proc.stdout + proc.stderr
    if "build failed" in out:
        raise SystemExit("MUTANT DID NOT COMPILE (harness bug, not a verdict):\n" + out[:2000])
    # rev 3: Go prints "--- FAIL: Name (0.00s)" -- the colon is glued to FAIL, so
    # the field is index 2 and the prefix must include the colon. rev 2 matched
    # "--- FAIL " and therefore reported every mutant as red_count=0, which is
    # the shape of a blind ruler, not of a green package.
    top = sorted({ln.split()[2].rstrip(":") for ln in out.splitlines()
                  if ln.startswith("--- FAIL: ")})
    skips = sorted({ln.split()[2].rstrip(":") for ln in out.splitlines()
                    if ln.startswith("--- SKIP: ")})
    return proc.returncode, top, skips, out


def main():
    wanted = set(sys.argv[1:])
    log = []
    rc, red, skips, _ = run(write_overlay({}, "baseline"))
    line = "BASELINE(no mutation) rc=%d red=%s skip=%s" % (rc, red, skips)
    print(line, flush=True)
    log.append(line)
    if rc != 0:
        print("baseline not green - stop, the harness is not measuring what it claims")
        return 1
    if "TestLiveWasapiSmoke" not in skips:
        print("baseline does not reproduce the known SKIP - the -v parser is blind")
        return 1
    for entry in MUTANTS:
        name, old = entry[0], entry[1]
        new = entry[2] if len(entry) > 2 else None
        if wanted and name not in wanted:
            continue
        mpath = build(name, old, new)
        rc, red, skips, out = run(write_overlay({ABS_SRC: mpath}, name))
        line = "%-32s rc=%d red_count=%d red=%s" % (name, rc, len(red), red)
        print(line, flush=True)
        log.append(line)
        if rc == 0:
            print("   ^^ mutant SURVIVED (whole package stayed green)", flush=True)
            log.append("   SURVIVED")
        with open(PROBE + "/" + name + ".full.txt", "wb") as fh:
            fh.write(out.encode("utf-8"))
    # Self-check on the instrument itself: the comment-only control must survive
    # and at least one real mutant must redden. If both fail the harness is blind.
    # Only meaningful on a full run; on a filtered run the control is not in the
    # log, so say so instead of crying BLIND (that false alarm happened at
    # 14:21 on the single-mutant M13 run and is recorded in the log).
    ctl_line = next((ln for ln in log if ln.startswith("C0-comment-only-control")), "")
    real = [ln for ln in log if "red_count=" in ln and not ln.startswith("C0")]
    n_red = sum(1 for ln in real if "red_count=0 " not in ln)
    if not ctl_line:
        line = "INSTRUMENT-SELFCHECK NOT-APPLICABLE (filtered run: control C0 was not executed in this invocation)"
    else:
        line = "INSTRUMENT-SELFCHECK control_survived=%s real_mutants_with_red=%d/%d => %s" % (
            "red_count=0 " in ctl_line, n_red, len(real),
            "OK" if ("red_count=0 " in ctl_line and n_red > 0) else "BLIND-HARNESS-ALERT")
    print(line, flush=True)
    log.append(line)
    with open(PROBE + "/mutation-v1-readings.txt", "ab") as fh:
        for ln in log:
            fh.write((ln + "\n").encode("utf-8"))
    return 0


if __name__ == "__main__":
    sys.exit(main())

# NOTE (241-v1): this harness had TWO bugs of its own whose readings would have
# been wrong in a flattering direction. Both are recorded here, not quietly fixed:
#   rev1 -> rev2: M1's anchor deleted the only use of the loop variable, so the
#     mutant failed to COMPILE; the harness reported "red_count=0", which reads
#     exactly like "the tests did not catch it". It was never a verdict.
#   rev2 -> rev3: the -v parser matched "--- FAIL " while Go prints
#     "--- FAIL: " (colon glued to the word), so EVERY mutant looked green. The
#     blind parser was only caught because BASELINE reported skip=[] while the
#     package is known to skip TestLiveWasapiSmoke.
# Protections added after that: baseline must be rc=0 AND must reproduce that
# known SKIP; a "build failed" outcome aborts loudly instead of being counted as
# green; C0 is a comment-only control that must SURVIVE (a harness that reddens
# on a comment is not measuring the scale); the trailing INSTRUMENT-SELFCHECK
# line fails loudly if the control dies or if no real mutant reddens.
