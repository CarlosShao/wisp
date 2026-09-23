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

## 7. 门禁读数（票面 AC#5）— 本接续方 `worker-ticket124-ac2b-2-r2` 补（2026-09-23 19:4x–20:3x +8，即 `date -u` 11:4x–12:3x）

**测的是哪棵树 / 共树漂移**：开工首读 `git rev-parse --short HEAD` = **`e113b1a`**（与派单给的起点一致，本方一切读数锚在这枚）。测量途中 HEAD 先后漂到 `b45d74d`、`00d3f34`（兄弟票 133 的 docs 与 A128 台账），已核 `git diff --name-only e113b1a HEAD -- internal/ cmd/ go.mod go.sum` **零 hunk**、`git diff --stat dfa3dc4 e113b1a -- internal/` **空** ⇒ 本方量的四包与 §0 那台 POST 快照（`dfa3dc4`）是同一版本码；`e113b1a` 相对 `dfa3dc4` 唯一代码差在 `cmd/wisp/leg_dispatch_gate_133_test.go`（票 133 地界，非本批）。工作树里别人的半成品一个字都不读。
- 快照：`/d/tmp/wisp124-2b2-r2-post` = `git archive e113b1a | tar -x`（**1048 文件**；`dfa3dc4` 那台是 1046，+2 系兄弟票 docs/报告件）。对照快照 `/d/tmp/wisp124-2b2-r2-pre` = `git archive 5417a3c`（1039 文件，§0 的 pre 锚，只用于 §7.4 的归因那一发）。
- 容器 `golang:1.27`（`go1.27.1 linux/amd64`、`CGO_ENABLED=1`），复用票 119 的命名卷 `ac119-gomodcache`/`ac119-gocache` ⇒ 离线可编；挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`；每枚样本进容器第一件事 `ls -l /src/go.mod` = `-rwxrwxrwx 1 root root 883`、`md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum /src/internal/winsec/resolve.go` = `b6876a5efe759f6e17434d1b50a129c3` ⇒ 与 §0 逐字同字、非空挂自证。
- 跑法（新建，未改前任任何脚本）：`/d/tmp/wisp124-2b2-r2-gate.sh`（`bash /gate.sh <link|plain|none> <tag>`，`COUNT=2`、`-timeout 60m`）与 `/d/tmp/wisp124-2b2-r2-vet.sh`（容器内 gofmt/vet）。**每形一枚新容器**：软链形容器 `abd20fa8ef62`（19:57:49→20:28:39 +8）、普通形容器 `3e89bc848b1f`（20:28:42→20:29:45 +8）。形状硬断言沿用 §0（`exit 96/97/98/99` 两发均未命中），原文读数：软链形 `SHAPE=link TMPDIR=/varlink/w124tmp resolved=/realpriv/w124tmp` + `lrwxrwxrwx 1 root root 9 ... /varlink -> /realpriv`；普通形 `SHAPE=plain TMPDIR=/plainroot/w124tmp resolved=/plainroot/w124tmp`（`/varlink` 根本不存在、`/plainroot` 是真目录）。日志：`/d/tmp/wisp124-2b2-r2-logs-link/`、`/d/tmp/wisp124-2b2-r2-logs-plain/`、`/d/tmp/wisp124-2b2-r2-logs/`（临时件，只建不删）。
- ⚠ 为什么必须 60m 且每形新容器：见 §9 未验证项 4——前任 §7 待量期间遗留的那发（`/d/tmp/wisp124-2b2-logs/POST-L2.*`，19:19–19:21）用 `-timeout 25m` ⇒ `panic: test timed out after 25m0s`、四数**被截断**；同一次容器复用又让普通形那一发撞 `exit 99`（`RUN-COUNT2.out.txt` 第二段）⇒ 普通形 `-count=2` 前任**没有读数**。那是补空格，不是复跑 §2-§6；本方这一发 `test timed out` 命中 **0**。
- 开测前查在飞（本机就是 self-hosted runner）：19:4x 见 `35856629513`（`docs(A128)`，11:47:56z 起）在跑 ⇒ 不计时的静态读数（7.2/7.3/7.4 容器那一发）先跑，**计时敏感的 `-count=2` 两发等它 `completed`（11:57:29z；`test-core=success`、`test-windows=failure`）之后才起**，期间宿主零编译动作。

### 7.1 `-count=2 -v` 四数（四包 × 两形；只能从 `-v` 量，`-count=2` 才不缓存）

| 台 | 形状 | `internal/tools` | `internal/llm` | `internal/perm` | `internal/agent/approval` |
|---|---|---|---|---|---|
| `GATE-L2` | 软链 | `rc=1 RUN=158 PASS=118 FAIL=6 SKIP=6 SUBPASS=28 SUBFAIL=0 SUBSKIP=0`（包级 `1806.869 s`） | `rc=0 RUN=144 PASS=126 FAIL=0 SKIP=0 SUBPASS=18 SUBFAIL=0`（`37.009 s`） | `rc=0 RUN=28 PASS=28 FAIL=0 SKIP=0`（`0.253 s`） | `rc=0 RUN=98 PASS=58 FAIL=0 SKIP=2 SUBPASS=38 SUBFAIL=0`（`0.560 s`） |
| `GATE-P2` | 普通 | `rc=0 RUN=158 PASS=124 FAIL=0 SKIP=6 SUBPASS=28 SUBFAIL=0`（`18.707 s`） | `rc=0 RUN=144 PASS=126 FAIL=0 SKIP=0 SUBPASS=18 SUBFAIL=0`（`38.881 s`） | `rc=0 RUN=28 PASS=28 FAIL=0 SKIP=0`（`0.298 s`） | `rc=0 RUN=98 PASS=58 FAIL=0 SKIP=2 SUBPASS=38 SUBFAIL=0`（`0.561 s`） |

⇒ **每包八个数都是 §2 那一发 `-count=1` 的正好 2 倍**（tools 软链 `79/59/3/3`+`14/0` → `158/118/6/6`+`28/0`、tools 普通 `79/62/0/3`+`14/0` → `158/124/0/6`+`28/0`；llm `72/63/0/0`+`9/0` → `144/126/0/0`+`18/0`；perm `14/14/0/0` → `28/28/0/0`；approval `49/29/0/1`+`19/0` → `98/58/0/2`+`38/0`）⇒ 无缓存读数、无 flake，**与前任 §2/§4 的读数无分歧**（本方未复跑 §2-§6，也没有改它们一行）。
⇒ **红名逐名**（软链形 6 行 = 票 123 那三枚 300 s 腿 × 2 counts）：`TestL1WriteGoesThroughTheRealBlockWindow`（`300.00s` / `300.01s`）、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`（`300.02s` / `300.03s`）、`TestFSReadOnlyNeverOpensACard`（`300.03s` / `300.04s`）⇒ **零枚来自本批 43 枚名册**；tools 包级 `rc` 仍为 `1` 就是这三枚 ×2 造成的，其余三包 `rc=0`。普通形四包 `FAIL=0`、`SUBFAIL=0`。
⇒ **SKIP 名册两形逐名同一组 4 枚 ×2 counts**（tools `TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`、`TestWorkspaceSwitchRefusesAJunctionToOutside` + approval `TestDefaultDeadlineWallClockMeasurement`），与 §4"形状无关、改前改后一字不变"吻合；两形 `RUN` 逐包相等（`158/144/28/98` 对 `158/144/28/98`）⇒ 门禁这一发同样没有"被跳过冒充变绿"。
⇒ 机制字串计数（本方这一发四包日志，四包合计的行计数）：`not provably resolved` **0**、`refusing to seal` **0**、`审批未通过` **0**、`审批通道尚未接入` **0**（普通形侧亦全 0）。
⚠ 计时账：软链形 tools 单包 `1806.869 s`，六枚 300 s 腿自身就占 1800 s ⇒ `-count=2` 下包级 `-timeout` **≥ 50m 是硬要求**（25m 必然截断，前任那发的截断读数已按名保留在 `/d/tmp/wisp124-2b2-logs/POST-L2.summary.txt`，本方不改它）。

