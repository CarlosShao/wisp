# 400-r1 交付件：票 400 产品级红三件（panel 契约对齐／35r8 流日志有界／winsec 342r1 竞态修）

- 腿＝`400-r1`（写腿·实现者）；开工锚＝`d26a29083`（dev，索引零已跟踪改动、写面无他人脏）；交件时 HEAD 漂到 `0b6a6b1df`（398-r1 两笔，与本腿写面零相撞，`git diff d26a29083..HEAD -- internal/panel internal/winsec frontend/src/lib/panel.ts`＝仅 398 的 evidence 件）。
- 写面（实落 9 枚）＝`frontend/src/lib/panel.ts`＋`internal/panel/{pump.go,subagent_roster_197.go,backpressure_bound_35r8_test.go}`＋`internal/winsec/{winsec.go,pinnedleaf_342_windows.go,pinnedleaf_342_unix.go,pinnedleaf_342_fallback.go,ancestorunlink_348_windows_test.go}`。⛔ 未动：C17 桥（bridge.go 0 字节）、frontend 模板/样式、thresholds/golden、D43、risk/session/speech/audio、票面/台账/PLAN/specs。

## AC#1 panel 契约对齐（红句逐字 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`）

- **修法＝PanelSnapshot 声明补齐**（派单选项一）。依据（三腿照录后本腿核死的两点）：①C17 桥无第二条任务状态出站通道——`knownComposerMethod`（bridge.go:154，7 方法）无任何任务入口，`publishPanelSnapshot` 是唯一快照出站面（30-w-v1 实锤），且 C17 契约文本自带「任务状态」推送（PLAN.md:1373）⇒ **"Go 侧收敛"（让 Go 少发键）被契约本身否决**；②missing-features-2026-09-28-v3 §3.3 逐字「**转绿只有一跳**：界面那侧 `frontend/src/lib/panel.ts` 的 `PanelSnapshot` 补 `tasks?:` …**不许**为过尺让 Go 少发一个键、也不许改尺或加豁免」＋台账 A273（owner「批了」给 PanelSnapshot 新增字段）。
- **落形（纯类型增量，零行为面）**：`PanelSnapshot`＋`instructions?: InstructionsSection`／`tasks?: TaskRosterSection`（与 Go 侧 omitempty 指针同构）；`ComposerState`＋`git?/currentModel?/modelKnown?/credentialState?/credentialKnown?`（Go 侧自 181/187/248 起已在发、TS 名册一直未跟——同一病、同一红枚 TestComposerContractTypesMatchFrontend 的第二句红句一起治）；四节真形状接口镜像（InstructionsSection/ProjectInstructionFile/TaskRowView/TaskRosterSection/GitView/GitWorktree/CredentialState）。
- **读数**：`TestApprovalCardViewJSONKeysMatchFrontendTypes` 与 `TestComposerContractTypesMatchFrontend` 全绿，逐对 reconciled：`Snapshot<->PanelSnapshot: 6`、`ComposerState<->ComposerState: 11`、`ModeView: 3`、`WorkspaceView: 6`、`AttachmentRef: 9`、`ResultChunk: 3`。前端 `tsc -b` rc=0、oxlint rc=0（35 warnings 0 errors，存量）。
- **⚠ 具名判断（编排骨）**：写面行写 `internal/panel/**＋internal/winsec/**`，而选项一的落点在 `frontend/src/lib/panel.ts`。本腿按派单 AC#1 的第一选项执行（理由如上：唯一诚实修法＋A273 预批＋纯类型增量＋不碰模板/样式/桥/契约文本）；审查如不认，revert 这一枚文件即回红、本腿其余两件不受影响。

## AC#2 35r8 流日志有界（4x flood→4.0x 驻留）

