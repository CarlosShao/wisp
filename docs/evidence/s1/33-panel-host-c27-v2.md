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

### 2.6 突变循环之后的终态复跑（同一命令口径，逐字）

⛔ 这不是"再取一次前人读数"，是**自证那六枚突变没留痕**——还原只靠 md5 相同不够，还要行为与名册相同。

| 尺 | 命令 | 终态读数 |
|---|---|---|
| 整包第二发 | 同 §2.2 那条，`17:14` 起／`17:18:32` 终，台件 `.scratch/wisp/probes/33/v2/fullpack-final.txt` | 逐字 `ok  	github.com/CarlosShao/wisp/cmd/wisp	254.271s` ＋ `RC=0`；顶层 **159** PASS／**0** FAIL／**0** SKIP，含子项 **240**／0／0 ⇒ **与 §2.2 那一发四数全同**（耗时 203.831 → 254.271，同一口径下负载差，⛔ 不据此判任何格） |
| 名册逐名 | `diff rosters/top-mine.txt rosters/top-final.txt` | **IDENTICAL**（159 枚逐名同集 ⇒ 六枚突变进出没吞掉、没新增任何用例） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**（终态件 `.scratch/wisp/probes/33/v2/d22scan-final.txt`；分母逐字不变：bans #1-5 `internal/=224` `cmd/=34`，ban #8 `internal/=476` `cmd/=81`） |
| gofumpt | `"D:/work/base/gopath/bin/gofumpt" -l cmd/ internal/` | **0 行输出** |
| 写面 | `git status --porcelain -- cmd internal docs/evidence` | **0 行**（⛔ 产码与测试文件一字未留改动；⛔ 零枚删除、零枚改名） |
| 桌面 | `tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV \| grep -c` | 起手 12 → 终态 **12**；`go.exe` **0** 枚（没留孤儿子进程） |
| 本腿提交枚数 | `git log --oneline 096aafad..HEAD` | `8fb99755`（骨架）／`5cae4cf8`（门禁读数）／`eeaded86`（13 格判语）／`5ce429ff`（突变台账＋CI＋过期指认＋时延＋总裁）＝**四枚，逐枚只带 `docs/evidence/s1/33-panel-host-c27-v2.md` 这一枚 pathspec**，⛔ 未 push |

### 2.7 ⛔ 本腿写面自证（这一节的标题在我插 §2.6 时被吞掉过一次，本节现按盘上重编为 2.7；内容未改一句）

本腿**未改**任何产码或测试文件。所有定向突变（§5）都用 `git cat-file blob HEAD:<path> > <path>` 还原，
还原前后各一次 `certutil -hashfile <path> MD5`，且每次还原后核 `git diff --stat HEAD` 为空。逐枚读数在 §5 每节末尾。
终态复跑见 §2.6。

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

尺：`git show --numstat <枚>` ＋ `git diff <枚>^ <枚> -- <路径>` 逐名读删除列。原话都是死腿写的、由编排者落盘——**本腿问的是"这些增量会不会让下一程读错"**。

