# 177-m1 —— 甲形（精确豁免）变异台件现跑读数：把 `177-c1` 表 §4.1 那张〔读码推导〕红名单变成逐名实测

- 派单＝`.scratch/wisp/dispatches/2026-09-27-220x-mutation-177-m1-shape-a-redlist-run-hot-then-restore.md`（本程＝AC#1b 的欠账）
- 票面＝`.scratch/wisp/issues/177-stamping-external-content-and-letting-the-model-re-read-a-spilled-artifact-are-mutually-exclusive-today-because-the-stub-embeds-the-host-minted-path-and-r4-scans-path-arguments.md`
- 上游被审件＝`docs/evidence/s1/177-c25-r4-path-exemption-c1.md`（195 行；**它是被告的推导，本表是去量它的**）
- 起手（`date` 现量）：`Sun Sep 27 22:29:27 CST 2026`｜branch `dev`｜`git rev-parse HEAD` = `9b4107530df050a2a8d6a918b2b38c1ccd4fbce0`
- 台件与原始日志＝`.scratch/wisp/probes/177/m1/**`（本目录只存台件文本＋日志，产码一枚不留在 `internal/`）
- 每格标〔现跑〕／〔现读〕／〔未证〕。**本表没有任何一格是从 `177-c1` 抄来的。**

---

## 0. 起手五件与排他约束〔现跑〕

| 件 | 命令 | 读数 |
|---|---|---|
| 时刻 | `date` | `Sun Sep 27 22:29:27 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` |
| 锚点 | `git rev-parse HEAD` | `9b4107530df050a2a8d6a918b2b38c1ccd4fbce0` |
| 树净 | `git status --porcelain -- internal/ cmd/` | 空（无输出） |
| 基线 | `go test -count=1 ./internal/risk/ ./internal/tools/` | `ok internal/risk 5.616s`／`ok internal/tools 17.980s` |

名册口径分开报：顶层包 **2 枚 PASS**；`-v` 逐名口径尺 `go test -count=1 -v ./internal/risk/ 2>&1 | grep -cE '^=== RUN'` → **172**，同包 `--- PASS` **171** ＋ `--- SKIP` **1**（`--- SKIP: TestSyncRegistryProbeLive`，平台条件跳过，非本程所造）；`./internal/tools/` 的 `=== RUN` → **173**，`--- FAIL/SKIP` → **0**。
**排他约束满足**：`git log --oneline -3` 顶格是 `9b410753 票 174 Progress log 追加 174-v1 验收行`，其前是 `1bd953b6 票 174 AC#2 验收（174-v1，非实现者）` ⇒ `174-r1`／`174-v1` 都已入库；本程变异窗口内编队无第二枚写手。

## 0.5 票面那一格＋那六条追加（派单要求抄进本表；逐字抄，未改一字，删节处标 `……`）

