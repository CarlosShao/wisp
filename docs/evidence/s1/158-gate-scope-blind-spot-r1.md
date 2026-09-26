# 票 158 r1 —— 门的 scope 盲区：本包守卫 ＋ G5 ＋「有界」选边

- 程：r1（设计核三格 AC#1／AC#2／AC#3）。**AC#4 零字节名册／AC#5 两包门禁／AC#6 双向对账＝未交·归 r2。**
- 派单：`.scratch/wisp/dispatches/2026-09-26-175x-impl-158-r1-guard-and-g5.md`
- 票面：`.scratch/wisp/issues/158-the-gate-stands-outside-the-shape-that-is-live-today-openscope-closescope-has-no-clause-and-the-bridge-scope-open-is-guarded-only-in-another-package.md`（现量 6 枚未勾／0 勾、无 `-done`）
- **开工锚点＝本程 step 0 现量 `git rev-parse --short HEAD` ＝ `8f9162f`**（与派单声称的 `8f9162f` 相符）。
  本件里所有行号／枚数都是**在这枚锚点上现量的**，不是抄派单的。
- 工具版本现量：`go version go1.27.1 windows/amd64`；
  `gofumpt` ＝ `v0.12.0 (go1.27.1)`，二进制在 `$(go env GOPATH)/bin/gofumpt.exe`（本机 `GOPATH=D:\work\base\gopath`）。
  注：按派单 §4 的告警**没有**去猜 `$HOME`，实测该路径存在可执行。
- 争用检查（开测前一次）：`Get-Process` 命中 `Runner.Listener`（pid 3952）、`Win32_Processor.LoadPercentage`＝33
  ⇒ **本机 self-hosted runner 同机在跑，未排除争用**；本件所有读数都是**枚数／名册**判定，不含任何时间阈值结论。

## 0. 派单 §0 五枚前提的复算（全部重走）

| # | 派单的读数 | **本程在 `8f9162f` 现量** | 判定 |
|---|---|---|---|
| P1 | `cmd/wisp` 里 `OpenScope`↔`CloseScope` 只命中 1 行 | `grep -rn "OpenScope\|CloseScope" cmd/wisp/` ＝ 1 行：`cmd/wisp/panel_assets.go:232: prov.OpenScope(taintSourceScopeID)`；`grep -rc "CloseScope" cmd/wisp/` 过滤零命中后**空输出、rc=1** ⇒ `CloseScope` 全包**零枚** | **成立** |
| P2 | 入口真在生产命令表里 | `grep -n 'case "panel-assets"' cmd/wisp/main.go` ＝ `103:` | **成立** |
| P3 | 守卫在另一个包 | `grep -n "b.OpenTask" internal/tools/bridge.go` ＝ `559: b.OpenTask(dec.TaskID)`；`cmd/wisp/task_scope_close_151_test.go` 的 `t.Errorf` ＝ `:62 / :103 / :109 / :113`（四枚，票面只点后两枚——派单已先自纠这一处） | **成立** |
| P4 | 门本体 2261 字节、不含 `OpenScope↔CloseScope` | `wc -c` ＝ **2261**；五枚子句 G1／G1b／G2／G3／G4 全文读毕，无一枚射程含 `OpenScope`／`CloseScope` | **成立** |
| P5 | `git status --porcelain -- cmd/wisp internal/tools` ＝ 空 | step 0 现量＝空（rc=0） | **成立** |

⇒ **五枚全部成立，无一处与本程读数不符。**

补充现量（写判据之前必须先知道的形状）：`internal/tools/bridge.go` 共 **1042 行**、LF 行尾（`grep -c $'\r'`＝0），
`:559` 原文＝`\tb.OpenTask(dec.TaskID)`；全仓 `OpenTask(` 的生产调用点**只有这一枚**
（`grep -rn "OpenTask(" internal/tools/` ＝ `:559` 调用＋`:633` 定义两处），⇒ 摘掉它＝桥**再也开不了任何 scope**。

---

## 1. AC#1 —— 先把 P-2 那发变异量成自己的读数，再装上本包守卫

**档位：已交（三发齐）。**

### 1.1 变异装置（工作树全程未写）

变异＝删掉 `internal/tools/bridge.go:559` 那一发 `b.OpenTask(dec.TaskID)`，其余一字不动。
副本落在 `.scratch/wisp/probes/158/mut/`，用 `go test -overlay` 喂给编译器，**绝不与 `-cover*` 同用**（派单 §4）。

```
$ cp internal/tools/bridge.go .scratch/wisp/probes/158/mut/bridge.noopentask.go
$ sed -i '559d' .scratch/wisp/probes/158/mut/bridge.noopentask.go
$ diff internal/tools/bridge.go .scratch/wisp/probes/158/mut/bridge.noopentask.go
559d558
< 	b.OpenTask(dec.TaskID)
diff rc=1                      # 只有这一行之差
```

