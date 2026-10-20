#!/bin/sh
# gw-1 具名连带逐枚色（派单 ②③ 都点名那味＝⛔ 只报"新增红 0"、要逐枚点名两侧色）。
# 尺＝在 `-v` 原始件里锚定行首 `^(    )*--- (PASS|FAIL|SKIP): <NAME>`；两侧同尺；
# 名⛔ 出现在该档＝ABSENT（＝该 build tag 在这一档⛔ 编译＝弃权形，⛔ 算绿⛔ 算红）。
# ★⛔ 拿四数顶逐枚；★正控＝同一条尺打一枚已知在册名，读出色＝尺⛔ 瞎。
set -u
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
NAMES="$R/names.txt"
: > "$NAMES"
cat >> "$NAMES" <<'EOF'
TestAC5FailedSealRefusesTheWrite
TestAC5FailureIsNotSwallowedByTheHappyPath
TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1
TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1
TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone
TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth
TestAC1POSIXSealDirThroughASymlinkRefuses
TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing
TestAC1POSIXSealFileThroughABackslashNamedLink
TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree
TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks
TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator
TestAC2WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1
TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1
TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1
TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1
TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1
TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1
TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1
TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1
TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2
TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2
TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1
TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1
TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1
EOF
echo "NAMES=$(grep -c . "$NAMES")"
for SRC in "$LOGS"/p352w-s1-gwB-winsec.raw.txt "$LOGS"/p352w-s2-gwL-winsec.raw.txt "$LOGS"/p353-s1-gwW-winsec.raw.txt "$LOGS"/p353-s2-gwL-winsec.raw.txt "$LOGS"/c352w-s1-gwB-winsec.raw.txt "$LOGS"/c352w-s2-gwL-winsec.raw.txt "$LOGS"/c353-s1-gwW-winsec.raw.txt "$LOGS"/c353-s2-gwL-winsec.raw.txt; do
  [ -f "$SRC" ] || { echo "SKIP_SRC(absent) $SRC"; continue; }
  TAG=$(basename "$SRC" .raw.txt)
  echo "=== $TAG ==="
  while IFS= read -r N; do
    [ -n "$N" ] || continue
    HIT=$(grep -E "^(    )?--- (PASS|FAIL|SKIP): $N(\([0-9]+\))? " "$SRC" | head -4)
    if [ -z "$HIT" ]; then
      echo "  ABSENT  $N"
    else
      printf '%s\n' "$HIT" | sed -E 's/ *(\([0-9.]+s\))$//' | sed "s/^/  GOT     $TAG /"
    fi
  done < "$NAMES"
  echo "  -- 子测试逐枚（private_fail 那五枚判定所住的宿主用例）--"
  grep -E "^    --- (PASS|FAIL|SKIP): " "$SRC" | sed -E 's/ *(\([0-9.]+s\))$//' | sort -u | sed "s/^/  SUB $TAG /"
done
echo
echo "== 正控（尺⛔ 瞎）＝同一条尺对一枚已知在两档都存在的名读出两行 =="
echo "  （＝上面 GOT 行数非零即正控；另加一枚刻意不存在的名应读 ABSENT）"
grep -E "^(    )?--- (PASS|FAIL|SKIP): DefinitelyNoSuchTestGW1 " "$LOGS"/p353-s2-gwL-winsec.raw.txt 2>/dev/null | sed 's/^/  SHOULD_BE_EMPTY /'
echo "  NEGCTRL_HITS=$(grep -cE '^(    )?--- (PASS|FAIL|SKIP): DefinitelyNoSuchTestGW1 ' "$LOGS"/p353-s2-gwL-winsec.raw.txt 2>/dev/null || true)"