- **AC#1b**（票面 `:35`）：「**AC#1b 现跑变异台件（09-27 21:5x 由 `177-c1` 自己点的缺口立；本格是 AC#1 那句"不许推理，要跑台件"的欠账）**：在真树上把**甲形**（`internal/risk/taintmatch.go:111-129` 建索引那一圈＋`provenance.go` Mark 侧的跨度排除）**临时落一次**，量出 `177-c1` 表 §4.1 那六行里**哪几枚真红**、把逐名红集合抄进证据件，然后**逐文件还原**（`git show HEAD:<path> > <path>`；禁 `checkout`／`restore`／`revert`／`clean`），还原后 `go test -count=1 ./internal/risk/ ./internal/tools/` 复绿、`git status --porcelain -- internal/ cmd/` 为空。⚠ **过界哨（它点名、我认可并把成 AC）**：`internal/risk/provenance_test.go:136`／`:139` 那两枚 fail-closed **在甲形下不该红**——**它们若红＝豁免越界，就地停手回禀，不许改判据**……」
- **21:5x 六条追加**（票面 Progress log `:52`–`:57`，标题逐字：`### 2026-09-27 21:5x 编排者追加：\`177-c1\` 交件之后的六格判定（票面一字未抹，全部以追加为准）`）：
  1. 「票面 `:14` 那句行号已漂：原文「`task.go:229` 逐字 `…全文见 %s…`」——现读 `internal/tools/task.go:229` 是 `totalBytes := len(full)`，那句格式串在 **`:245-246`**……谁照 `:229` 去改就会改错行」。
  2. 「`:20` 末句"换窄形（只加 `task.output`）红这 2 枚"要补一枚限定：那两枚红是**票 175 落地之后**才成立的形状。**今天的树是绿的，因为压根不盖戳**——`internal/tools/bridge.go:551-553` 首行就是 `!risk.IsSensitiveSource(dec.Tool) → return`，而 `sensitiveSourceTools`（`internal/risk/provenance.go:96-101`）八枚里没有 `task.output`……」
  3. 「AC#1 不翻勾，缺口立成 **AC#1b**（见上）……`TaskRoster.Record` 生产零调用点……／`ToolResultLog.Artifact` 一写零读……／`provenance.go:662` 那一支逐字／八枚名册逐枚……」
  4. 「**"那枚集合不存在"它裁成了"不需要新名册"**：甲形不做成一张 path 列表，而是"在宿主自产的那一枚 mark 里，P 所占的窗口不入片段索引"，随 mark 一起生死……⚠ 反面它写明、我抄成写手腿的 AC：**一旦实现做成"按 scope 维护一条 path 列表"＝新造名册**；且排除跨度必须证明随 `Scope.Close` 一起掉，否则跨任务存活＝一枚新的洗戳面」。
  5. 「它纠正了我两处前提（都对，见 ①②），⚠ 但它 readings §1 有一处**读数为真、注解为假**要留此：把工作树干净注成"（174-r1 此刻不在树里）"——174-r1 的产码 `f55ddd3d` 在它起手锚 `640d30c5` **之前**就已在树里……**不要照那句注解去重算行号**」。
  6. 「交件史（与本票内容无关，但必须留）：本票的只读腿被**一封假"完成通知"**提前报过喜……**判据不变：通知不是交件，盘上才是**。」

本程对这六条的关系：第 4 条（"随 `Scope.Close` 一起掉"）本程量了＝台件 F 发，绿〔现跑〕；第 1/2/3/5/6 条是行号与排程账，本程第 1 条要补一个更精确的读数（见 §8 末）。**本程一枚框都没勾。**

---

## 1. 本程没测什么（先立这条，免得下面的绿被读成"都量过了"）

1. **端到端那一发没测**：`internal/tools/bridge.go:560` 与 `internal/tools/task.go:245` 本程**一字未改**——派单 §2 第 4 项把载具裁成"包内最小通道"，而未导出的 `markWithHostPath` 桥够不着。⇒ 夹具文件 `probes/177/c1/fixture-reverse-criterion.md` 要求的"真桥＋真 `fs.read`＋真临时目录"那三发（R-1／R-2／R-3 落在 `internal/tools/`）**本程是在 `risk` 层用同一条 P 复现的**，不是桥层读数。票 164 那两枚续读腿（`internal/tools/task_output_leg_test.go:165`／`:346`）本程**没量**（它们撞的是"175 落地＋`task.output` 进名册"那一天，`177-c1` §3 的这条判断被本程基线读数支持：今天两包全绿）。
2. **C1 面（`internal/tools/tool.go` 的 `Result`）没碰**：派单禁，本程照禁。⇒ 上游 §8 b 支"两条载具二选一"那一问本程**没有加票**。
3. **丙形（每次续读一张 L2 卡）零读数**（票面 AC#0 原状），本程没为它补台件。
4. **"整段目录级豁免"那一形没测成**：本程的 S3c 是按"窗口与前缀串重叠"建模的，路径尾部那 ≥8 枚 runes 仍在索引里，所以 R-3 那一发**没红**（关掉位置记录后转绿，见 §5 末注）。真正的"整条值命中前缀就整条放行"那一形本程**没造** ⇒ R-3 能不能打红目录级豁免＝**〔未证，仍是推导〕**。
5. **`matchText`／`Inspect`／`writeGate`／`isPathKey`／`pathTargets` 一律未动**（上游 §6 的最小路线要求如此）⇒ §4.2 那六行匹配侧红名单本程**一枚没量**，它们仍是被告推导。
6. **`agent/`、`memory/`、`cmd/` 的包级测试没跑**（派单 §0 逐包口径禁全仓 `./...`）；`BenchmarkMarkFullSource` 等基准没跑。
7. **性能/内存那一格没量**：删窗口只减不增哈希条目，但本程没测 `skip` 带来的建索引开销变化。

