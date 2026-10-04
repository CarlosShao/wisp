# 174-a4 只读普查 — 仓里哪一些错误文本今天真的能走到模型眼前

- 起手锚：时间 `2026-10-04T09:25:03+08:00`（`date -Iseconds`）· 起手 HEAD `79b3982b`。
- 起手 `git status --porcelain internal/tools cmd/wisp tools/d22scan`＝**空**（派单说有两枚写腿在飞；此刻这三条路径上工作树与 HEAD 一致，读数无漂移风险；本程交件前复量一次）。
- 骨架提交 `0693a077`（第 9 轮，先于填表）。本程只读：产码／`docs/**`／票面 AC 框**零字节**，只写这一枚文件。
- 口径声明：以下每一跳都是**读码**读出来的，一条 Go 命令都没跑（禁 `go test`／`build`／`vet`／`run`）。
- 本文所有路径为**全路径**。所有 `file:line` 按 HEAD `79b3982b` 已提交面复量。

**这一程换了个方向数**：174-a3 数的是"显式调用者枚数"，它答不了 owner 那句问句。本程按项目定式**从终点倒读**——先钉死"模型眼前"到底是哪几个入口，再问哪些错误文本今天真走得进去。

## §0 三个终点的定义处（先立尺，再量名册）

| 终点 | 入口定义处（唯一） | 推导式 |
|---|---|---|
| ① 模型可见文本 | `internal/agent/loop.go:1080-1084` `toolResultMessage()` → `internal/agent/loop.go:722` `l.append(...)` → `:573` `l.asm.Build(in, l.History())` → C5 provider，**出网** | 回执进上下文的**唯一**一句是 `loop.go:722`；除此之外模型读不到工具的任何东西 |
| ② 面板可见文本 | `internal/agent/approval/gate.go:607` `Paths:` 装进 `Prompt`；卡片视图 `internal/panel/approval.go:39-57` `ApprovalCardView` **没有 Paths 字段**，`internal/panel/pump.go:103` 只映 `Args` | 见 §2② |
| ③ 日志／落盘 | 日志＝`cmd/wisp/logsink.go:149` 装的 JSONL 滚动文件（`teeHandler` 同时镜像 stderr）；落盘＝`internal/memory/schema.go:69-86` `tool_call` 表 + `internal/memory/artifacts.go` | 见 §2③ |

**终点①的两条注入喉（今天都活着）**
1. `internal/agent/loop.go:698` `log.Text = out.Text`（工具回执正文）。
2. `internal/agent/loop.go:696` `log.Text = execErr[i].Error()`（工具返回的 Go error 正文）。

**喉上没有任何脱敏**：`grep -rn 'Redact\|Sanitize' --include=*.go internal/agent internal/tools | grep -v '_test.go'` → **1 枚**（`internal/tools/bridge.go:1166 sanitizeArgs`，管的是 SQLite 的 `args_json` 列，**不在**出网那一跳上）。⇒ 回执正文原样出网。

## §1 名册：错误串里带变量位的构造点

### §1.1 全仓分母（口径与推导式，原样命令）

| 尺 | 命令（原样） | 读数 |
|---|---|---|
| 含 `%w` 的非测试行 | `grep -rn '%w' --include=*.go internal cmd \| grep -v '_test.go' \| wc -l` | **272** |
| 其中**同行另有插值动词**（`%s/%q/%v/%d`）＝错误串自己带变量位 | 上接 `\| grep -E '%[sqv]' \| wc -l` | **136** |
| 纯转发（只有 `%w`、不另插值） | 272 − 136 | 136 |
| 非测试文件里 `err.Error()` 拼进字符串的**装配点**（tools+agent，含审批门） | `grep -rn 'err\.Error()' --include=*.go internal/tools internal/agent \| grep -v '_test.go' \| grep -vE ':[0-9]+:\s*//' \| wc -l` | **52** |
| └ 其中拼进 `Result{Text:…}`（＝直通终点①） | `… \| grep 'Text:.*err\.Error()' \| wc -l` | **38** |

