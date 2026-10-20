#!/bin/sh
# ★本件＝笔④之前的去牙草稿（打 treeL＝上一版台面），已被 sh/24-final-mutations.sh 取代（打 finalL＝HEAD blob 重 archive 的终量台面）；
#   按"临时件只建不删"留在盘上入库当形状记录＝它的三发读数（同形、同还原哈希）见 readings-ac2-baseline.md 时代，交件以 sh/24 为准。
# 353-r1 敏感度（去牙）＝只打**仓外导出树**的物理字节；★⛔ -overlay（那形对读盘的尺失明）；★内容锚、⛔ 绝对行号；
#   落刀前后各验锚命中数（⛔ 空操作伪装成一次成功复跑＝公共块 突变那一节）；每发**当场还原**并报还原哈希；
#   ★两发串行（跨 shell 的 trap/finally 有竞态）；成对两发落点⛔ 同。
# 三发＝票面 AC#2 命名的成对两发（logs/04 §四＝`342-v3` 打过、★本腿自己重打、只押形状与量级⛔ 要求逐字同号）
#       ＋一枚随乙新增的"复验那步"发（＝证明 damage 那枚 0 是**复验拦的**、⛔ 是别的偶然）：
#   1 ＝摘 `unlinkPinnedLeaf` 里那步删除（`_ = p.root.Remove(base)`）⇒ 期望 damage 0、entered 照 3、★正控转红
#   2 ＝`openPinnedLeaf` 那发开句柄换成基线那形**按名开** ⇒ 期望 entered 3→0（＝本仪器那条"⛔ 进得去"的红）
#   3 ＝摘掉乙的同一性复验那一步（让它恒通过）⇒ 期望 damage 回到 3（＝拦住删的就是这枚复验）
# 用法：sh 09-mutations.sh <1|2|3>
set -u
ROOT=/d/tmp/wisp/20261020/353r1
L="$ROOT/logs"
TREE="$ROOT/treeL"                      # ★⛔ 碰 treeB；★⛔ 碰共享工作树字节
TARGET="$TREE/internal/winsec/pinnedleaf_342_windows.go"
PRISTINE="$ROOT/pristine-landed-pinnedleaf_342_windows.go"
MUTDIR="$ROOT/mut-treeL"
TOOTH="$1"

# pristine＝乙已落那批字节（与 pair2-treeL 同串＝先验一次）
cp -f "$ROOT/incoming/treeL/pinnedleaf_342_windows.go" "$PRISTINE"
H0=$(sha256sum "$PRISTINE" | cut -d' ' -f1)
H1=$(sha256sum "$TARGET" | cut -d' ' -f1)
echo "PRISTINE_SHA256=$H0"
[ "$H0" = "$H1" ] || { echo "TREE_L_NOT_PRISTINE=$H1 ABORT"; exit 8; }

case "$TOOTH" in
  1) OLD='	_ = p.root.Remove(base)'
     NEW='	_ = base; _ = now // MUT-1: removal step removed' ;;
  2) OLD='	f, err := p.root.OpenFile(base, flags, perm)'
     NEW='	f, err := os.OpenFile(full, flags, perm) // MUT-2: by-name open (baseline shape)' ;;
  3) OLD='	if err != nil || !os.SameFile(pinned, now) {'
     NEW='	if err != nil && !os.SameFile(pinned, now) { // MUT-3: 乙 identity re-check neutered' ;;
  *) echo "BAD_TOOTH=$TOOTH"; exit 8 ;;
esac

echo "== TOOTH=$TOOTH @ $(date '+%F %T %z') =="
AB=$(grep -cF "$OLD" "$TARGET" || true)
echo "ANCHOR_BEFORE=$AB (期望 1)"
[ "$AB" = 1 ] || { echo "ANCHOR_BEFORE_NOT_1 ABORT"; exit 7; }

rm -rf "$MUTDIR" 2>/dev/null || true
cp -a "$TREE" "$MUTDIR"
T2="$MUTDIR/internal/winsec/pinnedleaf_342_windows.go"
awk -v old="$OLD" -v new="$NEW" 'index($0, old) == 1 && !done { print new; done = 1; next } { print }' "$T2" > "$T2.new"
mv "$T2.new" "$T2"

MA=$(grep -cF "$NEW" "$T2" || true); OA=$(grep -cF "$OLD" "$T2" || true)
echo "MUT_ANCHOR=$MA (期望 1)  OLD_ANCHOR_AFTER=$OA (期望 0)  MUT_SHA256=$(sha256sum "$T2" | cut -d' ' -f1)"
[ "$MA" = 1 ] && [ "$OA" = 0 ] || { echo "MUTATION_VERIFY_FAILED ABORT"; exit 6; }

OUT="$L/mut-$TOOTH.raw.txt"
TS=$(date +%Y%m%d-%H%M%S)-$$
GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-mut$TOOTH-$TS"
mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 5; }
ROOTW='D:\tmp\wisp\20261020\353r1'
DLLP=$(cygpath -w "$ROOTW\\mut-treeL\\build"); SHP=$(cygpath -w "$ROOTW\\mut-treeL\\third_party\\sherpa-onnx")
{ echo "=== 353r1 mut-$TOOTH @ $(date '+%F %T %z') GOTMP=$GOTMPW"; } > "$OUT"
( cd "$ROOTW\\mut-treeL" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
    GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 40m \
    -run 'TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1|TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1|TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2' \
    ./internal/winsec/ ) >> "$OUT" 2>&1
RC=$?
{
  echo "RC=$RC"
  echo "RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
  echo "-- 红册（尺＝^--- FAIL）--"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort
  echo "-- 分子行 --"
  grep -E '353r1 AC#2 over|AC#1-one over|353r1 AC#2 CONTROL|did not ENTER' "$OUT"
  echo "-- 当场还原（⛔ git checkout，靠仓外 pristine 拷贝）--"
  cp -f "$PRISTINE" "$TARGET"
  HR=$(sha256sum "$TARGET" | cut -d' ' -f1)
  echo "RESTORE_SHA256=$HR"
  echo "RESTORE_MATCHES_PRISTINE=$([ "$HR" = "$H0" ] && echo YES || echo NO)"
  echo "TREEB_UNTOUCHED_SHA256=$(sha256sum "$ROOT/treeB/internal/winsec/pinnedleaf_342_windows.go" | cut -d' ' -f1)"
} | tee -a "$OUT"
