# 票 237 验收腿 237-v1（非实现者判语）——三格＋归因链＋生产面读数

- 被裁：`d8a4d913`（产码/测试面：`subagent_197_test.go` +1、`subagent_222_test.go` +136/-47；零产码）／
  `f3dc1389`（交付件：`probes/237/r1/evidence.md` 84 行＋锚/现量/变异/门禁读数＋票尾 bullet）。
- 锚：HEAD `bfb62488`（送发后一枚编排者文档提交，只动 `.scratch/orch/**`、`HANDOVER.md`、台账；代码面＝`d8a4d913` 态）。
  两枚 `_test.go` 现量 `d2d52bff…`／`ba852ebf…`，`subagent_197.go` `b72f6740…`（零改），`bridge.go` `4f8f6e90…`。
  GOTMPDIR `C:\tmp\wisp237v1gotmp`；逐发窗口与 tasklist 见 `probes/237/v1/logs/00-anchor.txt`。
- **结论先行**：

| 格 | 判语 | 一句话凭据 |
|---|---|---|
| AC#1（选甲＝等待条件改成真正依赖的事） | **成立** | 自造回退形 `-race -count=60` 一发红且红句＝"父任务没派生成功"族（`logs/10-ac1.txt` [1]）；交付态同命令 ×2 零红；diff 读＝逐枚 parked 因果等待，非"多等一会儿" |
| AC#3（选甲＝收紧成逐枚） | **成立** | 四形逐枚复跑与 r1 表同形（§2）；C25 scope 补丁双向坐实不是放宽（一枚绕过仍红＋摘回 no-op 三红） |
| AC#2（注释＋1 行、禁断言注释文字测试） | **成立** | `git show d8a4d913 -- …197_test.go` 现读＝纯 `//` 注释 +1/-0，无断言变动 |
| 零产码 | **成立** | 两笔名册无第三枚代码件；`git diff d8a4d913~1..f3dc1389 -- …197.go …bridge.go` 空 |
| unbound-scope L2 归因链 | **成立** | `dump5-leg2.txt` 现读＋源码逐行（§3） |
| ★生产面（r1 留给编排者的具名项） | **不是"只在 harness 缺、生产无恙"——生产同一机制可达** | §4；`OpenTask` 非测试调用点全树仅 `bridge.go:649`（mark 懒开），生产装配根不开任务 scope |

## 1. AC#1（成立）

- 自造回退形＝**交付态现件、只把 `parkParents` 函数体退回** `d8a4d913~1` 的快照形（diff 即那 12 行＋4 行头注）：
  `probes/237/v1/mut/x-parkparents-snapshot_222_test.go` ＋ `overlay-parkparents-snapshot.json`（`-overlay`，工作树零突变字节）。
- 发1（`-overlay … -race -count=60 -run Test222 ./internal/tools/`）09:54:03→09:54:16 **rc=1**，1 枚 FAIL：
  `subagent_222_test.go:516: 只有 3/4 枚父任务到达等待点（父任务没派生成功，等待读数就无从谈起）`，DATA RACE **0**
  （:516 属变异件＝交付件对应行 +4 行头注；r1 改前读数形状 3 枚/1 枚红，我的 1 枚/60 落在同族）。
- 交付态同命令 ×2：09:54:35→09:54:45 rc=0 `ok 3.900s`；09:55:05→09:55:16 rc=0 `ok 4.285s`。两条合计 0 FAIL／0 DATA RACE。
- 读 diff（`d8a4d913` 的 `parkParents`）判"等的确实是它依赖的那件事"：旧形在收满 n 枚 `started` 后**立刻拿
  `len(h.parked)` 当快照**；产码里 `RunAsync`（`subagent_197.go:347`）**早于** `giveBackWhileWaiting`（:368）与
  "已派生" delta（:369-371）⇒ 无蕴含。新形＝**逐枚收 `parked` 通道**（每枚父任务一条、通道收信即因果），
  并以"父任务带着结果提前现身"为第二条证伪支；**无降 count、无加长等待**（30s 只在真挂住时开火）。回退形红 ⇒ 甲成立。

## 2. AC#3（成立；四形逐枚复跑）

