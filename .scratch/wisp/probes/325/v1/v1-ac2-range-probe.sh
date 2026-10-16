#!/bin/sh
# 325-v1 AC#2 射程那一格：d22scan 这档今天到底够不够得着 internal/risk 那四枚
# 静态三把尺 ＋ 动态一发（CI 直呼本 runner 的那步形状，在未修混形码上跑）
# 用法：sh v1-ac2-range-probe.sh <repo> <mixed-tree> <out>
set -u
repo=${1:?repo}; mixed=${2:?mixed tree}; out=${3:?out}
M='tools/d22scan/runtests.sh'

native_hits() { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}\\d22scan-' || true; }
mixed_hits()  { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}/d22scan-' || true; }

{
echo "V1 AC#2 range probe (is ./internal/risk/ reachable from the d22scan gate today?)"
echo "start=$(date '+%F %T %z')   repo=$repo"
echo
echo '===== 1. STATIC: what does the only gate on this entry hand to the runner?'
grep -n 'runtests.sh' scripts/d22scan.sh
echo
echo '===== 2. STATIC: any mklink inside tools/ (the entry scope)?'
echo "  mklink_hits_in_tools=$(grep -rnF 'mklink' tools | wc -l)"
echo '  (fixed-string ruler; control: same ruler over internal/ gives)'
echo "  mklink_files_in_internal=$(grep -rlF 'mklink' internal | wc -l)   files, listed:"
grep -rlF 'mklink' internal | sed 's/^/    /'
echo
echo '===== 3. STATIC: the CI step that calls this runner against internal/risk'
grep -n 'runtests.sh' .github/workflows/ci.yml
echo '  --- the line(s) around it:'
n=$(grep -n 'runtests.sh' .github/workflows/ci.yml | head -1 | cut -d: -f1)
sed -n "$((n - 6)),$((n + 6))p" .github/workflows/ci.yml
echo
echo '===== 4. DYNAMIC: run THAT step shape on the MIXED code (mutate first, gate, restore after)'
mline=$(grep -n 'gotmp_run=' "$mixed/$M" | head -1 | cut -d: -f1)
sed -i "${mline}"'s/}\\\\d22scan-/}\/d22scan-/' "$mixed/$M"
echo "  after mutate: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed")"
if [ "$(mixed_hits "$mixed")" -ne 1 ] || [ "$(native_hits "$mixed")" -ne 0 ]; then
    echo '  MUTATION DID NOT LAND - aborting'; exit 3
fi
echo "  emitted root line + verdict for the CI -run name on the mixed code:"
tmp="$out.ci-shape.raw"
( cd "$mixed" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT PATH="$mixed/third_party/sherpa-onnx:$PATH" \
      sh tools/d22scan/runtests.sh ./internal/risk/ -run '^TestPathResolverJunctionWindows$' ) > "$tmp" 2>&1
ci_rc=$?
grep -E 'GOTMPDIR=|^--- |^runtests.sh: |^ok |^FAIL' "$tmp" | sed 's/^/    /'
echo "  CI_SHAPE_RC=$ci_rc"
echo
echo '===== 5. and the SAME -run name on the four 116/105 cases is NOT in that pattern (why the gate stays green)'
( cd "$mixed" && GOFLAGS= go test ./internal/risk/ -list '^TestPathResolverJunctionWindows$' 2>/dev/null | sed 's/^/    /' )
git -C "$repo" show 428ef21e:$M > "$mixed/$M"
echo "  after restore: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed") sha=$(sha256sum "$mixed/$M" | cut -d' ' -f1)"
echo "end=$(date '+%F %T %z')"
} > "$out" 2>&1
echo "rc=$? out=$out"
