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

## 7. 门禁读数（票面 AC#5 的形状，本批只裁本批的包 `./internal/agent/`）

⚠ **测的是哪棵树 / 共树漂移（读数前先核）**：本批全部门禁读数都在 `/d/tmp/wisp124-2b3a-post`（`git archive 62dda11`，1050 枚文件）里量，静态对照那发在 `/d/tmp/wisp124-2b3a-pre`（`git archive 4ea0db2`，1049 枚）。量完后共享树又前进：兄弟代理 `worker-ticket124-ac2b-3b` 的 `8ced405`（`internal/winsec/**` 6 枚文件）落进 `dev`。已按名核清：**`8ced405` 不是 `62dda11` 的祖先**（`git merge-base --is-ancestor 8ced405 62dda11` rc≠0），且 post 快照的 `internal/winsec/` 里**没有** `tempdir_resolved_124_other_test.go`（`ls` 已核）⇒ 本批九发读数与两发 d22scan 台账都量的是**只有本批改件**的那棵树，3b 的账不进本批、本批也不替它裁。（§0 那句「非 `internal/agent/` 命中 0 枚」是 HEAD 还在 `857a5fe` 时核的，当时成立；此后唯一进入 `internal/` 的漂移就是 3b 这一枚，出处同上。）

### 7.1 `-count=2 -v` 四数（本批的包 × 两形；只能从 `-v` 量，`-count=2` 才不缓存）

| 台 | 形状 | `./internal/agent/` |
|---|---|---|
| `GATE2-LINK` | 软链 | `rc=0 RUN=150 PASS=116 FAIL=0 SKIP=0 SUBPASS=34 SUBFAIL=0 SUBSKIP=0`（包级 `2.426 s`） |
| `GATE2-PLAIN` | 普通 | `rc=0 RUN=150 PASS=116 FAIL=0 SKIP=0 SUBPASS=34 SUBFAIL=0 SUBSKIP=0`（`2.344 s`） |

⇒ **每个数都是 §2 那一发 `-count=1` 的正好 2 倍**（`75/58/0/0`+`17/0` → `150/116/0/0`+`34/0`）⇒ 无缓存读数、无 flake、与 §2 无分歧；两形 `FAIL=0 SUBFAIL=0`、`SKIP=0`、两形八数逐数相同。`-timeout 30m` 下 `test timed out` 命中 **0**（本包无票 123 那族 300 s 腿 ⇒ 批次 2 那条「≥ 50m」的硬要求在本包不适用，见 §9）。每形一枚**新**容器（`wisp124-2b3a-gatelink` / `-gateplain`），普通形那发的 `exit 99` 断言（`/varlink` 必须根本不存在）成立。

### 7.2 `gofmt -l`（整包 + 全树参照；宿主与容器各一组）

| 调用（均在锚点纯净快照里） | 结果 |
|---|---|
| 宿主 `gofmt -l internal/agent` | **0 行**，`rc=0`（`/d/tmp/wisp124-2b3a-logs/STATIC-HOST.txt`） |
| 宿主 `gofmt -l internal cmd`（全树参照） | **0 行**，`rc=0` |
| 容器 `gofmt -l internal/agent`、`gofmt -l internal cmd` | **各 0 行**，`rc=0`（`VET-LINUX-CONTAINER.txt` 末两段） |

⇒ 改动的 8 枚文件与整包、整树都干净；票面「注释零 emoji／`_test.go` 也算」那一族由 §7.5 那把尺覆盖。

### 7.3 `"$(go env GOPATH)/bin/gofumpt.exe" -l`（票面 AC#5：真跑；写「未跑」必须引错误原文）

版本行逐字 **`v0.7.0 (go1.27.1)`** ⇒ 工具存在，本格**不是**「未跑」，没有错误原文要引。宿主对 `internal/agent` 与 `internal cmd` 各一发 ⇒ **均 0 行**、`rc=0`（`STATIC-HOST.txt`）。**并与 CI 逐字同形再跑一发**：在 post 快照根目录 `gofumpt -l . tools/d22scan tools/mockllm` ⇒ **0 行**、`rc=0`（`GOFUMPT-CI-SHAPE.txt`；同文件里 `gofmt -l .` 亦 0 行）⇒ 门禁自检路径与 CI 那一步同源，不是只扫了本批那几个目录。⚠ 如实登记两件事：① **CI 里确实跑 gofumpt，但版本没钉**——`.github/workflows/ci.yml:111-114` 那一步逐字是 `go install mvdan.cc/gofumpt@latest` 然后 `gofumpt -l . tools/d22scan tools/mockllm`，`@latest` 意味着门禁用的是"当天最新版"而不是某一枚被冻结的版本 ⇒ 这里写的 **`v0.7.0`** 是本方用的那一枚（与批次 1／2 同一枚），复现者拿到不同版本时以这一步的输出为准；本批改动全在 `internal/`（在那条命令的 `.` 作用域内）、且宿主 `gofumpt -l internal cmd` 为 0 行 ⇒ 对 `v0.7.0` 成立；② `golang:1.27` 容器里**没装** gofumpt 且带 `GOPROXY=off`（不联网装第三方工具）⇒ 本批 gofumpt 只有宿主读数，容器那一侧由 §7.2 的 gofmt + §7.4 的双 GOOS vet 覆盖。

