#!/bin/sh
# 228-r3c stage 3/4: ONE mutation shot — serial, self-restoring, all inside a SINGLE shell call.
# 规矩：echo 的"已还原"⛔ 当证据。每步落两把哈希尺（物理字节 sha256 ⇄ blob sha1），
# 还原后拿 pristine 对撞；对不上 ⇒ exit 7 ⇒ 下一发绝不开跑（串行硬规矩由这道闸保证，不靠人记得）。
# 跨 shell 的 trap/finally 有竞态（＝R-228-3b-1 的病根）⇒ 本脚本⛔ 用 trap，改"同一次调用内顺序执行＋闸"。
#
# Usage: mutate-shot.sh <M1|M2> <logfile-name-only> <run-pattern>
set -u
ROOT=/d/tmp/wisp/20261018/228r3c
EXPORT=$ROOT/tree
REPO='/d/work/workspace/projects plans/Wisp'
TARGET_REL=internal/proc/boot_windows.go
TARGET=$EXPORT/$TARGET_REL
ANCHOR=c11163c896b09e7b6863d347dadbcb65702585cc
PRISTINE_SHA=3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7
PRISTINE_BLOB=531060a55c080ebd535b426b8c60dcd9011785fc

SHOT="$1"; shift
LOGNAME="$1"; shift
RUN="$1"

LOG=$ROOT/logs/$LOGNAME
mkdir -p "$ROOT/logs"
export PATH="$EXPORT/third_party/sherpa-onnx:$EXPORT/build:$PATH"
export GOTMPDIR=$ROOT/gotmp
cd "$EXPORT" || { echo "cd export FAILED" >"$LOG"; exit 9; }

read_sha() {
  printf 'sha256=%s blob1=%s bytes=%s lines=%s\n' \
    "$(sha256sum "$1" | cut -d' ' -f1)" "$(git hash-object "$1")" "$(wc -c <"$1")" "$(wc -l <"$1")"
}
# 另一枚落点"这一发期间未被改"的尺（该行逐字节哈希）
line_sha() { grep -F "$1" "$2" | sha256sum | cut -d' ' -f1; }

if [ "$SHOT" = "M1" ]; then
  ANCHOR_LINE='	stopSignal := rt.StopSignal()'
  OTHER_TOKEN='rt.stopReq = newStopRequest()'
  MUT_LINE='	stopSignal := (<-chan struct{})(nil) // MUTATION-228r3c-M1 severed-loop-stop-check'
  LABEL='M1 = 断环路 stop 检查（nil 通道永不 ready ⇒ 记录下的请求转不成 return）'
else
  ANCHOR_LINE='	rt.stopReq = newStopRequest()'
  OTHER_TOKEN='stopSignal := rt.StopSignal()'
  MUT_LINE='	// MUTATION-228r3c-M2 boot-builds-no-stop-bit  rt.stopReq = newStopRequest()'
  LABEL='M2 = 令 Boot 不建停机位（StopSignal() 答 nil）'
fi

{
  echo "### ===== 228-r3c mutation shot $SHOT =====  $(date '+%F %T %z')"
  echo "### $LABEL"
  echo "### anchor=$ANCHOR"
  echo "### surface=c（默认 C 卷临时根）GOTMPDIR=$GOTMPDIR pwd=$PWD"
  echo "### target=$TARGET_REL（仓外导出树物理字节；⛔ 动跟踪件、⛔ -overlay）"
  echo "### content-anchor-line=[${ANCHOR_LINE}]（整行等值替换，⛔ 绝对行号）"
  echo "### other-landing-token=[${OTHER_TOKEN}] 本发⛔ 碰"
  echo "### run-pattern=$RUN"
  echo "### --- [1] 本发突变前哈希 ---"
  printf 'HASH-BEFORE-MUTATION '; read_sha "$TARGET"
} >"$LOG" 2>&1

PRE=$(sed -n 's/^HASH-BEFORE-MUTATION //p' "$LOG")
case "$PRE" in
  "sha256=$PRISTINE_SHA blob1=$PRISTINE_BLOB"*) ;;
  *) echo "### ABORT: tree NOT pristine before $SHOT ⇒ ⛔ 开跑" >>"$LOG"; exit 8 ;;
esac
OTHER_PRE=$(line_sha "$OTHER_TOKEN" "$TARGET")

awk -v old="$ANCHOR_LINE" -v new="$MUT_LINE" \
  '$0 == old {print new; hit++; next} {print} END {if (hit != 1) exit 6}' "$TARGET" >"$TARGET.new" \
  && mv -f "$TARGET.new" "$TARGET"