| 枚 | 盘上现量 | 本腿判语 |
|---|---|---|
| `a842e5b8`（r5 六件正文） | **1 枚文件**：`docs/evidence/s1/33-panel-host-c27-r5.md 95/8`。删除列逐名读＝**全是它自己骨架里的小标题与占位句**（`## 六件（正文在各节；本腿硬顺序＝…）` 那行、`（各件"改前读数—改后读数—反控读数"在下面的 §六件正文，本骨架先交 ⑤⑥⑦。）`），⛔ 零枚代码／票面／他人件 | **承重的部分不在这枚里**：r5 自己 `:173` 逐字写着"没写完的是「终跑名册」与「收尾三把尺」两节"，而 `grep -c "待填\|（填写中）" = 0` ⇒ **它不是留占位符死的，是留整节死的**，编排者那句"我没代填"复认成立。⇒ 这两节缺的那一维（终跑名册）由本表 §2.2／§2.3 **自己取数补上**，⛔ 不冒充 r5 的读数 |
| `84c67e7a`（r6 三枚未提交改动） | 3 枚文件：`panel_host_windows.go 45/41`／`panel_host_windows_test.go 11/1`／`panel_resident_windows_test.go 5/66`。逐名核到：① `mibTcp…` 无关；② **窄形排干收回**＝`panel_host_windows.go` 那 45/41（`drainStaleQuitBeforeCreate` 回到 `[WM_QUIT,WM_QUIT]`＋帽 64＋`slog.Warn`）；③ **harness 删掉 `runtime.UnlockOSThread()`**＝`panel_host_windows_test.go` 唯一那行删除，替换成 11 行注释（逐字 `DO NOT runtime.UnlockOSThread() here…`）；④ **删掉它自己刚加的用例**＝`git diff` 里逐字 `-func TestAC13BringUpSurvivesAStaleCloseOnAReusedThread(t *testing.T) {` | **承重，而且承的是本案那一枚的命**：②③④ 三件是一件事——"还池前不许欠消息"。⚠ **但"这枚只动了空白"那句话要分两读**：若指"编排者那次 `gofumpt -w` 只动空白"——**成立**，我复跑两把尺逐字 `5 66` vs `-w` 后 `3 64`，差值 2/2 全落在空白；若指"`84c67e7a` 这枚提交只动空白"——**不成立**，同发还删了 64 行内容（其中一整枚用例）。⇒ 请只在第一读下引用它。⚠ ④ 的删除我判**正当而非遮掩**：它期望的是"带着 stale close 也能建起来"，而 `f7d28ef0` 把产码改成"具名拒绝"，那枚钉与产品互相矛盾，留着必红；替代钉 `TestAC13BringUpRefusesAThreadWithAQueuedClose` 的牙由我 M2 亲验 |
| `096aafad`（r7 那 89 行正文） | **1 枚文件**：`33-panel-host-c27-r7.md 89/7`。7 行删除**逐名读**（不采信"全是小标题/旧判语被替换"那句）：`### (b) 整包那一发里"谁先跑、留下什么"：TID 逐枚对上`／`### 发 1（…）`／`### 发 2（…）` 三行是**小标题改写**；另 4 行是它自己那两枚 caveat 的**同句重写**——`+` 侧逐字仍在：`若错（真凶是别的还池点，或根本是 dataPath 复用/并发建窗），判据④仍会红，且我的修法（甲封点）不解决问题；` 与 `是**归因**` | ⚠ **这一枚我要具名更正一处普遍会读错的形状**：它把 `(b)` 那节从"整包那一发"改成"家族 A/B 一发里 TID 逐枚对上（**不是整包那发**）"，并给了 `tid=16272` 那种逐枚读数——**方向是收紧不是放宽**（原来说"凭具名 A/B ＋时刻相邻"，现在说"我 TID 对上了，但用的是我另造的 A/B 发"）。⇒ **判语：不承重为缺陷、承重为读数升级**；但那两行 TID 读数在本腿**不可复跑**（载体 `zz_33r7_probe_windows_test.go` 已改名成 `.bak` 封存，不在树里，⛔ 我不删不改名去复活它）⇒ 引用它一律带〔死腿自述＋载体在盘但不在构建图〕 |
| 旁证：`f7d28ef0`（产码那枚，非代提） | `66/22 panel_host_windows.go` ＋ `194/3 panel_resident_windows_test.go` | 建窗前提（`staleCloseQueued`＋具名拒绝）落在这枚，**它的牙由我 M1/M2 两发分别验过**（§5） |

**另核一项派单 §4(甲) 要的话**：那两枚 `1/1` 与 `11/1` 的删除列**里没有票面 AC 框**（`git diff --name-only` 三枚全是 `cmd/wisp/*.go`），与 `33-r4` 自报的"票面 3 枚删除不是我的"不冲突。

## 5. 定向突变台账（每把尺都可能恒绿吗）

统一流程：取 md5 → 突变 → `grep -c MUT-` 证落地 → 指名用例单跑 → `git cat-file blob HEAD:<path> > <path>` 还原 → 再取 md5 → 核 `git diff --stat HEAD -- cmd internal` 为空。⛔ 零枚 `checkout/restore/stash/clean`，⛔ 零枚删除，⛔ 零枚突变体进提交。

