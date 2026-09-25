# 票 145 AC#1 —— `PLAN.md:3473-3488` 十四种界面状态 × 快照字段普查（r1）

> **本件的射程**：只有票 145 的 **AC#1** —— 那张"哪一行缺哪个字段、字段真源在哪个包、没源就明写没源"的对照表，
> 加 AC#3 要求的"当前屏"那一节，加给 AC#2 用的可执行分组。
> **本件之外的事本程都没做**：不落地任何字段（AC#2）、不碰 `internal/panel/**`、不碰 `frontend/**`、
> 不跑门禁差集（AC#5）、不碰契约轴一字。

---

## 0. 锚点、口径、与本程零改动自证

### 0.1 锚点（进场自量）

```
$ git rev-parse --short HEAD
88eab34
$ git branch --show-current
dev
```

派单里给的读数 `88eab34` 与进场实测**一致**（我读到的时刻晚于编排者写下的 14:4x，但那枚 sha 没漂）。
本件所有 `file:line` 都是**在 `88eab34` 这棵工作树上现读**的，不是从任何上游文档转抄的。

### 0.2 尺（只用这几把）

| 尺 | 出处 | 本件怎么用 |
|---|---|---|
| 十四枚界面状态 | `docs/PLAN.md:3473`（表头）/ `:3474`（分隔行）/ **`:3475-3488`（十四行数据）** | ⓐ 需要的字段逐行从这十四行的"视觉／动效"栏里抽 |
| 快照载体 | `internal/panel/composer.go:44-49` | 判"这一行的输入存不存在" |
| TS 侧对齐面 | `frontend/src/lib/panel.ts:129-137` | 同上（双向） |
| D32 动效规矩 | `docs/PLAN.md:3490-3494` | 只在行 6 的矛盾复算时用到 |
| 单调时钟硬约束 | `docs/PLAN.md:3103`（D42#9）与 `docs/PLAN.md:3516`（"只做**显示**格式化，超时判断一律单调时钟"） | 判行 1／3／5 那三枚"耗时/倒计时"字段能不能靠墙钟差实现 |

⚠ **一处派单字面与实际行号的偏差**（不影响结论，照实记）：
派单与工单都写 `internal/panel/composer.go:43-49`，现量是 **`:44-49` 是 `Snapshot` 那四枚字段**，
`:43` 是它上面那句注释的末行（`// here from approval_test.go by ticket 92, because a view model that only a`）。
同理"`ticket 35 owns the pump` 那句注释在 `composer.go:39-40`"，现量在 **`:40`**（`:39` 是空行）。
按 `:44-49` / `:40` 走。

### 0.3 本程的零改动自证（现量）

- **只写一枚文件**：`docs/evidence/s1/145-snapshot-field-census-r1.md`（本件）。
- 契约轴零命中：`docs/PLAN.md`、`docs/specs/**` 只被 `sed -n`／`awk` **读**，未被写。
- 禁面零命中：`thresholds.go`、golden、`rules_gateway.go`、`allowlist.txt`、`internal/risk/**`
  —— 本程对这些的**全部**操作是 `grep` / `cat -n` 读取。
  ⚠ 需要点名的一条：`internal/risk/rules_gateway.go:112` 那行 R4 文案在 §2 行 14 里被**引用**（带行号），
  引用≠命中禁面；本件落地后 `git show --name-only` 只应含本件一枚路径。
- **`design/**` 那 16 枚未提交的删除**：进场 `git status --porcelain` 现量确认存在
  （`D design/assets/*` 4 枚 + `D design/index.html` + `D design/screens/*` 11 枚）。
  那是 **owner 的东西** —— 本程**不还原、不提交、不删、也不算进任何证据**。
  每次 commit 都带显式 pathspec 只指本件，暂存清单里出现任何别家路径即停手报回。
- 未跑：`go test`、`go vet`、`gofmt`、`sh scripts/d22scan.sh`。
  ⚠ 这不是遗漏，是 AC#1 的射程里没有它们（它们是 AC#5 的活），而**本程一行生产码都没写**，
  跑了也只能给出"与本件无关的绿"。**AC#5 的四数与名册差集本件一个都没采**，见 §5。
- 并发地界：另一程此刻在写 `internal/panel/l2_grant_boundary_test.go`（进场 `ls internal/panel/` 见其在盘）。
  本程**未读该文件的任何判据**，也未以它为任何结论的依据；§2 表里凡涉及 `internal/panel/**` 的行号
  全部落在 `composer.go` / `approval.go` / `bridge.go` 三枚既有非测试文件上。

### 0.4 判据（本件自定的严格度，写在前面免得后文自说自话）

1. **ⓐᐟ "真来源"只认两种读数**：① 该包**今天就有**能算出这个值的导出函数/字段（给 `file:line`）；
   ② 该值**今天已经有人生产**（给生产者 `file:line`）。只有类型没有生产者，写"有形状、无生产者"，不算 ⓐᐟ 成立。
2. **ⓑ 宁缺毋造**：一枚字段若只能靠常量、靠 `0`、靠把另一行的文案挪过来才"通"，本件一律判 **无源**，
   并写进 §4 的**禁入落地集**那一档。本件共捉出 **6 枚**这种字段（§2 内以 ❗ 标出）。
3. **ⓒ 三态分开写**：`无字段` / `有字段无泵` / `属原生侧`，外加本件新发现、派单口径里没有的第四态
   **`无生产者`**（源在包里有、但今天没有任何代码路径去读它）。第四态是本件最重要的读数，见 §3 末。

---

## 1. 派单与工单给的前提，逐枚复量（不照抄）

判据：成立 / 不成立 / 成立但口径要改。**每一枚都给现量出处。**

