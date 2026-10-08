# 01 配方：真机读数「一枚任务真跑到需批准的那一步 ⇒ 门调 `ui.Prompt` ⇒ 界面上出现一张卡」

- 腿：`card-proof-prep-1`（只读准备腿，**本腿一发未跑**）。锚＝`f13e7c372cd5ed00a1525bcf9d38450ec720cfac`（详见同目录 `00-anchor.md`）。
- 用途：票 246 `AC#7` 的另一半读数（台账 `A732` 裁定：`AC#7` 的勾只覆盖"起一条任务管线"那半；"常驻进程会不会真给你举出一张审批卡"要**行为凭据**，今天零凭据、待人派）。本配方把那一发写成**照着敲就能做、做完能判**。
- 行号纪律：下面每一处「文件:行号」都是本腿从 `git show HEAD:<path>` 现取（2026-10-08 18:5x–19:0x），引句与行号同一发现取。⛔ 本件不动任何契约、票面、产码、台账。
- 判据原文（经 `probes/246/raisercensus2/01-census.md` §4 转引票面 `:60`，本腿未直读票面）：「完成判据＝常驻那条腿真起一条任务管线（或经装配根注入一条最小任务源），一发真机走通"起任务 → 举卡 → Esc 否决 → 任务被取消"四步；⛔ 不许用测试构造的任务源冒充。」

---

## 1. 前置条件（每条给判据：怎么知道它在不在）

**P1 判定级必须走到读点。** 卡只在工具侧读点被问到才出现：`Bridge.Execute`（`internal/tools/bridge.go:268`）→ `b.route(ctx, &dec, sil, grantID)`（`:352`）→ `switch sil.Level`（`:424`）里 **L1 且 `grantID==0`** 走 `b.gate.PendingWindow(ctx, *dec)`（`:443`）、**L2** 走 `b.gate.PendingApproval(ctx, *dec)`（`:458`）。
判据：跑起来后进程自己的审计行 `audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason="R1: 工具声明为下界（L1）"`（实跑样本逐字在 `.scratch/wisp/probes/246/r2/live-four-step-run1.txt`）。⚠ L1 若已有 D45 会话授权（本进程里有人答过 `session`），`:434-442` 直接 `DecisionAllowGrant` 短路 ⇒ **不举卡**；判据＝用**新起进程**。L2 没有现成台件造（见 §2 末具名）。

**P2 数据根有 config.toml（且凭据可解析）。** 缺了会降级为"任务管线未装配"（不是崩溃）：读数逐字 `wisp run: 配置未就绪（Unconfigured）：…`＋`wisp: 任务管线未装配（退出码 2，…）`（`.scratch/wisp/probes/246/r2/manual-boot-injected.txt`，本腿现读）。
判据：启动输出出现 `wisp: 审批门已装配进本进程（取消通道：Esc 已加载；等待中的确认项：0）; 任务来源：…`（模板 `cmd/wisp/resident_approval_windows.go:802`）。A 道（§2 甲）由台件自造：`prepareResidentHarness`（`cmd/wisp/resident_task_source_live_246_windows_test.go:324`）写 config.toml ＋ `secret.NewStore(...).Store("dpapi:acme", fakeStoreKey)` 的 DPAPI blob（密钥只以常量名出现，不打印）。B 道（§2 乙）要 owner 自己录真凭据——今天本机 `%APPDATA%\wisp` 没有 config.toml／没有 secrets 子目录（`docs/evidence/s1/246-resident-task-source-v2.md` §3.C 第二栏现读），B 道今天**判不了**。

