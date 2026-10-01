# 票 33 面板宿主 —— 对抗验收表 `33-v2`（非实现者终裁）

**本腿代号**：`33-v2`（对抗验收腿，非实现者；产码与测试文件一字未改，见 §9 自证）
**起手锚点（本腿现取，`2026-10-01 16:47:37 +08`）**：HEAD＝`3216ba6da5e53b5239e6dfb3d1f4ed9f3e54959e`，分支 `dev`。
**票面**：`.scratch/wisp/issues/33-panel-host-c27.md`（340 行）。本腿现读复选框＝**13 未勾／1 已勾**（`grep -c "^- \[ \]"`＝13、`grep -c "^- \[x\]"`＝1；已勾那枚是 AC#11，编排者 10-01 11:1x 翻的）。
**归属规矩**：⛔ 勾与不勾一律归编排者。本表只出判语（成立／不成立／无法判），不翻任何框、不往票面写任何东西。

---

## 0. 本腿的取数口径（钉死，引用本表任何读数都带上它）

- 整包命令（唯一口径）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m`。
- 红绿只认 `--- FAIL`／`--- PASS`／`--- SKIP` 行（`t.Logf` 也带 `file:line:` 前缀，不当失败）。
- 前人读数一律具名，且注明本腿是否复跑。⛔ 本腿不复跑前人同格读数当凭据（派单 §0），但为差集归因所必需的那一发除外，且会写明。
- `staticcheck` 本机版与 CI 钉版不同 ⇒ 本表凡涉及它一律标〔未复认〕，不拿它出结论。
- 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）本表若引用，一律标「这是射程判断，不是把冻结件内容当产品凭据」。
- `frontend/**`／`design/**`：⛔ 未读、结论未引（两层禁令）。

## 1. 逐格判语表（票面 13 枚未勾，一枚不落）

| 格 | 票面行 | 判据原文摘要 | 判语 | 凭据（判据行＋我跑了什么＋逐字读数） | 若这格是假的，最坏会被谁误信 |
|---|---|---|---|---|---|
| AC#1 | `:72-73` | Show/hide/destroy 生命周期：会话内二次 show 复用同扇窗（进程树子进程数稳定）；会话 dispose 拆窗＋WebView 子进程 ≤2s 退净 | **不成立**（按票面原句：四段里一段尺会误归因、一段本腿未量） | ① **同 HWND 复用**有牙且绿：`cmd/wisp/panel_host_windows_test.go:580-582` 比对两枚 HWND，我这发十枚样本逐字 `same HWND ... across hide->re-show=true`（台件 `.scratch/wisp/probes/33/v2/latency-10x.txt`）。② **"process-tree child-count stable"** 这一维**误归因**：同一函数 `:584-587` 拿本树 webview 枚数当代理，十发里**两发（样本 6、10）当场红**，逐字 `hide -> re-show changed the msedgewebview2 set inside THIS process tree: 7 -> 6 (tree root pid 19976; single-window reuse broken)`，**而同发 HWND 身份断言是 `true`** ⇒ "reuse broken" 那句是错的，掉的是浏览器自己的子进程。③ **dispose 拆窗**：`panel_host_gate_test.go:177` 今天不再 skip（`ctorHits>0`＝`resident_windows.go:142` ＋`panel_resident_windows.go:193`；`managerDestroyHits`＝teardown 里 `rp.mgr.Destroy()`），我这发 `--- PASS`。④ **≤2s 退净**已按编排者 12:22 裁定乙搬进 `winlive`（`panel_host_windows_live_test.go:36`），CI 零岗位、**本腿未跑那一档**＝〔未量〕。⛔ ①②③我未做定向突变（N9） | 下一枚拿"名册全绿"当 AC#1 结了的人：看不见 ② 是一枚**会把正常行为判红**的尺、看不见 ④ 不在任何自动档里跑 |
| AC#2 | `:74` | 时延：冷 ≤1500ms／热 ≤200ms（10 发 P50/P95 进 SLO appendix） | **成立（带两注）** | 本腿自发 `-count=10` 同进程连续，逐字 `AC#2 percentiles over 10 runs (nearest-rank): cold P50=291.882 P95=852.634 (budget 1500) \| hot P50=32.506 P95=66.987 (budget 200)`、`P11 line ... max cold observed 852.634 ms`（未越 2000）。**断言不是 log**：`:549`／`:600` 单发、`:814`／`:817` 按 P95 尾部各 `t.Errorf`。33-r4 那枚"硬编码 HEAD"已修：逐字 `HEAD 5cae4cf8 at read time`＝本腿 HEAD。⛔ 阈值／`thresholds.go`／`docs/SLO.md` 一字节未动，我没挑好看的发（另有一发单样本 cold 732.598／hot 34.595＝`latency-1.txt`）。**注 1**＝"into SLO appendix" 的落点：`grep -nE "P50\|P95" docs/SLO.md`＝6 命中，`:117`／`:118` 是面板两行 PASS，**但那是 S0 尖峰（2026-09-19）的数**；本票的 10 发在 `docs/evidence/s1/`（编排者 11:1x 裁定 1 已改落点）⇒ ⛔ 谁都不许拿 `:117/:118` 当票 33 凭据（N7）。**注 2**＝负载敏感，见 §8 | 拿"P95 66ms"去承诺"任何负载下 ≤200ms"的人 |
| AC#3 | `:75` | 内嵌资源离线可给；不存在监听套接字（netstat 断言） | **成立（带一注）** | **突变亲验（M6）**：`panel_host_windows_test.go:368` 把 `mibTcpStateListen` 由 2 改回 10 ⇒ 逐字 `IPv4 table rows=212, this pid owned 0 before / 0 while 127.0.0.1:57505 is listening` ＋ `--- FAIL: TestAC3ListeningSocketRulerSeesItsOwnListener`；还原（md5 `38e113669fdde37a23a3032f2eae17e8` 与起手逐字节相同）后同尺逐字 `0 -> 1`（IPv4／IPv6 两族各 1）`and stops after Close (0)` ＋ `--- PASS` ⇒ 那两枚常量**承重**，33-v1 判的恒绿已真修掉。离线那一半：`TestPanelHostOpensNoListeningSocketL1`（go/ast 禁 `net`／`net/http`）＋"文档来自 embed"由 AC#13 正面证（M3 摘次序即红）。**注**＝产品断言只对 `AF_INET` 开（裁定 3 甲），且只有 windows cli 腿有分母 | 以为"L1 import 扫过"＝"运行时确实没 socket"的人 |
| AC#4 | `:76` | 焦点往返：编辑器 → 面板 → 隐藏 → 焦点回编辑器（自动＋手动） | **不成立**（自动那一半我这发真绿；票面原句两处不齐） | 自动：`TestAC4FocusReturnToPriorWindowGap33r5` 本发 `--- PASS`，读数链逐字 `foreground before any panel 0x12084c \| 编辑器 0xa50f50 \| after Show 0xa900da2 \| panel hwnd 0xa900da2 \| prevFocus recorded at Show 0xa50f50 \| after Hide 0xa50f50 \| Hide attempted restore to 0xa50f50 (SetForegroundWindow 1, SetFocus 10817360)` ⇒ 焦点确实还给编辑器；产码三处（`:368` 先取 prior／`:400-408` 拒绝自记／`:456-479` 回还并记两枚返回码）**在位**。**⚠ 但突变 M4 活了下来**：把 `setPriorFocusLocked` 的拒绝摘掉（`:404` 加 `&& false`）⇒ **两枚 AC#4 用例逐枚 `--- PASS`**（0.82s／0.74s）。机制我量到了：`:838-840` 那发 `SetForegroundWindow` 在本桌面会话**没把面板顶到前台**，`Show` 采到的 prior 仍是编辑器 ⇒ 拒绝那支从未被喂到。⇒ 用例注释 `:791-799` 自称"不依赖前台权限"**不成立**，那一半是**条件性恒绿**。手动那一半只有 owner 能报（N2） | 拿 `...RefusedPanelSample` 的绿当"摘掉拒绝也红"的人；**注释还写着它有牙，下一程会照着注释放心** |
| AC#5 | `:77` | 运行库缺失 fixture（改名/遮罩加载器）→ 降级事件＋App 存活＋原生 L2 标志开 | **不成立** | 三条逐枚现量：① `grep -rn "panel.unavailable" --include=*.go`（排 `.scratch`）＝**0 命中**，无生产者、无改名/遮罩 fixture 用例（159 枚名册无此形）；② App 存活只有 `panel_host_windows.go:276-278` 那句 nil 时的 error，**零枚用例喂得到**（要喂就得遮罩运行库，N3）；③ 原生 L2 标志无生产者。⚠ 编排者 11:1x 裁定 2 已具名更正"不走会 fatal 的 API"**不成立**（`webview.go:115/120/125` 三枚 `log.Fatal`）⇒ 连"能不能安全造出来"都没定 | 看见 `resident_windows.go:146` 那句 `panel host unavailable (%v)` 就打勾的人——那是**装配失败**的话，与运行库缺失两回事 |
| AC#6 | `:78` | 所服务文档带 CSP（响应侧检），含 `connect-src 'none'` | **不成立**（今天零进度） | `grep -rniE "content-security\|csp\|connect-src" --include=*.go cmd/ internal/`＝**0 命中**（含注释）。票面 `:65-66` 要"注入点＋默认严格 CSP"；今天落地形是 `SetHtml`（`:350`）＝**没有响应**、没有响应头可检；AC#13 选了"调换两发次序"那一支（`:308-312`），⛔ 没走 `AddWebResourceRequestedFilter`，连注入点都没长出来 | 把"离线、无 socket"读成"CSP 也齐了"的人；SPEC-06 §9 那一行到今天一个字没实现 |
| AC#7 | `:79-86` | `SnapshotPump.Snapshot()` 跨两瞬间的形状要有会响的检（`-race` 今天干净） | **不成立**（有钩子、无仪器；票面那枚函数名会把下一程引错） | 票面前件写的是"`Marshal()` 的第一枚生产调用者"。现量 `.Marshal()` 非 test 调用者＝**0**。⚠ **但实质前件今天已成立**，走的是另一扇门：`cmd/wisp/panel_pump.go:405` 调 `rt.pump.Publish()` → `internal/panel/pump.go:288` `snap := p.Snapshot()` → `:289 json.Marshal(snap)` ⇒ **"生产里读这份包全量字节"为真**，而 AC#7 要的那发"两个触发点之间状态变了 ⇒ 必须能红"**盘上没有**（`internal/panel/pump_test.go:65/:129/:155` 三枚都是单瞬间）。跨两瞬间那对读点仍在：`pump.go:197 p.src.Verdicts()`／`:206 p.src.Results()`，泵锁不护这一对 | 拿"`Marshal` 没有生产调用者"当"这格不用做"的人 |
| AC#8 | `:87-94` | 两枚零执行分支（`panel_pump.go` 摘要超 440 字符丢 ids／`EvToolStart`）各要一发"拿掉就红"或用例或删 | **不成立**（两枚分支都还在，且没找到能为它们红的用例） | ⓐ 分支还在、只是常量化：`cmd/wisp/panel_pump.go:367 const summaryClamp = 440`，用点 `:354`；`grep -rn "summaryClamp" cmd internal` **只命中产码两行** ⇒ 零枚判据件引用它。ⓑ `EvToolStart` 消费点还在（票面 `run.go:738` 已漂到 **`:1243`**），发射者现量仍＝**0**：`grep -rn "EvToolStart" --include=*.go cmd internal tools` 只有 `internal/agent/sink.go:24`（注释）／`:25`（定义）／`cmd/wisp/run.go:1243`（消费）。⇒ 票面"不许留着让人以为它们在守什么"仍未兑现。⚠ **我未做定向突变**（N9）：ⓐ 把 clamp 改成天文数字、ⓑ 摘掉 `:1243` 那支的 `changed = true` 再跑整包——该发我没跑，⛔ 不许把"我没跑"读成"跑了没红" | 以为这两行在守什么的人：ⓐ 走不到、ⓑ 不可达 |
| AC#9 | `:96-124` | `ComposerDispatch.Handle` 被非 test 文件真调用＋本机端到端留读数 | **成立（带一注）** | 生产调用点两枚：`cmd/wisp/panel_inbound.go`（CLI 缝，票面 `:113` 已登记落码）＋**今天新长的** `panel_host_windows.go:519 return m.disp.Handle(ctx, raw)`（在 `Bind(panelDispatchBinding, ...)` 回调里，`:287-290`），非 test 到达路径＝`resident_windows.go:142/148` → `startResidentPanel` → 线程上 `mgr.Show` → `Bind`。名册我这发逐名全绿：`TestAC9ComposerDispatchHasAProductionCaller`／`TestAC9InboundFlagSurfaceIsNarrow`／`TestAC9InboundLegFromStdinReachesTheWriteLeg`／`TestAC9InboundLegRefusesRosterMethodWithNoHandler`／`TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt`。**注**＝"端到端读数"那半票面 `:199-201` 指的是 **CLI 缝**（`config.toml` 档位真变），页面那一跳今天由 AC#14 两枚钉另发（两维不互替，裁定 P2）；反转钉 `TestPanelHostIsAttachedAndNamesTheWindowHops` 在 `internal/panel` 包，不在本腿名册射程，其牙属旧账不重判 | 拿"窗口能开"当"界面点一下后端真收到"的人：H10 到了，**H1（前端发送腿）在 `frontend/**`，本腿禁读、不可判** |
| AC#10 | `:245` | `[panel]` 的 `width`／`height`／`scale` 有生产者（真去设窗口边界，`hot` 生效） | **不成立** | 定义处还在（`internal/config/schema.go:530-538`，票面 `:525-535` 已漂）。读取方现量 `grep -rn "Panel\.\|PanelCfg\|cfg\.Panel" --include=*.go cmd/wisp/ \| grep -v _test`＝**0 命中**；宿主把尺寸**写死**：`panel_host_windows.go:272-273 Width: 420, Height: 260`。⇒ 票面要区分的两边**都不是**（既没读也没设）；`Enabled`／`KeepAliveInSession` 同样零读取方 | 看见窗口真开出来就以为尺寸是被配的——用户改 `[panel] width` 今天**没有任何后果** |
| AC#12 | `:37` | "能构建"≠"有页面可发"：判据要区分"embed 只匹配到占位文件"与"真有一包页面产物" | **无法判**（尺齐了；勾不了——它自己的完成条件是归口未定） | 尺这半齐：`panel_host_gate_test.go:77 TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 问 Go 侧能力（Built／Resolve／Check／Manifest）＋一条 tracked-vs-ignored 出处轴，本发 `--- PASS`，两形分开报：工作树＝**有真页包**（`entryIDProbes` 逐字 `AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]`），而 d22scan 那行 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]` **独立证明它是未跟踪的工作树产物**（库里 `frontend/dist` 仍只有 `.gitkeep`）。勾不了的原因**是票面 `:37` 自己写的**："在'谁把 dist 填上'落定前勾不了，与票 248 AC#9 同一枚归口"。⚠ 我未对该尺突变（N9；33-r4 自报做过，本腿不采信） | 拿"AC#12 的绿"当"面板有内容"的人：库里那一形**没有内容** |
| AC#13 | `:39-40` | 冷启动的往返探测不许把真页面盖掉（次序/过滤器＋会响的断言＋反控） | **成立（带一注）** | 三条完成判据逐条对上：① 次序已重排——`panel_host_windows.go:308` 先 `firstRoundTripLocked`、`:310` 后 `serveEntry`（票面 `:39` 指认的旧序反过来）；② 断言问的是**文档里有没有 embed 入口的真元素**、由页面自己报回（`panel_resident_windows_test.go:294-331`），不是"`SetHtml` 被调过"；③ **反控本腿亲跑（M3）**：换回旧序 ⇒ 逐字 `AC#13 page answer (head 5cae4cf8): "0" - 0 of 1 probe id(s) present in the live document` ＋ `--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.70s)`；还原后 md5 `00336ed0704e51f4cdb043a6d0b3b852` 逐字节相同。**注（弱处不算缺陷）**＝入口只声明得到 **1 枚 id（`root`）** ⇒ 它证"最终文档不是那枚探测页"，证不了"是那一包页面" | 拿这一格当"面板已能显示真内容"的人 |
| AC#14 | `:42-44` | "Go→页面"这一跳由谁投递；一枚会响的用例＝真页面调绑定并 await，断回话真到达页面 | **成立** | 投递者有名字：`cmd/wisp/panel_resident_windows.go:239 w.Run()`（专用 STA 线程把泵交给库；`:262-284 handOverPump` 先发布 `wRef` 再最后一次排干通道）。**两枚钉分维＋配对突变亲验（M5）**（`post()` 不再走 `Dispatch` ＋ 把 `w.Run()` 换成 Go 侧自泵）⇒ 逐字 `AC#14 nail 1 (reply hop), page's own words: "TIMEOUT-2S,TIMEOUT-2S,TIMEOUT-2S" (Go's handler was reached by 3 of the 3 real requests)` ＋ `--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (7.09s)`，**同一发** `AC#14 nail 2 (Eval push hop) title="PUSHED-33R5-OK"` ＋ `--- PASS`。⇒ 回话确实**只有 `Run()` 送得到**、推送不受它影响——正是票面 `:43`"两半要两枚钉"的正面读数。裁定 P2 拒绝的三种形状在测试文件 `:20-23` 被逐字写成不许用 | 拿"Go 收到了 3 枚请求"当"回执到页面了"的人——这格最容易被冒充，而 M5 正好把冒充形打红 |

## 2. 门禁读数（本腿自己取，逐条带命令）

### 2.1 写面与锚点卫生

| 尺 | 命令（逐字） | 读数 |
|---|---|---|
| 起手 HEAD | `git log -1 --format=%H` | `3216ba6da5e53b5239e6dfb3d1f4ed9f3e54959e`（`2026-10-01 16:47:37 +08`），分支 `dev` |
| 产码自我锚点起未动 | `git diff --numstat 096aafad..HEAD -- cmd internal` | **空**（⇒ 本腿测的 `cmd`／`internal` 与编排者 16:37 那一发**逐字节同码**，两发读数可直接对表） |
| 工作树写面 | `git diff --stat HEAD -- cmd internal` | **空** |
| 本腿名下提交 | `git log --oneline 096aafad..HEAD` | `8fb99755`（骨架）＋后续逐节提交；`3216ba6d` 是编排者的 |
| 票面复选框 | `grep -c "^- \[ \]"`／`grep -c "^- \[x\]"` | **13／1**（与派单一致；本腿一枚未碰） |
| 台件目录未被占 | `ls docs/evidence/s1/ \| grep 33-panel-host`；`ls .scratch/wisp/probes/33/` | 本腿到达时名册＝`r1 r4 r5 r6 r7 v1`／`a1 a2 h1 p1 r1..r7`，**无 `v2`** ⇒ 代号可用 |

### 2.2 整包（本腿自跑一发，独占机器）

命令逐字：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m`
起跑 `16:49:45`／终 `16:53:22`，逐字末两行：`ok  	github.com/CarlosShao/wisp/cmd/wisp	203.831s` ＋ `RC=0`。
台件＝`.scratch/wisp/probes/33/v2/fullpack-1.txt`（1,331 行）。

| 口径 | 本腿 | 编排者 16:37（`33r7-head-roster-1.txt`） | 编排者 14:07（`33r5-head-roster.txt`） |
|---|---|---|---|
| `grep -c '^--- PASS'`（**只顶层**） | **159** | 159 | **156** |
| `grep -c -- '--- PASS'`（顶层＋子项） | **240** | 240 | 237 |
| `^--- FAIL`／`-- '--- FAIL'` | 0／0 | 0／0 | 1／1 |
| `^--- SKIP` | **0** | 0 | 0 |
| 末行 | `ok ... 203.831s` rc=0 | `ok ... 238.545s` rc=0 | `FAIL ... 235.183s` rc=1 |

**⚠ 派单 §4(丁) 那句"156 涨到 240，口径＝`grep -c -- '--- PASS'`，同 `-v`"经复量＝口径混了，差集不是 84 枚而是 3 枚。**
同一份文件两把尺各量一次：`33r5-head-roster.txt` 顶层 **156**／含子项 **237**；`33r7-head-roster-1.txt` 顶层 **159**／含子项 **240**。
⇒ 156 那发用的是 `^--- PASS`（顶层），240 用的是含子项那一把。**涨的真实枚数＝顶层 156→159（＋3），含子项 237→240（＋3）**——不是 84 枚新用例，也不是有用例被洗掉。
（这一条不影响"逐名对过"的必要性，只影响差集该有几行；逐名见 §2.3。）

### 2.3 名册逐名差集（本腿自建的三枚名册，同一命令口径）

尺：`grep '^--- PASS\|^--- FAIL\|^--- SKIP' <log> | sed 's/^--- \([A-Z]*\): \([A-Za-z0-9_]*\).*/\2/' | sort -u`
台件：`.scratch/wisp/probes/33/v2/rosters/top-33r5.txt`（157 名）／`top-33r7.txt`（159 名）／`top-mine.txt`（159 名）。