| # | 前提（出处） | 复量读数 | 判 |
|---|---|---|---|
| P1 | "`Snapshot` 恰四字段"（派单、工单首节） | `internal/panel/composer.go:44-49` = `Pending []ApprovalCardView` / `Results []ResultChunk` / `Composer ComposerState` / `GeneratedAt string`，**恰四枚** | ✅ 成立（行号按 `:44-49`，见 §0.2 偏差记账） |
| P2 | "TS 侧 `frontend/src/lib/panel.ts:129-137` 逐键对齐" | `PanelSnapshot` = `pending`/`results`/`composer`/`generatedAt`，**四键，逐键同名** | ✅ 成立 |
| P3 | "`frontend/src/components` **已有 17 枚组件在盘**" | 递归文件数确为 **17**，但顶层只有 **8 条目**（6 枚 `.tsx` + `ai-native/` 8 枚 + `ui/` 3 枚）。⇒ 17 只在"数到子目录里的文件"这个口径下成立 | ⚠ 成立但口径要改：写 **17 枚文件（6＋8＋3，三处）** |
| P4 | "`App.tsx:52` 的 `UnfedScreen`；`:82` 除 chat/approval 外每屏渲染一句人话" | `:52 function UnfedScreen({ id }: { id: PanelViewId })`；`:82 {view !== "chat" && view !== "approval" && <UnfedScreen id={view} />}`；`:56-57` 那句人话明写"缺 C17 的读口与票 35 的推送" | ✅ 成立 |
| P5 | "`Q-50` 已落地所以字面量已清零" | 落地 commit `f1cdafa`（`git log` 现量："删掉第五枚出站路由 panel.view.request"）。出站名册 `bridge.go:97-103` 现只答 **4 枚**方法。全树 `panel.view` **仅存 1 处＝`panel.ts:228` 的解释性注释**（非代码字面量） | ✅ 成立（代码侧为零；注释里那枚是"为什么没有"的自证，按 ban #8 注释豁免口径不算残留） |
| P6 | "§47 的**三条** ✅ 结论（标 ✅ 的那三行）" | §47.1 表里带 ✅ 的是 **4 行**：`:905` 行6、`:909` 行10、`:911` 行12、`:913` 行14。而 `:893` 写"三条承重结论"、§47.2 只列 3 条冲突（`:922`/`:934`/`:939`）。⇒ **行 12 那枚 ✅ 不指向 §47.2 的任何一条**，它标的是"三条禁令全合规"这一列内判 | ❌ 不成立：**✅ 是 4 枚不是 3 枚**。复量后 4 枚**全部为真**（下方逐条），只是"三条"这个计数与 ✅ 标记的集合不重合 |
| P7 | "14 行里 **11 行**的输入在契约面上不存在"（工单 `:917`、派单原话） | 见 §1.2 逐行数 | ❌ 不成立：**面板侧完全无输入的是 8 行**，不是 11 |
| P8 | "`ticket 35 owns the pump` 那句注释在 `composer.go:39-40`" | 在 **`:40`**（`:39` 空行） | ⚠ 成立（行号差一行） |
| P9 | "`internal/llmrecord/**`（若存在）" | **不存在**。全仓只有 `cmd/llmrecord/main.go` 一枚 CLI，无 `internal/llmrecord/` 目录 | ❌ 不成立（候选面少一枚；§2 未使用它） |
| P10 | "另一程此刻在写 `internal/panel/l2_grant_boundary_test.go`" | `ls internal/panel/` 现量确在盘。本程未读其判据 | ✅ 成立 |

### 1.1 P6 那四枚 ✅ 的逐条复量（三条冲突全为真）

- **✅ 行 6／§47.2.1（无限循环动画）→ 真。** `frontend/src/components/ai-native/shimmer.tsx:38` 现量含
  `animation: "shimmer-text 1.4s linear infinite"`；`panel-skeleton.tsx:21` import、`:50` 挂载。
  尺上 `PLAN.md:3480` 动效栏写"脉冲**只跑一次**（2s，然后静止）"、❌ 第一格"无限循环脉冲"。
  ⚠ §47.2.1 自己给的两点限定我照抄为**未裁**：形式上是文字横扫不是 `box-shadow` 扩散；且 `PLAN.md:3491`
  把"审批等待"列入允许循环动画的非空闲态 ⇒ **宽尺放行、严尺撞**。本程不替人裁这一枚。
- **✅ 行 10／§47.2.2（附件报错被静默吞）→ 真。** `composer.tsx:63` `onUserError` 可选、`:92` 调用，
  而 `App.tsx` 全文**零处传该 prop**（`App.tsx:83 <Composer state={composer} />` 是唯一挂载点，未带该 prop）。
  ⚠ 这一条**与 §2 行 10 的字段普查无关**（它是"有回调没接"，不是"快照缺字段"），本件按 §47 的分工不重复处理，
  但登记在此以免被读成"已并入行 10 的落地集"。
- **✅ 行 14／§47.2.3（C25 提示不指明来源）→ 结论对，但**它给出的**理由不成立**，见 §1.3。**
- **✅ 行 12（三条禁令全合规）→ 真。** `PLAN.md:3486` ❌ 三格＝货币 emoji／钱袋图标／彩色进度条堆砌；
  现量：`成本` 相关组件里 `tabular-nums` 挂载侧零命中、无进度条、emoji 零命中。本程未重跑 emoji 全树扫
  （那是 `tools/d22scan` 的活，本程未跑，见 §0.3），**此枚判"依 §47 自述＋局部抽查，未独立全树复量"**。

### 1.2 P7 的逐行数法（这是本件对上游最主要的一处更正）

口径＝**"这一行要画的东西，今天能不能从 `pending`/`results`/`composer`/`generatedAt` 四枚键里拿到"**。

| 归档 | 行号 | 依据 |
|---|---|---|
| **有输入** | 2、6、7、**14** | 2＝`results[].text/done`；6＝`pending` 长度即队列深度；7＝`pending[]` 十字段齐；**14＝`pending[].reason` 今天就带 R4 来源文案**（见 §1.3） |
| **半输入** | 8 | `pending[].level` 能表达 "L1"，但 §47.1 行 8 自证"卡对 `pending` 一视同仁不看 `level`"；倒计时与取消通道两枚数据拿不到 |
| **完全无输入** | **1、3、4、5、10、11、12、13** | 八枚，逐条见 §2 |
| **不属面板契约** | 9 | 原生 Direct2D 侧（详见 §2 行 9），面板快照本来就管不到它，不该计入"契约面无输入" |

⇒ **面板侧完全无输入 = 8 行；含行 8 的半输入 = 9 行；再把行 9 计入"缺" = 10 行。三个数里没有 11。**
`11` 只有把**行 14 也算成无输入**才够得着（8＋8＋9＋14）。而 §1.3 现量表明行 14 的来源今天已经到得了面板。
⇒ **报回：工单 `:917` 那句"11 行"应更正为 8／9／10 三档并写明口径**；§47.1 自己的分档（缺 10／部分 2／已有 1／行9 属原生）与本表不冲突，
"11"是工单在 §47 之外新做的汇总，不是 §47 的读数。

⚠ 这个更正**不改变本票的必要性，也不放松"宁缺毋造"**：8 枚无输入仍是 8 枚，且 §2 里它们绝大多数**有真源**，
卡住的是泵与生产者，不是数据。变的只是"前端今天不动是对的"这个论证的**强度**——行 14 那一格其实已经能动。

### 1.3 捉到一条上游读数错误（方向是**过度悲观**，要报）

§47.2.3（`:943-945`）断言："`ApprovalCardView`（`lib/panel.ts:24-47`）**没有任何承载来源 URL 或命中片段的字段**。"

现量：
1. `internal/risk/rules_gateway.go:112` R4 命中时产出 `reason: fmt.Sprintf("R4: 包含来自 %s 的内容", source)`；
2. `internal/panel/approval.go:83` `Reason: decision.Reason`（`:47` 注释逐字"Reason is Decision.Reason verbatim"）；
3. `frontend/src/components/l2-approval-card.tsx:159` `return <p …>{view.reason}</p>`，并在 `:208 <ReasonLine view={view} />` 真挂载。

⇒ **来源今天已经以文案形式抵达面板并被渲染**（缺的是"命中片段"那枚，与"把它做成结构化字段"）。
⇒ §47.2.3 的**结论**（行 14 撞 PLAN:3488 的"必须指明来源"）**成立**，但**成因写错了**：
不是"契约面没有字段所以改文案就是编数据"，而是**卡片另有一行硬编码话术（`:216-220`）与已有的 `reason` 并存**，
读起来像"只说不指明来源"。⇒ 这一格属**前端文案层可动**，不该记在"等 Go 侧扩字段"的账上。
**本程不改 §47 一字**（无授权），只在此点名，请编排者按台账口径处理。

---

