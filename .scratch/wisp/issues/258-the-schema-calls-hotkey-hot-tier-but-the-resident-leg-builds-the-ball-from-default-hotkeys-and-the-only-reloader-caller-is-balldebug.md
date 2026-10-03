# 票 258 — `[hotkey]` 在 schema 里被写成"hot 档：改了会重新注册"，而**常驻那条腿建球时根本没读配置**（用的是写死的 `DefaultHotkeys()`），全仓**唯一**的非测试 reloader 调用者在调试台件里

**立票时刻**：2026-10-02 15:0x，锚点 HEAD `9588138a`（`dev`）
**来路**：只读腿 `e2e-panel-1` 的第四格未定义项④（台账 `A531` §4 里我写过"待量完 `256-a1` 之后立票，⛔ 现在不自填修法"）＋ `256-a1` 已交件（`A534`）⇒ 这一枚就是那笔欠账的兑现。
**与既有票的关系**：**不并进票 255**——255 那族是"配置说了生效、没人读"，修法在**登记表与回执**；这一枚的修法在**装配**（谁把配置里的四枚热键交给球、谁在改过之后重新注册），两族的写面与判据形状都不一样。⚠ 它也不是票 245（那枚管的是"裸 `Esc` 被注册成全局热键"，射程＝绑与不绑；这一枚管的是**绑的哪一枚、改了以谁为准**）。

## 现量（编排者 15:0x 自跑，⚠ 引用前先重跑，别把这几行当常量）

1. **承诺那句话在 schema 里**：`internal/config/schema.go:176` 逐字 `// HotkeySection is [hotkey]; hot-tier (hotkeys re-register on change).`，四枚字段在 `:179` 起的 `HotkeySection`；`internal/config/manager.go:278` 的 hot 应用表里**确实有** `{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }}` ⇒ **内存里的值会跟着文件变**。
2. **但常驻那条腿建球时不吃它**：`cmd/wisp/resident_ball_windows.go:171` 逐字 `Hotkeys:  ball.DefaultHotkeys(),` ⇒ 球一开始用的就是**写死那四枚**，`Config().Hotkey` 从没进过构造参数。
3. **会重新注册的那台机器只接在调试台件上**：`cmd/balldebug/main.go:237` `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), …)`＋`:243 bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }` ⇒ **全仓 `NewHotkeyReloader` 的非测试调用点＝这一枚**（尺＝`grep -rn "NewHotkeyReloader" cmd internal tools --include=*.go`）。
4. **重载的引擎本身在常驻腿是活的**：`cmd/wisp/config_reload.go:153` 在跑 `rt.mgr.CheckAndReload()`，而 `:101-108` 那段注释逐字写着 `startConfigReload` 由装配根在门存在之后调用 ⇒ 所以缺的不是"重载有没有发生"，**是重载之后没人把新值交给球**。
5. **一票既有钉射程要先量清**（⛔ 本票不许为了落地去改它们）：`cmd/wisp/resident_ball_228_windows_test.go:45`、`cmd/wisp/resident_ball_live_228_windows_test.go:144/160/163` 那族读的是启动判决串里的 `hotkeys live %d/4`；线程形状受 10-01 那批裁定约束（球/面板都在自己的线程上，投 `ui-sta` 会冻外层泵——见票 33 十一裁）。

## 要建什么（本票只管"改了以谁为准"这一件；⛔ 不碰界面、不碰凭据）

