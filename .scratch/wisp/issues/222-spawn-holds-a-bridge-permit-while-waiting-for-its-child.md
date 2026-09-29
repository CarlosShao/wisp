# 222 — **`task.spawn` 一边占着桥的执行许可、一边等孩子跑完** ⇒ 父任务永远在"不等了"那条错误文本上收场；而**测试看不见这一形，因为测试给孩子的工具面绕过了桥**

- Status: **产码已交并静态核过（09-29 12:46，见文末"收 `222-r1`"那一节），⛔ 六格今天一枚没勾**——两枚前置未到：① **我自己的整包复跑**（被在飞写腿 `226-r1` 挡住：它正在改 `internal/config/allowdirs.go`／`loader.go`／`permmode.go`，而 `cmd/wisp` 与 `internal/tools` 都吃 `internal/config`⇒ 现在跑编译读到的是它半截的码）；② **非实现者终裁 `222-v1`**（`SPEC-12 §4.3`／`AGENTS.md` §0.3：验收必须由另一枚程做，实现者的自述不是翻勾依据）。⛔ **不翻 `-done`**，且**这枚票今天只做成半格的那一半要说给 owner：「派了子任务、父任务报错」这件事现在只修掉一半**。
- 来源＝只读腿 `211-c1` 的 ④ 节（**它推翻了票 211 的前提**）＋编排者自己逐环现读（下面每一行都我自己跑过尺，不是转述）。
- ⚠ 与票 211 的关系：**票 211 那枚"要真 8 枚得先改契约 D38d"的问题，被本票证明是问错了**——见 §"为什么 211 的前提塌了"。

## 现量（起手逐条复算，别信这里的行号）

**链条一环一环，全部现读：**