## 2. Q1：上游 §4.1 那六行，哪几枚真红（变异态 vs 还原态两向）

红集合的尺（每态同一条，逐名口径）：
`go test -count=1 -v -run 'TestFragmentMatchExactBoundary|TestFragmentIndexShortSourceNeverMatches|TestFragmentHashCollisionCannotFakeHit|TestMarkEmptyAfterNormalization|TestMarkUnknownSourceToolFailClosed|TestConcurrentMarkInspect|TestScopeCloseRefusesAnotherOwnersRegistration|TestM1Rig' ./internal/risk/ 2>&1 | grep -E '^(--- |ok |FAIL|PASS)'`
枚数类第三枚（tools 侧）另跑：`go test -count=1 -v -run 'TestFSReadReturnsTheFileAndTaintsIt' ./internal/tools/ 2>&1 | grep -E '^(--- |ok|FAIL)'`

| 上游行 | 点名的用例（本程现读锚点） | S1 甲形·精确〔现跑〕 | 还原态〔现跑〕 | 本程结论 |
|---|---|---|---|---|
| 1 | `TestFragmentMatchExactBoundary`（`taintmatch_test.go:29`，调用点在 `:31`）／`TestFragmentIndexShortSourceNeverMatches`（`:120`／调用点 `:124`）／`TestFragmentHashCollisionCannotFakeHit`（`:146`／调用点 `:148`） | 三枚**逐名 PASS** | 三枚 PASS | **不红**。它们的红**只挂在载具上**：S4 把 `newFragmentIndex` 的第三枚参数改成必选，`FAIL internal/risk [build failed]`，编译错 7 处（尺：`grep -cE 'not enough arguments in call to newFragmentIndex' .scratch/wisp/probes/177/m1/s4-pkg.log` → 7）＝`taintmatch_test.go:31/:51/:66/:108/:124/:134/:148`。⇒ 上游写"三枚"是**漏计**：同一次编译崩多带出 4 处调用点（属 `TestFragmentMatchAcrossTransforms`／`TestFragmentMatchCJK`／`TestNormalizeTaintDropsIgnorableChars`），且整个 risk 测试二进制一枚都跑不了 |
| 2 | `TestMarkEmptyAfterNormalization`（`provenance_test.go:121-126`，判据在 `:123-125`） | **PASS** | PASS | **不按上游那支机制红**：S5 把空判定挪到排除之后（正文恰好只剩 P 就 `return false`），红的却是本程台件 E 发（`provenance_test.go:1096`），`:121` 仍绿——它传的正文是 `"   \n\t  "` 且**不带 P**，任何只作用于"宿主声明那一段"的排除都碰不到它。行 2 的因果要改写成"**只有改了 `Mark` 的导出签名，它才以编译红的身份红**" |
| 3 | `TestConcurrentMarkInspect`（`provenance_test.go:844`，判据 `:857` `!= 8`）／`internal/tools/fs_test.go:35`（`len(got) != 1`）／`provenance_test.go:916`（`TestScopeCloseRefusesAnotherOwnersRegistration` 里 `got != 1`） | 三枚**逐名 PASS** | 三枚 PASS | **不红**，与派单那句"三枚枚数类不响"**相符**〔现跑〕。原因不是它们钝，而是甲形**没为 P 另加一枚 mark**——本程用台件 E 发的第二枚断言（`ScopeTaints` 枚数＝1）正面钉住"正文只剩 P 也只记一枚"。⇒ 行 3 那句"**只有把豁免做成'为 P 另加一枚 mark'才红**"本程**没能反证**：那枚做法在包内载具下不产生任何被测量的 mark（见 §1 第 1 条） |
| 4 | `TestMarkUnknownSourceToolFailClosed` 的两枚 fail-closed（`:136`／`:139`） | **两枚 PASS** | 两枚 PASS | 见 §4，**四枚可测态全绿** |
| 5 | 契约文字面（`SPEC-06.md:70`／`:37`＋`taintmatch.go:11-15` 那句 REJECTED） | 未动一字（本程没改任何文档与注释里的契约句） | 同 | **不是测试**，本程不裁；但读数摆明：`provenance.go:472` 那句"Returns false only when the content has no matchable characters after normalization"在 S1 下**仍然逐字为真**（台件 E 发绿＝只有 P 的正文仍被记录），在 S5 下**为假**（E 发红）⇒ 这句要不要改，取决于实现是否把"可匹配字符"重定义为"排除 P 之后剩下的字符"，那是人工批准面，本程只报机制 |
| 6 | d22scan 禁面（`AGENTS.md §1.2`／`PLAN.md:1288`） | 变异态 `sh scripts/d22scan.sh` → **rc=0** | rc=0 | **不红**〔现跑〕。本程定位 P 用的是 `runeIndexOf`（纯 rune 比较），没在 `risk.PathResolver` 之外用 `filepath.Clean\|Abs`；两向 `examined` 枚数逐字相同，见 §6 |

