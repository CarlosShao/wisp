#!/usr/bin/env bash
# 328-v1 AC#3 三面尺＋AC#2① 逐跳 file:line＋AC#2② 任务号半边归属（全只读；颜色一律在仓外树，本件不跑 go）。
set -u
REPO='D:/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 1
BEFORE='648eda04'; AFTER='46d815d9'
# ★固定裁决锚＝env V1ANCHOR（不入 env 时现认 git log -1）。本腿所有只读尺＋那棵 after 树都用同一枚具名锚，
#   免得 HEAD 在飞（本波另五枚腿在落笔）让两把尺量到两份字节。
MYANCHOR="${V1ANCHOR:-$(git log -1 --format=%H)}"
LOGD='.scratch/wisp/probes/328/v1/logs'

echo "date=$(date '+%F %T %z') my_anchor=$MYANCHOR"
echo

echo '## AC#3① 那 4 枚 .go 到底动了什么（尺＝git diff -U0 --numstat ⇄ 逐枚函数体具名）'
for f in internal/panel/subagent_roster_197.go internal/panel/subagent_blocked_220_test.go internal/panel/subagent_corr_task_prefix_328_test.go internal/agent/corr_roster_join_328_test.go; do
  ns=$(git diff -U0 "$BEFORE".."$MYANCHOR" --numstat -- "$f" | awk '{print "+"$1" -"$2}')
  echo "  $f ${ns:-未改动}"
done
echo '--- 名册那枚：新增/改动的标识符逐枚点名（尺＝diff -U0 只取 + 行的 func/for/if） ---'
git diff -U0 "$BEFORE".."$MYANCHOR" -- internal/panel/subagent_roster_197.go | grep '^+' | grep -nE 'func |HasPrefix|splitClashSuffix|set\[|for task' | head -20
echo '--- 数字臂（splitClashSuffix 函数体）对锚零改动那把尺（尺＝取函数体行段 hash 对撞）---'
for A in "$BEFORE" "$MYANCHOR"; do
  body=$(git show "$A:internal/panel/subagent_roster_197.go" | awk '/^func splitClashSuffix/,/^}/' | sha256sum | cut -c1-16)
  echo "  splitClashSuffix_body_sha16 anchor=$A $body"
done
echo '--- 两枚新仪器文件＝改前侧不存在（尺＝cat-file -t）---'
for A in "$BEFORE"; do
  for f in internal/panel/subagent_corr_task_prefix_328_test.go internal/agent/corr_roster_join_328_test.go; do
    echo "  $A $f exists=$(git cat-file -e "$A:$f" 2>/dev/null && echo YES || echo 'NO(新增)')"
  done
done
echo

echo '## AC#3② 门控档位与 L2 许可面零字节（blob 对撞，起手锚 ⇄ 腿末笔 ⇄ 我的锚）'
for path in internal/risk/mode.go internal/risk/provenance.go internal/agent/approval/gate.go internal/agent/approval/queue.go; do
  b1=$(git rev-parse "$BEFORE:$path" 2>/dev/null); b2=$(git rev-parse "$AFTER:$path" 2>/dev/null); b3=$(git rev-parse "$MYANCHOR:$path" 2>/dev/null)
  v=ZEROBYTE; [ "$b1" != "$b3" ] && v=CHANGED
  echo "  $path before=$b1 legEnd=$b2 myAnchor=$b3 => $v"
done
echo '--- internal/risk/** 整面行数尺 ---'
echo "  git diff -U0 $BEFORE..$MYANCHOR -- internal/risk | wc -l = $(git diff -U0 "$BEFORE".."$MYANCHOR" -- internal/risk | wc -l)"
echo "  git diff -U0 $BEFORE..$MYANCHOR -- internal/agent/approval | wc -l = $(git diff -U0 "$BEFORE".."$MYANCHOR" -- internal/agent/approval | wc -l)"
echo "  git diff -U0 $BEFORE..$MYANCHOR -- internal/statemachine | wc -l = $(git diff -U0 "$BEFORE".."$MYANCHOR" -- internal/statemachine | wc -l)"
echo

echo '## AC#3③ 冻结件／thresholds／golden／go.mod／go.sum／ci.yml 命中 0（带正控）'
FROZEN='internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go'
for f in $FROZEN; do
  n=$(git diff "$BEFORE".."$MYANCHOR" -- "$f" | wc -l)
  echo "  diff_lines $f = $n"
done
echo '--- 正控（先证我取出的东西非空再信 0）：三枚冻结件在我的锚里真存在且非空 ---'
for f in $FROZEN; do
  echo "  $f bytes=$(git cat-file -s "$MYANCHOR:$f" 2>/dev/null) lines=$(git show "$MYANCHOR:$f" | wc -l)"
done
echo '--- 面值尺（thresholds.go／golden／go.mod／go.sum／ci.yml）---'
for face in 'internal/**/thresholds.go' 'testdata/**' go.mod go.sum .github/workflows/ci.yml; do
  echo "  diff_lines($face)=$(git diff "$BEFORE".."$MYANCHOR" -- "$face" | wc -l)"
done
echo '--- 正控：thresholds.go 与 golden 名册在我的锚非空（尺＝ls-tree）---'
git ls-tree -r --name-only "$MYANCHOR" | grep -E 'thresholds\.go$' | head -5
echo "  golden_files_count=$(git ls-tree -r --name-only "$MYANCHOR" | grep -c 'golden' || true)"
echo "  ci_yml_bytes=$(git cat-file -s "$MYANCHOR:.github/workflows/ci.yml" 2>/dev/null || echo NOFILE)"
echo