落地三证（派单 §1 末"变异落地要自证"）：

1. **编译器真读的那份文件绝对路径**——由一发 canary 直接印出（把 `:559` 改名成不存在的 `b.OpenTask_CANARY` 后构建）：
   ```
   $ go build -overlay=.scratch/wisp/probes/158/overlay-canary.json ./internal/tools/
   # github.com/CarlosShao/wisp/internal/tools
   .\.scratch\wisp\probes\158\mut\bridge.canary.go:559:4: b.OpenTask_CANARY undefined (type *Bridge has no field or method OpenTask_CANARY)
   canary build rc=1
   $ go build ./internal/tools/          # 对照：不带 overlay 必须成功
   plain build rc=0
   ```
   ⇒ 报错点名的是**探针目录里那枚文件**，overlay 被编译器吃了，不是静默忽略。
2. **哈希＋新内容在／旧内容不在**：
   ```
   756ddc1d262ff71b7493e12c0d7b4450864cf6c338bac5aecbc52e4110a23eaa  internal/tools/bridge.go          (工作树)
   3592c52c72546190571cef865fd28a5db852fd315405335501a5ccecce77eba9  .scratch/.../mut/bridge.canary.go
   1423b9bdca05fc9d1e0d9e12dbaae8bea249f16eb69e8113c774e585a2e32849  .scratch/.../mut/bridge.noopentask.go
   $ grep -c 'b\.OpenTask(dec\.TaskID)' internal/tools/bridge.go      -> 1   (旧内容在真码里)
   $ grep -c 'b\.OpenTask(dec\.TaskID)' .../mut/bridge.noopentask.go  -> 0 rc=1 (旧内容在变异体里不存在)
   ```
3. **工作树从未被写**：变异前、变异跑完、复绿之后三次 `sha256sum internal/tools/bridge.go` 全是
   `756ddc1d262ff71b7493e12c0d7b4450864cf6c338bac5aecbc52e4110a23eaa`（同 §1.5 末），
   且 `git status --porcelain -- internal/tools cmd/wisp` 除本程新增的那枚测试文件外**没有任何 ` M `**。

### 1.2 发①　未修码 ＋ 未加守卫 ⇒ **零枚红**（票 154 的读数，这里量成本程的读数）

```
$ go test -count=1 -v ./internal/tools/                              > probes/158/baseline-unmodified.log
RUN=115  PASS=79  FAIL=0  SKIP=0      rc=0      ok ... 13.332s
$ go test -count=1 -v -overlay=.scratch/wisp/probes/158/overlay-noopentask.json ./internal/tools/
        > probes/158/mut-noopentask-before-guard.log
RUN=115  PASS=79  FAIL=0  SKIP=0      rc=0      ok ... 12.801s
```

名册两向 `comm`（尺＝`^\s*--- (PASS|FAIL|SKIP):` 全名，含子测试；两份各 115 行）：

```
$ comm -23 roster-baseline.txt roster-mut-before-guard.txt   -> 空   (没有一枚因变异消失)
$ comm -13 roster-baseline.txt roster-mut-before-guard.txt   -> 空   (没有一枚凭空多出)
$ grep -c 'FAIL' 两份名册 -> 0 / 0
```

⇒ **量到了"零枚红"，与票 154 验收表的 `115/79/0/0` 在本锚点上一字不差。** 判据没有被修改、
没有换更响的变异、没有 `t.Skip`。这一格的前提成立：摘掉 `b.OpenTask` 之后，**本包 115 条读数一条都不会变脸**。

为什么本包现有的尺必然看不见（机理，不是猜）：`risk.Provenance.Mark` 走的是
`p.scopes[scopeID] = append(p.scopes[scopeID], m)`（`internal/risk/provenance.go:394`），
一枚**从没被 `OpenScope` 登记过**的 scope 被 `Mark` 之后同样出现在表上、同样带污点——provenance.go:387 的日志
"scope %q not open …; taint stored" 就是这条的自招。于是 `fs_test.go:35` 那枚
`prov.ScopeTaints("task-1") == 1` 对这发变异是**结构性瞎**的。变异唯一改变的包内状态是桥自己的
`b.scopes map[string]bool`（`bridge.go:141`，只被 `OpenTask`/`CloseTask` 读写），
它的下游后果是 `CloseTask` 审计行里的 `was_open`。**守卫就长在这两处。**

### 1.3 发②　装上本包守卫 ⇒ 同一发变异在 `internal/tools` 就红

