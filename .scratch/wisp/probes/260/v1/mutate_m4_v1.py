#!/usr/bin/env python3
"""260-v1 M4 mutation driver.

Mutates cmd/wisp/resident_approval_windows.go's residentCancelKeySpelling so the
"read the ball's hotkey receipt" step is REMOVED (the function answers straight
from the ball's default constant). Runs the named rulers, then restores the file
from the bytes this script saved before mutating, and prints the md5 check.

⛔ never leaves the tracked file in the mutated state: the restore runs in a
finally block and the md5 of the restored file is compared to the starting md5.
"""
import hashlib
import os
import subprocess
import sys

# this script is always run with the repo root as cwd (see verdict.md §2 for the
# exact command line); os.getcwd() beats a depth-guessed __file__ walk here.
ROOT = os.path.abspath(os.getcwd())
TARGET = os.path.join(ROOT, "cmd", "wisp", "resident_approval_windows.go")
LOGDIR = os.path.join(ROOT, ".scratch", "wisp", "probes", "260", "v1", "logs")

ORIGINAL_BODY = """func residentCancelKeySpelling(b *ball.Ball) string {
	if b != nil {
		for _, bd := range b.HotkeyReport().Bindings() {
			if bd.Name == "cancel" && bd.Binding != "" {
				return bd.Binding
			}
		}
	}
	return ball.DefaultHotkeys().Cancel
}
"""

MUTANT_BODY = """func residentCancelKeySpelling(b *ball.Ball) string {
	_ = b
	return ball.DefaultHotkeys().Cancel
}
"""


def md5(path):
    with open(path, "rb") as fh:
        return hashlib.md5(fh.read()).hexdigest()


def main():
    with open(TARGET, "r", encoding="utf-8", newline="") as fh:
        src = fh.read()
    if ORIGINAL_BODY not in src:
        print("M4 SETUP RED: the exact receipt-reading body was not found; refusing to guess.")
        return 2
    before = md5(TARGET)
    with open(os.path.join(LOGDIR, "M4-md5-before.txt"), "w", encoding="utf-8") as fh:
        fh.write("%s  %s\n" % (before, TARGET))

    mutant = src.replace(ORIGINAL_BODY, MUTANT_BODY)
    rc = 0
    try:
        with open(TARGET, "w", encoding="utf-8", newline="") as fh:
            fh.write(mutant)
        print("M4 mutated: the receipt-read step is gone; md5 now", md5(TARGET))
        for name, args in (
            ("M4-run-targeted-260", ["go", "test", "-count=1", "-v", "-run", "TestTicket260", "./cmd/wisp/"]),
            ("M4-run-full-cmdwisp", ["go", "test", "-count=1", "./cmd/wisp/"]),
        ):
            proc = subprocess.run(args + ["-timeout", "600s"], cwd=ROOT,
                                  capture_output=True, text=True, encoding="utf-8", errors="replace")
            out = (proc.stdout or "") + "\n--- STDERR ---\n" + (proc.stderr or "")
            with open(os.path.join(LOGDIR, name + ".txt"), "w", encoding="utf-8") as fh:
                fh.write("rc=%d\n%s" % (proc.returncode, out))
            print("[%s] rc=%d" % (name, proc.returncode))
            for line in out.splitlines():
                if line.startswith("--- ") or line.startswith("ok") or line.startswith("FAIL") \
                        or line.startswith("PASS") or "exit status" in line:
                    print("   " + line)
            if proc.returncode != 0:
                rc = proc.returncode
    finally:
        with open(TARGET, "w", encoding="utf-8", newline="") as fh:
            fh.write(src)
        after = md5(TARGET)
        print("M4 restored: md5 %s ; matches starting md5 = %s" % (after, after == before))
        with open(os.path.join(LOGDIR, "M4-md5-restored.txt"), "w", encoding="utf-8") as fh:
            fh.write("%s  %s  matches_before=%s\n" % (after, TARGET, after == before))
    return rc


if __name__ == "__main__":
    sys.exit(main())
