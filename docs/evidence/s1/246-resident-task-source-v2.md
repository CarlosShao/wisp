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

## §4 突变台账（改了哪一行 → 指名哪一枚用例红 → 还原凭据 → 复跑回到绿）

（本节随各发突变推进补满）

## §5 我判不动的地方

## §6 我对实现者表的更正

## §7 收尾三把尺

## §8 next