- **病灶（读码定案）**：`Append` 只在 `s.truncated`（键数越界）时 clamp——单键 flood（键数≤32）时 `cur.Text += text` 无界；35-r8 的两枚**刻意红交付**自证（3 流×400K/1.6M runes 全量钉住）。
- **修法＝对照 StreamLog 既有有界先例原样下探一层**：①**行级体积界限** `streamRowKeepRunes=4096`（head/tail 各半，键数之下也 clamp，elision 计数入同一本账；截断档收紧回 256 的既有行为逐字不动，`TestTruncationFactsRideThePacket` 的 600-rune 行 elided=88 读数不变）；②**具名台账有界** `streamDroppedNamesCap=64`（Ring 保最近名，`DroppedCount()` 总数 accessor 补诚实余项——3/46 名小案全容纳，4000/16000 名大案有界）；wire 注释（TaskRosterSection.DroppedStreamKeys）据实更正。
- **两枚红交付按其自写转绿条件重切成绿尺**（名号世系保留：`…IsNotBounded35r8`→`…IsBounded35r8`，注释记明来历；红腿明令"⛔ 缩 flood 凑绿"全程遵守——flood 形状 3×200/800×2000 与 4000/16000×1×1 一字未动）。
- **读数（落地树）**：panel 全包绿；fan-out 读数 `1x text=35072 rows=64 names=64 heap+410872B | 4x text=35072 rows=64 names=64 heap+411112B | cap text=36864 rows=64`——**留存与 flood 规模无关**；below-bound 绿尺读数 `1x text=24690 | 4x text=24693 | cap 36864`（差 3＝marker 位数）＋行级 `TruncationFor` 每行 elided 恰＝该流 flood−8192。

## AC#3 winsec 342r1 竞态修（TestAC1..342r1 祖先换位窗 2/400）

- **复现（本机，开工首跑）**：`junction truncated=1/400`＋`symlink-dir truncated=1/400`＝恰 ci-r1 的 2/400；8 发连跑再证（1/1/2/1/1 命中率量级稳定）。损伤 body 逐字＝`content "raced private payload 342r1 windows ancestor"`（44 字节产码载荷）＝真产码漏（与 342-v2 裁定同形）。
- **漏窗（代码序证明）**：`risk.Resolve` 管线序＝READ① `reparseComponents`（reparse 判定）**先**于 READ② `resolveHandle`（开全路径取 GetFinalPathNameByHandle＝答案拼写）；换位落两读之间 ⇒ 判定"干净"、答案穿链成**外来拼写** ⇒ 路由对答案树全链自洽钉住/开封/写入，每道门都"通过"。封链的 per-component 身份门（pin-then-verify）在其输入本身被污染时结构上无牙。
- **修法（封链语义内，⛔ 不动 risk/C26/入口签名/旗标/封链形状）**：`privateFile` 传调用方原拼写——`pinnedParent(path, resolved)`；Windows 侧 walk 改下探 **caller 的组件**（`pathComponents(caller)`），答案 `resolved` 保留两职＝traversal verdict＋全部拒答文本锚＋体积根来源（`pinVolumeRoot(resolved)` 不变）；caller 非绝对拼写具名拒（RemoveUnlinked 先例同理由）。POSIX 侧 `resolved==path`（floor/C26 均直通不改写）⇒ 走 `resolved` 逐字不变，零行为变化（签名并列、注释具名"taken and ignored"）。未许诺未做到的照旧具名：叶子腿 322 残差、seal-by-name（票 120 Windows 腿）、FILE_SHARE_DELETE 顶名窗（342-v3 新票族）。
- **成对读数**：修前＝上（2/400，liveness 在场）；**修后 8 发×2 形×400＝6400 attempts，truncated=0，liveness 逐发在场**（provedReparse 12–75/发、formBLegRefusals 4–42/发）＝票 342 AC#1 口径两半齐（修前分子>0、修后 0 且 liveness 非零）。

## AC#4-6 反向判据三发（各对应红复现＋还原核）

| 发 | 突变 | 手法 | 红（照照） | 还原核 |
|---|---|---|---|---|
| M1 | panel.ts 删 `instructions?`/`tasks?` 两键 | 就地改＋还原（运行时读盘文件，-overlay 不达） | 红句**逐字**＝`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`（ci-r1 原句） | SHA256 before=after=`901D0106…080F89`；还原后绿 |
| M2 | pump.go 回退行级 clamp（`if s.truncated` 守回） | 仓外副本＋`-overlay` | `TestStreamLogFloodBelowKeyBoundIsBounded35r8` 红：`pinned text 4800000 runes > the cap 36864; live heap + 4827576 bytes`＋`retention scaled with the flood: 1x 1200000, 4x 4800000` | overlay 突变不触盘（工作树 0 状态变化）；去 overlay 全绿 |
| M3 | winsec.go `pinnedParent(resolved, resolved)`（走回答案拼写＝修前形） | 仓外副本＋`-overlay` | `TestAC1..342r1` 红：`truncated=4/400`（run1 symlink-dir）＋`truncated=2/400`×2 形（run2）＝修前窗原样回归，liveness 同发在场 | 同上 |
- 正控另有两枚在场：`TestFloodRulerFiresOnAnUnboundedSink35r8`（绿尺判语区分有界/无界 sink）与 `TestPlantedRendererDoorShapesGoRed` 族（契约尺正控）。HEAD 导出树配对正控＝同环境同命令打 HEAD 树：panel 4 枚红全在（本腿修的三件＋既有外部面 2 枚）、winsec TestAC1 红 ⇄ 落地树全绿＝修的正是这三件、非环境漂移。

