# 票 272 · 腿 `272-v2`（非实现者验收腿）— AC#5 门禁四数 ＋ AC#7 单份 post-mark 欢迎语形：收尾判语

派单＝`.scratch/wisp/probes/orch/20261013/txt/272v2-dispatch.txt`。本腿面＝**只读＋重跑**：
⛔ 零产码/测试/断言改动、⛔ 未翻框、⛔ 未动 `Status:`/票名、⛔ 未 push；突变一律 `-overlay`（仓外 `D:/tmp/wisp272v2/`）。

## §0 起手锚（本腿自跑；锚件 `probes/272/v2/00-anchor.md`，提交 `7a42e168`）

- `date` 13:46:33；HEAD 起手 `11cfbc93`（recon 期间他腿交件推进 `eaf846a0→4cd874b7→11cfbc93`，其后 `78805129`/`c92f4dd5`/`f1676bfe`/`c2b00de4`——**均未触 `cmd/internal` 产码**）。
- 起手 `git status --porcelain -- cmd internal`＝**0 行**；被审尺 md5 `2ea71578ba9ca015c535673bc9792d09`＝HEAD blob、产码 md5 `348f5a67a73442d45f751079e3b80a6c`＝HEAD blob。
- 票 272 框现量：**2 未勾（AC#5／AC#7）／5 已勾**（⛔ 本腿不翻）。`gofumpt v0.12.0 (go1.27.1)`；GOTMPDIR 一律反斜杠形 `C:\tmp\wisp272v2gotmp`。

## §1 结论先行

- **AC#5：成立**（按派单口径：四数各有本腿现跑 rc；红名册每一枚都归因到在册族，无不可归因红；314 死区未复现）。
  字面「作差＝0」的落点：**cmd/wisp 整包 run2 与编排者 12:0x 基线同集**（唯一红 `TestAC2SealNotice…`）；
  **internal 两发对称差恰为在册间歇两枚**（`TestResolvePerCallBudget` ⇄ `TestSessionListByWorkspace`）。
- **AC#7：成立**（残留假绿**活着**，本腿自建突变复读双色）。单份 post-mark 欢迎语形在现行 HEAD 落地尺上 **PASS 4.22s**；
  零份对照 **FAIL 44.97s，红在 `:623` await**（⇒ 那枚绿不是空转）。**翻框建议：AC#7 可翻**，留档句见 §3。

## §2 AC#5 四数（逐条现跑；窗口逐发；rc 逐字）

| # | 命令要点 | rc | 窗口 | 读数（件在 `probes/272/v2/logs/`） |
|---|---|---|---|---|
| 1 | `PATH="$PWD/third_party/sherpa-onnx:$PATH" sh scripts/d22scan.sh` | **0** | 13:48:49→13:49:15 | 正控先绿 `runtests.sh: OK … top-level: PASS=35 FAIL=0 SKIP=0`；全扫 `clean`；`ban #8 cmd/ examined 125 Go files`、`internal/=550`（`G1-d22scan.txt`） |
| 2 | `gofmt -l`／`gofumpt v0.12.0 -l`，**只喂 `.go`** | **0** | 13:50:23→13:50:24 | 被审两枚（`config_reload.go`＋`config_reload_223_test.go`）**两把尺皆空**；包域存量（在册，非本票）：`gofmt -l cmd/wisp`＝`models.go`（CRLF 假枚）、`gofumpt -l cmd/wisp`＝`models.go`＋`panel_transport_35r2_test.go`（票 298 AC#5 账）（`G2-fmt.txt`） |
| 3 | `GOFLAGS= go vet ./cmd/wisp/` | **0** | 13:50:36→13:50:37 | 净（`G3-vet.txt`） |
| 4 | 名册两包 × 两发（见 §2 表） | 各 rc 见下 | 逐发 | internal `rc=1`（`G4a`/`G4g`）、cmd/wisp `rc=1`（`G4d`/`G4f`） |

### §2.1 `./internal/...` 两发成对（run1 13:50:59→13:53:22 @`78805129`；run2 14:19:27→14:21:23 @`c2b00de4`；rc=1 各发）