| 形状（r1 的 `mut/`＋`overlay-*.json`，我现跑 `-race -count=1 -run Test222 -v`） | 我的读数（rc；逐腿） | 与 r1 表比对 |
|---|---|---|
| `x-partial-after`（当前件＋一枚孩子绕开真桥）09:55:17→09:55:29 | rc=1；leg1 **3.00s** 红／leg2 **0.00s** 红／leg3 **0.00s** 红（"到达 3 枚、已 4 枚返回"） | 同形 |
| `x-partial-before`（HEAD~ 件＋同一枚绕过）09:55:48→09:56:27 | rc=1；leg1 **30.00s** 红（旧护栏"只等到 2/4"）／leg2 **PASS**／leg3 **PASS** | 同形（235-v1 点名失明形） |
| `x-full-after`（全部孩子绕开＝M3 回归）09:56:53→09:57:04 | rc=1；三腿全红 **0.00s**（"到达 0 枚、已 1 枚返回"） | 同形 |
| `x-noopadmit-after`（C25 scope 布线摘回 no-op）09:55:39→09:55:48 | rc=1；三腿全红（leg1 1<2／leg2 3<4／leg3 2<3，同一族红句）；leg1 时长 **0.01s**（r1 表写 3.00s） | 判据落；**时长有差**见下 |

- ★关键判语＝**harness 的 C25 scope 补丁不是放宽**：① 补丁在位时"一枚孩子绕开真桥"仍必红（partial-after）；
  ② 把补丁摘回 no-op 三腿必红（noopadmit）。双向坐实，红来自到达/返回覆盖计数本身。
- 具名差异：noopadmit 的 leg1 时长我读 0.01s、r1 表写 3.00s——红句同一族（到达数少于返回数），只是首个父返回出现的时点
  不同；判据"三枚仍全红"成立，不追平也不视作不一致。

## 3. unbound-scope L2 归因链（成立）

- `probes/237/r1/logs/dump5-leg2.txt`（68,968 B，留盘）现读：`VERDICT237 <孩子 task uuid> level L2 rules [R4]
  reason R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容`
  （:54/:56/:58/:94）＋`BEXIT237 probe.work222 <同一 uuid> … isErr true text L2 审批通道尚未接入（票 21），已拒绝执行`——
  后一句**逐字节＝`gate.go:144`**（NoGate 的 L2 拒词）。同 dump 里先行的 `VERDICT237 … level L0` 是"mark 尚未落地"的
  第一枚孩子；`PROV237 parent-task-222 …` 是结论 stamp 落父名。
- 源码链逐枚核过（行号自现读）：`subagent_197.go:487` `Provenance.Mark(parentID, SrcTaskOutput, …)`（stampConclusion
  :479-490，首枚收口即产生 mark；`Mark` 对"closed/never-opened scope"照记，`provenance.go:558-575` 文＋`:625-629` 建 reg）
  → 孩子未注册 task id 走桥，`assessorFor(taskID)` 挂 `prov.Detector(taskID)`（`bridge.go:996-1005`）→
  `scopeMarks` 见"本 scope 未注册、别处有 mark"即 fail-closed（`provenance.go:694-711`，log :706；`Inspect` 建 Hit :726-731）
  → R4（`assessor.go:18/:81`：C25 taint hit → L2，不可被会话授权覆盖）→ `bridge.go:458` `gate.PendingApproval`
  → NoGate → `gate.go:144` 拒。**判语：归因成立。**

## 4. ★生产面读数：`OpenTask` 谁开、何时开（判"harness-only"还是"生产真形"）

- **谁开：生产没有"任务开始开 scope"这个调用点。** 全树非测试 `OpenTask(` 调用点**只有 1 个**：
  `internal/tools/bridge.go:649`（`b.mark` 内、实参 `dec.TaskID`；且仅在 `risk.IsSensitiveSource(dec.Tool)` 的**结果收口**时
  走到，`bridge.go:642`）。另有定义 `:847`。测试面 `.OpenTask(` 仅 3 处（含一条注释）。
- **任务边界钩子不开**：`cmd/wisp/run.go:953-962` `admitTask`＝noteTask＋`gate.AdmitTextTask(:955)`＋revoke{revoke,
  `bridge.CloseTask(:960)`}——**只关不开**；`internal/agent/loop.go:373-379` 任务开始调它、defer revoke；子环路
  `subagent_197.go:311-320` 包一层 `baseAdmit`——也不开。
