# 票 219 AC#6 追补 —— 名册池号＋流截断三数的界面落点（写腿收件表）

> 2026-11-08 19:2x，前端写腿（frontend 波，票 219 AC#6 的**追补半边**）。
> 上一腿（15:1x，`219-ac6-task-roster-ui.md` §4）具名报了两处「数据已到页面但界面无落点」：
> ① 名册池号（`inFlightSlots/poolCap`）；② 流截断三数（行级 `streamTruncated/streamElidedRunes/streamDropped`）。
> 本波把这两组数画进任务监控弹层。本表只记**盘上可复跑的读数**。

## 0. 起手三核（现跑，不信转述）

1. **Go 侧真字段名**：`internal/panel/subagent_roster_197.go` 全文现读——池号＝`TaskRosterSection.InFlightSlots/PoolCap`（:162-163，json `inFlightSlots`/`poolCap`，**读名册自己的计数、不从 Rows 现数**，Go 侧注释写明槽在行之前被占住，两数是两个真相）；行级三数＝`TaskRowView.StreamTruncated/StreamElidedRunes/StreamDropped`（:140-142）；另有 section 级整卷三数 `StreamTruncated/StreamElidedRunes/DroppedCount`（:166-183）。**字段全在，不缺**——上一腿担心的「Go 侧没这字段」不成立。
2. **前端声明已齐**：`frontend/src/lib/panel.ts:224` `TaskRosterSection` 七键全声明（219/400 两波落的），本次**零改动契约文件**。
3. **接线现状**：`App.tsx` 现读——`agents={monitorAgents(snapshot.tasks)}`/`rosterPresent` 在（上一腿落的），池号与截断三数确实无人消费。

## 1. 落了什么

| 件 | 性质 | 内容 |
|---|---|---|
| `frontend/src/components/harness/monitor-popover.tsx` | 生产改形 | ① **池号**：「子智能体」组标题行右侧 `在跑 3 / 上限 8`（字形照 env 行 value 的 `font-mono tabular-nums` 既有形；读 wire `inFlightSlots/poolCap`，不从行数现数）；② **行级截断标**：行尾小字 `省 N 字`（折界）或 `流已丢`（丢弃，压过折界），叠加在状态字旁、**不参与状态判断**（显示层，不改 dot/档）；③ **section 级告示**：列表底一行 `流已截断：省略 N 字，已丢 M 条`（两损同句；仅折界⇒只省略数；仅丢弃⇒`已丢 M 条流`）；④ **诚实空态**：`tasks` 键缺席 ⇒ 三处全不画，**绝不显示假 0**（`poolInFlight===undefined` ⇒ 无池号；`streamTruncated` 缺省 false、`droppedCount` 缺省 0 ⇒ 无告示） |
| `frontend/src/App.tsx` | 生产接线 | `monitorAgents` 派生补带行级三数（verbatim，照 wire json 名）；`MonitorPopover` 新传五 props 全部 `snapshot.tasks?.` 直传（键缺席自然落 undefined，页面不造数） |
| `frontend/scripts/render-task-roster.tsx` | 牙扩展 | 四牙扩为九牙：池号牙、截断三态（双损/仅折界/仅丢弃）、干净日志**零告示**牙、行级两态（折界/丢弃/整行干净无标）、无键诚实空态（不画「在跑」不画截断）、App 级截断名册不崩，另三枚源级接线钉（池号/整卷三数/行级三数必须直传自 `snapshot.tasks`） |
| `docs/evidence/s1/219-ac6-roster-pool-truncation-live.png` | 活体证据 | 同一 bundle 供真 Chromium 推快照＋点开任务监控的截图，详 §3 |
| 本件 | 裁决表 | 读数与设计裁定 |

⛔ 没动：`cmd/wisp/**`、`internal/**` 产码（只读）、`internal/ball/**`、`design/**`、`PLAN.md`/specs、thresholds/golden、`lib/panel.ts`（契约声明本来就齐）。

## 2. 设计裁定（两条，理由在此）

1. **流截断只在真截断时显示**（派单给的两个选项里选后者）。理由：常态显示「省略 0 字」是告警噪音，会训练用户忽略它——与本仓 `App.tsx` 头注释对审批 chip 的同一论据（"A permanent row for something empty 99% of the time trains you to ignore it"），也是主流 harness 的告警做法（badge 只在异常点亮）。干净日志与整行干净都渲染零截断文字，牙里钉死。
2. **池号在键存在时恒显**（不属异常告示，是常驻容量读数——主流 harness 的 task/runner 面板把「running/capacity」当常驻元信息）。读 wire 计数不从行数现数：Go 侧 `TaskRosterState` 写明槽在行之前被占住、释放有窗口，页面不替它挑一个真相。

