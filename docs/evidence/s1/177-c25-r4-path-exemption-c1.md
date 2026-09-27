# 177-c25-r4-path-exemption-c1 —— 「宿主自己写下过的路径」精确豁免：可行性与射程裁定（只读设计核，零产码）

- 派单：`.scratch/wisp/dispatches/2026-09-27-214x-readonly-177-c1-where-does-the-host-minted-path-set-live.md`（锚 `640d30c5`）
- 票面：`.scratch/wisp/issues/177-stamping-external-content-and-letting-the-model-re-read-a-spilled-artifact-are-mutually-exclusive-today-because-the-stub-embeds-the-host-minted-path-and-r4-scans-path-arguments.md`
- 起手：2026-09-27 21:39:56 CST｜branch `dev`｜HEAD `640d30c52f66359023694c2e7c252442d53fb0b1`｜`internal/tools/`＋`internal/risk/`＋`cmd/wisp/` **起手干净**（174-r1 不在树里）
- ⚠ 这张表的名字历史上被伪造过一次（台账 `A347` 那段自纠＝票面 `:35`）。**这次它真在盘上**，落点见 §8。
- 每格标〔现跑〕／〔现读〕／〔读码推导〕／〔票面自陈，未复测〕。〔读码推导〕＝本程没跑台件，写手腿必须自己量。

---

## 1. 票面五格（连行号原样抄，AC#1 至今未勾）

| 格 | 票面行 | 逐字要点（删节号处原句更长） | 本程动了没 |
|---|---|---|---|
| AC#1 | `:34` | 「只读设计核：把四支候选各量一遍"要动哪一枚冻结件、动了之后哪几枚既有用例红"（**不许推理，要跑台件**…）」＋`:35` 编排者自纠「本格没有裁完，且我一度把不存在的结果写了进来…本票至今**一枚候选都没被独立量过**」 | 部分：三件裁了，**红名单是〔读码推导〕不是现跑**（本程禁写产码 ⇒ 跑不了变异）。缺口记在 §8 `next=` a 支 |
| AC#2 | `:36` | 「反向判据（这枚最值钱，未修码上今天不可能响）：无论最终选哪一支，落下的判据必须包含"**外来内容里出现某条路径、模型随后去读它 => 仍然要命中 R4**"那一发…任何修法若让这一发也变安静，**就是洗戳**」＋「与 AC#3 成对，缺一形＝装饰」 | 裁完：能写成常驻件，形状与夹具在 `probes/177/c1/fixture-reverse-criterion.md` |
| AC#3 | `:37` | 「正向判据：修完之后，`task.output` 的成功结果**必须有来源标记**（票 175 的 canary 本体…9301 字节），且票 164 那两枚续读腿**复绿**。**两向都要贴。**」 | 裁完：两向成对才是判据（§5 与台件里的 R-2／R-3 两支） |
| AC#4 | `:38` | 「零放宽自证：本次交付**不许**新增任何配置项／目录级豁免／`allowlist.txt` 条目；`thresholds.go`、golden、审批超时常量一字节不动…」 | 本程零放宽〔现跑〕：`sh scripts/d22scan.sh` rc=0，未碰任何冻结件 |
| AC#6 | `:40` | 「排程约束继承：本票落地那天，**票 176（后台起跑口）不得早于本票合入**」 | 本程发现一枚**同向硬约束**：176 之前，`wisp run` 端到端腿根本跑不出来（§7） |

派单地界（`:5`）复述一遍以免走偏：批的只到「**精确豁免**」这一形状，**目录级／前缀级／配置项级一律不在批准范围内**。
本程结论没有踩出这一形（§6 甲形），但**它要动 `internal/risk/**` 与 `internal/tools/bridge.go`，且极可能要动 `internal/tools/tool.go` 的 `Result`**（C1 面）⇒ 那是人工批准，不是写手腿自己决定（§8）。

---

## 2. 三件裁第一件：「宿主自己写下的路径集合」今天存不存在？