**136 枚生产者按文件全列**（尺＝`grep -rn '%w' … | grep -E '%[sqv]'` 后按文件聚合；这是"错误串自带变量位"的完整名册，一行不漏）：

```
cmd/llmrecord/main.go 100 | cmd/wisp/doctor.go 294 | cmd/wisp/firstrun.go 80,86
cmd/wisp/panel_config_store.go 272 | cmd/wisp/panel_inbound.go 232,264
internal/audio/wavinjector.go 55 | internal/llm/golden/replay.go 152 | internal/llm/probe_health.go 277
internal/memory/artifacts.go 82,128,155,164,205,210,215,265,287,289,308,339,349 (13)
internal/memory/dao_misc.go 148,216,231,276 | dao_profile.go 59,126 | dao_providerhealth.go 208
internal/memory/dao_tasklog.go 41,66,71,85,121 | internal/memory/open.go 337,487,519
internal/models/bridge.go 47,60 | downloader.go 239,590,640
internal/models/manifest.go 118,133,137,140,143,146,149,152,159,162,165,169,172,177,180 (15)
internal/panel/assets.go 88,91 | attachments.go 270,290,393,398 | bridge.go 129,133,136,140
internal/panel/composer_dispatch.go 219,231 | composer_handlers.go 118,129,134
internal/panel/config_handlers.go 271,282,295,310,316,322,329,335,350 (9)
internal/proc/envfork.go 222 | shutdown_hooks.go 96 | singleinstance_windows.go 54,56,72,134,138
internal/risk/pathresolver.go 95 | internal/secret/configrefs.go 29,33 | list.go 41,53
internal/secret/migrate.go 93,97,143,170,173,175,190,194 | refs.go 34,40,44
internal/secret/store.go 50,110,115,117,121,143,163,165
internal/tools/fs_write.go 103 | internal/tools/recycle_windows.go 193
internal/tools/registry.go 122,130,133
internal/winsec/placement_windows.go 34,38,46,66 | resolve.go 321,610,613,622,632,669,672
internal/winsec/winsec.go 258,262,373 | winsec_other.go 158 | winsec_windows.go 421,604,658
```

### §1.2 名册里"变量位＝路径／配置值／用户输入"的关键枚（逐枚：变量位是什么 · 今天被谁接住）

