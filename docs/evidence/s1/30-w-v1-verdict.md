# 票 30 装配缝终裁腿 `30-w-v1` —— D10 呈现路由的常驻装配（终裁裁决表，非实现者）

> **尺**：被裁实现＝`c15888b28`（6 枚 1085+/20−）＋`f8a951755`（口径精化）。授权来源＝`30-v1-verdict.md` §9.3 具名交回（逐字在）；材料＝`30-r1-impl.md` §8①②＋`internal/present/dispatch.go` 端口文档。★全部实验出自**仓外导出树** `D:\tmp\wisp\20261031\30wv1\export`（`git archive f8a951755` 18,138 文件；15 枚解包失败全在 `.scratch` 异形文件名探针日志、非 .go）；staging＝staged 二进制（自导出树现建）＋三 DLL 同目录（A1308 正解）＋六 env＋manifest 拷入＋CWD＝树根；涉 sherpa 补在册 PATH＝A1307 补（E2E 发已注入；六族发走 DLL 贴二进制——A1308 具名「PATH 注入无效」）。★工作树零触碰：`git status --porcelain -- cmd/wisp internal/panel internal/present internal/risk` 全程空；在飞 59-w-r1/387-r1/386-r1 脏面未碰、未归因。★`unset GOTMPDIR` 每跑、rc 全重定向现取。

## 1 ★★★ 七通道逐通道读码核（必跑①）

| 通道 | 判 | 依据（本腿现读） |
|---|---|---|
| Speak＝wire-r1 同一函数值、零新闸 | ✅ | `resident_windows.go:580`＝`src.run.voiceAnnounce = rv.announce`（wire-r1 挂点），装配段 ：598 在其后读同一字段 ⇒ `ch.Speak = presenterSpeaker(run.voiceAnnounce)` 捕获的**就是同一个函数值**；适配器（:166-168）纯转发零闸；ChTTS 闸真身＝`speak()` HOP 1（`resident_voice_engines_windows.go:875`，合成前），端口文档（dispatch.go:17-19）「NOT duplicated here」逐字核 |
| Clipboard＝写真写第一行 CheckText＋GMEM 真写半边 | ✅ | `WriteClipboard`（resident_present_windows.go:89-101）：nil 守卫后**第一行业务调用即 `c.prov.CheckText(presentClipboardScope, risk.ChClipboard, text)`**（:93），命中即具名拒绝、永不达 ：100 真写；端口义务原文（dispatch.go:28-31「MUST run … before the real write」）逐字核；真写半边＝Win32 `OpenClipboard/EmptyClipboard/GlobalAlloc(GMEM_FIXED=0x0000)+RtlMoveMemory/SetClipboardData(CF_UNICODETEXT=13)`，`CheckText` 真签名（provenance.go:826）＋`ChClipboard="clipboard.write"`（:123）现核 |
| TaskList 具名 nil 合法性 | ✅ | **「面板任务表无单条入口 API」为真**：任务表＝快照派生（`PumpSources.Tasks func() TaskRosterState` 整表来源，pump.go:180/192/316 → `taskRosterSectionFrom`），composer 白名单 `knownComposerMethod`（bridge.go:154）7 方法里无任何任务单条入口 ⇒ 具名 nil＝30-r1 §8① 预留的跳过形状，路由器记 skipped 不炸 |
| Panel＝C17 桥结果展示、桥一字未动 | ✅ | **`git diff 18a3f5055 c15888b28 -- internal/panel` ＝ 0 行**（present/risk/agent 同查皆 0）；适配器 `showPresentResult`（:190-198）＝`rt.stream.Append`（工件裸名键＋`全文已存：`指针行，`present.ArtifactName` 铸不发明第二套拼写）＋`publishPanelSnapshot`（panel_pump.go:447，唯一 C17 快照出站面） |
| Notify／Badge／Artifacts／Logger | ✅ | `run.notify`（runSpec 缝）／`rb.b.SetBadgeText`（nil 守卫＝本文件既有球面）／`run.store` 原样即插（`PutArtifact/ArtifactsDir` 方法集吻合零适配器）／`slog.Default()`；boot 行逐枚点名、每枚 nil 端口不替顶（:204-239） |

## 2 ★★★ 「让位」决策裁（必跑②——装配说明没覆盖、腿自拍的决定）

