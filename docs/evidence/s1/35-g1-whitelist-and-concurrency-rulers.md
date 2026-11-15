# `35-g1` — 票 35 两格补尺：`:42` 白名单模糊（50 名全拒＋全记日志＋在册 8 名逐一派发）与 `:44` 并发路由（10 路乱序配对）

- 程：`35-g1`（写腿，`cmd/wisp/**` 本轮唯一写手；`internal/panel/**` 本件**零写**，只读消费其导出面）
- 工序来源：`.scratch/wisp/probes/orch/20261115/txt/35-g1-rulers-dispatch.txt`（逐条照做）
- 侦察读数来源：`docs/evidence/s1/35-recon-20261115.md` 的 ① ② 两节（本件不重抄全文，只引判据）
- 写面：新增两枚测试文件（本件之外零改动）：
  - `cmd/wisp/panel_inbound_fuzz_35g1_test.go`
  - `cmd/wisp/panel_rpc_pairing_35g1_test.go`
  - 本件：`docs/evidence/s1/35-g1-whitelist-and-concurrency-rulers.md`
- 硬约束逐条：⛔ 未动 `thresholds.go`／golden／SLO／allowlist／`docs/PLAN.md`／`docs/specs/**`／D43／`frontend/**`／`design/**`；⛔ 未放宽任何既有断言；⛔ 未 push；⛔ 未改任何产码；两枚新测试均无新增产码协程（`observe.Registry.Spawn` 无新调用点，因为本波没造生产协程）。
- 勾选状态：本件**不勾票面任何 `- [ ]`**（勾归编排者与非实现者验收腿）。

---

## 一句话结论（先说最值钱的）

`：42` 补到票面口径：**50/50** 随机/禁止名经 `Handle` 全拒且**审计里恰 50 条、逐条带自己的方法名**；**在册 8/8** 逐一派发到**各自**处理器恰一次。`：44` **配对键不是 requestId，是 RPC 帧的整数 `id`**（`cmd/wisp/panel_control_windows.go:793` 的 `id := strconv.Itoa(d.ID)` → `:823` 的 `window._rpc[id].resolve(...)`）；已按这个**真实键**造出 10 路乱序配对尺 `TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1`，并用两条**故意错的配对规则**作反证（`byArrival`、`singleSlot` 都被判红），证明这把尺不是同步桩下的恒真尺。两枚新测试**定向现跑全绿**（见 §4 逐条读数）。

---

## 1. `：42` 白名单模糊 — 现状 → 补到票面口径

### 1.1 现状（侦察读数，我现读复核）

| 项 | file:line（今天真身） | 侦察读数 |
|---|---|---|
| 方法白名单守卫 | `internal/panel/bridge.go:172-174`（`if !knownComposerMethod(r.Method) { return ..., ErrComposerRequest }`） | 有牙 |
| 闭集名册（18 行钉，8 名） | `internal/panel/inbound_roster_253_test.go:430`（`TestFullInboundMethodRosterIsClosed`，`wantFullInboundRosterSize253 = 8` 在 `:68`） | 有牙 |
| 拒请求审计线 | `internal/panel/composer_dispatch.go:298`（`record` 的 `INBOUND-DISPATCH ... method=%q ...`） | 有牙 |
| 产码派发接线（七门） | `cmd/wisp/panel_inbound.go:384-403` | 已接 |
| **「在册逐一派发」覆盖** | 侦察 §①-c：仅 **3/8**（workspace／attachment／message 有断，mode／config.get／config.set／config.models.list／cost.summary 无"就是这台处理器"的断） | 缺口 |
| **最大一枚模糊尺** | 枚举 **11 名**且**直接调 `knownComposerMethod`**（绕开 `Handle`，不到审计线） | 缺口 |

⇒ 票面 `:42` 要求"50 随机/禁止名全拒＋全记日志 ⇒ 在册各枚逐一派发"，现状两半都没到。

### 1.2 我补的两枚（`cmd/wisp/panel_inbound_fuzz_35g1_test.go`）

