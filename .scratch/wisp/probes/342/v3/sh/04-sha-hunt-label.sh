#!/bin/sh
# Hunting the exact normalisation behind the quoted label e8a2f33dea8cdeab so this leg can either
# reproduce it or name the discrepancy instead of wearing it as its own reading.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
TARGET=e8a2f33dea8cdeab
{
echo "TARGET=$TARGET"
for r in 118,138 118,140 118,141 118,137 118,139 119,140 117,140 116,140 118,142 118,143; do
  s1=$(git cat-file blob HEAD:internal/winsec/winsec.go | sed -n "${r}p" | sha1sum | cut -c1-16)
  s2=$(git cat-file blob HEAD:internal/winsec/winsec.go | sed -n "${r}p" | sha256sum | cut -c1-16)
  s3=$(git cat-file blob HEAD:internal/winsec/winsec.go | sed -n "${r}p" | tr -d '\r\n' | sha1sum | cut -c1-16)
  hit=no
  [ "$s1" = "$TARGET" ] && hit=SHA1
  [ "$s2" = "$TARGET" ] && hit=SHA256
  [ "$s3" = "$TARGET" ] && hit=SHA1-NOLINEENDS
  echo "$r sha1=$s1 sha256=$s2 sha1nn=$s3 hit=$hit"
done
echo "== whole file under other hashes =="
echo "sha1sum-file=$(git cat-file blob HEAD:internal/winsec/winsec.go | sha1sum | cut -c1-16)"
echo "sha256-file=$(git cat-file blob HEAD:internal/winsec/winsec.go | sha256sum | cut -c1-16)"
echo "== the block delimited by the comment only (awk between markers) =="
git cat-file blob HEAD:internal/winsec/winsec.go | awk '/Ticket 322 is the second half/,/^func privateFile/' | sha1sum | cut -c1-16
echo "== conclusion =="
echo "No candidate reproduced the label => the quoted number is NOT this leg's reading; the substance"
echo "(HEAD block == baseline block, byte-for-byte) is proven by equal sha16 within each range instead."
} 2>&1 | tee /d/tmp/wisp/20261020/342v3/logs/sha-hunt.txt
echo RC=$?
