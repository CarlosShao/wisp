#!/usr/bin/env python3
"""287-v1b 突变台件（验收腿·只写仓外导出树）：逐字节替换一枚锚，命中数不是 1 就拒绝落种。
用法：python mutate.py <file> <old-anchor-file> <new-anchor-file>
锚从文件读（⛔ 走 shell 引号，中文与反引号一律不进命令行）。
打印 anchor_hits / 改行数 / hunks（尺＝diff -u 后 grep -c）。"""
import sys, hashlib

path, oldf, newf = sys.argv[1], sys.argv[2], sys.argv[3]
with open(path, "rb") as fh:
    data = fh.read()


def rd(p):
    with open(p, "rb") as fh:
        b = fh.read()
    # 锚文件由 heredoc 写成，末尾会带一枚换行；锚一律按"不含行尾"比对，
    # 这样 LF/CRLF 两形都不影响命中判定（本台件的尺另有 CR 计数，见下）。
    if b.endswith(b"\n"):
        b = b[:-1]
    return b


old = rd(oldf)
new = rd(newf)

hits = data.count(old)
print("anchor_hits=%d" % hits)
if hits != 1:
    print("REFUSE-TO-SEED: anchor must hit exactly once")
    sys.exit(7)
out = data.replace(old, new, 1)
with open(path, "wb") as fh:
    fh.write(out)
print("WROTE %s sha256=%s" % (path, hashlib.sha256(out).hexdigest()[:16]))
print("cr_before=%d cr_after=%d" % (data.count(b"\r"), out.count(b"\r")))
