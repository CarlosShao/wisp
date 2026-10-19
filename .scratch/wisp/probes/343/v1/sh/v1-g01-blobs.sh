#!/bin/sh
# 起手尺：确认 blacklist.go 的两枚 blob（起手 c1799cd2 ⇄ 交件 62f2dddb）真在盘上，
# 并现量 343-r1 起手锚 29d258f9 与交件那笔 31d4a476 各自树里该文件 blob。
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
LOG=/d/tmp/wisp/20261019/343v1/logs
{
echo "== A. HEAD blacklist.go blob =="
git rev-parse HEAD:internal/risk/blacklist.go
echo "== B. start-anchor 29d258f9 blacklist.go blob =="
git rev-parse 29d258f9:internal/risk/blacklist.go
echo "== C. impl commit 31d4a476 blacklist.go blob =="
git rev-parse 31d4a476:internal/risk/blacklist.go
echo "== D. parents of 31d4a476 =="
git log -1 --format='%H %P %s' 31d4a476
echo "== E. cat-file -e both claimed blobs (short->full) =="
git rev-parse c1799cd2 62f2dddb
echo "== F. type =="
git cat-file -t c1799cd2
git cat-file -t 62f2dddb
echo "== G. working tree blacklist.go blob == HEAD blob ? =="
git hash-object internal/risk/blacklist.go
echo "== H. anchors present =="
git rev-parse --verify '29d258f9^{commit}'
git rev-parse --verify 'd5ba5352^{commit}'
echo "== I. issue 343 name/status unchanged =="
git ls-files -- .scratch/wisp/issues | grep -c 'issues/343-'
echo "RC_A=$?"
} > $LOG/g01-blobs.log 2>&1
echo "RC=$?"