| # | file:line | 变量位是什么 | 今天被谁接住 |
|---|---|---|---|
| P1 | `internal/risk/pathresolver.go:95` | **两枚路径原文**（`r.Spelling`＋`r.Canonical`）＋改写构造名 | **只被 `Result.Actable()` 的调用者接**：`internal/panel/workspace.go:85,87`、`internal/tools/paths_workspace.go:64`、`internal/risk/syncdirs.go:208,226`、`internal/risk/winsec_c26.go:44,58`。⚠ **`internal/tools/paths.go:104-107` 的 `resolve()` 不调 `Actable()`**（a3 这一条复认成立）⇒ fs 工具那一条锥**取不到** P1 |
| P2 | `internal/winsec/placement_windows.go:66` | 路径原文＋reparse 组件名 | winsec 落封内层 → 上层 `observe.Wrap` 或 `refusing to seal %s: %w`（P3）→ 见 §2③ |
| P3 | `internal/winsec/resolve.go:610,613` | **`input`＝调用方递给封存的那条路径原文** | 同上；生产入口＝`spill.go:152,166,343`／`memory/artifacts.go:75`／`memory/open.go:176,188,503`／`secret/store.go:49,79`／`config/parse.go:215`／`config/migrate.go:93`／`secret/migrate.go:169,174,189` |
| P4 | `internal/winsec/resolve.go:622,632` | `input`＋`resolverLabel(r)`＋**答复串 `%q`（＝C26 改写后的路径）**＋`problem` | 同上；`resolverLabel` 是注册名（枚举名），不是用户数据 |
| P5 | `internal/winsec/resolve.go:669,672` | `input`＝路径原文＋`problem` | 同上（经 `risk/winsec_c26.go` 的 C26 管道回填） |
| P6 | `internal/winsec/resolve.go:321` | `resolverLabel(r)`＝** resolver 注册名**（枚举名） | 未装 C26 管道时才响；生产已装（`internal/risk/winsec_c26.go:21` `winsec.SetPathResolver`）⇒ 今天不可达 |
| P7 | `internal/tools/paths_workspace.go:70` | `res.Canonical`＝**工作区路径原文** | `internal/panel/workspace.go:81` → 审计行 `:93-96`（日志＋stderr）＋ `:97` `fmt.Errorf("panel: workspace switch refused: %w", err)` |
| P8 | `internal/tools/paths_workspace.go:74` | `res.Canonical`＋**配置项名字本体**（`[fs] allowed_dirs`） | 同 P7 |
| P9 | `internal/tools/paths_workspace.go:86` | `root`＝`%q` 路径原文 | 同 P7 |
| A1 | `internal/tools/fs.go:144` `fs.read`／`:208` `fs.list` | 转发 `Canonicalize` 的错误（今天＝静态集，见 §1.3） | 直通终点① |
| **A2** | **`internal/tools/fs.go:148`／`:158`／`:212`／`:218`** | **`*fs.PathError` 正文——内含 `os.Open` 用的那条 `canon`＝C26 展开后的真实绝对路径** | **直通终点①，今天就在走** |
| A3 | `internal/tools/fs_edit.go:138` | `*fs.PathError`（目标路径原文） | 直通终点① |
| **A4** | **`internal/tools/fs_write.go:343`／`:477`／`:540`** | **`os.Rename`／跨卷的错误——`rename <旧> <新>: …`＝两条路径原文同时在场** | 直通终点① |
| A5 | `internal/tools/fs_write.go:292`／`:325`／`:505`／`:514`／`:559`／`:620`／`:408` | 临时文件创建／Close／读源／建目标／删源／回收站／删除的 `*PathError`（路径原文） | 直通终点① |
| A6 | `internal/tools/fs_write.go:103` | `已停止于 %s：%w`（步骤名＋上游 why，why 里含 A5 的路径） | `fs_write.go:88 s.stop(step, err.Error())` → `head()` → 回执 |
| A7 | `internal/tools/fs.go:137`／`:201`、`fs_edit.go:108`、`fs_write.go:247`／`:382`／`:453`／`:604`、`task.go:509`／`:693`、`subagent_197.go:244`、`agent/tools.go:171`、`bridge.go:289` | **用户输入的碎片**（`encoding/json` 的报错会把出错处的字面带出来） | 直通终点① |
| A8 | `internal/tools/bridge.go:561` `工具内部故障：` | 转发工具 Execute 返回的 error（＝A1–A7 全集的兜底汇流） | 直通终点①；`internal/agent/loop.go:660` → `:696` |
| A9 | `internal/tools/task.go:837` | 转发 `Canonicalize` 错误（今天＝静态集） | 直通终点①；**接缝今天接不上**（§3.3） |
| A10 | `internal/tools/task_backfill.go:136` | **转发 spill 的整条错误锥——含 P1/P2/P3/P5 的路径原文** | `cmd/wisp/run.go:1110` → `:1112-1114` `rt.auditf(...)` → `run.go:931-935`（stderr ＋ JSONL 日志）。**不进终点①**（`Backfill` 在 `bg.Wait()` 之后） |
| A11 | `internal/agent/approval/gate.go:255`／`:505` | `admitCheck` 的错误原文 | `:254/:504` → `AnswerReject` 的 `why` → `bridge.go:1020 Text: why` → **终点①** |
| A12 | `internal/agent/approval/gate.go:297`／`:510`／`:514`／`:537` | 确认界面／审批队列／令牌签发的错误原文 | 同 A11（`:537` 走 `q.abandon` → `queue.go:501-502` → 同 `why`） |
| A13 | `internal/tools/bridge.go:942` ` (无法规范化: ` | `Canonicalize` 错误原文，拼进 `Decision.Paths` | `bridge.go:296` → `gate.go:607 Prompt.Paths` → 面板侧（§2②）＋审计行 |
| A14 | `internal/tools/cancel.go:76` | 上下文取消错误（今天静态） | 回执文本 |
| A15 | `internal/tools/subagent_197.go:322` | 子代理环路装配错误原文 | 直通终点① |
| A16 | `internal/agent/loop.go:821` | 工具目录不可用错误原文 | `Text:` → 终点① |
| A17 | `internal/agent/loop.go:920-921` | 任务错误原文（`res.Message` ＋ `EvError.Text`） | `:1284-1286` 终端印 `e.Err.Detail`；`EvError` 不入上下文 |

