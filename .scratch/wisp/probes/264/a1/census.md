# 票 264 — 264-a1 只读普查：把三形的代价量到能裁

**腿**：`264-a1`（只读普查腿，**本轮零 Go 命令**）
**普查时刻**：2026-10-04 09:5x +08，落骨架时锚点 HEAD `447e0d71`（`dev`）
**票面**：`.scratch/wisp/issues/264-error-text-carries-absolute-paths-into-model-visible-text.md`（**整读过**，AC 复选框一枚未动）
**这一格量什么**：三选一（甲＝发出去前脱敏／乙＝行为不动只装仪器／丙＝不做了）的**代价读数**。
**⛔ 本件不裁形**——裁形归机主（人工批准面，`AGENTS.md` §0.2）。全文不出现"建议选甲／乙／丙"。
**与票 174 的分工**：174 管**可达性**，264 管**外发内容的形状**；`174-a4` 已给的名册（`err.Error()` 装配点 52／能进模型可见文本 48）**本件不重跑**，只做它没覆盖的增量（见 §2）。

---

## §1 那枚现成的开关到底管不管这一侧（`[privacy] redact_paths` / `observe.Redactor`）

### 1.1 `Redactor` 被谁构造／谁调用（产码 vs 测试分开数）

`Redactor` 的定义＝`internal/observe/redact.go:45-49`；四个方法＝`Attr`（`:107`）、`String`（`:133`）、`Message`（`:145`，＝`String` 的别名）、`ConfigSnapshot`（`:221`）。

**产码构造点＝3 枚，全部在 `internal/observe` 包内，包外零枚：**

| # | `file:line` | 产码／测试 | 构造出来喂给谁 |
|---|---|---|---|
| P-1 | `internal/observe/logging.go:88` `red: Redactor{RedactPaths: cfg.RedactPaths}` | 产码 | 塞进 `redactHandler{}`（`:87`，类型定义 `:123-126`）→ 挂在 `slog.NewJSONHandler` 外面（`:89-91`）→ `Handler()`（`:104`）／`InstallAsDefault()`（`:107`）→ **只有 slog 记录这一条路**：`Handle`（`:132-139`）对每条记录的 message 走 `red.Message`、对每个 attr 走 `red.Attr`，`WithAttrs`（`:141-147`）连预烤的 attr 也过 |
| P-2 | `internal/observe/diagnostics.go:120` `red := Redactor{RedactPaths: o.RedactPaths}` | 产码 | `:121` `red.String(string(data))` → 把已经落盘的 `wisp-*.jsonl` **再洗一遍**写进诊断包 zip 的 `logs/`（`:122`） |
| P-3 | `internal/observe/diagnostics.go:138` 同上 | 产码 | `:139` `red.ConfigSnapshot(...)` → zip 里的 `config.redacted.toml`（`:144`）；洗不干净就 fail-closed 不出包（`:140-143`） |

方法调用点（包外）：**0 枚**。`grep -rn "Redactor" --include=*.go` 的包外命中只有测试与 `.scratch` 里的探针副本。

**测试构造点＝9 枚字面量＋3 枚经由 config 字段**（分开数，别混进产码）：

| # | `file:line` | 形状 |
|---|---|---|
| T-1…T-8 | `internal/observe/logging_test.go:31, :53, :64, :80, :93, :106, :110, :121` | `Redactor{}` ／ `Redactor{RedactPaths: true}`（只有 `:110` 打开） |
| T-9 | `internal/observe/diagnostics_test.go:121` | `Redactor{}.ConfigSnapshot(cfg)` |
| T-10 | `internal/observe/logging_test.go:23`（`redactHandler{red: …}`，经 `newTestHandler` `:20`） | handler 层形状 |
| 经字段 | `internal/observe/logging_test.go:243`（`LogConfig{… RedactPaths: true}`）、`internal/observe/diagnostics_test.go:39`（`BundleOptions{RedactPaths: true}`） | 这 2 枚＋`:110` 是全仓**仅有**的"路径掩码开着"的用例 |
| 配置本体 | `internal/config/boundary_test.go:152` `c.Privacy.RedactPaths = true` | **唯一**一枚直接吃 `PrivacySection.RedactPaths` 的用例（测的是 config 边界，不是掩码效果） |
| 件外 | `.scratch/wisp/probes/153/accept-r1/zzaccept153_ondisk_test.go:20`（`TestAccept153TraceOnDiskThroughRedactor`） | 票 153 的验收探针**副本**，在 `.scratch` 下，不算包内钉 |

### 1.2 它的调用面是不是只在日志那一侧

**是——而且是"只有"。** 三条投递线逐条具名：

| 线 | 装配点 | 判 |
|---|---|---|
| slog → JSONL 滚动日志 | `cmd/wisp/logsink.go:149` `observe.InitLog(observe.LogConfig{Dir: dir, Level: logSinkLevel})` → `:156` `teeHandler{primary: p.Handler(), mirror: …stderr}`；`cmd/wisp/slo_windows.go:272-275` `InitLog(LogConfig{Dir:…, Level:"info"})` → `:282 pipeline.InstallAsDefault()` | 日志侧（落盘那一条被 `redactHandler` 包着；`logsink.go:154-155` 的注释自己写着"the file side keeps the non-disableable redaction"） |
| 诊断包 zip | `internal/observe/diagnostics.go:62 BuildDiagnosticsBundle` | **产码调用者 0 枚**：`grep -rn "BuildDiagnosticsBundle" --include=*.go`（排 `.scratch`）＝定义 `:35/:61/:62` ＋ `diagnostics_test.go:35/:125/:140`。⇒ P-2／P-3 那两枚今天**在生产里到不了**（归口：这条线属票 45「diagnostics guards」） |
| 模型可见正文 | 见 §1.3 | **一处都没有** |

顺带点名一枚同名不同物的东西，免得后人当它是开关：`internal/llm/content.go:106 RedactContent` 的名字里有 Redact，注释逐字写着"for LOGGING ONLY"（`:102`），而它**产码调用者也是 0 枚**（只有 `internal/llm/llm_test.go:39` 与两处注释 `:15/:69` 提到它）。

