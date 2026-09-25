# 35 — 面板快照泵 r1：**独立对抗验收 r1**（非实现者程）

> 被验物：`5821e24`（`internal/panel/pump.go` 新 341 ／ `pump_test.go` 新 309 ／
> `internal/agent/approval/pending_read.go` 新 60 ／ `internal/panel/bridge.go` 8-1 只改注释）
> ＋ `37a4705`（`cmd/wisp/panel_pump.go` 新 252 ／ `panel_pump_test.go` 新 430 ／ `run.go` 88-6）
> ＋ `1485921`（`panel_pump_test.go` 1-1，编排者修自己造的恒红）
> ＋ 自述件 `docs/evidence/s1/35-panel-snapshot-pump-r1.md`。
> 本件由**另一程**写，不复用被验程的尺。总裁在最后一节。

---

## 0. 锚点与口径（自量，不沿用简报）

| 项 | 读数 | 怎么量的 |
|---|---|---|
| 锚点 | `9ed209873fa71e7657d368d77bbbb83edfd770f0` | `git rev-parse 9ed2098`（本机） |
| 判据树 | `D:\tmp\wisp-accept-35pump\anchor` | `git -c core.autocrlf=false -c core.eol=lf archive 9ed2098 \| tar -x -C <副本>` |
| 本程取数时刻 | `2026-09-25 17:20 +0800` | `date "+%Y-%m-%d %H:%M%z"` |
| 本机 HEAD | `c192a4f` —— **不等于锚点**，编队仍在提交 | `git rev-parse HEAD` |
| 被验物是否被我改过 | 未改（本程只写本文件一枚） | 每枚 commit 只带本文件路径、删除列 0 |

**口径纪律（本程一律带）**：锚点归档副本 ／ 本机工作树是两把不同的尺。§5.2 那枚
`TestC21DesignTokensFourWayAgree` 红按上一程是"副本 0 枚／工作树 1 枚，两个数都对"，本程沿用该口径、不拼单值。

**一处与派单简报不符，先报名（不按简报的样子写）**：简报说自述件 **747 行**，锚点副本现量
`wc -l docs/evidence/s1/35-panel-snapshot-pump-r1.md` = **760 行**。差的 13 行正是锚点那一枚
`9ed2098 docs(35 快照泵 §8 处置)` 追加的编排者处置段（`§8` 末"我按甲做了"那块，现量在 `:728-746`）。
⇒ 简报那枚 747 是 `a45c18e` 时代的数（台账 `a45c18e` 逐字写着"证据件 747 行完整"），**不是本锚点的数**。

---

## 1. 第一格（最值钱）：**"这次没有接上任何一条入站路由"——本程自己造尺重走**

### 1.1 被验那句的原样与它的自述方法

被验对象是 `35-panel-snapshot-pump-r1.md:690-698`（§8 给 owner 段）里的这句：

> "这次**没有**顺手接上任何一条'从面板进来的批准'……我用一把**会亮的尺**量的：同一条尺打在历史上
> 真接过入站路由的两枚提交上分别亮 **6 次**和 **3 次**，打在这一批上只亮 **1 次**，而那 1 次是一枚
> **只读**的枚举（`Queue.LiveApprovals()`）。"

**先落一枚程序性裁定：这一句在自述件里没有测量记录。** 本程在锚点副本全文（760 行）里找它的支撑：

| 要找什么 | 现量 | 命令 |
|---|---|---|
| "亮 6 次"的尺定义／命中清单 | **0 处** | `grep -n "入站\|亮 6\|会亮的尺\|正控" <自述件>` → 只有 `:690` 那一句自己 |
| 被当正控的两枚提交的 sha | **未点名** | 同上；文件里出现过的 sha 全表（`grep -oE '[0-9a-f]{7,40}' \| sort -u`）15 枚，无一被这句引用 |
| 文件里唯一带"正控＋命中数"的那格 | `:232` "正控：同一条尺打在 `94071ff` 上是 **7 枚命中**" | ⇒ 那是 **§2.2 的 512 截断尺**的正控，**不是入站路由尺**；`94071ff` 真存在（`git cat-file -t` → `commit`） |