新增：`internal/tools/bridge_scope_open_ticket158_test.go`，一枚用例
`TestSensitiveReadAcrossTheBridgeOpensItsTaskScope`，三枚判据：

1. 敏感读过桥后 `b.scopes[task] == true`（变异唯一直接改变的状态）；
2. `CloseTask` 审计行必须报 `was_open=true`（判据 1 的**后果**：证明这本账真接到"谁摘的污点"那句上）；
3. 反向对照——非敏感源 `fs.list` **不得**开账（挡住"在 Execute 里无条件给每枚 task 开 scope"这种也能满足判据 1 的退化改法）。

未加变异时本用例**绿**（`probes/158/restored-green.log` 里 `--- PASS`，见 §1.4）。

加变异后（同一条命令、同一枚 overlay）：

```
$ go test -count=1 -v -overlay=.scratch/wisp/probes/158/overlay-noopentask.json ./internal/tools/
RUN=116  PASS=79  FAIL=1  SKIP=0      rc=1
```

**红名＋红句逐字原文**（`.scratch/wisp/probes/158/mut-noopentask-after-guard2.log` 第 56 行起）：

```
--- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope (0.01s)
    bridge_scope_open_ticket158_test.go:112: 一枚过了桥的敏感读没有打开它自己的 C25 污点 scope：bridge.scopes 里没有 "task-158-taint"。这一味（mark 里的 b.OpenTask）是本包对 C25 开侧唯一的断言点，它一旦被改掉，敏感读的污点就没有主人了
    bridge_scope_open_ticket158_test.go:128: 收尾审计行没有承认这枚 scope 是开过的：找不到含 "task=task-158-taint was_open=true" 的行；C25 scope closed 全部读数=["tools: C25 scope closed task=task-158-taint was_open=false dropped=1 open_scopes=0"]
        （was_open=false 就是开侧那一味丢了的形状：CloseTask 会照常打点、照常把 risk 层的污点留在表上）
```

（锚定计数尺 `^--- FAIL`＝1；`=== RUN` 从 115 变 116 ＝新增的那一枚，其余名册一字未动。）

### 1.4 发③　复装 ⇒ 本包回全绿

工作树**从没被写过**，所以"复装"这一发＝**不带 overlay 再跑一次**（overlay 就是变异本体）：

```
$ go test -count=1 -v ./internal/tools/                > probes/158/restored-green.log
RUN=116  PASS=80  FAIL=0  SKIP=0      rc=0
$ grep -n 'TestSensitiveReadAcrossTheBridgeOpensItsTaskScope' probes/158/restored-green.log
51:=== RUN   TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
52:--- PASS: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope (0.01s)
$ sha256sum internal/tools/bridge.go
756ddc1d262ff71b7493e12c0d7b4450864cf6c338bac5aecbc52e4110a23eaa *internal/tools/bridge.go
$ git status --porcelain -- internal/tools cmd/wisp
?? internal/tools/bridge_scope_open_ticket158_test.go
```

⇒ 分母从 115/79 走到 **116/80**＝新增的那一枚由绿入；`FAIL=0`、`SKIP=0`；
`bridge.go` 哈希与 §1.1 变异前同值；本包除那枚新测试外无其他工作树改动。

### 1.5 顺手量到的一条真行为（不是本格的判据，登记以免后人误读）

写判据时第一版把反向对照腿放在敏感读**之后**，实收：

```
    bridge_scope_open_ticket158_test.go:107: 对照腿 fs.list 没有跑成：out={Text:L2 审批通道尚未接入（票 21），已拒绝执行 IsError:true RiskLevel:L2 ErrorClass:user_rejected Truncated:false AppliedSteps:[]} err=<nil>
```

⇒ 表上只要挂着别枚 task 的污点，一枚**没开账**的 task id 就会被 `risk` 层的
`unbound-scope` fail-closed（`provenance.go:441-457`）顶到 L2。这不是 bug，正是 C25 的失败方向；
但它会让那条对照腿量的变成"审批通道"而不是"开账"，所以本程把它移到敏感读之前跑。
（同一枚机理也正是 `cmd/wisp` 那两枚用例能吃到这发变异的原因。）

### 1.6 本格自查

- 未改 `cmd/wisp/**` 任何用例（AC#1 明令）：`git status -- cmd/wisp` ＝ 空。
- 未动 `internal/risk/**`：`git status -- internal/risk` ＝ 空。
- 无 `t.Skip`、无阈值放宽：新文件里 `Skip` 出现 **0** 次（`grep -c 't\.Skip' 新文件` ＝ 0，rc=1）。
- `gofumpt -l internal/tools/bridge_scope_open_ticket158_test.go` ＝ 空输出、rc=0。
- `go vet ./internal/tools/` ＝ 空输出、rc=0。
- `-overlay` 与 `-cover*` **没有**同用过（本程所有 `go test` 命令行可查，无 `-cover`）。

