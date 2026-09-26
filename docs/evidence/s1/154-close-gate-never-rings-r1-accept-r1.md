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