**一句话答 Q1**：在派单裁定的那枚最小载具（未导出参数、导出面与桥不动）下，**上游 §4.1 那六行一枚都不真红**；其中行 1、行 2 的红是**载具红**（改了签名就编译红），不是甲形的语义红，而行 3 的前提甲形根本不满足。⇒ 这张六行表作为"甲形的代价清单"是**空的**，作为"载具选择的警告清单"才是成立的。

## 3. Q3 形状自证：宽一档的变异真的响了吗（防"恒不变＝装饰"）

| 态 | 放宽法 | 响的是谁（逐名） | 命令（同一把尺） |
|---|---|---|---|
| S2 | 排除窗口＝**整段 mark** | `--- FAIL: TestM1RigBRestOfHostMarkStillIndexed`（同态 A/C/D/E/F 全 PASS；包级 `FAIL risk 5.450s`／`ok tools 15.881s`） | §2 那条 `-v -run` |
| S3a | 排除**按位置**推给全树所有 mark | 六发**全绿** ⇒ 这放宽法打不动任何一发（本身是条否定读数：位置粗筛漏得多） | 同上 |
| S3b | 排除**按值**（名册形）推给所有 mark | `--- FAIL: TestM1RigCLegForeignMarkStillHits`（`provenance_test.go:1065` "R-1: P carried by FOREIGN content must still hit — the exemption left its own mark"）；A/B/F PASS、D PASS | `go test -count=1 -v -run 'TestM1Rig\|...' ./internal/risk/` |
| S3c | 名册里放**目录前缀**（同一棵树里 D 发曾 FAIL） | D 发 FAIL 是**位置记录与值记录叠加**造成的：把 `m1RecordSpan` 的位置记录关掉后 D 转绿 | 同 S3b |
| S5 | 空判定排在排除**之后** | `--- FAIL: TestM1RigEPathOnlyMarkStillRecorded`（`provenance_test.go:1096`） | 同 §2 |
| F 发 | 生命周期（未放宽，随 `Close` 掉） | `--- PASS: TestM1RigFExemptionDiesWithScope`（四态皆 PASS） | 同 §2 |

⇒ **判据不恒真**：本程至少造出四枚各自当场响的放宽法（S2 响 B、S3b 响 C、S5 响 E，加上"什么都不豁免"会响 A——A 发在未加豁免的树上等价于今天的红，见 §4 末）。甲形（S1）六发全绿＋现有 13 枚点名用例全绿。
⚠ 但**"R-1/R-3 那一族在整段放宽下响"这一句不成立**：整段放宽响的是"同一段 mark 的其余片段"那一发（B），R-3（D）在两种粗筛下都不响。上游 §5 那句"整片都豁免（红 R-1/R-3）"里，**R-1 对按值豁免成立、R-3 本程未证**。