### 1.3 模型可见文本那条链有没有任何一处过这个 Redactor（负向句）

负向句，按规矩先读满三处再落笔。

| 谁产出（工具侧 → `Result.Text`） | 谁投递（进 `llm.RoleTool` → HTTP 正文） | 谁落盘 |
|---|---|---|
| `internal/tools/fs.go:137/144/148/158/201/208/212/218`、`fs_edit.go:108/119/128/132/138`、`fs_write.go:247/260/292/325/343/382/389/408/453/460/464/477`、`bridge.go:289/431/561`、`cancel.go:76`、`task.go:509/693/837`、`task_backfill.go:136`、`subagent_197.go:244/322` 各自 `return Result{Text: …}`（`tool.go` 的 C1 Result） | `bridge.go:576-581 out := agent.ToolOutcome{Text: res.Text（`:577`）…}` ⇒ `internal/agent/loop.go:698 log.Text = out.Text` ⇒ `:710 l.sp.Prepare(c.ID, log.Text)`（spill 可能改写这句并把 `全文见 <artifact 路径>` 拼进去，`spill.go:185-189`）⇒ `:722 l.append(toolResultMessage(c.ID, log.Text, …))` ⇒ `:1080-1083 toolResultMessage` 造 `llm.Message{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{… llm.TextPart{Text: text}}}}` ⇒ 三枚适配器各自把它序列化进请求体：`internal/llm/openaichat/request.go:273 encodeToolResult`、`internal/llm/anthropic/request.go:370 encodeToolResult`、`internal/llm/openairesponses/request.go:218-226` | 两个落盘点，**都不经过 Redactor**：① `tool_call` 行（`internal/memory/dao_toolcall.go:17 InsertToolCall`／`:69 FinishToolCall(outcome, errorClass)`）落的是**结论与等级，没有正文**；② 只有超长输出被 spill 写成 artifacts 文件（`internal/agent/spill.go:153-174`）。`history` 正文不落盘——`internal/agent/compress.go:59` 逐字写着"[privacy] keep_transcript forbids putting history text on this…"，而 `keep_transcript` 是硬编码 false（`internal/config/schema.go:517`／`PLAN.md:2741`）。⇒ 这条链上的失败句**只在内存里进正文、不落盘** |

**三处读完的结论（负向句本体）**：上面 15 个产出的 `file:line`、5 跳投递、2 个落盘点里，**没有一处出现 `Redactor`、`redactHandler`、`pathRe` 或任何路径掩码**。
独立复跑票面现量第 3 条（⚠ 我跑的是我自己的尺，不是抄它）：`grep -rni "redact" --include=*.go internal/agent internal/tools | grep -v _test` ＝ **0 命中**；`internal/llm` 侧 7 命中全部点名过：`content.go:15/:69/:102/:106/:116`（`RedactContent`，零产码调用者）、`errors.go:168`（注释"log-safe: redaction still applies upstream"——这句注释说的"upstream"指的是 **observe 那条日志线**，不是外发线）、`anthropic/stream.go:191`（`redacted_thinking` 是 Anthropic 协议字段，与路径无关）。

### 1.4 配置键名、默认值、有没有产码填过它

| 项 | 读数 |
|---|---|
| 键与 tag | `internal/config/schema.go:509` 逐字：`RedactPaths bool \`toml:"redact_paths" default:"false"\``，挂在 `PrivacySection`（`[privacy]`，定义 `:504-505` 的注释"PrivacySection is [privacy]"）⇒ **默认值就是 `false`，且是字面写死的 `default:"false"`**（不是 `default:""` 那种空串形状） |
| 注释自己的口径 | `schema.go:508`："masks filesystem paths in **user-visible text**"；`redact.go:46-47` 同词。⚠ **这句话与实际射程不符**：实际只作用在 slog 记录与诊断包文本上（§1.2/§1.3），模型可见正文不是"user-visible"，但它才是今天真出去的那一侧 |
| 两处 `LogConfig` 有没有填 | **都没填**（逐枚点名）：`cmd/wisp/logsink.go:149` 只有 `Dir`＋`Level`；`cmd/wisp/slo_windows.go:272-275` 只有 `Dir`＋`Level`。⇒ `RedactPaths` 走零值 `false` ⇒ **今天连日志那一侧的路径掩码都是关着的**（其余五类掩码不可关，`redact.go:13-20`／`logging.go:86`） |
| 有没有产码从 `cfg.Privacy` 读它 | **没有**。仓里已经把这句话说在产码里了：`cmd/wisp/config_readers_255.go:134` 逐字 `"privacy": hotClaimNoReader + "internal/observe/logging.go:50 [RedactPaths] - the mirror field exists and is filled by callers, never from cfg.Privacy"`——这一行是**产码里的一张名册**，且**每轮跑都被重新量一遍**：`cmd/wisp/config_receipt_255_test.go:166`（`TestTicket255RosterStillMatchesTheActualReadSites`）＋`:179`（`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`：file 必须存在、行号必须在range内、那一行必须还带着被引的 token，且"每个段的生产读文件集合必须正好等于名册点名的集合"） |
| `unwired.go` 里有没有它 | **没有**。`internal/config/unwired.go:61-105` 那枚 `unwiredKeys` 只覆盖锁定段（文件头 `:11-19` 三条准入条件之一＝"sits in a locked section (…`[risk]`, `[fs]`, `[net]`, `[plugins]`)"），`[privacy]` 不在锁定段 ⇒ 想靠这张表把 `redact_paths` 说成"未接线"是**接不进去的**（§6.3 的落点因此不是这里） |
| `BundleOptions.RedactPaths` | `diagnostics.go:44` 字段在，但整条 bundle 线零产码调用者（§1.2 第二行）⇒ 这枚默认值今天同样无人填、无人到 |

### 1.5 一句话结论