- `diff top-33r5 top-33r7` ＝**恰两行新增、零行删除**：
  - `TestAC13BringUpSurvivesAReusedThreadQuit`（33-r6 落，`cmd/wisp/panel_resident_windows_test.go:445`）
  - `TestAC13BringUpRefusesAThreadWithAQueuedClose`（33-r7 落，同文件 `:522`）
- `diff top-33r7 top-mine` ＝**空**（本腿这发与编排者那一发逐名同集）。
- **静悄悄消失＝零枚**；**静悄悄多出来＝上面两枚**，都在票 33 射程内（AC#13 那族顺序依赖），归因见 §4。
- ⚠ 33-r6 还删过一枚 `TestAC13BringUpSurvivesAStaleCloseOnAReusedThread`（派单 §4(甲) 点名要复核的那枚）。本腿现量：`grep -rn "StaleCloseOnAReusedThread" cmd/ internal/` ⇒ **只在注释里剩指认**，函数名在树里 **零枚**，且它**从未进过任何一发名册**（r5 那发的 157 名里没有它）⇒ 删除不产生名册缺口，但它是 33-r6 自己刚加的，见 §4(甲) 判语。

### 2.4 其余门禁（终态）

| 门 | 命令 | 读数 |
|---|---|---|
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；分母：bans #1-5 `internal/=224` `cmd/=34`；ban #6 `frontend/=85`；ban #7 `internal/tools/=23`；ban #8 `design/=39` `frontend/=85` `internal/=476` `cmd/=81`（⚠ 分母＝文件枚数，不是违规数）。另记：`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`——**这行本身证明工作树的 `frontend/dist` 里真有一包 assets**（我只读 git/d22scan 的元数据，⛔ 未读 `frontend/**` 任何内容） |
| gofumpt | `"$GOPATH/bin/gofumpt" -l cmd/wisp/ internal/panel/`（GOPATH＝`D:\work\base\gopath`） | **空输出** ⇒ 名下无未格式化件 |
| staticcheck | 未跑 | 〔未复认〕：本机版与 CI 钉版不同，派单 §0 禁止拿它出结论 |
| `go build ./...` / `go vet` | 未跑 | 口径偏离由编排者补跑销账（派单 §0）；本腿只用整包那一发作凭据 |
| `go mod tidy` / `go get` | **未跑** | 派单 §0 禁（HEAD 上 tidy exit 1）；相关债务裁见 §7 |
| 桌面卫生（起手 16:49:45 ／ 终态 16:55:23） | `tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV \| grep -c msedgewebview2` | **12 → 12**（零增量）。逐枚按 `--user-data-dir` 归因（PowerShell `Get-CimInstance Win32_Process`）：5 枚 `MicrosoftWindows.Client.CBS_*`（系统自带）＋6 枚 `com.clipsync.desktop`（第三方应用）＋1 枚 CBS 子项＝12；**零枚指向本仓**（既无 `%TEMP%\wisp-33r5-panel-profile` 也无 `%APPDATA%\wisp`）。⛔ 不指控谁泄漏，只报数 |
| 孤儿子进程 | `tasklist //FI "IMAGENAME eq wisp.test.exe"`／`eq go.exe` | **0／0** |