### 7.4 `go vet` 双 GOOS（逐错误行归因；宿主交叉那一发既不算破口也不算清白）

| 调用形状 | rc | 输出 |
|---|---|---|
| `GOOS=windows go vet ./...`（宿主原生、全树、post 快照） | **0** | 0 行（`VET-HOST.txt`） |
| `GOOS=windows go vet ./internal/agent/`（宿主原生） | **0** | 0 行 |
| `go vet ./...`（**容器 linux 原生**、`CGO_ENABLED=1`、`go1.27.1 linux/amd64`、同一快照） | **0** | 0 行（`VET-LINUX-CONTAINER.txt`）⇒ **这一发才是 linux 的真类型读数** |
| `go vet ./internal/agent/`（容器 linux 原生） | **0** | 0 行 |
| `GOOS=linux CGO_ENABLED=0 go vet ./...`（**宿主交叉**、全树） | **1** | 3 行，逐字见下 |
| `GOOS=linux go vet ./internal/agent/`（宿主交叉） | **0** | 0 行 ⇒ 本批那枚包在 linux 目标下过 vet |
| `GOOS=linux CGO_ENABLED=0 go vet $(go list ./... \| grep -vE 'cmd/(wisp\|balldebug)$')`（剔两枚既有形状，余 31 枚；`go list ./...` 本树列得 **33** 枚） | **0** | 0 行（`VET-HOST-EXCL.txt`）⇒ 交叉 `rc=1` **全部**由那两枚贡献 |
| `GOOS=linux go vet ./cmd/balldebug/`（宿主交叉、单点名） | **1** | `package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in D:\tmp\wisp124-2b3a-post\cmd\balldebug` |
| 同一支宿主交叉全树命令，改在 **pre 快照 `4ea0db2`** 上跑 | **1** | **逐字节相同**的 3 行（`VET-HOST.txt` 末段）⇒ 改前就在，与本批无关 |

宿主交叉全树那一发的全部输出（逐字）：

```
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
```

⇒ 死在**包加载**阶段（外部模块 `sherpa-onnx-go-linux@v1.13.8` 的 build constraints），结构上到不了类型检查 ⇒ 按票面纪律**既不算破口也不算清白**；清白由上面那发容器原生 `rc=0` 提供。错误文本点到的是 `cmd/wisp`（2b-4 地界）与一枚外部模块目录，**零行指向 `internal/agent`**。与批次 2 §7.4 那本账一致（那发亦报 1 处；报几处取决于加载顺序，把 `cmd/balldebug` 单独点名它照样红，本方已逐字量）。

### 7.5 一次全仓仪器：`sh scripts/d22scan.sh`（作用域整个 `internal/`；按包门禁结构性看不见它）

调用形状（唯一受支持的形状，本方**没撞到** `main module does not contain package`）：在纯净快照里 `cd /d/tmp/wisp124-2b3a-{pre,post} && sh scripts/d22scan.sh`（宿主 Git Bash；脚本从自身位置推导仓根，两步 `set -eu` 依次硬跑：① `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` 播种违规的正向对照 ② `go run . -root <推导出的根>` 真扫）。

| 台 | rc | 正向对照 | 真扫规模 |
|---|---|---|---|
| `D22-PRE`（`4ea0db2`，1049 文件） | **0** | `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` | `examined 225 production Go files under internal/ and cmd/` |
| `D22-POST`（`62dda11`，1050 文件） | **0** | 同一行逐字相同（⇒ 这把尺**能红**，"clean" 不是"眼睛瞎"） | 同一规模 |

台账八 scope 逐数对照（原文 `/d/tmp/wisp124-2b3a-logs/D22-{PRE,POST}.txt`）：

