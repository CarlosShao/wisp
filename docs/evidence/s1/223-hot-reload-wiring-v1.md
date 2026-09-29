# 票 223 非实现者对抗验收 `223-v1` —— CheckAndReload 热加载接线

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 被裁的实现件：`docs/evidence/s1/223-hot-reload-wiring-r1.md`（实现腿 `223-r1` 死于 150 轮上限；AC#5 由编排者代跑；五枚文件未提交改动由编排者代落＝commit `248095d1`）
- 落点普查前案：`.scratch/wisp/probes/223/c1/census.md`
- 跨票前案：`.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`／`docs/evidence/s1/226-config-write-no-clobber-v2.md`／`internal/config/writeguard.go`
- 本腿起手锚点：`git rev-parse --short HEAD` = **`577ae8c7`**，分支 `dev`，`date` = 见下方起手时刻表
- 台件目录：`.scratch/wisp/probes/223/v1/`（原始输出全量落盘，绝不接 `| head`／`| tail`）
- 本腿性质：**只裁、只跑台件、不产码、不 commit 任何源码改动**
- ⛔ 零读零写：`frontend/**`、`design/**`（本腿对这两目录不读、不引、不转述）
- ⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件

## 入库清单（本腿会 commit 的路径，逐枚具名）

| 路径 | 是什么 |
|---|---|
| `docs/evidence/s1/223-hot-reload-wiring-v1.md` | 本表 |
| `.scratch/wisp/probes/223/v1/` | 全部原始读数与台件（含 `full-v.txt`、突变日志、loader 对抗输入台件） |

---

## 起手读数（本腿现跑，不复用他人数字）

- 起手时刻：`date` = **Tue Sep 29 15:44:04 CST 2026**；锚点 `git rev-parse --short HEAD` = **`577ae8c7`**（骨架提交后＝`4fc03f21`，本腿所有读数都在 `4fc03f21` 这枚含代落产码的锚点上跑）。
- 在飞检查：`ps -W | grep -i -E "go\.exe|go-test|test\.exe"` 起手＝**空**；`git diff --cached --name-only` 起手＝**空**（索引干净）。
- 别人地界的脏改动本腿一律不碰：`.gitignore`、`design/**`（16 枚删除＋若干新增）、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`（都在起手 `git status --porcelain` 里现读到，未读其内容）。

| 尺（本腿现跑，原文见 `.scratch/wisp/probes/223/v1/`） | 读数 |
|---|---|
| `grep -rn "CheckAndReload()" --include=*.go cmd/ internal/ \| grep -v _test.go` | **2 枚生产调用者**：`cmd/balldebug/main.go:243` ＋ **`cmd/wisp/config_reload.go:153`**；另 `internal/ball/hotkey_reload.go:21` 是注释。（实现腿表里写 `:150`，现读 **153** ⇒ 行号漂移，见"票面写的 vs 我读到的"） |
| `grep -rn "ConfirmLocked = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/wisp/config_reload.go:114`（实现腿写 `:111`，现读 114）；含测试则 11 枚命中，全在 `*_test.go` 与本腿台件目录 |
| `grep -rn "OnRestartPending = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/wisp/config_reload.go:115`（实现腿写 `:112`）；测试命中 3 枚（`manager_test.go:82/94`） |
| `grep -rn "OnReload = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/balldebug/main.go:244` ⇒ **`wisp run` 仍没有 reload 档的生产赋值点**（与实现腿"没做完"第 3 条一致，现读复认） |
| `grep -rnE "^[[:space:]]*go [a-zA-Z_]" --include=*.go cmd/ internal/ \| grep -v _test.go`（ban #1 连具名 `go f()` 一起扫） | **1 枚**＝`internal/observe/goroutine.go:281`，正是 d22scan 按文件路径豁免的那枚注册表实现 ⇒ **本票新起的 tick 不是裸协程** |
| `grep -n "time.Now\|\.Sub(\|time.Since\|Unix(" cmd/wisp/config_reload.go` | 只命中 `:37` 那一行**注释**（`//` 开头）。尺的判定：`tools/d22scan/main.go:759-770` 对每行 `TrimSpace` 后 `HasPrefix(code,"//")` 就 `continue` ⇒ 注释行不进 `wallclockRe`；非注释行**零命中** |
| `grep -n "time.Now\|\.Sub(\|time.Since" internal/config/manager.go` | **零命中** |

## AC#1 生产里真有人在轮询

〔待填〕

## AC#1 生产里真有人在轮询 —— **成立（带注：这圈 tick 没被 join，是同类违约的第二枚实例）**

- **生产调用者点数（本腿现跑，口径＝`grep -rn "CheckAndReload()" --include=*.go cmd/ internal/ | grep -v _test.go` 的命中行数）**：改前 1 枚（`cmd/balldebug/main.go:243`，旁支程序）→ 现测 **2 枚**＝`cmd/balldebug/main.go:243` ＋ **`cmd/wisp/config_reload.go:153`**（`reloadOnce` 里 `rep, err := rt.mgr.CheckAndReload()`）。
  ⇒ `cmd/wisp` 这一包从 0 枚变 1 枚，**且不是拿 `cmd/balldebug` 或测试交差**：`cmd/balldebug` 那枚本腿点数时单列、不算进"≥1 枚在 `cmd/wisp`"里。
