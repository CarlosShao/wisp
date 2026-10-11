# 309-v1 · `AC#5` 四本账对账（⛔ 同一件事记两次；按四张票**自己的正文**判，⛔ 按实现腿的摘要判）

尺一律具名；原始＝`logs/ac5-foursome.log`／`ac5-texts.log`／`ac5-bookings.log`／`ac5-gate154.log`。

## 1. 票 307（`307-…-done.md`，Status done，六格全勾）
- 它自己第一行逐字＝"`syncdirs` 那道…兜底**只折一边**：候选走 C26 解析器折成长形，`home` 只小写留短形 ⇒ 资料目录下的写读成"不落在任何同步根"（**fail-open**，方向与票 252 相反）"；性质行逐字＝"**缺陷票**，方向＝少报"。
- 与本票关系＝**同向不同因**：307 少报在**路径折叠**那一层，309 少报在**探测器挂没挂**那一层；309⛔ 动 `internal/risk/**`（票面边界），307 已 `-done` 且残余逐枚落到票 72／票 40／编排者车道（它自己 Status 行自陈）。
- ⛔ 同一件事记两次：309 正文⛔ 重述 307 的读数，只在 `AC#5` 里点名字。✅ 未重复记账。

## 2. 票 40（`40-security-privacy-cost-pages.md`，Status ready-for-agent，0 勾／8 未勾）
- 它自己逐字（`:13`／`:19`）＝"Security page (live session grants with revoke, blacklist status, taint-hit …)"／"taint-hit (R4) history from"⇒ **票 40 是 R4 的显示面**，⛔ 是 R4 的判定面。
- "零消费者"那枚归口⛔ 住在票 40 正文里（尺＝`grep -n "零消费者\|消费者" 40-*.md` ＝ 0 命中），它住在**台账 `A838` §3**（`A838` 标题逐字："`SyncRoots()`＋`SyncDetectionComplete()` 生产调用者＝**0** ⇒ 那两处注释把一枚还没建的页面写成现在的消费者；归口＝票 40"）。
- ⇒ 309 与票 40 **⛔ 同一条路径**（一个改判定输入的可核性、一个建读面）；309⛔ 把票 40 的活揽进来（`git show --name-only` 逐笔里⛔ 一枚面板文件）。✅
- 关系一句话（写进裁决表）：309 的 Ⓑ 让"空 id 那次调用⛔ 报 R4"这件事**可核**；票 40 决定"报了 R4 的历史用户看不看得见"。两格互⛔ 顶。

## 3. 票 197（`197-…-done.md`？——实际⛔ 带 `-done`；2 勾／5 未勾）
- ★**腿的摘要有一处⛔ 对不上票 197 的正文**：`AC#5` 说"票 197＝`subagent_197.go` 那支拒绝形的出处"。尺＝`grep -n "没有宿主给的任务 id|parentID == ""|TaskID(ctx)" 197-*.md` ＝ **0 命中** ⇒ 那句拒绝文案**只在产码注释里**（`internal/tools/subagent_197.go:263-266`，本腿逐字读过），⛔ 在票 197 正文里。
- ⇒ 出处其实是**文件名对票号的命名**（`subagent_197.go`）＋产码文本。这⛔ 是"记两次"，倒是**一枚口径缺陷**：裁决表按"产码文本＝③ 的出处，票 197 正文⛔ 谈这枚拒绝"落一行。
- 顺带（⛔ 本票账）：今天 `internal/agent` 之外唯一**转发** caller 的形状＝`internal/tools/subagent_197.go:581 return p.inner.Execute(ctx, req)`（尺＝`grep -rn "\.Execute(" internal cmd --include=*.go` 滤测试＝4 命中：`loop.go:760`/`:764`/`bridge.go:545`/`subagent_197.go:581`）。它是**转发**⛔ 构造 ⇒ ⛔ 撞枚数钉；它把 `gate-clauses.sh` 的 **G1b** 顶到会响那枚读数，`A844` §3 已具名归口＝**票 154 那一族**、逐字"⛔ 在 309 里重登记 `want_n`"。✅ 309 遵住了。

