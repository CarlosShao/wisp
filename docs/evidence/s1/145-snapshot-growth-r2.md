# 票 145 r2 —— 快照生长（写腿）：现量、落地判定与**停手报回**

> **本件的射程**：派单 `.scratch/wisp/dispatches/2026-09-28-174x-impl-145-r2-snapshot-growth-with-nail-preapproval.md`
> 的三格（`AC#2` 落地有源集／`AC#6` 反向判据／`AC#2b` 两维），加派单 `起手必做` 的四条现量，加前端那份
> "应当长这样"的逐键清单。**票面三格本程一律不自勾**（勾要非实现者表）。
>
> **一句话结论（先说，免得被读成"做了半程"）**：**本程落地集＝空，且本程一行生产码都没改。**
> 不是找不到源——**是"填真值"那一跳不在这次批准的写面里**：快照的值全部来自装配根的读口，
> 而装配根是 `cmd/wisp/run.go:442`（派单具名放开面＝`internal/panel/composer.go`＋`internal/panel/pump.go`＋
> `cmd/wisp/panel_pump.go`＋`pump_test.go` 的 `:124`／`:276`；`run.go` 不在其中）。
> 第二把锁：`AC#6` 要"每一枚新字段答得出**哪一枚用例**断言它来自真源"，而**写面里没有一枚可落断言的测试文件**
> （`pump_test.go` 只放开那两行）。⇒ 按派单自己第 1 条"答不出的字段直接不要"，**任何一枚字段本程都不该进树**。
> 两把锁都是 r1（`docs/evidence/s1/145-snapshot-fields-landed-r1.md` §8 的 R2／R1）具名报过、
> 编排者 09-26 明确"暂不给"的那两枚；**这一程的派单没有解开它们**——这就是停手的内容，具名在 §4。

---

## 0. 锚点 · 起手名册 · 派单要求的四条现量

### 0.1 进场现读（不抄派单）

```
$ date                                  2026-09-28 17:42 +08
$ git log --oneline -1                  318d4bb3 33-r3(片④): 证据件 33-inbound-listener-r3.md ＋台件 probes/33/r3/**
$ git branch --show-current             dev
```

派单给的编排者锚 `318d4bb3` 与进场实测**一致**（派单自己写明"编排者锚≠你的锚"）。
**锚点在本程内漂过两次，都不是本程所为**：`318d4bb3` → `38fc7c0e`（台账 `A387`–`A388` 收 33-r3＋派本程）
→ `971b7dad`（本程取门禁读数时又落了一枚）。⇒ 本件凡引行号，**按本程那一版树现读**，见 §0.3。

