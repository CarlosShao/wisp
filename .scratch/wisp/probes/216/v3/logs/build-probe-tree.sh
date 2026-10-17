#!/usr/bin/env bash
# 216-v3：把只读探针（probe-216v3_test.go.txt）拷进【另一枚】仓外导出树跑一次，
# 用它打印 formsOf/anchorForms/uncertainAnchorMiss 的实值。⛔ 突变树（D:\W216V3MUT）
# 保持原样、⛔ 跟踪字节、⛔ 共享工作树。两把台面各跑一发，落点差异就是判据本身。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
HERE="$(cd "$(dirname "$0")/.." && pwd)"
ANCHOR=7a94c3bd
ROOT=/d/W216V3PROBE
EIGHT="go.mod go.sum internal cmd tools scripts frontend models"
mkdir -p /d/W216V3TMP "$HERE/logs"
export GOPATH='D:\work\base\gopath'
echo "GOTMPDIR_in_env=[${GOTMPDIR:-<unset>}]"
echo "########## 0 建探针树（锚 $ANCHOR）"
rm -rf "$ROOT"; mkdir -p "$ROOT"
( cd "$REPO" && git archive "$ANCHOR" $EIGHT ) | tar -x -C "$ROOT"
( cd "$REPO" && git archive "$ANCHOR" docs ) | tar -x -C "$ROOT"
mkdir -p "$ROOT/build" "$ROOT/third_party/sherpa-onnx" "$ROOT/frontend/dist"
cp "$REPO"/build/*.dll "$ROOT/build/" 2>/dev/null
cp "$REPO"/third_party/sherpa-onnx/*.dll "$ROOT/third_party/sherpa-onnx/" 2>/dev/null
cp -r "$REPO"/frontend/dist/. "$ROOT/frontend/dist/" 2>/dev/null
echo "go_files=$(find "$ROOT/internal" "$ROOT/cmd" "$ROOT/tools" -name '*.go' | wc -l)"
cp "$HERE/logs/probe-216v3_test.go.txt" "$ROOT/internal/risk/zz_probe_216v3_test.go"
echo "probe_installed=$ROOT/internal/risk/zz_probe_216v3_test.go"
( cd "$ROOT" && go vet ./internal/risk/ >/dev/null 2>"$HERE/logs/probe-216v3-vet.log"; echo "go_vet_rc=$?" >>"$HERE/logs/probe-216v3-vet.log" )
echo "########## 1 探针·C 卷默认台面（t.TempDir() 在主目录祖先链内）"
( cd "$ROOT" && go test ./internal/risk/ -count=1 -timeout 20m -v -run TestProbe216V3FailClosedChain ) >"$HERE/logs/probe-216v3-chain-Cdefault.log" 2>&1
echo "rc=$?" >>"$HERE/logs/probe-216v3-chain-Cdefault.log"
echo "########## 2 探针·D 台面（t.TempDir() 在主目录之外）"
( cd "$ROOT" && env TMP='D:\W216V3TMP' TEMP='D:\W216V3TMP' go test ./internal/risk/ -count=1 -timeout 20m -v -run TestProbe216V3FailClosedChain ) >"$HERE/logs/probe-216v3-chain-Dtmp.log" 2>&1
echo "rc=$?" >>"$HERE/logs/probe-216v3-chain-Dtmp.log"
echo "########## 3 探针树里被裁那两枚件仍等于锚上 blob（突变没串台）"
for f in rules_shell.go blacklist.go allow_substring_pin_216_test.go name_claim_two_arm_216_test.go; do
  echo "$f tree=$(git hash-object "$ROOT/internal/risk/$f") blob=$(cd "$REPO" && git rev-parse "$ANCHOR:internal/risk/$f")"
done
echo "########## 4 共享工作树仍空（应空）"
( cd "$REPO" && git status --porcelain -- internal/risk )
echo "########## done"
