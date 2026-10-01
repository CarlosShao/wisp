# 票 246 AC#7 — 非实现者验收腿 `246-v2` 裁决表

- 腿：`246-v2`（对抗验收，只裁 **AC#7** 一格；不写产码、不修任何东西、**不勾任何 AC 框**）
- 起手 HEAD：`aec06d24`（`git rev-parse --short HEAD` 现跑，2026-09-30 23:43；含台账 `A486` 那枚 `45202bb0`，满足"起手＝`A486` 或更新"）
- 被验收的产码：`bb37fac2`／`3edc11d4`／`d8f3b27e`／`2071f59e`／`8af4347d`（写腿 `246-r2`，死于 150 轮上限）
- 判据原文出处：`.scratch/wisp/issues/246-*.md`「AC#7（新补）」那一格
- 上一枚验收腿：`246-v1`（裁 AC#1..AC#6，已被翻勾，不在本表射程）
- ⛔ 本表不引 `246-r2` 的自述当凭据：它表内 §5（突变与正控）与 §6/§7 未写完，那些格由本腿**自己现跑取数**。
- ⛔ 本表不含任何凭据值（密钥一律只写变量名）。

## §1 门禁读数

### §1.1 进程与前置（本腿现跑，2026-09-30 23:43 +08）

| 尺 | 命令（逐字） | 读数 |
|---|---|---|
| 球调试进程占位 | `tasklist //FI "IMAGENAME eq balldebug.exe"` | `INFO: No tasks are running which match the specified criteria.` ＝ **0 枚** |
| Wisp 进程占位 | `tasklist //FI "IMAGENAME eq wisp.exe"` | 同上一行一字不差 ＝ **0 枚** |
| 起手 HEAD | `git rev-parse --short HEAD` | `aec06d24`（含 `45202bb0 ledger(A486)` ⇒ 满足"起手＝`A486` 或更新"） |
| 工作树（本腿地界） | `git status --porcelain -- cmd internal` | **0 行**（起手即干净；脏项全在 `.gitignore`／`design/**`／`probes/**`／别人名下两枚 evidence 文件，一枚未动） |
| 真机用例前置 | `GOFLAGS= go build -o build/wisp.exe ./cmd/wisp` | rc=**0**；`build/wisp.exe` mtime `23:43:42` > `cmd/wisp/run.go` mtime `23:37:24` ⇒ 不是拿旧 exe 读数 |

### §1.2 起手基线（本腿本人这轮取，⛔ 不引编排者的数当凭据）

命令逐字：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1`

| 项 | 读数 |
|---|---|
| 起跑 | `2026-09-30 23:44:09 +0800` |
| 终态 | `2026-09-30 23:47:06 +0800` |
| 结果 | **`ok github.com/CarlosShao/wisp/cmd/wisp 173.655s`＝全绿，`--- FAIL` 计数 0** |
| rc | 0 |

口径对照（编排者 23:37–23:41 那发，锚 `18270ad7`）＝ `ok ... 209.743s`。**两发都是"各自那一次的读数"**，不是上一版 vs 新版的关系；本腿只把自己这发当凭据。

PATH 坑本腿复认过它的形状：不带那两枚目录时 `go test` 报 `exit status 0xc0000135`、`0.0xxs`、**一行 `--- FAIL` 都没有**＝一枚用例都没跑（仪器坑，不是红）。

### §1.3 起手红名册归因（AC#5 的同一把尺，本腿重跑）

起手基线 **0 枚红** ⇒ `A482` 那枚 `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 在本腿起手这一发**未复现**。它是 boot 时序 flake，命中率不因这一发改变；**归口＝票 228（AC#8）**，本腿不顺手修、不当回归报。

（本节以下随各攻推进补满；交付时本表每一节都必须是有内容的格子，⛔ 不许留空段与占位句——本票已连续两枚腿死在最后两节，那是我最在意的失败形状）

## §2 逐格判语（AC#7）

**结论：AC#7 成立——附两条具名条件，不退回。**

判据三截，逐截给凭据：

1. **「常驻那条腿真起一条任务管线」＝成立。**
   凭据＝能力尺，不是文案：`startResidentTaskSource` 有 **1 枚非测试调用者** `cmd/wisp/resident_windows.go:177`，那枚函数在 `main.go:66` 的无参数分支里（＝票面说的"双击图标起来的那个进程"）；它经**全仓唯一一处装配** `assembleRuntime`（`resident_task_source_windows.go:265`）把 `ra.gate`／`ra.ui`／`ra.cards`／`ra.root` 注进去，再由 `src.run.execute`（`:433`）起真环路，卡片走 `run.go:716-717` 的 `Gate: rt.gate` 进桥。**MUT-4 把这枚调用改成 `var src *residentTaskSource` ⇒ 两枚真机用例立刻红**（`resident_task_source_live_246_windows_test.go:120`「never raised a card」＋`:255`「no streamed text ever arrived」）⇒ 四步读数的因果线确实挂在这枚调用者上，不是别处造出来的。
   本腿另现读 `objdump -p build/wisp.exe` ＝ `Subsystem 00000003 (Windows CUI)` ⇒ 今天这枚构建是**控制台子系统**，双击进程真有一枚控制台输入缓冲，`interactiveStdin()` 在那条路上返回 `os.Stdin`＝裁定 2.1 指定的正路**今天是通的**（这一枚读数同时把 §5 第 6 条那笔"票 244 落地后会静默断掉"的前瞻账圈定为**未来**问题而非今天的缺陷）。

