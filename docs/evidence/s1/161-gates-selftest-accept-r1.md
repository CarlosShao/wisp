# 161-v1 非实现者验收表：票 161 的 AC#1／AC#2／AC#3 三格（对抗验收，逐格独立读数）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-092x-accept-161-v1-selftests-ac1-ac2-ac3.md`（09-27 09:2x）
- 本程性质＝**非实现者裁决**。三格的实现者分别是 161-r1（AC#1）／161-r3 写码＋编排者代提复跑（AC#2）／161-r2（AC#3）。
- 锚＝**`8efece378a38e74bd95f6cf7aabd971fb3dfc3eb`**（`git rev-parse HEAD` 现量，`git cat-file -t` = `commit`，2026-09-27 09:20）。
- ⚠ 并发事实：另程此刻在写 `probes/154|161/**`；本程**全部被验读数取自仓外副本**——
  `git ls-files` 全清单（2385 枚）`tar -T` 到 `/tmp/wisp161v1/tree`，再把 22 枚工作树相对锚有出入的已跟踪路径
  逐一 `git show <锚>:<path>` 复写回锚字节，副本内 `git init`＋全量 add/commit 建索引。
  副本内 `git ls-files '*.go'` = **531**（与 r5 的甲形分母同数）。锚到 HEAD 之间 `tools/d22scan`、`.github/workflows/ci.yml`、
  `docs/evidence/s1`、`internal`、`cmd` 的 `git diff --stat` 为**空**。
- 本机仪器：`go1.27.1 windows/amd64`、`gofumpt v0.12.0 (go1.27.1)`（两枚都现跑 `--version`，日志 `probes/161/v1/logs/go-version.txt`）。
- 档位约定：本件默认〔我本轮现跑过〕；凡引用实现者的话都显式降档。

---

## 0. 本程**没**测什么／哪些读数受并发影响（先写这节）

1. **乙形（工作树归因）没在真树复跑**：真树 `probes/**` 正被另一程写，任何"此刻跑 `gofumpt -l .`"的读数都不属于任何锚。
   本程只在仓外锚副本上跑了**甲形**（§2.4）。⇒ 编排者 09:0x 那句"乙形 9 行全归票 161、rc=0"的档位在本件里是
   〔仅自述，我未复算〕——要钉它需要一次安静窗口的真树现跑。
2. **CI runner 一次没跑**：GitHub 上任何一步的真实 conclusion 没量过；§2.5 的可达性是**按 ci.yml 锚字节的作业/步骤结构**现读的。
3. **B3 盲区（bans #1–#5 不看 `_test.go`）没独立复算**：派单只要求 ≥2 枚盲区复算，我做了 B1/B4/B5/B6 四枚＋B2 的扩张取样（§1.3），
   B3 属"未定义即停"族，留给 owner 的话。
4. **r2 的尺 B/C/D 本体（AST 尺、链接器尺）没重跑**：AC#3 我换的是**自己的 grep 尺**（独立于它的四把），它的尺只抽验了 file:line 对得上。
5. 门铃的**运行期端到端**（真起进程看卡）没测——那本来就是 r2 §2 N1 自报未测，本格判据不含它。

---

## 1. AC#1 裁决：〔成立〕——"8 枚全点得到＋六枚盲区登记"经独立复算成立；本件另登记两枚 r1 未取样的相邻形状

### 1.1 名册现读（锚字节，不抄话）

`tools/d22scan/main.go:5–49` 的 `// Bans` 块：编号禁令 **8 枚**（bare-goroutine / pathresolver-bypass / plaintext-key /
wallclock-timeout / mirror-hash / panel-approval / internal-artifact-tool / emoji），`unparseable` 是发现类型不算分母。
⇒ 与 r1 §0.5、票面 09-26 21:3x 更正一致。**"九条"作废**这一条成立。

### 1.2 我自己的新样本（27 发之外、r1 名册没出现过的形状）

