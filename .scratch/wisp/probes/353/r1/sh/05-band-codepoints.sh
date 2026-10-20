#!/bin/sh
# 353-r1 起手自查（第二把尺）：新件里是否含 d22scan 实际扫的高位码位（★⛔ 拿 AGENTS.md 那句"注释豁免"当凭据＝
# tools/d22scan/main.go:1068 逐字写着 emojiRe 打 RAW 行、⛔ 剥注释）。
# 尺＝python 逐枚码位＋射程＝该件 ⇄ internal/winsec/*.go 全体；决定性那把＝真跑 sh scripts/d22scan.sh（见 06-gates.sh）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
python - <<'PY'
import collections, glob
BANDS = [(0x1F000,0x1FAFF),(0x2200,0x22FF),(0x2600,0x27BF),(0x2B00,0x2BFF),(0xFE00,0xFE0F),(0x1F1E6,0x1F1FF)]
def scan(p):
    t = open(p, encoding='utf-8').read()
    c = collections.Counter(ch for ch in t if ord(ch) > 127)
    hits = collections.Counter()
    for ch, n in c.items():
        for lo, hi in BANDS:
            if lo <= ord(ch) <= hi:
                hits['U+%04X' % ord(ch)] += n
    return c, hits
mine = 'internal/winsec/nametopup_353_windows_test.go'
c, h = scan(mine)
print('MYFILE distinct nonascii =', sorted('U+%04X' % ord(k) for k in c))
print('MYFILE scanned-band hits =', dict(h) or 'none')
print('--- 既有 internal/winsec/*.go 的 scanned-band 命中（先例尺）---')
for f in sorted(glob.glob('internal/winsec/*.go')):
    cc, hh = scan(f)
    if hh:
        print('%s -> %s' % (f, dict(hh)))
print('--- 全仓 tracked .go 的 scanned-band 命中（前 12 枚）---')
import subprocess
files = subprocess.run(['git','ls-files','-z','*.go'], capture_output=True).stdout.split(b'\0')
n = 0
for b in files:
    if not b: continue
    p = b.decode('utf-8')
    try:
        cc, hh = scan(p)
    except Exception:
        continue
    if hh:
        print('%s -> %s' % (p, dict(hh)))
        n += 1
        if n >= 12: break
print('TRACKED_GO_WITH_BAND_HITS_SHOWN =', n)
PY
