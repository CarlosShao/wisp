#!/usr/bin/env python3
"""260-r2 mutation driver: flip one defined line, run the default-档 suite, restore.

Usage: python mutate260r2.py <M1|M2|M3>
Each mutation is a single literal string replacement in a single file, restored by
the same script. The caller verifies the md5 of both files against
md5-before.txt / md5-restored.txt.
"""
import subprocess
import sys

ROOT = r"D:\work\workspace\projects plans\Wisp"
EXPECT = r"internal\ball\hotkey_cancel_borrow_expect_260r2_test.go"
PROD = r"internal\ball\hotkey_windows.go"

MUTATIONS = {
    # M1: the expectation stops following the seed - it goes back to being the
    # pre-260 constant. This is the exact shape this round replaced.
    "M1": (EXPECT,
           "\treturn borrowWish260r2{Seed: bind, Acc: acc, Default: acc == wantDefaultBorrow260r2}",
           "\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true}"),
    # M2: production goes back to ignoring [hotkey] cancel (the ticket's defect,
    # the pre-260 shape), while the ruler still expects the configured key.
    "M2": (PROD,
           "\treturn cancelBorrow{Binding: bind, Acc: acc}",
           "\treturn cancelBorrow{Binding: bind, Acc: escBorrowAcc()}"),
    # M3: a refused seed stops being reported - the expectation would then nod at
    # the default borrow without naming the substitution (condition ③).
    "M3": (EXPECT,
           "\t\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true, Refused: err}",
           "\t\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true}"),
}


def main() -> int:
    key = sys.argv[1]
    rel, old, new = MUTATIONS[key]
    path = ROOT + "\\" + rel
    with open(path, encoding="utf-8") as fh:
        src = fh.read()
    if src.count(old) != 1:
        print("%s: anchor found %d times, refusing to mutate" % (key, src.count(old)))
        return 2
    with open(path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(src.replace(old, new, 1))
    print("== %s applied to %s ==" % (key, rel))
    run = subprocess.run(
        ["go", "test", "-count=1", "-v", "./internal/ball/"],
        cwd=ROOT, capture_output=True, text=True, encoding="utf-8", errors="replace")
    lines = [l for l in (run.stdout + run.stderr).splitlines()
             if l.startswith("--- FAIL") or l.startswith("--- PASS") or l.startswith("    ")]
    for l in lines:
        if "FAIL" in l or "260r2" in l or "260" in l:
            print(l)
    print("package exit=%d" % run.returncode)
    with open(path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(src)
    print("== %s restored ==" % key)
    return 0


if __name__ == "__main__":
    sys.exit(main())
