# 票 33 lane-2：CSP 真实响应头三枚／netstat 按 PID 零 LISTENING／AC#5 夹具／AC#7/8/9 读数／延迟 10 发

- **日期**：2026-11-16 12:5x–13:0x +08
- **派单**：`.scratch/wisp/probes/orch/20261116/txt/33-lane2-dispatch.txt`（编排者）
- **起手 HEAD**：`e1fa5e53033ae79780b4dd2afcb77788d1b9aa60`（读数日志里 winlive 台件自报 `head e1fa5e530`）
- **写面**：`cmd/wisp/**`（产码 2 枚文件＋新测试 1 枚）、`docs/evidence/s1/33-*`（本件）、`docs/SLO.md`（仅追加）
- **环境**：DESKTOP-LVS7839；Wisp 球常驻进程 PID 49348 全程在跑（未动它）；`frontend/dist` 先重建（新哈希 `index-DZQOdBdK.js`）跑完真页面那一发后清回 `.gitkeep`，终态＝anchor-only

## 0. 产码改动（唯一一处行为变化）

`serveNotBuiltNotice`（`cmd/wisp/panel_host_windows.go`）：未构建通知从 `SetHtml`（about:blank 基址，**无响应可带 CSP 头**）改为 `Navigate(panelNotBuiltURL = https://wisp.panel/__notbuilt__)`，通知文档由资源过滤器应答（`panel_assetLookup` 新增通知路由，字节与原文逐字相同），因此与入口页和全部资产一样**以真实过滤器响应投递、响应头带 CSP**。动机：派单第 1 条要求"未构建通知页有 CSP"，而 SetHtml 文档没有响应头通道——这是本轮量出的真缺口，不是假设。

## 1. CSP 三枚（判据＝response inspection：浏览器对真实响应的执行读数）

**仪器**（`cmd/wisp/panel_csp_response_33lane2_windows_test.go`，真窗）：页面里注册 `securitypolicyviolation` 监听 → 发一发同源 fetch（`connect-src 'none'` 必须在网络层之前掐断）→ 通过在册入向门 `window.wispDispatch` 把 `location.href` 与每条违规的 `effectiveDirective`、`originalPolicy`（**浏览器从真实响应解析出的策略原文**）报回 Go。fork 只映射了 `AddWebResourceResponseReceived` 的 vtbl 槽、无 Go 包装，所以"直接读响应头字节"在包外不可达；`originalPolicy` 是浏览器解析真实响应头后给出的策略文本，是本包可达的最强 response-inspection 读数。33r14 那两枚构造侧尺保持在册不变。

| 枚 | 读数 | 判 |
|---|---|---|
| 真页面有 CSP（构建包在场的二进制，载具第 2 发，12:58） | `href="https://wisp.panel/index.html"`，2 条 connect-src 违规；其中**响应头那份策略**逐字 `default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data: blob:; font-src 'self'; connect-src 'none'; frame-src 'none'; object-src 'none'; base-uri 'none'`；另一份＝页面自带 `<meta>` CSP（`frontend/index.html`，无显式 connect-src，掐断走 default-src 兜底），如实记录不判 | **PASS** |
| 未构建通知页有 CSP（anchor-only 二进制，载具第 1 发，12:55） | `href="https://wisp.panel/__notbuilt__"`（＝过滤器路由，不是 about:blank），1 条 connect-src 违规，`originalPolicy` 与上同串、逐字含 `connect-src 'none'` | **PASS**（产码改动后） |
| connect-src 逐字为 `'none'` | 两发读数都取自浏览器解析出的 `originalPolicy` 文本，均含 `connect-src 'none'` 原文 | **PASS** |

**反向控制（M1，overlay 变异，产码还原未落盘）**：把 `serveNotBuiltNotice` 退回 `SetHtml(panelNotBuiltDocument)`，anchor-only 内嵌重建后单跑通知页那发 ⇒ **红**，红句两条：`href is "about:blank", want the filter-served notice route "https://wisp.panel/__notbuilt__"` ＋ `ZERO securitypolicyviolation events ... no response-header CSP reached this document`。变异体零枚进提交（overlay 件在 `.scratch/wisp/probes/33/lane2/mut/`）。

