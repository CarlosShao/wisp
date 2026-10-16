#!/bin/sh
# 322-v1 (acceptance leg, read-only): Windows teeth-out + reverse-shape mutations, run
# PHYSICALLY inside a git-archive export tree outside the repo. Never -overlay, never the
# shared worktree (a writer leg is live in cmd/wisp this wave).
# Rules copied from recorded accidents, not preferences: sha256 before/after, .bak kept,
# trees never removed, rc from $? directly, one anchor uniqueness check before any edit,
# GOTMPDIR native backslash form on the C: volume, sherpa dll dir on PATH.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322v1
LOGS="$REPO/.scratch/wisp/probes/322/v1/logs"
mkdir -p "$S" "$LOGS"
export GOTMPDIR='C:\Users\swq\AppData\Local\Temp\gotmp322v1'
mkdir -p "C:/Users/swq/AppData/Local/Temp/gotmp322v1"
export PATH="$REPO/third_party/sherpa-onnx:$PATH"
PAT='TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1|TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'
: >"$LOGS/mut-win-v1-summary.txt"
say() { echo "$@" | tee -a "$LOGS/mut-win-v1-summary.txt"; }
say "### 322-v1 Windows teeth run, $(date '+%F %T %z') HEAD=$(git -C "$REPO" rev-parse HEAD)"

T=$(mktemp -d "$S/win.XXXXXX")
( cd "$REPO" && git archive HEAD | tar -x -C "$T" )
say "TREE $T  .go=$(find "$T" -name '*.go' | wc -l | tr -d ' ') (counted before any number is trusted)"
W=internal/winsec/winsec_windows.go
BEFORE=$(sha256sum "$T/$W" | cut -d' ' -f1)
say "WIN-BEFORE sha=$BEFORE"
say "ANCHOR hits of the guard func line = $(grep -c '^func preOpenLeafGuard(path string) error {' "$T/$W")"
[ "$(grep -c '^func preOpenLeafGuard(path string) error {' "$T/$W")" = "1" ] || { say "ANCHOR NOT UNIQUE - ABORT"; exit 3; }
cp "$T/$W" "$T/$W.orig"

run() { # <label>
	label=$1; rc=0
	( cd "$T" && go test ./internal/winsec/ -count=1 -v -run "$PAT" ) >"$LOGS/$label.log" 2>&1 || rc=$?
	printf '=== %s rc=%s RUN=%s topPASS=%s topFAIL=%s topSKIP=%s\n' "$label" "$rc" \
		"$(grep -c '^=== RUN' "$LOGS/$label.log" || true)" "$(grep -c '^--- PASS' "$LOGS/$label.log" || true)" \
		"$(grep -c '^--- FAIL' "$LOGS/$label.log" || true)" "$(grep -c '^--- SKIP' "$LOGS/$label.log" || true)" | tee -a "$LOGS/mut-win-v1-summary.txt"
	say "tooling words: walk=$(grep -ci walk "$LOGS/$label.log" || true) buildfailed=$(grep -ci 'build failed' "$LOGS/$label.log" || true) noteststorun=$(grep -ci 'no tests to run' "$LOGS/$label.log" || true) setupfailed=$(grep -ci 'setup failed' "$LOGS/$label.log" || true) rc0run=$(grep -c '^RUN=0' "$LOGS/$label.log" || true)"
	grep -h 'AC#3 local paired reading' "$LOGS/$label.log" | tee -a "$LOGS/mut-win-v1-summary.txt" || say "NO PAIRED ROW in $label"
	say "--- RED lines verbatim ($label):"
	grep -h 'RED' "$LOGS/$label.log" | tee -a "$LOGS/mut-win-v1-summary.txt" || say "(no RED line)"
}

# ---------- w0: unmutated export tree baseline (must be green) ----------
run w0-baseline

# ---------- w1: teeth OUT -- guard always answers "not a link" ----------
awk '{print} /^func preOpenLeafGuard\(path string\) error \{$/{print "\treturn nil // MUTANT w1: guard always answers OK"}' "$T/$W.orig" > "$T/$W"
say "W1-AFTER line: $(grep -n 'MUTANT w1' "$T/$W")"
run w1-teethout
cp "$T/$W.orig" "$T/$W"
AFTER=$(sha256sum "$T/$W" | cut -d' ' -f1)
say "W1-RESTORE sha=$AFTER (must equal WIN-BEFORE $BEFORE)"
run w1-restore

# ---------- w2: REVERSE shape -- guard refuses EVERY name ----------
awk '{print} /^func preOpenLeafGuard\(path string\) error \{$/{print "\treturn preOpenRefusal(path, \"MUTANT w2: blanket refusal\") // MUTANT w2"}' "$T/$W.orig" > "$T/$W"
say "W2-AFTER line: $(grep -n 'MUTANT w2' "$T/$W")"
run w2-blanket
cp "$T/$W.orig" "$T/$W"
AFTER2=$(sha256sum "$T/$W" | cut -d' ' -f1)
say "W2-RESTORE sha=$AFTER2 (must equal WIN-BEFORE $BEFORE)"
run w2-restore
say "### done $(date '+%F %T %z') TREE kept: $T"
