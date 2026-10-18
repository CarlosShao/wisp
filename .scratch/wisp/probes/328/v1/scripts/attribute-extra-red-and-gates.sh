#!/usr/bin/env bash
# 328-v1 第二发台面：① 改后多出的那枚红名归位（是不是 328 的字节造成的）② 精确定点（四形＋票 220 那族逐名）
# ③ 门禁四数（共享工作树）。全部仓外树跑 go test；⛔ ./... ⛔ ./cmd/wisp。
set -u
BASE='D:/tmp/wisp/20261018/328v1'
REPO='D:/work/workspace/projects plans/Wisp'
LOGD="$BASE/logs"
export GOTMPDIR="$BASE/gotmp" GOPATH='D:\work\base\gopath'
V1ANCHOR="${V1ANCHOR:?需要裁决锚}"
ROSTER='internal/panel/subagent_roster_197.go'

fours() { local L="$1" G="$2" RC="$3"
  echo "$L: RUN=$(grep -c '^=== RUN' "$G" 2>/dev/null || true) PASS_TOP=$(grep -c '^--- PASS' "$G" 2>/dev/null || true) FAIL_TOP=$(grep -c '^--- FAIL' "$G" 2>/dev/null || true) SKIP=$(grep -c '^--- SKIP' "$G" 2>/dev/null || true) rc=$RC"; }
go_t() { local T="$1" P="$2"; shift 2; ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -v "./$P/" -count=1 -timeout 30m "$@" ); }

echo "date=$(date '+%F %T %z') adjudication_anchor=$V1ANCHOR"
echo "tree_roster_blob=$(cd "$BASE/after" && git hash-object "$ROSTER") anchor_blob=$(git -C "$REPO" rev-parse "$V1ANCHOR:$ROSTER")"
echo "legEnd_f705b082_blob=$(git -C "$REPO" rev-parse "f705b082:$ROSTER")"
echo

echo '## A. 那枚多出来的红名：文件归属尺（它是不是 328 射程）'
for f in internal/panel/backpressure_bound_35r8_test.go; do
  echo "  diff_lines $f 648eda04..$V1ANCHOR = $(git -C "$REPO" diff "648eda04".."$V1ANCHOR" -- "$f" | wc -l)"
  echo "  blob before=$(git -C "$REPO" rev-parse "648eda04:$f" 2>/dev/null) after=$V1ANCHOR=$(git -C "$REPO" rev-parse "$V1ANCHOR:$f" 2>/dev/null)"
done
echo '  328 落笔四枚文件名册（尺＝git show --name-only 逐枚）：'
for c in b2d55342 1696d01b ffebb6d1 f705b082; do echo "    $c => $(git -C "$REPO" show --name-only --format='' "$c" | tr '\n' ' ')"; done
echo

echo '## B. 归位实验一：改前侧整包再跑第二发（同树同配方；那枚红名是否本来就会自己冒出来）'
go_t "$BASE/before" internal/panel > "$LOGD/before-panel-full-rerun2.log" 2>&1; rc=$?
fours 'before-panel-full-rerun2' "$LOGD/before-panel-full-rerun2.log" "$rc"
grep '^--- FAIL' "$LOGD/before-panel-full-rerun2.log" | sed 's/ ([0-9.]*s)//' | sort -u > "$LOGD/red-before-panel-rerun2.txt"
echo '  改前第二发红名册（⇄ 改前第一发）：'; diff "$LOGD/red-before-panel.txt" "$LOGD/red-before-panel-rerun2.txt" && echo '  SAME_AS_FIRST_BEFORE'
echo '  改前第二发 ⇄ 改后第一发 逐名作差：'; diff "$LOGD/red-before-panel-rerun2.txt" "$LOGD/red-after-panel.txt" || true
echo

echo '## C. 归位实验二：那枚红名单跑（两侧各两发；堆量尺在整包里受同场分配影响）'
for side in before after; do
  for i in 1 2; do
    go_t "$BASE/$side" internal/panel -run 'TestFloodRulerFiresOnAnUnboundedSink35r8' > "$LOGD/floodruler-$side-$i.log" 2>&1; rc=$?
    fours "floodruler-$side-$i" "$LOGD/floodruler-$side-$i.log" "$rc"
    grep -E 'backpressure_bound_35r8_test.go:(3[0-9][0-9])' "$LOGD/floodruler-$side-$i.log" | head -3
  done
done
echo

echo '## D. 精确定点：四形实测（改后侧，逐名 RUN/结果）'
go_t "$BASE/after" internal/panel -run 'TestProductionCorrelationShapesLightTheirRows|TestEachCorrelationLightsOnlyTheRowItIsBuiltOn|TestAnL1WindowCorrelationIsMatchedByTheSameArm|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestASecondCardForTheSameTaskStillLightsItsRow' > "$LOGD/after-panel-fourshapes-precise.log" 2>&1; rc=$?
fours 'after-panel-fourshapes-precise' "$LOGD/after-panel-fourshapes-precise.log" "$rc"
grep -E '^(--- |    --- )(PASS|FAIL)' "$LOGD/after-panel-fourshapes-precise.log" | sed 's/ ([0-9.]*s)//'
echo '  --- 红句逐字（若有）---'
grep -E '^\s+\S+_test\.go:[0-9]+:' "$LOGD/after-panel-fourshapes-precise.log" | head -8
echo

echo '## E. 改前侧同一条 -run（ⓑ 那枚断言的起手形状：child-1#retry 当时钉 false）'
go_t "$BASE/before" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestASecondCardForTheSameTaskStillLightsItsRow' > "$LOGD/before-panel-220table-precise.log" 2>&1; rc=$?
fours 'before-panel-220table-precise' "$LOGD/before-panel-220table-precise.log" "$rc"
grep -E '^(--- |    --- )(PASS|FAIL)' "$LOGD/before-panel-220table-precise.log" | sed 's/ ([0-9.]*s)//'
echo

echo '## F. 门禁四数（只在共享工作树；我零 .go，尺打的是腿那四枚文件＋两包）'
cd "$REPO" || exit 1
/d/work/base/gopath/bin/gofumpt.exe -l internal/panel/subagent_roster_197.go internal/panel/subagent_blocked_220_test.go internal/panel/subagent_corr_task_prefix_328_test.go internal/agent/corr_roster_join_328_test.go; echo "gofumpt_rc=$?"
sh scripts/d22scan.sh 2>&1 | tail -6; echo "d22scan_rc=${PIPESTATUS[0]}"
sh scripts/check-path-length-budget.sh 2>&1 | tail -6; echo "length_rc=${PIPESTATUS[0]}"
go vet ./internal/panel/ ./internal/agent/ 2>&1 | tail -8; echo "vet_rc=${PIPESTATUS[0]}"
echo '--- 共享树那两包字节与裁决锚一致性（跑门禁前现量，别腿若已落笔要具名）---'
git status --porcelain -- internal/panel internal/agent cmd/wisp; echo "status_rc=$?"
echo "shared_tree_vs_anchor_internal_panel_agent_diff_files=$(git diff --name-only $V1ANCHOR -- internal/panel internal/agent | wc -l)"
echo "done_rc=0"