⇒ 所以这一格**不能按"它量过了、我复算它的数"来结**，只能**另造一把尺现量**。下面 1.2–1.6 就是另造的那把。
**这一条本身记为 F-PUMP-1**（见 §4）：**给 owner 的安全结论不许只有人话段里没有测量记录的一句话**。

### 1.2 本程自造的尺 R-A（"入站授权边"尺）——定义

一句话：**在 `internal/panel/**` ∪ `cmd/wisp/**` 的**非测试**代码里，有没有任何一行，
从"页面能敲的那扇门"或"宿主可变的状态"这一族符号里，伸进"能把一票批准结算掉"的那一族符号。**

三族符号是**我点名**的，逐族给来路（不是从被验件抄）：

- **结算族（碰一次就能把一票变成 allow/reject/烧 nonce）** —— 现量取自 `internal/agent/approval/`：
  `q.allow(`（`gate.go:581`、`:615`；`queue.go:343` 是定义）· `q.reject(` ·
  `DecideFromNative`（`gate.go:610`，**allow 侧的路由**）· `DecideFromPanel`（`gate.go:622`）·
  `.Veto(` · `grantNonce(`（`queue.go:326`）· `revokeGrants(`（`queue.go:370`）·
  `.deliver(` · `.expire(` · `.abandon(` · `.Allow(`（`nativeAPI.Allow`）。
- **页面门族（今天唯一能从渲染侧进来的名字）** —— `ParseComposerRequest`（`bridge.go:84`）·
  `HandleModeRequest`／`HandleWorkspaceRequest`／`HandleAttachmentAdd`／`HandleMessageSend` ·
  `DecodeAttachmentPayload` · `RequestWorkspaceSwitch`（`workspace.go:76`）· `ComposerRequest`。
- **宿主可变状态族（改了就不是原样）** —— `ModeWriteHandler` · `perm.Store.Set`（`store.go:175`）·
  `SetWorkspaceRoot`（`workspace.go:44`、`:87`）· `PermissionMode`（`store.go:155`，**读**）·
  `confirmModeSwitch`／`PendingApproval`（`run.go:472`／`:479`，**发起问、不是答题**）。

**尺的形状**：对每枚提交 C，取 `C^` 与 `C` 两棵树的命中集（**去 `_test.go`、去注释行**），
报 `ADDED = |hits(C) − hits(C^)|`。**要的是"这一枚往门里伸了几行"**，绝对数没意义（绝对数只在两棵树之间比）。

驱动件：`D:\tmp\wisp-accept-35pump\ruler.py` ＋ 本程两发内联 python（仓外，路径见 §5）。

### 1.3 正控：这把尺是活的（**含本程自己另找的三枚真提交**）

派单写"拿一把已知会亮的旧提交当正控，它给了两枚、你自己 `git log` 另找或核它给的那两枚真存在"。
被验件没给 sha ⇒ **本程自己用 `git log` 找真接过路由的提交**，逐枚具名：

| 提交 | 它当时做了什么（现量 ADDED 行逐枚点名） | 尺亮几次 |
|---|---|---|
| **`12d8e28`** `feat(114,AC#2)` | 新增 `internal/panel/composer_handlers.go`：`type ModeWriteHandler struct`、`func (h *ModeWriteHandler) HandleModeRequest(...)`、**`if err := h.Modes.Set(ctx, to, PanelModeOrigin, h.actor()); err != nil {`（`:141`）**、`PermissionMode()`×2、`record(...)`、`actor()` | **+7** |
| **`8e10095`** `feat(92,AC#1-#7)` | 新增 `ParseComposerRequest`（`bridge.go`）＋ `workspace.go` 的 `RequestWorkspaceSwitch` 与 **`scope.SetWorkspaceRoot(actable)`（`:87`）**＋ `PathScope.SetWorkspaceRoot` 接口声明 | **+4** |
| **`11f3927`** `feat(114,AC#2)` | 把门装进装配根：`cmd/wisp/run.go` 的 `modeWrites *panel.ModeWriteHandler` ＋ **`rt.modeWrites = &panel.ModeWriteHandler{`（现锚点 `:396`）** | **+2** |
| **`63ef895`** `feat(77,AC#2,#5,#7)` | 只加了 `cmd/wisp/panel_assets.go:68` 那枚 `panel.NewApprovalCardView` **建造者**调用者 | **+0** |

