#!/bin/sh
# gw-1 五枚 ErrNotSealable 判定 ⇄ 两枚具名族＝逐枚点名两侧色的干净尺（前一发 awk 取字段取错＝$3 是名字不是色，
# 本发改用"色在前"的 sed 形，★正控＝同尺打一枚不存在的名须 0、打一枚在场的名须非 0）。
set -u
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
pick() { # $1=raw $2=名（宿主或前缀）
  grep -E "^(    )?--- (PASS|FAIL|SKIP): $2([^A-Za-z0-9]|$)" "$1" \
    | sed -E 's/^( *)(--- (PASS|FAIL|SKIP): ).*/\3/' 2>/dev/null | tr '\n' ' '
}
color() { # $1=raw $2=完整名（含子测试路径）
  grep -E "^(    )?--- (PASS|FAIL|SKIP): $2 " "$1" | sed -E 's/.*--- (PASS|FAIL|SKIP): .*/\1/'
}
echo "== 五枚判定逐枚点名（色＝顶层用例 ⇄ 四枚子测试；:72⇄:75 同住 seal_file_and_dir_direct）=="
printf '%-26s %-9s %-9s %-9s %-9s %-9s\n' SOURCE TOP exclusive replacement dirchain sealfiledir
for SRC in p352w-s1-gwB-winsec p352w-s2-gwL-winsec p353-s1-gwW-winsec p353-s2-gwL-winsec \
           c352w-s1-gwB-winsec c352w-s2-gwL-winsec c353-s1-gwW-winsec c353-s2-gwL-winsec; do
  F="$LOGS/$SRC.raw.txt"
  [ -f "$F" ] || { echo "$SRC MISSING"; continue; }
  T=$(color "$F" 'TestAC5FailedSealRefusesTheWrite')
  A=$(color "$F" 'TestAC5FailedSealRefusesTheWrite/exclusive_artifact')
  B=$(color "$F" 'TestAC5FailedSealRefusesTheWrite/replacement_write')
  C=$(color "$F" 'TestAC5FailedSealRefusesTheWrite/directory_chain')
  D=$(color "$F" 'TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct')
  printf '%-26s %-9s %-9s %-9s %-9s %-9s\n' "$SRC" "${T:-ABSENT}" "${A:-ABSENT}" "${B:-ABSENT}" "${C:-ABSENT}" "${D:-ABSENT}"
done
echo
echo "== TestAC5FailureIsNotSwallowedByTheHappyPath（同文件另一枚 fs.ErrExist 判定）=="
for SRC in p352w-s1-gwB-winsec p352w-s2-gwL-winsec p353-s1-gwW-winsec p353-s2-gwL-winsec \
           c352w-s1-gwB-winsec c352w-s2-gwL-winsec c353-s1-gwW-winsec c353-s2-gwL-winsec; do
  F="$LOGS/$SRC.raw.txt"; [ -f "$F" ] || continue
  printf '  %-26s %s\n' "$SRC" "$(color "$F" 'TestAC5FailureIsNotSwallowedByTheHappyPath')"
done
echo
echo "== 读 ErrUnresolvedPath 的两族逐枚（openleaf_322 ⇄ placement_symlink_113）=="
for N in TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1 \
         TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone \
         TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth \
         TestAC1POSIXSealDirThroughASymlinkRefuses \
         TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing \
         TestAC1POSIXSealFileThroughABackslashNamedLink \
         TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree \
         TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks \
         TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator \
         TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1 \
         TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1 \
         TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1; do
  LINE=""
  for SRC in p352w-s1-gwB-winsec p352w-s2-gwL-winsec p353-s1-gwW-winsec p353-s2-gwL-winsec \
             c352w-s1-gwB-winsec c352w-s2-gwL-winsec c353-s1-gwW-winsec c353-s2-gwL-winsec; do
    F="$LOGS/$SRC.raw.txt"; [ -f "$F" ] || continue
    C=$(color "$F" "$N"); LINE="$LINE ${C:-ABS}"
  done
  printf '  %-70s%s\n' "$N" "$LINE"
done
echo "  列序＝pB-host pL-host pW-host pL-host(353) cB-cont cL-cont cW-cont cL-cont(353)"
echo
echo "== ★正控两把（同尺）=="
echo "  POSCTRL(在场名应读 PASS)=$(color "$LOGS/p353-s2-gwL-winsec.raw.txt" 'TestAC5FailedSealRefusesTheWrite')"
echo "  NEGCTRL(不存在名应空)=\"$(color "$LOGS/p353-s2-gwL-winsec.raw.txt" 'DefinitelyNoSuchTestGW1')\""
echo
echo "== 具名连带的两侧字节同一性（尺＝git rev-parse <锚>:<path>，量的）=="
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for p in internal/winsec/private_fail_test.go internal/winsec/openleaf_overwrite_322_test.go \
         internal/winsec/openleaf_322_other_test.go internal/winsec/placement_symlink_113_other_test.go \
         internal/winsec/openleaf_guard_322_windows_test.go internal/winsec/nametopup_353_windows_test.go; do
  printf '  %-62s 0fb399c3=%s 3a509e50=%s a19d5bf2=%s\n' "$p" \
    "$(git rev-parse --short=8 0fb399c3:$p 2>/dev/null)" \
    "$(git rev-parse --short=8 3a509e50:$p 2>/dev/null)" \
    "$(git rev-parse --short=8 a19d5bf2:$p 2>/dev/null)"
done