### 7.2 `gofmt -l`（票面 AC#5 的"整包"：四包各自一发；宿主 + 容器各一组）

| 调用（均在锚点纯净快照里） | 结果 |
|---|---|
| 宿主 `gofmt -l internal/tools`、`internal/llm`、`internal/perm`、`internal/agent/approval` | **各 0 行**，`rc=0`（`/d/tmp/wisp124-2b2-r2-logs/GOFMT-HOST.txt`） |
| 宿主 `gofmt -l internal cmd`（全树参照，非本格要求） | **0 行** |
| 容器 `gofmt -l` 四包各自 + `gofmt -l internal cmd` | **各 0 行 / 0 行**（`/d/tmp/wisp124-2b2-r2-logs/VET-CONTAINER-RUN.txt` 第三段） |

### 7.3 `"$(go env GOPATH)/bin/gofumpt.exe" -l`（票面 AC#5：真跑；写「未跑」必须引错误原文）

版本行 **`v0.7.0 (go1.27.1)`** ⇒ 工具存在，本格**不是**「未跑」，没有错误原文要引。宿主对四包各自 `gofumpt -l <pkg>` ⇒ **各 0 行**（`/d/tmp/wisp124-2b2-r2-logs/GOFUMPT-HOST.txt`）。

### 7.4 `go vet` 双 GOOS（逐错误行归因；宿主交叉那一发既不算破口也不算清白）