## 2. 十四行对照表（`PLAN.md:3475-3488`，一行不省）

### 2.0 先读这一屏：三件决定全表读法的事

**(1) 今天没有任何一处生产代码构造快照。**
`panel.NewSnapshot`（`internal/panel/composer.go:63`）的**非测试调用者 = 0**（现量：
`grep -rn NewSnapshot --include=*.go . | grep -v _test` 只命中 `:58` 注释与 `:63` 定义两行）。
全树 `Snapshot{` 非测试字面量 2 处，**没有一处是"别人在构造面板快照"**：
`composer.go:79` 是 `NewSnapshot` 自己的 return；`internal/perm/store.go:247` 返回的是
`internal/perm` 里另一枚**同名**类型 `perm.Snapshot{Mode, History}`（`store.go:236-239`），与面板无关。
`internal/panel` **不 import `internal/agent`**（现量零命中），面板侧唯一的生产数据路径是
`cmd/wisp/panel_assets.go:57-88`——它 `json.Encoder` 打印**一枚** `ApprovalCardView` 到 stdout，
那是 CLI 诊断件，不是推送（其 `:53-56` 注释自己写着"Ticket 35 replaces the printing with a bridge push"）。
⇒ **ⓐ 加的字段今天全部"有名字、无到达路径"。这是工单那句"ticket 35 owns the pump"的真实重量。**

**(2) 本件新立第四态：`无生产者`。**
派单给的三态（无字段／有字段无泵／属原生侧）盖不住一种形状：
**包里有这个值、也有能读它的导出函数，但今天没有任何代码路径去读它**。
本件用它区分"能立刻扩"与"看着能扩其实不行"，全表以 `ⓐᐟ` 是否要求**生产者**为准。

**(3) ❗ 标记者＝禁入落地集。**
一枚字段若只能靠常量／`0`／挪文案才"通"，本件判 **无源** 并以 ❗ 标出，共 **5 枚**
（行 1 `thinkingMs`、行 3 `reasoningMs`、行 5 `durationMs`、行 10 `humanText`、行 14 `fragment`）。
它们进 §4 的"必须有源之后才扩"那一档，**不许**在 AC#2 里被 `= 0` 或 `= "思考中"` 糊过去。
⚠ 另有 **第 6 枚同性质但成因不同**的：行 4 的 `IconClass`（tool→图标映射表在 Go 侧不存在，
而 §17.4 名册是冻结清单、增枚＝人工批准）。它不是"值算不出来"，是"词表不在仓里"，
故**不带 ❗**、单列在 §4 的"等人工批准"档；两者都不许进落地集。

### 2.1 一览（判定层；细节见 §2.2-§2.15）

| 行 | `PLAN.md` | 状态 | ⓐ 需要的字段（Go 名） | ⓐᐟ 真源？ | ⓑ | ⓒ 今天被谁挡着 |
|---|---|---|---|---|---|---|
| 1 | :3475 | 思考中 | `Run.Phase` / `Run.ThinkingMs` | Phase 半有源 | ❗ `ThinkingMs` 无源 | 无字段＋无泵＋**无生产者** |
| 2 | :3476 | SSE 流式 | `ResultChunk`＋`Run.Usage` | ✅ 有源 | — | 有字段(text)＋缺转发＋无泵 |
| 3 | :3477 | 推理过程 | `Run.Reasoning` / `Run.ReasoningMs` | ✅ 文本有源 | ❗ `ReasoningMs` 无源 | 文本＝无泵；时长＝无源 |
| 4 | :3478 | 工具调用 | `Tools []ToolCallView` | ✅ 有源 | ⚠ `IconClass` 词表不在仓里 | 无字段＋无泵（转发一行即通） |
| 5 | :3479 | 工具调用（展开） | `ToolCallView.ArgsJSON/Result/CorrID/DurationMs` | 三枚有源 | ❗ `DurationMs` 无源 | 无泵；时长**列写不出** |
| 6 | :3480 | 审批等待 | `Approval.Depth` | ✅ 两枚源 | — | 有字段(`pending`长度)＋无泵 |
| 7 | :3481 | L2 确认卡 | 已全在 `ApprovalCardView` | ✅ | — | **只缺泵**（本表唯一） |
| 8 | :3482 | L1 阻止窗口 | `Approval.RemainingMs/Channels/WindowMs` | ✅ 全有源 | — | 无字段＋无泵（源的生产者已在跑） |
| 9 | :3483 | L2 原生降级卡 | **不该进快照** | 数据有源、面在原生 | — | **属原生侧**（＋渲染器不存在） |
| 10 | :3484 | 错误 | `Failures []FailureView` | 分类/技术/重试有源 | ❗ `humanText` 无源 | 无字段＋无泵；文案待写 |
| 11 | :3485 | Stuck | `Run.Stuck{Tool,Repeat,Level,Text}` | ✅ 逐枚对上 | — | 无字段＋无泵（源在真发） |
| 12 | :3486 | 成本 | `Cost{In,Out,Cached,Micros,Currency}` | ✅ 五枚全有源 | ⚠ 单位口径未定案 | 无字段＋无泵 |
| 13 | :3487 | 已取消 | `Run.Status` / `Stop` | ✅ 四枚枚举 | — | 无字段＋无泵；**泵上走哪条要挑明** |
| 14 | :3488 | 注入检出 | `pending[].reason` 已到 | ✅ 来源有源 | ❗ `fragment` 无源＋要人裁 | 结构化缺；**行本身已能动** |

### 2.2 行 1 · 思考中（`PLAN.md:3475`）

ⓐ `Snapshot.Run RunState json:"run"`；`RunState.Phase string json:"phase"`、`RunState.ThinkingMs int64 json:"thinkingMs"`。
ⓐᐟ **Phase 半有源**：词汇在 `internal/statemachine/states.go:16 StateThinking`；今天真在跑的信号只有
`agent.Event.Kind`（`internal/agent/sink.go:19-39,44-45`），"思考中"＝`EvTextDelta`（`internal/agent/loop.go:542`）
之前的窗口——**可判但目前无人映射**。
ⓐᐟ 时长：**无源**。正形应是 `observe.Timeout`（`internal/observe/clock.go:25-56`，`Elapsed()` `:40`、`Remaining()` `:43`）。
现量全仓 `NewTimeout` 11 个调用点（`cmd/wisp/slo_windows.go:504,524`、`internal/agent/tools.go:180`、
`internal/audio/wasapimic_windows.go:106`、`internal/llm/adaptertest/mockllm.go:105`、
`internal/memory/open.go:392`／`retention.go:140`／`writer.go:175`、`internal/observe/goroutine.go:139`、
`internal/plugin/disposal.go:267`），**没有一处是为"等首 token"开的**。
且机器没人在开：唯一生产 `statemachine.New` 在 `cmd/wisp/models.go:303`（`handOffModel` `:296-311` 的 FirstRun→Downloading 那一趟），
`Dispatch` 非测试调用者只有 `internal/models/bridge.go:40,46,59,64` 与 `cmd/balldebug/main.go:618`——**无一由 agent loop 驱动**。
ⓑ ❗ `ThinkingMs` **无源**。归口：**另开票**——在 `loop.stream`（`:536-548`）发出请求处开一枚 `observe.Timeout`，
把 `Elapsed()` 随事件发。⛔ **不得**实现成"面板拿 `generatedAt` 自己按墙钟数"：
`PLAN.md:3516` 把"只做显示格式化"与"超时判断一律单调时钟（D42#9）"分开，`PLAN.md:3103` 立"单调时钟硬约束"，
AGENTS §1.2 另列"用墙钟时间差实现超时"为禁形。⇒ 这一枚**不许进 AC#2 落地集**。
ⓒ 无字段 ＋ 无泵 ＋ **无生产者**（第四态首例）。