**对照组（证明这把尺分得开两类，全部零变量位）**：`internal/risk/pathresolver.go:36`、`:49`、`internal/agent/spill.go:229`、`:232-234`、`internal/tools/paths.go:121`（`tools: empty path`）、`internal/winsec/resolve.go:54`。
尺的推导式：`grep -n 'errors.New' internal/risk/pathresolver.go internal/agent/spill.go internal/tools/paths.go` → 逐字读出无 `%`。

### §1.3 `Canonicalize` 今天的可达错误集（A1/A9/A13 的上游，闭合枚举）

`internal/tools/paths.go:122-126` → 只有两支：`fmt.Errorf("tools: empty path")`（`:121`，静态）与 `p.resolve()` 透出的 `risk.ErrReparseDenied`（`internal/risk/pathresolver.go:123` 产出，静态）。**`Actable()` 不在这一条路上** ⇒ P1 今天进不了 fs/task/bridge 那三枚装配点。

## §2 三种终点分开判

**默认档自己复跑的读数（不照抄 a3）**：
- 日志级别默认＝**info**。三处独立证据：`internal/observe/logging.go:79` `level := slog.LevelInfo`（`cfg.Level` 为空时就是它）、`cmd/wisp/logsink.go:82` `const logSinkLevel = "info"`、`cmd/wisp/slo_windows.go:274` `Level: "info"`。⇒ `Warn`／`Error` **一律落盘**。
- **路径脱敏默认关、且今天根本接不上**：`internal/observe/redact.go:138` `if r.RedactPaths { …mask… }` 是唯一的路径掩码点；`internal/config/schema.go:509` `RedactPaths bool … default:"false"`；两处 `LogConfig{}` 字面量（`cmd/wisp/logsink.go:149`、`cmd/wisp/slo_windows.go:272`）**都没有写 `RedactPaths`** ⇒ 零值 `false`；`cmd/wisp/config_readers_255.go:134` 逐字记着 `the mirror field exists and is filled by callers, never from cfg.Privacy`。⇒ **日志里的路径原文不会被掩掉**。
- 内联密钥掩码（`redact.go:137 maskInlineSecrets`）**与路径无关**，它只盖 `sk-`/`Bearer`/`api_key=` 形状。

