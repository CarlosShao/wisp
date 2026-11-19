# 167 片B 落地证据 — 忙时能停 / 失败原文单独存住 / 崩溃自救最小闭合（写腿 agent-143，2026-11-19）

- 票：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md` AC#4＋AC#5＋AC#6
- 派单：`.scratch/wisp/probes/orch/20261119/txt/167-b-stop-draft-rescue-dispatch.txt`（存档提交 `2c43ec586`，该提交只含派单存档本身，实现零行）
- 写面：`cmd/wisp/**`、`internal/panel/**`（composer/pump 新增字段一寸）、`docs/evidence/s1/`；`internal/agent/**` 零字节（既有取消语义全部复用，没有可写处）；`frontend/**` 零字节
- 硬约束核过：thresholds/golden/allowlist/PLAN/specs/D43/design 零字节；无断言放宽；未 push；零 emoji（d22scan 0 findings 现跑见 §5）；`git add`/`commit` 显式 pathspec

## 0. 先量后写 — 现量到的既有出口（报回，未造第二套）

| 枚 | 既有出口 | 凭据（全部现量，2026-11-19 本腿） |
|---|---|---|
| 停止·入向方法名 | `panel.MethodRosterCancel`＝"roster.cancel" 已在白名单与派发表 | `internal/panel/bridge.go:266`（值）、`:535`（白名单 case）、`internal/panel/composer_dispatch.go:493-497`（派发）——票 167 §9 行③写的"需第五枚入向方法名＝Q-76"已被票 38 AC#4（2026-11-17 readface wave 2）落地覆盖，`roster.cancel` 就是那枚方法名，Q-76 那一寸不再是本格缺口 |
| 停止·门＋缝 | `rosterCancelDoor`（late-bound seam）＋`cancelRosterRowWave2`：根任务走 `agent.RunningTask.Cancel`（`internal/agent/loop.go:341-342`），名册行走 `tools.TaskRoster.Cancel`（`internal/tools/task.go:569`） | `cmd/wisp/panel_roster_cancel_wave2.go:119-149,204-233`；装配绑定 `cmd/wisp/resident_windows.go:574`（bindPanelRosterCancel） |
| 停止·根任务句柄 | 根任务在飞时 `setCurrentRoot(bg)`、`Wait` 后清掉 | `cmd/wisp/run.go:2429/:2431`（现号 :2440/:2442 一带，本腿未动该段） |
| 忙态停止有效的牙 | 票 38 的两发已有：真子代理停到 D43 Muted（M16 假停变异被终态字钉红）＋真根任务忙时带位停到 exit=1 | `cmd/wisp/panel_roster_cancel_wave2_test.go:97-252`（:216 Muted 断言）、`:265-335`（:331 exit=1 断言）——本腿现跑复跑两发均 PASS（§5） |
| 草稿·入向承载位 | `ComposerRequest.Text`（tag `json:"text,omitempty"`）＋`panel.message.send` 名册/白名单/派发三处齐 | `internal/panel/bridge.go:98` 一带、`composer_dispatch.go` MethodMessageSend case |
| 草稿·发送门 | 453 门已在产：preview 转发＋纯消息经 OutgoingMessage.ForAgent → relay | `cmd/wisp/panel_message_453.go:121-208`；resident 中继绑定 `resident_windows.go:565`（bindPanelMessageRelay453）——票 167 §9 行④"生产装配 Message: nil"已被 453 落地覆盖 |
| 草稿·出向 | **零枚（量后才写的就是这一枚）** | 发送失败后拿不回原文——本腿 §2 补的就是这枚出向（composer 区新增两键） |
| 崩溃·现场记录 | panic sink 产/投/落三处齐：`defaultPanicSink` → slog → JSONL（tee：primary＝过脱敏流水线文件，mirror＝stderr） | `internal/observe/goroutine.go:167-177,236`；`cmd/wisp/logsink.go:215-223`——复用，未造第二记录面 |
| 崩溃·截断口径 | 盘上单串上限 512（rule 4），完整栈只在 stderr，GUI 腿无控制台 | `internal/observe/redact.go:36-37`（`MaxLoggedString = 512`）、`:141` |
| 崩溃·诊断包 | `BuildDiagnosticsBundle` 已有产码调用者（167-a2 时为 0，后被 180P/40 接上）：CLI `wisp diagnostics export`＋面板 diagnostics 门 | `cmd/wisp/diagnostics_cmd.go:286`、`cmd/wisp/panel_diagnostics_40.go:186`——自救指令只指向这两个既有出口 |
| 崩溃·doctor 读 logs | **零枚（量后才写的就是这一枚）** | `cmd/wisp/doctor.go` 28 枚检查项无一读 `logs\`（167-a2 §9 行⑦读数复认）——本腿 §3 补的就是这枚接线 |

## 1. AC#4 停止 — 牙名/读数

生产码**零改动**（出口全在，造第二套是被禁的）；本格交的是 167 自己那发牙（164 的牙钉名册行回填尾巴，`internal/tools/task_164_ac4_cancel_tail_test.go`；本发钉**根任务流尾巴**）：

- **牙** `Test167AC4BusyStopLandsAndLeavesNoWritingTail`（`cmd/wisp/panel_stop_tail_167_test.go`）：
  - **忙态下停止有效**：起真根任务（mock 300ms 延迟长回声），**等到根任务自己的流行已有 ≥1 条 delta 才发停**（忙是被测前验证过的事实，不是假设），经真 `ComposerDispatch` 门（`roster.cancel`＋确认位）→ 缝 → `RunningTask.Cancel`。读数：回执"已对任务 … 发出停止"、exit=1（取消类，假停会是 0）、task_log 落"cancelled"闭行（根名册行对无正文终态刻意不落 D43 维——backfill 空文本闸拒收，D43-Muted 读数属 164 的名册支，本格凭据用闭行＋exit）。
  - **停止后不留还在写的尾巴**：结算后（EvDone 已 Close）读 results 面（`rt.stream.Chunks()`），静置 600ms 再读——**一字不许长**。牙句逐字："停止后仍有 chunk 落地：settle 时 N 字，之后 M 字（+…）——取消之后不许有还在写的尾巴"。
  - **变异证 M1（红）**：`consoleSink` EvDone 分支埋 200ms 延迟续写（"还在写"形状）⇒ 该牙红，读数逐字："停止后仍有 chunk 落地：settle 时 43 字，之后 54 字（+\"M1-tail-167\"）"（vehicle rc=1）。逐字改回，`git diff` 复核零残留。

## 2. AC#5 草稿 — 牙名/读数

- **生产码（Go 半边）**：
  - `cmd/wisp/panel_failedsend_167.go`：`failedSendLedger`——单槽、逐字、无第二上限（入向封套上限就是上限；截原文会让"找回你发的话"在最需要它的场景变成假话）。写入只发生在两枚真失败支：中继未接入、中继拒收（`panel_message_453.go` 两分支各一行 `record`，空白不记账）。
  - 出向：`ComposerState` 新增 `failedSendText`（omitempty）＋`failedSendKnown`（omitempty）两键——**不在**票 145 已批解冻范围之外的任何面；`PumpSources.FailedSend` 读者＋pump 填充；无读者 ⇒ 两键都不上线（reader-gated absence）。
  - **反覆盖规则的 Go 半边＝分离本身**：失败原文走自己的键，失败路径不碰 composer 区任何其他字段（不发"整块替换 composer 内容"的指令——这条线上不存在这样的指令）；成功不抹槽（正控）；新失败**替换**为自己的原文、**绝不并叠**。页面半边（恢复的呈现与消费）归界面腿。
  - 接线三处：run.go pump literal（`FailedSend: rt.failedSendReading`）、resident boot 绑定（`bindFailedSendLedger453(rp, src)`，resident_windows.go）、降档泵（`residentPumpSources65` 从链上 message 门直读）。
  - 契约轴：新键在 `internal/panel/composer_test.go` `composerTSGapLedger` 注册（judged page-side gap，判语在册），两把 reconciliation 尺全绿。
- **牙**（`cmd/wisp/panel_failedsend_167_test.go`，四发全 PASS）：
  - `Test167AC5FailedSendKeepsItsOwnWordsSeparately`：失败原文逐字入槽（含换行）；**正控：成功后槽里仍是原文甲（"不失败时草稿不被抹"）**；新失败槽＝"失败丙"整串（不许并叠）；空白失败不挪槽。
  - `Test167AC5FailureTouchesNoOtherComposerField`："失败时整块回填"钉子——失败（与成功）前后 composer 区**除两枚新键外逐键 byte 相等**，回填进任何别的字段都红。
  - `Test167AC5WireShapeIsReaderGated`：无读者两键缺席；有读者恒发 flag；无失败不发 text（空恢复邀约＝撒谎键族）。
  - `Test167AC5FailedSendFacesAreWiredAtBothPumps`：两台泵的接线钉（拔线即红）。
  - **变异证 M2（红）**：成功分支抹槽 ⇒ 正控行红，读数逐字："成功发送抹掉了草稿：slot=\"\" want=\"原文甲：这一段是用户先写的字\\n第二行也在此\"（不失败时草稿不被抹）"。**变异证 M2b（红）**：`record` 改并叠 ⇒ 槽读成"原文甲…失败丙"，无并叠行红（"want exactly \"失败丙\" … no merge with the older draft"）。逐字改回，`git diff` 复核零残留。

## 3. AC#6 崩溃自救 — 记一笔可复制现场＋一条明确自救指令；明写未裁＋最小闭合集合

- **生产码**：`cmd/wisp/crash_selfrescue_167.go`——`doctorCrashScenes`（扫 `logs\wisp-*.jsonl`，按 mtime 新到旧、命中即停、上限 64 文件、逐行容错）＋`crashSelfRescueLine`（**可复制现场**＝文件路径＋记录自身 at/recovered/goroutine＋诚实截断注"盘上栈最多 512 字（日志管线截断规则，完整现场的落盘口径未裁）"；**一条自救指令**＝"把该文件里这一行 goroutine panic recovered 整行复制出来交给维护者；更完整的现场用 wisp diagnostics export 导出诊断包（config.toml [privacy] diagnostics_opt_in 关着时先改为 true 再跑）"）。接线：`cmd/wisp/doctor.go` 新增 INFO 检查项"crash scene (logs)"（过去式崩溃不是今天红闸，INFO 非 critical；数据根不可读也明说）。
- **复用，未造第二记录面**：现场本体＝既有 panic sink 落的 JSONL 行；指令指向的两个出口（复制该行 / diagnostics export）都是既有出口。
- **明写未裁（这一格的未裁清单，⛔ 本腿一件都没做）**：
  1. 盘上栈只有栈头 512 字（`MaxLoggedString`，rule 4）——放宽它＝改脱敏规矩（甲形），未裁；
  2. 另开一枚不走脱敏管线的崩溃现场文件（乙形，要先答"栈里会不会带密钥"）——未裁，与票 264 同轴；
  3. "禁用插件后重启"——插件地基票 50/51 没建，**指令里不写、不承诺**（有牙钉它的缺席：自救行含"禁用插件"⇒红）。
- **最小闭合集合（本格已闭合的全部）**＝既有 sink ＋ doctor 把既有记录变成可找到、可复制的 ＋ 一条只指向既有出口的明确指令。
- **牙**（`cmd/wisp/crash_selfrescue_167_test.go`，四发全 PASS）：找到现场并给出含路径/复制/diagnostics export/opt_in 的指令行（"禁用插件"缺席钉）、空目录与缺目录两形安静行必带所查路径、新旧文件最新者胜、doctor 接线钉（拔了 `doctorCrashScenes(dir)` 调用即红——"话说出口没接线"那枚普查缺口从此有尺）。
  - **变异证 M3（红）**：扫描键改成不匹配的 msg ⇒ "the planted panic record … was not found: the scene the sink already wrote stays unreachable"（两发红，vehicle rc=1）。逐字改回，`git diff` 复核零残留。

## 4. 边界（本腿没做、也不许做的）

- 页面侧渲染/恢复/消费（草稿槽怎么展示、何时清）＝界面腿；`frontend/**` 零字节。
- `internal/agent/**` 零字节：忙停、取消、流尾全部走既有语义（RunningTask.Cancel / loop 的 StatusCancelled / terminalWriteCtx 取证写），没有可写处。
- "面板真能按下去"若还需要新入向方法名（C17 白名单变动）＝契约面，另走批准，本腿零枚新增方法。

## 5. 三样现跑读数（全在本腿会话内现跑）

1. **d22scan**：`./scripts/d22scan.sh` ⇒ **0 findings**（写证据件前跑过 1 finding＝doctor 注释引用了本件路径而本件未建，本件落盘后复跑归零；末行："d22scan: 0 finding(s)"，rc=0）。
2. **载具**（`scripts/wisp-cli-staged-tests.sh`，串行一次一发，跑前 `Get-Process wisp.test` 无在飞）：
   - `Test167…` 11 发全 PASS（AC#4 一发＋AC#5 四发＋AC#6 四发＋AC#2/AC#3 既有两发复跑）⇒ **vehicle rc=0**；
   - 邻族回归 `TestRosterCancelDoor|Test453|Test192MessageDoor|Test214|Test38ForceChip|TestTaskCardFields|Test65ResidentPump|TestThePumpIsDriven|TestAC8SnapshotClamp` 全 PASS ⇒ rc=0；
   - 变异三次均按预期红（rc=1，读数见 §1/§2/§3），改回后复跑绿。
3. **`go build ./...`** ⇒ **rc=0**。
4. 附属：`go test ./internal/panel/` ⇒ ok 94.6s rc=0（reconciliation 两把尺带新注册键全绿）；新文件 gofmt 归零（run.go/composer_test.go 各有**既存**未格式化行，非本腿引入，未动）。
