#!/usr/bin/env bash
# 342-v1 host (Windows) phase: mandatory (i) Windows rows, 1500 attempts x3 repeats.
set -u
S=/d/tmp/wisp/20261018/342v1
L=$S/logs
cp "$S/probe/v342v1-windows-forms.go.txt" "$S/treeW/internal/winsec/v342v1win_test.go"
echo "probe sha = $(sha256sum "$S/treeW/internal/winsec/v342v1win_test.go" | cut -d' ' -f1)"
for r in 1 2 3; do
  ( cd "$S/treeW" && unset GOTMPDIR TMPDIR TEMP TMP; V342V1_ATTEMPTS=1500 go test -count=1 -timeout 40m -v \
      -run TestV342v1WindowsFormsRace ./internal/winsec/ > "$L/31-win-race-r$r.log" 2>&1; echo "RC=$?" >> "$L/31-win-race-r$r.log" )
  echo "== r$r rows:"
  grep -E 'V342V1 W[A-Z]+: attempts' "$L/31-win-race-r$r.log" || true
done
echo "== deterministic Windows rows (W-LEAFSWAP + W-ROOTSYM/W-ROOTJUNC), 3 repeats:"
for r in 1 2 3; do
  ( cd "$S/treeW" && unset GOTMPDIR TMPDIR TEMP TMP; go test -count=1 -timeout 20m -v \
      -run 'TestV342v1WindowsLeafSwap|TestV342v1WindowsOsRoot' ./internal/winsec/ > "$L/32-win-det-r$r.log" 2>&1; echo "RC=$?" >> "$L/32-win-det-r$r.log" )
  grep -E 'V342V1 W-' "$L/32-win-det-r$r.log" || true
done
echo "HOST-W-DONE"