2. **「一发真机走通四步」＝成立，但只在〔接缝注入〕这一栏成立。**
   凭据＝§3.C 那 6 发（`-count=1` 一发＋`-count=2` 两发；0 skip；卡片编号每发不同）。〔真模型〕那一栏本机**为空**（23:44 现读 `%APPDATA%\wisp` 无 `config.toml`、无 secrets）。
   ⇒ 任何人（含编排者翻勾时）把这一格写成"任务管线已跑通"而不带栏头，就是**本表反对的那句话**；票面 AC#2 那句"可见证据只到日志＋状态位"与这里同形，翻勾措辞请照 §3.C 的两栏结构写。

3. **⛔「不许用测试构造的任务源冒充」＝没有触犯；派单那条"被造出来就该退回"的规矩在这一格不触发。**
   我的判据是一枚可操作的问句，不是观感：**把产码里那枚消费者删掉，四步读数会不会死？** MUT-4 的答案是"会，两枚真机用例立刻死"。若任务源真是测试在进程内构造的（`assembleRuntime` 直接由用例调、或 `AskOnTaskRoot` 由用例调），MUT-4 那一发**应当全绿**——它没有。
   这一族里被注入的只有**两枚输入**，且都是 AGENTS.md §1.3 点名的接缝：① 一发任务文本走 `WISP_TEST_TASK_TEXT`（只在 `WISP_ENV=test` ＋ 数据根就是台件那枚时被受理，MUT-6 证明这枚谓词是活的；拒绝侧由 `TestAC246DevLegIgnoresTheTestTaskInjection` 钉住）；② 模型字节走 `internal/llm/golden` 的 C5 golden SSE 经本地 `httptest` HTTP。**被注入的从来不是机构本身**：门、界面、球窗、Esc 借键、否决路由、取消、退出十步——全在那枚真起的 `build/wisp.exe` 里，由真第二个 OS 进程真按一次 Esc 读出（`keydown 1->0->1`，且 `-steal` 正控先跑）。
   还有一条反证值得记：这张网**没有把自己打扮成正路**。受理注入时产码逐字打印「这是只在本机测试台件里受理的一发任务文本，本机没有可交互控制台」，boot 报告的 `任务来源` 字段在该分支取值「仅本发的测试注入位」（`resident_task_source_windows.go:502-513`）＝它自己把这一栏标成了注入栏。

**两条具名条件（不成立的部分，具名、不掩盖）**：

- **条件①｜真模型四步读数＝欠账。** 要 owner 本人往 `%APPDATA%\wisp` 录入凭据后重跑同一对用例。本腿不代填、不碰任何 key 值。这一条不销，AC#7 的"能干活"就只证到接缝层。
- **条件②｜裁定 2.1 指定的正路（真控制台键盘那一支）零自动覆盖。** 本腿现读确认：没有任何用例（默认层或 winlive）把一行 `task <文本>` 送进 `runConsoleLoop`——`grep -rn "runConsoleLoop" --include=*_test.go cmd/` ＝ **0 处**。今天对它只有**反向**覆盖（MUT-1：中和闸门 ⇒ `TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry` 在 `:307` 红，红句里冒出的正是那句谎）＋ `246-r2` §1.3 的手测。"人真坐在终端前敲一行起出一张卡"这一形**没有读数**。是否为此立一枚伪控制台台件，是编排者的裁（§5 第 2 条），本腿不自扩判据、也不因此把 AC#7 判不成立——因为 AC#7 票面写的是"一发真机走通四步"，没说"必须走控制台那一支"。

**顺手复认 AC#8 的归口（不在我射程，只裁归口对不对）**：**对。** 用例名 `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 是票 228 的名字，落点那段启动窗口（`resident_windows.go:106-108` 的 `bootExit := make(chan os.Signal, 1)` ＋ `signal.Notify` 排在 `installLogSink` 之后）本票五枚提交一行未动（尺＝`git show --numstat` 逐枚：无一枚列出该文件该窗口的行）。本腿**没有**跑 `-count=10`（那是票 228 的判据），且本腿起手那一发整包 0 红，既不构成复现也不构成推翻。**AC#8 保持未勾、保持归票 228。**

## §3 六攻读数（A 生产调用者／B 一进程一枚门／C 四步真机读数／D 定向突变／E 重锚的钉／F 零新依赖边）

### §3.A 生产调用者（能力，不问文案）

尺（本腿现跑，⛔ 不是"import 了就算装配"，问的是**跑着的那条腿会不会真调**）：

```
$ grep -rn "startResidentTaskSource\|newResidentApproval\|AskOnTaskRoot\|askConfirmation" --include=*.go cmd internal | grep -v _test.go
```