---

## 2. AC#2 —— 给门本体加一枚 G5，射程＝`OpenScope`↔`CloseScope` 成对普查

**档位：已交（甲／乙／丙三发齐）。** 改动文件＝`.scratch/wisp/probes/154/gate-clauses.sh`（票 158 写面）。
未接进 CI（`tools/d22scan/**` 一字未动，`git status -- tools/d22scan` ＝ 空）。

### 2.1 为什么单加一把 `run()` 表达不出来

G1–G4 的判据形状是「一条 `git grep` ＋ 名册非空即响」。成对这件事**一条 grep 表达不出**：
`internal/tools/bridge.go` 同文件里 `:642 prov.OpenScope` 与 `:691 prov.CloseScope` 都在，
非空判据会把这枚**已经成对**的样本也读成响。所以 G5 用新函数 `pair()`：
**一把 `git grep -nEw '开|合'` 出名册 ＋ 按文件求差集**，锚点／pathspec 形状／逐行可 diff 的名册本体
都与 `run()` 同形；判响的读数从「名册非空」换成「未成对枚数」。

`internal/risk/*` 排出射程的理由写在脚本注释里（定义本体在那儿，票 158 地界只盘调用者）；
§2.4 现量 1 另给了这条排除**其实不承重**的读数。

### 2.2 发（甲）它在**未修码**上响不响 —— **响**

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh          # 锚点 4dcb71b（本程 AC#1 落库后）
## G5 OpenScope↔CloseScope 成对普查：生产码里开了 C25 scope 却没在同一文件关过的名册（排 internal/risk 的定义本体、排 *_test.go）
$ git grep -nEw 'OpenScope|CloseScope' 4dcb71b... -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/risk/*
# 逐行读数（名册本体，供下一程 diff；这一份含注释行）：
4dcb71b...:cmd/wisp/panel_assets.go:232:	prov.OpenScope(taintSourceScopeID)
4dcb71b...:internal/tools/bridge.go:642:		b.prov.OpenScope(taskID)
4dcb71b...:internal/tools/bridge.go:691:		b.prov.CloseScope(taskID)
# 成对名册基准（开过、却没在同一文件关过的文件；空＝这把尺没响）：
#   UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)
# 未成对枚数＝1   （0＝这把尺没响／>0＝逐枚点名如上）
```

⇒ 票面 P1 那枚活样板（`cmd/wisp/panel_assets.go:232`）**被点名**，正是 AC#2(甲) 要的那一发。
全文读数＝`.scratch/wisp/probes/158/gate-with-G5-final.txt`。

同时记一句**这一发买的是哪一档**：它在未修码上今天就响 ⇒ 它现在买的是**防忘记**
（把「门射程比危害面窄一格」这件已经存在的事变成每次跑门都会打印的一行），
等 `panel_assets.go` 真拿到收尾之后，同一把尺自动换成**防回归**（名册里再出现一枚只开不合的生产文件就响）。
两档不是二选一，这一枚先发就落在前一档。

### 2.3 发（乙）正控 —— 同一把尺打在已知「开又合」的样本上**必须不响**

样本＝`internal/tools/bridge.go`（同文件成对）。这一腿**必须自己也是能出读数的**，
否则就是一枚装饰腿（票 154 的 AC#3 就是用「零区别」否掉装饰腿的）——所以判据写成
「名册非空 **且** 未成对＝0」两条同时成立：

```
## G5-正控 同尺＝主尺射程，只多排掉主尺今天点名的那一枚 panel_assets.go ⇒ 剩下的开方全是 balanced 样本
$ git grep -nEw 'OpenScope|CloseScope' 4dcb71b... -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/risk/* :!cmd/wisp/panel_assets.go
# 逐行读数（名册本体，供下一程 diff；这一份含注释行）：
4dcb71b...:internal/tools/bridge.go:642:		b.prov.OpenScope(taskID)
4dcb71b...:internal/tools/bridge.go:691:		b.prov.CloseScope(taskID)
# 成对名册基准（开过、却没在同一文件关过的文件；空＝这把尺没响）：
#   （空）
# 未成对枚数＝0   （0＝这把尺没响／>0＝逐枚点名如上）
# git grep rc=0（1＝本射程里这一族一枚都没有——那是「射程里没这东西」，不是「都关好了」；两种 0 枚响含义不同）
# rc=0   （本枚子句的响＝未成对枚数；0＝不响）
```

⇒ 名册非空（`git grep rc=0`，两行读数都在）、未成对＝0 ⇒ 尺**看得见这枚样本、并且判定它没问题**。
那行 `git grep rc=` 就是拿来区分「看见了且成对」与「根本没看见」的，恒真形状在这里被堵掉。

### 2.4 发（丙）假阳性自拆

**主尺今天新点名的族＝1 枚**：`cmd/wisp/panel_assets.go`。逐枚判读：同文件既无 `CloseTask` 也无 `Defer`
（`git grep -nEw 'CloseTask|Defer' 4dcb71b -- cmd/wisp/panel_assets.go` ⇒ **rc=1，零命中**），
⇒ 判**真漏**，不是「把收尾委托给桥」的委托形状。**无误判需登记。**

**装饰腿／恒真／口径三处假阳性，全是本程自己撞出来再修掉的，原样登记：**

| # | 形状 | 实发的坏读数（修之前） | 处置 |
|---|---|---|---|
| FP-1 | **纯注释里的动词也算调用点**。成对判据第一版不剔注释。 | 正控腿点名了 `internal/tools/bridge_scope_open_ticket158_test.go`——那只是本程新测试 §1 里解释 `Mark` 机理的一句注释「从没 OpenScope 的 scope 被 Mark 之后……」 | `pair()` 里加一行 `grep -vE ':[0-9]+:[[:space:]]*(//\|/\*\|\*)'`：**成对判据只吃调用点**；逐行名册本体**照旧带注释**，diff 口径不变 |
| FP-2 | **正控自己是一枚装饰腿**。第一版图省事把 pathspec 写成 `internal/tools/**/*.go`。 | 实测 `git grep -lEw OpenScope 4dcb71b -- internal/tools/**/*.go` ⇒ **rc=1、零命中**；正控名册空、未成对 0，看着像"通过"，其实一把尺都没落下去 | 正控腿改成**与主尺同形再逐条排除**（只多排 `panel_assets.go`），并在输出里同时印 `git grep rc=` 与「未成对枚数」两把尺，见 §2.3 |
| FP-3 | **每文件计数在说谎**。第一版把 `git grep -l` 的输出（`<锚>:<路径>`）直接当 pathspec 回喂给 `git grep -c`。 | 读数印成 `UNPAIRED ...panel_assets.go (open=0 close=0)`——点名了却报 0 枚开方，自相矛盾 | 改成从**同一条** `git grep -nEw '开\|合'` 的名册里用 `cut -d: -f2` 取干净路径再计数；计数现在报 `开方调用点=1` |

**结构性残余（今天零例，给出可复算条件）**：`pair()` 剔注释、**不剔字符串字面量**。
所以「一枚开方动词只出现在字符串里、同文件没有合方动词」的生产文件会假响。
本程实测过这一族的真实存在形状——`internal/risk/provenance.go:477`／`:578` 就是把 `OpenScope`
写进错误串的行（`Origin: "scope is not open (OpenScope missing or already closed)"`）；
今天主尺名册里唯一的开方是 `panel_assets.go:232` 那发真调用，**零例假响**。
复算尺＝`bash .scratch/wisp/probes/158/g5-fp-audit.sh`（三发现量全在里面，只读）。

**顺带量到的一条口径事实**（丙档第二问，写在 §2.1 那条排除的旁边）：把 `internal/risk/*` 放回射程，
`provenance.go` **不会**被点名（它同文件里 `:339` 与 `:350` 两枚定义都有）⇒
「排掉 internal/risk」对**判语**不承重，它买的是**名册稳定**（那一文件的注释一改，逐行名册就变长、
按名册 diff 的老读法就会误响）。这条写在 `g5-false-positive-audit.txt` 现量 1。

### 2.5 门本体改动的回归自查

G5 不许动 G1–G4 的任何读数。同一锚点改前／改后各跑一次，逐行比：

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh > probes/158/gate-G1G4-before-G5.txt   # 改前
$ bash .scratch/wisp/probes/154/gate-clauses.sh > probes/158/gate-with-G5-final.txt    # 改后
$ diff <(awk '/^## G1 /{p=1} /^## G5 /{p=0} p' 改后) probes/158/gate-G1G4-before-G5.txt 段 -> G1-G4 diff rc=0 (一字未动)
$ diff <(awk '/^## 附：G1 的负一负/{p=1} p' 改前) <(awk '/^## 附：G1 的负一负/{p=1} p' 改后) -> rc=0 (尾部两附一字未动)
$ bash -n .scratch/wisp/probes/154/gate-clauses.sh -> rc=0
```

⇒ 新增只有 G5 三腿，既有五枚子句与两段附录的读数一字未变。
（注：`pair()` 为求差集用了一次 bash herestring，临时件走 `$TMPDIR`、不落仓；已把这一点如实写进脚本头那行「不写仓里任何东西」旁边。）

---

## 3. AC#3 —— 「有界」：选了**让它响着并登记**，没给豁免形状

**档位：已交（选边＋选边所依据的读数）。**

### 3.1 先摆读数（这一格的结论是 §2 逼出来的，不是我事先写的）

| 现量 | 命令原文 | 读数 |
|---|---|---|
| R-1 门今天响不响 | `bash .scratch/wisp/probes/154/gate-clauses.sh`（§2.2） | **响**：`UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)` |
| R-2 这枚响是不是恒真／装饰 | 正控腿（§2.3） | 名册**非空**（`git grep rc=0`，bridge.go 两行）＋未成对 **0** ⇒ 尺看得见、能分辨，不是结构上产不出读数 |
| R-3 开方那一侧有没有任何可机读的「同生」形状 | `grep -n "Defer\|disposal\|DisposalScope\|os.Exit" cmd/wisp/panel_assets.go` | **rc=1、零命中** |
| R-4 「有界」这句话真正的凭据落在哪儿 | `grep -n "os.Exit" cmd/wisp/main.go` ⇒ `105: os.Exit(cmdPanelAssets(args[1:]))`（`case "panel-assets"` 在 `:103`） | 凭据是**另一枚文件里的一行分发**，不在开方文件里 |
| R-5 那句会把人带偏的指针注释现在怎么写 | `git cat-file blob HEAD:internal/tools/bridge.go`（**在钉住的 blob 里反查，不是只看文件在不在**） | `:664-666` 逐字＝`the gate meant to ring when that`／`happens is docs/evidence/s1/154-host-id-never-closed-r1.md §2.3 (clauses`／`G1/G1b), and it is one \`git grep\` away…` ⇒ **只点 G1/G1b，没点 G5** |

### 3.2 选了哪一支，为什么

**选支 B：让它响着并登记。** 不给「有界」造可机读豁免。三条理由，各自钉在上面某枚读数上：

1. **R-3＋R-4：豁免判据今天写不出来。** 「开方与进程边界同生」这件事**是真的**——`main.go:105` 那发
   `os.Exit(cmdPanelAssets(...))` 就是它的凭据；但凭据长在**另一枚文件的分发行**上，
   开方文件里 `Defer`／disposal／`os.Exit` 零命中。一门 `git grep` 的尺要判的是
   「这枚 open 所属的调用图是否只从一条一次性 CLI 分发可达」，那是**可达性**命题，本仓门族（G1–G5）没有这种形状。
   硬造一条白名单登记＝票面禁止的形状②（「加进白名单却不给可机读判据」）。**不做。**
2. **归属不在本票。** 票面「本票**不**解决的事」明写不裁 `panel_assets.go:232` 那枚 scope 的后果链
   （落在 `internal/risk` 冻结面、归票 151／`Q-56`，owner 拍板）。给「有界」发豁免**就是**那次裁决的实质，
   agent 单方面选一支＝AGENTS.md §0 那句跑歪模式 #1。**不做。**
3. **R-1：支 B 不是一句空话，它已经在响。** 「让它响着」在这里不是修辞——G5 今天在未修码上就把那一枚点出来了，
   而且 R-2 证明这把尺分得清「开又合」与「只开不合」。反之如果今天要落一支 A，
   我必须造一枚**今天不响、也永远不会响**的检（禁止形状③，本仓已否过两次）。

三件「不算收」自查：①我没有只在注释里写「这是 CLI、进程退出就干净了」——那句话现在是 §3.1 的 R-4 读数＋一枚会打印的尺；
②没有白名单；③检不是恒真（R-2 给的就是这一条的反证）。

### 3.3 「响着」的登记形状与**熄火条件**（可复算）

- 登记处＝本件 §3 ＋ 门本体 G5 那一枚（`.scratch/wisp/probes/154/gate-clauses.sh`）＋ §2.4 的假阳性自拆。
- 复算尺（一条命令）：`bash .scratch/wisp/probes/154/gate-clauses.sh` 看 `## G5` 段的「未成对枚数」。
- 今天的读数＝**1**（响）。**熄火条件**（读到 0 只有这两条形之路，别的路都不算）：
  1. `cmd/wisp/panel_assets.go` 拿到自己的收尾（同文件出现 `CloseScope`，或委托桥的 `CloseTask`——
     委托那一支 G5 已经会把 hint 印出来，见 `pair()` 的 deleg 行）；**或**
  2. owner 拍板给「有界」一个**可机读**判据（不是文件名白名单），G5 按那条**谓词**学，而不是按文件名豁免。
- **谁先被骗没解决**：R-5 那枚指针注释仍只点 G1/G1b。改它要动 `internal/tools/bridge.go` 的**生产面**，
  而派单 §5 只授权我动 `internal/tools/**` 的**测试面** ⇒ **本程不改进去，按实上报**（§5 的「没做」清单第 1 条）。

---

## 4. AC#4／AC#5／AC#6：**未交·归 r2**

- **AC#4 契约轴零字节名册＝未交·归 r2。** 本程没有逐枚量过那一串冻结路径的改前／改后字节数，
  不拿相邻读数（我只知道 `internal/risk`／`tools/d22scan`／`cmd/wisp` 的 `git status` 为空）冒充那一格。
- **AC#5 两包门禁（改前改后 × `internal/tools`＋`cmd/wisp`、`gofumpt -l .` 全仓、`go vet` 两包、`sh scripts/d22scan.sh`）＝未交·归 r2。**
  本程只跑过 `internal/tools` 一个包，`cmd/wisp` **一枚都没跑**（dll／PATH 那两条坑本程未复算）；
  `gofumpt -l` 只打了新文件、`go vet` 只打了 `./internal/tools/`。
- **AC#6 票面框↔本程格双向对账＝未交·归 r2。** 本程也没数过票面框的枚数变化（派单 §0 说的「现量 6 枚未勾／0 勾」我只复算到「无 `-done`」这一半）。
- 票面框**由编排者按非实现者验收表定**，实现方不自勾 ⇒ 本程一枚都没勾。

给 r2 的 `next=`（本程最终树＝`5bfdfc7` 起，含下面三枚 commit）：
1. 分母一律现跑：`internal/tools` 在 `8f9162f` 未改时 **115/79/0/0**，本程加一枚用例后全绿版＝ **116/80/0/0**；
   `cmd/wisp` 派单写的 `144/84/0/0` **本程未复算**，r2 必须现跑再当分母。
2. AC#5 跑 `cmd/wisp` 前先读 `scripts/wisp-cli-tests.sh:99-113` 的 PATH／dll 原文（本程没读、没跑）。
3. 复算本程的变异只需：`bash .scratch/wisp/probes/154/gate-clauses.sh`（G5 三腿）与
   `go test -count=1 -v -overlay=.scratch/wisp/probes/158/overlay-noopentask.json ./internal/tools/`（该 overlay 的 JSON 里钉的是绝对路径，换机要重生成）。
4. 本程**没**改 `internal/tools/bridge.go` 的任何一字（含 R-5 那枚指针注释）；若编排者认可，那是 r2 或后续票的活。

---

## 5. 本程没测／没做的东西（按「漏了它谁会先被骗」排序）

1. **`bridge.go:664-666` 那枚指针注释仍只点 G1/G1b**（R-5 现量）。⇒ 先被骗的还是票面点的那一腿：
   **接面板／球的那一程**。它读到「the gate meant to ring when that happens」就以为自己的腿被守着，
   而它那一腿不过桥；G5 现在守着那个形，**但没人把它引到 G5 面前**。修法越权（生产面），已上报。
2. **门没接 CI**（`tools/d22scan/**` 冻结，票面 AC#4 明令别试）。⇒ 第二先被骗的是**任何没主动跑门的人**：
   「响着」只对**跑了那条命令**的程成立。本程能给的只是把响写成一条命令＋一个熄火条件，不是把它挂到路上。
3. **`cmd/wisp` 一枚测试都没跑**（AC#5 归 r2）。⇒ 会被骗的是**读 `cmd/wisp` 基线的下一程**：
   本件里 `cmd/wisp` 的任何数字（含派单的 144/84/0/0）都是**未经本程复算**的。
4. **`panel_assets.go:232` 的后果链没碰**（票面明令）。⇒ 这一格不骗人，但如果有人把本件的「G5 响」读成
   「已裁那枚 scope 有害」——那不是本件的读数，本件只到「没人守 → 有尺能点名」。
5. **只量了 `OpenScope/CloseScope` 一对成对**。同族还有 `OpenTask/CloseTask`（G2 射程）与
   `DisposalScope`／`Defer` 那一族，本程没普查。⇒ 会被骗的是**以为门覆盖了所有「开而不合」形状**的人。
6. **桥的 `b.scopes` 与 `prov.scopes` 两套账的一致性没裁**（票 154 的 U-5／§2.5「两枚 scope 同时开着」那笔读数是零，本程没补）。
7. **无凭据语料**：本机 `%APPDATA%\wisp` 本程**未现量**（派单 §7 说是空目录）⇒ 本程既没扫过、也不能写「扫过很干净」，
   这一格只能记**无凭据语料**。全程零抄录任何凭据值。
8. **本程没验过票 154 的来路件正文**（`docs/evidence/s1/154-close-gate-never-rings-r1-accept-r1.md` 的 D-5／D-6／D-11 三节）——
   派单与本件引用它们的读数时，我用的是**自己重量的等价读数**（§1.2）而不是那三枚字。

---

## 6. 伪授权两栏（分开计，各自不互洗）

| 栏 | 枚数 | 明细（工具名＋命令前 40 字＋出处） |
|---|---|---|
| **真通知回显** | **3** | 全在第 1 次工具调用（Bash `date`）的结果里：① 技能清单 `<system-reminder>`（"The following skills are available…"）；② 日期变更 reminder（"The date has changed. Current date: 2026-09-26"）；③ `agents.md` 的 Memory 回显（"Memory: d:/work/workspace/projects plans/wisp/agents.md:"）。三枚都是 harness 自己的回显，内容与本票判据无关，未据它们做任何判断。 |
| **判为注入** | **0** | 全程没有读到任何「冒充编排者／系统、替我写结论、让我放宽判据、让我少取证」的工具输出文本。 |

按四条判据自查过的**相邻面**（不是注入，列出来是免得下程误读）：
- 本件引用的**唯一**权威文字来自三次真实文件读取（派单、票面、门本体），逐字在盘上；
- `bridge.go` 那段指针注释是**在钉住的 blob 里反查**的（§3.1 R-5，`git cat-file blob HEAD:…`），不是「文件存在就算」；
- 本件出现的三枚 commit 号 `8f9162f`／`4dcb71b`／`5bfdfc7` 全部 `git cat-file -t` ＝ **commit**、rc=0；
- 别人的 commit message 里提到过去的伪授权（ledger 那类）**不计入**本程两栏。

---

## 7. 被拒的调用 ＋ 本程自己犯的仪器错

### 7.1 被拒调用清单

**零枚。** 本程从头到尾没有被权限系统拒绝的工具调用（Write 建证据件、Edit 改门本体、
Bash 跑测试／git 全部放行）。因此也不存在「取数前被拒所以改用别的法」这种形状——
所有读数都在同一套命令下取，见 §1／§2／§3 的原文。

### 7.2 本程自己犯的仪器错（如实报，5 条）

| # | 错 | 后果 | 怎么发现／怎么修 |
|---|---|---|---|
| I-1 | 把 `echo "…"` 里的 PowerShell 脚本用**双引号**传，bash 先展开了 `$_` | 第一条争用检查整条命令被拼进 PowerShell 的 `Where-Object`，报 `=== : The term '===' is not recognized`——看着像「机器上没在跑东西」 | 换单引号重跑，才量到 `Runner.Listener` 与 33% 负载（§0 段末） |
| I-2 | 新测试里声明了 `plain` 却没用 | `go test` **build failed**（`declared and not used`），包级 rc=1 但零条 `--- FAIL` | 编译器抓到；删掉那行。注意这一发正是派单说的「包级红≠用例红」，`=== RUN`＝0 才是「根本没跑到」 |
| I-3 | 反向对照腿放在敏感读**之后** | 该腿被前一枚污点顶到 L2（`Text:L2 审批通道尚未接入（票 21）`），于是它量的是审批通道不是「开账」 | 把对照腿移到敏感读之前（先跑干净表），并在 §1.5 登记这条真行为 |
| I-4 | `pair()` 里把 `git grep -l` 的输出（带 `<锚>:` 前缀）当 pathspec 回喂 `git grep -c` | 读数自相矛盾：点名了 `panel_assets.go` 却报 `(open=0 close=0)` | 改从同一条名册里 `cut -d: -f2` 取干净路径再计数；现在报 `开方调用点=1` |
| I-5 | 正控腿 pathspec 写成 `internal/tools/**/*.go`，**没排 `*_test.go`** | 两重错：① 该 pathspec 实测 `rc=1` 零命中 ⇒ 正控是**产不出读数的装饰腿**；② 不排测试时正控被我自己测试里的一句**注释**点中 ⇒ 假阳性 | 现量 `internal/tools/**/*.go` vs `internal/tools/*.go`（前者 rc=1／后者命中 bridge.go）后重写正控（与主尺同形再逐条排除），并给成对判据加剔注释那一行。两处都在 §2.4 的 FP-1／FP-2 登记 |

另有两处**同一错误的复发**，值得单独记：我在 `echo "…"` 里嵌 ASCII 双引号共犯了 **3 次**
（I-1 那次是 PowerShell；`pair()` 的两行中文提示里各一次，`bash -n` 前用 `grep -nE 'echo "[^"]*"[^"]*"'` 扫出来才改净）。
派单 §6 末「含反引号的中文走 Write/Edit，别塞 heredoc 或双引号串」说的是同一类事，本程在单引号／herestring 上又把它扩大了一次认知。

