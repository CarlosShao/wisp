#!/usr/bin/env python
"""Ticket 155: build a BYTE-EXACT out-of-repo snapshot of a commit's tree.

Why not `git archive | tar -x`: this repo has `.gitattributes` = `* text=auto`
and local `core.autocrlf=true`, with no `eol=lf` rule on txt/log/tsx/json, so an
archive tree is NOT byte-equal to the worktree (a real .log blob measured 3391
vs 3436 bytes). So: `git ls-tree -r <anchor>` for the blob ids, then one
`git cat-file --batch` process for the bytes, then a local sha1 recomputation
(`blob <len>\0` + data) compared against the ls-tree id for EVERY file.
Any mismatch aborts with a non-zero exit code.

NB: stdin of `git cat-file --batch` is a FILE, not a pipe. The first version
piped the sha list in and deadlocked (child blocked writing its stdout pipe
while the parent was still blocked writing the stdin pipe); it ran 7 min and
landed 0 files. Registered in the evidence file as an instrument finding.

usage: sync_blob_tree.py <repo> <anchor> <dest>
"""
import hashlib
import os
import subprocess
import sys

repo, anchor, dest = sys.argv[1], sys.argv[2], sys.argv[3]
os.makedirs(dest, exist_ok=True)
scratch = os.path.dirname(os.path.abspath(dest))
listfile = os.path.join(scratch, "sha-list-%s.txt" % anchor)

ls = subprocess.run(["git", "-C", repo, "ls-tree", "-r",
                     "--format=%(objecttype) %(objectname) %(path)", anchor],
                    capture_output=True, text=True, encoding="utf-8", check=True).stdout
entries = []
for line in ls.splitlines():
    kind, sha, path = line.split(" ", 2)
    if kind != "blob":
        raise SystemExit("unexpected entry: " + line)
    entries.append((sha, path))
print("ls-tree blobs: %d" % len(entries))

with open(listfile, "w") as lf:
    lf.write("\n".join(s for s, _ in entries) + "\n")

bad = []
written = 0
with open(listfile, "rb") as lin:
    proc = subprocess.Popen(["git", "-C", repo, "cat-file", "--batch"],
                            stdin=lin, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    for sha, path in entries:
        header = proc.stdout.readline()
        if not header:
            proc.kill()
            raise SystemExit("batch truncated at %s (written=%d)" % (path, written))
        parts = header.decode().strip().split(" ")
        if len(parts) != 3 or parts[1] != "blob":
            proc.kill()
            raise SystemExit("bad header for %s: %r" % (path, header))
        size = int(parts[2])
        data = proc.stdout.read(size)
        nl = proc.stdout.read(1)
        assert nl == b"\n", "missing delimiter after %s" % path
        got = hashlib.sha1(b"blob %d\x00" % size + data).hexdigest()
        if got != sha or size != len(data):
            bad.append("%s: want %s got %s (%d/%d bytes)" % (path, sha, got, len(data), size))
            continue
        dst = os.path.join(dest, path.replace("/", os.sep))
        os.makedirs(os.path.dirname(dst), exist_ok=True)
        with open(dst, "wb") as f:
            f.write(data)
        written += 1
    proc.stdout.read()
    proc.wait()

print("written: %d   mismatched: %d" % (written, len(bad)))
for b in bad[:20]:
    print("  MISMATCH " + b)
print("SNAPSHOT %s -> %s : %s" % (anchor, dest, "EXACT" if not bad else "NOT EXACT"))
sys.exit(1 if bad else 0)