| 符号 | 定义处 | **非 `*_test.go` 调用者** | 这条腿跑起来会不会真调 |
|---|---|---|---|
| `runResident` | `resident_windows.go` | `main.go:66`（无参数分支，＝双击那个二进制） | 会 |
| `newResidentApproval` | `resident_approval_windows.go:105` | **1 枚**：`resident_windows.go:123` | 会 |
| `startResidentTaskSource`（AC#7 新增的入口） | `resident_task_source_windows.go:215` | **1 枚**：`resident_windows.go:177` | 会（无控制台／注入被拒时它自己 return nil，见 §3.C） |
| `residentTaskSource.submitTask` | `:401` | **2 枚**：`:322`（注入分支）、`:363`（控制台 `task <文本>` 动词） | 会 |
| `runConsoleLoop` | `:345` | **1 枚**：`:319`（`observe.Default.Spawn` 的协程体） | 仅当真有控制台；测试起的子进程进不到这一支（见 §5 第 3 条） |
| `ra.vetoByEsc` | `resident_approval_windows.go:171` | **1 枚**：`resident_windows.go:136` 作函数值注入球宿主 | 会 |
| `ra.cancelTaskRoots` | `:283` | **1 枚**：`resident_windows.go:157`（D38(e) 第 3 步钩子） | 退出时会 |
| `src.closeDoor` / `src.drainAndClose` | `:458`／`:478` | **各 1 枚**：`resident_task_source_windows.go:308`／`:309`（第 1／第 7 步钩子） | 退出时会 |
| `src.taskPosture` | `:502` | **2 枚**：`resident_windows.go:187`／`:201` | 会 |
| `ra.consoleVetoChannel` | `:520` | **1 枚**：`resident_task_source_windows.go:288` | 仅控制台分支 |
| **`AskOnTaskRoot`** | `resident_approval_windows.go:261` | **0 枚**（只有注释 `resident_task_source_windows.go:13` 提到它） | **不会** |
| **`askConfirmation`** | `:205` | **1 枚**：`:264`，调用者是 `AskOnTaskRoot` ⇒ **传递不可达**（从 `main` 走不到） | **不会** |

**本腿对最后两行的处置（不算缺陷、但要写清为什么）**：AC#7 票面把缺的那一跳量成"`AskOnTaskRoot`／`askConfirmation` 产码调用者＝0 ⇒ 没有任务会去举卡"。本腿现读到：**举卡这一跳今天真的走了，但走的是另一条既有通道**——`startResidentTaskSource` → `assembleRuntime`（`resident_task_source_windows.go:265`）→ `run.execute`（`:433`）→ agent 环路 → `tools.Bridge`（`run.go:716 Gate: rt.gate`、`:717 Cancel: rt.gate.ToolsCancelBus()`）→ 注入的那枚 gate → `ballCardUI.Prompt`。也就是说卡片现在由**模型那一次真工具调用**举起，而不是由宿主自己 `askConfirmation` 举起。`askConfirmation`／`AskOnTaskRoot` 是 246-r1 那一段为"宿主发起的卡片"留的接缝，AC#7 没用它、也没把它删掉 ⇒ 它今天仍是**产码里两枚零真消费者的接缝**。本腿判：AC#7 要防的结局（"门在、卡进得来、还没有任务会去举它"）**没有发生**（凭据 §3.C 四步读数），所以不构成退回理由；但"两枚接缝零消费者"这件事**没有任何尺在钉**，票面下一格若再拿它当"已经接上了"来读就是过期指认 ⇒ 登记进 §5 交裁。

### §3.B 一进程一枚门

尺（本腿现跑）：

```
$ grep -rn "approval.New(" --include=*.go cmd internal | grep -v _test.go
cmd/wisp/resident_approval_windows.go:109:	ra.gate = approval.New(approval.Options{
cmd/wisp/run.go:586:		rt.gate = approval.New(approval.Options{
internal/agent/approval/doc.go:28://   - compose approval.New(...) as tools.Options.Gate (replacing NoGate);
```

⇒ `cmd/wisp` 产码里 **2 处** `approval.New(`（doc.go 那处是注释，不算）。常驻那条路**复用注入的那一枚**，机制在 `run.go:574-602`：`if s.gate != nil { rt.gate = s.gate; rt.ui = nil; rt.liveCards = &nativeCards{h: s.cards} } else { … approval.New(…) }`——注入分支**不进 else**，所以同一个进程里不会同时存在两枚门。装配仍是唯一一处（`assembleRuntime`，`resident_task_source_windows.go:265` 是本票新增唯一一处调用带注入 gate 的装配）。

**"注入的门被无视"今天会不会响**：本腿做了 MUT-C（`run.go:574` 改 `if s.gate != nil && false`），**两枚用例点名红**，见 §4 MUT-C ⇒ 这件事有人钉，不是装饰。

顺带本腿现读的接缝默认值（实现者 §7.4 那两条代价，本腿自己核过、判它说得对）：`approval/gate.go:139-145` 的钳位 `win <= 0 → DefaultL1Window`；`approval/queue.go:107` `DefaultApprovalTimeout = 300s`、`:115` `DefaultL1Window = 3s`、`:122` `MaxL1Window = 3s`。⇒ 常驻注入的那枚门没吃 `cfg.Risk.L1WindowSec`/`ConfirmTimeoutSec`，跑 3s／300s；**窗口这一项即便接了 config 也钳在 3s，所以只有超时那一项是真差异**——与它表内那句一致。

### §3.C 四步真机读数（分两栏）

**〔接缝注入的四步读数〕**——本腿本人这轮现跑，⛔ 未引 `246-r2` 的自述：

命令逐字（两次）：
```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= \
  go test -tags winlive ./cmd/wisp -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v -count=1
PATH=… 同上 … -count=2
```

第 1 发（`-count=1`，起跑 23:52:03／终态 23:52:37，`ok … 31.287s`，rc=0，**0 枚 skip、0 枚 FAIL**）逐字读数：

