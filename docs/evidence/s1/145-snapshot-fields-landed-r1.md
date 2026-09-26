# 票 145 —— 快照字段落地（r1，实现程 · 代码那半）

> **本件的射程**：票 145 的 **AC#2／AC#3／AC#5／AC#6**，按 `.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md`
> 解冻的那一小块（`internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`）执行。
> AC#1 的对照表是另一枚只读程交的（`docs/evidence/s1/145-snapshot-field-census-r1.md`，锚 `88eab34`），本件在它上面**复量并更正**。
>
> **一句话结论（先说，免得被读成"做完了"）**：**落地集＝空。** 不是"没找到有源的字段"，是
> **有源的字段全部卡在四把锁上，而那四把锁的钥匙一枚都不在本程的放开面里**。
> 本程把这四把锁逐枚**量出来**（含一次仓外全量闭合实验：把字段真落下去、真填上活对象、真加断言，看它能不能绿），
> 并把**验证过的闭合补丁**交到 `.scratch/wisp/probes/145/patch/`（`git apply --check` 在 `e24ae28` 上干净）。
> **票面六格本程一格不自勾。**

---

## 0. 锚点 · 地界 · 零改动自证

### 0.1 锚点（进场现读，非抄派单）

```
$ git rev-parse --short HEAD        # 进场 09:4x
720cae6
$ git rev-parse --short HEAD        # 写本件时（10:08）
e24ae28
$ git branch --show-current
dev
```

锚点从 `720cae6` 漂到 `e24ae28` 是**共享树里别家在提交**，不是漂移。`720cae6..e24ae28` 现量 **17 枚**：
`f58c5136` `5d46f241` `bec17272` `4e169760` `eed7c998` `63e228fc` `9a8766e6` `10e3585c` `c0778696`
`6550dc42` `81673ddc` `e9ce4198` `dd2b277b` `6de3d1c5` `5429c0d4` `da8b6677` `e24ae28d`
（票 153 证据件／票 151 交件与收表／票 152 交件与探针／台账 A273–A279／三枚验收派单存档）。

**它们有没有碰本程那三枚文件——现量，不推理**：

```
$ for h in $(git log --format="%h" 720cae6..HEAD); do git show --name-only --format= $h \
      | grep -E "internal/panel/|cmd/wisp/panel_pump"; done
（无输出）
$ git status --porcelain -- internal/panel/composer.go internal/panel/pump.go cmd/wisp/panel_pump.go
（无输出 = 工作树与 HEAD 一致，且无人在写）
```

⇒ 那 17 枚落的是 `cmd/wisp/run.go`／`slo_windows.go`／`internal/tools/bridge.go`／`docs/**`／`.scratch/**`，
**与本程的放开面零重叠**（编排者在 `A273②` 也是这么划的；他那句"不碰 151 的 `run.go`"本件要更正，见 §1 P3）。

### 0.2 本程写了什么、没写什么

