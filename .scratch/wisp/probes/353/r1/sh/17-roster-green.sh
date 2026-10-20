#!/bin/sh
# 353-r1 名册照绿尺＝派单点名的两枚按 `why` 子串作差的归因名册用例 ⇄ 两枚⛔ 许受影响的既有钉，
# 在【treeB2＝产码基线＋同一对测试字节】⇄【treeL＝乙已落】两侧逐名查（★尺＝^--- (PASS|FAIL|SKIP): 前缀名，⛔ 子串）。
set -u
cd /d/tmp/wisp/20261020/353r1/logs || exit 9
NAMES='TestAC1LeafRaceWindowsBareSiblingReachesOpenPinnedLeaf342r2|TestAC2WindowsAncestorWorldAttributesEveryRefusalToALeg342r1|TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1|TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2|TestAC5FailedSealRefusesTheWrite|TestAC5FailureIsNotSwallowedByTheHappyPath|TestAC1WindowsNameTopUp|TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1'
for f in pername-s1-treeB2.raw.txt pername-s2-treeL.raw.txt fullpair-s1-treeB.raw.txt fullpair-s2-treeL.raw.txt; do
  [ -f "$f" ] || continue
  echo "== $f =="
  grep -E "^--- (PASS|FAIL|SKIP)" "$f" | grep -E "$(echo "$NAMES" | sed 's/|/\\|/g' | sed 's/|/|/g')" \
    | sed 's/^/   /; s/ *([0-9.]*s)$//'
done
