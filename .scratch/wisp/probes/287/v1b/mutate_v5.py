#!/usr/bin/env python3
"""287-v1b V5：把深度判定那一支整体搬到池位占用之后（＝"次序换过来"那一发，
钉夹具 :197-199 那句"那一支必须在 TryAcquireSubagentSlot 之前返回"）。
按内容切块（⛔ 行号锚）：两块各自的起止行都逐字验过才动手，否则拒绝落种。"""
import sys, hashlib

path = sys.argv[1]
with open(path, "rb") as fh:
    raw = fh.read()
text = raw.decode("utf-8")
lines = text.split("\n")

D0 = "\tif rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {"
P0 = "\trelease, ok := t.d.Roster.TryAcquireSubagentSlot(MaxConcurrentSubagents)"
CLOSE = "\t}"

if lines.count(D0) != 1 or lines.count(P0) != 1:
    print("REFUSE-TO-SEED: block anchors not unique (D=%d P=%d)" % (lines.count(D0), lines.count(P0)))
    sys.exit(7)
di = lines.index(D0)
pi = lines.index(P0)
if di >= pi:
    print("REFUSE-TO-SEED: depth block is not before the pool block")
    sys.exit(7)
# 深度块＝D0 起，到其后第一枚 CLOSE 为止（含）
dj = di
while lines[dj] != CLOSE:
    dj += 1
# 池位块＝P0 起，到其后的 `if !ok {` 的收尾 CLOSE 为止
pk = pi
while True:
    pk += 1
    if lines[pk] == CLOSE:
        break
    if pk - pi > 40:
        print("REFUSE-TO-SEED: pool block runaway")
        sys.exit(7)

depth = lines[di:dj + 1]
pool = lines[pi:pk + 1]
out = lines[:di] + pool + depth + lines[pk + 1:]
blob = "\n".join(out).encode("utf-8")
with open(path, "wb") as fh:
    fh.write(blob)
print("MOVED depth_block_lines=%d pool_block_lines=%d hunks_shape=1-move" % (len(depth), len(pool)))
print("WROTE %s sha256=%s" % (path, hashlib.sha256(blob).hexdigest()[:16]))
print("cr_before=%d cr_after=%d" % (raw.count(b"\r"), blob.count(b"\r")))
