#!/bin/sh
# AC#2 结构尺（独立再跑）：blacklist.go 起手 blob c1799cd2 ⇄ 交件 blob 62f2dddb，
# 剥掉注释（sed 's://.*::'）＋去空白＋去空行后，两侧可执行字节 sha256 必须相同。
# ⛔ 采信任何人给的哈希串 —— 这里两串都是我自己现算。
set -e
cd '/d/work/workspace/projects plans/Wisp' || exit 90
LOG=/d/tmp/wisp/20261019/343v1/logs
B_A=c1799cd24afce19e6accef23db23a7d7ee746dec
B_B=62f2dddbde7c72eb12cbeb4cf1ec1edc2708254b
D=/d/tmp/wisp/20261019/343v1/ac2
mkdir -p $D

# 1) raw blob dump from object store (not from working tree)
git cat-file blob $B_A > $D/raw-start.go
git cat-file blob $B_B > $D/raw-head.go

{
echo "== 0. full blob hashes (start / head) =="
git rev-parse $B_A $B_B
echo "== 1. raw sha256 (different expected: comment line changed) =="
sha256sum $D/raw-start.go $D/raw-head.go
echo "== 2. raw line counts / byte counts / CR count =="
wc -l -c $D/raw-start.go $D/raw-head.go
printf 'CR_start='; tr -cd '\r' < $D/raw-start.go | wc -c
printf 'CR_head='; tr -cd '\r' < $D/raw-head.go | wc -c
} > $LOG/g02-ac2-01-raw.log 2>&1

strip() { # 剥注释 -> 去水平空白 -> 去空行
  sed 's://.*::' | tr -d ' \t\r' | sed '/^$/d'
}

strip < $D/raw-start.go > $D/exe-start.txt
strip < $D/raw-head.go  > $D/exe-head.txt

{
echo "== 3. stripped (executable-only) sha256 BOTH SIDES =="
sha256sum $D/exe-start.txt $D/exe-head.txt
echo "== 4. stripped line counts (orchestrator claimed 232 each) =="
wc -l $D/exe-start.txt $D/exe-head.txt
echo "== 5. cmp (0 = identical) =="
cmp $D/exe-start.txt $D/exe-head.txt && echo CMP=IDENTICAL || echo CMP=DIFFERENT
echo "== 6. byte counts of the two stripped files =="
wc -c $D/exe-start.txt $D/exe-head.txt
echo "== 7. symdiff of stripped sides (must be empty) =="
diff $D/exe-start.txt $D/exe-head.txt && echo DIFF=EMPTY
echo "== 8. second, independently-shaped stripped gauge: delete whole //-lines then strip inline =="
echo "-- start --"
sed -e '/^[[:space:]]*\/\//d' -e 's://.*::' $D/raw-start.go | tr -d ' \t\r' | sed '/^$/d' | sha256sum
echo "-- head --"
sed -e '/^[[:space:]]*\/\//d' -e 's://.*::' $D/raw-head.go | tr -d ' \t\r' | sed '/^$/d' | sha256sum
echo "== 9. what actually differs in the raw blobs (proof: only one comment line) =="
diff $D/raw-start.go $D/raw-head.go || true
} > $LOG/g02-ac2-02-stripped.log 2>&1

{
echo "== 10. git grep -c 'func aEquiv' HEAD -- internal/risk (expect 0 hits) =="
git grep -c 'func aEquiv' HEAD -- internal/risk || true
echo "== 11. git grep -n aEquiv HEAD -- internal/risk (delivery side, expect 0 hits) =="
git grep -n aEquiv HEAD -- internal/risk || true
echo "-- positive control for the same grep shape: aEq =="
git grep -c 'aEq' HEAD -- internal/risk/blacklist.go || true
echo "== 12. start-side control: aEquiv DID exist at 29d258f9 =="
git grep -n aEquiv 29d258f9 -- internal/risk || true
} > $LOG/g02-ac2-03-grep.log 2>&1
echo "DONE"
