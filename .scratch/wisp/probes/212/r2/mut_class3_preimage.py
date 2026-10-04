"""212-r2: second half of the mutation self-proof, run against the PRE-IMAGE tree.

Same patch as mut_class3_exemption.py (ban #9's class-3 exemption neutered), but
the scan target is the extracted `5e8748b3^` tree instead of the self-test
fixture, so the two readings together show what 裁 (a) costs in findings:

  fixed scanner    -> N phantom tokens on the pre-image tree
  exempted removed -> N+? phantom tokens, the difference being the abbreviated
                      spellings the pre-image comments actually used

The tree is a `git archive` extraction in the OS temp dir (never inside the
repository), so untracked-but-present files are missing from it; that is a
method artifact named in .scratch/wisp/probes/212/r2/fix-and-readings.md.
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
DEST = pathlib.Path(sys.argv[1])


def scan(tag: str) -> str:
    p = subprocess.run(
        ["go", "run", ".", "-root", str(DEST)],
        cwd=str(ROOT / "tools" / "d22scan"),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    out = "rc=%d\n--- stdout ---\n%s\n--- stderr ---\n%s\n" % (p.returncode, p.stdout, p.stderr)
    (LOGDIR / ("preimage-tree-scan-%s.txt" % tag)).write_text(out, encoding="utf-8")
    toks = sorted({
        l.split('repo path "')[1].split('"')[0]
        for l in out.splitlines() if 'phantom-citation] comment cites repo path "' in l
    })
    return "rc=%d distinct-phantom-tokens=%d\n%s" % (p.returncode, len(toks), "\n".join(toks))


def main() -> int:
    pristine = TARGET.read_bytes()
    before = hashlib.md5(pristine).hexdigest()
    text = pristine.decode("utf-8")
    if text.count(ORIG) != 1:
        print("ABORT: mutation anchor count=%d" % text.count(ORIG))
        return 3
    try:
        fixed = scan("fixed-here")
        TARGET.write_text(text.replace(ORIG, MUT), encoding="utf-8")
        mutated = scan("mutated-exemption-removed")
    finally:
        TARGET.write_bytes(pristine)
    after = hashlib.md5(TARGET.read_bytes()).hexdigest()
    (LOGDIR / "preimage-comparison.txt").write_text(
        "fixed:\n%s\n\nmutated (class-3 exemption removed):\n%s\n\nmd5 before=%s after=%s identical=%s\n"
        % (fixed, mutated, before, after, before == after), encoding="utf-8")
    print(fixed)
    print("-----")
    print(mutated)
    print("-----")
    print("md5 before=%s after=%s identical=%s" % (before, after, before == after))
    return 0


if __name__ == "__main__":
    sys.exit(main())
