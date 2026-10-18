#!/usr/bin/env bash
# 328-v1 第三发只读尺：① 上一发那个空读数的追因（TaskRosterState 生产构造点）② 零新导出名／零新 JSON 键
# ③ 正控那一条链的逐跳 file:line（含在册任务号那半边的归属）④ leg 四枚文件的 blob 逐枚对撞。
set -u
REPO='D:/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 1
V1ANCHOR="${V1ANCHOR:?需要裁决锚}"
BEFORE='648eda04'
echo "date=$(date '+%F %T %z') adjudication_anchor=$V1ANCHOR before=$BEFORE"
echo

echo '## 1. 追因：上一发问C 那把尺回空——先看 stderr 再判（规矩块读数纪律）'
echo '$ git grep -n "TaskRosterState{" <锚> -- internal cmd   （全含测试）'
git grep -n 'TaskRosterState{' "$V1ANCHOR" -- internal cmd > /tmp/g1.txt 2>/tmp/e1.txt; echo "  rc=$? hits=$(wc -l < /tmp/g1.txt) stderr_bytes=$(wc -c < /tmp/e1.txt)"
sed 's/^[0-9a-f]\{40\}://' /tmp/g1.txt | head -20
echo '  ⇒ 上一发我把 grep -v _test.go 直接接在同一根管子上并只打非测试 ⇒ 需逐枚点名：'
echo '$ 非测试构造点（量的）：'
awk -F: '/_test\.go/ {next} {print}' /tmp/g1.txt | sed 's/^[0-9a-f]\{40\}://' | head -20
echo "  非测试命中枚数=$(awk '/_test\.go/ {next} {c++} END{print c+0}' /tmp/g1.txt)"
echo '$ 那么名册那批行在生产是谁塞进去的（TaskRosterState 作参数/字段的所有者，非测试）'
git grep -n 'TaskRosterState' "$V1ANCHOR" -- internal cmd | grep -v '_test\.go' | sed 's/^[^:]*:[^:]*://' | head -20
echo

echo '## 2. ★在册那半边（任务号）的产码链：cmd/wisp 侧真装配者（现量行号）'
echo '$ git grep -n "TaskRosterSectionFrom(" 非测试'
git grep -n 'TaskRosterSectionFrom(' "$V1ANCHOR" -- internal cmd | grep -v '_test\.go' | sed 's/^[^:]*:[^:]*://;s/:[0-9]*:/: /' | head
echo '$ roster 行的 TaskID 在宿主侧从哪来（TaskRosterState 装配那枚函数）'
git grep -n 'rosterState\|TaskRosterState{\|Rows: *\[\]\|tasks\.Look' "$V1ANCHOR" -- cmd/wisp | grep -v '_test\.go' | sed 's/^[^:]*:[^:]*://' | head -20
echo

echo '## 3. 零新导出名／零新 JSON 键（尺＝对撞两枚锚的标识符名册）'
for f in internal/panel/subagent_roster_197.go internal/panel/subagent_blocked_220_test.go internal/panel/subagent_corr_task_prefix_328_test.go internal/agent/corr_roster_join_328_test.go; do
  echo "  blob $f before=$(git rev-parse "$BEFORE:$f" 2>/dev/null || echo NOFILE) after=$(git rev-parse "$V1ANCHOR:$f")"
done
echo '--- 导出函数名册作差（panel 那枚产码文件）---'
git show "$BEFORE:internal/panel/subagent_roster_197.go" | grep -oE '^func [A-Z][A-Za-z0-9_]*' | sort -u > /tmp/exp-b.txt
git show "$V1ANCHOR:internal/panel/subagent_roster_197.go" | grep -oE '^func [A-Z][A-Za-z0-9_]*' | sort -u > /tmp/exp-a.txt
diff /tmp/exp-b.txt /tmp/exp-a.txt && echo '  exported_func_diff=0（同集）'
echo '--- 全仓新增导出标识符尺（四枚文件 + 行）：diff 的 + 行里抓 func/type/const 顶层导出名 ---'
git diff "$BEFORE".."$V1ANCHOR" -- internal/panel internal/agent | grep '^+' | grep -E '^\+\s*(func|type|const|var) [A-Z]' | sed 's/^/    /'
echo '--- JSON 键名册作差（四枚文件里所有 json: 标签）---'
git show "$BEFORE:internal/panel/subagent_roster_197.go" | grep -oE 'json:"[^"]*"' | sort -u > /tmp/jk-b.txt
git show "$V1ANCHOR:internal/panel/subagent_roster_197.go" | grep -oE 'json:"[^"]*"' | sort -u > /tmp/jk-a.txt
diff /tmp/jk-b.txt /tmp/jk-a.txt && echo '  json_keys_diff=0（同集）'
echo '--- approval.go 卡片字段名册（腿那句 grep -c TaskID=0 的复认）---'
echo "  ApprovalCardView 里 TaskID 命中=$(git show "$V1ANCHOR:internal/panel/approval.go" | grep -c 'TaskID' || true)（尺＝同一枚文件逐枚数）"
git show "$V1ANCHOR:internal/panel/approval.go" | grep -nE 'json:"' | wc -l
echo

echo '## 4. ⓭ 那句自陈是否改到说实话（尺＝现取那两行段逐字）'
git show "$V1ANCHOR:internal/panel/subagent_roster_197.go" | sed -n '184,200p'
echo '  --- splitClashSuffix 文档段 ---'
git show "$V1ANCHOR:internal/panel/subagent_roster_197.go" | sed -n '316,332p'
echo

echo '## 5. 腿件①里那句"唯一读者 :308"⇄现量（splitClashSuffix 直呼尺）'
git grep -n 'splitClashSuffix' "$V1ANCHOR" -- internal cmd | sed 's/^[^:]*:[^:]*://'
echo "  测试直呼枚数=$(git grep -n 'splitClashSuffix' "$V1ANCHOR" -- internal cmd | grep -c '_test.go' || true)"
echo "  正控＝同尺族 indexWaitingKey 命中 $(git grep -n 'indexWaitingKey' "$V1ANCHOR" -- internal cmd | wc -l) 行"
echo

echo '## 6. 票面 AC 框数⇄leg 那句"七格"复认（同一把尺，量的）'
grep -n '^[[:space:]]*- \[' .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md | sed 's/\(.\{40\}\).*/\1…/'
echo "  五枚 AC 框＝AC#0 AC#1 AC#2 AC#3 AC#4；:50/:59 那两条是 bullet 不是框（grep 形状同族被计入＝boxes_total=7 的来源，具名分清）"
echo "  AC 标题行数（尺＝grep -c 'AC#[0-4] 现量普查\|AC#1 选一形\|AC#2 ★正控\|AC#3 判定与批准\|AC#4 门禁四数'）=$(grep -cE '\*\*AC#[0-4] ' .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md)"
