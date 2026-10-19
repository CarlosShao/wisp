#!/bin/bash
# 取证：本腿 full-B 里编排者那 15 枚"只在编排者侧红"的名字各自的颜色＋包级失败行＋死亡标记。
set -u
L=/d/tmp/wisp/20261019/227v2/full-B.log
echo '== package-level FAIL lines (full-B) =='
grep -E '^FAIL[[:space:]]+github' "$L"
echo '== death markers in full-B =='
grep -nE 'exit status|signal: |fatal error|test timed out|0xc0000' "$L" | head -20
echo '== verdict of each name that was RED only in orchestrator A1053 shot =='
for n in Test20r2JunctionRefusalIsNeverAskable \
         TestBridgeRefusesARealJunctionOnTheReadRoute \
         TestBridgeWritesNothingThroughARealJunction \
         TestSweepNeverDeletesThroughARealJunction \
         TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing \
         TestSyncAncestorActableLegFailsClosedOnMovedAncestor116 \
         TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116 \
         TestSyncAncestorLegReResolvesThroughJunctionToRealTree \
         TestSyncAncestorLegNotExemptedStaysFailClosed \
         TestDirOfBaseOfCutOnThePlatformSeparator \
         TestGateJudgesThePrivateSetByResolvedSID \
         TestSealNarrowsAndNamesThePrincipalItRemovedBySID \
         TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop \
         TestAC4InstallDirWidthIsDeliberateAndHasABreakLine \
         TestAC2SealNoticeLandsInTheRunLegLogFile ; do
  v=$(grep -m1 -E "^(    )*--- (PASS|FAIL|SKIP): ${n}( |\$|/)" "$L" | sed 's/^ *//')
  [ -z "$v" ] && v="NOT-RUN-IN-THIS-SHOT"
  printf '%s\n' "$v   << $n"
done
echo '== my 5 new reds: first body line (go test prints body BEFORE the --- FAIL line) =='
for n in Test258BridgeMutationNoSrcKeepsOldBinding Test258OccupiedCombinationNamesTheNewValue \
         Test258V1ProbeSummonEditRebindsLiveBall TestAC246DevLegIgnoresTheTestTaskInjection \
         TestTicket252R1AlignmentAddsNoAuthorization ; do
  ln=$(grep -nE "^--- FAIL: ${n} " "$L" | head -1 | cut -d: -f1)
  echo "--- $n at log line $ln ---"
  if [ -n "$ln" ]; then sed -n "$((ln-6)),${ln}p" "$L"; fi
done
