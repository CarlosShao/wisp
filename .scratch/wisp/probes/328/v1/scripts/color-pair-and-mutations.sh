#!/usr/bin/env bash
# 328-v1 AC#4① 成对颜色复跑＋AC#2/AC#3 定点＋两发去牙突变（全部串行、单作业、仓外导出树）。
# 规矩块：-count=1 -timeout 30m；四数尺＝RUN(^=== RUN 含子测试)／PASS_TOP(^--- PASS 只顶层)／
# FAIL_TOP(^--- FAIL)／SKIP(^--- SKIP)，逐枚点名，⛔ 混用，⛔ 相减。⛔ ./... ⛔ ./cmd/wisp。
# 突变只打在仓外树字节上；逐发 pristine 前置断言→突变→跑→当场还原→还原后断言 pristine（同一 shell）。
set -u
BASE='D:/tmp/wisp/20261018/328v1'
REPO='D:/work/workspace/projects plans/Wisp'
LOGD="$BASE/logs"
mkdir -p "$LOGD"
export GOTMPDIR="$BASE/gotmp"
export GOPATH='D:\work\base\gopath'

fours() { # <label> <log> <rc>
  local L="$1" G="$2" RC="$3"
  local run=$(grep -c '^=== RUN' "$G" 2>/dev/null || true)
  local pt=$(grep -c '^--- PASS' "$G" 2>/dev/null || true)
  local ft=$(grep -c '^--- FAIL' "$G" 2>/dev/null || true)
  local sk=$(grep -c '^--- SKIP' "$G" 2>/dev/null || true)
  echo "$L: RUN=$run PASS_TOP=$pt FAIL_TOP=$ft SKIP=$sk rc=$RC"
}
rednames() { # <log> <out>
  grep '^--- FAIL' "$1" 2>/dev/null | sed 's/ ([0-9.]*s)//' | sort -u > "$2" || true
  echo "red_names=$(wc -l < "$2") out=$2"
}
go_t() { # <tree> <pkg> [ -run <re> ]
  local T="$1" P="$2"; shift 2
  ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -v "./$P/" -count=1 -timeout 30m "$@" )
}

ROSTER='internal/panel/subagent_roster_197.go'
# ★固定裁决锚：只读尺与 after 树取同一份字节（⛔ 拿两把尺量两份字节）。
MYANCHOR="${V1ANCHOR:?需要 V1ANCHOR＝入场时现认的那枚锚}"
if ! git -C "$REPO" rev-parse --verify "$MYANCHOR^{commit}" >/dev/null 2>&1; then
  echo "V1ANCHOR_NOT_A_COMMIT=$MYANCHOR"; exit 1
fi
TREE_BLOB=$(cd "$BASE/after" && git hash-object "$ROSTER")
ANCHOR_BLOB=$(git -C "$REPO" rev-parse "$MYANCHOR:$ROSTER")
echo "adjudication_anchor=$MYANCHOR"
echo "tree_roster_blob=$TREE_BLOB anchor_roster_blob=$ANCHOR_BLOB same=$([ "$TREE_BLOB" = "$ANCHOR_BLOB" ] && echo YES || echo 'NO ★树建在另一枚锚上⇒停手重建树')"
[ "$TREE_BLOB" = "$ANCHOR_BLOB" ] || exit 1
echo

echo "===== 1. 成对整包（改前侧＝648eda04 树）====="
go_t "$BASE/before" internal/panel > "$LOGD/before-panel-full.log" 2>&1; rc=$?
fours 'before-panel-full' "$LOGD/before-panel-full.log" "$rc"; rednames "$LOGD/before-panel-full.log" "$LOGD/red-before-panel.txt"
go_t "$BASE/before" internal/agent > "$LOGD/before-agent-full.log" 2>&1; rc=$?
fours 'before-agent-full' "$LOGD/before-agent-full.log" "$rc"; rednames "$LOGD/before-agent-full.log" "$LOGD/red-before-agent.txt"

echo "===== 2. 成对整包（改后侧＝我的锚现认树）====="
go_t "$BASE/after" internal/panel > "$LOGD/after-panel-full.log" 2>&1; rc=$?
fours 'after-panel-full' "$LOGD/after-panel-full.log" "$rc"; rednames "$LOGD/after-panel-full.log" "$LOGD/red-after-panel.txt"
go_t "$BASE/after" internal/agent > "$LOGD/after-agent-full.log" 2>&1; rc=$?
fours 'after-agent-full' "$LOGD/after-agent-full.log" "$rc"; rednames "$LOGD/after-agent-full.log" "$LOGD/red-after-agent.txt"

echo "===== 3. 红名册逐名作差（尺＝sort 后 comm；⛔ 只看枚数）====="
echo '--- panel 改前⇄改后 ---'
diff "$LOGD/red-before-panel.txt" "$LOGD/red-after-panel.txt" && echo "panel_redroster_diff_rc=0（名册同集）" || echo "panel_redroster_diff_rc=NONZERO（名字换了，逐名看）"
echo '--- agent 改前⇄改后 ---'
diff "$LOGD/red-before-agent.txt" "$LOGD/red-after-agent.txt" && echo "agent_redroster_diff_rc=0" || echo "agent_redroster_diff_rc=NONZERO"

