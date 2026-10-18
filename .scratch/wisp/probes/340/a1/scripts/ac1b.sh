#!/bin/sh
# 票 340 普查腿 340-a1 · AC#0 c 承重名册原料＋AC#1 决定性那一跳的现量（只读，具名提交树）
# 逐字重跑：sh .scratch/wisp/probes/340/a1/scripts/ac1b.sh
A=feac8f77c2fbb15addbc36fb1bbaf13ed56588be

echo '### 1. 交集硬尺：同一枚测试文件里既有真 callCorr 侧调用形又有真队列侧调用形'
echo '-- 支 L（真铸造侧的调用形，⛔ 注释词）'
git grep -lE "callCorr\(|agent\.New\(|h\.run\(|loop\.Run\(" "$A" -- internal cmd tools | grep '_test\.go$' | sed "s|^$A:||"
echo '-- 支 P（真队列侧的调用形，⛔ 注释词）'
git grep -lE "q\.push\(|\.PendingApproval\(|g\.Queue()\.LiveApprovals()|gate\.Replay\(" "$A" -- internal cmd tools | grep '_test\.go$' | sed "s|^$A:||"
echo '-- 两支同文件交集（尺＝comm，⛔ 空交集当"尺没响"）'
comm -12 \
  <(git grep -lE "callCorr\(|agent\.New\(|h\.run\(|loop\.Run\(" "$A" -- internal cmd tools | grep '_test\.go$' | sed "s|^$A:||" | sort -u) \
  <(git grep -lE "q\.push\(|\.PendingApproval\(|g\.Queue\(\)\.LiveApprovals\(\)|gate\.Replay\(|approval\.New\(" "$A" -- internal cmd tools | grep '_test\.go$' | sed "s|^$A:||" | sort -u)
echo '   => 交集枚数即"真 loop 与真队列同在一枚用例里"的文件枚数（本腿据此判端到端通路是否已在盘上）'
echo '-- 正控（这把尺会响）：队列侧命中文件名册逐枚'
git grep -lE "q\.push\(|\.PendingApproval\(|g\.Queue\(\)\.LiveApprovals\(\)|gate\.Replay\(|approval\.New\(" "$A" -- internal cmd tools | grep '_test\.go$'

echo
echo '### 2. 决定性那一跳：approval 包内既有测试直呼非导出 push 的枚数（甲形的盘上现状）'
echo '-- 包名册（package approval ⇄ package approval_test）'
for f in $(git ls-tree --name-only "$A:internal/agent/approval" | grep "_test\.go$"); do
  printf '%s  ' "$f"; git show "$A:internal/agent/approval/$f" | head -1
done
echo '-- 直呼 q.push( 的文件及其包名（逐枚）'
for f in $(git grep -l "q\.push(" "$A" -- internal/agent/approval | sed "s/^$A://"); do
  printf '%s  ' "$f"; git show "$A:$f" | head -1
done
echo '-- internal/agent 侧：23 枚测试全 in-package ⇒ callCorr 今天是否已被测试直呼'
git grep -n "callCorr(" "$A" -- internal/agent | grep '_test\.go'

echo
echo '### 3. 全导出名的端到端通路（乙形前提"要跨包入口"是否只有一条形）'
echo '-- push 的两枚导出入口（Gate 侧）逐字'
git grep -n "^func (g \*Gate) PendingApproval\|^func (g \*Gate) Replay\|^func (g \*Gate) Queue" "$A" -- internal/agent/approval/gate.go
echo '-- 改写输出的导出读出口（测试拿叠改后的 corr 走这扇）'
git grep -n "^func (q \*Queue) LiveApprovals" "$A" -- internal/agent/approval/pending_read.go
git grep -n "CorrelationID: it\.Corr" "$A" -- internal/agent/approval/pending_read.go
echo '-- Decision.corr 的生产搬运工（req→dec→push 全导出字段）'
git grep -n "CorrelationID: req.CorrelationID\|CorrelationID: orDefault" "$A" -- internal/tools/bridge.go
echo '-- D47 前置（PendingApproval 第一跳的 admitCheck ⇄ 导出登记口）'
git grep -n "^func (g \*Gate) AdmitTextTask" "$A" -- internal/agent/approval/gate.go
echo '-- 名册侧导出入口'
git grep -n "^func TaskRosterSectionFrom" "$A" -- internal/panel/subagent_roster_197.go

echo
echo '### 4. push 撞名支的完整字节（含两枚操作数与空 corr 回退，丙形看不见的语义在这一段）'
git show "$A:internal/agent/approval/queue.go" | sed -n '154,184p' | nl -ba -v154
echo '-- callCorr 两臂格式串逐字（铸造侧的另一枚 #）'
git show "$A:internal/agent/loop.go" | sed -n '603,608p' | nl -ba -v603

echo
echo '### 5. 臂承重名册的原料（⛔ 本腿跑，只逐枚点名表项与两枚臂的字节判据）'
echo '-- 数字臂的判据逐字（LastIndexByte + 尾段全数字）'
git show "$A:internal/panel/subagent_roster_197.go" | sed -n '333,344p' | nl -ba -v333
echo '-- 两枚名册里全部带 # 的表项（逐枚，file:line＋字面串＋want）'
git grep -nE '\{"[^"]*", "[A-Za-z0-9_.:-]*#?[^"]*", (true|false)\}' "$A" -- internal/panel/subagent_corr_task_prefix_328_test.go internal/panel/subagent_blocked_220_test.go
echo '-- 表项之外的 card197/L1WindowWait 直投形（逐枚）'
git grep -nE 'card197\("[^"]*"|CorrelationID: "[^"]*"' "$A" -- internal/panel/subagent_corr_task_prefix_328_test.go internal/panel/subagent_blocked_220_test.go

echo
echo '### 6. 丙形的"有没有牙"原料（谁读改写输出）'
echo '-- approval 测试里带 # 的 corr 字面串断言（预期 0＝格式变了没有既存用例红）'
git grep -nE '"[A-Za-z0-9_.-]+#[A-Za-z0-9_.-]*"' "$A" -- internal/agent/approval | grep '_test\.go' || echo '  => 0 命中 rc=1'
echo '-- 正控：同尺族打到 internal/panel（该响）'
git grep -nE '"[A-Za-z0-9_.-]+#[A-Za-z0-9_.-]*"' "$A" -- internal/panel | grep '_test\.go' | wc -l
echo '-- approval 测试读 it.Corr 的符号化断言枚数（读变量、不比字面串）'
git grep -n "\.Corr\b" "$A" -- internal/agent/approval | grep -c '_test\.go'
echo '-- 既存 replay 那一发怎么断言新 corr（符号化 ⇄ 字面串）'
git grep -n "corr == p.CorrelationID\|fresh, corr, err" "$A" -- internal/agent/approval/queue_test.go

echo
echo '### 7. AC#4 去牙落点候选（具名，⛔ 同一枚断言行两处落点）'
git grep -n "%s#%d\|#replay" "$A" -- internal/agent/approval/queue.go
git grep -n 'HasPrefix(key, task+"#")' "$A" -- internal/panel/subagent_roster_197.go
git grep -n "if _, clash := q.byID\[corr\]; clash" "$A" -- internal/agent/approval/queue.go
