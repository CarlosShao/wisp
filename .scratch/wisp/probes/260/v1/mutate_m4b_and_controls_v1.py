#!/usr/bin/env python3
"""260-v1 mutation rig: M4b (read the WRONG receipt line) + CTRL (drop the
assembly-root injection) + CTRL2 (make the reader answer a fixed literal).

Every mutation is applied to the tracked file, tested, then restored from the
bytes held in memory; the restored md5 is compared against the starting md5 and
printed. ⛔ never leaves the tracked file mutated.
"""
import hashlib
import os
import subprocess
import sys

ROOT = os.path.abspath(os.getcwd())
TARGET = os.path.join(ROOT, "cmd", "wisp", "resident_approval_windows.go")
LOGDIR = os.path.join(ROOT, ".scratch", "wisp", "probes", "260", "v1", "logs")

RECEIPT_LOOP = """	if b != nil {
		for _, bd := range b.HotkeyReport().Bindings() {
			if bd.Name == "cancel" && bd.Binding != "" {
				return bd.Binding
			}
		}
	}
"""

INJECT = "\tapproval.SetCancelKeySpelling(ra.cancelKeySpelling)\n"


def md5(path):
    with open(path, "rb") as fh:
        return hashlib.md5(fh.read()).hexdigest()


def run(tag, full):
    args = ["go", "test", "-count=1", "-timeout", "900s"]
    if not full:
        args += ["-v", "-run", "TestTicket260"]
    args += ["./cmd/wisp/"]
    proc = subprocess.run(args, cwd=ROOT, capture_output=True, text=True,
                          encoding="utf-8", errors="replace")
    out = (proc.stdout or "") + "\n--- STDERR ---\n" + (proc.stderr or "")
    name = "%s-run-%s.txt" % (tag, "full" if full else "targeted")
    with open(os.path.join(LOGDIR, name), "w", encoding="utf-8") as fh:
        fh.write("rc=%d\n%s" % (proc.returncode, out))
    lines = [l for l in out.splitlines()
             if l.startswith("--- ") or l.startswith("ok") or l.startswith("FAIL")
             or l.startswith("PASS") or "exit status" in l]
    print("[%s / %s] rc=%d" % (tag, "full" if full else "targeted", proc.returncode))
    for l in lines:
        print("   " + l)
    return proc.returncode


def apply_and_test(tag, mutated, full):
    before = md5(TARGET)
    with open(TARGET, "w", encoding="utf-8", newline="") as fh:
        fh.write(mutated)
    print("  %s: mutated md5=%s" % (tag, md5(TARGET)))
    try:
        return run(tag, full)
    finally:
        src = load()
        with open(TARGET, "w", encoding="utf-8", newline="") as fh:
            fh.write(src)
        after = md5(TARGET)
        print("  %s: restored md5=%s matches_start=%s" % (tag, after, after == before))


def load():
    with open(TARGET, "r", encoding="utf-8", newline="") as fh:
        return fh.read()


def main():
    src = load()
    start = md5(TARGET)
    print("starting md5 =", start)
    if RECEIPT_LOOP not in src or INJECT not in src:
        print("SETUP RED: anchors not found, refusing to guess")
        return 2

    # M4b: the read step stays, but it reads the WRONG receipt line. If any ruler
    # actually watched the ball's cancel line, naming the summon line must move it.
    m4b = src.replace(RECEIPT_LOOP, RECEIPT_LOOP.replace('bd.Name == "cancel"', 'bd.Name == "summon"'))
    apply_and_test("M4b-wrong-receipt-line", m4b, full=False)
    apply_and_test("M4b-wrong-receipt-line", m4b, full=True)

    # CTRL: drop the assembly-root injection (r4 claims this one goes red).
    ctrl = src.replace(INJECT, "\t// CTRL: injection removed by 260-v1 rig\n")
    apply_and_test("CTRL-no-assembly-injection", ctrl, full=False)

    # CTRL2: re-hardcode "Esc" in ONE of the two veto sentences. This is the shape
    # AC#4 ② exists to catch, so a rig that cannot see it is a blind rig and its
    # green on M4b would prove nothing. (r3's cell ② claims the teeth; this is the
    # 260-v1 re-check of that claim, not a restatement of it.)
    OLD_DONE = '\t\tkey, card.CorrelationID, card.Level, card.Tool)'
    ctrl2 = src.replace(OLD_DONE, '\t\t"Esc", card.CorrelationID, card.Level, card.Tool)')
    if ctrl2 == src:
        print("  CTRL2 SETUP RED: vetoDoneLine argument anchor not found")
    else:
        apply_and_test("CTRL2-rehardcoded-esc-in-vetoDoneLine", ctrl2, full=False)

    print("final md5 =", md5(TARGET), "matches_start =", md5(TARGET) == start)
    return 0


if __name__ == "__main__":
    sys.exit(main())
