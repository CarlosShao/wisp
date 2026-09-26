# 154 对抗验收程（非实现者）— 「关」只长在环路那一枚 id 上：RESERVED＋触发门那一票的六格

- 派单＝`.scratch/wisp/dispatches/2026-09-26-132x-accept-154-r1.md`（61 行，本程 13:2x 逐字读完）
- 工单＝`.scratch/wisp/issues/154-the-close-only-covers-the-loop-task-id-so-host-supplied-task-ids-still-have-no-owner-and-concurrent-tasks-have-zero-readings-reserved-with-a-trigger-gate.md`（88 行，`git ls-tree --name-only 5766be9 .scratch/wisp/issues/ | grep -c "^\.scratch/wisp/issues/154-"` ⇒ **1**，本程不改它一字、不加 `-done`）
- 被验交付物＝`docs/evidence/s1/154-host-id-never-closed-r1.md`（锚点 **`5766be9`**，**695 行**）＋ `.scratch/wisp/probes/154/**`（锚点上 **34 枚**，现读 `git ls-tree -r --name-only 5766be9 .scratch/wisp/probes/154 | wc -l`＝34）
- 本程角色＝**验收程**：裁决者≠实现者（`AGENTS.md §0.3`、`SPEC-12 §4.3` #1/#3）。本程**不新造规矩**，六格各自要换的那把尺逐条来自派单的表。
- 本程写面（**只有两枚**）＝本文件 ＋ `.scratch/wisp/probes/154-accept/**`；临时件在 `D:/tmp/wisp-154-accept/**`，**只建不删**。

---

## 第 0 节 锚点、脏树分离、工具链

**被验锚点＝`5766be9`，全程不拿工作树当被验版本。** 现读三条：

```
$ git cat-file -t 5766be9 / 79551d1                      →  commit / commit
$ git rev-parse 5766be9:docs/evidence/s1/154-…-r1.md      →  330c18d463adb4163d38aed31fff58ba610952c8
$ git rev-parse HEAD:docs/evidence/s1/154-…-r1.md         →  330c18d463adb4163d38aed31fff58ba610952c8
  ⇒ 这枚证据件在锚点与 HEAD 是同一枚 blob（编排者另已 cmp 过前 520 行），本程仍一律 `git show 5766be9:<路径>` 取版。
```

**共享树确实在别人手里往前走**（派单点名的并发，本程自己量，不抄派单的数）：

```
$ git status --porcelain -- internal cmd                 →  空（13:2x）
$ git diff --name-status 5766be9 -- . ':!docs' ':!design' ':!.scratch'
	A	cmd/wisp/slo_exit_os_156_windows_test.go
	M	cmd/wisp/slo_report_144_windows_test.go
	M	cmd/wisp/slo_windows.go
$ git diff --name-status 5766be9 -- internal             →  空      ← 关键：internal/** 锚点＝工作树
$ git ls-files --others --exclude-standard -- internal cmd →  空（无未跟踪 .go 混进编译）
```

⇒ 三条后果，都写死：①**`cmd/wisp` 的工作树不是 `5766be9`**（156 的 WIP 已由编排者代提在 `0e95353`，本程进场时 HEAD 又走到 `7a56af5`）⇒ 本程复跑 `cmd/wisp` 若不做取版，读到的就是别人的半成品，**那种跑只能记"没判据"**；
②`internal/**` 在锚点与工作树逐字节相同 ⇒ 本程 AC#5 那一发包选在 `internal/tools`，**原地跑＝跑锚点树**（这条是本程能用便宜尺的原因，先钉住再跑）；
③`design/**`／`frontend/**` 的未提交删除与前端会话的活**不进本程任何零命中宣称**（票面 AC#4 的 ⚠ 明写）。

工具链现读：`go version` ⇒ `go1.27.1 windows/amd64`；`gofumpt` 本程没用到（不写码）；`third_party/sherpa-onnx` 的 dll 是**未跟踪**本机件（`git ls-files third_party | wc -l`＝0 ⇒ 任何"干净树"复跑都拿不到它，这点与 §0 第 0 节 r1 的读数一致）。

**引号先验号**（本程全程执行）：`245d5f4` 与 `782b01c` 两枚**按"不存在"处理**——`git cat-file -t` 各 rc=128（本程 13:2x 现量，与派单自述一致），本程不据它们判任何事、也不替编排者圆。
六枚 154 commit（`1424aa7`·`352d3d8`·`76c89b2`·`304199f`·`e9a1db0`·`79551d1`）＋ `5766be9`＝**七枚全部 `git cat-file -t`＝commit**（现读）。

---

## 第 1 节 格 AC#1 — 触发门：复跑他们的尺 ＋ 换一把尺判"判不出"那两枚

### 1.1 复跑他们的 `gate-clauses.sh`（派单点名那一发）

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh 5766be9   →  脚本 rc=0
  # rc= 序列（G1 / G1b / G2 / G3 / G4）＝ 1 / 1 / 0 / 1 / 1      ← 与 r1(1424aa7)、r2(be603fd) 逐枚相同