APRC=$?
{
  echo "### --- [2] 突变落笔 rc=$APRC（hit 必须恰 1 次，否则 awk exit 6）---"
  if [ "$APRC" != "0" ]; then echo "### ABORT: anchor not matched exactly once ⇒ ⛔ 继续"; fi
  printf 'HASH-AFTER-MUTATION '; read_sha "$TARGET"
  echo "### mutated-line-readback:"
  grep -n 'MUTATION-228r3c' "$TARGET"
  echo "### 另一枚落点未被改的尺（行哈希前后 ⇄ 整文件哈希见上两行）:"
  printf 'OTHER-ANCHOR-LINE-SHA256-PRE  %s\n' "$OTHER_PRE"
  printf 'OTHER-ANCHOR-LINE-SHA256-POST %s\n' "$(line_sha "$OTHER_TOKEN" "$TARGET")"
  printf 'OTHER-ANCHOR-HITS %s (expect 1)\n' "$(grep -cF "$OTHER_TOKEN" "$TARGET" || true)"
} >>"$LOG" 2>&1
[ "$APRC" = "0" ] || exit 6

go test ./cmd/wisp/ -count=1 -timeout 30m -v -run "$RUN" >>"$LOG" 2>&1
RC_TEST=$?
# 此刻 log 里只有"突变那一发"的输出（复跑段还没落笔），故整件 ^--- FAIL 计数＝该发顶层红数
RED=$(grep -c '^--- FAIL' "$LOG" || true)
NAMED=$(grep -E '^(--- FAIL:|    [A-Za-z0-9_]+\.go:[0-9]+: AC#11 RED)' "$LOG" || true)
{
  echo "### --- [3] 点名红读数  go-test-rc=$RC_TEST（期望 1） $(date '+%F %T %z') ---"
  echo "### named-red-lines:"
  [ -n "$NAMED" ] && printf '%s\n' "$NAMED" || echo "### (NO NAMED RED ⇒ 判据没响，本发不算，交件具名)"
  printf '### NAMED-FAIL-TOP-COUNT=%s\n' "$RED"
} >>"$LOG" 2>&1

{
  echo "### --- [4] 还原＝按具名 blob 重新 git archive 覆盖（⛔ 信 echo）---"
  ( cd "$REPO" && git archive "$ANCHOR" "$TARGET_REL" ) | tar -x -C "$EXPORT"
  echo "### restore-command-rc=$?"
} >>"$LOG" 2>&1

{
  echo "### --- [5] 还原后哈希（只认这把对撞；三段齐＝本发/前/后）---"
  printf 'HASH-AFTER-RESTORE '; read_sha "$TARGET"
} >>"$LOG" 2>&1
POST=$(sed -n 's/^HASH-AFTER-RESTORE //p' "$LOG")
case "$POST" in
  "sha256=$PRISTINE_SHA blob1=$PRISTINE_BLOB"*) ;;
  *) echo "### RESULT: NOT RESTORED ⇒ 停，下一发⛔ 开跑" >>"$LOG"; exit 7 ;;
esac

# 复跑段单独可分：记下当前行数，只 tail 这一段——否则上一发（突变态）的
# --- FAIL 会被 tail 一起带进来＝读数歧义（M1 那发的尺局限已在提交里具名）。
RERUN_FROM=$(( $(wc -l <"$LOG") + 1 ))
go test ./cmd/wisp/ -count=1 -timeout 30m -v -run "$RUN" >>"$LOG" 2>&1
RC_RERUN=$?
RERUN_COLOUR=$(sed -n "${RERUN_FROM},\$p" "$LOG" | grep -E '^(--- (PASS|FAIL):|=== RUN|ok  |FAIL\t)' || true)
RERUN_FAIL=$(sed -n "${RERUN_FROM},\$p" "$LOG" | grep -c '^--- FAIL' || true)
RERUN_PASS=$(sed -n "${RERUN_FROM},\$p" "$LOG" | grep -c '^--- PASS' || true)
{
  echo "### --- [6] 还原后复跑同一用例（期望 rc=0 绿）  go-test-rc=$RC_RERUN  $(date '+%F %T %z') ---"
  echo "### rerun-only-segment (lines ${RERUN_FROM}-end of this file):"
  printf '%s\n' "$RERUN_COLOUR"
  printf '### rerun-own-counts: PASS_TOP=%s FAIL_TOP=%s (expect FAIL_TOP=0)\n' "$RERUN_PASS" "$RERUN_FAIL"
  printf 'HASH-FINAL '; read_sha "$TARGET"
  printf 'OTHER-ANCHOR-LINE-SHA256-FINAL %s (expect == PRE)\n' "$(line_sha "$OTHER_TOKEN" "$TARGET")"
  echo "### ===== DONE shot=$SHOT rc_test=$RC_TEST rc_rerun=$RC_RERUN ====="
} >>"$LOG" 2>&1
echo "DONE shot=$SHOT rc_test=$RC_TEST rc_rerun=$RC_RERUN log=$LOG"
