# 票 35 余 2 格写腿：伪 allow 必拒（35g5）＋ resync 无态（35g6）— 实现腿交件读数

**腿**：agent-92（`auto/glm-5.3-flash`），写面 `cmd/wisp/**` ＋ 本证据件；**零产码改动**（见 §5）。
**起手锚**：`e24abd8e`（dev）。**交件**：见 §6。⛔ 全程未 push、⛔ 未动 `frontend/**`、⛔ 未开真窗、⛔ 未动 C17 白名单（**零新增名**，见 §1 的理由）、⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design、⛔ 未提交 dist 字节。

---

## 0. 现读的今天桥面（派单要求的"先现读"，⛔ 不照票面旧行号）

- 名册＝`internal/panel/bridge.go` 的 `Method*` 常量，**8 枚**：`panel.mode.request` / `panel.workspace.request` / `panel.attachment.add` / `panel.message.send`（`:42-45`）＋ `config.get` / `config.set`（`:66-67`）＋ `config.models.list`（`:75`）＋ `cost.summary`（`:107`）；case 表在 `:199`。
- `ParseComposerRequest`（`bridge.go:172-190`）守卫次第：解析 `:174` → 名册 `:178-179` → 来源 `:181-184` → requestId `:185-188`。**名册守卫最前**。
- 调度＝`internal/panel/composer_dispatch.go` `Handle`（`:197`）→ `dispatch`（`:219`，8 个 case ＋ default `rosterMismatch`）；生产装配 `cmd/wisp/panel_inbound.go`（`newComposerDispatchChain`，8 门全接真腿）。
- `approval.decide` **不在册** ⇒ 从 webview 发的 `approval.decide{allow:true}` 结构上只能死在名册守卫 `:179`。
- resync＝`cmd/wisp/panel_resident_windows.go` `showOnThread`（`:560`，成功 Show 后 `publishOnShow()` `:566`）→ `publishOnShow`（`:594`）→ 引导线选定的 publisher（`resident_windows.go` 默认支 `:668`＝`newResidentPanelPump(...).publish`）→ `internal/panel/pump.go` `Publish`（`:583`，**每次现读全部 sources**）→ `PanelManager.evalPanelJS`（`panel_host_windows.go:661`，推队列 cap=1）→ `deliverPanelPush` → `Eval`。

## 1. 格一：伪 allow 必拒（票面 `:47`，新件 `cmd/wisp/panel_forged_allow_35g5_windows_test.go`）

**结构性拒绝的机理**（不是前端自觉）：`approval.decide` 不在 C17 名册，任何来源（包括完全伪造的合法 source+requestId）都被 `ParseComposerRequest` 的名册守卫具名拒绝。**C17 白名单零新增**——加名＝改契约＋撞 d22scan ban #6 与冻结的 grant-boundary 仪器（`grantRouteWords` 含 approval/decide），本腿的判据断言的恰是"这个门不开"。

**① 具名拒绝（真桥调用面）**：`TestPanelForgedApprovalDecideAllowIsRefusedByNameOnPageEdge35g5` —— 信封＝票面字面形状的**最强伪造**：`{"method":"approval.decide","correlationId":"pc-35g5-allow-1","allow":true,"source":"panel-composer","requestId":"pc-35g5-allow-1"}`（合法 source＋requestId，把后面两道守卫全部喂饱，孤立出名册守卫这一道）；入口走页面那条边（`installPanelTransport` → `openDocument` → `pagePost`＝`window.chrome.webview.postMessage`，与 (d) 三支同形，⛔ 未直调 `dispatchRaw`、⛔ 未 `Eval` 扮页面）。

**红句逐字**（页面收回的 reply，载具定向发 PASS 读数）：

```
面板请求被拒绝 [approval.decide pc-35g5-allow-1]：panel: composer request refused: 方法 "approval.decide" 不是面板 composer 通路的能力入口
```

