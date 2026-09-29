# 票 222（丙形）——等待孩子的那段时间不占桥的执行许可：r1 交付证据件

- 写腿：`222-r1`。开工定锚 `49abb593`（`git log -1 --format=%h` 现取）；本件首次落盘时刻 `2026-09-29 12:12 +08`，`date` 现跑。
- 射程＝工单 `.scratch/wisp/issues/222-spawn-holds-a-bridge-permit-while-waiting-for-its-child.md` 的**丙形**（等孩子的这段时间不占许可）。甲形（`task.spawn` 立即返回句柄）**未做**，等 owner 批。
- ⛔ 本件不新造规矩：D1–D47／C1–C32／SLO／golden／`thresholds.go`／`PLAN.md`／`docs/specs/**` 一字未动。桥的 `MaxToolConcurrency = 4` 与池的 `MaxConcurrentSubagents = 4` **一枚都没改**（§④ 有两枚常量的修复后现读）。

---

## ① 现读链条（行号全部 2026-09-29 12:0x 现跑 `grep -n`，非引用任何归档快照）

| # | 事实 | 现读出处 |
|---|---|---|
| 1 | 桥在调用工具**之前**取许可、整个调用期间握着 | `internal/tools/bridge.go:438-440`（`select { case b.sem <- struct{}{}: defer func() { <-b.sem }()`）；池 `:171` `sem: make(chan struct{}, ceiling)`；常量 `:24` `MaxToolConcurrency = 4` |
| 2 | `task.spawn` 是桥内工具，`Execute` 里阻塞等孩子 | 注册 `internal/tools/subagent_197.go:218`；阻塞 `:356-368` 那个 `select { case <-ctx.Done(): … case res := <-done: … }` |
| 3 | 等待用的 ctx 带 per-tool 超时 | `bridge.go:449` `ectx, cancel := b.execContext(ctx, timeout)`（建在取许可**之后**）→ `:470` `entry.Tool.Execute(ectx, …)`；`timeoutFor` `:550-555`；`TaskSpawnDecl()`（`subagent_197.go:202-211`）没有 `Timeout` 字段 ⇒ 生产走 `cmd/wisp/run.go:562` 的 `DefaultTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) * time.Millisecond`，未配置＝`bridge.go:28` `DefaultToolTimeout = 30 * time.Second` |
| 4 | 孩子在生产上用同一枚桥递工具 | `cmd/wisp/run.go:581` `ParentTools: rt.bridge` ＋ `subagent_197.go:277` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` |
| 5 | 孩子的 ctx 没有 deadline | `subagent_197.go:327` `context.WithCancel(context.WithoutCancel(ctx))` |
| 6 | 池帽＝桥帽＝4 | `subagent_197.go:78` ＋ `bridge.go:24` |
| 7 | 测试看不见这一形：孩子拿的是假目录 | `internal/tools/subagent_197_test.go:183` `ParentTools: h.dir`、`:186` `Tools: h.dir`、`:159` `dir: &fake197Dir{}`；父侧 `:216` 却是 `h.bridge.Execute`。全文件 `grep -n "Tools: h.bridge"` ＝ **零命中**（现跑） |

⚠ **这七行号＝开工锚点 `49abb593` 的现读**（12:0x 那一发 `grep -n`）。本腿产码在 `bridge.go` 的 `:437-453` 之间插了 4 行、并在 `:629-700` 追加了载体，所以**修后**同一批位置落在：取许可 `:440-442`、`execContext` 调用 `:451`、`entry.Tool.Execute` `:473`、超时前查 `:478`（后查 `:498`）、`timeoutFor` `:553`、`execContext` 定义 `:545-550`。`subagent_197.go` 插了两段注释共 26 行，修后：常量注释块 `:60-85`、那个 `select` 块 `:382-394`、"可以单独停它"文本 `:389`、`Description()` `:196`。两枚常量未动（`bridge.go:24`、`subagent_197.go:85`）。§③ 那张表给的就是修后的行号。

⚠ **对票面 §"推出来的真实形状"的一处现读更正**：父任务真正收到的**不是** `subagent_197.go:360-365` 那段「父任务这一侧已经不等了」文本，而是**桥把它换成了超时文本**——`bridge.go:475-483`（修后 `:478`）在 `entry.Tool.Execute` 返回**之后**先查 `errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil`，命中就丢回「工具 task.spawn 超时（1500ms），已协作式中止」。§② 那发红名的原始读数就是这句。结论形状不变（`IsError: true`、父侧拿不到孩子结论），但**用户/模型看到的字面是桥的超时文案**，票面与 `A420` 那句"收到的错误文本就是父任务这一侧已经不等了"按现读应改成上面这句。

## ② 修复前的读数（三枚新用例，产码＝HEAD 未动，`-count=1`）

⚠ 本节记**最终版用例**对未修码那一发（12:26 现跑）。本腿更早还有一发 12:11 的红（用例第一版，台件里少了"四枚孩子的桥上调用的那枚 gate"），形状相同（前两枚 0.00s 判红、第三枚 1.50s 判红），具名存在过但不作为交付读数。

时刻 `2026-09-29 12:26 +08`（`date` 现跑）。做法：`internal/tools/bridge.go` 与 `internal/tools/subagent_197.go` 用 `git show HEAD:<path> >` 还原成未修状态（修后版本备份在仓外 `/tmp/222r1/*.fixed`，跑完逐字节还原并复跑），**测试件与台件一字未动**，命令：

```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
go test ./internal/tools -run 'Test222' -v -count=1
```

逐名结果：

| 用例 | 修复前 | 用时 | 红在哪一枚读数 |
|---|---|---|---|
| `Test222SpawnConclusionArrivesThroughRealBridgeChildren` | **FAIL** | 3.00s | 正文：父任务全收到桥的超时文案；孩子工具调用起跑时已有 3/4 枚父任务返回 |
| `Test222WaitingParentHoldsNoBridgeSlot` | **FAIL** | **0.00s** | `len(bridge.sem)`＝4 |
| `Test222CeilingStillCapsExecutedCallsWhileParentsWait` | **FAIL** | **0.00s** | 同一枚 `len(bridge.sem)` 读数 |

关键原始行（逐字，节选）：

```
subagent_222_test.go:440: 等待孩子的父任务占着 4 枚桥位，want 0：全桥 4 枚，父任务一边干等一边占位＝票 222 现量链的形状，孩子只能排在 sem 后面，直到父任务自己的 per-tool 超时松绑
subagent_222_test.go:474: 等待孩子的父任务占着 4 枚桥位，want 0（反控的前半：先把许可还不回来）
subagent_222_test.go:377: 孩子 c1 的工具调用起跑时已有 3 枚父任务返回，want 0：它是等父任务的 per-tool 预算到点才挤上桥的
subagent_222_test.go:377: 孩子 c2 的工具调用起跑时已有 4 枚父任务返回，want 0：它是等父任务的 per-tool 预算到点才挤上桥的
subagent_222_test.go:392: 父任务 0 收到的是错误而不是结论："工具 task.spawn 超时（3000ms），已协作式中止"
subagent_222_test.go:410: 结论 "子代理 1 的结论正文 222" 出现在 0 枚父任务回复里, want 1
--- FAIL: Test222SpawnConclusionArrivesThroughRealBridgeChildren (3.00s)
--- FAIL: Test222WaitingParentHoldsNoBridgeSlot (0.00s)
--- FAIL: Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/tools	3.037s
```

读数解释（每一枚都问过了"把修复拿掉它会不会红"）：
- `Test222WaitingParentHoldsNoBridgeSlot` 直接读 `len(bridge.sem)`：等待中的 4 枚父任务**占着 4 枚**，全桥就是 4 枚 ⇒ 孩子一枚都拿不到。这一发 **0.00s 判红**，不靠超时、不靠挂死。
- `Test222CeilingStillCapsExecutedCallsWhileParentsWait` 同一枚读数做前半，后半的反控（同时执行数＝4、宿主调用起跑时 `len(parents) == 0`）在修前被前半直接 `Fatalf` 挡住。
- `Test222SpawnConclusionArrivesThroughRealBridgeChildren` 量的是**内容**：孩子的工具调用只能等父任务放弃之后才挤上桥（`returnedYet` 读出 3 与 4），四枚父任务全收到桥的超时文案、四枚孩子的结论在正文里各出现 **0** 次。

## ③ 修复形状（丙：等孩子期间不占许可；不动任何冻结数字）

产码两枚文件（现读行号＝commit `bfc55b5b` 之后；本件写完后又只改过测试件的 helper 命名，见 §⑦）：

| 落点 | 是什么 | 为什么是这个形状 |
|---|---|---|
| `internal/tools/bridge.go:438-442` | run 取到许可之后把令牌交给一枚 `inFlightSlot`，`defer slot.giveBack()` 仍是兜底释放 | `sync.Once` 保证同一枚令牌只被扣一次：工具提前交还之后，run 返回时的 defer 走同一条 once ⇒ 不会在空信号量上永久阻塞。"修完才挂"这一形另有尺（§④ 最后一枚读数） |
| `internal/tools/bridge.go:453` | `ectx = withInFlightSlot(ectx, slot)` | 与 177 的 `hostPathBox`（`:594-596` 那一族）同形：挂在调用自己的 ctx 上 ⇒ 零新增导出名、零新契约字段、C1 的 `Result` 一字不动 |
| `internal/tools/bridge.go:629-700` | `inFlightSlot` 类型 + `take`(`:664`) + `giveBack`(`:670`) + `inFlightSlotKey`(`:682`) + `withInFlightSlot`(`:684`) + `giveBackWhileWaiting`(`:694`)，**全部包内未导出** | 交还是"我接下来只在等"这一句陈述，只有真的只等的那一枚工具才诚实；今天满足的只有 `task.spawn`。做成导出 seam＝把"何时可以不占许可"交给每个插件自己声明 |
| `internal/tools/subagent_197.go:338-356` | 在 `child.RunAsync` 之后、两条「已派生」文案之前 `giveBackWhileWaiting(ctx)`（`:356`） | 位置是承重读的：看到了这一枚调用的 spawn delta 就等于看到了交还已经发生 ⇒ §② 那枚 0.00s 的红/绿不靠超时。**池位不在这里还**——孩子没 join 位就得留着（`Test197FullPoolRefusesNextSpawnWithReadableReason` 钉的就是那一条） |
| `internal/tools/subagent_197.go:63-85` | `MaxConcurrentSubagents` 的**注释**重写；常量本身仍是字面 4（`:85`） | 旧注释的理由（"spawn 全程占着一枚桥位"）在丙形落地后不再成立，留着就是一枚会说谎的注释。新理由＝天花板约束的是"同时在执行的工具调用"；池要不要写 8 是 AC#5 那道独立决定，不由这枚文件自裁（§⑥） |

交还之后**不重新取得**，理由写在 `bridge.go:641-651` 那段：重新排队会把父任务放回它刚离开的那条队，等待结束时只会复制出同一处阻塞的更隐蔽版本；而交还之后这一枚调用剩下的全是簿记（名册行、C25 盖戳、journal 行），不是能力执行 ⇒ 那句"同时执行的工具调用数＝4"逐字仍然为真。

## ④ 修复后的读数（同名三枚，产码＝`bfc55b5b`，helper 改名后复跑）

时刻 `2026-09-29 12:32 +08`（`date` 现跑）：`go test ./internal/tools -run 'Test222' -count=25 -v`

```
     25 --- PASS: Test222SpawnConclusionArrivesThroughRealBridgeChildren (0.00s)
     25 --- PASS: Test222WaitingParentHoldsNoBridgeSlot (0.00s)
     25 --- PASS: Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)
```

逐枚读数（修后 / 修前，两发都是同一枚用例、同一份台件）：

| 读数 | 修后 | 修前 |
|---|---|---|
| 等待中的 4 枚父任务占着的桥位 `len(bridge.sem)` | **0** | **4**（全桥就是 4 枚 ⇒ 孩子一枚都拿不到） |
| `cap(bridge.sem)`（那枚冻结数字的现读） | 4 | 4（未动） |
| 第 4 枚孩子的工具调用起跑时 `len(bridge.sem)` | **4**＝四枚孩子同时跑在四枚许可上 | 到不了那一步（孩子在 sem 后面排队直到父任务超时） |
| 工具调用自记"起跑时已返回的父任务数" | **0** | **3 与 4** |
| 同时在执行的工具调用数 `maxSeen` | 4，从未超过字面 4 | 4（修前也 ≤4——这枚是反控，不是判红源） |
| 父任务正文 | 各自孩子的结论，且逐枚只出现在 1 枚回复里 | 「工具 task.spawn 超时（3000ms），已协作式中止」×4，四枚结论各出现 **0** 次 |
| 名册池位 `InFlightSubagents()` 等待期间 / 终态 | 4 / 0（池没被"修没"） | 4 / 0 |
| 终态 `len(bridge.sem)` | 0（once 没被双扣） | 0 |

撞钉名册（工单派单给的那一批）逐名读数与整包终态在 §⑤。

## ⑤ 门与原始输出

全部带 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（少了这步 `cmd/wisp` 会以 0xc0000135 假失败且没有任何 `--- FAIL` 行）。每发之前现跑 `date`。

| 门 | 时刻 | 原始读数 |
|---|---|---|
| `go build ./...` | 12:27 / 12:29 / 12:31 | 空输出（净） |
| `go test ./internal/tools ./cmd/wisp -count=1` | 12:27→12:29 | `ok github.com/CarlosShao/wisp/internal/tools 12.730s`／`ok github.com/CarlosShao/wisp/cmd/wisp 77.622s` |
| `gofumpt -l`（v0.12.0，`$GOPATH/bin/gofumpt.exe`）对 `bridge.go` `subagent_197.go` `subagent_222_test.go` | 12:32 | 空列表（三枚都无需再格式化；第一发曾点名测试件，已 `-w` 改掉后复量） |
| `sh scripts/d22scan.sh` | 12:27 | 阳性对照：`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0`；实扫：`d22scan: clean - no D22 ban violations`，`examined 246 production Go files under internal/ and cmd/`，bans#1-5 internal/=218、cmd/=28、ban#6 frontend/=85、**ban#7 internal/tools/=22**、ban#8 internal/=457（含 `_test.go` 与注释） |
| `go test ./internal/tools -list '.*'` | 12:27 | 测试函数 **161** 枚＝原有 158（同名文件逐名排除法现跑）＋本票新增 3 枚；原 158 枚**一字断言未动**（`git show --stat` 证明本腿只碰三枚文件） |
| 13 枚撞钉预检名册逐名 `-run … -v` | 12:33 | 13/13 `--- PASS`：`TestToolConcurrencyCeilingIsFour`(0.05s) `TestBridgeIsTheLoopToolProvider` `Test197SubagentPoolNeverExceedsBridgeCeiling` `Test197SubagentPoolCapsAtBridgeCeiling` `Test197SpawnDescriptionNamesTheRealPoolCap` `Test197FullPoolRefusesNextSpawnWithReadableReason` `Test197RowExistsBeforeFirstChildModelCall` `Test197SpawnPublishesIdentityRow` `Test197CancelIsPerRowAndNeverCascades` `Test197SubagentHasNoSelfApprovalOutlet` `Test197ChildCannotDeriveSubagent` `Test197ConclusionCarriesTaskOutputTaint` `Test197StreamKeyShapeIsLiteral`（其余各 0.00s） |
| 整包终态 `go test ./cmd/wisp ./internal/... -count=1`（不是 `-run` 单跑） | 12:29→12:32 | **24 枚包 ok／5 枚 FAIL**，逐名列名（全部是派单已告知的历史红，一枚没修一枚没新增）：`internal/ball` → `TestC21TableColourRowsMatchTokensCSS`（`tokens_table_test.go:1468` 逐字 `read design/assets/tokens.css: ... cannot find the path specified`）；`internal/panel` → `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`（`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`）／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（`a second style source appeared in code the bundle actually ships`）／`TestC21DesignTokensFourWayAgree`（同一枚缺失 tokens.css，`tokens_fourway_test.go:441`）。`internal/tools 16.464s ok`、`cmd/wisp 84.766s ok` |
| 同名三枚 `-count=25 -v`（helper 改名之后复量） | 12:32 | 25 发全 `--- PASS (0.00s)`，见 §④ |

⚠ 一处工具链事实（记给自己也记给下一位）：整包终态这一发跑在 12:29，测试件里 `awaitTokens` 那枚 helper 是 12:31 才改名成 `await222Tokens` 的——改名只动标识符、不动任何断言，改名后 25 发与三包门均已复跑；第二发整包终态（改名后）见本件末尾的 §⑤bis。

## ⑤bis 改名后的整包终态复跑（12:32→12:34，`-count=1`，原始日志 `/tmp/222r1/gates-full2.txt` 67 行，仓外临时件只建不删）

`go test ./cmd/wisp ./internal/... -count=1` ⇒ **24 枚包 ok／2 枚包 FAIL（共 5 枚用例）**，逐名：

```
ok  	github.com/CarlosShao/wisp/cmd/wisp	83.788s
FAIL	github.com/CarlosShao/wisp/internal/ball	0.207s        （TestC21TableColourRowsMatchTokensCSS）
FAIL	github.com/CarlosShao/wisp/internal/panel	2.516s       （TestApprovalCardViewJSONKeysMatchFrontendTypes／
                                                               TestComposerContractTypesMatchFrontend／
                                                               TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme／
                                                               TestC21DesignTokensFourWayAgree）
ok  	github.com/CarlosShao/wisp/internal/tools	15.129s
```

其余 21 枚 ok（agent 2.801s／approval 0.459s／audio 15.822s／buildinfo 0.038s／config 0.969s／llm 32.342s／adaptertest 2.052s／anthropic 14.703s／golden 0.090s／openaichat 21.342s／openairesponses 18.895s／memory 13.221s／models 6.916s／observe 3.701s／perm 0.651s／plugin 0.122s／proc 1.557s／projctx 0.286s／risk 5.071s／secret 0.515s／statemachine 0.239s／winsec 12.025s）。

与派单给的"历史 5 枚红"名册**逐名相同、零新增**；五枚的逐字原因全指另一队地界（`read design/assets/tokens.css: ... cannot find the path specified` 两枚、`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`、`a second style source appeared in code the bundle actually ships`），本腿一枚没修、也没把它们改绿。前一发同命令的读数（12:29→12:32，`/tmp/222r1/gates-full.txt`）＝24 ok／同样 5 枚 FAIL，两发名册一致。

## ⑥ 没做的格与卡点（逐格，含"卡在哪一句"）

| AC 格 | 状态 | 卡点（具名） |
|---|---|---|
| AC#1 生产形状第一次有尺 | **做了** | `subagent_222_test.go` 三枚用例的孩子工具面＝真桥（`ParentTools: h.bridge`，同 `cmd/wisp/run.go:581`）；修前/修后两发读数在 §②／§④，修前红不靠挂死（两枚 0.00s 内容判红、一枚 3.00s 正文判红） |
| AC#2 占着许可干等被永久钉住 | **做了** | 正读＝`len(bridge.sem)` 在等待点＝0（修前 4）；反控＝同时执行数 ≤ 字面 4 且宿主第 5/6 枚调用仍在排队；终态 `len(sem)=0` 兼作"同一枚许可被扣两次"那形的尺 |
| AC#3 父侧结论通道要么真要么明说不做 | **半做，剩半支等 owner** | 丙形把"孩子被自己的父饿死"这一半消掉了；但**父侧等待仍受 per-tool 超时约束**：孩子跑得比 `cfg.Agent.PerToolTimeoutMS`（未配置 30s，`bridge.go:28`）久时，父任务仍收到桥的超时文案、孩子继续跑并正常落名册。要做乙（等待不受 tool 超时约束）就得给这一枚调用一枚 C22 豁免或独立预算——`bridge.go:30-31` 逐字 "nothing here can switch enforcement off"、`execContext`（修后 `:545-550`）是契约面 ⇒ **不自裁**。甲（立即返回句柄）票面自己写了要先落 `A##` ⇒ **等 owner**。⚠ 所以"派了个子任务、父任务报错"这一形**今天只修掉一半**，剩下一半的开关在 AC#3 那一格 |
| AC#4 不许留对模型许诺了但没接的 | **未做，有意不做** | `subagent_197.go:196`（Description："可以用 task.cancel 单独停它"）与 `:389`（等待失败文本："可以单独停它"）两处一字未动。`task.cancel` 不存在＝票 221 的射程，且两枚票同撞 `subagent_197.go` ⇒ 按串行规矩交给 221 |
| AC#5 池帽与桥帽的关系要现读定性 | **定性已给（下面 §⑥bis），常量一枚没动** | 池抬到 8 需要 owner 拍板＋两枚现存钉的重写，不在本腿射程 |
| AC#6 整包终态读数 | **做了** | §⑤ 那行：24 ok／5 枚历史红逐名，零新增红；gofumpt v0.12.0 与 d22scan 均净 |

## ⑥bis AC#5 的现读定性（不靠"测试绿了"过）

- **修前**池＝4 是被迫的：`subagent_197.go` 旧注释的理由是"一枚在跑的孩子＝一枚被父任务占住的桥位"，池写 8 只会让 4 枚孩子排在桥里而名册说它们在跑。
- **修后**这枚关系换了性质：桥的 4 枚许可约束的是"**同一瞬间在执行工具调用的孩子数**"，不再约束"在等的父任务数"。于是——
  - 池若写 8：**同一瞬间仍最多 4 枚孩子在执行工具**，另外至多 4 枚在 `bridge.go:439` 那枚 sem 后面轮转。这不违反 `PLAN.md:2849` 那句理由（它防的是"LLM 一次发 50 个 tool call 打爆机器"，同时执行的数量没有变），也**不需要动那枚 4**。
  - 但今天抬到 8 会立刻红两枚现存钉：`Test197SubagentPoolNeverExceedsBridgeCeiling`（`:403`，比的就是 `MaxConcurrentSubagents > MaxToolConcurrency`）与 `Test197SubagentPoolCapsAtBridgeCeiling`（`:441`，它要求 N 枚孩子**同时 started**）；更要紧的是名册会同时显示 8 枚「在跑」而其中 4 枚此刻递不动工具——那正是 197 立这两枚钉要拦的"名册与机器不是一回事"。⇒ 抬池值＝owner 一句话＋那两枚钉的重写，**不是改一个常量**。
- **一句话结论**：池＝4 在本腿交付之后**仍然必要**（钉在测试上、理由仍成立），但它不再是"等孩子就得占许可"逼出来的假诚实；丙形把"要不要 8 枚孩子"从契约问题降级成了名册诚实度问题。

## ⑦ 新增名字具名清单（要新增导出名/字段名的话，就在这儿）

- **导出名：零枚新增。** 没有新导出类型、方法、字段、函数；`Bridge` 的公开面（`Execute/Tools/Registry/Paths/CloseTask`）一字未动；`SubagentDeps` 仍是 5 枚字段（`Test197SubagentHasNoSelfApprovalOutlet` 那枚字段枚举守卫逐字仍绿，§⑤ 的 13 枚名册里有它）。
- 包内未导出新增（`internal/tools`）：`inFlightSlot`（type，`bridge.go:657`）、`(*inFlightSlot).take`（`:664`）、`(*inFlightSlot).giveBack`（`:670`）、`inFlightSlotKey`（type，`:682`）、`withInFlightSlot`（`:684`）、`giveBackWhileWaiting`（`:694`）。
- 测试件内新增（只在 `_test.go`）：`parent222`、`probe222Tag`、`h222CeilingLiteral`、`h222PreFixBudget`、`h222SafetyBound`、`h222`、`probe222`、`probeRun222`、`child222`、`newH222`、`(*h222).wireSpawn`、`(*h222).launchParents`、`(*h222).parkParents`、`await222Tokens`、三枚 `Test222*` 用例名；工具名 `probe.work222`（只在测试里注册）。
- 逐枚对照 `internal/panel/l2_grant_boundary_test.go:1882-1886` 那 12 词根×两拼（`outcome`/`allow`/`allowOnce`/`approved`/`grant`/`verdict`/`decision`/`decide`/`bypass`/`override`/`permit`/`authorize`）与 `:192-195` 那批路由名词根：**一枚都不含**（新名全部落在 slot/handoff/giveBack/inFlight/probe/park/parent/child/waiting 这些词根上）。那枚冻结件一字未动，`internal/panel` 里也没有新增可解码结构体（`:1855-1860` 那条 registry 腐蚀守卫因此碰不到）。
- d22scan ban #7 的正则 `"(spill|internal[._-][a-z0-9_.-]+)"` 与 `probe.work222` 不命中（实扫 clean，见 §⑤）。
- **想新增却没敢造的：无。** 若下一步要做 AC#3 乙（等待不受 tool 超时约束）或"等待结束后重新排队拿回许可"，那都要新增导出面或契约位，属要人批的那一格，本腿没碰。

## ⑧ 留给下一位的三件（本腿没动手、但已被现读证伪的文字）

1. `subagent_197_test.go:398-402` 那段注释仍在说"one in-flight spawn holds one bridge slot for its child's whole life (bridge.run keeps the semaphore held across entry.Tool.Execute)"——丙形落地后这句**不再为真**。它挂着的**断言**没受影响（§⑤ 的 13 枚名册逐名绿），所以本腿没改那枚文件（禁区：别人的测试件断言一字不动；注释归下一枚合法写腿顺手改）。
2. 同文件 `:516-525` 那段"now that the pool EQUALS the ceiling, N real children also occupy all N bridge slots, so a (N+1)-th call could only queue behind them"同样过期：修后第 5 枚调用能直接到达工具并拿到硬拒理由。断言本身照旧绿。
3. `internal/tools/subagent_197.go:196`/`:389` 的「可以单独停它」＝票 221 的射程（AC#4），同文件串行。
