#!/usr/bin/env python3
"""260-r2 mutation driver, second attempt (v2).

mutate260r2.py M1/M2 were REJECTED by their own run: both replacements left a local
variable unused, so `go test` died at build time instead of at an assertion - that
is a compile error, not a red ruler, and it proves nothing. v2 keeps every parsed
value used and forces only the VALUE under test, so each mutation has to be caught
by an assertion or not at all. The first driver stays on disk (repo rule: 临时件
只建不删); its readings are in mutation-M1M2M3.txt.

Usage: python mutate260r2b.py <N1|N2|N3|N4>   (full go test -v output -> mutation-<key>-run.txt)
"""
import subprocess
import sys

ROOT = r"D:\work\workspace\projects plans\Wisp"
EXPECT = r"internal\ball\hotkey_cancel_borrow_expect_260r2_test.go"
PROD = r"internal\ball\hotkey_windows.go"

MUTATIONS = {
    # N1: the expectation stops following the seed and always names the default
    # pair (the pre-260 constant this round removed), while still consulting the
    # parser for the Default flag so the file keeps compiling.
    "N1": (EXPECT,
           "\treturn borrowWish260r2{Seed: bind, Acc: acc, Default: acc == wantDefaultBorrow260r2}",
           "\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: acc == wantDefaultBorrow260r2}"),
    # N2: production's configured branch forces VK_ESCAPE (the ticket's defect:
    # the configured key never reaches RegisterHotKey) but keeps the parsed mods.
    "N2": (PROD,
           "\treturn cancelBorrow{Binding: bind, Acc: acc}",
           "\treturn cancelBorrow{Binding: bind, Acc: Accelerator{Mods: acc.Mods, VK: vkEscape}}"),
    # N3: a refused seed stops being reported, so the ruler would nod at the
    # default borrow without naming the substitution (condition ③).
    "N3": (EXPECT,
           "\t\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true, Refused: err}",
           "\t\treturn borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true}"),
    # N4: the default cell is NOT the shipped pair any more - condition ①'s zero
    # drift has to be what goes red here, not the configured cell.
    "N4": (EXPECT,
           "var wantDefaultBorrow260r2 = Accelerator{Mods: 0x4000, VK: 0x1B}",
           "var wantDefaultBorrow260r2 = Accelerator{Mods: 0x4000, VK: 0x1C}"),
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
    with open(ROOT + r"\.scratch\wisp\probes\260\r2\mutation-%s-run.txt" % key, "w",
              encoding="utf-8", newline="\n") as fh:
        fh.write("== %s: %s ==\n" % (key, rel))
        fh.write("command: go test -count=1 -v ./internal/ball/\n")
        fh.write("package exit=%d\n\n" % run.returncode)
        fh.write(run.stdout or "")
        fh.write("\n--- stderr ---\n")
        fh.write(run.stderr or "")
    fails = [l for l in (run.stdout or "").splitlines() if l.startswith("--- FAIL")]
    print("package exit=%d" % run.returncode)
    print("failing tests: %s" % (", ".join(fails) if fails else "NONE"))
    with open(path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(src)
    print("== %s restored ==" % key)
    return 0


if __name__ == "__main__":
    sys.exit(main())
