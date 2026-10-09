# 01 A 道真机读数 — 常驻进程真举出了一张卡（票 277 AC#1／AC#2／AC#3）

腿：`277-r1`。锚＝HEAD `bf9974c2`（起手时）＋本腿第 1 笔 commit `9d7bae83`。时刻＝2026-10-09 09:49–09:54 +08:00。
命令原文＝直读 `cmd/wisp/resident_task_source_live_246_windows_test.go:38-41`（配方 §2 甲、编排者转述三方同文）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp \
  -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v
```

⛔ 产码零改动、⛔ 台件零改动、⛔ `build/wisp.exe` 未重建未覆盖（见文末复核）、⛔ 零 push、⛔ 未翻任何框。

---

## 一句话结论

**绿。** 常驻进程（真 exe、真悬浮球窗口、真审批门）从它自己的任务源起了一发任务，
**真举出一张 L1 卡片**，同发的台账卡记录带 `esc_borrowed=true`／`orb_state=Confirming`，
注入的裸 Esc 在窗口期内把它否决、目标文件没落盘，两枚用例 `--- PASS`、`RC=0`，
配方 §4 的红句与「判不了」句**一发都没出现**。

---

## 三发执行记录（⛔ 不挑对你有利的：全部留档，含一发是假读数）

| 发 | 命令差异 | rc | 是不是真跑 | 读数 |
|---|---|---|---|---|
| run1 | 配方原文，无 `-count=1` | 0 | **真跑**（首次执行，32.748s） | 卡行 corr `195a2c6f-152e-4316-b75a-829d6d207194#call_246r2`；两枚 `--- PASS` |
| run2 | 同上重打 | 0 | **不是真跑＝go test 缓存回放**（见下） | 与 run1 **逐字节同**：同 `10.95s`/`21.60s`、同 corr `195a2c6f…`、同临时目录名 `…3405263550\002`，且 `ok … (cached)` |
| run3 | 原文＋`-count=1`（强制真跑；未动任何断言/等待常量） | 0 | **真跑**（32.968s） | 卡行 corr `0455ee9f-5fc7-44b4-9d28-e3384163622c#call_246r2`；两枚 `--- PASS`；**台账两栏在这发取到手**（见 §AC#1(b)(c)） |

具名一笔自纠：run2 我原本打算当"第二发独立读数"，它其实是缓存回放，**没有新进程、新桌面动作**；
本件的"两发真跑"＝ run1 ＋ run3（编号每发不同，与配方 §5 记的历史形状一致）。
run3 另加 `-count=1` 是**重跑旋钮**，不是断言/等待常量的改动（配方 §5 自己就记着历史读数用 `-count=1`／`-count=2`）。

**为何 run1 缺 (b)(c)、run3 才齐**：那两栏是产码 `slog` 行，落在**数据根 sink 文件**里，
路径由用例自报（`resident_task_source_live_246_windows_test.go:181`），
而数据根＝`t.TempDir()` ⇒ 用例收尾即删。run1 跑完后我按路径去找，目录已不存在
（`ls -d /c/Users/swq/AppData/Local/Temp/TestLive246*` → No such file）。
run3 的做法是**同一条命令重跑**、同时在进程存活期把 sink 文件只读复制到仓外 `/tmp/277sinks3/`
（⛔ 未改台件、⛔ 未改产码、复制件不入库，仓内只贴关键行逐字）。

---

## AC#1 逐判据（判据原文＝配方 §4 三栏，未换写法）

### 绿栏 (a) stdout 卡行

模板（本腿从 `git show`/Read 现读产码）：`cmd/wisp/resident_approval_windows.go:886`
`fmt.Printf("wisp: 卡片挂起：%s %s（编号 %s）\n", p.Level, p.Tool, p.CorrelationID)`。

run3 逐字（用例 `:145` 的 `t.Logf` 转报它读到的子进程 stdout 行）：

```
resident_task_source_live_246_windows_test.go:145: AC#7 LIVE steps 1-2: wisp: 卡片挂起：L1 fs.write（编号 0455ee9f-5fc7-44b4-9d28-e3384163622c#call_246r2）
```