| 名册组 | ①模型可见 | ②面板可见 | ③日志／落盘 |
|---|---|---|---|
| **A2／A4／A5／A3（OS 级 `*PathError`，含真实路径）** | **会，今天就在走**。生产注册点＝`cmd/wisp/run.go:515` `tools.BuiltinFSEntries(tools.FSDeps{…})`，而该函数（`internal/tools/fs.go:325-331`）**逐字 append 了整个写入半族** ⇒ `fs.read`/`fs.list`/`fs.write`/`fs.edit`/`fs.trash`/`fs.move` 全在册；`fs.delete` 仅当 `[fs] delete_enabled=true`（`internal/tools/fs.go:320-323` 注释逐字）。装配点＝`fs.go:148/158/212/218`、`fs_edit.go:138`、`fs_write.go:292/325/343/408/477/505/514/540/559/620` | 不会（回执不上卡；卡只映 `Args`，`pump.go:103`） | 不会（没有一行日志写回执正文——`grep -rn 'slog\…' internal/agent internal/tools \| grep -i text` → **0 枚**；终端 `cmd/wisp/run.go:1278-1279` 只印工具名＋结果词） |
| A1／A9／A13／A14（`Canonicalize` 转发，今天＝静态集） | 会（不含路径，含 C26 词汇本体 `reparse_point_exceptions`／`(junction/symlink)`） | A13 会（`gate.go:607`） | 不会 |
| A7（JSON 报错，含用户输入碎片） | 会 | 不会 | 部分会：`bridge.go:1115 sanitizeArgs(req.Args)` → SQLite `tool_call.args_json`（`internal/memory/schema.go:74`，2048 字节内**不掩路径**） |
| A8／A15／A16（兜底汇流） | **会**（继承全部上游变量位） | 不会 | 不会 |
| A11／A12（审批拒绝 `why`） | **会**（`bridge.go:1020 Text: why`） | `why` 同时是卡上的 `Reason`（`internal/panel/approval.go:47-48`） | 会（`gate.go:251 g.logf` 等） |
| A10（spill 错误 → `Backfill` 的 `why`） | **不会**（在 `bg.Wait()` 之后，`cmd/wisp/run.go:1106`→`:1110`） | 不会 | **会**：`run.go:1113 rt.auditf` → `run.go:931-935` 同时写 stderr 与 JSONL 日志文件 |
| **P1（`pathresolver.go:95`，两枚路径原文）** | 不会（`Actable()` 不在工具锥上，§1.3） | **半会**：拒绝时面板侧拿到的是 `internal/panel/workspace.go:99` 那句 `%w`（内含两枚路径原文）；但**渲染成人的那句话**不含路径——改写分支的 `Reason`（`:110-111`）只印改写构造名，正常拒绝走 `currentAfter`（`:114-122`）。**例外形状具名**：`internal/panel/workspace.go:118-122` 在"拒绝后范围真的变了"时把两枚路径 `%q -> %q` 印进 `Reason` | **会**：`internal/panel/workspace.go:93-96` 审计行逐字 `err=%v … detail=%q`，成功分支 `:100-102` 另把 `spelling=%q` 落盘 |
| P2／P3／P4／P5／P6（winsec 封存锥） | 不会（winsec 错误到不了回执；见 §3.2） | 未量到（`internal/panel/config_handlers.go` 那 9 枚会不会把 P3 原文上界面——没读） | **会**：`internal/agent/loop.go:711 Warn("agent: spill failed","err",err)` → JSONL；＋ A10 那一条；＋ winsec 自记 |
| P7／P8／P9（工作区拒绝，含路径＋`allowed_dirs` 键名） | 不会（工作区切换是**面板侧**动作，`internal/tools/paths.go:42-47` 注释逐字 "a panel-side workspace switch"；全仓无工具暴露它） | **会**（`RequestWorkspaceSwitch` 的返回错误由面板侧处理） | **会**（`workspace.go:93-96` 审计行 `from=%q requested=%q err=%v`） |

## §3 那一跳到不到的判据（逐跳到底）

### 3.1 A2／A4／A5（OS 级错误 → 模型）：**接到，今天就在走**
产：`internal/tools/fs.go:146` `os.Open(canon)`（`canon` 来自 `fs.go:142 → :93 FSDeps.open → internal/tools/paths.go:122 Canonicalize → internal/risk/pathresolver.go:105 Resolve`，**C26 已把 `~`／`%APPDATA%`／8.3 短名展开成真实路径**）。
拼：`internal/tools/fs.go:148` `"打开失败：" + err.Error()`（Go 的 `*fs.PathError.Error()`＝`open <整条路径>: <系统消息>`）。
投：`internal/tools/bridge.go:577 Text: res.Text` → `internal/agent/loop.go:660` → `:698 log.Text = out.Text` → `:707-713` spill（不吞正文，`loop.go:712`）→ `:722 l.append(toolResultMessage(…))` → `:573 asm.Build` → C5 provider，**出网**。
旁路（本程新抓到的一跳）：`internal/tools/bridge.go:646` `b.prov.MarkWithHostPath(dec.TaskID, dec.Tool, origin, res.Text, hostPath)`——**同一枚 `res.Text`（含路径）还被喂进 C25 溯源面**；这一支的下游本程没读完（见 §5 第 5 条），但它说明带路径的回执不止出网一条去路。
判：**没有断点**。这一条不需要等票 174。