| # | 突变（文件:行 ＋ 改坏的那一处） | 指名用例 | 读数（逐字） | 判语 |
|---|---|---|---|---|
| **M1** | `panel_host_windows.go:261` 把 `drainStaleQuitBeforeCreate()` 塞进 `if false` | `TestAC13BringUpSurvivesAReusedThreadQuit` | `--- FAIL … (0.00s)`，红因逐字 `err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x0` | **有牙，但抓到一枚新缺陷**：它没有按注释预测的那样 panic，而是被**后加的 close 检查**拦下来——检查报的是 `WM_CLOSE … hwnd 0x0`，而那枚线程里躺的是 `PostQuitMessage` 种的 **WM_QUIT**。⇒ `staleCloseQueued()` **不专指 WM_CLOSE**：队列里有 latched quit 时它返回 `(true, 0)`，红句会点一枚**不存在的 hwnd**（细节与后果见 §5b） |
| **M2** | `panel_host_windows.go:262-264` 把 `staleCloseQueued` 那段拒绝改成死枝 | `TestAC13BringUpRefusesAThreadWithAQueuedClose` | `--- FAIL … (0.80s)`，红句逐字 `bringUp on a thread carrying a queued WM_CLOSE runtime error: invalid memory address or nil pointer dereference instead of refusing it by name` | **有牙**，且顺带把票面那枚因果链复认成读数：摘掉拒绝 ⇒ **真的**在 `chromium.go:131`  deref nil（不是"15 秒超时"那一形） |
| **M3** | `panel_host_windows.go:308-312` 两发手递换回旧序（先 `serveEntry` 后探测） | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | `--- FAIL … (0.70s)`，逐字 `AC#13 page answer (head 5cae4cf8): "0" - 0 of 1 probe id(s) present in the live document` | **有牙**（这就是票面 `:39` ③ 要的那发反控，本腿亲跑）。同时暴露弱处：入口只声明得到 1 枚 id `root` ⇒ 尺的分辨率是"1 枚"，见 §1 AC#13 注 |
| **M4** | `panel_host_windows.go:404` `setPriorFocusLocked` 的自记拒绝加 `&& false` | `TestAC4PriorFocusSurvivesARefusedPanelSample` ＋ `TestAC4FocusReturnToPriorWindowGap33r5` | **两枚逐枚 `--- PASS`**（0.82s／0.74s），`ok … 1.720s` | ⛔ **拿掉仍绿＝这把尺是瞎的**。瞎在哪：`...RefusedPanelSample:838-840` 要先用 `SetForegroundWindow` 把面板顶到前台才喂得到"样本就是面板自己"那一支，而这台桌面会话没让它生效 ⇒ `Show` 采到的 prior 仍是编辑器，拒绝那支**从未被执行**，断言自然两边都过。用例注释 `:791-799` 自称"不依赖前台权限"——**这句是错的**（`33-r5` 自己在 `:70` 就记过"发 1 ⇒ ③红"，说明它确实依赖） |
| **M5** | `panel_resident_windows.go:239` 把 `w.Run()` 换成 Go 侧自泵（`drainTasks`＋`pnlPumpOnce`），同发 `post()` 不再走 `Dispatch`（`:297 && false`） | 两枚 AC#14 钉**配对** | nail 1：`--- FAIL … (7.09s)`，逐字 `page's own words: "TIMEOUT-2S,TIMEOUT-2S,TIMEOUT-2S" (Go's handler was reached by 3 of the 3 real requests)`；**同一发** nail 2：`--- PASS … title="PUSHED-33R5-OK"` | ⭐ **本表最强的一枚**：两维当场分家——页→Go 到、Go→页的回话不到，正是裁定 P2 要的形状；且证明 `Run()` 那一行**承重**、两枚钉不互相冒充 |
| **M6** | `panel_host_windows_test.go:368` `mibTcpStateListen` 由 2 改回 10（判据侧突变） | `TestAC3ListeningSocketRulerSeesItsOwnListener` | `--- FAIL … (0.01s)`，逐字 `this pid owned 0 before / 0 while 127.0.0.1:57505 is listening` ＋ `the ruler's IPv6 layout is blind` | **有牙**；还原后同尺逐字 `0 -> 1` 两族 ＋ `--- PASS` ⇒ 33-r4 修的那两枚常量不是装饰。⚠ 同发 `TestPanelHostOpensNoListeningSocketL1` **照旧绿** ⇒ 再次证明"import 扫"看不见"运行时真开 socket" |

**还原自证（逐枚，五把尺）**：`panel_host_windows.go` 三次还原后 md5 均＝起手 `00336ed0704e51f4cdb043a6d0b3b852`；`panel_resident_windows.go` 还原后＝起手 `1c8decbc94efea82e4809bc1baf72918`；`panel_host_windows_test.go` 还原后＝起手 `38e113669fdde37a23a3032f2eae17e8`；每次还原后 `git diff --stat HEAD -- cmd internal` **空**；`grep -c MUT- cmd/wisp/*.go` 非零者**零枚**。终态另发整包一发确认（§2.2 之后又跑一次同命令＝名册一致，见 §9）。

