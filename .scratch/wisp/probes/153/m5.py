"""Ticket 153 probe helper: apply / revert the M5 mutation from 139 acceptance
r1 s2.2 inside an OUT-OF-REPO snapshot.

M5 = replace the fold guard `if rep.Ran {` with `if c.Need(hist) {`, i.e. the
"print whenever the threshold was crossed" shape.

Files are opened with newline="" so the round trip is byte-exact: a mutation run
must not leave a line-ending change behind, or "restored identical" is a lie.

Usage: python m5.py apply|revert|append <snapshot-root> [payload-file]
"""
import sys
from pathlib import Path

TARGET = "internal/agent/compress.go"
OLD = "\tif rep.Ran {\n"
NEW = "\tif c.Need(hist) {\n"


def read(p: Path) -> str:
    with p.open(encoding="utf-8", newline="") as f:
        return f.read()


def write(p: Path, s: str) -> None:
    with p.open("w", encoding="utf-8", newline="") as f:
        f.write(s)


def main() -> int:
    mode, root = sys.argv[1], Path(sys.argv[2])
    if mode == "append":
        dst = root / "internal/agent/compress_trace_test.go"
        payload = Path(sys.argv[3]).read_text(encoding="utf-8")
        src = read(dst)
        marker = "func TestCompressionTraceSilentWhenNothingFoldableOverThreshold"
        if marker in src:
            print("append refused: payload already present")
            return 2
        write(dst, src + "\n" + payload)
        print(f"appended {len(payload.splitlines())} lines -> compress_trace_test.go "
              f"({len(read(dst).splitlines())} lines total)")
        return 0
    p = root / TARGET
    src = read(p)
    if mode == "apply":
        if NEW in src or src.count(OLD) != 1:
            print(f"apply refused: guard occurrences = {src.count(OLD)}")
            return 2
        write(p, src.replace(OLD, NEW, 1))
        print("M5 applied: `if rep.Ran {` -> `if c.Need(hist) {`")
        for n, ln in enumerate(read(p).splitlines(), 1):
            if "if rep.Ran {" in ln or "if c.Need(hist) {" in ln:
                print(f"  landed at compress.go:{n}: {ln.strip()}")
        return 0
    if mode == "revert":
        if NEW not in src:
            print("revert refused: mutated guard not present")
            return 2
        write(p, src.replace(NEW, OLD, 1))
        print("M5 reverted: guard back to `if rep.Ran {`")
        return 0
    print("unknown mode")
    return 2


if __name__ == "__main__":
    sys.exit(main())