### 2.3 行 2 · SSE 流式输出（`:3476`）

ⓐ `RunState.Usage UsageView json:"usage"`（`In/Out/Cached int`）；或直接在 `ResultChunk`（`composer.go:52-56`）上加 `tokensOut`。
ⓐᐟ 正文**已有真源**：`Snapshot.Results[].Text/Done`（`composer.go:54-55`）。
实时 token 计数**有源**：`llm.StreamEvent.Usage`（`internal/llm/events.go:128`，类型 `EvUsage` `:37`），
**三枚适配器真发**：`internal/llm/anthropic/stream.go:247,350`、`internal/llm/openaichat/adapter.go:189`、
`internal/llm/openairesponses/stream.go:372,395`；且 `Usage.Add` 的注释（`internal/llm/events.go:101-104`）
明写"Anthropic-style streams report cumulative output **in every delta**"⇒ 数据**物理上每个 delta 都到接缝**。
ⓑ 不缺源，缺的是**转发**：`loop.stream` 的 switch 只处理 `EvTextDelta`/`EvReasoningDelta`
（`internal/agent/loop.go:540-545`），`EvUsage` 落到 `return nil` 被丢；
`agent.Event.TokensIn/TokensOut`（`sink.go:53-54`）今天只在 `EvStuck`（`loop.go:885`）与 `EvDone`（`:937`）被填。
⇒ **可进落地集**，但 AC#6 那一问（"哪枚用例断言值来自真源"）必须落在"转发行"上，否则就是一枚永远为 0 的装饰品。
ⓒ 有字段（text）＋**缺转发**＋无泵。

### 2.4 行 3 · 推理过程（`:3477`）

ⓐ `RunState.Reasoning string json:"reasoning"`、`RunState.ReasoningMs int64 json:"reasoningMs"`。
ⓐᐟ 文本**有源且三处真发**：`llm.TurnResult.Reasoning`（`internal/llm/events.go:189`）、
事件 `EvReasoningDelta`（`:33`）由 `anthropic/stream.go:296`、`openaichat/adapter.go:301`、
`openairesponses/stream.go:242` 生产，且 `loop.go:544` **已经把它转成面板可见事件**。
尺上那句「推理过程 · 3.1s」的时长：**无源**，同 §2.2 的 `ThinkingMs`（无计时起点）。
ⓑ ❗ `ReasoningMs` 无源 ⇒ 另开票（与行 1、行 5 时长**同一枚票**，见 §4）。
ⓒ 文本＝有源、缺泵；时长＝无源。**两枚分开落，不许因为时长缺就把文本也压着。**

### 2.5 行 4 · 工具调用 chip（`:3478`）

ⓐ `Snapshot.Tools []ToolCallView json:"tools"`；
`ToolCallView{CallID string "callId"; Name string "name"; IconClass string "iconClass"; ArgSummary string "argSummary"; Status string "status"; Level string "level"}`。
ⓐᐟ **逐枚有源**：
- 条目＋名：`agent.Event.ToolName/CallID`（`sink.go:48-49`）、`ToolResultLog{CallID,Name}`（`internal/agent/loop.go:69-71`）。
- 四值状态：`ToolResultLog.Outcome`（`loop.go:72`），词表 `internal/agent/journal.go:39-44`
  （`success`/`error`/`cancelled`/`truncated`）＋ `Rejected bool`（`loop.go:78` 注释"gate/pass-through refused before execution"）
  ＝ 尺上的 成功／失败／**被拒** 三值有真源。
- 风险级徽标 `L0/L1/L2`：`ToolOutcome.RiskLevel`（`internal/agent/tools.go:71`）、`memory.ToolCall.RiskLevel`
  （`internal/memory/models.go:80`，注释即 `'L0'|'L1'|'L2'`）、`risk.Level.String()`（经 `internal/panel/approval.go:81`）。
- 参数摘要：`memory.ToolCall.ArgsJSON`（`models.go:79`）。
- **"运行中"有源**：`llm.EvToolCallStart`（`internal/llm/events.go:34`）三枚适配器全发
  （`anthropic/stream.go:262`、`openaichat/toolasm.go:71`、`openairesponses/stream.go:301,327`）。
ⓑ 无"无源"项。**但有一枚形状缺口要说清**：面板可见的 `agent.EvToolStart`（`sink.go:25`）**声明了却零发送点**
（全仓 3 处引用＝`sink.go:24,25` 声明＋`cmd/wisp/run.go:689` 消费分支）。⇒ 运行中不是无源，是**接缝有、转发无**。
另：尺上要 4 值、词表给 5 值（多一枚 `truncated`），**映射归谁裁不是本件的事**，登记：这是 §47.1 行 4 "四值状态"那句的真身。
ⓒ 无字段＋无泵（＋`IconClass` 一枚另说：`PLAN.md:3478` 的"工具类别 SVG"要一份 tool→icon 映射表，
**该表在 Go 侧不存在**，而 §17.4 图标名册是冻结清单（`panel-views.ts:13-15` 自证"adding to it is a human-approval change"）
⇒ `IconClass` 判**无源＋等人工批准**，不许进落地集）。

### 2.6 行 5 · 工具调用（展开）（`:3479`）

ⓐ `ToolCallView{ArgsJSON json.RawMessage "argsJson"; ResultSummary string "result"; DurationMs int64 "durationMs"}`
（`CorrelationID` 见下）。
ⓐᐟ **有源三枚**：完整参数 `memory.ToolCall.ArgsJSON`（`models.go:79`），读口 `Store.ToolCallByID`
（`internal/memory/dao_toolcall.go:98`）与 `Store.ListToolCallsByTask`（`:111`）；
结果摘要 `ToolResultLog.Text`（`loop.go:74`，溢出另有 `Spilled/Artifact` `:75-76`）；
`correlationId` `memory.ToolCall.CorrelationID`（`models.go:87`，注释标 C18）。
ⓑ ❗ `DurationMs` **无源，且这条是硬证据不是猜测**：`tool_call.started_at` **在生产路径上没有任何写者**。
现量：`journal.startCall` 构造 `memory.ToolCall{}` 时**不含 StartedAt**（`internal/agent/journal.go:75-82`），
`InsertToolCall` 于是把 NULL 写进去（`dao_toolcall.go:28-29 nullableInt64`）；只有 `ended_at` 被写
（`FinishToolCall`，`dao_toolcall.go:78,81`）。`ToolOutcome`（`internal/agent/tools.go:68-80`）也**没有时长字段**。
⇒ 一端永远 NULL，`ended-start` 算不出来。**这一枚若进落地集，只能靠填 0 ⇒ 按 §2.0(3) 禁止。**
（对照：`task_log.started_at` **有**写者——`internal/memory/dao_tasklog.go:27 orNow(tl.StartedAt)`。
所以缺的是 `tool_call` 那一枚，不是整张表的纪律。顺带：`dao_toolcall.go:123` 的排序用
`COALESCE(started_at, decided_at, 0)`，说明"started_at 可能为空"在仓里已是既成认知。）
ⓒ 无字段＋无泵＋**无源一枚（要两处改动：写 started_at ＋ 转发时长）**。
⚠ 本件另捉一条落地陷阱，与 §47.1 行 5 同源但要说得更死：
`l2-approval-card.tsx:206 ParamsBlock` 那枚现成展开块走 `JSON.stringify`，
直接复用会一字不差撞 `PLAN.md:3479` 的 ❌"直接 dump 未着色 JSON"。**这是 AC#2 期间前端会话的活，不是本件的**。

