# 票 23 / 遗产腿收口 — AC#2 孤儿测试落地 + AC#4 真剪贴板往返 + AC#3 谓词/绑定键余量（实现腿交付记录）

- 票面：`.scratch/wisp/issues/23-system-window-input-tools.md`（Blocked by 20-host-bridge-fs-tools，已 done）
- 基线 HEAD：`aa7159b29`（派单写 ac7d14c7f，其后仅台账/证据追加 `aa7159b29`、`aa7159b29` 前的 wave44 收账 `A1588`，写面无漂移）
- 本腿写面：`internal/tools/**`（三枚 `*_test.go`）+ 本文件 —— 零产码改动
- 本文件是**实现证据**，不是裁决表：AC 的对抗验收须由非实现者另做（SPEC-12 §4.3）

## 1. 起手三核（动手前现跑）

| 核 | 问题 | 读数 |
|---|---|---|
| (a) 票面全文读完 | 五条 AC、Out of scope、进度日志 | 读毕。进度日志（2026-11-03）点名残余：AC#2 C7 接缝、AC#3 fixture/UIPI/确认卡、AC#4 真剪贴板往返 |
| (b) 点名行号是否漂移 | PLAN.md:2556（bytes_ref 内存 only）/ :3234（UIPI 句）/ :3236（DRM 黑帧行）/ D34 表头 | **现测**：2556=screen.capture 注「只存内存 bytes_ref，不落盘、不进日志、不进 DB」✓；3234=「目标程序以管理员身份运行，Wisp 无法向其输入」✓；3236=「对 DRM 保护窗口截到黑屏…明示即可，不修」✓；D34 表头在 2522 ✓ |
| (c) 有没有吃别人的在飞改动 | `git status --porcelain -- internal/tools` | **占用**：未跟踪 `internal/tools/ticket23_image_channel_test.go`（12261 字节，死腿遗产，其头注自述针对 `a67e1d335` 的接缝写）。⇒ 按规则 21 走 `git clone --shared`，落回主树时 SHA256 字节同（`82C0E007…D07DF`，见 §4） |

前提复核：`a67e1d335`（AC#2 loop+llm 半边）与 `111bfdb04`（AC#2 tools 半边实现：`Result.Images` 通道 + `bridgeImages` + system.go 文案改写）均为 HEAD 祖先（`git merge-base --is-ancestor` rc=0）——**AC#2 实现已全在盘上，孤儿测试正是量它的**。孤儿件逐形状核对：`Result.Images`（tool.go:50-56）、`bridgeImages`（bridge.go:704）、`agent.ImageAttachment{MimeType,Data,Alt}`（agent/tools.go:93）、`agent.Journal` 五方法（agent/journal.go:18）、`GateFuncs`/`AnswerAllow`（gate.go:126/94）、`memory.ToolCall{Tool,ArgsJSON,Outcome}`（memory/models.go:74）全部存在且形状一致。

## 2. 本腿落了什么（三枚测试文件 + 本文件，零产码）

1. **`internal/tools/ticket23_image_channel_test.go`（代提，字节同）** — 死腿留下的 AC#2 tools 半边 5 腿：像素进 `Result.Images`（bytes/mime/alt）且 Text 只带元数据；桥 `ToolOutcome.Images` 转发不丢；tool_call 行 + 审计日志四编码零像素泄漏；负控（system.get/空后端不凭空生图）；DRM 黑帧在 alt 与文本双点名。`a67e1d335` 的证据文件 §7 点名的三个洞，由 `111bfdb04` 补了实现、由本件钉住回归。
2. **`internal/tools/ticket23_clipboard_roundtrip_windows_test.go`（新）** — AC#4 残余行。真剪贴板往返腿 `Test23LiveClipboardRoundTripThroughRealOS` 走生产入口 `realClipboardWrite→realClipboardRead`，**门闩 `WISP_T23_LIVE_CLIPBOARD=1`**（EmptyClipboard 清全部格式，故默认门禁跑必 skip、恢复原内容，与 *winlive* 禁令同构）；纯内存腿 `Test23ClipboardReadSourceFeedsR4ThroughWrite` 用本族自己的源名 `risk.SrcClipboardRead` 落污点 → `clipboard.write` 判 R4、通道=ChClipboard（前腿用 SrcFSRead 打的同一断言，这腿补上剪贴板读源）。
3. **`internal/tools/ticket23_input_uipi_windows_test.go`（新）** — AC#3 后端余量，收窄到无 GUI 可诚实断言的部分：UIPI 门在 `inputTarget` 于**任何字节发出之前**回 `ErrElevatedTarget`（seam 驱动，runner 无关）；拒绝文本必须带「不会尝试提权绕过」承诺、不得出现 runas/bypass/shellexecute 指引；真谓词 `realProcessElevated` 对**本进程自己的控制台窗口**（GetConsoleWindow，无创建无聚焦）放行并解析出 hwnd+进程名（headless 形走 pid 桩，被问的仍是本进程）；`realProcessName` 产出 D45 绑定键（小写镜像名、无目录分隔）——确认卡「本会话内允许向此应用」将来写入 approval_grant.pattern 的就是它。

