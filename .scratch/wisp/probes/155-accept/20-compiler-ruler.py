#!/usr/bin/env python
"""Ticket 155 ACCEPT pass, cell 1/cell 2 ruler: ask the compiler, not grep.

Three rulers the impl pass did NOT use for "can Compress reach the task id":
insert one statement that *references* the candidate channel and let
`go build ./internal/agent/` decide. A compile error naming the missing
thing is a measurement of absence; the same patch compiling on the post-change
tree is the positive control that the instrument can answer "present".

  compress.go, right after the Compress header:
    C-env  _ = ctx                     positive control (must build at both)
    C-b    _ = c.b                     positive control: receiver field exists
    C-f1   _ = taskID                  channel "parameter / any in-scope name"
    C-f2   _ = c.taskID                channel "receiver field"
    C-f3   _ = traceTaskID(ctx)        channel "ctx carries a readable task id"
  loop.go, right before the Compress call:
    L-f1   _ = taskID                  is the param in scope at the call?
    L-f2   _ = l.taskID                does the Loop object hold an id?
    L-f3   _ = root.ID                 second holder alive at that moment?
    L-f4   _ = res.TaskID              third holder (and is it BEFORE the call?)
    L-f5   _ = j.taskID                fourth holder (journal), unlisted by 155

usage: compiler_ruler.py <snapshot_root> <label>
"""
import os
import subprocess
import sys

SNAP, LABEL = sys.argv[1], sys.argv[2]
AGENT = os.path.join(SNAP, "internal", "agent")
FILES = {"compress": os.path.join(AGENT, "compress.go"),
         "loop": os.path.join(AGENT, "loop.go")}
PRISTINE = {k: open(v, "rb").read() for k, v in FILES.items()}

PROBES = [
    ("compress", "compress.go 头行之后", "C-env", "_ = ctx"),
    ("compress", "compress.go 头行之后", "C-b", "_ = c.b"),
    ("compress", "compress.go 头行之后", "C-f1", "_ = taskID"),
    ("compress", "compress.go 头行之后", "C-f2", "_ = c.taskID"),
    ("compress", "compress.go 头行之后", "C-f3", "_ = traceTaskID(ctx)"),
    ("loop", "loop.go 压缩那一发之前", "L-f1", "_ = taskID"),
    ("loop", "loop.go 压缩那一发之前", "L-f2", "_ = l.taskID"),
    ("loop", "loop.go 压缩那一发之前", "L-f3", "_ = root.ID"),
    ("loop", "loop.go 压缩那一发之前", "L-f4", "_ = res.TaskID"),
    ("loop", "loop.go 压缩那一发之前", "L-f5", "_ = j.taskID"),
]
ANCHOR = {
    "compress": "func (c *Compressor) Compress(",
    "loop": "nh, rep, err := l.comp.Compress(",
}


def restore():
    for k, v in FILES.items():
        with open(v, "wb") as fh:
            fh.write(PRISTINE[k])


def run(which, where, expr):
    restore()
    raw = PRISTINE[which].decode("utf-8")
    lines = raw.split("\n")
    hit = [i for i, l in enumerate(lines) if ANCHOR[which] in l]
    assert len(hit) == 1, "anchor not unique in %s: %r" % (FILES[which], hit)
    at = hit[0]
    if which == "compress":          # insert after the header line
        lines.insert(at + 1, "\t" + expr)
        where_line = at + 2
    else:                            # insert before the call line
        lines.insert(at, "\t\t" + expr)
        where_line = at + 1
    with open(FILES[which], "w", encoding="utf-8", newline="") as fh:
        fh.write("\n".join(lines))
    proc = subprocess.run(["go", "build", "./internal/agent/"], cwd=SNAP,
                          stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                          env=dict(os.environ, GOPROXY="off", GOFLAGS="-mod=mod"))
    out = proc.stdout.decode("utf-8", "replace").strip().replace("\r\n", "\n")
    return proc.returncode, where_line, out


print("### 编译器那把尺 · %s（快照 %s）" % (LABEL, SNAP))
env = subprocess.run(["go", "version"], cwd=SNAP, capture_output=True, text=True).stdout.strip()
print("go version =", env)
for which, where, name, expr in PROBES:
    rc, line, out = run(which, where, expr)
    verdict = "BUILD-OK（这一发在码里存在）" if rc == 0 else "BUILD-FAIL（编译器判它不存在）"
    print("\n--- %s  插在第 %d 行：%s %s" % (name, line, expr, verdict))
    print("    rc=%d" % rc)
    if out:
        for l in out.split("\n")[:6]:
            print("    " + l)
restore()
print("\n--- 还原自证 ---")
for k, v in FILES.items():
    same = open(v, "rb").read() == PRISTINE[k]
    print("    %s restored-identical=%s" % (k, same))