| # | 事实 | 出处 |
|---|---|---|
| 1 | 桥在**调用工具之前**取得一枚许可，并且**整个调用期间一直握着**（`defer` 释放） | `internal/tools/bridge.go:439` `case b.sem <- struct{}{}:` ＋ `:440` `defer func() { <-b.sem }()`；池＝`make(chan struct{}, ceiling)`（`:171`），`ceiling` 来自 `MaxToolConcurrency = 4`（`:24`） |
| 2 | `task.spawn` 是**一枚桥内工具**，它的 `Execute` 里**阻塞等孩子的结论** | 注册＝`subagent_197.go:218`（`BuiltinSubagentEntries`）；阻塞＝`:356-369` 那个 `select { case <-ctx.Done(): … case res := <-done: … }` |
| 3 | 等的时候用的 ctx 就是桥给的**带 tool 超时的那个 context** | `bridge.go:449` `ectx, cancel := b.execContext(ctx, timeout)`（**在取得许可之后**才建），`:470` 一带 `entry.Tool.Execute(ectx, …)`；超时取值＝`timeoutFor`（`:550-555`：`Decl.Timeout>0` 用它，否则 `b.defTm`），而 `TaskSpawnDecl()`（`:202`）**没设 Timeout** ⇒ 生产值走 `run.go:548` 的 `DefaultTimeout = cfg.Agent.PerToolTimeoutMS`（未配置时 `DefaultToolTimeout = 30 * time.Second`，`bridge.go:28`） |
| 4 | 孩子在生产上**用同一枚桥**递它自己的工具调用 | `run.go:567` `ParentTools: rt.bridge` ＋ `subagent_197.go:277` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)`；注释逐字（`:129-131`）「the child receives it wrapped」 |
| 5 | 孩子的上下文**没有 deadline**（父取消不级联那一族的实现） | `subagent_197.go:327` `context.WithCancel(context.WithoutCancel(ctx))` |
| 6 | 池帽与桥帽**同为 4**（票 211 甲落地的形状） | `subagent_197.go:78` `MaxConcurrentSubagents = 4` ＝ `bridge.go:24` |
| 7 | ⛔ **测试看不见这一形**：台件给**父**用的是真桥，给**孩子**用的是一枚假目录 | 父＝`subagent_197_test.go:216` `h.bridge.Execute(...)`；孩子＝`:183` `ParentTools: h.dir` ＋ `:186` `opt.Tools = agent.Options{… Tools: h.dir}`，而 `h.dir`＝`:159` 的 `&fake197Dir{}`。**全文件 grep `Tools: h.bridge` ＝ 零命中** |

**推出来的真实形状（我自己按上面七环推，标口径）**：
- 孩子只要跑得比那枚 **per-tool 超时**久（真干活都是几分钟），父任务这次 `task.spawn` **必然**落到 `:358-365` 那条分支，收到的是**错误文本**（`IsError: true`）：「父任务这一侧已经不等了（…）。子代理 X **没有被级联取消**，它仍在名册里，**可以单独停它**：它的流键是 …」。⇒ **"同步派一枚子代理并拿回结论"这条路今天是断的**；而那句"可以单独停它"同时踩了 **票 221**（`task.cancel` 不存在）。
- 同时**孩子本身并不死**——它由 finish watcher（`:346-354`，`observe.Default.Spawn`，上下文 `WithoutCancel`）跑完并落到名册那一行 ⇒ **名册上是好的、父任务里是错的**，这正好是 owner 要"看得见状态"要看的那种分裂。
- **并发一上去更糟**：4 枚 `spawn` 各占一枚许可 ⇒ 桥的 4 枚全被"正在等"的动作占满，孩子递工具时在 `:439` 排队；这条**互等**只在父的 per-tool 超时陆续到点后才松开。⚠ **"永挂"这一支我不下结论**（腿报"无限挂住"，但父侧有 `ectx` 超时兜着 ⇒ 真形状是**周期性松绑＋父侧一律报错**，不是死锁）。**归因未做，留给验收腿现跑。**

## 为什么 211 的前提塌了

- **那枚 4 是字面冻结的，不是注释**：`docs/PLAN.md:2849`（在 `:2818` 的 D38、`:2845` 的「（d）背压」之下）逐字「**工具并发上限 = 4**（防 LLM 一次发 50 个 tool call 打爆机器；也与 D45 的批量聚合配合 —— 4 路并发足以让批量场景快起来）」；同值另见 `docs/specs/SPEC-01-architecture.md:134`／`:26`、`SPEC-05-agent-core.md:70`。⇒ **改 4＝契约改动**，这一条腿说得对。
- ⛔ **但"把 4 抬到 8"买不到任何东西**：池 8 时 8 枚 `spawn` 一样各占一枚许可，孩子照样排队 ⇒ **今天的"诚实数"不是 4 也不是 8，而是"父侧全部超时"**。
- ✅ **真正的一条路（丙）＝"等孩子的这段时间不占执行许可"**，而且**它让冻结文字逐字仍然为真**——那句理由的射程是**正在执行的工具调用**（防一次发 50 个打爆机器），而"占着许可干等"本来就不是执行。拆掉这层 ⇒ 桥仍是 4、`PLAN.md` 一字不动、池能不能写 8 变成**另一个独立问题**（见下）。
- ⚠ **别家分设两枚帽**（这条支持"8 不与 4 抢"的形状）：Step-Code `step-subagent.ts:72-73`＝**8 枚存在／4 枚在跑**；deepseek-harness `constants.ts:6` 工具 10 对 `workflow-ptc/types.ts:16` 孩子最多 16 ⇒ **孩子比工具帽多**；minimax 工具 5 per-server、孩子帽零命中；openchamber 两枚皆零命中。〔全部＝腿报，我未逐枚 `sed` 复现 ⇒ 只当"存在这种形状"用，不许当判据〕

## 判据（每格都要现跑读数；⚠ 本票的判据有一条是"改前必红、且不靠挂死来红"）

- [ ] **AC#1 生产形状第一次有尺**：新增一枚判据，**孩子的工具面也走真桥**（今天 `:183`／`:186` 用的是 `h.dir`）。它必须**在修之前红、修之后绿**；⚠ 不许写成"等超时才红"——**未修码那一发要能在秒级内判定**（判据形态自己选：数许可占用数／断言父侧不收到 `ctx.DeadlineExceeded` 文本／断言父侧拿到的是孩子结论而不是"不等了"）。
- [ ] **AC#2 "占着许可干等"这一形被永久钉住**：种一枚"父在等孩子、期间桥的可用许可数＝4-（父之外正在执行的真工具数）"的读数 ⇒ 修后父等待**不占**许可；并一枚反控：**同时执行的工具调用数仍 ≤ 4**（这条是把 `PLAN.md:2849` 的理由保住，不许为了并发把它说破）。
- [ ] **AC#3 父侧结论通道要么真、要么明说不做**：二选一并写进票面——(甲) `task.spawn` **立即返回**句柄（子代理 id＋流键），父任务之后用**现成的 `task.output`** 读结论（**别家形状**：DSH／Step-Code 都是"存在≠在跑"两枚帽＋后台句柄）；(乙) 维持"等完"，但**等待不受 tool 超时约束**且父侧文本不再报"不等了"。⚠ **我（编排者）倾向 (甲)**：它同时消掉 AC#4 那句"可以单独停它"的谎（票 221），且不新增任何导出名。**走 (甲) 要先落一枚 `A##`**（`task.spawn` 的语义属于 D34 那行的批准面，改行为不改冻结文字）。
- [ ] **AC#4 不许留任何"对模型许诺了但没接"**：`:362-363` 与 `Description()`（`:189`）两处都写着"可以单独停它"⇒ **随票 221 一起收**（两枚文件同撞 `subagent_197.go` ⇒ **必须串行、同批只一枚碰它**）。
- [ ] **AC#5 池帽与桥帽的关系要现读定性**：修完之后**池＝4 是不是仍然必要**？池若可以写 8，要具名写出"8 枚孩子在 4 枚工具许可上**轮转**、单轮并发仍是 4"这条不违反 `PLAN.md:2849` 的理由——**这一格不许靠"反正测试绿了"过**。
- [ ] **AC#6 整包终态读数**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 跑到终态＋逐名比红名集合（**`-run` 单跑不算终态**，这条规矩来自 `A415`）；gofumpt 用 `"$GOPATH/bin/gofumpt.exe"`（v0.12.0）。