### 3.2 P2／P3／P5（winsec 路径原文 → 模型）：**断在"投递"这一跳**
产：`internal/winsec/resolve.go:610`／`placement_windows.go:66`。投的候选只有两条，两条都断：
- 断口 A：`internal/agent/spill.go:152-153`（`PrivateDirAll` 失败 → `observe.Wrap(observe.ClassResource, err, "agent: create artifacts dir")`）→ 调用者是 `internal/agent/loop.go:709 sp, err := l.sp.Prepare(…)`，而 `:710-711` 的处理是 **`Warn` 进日志、不进 `log.Text`**；`log.Text` 保留未落盘的全文正文。⇒ 模型拿到正文，**拿不到这句错误**。
- 断口 B：`internal/tools/task_backfill.go:130-136`（`b.Spills.Prepare` 失败 → `spillWhy` 拼上 `err.Error()`）→ `cmd/wisp/run.go:1110` 的 `Backfill` 位于 `:1106 res := bg.Wait()` **之后**，环路已收，随后进程退出；`why` 只被 `:1113 rt.auditf` 消费。⇒ **同一形状，落在日志侧**。
具名断点：**`cmd/wisp/run.go:1106`（`bg.Wait()`）与 `:1110`（`Backfill`）的先后顺序**，以及 **`internal/agent/loop.go:711` 走日志不走正文**。⛔ 不写成"永远不会"：AC#2b 把读回接进生产、或 `loop.go:711` 改成把 `err` 写进回执，这两跳任一动一根就通。

### 3.3 A9（`task.output` 注意句 → 模型）：**接缝今天接不上（复认 a3）**
唯一生产者 `internal/tools/task_backfill.go:138` `rec.ArtifactPath = sp.Path`；`grep -rn 'ArtifactPath *=\|ArtifactPath:' --include=*.go internal/ cmd/ | grep -v '_test.go'` → **恰 1 行**（复跑到 HEAD `79b3982b` 仍成立）。它唯一的生产调用点 `cmd/wisp/run.go:1110` 在 `bg.Wait()` 之后 ⇒ `pointerNotice`（`internal/tools/task.go:829-851`）今天**不被喂非空路径**。子代理行走 `internal/tools/subagent_197.go:424` 不带 `ArtifactPath`。
⚠ 本程修正 a3 一处归因：A9 那一条**不是**"这一句漏不漏"的孤例，它只是 §2① 那一大堆装配点里最晚接通的一枚。

### 3.4 P1（`pathresolver.go:95`）：**模型侧断在生产这一跳，日志／面板侧通**
生产确实造它：`internal/panel/workspace.go:85`、`internal/tools/paths_workspace.go:64`、`internal/risk/syncdirs.go:208,226`、`internal/risk/winsec_c26.go:44,58`。
断口：工具锥上的 `Canonicalize`（`internal/tools/paths.go:122`）**不调 `Actable()`** ⇒ P1 进不了 §1.3 那个闭合枚举 ⇒ 到不了 `fs.go:144`／`task.go:837`／`bridge.go:942`。
通的一面：`internal/panel/workspace.go:93-96` 把它整句（含两枚路径）写进审计行 → 日志文件＋stderr。

