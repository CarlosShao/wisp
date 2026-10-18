#!/bin/sh
# 票 340 普查腿 340-a1 · AC#0 起手现量（只读；一律取具名提交树，⛔ 工作树字节）
# 逐字重跑：sh .scratch/wisp/probes/340/a1/scripts/ac0.sh
# 入场锚＝本腿 2026-10-18 11:3x `git log -1` 现认；票面记锚与编排者裁决锚一并撞。
ANCHOR=feac8f77c2fbb15addbc36fb1bbaf13ed56588be
TICKET=f705b082f452754af477dff096e7000e590c5d82
ORCH=5db0b538dad3bf2841adec834a130e4a6b28d8df

echo "### 0. 锚自证（⛔ 打印 HEAD 标题＝那字节属在飞别腿，不进本腿件）"
echo "本腿入场锚（11:3x 现认，全尺钉在这一枚提交树上）=$ANCHOR"
echo "跑动时刻的 HEAD（仅供对照，⛔ 本腿读数来源）:"
git rev-parse HEAD

echo
echo '### AC#0 a-1 票面窄尺：真产叠改形枚数，射程＝internal/agent/approval'
for R in "$ANCHOR" "$TICKET" "$ORCH"; do
  printf 'rev=%s  命中行数=' "$R"
  git grep -nE '"[^"]*#[0-9]+"' "$R" -- internal/agent/approval | wc -l
done
echo '-- 逐字（入场锚，预期无输出）'
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal/agent/approval || echo '  => 0 命中 rc=1'

echo
echo '### AC#0 a-1b 这把窄尺的响应性（正控＝同一枚字节序列、只换射程）'
echo '-- 正控 1：尺逐字节不变，只换射程＝internal/panel ⇒ 该响（入场锚行数）'
printf 'panel@ticket rows=' ; git grep -nE '"[^"]*#[0-9]+"' "$TICKET" -- internal/panel | wc -l
printf 'panel@anchor rows=' ; git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal/panel | wc -l
echo '-- 正控 1 逐字（入场锚，internal/panel）'
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal/panel
echo '-- 正控 2：尺逐字不变，射程＝全仓枚数（证明这把尺在别处看得见东西）'
git grep -nE '"[^"]*#[0-9]+"' "$ANCHOR" -- internal cmd tools scripts | wc -l
echo '-- 反形（⛔ 一把两形都放行的尺）：只改字符类 [0-9] 到 [A-Z]，approval 射程'
git grep -nE '"[^"]*#[A-Z][A-Za-z0-9]*"' "$ANCHOR" -- internal/agent/approval || echo '  => 0 命中 rc=1（说明 approval 射程内没有任何 "#字母" 字面串，⛔ 拿这形当正控）'

echo
echo '### AC#0 a-3 自加的宽尺：叠改形字面串（两枚 #）产码档，⛔ 测试档／⛔ .scratch 档'
echo '-- 档 A 产码（排除 *_test.go）＝票面那一格的真产读数'
git grep -nE '"[^"]*#[^"]*#[0-9]+"' "$ANCHOR" -- internal cmd tools scripts | grep -v '_test\.go' || echo '  => 0 命中 rc=1'
echo '-- 档 B 同一枚尺不过滤（正控，必须响；逐枚点名它们全是测试/注释）'
git grep -nE '"[^"]*#[^"]*#[0-9]+"' "$ANCHOR" -- internal cmd tools scripts

echo
echo '### AC#0 a-4 真产链的构成条件现读（两枚铸造点各自的生产调用者枚数）'
echo '-- Queue.push 调用点：产码档逐字'
git grep -n '\.push(' "$ANCHOR" -- internal/agent/approval | grep -v '_test\.go'
echo '-- Queue.push 调用点：含测试总枚数'
git grep -n '\.push(' "$ANCHOR" -- internal/agent/approval | wc -l
echo '-- callCorr：产码档逐字（def ⇄ 唯一调用点）'
git grep -n 'callCorr' "$ANCHOR" -- internal cmd tools | grep -v '_test\.go'
echo '-- panel／cmd 的测试有没有 import 过 approval（跨包真队列的现量门槛）'
git grep -n 'wisp/internal/agent/approval' "$ANCHOR" -- internal/panel | wc -l
git grep -n 'wisp/internal/agent/approval' "$ANCHOR" -- internal/agent ':!internal/agent/approval' || echo '  => 0 命中 rc=1（internal/agent 侧零测试 import approval）'
echo '-- 正控（同尺族证明这把 import 尺会响）：cmd/wisp import approval 枚数'
git grep -l 'wisp/internal/agent/approval' "$ANCHOR" -- cmd/wisp | wc -l

echo
echo '### AC#0 b 第三臂侧叠改形那枚子用例的具名落点（⛔ 本腿跑它，只点名）'
git grep -n 'arm B with the queue clash rewrite stacked' "$ANCHOR" -- internal/panel

echo
echo '### AC#0 c 数字臂 ⇄ 第三臂：两枚臂的本体与全部消费点'
echo '-- 第三臂本体（rostered-task arm）逐字'
git grep -n 'HasPrefix(key, task' "$ANCHOR" -- internal/panel
echo '-- 数字臂本体 splitClashSuffix：def ⇄ doc ⇄ 唯一读者'
git grep -n 'splitClashSuffix' "$ANCHOR" -- internal/panel
echo '-- 数字臂的字符判据逐字（LastIndexByte 那枚尺）'
git grep -n 'LastIndexByte' "$ANCHOR" -- internal/panel
echo '-- 全仓测试名册里带两枚 # 的字面串（叠改形夹具，逐枚点名）'
git grep -nE '"[A-Za-z0-9_-]+#[A-Za-z0-9_-]+#[0-9]+"' "$ANCHOR" -- internal cmd
echo '-- cmd/wisp 卡片 corr 带 # 的字面串枚数（裁决件那把界，本腿复跑）'
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' "$ANCHOR" -- cmd/wisp || echo '  => 0 命中 rc=1'
echo '-- 正控：同尺族打到 internal/panel 的行数'
git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' "$ANCHOR" -- internal/panel | wc -l

echo
echo '### AC#0 d 两枚铸造点的非导出性 + 那枚改写格式串逐字（丙形承重）'
echo '-- callCorr def（非导出＝首字母小写，逐字）'
git grep -n '^func callCorr' "$ANCHOR" -- internal/agent
echo '-- Queue.push def（非导出，逐字）'
git grep -n '^func (q \*Queue) push' "$ANCHOR" -- internal/agent/approval
echo '-- 两枚格式串逐字（带行号）'
git grep -n '%s#%d' "$ANCHOR" -- internal/agent/approval/queue.go
echo '-- queue.go 整枚 blob 哈希跨锚对撞（同串＝三锚的这枚字节未变）'
for R in "$ANCHOR" "$TICKET" "$ORCH"; do
  printf '%s  ' "$R"; git show "$R:internal/agent/approval/queue.go" | git hash-object --stdin
done
echo '-- approval 包导出 func 名册里有 push 吗（乙/甲越界的现量凭据）'
git grep -nE '^func \(q \*Queue\) [A-Z]' "$ANCHOR" -- internal/agent/approval

echo
echo '### AC#0 e 票面那行第三臂 PASS 凭据：⛔ 本腿零 go（只读普查）'
echo '本腿⛔ 跑 go test／go build ⇒ 那一行的颜色＝读码推的，凭据＝裁腿判语件第一节 1 条逐字（docs/evidence/s1/328-v1-verdict.md:11）'
