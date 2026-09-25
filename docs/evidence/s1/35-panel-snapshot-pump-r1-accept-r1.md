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
   ⚠ 本程又见到一次（`grep -c "WARNING: DATA RACE" race_cmd.log` 打出 `0` 之后多出一句那句），
   处置相同：该格以同一次跑里的 `ok … cmd/wisp 75.715s` ＋ `rc=0` 两行为准。

---

## 3. 第三格：§7 第 2 条那枚并发形状——**给裁，不记"残余"**

### 3.1 派单摆出来的那个假设，先答它：**不成立**

派单写："若全程单 goroutine，那'混合两个瞬间'在生产里根本发生不了，就该写成**射程边界**而不是洞"。

**本程现量：这条不是射程边界，它可达。** 链子逐环给 file:line（锚点副本）：

| 环 | 位置 | 现量 |
|---|---|---|
| 1 每次工具调用**各起一枚 worker** | `internal/agent/loop.go:647` `h := l.reg.Spawn("tool-exec-"+taskID, "agent", root, func(c context.Context){ sem <- struct{}{} … l.dispatch(...) })` | 一个 `turn.ToolCalls` 一枚，带 owner＋recover（不是裸 `go func`） |
| 2 并行上限**不是 1** | `loop.go:669` `sem := make(chan struct{}, guard.Concurrency())` → `guard.go:129 return g.cfg.ToolConcurrency` → `guard.go:107-108` 未设或越界即取 `MaxToolConcurrency` → `budgets.go:47 const MaxToolConcurrency = 4` | `cmd/wisp/run.go:560-567` 的 `agent.Config{…}` **一个字都没设 `ToolConcurrency`** ⇒ 生产那把闸 = **4** |
| 3 worker 会走到"问用户" | `internal/tools/bridge.go:395` `a, why := b.gate.PendingApproval(ctx, *dec)`（L2 分支；`:380` 是 L1 `PendingWindow`） | 在 `l.dispatch` 之下 ⇒ **在 worker goroutine 上** |
| 4 问用户会驱动泵 | `gate.go:486` `g.ui.Prompt(ctx, p)` ＋ `:504/:516/:523` `g.ui.Update(...)` → `cmd/wisp/run.go` 的 `consoleApprovalUI` → `run.go` 那两处 `u.publish()`（`Prompt` 末尾、`Update` 的 `EventDismissed`/`EventStarted`）→ `panel_pump.go:237 rt.pump.Publish()` → `pump.go:136 Snapshot()` | 全链无一处换 goroutine、也无一处串行化 |

⇒ **只要一个 turn 里 Declare 出 ≥2 枚工具调用、其中 ≥2 枚撞到 L2 卡片，
`Snapshot()` 就有并发调用者**，队列本身还为此设计了 FIFO 与 depth badge（`pending_read.go:28-29` 的
`Position` 就是"第几张卡"）。⇒ **这一格判"可达"，不判"射程边界"。**

⚠ **这一句是对编排者简报前提的复算不符，按规矩报名、不按简报的样子改**：
简报里那句"三个触发点是不是同一 goroutine 顺序调用"的**倾向答案（是）在盘上不对**。

### 3.2 但"可达"分两种，本程把它们量开了：**混合两个瞬间 ≠ 数据竞争**

**(a) 内存安全这一支：`-race` 今天全绿——这是 §7 第 2 条自陈"没跑 `-race`"，本程补了。**

| 发 | 范围 | 读数 |
|---|---|---|
| `-race ./internal/panel/` **整包** `-count=1` | 全 59 枚 | **`ok … 4.809s`、`rc=0`、`WARNING: DATA RACE` 0 次、`--- FAIL:` 0 枚** |
| `-race ./cmd/wisp/` **整包** `-count=1` | 全 63 枚 | **`ok … 75.715s`、`rc=0`、DATA RACE 0、FAIL 0**（含那枚 `TestAC1ResidentLeg…`，这一发**没红** ⇒ 印证 §2.5 第 1 条那枚红是并发跑包造成的，不是用例的） |
| `-race` 两包合跑、`-run 'Pump\|Snapshot\|Panel\|Run\|Card\|Mode\|Composer'` | 泵相关名册 | `ok 8.106s` ＋ `ok 16.311s`、rc=0 |

