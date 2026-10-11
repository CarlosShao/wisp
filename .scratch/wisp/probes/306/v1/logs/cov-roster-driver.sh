#!/usr/bin/env bash
# 306-v1 coverage + full-package roster driver. Runs ONLY inside the out-of-repo export tree.
# Usage: bash logs/cov-roster-driver.sh <EXPORT_TREE>
set -u
TREE="$1"
REPO="D:/work/workspace/projects plans/Wisp"
LOG="$REPO/.scratch/wisp/probes/306/v1/logs/sens"
mkdir -p "$LOG"
INJ="$TREE/internal/audio/wavinjector.go"
T306="$TREE/internal/audio/wavinjector_extensible_float_306_test.go"
PRIST="$TREE/_306v1-pristine"
HIDE="$TREE/_306v1-hidden"
mkdir -p "$HIDE"
SUM="$LOG/ZZ-cov-summary.txt"
export GOFLAGS=-mod=mod
BLOCKS='wavinjector\.go:(18[5-9]|19[0-9]|2[0-2][0-9])\.'

restore() { cp "$PRIST/wavinjector.go" "$INJ"; cp "$PRIST/test.go" "$T306"; }
m192() { perl -i -pe 's/if fmtTag == 0xFFFE \{/if fmtTag == 0xFFFD {/' "$INJ"; }
hidefix() { mv "$T306" "$HIDE/wavinjector_extensible_float_306_test.go.hidden"; }
showfix() { [ -f "$HIDE/wavinjector_extensible_float_306_test.go.hidden" ] && mv "$HIDE/wavinjector_extensible_float_306_test.go.hidden" "$T306"; }

echo "=== cov/roster driver start $(date) ===" > "$SUM"
{ echo "tree_head_tree=$(ls "$TREE" | wc -l)"; echo "go_version=$(go version)"; } >> "$SUM"

# ---------- C1 coverage: HEAD tree, FULL package (fixture present) ----------
restore; hidefix 2>/dev/null; showfix
(cd "$TREE" && go test ./internal/audio/ -count=1 -coverprofile="$LOG/c1-head-fullpkg.cover") >"$LOG/C1-head-fullpkg.txt" 2>&1
echo "C1 rc=$? coverage_total=$(grep -c . "$LOG/c1-head-fullpkg.cover")" >> "$SUM"
grep -E "$BLOCKS" "$LOG/c1-head-fullpkg.cover" | sed 's#.*wavinjector.go:#  C1 block #' >> "$SUM"

# ---------- C2 coverage: fixture moved away, FULL package (the "before" world) ----------
restore; hidefix
(cd "$TREE" && go test ./internal/audio/ -count=1 -coverprofile="$LOG/c2-nofixture-fullpkg.cover") >"$LOG/C2-nofixture-fullpkg.txt" 2>&1
echo "C2 rc=$?  topPASS_in_C2=$(grep -c '^ok' "$LOG/C2-nofixture-fullpkg.txt")" >> "$SUM"
grep -E "$BLOCKS" "$LOG/c2-nofixture-fullpkg.cover" | sed 's#.*wavinjector.go:#  C2 block #' >> "$SUM"
showfix

# ---------- C3 coverage: ONLY the two new tests ----------
restore
(cd "$TREE" && go test ./internal/audio/ -count=1 -run ExtensibleFloat32306 -coverprofile="$LOG/c3-newtests-only.cover") >"$LOG/C3-newtests-only.txt" 2>&1
echo "C3 rc=$?" >> "$SUM"
grep -E "$BLOCKS" "$LOG/c3-newtests-only.cover" | sed 's#.*wavinjector.go:#  C3 block #' >> "$SUM"

