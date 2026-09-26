#!/usr/bin/env python
"""Ticket 155 ACCEPT pass: build / verify a byte-exact out-of-repo snapshot of a commit tree.

Independent re-implementation (not a copy of the impl pass' sync_blob_tree.py):
  * `git ls-tree -r <anchor>` for blob ids            -> manifest in
  * `git cat-file --batch` with stdin from a FILE     -> bytes out
    (the pipe-fed form deadlocks; registered by ticket 153/155 as an instrument pit)
  * per-blob sha1 over "blob <len>\\0"+data recomputed and compared to ls-tree
  * nothing is trusted: a manifest is written and `verify` mode re-reads a
    directory against an existing manifest without touching git objects.

usage:
  mksnap.py make   <repo> <anchor> <dest> [<manifest>]
  mksnap.py verify <dest> <manifest>
  mksnap.py cmp    <dirA> <dirB>
"""
import hashlib
import json
import os
import subprocess
import sys


def git(repo, *args, stdin_data=None):
    return subprocess.run(["git", "-C", repo, *args], input=stdin_data,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=True).stdout


def blob_sha(size, data):
    return hashlib.sha1(b"blob %d\x00" % size + data).hexdigest()


def make(repo, anchor, dest, manifest):
    listing = git(repo, "ls-tree", "-r", "--format=%(objecttype) %(objectname) %(path)",
                  anchor).decode("utf-8")
    entries = []
    for line in listing.splitlines():
        kind, sha, path = line.split(" ", 2)
        assert kind == "blob", "non-blob entry: " + line
        entries.append((sha, path))
    print("anchor=%s blobs=%d" % (anchor, len(entries)))
    listfile = os.path.join(os.path.dirname(os.path.abspath(dest)), "shalist-%s.txt" % anchor)
    os.makedirs(os.path.dirname(listfile), exist_ok=True)
    with open(listfile, "w", encoding="utf-8", newline="\n") as fh:
        fh.write("".join(s + "\n" for s, _ in entries))
    rows, bad = [], []
    with open(listfile, "rb") as lin:
        proc = subprocess.Popen(["git", "-C", repo, "cat-file", "--batch"], stdin=lin,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        for sha, path in entries:
            header = proc.stdout.readline()
            if not header:
                raise SystemExit("batch stream ended early at %s" % path)
            _hsha, htype, hsize = header.decode().strip().split(" ")
            assert htype == "blob", path
            size = int(hsize)
            data = proc.stdout.read(size)
            assert proc.stdout.read(1) == b"\n", "delimiter missing at " + path
            if len(data) != size or blob_sha(size, data) != sha:
                bad.append(path)
                continue
            dst = os.path.join(dest, path.replace("/", os.sep))
            os.makedirs(os.path.dirname(dst) or dest, exist_ok=True)
            with open(dst, "wb") as out:
                out.write(data)
            rows.append({"path": path, "sha": sha, "size": size})
        proc.stdout.read()
        proc.wait()
    with open(manifest, "w", encoding="utf-8", newline="\n") as fh:
        json.dump({"anchor": anchor, "files": rows}, fh, indent=0, sort_keys=True)
    print("written=%d mismatched=%d manifest=%s" % (len(rows), len(bad), manifest))
    for b in bad[:10]:
        print("  MISMATCH " + b)
    return 1 if bad else 0


def verify(dest, manifest):
    rows = json.load(open(manifest, encoding="utf-8"))["files"]
    drift, gone = [], []
    for r in rows:
        p = os.path.join(dest, r["path"].replace("/", os.sep))
        if not os.path.exists(p):
            gone.append(r["path"])
            continue
        data = open(p, "rb").read()
        if blob_sha(len(data), data) != r["sha"]:
            drift.append(r["path"])
    extra = []
    for root, _dirs, files in os.walk(dest):
        for f in files:
            rel = os.path.relpath(os.path.join(root, f), dest).replace(os.sep, "/")
            if rel not in {r["path"] for r in rows}:
                extra.append(rel)
    print("manifest files=%d drift=%d missing=%d extra_in_tree=%d" %
          (len(rows), len(drift), len(gone), len(extra)))
    for x in drift[:20]:
        print("  DRIFT " + x)
    for x in gone[:20]:
        print("  MISSING " + x)
    for x in sorted(extra)[:40]:
        print("  EXTRA " + x)
    return 1 if (drift or gone) else 0


def cmp_dirs(a, b):
    fa = {os.path.relpath(os.path.join(r, f), a).replace(os.sep, "/")
          for r, _d, fs in os.walk(a) for f in fs}
    fb = {os.path.relpath(os.path.join(r, f), b).replace(os.sep, "/")
          for r, _d, fs in os.walk(b) for f in fs}
    diff = []
    for rel in sorted(fa & fb):
        if open(os.path.join(a, rel.replace("/", os.sep)), "rb").read() != \
           open(os.path.join(b, rel.replace("/", os.sep)), "rb").read():
            diff.append(rel)
    print("shared=%d byte-differ=%d onlyA=%d onlyB=%d" % (len(fa & fb), len(diff),
          len(fa - fb), len(fb - fa)))
    for x in diff[:20]:
        print("  DIFF " + x)
    for x in sorted(fa - fb)[:20]:
        print("  ONLY-A " + x)
    for x in sorted(fb - fa)[:20]:
        print("  ONLY-B " + x)
    return 1 if diff else 0


if __name__ == "__main__":
    cmd = sys.argv[1]
    if cmd == "make":
        sys.exit(make(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5]))
    if cmd == "verify":
        sys.exit(verify(sys.argv[2], sys.argv[3]))
    if cmd == "cmp":
        sys.exit(cmp_dirs(sys.argv[2], sys.argv[3]))
    raise SystemExit("unknown cmd")