echo '## AC#2① 正控逐跳 file:line（现跑，行号只当读数；射程＝我的锚具名提交树）'
echo '--- 跳1 铸造点 def＋唯一非测试调用 ---'
git grep -n 'func callCorr' "$MYANCHOR" -- internal/agent
git grep -n 'callCorr(taskID' "$MYANCHOR" -- internal/agent/loop.go
echo '--- 跳2 ToolRequest 生产构造点（唯一一枚非注释）---'
git grep -n 'ToolRequest{' "$MYANCHOR" -- internal cmd | grep -v '_test.go'
echo '--- 跳3 executeCalls ← run 的父链 ---'
git grep -n 'func (l \*Loop) executeCalls\|executeCalls(' "$MYANCHOR" -- internal/agent/loop.go
git grep -n 'func (l \*Loop) run' "$MYANCHOR" -- internal/agent/loop.go
echo '--- 跳4 卡片搬运 NativeVerdict.CardView → approval.go ---'
git grep -n 'func (v NativeVerdict) CardView\|CardViewFromDecision' "$MYANCHOR" -- internal/panel | head -5
echo '--- 跳5 名册渲染入口 TaskRosterSectionFrom def ⇄ 宿主真调用点 ---'
git grep -n 'func TaskRosterSectionFrom\|func taskRosterSectionFrom' "$MYANCHOR" -- internal/panel
git grep -n 'TaskRosterSectionFrom(' "$MYANCHOR" -- cmd/wisp internal/panel/pump.go | grep -v '_test'
echo '--- 跳6 第三臂本体＋判据行 ---'
git grep -n 'HasPrefix(key, task+"#")' "$MYANCHOR" -- internal/panel/subagent_roster_197.go
git grep -n 'func rosteredTaskIDs' "$MYANCHOR" -- internal/panel/subagent_roster_197.go
echo

echo '## AC#2② 任务号半边（T）是谁给的——三问逐问现跑'
echo '--- 问A：正控里喂给名册的 TaskID 取值来自哪里（夹具文件逐字节，射程＝该文件）---'
git show "$MYANCHOR:internal/agent/corr_roster_join_328_test.go" | grep -n 'res.TaskID\|rosterWithTwoRows\|sibling-of\|TaskID:' | head -12
echo '--- 问B：真产 TaskID 的那枚函数（newTaskID）---'
git grep -n 'func newTaskID\|TaskID:  *= *newTaskID\|taskID := newTaskID' "$MYANCHOR" -- internal/agent | head -5
echo '--- 问C：生产侧那枚名册装配者（真把行塞进 TaskRosterState 的代码在哪家）---'
git grep -n 'TaskRosterState{' "$MYANCHOR" -- internal cmd | grep -v '_test.go' | head -10
echo '--- 问C补：TaskRosterState 的字段握法＋在册那把尺的入参来源 ---'
git grep -n 'func rosteredTaskIDs' -A6 "$MYANCHOR" -- internal/panel/subagent_roster_197.go | sed 's/^/    /'
git show "$MYANCHOR:internal/panel/subagent_roster_197.go" | awk '/^func rosteredTaskIDs/,/^}/' | sed 's/^/    /'
echo

echo '## 必答①那格的零读数尺（叠改形真产过没有）＋正控'
echo "  尺A git grep -nE '\"[^\"]*#[0-9]+\"' $MYANCHOR -- internal/agent/approval | wc -l = $(git grep -nE '"[^"]*#[0-9]+"' "$MYANCHOR" -- internal/agent/approval | wc -l)"
echo "  正控 尺B 同族放宽到 internal cmd：$(git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' "$MYANCHOR" -- internal cmd | wc -l) 行"
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' "$MYANCHOR" -- internal cmd | sed 's/\(.\{100\}\).*/\1…/' | head -15
echo '--- 队列改写那半边（真产 <corr>#<seq> 的代码）逐枚点名 ---'
git grep -n 'seq\b' "$MYANCHOR" -- internal/agent/approval/queue.go | grep -E '#|fmt\.Sprintf' | head -8
echo '--- 两枚铸造点的可导性（要不要新导出名才链得起来）---'
git grep -n 'func (q \*Queue) push\|func callCorr' "$MYANCHOR" -- internal/agent internal/agent/approval
echo

echo '## 必答② 用户可见那一跳：frontend 命中枚数（带尺＋射程＋正控）'
for r in frontend frontend/src; do
  echo "  git grep -l 'blockedOnApproval' $MYANCHOR -- $r | wc -l = $(git grep -l 'blockedOnApproval' "$MYANCHOR" -- "$r" | wc -l)"
  echo "  git grep -il 'blockedonapproval' $MYANCHOR -- $r | wc -l = $(git grep -il 'blockedonapproval' "$MYANCHOR" -- "$r" | wc -l)"
done
echo '--- 正控（同尺族 correlationId；量的是文件枚数）---'
echo "  correlationId -- frontend = $(git grep -l 'correlationId' "$MYANCHOR" -- frontend | wc -l) 枚文件"
echo "  correlationId -- frontend/src = $(git grep -l 'correlationId' "$MYANCHOR" -- frontend/src | wc -l) 枚文件 / $(git grep -n 'correlationId' "$MYANCHOR" -- frontend/src | wc -l) 处出现"
echo '--- 接收器自陈（现量行号）---'
git grep -n 'wispOnPanelSnapshot' "$MYANCHOR" -- cmd/wisp | head -6
git grep -n 'blockedOnApproval is never asserted' "$MYANCHOR" -- cmd/wisp
