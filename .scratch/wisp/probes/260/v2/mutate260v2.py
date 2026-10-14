#!/usr/bin/env python3
"""Ticket 260 AC#2 acceptance leg 260-v2: counter-evidence mutations.

Four mutations, each applied alone against the same start bytes:

  M-A  drop the failure-receipt write (ball_windows.go, the ticket's own
       criterion shape; r5 planted this, re-planted here independently)
  M-B  keep the write but make it silent (receipt = cancelIdleLine) -- a
       shape neither r5 leg planted
  M-C  drop takeEscWithAcc's slog.Error (hotkey_windows.go) -- the second
       symptom channel; r5 never mutated this file in either leg
  M-D  keep the failed receipt but claim the key live in the registration
       set -- the stronger read of the :896 mirror the dispatch asked about
  M-E  drop only the :896 registration-set mirror (r5's M4) -- re-planted
       here so the sensitivity boundary in both directions is first-hand

Hygiene (per dispatch):
  - restore source = the out-of-repo backup taken BEFORE any mutation
    (C:\\tmp\\wisp260v2backup) AND the HEAD blob copies in logs/;
  - start interlock: sha256(working) == sha256(backup) == sha256(head blob)
    or the script refuses to run (someone else's byte is in the file);
  - sha256, not md5, for the byte-identity chain; md5 kept only as a
    cross-check against r5's own ledger;
  - every mutation window is timestamped start->end in the ledger.
"""
import hashlib
import os
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.normpath(os.path.join(HERE, "..", "..", "..", "..", ".."))
LOGS = os.path.join(HERE, "logs")
BACKUP = r"C:\tmp\wisp260v2backup"

BALL = os.path.join(ROOT, "internal", "ball", "ball_windows.go")
HOTKEY = os.path.join(ROOT, "internal", "ball", "hotkey_windows.go")
BALL_BLOB = os.path.join(LOGS, "ball_windows.go.head-blob")
HOTKEY_BLOB = os.path.join(LOGS, "hotkey_windows.go.head-blob")
BALL_BAK = os.path.join(BACKUP, "ball_windows.go.orig")
HOTKEY_BAK = os.path.join(BACKUP, "hotkey_windows.go.orig")

ENV = dict(os.environ, GOTMPDIR=r"C:\tmp\wisp260v2gotmp")

# ---------------------------------------------------------------- mutations

RECEIPT_WRITE = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
)
SILENT_RECEIPT = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
)
SILENT_DONE = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelIdleLine(b.cancelBinding))\n"
)
FAILED_BLOCK = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
    "\t\t\tb.registeredHotkeys = b.hotkeyReport.Live()\n"
    "\t\t\treturn\n"
)
FALSE_LIVE_BLOCK = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
    "\t\t\tb.registeredHotkeys = b.hotkeyReport.withCancel(cancelBorrowedLineFor(borrow)).Live()\n"
    "\t\t\treturn\n"
)
LOG_DROP = (
    '\t\tslog.Error("cancel takeover refused by Win32; Confirming will have no cancel key",\n'
    '\t\t\t"hotkey", "cancel", "binding", b.Binding, "err", err)\n'
)
MIRROR_DROP_BLOCK = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
    "\t\t\tb.registeredHotkeys = b.hotkeyReport.Live()\n"
)
MIRROR_DROP_KEEP = (
    "\t\t\tb.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLineFor(borrow, err))\n"
)

MUTATIONS = [
    # name, file, old, new, full-package?
    ("M-A-drop-receipt-write", BALL, RECEIPT_WRITE, "", True),
    ("M-B-silent-receipt", BALL, SILENT_RECEIPT, SILENT_DONE, False),
    ("M-C-drop-refusal-log", HOTKEY, LOG_DROP, "", False),
    ("M-D-receipt-failed-set-claims-live", BALL, FAILED_BLOCK, FALSE_LIVE_BLOCK, False),
    ("M-E-drop-mirror-only", BALL, MIRROR_DROP_BLOCK, MIRROR_DROP_KEEP, False),
]

# ---------------------------------------------------------------- helpers


def sha256(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 16), b""):
            h.update(chunk)
    return h.hexdigest()


