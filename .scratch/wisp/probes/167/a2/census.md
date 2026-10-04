# 167-a2 只读普查：票 167 AC#1 剩下的两格——**草稿不吞（AC#5）**与**崩溃自救（AC#6）**

**取数时刻**：2026-10-04 09:4x–09:5x（腿 `167-a2` 于 09:5x 前交回读数；编排者复跑与落笔同发，见 §7）
**锚点**：腿自报起手 HEAD（未复述），编排者代落此件时 `git log -1`＝见 §7 那行。
**⚠ 本件由编排者代落**：派单把交付写面给了 `167-a2`，而该腿的运行类型**没有写文件的工具**（只读检索型）⇒ 它**零 commit、census 未建立**，全文以回报正文交回。定式（本仓既有）：**要腿写文件，先确认该类型有写工具；没有＝编排者代落并逐行自跑尺**。本件按那把尺办：**§1–§6 的每条读数都注明〔编排者复跑〕或〔仅自述〕**。

## 1. AC#5 草稿不吞：Go 侧今天**没有任何一格留住发送失败的原文**〔编排者复跑：`cmd/wisp/panel_inbound.go:272-277` 逐字 `Message:    nil`，其上注释逐字 `workspace = ticket 186, attachment = ticket 92, message = ticket 35`〕

链路（后四跳＝〔仅自述〕，编排者未逐跳复跑，归 `167-r1` 起手现读）：
`panel.message.send` ⇒ `internal/panel/composer_dispatch.go:194` ⇒ `:230 unattached` ⇒ `:240 record`，而那行审计**不含 `req.Text`**（`:244` 的 detail 是固定句子）；回给页面的只有 `[method requestId] ＋ err`（`internal/panel/bridge.go:157-166`）。
⇒ **判语：这一格今天落在页面那一侧**——Go 侧既无字段承载、也无"发送失败"这条分支的原文留存。⚠ 「页面今天自存不自存」**量不到**（`frontend/**` 不属本编队射程，票 167 Status 逐字 `skipped=frontend(owner-delegated)`）⇒ **归前端那支，由机主自己带给他用的那枚 agent**。

## 2. Go 侧承载位穷举（〔仅自述〕；口径＝`internal/panel` 非测试 struct 39 枚逐枚看）

