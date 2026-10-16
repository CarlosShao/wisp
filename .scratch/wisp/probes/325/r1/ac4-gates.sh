#!/bin/sh
# 325-r1 AC#4 gates. Each gate prints its own rc; nothing is chained through a pipe.
# Surfaces are stated per gate: the three tree-wide gates run on the OUTSIDE-THE-REPO
# export of 428ef21e (a sibling leg is writing internal/winsec in the shared worktree);
# the path-length gate needs a git index, which `git archive` does not carry, so that
# one runs read-only in the shared worktree and says so.
LOGD="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/325/r1/logs"
REPO="/d/work/workspace/projects plans/Wisp"
POST=/c/tmp/wisp325r1/post

cd "$POST" || exit 9

echo "== G1 bash -n tools/d22scan/runtests.sh (export tree) =="
bash -n tools/d22scan/runtests.sh
echo "G1-rc=$?"

echo "== G2 sh scripts/d22scan.sh (export tree, GOTMPDIR/WISP_TEST_GOTMP_ROOT drained) =="
env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT sh scripts/d22scan.sh > "$LOGD/a4-g2-d22scan.txt" 2>&1
echo "G2-rc=$?"

echo "== G3 sh scripts/check-path-length-budget.sh --with-self-test (SHARED WORKTREE, read-only) =="
cd "$REPO" || exit 9
sh scripts/check-path-length-budget.sh --with-self-test > "$LOGD/a4-g3-length.txt" 2>&1
echo "G3-rc=$?"

echo "== G4 GOFLAGS= go build ./... (export tree) =="
cd "$POST" || exit 9
GOFLAGS= go build ./... > "$LOGD/a4-g4-gobuild.txt" 2>&1
echo "G4-rc=$?"

echo "== G5 (not a gate, a finding) which gofumpt on this host =="
which gofumpt
echo "G5-which-gofumpt-rc=$?"
