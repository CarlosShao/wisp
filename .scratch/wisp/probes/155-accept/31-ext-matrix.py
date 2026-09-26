#!/usr/bin/env python
"""Ticket 155 ACCEPT pass, cell 3 driver: re-run ext-A/B/C/D on my own snapshot.

Snapshot  = D:/tmp/wisp155accept/snap6de  (blob-exact tree of 6de3d1c5; the
            three files under test are byte-identical to the delivered 153
            commit b23c7f7, see probe 11).
Mutation  = the pre-existing driver .scratch/wisp/probes/153/mut153.py, read
            from anchor ac7fb00 and NOT modified: it refuses to land unless its
            anchor string occurs exactly once, and prints LANDED otherwise.
Probes in = zzaccept155_parsed-ondisk  (mine, JSON key/value ruler)
            zzaccept153_ondisk_test.go (ticket 153's two-leg probe, from anchor
            777d6cc - used here to test 155's claim that it is structurally
            blind to `unloop`)
Nothing is written into the repo working tree; -overlay and -cover* are not used.
"""
import os
import re
import shutil
import subprocess
import sys

SNAP = "D:/tmp/wisp155accept/snap6de"
AGENT = os.path.join(SNAP, "internal", "agent")
PRISTINE = "D:/tmp/wisp155accept/pristine"
MUT = "D:/tmp/wisp155accept/mut153.py"
MYPROBE = os.path.join(os.path.dirname(os.path.abspath(__file__)),
                       "30-zzaccept155-parsed-ondisk_test.go")
T153PROBE = "D:/tmp/wisp155accept/zzaccept153_ondisk_test.go"
ENV = dict(os.environ, GOPROXY="off", GOFLAGS="-mod=mod")
FAILRE = re.compile(r"^\s*--- FAIL: (\S+)", re.M)
RUNRE = re.compile(r"^=== RUN\s+(\S+)", re.M)


def sh(*args):
    return subprocess.run(args, cwd=SNAP, stdout=subprocess.PIPE,
                          stderr=subprocess.STDOUT, env=ENV).stdout.decode("utf-8", "replace")


def restore():
    print("$ python mut153.py restore %s %s" % (SNAP, PRISTINE))
    print("   " + sh(sys.executable, MUT, "restore", SNAP, PRISTINE).strip().replace("\n", "\n   "))


def mutate(op):
    print("$ python mut153.py %s %s %s" % (op, SNAP, PRISTINE))
    out = sh(sys.executable, MUT, op, SNAP, PRISTINE).strip()
    print("   " + out.replace("\n", "\n   "))
    assert "LANDED" in out or op.startswith("restore"), "op did not land: " + op


def test(label, runfilter):
    out = sh("go", "test", "./internal/agent/", "-count=1", "-v", "-run", runfilter)
    runs = RUNRE.findall(out)
    fails = FAILRE.findall(out)
    print("--- %s : RUN=%d FAIL=%d %s" % (label, len(runs), len(fails),
                                          ["--- FAIL: " + f for f in fails] if fails else ""))
    for line in out.splitlines():
        if ("PROBE-155-ACCEPT" in line or "PROBE-153 DISK" in line or FAILRE.match(line)
                or line.startswith("ok  ") or line.startswith("FAIL") or "tagged record" in line
                or "want 2" in line):
            print("    " + line.strip())
    return runs, fails


print("### 155 对抗验收程 · ext 矩阵（我自己的快照 + 我自己的尺）")
print("go version =", sh("go", "version").strip())
for src in (MYPROBE, T153PROBE):
    shutil.copyfile(src, os.path.join(AGENT, os.path.basename(src)))
    print("probe copied into snapshot:", os.path.basename(src))

print("\n================ 正控①：未变异，这棵树活着 ================")
restore()
test("delivered TestCompressionTrace*", "TestCompressionTrace")

print("\n================ 正控②：已知会红（drop-attr） ================")
mutate("drop-attr")
test("delivered TestCompressionTrace*", "TestCompressionTrace")

print("\n================ ext-B（摘 N3＝drop-attr），本程的尺 ================")
test("ext-B parsed on-disk", "TestAccept155ParsedRecordsOnDisk")

print("\n================ ext-A（交付码基线），本程的尺 ================")
restore()
test("ext-A parsed on-disk", "TestAccept155ParsedRecordsOnDisk")

print("\n================ ext-C（摘 N2＝unloop），本程的尺 + 153 两腿探针 ================")
mutate("unloop")
test("ext-C parsed on-disk", "TestAccept155ParsedRecordsOnDisk")
test("ext-C 153 两腿到盘探针（宣称：结构上看不见这一发）", "TestAccept153TraceOnDiskThroughRedactor")

print("\n================ ext-B 复跑第二回（两回互校）+ 153 两腿探针应红 ================")
restore()
mutate("drop-attr")
test("ext-B parsed on-disk 第二回", "TestAccept155ParsedRecordsOnDisk")
test("ext-B 153 两腿到盘探针（应红）", "TestAccept153TraceOnDiskThroughRedactor")

print("\n================ ext-D（摘 N1＝AC#1 那枚判据用例），本程的尺 ================")
restore()
mutate("drop-test:TestCompressionTraceSilentWhenNothingFoldableOverThreshold")
test("ext-D parsed on-disk", "TestAccept155ParsedRecordsOnDisk")
print("\n---- ext-D 同一棵树上再施 M5（句①那一形：摘证人⇒逃逸回来） ----")
mutate("m5")
test("ext-D + m5 delivered suite", "TestCompressionTrace")

print("\n================ 收口：还原自证 ================")
restore()
out = sh("git", "hash-object", "internal/agent/compress.go", "internal/agent/loop.go",
         "internal/agent/compress_trace_test.go")
print("$ git hash-object 三枚被测文件（应等于 6de3d1c5/b23c7f7 的 blob 号）")
for line in out.split():
    print("   " + line)
for f in ("compress.go", "loop.go", "compress_trace_test.go"):
    a = open(os.path.join(AGENT, f), "rb").read()
    b = open(os.path.join(PRISTINE, "internal", "agent", f), "rb").read()
    print("   RESTORED-IDENTICAL %-26s %s" % (f, a == b))
print("   快照里本程新增的 _test.go 枚数：",
      [n for n in os.listdir(AGENT) if n.endswith("_test.go") and ("155" in n or "accept153" in n)])
