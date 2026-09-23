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

## 3. 判据②：普通形一枚都不许多红 + `git diff` 证判定分支一字未动

**逐数相同**：POST-P 与 PRE-P **八个数逐数相同**——`RUN=75 PASS=58 FAIL=0 SKIP=0` + `SUBPASS=17 SUBFAIL=0 SUBSKIP=0`（连包级耗时都同量级 `1.238 s` / `1.222 s`）。普通形红名数 **0 → 0**，一枚都不许多红。用例名集合改前改后普通形**逐名相等**（`comm -3` 0 行 ⇒ 既没多一枚也没少一枚）。

`git show 62dda11` 的**删除侧全文只有 18 行**（`git show --numstat` ＝ 8 枚文件、删除列合计 18），逐形抄（`-` 号省略、`^I` 是一个制表符）：

- `^Idir := t.TempDir()` — **14 行**（forensics 2 + guard 2 + loop_golden 1 + spill_name_injectivity 4 + spill_test 4 + truncation 1）
- `^Iroot := t.TempDir()` — **3 行**（spill_path_invariant 的 `:103` / `:200` / `:345`）
- `^IdirBig, dirTiny := t.TempDir(), t.TempDir()` — **1 行**（spill_test `:56`，一行两处调用）

**新增侧**（除那枚新文件外只有 18 行）：18 行全部是同名替换（`sealableTempDir124(t)` 形）。逐文件 `+/-` 严格对称（`2/2`、`2/2`、`1/1`、`4/4`、`3/3`、`5/5`、`1/1`），第 8 枚文件是 `26/0` 的新委托 ⇒ **判定分支、断言、阈值、golden 一行未动**：删除行与新增行一一对应，全部落在「递给底线／`NewSpiller` 的根怎么拼」这一件事上。⇒ 票面 AC#3「不许拿放行侧放宽换绿」在本批成立：解析层在调用方（测试）这一侧，`internal/winsec` 与 `internal/risk` 的判定码一个字没碰。

范围核对（脚本，不靠目测）：

- `git show --numstat 62dda11` ＝ 8 枚文件、**全部 `_test.go`**（7 枚改点 + 1 枚新委托）；`git diff --name-only 4ea0db2 62dda11 -- internal/ cmd/` 里非 `_test.go` 命中 **0 枚**（本批零生产码 hunk）。
- 禁改列**零 hunk**：`internal/winsec/**`（那是并行的 3b）、`internal/risk/**`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、任何阈值／golden／`thresholds.go`、`frontend/**`、`docs/PLAN.md`、`docs/specs/**`。
- 批次 1／2 已交的文件**零 hunk**：`internal/memory/**`、`internal/config/**`、`internal/tools/**`、`internal/llm/**`、`internal/perm/**`、`internal/agent/approval/**`（`git diff --name-only 4ea0db2 62dda11` 的 8 枚路径全在 `internal/agent/` 顶层）。
- `cmd/wisp/**` 本批零 hunk（2b-4 地界）。

⚠ **本包内留着没改的同类站点（逐名登记，与批次 2「非红用例的站点不顺手改」同处理）**：

| 站点 | 为什么没改 |
|---|---|
| `internal/agent/spill_test.go:157`（`sp := NewSpiller(t.TempDir(), Budgets{RawOutputCapBytes: 100})`） | 属 `TestSpillRawHardCap`（`:156`），**两形都不红**（它只在内存里算 `CapRaw`，没让底线碰这个根）⇒ 不在账上 26 枚内，不改 |
| `internal/agent/harness_test.go:107`（`ArtifactsDir: t.TempDir()`） | 该枚在 PRE-L 是 PASS（两形都不红）⇒ 不在本批账面内 |
| `internal/agent/spill_acl_windows_test.go:20` | 带 `//go:build windows`，POSIX 容器里**连编译都不参与** ⇒ 本票形状无分母 |
| `internal/agent/tempdir_resolved_124_test.go:25` | 就是委托本体里那一枚 `t.TempDir()`（被解析的那一枚），不是递根点 |

⇒ 改后 `grep -n "t\.TempDir()" internal/agent/*_test.go` 只剩上面 4 行（其中 3 行是本批范围外的站点、1 行是委托本体），其余 18 行全部换成同名委托。