- [x] **AC#0（本票第一格，且是闸门）＝先把三问答出来，⛔ 不许直接开写**：① 常驻腿要把配置里的四枚热键交出去，**构造期**该改哪一处、`DefaultHotkeys()` 那枚默认还要不要留作"配置缺失时的那一份"；② **改过之后**谁去重新注册——是把 `ball.NewHotkeyReloader` 那台机器接进常驻腿，还是让 `config_reload.go` 那一跳多调一次注册，两形的**代价与线程约束**各是什么；③ 四枚里有一枚注册失败（被别的程序占了）时，**今天那句 `hotkeys live %d/4` 会怎么说、应该怎么说**。交件判据＝逐处带 `file:line`＋尺读数；量不到的**具名说量不到**，⛔ 不许用"应该没问题"填空；⛔ 不许改任何产码。**〔16:4x 编排者翻勾：`258-a1` 交件 `1286382d`（104 行，占位 0，禁 Go 遵守），我抽验两把尺（`ApplyHotkeyDefaults` 真身 `hotkey_windows.go:87`＋唯一非测试调用者 balldebug:239；`NewHotkeyReloader` 全仓非测试调用点确实只有 balldebug:237）逐字命中；选形与裁定在下面第 7 节〕**
- [x] **AC#1 只在 AC#0 交完、并由编排者落一枚具名 `A##` 批准选形之后才许动**：`wisp run` 与**常驻腿**两条入口建出来的球，四枚热键**以 `config.toml` 的 `[hotkey]` 为准**；配置文件缺失／那一节缺失时**退回 `DefaultHotkeys()` 并说得出这句话**（不许静默换成另一套）。
- [x] **AC#2 说实话的判据（负向必配正控）**：种一发"把 `summon` 改成别的组合键" ⇒ **下一次建球注册的就是新值**（正控）；⛔ 反向不许做成词面尺（不许只扫注释里有没有 `hot-tier` 那个词），要问能力。⚠ 若量出来"重新注册必须回到 `ui-sta` 才安全"，那**属票 33 那批线程裁定的射程**，停下来上报、不许自填。
- [x] **AC#3 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。⛔ 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不许动。

## 禁区（本票全程）

- ⛔ **未定义即停**：票 33 十一裁（面板专用 STA 线程＋库 `Run()` 泵；投 `ui-sta` 会冻外层泵）是既有裁定，动它＝**人工批准**。
- ⛔ 不许为变绿放宽任何断言；不许 `t.Skip`；不许把 SKIP 读成通过（`tools/d22scan/runtests.sh:98` 把 SKIP 判红）。
- ⛔ 凭据值绝不进对话／日志／表（只写变量名）。
- Git：只 commit 不 push；**add 与 commit 同发一条命令、commit 必带显式 pathspec**（`A532` 那枚归属事故就是这条没做到）；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。
- `frontend/**`／`design/**` 两层禁令；`grep`/`find` 显式根（`cmd internal tools docs scripts .scratch`）。
- ⛔ **界面侧那一半不在本票**：设置页里"改热键"那个输入面属界面那支（机主自己带给他在用的那枚 agent）；本票只交 Go 侧"改了以谁为准"。**C17 白名单既有名字不动**；⛔ 不新增 D34 工具行。

## 排程

写面＝`cmd/wisp`＋`internal/ball` ⇒ ⛔ 与 `198-r2`（在飞，写 `cmd/wisp`）**串行**。AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与写腿的整包测量同机即互洗）。
**默认不排落地腿**：本票排在票 198 与票 248 那两片交完之后。撤销口令「**258 撤**」。

---

## 7. `258-a1` 收档 ⇒ AC#0 翻勾＋**选形＝形 A**（编排者裁定，账 `A538`，2026-10-02 16:4x，现量 HEAD `3df8b82b`）

**1. 交件核过**：`1286382d`（76/8），census **104 行／17,841 字节**、占位 0、全程零 Go 命令、写面 porcelain＝0。**它没推翻我票面任何一句**（五条现量逐字复认全中），两处补充：① **票面"两条入口建球"对 `wisp run` 是空集**——`wisp run` 产码零 `ball.` 引用、**不建球**（这处补充记我：我立票时把 run 腿当成了建球入口之一）；② "从没进过球"的机制更具体＝常驻进程里有**两枚** config.Manager（任务管线＋面板链），都进不了 `startResidentBall` 的参数表。

