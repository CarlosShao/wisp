# ci-red-1 —— 今晚两枚 push 触发的 CI 红，逐名归因（只读腿台件）

- 本腿起笔时刻 `2026-09-29 19:56 +08`，交件时刻 `2026-09-29 20:1x +08`（`date` 现跑＝`2026-09-29 20:12:12 +0800`）。
- 本腿起手锚点 `git rev-parse --short HEAD` = **`29215a93`**；交件时 HEAD = **`f6186158`**（其间别的腿进了 6 枚，见 §6 那条"树不同"的警告）。
- 被验两发（都 `event=push`、`headBranch=dev`、workflow `ci`）：
  - **`36559498617`** ＝ sha **`a7993a9b`**（A449 那发 ledger），`2026-09-29T11:04:39Z` 起跑（＝本机 19:04 +08），job 终态 11:06–11:17Z 之间，整体 failure。
  - **`36556778053`** ＝ sha **`4d01e0b2`**（A448 那发 ledger），`2026-09-29T10:38:10Z` 起跑（＝本机 18:38 +08），failure。
- 对照发（都在 `4d01e0b2` 之前、远程 tip 还是 `f7478d37` 的时候）：
  - **`36500353477`** ＝ `ci` **schedule** 触发，sha `f7478d37`，`2026-09-28T23:53:41Z`，failure。⇒ 这是"今晚推送之前 CI 看到的最后一发全形状"。
  - **`36149256584`** ＝ push，sha `64858d68`，`2026-09-25T14:42:15Z`，failure。⇒ 更早的第二把锚。
  - **`36206542728`** ＝ push，sha `f7478d37`，`2026-09-26T00:54:21Z`，failure（只取步级结论，没抽名册）。
  - **`36094258734`** ＝ push，sha `cd87354`，`2026-09-25T04:23:42Z`，failure。⇒ 台账 `A237②` 那支名册的锚，本腿**直接下日志复算**过它（见 §5）。
- 取数手段＝只有 `gh run view`（`--json` / `--log-failed` / `--log --job`）＋ `gh api …/artifacts` ＋ 读仓里已有文件。**零编译、零 `go test`、零 `go build`、零 `go vet`**；没有写任何跟踪文件；没有 commit。
- 原始件都在本目录 `logs/` 下（只建不删）：
  - `run-36559498617-log-failed.txt` 1,723,889 B ／ `run-36556778053-log-failed.txt` 1,654,814 B
  - `run-36500353477-log-failed.txt` ／ `run-36149256584-log-failed.txt` ／ `run-36094258734-log-failed.txt`
  - `run-36559498617-slofull-job.txt` 29,704 B ／ `run-36556778053-slofull-job.txt` 33,003 B ／ `run-36559498617-slosmoke-job.txt` ／ `run-36559498617-testwindows-job.txt` 943,228 B
  - `roster-baseline-f7478d37.txt`（17 名）／ `sc-baseline.txt`（47 条）／ `sc-tonight.txt`（50 行含非 finding 行）

---

## 0. 先回答三件被点名的口径问题

### 0.1 那 11 枚名册：不成立。真值＝**顶层 23 枚**（含缩进子测试 25 枚）

两发（`36559498617` 与 `36556778053`）的 `--- FAIL` 名册**逐字相同**，我用两把互相独立的尺量到同一个数：

| 尺 | 读法 | `36559498617` | `36556778053` |
|---|---|---|---|
| 甲（缩进敏感） | `grep -a -E '^\s*--- FAIL\b'` 打在剥掉 `job⇥step⇥时间戳` 前缀之后的正文行上 | **25 行**＝23 顶层＋2 缩进子测试 | 同一批 25 行 |
| 乙（仪器自报，不依赖我的 grep） | 各步末尾那行 `portable-tests.sh: four numbers (all from -v output)` | test-core `RUN=1357 PASS=913 FAIL=4 SKIP=0`／cmd/wisp `RUN=191 PASS=110 FAIL=7 SKIP=0`／portable windows `RUN=454 PASS=311 FAIL=12 SKIP=1` ⇒ **4+7+12=23** | **同三行逐字相同**（`FAIL=4 / 7 / 12`） |

