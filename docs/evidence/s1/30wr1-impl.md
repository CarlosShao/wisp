# 票 30 装配缝腿 `30-w-r1` —— D10 结果呈现路由的常驻装配（写腿）

> **尺**：`HEAD=f6928ea01`（`dev`，现跑 2026-10-31 15:2x–16:3x +08）。★开工三核：①分支＝`dev`、HEAD＝`f6928ea01`（派单锚点即此提交，无偏移）；②索引脏面＝**他人 staged 的 `.scratch/wisp/issues/28-session-scope-warm-done.md` 一枚**＋散在 untracked（都不在本腿写面；本腿全程显式 pathspec，未触碰）；③`git status --porcelain -- cmd/wisp` 开工前**空** ✓。★环境：staging 全套按 `381-v3-verdict.md` §1 清单（staged 测试二进制＋三 DLL 同目录＋`WISP_SILERO_MODEL/SPIKE_MODELS/PUNC_MODEL/TTS_DIR/KWS_DIR/15E_TMP`＋CWD＝树根）；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`；rc 一律重定向现取；`unset GOTMPDIR` 每跑；变异全走**仓外 mutant ＋ `-overlay`**；⛔ push/amend/reset/rebase/stash/checkout。
> **写面**：`cmd/wisp` **新增 3 枚／33,385 B**（`resident_present_windows.go` 10,377 B 生产装配＋适配器＋真剪贴板写半边；`present_wiring_30w1_windows_test.go` 20,550 B E2E 读数件；`present_clipboard_30w1_live_windows_test.go` 2,458 B winlive 往返件）＋**改 2 枚**（`run.go`：import＋`presenter` 字段＋execute() 调用点与旧路的和解门，59+/20− 合计；`resident_windows.go`：`raudio.attachVoice(rv)` 之后的**兄弟装配段** 13 行）＋本证据件。**⛔ 零字节**：`internal/present/**`（19 PASS 不动）、`internal/agent/**`、`internal/risk/**`、`internal/panel/**`（C17 桥、panel.* 白名单、Snapshot 四键零碰）、`internal/ball/**`、**43-r1/384-r1 的挂接行**（逐行未动，兄弟段插入其前后空隙）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、specs、PLAN。零配置键、零 thresholds、零 golden、零 `DEFERRED` 标记、零裸 `go func(`、零墙钟超时（`observe.NewTimeout` 一枚，测试内 readiness 等待）。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\30wr1\{logs,mut,bin,export,tmp15e,overlay-m1.json,overlay-m2.json}`。

## 0 一句话

票 30 收口残余③执行：`present.New(present.Channels{...})` **一次装配**进 `runResident`（`assembleResidentPresenter`，`resident_present_windows.go`），把票 30 落的库（生产调用者 0 → 本腿后 1）接进常驻进程；调用点＝`execute()` 里 **loop.Run 返回后对 `res.Status==completed` 的结果** `Dispatch(ctx, present.Input{TaskID: res.TaskID, Text: res.Text})`（run.go，30-v1 §9.3 具名交回逐字）。旧 D10 通知＋wire-r1 播报路在 presenter 已装配的宿主上**让位**给路由器（否则 completed 结果双 toast＋双播报、长结果发通知违背 D10 档表）；presenter 为 nil 的宿主（＝CLI `wisp run`）**字节同行为**。

## 1 ★★★ 装配逐通道（30-r1 §8① 交回清单 ↔ 实装对照）

| 通道 | 交回原话 | 实装 | 真对象 |
|---|---|---|---|
| Speak | `rt.voiceAnnounce` 同款闭包、⛔ 重复设闸 | `ch.Speak = presenterSpeaker(run.voiceAnnounce)`——**就是 wire-r1 在 resident_windows.go:580 接的 `rv.announce` 同一函数值**；ChTTS 闸在 announce 路径（26-r3），本腿零新闸 | `residentVoice.announce` |
| Notify | `rt.notify` | `ch.Notify = presenterNotifier(run.notify)`——runSpec 通知缝，resident 装配即 `postSystemNotification` | 真 shell 气泡 |
| Badge | `ball.SetBadgeText` | `ch.Badge = presenterBadger(rb.b.SetBadgeText)`；`rb.b != nil` 守卫＝本文件既有球面形状（`SetBadgeText` 内部 STA post，任意 goroutine 可调） | 悬浮球角标 |
| TaskList | 面板任务表（nil＝跳过） | **具名 nil**：面板任务表是快照派生的（`PumpSources.Tasks` 读名册），没有单条入口 API（票 36 的面）——30-r1 §8① 预留的跳过形状；receipt 记 `task skipped: not wired` | — |
| Panel | C17 桥结果展示（桥一字未动） | `agentRuntime.showPresentResult`：结果区 `rt.stream`（**artifact != "" 时按工件裸文件名加一行 `全文已存：<路径>`**——D10「落文件必须在面板可见」的增量；display 文本本身已由流式增量进同一结果区）＋ `publishPanelSnapshot`（唯一 C17 准许的出站面＝快照包） | 面板结果区 |
| Clipboard | 真写适配器＋**先 CheckText** | `presentClipboard.WriteClipboard`：**第一行即 `c.prov.CheckText(presentClipboardScope, risk.ChClipboard, text)`**（端口文档义务，结构上不可绕过）；拒写具名上报；真写半边＝Win32 `OpenClipboard/EmptyClipboard/GlobalAlloc(GMEM_FIXED)+RtlMoveMemory/SetClipboardData(CF_UNICODETEXT)` | 系统剪贴板 |
| Artifacts | `*memory.Store` 原样即插 | `ch.Artifacts = run.store`——`PutArtifact(ctx,name,[]byte)/ArtifactsDir()` 方法集逐字吻合，**零适配器** | SQLite 店的工件树 |
| Logger | （端口第七字段） | `slog.Default()`——resident 由 `installLogSink` 接管 ⇒ `present: routed`／`present: artifact written` 两行落盘 | 滚动 JSONL |

★**Clipboard 闸的 scope＝进程级 `"present-clipboard"`**（`voiceAnnounceScope` 同款先例，理由逐字同源：任务自己的 C25 scope 由 `admitTask` 的 revoke 在 loop.Run 返回**之前**关闭（run.go:1132-1177 的钩子注释自钉），`CheckText` 对未绑 id fail-closed（provenance.go:826），指给已关的任务 scope 会把每一次写都错杀。**具名残余**（与 announce 路径同一枚，resident_voice_engines_windows.go:88-97 已留痕）：任务污点不跨 scope 关闭存活，闸看到的污点集是进程级的；把任务污点带到呈现时点是 `internal/risk` scope 生命周期的变更，在本腿写面之外。闸的结构性义务（写真写前 CheckText 必被调、命中即拒且具名）由负控钉死，§4。）

★**assembleResidentPresenter 的 nil 取舍**：每枚 nil 端口（无球/无声链/无店/闸打不开）都是路由器自己的 skipped 注记，装配从不用替身顶上——boot 报告行逐枚点名接了什么（`wisp: present 装配：Notify Speak Badge Panel Artifacts Clipboard TaskList=nil` 形状）。

## 2 ★★★ 七通道端到端（真 loop×四档；`Test30Wr1*` 全绿，plain＋`-race` 双 rc=0）

真装配根（`assembleRuntime`）×真环路×真通道适配器。**mockllm×3＋golden C5 缝×1**（对照 30-r1 的 19 PASS 形状）：

- **为什么不是回声模式**：cmd/wisp 的 prompt 恒带 scene 后缀（`ContextMessage` 的当前时间块），mockllm 回声回的是**后缀块**（首跑读数逐字：`echo: ## scene\n当前时间：2026-10-31 07:56 +00:00`）——30-r1 的「回声＝样本本身」成立在它直驱 `agent.New` 的零 scene 形状上。解法＝**同一个真 mockllm 子进程的 golden 选择器**（请求 model `golden/<名>`，`server.go:216` 的公开选择路），golden 目录由本腿自备（共享 `testdata/` 零写入），应答文本逐字可控、**服务器侧 `RouteCount("chat")==1` 照读**；tier-2 走 `golden.ServerFromBytes`（C5 缝，AGENTS §1.3 许可面，30-r1 §8③ 同款具名绕行）。

| 案 | 走法 | 读数（现跑逐字） |
|---|---|---|
| tier-4 短句 | 真 mockllm 子进程×golden 单元 | `RouteCount==1`；剪贴板**恰 1 次**且内容＝应答逐字 `把这句复制一下`；通知 **0**（旧路已让位）；工件 0 |
| tier-3 中长 | 真 mockllm（100 汉字） | `RouteCount==1`；通知**恰 1 次**（双发＝红）；剪贴板 0；工件 0 |
| tier-1 表格 | 真 mockllm（markdown 表） | `RouteCount==1`；**工件真落盘**（内容＝display 逐字）；面板指针行 `全文已存：<路径>`（CorrelationID＝工件裸名）；通知 **0**（长结果不发通知＝D10）；剪贴板 0 |
| tier-2 超长＋摘要 | golden C5 缝（>400 汉字＋`口播摘要：长文已归档，请查收`） | `rep.Requests==1`（单请求规则）；工件真落盘且**剥掉摘要行**；指针行；通知 0；剪贴板 0 |
| 反向另一半 | 无 presenter 宿主（CLI 形状） | completed 结果仍走旧通知路**恰 1 次**（`Test30Wr1UnwiredHostKeepsTheLegacyPresentation`——和解门没把旧路打死） |

★`-race` 一发＝**rc=0、0 DATA RACE**（首跑曾报 4 条 race——**全部落在本腿测试 launcher 自己身上**：轮询 `bytes.Buffer` 与 exec 的拷贝 goroutine 并发；改为 `StdoutPipe`＋registry 收割的 adaptertest 同款形状后归零。产品码零 race。）

## 3 ★★★ 反向判据：两发变异（仓外 mutant ＋ `-overlay`）＋ 还原核

| 发 | 变异 | rc | 逐字红句（现跑节选） |
|---|---|---|---|
| **M1** | `run.go` 摘 Dispatch 调用点（条件置 `false &&`，调用点死） | **1** | tier-4 `clipboard writes = 0, want 1 on tier 4`；tier-1 `artifacts = [], want exactly 1 landed file on tier 1`；tier-2 同句；clean `clipboard writes = 0, want exactly 1`——四档呈现塌回无路由形状 |
| **M2** | `resident_present_windows.go` 摘 CheckText 闸（`refused = false` 硬置） | **1** | poisoned 子案 `clipboard writes = [凭据是 SK-wisp-30wr1-tainted-secret], want 0: the ChClipboard gate must refuse before the real write`——**污句直抵真写半边**，闸的牙是真的 |

★**还原核**（两发之后，同一会话）：①无 `-overlay` 重跑 rc=0（`ok … 12.012s`）；②`git hash-object` 与变异前逐字同值——`run.go=514f73e1…`、`resident_present_windows.go=e61c82db…`；③树内只有本腿 5 件（2 改 3 新），无变异残留。

## 4 ★★ Clipboard 闸核（端口义务的三层读数）

1. **负控**：`rt.c25.Mark(presentClipboardScope, risk.SrcFSRead, …, "凭据是 SK-wisp-30wr1-tainted-secret")`（postWire 缝，闸 scope 打开之后）＋ 回执文本含污源 ⇒ **真写半边恰 0 次调用**＋运行面具名 `present: 剪贴板写入被 ChClipboard 闸拦下（来源 …）`＋持久日志 `"ok":false`（receipt 具名，不被吞）。任务本身仍 exit 0——拒绝是呈现层降级，不是任务失败。
2. **正控**：干净短句 ⇒ 恰 1 次真写、内容逐字——闸活着，不是常拒的死闸。
3. **live 往返**（`-tags winlive` 单跑）：真 Win32 写→读回**逐字同**（`TestLive30Wr1RealClipboardWriteRoundTrip` PASS）；写在真机留下的就是探针标记（`Get-Clipboard` 现读证实）。读半边首跑撞上「写入后 shell 查看器抢 OpenClipboard」的 Win32 既有行为——加**固定 40 次×25ms 尝试预算**（无任何钟差参与判定）＋逐调用具名后绿。

## 5 ★★★ 零回归（staging 全套；381-v3 完整清单）

| 面 | 读数 |
|---|---|
| **六族＋WireR1**（staged 二进制＋三 DLL 同目录＋WISP_* 六 env＋CWD=树根） | 合跑 **35 PASS**（`Test15e 6/6＋Test26R3 3/3＋Test41a2r1 5/5＋Test381r2 5/5＋Test43 9/9＋Test384 2/3＋TestWireR1 5/7`）＋两族按注②③单跑补齐（`Test384 e2e`、`TestWireR1` 两枚）⇒ **并集 38 PASS／0 FAIL／0 SKIP**（同一 staged 生产码；两族单跑发生在 staged 二进制重跑，非另一份产品码） |
| `LATENCY`（15e e2e 打印） | `openDelay=262ms chunk=640ms ballLatency=902ms` ——与 26r3/381-v2/381-v3 基线**逐位同**（第 21 处定式要求的 SLO 读数） |
| present 全套 | `go test ./internal/present -count=1` ＝ **19 PASS／0 FAIL／0 SKIP** rc=0（19 PASS 不动） |
| 邻接 | `go build ./cmd/wisp` rc=0；`go vet ./cmd/wisp` rc=0；`go vet -tags winlive ./cmd/wisp` rc=0；E2E `-race` rc=0 |

★staging 途中两处环境形状（具名，非产品）：①**System32 的 onnxruntime 1.17.1 遮蔽 PATH 上的 1.28**（28-r1 §6 已测的加载序）——首跑 15e/41a2r1/WireR1 加载族 9 枚秒红；改用 **staged 二进制＋DLL 同目录**（381-v3 清单的原始形状）后归零。②`WISP_PUNC_MODEL` 首跑给成**目录**，端口要的是 **onnx 文件路径**——WireR1 两枚「标点级未接上」；改 `…\model.int8.onnx` 后 7/7。③`Test384ResidentWatchdogEndToEnd` 合跑里 3 秒「树回落 Armed 帽」时窗未中（合跑前序用例的残留 goroutine 抬树）——**单跑 PASS（1.42s，全套 384-READ 读数齐：unload passes=1、D43#9 kws.keep-alive-alert=1、settle failures=3 by contract、D43#35 ui.one-click-restart=1）**。三处均与产品码无关：本腿改动（装配段在 `runResident`、调用点在 `execute()`）不在这些测试的路径上。

## 6 ★ 四数 ＋ 两尺 ＋ d22scan（第 21 处）

| 面 | 读数 |
|---|---|
| `cmd/wisp` 改前（HEAD `f6928ea01`） | `git ls-tree -r HEAD -- cmd/wisp` ＝ **167 文件** |
| `cmd/wisp` 改后（盘上） | **170 .go 文件**（＋3 新件 33,385 B：§页首逐枚）＋2 件改（`run.go`、`resident_windows.go`，合计 59+/20−，diff 现取） |
| 甲 `gofmt -l cmd/wisp` | **空**，rc 0 |
| 乙 `gofumpt -l cmd/wisp`（绝对路径） | **空**，rc 0 |
| d22scan 全仓 | `go run . -root <repo>` **rc=0 clean**；分母：bans#1-5 internal/=285、cmd/=50、ban#8 internal/=723、cmd/=174、design/=39、frontend/=116、tools/=25 |
| d22scan 导出树（`git archive HEAD` 18,134 文件＋`.scratch` 补入＝20,835） | **rc=0 clean**（bans#1-5 internal/=283、cmd/=49、ban#8 internal/=719、cmd/=171） |
| 正控 | 导出树 `internal/present/zz_d22scan_probe.go` 植真 U+2705 字符串 ⇒ **rc=1**，逐字 `internal/present/zz_d22scan_probe.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)`；删探针 ⇒ **rc=0**（可逆） |

## 7 ★ 边界／交回父代理（逐条自陈）

1. **旧路与路由器的和解是本腿的具名决策**：装配说明只给了装配形状与调用点，没说旧 notify/announce 行怎么办。保留原样 ⇒ completed 结果双 toast＋双播报（且 tier-1/2 发通知直接违背 D10 档表）；本腿的和解＝`rt.presenter != nil && res.Status == agent.StatusCompleted` 才 Dispatch，else 分支**逐字保留**旧 notify＋wire-r1 announce（失败/控制词/截断形状永远走旧路——`present.Input` 文档自钉路由器只见 completed）。CLI（presenter 恒 nil）字节同行为，由 `Test30Wr1UnwiredHostKeepsTheLegacyPresentation` 钉住。
2. **Panel 适配器的映射具名**：display 文本不重复追加（流式增量已在同一结果区，重复＝面板双行）；适配器加的是流式行带不动的增量——工件指针行（键＝工件裸文件名，`present.ArtifactName` 铸的、`memory.ValidArtifactName` 守的，**不发明第二套 stream key 拼写**）＋快照发布。C17 桥一字未动。
3. **Clipboard scope 残余**（§1 ★）：与 announce 路径同一枚已登记残余，归 `internal/risk` scope 生命周期，非本腿写面。
4. **TaskList＝具名 nil**：票 36 的面没有单条入口 API；路由器记 skipped 不炸。
5. **live 往返件**会真改本机剪贴板一次——按 43/228 live 层规矩 `-tags winlive` 单跑、不进默认套件；本腿跑了一发（PASS，剪贴板留在探针标记上）。
6. **零 goroutine／零墙钟超时（产品码）**：Dispatch 全同步跑在任务 goroutine；慢通道是 announce 队列与 STA post（既有宿主机制）。测试码两处计时（readiness `observe.NewTimeout`、live 读 OpenClipboard 固定尝试预算）均无钟差判定。零配置键、零 emoji、`filepath` 只用 `Join/Base`（无 `Clean/Abs`）。
7. **`DEFERRED` 面**：零新增标记，无需登记表动作。

## 8 未定义即停：未触发

①七通道真对象全部在（逐枚 §1 右列）；②C17 桥形状与 present 端口对得上（`showPresentResult` 走快照出站面，桥/白名单/四键零碰）；③Clipboard 义务可满足（CheckText 在公开面，scope 是宿主事实且已具名）。开工三核的索引偏差（他人 staged 的 .scratch 一枚）已页首具名。