⇒ **尺是活的**（三枚真路由各亮 7／4／2），**且它不是大锤**：第四枚同样是"往 `cmd/wisp` 里加了一个对
`internal/panel` 的生产调用"，但它调用的是**纯建造者**、不是门也不是结算，尺**不亮**。
这正是我要的形状——**能区分"多一个建造者的调用者"与"多一条入站授权边"**。

（与被验件自述的 6／3 相比：枚数不同、我这三枚也不同。**这两组数不抵账**，本程只认自己这一发的 7／4／2，
因为它逐枚点得出名、且我核过 `12d8e28`／`8e10095`／`11f3927` 三枚的 ADDED 行确实是我那三族的符号。）

### 1.4 打在被验这一批上：**ADDED = 0 / 0 / 0**

```
5821e24: hits 17->17  ADDED=0 REMOVED=0     （pump.go 新 341 / pending_read.go 新 60 / bridge.go 注释）
37a4705: hits 17->18  ADDED=1 REMOVED=0     ← 见下面这一行
           + cmd/wisp/run.go | PermissionMode | Mode:      rt.modes.PermissionMode,
1485921: hits 18->18  ADDED=0 REMOVED=0
9ed2098: hits 18->18  ADDED=0 REMOVED=0
```

**那唯一一次亮，就是这句 `cmd/wisp/run.go:423`。** 本程逐层把它拆到底，判它**不是**入站授权边：

1. 它是 `panel.PumpSources.Mode` 字段，类型是 **`func() risk.Mode`**（`internal/panel/pump.go:97`）——
   **一个读数器，不是一次调用**。写档位的是另一个方法 `perm.Store.Set`（`internal/perm/store.go:175`），
   全批没有新增对它的引用（`rt.modes` 的非测试引用锚点现量只有三处：`run.go:386` 赋值、`:423` 交读数器、
   `:438 Modes: rt.modes` —— 第三处是 `11f3927` 时代就有的门装配，**不是本批新增**）。
2. `PermissionMode()` 本体（`store.go:155-164`）逐行读完：只 `s.mgr.Config()` → `c.PermissionMode()`，
   **无写、无副作用**，读不到时 fail-closed 回 `risk.DefaultMode()`。
3. 泵里这一支唯一消费点是 `pump.go:159-162`：`mode = p.src.Mode(); modeReadable = mode.Valid()`，
   结果只流向 `NewComposerState(...)` 与 `composer.Mode`（`pump.go:171-176`）——**面板档位显示**。
   它不进 `perm.Set`、不进 `q.allow`、不进任何 `Request{Allow:...}`。

⇒ **与被验件的自述相对照**：它说"只亮 1 次，那 1 次是 `Queue.LiveApprovals()`（只读枚举）"。
**我这把尺亮的那 1 次是另一枚符号**（`rt.modes.PermissionMode`，只读档位），
而 `LiveApprovals` 在**我这一族里根本不出现**（它不在结算族也不在可变状态族里，只是一份 `[]LiveApproval` 拷贝）。
⇒ 两把尺各自只亮 1 次、且亮的都是"读"。**结论方向一致，但枚数与所指符号不同，两发不抵账，本程记自己这一发。**

### 1.5 (ⓒ) 逐枚答：新增的 4 枚导出符号里，有没有任何一枚的调用点能反过来喂进授权决定

本程现量的**非测试调用点全表**（锚点副本，`grep -rnE --include=*.go`，去 `_test.go`）：

