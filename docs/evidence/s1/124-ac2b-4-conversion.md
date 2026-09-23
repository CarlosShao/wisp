# 票 124 AC#2b 批次 4（`worker-ticket124-ac2b-4`）— `cmd/wisp` 那 5 枚接上票 119 那条纪律

**日期**：2026-09-23（测量段 21:38–21:5x +8；本条落笔前的 `date -u` = 13:54z ⇒ +8 = 21:54）· **性质**：转换方交件（改 `.go`，只改 `_test.go`）
**本批范围**：清点账 `docs/evidence/s1/124-ac2a-leg-classification.md` §5.9（行「128-132」，共 5 枚）＝`cmd/wisp` 一整包，账上全标「可转」；这是 AC#2b 的**最后一批**。
**不在本批**：全树终判据「软链形红名数＝0」的复算（编排者明写另安排，因票 133 的 AC#2 修复也在改 `cmd/wisp`）——本件 §9 的 `next=` 只把豁免名单逐名钉死，不复算。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2b 的放行与分批」（09-23 16:33 编排者）那张表 + 五条共用结案判据；批次 3b 的 `next=` 给这一批的原话在 §1 逐条复核。

## 0. 锚点 / 快照 / 跑法（可复核）

- 开工首读 `git rev-parse --short HEAD` = **`54123e0`**（下文记 **pre 锚**；全部改前读数在它上面量）。分支 `dev`。
- 改件 commit：**`5265c3a`**（`test(124,AC#2b-4)`，3 枚文件、+40/-2；删除侧 2 行**全是 `t.TempDir()` 那一个表达式**，见 §3）。
- 开工时 `git status --porcelain` 的五项脏/未跟踪（`M .scratch/…/133-*.md`、`M .scratch/…/135-*.md`、`?? .scratch/…/137-*.md`、`?? docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`、`?? internal/observe/sampler_zerosample_136_test.go`）**没有一项落在 `cmd/wisp/`** ⇒ 首读的 HEAD 与本批被验包同一版本。
- 共树漂移：测量途中 HEAD 经 `4fc65dd`/`a122240` 到本方改件落库的 `5265c3a`，落库后又到 `6da6710`。已按名核清：`git diff --name-only 54123e0 5265c3a -- cmd/wisp/` 只有本批这 3 枚路径 ⇒ pre 四台与本件所有改前读数都是「只缺本批这 3 枚改动」的那棵树，改后读数都是只多这 3 枚的那棵树。`54123e0` **是** `5265c3a` 的祖先（`git merge-base --is-ancestor` rc=0）。
- ⚠ **兄弟在飞同包**：`worker-ticket133-ac2-fix` 独占 `cmd/wisp/leg_dispatch_gate_133_test.go`。本程对该文件**零 hunk**（`git show --numstat 5265c3a` 只列出本批 3 枚路径），且全程未 `git add` 它（每次提交前 `git diff --cached --name-only` 逐枚点过）。落库后该文件在工作树里出现未提交改动（`M cmd/wisp/leg_dispatch_gate_133_test.go`），**不影响本批任何一发**——四台读数全部跑在 `git archive` 纯净快照里，不读工作树。
- 全部测量都在 `git archive <sha> | tar -x` 的纯净快照里做；脏工作树一个字都不读；**仓内未建 worktree、未 checkout**。

| 快照目录 | 内容 | 文件数 |
|---|---|---|
| `/d/tmp/wisp124-2b4-pre` | `git archive 54123e0`（改前基线） | 1054 |
| `/d/tmp/wisp124-2b4-post` | `git archive 5265c3a`（改后） | 1058（本批 +1＝新委托；另 3 枚是共树漂移带进来的他票文件） |
| `/d/tmp/wisp124-2b4-probe` | 改后 + 判据④取字串探针（6 处纯打印插入 + 1 枚只打印 helper，**未入库**） | 1058 |
| `/d/tmp/wisp124-2b4-mut` | 改后 + 判据⑤变异（把委托退回旧实现，**未入库**） | 1058 |