**过严断言的一笔（诚实记录）**：首版断言要求**每条** connect-src 违规的策略都含 `connect-src 'none'`，被页面自带 meta CSP 的违规打红（那份没有 connect-src 指令）。裁定＝AC 行管的是**被服务文档的响应头**，断言改为"浏览器解析出的策略集中存在含 `connect-src 'none'` 的那份"，其余策略逐条记录。⛔ 未放宽牙：零违规仍然红、含 'none' 的策略缺席仍然红。

## 2. 无监听 socket（netstat -ano 按 PID）

- **命令与读数**：真窗活着期间外部轮询 `netstat -ano -p tcp`（每 2 秒，日志 `.scratch/wisp/probes/33/lane2/netstat-watch.txt`），两发抓到窗活样本：
  - `12:56:26.810 pid=16336 totalListening=44 mine=`（空＝该 PID 零 LISTENING 行）
  - `12:58:26.078 pid=4304 totalListening=44 mine=`（同）
  全机 44 条 LISTENING 全属其他进程；wisp.test 进程自身两条样本均**零 LISTENING**。
- **在册尺牙**（`TestAC3ListeningSocketRulerSeesItsOwnListener`，载具第 1 发 PASS）：IPv4 真监听 0→1、IPv6 0→1、Close 后 0——尺不是恒真。按-PID 过滤＋状态常量（表类 5／LISTEN=2）沿用 33-r4 修好的 `GetExtendedTcpTable` 读数（与 netstat 同一数据源）。

## 3. AC#5 运行时缺失夹具（现成夹具复跑，载具第 1 发，全部 PASS）

- `TestAC5RuntimeMissingRefusesNamedSetsTheFallbackFlag33ac5`：探针缝遮蔽（`panelRuntimeVersion` 返回 `("", nil)`，夹具注释已具名：loader 内嵌进二进制、改磁盘 DLL 造不出缺失世界）⇒ 具名拒绝 `errPanelRuntimeMissing`＋Evergreen 指引＋窗口未创建＋`Unavailable()` 旗亮带 reason＋`panel.unavailable` 事件发出（ERROR 日志逐字在案）。
- `TestAC5FreshManagerFlagStartsOff33ac5`：负控 `(false, "")`。
- `TestAC5NonRetiringRefusalIsObservableOnTheResidentThread33ac5`：真常驻线程上热键触发 ⇒ **拒绝 10.1633ms 可读**、线程存活、`startUp==nil`（app-alive 读数）。
- `TestAC5ProbeMapsEveryMissingAnswerToTheNamedRefusal33ac5`：空版本/loader 错误→具名拒绝（fail-closed）、真版本→放行。

## 4. AC#7／AC#8／AC#9 读数

**AC#7**（前提读数＋实验件复跑，`go test ./internal/panel/ -run TestAC7`，12:54，rc=0）：
- 前提：导出 `Marshal()` 的**直接**非测试调用者今天仍 **0**（`Publish` 内联自己的 `json.Marshal`，`internal/panel/pump.go:583-588`）；但"生产里有人读这份包的全量字节"已由票 329 的推送跳成立：`rt.pump.Publish()`（`cmd/wisp/panel_pump.go:511`）→ booked bytes → `pushPanelSnapshotToHost`（`cmd/wisp/panel_snapshot_deliver_329.go:91-105`，`lastPanelSnapshot` 的生产读者）逐字推给页面。
- `TestAC7TornPacketIsObservable` **PASS**＝那枚"会响的红"在册（确定性撕裂：一包同时载 T1 的卡与 T2 的结果）。
- `TestAC7PumpLockDoesNotCloseTheTear` **PASS**＝票面处方"持锁一次性快照"**量出不成立**（`SnapshotPump.mu` 只护 `Publishes()/count()`，每个源各持自己的锁）。
- `TestAC7PumpLockGuardsOnlyTheCounters` **PASS**＝发现以类型性质在册。