# ---------- B1/B2 full package verbose roster, two shots, HEAD tree ----------
restore
(cd "$TREE" && go test ./internal/audio/ -count=1 -v) >"$LOG/B1-fullpkg-shot1.txt" 2>&1
echo "B1 rc=$? topFAIL=$(grep -c '^--- FAIL' "$LOG/B1-fullpkg-shot1.txt") topPASS=$(grep -c '^--- PASS' "$LOG/B1-fullpkg-shot1.txt") topSKIP=$(grep -c '^--- SKIP' "$LOG/B1-fullpkg-shot1.txt") subPASS=$(grep -c '^    --- PASS' "$LOG/B1-fullpkg-shot1.txt") subFAIL=$(grep -c '^    --- FAIL' "$LOG/B1-fullpkg-shot1.txt") okline=$(grep -c '^ok' "$LOG/B1-fullpkg-shot1.txt")" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP|FAIL|ok )' "$LOG/B1-fullpkg-shot1.txt" | sed 's/^/  [B1] /' >> "$SUM"
restore
(cd "$TREE" && go test ./internal/audio/ -count=1 -v) >"$LOG/B2-fullpkg-shot2.txt" 2>&1
echo "B2 rc=$? topFAIL=$(grep -c '^--- FAIL' "$LOG/B2-fullpkg-shot2.txt") topPASS=$(grep -c '^--- PASS' "$LOG/B2-fullpkg-shot2.txt") topSKIP=$(grep -c '^--- SKIP' "$LOG/B2-fullpkg-shot2.txt") subPASS=$(grep -c '^    --- PASS' "$LOG/B2-fullpkg-shot2.txt") subFAIL=$(grep -c '^    --- FAIL' "$LOG/B2-fullpkg-shot2.txt")" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP|FAIL|ok )' "$LOG/B2-fullpkg-shot2.txt" | sed 's/^/  [B2] /' >> "$SUM"
echo "B1vsB2 roster diff (FAIL+SKIP names):" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP)' "$LOG/B1-fullpkg-shot1.txt" | sed 's/ (.*//' | sort > "$LOG/B1-roster.txt"
grep -E '^(--- FAIL|--- SKIP)' "$LOG/B2-fullpkg-shot2.txt" | sed 's/ (.*//' | sort > "$LOG/B2-roster.txt"
diff "$LOG/B1-roster.txt" "$LOG/B2-roster.txt" | sed 's/^/  diff /' >> "$SUM"
echo "rc_B1_B2_roster_diff=$?" >> "$SUM"

# ---------- B3 full package with the new fixture MOVED AWAY (the "before" green) ----------
restore; hidefix
(cd "$TREE" && go test ./internal/audio/ -count=1 -v) >"$LOG/B3-nofixture-fullpkg-v.txt" 2>&1
echo "B3 rc=$? topFAIL=$(grep -c '^--- FAIL' "$LOG/B3-nofixture-fullpkg-v.txt") topPASS=$(grep -c '^--- PASS' "$LOG/B3-nofixture-fullpkg-v.txt") topSKIP=$(grep -c '^--- SKIP' "$LOG/B3-nofixture-fullpkg-v.txt")" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP|FAIL|ok )' "$LOG/B3-nofixture-fullpkg-v.txt" | sed 's/^/  [B3] /' >> "$SUM"
showfix

# ---------- B4 full package, fixture present, :192 mutated (attributing the red roster) ----------
restore; m192
(cd "$TREE" && go test ./internal/audio/ -count=1 -v) >"$LOG/B4-mut192-fullpkg-v.txt" 2>&1
echo "B4 rc=$? topFAIL=$(grep -c '^--- FAIL' "$LOG/B4-mut192-fullpkg-v.txt") topPASS=$(grep -c '^--- PASS' "$LOG/B4-mut192-fullpkg-v.txt") topSKIP=$(grep -c '^--- SKIP' "$LOG/B4-mut192-fullpkg-v.txt")" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP)' "$LOG/B4-mut192-fullpkg-v.txt" | sed 's/^/  [B4] /' >> "$SUM"
grep -E 'unsupported wav format' "$LOG/B4-mut192-fullpkg-v.txt" | head -4 | sed 's/^/  [B4 msg] /' >> "$SUM"
showfix

# ---------- B5 full package, fixture moved away, :192 mutated (the ticket's rc=0 baseline) ----------
restore; hidefix; m192
(cd "$TREE" && go test ./internal/audio/ -count=1 -v) >"$LOG/B5-nofixture-mut192-fullpkg.txt" 2>&1
echo "B5 rc=$? topFAIL=$(grep -c '^--- FAIL' "$LOG/B5-nofixture-mut192-fullpkg.txt") topPASS=$(grep -c '^--- PASS' "$LOG/B5-nofixture-mut192-fullpkg.txt") topSKIP=$(grep -c '^--- SKIP' "$LOG/B5-nofixture-mut192-fullpkg.txt')" >> "$SUM"
grep -E '^(--- FAIL|--- SKIP|FAIL|ok )' "$LOG/B5-nofixture-mut192-fullpkg.txt" | sed 's/^/  [B5] /' >> "$SUM"

restore
cmp -s "$INJ" "$PRIST/wavinjector.go"; echo "rc_final_cmp_inj=$?" >> "$SUM"
cmp -s "$T306" "$PRIST/test.go"; echo "rc_final_cmp_test=$?" >> "$SUM"
echo "=== cov/roster driver done $(date) ===" >> "$SUM"