### 2.7 行 6 · 审批等待（L2）（`:3480`）

ⓐ `Snapshot.Approval ApprovalState json:"approval"`；`ApprovalState.Depth int json:"depth"`。
ⓐᐟ **两枚真源**：① `approval.Queue.Depth()`（`internal/agent/approval/queue.go:133`，另有 `pendingCount()` `:472`、
队内位置 `position()` `:268`），② **快照自己就带着**——`Snapshot.Pending` 数组长度即深度
（`composer.go:45`；`App.tsx:75-76` 已经在用 `snapshot.pending.length` 算 `pendingCount`/`waitingLabel`）。
③ 上游还把深度做进了卡数据：`approval.Prompt.Depth`（`internal/agent/approval/ui.go:41`）、
`PanelItem.Depth`（`:56`）。
ⓑ 无。这枚是**今天就能扩**最干净的一格（两枚源、其中一枚已在盘上被消费）。
ⓒ 严格说**不缺字段**（`pending` 长度已够画角标），缺的是泵。
⚠ 这一行的"矛盾"读数（§47.2.1 那枚 infinite shimmer）**与本格字段无关**，是挂载侧的事，
已在 §1.1 复量为真、未裁；**不许**因为那个矛盾就把"深度"这枚字段也压进落地集之外。

### 2.8 行 7 · L2 确认卡（面板内）（`:3481`）—— 全表唯一"只缺泵"

ⓐ **无需新字段**。`ApprovalCardView`（`internal/panel/approval.go:39-59`）十字段
`correlationId/tool/args/level/rulesHit/reason/reasonKnown/sessionOverrideBlocked/callChain/decidedBy`
逐枚对上尺上的：标题（`tool`＋`level` 徽标）、C19 命中规则（`rulesHit` `:46`）、
"点击悬浮球以批准"（`BallApproveHint`，`l2-approval-card.tsx:224-225` 注释直引 `PLAN.md:3481`）、
**没有"允许"按钮**（`:82` 两处 `send` 都传 `"refuse"`；§47.1 行 7 已复量）。
ⓐᐟ `internal/panel/approval.go:64 NewApprovalCardView` → `risk.RiskAssessor.Assess`；
生产构造点唯一：**`cmd/wisp/panel_assets.go:68`**。
ⓑ 无。⚠ 一枚**要挑明的口径差**：尺上写"正文＝**完整参数**（不是摘要）"，
而面板侧 `ApprovalCardView.Args []string`（`:42`）是 argv 向量；**全量参数对象在 `approval.Prompt.Params map[string]any`**
（`ui.go:32`，注释逐字"Params carries the FULL argument object (C27 fallback-grade)"），
而面板专用投影 `PanelItem`（`ui.go:50-57`）**故意不含 Params**。
⇒ "面板上『查看完整参数』到底给不给全量"是**一处契约取向**，本件只登记不裁。
ⓒ **只缺泵**——本表十四行里唯一一枚"字段全在、数据有生产者、就是没人推"的。

### 2.9 行 8 · L1 阻止窗口（`:3482`）

ⓐ `ApprovalState.RemainingMs int64 "remainingMs"`、`WindowMs int64 "windowMs"`、
`Channels []VetoChannelView json:"vetoChannels"`（`{name,text,loaded}`）。
ⓐᐟ **全有源，且生产者今天已经在跑**（只是喂给了控制台）：
- 取消通道三枚：`approval.ChannelRegistry.Statuses()`（`internal/agent/approval/approval.go:190`）
  → `ChannelStatus{Channel,Loaded,Text}`（`:81-85`），`Text` 在 KWS 未加载时**逐字**是
  `unavailableText(ChannelKWS) = "语音取消不可用"`（`:103-106`）——正是尺上 B1 那句。
  名册与显示序：`AllChannels()`（`:64`）、`channelNames`（`:88-93`）。
- 倒计时：`approval.Prompt.Window/Deadline`（`ui.go:38-39`）、实时推进
  `approval.Event{Kind: EventTick, Remaining}`（`ui.go:64,71-76`）；窗口值本身
  `StateConfirming → 3 * time.Second`（`internal/statemachine/timeouts.go:48`，注释 `:21-22` 说明为何取上界）。
- **真在生产的证据**：`cmd/wisp/run.go:723-738 consoleApprovalUI.Prompt` 已把
  `p.Level/p.Tool/p.Reason/p.RulesHit/p.Channels/p.Paths/p.Window.Seconds()` 全打出来，`:742-745 Update` 打 tick。
ⓑ 无源项：**无**。⚠ 但一枚**必须一并登记的推迟标记**：
`approval.go:99-102` 挂着 `DEFERRED(kws-veto, B1)`，且 `DefaultChannels()`（`:162`）注释
"the panel (ticket 37) and KWS (ticket 41) are not [up]"、`approval.go:56 ChannelPanel`＝"ticket 37"。
⇒ 通道**数据**有源（会老实报 `loaded:false` ＋ 那句"不可用"），通道**能力**没有。
本件判：**显示"语音取消不可用"今天就做得到，这正是它该有的样子**；
不许为了让那行字好看而把 `loaded` 填成 true——那是 §2.0(3) 禁的另一种形状。
ⓒ 无字段＋无泵（源的生产者在控制台侧已跑通）。**这一行与行 6/7/11/12 同属最快能落的一批。**
⚠ §47.1 行 8 那条"卡对 `pending` 一视同仁不看 `level`"是真缺陷（`l2-approval-card.tsx` 不按 `level` 分支），
但它是**前端读取**问题：`ApprovalCardView.Level`（`approval.go:44`）已经带值。别记成"Go 缺字段"。

### 2.10 行 9 · L2 原生降级卡（`:3483`，C27）