- **可达链路（逐枚行号，本腿现读）**：`cmd/wisp/main.go:83-85`（`case "run": os.Exit(cmdRun(...))`，一次性进程）→ `cmd/wisp/run.go:154 runTextTask` → `:207 rt, code := assembleRuntime(s)` → `:329 func assembleRuntime` → **`:620 rt.startConfigReload()`** → `:105 startConfigReload` → `:119 observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)` → `:131 configReloadTick` → `:144 rt.reloadOnce()` → `:153 CheckAndReload()`。`startConfigReload` 在生产里**只有这一枚调用点**（现跑 `grep -rn startConfigReload --include=*.go cmd/ | grep -v _test`＝定义 `:105` ＋ `run.go:620`）。
- **触发形状（具名）**：**常驻 watchdog tick**，`time.NewTicker(1 * time.Second)`（`config_reload.go:88` 常量、`:132` 建、`:133 defer t.Stop()`、`:143 case <-t.C`）。不是文件通知（`manager.go:15-19` 自陈 no fsnotify），也不是"只有显式命令"。协程名 `watchdog` 是 D38b 名册里的在册常驻名（`internal/observe/goroutine.go:38-40` 的 `ResidentBaseline = 6`），本腿现读 `Spawn` 的签名是 `(name, owner, root, fn)`（`goroutine.go:262`），recover 边界与名册都在 Registry 里。
- **ban #1（裸 `go`）**：本腿按"实际扫什么"复跑尺，而不是按注释侥幸——`grep -rnE "^[[:space:]]*go [a-zA-Z_]" --include=*.go cmd/ internal/ | grep -v _test.go`（连具名 `go f()` 一起抓，对应 `tools/d22scan/main.go:695-711` 的 `*ast.GoStmt` 分支）⇒ **全仓生产只剩 `internal/observe/goroutine.go:281` 一枚**，正是 d22scan 按**文件路径**豁免的那枚 Registry 实现。⇒ 新起的 tick **不是裸协程**，有名字有 owner。
- **ban #4（墙钟时间差实现超时）**：`grep -n "time.Now\|\.Sub(\|time.Since\|Unix(" cmd/wisp/config_reload.go` ⇒ 唯一命中 `:37`，是 `//` 开头的注释行；`internal/config/manager.go` ⇒ **零命中**。仪器判定本腿现读＝`tools/d22scan/main.go:759-763`（逐行 `TrimSpace` 后 `code == ""` 或 `HasPrefix(code, "//")` 就 `continue`，然后才 `wallclockRe.MatchString(code)`）⇒ 注释行不进射程，而 `:37` 那行注释里逐字写着 `.Sub(time.Now())` 是**给禁令举反例**，仪器不报。**结论：这条路径没有任何计时判定**——是否重读由 `manager.go:156` 的 mtime+size 指纹相等性决定（内容指纹，不是"过了多久"），答复由通道结果决定；唯一的 deadline 属于 C18 gate 的 `ApprovalTimeout`，不在本票代码里。
- **"拿掉实现会不会红"＝会，且本腿有独立读数（不转述实现腿）**：`cmd/wisp/config_reload_223_test.go` 全文**既不 assign 三枚钩子、也不自己调 `CheckAndReload`**（现跑 `grep -n "CheckAndReload|ConfirmLocked|OnRestartPending|OnReload"` 在该文件里只命中 `:8` 的一句**注释**）。所以 `TestTicket223RunArmsTheReloadTick` 那句 `r.rt.mgr.Config().Ball.Size == 64` 与 `r.rt.cfg.Ball.Size != 64`（快照没动）**只能由生产 tick 自己跑出来**——没有 tick 就没有第二行。本腿那发全量 `-v` 里它 **PASS 2.22s**（原始件 `full-v.txt`），即"通电"是被观测到的，不是被断言在场。
- **注（本票新增的第二枚同类违约，具名登记）**：这圈 tick 起来之后**没有被 join**。`rt.reloadHandle`（字段 `run.go:273`、赋值 `config_reload.go:119`）**全仓再无读取点**（现跑 `grep -rn "reloadHandle|replyHandle" --include=*.go cmd/ internal/`＝4 处：2 枚字段声明＋2 枚赋值），`observe.Handle` 明明提供 `Done()`/`Err()`（`internal/observe/goroutine.go:200-210`）却没人用；`rt.close()`（`run.go:694-707`）只 `rt.reloadRoot.Cancel()`（`:701-703`），且 `:688` 的注释逐字承认「The cancel is not a join」。
  对照 `PLAN.md:2840-2843` D38(c)：「**任务完成必须等所有派生 goroutine 退出（WaitGroup）才算完成**」。⇒ **严格按字面，本票新起的 tick 与 `replyHandle` 是同一个毛病的第二枚实例**：没有取消之外的等待，进程收口时它可能仍在一发 `CheckAndReload` 里。
  **但这不推翻 AC#1 的字面判据**（票面只要"生产调用者 ≥1 枚＋触发形状具名＋不是墙钟超时"）。且这笔账**早已在册**：台账 `docs/reports/pending-and-issues.md:9366` 把 `rt.close()` 从不 join 记为 D38(c) 违约并派给**票 228 AC#4**（`run.go` 行号已从那时的 `:676-683` 漂到 `:694-707`）。⇒ **本腿的处理＝不重复立罪、只登记"223 新增了第二枚实例，票 228 的射程要从 1 枚变 2 枚"**。
  另：`config_reload.go:140-141` 在 ctx 退出时会补一行 `state=stopped` 审计，所以"退出时它是否真停了"是可查的，不是全黑。

## AC#2 三档生效级别各有读数 —— **成立但带注（立即档成立；reload 档缺生产 `OnReload` 赋值点＝半格，与实现腿自述一致；没有人造第四档）**

