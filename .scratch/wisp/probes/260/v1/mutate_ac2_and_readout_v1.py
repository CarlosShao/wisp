#!/usr/bin/env python3
"""260-v1 driver v2 (save-restore fixed).

Saves the ORIGINAL bytes of every file it may touch BEFORE the first write and
restores from those saved bytes in a finally block - the v1 rig re-read the file
inside the finally block and therefore restored the mutant (booked in verdict
section 9). Here the starting md5 is asserted at the end and printed.

Mutations:
  A  internal/ball/ball_windows.go  TakeEscForCancel: drop the failed-borrow
     receipt write  -> does any non-winelive ruler notice? (AC#2)
  B  internal/ball/hotkey_windows.go  Problems(): stop naming an attempted-but-
     failed line  -> CONTROL, must go red (proves ruler A is not blind)
  C  cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go: break the
     default-档 constant so the red line PRINTS the shipped sentence (read-out,
     not an assertion change - restored immediately)
"""
import hashlib
import os
import subprocess
import sys

ROOT = os.path.abspath(os.getcwd())
LOGDIR = os.path.join(ROOT, ".scratch", "wisp", "probes", "260", "v1", "logs")
BALL_WIN = os.path.join(ROOT, "internal", "ball", "ball_windows.go")
BALL_HOT = os.path.join(ROOT, "internal", "ball", "hotkey_windows.go")
R3_TEST = os.path.join(ROOT, "cmd", "wisp", "resident_cancel_key_wording_260r3_windows_test.go")

FAILED_WRITE = """		if err := takeEscBorrow(b.hwnd, borrow); err != nil {
			// The borrow was refused: say it on the report (Problems() names it)
			// rather than leaving a card with no cancel key and no explanation.
			b.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))
			b.registeredHotkeys = b.hotkeyReport.Live()
			return
		}
"""
FAILED_WRITE_MUT = """		if err := takeEscBorrow(b.hwnd, borrow); err != nil {
			// 260-v1 A: the receipt write is deleted (silent borrow failure)
			return
		}
"""

PROBLEM_DEFAULT = """		default:
			out = append(out, fmt.Sprintf(
				"hotkey %s = %q was not registered: %v", b.Name, b.Binding, b.Err))
"""
PROBLEM_DEFAULT_MUT = """		default:
			// 260-v1 B: an attempted-but-refused line stops being reported
			continue
"""

R3_OLD_DONE = '\told260r3DoneSaid     = "按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计"'
R3_OLD_DONE_MUT = '\told260r3DoneSaid     = "260-v1 READ-OUT SENTINEL"'


def md5(p):
    with open(p, "rb") as fh:
        return hashlib.md5(fh.read()).hexdigest()


def load(p):
    with open(p, "r", encoding="utf-8", newline="") as fh:
        return fh.read()


def run(tag, args, pkg):
    cmd = ["go", "test", "-count=1", "-timeout", "900s"] + args + [pkg]
    proc = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True,
                          encoding="utf-8", errors="replace")
    out = (proc.stdout or "") + "\n--- STDERR ---\n" + (proc.stderr or "")
    with open(os.path.join(LOGDIR, tag + ".txt"), "w", encoding="utf-8") as fh:
        fh.write("rc=%d\n$ %s\n%s" % (proc.returncode, " ".join(cmd), out))
    print("[%s] rc=%d" % (tag, proc.returncode))
    for l in out.splitlines():
        if l.startswith("--- FAIL") or l.startswith("--- SKIP") or l.startswith("    ") \
                or l.startswith("ok ") or l.startswith("FAIL") or "exit status" in l:
            if l.startswith("    ") and ("level=" in l or "msg=" in l):
                continue
            print("   " + l[:400])
    return proc.returncode, out


def main():
    originals = {p: load(p) for p in (BALL_WIN, BALL_HOT, R3_TEST)}
    start = {p: md5(p) for p in originals}
    for p, m in start.items():
        print("start md5 %s %s" % (m, os.path.basename(p)))

    ok = True
    for p, needle in ((BALL_WIN, FAILED_WRITE), (BALL_HOT, PROBLEM_DEFAULT), (R3_TEST, R3_OLD_DONE)):
        n = originals[p].count(needle)
        if n != 1:
            print("SETUP RED: anchor count=%d (want exactly 1) in %s" % (n, p))
            ok = False
    if not ok:
        return 2

    try:
        # A: silent borrow failure at the production wiring
        a = originals[BALL_WIN].replace(FAILED_WRITE, FAILED_WRITE_MUT, 1)
        with open(BALL_WIN, "w", encoding="utf-8", newline="") as fh:
            fh.write(a)
        print("  A applied md5=%s" % md5(BALL_WIN))
        run("A-silent-borrow-failure-ball", [], "./internal/ball/")
        with open(BALL_WIN, "w", encoding="utf-8", newline="") as fh:
            fh.write(originals[BALL_WIN])

        # B: control - Problems() stops naming the attempted/failed line
        b = originals[BALL_HOT].replace(PROBLEM_DEFAULT, PROBLEM_DEFAULT_MUT, 1)
        with open(BALL_HOT, "w", encoding="utf-8", newline="") as fh:
            fh.write(b)
        print("  B applied md5=%s" % md5(BALL_HOT))
        run("B-control-problems-goes-silent", ["-run", "TestCancelBorrowFailureIsAProblemLine|TestBorrow"], "./internal/ball/")
        with open(BALL_HOT, "w", encoding="utf-8", newline="") as fh:
            fh.write(originals[BALL_HOT])

        # C: read-out of the shipped default-档 sentence
        c = originals[R3_TEST].replace(R3_OLD_DONE, R3_OLD_DONE_MUT, 1)
        with open(R3_TEST, "w", encoding="utf-8", newline="") as fh:
            fh.write(c)
        print("  C applied md5=%s" % md5(R3_TEST))
        run("C-readout-default-sentence", ["-v", "-run", "TestTicket260R3DefaultWordingIsTheOldSentence"], "./cmd/wisp/")
        with open(R3_TEST, "w", encoding="utf-8", newline="") as fh:
            fh.write(originals[R3_TEST])
    finally:
        bad = []
        for p, src in originals.items():
            with open(p, "w", encoding="utf-8", newline="") as fh:
                fh.write(src)
            now = md5(p)
            print("restore %s md5=%s matches_start=%s" % (os.path.basename(p), now, now == start[p]))
            if now != start[p]:
                bad.append(p)
        if bad:
            print("!!! RESTORE FAILED for: %s" % bad)
            return 3
        print("all tracked files back at starting md5")
    return 0


if __name__ == "__main__":
    sys.exit(main())
