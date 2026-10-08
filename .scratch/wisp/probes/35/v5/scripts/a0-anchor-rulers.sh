#!/bin/bash
# 35-v5 opening anchors. Independent re-measurement of everything the orchestrator restated.
# Nothing here is transcribed: every number is produced by the command next to it.
# Every file ends with its own rc= line. No tracked file is written by this script.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/35/v5/logs"
cd "$REPO" || exit 1
FIX=cmd/wisp/panel_transport_35r2_test.go
LIVE=cmd/wisp/panel_transport_live_35v2_windows_test.go
PROD=cmd/wisp/panel_host_windows.go
TICKET=.scratch/wisp/issues/35-panel-bridge-c17.md
mkdir -p "$LOGS"

{ echo "=== identity ($(date '+%F %T%z')) ==="
  echo "HEAD=$(git log -1 --format=%H)"
  echo "branch=$(git rev-parse --abbrev-ref HEAD)"
  echo "write-leg commits since daf1f5a0: $(git log --oneline daf1f5a0..HEAD | wc -l)"
  git log --format='%h %s' daf1f5a0..HEAD | cut -c1-90
  echo "rc=0"; } > "$LOGS/a01-identity.txt" 2>&1

{ echo "=== ticket AC#8 box census (must be 8 unchecked / 2 checked, both before and after this leg) ==="
  echo "unchecked=$(grep -cE '^[[:space:]]*- \[ \]' "$TICKET")"
  echo "checked-strict=$(grep -cE '^[[:space:]]*- \[x]' "$TICKET")"
  echo "--- AC#8 body line range ---"
  grep -n 'Behavioural yard: the third face' "$TICKET"
  echo "rc=0"; } > "$LOGS/a02-boxes.txt" 2>&1

{ echo "=== (ii) prefix census on the CURRENT fixture (re-measured, not restated) ==="
  echo "line-count fixture=$(wc -l < "$FIX")"
  for p in 'M-B RED' 'M-B1' 'M-B2' 'M-B3' 'M-A RED'; do
    echo "'$p'  lines=$(grep -c "$p" "$FIX")  occurrences=$(grep -o "$p" "$FIX" | wc -l)"
  done
  echo "--- the trailing-space ruler form the write leg used ---"
  echo "M-B1space=$(grep -c 'M-B1 ' "$FIX") M-B2space=$(grep -c 'M-B2 ' "$FIX") M-B3space=$(grep -c 'M-B3 ' "$FIX")"
  echo "=== same ruler on the daf1f5a0 (pre-35-r5) blob ==="
  git show "daf1f5a0:$FIX" > "$LOGS/before-fixture.copy"
  for p in 'M-B RED' 'M-B1' 'M-B2' 'M-B3' 'M-A RED'; do
    echo "'$p'  lines=$(grep -c "$p" "$LOGS/before-fixture.copy")"
  done
  echo "before line-count=$(wc -l < "$LOGS/before-fixture.copy")"
  echo "=== assertion-strength ruler ==="
  echo "t.Fatalf before=$(grep -c 't\.Fatalf' "$LOGS/before-fixture.copy") after=$(grep -c 't\.Fatalf' "$FIX") delta=$(( $(grep -c 't\.Fatalf' "$FIX") - $(grep -c 't\.Fatalf' "$LOGS/before-fixture.copy") ))"
  echo "t.Errorf before=$(grep -c 't\.Errorf' "$LOGS/before-fixture.copy") after=$(grep -c 't\.Errorf' "$FIX")"
  echo "t.Logf before=$(grep -c 't\.Logf' "$LOGS/before-fixture.copy") after=$(grep -c 't\.Logf' "$FIX")"
  echo "func Test before=$(grep -c '^func Test' "$LOGS/before-fixture.copy") after=$(grep -c '^func Test' "$FIX")"
  echo "rc=0"; } > "$LOGS/a03-prefix-census.txt" 2>&1

