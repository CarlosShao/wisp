# 228 — 球和托盘**根本不在能干活的那条进程里**：`wisp run` 一发即退，常驻那条腿又没有审批门／没有任务 ⇒ 票 201 剩余四格缺的不是"接一根线"，是**一枚宿主**

- Status（**09-29 16:1x 收 `228-a1`，台账 `A439`）：AC#0 成立并翻勾**（代价表＝`.scratch/wisp/probes/228/a1/cost-table.md`，69,717 字节／354 行／10 枚提交，编排者自己复跑六把尺全部对上）；**方向由编排者按 owner 长期原话裁死＝走甲**（把审批门与任务通路接进**常驻**那条腿），⛔ 不是"待他选"——他的原话是「**别人怎么做的，你就怎么做**……**必须特么做完整功能**」，而 D2（`PLAN.md:74/83-88/:103`）本来就把球＋托盘＋快捷键＋Job 画在同一枚常驻主进程里。**撤销口令：「228 走乙」**（要改回乙只需说这四个字，本票现量与代价表两支都还在）。⛔ **AC#1..AC#6 一格不勾**（它们是实现格，实现腿还没动）。⚠ **三格判不了按未定义即停上报**（`cost-table.md` §8 的 C-1／C-2／C-3）：C-2 归票 42 那枚 `DEFERRED(watchdog loop/thresholds)`（`internal/watchdog/doc.go:18`，本票不吞）；C-1 与 C-3 **都是"乙才需要回答"的问题**（CLI 进程算不算常驻／CLI 与常驻实例能否并存），走甲之后**今天零行为后果**，故不摆给 owner，具名留在票里备查。⚠ 顺带两笔账落别处：`rt.close()` 不 join 的协程**分母 1→2**（见下面 AC#4 已改）；票 226 的"nothing polls"前提过期⇒ 补格落在**票 227 AC#6**。
- Status（**上一口径，保留不抹**）：**待派，排在队尾**（票 223 → `222-v1` → 221 → 224 → 220 → 213 → 214 → 219 那一小格 → **本票** → 227）。⛔ **它不是小活**：票 201 那四格（AC#1／AC#2／AC#3a／AC#6b）我之前写成"缺宿主原生入口"，**这说法被普查腿 `201-c2` 判为不完整**——真正缺的是"**一枚同时有审批门、有任务、又活着的进程**"，而今天这两条腿各缺另一半。批准与口径见台账 `A437`。
- 来源：只读普查腿 `201-c2`（`.scratch/wisp/probes/201/c2/census.md`，编排者代落盘＋逐条自己复跑尺）第 **①②③⑤⑦** 节。
- ⚠ **这不是安全漏洞，照这三行读**：①**现象在哪**＝产品界面（球／托盘）在真实运行时不出现，用户只能对着**控制台的文字**答复；②**有没有本机被入侵的证据**＝**没有**；③**最坏后果是什么形状**＝**该问的那一发在纯命令行里只能超时**（L2 超时＝自动拒绝），用户看不到那张本该点的卡。**功能与可用性缺口，不是被攻破。**

## 现量（起手逐条复算，别信这里的行号；票面行号已知会漂）

| 事实 | 读数 | 尺 |
|---|---|---|
| ⛔ **球只在旁支程序里** | `grep -rln "wisp/internal/ball" --include=*.go .`（剥 `.scratch`）＝ **只有 `cmd/balldebug/main.go` 一枚** | 编排者 09-29 14:1x **自己跑过同一把尺** |
| ⛔ **`wisp run` 跑完一发任务就退** | `cmd/wisp/main.go:85 os.Exit(cmdRun(args[1:]))` → `cmd/wisp/run.go:220 return rt.execute(task)` → `:855 res := bg.Wait()` → `:914 return`；`ctx` 本身就是一次性超时（`:834`），**没有 tick／队列／事件循环** | 两枚行号编排者复跑复认 |
| 常驻那条腿在，但**它没有门** | `cmd/wisp/resident_windows.go:24 runResident()` → `:33 proc.Boot` → `:83 rt.RunEventLoop()`；**这条腿没有 gate、没有 bridge、没有任务**（`approval_reply.go:48` 与票 201 §47 各自自陈一致） | 〔普查腿现读；`internal/proc` 那包**它自己没打开**＝见"先读它"那格 AC#0〕 |
| 球是**原生**的，接它**不需要新依赖** | `internal/ball/renderer_windows.go:4-9` 逐字：GDI 内存 DC 上的 top-down 32bpp DIB ＋ `ID2D1DCRenderTarget` ＋ `UpdateLayeredWindow`；窗口形 `ball_windows.go:196/212`＝`WS_EX_LAYERED\|TOPMOST\|TOOLWINDOW\|NOACTIVATE` ＋ `WS_POPUP`。**`go.mod` 里 webview 类依赖＝0 命中**（`go-webview2` 只活在旁支模块 `scripts/spike/go.mod:9`） | `grep -icE 'webview\|systray\|walk\|getlantern' go.mod`＝**0**（编排者复跑） |
| 托盘今天只有四枚、**没有任何一条路径通向答复** | `internal/ball/tray_windows.go:20-23`（id 1..4）＋`:86-91`（「打开面板」「静音」「暂停唤醒」「退出」）；点击→`ball_windows.go:664-681` 的 `Events.OnTray*`（宿主传的裸 `func()`，`:49-60`）⇒ **`internal/ball` 不 import 审批包** | 编排者 `sed` 复跑复认 |
| ⚠ **"托盘有个空 stub"这句是错的，别照它找落点** | ball 侧**没有 stub**，它已经把选择发出来了；那句 `fmt.Println("tray: open panel (stub, ticket 33)")` 在**消费者** `cmd/balldebug/main.go:199` | 普查腿顶回票 201 §13/§55 的口径，我复读原文复认（`census.md` 的 B-4 行） |
| 答复侧**已经全齐**（这是好消息，写清楚免得下一位重造） | `Gate.DecideFromNative`（`internal/agent/approval/gate.go:626`）／`DecideFromPanel`（`:638`）／`Veto`（`:387`）——**生产调用者 5 枚，全在 `replies.go:325/372/377/401/427`**；`Replies` 那一层由 `cmd/wisp/approval_reply.go:210/249/251/278/310` 调；入口枚数仍＝**1（控制台 stdin）** | 编排者复跑（多出的一行属 `Replies` 层，见 `census.md` 第 0 节末行） |
| 待答卡的数据**原生面直接读得到**，不用新造面 | `Replies.AwaitingHuman()`（`replies.go:249`）／`Pending()`（`:224`，按显示序）／`ReplyCard` 带 `CorrelationID/Level/Tool/Grant/Paths`（`:70-87`）；L2 另有 `Queue.LiveApprovals()`（`pending_read.go:104`，生产调用者 `cmd/wisp/panel_pump.go:62`） | 〔已证〕⚠ 但 **L1 的 `Gate.windows` 全仓零枚枚举口**（`grep -n 'g\.windows'` 只有 `:321/:324/:330/:393` 四处＝查重/写/删/按 corr 查，编排者复跑）⇒ 要枚举 L1 那一发就撞**票 220 的 AC#2 甲形**，先裁甲乙 |
| ⛔ **AC#2 那一格今天不可达，原因很具体** | `runSpec.replyVeto`（`cmd/wisp/run.go:146` 声明、`:601` 读）**生产赋值点＝0**；且 `run.go:446` 传的是**空参** `approval.NewChannels()` ⇒ 全部 `loaded=false`。⚠ **现成半成品就在场**：`internal/agent/approval/approval.go:163 DefaultChannels()` 逐字 `return NewChannels(ChannelBall, ChannelEsc)`，**生产零调用者**（`gate.go:108` 那个兜底只有 `Options.Channels == nil` 才走到）⇒ **不是"没建"，是装配处主动关掉了** | 编排者自己跑三把 grep 复认，并把"主动关掉"这句加重 |
| ⚠ **D38(c) 那句"必须等所有派生 goroutine 退出"今天已经违约** | `rt.close()`（`run.go:676-683`）只 `replyRoot.Cancel()`，`replyHandle`（`:265`，Spawn 于 `approval_reply.go:421`）**从不被等**；码内注释 `:672-675` 自陈「the goroutine dies with the process」 | 〔读到码＝已证；**"有没有仪器在查这条"未查**＝普查腿自陈，别当有仪器〕 |