**AC#8**（载具第 1 发，5+1 枚全 PASS）：
- ⓐ 超 440 丢 ids 分支：`TestAC8SnapshotClampBranchDropsIdsAndKeepsTheDigest33r14`（`pending=-`＋depth＋digest，行 440 内）＋正控 `TestAC8SnapshotShortPacketKeepsItsIds33r14`（短包保 ids，把"分支真射了"与"函数本来就写 pending=-"分开）＋经生产泵入口的三枚（`BooksThroughThePumpEntry`／`TryEntryRefusesWhileTheBookingLockIsHeld`／`BlockingEntryStillBooks`）。"拿掉就红"的变异四读数在 `docs/evidence/s1/33-panel-host-c27-r14.md`（33r14 在册）。
- ⓑ `EvToolStart` 分支：**"不可达"已过期**——`internal/agent/loop.go:654-660` 现在把适配器的 tool_call_start 转发成 `EvToolStart`（该 kind 的首个生产发射者），分支消费进 `runFeed.ToolStarted`；`Test145RunningUsageAndDoneTravelTheWire` PASS（线上出现 running 行、end/done 退役）。摘掉 `ToolStarted` 那两行即该发红（分支与 wire 断言直连），无需再造仪器。

**AC#9**（载具第 1 发，4 枚 PASS）：
- `TestAC9ComposerDispatchHasAProductionCaller` PASS：非 test 生产调用者＝`cmd/wisp/panel_inbound.go`（33-r3 落码）。
- `TestAC9InboundLegFromStdinReachesTheWriteLeg` PASS（stdin 封套→`Handle`→真写档，`config.toml` 档位真变）。
- `TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt`／`TestAC9InboundLegRefusesRosterMethodWithNoHandler` PASS。
- ⚠ 框不翻：AC#9 的勾按票面要非实现者裁，本腿只交读数。

## 5. 延迟（票面 `:74`，winlive 档，载具第 3 发，13:01:50–13:02:09）

口径＝同进程连续 10 发（`-test.count=10`，lifecycle+percentile 逐轮配对，nearest-rank），HEAD `e1fa5e530`，本机 DESKTOP-LVS7839，期间 Wisp 球 PID 49348 常驻：

| 样本 | cold ms | hot ms |
|---|---|---|
| 1 | 1332.472 | 107.115 |
| 2 | 1164.654 | 70.004 |
| 3 | 1164.660 | 47.717 |
| 4 | 1193.290 | 92.224 |
| 5 | 1170.966 | 33.625 |
| 6 | 1203.090 | 50.681 |
| 7 | 1219.985 | 50.241 |
| 8 | 1247.368 | 120.657 |
| 9 | **1693.016** | 38.174 |
| 10 | 1265.587 | 37.274 |

- **hot：P50=50.241／P95=120.657（预算 200）＝PASS**，10 发最大 120.657 全在内。
- **cold：P50=1203.090 在内；P95=1693.016（预算 1500）＝FAIL（如实红）**。样本 9（1693.016）单发超线并当场红（该轮 lifecycle FAIL），聚合 P95 随之红。⛔ 未挑发、未重跑拧绿、未动阈值一字节。
- P11（冷拉起 >2s 才重评 L2 回原生）：**未触发**（1693.016 < 2000）。
- 归因注记：样本 9 前后无本腿负载；与 10-01 `A501` 那枚"负载下发抖"同形状（同代码历史桌面读数 964/1098/1332 内），真相归验收腿/编排者判，本腿只交原始 10 发。

## 6. 门禁终态

- `go vet ./cmd/wisp/` rc=0；`gofmt -l`（三枚写面文件）零行；`d22scan` clean（见提交信息）。
- `scripts/wisp-cli-staged-tests.sh` 三发（第 1 发 rc=0 12 枚名选；第 2 发 rc=0 真页面 CSP；第 3 发 winlive rc=1＝上表如实红）。
- `frontend/dist` 终态＝`.gitkeep`（`git status --porcelain frontend/dist` 零行）；⛔ dist 字节零提交。
- ⛔ 未动：thresholds.go／golden／allowlist.txt／PLAN.md／specs／D43／design／C17 白名单（bridge.go 零字节）／票面 AC 框（一枚未翻）。
- 仪器坑（记一笔）：pwsh 直接向 exe/bash 传 `-test.run 'a|b'` 会被截成 `-test`（载具第 1 发首跑与变异二进制首跑都踩了）；改 `bash -c`／`cmd /c` 引用传参解决。