| 四步中的哪一步 | 出自哪枚用例 | 逐字读数 |
|---|---|---|
| **① 起任务** | `TestLive246ResidentPipelineRaisesACardAndEscVetoesIt`（用例文件 `:124-127`） | 断言 `taskEntryInjectedClaim` 在场并通过；随后同一枚用例打印 `AC#7 LIVE steps 1-2: wisp: 卡片挂起：L1 fs.write（编号 cf73a8eb-67f4-4229-8c21-e685f35130b5）` |
| **② 举卡** | 同一枚（`:116` 轮询 `cardRaisedClaim`） | `wisp: 卡片挂起：L1 fs.write（编号 cf73a8eb-67f4-4229-8c21-e685f35130b5）`；台账侧 `esc_borrowed=true`、`orb_state=Confirming`（`:172-183` 断言，通过） |
| **③ 按 Esc 生效** | 同一枚（`:130-141`＋`:224-225`） | 第二个进程 `observer keydown 1->0->1 (steal control 0)`；审计行逐字 `audit: approval: ANSWER-VETO corr=cf73a8eb-67f4-4229-8c21-e685f35130b5 tool=fs.write channel=esc decision=veto`；否决句 `按 Esc 否决了卡片` 轮询到；被否决那发 `fs.write` 的目标文件 `os.Stat` 仍为 not-exist（`:144-147`） |
| **④ 任务被取消（D43 第 22 行那一形＝取消那一发调用）** | 同一枚（`:151-160`） | `wisp run: 任务 cf73a8eb-67f4-4229-8c21-e685f35130b5 结束（completed，2 轮，1 次工具调用，成本 0 CNY（未计价））` |
| **④ 任务被取消（D38(e) 第 3 步那一形＝取消在跑的任务）** | `TestLive246ExitCancelsARunningTask`（`:267` 断含 `结束（cancelled`，`:284` 打印最后一行） | `wisp run: 任务已取消`；且第 3 步完成记录 `resident-approval: 退出第 3 步完成` 落盘（`cancelStepBookedMsg`）、第 7 步 `退出第 7 步完成` 落盘 |
| 本进程 owns 的两枚退出步 | `TestLive246ResidentPipelineRaisesACardAndEscVetoesIt:188-193` | boot 报告含 `1:scheduler-close` 与 `7:flush-logs-close-db`（断言通过） |

第 2 发 `-count=2`（起跑 23:53:35／终态 23:54:39，`ok … 61.093s`，rc=0）＝**4 发全 PASS、0 枚 FAIL、0 枚 SKIP**，四步逐条重复成立：

| 发次 | 卡片编号 | ② 举卡 | ③ Esc | ④（D43 形） | ④（D38e 形） |
|---|---|---|---|---|---|
| 1 | `0b848087-…` | `wisp: 卡片挂起：L1 fs.write（编号 0b848087-…）` | `keydown 1->0->1 (steal control 0)` ＋ `ANSWER-VETO … channel=esc decision=veto` | `结束（completed，2 轮，1 次工具调用…）` | `wisp run: 任务已取消` |
| 2 | `dbd643a9-…` | 同形，编号不同 | 同形 | 同形 | 同形 |

⇒ **本腿合计 6 发真机（`-count=1` 一发＋`-count=2` 两发×两枚用例＝各自 3 发），零 skip、零 FAIL、编号每发不同**（不是同一张卡被读第二次）。

**〔真模型四步读数＝欠账〕**：**本机拿不到**。尺＝本腿 23:44 现跑 `ls "$APPDATA/wisp"` ⇒ 目录内**没有 `config.toml`、没有 secrets 子目录**（与编排者 22:5x、`246-r2` 23:2x 两枚读数同形，本腿独立复认一次）。上面那一族跑的是 `internal/llm/golden` 的 golden SSE 经 `httptest` 本地 HTTP + 真 openai-chat 适配器 + 真 DPAPI blob（密钥只以本包既有常量 `fakeStoreKey` 的名字出现，**本表不写其值**）＝AGENTS.md §1.3 允许的 **C5 接缝**。**这一栏为空是事实，不是通过。**

### §3.F 零新依赖边

尺（本腿现跑两发名册，逐名比差集）：落地前＝`cf540f8a`（`246-r2` 起手 HEAD，AC#7 五枚提交之前）用 `git archive` 导到**仓库外**临时目录跑（⛔ 不在仓内建 worktree/checkout），落地后＝HEAD 本树。

```
$ cd <仓外临时目录>  && GOFLAGS= GOOS=windows GOPROXY=off go list -deps ./cmd/wisp | wc -l   # 前：276
$                                          GOFLAGS= GOOS=windows GOPROXY=off go list -deps ./cmd/wisp | wc -l   # 后：276
$ diff <(sort deps-base.txt) <(sort deps-head.txt)   ⇒ 无输出（IDENTICAL）
$ comm -13 … ⇒ 新增 0 名      $ comm -23 … ⇒ 减少 0 名
```

⇒ **名册 276 → 276、逐名相同、差集两侧皆空**＝零新包级正向边成立。两枚名册原文入盘：`.scratch/wisp/probes/246/v2/deps-base.txt`／`deps-head.txt`。

### §3.E 那枚重锚的钉（票 127／130 的字面量地标）

**逐字对上**（两边现读）：

| 侧 | 文件:行 | 逐字 |
|---|---|---|
| 钉（改锚后） | `cmd/wisp/resident_sink_nail_127_windows_test.go:119` | `residentReachedLoop = "resident event loop running"` |
| 产码真打印 | `cmd/wisp/resident_windows.go:200` | `fmt.Printf("wisp: resident event loop running (task source: %s); %s (Ctrl+C exits cleanly)\n", …)` |

