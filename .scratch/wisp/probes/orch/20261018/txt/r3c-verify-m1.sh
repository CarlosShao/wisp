#!/bin/bash
# 编排者收件复跑（228-r3c 决定性那一发）：M1 突变 ⇒ TestAC11TrayExitTriggerWalksTheTenSteps 必须点名红
# 落点守新铁律：D:\tmp\wisp\20261018\orch\r3c-verify （⛔ 盘符根）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/r3c-verify"
ANCHOR=c11163c896b09e7b6863d347dadbcb65702585cc
F=internal/proc/boot_windows.go
LOG="$W/logs/orch-m1-recheck.log"
mkdir -p "$W/tree" "$W/gotmp" "$W/logs" "$W/tdst"
cd "$REPO" || exit 9
{
echo "=== [1] 建树（八件套＋docs＋手工三件） ==="
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs > "$W/tree.tar"
tar -xf "$W/tree.tar" -C "$W/tree"
echo "go files = $(find "$W/tree" -name '*.go' | wc -l)"
mkdir -p "$W/tree/build" "$W/tree/third_party/sherpa-onnx"
cp -f build/*.dll "$W/tree/build/" 2>/dev/null || echo "WARN build dll"
cp -f third_party/sherpa-onnx/*.dll "$W/tree/third_party/sherpa-onnx/" 2>/dev/null || echo "WARN sherpa dll"
[ -d frontend/dist ] && cp -r frontend/dist "$W/tree/frontend/" 2>/dev/null
echo "=== [2] pristine 断言（突变前） ==="
sha256sum "$W/tree/$F" | cut -d' ' -f1
git hash-object "$W/tree/$F"
echo "expect sha256=3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7 sha1=531060a55c080ebd535b426b8c60dcd9011785fc"
echo "=== [3] 正控：未突变树上该用例必须绿 ==="
cd "$W/tree" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
export GOTMPDIR="$W/gotmp" TMPDIR="$W/tdst" TEMP="$W/tdst" TMP="$W/tdst"
go test ./cmd/wisp/ -run 'TestAC11TrayExitTriggerWalksTheTenSteps' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL|    resident_tray)' | head -8
echo "PRISTINE_RC=${PIPESTATUS[0]}"
echo "=== [4] M1 突变（逐字节＝腿 mutate-shot.sh 的 MUT_LINE） ==="
awk '{ if ($0 ~ /^\tstopSignal := rt\.StopSignal\(\)[[:space:]]*$/) print "\tstopSignal := (<-chan struct{})(nil) // MUTATION-228r3c-M1 severed-loop-stop-check"; else print }' "$F" > "$F.mut"
mv "$F.mut" "$F"
sha256sum "$F" | cut -d' ' -f1; git hash-object "$F"; wc -c < "$F"
echo "expect sha256=3856c24aa1c1ee25b1a351d99fb91ec8ba7bb475cabeea793796944cfcd57fb3 sha1=54e0b8801d025be8d35659eeeb56f9e8c980f832 bytes=12969"
echo "=== [5] 该用例必须点名红 ==="
go test ./cmd/wisp/ -run 'TestAC11TrayExitTriggerWalksTheTenSteps' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|FAIL|    resident_tray_stop)' | head -8
echo "MUT_RC=${PIPESTATUS[0]}"
echo "=== [6] 还原＋复验哈希＋复跑绿 ==="
cd "$REPO" || exit 9
git -C "$REPO" cat-file blob "$ANCHOR:$F" > "$W/tree/$F"
cd "$W/tree" || exit 9
sha256sum "$F" | cut -d' ' -f1; git hash-object "$F"
go test ./cmd/wisp/ -run 'TestAC11TrayExitTriggerWalksTheTenSteps' -count=1 -timeout 30m -v 2>&1 | grep -E '^(--- (PASS|FAIL)|ok  |FAIL)' | head -5
echo "RESTORED_RC=${PIPESTATUS[0]}"
} > "$LOG" 2>&1
echo "LOG=$LOG"
tail -60 "$LOG"