## AC#7 零回归＋门禁

- **四数（落地树，`^=== RUN`/`^--- PASS`/`^--- FAIL`/`^--- SKIP` 口径）**：`internal/panel` RUN=271 PASS_TOP=173 FAIL_TOP=0 SKIP=0；`internal/winsec` RUN=127 PASS_TOP=76 FAIL_TOP=0 SKIP=0；`internal/tools` ok 16.79s 全绿；**`cmd/wisp` RUN=574 PASS_TOP=391 FAIL_TOP=41 SKIP=2（556.9s）——41 枚红全部落在音频/语音/窗测族**（15e/41a2r1/26P/392/381r2/390r1/26R3/382WR2/WireR1/12r2/16r1/33r5/AC14/256），红因逐字＝`speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)`＋`WISP_SPIKE_MODELS is not set`＝本 shell 的 DLL 阴影/语音台件环境形；**HEAD 导出树抽四枚代表（33r5/15e/26P/12r2）同命令同红＝环境形坐实、零枚可归因本腿**；本腿改动的真实消费面（cmd/wisp 的 panel/roster/streamlog/instructions 族：Test197×13／Test200r2／Test220／Test328/329／panel_pump 族）＝**PASS=17 FAIL=0 SKIP=0**。
- **d22scan.sh rc=0 clean**（分母与派单起点对上：bans#1-5 internal/=287、cmd/=52、ban#7 tools=25、ban#8 internal/=732〔起点 730＋本腿新枚〕、cmd/=183）；**d22scan self-test 全绿**（TestScanDetectsAllSeededViolations…TestBan8 系列）；**gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`）panel+winsec 0 命中 rc=0**；**path-length budget VERDICT GREEN**（62=62）；`go vet ./internal/panel ./internal/winsec ./cmd/wisp` 各 rc=0；前端 `tsc -b` rc=0、oxlint rc=0。
- 环境注：本 shell 跑 `go test ./cmd/wisp` 首发报 `exit status 0xc0000135`（测试二进制加载期 DLL not found——二进制依赖仓根的 `onnxruntime.dll`/`sherpa-onnx-*-api.dll`，`go test` 以包目录为 cwd 而 cwd/PATH 均无）；**HEAD 导出树同命令同败＝环境形**，非本腿引入；PATH 前置仓根后二进制可加载、全包照跑，语音族因 DLL 阴影（System32 的旧版 onnxruntime 优先于 PATH）照红——该形在 HEAD 树同现，与本腿零因果。
- census 认领面（ci-r1 刚认领 present/watchdog/aecspike）零触碰：本腿 `git status` 生产改动面仅上列 9 枚，无一枚在认领面。

## 纪律格

- ⛔ 翻框：票面/台账/PLAN/specs/阈值/golden/testdata/scripts/.github/C17 桥/D43——一字未动；⛔ push（本腿只 commit）；`add`/`commit` 显式 pathspec 逐枚手写；临时件只建不删（`D:\tmp\wisp\20261101\400r1\**`：突变副本/overlay.json/导出树/logs）。
- follow-up（⛔ 本腿立票，报编排者）：①`TaskRosterState`/wire 若要携带 `DroppedCount` 需动 cmd/wisp 的装配点（panel_pump.go:250），本腿写面不含 cmd/wisp，未动；②342 票 AC#1/AC#4 两格仍开（Windows 半区的终态判读归终裁腿）——本腿交付的是修后成对读数与反向判据，不勾格。