| 调用形状 | rc | 输出 |
|---|---|---|
| `GOOS=windows go vet ./...`（宿主原生、全树、锚点快照） | **0** | 0 行（`VET-WINDOWS-HOST-ANCHOR.txt`） |
| `GOOS=windows go vet ./internal/tools/`、`./internal/llm/`、`./internal/perm/`、`./internal/agent/approval/` | **0** ×4 | 0 行（`VET-WINDOWS-HOST-PKGS.txt`） |
| `GOOS=linux go vet ./...`（**容器原生**、`CGO_ENABLED=1`、`go1.27.1 linux/amd64`、同一快照） | **0** | 0 行（`VET-LINUX-NATIVE-ALL.txt`；原文贴在 `VET-CONTAINER-RUN.txt` 末段）⇒ **这一发才是 linux 的真类型读数**；四包单跑亦各 `rc=0`、`VET-LINUX-NATIVE-PKG{1,2,3,4}.txt` 全空 |
| `GOOS=linux CGO_ENABLED=0 go vet ./...`（**宿主交叉**、全树） | **1** | 4 行，逐字见下 |
| `GOOS=linux CGO_ENABLED=0 go vet <pkg>`（四包各自，宿主交叉） | **0** ×4 | 0 行（`VET-LINUX-HOST-CROSS-PKGS.txt`）⇒ 本批四包在 linux 目标下过 vet |
| `GOOS=linux go vet ./cmd/balldebug/`（宿主交叉、单跑） | **1** | `package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in <tree>\cmd\balldebug`（该 main 全文件带 `//go:build windows`）；同一发在 pre 锚 `5417a3c` 快照里**只差路径**、其余逐字相同（`VET-LINUX-HOST-CROSS-BALLDEBUG-wisp124-2b2-r2-{post,pre}.txt`） |
| `GOOS=linux go vet $(go list ./... \| grep -vE 'cmd/(wisp\|balldebug)$')`（剔掉那两枚包，余 30 枚） | **0** | 0 行（`VET-LINUX-HOST-CROSS-EXCL.txt`；同一发里 `go list ./...` 列得 32 枚包）⇒ 交叉 `rc=1` **全部**由 `cmd/wisp`、`cmd/balldebug` 两枚既有形状贡献 |