## 冻结文字怎么说的（**只读，一字不改**；抄原文给下一位省一次翻文件）

- **D2（`docs/PLAN.md:73-99`）**：`:74`「**常驻原生进程树，空闲态无 WebView。**」；`:83-88` 把「悬浮球：原生分层窗口」「全局快捷键注册」「系统托盘」「Job Object 持有者」「KWS 唤醒词」画在**同一枚常驻主进程**里；`:103`「连锁推论：**Agent 核心循环必须在原生侧**」。⇒ **它没有一句话按"`wisp run` 这条 CLI 子命令"分别许可或禁止**——它是拓扑图。
- **D38（`docs/PLAN.md:2818-2848`）**：`:2824`「一个 STA/UI 线程｜拥有：分层窗口＋托盘＋WebView2 窗口＋所有 Win32 消息循环｜…**悬浮球绘制与面板必须在同一 STA 线程**；决定：**共用一个 STA 线程，不建第二个 D2D factory**」；`:2831-2833` 常驻 6 枚名额含 `ui-sta`（`internal/observe/goroutine.go:44` 在册）；`:2838`「总数上限：常驻 6＋每任务 3，超出即为泄漏征兆」；`:2845`「任务完成必须等所有派生 goroutine 退出（WaitGroup）才算完成」。
- ⚠ **两条都没写的东西＝本票的"未定义即停"点**：**上限是按"常驻"计名的**，而 `wisp run` 不是常驻进程 ⇒ **"在一发即退的 CLI 腿里起 `ui-sta` 会不会被看门狗判成泄漏"今天判不了**（`observe.Registry` 是否按进程角色分别计上限，普查腿没读到，`RosterReport.Unknown` 的判据也没读）。**这一格是 AC#0 的前置，不许按口味决定。**

## 为什么不许"就在 `wisp run` 里挂个球"交差了事（三支里我否掉的那两支）

1. **它要改退出语义**：球要活着，`execute` 后面就得挂等待（`run.go:220`）——那是把 CLI 腿改成常驻，**改的是 D2 的拓扑，不是加一行**。
2. **D38(c) 那枚违约还没修就先挂 GUI**：GUI 线程收口要 `WM_QUIT`（`sta_windows.go:47` 的 `start` 阻塞到它）、`Ball.Close()` 在 `ball_windows.go:904`；今天连一根答复 goroutine 都不 join ⇒ **进程退出时球窗口是怎么没的，没人说得清**。这一格必须**先修再接**。
3. **资源与 SLO 代价没量过**：D32 的 25MB/40MB 常驻口径与 handle 数那扇门，在"CLI 腿里多一枚 STA 线程＋一块 D2D 面"上是多少，**今天没有读数**（`cmd/balldebug/main.go:78-85 handleCount()` 是现成量具，但在 debug 件里）。

## 判据（草稿——每格都要 `file:line` 与正控；措辞由写腿自己按现跑读数重写）