**未做突变的尺（⛔ 不许被读成"我验过没红"，对应 N9）**：
- AC#1 的 `same HWND` 断言与 dispose 那枚 AST 尺（我只做了静态复认＋读了我自己那发的两口径读数）。
- AC#11 `TestCleanCheckoutBuilds_AC11`（`33-v1` 自报跑过反向证，本腿不重跑前人读数＝派单 §0）。
- AC#12 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`（`33-r4` 自报做过，不采信）。
- AC#8 两枚分支（ⓐ 把 `summaryClamp` 改成天文数字、ⓑ 摘 `run.go:1243` 的 `changed = true`，再跑整包看名册）。
- AC#2 那四行预算断言我只**读码**确认它们是 `t.Errorf` 而非 `t.Logf`（`:549/:600/:814/:817`），没做"把预算读成 log"的突变。
- AC#4 端到端那枚（`...Gap33r5`）：M4 只证明"拒绝自记"那一半无牙，⛔ 不能反推这枚也瞎——它在我这发里给的是完整句柄链读数。

## 5b. 新落修法（建窗前那道检查）的可达面——派单 §4(乙) 两问

### 问①：拒绝之后还有没有恢复路径？⇒ **没有。今天这形是"永久拒绝"，我判它是缺陷，不是正确的 fail-closed**

构造与读数（不靠猜，逐处指行）：`Show` 被拒 ⇒ `panel_host_windows.go:262-264` 返回 error，**队列一字节没动**（`staleCloseQueued` 是 `PM_NOREMOVE`，`:659-667`；这一点由在册用例 `panel_resident_windows_test.go:563/:582-584` 正面钉住"拒绝之后那枚 close 还在"）。线程那一侧：`panel_resident_windows.go:209-224` 的循环条件是 `rp.isCreated() || rp.startUpErr() != nil`，**被拒既不置 created、也不置 startUp**（`RequestShow:324-327` 只 `slog.Error`＋`Printf`，见 §6 那格欠账）⇒ 线程留在循环里等下一条任务；此后每一次 `Show` 都走同一条拒绝路。**盘上没有任何一段产码把该线程的队列派到空**：`pumpThreadToQuiet`／`releaseThreadClean` 只存在于 `panel_resident_windows_test.go:374/:409`（⛔ 零枚产码调用者，`grep -rn "pumpThreadToQuiet" cmd/wisp` 只命中测试文件）。⇒ **"每次都被拒、而拒绝不清队列 ⇒ 面板永久开不出来"这一形在生产路径可达**，且对用户的形状是"按了没反应"。

那它算不算正确姿势？分两层判：
- **"不吃别人的消息"这一层成立**：M1/M2 两发把替代形都量过了——摘掉检查 ⇒ `chromium.go:131` nil deref（M2 红句），摘掉排干 ⇒ 拒绝误触发（M1 红句）。remove-only 那一形 33-r7 量到"每次建窗漏一扇窗＋3 发挂 1 发"（`modes-grid.txt` mode=wide：`iter=2 rc=2`）。所以产码"只检查、只具名拒绝"是当前读数支持的**最小安全动作**，⛔ 不许谁顺手改成"替它排干"。
- **"没有出口"这一层是缺陷**：一个进程一旦落到该形就**再也开不出面板**，而它手上明明有两条读数支持的出口——(i) `stop()` 已能终止线程（`:363-401`，`TestPanelThreadIsSTAAndExitsCleanly` 钉着），线程重起就是干净队列；(ii) 拒绝时把 `startUp` 置上，让 `post()` 当场回 `false`（`:293-295`）而不是让人等 15 秒超时（§6）。⇒ 最小修法＝**把"被具名拒绝"升成线程级致命（或让 `stop()` 之后允许一次重起）**，⛔ 不是"让 `bringUp` 去排别人的队"。**归口**：`cmd/wisp/**` 的下一枚落地腿（票 248 或 228 那两发之一，与编排者 §3 条第 4 格那枚欠账同一条线）；本腿不改产码，只出这一句判语。

### 问②：四形对照（`modes-grid.txt`，每形三发）有没有假绿形状？⇒ **有一枚，且它问的正是"泵到空会不会拆别人的窗"**

先复认那格的读数没写歪（我逐行读了 83 行那份）：`none` 3/3 panic、`dispatch` 3/3 panic（`pre-create peek: filtered WM_QUIT found=false … unfiltered head=WM_USER+1(0x401)`）、`wide` 2/3 建起且 `thread windows now=4`＋`iter=2 rc=2`（挂死那发）、`pumpAfterDestroy` 3/3 建起且 `thread windows now=0`→再建回 3。⇒ **"只有泵到空成功"这一句在它那三发里是真的**。