run1 同句、编号 `195a2c6f-152e-4316-b75a-829d6d207194#call_246r2`。⇒ (a) 中。

### 绿栏 (b) 台账卡记录（结构化字段）

逐字（run3 sink `…\TestLive246ResidentPipeline…4039553199\002\logs\wisp-20261009-001.jsonl`）：

```
{"time":"2026-10-09T09:53:55.1616661+08:00","level":"INFO","msg":"approval: 常驻进程显示一张确认卡片","corr":"0455ee9f-5fc7-44b4-9d28-e3384163622c#call_246r2","level":"L1","tool":"fs.write","orb_state":"Confirming","esc_borrowed":true,"channels":"ball=unloaded | esc=loaded | panel=unloaded | kws=unloaded"}
```

`orb_state=Confirming`（L1 应有值）、`esc_borrowed=true`、`channels` 只有 `esc=loaded`
⇒ 与配方 §1-P5「本票只落 Esc 一条通道」逐字对得上。⇒ (b) 中。

### 绿栏 (c) 同一 corr 与同发 `风险判定 … level=L1` 同链

同文件逐字两行（时间上在卡行**之前 603µs**，同一进程同一流）：

```
{"time":"2026-10-09T09:53:55.1610658+08:00","level":"INFO","msg":"audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""}
```

诚实边界（⛔ 不许把链说满）：`风险判定` 这一行**本身不带 `corr` 字段**（它带的是 tool/level/rules/reason）。
"同链"是靠同一 record 流＋同一 tool/level＋紧邻时序，再由**带 corr 的那行**把三者钉死：

```
{"time":"2026-10-09T09:53:55.1958445+08:00","level":"INFO","msg":"audit: tools: call kind=refused task=0455ee9f-5fc7-44b4-9d28-e3384163622c corr=0455ee9f-5fc7-44b4-9d28-e3384163622c#call_246r2 tool=fs.write risk=L1 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）\""}
```

⇒ (c) 中（判定链 L1／R1／`grant_id=0` 与卡行的 corr 同一家族：task=`0455ee9f…`、corr=`…#call_246r2`）。

### 绿栏 (d) 用例 `--- PASS` 且无红句

```
--- PASS: TestLive246ResidentPipelineRaisesACardAndEscVetoesIt (11.57s)
--- PASS: TestLive246ExitCancelsARunningTask (21.32s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	32.968s
```

### 红栏／判不了栏：逐句扫过，**零命中**

对 run3 全量 grep（`grep -nE`）以下逐字句：`never raised a card`／`任务入口未启用`／`任务管线未装配`／
`在当前环境不被受理`／`卡片无处呈现`／`somebody else owns the key on this desktop, so the borrow cannot be measured here`／
`the observer rig is blind`／`取消键未借到`／`--- FAIL`／`--- SKIP` ⇒ **全部 0 命中**（run1 同样 0 命中）。

### 四步全走的顺带读数（票 246 `AC#7` 判据那四步，本票只作顺带记录、不代它下结论）

- 第 3 步否决（run3 sink 逐字）：
  `audit: approval: ANSWER-VETO corr=0455ee9f-…#call_246r2 tool=fs.write channel=esc decision=veto`
  ＋ `按 Esc 否决了卡片 0455ee9f-…#call_246r2（L1 / fs.write）：该调用未执行，答案已入审计`
- 观察台件自证（run1/run3 stdout 逐字）：`AC#7 LIVE: observer keydown 1->0->1 (steal control 0)`
  ⇒ Esc 键**先归别家（1）、卡片窗口期被借走（0）、收口归还（1）**，借到没被抢。
- 第 4 步：`wisp run: 任务 0455ee9f-5fc7-44b4-9d28-e3384163622c 结束（completed，2 轮，1 次工具调用，成本 0 CNY（未计价））`
  （run1 逐字同形，编号不同）
- 目标文件不落盘＝用例内 `os.Stat` 断言（`:162`），它没报错 ⇒ PASS 已含这一条，本腿不另报读数值。
- 第二枚用例 `TestLive246ExitCancelsARunningTask`：`AC#7 LIVE (exit cancels a running task): wisp run: 任务已取消`。

