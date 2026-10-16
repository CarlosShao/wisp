#!/bin/sh
# 325-v1 AC#2 成对反控（三形同屏）
#   A 正形  ＝ HEAD 码原样（}\d22scan-）              → 预期绿
#   B 混形  ＝ 导出树上把那一行物理退回 }/d22scan-    → 预期四红（本票的牙）
#   C 纯正  ＝ GOTMPDIR 整根全正斜杠 C:/…            → 预期 mklink Invalid switch（票 317 那对形）
#   还原    ＝ git show 428ef21e 的字节灌回去 ＋ sha 复尺 ＋ 复跑绿
# ⛔ -overlay：这族尺运行时读盘，overlay 恒绿＝假读数。
# 用法：sh v1-ac2-pair.sh <repo> <native-tree> <mixed-tree> <out>
set -u
repo=${1:?repo}; native=${2:?native tree}; mixed=${3:?mixed tree}; out=${4:?out}
NAMES='^(TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed)$'
M='tools/d22scan/runtests.sh'

sha() { sha256sum "$1/$M" | cut -d' ' -f1; }
# 形状尺：只数 gotmp_run 那一行（固定串，⛔ 正则形）
shape() {
    line=$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)
    txt=$(sed -n "${line}p" "$1/$M")
    printf 'line=%s native_tail_hits=%s mixed_tail_hits=%s' "$line" \
        "$(printf '%s' "$txt" | grep -cF '}\\d22scan-' || true)" \
        "$(printf '%s' "$txt" | grep -cF '}/d22scan-' || true)"
}
runfour() { # $1=tree $2=label  $3=extra env assignment (may be empty)
    echo "@@@@ RUN $2  tree=$1  start=$(date '+%F %T %z')"
    if [ -n "${3:-}" ]; then
        ( cd "$1" && env -u WISP_TEST_GOTMP_ROOT PATH="$1/third_party/sherpa-onnx:$PATH" \
              GOTMPDIR="$3" sh tools/d22scan/runtests.sh ./internal/risk/ -run "$NAMES" )
        echo "RC=$?"
    else
        ( cd "$1" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT \
              PATH="$1/third_party/sherpa-onnx:$PATH" \
              sh tools/d22scan/runtests.sh ./internal/risk/ -run "$NAMES" )
        echo "RC=$?"
    fi
    echo "@@@@ END $2"
    echo
}

{
echo "V1 AC#2 three-shape counter-control"
echo "start=$(date '+%F %T %z')"
echo "repo=$repo  native=$native  mixed=$mixed"
echo
echo '===== 0. PRISTINE STATE'
echo "native sha  =$(sha "$native")"
echo "mixed  sha  =$(sha "$mixed")  $(shape "$mixed")"
echo "git 428ef21e blob sha =$(git -C "$repo" show 428ef21e:$M | sha256sum | cut -d' ' -f1)"
echo
echo '===== 1. NAME SELF-PROOF (-list must hit 4/4, else the run below is a false reading)'
for n in TestSyncAncestorActableLegFailsClosedOnMovedAncestor116 \
         TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116 \
         TestSyncAncestorLegReResolvesThroughJunctionToRealTree \
         TestSyncAncestorLegNotExemptedStaysFailClosed; do
    hits=$( cd "$native" && GOFLAGS= go test ./internal/risk/ -run "^${n}\$" -list "^${n}\$" 2>/dev/null | grep -c "^Test" || true )
    printf 'list %-58s hits=%s\n' "$n" "$hits"
done
echo
echo '===== 2. SHAPE A: native code as committed (expect GREEN)'
runfour "$native" A-native
echo
echo '===== 3. MUTATE mixed tree: physical byte revert of the tail separator'
mline=$(grep -n 'gotmp_run=' "$mixed/$M" | head -1 | cut -d: -f1)
sed -i "${mline}s/}\\\\d22scan-/}\/d22scan-/" "$mixed/$M"
echo "mutation line=$mline"
echo "after mutate: $(shape "$mixed")   sha=$(sha "$mixed")"
echo "MUTATION_SELF_PROOF: mixed-shape-hits on the target line must be 1 and native 0 (else sed no-oped)"
echo
echo '===== 4. SHAPE B: mixed tail, same four names (expect FOUR RED)'
runfour "$mixed" B-mixed
echo
echo '===== 5. SHAPE C: pure forward-slash root on native code (ticket 317 shape)'
psxroot="C:/tmp/wisp325v1/psxroot-$(date +%H%M%S)-$$"
mkdir -p "$psxroot"
echo "psx root=$psxroot"
runfour "$native" C-pure-slash "$psxroot"
echo
echo '===== 6. RESTORE mixed tree from 428ef21e bytes, re-measure, re-run green'
git -C "$repo" show 428ef21e:$M > "$mixed/$M"
echo "restored: $(shape "$mixed")   sha=$(sha "$mixed")"
echo "restore_sha_equals_committed=$([ "$(sha "$mixed")" = "$(git -C "$repo" show 428ef21e:$M | sha256sum | cut -d' ' -f1)" ] && echo YES || echo NO)"
runfour "$mixed" D-restored-green
echo "end=$(date '+%F %T %z')"
} > "$out" 2>&1
echo "wrote $out rc=$?"
tail -5 "$out"