- 顶层 23 枚名册见 §1 表一（逐枚列了）。含子测试的 25 枚＝再加 `TestL1VetoNeedsAChannelTheHostReallyWired` 的两枚子臂。
- `FAIL\t<pkg>` 那一枚口径**单列**：两发都只有 **3 行**＝`internal/risk`（8.839s / 7.280s）·`internal/panel`（0.944s / 0.816s）·`cmd/wisp`（540.111s / 435.229s）。**包级 3 ≠ 名级 23**，两个口径不许混。另外那 6 行光秃秃的 `FAIL`（每枚红包 2 行）**就是台账里 warned 过的 `^FAIL` 双计形**，本腿没把它当分母。
- **11 这个数我复现不出来**，任何一条单一过滤器都给不出 11。我能指出的一条具体坏尺（本腿实测过）：`grep -E '^[[:space:]]*--- FAIL'` 直接打在**原始日志行**上是**零命中**——原始每行的行首是 `test-core`/`test-windows` 这类 job 名，不是缩进＋`--- FAIL`；必须先剥掉 `job⇥step⇥<ts>Z ` 前缀。另一条形＝`--log-failed` 是**按 job 分块**流出的，只要上游接了 `head` 或只喂了一个 job，就会少一整块（本腿第一把也栽过：`test-core` 只有 4 名，单看它会以为"整发只红 4 枚"）。按本仓那条老规矩——**"少命中/0 命中"必须先拿已知正控打一遍**——这里的正控就是乙尺（脚本自报的 `four numbers` 行），它和我的甲尺互相咬合，所以 **23 是硬的、11 不硬**。
- 顺手一枚新坏的尺（要登记，别当已知好尺用）：**`gh run view … --log-failed/--log` 的"步名"归属在 `36556778053` / `36500353477` / `36149256584` 三发上全部打成 `UNKNOWN STEP`**，只有 `36559498617` 给得出步名。⇒ 步级归因只能从**正文**复原（本腿靠 `portable-tests.sh: four numbers` 行、`##[group]Run …` 回显、以及 `--- FAIL` 行里点名的测试文件来定位是哪一步）。谁再拿 `--log-failed` 的第二列当步名分母，会量出"一步没跑"。

### 0.2 台账那两枚名（`TestTheRendererHoldsExactlyOneDoorToTheHost`／`TestPlantedRendererDoorShapesGoRed`）：它们真转绿了，顶上来的是另外三枚 panel 名

不是口径不同，也不是 runner 差——**是两枚都变绿了，`internal/panel` 那一步换了一批红名**。逐发直接读日志（不是转述台账）：

| run | sha | 该两枚的状态 | 同一步（`test-core::Portable package tests`）的 `--- FAIL` 名册 |
|---|---|---|---|
| `36094258734`（09-25 12:2x） | `cd87354` | **两枚都 FAIL** | 恰这 2 枚（`RUN=1139 PASS=721 FAIL=2 SKIP=0`）＝`A237②` 记的那支 |
| `36149256584`（09-25 14:42） | `64858d68` | 不再出现在红名册 | 1 枚＝`TestC21DesignTokensFourWayAgree`（`FAIL=1`） |
| `36500353477`（09-28 23:53） | `f7478d37` | **两枚都 PASS**（逐名 `--- PASS` 取到） | 1 枚＝`TestC21DesignTokensFourWayAgree` |
| `36559498617`（今晚 19:0x） | `a7993a9b` | **两枚仍 PASS** | 4 枚＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree` |

⇒ 台账 `A237②`（`pending-and-issues.md:6173`）**已过期**，但它自己那句预言兑现了：原文写的"甲落地之后它应当自己转绿"＝**已转绿**，绿最早出现在 `36149256584`（`64858d68`）。两枚定义所在文件在 `f7478d37`/`a7993a9b`/`HEAD` 三版都还在（`internal/panel/composer_test.go`，`git grep` 现取），不是用例被删。

### 0.3 C21 那枚"已知常红读"成立，出处不是转述

`TestC21DesignTokensFourWayAgree` 今晚红、`f7478d37` 与 `64858d68` 两发也红；`36094258734`（09-25 上午）那一发它**还是 PASS**。owner 的定档在 `Q-52`（撤回那条，`pending-and-issues.md:1089`）：②"**本编队此后不动 `internal/panel/tokens_fourway_test.go` 一字**（不修、不'顺手让它安静'、不换基准、不注释掉），那枚 `TestC21DesignTokensFourWayAgree` **一律当已知常红读**"。红因（CI 侧、不引 design/frontend 内容）＝四方对账里设计 token 源声明的**若干行**没进生成主题，末尾那句 `only 0 of 78 colour rows completed all four legs`；同一条在本机是**另一枚红因**（读不到那个 token 源文件），`Q-52` 明写"两枚红、两句红因，别让下一位把它们读成同一枚"。

---

## 1. 表一：逐名归因（23 枚名级 ＋ 3 枚步级）

档位定义照派单：〔先存在〕必须有两把更早的读数；〔今晚新增〕＝第一发红、对照发绿（或该用例在对照发上根本不存在＝标"首跑无对照绿"）；〔环境差异〕；〔判不了〕。
"机制"一列是我从日志红句里读到的东西，**它和档位是两回事**，别合起来当结论。

### 1.1 `test-core :: Portable package tests (core scope…)`（ubuntu-latest，`internal/panel` 那 4 枚）

| # | 用例名 | 档位 | 证据（更早读数／对照 run） | 机制（红句形状，不引前端/设计文件内容） |
|---|---|---|---|---|
| 1 | `TestApprovalCardViewJSONKeysMatchFrontendTypes` | **今晚新增** | `36149256584` 与 `36500353477` 逐名 `--- PASS`；今晚两发 `--- FAIL` | Go 侧 `Snapshot` 发出的两个键在前端那个接口里没有声明（`approval_test.go:129`）＝**Go↔TS 字段面漂移** |
| 2 | `TestComposerContractTypesMatchFrontend` | **今晚新增** | 同上两发 `--- PASS` | 同上，两支：`Snapshot` 多两枚键、`ComposerState` 多三枚键（`composer_test.go:74`） |
| 3 | `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | **今晚新增** | 同上两发 `--- PASS` | 一枚"随包出货"的前端组件文件里出现了主题外的色值（`frontend_hygiene_test.go:216` 点名到行；内容不转述） |
| 4 | `TestC21DesignTokensFourWayAgree` | **先存在** | `36149256584`＋`36500353477` 同枚名 `--- FAIL`；`36094258734` 还是 PASS ⇒ 变红发生在 09-25 下午，不是今晚 | 四方对账不齐（`tokens_fourway_test.go:532/:547`）；owner 已定"已知常红读"（`Q-52` 撤回条） |