{ echo "=== (i) the new case: exact line anchors ==="
  grep -n 'func TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable' "$FIX"
  echo "--- first assertion (red sentence) line ---"
  grep -n 'M-B3 RED (guard reduced to a truthiness test)' "$FIX"
  echo "--- 台件 (the world-staging line) ---"
  grep -n 'bf.window.set(panelDispatchBinding, "not-a-function")' "$FIX"
  echo "--- the production guard line 680 verbatim ---"
  sed -n '680p' "$PROD"
  echo "prod blob worktree=$(git hash-object "$PROD") HEAD=$(git rev-parse "HEAD:$PROD") daf1f5a0=$(git rev-parse "daf1f5a0:$PROD")"
  echo "zero-prod-change numstat (expect NO output):"
  git diff --numstat daf1f5a0..HEAD -- "$PROD"; echo "numstat-rc=$?"
  echo "rc=0"; } > "$LOGS/a04-newcase-lines.txt" 2>&1

{ echo "=== (iv) the 〔仅本机可量〕 label ==="
  echo "grep -c '仅本机可量' on $LIVE = $(grep -c '仅本机可量' "$LIVE")"
  echo "file length=$(wc -l < "$LIVE")"
  echo "--- is every added header line a comment? (non-comment lines in 1..25, expect 1 = the //go:build line at most) ---"
  sed -n '1,30p' "$LIVE" | grep -vn '^\s*\(//\|$\)' || echo "(all comment/blank in 1..30)"
  echo "--- four elements + revocation ruler ---"
  for k in '仅本机可量' 'WHAT IS MACHINE-LOCAL' 'WHO OWNS IT' 'RE-RUN CADENCE' 'PRICE PAID' 'Revocation' 'ORCHESTRATOR'; do
    echo "'$k' -> $(grep -c "$k" "$LIVE") hit(s)"
  done
  echo "rc=0"; } > "$LOGS/a05-iv-label.txt" 2>&1

{ echo "=== Q3 side-check: is any LIVE instrument grepping the literal M-B-plus-RED? ==="
  echo "--- (1) the .sh/.py/.ps1/.go instruments in the repo that mention it ---"
  PAT="M-B"" RED"
  git grep -l "$PAT" -- '*.sh' '*.py' '*.ps1' '*.go' '*.yml' 2>/dev/null || echo "(none)"
  echo "--- (2) every .sh instrument that mentions it, with its lines ---"
  for f in $(git grep -l "$PAT" -- '*.sh' 2>/dev/null); do
    echo "## $f"
    grep -n "$PAT" "$f"
  done
  echo "--- (3) do CI / tools / scripts reference any red-sentence prefix at all? ---"
  git grep -n 'M-B\|M-A RED' -- .github tools scripts 2>/dev/null || echo "(no CI/tools/scripts hit)"
  echo "--- (4) is any .sh grepping the WORKING-COPY fixture rather than a saved blob copy? ---"
  for f in $(git grep -l "$PAT" -- '*.sh' 2>/dev/null); do
    echo "## $f : $(grep -c "grep.*\$FIX\b\|grep.*panel_transport_35r2_test.go" "$f") working-copy read(s), $(grep -c "before-35r2-fixture.copy\|git show" "$f") blob-anchored read(s)"
  done
  echo "rc=0"; } > "$LOGS/a06-live-instruments.txt" 2>&1

{ echo "=== roster the write leg claims 13 PASS: which top-level tests match -run ==="
  RUN='TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope'
  echo "--- matching top-level funcs across the whole package's NON-windows test files ---"
  grep -rnE "^func ($RUN)" --include='*_test.go' cmd/wisp | grep -v '_windows_test.go' | sed 's/(t \*testing.T).*(//'
  echo "count=$(grep -rhE "^func ($RUN)" --include='*_test.go' cmd/wisp | wc -l)"
  echo "rc=0"; } > "$LOGS/a07-run-roster.txt" 2>&1

echo "--- summary ---"
for f in a01-identity a02-boxes a03-prefix-census a04-newcase-lines a05-iv-label a06-live-instruments a07-run-roster; do
  echo "### $f"; cat "$LOGS/$f.txt"; echo
done
tail -1 "$LOGS/a01-identity.txt"