`grep -c "resident event loop running" cmd/wisp/resident_windows.go` ＝ **2**：`:190` 是注释、`:200` 是唯一一处真打印 ⇒ 地标句子全仓产码只此一处，与改锚前"empty event loop running"的唯一性同形。during-boot-exit 那一支（`:197-198`）打印的是 `"…arrived during boot; the event loop was never entered"`，**不含**该 token ⇒ 它表注释那句"新 token 在两个分支都真、during-boot-exit 分支为假"与本腿现读一致（两支＝`:200` 与 `:186` 之外的 `default` 分支；`:186` 那行 boot 报告不含该 token，所以"进了环"这件事仍只由 `:200` 一处声明）。

**两重职责有没有被改弱**：

1. **改锚只动了 const 块**。尺＝`git show 3edc11d4 -- cmd/wisp/resident_sink_nail_127_windows_test.go`，hunk 头逐字 `@@ -94,7 +94,29 @@`，numstat ＝ **23 增／1 删**，全部落在 `const (…)` 里那枚常量的注释与字面量上；文件其余部分（含 `:415`、`:497`、`:520`、`:530`、`:565` 那些断言）**一行未动**。⇒ 职责 2（"没有早于它的启动记录"，由该文件上面的记录序号断言承载）在文本层面不可能被这发增量改弱——本腿不需要再动 `logsink.go` 去做第七发突变来反证这一条。
2. **反向证（MUT-5，本腿现跑）**：把产码 `:200` 那句改成 `resident loop runs` ⇒ **三枚消费者全部点名红**：
   - `TestAC1ResidentLegInstallsItsLogListenerOnDisk` → `resident_sink_nail_127_windows_test.go:498` 逐字「the child never reached the event loop, so the disk read above proves little.」（＝职责 1：sink 安装记录这块地标的**效力**还在）
   - `TestAC1ResidentLegOutlivesItsOwnLogFailure` → `:539` 逐字「the refusal stopped the app: the child never reached the event loop.」
   - `TestAC228ResidentLegReportsAndBooksItsBall` → `resident_ball_228_windows_test.go:83` 逐字「the ball sentence did not arrive on the same boot as the event-loop sentence ("resident event loop running")」
   还原后同发复跑 **3/3 绿**（`ok … 13.863s`）。
3. **消费者枚数只增不减**：`grep -rn residentReachedLoop cmd/` 现读 **6 处引用**（`resident_ball_228_windows_test.go:82`、`resident_sink_nail_127_windows_test.go:497/:530`、`resident_task_source_246_windows_test.go:303/:370/:430`）＝改锚前 3 枚、改锚后 6 枚。

⇒ **钉还在咬，没有被改弱。** 本腿不因为它"绿着"就判成立——上面第 2 条就是它的红证。

## §4 突变台账（改了哪一行 → 指名哪一枚用例红 → 还原凭据 → 复跑回到绿）

**纪律遵守声明**：**一次只改一枚**；读完红名**立刻**用 `git cat-file blob HEAD:<path> > <path>` 还原（⛔ 全程未用 `checkout`／`restore`／`stash`）；**六枚突变体一行都没有提交**（`git status --porcelain -- cmd internal`＝0 行、`git diff --numstat -- cmd internal`＝0 行，读数见 §7）。还原后每枚都同发复跑被点名的用例证明回到绿。突变体里都带 `MUT-<n>` 行尾注释，便于与本仓历史注释里那些 `MUT-` 指认区分（现读 `grep -rn "MUT-" --include=*.go cmd internal`＝15 处，**全部是早先票号留在测试注释里的历史指认，产码 0 处**，本腿逐条看过）。