造成 1/2 那两枚 Go↔TS 缝的 Go 侧提交都在今晚那一批里（`git log f7478d37..a7993a9b -- internal/panel` 现取）：`fccfe752`(09-28 15:11，git 那一维)·`1ba16de0`(09-28 18:18，`currentModel`/`modelKnown`)·`84feec44`/`a90677e9`(09-28 19:30/20:12，instructions 那一维)·`335b8d2b`(09-28 21:53，名册进快照)。**代码不是今晚写的，但 CI 第一次看到它是今晚 19:04/18:38 那两推**——`推送解除按住`之前 `origin/dev` 一直停在 `f7478d37`（`HANDOVER 4.0w` 补三/补四都记着"继续按住"）。第 3 枚的触发件是 09-26 那批前端提交里的一枚组件文件（路径历史现取：`e95caf1d 09-26 11:47`→`611ae8b9 09-26 21:26` 共 5 枚），同样"代码早就在、今晚才进 CI"。

### 1.2 `test-windows :: Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)`（windows-latest，全在 `internal/risk`）

**12 枚全部＝〔先存在〕，且机制全部＝〔环境差异〕。** 硬证据：同一枚名册在 `36149256584`(09-25 14:42) 与 `36500353477`(09-28 23:53) 上**同名同句**（我把两发的断言行做了归一化差集，**差集为空**）；本机（Windows 真机，用户 `swq`）最近一发 `.scratch/wisp/probes/221/v1/logs/ac5-full.txt`（09-29 20:04，非 `-v`）里 `ok github.com/CarlosShao/wisp/internal/risk 7.450s`＝**这 12 枚在本机全绿**。

| # | 用例名 | 档位 | 机制（红句） |
|---|---|---|---|
| 5 | `TestClassifyAnchorSpellingIsNotVerdict` | 先存在＋环境差异 | 临时根目录带 8.3 短名 `C:\Users\RUNNER~1\…`，`Resolve` 后同一枚路径拿到两种拼写（`pathresolver_anchor_spelling_windows_test.go:57`"test premise broken"） |
| 6 | `TestCanonicalInputGainsNoSecondForm` | 先存在＋环境差异 | 同一枚 canonical 产出 **2 种比较形态**（`:140`，want 1） |
| 7 | `TestAListWinsWhereBothTablesHit` | 先存在＋环境差异 | `:238` "B-only control must stay overridable by design"（同族，锚被短名打成 B） |
| 8 | `TestBListDefaultDenyAndOverride` | 先存在＋环境差异 | 同族 |
| 9 | `TestPathResolverShortNameAListDenied` | 先存在＋环境差异 | 同族（短名那一支） |
| 10 | `TestPathResolverUNCAListDenied` | 先存在＋环境差异 | 同族 |
| 11 | `TestPathResolverExtendedLengthPrefixAListDenied` | 先存在＋环境差异 | 同族 |
| 12 | `TestSyncFixtureFallbackAndMatch` | 先存在＋环境差异 | `syncdirs_test.go:162` "write under fixture root must be sync"；同发 `syncdirs_redteam_windows_test.go:220` 明写 **`no live sync root on this machine (detected roots: [])`** |
| 13 | `TestSyncFallbackNotDisarmableByWeakRoot` | 先存在＋环境差异 | `:233` "suspect fallback must stay armed" |
| 14 | `TestSyncUnverifiedRootKeepsFallback` | 先存在＋环境差异 | `:213` 四臂（`""`/`default`/`fixture`/`options`）"under-profile fallback must stay armed" |
| 15 | `TestSyncSuspectFallbackIsComponentBounded` | 先存在＋环境差异 | `:355` "under-profile must be suspect" |
| 16 | `TestSyncSuspectFallbackWhenUndetectable` | 先存在＋环境差异 | `:402` "under-profile path must be sync-suspect, got {Sync:false …}"；`:494` 那句 `P12 evidence: no sync location detected on this machine` |