- [x] **AC#0 先把那枚选择摆清（不许直接开写）**：读 `internal/proc`（`RunEventLoop` 到底能挂什么）＋`internal/observe/goroutine.go` 全文＋`RosterReport.Unknown` 判据，交回**两支的代价表**：**甲＝把审批门与任务通路接进常驻那条腿（贴 D2 的原意）**；**乙＝让 `wisp run` 在一发任务期间宿主球并正确收口**。**具名答三问**：常驻名额上限怎么算／看门狗会不会误报泄漏／谁的进程持有 Job Object。**未定案即停上报，不许按口味选。**
- [x] **AC#1 宿主存在性有生产点数**：改前 `internal/ball` 的生产引用者＝**1（`cmd/balldebug`）**，改后**必须有 ≥1 枚来自 `cmd/wisp` 那条真正会跑任务的腿**；尺＝同一条 grep，逐名报命中文件。⛔ 不许用"能编译"充当"有人起它"。
- [ ] **AC#2 托盘四枚按钮各映射到现成枚举、且有一枚真执行者**：「允许一次」→`Replies.Allow`（`replies.go:313`）；「拒绝」→`Replies.Reject`（`:336`）；「长期允许」→`replySurface.always`（`approval_always.go:70`，**它自带第二张 L2 卡**，不许绕过）；「本次会话内允许」⛔ **今天不许做**——`Answer` 词表只有 `allow/reject/veto/timeout` 四枚（`internal/tools/gate.go:75-89`，编排者复跑），会话档归**票 224**。判据形状＝每枚按钮断"点下去 ⇒ `Gate` 收到一条带正确 `Channel` 的答复"，并把回调接反/接空必须能判红。
- [ ] **AC#3 等待态真被驱动**：`bookWaitingState`（`cmd/wisp/approval_always.go:171-184`）**今天只写审计日志、不碰球**（我逐行读过）。本票要把球推进 D43 的"等人"那一态——⚠ **注意票 201 拆格口径**：AC#6a「态可计算可审计」**已成立、不许重做**；本票只做 AC#6b「球真进那一态」。正控＝**把驱动那一态的调用拿掉，判据必须红**（现在这一格是空的，任何突变都不会红——**这就是它必须被写成一格的理由**）。
- [ ] **AC#4 goroutine 收口先补旧账**：`rt.close()`（`run.go:676-683`）必须**等**它 spawn 的每一枚（含 `replyHandle`），且 GUI 腿走 `observe.Registry.Spawn`（`ball_windows.go:170` 那形，`ui-sta` 名额现成）。**具名写出：谁 join 谁、超时从哪来。**⚠ **分母已从 1 枚涨到 2 枚（09-29 由非实现者验收腿 `223-v1` 现测，台账 `A439`）**：除 `replyHandle` 之外，票 223 新起的那枚配置轮询协程 `reloadHandle`（`cmd/wisp/config_reload.go:119` Spawn，名字占用 D38b 常驻名册里的 `watchdog`）**全仓零读取点**，`rt.close()`（`run.go:694-707`，行号已从票面原记的 `:676-683` 漂下来）对它同样只 `Cancel` 不 join，注释逐字承认「The cancel is not a join」⇒ **同一个违约的第二枚实例**，本票要**两枚一起 join**、且判据得能证"没 join 就红"。⚠ 这一格不是新功能——它同时是 **D38(c)（`PLAN.md:2845`）那条违约的自救**。
- [ ] **AC#5 面板侧来源的 L2「允许」仍必须被判红**：`AGENTS.md` §1.2 那条硬禁不因"加了中国原生入口"而松（`Q-49` 丙那批判据）。正控＝从界面方法名送一发"允许"，路由必须拒；⛔ **不许为了便利把 `DecideFromPanel` 改宽**。
- [ ] **AC#6 整包终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1` 到终态＋逐名比红名册（历史在册 `internal/ball` 1＋`internal/panel` 4 属别人地界，不算本票新增、不许顺手修；`internal/risk TestResolvePerCallBudget` 是争用型假红，安静 `-count=3` 为准，见 `A432`／票 226 AC#6 那格）。**⚠ 本票动的是 Windows 原生面：POSIX 半边只能在全仓仪器里量到，别拿"编译过"当"跑过"。**

## 禁区（写腿动手前逐条读，含三枚按符号名办事的雷）

- ⛔ **禁跑在别的写腿在飞时**：与**票 220（同文件 `internal/agent/approval/gate.go`）**、**票 223／224／214（同文件 `cmd/wisp/run.go`）**、**票 226（同文件 `cmd/wisp/approval_always.go`＋`internal/config/allowdirs.go`）**同文件互斥；与票 221／222 包级串行（`cmd/wisp`）。**不撞的只有票 213。**
- ⛔ **`tools/d22scan` 三条会天然打中本票**（现读 `tools/d22scan/main.go`，全部现跑复认）：**ban #1**（`:697-711`）射程＝**任何 `go <anything>`**，具名调用也算，唯一按文件豁免＝`internal/observe/goroutine.go`（`allowlist.txt:7` 逐字「Exempted BY FILE PATH only (never by call shape)」）⇒ 宿主 goroutine 必须走 `observe.Registry.Spawn`；**ban #4**（`:144-146`＋判定循环 `:758-773`）只认 `\.Sub\(time\.Now\(\)\)` 与"同一行 `time.Now().Unix*()`＋`timeout|deadline|expire|ttl|budget|until`"，⚠ **它只跳过整行 `//` 注释、行尾注释仍算进该行** ⇒ 做"卡片剩余秒数"必须用 `context.WithDeadline`／monotonic `Duration`，**别写减法**；**ban #8**（`emojiRe` 整行在 `:164`）**实际扫的 6 段码位＝`U+1F000–1FAFF`／`U+2200–22FF`／`U+2600–27BF`／`U+2B00–2BFF`／`U+FE0F`／`U+1F1E6–1F1FF`** ⇒ **⚠/✅/⛔/✓/✗ 放进任何非注释文本必红**，而**托盘菜单标签是字符串字面量＝在射程内**（`tray_windows.go:86-91` 那种 `appendItem(id, "打开面板", false)`）；`← ⇒ ≥ — 「」§` 与全部 CJK 安全。
- ⛔ **命名与落点**：新增可解码结构体**不要放 `internal/panel`**（`l2_grant_boundary_test.go` 是冻结件，且它的 12 词根×两拼禁名词表起于 `:1881`——⚠ 票 219 写的 `:1882` 差一行；"它只扫 `internal/panel`"这句**只有注释级自证、未行级证**，保守处理）；答复字段**只叫 `reason`**；不新增 C17 方法名（新增要先落 `A##`）。
- ⛔ 不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；三枚冻结件（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`ticket90_persist_test.go`）一字不动；**`frontend/**`／`design/**` 零写、连内容都不转述**（界面侧那三枚按钮属 owner 自己带给他的另一枚 agent，本仓只交 Go 侧宿主与状态生产者）。
- ⛔ **不许顺手把"会话内允许"做出来**（那是票 224，且它的"什么算同一次会话"已由 `A435` 定死＝宿主随机铸造、不许派生、结束点今天＝进程退出）；**不许把 `allowed_dirs` 做成执行时硬边界**（它是判级输入，这条已有裁定）。