### 配方 §1 前置在本发的真机自证（同一 sink，run3 逐字）

```
{"time":"2026-10-09T09:53:54.736943+08:00",…,"msg":"winsec: sealing path resolver installed","resolver":"risk.c26Pipeline","probes_passed":2}
{"time":"2026-10-09T09:53:54.7…",…,"msg":"ball: the resident leg created the floating ball window"}
{"time":"…","level":"INFO","msg":"audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"}
```

- **P4 球窗口在场**＝第一手（上一行逐字）。
- **P5 通道只有 Esc**＝上一行后半句逐字（与卡行 `channels` 里的 `esc=loaded` 对上）。
- **P3 任务入口＝注入形**：那一句（`wisp: 任务入口经测试注入位打开…`，配方 §1-P3）是**子进程 stdout**、不是 slog 行，
  所以**不在**我复制到的 sink 里，用例也只把**卡行**tail 进 `t.Logf` ⇒ **本腿没有它的逐字句**。
  间接凭据＝用例 `:142-144` 对 `taskEntryInjectedClaim` 是硬 `t.Errorf` 断言（"任务来自注入但 boot 没说过"就红），该用例 PASS。
  ⛔ 我不把这一条写成"逐字到手"。

### 两栏分写（配方 §4「灰」：不许并成一格）

- **〔接缝注入〕栏**：以上全部读数都是这一栏。任务文本走 `WISP_TEST_TASK_TEXT`（`resident_task_source_boot` 侧受理，
  配方 §1-P3 注入形），"模型"是 golden SSE 脚本 `residentCardGolden`，**不是真凭据真模型**。
  注入形算不算票面那句"测试构造的任务源冒充"＝ `A732` 已登记**待人派**，本腿不自裁。
- **〔真凭据／真模型〕栏**：**空**。这一格 A 道给不了，见 AC#2。

---

## AC#2 B 道前置记账（本腿现跑的尺，⛔ 不是转述配方）

**尺 1 真凭据／config.toml 在不在**（本腿 09:52 现跑）：

```
$ ls -la "$APPDATA/wisp"
drwxr-xr-x 1 swq 197609 0 Oct  7 09:33 .
drwxr-xr-x 1 swq 197609 0 Oct  7 19:02 ..
drwxr-xr-x 1 swq 197609 0 Oct  7 09:33 logs

$ ls -la "$APPDATA/wisp/secrets"
ls: cannot access 'C:\Users\swq\AppData\Roaming/wisp/secrets': No such file or directory
```

读数＝数据根只有 `logs`，**没有 `config.toml`、没有 `secrets` 子目录** ⇒ 配方 §1-P2 的 B 道缺料在本机**今天仍然成立**；
B 道（真控制台敲 `task`）今天**判不了**（无凭据 ⇒ 会降级成 `配置未就绪（Unconfigured）`／`任务管线未装配`，配方 §1-P2 引的 r2 读数形状）。

**尺 2 `build/wisp.exe` 旧于 HEAD**（`stat` 现读＋`git log` 现读）：

```
$ stat -c '%y %s %n' build/wisp.exe
2026-10-07 11:57:46.007254300 +0800 31074357 build/wisp.exe

$ git log -1 --format='%h %cI' -- cmd/wisp
aa6c1881 2026-10-08T19:20:48+08:00
```

读数＝exe 比**最近一次碰 `cmd/wisp` 的提交**旧 **31 小时 23 分**，更旧于 HEAD（`bf9974c2` 2026-10-09 09:44:44 +0800）
⇒ 配方 §5 那一栏本机复核成立。⚠ 本腿**没有**重建它（机主未授权；A 道吃的是 `buildWispForTest` 在临时目录现建的 exe）。
B 道真手测那一发要 owner 先重出 exe（`GOFLAGS= go build -o build/wisp.exe ./cmd/wisp`，姿势在 `docs/evidence/s1/246-resident-task-source-v2.md:21`）＋自录凭据。

**尺 3「人坐终端敲 `task`」这一形有没有自动台件**：

```
$ grep -rn "runConsoleLoop" --include=*_test.go     （本腿＝Grep 工具，glob *_test.go，全仓）
No matches found
```

