#!/usr/bin/env bash
# 287-v1b 门禁复跑（派单第 0 步第 1 项的五件＋格式两把）
# 射程＝HEAD 的工作树（仓根）；每把尺各写一枚 .log ＋一枚 .rc
# 用法：sh .scratch/wisp/probes/287/v1b/gates.sh <outdir>
set -u
REPO="D:/work/workspace/projects plans/Wisp"
OUT="${1:?outdir required}"
cd "$REPO" || exit 99
mkdir -p "$OUT"

run() { # name, then command...
  local name="$1"; shift
  local log="$OUT/$name.log"
  "$@" >"$log" 2>&1
  local rc=$?
  printf '%s\n' "$rc" >"$OUT/$name.rc"
  printf '%s rc=%s log=%s bytes=%s\n' "$name" "$rc" "$log" "$(wc -c <"$log" | tr -d ' ')"
}

run v1b-01-build       env GOFLAGS= go build ./...
run v1b-02-vet-tools   go vet ./internal/tools/
run v1b-03-d22scan     sh scripts/d22scan.sh
run v1b-04-pathlength  sh scripts/check-path-length-budget.sh --with-self-test
run v1b-05-gofmt-two   gofmt -l internal/tools/subagent_depth_fixture_287_test.go internal/tools/subagent_197.go
run v1b-06-gofmt-dir   gofmt -l internal/tools
echo "HEAD=$(git rev-parse HEAD)"