⇒ 两族机制都写在红句里：**(甲) GitHub 托管 windows runner 的 profile 是 `runneradmin` 且有 8.3 短名别名 `RUNNER~1`；(乙) 那台 runner 上没有真同步根（OneDrive/Dropbox 未装配）**。本机两条都不成立（`C:\Users\swq\…`、有真根）⇒ 这 12 枚**不是代码回归**，也**不是"只能在 linux 容器里量"那类**（它们是 `*_windows_test.go`，本机能量且能量绿）。
旁证一枚（不是红）：`ci.yml:519-524` 那段注释还写着"这一步（`PathResolver junction placeholder`）今天在 windows-latest 会因 `RUNNER~1` 判 B 而**红**"，但今晚两发那一步的结论是 **success**，日志实读到 `--- PASS: TestPathResolverJunctionWindows (0.03s)`、`runtests.sh: OK … PASS=1 FAIL=0 SKIP=0` ⇒ **那段注释已过期**（同一条 8.3 短名仍在打 5–11 那几枚）。本腿不动它，只登记。

### 1.3 `test-windows :: cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)`（windows-latest，`cmd/wisp`）

| # | 用例名 | 档位 | 证据 | 机制 |
|---|---|---|---|---|
| 17 | `TestComposedGateBlocksAWriteForTwoSeconds` | **先存在** | `36149256584`／`36500353477` 同名同句（`run_test.go:378`），今晚耗时 **304.08s**、09-25 那发 301.57s、09-28 那发同形 | 期望"没人否决的 L1 窗口＝执行"，实得"审批超时（**300 秒**未确认），C18 一律判拒绝"。CI 上把这枚 2 秒窗口的判据等成了 300 秒默认超时 |
| 18 | `TestTicket101ManualSwitchSurvivesRestart` | **先存在** | 两发同名同句（`run_mode101_test.go:309`＋`:320`） | 期望那枚 silenced L1 写"该执行"，实得"审批超时（1 秒未确认）…已自动拒绝"；`:320` 文件没落盘 |
| 19 | `TestTicket101UntouchedConfigRestartsAtDefault` | **先存在** | 同上 | `:482` "control write should NOT error under auto_approve … the control half is broken, so the refusal above proves nothing"＝**半数是坏的对照组**，红句自己这么说 |
| 20 | `TestTicket101SessionGrantDoesNotCrossRestart` | **先存在** | 同上 | 同族 |
| 21 | `TestL1VetoNeedsAChannelTheHostReallyWired`（含 2 枚子臂） | **今晚新增（首跑，无对照绿）** | 该测试文件在 `f7478d37` **不存在**（`git grep -l "func TestL1Veto…(" f7478d37` 空、`a7993a9b` 有），入库件＝`58bf0158`（09-29 11:18）；CI 从来没有过它的绿读数 | `approval_reply_201_test.go:540` "the L1 timeout polarity changed (this is the D4 contract, not this leg)"＋`:632` "audit is missing the veto line"；两枚子臂分别 38.64s／1.26s。**机理＝判不了**（见 §4） |
| 22 | `TestTicket223PermissionDeniedSitsInItsOwnSentence` | **今晚新增（首跑，无对照绿）** | 文件 `config_reload_perm_223_windows_test.go` 入库 `ec7a034d`（09-29 15:07），`f7478d37` 查无此名 | `:95` "**icacls denied nothing**, so this case cannot show the permission sentence"＝那台 runner 上造不出"权限被拒"这个被试物 |
| 23 | `TestRunPacketCarriesTheLoadedInstructionFiles` | **今晚新增（首跑，无对照绿）** | 文件 `instructions_200r2_test.go` 入库 `a90677e9`（09-28 20:12） | `:167` 期望键与实得键**是同一枚路径的两种拼写**（`…\RUNNER~1\…` vs `…\runneradmin\…` 全小写）＝与 5–11 同族的 8.3 短名机制 |

**步级 3 枚：**

| 步 | 档位 | 证据 | 机制 |
|---|---|---|---|
| `lint :: staticcheck` | **先存在**（今晚两发红、`36149256584`／`36206542728`／`36500353477` 同步步也红） | 台账在册＝票 85（钉版本）＋票 122（清 finding，"不许为变绿放宽"） | 自报 `version=staticcheck 2026.2.1 modules=3 packages=36 findings=49 exit=1`；基线 `36500353477` 是 `packages=34 findings=47`。逐条差集：净＋2＝**3 枚新 finding**（`cmd/wisp/approval_reply.go:91` `maxTrackedCards` unused／`internal/panel/git_test.go:448` `composerMethodWhitelist` unused／`internal/tools/subagent_197.go:294` S1016）**－1 枚已消失**（`internal/config/parse.go:233` `fileMissing` unused）；其余"看起来变了"的全是行号漂移（`run.go:110→111`、`:252→323`、`:324→395`、`queue.go:472→504`）。⇒ finding 增量是这批带进来的**未接线残留**，不是尺坏了 |
| `lint :: gofmt (gofumpt)` | **今晚新增**（`36500353477`/`36206542728`/`36149256584` 三发该步 **success**；今晚两发 **failure**） | 触发件唯一：`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1 imports must appear before other declarations`，`gofumpt` 以 **exit 2** 死在解析上；那枚件由 `4813567e`（**09-28 14:42**）入库，`git merge-base --is-ancestor 4813567e f7478d37`＝**NO**、`…4d01e0b2`＝**YES** ⇒ 今晚 18:38 那推第一次把它送进 CI | **不是今晚写的码**：`ci.yml:152-172` 那 20 行注释已经把这件事定成结构性空窗——"`gofumpt -l . tools/… must be EMPTY` 在 bench 树上**字面不可满足**（gofumpt 会走 `.scratch/**`，那里按定义放着故意坏的样本）"。先例＝`36124009826`（sha `7b2a8703`，`2026-09-25T10:26:43Z`＝本机 18:26 +08）该步也红过，之后 `36140078431`/`36149256584`/`36206542728`/`36500353477` 连续四发绿 ⇒ **这枚步红是"有没有坏样被推进去"的函数，不是"码好不好"的函数** |
| `slo-full :: Upload SLO report` | **今晚新增（只在后一发红）**：`36556778053` 该步 **success**、`36559498617` 该步 **failure**；更早各发均 success | `gh api …/artifacts` 实读：今晚 `36559498617` **只有 `slo-smoke-report`（7200 B）**、**没有 `slo-full-report`**；前一发有 `slo-full-report` 20,701 B（`created 2026-09-29T10:40:19Z`） | 红句＝`##[error]Failed to CreateArtifact: Unable to make request: ECONNRESET`，Actions 自己跟一句"self-hosted runner 请确认能访问所有 GitHub 端点"。该步是 `if-no-files-found: warn`，而日志明写 `With the provided path, there will be 1 file uploaded` ⇒ **文件是有的，红在上传那一跳** ⇒ **环境差异／瞬时网络**（runner＝`wisp-selfhosted-01`，机器名 `DESKTOP-LVS7839`，即本机），**零代码含义** |