## 与其它票的关系（别在本票里顺手做）

- **票 201**：本票是它 AC#1／AC#3a／AC#6b 三格的**唯一前置**；**AC#2 的那半格可以早做**——只差把 `runSpec.replyVeto`（`run.go:146`）设上、并按 `approval_reply.go:415-417` 那道 `SetLoaded` 的门，三件现成品都在（`approval.go:147`／`:163`、`gate.go:139`），**缺的只是那一枚赋值**。要不要单独拆一小腿先做，由排程决定，但**它撞 `run.go` ⇒ 仍须排在票 223 之后**。
  ⛔ **上面那句"缺的只是那一枚赋值"是我（编排者）自己写错的，09-30 15:4x 现读推翻，原句不抹、账 `A472`**：`cmd/wisp/run.go:141-147` 逐字写着「Production leaves it empty because a console run wires **none** of SPEC-06 §2's four channels, and an empty value makes `Gate.Veto` say that back **instead of letting this process claim a cancel path it does not have**」＋「A host that DOES own one - the ball click, the global Esc hook - names it here **and marks it loaded in the same step**」。
  **两把现量尺**（我现跑）：① `grep -rn "ChannelEsc" --include=*.go internal cmd \| grep -v _test.go` ⇒ 命中**全在 `internal/agent/approval/` 自己的定义里**（`approval.go:55/:61/:90/:111/:163`），**没有任何宿主把它设上**；② 全局热键那套东西**在场但只在球包里**——`internal/ball/hotkey_windows.go:434 takeEsc`（逐字 "binds VK_ESCAPE as the cancel hotkey (B1 Confirming takeover)"）、`:374 registerAll`，而**唯一把它跑起来的宿主是旁支程序 `cmd/balldebug`**（`:224`、`:255`）。
  ⇒ **结论（这一处是"照我写的做会造出一枚谎"）**：在 `wisp run` 里把 `replyVeto` 设成 `ChannelEsc`／把 `NewChannels()` 换成 `DefaultChannels()` ＝ **让一条纯命令行进程对外声称自己能按 Esc 取消、并往卡片上印"按 Esc 键"**，而那形里既没有球也没有热键钩子。`run.go:501` 那段注释把这枚留空写成 **load-bearing 的诚实选择**，不是漏了赋值。
  ⇒ **改后的落点**：这一半**必须在"球进常驻那条腿"之后才有对象可声明**（宿主先真的有 ball click／`takeEsc` 那条路，再按 `run.go:147` 那句"names it here and marks it loaded in the same step"同一步设上＋`SetLoaded`）。⛔ **不许再有人拿"只差一行赋值"来拆这一小片**；票 201 `AC#2` 那半格的真实前置＝**本票的载体主体**。
- **票 220**：谁先造"L1 那一发的只读枚举口"要先裁（它 AC#2 甲形与本票第⑤问同一枚位置）——**两腿各造一次＝两份枚举器**。
- **票 227**：配置写路径的覆盖面补全，排在**本票之后**（本票会真去写 `allowed_dirs`，那一天的基线才作数）。

## 编排者增量 — 09-30 **14:47:10**：接一枚从票 241 挖出来的**过期落点指认**（本票因此多一格）

- [x] **AC#7 「resident_windows.go:81」那行注释今天是假话，且它已经被三程原样传递**：现读逐字「wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)」——**票 07 早已结案**（盘上文件「07-ball-state-machine-core-done.md」），"要等票 07 才来"这句把本票（228）该干的活指给了一个不存在的未来。传递链：该注释 → 普查「240-c1 §3.1 第 1 行」→ 「241-r1」票面正文两处归属 → 均由「241-v1」当场挖出并上交（它按硬约束没碰「cmd/wisp」）。**完成判据**＝① 注释改成带条件的事实句（说清"今天这行进程里确实没有球，把球装进来是本票的活"）；② 票面正文那两处「欠票 07」的作废声明由我 09-30 在票 241 落笔（已落）；③ ⛔ **不新增仪器去扫注释里的票号**——词面型尺会误伤，这条债的归属靠台账与票面，不靠新门。凭据出处：「docs/evidence/s1/241-audio-level-producer-v1.md」§2.3。
- [ ] **AC#7 的射程更正（09-30 15:3x，账 `A471`；这条是上面那一格的正文补，不是第二格）**：本格原来只写「resident_windows.go:81」**那一行，太窄**。我现跑 `grep -rn "ticket 07" --include=*.go internal cmd scripts tools`＝**28 处命中**，按"谁读得到"分两堆，两堆处理方式不一样：
  ① **3 处不在注释里，是打给人看／写进日志的字符串** ⇒ **这三处才是本格的完成判据**（逐处改成带条件的事实句，说清"今天这条进程里确实没有球，把球装进来是 228 的活"）：
  「cmd/wisp/main.go:25」⚠ **它落在 `const usage = ` 那枚反引号原始串里＝`wisp -h` 打给用户的正文**（逐字「the floating ball window is ticket 07」，我 `sed -n '20,30p'` 复认），**不是注释、用户今天真读得到**；
  「cmd/wisp/resident_windows.go:81」是 `fmt.Printf` 的运行时输出（本格起手记的那一处）；
  「internal/proc/boot_windows.go:127」是 `slog.Info` 的运行时日志（逐字「activation requested by second launch (ball bring-to-front lands with ticket 07)」）。
  ⚠ 第三处在 `internal/proc`、跨了本票地界：改的是一枚字符串、**不新增依赖边**，所以本票可以做，但交件时必须具名写清"动了 `internal/proc` 那一行"。
  ② 其余 **25 处是注释**（其中 **6 处在 `*_test.go`**）——**大部分是合法的历史出处**（「Implemented by ticket 07」＝过去时、「SPEC-08 §2, ticket 07 constraint」＝引来源），少数是**还在承诺未来**。逐处判"历史出处 vs 过期承诺"⛔ **不塞进本格**（那会让 AC#7 永远勾不上），另立**票 243**（只读普查，零产码）。
  ⛔ 原判据第 ③ 条**继续有效、且对两堆都成立**：**不新增扫注释里票号的词面型仪器**——上面那把 `grep` 是**派单与台账用的尺，不是门**。

