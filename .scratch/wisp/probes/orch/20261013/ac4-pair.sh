#!/usr/bin/env bash
# Ticket 299 AC#4, the last missing half: the paired whole-package comparison.
#   before side = git archive 6985de48  (the commit right before the 299 fix landed)
#   after  side = the shared working tree at HEAD (44f5bce9)
# Two shots per side, rosters compared by NAME only (durations stripped).
# Out-of-repo trees under the 专用根; the real repo's third_party DLLs are put on PATH.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT="/c/Users/swq/AppData/Local/wisp-test-gotmp/orch299ac4-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$ROOT/logs" "$ROOT/before" || exit 9
LOG="$ROOT/logs/paired.txt"
say() { printf '%s\n' "$*" >> "$LOG"; }
cd "$REPO" || exit 9

say "orch-299-ac4 TS=$(date '+%Y-%m-%d %H:%M:%S %z') root=$ROOT"
say "after_head=$(git rev-parse HEAD)"
say "before_commit=6985de48 $(git rev-parse 6985de48)"

# --- before tree: out-of-repo export of 6985de48 ---
git archive 6985de48 go.mod go.sum internal cmd tools scripts frontend | tar -x -C "$ROOT/before" || exit 9
say "before_tree_go_files=$(find "$ROOT/before" -name '*.go' | wc -l)"
say "before_tree_six_pack=$(for d in go.mod go.sum internal cmd tools scripts frontend; do [ -e "$ROOT/before/$d" ] && printf '%s ' "$d"; done)"

run_shot() {
  tag="$1"; workdir="$2"; shift 2
  w="$ROOT/w-$tag"; mkdir -p "$w"
  say ""
  say "=== shot $tag ==="
  say "workdir=$workdir"
  say "command_line=GOTMPDIR=$w PATH=<repo>/third_party/sherpa-onnx:<repo>/build:\$PATH go test -count=1 -v ./cmd/wisp/ ./internal/panel/"
  ( cd "$workdir" && GOTMPDIR="$w" PATH="$REPO/third_party/sherpa-onnx:$REPO/build:$PATH" GOFLAGS= \
      go test -count=1 -v ./cmd/wisp/ ./internal/panel/ > "$ROOT/logs/full-$tag.txt" 2>&1 )
  rc=$?
  say "GO_TEST_RC=$rc"
  say "lines=$(wc -l < "$ROOT/logs/full-$tag.txt")"
  say "census=$(grep -c -- '--- PASS' "$ROOT/logs/full-$tag.txt" || true) PASS / $(grep -c -- '--- FAIL' "$ROOT/logs/full-$tag.txt" || true) FAIL / $(grep -c -- '--- SKIP' "$ROOT/logs/full-$tag.txt" || true) SKIP / $(grep -c -- '=== RUN' "$ROOT/logs/full-$tag.txt" || true) RUN"
  grep -E '^(ok|FAIL)' "$ROOT/logs/full-$tag.txt" | sed 's/^/  pkg: /' >> "$LOG"
  # roster: names only
  grep -oE '^[[:space:]]*--- (PASS|FAIL|SKIP): [^ ]+' "$ROOT/logs/full-$tag.txt" | sed 's/^[[:space:]]*//' | sed 's/: / /' | sort > "$ROOT/roster-$tag.txt"
  say "roster_lines=$(wc -l < "$ROOT/roster-$tag.txt")"
}

run_shot before-1 "$ROOT/before"
run_shot before-2 "$ROOT/before"
run_shot after-1  "$REPO"
run_shot after-2  "$REPO"

say ""
say "=== per-side intersection (>=2 shots each, keeps flakes out of BOTH sides) ==="
for side in before after; do
  comm -12 "$ROOT/roster-$side-1.txt" "$ROOT/roster-$side-2.txt" > "$ROOT/common-$side.txt"
  say "$side common_lines=$(wc -l < "$ROOT/common-$side.txt")"
done

say ""
say "=== reds on each side (from the intersection) ==="
grep '^FAIL ' "$ROOT/common-before.txt" | sed 's/^/  before: /' >> "$LOG"
grep '^FAIL ' "$ROOT/common-after.txt"  | sed 's/^/  after : /' >> "$LOG"

say ""
say "=== NEW reds = after_reds - before_reds (the AC#4 question) ==="
comm -13 <(grep '^FAIL ' "$ROOT/common-before.txt" | sort) <(grep '^FAIL ' "$ROOT/common-after.txt" | sort) | sed 's/^/  NEW_RED: /' >> "$LOG"
say "new_red_count=$(comm -13 <(grep '^FAIL ' "$ROOT/common-before.txt" | sort) <(grep '^FAIL ' "$ROOT/common-after.txt" | sort) | wc -l)"
say ""
say "=== gone reds (before minus after), for the record ==="
comm -23 <(grep '^FAIL ' "$ROOT/common-before.txt" | sort) <(grep '^FAIL ' "$ROOT/common-after.txt" | sort) | sed 's/^/  GONE: /' >> "$LOG"
say "rc=0"