---

## 2. 表二：这两发的"步的覆盖形状"——我采不到读数的步

两发的 job 集合都是 6 枚（`lint` / `lint-frontend` / `test-core` / `test-windows` / `slo-smoke` / `slo-full`），**没有整枚 job 没跑、也没有整枚 job skipped**。`slo-fresh.yml` 只在 schedule／dispatch 触发，push 不跑（今日 `05:52:58Z` 那发 success，与这两发无关，别混进归因）。

**采不到＝红步之后被吃掉的步（两发逐字同）：**

| job | 步名 | 结论 | 为什么采不到 |
|---|---|---|---|
| `lint` | `gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)` | **skipped** | 它自己**没有** `if:`（`ci.yml:146-172` 明写"不许给它 `if:`、不许 continue-on-error"），被上一步 gofmt 的红吃掉 ⇒ 本仓登记过的"一步红吃掉后续步"形状原样复发；而且那 20 行注释自己承认它"**adds NO REACH TODAY**" |
| `lint` | `go vet (module)` | **skipped** | 同因。⇒ **今晚两发都没有任何 `go vet ./...` 读数** |
| `lint` | `go vet (tools/d22scan module)` | **skipped** | 同因 |
| `test-core` | `Post Run actions/setup-go@v5` | skipped | runner 自带收尾，无判据，无害 |
| `test-windows` | `Post Cache third_party (deps.toml key)` / `Post Run actions/setup-go@v5` | skipped | 同上，无害 |

`lint` 里排在 gofmt **之后**却仍拿到读数的只有两步＝`staticcheck`（带 `if: ${{ !cancelled() }}`，`ci.yml:231`）与 `mockllm module vet`（`:268`），两发都自报了结论 ⇒ R-4 那一族在这一发上没复发。gofmt **之前**的三步（D22 正控／D22 自测／D22 七 ban＋emoji 扫）两发全 success、读数可采。

**"步层绿但零读数"这一类，今晚抓到一枚真的（这是本表最要紧的一行）：**

| job | 步 | 结论 | 实读 |
|---|---|---|---|
| `slo-full`（self-hosted） | **`SLO full gate (six states + settle + leak)`** | **success** | **该步名在整份 job 日志里出现 0 次、正文 0 行**。前一步 `Build wisp.exe` 最后一行 11:05:18.867Z，后一步 `Upload SLO report` 第一行 11:07:24.480Z ⇒ **中间 2 分 06 秒没有任何日志**，而 `-Subset full -SecondsPerState 6`＝6 态×6s＋settle 10s＋leak 10s ≈ 66s＋起进程，时长对得上"跑了但没写日志"。同发对照正控：`slo-smoke`（托管 windows）那一步**照打了全量**（precheck ok → `Sleeping/Warm` → settle → leak `flipped_to_fail=True`），前一发 `36556778053` 的这同枚 `SLO full gate` 也照打全量（六态 `pass=True`＋`settle pass=True`＋`leak exit=1 flipped_to_fail=True`＋`report written … all_pass=True`）⇒ **零输出不是 gh 日志的普遍毛病，是这一发这一步独有的** |

⇒ 直接后果（照派单要求具名"我采不到的步"）：**`slo-full::SLO full gate` 在 `36559498617` 上＝绿而不可读**；再加 `36559498617::slo-full::Upload SLO report` 红，那一发**没有任何可证的 slo-full 样本落进工件**（API 实读只有 `slo-smoke-report`）。仓里 `build/slo/slo-report.json` 那份字节今晚确实存在（所以 upload 数到"1 file"），但**从日志无法证明它是 19:0x 这发写的，而不是 18:4x 那发留下的**（同一台机器、同一个 `_work/wisp/wisp` 工作目录、`build/` 不被 checkout 清掉）。⇒ **明早任何人引用"今晚 slo-full 绿"都拿不出读数支撑**，也不许写成"D32 达标"。

