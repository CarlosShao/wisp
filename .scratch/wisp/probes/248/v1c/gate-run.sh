#!/usr/bin/env bash
# 248-v1c gate rulers: sequential, one heavy job at a time (timing-sensitive tests).
set -u
ROOT=$(cd "$(dirname "$0")/../../../../.." && pwd)
cd "$ROOT" || exit 1
OUT=.scratch/wisp/probes/248/v1c
mkdir -p "$OUT"
export PATH="$ROOT/third_party/sherpa-onnx:$ROOT/build:$PATH"
GOF="$(go env GOPATH)/bin/gofumpt"
note() { echo "@@@ $* @@@"; }

note "start"; date -Iseconds; git rev-parse HEAD; pwd

note "GOFLAGS= go build ./..."
GOFLAGS= go build ./... > "$OUT/build.txt" 2>&1; note "build rc=$?"; tail -5 "$OUT/build.txt"

note "go vet ./cmd/wisp ./internal/panel ./internal/config"
GOFLAGS= go vet ./cmd/wisp ./internal/panel ./internal/config > "$OUT/vet.txt" 2>&1; note "vet rc=$?"; tail -5 "$OUT/vet.txt"

note "gofumpt -l cmd/wisp internal/config internal/panel"
"$GOF" -l cmd/wisp internal/config internal/panel > "$OUT/gofumpt-touched.txt" 2>&1
note "gofumpt rc=$? list-count=$(wc -l < "$OUT/gofumpt-touched.txt")"; cat "$OUT/gofumpt-touched.txt"

note "d22scan"
sh scripts/d22scan.sh > "$OUT/d22scan.txt" 2>&1; note "d22scan rc=$?"; tail -3 "$OUT/d22scan.txt"

note "go test ./internal/panel ./internal/config -count=1"
GOFLAGS= go test ./internal/panel ./internal/config -count=1 > "$OUT/test-panel-config.txt" 2>&1; note "panel+config rc=$?"
grep -E '^(---|\s+---) FAIL|^FAIL|^ok ' "$OUT/test-panel-config.txt" | head -30

note "go test ./cmd/wisp -count=1 (full package)"
GOFLAGS= go test ./cmd/wisp -count=1 -timeout 25m > "$OUT/test-cmdwisp-full.txt" 2>&1; note "cmd/wisp rc=$?"
grep -E '^(---|\s+---) FAIL|^FAIL|^ok ' "$OUT/test-cmdwisp-full.txt" | head -40

note "done"; date -Iseconds