## 编排者收表 — 09-30 **17:0x**（账 `A477`）：非实现者裁决表 `228-v1` 交回 ⇒ **AC#1／AC#7 翻勾**，另补两格它挖出来的账

裁决表：`docs/evidence/s1/228-resident-ball-v1.md`（**346 行／52,643 字节、真占位符 0 枚**，我现跑核；三枚 commit「48c15f71 骨架 → caa00c6a 全读数 → 8dc60090 补正」逐枚 `git log` 复认存在）。验收者≠实现者：实现是 `228-r1`，本表出自另一程。

- [x] **AC#1 翻勾凭据**：判语**成立**。承重证据不是"`internal/ball` 被 import 了"（那正是我派单里点名的假绿形状），而是**定向突变 M1**：摘掉生产调用那一行 ⇒ `TestAC228ResidentLegIsTheBallHost` 与 `TestAC228ResidentLegReportsAndBooksItsBall` **双双点名红**；另有真机七发拉起无参数 `wisp.exe`（hwnd／pid 在册、hotkeys live 4/4、退出后窗口不在）。**我自己复跑的两把尺**：`grep -n "startResidentBall\|runResident(" cmd/wisp/*.go | grep -v _test` ⇒ 生产链完整（`main.go:64` 无参数分支 → `resident_windows.go:27 runResident` → `:118 rb := startResidentBall(rt.Registry)`）；`grep -rln "internal/ball" cmd/wisp/*.go | grep -v _test` ⇒ 三枚产码文件在。**一枚只有突变看得见的精度**：nil 的 `statusLine()` 照样含"NO floating ball window"，所以控制台那句不红、红在磁盘记账那一支——本票的钉里"打印面"比"记账面"松，这一条我登记为已知缺口（归 228 后续片，不追改）。
- [x] **AC#7 翻勾凭据**：判语**成立**，且**只在 ①那三枚人读字符串这一档射程内成立**（三处逐处现读均改成带条件事实句：`main.go:24-29`／`resident_windows.go:131`／`internal/proc/boot_windows.go:132-134`，boot 分支那句拿真机 console 原文证的）。⚠ **本格翻勾不等于票 07 那批注释都清了**——见下面 AC#8。

- [ ] **AC#8（新补，`228-v1` 的 F-1 ⇒ 归属由我裁）`cmd/wisp/main.go:8` 那句 godoc 注释今天还是假的**：现读逐字「// frozen D38(e) 10-step shutdown order. **The floating ball GUI is ticket 07.**」——球已经在同一条进程里被装起来了，这句还在承诺一个已经结案的人。**完成判据**＝改成带条件的事实句（说清"球由本文件所在的常驻腿装配，见 `resident_ball_windows.go`"），⛔ 不许新增扫注释票号的词面型仪器。**残缺表现**＝读代码的人（含下一程腿）按这句去找"票 07 还没干的活"，会白跑或把债记到别人名下。
  **归属怎么定的（这是我的裁定，不是盘上现成的答案）**：`228-v1` 报"F-1 与 AC#7 ② 的归属互斥、未裁"。我判**归本票**，三条现读凭据：① 文件在 228 地界（`cmd/wisp`）；② 是**本票这一程把那句变成假话的**（把主张说反的人负责改它）；③ 票 243 普查表 `census.md:75` 第 24 行的归属列**逐字写的就是「票 228」**，与 ② 一致 ⇒ 不是互斥，是我先前那句"② 全归 243"写得太粗。**票 243 保持结案不动**（它的活是普查＋指认，那部分真交付了）。
- [ ] **AC#9（新补，`228-v1` 的 F-6）球起不来时那句实话会被说两遍**：`resident_ball_windows.go:96` 与 `internal/proc` 的 boot 报告各说一次，而成功支不说第二遍 ⇒ 控制台"恰好说 1 次"那条不变式被强推的失败支撞红。⚠ **失败支今天从未被自然走到**（`startResidentBall` 把 `ball.New` 写死、无注入面），所以这条**只有突变能现形**、CI 里零读数（`winlive` 在 `ci.yml` 零命中）。**完成判据**＝两处只留一处作者（另一处改成引用或删句）＋给失败支一枚能被自然走到的注入面或台件。