- **档位只有三枚，本腿现读复认**：`internal/config/schema.go:33-44`＝`TierHot "hot"` / `TierReload "reload"` / `TierRestart "restart"`；`grep -rniE "next.?session|下次会话" internal/config/ cmd/wisp/`（排除 `_test.go`）＝**零命中** ⇒ **没人偷偷造出"下次会话"那一档**（票面那一档今天不存在这条，实现腿已顶回、编排者已复认，本腿只确认"没人造新档"）。
- **立即档＝行为读数，不是字段读数（成立）**：`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 本腿全量 `-v` 里 **PASS 2.39s**，原始件里有它自己那行生产日志逐字 `msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]`。用例形状（现读 `config_reload_223_test.go:404-425`）：先断前提 `rt.modes.PermissionMode() == ask_every_step` → 种 `[risk] permission_mode = "auto_approve"` → **等卡期间再断一次仍是 ask_every_step**（`:416-418`，放宽没静默生效）→ 答 `yes` → `awaitLive` 轮询到 `PermissionMode() == auto_approve`（`:421-423`）。⇒ "改了→不重启→**决策链的档位真的变了**"。这条消费路径是真的：`internal/perm/store.go` 每次判定都读 `mgr.Config()`。
  另有 hot 档读数 `TestTicket223RunArmsTheReloadTick`（`ball.size` 56→64 而 boot 快照不动）与 `TestTicket223TighteningRaisesNoCard`（收紧热生效）同发 PASS。
- **第二／三档要"明确告知未生效＋为什么"**：重启档见 AC#7（独立用例，成立但带注）；**reload 档＝半格**——现跑 `grep -rn "OnReload = " --include=*.go cmd/ internal/ | grep -v _test.go` 只有 **`cmd/balldebug/main.go:244`** 一枚 ⇒ **`wisp run` 今天仍然没有 `OnReload` 的生产赋值点**，`manager.go:198-200` 那个分支在 `wisp run` 里恒不触发。审计行 `config: HOT-RELOAD state=applied hot=%v reload=%v restart=%v`（`config_reload.go:168`）里 `reload=` 这一列是有的，但**没有任何消费者**。⇒ **本格按票面字面判据不红**（票面只要求"立即档观测到变了＋重启档观测到明确告知"），但**"三档各有读数"这句做不到完整**，具名列入"没做完"第 3 条的复认，并建议另立票（语音管线重载属 S4/S5 射程，不是 223 能独立闭合的）。
- **拿掉实现会不会红**：立即档会（本腿 AC#6 那发突变把裁决挪回锁内时，同一条 tick 链路仍然在跑，所以这一档不是靠突变证的；它靠"用例自己不调 CheckAndReload、也不赋钩子"这一条——见 AC#1 的独立论证）。

## AC#3 放宽必带 L2 复确认（D33 正控）—— **成立**

- **正控走的是生产赋值路径**（这是票面最硬的一条，本腿三重现跑验证）：
  ① 生产钩子赋值点只有 `cmd/wisp/config_reload.go:114`；② `cmd/wisp/config_reload_223_test.go` 里**零处** `mgr.ConfirmLocked =`（现跑 `grep -rn "ConfirmLocked = " cmd/wisp/` ⇒ 唯一命中是产码 `:114`）；③ 那枚钩子真的会自己去读活配置：`confirmLockedLoosening` 在弹卡之前调 `rt.permissionMode()`（`config_reload.go:245`），而 `rt.permissionMode()` 的宿主是 `mgr.Config()`（`approval_always.go:151-156`）⇒ 它是**活的宿主路径**而不是测试桩。
- **实测读数（本腿那发全量 `-v`）**：`TestTicket223HandEditedFsLooseningCostsAnL2Card` **PASS 2.22s**，原始件里同段生产日志逐字 `msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]`；用例断的是 `card.Level == "L2"`（`:265-267`）、控制台出现 `[确认 L2 config.reload]` 与 `fs.allowed_dirs`（`:271-276`）、**答卡前内存里没有那条新目录**（`:291-300`）、答 `yes` 后审计 `config: D36-CONFIRM … result=allow` ＋ `config: D36-SECTION section=fs direction=loosen … effect=applied-after-L2` ＋ `awaitLive` 到内存真有 2 枚目录（`:306-316`）。**一张卡＝一次放宽**由 `:301-303` 的 `windowCount()==1` 钉住。
- **把复确认拿掉必须红＝会**：实现腿的 m2 已量（本腿不转述当读数），但本腿有**独立同向读数**＝任务二那发**锁内突变**——它没有摘钩子，只把调用位置挪回持锁段，`TestTicket223HandEditedFsLooseningCostsAnL2Card` 立刻转红（读数见 AC#6），这证明**这枚用例确实依赖那枚生产钩子的行为，而不是依赖"字段在场"**。
- **收紧不许弹卡＝现读＋实测都在**：`manager.go:313-332` 的 `switch` 按方向分三支——`len(loosen) > 0` 才 `plan.confirm = append(...)`（`:318`）；`len(tighten) > 0` 走 `plan.set`（`:324`）**根本不进 confirm 队列**；`default`（neutral）同。方向判定是逐键的真比较，不是"任何改动都弹卡"：`fsDirection`（`:457-466`，`setDirection` 加算 loosen、减算 tighten）、`riskDirection`（`:421-453`，含 `permission_mode` 按 rank 比大小）、`netDirection`（`:470-486`）、`pluginsDirection`（`:490-517`）。实测 `TestTicket223TighteningRaisesNoCard` **PASS 2.25s**，原始件同段逐字 `msg="config: locked section tightened, hot-applied" section=fs keys=[fs.allowed_dirs]`，用例还负断言整条审计里**没有** `config: D36-CONFIRM`（`:381-383`，即钩子压根没被问）与 `windowCount()==0`。
- **拒绝那一支 fail-closed**：`TestTicket223RefusedLooseningKeepsOldValues` **PASS 5.30s**（答 `no` ⇒ `result=deny` + `effect=kept-old-values` + 内存仍只有旧目录 + **文件里那行没被抹掉**）。生产代码里"除显式 allow 之外全是 deny"的四个出口现读齐：`rt == nil`（`:222-226`）、无 gate（`:227-232`）、无 root（`:233-238`）、`ans != tools.AnswerAllow`（`:256-263`）。

## AC#7 "重启后生效"那一档要有出口 —— **成立但带注（红的那一枚本票用例就是它）**

- 点数：改前 0 ⇒ 现测 **1**＝`cmd/wisp/config_reload.go:115`（`rt.mgr.OnRestartPending = rt.reportRestartPending`），生产读取点 `manager.go:201-203`（`len(rep.Restart) > 0 && m.OnRestartPending != nil`）。
- 出口内容（现读 `config_reload.go:278-294`）：两行审计（`state=restart-pending sections=%v effect=next-process-start tier=restart` ＋ `RESTART-PENDING detail=` 里带**为什么**与 `restartTierKeys`＝`app.language / app.autostart / app.single_instance`，`:299-301`）＋ 一行 stdout 句子（"这些段的改动本次运行不会生效……原因：……涉及：……"）。**不是"静默不生效"**：`planApp`（`manager.go:351-356`）把 restart 键**留在旧值**并进 `rep.Restart`，注释逐字承认"the silence here is the bug that file closes"。
- 独立性：它是**单独一枚用例**（`TestTicket223RestartTierSaysItWillNotApply`），没与 AC#2 立即档合并——票面这条禁合的要求满足。
- **注（＝AC#5 那枚红的归属）**：这枚用例**自身有读序竞态**，本腿安静单跑 5 发 **4 PASS／1 FAIL**（详见 AC#5）。失败的是**断言写法**（`config_reload_223_test.go:452` 对 stdout 单发读，而被等的最后一枚信号是 `:442` 的 stderr 审计行，生产代码 `config_reload.go:284`→`:288` 之间就是窗口），**不是产品句子缺失**——同发里两行 restart 审计都已到位。**判语：AC#7 的产品事实成立、验收用例不稳**。这一格要翻勾，必须先把那枚断言改成轮询（**本腿不修**：非实现者腿不产码）。
〔待填〕

## AC#5 整包终态读数（任务一：带 `-v`）—— **不成立（终态有一枚本票自己的红，且它不是争用）**

- 尺（逐字）：先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，再
  `go test ./cmd/wisp ./internal/... -count=1 -v -timeout 30m`。
  **起手 `date`＝`Tue Sep 29 15:46:22 CST 2026`／终态 `date`＝`Tue Sep 29 15:48:28 CST 2026`**（两个时刻都写进了原始文件首尾）。
- 原始输出全量落盘＝`.scratch/wisp/probes/223/v1/full-v.txt`（**7,057 行／841,116 字节，未接 `| head`／`| tail`**），起手锚点 `4fc03f21`（＝`577ae8c7`＋本腿的表；产码与 `248095d1` 逐字相同），`GATE_EXIT=1`。
- **起手前先自证没有并发门**：`ps -W | grep go.exe|test.exe` 空；跑期间该机只有这一发。
- 口径（**引用数字必带**）：顶层 `=== RUN` **1737**；顶层 `--- PASS` **1158**；顶层 `--- FAIL` **7**；`--- SKIP`（含缩进层）**7**；子测试 `--- PASS` **565**；子测试 `--- FAIL` **0**。
  包级：**26 枚包＝22 `ok`／4 `FAIL`**（`cmd/wisp`、`internal/ball`、`internal/panel`、`internal/risk`）。
- **红名册（用例名集合，7 枚，不用枚数）**：

| 用例名 | 包 | 归谁 |
|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | internal/ball | 历史在册（别人地界） |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | internal/panel | 历史在册 |
| `TestComposerContractTypesMatchFrontend` | internal/panel | 历史在册 |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | internal/panel | 历史在册 |
| `TestC21DesignTokensFourwayAgree` | internal/panel | 历史在册 |
| `TestResolvePerCallBudget` | internal/risk | 争用型假红（本腿复量见下） |
| **`TestTicket223RestartTierSaysItWillNotApply`** | **cmd/wisp** | **本票新增＝AC#7 那一枚用例自己** |

- 本票 12 枚具名用例在这发里的逐名读数：`TestTicket223RunArmsTheReloadTick` PASS 2.22s／`…HandEditedFsLooseningCostsAnL2Card` PASS 2.22s／`…RefusedLooseningKeepsOldValues` PASS 5.30s／`…TighteningRaisesNoCard` PASS 2.25s／`…ModeLooseningChangesTheRunningModeAfterAllow` PASS 2.39s／**`…RestartTierSaysItWillNotApply` FAIL 2.10s**／`…FailureSentencesAreDistinct` PASS 8.62s（4 枚子例全绿）／`…PanelInboundSaysHotReloadIsDisabled` PASS 0.01s／`…PermissionDeniedSitsInItsOwnSentence` PASS 2.15s（真 ACL 这台上跑通）／`TestConfirmLockedRunsOutsideTheManagerLock` PASS 0.03s／`TestCheckAndReloadAsksOncePerFileChange` PASS 0.04s。

- **那一枚红的复量＝它不是争用**（尺：`go test ./cmd/wisp -count=5 -v -run TestTicket223RestartTierSaysItWillNotApply`，单包、安静机、**15:50:33→15:50:45**，落 `.scratch/wisp/probes/223/v1/restart-tier-recheck.txt`）：
  **4 PASS／1 FAIL**（FAIL 在 `config_reload_223_test.go:453`，内容逐字＝"the operator is not told the edit will not land this run"，且 dump 出来的 stdout 里只有"答复监听已接入"＋"配置热加载已接管"两行，**没有重启档那句**）。
  ⇒ 同一枚用例在**没有别的门在飞**时仍然 5 发里红 1 发 ⇒ **争用解释被排除，这是用例自身的读序竞态**：`cmd/wisp/config_reload.go:282` 与 `:284` 两行审计**先**写、`:288` 的 stdout 句子**后**写，而用例 `config_reload_223_test.go:438/:442` 用 `awaitAudit` 轮询 stderr 到第二行就返回，接着 `:452` 对 stdout **只读一次、不轮询**（同文件里别的用例读 stdout 都走 `awaitLive`/`awaitAudit` 之后再读）⇒ 存在真实窗口。**本腿不修它**（禁区＋非实现者腿不产码），登记为 AC#7 的注与"该翻勾前必须补的一发"。
- **`TestResolvePerCallBudget` 的处置照派单走完**（尺：`go test ./internal/risk -count=3 -v -run 'TestResolvePerCallBudget'`，15:52:20→15:52:25，落 `budget-recheck.txt`）：
  **3/3 PASS，`ok github.com/CarlosShao/wisp/internal/risk 3.839s`**，逐发读数 `248908/252504/246368 ns/op`（预算 1.000 ms、约 4900 样本）。⇒ 这枚是争用型假红，**两发读数都写进表**：本腿 `3.839s`；编排者 `3.804s`／`3.813s`。**`thresholds.go` 与 golden 一字节未动（本腿全程未打开过这两枚文件）**。

- **复核编排者那一发（AC#5 里"由编排者代跑补齐"那节）对不对**：
  ① 它说 **23 包 ok／3 包 FAIL 共 6 例** ⇒ 我这发是 **22 ok／4 FAIL 共 7 例**。**包计数与例数都不同**，差异可由**一枚**用例完全解释＝`TestTicket223RestartTierSaysItWillNotApply`：它红 ⇒ `cmd/wisp` 从 `ok` 变 `FAIL`（包计数 23→22、FAIL 包 3→4），例数 6→7。**编排者那一发的读数在它自己那一刻是自洽的**（那枚 flaky 没命中），但"6 例"不是稳定终态。
  ② 它说红名册＝5 枚历史红＋`TestResolvePerCallBudget` ⇒ 与我这发的集合**逐名相符**（那 5 枚名字一字不差）。
  ③ 它说复量两发 `ok 3.804s`／`ok 3.813s` ⇒ 与我 `3.839s` 同向，**归因（争用）复认**。
  ④ **缺 `-v` 会不会改变结论**：**不会改变红名册**（Go 不带 `-v` 也打 `--- FAIL`，它那发的 6 例逐名可读、我已复认），**但结论本身要改**——它缺的正是"逐名绿册"（我这发才采到 1158 顶层＋565 子测试 PASS），而且**带 `-v` 的终态多抓出一枚本票自己的红**，那枚红在它那发不存在。⇒ **它那发不足以定"终态干净"，这一条它自己已在表里具名承认，本腿复认它的诚实、但读数口径要按我这发替换。**
- **"拿掉实现会不会红"**：会，且不是靠推断＝编排者/实现腿的 m1（注释掉三枚接线）与**本腿任务二那发锁内突变**都会让 `cmd/wisp`／`internal/config` 转红（见 AC#6）。本格自身的不成立点是**用例有竞态**，不是"没有实现"。

## AC#6 裁决不在锁内（任务二：专属突变）—— **成立（这一格今天有牙了）**

- 点数复认：`ConfirmLocked` 生产赋值点 改前 0 ⇒ 现测 1（`cmd/wisp/config_reload.go:114`，见起手读数表）。
- 现读形状（`internal/config/manager.go`）：`reloadMu` 全程持有（`:146-147`）→ `mu` 持有段内 `Stat` 指纹判定 `:155-159` → `LoadFile` `:160` → `plan(fresh)` `:169`（**只算不改**，`:243-297`）→ **钩子在 `mu` 之外被调**（`:172-173` 快照后 `Unlock`，`:178-190` 逐条裁决）→ 重新 `mu.Lock()` `:192` → `plan.commit()` `:193`。⇒ 头注释 `:21-31` 那句"Callbacks run OUTSIDE the lock so they may call Config()"现在对三枚钩子都成立，旧的 `applyLocked` 形状已被删除（全仓 `grep -rn "applyLocked" --include=*.go .` 在产码里零命中，本腿现跑）。

- **突变（本腿自己种的，`15:58:48 → 15:59:41 CST 2026`，原始件 `.scratch/wisp/probes/223/v1/ac6-mutation-lock-inside.txt`）**：
  改的对象＝**产码 `internal/config/manager.go`**，不是测试。派单猜"突变对象大概率在两枚 `*_223_test.go`"——本腿先打开读了它们，结论是**改测试证不了这件事**：那两枚用例断的正是"钩子被调时 `mu` 没被持有"，只有把产码挪回持锁段才是它们的反事实。
  具体两处（逐字）：① 删掉 `:170-173` 那三行（快照注释＋`confirm := m.ConfirmLocked`＋`m.mu.Unlock()`），并把循环体里的调用改回 `m.ConfirmLocked != nil && m.ConfirmLocked(c.section, c.keys)`＝**在 `mu` 仍被持有时裁决**（正是票面第 ② 格描述的旧 `applyLocked` 形状）；② 在循环之后、`m.mu.Lock()` 之前补 `m.mu.Unlock()`。
- **两枚声称"`Config()` 在等卡期间仍可读"的用例＝都转红**（口径＝`-count=1 -v`，逐名）：

| 用例 | 读数 | 红在哪一行、报的是哪句话 |
|---|---|---|
| `internal/config/TestConfirmLockedRunsOutsideTheManagerLock` | **FAIL 6.02s**（`FAIL github.com/CarlosShao/wisp/internal/config 6.064s`，`EXIT_CONFIG=1`） | `manager_223_test.go:115` 逐字 **"ConfirmLocked could not call Config(): the hook still runs inside the Manager lock"** ⇒ **这一枚就是 AC#6 那一格的专属牙** |
| `cmd/wisp/TestTicket223HandEditedFsLooseningCostsAnL2Card` | **FAIL 41.12s**（`FAIL github.com/CarlosShao/wisp/cmd/wisp 41.173s`，`EXIT_CMDWISP=1`） | `config_reload_223_test.go:264` "no \"config.reload\" card was displayed within 40s" ⇒ 生产钩子在 `mu` 里调 `rt.permissionMode()`（`config_reload.go:245` → `approval_always.go:151-156` → `mgr.Config()` → `manager.go:116` 抢自己已持有的 `mu`）**自死锁**，卡根本弹不出来 |

  同发里 `TestCheckAndReloadAsksOncePerFileChange` **PASS 0.02s**（它的钩子不读 `Config()`，所以锁内形状照样过）⇒ **本腿具名：这枚用例对"锁内/锁外"不敏感，它只钉 `reloadMu` 序列化，不能拿来充当 AC#6 的证据**（实现腿表 AC#6 一节把它和锁外证据并列写在一起，读起来会让人以为它也证了锁外——那是它表里的一处**口径含糊**，见推翻清单第 5 条）。
- ⛔ **不是装饰**：两枚红都拿到了，逐名可查，所以 AC#6 那一格的判据成立且被证过有牙。
- **还原协议（照派单逐字走完，`certutil -hashfile … SHA256`）**：
  - 改之前基线＝`32dd19893f8a0cfbc7f1461a7d1d4f39b6c67dcb0a6f9a6b3e0fd00bed41dc34`（blob `git rev-parse HEAD:internal/config/manager.go` ＝ `4cde0f491f107441a1714554841b670158b380a7`）
  - 还原＝`git cat-file blob HEAD:internal/config/manager.go > internal/config/manager.go`（**未用** `checkout`／`restore`／`reset`／`stash`／`clean`）
  - 还原后＝`32dd19893f8a0cfbc7f1461a7d1d4f39b6c67dcb0a6f9a6b3e0fd00bed41dc34` ⇒ **两个哈希逐字相同**
  - 收工自查＝`git status --porcelain -- internal/config cmd/wisp` 与 `git diff --stat -- internal/config cmd/wisp` **皆空**；本腿全程**没有 commit 任何源码**。
- **另附注（产品事实，与突变无关）**：卡片挂在 `rt.reloadRoot.Ctx` 上（`config_reload.go:246`），`rt.close()` 取消该根（`run.go:701-703`）⇒ 收口时未答的卡被放弃、落拒绝，不会"顺手变宽"。但**放弃 ≠ 等待**：`reloadHandle` 无人 join（见 AC#1 的注）。

## AC#4 不生效与读不到是四句话 —— **成立但带注（注＝有一句会说反，见推翻清单第 4 条）**


- 台件＝**`go test -overlay`**（物理件全在 `.scratch/wisp/probes/223/v1/overlay/`：`zz_v1probe_config_test.go`／`zz_v1probe_cmdwisp_test.go`／`overlay.json`；**`internal/config` 与 `cmd/wisp` 目录里不存在这两枚文件，本腿零源码改动**）。
  尺：`go test -overlay .scratch/wisp/probes/223/v1/overlay/overlay.json ./internal/config -count=1 -v -run TestV1ProbeLoaderClassification` ＋ 同 overlay 的 `./cmd/wisp -run TestV1ProbeProductionSentence`，15:54:12→15:54:17，原始输出 `.scratch/wisp/probes/223/v1/ac4-probe.txt`，两发 `EXIT=0`（台件只读数、不判等，所以它绿不代表实现绿）。
- **票面要求的四句话各有真实出口且互不共用**（全量门里逐名绿）：缺失（`os.Remove` 真文件）／语法错（写 `this is not toml [[[`）／权限不够（**真 ACL `icacls … /deny *S-1-1-0:(R)`**，`TestTicket223PermissionDeniedSitsInItsOwnSentence` PASS 2.15s）／热加载被禁用（第二宿主 `wisp panel-inbound`，`cmd/wisp/panel_inbound.go:211` 逐字 `auditf("%s", hotReloadDisabledPanelInbound)`，这枚宿主确实 `config.NewManager` 开同一个 `config.toml`（`:200`）且 `Confirm: nil`（`:220`）⇒ "它没法弹卡"这句是真的）。
- **`peekSchemaVersion` 的解析错误确实不再被丢掉**（这是＋72/−6 的核心断言，实测）：17 发对抗输入里凡是"读不出版本"的形状（`B2/C1 之外的`…逐条见原始件）产出的错误文本都以 `config.toml parse: toml: …` 开头（＝loader.go:86 那枚 `observe.Wrap(…, peekErr, "config.toml parse")`），生产句子＝`cause=syntax`。**改前**这一支不可能出现（错误被 `_ =` 丢掉 ⇒ ver=0 ⇒ 一律进 `applyMigrations` 报"cannot migrate from schema version 1"）。
- **但实现腿自述的边界规则被实测推翻了一半**：它写的是"**读不出版本＝语法错；读得出版本＝交迁移管线继续说话**"。现测——**"读得出版本"且版本恰等于 `SchemaVersionCurrent=2`（`internal/config/schema.go:28`）时，`loader.go:93` 的 `ver != SchemaVersionCurrent` 不成立，迁移管线根本没被调用**，错误由 `decodeStrict`→`formatDecodeError`（`internal/config/parse.go:104-113` 的 `toml.DecodeError` 分支）给出 `config.toml: line N, col M: expected character =`，再被 `describeReloadFailure`（`cmd/wisp/config_reload.go:352` 的 `HasPrefix("config.toml:")` 分支）归入 **`cause=invalid`**，而那一句中文逐字是「config.toml **语法没问题**，但内容被校验拒绝（值不合法或引用解不开）」。**文件真正的毛病就是语法错**（`expected character =`）⇒ **这句话在这个形状上是假的＝"语法错"与"内容不合法"说反**。命中发数（17 发中 4 发）：
  | 对抗输入 | declaredSchemaVersion | 生产句子 | 该不该是语法错 |
  |---|---|---|---|
  | C1 注释行以 `[` 开头（`# [fs] …`）＋版本 2＋坏行 | `2,true`（注释被正确跳过） | `cause=invalid` "语法没问题" | **说反** |
  | E1 CRLF＋版本 2＋坏行 | `2,true`（`\r` 被正确剥掉） | `cause=invalid` | **说反** |
  | G1 `schema_version=2` 无空格＋坏行 | `2,true` | `cause=invalid` | **说反** |
  | L1 行首缩进的版本行＋坏行 | `2,true` | `cause=invalid` | **说反** |
  | J1 版本 99＋坏行 | `99,true` | `cause=invalid`（detail 本身是"newer build"，诚实；句子前半句仍多说了"语法没问题"） | 半说反 |
  | H1 同名键两次 | `2,true`（取第一枚） | `cause=syntax` | 对（靠 `formatDecodeError` 的兜底 `Wrap("config.toml parse")` 恰好救回） |
  | A2 版本 1＋坏表头（`migrate_test.go:123` 那枚既有断言的形状） | `1,true` | `cause=migration` | 对，且既有断言一字未放宽、仍在绿包里 |
  | B1 版本写在 `[app]` 里 | 交迁移管线→`unknown key "app.schema_version" at line 4` | `cause=unknown-key` | 对 |
  | D1/D2 多行字符串里有一行以 `[` 开头 | `0,false`（扫描器撞假表头就停） | `cause=syntax` | D1 对；**D2 说反方向＝"文件声明了版本却因为假表头读不出"，本腿判它对**（它确实是语法错，且它没有把两句调换成迁移） |
  | F1 BOM＋版本 2 | `0,false`（BOM 让 `schema_version` 匹配不上） | `cause=syntax`（detail 逐字 `invalid character at start of key: ï`） | 对 |
  | I1 首行就是 `[fs]` | `0,false` | `cause=syntax` | 对 |
  | K1 版本被写成字符串 | `0,false` | `cause=syntax` | 对（严格讲该归"值不合法"，但本腿不据此加码） |
  | M1 合法 v1 | 走真迁移 | `APPLIED hot=[] reload=[] restart=[] locked=0`，文件被重写并留 `.bak-1` | 对（唯一一发改了文件，符合 SPEC-03 §4.4） |
- ⇒ **结论**：AC#4 的字面判据（四件各一句、不许合成一句）**成立**；`语法错` 这句话从"生产里不可能出现"变成"可达且有 9 发读数"，这是真修好。带注的部分＝决定句子归谁的**实际规则不是"能不能读出版本"，而是"读出的版本等不等于 2 ＋ go-toml 抛的是哪一种错误对象"**，于是一批"声明了当前版本又语法坏"的文件被告知"语法没问题"。这条具名进推翻清单第 4 条，**建议另立票**（修法要动 `describeReloadFailure` 的分类或 `formatDecodeError` 的 Detail 前缀，都不是本腿权限）。
- ⛔ 本腿**没有**为了让它绿去改判据，也没有碰 `internal/config/migrate_test.go:123` 那句既有断言（现读该行仍在，包仍绿）。

## AC#6 裁决不在锁内（任务二：专属突变）

〔待填〕

## AC#7 "重启后生效"那一档要有出口 —— **成立但带注（AC#5 里那枚红就是它的用例）**

〔待填〕

## 编排者代落＋代跑这两件事的独立裁决

**代落＝commit `248095d1`（五枚文件的未提交改动，含既有产码 `internal/config/loader.go` ＋72/−6）**

- 合规判定：**程序上可接受，但要记一笔角色重叠**。三点现读依据：
  ① 码**不是编排者写的**——`git show --stat 248095d1` 现读：6 枚文件、＋284/−74，全部是 `223-r1` 的现场；"实现者"仍是那枚死腿 ⇒ **D22 双角色（裁决者≠实现者）没有被破坏**，因为我（裁决者）是第三枚。
  ② 禁区实测复认：**`git show 248095d1 | grep -E "^\+(func|type|var|const) [A-Z]"` ＝零命中** ⇒ "没新增导出名"这句在**产码**上为真；`ec7a034d` 里新增的 11 枚顶层导出函数**全部是 `_test.go` 里的 `Test*` 用例入口**，不是产品 API（口径写清，别让下一位以为等价）。
  ③ `config.reload` 现跑 `grep -rn "config\.reload" --include=*.go .`（排除 probes 与 `_test.go`）⇒ **只命中 `cmd/wisp/config_reload.go:86` 那枚常量**，在 `internal/agent`／`internal/tools`／`internal/panel` 里**零命中** ⇒ "不是注册工具、不是 C17 方法"这句为真，C17 方法名册没被扩。
- ⚠ **要记的那一笔**：同一枚编排者在这票里同时是**派单者＋代落者＋代跑者＋翻勾人**。AGENTS.md §0 第 3 句与 `SPEC-12 §4.3` #1/#3 立双角色，就是为了不让"接线的人"同时是"宣布接通的人"。他这次没有翻勾、把缺 `-v` 的口径差具名写进表里（`docs/evidence/s1/223-hot-reload-wiring-r1.md:103` 逐字承认"仍缺一发带 -v 的终态读数"），**处置是合规的**；但如果"代跑"变成常例，双角色就名存实亡。**建议**：把"编排者代跑"限成一枚具名例外（先例 `201-r1`／票 35 的"代落"可以延用，"代跑读数"不要延用），并在 `AGENTS.md` §1.4 之外另记一条。

**代跑＝AC#5 那一发没有 `-v` 的门**

- 他那发的数字**自洽但不稳定**：现读他的原始件 `223/wisp run`…（口径：23 包 ok／3 包 FAIL／6 例）；我这发同一条尺加 `-v`＝**22 包 ok／4 包 FAIL／7 例**。差的那一枚**全部**由 `TestTicket223RestartTierSaysItWillNotApply` 解释（它红 ⇒ `cmd/wisp` 由 `ok` 变 `FAIL`）。⇒ **他的读数不是错，是"单发"**；错处在于把单发当终态写进表标题（"由编排者代跑补齐"），**终态其实没被补齐**——见 AC#5 的 5 发复量（4 PASS／1 FAIL，安静机）。
- **缺 `-v` 会不会改变结论＝会**。（a）红名册部分不会（Go 无 `-v` 也打 `--- FAIL`，他那发逐名可读、我复认）；（b）**"逐名比绿"那一半只有 `-v` 才有**（我这发才采到 1158 顶层 PASS＋565 子测试 PASS＋7 SKIP）；（c）**带 `-v` 的这一发抓到了他那发没抓到的一枚本票红**——所以"补一发带 `-v`"不是形式主义，它确实换了结论。
- 他那发的复量 `ok 3.804s`／`ok 3.813s` 与我这发 `ok 3.839s` 同向，**争用归因复认**；`thresholds.go`／golden 全程零接触（本腿连文件都没打开过）。

## 跨票攻击：票 226 的"手改／不认领"在真轮询路径下是否仍成立

**先说定性三行（涉及"绕过／说谎"字样）**：①现象出现在哪＝本机一个 `%APPDATA%\wisp\config.toml` 的**读取时序**与**一枚测试用例的断言方向**；②有没有本机被入侵的证据＝**没有**，这是我们自己的轮询与别人自己的用例撞了语义；③最坏后果是什么形状＝**用户在文件里手改的一枚可热加载键，可能在下一次启动前就被这个进程认领走**（值本来就该热生效，不是放宽权限、不是越门），以及**一条绿色用例其实在靠"跑不满 1 秒"过关**（读数失真，不是行为失真）。

- **票 226 §⑥ 第 8 条的原话本腿现读复认**（`docs/evidence/s1/226-config-write-no-clobber-v2.md:170`）：那句"生产侧尚无轮询者"在**当时**是对的；**现在错了**——223 已把轮询接上（AC#1）。
- **"写完之后再手改一枚键、下一次轮询要能看见"这一跳今天有没有常驻用例？＝没有。**
  尺：`grep -rn "CheckAndReload|live(|rtHook|newReplyHost" --include=*226*.go cmd/wisp/ internal/config/` 现跑 ⇒ 226 的 14 枚用例里，**包内手动调 `m.CheckAndReload()` 的仍是那 5 处**（`internal/config/writeguard_226_test.go:193/203/219/254/295`），全都在裸 Manager 上、且 `:253` 手动把钩子设成 `return false`；走真宿主的只有 `cmd/wisp/always_write_no_clobber_226_test.go:39` 一枚，而它证的**不是**"下一次轮询看得见"。⇒ **"生产路径无钉"这一条到今天我仍然要写出来**。
- **更硬的一枚现测（本腿新增，票面没写）**：那枚唯一走真宿主的 226 用例，**它的断言方向现在和 223 的 tick 相反**，而它靠时序侥幸过关。
  - 它在装配之后往文件里追加 `[app] theme = "light"`，注释逐字写着「while the process is running and **nothing polls**」（`always_write_no_clobber_226_test.go:43-44`）——**这句话从 `248095d1` 起就是假的**，`wisp run` 现在每 1s 轮询。
  - `[app] theme` 在 D36 里是**热档**（`internal/config/manager.go:343-349`：`themeChanged` → `plan.set` → `rep.Hot`）⇒ 只要这圈 tick 抢到一次，内存里的 `App.Theme` 就会变成 `light`。
  - 而它在 `:102-105` 断的是「`rt.mgr.Config().App.Theme` 必须仍是 `dark`，否则就是"受守卫写把用户手改进的档塞进了运行中的进程"」。
  - **实测**：尺 `go test ./cmd/wisp -count=3 -v -timeout 6m -run TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（16:03:17→16:03:22，原始件 `cross-226-tick-claims-handedit.txt`）⇒ **3/3 PASS，但三发耗时只有 1.11s／0.91s／0.93s**——**整条 run 从手改到收口都跑不满那 1s 的 tick**，所以"没被认领"是**时序巧合**，不是行为保证。同发里 writeguard 那句 `kept_in_file_not_in_memory=[app.theme]` 三发都在（`:14/:26/:38`），说明"不认领"发生在**写路径**那一侧、而不是**轮询**那一侧。
  - 对照：本票自己的 `TestTicket223*` 用例耗时 2.1–2.4s，`state=applied`/`restart-pending` 都在约 1s 内被观测到 ⇒ **同一枚 tick 在跑得久一点的宿主里确实会认领热档手改**。
  - ⇒ **这不是 223 的洞，但必须登记**：223 的 tick 让 226 的一枚用例的**前提失效**（"nothing polls"）且**断言语义与 D36 热档相反**，今天只有时序在替它兜底。票 227 的派单已经预感到这件事（`.scratch/wisp/issues/227-...md:40` 逐字「必须排在票 223 之后（223 一旦接上生产轮询，AC#2 那一支的行为面会变，现在钉的基线会被它冲掉）」），但 227 的 AC#1–AC#5 里**没有一格**收这一笔。⇒ **建议落点＝票 227 的 AC#3 矩阵加一格"改值（热档）→下一次真轮询可见"，或另立一票**；本腿**不改 226 的任何代码**，只裁与登记。

## 推翻清单（实现腿＋编排者）

〔待填〕

## 票面写的 vs 我读到的（不一致单列）

〔待填〕

## 没做完／判不了（具名清单）

〔待填〕

## 门禁与尺：时刻表

〔待填〕