宿主交叉全树那一发的全部输出（`VET-LINUX-HOST-CROSS-ANCHOR.txt` 逐字）：

```
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
```

⇒ 死在**包加载**阶段（外部模块 `sherpa-onnx-go-linux@v1.13.8` 的 build constraints），结构上到不了类型检查 ⇒ 按票面纪律**既不算破口也不算清白**；清白由上表容器原生那发 `rc=0` 提供。归因：同一支命令在同一台机器上对 pre 锚 `5417a3c` 的快照输出**逐字节相同**（`diff` 前 3 行 rc=0，`VET-LINUX-HOST-CROSS-PRE5417.txt` vs `VET-LINUX-HOST-CROSS-ANCHOR.txt`）⇒ 与本批无关的既有形状；错误文本点到的是 `cmd/wisp`（2b-4 地界）与一枚外部模块目录，**零行指向本批四包**。
⚠ 与**批次 1** §7 的交叉读数枚数不同（那发报 2 处：`cmd/wisp` 链 + `cmd/balldebug`；本方 `./...` 这发只报 1 处）——`go vet ./...` 在包加载阶段遇到第一枚不可加载的包即中止，报几处取决于加载顺序；把 `cmd/balldebug` 单独点名时它照样红（本方已逐字量，见上表倒数第二行）。⇒ 这是**调用形状不同**，不是同一支命令的读数分歧；两发都不构成对本批四包的判据。

### 7.5 一次全仓仪器：`sh scripts/d22scan.sh`（作用域是整个 `internal/`；按包门禁结构性看不见它）

调用形状（唯一受支持的形状；票 67 AC#2 把它做成脚本，就是为了堵"从仓根跑独立 module"与"默认 `-root .` 走空树"那两种假成功——本方**没有**撞到 `main module does not contain package`）：在**纯净快照**里 `cd /d/tmp/wisp124-2b2-r2-post && sh scripts/d22scan.sh`（宿主 Git Bash）。脚本自己从所在位置推导仓根，两步依次硬跑（`set -eu`、无跳步、无 `|| true`）：① `sh tools/d22scan/runtests.sh -C tools/d22scan ./...`（播种违规的正向对照）② `go run . -root <推导出的仓根>`（真扫）。

- **rc=0**（全量原文 `/d/tmp/wisp124-2b2-r2-logs/D22-ANCHOR.txt`；脚本推导出的根逐字为 `d22scan.sh: scan of /d/tmp/wisp124-2b2-r2-post` ⇒ 扫的是快照，不是脏工作树）
- 正向对照那一步：`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` ⇒ 这把尺**能红**，"clean"不是"眼睛瞎"
- 真扫规模：`examined 225 production Go files under internal/ and cmd/ of D:/tmp/wisp124-2b2-r2-post`
- 台账八 scope 逐数对照（本方 `e113b1a` vs 批次 1 §7 在 `fafe2b4` 那发记的数）：

| scope | 批次 1 §7 | 本方（`e113b1a`） | 变化与出处 |
|---|---|---|---|
| bans #1-5 `internal/` | 203 | 203 | 0 |
| bans #1-5 `cmd/` | 22 | 22 | 0 |
| ban #6 `frontend/` | 40 | 40 | 0 |
| ban #7 `internal/tools/` | 18 | 18 | 0 |
| ban #8 `design/` | 16 | 16 | 0 |
| ban #8 `frontend/` | 40 | 40 | 0 |
| ban #8 `internal/` | 392 | **397** | **+5 ＝ 本批五枚单行委托 `_test.go` 进树**（`git diff --name-status fafe2b4 e113b1a -- internal/` 里 `^A` 恰为这 5 枚，逐名可点） |
| ban #8 `cmd/` | 37 | **38** | **+1 ＝ 票 133 新增的 `cmd/wisp/leg_dispatch_gate_133_test.go`**（同一命令 `^A` 清单里唯一一枚 cmd 文件，**非本批**） |

