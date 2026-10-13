#!/usr/bin/env bash
# Orchestrator's own ticket-295 AC#4 matched-pair collection.
# Same worktree @28b083d4 for all 4 shots: a,b = 改前 (overlay reverts exactly the 4 fix files
# to their ef8d16b3 = d5ef9cb9^ blobs), c1,c2 = 改后 (no overlay). Then the four gates.
# Ruler: content anchors + per-name roster diff; never absolute line numbers.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT="/c/Users/swq/AppData/Local/wisp-test-gotmp/orch295-pre-20261013"
cd "$REPO" || exit 9
LOG="$ROOT/logs/orch-295-pre-collect.txt"
say() { printf '%s\n' "$*" >> "$LOG"; }

say ""
say "=== launch TS=$(date '+%Y-%m-%d %H:%M:%S %z') ==="
say "head=$(git rev-parse HEAD)"
say "tasklist_family_hits=$(tasklist | grep -icE 'wisp|balldebug' || true) (0 expected)"

# --- alt-file hash proof vs 295-r1's recorded pre-fix hashes (05-pristine-hashes.txt first block) ---
say "=== alt hashes (pre-fix shape) ==="
for f in resident_audio_windows.go resident_audio_247_windows_test.go resident_audio_247_live_windows_test.go resident_mute_290_windows_test.go; do
  say "$(sha256sum "$ROOT/alt/$f" | cut -d' ' -f1)  $f"
done

say "=== 4 files worktree vs HEAD blob at launch ==="
for f in cmd/wisp/resident_audio_windows.go cmd/wisp/resident_audio_247_windows_test.go cmd/wisp/resident_audio_247_live_windows_test.go cmd/wisp/resident_mute_290_windows_test.go; do
  a=$(sha256sum "$f" | cut -d' ' -f1); b=$(git show "HEAD:$f" | sha256sum | cut -d' ' -f1)
  say "$f $([ "$a" = "$b" ] && echo MATCH || echo DIFF) $a"
done

run_shot() {
  tag="$1"; shift
  w="$ROOT/w-$tag"; mkdir -p "$w"
  say ""
  say "=== shot $tag ==="
  say "command_line=GOTMPDIR=$w PATH=... go test -count=1 -v ./cmd/wisp/ ./internal/audio/ $*"
  GOTMPDIR="$w" PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
    go test -count=1 -v ./cmd/wisp/ ./internal/audio/ "$@" > "$ROOT/logs/full-$tag.txt" 2>&1
  say "GO_TEST_RC=$?"
  say "lines=$(wc -l < "$ROOT/logs/full-$tag.txt") sha256=$(sha256sum "$ROOT/logs/full-$tag.txt" | cut -d' ' -f1)"
  say "census=$(grep -c -- '--- PASS' "$ROOT/logs/full-$tag.txt" || true) PASS / $(grep -c -- '--- FAIL' "$ROOT/logs/full-$tag.txt" || true) FAIL / $(grep -c -- '--- SKIP' "$ROOT/logs/full-$tag.txt" || true) SKIP / $(grep -c -- '=== RUN' "$ROOT/logs/full-$tag.txt" || true) RUN"
  grep -E '^(ok|FAIL|---)\s+github.com/CarlosShao/wisp' "$ROOT/logs/full-$tag.txt" | head -5 | sed 's/^/  /' >> "$LOG"
}

run_shot a "-overlay=$ROOT/ov-pre.json"
run_shot b "-overlay=$ROOT/ov-pre.json"
run_shot c1
run_shot c2

say ""
say "=== rosters ==="
for t in a b c1 c2; do
  grep -oE '^[[:space:]]*--- (PASS|FAIL|SKIP): [^ ]+' "$ROOT/logs/full-$t.txt" | sed 's/^[[:space:]]*//' | sed 's/: / /' | sort > "$ROOT/roster-$t.txt"
  say "$t roster_lines=$(wc -l < "$ROOT/roster-$t.txt")"
done

comm -12 "$ROOT/roster-a.txt" "$ROOT/roster-b.txt" > "$ROOT/roster-pre-intersect-mine.txt"
comm -12 "$ROOT/roster-c1.txt" "$ROOT/roster-c2.txt" > "$ROOT/roster-post-intersect-mine.txt"
say "pre_intersect lines=$(wc -l < "$ROOT/roster-pre-intersect-mine.txt") PASS=$(grep -c '^--- PASS' "$ROOT/roster-pre-intersect-mine.txt" || true) FAIL=$(grep -c '^--- FAIL' "$ROOT/roster-pre-intersect-mine.txt" || true) SKIP=$(grep -c '^--- SKIP' "$ROOT/roster-pre-intersect-mine.txt" || true)"
say "post_intersect lines=$(wc -l < "$ROOT/roster-post-intersect-mine.txt") PASS=$(grep -c '^--- PASS' "$ROOT/roster-post-intersect-mine.txt" || true) FAIL=$(grep -c '^--- FAIL' "$ROOT/roster-post-intersect-mine.txt" || true) SKIP=$(grep -c '^--- SKIP' "$ROOT/roster-post-intersect-mine.txt" || true)"
say "FAIL names per shot:"
for t in a b c1 c2; do say "  $t: $(grep '^--- FAIL' "$ROOT/roster-$t.txt" | tr '\n' '; ')"; done

say ""
say "--- flake check: a vs b ---"; diff "$ROOT/roster-a.txt" "$ROOT/roster-b.txt" >> "$LOG" 2>&1 || true
say "--- flake check: c1 vs c2 ---"; diff "$ROOT/roster-c1.txt" "$ROOT/roster-c2.txt" >> "$LOG" 2>&1 || true

say ""
say "=== PRE-INTERSECT vs POST-INTERSECT (per name+status) ==="
diff "$ROOT/roster-pre-intersect-mine.txt" "$ROOT/roster-post-intersect-mine.txt" >> "$LOG" 2>&1 || true
say "--- post-only lines ---"; comm -13 "$ROOT/roster-pre-intersect-mine.txt" "$ROOT/roster-post-intersect-mine.txt" >> "$LOG" 2>&1 || true
say "--- pre-only lines ---"; comm -23 "$ROOT/roster-pre-intersect-mine.txt" "$ROOT/roster-post-intersect-mine.txt" >> "$LOG" 2>&1 || true

say ""
say "=== gates on final tree ==="
GOFLAGS= go build ./... > "$ROOT/logs/gate-build.txt" 2>&1; say "go_build_rc=$?"
GOPATHBIN="$(go env GOPATH)/bin/gofumpt.exe"
"$GOPATHBIN" -l cmd/wisp > "$ROOT/logs/gate-gofumpt.txt" 2>&1; say "gofumpt_l_cmdwisp_rc=$? lines=$(wc -l < "$ROOT/logs/gate-gofumpt.txt")"
say "gofumpt_names:"; sed 's/^/  /' "$ROOT/logs/gate-gofumpt.txt" >> "$LOG"
sh scripts/d22scan.sh > "$ROOT/logs/gate-d22scan.txt" 2>&1; say "d22scan_rc=$?"
say "d22scan_tail:"; tail -6 "$ROOT/logs/gate-d22scan.txt" | sed 's/^/  /' >> "$LOG"

say ""
say "rc=0 done TS=$(date '+%Y-%m-%d %H:%M:%S %z')"