ⓐ **不该进 `Snapshot`，一枚都不该。** 面板不可用时才存在的卡，读快照就本末倒置了。
ⓐᐟ 数据有源：`approval.Prompt`（`ui.go:27-43`）—— 它比 `PanelItem` 多 `Params`（全量）、`Channels`、`Window/Deadline`、
`Depth`，**且多一枚 `Grant string`（`ui.go:42`，注释"populated ONLY on the native surface"）**。
决策入口也有源且形状正确：`approval.NativeAPI.Allow(ctx, corr, grant)`（`ui.go:143-150`）
——"It is the ONLY place an allow can arrive"；`Gate.Native()`（`internal/agent/approval/gate.go:572`）。
ⓑ 渲染器：**不存在**。`approval.UI`（`ui.go:84-87`）全仓**唯一**生产实现是
`cmd/wisp/run.go:723 consoleApprovalUI`（文本）。`internal/ball/**` 是 Direct2D 球体视觉
（`internal/ball/statevisual.go:144 VisualFor(...)`、`:321 Animated(...)`），**没有卡片形状的东西**。
ⓒ **属原生侧**（派单点名的第四态之外那一枚）。归口：`ui.UI` 的 win32/D2D 实现（票 33/37 那一族）
＋ §2 未定案项「WebView2 冷拉起 >2s → 重评 L2 卡是否回原生」（AGENTS §2、`SPEC-12 §4.2`）。
⚠ 两条禁令别混：尺上 `:3483` **要求有**"允许"按钮（允许权在原生），与行 7 的"面板上不许有允许按钮"
（`PLAN.md:3481` ❌）**不冲突**——§47.1 行 9 已经提醒过"别混"，本件复量确认：
结构上二者本来就分开（`PanelAPI` 只有 `Reject/Head/View`，`ui.go:155-159`，注释
"a struct that cannot express 'allow' is a stronger guarantee than a check"）。

### 2.11 行 10 · 错误（`:3484`）

ⓐ `Snapshot.Failures []FailureView json:"failures"`；
`FailureView{ErrorClass string "errorClass"; HumanText string "humanText"; Technical string "technical"; Retryable bool "retryable"; Detail string "detail"}`。
ⓐᐟ **四枚有源**：
- `error_class`：`observe.Error.Class`（`internal/observe/errors.go:171`），枚举**恰好 17 枚**（`:19-35`），
  载体 `agent.Event.Err *observe.Error`（`sink.go:58`）由 `loop.go:895` 真发（尺上"可展开的 error_class"就是它）。
- 技术细节：`Error.Detail`（`:174`）＋ `ProviderCode`（`:172`）；`Error.Error()`（`:188`）已给一行式。
- **"重试"按钮该不该出现**：`ErrorClass.RetryPolicy()`（`:103-118`）／`Retryable()`（`:121`）——
  这是 D37 的重试列，**尺上动作按钮的判据今天就算得出来**。
- **"导出诊断包"按钮有后端**：`observe.BuildDiagnosticsBundle(BundleOptions)`（`internal/observe/diagnostics.go:62`）。
- 状态映射另有一枚：`ErrorClass.MappedStates()`（`:131-140`）把 class 映到 D43 BallState 名。
ⓑ ❗ `HumanText`（"人话文案"）**无源**：全仓**没有** class→中文用户文案的映射
（现量：`grep` "人话"/`UserMessage`/`classText` 于 `internal/observe|internal/agent` **零命中**；
唯一像的 `internal/llm/retry.go:70 "重试中 (n/3)"` 是重试进度语，不是错误文案）。
⇒ 归口：**另开票**（写 17 枚文案＝产品决策＋人工过目），**不是**"扩快照"能顺手带过的。
按 §2.0(3)，`HumanText` **不许进落地集**（拿 class 字符串或空串顶上＝尺上 ❌"只说出错了"那一格）。
ⓒ 无字段＋无泵＋**一枚待写文案**。⚠ 派单口径提醒：§47.2.2 那枚"附件报错被静默吞"虽然被 §47.1 挂在行 10，
它其实是 `composer.tsx:63/92` 与 `App.tsx` 未传 `onUserError` 的事，**与本行字段无交集**（见 §1.1）。
ⓐ 备注：尺上第三枚动作「切文字模式」**无对应入站指令**——`agent.ControlVerb` 词表只有
`stop/cancel/repeat/louder/confirm`（`internal/agent/control.go:22-30`）。它属 ⓑ（面板→Go）那一侧，见 §3。

### 2.12 行 11 · Stuck（`:3485`，C22 梯度刹车）

ⓐ `RunState.Stuck *StuckView json:"stuck,omitempty"`；
`StuckView{Tool string "tool"; Repeat int "repeat"; Level int "level"; Text string "text"; Brake string "brake"}`。
ⓐᐟ **有源，且尺上那句话逐字就有生产者**：`agent.Reminder{Level,Repeat,Tool,Text,Stuck}`
（`internal/agent/guard.go:47-58`）——`Tool`＋`Repeat`＋`Stuck` 三枚正是尺上
「`web.fetch` 已连续调用 8 次，参数相同」的形状；阶梯 `[3,5,8]` 在 `guard.go:14-16,28-30`。
真发两处：`loop.go:470-473`（`EvReminder`）与 `:475-478,883-886`（`EvStuck`，带 `RepeatLevel`）。
刹车归属另有 `Result.Brake`（`loop.go:96`，词表 `guard.go:63-69` `repeat/token_budget/rounds/tool_timeout`）
与 `ReminderLevels []int`（`:98`）。消费端已存在：`cmd/wisp/run.go:689` 那一族 switch。
ⓑ 无。**尺上三按钮（继续/换个方式/放弃）不属本行字段**：它们是入站指令，
而 `ControlVerb`（`control.go:22-30`）里**只有 `repeat` 与 `stop/cancel` 近似**，"换个方式"无词表项
⇒ 归口 §3 的 ⓑ 那一侧（另开票），**不许**为凑按钮去 `Snapshot` 里塞三个布尔。
ⓒ 无字段＋无泵。**源的两侧都在跑**（loop 真发、控制台真打），这一格是"只差泵"第二例（第一例是行 7）。

### 2.13 行 12 · 成本（`:3486`，C23）

ⓐ `Snapshot.Cost CostView json:"cost"`；
`CostView{TokensIn int "tokensIn"; TokensOut int "tokensOut"; CachedTokens int "cachedTokens"; Micros int64 "micros"; Currency string "currency"}`。
ⓐᐟ **五枚全有源，一一对上尺上的 `IN`/`OUT`/`CACHED` ＋ `tok` ＋ 金额**：
- `agent.Cost{Usage, ToolCalls, Micros, Currency}`（`internal/agent/cost.go:23-30`）；
- `llm.Usage{InputTokens,OutputTokens,CachedTokens}`（`internal/llm/events.go:92-96`）——**正好三枚**；
- 定价 `PriceMicros`（`cost.go:50-63`）、`Result.CostMicros/Currency`（`loop.go:92-93`）、
  现成句子 `Guard.Consumption()`（`guard.go:207-213`，"本次任务已用 X token / 约 Y（估算 Z）"）；
- 历史/聚合读口：`Store.CostDay`（`internal/memory/dao_misc.go:154`）、`Store.ListCostDaily`（`:169`）、
  写口 `BumpCostDay`（`:132`）、任务清单 `Store.ListTaskLogs`（`dao_tasklog.go:93`）。
ⓑ **不缺源，缺一条单位口径**：`cost.go:16-20` 逐字登记着
"config.Price is documented as **micro-USD** per 1M tokens while task_log.currency defaults to **CNY**…
it never invents an exchange rate"。⇒ 尺上那个 `¥0.31` 的 `¥` **正是这枚雷**：
数字是 micro-口径、标签是 `Config.Currency` 原样。本件判：**字段可以扩，`¥` 不许由面板自己加**；
`Currency` 必须与 `Micros` 成对下发、由 Go 侧原样带（`Result.Currency` `loop.go:93` 已经是这么做的）。
归口：单位定案（工单原话"open question registered in the ticket report"，`cost.go:16-17`）——**等哪一切片？本件说不出**，
只登记为"AC#2 落地时必须带着 `Currency` 走，单独落 `Micros` 会造出假 ¥"。
ⓒ 无字段＋无泵＋**一枚单位口径未定案**。