## 4. 票 242（`242-…-done.md`，Status 段 3 勾／0 未勾）
- 它自己正文逐字（`:14`／`:19`／`:74`）讲的是 correlation id 从 `taskID` 改成每调用一枚，且 `:74` 追加更正具名"拆它的是**票 259 那笔 `77150dcc`**"。
- 与 309 的关系＝**前因**：`bridge.go:926-929` 那句"the binding is per scope, hence per task"能成立，正因为 id 已经⛔ 再是 per-call 的 correlation（`task.go:704-707` 注释逐字"id went per-call in ticket 242 and no longer names a task"）。
- `AC#3` 要求⛔ 改契约声称 ⇒ 那句**一字未动**（尺＝diff 的 ± 行清单里⛔ 它，见 `logs/ac3-changed-lines.txt`）。✅ 309⛔ 抢 242 的账、也⛔ 把 242 的活搬进本票。

## 5. 台账那一枚欠读数⛔ 在本票重记
`A842` 标题逐字＝"★结掉 `A838` §3 我自己登记那枚欠读数：`bridge.go:925` 的 `taskID == ""` 那一肢今天不可达…"；票 309 `AC#5` 逐字"由 `A842` 结掉，⛔ 在本票再记一遍"，正文也确实只**指**⛔ **复述**。✅

## 6. ★本腿新查到一枚"⛔ 该由 309 记第二次"的先例（重要，直接影响 `AC#2` 的判语）
`.scratch/wisp/probes/154-accept/ac1-g1-blindspot-synthetic.txt`＝**票 154 的非实现者验收程**（票 154 已 `-done`）已经在**一棵合成树上量过**"宿主自造 id 走到 `Execute`"的**五形**：
```
A  var req agent.ToolRequest; req.TaskID = "…"; b.Execute(ctx, req)
B  type dispatch = agent.ToolRequest  → 字面量文本是 dispatch{，⛔ 是 ToolRequest{
C  req := agent.NewRequest()（字面量在别的包里）
D  req2 := seed; req2.TaskID = "…"（拷贝既有值再改 id）
E  fn := (bridgeImpl{}).Execute; fn(req)（方法值/函数指针间接调）
```
它那把尺（`grep -rnE 'ToolRequest\{' --include='*.go'`）打上去的结论逐字＝"A／B／C／D／E 五枚**宿主函数本身**没有一枚被 G1 吃到"。
⇒ 本腿在 `AC#2` 用真仓**重打**了其中两形（形 2 声明＋赋值、形 3 别名）并同样拿到"钉⛔ 响"（`logs/nail-2-decl.log`／`nail-3-alias.log` 各 rc=0）。
⇒ **结论两重**：① 枚数钉的形状洞⛔ 是 309 新发现的债，它是 `A`/`B`/`C`/`D` 那一族**已登记先例**的一次复现 ⇒ 归口＝编排者（票 154 族／Q-56／新票），309 侧只准**引用**⛔ 准另立；
② `AC#2` 那格正文写的是"**任一** caller 手搭 `ToolRequest` 递进空 `TaskID` ⇒ 在盘上可核地响"，而先例＋本腿复现一致说"任一"⛔ 成立 ⇒ 本格⛔ 判"成立"（见 `90-final.md`）。

## 7. 本票⛔ 揽进来的东西（自查，逐名）
`git show --name-only` 五笔并排＝`internal/tools/` 三枚（1 产码＋2 测试）／`probes/309/r1/**`／票 309 末尾 1 行 ⇒ ⛔ `cmd/wisp`、⛔ 面板件、⛔ `internal/risk/**`、⛔ 票 154 的门脚本、⛔ 票 298 的 `gofumpt` 债（`A846` 已承接）。
