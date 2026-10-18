#!/bin/bash
# 227-v1 MU-1: delete the schema_version line inside const matrixBase (manager_test.go:38)
# in the off-repo export tree, then re-run internal/config whole package.
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR='D:/tmp/wisp/20261018/227v1/gotmp'
unset TMPDIR
cd "$root/mu" || exit 2
echo "PRISTINE sha256 manager_test.go = $(sha256sum "$root/deliv/internal/config/manager_test.go" | cut -d' ' -f1)" > "$out/03-mu1.txt"
echo "MUTATED  sha256 manager_test.go = $(sha256sum internal/config/manager_test.go | cut -d' ' -f1)  bytes=$(wc -c < internal/config/manager_test.go)" >> "$out/03-mu1.txt"
echo "ANCHOR (content, not line number): the deleted line is the one inside const matrixBase; grep of the surviving line in deliv:" >> "$out/03-mu1.txt"
grep -n 'schema_version = 2' "$root/deliv/internal/config/manager_test.go" >> "$out/03-mu1.txt" || echo "NONE" >> "$out/03-mu1.txt"
grep -n 'const matrixBase' "$root/deliv/internal/config/manager_test.go" >> "$out/03-mu1.txt"
go test ./internal/config/ -count=1 -timeout 30m -v > "$out/03-mu1-full.txt" 2>&1
rc=$?
RUN=$(grep -c '^=== RUN' "$out/03-mu1-full.txt" || true)
PASSTOP=$(grep -c '^--- PASS' "$out/03-mu1-full.txt" || true)
FAILTOP=$(grep -c '^--- FAIL' "$out/03-mu1-full.txt" || true)
SKIP=$(grep -c '^--- SKIP' "$out/03-mu1-full.txt" || true)
echo "MUTATED tree=$root/mu rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP" >> "$out/03-mu1.txt"
echo "---- named FAIL roster ----" >> "$out/03-mu1.txt"
grep '^--- FAIL' "$out/03-mu1-full.txt" >> "$out/03-mu1.txt" || echo "NONE" >> "$out/03-mu1.txt"
echo "---- the :393 case red sentences, verbatim with file:line ----" >> "$out/03-mu1.txt"
grep -n -E 'writeguard_226_test\.go:[0-9]+|manager_test\.go:[0-9]+|writeguard_handedit_227_test\.go:[0-9]+' "$out/03-mu1-full.txt" >> "$out/03-mu1.txt" || echo "NONE" >> "$out/03-mu1.txt"
echo "MUT sha now:" >> "$out/03-mu1.txt"
sha256sum internal/config/manager_test.go >> "$out/03-mu1.txt"
echo "rc=$rc"