**一句话：不存在"可供 C25 在这一轮里查询"的集合。** 有两处**局部痕迹**、一处**明令不许拿来做集合**：〔现读〕

| 在场结构 | 位置 | 它是什么 | 它能不能当那枚集合 |
|---|---|---|---|
| `ToolResultLog.Artifact` | 字段 `internal/agent/loop.go:76`；唯一写入点 `:708`（`log.Artifact = sp.Path`，来自 `internal/agent/spill.go:145`）；累积在 `Result.ToolLog`（`:714`） | 每一次 D15(3) 落盘，宿主都在这条 per-call 记录上写下过那条路径 | **不能直接用**：本程现量（尺见 `probes/177/c1/readings.md` §3）——该字段**一写、零读**（`internal/`＋`cmd/` 生产码里除声明外无任何读取点），且它是**一次 Run 的结果面**，运行中拿不到；C25 要在 `Inspect` 当时判定 |
| `TaskRoster.byTask` 的 `TaskOutput.ArtifactPath` | `internal/tools/task.go:110-113`＋`:98-101`；消费点 `:242-246`（桩里那句 `…全文见 %s…`）；`:90-92` 逐字「the host - not the model - fills it」 | 最接近"宿主写下过哪几条"的**现成名册**，语义正是这件事 | **三个不**：① 覆盖面只有 `task.output` 的后台任务，不含任意工具调用的 D15(3) 落盘；② **生产码零填充**——`Record`（`:123`）无生产调用点（现量尺见 `probes/177/c1/readings.md` §3：`.Record(` 排 `_test.go` 后**为空**），`cmd/wisp/run.go:361` 只 `NewTaskRoster()`、`:362` 只注册条目 ⇒ 今天真机上这枚集合**恒空**（那就是票 176 缺的起跑口）；③ **没有枚举器**——公开面只有 `Look(taskID)`（`:139`）与 `Count()`（`:152`），而 `Count` 的注释（`:149-152`）逐字写明它存在只为宿主自查、「giving task.output a "list everything" mode would smuggle that row back in」，那张被不许捞回的表＝`task.list`，DEFERRED 于 `PLAN.md:1531` |
| `Spiller` | `internal/agent/spill.go:34-43`（字段只有 `dir`/`Budgets`/`mu`/`sequence`）、`:105-106`（`name`/`path` 算出来就地写文件）、`:150-155` | 它**算出**名字（`artifactName` `:184-199`）但**不留账** | 不能：它今天不持有任何历史路径集合，只持有一个计数器 |
| `memory.Store.ListArtifacts` | `internal/memory/artifacts.go:96`／`:112 listArtifactsDir` | 扫 artifacts 目录得到全量文件 | **禁止用作集合**：派单 `:12` 明令「不许提出"扫目录当作豁免集合"」——那是把豁免从"宿主写下过的几条"悄悄放大成"整棵树"，AC#4 的目录级豁免正是同一形状 |

依赖方向（决定"集合必须被推进 risk、不能被 risk 捞"）：〔现读〕`internal/risk/**` 的 import 只落在 `observe`/`plugin`/`winsec`
（尺：`grep -rn "wisp/internal" internal/risk/*.go | grep -o "wisp/internal/[a-z]*" | sort -u`），**不 import `internal/tools`／`internal/agent`**。

**结论（派单 §1.1 要的那句"能不能不新增名册"）**：**能不新增名册**，但**不是**靠上面任何一枚现成结构当集合——
而是**根本不要"集合"这个形状**：把"宿主写下过 P"表达成「**在宿主自己产出这枚桩的那一枚 mark 里，P 所占的那一段窗口不入片段索引**」。
它是 mark 局部的、随 mark 一起生死（`taintMark` `provenance.go:233-238`，scope 关闭 `:425`／`DisposalScope` 见 `provenance_test.go:305`），
天然不会跨 scope 泄漏（`TestScopesNeverInherit` `:296` 的形状不破），也不需要给 `TaskRoster`/`Spiller`/`Loop` 加收任何枚举方法。
⚠ 反过来：一旦实现做成「按 scope 维护一条 path 列表」（`Provenance:242-262` 加字段、`scopeReg:263-267` 加一张表），
**那就是新造名册**，并且要额外回答生命周期与"谁有权往里塞"两问 —— 本程明确不推荐。

