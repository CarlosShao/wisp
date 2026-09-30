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

**其余判语我照表收（不重裁，那是下一枚验收腿的活）**：② 失败支进程选型"响亮但不致命"**成立**（与本仓票 117 同腿先例／票 128 射程／不堵死联调三条相合，不判应拒绝启动）；③ `bootExit` 第二发是**被丢不是被打死**（甲形 5/5 exit 0、0xc000013a 零枚；乙形也不重复触发）＝F-4 具名未修；④ 0xc000013a **亲手复现两枚**且票 127 那枚钉**断言一字未放宽**；⑥ "零新增模块依赖"成立（`go.mod`/`go.sum` 四枚间 diff 空），**"零新增依赖边"要分开读**＝`cmd/wisp` 生产 import 集 16→17、新增恰好只有 `internal/ball` 一条（这正是 AC#1 的尺本身），唯一新协程在册名 `ui-sta` 走 `Registry.Spawn`。
**它自己列的判不动／未做（我不替它补）**：winlive 只跑了那一枚未跑整包；未量句柄／RSS（235ms 只当空窗旁证，**不入性能账**）；F-6 在 CI 会不会响零读数；球的十枚手势一次没真按过（**不替 AC#2 背书**）；GitHub runner 有没有能建窗口的桌面＝判不动。
**本票状态**：AC#1／AC#7 已勾，AC#2..AC#6 与新增 AC#8／AC#9 未勾 ⇒ **不加 `-done`**。队列不变：本票后续片排在 **245-r1 → 197-r3 → 224-r3 → 242 → 244** 之后。
