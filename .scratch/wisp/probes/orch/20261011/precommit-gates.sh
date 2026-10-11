#!/bin/sh
# 编排者 commit 前门禁三把（本轮只动 .md：票 72 追加节＋票 308 新立＋台账 A838-A840＋HANDOVER 4.0cf）
cd "D:/work/workspace/projects plans/Wisp" || exit 99
OUT=.scratch/wisp/probes/orch/20261011/logs/orch-precommit-gates.txt
: > "$OUT"
GOFLAGS= go build ./... > .scratch/wisp/probes/orch/20261011/logs/tmp-build.txt 2>&1
echo "build_rc=$?" >> "$OUT"
go vet ./internal/risk/ > .scratch/wisp/probes/orch/20261011/logs/tmp-vet.txt 2>&1
echo "vet_rc=$?" >> "$OUT"
sh scripts/d22scan.sh > .scratch/wisp/probes/orch/20261011/logs/tmp-d22scan.txt 2>&1
echo "d22scan_rc=$?" >> "$OUT"
gofmt -l internal cmd > .scratch/wisp/probes/orch/20261011/logs/tmp-gofmt.txt 2>&1
echo "gofmt_rc=$? gofmt_n=$(grep -c . .scratch/wisp/probes/orch/20261011/logs/tmp-gofmt.txt)" >> "$OUT"
echo "go_side_diff=$(git diff --numstat -- internal cmd | wc -l)" >> "$OUT"
wc -c "$OUT" >> "$OUT"
cat "$OUT"