**另一枚采集面损失（口径级，不是步级）**：`--log-failed` 只给**红步**正文，绿步一律取不到——本表里所有"绿但有输出"的判断都是我用 `--log --job <id>` 逐 job 另拉的（`slo-full`／`slo-smoke`／`test-windows`）。只拉 `--log-failed` 的人会以为"绿步＝没内容"。

---

## 3. 表三：与本地（Windows 真机）红名册的差集

本地读数取**最新一发** `.scratch/wisp/probes/221/v1/logs/ac5-full.txt`（09-29 20:04:36，非 `-v`，`cmd/wisp 118.938s`）：**名级红恰 6 枚**，台账自己那份判语在 `f6186158`（"23 ok／3 FAIL／5 无测试、名级红 6 枚、cmd/wisp 那枚隔离复量 3/3 绿＝计时红"）：
`TestC21TableColourRowsMatchTokensCSS`（`internal/ball`）·`TestApprovalCardViewJSONKeysMatchFrontendTypes`·`TestComposerContractTypesMatchFrontend`·`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`·`TestC21DesignTokensFourWayAgree`（`internal/panel`）·`TestTicket223RefusedLooseningKeepsOldValues`（`cmd/wisp`，计时红）。
（对照：09-29 13:01 那发 `.scratch/wisp/probes/222/r1/gates-full.txt` 同形，只是 `cmd/wisp` 当时整包 `ok 103.330s`、多一枚 `internal/risk/TestResolvePerCallBudget` 计时红。）

### 3.1 方向 A：CI 有而本地没有 = **19 枚名 + 2 枚步**

19 枚名 = §1.2 的 **12 枚 `internal/risk`**（`*_windows_test.go`，**本机能量且量到绿**，`ok internal/risk 7.450s`）+ §1.3 的 **7 枚 `cmd/wisp`**（`TestComposedGateBlocksAWriteForTwoSeconds`、`TestTicket101ManualSwitchSurvivesRestart`、`TestTicket101UntouchedConfigRestartsAtDefault`、`TestTicket101SessionGrantDoesNotCrossRestart`、`TestL1VetoNeedsAChannelTheHostReallyWired`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`、`TestRunPacketCarriesTheLoadedInstructionFiles`；本机这一发只红 `TestTicket223RefusedLooseningKeepsOldValues` 那一枚，7 枚全不在名册里）。
2 枚步 = `lint::gofmt (gofumpt)`、`slo-full::Upload SLO report`。

⚠ **一枚口径警告（别把这句读成"CI 独有问题"）**：这 19 枚**零枚**属于 `*_other_test.go`／`//go:build !windows` 那种"只能在 linux 容器里量"的形状。它们是 `*_windows_test.go` 或无平台后缀的 `cmd/wisp` 用例，本机（Windows）**能量、且量出来是绿** ⇒ 差异方向是 **runner 环境**（`RUNNER~1` 8.3 短名／没有真同步根／`icacls` 造不出拒绝／无交互审批通道），不是平台覆盖缺口。
⚠ 第二枚口径警告：本机那两发跑在 `HEAD≈29215a93/f6186158` 的**工作树**上，比被验两发（`4d01e0b2`/`a7993a9b`）**多 6 枚提交**且其中含 `cmd/wisp/run.go`、`internal/tools/task.go`、`internal/tools/subagent_197.go` ⇒ "本机绿"不能当成"`a7993a9b` 那棵树也绿"的证明，只能当成"这台机器不产生这一形红"。
⚠ 第三枚：`gofmt` 这一枚在本地门禁里出现过"空输出"的自报（`abaedd7a` 那节的"gofumpt v0.12.0 空输出"），而 CI 红——**两者口径不同**（CI 是 `gofumpt -l . tools/d22scan tools/mockllm` 整棵树，含 `.scratch/**` 坏样）。别把本地那句读成"格式门本地也绿"。

### 3.2 方向 B：本地有而 CI 没有 = **2 枚名**