**判：正确。** ①D10 档表现核（`PLAN.md:242` 起）：tier-1/2 通道＝落文件＋面板＋播报摘要与路径（**无系统通知**）、tier-3＝系统通知＋面板、tier-4＝TTS＋球＋剪贴板（**无通知**）；旧路对 completed 结果**无条件** toast＋announce ⇒ 保留原样必致 tier-3 双 toast、tier-4 双播报、tier-1/2 违档表发通知——「让位」不是口味是档表服从。②和解门形状（run.go :1504-1532 diff 现读）：`rt.presenter != nil && res.Status == agent.StatusCompleted` 才 Dispatch，else 分支**逐字保留**旧 notify＋wire-r1 announce（失败/控制词/截断形状永远走旧路）；CLI（presenter 恒 nil）字节同行为由 `Test30Wr1UnwiredHostKeepsTheLegacyPresentation` 钉住（本腿复跑 PASS：旧路恰 1 通知）。★观察（非阻）：presenter 已装配宿主上**失败结果**走 else 旧路一形无专测——条件结构平凡（非 completed 即 else、else 为原码保留）＋present 包内测覆盖路由侧，登记为低危观察不阻结。

## 3 ★★★ Clipboard scope 裁（必跑③）

**判：进程级 `present-clipboard` 选型正确；「任务污点跨 scope 关闭存活＝382 形 B 设计行为非洞」定性成立（带一处措辞精化）。** ①机理现核：任务 scope 由 `admitTask` 的 revoke 在 loop.Run 返回**之前**关闭（run.go 钩子注释自钉），`CheckText` 对未绑 id fail-closed ⇒ 指给已关任务 scope＝每次写真写被错杀；`voiceAnnounceScope`（resident_voice_engines_windows.go:88-97）同款先例**逐字同源**。②负控钉闸牙：污句打进进程 scope ⇒ 真写 0 次＋具名拒绝＋`"ok":false` 落盘（§5）。③382 对照：形 B＝`RecordTaintInheritance`（记点拷贝 NORMALIZED 内容）＋`CheckTextWithInheritance`（播报点重推继承链、opt-in、深度上限 16）——**污点跨关闭存活在 risk 侧是有意设计非洞**；精化：形 B 恰恰**不改** scope 寿命（缝注释把修复类归「scope lifetime 变更」偏窄），但两形修复都在 `internal/risk`＋任务级积累钩、非本缝写面。④生产现状：TTSGate 继承机制**零接线**（announce :875 仍裸 `CheckText`）⇒ 两闸同担已登记残余 **R26-5**（tri-v1 立号、26-v1 记「仍在」），非本腿新洞。

## 4 ★★★ E2E 四档复刻（必跑④；staged 二进制＋真 mockllm 子进程 golden 选择器×3＋golden C5 缝×1）

`Test30Wr1*` 六案 **6/6 PASS rc=0**（本腿自己的导出树 staging 发）：tier-4 剪贴板**恰 1**（内容逐字）／通知 **0**／工件 0；tier-3 通知**恰 1**（双发即红）／剪贴板 0；tier-1 工件真落盘（内容逐字）＋面板指针行 `全文已存：`／通知 0；tier-2 `rep.Requests==1`＋工件**剥掉摘要行**＋指针行；`RouteCount("chat")==1` 逐案断言在（服务器侧计数）。

## 5 ★★★ mutation 两发复刻（必跑⑤；仓外 mutant＋`-overlay`，导出树零写）

| 发 | 变异（本腿自制 mutant，非复用实现件文件） | rc | 逐字红句（现跑） |
|---|---|---|---|
| **M1** | `run.go` Dispatch 调用点摘除（`_ = rt.presenter` 替调用） | **1** | tier-4 `clipboard writes = 0, want 1 on tier 4`；tier-3 `notifications = 0, want exactly 1`；tier-1/2 `artifacts = [], want exactly 1 landed file`；**CLI 案独绿**（旧路没被打死） |
| **M2** | `resident_present_windows.go` 摘 CheckText（`refused = false` 硬置） | **1** | poisoned 子案 `clipboard writes = [凭据是 SK-wisp-30wr1-tainted-secret], want 0: the ChClipboard gate must refuse before the real write`——**污句直抵真写半边**；clean 子案不受累 |

★**还原核**（两发之后）：pristine 二进制重跑 **rc=0 6/6**；`git hash-object` 导出树两件＝提交 blob 逐字同值（`run.go=514f73e1…`／`resident_present_windows.go=e61c82db…`，与 30-w-r1 所记同值）；「MUTANT 30wv1」仓内 **0 命中**、导出树 **0 命中**。

## 6 ★★ 闸核复刻（必跑⑥）