---

## 3. 撞车到底撞在哪一发（先把病灶钉准，两侧红名单才有针对性）〔现读〕

- `internal/risk/provenance.go:662` 那一支逐字：`if !gateOpen && !isPathKey(k) { continue }` ——
  **只有"非路径键"在安静本地写时被跳过**，路径类参数**从不豁免**。派单前提**成立**。（`isPathKey` 键表 `:152`＋`:821-829`）
- 盖戳发生在**桥**里：`internal/tools/bridge.go:551-564`（`b.prov.Mark(dec.TaskID, dec.Tool, origin, res.Text)` 在 `:560`），
  且 `:552` 的门是 `!risk.IsSensitiveSource(dec.Tool)` 就直接 return。
  ⚠ **今天 `task.output` 不在那份名册里**：`sensitiveSourceTools`＝`provenance.go:96-101` 八枚（`:83-91` 常量表），
  逐枚是 `fs.read/search.content/clipboard.read/system.get/web.fetch/doc.read/screen.capture/asr.transcript`。
  ⇒ 所以**今天的基线是绿的**（§7 现跑），冲突是**票 175 那一步落地之后才成立**的未来形状，不是当下的红。
  这条是本程补上的一格事实，票面 `:20` 只说"窄形（只加 `task.output`）红这 2 枚"，没说清"今天压根不盖"。
- 两条"宿主写下路径"的产地，**只有一条真会进 mark**：
  - `internal/agent/spill.go:141-142`（D15(3) 中括号）：桥在**全文**上先盖戳（`bridge.go:560` 拿到的是完整 `res.Text`），
    循环里之后才由 `internal/agent/loop.go:702` 落盘并把桩塞回上下文 ⇒ **这条路径今天并不在索引里**；
    外来正文一旦逐字提到某条 artifacts 路径，那才是 `Inspect` 命中的正解（AC#2 的形状，本该命中）。
  - `internal/tools/task.go:245-246`：`task.output` **把桩当结果交回去** ⇒ 一旦 175 把 `task.output` 纳入盖戳，
    桩正文里的 `rec.ArtifactPath` 就真的进了外来索引 ⇒ **这就是撞车那一发**。
  ⇒ 结论：豁免要救的那一发**只在 `task.output` 这一支**，射程天然比票面想象的小一格。
- ⚠ **行号漂移**（要报）：票面 `:14` 写「`task.go:229` 逐字 `…全文见 %s…`」——现读 `task.go:229` 是 `totalBytes := len(full)`；
  那句格式串在 `:245-246`。派单 §1.1 指 `spill.go` 没错，但谁照 `:229` 去改就会改错行。

---

## 4. 三件裁第二件：索引侧 vs 匹配侧（两侧红名单）

### 4.1 索引侧（甲形：只在那一枚宿主自产的 mark 里，把 P 的窗口排除在片段索引之外）

会动的行（**只列位置，不写码**）：
`internal/risk/taintmatch.go:101-131`（`newFragmentIndex` 建哈希那圈 `:111-129`；`:89-94` 结构体保持 `src` 全文，
命中仍走 `strings.Contains(f.src, w)` 复核 `:148-153` ⇒ **不会产生假命中**，被删的只是哈希项），
`internal/risk/provenance.go:468-473`（`Mark` 的"Returns false only when …"那句契约文字，属 SPEC-06 语义面）、
`:474`（签名）、`:480-488`（**必须先归一化再定位 P 的跨度**：`normalizeTaint` `taintmatch.go:49-65` 会丢空白与零宽字符，
用原始下标定位一定偏），`:492-496`（构造 `taintMark`），
`internal/tools/bridge.go:551-564`（P 怎么到 `mark`），
`internal/tools/task.go:242-246`（P 已知处）。