## 4. Q2 过界哨：`provenance_test.go:136`／`:139` 逐字读数

现读锚（本程逐字读过盘上那两行）：`:136` `t.Fatal("unknown-source marks must still be recorded (fail-closed)")`／`:139` `t.Fatal("taint from an off-list source must still gate")`，二者同属 `TestMarkUnknownSourceToolFailClosed`（`:128`起）。

| 态 | `go test -count=1 -v -run 'TestMarkUnknownSourceToolFailClosed' ./internal/risk/` |
|---|---|
| S1 甲形·精确 | `--- PASS: TestMarkUnknownSourceToolFailClosed (0.01s)` |
| S2 整段 mark | 同上 PASS（S2 那次 `-v` 全量点名里逐名可见） |
| S3a/S3b/S3c 合并态 | PASS |
| S3b 归因复跑 | `--- PASS: TestMarkUnknownSourceToolFailClosed (0.01s)` |
| S5 | PASS |
| 还原态 | 包级 `ok internal/risk 5.513s`，点名 PASS |

**没有触发停手条件**：两枚 fail-closed 在每一枚可测态下都绿，`markWithHostPath` 从未在 P 之外的地方删窗口。
⚠ **同时报一条对上游的判断**：`177-c1` §4.1 行 4 写"`:136` 红＝'长得像路径的串都不入索引'；`:139` 红＝'按值/按集合在匹配时放行'"。本程把这两枚过界形都造出来了（S3b 就是按值名册），那两枚**仍然绿**——理由现读得出：该用例喂的值是 `marker`（`provenance_test.go:16` `MARKER-QZX-98WVE7-TAINTED-SOURCE`，里面没有 `/`、没有路径段），而且它压根不经过任何声明 P 的载具。⇒ **那两枚不是"豁免越界"的哨**，它们守的是"名册外的工具仍要盖戳、盖了仍要命中"那一格；越界要看得住，靠的是 §3 的 C 发（R-1）那一族。**本程未因此放宽或改动那两枚的任何一字**（AGENTS §1.1＋票面 AC#4）。

## 5. 次菜（乙形那一发结构性冲突）：已证〔现跑〕，但证在 risk 层

上游 §4.2 行 6 的断言是：`Inspect` 只看得到候选值 ⇒ 合法续读与外来说到同一条 P 时**候选值逐字相同**，按值放行必然同时放过外来那一发。
本程用 S3b 直接量到它：两发的候选值都是同一枚常量 `m1Path`（A 发 `Inspect("task-1","notify",{"text": m1Path})` 来自宿主自产 mark，C 发同一枚 `Inspect` 同一串文本，外来 mark 逐字带它）。在"按值豁免"的名册下：A 发**如期干净**（PASS＝放行成功），C 发**同时没了命中**（FAIL 文本 `R-1: P carried by FOREIGN content must still hit`）——**同一枚候选值，两次调用给出同一个"放行"结论**。⇒ 乙形的那枚冲突不是推导，是读数。
末注（诚实标）：D 发（R-3 过界哨）在同一枚名册下**没红**，因为本程把"目录前缀"建模成"与前缀串重叠的窗口才删"，路径尾部仍留 ≥8 runes 在索引里；把建模改成"整条值放行"才会红，而那一步本程没造 ⇒ **R-3 能打红目录级豁免＝〔未证〕**，别再把它当读数引用。

## 6. Q4＋§5 门禁：两把尺两向（变异态 vs 还原态），名册口径分开

