#!/usr/bin/env python3
"""Ticket 260 AC#2, leg 260-r5 (2026-10-05 recheck): the sensitivity-boundary probe.

Question this answers by measurement instead of by inheriting instrument.md: if
ONLY the registration-set mirror write at ball_windows.go:896 is deleted (the
report line at :895 and the return stay), does this leg's ruler go red?

Expected answer, stated before running: NO - the seed set never holds the cancel
slot idle, and a refused borrow must not hold it either, so liveAfter reads the
same thing with or without that line. That is the honest boundary of what the
wiring-layer ruler covers, and it belongs in section 5 of the evidence file.

Same hard rules as the other rigs: restore comes from a HEAD copy taken before
anything is touched, md5 printed at every step, start interlock refuses to run
if the working file is not HEAD.
"""

import hashlib
import os
import subprocess
import sys

REPO = r"D:\work\workspace\projects plans\Wisp"
TARGET_REL = "internal/ball/ball_windows.go"
TARGET = os.path.join(REPO, TARGET_REL)
LOGS = os.path.join(REPO, ".scratch", "wisp", "probes", "260", "r5", "logs-recheck")
HEAD_BLOB = os.path.join(LOGS, "ball_windows.go.head-blob")

CHECK_OPEN = b"\t\tif err := takeEscBorrow(b.hwnd, borrow); err != nil {\n"
COMMENT = (
    b"\t\t\t// The borrow was refused: say it on the report (Problems() names it)\n"
    b"\t\t\t// rather than leaving a card with no cancel key and no explanation.\n"
)
WRITE_LINE = b"\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
MIRROR_LINE = b"\t\t\tb.registeredHotkeys = b.hotkeyReport.Live()\n"
RETURN_LINE = b"\t\t\treturn\n"
CLOSE = b"\t\t}\n"
IF_BLOCK = CHECK_OPEN + COMMENT + WRITE_LINE + MIRROR_LINE + RETURN_LINE + CLOSE
BLOCK_NO_MIRROR = CHECK_OPEN + COMMENT + WRITE_LINE + RETURN_LINE + CLOSE


def md5_bytes(b):
    return hashlib.md5(b).hexdigest()


def md5_file(p):
    with open(p, "rb") as f:
        return md5_bytes(f.read())


def read(p):
    with open(p, "rb") as f:
        return f.read()


def write(p, data):
    with open(p, "wb") as f:
        f.write(data)


def stamp():
    return subprocess.run(["date", "+%Y-%m-%d %H:%M:%S %z"], capture_output=True, text=True).stdout.strip()


def main():
    os.makedirs(LOGS, exist_ok=True)
    lines = []

    def say(m):
        print(m)
        lines.append(m)

    say("# 260-r5 M4 sensitivity-boundary probe, started %s" % stamp())
    blob = subprocess.run(["git", "cat-file", "blob", "HEAD:" + TARGET_REL], cwd=REPO, capture_output=True)
    if blob.returncode != 0:
        say("FATAL: cat-file failed")
        return 1
    write(HEAD_BLOB, blob.stdout)
    base = read(HEAD_BLOB)
    start = md5_file(TARGET)
    say("start: md5(working)=%s md5(HEAD blob)=%s matches_start=%s" % (start, md5_bytes(base), start == md5_bytes(base)))
    if start != md5_bytes(base):
        say("FATAL: working file is not HEAD - another leg's edit is in flight. Refusing.")
        return 1

    cur = read(HEAD_BLOB)
    if cur.count(IF_BLOCK) != 1:
        say("FATAL: if-block found %d times" % cur.count(IF_BLOCK))
        return 1
    mutated = cur.replace(IF_BLOCK, BLOCK_NO_MIRROR, 1)
    if mutated == cur:
        say("FATAL: no-op")
        return 1
    write(TARGET, mutated)
    say("== M4-drop-mirror-only applied at %s md5(mutated)=%s matches_start=%s"
        % (stamp(), md5_file(TARGET), md5_file(TARGET) == start))
    began = stamp()
    proc = subprocess.run(["go", "test", "-count=1", "-v", "-run", "260r5", "./internal/ball/"], cwd=REPO, capture_output=True)
    ended = stamp()
    out = os.path.join(LOGS, "mutation-M4-drop-mirror-only-go-test-v.txt")
    with open(out, "wb") as f:
        f.write(("#### mutation=M4-drop-mirror-only  window %s -> %s  (targeted -run 260r5)\n" % (began, ended)).encode())
        f.write(proc.stdout)
        f.write(b"\n===== stderr =====\n")
        f.write(proc.stderr)
    say("   go test rc=%d  window=%s -> %s  reading=%s" % (proc.returncode, began, ended, os.path.relpath(out, REPO)))
    write(TARGET, base)
    now = md5_file(TARGET)
    say("   restored at %s: md5=%s matches_start=%s" % (stamp(), now, now == start))
    say("FINAL md5=%s matches_start=%s" % (md5_file(TARGET), md5_file(TARGET) == start))
    with open(os.path.join(LOGS, "mutation-md5-ledger.txt"), "a", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")
    return 0 if now == start else 1


if __name__ == "__main__":
    sys.exit(main())
