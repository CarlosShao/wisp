#!/usr/bin/env bash
# Orchestrator's own before/after whole-package roster for ticket 303 AC#3's gate.
# Runs BOTH sides in the same external clone (like-for-like surface); the main
# worktree is never checked out (shared worktree: other legs read it).
# Read-only w.r.t. the repo; creates only files under probes/303/orch/.
set -u
CL="$HOME/wisp-303-bisect"
REPO="/d/work/workspace/projects plans/Wisp"
OUTDIR="$REPO/.scratch/wisp/probes/303/orch"
export PATH="$REPO/third_party/sherpa-onnx:$REPO/build:$PATH"

echo "start_local=$(date '+%F %T %z')"
echo "clone_head_before=$(git -C "$CL" log -1 --format='%h')"
echo "main_head=$(git -C "$REPO" log -1 --format='%h')"
echo "scoped_porcelain_main=$(git -C "$REPO" status --porcelain -- scripts internal .github docs cmd | wc -l)"

# ---- BEFORE = the commit the fix's parent points at (pre-fix code) ----
git -C "$CL" checkout -q --detach cf46c24a
echo "detached=$(git -C "$CL" log -1 --format='%h') panel_host_lines=$(git -C "$CL" show HEAD:cmd/wisp/panel_host_windows.go | wc -l)"
go test -C "$CL" ./cmd/wisp/ -count=1 -timeout 900s -v > "$OUTDIR/r5-orch-package-before.txt" 2>&1
echo "rc_before=$?" >> "$OUTDIR/r5-orch-package-before.txt"
grep -E '^--- (FAIL|SKIP): ' "$OUTDIR/r5-orch-package-before.txt" | sed 's/ (.*//' | sort > "$OUTDIR/r5-orch-package-before-names.txt"
echo "before_fail=$(grep -c '^--- FAIL' "$OUTDIR/r5-orch-package-before.txt") before_skip=$(grep -c '^--- SKIP' "$OUTDIR/r5-orch-package-before.txt")"

# ---- AFTER = the fix ----
git -C "$CL" checkout -q --detach 807497c1
echo "detached=$(git -C "$CL" log -1 --format='%h')"
go test -C "$CL" ./cmd/wisp/ -count=1 -timeout 900s -v > "$OUTDIR/r6-orch-package-after.txt" 2>&1
echo "rc_after=$?" >> "$OUTDIR/r6-orch-package-after.txt"
grep -E '^--- (FAIL|SKIP): ' "$OUTDIR/r6-orch-package-after.txt" | sed 's/ (.*//' | sort > "$OUTDIR/r6-orch-package-after-names.txt"
echo "after_fail=$(grep -c '^--- FAIL' "$OUTDIR/r6-orch-package-after.txt") after_skip=$(grep -c '^--- SKIP' "$OUTDIR/r6-orch-package-after.txt")"

echo "=== NEW REDS (after minus before, per-name) ==="
comm -13 "$OUTDIR/r5-orch-package-before-names.txt" "$OUTDIR/r6-orch-package-after-names.txt" | grep 'FAIL' || echo "(none)"
echo "=== FIXED (before-red minus after-red) ==="
comm -23 <(grep 'FAIL' "$OUTDIR/r5-orch-package-before-names.txt") <(grep 'FAIL' "$OUTDIR/r6-orch-package-after-names.txt") || echo "(none)"
echo "=== restored clone to origin/dev ==="
git -C "$CL" checkout -q --detach "$(git -C "$REPO" rev-parse origin/dev)"
echo "clone_head_end=$(git -C "$CL" log -1 --format='%h')"
echo "end_local=$(date '+%F %T %z')"