**A. `TestWhitelistFuzz50NamesAllRefusedAndAudited35g1`**
- 表驱动 **恰 50 名**（`g1ForbiddenNames()`；落表先断言 `len == 50`，少一名即 `Fatalf`，防止悄悄降标）。
- 50 名覆盖五形：真页面/邻票名（`approval.decide`、`panel.approval.request`、`panel.resync`、`tasks.list`…）、在册名近邻（`panel.mode.request2`、`config.Get`、`cost.summary2`…）、异命名空间（`models.delete`、`grants.revoke`…）、形状噪声（`panel.`、`.hidden`、`a.b`…）、以及**故意保证没有一名 TrimSpace 后等于在册名**（否则守卫会放行、这把尺会撒谎——守卫先 `TrimSpace` 再比）。
- 每名都**经生产入口 `panel.ComposerDispatch.Handle`**（不是直接调守卫）：断言 ① 必错；② 错误**包 `panel.ErrComposerRequest`**（用 `isComposerRequestRefusal` 沿 `Unwrap` 走，非"任意错误"）；③ 回用户的句子**非空且带自己的 requestId**；④ mode 处理器被调 **0** 次；⑤ 审计 sink **恰 50 条**、第 i 条**含 `method=<该名>`** 且带 `INBOUND-DISPATCH`。
- **正控**：一个在册名（`panel.mode.request`，挂处理器）必须被受理且处理器恰 1 次，证明这把尺能区分合法与伪造、不是"Handle 永远报错"。

**B. `TestEveryRosterMethodDispatchesToItsOwnHandler35g1`**
- 对**在册 8 名**逐一 `Handle` 一次，七个处理器槽各挂 spy；断"**本槽 +1、其余六槽 +0**"（`config.get` 与 `config.set` 共用 config 槽，各记一次 delta）。
- 另断：8 名都仍被**运行中的**守卫受理（经 `panel.ParseComposerRequest` 问，不重实现守卫）。

**行为牙（A1880）**：
- 删/改 `bridge.go:172-174` 守卫 ⇒ A 的红（名落到派发或默认支，审计条数变化）。
- 删 `composer_dispatch.go` 任一 `case`（如 `MethodModelsList`）⇒ B 红（对应 spy 不涨、走 `rosterMismatch` 报错）。
- 从 `knownComposerMethod` 删一名 ⇒ 该名在 `ParseComposerRequest` 即被拒，B 红。
⇒ 都是"删/改被测**行为**会红"，不是字符串断言。

**定向现跑读数**（日志 `.scratch/wisp/probes/orch/20261115/g1-targeted-4.log`）：
```
=== RUN   TestWhitelistFuzz50NamesAllRefusedAndAudited35g1
    ... whitelist fuzz: 50/50 forbidden names refused by the ErrComposerRequest door, 50 audit line(s), each naming its method; positive control accepted
--- PASS: TestWhitelistFuzz50NamesAllRefusedAndAudited35g1 (0.00s)
=== RUN   TestEveryRosterMethodDispatchesToItsOwnHandler35g1
    ... roster dispatch: 8/8 roster names each reached their own handler exactly once and no other socket moved
--- PASS: TestEveryRosterMethodDispatchesToItsOwnHandler35g1 (0.00s)
```

---

## 2. `：44` 并发路由 — **配对键到底是什么（先把差异具名）**

### 2.1 判据（file:line，逐字）

| 事实 | 出处 |
|---|---|
| 回执按 **RPC 帧 `id`** 定位 | `cmd/wisp/panel_control_windows.go:793` `id := strconv.Itoa(d.ID)`；`:823` `c.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")` |
| 帧 `id` 由页面侧 `window._rpc` 序号章铸 | `:897-914`（`Bind` 注入的 stub：`var seq = RPC.nextSeq++;` 注册 `RPC[seq]`，再 `window.external.invoke(JSON.stringify({ id: seq, method: name, params: ... }))`） |
| 帧结构 | `:764-768` `ctlRPCMessage{ ID int `json:"id"`; Method string; Params []json.RawMessage }` |
| `requestId` **不在回执路径** | `internal/panel/bridge.go:179-182` 只要求它非空；`NewRequestID`（`bridge.go:120`）在本树**非测试调用点 = 0**。它是审计/日志关联键，不是配对键。 |