| 尺 | 变异态（S1） | 还原态 | 命令 |
|---|---|---|---|
| `sh scripts/d22scan.sh` | rc=0，`clean - no D22 ban violations`；`bans #1-5 internal/=207`、`bans #1-5 cmd/=23`、`ban #6 frontend/=85`、`ban #7 internal/tools/=20`、`ban #8 internal/=426`、`ban #8 cmd/=45`、`ban #8 design/=39` | rc=0，同上一枚数逐字相同 | 日志 `.scratch/wisp/probes/177/m1/s1-d22scan.log`／`r-d22scan.log`；`tail -2` 两文件 |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | rc=0；`# 腿数＝14 声明与实测不符＝0`；`# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`；`# 腿数断言＝相符（名册上每一腿都记了账）` | rc=0，同一枚三行逐字相同 | `grep -E "腿数" .scratch/wisp/probes/177/m1/{s1,r}-gates.log` |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | rc=0；`--- PASS` **34**／`--- FAIL` **0**／`--- SKIP` **0**／`=== RUN` **76** | rc=0，四数逐字相同 | 尺（对两枚日志各跑一遍）：`for L in .scratch/wisp/probes/177/m1/s1-runtests.log .scratch/wisp/probes/177/m1/r-runtests.log; do grep -cE '^--- PASS' $L; grep -cE '^--- FAIL' $L; grep -cE '^--- SKIP' $L; grep -cE '^=== RUN' $L; done` |
| **名册差集** | — | 空 | `grep -oE '^--- PASS: [A-Za-z0-9_/]+' <log> \| sort` 两份后 `comm -3` → **无输出**（34 行 vs 34 行，`wc -l` 各 34） |

⚠ 两处自我更正（口径，不是上游的错）：
1. 我第一次用 `grep -c '\[no tests to run\]'` 得到 **1**，那枚命中其实是 `runtests.sh` 自己那行总结（`s4`/`r` 两份日志第 236 行逐字为 `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`）。按脚本自己的记账，**`[no tests to run]`＝0**，与 `177-c1` §7 相符。
2. `177-c1` §7 那句"未跑 `gate-clauses.sh`（别家在树里）"在本程不成立：本程独占树，两向都跑了。

包级四数（还原态，逐包，禁 `./...`）：`go test -count=1 ./internal/risk/ ./internal/tools/` → `ok risk 5.513s`／`ok tools 15.410s`＝**顶层 2 枚 ok、0 枚 FAIL**；`-v` 口径 risk `=== RUN` 172／`--- PASS` 171／`--- SKIP` 1（`TestSyncRegistryProbeLive`），tools `=== RUN` 173／FAIL+SKIP 0。锚点 `9b4107530df050a2a8d6a918b2b38c1ccd4fbce0`。

## 7. 还原证明（这一步先于本文件落盘）