## 4. 判据③：两形各取一枚 + RUN/SKIP 差逐包解释

两形各一枚已完成（§2 的 POST-L / POST-P；另 `-count=2` 两形各一枚见 §7.1）。本批只裁**一枚包** `./internal/agent/`，逐包解释：

- **`RUN` 两形逐台相等、差为 0**：PRE-L `75` = PRE-P `75` = POST-L `75` = POST-P `75`。⇒ 票面 15:3x 那条「软链形自己会缩小分母」的附带发现（AC#1 记的是 `memory 36→67`、`tools 77→79`、`config 99→101`、`winsec 45→52` 四枚包）**在本包不复现**，`internal/agent` 从来不在那四枚里。**机制差别**：memory 那一族是父用例死在 `memory.Open` ⇒ `t.Run` 的 31 枚子测试从未被创建；本包 26 枚的红全部落在**子测试体内各自的那次 `Prepare`/`NewSpiller` 调用**上（父用例先建好子测试、子测试自己拿错），所以四台的子测试名数恒为 `SUBPASS+SUBFAIL = 17`（PRE-L `9+8`、PRE-P `17+0`、POST-L `17+0`、POST-P `17+0`）⇒ 分母一枚没缩。
- **`SKIP` 两形逐台相等、恒为 0**（顶层 + 子测试，四台八数全 0）。⇒ 本包**没有**形状自带的 SKIP：AC#1 点名的那 4 枚「从跑变成 SKIP」全在 `internal/winsec`（3 枚票 125 seam 自拒探针）与 `internal/config`（1 枚），**一枚都不在本包**（那三枚 winsec 探针属 3b 的地界与解释）。
- **「变绿」还是「被跳过」逐枚分清**：① §2 表 POST-L 列 26 枚**逐枚写的是 `PASS`**，取自 `-v` 日志的逐名 `--- PASS` / `    --- PASS` 行（非包级 `ok`）；② 四台 `SKIP` 计数全 0 ⇒ 没有任何一枚可以用 SKIP 蒙绿；③ 用例名集合 PRE-L ↔ POST-L `comm -3` = **0 行** ⇒ 改后软链的可见用例集合与改前软链**同一批**，不是「少跑了几枚所以绿」；④ 红串的机制字串（`not provably resolved` / `refusing to seal` / `/varlink`）**29 → 0**，是同一枚用例走通了 setup，不是没走。
- 普通形这一发同样 `FAIL=0 SUBFAIL=0 SKIP=0` ⇒ 判据②（普通形一枚都不许多红）与判据③（两形各一枚）用的是同一台架、同一脚本，只差 `TMPDIR` 一个变量与 `exit 99` 那枚防「容器复用把 `/varlink` 留下」的断言。

## 5. 判据④：本批「另一种拒」复算 — 逐枚实拿被拒字符串

**取字串仪器**：`/d/tmp/wisp124-2b3a-probe` ＝ 改后快照 + **9 处纯打印插入**（`probe-patch.py` 生成，只加 `t.Logf("ACQ124-2B3A …")`，一处不动断言／阈值；行号锚不中即停手不写盘）。**探针台 PROBE-L 四数与 POST-L 逐数相同**（`rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0 SUBSKIP=0`）⇒ 探针没改变任何结局。以下字符串**逐字**抄自该台软链形 `-v` 日志的 `ACQ124-2B3A …` 行（原文全量在 `/d/tmp/wisp124-2b3a-logs/PROBE-L.__internal_agent_.txt`）。

本批 26 枚里，断言方向是「必须拿到某个拒／某个判决」的共 **8 个入口**，逐枚实拿：