| # | 改哪一行（→ 改成什么） | 对应派单裁定 | 点名的红（文件:行 ＋ 逐字红句主干） | 还原凭据 | 复跑 |
|---|---|---|---|---|---|
| **MUT-1** | `cmd/wisp/resident_task_source_windows.go:218`：`console := interactiveStdin()` → `console := io.Reader(os.Stdin)` | 裁定 2.1（真控制台闸门） | `--- FAIL: TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry (14.04s)`；`resident_task_source_246_windows_test.go:307` 逐字「AC#7 RED: the shipped leg with no console said **false** that its entry is off and true that it reached the loop.」；子进程 stdout 里冒出的正是那枚谎句「wisp: 任务入口已启用（控制台键盘）…」（stdin 明明断着） | porcelain 0 行；`grep -n "console := interactiveStdin()"` 复现 `:218` | 该用例随 MUT-4 复跑列回到绿（同族一发，见 §7） |
| **MUT-2** | `cmd/wisp/run.go:574`：`if s.gate != nil {` → `if s.gate != nil && false {` | 裁定 2.2（一进程一枚门） | 两枚点名红：`TestAC246ResidentPipelineAsksThroughTheOneGate` → `:90` 逐字「the assembled pipeline runs on a **DIFFERENT** gate than the resident leg built: 0xdbebb9ec2c0 vs 0xdbebb9ec0b0」；`TestAC246ResidentGateInjectionIsRefusedHalfAssembled` 两枚子用例 → `:161`「assembleRuntime accepted gate without ledger／without surface」。⚠ `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline` 在这一发**仍绿**＝默认层对"注入被无视"的钉是那两枚能力尺，不是文案尺 | porcelain 0 行；`grep -n "if s.gate != nil {"` 复现 `:574` | 同上（该枚绿／红配对本身就是凭据） |
| **MUT-3** | `cmd/wisp/run.go:1046`：`baseCtx := parent` → `baseCtx := context.Background()` | 裁定 2.3（任务根挂在 `ra.root` 下） | `--- FAIL: TestAC246ResidentTaskRootCancelStopsTheModelCall (1.15s)`；`:219` 逐字「a cancelled task root still issued **1** provider requests, want 0」＋`:223`「a task cancelled before it started exited 0」。⚠ 同用例 `:199` 的正控「live root → execute code 0, provider chat requests 1」在这一发**是绿的**＝计数器没坏，红是从对照里长出来的 | porcelain 0 行；`sed -n '1046p'` 复现 `baseCtx := parent` | — |
| **MUT-4** | `cmd/wisp/resident_windows.go:177`：`src := startResidentTaskSource(rt, ra)` → `var src *residentTaskSource`（把"任务源没有真消费者"改回零调用者形状） | AC#7 票面那一枚"没有任务会去举它" | **四枚点名红，其中两枚是能力尺不是文案尺**：`TestLive246ResidentPipelineRaisesACardAndEscVetoesIt` → `resident_task_source_live_246_windows_test.go:120` 逐字「AC#7 LIVE RED (step 1 起任务 / step 2 举卡): the shipped resident process **never raised a card** from its own task source.」；`TestLive246ExitCancelsARunningTask` → `:255`「no streamed text ever arrived, so there was never a running task to cancel.」；再加默认层 `:307` 与 `:373`（「injection said=false pipeline=false loop=true」） | porcelain 0 行；`grep -n "src := startResidentTaskSource"` 复现 `:177` | **同发复跑四枚全绿**：`--- PASS` ×4，`ok … 42.605s`，rc=0 |
| **MUT-5** | `cmd/wisp/resident_windows.go:200`：地标句子 `resident event loop running` → `resident loop runs` | 攻 E（重锚的钉是否还咬） | 三枚点名红，逐字见 §3.E 第 2 条（`:498`／`:539`／`resident_ball_228_windows_test.go:83`） | porcelain 0 行；`grep -c` 复现 2 | **同发复跑三枚全绿**：`ok … 13.863s`，rc=0 |
| **MUT-6** | `cmd/wisp/resident_task_source_windows.go:185`：`case env != buildinfo.EnvTest:` → `case env != buildinfo.Env("MUT6-never"):`（test 分支永不成立） | 裁定 2.4（逃生口必须窄） | 两枚点名红：`TestAC246TestTaskInjectionPredicate` 子用例 `accepted_in_the_harness'_own_root` → `:245`「test env + harness-owned root must be accepted: text=""」＋子用例 `refused_when_the_root_is_not_the_harness'_own` → `:265`；`TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline` → `:373`「injection said=false pipeline=false loop=true」。⚠ `TestAC246DevLegIgnoresTheTestTaskInjection` 在这一发**仍绿**（它钉的是拒绝侧，MUT-6 把拒绝变成万能 ⇒ 它测不出，这是**该尺的射程**，不是缺陷） | porcelain 0 行；`grep -c "case env != buildinfo.EnvTest:"` 复现 1 | **同发复跑三枚全绿**：`ok … 7.730s`，rc=0 |

**六发红名的共同形状**（本腿要说的结论）：每一枚被中和的裁定都**至少有一枚能力尺**（指针身份／provider 请求计数／真机卡片与流式文本）响，而不是只有一枚文案尺响。⇒ **这一族的仪器不是装饰**；`246-r2` 死之前没来得及写的 §5（突变与正控）由本表这一节代位，且本腿的读数全部是自己取的。

## §5 我判不动的地方（具名交编排者裁，本腿不自行填、不自行修）

1. **AC#7 那句"任务被取消"与 D43 第 22 行不是一个射程——这是契约面，只能人拍。**
   转移表逐字（票面 `:12` 已引，本腿现读 `internal/agent/approval/approval.go:87-92` 那族通道名标签：`ChannelEsc: "按 Esc 键"`，票面写的 `:90` 落在 89——票面自己那句"别信行号"成立，⛔ 本腿未改该文件一字）：`Confirming`｜否决｜`Acting`＝**取消该调用并回给 LLM**。所以 Esc 的语义对象是"那一发调用"，不是"整条任务"。写腿两形都量了（§3.C 的第 ④ 行两栏），**没有**新造第五枚通道——本腿认同这是**唯一与冻结契约相容的读法**；但票面那一格要是想要后者，那就得动 C12／D43，那是人工批准，不是我能填的格。
2. **真控制台那一支要不要立"必须真终端"台件——本腿裁不了。**
   尺（现跑）：`grep -rn "runConsoleLoop" --include=*_test.go cmd/` ＝ **0 处**（产码里 3 处：定义 `:345`、协程体 `:319`、注释）。⇒ 裁定 2.1 指定的正路今天只有**反向**覆盖（MUT-1 红）＋写腿 §1.3 的手测。要不要为它造一枚伪控制台台件（人或 CI 附一个），是台件射程的裁，不是本票 AC 的射程；本腿不为此放宽 `interactiveStdin()`，也不主张"已覆盖"。
3. **`AskOnTaskRoot`／`askConfirmation` 今天是两枚零真消费者的产码接缝——留、删、还是改走它，本腿不裁。**
   本腿现读：`AskOnTaskRoot` 非测试调用者 **0**；`askConfirmation` 的唯一产码调用者就是 `AskOnTaskRoot`（`resident_approval_windows.go:264`）⇒ 传递不可达。今天举卡走的是桥→gate 那条既有路（§3.A 末段）。AC#7 不因此不成立（要防的结局没发生），但票面若把"接缝存在"读成"接上了"就是过期指认。`resident_approval_windows.go:252-256` 那段注释今天仍在说"它的调用者只有本包用例"——**这句仍然为真**，本腿特此具名，免得下一程反过来把它当缺陷清掉。