- [ ] **AC#10（新补，10-01 09:5x 编排者追加，来路＝`246-v2` §5 第 4 条 ⇒ 账 `A488`）票 246 AC#7 之后，球那条腿的三处"我没有任务管线"自述里有两处变成假话，其中一处是运行期真打进台账的句子**：`cmd/wisp/resident_ball_windows.go:21`（注释）与 **`:160`（真打印的 slog 字段）** 逐字「this leg has no task pipeline and no microphone」，`:203` 的 `ballGestureWhy` 常量逐字「this process has no task pipeline, no microphone and no panel host」。⚠ 关键点不是"注释过期"，是**打印时刻（`:136` 球起来）确实早于装配时刻（`resident_windows.go:177` 任务源接上）**，所以那一行在它打印的那一瞬间不算谎——**但它声明的是属性而非此刻**，读台账的人（含 owner 自己看日志）会当成永久事实，而票 246 那 6 发真机的卡片就是从这条腿举起来的。**完成判据**＝三处逐处改成带条件的事实句（说清"此刻还没有任务源，任务源由同进程 `resident_windows.go:177` 之后装配"），⛔ 不许新增扫注释票号的词面型仪器（与本票 AC#8 同规矩）；⛔ 不许靠"把打印挪到装配之后"来糊——那会破 D38(e) 的时序，属另一枚契约面。
  **归属（这是我的裁定）**：本票是"球进常驻那条腿"的真相源，这两句由**票 246 那一程说假**、但改口责任在球这一侧（谁的话谁负责），且落地时`cmd/wisp`同包只许一枚写腿 ⇒ **排在本票 AC#2/AC#3 之后、与 AC#8 同批改**，⛔ 不与票 33／248 的写腿并发。

- [ ] **AC#2 的射程被票 244 的 `244-a2` 普查撑宽了一格（10-01 10:3x 编排者追加，账 `A490`；这一格**并进 AC#2**，不新开 AC）**：托盘那四枚按钮里 **"退出"这一枚是 GUI 进程未来的唯一停机源**。现量（我 10:2x 亲自复认过两处）：唯一真停机入口是 `cmd/wisp/resident_windows.go:107` 的 `os.Interrupt`／`SIGTERM`，而 `cmd/wisp/resident_ball_windows.go` 的 `recordTrayExit` 只打 `outcome=ignored`。⇒ **票 244 切 `-H=windowsgui` 之前，AC#2 那一枚"真执行者"必须先落地**（否则双击起来的常驻进程只能靠任务管理器杀、D38(e) 十步一步都走不到）。⚠ 那段注释自称两条补法都是新契约（`internal/proc` 的 stop-request 钩子／第 7 枚常驻协程受 `observe.ResidentNames` 六枚冻结名册限制）——**用哪条形由我在 `244-a3` 读数回来后裁，⛔ 实现腿不许自己挑**；钩子形若可行，可复用 `A488` K7 那条已裁灰区（"往 `RegisterShutdownHook` 注册入口不算契约面"），协程形则要人工批准（动 D38(b) 名册）。

