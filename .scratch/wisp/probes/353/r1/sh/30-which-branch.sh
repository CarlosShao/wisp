#!/bin/sh
# 353-r1 具名连带尺＝那一枚 348 族仪器（TestAC2...NeverDeletesAForeignFile342r1）在两侧各 4 发里
#   **到底响在哪一支断言**＝逐字取回。⛔ 拿名册名当凭据而不看它为什么红＝本票⛔ 借别的族的绿/红。
set -u
L=/d/tmp/wisp/20261020/353r1/logs
for TREE in finalB finalL; do
  for s in 1 2 3 4; do
    f="$L/swing-$TREE-$s.raw.txt"
    [ -f "$f" ] || continue
    echo "== $TREE shot=$s =="
    awk '
      /^=== RUN   TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1$/ {p=1; next}
      p && /^--- (FAIL|PASS): TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1/ {print "  VERDICT: " $0; p=0}
      p && /ancestorunlink_348_windows_test.go:/ {print "  " $0}
    ' "$f" | grep -E 'RED|VERDICT' | sed 's/  */ /g'
    echo "  numerator: $(grep -oE 'form=symlink-dir attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$f" | head -1)"
  done
done