## 禁区

⛔ **不动 `PLAN.md:2849` 那个 4**，也不动 `docs/specs/SPEC-01:134`／`:26`／`SPEC-05:70`（改它＝契约改动，要 owner 一句话）；`thresholds.go`／golden／`allowlist.txt`／C17 白名单既有名字一字不动；不新增导出方法名（票 194 是名册补齐的射程）；`frontend/**`／`design/**` 零写（连内容都不转述）；
⛔ **不许用"把池调到 8"来交这枚票的差**——那是把"占着许可干等"从 4 枚放大到 8 枚，症状更糟；
⚠ 写面撞车：本票与 **票 201**（正在飞的写腿）、**票 221** 都碰 `internal/tools/subagent_197.go`／`cmd/wisp` ⇒ **一律串行**。

---

## 收 `222-r1`（编排者，2026-09-29 12:46，锚点 `beaeaeba`；上面原句一字不抹，本节只追加与更正）

**盘上事实（我自己跑的尺，非转述）**：三枚提交＝`92a4b5b7`（只 `docs/evidence/s1/222-spawn-permit-release-r1.md`，现读 23,258 字节）／`bfc55b5b`（只 `internal/tools/bridge.go`＋`internal/tools/subagent_197.go`＋`internal/tools/subagent_222_test.go`）／`beaeaeba`（只那枚证据件＋测试件）；文件清单与消息**逐枚相符**（`A429` 那起共享索引事故之后，这批每发 commit 都带了 pathspec）。测试件工作树 19,956 字节＝HEAD 内 19,956 字节 ⇒ `dd7c447f` 那次"被半路收走"的中间态**已闭合**。

**我逐枚现读复认的冻结面**：`bridge.go:24` `const MaxToolConcurrency = 4`、`:28` `const DefaultToolTimeout = 30 * time.Second`、`subagent_197.go:85` `MaxConcurrentSubagents = 4`；`git diff beaeaeba --stat -- PLAN.md docs/specs/ internal/risk/thresholds.go` ＝ **空**。⇒ 上面那句"真正的一条路（丙）……桥仍是 4、`PLAN.md` 一字不动"**已被兑现**。

**交接点（形状是真的，位置有理由）**：`subagent_197.go:356` `giveBackWhileWaiting(ctx)` 落在 `child.RunAsync`（`:335`）之后、两条"已派生"文案（`:357`／`:358-360`）之前；桥侧六枚新载体**全未导出**（`bridge.go:657`／`:664`／`:670`／`:682`／`:684`／`:694`，接线 `:438`／`:453`）⇒ 零新增导出名、`SubagentDeps` 仍 5 枚字段。

### ⚠ 两处更正（它顶回我，我现读 `bridge.go:473-486` 复现，认账）

1. **上面 §现量 推出的那句父侧文本，父任务其实收不到**。工具自己写的是 `subagent_197.go:389` 那段「父任务这一侧已经不等了……」，但 `bridge.go:478` 的 `errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil` 会在 `Execute` 返回之后把它**换成**「工具 task.spawn 超时（Nms），已协作式中止」（`:481-482`）。⇒ **形状不变**（仍 `IsError`、父侧仍拿不到孩子结论），**措辞要按这句改**；本仓判据没有一条依赖那句字面，所以不影响任何 AC。
2. **行号漂移**：上面 AC#4 引的 `:362-363`（"可以单独停它"）现读在 **`:389`**；Description 那处现读在 **`:196`**（票面写 `:189`）。AC#4 本身仍归票 221（同撞 `subagent_197.go`），只更正坐标。

