#!/usr/bin/env bash
# 336-v1 correction: the previous AC#4 surface script printed teeth_red_hits=8 with a
# WRONG SCOPE (it grepped the whole log, so === RUN / --- PASS / --- CONT lines naming
# the teeth file were counted). This ruler re-measures with the four-count scope only:
# a red is a line matching  ^[[:space:]]*--- FAIL:
set -u
REPO="D:/work/workspace/projects plans/Wisp"
SRC="${LOGSRC:-/d/tmp/wisp/20261018/336v1/logs}"
{
echo "### corrected ruler: FAIL lines only (top-level and subtest forms both start with --- FAIL:)"
for f in ac4full-ci-after ac4full-d-after ac4full-ci-before; do
  n_teeth=$(grep -E '^[[:space:]]*--- FAIL: ' "$SRC/$f.log" | grep -c 'anchor_fallback_teeth_336' || true)
  n_fail=$(grep -cE '^[[:space:]]*--- FAIL: ' "$SRC/$f.log" || true)
  n_c26=$(grep -E '^[[:space:]]*--- FAIL: ' "$SRC/$f.log" | grep -c 'TestC26' || true)
  printf "%s fail_lines_naming_teeth=%s fail_lines_total=%s (positive control: fail lines naming TestC26 = %s)\n" "$f" "$n_teeth" "$n_fail" "$n_c26"
done
echo "### what the 8 wrong-scope hits actually were (line kinds, first 12):"
grep -n 'anchor_fallback_teeth_336' "$SRC/ac4full-d-after.log" | head -12 | cut -c1-150 | sed 's/^/  /'
echo "### named surface blemish in run-ac4-surface.sh: the ci-shape GOTMPDIR line has a misplaced quote so the dir name ends with '}' (visible in ac4-surface-runner.log). Same volume, same shape, same TMP/TEMP handling => the four-count and roster readings stand; the script is kept as-run (no deletes) and corrected here."
echo "correction_rc=0"
} 2>&1 | tee "$SRC/ac4-teeth-red-correction.log"