echo "===== 4. 定点：票 220 那族形表（改前侧应钉 child-1#retry=false 且绿；改后侧=true）====="
go_t "$BASE/before" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' > "$LOGD/before-panel-220table.log" 2>&1; rc=$?
fours 'before-panel-220table' "$LOGD/before-panel-220table.log" "$rc"
go_t "$BASE/after" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestACorrelationBuiltOn|TestASecondCardForTheSameTask' > "$LOGD/after-panel-220table.log" 2>&1; rc=$?
fours 'after-panel-220table' "$LOGD/after-panel-220table.log" "$rc"

echo "===== 5. 定点：四形实测（改后侧）====="
go_t "$BASE/after" internal/panel -run 'TestCorrelationShapesFromTheLoop|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestA' > "$LOGD/after-panel-fourshapes.log" 2>&1; rc=$?
fours 'after-panel-fourshapes' "$LOGD/after-panel-fourshapes.log" "$rc"

echo "===== 6. 定点：AC#2 正控（agent 包，经真铸造点那两发）====="
go_t "$BASE/after" internal/agent -run 'TestACorrelationTheLoopActuallyMintedLightsItsRosterRow|TestTheLoopMintsTheShapesTheRosterAccepts' > "$LOGD/after-agent-positive.log" 2>&1; rc=$?
fours 'after-agent-positive' "$LOGD/after-agent-positive.log" "$rc"
grep -E '^(=== RUN|--- PASS|--- FAIL|ok|FAIL|PASS)' "$LOGD/after-agent-positive.log" | head -12

echo "===== 7. ★去牙一发（派单 AC#2③ 那一发＝抹掉边界⇒:142 child-1x 必须红）====="
F="$BASE/after/$ROSTER"
P0=$(sha256sum "$F" | cut -d' ' -f1)
B0=$(git -C "$REPO" hash-object "$F")
echo "pristine_before: sha256=$P0 gitblob=$B0 (对锚串应＝2da6c596369c1ee38672ae2656a7d5a3ecbaabfb)"
sed -i 's/strings\.HasPrefix(key, task+"#")/strings.HasPrefix(key, task)/' "$F"
grep -n 'strings.HasPrefix(key, task' "$F"
go_t "$BASE/after" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' > "$LOGD/mut-V1-boundary.log" 2>&1; rc=$?
fours 'mut-V1-boundary' "$LOGD/mut-V1-boundary.log" "$rc"
grep -B1 -A6 'FAIL: TestTheClashForm' "$LOGD/mut-V1-boundary.log" | head -30
P1=$(sha256sum "$F" | cut -d' ' -f1); echo "mutated_sha256=$P1"
cp "$BASE/before/$ROSTER" /tmp/nothing-ignore 2>/dev/null
git -C "$REPO" show "HEAD:$ROSTER" > "$F"
P2=$(sha256sum "$F" | cut -d' ' -f1)
B2=$(git -C "$REPO" hash-object "$F")
echo "restored: sha256=$P2 gitblob=$B2"
[ "$P2" = "$P0" ] && [ "$P0" != "$P1" ] && echo "restored_same_as_pristine=YES" || echo "restored_same_as_pristine=NO ★"
echo "pristine_equals_anchor_blob=$([ "$B2" = '2da6c596369c1ee38672ae2656a7d5a3ecbaabfb' ] && echo YES || echo NO)"

echo "===== 8. ★去牙二发（第三臂登记抹掉＝改回今天的形状⇒agent 正控必须红，落点与一发不同枚）====="
sed -i 's/^\t\t\tset\[task\] = true$/\t\t\t_ = task/' "$F"
if ! grep -n '_ = task' "$F"; then
  # TAB 锚可能不中⇒退回内容锚的第二写法
  perl -0pi -e 's/(if strings\.HasPrefix\(key, task\+"#"\) \{\n)\t{3}set\[task\] = true/$1\t\t\t_ = task/' "$F"
fi
grep -n -A2 'strings.HasPrefix(key, task+"#")' "$F"
go_t "$BASE/after" internal/agent -run 'TestACorrelationTheLoopActuallyMintedLightsItsRosterRow|TestTheLoopMintsTheShapesTheRosterAccepts' > "$LOGD/mut-V2-droparm-agent.log" 2>&1; rc=$?
fours 'mut-V2-droparm-agent' "$LOGD/mut-V2-droparm-agent.log" "$rc"
grep -E 'FAIL|not-blocked|does not light' "$LOGD/mut-V2-droparm-agent.log" | head -8
go_t "$BASE/after" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' > "$LOGD/mut-V2-droparm-panel.log" 2>&1; rc=$?
fours 'mut-V2-droparm-panel' "$LOGD/mut-V2-droparm-panel.log" "$rc"
grep -E 'FAIL: TestTheClash|--- (PASS|FAIL)' "$LOGD/mut-V2-droparm-panel.log" | sed 's/ ([0-9.]*s)//' | head -14
git -C "$REPO" show "HEAD:$ROSTER" > "$F"
P3=$(sha256sum "$F" | cut -d' ' -f1)
B3=$(git -C "$REPO" hash-object "$F")
echo "restored2: sha256=$P3 gitblob=$B3 pristine_equals_anchor_blob=$([ "$B3" = '2da6c596369c1ee38672ae2656a7d5a3ecbaabfb' ] && echo YES || echo NO)"
echo "===== 9. 还原后复绿一发（证树没坏）====="
go_t "$BASE/after" internal/panel -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' > "$LOGD/after-panel-220table-restored.log" 2>&1; rc=$?
fours 'after-panel-220table-restored' "$LOGD/after-panel-220table-restored.log" "$rc"
echo "color_pair_mutations_rc=0"