- 容器 `golang:1.27`（`go1.27.1 linux/amd64`、`WSL2 6.6.114.1-microsoft-standard`、`CGO_ENABLED=1`、`GOPROXY=off`），复用票 119 的命名卷 `ac119-gomodcache`/`ac119-gocache` ⇒ 离线可编。
- 挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`，`/src` 以 `:ro` 挂；**每枚样本进容器第一件事**打 `ls -l /src/go.mod` = `-rwxrwxrwx 1 root root 883`、`md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum /src/internal/winsec/resolve.go` = `b6876a5efe759f6e17434d1b50a129c3` ⇒ 与 AC#1／AC#2a／批次 1／批次 2／3a／3b **逐字同字**，非空挂自证（Git Bash 下 `docker run -v "C:\…"` 会静默挂空且 rc=0＝假绿）。另加一发 `ls -d /src/cmd/wisp` 打真实包目录，防「挂对文件、跑错包」。
- 形状硬断言沿用批次 1／2／3a：软链形 `exit 96/97/98`（链接必须真是链接、`readlink -f /varlink/w124tmp` 必须等于 `/realpriv/w124tmp`）；普通形 `exit 99`（`/varlink` 必须**根本不存在**、`/plainroot` 不许是链接）。批次 3b 自撞那一发新增的 `exit 90`（容器内 `cd /src` 失败即停，防「四数全 0 被当成没红」）**本批写进脚本**。**本批 8 发无一命中**。软链形逐字形状行：`SHAPE=link TMPDIR=/varlink/w124tmp resolved=/realpriv/w124tmp` + `lrwxrwxrwx 1 root root 9 ... /varlink -> /realpriv`；普通形：`SHAPE=plain TMPDIR=/plainroot/w124tmp resolved=/plainroot/w124tmp`。
- **每形一枚新容器**（批次 2 §9 那本账）。容器名逐枚：`wisp124-2b4-{prelink,preplain,postlink,postplain,probelink,mutproof,mutlink,mutplain}` ＋ 门禁两发（全部 `--rm`）。
- 跑法脚本（新建，未改前三批任何一支）：`/d/tmp/wisp124-2b4-h.sh`（`bash /h4.sh <link|plain|none> <tag> [pkgs]`，`COUNT=2` 切门禁那一发、`TIMEOUT` 默认 50m）、`/d/tmp/wisp124-2b4-docker.sh`（宿主侧一发一容器）、`/d/tmp/wisp124-2b4-roster.sh`（把 `-v` 日志折成逐名名册）、`/d/tmp/wisp124-2b4-convert.py`（2 处行号锚定替换，锚不中即 `ANCHOR MISS` 且不写盘）、`/d/tmp/wisp124-2b4-probe-patch.py`、`/d/tmp/wisp124-2b4-mutate.py`、`/d/tmp/wisp124-2b4-mutproof.sh`。日志与名册全在 `/d/tmp/wisp124-2b4-logs/`。
- 包级 `-timeout 50m`：沿用批次 2 那本「软链形含 300 s 腿 ⇒ ≥ 50m」的账。⚠ **本包实测无需要**：`cmd/wisp` 的票 123 三枚 300 s 腿一枚都没有（那三枚全在 `internal/tools`），`-count=1` 两形包级耗时 `10.1 s`／`10.4 s`，`test timed out` 在 8 发 `-v` 日志里命中 **0**。留 50m 是防那 28 枚两形都红的命令面腿里有隐藏的长等待把它截成"没跑完"，不改变任何读数。
- ⚠ **四数之外还比名册差集**（一条用例 panic 会吞掉同包其余读数）：本批 8 发 `-v` 日志 `grep -ci panic` **全 0**，逐名名册四台恒 71 行，`comm -3` 两两 0 行（见 §2）。
- 开测前查在飞（`/d/tmp/wisp124-2b4-logs/GH-RUN-LIST.txt`）：**1 枚 `in_progress`**（`35870138533`，13:52:35z 起，docs push）＋ 3 枚 completed（全 `failure`）。⇒ 本机就是 self-hosted runner，存在同机抢 CPU；本批判据全是**红绿名册**而非耗时（两形同容器同码只差 `TMPDIR` 一个变量），按 AC#2a §0 同处理：**影响单枚耗时、不影响红绿**，照跑并在 §9 登记。
- 时间戳：先 `date -u`（13:37z 开工、13:54z 本段落笔）再 +8。凭据值不进本报告（本包 `secret_test.go` 的文件头即明写所有 key 值是假占位符 D36#5/R7，本报告未抄任何形似密钥的串）。

## 1. 名单与数字复核（派单给的每个数字都是断言，包括那个 5）

**① 分母自核**：从 AC#2a 清点表取 `cmd/wisp/**` 那一档（`124-ac2a-leg-classification.md` §2 表行 `./cmd/wisp/` ＝ `5/5/0/0`；§5.9 逐枚行号 128-132）：

| # | 用例名 | 递根点 | 账上判定 |
|---|---|---|---|
| 1 | `TestProvidersProbeRecordsMeasuredThinkingFalse` | `providers_test.go:133` → `newProvidersFixture(t)` | 可转 |
| 2 | `TestProvidersProbeRecordsMeasuredThinkingTrue` | `providers_test.go:179` → 同一枚 fixture | 可转 |
| 3 | `TestProvidersDiscoverListsWhatTheServerServes` | `providers_test.go:208` → 同一枚 fixture | 可转 |
| 4 | `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless` | `providers_test.go:228` → 同一枚 fixture | 可转 |
| 5 | `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)` | `secret_test.go:667`（断言在 `:679`） | 可转 |

**② 实测复算（pre 锚两形各一枚 `-count=1 -v`，台 PRE-L／PRE-P）**。only-in-link ＝「PRE-L 名册里的 FAIL 名」减「PRE-P 的 FAIL 名」，脚本 `comm -23`，**5 枚、逐名与上表枚枚对上**（`/d/tmp/wisp124-2b4-logs/only-in-link.txt`）；反向差集（only-in-plain）**0 枚**。⇒ **与派单给的 5 相符，登记差 0**；形状也如批次 3b 所述：**一处 fixture 覆盖 4 枚顶层 ＋ 一处子测试 = 2 个递根点**。

**③ 28 枚两形都红**：`comm -12 PRE-L.red PRE-P.red` ＝ **28 枚**（`/d/tmp/wisp124-2b4-logs/both-red.txt`）⇒ 与批次 3b `next=` 那句「同包另有 28 枚两形都红」**逐数相符、差 0**，且与 AC#1 §5.3 排除在本票 132 之外的那本账同源。名单（本批一枚未碰）：`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` 顶层与其 4 枚子测试（`/cmdDoctor`、`/cmdModels`、`/cmdProviders`、`/runTextTask`）、`TestAC2…128` 之外另两枚 AC2、`TestAC4EveryLegIsNailedOrRuled`、`TestComposedGateBlocksAWriteForTwoSeconds`、`TestHostDispatchThroughTheAssembledBridge`、`TestMissingBlobFailsUnconfiguredNeverSilently`、`TestRunTextTask{FailNextIsClassified,KeyResolvesInTheStore,TextPathEndToEnd}`、`TestSecretEndToEndConfigRefResolvesAtRequestTime`、`TestSecretFailurePathsLogAndPrintNoPlaintext` 顶层与其 3 枚子测试（`bad_dpapi_decrypt_(corrupted_blob,_non-portable)`、`bad_dpapi_decrypt_(portable,_P13)`、`store_write_failure_surfaces_the_ref_only`）、`TestSecretFromStdinWritesNoIntermediateFile`、`TestSecretOverwriteIsAnnounced`、`TestSecretPortableModeUsesTicket06Seam`、`TestSecretSameNameUnderThreeEnvsIsThreeBlobs`、`TestSecretSetGetListRoundTrip`、`TestSecretUnsetRefusesWhileReferenced`、`TestTicket101{ManualSwitchSurvivesRestart,ModeSwitchUsesTheRealL2Gate,SessionGrantDoesNotCrossRestart,UnreadableModeFailsLoudlyAndStrict,UntouchedConfigRestartsAtDefault}`。
⚠ 这 28 枚的包级 `rc` 在改前改后都是 **1**（`FAIL=21 SUBFAIL=7` 恒在）⇒ 本批判据**逐枚点名**，绝不以包级 rc 论绿论红。

**④ 红因复算**：PRE-L 日志里那 5 枚的失败串**只有一种拼写**（`winsec: refusing to seal /varlink/…` ＋ `winsec: path is not provably resolved … reaches it through the link at /varlink`），且点名的是**用例自己的** `t.TempDir()` 路径（不是任何用例自己种的链接）。逐台机制字串行计数：

| 台 | `not provably resolved` | `refusing to seal` | `/varlink` | `test timed out` | `panic` |
|---|---|---|---|---|---|
| PRE-L | 26 | 26 | 27 | 0 | 0 |
| PRE-P | 1 | 1 | **0** | 0 | 0 |
| POST-L | 21 | 21 | 22 | 0 | 0 |
| POST-P | 1 | 1 | **0** | 0 | 0 |

⚠ **本包的机制字串在改后不归零，且这是对的**——与 `internal/agent`（3a 那发 `29→0`）不同形状，如实登记：软链形下那 28 枚两形都红的用例**自己也**在往底线递未解析根（它们红在别处、不红在这一句；那些根没换，按地界也不该由本批换）。三条核过：
1. **逐名归因**：POST-L 里 19 个打出 `refusing to seal /varlink/w124tmp/<TestXxx>` 的父用例名，全部落在 28 枚名册内（`comm -23 postl-seal-parents.txt`「both-red 的父名侧」⇒ **0 枚差集**）；本批那 5 枚的名字**一个都不再出现**在该串里。
2. **差值恰为本批**：PRE-L `26` − POST-L `21` ＝ **5**、`/varlink` `27` − `22` ＝ **5**，与转绿的枚数逐数吻合。
3. **普通形逐数不变**：PRE-P 与 POST-P 都是 `1/1/0`（那 1 行来自两形都红的 `…/store_write_failure_surfaces_the_ref_only`，不含 `/varlink`、在普通形也拒，属 28 那一档、本批未碰）。

⇒ 本批「归零」的口径是**逐枚点名的那 5 枚红名归零**（§2 表），**不是**整包机制字串归零；终判据复算方拿这句时要按名册拿，别拿 `grep -c`。

**⑤ 批次 3b 给这一批的三条事实自核**（我的转述也是断言，逐条重走）：
1. **名册 5 枚** ⇒ 见上面 ②，**成立**。
2. **同包另有 28 枚两形都红，一枚都不许顺手修绿** ⇒ 见上面 ③，**成立**；本批改动的 2 个递根点**没有一处**落在该 28 枚的任一可达路径上（改后名册逐名相同、28 枚红名 `diff` 前后 **0 行差**，见 §2 末）。
3. **`cmd/wisp` 可以直接接 `proc.SealableRoot`，只有 `internal/winsec` 是例外** ⇒ **成立且本批独立复看**：`grep -n "proc.SealableRoot" cmd/wisp/*.go` 在**生产码**里已有两处（`doctor.go:263` 的 `resolveDataDir`、`secret.go:148` 的 `secretLayoutOf`），两处的注释逐字写明这就是票 119「调用方解析 OS 给的答案」那条纪律；本批把**测试侧**接上同一枚叶子，不新增第三条路数。`internal/proc` 早于 `main` 且 `cmd/wisp` 本就依赖它 ⇒ 无环、无新依赖边（§7.4 双 GOOS vet 各 rc=0 为证）。**本批不重走 3b 那套包内 `resolveProbeRoot`。**

**⑥ pre 锚两形四数（`-count=1 -v`）与改后两形（同一台架、同一脚本、同一容器镜像）**：

| 台 | 形状 | `./cmd/wisp/` |
|---|---|---|
| PRE-L | 软链 | `rc=1 RUN=71 PASS=16 FAIL=25 SKIP=0 SUBPASS=22 SUBFAIL=8 SUBSKIP=0`（包级 `10.142 s`） |
| PRE-P | 普通 | `rc=1 RUN=71 PASS=20 FAIL=21 SKIP=0 SUBPASS=23 SUBFAIL=7 SUBSKIP=0`（`10.377 s`） |
| POST-L | 软链 | `rc=1 RUN=71 PASS=20 FAIL=21 SKIP=0 SUBPASS=23 SUBFAIL=7 SUBSKIP=0`（`11.019 s`） |
| POST-P | 普通 | `rc=1 RUN=71 PASS=20 FAIL=21 SKIP=0 SUBPASS=23 SUBFAIL=7 SUBSKIP=0`（`10.334 s`） |

⇒ 本批软链形红名 **5 → 0**（`FAIL 25→21`、`SUBFAIL 8→7`，恰 −4 顶层 −1 子测；`PASS 16→20`、`SUBPASS 22→23` 恰 +4 +1）；普通形**八个数改前改后逐数相同**。包级 `rc=1` **不是**本批的判据（§1 ③ 那 28 枚贡献的，两形恒在）。

## 2. 判据①：逐枚转绿且点名（`-v` 才有 PASS 名；四数只能从 `-v` 量）

四台对照（同一批用例，`-count=1 -v`）。PRE-L=改前软链、PRE-P=改前普通、POST-L=改后软链、POST-P=改后普通。逐名状态取自各台 `-v` 日志的 `--- (PASS|FAIL|SKIP)` / `    --- (PASS|FAIL|SKIP)` 行（名册 `/d/tmp/wisp124-2b4-logs/{PRE-L,PRE-P,POST-L,POST-P}.names.txt`，每份 71 行）。下表**程序化生成**（`python` 读四份名册，非手抄）。

| 用例名（子测试写作 `父/子`） | PRE-L | PRE-P | POST-L | POST-P |
|---|---|---|---|---|
| `TestProvidersProbeRecordsMeasuredThinkingFalse` | FAIL | PASS | **PASS** | PASS |
| `TestProvidersProbeRecordsMeasuredThinkingTrue` | FAIL | PASS | **PASS** | PASS |
| `TestProvidersDiscoverListsWhatTheServerServes` | FAIL | PASS | **PASS** | PASS |
| `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless` | FAIL | PASS | **PASS** | PASS |
| `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)` | FAIL | PASS | **PASS** | PASS |

⇒ 5 枚逐枚点名，POST-L 列全部是 **`PASS` 字样**（取自 `-v` 的逐名 `--- PASS` 行，不是包级 `ok`、不是 SKIP）。

改前那 5 枚**各自拿到**的句子（逐字抄自 `/d/tmp/wisp124-2b4-logs/PRE-L.__cmd_wisp_.txt`，`providers_test.go:133/179/208/228` 与 `secret_test.go:679` 是报出行号）：

- `providers_test.go:133`（`ThinkingFalse`）：`memory: create data dir: winsec: refusing to seal /varlink/w124tmp/TestProvidersProbeRecordsMeasuredThinkingFalse10120219/002: the installed risk.c26Pipeline answered "…", a spelling the floor itself refuses: winsec: path is not provably resolved, refusing to seal: … reaches it through the link at /varlink, which is not the tree this call names`
- `providers_test.go:179`（`ThinkingTrue`）／`:208`（`Discover`）／`:228`（`UnconfiguredRef`）：**同一族文案、各自的 `t.TempDir()` 路径**（`…/TestProvidersProbeRecordsMeasuredThinkingTrue2931414033/002`、`…/TestProvidersDiscoverListsWhatTheServerServes2269718987/002`、`…/TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless2464220708/002`）。
- `secret_test.go:679`（`unset name (no such blob)`）：`the failure must name the ref, got: wisp secret unset: secret: create /varlink/w124tmp/TestSecretFailurePathsLogAndPrintNoPlaintextunset_name_(no_such1097638195/001/secrets: winsec: refusing to seal …: the installed risk.c26Pipeline answered "…", a spelling the floor itself refuses: winsec: path is not provably resolved …` ⇒ 与 AC#2a 账 §5.9 行 132 摘的实拿逐字同形（它要的是「点名 ref」，拿到的是「密封拒」）。

**逐台核对（脚本 `comm`/`diff`，不靠目测）**：

- 用例名集合四台**逐名相等**：`PRE-L↔PRE-P`、`PRE-L↔POST-L`、`PRE-P↔POST-P`、`POST-L↔POST-P` 四对 `comm -3` 各 **0 行** ⇒ 既没多出一枚用例，也没少一枚（分母恒定 71 名）。
- 四台 `SKIP` 计数（顶层 + 子测试）**全 0**。
- ⚠ **28 枚两形都红的名册改前改后逐名相同**：`diff /d/tmp/wisp124-2b4-logs/both-red.txt both-red-post.txt` ⇒ **0 行、rc=0**。即本批既没把它们修绿、也没多红一枚、也没换名。
- 改后 only-in-link **0 枚**（`comm -23 POST-L.red POST-P.red` 空）；`POST-L.red.txt` 与 `POST-P.red.txt` `diff` **逐名 IDENTICAL**（各 28 名）。
