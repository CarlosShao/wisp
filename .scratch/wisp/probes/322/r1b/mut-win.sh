#!/bin/sh
# 322-r1b: teeth check (去牙反控) for the Windows half of ticket 322, on this machine.
#
# Same rules as mut-posix.sh: git-archive export tree, PHYSICAL edits, .bak kept, sha256
# before/after, restore run must hash equal and go green. Nothing here touches the shared
# working tree, so no other leg can see a mutant.
#
# WHAT EACH SHAPE IS SUPPOSED TO DO (this leg does not decide it, it measures it):
#   w-fixed : HEAD bytes -> both 322r1 Windows cases green;
#   w1      : the probe body is made to answer "not a link" for every name (牙摘掉).
#             TestAC3WindowsLeafGuard... must go RED on all three link names AND
#             TestAC3WindowsPrivateFileRoute... must go RED on its instrument line.
#   w2      : the probe refuses EVERY name (反形). If a case stays green here while it was
#             green in fixed, that case does not measure the guard's JUDGMENT, only that
#             some error came back - which is the 恒真 shape the dispatch asks about.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322r1b
LOGS="$REPO/.scratch/wisp/probes/322/r1b/logs"
mkdir -p "$S" "$LOGS"
W=internal/winsec/winsec_windows.go
PAT='TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1|TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1'
: >"$LOGS/mut-win-summary.txt"
say() { echo "$@" | tee -a "$LOGS/mut-win-summary.txt"; }
say "### 322-r1b Windows teeth run, started $(date '+%F %T %z') HEAD=$(git -C "$REPO" rev-parse HEAD)"
say "host go: $(go version) / GOTMPDIR=${GOTMPDIR:-<unset>} / tmp volume=$(cygpath -w /tmp)"

newtree() {
	tree=$(mktemp -d "$S/$1.XXXXXX")
	( cd "$REPO" && git archive HEAD | tar -x -C "$tree" )
	say "TREE $tree: .go=$(find "$tree" -name '*.go' | wc -l | tr -d ' ') winsec/*.go=$(find "$tree/internal/winsec" -name '*.go' | wc -l | tr -d ' ')"
}

hostrun() { # <label> <tree> <runfilter>
	label=$1; tree=$2; pat=$3
	rc=0
	( cd "$tree" && go test ./internal/winsec/ -count=1 -v -run "$pat" ) >"$LOGS/$label.body.log" 2>&1 || rc=$?
	{
		echo "=== $label rc=$rc filter='$pat'"
		printf 'RUN=%s topPASS=%s topFAIL=%s topSKIP=%s\n' \
			"$(grep -c '=== RUN' "$LOGS/$label.body.log" || true)" \
			"$(grep -c '^--- PASS' "$LOGS/$label.body.log" || true)" \
			"$(grep -c '^--- FAIL' "$LOGS/$label.body.log" || true)" \
			"$(grep -c '^--- SKIP' "$LOGS/$label.body.log" || true)"
		grep -hE '^--- (FAIL|PASS|SKIP)' "$LOGS/$label.body.log" || true
		grep -h 'AC#3 local paired reading' "$LOGS/$label.body.log" || echo "NO AC#3 PAIRED ROW in $label"
		echo "--- RED lines verbatim:"
		grep -h 'RED' "$LOGS/$label.body.log" || echo "(no RED line)"
	} | tee -a "$LOGS/mut-win-summary.txt"
}

inject() { # <tree> <statement> <label> : one line added at the top of preOpenLeafGuard's body
	tree=$1; stmt=$2; label=$3
	n=$(grep -c '^func preOpenLeafGuard(path string) error {' "$tree/$W" || true)
	[ "$n" = "1" ] || { say "ANCHOR NOT UNIQUE for $label: $n hits - ABORT"; exit 3; }
	line=$(grep -n '^func preOpenLeafGuard(path string) error {' "$tree/$W" | cut -d: -f1)
	cp "$tree/$W" "$tree/$W.$label.bak"
	say "$label ANCHOR $W:$line; BEFORE sha=$(sha256sum "$tree/$W" | cut -d' ' -f1)"
	awk -v at="$line" -v s="$stmt" 'NR==at+1 {print "\t" s} {print}' "$tree/$W.$label.bak" > "$tree/$W"
	say "$label AFTER line $(($line + 1)): $(sed -n "$((line + 1))p" "$tree/$W")"
	say "$label AFTER sha=$(sha256sum "$tree/$W" | cut -d' ' -f1)  hits-of-injected=$(grep -c "$stmt" "$tree/$W" || true)"
}

restore() { # <tree> <label>
	cp "$1/$W.$2.bak" "$1/$W"
	say "$2 RESTORE sha=$(sha256sum "$1/$W" | cut -d' ' -f1) (must equal its BEFORE sha)"
}

# ---------- fixed ----------
FX=$(mktemp -d "$S/wfixed.XXXXXX"); ( cd "$REPO" && git archive HEAD | tar -x -C "$FX" )
say "TREE $FX: .go=$(find "$FX" -name '*.go' | wc -l | tr -d ' ') winsec/*.go=$(find "$FX/internal/winsec" -name '*.go' | wc -l | tr -d ' ')"
say "FIXED-BYTES $(sha256sum "$FX/$W" | cut -d' ' -f1)"
hostrun w-fixed "$FX" "$PAT"

# ---------- w1: probe answers "not a link" for every name (teeth out) ----------
T1=$(mktemp -d "$S/w1-teethout.XXXXXX"); ( cd "$REPO" && git archive HEAD | tar -x -C "$T1" )
say "TREE $T1: .go=$(find "$T1" -name '*.go' | wc -l | tr -d ' ') winsec/*.go=$(find "$T1/internal/winsec" -name '*.go' | wc -l | tr -d ' ')"
inject "$T1" 'return nil // MUTANT w1: guard removed' w1
hostrun w1 "$T1" "$PAT"
restore "$T1" w1
hostrun w1-restore "$T1" "$PAT"

# ---------- w2: probe refuses every name (reverse shape) ----------
T2=$(mktemp -d "$S/w2-blanket.XXXXXX"); ( cd "$REPO" && git archive HEAD | tar -x -C "$T2" )
say "TREE $T2: .go=$(find "$T2" -name '*.go' | wc -l | tr -d ' ') winsec/*.go=$(find "$T2/internal/winsec" -name '*.go' | wc -l | tr -d ' ')"
inject "$T2" 'return preOpenRefusal(path, "MUTANT w2: blanket refusal")' w2
hostrun w2 "$T2" "$PAT"
restore "$T2" w2
hostrun w2-restore "$T2" "$PAT"

say "### done $(date '+%F %T %z')"
say "TREES kept (never removed): fixed=$FX w1=$T1 w2=$T2"