| 新导出符号 | 全部非测试调用点 | 终点是什么 | 能不能喂进授权／`bridge.Execute` |
|---|---|---|---|
| `panel.NewSnapshotPump` | **1 枚**：`cmd/wisp/run.go:421` | `assembleRuntime` 里的装配语句 | ❌ 装配根自己，且不在任何页面门之后 |
| `(*SnapshotPump).Publish` | **1 枚**：`cmd/wisp/panel_pump.go:237`（在 `publishPanelSnapshot` 内） | → `p.Snapshot()` → `src.Out` = `rt.bookPanelSnapshot`（`panel_pump.go:151-157`：存 `rt.lastSnap` ＋ 打一行 ledger `Info`） | ❌ 出口是**日志＋进程内留档**，两处都不碰队列 |
| `(*SnapshotPump).Snapshot`／`Marshal` | `Snapshot`：`pump.go:191`、`:205`（泵自己）；`Marshal`：**非测试调用者 0 枚** | 只构造值 | ❌ |
| `approval.Queue.LiveApprovals` | **1 枚**：`cmd/wisp/panel_pump.go:61` | 返回 `[]LiveApproval`（`Decision` 原样拷贝），调用方只读字段 | ❌ 定义体 `pending_read.go:42-60` 只 `q.mu.Lock()` ＋ 遍历 `q.pending`，**不触 `deliver/allow/reject/dropLocked/grantNonce/revokeGrants` 中任何一枚**（本程逐行读完） |
| `panel.CardViewFromDecision` | `internal/panel/approval.go:66`（`NewApprovalCardView` 委托，**本批之前就有的 1 枚**）＋ `pump.go:72`（`NativeVerdict.CardView`）＋ `cmd/wisp/panel_assets.go:68` 经 `NewApprovalCardView` 间接 | 纯函数 → `ApprovalCardView` 值 | ❌ 无副作用、无写点 |

**反向也查了**（这才是"接错一次就自己批自己"的真正形状）：结算族在锚点的非测试命中集
= **20 行，全部落在 `internal/agent/approval/{gate,queue}.go` 之内**，`internal/panel/**` 与
`cmd/wisp/**` **各 0 枚**。`5821e24^` 的同一把尺是 **20 行、同一批行**（只有 `bridge.go:77→:84`、
`run.go:393→411`、`:437→472`、`:444→479` 这类**行号漂移**）。⇒ **这一批没有把结算族向 panel/cmd 方向挪动一行。**

**并且：`Native().Allow`／`DecideFromNative`／`DecideFromPanel` 在全仓的非测试调用者是 0 枚**
（锚点现量：命中的三行都在 `gate.go` 里，是它们**自己的定义**）。
⇒ 今天**没有任何一条**路径能把一次答案送进队列，无论原生侧还是面板侧；
`wisp run` 的 L2 卡片只能超时转拒绝（`run.go:462-465` 的注释逐字写着"A console run has no native channel,
so this refuses"，与 `gate.go:487`/`:522` 的 abandon／expire 路径对得上）。

### 1.6 硬禁线"由面板侧来源的 L2 允许"：碰没碰

**裁定：没碰，且今天这根线两侧都是空的——但请把下面这三行读成"射程"，不要读成"已被接线挡住"。**

| 侧 | 现量（锚点副本） | file:line |
|---|---|---|
| 页面进来的那侧 | `ParseComposerRequest` 的**非测试调用者 0 枚**（全仓 `.go` 里只有定义 `bridge.go:84` ＋ 三处注释 `run.go:225`、`bridge.go:16`、`composer_handlers.go:37`） | `cmd/wisp/run.go:225` 那句 "Nothing calls it yet" **未被本批推翻** |
| 宿主答题的那侧 | `DecideFromPanel` / `Native().Allow` 非测试调用者 **0 枚** | 定义在 `internal/agent/approval/gate.go:622` / `:580` |
| 类型名册 | `type PanelBridge` 全仓声明 **0 枚**（`grep -rn "type PanelBridge" --include=*.go . \| wc -l` → `0`）；18 枚 C17 名册抽查 3 枚（`panel.resync`／`approval.decide`／`cost.summary`，固定串 `-F`、去 `_test.go`、去 `tools/d22scan`）非测试命中 **各 0** | — |