| scope | pre（`4ea0db2`） | post（`62dda11`） | 变化与出处 |
|---|---|---|---|
| bans #1-5 `internal/` | 203 | 203 | 0 |
| bans #1-5 `cmd/` | 22 | 22 | 0 |
| ban #6 `frontend/` | 40 | 40 | 0 |
| ban #7 `internal/tools/` | 18 | 18 | 0 |
| ban #8 `design/` | 16 | 16 | 0 |
| ban #8 `frontend/` | 40 | 40 | 0 |
| ban #8 `internal/` | 397 | **398** | **+1 ＝ 本批那枚单行委托 `internal/agent/tempdir_resolved_124_test.go` 进树**（post 快照文件数亦 +1，见 §0 表） |
| ban #8 `cmd/` | 38 | 38 | 0 |

⇒ **八数一枚都不降**；唯一变化是"多扫了一枚文件"。pre 那发的 397 与批次 2 §7.5 在 `e113b1a` 记的 397 **逐数相同** ⇒ 跨批台账连续、无回退。

**AC#5 小结（本批只裁本批的包）**：`-count=2 -v` 两形四数齐且各为 `-count=1` 那发的 2 倍（两形零红、零 SKIP）；`gofmt -l` 整包 + 整树 × 宿主／容器全 0 行；`gofumpt v0.7.0 (go1.27.1)` 真跑、0 行（版本没钉在 CI ⇒ 写明用的是哪一枚，容器无此工具）；`go vet` 双 GOOS——windows 全树 `rc=0`、linux **容器原生**全树 `rc=0`（真类型读数），宿主交叉 `rc=1` 的 3 行逐行归因到两枚既有形状（pre 快照逐字节同、剔掉即 `rc=0`、本包单点名 `rc=0`）；`sh scripts/d22scan.sh` pre／post 各 `rc=0`、正向对照 21/0/0、台账八 scope 不降。**本格不为 `internal/agent` 翻 AC#2 那一格，也不翻 AC#5**（票面 16:33 已定终判据留到 2b-4 交回后一次性复算；AC#5 要等所有受影响包——含 3b 的 `internal/winsec` 与 2b-4 的 `cmd/wisp`——都进门禁账）。

## 8. 改动面与本批用的 helper

- 用的既有 helper = **`proc.SealableRoot`**（`internal/proc/envfork.go:160`），即票 119 生产路 `TestDataDir` / `DefaultLayout` / `cmd/wisp` 的 `resolveDataDir` 走的那一枚；批次 1／2 同法。**没另发明第二套**：本批只新增**一枚**委托。
- 为什么本批只需一枚委托（批次 2 的 `next=` 那句断言，本方复核成立）：涉事 7 枚文件 `grep '^package '` 逐枚核过，全部是 `package agent`（含 `spill_acl_windows_test.go`）⇒ 不存在 tools 那种内／外双测试包作用域问题，因此**不需要** tools 的第二枚 `124x`。
- 新增面：`internal/agent/tempdir_resolved_124_test.go`（26 行：包注释 + 一段"为什么"的说明 + 一枚 3 行函数）。**零新解析路数**：`R-125-2` 那本副本账（票 125 现计三枚：`proc.SealableRoot`、`internal/tools/paths.go`、`winsec/resolve.go` 的 `resolveProbeRoot`）**没有增行**——委托体就是一行 `return proc.SealableRoot(t.TempDir())`。
- 递根点 **18 行 / 19 处调用**，覆盖 26 枚。最集中的三处：`spill_path_invariant_test.go:103`（1 处覆盖 6 枚＝顶层 + 5 子测试）、`spill_name_injectivity_test.go:100`（1 处覆盖 4 枚）、`spill_test.go:56`（1 行 2 处覆盖 1 枚）。⇒ 与批次 2 预期的"远小于枚数"部分相符，差异与原因登记在 §1。
- **`internal/winsec/**` 本批零 hunk**（那是并行的 3b），并且本批没有让任何 winsec 码去调 `proc.SealableRoot`：`grep -n "SealableRoot(" internal/winsec/*.go` 的非注释命中数 **0** ⇒ `internal/proc/envfork.go:148` 那句边界话（"…and nothing in internal/winsec calls it"）**不因本批腐坏**。本批把调用方从 8 枚增到 **9 枚**（`memory`、`config`、`tools` ×2、`llm`、`perm`、`agent/approval`、`agent`），全在 winsec 之外。
- 未动（本包内）：`spill_test.go:157`、`harness_test.go:107`、`spill_acl_windows_test.go:20`（逐名理由见 §3 末表）；`compress_test.go`、`control_test.go`、`prompt_test.go`、`testtools_test.go`、`harness_test.go`、`doc.go` 与全部非测试码。
- 未动（本包外）：批次 1／2 已交面、`internal/winsec/**`、`internal/risk/**`、`cmd/wisp/**`、`rules_gateway.go`、`allowlist.txt`、任何阈值／golden／`thresholds.go`、`frontend/**`、`docs/PLAN.md`、`docs/specs/**`、`tools/d22scan/**`。

