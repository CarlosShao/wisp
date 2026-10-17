#!/bin/sh
# 333-v1 AC#1 (three starting rulers, re-run independently) + AC#7 (call-chain file:line
# re-recognised with sed, front-end receiver counts). Every tree reading is from the
# NAMED delivery tree exported to D:\W333V1-SuqCO8 (git archive 859026d2), and the
# pre-anchor readings from D:\W333V2-hLzM9f (git archive ae710c49). Never the shared
# working tree.
A=/d/W333V1-SuqCO8
P=/d/W333V2-hLzM9f
echo "### RULER 1: literal codepoint roster in the ANCHOR version of isDisplayWashCf"
awk '/^func isDisplayWashCf/,/^}$/' "$P/internal/panel/display_scrub_216.go" | head -30
echo "-- hex literals counted one by one (awk over the function body only) --"
awk '/^func isDisplayWashCf/,/^}$/' "$P/internal/panel/display_scrub_216.go" \
  | grep -o -E '0x[0-9A-Fa-f]{2,4}' | sort -u | tee /tmp/anchor-roster.txt | wc -l
echo "-- case-line count inside the body (cross-line case trap) --"
awk '/^func isDisplayWashCf/,/^}$/' "$P/internal/panel/display_scrub_216.go" | grep -c 'case 0x' || true
echo "-- whole-file grep -c 'case 0x' (the trap the ticket names: reads 1) --"
grep -c 'case 0x' "$P/internal/panel/display_scrub_216.go" || true
echo
echo "### RULER 3: washed fields + production call sites (delivery tree)"
echo "-- non-test hits of scrubComposerForDisplay --"
grep -rn 'scrubComposerForDisplay' "$A/internal" "$A/cmd" "$A/tools" 2>/dev/null | grep -v '_test.go'
echo "-- test hits (for the denominator) --"
grep -rln 'scrubComposerForDisplay' "$A/internal" "$A/cmd" 2>/dev/null | grep '_test.go' | wc -l
echo "-- composer.go:145 verbatim --"
sed -n '143,147p' "$A/internal/panel/composer.go"
echo "-- display_scrub_216.go:218 / :224 verbatim --"
sed -n '217,219p;223,225p' "$A/internal/panel/display_scrub_216.go"
echo
echo "### AC#7 chain, each file:line re-read with sed (delivery tree)"
for spec in "cmd/wisp/run.go:714" "cmd/wisp/panel_pump.go:451" "internal/panel/pump.go:354" \
            "internal/panel/pump.go:309" "internal/panel/composer.go:106" "internal/panel/composer.go:145" \
            "internal/panel/display_scrub_216.go:218" "internal/panel/display_scrub_216.go:224" \
            "internal/panel/pump.go:359" "internal/panel/composer.go:231" "internal/panel/attachments.go:86"; do
  f=${spec%%:*}; n=${spec##*:}
  printf '%s -> ' "$spec"
  sed -n "${n}p" "$A/$f"
done
echo
echo "### AC#7 front-end: renderer vs receiver (delivery tree, frontend/src)"
echo "-- renderer: canonical reads --"
grep -rn -i 'canonical' "$A/frontend/src" | head -20
echo "HITS_canonical=$(grep -rn -i 'canonical' "$A/frontend/src" | wc -l)"
echo "PC_HITS_composer=$(grep -rn -i 'composer' "$A/frontend/src" | wc -l)"
echo "-- receiver candidates --"
for pat in 'onmessage' 'chrome.webview' 'chrome\.webview' 'addEventListener' 'hostObject' 'window\.external'; do
  printf '%s -> ' "$pat"; grep -rn -i -- "$pat" "$A/frontend/src" | wc -l
done
echo "-- addEventListener(\"message\") specifically --"
grep -rn 'addEventListener("message"' "$A/frontend/src" | wc -l
grep -rn 'addEventListener' "$A/frontend/src" | head -10
echo "-- main.tsx mount line --"
sed -n '90,100p' "$A/frontend/src/main.tsx"
echo "-- App.tsx EMPTY default --"
sed -n '85,95p' "$A/frontend/src/App.tsx"