| 红名（包） | run1 | run2 | 归因 |
|---|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS`（ball） | R | R | 在册 C21 族（272-v1 基线即在） |
| `TestAC4InstallDirWidthIsDeliberateAndHasABreakLine`（models） | R | R | **环境形（GOTMPDIR=C:\tmp 族）★新名**：同名单跑 mandated rc=1／默认 `%TEMP%` rc=0（§2.3） |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes`（panel） | R | R | 在册 panel 契约漂移 |
| `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`（panel） | R | R | 在册 35r8（票 35 `:49` 有意保留） |
| `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（panel） | R | R | 同上 |
| `TestComposerContractTypesMatchFrontend`（panel） | R | R | 在册 panel 契约漂移 |
| `TestC21DesignTokensFourWayAgree`（panel） | R | R | 在册 C21 族 |
| `TestResolvePerCallBudget`（risk） | R | **G** | 在册计时带载（`A675/A676`：出现/消失都不计账）；非断言放宽对象 |
| `TestSessionListByWorkspace`（session） | **G** | R | 在册 session 间歇（CI `test-core` 6 之一，归 166 侧） |
| `TestGateJudgesThePrivateSetByResolvedSID`（winsec） | R | R | **环境形（GOTMPDIR 根）**：默认档 rc=0（§2.3）；`TempDir RemoveAll cleanup … Access is denied` |
| `TestSealNarrowsAndNamesThePrincipalItRemovedBySID`（winsec） | R | R | 同上 |

⇒ **对称差＝{`TestResolvePerCallBudget`, `TestSessionListByWorkspace`}**＝在册间歇对，其余 9 名两发同色。
起手基线对拉：272-v1（10-07@`e18e32da`）6 名红全部仍在册解释；新增的 4 名＝3 枚环境形＋2 枚 35r8 中……逐名如上，**无不可归因红**。

### §2.2 `./cmd/wisp/` 两发成对 ＋ 与编排者 12:0x 基线对拉

| 发 | 窗口 | HEAD 起→止 | 结果 | 红名 |
|---|---|---|---|---|
| run1（`G4d`） | 13:58:29→14:04:05 | `f1676bfe`→`c2b00de4` | `FAIL … 330.359s` rc=1 | `TestAC2SealNoticeLandsInTheRunLegLogFile (0.06s)`／`TestAC14GoSideEvalPushReachesThePage (0.57s)`／`TestBallPanelGesturesReachThePanelThread (19.55s)` |
| run2（`G4f`） | 14:07:25→14:12:16 | `c2b00de4`→`c2b00de4` | `FAIL … 286.119s` rc=1 | **仅** `TestAC2SealNotice… (0.06s)` |
| 编排者 12:0x（`A915`） | 296.812s／333.948s | — | rc=1 各 | **仅** `TestAC2SealNotice…`（逐字同） |

- ★**314 死区两发均未复现**（330.359s／286.119s 全程跑完、无整进程消失、无 157–215s 死形状）。
- 归因：`TestAC2SealNotice…`＝**环境形**（`GOTMPDIR=C:\tmp` 族；186-v1/A915 先例＋本腿 A/B §2.3：mandated rc=1／默认 rc=0）；
  `TestAC14GoSideEvalPushReachesThePage`＝**在册真窗族**（颜色随会话漂；台账 `:14768`/`:17057`/`:17354`/`:17382` 具名）；
  `TestBallPanelGesturesReachThePanelThread`＝**真窗族 ★新名**（CI run327/328 10-13 皆 `PASS 0.29s`；本腿单跑 B/E 皆 PASS 0.81/0.69s；
  整包 run1 红 19.55s＝`timed out after 15s waiting for the gesture's show request to end in a shown window - this is a failed measurement`，run2 未现）
  ⇒ 同一对整包内**翻色**＝该族定义行为（单发名册不稳定，判据取成对）。
- ★两枚**族在册、名未入册**的「新名」（`TestAC4InstallDirWidth…`、`TestBallPanelGestures…`）已具名上报；**无不可归因新名**，故未触发「停手」。
- 被审用例（`TestTicket223RestartTierSaysItWillNotApply`）在两发整包内均未入红名册；家族 `-run 'TestTicket223'`＝316-r1 今日 rc=0 PASS=10（新鲜读数在案）。

### §2.3 环境形 A/B（本腿自跑，逐单测）

| 单测 | mandated `C:\tmp\wisp272v2gotmp`（反斜杠形） | 默认 `%TEMP%`（unset） | 件 |
|---|---|---|---|
| `TestAC2SealNoticeLandsInTheRunLegLogFile` | **rc=1**（`records naming "winsec: seal cleared…" = 2, want 1`） | **rc=0** PASS 0.08s | `G4e` A/D |
| `winsec` test pair | **rc=1**（`TempDir RemoveAll cleanup … Access is denied`） | **rc=0** PASS 0.02/0.11s | `G4bc` A/B |
| `TestAC4InstallDirWidthIsDeliberateAndHasABreakLine` | **rc=1**（inherit 的 `BUILTIN\Users (I)(OI)(CI)(RX)` 留下） | **rc=0** PASS 1.52s | `G4bc` A/B |
| `TestBallPanelGesturesReachThePanelThread` | rc=0 PASS 0.81s | rc=0 PASS 0.69s | `G4e` B/E |
| `TestAC14GoSideEvalPushReachesThePage` | rc=1（title=""） | rc=1（title=""） | `G4e` C/F |

⇒ 三枚环境形红全部由 **GOTMPDIR 根**决定（默认档全绿）；真窗族两枚与 GOTMPDIR 无关（随会话/负载漂）。

