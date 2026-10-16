#!/bin/sh
# 325-v1 AC#3：整档 bash scripts/portable-tests.sh --scope=windows（仓外导出树，干净台面）
# 读数：全量 stdout 落盘（含每一枚 === RUN / --- PASS / --- FAIL / --- SKIP 与逐包结果行）
set -u
tree=${1:?need export tree root}
out=${2:?need out file}
cd "$tree" || exit 9
{
    echo "V1 AC#3 whole-windows-scope run"
    echo "start=$(date '+%F %T %z')"
    echo "tree=$tree"
    echo "go=$(go version)"
    echo "GOOS=$(go env GOOS)"
    echo "LOCALAPPDATA=$LOCALAPPDATA"
    echo "sherpa PATH prefix=$tree/third_party/sherpa-onnx"
    ls "$tree/third_party/sherpa-onnx" 2>/dev/null | tr '\n' ' '
    echo
    echo "----"
} > "$out"
PATH="$tree/third_party/sherpa-onnx:$PATH" \
    env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT \
    bash scripts/portable-tests.sh --scope=windows >> "$out" 2>&1
rc=$?
echo "RC=$rc" >> "$out"
echo "end=$(date '+%F %T %z')" >> "$out"
echo "rc=$rc"