**(b) 为什么 `-race` 打不到它**——四个读数器**各自带锁**，本程逐枚读到行：

| 泵读的 | 锁在哪 |
|---|---|
| `Queue.LiveApprovals()` | `pending_read.go:46-47` `q.mu.Lock(); defer q.mu.Unlock()` |
| `StreamLog.Chunks()` | `pump.go:318-319` `s.mu.Lock(); defer s.mu.Unlock()` |
| `PathCanonicalizer.WorkspaceRoot()` | `paths_workspace.go:39-40` `p.mu.RLock(); defer p.mu.RUnlock()` ⇒ ⚠ 顺带纠一处上一程措辞：§6.1′ 写"`:40 WorkspaceRoot()` 直接 `return p.workspace`"，`:40` 那行**是** return，但它头上 `:39` 是 RLock——**不是无锁直读**；它的 finding（拿到的是 fold 过的键）不受影响 |
| `Store.PermissionMode()` | `store.go:155` 本体不持锁，但它读的是 `s.mgr.Config()`，`config/manager.go:92-95` 持 `m.mu` |
| 出口 `bookPanelSnapshot` | `panel_pump.go:152-154` `rt.snapMu.Lock()` 包住 `lastSnap`/`lastSnapBytes`/`snapSeen` |

⇒ **`SnapshotPump.mu` 只守 `publishes` 那一个计数器**（`pump.go:230-238`），**不守装配** ——
所以两发 `Publish()` 可以各自拿到"一半是 t、一半是 t+Δ"的字段组合，而**每一枚字段自己都不是脏读**。
这正是 `pending_read.go:38-41` 承认的那件事。**"承认了 ≠ 量过"这一句被验方说得对，本程把"量"补上了：
可达性＝已量（§3.1 那四行链），内存安全＝已量（上面三发 `-race`），**唯一没量的是"混出来的包到底长什么样"**
（要造它得在两个读数器之间插一次状态变化，今天没有这样的钩子——见 F-PUMP-4 的闭合集合）。

### 3.3 本程量的是哪条腿，那条腿之外的今天是什么形状

- **量的是 `wisp run` 这条腿**（`assembleRuntime` → `execute` → `loop.Run`）。
- **§7 第 6 条复算为真**：常驻腿今天**不 import** `internal/panel`——
  `cmd/wisp/resident_other.go` **0 命中**、`cmd/wisp/resident_windows.go` **0 命中**、
  `cmd/wisp/resident_sink_nail_127_windows_test.go` **0 命中**；并且更硬的一层：**两枚
  `runResident()` 都不调 `assembleRuntime`**（`resident_windows.go:24`、`resident_other.go:22`
  里 `assembleRuntime`/`rt.pump`/`publishPanelSnapshot` 各 **0 命中**）⇒ 常驻进程里**没有那台泵**，
  连"顺手被装配到"的可能都没有。**被接上的确实只有 `wisp run`。**

**这一格的裁定：§7 第 2 条从"没测什么"**升格**为一条已定性的洞（F-PUMP-4），
严重度＝**显示层撕裂（混合两个瞬间）已可达、无内存不安全、无授权后果**。**
最坏后果形状（写清，别让人以为更狠）：票 33 把那根管子接上之后，那块面板可能显示
"一张其实刚被拒绝/刚超时的卡片还挂在待确认里"、或"卡片是 12:00:03.1 的队列、档位是 12:00:03.4 的档位"。
**它不会多签一次批准**（§1 那一格已裁：整条结算链非测试调用者 0 枚），
**它会让一张安全卡在这一秒说的话和下一秒做的事对不上**——所以修它要趁票 33 之前，不是之后。

---