### 3.5 A11／A12（审批 `why`）：**接到**
`gate.go:255`（或 `:297/:510/:514`）→ `tools.AnswerReject, why` → `internal/tools/bridge.go:358/:283 b.reject(…, why, …)` → `bridge.go:1020 Text: why` → 与 3.1 同一条投递链 → 出网。今天可达与否取决于 `admitCheck`／队列／UI 传输这三枚错误今天带不带变量位——**本程未逐枚读到它们的产出点**（见 §5）。

## §4 残余形状（AC#2b 那格）

**那 11 枚在哪（自己数，不复述 a3）**：`internal/tools/task_output_pointer_notice_test.go:474-489`。尺＝`grep -n 'for _, banned := range' internal/tools/task_output_pointer_notice_test.go` → **两枚命中：`:219`（6 项，另一格）与 `:474`（11 项）**。`:474` 逐枚：`filed`、`leakProbeSegment174r3`（`:392`＝`ALLOWLIST-LEAK-PROBE-174r3`）、`leakProbeResolved174r3`（`:397`）、`rootA`、`rootB`、`读得回来`、`随时可读`、`可以读回`、`不在你被授权的目录范围内`、`未接线`、`不可找回` ＝ **11 枚**。

**接上那天会不会响？不会响——但 a3 给的理由要改。**

- a3 的理由是"11 枚里不含 C26 词汇本体"。这句**对**（`reparse_point_exceptions`、`junction/symlink`、`risk:` 前缀**都不在表里**），但它只解释了"词汇泄漏不响"这一半。
- **表里有 3 枚是路径字面**（`filed`、`rootA`、`rootB`）＋2 枚毒值字面。⇒ 这枚钉**对路径泄漏是有牙的**：AC#2b 当天若 `Canonicalize` 返回的是一条**会把路径原文抄进错误串**的错误（P1/P2/P3 那种形状），注意句里就会出现 `filed`，`:474` 当场红。
- 今天**不响**的确切原因＝ §1.3 那个闭合枚举：`Canonicalize` 今天只会返 `tools: empty path` 或静态的 `ErrReparseDenied`，两句都不含 11 枚中的任何一枚。所以**"不响"是"可达错误集恰好全静态"的结果，不是"这枚钉射程不到"的结果**。

**要钉在哪一枚才叫有牙（只给候选，不替它写码）**：

| 候选 | file:line | 牙在哪 |
|---|---|---|
| **甲（首选）** | `internal/tools/paths.go:122-126`（`Canonicalize` 的返回支） | 钉"返回给工具层的错误**只能**来自一个闭合枚举"——今天把它钉成 `{tools: empty path, risk.ErrReparseDenied}`，将来任何人往里加一枚带 `%s` 的支就当场红。这一枚把 §3.1 那种"无界开口"从**约定**变成**判据** |
| 乙 | `internal/tools/task.go:837`（拼接点） | 钉"注意句区间内不得出现 `risk:`／`winsec:`／`tools:` 前缀与 C26 词汇本体"。补的是 `:474` 那 11 枚**没覆盖的词汇那一半** |
| **丙（本格真正的洞）** | `internal/tools/fs.go:144`＋`:148`、`internal/tools/fs_write.go:343` | 现有 11 枚**射程只在 `task.output` 的注意句上**；§3.1 那条今天就在漏的路径**没有任何一枚钉**。要牙就钉在这里——不写码，只登记这一格为空 |

## §5 我读到的 vs 我没量到的

**读到的**：§0 三终点的定义处与"喉上无脱敏"；§1 全仓分母 272／136／52／38 及其推导式；§1.2 逐枚变量位性质与接住者；§1.3 `Canonicalize` 可达错误集；§2 默认档三项现量（info／`redact_paths=false`／两处 `LogConfig` 都没写该字段）；§3.1–§3.5 五跳到底；§4 11 枚逐枚与两处射程判断。

