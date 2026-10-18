#!/bin/sh
# 票 338 落地腿 338-r1 · AC#2 交件侧四道门（两枚都已归位；lead 甲形补偿条件④也在本枚落）
# ★先 unset GOTMPDIR TMPDIR（票 311 AC#1 shape B 卷闸＝A1027 第⑤节定式）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
unset GOTMPDIR TMPDIR
echo "GOTMPDIR=[${GOTMPDIR:-unset}] TMPDIR=[${TMPDIR:-unset}]"
echo "HEAD=$(git rev-parse --short HEAD)"

echo '===== [1] gofumpt -l internal cmd（交件侧名册，逐名） ====='
/d/work/base/gopath/bin/gofumpt.exe -l internal cmd | sed 's/\\/\//g' | sort
echo "rc_gofumpt=$?"

echo '===== [2] sh scripts/d22scan.sh ====='
sh scripts/d22scan.sh > /d/tmp/wisp/20261018/338r1/ac2b-d22scan.txt 2>&1
echo "rc_d22scan=$?"
tail -6 /d/tmp/wisp/20261018/338r1/ac2b-d22scan.txt

echo '===== [3] sh scripts/check-path-length-budget.sh ====='
sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261018/338r1/ac2b-len.txt 2>&1
echo "rc_len=$?"
tail -4 /d/tmp/wisp/20261018/338r1/ac2b-len.txt

echo '===== [4] go vet ./internal/tools/ ====='
go vet ./internal/tools/ > /d/tmp/wisp/20261018/338r1/ac2b-vet-tools.txt 2>&1
echo "rc_vet_tools=$?"; cat /d/tmp/wisp/20261018/338r1/ac2b-vet-tools.txt

echo '===== [5] go vet ./cmd/wisp/（本机档） ====='
go vet ./cmd/wisp/ > /d/tmp/wisp/20261018/338r1/ac2b-vet-wisp.txt 2>&1
echo "rc_vet_wisp=$?"; cat /d/tmp/wisp/20261018/338r1/ac2b-vet-wisp.txt

echo '===== [6] GOOS=linux go vet ./cmd/wisp/（lead 甲形条件④第二档） ====='
GOOS=linux go vet ./cmd/wisp/ > /d/tmp/wisp/20261018/338r1/ac2b-vet-wisp-linux.txt 2>&1
echo "rc_vet_wisp_linux=$?"; head -30 /d/tmp/wisp/20261018/338r1/ac2b-vet-wisp-linux.txt

echo '===== [7] go build ./...（lead 甲形条件④第三档；★派单禁在共享工作树跑 go build ⇒ 落仓外导出树） ====='
NEST=$(cat /d/tmp/wisp/20261018/338r1/nest-v2.txt)
TB2="$NEST/tree-build"
mkdir -p "$TB2"
git archive HEAD go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$TB2"
git archive HEAD docs | tar -x -C "$TB2"
mkdir -p "$TB2/build" "$TB2/third_party/sherpa-onnx" "$TB2/frontend/dist"
cp -r build/. "$TB2/build/" 2>/dev/null
cp -r third_party/sherpa-onnx/. "$TB2/third_party/sherpa-onnx/" 2>/dev/null
[ -d frontend/dist ] && cp -r frontend/dist/. "$TB2/frontend/dist/"
echo "tree-build go=$(find "$TB2" -name '*.go' | wc -l)"
( cd "$TB2" && go build ./... ) > /d/tmp/wisp/20261018/338r1/ac2b-build.txt 2>&1
echo "rc_build=$?"; head -20 /d/tmp/wisp/20261018/338r1/ac2b-build.txt

echo '===== [8] 六枚在册件＋CRLF 四枚：写面自证 ====='
git status --porcelain -- cmd/wisp internal/tools internal/agent internal/risk internal/config | cat
echo "status_rows=$(git status --porcelain -- cmd/wisp internal/tools internal/agent internal/risk internal/config | grep -c . || true)"