## 4. 顺手核的另外两件事（不展开，但每条现量）

### 4.1 §6.4 那格"0 枚调用者被改掉了什么"——**两句都要现量，因为编排者要拿它改 `A248③`**

| 那一句 | 本程现量（锚点副本，非测试） | 判 |
|---|---|---|
| "`Snapshot` 现在有生产可达建造者（`run.go:421` → `Publish` → `Snapshot()` → `NewSnapshot`）" | `run.go:421 panel.NewSnapshotPump(...)`（**全仓唯一非测试调用点**）→ `panel_pump.go:237 rt.pump.Publish()` → `pump.go:205 p.Snapshot()` → `pump.go:182 NewSnapshot(cards, results, composer, now())` | **成立**，链子逐环有行 |
| "`PanelBridge` 声明 0 枚" | `grep -rn "type PanelBridge" --include=*.go . \| wc -l` → **0** | **成立** |
| "C17 名册 18 枚全 0" | 抽 3 枚（`panel.resync`／`approval.decide`／`cost.summary`，`-F` 固定串、去 `_test.go`）→ **各 0** | **抽查内成立**（其余 15 枚沿用被验件 §6.3 的全量，本程未复跑） |
| "`ParseComposerRequest` 非测试调用者 0" | 全仓 `.go` 非测试命中 6 行，**逐枚读**：`bridge.go:84` 是定义、`run.go:225`／`bridge.go:16`／`composer_handlers.go:37` 是注释 ⇒ **调用者 0**，且 `run.go:225` 那句 "Nothing calls it yet - the WebView2 event -> ParseComposerRequest hop does not exist in this tree" **未被本批推翻** | **成立** |
| "`Marshal()` 非测试调用者 0"（＝"全量字节仍无人接"那一格的一半） | `pump.go:190` 定义；调用点：`Snapshot()` 内部两处 `json.Marshal`（**不是 `Marshal()`**）、`Publish()` 里也是 `json.Marshal(snap)` ⇒ `Marshal()` 自己 **0 枚非测试调用者** | **成立** |
| "`lastPanelSnapshot()`／`snapshotCount()` 非测试调用者 0" | 两枚各自只有"注释行＋定义行"，**0 枚调用** | **成立** |
| "Go→面板通道仍不存在"（本程**另造**四条，不复用 §1.2/§2.1.1） | ① `http.Server`／`ListenAndServe` 非测试 **0**；② `go.mod` 里 webview 依赖 **0**（`grep -ci webview go.mod` → 0）；③ 全仓非测试的 `PostMessage` 命中 **只有 `internal/ball/` 那三枚 Win32 `PostMessageW`**（`sta_windows.go:142/:159`、`tray_windows.go:105`、`win32_windows.go:33` 声明）——**那是球的消息泵，不是 WebView2 的 host→page**；④ `text/event-stream` 命中 7 处全在**出站** LLM 适配器与 golden 夹具（`anthropic/adapter.go:152`、`openaichat/adapter.go:87`、`openairesponses/adapter.go:123`、`golden/golden.go:80/:218`、`cmd/llmrecord/main.go:95`），**没有一枚是服务端往页面推** | **成立**（另：`internal/proc/shutdown.go:55` 那枚 `"destroy-panel-webview"` 只是**步骤名册里一个还没有模块的名字**，实测 `wisp run`/resident 测试日志里它回的是 `shutdown step skipped (module not present)`） |

⇒ **给编排者改 `A248③` 用的那句话**：这两句**都现量成立**，方向相反、不互相抵消——
**"那份包在生产里被造出来并被驱动"＝已成立；"那份包能走到页面"＝仍不成立。**

### 4.2 `bridge.go` 那处注释的新指向——**两枚真存在、真在那两行**