**最坏后果形状（写给 owner，与自述件 §8 同一条、我独立复算过）**：
如果将来有人在这条线上接错一次——**把页面能敲的那扇门接到 `DecideFromNative`（或直接 `q.allow`）而不是
`DecideFromPanel`**——那块面板就能自己批自己：一次 `postMessage` 就能签发一次 L2 允许。
今天两侧都没接，所以这一发**不是"被挡住了"，而是"还没人接到会被挡的那一步"**。
⇒ 这一点很重要：**`DecideFromPanel` 的拒 allow 逻辑（`gate.go:622-636`，含真 grant 也烧掉
`revokeGrants`）今天零生产调用者**。它是一枚**装好但未通电**的保险丝。
现有会响的仪器只有 `internal/panel/l2_grant_boundary_test.go`（四面两向，**纯测试层**）。

### 1.7 第一格结论

**成立**——"这一批没有接上任何一条入站路由"这句话，在本程**另造**的尺上复算为真：
`5821e24`／`37a4705`／`1485921`／`9ed2098` 四枚的入站授权边 **ADDED = 0／1／0／0**，
唯一那一次亮是 `cmd/wisp/run.go:423` 交出一个**只读档位读数器**（`func() risk.Mode`），
本程逐层跟到 `composer.Mode` 显示字段为止；结算族在两棵树之间**逐行相同、全部留在
`internal/agent/approval` 之内**。正控三枚真提交（`12d8e28` +7／`8e10095` +4／`11f3927` +2）
证明这把尺会亮，`63ef895` +0 证明它不是大锤。

**附一枚不成立的自述**（记债 F-PUMP-1）：§8 那句"两枚提交分别亮 6 次和 3 次"在自述件里
**没有尺定义、没有 sha、没有命中清单**，本程无法复算它那一发，只能另造 ⇒ 结论方向对，
**但那一发今天不构成凭据**。

---

## 2. 第二格：§4 的逐字段承重自证——三发自己复现 ＋ 三发它没打的

### 2.0 台件与取数形状（先自报我这一发的尺是什么形状）

| 项 | 读数 |
|---|---|
| 靶树 | `D:\tmp\wisp-accept-35pump\mut` = `git archive 9ed2098` 解出来的副本 ＋ 从工作树拷进去的 `third_party/sherpa-onnx/*.dll` **三枚**（锚点副本里 0 枚，`git ls-files third_party/sherpa-onnx \| wc -l` → `0`，**它们是未跟踪件**） |
| 用哪支正路 | **不是 CI 那支**。是"把三枚 DLL 放进 `PATH`、逐包 `go test -count=1 -v`"那一支，与被验件 §5.1 选的是同一支（**不经 GUARD A/B、不经 `runtests.sh`**）⇒ 本节的表**不等于 CI 会打出的表**，这一条与被验件 §7 第 1 条同罪，本程不替它免 |
| 判红绿 | 只认 `^--- FAIL:`；`[build failed]` / `setup failed` **单列、不算红**（这一条是本程自己踩出来的，见 §2.5 仪器记录） |
| 基线（锚点副本、**逐包**跑） | `./internal/panel/` **ok、0 红** ／ `./cmd/wisp/` **ok、0 红** |
| 被验物有没有被我改过 | 六枚被碰文件逐枚 `git hash-object --no-filters` 对 `9ed2098:<path>` **全等**（`pump.go`／`pump_test.go`／`pending_read.go`／`panel_pump.go`／`run.go`／`bridge.go`） |

### 2.1 复现它表里那三发（我自己造、自己取红名红句，没读它的日志）

