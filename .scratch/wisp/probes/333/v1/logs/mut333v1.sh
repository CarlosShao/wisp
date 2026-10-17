#!/bin/sh
# 333-v1 mutation runs (teeth). All mutations land on the OUT-OF-REPO export tree
# M (git archive 859026d2); tracked files are never touched. Every shot restores in a
# finally-style trap and prints the restore sha256 next to the PRE sha256.
set -u
LOGD="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/333/v1/logs"
A=/d/W333V1-SuqCO8
M=$(cat /d/333v1-tree-m.txt)
GO=internal/panel/display_scrub_216.go
TT=internal/panel/display_scrub_327r2_test.go
RUNSET='TestDisplayWashRoster|TestDisplayWashJudgeAndWriter|TestDisplayWashFastPath|TestDisplayWashEscapes|TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut'
export TMP="D:/wisp-333v1-tmp/tmp" TEMP="D:/wisp-333v1-tmp/tmp" GOTMPDIR="D:/wisp-333v1-tmp/gotmp"
mkdir -p "$TMP" "$GOTMPDIR"

# the anchor's 15-rune switch body, lifted verbatim from the pre-anchor tree P
ROSTER_BODY=$(awk '/^func isDisplayWashCf/,/^}$/' /d/W333V2-hLzM9f/$GO | sed '1d;$d')

pre_h() { sha256sum "$M/$GO" "$M/$TT" | cut -c1-16 | tr '\n' ' '; }
restore() {
  cp "$A/$GO" "$M/$GO"; cp "$A/$TT" "$M/$TT"
  echo "RESTORE_HASHES $(sha256sum "$M/$GO" "$M/$TT" | cut -c1-16 | tr '\n' ' ')"
}
shot() { # $1 label
  L="$1"; OUT="$LOGD/mut-$L.log"
  { echo "# 333-v1 mutation shot $L  tree=$M  anchor=$(cat $M/.anchor-full)";
    echo "# PRE  $PRE"; echo "# -run $RUNSET"; } > "$OUT"
  ( cd "$M" && go test -v -count=1 -timeout 30m -run "$RUNSET" ./internal/panel/ >> "$OUT" 2>&1; echo "rc=$?" >> "$OUT" )
  RUN=$(grep -c '^=== RUN' "$OUT" || true); PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true); SK=$(grep -c '^--- SKIP' "$OUT" || true)
  echo "FOUR $L: RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK"
  grep '^--- FAIL' "$OUT" | sed 's/^--- FAIL *//' | awk '{print $1}' | sort > "$LOGD/mut-$L.failnames.txt"
  restore >> "$OUT"
  echo "POST $(pre_h)" >> "$OUT"
}

cd "$M" || exit 1
cp "$A/$GO" "$M/$GO"; cp "$A/$TT" "$M/$TT"
PRE=$(pre_h); echo "PRE-BASELINE $PRE"

echo "=== M0 pristine control (no mutation) ==="
shot M0-control

echo "=== M1 product predicate back to the 15-rune switch (test untouched) ==="
perl -0pi -e 's/\treturn unicode\.Is\(unicode\.Cf, r\)\n/'"$ROSTER_BODY"'\n/' "$M/$GO"
echo "M1 body now:"; awk '/^func isDisplayWashCf/,/^}$/' "$M/$GO" | head -14
echo "MUT $(pre_h)"; shot M1-roster

echo "=== M2 fast path short-circuited to false ==="
cp "$A/$GO" "$M/$GO"
perl -0pi -e 's/(func needsDisplayWash\(s string\) bool \{\n)/$1\treturn false \/\/ MUTANT 333v1 M2\n/' "$M/$GO"
echo "MUT $(pre_h)"; shot M2-fastpath

echo "=== M3 self-copy reverse form: test actual side copied from expected + product back to roster ==="
cp "$A/$GO" "$M/$GO"
perl -0pi -e 's/\treturn unicode\.Is\(unicode\.Cf, r\)\n/'"$ROSTER_BODY"'\n/' "$M/$GO"
perl -0pi -e 's/washSetOver\(isDisplayWashCf\)/washSetOver(func(r rune) bool { return unicode.Is(unicode.Cf, r) })/' "$M/$TT"
echo "MUT $(pre_h)"; shot M3-selfcopy

echo "=== M4 self-copy reverse form alone: test copied, product left CORRECT ==="
cp "$A/$TT" "$M/$TT"
perl -0pi -e 's/washSetOver\(isDisplayWashCf\)/washSetOver(func(r rune) bool { return unicode.Is(unicode.Cf, r) })/' "$M/$TT"
echo "MUT $(pre_h)"; shot M4-selfcopy-only

echo "=== final byte-identity check vs pristine tree A ==="
diff -r "$A/internal/panel" "$M/internal/panel" && echo "TREES IDENTICAL AFTER RESTORE"
echo "sha A=$(sha256sum "$A/$GO" | cut -c1-16) sha M=$(sha256sum "$M/$GO" | cut -c1-16)"