1. `TestFailedTaskBooksOpenCallRowWithDecision`（要的是那行落下闸门判决，不是执行）实拿：`P1 forensics row: outcome="error" decision="reject" error_class="network" ended_at_null=false args_json="{\"text\":\"前半段"` ⇒ 拿到的正是它断言的 `DecisionReject` + `ClassNetwork`。
2. `TestPerToolTimeoutFires`（要「工具被 60ms 协作式中止并告诉模型」）实拿：`P2 per-tool-timeout row: outcome="error" error_class="tool" text="工具 sleep 超时（60ms），已协作式中止"`。
3. `TestPerToolTimeoutOfContractHonestToolIsToolClass`（决定性断言：工具的自有错不许偷走 C22 类别）实拿：`P3 contract-honest-timeout row: outcome="error" error_class="tool" text="工具 sleep 超时（60ms），已协作式中止"` ⇒ 类别仍是 `tool`。
4. `TestMaxTokensFailsAllToolCallsOfThatMessage`（要两行都是 truncated + reject + loop）实拿：`P4 truncation row: tool="echo" outcome="truncated" decision="reject" error_class="loop"`（两枚子行逐字相同，`sort -u` 折成一行）。
5. `TestWriteFileExclusiveIsExclusive`（要第二次独占写**被拒**且前任字节不动）实拿：`P5 exclusive second write refusal: err=open /realpriv/w124tmp/TestWriteFileExclusiveIsExclusive1049974371/001/artifact.txt: file exists errors.Is(fs.ErrExist)=true` ⇒ 拿到的是 `fs.ErrExist` 那一族拒，且串里的根已是解析后的 `/realpriv/…`（这层解析起效的直接旁证）。
6. `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`（要 6 种「调用方点名的路径」形状全被 `DeleteArtifact` 拒）实拿 6 行，逐字样例：`P8 caller-named delete refusal: shape="nested/canary.txt" err=memory: invalid artifact name "nested/canary.txt" (bare file names only)`；其余五枚形状 `../canary.txt`、`..\canary.txt`、`nested\canary.txt`、绝对路径 `/realpriv/w124tmp/…/diary.txt`、UNC `\\localhost\/$ealpriv/…` **各拿到同一族 `memory: invalid artifact name … (bare file names only)`**、无一枚拿到「未解析根」那句。
7. `TestSpillContainmentByDirectoryListing`（要的是「每种恶意 callID 都 spill 成功、目录树只在 artifacts 下长」）实拿：`P7 containment collected prepare errors: total=0` ⇒ 一枚 Prepare 错误都没拿到（同一枚用例改前软链形名下的 `refusing to seal` 行数是 **7**，台 PRE-L）。
8. `TestSpillCallIDHostileShapesSanitizedToBareNames` 的 4 枚子测试（名字带 Hostile，但断言方向是**编码而不是拒绝**，源码 `:99` 那句明写）实拿：`P6 hostile callID shape="separator" prepare_err=<nil> spilled=true name="tool-output-p%2Fq.txt" path="/realpriv/…/data/separator/artifacts/tool-output-p%2Fq.txt"`，另有 `dotdot` → `tool-output-%2E%2E%2F%2E%2E%2Fescape.txt`、`drive_letter` → `tool-output-%43%3A%5C%57indows%5C%53ystem32%5Cdrop.txt`、`unc` → `tool-output-%5C%5Cfileserver%5Cshare%5Cpayload.txt` ⇒ 四枚都拿到 `prepare_err=<nil>` + `spilled=true` + 净化后的裸名，正是它们各自要的那一形。

**其余 15 枚**（上面 8 个入口共点到 **11 枚**具名用例：7 枚顶层 + `TestSpillCallIDHostileShapesSanitizedToBareNames` 的 4 枚子测试；26 减 11 余 15）断言方向全是「必须成」，不涉任何拒：`TestCancelledTaskPersistsTerminalRows`、`TestTaskLogAndToolCallRows`、`TestSpillThresholdScalesWithWindow`、`TestSpillTokenBoundary`、`TestSpillArtifactAndStubShape`、`TestSpillThroughLoop`、`TestSpillArtifactRespectsRawCap`、`TestSpilledBytesSurviveANameThatUsedToCollide` 顶层与其 3 枚子测试、`TestSpillSameIDRetryOverwrites`、`TestSpillAcrossRestartsKeepsRetrySemantics`、`TestSpillCallIDHostileShapesSanitizedToBareNames` 顶层与 `encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence` 子测试。

