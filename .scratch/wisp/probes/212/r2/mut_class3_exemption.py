"""212-r2 mutation self-proof for ticket 212 裁 (a).

Question answered: does the new "abbreviated path token is class 3, never
convicted" exclusion in tools/d22scan/main.go actually have teeth?

Method: patch the ONE line that feeds the exclusion into ban #9's loop so the
exemption map is always empty (the ban behaves as it shipped on 2026-10-03),
run `go run . -self-test`, then restore the pristine bytes in a finally block.
Nothing here touches git; the working tree is restored from the bytes read at
the start, and the md5 is re-measured after the run.
"""

import hashlib
import pathlib
import subprocess
import sys

ROOT = pathlib.Path(r"D:\work\workspace\projects plans\Wisp")
TARGET = ROOT / "tools" / "d22scan" / "main.go"
LOGDIR = ROOT / ".scratch" / "wisp" / "probes" / "212" / "r2" / "logs"
ORIG = "\t\t\tshort := shorthandPathStarts(c.Text)\n"
MUT = "\t\t\tshort := map[int]bool{} // MUTATION 212-r2: class-3 exemption removed\n"


def run_selftest(tag: str) -> tuple[int, str]:
    p = subprocess.run(
        ["go", "run", ".", "-self-test"],
        cwd=str(ROOT / "tools" / "d22scan"),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    out = "rc=%d\n--- stdout ---\n%s\n--- stderr ---\n%s\n" % (p.returncode, p.stdout, p.stderr)
    (LOGDIR / ("mut-%s.txt" % tag)).write_text(out, encoding="utf-8")
    return p.returncode, out


def main() -> int:
    pristine = TARGET.read_bytes()
    digest_before = hashlib.md5(pristine).hexdigest()
    text = pristine.decode("utf-8")
    if text.count(ORIG) != 1:
        print("ABORT: mutation anchor not found exactly once (count=%d)" % text.count(ORIG))
        return 3
    try:
        rc_clean, out_clean = run_selftest("00-pristine")
        TARGET.write_text(text.replace(ORIG, MUT), encoding="utf-8")
        rc_mut, out_mut = run_selftest("01-class3-exemption-removed")
    finally:
        TARGET.write_bytes(pristine)
    digest_after = hashlib.md5(TARGET.read_bytes()).hexdigest()

    def ring_lines(out: str) -> list[str]:
        return [l for l in out.splitlines() if "phantom-citation" in l and ("FAIL" in l or "silent OK" in l)]

    print("pristine md5 before = %s" % digest_before)
    print("pristine md5 after  = %s" % digest_after)
    print("restored identical  = %s" % (digest_before == digest_after))
    print("self-test rc pristine = %d" % rc_clean)
    print("self-test rc mutated  = %d" % rc_mut)
    print("pristine ban #9 lines:")
    for l in ring_lines(out_clean):
        print("  " + l[:150])
    print("mutated ban #9 lines:")
    for l in ring_lines(out_mut):
        print("  " + l[:150])
    print("pristine tail: " + [l for l in out_clean.splitlines() if "direction checks" in l or "direction(s) failed" in l][-1])
    print("mutated tail: " + [l for l in out_mut.splitlines() if "direction checks" in l or "direction(s) failed" in l][-1])
    return 0


if __name__ == "__main__":
    sys.exit(main())
