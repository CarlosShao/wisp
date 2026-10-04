#!/usr/bin/env python3
"""Ticket 260 AC#2, leg 260-r5: the mutation rig for the refused-borrow ruler.

Two hard rules this script is written around (the repo nearly lost production
code to the second one today):

  1. The restore source is a copy pulled out of HEAD BEFORE any mutation
     ("git cat-file blob HEAD:internal/ball/ball_windows.go"). Nothing here ever
     re-reads the working file and writes it back.
  2. Every write prints its md5 next to the start value, and the final md5 has
     to equal the start value or the script exits non-zero.

Mutations (each is an exact-string replace, and the script refuses to run if the
fragment is not found exactly once):

  M1  the ticket's own positive control: delete the failure-receipt write at
      ball_windows.go:895. A ruler that stays green here is not a ruler.
  M2  the receipt is CANNED: the line is written with an error this file made up
      instead of the errno user32 returned. A ruler that only asks "is there a
      failure line" passes; this one must not.
  M2b the receipt names the WRONG KEY: production's default-档 line
      (cancelFailedLine) instead of the line built from the borrow that was
      really attempted (cancelFailedLineFor).
  M3  the refusal is not a door: the "return" is dropped, so the ball falls
      through and reports the key as borrowed anyway.
"""

import hashlib
import os
import subprocess
import sys

REPO = r"D:\work\workspace\projects plans\Wisp"
TARGET_REL = "internal/ball/ball_windows.go"
TARGET = os.path.join(REPO, TARGET_REL)
LOGS = os.path.join(REPO, ".scratch", "wisp", "probes", "260", "r5", "logs")
HEAD_BLOB = os.path.join(LOGS, "ball_windows.go.head-blob")

WRITE_LINE = b"\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
IF_BLOCK = (
    b"\t\tif err := takeEscBorrow(b.hwnd, borrow); err != nil {\n"
    b"\t\t\t// The borrow was refused: say it on the report (Problems() names it)\n"
    b"\t\t\t// rather than leaving a card with no cancel key and no explanation.\n"
    + WRITE_LINE
    + b"\t\t\tb.registeredHotkeys = b.hotkeyReport.Live()\n"
    b"\t\t\treturn\n"
    b"\t\t}\n"
)

MUTATIONS = [
    ("M1-drop-receipt-write", IF_BLOCK, IF_BLOCK.replace(WRITE_LINE, b"")),
    (
        "M2-canned-error-receipt",
        WRITE_LINE,
        b"\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, "
        b'fmt.Errorf("260-r5 M2: a canned refusal this file invented")))\n',
    ),
    ("M2b-canned-default-key-line", WRITE_LINE, b"\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLine(err))\n"),
    ("M3-no-return-after-refusal", IF_BLOCK, IF_BLOCK.replace(b"\t\t\treturn\n", b"")),
]


def md5_bytes(b: bytes) -> str:
    return hashlib.md5(b).hexdigest()


def md5_file(p: str) -> str:
    with open(p, "rb") as f:
        return md5_bytes(f.read())


def read(p: str) -> bytes:
    with open(p, "rb") as f:
        return f.read()


def write(p: str, data: bytes) -> None:
    with open(p, "wb") as f:
        f.write(data)


def stamp() -> str:
    return subprocess.run(["date", "+%Y-%m-%d %H:%M:%S %z"], capture_output=True, text=True).stdout.strip()


def run_tests(tag: str) -> tuple:
    """Full go test -v of internal/ball, saved under logs/. Returns (rc, path)."""
    out_path = os.path.join(LOGS, "mutation-%s-go-test-v.txt" % tag)
    proc = subprocess.run(
        ["go", "test", "-count=1", "-v", "./internal/ball/"],
        cwd=REPO,
        capture_output=True,
    )
    with open(out_path, "wb") as f:
        f.write(proc.stdout)
        f.write(b"\n===== stderr =====\n")
        f.write(proc.stderr)
    return proc.returncode, out_path


def main() -> int:
    os.makedirs(LOGS, exist_ok=True)
    lines = []

    def say(msg: str) -> None:
        print(msg)
        lines.append(msg)

    say("# 260-r5 mutation rig, started %s" % stamp())

    # Rule 1: the restore source comes out of HEAD, before anything is touched.
    blob = subprocess.run(
        ["git", "cat-file", "blob", "HEAD:" + TARGET_REL],
        cwd=REPO,
        capture_output=True,
    )
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
        rc, path = run_tests(name)
        say("   go test rc=%d  reading=%s" % (rc, os.path.relpath(path, REPO)))
        # Rule 1 again: restore from the HEAD copy, never from the current file.
        write(TARGET, base)
        now = md5_file(TARGET)
        say("   md5(after restore)=%s matches_start=%s" % (now, now == start_md5))
        results.append((name, "applied", now, rc))

    say("")
    say("== restored control run (no mutation in the file) at %s" % stamp())
    rc, path = run_tests("restored")
    say("   go test rc=%d (rc=1 is the known environment red TestC21TableColourRowsMatchTokensCSS) reading=%s"
        % (rc, os.path.relpath(path, REPO)))

    final = md5_file(TARGET)
    say("")
    say("FINAL md5(working ball_windows.go)=%s matches_start=%s" % (final, final == start_md5))
    for name, now, rc in [(a, c, d) for a, b, c, d in results]:
        say("   %s: md5-after-restore=%s test-rc=%s" % (name, now, rc))
    with open(os.path.join(LOGS, "mutation-md5-ledger.txt"), "a", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")
    return 0 if final == start_md5 else 1


if __name__ == "__main__":
    sys.exit(main())