红名单（〔读码推导〕，写手腿要现跑复核）：

| # | 会红的既有用例 | 位置 | 在什么做法下才红 |
|---|---|---|---|
| 1 | 片段索引直测三枚：`TestFragmentMatchExactBoundary`／`TestFragmentIndexShortSourceNeverMatches`／`TestFragmentHashCollisionCannotFakeHit` | `internal/risk/taintmatch_test.go:29`／`:120`／`:146` | 只要 `newFragmentIndex` 签名/返回形状变——三枚**直接编译红**（它们手搓这枚函数） |
| 2 | `TestMarkEmptyAfterNormalization` | `provenance_test.go:121-126` | 若把排除做在**空判定之前**（`:475-479` 那一支），正文恰好只剩 P 的 mark 会返回 false，把"`Mark` 只在归一化后无字符才 false"那句改掉 |
| 3 | 枚数类：`TestConcurrentMarkInspect`（`:857` `!= 8`）＋`internal/tools/fs_test.go:35`（`len(got) != 1`）＋`provenance_test.go:916`（`!= 1`） | 见左 | **只有把豁免做成"为 P 另加一枚 mark"** 才红（mark 计数被抬高）。同一段注释在 `internal/tools/bridge_scope_open_ticket158_test.go:22` 已把 `fs_test.go:35` 这把尺的敏感面点过一次 |
| 4 | **`TestMarkUnknownSourceToolFailClosed` 的两枚 fail-closed** | `provenance_test.go:128-144`：`:136` 「unknown-source marks must still be recorded (fail-closed)」／`:139` 「taint from an off-list source must still gate」 | **甲形下这两枚不该红**（豁免只作用于宿主声明的那一段，`mystery.reader`/`marker` 那两发与 P 无关）。它们红＝豁免越界的信号：`:136` 红＝"长得像路径的串都不入索引"；`:139` 红＝"按值/按集合在匹配时放行"。⇒ **它们是那把闸，不是 casualties**；红了就回禀，不许改测试（AGENTS §1.1、票面 `:24`） |
| 5 | 契约文字面（**人工批准，不是测试**）：`SPEC-06.md:70`（§5 C25 六个外泄通道）＋`:37`（R4 行「不受 D45 会话授权覆盖」）＋`taintmatch.go:11-15` 那句「Exact per-token taint tracking is REJECTED (16.9#1) and must not be resurrected here」 | 见左 | 索引侧**按内容选择性删窗口**，天然贴着那句 REJECTED。这一格必须摆给 owner 看，不能由写手腿自行解释成"只是局部优化" |
| 6 | d22scan 禁面 | `AGENTS.md §1.2`／`PLAN.md:1288` | 定位 P 的跨度若用手写路径归一（`filepath.Clean/Abs` 在 `risk.PathResolver` 之外）＝**ban 2 当场红**；`scripts/d22scan.sh` 是本程现跑 rc=0 的那把尺 |

### 4.2 匹配侧（乙形：`Inspect`/`matchText` 遇到等于宿主写过的那几条路径就放行）

会动的行：`provenance.go:598-692`（`Inspect` 全体）、`:676-680`（候选循环第一个命中就 return）、
`:713-727`（`matchText`：`:721-724` **按 oldest-first 逐枚比对、第一枚命中即返回**，见 `:589` 那句确定性说明）、
`:242-262`（`Provenance` 得挂一张 per-scope 路径表＝**新名册**）、`:263-267`（`scopeReg` 生命周期）、
`:731-746`（`Detector.TaintHit` 用 `tool=""` 走通用扫描）、`:781-818`（`writeGate`）＋`:821-829`／`:920-937`（路径键面）。

红名单（〔读码推导〕）：

| # | 会红的既有用例 | 位置 | 为什么 |
|---|---|---|---|
| 1 | `TestDetectorUnknownToolScansEverything` | `provenance_test.go:281-292` | 通用面逐字承诺"任意键名/嵌套/字节都扫，fail-closed"；按值放行就是给这条缝开一个洞 |
| 2 | `TestNamedChannelParamEscape`（10 支子用例） | `provenance_test.go:391-…`（子例 `:402-411`，其中 renamed/abbreviated/extra key/nested/array/byte/deep map 四支最易被"按键名/按路径形状"放行漏掉） | 与 adversarial 报告 B-2／M-7／C-3 同族：**豁免不许跟着参数名或参数形状旅行**（`:591-596`、`:755-780` 两段注释逐字钉着） |
| 3 | `writeGate` 四枚：`TestWriteGateNotSelectedByPayloadKey`／`TestWriteGatePlainLocalWriteNotFlagged`／`TestWriteGateEveryPathTargetJudged`／`TestWriteGateAllPathsNonSyncStaysExempt` | `provenance_test.go:485`／`:557`／`:624`／`:686` | `:662` 那一支已经拿 `isPathKey` 做判据；再在匹配侧加一层"按路径值放行"＝两处路径知识互相打脸，M-7/C-3 那两个反例形状会重开 |
| 4 | R4 端到端与升级不被静默五枚：`TestR4EndToEndViaAssessor`／`TestFourChannelExfilSuite`／`TestRuleTaintR4`／`TestFusionR4BlocksSessionOverride`／`TestTicket90TaintEscalationNeverSilenced` | `provenance_test.go:229`／`:169`、`rules_test.go:127`、`assessor_test.go:78`、`internal/tools/ticket90_test.go:357` | 匹配侧的"放行"若落在 `assessor`/`ruleTaint`（`rules_gateway.go:101-115`）一层，就是把 R4 做成可协商档；票面 `:24` 与 AC#4 禁的正是"为了让门安静去改门" |
| 5 | **`provenance_test.go:136`／`:139` 那两枚 fail-closed** | 同上 | 同样**不该红**（值 `marker` 不是 P）。它们红＝匹配侧被做成了"路径形状一律放行"或"集合命中即放行"这种过界形 |
| 6 | 结构性致命伤（**不是红名单，是判据不可能成立**） | — | `Inspect` 只看得到候选值：合法续读与外来说到同一条 P，**候选值逐字相同** ⇒ 按值放行的匹配侧**无法同时满足 AC#2**（票面 `:36`）。要在匹配侧救 AC#2，就得区分"P 这次来自哪一枚 mark"⇒ 等于把索引侧的排除重做一遍，还要多背一张名册 |

---

## 5. 三件裁第三件：反向判据能不能常驻（结论：能）

- 能。形状、三发断言（R-1 反向／R-2 正向对照／R-3 "同名不同目录"过界哨）、夹具字节、接缝可用性，全部写在
  `.scratch/wisp/probes/177/c1/fixture-reverse-criterion.md`。
- 落点：`internal/tools/`（与票 164 两枚续读腿同族，复用真桥具 `task_output_leg_test.go:35-50`，零 mock）。
  **三面注入接缝（C5 golden／C17／CLI `wisp run`）里今天一发都跑不出这枚端到端**，理由见 §7 b 支。
- 不恒真的机制（一句）：反向那发的命中必须来自**另一枚 mark**（外来那枚），
  而豁免只碰宿主自产那枚 mark 里 P 的窗口；两发同批、走同一枚 matcher，
  才让"什么都不豁免"（红 R-2）与"整片都豁免"（红 R-1/R-3）两种偷懒各自当场响。
  票面 `:36`/`:37` 的"缺一形＝装饰"就是这个配对要求。

---

## 6. 推荐

**推荐索引侧的甲形（mark 局部排除），明确不推荐匹配侧。** 三条理由（每条都指得到行号）：
1. 匹配侧按值放行与 AC#2 结构冲突（§4.2 第 6 行），做到最后要么洗戳要么把索引侧重做一遍；
2. 匹配侧要新背 `Provenance`/`scopeReg` 两处状态（`provenance.go:242-262`/`:263-267`）＝派单 §1.1 不许的"新名册"，
   还多一个"谁能往集合里塞"的授权面（AGENTS §1.2「由面板侧来源的 L2『允许』」禁令同一族）；
3. 索引侧的红名单集中在**三枚直测函数**（`taintmatch_test.go:29/120/146`）＋一枚语义句（`:471-473`），
   是**可以逐枚摆给 owner 的清单**；匹配侧的红名单集中在**契约承诺**（`:591-596`、`:755-780` 两段对抗结论），不是清单。

最小改动落点（**行，不是码**）：
`internal/risk/taintmatch.go:101-131`（建窗口那一圈 `:111-129`；`:89-94`/`:148-153` 的 src 复核面**不动**）
→ `internal/risk/provenance.go:474`＋`:480-488`（先归一化后定位跨度）＋`:492-496`＋`:468-473`（那句契约文字要改，属人工批准）
→ `internal/tools/bridge.go:551-564`（P 送达；两条载具见 §7 b 支）
→ `internal/tools/task.go:242-246`（P 已知处）。
`Inspect`／`matchText`／`writeGate`／`isPathKey`／`pathTargets` **一律不动** ⇒ §4.2 那六行不适用。

---

## 7. 门禁读数（本程自己重跑，逐包，禁全仓 `./...`）〔现跑〕

- `go test -count=1 ./internal/risk/` -> **ok 5.710s**（FAIL 0）
- `go test -count=1 ./internal/tools/` -> **ok 16.651s**（FAIL 0；起手 `internal/tools/` 未被在飞改动弄脏）
- `sh scripts/d22scan.sh` -> **rc=0**；examined 230 production Go files（internal/=207、cmd/=23、ban#7 internal/tools/=20）
- `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` -> **rc=0**；top-level **PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0**, ok 31.440s
- **未跑**：`.scratch/wisp/probes/154/gate-clauses.sh`（派单 §0 明令：别家在树里，数不可归因）。它的**文本**引用：
  腿名册唯一处 `:69 LEGS_EXPECT='G1 G1b G2 G3 G4 G5 G5pos G5neg G6 G6pos G6neg G7 G7pos G7neg'`；
  `:347-350` G1（生产码构造 `ToolRequest`，`want quiet`/`want_n 0`，排 `*_test.go` 与 `internal/agent/`）；
  `:352-355` G1b；`:357-361` G2（`want ring`/`want_n 2`，`:345` 记录锚 `f1b99a70` 现量 G1=0/G1b=0/G2=2/G3=0/G4=0）；
  `:363-366` **G3（`want quiet`/`want_n 0`，pattern `^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID`）**
  ⇒ 与本程结论同向：**别给 `Loop` 加收带 taskID 的导出方法**（甲形不需要）；
  `:368-371` G4。**这些是文件内容，不是本程读数。**
- **没测的**：见 `probes/177/c1/readings.md` §4 七条（其中第 1 条＝"两侧红名单是推导、不是现跑"，本程最大缺口）。

---

## 8. 交件报告剩下的格子

**被拒／没成功的调用**：无（本程 34 次工具调用全部成功，未触发权限拒绝；无取数前失败）。
**有没有跑过删除命令**：没有。没跑 `rm`/`git clean`/`restore`/`checkout .`；没建 worktree；没 push；`git add -A`／`git add .`／`commit -a` 一律未用。
**凭据值抄录**：零枚。

**伪授权两栏（各带出处）**

| 栏 | 内容 | 出处 | 本程态度 |
|---|---|---|---|
| 甲：转述在场、**我未复核** | "owner 刚批了甲（精确豁免）"——派单 `:1`/`:5` 与票面 `:28` 都这么写，MEMORY 也记"09-27 21:3x 四枚已批" | `.scratch/wisp/dispatches/2026-09-27-214x-readonly-177-c1…md:1,5`；票面 `:28` | 我**没读** `docs/reports/pending-and-issues.md`（不在写面）⇒ 台账那枚 `A##` 在不在盘上我**不声称已核**。我只把这"批准"当作**形状边界**（目录／前缀／配置项级一律不做），不拿它授权任何产码 |
| 乙：在场而互相矛盾 | 票面 `:3` 写"写手腿**按住**等 owner 对 `Q-61` 的一句话"、`:33` AC#0 写"丙形今天没有任何读数…**不许被当成已裁**"、`:34` AC#1 写"本票至今**一枚候选都没被独立量过**"；派单却写"owner 已批甲" | 票面 `:3`/`:33`/`:34`（与 `:35` 那段自纠）对比派单 `:1` | 我按**更保守的那一面**做：只裁形状、不动任何冻结件、不勾任何框、并把 AC#0（丙腿零读数）与 AC#1（要求现跑）两格**明写成本票未裁完** |

**next=（含一句：177-r1 写手腿派之前还缺什么）**

写手腿**现在还不能派**，缺这七件，按轻重排：

- a. **缺现跑的红名单**：AC#1（票面 `:34`）逐字要"跑台件"，本程只能〔读码推导〕。要在 scratch 合成树上对"甲形（mark 局部排除）"跑一次变异台件，
  实测 §4.1 六行里哪几枚真红、`provenance_test.go:136`/`:139` 是否仍绿——**这两枚若红，就地停手回禀**，不许改判据。
- b. **缺 P 送达桥的载具裁定**（这是唯一可能踩到人工批准的口子）：`internal/tools/tool.go:39-60` 的 `Result` 今天没有承载"宿主刚写下的路径"的字段
  （`Origin` `:50-53` 语义是"内容从哪来"，不是"我把它抄到了哪"）。两条路：① 给 `Result` 加收字段＝**C1 契约面变更＝人工批准**；
  ② 由 `bridge.mark` 用 `dec.TaskID` 去 `TaskRoster` 查（桥在 `internal/tools` 内，够得着）＝不加 C1 面，但要求 `Record` 先在场（见 c）。**要 owner/编排者二选一，别让写手腿自选。**
- c. **缺 176 的起跑口**：`TaskRoster.Record` 生产码零调用点（现量），`cmd/wisp/run.go:361` 只建不填 ⇒
  CLI `wisp run` 端到端腿（派单 §1.3 三面之一）**在本票落地那天仍然跑不出来**。票面 `:40` AC#6 要 176 不早于本票合入，
  这条要写进排程，别在 AC 里假装端到端已通。
- d. **缺生命周期那一问的答案**：排除跨度必须证明"随 scope 一起掉"（`provenance.go:425` `Scope.Close`、`provenance_test.go:305` `TestDisposalScopeClearsTaints`），
  否则跨任务存活＝一枚新的洗戳面。派给写手腿时把这一格当 AC 写，别留给"以后加固"。
- e. **缺契约文字那一枚批准**：`provenance.go:468-473` 那句"Returns false only when the content has no matchable characters after normalization"、
  `taintmatch.go:11-15` 那句"per-token taint tracking is REJECTED (16.9#1)"、`SPEC-06.md:70`/`:37` —— 三处文字要在产码**之前**摆给 owner，属"改契约＝人工批准"。
- f. **缺行号重锚**：票面 `:14` 的 `task.go:229` 现读是 `totalBytes := len(full)`，那句 `…全文见 %s…` 在 `:245-246`（另一枚同款在 `spill.go:141`）。
  派单文案若继续引 `:229`，写手会改错行。
- g. **缺 AC#0**：丙形（每次续读一张 L2 卡）今天仍**零读数**（票面 `:33`），本票没裁它；派单 `:6` 也明令我别为它补台件。
  若 owner 要的是"三支都量过"，AC#0 是一枚独立的只读腿，别让本表冒充它。