**2. 三问判语（读数全在 census §1–§3，都带 file:line）**：① 形状四枚**同名同序**但**缺省值不一致**——schema 只有 `cancel` 带 `default:"Esc"`（`schema.go:182`），summon/mute/panel 无 tag ⇒ 配置路产物三枚 `""`，缝由 `ApplyHotkeyDefaults`（`hotkey_windows.go:87`）补，**唯一非测试调用者＝balldebug:239**；② `config_reload.go:153` 之后**无任何产码**把新值交给球（`OnReload` 只吃 reload 档、hot 档不触发；`manager.go:278` 闭包只换内存值无人被通知）；③ 失败句链条完整（`HotkeyTaken` 逐行进 `Problems()`＋verdict `live 2/4` 形），但 **`Problems()` 那句"去 [hotkey] 挑组合"指向一扇常驻腿今天没有的门**——诚实缺口，AC#2 的尺要顺手钉它。

**3. 两形代价（census §2，本编排者裁）**：**形 A**＝把 `ball.NewHotkeyReloader` 那台现成机器接进常驻腿（balldebug:237-253 有跑通的全套：diff/panic 隔离/轮询闭包）；代价＝rb 与 mgr 两对象不同层（建球 `resident_windows.go:163` 早于任务管线 `:206`），要在装配次序里给桥一个位置。**形 B**＝`config_reload.go` 那一跳多调一次注册；代价＝`agentRuntime` 无 ball 字段（要动装配根结构）、run 腿无球、"reload tick 在非 ui-sta 线程调 rebind"的新线程形状、`RebindHotkeys` 全量 churn。**线程形状两形都合法**（rebind 经 `uiRun` post-and-wait 到球自己的 ui-sta＝D38b 冻结名第 1 枚；`Check()` 注释"Safe from any goroutine EXCEPT the ball's UI thread"）⇒ **不触票 33 十一裁的停手上报线**。

**4. ★ 选形＝形 A（具名裁定五样齐）**：**文件**＝`cmd/wisp/resident_ball_windows.go`（`Hotkeys: ball.DefaultHotkeys()` 那一发改吃 `ApplyHotkeyDefaults(mgr.Config().Hotkey)` 那条链）＋`cmd/wisp/resident_windows.go`（桥的装配位）；**理由**＝零新发明（balldebug 全套现成、票 246/228 的装配次序判例都在）、不动 `agentRuntime` 结构、不新增线程形状；**边界**＝⛔ 两形共有的已登记残留**不归本票修**：rebind 丢 in-flight Esc borrow（`ball_windows.go:807-812`，票 245 在册）——落地腿不许顺手做；⛔ `Problems()` 那句指向"没有的门"只许在 AC#2 的尺里钉住诚实形（文案改动若要做＝另落一格，不许夹带）；**撤销口令**＝**「258 改形 B」**。

**5. 落地腿（`258-r1`）判据预告（AC#1／AC#2，暂不派——cmd/wisp 测试面被 `253-r1` 占着，包级互斥）**：AC#1＝常驻腿建出的球四枚热键以 `config.toml` 为准、缺失时退 `DefaultHotkeys()` 且说得出那句话；AC#2＝正控（种一发改 `summon` ⇒ 下次建球注册新值）＋反向不词面尺＋**三枚既有钉全盖不到这一格**（census §3：`:45` 是常量、`:144` 无断言、`:160/:163` 只钉两极端）⇒ 尺新造不撞钉；正控还要钉"改配置后 reloader 真把新值交给球"（现在零产码路径）。**写面＝`cmd/wisp`＋`internal/ball`（只读后者）**。

## 6. 收只读普查腿 `258-a2`＋落地腿必配四枚（2026-10-03 09:4x；台账 `A560` §1；票面原话一字未改）

