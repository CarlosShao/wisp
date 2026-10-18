#!/bin/bash
# 339-r1 AC#4 名册逐名作差尺（可重跑）＝起手⇄交件，两平台分别 comm/diff。
# ⚠ 第一稿我用内联命令跑，本件把它钉成尺；logs 里的 lb/lf/wb/wf.txt＝两平台起手⇄交件的顶层名册排序件（只建不删）。
set -u
cd "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/339/r1/logs" || exit 9
names() { grep -E "^ *--- (PASS|FAIL|SKIP)" "$1" | sed -E "s/ \(.*//" | sort; }
echo "=== LINUX 起手(ac0)⇄交件(final) 逐名作差 ==="
comm -23 <(names linux-ac0-fullpkg.log) <(names linux-final-fullpkg.log)
echo "--- 交件新增 ---"
comm -13 <(names linux-ac0-fullpkg.log) <(names linux-final-fullpkg.log)
echo "=== WINDOWS 起手⇄交件（期望＝空＝名册逐字相同）==="
diff <(names win-ac0-fullpkg-psnative.log) <(names win-final-fullpkg-psnative.log) && echo "IDENTICAL"
echo "=== 关键两枚＋绊线＋C26 的逐名色（final 两平台）==="
for t in TestAnchorMissReasonNamesCandidateAndMissingAnchor TestAnchorMissKeepsNonOverridableClassAndNeedL2 TestAnchorMissHasNoWorldOnThisPlatform339 TestC26RewriteAccountIsConsumedAtEverySecurityLeg TestFileBaseNameClaimChangesNoClassOrReason TestAnchorFallbackTeethPerAnchor336 TestFallbackBranchHasNoWorldOnThisPlatform336; do
  for p in linux-final-fullpkg win-final-fullpkg-psnative; do
    line=$(grep -E "^--- [A-Z]+: $t" "$p.log" | head -1 | sed -E 's/ \(.*//')
    [ -n "$line" ] && printf "%-60s %-26s %s\n" "$t" "$p" "${line#--- }"
  done
done