**结论（给第一行）**：**配对键是 RPC 帧的 `id`（即页面 `window._rpc` 的 `nextSeq` 序号），不是 `requestId`。**
- RPC `id`：`cmd/wisp/panel_control_windows.go:793`（取值）＋`:823`（落槽）。
- `requestId`：`internal/panel/bridge.go:133`（仅是封套字段）＋`:179-182`（仅校验非空），**没有任何回执按它落槽**。

⇒ 票面写"correlationId routing"是**措辞过期**：语料里 correlation 的实际载体是帧 `id`。**我没有假造一个 requestId 配对去迎合票面文字**（那是 dispatch 明令禁止的）；本尺按今天**真实的配对键 `id`** 造。

### 2.2 尺：`TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1`

**为什么是 headless（诚实边界，先声明）**：具体 `*panelControl` 的回执路径**没有可观测接缝**——`msgcb` 的回执就是交给 `c.Eval` 的一串字符串，而 `panelControl.Eval` 是 `*edge.Chromium` 的 COM 调用（字段未导出、无可注入 core）。加接缝＝**改产码**，本 dispatch 明令不做。改开**第二个真窗**又撞本包**已量到的**抖动：`cmd/wisp/panel_unbound_253r6_windows_test.go` 头注记"a second panel window created right after a previous one intermittently dies mid-create (measured … on the pristine tree)"。故本尺用：
- **页面侧**：真异步 `_rpc[seq]` **槽表**（`g1Page`）——正是侦察说现夹具缺的东西（同步桩会让任何配对恒真）；
- **宿主侧**：`g1Pipe.msgcb`，逐行照 `panel_control_windows.go:782-826` 的"已绑定门成功支"（含 `id := strconv.Itoa(d.ID)`、`_rpc[id]`），且**只排队不即时落槽**（把异步槽建出来）；
- **答案**：走**真的生产门** `panel.ComposerDispatch.Handle`（封套 `panel.bogus.<id>` 被拒，拒句自带 `rid-k<id>`，故每条答案自证身份）。

**断言**：10 个槽先全开（10 路在飞），10 帧交付，回执按**乱序**（`[7,2,9,1,10,4,6,3,8,5]`）落槽；断**每槽落到的答案只叫自己的 `rid-k<id>`、恰 settle 一次**。

**反证（让尺有牙，不是恒真）**：同一批 10 帧再用两条**故意错的配对规则**重放，**必须被判红**：
- `byArrival`（第 k 个回执塞进第 k 槽，忽略 id）⇒ **9/10 槽错**（被捕获）；
- `singleSlot`（全塞槽 1）⇒ **10/10 槽错**（被捕获）。
若这两条控制**没有**错，测试 `Fatal`（那把尺就是恒真的）。

**生产锚（A1880 注明是字面钉、不是行为牙）**：`assertProductionPairsByFrameID35g1` 读 `cmd/wisp/panel_control_windows.go`，若要害表达式不再是 `id := strconv.Itoa(d.ID)` 或 `window._rpc[" + id + "].resolve(` 即红——这是"模型↔产码"的系绳，**行为牙在模型＋两条控制**，此锚只防模型腐烂。

**定向现跑读数**（`g1-targeted-4.log`）：
```
=== RUN   TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1
    ... control "byArrival": 9/10 slots mismatched (caught, as required)
    ... control "singleSlot": 10/10 slots mismatched (caught, as required)
    ... 10-way out-of-order pairing: all 10 slots settled by frame id, each naming its own rid-k<id>; control byArrival and singleSlot both mismatched as required
--- PASS: TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1 (0.00s)
```

**诚实标注本格的缺口（不遮掩）**：本尺是 **模型尺**（页面槽表＋照抄产码的 `msgcb` 支），**不是真窗端到端**。真窗 `:44` e2e（真 WebView2 上 10 路 `window.wispDispatch` 乱序回执）**本波未造**，理由如上（无可观测接缝＋已知连窗抖动），属〔仅本机可量〕，留给后续腿或编排者开窗窗口——**这把尺不声称覆盖了真窗那一发**。

---

## 3. 我改了什么、没改什么