交件凭据（盘上尺复量）：`.scratch/wisp/probes/258/a2/census.md` **606 行／55,101 字节**，六枚 commit `3736f0dd`→`06d89e25`→`5c3c22d8`→`761b5d45`→`0c72c253`→`5d4f5343`；只读、零 Go 命令、AC 框未碰。⚠ 它头注释那一处推翻我票面：`internal/ball/hotkey_reload.go:16` 自述 `b.HotkeyConfig()` **仓里根本没有这枚方法**，真名 `ConfiguredHotkeys`（`internal/ball/ball_windows.go:840`）——⛔ 照注释抄即编译不过（该文件属本票只读面，⛔ 不改字，落地腿按真名接）。

**★ 我亲自复量过的三条（不是转述）**：
1. **`cmd/wisp` 连 `RebindHotkeys` 的边都是零**：`grep -n RebindHotkeys cmd/wisp/*.go | grep -v _test` ＝ **0 枚**；非测试调用者只有 `cmd/balldebug/main.go:237` 一处（测试三枚 `internal/ball/hotkey_live_test.go:159`、`hotkey_status_test.go:458/:511`）。⇒ **比 §45 那句"唯一 reloader 调用者是 balldebug"更硬：常驻腿这侧不是"没接对"，是一条边都没有。**
2. **库静默把零值补成默认**：`internal/ball/ball_windows.go:154` 逐字 `opts.Hotkeys = DefaultHotkeys()`（`:67` 注释自述 zero value = DefaultHotkeys）⇒ **AC#1 那句"退回了默认"只能由 `cmd/wisp` 说**（库里说不出"我退回了"，因为它补得无声），挂载点候选 `cmd/wisp/resident_ball_windows.go:202-208`／`:316-326`。
3. **`pRegisterHotKey` 全仓产码只有一枚调用点**（`internal/ball/hotkey_windows.go:391`；`win32_windows.go:42` 是 proc 声明），而它只被两处用：跳过 cancel 的 idle 遍（`:458` 注释逐字"The one slot an idle ball must not register"）与**写死 `VK_ESCAPE` 的借用遍**（`:118` `vkEscape = 0x1B`、`:521` `escBorrowAcc`）。⇒ 推论（我复量后认它成立）：**配置里那枚 `[hotkey] cancel` 组合键今天在任何路径都不会成键**；AC#1 那句"四枚以 config 为准"按现读码只对三枚成立。⚠ **这一条不在本票射程**（票 245 裁的是"idle 不绑"＝设计；本条是"借来的键不吃配置"＝另一枚洞），已登记为**待独立立票**（`A560` §6），落地腿⛔ 不许顺手修、也⛔ 不许把 AC#1 的文案写成"四枚都吃配置"。

