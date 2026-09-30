# 246 — 反向缺口：**会跑任务的那条腿不在带着球的常驻进程里** ⇒ 审批卡片进不来，D43 承诺的四条否决通道（单击球／Esc／语音否决词／面板拒绝）**在生产进程里一枚都没有执行者**

- Status: **已立，未派**（09-30 18:2x，编排者立；账 `A480`）。触发＝验收腿 `245-v1` 交回的一条事实，我自己现跑复认过。
- 与票 228 的关系：票 228 修的是**半边**（"球与托盘不在跑任务的那个进程里"⇒ 把球搬进常驻进程，已落）。**这一枚是另半边**：常驻进程今天仍然**不跑任务、没有审批门、没有面板宿主** ⇒ 球进来了但**没有东西可否决**。两枚合起来才是那句产品事实："**双击图标起来的那个进程，能干活、能被打断**"。

## 现量（09-30 18:23 编排者自己跑，别信行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 常驻那条腿自己逐字承认没有管线 | `cmd/wisp/resident_ball_windows.go:22` 逐字「pipeline, no microphone, no approval gate and no panel host today, so this…」；`:126` 逐字「gestures", "recorded only: this leg has no task pipeline, no microphone and no approval gate」；`:163` 逐字「…so the gesture has no executor here」 | `grep -n "Confirming\|approval" cmd/wisp/resident_ball_windows.go` 我现跑 |
| "借 Esc"那一支在生产里零调用者 | `internal/ball/ball_windows.go:873 TakeEscForCancel`／`:895 ReleaseEscAfterSession` 的非测试消费者**只有旁支程序** `cmd/balldebug/main.go:635`／`:638` | `grep -rn "TakeEscForCancel\|ReleaseEscAfterSession" --include=*.go internal cmd \| grep -v _test` 我现跑 |
| 契约要求的是**四条**否决通道 | `docs/PLAN.md:3082`（**D43 转移表＝C12 冻结**）第 22 行逐字「否决（单击球 / `Esc` / KWS 否决词 / 面板拒绝，B1）」；`internal/agent/approval/approval.go:90` 逐字 `ChannelEsc: "按 Esc 键"` | 票 245 已逐字读过，本票**不改那两枚文件一字** |

## 前段（派单第一枚只能是普查，不许直接开写）

- [ ] **AC#0 先把"卡片怎么进常驻进程"的两形代价摆开（不许直接开写）**：甲＝常驻腿**自己起一条 loop／审批门**（新依赖边？谁装配？）／乙＝**装配根 `cmd/wisp` 注入**（与台账「票 197 段：装配根是唯一的接缝」一致）。完成判据＝两形各带"要动哪几枚文件＋新增哪几条依赖边＋会不会破 `internal/proc` 的 D38(e) 十步顺序＋与票 228 后续片（`config.toml` 未接、托盘「退出」无执行者）谁先谁后"，由编排者裁后再派落地腿。⛔ **普查腿不许改任何产码**；写点只准落在 `.scratch/wisp/probes/246/a1/` 与本票面。

## 判据（前段裁完之后才许补，⛔ 现在不填）

- [ ] **AC#1.. 待 `246-a1` 交回后由编排者落**（现在填＝凭想象写判据，本仓已因此白跑过两次）。

## 禁区

- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；D43 转移表那句"四条否决通道"**本票只能去实现它，不能去改它**。
- ⛔ 不新造第五枚否决通道、不改冻结契约 C12／C18 的形状；`frontend/**`／`design/**` 零读零写零转述。
- ⛔ 面板侧要的东西**只写进本票面与台账**，由 owner 自己带给他用的那枚前端 agent；本编队永不调用跨会话工具去联系任何会话。
- ⚠ **下一枚动 `cmd/wisp` 的腿必须先归因一枚已知的红**：`245-v1` 自报该包默认层「1 红 2 绿」但用例名丢了（`A480` ④）⇒ 起手先跑默认层把红名取出来、判断是不是本票造成的，⛔ 不许当已知常红略过。
- ⛔ **不新增依赖边要先经编排者裁**（本仓有"两枚正向边一律不开、改注入"的先例：票 238 第 1 刀）。
- git：只 commit 不 push；显式 pathspec；禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 排程与串行

- 排在**队列最前**（先于 `197-r3`）：owner 09-30 原话「你可千万别忘了要紧事儿，**你得赶紧干后端的活儿哈**」——本票是后端主链（"能干活＋能被打断"），票 197／224 是仪器精度。
- ⛔ 与票 228 后续片、票 245 的 AC#6..AC#9 同撞 `cmd/wisp`／`internal/ball` ⇒ **一律串行**，跑突变的腿绝不并发；派单里写死"此刻哪几枚包有别的写腿"。
- 测量坑：跑真机热键／桌面用例前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 计数必须为 0；`cmd/wisp` 缺 sherpa PATH 会 `0xc0000135` 且**无 `--- FAIL`**＝根本没跑。

## Progress log (append-only, newest last)
- [2026-09-30 18:2x +08] agent=246-a1 did=前段只读普查交件 `.scratch/wisp/probes/246/a1/census.md`(206 行/31,880 字节, 占位符 0, 骨架先落 f9cbe725→填满 63508c50)。三条现读凭据独立复认成立; 两形代价表各 ①-⑥; 依赖边差集为空(甲乙合规做都 0 条新增包级边, cmd/wisp 已 import agent/approval/ball/proc); D38(e) step3 cancel-task-roots 是 proc 侧两形共用缺口(boot_windows.go:151-161 只填 CloseJob, step1-7 零生产者)。建议乙形(装配根注入, 守"唯一接缝"、不造第二真相源), 最短链『卡片进得来+Esc 真能否决一次』可只注入裸 gate(不需 loop/provider/config/麦)。硬契约改动≈0(D43/C12 一字未动), 需告知功能 3-4 枚, proc hook 注册入口是否算契约面待 owner 裁。go build/vet rc=0; 未跑 go test; A480④ 未归因红留落地腿。零产码、未碰任何 AC 框。next=编排者裁甲/乙后派落地腿