**改**：新增两枚测试文件（§1.2 / §2.2）＋本件。
**没改**：⛔ 任何产码（`cmd/wisp/**` 的非 `_test.go` 一字节未动；`internal/panel/**` 一字未动）；⛔ 任何既有测试；⛔ 票面；⛔ 台账；⛔ push。

---

## 4. 复跑方式（可核命令）

定向（快，二进制已在 `.scratch/wisp/stage-cli/`，DLL 已就地）：
```
cd cmd/wisp
D:\\work\\workspace\\projects\\plans\\Wisp\\.scratch\\wisp\\stage-cli\\wisp.test.exe -test.run=35g1 -test.v -test.count=1
```
或走载具（会重建并用 pin DLL 台件）：
```
D:\\work\\soft\\Git\\bin\\bash.exe scripts/wisp-cli-staged-tests.sh -test.run=35g1
```

---

## 5. 整包读数（A1882 三件套）

**发跑**：`scripts/wisp-cli-staged-tests.sh`（无参，整包 `./cmd/wisp`，pin DLL 就位、cwd=`cmd/wisp`），日志 `.scratch/wisp/probes/orch/20261115/g1-wholepkg-1.log`。

**我怎么停的（A1886，先具名）**：整包里一串"嵌套重编译/干净 checkout 构建"的用例极慢（`Test219AC2Mutation*` 会 fork 一个 `go test -overlay` 子进程把整包重编一遍；`TestCleanCheckoutBuilds_AC11` 会做干净 checkout + 构建），我用**自己的上限**停掉，不让它无界跑。停时读数如下——这**不是**"整包绿"，是一个**截断点**。

| A1882 三件 | 现量（`g1-wholepkg-1.log`，停在用户上限处） |
|---|---|
| **终点锚 `TestZZ315A1`** | **未命中**（该锚在文件序 index 670/673，本发没跑到） |
| **最后一条 `=== RUN`** | `TestRunBooksWithASnapshotOfItsLiveQueue`（文件序 index **327/673**） |
| **本波新测试在截断点前/后** | `:42` 两枚（index **278/279**）在**截断点前**，且**整包里已跑并 PASS**；`:44` 一枚（index **361**）在**截断点后**（327 < 361，未跑到）⇒ **整包对 `:44` 的覆盖＝0，`:44` 的绿只来自定向发**。 |

- 停时 PASS/FAIL：**PASS=326，FAIL=1**。唯一 FAIL＝`TestAC4FocusReturnToPriorWindowGap33r5`（4.96s，文件序 index **270**，早于我两枚 `:42` 测试的 278/279）；它**不是我这枚文件里的测试**（我两枚是 `*35g1`，本发里均 PASS），且本腿只**新增测试文件、零改产码**，按构造不可能是它变红的成因——发跑兄弟腿的派单把它列为**本机已知红**之一。**41 枚语音红在更后面**（voice 用例文件序更后），本发窗口内**尚未出现**——即本波窗口里没有拿"41 红"当借口，也没被它掩盖。
- **文件序权威读数**（`wisp.test.exe -test.list=.*`，共 **673** 枚）：
  - `TestWhitelistFuzz50NamesAllRefusedAndAudited35g1` = **278**
  - `TestEveryRosterMethodDispatchesToItsOwnHandler35g1` = **279**
  - `TestAC4PriorFocusSurvivesARefusedPanelSample`（已知 AV 崩溃点）= **360**
  - `TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1` = **361** ⇒ **紧挨 AV 崩溃点之后**
  - `TestZZ315A1RootLoopFirstBridgeCall`（终点锚）= **670**
  ⇒ 关键结论（A1882 的正题）：**`:44` 那枚恰好排在已知 AV 崩溃点的下一个**，真崩的那天它会落在截断点之后、整包对它覆盖为 0——这正是 A1882 说的"PASS=N 不讲边界"。所以 `:44` 的绿**必须**靠定向发来证明，本件就是按这个口径报的。

**本腿三枚在"截断点前"的读法（dispatch 要的就是这句）**：`:42` 两枚（278/279）**在截断点前全绿**（整包 PASS＋定向 PASS）；`:44` 一枚（361）**在截断点后**，整包未覆盖，**定向 PASS**（见 §1.2/§2.2 逐条日志行）。