## 9. 未验证项（照实列，不含"我相信"）

1. **macOS 那一半仍然没有实测。** 这一族在票面《事实》第 3 条明写"代表的是 macOS 的真实形状"，而本批全部读数都是 **Linux 容器 `golang:1.27` + `ln -s /realpriv /varlink` 的代理形状** ⇒ "这 26 枚在 macOS 真机上会同样转绿"是**外推**。与 AC#1／AC#2a／批次 1／批次 2 同一本未付账。
2. **CI 两腿对本批零信号，且本批没有可引的 CI run。** 只 commit 未 push ⇒ `dev` 上**不存在**以 `62dda11` 为 `head_sha` 的 run。开工前那发 `gh run list --limit 4` 读到 1 枚 `in_progress`（`35863367027`，docs push）＋ 3 枚 completed（其中两枚 `failure`），本机就是 self-hosted runner ⇒ 存在同机抢 CPU；本批判据是红绿名册而非耗时，按 AC#2a §0 同处理（影响单枚耗时、不影响红绿），但**没有**跑完的 CI 读数可引。CI 看不见这一族的结构性原因照旧：ubuntu 腿 `/tmp` 是真目录、windows 腿 `%TEMP%` 不带链接。
3. **宿主（Windows）侧没有 `-count=2` 四数**：AC#5 的四数在容器两形取；Windows 的 `%TEMP%` 不属本票形状，未量（票面亦未要求）。同理宿主那发 `GOOS=linux go vet` 是**交叉**，结构上到不了类型检查（§7.4 已逐行归因），真读数在容器那发。
4. **gofumpt 只有宿主读数**（容器无该工具、`GOPROXY=off` 不联网装），且版本没钉在 CI ⇒ 写的是本方用的 `v0.7.0 (go1.27.1)`。
5. **判据④的"实拿"覆盖 11 枚具名用例，另 15 枚是按断言原文读码定为「必须成」**（票面 16:33 列的 6 枚"另一种拒"里落在 `internal/agent` 的是 0 枚，所以没有账上点名的枚可拿）。这 15 枚的"没拿到未解析根那句"由两枚**包级**读数支撑（`not provably resolved`／`refusing to seal`／`/varlink` 在 POST-L 与 PROBE-L 命中均为 0、四台 `FAIL=0`），不是逐枚探针 ⇒ 若要把这 15 枚也做成逐枚实拿，得再开一枚探针台，本批没做。
6. **票 123 那三枚 300 s 审批超时腿本批零枚命中**（`TestL1WriteGoesThroughTheRealBlockWindow` / `TestLateVetoRendersTheApprovalLayersAppliedStepsReport` / `TestFSReadOnlyNeverOpensACard` 全在 `internal/tools`）⇒ 本批**没有**豁免项、也没顺手修它们；批次 2 那条「`-count=2` 包级 `-timeout` ≥ 50m」在本包不必要（实测两形各 `2.4 s`，最长腿 < 0.02 s）。
7. **全树终判据未复算**（票面 16:33 那条"软链形红名数＝0"）——票面已定"留到 2b-4 交回后一次性复算" ⇒ **AC#2 本格不翻**、AC#2b 各批格亦不由本方翻；**AC#5 本格亦不翻**（3b 的 `internal/winsec` 与 2b-4 的 `cmd/wisp` 还没进门禁账）。
8. **本批只交批次 3 的前一半**：`internal/winsec` 那 17 枚在 `worker-ticket124-ac2b-3b` 手上（其 `8ced405` 已在共享树里，见 §7 抬头）；本程没读它的证据、没评它、也不替它记账。

**`next=`（写给 2b-4，并给终判据复算方；交编排者）**

