# 票 124 AC#2b 批次2（`worker-ticket124-ac2b-2`）— tools 20 + llm 17 + perm 5 + approval 1 接上票 119 那条纪律

**日期**：2026-09-23 · **性质**：转换方交件（改 `.go`，但只改 `_test.go`）
**本批范围**：清点账 `docs/evidence/s1/124-ac2a-leg-classification.md` §5.8（tools 行 107-124、126-127，除行 125 归票 123）、§5.5（llm 行 52-68）、§5.6（perm 行 69-73）、§5.3（approval 行 44），共 **43 枚**，账上全标「可转」。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2b 的放行与分批」（09-23 16:33 编排者）那张表 + 五条共用结案判据；批次 1 的 `next=` 两条本批提醒照办。

## 0. 锚点 / 快照 / 跑法（可复核）

- 开工首读 `git rev-parse --short HEAD` = **`5417a3c`**（与派单给的锚点一致；下文记 **pre 锚**，全部改前读数在它上面量）。
  ⚠ 共树漂移：测量途中 HEAD 先后经 `4d266f9`、`38c09df` 到本批改件 `dfa3dc4`。已核 `git diff --stat 5417a3c dfa3dc4^ -- internal/ go.mod go.sum` **零 hunk**；`git diff --name-only 5417a3c dfa3dc4^` = 9 枚文件，全部在 `.scratch/wisp/issues/`、`docs/`、`cmd/wisp/**`（兄弟代理的 docs 与 2b-4 地界的 `cmd/wisp` 码）⇒ 与本批四包无交集，pre 读数量的是被改前同一版本。