**但整张表里"上一任"永远是同一个人**：每形的前置都是 `bringUp#1 ok … after Destroy (undispatched): head=WM_CLOSE(0x10) hwnd=…`——那扇窗的 **Go owner 已经先被 `Destroy` 掉了**，所以"泵到空"拆的是一枚**没人在认领的窗**（读数里 `thread windows 3 -> 0` 是它的效果）。**表里没有第五形**：*"这线程上有一扇窗，它的 Go owner 还活着、没被 Destroy，队列里另躺着一枚 WM_CLOSE"*。而产品要的恢复动作（问① 那条出口）一旦落地，就必然要在**不区分这两形**的前提下泵到空——那时"泵"拆掉的就是**还活着的合法 owner 的窗**。⇒ **判语：这不是测试卫生问题，是产品形状缺口**——那张四形表支撑的是"检查＋拒绝"，⛔ 它**不支撑**"任何腿可以据它实现自动排干恢复"。谁要做问① 那条恢复出口，**必须先补第五形**（一扇活窗＋一枚外来 WM_CLOSE，断言"恢复动作不许拆那扇活窗"），否则它就是拿一张没测过的表当许可证。
⛔ 本腿不实现那一形（要动 `cmd/wisp` 写面或另起探针载体并自建真窗，且派单 §4 只要求我出判语）；⚠ 我要过它就得自己重造载体——r7 那枚 `.bak` 我不能改名复活（§7 禁删纪律同样禁我改名他人件）。

## 6. CI 那格（静态必答①）——结论：**推送之后 windows cli 腿会红，而且红的成因不是运行库**