## §3 AC#7 两读数（本腿**自建** overlay 突变；与在案读数对拉）

- 补丁基准＝现行 HEAD 产码（`348f5a67…`）：(a) 删 `reportRestartPending` 的真句（`config_reload.go:321-327` 的 `Fprintf`，两条 `auditf` 原样）＋(b) 在 `reportReload` 的 `state=applied` auditf（`:182-184`）之后注入一句与本次 plant 无因果的「欢迎语」打印（带五字＋全六 needle）。生成件仓外 `/d/tmp/wisp272v2/`（`mut-welcome`／`mut-zero`＋两枚 overlay，`AC7-gen.txt`）；两版 left-over 计数：welcome＝1（＝注入句）／zero＝0／HEAD＝1。
- **welcome（含注入）**：`--- PASS: TestTicket223RestartTierSaysItWillNotApply (4.22s)`，rc=0（14:17:39→14:17:52，`AC7-welcome.txt`）。
- **zero（只删不注入）**：`config_reload_223_test.go:623: stdout never carried "本次运行不会生效" AFTER the plant within 40s; window since the mark:` → `--- FAIL … (44.97s)`，rc=1（14:17:52→14:18:47，`AC7-zero.txt`）。
- **overlay 进编译自证**：两发只差注入一处删除相同 ⇒ 颜色差（红⇄绿）只能由注入造成——零份红＋含份绿＝双色互证，绿不是空转。
- ⇒ 与在案读数同向（272-v1 `PASS 2.13s`；编排者 `PASS 2.36s`＋零份 `FAIL 42.10s`@`:623`），**本腿读数非抄**。
- **判语：成立**（残留＝「单份 post-mark 欢迎语」仍假绿；arm1 只堵种前份、arm2 只堵窗内 ≥2 份）。
- **留档句（建议原样写进 AC#7 框内）**：「成立（272-v2 自建 overlay 复读：welcome PASS 4.22s、零份对照 FAIL 44.97s 红在 `:623` await）。残留＝单份 post-mark 欢迎语仍假绿：arm1 只堵种前份、arm2 只堵 ≥2 份；钉死需产码侧同流因果锚点（动文案＝本票禁区＋人工批准），本票不选形、不修——前向风险据实留档。」
  ⚠ 边界（不是缺陷、是形的定义）：本腿注入点仍走 `reportReload`（plant 路径上的一次重载）；此形语义＝「任何一次重载都会打」＝票面 `:12` 指的前向形，故绿是**假绿**的正确复现。

## §4 具名没做／交编排者

- **没做**：`-race` 档；`winlive` 档；AC#2 三形乙/丙**未重复自跑**（316-r1 今日已在现行 HEAD 现读，本腿未并跑）；`thresholds.go`/golden/`PLAN.md`/`specs/**` 一字未碰；`TestBallPanelGestures…` 的 15s 超时机制未深挖（归 314 面）。
- **环境纪律自报**：GOTMPDIR 全用反斜杠形；全部突变 `-overlay`；本腿对 `cmd/internal` **零字节**写入（§5）。
- ★**具名自曝（仓外笔误，两笔）**：本腿起手 `mkdir` 与两发短命令的 cwd 笔误（`…plans Wisp` 少一个 `/`）把空目录与两枚小件写进**仓外**同名邻近目录 `D:\work\workspace\projects plans Wisp`（**生于 12:20:08、非本腿创建**）：该目录下 `…/272/v2/logs/G4b-attr-mandated-gotmp.txt`（362B，误发次产物）与 `G4d-cmdwisp-suite.txt`（579B，go.mod 空转记录）。按「临时件只建不删」**不删、如实登记**；本腿全部真实读数件在仓内 `probes/272/v2/logs/`（`G4e` 文首亦保留了那次失败段，判语只引其中 `RETRY` 段）。
- **他腿在飞披露**：236-r7（`internal/tools`，我 internal run1 窗口内其 `M …teeth_test.go` 在场——该包仍 `ok 26.689s`）；317-r1（`scripts`）；253-r6（`cmd/wisp`：`?? panel_unbound_253r6_windows_test.go` 14:10:55 露面＝在我 cmd/wisp run2 编译**之后**、未进编译；`M panel_control_windows.go` 14:21:39 在我两发整包**之后**）。每发窗口/HEAD/status 已逐发落盘。

## §5 终态自证

- 被审两件 md5 终态复量＝起手：`348f5a67…`／`2ea71578…`（`git show HEAD:` 同值）；`git status --porcelain -- cmd internal` 中本腿零项（出现的 `??`/`M` 均为他腿在飞件、已具名）。
- 本腿只写：`probes/272/v2/**`、本证据件、票 272 尾 1 条 bullet；commit 显式 pathspec、只 commit 不 push。