| 用例名 | CI 今晚 | 本机 | 为什么不算 CI 的事 |
|---|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS`（`internal/ball`） | **PASS**（`36559498617` 与 `36500353477` 都逐名取到 `--- PASS`；`./internal/ball/` 在 linux 与 windows 两条腿的名册里都真跑了） | **FAIL** | 红因＝本机工作树里那个设计 token 源文件被挪出工作树（`Q-52` 撤回条与 `A208③` 都写着"本机红句是 `read design/assets/tokens.css: The system cannot find the path specified`"；CI 从提交树取，文件在，所以绿）⇒ **工作树状态，不是 CI 独有问题，也不是回归** |
| `TestTicket223RefusedLooseningKeepsOldValues`（`cmd/wisp`） | **PASS**（今晚逐名 `--- PASS`） | **FAIL（计时红）** | 台账 `f6186158` 自己判"隔离复量 3/3 绿＝计时红"；且本机那一发在别的腿并发的窗口里 ⇒ **别当真红登记**。⚠ 该用例文件在 `a7993a9b` 存在（`cmd/wisp/config_reload_223_test.go`，`git grep` 现取），所以 CI 确实跑了它、确实绿 |

另有 1 枚"两台机器都能出"的形，不进差集：`internal/risk/TestResolvePerCallBudget`——本机争用时红（`A428`/`A432` 已归因为 CPU 争用假红、`thresholds.go` 零字节未动），而 CI 三发逐名都取到 `--- PASS`：`36559498617`＝1.13s 与 2.00s 两枚（linux 腿＋windows 腿各一枚）、`36556778053`＝1.98s 与 4.35s、`36500353477`＝2.00s 与 2.81s ⇒ **CI 这一侧从来不复现它**。

**"只能在容器量"那一类**：本腿在这两发的红里数到 **零枚**（理由见 3.1 的第一条警告）。反向的"只能在 Windows 量"那三枚（`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`／`TestD34WriteMatrix`／`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`）由 `portable-tests.sh` 自己的 ledger 在 linux 上 `-skip` 掉，日志逐行点名并给理由（fixture 需要 windows 独有的对象），**它们今晚没红，也不是本轮的债**。

---

## 4. 判不了的格子（缺哪份证据，写死）

1. **`36559498617::slo-full::SLO full gate` 绿而零输出**——分不清"跑了但 stdout 没进日志"与"这一步根本没执行/日志件丢了"。缺的证据＝self-hosted runner 本机 `E:\work\base\actions-runner\…\_diag\Worker_*.log`（在盘上、不在仓里、`gh` 取不到），或 `build/slo/slo-report.json` 的 mtime（同样在 runner 工作目录里，本腿不去读别人机器状态、也不动它）。**这一格不许被写成"slo-full 绿"**。
2. **今晚那枚 `slo-full-report` 工件不存在**，最近一枚有效样本停在 `36556778053` 的 `10:40:19Z`（约 19:40 +08）。`ci.yml:639-641` 写明这份工件的存在是 `scripts/slo-freshness.sh` 探针 P3 的输入 ⇒ 缺的证据是"下一次成功产样在什么时候"，本腿不预测。
3. **`TestL1VetoNeedsAChannelTheHostReallyWired` 的机理**——我只证到"CI 首跑即红＋本机这一发不在名册里＋红句自己点名 D4 极性那一支与缺一条审计行"。到底是"CI 那台 runner 没有可接的审批通道/没有 TTY"还是"D4 那支语义在两条路上不一样"，需要**同一枚用例在两个环境各跑一遍**（本腿禁编译）才能定。
4. **§1.3 那 4 枚先存在的 `cmd/wisp` 红的机理**——同名同句在 09-25/09-28/今晚三发一致（这一条已证），但"为什么托管 windows runner 上那枚窗口会走满默认超时"我拿日志证不了，需要一次定向复现。
5. **`gofmt` 那一枚的"是否只红在 `fs_broken.go`"**——步正文只打了这一行 `##[error]` 就以 exit 2 收；`gofumpt -l` 的完整清单在这条 `OUT="$(…)"` 的命令替换里**没被打出来**（步骤是 `bash -e`，命令替换返回 2 就直接终止，永远走不到 `echo "$OUT"`）。⇒ **"还有没有第二枚坏样"这一格判不了**，缺的证据是一次能同时给 rc 和清单的跑法（属修步形状，超出本腿射程，本腿不动）。

---

## 5. 台账里那几条在册红的"锚点是哪一发"逐条现读结果

| 台账条目 | 它写的锚 | 本腿现读 | 状态 |
|---|---|---|---|
| `A237②`（`pending-and-issues.md:6173`）：`test-core::Portable package tests` 名册恰 2 枚＝renderer-door 那两枚 | run `36094258734`（sha `cd87354`） | 该发日志实证：`RUN=1139 PASS=721 FAIL=2`，两枚逐名 `--- FAIL` ⇒ 当时**为真** | **已过期**（`64858d68` 起转绿，今晚两发逐名 `--- PASS`）。它自己那句"甲落地之后它应当自己转绿"＝**兑现**。同一行末尾"`lint`/`test-windows` 两枚红的步名与昨天基线逐字同 ⇒ 先存在"——**今晚仍成立**（步名逐字同） |
| `Q-52`（撤回条，`:1089`）：`TestC21DesignTokensFourWayAgree` 一律当已知常红读；编队不动 `tokens_fourway_test.go` 一字 | 09-25 19:2x | 今晚红、`f7478d37`/`64858d68` 红、`cd87354` 及更早绿 | **仍有效**，且是本轮唯一被 owner 定档"常红读"的名 |
| 票 85 / 票 122（`:1232`、`:1486`、`:3467`、`:6324`）：`lint :: staticcheck`＝先存在、清 finding 单独排票 | 多轮 | 今晚 `findings=49` 红；基线 `36500353477` `findings=47` 红；09-25/09-26 亦红 | **仍有效**；增量＝净 +2（3 枚新、1 枚掉），属票 122 那张票的账，**本腿不放宽任何断言** |
| `:9551`（今晚那条 🔴）：两发都红、"11 枚唯一 `--- FAIL` 名"、"与 `A237②` 名册不重合＝要么口径不同要么真转绿" | 本腿起手时 | **转绿那一支为真、口径那一支为假**；"11 枚"这个数**不成立**（真值 23 顶层／25 含子测试，两把尺互咬） | 该条的"〔待归因〕"前缀可以摘，但**摘的人不是我**（双角色规矩：登记与翻勾归编排者） |
| `ci.yml:519-524` 那句"该步今天在 windows-latest 会红" | run `35551819606` | 今晚两发该步均 **success**、日志实读 `--- PASS: TestPathResolverJunctionWindows (0.03s)` | **注释过期**（登记，不改） |

