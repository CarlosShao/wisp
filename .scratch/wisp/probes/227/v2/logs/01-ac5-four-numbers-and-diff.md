# 227-v2 读数件①＝AC#5 整包终态四数⇄逐名作差（对编排者 `A1053` 那发）

本腿现认锚＝起手 `53b8fc6b`（19:33:02）⇄ 发整包时 `d86bf31c`（19:35:03）。★尺＝`git diff --name-only 53b8fc6b..HEAD -- '*.go'`＝**0** ⇄ `7c44446a..HEAD -- '*.go'`＝**0** ⇒ 编排者那发与本腿两发的目标字节同一。

## 台规（三味都自证）

- 台面＝本腿自己 `git archive d86bf31c` 八件套＋仓根 `docs`＋手工三件（`build/*.dll` 3＋`third_party/sherpa-onnx/*.dll` 3＋`frontend/dist`）＝`go_files=775`／`dll_build=3`／`dll_third=3`／`dist_files=5`／`docs/PLAN.md=yes`（`scripts/01-setup.sh`、`04-setup-treeB.sh` 当场落的行）。★⛔ 复用编排者 `orch227-b` 那棵。
- GOTMPDIR＝`C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-b`＝**短＋与 %LOCALAPPDATA% 同卷＋反斜杠拼法**（先 `mkdir -p`）。
- PATH＝pwsh 原生整条 Windows 形＋**分号**（`$env:PATH.Split(';')`＝68 段、首段 `C:\Program Files\PowerShell\7`）⇒ ⛔ POSIX 冒号形混入＝`A1042` 定式那形没犯。
- 退码＝裸命令 `*>` 重定向到文件后立刻取 `$LASTEXITCODE`（件里逐字 `RC=1` 行），⛔ 管进 tail。

## 四数（三把尺⛔ 混用⛔ 相减）

| 发 | 树名 | RUN | PASS_TOP | FAIL_TOP | SKIP | rc |
|---|---|---|---|---|---|---|
| 本腿 full-A（弃用＝自造台损，见下） | `…\227v2\tree` | 2537 | 1706 | 27 | 13 | 1 |
| **本腿 full-B（决定性）** | `…\227v2\deliv-tree-ac5`（8.3＝`DELIV-~1` 现量） | **2552** | **1714** | **19** | **13** | **1** |
| 编排者 `A1053` | `…\orch227-b\tree`（8.3＝`ORCH22~1`） | 2317 | 1529 | 29 | 8 | 1 |

★full-A 那发的 8 枚 `internal/tools` 票 252 短名探针红＝**本腿自造台损**：我把导出树命名为 `tree`＝整条路径每段都≤8 字符、无需别名 ⇒ `GetShortPathNameW` 回显＝长名同串 ⇒ 那族前置"两形对比无从谈起"响。改树名后（full-B）同族只剩 1 枚（下）。

## ★编排者那发的名册是**截断名册**（这点最要紧）

- 尺①＝它那本 `full.log` 里 `cmd/wisp` 包块：`FAIL github.com/CarlosShao/wisp/cmd/wisp 176.237s` 那一行之前，最后一枚是 `=== RUN TestPanelThreadIsSTAAndExitsCleanly`、**该枚⛔ 任何结果行**＝整进程在包中途死亡（公共块颜色尺那条已知窗区形状，⛔ 计判据、但★计名册射程）。
- 尺②＝源码⇄实跑（`scripts/12/13`）：`cmd/wisp` 定义顶层用例 **356** 枚 ⇄ 编排者那发只跑到 **154** 枚 ⇄ 本腿 full-B 跑到 **333** 枚。块内 `=== RUN`：orch **208** ⇄ mine **448**。
- 尺③＝★**权威分母尺**＝本腿在导出树上跑 `go test ./cmd/wisp -list '.*'`（编译后的名册，⛔ 靠 grep 源码行）＝**333** 枚；★full-B 逐名比对本腿实际跑到的 `^=== RUN` 名集 ⇒ **差 0 枚＝333/333 全跑到＝终态坐实**。同一把尺打编排者那发 ⇒ 它那发只跑到 **154/333（差 179 枚）**＝名册截断的量化凭据。
- 尺④＝本腿那条 `grep -h '^func Test' cmd/wisp/*_test.go`（分母 **356**）⇄ 权威 333 的差＝**23 枚**逐名归属＝**19 枚 `//go:build windows && winlive`**（未带 tag＝⛔ 参与编译）＋**3 枚 `//go:build !windows`**（`secret_dataroot_119b_test.go`）＋**1 枚 `TestFixtureAlphaDrivesTheLeg`＝⛔ 真用例**（`dispatch_hop_selfproof_133_test.go:228` 那行落在反引号字符串字面量内＝票 133 自证台的夹具源码）。⇒ ⚠ 那把 grep 尺⛔ 感知 Go 词法＝分母虚高 1 枚，具名登记（`14-blob-eol.sh` C 段的反引号奇偶尺＝粗尺但指中了那一枚）；★AC#5 的"终态"判语⛔ 踩在这把坏分母上，踩的是 `-list` 那把。
⇒ 结论＝编排者那发那 29 枚⛔ 是"整包终态"的完整红名册；四数差（`2317`⇄`2552`＝235 枚 RUN）主体＝`cmd/wisp` 那包在它中途死掉。⛔ 用四数差当"改了 N 枚测试"，名集尺为准。

