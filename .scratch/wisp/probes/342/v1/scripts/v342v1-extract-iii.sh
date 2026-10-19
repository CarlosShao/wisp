#!/usr/bin/env bash
# 342-v1: extract the yinofb mis-shape control's verbatim red sentences (mandatory (iii))
# and the A/B evidence that the two baseline reds on my A container config were a bench
# defect of this leg, not a reading about ticket 342.
set -u
S=/d/tmp/wisp/20261018/342v1
R="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/342/v1/logs"
mkdir -p "$R"
f="$S/logs/17-colorB-col-yinofb.log"
{
  echo "RAW=$f"
  echo "RECOMPUTED: RUN=$(grep -c '^=== RUN' "$f" || true) PASS_TOP=$(grep -c '^--- PASS' "$f" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$f" || true) SKIP=$(grep -c '^--- SKIP' "$f" || true)"
  echo "--- top-level red names (verbatim, 4 expected):"
  grep -E '^--- FAIL' "$f" | sed 's/ ([0-9.]*s)//'
  echo
  echo "--- per-test verbatim RED sentences:"
  for n in TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1 \
           TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1 \
           TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3 \
           TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3; do
    echo "===== $n"
    awk -v want="$n" 'BEGIN{inside=0} /^=== RUN   /{inside=(index($0,want)>0)} inside && /AC#[0-9]/{print}' "$f" | head -8
  done
  echo
  echo "===== baseline contrast: top-level FAIL count on the SAME container config, pristine treeC:"
  grep -cE '^--- FAIL' "$S/logs/17-colorB-treeC.log" || true
  echo "===== A-config (this leg's mistake) baseline treeC reds, kept so the attribution is checkable:"
  grep -E '^--- FAIL' "$S/logs/13-color-treeC.log" | sed 's/ ([0-9.]*s)//' || true
  echo "A-config four numbers were RUN=71 PASS_TOP=41 FAIL_TOP=2 SKIP=0 (13-color-treeC.log, two invocations: 2 and 3 reds); B-config (matched to the leg's driver env, no GOTMPDIR bind) is 71/43/0/0 three times."
} > "$R/18-yinofb-and-config-attribution.txt"
wc -c "$R/18-yinofb-and-config-attribution.txt"
