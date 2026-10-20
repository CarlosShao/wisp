#!/bin/sh
# 353-r1 交件前去牙终量＝★打**仓外 finalL 树**的物理字节（＝已入库那批 HEAD blob 取出来的）；
#   ⛔ -overlay；★内容锚、⛔ 行号；落刀前后各验锚命中数；每发当场还原并报还原哈希；★串行。
# 三发＝票面 AC#2 的成对两发（logs/04 §四；★本腿自己重打、⛔ 采信 342-v3 的数）＋乙新增那枚"复验那步"发。
set -u
R=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$R/logs"
TREE="$R/finalL"
TARGET="$TREE/internal/winsec/pinnedleaf_342_windows.go"
PRISTINE="$R/pristine-final-pinnedleaf_342_windows.go"
RUNPAT='TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1|TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1|TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2'
cp -f "$R/incoming/treeL/pinnedleaf_342_windows.go" "$PRISTINE"
H0=$(sha256sum "$PRISTINE" | cut -d' ' -f1)
for k in 1 2 3; do
  case "$k" in
    1) OLD='	_ = p.root.Remove(base)'
       NEW='	_ = base; _ = now; _ = pinned // MUT-1: removal step removed' ;;
    2) OLD='	f, err := p.root.OpenFile(base, flags, perm)'
       NEW='	f, err := os.OpenFile(full, flags, perm) // MUT-2: by-name open (baseline shape)' ;;
    3) OLD='	if err != nil || !os.SameFile(pinned, now) {'
       NEW='	if err != nil && !os.SameFile(pinned, now) { // MUT-3: 乙 identity re-check neutered' ;;
  esac
  echo "===== TOOTH=$k ====="
  HT=$(sha256sum "$TARGET" | cut -d' ' -f1)
  [ "$HT" = "$H0" ] || { echo "NOT_PRISTINE_BEFORE=$HT (restore first)"; cp -f "$PRISTINE" "$TARGET"; }
  AB=$(grep -cF "$OLD" "$TARGET" || true)
  echo "ANCHOR_BEFORE=$AB (期望 1)"
  [ "$AB" = 1 ] || { echo "ABORT_ANCHOR"; continue; }
  MUTDIR="$R/mut-final-$k"
  rm -rf "$MUTDIR" 2>/dev/null || true
  cp -a "$TREE" "$MUTDIR"
  T2="$MUTDIR/internal/winsec/pinnedleaf_342_windows.go"
  awk -v old="$OLD" -v new="$NEW" 'index($0, old) == 1 && !done { print new; done = 1; next } { print }' "$T2" > "$T2.new"
  mv "$T2.new" "$T2"
  MA=$(grep -cF "$NEW" "$T2" || true); OA=$(grep -cF "$OLD" "$T2" || true)
  echo "MUT_ANCHOR=$MA (期望 1) OLD_ANCHOR_AFTER=$OA (期望 0) MUT_SHA256=$(sha256sum "$T2" | cut -d' ' -f1)"
  [ "$MA" = 1 ] && [ "$OA" = 0 ] || { echo "ABORT_MUT_VERIFY"; continue; }
  OUT="$LOGS/final-mut-$k.raw.txt"
  TS=$(date +%Y%m%d-%H%M%S)-$$
  GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1fm-$TS"
  mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP"; continue; }
  DLLP=$(cygpath -w "$ROOTW\\mut-final-$k\\build"); SHP=$(cygpath -w "$ROOTW\\mut-final-$k\\third_party\\sherpa-onnx")
  ( cd "$ROOTW\\mut-final-$k" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
      GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 40m -run "$RUNPAT" ./internal/winsec/ ) > "$OUT" 2>&1
  RC=$?
  echo "RC=$RC RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
  echo "RED_NAMES:"; grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: /   /; s/ *([0-9.]*s)$//'
  grep -E '353r1 AC#2 over|AC#1-one over|353r1 AC#2 CONTROL|did not ENTER' "$OUT" | sed 's/^/   NUM| /'
  cp -f "$PRISTINE" "$TARGET"
  HR=$(sha256sum "$TARGET" | cut -d' ' -f1)
  echo "RESTORE_SHA256=$HR MATCHES_PRISTINE=$([ "$HR" = "$H0" ] && echo YES || echo NO)"
done
echo "===== 还原后 finalL 与 finalB 各自的 pinnedleaf（期望 finalB 全程⛔ 动）====="
sha256sum "$TARGET" "$R/finalB/internal/winsec/pinnedleaf_342_windows.go" | sed 's/^/  /'