- [ ] **AC#11（新补，10-01 10:5x 编排者追加，来路＝`244-a3` §② 乙格主表＋我自己的裁，账 `A493`）「退出」那枚托盘按钮的真执行者只能长成一个形状：记一笔请求、立刻返回，由常驻腿自己去结束环路。** 底座读数（`244-a3` 现量，我逐条复读过原文）：常驻腿的**有序退出已经全部挂在 `runResident` 的正常返回上**（`cmd/wisp/resident_windows.go:72-82` 的 `defer rt.Shutdown(false)` 走十步＋`:137 defer rb.stop()`＋`:144 defer ra.detachBall()`＋`:70 defer sink.close()`）⇒ **要做的只有一件事：让 `rt.RunEventLoop()`（`:202`）返回**，十步一字不用改。
  ⛔ **反面判据（谁做成这样谁就是错的，不必等验收腿发现）**：在托盘回调里直接 `rt.Shutdown(…)`＋`os.Exit(…)`＝**跳过 `os.Exit` 之前的全部 defer**——十步一步不走、日志不刷、球窗口留着。⚠ 另一枚必读机制：托盘回调跑在 **`ui-sta` 线程内部**（`internal/ball/ball_windows.go:669-686` 的 wndproc），而 `Ball.Close()` 是 `sta.PostTask(…) + <-done`（`:904-935`）⇒ **在回调里同步收口＝自等待死锁**。合法形状只有"记一笔请求、立刻返回"，而这正是 `vetoByEsc` 已有的形状（`cmd/wisp/resident_approval_windows.go:168-170` 逐字「It runs ON the ui-sta thread…it waits on nothing」）。
  **修法＝乙-1（`internal/proc` 的 stop-request 位），⛔ 不是协程形**。三条读数定死这件事：① `RunEventLoop` 的 `sigCtx` 是**函数内局部变量**（`boot_windows.go:128`）且 `Runtime` 今天没有任何 stop 位（`grep -rn "RequestStop\|StopRequest" internal/proc/`＝零命中）⇒ **第 7 枚协程那一形单独完不成停机**（它要么仍然调一枚 proc 侧方法＝回到乙-1，要么跳过 defer，要么在 ui-sta 上死锁）；② 动名册那一形（乙-2c）**在本票射程内做不到自洽**——它要改 `docs/PLAN.md:2831` 的冻结正文（＝人工批准），还要为让 SLO 那句继续说真话去改 `internal/observe/thresholds.go:36`，而那枚文件是**票 244 的禁区**；③ 现有两枚真机台件（`resident_ball_228_windows_test.go:184`／`resident_sink_nail_127_windows_test.go:565）**只会打"控制台信号"那一支**（`buildWispForTest` 不带 `-ldflags` ⇒ 台件二进制永远 CUI）⇒ 托盘那一支落地后**它们不会自动变红也不会自动变绿**。
  **本格的三条完成判据（缺一格不算落地）**：ⓐ 停机请求位**必须在 boot 时建好**，并**同时接进 `RunEventLoop` 之前那一段空窗**（`resident_windows.go:106-108 signal.Notify(bootExit)` 与 `:195-203` 那个 `select`）——⚠ 这是本形唯一隐藏的坑：只接进环路的话，球起来之后、环路进入之前按下"退出"会**被吞掉**，而 `resident_ball_228_windows_test.go:172-215` 那枚用例的存在理由就是这段空窗；ⓑ **必须新加一枚会响的仪器盯"托盘那一支触发的退出也走满十步＋审计顺序仍是 1..10"**（现成可照的形状＝`resident_approval_246_windows_test.go:132`／`:184` 那两枚"钩子真挂上序列、出货进程真持有那一步"），⛔ 不许把现有控制台台件当这半格的凭据（那是 `A490` G5 那句"看着验过、其实没验这一维"的重演）；ⓒ `resident_ball_windows.go:246-252` 那段"两形并列都能给执行者"的注释**是不准确的**（见①）⇒ 落地时随 AC#10 一起改口，具名说"协程形单独完不成停机"。
  **契约面这一格我自己裁（这是我的推广，不是盘上现成的答案）**：`A481`／`A488` K7 那句"往十步里挂**注册入口**不算契约面"罩的是**钩子注册**；本形动的是 `RunEventLoop` 的 doc 里被**逐字枚举过**的退出触发名册（`boot_windows.go:122-126`「Ctrl+C / console close / WM_ENDSESSION landing in ticket 43」）。我判：**加一枚触发同样不改那十步的顺序**，与 K7 同理 ⇒ **不算契约变更、不摆 owner**，但落地那发必须①把那句 doc 改口成**完整**触发名册（⛔ 不许留过期句子），②本条以一枚具名 `A##` 记录在案（＝`A493`），③若实现时发现必须动 `hookableRoster` 那**闭集八枚**（`shutdown_hooks.go:47-56`）才能做成——那就从"加触发"滑成"改契约"，**当场停手上报**，不许顺手加第九枚槽。

  ⚠⚠ **判据 ⓑ 我自己写空了，就地更正（10-01 11:3x，来路＝只读普查 `228-a2` §⑦ F5＋§③ D-3，账 `A495`；上面 ⓑ 那句原话不抹）**：`228-a2` 量到 **"审计顺序仍是 1..10"这一维在真机盘面上结构性读不到**——`internal/proc/shutdown.go` 只为 **skipped／失败／fast** 出声（`:130`/`:145`/`:147`/`:172` 四支出声，`:178-180` 与 `:186` **零日志**），所以成功的那几步在台账里根本没有行；而 `stepsOf127` 那组读数今天只进 `t.Logf`（`cmd/wisp/resident_sink_nail_127_windows_test.go:640`）。⇒ 我原来那句"一枚仪器盯『托盘那一支触发』＋『顺序还是 1..10』"**把两维压在了一枚测不到的钉上**，照它做只有一条出路＝给 `internal/proc` 每一步补一行日志，而那一支会同时改掉 127 一族对**记录集合**敏感的断言（`:619-641` 那组"installIdx==1／最后一条必须是 shutdown 记录"）＝**动 D38(e) 的观测面**，本票不许。**ⓑ 改成两枚分开的钉**：
  - **ⓑ-1（进程内，不需要窗口）**＝照抄现成形状 `cmd/wisp/resident_approval_246_windows_test.go:145-158`（真 `proc.Boot`＋真 `RegisterShutdownHook`＋`rt.Shutdown(false)` 之后**在进程内读满 10 条 `[]StepRecord` 并逐名断顺序**）。它钉的是"**触发之后十步真走满且顺序对**"这一维。⛔ 不许用它去声称"托盘那一次点击被验过"（它没有窗口腿）。
  - **ⓑ-2（真机，只钉可观测的那半）**＝离它最近的可照形是普查挖出的第四枚成员 `cmd/wisp/resident_ball_live_228_windows_test.go:89`（真机进程×跨进程真窗句柄，带 `windows && winlive`）。这一枚只许钉"**那一次触发导致进程以退出码 0 结束＋审计里仍有 install 记录与 D38(e) 尾巴**"，⛔ **不许写成"读到 1..10 全序"**（见上，结构性读不到；写上去就是逼下一程去动日志面）。
  - ⛔ **一条被证伪的取巧路就地关掉**：「从测试里投一枚 `WM_COMMAND`＋菜单 id 去模拟托盘选择」**不成立**——`internal/ball/tray_windows.go:15-16` 明写**不用** `WM_COMMAND`，选择值只从 `TrackPopupMenu` 的**同步返回**拿（`:101-106`）。⇒ 任何腿不许再造这一形；要验"执行者被挂上"就走 ⓑ-1 那一形（在装配处拿那枚被注入的函数值直接调），要验"真点一次"就只许走 ⓑ-2 那枚 `winlive`。
  - **附带一条我自己复算出来的更正（这条普查腿交回的一句我核后打折）**：它写「第七枚常驻协程 ⇒ `boot_windows.go:115-118` 直接拒 Boot」。**只在挪用**在册 resident 名字**时**成立：`internal/observe/goroutine.go:422` 是 `rep.Resident > ResidentBaseline`，而 `Resident` 只累加 `Category == CategoryResident` 的条目（`:409-411`），`ClassifyGoroutine`（`:64-67`）判"在册"用的是 `slices.Contains(ResidentNames, name)` ⇒ **名册外的名字落 `rep.Unknown`（`:419`）、不进 `Resident`、碰不到这扇门禁**（这与 `244-a3` R25／`A493` J6 的读数一致）。⛔ 下一程不许拿"会被拒 Boot"去排除或论证任何协程形——那一维的真实形状是"名册外不拦，只留一行 WARN＋一条 `Unknown`，代价见 `A493` J6 乙"。
  - **F6 那一格我答**：复用现成的 `OnTrayExit`（不新增 `Events` 键）是**可接受**的，因为那枚名册钉的措辞是"每个 gesture callback 都要有执行者"（`cmd/wisp/resident_ball_228_test.go:315-318`），**不是**"每个执行者都要有行为"；⇒ 名册不响**不等于**有行为。行为那一维**另钉**＝上面的 ⓑ-1。⛔ 不许为了"让钉响"去扩那枚名册断言到行为维——那是改别人已勾判据的形状（本票 AC#1 已由 `228-v1` 判过）。

