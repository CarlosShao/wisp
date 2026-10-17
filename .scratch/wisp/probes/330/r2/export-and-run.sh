#!/usr/bin/env bash
# 330-r2: export-tree builder + whole-package colour run (ticket 330 G2 leg).
# usage: sh export-and-run.sh <ref> <treedir-name-under-D-tmp>
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
REF="$1"
NAME="$2"
TREE="/d/tmp/${NAME}"
mkdir -p "$TREE" "$TREE/logs" || exit 9
echo "ref=${REF}"
git archive "$REF" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$TREE" || exit 8
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx" "$TREE/frontend/dist"
cp build/*.dll "$TREE/build/" 2>/dev/null
cp third_party/sherpa-onnx/*.dll "$TREE/third_party/sherpa-onnx/" 2>/dev/null
cp -r frontend/dist/. "$TREE/frontend/dist/" 2>/dev/null
echo "go_files=$(find "$TREE" -name '*.go' | wc -l)"
echo "docs_present=$([ -d "$TREE/docs" ] && echo yes || echo no)"
mkdir -p /d/tmp/w330tmp || exit 7
cd "$TREE" || exit 9
export GOTMPDIR='D:\tmp\w330tmp'
go test ./internal/agent/ -count=1 -timeout 30m -v > "$TREE/logs/full-v.txt" 2>&1
echo "rc_full=$?"
echo "RUN=$(grep -c '^=== RUN' "$TREE/logs/full-v.txt")"
echo "PASS_TOP=$(grep -c '^--- PASS' "$TREE/logs/full-v.txt")"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$TREE/logs/full-v.txt")"
echo "SKIP=$(grep -c '^--- SKIP' "$TREE/logs/full-v.txt")"
echo "PANIC=$(grep -c '^panic' "$TREE/logs/full-v.txt")"
grep -E '^--- FAIL' "$TREE/logs/full-v.txt" | head -20