断言面：`doorRounds==1`（门答了，不是静默丢弃）＋ reply 非空且含 `面板请求被拒绝`＋含 `approval.decide`（具名）＋含名册红句 `不是面板 composer 通路的能力入口`（与 (d) 三支共用同一常量 `guard35r3Roster`，防两处漂移）＋**不含**另两支红句（`按伪造/串台拒绝`／`缺少 requestId`）＋ `spy.reqs==0`（没有任何处理器执行了 allow）＋ 审计恰 1 条 `INBOUND-DISPATCH` 且带 `method="approval.decide"`（不是静默丢弃的第二面：审计落账）。

**② 突变必红 / 还原复绿**（`go test -c -overlay`，突变副本在仓外 `D:/tmp/wisp35g6/`，跑完 `git diff -- internal/panel/ cmd/wisp/` 空、四枚产码 blob worktree==HEAD）：

| 突变 | 改了哪一行 | 颜色 | 红句逐字（截头） |
|---|---|---|---|
| **M-forged-1** 名册开一格 | `bridge.go:199` case 行尾追加 `, "approval.decide"` | 35g5 主测 **FAIL**／正控 PASS | `35g5 RED: the reply does not carry the roster guard's sentence "不是面板 composer 通路的能力入口" ... reply="面板请求被拒绝 [approval.decide pc-35g5-allow-1]：panel: 方法在白名单内而派发表未列出: 名册与派发表不一致（requestId=\"pc-35g5-allow-1\"），处理器未调用"` |
| **M-forged-2** 真·拒→放 | 同上 ＋ `composer_dispatch.go` dispatch 加 case 把 `approval.decide` 路进 Mode 处理器（**allow 被执行**） | 35g5 主测 **FAIL**／正控 PASS | `35g5 RED: the door answered with an EMPTY reply - a forged panel allow must come back with the named refusal sentence, not a silence. nativeExitCalls=1 ... doorRounds=1 ... resolved=1 rejected=0` |
| 还原（纯净载具复跑） | — | **3/4 PASS**（35g5×2＋35g6×1＋原生路 201） | reply 回到上面的逐字红句，`ok github.com/CarlosShao/wisp/cmd/wisp 4.00s` rc=0 |

**③ L1/L2 合法路径不受影响（PASS 佐证）**：
- 原生二次确认那条路：既有具名尺 `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject`（`approval_reply_201_test.go:344`：面板投影无令牌 → 面板路 allow 按路由拒 → 原生令牌面 → 面板路 reject 真拒）——本波经载具复跑 **PASS 2.77s**。
- 同文件正控 `TestPanelForgedAllowPositiveControlLegitRequestStillDelivered35g5`：同一台件、同一信封形状、只换在册名 `panel.mode.request` ⇒ 处理器恰 1 次、requestId 原样到达、0 条拒绝审计行 —— 门按**方法**给不同裁决，不是一刀切拒绝。
- 页面今天真发的 `panel.approval.request`（`panel.ts:630-633`）＋`outcome` 形状已由 (d)-1（`panel_inbound_guards_35r3_test.go:150`）在册钉着，本腿引用不重测。

## 2. 格二：resync 无态（票面 `:49`，新件 `cmd/wisp/panel_resync_stateless_35g6_windows_test.go`）

**真桥链**（逐跳产码，仅两处 stand-in＝35r9 已立的那条缝）：`RequestShow` → `post`→`drainTasks`（面板线程用 35r9 `resume()` 的同步虚构）→ `showOnThread` → `mgr.Show` 热路（`attachPanelWindow35r9` 把 stand-in 页面接在 WebView2 COM 边界）→ `publishOnShow` → 引导默认支 publisher → `SnapshotPump.Publish`（现读）→ `evalPanelJS` 推队列 → `deliverPanelPush` → 页面 `Eval`。**零新协程**（无 Spawn 欠账）、**零真窗**。

**① 新旧值逐字段**（载具 PASS 读数，状态变化＝经链上 Manager 自己的 `SetBallStyle` 真写，词表取 schema 的 glass|jarvis）：

