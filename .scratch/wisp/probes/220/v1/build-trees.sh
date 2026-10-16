#!/usr/bin/env bash
# 220-v1 台面建树（裁决腿自用·仓外导出树）
# 定式来源＝编排者 220v1-dispatch 更正 (2)(3)＋A1002；本脚本自带三枚台面缺件补齐：
#   ① 建树七件套带 frontend（internal/panel/assets.go:25 引用 frontend 包，缺＝setup failed 假红）
#   ② build/ 与 third_party/sherpa-onnx/ 零跟踪件 ⇒ 手工拷三枚 .dll，另拷一份进 cmd/wisp/ 旁
#   ③ frontend/dist 只 1 枚跟踪件 ⇒ 整目录从工作树拷，否则真窗用例读到空白页
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ROOT=/tmp/wisp220v1
PROBE="$REPO/.scratch/wisp/probes/220/v1"
mkdir -p "$PROBE/logs"
LOG="$PROBE/logs/00-tree.txt"
: > "$LOG"

build_tree() { # ref out
  local ref="$1" out="$2"
  rm -rf "$out" 2>/dev/null
  mkdir -p "$out"
  ( cd "$REPO" && git archive "$ref" go.mod go.sum internal cmd tools scripts frontend ) | tar -x -C "$out"
  cp -r "$REPO/frontend/dist/." "$out/frontend/dist/" 2>/dev/null
  mkdir -p "$out/build" "$out/third_party/sherpa-onnx"
  cp -p "$REPO/build/"*.dll "$out/build/" 2>/dev/null
  cp -p "$REPO/third_party/sherpa-onnx/"*.dll "$out/third_party/sherpa-onnx/" 2>/dev/null
  cp -p "$REPO/build/"*.dll "$out/cmd/wisp/" 2>/dev/null
  local n gomod
  n=$(find "$out" -name '*.go' | wc -l)
  gomod=$(find "$out" -name '*.go' -path '*internal*' | wc -l)
  {
    echo "tree=$out ref=$ref"
    echo "  go_files_total=$n go_files_internal=$gomod"
    echo "  dist_files=$(find "$out/frontend/dist" -type f | wc -l) third_party_dll=$(ls "$out/third_party/sherpa-onnx"/*.dll 2>/dev/null | wc -l) build_dll=$(ls "$out/build"/*.dll 2>/dev/null | wc -l) cmdwisp_dll=$(ls "$out/cmd/wisp"/*.dll 2>/dev/null | wc -l)"
  } >> "$LOG" 2>&1
}

mk_env() { # tree -> prints GOTMPDIR native
  local got="$1/.gotmp"
  mkdir -p "$got"
  cygpath -w "$got"
}

echo "=== build start $(date '+%F %T %z') HEAD=$(cd "$REPO" && git rev-parse HEAD) ===" >> "$LOG"
build_tree "$(cd "$REPO" && git rev-parse HEAD)" "$ROOT/bench"
build_tree "$(cd "$REPO" && git rev-parse HEAD)" "$ROOT/full"
build_tree fd01e7962 "$ROOT/pre"
echo "=== go build check $(date '+%F %T %z') ===" >> "$LOG"
for t in bench full pre; do
  ng=$(mk_env "$ROOT/$t")
  ( cd "$ROOT/$t" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
      go build ./cmd/wisp/ >"$ROOT/$t/build.log" 2>&1; echo "  tree=$t go_build_rc=$?" >> "$LOG" )
  tail -3 "$ROOT/$t/build.log" >> "$LOG" 2>&1
done
echo "=== build end $(date '+%F %T %z') ===" >> "$LOG"
cat "$LOG"