**结论**：本批 26 枚没有一枚的断言被「未解析根」那句替代。探针台与改后台的软链形日志里 `not provably resolved` / `refusing to seal` / `/varlink` 命中数 = **0**（§2 表）。账 §5.2 那句「可转 ⇒ 接解析后仍成立」在本批 26 枚上逐枚复核成立；票面 16:33 那节列的 6 枚「另一种拒」中落在 `internal/agent` 的是 **0 枚**（那 6 枚分布在 winsec 3、memory 2、tools 1、llm 1 上），所以本批按上面 8 个入口自己找、自己实拿，没有靠清点账的结论推断。

## 6. 判据⑤：变异自证（票面 AC#4）— 三态原文

变异 = **把接上去那一步退回旧实现**：唯一一枚委托 `sealableTempDir124` 的函数体不再解析，把 `t.TempDir()` 原样交回去（`MUTATION-124-2B3A`；`proc` 引用留一行 `_ =` 保编译）。18 处调用点不动 ⇒ 整包行为逐字回到改前。台子 `/d/tmp/wisp124-2b3a-mut`（由 post 快照复制），跑法同一支 `/d/tmp/wisp124-2b3a-h.sh`。

**变异先证落地，再读数**（容器 `wisp124-2b3a-mutproof2`，原文 `/d/tmp/wisp124-2b3a-logs/MUT-PROOF.txt`）：

```
== mutation landed (grep -n) ==
25:	_ = proc.SealableRoot // MUTATION-124-2B3A: the resolution layer this batch added is removed here
26:	return t.TempDir() // MUTATION-124-2B3A: hand the OS's unresolved root straight back
return-proc-SealableRoot-remaining: 0
== go build ./... ==
go build rc=0
== go vet ./internal/agent ==
go vet rc=0
```

⇒ 被改后那一行 `grep -n` 命中两行、`return proc.SealableRoot(t.TempDir())` 命中 **0**、`go build` rc=0、`go vet` rc=0，四条先于任何红名读数。⚠ 一处如实登记：插入后 `gofmt -l` 曾点到这枚被改文件（注释对齐），本方对**变异台**跑了 `gofmt -w` 再进容器（`_ = proc.SealableRoot` 与 `return t.TempDir()` 变成对齐的两枚尾注释，语义一字未变）；入库的那棵树从未被 gofmt 碰过（§7.2 的 `gofmt -l` 全 0 行是在 post 快照上量的）。

三态读数（同快照血统、同容器镜像、同形状断言）：

| 台 | 形状 | 四数 | 红名 |
|---|---|---|---|
| **POST（未变异）** | 软链 | `rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0` | 本批 **0 枚** |
| **MUT-L（拆掉解析层）** | 软链 | `rc=1 RUN=75 PASS=40 FAIL=18 SKIP=0 SUBPASS=9 SUBFAIL=8` | **26 枚全部回归**，红名与改前基线 PRE-L **逐名 IDENTICAL**（`diff RED-PRE-L.txt RED-MUT-L.txt` rc=0，且连 `(0.00s)` 时延一起比的整行 `--- FAIL` 名册亦 IDENTICAL）；机制字串回到 `not provably resolved`=29、`refusing to seal`=29、`/varlink`=29（与 PRE-L **逐数相同**），样例逐字：`forensics_test.go: memory.Open: memory: create data dir: winsec: refusing to seal /varlink/w124tmp/TestFailedTaskBooksOpenCallRowWithDecision2418577891/001: the installed risk.c26Pipeline answered "/varlink/…", a spelling the floor itself refuses: winsec: path is not provably resolved, refusing to seal: … reaches it through the link at /varlink, which is not the tree this call names` |
| **MUT-P（同一发变异）** | 普通 | `rc=0 RUN=75 PASS=58 FAIL=0 SKIP=0 SUBPASS=17 SUBFAIL=0` | 无（`FAIL` 行 = 0；八数与 PRE-P／POST-P **逐数相同** ⇒ 这层解析只在软链形起作用，普通形是恒等操作） |

⇒ 红名**点到用例自己**（26 枚逐枚点名，不是包级 `FAIL`），且变异先证落地（`grep -n` 出被改后那一行 + `go build` rc=0）后才读数。**这一发做了。**