## 3. 门禁读数（全部本波现跑）

| 门禁 | 工作树（commit 前） | 导出树（`git clone --shared` @ `37df1ec64`） |
|---|---|---|
| `npm run typecheck` | **rc0** | **rc0** |
| `npm run build` | **rc0**（629.51 kB chunk 警告为既有） | **rc0** |
| `npm run render:roster` | **rc0，九牙全绿** | **rc0，九牙全绿** |
| `npm run lint`（oxlint） | **rc0**（0 errors；36 warnings 皆既有类） | **rc0** |
| `sh scripts/d22scan.sh` | 不适用（导出树为准） | **clean rc0**（frontend 121 text files；日志 `/d/tmp/wisp/20261108/219ac6b-d22scan2.log`） |
| `gofumpt` | **不适用**（本腿零 Go 产码） | 同 |

**d22scan 第一次红的自查与根因（如实记）**：导出树首跑 `TestLedgerCountsMatchAnIndependentWalk` 红（frontend 自报 121 vs 独立 walk 122）。排查：主仓同测试 PASS（121=121）⇒ 变量在导出树特有物＝本腿为跑 npm 门禁建的 `frontend/node_modules` **junction**：Windows reparse point 在独立 walk 的 `filepath.WalkDir` 里不算目录，被当**一个文件**数进（`node_modules/` 目录规则对「文件」不失配），而 scanner 的目录级判定正确跳过 ⇒ 121 vs 122。删 junction（仅删链接，真身 `frontend/node_modules` 现跑 `Test-Path` 证实完好）后复跑全量 **clean rc0**。教训归仪器的 Windows junction 盲区，不属本票写面，报编排者记账。
`go build -o wisp.exe ./cmd/wisp` **rc0**（51,013,194 字节；bundle 重嵌；环境变量按纪律 unset GOTMPDIR/TMPDIR/TEMP/TMP＋PATH 加 build/third_party/sherpa-onnx；未设 CGO_ENABLED=0）。
LF：三枚改动源文件现跑字节核 CR=0。

## 4. 活体证据（同一 bundle 供真 Chromium）

本机 wisp 的 WebView2 面板白窗属环境问题（账 A1631 一族），按派单用浏览器供活体：
`npm run build` 的 dist 经 `vite preview`（localhost:4199）供出，真 Chromium 加载后对页面调 `window.wispOnPanelSnapshot(<快照>)`（**同一扇接收门**）推一枚快照：`inFlightSlots=3, poolCap=8`；一行 `streamTruncated:true, streamElidedRunes:900`、一行 `streamDropped:true`、一行 `blockedOnApproval:true`；section `streamTruncated:true, streamElidedRunes:421, droppedCount:2`。点顶栏「任务监控」，截图：

![任务监控显示池号与截断](219-ac6-roster-pool-truncation-live.png)

- 弹层实文（Chromium `innerText` 现读）：「**在跑 3 / 上限 8**」「流被折界的子任务 / **省 900 字** / 执行中」「流页被丢的子任务 / **流已丢** / 思考中」「等批准的子任务 / 等批准」（上一腿的档未被本波破坏）＋底部「**流已截断：省略 421 字，已丢 2 条**」——五个落点一次上屏，零 emoji；
- **负例同跑**（第二枚快照，无 `tasks` 键）：弹层只显示「还没收到名册（宿主没带任务这一节）」，程序化断言 `hasFakePool=false, hasFakeTruncation=false`——**无键 ⇒ 无假 0**，与牙同判；
- 边界如实：这是**页面消费跳**的活体，不是 wisp.exe 端到端——快照由证据程序推入，未经 Go 泵真推（wisp.exe 本身已重编 rc0 供 bundle 嵌入账；Go→页面那跳的活体归 329 的账，WebView2 白窗见 A1631）。

## 5. 裁决建议（裁决者≠实现者）

- 池号、行级截断标、整卷截断告示、诚实空态四格盘上读数齐（§2/§3/§4）；
- 上一腿 §4 具名的「名册池号与流截断三数未画」缺口**本波闭合**；
- AC#6 框仍不由本腿翻：等非实现者按两份证据表现跑复认后并裁。