⇒ **八数一枚都不降**；唯二变化都是"多扫了一枚文件"而非"少扫一枚"，且各有可核出处。

**AC#5 小结**：四包两形 `-count=2` 四数齐（软链形红名只剩票 123 那三枚 ×2、普通形零红）；`gofmt -l` 整包四包 × 宿主/容器全 0 行；`gofumpt v0.7.0` 真跑、0 行；`go vet` 双 GOOS——windows 全树 `rc=0`、linux **容器原生**全树 `rc=0`（真类型读数）、linux 宿主交叉 `rc=1` 已逐行归因到两枚与本批无关的既有形状（剔掉即 `rc=0`）；`sh scripts/d22scan.sh` 纯净快照 `rc=0`、正向对照 21 PASS/0 FAIL/0 SKIP、台账八 scope 不降。**本格不为四包翻 AC#2 那一格**（票面 16:33 已定：终判据留到 2b-4 交回后一次性复算）。

## 8. 改动面与本批用的 helper

- 用的既有 helper = **`proc.SealableRoot`**（`internal/proc/envfork.go`），即票 119 生产路 `TestDataDir` / `DefaultLayout` / `cmd/wisp` 的 `resolveDataDir` 走的那一枚；批次 1 同法。
- 新增五枚**单行委托**（`tools`、`tools_test`、`llm_test`、`perm`、`approval_test` 各一枚；tools 因内/外两枚测试包作用域不同而拆两枚，**不是**第二份解析实现）+ 一枚 `sealableTempCanonical124`（= `mustCanonical(sealableTempDir124(t))`，被测试的规范化器照旧跑，变的只是喂进去的根拼写）。零新解析路数：`R-125-2` 那本副本账不增行。
- 递根点 30 处（tools 22 行 / llm 1 / perm 6 / approval 1），对应 43 枚（llm 的 1 处 fixture 覆盖 17 枚；tools 的 fs_write 3 处覆盖 4 枚含两子测试）。
- 未动：`internal/tools/wiring_test.go`（票 123 两枚）、`internal/tools/paths_ticket107b_probes_test.go` 的 ProbeA/ProbeB 与 `ticket107_portable` 非红枚、`fs_write_test.go`/`paths_workspace_test.go`/`pathshape_portable_test.go` 里服务非红用例的 `tempRaw`/`t.TempDir()` 站点、`internal/winsec/**`、`internal/risk/**`、`internal/memory/**`、`internal/config/**`、`cmd/wisp/**`。

## 9. 未验证项 + `next=`（本接续方 `worker-ticket124-ac2b-2-r2` 补，2026-09-23 20:3x +8）

**未验证项（照实列，不含"我相信"）**：

