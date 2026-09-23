# 票 124 AC#2b 批次 3a（`worker-ticket124-ac2b-3a`）— `internal/agent` 26 枚接上票 119 那条纪律

**日期**：2026-09-23（测量段 20:51–21:10 +8，即 `date -u` 12:51–13:10）· **性质**：转换方交件（改 `.go`，只改 `_test.go`）
**本批范围**：清点账 `docs/evidence/s1/124-ac2a-leg-classification.md` §5.2（行 18-43，共 26 枚）＝`internal/agent` 一整包，账上全标「可转」。
批次 3 的后一半（`internal/winsec` 17 枚）**不在本批**——那是并行代理 `worker-ticket124-ac2b-3b` 的地界，本程 `internal/winsec/**` 零 hunk。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2b 的放行与分批」（09-23 16:33 编排者）那张表 + 五条共用结案判据；批次 2 的 `next=` 给这一批的原话是「`internal/agent` 的 26 枚可沿用本批形状，且只需一枚委托」。

## 0. 锚点 / 快照 / 跑法（可复核）

- 开工首读 `git rev-parse --short HEAD` = **`4ea0db2`**（下文记 **pre 锚**；全部改前读数在它上面量）。分支 `dev`，开工时工作树除一枚来源未明的未跟踪文件（`docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`，本程不提交、不改、不删、不据它开票）外**零脏项** ⇒ 首读的 HEAD 与被验版本同一。
- 改件 commit：**`62dda11`**（`test(124,AC#2b-3a)`，8 枚文件、+44/-18；删除侧 18 行全是递根点，见 §3）。
- 共树漂移：测量途中 HEAD 先后经 `1bf88db`（A130 台账）、`857a5fe`（3b 的分类文档）到本方读数收口时。已核 `git diff --name-only 4ea0db2 HEAD -- internal/ cmd/` 里非 `internal/agent/` 命中 **0 枚**（全部漂移都在 `docs/**`）⇒ §2 那四台量的是同一版本码，pre 读数是改前同一版本。
- 全部测量都在 `git archive <sha> | tar -x` 的纯净快照里做；脏工作树一个字都不读。

| 快照目录 | 内容 | 文件数 |
|---|---|---|
| `/d/tmp/wisp124-2b3a-pre` | `git archive 4ea0db2`（改前基线） | 1049 |
| `/d/tmp/wisp124-2b3a-post` | `git archive 62dda11`（改后） | 1050（+1＝本批那枚单行委托 `_test.go`） |
| `/d/tmp/wisp124-2b3a-probe` | 改后 + 判据④取字串探针（9 处插入、只 `t.Logf`，**未入库**） | 1050 |
| `/d/tmp/wisp124-2b3a-mut` | 改后 + 判据⑤变异（把委托改回原样，**未入库**） | 1050 |