- **触发条件**（`provenance.go:694-711`＋`Detector` :858-869 每次调用都 Inspect）：调用方 scope 未注册＋**任一别的 scope
  有 mark** ⇒ 该任务**每一次**工具调用都被判 unbound-scope R4/L2。而"别的 scope 有 mark"在生产现成：R 派孩子 c1..c4
  并发，c1 先收口 ⇒ `stampConclusion` 把 task.output mark 落 **R 名下**（`:487`）⇒ 此后任何自己的 C25 scope 尚未注册的任务
  （每个孩子在被自己第一枚 mark 打开前都是）下一次走桥即 L2。孩子的工具面＝父桥（`subagent_197.go:296`），dispatch 带
  自己的 task id（`loop.go:675-676`）⇒ 形状与 h222 完全同构。另有常驻一枚：`run.go:1052` 为项目说明开**非任务** scope
  `run:project-instructions`（defer :1053 关），装载成功即有 mark（`projctx.go:349`），整跑期间构成"别的 scope 有 mark"。
- **症状落到真门**：生产 L2 → 真审批门（`run.go:612` `approval.New`，console UI＋超时可拒）⇒ 假 L2 审批卡/超时拒绝；
  h222 里是 NoGate 直接打回（`gate.go:144`）。机制同、终端不同。
- 旁证：151 裁决件已登记"`OpenTask` 只在 `mark` 里懒开"且**刻意不接**"任务一开始就 OpenScope"
  （`docs/evidence/s1/151-task-scope-never-closed-r1.md:85/:197`）；测试先例 `pointer_185_cli_seam_test.go:346` 正是拿
  `b.OpenTask` 补这一格（注释自认 composition root 会在任务开始开 scope——**生产里没有这个调用点**）。
- **判语：不是"只在 harness 缺、生产无恙"；同一机制在生产链路静态可达（含逐行凭据）。**
  本腿未跑 `cmd/wisp` 活体复量（写面/纪律：313-v1 在飞；本腿面＝`internal/tools`）——建议编排者落 `A##`/`Q##` 时先发一发
  cmd/wisp 复量（注意 `ProjectInstructionsEnabled` 默认 `true`，`internal/config/schema.go:444`）。

## 5. AC#2 与零产码复核（成立）

- `git show d8a4d913 -- internal/tools/subagent_197_test.go` 现读＝**+1/-0**：唯一新增行是 `//` 注释（守池注释段尾，
  交付件 :411），删除列为零，**无断言变动**；未新写"断言注释文字"的测试（`git diff --stat` 197：+1/-0）。
- 零产码：`git diff d8a4d913~1..f3dc1389 -- internal/tools/subagent_197.go internal/tools/bridge.go`＝**空**（复跑）；
  两笔完整名册＝d8a4d913 只两枚 `_test.go`；f3dc1389＝票＋evidence.md＋4 枚小 log。无第三枚代码件。

## 6. 门禁抽核与全包（rc；窗口 10:01:42→10:02:17／10:03:12→10:05:25）

| 项 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | 0 | 零输出 |
| vet | `GOFLAGS= go vet ./internal/tools/` | 0 | 零输出 |
| d22scan | `sh scripts/d22scan.sh` | 0 | `d22scan: clean`（23,749 B 原文留盘＋入库） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN；over-budget=58；单位互核同 58 |
| 全包 race 发1 | `go test -count=1 -race ./internal/tools/`（GOTMPDIR 正斜杠拼法） | 1 | 7 FAIL **全为环境前置**（junction 造不出/分隔符 want）；根因＝GOTMPDIR 递进测试进程临时根后 `mklink` 参数带 `/`；**DATA RACE 0**、与 237 断言无关 |
| 全包 race 发2 | 同上，`GOTMPDIR='C:\tmp\wisp237v1gotmp'`（派单例形） | 0 | `ok 21.798s`；0 FAIL／0 DATA RACE |

## 7. 具名没做与边界

- 未跑 `cmd/wisp` 生产活体（§4 是"静态链逐行＋h222 现量"的合成判语；活体复量建议编排者先跑）。
- 未复跑 r1 的改前"红两发"（其 `pre-fix-race60-shot1/2.txt` 现量 770 B/383 B 与原 anchor hash 同态，且本腿自造的
  回退形已独立复现同族红——不改写实现者件）。
- 未跑 r1 的 instrumentation overlay（`dump*`）；归因链用其留盘 dump5＋现读源码核。
- 临时件只建不删：raw 逐发日志（`fallback-shot1` 400 B、`delivery-*` 54/54 B、`mut-*` 2,184/2,257/720/2,257 B、
  `fullpkg-race1` 24,617 B、`fullpkg-race2` 55 B）留盘；入库件＝本文件＋`probes/237/v1/` 6 枚（锚、3 枚合订 log、
  回退形变异件、overlay）。