$ G2 那两行＝5766be9:cmd/wisp/run.go:546（注释）＋ :563（唯一生产调用）   ← 与 §2.3 表逐字相同
$ 附段：摘掉排测试 ⇒ 'ToolRequest\{' 命中 14；.Execute( 全名册＝3 行（loop.go:730/:734、bridge.go:463）
```

台件＝`.scratch/wisp/probes/154-accept/ac1-gate-clauses-rerun-by-accept.txt`。
**名册等价性不是"看着一样"**：本程把三枚读数各自剥掉锚点前缀与时刻行后 `cmp`——

```
$ cmp <(归一化 本程@5766be9) <(归一化 r1@1424aa7)      →  rc=0
$ cmp <(归一化 本程@5766be9) <(归一化 r2@be603fd)      →  rc=0
（归一化尺＝把行首的 sha 一律换成 ANCHOR:、删掉"# 生成时刻/# 票 154 AC#1 触发门"标题行与"\$ git grep"回显行；
 台件＝probes/154-accept/norm-{accept,r1,r2}-roster.txt）
```

⇒ **他们这一格给的名册，在三个锚点上是一枚东西**；门"今天不响"是可复算的，不是那一刻的运气。

### 1.2 他们的正控可不可复算（派单"不许靠做不到免检"的第一半）

`/d/tmp/wisp154-gatectl` 那棵合成树**还在**（`git log` 三枚 `ctl`/`ctl2`/`ctl3`；`git ls-files`＝`cmd/wisp/host_probe.go`·`cmd/wisp/close_owner_probe.go`·`internal/agent/loop_with_taskid_probe.go`）⇒ 本程原样复跑他们的 `gate-positive-control.sh`：G1/G2/G3 三枚在合成树上**各自真响**，本仓侧仍空。台件＝`probes/154-accept/ac1-their-positive-control-rerun.txt`。
⚠ 一处**复算面**的真相：那枚脚本读的是 `ANCHOR=$(git rev-parse HEAD)`（不是参数化的锚点）⇒ 本程这一发 HEAD＝`7a56af5`，**不是** `5766be9`。本程核过 G1/G2/G3 三条尺在 `internal`＋`cmd` 上与锚点无差（§0 的 `git diff 5766be9 -- internal` 空、`cmd` 只有 156 的三枚文件且都不含这些符号）⇒ 这一发的等价性成立，但**下一位复跑时锚点会再漂**：正控可复算的前提是那棵合成树没被清掉，而它在 `D:/tmp`（仓外、不入库）。这是"门本体"的一处真实脆弱，登记给编排者（不是本格的判语）。

### 1.3 换尺判"判得出／判不出"（本程自己的尺＝**编译**，不是 grep）

§2.2 第一行的承重句是：**"要带着身份过桥，必须在某个文件里出现 `ToolRequest{` 这个字面量。它是唯一的语法位置"**。本程不数命中，**把这句话当定理打**：在 `D:/tmp/wisp-154-accept/g1blind` 造一棵五形合成树，逐形 `go build ./...`（**BUILD OK**，Go 1.27.1），再把他们的两条 pattern 原样打上去。

| 形 | 写法 | 编得过？ | G1 `'ToolRequest\{'` | G1b `'\.Execute\('` |
|---|---|---|---|---|
| A | `var req agent.ToolRequest; req.TaskID="host…"; b.Execute(ctx,req)` | **是** | **不响**（零值＋字段赋值，无字面量） | 响 |
| B | `type dispatch = agent.ToolRequest` ⇒ `dispatch{TaskID: …}` | **是** | **不响**（别名不是那个词） | 响 |
| C | `req := agent.NewRequest()`（字面量在别的包）＋改字段 | **是** | 响的只是 helper 那个文件，宿主函数**不响** | 响 |
| D | `req2 := seed; req2.TaskID="host…"` | **是** | **不响** | 响 |
| E | `fn := (bridgeImpl{}).Execute` ⇒ `fn(req)` | **是** | **不响** | **不响**（`Execute` 后面没有 `(`） |

⇒ **三条判读，方向分开：**
1. **"构造点是唯一语法位置"这句在 Go 里不成立**（A/B/C/D 四形都带着宿主自造的 id 走到 `Execute`，而 G1 一字不吃）。本格不至于翻，因为 §2.3 明写"与 G1 任一枚响即同判"——**是 G1∪G1b 大致补上了 A–D**，不是 G1 单独。但**这句话一旦被下一位单独引用**（引文很容易只引 G1 那一行），"门空"就会被读成"没人能带 id 过桥"。⇒ 记为**缺陷 D-1**（措辞过强，且过强的方向＝漏报方向）。
2. **E 形（方法值／函数指针）两枚子句都不响**，而它是 Go 里很常见的接线形状（把 `Bridge.Execute` 存进 struct 字段再调）。本程在本仓现量这形今天**不存在**：`git grep -nE ':?= *[a-zA-Z0-9_.(){}]*\.Execute$' 5766be9 -- internal cmd ':!*_test.go'` ⇒ **rc=1**。⇒ 记为**缺陷 D-2（门的静默形，今天无实例）**：与 D-1 同源——门吃的是**词法**，票面 AC#1① 那条候选形状要判的是**语义**。
3. 本程**没有**为此新增任何断言，也不要求他们新增（票面 AC#6 明禁恒真检）。两形的正确处置是**门本体加一句"本门吃字面量与直接调用，不吃别名/零值赋值/方法值"**——那是文案，不是检。

台件＝`probes/154-accept/ac1-g1-blindspot-synthetic.txt`（五形源码在 `D:/tmp/wisp-154-accept/g1blind/**`，仓外、不删）。

### 1.4 他们承认"判不出"那两枚，是真的判不出吗（派单点名第二半）

| 他们判"判不出"的半条 | 本程现量 | 判定 |
|---|---|---|
| "`TaskID` 非空"＝值域问题，grep 给不出答案 | 本程没有能答它的仪器；但**有一条他们没写的近似**：把名册打到**赋值的右值形态**上——现读 `git grep -nE '\.TaskID\s*=' 5766be9 -- internal cmd ':!*_test.go'` ⇒ 生产里只有 `internal/agent/approval/gate.go:649` 一枚 `d.TaskID = taskID`（且它是 `Decision`、不是 `ToolRequest`）；`TaskID: ` 的字面构造点全仓生产只有 `internal/agent/loop.go:646` 一枚。⇒ 近似尺能给"名册"，不能给"非空"。 | **成立（真判不出）**，但**近似命令他们少给了一条**：票面 AC#1① 要的是"最接近的近似命令"，§2.2 给的是排目录/排值域两条**方向说明**，没给这条可跑的右值尺。记为**缺陷 D-3（该给而未给的近似命令，不是假门）** |
| "调用者不在 `Loop.run` 边界内"＝跨函数可达性，本仓无 callgraph 仪器 | `go.mod` 里**没有** `golang.org/x/tools` 的 require（`git show 5766be9:go.mod | grep -nE '^\s*(require\s+)?golang.org/x/tools'` ⇒ rc=1）；`go.sum:37-38` 有 v0.48.0 两行＝只在图上不在 require；**且它不在本机模块缓存**（`ls $GOMODCACHE/golang.org/x/tools@v0.48.0` ⇒ rc=2）⇒ 上这枚仪器**还要联网拉取**，不是本地改改就有。`tools/d22scan/main.go` 的 import 现读只有 `go/ast`·`go/parser`·`go/token`（单文件解析、无类型信息、无 callgraph）。 | **成立（真判不出）**。⚠ 但**他们给的理由里有一处不成立**：§2.2 写"新增依赖＋改 `tools/d22scan/**`，**两者都在本票零字节面里**"。现量：`go.mod` **不在**票面 AC#4 点名的 13 支里（`sed -n '39,41p' 工单 \| grep -o 'go\.mod'` ⇒ rc=1），新工具目录也**不等于** `tools/d22scan/**`（那一支只框 d22scan 自己）。⇒ 记为**缺陷 D-4（免检理由引错了面）**；注意它的**结论仍然对**——模块缓存里没有 ⇒ 这一程本来就不该联网上仪器。所以这是**理由的缺陷，不是判断的缺陷**。 |

### 1.5 本程新盘到的一枚门洞（本格最该被读的一条）

他们的门五枚子句**不覆盖 `risk.Provenance.OpenScope`/`CloseScope` 这一对**，而这一对恰恰是他们自己普查里**第二枚真"只开不合"**（§4.1 表第 2 行）。锚点现量（同一把 grep 尺、**判得出**，不需要数据流）：

```
$ git grep -nE 'OpenScope\(|CloseScope\(|prov\.Mark\(' 5766be9 -- internal cmd ':!*_test.go' ':!internal/risk/*'
5766be9:cmd/wisp/panel_assets.go:232:	prov.OpenScope(taintSourceScopeID)      ← 全仓没有任何 CloseScope 指向它（§4.1 第 2 行自陈）
5766be9:cmd/wisp/panel_assets.go:234:	prov.Mark(taintSourceScopeID, …)
5766be9:internal/tools/bridge.go:642:	b.prov.OpenScope(taskID)
5766be9:internal/tools/bridge.go:691:	b.prov.CloseScope(taskID)
（rc=0；台件 probes/154-accept/ac1-openscope-clause-not-written.txt）
```

⇒ 后果可判：**面板/球今天真正会用的那一形（自己拿 `Provenance` 开一枚 scope、不经桥）出现时，五枚子句全静默**——G1/G1b 只看过桥的派发，G2 只盯 `OpenTask|CloseTask` 这一对名字，G3 只看 `Loop` 的方法签名，G4 只看 L0 直通。
而 §6.4② 把"响的条件"明文**只从 G1–G4 取**，于是这形被排除在门本体之外。
⇒ 判读：这**不是**假门（五枚各自都响得住，§1.2/§1.3 有正控），是**射程比本票的危害面窄一格**；票面 §"为什么这不是再抠一次文案" 说的受害方（面板/球上线）今天已有一枚**在册的活样板**（`panel_assets.go:232`），普查量到了、门没接住。
⇒ 记为**缺陷 D-5（建议新增第六枚子句 `G5＝OpenScope/CloseScope 生产出现点名册（排 internal/risk 自身）`，今日基线 4 行、判据＝"开那侧枚数 > 关那侧枚数"）**。本程**不代他们写进门本体**（那是票面与注释面，本程写面只有两枚）。

### 1.6 本格判语

**AC#1＝成立（可判的门、承认判不出、无假门、正控可复跑），但带五枚必须登记的缺陷 D-1…D-5。**
其中 **D-5 是实质缺陷**（门的射程不含本票普查自己找到的第二枚实例），**D-1/D-2 是同一枚措辞‑射程缺陷的两面**（"唯一语法位置"过强＋方法值静默形），**D-4 是引错面的免检理由（结论不受影响）**，**D-3 是少给一条可跑的近似尺**。
本程**没有**逼实现程造恒真检，**也没有**代它造：本程新增的只有"编译正控"（证明形状的语法可能性）与"未写子句的读数"（证明门不覆盖），两者都不落进任何断言面。

⚠ 一条**不能说过头**的话（免得本程自己犯本格要抓的错）：以上读数判的是**锚点 `5766be9` 的门本体**；`bridge.go` 注释里那句"the gate … is one `git grep` away"指的正是这套子句，它没有承诺"所有形都能被抓"，所以 D-1/D-5 **不与注释文案矛盾**，只与证据件 §2.2 那句"唯一的语法位置"矛盾。

---

## 第 2 节 格 AC#2 — 文案落点：逐句指回盘上事实 ＋ 每个引用当场解析

派单要本程换的尺：**逐句判这三段是不是打发**，且**注释里不许预先引用尚未产出的读数**——"逐条 `ls`／`git cat-file` 它引的每个路径与号"。

### 2.1 改动形状（换尺＝读 hunk 头，不读他们的 grep 过滤）

r1 在 §3.3 用的尺是 `git diff -U0 | grep '^[+-]' | grep -v '^[+-]\s*//'`（文本过滤）。本程换成 **hunk 结构**（台件＝`probes/154-accept/ac2-bridge-comment.diff`）：

```
$ git diff -U0 --no-prefix 76c89b2^ 76c89b2 -- internal/tools/bridge.go
@@ -628,0 +629,4 @@    ← OpenTask 注释：纯插入 4 行，0 删除
@@ -650,0 +655,21 @@   ← CloseTask 注释：纯插入 21 行，0 删除
$ git diff --numstat 76c89b2^ 76c89b2 -- internal/tools/bridge.go   →  25  0
$ 非注释增删行枚数（grep -vE '^\+[[:space:]]*//'）                  →  0
```

⇒ 两枚 hunk 都是 `-,0`（**纯插入**）⇒ §3.1 那句"不抹原句、不改原句一字"**在结构上成立**，且插入点不落在任何原句中间（`-628,0` 落在原 OpenTask 注释末行之后、`-650,0` 落在原 8 行段末与 audit 段之间）——这条**只有 hunk 头能给**，他们的 grep 尺给不出"有没有把原句劈开"。
本程另核 `git log --oneline 76c89b2..5766be9 -- internal/tools/bridge.go` ⇒ **空**（这一枚文件在 154 之后没人动过）⇒ §3.2 贴的行号（`626-632`／`646-679`）在**本程锚点**仍逐位对得上：`func (b *Bridge) OpenTask` ＝ `:633`、`func (b *Bridge) CloseTask` ＝ `:680`（现读）。

### 2.2 注释里每一个引用，本程当场解析（两枚锚点各一次：写它的那枚 commit ＋ 被验锚点）

| 注释里的引用 | 解析命令 | 读数（`76c89b2` ／ `5766be9`） | 判定 |
|---|---|---|---|
| `internal/agent/loop.go:366`（原句引，仍在原位） | `git show <锚>:… \| sed -n '366p'` | 两枚锚点都是 `defer revokeAdmission()` | **对** |
| `Run/RunAsync -> newTaskID, internal/agent/loop.go:332/:321` | `sed -n '332p;321p'` | `:332 func (l *Loop) Run(` ／ `:321 func (l *Loop) RunAsync(` | **对，且顺序没写反**（票面 §现量第 4 条也是这个配对） |
| `newTaskID` 未导出 | `git grep -n newTaskID` | 定义 `:1107`，生产调用只有 `:322`/`:333` 两枚（都在 `Loop` 内部） | **对** |
| `docs/evidence/s1/154-host-id-never-closed-r1.md §2.3 (clauses G1/G1b)` | `git show 76c89b2:<该件> \| grep -nE '^### 2\.'` | §2.1–2.5 **全在**（`### 2.3 本程采用的门＝五枚子句` 在 `:137`） | **对：指针落点当时已在盘上** |
| 同件 `§2.5` | 同上 | `### 2.5 "并发那一形零读数"` 在 `:177` | **对** |
| `Q-56` | `git show <锚>:docs/reports/pending-and-issues.md \| grep -c Q-56` | `76c89b2`＝**9**、`5766be9`＝在 `:1095` 有整行登记 | **对**，且注释只引编号、**没替它选支**（ⓐ/ⓑ/ⓒ 一字未提） |
| `.scratch/wisp/probes/154/gate-clauses.sh`（"one `git grep` away" 的落点） | `git cat-file -e 76c89b2:.scratch/…/gate-clauses.sh` | 存在（rc=0） | **对** |
| `DEFERRED(C25-loop-wiring) item (1)`（原句） | `git grep -n` | `internal/risk/provenance.go:55` 真有这枚标记 | **对** |
| `ticket 151` / `ticket 154` / `ticket 19` | `git cat-file -t 5d46f24`＝commit；`ls` 三枚票面 | 全在（19 的票面 `:62` 就是那枚 DEFERRED 的登记行） | **对** |

⇒ **"不许预先引用尚未产出的读数"这一条：成立**——注释没有引用任何"以后才会有"的数（它不含任何读数型句子：没有"115/79"、没有"两跑皆红"、没有枚数），指向的两件（证据件小节、门脚本）在**写它的那枚 commit 上就已入库**。本程特意去数了注释全文里的数字：**只有行号与票号**（`366`/`332`/`321`/`151`/`154`/`56`/`2.3`/`2.5`），无一枚是"读数"。

### 2.3 三段文案逐段判"是不是打发"（判据＝票面 AC#2 那句"只读注释的人能不能答出'我这条腿要不要自己关'"）

| 段 | 承重的句子 | 指回的盘上事实（本程现量，不是他们自述） | 判 |
|---|---|---|---|
| A `OpenTask` 指路句 | "The open side is per CALL (mark above), the close side is per TASK and lives in another package, so nothing here keeps the two sides in step by construction" | 开＝`bridge.go:559` 在 `mark()` 里（每次调用）；合＝`cmd/wisp/run.go:563`（每枚任务）——本程用 **AST 名册**独立数过：`OpenTask` 生产调用者 **0 枚**（只被同文件的 `mark` 调）、`CloseTask` 生产调用者 **1 枚**（`run.go:563`） | **不是打发**：这一句给的是"为什么编译器帮不了你"，且两侧枚数就是它的证据 |
| B1 `WHICH LEGS STILL HAVE NO OWNER HERE` | 自测问句 **"does the TaskID you dispatch with come from that boundary?"** ＋ 三枚举例（panel/ball host 自造 id、retry wrapper、scheduled task 跨轮带同一枚 id）＋ **"nothing in this tree closes it, and no test will tell you so"** | "没有别人关"＝G2 名册 2 行（注释＋唯一调用）本程复算逐字同；"**no test will tell you so**"＝本程把全仓 `_test.go` 里 `OpenScope/OpenTask/CloseTask/open_scopes` 的命中逐枚点开（**33 枚**＝仓内 28 ＋ `.scratch` 探针文本 5；后者目录以点开头，`go` 工具结构性忽略 ⇒ 不进任何包，台件＝`probes/154-accept/ac2-test-mention-scan.txt`）——**没有一枚是"普查式"断言**（`internal/risk/*_test.go` 那些是自己开自己关的用例；`cmd/wisp/task_scope_close_151_test.go` 两枚断言的是 revoke 那一形） | **不是打发**：读者拿这一句能自答"我这腿要不要自己关"＝**要**，且知道没有测试会提醒他 |
| B2 `Being first is also an API decision` | "**Loop has no exported method that takes a caller-supplied task id**, so a host cannot route its id through the boundary that owns the close. That choice is `Q-56` and **this file does not answer it**." ＋ "closing your leg is only half of what ticket 154 owes - the reading for two tasks with scopes open at once is **still zero**, and it has to be measured in the same change (§2.5)" | 本程**不用他们的正则**证这一句：把 `*Loop` 的**全部**导出方法枚举出来（`probes/154-accept/ac2-loop-exported-methods.txt`）＝`Budgets`/`History`/`Reset`/`Steer`/`RunAsync`/`Run` 六枚，**无一枚收 taskID**；收 taskID 的 `l.run(ctx, taskID, input)` 未导出 ⇒ 句子成立。"并发读数仍为零"＝本程的 **AST 同函数序列扫描**独立复算（见 §2.4） | **不是打发**，且这一句是本票最值钱的一句：它把"缺的不只是 defer，而是一枚 API"写死在读者必经的位置，同时**没有**替 `Q-56` 选支 |

### 2.4 本程替 AC#2 补的一发校验：'still zero' 那句到底零不零

r1/r2 的 §2.5 用的是"数 `unbound-scope` 命中"这把尺。本程换尺＝**按函数体扫 `OpenScope`/`CloseScope`/`OpenTask`/`CloseTask` 的调用序列**（AST，457 枚 `.go` 全解析成功、0 枚 parse-error）：

```
同函数内 ≥2 枚开/合的全部用例：
  internal\risk\provenance_test.go#TestDisposalScopeClearsTaints        open=2 close=1   （session-1 先 Dispose 才开 session-2 ⇒ 串行）
  internal\risk\provenance_test.go#TestInspectUnknownScopeIsEmptyStore  open=2 close=2   （task-1 关后才开 task-2 ⇒ 串行）
  internal\risk\provenance_test.go#TestFragmentThresholdStricterAllowed open=2 close=0   （两枚是 p 与 p2 两枚**引擎实例**，不是同树两枚 scope）
  internal\risk\taintmatch_test.go#BenchmarkMarkFullSource              open=2 close=1   （同一枚 id "s"）
⇒ "两枚任务同时各开一枚 scope 且都活着"这一形：全仓**零枚**用例 ⇒ 注释那句 "still zero" 成立。
```
⚠ 但本程要给 §2.5 那句"**全仓没有一枚用例让两枚 scope 同时开着**"记一枚措辞缺陷 **D-7**：`internal/risk/provenance_test.go:296` 的 `TestScopesNeverInherit` 一边 `OpenScope("task-2")`、一边让 `task-1` 手里有污点（`Mark("task-1", …)`，且 `provenance.go` 的 `Mark` 注释明写"Marking a closed/never-opened scope **still records the taint**"），而那枚用例自己的失败文案写的是 **"a second concurrent scope must never see task-1's taint"**。⇒ 按"注册过且非空的 scope 有两枚"读，§2.5 成立；按"引擎里同时存在两枚带状态的 scope"读，**有一枚用例就在做这件事、还自称 concurrent**。这一处不影响任何判语（它量的正是同一形），但下一位若拿"全仓没有并发读数"去推"并发无仪器"会被这句绊一下。

### 2.5 本格抓到的一枚文案缺陷（D-6，实质）

`CloseTask` 新段落的标题句是 **"WHICH LEGS STILL HAVE NO OWNER HERE"**，但列出的三枚全是**假设形**（host 自造 id／retry wrapper／scheduled task）。**今天真的在开着、且没人关的那一形不在名单上**：

```
$ git grep -nE 'OpenScope\(|CloseScope\(' 5766be9 -- cmd ':!*_test.go'
5766be9:cmd/wisp/panel_assets.go:232:	prov.OpenScope(taintSourceScopeID)     ← cmd/ 下 CloseScope 枚数 = 0
$ git show 5766be9:cmd/wisp/panel_assets.go \| sed -n '163,166p'
  "…a CLI probe has no task, so it names one and **closes nothing** - the engine lives and dies inside this process."
$ 入口现读＝cmd/wisp/main.go:103  case "panel-assets":（生产命令面，不是测试夹具）
```

⇒ 判读分两面写，免得本程说过头：
- **不算谎**：注释里"legs"那三枚举例的限定语是"dispatch on the bridge"，而 `panel-assets` 那一枚**不过桥**（它直接调 `Provenance.OpenScope`），所以严格讲它不属注释正在答的那一形；且它的后果在它**自己的落点**（`panel_assets.go:163-166`）写得很老实、也有界（进程一次性退出、方向 fail-closed）。
- **但确实是缺**：本票 AC#2 选 `CloseTask` 的理由是"那是会被读到的地方"，而票 154 自己的普查（§4.1 表第 2 行）把这一枚列为**"同族第二枚真只开不合"**。一份标题写着"哪些腿还没有 owner"的清单里**没有今天唯一的实际实例**，下一位容易读成"这一族还只是假设"。
⇒ 记 **D-6**：建议补一句具名指路（"今日已存在一枚不过桥、直接开 Provenance scope 的诊断腿＝`cmd/wisp/panel_assets.go:232`，它有界且自带说明，但同族"）。**本程不改注释面**（写面只有两枚）。
⇒ 与 §1.5 的 **D-5 同源**：门（G1/G1b）不吃不过桥那一形、文案（B1）也没点名它——**两件交付各自缺了同一枚形状**，这是本格最该被下一位读到的一条。

### 2.6 AC#2 判语

**成立。** 三段都不是打发：每句能指回一件盘上事实（§2.3 那栏的右边全是本程现量，行号按两枚锚点各取一次）；"不许预引未产出读数"这一条**通过**（注释零枚读数型句子，两处指针在写它的那一刻已入库）；票面那句判据（只读注释的人能否自答"我这腿要不要自己关"）**能答**，因为自测问句与"no test will tell you so"都在。
带两枚缺陷：**D-6**（"还没 owner 的腿"清单缺今天唯一那枚实例）与 **D-7**（"没有一枚用例让两枚 scope 同时开着"的措辞比它的判据宽）。

---

## 第 3 节 格 AC#3 — 同族普查：换尺（AST）重扫 ＋ 判它有没有漏族

派单要本程"换尺重扫（它用 `git grep`；你换 `--name-only`×逐枚、或 AST／`go list` 级）＋判它有没有漏族：注册/注销、加/解监听、`Open`/`Close`、`Acquire`/`Release`、`Defer(` 挂在今天不可达的边上"。

### 3.1 本程的尺与它的射程（先把尺写死，免得"零命中"没主语）

```
取版＝git ls-tree -r --name-only 5766be9 -- internal cmd → 457 枚 .go
     ＋ git cat-file --batch（stdin 走文件、不走管道，按派单的坑②办）→ 457 枚 blob 全部落 D:/tmp/wisp-154-accept/tree/
尺＝go/parser 单文件语法树（stdlib、离线）：枚举 ①名字落在 lifecycle 动词表里的声明（含接收者类型）
    ②所有 selector 调用 x.Foo() 的 (Foo, 接收者表达式文本, 所在函数, 生产/测试, 文件:行)
读数＝声明 122 枚／调用点 2108 枚／**parse-error 0 枚**（台件 probes/154-accept/ac3-ast-roster-5766be9.txt，2260 行）
第二棵＝tools/**（他们 §4.3 明写"未穷举"的那支）：15 枚 .go，同一把尺，parse-error 0
        （台件 D:/tmp/wisp-154-accept/ast-roster-tools.txt，读数转记在 §3.4 第 6 条）
```

⚠ **本尺的三处结构性吃不到，先于任何结论写出来**：①只认 `x.Foo()` 形状 ⇒ **裸标识符调用**（`unregisterAll(...)`、`closeJob(...)`）不在名册里；②**接收者只到表达式文本，不到静态类型** ⇒ 配对表是"按方法名"配的，一枚 `Close` 会被算给所有叫 Close 的东西（`llm.Usage.Add` 那种累加器因此会被误点名为"只开不合"）；③动词表**没含 `Notify`** ⇒ `signal.Notify` 这一族本尺整个漏掉，是 §3.4 第 3 条用派单点名的族名去手查才捞回来的。
⇒ 第③条与本票 §4.1 末段 r1 自己那枚"整词尺吃不到 `pUnregisterHotKey`"是**同一类错的两种形**：本程不据此说他们的尺差，只说**两把尺在这里同盲**（`unregisterAll` 我的 AST 也看不见），所以他们的第 13 对本程**无法用 AST 尺复核**，只能用他们的 grep 尺换个旗重跑（下 §3.3 第 13 行）。

### 3.2 与他们的表交叉核对：两把尺都覆盖的行，读数相不相同

| 他们的行 | 他们的数（grep 尺） | 本程的数（AST 尺） | 判定 |
|---|---|---|---|
| 1 `OpenTask`↔`CloseTask` | 开 1（`bridge.go:559`）／合 1（`run.go:563`） | 开 1（`C OpenTask on=b in=mark PROD internal\tools\bridge.go:559`）／合 1（`C CloseTask on=rt.bridge in=admitTask PROD cmd\wisp\run.go:563`） | **同**，且 AST 还给出他们没写的"所在函数名"（`mark`／`admitTask`）⇒ "开按 call、合按 task"那句有第二把尺撑着 |
| 2 `OpenScope`↔`CloseScope` | 开 2／合 1 | 开 2（`bridge.go:642`、`panel_assets.go:232`）／合 1（`bridge.go:691`） | **同**（本程另在 §1.5 指出门不吃这一对） |
| 6 `Registry.Register` | 1（`run.go:345`）／`RegisterProvider` 0 | `C Register on=reg in=assembleRuntime PROD cmd\wisp\run.go:345` ＝ 1 | **同** |
| 8 `StartInJob`↔`Close` | 开 1／合 3 | 开 1（`slo_windows.go:498`） | **同**（合那侧本尺按名字吃到 88 枚，不可用＝§3.1 第②条） |
| 4/5 `DisposalScope` | `NewDisposalScope`/`Defer` 生产 0 | 本尺点名 `DeferNamed` 有 1 枚 PROD ⇒ 点开＝`disposal.go:187`，是 `Defer`→`DeferNamed` 的**包内委托**；`NewDisposalScope` 生产调用 **0 枚**（只有 `retention_test.go:193`、`disposal_test.go:17/168`、`provenance_test.go:307` 三处测试） | **他们的行成立**；本尺那一枚是**假阳性**，本程自己拆掉（放水两问之②的答案在这里） |

⇒ 结论：**他们表里本程能复核的行，两把尺数一致**（5/5 行）；本程没有一行是拿他们的数抄的。

### 3.3 本尺新点名、而他们的 14 对里没有的族（逐枚判读，含"这不是缺口"的）

| # | 族（AST 尺点名） | 开那侧生产出现点 | 合那侧生产出现点 | 判读 |
|---|---|---|---|---|
| N1 | `memory.Open` ↔ `(*Store).Close` | **2**：`cmd/wisp/run.go:313`、`cmd/wisp/providers.go:180` | **2**：`run.go:517`（`_ = rt.store.Close()`）、`providers.go:185`（`defer … store.Close()`） | **对称、不是缺口**；但⚠ `run.go` 的开与合**不在同一枚函数**（开在 `assembleRuntime`、合在 `execute`）⇒ 本程"同函数窗口"那把小尺**第一发把它读成 0 枚 Close＝假缺口**（见 §3.5 的自抓 bug）。真读数以这两枚点名为准。他们表里缺这一行＝完整性问题，不是漏判 |
| N2 | `Store.StartTaskLog` ↔ `FinishTaskLog` | **2**：`cmd/wisp/run.go:619`、`internal/agent/journal.go:129` | **4**（`FinishTaskLog`，尺＝`git grep -c 'FinishTaskLog('` 生产侧） | **有 owner、不必修**（任务行的开合都有人）；他们表里缺这一行。本程没去数"每条出口都调了 Finish 吗"（那是状态机那一票的账，§3.6） |
| N3 | `signal.Notify` ↔ `signal.Stop` | **3**：`cmd/balldebug/main.go:293`（`*hold` 分支）、`:306`（`*stay` 分支）、`:564`（`waitForExit`） | **1**：`cmd/balldebug/main.go:571` | **真不对称**（3 开/1 合），但两枚未关的都紧接进程退出 ⇒ 后果与他们的 pair 6/13 同形（注册表随进程寿命、调试宿主、非产品路径）。登记、不判"必修"；**这一族他们的表里没有**，而它正是派单点名的"加/解监听" |
| N4 | `llm.RegisterProtocol`（`anthropic/openaichat/openairesponses` 三枚 adapter 的 `init()`） | **3**（`adapter.go:100/:38/:75`） | **0**（`Unregister*` 全树 AST 零枚） | 与 pair 6 同判：**启动期一次填、进程寿命**，不必修；但它是**注册即无注销**的第四枚实例，他们表里没这一行 |
| N5 | `llm.BucketLimiter.Acquire` ↔（无 Release） | **1**：`internal/llm/ratelimit.go:223`（`Stream` 内） | **不存在该名字** | **不是缺口**：令牌桶没有"归还"语义；本尺按名字点了名，判读归"不成族"。写出来是因为它长得像 `Acquire/Release` |
| N6 | `audio`：`wasapiOpener.Open` ↔ `deviceStream` | `m.opener.Open` **1**（`wasapimic_windows.go:263`）＋ 全树生产 `Open` 共 16 枚 | 合那侧**不叫 Close**（interface 那侧是 `Stop`/release 闭包） | 本尺的**形状盲**（§3.1 第②条）：按名字配不上 ⇒ 不构成对 pair 10 的反驳。他们的第 10 对（构造函数零生产持有者）本程无法用 AST 尺独立复核，只能记"两把尺各判各的" |
| N7 | `llm.Usage.Add` / `*Guard.AddUsage` / `*Cost.AddCost` | Add 系生产 **48** 枚（其中 `Usage.Add` 40） | 0 | **假阳性**（累加器不是资源）。列出来是为了让下一位不要把本尺的输出当判据 |

### 3.4 派单点名的五族，逐族答一句（盘得到的写数，盘不到的写命令）

1. **注册/注销**：`Register` 生产 4 枚（N4 的 3 ＋ `run.go:345`）／`Unregister` 生产 **0 枚**（AST 名册里 selector 名以 `Unregister` 开头的枚数＝0）；ball 的 `unregisterAll` 属裸标识符调用、**两把尺都看不见** ⇒ 这一子句**盘不到**，命令＝`git grep -inE "unregister" 5766be9 -- internal cmd ':!*_test.go'`（21 行，r1 的尺）。判读：注册那侧今天全是启动期，注销不存在＝有界。
2. **加/解监听**：`signal.Notify` 3/1（N3）；`slog`/事件总线侧本程只盘到 `internal/observe` 的订阅不是配对外形 ⇒ 命令＝`git grep -nE "\.Subscribe[A-Za-z]*\(|AddListener|RemoveListener" 5766be9 -- internal cmd`（**零命中，尺已记 rc**）。**没有"应当没有"这句话**。
3. **`Open`/`Close`**：生产 `Open` 16 枚逐枚点名（本程窗口尺 15 枚同函数有 Close、1 枚跨函数有 Close＝N1、1 枚是 deviceStream 形状＝N6）⇒ **本族本程没有新缺口**，但结论的来源是本程那把有 bug 的小尺＋人工点开，不是名册本身。
4. **`Acquire`/`Release`**：`proc.AcquireSingleInstance`↔`Release` 1/1（与 pair 9 同）；新增 N5（不成族）。
5. **`Defer(` 挂在今天不可达的边上**：他们的 pair 4/5 是本票最有价值的一行（`provenance.go:56/:348` 那句义务写在 `DisposalScope` 上、而 `NewDisposalScope` 生产零使用者）。本程 AST 复核＝**成立**（`Defer` 系生产仅包内委托 1 枚、`Dispose` 生产 0 枚），且**没有新发现第四枚**。

### 3.5 本程自己弄错又改回来的一枚（不写这句，§3.3 的 N1 就是假缺口）

本程给"同函数体内有没有 `.Close(`"那把小尺的第一发读数，把 **两枚** 站点报成 0：`cmd/wisp/run.go:313`（真因：Close 在别的函数）与 **`internal/tools/fs_write.go:502`**（假因：**函数签名跨两行时，窗口在 `{` 出现之前就 `depth<=0` 返回了**）。
现读更正＝`fs_write.go:506` 明写 `defer src.Close() //nolint:errcheck // read-only handle` ⇒ **不是缺口**，是本程尺的 bug。修法是"只在见到第一个 `{` 之后才开始判 `depth<=0`"，修完重跑 16 枚站点的读数见 §3.3 的 N1（`providers.go:180` 也从 0 变 1）。
⇒ 与本票 §4.1 末段 r1 那枚（整词尺吃不到 Win32 名）**同形不同人**：本程把它写进证据件，是因为**下一位拿本程的尺复算时会被同一枚坑咬**。

### 3.6 AC#3 判语

**成立。** 判据本体（票面："逐枚答'它的关那侧有几枚生产调用者'，并给命令"）他们给了：14 对 × 逐枚生产枚数 ＋ 5 道可见性锁 × 现量命令，且**"盘不到"那三处写的是命令与原因、不是"应当没有"**（§4.3 原文三条）。本程换尺复核了能复核的 5 行，**逐行同数**；新点名 7 族，其中 **N3（`signal.Notify` 3 开/1 合）与 N4（`llm.RegisterProtocol` 3 开/0 合）是他们表里确实没有的两行**，N1/N2 是完整性缺行（两侧都有 owner），N5/N6/N7 是本程尺的假阳性或形状盲，全部具名拆掉。
⇒ **没有一枚足以把本格判成不成立**：漏的都是"启动期一次填、随进程寿命"那一类有界形，与本票已登记的 pair 6/13 同判；真正需要下一位做的是把 N1–N4 补成表里的 4 行（本程不代他们写票面）。
⇒ 本程**补上了他们 §4.3 第 3 条自己承认没穷举的两处**：`cmd/**` 四枚宿主（AST 全量，本程树覆盖）＋ `tools/**` 15 枚 `.go`（第二棵尺，读数＝mockllm 的 `listen` 拿到的 `ln` 生产侧无 `Close`，但它 `Serve(ln)` 阻塞到进程退出＝有界；`d22scan`/`signmodels` 无配对 API）。⇒ 那一条现在可以写成"盘到了、命令是这些"，不再是"未穷举"。
