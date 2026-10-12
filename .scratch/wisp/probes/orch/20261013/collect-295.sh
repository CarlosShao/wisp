#!/usr/bin/env bash
# Orchestrator's own re-run of 295-v1's three main dishes (one clean room, three shots).
# Ruler for every line: content anchors, never absolute line numbers.
# The leg's own file for reference: .scratch/wisp/probes/295/v1/logs/roster-diff-m1-p-m2.txt
set -u
REPO="D:/work/workspace/projects plans/Wisp"
SRC="cmd/wisp/resident_audio_windows.go"
ROOT="/c/Users/swq/AppData/Local/wisp-test-gotmp/orch295-$(date +%Y%m%d-%H%M%S)"
ANCHOR=$'\tra.started.Store(gate.Open())'
cd "$REPO" || exit 9
mkdir -p "$ROOT/logs" || exit 9
LOG="$ROOT/logs/orch-295-collect.txt"
say() { printf '%s\n' "$*" >> "$LOG"; }

say "orch-295-collect TS=$(date '+%Y-%m-%d %H:%M:%S %z') root=$ROOT"
say "head=$(git rev-parse HEAD)"
say "worktree_vs_head_for_src=$(git diff --numstat -- "$SRC" | wc -l) (0 = untouched)"

# --- pristine copy + two mutants, all out of the repo, all 425 lines ---
git show "HEAD:$SRC" > "$ROOT/pristine.go" || exit 9
cp "$ROOT/pristine.go" "$ROOT/m1.go"
cp "$ROOT/pristine.go" "$ROOT/m2.go"
python - "$ROOT" <<'PY'
import sys, io
root = sys.argv[1]
anchor = "\tra.started.Store(gate.Open())\n"
for name, repl in (("m1", "\t_ = gate.Open() // 295 MUT-M1 delete-write: the boot store stays, this path stops updating\n"),
                   ("m2", "\tra.started.Store(true) // 295 MUT-M2 tautology: always true regardless of the gate\n")):
    p = f"{root}/{name}.go"
    s = io.open(p, "r", encoding="utf-8", newline="").read()
    n = s.count(anchor)
    if n != 1:
        raise SystemExit(f"ANCHOR COUNT {n} in {p} - refusing")
    io.open(p, "w", encoding="utf-8", newline="").write(s.replace(anchor, repl))
    print(f"{name}: anchor replaced once")
PY
RC=$?; say "mutants_built_rc=$RC"; [ $RC -ne 0 ] && exit 9
for f in pristine m1 m2; do
  say "$f sha256=$(sha256sum "$ROOT/$f.go" | cut -d' ' -f1) lines=$(wc -l < "$ROOT/$f.go")"
done

# overlays (module path form)
ABS_SRC="D:/work/workspace/projects plans/Wisp/$SRC"
python - "$ROOT" "$ABS_SRC" <<'PY'
import sys, json, os
root, abs_src = sys.argv[1], sys.argv[2]
for name in ("m1", "m2"):
    ov = {"Replace": {abs_src: f"{root}/{name}.go"}}
    json.dump(ov, open(f"{root}/ov-{name}.json", "w"))
    print("wrote", f"{root}/ov-{name}.json")
PY

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
  say "package_lines="; grep -E '^(ok|FAIL|---)\s+github.com/CarlosShao/wisp' "$ROOT/logs/full-$tag.txt" | head -5 | sed 's/^/  /' >> "$LOG"
}

run_shot pristine
run_shot m1 "-overlay=$ROOT/ov-m1.json"
run_shot m2 "-overlay=$ROOT/ov-m2.json"

say ""
say "=== roster diff (name only, no durations; PASS/FAIL/SKIP) ==="
for t in pristine m1 m2; do
  grep -oE '^[[:space:]]*--- (PASS|FAIL|SKIP): [^ ]+' "$ROOT/logs/full-$t.txt" | sed 's/^[[:space:]]*//' | sed 's/: / /' | sort > "$ROOT/roster-$t.txt"
  say "$t roster_lines=$(wc -l < "$ROOT/roster-$t.txt")"
done
say ""
say "--- pristine vs m1 ---"; diff "$ROOT/roster-pristine.txt" "$ROOT/roster-m1.txt" >> "$LOG" 2>&1 || true
say "--- pristine vs m2 ---"; diff "$ROOT/roster-pristine.txt" "$ROOT/roster-m2.txt" >> "$LOG" 2>&1 || true
say ""
say "rc=0"