| 指针对 | 现量（锚点副本 `sed -n '502p;533p'`） | 判 |
|---|---|---|
| `composer_test.go:502` | `func TestTheRendererHoldsExactlyOneDoorToTheHost(t *testing.T) {` | **真存在、行号逐字对** |
| `composer_test.go:533` | `func TestPlantedRendererDoorShapesGoRed(t *testing.T) {` | **真存在、行号逐字对** |
| 旧指向（那枚不存在的名字） | `grep -rn "TestComposerMethodNamesMatchFrontend" --include=*.go . \| wc -l` → **0**（⚠ 不加 `--include=*.go` 是 **16**，那 16 枚全在 `docs/evidence/**` 的历史自述里 ⇒ **口径必须写"Go 源码 0 枚"**，被验件 §3.3 那句"仍全仓 0 处"是**过宽**、但不是假） | 空指针已修、没留同型 |

### 4.3 §6.1 与 §6.3 抽查——**一枚算术不符，报名**

本程自造一把区分**同包裸调用**与**跨包带点调用**的尺（`(^|[^A-Za-z0-9_.])<sym>\(` 与 `panel\.<sym>\(`，
两向都去 `_test.go`）：

| 建造者（`5821e24` message 点名的那五枚） | 同包非测试调用点 改前→改后 | **跨包**非测试调用点 改前→改后 |
|---|---|---|
| `NewSnapshot` | 0 → 2（`pump.go:138`、`:182`） | 0 → **0** |
| `NewComposerState` | 0 → 1（`pump.go:171`） | 0 → **0** |
| `NewModeView` | **1 → 1**（只有 `composer.go:202`） | 0 → **0** |
| `ModeUnknownView` | 0 → 1（`pump.go:175`） | 0 → **0** |
| `UnsetWorkspaceView` | 2 → 3（新增 `pump.go:163`） | 0 → **1**（`cmd/wisp/panel_pump.go:84`） |

⇒ 锚点上"跨包新增的非测试调用点"**总共 4 枚**：`panel_pump.go:84`（`UnsetWorkspaceView`）、
`panel_pump.go:86`（`WorkspaceViewFromRoot`）、`run.go:420`（`NewStreamLog`）、`run.go:421`（`NewSnapshotPump`）。
后两枚是**新符号自己**、不是"某枚既有建造者 gained 包外调用者"。
⇒ **所以 §6.1 结论句"五枚里只有 2 枚真 gained 包外调用者"是算错的**：
**按它自己那张表，五枚里只有 1 枚（`UnsetWorkspaceView`）；第 2 枚 ✅ 打在了 `WorkspaceViewFromRoot` 上，
而那枚不在 commit message 点名的五枚里。** 表行的读数都对、**汇总句的分子/分母配错了**。记 **F-PUMP-5**（措辞级）。
⇒ 它对**一句话不成立**的复算（`NewModeView` 前后各 1 枚、调用点没动过）**本程独立复算为真**。

### 4.4 顺手量出一枚**没被任何一格登记**的零执行分支（同 §7 第 3 条那一族）

`37a4705` 的 commit message 写着触发点之一是 "`consoleSink.Publish` 的 **tool-start**/end/stuck/error/done"。
现量：`agent.EvToolStart` 在全仓 `.go` 里只有 **3 命中**——`internal/agent/sink.go:24`（注释）、
`:25`（常量定义）、`cmd/wisp/run.go:738`（`case` 消费），**`internal/agent` 里 0 枚 publish**
（同一把尺逐枚量：`EvToolEnd` 3 枚、`EvStuck` 2、`EvError` 2、`EvDone` 1、`EvTextDelta` 1、`EvControl` 1，
**只有 `EvToolStart` 是 0**）。
⇒ 这一批在 `run.go:741` 那行新加的 `changed = true` **今天零次执行**，且**没有任何仪器**管得住
"一枚 sink 事件种类从来没人发"这件事（它不是分支覆盖问题，是**词表里一枚死名**）。
⇒ 后果很小（那行只是"少推一次包"），但它和被验件 §7 第 3 条给 `panelSnapshotSummary` 那格的是
**同族**，而 §7 那张表**没有这一条**，commit message 还把它列成了活的触发点。记 **F-PUMP-6**。

---