4. **两处被 AC#7 自己改假的自述句，本腿只登记不动。**
   - `cmd/wisp/resident_ball_windows.go:21`（注释）与 **`:160`（运行期真打印的 slog 字段）** 逐字「this leg has no task pipeline and no microphone」；`:203` 的 `ballGestureWhy` 常量逐字「this process has no task pipeline, no microphone and no panel host」。⇒ 前两条里 **`:160` 是打进台账的句子**，而 AC#7 之后这条腿**有**任务管线（真机六发的卡片就是从它举起来的）。打印时刻（`:136` 球起来）确实早于装配时刻（`:177`），所以那一行在它打印的那一瞬间不算谎，但它声明的是**属性**而非**此刻**，读台账的人会当成永久事实。票面「归口」把过期注释指给票 228 后续片与票 245，可**这一枚是本项目自己刚造成的** ⇒ 归谁，请裁。
   - `cmd/wisp/main.go:61` 注释仍逐字「park in the empty event loop until an exit signal arrives」，而 `3edc11d4` 已把同一文件 `:24-29` 的 usage 句子改成实际行为。写腿 §7.8 自报改了 usage、没说改了这行注释——本腿现读确认这行**没改**。
5. **票 244（GUI 子系统构建）一落地，裁定 2.1 指定的任务源会静默关掉——今天不炸，是本腿这轮挖出的最重一笔。**
   现读凭据三把：① `objdump -p build/wisp.exe` ＝ `Subsystem 00000003 (Windows CUI)` ⇒ **今天**双击进程有真控制台输入缓冲，正路是通的；② `cmd/wisp/console_windows.go:33-43` 的 `attachParentConsole` 只会 `AttachConsole(ATTACH_PARENT_PROCESS)`＋重绑 std handle，**Explorer 拉起的 GUI 进程没有父控制台**，且函数开头 `if out != 0 { return }`；③ `docs/specs/SPEC-11-build-deployment-containerization.md:50` 逐字要求"无参＝GUI（`-H=windowsgui`）"。⇒ 三线合起来的后果：`-H=windowsgui` 一旦落地，`interactiveStdin()` 在双击那条路上返回 nil，**常驻腿重新变成"有门、有球、没有东西举卡"**，而 AC#7 的仪器（真机那一族）用的是 `WISP_ENV=test` ＋注入位，**不会因此变红**。这是"票面声称要防的结局会在打包那天回来、而仪器不响"的形状。票 244／248／33 的射程都不含这一格 ⇒ 请在票 244 落地前定案（本腿不派单、不立票、不改 `docs/specs/**` 一字）。
6. **注入 gate 的形状带来的两处真实代价，要不要消，等于是不要改 AC#1 裁过的乙形次序。**
   本腿自己核过写腿说得对（现读 `approval/gate.go:135-145` 的 `win <= 0 → DefaultL1Window` 钳位、`approval/queue.go:107/115/122` ＝ 300s／DefaultL1Window 3s／MaxL1Window 3s）：常驻那枚 `approval.New` 在会话账本之前建 ⇒ `Options.Grants` 为 nil ⇒「本会话内允许」放行但不落盘并写 `GRANT-DROPPED`；且它吃不到 `[risk] l1_window_sec`／`confirm_timeout_sec`（窗口那一项即便接了也钳在 3s ⇒ 只有超时是真差异）。要消掉就得把 `approval.New` 移进 `assembleRuntime`＝改 AC#1 的裁定次序，本腿不动、也不认为写腿有权动（它没动，正确）。