### 0.2 起手名册（`git status --porcelain` 逐枚抄，17:42）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt … flip-6.txt … flip-baseline.txt … flip-restored.txt   (8 枚)
 D design/assets/base.css  D design/assets/icons.js  D design/assets/theme.js  D design/assets/tokens.css
 M design/doubao/README.md  M design/doubao/demo/app.js  M design/doubao/demo/index.html  M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/{approval,ball,chat,config,cost,firstrun,palette,privacy,security,states,tasks}.html   (11 枚)
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/  ?? .scratch/wisp/dispatches/2026-09-28-174x-impl-145-r2-…md  (本程派单，已入库 38fc7c0e)
?? .scratch/wisp/probes/{139/accept-r1,156/…,158/r2,161/r2,161/r5,162/r4,176/r1/logs,183/accept-v1,185,33/r2,999}
?? .zcodeignore
?? design/doubao/01-ball-states.jpg  design/doubao/demo/{lib/,rb-*.js,rightbar.js,sidebar.js,screens/home.js,screenshots/}
?? design/old/
?? part1-state1-{fixed,pristine}.txt part1-state2-{fixed,pristine}.txt part2-nog6-{fixed,pristine}.txt part3-stale-{fixed,pristine}.txt
```

⚠ **派单点名的别家脏件，本程现量与名册不符（照实报，不据此扩权也不据此删）**：
派单正文与编排者的话都写着别家在飞的是 `internal/risk/**`、`internal/tools/pointer_185_cli_seam_test.go`、
`cmd/wisp/panel_inbound*.go`——**起手 `git status --porcelain` 里这三组一枚都不在**（`internal/risk/**` 干净、
`pointer_185_cli_seam_test.go` 不在名册、`cmd/wisp/panel_inbound.go` 已由 `6609e7e1` 提进树，
工作树里只剩它的判据件 `panel_inbound_33_test.go` 也在树内）。
⇒ 现量：`git status --porcelain -- internal/panel/ cmd/wisp/ docs/PLAN.md docs/specs/ go.mod go.sum` = **空**
（本程的写面起手就是干净的，且别家没在里面写）。
**终态名册**在本件末节 §6.3；本程新增只有 `docs/evidence/s1/145-snapshot-growth-r2.md`＋`.scratch/wisp/probes/145/r2/**`＋
票 145 的 Progress log 追加，全部逐枚具名 pathspec 提交。

### 0.3 派单 `起手必做` 第 2 条那四条，逐条现量

| # | 派单给的尺 | 本程现量 | 判 |
|---|---|---|---|
| 1 | `grep -c "" internal/panel/composer.go` | **290** | ✅ |
| 2 | "`PanelSnapshot` 当前字段枚数（用结构体数，别抄票面）" | Go 侧那枚类型叫 **`Snapshot`**（`internal/panel/composer.go:57`），字段 **恰四枚** `Pending`／`Results`／`Composer`／`GeneratedAt`（`:58-61`）。`PanelSnapshot` 是 **TS 侧**的 interface 名（`composer_test.go:60` 的配对行 `{"Snapshot", Snapshot{}, "PanelSnapshot"}` 就是把这两个名字对接起来的那行） | ⚠ 成立但**名字要更正**：Go 侧无 `PanelSnapshot` 这枚标识符 |
| 3 | `grep -rn "PanelSnapshot{" --include=*.go internal/ cmd/ \| grep -v _test.go` ＝**生产构造点枚数** | **0 枚**（`rc=1`，无输出）。按语义改成真身那枚：`grep -rn "Snapshot{" … \| grep -v _test.go` ⇒ **5 枚命中**，其中**面板的只有 3 枚且都不是装配点**：`composer.go:97`（`NewSnapshot` 自己的 return）、`pump.go:228`（无泵时的零值＋错误）、`cmd/wisp/panel_pump.go:248`（nil receiver 的零值读口）；另两枚是**同名不同物**（`internal/perm/store.go:247` 的 `perm.Snapshot`、`internal/ball/liquid.go:276` 的 `motionSnapshot`）。真·生产构造链＝`cmd/wisp/run.go:442` 装配 → `internal/panel/pump.go:207` 调 `NewSnapshot` → `composer.go:97` return | ❌ **派单这条尺按字面量是空的**；具名构造点见下表 §1-P2 |
| 4 | `go test -count=1 ./internal/panel/`＝恰三枚在册红 | `rc=1`、`=== RUN` **137**、`--- PASS` **81**、`--- FAIL` **3**、`--- SKIP` **0**；红名册逐名＝`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree` | ✅ **与派单逐名一致**（本件 §6.1 留原始件） |

**真·生产构造点具名（派单 AC#2 要的那张"谁、几处"表，全部 `grep` 现量）**：

| 层 | `file:line` | 射程 | 在不在这次放开面里 |
|---|---|---|---|
| 装配根（**唯一一枚**） | `cmd/wisp/run.go:442`（`rt.pump = panel.NewSnapshotPump(panel.PumpSources{`），读口行 `:443-448`＝`Verdicts`／`Mode`／`Workspace`／`Git`／`Results`／`Out` | 决定"每一段有没有人喂" | ❌ **不在** |
| 快照构造 | `internal/panel/pump.go:207`（`return NewSnapshot(cards, results, composer, now())`）；`pump.go:158` 是无泵兜底那一发 | 把读口喂进来的值组装成包 | ✅ 在 |
| 段构造 | `internal/panel/composer.go:97`（`Snapshot{…}` return）、`composer.go:224`（`ComposerState{…}` return） | 字段落位 | ✅ 在 |
| 出口 | `cmd/wisp/panel_pump.go:182 bookPanelSnapshot`／`:264 publishPanelSnapshot` | 包构造完的去向（账本＋`lastSnap`） | ✅ 在 |
| 全树 `NewSnapshotPump(` 生产调用者枚数 | **1**（`cmd/wisp/run.go:442`）；其余 7 枚全在 `_test.go`（`internal/panel/pump_test.go` 5 枚＋`git_test.go` 2 枚） | — | — |

---

## 1. 派单与票面给的前提，逐枚复量（不照抄）

| # | 前提（出处） | 现量 | 判 |
|---|---|---|---|
| P1 | "`PanelSnapshot` 那四字段与生产构造点"都在 `internal/panel/composer.go`（派单 §解冻面） | 字段在 `composer.go:57-62`；**构造点不在 composer.go 里被喂**——composer.go 只是 return，值从 `pump.go:207` 的实参来，实参从 `run.go:442` 的读口来 | ⚠ **成立但不足**：按这一句的字面做，加出来的键今天**没有值**（与 r1 §1 P2 同一枚，r1 已用 E1b 实测过：`{"depth":0,"windowMs":0,"vetoChannels":null}` 而同包的 `pending`/`results` 是活的） |
| P2 | "在生产构造点填真值"（派单要落的三格 ①） | 唯一生产装配根 `cmd/wisp/run.go:442` **不在放开面**（`A388` 那句"在飞＝`145-r2`（写·composer.go＋pump.go＋panel_pump.go＋pump_test.go 那两行）"逐枚点名，无 `run.go`） | ❌ **不成立** ⇒ 停手，见 §4 R-1 |
| P3 | 解冻 `pump_test.go:124`／`:276` 两行（派单＋`A388`，理由＝加字段必被打红） | 两行现量仍在（本程逐字读 `:109-126`、`:264-277`），内容确为写死名单 `!= "composer,generatedAt,pending,results"`；**且 `:111-112` 的注释自己写着"the four keys are the contract"** | ✅ 成立。**本程未改这两行**：本程零枚新字段 ⇒ 这两行**不会红**，改了等于把"恰好四枚"换成"随结构体"而**什么都买不到**（派单自己那一句"没有正控＝把门换成好看"在此更严：连要保护的第五枚键都没有）。改动方案与正控形状写进 §4 R-2，留给有字段的那一程 |
| P4 | 在册三枚红（派单硬禁区） | §0.3 第 4 行：逐名一致 | ✅ 成立。本程未修、未放宽、未豁免；`TestComposerContractTypesMatchFrontend` 现在的红句仍只有一枚键名 `git`（`composer_test.go:74`，票 181 那枚 `ComposerState.git`），**本程没让它更红也没让它变绿** |
| P5 | `AC#2b` ① "当前可选模型清单，来源＝配置目录 ∩ `enabled=true`" | 目录与位都在：`internal/config/schema.go:345`（注释逐字 "ModelSpec is llm.providers.<name>.models.<model-id>: one catalog entry"）、`:348 type ModelSpec struct`、`:369-371`（`Enabled bool toml:"enabled" default:"true"`，注释"removes the model from discovery/selection"）、`:400 Models map[string]ModelSpec`。装配这枚目录的对象**本进程就有**：`cmd/wisp/run.go:209 rt.cfg *config.Config`。**但"∩ enabled 后拿给人看"这件事今天没有任何生产读者**（`grep -rn "ModelSpec" --include=*.go internal/ cmd/ \| grep -v _test.go` ⇒ 只有 `schema.go`／`validate.go:179,193,195`／`internal/llm/discover.go:108,117`——`discover.go` 那两处是**写**目录，不是列目录） | ⚠ **有源、无生产者**（普查 §0.4 的第四态）。落到快照仍要 `run.go` 那一行读口 |
| P6 | `AC#2b` 票面点 5 "正解形状＝复用 `cmd/wisp/models.go:104 cmdModels` 的取数" | 现量：`cmdModels`（`:104`）的 `list` 腿 `modelsList`（`cmd/wisp/models.go:205-225`）打的是 **`store.manifest.Models`**——签名验过价的**下载清单**（`id/purpose/quant/sizeBytes/license/status`），**不是** `llm.providers.<name>.models.<id>` 那枚配置目录；`modelsList` 上面 `:206-208` 的注释逐字写着"listing is a question about the manifest, **not about this boot's [models] section**" | ❌ **票面这一枚点错对象**：两枚"目录"同名不同物。要落 ①，取数得新写（读 `rt.cfg.LLM.Providers[*].Models` ∩ `Enabled`），**复用 `cmdModels` 复不出来** ⇒ §4 R-3 |
| P7 | `AC#2b` ② "第一维＝`Capabilities.Thinking`（`internal/panel/pump.go` 侧的快照位）" | `internal/config/schema.go:324 type Capabilities struct`、`:330 Thinking bool`；位是**人填的**（`:322` 注释："声明可由人填，但必须由探测核实（票 11），探测结果本身住在 SQLite `provider_health`、这里刻意不留字段"）。真源那一侧现量：探测确实**在写**（`internal/llm/probe.go:139` 用 `ThinkingIntensity:"low"` 探；`internal/llm/probe_health.go:112 HealthSink`；生产 sink＝`cmd/wisp/run.go:700-704 storeHealthSink`，装配点 `cmd/wisp/providers.go:194`），SQLite 读口 DAO **存在但零枚非测试调用者**（`internal/memory/dao_providerhealth.go:170`／`:178` 两枚 SELECT，`grep ProviderHealth --include=*.go internal/ cmd/ \| grep -v _test \| grep -v dao_providerhealth` ⇒ 无输出）。而发现式录入（`discover.go:116-121`）把 `Capabilities` 全填 false ＝"unknown until probe" | ⚠ **半有源**：配置声明位有源（人填），**已核实那一版有数据、无读者**。快照送出去只能是"人填的那一版"，这一点必须在票面写清，否则界面把"未探测"画成"不支持思考" |
| P8 | `AC#2b` ② 第二问"该适配器实际接受的档位词表"——票面点 4 要求"逐家量出来，不许把四档原样送出去当作每家都支持" | **本程逐家量了**（见 §3.2 那张表）：配置词表四枚（`schema.go:85-88` 导出常量 `off/low/medium/high`，内部清单 `schema.go:98 thinkingIntensity`，枚举校验 `validate.go:202,222,231`）；`openairesponses/request.go:63 effortLevels`＝**三枚**（`low/medium/high`），`off` 与 `""` 走 `:127` 的外层判断＝**不发这个参数**（`:128-131` 对未知值报 `unknown thinking_intensity`）；`anthropic/request.go:92 thinkingBudgets`＝**三枚**，`off`/`""` 同样不发参数（`:112-116`）；`openaichat`＝**零枚**——`adapter.go:14-17` 逐字"Request.ThinkingIntensity is intentionally NOT mapped: the chat completions wire has no portable field" | ❌ **这一维今天没有可导出的源**：两张映射表都是**包内私有 `var`**（`effortLevels`／`thinkingBudgets`），全仓**没有任何导出访问器**（`grep -rn "func .*(Levels|Intensities|AcceptedThinking|Capabilities())" internal/llm/ internal/config/ \| grep -v _test` ⇒ 只命中 `llm.AllProbeCapabilities`，那是探测能力名，不是档位词表）。⇒ 快照要带"这家接受哪几档"，只有两条路：**在 `internal/llm/*` 导出访问器**（不在放开面、且属另立契约），**或在 panel/main 另抄一份词表**＝普查 §0.4(2) 禁的"第二份词汇表"（它会在适配器改表时静默漂移）。⇒ **本程不落，且这一格必须摆给编排者定形** |
| P9 | "别人在飞／脏件不碰"名册（派单硬禁区） | 见 §0.2：点名的三组此刻**不在工作树**；本程对所有别家路径**零读写** | ✅ 照做（附一条读数更正，不改纪律） |
| P10 | 不新增任何 `panel.*`／方法名；不做审批出口 | 本程零码改动 ⇒ 结构性未犯。现量入站名册仍 4 枚：`internal/panel/bridge.go:104-109 knownComposerMethod`（`case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend` 在 `:106`，本程只读） | ✅ |

---

## 2. `AC#2`＋`AC#6`：逐维判（"真来源在哪一行"与"答不出的就不要"）

**口径**（派单第 1 条＋普查 §0.4 原尺）：一列"真源"只认①该包今天就算得出来、②今天已有人生产；
"答句"＝**哪一枚用例断言它的值来自真源**。
"喂法"＝要从真源抵达 `Snapshot`，需要动哪一枚文件；❌＝那枚文件不在放开面。

| 维（候选段） | 真源 `file:line` | 喂法要动的文件 | **AC#6 答句** | 判 |
|---|---|---|---|---|
| `composer.models`（AC#2b ①） | 目录∩`Enabled`：`internal/config/schema.go:348,369-371,400`；数据在 `cmd/wisp/run.go:209 rt.cfg` | **`cmd/wisp/run.go:442`**（新读口行）❌ ＋ `pump.go:191` 段构造 ＋ `composer.go:205` 新字段 ＋ 一枚 `_test.go` ❌ | **答不出**：写面内无可落断言之处（`pump_test.go` 只放开 `:124`／`:276` 两行，加不进新断言；`cmd/wisp/panel_pump_test.go` 不在写面） | **不要**（装饰品——加了键也没人喂） |
| `composer.efforts`（AC#2b ②） | 声明位 `schema.go:330`；档位词表 `openairesponses/request.go:63`／`anthropic/request.go:92`／`openaichat/adapter.go:14-17`（**全为包内私有或"不映射"**） | 除上面那些，还要 `internal/llm/**` 导出访问器 ❌ | **答不出**，且**源本身不可导出**（P8） | **不要**，并报形（§4 R-3） |
| `approval.windowMs` | `approval.Gate` 的窗口（r1 实测 `2000ms`，`internal/agent/approval/gate.go`） | `run.go:442` ❌（或把 gate 状态塞进已有读口的返回结构＝**改语义**，批准只到"新增字段"） | 答不出（同上两把锁） | **不要** |
| `approval.vetoChannels` | `ChannelRegistry.Statuses()`（`internal/agent/approval/approval.go:190`）＋ B1 原句（`:103-106`） | 同上 ❌ | 答不出 | **不要** |
| `approval.remainingMs` | ❗ **无源**（r1 P6 已把普查甲组这一枚挪进乙组：`approval.EventTick` 声明了全仓零生产者；倒计时本体是 `gate.go:257 clock.After`，无"还剩多少"读口） | — | 答不出：任何实现只能填静态窗口长 | **不要**（宁缺毋造） |
| `run.reasoning` | `agent.EvReasoningDelta` 三适配器真发、`internal/agent/loop.go` 已转发；**汇在 `cmd/wisp/run.go:800-806` 被打印但没被记**，那行注释自己就写着"the reasoning field is ticket 145 row 3, which is a key the snapshot does not have" | `run.go`（记录点＋读口）❌ | 答不出 | **不要** |
| `run.usage`（实时） | `llm.EvUsage` 三适配器真发，但 agent 侧 switch 不转 | `internal/agent/loop.go` ❌（票 153 地界）先行，再 `run.go` ❌ | 答不出 | **不要** |
| `run.status`（取消／完成） | `agent.StatusCancelled` 真发三处；`EvDone` 不填 `Status`／`Stop` | 取向未定（改 `loop.go` 还是从 `Result` 组包）＋ `run.go` ❌ | 答不出 | **不要** |
| `tools[]` | `memory.ToolCall`（name/outcome/riskLevel/argsJson/correlationId，读口 `internal/memory/dao_toolcall.go:98,111`）；taskID 只活在 `run.go` 的 sink 与 `StreamLog` 的键里 | `run.go` ❌（要把 taskID 随包带出） | 答不出 | **不要** |
| `cost` | `agent.Cost`／`llm.Usage`／`Result.CostMicros/Currency`；**单位口径未定案**（`internal/agent/cost.go:16-20` 自陈 micro-USD vs CNY，"never invents an exchange rate"） | `run.go` ❌ ＋口径待裁 | 答不出 | **不要**，且**绝不许面板自己加货币符号** |
| `failures[]`（除人话文案） | `observe.Error.Class/Detail/ProviderCode`＋`Retryable()`；载体 `agent.Event.Err` 真发到 `run.go:816` 那一支，sink 不留存 | `run.go` ❌ | 答不出 | **不要** |
| `failures[].humanText` | ❗ 无源（class→中文文案映射不存在，普查乙组） | — | 答不出 | **不要** |
| `view`（AC#3 ⓐ） | ❗ Go 侧九枚屏 id 零命中；名字无校验 | 还要先在 Go 立枚举＋双向尺 | 答不出 | **不要**（维持票面 AC#3 的硬答案："ⓐ 单独存在换不了屏"） |
| `approval.depth` | `len(pending)` 的复述（`App.tsx:75-76` 今天已在自己算） | `pump.go` 内可派生 ✅ | 只能答"与 `pending` 长度相等"——**同一枚事实的第二份拷贝** | **不要**（装饰） |

⇒ **本程落地集＝空**，与 r1 同一枚结论、同一把尺，但**卡点更窄更明**：r1 那程数出四把锁（B 前端面／C 键集钉＋两把尺／D `run.go`／E 无测试文件），
这一程派单已经**自己解开 C**（`pump_test.go` 两行）**并把 B 判成"更红是预期后果"**；
**剩下的 D 与 E 两把没解，而这两把恰好是"填真值"与"答 AC#6"的那两跳本身**——不是可以再挑几枚字段绕开的形状。

---

## 3. `AC#2b` 两维的取数现状（给下一程的施工图，本程不落）

### 3.1 ① 可选模型清单

- 真源：`rt.cfg.LLM.Providers[<name>].Models[<model-id>]`（`internal/config/schema.go:400`）∩ `ModelSpec.Enabled`（`:369-371`）。
- 该带出的列：provider 名、model id、`Display`（`:348` 起那枚结构体的人看名）。
- ⚠ **不许把未启用的算进去**（票面原文）；`Enabled` 的默认是 `true`，所以"没写这行"＝启用，"写了 `enabled=false`"＝剔除——读法要按 TOML 解码后的字段值，**不是**按"配置里出现过"。
- ⚠ `AC#2b` 点 5 那句"复用 `cmdModels`"**复不出来**（P6：那是签名下载清单）。下一程的取数是**新写一段**，还是把 `wisp models list` 那腿改造成两个都列，属**取向**，本程不替编排者定。

### 3.2 ② 思考档位：两问分开答（本程逐家量）

| 家 | 接受的档位词表 | `off`／`""` 的行为 | 未知值 | 可导出？ |
|---|---|---|---|---|
| `internal/llm/openairesponses` | `low`／`medium`／`high`（`request.go:63 effortLevels`） | **不发** `reasoning` 参数（`request.go:127` 外层 `!= "" && != "off"`） | 报错 `unknown thinking_intensity`（`:128-131`） | ❌ 私有 `var` |
| `internal/llm/anthropic` | `low`／`medium`／`high`（`request.go:92 thinkingBudgets`＝1024／4096／8192） | **不发** `thinking` 参数（`request.go:112`） | 报错（`:113-116`） | ❌ 私有 `var` |
| `internal/llm/openaichat` | **零枚**——刻意不映射（`adapter.go:14-17`，理由："chat completions wire has no portable field，给严格端点发私有参数会 400"） | 无此参数 | — | — |

⇒ 票面点 4 要的那枚"四档 vs 三档真分歧"**量到了**，而且比"三档"更硬的一面是：**第三家一档都不收**。
⇒ 快照若带这一维，**必须是"每 provider／每模型各自的那一版"**，且"这一版从哪儿读"要先有人把词表**导出**——这是 `internal/llm/**` 的形状决定，不在这次批准里。

---

## 4. 停手报回（三枚，按一句话能否办完排序；本程一律未擅自办）

**R-1 · "填真值"那一跳在写面外：`cmd/wisp/run.go:442`（唯一生产装配根）。**
这是 r1 的 R2 同一枚锁；09-26 编排者判"不给"，理由之一是"就算给了 R2，没有 `Q-51` 照样落不了地"——
**今天这一条已经变了**：这一程的派单把 `Q-51` 那一侧改成"更红是预期后果、不必修"，所以**只差 `run.go` 那一行读口**。
一句可办完的形：批准在 `cmd/wisp/run.go:442-449` 的 `panel.PumpSources{…}` 字面量里**只加读口行**（不加逻辑、不动 sink 之外的任何一行），
并把对应的记录点（`run.go:800-816` 那三支 case 里"把已打出来的值记到 `rt` 上"）划进同一枚地界——r1 的 R3 已经量过：**泵已被这些事件驱动，缺的只是"记下来"**。
⇒ 本程**未动 `run.go` 一字**（它在写面外，且派单的写面名册里没有它）。

**R-2 · `AC#6` 的用例之家在写面外：放开面里没有一枚可写新断言的 `_test.go`。**
这是 r1 的 R1，编排者 09-26 明确"暂不给"（理由：给开门换来的可能是"看起来绿、从不跑"的断言，那枚仪器缺口没有尺）。
⇒ 于是"每一枚新字段答得出哪一枚用例断言它来自真源"这句话，**在批准的写面里结构上无定义域**（不是措辞问题、不是欠一件小账）。
本程按派单第 1 条执行：**答不出的字段一枚不要** ⇒ 落地集空 ⇒ **`pump_test.go:124`／`:276` 那两行本程未改**（没有第五枚键，改它只是把"恰好四枚"换成"随结构体"，买不到任何东西，还白丢一枚钉）。
下一程若拿到 R-1＋R-2，这两行的改法与正控应当是：
名单**由 `Snapshot` 的 `json` 标签现推**（不在第二处抄号），正控＝在**仓外副本**里把某一枚键的标签临时摘掉／把某一枚字段临时删掉 ⇒ `:124` 与 `:276` 必须红（红句具名：`snapshot JSON keys = …`／`exit bytes carry keys …`）；
⚠ 顺带一枚仪器缺口留在案上（r1 §2.7、R6）：两把契约尺读 `json:` 标签，**无标签的导出字段它们看不见**，而 `encoding/json` 照发；今天唯一的兜底就是这两枚键集钉——所以**改成结构体推导时必须连标签齐不齐一起钉**，否则门反而变松。

**R-3 · `AC#2b` 的两枚前提与本程现量不符（照字面停手，不替它改票面）。**
① 票面点 5"复用 `cmdModels` 的取数"指向的是**签名下载清单**，不是 `llm.providers.*.models.*` 配置目录（P6）。
② 票面点 4 要求"逐家量出来再画进快照"——量了：**两家只收三档、第三家一档不收，而这三张词表全是包内私有 `var`、全仓无导出访问器**（P8）。
⇒ ②不是"缺字段"，是**缺一个源**：要快照带"这家接受哪几档"，得先在 `internal/llm/**` 开口子（另立形状），或者面板显示"未知／不适用"（那就不是票面要的那一版）。
**本程既没抄第二份词表、也没改票面一字。**

---

## 5. 前端"应当长这样"逐键清单（**只写清单，零 `frontend/**` 改动**；由 owner 自己带给他用的那枚 agent）

⚠ 口径：本清单是**建议形状**，不是已落地事实——**Go 侧今天一枚都没发**。前端若照它改，那两把双向尺
（`composer_test.go:48`／`approval_test.go:105`）会**反向**报"TS 读了 Go 不发的键"。
所以对齐的正确次序＝**同一枚 commit 同时含 `internal/panel/**` 与 `frontend/src/lib/panel.ts`**（r1 §4.4 的三个读法里那一枚，需编排者／owner 点头）。

| TS 键（挂在 `PanelSnapshot.composer` 下） | 类型 | 含义 | 供它的 Go 字段／真源 | 缺省／未接时的显示纪律 |
|---|---|---|---|---|
| `models` | `ComposerModelOption[]`，元素 `{ provider: string; modelId: string; display: string }` | 当前**可选**模型清单，只含 `enabled=true` | **尚未发**。将来自 `config.Provider.Models` ∩ `ModelSpec.Enabled`（`internal/config/schema.go:348,369-371,400`） | 空数组＝**按钮不画**（今日现状，P9"不得拿假件当真实字段"），**不许**画一条写死的模型名 |
| `currentModel` | `string`（`"provider/model-id"`，或 `""`） | 这一轮真正在用哪枚模型 | **尚未发**。源＝`rt.provs`／`rt.names`（`cmd/wisp/run.go:314-315`，今天只有一枚端点）——⚠ 本程未把它判成"可选清单"的源，二者不是一回事 | `""`＝显示"未选模型"，**不许**回落成清单第一枚 |
| `efforts` | `string[]` | **这一枚模型真被接受的档位词表**（每家各自的一版） | **尚未发**，且**今天不可导出**（§3.2：两张私有 `var`＋一家完全不映射） | 空数组＝**"思考档位"这一维不画**，不是"只有 off"；两者必须在文案上分得开 |
| `effortVerified` | `boolean` | 这一维的"支不支持"是人声明的（`Capabilities.Thinking`）还是探测核实过的（SQLite `provider_health`） | 声明位 `internal/config/schema.go:330`；核实位**有数据、零枚非测试读者**（`internal/memory/dao_providerhealth.go:170,178`） | 缺省＝`false`，界面须能表达"未探测"；**不许**把"未探测"画成"不支持" |

（`AC#2b` 明写：**本格只做"显示"**；任何"面板改档位／换模型"的写回＝新增 `panel.*` 方法＝C17 契约变更，已并入 `Q-64`，默认不做。本程同样零方法名。）

---

## 6. 门禁读数（派单"门禁与终态"那一节，逐把原样）

⚠ **基线用哪一支，先说清**（票面 AC#5 点名）：`cmd/wisp` 一支按派单给的 DLL 注入直跑
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v ./cmd/wisp/`（未走 `scripts/wisp-cli-tests.sh`），
"跑没跑到"只认 `=== RUN` 枚数（r1 P8 同一把尺）。

### 6.1 改前（＝本程进场那棵树，锚 `318d4bb3`／取数时 `971b7dad`，别家提交非本程）

| 包 | rc | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | 红名册 |
|---|---|---|---|---|---|---|
| `./internal/panel/` | 1 | **137** | 81 | **3** | 0 | `TestC21DesignTokensFourWayAgree`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` |
| `./cmd/wisp/` | 0 | **149** | 89 | 0 | 0 | — |

原始件：`.scratch/wisp/probes/145/r2/gate-before-panel.txt`＋`gate-before-cmdwisp.txt`（名册 `roster-before-panel-fail.txt`）。
⚠ 与 r1（锚 `e24ae28`）那两行四数比：`panel` 105→137 RUN、`cmd/wisp` 139→149 RUN——**是票 33／181／185 这三天新加的用例**，不是本程动的。

### 6.2 其余三把仪器

| 仪器 | 读数 |
|---|---|
| `"$(go env GOPATH)/bin/gofumpt" -l internal/panel/ cmd/wisp/` | **空**（`probes/145/r2/gofumpt.txt` 0 行） |
| `sh scripts/d22scan.sh` | **rc=0、clean**；分母现量：bans #1-5 `internal/=211`＋`cmd/=24`（examined 235 production Go files）、ban #6 `frontend/=85`、ban #7 `internal/tools/=21`、ban #8 `design/=39`／`frontend/=85`／`internal/=441`／`cmd/=47`；`skipped as git-ignored: 1 file(s) under frontend/dist/assets/`。全文 `probes/145/r2/d22scan.txt` |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | **本程未跑**，原因照实写：本程**零生产码改动**，而那把尺读的是 `git grep` 在指定锚上的形状；它要防的是"改了触发门却没登记"。⚠ 这是**本件的一处已知欠读**（§7 N6），不是"跑过且绿" |
| `go test -count=1 ./internal/panel/ ./cmd/wisp/`（终态复跑） | 见 §6.3 |

⚠ `tools/d22scan` 是独立 module，本程按 `scripts/d22scan.sh` 跑，**没有**在根目录 `go vet ./tools/d22scan/`；
**没有**跑过 `probes/161/r6/flip-declaration.sh`（派单明令禁止）。

### 6.3 契约轴零命中＋终态名册（只认本程自己的 commit 集）

本程的 commit 集（逐枚 `git show --name-only`，交件时续记）：路径并集**只应有**
`docs/evidence/s1/145-snapshot-growth-r2.md`＋`.scratch/wisp/probes/145/r2/**`＋`.scratch/wisp/issues/145-…md`。
`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／
`frontend/**`／`design/**`／`cmd/wisp/run.go`／`bridge.go`／`internal/panel/**.go`／`pump_test.go` —— **一枚都不该出现在里面**。
终态 `git status --porcelain` 与起手名册（§0.2）的逐枚具名差集＝**本程那三枚路径**，其余别家一枚不动。

---

## 7. 本程**没**测什么（不假装核过）

| # | 没测的事 | 为什么 | 影响读法 |
|---|---|---|---|
| N1 | **没落任何字段**，所以派单"加字段会让 `TestComposerContractTypesMatchFrontend` 更红"那一发**未发生** | §2 的判定：落地集＝空 | 现在那枚红仍只点名 `git` 一枚键；下一程落第一枚键时它会点名新一集 |
| N2 | **没做仓外闭合实验**（r1 的 E1–E4 那一套） | 本程的卡点是"写面里没有装配根、没有测试之家"，这两条**读代码即可具名**（§0.3 构造点表＋`A388` 写面名册），且 r1 已经实测过同一形状（E1b 那发：不动装配根 ⇒ 新段全 0） | 本件的"没人喂"是**引用 r1 实跑＋本程现量名册**，不是本程新跑的一发 |
| N3 | **没跑 `gate-clauses.sh`**（§6.2） | 零码改动 | BAD 腿名册今天该是什么，本件**没有读数**，请下一程或编排者自取 |
| N4 | **没测真机 `wisp run`**（不起真进程） | 本程无行为改动 | 快照今天仍在"泵已装配、最后一公里未接"那一态（`panel_pump.go:153-181` 自证） |
| N5 | **像素／视觉／截图 0 次**；十四态能不能真画，本程一字未判 | 无界面动作 | 见 §8 |
| N6 | `provider_health` 那两枚 SELECT 的**内容形状**未读（只证到零枚非测试调用者）；`ModelSpec` 里除 `Display`／`Enabled`／`Capabilities` 之外的列未逐枚点名 | 超出本程三格 | §3／§5 那两张表的"该带哪几列"是**建议**，不是清单全量 |
| N7 | 台账 `A##`／`Q##` 一条未登记、票面一格未勾（归编排者） | 派单口径 | §4 那三枚停手要落账，请编排者按 `A388` 续号 |

---

## 8. 一句人话（给不读术语的人）

要给面板多送的那两样（"有哪些模型可选"、"这个模型支持哪几档思考"），**东西确实在后台**：
模型目录在配置文件里、每家支持哪几档在适配器代码里。但这一趟**我一样也没加进快照**，原因是两件很具体的事：

1. **加进去的那个"数字接口"不在这一次批准我能碰的文件里。** 后台有个地方专门决定"每样东西从哪个真对象取"，
   它现在只接了六根线（审批队列、权限档、工作区、git、流式文本、出口）。新加一样东西＝**在那个地方多接一根线**，
   而那一个文件这次没划给我。多接线的代价我量过：不接，加上去的字段出门时永远是 `0` 或空——**这正是这张票要防的假数据**。
2. **"哪一枚测试证明它来自真东西"这句话，这一次没有能写下它的地方。** 批准的三个生产文件＋只放开两行的那个测试文件里，
   写不出一枚新断言。⇒ 所以按派单第 1 条："答不出的字段直接不要"。

**还量到一件比"三档／四档"更硬的事**：三家模型接口里，**两家只认三档、第三家一档都不认**（它刻意不发这个参数），
而那几张"认哪几档"的表**全是各家代码里的私有变量，外面读不到**。
所以界面上那个"思考档位"按钮要显示**真的**档位表，得先有人在家那一侧开个口——**这不是加字段能顺手带过去的**。
