#!/usr/bin/env bash
# 342-v1 race re-run: mandatory question (i). One container invocation runs all five
# shapes against the same world + racer, so the numbers are comparable within the run.
# Usage: bash v342v1-race.sh <tree> <tag> ; repeats are separate invocations.
set -u
S=/d/tmp/wisp/20261018/342v1
TREE="$1"
TAG="$2"
N="${3:-1500}"
cp "$S/probe/v342v1-posix-race.go.txt" "$S/$TREE/internal/winsec/v342v1race_test.go"
echo "== pre-hash =="; sha256sum "$S/$TREE/internal/winsec/v342v1race_test.go"
bash "$S/scripts/v342v1-docker.sh" "$TREE" \
  "V342V1_ATTEMPTS=$N go test -count=1 -timeout 40m -v -run TestV342v1PosixAncestorRaceRows ./internal/winsec/" \
  "10-race-$TAG.log"
echo "-- rows:"
grep -E '^    v342v1race_test.go:[0-9]+: V342V1 V(BASE|PROD|JIA|YI|BING): attempts' "$S/logs/10-race-$TAG.log" || true
echo "-- damage rows:"
grep -E 'V342V1 V[A-Z]+ (leg|damageRow)' "$S/logs/10-race-$TAG.log" | head -80 || true