> **"仓里已有脱敏开关"这句只管日志侧，一点都管不到外发正文侧**——而且今天**日志侧那一半也是关着的**（两处生产 `InitLog` 都没填 `RedactPaths`，schema 默认 `false`，产码里没有任何一处从 `cfg.Privacy` 读它）。
> 更贴身的补刀：**就算把它打开，也遮不住这一侧**——`redact_paths` 的开关只存在于 `observe.Redactor` 里，而 §1.3 那条链（`bridge.go:577 → loop.go:698/710/722 → :1080 → llm/*/request.go`）一枚 `Redactor` 都不碰。
> 所以票面"不做"那一栏写的"日志侧与模型侧都无掩码"**读数成立**（我复跑确认），并且还要再加一句：**`redact_paths` 即便被接上并打开，模型侧照旧无掩码**。

附一条本腿新量到的、`174-a4` 明确留白的东西（它自己列在第 8 条"没逐字读那条正则"）——`pathRe` 逐字（`internal/observe/redact.go:81`）：

```go
pathRe = regexp.MustCompile(`(?i)(?:\\{1,2}\?\\)?(?:[a-z]:\\|\\\\)[^\s"',;:)>\]]+|(?:(?:/home|/users|/root)/[^\s"',;:)>\]]+)`)
```

我的**读法**（⛔ 不是执行级判断，见 J-1）：两支都要求**反斜杠**接在盘符后（`[a-z]:\\`）或 UNC 的 `\\`，或 POSIX 侧的 `/home`｜`/users`｜`/root` 开头；
按字面读，**斜杠正写的 Windows 绝对路径（`C:/Users/…`、`D:/work/…`）落不进这两支**，而 `rules_gateway.go:42/:73` 的 `%q` 与 `risk.Resolve` 的 `lexCanonical` 产出的正是反斜杠形（`fs.go:144` 那句吃的是 `Canonicalize` 的错，见 §2.3）。
⇒ "打开后遮得住哪些真串"＝〔待验，需跑一遍才知道〕，本腿只把正则字面与读法摆这儿。

---

## §2 除了 `fs` 族，还有哪些族的失败句今天真把路径带出去（增量，⛔ 不重跑 174-a4）

`174-a4` 件＝`.scratch/wisp/probes/174/a4/census.md`（187 行／26,592 字节，本腿 `wc -l -c` 现量），其分母（`err.Error()` 装配点 52／能进模型可见文本 48／`%w` 272／插值动词 136）**按令未复算**，属〔仅自述〕。本节只补两块。

### 2.1 `internal/tools` 里非 `fs`／`fs_write` 的族

尺：`grep -rn 'err\.Error()' internal/tools/*.go | grep -v _test` ＝ **42 枚**（本腿现量，含 `fs.go`／`fs_write.go`）；`fs.go`＋`fs_write.go` 之外的 **13 枚**逐枚判如下。⚠ 另外三族（`fs_edit`／`fs_staging`／`bridge` 的门禁句）不在这 13 枚里也要单独看，见 2.3／2.4——**票面 §现量 只点了 4 行（`fs.go:144`、`:148`、`fs_write.go:343`、`:345`），量到的远不止 4 行**。

| `file:line` | 句子（逐字截） | 该 err 来自带路径的调用吗 | 落进 `Result.Text`（＝模型可见）吗 |
|---|---|---|---|
| `internal/tools/fs_edit.go:108` | `参数解析失败：`＋err | 否（`json.Unmarshal`，带字段名不带路径） | 是 || `internal/tools/fs_edit.go:119` | `路径无法解析（按 fail-closed 拒绝）：`＋err | **否**——见 2.3：这一支的 err 集合里没有带路径的那枚 | 是 |
| `internal/tools/fs_edit.go:128` | `refusal("无法确认目标是否存在…："＋err)` | 需 `stat` 失败才走到（`os.Stat` 的 `*fs.PathError` 带路径） | 是（`refusal()` 见 `:92-95`，整句拼进 `Text`） |
| **`internal/tools/fs_edit.go:132`** | `目标文件不存在，没有可定位的内容：`**＋`target`** | **是，而且根本没有 err**——这句把 canonical 后的目标路径**自己写进正文** | 是 |
| **`internal/tools/fs_edit.go:138`** | `读取目标失败：`＋err＋`（目标未被改动）`，并带 `Origin: target` | **是**（`os.ReadFile(canon)` 的 `*fs.PathError`＝逐字绝对路径） | 是 |
| `internal/tools/fs_edit.go:142/160/163/168/172/178/182/193/209/254/263` | 其余 `refusal(...)` 十余枚 | 否（都是字节数／编辑定位类，路径不在串里） | 是 |
| `internal/tools/bridge.go:289` | `参数不是合法的 JSON 对象：`＋err | 否（同上） | 是（`reject()` `:1019-1024` 把 `why` 原样当 `Text`） |
| **`internal/tools/bridge.go:431`** | `该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：`＋`dec.Reason` | **是**——`dec.Reason` ＝ `verdict.Reason`（`:301`），而 R2/R3 的理由句**逐字带路径字面**：`internal/risk/rules_gateway.go:42` `R2: 路径无法规范化，按越界处理（fail-closed）: %q`、`:73` `R3: 路径无法规范化，敏感检查 fail-closed: %q`（用例把形状钉死了：`internal/risk/assessor_test.go:172` 的 `want` 里就有 `… fail-closed）: "C:/broken"`） | 是（同 `reject()`） |
| **`internal/tools/bridge.go:561`** | `工具内部故障：`＋err | **放大器**：任何工具 `Execute` 直接 `return Result{}, err` 都从这里出去，而 err 常是 `os.*` 的 `*fs.PathError`（例：`fs.go:132` ctx 之外还有各条 `Result{}, err` 支） | 是 |
| `internal/tools/cancel.go:76` | `上下文已结束：`＋err | 否（`context.Canceled`／`DeadlineExceeded` 定形句） | 是 |
| `internal/tools/task.go:509`、`:693` | `参数解析失败：`＋err | 否 | 是 |
| **`internal/tools/task.go:837`** | `注意：这条路径现在读不到，C26 连规范化都没通过（`＋err＋`）` | 这支的 err 就是 §2.3 那枚**不带路径**的集合；路径本身由句子里的 `canonicalizeFailsHere`/pointer 段带出（见 2.5 的用例） | 是，但**今天有没有真接线是分裂的**：`task_output_canonicalize_fail_174_test.go:130-148` 自己就分 wired／unwired 两支各断一句 ⇒ 生产装配走哪一支＝运行时问题（J-10） |
| **`internal/tools/task_backfill.go:136`** | `超长输出没能落盘，名册里不许诺任何路径：`＋err | **是**：这句的 err 来自 spill／artifacts 写盘（`internal/agent/spill.go:153-174` 的 `observe.Wrap(ClassResource, err, "agent: write spill artifact")`，链里是 `*fs.PathError`） ⇒ **一句话自己说不许诺路径、后半截却把路径接上去** | 是 |
| `internal/tools/subagent_197.go:244` | `参数解析失败：`＋err | 否 | 是 |
| `internal/tools/subagent_197.go:322` | `子代理环路装配失败：`＋err | 装配类（`observe.New(ClassConfig, "agent: no LlmProvider (C5 seam) wired")` 这类定形句为主）；不能排除链里混进 `os.*` | 是 |

`internal/tools/fs_staging.go`：**这一族今天一枚 `err.Error()` 都没有**（`grep -n 'err.Error()\|Errorf\|Text:' internal/tools/fs_staging.go`＝空）。它对外说话走 `fs_write.go` 的句子。

**路径还有第二条腿（⛔ 不是"失败句"，但同一条投递链）**：`AppliedSteps` 台账里的 `se.record` 语句带目录与临时文件名——`fs_write.go:295`（`在 %s 创建临时文件 %s`，`parent` 是绝对目录）、`:347`（`原子重命名 %s → %s`）、`:412`（`通过 %s 将 %s 放入回收站，并已核对还原记录 %s（位于 %s）`）。这些串经 `bridge.go:581 AppliedSteps: …res.AppliedSteps…` → `internal/agent/approval/report.go:98/:145`（逐条拼进 `txt`）→ `bridge.go:594-595 out.Text = orDefault(out.Text, "调用已中止") + "\n\n" + txt`。⇒ **调用被中止／否决那一条路上，台账把目录名一并发给服务商**，而且这条支的开关是"有没有 AppliedSteps"，与失败句无关。（`fs_write.go:115 stoppedResult` 也带着它。）

### 2.2 C25 那条：`internal/tools/bridge.go:646` 到底把什么喂给了谁

复认结果（三处读满）：

| 谁产出 | 喂进去的是什么 | 谁消费／落到哪 |
|---|---|---|
| `bridge.go:638 func (b *Bridge) mark(dec Decision, res Result, hostPath string)`；`:641 origin := res.Origin`（`tool.go:50-53`：`Origin` 是"为 C25 溯源标记指明这条结果来自哪儿"的字段，`fs.go:148/158/170/212/218/238`、`fs_edit.go:138/272`、`fs_write.go:263` 都会填 canonical 后的路径）；`:642-643` 若 `origin` 为空则退到 `dec.Paths[0]`（＝`displayPaths` 的产物，见下） | **`:646 b.prov.MarkWithHostPath(dec.TaskID, dec.Tool, origin, res.Text, hostPath)`——第 4 个实参逐字是 `res.Text`**，也就是 §2.1 那些**含绝对路径的失败句整串**；第 3 个实参 `origin` 单独又是一枚路径 | `internal/risk/provenance.go:558 MarkWithHostPath`（`Mark` `:479` 是它的 `hostPath=""` 特例）→ 建 taint 片段索引（`taintmatch.go:91-144`）。**落盘＝内存态 scope 账**，不是外发正文；它的用途是**出网前 R4 粗粒度匹配**（`PLAN.md:1375` C25 契约逐字："工具结果的来源标记（`{tool, origin, sensitive}`）＋ 出网前的**粗粒度污染匹配**（规范化后 ≥8 字符连续片段命中即升 L2）"）。`bridge.go:648-649` 匹配不上时只 `b.log(...)` 一句"produced no matchable taint fragment (origin=%q)" |

⇒ 这一条的**性质要说准**：`bridge.go:646` 不是"第二条外发腿"，它是**把含路径的回执喂给污染台账**（溯源／升 L2 用）。⚠ 但正因为喂的是 `res.Text`，**甲形若落在工具侧（§4 候选①）就会改这个实参**，连带改变 C25 片段的形状——这既是它的代价，也是它会顶到人的一处（见 §5.6）。

### 2.3 `displayPaths`／卡片那一侧（顺带量到，具名不裁）

`bridge.go:942 out = append(out, p+" (无法规范化: "+err.Error()+")")`（在 `displayPaths` `:934-949` 里）：这一枚**同时带用户的原始拼写 `p` 和 err**。它的去向不是模型正文：`:296 dec.Paths = b.displayPaths(rawPaths)` → 审批卡片与 `tool_call` 审计行（`grant.go:76` 注释、`bridge.go:1092`）。用例把这处形状钉成了**故意的**：`internal/tools/bridge_junction_windows_test.go:447`（注释逐字："<原样路径> (无法规范化: …)，批准它的人看不到自己批准的是什么"）＋`:459-460`（断言卡面**必须**带着"无法规范化"标记与那枚原样路径）。另有名册级断言：`internal/session/grants_test.go:267`、`:320`。⇒ **这一侧带路径是有明文理由的（人要看得见自己批准了哪儿），与 §1.5 那侧不是同一件事**；本腿只点名，不裁。

### 2.4 ★ 对票面 §现量 第 1 条的一处更正（读码级，具名）

票面说 `fs.go:144` 那句的 err"消息形在 `internal/risk/pathresolver.go:95`（`%w: %s expands to %s (%s)`）"。本腿顺着 `t.d.open` 读满：`internal/tools/fs.go:93-104 open()` 只做三件事——`d.Paths == nil` ⇒ 定形句 `fs: no C26 resolver configured (fail-closed)`；`d.Paths.Canonicalize(raw)`；结果含 `??` ⇒ `fs: unresolved path component`。而 `internal/tools/paths.go:119-148 Canonicalize` 的错误来源只有两枚：`tools: empty path`（`:121`）与 `p.resolve(raw)` 的 err；`p.resolve`（`:103-109`）透传 `risk.Resolve`（`pathresolver.go:105-139`），后者的**唯一 `return res, err`** 是 `:123 return res, ErrReparseDenied`（定形句，`:36`，不含路径）。
⇒ **按读码：今天 `fs.go:144` 与 `fs_edit.go:119` 这两枚"路径无法解析"句里，量不到 `expands to` 那一形**；`Actable()`（`pathresolver.go:93-99`，那句带 `%s expands to %s` 的）在生产里的调用者是 `internal/panel/workspace.go:85` 与 `internal/risk/syncdirs.go:208/:226`，**`internal/tools` 零枚**（`grep -rn "Actable()"` 现量）。
⛔ 这不是"那一侧没问题"：`fs.go:148/158/212/218`、`fs_edit.go:138`、`fs_write.go:292/325/343/408/477` 吃的是 `os.Open`／`os.ReadFile`／`os.Rename` 的 `*fs.PathError`＝**逐字绝对路径**，票面点名的 `:148` 与 `:343` 两支**复认成立**。⇒ 归口：票面 §现量 第 1 条里"`:144` 的消息形在 `pathresolver.go:95`"这一句需要更正（已推历史不改写，追加一条 `A##` 即可）。

### 2.5 增量小结（一句）

> 除 `fs`／`fs_write` 外，今天真把路径带进模型可见正文的还有**四族**：**`fs_edit`**（`:128`/`:132`/`:138`，其中 `:132` **不经任何 err** 直接写 `target`）、**门禁理由句**（`bridge.go:431` ← `rules_gateway.go:42/:73` 的 `%q`）、**spill／落盘失败句**（`task_backfill.go:136`，句子自己否认带路径）、**AppliedSteps 台账**（`fs_write.go:295/347/412` → `report.go:145` → `bridge.go:594-595`）。另：**票面点名的 4 行只是冰山的一角**，同形状行号在 `fs`／`fs_write` 两文件里本腿就数到 **21 枚**（`fs.go` 8＋`fs_write.go` 13，见 §1.3 产出列与 2.1）。
---

## §3 会不会撞已定案契约（**读了 `docs/PLAN.md`／`docs/specs/**`，一枚字未改**）

读的范围（具名，全为只读）：`docs/PLAN.md:1057-1112`（D30）、`:1375`（C25 契约行）、`:2514-2545`（D34 权威表）、`:2662`（D46 提取器行）、`:2741`（D36 `[privacy]` 行）、`:2764-2812`（D37 错误模型）；`docs/specs/SPEC-06-security-gatekeeping.md:60-92`、`SPEC-03:37`、`SPEC-04:77-86`、`SPEC-05:106`；`SPEC-10-testing-acceptance.md`（`grep -n "脱敏\|失败句\|正文\|绝对路径"`＝**0 命中**）。

### 3.1 有没有任何已定案文字要求"失败句必须原样给模型"或反之"必须脱敏"

**先给结论的形状**：⛔ **没有一句要求"原样"，也⛔ 没有一句要求"外发正文必须脱敏路径"**。但有 **5 处已定案文字与这一格相邻到不能不点名**，逐字引：

| # | 出处 `file:line` | 原文（逐字） | 与本格的相干性（只报形状，不裁） |
|---|---|---|---|
| K-1 | `docs/PLAN.md:2788`（D37 表 `tool` 行） | `\| \`tool\` \| 工具内部失败、参数不合法 \| **回给 LLM 让它自纠，不打断用户** \| \`Acting\` \| 由 Agent 决定 \|` | 它规定的是**目的**（让模型自纠），不是**内容形状**。甲把路径换成 `<path>`／相对写法，算不算还在"让它自纠"的射程内＝一个解释问题（J-4） |
| K-2 | `docs/PLAN.md:2789`（D37 表 `permission_denied` 行） | `\| \`permission_denied\` \| 越界目录、未声明 capability、黑名单、reparse point \| 回给 LLM + **面板可见** \| \`Acting\` \| 否 \|` | §2.1 的 `bridge.go:431`＋`rules_gateway.go:42/:73` 走的**正是这一行**：它同时要求"回给 LLM"与"面板可见"，而面板那一侧已被用例钉成**要看得见原样路径**（`bridge_junction_windows_test.go:447/:459-460`）⇒ 甲落点若选在"两个消费者共用的一串"，这两半会开始互相拽 |
| K-3 | `docs/PLAN.md:1375`（C25 契约行） | `\| **C25** \| **\`Provenance\`** \| 工具结果的来源标记（\`{tool, origin, sensitive}\`）+ 出网前的**粗粒度污染匹配**（规范化后 ≥8 字符连续片段命中即升 L2）。D30① 组合闸门的实现契约 \| D33/F4 \|` | **本格唯一一处带数字的已定案射程**（"≥8 字符连续片段"）。甲落**工具侧**（§4 候选①）改的正是 `MarkWithHostPath` 的第 4 实参 `res.Text`（`bridge.go:646`）⇒ 片段集合直接跟着变。C25 是 **C 编号** ⇒ 触碰它＝人工批准面（`AGENTS.md` §0.2） |
| K-4 | `docs/PLAN.md:1101`（D30 第④层） | `工具输出进上下文时用明确边界包裹 + system prompt 声明「这是数据不是指令」` | 它要求的是**包裹**，不是**替换**。紧接的同一段还写着"⚠ 业界共识：这是**弱防御**……不得当作主力" ⇒ 想引这一条给甲背书，得连这句一起引 |
| K-5 | `docs/PLAN.md:2741`（D36 `[privacy]` 行）＋ `docs/specs/SPEC-03-config-secrets-envs.md:37` | `PLAN.md:2741`：`\| \`[privacy]\` \| \`redact_paths\` \`diagnostics_opt_in\` \`retention_days\` · **\`keep_transcript\`/\`keep_audio\` 为硬编码 false…** \| \`hot\`（硬编码项不可改） \|`；`SPEC-03:37`：`\`[privacy]\` \| \`redact_paths(bool)=false\` …` | **默认值 false 是被冻结契约写死的**。⇒ "读它"不改这两处一字；"**把默认改成 true**"改的正是这两处 ⇒ 明确是契约变更。（⚠ 按既有记忆坑：`PLAN.md:3`／`SPEC-08:85` 那类计数已腐坏但冻结——这里引的是**行内文字**不是计数） |
| 旁 | `docs/PLAN.md:2662`（D46 Tier-1 行） | `原始 stdout 不直接进 LLM 上下文（截断上限 ≤32KB）` | 形似而射程不同：管的是**插件命令 stdout**，不是内置 fs 失败句 |
| 旁 | `docs/specs/SPEC-04-voice-pipeline.md:86` | `模式匹配脱敏（密码/卡号/token 样式掩码）：尽力而为，UI 必须标注「非完全脱敏」` | **语音转写那一侧的脱敏**，与外发正文不是一条线；⛔ 别当"已经有一处要求脱敏"的证据引 |
| 旁 | `docs/specs/SPEC-06-security-gatekeeping.md:91-92` | `【SPEC】URL 长度上限 2048 字符（防把文件内容塞进 query 外发）` ＋ `数据/指令边界标记：工具输出进上下文用明确边界包裹…` | SPEC-06 全篇**没有一句**规定失败句／外发正文的路径形状（`grep "外发\|发给模型\|进上下文\|上下文可见"` 于 SPEC-06＋SPEC-05＝上述 2 命中＋SPEC-05:106 一条无关的排序说明） |

**产码里两句"已定案的实现自述"**（不是 D/C 编号，但会拦住想顺手改的人，具名）：`internal/tools/bridge.go:75` 注释 `card renders RulesHit/Reason VERBATIM, so they must leave the bridge as …`；`internal/tools/task_output_leg_test.go:425-426` 断言逐字 `the refusal echoed a canonicalized path`——**这一枚已经在往"refusal 不许回声 canonical 路径"的方向钉了**（方向与本格的甲一致，射程只覆盖 `task.query` 那一条）。

### 3.2 若走甲，算不算契约变更（只报撞／不撞）

**不撞**（读了原文，找不到任何一句要求路径必须原样出去）：`D34` 权威表的表体（`PLAN.md:2529-2545` 管的是工具集／RiskLevel／capability／落切片，没有一行管失败句内容）、`D30` 五层（`:1085-1103`）、`C26`（SPEC-06 §4）、`D37` 的 `error_class` 枚举本身、`D36` 的三档生效级别。

**撞／相邻到必须点名**（甲三形的落点不同，撞面就不同——这是"落点选择"问题，不是"甲本身"问题）：

| 甲落点 | 撞哪一条 | 性质 |
|---|---|---|
| 候选① 工具侧（改 `res.Text`） | **K-3 `PLAN.md:1375` C25 的"≥8 字符连续片段"** ＋ K-2 | C 编号 ⇒ 人工批准面 |
| 候选② 投递侧（`loop.go:710→:722` 之间统一过一道） | **不撞任何 D/C 字面**（在 mark 之后、在 C5 之前）；但会把票 174 已交付的"可读指针"一起遮掉 ⇒ 撞的是**已验收用例**，见 §5.2 那 6 枚 | 射程外／批准面内，待裁 |
| 候选③ 接缝侧（三枚 `llm/*/request.go`） | **不撞任何 D/C 字面**；但 C5（`AGENTS.md` §1.3 列的注入面＝golden SSE）是冻结契约，改适配器＝**golden 测不到**（现量：`find` 到的 6 个 testdata 目录里**零枚文件含绝对路径**）⇒ 这一形买到"最少副作用"，也买到"这一侧从此没有钉" | 待裁 |
| 任何一形，只要"把默认值改成 true"或"改 `default:` 措辞" | **K-5 两处字面**（`PLAN.md:2741`／`SPEC-03:37`），且票面禁区第 3 条已另写"⛔ 不许改 `internal/config` 里任何 `default:` tag 的措辞来『显得更安全』" | **明确是**契约变更 |
| 任何一形，若用 `filepath.Rel/Clean/Abs` 做"绝对→相对" | `AGENTS.md` §1.2 禁止清单第 2 条（`PLAN.md:1288`，CI 自动扫：尺在 `tools/d22scan/main.go:15` ban #2 `pathresolver-bypass`，扫 `internal/`＋`cmd/`）；而**解药 `tools/d22scan/allowlist.txt` 恰在票面 AC#3 的越界黑名单里** ⇒ 加白名单这条路对本票是堵死的 | 撞禁止清单，非契约 |

> 报一句就收口（⛔ 不裁）：**甲本身不要求改任何 D/C 的字面**——真正会撞 C 编号的是"**落点①**"，真正会撞两处冻结字面的是"**改默认**"；这两件事都可以不做，甲也可以做。

---

## §4 甲形的落点候选与各自代价（只报形状，不裁）

| 候选 | 落点 `file:line` | 改动面枚数（本腿数到的） | 牵不牵动别的包 | 会不会让模型看不见是哪个文件而多试一轮 |
|---|---|---|---|---|
| **① 工具侧**：产出 `Result.Text` 时就不带路径 | 必改＝带路径那 4 族 5 枚文件：`internal/tools/fs.go`（`:148/:158/:212/:218`）、`fs_write.go`（`:292/:325/:343/:408/:477` ＋ 台账 `:295/:347/:412`）、`fs_edit.go`（`:128/:132/:138`）、`bridge.go`（`:431` 经 `dec.Reason`、`:561` 万能catch）、`task_backfill.go`（`:136`）＝**5 枚文件／约 16 枚行**；`fs.go:144`、`fs_edit.go:119`、`fs_write.go:260/389/460/464` 那几枚按 §2.4 **今天不带路径**，不必动 | **会，且是三形里最重**：`res.Text` 同时是 `bridge.go:646` 喂 C25 的实参 ⇒ 污染片段跟着变（K-3）。现成的低代价替代形状已经存在：`Result.Origin`（`tool.go:50-53`）本就是"路径走另一条通道"的字段，`fs.go:148` 已经在填它 ⇒ "**Text 不带、Origin 带**"是包内既有形状，不是新造 | **会**（`fs.go:148` 整句只有"打开失败："＋路径；掩掉后模型只剩动词）。⚠ 仓里**没有现成的"绝对→相对"能力**：`Canonicalize`（`tools/paths.go:119-148`）只会返回绝对 canonical，包内现成的只有 `baseOf()`（`fs_write.go:340/:347` 在用）取 basename。⇒ "换成相对写法"这一栏的字面承诺，落点①要新造东西；换成 basename 是更便宜的一形 |
| **② 投递侧**：`loop.go:722` 那一跳之前统一过一道 | **1 枚文件**：`internal/agent/loop.go`（夹在 `:710` spill 之后、`:722 append` 之前；或改 `:1080-1083 toolResultMessage` 一处）。helper 可新建在同包 | **会**：这一跳的串同时被三个消费者拿走——`:720 res.ToolLog`（DB 侧无读者，见 §1.3 落盘列）、`:722 toolResultMessage`（正文）、`:723-725 l.publish(Event{Text: log.Text})`（面板）。⇒ **一把掩码会把票 174 交付的"全文见 `<artifact 绝对路径>`"（`spill.go:186`）一起遮掉**，而那正是票 174 的命门（见 §5.2 的 6 枚正钉）。票面排程已写明这一形"写面跨 `internal/agent`＋`internal/tools`，要**具名解冻**并先跑撞钉预检" | **不会全瞎**：`log.Text` 里除路径还有工具自己写的其它内容；但 `fs.go:148` 那种"只有动词"的句子，落点②与落点①一样瞎 |
| **③ 接缝侧**：`llm` 适配层发请求前统一替换 | **3 枚文件**：`internal/llm/openaichat/request.go:273-…`、`internal/llm/anthropic/request.go:370-…`、`internal/llm/openairesponses/request.go:218-226`（三处 `encodeToolResult`）。可选第 4 处：`internal/llm/provider.go:126`（`case ToolResultPart`） | **最轻**：在 C25 mark 之后、也在面板与 DB 之后 ⇒ 用户看得见的一切不变（面板／卡片／审批／`tool_call` 行全不受影响，与票面"乙＝用户看得见的一切不变"那句同侧）。代价在别处：三形各写一份 ⇒ **三份实现会漂移**；且 `tools/mockllm` 不吃真适配器 ⇒ 测试里看不见这道掩码（C5 那条注入面看不见它，golden 今天也零枚含路径） | 同②；⛔ 但这一形"看不见是哪个文件"这件事**没有任何现成仪器会报**——它发生在最外面那一跳，仓内的尺一律够不着 |

### 4.1 现成用例里"失败句要有路径／要有 canon"的断言（甲形会顶到的）

见 §5.2 的具名清单（**正文级 6 枚**是"甲一落地就红"的那批；另外卡片级 2 枚与 C25 stub 级 3 族是"落点①才碰得到"的那批）。一句话：**这批钉全部长在票 174 那一侧的"可读指针"上，没有一枚是为 fs 族的失败句写的**——失败句那一侧的钉是 §5.3 的**负向**钉（不许回声路径），枚数比正向少，而且只覆盖 `task.query`／spill notice 两处。

---

## §5 既有钉名册（穷举，别抽样）

### 5.1 尺与命中集

- 尺（原话）：`grep -rln 'ToolResult\|IsError\|Result{Text' --include=*_test.go cmd internal`
- 命中枚数：待填

### 5.2 ① 断言"失败句里含某路径／含 `expands to`"的用例（甲形一落地会红的）

| # | `file:line` | 断言原文 | 甲形落地后 |
|---|---|---|---|
| | | | |

### 5.3 ② 断言"失败句不含某物"的用例

| # | `file:line` | 断言原文 |
|---|---|---|
| | | |

### 5.4 ③ `internal/risk/pathresolver.go:95` 那句 `%w: %s expands to %s (%s)` 被谁引用

| # | `file:line` | 性质 |
|---|---|---|
| | | |

### 5.5 具名清单枚数

> 待填。

---

## §6 排程结论（不裁形）

### 6.1 若答甲：最小可落地腿的写面＋必须先解冻／必须改的既有钉枚数

> 待填。

### 6.2 若答乙：那枚能力尺该钉在哪一跳（复认三候选今天哪一处能被断言）

| 候选 | `file:line` | 今天能被断言吗 |
|---|---|---|
| 甲 | `internal/risk/paths.go:122-126` | |
| 乙 | `internal/tools/task.go:837` | |
| 丙 | `fs.go:144+148` ＋ `fs_write.go:343` | |

### 6.3 若答不做：缺的那句"说到明处"该写在哪

> 待填。

---

## 门禁读数

**本节先于六格填表写满并 commit**（本轮硬闸门）。以下读数全部是本腿自量，⛔ 不是转述。

### G-1 零 Go 命令（本腿最重要的一条门禁）

本腿**全程没有执行任何 Go 工具链命令**。逐名列出我实际用过的命令形状：

| 用过 | 形状 |
|---|---|
| 文本检索 | `grep -r`／`grep -rn`／`grep -rln`／`grep -o`（全部经 Grep 工具或 Git Bash） |
| 定点读码 | `sed -n 'A,Bp'`（只读窗口）、`Read` 工具、`ls`／`find` |
| 度量 | `wc -l`／`wc -c` |
| Git | `git status --short`／`git log --oneline`／`git rev-parse`／`git add`＋`git commit -F`（同发、显式 pathspec） |

**没有执行**（一枚都没有，具名点清）：`go build`、`go vet`、`go test`、`go run`、`go mod`／`go get`、
`go build ./tools/d22scan`，也**没有运行**仓里那枚现成的 `tools/d22scan/d22scan.exe`（盘上存在＝我读到了它，
没跑它）。原因＝同机有写码腿 `260-r1`（`internal/ball/**`）与 `262-r1`（`scripts/`＋`.github/workflows/ci.yml`）
正在改动，任何一次编译／测试都会互洗读数。
⇒ **推论（必须说到明处）**：本件所有"某枚用例会不会变红""某个开关打开后遮不遮得住"的**能力类**判断，
**只能是读码级判断**，一律标〔待验，需跑一遍才知道〕，⛔ 不许写成"做不到"，也不许写成"已验证"。

### G-2 本件的自我度量

| 项 | 读数 |
|---|---|
| `wc -c census.md` | 自指量不入正文：以每枚 commit 的 `git show --stat` 与本腿回报里的实测值为准（写入数字会让下一次 `wc -c` 立刻过期） |
| 占位符（`待填`）枚数 | 本枚 commit 时＝填表尚未开始，六格里仍留占位；**终态枚数必须＝0**（终稿会重测并在 G-2 里改成本读数） |
| `（填写中）`枚数 | 0（本腿不用这一形状） |
| 写面枚数 | 1＝`.scratch/wisp/probes/264/a1/census.md`（另：`.scratch/commit-msg-264a1-*.txt` 是 commit message 临时件，按"只建不删"留在 `.scratch/`，不入 commit） |
| 产业代码／测试件／票面／台账改动 | **0**（`git diff` 见每枚 commit 的 stat） |
| `docs/PLAN.md`／`docs/specs/**` | **读了、未改**（§3 的读数：`docs/PLAN.md:1057-1112`、`:1375`、`:2514-2545`、`:2662`、`:2764-2812`；`docs/specs/SPEC-06-security-gatekeeping.md:60-92`、`SPEC-05:106`、`SPEC-04:77-86`、`SPEC-03:37`） |
| `frontend/**`／`design/**` | **一枚没读、一枚没写**（§4 需要的"用户看得见的那一侧"我从 `internal/tools`＋`internal/panel`＋`cmd/wisp` 的产码与用例读到了，没进这两个目录） |
| AC 复选框 | 264 票面 `[ ]` 三枚全部原样未动 |

### G-3 越界自查（对照票面 AC#3 的黑名单）

`git diff --stat` 里**不出现**以下任一路径：`docs/PLAN.md`／`docs/specs/**`／
`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`。
本腿写面只有 G-2 那一枚，故这一条是**结构性满足**，不是努力满足。

---

## 判不动的地方

逐条具名＋归口。凡本腿读码读不出来的，写在这里，⛔ 不用"应该没问题"填空，⛔ 不写成"做不到"。

| # | 判不动的具体问题 | 为什么判不动 | 归口（该由谁／哪一枚腿判） |
|---|---|---|---|
| J-1 | `redact_paths` 打开后，`pathRe` 到底遮住／漏掉**今天真出去的哪几枚串** | 正则字面我逐字抄到了（§1.5），但"某一枚具体失败句会不会命中"是**执行级**判断；本腿零 Go 命令 | 乙／甲落地的写码腿＋AC#2 正控（种一发带绝对路径的失败句⇒必响） |
| J-2 | 甲形会不会让模型**多试一轮**才改对 | 要真跑一轮 provider 循环；mockllm golden 只能测形状测不了"模型会不会因此重试" | 机主（这一枚正是甲形"机主看得见的后果"那一栏的代价本体） |
| J-3 | 这些失败句**今天多久出去一发**（发生率／覆盖面） | 静态读码量不出运行时频率 | 乙形的仪器；⛔ 不许拿"读码推定的高频"当读数用 |
| J-4 | 甲形算不算**改契约** | §3 只报得出"哪一条已定案文字与它相邻／相抵"（`PLAN.md:2788`、`:1375`、`:1101`），"相抵到什么程度算改"是人工批准面 | 机主＋编排者（`AGENTS.md` §0.2） |
| J-5 | 逐枚用例"变红／不变红"的**最终判语** | 我只给得出"这一枚断言吃的是哪一句串"；变红与否要跑 | 写码腿的撞钉预检（形甲落地前必须先跑，票面排程已写"要具名解冻并先跑撞钉预检"） |
| J-6 | `174-a4` 的分母（装配点 52／能进模型可见文本 48／`%w` 272／自带插值动词 136） | **按令不重跑**（避免与它互洗）；本腿只在自己的格子里复认它没覆盖的两块 | 若要复核：票 174 自己的腿，或另派一枚带尺的腿 |
| J-7 | `internal/risk/paths.go:122-126` 这一枚候选 | **坐标今天不存在**：`internal/risk/` 下没有 `paths.go`（`ls internal/risk/*.go` 实测，文件清单见 §6.2）。是"坐标失效"，⛔ 不是"这一侧没问题" | 编排者：票 174 与票面排程引用的这枚 file:line 需要更正（追加，不改已推历史） |
| J-8 | 外发正文里**除工具失败句之外**还有哪几族带路径（如 system prompt／记忆注入段／spill 回执） | 本格射程＝工具失败句那一条链；其余族我只在 §2 增量里点名了我查到的，没有做全仓穷举（那是票 174 的射程） | 票 174 的名册轴／或另立票 |
| J-9 | 面板／球侧今天把哪一句显示成什么 | `frontend/**`、`design/**` 我没读（G-2 具名）；我只读到产码侧 `bridge.go:374 Args: req.Args` 与 `displayPaths` 这两条喂卡片的线 | 如需"用户看得见"的像素级读数：派一枚前端腿（可只读） |
| J-10 | `internal/tools/task.go:837` 那一枚接缝**今天到底接没接** | 它旁边的用例（`task_output_canonicalize_fail_174_test.go:130-148`）自己就分了 wired／unwired 两支断言，"生产装配走哪一支"是运行时选择 | 写码腿或乙形仪器；§6.2 只报"这一处今天有没有可断言的形状" |