1. **macOS 那一半仍然没有实测。** 仓库只注册一枚 self-hosted runner（`gh api repos/{owner}/{repo}/actions/runners` ⇒ `wisp-selfhosted-01`，`os=Windows`，`status=online`），没有 macOS runner。而票面《事实》第 3 条明写这一族"代表的是 **macOS 的真实形状**（`TMPDIR` 在 `/var` 之下，而 `/var` 是符号链接）"⇒ 本批（前任 §2-§6 ＋ 本方 §7）全部读数都是 **Linux 容器 `golang:1.27` + `ln -s /realpriv /varlink` 的代理形状**，不是 macOS 真机。"这 43 枚在 macOS 上会同样转绿"是**外推**。与 AC#2a §7 未验证项 1、票 119 `R-119-8` 同账。
2. **CI 两腿对本票判据零信号（不等于"CI 清白"）。** 本批四包都在 CI 清单里（`scripts/portable-tests.sh:173-179` 的 `--scope=core` 逐条含 `./internal/agent/ ./internal/agent/approval/`、`./internal/llm/...`、`./internal/tools/...`、`./internal/perm/`），但 ubuntu 腿 `/tmp` 是真目录、windows 腿 `%TEMP%` 不带链接 ⇒ 软链形在 CI 上**没有分母**，这一族在 CI 恒绿。实测账：近 6 枚 push 触发的 `ci` run（`35840958334`/`35843139130`/`35845195469`/`35848845645`/`35848981451`/`35856629513`）`test-core` **6/6 `success`**、`test-windows` **6/6 `failure`**；其中三枚的 `headSha`（`4a23de5`/`3d43c3f`/`b45d74d`）已用 `git merge-base --is-ancestor dfa3dc4 <sha>` 逐枚核为**带着本批改件码**。本方逐 job 读的那枚（`35856629513`）红在 `lint`（`staticcheck` 一步）与 `test-windows`（逐名只两枚包：`cmd/wisp`、`internal/risk`），**无一枚属本批四包**。⚠ `e113b1a` 已随 `b45d74d` 那次 push 进了两枚远端（`git branch -r --contains e113b1a` ⇒ `origin/dev`、`cnb/dev`），但一次 push 只跑尖端 ⇒ CI 上**不存在**以 `e113b1a` 为 `head_sha` 的 run（最新一枚 `35859261431`，`headSha=03f87f4`，12:14:38z，整体 `failure`）；⇒ "CI 看不见这一族"成立的原因是**形状不在 CI 的任何一台 runner 上**，不是"没推上去"。
3. **票 123 的 300 s 腿本批一枚都没修，§7 里它们仍是 6 行红。** `TestL1WriteGoesThroughTheRealBlockWindow`、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`（票面点名的两枚）、`TestFSReadOnlyNeverOpensACard`（§1 新登记的同族第三枚）× `-count=2` ⇒ 各 `300.0x s`。⇒ 已知遗留，不是"43 枚没转绿"，也不是门禁破口；**归口票 123**。
4. **`-count=2` 的 tools 软链形是计时读数，包级 `-timeout` ≥ 50m 是硬要求。** 前任 §7 待量期间遗留那一发（`/d/tmp/wisp124-2b2-logs/POST-L2.*`，19:19–19:21，`-timeout 25m`）以 `panic: test timed out after 25m0s` 收（该日志 `:390`），四数**被截断**：`RUN=156`（而非 `2×79=158`）、`FAIL=4`（而非 6）、`PASS=117`（而非 118）⇒ 那是**分母缩小**、不是变绿，本方不改它、只按名保留。同一次容器复用还让普通形那发撞了 `exit 99`（`RUN-COUNT2.out.txt` 第二段：`SHAPE_ASSERT_FAILED plain form must not have /varlink at all`）⇒ **普通形 `-count=2` 前任没有读数**。本方 §7 用 60m + 每形一枚新容器补上（`test timed out` 命中 0）。⚠ 这是补 §7 的空格；**§2-§6 的读数一枚未复跑、一行未改**，两形四数与 §2 逐数吻合（正好 2 倍）⇒ 无分歧。
5. **宿主（Windows）侧没有四包的 `-count=2` 四数**：AC#5 的四数在容器两形取；Windows 的 `%TEMP%` 不属本票形状，未量（票面亦未要求）。同理宿主那一发 `GOOS=linux go vet` 是**交叉**，结构上到不了类型检查（§7.4 已逐行归因），真读数在容器那发。
6. **全树终判据未复算**（票面 16:33 那条"软链形红名数＝0"）——票面已定"留到 2b-4 交回后一次性复算"⇒ **AC#2 本格不翻**、AC#2b 各批格也不由本方翻。
7. **判据④（探针台）与判据⑤（变异台）的原始台子本方未复跑、未复核**（派单禁令）。§10 的临时件清单里前任仍写着"（探针/变异台目录待量后登记）"——本方**不动那一节**（它不归本接续方），只按名保留盘上目录 `/d/tmp/wisp124-2b2-{probe,mut}`。

**`next=`（写给 2b-3 与 2b-4；交编排者）**

- **2b-3＝`internal/agent` 26 ＋ `internal/winsec` 17 = 43 枚。**
  - ⚠⚠ **`internal/winsec` 那 17 枚不许照抄本批这枚委托**（批次 1 的 `next=` 点过；本方把它三条实据逐字核了一遍，全部成立）：
    ① 票 119 面 `R-119-9`（`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md:73`）原文——"**不许拿被测函数算 fixture。** 一枚用例的期望值只能来自文件系统或调用方自己声明的字面值，不能来自被测函数（含其幂等组合）：`SealableRoot(base)` 当期望值 ⇒ 把'声明的树原样返回'这条纪律抹掉也照样绿"；
    ② 这条纪律**已经写进 winsec 的码**：`internal/winsec/dataroot_symlink_119_other_test.go:108-111` 逐字——`cleanSpelling119` "It is **deliberately filepath.EvalSymlinks and not proc.SealableRoot**: the latter is what several of these cases exist to test, and an expectation computed by the function under test cannot fail."（同形还有 `placement_leaf_118_other_test.go:37` 的 `leafLinkTo118` 与 `seam_probe_root_125_other_test.go:84`）；
    ③ `internal/proc/envfork.go:148` 的边界话是 "…never case-folds, and **nothing in internal/winsec calls it**" ⇒ 一旦 2b-3 让 winsec 的测试去调 `proc.SealableRoot`，这句注释当场腐坏（改注释＝动被冻结的边界话，属于要报的偏离）。
    ⇒ 那一包**要么**沿用票 125 的先例、自带走法（只用文件系统自己给的答案），**要么停下报编排者**。别为了"四批统一形状"把 R-119-9 那一族恒真用例做回来。
  - ⚠ winsec（与 config）的 POSIX 用例全在 `*_other_test.go`（`//go:build !windows`）⇒ 在宿主上**连编译都不参与**；改完只在容器里读，别拿宿主 `go build ./...` 的绿当清白。
  - ⚠ 两形四数要预留**形状自带 SKIP**：AC#1 附带发现 1 记录软链形有 4 枚从"跑"变成"SKIP"（`internal/winsec` 3 枚 `TestAC2POSIX…125` seam 探针 + `internal/config` 1 枚，后者已由批次 1 §4 逐名解释）。2b-3 若见 winsec 两形 `RUN`/`SKIP` 不等（AC#1 那本是 `45→52`），那是**形状副作用**、不是本批破口，但必须逐名写进解释，不许算进"归零"也不许算进"红"。
  - `internal/agent` 那 26 枚可沿用本批形状，且**只需一枚单行委托**：涉事 7 枚文件（`spill_path_invariant`/`spill_name_injectivity`/`spill_test`/`forensics`/`guard`/`truncation`/`loop_golden`，本方逐枚 `grep '^package '` 核过）全在 `package agent` 内 ⇒ 不存在 tools 那种内/外双包问题。账 §5.2 已写明 26 枚的红只有两种串、且全部落在 setup（`memory.Open` / `NewSpiller`+`Prepare`）⇒ 递根点数应远小于枚数（照 llm 那本账"1 处 fixture 覆盖 17 枚"预期）。
  - ⚠ 跑法账：把 agent + winsec 并进软链形时，包级 `-timeout` 按**包内最长腿 ×2** 抬（本批 tools 软链形 `-count=1` 就要 `903.5 s`、`-count=2` 要 `1806.9 s`）；判据③要两形各一枚，且**每形一枚新容器**（§9 未验证项 4 那一发就是复用容器撞的 `exit 99`）。