def md5(path):
    h = hashlib.md5()
    with open(path, "rb") as f:
        h.update(f.read())
    return h.hexdigest()


def read(path):
    with open(path, "rb") as f:
        return f.read()


def write(path, data):
    with open(path, "wb") as f:
        f.write(data)


def stamp():
    return time.strftime("%Y-%m-%d %H:%M:%S")


def now_ms():
    return time.strftime("%H:%M:%S.") + f"{int(time.time() * 1000) % 1000:03d}"


def run_tests(full, log_path, header):
    args = ["go", "test", "-count=1", "-v"]
    if not full:
        args += ["-run", "260r5|BorrowRefused"]
    args += ["./internal/ball/"]
    with open(log_path, "w", encoding="utf-8", newline="\n") as lf:
        lf.write(header)
        lf.write(f"command: {' '.join(args)}\n")
        lf.write(f"GOTMPDIR: {ENV['GOTMPDIR']}\n")
        lf.flush()
        p = subprocess.run(args, cwd=ROOT, env=ENV, stdout=lf, stderr=subprocess.STDOUT)
    return p.returncode


def main():
    ledger_path = os.path.join(LOGS, "mutation-sha256-ledger.txt")
    ledger = open(ledger_path, "w", encoding="utf-8", newline="\n")

    def book(line):
        ledger.write(line + "\n")
        ledger.flush()
        print(line, flush=True)

    book("### 260-v2 mutation ledger (sha256; md5 only as r5 cross-check)")
    book(f"script started {stamp()}")

    # start interlock: three copies byte-identical per file
    for label, work, blob, bak in (
        ("ball_windows.go", BALL, BALL_BLOB, BALL_BAK),
        ("hotkey_windows.go", HOTKEY, HOTKEY_BLOB, HOTKEY_BAK),
    ):
        s = (sha256(work), sha256(blob), sha256(bak))
        ok = s[0] == s[1] == s[2]
        book(f"start {label}: working={s[0]} blob={s[1]} backup={s[2]} identical={ok}")
        if not ok:
            book("REFUSE: start bytes differ (someone else's edit is in the file); no mutation run")
            ledger.close()
            sys.exit(2)

    orig_hashes = {BALL: sha256(BALL), HOTKEY: sha256(HOTKEY)}
    orig_data = {BALL: read(BALL), HOTKEY: read(HOTKEY)}

    for name, path, old, new, full in MUTATIONS:
        src = orig_data[path].decode("utf-8")
        n = src.count(old)
        if n != 1:
            book(f"{name}: ABORT anchor found {n} times, want 1")
            ledger.close()
            sys.exit(3)
        log_path = os.path.join(LOGS, f"mutation-{name}-gotest-v.txt")
        t0 = stamp()
        write(path, src.replace(old, new).encode("utf-8"))
        mmd5, msha = md5(path), sha256(path)
        rc = run_tests(full, log_path,
                       f"### {name} start={t0} mutated-sha256={msha} full={full}\n")
        t1 = stamp()
        write(path, orig_data[path])  # restore from the pre-mutation bytes
        rsha = sha256(path)
        book(f"{name}: window {t0} -> {t1} rc={rc} "
             f"mutated-sha256={msha} mutated-md5={mmd5} restored-sha256={rsha} restored-ok={rsha == orig_hashes[path]}")
        if rsha != orig_hashes[path]:
            book("ABORT: restore did not return the start bytes")
            ledger.close()
            sys.exit(4)
        # per-mutation verdict counts, straight off the log
        text = open(log_path, encoding="utf-8").read()
        p = text.count("\n--- PASS:")
        f_ = text.count("\n--- FAIL:")
        s_ = text.count("\n--- SKIP:")
        reds = [ln.split("--- FAIL: ")[1] for ln in text.splitlines() if ln.startswith("--- FAIL: ")]
        book(f"    verdicts PASS={p} FAIL={f_} SKIP={s_} reds={reds}")

    # final interlock
    for label, path in (("ball_windows.go", BALL), ("hotkey_windows.go", HOTKEY)):
        s = sha256(path)
        book(f"FINAL {label}: sha256={s} matches_start={s == orig_hashes[path]}")
    book(f"script ended {stamp()}")
    ledger.close()


if __name__ == "__main__":
    main()