### 2.14 行 13 · 已取消 / 被打断（`:3487`）

ⓐ `RunState.Status string json:"status"`（或 `Stopped string json:"stopReason"`）。
ⓐᐟ **四枚枚举全有源**：`agent.StatusCancelled`（`internal/agent/loop.go:46`），
真发三处 `loop.go:378,422,494`（都带 Message「任务已取消」）；`llm.StopCancelled`（`llm/events.go:77`，
判定在 `loop.go:417`）；`observe.ClassCancelled`（`observe/errors.go:30`）；
`agent.OutcomeCancelled`（`journal.go:42`）。⇒ 与 §47.1 行 13 那句"`done` 是写完了、不表达被取消"一致：
`ResultChunk.Done`（`composer.go:55`）确实表达不了，**但上游有四处能表达**。
ⓑ 无。**一处必须挑明的取向**：`agent.Event` 虽然有 `Stop llm.StopReason`（`sink.go:52`），
`EvDone` 那枚发布（`loop.go:935-938`）**只填 Kind/TaskID/Text/TokensIn/TokensOut，不填 Stop**；
今天取消与完成在**事件面**上只靠 `Text` 里那句"任务已取消"区分，`Result.Status`（`loop.go:84`）
只在**同步返回**里。⇒ AC#2 落地时必须选一条：**给 `EvDone` 补 `Stop`/`Status`（改 `loop.go:935`）
还是从 `Result` 组快照**。本件判前者更省，但**这是取向，不是读数**，交编排者定。
ⓒ 无字段＋无泵（＋一条"走哪条路"待定）。

### 2.15 行 14 · 注入检出（`:3488`，C25/R4）

ⓐ 若要结构化：`Pending[].Taint *TaintView json:"taint,omitempty"`，`TaintView{Source string "source"; Channel string "channel"; Fragment string "fragment"}`。
ⓐᐟ **来源这一枚，今天就抵达面板**（§1.3 已证）：
`risk.Hit{ScopeID,Channel,SrcTool,Origin,Fragment}`（`internal/risk/provenance.go:246-253`）＋
`Hit.Source()`（`:257-262`，注释逐字"the `<源>` the confirmation card must name（"本次操作包含来自 `<源>` 的内容", SPEC-06 §5 / R4 reason）"）
→ `rules_gateway.go:112` 产 `reason` → `approval.go:83` verbatim → `l2-approval-card.tsx:159`＋`:208` 真渲染。
ⓑ ❗ `Fragment` **无源可给面板**，两道独立证据：
① 它自己的字段注释划了界：`provenance.go:251-252`
"the matched fragment (**native confirmation card only — NEVER write this into the DB/args_json: the fragment IS sensitive text**)"；
② `risk.Hit` **在 `internal/risk/**` 之外零消费者**（现量：`grep -rn "risk\.Hit|ScopeTaints|TaintHit"` 排除 `internal/risk/` 与测试后**无输出**）；
③（旁证）`approval` 自己的面板投影 `PanelItem` 注释（`ui.go:46-49`）逐字写着面板可见面
"There is deliberately no grant field, no Allow capability **and no way to name a source**"。
⇒ 归口：**等人工批准**——把敏感片段放进 WebView2 快照是**安全取向**，不是字段增补。
尺上"可展开看命中的片段"这一格今天**做不到也不该硬做**。
ⓒ 三分：来源＝**已通（缺的只是别盖住它）**；结构化 `Source/Channel`＝无源（`Hit` 不外露，要 `risk` 侧开口子，
⚠ 而 `internal/risk/**` 在本票零命中面上 ⇒ **本程不碰、AC#2 也不该顺手碰**）；片段＝**等人工批准**。
⇒ **本行是 §2 里唯一一枚"ⓐ 判已有、ⓑ 判禁入"同时成立的**，读表时别被 §47 那句"没有任何承载字段"带偏。

---

## 3. "当前屏"这一维（AC#3）：ⓐ 与 ⓑ 是两件事，且 ⓐ 单独存在**换不了屏**

### 3.1 先把两件事命名清楚

| | 方向 | 谁发起 | 今天的状态 |
|---|---|---|---|
| **ⓐ** | **Go → 面板**：快照里带"现在该显示哪一屏" | 宿主 | **字段不存在**（`Snapshot` 四枚键里没有 `view`，`composer.go:44-49`） |
| **ⓑ** | **面板 → Go**：用户点了导航条，请 Go 换屏 | 用户 | **曾经存在、已被删除**（`Q-50` 甲，commit `f1cdafa`） |

⚠ 本票按 AC#3 只做 ⓐ。下面 3.2 与 3.3 分开写，**3.4 是硬答案**。

### 3.2 ⓐ Go→面板：唯一一枚"TS 侧不用改"的字段——但**值无源**

**消费端今天就位**（这是本件在 ⓐ 上的新读数，工单没写）：
`frontend/src/lib/panel-views.ts:145-153 currentView(snapshot)` **已经在读 `snapshot.view`**——
签名是 `PanelSnapshot & { view?: string }`（`:145`），`named === undefined` 时回落到
`DEFAULT_VIEW = "chat"`（`:127`，注释逐字"it is a fallback for a missing field, **not a memory** of where the user last was"），
命中 `PANEL_VIEWS` 就返回（`:148`），不命中则 `console.error` 报契约漂移再回落（`:149-152`）。
`App.tsx:71 const view = currentView(snapshot)` 是**唯一**读取点，全树无组件自己持有或推导 view
（`panel-views.ts:142-144` 逐字："no component holds a view, and none derives one either"）。
⇒ **Go 侧加一枚 `View string json:"view"`，前端一行都不用动就能被读到。**（对齐 `panel.ts:129-137` 时另计。）

**但值本身无源，且没有任何东西能校验它。** 现量两路：
1. `grep -rn 'palette|privacy|"tasks"|"cost"|"security"' --include=*.go internal/panel/ cmd/wisp/`（排除测试）
   → **零命中**。⇒ **Go 全仓不知道这九枚屏的存在。**
2. 九枚 id 的出处是**纯前端表**：`panel-views.ts:32-41` 的 `PanelViewId` 联合类型与 `:68-122` 的 `PANEL_VIEWS`，
   其头部注释（`:6-7`）逐字说它们"taken entry-for-entry from the demo's own list"，
   现量核实于 `design/doubao/demo/app.js:17-27`（九行 `NAV_ITEMS`，id 与前端逐枚相同）。

⇒ ⓐ 落地形态只能是"`View string`、值由宿主随便给"：**一枚无校验的自由字符串进入 C17 对外形状**，
拼错的表现是 `panel-views.ts:151` 那句 `console.error` ＋ 回落 chat，**不是一枚测试红**。
⇒ **本件判：ⓐ 与 §2 行 4 的 `IconClass` 同性质——词表不在仓里。**
建议（**不替编排者定，只是把代价摆出来**）：ⓐ 要落就得同时在 Go 侧立一枚具名枚举
（`internal/panel` 里一份 `ViewNames()`，与 `panel-views.ts:32-41` 用 `TestComposerContractTypesMatchFrontend`
同一条尺双向钉住，`internal/panel/composer_test.go:48-81`），否则 ⓐ 比 §2 里任何一枚 ❗ 更坏：
❗ 是"值假"，ⓐ 是"**名字都没定**"。