台件＝`.scratch/wisp/probes/161/v1/{samples,logs}/**`；跑法＝副本内 `go build` 出的 `d22scan-v1.exe -root <假根>`，
四枚假根在 `/tmp/wisp161v1/fakeroots/{rootA..rootD}`（不进生产树）。**假根下限读数先撞了一发好门**：
rootA 骨架（3 枚 .go）被 main 的"production .go < 10 ⇒ rc=2 拒答"守卫挡下 ⇒ 我的"干净基线"对照由
**同跑同 tag 的必须响样本**承担（见每格控制列），这一发顺带证明 -root 侧也不存在"递给尺 3 枚文件报 clean"的恒真形状。

| 发 | 形状（新造的） | 预期 | 实测（rootB rc=1，finding 逐条在 `logs/v1-rootB.txt`） | 控制（同跑同 tag） |
|---|---|---|---|---|
| s1 | 方法值裸协程 `go s.worker()` | 响 | **响**（`bare \`go s.worker(...)\` is banned ... R16: named calls count too`） | — |
| s6 | 裸 `filepath.Abs/Clean` | 响 | **响**（两臂各一条） | — |
| s7 | `const probeAPIKey = "PROBE161V1..."` | 响 | **响** | — |
| s8 | 同行 mirror+sha256 | 响 | **响**（2 条） | — |
| s9 | `now.Sub(time.Now())` | 响 | **响** | — |
| s5 | 三行三字形：`→`／`①`／`✓` 各占一行 | 只 `✓` 响 | **只 `s5.go:7`（✓ 行）响**，3/5 行静默 | ✓ 行响 |
| s2 | `import fp "path/filepath"` + `fp.Clean` | （盲区）不响 | **不响** | s6 同跑响 |
| s3 | `time.Now().Sub(deadline)`（反向操作数） | 我认为它漏 | **不响** | s9 同跑响 |
| s4 | map 字面量 `"apiSecret": "PROBE161V1..."` ＋ `cache["apiKey"] = "..."` 索引赋值 | 我认为它漏 | **不响** | s7 同跑响 |
| unp | rootD：故意不解析的 .go 内含 `go func(){` | 只响 unparseable | **该文件只有 `[unparseable]`**，全树 `bare-goroutine` 计数=1（是 s1.go 的，不是它的） | s1 同跑响 |
| al | rootC：与 rootB 同文件集＋allowlist 条目 `pathresolver-bypass→internal/` | 被压掉且零打印 | **pathresolver-bypass 0 条**、bare-goroutine 照响、全输出 `grep -ci suppress`=**0** | rootB 同文件响过 |

**r1 表格的复算结论**：八枚编号禁令里我在假根上亲手打响 6 枚（#1 #2 #3 #4 #5 #8）；#6/#7 两枚经 §2 的
`-self-test`（AC#2 表，本程现跑 rc=0，#6 两响一静、#7 两响一静）在我本轮过了一遍 ⇒ **"8 枚今天都点得到"经独立路径成立**。
登记成盲区的六枚里我独立复算 **B1（s2）／B4（s5）／B5（rootD）／B6（rootC）四枚，全部与其读数一致**；
B2 被我的 s4 扩张（map 字面量与 IndexExpr 左值赋值同族两形今天都不响——r1 §2 N3 自己标了"没取样"，这一格从"没量"变成"量过、静默"）。

### 1.3 本件新增登记（r1 表格里没有的两枚，不算退回、算续表）

- **V1-N1** ban #4 反向操作数：`time.Now().Sub(deadline)` 不响，而禁令文字是"timeout/deadline 逻辑里的墙钟差"。
  会响条件（若哪天扩它）：上表 s3 翻响。今天它是**登记外**的相邻形状，r1 的"射程比文字窄"结论因它**更强**而不是被推翻。
- **V1-N2** ban #3 的 map 索引赋值：`cache["apiKey"] = "<24 字符>"` 不响（左值是 IndexExpr，AST 两支都不看）。
  同族同因（只认 ValueSpec 与 Ident 左值），会响条件同 B2 那条。
- 顺手一枚反向哨兵（不是洞）：s1 的 `go s.worker()` 证明 r1 §2 N3 列的"没量的相邻形状"之一今天**在射程内**。

