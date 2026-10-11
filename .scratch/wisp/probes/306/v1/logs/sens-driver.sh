#!/usr/bin/env bash
# 306-v1 sensitivity driver. Runs ONLY inside an out-of-repo export tree (never the shared worktree).
# Usage: bash logs/sens-driver.sh <EXPORT_TREE>
set -u
TREE="$1"
REPO="D:/work/workspace/projects plans/Wisp"
LOG="$REPO/.scratch/wisp/probes/306/v1/logs/sens"
mkdir -p "$LOG"
INJ="$TREE/internal/audio/wavinjector.go"
T306="$TREE/internal/audio/wavinjector_extensible_float_306_test.go"
PRIST="$TREE/_306v1-pristine"
mkdir -p "$PRIST"
cp "$INJ" "$PRIST/wavinjector.go"
cp "$T306" "$PRIST/test.go"
SUM="$LOG/ZZ-summary.txt"
export GOFLAGS=-mod=mod

restore() { cp "$PRIST/wavinjector.go" "$INJ"; cp "$PRIST/test.go" "$T306"; }

m192() { perl -i -pe 's/if fmtTag == 0xFFFE \{/if fmtTag == 0xFFFD {/' "$INJ"; }
m218() { perl -i -pe 's/case fmtTag == 3 && bits == 32:/case fmtTag == 2 && bits == 32:/' "$INJ"; }
m196b() { perl -i -pe 's/data\[body\+24 : body\+26\]/data[body+26 : body+28]/' "$INJ"; }
m196l() { perl -i -pe 's/data\[body\+24 : body\+26\]/data[body+26 : body+26]/' "$INJ"; }

# run <name> <filter>
run() {
  local name="$1" filt="$2" out="$LOG/$1.txt"
  (cd "$TREE" && go test ./internal/audio/ -count=1 -vet=off -run "$filt") >"$out" 2>&1
  local rc=$?
  echo "SCENARIO=$name rc=$rc topFAIL=$(grep -c '^--- FAIL' "$out") topPASS=$(grep -c '^--- PASS' "$out") topSKIP=$(grep -c '^--- SKIP' "$out") paniclines=$(grep -c '^panic:' "$out")" >> "$SUM"
  grep -E '^(--- FAIL|--- PASS|panic:|ok |FAIL)' "$out" | sed "s/^/  [$name] /" >> "$SUM"
}

echo "=== 306-v1 driver start $(date) TREE=$TREE ===" > "$SUM"

# ---------- Group P: pristine test file, mutated production ----------
restore; run P0-pristine ExtensibleFloat32306
restore; m192; echo "P1 line192:$(sed -n '192p' "$INJ" | sed 's/^[[:space:]]*//')" >> "$SUM"; run P1-mut192 ExtensibleFloat32306
restore; m218; echo "P2 line218:$(sed -n '218p' "$INJ" | sed 's/^[[:space:]]*//')" >> "$SUM"; run P2-mut218 ExtensibleFloat32306
restore; m196b; echo "P3 line196:$(sed -n '196p' "$INJ" | sed 's/^[[:space:]]*//')" >> "$SUM"; run P3-mut196-boundsshift ExtensibleFloat32306
restore; m196l; echo "P4 line196:$(sed -n '196p' "$INJ" | sed 's/^[[:space:]]*//')" >> "$SUM"; run P4-mut196-literalwording ExtensibleFloat32306

# ---------- Group VA: value expectations made tautological (err gates stay live) ----------
mkVA() {
  cp "$PRIST/test.go" "$T306"
  perl -i -pe 's/\tif rate != 48000 \{/\tif false && rate != 48000 {/; s/\tif chans != 2 \{/\tif false && chans != 2 {/; s/\tif len\(samples\) != len\(w306WantI16\) \{/\tif false \&\& len(samples) != len(w306WantI16) {/; s/\t\tif samples\[i\] != want \{/\t\tif false && samples[i] != want {/; s/\tif len\(inj\.samples\) != len\(w306WantI16\) \{/\tif false \&\& len(inj.samples) != len(w306WantI16) {/; s/\t\tif inj\.samples\[i\] != want \{/\t\tif false && inj.samples[i] != want {/' "$T306"
}
restore; mkVA; cp "$T306" "$PRIST/VA.go"; echo "VA_inert_count=$(grep -c 'if false &&' "$T306")" >> "$SUM"
run VA-pristine ExtensibleFloat32306
restore; mkVA; m192; run VA-mut192 ExtensibleFloat32306
restore; mkVA; m218; run VA-mut218 ExtensibleFloat32306
restore; mkVA; m196b; run VA-mut196boundsshift ExtensibleFloat32306
restore; mkVA; m196l; run VA-mut196literal ExtensibleFloat32306

# ---------- Group VR: err gate inert in Test 1 only; value assertions live ----------
mkVR() {
  cp "$PRIST/test.go" "$T306"
  perl -i -pe 's/\tif err != nil \{/\tif false && err != nil {/ if $. == 176' "$T306"
}
restore; mkVR; cp "$T306" "$PRIST/VR.go"; echo "VR_lines=$(grep -n 'if false && err != nil' "$T306" | tr '\n' ';')" >> "$SUM"
run VR-pristine TestParseWavExtensibleFloat32306
restore; mkVR; m192; run VR-mut192 TestParseWavExtensibleFloat32306
restore; mkVR; m218; run VR-mut218 TestParseWavExtensibleFloat32306
restore; mkVR; m196b; run VR-mut196boundsshift TestParseWavExtensibleFloat32306
restore; mkVR; m196l; run VR-mut196literal TestParseWavExtensibleFloat32306

# ---------- Group FLIP: reverse form, one assertion at a time, pristine production ----------
for ln in 171 176 181 184 187 192 212 217 221; do
  restore
  perl -i -pe "s/!=/==/ if \$. == $ln" "$T306"
  echo "FLIP line$ln => $(sed -n "${ln}p" "$T306" | sed 's/^[[:space:]]*//')" >> "$SUM"
  run "FLIP-line$ln" ExtensibleFloat32306
done
restore; perl -i -pe 's/!=/==/ if $. >= 169 && $. != 207' "$T306"
run FLIP-all-from-169 ExtensibleFloat32306

restore
cmp -s "$INJ" "$PRIST/wavinjector.go"; echo "rc_final_cmp_inj=$?" >> "$SUM"
cmp -s "$T306" "$PRIST/test.go"; echo "rc_final_cmp_test=$?" >> "$SUM"
echo "=== driver done $(date) ===" >> "$SUM"
