#!/bin/sh
# 260-r1 mutation driver:摘掉新逻辑 ⇒ 指名用例必红；还原 ⇒ 绿且 md5 一致、git 面干净。
# Each mutation is one exact string replacement in internal/ball/hotkey_windows.go,
# applied over the FINAL committed shape and restored from the byte-identical backup
# (.scratch/wisp/probes/260/r1/hotkey_windows.go.260r1-finalbackup) right after.
# Temp files are created, never deleted (issues/README rule 8).
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../../../" && pwd)
cd "$root"
if [ ! -f go.mod ]; then
    echo "mutation-driver.sh: derived root $root has no go.mod - refusing" >&2
    exit 2
fi

prod=internal/ball/hotkey_windows.go
names='TestBorrowDefaultIsBitIdentical260|TestBorrowFollowsConfiguredBinding260|TestBorrowUnparsableFallsBackLoudly260|TestBorrowReceiptSharesOneSourceWithRegistration260|TestBorrowFailureNamesTheKeyItTried260|TestConfiguredCancelStillNeverBoundWhileIdle260|TestLiveBorrowHelperPremiseIsConfigDependent260'

run_case() {
    tag=$1
    python "$here/mutate.py" "$tag" apply
    echo "=== $tag: build ==="
    go build ./internal/ball/ 2>&1
    echo "=== $tag: go test -count=1 -v -run '260' ./internal/ball/ ==="
    go test -count=1 -v -run "$names" ./internal/ball/ 2>&1 || true
    cp "$here/hotkey_windows.go.260r1-finalbackup" "$prod"
    echo "=== $tag: restored, md5 must equal the baseline ==="
    md5sum "$prod"
}

run_case M1
run_case M2
run_case M3

echo "=== final: all three mutations restored, gate re-read ==="
md5sum "$prod" internal/ball/ball_windows.go
git status --short internal/ball/ || true
go test -count=1 -v -run "$names" ./internal/ball/ 2>&1 | grep -E '^(--- |ok|FAIL)'