| 我的编号 | 摘掉的真来源读法（file:line） | 我量到的 | 它的编号／它的读数 | 判 |
|---|---|---|---|---|
| **A1** | `cmd/wisp/panel_pump.go:70` `Reason: it.Decision.Reason` → 常量 `"A1 hardcoded"` | 红 **1 枚** `TestRunBooksWithASnapshotOfItsLiveQueue`，红句 `panel_pump_test.go:240: reason: console="R2: 目标路径在授权目录之外: C:\\Windows\\win.ini" packet="A1 hardcoded", want the same native reason on both surfaces` | M1，红句同一枚 `:240` | **复现** |
| **A2** | `cmd/wisp/run.go:423` `Mode: rt.modes.PermissionMode` → `func() risk.Mode { return risk.ModeAskEveryStep }` | 红 **1 枚** 同名用例，红句**两枚**：`:293 composer.mode.current = "ask_every_step", want the boot read "ask_high_risk"` ＋ `:296 …, want the non-default档 this config set`；booked record 自报 `mode:ask_every_step` | M3，`:293` 与 `:296` 都响 | **复现**（且这一发是在 `1485921` 已补夹具的树上打的 ⇒ 与被验件 §4.5 的 `M3` 同形，不是 `M3_no_wfix` 那形） |
| **A3** | `internal/panel/pump.go:178-181` 整块（`now := time.Now` **与** `if p.src.Now != nil` 两支一起）→ 恒常量钟 | `internal/panel` 红 **2 枚**：`TestThePumpBuildsThePacketFromWhatTheHostHolds`、`TestPublishHandsTheBytesToTheAttachedExit` ／ `cmd/wisp` **0 红、rc=0** | M4 响两枚（它给的红句在 `pump_test.go:102`、`:292`，枚数与名一致）／ §4.4 M4b 说生产那一支换成常量**全绿** | **两发一次复现**，见 §2.2 |

### 2.2 A3 这一发顺手把 §4.4 那枚"装饰品"裁定也独立量了

A3 是 M4 与 M4b 的**合体**（注入钟读法与生产默认钟一起摘）。分腿读数：

- `internal/panel` **2 枚红** ⇒ M4 那一支是承重的（"注入的钟有没有被读"有钉子）。
- `cmd/wisp` **0 枚红** ⇒ §4.4 那句"`generatedAt` 在 `wisp run` 的实际装配上不被任何用例管"**复算为真**：
  生产的包即便钟被焊死成常量也不会红。红名/红句都不是 `generatedAt`（本程现量 `cmd/wisp` 那三枚
  泵用例只断言非空与 bytes/sha/depth/mode，不比对 `at`）。

⇒ **这一格我给的不是"照抄它的表"，是"它的表里两行都响对了，且合体的那一发把它们分开来了"。**

### 2.3 自加三发它表里没有的

**B1（替换级，但换的是**装配根**那一层，不是它换的库层）**
摘掉 `cmd/wisp/panel_pump.go:68-71` 那四行"原样搬运"（`Level`／`RulesHit`／`Reason`／
`SessionOverrideBlocked` 一起不抄，别的都不动）。
⇒ 红 **1 枚** `TestRunBooksWithASnapshotOfItsLiveQueue`，红句两枚：
`:234 level: packet="L0" console="L2", want the same native level on both surfaces`、
`:240 reason: console="R2: …" packet="", want the same native reason on both surfaces`。
⇒ **为什么这发不是 M5 的第二遍**：M5 打的是 `internal/panel/pump.go:72-85`（把搬运换成
"面板侧重算"），响在**两包共 6 处**；B1 打的是 `cmd/wisp` 那一层的**赋值语句**，
`internal/panel` 一行没动 ⇒ 它验的是"**装配根有没有把 verdict 抄进泵**"这一格。
**结论：搬运在两层各有一枚独立的钉，摘任何一层都会红**（库层 M5 响两包、装配层 B1 只响 `cmd/wisp`）。
本程未复跑 M5，所以"库层那一发"仍只算〔仅自述〕；**B1 是本程现量**。

**B2／B3（移除级——§7 第 10 条自陈没做那一发，今天补）**