### 逐格状态（⛔ 全不翻勾，等 `222-v1`＋我的整包复跑）

| AC | 实现者自述 | 我今天的读数 |
|---|---|---|
| AC#1 | 做了 | **未翻**：三枚用例真名在盘（`subagent_222_test.go:362`／`:434`／`:469`），"孩子工具面走真桥"这半条我从码上认（`ParentTools: h.bridge`）；**但修前红与修后绿两发都是实现者自跑**，我未复跑 |
| AC#2 | 做了 | **未翻**：同上；`giveBack` 走 `sync.Once`、终态 `len(sem)=0` 兼作"双扣"那形的尺——形状我读了，读数我没跑 |
| AC#3 | **半做** | ⛔ **只能算半格，且它自己声明剩半支等 owner**：丙形消掉"孩子被自己的父饿死"；**等待仍受 per-tool 超时约束**（孩子跑得比 `cfg.Agent.PerToolTimeoutMS`／未配置 30s 久⇒父任务仍收到桥的超时文案，孩子继续跑并正常落名册⇒**名册对、父任务错**这枚分裂仍然存在）。乙＝给这枚调用开 C22 豁免，`bridge.go:30-31` 逐字 "nothing here can switch enforcement off" ⇒ **契约面，谁都不自裁**；甲＝`task.spawn` 立即返回句柄＋之后用现成 `task.output` 读结论 ⇒ **要先落 `A##`**。**大白话：派了子任务、父任务报错这件事，今天只修掉一半。** |
| AC#4 | 未做（有意） | 归票 221（`221-c1` 只读普查已派，12:46） |
| AC#5 | 只给定性 | 它量到"池抬到 8 会立刻红 `Test197SubagentPoolNeverExceedsBridgeCeiling` 与 `Test197SubagentPoolCapsAtBridgeCeiling`，且名册会出现 4 枚递不动工具的'在跑'"⇒ 这一格**不许靠测试绿了过**的要求已满足为"没靠绿过"，但仍属〔仅自述〕未复核 |
| AC#6 | 做了 | **未翻**：它的两发原始整包日志我读了（`/tmp/222r1/gates-full.txt`、`gates-full2.txt`，各 67 行）⇒ 两发均 **24 包 ok／FAIL 只有 `internal/ball` 1 例＋`internal/panel` 4 例**，逐字原因全指另一队地界，零新增红；**`internal/risk` 这两发为绿**（6.122s／5.071s）⇒ 我在 `A428` 给的"CPU 争用型假红"归因**多了两次独立样本**。**但这不是我跑的那一发。** |

### 给下一位的两枚待办（顺序有依赖）

1. **等 `226-r1` 交件后**，我本人复跑并落到 `.scratch/wisp/probes/222/r1/gates-full.txt`（**绝不接 `| tail`**）：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"; go test ./cmd/wisp ./internal/... -count=1`，与 `.scratch/wisp/probes/201/r2/gates-full.txt` **逐名比红名集合**（基线＝`internal/ball` 1 例＋`internal/panel` 4 例；`internal/risk` 若再红先按争用复量再归因）。
2. **`222-v1`（非实现者对抗验收）**，判据至少三发变异：① 删掉 `subagent_197.go:356` 那一行 ⇒ 三枚用例必须红（正控）；② 把 `giveBack()` 的 `sync.Once` 拆掉 ⇒ "同一枚许可被扣两次"那形必须红；③ 核 `MaxToolConcurrency`／`MaxConcurrentSubagents`／`PLAN.md:2849` 未被为了并发而说破（**同时执行数 ≤ 4 是反控，不许被改绿来迁就修复**）。锚点交给它自取（`git rev-parse --short HEAD`），⛔ 不许引用 `dd7c447f` 作为 `internal/tools` 的"已核过绿"锚（`A429` 事故那发）。

> owner 09-28 原话（本票起手时挂在 Status 行，交付后移到这里保留）：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面**……这是主流 harness 必做的，不要偷懒」。**这一枚修掉的是"父任务里看不到孩子的结论"那一半；"点进去看各自的流"属票 197／票 220 的地界，不因本票结案。**