读数＝**测试调用＝0**，与配方 §2 乙一致 ⇒ 控制台键盘那一支今天没造台件（`consoleApprovalUI`／`runConsoleLoop` 只有产码路径）。

**尺 4 A 道四前提（配方 §5 缺谁就打不了）**：本腿 09:49–09:50 现尺，逐条在 `00-anchor.md`：
go＋gcc 在 PATH、`third_party/sherpa-onnx` 三枚 dll 在盘、真桌面＋`balldebug.exe`／`wisp.exe` 各 0 枚、Esc 空闲（由读数自证：`keydown 1->0->1`）。

---

## AC#3 与票 246 的边界

- 本腿**没动**票 246 任何框、**没写**票 246 的结论。票面 `AC#7` 至今是 `- [x]`（本腿只读，`grep -n "AC#7"` 于 09:52）。
- 零触碰凭据（本腿现跑，输出为空＝未改）：

```
$ git status --porcelain -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go \
    internal/perm/ticket90_persist_test.go .scratch/wisp/issues/246* .scratch/wisp/issues/277* \
    docs/reports/pending-and-issues.md
（空）
```

- 结论只写在本件（＋`00-anchor.md`）。翻框、台账、票面回写一律归编排者。
- 本件能给票 246 的那半句话，**以读数形状说、不以结论说**：`AC#7` 判据的四步在〔接缝注入〕栏今天**有一发真机走通的原始记录**，
  它覆盖不覆盖 `AC#7` 那格＝编排者裁（`A732` 已把"真举卡"另立为本票）。

---

## 本腿具名：没做的／判不动的／顶回来的转述

**没做的**
1. B 道（真控制台 `task`）一发未跑——缺凭据＋缺重出的 exe，条件不齐⛔ 不派（票面 AC#2 原文）。
2. L2 卡一发未跑——配方 §2 丙具名"没造"，盘上没有把任务驱动到 L2 的场景。
3. sink 快照文件（两枚 jsonl）留在仓外 `/tmp/277sinks3/`，**没有**入库（写面只允许 `.md`；仓里也不许新建 `.txt`／`.out`）。

**判不动的**
1. "面板在场时卡真的画在界面上"——本腿无面板（`resident_task_source_windows.go:273-275`「leg has a ball but no panel」），
   观测面＝stdout＋sink 台账＋用例 `t.Logf`；票面 AC#1 明写**不要求**面板在场。
2. run1 的 (b)(c) 两栏**永久取不到**（临时数据根已删）；只有 run3 有。本件按发分栏写明，不混算。

**顶回来的转述／盘上不符（以盘上原文为准）**
1. 编排者与票 277 都写卡行＝`resident_approval_windows.go:885`、审计行＝`:880`。
   **本腿 Read 现读：卡行在 `:886`，`slog.Info` 那五条 attr 在 `:881-885`**（`:885` 现在是 `"channels", channelRosterText(...)` 那一行）；
   无球 fail-closed 句在 `:868`（票面/配方写 `:864-869`/`:867`）。行号整体**下移 1**，引句一字未变。
2. 编排者写「写面＝只新建 `probes/277/r1/*.md`」。**仓里没有顶层 `probes/`**；既有各腿一律是
   `.scratch/wisp/probes/<腿名>/…`（`probes/246/r2`、`probes/card-proof-prep-1`、编号目录 111/114/…）。
   本腿按仓内约定落 `.scratch/wisp/probes/277/r1/`，已在 `00-anchor.md` 具名。
3. 编排者硬坑第 2 条（长跑输出先落文件）在本仓**还多一坑**：`go test` 的**结果缓存**会把"重跑"变成**逐字节回放**
   （run2＝`ok … (cached)`，同 corr、同时序、同临时目录名）。**不带 `-count=1` 的第二发不是第二次测量**。
   配方 §5 记着历史读数用了 `-count=1`，但没把这条写成坑。
4. `风险判定` 审计行**不带 corr 字段**（配方 §4 (c) 那句"该 corr 与同发的 `风险判定` 审计行同链"在字段层面要靠
   `tools: call kind=refused` 那行才闭合）——本件按实际形状写，没有假装它自带 corr。