| 发 | 摘掉的是什么 | `internal/agent/approval` | `internal/panel` | `cmd/wisp` |
|---|---|---|---|---|
| **B2** | 整枚 `internal/panel/pump.go` ＋ `pump_test.go`（只删实现会 `[setup failed]` 把整包读数吃掉，所以两枚一起摘，量"还有谁notice"） | 未测（不引用） | **ok、0 红、用例少 6 枚** | **`[build failed]`**，5 条 `undefined:`（`run.go:235 panel.SnapshotPump`、`:238 panel.StreamLog`、`:420 panel.NewStreamLog`/`DefaultStreamKeys`、`:421 panel.NewSnapshotPump`），**测试红 0 枚** |
| **B3** | 整枚 `internal/agent/approval/pending_read.go`（`Queue.LiveApprovals` 的定义） | **ok、0 红** | **ok、0 红** | **`[build failed]`**（`panel_pump.go:61` 引用不到） |

⇒ **对 §7 第 10 条那一问的答案**（"删掉整枚 `pump.go` 会不会有别的包变红"）：
**不会有任何别的包"变红"——只有一个别的包"编译不过"。**
按本仓"改后全绿可能是那条新路径零次执行"那一族的规矩，这条**必须写成形状而不是写成"有编译器守着"**：
`[build failed]` 在 `runtests.sh` 那把尺里是**另一色**，它不等于"某枚用例守住了这件事"。

⇒ **顺手量出一枚它没登记的事实**：`LiveApprovals` 在**它自己的包里零枚用例**
（`grep -rn "LiveApprovals" internal/agent/approval/ \| wc -l` → **2**，两枚都在
`pending_read.go` 自己身上：`:36` 注释、`:42` 定义 ⇒ 该包**没有**任何 `*_test.go` 引用它）。
它的行为今天只被**别人家的**用例间接管（`internal/panel/pump_test.go` 用自造的 `NativeVerdict`、
`cmd/wisp/panel_pump_test.go` 经真队列）。B3 里 `internal/agent/approval` **ok、0 红**就是这一句的实测。
记 **F-PUMP-2**。

### 2.4 第二格结论

§4 那张表**抽的三发全部复现**（M1→A1 同枚同名同红句、M3→A2 两枚红句同两处、M4→A3 两枚红名同两名），
§4.4 那枚"生产钟是装饰品"顺带被 A3 的另一腿**独立证实**。
自加三发里：**B1 新增一格承重（装配层搬运有自己的钉）**、**B2/B3 答掉 §7 第 10 条那一发**、
并量出 **`LiveApprovals` 自己包里零用例**这一枚未登记事实。

### 2.5 本程自己踩的仪器坑（写下来给下一程，不是我修的洞）

1. **两包合跑会凭空多一枚红**：`go test ./internal/panel/ ./cmd/wisp/` 并发跑，
   `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 红在
   `resident_sink_nail_127_windows_test.go:569`"the child exited through the shutdown path
   with exit status 0xc000013a"；**同一枚用例逐包单跑 ⇒ PASS**（`:549 resident leg exit: exited 0`）。
   ⇒ 本程基线一律**逐包**跑。与被验件 §5.4 末尾记的那枚"合成一发跑会量出荒谬读数"是同一族、
   但**这一枚更狠**：它不是归包错，是**多出一枚不存在的红**。
   ⚠ 那枚用例与本批**零关系**（`git log 5821e24^..9ed2098 -- cmd/wisp/resident_sink_nail_127_windows_test.go cmd/wisp/resident_windows.go` → **0 枚**）。
2. **A3 第一版我把 `time` 的 import 一起删了** ⇒ 三行 `undefined: time`、`--- FAIL:` 计数 0。
   那**不算**变异读数（是变异本身非法）。修正成"只换钟、不动 import"才是上面那一发。
   与被验件 §4.1 的 M2 首发 `[setup failed]` 是同一族失手，本程照它的处置：**那发不算数、留着说明哪发不算**。
3. 沿用上一程的口径：**`grep -c` 为 0 时回显会冒出非 bash 形状的 `No matches found`** ⇒
   本程关键计数一律 `| wc -l` 或 python 复算（§1.6 的三枚 0、§2.3 的 `LiveApprovals` 枚数都是这么来的）。

---
