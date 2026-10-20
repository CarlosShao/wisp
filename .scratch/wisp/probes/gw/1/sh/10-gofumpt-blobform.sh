#!/bin/sh
# gw-1 味① 第 4 把的**blob 形**＝"CI 看到的是入库字节"那一问。
# 本地共享树的工作树字节带 CRLF（autocrlf）⇒ 35 行 TRACKED-DIRTY 里可能混着"只在本地脏"的枚。
# ⇒ 逐枚取 HEAD blob 落仓外同名件，再打同一把尺（gofumpt -l ⇄ gofmt -l 并排）。
# ★⛔ 洗任何件、⛔ -w、⛔ 动工作树；射程＝仓外 D:\tmp\wisp\20261020\gw1\blobcheck
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
B="$R/blobcheck"
unset GOTMPDIR TMPDIR TEMP TMP
echo "== START anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
mkdir -p "$B"
echo
echo "== 1) 喂料＝从 --tracked-only 那把尺的输出取 TRACKED-DIRTY 的 35 枚路径 =="
grep 'TRACKED-DIRTY' "$LOGS/attest-tracked.txt" | cut -f2 | tr '\\' '/' > "$LOGS/g4-35-paths.txt"
echo "PATHS=$(grep -c . "$LOGS/g4-35-paths.txt" || true)"
n=0
while IFS= read -r p; do
  [ -n "$p" ] || continue
  n=$((n+1))
  mkdir -p "$B/$(dirname "$p")"
  git cat-file blob "HEAD:$p" > "$B/$p"
done < "$LOGS/g4-35-paths.txt"
echo "FED=$n"
echo "BLOB_GO_FILES=$(find "$B" -name '*.go' | wc -l)"
echo
echo "== 2) 同一把尺打在 blob 字节上（并排两把）=="
/d/work/base/gopath/bin/gofumpt.exe -l "$B" > "$LOGS/g4-blobform-gofumpt.txt" 2>&1
echo "GOFUMPT_BLOB_RC=$?"
echo "GOFUMPT_BLOB_LISTED=$(grep -c . "$LOGS/g4-blobform-gofumpt.txt" || true)"
gofmt -l "$B" > "$LOGS/g4-blobform-gofmt.txt" 2>&1
echo "GOFMT_BLOB_RC=$?"
echo "GOFMT_BLOB_LISTED=$(grep -c . "$LOGS/g4-blobform-gofmt.txt" || true)"
echo
echo "== 3) 逐枚名册＝blob 形仍脏（＝CI 也红）⇄ 只在本地工作树脏（＝CRLF 假脏）=="
grep '\.go' "$LOGS/g4-blobform-gofumpt.txt" | tr '\\' '/' | sed "s|^$B/||" | sort > "$LOGS/g4-blob-dirty.names"
sort "$LOGS/g4-35-paths.txt" > "$LOGS/g4-35-paths.sorted"
echo "-- STILL_DIRTY_IN_BLOB=$(grep -c . "$LOGS/g4-blob-dirty.names" || true) --"
grep . "$LOGS/g4-blob-dirty.names" | sed 's|^|  BDIRTY |'
echo "-- LOCAL_ONLY_CRLF=$(comm -23 "$LOGS/g4-35-paths.sorted" "$LOGS/g4-blob-dirty.names" | grep -c . || true) --"
comm -23 "$LOGS/g4-35-paths.sorted" "$LOGS/g4-blob-dirty.names" | sed 's|^|  LOCALONLY |'
echo "-- 生产面（非 .scratch）在 blob 形里仍脏的枚数 --"
grep -v '^\.scratch/' "$LOGS/g4-blob-dirty.names" | grep -c . | sed 's|^|  PRODUCTION_BLOB_DIRTY=|' || true
grep -v '^\.scratch/' "$LOGS/g4-blob-dirty.names" | sed 's|^|  PRODIRTY |'
echo
echo "== 4) 第二把（gofmt）对 blob 形的名册差＝两把独立形对撞 =="
grep '\.go' "$LOGS/g4-blobform-gofmt.txt" | tr '\\' '/' | sed "s|^$B/||" | sort > "$LOGS/g4-blobform-gofmt.names"
echo "ONLY_GOFUMPT=$(comm -23 "$LOGS/g4-blob-dirty.names" "$LOGS/g4-blobform-gofmt.names" | grep -c . || true)"
comm -23 "$LOGS/g4-blob-dirty.names" "$LOGS/g4-blobform-gofmt.names" | sed 's|^|  GONLY |'
echo "ONLY_GOFMT=$(comm -13 "$LOGS/g4-blob-dirty.names" "$LOGS/g4-blobform-gofmt.names" | grep -c . || true)"
comm -13 "$LOGS/g4-blob-dirty.names" "$LOGS/g4-blobform-gofmt.names" | sed 's|^|  MONLY |'
echo
echo "== 5) 正控＝这 35 枚里取一枚已知"故意写坏"的样本，证明 blob 形尺读得出它 =="
head -1 "$LOGS/g4-blob-dirty.names" | sed 's|^|  POSCTRL_PICK=|'
echo "== END anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
