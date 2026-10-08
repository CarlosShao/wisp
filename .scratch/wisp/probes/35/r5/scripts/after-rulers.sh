#!/bin/bash
# 35-r5 after-rulers for cells (ii), (iii), (iv) plus the zero-production-change proof.
# Every before-reading is repeated from the daf1f5a0 blob (`git show`), so nothing here is
# transcribed from the dispatch; every file ends with its own rc= line.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/35/r5/logs"
cd "$REPO" || exit 1
FIX=cmd/wisp/panel_transport_35r2_test.go
LIVE=cmd/wisp/panel_transport_live_35v2_windows_test.go
PROD=cmd/wisp/panel_host_windows.go
git show "daf1f5a0:$FIX" > "$LOGS/before-35r2-fixture.copy"

{ echo "=== (ii) shared red-sentence prefix, before vs after ($(date '+%F %T%z')) ==="
  echo "BEFORE (daf1f5a0 blob): 'M-B RED'=$(grep -c 'M-B RED' "$LOGS/before-35r2-fixture.copy")  [ruler asked for 6]"
  echo "AFTER  (working copy) : 'M-B RED'=$(grep -c 'M-B RED' "$FIX")  [want 0]"
  echo "AFTER per-face token counts: M-B1=$(grep -c 'M-B1 ' "$FIX") M-B2=$(grep -c 'M-B2 ' "$FIX") M-B3=$(grep -c 'M-B3 ' "$FIX")"
  echo "AFTER 'M-B RED'-family line map (prefix -> which case owns it):"
  grep -n 'M-B[0-9] ' "$FIX"
  echo "--- SCOPE PROOF: fold M-B1/M-B2/M-B3 back to M-B, diff vs the before blob. Expect ONLY"
  echo "    the AC#8(i) addition and the AC#8(iii) wording changes, nothing else. ---"
  sed 's/M-B1 /M-B /g; s/M-B2 /M-B /g; s/M-B3 /M-B /g' "$FIX" > "$LOGS/folded-prefix-copy"
  diff "$LOGS/before-35r2-fixture.copy" "$LOGS/folded-prefix-copy" > "$LOGS/scope-diff-ii.txt"
  echo "scope-diff changed-line-count=$(grep -c '^[<>]' "$LOGS/scope-diff-ii.txt") (lines: (i) case added + (iii) three claims)"
  echo "assertion-strength sanity: t.Fatalf count before=$(grep -c 't\.Fatalf' "$LOGS/before-35r2-fixture.copy") after=$(grep -c 't\.Fatalf' "$FIX")"
  echo "test-func count before=$(grep -c '^func Test' "$LOGS/before-35r2-fixture.copy") after=$(grep -c '^func Test' "$FIX")"
  echo "rc=0"; } > "$LOGS/after-ii.txt" 2>&1

{ echo "=== (iii) three browser-fact claims, before vs after ($(date '+%F %T%z')) ==="
  echo "--- BEFORE, verbatim from the daf1f5a0 blob ---"
  grep -n 'Chrome and WebView2 answer an\|the way a browser.s host method does\|A real browser answers that with Illegal invocation' "$LOGS/before-35r2-fixture.copy"
  echo "before-claim-count=$(grep -c 'Chrome and WebView2 answer an\|the way a browser.s host method does\|A real browser answers that with Illegal invocation' "$LOGS/before-35r2-fixture.copy") [want 3]"
  echo "--- AFTER, verbatim from the working copy ---"
  grep -n 'THIS FIXTURE answers\|THIS RULER CANNOT PRODUCE\|modelled after what Chrome\|modelled after Chrome' "$FIX"
  echo "after-unqualified-claim-count=$(grep -c 'Chrome and WebView2 answer an\|the way a browser.s host method does\|A real browser answers that with Illegal invocation' "$FIX") [want 0]"
  echo "--- the three after-strings in full (each must still name its case's red) ---"
  grep -n -A1 'THIS FIXTURE answers\|THIS RULER CANNOT\|the rule THIS FIXTURE answers with' "$FIX"
  echo "M-A RED prefix still present: $(grep -c 'M-A RED' "$FIX") occurrences"
  echo "rc=0"; } > "$LOGS/after-iii.txt" 2>&1

{ echo "=== (iv) the 〔仅本机可量〕 label on its own file, before vs after ($(date '+%F %T%z')) ==="
  echo "--- BEFORE (daf1f5a0 blob) : the family grep the ticket said to re-measure ---"
  git show "daf1f5a0:$LIVE" > "$LOGS/before-livedoc.copy"
  grep -n '仅本机\|machine-local\|CI 永看不见' "$LOGS/before-livedoc.copy"; echo "before-grep-rc=$? [1 = zero hits, as AC#8(iv) predicted]"
  echo "before-count=$(grep -c '仅本机\|machine-local\|CI 永看不见' "$LOGS/before-livedoc.copy")"
  echo "--- AFTER (working copy) ---"
  grep -n '仅本机\|machine-local\|CI 永看不见' "$LIVE"; echo "after-grep-rc=$?"
  echo "after-count=$(grep -c '仅本机\|machine-local\|CI 永看不见' "$LIVE")"
  echo "--- four elements present in the header block? ---"
  for k in '① WHAT IS MACHINE-LOCAL' '② WHO OWNS IT' '③ RE-RUN CADENCE' '④ THE PRICE PAID'; do
    echo "$k -> $(grep -c "$k" "$LIVE") hit(s)"
  done
  echo "orchestrator-owner=$(grep -c 'ORCHESTRATOR' "$LIVE")  probes-35-cadence=$(grep -c '.scratch/wisp/probes/35' "$LIVE")  revocation=$(grep -c 'Revocation' "$LIVE")"
  echo "file length before=$(wc -l < "$LOGS/before-livedoc.copy") after=$(wc -l < "$LIVE")"
  echo "rc=0"; } > "$LOGS/after-iv.txt" 2>&1

{ echo "=== zero production-code change (self-proof) ($(date '+%F %T%z')) ==="
  echo "git diff --numstat daf1f5a0 -- $PROD  ->  [expect NO output]"
  git diff --numstat daf1f5a0 -- "$PROD"
  echo "numstat-rc=$? (rc 0 with empty output = that file has no line changed)"
  echo "prod blob now=$(git hash-object "$PROD") at-daf1f5a0=$(git rev-parse "daf1f5a0:$PROD")"
  echo "changed tracked files on this leg:"
  git diff --name-only daf1f5a0 -- cmd/wisp
  echo "full name-status of everything 35-r5 stages (must be only the fixture, the winlive doc and probes/35/r5):"
  git status --porcelain -- cmd/wisp .scratch/wisp/probes/35/r5
  echo "rc=0"; } > "$LOGS/zero-prod-change.txt" 2>&1

for f in after-ii after-iii after-iv zero-prod-change; do echo "--- $f ---"; tail -4 "$LOGS/$f.txt"; done
