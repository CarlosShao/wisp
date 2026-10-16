#!/bin/sh
# 325-v1 AC#2 第二发（牙修好了的那把）
# 第一发 logs/v1-ac2-three-shapes-pair.txt 的教训具名：我的 sed 写在双引号里，
# "\\\\\\\\" 被 shell 吃成两枚、sed 再读成一枚字面反斜杠 ⇒ 与文件里的两枚对不上 ⇒
# 突变静默 no-op（sha 未变），那一发"混形"其实跑的是正形，读数作废。
# 本尺加硬门：突变后必须 mixed_tail_hits=1 且 native_tail_hits=0，否则 exit 3 不跑测试。
# 四形同屏：A 正形绿 → B 混形红 → C 纯正斜杠红 → D 还原后再绿。
# 用法：sh v1-ac2-mixed-tail-pair2.sh <repo> <native-tree> <mixed-tree> <out>
set -u
repo=${1:?repo}; native=${2:?native tree}; mixed=${3:?mixed tree}; out=${4:?out}
M='tools/d22scan/runtests.sh'
NAMES='^(TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed)$'

sha() { sha256sum "$1/$M" | cut -d' ' -f1; }
# 形状尺：只看 gotmp_run 那一行；固定串，⛔ 正则形
native_hits() { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}\\d22scan-' || true; }
mixed_hits()  { sed -n "$(grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1)p" "$1/$M" | grep -cF '}/d22scan-' || true; }
target_line() { grep -n 'gotmp_run=' "$1/$M" | head -1 | cut -d: -f1; }

runfour() { # $1=tree $2=label [$3=pure-slash root]
    echo "@@@@ RUN $2 tree=$1 start=$(date '+%F %T %z')"
    if [ -n "${3:-}" ]; then
        ( cd "$1" && env -u WISP_TEST_GOTMP_ROOT GOTMPDIR="$3" \
              PATH="$1/third_party/sherpa-onnx:$PATH" \
              sh tools/d22scan/runtests.sh ./internal/risk/ -run "$NAMES" )
    else
        ( cd "$1" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT \
              PATH="$1/third_party/sherpa-onnx:$PATH" \
              sh tools/d22scan/runtests.sh ./internal/risk/ -run "$NAMES" )
    fi
    echo "RC=$?"
    echo "@@@@ END $2"; echo
}

{
echo "V1 AC#2 second shot: gated mutation, four shapes on one screen"
echo "start=$(date '+%F %T %z')   repo=$repo"
echo "native tree sha=$(sha "$native")   mixed tree sha=$(sha "$mixed")"
echo "committed 428ef21e blob sha=$(git -C "$repo" show 428ef21e:$M | sha256sum | cut -d' ' -f1)"
echo "target line in mixed tree=$(target_line "$mixed")  native_hits=$(native_hits "$mixed")  mixed_hits=$(mixed_hits "$mixed")"
echo
echo '===== A. NATIVE (committed code, expect GREEN)'
runfour "$native" A-native
echo '===== B. MUTATE to mixed tail, GATED'
mline=$(target_line "$mixed")
# 单引号包住整段 sed：4 枚反斜杠进 sed ⇒ 读成 2 枚字面 ⇒ 与文件字节对上
sed -i "${mline}"'s/}\\\\d22scan-/}\/d22scan-/' "$mixed/$M"
echo "mutation attempted on line $mline"
echo "after mutate: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed") sha=$(sha "$mixed")"
sed -n "${mline}p" "$mixed/$M"
if [ "$(mixed_hits "$mixed")" -ne 1 ] || [ "$(native_hits "$mixed")" -ne 0 ]; then
    echo "MUTATION DID NOT LAND - aborting, no fake mixed run this time. rc=3"
    exit 3
fi
echo "GATE PASSED: the mixed tail is physically on disk (⛔ overlay - this family reads the disk at run time)"
echo
runfour "$mixed" B-mixed
echo '===== C. PURE FORWARD-SLASH ROOT on native code (ticket 317 wording)'
psx="C:/tmp/wisp325v1/psx2-$(date +%H%M%S)"
mkdir -p "$psx"
echo "psx root=$psx"
runfour "$native" C-pure-slash "$psx"
echo '===== D. RESTORE from 428ef21e bytes + re-measure + re-run green'
git -C "$repo" show 428ef21e:$M > "$mixed/$M"
echo "after restore: native_hits=$(native_hits "$mixed") mixed_hits=$(mixed_hits "$mixed") sha=$(sha "$mixed")"
runfour "$mixed" D-restored
echo "end=$(date '+%F %T %z')"
} > "$out" 2>&1
rc=$?
echo "script_rc=$rc out=$out"