## 3. 盘上读数（全部现跑）

门禁（共享克隆内）：

```
gofumpt -l internal/tools        → 空（rc=0）
gofmt -l internal/tools          → 空（rc=0）
go vet ./internal/tools          → rc=0
go test ./internal/tools -count=2 -timeout 20m
                                 → ok  github.com/CarlosShao/wisp/internal/tools  57.554s（rc=0）
sh scripts/d22scan.sh            → clean - no D22 ban violations（rc=0）
```

点名 9 腿（`-run` 现跑，全部 PASS）：`Test23ScreenCaptureCarriesPixelsOnTheImageChannel` / `Test23ScreenCaptureImageCrossesTheBridge` / `Test23CapturePixelsNeverReachRowOrLog` / `Test23NonCaptureResultsCarryNoPhantomImage` / `Test23CaptureAltNamesDRMBlackFrame` / `Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte` / `Test23UIPIPredicateLetsLiveNonElevatedTargetThrough` / `Test23ProcessNameIsTheD45BindingKey` / `Test23ClipboardReadSourceFeedsR4ThroughWrite`。

AC#4 真剪贴板往返（操作者点名读数，`WISP_T23_LIVE_CLIPBOARD=1 go test -run Test23LiveClipboardRoundTripThroughRealOS -count=1`）：

```
--- PASS: Test23LiveClipboardRoundTripThroughRealOS (0.01s)
（写入含中文+纳秒戳的 marker → 读回逐字一致 → 原内容恢复 → 恢复读回一致）
```

开发期一次真红（诚实记录）：UIPI 腿初版把「本工具不会尝试提权绕过」这句**冻结文案里的不绕过承诺**误判成绕过提示 —— 断言改为「必须带承诺 + 只禁 runas/bypass/shellexecute 指引」，产品文本一字未动。

## 4. 落回主树（规则 21，字节同）

- 工作克隆：`D:\work\workspace\projects plans\wisp-t23-legacy-leg`（`git clone --shared`，仓外，不在仓内建 worktree）。
- 孤儿件落回 = `git add internal/tools/ticket23_image_channel_test.go`（显式 pathspec），落回前后 SHA256 均为 `82C0E007D1516F5E971765C153C8E1ABA808A95F973CBEB29CF4F3B6558D07DF`、12261 字节——与克隆内副本、与死腿原盘上件三方字节同。
- 其余三件（两枚新测试 + 本文件）以 `git diff` 补丁形式落到主树后同 commit 收口（显式 pathspec，无 add -A）。

## 5. AC 台账（本腿后的残余，供裁决者核对）

| AC | 状态 | 依据 |
|---|---|---|
| AC#1 风险矩阵 | 已闭（2026-11-03 腿） | system_23_test.go 全表断言 |
| AC#2 C7 ImagePart | **本腿闭** | 实现= `a67e1d335`（loop+llm 半）+ `111bfdb04`（tools 半）；回归= 孤儿 5 腿本腿代提落地，全绿 |
| AC#3 input.type | 部分闭 | 已闭：L2 声明、UIPI 拒绝序+文本、D45 绑定键、session 绑定键生产端。**未闭（写面外）**：真 fixture app 输入 300 字符、真提权窗口实机拒绝、确认卡三选项渲染（=前端/机主外部 agent；真机验证还要求开窗，本腿禁令不允许） |
| AC#4 剪贴板 | **本腿闭到可复跑** | taint→R4（源=clipboard.read）+真往返腿（点名跑 PASS）；默认门禁跑 hands-off |
| AC#5 关机确认词 | 已闭（2026-11-03 腿） | system_23_test.go |

零产码：本腿 diff 全为 `*_test.go` + evidence md；`cmd/wisp`（另有未跟踪件，未触碰）、`frontend/**`、`design/**`、`docs/reports/pending-and-issues.md` 全程零触碰；无 GUI、无 wisp.exe、无 *winlive* 跑。

撤销口令：`git revert <本腿hash>` 一枚即整体复原（四枚文件全是本腿新增，无既有行被改）。