**默认档里"真开窗户"的用例枚数（现量，名册口径逐字）**：`grep -cE "TestPanelHost|TestAC13|TestAC14|TestPanelThread|TestBallPanel|TestAC4" .scratch/wisp/probes/33/v2/rosters/top-mine.txt`＝**14 枚命中**，扣掉三枚不起窗的（`TestPanelHostOpensNoListeningSocketL1`＝go/ast 扫 import、`TestPanelThreadNameIsNotInResidentRoster`＝名册比对、`TestPanelHostLatencyPercentilesAC2`＝复用别人记的样本）⇒ **11 枚今天会在这台机器上真建 WebView2 窗**，逐名：`TestPanelHostRealWindowHopAndLifecycle`／`TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC4PriorFocusSurvivesARefusedPanelSample`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC13BringUpSurvivesAReusedThreadQuit`／`TestAC13BringUpRefusesAThreadWithAQueuedClose`／`TestAC14AwaitedBindingReplyReachesThePage`／`TestAC14GoSideEvalPushReachesThePage`／`TestPanelThreadIsSTAAndExitsCleanly`／`TestBallPanelGesturesReachThePanelThread`（＋ `TestAC13ColdStart` 家族里的 dispose 复用那一枚）。⚠ 它们**没有 build tag**：`panel_resident_windows_test.go:25-28` 逐字写着"Real windows, default windows tier like the rest of this family: **a machine that cannot create one is red, not skipped**"。

**进不进 CI**：进。`ci.yml:388 runs-on: windows-latest` → `test-windows` 第 4 步 `bash scripts/wisp-cli-tests.sh`（`ci.yml` 里 `shell: bash`）→ 末行逐字 `bash "$portable" --scope=cli` → `portable-tests.sh:191-197 cli) scope=(./cmd/wisp/)` → `tools/d22scan/runtests.sh`。⇒ **这 11 枚全部在 CI 分母里**。

**三条独立会红的成因，按"最坏产物"排序**：
1. **干净检出没有页面产物 ⇒ 必有一枚 SKIP ⇒ 腿红（与运行库无关，这一条今天已经成立）**。链条：CI 用 `actions/checkout@v4` ⇒ `frontend/dist` 里只有被跟踪的那枚 `.gitkeep`（票面 `:37`）⇒ `panel.BuiltinAssets()` 走 `shape=anchor-only` 那一支（`33-r4` 两口径读数一致）⇒ `entryIDProbes` 返回 nil ⇒ `panel_resident_windows_test.go:316` 命中 `t.Skipf("AC#13 has no subject in this tree…")` ⇒ `runtests.sh:88` 数 `^--- SKIP`、`:98` 逐字 `if [ "$skipped" -ne 0 ]; then` 判红。**本腿在这台机器上量到 0 枚 SKIP，正因为这是工作树、不是干净检出**（`git ls-files` 那把尺只能证库里有什么，⛔ 我没读 `frontend/**` 内容）。
2. **运行库有没有**＝盘上读不出（N1）：`ci.yml` 全文 `webview`／`msedge` **0 命中**，两步 `go install` 是 gofumpt 与 staticcheck ⇒ 没有任何一步装它。若镜像不带 Evergreen 运行库 ⇒ 那 11 枚当场红（`NewWithOptions` 返回 nil 那支还会撞 `webview.go:115/120/125` 的 `log.Fatal`＝**整个 test 二进制挂掉**，比红更糟：名册采不到）。
3. **负载敏感那一族**：AC#1 的树枚数代理尺我自己 10 发里红 2 发（§1 AC#1），且 `slo-full` 就跑在这台机器上、每次 push 自启抢 CPU ⇒ 推送后这一枚可能随机红。⚠ 这条我**不建议**用重试或降档处理（编排者 12:12 裁定 1 那条禁区照旧），要处理就处理**尺**（那句 "single-window reuse broken" 的错归因）。

**`winlive` 零岗位复认（本腿自己跑，不引票面）**：`grep -rn "winlive" .github/workflows/ scripts/` ＝ **0 命中**；而默认档那 11 枚**全在 CI**、`winlive` 那 1 枚（`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`）**只在带 tag 的本机跑得到** ⇒ 今天真实的错位是**反的**：票面担心的"真机时序断言被搬到别人机器"确实搬了（11 枚在托管镜像上建真窗），而唯一被搬去 `winlive` 的那一枚反而是最需要安静桌面的那一枚。**登记，不重开**（裁定 7 已定：单开 workflow＝契约级）。

## 7. 两处过期指认（静态必答②）——判语：**都算该登记的债，且都归落地腿顺带改**

### ① `cmd/wisp/approval_always.go:165` 那句（逐字）

> the ball is built only by cmd/balldebug, and this binary links no WebView2 host (internal/panel/pump.go:16, ticket 33 unclaimed)

**一句话里两枚断言，今天对 `cmd/wisp` 两枚都不成立**：
- "ball is built only by cmd/balldebug"：现量 `grep -rn "ball\.New(" --include=*.go cmd/ \| grep -v _test` ＝ **两枚命中**，`cmd/balldebug/main.go:189` ＋ **`cmd/wisp/resident_ball_windows.go:165`**。⇒ 常驻进程今天自己造球。
- "this binary links no WebView2 host"：`go.mod:19 github.com/jchv/go-webview2 …` ＋ 非测试调用点 `cmd/wisp/panel_resident_windows.go:193 return NewPanelManager(...)`、`cmd/wisp/resident_windows.go:142 rp, rpErr := newResidentPanelManager(...)`。⇒ **链接且装配了**。
- **对 `cmd/balldebug` 仍成立**：`main.go:199` 逐字 `OnTrayPanel: func() { fmt.Println("tray: open panel (stub, ticket 33)") }`、`main.go:611` 逐字 `fmt.Println("hotkey: panel (stub, ticket 33)")`，且 `grep -c webview cmd/balldebug/main.go`＝**0**。⇒ 那两枚 stub 是真的，⛔ 谁都不许把本条判语说成"balldebug 也接上了"。
- **连带第三处**（同一族的过期指认，我在读它的出处时撞到的）：`internal/panel/pump.go:15-17` 逐字 "There is no Go -> page channel in this tree today - **no WebView2 host (ticket 33 is unclaimed)**, no postMessage writer…" ⇒ 前半今天由 AC#14 那两枚钉正面推翻（M5 读数），"ticket 33 unclaimed" 也已不成立（票面有 33-r1..r7 六段进度）。
- **判语＝债，该登记**：这类"词面型过期指认"在本仓只可**抽样纠**、不可当门（记忆第 77 条那条反例）。修法＝一枚动 `cmd/wisp` 或 `internal/panel` 的落地腿**顺带改口这三处注释**；按记忆第 83 条定式，"改口那句注释"本身要写成那一腿的**一格判据**，否则下一程还会照着它推理。⛔ 本腿一字未改（含注释）。

### ② `go.mod` 把直接依赖标在 `// indirect` 那一节

现量：`go.mod:19`＝`github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`，而 `cmd/wisp/panel_host_windows.go:63` 逐字 `webview2 "github.com/jchv/go-webview2"` 是**非测试文件的直接 import**。⇒ 注释与事实不符，且**不会自愈**：`go mod tidy` 在 HEAD 上 exit 1（票面 `:267` 编排者复认，并写死"派单里要写死不许跑"；本腿⛔ 未跑）。
- **判语＝债，但不影响任何正确性**：Go 的构建/版本选择**不看这枚注释**（它只是给人与 `go mod tidy` 读的标签），今天编译、链接、整包 159/0/0 都是在标着 indirect 的状态下发生的；风险只在"有人拿那一节当依赖图读"——比如据此说"这枚依赖不是我们的直接依赖，可以随便换/升"。
- **不跑 tidy 的表达方式（两档，都归编排者点，本腿都不做）**：
  - **甲（推荐，零文件面）**：在真相源台账 `docs/reports/pending-and-issues.md` 记一条 `A##`，字段写全——"现象＝`go.mod:19` 标 indirect 而 `panel_host_windows.go:63` 直接 import；成因＝tidy 在 HEAD 上 exit 1 故注释不自愈；后果＝词面读者会低估这条依赖边；解除条件＝任何一次使 tidy 回到 exit 0 的变更落地时**由那次变更顺带消掉**，⛔ 不许为它单跑 tidy 造出一枚不属于它的 diff"。（这条是账目归位，按记忆第 14 款不上 owner 清单。）
  - **乙（要新判据）**：加一枚**不调 tidy** 的静态尺——`go/parser` 扫全仓非 test import 集合，与 `go.mod` require 段每行的 `// indirect` 注释比对，"直接 import 却标 indirect" ⇒ 红；自带正控（把某枚真直接依赖的注释塞回去必须红）。⚠ 这是**新增判据件＝产码面变更**，且它一旦落地就会把今天这条债变成门红，⇒ **必须编排者拍**，本腿只把它列成选项，不代做、不代选。

## 8. 时延分位复跑＋那两枚热重显红怎么判（派单 §4(丙)）

### 8.1 本腿自己的一发（⛔ 不与 r4/r7 合并成一句）

起跑前安静度（`17:04` 现量）：`go.exe`＝**0 枚**、`msedgewebview2.exe`＝**12 枚**（全部归因于系统 CBS 与第三方 clipsync，§2.4）。命令逐字：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=10 -v -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'`
台件 `.scratch/wisp/probes/33/v2/latency-10x.txt`。**口径**＝nearest-rank ＋ **同进程连续 10 发**（量的是热重复，不是冷机首启）＋ HEAD 取数时刻逐字 `5cae4cf8`（＝本腿那枚门禁读数提交，⛔ 不是硬编码）。

| 维 | 本腿 10 发 | 对照：本腿单发（17:05，`latency-1.txt`） | 对照：`33-r4` §⑦（表内自述，本腿未复跑） |
|---|---|---|---|
| cold | **P50 291.882／P95 852.634**（max 852.634，budget 1500） | 732.598 | P50 633.534／P95 815.826 |
| hot | **P50 32.506／P95 66.987**（max 66.987，budget 200） | 34.595 | P50 32.913／P95 45.058 |
| P11（冷 >2000ms） | **未触发**（逐字 `max cold observed … 852.634 ms`） | 未触发 | 未触 |

⛔ 三行各自带时刻＋HEAD，⛔ 不许并成"冷启 P95"一枚数（裁定 6 的引用规矩）；⛔ 阈值、`thresholds.go`、`docs/SLO.md`、判据一字未动。

### 8.2 那两枚热重显红算不算 D32／AC#2 那一行没达标

**判语：那一发算"没达标"，而这一格今天算"达标"——两句话不冲突，因为票面没定"在什么负载态下判"。** 逐条给凭据：
1. `292.1 ms` 不是仪器假读数：断言是 `t.Errorf` 不是 `t.Logf`（`panel_host_windows_test.go:600` 逐字 `hot re-show %.1f ms exceeds D32 panel hot budget 200 ms`），且它取的是 `mgr.LastHotMs()`＝**产码自己量的 `HotShow` 全程**（`:429-438`），没有中间件能把 292ms 造出来。我**不能**判它假。
2. 我也**不能**判它是本票缺陷：我这 10 发在较安静的机器上最大 66.987ms，同一枚断言一次没响；`33-r7` 自报那一发同发记录"整包慢两成、机器 webview 进程 6→19"。⇒ 现有读数只支持"**这一行负载敏感**"，不支持"热路径超预算"这一结论，也不支持"那一发是仪器噪声"。⛔ 我不拿我这一发去抹它那一发（派单 §4 丙明令不许挑一发好看的代替）。
3. **缺口在口径不在阈值**：票面 `:74` 逐字要的是 **10 发 P50/P95**，而今天同一格上挂着**两把不同口径的尺**——单发门（`:549`／`:600`）与 P95 尾门（`:814`／`:817`）。单发门会让**一整包的负载态**去否决 AC#2，而 P95 门（票面真正要的那一枚）在我这 10 发里 66.987 < 200 从容通过。⇒ **要裁的是"AC#2 按哪一把尺结"**：我建议按票面原句＝**P95 那把算 AC#2，单发那两把是回归探针**，⛔ 但这要改判据语义＝**归编排者裁，本腿不动一字**。
4. ⚠ 顺带一枚只有整包才看得见的事实（对本票不利，具名）：`-count=10` 那发里 **`TestPanelHostRealWindowHopAndLifecycle` 红了 2 次（样本 6、10）而 `TestPanelHostLatencyPercentilesAC2` 十次全绿** ⇒ 喂给它样本的那枚用例失败时，百分位那把尺照样报"达标"。**AC#2 的绿与 AC#1 的红可以同时为真**，引用时不许互相顶。

## 9. 总裁

**票 33 今天不能结。** 13 枚未勾逐枚判语枚数：**成立 5／不成立 7／无法判 1**（勾不勾归编排者，本表不翻任何框）。

- **成立（带注）**：AC#2、AC#3、AC#9、AC#13、AC#14。其中 **AC#13／AC#14／AC#3 三枚的牙是我自己突变出来的**（M3／M5／M6），不引任何前人读数。
- **不成立**：AC#1（票面原句四段里，"进程树枚数稳定"那一维的尺**会把正常子进程回收判红并错归因成 reuse broken**，10 发里 2 发；"≤2s 退净"已进 `winlive` 且本腿未量）、AC#4（端到端在我这发真绿，但"M4 摘掉拒绝⇒两枚用例照绿"证的那一半**无牙**，且票面要的 manual 那半无人交读数）、AC#5（`panel.unavailable` **0 命中**、无 fixture、无 L2 标志生产者）、AC#6（CSP 连注入点都没有，**grep 0 命中**）、AC#7（实质前件已成立但用例不存在）、AC#8（两枚分支都还在、都无尺）、AC#10（三个字段零生产者，尺寸写死 420x260）。
- **无法判**：AC#12（尺齐了；票面 `:37` 自己把它绑在"谁把 dist 填上"那枚归口上，那不是我判得动的）。
- **决定性那几枚的凭据**：AC#6／AC#10／AC#5 三枚是**词面＋生产者枚数**级读数（各一把 grep，逐字在本表），最硬；AC#4 的"无牙"是**突变存活**，次硬；AC#1 的错归因是**本腿 10 发名册里的两枚红句原文**。
- **今天这六格新落地的东西我认它承重**：`Run()` 那一行、建窗前的检查、AC#13 的次序、AC#14 的两枚分维钉、netstat 那两枚常量、harness 不再 `UnlockOSThread`——每一枚我都指到行并给了一发读数或一次突变。
- **最该被编排者看见的三枚**（不是清单，是后果）：① **CI windows cli 腿一推送就红**，成因第一枚是"干净检出没页面产物 ⇒ AC#13 那枚 `t.Skipf` ⇒ `runtests.sh:98` 把 SKIP 判红"，与 WebView2 运行库**无关**（§6）；② **派单 §4(丁) 那句"156 涨到 240"是口径混用**（顶层 vs 含子项），真实名册增长＝**+3 枚、零枚静悄悄消失**（§2.2/§2.3）；③ **"被具名拒绝"在读数上仍表现为 15 秒超时**这条欠账今天没被任何一格踩到（本腿整包 0 枚 15 秒红），但它在 AC#13／AC#14／AC#4／AC#1(手势) 四族里**都可能被踩**——凡见这些族里 `--- FAIL … (15.0x s)` 一律先按"拒绝"读、不按"卡住"读（历史上 `33-r5` 那一发 15.01 s 正是这一形）。

### 9.1 一格已定案欠账的读法核对（派单 §6）

`showAndWait`（`panel_resident_windows_test.go:103-117`）等的是 `IsShown() || startUpErr()`，而 `Show` 的拒绝只落日志（`panel_resident_windows.go:324-327`）⇒ 被拒时 `waitPanelTrue` 在 `panelThreadWait = 15s` 处 `t.Fatalf("timed out after 15s waiting for the panel thread to finish Show - this is a failed measurement, not a pass")`。
**本腿答案：今天没有任何一格判语被这条错读过**——我这发整包 0 枚 FAIL、6 次突变里 5 次给出的都是**具名**红句（M2 直接点 nil deref、M3 点 "0 of 1 probe id(s)"、M5 点 "TIMEOUT-2S"、M6 点 "owned 0 before / 0"），没有一枚落成 15 秒超时。**但两条标注规矩**要写进后续每一张表：① AC#13／AC#14／AC#4／AC#1(手势) 四族的 `--- FAIL … (≈15 s)` ＝"被具名拒绝的伪装"，⛔ 不许写成"卡住"；② 派单 §6 那句"编排者裁＝记欠账、不单开票"我遵守——**本腿没有加 `lastShowErr`、没有改 `showAndWait`、没有动 `RequestShow` 那三行日志**，只登记"这一形会让读数说谎"。