- 改件 commit：**`dfa3dc4`**（`test(124,AC#2b-2)`，18 枚文件，+161/-30；其中删除侧 30 行全是递根点，见 §3）。
- 全部测量都在 `git archive <sha> | tar -x` 的纯净快照里做；工作树里别人的半成品（`.scratch/.../134-*.md`、`docs/reports/*.md`、`cmd/wisp/leg_dispatch_gate_133_test.go`、`cmd/wisp/part133b.go.txt`）一个字都不读。
- 容器 `golang:1.27`（`go1.27.1 linux/amd64`，容器内 `CGO_ENABLED=1`），复用票 119 的命名卷 `ac119-gomodcache` / `ac119-gocache` ⇒ 离线可编。
- 挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`；每枚样本进容器第一件事打 `ls -l /src/go.mod` + `md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum resolve.go` = `b6876a5efe759f6e17434d1b50a129c3`（与 AC#1/AC#2a/AC#2b-1 逐字同字）⇒ 非空挂自证。
- 形状硬断言沿用批次 1：软链形 `exit 97`（`/varlink` 不是 symlink）/ `exit 98`（`readlink -f /varlink/w124tmp != /realpriv/w124tmp`）；普通形 `exit 99`（`/varlink` 必须根本不存在、`/plainroot` 不许是链接）。
- 跑法脚本（新建，未改 2b-1 那枚）：`/d/tmp/wisp124-2b2-h.sh`（`bash /h.sh <link|plain> <tag> [pkgs]`，`COUNT=2` 切门禁那一发；`-timeout 25m`——tools 软链形含票 123 那两枚 ~300.0x s 的已知腿，12m 会掐）。变异台：`/d/tmp/wisp124-2b2-m.sh`（判据⑤：先 grep 证落地 + `go build` rc=0 再读红名）。转换件：`/d/tmp/wisp124-2b2-convert.py`（30 处行号锚定替换，锚不中即 `ANCHOR MISS` 停手）。日志全在 `/d/tmp/wisp124-2b2-logs/`。
- 开测前查在飞 run（本机就是 self-hosted runner）：18:13 本地（`10:13z`）`gh run list --limit 5` 最近 5 枚全部 `completed` ⇒ 无在飞争用。
- 快照表：

| 快照目录 | 内容 | 文件数 |
|---|---|---|
| `/d/tmp/wisp124-2b2-pre` | `git archive 5417a3c`（改前基线） | 1039 |
| `/d/tmp/wisp124-2b2-post` | `git archive dfa3dc4`（改后） | 1046 |
| `/d/tmp/wisp124-2b2-probe` | 改后 + 判据④取字串探针（5 处插入、只 `t.Logf`，**未入库**） | 1046 |
| `/d/tmp/wisp124-2b2-mut` | 待量：改后 + 判据⑤变异（**未入库**） | - |

## 1. 对派单数字的复核（派单给的每个数字都是断言）

- 清点账点名（读 `124-ac2a-leg-classification.md` 逐行）：tools §5.8 = 行 107-127 共 21 枚，其中行 125（`TestL1WriteGoesThroughTheRealBlockWindow`）标「归因待票 123」不计本批 ⇒ **20**；llm §5.5 = 行 52-68 ⇒ **17**；perm §5.6 = 行 69-73 ⇒ **5**；approval §5.3 = 行 44 ⇒ **1**。合计 **43**。⇒ 与派单「20+17+5+1=43」**一致，登记差 0**。
- **实测复算（pre 锚软链形 `-v` 逐名，台 PRE-L）**：tools 顶层 `--- FAIL` **21** + 子测试 `--- FAIL` **2** = 23 枚；llm 顶层 **11** + 子 **6** = 17 枚；perm 顶层 **5**；approval 顶层 **1**（其子测试 19 枚全 PASS）。⇒ 本锚 only-in-link 共 **46** 枚，**减票 123 那三枚 300 s 腿 = 43**，与账上名册逐名对得上（名册差集脚本核，`/d/tmp/wisp124-2b2-logs/ROSTER-*.txt`）。
- ⚠ **登记一枚账上没有、本锚量到的新红（+1 差，按实测做）**：`internal/tools/TestFSReadOnlyNeverOpensACard`（`wiring_test.go:300`）软链形 `--- FAIL (300.04s)`、普通形 `--- PASS (0.00s)`，红串逐字 `out={Text:审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝 IsError:true RiskLevel:L2 ErrorClass:user_rejected}` —— 与账上已单列的 `TestL1Write`（300.04s）、`TestLateVeto`（300.01s）**同族同文案**（C18 审批超时腿，非密封拒）。
  **为什么账上没有**：AC#2a 那台 tools 软链形跑到 **720 s 被包级 `-timeout` panic**（其日志 `PRE-L` 对应物 `A2L.__internal_tools_.txt` 尾行 `FAIL github.com/CarlosShao/wisp/internal/tools 720.105s`），三枚 300 s 腿只跑完两枚就被掐，第三枚从未计入红名；本批改跑法 `-timeout 25m`（900.3 s 跑完全部）才显形。⇒ 它与另两枚同归口**票 123**，**不在本批 43 枚清零目标内**；`internal/tools/wiring_test.go` 本批一字未动，改后软链形它照常红（§2 如实登记，不算本批转绿，也不算破判据②——普通形它本就 PASS 0.00s）。
- ⚠ 票 123 那两枚已知红（`TestL1Write` / `TestLateVeto`）**不在本批 43 枚**：`internal/tools/wiring_test.go` 本批一字未动（§3 的 `git show dfa3dc4` 文件清单可核），它们在改后软链形照常红，逐名登记为已知遗留。
- pre 锚两形四数（`-count=1 -v`）：

| 台 | 形状 | tools | llm | perm | approval |
|---|---|---|---|---|---|
| PRE-L | 软链 | `rc=1 RUN=79 PASS=41 FAIL=21 SKIP=3 SUBPASS=12 SUBFAIL=2` | `rc=1 RUN=72 PASS=52 FAIL=11 SKIP=0 SUBPASS=3 SUBFAIL=6` | `rc=1 RUN=14 PASS=9 FAIL=5 SKIP=0` | `rc=1 RUN=49 PASS=28 FAIL=1 SKIP=1 SUBPASS=19 SUBFAIL=0` |
| PRE-P | 普通 | `rc=0 RUN=79 PASS=62 FAIL=0 SKIP=3 SUBPASS=14 SUBFAIL=0` | `rc=0 RUN=72 PASS=63 FAIL=0 SKIP=0 SUBPASS=9 SUBFAIL=0` | `rc=0 RUN=14 PASS=14 FAIL=0 SKIP=0` | `rc=0 RUN=49 PASS=29 FAIL=0 SKIP=1 SUBPASS=19 SUBFAIL=0` |

## 2. 判据①：逐枚转绿且点名（`-v` 才有 PASS 名）

四台对照（同一批用例，`-count=1 -v`）。PRE-L=改前软链、PRE-P=改前普通、POST-L=改后软链、POST-P=改后普通。逐名状态取自各台 `-v` 日志的 `--- (PASS|FAIL|SKIP)` 行（名册 `/d/tmp/wisp124-2b2-logs/ROSTER-*.txt`，16 枚名册文件）。

| 包 | 用例名 | PRE-L(改前软链) | PRE-P(改前普通) | POST-L(改后软链) | POST-P(改后普通) |
|---|---|---|---|---|---|
| `tools` | `TestRealToolCallWritesRewriteAccountIntoAudit` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestDeclaredRiskIsOnlyAFloor` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestToolCallRowsAreComplete` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSReadReturnsTheFileAndTaintsIt` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSReadTaintFeedsR4` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSListSummarizesADirectory` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSListHonoursItsCap` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestAtomicWriteKillsMidWrite` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSTrashGoesToTheRecycleBin` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestFSMoveSameVolume` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestPathCanonicalizerAccountsForRewrittenRoots` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestTicket107AllowlistJudgmentTwoShapes` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestTicket107AllowlistBoundaryIsComponentWise` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestCanonicalizeReturnsAPathTheOSCanOpen` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestLoopPassesDeclaredL1WriteThroughTheGate` | FAIL | PASS | **PASS** | PASS |
| `tools` | `TestLoopStillRefusesL1WhenNoGateIsRegistered` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteMeasuresBrokenFC` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteMeasuresBrokenVision` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteHonestProviderRecordsNoMismatch` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteHonestNegativeIsNotAMismatch` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteBrokenOnAllThreeDialects` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteBrokenOnAllThreeDialects/openai-chat` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteBrokenOnAllThreeDialects/anthropic` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteBrokenOnAllThreeDialects/openai-responses` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteCapableOnAllThreeDialects` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteCapableOnAllThreeDialects/openai-chat` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteCapableOnAllThreeDialects/anthropic` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteCapableOnAllThreeDialects/openai-responses` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteRefusesSilentRuns` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteSinkFailurePropagates` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteAudioStaysUnprobed` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteMeasuresBrokenThinking` | FAIL | PASS | **PASS** | PASS |
| `llm` | `TestProbeSuiteThinkingCapableIsNotTheSameAsBroken` | FAIL | PASS | **PASS** | PASS |
| `perm` | `TestTicket90ManualSwitchSurvivesRestart` | FAIL | PASS | **PASS** | PASS |
| `perm` | `TestTicket90UntouchedConfigStartsAtTheDefault` | FAIL | PASS | **PASS** | PASS |
| `perm` | `TestTicket90SessionGrantDoesNotSurviveRestart` | FAIL | PASS | **PASS** | PASS |
| `perm` | `TestTicket90ConfigKeyChangesWhatTheChainAsks` | FAIL | PASS | **PASS** | PASS |
| `perm` | `TestTicket90HandEditLooseningGoesThroughD36` | FAIL | PASS | **PASS** | PASS |
| `approval` | `TestTenOpsInOneToolCallGetOneConfirm` | FAIL | PASS | **PASS** | PASS |

**逐包四数**：

| 台 | 形状 | `internal/tools` | `internal/llm` | `internal/perm` | `internal/agent/approval` |
|---|---|---|---|---|---|
| PRE-L（改前） | 软链 | `rc=1 RUN=79 PASS=41 FAIL=21 SKIP=3 SUBPASS=12 SUBFAIL=2` | `rc=1 RUN=72 PASS=52 FAIL=11 SKIP=0 SUBPASS=3 SUBFAIL=6` | `rc=1 RUN=14 PASS=9 FAIL=5 SKIP=0` | `rc=1 RUN=49 PASS=28 FAIL=1 SKIP=1 SUBPASS=19 SUBFAIL=0` |
| PRE-P（改前） | 普通 | `rc=0 RUN=79 PASS=62 FAIL=0 SKIP=3 SUBPASS=14 SUBFAIL=0` | `rc=0 RUN=72 PASS=63 FAIL=0 SKIP=0 SUBPASS=9 SUBFAIL=0` | `rc=0 RUN=14 PASS=14 FAIL=0 SKIP=0` | `rc=0 RUN=49 PASS=29 FAIL=0 SKIP=1 SUBPASS=19 SUBFAIL=0` |
| POST-L（改后） | 软链 | `rc=1 RUN=79 PASS=59 FAIL=3 SKIP=3 SUBPASS=14 SUBFAIL=0` | `rc=0 RUN=72 PASS=63 FAIL=0 SKIP=0 SUBPASS=9 SUBFAIL=0` | `rc=0 RUN=14 PASS=14 FAIL=0 SKIP=0` | `rc=0 RUN=49 PASS=29 FAIL=0 SKIP=1 SUBPASS=19 SUBFAIL=0` |
| POST-P（改后） | 普通 | `rc=0 RUN=79 PASS=62 FAIL=0 SKIP=3 SUBPASS=14 SUBFAIL=0` | `rc=0 RUN=72 PASS=63 FAIL=0 SKIP=0 SUBPASS=9 SUBFAIL=0` | `rc=0 RUN=14 PASS=14 FAIL=0 SKIP=0` | `rc=0 RUN=49 PASS=29 FAIL=0 SKIP=1 SUBPASS=19 SUBFAIL=0` |

⇒ 本批红名数 **43 → 0**（软链形）。tools 的 `rc` 仍为 `1`，其 `FAIL=3` **逐名**是票 123 那三枚 300 s 腿（`TestL1WriteGoesThroughTheRealBlockWindow` / `TestLateVetoRendersTheApprovalLayersAppliedStepsReport` / `TestFSReadOnlyNeverOpensACard`，§1 已登记），零枚属本批名单；其余三包 `rc` 从 `1` 到 `0`。日志里的机制字串（四包合计，行计数）：`not provably resolved` **23 → 0**、`refusing to seal` **23 → 0**、`审批未通过` **4 → 0**、`审批通道尚未接入` **4 → 0**。

## 3. 判据②：普通形一枚都不许多红 + `git diff` 证判定分支一字未动

**逐数相同**：POST-P 与 PRE-P 四包**十六个数逐数相同**——tools `79/62/0/3`+`14/0`、llm `72/63/0/0`+`9/0`、perm `14/14/0/0`、approval `49/29/0/1`+`19/0`。普通形红名数 **0 → 0**，一枚都不许多红。⚠ 普通形本来 PASS 的枚（含票 123 那三枚的普通形侧 `PASS 0.00s/3.00s/3.01s`）改后仍 PASS，逐名见 ROSTER-PRE-P/POST-P 差集（脚本核，四包**逐名相等**、无新增名无丢失名）。

`git show dfa3dc4` 的**删除侧全文只有 30 行**，逐行抄（`-` 号省略；`git show --numstat` = 18 枚文件、删除列合计 30）：

- `dir := t.TempDir()` — 10 行（tools 的 ticket105/loop_approval/pathshape/paths_rewrite 站点 + perm 5 + approval 1）
- `	dir := t.TempDir()` — 1 行（perm `TestTicket90ConfigKeyChangesWhatTheChainAsks` 循环体内，缩进一级）
- `root := tempCanonical(t)` — 6 行；`outside := tempCanonical(t)` — 2 行（tools fs_test 4 + bridge_test 4）
- `root := tempRaw(t)` — 2 行；`	root := tempRaw(t)` — 3 行（tools fs_write 站点，子测试内缩进一级）
- `existing := t.TempDir()` — 3 行；`base := t.TempDir()` — 2 行
- `store, err := memory.Open(filepath.Join(t.TempDir(), "data"))` — 1 行（llm `newProbeFixture` 唯一 fixture 根）

**新增侧**（除五枚新文件外只有 30 行）：30 行全部是同名替换（`sealableTempDir124(t)` / `sealableTempCanonical124(t)` 形）。**判定分支、断言、阈值、golden 一行未动**——删除行与新增行一一对应，全部落在「递给底线/规范化器的根怎么拼」这一件事上。⇒ 票面 AC#3「不许拿放行侧放宽换绿」在本批成立：解析层在调用方（测试）这一侧，`internal/winsec` 与 `internal/risk` 的判定码一个字没碰。

范围核对：`git show --numstat dfa3dc4` = 18 枚文件、全部 `_test.go`（13 枚改点 + 5 枚新委托）；`git diff --name-only 5417a3c dfa3dc4` 里 `internal/`+`cmd/` 非 `_test.go` 命中 **0 枚**。禁改列 `internal/winsec/**`、`internal/risk/**`、`rules_gateway.go`、`allowlist.txt`、`frontend/**`、批次 1 已交的 `internal/memory/**`、`internal/config/**` **零 hunk**；`internal/tools/wiring_test.go` **零 hunk**（票 123 三枚腿所在文件）；`cmd/wisp/**` 本批零 hunk（2b-4 地界；§0 登记的 9 枚漂移文件里有 `cmd/wisp` 码，那系兄弟代理所改，不来自本批 commit）。

## 4. 判据③：两形各取一枚 + RUN/SKIP 差逐包解释

两形各一枚已完成（§2 的 POST-L / POST-P，另 `-count=2` 各一枚见 §7）。逐包解释：

- **本批四包两形 `RUN` 逐包相等、差为 0**：tools `79=79`、llm `72=72`、perm `14=14`、approval `49=49`（改前改后四个数都相等）。⇒ 批次 1 那种「父用例死在 `Open` ⇒ 子测试从未被创建」的分母收缩**没有发生在本批**：本批 43 枚的红全部落在**用例体内部**（bridge 执行步 / `SaveFile` / `RunProbeSuite` / `Veto`），`t.Run` 已建子测试后其父才判负，四台 `SUBPASS+SUBFAIL` 恒为 tools `14`、llm `9`、approval `19`、perm `0`。
- **`SKIP` 名册四台逐名同一组 4 枚**（脚本对照 ROSTER-*.\*.txt）：tools `TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`、`TestWorkspaceSwitchRefusesAJunctionToOutside`（三枚都要 Windows 侧的回收站/卷/junction 能力，POSIX 自跳）+ approval `TestDefaultDeadlineWallClockMeasurement`（墙上时钟前提自跳）。**形状无关、改前改后一字不变**，与本批解析层无关，也不在任何名册差集里。
- **怎么排除「被跳过冒充变绿」（逐枚）**：§2 表 POST-L 列 43 枚**逐枚写的是 `PASS`**（取自 `-v` 日志的逐名行，非包级 `ok`），名册里 43 枚**无一枚 `SKIP`/`FAIL`**；再核改后软链的**用例名集合**与改前普通形**逐名相等**（`comm` 双向为空 ⇒ 既没多出新用例名，也没少一枚）。
- ⚠ **AC#2a 那本账的软链形「缩分母」现象在本批不复现**（tools `RUN=77`→本批 `79`、winsec 的 125 探针族不在本批四包）——本批跑法 `-timeout 25m` 让三枚 300 s 腿跑完，`RUN` 不再被包级超时截断。

## 5. 判据④：本批「另一种拒」复算 — 逐枚实拿字符串

**取字串仪器**：`/d/tmp/wisp124-2b2-probe` = 改后快照 + 5 处**纯打印**插入（`probe-patch.py` 生成，只加 `t.Logf("ACQ124-2B2 …")`，一处不动断言/阈值）。探针台四数与 POST-L **逐数相同**（tools `79/59/3/3`+`14/0`、llm `72/63/0/0`+`9/0`、perm `14/14/0/0`、approval `49/29/0/1`+`19/0`）⇒ 探针没改变任何结局。以下字符串逐字抄自该台软链形 `-v` 日志的 `ACQ124-2B2 …` 行。

账上点名的两枚（124-ac2a §7 提醒②里落在本批的两枚）：

1. `TestFSTrashGoesToTheRecycleBin`（tools，要的是「POSIX 无回收站时 fs.trash 必须拒并说出带『回收站』」）实拿：`trash supported=false out.IsError=true out.Text="本平台无 Shell 回收站 API，fs.trash 拒绝执行（绝不退化成删除）：/realpriv/w124tmp/TestFSTrashGoesToTheRecycleBin2188308832/001/trashme.txt"` ⇒ 拿到的仍是**回收站那一族的拒**，不是未解析根那句。
2. `TestProbeSuiteRefusesSilentRuns`（llm，要的是「probe 不许静默跑」的两枚校验拒）实拿：`refusenosink err="llm: probe suite requires a HealthSink (an unread probe is not a measurement)"` 与 `refusenomismatch err="llm: probe suite requires a mismatch notice (SPEC-05 sec 3.1: declared-vs-measured must be visible)"` ⇒ 两种拒各自照拿（该枚在名册里是 PASS，探针是在断言之外**再跑一次同形调用**取字串，断言与请求计数检查未动）。

本批其余断言里带「拒 / Rejects / Refuses / still refuses / 否决 / want L2 / want L1」的枚，逐枚实拿：

3. `TestFSReadTaintFeedsR4`（tools，要 L2+R4 这一档**风险判级**）实拿：`r4 level=L2 rules=[R4] overrideBlocked=true reason="R4: 包含来自 fs.read /realpriv/w124tmp/TestFSReadTaintFeedsR44178382413/001/id.txt 的内容"` ⇒ 判级要的那档拒拿到的是自己的拒；字符串里根已是解析后的 `/realpriv/…`（这层解析起效的直接旁证）。
4. `TestLoopStillRefusesL1WhenNoGateIsRegistered`（tools，要「无门时 L1 写被拒」的闸门语义）实拿：`nogate rejected=true errorClass="permission_denied" text="该操作属于 L1 级，审批通道尚未接入（ticket 21），已拒绝执行" cards=0` ⇒ 拒因还是**缺审批通道**，不是密封底线。
5. `TestTenOpsInOneToolCallGetOneConfirm`（approval，要「L1 留级 + 10 项并一次确认 + 一次否决取消整批」）实拿：`batch level="L1" batchTotal=10 out.IsError=true out.Text="用户在 L1 确认窗口中通过「按 Esc 键」否决了本次操作（取消窗口，非撤销已写出的内容）"` ⇒ 级别回到 L1（改前软链形它拿到的是 `level="L2"`，账 §5.3 逐字），它要的拒是那扇 Esc 否决。

**结论**：本批 43 枚没有一枚的断言被「未解析根」那句替代；全批软链形日志里 `not provably resolved` / `refusing to seal` 命中数 = **0**（§2）。账上「可转 ⇒ 接解析后仍成立」这一判在本批 43 枚上逐枚复核成立。

## 6. 判据⑤：变异自证（票面 AC#4）— 三态原文

变异 = **拆掉本批新加的那一层解析**：五枚委托文件的函数体改回把 `t.TempDir()` 原样交回去（`MUTATION-124-2B2`，`proc` 引用留一行 `_ =` 保编译）。`sealableTempCanonical124` 走 `sealableTempDir124`，同一发变异把它一起退回。台子 `/d/tmp/wisp124-2b2-mut`，跑法 `/d/tmp/wisp124-2b2-m2.sh`。

```
== mutation landed (grep) ==
internal/tools/tempdir_resolved_124_test.go:28:	_ = proc.SealableRoot // MUTATION-124-2B2: the resolution layer this batch added is removed here
internal/tools/tempdir_resolved_124_test.go:29:	return t.TempDir() // MUTATION-124-2B2: hand the OS's unresolved root straight back
internal/tools/tempdir_resolved_124x_test.go:21:	_ = proc.SealableRoot // MUTATION-124-2B2: …
internal/tools/tempdir_resolved_124x_test.go:22:	return t.TempDir() // MUTATION-124-2B2: …
internal/llm/tempdir_resolved_124_test.go:22:	_ = proc.SealableRoot // MUTATION-124-2B2: …
internal/llm/tempdir_resolved_124_test.go:23:	return t.TempDir() // MUTATION-124-2B2: …
internal/perm/tempdir_resolved_124_test.go:23:	_ = proc.SealableRoot // MUTATION-124-2B2: …
internal/perm/tempdir_resolved_124_test.go:24:	return t.TempDir() // MUTATION-124-2B2: …
internal/agent/approval/tempdir_resolved_124_test.go:24:	_ = proc.SealableRoot // MUTATION-124-2B2: …
internal/agent/approval/tempdir_resolved_124_test.go:25:	return t.TempDir() // MUTATION-124-2B2: …
resolution delegations still present (want 0):
internal/tools/tempdir_resolved_124_test.go:0
internal/tools/tempdir_resolved_124x_test.go:0
internal/llm/tempdir_resolved_124_test.go:0
internal/perm/tempdir_resolved_124_test.go:0
internal/agent/approval/tempdir_resolved_124_test.go:0
go build rc=0
go vet rc=0
```

三态读数（同快照血统、同容器、同形状断言；变异先证落地再读数）：

| 台 | 形状 | 四数 | 红名 |
|---|---|---|---|
| **POST（未变异）** | 软链 | tools `rc=1 RUN=79 PASS=59 FAIL=3 SKIP=3 SUBPASS=14 SUBFAIL=0`；llm `rc=0 72/63/0/0`+`9/0`；perm `rc=0 14/14/0/0`；approval `rc=0 49/29/0/1`+`19/0` | 本批 **0 枚**（tools 那 3 枚系票 123 已知腿） |
| **MUT-L（拆掉解析层）** | 软链 | tools `rc=1 RUN=79 PASS=41 FAIL=21 SKIP=3 SUBPASS=12 SUBFAIL=2`；llm `rc=1 RUN=72 PASS=52 FAIL=11 SKIP=0 SUBPASS=3 SUBFAIL=6`；perm `rc=1 RUN=14 PASS=9 FAIL=5 SKIP=0`；approval `rc=1 RUN=49 PASS=28 FAIL=1 SKIP=1 SUBPASS=19 SUBFAIL=0` | **43 枚回归**，逐名与改前基线 PRE-L **四包全部 IDENTICAL**（`tools/llm/perm/agent_approval: IDENTICAL red name sets`，23/17/5/1 枚）；机制字串回到 `not provably resolved`=23、`refusing to seal`=23、`审批未通过`=4、`审批通道尚未接入`=4，样例逐字：`probe_health_test.go:151: memory: create data dir: winsec: refusing to seal /varlink/w124tmp/TestProbeSuiteMeasuresBrokenFC3118190525/002/data: the installed risk.c26Pipeline answered "/varlink/…", a spelling the floor itself refuses` |
| **MUT-P（同一发变异）** | 普通 | tools `rc=0 RUN=79 PASS=62 FAIL=0 SKIP=3 SUBPASS=14 SUBFAIL=0`；llm `rc=0 72/63/0/0`+`9/0`；perm `rc=0 14/14/0/0`；approval `rc=0 49/29/0/1`+`19/0` | 无（`FAIL` 行四包 = 0；四数与 PRE-P/POST-P **逐数相同** ⇒ 这层解析只在软链形起作用，普通形是恒等操作） |

⇒ 红名**点到用例自己**（43 枚逐枚点名 + 3 枚已知腿一并回归，不是包级 `FAIL`），且变异前 `grep` 证落地、`go build` / `go vet` rc=0 先于读数。这一发**做了**。

## 7. 门禁读数（票面 AC#5）

待量。

## 8. 改动面与本批用的 helper

- 用的既有 helper = **`proc.SealableRoot`**（`internal/proc/envfork.go`），即票 119 生产路 `TestDataDir` / `DefaultLayout` / `cmd/wisp` 的 `resolveDataDir` 走的那一枚；批次 1 同法。
- 新增五枚**单行委托**（`tools`、`tools_test`、`llm_test`、`perm`、`approval_test` 各一枚；tools 因内/外两枚测试包作用域不同而拆两枚，**不是**第二份解析实现）+ 一枚 `sealableTempCanonical124`（= `mustCanonical(sealableTempDir124(t))`，被测试的规范化器照旧跑，变的只是喂进去的根拼写）。零新解析路数：`R-125-2` 那本副本账不增行。
- 递根点 30 处（tools 22 行 / llm 1 / perm 6 / approval 1），对应 43 枚（llm 的 1 处 fixture 覆盖 17 枚；tools 的 fs_write 3 处覆盖 4 枚含两子测试）。
- 未动：`internal/tools/wiring_test.go`（票 123 两枚）、`internal/tools/paths_ticket107b_probes_test.go` 的 ProbeA/ProbeB 与 `ticket107_portable` 非红枚、`fs_write_test.go`/`paths_workspace_test.go`/`pathshape_portable_test.go` 里服务非红用例的 `tempRaw`/`t.TempDir()` 站点、`internal/winsec/**`、`internal/risk/**`、`internal/memory/**`、`internal/config/**`、`cmd/wisp/**`。

## 9. 未验证项 + `next=`

待量后写。

## 10. 临时件清单（只建不删）

`/d/tmp/wisp124-2b2-{pre,post}/`、`/d/tmp/wisp124-2b2-logs/`、`/d/tmp/wisp124-2b2-h.sh`、`/d/tmp/wisp124-2b2-convert.py`（探针/变异台目录待量后登记）。
