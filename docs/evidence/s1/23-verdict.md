# 票 23 — 非实现者裁决表（system/window/input 工具族 · 遗产腿 e3abe0ebb）

- **角色**：本表出自**非实现者**（裁决者≠实现者，SPEC-12 §4.3 #1/#3）。实现证据＝
  `docs/evidence/s1/23-legacy-ac2-tests-ac34-residue.md`；实现提交＝ `e3abe0ebb`（四枚文件全为测试＋该证据件，零产码）。
  AC#2 实现基础＝ `a67e1d335`（loop+llm 半）＋ `111bfdb04`（tools 半：`Result.Images`＋`bridgeImages`＋system.go 文案）；两者均为 HEAD 祖先（`git merge-base --is-ancestor` rc=0，本表现跑）。
- **裁决锚点**：HEAD `7393cc604`（≥ 任务书给的 `1cc67664b`；`1cc67664b..7393cc604` 仅台账一笔 7 行，零 Go 面漂移——`git diff --name-only` 现量，无 `internal/tools`、无 `internal/agent`、无 `internal/llm`）。
- **隔离载体**：`git clone --shared` 工作克隆 `D:\work\workspace\projects plans\wisp-t23-verdict-leg`（仓外、不在仓内建 worktree，规则 21）。主树 `internal/tools` 未跟踪在飞件＝**0 枚**（孤儿件已随 e3abe0ebb 入库），`cmd/wisp` 的 323 未跟踪件不进克隆、与本裁决的包无关。
- **本表零写面**：未编辑票面、未编辑任何产码/测试、未碰 `frontend/**`、`design/**`、`docs/reports/pending-and-issues.md`；未设置 `WISP_T23_LIVE_CLIPBOARD`；零 GUI、零 wisp.exe、零 *winlive*、零 push。

## §1 起手三核（现跑）

| 核 | 读数 |
|---|---|
| 票面全文＋两份证据件读毕 | 五条 AC、Out of scope、进度日志（2026-11-03 残余点名）、遗产腿证据 §5 台账 |
| 实现提交在锚点上 | `git show --stat e3abe0ebb`：四文件 690 行新增；`111bfdb04`/`a67e1d335` 均祖先 |
| 在飞占用 | `git status --porcelain -- internal/tools` ＝空（克隆内亦空）⇒ 共享克隆内跑读数成立 |

## §2 门禁四把（共享克隆内全量现跑，fresh）

| 把 | 命令（逐字） | 读数 | 判 |
|---|---|---|---|
| ① | `gofumpt -l internal/tools`（`/d/work/base/gopath/bin/gofumpt.exe`） | **空输出**，rc=0 | ✅ |
| ①b | `gofmt -l internal/tools`（任务书字面那把） | **空输出**，rc=0 | ✅ |
| ①c | `gofmt -l .`（全树） | 15 行，**全部**落在 `.scratch/wisp/probes/**` 台件（历史已知，票 236 AC#6 族，`internal/`、`cmd/`、`tools/` 零命中） | ✅（域内） |
| ② | `go vet ./internal/tools` | 无输出，rc=0 | ✅ |
| ③ | `go test ./internal/tools -count=2 -timeout 20m` | `ok github.com/CarlosShao/wisp/internal/tools 67.801s`，rc=0 | ✅ |
| ④ | `sh scripts/d22scan.sh` | 正控 `PASS=45 FAIL=0 SKIP=0`；主扫 `d22scan: clean - no D22 ban violations`，rc=0 | ✅ |