| 字段 | show #1（改前） | show #2（改后） | 判 |
|---|---|---|---|
| `ballStyle.style` | `"glass"` | `"jarvis"`（**新值**） | 页面拿到的是新状态 |
| `ballStyle.style`（包1 复核） | `"glass"`（保持改前值） | — | 包 1 没有事后移动 |
| `generatedAt` | `"2026-11-16T11:44:54Z"` | `"2026-11-16T11:44:55Z"` | 每次发布现打戳 |
| `modelSwitch.current`（未动维度） | 两包相等 | 两包相等 | 逐字段比较，不是"字节不同就算" |
| 包数 | — | 恰 2 | 一 show 一整包，无缓存无拼接 |

（两包 sha256 亦不同：`94ff68d7e5014b28` → `994fc67d4141381e`，SNAPSHOT 摘要行。）

**② 反向牙（摘掉"重新取快照"→返回缓存）**：

| 突变 | 改了哪一行 | 颜色 | 红句逐字（截头） |
|---|---|---|---|
| **M-resync-cache** | `panel_resident_pump.go` 的 `Out`：首包字节存 `resyncCachedBytes35g6`，此后只发缓存包 | 35g6 **FAIL** | `35g6 RED: packet 2's ballStyle.style = "glass", want the NEW value "jarvis" - the page received a stale state across hide/mutate/show (cross-show cache or a stale packet)` ＋ `35g6 RED: both packets carry ballStyle.style = "glass" ...` ＋ `35g6 RED: generatedAt is identical across the two shows ("\"2026-11-16T11:44:10Z\"") - the second packet was not built at show time` |
| **M-resync-noop**（加发） | `panel_resident_windows.go` `publishOnShow` 的 `if fn != nil` → `if false && fn != nil`（触发器整个摘掉） | 35g6 **FAIL** | `35g6 RED: after one show the page holds 0 packet(s), want exactly 1: []` |
| 还原（纯净载具复跑） | — | **PASS 1.12s** | 上表逐字段读数回到新值 |

**③ SDK 侧不缓存跨 show（file:line 证据，前端零改动）**：
- `frontend/src/lib/panel.ts:560-578`（`installPanelSnapshotReceiver`）：门"哑"——每个推包**原样**交给回调，替换语义在调用方的 state 更新，"no merging happens anywhere on this path"；卸载器只在仍持有**同一个**回调时才摘全局（StrictMode 重挂不互相清）。
- `frontend/src/lib/panel.ts:536-548`：每条 push 是**整包**，"the receiver replaces, it never merges … always ending up showing the latest package"。
- `frontend/src/lib/panel.ts:611-617`（`requestApprovalResolution`）：响应只以"fresh PanelSnapshot push"到达，**"a panel that keeps a copy of the answer would be a second state holder"** —— SDK 明文禁缓存。
- `frontend/src/lib/panel.ts:13`：文件头逐字 **"There is no cache, no store, no rehydration."**
- 架构钉（lint 半格已存在且非恒真）：`internal/panel/frontend_hygiene_test.go:169` `TestPanelFrontendIsStateless`（7 枚禁存面正则 `:48-56`，空文件集即 `t.Fatal`）；`internal/panel/pump.go:583` `Publish` 每次 `Snapshot()` 现读全部 sources＝Go 侧无跨包缓存。

## 3. 门禁读数