`Test30Wr1ClipboardGateRefusesTaintedSentence`（含在 §4 的 6/6 里）poisoned 子案：真写**恰 0**＋运行面具名 `present: 剪贴板写入被 ChClipboard 闸拦下`＋持久日志 `"ok":false`；clean 子案恰 1 次逐字（闸活着非常拒）。**winlive 真写往返**：`-tags winlive` 单跑 `TestLive30Wr1RealClipboardWriteRoundTrip` **PASS rc=0**（真 Win32 写→读回逐字同；探针标记与 30-w-r1 留下的同一串 42 字，用户剪贴板净变化＝零）。

## 7 ★★ 零回归复跑（必跑⑦）

| 面 | 读数（本腿现跑） |
|---|---|
| 六族＋WireR1（staged） | 合跑 **37 PASS／1 FAIL**——唯一红＝`Test384ResidentWatchdogEndToEnd`（3 秒回帽时窗未中＝环境形状③同形）；单跑 **PASS（1.28s）**、384-READ 四读数齐（unload passes=1／D43#9=1／settle failures=3 by contract／D43#35=1）⇒ **并集 38 PASS／0 FAIL／0 SKIP**（本腿合跑分劈 37+1 vs 30-w-r1 的 35+3，同一并集；本腿 PUNC env 一次给对 ⇒ WireR1 合跑即 7/7） |
| LATENCY | `openDelay=262ms chunk=640ms ballLatency=902ms`——**逐位同** 26r3/381-v2/381-v3/30-w-r1 基线 |
| present | `go test ./internal/present -count=1` **19 PASS／0 FAIL／0 SKIP rc=0**；`-race` rc=0 |
| E2E `-race` | `Test30Wr1` 六案 **rc=0、0 DATA RACE**（归零核：产品码零竞争） |
| 邻接尺 | `go build ./...` rc=0；`go vet ./cmd/wisp` rc=0；`go vet -tags winlive` rc=0；`gofmt -l`／`gofumpt -l` 皆空 |

## 8 ★ d22scan 导出树＋正控（必跑⑧）

- 导出树（`git archive f8a951755`）：**rc=0 clean**，分母 **bans#1-5 internal/=283、cmd/=50、ban#8 internal/=719、cmd/=174**。cmd/ 两数与 30-w-r1 **全仓**行逐字同；internal 的 285/723（其全仓）比 committed 树多 2/4＝在飞未跟踪件所致（只记机理、不归因；其导出行 283/719 与本腿逐字同，cmd 侧 49+1/171+3 恰＝本缝 3 枚新件，自洽）。design/=14 同其导出行（A207 分母机依）。
- 正控：`internal/present/zz_d22scan_probe.go` 植真 U+2705 ⇒ **rc=1** 逐字 `zz_d22scan_probe.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23)…Q-46(c)`；删探针 ⇒ **rc=0 clean**（可逆）。

## 9 ★★★ 终判（必跑⑨）

1. **残余③落地＝成立**：装配七通道逐通道真（§1），调用点＝`execute()` loop.Run 返回后对 completed 结果 Dispatch（§9.3 交回逐字），两发 mutation 牙真（§5），四档 E2E／闸核／winlive／零回归／d22scan 全部本腿复跑全绿（§4-§8）。
2. **销账资格＝具备**：票 30 已 -done 带残余；本腿确认残余③（`cmd/wisp` 装配缝＝§9.3 交回清单的执行）**已落地并对抗验收通过** ⇒ 编排者可在台账把残余③销账。
3. **残余定性三条**：①**TaskList 具名 nil＝合法非缺陷**——面板任务表确无单条入口 API（§1），真入口归票 36 的面；②**任务污点跨 scope＝已登记残余 R26-5 同源非新洞**——进程级 scope 是正确选型，修复形状＝382 形 B 继承接线（`CheckTextWithInheritance` 签名 channel 泛型）＋任务级积累钩，归 `internal/risk`／票级，非本缝写面；③**环境账三条＝零产品判语**——System32 onnxruntime 1.17.1 遮蔽（解＝DLL 贴二进制，A1308 正解）、`WISP_PUNC_MODEL` 要文件路径非目录、Test384 合跑回帽时窗（单跑 PASS）——本腿三处全部复现，A1322 已入册。

★未定义即停：未触发（派单九项全部可裁；tar 解包 15 枚 `.scratch` 异形名失败不影响任何读数，已具名）。★写面：本文件一件；临时件全在 `D:\tmp\wisp\20261031\30wv1\`（export/bin/logs/mut/tmp15e/tree.tar，只建不删）。