- 容器 `golang:1.27`（`go1.27.1 linux/amd64`、`CGO_ENABLED=1`），复用票 119 的命名卷 `ac119-gomodcache` / `ac119-gocache` ⇒ 离线可编（每发都带 `-e GOPROXY=off`）。
- 挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`，`/src` 以 `:ro` 挂；**每枚样本进容器第一件事**打 `ls -l /src/go.mod` = `-rwxrwxrwx 1 root root 883`、`md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum /src/internal/winsec/resolve.go` = `b6876a5efe759f6e17434d1b50a129c3` ⇒ 与 AC#1／AC#2a／批次 1／批次 2 逐字同字，非空挂自证（Git Bash 下 `docker run -v "C:\…"` 会静默挂空且 rc=0＝假绿）。
- 形状硬断言沿用批次 1／2：软链形 `exit 96/97/98`（链接必须真是链接、`readlink -f /varlink/w124tmp` 必须等于 `/realpriv/w124tmp`）；普通形 `exit 99`（`/varlink` 必须根本不存在、`/plainroot` 不许是链接）。**本批 9 发无一命中**。软链形每次打出的形状行逐字：`SHAPE=link TMPDIR=/varlink/w124tmp resolved=/realpriv/w124tmp` + `lrwxrwxrwx 1 root root 9 ... /varlink -> /realpriv`；普通形：`SHAPE=plain TMPDIR=/plainroot/w124tmp resolved=/plainroot/w124tmp`。
- **每形一枚新容器**（批次 2 §9 未验证项 4 那本账：复用容器会让普通形撞 `exit 99`）。本批 **9 发 `-v` 读数**（PRE-L／PRE-P／POST-L／POST-P／PROBE-L／MUT-L／MUT-P／GATE2-LINK／GATE2-PLAIN）各一枚新容器，另有 3 发非测试容器（变异落地自证 2 发、linux 原生 `go vet` 1 发）。容器名逐枚：`wisp124-2b3a-{prelink,preplain,postlink,postplain,probe-link,mutlink,mutplain,gatelink,gateplain,mutproof,mutproof2,vetlinux}`（全部 `--rm`，跑完即退）。
- 跑法脚本（新建，未改批次 1／2 那两支）：`/d/tmp/wisp124-2b3a-h.sh`（`bash /h3a.sh <link|plain|none> <tag> [pkgs]`，`COUNT=2` 切门禁那一发、`TIMEOUT` 默认 30m）；转换件 `/d/tmp/wisp124-2b3a-convert.py`（18 行行号锚定替换，锚不中即 `ANCHOR MISS` 且不写盘）；探针件 `/d/tmp/wisp124-2b3a-probe-patch.py`；变异件 `/d/tmp/wisp124-2b3a-mutate.py`。日志与逐名表全在 `/d/tmp/wisp124-2b3a-logs/`。
- 包级 `-timeout`：本包最长腿是 `TestSpillThroughLoop`（两形都 < 0.02 s），`30m` 远超需要；**`test timed out` 在本批 7 发 `-v` 日志里命中 0**（逐台计数见 §2 末）。批次 2 那本「软链形含三枚 300 s 腿 ⇒ ≥ 50m」的账在本包不适用：`internal/agent` 零枚 300 s 腿（票 123 那三枚全在 `internal/tools`）。
- 开测前查在飞：`gh run list --limit 4` 见 **1 枚 `in_progress`**（`35863367027`，12:53:23z 起，触发自一枚 docs push）＋ 三枚 completed。⇒ 本机就是 self-hosted runner，存在同机抢 CPU；本批判据全是**红绿名册**而非耗时（两形同容器同码只差 `TMPDIR` 一个变量），按 AC#2a §0 那条同一处理：**影响单枚耗时、不影响红绿**，照跑并在 §9 登记。

## 1. 对派单数字的复核（派单给的每个数字都是断言）

- **账上点名**（读 `124-ac2a-leg-classification.md` §5.2 逐行，行 18-43）：顶层 18 枚（行 18-35）+ 子测试 8 枚（行 36-43）= **26**。⇒ 与派单那句「`internal/agent` 那 26 枚」**一致，登记差 0**。
- **实测复算（pre 锚软链形 `-v` 逐名，台 PRE-L）**：`./internal/agent/` 顶层 `--- FAIL` **18** + 子测试 `--- FAIL` **8** = **26 枚**。逐名与账上名册**枚枚对上**（名册差集见 §2 的 `RED-PRE-L.txt`，26 行）。⇒ **无新增、无缺失**：本批没有像批次 2 那样量到「账上没有、本锚多出来」的红（那发是 `internal/tools` 的 `-timeout` 截断造成的，本包无 300 s 腿、无截断）。
- **红因复算**：PRE-L 日志里机制字串只有两种拼写、且成对出现——`not provably resolved` **29 行**、`refusing to seal` **29 行**、`/varlink` **29 行**（29 > 26 是同一枚用例把同一句打了两次：`forensics`／`guard` 的 `t.Fatalf` 里包着一层 `memory.Open` 的包装串）。⇒ 与账 §5.2 那句「软链形失败串只有两种、无第三因」**逐字复核成立**。
- **递根点数与派单预期的差**：批次 2 的 `next=` 预期「递根点数应远小于枚数」。**实测 19 处调用（18 行）覆盖 26 枚**——比 llm 那本「1 处 fixture 覆盖 17 枚」大，因为 `agent` 的红不是一枚 fixture 被拒、而是**每枚用例各自 `t.TempDir()`**（spill／guard／forensics 各写各的根）。⇒ 登记这条**与派单预期不一致**之处：形状沿用、枚数不变，只是「一处覆盖多枚」在本包只部分成立（`spill_path_invariant_test.go:103` 一处覆盖 6 枚、`spill_name_injectivity_test.go:100` 一处覆盖 4 枚、`spill_test.go:56` 一行两处覆盖 1 枚）。
- pre 锚两形四数（`-count=1 -v`）与改后两形（同一台架、同一脚本、同一容器镜像）：

| 台 | 形状 | `./internal/agent/` |
|---|---|---|
| PRE-L | 软链 | `rc=1 RUN=75 PASS=40 FAIL=18 SKIP=0 SUBPASS=9 SUBFAIL=8 SUBSKIP=0`（包级 `0.764 s`） |
| PRE-P | 普通 | `rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0 SUBSKIP=0`（`1.238 s`） |
| POST-L | 软链 | `rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0 SUBSKIP=0`（`1.237 s`） |
| POST-P | 普通 | `rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0 SUBSKIP=0`（`1.222 s`） |

⇒ 本批软链形红名 **26 → 0**；普通形四数改前改后**逐数相同**。

## 2. 判据①：逐枚转绿且点名（`-v` 才有 PASS 名）

四台对照（同一批用例，`-count=1 -v`）。PRE-L=改前软链、PRE-P=改前普通、POST-L=改后软链、POST-P=改后普通。逐名状态取自各台 `-v` 日志的 `--- (PASS|FAIL|SKIP)` 行（名册 `/d/tmp/wisp124-2b3a-logs/{PRE-L,PRE-P,POST-L,POST-P}.names.txt`、合成表 `TABLE-26.txt` / `TABLE-26.md`）。

| 用例名（子测试写作 `父 / 子`） | PRE-L | PRE-P | POST-L | POST-P |
|---|---|---|---|---|
| `TestCancelledTaskPersistsTerminalRows` | FAIL | PASS | **PASS** | PASS |
| `TestFailedTaskBooksOpenCallRowWithDecision` | FAIL | PASS | **PASS** | PASS |
| `TestMaxTokensFailsAllToolCallsOfThatMessage` | FAIL | PASS | **PASS** | PASS |
| `TestPerToolTimeoutFires` | FAIL | PASS | **PASS** | PASS |
| `TestPerToolTimeoutOfContractHonestToolIsToolClass` | FAIL | PASS | **PASS** | PASS |
| `TestSpillAcrossRestartsKeepsRetrySemantics` | FAIL | PASS | **PASS** | PASS |
| `TestSpillArtifactAndStubShape` | FAIL | PASS | **PASS** | PASS |
| `TestSpillArtifactRespectsRawCap` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames / dotdot` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames / drive_letter` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames / encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames / separator` | FAIL | PASS | **PASS** | PASS |
| `TestSpillCallIDHostileShapesSanitizedToBareNames / unc` | FAIL | PASS | **PASS** | PASS |
| `TestSpillContainmentByDirectoryListing` | FAIL | PASS | **PASS** | PASS |
| `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir` | FAIL | PASS | **PASS** | PASS |
| `TestSpillSameIDRetryOverwrites` | FAIL | PASS | **PASS** | PASS |
| `TestSpillThresholdScalesWithWindow` | FAIL | PASS | **PASS** | PASS |
| `TestSpillThroughLoop` | FAIL | PASS | **PASS** | PASS |
| `TestSpillTokenBoundary` | FAIL | PASS | **PASS** | PASS |
| `TestSpilledBytesSurviveANameThatUsedToCollide` | FAIL | PASS | **PASS** | PASS |
| `TestSpilledBytesSurviveANameThatUsedToCollide / bare_id_then_backslash_id` | FAIL | PASS | **PASS** | PASS |
| `TestSpilledBytesSurviveANameThatUsedToCollide / dotdot_id_then_word_id` | FAIL | PASS | **PASS** | PASS |
| `TestSpilledBytesSurviveANameThatUsedToCollide / slash_id_then_bare_id` | FAIL | PASS | **PASS** | PASS |
| `TestTaskLogAndToolCallRows` | FAIL | PASS | **PASS** | PASS |
| `TestWriteFileExclusiveIsExclusive` | FAIL | PASS | **PASS** | PASS |

⇒ 26 枚逐枚点名，POST-L 列全部是 `PASS` 字样（不是 `SKIP`、不是包级 `ok`）。

**逐台核对（脚本 `comm`/`diff`，不靠目测）**：

- 用例名集合四台**逐名相等**：`PRE-L vs PRE-P`、`PRE-L vs POST-L`、`PRE-L vs POST-P` 三对 `comm -3` 都 **0 行** ⇒ 既没多出一枚用例，也没少一枚（分母恒定 75 名）。
- 四台 `SKIP` 计数（顶层 + 子测试）**全 0** ⇒ 本包不存在「被跳过冒充变绿」的形状。
- 机制字串行计数（同一台架逐台 `grep -c`）：

| 台 | `not provably resolved` | `refusing to seal` | `/varlink` | `test timed out` |
|---|---|---|---|---|
| PRE-L | 29 | 29 | 29 | 0 |
| POST-L | **0** | **0** | **0** | 0 |
| PRE-P | 0 | 0 | 0 | 0 |
| POST-P | 0 | 0 | 0 | 0 |

⇒ 本批软链形红名数 **26 → 0**，且日志里那两句「未解析根」的机制字串同时归零。