7. **写腿注册了 D38(e 第 1 步与第 7 步）两枚钩子，AC#4 的射程只有第 3 步——算不算越界，请裁。**
   本腿现读确认十步顺序一字未动（`resident_task_source_windows.go:303-315` 只是往 `RegisterShutdownHook` 里挂两枚函数），且真机用例断言退出是 `10 steps, 0 failed` 并且 boot 报告里出现 `1:scheduler-close`／`7:flush-logs-close-db`（`resident_task_source_live_246_windows_test.go:188-193`，本腿 6 发都过）。本腿**倾向**这是 AC#7 的必要条件（没有第 1 步 Ctrl+C 后控制台还在接新任务；没有第 7 步存储会在还在写行的任务下面被关），但"倾向"不是裁——写腿 §7.5 自己问了，本腿复问一次。
8. **`task <文本>` 这枚动词是票面没有的**（写腿 §7.3）。本腿读过实现：裸行一律交 `replySurface.handle` 现成动词表，未知指令打印拒绝＋复读帮助行（`resident_task_source_windows.go:376-387`），**不会**把打错的答复变成一次真工具调用。本腿判这个选择方向安全，但它确实是**新增语法**（也是 `:333` 那枚 `consoleEntryHelpLine` 常量的全部内容），要否认可不认可请裁；本腿不自扩判据。
9. **常驻腿一次只接一发任务**（`:411-414` 第二行响亮拒绝）是设计选择，票面没写。本腿没有尺能判它对错，只确认它**打印**而不是静默丢弃。
10. **winlive 这一族在 CI 没有 `-tags winlive` 档** ⇒ 本腿这 6 发真机读数**只有本机有**。这是票面末段记的第三次出现，仍无人立尺。本腿不立票。
11. **`internal/ball` 那枚红的归属**：见 §7 第三把尺——它是工作树里 `design/assets/tokens.css` 被**别人未提交的删除**造成的，不是本票造成的，也不是 flake（本腿 `-count=2` 2/2 复现）。本腿既不恢复该文件（那是别人地界，且 `design/**` 对本腿是禁令层），也不替它找归口。

## §6 我对实现者表的更正（⛔ 不照抄它的自述；以下是本腿现跑与现读推翻或补齐的部分）

先说清楚**它表里哪些格子本腿用自取数据代位填了**：它的 §0（起手名册／主张名册）、§1.4（新尺首发与钉子复跑读数）、§1.5（三把尺终态）三处逐字写着「（填写中）」；它的 §2／§3.2／§3.4／§4／§5／§6 六枚小节号被正文引用，**文件里根本没有对应标题**（`grep -n "^## §\|^### §"` 现读只有 §0／§1（含 1.1-1.5）／§7 三节）。⇒ 本表的 §1／§3／§4／§7 就是那些格子的**代位凭据**，全部本腿自己取。

| # | 它的原句／原状态 | 本腿的更正 | 尺 |
|---|---|---|---|
| 1 | §1.1「新用例默认层首发 7 枚里 6 枚 PASS、1 枚 FAIL：`TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline`」 | **不适用于 HEAD**：那一发是在 `2071f59e` 之前的中间态量的（印序缺陷它随后自己修了）。本腿在起手 HEAD `aec06d24` 现跑：**15 枚 `TestAC246*` 全绿**（8 枚前腿的＋7 枚本票新增），0 红。⚠ 原句留在表里会被下一程读成"HEAD 有一枚常红"——那是假账，故更正 | `grep "^--- PASS: TestAC246" roster.log` ＝ 15 行；名册原文 `.scratch/wisp/probes/246/v2/roster-cmd-wisp-green.txt` |
| 2 | §7.1「run1 那枚红的名字是本腿丢的……本腿没能闭合这条归因」 | **本腿也没能找回那枚名字，但可以把这桩悬案关掉**：本腿起手整包 **0 红**（`ok … 173.655s`）⇒ HEAD 上没有待归因的红；唯一已知候选仍是票 228 那枚 boot flake（`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`），归口不变（§2 末段）。丢失的名字本腿**不**替它编一个 | `grep -c "^--- FAIL" /tmp/246v2/baseline.log` ＝ 0；`-v` 名册 141 枚顶层 PASS／0 枚 FAIL／0 枚 SKIP |
| 3 | 它表头与 §7.9 说"真凭据那一发＝欠账，记在 §2" | **§2 不存在**（见上）。⇒ `cmd/wisp/resident_task_source_live_246_windows_test.go:26-28` 那句"booked as owed in `docs/evidence/s1/246-resident-task-source-r2.md §2`"是**指向不存在小节的断链凭据**。本腿把它在本表 §3.C 第二栏补成实存在；**那两枚文件的字本腿一枚未改**（产码禁区） | `grep -n "^## §" r2.md` ＝ 只有 §0/§1/§7；`grep -n "§2" 246-resident-task-source-r2.md` 命中 2 处引用、0 处标题 |
| 4 | 提交 `8af4347d` 的标题「四步真机读数拿到」 | **措辞缺陷，不是产码缺陷**：不带栏头的"四步真机读数拿到"离"任务管线已跑通"只差一次转述。本腿在 §3.C 用两栏结构重写；票面翻勾时请照栏头写，不要照提交标题写 | 该提交 numstat 只含测试文件一枚，正文引本腿 §3.C |
| 5 | §1.2 第 2 条「`resident_approval_246_windows_test.go:271/324/367/374` 那两句与那串 forbidden 词本腿一字未改」 | **复核为真**：`git show --numstat` 五枚提交里**没有一枚**列出该文件；且 `TestAC246StatusLineSaysWhatTheLegDoesNot` 与 `TestAC246ShippedResidentProcessOwnsItsCancelStep` 在本腿名册里现跑绿（3.20s／0.00s）。⇒ "没有任务源那一层诚实"没有被搬薄 | 五枚 numstat 逐枚现读；roster 同名册行 |
| 6 | §1.2 第 3 条「`runSpec{` 六处测试构造点全是加字段，零值语义不变」 | **复核为真**，并由 MUT-2 加一条它没测的：注入分支被中和时，`wisp run` 那五枚用例**不受影响**（本腿 MUT-2 那一发 `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline` 仍绿＝默认路径确实与注入路径解耦） | §4 MUT-2 行 |
| 7 | §5／§6 整节缺失（它死于 150 轮上限，四枚 `MUT-` 留在工作树） | **本腿另跑六枚（MUT-1…MUT-6），四枚对应派单点名的四条裁定，另两枚是 E 的反向证与谓词窄射程**；全部逐枚还原、零枚入仓；`git status --porcelain -- cmd internal` 收尾 0 行 | 本表 §4 全表 ＋ §7 |
| 8 | 它 §7.8 自报"改了 `main.go` usage 那段文案" | **属实且必要**（AC#7 之后"takes no task yet"那句是假话）；但它**没提**同一文件 `:61` 那行注释仍是"empty event loop" ⇒ 见 §5 第 4 条第二枚 | 本腿现读 `main.go:24-29` 与 `:61` |

## §7 收尾三把尺

## §8 next