**P3 任务入口二者至少居一（不然绝不举卡）。**
- 控制台形：`interactiveStdin()`（`cmd/wisp/approval_reply_stdin_windows.go:41`，判据＝`GetConsoleMode` 成功 ⇒ 真控制台输入缓冲）≠ nil；启动印 `wisp: 任务入口已启用（控制台键盘）：输入 task <文本> 起一发任务；…`（`cmd/wisp/resident_task_source_windows.go:120`＋`:253-254`＋`:362`）。
- 注入形（测试窄门）：`WISP_ENV=test` ＋ `WISP_TEST_DATA_DIR=<数据根>` ＋ `WISP_TEST_TASK_TEXT=<文本>` 三条全中才受理（判据源码 `resident_task_source_windows.go:192-206`；变量名 `:137`）；启动自报 `wisp: 任务入口经测试注入位打开（WISP_TEST_TASK_TEXT）：…`（`:257-258`）。设了但环境不是 test 会被响亮拒：`wisp: WISP_TEST_TASK_TEXT="…" 已设置，但本进程的环境是 dev：这条注入位在当前环境不被受理（只有 test 受理任务文本注入）`。
- 两条都没有 ⇒ `wisp: 任务入口未启用（本机没有可交互控制台）：本进程仍然带球常驻…但没有任何东西会去举一张卡`（`:101`＋`:237-243`），boot 报告 `任务来源：无（任务入口未启用，理由见上面那行）`（`:122`）。

**P4 球窗口在场。** `ballCardUI.Prompt` 无球即 fail-closed（`cmd/wisp/resident_approval_windows.go:864-869`，`errResidentNoBall` `:88`，slog「卡片无处呈现…」`:867`）。判据＝启动行 `the floating ball window is up in this process (tray icon added, hotkeys live 3/4: …)`。

**P5 否决/允许入口今天有几条（照实说）。** 否决＝**Esc 一条**，只在 L1 窗口期借、收口即还（装 `:567` `Channels().SetLoaded(approval.ChannelEsc, true)`、卸 `:787`）；装配时的审计句逐字「本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入」（`:574-575`）。球/托盘**没有允许入口**（托盘项＝open-panel / mute / pause-wake / exit，`cmd/wisp/resident_ball_windows.go:284-286`；exit 只有记录无停机路径 `:452-462`）；**允许**在 B 道走控制台答复动词 `yes / session / no / veto / always / head / view`（帮助行逐字 `resident_task_source_windows.go:362`）。

**P6 桌面独占 ＋ Esc 键空闲。** A 道用例头部逐字「Run it with no other Wisp or balldebug process on the desktop」（`:39`）；跑前判据＝`tasklist //FI "IMAGENAME eq balldebug.exe"`（`wisp.exe` 同）皆 0（姿势抄 `docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md` §起手）。Esc 被别家占着 ⇒ **判不了**（用例自己的正控红句，见 §4）。

---

## 2. 怎么把一发任务喂到"需要批准"那一步（只列盘上已存在的入口）

**甲（推荐，确定性，今天就能跑）：winlive 台件。**
命令原文（文件头逐字，`cmd/wisp/resident_task_source_live_246_windows_test.go:40-41`）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp \
  -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v
