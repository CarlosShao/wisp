#!/bin/sh
# ★本件＝笔④后台面草稿，已被 sh/23-final-verify.sh 取代（后者重 archive 两棵树、同一把合一尺）；
#   按"临时件只建不删"入库当形状记录；交件读数一律以 sh/23⇄sh/24⇄sh/25 为准。
# 353-r1 交件前终量尺＝★随笔④动过字节 ⇒ 所有读数重取一遍（⛔ 拿笔④之前的数顶替＝公共块读数纪律）。
# 一次跑完：① 定向成对（修前红 ⇄ 落地绿）② 去牙三发 ③ 整包 host 成对＋逐名矩阵 ④ 整包 container 成对。
# ★串行；★⛔ 碰 treeB 的产码字节（它＝基线）；★⛔ 共享工作树跑 go test。
set -u
ROOT=/d/tmp/wisp/20261020/353r1
cd "$ROOT" || exit 9
echo "=== FINAL-1 同步终字节进 treeL（产码三枚＋仪器件）==="
cp -f incoming/treeL/pinnedleaf_342_windows.go treeL/internal/winsec/
cp -f incoming/treeL/pinnedleaf_342_unix.go treeL/internal/winsec/
cp -f incoming/treeL/pinnedleaf_342_fallback.go treeL/internal/winsec/
cp -f incoming/treeL/winsec.go treeL/internal/winsec/
cp -f incoming/treeL/nametopup_353_windows_test.go treeL/internal/winsec/
cp -f incoming/treeL/nametopup_353_windows_test.go treeB/internal/winsec/
cp -f incoming/treeL/nametopup_353_windows_test.go treeB2/internal/winsec/
cp -f incoming/treeL/ancestorwindow_342r2_windows_test.go treeB2/internal/winsec/
for t in treeB treeB2 treeL; do
  echo "  $t pinnedleaf=$(sha256sum $t/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-12) test=$(sha256sum $t/internal/winsec/nametopup_353_windows_test.go | cut -d' ' -f1 | cut -c1-12)"
done
echo "  treeB  pinnedleaf=$(sha256sum treeB/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-12) (期望 b9dc39ae＝基线未动)"
echo "  treeB2 winsec=$(sha256sum treeB2/internal/winsec/winsec.go | cut -d' ' -f1 | cut -c1-12) pinnedleaf=$(sha256sum treeB2/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-12)"
echo "  treeL  winsec=$(sha256sum treeL/internal/winsec/winsec.go | cut -d' ' -f1 | cut -c1-12) pinnedleaf=$(sha256sum treeL/internal/winsec/pinnedleaf_342_windows.go | cut -d' ' -f1 | cut -c1-12)"
echo "=== FINAL-2 定向成对（treeB 修前 ⇄ treeL 落地）==="
rm -rf pair1-treeB-final pair2-treeL-final 2>/dev/null || true
sh ./sh/02-run-pair.sh host final "TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1|TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1|TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2" nametopup_353_windows_test.go pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go 2>&1 | grep -E 'RC=|RED_NAMES|NUM\|'
echo "=== FINAL-3 去牙三发 ==="
for k in 1 2 3; do sh ./sh/09-mutations.sh "$k" 2>&1 | grep -E 'PRISTINE_SHA256|ANCHOR_BEFORE|MUT_ANCHOR|RC=|RUN=|RED |AC#2 over|AC#1-one over|CONTROL|RESTORE_|TREEB_UNTOUCHED'; done
echo "=== FINAL-4 整包 host 逐名矩阵（treeB2 ⇄ treeL，各 2 发）==="
sh ./sh/14-matrix.sh 2 2>&1 | tail -14
echo "=== FINAL-5 整包 container 成对 ==="
sh ./sh/15-container-pair.sh colfinal treeB2 treeL 2>&1 | grep -E 'RC=|RED_NAMES|ONLY_IN|SYMDIFF'
