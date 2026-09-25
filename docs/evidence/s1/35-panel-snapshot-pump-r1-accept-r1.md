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
