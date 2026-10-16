#!/bin/sh
# 165-v2 裁决腿突变器：把一枚形打进去 -> 跑 scoped 那一发 -> 落日志 -> 当场还原 -> 报两枚 sha256。
# 台面＝仓外导出树 /d/tmp/165v2/t1（git archive 44f51fc8 八件套＋build/third_party dll＋frontend/dist）。
# ⛔ 碰共享工作树的 internal/**／cmd/** 字节：本脚本只读 $BASE 两枚 pristine，只写 $TREE。
TREE=/d/tmp/165v2/t1
LOGS=/d/tmp/165v2/logs
BASE=/d/tmp/165v2/pristine
SSE_REL=internal/agent/testdata/golden/plan-present.sse
GO_REL=internal/agent/plan_present_doorbell_165r3_test.go
RUN='^TestAC|^TestPlanPresent'

# 用法：sh mut.sh <形名> <sed脚本> <被sed的文件：sse|go>
name="$1"; script="$2"; which="${3:-go}"
tgt="$SSE_REL"; [ "$which" = "go" ] && tgt="$GO_REL"

cp "$BASE/plan-present.sse.base" "$TREE/$SSE_REL"
cp "$BASE/doorbell.go.base"     "$TREE/$GO_REL"
sed -i "$script" "$TREE/$tgt" || { echo "sed rc=$?"; exit 9; }
echo "== 形 $name ｜ sed '$script' ｜ 靶 $(basename $tgt) =="
wc -l "$TREE/$tgt"

cd "$TREE" || exit 9
GOTMPDIR='D:\tmp\165v2\gotmp' go test ./internal/agent/ -run "$RUN" -count=1 -timeout 30m -v \
  > "$LOGS/mut-$name.txt" 2>&1
rc=$?
echo "rc=$rc  RUN=$(grep -c '^=== RUN' $LOGS/mut-$name.txt) PASS_TOP=$(grep -c '^--- PASS' $LOGS/mut-$name.txt) FAIL_TOP=$(grep -c '^--- FAIL' $LOGS/mut-$name.txt)"
echo "-- 红句（逐字）--"
grep -n 'plan_present_doorbell_165r3_test.go:[0-9]*:' "$LOGS/mut-$name.txt" | head -20
grep -n 'plan-present.sse\|load golden' "$LOGS/mut-$name.txt" | head -5

# 还原＋自证
cp "$BASE/plan-present.sse.base" "$TREE/$SSE_REL"
cp "$BASE/doorbell.go.base"     "$TREE/$GO_REL"
echo "-- 还原后 sha256（两枚必回基线）--"
sha256sum "$TREE/$SSE_REL" "$TREE/$GO_REL" "$TREE/$tgt"
