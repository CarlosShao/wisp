#!/usr/bin/env bash
# 342-a1 HOST driver (g3): runs docker per phase. ⛔ the shared working tree; exported
# trees under /d/tmp/wisp/20261018/342a1 only. MSYS_NO_PATHCONV=1 is set by the caller.
set -u
S=/d/tmp/wisp/20261018/342a1
OUT=$S/logs
IMG=golang:1.27

crun () { # $1 = host tree dir, $2 = sh command, $3 = log
  MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
    -v "$1:/src" -v /d/work/base/gopath/pkg/mod:/gomod:ro -v "$S/gobuild:/gobuild" -w /src "$IMG" \
    sh -c "GOMODCACHE=/gomod GOCACHE=/gobuild $2" > "$3" 2>&1
  echo "rc=$?"
}

four () {
  echo "RUN=$(grep -c '^=== RUN' "$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$1" || true) SKIP=$(grep -c '^--- SKIP' "$1" || true)"
}

echo "=== (a) hypothetical landing v2 FULL roster (⚠ the descent names its failures open <comp>/ELOOP-wording = NOT the landing-leg shape; this row exists to show WHICH pins a mis-shaped landing breaks) ==="
crun "$S/mut-jia-linux-v2" "go test -count=1 -timeout 40m -v ./internal/winsec/" "$OUT/26-container-simlanding-jia-v2-full.log"; four "$OUT/26-container-simlanding-jia-v2-full.log"
grep -E '^--- FAIL' "$OUT/26-container-simlanding-jia-v2-full.log" | sed 's/ ([0-9.]*s)//' > "$OUT/26-jia-v2-failnames.txt" || true
echo "failnames:"; cat "$OUT/26-jia-v2-failnames.txt" || true

echo "=== (b) cheap variant: pinned parent + leaf O_EXCL (乙's flag-face wrong way) ==="
crun "$S/mut-yinofb-linux" "go test -count=1 -timeout 40m -v ./internal/winsec/" "$OUT/27-container-simlanding-yinofb-full.log"; four "$OUT/27-container-simlanding-yinofb-full.log"
grep -E '^--- FAIL' "$OUT/27-container-simlanding-yinofb-full.log" | sed 's/ ([0-9.]*s)//' > "$OUT/27-yinofb-failnames.txt" || true
echo "failnames:"; cat "$OUT/27-yinofb-failnames.txt" || true

echo "=== (c) landed-form paired racer: production PrivateFile on the two POSIX landings ==="
for form in jia-linux yi-linux; do
  cp "$S/probe/census342a1_probe_test.go" "$S/mut-$form/internal/winsec/" || { echo "cp failed $form"; continue; }
  crun "$S/mut-$form" "C342_ATTEMPTS=1500 go test -count=1 -timeout 30m -v -run TestCensus342a1A2 ./internal/winsec/" "$OUT/28-container-A2-landed-$form.log"
  echo "--- mut-$form:"
  grep -E 'CENSUS A2|--- (PASS|FAIL)' "$OUT/28-container-A2-landed-$form.log" | head -14 || true
  sha256sum "$S/mut-$form/internal/winsec/winsec.go"
done

echo "=== (d) consumers on the two POSIX landing faces ==="
for form in jia-linux yi-linux; do
  crun "$S/mut-$form" "go test -count=1 -timeout 40m ./internal/risk/ ./internal/secret/ ./internal/memory/ ./internal/config/ ./internal/agent/" "$OUT/29-container-consumers-$form.log"
  echo "--- consumers mut-$form:"; tail -7 "$OUT/29-container-consumers-$form.log"
done
echo "DRIVER-G3-DONE"