| | 读数 |
|---|---|
| 生产码改动 | **零枚文件的逻辑改动**。只有 `internal/panel/composer.go` 与 `internal/panel/pump.go` 的**头注释**各加一段**指针**（把本件量到的四枚钉、四把锁与"无生产者"那枚字段钉在载体上），**不改行为、不新增字段**——见 commit `eb38c97`，内容在 §5.5 与 §2 |
| 测试码改动 | **零**。放开面里没有 `*_test.go`，本程不假装能写（⇒ AC#6 的锁，见 §3.4） |
| `frontend/**`／`design/**` | **一个字没读、没写**。唯一例外＝按普查 §4.4 自己要求的取证，留了一行 `git status --porcelain -- frontend/src/lib/panel.ts`（读数：**空**，即那把尺读的是与 HEAD 一致的版本，见 §5.3） |
| 单独保留的三枚 | `tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go` — **零字节**，且本程没为它们中任何一枚的变绿做任何事 |
| 契约轴 | `docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／`thresholds.go`／golden／`allowlist.txt`／`rules_gateway.go` — 本程的全部操作是 `grep`／`sed -n` **读** |
| 台账／工单 | `docs/reports/pending-and-issues.md` 与 `.scratch/wisp/issues/145-…` **零字节**（勾与账归编排者） |
| push | **零**。三枚 commit 全在本地 `dev` |
| 临时件 | 仓外副本 `D:\tmp\wisp145-e1`（**只建不删**，闭合实验的现场）＋仓内原始读数 `.scratch/wisp/probes/145/`（E1–E4、基线、名册、补丁） |

### 0.3 尺（本件只认这几把）

| 尺 | 出处 | 用来判什么 |
|---|---|---|
| 载体 | `internal/panel/composer.go:44-49` | 有没有那一枚键 |
| 装配根 | `cmd/wisp/run.go:421-427`（`panel.NewSnapshotPump(panel.PumpSources{…})`） | 那一枚键**谁填** |
| 事件汇 | `cmd/wisp/run.go:751-785`（`consoleSink.Publish`） | `run.*`／`failures` 的数据今天落在哪 |
| 双向对账 | `internal/panel/composer_test.go:48`／`approval_test.go:105` | Go 加了键而 `panel.ts` 没跟上 ⇒ **谁响**（本件实跑，见 §2） |
| 键集钉 | `internal/panel/pump_test.go:111-124`／`:270-276` | 快照的 JSON 键集＝契约本身，第五枚键进来就红 |
| 屏词汇 | `frontend/src/lib/panel-views.ts`（已提交版）＋全树 grep | AC#3 ⓐ 的值有没有名字 |
| 单调时钟 | `internal/observe/clock.go`＋`internal/agent/approval/gate.go:257`＋AGENTS §1.2 | `remainingMs` 这类计时字段能不能落 |

⚠ **方法学一条**（承普查 §7.2）：本仓所有 commit 的 author 同名，**"我碰了什么"只能按自己的 commit hash 集 `--name-only` 现量**；
`720cae6..HEAD` 的区间 diff 里 `cmd/wisp/run.go`／`internal/tools/bridge.go` 各 1 枚，**都不是本程的**。

---

## 1. 派单与普查给的前提，逐枚复量（不照抄）

| # | 前提（出处） | 现量读数 | 判 |
|---|---|---|---|
| P1 | "`Snapshot` 恰四字段，`composer.go:44-49`；生产构造点 `:79`"（093x 特批 §1） | `:44-49` 确为四枚；`NewSnapshot` 在 `:63`、其 `return` 在 `:79` | ✅ 成立（但"构造点"这词有第二层，见 P2） |
| P2 | 同上，被当成"改一处构造点就填得上真值" | 真值不来自 `NewSnapshot`，来自**装配根的读口**：`cmd/wisp/run.go:421-427` 那枚 `PumpSources` 字面量（全树 `NewSnapshotPump(` 生产调用者**唯一一枚**＝此处）。`pump.go:182`／`panel_pump.go` 只是它的下游 | ⚠ 成立但**不足**：按 P2 的字面做，加出来的键今天**没人填**（实测见 §2.3） |
| P3 | 编排者 `A273②`："四组两两不重叠……**不碰 151 的 `run.go`**" | 候选落地集 5 枚段里 **4 枚的读口必须出现在 `run.go:421-427`**（`run`／`tools`／`cost`／`failures`；其中 `run.usage` 还要 `internal/agent/loop.go` 先转发 `EvUsage`，那是票 153 的地界）。只有 `approval.*` 一枚能走已接线的 `Verdicts` 读口旁通（本程在仓外副本真做通了，§2.4） | ❌ **不成立**：AC#2 的"填真值"这一半结构性需要 `run.go`。⇒ **报回**，不硬改 |
| P4 | 093x 特批 §2：`加字段＝扩契约` 的人工批准"只到新增字段为止" | 本程**接受为授权**，并严格按字面用：不加字段＝本程不落键；**未**据此动任何现有字段语义、未删字段、未动 `C17` 白名单条目名/类型 | ✅ 授权在效（它解的是契约批准那把锁，没解 §3 那三把） |
| P5 | 普查 §4.4："Go 加键而 `panel.ts` 没跟上 ⇒ **那把尺**会红"（单数，指 `composer_test.go:73-78`） | 实跑（仓外副本 E1）：响的是 **4 枚**用例，不是 1 枚——`approval_test.go:129`、`composer_test.go:74`、`pump_test.go:124`、`pump_test.go:276`（后两枚是票 35 的**键集钉**，普查没计） | ⚠ 成立但**数少了**：闭合的爆炸半径是 4 枚测试＋1 枚前端文件（§2.1/§2.2） |
| P6 | 普查 §2.8 行 8："`approval` 那一组**全有源，且生产者今天已经在跑**"，含倒计时 `remainingMs` | 逐枚分：**`windowMs` 有源**（`gate.Window()`，实测 `2000ms`）；**`vetoChannels` 有源**（`ChannelRegistry.Statuses()`，实测四行含 B1 原句）；**`remainingMs` 无源**——`approval.EventTick`（`ui.go:64`）**声明了、全仓零生产者**，真发的三种 Event（`gate.go:274`／`:387`／`:510`）带的 `Remaining` 是**静态窗口长／预警提前量**，不是倒计时；倒计时本体是 `gate.go:257` 的 `g.clock.After(g.window)`，**没有"还剩多少"的读口** | ❌ **部分不成立**：按普查自己的 §0.4 判据①（"只有类型没有生产者，不算 ⓐᐟ 成立"），`remainingMs` 应从**甲组挪进乙组**，与 `thinkingMs`/`reasoningMs`/`durationMs` 同族（另开票装计时器） |
| P7 | 普查 §3.2：ⓐ"消费端今天就位，Go 加一枚 `view` 前端一行都不用动" | 已提交版 `panel-views.ts` 的 `currentView()` 确读 `snapshot.view`；而 Go 侧九枚屏 id 现量**零命中**（`grep -rn --include=*.go` 限 `internal/panel/`＋`cmd/wisp/`、剔 `_test.go` → 无输出） | ✅ 成立，且**加重**：ⓐ 落地的话值只能由宿主随便给（§4.1） |
| P8 | 派单共同纪律："判跑没跑到只认 `=== RUN` 枚数" | 本程每一发门禁都带 `=== RUN` 计数（§5）：仓内基线 `internal/panel` RUN=105、`cmd/wisp` RUN=139；仓外副本各发同样带计数 | ✅ 照做 |
| P9 | "同树三枚别家在跑" | 现量：票 151 **已交件**（`63e228fc` 收表）⇒ `run.go` 现在**没人写**，但**仍不在本程放开面**（地界由特批划，不由"没人抢"划） | ✅ 成立，附一条更新 |
| P10 | 记忆里那条 owner 排序信号"前端那儿就等后端完事儿了" | 与实测冲突的一面：**前端那一侧今天缺的不是这五枚键**——五枚键落下去也仍要 `panel.ts` 同批＋泵把包送到页（最后一公里在票 33/35 名下，`pump.go:15-23` 自证"无 Go→页通道"） | ⚠ 不当授权用；登记以免下一位以为"落了字段前端就动了" |

⇒ **两条要报回的**：P3（`A273②` 那句"不碰 run.go"对 AC#2 不成立）与 P6（`remainingMs` 无生产者）。
本件把它们写成读数与出处，**不据它们硬改，也不据它们扩权**。

---

## 2. "加字段谁会响"——AC#4③ 要的那一发，实跑（普查那发是推论，本程补成读数）

**现场＝仓外副本 `D:\tmp\wisp145-e1`**（整树拷贝，剔 `.git`/`node_modules`/`.scratch`/`design`/`third_party`；
副本里怎么改都碰不到共享树，副本**只建不删**）。候选落地集＝`run` / `tools` / `approval` / `cost` / `failures`
五枚段、9 枚键、8 枚 view 类型（定义见 `.scratch/wisp/probes/145/patch/`）。

原始读数：`probes/145/copy-E0-baseline-rulers.txt`、`wisp145-e1-E1-panel.txt`、`wisp145-e1-E2-panel.txt`、
`wisp145-E3-cmdwisp.txt`、`wisp145-E3b-cmdwisp.txt`、`E4-removal-check.txt`。

### 2.1 E0 副本基线（两把尺绿，逐键数＝4）

```
approval_test.go:134:  Snapshot <-> PanelSnapshot: 4 JSON keys reconciled
composer_test.go:79:   Snapshot <-> PanelSnapshot: 4 JSON keys reconciled
```

### 2.2 E1 只加 Go 侧九枚键、`panel.ts` 不动 ⇒ **响四枚，不是一枚**

`go test -count=1 -v ./internal/panel/` → `RUN=105`、`FAIL=4`＋1 枚副本产物红：

| 用例 | 红句（原文） |
|---|---|
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `approval_test.go:129: Go Snapshot emits [run tools approval cost failures] that interface PanelSnapshot does not declare` |
| `TestComposerContractTypesMatchFrontend` | `composer_test.go:74:` 同一句 |
| `TestThePumpBuildsThePacketFromWhatTheHostHolds` | `pump_test.go:124: snapshot JSON keys = [pending results run tools approval cost failures composer generatedAt], want exactly the four PanelSnapshot declares` |
| `TestPublishHandsTheBytesToTheAttachedExit` | `pump_test.go:276: exit bytes carry keys […], want the four PanelSnapshot declares` |

⚠ 第 3、4 枚是**票 35 的键集钉**（`pump_test.go:111-112` 原话："the four keys are the contract, so a fifth arriving here is
a contract change dressed as a bug fix"）。普查 §4.4 只点了前两家 ⇒ **P5 的更正**：闭合的爆炸半径是 **4 枚测试＋1 枚前端文件**，
且其中 `pump_test.go` 与两把尺**都不在本程放开面**。
⚠ 同发里 `TestC21DesignTokensFourWayAgree` 也红，但**那是副本产物**（`read design/assets/tokens.css: open D:\tmp\wisp145-e1\design\…`
——本程拷树时剔了 `design/`），**不计入读数**；仓内那一枚另有其既有红因（`Q-52` 撤回令，见 §5.2）。

### 2.3 E2 把 `panel.ts` 逐键补上（九枚键＋7 枚 interface）⇒ 两把尺转绿，两枚键集钉仍红

```
composer_test.go:79:  Snapshot <-> PanelSnapshot: 9 JSON keys reconciled   --- PASS
approval_test.go:134: Snapshot <-> PanelSnapshot: 9 JSON keys reconciled   --- PASS
pump_test.go:124:     snapshot JSON keys = […9 枚…], want exactly the four  --- FAIL
```
⇒ 键集钉要改的是**两行硬编码字符串**（`pump_test.go:123`／`:275`）。**这是"改断言"，不是"放宽断言"**——
它跟着契约走、方向与契约同批；但按字面它仍属"动别人的测试"，本程没在仓内做。

### 2.4 E1b **只加键、不动装配根**——那九枚键今天会带什么出门（本件最硬的一发）

副本里 `cmd/wisp/run.go:421-427` **原封不动**（就是仓内现状那五枚读口），把包构造出来看新段：

```
run      = {"reasoning":"","status":"","usage":{"input":0,"output":0,"cachedRead":0}}
tools    = null
approval = {"depth":0,"windowMs":0,"vetoChannels":null}
cost     = {"tokensIn":0,"tokensOut":0,"cachedTokens":0,"micros":0,"currency":""}
failures = null
pending len=1 results len=1 (both live)
```

⇒ **五枚段零枚被填**，而同一发包里 `pending`/`results` 是活的。也就是说：**今天往 `Snapshot` 上加字段，
加出来的正是本票自己要防的那形——一串永远为 0/空的常量**（普查 §2.0(2) 的第四态"无生产者"，只不过这次是"有键无填"）。
这一发是 AC#6 的反证，**不需要任何推测**：跑的是生产装配根那一份读口清单。

### 2.5 E3 全量闭合（副本内）：**能绿，且能真填**——但要用到四张本程没有的钥匙

副本里同时动：`composer.go`（键）＋`pump.go`（`PumpSources.Approval` 读口＋装配）＋
`cmd/wisp/panel_pump.go`（`liveGateState()` 读 `rt.gate.Window()`／`rt.gate.Channels().Statuses()`）＋
`cmd/wisp/run.go`（字面量加一行 `Approval: rt.liveGateState,`）＋`frontend/src/lib/panel.ts`（逐键）＋
`internal/panel/pump_test.go`（两行键集）＋**一枚新用例** `cmd/wisp/panel_pump_gate145_test.go`（AC#6 那问的答句）。
跑生产装配根（`go test -run TestSnapshotGateSectionComesFromTheLiveGate ./cmd/wisp/`，`RUN=1`、rc=0）：

```
packet approval section = {"depth":0,"windowMs":2000,
  "vetoChannels":[{"name":"ball","text":"单击悬浮球：悬浮球取消不可用","loaded":false},
                  {"name":"esc","text":"按 Esc 键：Esc 取消不可用","loaded":false},
                  {"name":"panel","text":"面板拒绝：面板取消不可用（票 37 未接入）","loaded":false},
                  {"name":"kws","text":"说取消词：语音取消不可用","loaded":false}]}
--- PASS
```
⇒ 值**来自这个进程真在跑的那两枚对象**（gate 的窗口、通道注册表的原句），B1 那句"语音取消不可用"逐字抵达；
`./internal/panel/` 同发除 `design/` 副本产物外**全绿**。

### 2.6 E4 承重判据（把 E3 那枚填值摘掉）

副本里删 `pump.go` 那一行 `snap.Approval = approval`（其余一字不动）：

```
--- FAIL: TestSnapshotGateSectionComesFromTheLiveGate
    windowMs = 0, want the gate's own 2000
    windowMs = 0: the section arrived unfilled
    packet carries 0 channel rows, registry reports 4
```
⇒ **摘掉填值那一行，断言就不响**——这一枚段的"承重"是量出来的，不是声称的。
按同一把尺，其余四枚段（`run`/`tools`/`cost`/`failures`）今天**连承重都无从谈起**：没有读口、没有值、没有用例（§3）。

### 2.7 顺手捉到一枚**仪器缺口**（登记，本程未使用）

`jsonKeysOf`（`internal/panel/approval_test.go:174-191`）读的是 `json:` 标签：
**没有标签的导出字段它看不见**，而 `encoding/json` 会按 Go 字段名把它发出去。实测（副本）：
`Snapshot` 加一枚 `LeakProbe string`（不带标签）——

```
probe: encoding/json emits LeakProbe = true ; jsonKeysOf(Snapshot) = [pending results composer generatedAt run tools approval cost failures]
TestComposerContractTypesMatchFrontend            --- PASS   ("9 JSON keys reconciled")
TestApprovalCardViewJSONKeysMatchFrontendTypes    --- PASS
TestThePumpBuildsThePacketFromWhatTheHostHolds    --- FAIL   (键集钉抓到了 LeakProbe)
```

⇒ 那两把"契约尺"**看不见无标签导出字段**，唯一的兜底是票 35 的字节级键集钉。
本程**没有**用这条口子塞任何字段（用＝绕开契约尺，正是本仓定的假绿形），只登记：
**将来谁给 `Snapshot` 加字段而忘了写标签，两把尺会照样绿**。归口＝编排者决定要不要为它补一枚钉（`AGENTS §1.1` 的"改断言"射程外，属新增判据）。

---

## 3. 逐字段判（AC#2 的落地集 ＋ AC#6 的反问）

**口径**：一列"真源"只认两种读数（① 该包今天就能算出这个值；② 今天已有人在生产它）——普查 §0.4 原尺；
本程自己实跑/现量的标 **[量]**，沿用普查未复量的标 **[普查]**。
"锁"＝落地还缺谁：**B**=`frontend/src/lib/panel.ts`（Q-51 未答，本程禁写）· **C**=`pump_test.go` 两行键集＋两把尺（禁面）·
**D**=`cmd/wisp/run.go` 装配根/事件汇（票 151 地界，特批未放开）· **D′**=`internal/agent/loop.go`（票 153 地界）·
**E**=一枚 `*_test.go` 写 AC#6 的断言（**放开面里根本没有测试文件**）。

| 字段（候选） | 真源 | 生产者可达？ | 锁 | **AC#6 答句**（哪枚用例断言它来自真源） | 判 |
|---|---|---|---|---|---|
| `approval.depth` | **[量]** 就是本泵刚建出来的 `pending` 卡片数 | ✅ `pump.go` 内派生 | B·C·E | 副本里能答（`depth` 与队列实测）；**但它是 `len(pending)` 的复述**，`App.tsx:75-76` 今天已经在自己算 | **不落**：装饰（同一枚事实的第二份拷贝），且缺 B/C/E |
| `approval.windowMs` | **[量]** `gate.Window()`（实测 `2000ms`） | ✅ 已在跑；填它要 `run.go:426` 加一枚读口（或旁通 `Verdicts`） | B·C·D·E | 副本 E3 那枚用例＝答句；仓内**写不出**（E） | **暂不落**（唯一实质阻塞是 B/C/D/E，不是数据） |
| `approval.vetoChannels` | **[量]** `ChannelRegistry.Statuses()`，四行含 B1 原句 | ✅ 同上 | B·C·D·E | 同上 | **暂不落** |
| `approval.remainingMs` | ❗ **无源**（P6）：`EventTick` 零生产者，倒计时是 `gate.go:257 clock.After`，无"还剩多少"读口 | ✗ | — | 答不出：任何实现都只能填静态窗口长 | **禁入**（挪进乙组；另开票装单调计时器） |
| `run.reasoning` | **[量]** `loop.go:545` 真转 `EvReasoningDelta`；`run.go:759-765` 收到并打印，**注释自己写着**"the reasoning field is ticket 145 row 3, which is a key the snapshot does not have" | ✅ 到得了 sink，sink 不留存 | B·C·D·E | 写得出用例，但要 D＋E | **暂不落** |
| `run.usage`（终值） | **[量]** `EvDone` 已带 `TokensIn/TokensOut`（`loop.go:937` 实填） | ✅ 到得了 sink | B·C·D·E | 同 E3 形状 | **暂不落** |
| `run.usage`（实时） | **[量]** `llm.EvUsage` 三适配器真发，但 `agent` 侧 `grep EvUsage` **0 命中**——`loop.go:543-546` 的 switch 只转 text/reasoning | ✗ 今日无转发 | D′ 先行 | 答不出（值根本不到 agent 面） | **禁入本票**（属 `internal/agent/**`，票 153 地界） |
| `run.status`（取消/完成） | **[量]** `EvDone` 只填 `Kind/TaskID/Text/TokensIn/TokensOut`，**不填 `Status`/`Stop`**；`Result.Status` 只在同步返回里 | 半：sink 那侧要记 | B·C·D（＋取向：改 `loop.go:937` 还是从 `Result` 组包） | 答不出，取向未定 | **暂不落**，取向摆编排者（普查 §2.14 同判） |
| `run.phase` | ❗ **[量]** `agent.Event` → `statemachine.State` 的映射表**在仓里不存在**（本程在 `internal/panel/`＋`cmd/wisp/` 剔测试 grep 状态名：零命中） | ✗ | — | 答不出：只能填常量或猜 | **禁入**（第四态"无生产者"，另开票） |
| `tools[]`（除 `durationMs`/`iconClass`） | **[普查]** `memory.ToolCall` 的 name/outcome/riskLevel/argsJson/correlationId，读口 `Store.ListToolCallsByTask`（本程现量在 `dao_toolcall.go:111`）；taskID 今天只活在 `run.go` 的 sink（`e.TaskID`）与 `StreamLog` 的键里 | 半：`rt.store` 可达，**taskID 未随包带出** | B·C·D·E | 写得出用例，但要 D＋E | **暂不落** |
| `cost`（`micros` 必与 `currency` 成对） | **[普查]** `agent.Cost`／`llm.Usage`／`Result.CostMicros/Currency`；历史读口 `task_log`（本程现量 `memory/models.go:65-68` 确有 `CostTokensIn/Out/Currency`） | 半，同 `tools[]` | B·C·D·E ＋ **单位口径未定案**（`cost.go:16-20` 自陈 micro-USD vs CNY） | 答不出：单位没定 | **暂不落**，且**绝不许面板自己加"￥"** |
| `failures[]`（除 `humanText`） | **[量]** `EvError` 到 sink（`run.go:775-778` 打的正是 `e.Err.Class`/`e.Err.Detail`），sink 不留存 | ✗ 未记录 | B·C·D·E | 写得出用例，但要 D＋E | **暂不落** |
| `view`（AC#3 ⓐ） | ❗ **[量]** Go 侧九枚屏 id 零命中（§1 P7）；消费端 `currentView()` 今天就位 | ✗ 值无源、名字无校验 | B·C·D·E ＋ 要先在 Go 立枚举＋双向尺 | 答不出 | **禁入**（详见 §4） |
| `thinkingMs` / `reasoningMs` / `durationMs` / `humanText` / `fragment` / `iconClass` | ❗ 普查乙组六枚，本程**逐条复查未推翻任何一条**（`tool_call.started_at` 零写者、class→中文文案映射不存在、`risk.Hit` 包外零消费者、图标名册冻结…） | ✗ | — | 答不出 | **禁入**（宁缺毋造） |

### 3.4 一句话总结这张表

**没有一枚字段是"只差一个键"**：每一枚实质候选都同时卡在 **B（前端契约面）＋ C（三枚禁面测试）＋ D（票 151 的 `run.go`）＋ E（放开面里没有测试文件）** 上。
⇒ **落地集＝空**，且这是**量出来的空**（E1b 那发证明"只加键"产出的就是装饰；E4 那发证明"能承重"的形今天需要四张钥匙）。
按票面骨头"宁缺毋造"，**空的落地集优于十二枚常量**。

⇒ **AC#6 逐枚答句**在上面最后一列：**一枚都没答出来**，因为答它需要的那枚用例**写在本程写不了的地方**。
这不是措辞问题——票 145 的 AC#6 与本程的放开面**互斥**，见 §8 报回。

---

## 4. AC#3："当前屏"这一维，两种角色一次定死（ⓐ 判不落，ⓑ **停手上报**）

### 4.1 ⓐ Go→面板（快照带"现在该显示哪一屏"）——**不落**，且本程把它判成"比 ❗ 更坏"那一形

本程自己现量三发（不沿用普查）：

```
$ grep -rn --include=*.go "palette|privacy|"tasks"|"cost"|"security"" internal/panel/ cmd/wisp/  (剔 _test.go)
（无输出）
$ grep -rn "panel\.view" --include=*.go .
（无输出）
$ sed -n '/func knownComposerMethod/,/^}/p' internal/panel/bridge.go
case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:   ← 四枚，无第五枚
```

⇒ **Go 全仓不知道那九枚屏的名字**：`view` 落进 `Snapshot` 就是"一枚无校验的自由字符串进入 C17 对外形状"
（拼错的表现是前端 `currentView()` 的 `console.error` ＋ 回落 chat，**不是一枚测试红**）。
而它同样要吃 §3 那四把锁（B/C/D/E）——**ⓐ 不是"顺手的一枚键"**。
⇒ **判**：ⓐ 要落，必须**同时在 Go 侧立一枚具名枚举＋与 `panel-views.ts` 双向钉住**（普查 §3.2 同一建议），
且那一步属**契约形状设计**，不属本程那三枚文件能自决的范围 ⇒ **不落，交编排者定形**。

**那句硬答案（普查 §3.4 的判，本程各补一发读数，原样保留）**：

1. **快照今天到不了页**：`internal/panel/pump.go:15-23` 自己写着"树里没有 Go→页通道（票 33 未领）、
   最后一公里不在这里"；§2.4 那发又量到"装配根不动 ⇒ 新段全为 0"。⇒ 加 ⓐ ＝ **加一枚永远读不到值的键**。
2. **ⓐ 单独存在换不了屏**：它只让宿主有权指屏，用户点击要的是 ⓑ 那条回路（§4.2）。
   ⇒ 下游任何读数**不许**写"换屏已通"；本程连 ⓐ 都没落，这句只作**判词**留存。

### 4.2 ⓑ 面板→Go（用户点导航条请 Go 换屏）——**停手上报，本程零动作**

- 批准射程：`Q-50` 甲批的是**显示与发起请求**；回写"当前屏"是一枚**新的入站方法**。
- 它撞的闸门具名：`AGENTS.md §2` 未定义即停清单里那条 **"C24 GojaHostAPI 初始集与 C17 方法白名单定稿"**。
- 现量它今天**确实不在名单上**（§4.1 第三发：`knownComposerMethod` 只答四枚），且门**已经在响过**：
  `internal/panel/composer_test.go:522` —— `"the renderer names a route the Go side does not answer"`（本程在锚 `e24ae28` 现读到该行仍在）。
- 编排者自己也已把它排在"动手前要先摆给 owner"那一格（`A273 next=④`：票 114 那根"网页事件→Go"的线，
  **`Q-50` 甲只批了显示与发起请求，动手前要先摆给 owner**）。

⇒ **本程动作＝无**：不加方法、不改白名单、不改 `bridge.go`（它也不在放开面）、不"先落 ⓐ 等 ⓑ 再说"。
⇒ **报回内容**：ⓐ 与 ⓑ 是两件事，**ⓑ 不是泵能顺带解决的**（泵能送达 ⓐ，泵造不出被 `f1cdafa` 删掉的那条路由）。
要 ⓑ 就要 `C17` 白名单的人工批准——那是 owner 的一句话，不是本程的一次判断。

---

## 5. AC#5 门禁：四数、名册差集、三件仪器（改前／改后各一次）

⚠ **基线用哪支，先说清**（票面 AC#5 点名要写）：本程**没走** `scripts/wisp-cli-tests.sh`——它是
`portable-tests.sh` 外面一层 PATH staging，而本票要的是**逐名名册**。本程用同一条 DLL 注入直跑：
`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./cmd/wisp/`（**MSYS 形路径，不是 `D:/…` 正斜杠形**）。
每一发都带 `=== RUN` 计数，**只认枚数判"跑没跑到"**：本程第一次取 `cmd/wisp` 基线时把 `PATH` 弄坏（`unset PATH`），
那发 `rc=127`、`=== RUN` 无从谈起 ⇒ **作废、不计入任何读数**，重跑的那发才算（读数在下面）。

### 5.1 四数与名册差集

| 包 | 时刻／版本 | rc | `=== RUN` | PASS | FAIL | SKIP |
|---|---|---|---|---|---|---|
| `./internal/panel/` | 改前 `e24ae28` 10:08 | 1 | **105** | 58 | 1 | 0 |
| `./internal/panel/` | 改后 `eb38c97` 10:1x | 1 | **105** | 58 | 1 | 0 |
| `./cmd/wisp/` | 改前 `e24ae28` 10:0x（DLL PATH） | 0 | **139** | 79 | 0 | 0 |
| `./cmd/wisp/` | 改后 `eb38c97` 10:2x（DLL PATH） | 0 | **139** | 79 | 0 | 0 |

**名册差集**（`--- (PASS|FAIL|SKIP): 名字` 排序后比集）：

```
$ diff probes/145/roster2-panel-before.txt probes/145/roster-after-panel.txt
（无输出 = 差集为空）
$ diff probes/145/roster2-cmdwisp-before.txt probes/145/roster-after-cmdwisp.txt
（无输出 = 差集为空，139 行名册逐名相同）
```

⇒ **差集为空**是本程**期望的**结果：唯一的生产码改动是两枚文件的**注释**（`eb38c97`，27 行 `+`、**0 行 `-`**），
零行为改动、零新字段 ⇒ 门禁颜色与名册都不该动。**"改后一样"在本票不是"没做活"的证据，是"没越界"的证据**——
做了活的读数在 §2（全部来自仓外副本，那是唯一能既真落字段又不脏共享树的形状）。

### 5.2 那一枚红的归因（不是本程的债，本程一字未动）

`TestC21DesignTokensFourWayAgree`（在**单独保留**的 `internal/panel/tokens_fourway_test.go` 里）。
今天这一发日志里它只暴露**一枚**红因：

```
tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css:
    The system cannot find the path specified.
```
⇒ `readRepoFile` 是 `t.Fatalf` 形（`:80-86`），第一枚文件读不到就**终止**，后面的四方比对**根本没跑**。
所以"这枚红今天只有一条可见病因、且不能据此说另一条不存在"——记忆里那条"现有两枚红因（放回文件也不会绿）"
本程**未复量、未推翻**，只登记本程看到的是哪一发。`Q-52` 撤回令在效 ⇒ **本程没为它变绿做任何事**（没还原那 16 枚删除、没换基准、没注释、没 Skip）。

### 5.3 那把尺读工作树这一形，今天不成立（普查 §4.4 末提醒的那条，本程留证）

```
$ git status --porcelain -- frontend/src/lib/panel.ts
（无输出）
```
⇒ 改前改后两发的双向对账读的都是**与 HEAD 一致**的 `panel.ts`——**没有**"对着别人未提交的版本报绿"那一形。
（另：本程任何"零命中"宣称**不含** `frontend/**` 与 `design/**` 的工作树状态，那半棵树归别人。）

### 5.4 三件仪器

| 仪器 | 读数 |
|---|---|
| `gofmt -l internal/panel/ cmd/wisp/` | **空** |
| `go vet ./internal/panel/ ./cmd/wisp/` | **rc=0、无输出** |
| `sh scripts/d22scan.sh` | **rc=0**；分母现量：production Go `internal/=205`＋`cmd/=23`（合计 examined 228）、ban #6 `frontend/=73`、ban #7 `internal/tools/=18`、ban #8 `design/=39`／`frontend/=73`／`internal/=413`／`cmd/=44`；`skipped as git-ignored: 1 file(s) under frontend/dist/assets/`。全文在 `probes/145/d22scan-after.txt` |

⚠ `tools/d22scan` 是独立 module，本程按 `scripts/d22scan.sh` 跑，**没有**在根目录 `go vet ./tools/d22scan/`。

### 5.5 契约轴零命中——只认本程自己的 commit 集

```
$ for h in ef07ba5 c655458 972ceba eb38c97; do git show --name-only --format= $h; done | sort -u
.scratch/wisp/probes/145/…      （E1-E4 读数、两包基线与名册、d22scan 全文、闭合补丁）
docs/evidence/s1/145-snapshot-fields-landed-r1.md
internal/panel/composer.go      （只有注释：13 行 +、0 行 -）
internal/panel/pump.go          （只有注释：14 行 +、0 行 -）
```
⇒ 枚数会随本件后面那两枚交件 commit 增长，但**路径集不变**（只有 `docs/evidence/s1/` 这一枚＋`probes/145/`）；
收口时按 `git log --author` 不可分（全同名），**按 hash 集现量**是本仓唯一可核形。
⇒ `docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／`thresholds.go`／golden／`allowlist.txt`／
`rules_gateway.go`／`frontend/**`／`design/**`／`pending-and-issues.md`／工单本体 —— **一枚都不在本程的 commit 里**。
单独保留的三枚测试文件同样零命中（`git show --name-only` 里不存在）。**未 push**（4 枚全在本地 `dev`）。

### 5.6 合并态复算（本程交件之后、别家又提了几枚，再跑一次同一把尺）

本程交完之后同树又落了别家的 commit（含 `5c28b3b` 前端 react-bits 落地、`777d6cc` 票 153 验收、
`091390c` 之后本程那枚探针件——**都非本程所写**）。在合并态 `777d6cc`／`60c47cf` 上按 §5.1 **同一把尺**重跑：

| 包 | 时刻／版本 | rc | `=== RUN` | PASS | FAIL | 名册与我 §5.1"改后"那一发 |
|---|---|---|---|---|---|---|
| `./internal/panel/` | `777d6cc` | 1 | 105 | 58 | 1（`TestC21DesignTokensFourWayAgree`，同因） | **逐名相同** |
| `./cmd/wisp/` | `60c47cf`（DLL PATH） | 0 | 139 | 79 | 0 | **逐名相同** |

⇒ 本程那两枚注释件在合并态没改任何颜色，也没被别家的改动盖掉。

### 5.7 一次**实际发生的**共享 index 险情（记下来，因为规矩救住了）

本程补探针名册那一枚 commit（`091390c`）之前，`git diff --cached --name-only` 里同时出现
**票 153 验收程的 28 枚 staged 文件**（`.scratch/wisp/probes/153/accept-r1/*` ＋ 它的证据件）——
那是别人 staged 在**共享 index** 里的活，与本程无关。
⇒ 本程按规矩做对两件事：① commit **带显式 pathspec 只指自己那两枚路径**，
实测 `git show --name-only 091390c` **只含** `roster2-{panel,cmdwisp}-before.txt` 两枚；
② **未 unstage、未动别人的暂存态**——那 28 枚随后由 153 验收程自己提交（`777d6cc`）。
这正是 `A272①` 那族的**第三次实发**：**风险不在 `git add -A`，在"别人已经 staged 在同一个 index 里"**，
而"我带了 pathspec"本身不是豁免——**每次 commit 前现量名册**才是那堵墙。本程七枚 commit 的路径并集
＝`probes/145/*` ＋ 本证据件 ＋ `internal/panel/composer.go` ＋ `internal/panel/pump.go`，**别家一枚都没有**。

---

## 6. 本程**没**核什么（不假装核过）

| # | 没核的事 | 为什么 | 影响读法 |
|---|---|---|---|
| N1 | **"字段抵达界面"这一整条**：没有任何用户可见变化 | 最后一公里不存在（`pump.go:15-23` 自证无 Go→页通道；票 33/35 未接），且本程连字段都没落 | 本件**不许**被读成"界面那条腿有东西可读了"——今天仍然没有 |
| N2 | **真机 `wisp run`** 一发没跑（没起真进程、没喂真 provider） | 副本 E3 走的是 `cmd/wisp` 测试里的装配根（fixture → `assembleRuntime` → 真 `Publish`），不是真 CLI | E3 那发是"生产装配根真填出了值"，**不是**"端到端跑过一次真任务" |
| N3 | 像素／视觉／截图 **0 次** | 本程无界面动作 | 行 7 那枚 2px 横条、行 12 的 `tabular-nums` 观感，本程一个字没判 |
| N4 | `cost` 的**单位口径**（micro-USD vs CNY）归口 | `cost.go:16-20` 自陈未定案，那份 ticket report 本程没去找 | §3 里 `cost` 那行"暂不落"含这一条，但**说不出该哪一切片定** |
| N5 | `DEFERRED(kws-veto, B1)` 与 `SPEC-12 §5` 登记表的 1:1 双向 | 那是 AGENTS §1.1 的独立约束；本程只引用代码里那枚标记 | 行 8 的"通道能力没有、显示有源"这句不依赖登记表是否对得上 |
| N6 | **单独保留的三枚测试文件的判据内容**未读 | 派单硬约束 | §5.2 只引用日志回显那一行 |
| N7 | 普查乙组其余 ❗（`thinkingMs`/`reasoningMs`/`durationMs`/`humanText`/`fragment`/`iconClass`）**只复查未推翻，未逐枚重新取证** | 本程新立的是 `remainingMs` 那一枚（P6，有实测） | 那六枚仍按普查的判语用；谁要落仍需自己复量 |
| N8 | 全树门禁／CI 颜色一枚未取 | 推送归编排者，且 `slo-full` 在本机自启会抢 CPU | §5 的四数**只覆盖** `./internal/panel/` 与 `./cmd/wisp/` 两枚包 |
| N9 | `tools[]`／`cost` 的 store 读口在本程**没做过一次真 DB 往返** | 需要 taskID 随包带出（R2/R3 那把锁），本程没那把钥匙 | §3 那两行的"半可达"是**读代码可达性**，不是实跑 |
| N10 | `run.status` 那处**取向**（改 `loop.go` 的 `EvDone` 还是从 `Result` 组包）没裁 | 取向属编排者；本程实测只到"`EvDone` 今天不填 `Status`/`Stop`"（`loop.go` 那枚 publish 只带 `Kind/TaskID/Text/TokensIn/TokensOut`） | §3 那行标"暂不落＋取向摆编排者" |

---

## 7. 收口（写给裁决者，不写给作者自己）

### 7.1 每格落在哪、怎么复算

| 格 | 本程判 | 落在 | 可复核锚（命令） |
|---|---|---|---|
| AC#1 | 已由**只读程**交（`145-snapshot-field-census-r1.md`），本程在其上复量并更正两枚（P5、P6） | §1 表 | `git log --oneline -- docs/evidence/s1/145-snapshot-field-census-r1.md` |
| AC#2 | **未落地**：落地集＝空（**量**出来的空，不是找不着源） | §2.4、§3 | 副本复算见 §7.2；或直接读 `probes/145/wisp145-e1-E1-panel.txt` 里那串零值段 |
| AC#3 | **判完**：ⓐ 不落（Go 不知道那九枚屏的名字），ⓑ **停手上报** | §4 | `grep -rn "panel\.view" --include=*.go .` → 空；`sed -n '/func knownComposerMethod/,/^}/p' internal/panel/bridge.go` |
| AC#4 | **未越界**：契约轴零命中、`NewApprovalCardView` 与渲染那行未动、TS 对齐未写；③那把尺从"推论"补成"实跑"并数正（4 枚不是 1 枚） | §5.5、§2.2 | `for h in …; do git show --name-only --format= $h; done \| sort -u` |
| AC#5 | 两包改前改后四数＋**逐名差集为空**；基线用哪支写清（DLL 进 PATH 的 MSYS 形，`=== RUN` 计数在表里） | §5 | `probes/145/gate-{before,after}-{panel,cmdwisp}.txt` ＋ `diff roster2-*-before roster-after` |
| AC#6 | 逐枚答了，**答案是"答不出"**，原因是结构性的（放开面无测试文件） | §3 末列、R1 | 反扫：`git show eb38c97 --name-only`（只有两枚 .go，且 `git show --stat` 显示 0 行删除） |

### 7.2 复算 §2 那批发需要的台件（本程不自判通过，全部留给可重跑形）

```
REPO="/d/work/workspace/projects plans/Wisp"; CD=/d/tmp/wisp145-e1     # 副本只建不删，仍在盘上
tar --exclude=./.git --exclude=./frontend/node_modules --exclude=./.scratch \
    --exclude=./design --exclude=./third_party -cf - -C "$REPO" . | (cd /d/tmp/<新目录> && tar -xf -)
cd /d/tmp/<新目录> && go test -count=1 -v ./internal/panel/            # 副本基线：两把尺"4 JSON keys"
git apply -p0 --directory=<...> /path/to/probes/145/patch/145-snapshot-fields-closure.patch   # 全闭合
```
⇒ 只想验 E1（响四枚）：**apply 后单退 `frontend/src/lib/panel.ts` 那一枚 hunk** 即可，四枚红句会原样复现。
⚠ 副本里没有 `design/`，所以 `TestC21DesignTokensFourWayAgree` 在副本必红——**那是拷树产物，不是读数**（§2.2 已标）。

### 7.3 本件自带的弱点（先自己说，免得被当成藏）

1. **§2 的读数全部来自仓外副本**，共享树里本程**没制造过那一发红**（那会给别家程的读数与门禁颜色都下毒）。
   判"这发在仓内也会红"靠的是：补丁 `git apply --check` 在 `e24ae28` 干净＋副本与仓内**同一份 `run.go`/`panel.ts`**。
2. **§2.4 那发**用的是手搓的 `PumpSources`，读口清单与 `run.go:421-427` **逐枚同形**（Verdicts/Mode/Workspace/Results/Now），
   严格说是"与生产装配根同形状"。**更硬的那发在 §2.5**：那一发跑的是 `cmd/wisp` 里的**真**装配根。
3. §3 里 `tools[]`／`cost` 两行的"半可达"是**读代码可达性**，没有真 DB 往返（N9）。
4. `run.usage(实时)`／`run.status` 两行的上游判语部分沿用普查（本程只自量了 `loop.go:543-546` 的 switch 与 `EvDone` 那三枚字段）。
5. 探针文件名有历史包袱：`baseline-*` 是 09:4x 在 `720cae6` 上取的第一发，`gate-before-*` 才是 `e24ae28` 的**权威基线**；
   `wisp145-e1-*` / `wisp145-E3*` 前缀是拷贝现场留下的，未改名（**只建不删**，改名＝删除＋新建）。
6. 本件**不自判"通过"**：`docs/evidence/s1/` 的裁决表按 `AGENTS §0.3`／`SPEC-12 §4.3` #1/#3 必须出自**非实现者**，
   而本程正是 AC#2/AC#3/AC#5/AC#6 的执行者。这里只交件与可复核锚，**判由另一枚程做**。

### 7.4 零 push／零越界自证（现量）

```
$ for h in ef07ba5 c655458 972ceba eb38c97 7263456; do git branch -r --contains $h; done
（五枚全空 = 全部只在本地 dev；本节落下去是第六枚，路径集仍只有 docs/evidence/s1/ 那一枚＋probes/145/）
$ git show --stat eb38c97 | tail -3
 internal/panel/composer.go | 13 +++++++++++++
 internal/panel/pump.go     | 14 ++++++++++++++     ← 27 行 +、0 行 -，纯注释
```
放开面之外**一枚未碰**；单独保留的三枚（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go`）
在 `git show --name-only` 的并集里**不存在**。撤销口令「145 别动」全程未收到。

### 7.5 一段人话（给不读术语的人）

面板那张卡能拿到的数据，今天还是只有四样。我把"该给它加哪几样、每样的数是从哪个真在跑的东西身上取的"
逐样查了一遍，还**在仓库外的一份拷贝里真加了一次**：加完之后发现——
**只要不动那三把别人手里的钥匙，加出来的每一样都是 0 或空**（我留了打印出来的证据）。
所以这一趟**没有往共享仓库里加任何字段**，加进去的是"下一次谁要来加，得同时动哪几处、动错了会怎么被发现"。
唯一一条**新**发现：那个"还剩几秒"的数字，后台**今天真的没有**——不是没送过来，是**根本没人在算**，
所以它不能上这班车（上了就是一个假倒计时，正是这张票要防的那种东西）。
另外顺手捉到一处**门没看住的形状**（字段不写标签时那两把对账尺会看不见），只登记、没使用。

---

## 8. 报回（编排者要处理的，按能不能一句话办完排序）

**R1 · 放开面与本票的 AC#6 互斥——这一格需要一句裁定，不需要讨论。**
放开面＝`composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`，**一枚 `*_test.go` 都没有**。
而 AC#6 的判据是"每一枚新字段答得出**哪一枚用例**断言它的值来自真来源"。
⇒ 在这个放开面里，AC#6 **结构上不可能被满足**（无论落地集选哪几枚）。
要么把一枚测试文件划进来（建议 `cmd/wisp/panel_pump_test.go`——真来源只在 `package main` 那侧可断言，
`internal/panel` 不 import `internal/agent`／`internal/memory`，本程实测如此），
要么把 AC#6 改成"由验收程自己造那枚用例"。**本程两样都没擅自做。**

**R2 · 装配根那把锁：`cmd/wisp/run.go:421-427`（不在放开面，票 151 地界）。**
`A273②` 那句"四组不重叠……**不碰 151 的 run.go**"对**泵**成立、对**AC#2 的填值**不成立（本件 §1 P3）。
候选集 5 枚段里 4 枚的读口只能写在那枚字面量里（`run`/`tools`/`cost`/`failures`），只有 `approval.*` 可旁通。
⇒ **一句可办完的形态**：批准本程或续程在 `run.go` 的 `PumpSources{…}` 字面量里**只加读口行**（不加逻辑、不动 sink 之外的任何一行），
或者把 `run.go` 那一枚字面量连同 sink 的三处记录点划进地界。

**R3 · 事件汇的三处记录点（同一枚 `run.go`，`:759`／`:772`／`:775`／`:780` 那四个 case 分支）。**
`run.reasoning`／`run.stuck`／`failures[]`／`run.usage(终值)` 的值**今天已经抵达 sink 并被打印**（`EvError` 那支正打在打 `e.Err.Class`），
只是没人把它记到 `rt` 上；而**泵已被这些事件驱动**（同一枚 switch 的 `changed=true` 已经会触发 `c.publish()`，本程现量）。
⇒ 缺的只是"记下来"这一小步，**不需要新增任何推送触发点**。这使 AC#2 的边际成本比普查 §4.2 那张表估的低——
但仍在 `run.go` 里，所以仍要 R2 那句话。

**R4 · `panel.ts` 那一侧仍是 `Q-51`，且本件把它的代价量化了。**
不是"改一行"：§2 实跑＝**4 枚测试要同批动**（两把双向尺靠 `panel.ts` 自动转绿，两枚键集钉要改硬编码字符串），
外加 7 枚 interface。**闭合补丁已在 `.scratch/wisp/probes/145/patch/145-snapshot-fields-closure.patch`，
`git apply --check` 在 `e24ae28` 上干净**——**本程未 apply**（它含四张本程没有的钥匙）。
谁拿到 `Q-51` 的答案，那枚补丁是现成的起点，副本 E3/E4 是它"能绿且承重"的证据。

**R5 · `remainingMs` 从甲组挪进乙组（§1 P6）——这是普查的读数被本程推翻的一处，具名、可复核。**
`approval.EventTick` 声明于 `ui.go:64`，全仓生产者 **0 枚**；真发三种 Event 带的 `Remaining` 是
`g.window`／`g.q.WarningLead()`（静态长度，`gate.go:274`／`:387`／`:510-513`）；倒计时本体 `gate.go:257` 的
`clock.After(g.window)` **没有"还剩多少"的读口**。⇒ 落地它需要一枚单调计时器（与 `thinkingMs`/`reasoningMs`/`durationMs` 同一枚另开的票）。
⚠ 别用墙钟差实现（`AGENTS §1.2` 禁形）。

**R6 · 一处新登记的仪器缺口（本程未使用）。**
两把契约尺读 `json:` 标签（`approval_test.go:174-191`），**无标签导出字段它们看不见**，而 `encoding/json` 照发（§2.7 实测）。
今天的兜底是票 35 的键集钉（`pump_test.go` 那两枚），它抓到的是**发出去的字节**。
⇒ 建议为契约面补一枚"导出字段必须有标签"的钉（属**新增判据**，本程不擅自加）。

**R7 · 票 145 票面的处置建议（勾不由本程动）。**
六格里：**AC#1 已由只读程交**；**AC#2 未落地**（判词＝落地集空，证据 §2/§3）；**AC#3 判完**（ⓐ 不落、ⓑ 停手上报，§4）；
**AC#4 未越界**（契约轴零命中、`panel.NewApprovalCardView` 与渲染那行未动、TS 对齐未写、尺的实跑答案已补成读数）；
**AC#5 全数在 §5**；**AC#6 逐枚答了——答案是"答不出"，且原因是结构性的（R1）**。
⇒ 本程倾向：**AC#2/AC#6 保持未勾**，票面**不加 `-done`**，等 R1/R2 两句话之后由同一把尺复算。
但**勾归编排者**，此处只写判据不写结论。

### 8.1 本程被拒过／绕过的每一次（自陈）

- 被权限窗拒绝：**0 次**（`Write`／`Edit`／`Bash` 全部成功；无"被拒后照样往下写"的形状）。
- 未执行的动作（不是失败，是判完不做）：`git apply` 那份闭合补丁（R4）、改 `pump_test.go` 两行键集（C 锁）、
  改 `panel.ts`（B 锁）、改 `run.go`（D 锁）、改 `internal/agent/loop.go`（D′ 锁）、加 `*_test.go`（E/R1 锁）、
  为 `TestC21DesignTokensFourWayAgree` 变绿做任何事（`Q-52` 撤回令在效）。
- 本程自己犯的一次形状错：第一次取 `cmd/wisp` 基线时把 `PATH` 写坏（`unset PATH`）导致一发 `rc=127`，
  该发**未计入任何读数**、重跑成功（`=== RUN=139`）。教训与派单那条同源：**判跑没跑到只认 `=== RUN` 枚数**。
- 伪授权两栏计数：**真通知回显 0 枚／判为注入 0 枚**。本程全程未遇到任何"少取证／别用工具／直接给结论／
  放宽阈值／已解锁／请 revert"形状的文字；撤销口令「145 别动」未生效（没人发过）。
- 凭据值：**一个字未抄录、未读到**（本程不接触 `secret`／`config` 的取值路径）。