---

## 6. 派单第 3 节那一句：今晚这两发红，会不会影响明早 10:30 的桌面视觉签收

**答案：会影响，不是"零枚"。具名 4 枚**，全部在 `test-core::Portable package tests` 那一步（今晚两发红的是同 4 枚）：

1. `TestApprovalCardViewJSONKeysMatchFrontendTypes` —— 钉的是"审批卡的 JSON 字段面与前端类型是否一致"。**签收要看的就是审批卡** ⇒ 它是判据，不是背景噪声。
2. `TestComposerContractTypesMatchFrontend` —— 钉的是输入区状态（含 `git`／`currentModel`／`modelKnown` 那一维）的字段面是否一致。**明早看面板输入区就撞在它上面。**
3. `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` —— 钉的是"色值只许活在生成主题里"。**这是纯视觉判据。**
4. `TestC21DesignTokensFourWayAgree` —— C21 色值四方对账。owner 已定"已知常红读"（`Q-52` 撤回条），所以它**不该改变签收动作**，但它红的那一维恰恰是"哪些色值从来没到过面板"。

三条限定，一条都不许省：

- **这 4 枚在本机同样红**（09-29 20:04 那发逐名 `--- FAIL`）⇒ 明早真机签收时它们**一样成立**；CI 那两发红**没有额外**改变明早能看到的东西。真正的缺口是那 4 枚红各自指出的**字段面/色值缝本身**，而 owner 已经拍过"这个界面还不是最终定稿……先等等"（`Q-52` 撤回条引的原话）。所以**结论＝"别把签收判成通过／不通过"，而不是"CI 挡住了签收"**。
- **球／托盘那一侧：红名册里零枚**。`internal/ball` 今晚 CI 两发全绿（含 `TestC21TableColourRowsMatchTokensCSS` 逐名 `--- PASS`），名册里也没有任何 tray 用例；本机那枚 ball 红是**工作树里 token 源被挪走**、不是界面行为。据以判断的步名＝`test-core::Portable package tests`、`test-windows::Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)`（后者的 linux 孪生是 test-core 那一步，两条都真跑了 `./internal/ball/`）。
- **另一侧风险（不属"红名册"、但直接关系签收判断的资源那一半）**：`slo-full::SLO full gate` 在 `36559498617` 上**绿而零读数**、同发 `Upload SLO report` 红且没落工件 ⇒ **今晚没有任何可证的 slo-full 样本**。明早若有人要说"球常驻那套资源达标／不达标"，**这份 CI 读数给不了他依据**（照旧规矩，也绝不许写"slo-full 绿＝D32 达标"）。据以判断的步名＝`slo-full::SLO full gate (six states + settle + leak)`、`slo-full::Upload SLO report`；有读数的最近一发是 `36556778053`（六态全 `pass=True`＋settle `pass=True`＋leak 自测 `flipped_to_fail=True`＋`all_pass=True`，工件 20,701 B 在）。

**"今晚新增"那一档到底几枚**：**名级 6 枚**＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（三枚有"更早绿"的对照发）＋`TestL1VetoNeedsAChannelTheHostReallyWired`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`、`TestRunPacketCarriesTheLoadedInstructionFiles`（三枚是**首跑即红、根本没有对照绿**，用例文件在 `f7478d37` 不存在）。**步级 2 枚**＝`lint::gofmt (gofumpt)`、`slo-full::Upload SLO report`。合计 8 枚。其余 17 枚名＋`lint::staticcheck`＝先存在（各有两把以上更早的同名同句读数）。

---

## 7. 本腿没做的事（免得被读成做了）

- 没跑任何编译／测试／门禁脚本；没动 `.gitignore`、`design/**`、`frontend/**`、任何跟踪文件；没 commit、没 push。
- 没读、没引、没转述 `frontend/**` 与 `design/**` 的内容；提到它们只给**测试函数名＋CI 步名＋测试自己点名的位置**。
- 没修任何东西，包括没"顺手让 `gofmt` 安静"、没给任何步加 `if:`、没动 `thresholds.go`／golden／任何断言。
- 台账登记、翻勾、以及 §5 那几条"过期"的正式改档，**一律不归本腿**（双角色：我是只读归因腿，不是裁决者，也不是编排者）。