## 逐名作差（`scripts/05-namesets.sh` 落盘）

**A) 只在编排者侧红的 15 枚 ⇄ 本腿 full-B 逐名颜色**（全部 `--- PASS`，逐字行在 `07-attribute.sh` 输出里）：
`Test20r2JunctionRefusalIsNeverAskable`／`TestBridgeRefusesARealJunctionOnTheReadRoute`／`TestBridgeWritesNothingThroughARealJunction`／`TestSweepNeverDeletesThroughARealJunction`／`TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing`／`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`／`TestSyncAncestorLegNotExemptedStaysFailClosed`／`TestSyncAncestorLegReResolvesThroughJunctionToRealTree`／`TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116`＝**junction 族 9**；
`TestGateJudgesThePrivateSetByResolvedSID`／`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`＝**SID 卷形 2**；
`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`＝**跨卷 1**；`TestAC4InstallDirWidthIsDeliberateAndHasABreakLine`＝**安装目录 ACL 出带 1**；`TestAC2SealNoticeLandsInTheRunLegLogFile`＝**形状敏感 1**。
⇒ 这 15 枚与 `A1053`② 的 17 枚台损**逐名对得上**（17 ＝ 15 枚＋缺 `.git` 2 枚；那 2 枚两发同红，见 B 段），★本腿⛔ 一句"差不多"。