### 3.3 ⓑ 面板→Go：被 `Q-50` 甲删掉的那条，且它**不是前端能自己加回来的**

- 出站名册现量：`internal/panel/bridge.go:97-103 knownComposerMethod` 只答 **4 枚**
  （`MethodModeRequest`/`MethodWorkspaceRequest`/`MethodAttachmentAdd`/`MethodMessageSend`）。
- 全树 `panel.view` **仅存 1 处＝`frontend/src/lib/panel.ts:228` 的注释**（解释"为什么故意没有第五枚"）。
- 门存在且**已经在响过**：`internal/panel/composer_test.go:521-523`
  ——`"the renderer names a route the Go side does not answer"`。`panel.ts:230-233` 记着它的战绩：
  "a fifth means the renderer names a route the native side does not answer - which is exactly what
  internal/panel/composer_test.go:522 exists to catch (**it was red on two tests before this comment was written**, green after)"。
- 用户侧的体感：`frontend/src/components/nav-rail.tsx` 每枚按钮 `aria-disabled`（`:61`）、`disabled`（`:67`）、
  `title={`${v.label}（换屏待接线）`}`（`:69`）。⇒ **点下去不会动，且它自己说出为什么。**
- ⚠ **ⓑ 不是"等泵"，是"等人工批准"**：`panel.ts:234-235` 逐字
  "Re-adding one is not a frontend change: **it reopens C17's method whitelist, a contract-level human-approval surface**"。
  而 `C17 方法白名单定稿`正列在 AGENTS §2／`SPEC-12 §4.2` 的**"未定义即停"清单**里
  （原文："`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿（S7/S5 切片卡批准）"）。
  ⇒ **本程不申请、不假设、不改工单**，只把"ⓑ 撞那闸门"这一条钉在这儿给下一位读。

### 3.4 硬答案（AC#3 要求的那一句，不许读成"换屏已通"）

> **只做 ⓐ，界面真换不了屏。**两条独立理由，任一成立就够：

1. **快照根本还没到过面板。** §2.0(1) 现量：`panel.NewSnapshot`（`composer.go:63`）**非测试调用者 = 0**，
   全树 `Snapshot{` 字面量只有两处、**都不是面板的那枚在被人构造**：
   `internal/panel/composer.go:79`（就是 `NewSnapshot` 自己的 return）、
   以及 `internal/perm/store.go:247`——后者返回的是 `internal/perm` 里**另一枚同名类型**
   `perm.Snapshot{Mode, History}`（`store.go:236-239`），与 `panel.Snapshot` 无关，只是**同名**。
   ⚠ 提这一枚是为了让下一位别把它误读成"perm 侧已经在造面板快照"；
   顺带它确是 `Composer.Mode` 的真上游（`Store.Snapshot()` `:242`），但那是票 92 已经接好的路，不是新读数。
   ⇒ 今天没有任何一条推送会把 `view` 送进去；`currentView` 会**永远**走 `panel-views.ts:147` 那行回落，
   屏**一直停在 chat**。加 ⓐ ＝ 加一枚**永远读不到的键**。
2. **即便泵存在（票 35 落地后），ⓐ 也只让 Go 有权"切到某一屏"，用户的点击仍然无处可去。**
   换屏这个动作的两个发起方是宿主与用户；ⓐ 只覆盖前者。
   后者要的那条路由在 §3.3：名册只有 4 枚、第五枚被 `f1cdafa` 删、加回来撞 C17 白名单人工批准。

⇒ **请照这句写进任何下游读数**：
**「ⓐ 落地后界面仍点不动；缺的是 ⓑ，而 ⓑ 不是泵（票 35）能顺带解决的——它是 `C17` 方法白名单，属人工批准。」**
工单 AC#3 原话是"ⓑ **留给泵（票 35）**"，本件据 3.3 现量**部分推翻这个归口**：
泵能送达 ⓐ，泵**造不出** ⓑ 那条被删掉的路由（那是白名单，不是推送）。

### 3.5 顺带一枚，且这枚对本票有用：`panel.resync` **不是"全仓找不到"**，它已在 `SPEC-08` 的推送表里挂了名

`panel-views.ts:136-138` 断言"`panel.resync` … is still **zero-hit repo-wide**"；
`frontend-session-log.md:684` 写"全仓 0"；`frontend-session-brief.md:230` 写"**全仓零命中**（规格里有名字、盘上没实现）"。

现量（`grep -rn "panel\.resync"`，排除本件自身）：**13 处命中、分布在 5 类位置**。
"zero-hit repo-wide" 这个写法**不成立**，成立的是它后半句"规格里有名字、盘上没实现"：

| 位置 | 读数 | 判 |
|---|---|---|
| **Go 代码**（`*.go`） | **0 处** | ✅ 未实现，这条为真 |
| **TS/TSX 代码** | **1 处，且是注释**：`frontend/src/lib/panel-views.ts:137`（就是那句自证自己） | ✅ 无实现 |
| **规格** | `docs/specs/SPEC-08-ui-ball-panel.md:150`（"每次 `show` Go 侧推 `panel.resync` 全量状态"）、**`:163` 推送表里独立一行 `| panel.resync | Go→前端推送 | — |`**、`docs/PLAN.md:2966`（"④ 前端无状态的强制手段：面板每次 show 必须发 `panel.resync`"） | **已命名、已规定** |
| **工单** | `.scratch/wisp/issues/35-panel-bridge-c17.md:14,18,29`（`:18` 把它列进 **C17 方法白名单**、`:29` "on EVERY show"）、`issues/34-frontend-scaffold.md:36` | **票 35 名下已立项** |
| **台账/报告** | `docs/reports/pending-and-issues.md:6009,6012`、`frontend-session-log.md:501,684`、`frontend-session-brief.md:230` | 已登记 |

⇒ **这一处口径差别对下一位是实质性的**：ⓐ（快照带"该显示哪一屏"）的**运载工具已经在 `SPEC-08:163` 挂名、
在 `PLAN.md:2966` 被定为强制**，它不是 §2 里那种"等哪一枚票"的悬案，而是**票 35 名下已立项、尚未落地**的那根线。
⇒ 所以本件对 ⓐ 的判语要读成"**泵是已立项的，缺的是泵里的这一枚键＋Go 侧那枚屏枚举**"，
不是"没人打算接"。

⚠ 顺带把 `SPEC-08:150-151` 那条"前端必须无状态"与 ⓐ 的落地形状钉在一起，免得下一位走出第三种形状：
`PLAN.md:2966` 与 `SPEC-08:150` 规定的是"**每次 show 全量推**"，
对应 `App.tsx:71` 每次从快照现读、`panel-views.ts:127` 那句"not a memory of where the user last was"。
⇒ ⓐ 落地时 `view` **必须由 Go 每枚快照都带上**；前端记住上一次的值会同时违 `SPEC-08:150`、`PLAN.md:2966` 与票 92 的 AC#4。