| 门 | 读数 |
|---|---|
| `gofmt -l`（两枚新件） | 空 |
| `go vet ./cmd/wisp/ ./internal/panel/` | rc=0 |
| `d22scan`（`tools/d22scan` 独立 module 起法） | rc=0 clean（ban #8 cmd/ 318 枚含两枚新件） |
| 载具定向发（35g5×2＋35g6×1＋201 原生路） | 突变前 **4/4 PASS**（rc=0）＋突变后还原复跑 **4/4 PASS**（rc=0），日志 `35-g5g6-logs/targeted-green-after-mutations.txt` |
| 整包 `go test ./cmd/wisp`（经载具一次一发） | 1417s，rc=1，红名册恰＝环境红 2 枚（TestAC4Focus／TestAC246，PID 49348 球在跑），新增红 0；日志 `35-g5g6-logs/fullpkg-35g5g6.txt` |
| 四发突变 | 各恰红具名那一枚（正控照绿），日志 `35-g5g6-logs/mut1-roster-opened.txt` / `mut2-allow-executed.txt` / `mutcache-resync-cached.txt` / `mutnoop-resync-removed.txt` |
| 仓零改动自证 | 四枚被突变产码件 blob worktree==HEAD（`cb8eca8c`/`94847ecc`/`a9e40ed3`/`8fa5b2d2`） |
| `go test ./cmd/wisp` 串行 | 每发前 `Get-Process wisp.test`（全程 0 命中），全部经 `scripts/wisp-cli-staged-tests.sh` 或其同款 staged-DLL 直跑形状（A1884：`D:\work\soft\Git\bin\bash.exe`），⛔ 裸跑 0 发 |

## 4. 硬约束逐条自查

- ⛔ thresholds/golden/allowlist/PLAN/specs/D43/design：零字节（本腿 diff 只有两枚新测试件＋本证据件＋日志）。
- ⛔ 放宽断言：零；新断言全部是新加的具名尺，未改任何既有尺一字（`git diff` 空）。
- ⛔ push：未 push（只 commit）。
- 新协程：**零个**（测试全程同步驱动），无 Spawn 欠账。
- ⛔ 反向动 C17 白名单：**零新增名**；35g5 判据断言的恰是名册门不开；两发"开一格"突变只在 overlay 副本里、只为证牙，跑完即弃（仓外）。
- ⛔ 提交 dist 字节：未碰 `frontend/dist`。
- ⚠ PID 49348（机主球）未动；`TestAC4Focus…`/`TestAC246…` 的环境红按派单口径记账（整包读数见 §5）。

## 5. 整包读数（A1882 口径，日志 `35-g5g6-logs/fullpkg-35g5g6.txt`，8075 行）

| 项 | 读数 |
|---|---|
| 整包 rc | **1**（rc=0 的 job 外壳读数是 `Out-File` 的，不是测试的；测试卷面＝`FAIL github.com/CarlosShao/wisp/cmd/wisp 1417.00s`） |
| 红名册（全量，恰 2 枚） | `TestAC4FocusReturnToPriorWindowGap33r5 (5.21s)` / `TestAC246DevLegIgnoresTheTestTaskInjection (39.86s)` —— **恰＝派单点名的环境红族**（⚠ Wisp 球 PID 49348 在跑），与既有已知红族作差＝**新增红 0 枚**，作差也无"消失"要记 |
| 终点锚 `TestZZ315A1RootLoopFirstBridgeCall` | **命中且 PASS**（`=== RUN` @ log:7901，`--- PASS (1.95s)` @ 7946；case B / case C 也 PASS） |
| 最后一条 `=== RUN` | `=== RUN   TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（log:8024，**在锚之后**，PASS 1.06s） |
| 新测试位置（截断点前/后） | **全部在锚（log:7901）之前**：`TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` @ :552（PASS 2.60s）、35g5 两枚 @ :2715/:2718（PASS 0.00s×2）、35g6 @ :3973（PASS 1.18s）——四枚新/佐证测试整包内全 PASS |
| 无声断掉（A1893） | 红名册逐枚具名在卷（两枚都带时长与 FAIL 行，不是无声消失）；无需单跑复量 |

## 6. 交件

- `cmd/wisp/panel_forged_allow_35g5_windows_test.go`（新增，约 230 行）
- `cmd/wisp/panel_resync_stateless_35g6_windows_test.go`（新增，约 290 行）
- `docs/evidence/s1/35-g5-forged-allow-and-resync-20261116.md`（本件）
- `docs/evidence/s1/35-g5g6-logs/`（定向绿＋四发突变＋整包，共 6 枚日志）