**B) 两发都在册（＝一致，⛔ 新增）**：`TestCleanCheckoutBuilds_AC11`＋`TestReadGitOnThisRepositoryIsSelfConsistent`（缺 `.git`＝台损，编排者归⑶）；CI 存量 8 枚逐名（`TestAC14GoSideEvalPushReachesThePage`／`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`／`TestComposerContractTypesMatchFrontend`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`）；`TestComposerRenderFixtureTellsTheTruth`（编排者归"导出树缺前端产物"，本腿拷了 `frontend/dist` **仍红**＝红句是缺 render fixture 的三块具名状态、`composer_test.go:855` 自己写"regenerate it with `npm run render:composer`"＝⛔ `frontend/dist` 那味能补＝**归类要更正**）；`TestResolvePerCallBudget`（争用型，本腿定向 3 发＝3/3 `--- PASS`，`targeted-budget-{1,2,3}.log`）。

**C) 只在本腿侧红的 5 枚 ⇄ 逐名定性**（⛔ 一枚是产品缺陷）：
1. `Test258BridgeMutationNoSrcKeepsOldBinding`／2. `Test258OccupiedCombinationNamesTheNewValue`／3. `Test258V1ProbeSummonEditRebindsLiveBall`＝红句逐字 `resident_hotkey_258_windows_test.go:361 the old binding did not survive a no-rebind host (VK=0x0): the premise moved` ⇄ `:453 live slots after the rebind = 1, want 2` ⇄ 产品自己那行 `hotkey occupied by another program, not registered … err="Hot key is already registered."`＝本机 `Ctrl+Alt+Q/M/P` 全被他程序占住。★安静窗定向 3 发＝**0P/4F ×3**＝稳定红、非摆动＝**桌面机器态**（编排者那发这 3 枚根本没跑到）。
4. `TestAC246DevLegIgnoresTheTestTaskInjection`＝红句里逐字 `wisp: another instance is running in this session; activated it; exiting`＋本腿现量盘上 `Get-Process wisp`＝**PID 21892 自 2026-10-15 19:03 活着**＝单实例互斥把这一发挡住了＝**桌面机器态**（⛔ 产品缺陷；本腿⛔ 动那枚进程）。
5. `TestTicket252R1AlignmentAddsNoAuthorization`＝红句 `paths_shortname_252_r1_test.go:166 前置条件缺失：D:\tmp\wisp\20261019\227v2 的短名回显…与长名同串`＝它取的是仓根的**父段**短名、我那一层目录名仍无别名＝**本腿台损**（编排者那发这枚在两侧名集里都出现＝同枚，⛔ 差集）。
ⓐ 另枚具名：`internal/winsec` 那族 342 落地件在整包⇄定向之间**摆动**（full-A 3 枚红、full-B 2 枚红；安静定向 3 发＝`0P/3F`→`1P/2F`→`1P/2F`，红句逐字 `ancestoropen_342_windows_test.go:48 … leg that answered = none: the call returned nil : 11 attempts` ⇄ `ancestorworld_342_348_windows_test.go:386 settle the world back to the real state: remove …`）＝编排者归的"落地续体 2"方向成立（342-v2 已判部分成立的 Windows 半区未完工作、⛔ 新增红）；第 3 枚 `TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1` 摆动、⛔ 定性为新增缺陷，具名交下程。

**D) 作差闭合（逐名⇄`A1053`②③的归类，⛔ 笼统"差不多"）** — 先把两发名集摆成一把尺：编排者 29 枚 ⇄ 本腿 19 枚 ⇒ **只在编排者侧红＝15 枚（本腿逐名全绿）**／**两发同红＝14 枚**／**只在本腿侧红＝5 枚**（15+14＝29、14+5＝19，两把都闭合）。
- 那 **15 枚**逐名＝`A1053`② 台损里的 junction 族 9（含 `Test20r2Junction…`）＋SID 卷形 2＋跨卷 1＋`AC4` 安装目录 1＋`AC2SealNotice` 形状敏感 1＋dirOf 1 ＝15 ⇒ **与本腿绿逐名相符**。
- 那 **14 枚同红**逐名＝CI 存量 8（上面 B 段列全）＋落地续体 2（342 那两枚）＋缺 `.git` 2＋`TestComposerRenderFixtureTellsTheTruth` 1＋`TestResolvePerCallBudget` 1 ⇒ 归类方向与 `A1053` 一致，★除 composer 那一枚（下面更正）。
- ★一处**归类要更正**（具名给编排者，⛔ 我改它台账）：`A1053`⑹ 把 `TestComposerRenderFixtureTellsTheTruth` 归"导出树缺前端产物"＝本腿拷了 `frontend/dist`（5 枚文件）在**两枚树**上都**仍红**、红句逐字 `the render fixture has no block for "no host attached" …` ⇄ `:902 render fixture verified across 0 painted states`，而 `composer_test.go:855` 自己写着 regenerate 靠 `npm run render:composer`（缺的产物＝render fixture、⛔ `dist`）⇒ 归类应改指"缺 render fixture 产物"＝**仍属台面形、⛔ 产品缺陷**，但补哪一味完全不同。
- ★一处**与 `A1053` 自洽性相反的事实**（派单说的"复跑撞到矛盾＝具名上报、⛔ 自行改判"）：`A1053` 末行那句闭合＝"CI 存量 8＋落地续体 2＋台损 17＋未定性 2＝29"，而 `window-facts.txt` 第 5-6 行那把台损细分**把 dirOf 计入 17**（`…＋dirOf 切在 gotmp 根上 1＋…`）、⑴ 又把 `20r2` 计入 junction 9 ⇒ ③那两枚"未定性"名字**已在 17 之内**＝同两枚被数了两次；两文件之间还差一枚（`window-facts` 有 dirOf⛔ AC2Seal、`A1053`② 有 ⑺AC2Seal⛔ dirOf）。⇒ 闭合的**枚数**（8＋2＋17＋2＝29）在两把细分尺下⛔ 同时成立；本腿⛔ 据此改任何一格结论，只具名上报请编排者自己重对那两枚的归属。**本腿的独立闭合是上面那把 15⇄14⇄5 尺**。
- ★一处**根因更正**（同形＝派单点名的"我⛔ 把你的读数当凭据"）：`A1053`⑵ 把 junction 那族归"长 D 卷 `GOTMPDIR` 让 `mklink` 建不出"。本腿现量＝**短 C 卷**根、只把拼法写成正斜杠，同一枚 `mklink /J` 照样 `Invalid switch - "Users"`（件②配对台）⇒ 根因＝**GOTMPDIR 的拼法**，⛔ 卷、⛔ 长短。归类方向不变（仍＝台损）。
- **只在本腿侧红的 5 枚**逐名定性：`Test258Bridge…`／`Test258Occupied…`／`Test258V1Probe…`＝本机 `Ctrl+Alt+Q/M/P` 被他程序占住（红句 `Hot key is already registered.`；安静窗定向 3 发＝`0P/4F`×3＝稳定、⛔ 摆动）；`TestAC246DevLegIgnoresTheTestTaskInjection`＝红句逐字 `wisp: another instance is running in this session; activated it; exiting`＋盘上现量 `wisp` PID **21892**（起于 2026-10-15 19:03、本腿⛔ 动它）＝单实例互斥挡住这一发；`TestTicket252R1AlignmentAddsNoAuthorization`＝本腿台损（它取仓根**父段**的短名、`D:\tmp\wisp\20261019\227v2` 那层仍无别名）。⇒ **⛔ 一枚产品缺陷**。
- ⇒ **终判（AC#5）**：⊖〔台损＋桌面机器态＋CI 存量＋落地续体〕后 **新增产品红＝零枚**；两发的差异全部可归到"编排者那发 `cmd/wisp` 中途死（名册截断）＋本腿两枚台面形状（别名父段／拼法）＋本机热键与单实例"。