**恒真/空心自查**：表内每一枚"不响"都有同跑同 tag 的"响"作控制（见控制列）；rootA 的骨架下限由 main 的 ≥10 文件守卫硬拒。
⇒ 本程没有一枚"永远不响的检"，也没有一枚"递给尺 0 枚文件报干净"。

---

## 2. AC#2 裁决：〔成立〕——默认立场"它可能是假的"没有被证实：全部读数由本程在仓外锚副本上亲跑

被验面＝`tools/d22scan/{selftest.go,selftestsamples.go,selftest_test.go,main.go}` ＋ `.github/workflows/ci.yml` 的 lint 新步。
**本程一字节未改**（三发单点回退全部打在 `/tmp/wisp161v1/tree` 的副本上，跑完 `cmp` 复原＝RESTORED-IDENTICAL）。

### 2.1 基线（档位＝〔我本轮现跑过〕）

`cd tools/d22scan && go run . -self-test`（副本内）⇒ **rc=0**，尾行
`clean - all 34 direction checks passed (19 expect-ring, 15 expect-silent)`；
名册行 `roster read from main.go = 8 numbered ban(s) [...] + 1 finding type(s); 34 cases, 9 tag(s) covered`。
⇒ 枚数与本程 §1.1 的 `// Bans` 现读一致；19/15 为工具自打的 `countWant`，本程按样本表逐 tag 手点复核同为 ring=19/silent=15。
**tag 级承重名册**（ring 枚数）：bare-goroutine **4**／plaintext-key 2／wallclock 2／panel-approval 2／artifact 2／emoji 4／
**pathresolver-bypass 1**／**mirror-hash 1**／**unparseable 1**——ring=1 的其实有**三枚 tag**，不止 mirror-hash 一族（见 §2.3 与 next=）。
日志 `probes/161/v1/logs/selftest-baseline.txt`。

### 2.2 `go test` 那一腿

