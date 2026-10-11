#!/usr/bin/env bash
# 306-v1 gate rulers. Read-only in the shared worktree; whole-package build only in the export tree.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE=$(cat "$REPO/.scratch/wisp/probes/306/v1/logs/TREE-PATH.txt" | tr -d '\r\n ')
LOG="$REPO/.scratch/wisp/probes/306/v1/logs"
SUM="$LOG/ZZ-gates.txt"
cd "$REPO" || exit 9
: > "$SUM"
note() { echo "$@" >> "$SUM"; }

note "=== G0 tree path = $TREE ==="
note "G0_head=$(git rev-parse HEAD)"

# G1 go vet internal/audio in the shared worktree (in-lane, read-only for the tree)
GOFLAGS= go vet ./internal/audio/ > "$LOG/g1-vet-worktree.txt" 2>&1; rc=$?
note "G1 rc_go_vet_internal_audio_worktree=$rc  lines=$(wc -l < "$LOG/g1-vet-worktree.txt")"

# G2 d22scan gate in the worktree (read-only scan; its own module)
sh scripts/d22scan.sh > "$LOG/g2-d22scan-worktree.txt" 2>&1; rc=$?
note "G2 rc_d22scan_worktree=$rc"
grep -nE 'ban|Ban|#8|scanned|files' "$LOG/g2-d22scan-worktree.txt" | tail -30 | sed 's/^/  [G2] /' >> "$SUM"

# G3 whole-package build ONLY in the export tree (worktree ./... is the orchestrator lane)
(cd "$TREE" && GOFLAGS= go build ./... ) > "$LOG/g3-build-all-exporttree.txt" 2>&1; rc=$?
note "G3 rc_go_build_all_EXPORTTREE=$rc  (scope: git archive e4740e35 tree, NOT the shared worktree)  lines=$(wc -l < "$LOG/g3-build-all-exporttree.txt")"
# G3b same in the export tree for internal/audio only
(cd "$TREE" && GOFLAGS= go vet ./internal/audio/ ) > "$LOG/g3b-vet-exporttree.txt" 2>&1; rc=$?
note "G3b rc_go_vet_internal_audio_EXPORTTREE=$rc"

# G4 gofmt ruler 1: worktree, scope internal/audio (recursive)
gofmt -l internal/audio > "$LOG/g4-gofmt-worktree-internalaudio.txt" 2>&1; rc=$?
note "G4 rc_gofmt_worktree=$rc roster_count=$(wc -l < "$LOG/g4-gofmt-worktree-internalaudio.txt") scope=worktree:internal/audio"
cat "$LOG/g4-gofmt-worktree-internalaudio.txt" | sed 's/^/  [G4] /' >> "$SUM"
gofmt -l . > "$LOG/g4b-gofmt-worktree-fullrepo.txt" 2>&1; rc=$?
note "G4b rc_gofmt_worktree_fullrepo=$rc roster_count=$(wc -l < "$LOG/g4b-gofmt-worktree-fullrepo.txt") scope=worktree:whole_repo_tracked_and_untracked_files"
grep -cE '^internal/audio/' "$LOG/g4b-gofmt-worktree-fullrepo.txt" | sed 's/^/  [G4b internal/audio hits] /' >> "$SUM"

# G5 gofmt ruler 2: HEAD blob + baseline blob, exported to out-of-repo dirs
B54="$TREE/_306v1-blobs"; mkdir -p "$B54/5480434f" "$B54/e4740e35"
git archive 5480434f internal/audio | tar -x -C "$B54/5480434f"; rc1=$?
git archive e4740e35 internal/audio | tar -x -C "$B54/e4740e35"; rc2=$?
note "G5 rc_archive_5480434f=$rc1 rc_archive_HEAD=$rc2"
for sha in 5480434f e4740e35; do
  gofmt -l "$B54/$sha" > "$LOG/g5-gofmt-blob-$sha.txt" 2>&1; rc=$?
  n=$(grep -c . "$LOG/g5-gofmt-blob-$sha.txt")
  tot=$(find "$B54/$sha/internal/audio" -name '*.go' | wc -l)
  note "G5 gofmt -l blob=$sha scope=internal/audio rc=$rc unformatted_count=$n total_go_files_in_scope=$tot"
  cat "$LOG/g5-gofmt-blob-$sha.txt" | sed "s/^/  [G5 $sha] /" >> "$SUM"
done

# G6 the two ticket-AC#4 test-side counts + ban#8 roster from d22scan output
grep -nE '#8|emoji|Emoji' "$LOG/g2-d22scan-worktree.txt" | sed 's/^/  [G6 ban8] /' >> "$SUM"
tail -25 "$LOG/g2-d22scan-worktree.txt" | sed 's/^/  [G6 d22tail] /' >> "$SUM"

note "=== done $(date) ==="