```

它内部（盘上已有，逐跳到行）：
① `buildWispForTest`（`cmd/wisp/secret_argv_windows_test.go:161`）现建一枚真 `wisp.exe` 到临时目录（`:174`）并拷 `third_party/sherpa-onnx/*.dll`（`:184`）——**不吃 `build/wisp.exe` 的旧**；
② 以 env `WISP_ENV=test`、`WISP_TEST_DATA_DIR=<temp>`、`WISP_TEST_TASK_TEXT=把这段说明写进 note.txt`（`:127`）拉起无参 exe ⇒ 落到 `main.go:66 runResident()` → `resident_windows.go:260 src := startResidentTaskSource(rt, ra)` → 任务源把 `ra.gate/ra.ui/ra.cards/ra.root` 经 `assembleRuntime`（`resident_task_source_windows.go:278`；spec 字面量 `:265-277`，`gate: ra.gate` `:269`）装配成环路；
③ golden SSE 脚本（`residentCardGolden` `:372`）让"模型"发**一次 `fs.write`** 到 `<dataDir>/resident-246r2-card.txt`（在 `[fs] allowed_dirs` 内 ⇒ R1 判 **L1**）⇒ `bridge.go:352→:443` → `gate.go:294 g.ui.Prompt` → `:864 ballCardUI.Prompt` **举卡**；
④ 第二个进程注入裸 Esc（`buildEscListener246` `cmd/wisp/resident_approval_live_246_windows_test.go:405`）⇒ 否决 → `按 Esc 否决了卡片…`，文件不落盘。

**乙（正路但今天缺料）：真控制台键盘。** 从真终端启动 `build/wisp.exe`（无参），等 `任务入口已启用（控制台键盘）`，敲 `task <文本>`（动词 `resident_task_source_windows.go:387`，整段余文＝任务文本）。它会落到同一条装配链（`:278`→`run.go:606 rt.gate = s.gate`→`:748 Gate: rt.gate`→`run.go:997 Tools: rt.bridge`→`bridge.go:352` 读点）。⚠ 缺料：真模型凭据＋"模型是否真发需批准调用"不在我们手里 ⇒ 非确定性；且 `runConsoleLoop` 的自动覆盖＝**0**（`docs/evidence/s1/246-resident-task-source-v2.md` §2 条件②现读：`grep -rn "runConsoleLoop" --include=*_test.go`＝0）＝**"人坐终端敲 task"这一形今天没造台件**。

**丙（没造，具名）：L2 卡。** 读点 `bridge.go:458` 与 `gate.go:536` 在码上，但盘上没有现成场景把一次任务驱动到 L2（现成 golden 只造 L1）；要 L2 得另造，本腿未见 ⇒ 写"没造"。
（`wisp run` 不是这一发的目标：那是 CLI 腿，自造门与 `consoleApprovalUI`，`cmd/wisp/run.go:611-612`。）

---

## 3. 观测点：卡出现的可观测信号（这一腿**没有页面**，"由界面/页面自己报回"落在进程自己的 stdout/台账）
具名：本腿没有面板（`resident_task_source_windows.go:273-275` 逐字「leg has a ball but no panel」）；界面侧只有悬浮球的状态（L1＝D43 的 `Confirming`）。既有台件谁报回＝winlive 用例自己的 `t.Logf`（`AC#7 LIVE steps 1-2: …`，`:145`）＋它读的进程 stdout 与数据根 `logs/`（sink）。

期望输出逐字（模板来自产码；样本行来自 r2 实跑读数 `.scratch/wisp/probes/246/r2/live-four-step-run1.txt` 与 `docs/evidence/s1/246-resident-task-source-v2.md` §3.C）：

1. **卡挂起（主信号，stdout）**：`wisp: 卡片挂起：L1 fs.write（编号 <corr>）`（`resident_approval_windows.go:885`；用例等待常量 `cardRaisedClaim = "wisp: 卡片挂起："`，`:66`）。r2 样本编号 `cf73a8eb-67f4-4229-8c21-e685f35130b5`。
2. **台账卡记录（slog）**：`approval: 常驻进程显示一张确认卡片`，attr＝`corr/level/tool/orb_state/esc_borrowed/channels`（`:880-884`）；期望 `orb_state=Confirming`（L1；L2 为 `AwaitingApproval`）、`esc_borrowed=true`（L1）。
3. **入口与任务来源自报**：§1-P3 那两句其一 ＋ boot 报告 `任务来源：控制台键盘（已接入）` 或 `仅本发的测试注入位`（`:120-121`，字段函数 `:528-534`）。
4. **否决（四步全走时）**：`wisp: 按 Esc 否决了卡片 <corr>（L1 / fs.write）：该调用未执行，答案已入审计`（`resident_approval_windows.go:592-593`；用例等待 `vetoDoneClaim = "按 " + ball.DefaultHotkeys().Cancel + " 否决了卡片"`，`:87`）＋台账 `audit: approval: ANSWER-VETO corr=… tool=fs.write channel=esc decision=veto`（r2 原文 `:170` 样本）；＋**目标文件不存在**（`os.Stat` IsNotExist）。
5. **反面（不是"举了"）**：无球 ⇒ `approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝`（`:867`）＋门侧 `确认界面不可达（…），已 fail-closed 拒绝执行`（`internal/agent/approval/gate.go:297`）；Esc 没借到 ⇒ `wisp: 卡片 <corr> 的取消键未借到（桌面已有占位者），按 Esc 不会否决它`（`:921`）。
6. 用例侧对照句：`AC#7 LIVE steps 1-2: …卡片挂起…`（`:145`）、`AC#7 LIVE: observer keydown 1->0->1 (steal control 0)`（r2 原文）。

---

## 4. 判红绿

**绿（"真举出来了"）＝同一发全中**：(a) stdout 有 `wisp: 卡片挂起：L1 <tool>（编号 <corr>）`；(b) 台账有 `approval: 常驻进程显示一张确认卡片` 且结构化字段带 `esc_borrowed=true`／`orb_state=Confirming`（L1）；(c) 该 corr 与同发的 `风险判定 … level=L1` 审计行同链；(d) A 道用例 `--- PASS` 且无任何 `never raised a card` 红句。
**红（"没举"）＝任一条**：出现 `任务入口未启用（本机没有可交互控制台）`／`任务管线未装配`／注入句 `在当前环境不被受理`；或轮询超时后 stdout 始终无 `卡片挂起：`——A 道红句逐字 `AC#7 LIVE RED (step 1 起任务 / step 2 举卡): the shipped resident process never raised a card from its own task source.`（用例 `:120`，r2/v2 原文）；或举卡路径落在 `卡片无处呈现`。
**判不了（⛔ 不许当绿也不许当红）**：(a) Esc 被别家占着——用例正控红句逐字含 `somebody else owns the key on this desktop, so the borrow cannot be measured here`（`AC#7 LIVE RED (desktop state, not our code)`，r2 原文 `:88-91` 一带）；(b) 没有真桌面/窗口起不来/观察台件自检 `the observer rig is blind`；(c) B 道无凭据（§1-P2）——只能给〔接缝注入〕栏读数，⛔ 不得写成"真模型跑通"。
**灰（照实登记，不许自裁）**：注入形算不算票面那句"测试构造的任务源冒充"——`A732` 已登记待人派；配方只要求**两栏分写**（接缝注入栏／真凭据栏），不得并成一格。

---

## 5. 诚实边界（这一发在本机今天跑得动/跑不动）

- **A 道跑得动的前提**（缺谁就打不了）：① go 工具链＋C 编译器在 PATH（台件自建 exe：`secret_argv_windows_test.go:161/:174`，失败句 `is mingw on PATH?` `:179`）；② `third_party/sherpa-onnx/*.dll` 在盘（`:184` 断言；本腿现读该目录有 `onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）；③ 真桌面＋无别的 wisp/balldebug 进程；④ Esc 键空闲（§1-P6）。
- **`build/wisp.exe` 约束**：现盘 mtime `Oct 7 11:57`（本腿 `ls -l` 现读），而最近 `cmd/wisp` 提交＝`65f4c968 2026-10-08 15:48` ⇒ **比 HEAD 旧**；A 道不受它牵制（自建 exe），但 **B 道手测那一发必须先重建**（`GOFLAGS= go build -o build/wisp.exe ./cmd/wisp`，姿势抄 `docs/evidence/s1/246-resident-task-source-v2.md:21`）。⛔ 本腿零 Go 命令，未复现该构建。
- **不带 PATH 的坑**：`go test ./cmd/wisp` 不带 `third_party/sherpa-onnx:build` 会以 `0xc0000135` 死在加载、且**没有一行 `--- FAIL`**（r2 表 §1.1 末段实测；本腿未复跑）⇒ 命令必带 PATH 前缀。
- **winlive 族在册、但 CI 不跑**：`cmd/wisp/resident_task_source_live_246_windows_test.go:1`＝`//go:build windows && winlive`；CI 无 `-tags winlive`（`docs/evidence/s1/246-resident-task-source-v2.md` §6 第 10 条）⇒ 读数**只有本机有**。历史读数：`v2` §3.C 记 6 发（`-count=1`＋`-count=2`×两枚用例）零 skip 零 FAIL、编号每发不同——**本腿未复跑，仅转述**；票面"真凭据/真模型"栏仍为空（欠账）。
- **没造的**：真控制台键盘那一支的自动台件（`runConsoleLoop` 测试调用＝0）；L2 场景台件；"允许"在球/托盘上的入口（结构性没有，见 §1-P5）。