**落地腿 `258-r1` 必配四枚（派单要逐条写进去）**：
- **(a) 同批改票 255 那把对齐尺**：`cmd/wisp/config_readers_255.go:110-113` 把 `cmd/wisp/resident_ball_windows.go:171` 当**证据行**引用，且 `:162` 的 `sectionReadSites["Hotkey"]` 今天**只允许 balldebug** ⇒ 形 A 一落地这两处必红。⚠ 该文件此刻**未被 tracked**（`?? `）、`config_reload.go` 是 `M`，**行号以 `255-r2` 交件之后重跑为准**，⛔ 不许拿这里的行号当常量。
- **(b) 收口挂死的真雷**：`uiRun` 是 post-and-wait（`internal/ball/ball_windows.go:741-752`），而 `internal/ball/sta_windows.go:240-247` 在 `hwnd==0` 时**丢任务不执行** ⇒ **Close 之后来一次 `Check` 会永久挂住退出序列**。落地腿要带 nil/已关守卫，判例形状现读 `cmd/wisp/resident_approval_windows.go:393-409/:429-434/:487-490`。
- **(c) rebind×借还那枚自钉只能装 `winlive` 层**：CI 的 `cmd/wisp` 测试步不跑该 tag（尺＝`grep -E 'tags|winlive' scripts/wisp-cli-tests.sh .github/workflows/ci.yml` 零命中〔**射程判断，非内容引用**〕），而非 winlive 的同 tag 球测试 `cmd/wisp/resident_approval_246_windows_test.go:309-310` 故意容忍无桌面，这枚钉三条断言全以"真借到键"为前提。⚠ 五枚零件今天在 `cmd/wisp` 写面内**全齐**（`startResidentBall`／`bindBallHost:139`／`AskOnTaskRoot:261`／导出的 `NewHotkeyReloader`／`observeEsc246:444`／`requireIdleCancelSlot246:350`），测试里自搭一枚桥即可闭合"改 config ⇒ Check ⇒ rebind"整条链 ⇒ **推翻 `245-c1` §1.3 那句"形 A 落地后残余才成为可发生的事"的"对钉也成立"那半**（对产线成立）。
- **(d) 丢借用不是"卡片少个键"，是"卡片到点必执行"**：`internal/agent/approval/gate.go:319→:332` 逐字 `ANSWER-EXPIRED decision=timeout->execute`（`:338`），窗口 3s＝`queue.go:116`；且 `Standby` 不进 `Problems()`（`hotkey_windows.go:365`）⇒ **零症状**。这正是 §4 里"⛔ 已登记残留不归本票"那一枚的**代价读数**，落地腿要在证据件里把它写成"为什么必须自钉"，⛔ 不许因为"今天产线不可达"就不装。
**今天可达性另加一条它的新事实**：票 255 的回执今天**已替这格说实话**（quiet 分支 `cmd/wisp/config_reload.go:189-194`），⚠ 但**主语写错**（把"吃默认键"归给不建球的 `wisp run`，实为常驻腿 `resident_ball_windows.go:171`）；且 `OnReload` 对 hot 档结构性不通（`internal/config/manager.go:198` 只吃 `rep.Reload`，`tiers.go:38` 把 hotkey 记成 `hot`）⇒ 归 `255-r2` 那一格的措辞，⛔ 本票不修它。

---

## 编排者翻勾节（2026-10-03 22:0x，锚 `bd5c049f`/`d9bce7ce`；258-v1 非实现者判决书收档）

**AC#1/AC#2/AC#3 三格翻勾**（口令「258 退回」），凭据＝`.scratch/wisp/probes/258/v1/verdict.md`（145 行，13 份读数档）＋v1 判语：
- **AC#1 成立但带条件**（我照判语翻勾，两枚条件写死）：①"节缺失⇒终值 DefaultHotkeys() 但档位词印 config"的措辞裁定点＝**归 258-r2 一句修**（印 "defaults" 才诚实）；②R1 洞（6 主判据 untracked）已由 `bd5c049f` 补 commit 封口。
- **AC#2 成立但带条件**：死腿"非真窗"自报**判错**（v1 纠正＝默认档就是真窗读数，PASS 非 SKIP）；真缺口＝两枚 winlive 尺子自身缺陷（rebind 枚 grep 等号格式 vs JSON sink、occupied 枚赌 Ctrl+Alt+U 未被占无 squat 兜底）＝**尺子修复归 258-r2**，不是行为缺陷。
- **AC#3 成立**：禁区零触碰、agentRuntime 零结构改动、Esc borrow 零顺手修（+19 行全是 Debug squat seam）、d22scan clean。
- **恒真四突变**：M1 摘桥红 2 枚具名（Rebinds＋Occupied）＝桥有牙；M0/M2/M3 默认档绿＝**判据覆盖洞**（6 主判据自带闭包注入不读装配根）——v1 具名登记，不动判据、留给后续程补"读装配根"那一发。
- **意外收获**：自写正控首发撞上 `Ctrl+Alt+W` 被第三方真占用 ⇒ AC#3 的占用语义被实测（`hotkeys live 3/4`＋`binding=Ctrl+Alt+W` 逐字）＝票 260 AC#0 关心的"占用时说什么"有了活体样本。
- **R2 纠我**：A576 记的"死腿报的红名是 C18 形"实测倒的是 `TestTicket223…`（solo PASS＝挤压间歇）——归因换名，间歇形结论不变。
