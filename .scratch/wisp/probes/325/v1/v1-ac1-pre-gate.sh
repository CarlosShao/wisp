#!/bin/sh
# 325-v1 AC#1 判据补全：票面要"改前/改后各跑一发 sh scripts/d22scan.sh"。
# 改后那发＝AC#4 里 logs/v1-ac4-gates.txt.d22scan-sh.txt（rc=0 clean）。
# 本尺补改前那发：在仓外导出树上把那一行物理退回混形（带门），跑同一枚门，再还原复尺。
# 用法：sh v1-ac1-pre-gate.sh <repo> <mixed-tree> <out>
set -u
repo=${1:?repo}; mixed=${2:?mixed tree}; out=${3:?out}
M='tools/d22scan/runtests.sh'
native_hits() { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}\\d22scan-' || true; }
mixed_hits()  { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}/d22scan-' || true; }

{
echo "V1 AC#1 pre-fix gate shot (d22scan.sh on the MIXED code)"
echo "start=$(date '+%F %T %z')"
echo "pristine sha=$(sha256sum "$mixed/$M" | cut -d' ' -f1)  native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed")"
mline=$(grep -n 'gotmp_run=' "$mixed/$M" | head -1 | cut -d: -f1)
sed -i "${mline}"'s/}\\\\d22scan-/}\/d22scan-/' "$mixed/$M"
echo "after mutate: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed") sha=$(sha256sum "$mixed/$M" | cut -d' ' -f1)"
if [ "$(mixed_hits "$mixed")" -ne 1 ] || [ "$(native_hits "$mixed")" -ne 0 ]; then
    echo 'MUTATION DID NOT LAND - aborting'; exit 3
fi
echo '--- sh scripts/d22scan.sh on the mixed code:'
( cd "$mixed" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT PATH="$mixed/third_party/sherpa-onnx:$PATH" \
      sh scripts/d22scan.sh ) > "$out.gate-pre.txt" 2>&1
echo "  PRE_RC=$?  bytes=$(wc -c < "$out.gate-pre.txt")"
grep -E 'GOTMPDIR=|^runtests.sh: OK|clean - no D22' "$out.gate-pre.txt" | sed 's/^/  | /'
git -C "$repo" show 428ef21e:$M > "$mixed/$M"
echo "after restore: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed") sha=$(sha256sum "$mixed/$M" | cut -d' ' -f1)"
echo '--- same gate on the restored (native) code, for the two-shots-on-one-screen read:'
( cd "$mixed" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT PATH="$mixed/third_party/sherpa-onnx:$PATH" \
      sh scripts/d22scan.sh ) > "$out.gate-post.txt" 2>&1
echo "  POST_RC=$?  bytes=$(wc -c < "$out.gate-post.txt")"
grep -E 'GOTMPDIR=|^runtests.sh: OK|clean - no D22' "$out.gate-post.txt" | sed 's/^/  | /'
echo "end=$(date '+%F %T %z')"
} >> "$out" 2>&1
echo "rc=$? out=$out"
