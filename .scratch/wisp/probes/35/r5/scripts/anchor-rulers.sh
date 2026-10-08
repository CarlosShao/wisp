#!/bin/bash
# 35-r5 opening rulers (起手尺). Every reading is taken live off the disk at run time;
# nothing here is transcribed from the dispatch. Each file ends with its own rc= line.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT="$REPO/.scratch/wisp/probes/35/r5/logs"
cd "$REPO" || exit 1
mkdir -p "$OUT"
A="$OUT/opening-rulers.txt"
{
  echo "=== 35-r5 opening rulers ($(date '+%F %T%z')) ==="
  echo "--- HEAD ---"
  git log --oneline -1
  echo "full-HEAD=$(git rev-parse HEAD)"
  echo "branch=$(git rev-parse --abbrev-ref HEAD)"
  echo "--- porcelain line count (others' work in flight; 35-r5 touches nothing of it) ---"
  echo "porcelain-total=$(git status --porcelain | wc -l)"
  echo "porcelain-on-target-two-files:"
  git status --porcelain -- cmd/wisp/panel_transport_35r2_test.go cmd/wisp/panel_transport_live_35v2_windows_test.go cmd/wisp/panel_host_windows.go
  echo "porcelain-on-target-count=$(git status --porcelain -- cmd/wisp/panel_transport_35r2_test.go cmd/wisp/panel_transport_live_35v2_windows_test.go cmd/wisp/panel_host_windows.go | wc -l)"
  echo "--- ticket 35 box census (AC#8 present, nothing checked by this leg) ---"
  echo "unchecked=$(grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-panel-bridge-c17.md)"
  echo "checked=$(grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/35-panel-bridge-c17.md)"
  echo "--- AC#8 ruler (where the box lives) ---"
  grep -n 'AC#8' .scratch/wisp/issues/35-panel-bridge-c17.md | head -20
  echo "--- fixture size / winlive size ---"
  wc -l cmd/wisp/panel_transport_35r2_test.go cmd/wisp/panel_transport_live_35v2_windows_test.go
  echo "--- (ii) ruler BEFORE: shared red-sentence prefix ---"
  echo "M-B-RED-count=$(grep -c 'M-B RED' cmd/wisp/panel_transport_35r2_test.go)"
  grep -n 'M-B RED' cmd/wisp/panel_transport_35r2_test.go
  echo "--- (iii) ruler BEFORE: the three browser-fact claims ---"
  grep -n 'Chrome and WebView2\|a browser.s host method\|A real browser answers' cmd/wisp/panel_transport_35r2_test.go
  echo "--- (iv) ruler BEFORE: machine-local family words in the winlive file header (expect ZERO) ---"
  grep -n '仅本机\|machine-local\|CI 永看不见' cmd/wisp/panel_transport_live_35v2_windows_test.go
  echo "iv-grep-rc=$? (1 = no match, as ticket 35 AC#8(iv) predicts)"
  echo "--- production guard, NOT to be edited (overlay target) ---"
  grep -n 'typeof window\.%\[1\]s !== "function"' cmd/wisp/panel_host_windows.go
  echo "prod-blob=$(git hash-object cmd/wisp/panel_host_windows.go) HEAD-blob=$(git rev-parse HEAD:cmd/wisp/panel_host_windows.go)"
  echo "--- fixture door construction point (Bind stub) ---"
  grep -n 'libBindStubScript35r2\|func (bf \*fakeDoc35r2) Bind' cmd/wisp/panel_transport_35r2_test.go
} > "$A" 2>&1
echo "rc=$?" >> "$A"
tail -3 "$A"