`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（副本内）⇒ **`PASS=34 FAIL=0 SKIP=0`、`=== RUN=76`、rc=0**
（`logs/runtests.txt`）。⚠ 一枚如实登记：`selftest_test.go:178` 有一发**条件性** `t.Skipf`（仅当 PATH 上没有 go 工具链时跳过
`TestSelfTestFlagIsWiredInTheBuiltBinary`）——它不豁免任何断言的语义，且 runtests 的零-SKIP 规则实测 SKIP=0。本程**没动它**（派单禁）。

### 2.3 单点回退进攻（三发，全在副本）

| 发 | 摘什么 | 预期 | 实测 | 日志 |
|---|---|---|---|---|
| ① | bare-goroutine 四枚 ring 里摘 `leaknamed.go` 一枚 | **仍绿**（tag 级承重，摘一枚符合设计） | **rc=0**：`clean - all 33 direction checks passed (18 expect-ring, 15 expect-silent)` | `logs/mut1-leaknamed-gone.txt` |
| ② | `mirror-hash` **唯一** ring 样本摘掉 | **必须红＋零行用例** | `go run` 退 **1**（二进制自身退 **2**，输出末行 `exit status 2`）：`FATAL the table does not cover the tool (1 hole(s); an unrun self-test is not a green self-test` ＋ `HOLE tag "mirror-hash" has only an expect-silent sample ...`，**一条 case 行都没打** | `logs/mut2-mirror-ring-gone.txt` |
| ③ | wallclock **唯一** silent 样本摘掉（恒真方向的对称回退） | 必须红 | 同形：`rc(got run)=1`／二进制 2，`HOLE tag "wallclock-timeout" has only an expect-ring sample - ... tautological-check shape this repo has rejected twice` | `logs/mut3-wallclock-silent-gone.txt` |

⇒ "摘任意一味是否有一发变异从此打不红"：对**每枚 tag 的末枚方向样本**，②③两形实测会红且**在跑任何用例之前**红（名册空心 ⇒ 拒出判语），
这一条今天由本程**实测**而非引用。摘非末枚（①）不红是设计而非缺陷，本件按操作定义把它和②③分开记。

### 2.4 空心保护（递给尺 0 枚文件必须硬退出）

- `-self-test-roster` 指向一份无 `// Bans` 头的文件 ⇒ 二进制 **rc=2**＋`FATAL ...: no `// Bans (` header ...`（`logs/hollow-roster.txt`）。
- CI 甲形那步（`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`）：锚副本上 ⇒ **rc=0、`(A) lines=0 files=0`、分母 531**
  （`logs/attrib-tracked.txt`，gofumpt `v0.12.0 (go1.27.1)` 由脚本自打）；
  再在一枚 `git init` 后**零枚 .go** 的空仓副本跑同一条 ⇒ **rc=2**＋
  `(A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing`（`logs/hollow-emptyrepo.txt`）。
  ⇒ 这一发同时钉住甲形步不是第二枚恒绿。
- `attrib.sh --self-test`（判据分类器双向自测）⇒ **cases=8 failures=0 rc=0**：5 枚必须响（含 `TRACKED<CI-RED>` 与两枚假票号）＋
  3 枚必须不响（含"另一枚真票 169 的路径"）（`logs/attrib-selftest.txt`）。

### 2.5 CI 可达性（那一问的答案，按 ci.yml 锚字节现量）

- 那步在 `lint` 作业**第 4 个 step / 第 2 个 run 步**（ci.yml:83–106：`run: go run . -self-test`、`working-directory: tools/d22scan`、
  解析出的键只有 `name/run/working-directory` ⇒ **无 `if:`、无 `continue-on-error`**）；它前面只有 checkout、setup-go、
  和"D22 scanner positive control"（`runtests.sh -C tools/d22scan`）三步。
- **`internal/panel/tokens_fourway_test.go` 的常红会不会吃掉这一步？**：不会——本程两向现量：
  ① 它（`TestC21DesignTokensFourWayAgree`）在锚副本上实测 **rc=1 常红坐实**（本程只判断、未修、未 Skip，尾部 40 行在
  `logs/fourway-tail40.txt`）；② 它的归口是 `scripts/portable-tests.sh`（core 清单第 141 行含 `internal/panel`）
  ＝**test-core / test-windows 作业**，而整份 ci.yml `grep -c 'needs:'` = **0** ⇒ 作业间无依赖链，GitHub 的
  "首个失败之后不跑"只在**同一作业内**生效。lint 作业没有任何一步编译或运行 `internal/panel`。
- lint 作业内**真正**能吃掉这一步的前置风险（登记给编排者，不是退回项）：第 3 步 runtests 含
  `TestScannerSelfScanOfRealRepoIsGreen`／`TestRealRepoLedgerIsHonest` 两枚真仓扫描用例——票 169 那一族的前端字形若再红，
  第 3 步先响、第 4 步的 self-test 读数就失踪。锚上实测两枚全绿（34/0/0），今天不构成缺口；这正是 AC#6"聚合退码/步骤归因"那一族。

### 2.6 对"无实现者自证"这一档的处置

`docs/evidence/s1/161-selftests-r3.md` 现量 **262 行、止于 §3**，AC#2 一节确实不存在 ⇒ 编排者 23:5x 那格的档位声明
〔编排者代提并复跑，非实现者自证〕**与盘相符**。本表 §2 以非实现者现跑补上验收腿；**代码本体经四路进攻（2.1–2.4）无一谎**。
⇒ AC#2 **成立**。实现者自证缺口本体不在验收权限内补，登记在 next=。

---

## 3. AC#3 裁决：〔成立〕——换尺重数，三个头条读数全部复算一致；零结论都先过了我这把尺的正控

被验物＝`docs/evidence/s1/161-doorbell-census-r1.md`（161-r2）。它的尺是 AST/链接器族；**我的尺是 grep 族**
（语句口径照它的分层定义自己取，不引它的中间件）。跑在锚副本（`/tmp/wisp161v1/tree`）上，全表复现命令与输出
已归档 `probes/161/v1/logs/ac3-rulers.txt`。

### 3.1 两把尺的枚数对比

| 格 | r2（AST 尺 B/C＋grep 上界＋nm） | v1（本程 grep 尺，锚副本现跑） | 一致？ |
|---|---|---|---|
| C 层「问题真到人眼前」生产发问点 | **2**：`approval/gate.go:260`、`:486` | **2**：同两枚 file:line 逐字相同 | ✅ |
| B 层路由点（`PendingWindow/PendingApproval` 非 test 调用） | **4** | **4**（bridge.go:380/:395、gate.go:236、run.go:479） | ✅ |
| D 层门铃 `Dispatch(EvApprovalNeeded,…)` 生产 | **0**；非 test 提及 3 枚＝1 声明＋2 表数据 | **0**（`Dispatch(.*EvApprovalNeeded` 非 test 空，rc=1）；提及同样 3 枚同一逐字 | ✅ |
| 答案侧 7 符号生产调用点 | **各 0**（它的 grep 上界 9 行全是声明/注释/fallback） | **各 0**：我的模式只命中 2 行，都是 `doc.go:32/33` 注释 | ✅ |

### 3.2 "0 命中"先拿已知正控打我的尺（派单 §1 AC#3 硬要求）

台件＝`/tmp/wisp161v1/ctlroot/internal/x/ctl.go`（一发 `ui.Prompt`、一发 `Dispatch(EvApprovalNeeded, nil)`、一发 `q.Native()`）。
同三条命令打上去 ⇒ **1 / 1 / 1 全非零**（归档件末段）。**归档件里留了我自己的一发尺错**：控制循环第一次用 `-E` 跑
`Dispatch(.*` 这个 BRE 形状模式 ⇒ `Unmatched (` 报错、那一行记成 0；改 `-F` 复跑 ⇒ 1。错行与复串行**都留在**
`ac3-rulers.txt` 里不抹（派单 §5 那一族仪器坑的自家再现：管命令形状的那一发，写错的那次不算读数）。
⇒ 上面三个 0 是"尺会响处无量"，不是"尺根本不会响"。
r2 自己那三枚正控台件（含"解析失败不许静默归零"的教训）档位＝〔它的读数，我只抽验了 file:line 归属，AST 尺本体未重跑，见 §0.4〕。

### 3.3 承重那句与"防忘记/防回归"

- **摘掉这把尺，外部可见读数会变吗？不会**——普查是证据件，没有任何 CI 腿或仪器消费它的数。它买的形状与 r2 自述一致：
  D 层与答案侧＝**防忘记**（树里今天零枚仪器盯 `Dispatch(EvApprovalNeeded…)` 枚数，今天接上或再拆掉都不会红）。
- 唯一的**半枚防回归**（C 层）我按原文复算成立：`cmd/wisp/panel_pump_test.go` 211–214 行实测有
  `if len(cardRecords) < 2 { t.Fatalf("depth=1 records ... want the run's own publish plus this case's sample" }`，
  且注释明写第一条来自 `consoleApprovalUI.Prompt` 自己入账 ⇒ "两枚 ui.Prompt 全部消失"这枚事故今天会红。
- **可复算的会响条件**：r2 §12 表把四个从 0 变 ≥1 的复算命令写成了可跑的 grep；本程以同形 grep 在锚上打出 2/4/0/0，
  正控打 1/1/1 ⇒ 会响条件**已被外部复算过一次（我）**，满足票面 AC#3"零枚答案必须带会响条件"的入账要求。

### 3.4 一条纪律登记（不改判、只归口）

r2 §10 自陈一枚 `rm -rf`（范围＝它自己本轮的两份副本目录）。档位＝〔仅自述〕——被删字节已不在任何锚上，本程**不可复核亦不复跑**。
它没有掩盖而是主动报并留了全文，处置（记台账与否）归编排者/owner。

---

## 4. 对编排者的不服（≥2 枚，直说）

1. **你派单 §1 AC#2① 的"三枚里的一枚"是过期数**：盘上现表 bare-goroutine 有 **4 枚** ring 样本（基线 19 expect-ring 逐 tag 手点）。
   结论方向没翻（摘一枚仍绿是 tag 级设计），但下一位若按"三枚"去核名册会对不上并怀疑表被改过。出处＝我 §2.1 的基线与 mut1 两本日志。
2. **"rc 0→2"与 CI 步骤形状混写**：派单②要求我贴"我自己的命令与输出"，而我实测两形并存——二进制退 **2**，
   `go run . -self-test`（CI 那一步的形状）退 **1**。你对 r4 的转述只给了 2。两形都红、判据不翻，但表格文字应把
   "binary rc"与"step rc"分开钉，否则拿 GitHub step conclusion=1 去核"0→2"的人会以为读数漂移。
3. **你的"它可能是假的"预设没有兑现，这一点我按读数认输**：AC#2 的码是真的、四路进攻全部按你说的形状动
   （包括"空心⇒拒出判语"今天在我手里第二次被实测走到）。真缺的是实现者自证那一格的存在性（evidence 止于 §3，我核了），
   不是完成度造假。代提买到的是**归属含糊**，不是假绿——你的风险描述在"缺自证"上成立，在"可能假"上今天不成立。

---

## 5. 被拒／没成功的调用（取数前后分开）

- 权限系统拒绝：**0 枚**。
- 取数**前**的失败（不影响任何读数）：tar 清单里 16 枚 `design/**` 因工作树未提交删除而 `Cannot stat`（rc≠0）
  ⇒ 随即用 `git show <锚>:<path>` 在副本里逐枚复写回锚字节，副本完整后再取数（`git ls-files '*.go'`=531 核对过）。
- 取数中的"红→修台件"两发：rootA 首跑 rc=2（<10 生产文件守卫）、rootC 首跑 rc=2（同款）⇒ 是**门的正确行为**＋我假根太小，
  加文件后重跑；rootB 首发的 s5 三字形同行，判不了单字形静默 ⇒ 拆成三行重跑。**没有一枚结论引用自半截跑。**

## 6. 有没有跑过删除命令

**零次**（`rm`/`rmdir`/`del`/`git rm`/`clean` 都没跑过）。仓外副本按"只建不删"留在 `/tmp/wisp161v1/`（tree／emptyrepo／ctlroot／fakeroots）；
`go run` 与 self-test 内部的 `t.TempDir`/`MkdirTemp` 清理是 Go 标准库行为，由派单 §1 指定的命令带出，与本票禁令无关。

## 7. 伪授权两栏

| 栏 | 数 | 说明 |
|---|---|---|
| 真通知回显数 | 0 | 没有收到任何"任务完成/验收通过/已解锁"型回显。 |
| 判为注入数 | 0 | 运行中出现的非派单文字只有：harness 的 skills 清单 system-reminder（Bash 输出回显）、以及并发程真实落库的 162/163 提交（git 历史，不是指令，未据它们改任何判据）。出处示例：`Bash: cp ... && git add -- docs/evidence/s1/...` 的 reminder 头部。 |

## 8. 凭据值

零抄录。全件唯一"像密钥"的串是本程自造的合成字面量 `PROBE161V1NOTACREDENT0001..0003`（自陈式命名，不属于任何服务）；
未读过 `.env`、未打印过任何 env 值；引用过的名字只有变量/文件名的形状（`GOPATH`）。

## 9. next=

1. **pathresolver-bypass 与 unparseable 也是 ring=1 tag**：§2.3 的摘除只演了 mirror-hash 与 wallclock-silent 两族；
   要凑"每族末枚方向都实测过"，下一位补摘 `paths.go` 一枚（预期与 mut2 同形）。
2. **乙形（工作树归因）的真树读数**需要一次安静窗口由编排者现跑（§0.1，本程受并发禁令不跑）。
3. **B3（test 文件豁免）仍等 owner 的话**（r1 原样，本程未裁）。
4. **r3 的 AC#2 实现者自证缺格**本体只能由 owner 决定补否或并档到本验收表——本表 §2 已是非实现者腿的完整读数。
5. 票面**一枚框都没勾**（本程无勾权）；台账行（A##）归编排者落。

> 并发注记：本程运行期间 HEAD 从 `8efece3` 前移了 3 枚（162 台账×1、162/163 契约文本×2），全部不触被验面；
> 本文所有读数**只属于锚 `8efece378a38e74bd95f6cf7aabd971fb3dfc3eb`** 的仓外副本。