**没量到的（具名，"应该没问题"不算结论）**：
1. **一条 Go 命令都没跑**（派单禁）⇒ §3.1"接到"是**读码链**，不是实测链。尤其 `*fs.PathError` 的正文形状是 Go 标准库的公开形状，我没在本仓任何一发实测里看到它落进回执。
2. **A11／A12 的上游产出者**（`admitCheck`、`q.push`、`grantNonce`、面板桥的传输错误）今天各带不带变量位——没逐枚读到产出点，只读到"会原样拼进 `why`"。
3. **P2–P5 会不会上界面**：`internal/panel/config_handlers.go:271-350` 那 9 枚、`internal/panel/attachments.go:270-398` 那 4 枚的**错误消费端**没读——§2② 那三行"未量到"就是这个。
4. **`Prompt.Paths` 会不会被序列化进面板 JSON**：读到 `gate.go:607` 装了、`queue.go:530`/`pending_read.go:65` 透传了，但**最终写出点**没读。⇒ A13 的②栏按"到 Prompt 为止"判，界面那一格未定案。
5. **C25 片段索引**会不会把带路径的回执本身收进 taint 片段再回喂（`internal/risk/provenance.go` 消费端）——a3 也登记为缺口，本程同样没读。
6. **136 枚逐枚的性质判断**：我只对其中 risk/winsec/panel/tools 四组（48 枚）逐枚读了变量位；`internal/memory/*`（26）、`internal/models/*`（17）、`internal/secret/*`（24）、`internal/proc/*`（7）四组＝**只数了枚数、没逐枚读实参**，它们今天到不到终点①未判（先验：这四组不在工具回执锥上，但"不在"是我推的，不是读到的）。
7. **门禁实测**（`tools/d22scan`、`scripts/*.sh`）——禁跑，一律未验。
8. `redact_paths` 设成 true 之后掩码正则的实际覆盖面（`redact.go` 的 `pathRe` 词形）——没逐字读那条正则。

## §6 一句话给 owner（零术语）

今天会发生的事：这个程序在替你做文件操作而失手时（打不开、读不动、改名失败、放不进回收站），会把系统给出的那句原因**原样转给云端模型看**，而那句话里**通常带着这条文件在电脑上的完整真实位置**——包括程序替你展开过的部分（你写"～"、它认到的是你的真实用户目录名）。今天不会发生的事：日志里的脱敏开关默认是关的，所以这些位置也**明明白白写在你硬盘上的日志文件里**；但你的**授权目录清单**、以及"为什么这条路被挡"那一句里目前只含一句固定英文，不含你的路径。哪一天开始可能会更多：本程没等来那一天——那一格（让模型自己去读长答案）今天确实还没接上，接上的那天，"为什么被挡"那句会连同它当时手里那整条路径一起出去，而**现有 11 枚禁词钉里对路径有牙、对词汇没牙**，并且**今天这些已经出去的路径那一条上，一枚钉都没有**。要不要动、动哪一枚，是你拍，不是我拍。

---

**交件前复量（只追加，上文一字不改）**

- HEAD 从起手 `79b3982b` 走到 `9d089603`（编排者在途中落了账）；派单说在飞的两枚写腿，本程起手与交件前两次 `git status --porcelain internal/tools cmd/wisp tools/d22scan` **都是空** ⇒ 本文引的那三条路径上读数无未提交漂移。
- 两处本文自查出的锚更正：①`internal/tools/paths.go` 那句 `tools: empty path` 在 **`:121`**（初稿误写 `:123`，已就地改）；②§1.3 那个"闭合枚举"**复量成立**——`Canonicalize`（`internal/tools/paths.go:119-149`）全文只有两枚 `return …, err`，其后的 `sameFormOfUnresolved(res.Canonical)`（`:149`）**只返字符串、不返错误**，所以它不可能带进第三枚错误形状。
- 生产注册这一格复核到底：`fs.delete` 的条件在册说明见 §2① 首行；`fs.edit` 随写入半族一并注册（`internal/tools/fs.go:330` → `internal/tools/fs_write.go:772`）。
