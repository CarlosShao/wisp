#!/usr/bin/env python3
"""Ticket 260 AC#2, leg 260-r5 (re-dispatch of 2026-10-05): re-run the mutation.

Why this file exists next to mutate260r5.py instead of being that script:

  * The ticket's own rule ("reference it before citing, re-run it") means this
    leg has to produce the reading at TODAY's anchor, not quote yesterday's log.
  * mutate260r5.py writes its readings into logs/mutation-*-go-test-v.txt, and
    instrument.md section 5 cites those files verbatim. Re-running it would
    overwrite the earlier leg's evidence. This rig logs to logs-recheck/ and
    touches nothing in logs/.

Same two hard rules as the original rig:

  1. The restore source is a copy pulled out of HEAD before anything is touched
     (git cat-file blob HEAD:internal/ball/ball_windows.go, written to
     logs-recheck/ball_windows.go.head-blob). Nothing here re-reads the working
     file and writes it back.
  2. Every step prints md5 next to the start value; the start interlock refuses
     to run at all unless md5(working) == md5(HEAD blob), so another leg's edit
     in flight cannot be carried away.

Mutations (exact-string replaces; the rig refuses if the fragment is not found
exactly once):

  M1   the ticket's judgment criterion, same shape as the original rig: delete
       the failure-receipt WRITE (ball_windows.go:895) and keep everything else.
  M1b  the stronger reading of "delete the write at ball_windows.go:892-899":
       delete BOTH writes inside the refused-borrow block (:895 report line and
       :896 registration-set mirror), keep the error check and the return.
  M3   the refusal stops being a door: drop the "return" at :897 so the borrow
       is reported as successful anyway.

Each mutation runs the FULL package (go test -count=1 -v ./internal/ball/) so
the red roster is the whole package's, not a -narrowed one.
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
BLOCK_NO_WRITE = CHECK_OPEN + COMMENT + MIRROR_LINE + RETURN_LINE + CLOSE
BLOCK_NO_WRITES = CHECK_OPEN + COMMENT + RETURN_LINE + CLOSE
BLOCK_NO_RETURN = CHECK_OPEN + COMMENT + WRITE_LINE + MIRROR_LINE + CLOSE

MUTATIONS = [
    ("M1-drop-receipt-write", IF_BLOCK, BLOCK_NO_WRITE),
    ("M1b-drop-both-writes", IF_BLOCK, BLOCK_NO_WRITES),
    ("M3-no-return-after-refusal", IF_BLOCK, BLOCK_NO_RETURN),
]


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


def run_tests(tag):
    """Full go test -v of internal/ball, logged under logs-recheck/ with stamps."""
    out_path = os.path.join(LOGS, "mutation-%s-go-test-v.txt" % tag)
    began = stamp()
    proc = subprocess.run(["go", "test", "-count=1", "-v", "./internal/ball/"], cwd=REPO, capture_output=True)
    ended = stamp()
    with open(out_path, "wb") as f:
        f.write(("#### mutation=%s  go test started %s  ended %s\n" % (tag, began, ended)).encode())
        f.write(proc.stdout)
        f.write(b"\n===== stderr =====\n")
        f.write(proc.stderr)
    return proc.returncode, out_path, began, ended


def main():
    os.makedirs(LOGS, exist_ok=True)
    lines = []

    def say(msg):
        print(msg)
        lines.append(msg)

    say("# 260-r5 recheck rig, started %s" % stamp())

    blob = subprocess.run(["git", "cat-file", "blob", "HEAD:" + TARGET_REL], cwd=REPO, capture_output=True)
    if blob.returncode != 0:
        say("FATAL: git cat-file blob HEAD:%s failed: %s" % (TARGET_REL, blob.stderr.decode("utf-8", "replace")))
        return 1
    write(HEAD_BLOB, blob.stdout)
    base = read(HEAD_BLOB)

    start_md5 = md5_file(TARGET)
    head_md5 = md5_bytes(base)
    say("start: md5(working ball_windows.go)=%s" % start_md5)
    say("start: md5(HEAD blob copy)          =%s  head_blob=%s" % (head_md5, HEAD_BLOB))
    matches_start = start_md5 == head_md5
    say("start: matches_start=%s (working tree must equal HEAD before any mutation)" % matches_start)
    if not matches_start:
        say("FATAL: the working file is not the HEAD version - somebody else's edit is in flight. Refusing.")
        return 1

    results = []
    for name, old, new in MUTATIONS:
        cur = read(HEAD_BLOB)
        n = cur.count(old)
        if n != 1:
            say("FATAL: %s: fragment found %d times, expected exactly 1" % (name, n))
            results.append((name, "NOT-APPLIED", "", -1))
            continue
        mutated = cur.replace(old, new, 1)
        if mutated == cur:
            say("FATAL: %s: the mutation is a no-op" % name)
            results.append((name, "NO-OP", "", -1))
            continue
        write(TARGET, mutated)
        say("")
        say("== %s applied at %s" % (name, stamp()))
        say("   md5(mutated)=%s matches_start=%s (must be False: the file really is mutated)"
            % (md5_file(TARGET), md5_file(TARGET) == start_md5))
        rc, path, began, ended = run_tests(name)
        say("   go test rc=%d  window=%s -> %s  reading=%s" % (rc, began, ended, os.path.relpath(path, REPO)))
        write(TARGET, base)
        now = md5_file(TARGET)
        say("   restored from HEAD blob copy at %s: md5(after restore)=%s matches_start=%s" % (stamp(), now, now == start_md5))
        results.append((name, now, rc, None))

    say("")
    say("== restored control run (no mutation in the file) at %s" % stamp())
    rc, path, began, ended = run_tests("restored")
    say("   go test rc=%d (rc=1 is the known environment red TestC21TableColourRowsMatchTokensCSS)"
        "  window=%s -> %s  reading=%s" % (rc, began, ended, os.path.relpath(path, REPO)))

    final = md5_file(TARGET)
    say("")
    say("FINAL md5(working ball_windows.go)=%s matches_start=%s" % (final, final == start_md5))
    for row in results:
        say("   %s: md5-after-restore=%s test-rc=%s" % (row[0], row[1], row[2]))
    with open(os.path.join(LOGS, "mutation-md5-ledger.txt"), "a", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")
    return 0 if final == start_md5 else 1


if __name__ == "__main__":
    sys.exit(main())