环境读数：`CGO_ENABLED=1`（显式设 1，绝无 0）、`go 1.27.1 windows/amd64`、`GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset。克隆不带 `build/`、`third_party/`（untracked）⇒ 实测不需要：③④两发均在无 DLL 的克隆内全绿（与遗产腿 §3 读数同形）。

## §3 票 23 点名腿（克隆内 fresh `-run`，默认环境）

十枚点名（九枚 `--- PASS` ＋ 一枚 `--- SKIP`，rc=0）：

```
--- PASS  Test23ScreenCaptureCarriesPixelsOnTheImageChannel
--- PASS  Test23ScreenCaptureImageCrossesTheBridge
--- PASS  Test23CapturePixelsNeverReachRowOrLog
--- PASS  Test23NonCaptureResultsCarryNoPhantomImage
--- PASS  Test23CaptureAltNamesDRMBlackFrame
--- PASS  Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte
--- PASS  Test23UIPIPredicateLetsLiveNonElevatedTargetThrough
--- PASS  Test23ProcessNameIsTheD45BindingKey
--- PASS  Test23ClipboardReadSourceFeedsR4ThroughWrite
--- SKIP  Test23LiveClipboardRoundTripThroughRealOS   ← WISP_T23_LIVE_CLIPBOARD 未设＝默认跳过
```

加跑 loop 侧两枚（克隆内 fresh）：`TestScreenCaptureImageReachesMockLLM23`、`TestImageBytesRefNeverPersistedInJournal23`（`internal/agent`）＝ **2/2 PASS**。

## §4 AC#2 五腿逐一核对（都真）

| 腿 | 测试（盘上可复跑） | 独立核对读数 | 判 |
|---|---|---|---|
| 1. 像素→`Result.Images` | `Test23ScreenCaptureCarriesPixelsOnTheImageChannel` | 生产链在码上：`renderSystemResult`（system.go:273）`len(res.Bytes)>0` 分支填 `out.Images`，Text 只带 `{mime,bytes_len,drm}` 元数据且断言零像素 | ✅ 真 |
| 2. 桥 `ToolOutcome.Images` | `Test23ScreenCaptureImageCrossesTheBridge` | `bridgeImages`（bridge.go:704）把 `Result.Images` 逐张映射进 `agent.ImageAttachment{MimeType,Data,Alt}`，`out.Images: bridgeImages(res.Images)`（bridge.go:661）——无过滤无静默丢 | ✅ 真 |
| 3. mock LLM 可见 | `TestScreenCaptureImageReachesMockLLM23`（agent 包） | 历史带 `llm.ImagePart`（`BytesRef`＝`mem:` 句柄），第二发请求 wire body 逐字含 `data:image/png;base64,…` 真像素（data-URL 形），且 `BytesRef` 本身不上 wire、不进日志 | ✅ 真 |
| 4. tool_call 行＋审计四编码零泄漏 | `Test23CapturePixelsNeverReachRowOrLog` | 四编码＝raw/base64/hex/十进制字节串（`leakForms23`）；断言打在**桥真写的行**（`image23Journal.InsertToolCall` 记账副本）与真审计行上；正控先行（确认本次调用真带了像素，"没进库"不是测空气） | ✅ 真 |
| 5. 负控＋DRM 黑帧 | `Test23NonCaptureResultsCarryNoPhantomImage`、`Test23CaptureAltNamesDRMBlackFrame` | 负控双支：`system.get` 不凭空生图、后端零字节时不生空图（防"模型描述没收到的画面"）；DRM＝alt 与文本双点名黑帧（system.go:282 真文案） | ✅ 真 |

配套（同一证据链）：`TestImageBytesRefNeverPersistedInJournal23` 钉 tool_call 行零 `mem:` 句柄；`Image` 类型本身无 Path/Ref 字段（tool.go:71-91），"只有 mem 句柄进历史/DB/日志"由类型形状保证。

## §5 AC#4 剪贴板：门闩逐字核对（我未设环境变量）

1. **门闩**：`liveClipboard23()` ＝ `os.Getenv("WISP_T23_LIVE_CLIPBOARD") == "1"`（clipboard 测试 :44）；真往返腿第一行 `if !liveClipboard23() { t.Skip(...) }`（:87-89）。⇒ **默认门禁跑必 skip**——§3 现跑读数即为 `--- SKIP`（默认环境）。同包另一枚真 OS 腿族先例（`system_23_clipboard_memory_windows_test.go` 头注 :20-24）明写"绝不碰系统剪贴板"（纯内存 GlobalAlloc/RtlMoveMemory 方向腿）。
2. **默认跑不碰操作者剪贴板**：skip 分支在任何 OpenClipboard/EmptyClipboard 之前返回——**成立**。
3. **点名跑真绿**（凭据＝遗产腿 §3 操作者点名读数 `--- PASS (0.01s)`，marker 含中文＋纳秒戳→读回逐字一致→原内容恢复→恢复复读一致）。本腿按禁令未设 `WISP_T23_LIVE_CLIPBOARD`、未代跑该点名；默认门下它结构性不可触发（skip-first），故门禁绿不依赖它。
4. **taint→R4（评判半）**：`Test23ClipboardReadSourceFeedsR4ThroughWrite` PASS（§3 现跑）——以族自己的源名 `risk.SrcClipboardRead` Mark→`clipboard.write` 判 R4→`Channel==risk.ChClipboard`、`SrcTool==clipboard.read`；与既有 `Test23ClipboardWriteFactsFuseToL2`（SrcFSRead 形）互补成对。
5. **诚实边界**：`EmptyClipboard` 清全部格式——恢复逻辑只承诺文本（原文本→复原；非文本原状→空板，绝不伪造替代）。代码注释 :17-24 逐字写明。

## §6 其余三条 AC 核对

- **AC#1 风险矩阵**（2026-11-03 腿交付，本票范围）：`system_23_test.go` 全表断言——`SystemRowLevels()` 逐枚（:127-137：`system.get` L0、`system.set` L1、`system.power` L2、`input.type` L2…）；`Test23WindowManageCloseIsL2ByFacts`（close 经 R8 事实升 L2，:243）；`system.get` Origin＝焦点窗口标题（taint 出处，:469）；`Test23PowerAndClipboardWriteCarryR8Classes` 负控防"全钩子恒 L2"。fresh count=2 绿 ⇒ **PASS**。
- **AC#5 关机确认词**（同上）：`Test23PowerShutdownRequiresTheTypedWord`（:391）——lock/sleep/logoff 无词可用，restart/shutdown 必须带 typed 词，错词拒绝、未知 action 拒绝。fresh count=2 绿 ⇒ **PASS**。
- **AC#3 input.type**（后端余量腿，e3abe0ebb 交付）：UIPI 门在 `inputTarget` 于**任何字节发出之前**回 `ErrElevatedTarget`（seam 驱动，确定性；`sentText/sentKeys` 断言先拒后注＝零）；拒绝文本带「不会尝试提权绕过」承诺、无 runas/bypass/shellexecute 指引；真谓词 `realProcessElevated` 对本进程控制台窗口放行（GetConsoleWindow，零开窗；headless 走 pid 桩）；`realProcessName` 产出 D45 绑定键（小写镜像名、无目录分隔）＝将来 `approval_grant.pattern` 的形状。⇒ 后端可无 GUI 断言的余量**全部闭合且 fresh 绿**。

## §7 逐 AC 终裁

| AC | 票面判据 | 终裁 | 依据与残余 |
|---|---|---|---|
| **AC#1** 风险矩阵 | 全表逐枚断言 | **PASS** | §6；count=2 fresh 绿。无残余 |
| **AC#2** C7 ImagePart | 截图以 ImagePart 到 mock LLM；bytes_ref 不落盘/库/日志；DRM 黑帧＋文档 | **PASS** | §4 五腿全真＋fresh 绿；bytes_ref 内存 only 由类型形状（无 Path/Ref）＋两枚 DB/日志腿＋`assertNoPNGOnDisk23` 钉住；DRM 黑帧 alt＋文本双点名。**无本域残余**（「真屏截图→真 LLM」的实机演示被 *winlive*/GUI 禁令挡住，属真机签收类，非本票写面欠账） |
| **AC#3** input.type | 300 字符入 fixture 焦点窗；提权窗显式拒绝；确认卡三选项渲染（danger 样式） | **PARTIAL** | 已闭：L2 声明、UIPI 拒绝序＋冻结句＋绑定键（本票，fresh 绿）。**未闭（具名）**：① 真 fixture app 输入 300 字符 ② 真提权窗口实机拒绝 ③ 确认卡三选项 UI 渲染（`frontend/` 源内三档文案零命中＝未实现；`design/old/screens/approval.html:848-871` 是设计稿、非生产面板） |
| **AC#4** 剪贴板 | 往返＋taint 进 R4 | **PASS** | §5：taint→R4 fresh 绿；真往返腿门闩正确（默认 skip-first、不碰操作者剪贴板）＋操作者点名读数 PASS（凭据明示）。无残余（live 腿按禁令不可由 agent 代跑，门已备好） |
| **AC#5** 关机确认词 | fake executor 单测 | **PASS** | §6；fresh 绿。无残余 |

## §8 残余清单与归属（closeable 判断）

| # | 残余 | 归属 | 何时可闭 |
|---|---|---|---|
| R1 | 真 fixture app 注入 300 字符（开窗＋聚焦） | **frontend/owner-open-window**（真机集成；与 `wisp run` 实机验收同批做） | owner 允许开窗的窗口期，由可开窗腿代跑 |
| R2 | 真提权窗口实机 UIPI 拒绝 | **frontend/owner-open-window**（需起提权进程＝开窗类） | 同上 |
| R3 | 确认卡三选项渲染＋「所有应用」danger 样式＋二次点击展开 | **frontend**（owner 外部 agent 的写面；本编队零写面） | 源内三档文案落地＋PanelBridge 卡数据接通后 |
| R4 | 「真屏截图→真 LLM」实机演示 | owner 真机签收类（非写面欠账） | *winlive*/GUI 禁令解除时 |

**closeable 判断**：票 23 的**后端/工具可闭合面全部闭合**（AC#1/#2/#4/#5 PASS、AC#3 后端余量 fresh 绿）。**票本体不可整体结案**——AC#3 三枚残余（R1–R3）全部在 frontend/真机地界，本编队无写面也无法开窗验证。按本仓惯例（同型：票 50 授权 UI 残余），处置＝**票保持 open，残余逐枚具名归属 frontend/owner-open-window，等真机窗口期一次收口**；本裁决不翻任何票面勾，勾的翻动属编排者/票面流程。

## §9 独立性自证（对抗性核对）

- 我不采信实现腿自报的任何读数：§2/§3/§5 的每一发都是本腿在共享克隆内现跑（命令逐字、读数逐行）。
- 孤儿件三方字节同：SHA256 `82C0E007D1516F5E971765C153C8E1ABA808A95F973CBEB29CF4F3B6558D07DF` / 12261 字节，本表现跑复验＝遗产腿 §4 声明值。
- 挑过刺的地方（诚实记录）：① 遗产腿 §3 把 `gofmt -l` 记作 `gofumpt -l` 同发——本表拆开两把各自现跑，均空；② 全树 `gofmt -l` 非空（15 枚 `.scratch` 台件）——域内（`internal/tools`）为空，CI 的 tracked 分母问题已有票 236 AC#6 在治，不记在本票账上；③ `cmd/wisp` 在主树有 323 枚未跟踪件——不进共享克隆、不影响本票任何读数。
- 撤销口令：本表＝单一证据文件新增，`git revert <本提交hash>` 一枚即整体复原；对实现腿 `e3abe0ebb` 的复原仍是其自己那一枚 revert。
