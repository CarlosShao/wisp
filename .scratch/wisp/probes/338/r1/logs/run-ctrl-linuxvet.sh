#!/bin/sh
# 票 338 落地腿 338-r1 · lead 甲形条件④：GOOS=linux go vet ./cmd/wisp/ 的存量归位尺
# 改后锚（HEAD）已跑＝rc=1（sherpa-onnx-go-linux 的 build constraints 排除全部 Go 文件）
# 本枚＝改前锚 0b1d1f0b 同配方导出树复跑＝坐实这枚红与我的字节无关（别票/平台存量）
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
TB="$NEST/tree-ctrl-linux"
mkdir -p "$TB"
cd "/d/work/workspace/projects plans/Wisp" || exit 9
git archive 0b1d1f0b go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$TB"
git archive 0b1d1f0b docs | tar -x -C "$TB"
mkdir -p "$TB/build" "$TB/third_party/sherpa-onnx" "$TB/frontend/dist"
cp -r build/. "$TB/build/" 2>/dev/null
cp -r third_party/sherpa-onnx/. "$TB/third_party/sherpa-onnx/" 2>/dev/null
[ -d frontend/dist ] && cp -r frontend/dist/. "$TB/frontend/dist/"
echo "tree-ctrl-linux go=$(find "$TB" -name '*.go' | wc -l)"
unset GOTMPDIR TMPDIR
echo '===== 改前锚（0b1d1f0b）GOOS=linux go vet ./cmd/wisp/ ====='
( cd "$TB" && GOOS=linux go vet ./cmd/wisp/ ) > "$ROOT/ac2b-ctrl-linuxvet.txt" 2>&1
echo "rc_ctrl_linux_vet=$?"
cat "$ROOT/ac2b-ctrl-linuxvet.txt"
echo '===== 同树本机档正控 go vet ./cmd/wisp/ ====='
( cd "$TB" && go vet ./cmd/wisp/ ) > "$ROOT/ac2b-ctrl-nativevet.txt" 2>&1
echo "rc_ctrl_native_vet=$?"
cat "$ROOT/ac2b-ctrl-nativevet.txt"