- **出向可装"未发出原文"的字段＝0 枚**。唯一文本槽是 `ResultChunk.Text`（助手流）；`ApprovalCardView{Reason,Args}`、`TaskRowView{Label,StatusReason}`、`SettingsView` 族、`ComposerState` 的 11 键**零自由文本位**。
- **入向文本载体＝2 枚**，且都是一次性封套：`internal/panel/bridge.go:98`、`internal/panel/composer.go:296`。
- **持久面＝0**：`internal/memory/schema.go` 建表 9 枚（`:23/:30/:40/:53/:69/:89/:102/:111/:136`），其中 `task_log.query_text:58` **只给已受理的任务**；`internal/store` 这个包**不存在**。
- 数据根下文件面名册＝`config.toml`／`logs\`（`cmd/wisp/logsink.go:88`）／`models\`／`artifacts\`／`panel-webview2\`／`secrets\`／`memory.db`。
- **涉及 "draft" 的测试＝9 行／5 文件，真草稿用例 0 枚**（全是"本文第一稿"这类散文，或 `internal/tools/fs_edit_ac34_test.go:75,86` 的 `status: draft` 桩文本）。同族"新的不许顶旧的"只有配置那两枚：`cmd/wisp/always_write_no_clobber_226_test.go:39`、`internal/config/writeguard_226_test.go`。

## 3. AC#6 崩溃自救：旧假话确认被推翻，但**新量到一枚咬人的**

- **推翻旧假话**〔编排者复跑：`internal/observe/goroutine.go:174` 逐字 `"stack", ev.Stack`；`:236` 逐字 `sink: defaultPanicSink`（`NewRegistry` 里就装了）；`:167` 是那个 sink〕⇒ **生产确实装了 panic sink 且逐字带栈**，10-01 那条"生产没装 ⇒ 盘上零栈"是假话（账已撤，票 249 撤回）。
- ★ **新读数：那枚栈进文件之前会被截到 512**〔编排者复跑，三处全中〕：
  - `internal/observe/redact.go:37` 逐字 `MaxLoggedString = 512`（其上 `:36` 注释逐字 `bounds any single string that reaches the log (rule 4)`）；
  - `internal/observe/redact.go:133-142` 的 `Redactor.String` 末尾逐字 `return truncate(s, MaxLoggedString)`（同函数里 `if r.RedactPaths { s = pathRe.ReplaceAllString(s, "<path>") }`＝**路径脱敏这枚开关的真实落点**）；
  - `cmd/wisp/logsink.go:153-160`：`teeHandler` 的 `primary: p.Handler()`（＝过脱敏流水线的 JSONL 文件侧）vs `mirror: slog.NewTextHandler(os.Stderr, …)`（**不过脱敏**），`:208-213` 把**同一枚 record** 交两边。
  ⇒ **盘上只有 512 rune 的栈头；完整栈只出现在 stderr**，而**双击／Explorer 启动没有控制台**（票 244 那条线）⇒ **机主手上拿不到现场**。
- ⚠ **要不要放宽这 512 那条＝动既有脱敏规矩（rule 4），属人工批准面**；编排者在此**不替机主裁**。合法备选（不裁，只登记形状）：甲＝放宽 `MaxLoggedString`；乙＝**另开一枚不受日志脱敏管线管的崩溃现场文件**（⚠ 这会新造一条"不掩码的外发/落盘面"，比甲更需谨慎）；丙＝不改规矩，只在回执里说实话（"盘上的栈只有前 512 字，完整现场需要带控制台的跑法"）。
- **自救指令那一半**〔仅自述〕：`cmdDoctor` 的检查项 10 枚（`cmd/wisp/doctor.go:29-131`）**零枚读 `logs\`**；`logSink.logDir()` 的产码调用者 **0 枚**；`SetPrintOrigin` 全仓 **0 命中**；`internal/watchdog` **只有 `doc.go`**（DEFERRED 票 42）⇒ **无守护/重启能力**。
- **产码 `recover()` 共 7 处**〔仅自述，枚数编排者未复跑〕：`internal/observe/goroutine.go:296`、`internal/memory/writer.go:144`、`internal/risk/assessor.go:278`＋`:292`、`internal/ball/hotkey_reload.go:146`、`internal/tools/bridge.go:823`＋`grant.go:70`＋`mode.go:40`；**`cmd/wisp` 产码零 recover**。
- **AC#6 三格判语（腿的原判语，编排者认其形状）**：**「记录通路已在、话说出口那一寸没接线（纯 Go 可闭合）；完整现场缺一整块（要动脱敏规矩＝须上报）；禁用某插件后重启应用缺一整块（票 50/51）」**。

## 4. 与前轮读数的一致性核对（只做一致性、不重取证）

`requestTaskStop` 在 Go 码 **0 命中**〔编排者复跑：`grep -rn 'requestTaskStop' cmd internal tools --include=*.go \| wc -l` ＝ **0**；台账散文里有 2 处提及，不算码〕；占用相关 grep 3 命中全为热键测试注释 ⇒ 前轮三条结论（占用零 reader／序号三真源在／停止＝4 枚写好零调用）**均立**。

## 5. 判不动的地方（逐条具名＋归口）

1. **页面今天自存不存草稿**——`frontend/**` 不在射程（票硬约束＋本编队 `skipped=frontend`）⇒ **归前端那支**。
2. **给草稿加一枚出向字段算不算越出票 145 已批的"局部解冻那一寸"**——现行范围文本不在我手上 ⇒ **编排者裁**（需先读 `A##` 那枚解冻的边界行）。
3. 〔待验，需跑一遍才知道〕**512 rune 实际剩几帧、够不够定位一次崩溃**——要一发真 panic 落盘再读 JSONL；⛔ **不是"做不到"**。
4. 〔待验〕**main goroutine 崩溃有没有 Windows／WER 兜底留痕**——源码面零 recover 已量，运行时面读不到。
5. 〔待验〕**这台机器 `%APPDATA%\wisp\logs\` 的实解析结果**——链已给（`internal/proc/envfork.go:50-69`，prod＝`wisp`／dev＝`wisp-dev`，`portable.txt` 覆盖在 `:226`），但不跑进程无法现量。

## 6. 给编排者的排程结论（不裁形，只报"能不能派"）

票 167 的 **AC#1 四格到此料齐**（占用／序号／停止＝前轮，草稿／崩溃自救＝本件）⇒ `167-r1` 的**现量闸门不再缺料**。但两格**今天不能开工**：① **AC#5 草稿**——Go 侧零承载位，真要做得先开出向字段（撞票 145 的解冻边界）且**页面那一半不归本编队**；② **AC#6 的"完整现场"**——要动 `MaxLoggedString` 那枚既有规矩＝**人工批准面**。⇒ **可派的那一寸＝AC#6 里"话说出口没接线"那一寸（纯 Go、不动规矩）＋ AC#2/AC#3/AC#4 那三枚（占用/序号/停止）**，写面要按包级互斥现量。

## 7. 门禁读数（本件自量）

- **本轮零 Go 命令**（腿自报，编排者复跑口径一致：只用 `grep`/`sed`/`Read`/`git`）。
- 编排者代落此件的同发读数：`date`＝**09:48:02＋08**，`git log -1`＝**`eb2c0173`**。
- 本件自量：`wc -c`＝**8,255**；占位符（"待填／填写中／TBD"字样）**0 枚**；七节齐、§5 那五条判不动项全带归口。