### 2.5 ⛔ 本腿写面自证

本腿**未改**任何产码或测试文件。所有定向突变（§5）都用 `git cat-file blob HEAD:<path> > <path>` 还原，
还原前后各一次 `certutil -hashfile <path> MD5`，且每次还原后核 `git diff --stat HEAD` 为空。逐枚读数在 §5 每节末尾。

## 3. 判不动的地方（逐条：甲／乙／不做 ＋现量＋为什么判不了＋具名交回）

| # | 哪一格／哪一问 | 判不了的原因（带现量） | 处置（甲＝我补取数／乙＝归口谁／不做＝代价） |
|---|---|---|---|
| N1 | **静态必答①：CI 那枚 windows 命令行上到底有没有 WebView2 Runtime** | 盘上可读的三件事，逐条现跑：① `ci.yml:388`（`test-windows`）与 `:533`（`slo-smoke`）＝ **`runs-on: windows-latest`**（GitHub 托管镜像，⛔ 不是本机 self-hosted；self-hosted 只有 `:591` 的 `slo-full`）。② `grep -noE "webview|msedge" .github/workflows/ci.yml` ＝ **0 命中** ⇒ 工作流里**没有任何一步提到 WebView2 运行库**。③ `grep -niE "install|winget|choco" ci.yml` ＝ **2 命中，逐枚读＝`:138 go install mvdan.cc/gofumpt@latest`／`:233 … staticcheck@2026.2.1`**，两枚都是 Go 工具链，**没有一枚装运行库**。镜像出厂带不带＝**只有 runner 自己能答**，而派单 §5 明令⛔ 不许"先推一次看看"，`git push` 归编排者 | **乙**：把"问 OS"那一发归给编排者（一枚一次性 job 打印 WebView2 注册表键即可结）。⚠ **但这一格其实左右不了本票结论**——真正的 CI 后果在 §6：`runtests.sh` 把任何 `--- SKIP` 判红 ＋ 干净检出里 `frontend/dist` 只有 `.gitkeep` ⇒ AC#13 那枚用例必 Skip ⇒ windows cli 腿红，**与运行库在不在无关** |
| N2 | **AC#4 的 `manual` 那半** | 票面 `:76` 逐字要 "automated **+ manual**"。手动那半＝人在真桌面上看焦点回没回编辑器，只有 owner 能报 | **乙**：交回编排者向 owner 取一句（⛔ 我不拿自动读数冒充手动那半，也⛔ 不许谁拿 §AC#4 的自动判语说"AC#4 结了"） |
| N3 | **AC#5 的真 fixture** | 票面 `:77` 要"改名/遮罩加载器"那一形。造它＝动本机 Edge/WebView2 运行库的安装态（owner 机器上的东西，本腿无权动），且依赖里 `webview.go:115/120/125` 三枚 `log.Fatal` 会把测试进程**当场打死**（票面 `:260` 编排者 10-01 已复认在位）⇒ 拿到的不是红用例而是 `exit status` | **乙**：先由编排者拍"允不允许在本机临时遮罩运行库目录"（这是一次真机变更，不是测试卫生）；拍了才谈得上判这格。本腿判语只到"有没有仪器"那一层（见 §1 AC#5 行） |
| N4 | **MTA panic 的成因** | 三档矩阵的现象有读数（STA 561ms／未初始化 531ms／MTA 当场 panic），但那枚 HRESULT 要自己实现未导出的 environment handler 才拿得到（`33-p1` ⑦-结-3，本腿无法在不改依赖的前提下取） | **不做**：引用这一族一律带"成因未定值、只定现象"（编排者裁定 P3 已这么写）。⛔ 不许任何程把这格升成"已解释" |
| N5 | **`Dispatch` 队列在真常驻里的真实长度** | `dispatchq` 未导出（`webview.go:58`），外部只能间接读（`33-p1` R25：4 枚 WM_APP 被取走／0 枚闭包执行） | **不做**：编排者裁定 P5 已定＝记欠账、不当门。本腿沿用，并带"未直接量"这句 |
| N6 | **`33-r8`（球侧 `ui-sta` 收摊那一形）今天会不会命中** | 写面在 `internal/ball/**`，⛔ 本腿不改产码；且今天 159 枚顶层名册里**零枚**用例真去"球线程收摊后紧接着建面板窗"（现量：`grep -n "sta_windows\|ui-sta" cmd/wisp/*_test.go` 无该形）⇒ 盘上没有可复现的红 | **乙**：归 `33-r8`。本腿只登记"没中 ≠ 不会中"（照抄编排者裁定 5 那句，不升格也不降格） |
| N7 | **AC#2 的"P95 进 SLO appendix"那一半：落点归谁** | 我第一版把这里读错了，现量更正：`grep -nE "P50\|P95" docs/SLO.md`＝**6 命中**，其中 **`:117`「面板冷拉起（会话内首次）≤1500ms ＝ P50 880–1126ms · P95 1042–1256ms（12 子进程真冷 ×2 轮）PASS」**、**`:118`「面板热拉起（同窗口 show）≤200ms ＝ P50 26–71ms · P95 50–80ms PASS」**。⇒ **appendix 是存在的，但那两行是 S0 尖峰（2026-09-19）的读数，不是本票宿主那 10 发**；本票那 10 发在同目录 `33-panel-host-c27-r4.md`（`grep -c "P95"`＝11）。`docs/SLO.md` 与 `thresholds.go` 都在⛔ 禁改清单 ⇒ 我这一发也只能落 `docs/evidence/s1/` | **甲＋乙**：甲＝本腿自己重发一发同口径读数（§8），⛔ 不与 `33-r4`／`33-r7` 合并成一句，各自带时刻＋HEAD；乙＝**"SLO.md 那两行要不要用本票宿主的读数换掉"是账目归位＋动禁改文件，归编排者拍**（按记忆第 14 款不上 owner 清单）。⚠ 引用 `:117/:118` 时不许说成"票 33 达标凭据"——它是尖峰期数字 |
| N8 | **`33-r5`/`33-r7` 证据件里我没复跑的那些读数** | 派单 §0 规定 `-v2` 禁重跑前人读数当凭据。r5 的"终跑名册／收尾三把尺"两节编排者说**没代填**；r7 的「交件判语」逐字写着"待填" | **甲**：本腿不代填也不覆盖，只在 §4 具名指出"这一节在盘上是空的"，判语一律改用**我自己那一发**（§2.2）＋突变（§5） |
| N9 | **本腿未做定向突变的格子**（见 §5 末尾清单） | 预算闸门：本腿在 100 轮内只兜得住若干发突变。没突变的尺我在判语里明写"未亲验其牙" | **乙**：具名交回编排者补派，⛔ 不许把我"没红"读成"我验过没红" |

## 4. 三枚代提增量的承重判断（派单 §4(甲)）

（填写中）

## 5. 定向突变台账（每把尺都可能恒绿吗）

（填写中）

## 6. CI 那格（静态必答①）

（填写中）

## 7. 两处过期指认＋`go.mod` 那格（静态必答②）

（填写中）

## 8. 时延分位复跑（派单 §4(丙)）

（填写中）

## 9. 总裁

（填写中）