- **2b-4＝`cmd/wisp` 5 枚。**
  - 名册（账 §5.9 逐枚）：4 枚顶层 `providers_test.go:133`/`:179`/`:208`/`:228`，全部经同一枚 `newProvidersFixture(t)` ⇒ **一处递根覆盖 4 枚**；另 1 枚子测试 `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`，断言原文在 `secret_test.go:679`，软链形实拿 `secret: create /varlink/…: refusing to seal … not provably resolved`。
  - ⚠ **同包另有 28 枚两形都红**（逐名清单 `docs/evidence/s1/124-ac1-denominator-readings.md:931`；那批错误原文里 `DPAPI is only available on Windows` 命中 20 处）⇒ 与本票形状**无关**（普通形同红），归票 119 `R-119-7` 那本账。2b-4 一枚都不许顺手"修绿"；反过来也不许因为包级 `FAIL≠0` 就把自己那 5 枚判成没转绿：**判据是逐枚点名，不是包级 rc**。
  - ⚠ 地界是**文件级**：`cmd/wisp/**` 兄弟票正在动——`cmd/wisp/leg_dispatch_gate_133_test.go` 系票 133 在 `fafe2b4..dfa3dc4` 之间新增（§7.5 台账 ban #8 `cmd/` 37→38 就是它，非本批）。票面明写 2b-4 仍排在 `acceptor-ticket131-r3` 交回之后；动码前先按名核清哪枚文件有谁在飞。