```
for f in internal/risk/taintmatch.go internal/risk/provenance.go internal/risk/provenance_test.go; do git cat-file blob "HEAD:$f" > "$f"; done
git status --porcelain -- internal/ cmd/     # 输出：（空）
git diff --numstat -- internal/ cmd/         # 输出：（空）
go test -count=1 ./internal/risk/ ./internal/tools/   # ok risk / ok tools
```
用的就是票面 AC#1b 指定的逐文件还原（`git cat-file blob HEAD:<path> > <path>`，等价于它写的 `git show HEAD:<path> > <path>`）；**没用** `checkout`／`restore`／`revert`／`clean`／`switch`／`merge`／`worktree`，**没跑任何删除命令**，没建第二份 `.git`。
现场那批脏件（`.gitignore`、`.scratch/wisp/probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`design/**`）本程未碰、未提交、未评论。

## 8. 交件报告剩下的格子

**被拒／没成功的调用**：取数前后各 0 枚被权限拒绝；失败 1 次——第一次给 `internal/risk/provenance_test.go` 追加台件时 `old_string` 锚点写错（我误把上一段的 `t.Fatalf("ID() = %q", ...)` 当成文件末尾），Edit 报 0 occurrences，改成文件真实末三行后成功。那次失败发生在**变异窗口内、产树未脏之外**（当时该文件已带着 S1 变异），随后一切读数都在成功写入台件之后取。
**有没有跑过删除命令**：没有。没跑 `rm`／`git clean`／`restore`／`checkout .`；台件临时件只建不删。
**凭据值抄录**：零枚（全文只出现工具名与键名，无密钥／令牌／URL 凭据）。

**伪授权两栏（各带出处）**

| 栏 | 内容 | 出处 | 本程态度 |
|---|---|---|---|
| 甲：转述在场、**我未复核** | "owner 原话'行，都按推荐就完事了'＝批准甲（精确豁免），台账 `A349`" | 派单 `:6`；票面 `:3`（21:5x 追加那句"owner 已批 `Q-61`＝甲"） | 我**没读** `docs/reports/pending-and-issues.md`（不在写面）⇒ 不声称 `A349` 在盘上。我只把"批准"当作**形状边界**用：目录级／前缀级／配置项级一枚没造，C1 面与 `Loop` 面一枚没加收 |
| 甲′：派单里我复述了但没验的数 | 派单称"HEAD 是 `9b410753`、编队安静、`174-v1` 已交" | 派单开头编排者注；本程 §0 现量 `git rev-parse HEAD` ＋ `git log --oneline -3` | **相符**，已升为〔现跑〕 |
| 乙：在场而互相矛盾 | 派单说上游表"§4.1 六行红名单"，本程现跑结果＝**那六行在裁定载具下一枚不红**；派单说"三枚枚数类不响"（**相符**）、"那两枚 fail-closed 不该红"（**相符**，但**理由不成立**：它俩对按值／按形状的豁免根本钝，见 §4 末）；`177-c1` §4.1 行 1 说"三枚直接编译红"（现跑＝**7 处调用点、整包 build failed**，且红的是载具不是形状） | `docs/evidence/s1/177-c25-r4-path-exemption-c1.md:84-93`；派单 §2／§7；本程 §2／§3／§4 | 我按**实测**报，不改判据、不为对齐任何一句去动测试。差异逐条写在 §2／§4 |
| 乙′：行号账 | 票面 `:14` 的 `task.go:229` 已漂（编排者 21:5x 追加第 1 条改成 `:245-246`）；本程现读：那句格式串**整枚在 `internal/tools/task.go:245`**，同款在 `internal/agent/spill.go:141`；`:229` 现读是别句 | 尺：`grep -n "全文见" internal/tools/task.go internal/agent/spill.go` | 追加更正（不动票面一字）：**精确到 `:245` 一枚行号**。这是比派单/票面更窄的一处，属"我量到的和写的不一样"，不是谁造假 |

**next=（含甲形落地前还缺哪几格）**

编排者这一轮的裁定可以下了：**甲形在"包内最小载具"下零既有用例红，形状本身没有隐藏的 casualties；真正的红全在载具上**。落地前还缺：

1. **载具那一枚批准**（`177-c1` §8 b 支，本程加一条实测）：改 `newFragmentIndex` 为必选参数＝**risk 测试整包编译红**（7 处）；改 `Mark` 导出签名＝还会把 `internal/tools` 一起拖红（本程**没测**这一支，见 §1 第 1 条）。⇒ 要么走"未导出参数＋桥侧新枚未导出包装"（本程量过，零红），要么走 C1 加收字段（人工批准）。
2. **桥侧端到端那一发**：本程的六发在 risk 层；落 `internal/tools/` 用真桥复现 R-1/R-2/R-3 才是 AC#2/AC#3 要的常驻判据（夹具文件已给形状）。
3. **R-3 对目录级豁免的打力**〔未证〕：要把"整条值命中前缀即放行"那一形造出来跑一发，才能说 R-3 是名册与目录的分界。
4. **行 2 那句契约文字**（`provenance.go:472`）到底改不改：本程给了机制读数（S1 下仍为真，S5 下为假），改契约＝人工批准，摆给 owner。
5. **票 176 的起跑口**与 **AC#0 丙形零读数**：两格本程未动，排程账仍在票面 `:40`／AC#0 那里挂着。
6. 本表 §2 的红名单可以直接把上游 §4.1 那张〔读码推导〕表**换成现跑**；但**别把"六行全不红"读成"零风险"**——风险集中在 §3 那四枚放宽法上，它们每一枚都对应一条要写进写手腿 AC 的反向判据。

**本程没勾任何框**（AC#1／AC#1b 由编排者裁）。