**其余判语我照表收（不重裁，那是下一枚验收腿的活）**：② 失败支进程选型"响亮但不致命"**成立**（与本仓票 117 同腿先例／票 128 射程／不堵死联调三条相合，不判应拒绝启动）；③ `bootExit` 第二发是**被丢不是被打死**（甲形 5/5 exit 0、0xc000013a 零枚；乙形也不重复触发）＝F-4 具名未修；④ 0xc000013a **亲手复现两枚**且票 127 那枚钉**断言一字未放宽**；⑥ "零新增模块依赖"成立（`go.mod`/`go.sum` 四枚间 diff 空），**"零新增依赖边"要分开读**＝`cmd/wisp` 生产 import 集 16→17、新增恰好只有 `internal/ball` 一条（这正是 AC#1 的尺本身），唯一新协程在册名 `ui-sta` 走 `Registry.Spawn`。
**它自己列的判不动／未做（我不替它补）**：winlive 只跑了那一枚未跑整包；未量句柄／RSS（235ms 只当空窗旁证，**不入性能账**）；F-6 在 CI 会不会响零读数；球的十枚手势一次没真按过（**不替 AC#2 背书**）；GitHub runner 有没有能建窗口的桌面＝判不动。
**本票状态**：AC#1／AC#7 已勾，AC#2..AC#6 与新增 AC#8／AC#9 未勾 ⇒ **不加 `-done`**。队列不变：本票后续片排在 **245-r1 → 197-r3 → 224-r3 → 242 → 244** 之后。

## Progress log（本小节由普查腿 `228-a2` 于本轮新增；此前票面无此节。⛔ 未改动任何 AC 复选框）

- `228-a2`（只读普查腿，HEAD `0c8b9fddd2c3da647a76f0757fd466ef6767e6c1`，branch dev，起手尺 `date` 逐字读数见 census §⑤ R1）：
  交付 `.scratch/wisp/probes/228/a2/census.md`。射程＝"托盘那一次点击导致十步有序退出"这一维今天能被哪几种仪器形状看见。
  三条承重现量：① `internal/ball/tray_windows.go:101-106` 用 `TrackPopupMenu(TPM_RETURNCMD)` 的**同步返回值**取选择，
  注释 `:15-16` 明写**不走 `WM_COMMAND`** ⇒ "PostMessage 一枚 `WM_COMMAND`＋id=4"这一支**在今天的码里不存在投递路径**；
  ② 全仓 `grep -n 'shutdownRequest|stopRequest'` **零命中**，`RunEventLoop` 唯一定义 `internal/proc/boot_windows.go:127`，
  其唯一输入 `:128` 的 `signal.NotifyContext(os.Interrupt, SIGTERM)` ⇒ AC#11 要加的位是**全新面**；
  ③ `internal/ball/interaction_live_test.go:12` 逐字承认 winlive 族测不了托盘菜单（"needs a real mouse in the notification area"）。
  十步审计侧结论：新加触发**不会**让 `internal/proc/shutdown_test.go:11` 自动变红（它手搓 `ShutdownHooks{}`，`shutdown.go:112` 签名不受影响）；
  但真机侧今天**无任何台件读满 10 条 `StepRecord`**（只看 `residentShutdownRecord` 存在与相对位置）。
  判不动的 5 格已按 甲／乙／不做 交裁：census §⑦ F1（托盘点击可注入面）、F2（哪一族能完整覆盖）、F3（空窗坑最小可测形状）、
  F4（`buildWispForTest` 是否产 `-H=windowsgui` 第二产物，本腿未现读）、F5（要不要新增跨进程读 10 条记录的仪器）。
- `228-a2` 第二发（终态）：`①②③④` 正文补完，并**自查出并改正了第一发的两处假读数**（详见 census §⑥ E4/E5/E11）。
  最重要的三格新读数：① 全仓 `*_test.go` 里 `wmAppTray`/`showMenu`/`menuExit`/`recordTrayExit` **零命中**（尺 R13）⇒ 托盘命令分流今天**四族台件皆未碰**；
  ② `internal/proc` 的日志面**只为 skipped／失败／fast 出声**（`shutdown.go:130/145/147/172`），成功步与第 8/10 步零日志（尺 R19）
  ⇒ **"审计记录里 1..10 都在"在真机盘面上结构性读不到**，不是缺台件；进程内唯一读满 10 条 `StepRecord` 的现成形状是
  `cmd/wisp/resident_approval_246_windows_test.go:145-158`（尺 R23）。
  ③ 普查挖出任务描述之外的第四枚成员 `cmd/wisp/resident_ball_live_228_windows_test.go:89`（真机进程 × 跨进程真窗句柄 × 要求窗灭），
  它是现存码里离 AC#11 ⓑ 最近的一枚；其触发器 `breakToLoop`（`resident_sink_nail_127_windows_test.go:253`）仍是控制台事件。
  会自动变红的两枚（其余不会，凭据 census §④ D-1）：第七枚常驻协程 → `internal/proc/boot_windows.go:115-118` **Boot 直接失败**；
  新增 `internal/ball.Events` 键 → `cmd/wisp/resident_ball_228_test.go:319-322`。
  交裁 7 格：census §⑦ F1（托盘命令可确定性送达）／F2（拒绝猜修法）／F3（空窗最小形状甲乙）／F4（`-H=windowsgui` 第二产物）／
  F5（真机侧十步读面）／F6（族② 沉默是否算缺陷）／F7（winlive 两包桌面互斥）。本腿零产码零构建，未 push。
