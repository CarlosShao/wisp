#!/bin/bash
# 编排者收件复跑（328-v1 决定性那一发＝AC#2 的"边界"单独敏感性）：⛔ 抄判语件的读数
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/v1-verify"
VANCHOR=5db0b538
F=internal/panel/subagent_roster_197.go
LOG="$W/logs/orch-v1-mut.log"
mkdir -p "$W/tree" "$W/gotmp" "$W/logs" "$W/tdst"
cd "$REPO" || exit 9
{
echo "=== [1] 建树 @ $VANCHOR ==="
git archive "$VANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs > "$W/tree.tar"
tar -xf "$W/tree.tar" -C "$W/tree"
echo "go files = $(find "$W/tree" -name '*.go' | wc -l)"
mkdir -p "$W/tree/build" "$W/tree/third_party/sherpa-onnx"
cp -f build/*.dll "$W/tree/build/" 2>/dev/null || echo "WARN build dll"
cp -f third_party/sherpa-onnx/*.dll "$W/tree/third_party/sherpa-onnx/" 2>/dev/null || echo "WARN sherpa dll"
[ -d frontend/dist ] && cp -r frontend/dist "$W/tree/frontend/" 2>/dev/null
echo "=== [2] 锚 blob 对撞（判语件说 2da6c596…） ==="
echo "anchor blob      = $(git rev-parse $VANCHOR:$F)"
echo "tree file sha256 = $(sha256sum "$W/tree/$F" | cut -d' ' -f1)"
echo "tree gitblob     = $(git hash-object "$W/tree/$F")"
echo "third-arm line   = $(grep -n 'HasPrefix(key, task' "$W/tree/$F")"
echo "=== [3] 正控：未突变树上该用例必须全绿（7 枚子用例） ==="
cd "$W/tree" || exit 9
export GOTMPDIR="$W/gotmp" TMPDIR="$W/tdst" TEMP="$W/tdst" TMP="$W/tdst"
go test ./internal/panel/ -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL)' | head -12
echo "PRISTINE_RC=${PIPESTATUS[0]}"
echo "=== [4] 突变：HasPrefix(key, task+\"#\") -> HasPrefix(key, task)（去掉边界） ==="
cd "$REPO" || exit 9
git -C "$REPO" cat-file blob "$VANCHOR:$F" | sed 's/strings\.HasPrefix(key, task+"#")/strings.HasPrefix(key, task)/' > "$W/tree/$F"
cd "$W/tree" || exit 9
echo "mutated line     = $(grep -n 'HasPrefix(key, task' "$F")"
echo "mutated sha256   = $(sha256sum "$F" | cut -d' ' -f1)"
echo "expect           = 31b8dcdc4cfca4ef86f5cf4cf8d8c9933b613b622013f73cf9f93add7c64c70a"
echo "=== [5] child-1x 那枚必须红、其余五枚必须绿 ==="
go test ./internal/panel/ -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(    --- (PASS|FAIL)|--- (PASS|FAIL)|ok  |FAIL|    subagent_blocked)' | head -14
echo "MUT_RC=${PIPESTATUS[0]}"
echo "=== [6] 还原＋复验＋复跑绿 ==="
cd "$REPO" || exit 9
git -C "$REPO" cat-file blob "$VANCHOR:$F" > "$W/tree/$F"
cd "$W/tree" || exit 9
echo "restored sha256  = $(sha256sum "$F" | cut -d' ' -f1)"
echo "expect pristine  = 450e1d1c47ee715fb74b56116a827988b0d9bfbc261fc06748c4704b4bd87028"
go test ./internal/panel/ -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(--- (PASS|FAIL)|ok  |FAIL)' | head -5
echo "RESTORED_RC=${PIPESTATUS[0]}"
} > "$LOG" 2>&1
echo "LOG=$LOG"; tail -50 "$LOG"