- **2b-4＝`cmd/wisp` 5 枚。** 名册沿用批次 2 §9 那段（4 枚顶层 `providers_test.go:133`/`:179`/`:208`/`:228` 全经同一枚 `newProvidersFixture(t)` ⇒ 一处递根覆盖 4 枚；另 1 枚子测试 `secret_test.go:679`）。补三条本批量到的、能给 2b-4 省事的账：
  1. **形状与工具链可直接沿用本批这一套**：`/d/tmp/wisp124-2b3a-h.sh`（`bash /h3a.sh <link|plain|none> <tag> [pkgs]`，`COUNT`／`TIMEOUT` 走环境变量，逐台自动产 `.names.txt` 名册）、`/d/tmp/wisp124-2b3a-convert.py`（行号锚定 + `ANCHOR MISS` 即停手不写盘）、`/d/tmp/wisp124-2b3a-probe-patch.py`（纯打印探针）、`/d/tmp/wisp124-2b3a-mutate.py`（判据⑤）。`/src` 以 `:ro` 挂即可，命名卷 `ac119-gomodcache`/`ac119-gocache` 离线可编。
  2. ⚠ **同包另有 28 枚两形都红**（DPAPI 族与命令面腿，`R-119-7` 那本账）⇒ 一枚都不许顺手修绿，也不许因包级 `FAIL≠0` 把自己那 5 枚判成没转绿：**判据是逐枚点名，不是包级 rc**。本批判据③那条"名集合 `comm -3` 两两 0 行"的做法在 `cmd/wisp` 更要紧（那 28 枚会一直在名册里，别把它们读成"没转绿"）。
  3. ⚠ 地界按**文件级**核：`cmd/wisp/leg_dispatch_gate_133_test.go` 系票 133。票面明写 2b-4 排在 `acceptor-ticket131-r3` 交回之后——动码前先按名核清哪枚文件有谁在飞。
- **给全树终判据复算方（2b-4 交回后一次性；六步形状见批次 2 §9，此处只加三条本批新增的可核断言）**：
  1. `./internal/agent/` 两形名册**应逐名相等**（本批实测四台都 `RUN=75`、`SKIP=0`、名集合 `comm -3` 全 0 行）⇒ 复算那发若见 agent 两形 `RUN` 不等或冒出 SKIP，就是**新破口**，不是形状副作用（本包没有形状自带的 SKIP 用例）。
  2. 本批 26 枚的红名与账 §5.2 名册**逐名相同**（`/d/tmp/wisp124-2b3a-logs/RED-PRE-L.txt`，26 行）⇒ 差集期望为空的那 131 枚里，agent 贡献 26 枚，可直接拿这枚文件当一侧。
  3. 豁免名单照批次 2 钉死的三条：票 123 那族 **三枚** 300 s 腿（含 §1 新登记那枚，三枚一起豁免否则得"多一枚红"的假破口）＋ 两形都红 29 枚（panel 1 + cmd/wisp 28）＋ 形状自带 SKIP 4 枚（winsec 3 + config 1）。⚠ **本批不给 agent 加任何豁免**：本包两形 `SKIP` 恒 0。
- 建议：**AC#2 本格与 AC#5 都继续不翻**，等 2b-4 交回后按"软链形红名数＝0（豁免逐名钉死）"一次性复算。

## 10. 临时件清单（只建不删；本程零 `rm`）

快照四台：`/d/tmp/wisp124-2b3a-pre/`（`4ea0db2`）、`/d/tmp/wisp124-2b3a-post/`（`62dda11`）、`/d/tmp/wisp124-2b3a-probe/`（改后 + 9 处纯打印探针，**未入库**）、`/d/tmp/wisp124-2b3a-mut/`（改后 + 判据⑤变异，**未入库**）。
仪器四支：`/d/tmp/wisp124-2b3a-h.sh`、`-convert.py`、`-probe-patch.py`、`-mutate.py`。
日志与名册：`/d/tmp/wisp124-2b3a-logs/` —— 9 份 `-count=1/-count=2 -v` 原始日志（`{PRE-L,PRE-P,POST-L,POST-P,PROBE-L,MUT-L,MUT-P}.__internal_agent_.txt` 与 `GATE2-{LINK,PLAIN}.*`）、同名 `.summary.txt`／`.roster.txt`／`.names.txt`／`.byname.txt`、红名册 `RED-PRE-L.txt` 与 `RED-MUT-L.txt`、逐枚四态表 `TABLE-26.txt`／`TABLE-26.md`、静态读数 `STATIC-HOST.txt`／`VET-HOST.txt`／`VET-HOST-EXCL.txt`／`VET-LINUX-CONTAINER.txt`／`MUT-PROOF.txt`、全仓仪器 `D22-PRE.txt`／`D22-POST.txt`、开工在飞查读 `GH-RUN-LIST.txt`。
被本文件引用的目录谁都别删（它们是 §2-§7 每个数的一手出处）。