- **全树终判据怎么复算（2b-4 交回后一次性；谁复算谁照这六步，别凭记忆）**：
  1. 取当时 HEAD 的 `git archive <sha> | tar -x` 纯净快照（**禁读脏工作树**），锚 sha 写进证据；并核 `git diff --stat <改件sha> <锚sha> -- internal/ cmd/ go.mod go.sum` 的 hunk 各归谁（本方这一发就是这么把 `cmd/wisp` 那一枚漂移甩给票 133 的）；
  2. 容器 `golang:1.27`；软链形 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp`，`exit 96/97/98` 断言照用；普通形**换一枚新容器** + `exit 99`（`/varlink` 必须根本不存在、`/plainroot` 不许是链接）；每枚样本进容器先 `ls -l /src/go.mod`（883 字节）＋ `md5sum`（`go.mod`=`f6ef661732b1851e5c3db348113cb605`、`internal/winsec/resolve.go`=`b6876a5efe759f6e17434d1b50a129c3`）证明确实挂上、不是空挂（Git Bash 下 `docker run -v "C:\…"` 会**静默挂空且 rc=0**）；
  3. 逐包 `-count=1 -v` 跑 **AC#1 那本全树分母**（`internal/**` + `cmd/**` 全部包，AC#1 的 log 记的是"全部 30 枚 `internal/**` + `cmd/**` 包"），从 `-v` 收软链形 `--- FAIL`（顶层 + 子测试）逐名；
  4. 与名册做差集：账 `124-ac2a-leg-classification.md` §5 那 **131 枚"可转"** ⇒ **期望差 = 空集**，这就是票面 16:33 那句"软链形红名数＝0（本票范围内）"；
  5. 允许的余红**逐名钉死**、不按枚数四舍五入：(a) 票 123 那族 300 s 审批超时腿 **三枚**——`TestL1WriteGoesThroughTheRealBlockWindow`、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`、`TestFSReadOnlyNeverOpensACard`（票面 16:33 只点了前两枚；第三枚是本票 §1 新登记的同族腿，**复算时三枚一起豁免，否则会得到"多一枚红"的假破口**）；(b) 两形都红的 29 枚（`internal/panel/` 1 ＋ `cmd/wisp/` 28，`R-119-7` 那本账）；(c) 形状自带的 SKIP（winsec 3 ＋ config 1）——**SKIP 既不是红也不是绿**，必须逐名写在读数旁边；
  6. 结论必须**两形并列**、`RUN`/`SKIP` 差逐包解释；"归零"不许拿包级 `ok` 冒充逐名点名；判据不成立（差集非空、或多一枚少一枚）就**报回来**，不许改数、不许放宽断言。

## 10. 临时件清单（只建不删）

`/d/tmp/wisp124-2b2-{pre,post}/`、`/d/tmp/wisp124-2b2-logs/`、`/d/tmp/wisp124-2b2-h.sh`、`/d/tmp/wisp124-2b2-convert.py`（探针/变异台目录待量后登记）。
