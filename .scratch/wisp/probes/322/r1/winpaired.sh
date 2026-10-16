#!/bin/sh
# 322-r1 Windows paired readings, run on THIS host (AC#3 is 〔仅本机可量〕: the container
# is linux and the CI runner's symlink-privilege tier has never been measured).
#
# Same rule as the container driver: two fresh snapshot trees from git archive HEAD, never
# deleted, one carrying the pre-ticket stub and one carrying the delivered bytes, and the
# same -run filter against both.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322r1
LOGS="$REPO/.scratch/wisp/probes/322/r1/logs"
mkdir -p "$LOGS"
export GOTMPDIR='C:/Users/swq/AppData/Local/wisp-test-gotmp/322r1gotmp'
mkdir -p "$GOTMPDIR"

PROD='internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go'
NEWFILES='internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'
PAT='TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1|TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1'

hostrun() {
	label=$1
	tree=$2
	rc=0
	(cd "$tree" && go test -count=1 -v -run "$PAT" ./internal/winsec/) >"$LOGS/$label.body.log" 2>&1 || rc=$?
	{
		echo "=== $label go-test rc=$rc (tree $tree)"
		printf 'RUN=%s\n' "$(grep -c '=== RUN' "$LOGS/$label.body.log" || true)"
		printf 'top-level (anchored ^---): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '^--- PASS' "$LOGS/$label.body.log" || true)" \
			"$(grep -c -- '^--- FAIL' "$LOGS/$label.body.log" || true)" \
			"$(grep -c -- '^--- SKIP' "$LOGS/$label.body.log" || true)"
		grep -h 'AC#3 local paired reading' "$LOGS/$label.body.log" || echo "NO AC#3 PAIRED ROW in $label"
	} | tee -a "$LOGS/win-summary.txt"
	grep -h -E "AC#3 RED" "$LOGS/$label.body.log" | head -8 | tee -a "$LOGS/win-summary.txt" || true
}

: >"$LOGS/win-summary.txt"
echo "### 322-r1 Windows paired readings, host $(go version)" | tee "$LOGS/win-summary.txt"

BASEW=$(mktemp -d "$S/basewin.XXXXXX")
(
	cd "$REPO" && git archive HEAD | tar -x -C "$BASEW"
)
for f in $NEWFILES; do
	cp "$REPO/$f" "$BASEW/$f"
done
cat >"$BASEW/internal/winsec/zz322r1_stub_pre_ticket_windows_only.go" <<'STUB'
//go:build windows

package winsec

// Pre-ticket Windows semantics: no leaf probe, oNoFollowPlatform contributes nothing.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB
printf 'TREE %s: .go=%s winsec/*.go=%s\n' "$BASEW" "$(find "$BASEW" -name '*.go' | wc -l | tr -d ' ')" "$(find "$BASEW/internal/winsec" -name '*.go' | wc -l | tr -d ' ')" | tee -a "$LOGS/win-summary.txt"
hostrun win-base "$BASEW"

FIXW=$(mktemp -d "$S/fixedwin.XXXXXX")
(
	cd "$REPO" && git archive HEAD | tar -x -C "$FIXW"
)
for f in $NEWFILES; do
	cp "$REPO/$f" "$FIXW/$f"
done
for f in $PROD; do
	cp "$REPO/$f" "$FIXW/$f"
done
printf 'TREE %s: .go=%s winsec/*.go=%s\n' "$FIXW" "$(find "$FIXW" -name '*.go' | wc -l | tr -d ' ')" "$(find "$FIXW/internal/winsec" -name '*.go' | wc -l | tr -d ' ')" | tee -a "$LOGS/win-summary.txt"
hostrun win-fixed "$FIXW"
echo "TREES kept: base=$BASEW fixed=$FIXW" | tee -a "$LOGS/win-summary.txt"
