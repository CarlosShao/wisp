#!/usr/bin/env python3
"""241-r1 AC#3 directed-mutation ruler (ticket 241, producer half).

Each mutant is applied to internal/audio/level.go, the level test name册 of
internal/audio is run, and the file is restored from the pristine copy this
script writes first. Nothing is deleted; the pristine copy and the readings
stay under .scratch/wisp/probes/241/r1/.

Run from anywhere:  python .scratch/wisp/probes/241/r1/mutation241.py
"""

import pathlib
import subprocess
import sys

REPO = pathlib.Path(r"D:\work\workspace\projects plans\Wisp")
LEVEL = REPO / "internal" / "audio" / "level.go"
OUT = REPO / ".scratch" / "wisp" / "probes" / "241" / "r1"
PRISTINE = OUT / "level.go.pristine"

MUTANTS = [
    (
        "X1-constant-numerator",
        "\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\treturn rms / LevelFullScale\n",
        "\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\t_ = rms\n\treturn math.Sqrt(0.25) // MUTANT X1: RMS numerator frozen to a constant\n",
    ),
    (
        "X2-tolerance-wide-enough-to-pass-anything",
        "const SineLevelTolerance = 5e-5\n",
        "const SineLevelTolerance = 1e9 // MUTANT X2: tolerance relaxed past every shape\n",
    ),
    (
        "X3-wrong-denominator",
        "const LevelFullScale = 32768.0\n",
        "const LevelFullScale = 16384.0 // MUTANT X3: the scale references the wrong full scale\n",
    ),
    (
        "X4-hidden-compression-curve",
        "\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\treturn rms / LevelFullScale\n",
        "\trms := math.Sqrt(float64(sumSquares) / float64(len(samples)))\n\treturn math.Sqrt(rms / LevelFullScale) // MUTANT X4: a gain curve folded into the producer\n",
    ),
]


def run_level_tests():
    proc = subprocess.run(
        ["go", "test", "./internal/audio/", "-count=1", "-v",
         "-run", "^(TestLevel|TestFrameLevel|TestSeamCodec)"],
        cwd=str(REPO), capture_output=True, text=True,
    )
    red = [ln.split(" ", 2)[2].split(" ")[0]
           for ln in proc.stdout.splitlines() if ln.startswith("--- FAIL:")]
    return proc.returncode, red, proc.stdout, proc.stderr


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    original = LEVEL.read_bytes()
    PRISTINE.write_bytes(original)
    log = []
    try:
        text = original.decode("ascii")
        for name, old, new in MUTANTS:
            if old not in text:
                raise SystemExit("anchor missing for %s - the ruler is stale" % name)
            LEVEL.write_bytes(text.replace(old, new, 1).encode("ascii"))
            rc, red, out, err = run_level_tests()
            line = "%-42s rc=%d red=%d %s" % (name, rc, len(red), ", ".join(sorted(red)) or "(NONE: the nail is decoration)")
            print(line)
            log.append(line)
            if err.strip():
                log.append("   stderr: " + err.strip().splitlines()[0])
    finally:
        LEVEL.write_bytes(original)
        rc, red, out, err = run_level_tests()
        tail = "restored pristine: rc=%d red=%d %s" % (rc, len(red), ", ".join(sorted(red)) or "(none)")
        print(tail)
        log.append(tail)
        (OUT / "mutation-readings.txt").write_text("\n".join(log) + "\n", encoding="ascii")


if __name__ == "__main__":
    main()
