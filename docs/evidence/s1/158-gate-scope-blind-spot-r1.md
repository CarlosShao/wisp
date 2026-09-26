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

### 2.6 三腿读数在**最终树**上复跑过一遍（免得 r2 以为只有 AC#1 落库那一次）

```
$ git rev-parse --short HEAD        -> 5bfdfc7      # AC#2 落库后
$ bash .scratch/wisp/probes/154/gate-clauses.sh | grep -E '^## G5|UNPAIRED|未成对枚数'
  G5 主尺   ：UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)  未成对枚数＝1   rc=1  （响）
  G5-正控   ：未成对枚数＝0   rc=0                              （不响）
  G5-负一负 ：UNPAIRED 5 枚（panel_assets + 4 枚 internal/risk 测试）
$ go test -count=1 ./internal/tools/  -> ok ... 12.796s   rc=0 （最终树全绿，非 -v 快跑一次）
```
（AC#3 那一枚 commit 只动本证据件，不改 Go 码也不改门本体 ⇒ 上面两行读数对 `7f37ea4` 同样成立。）

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

给 r2 的 `next=`（**别信本件正文里任何一枚写死的 commit 号，包括这一句**——本程的更正 commit 会把
"最终树"那一行自己刷过期，实发过两次。本程全部落库＝开工锚点 `8f9162f` 之后、下面这条命令列出的那些：

```
$ git log --format='%h %ad %s' --date=format:'%H:%M' 8f9162f..HEAD      # 只数本程的（编排者同期提交会混进来，看 subject 前缀「票 158 r1」）
$ git rev-parse --short HEAD                                            # 现量即最终树
```

本程写面（`git diff --name-status 8f9162f..HEAD` 实测只这三类）：
`internal/tools/bridge_scope_open_ticket158_test.go`（新用例）·
`.scratch/wisp/probes/154/gate-clauses.sh`（加 `pair()`＋G5 三腿，74 加 2 删）·
`docs/evidence/s1/158-gate-scope-blind-spot-r1.md`（本件）· `.scratch/wisp/probes/158/**`（只建不删）。）
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

---
---

# r2 追加节 · 票 158 r2 —— AC#4 零字节名册／AC#5 两包门禁／AC#6 双向对账 ＋ 一枚经批准的注释面扩权

- 程：**r2**。派单＝`.scratch/wisp/dispatches/2026-09-26-194x-impl-158-r2-redispatch-after-cancel.md`（19:4x 重派那一程）。
- **上面 r1 的第 0–7 节一字未改**（本件是续写不是改写；AC#1／AC#2／AC#3 三格 r1 已交并经编排者独立复算，本程不重做）。
- 本程只做三格：**§C(AC#4)／§D(AC#5)／§E(AC#6)** ＋ **§B 那一处经编排者批准的注释面扩权**。
- 上面 r1 §4 那节写的「AC#4／AC#5／AC#6＝未交·归 r2」现在由本程接手；那一节原文不抹。

## A. step-0 锚点（派单 §0 那四件的**原文输出**，本程一切行号／枚数以它为准）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-26 19:37 +0800

$ git rev-parse --short HEAD
bb08d1f

$ git status --porcelain                     # 全树
 M .scratch/wisp/probes/152/my152.py
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/

$ git log --oneline -3
bb08d1f 派单 19:4x: 票 158 r2 重派（前一程 18:33 被取消、零枚 commit）
bf046d1 ledger(A304)＋票158进度＋停车点＋158-r2 派单: r1 三格复算通过；一次由我发起的局部扩权
6f127f1 票 158 r1 交件：next= 里别再写死 commit 号，改成 8f9162f..HEAD 现量
```

工具版本现量（本程自己跑，不抄 r1 的）：`go version go1.27.1 windows/amd64` ·
`$(go env GOPATH)/bin/gofumpt.exe --version` ＝ **`v0.12.0 (go1.27.1)`**。

### A.1 与派单两条前提不符的地方（如实登记，没有一条让我改判据）

1. **派单 §0 说前一程「没有留下半成品」——不完全成立。** 现量：`git status --porcelain` 里那一枚
   `?? .scratch/wisp/probes/158/r2/` **确实存在**（11 枚文件，时间戳 18:29–18:32＝被取消那一程），
   里面是改前的门禁读数（`roster-before-*`／`gate-before-*`／`gofumpt-before.txt`／`vet-before.txt`）。
   **零枚 commit 那半句成立**（`git log` 里 18:33–19:35 之间没有任何提交）。
   ⇒ 处置：那些读数**全部不用**（锚点不是我的，且按临时件只建不删的规矩我不碰它们）；
   本程自己的读数一律落在**新建的** `.scratch/wisp/probes/158/r2b/`。派单那句「你从头做」我照做了。
2. **本程跑到一半树上又落了别人的提交**（派单 §1 预告的那一种）。现量序列：
   `bb08d1f`(19:36 我的 step-0) → `192ad56a`(19:40:36，编排者 ledger A305) → `ef34118f`(19:47:41，**本程那一枚注释 commit**) → `714f1698`(19:48:46，编排者 ledger A306＋停车点)。
   ⇒ 按派单 §1 的规矩**把改后那一跑重跑到新 tip**：§D 里 `cmd/wisp` 那一跑、`gofumpt`／`vet`／门禁脚本／`d22scan.sh` 全部在 **tip `714f169`** 上跑。
   ⇒ 并且机器核过这两枚编排者提交**与本程读数无关**：`git diff --name-only bb08d1f 192ad56 -- '*.go'` ＝ **空**（改前两跑的 Go 输入未变），
   `git diff --name-only 192ad56 714f169 -- ':!docs'` ＝ 只有 `internal/tools/bridge.go`（那是**我自己**那一枚）。

---

## B. 那一处注释面扩权（唯一一处，写面没有更宽）

| 项 | 内容 |
|---|---|
| 批准人 | **编排者**（不是我推断的、不是我从相邻授权外推的） |
| 批准出处 | `.scratch/wisp/dispatches/2026-09-26-194x-impl-158-r2-redispatch-after-cancel.md` §4；同一射程由编排者登记在票面 `Progress log [2026-09-26 18:2x]` 那一长条里 |
| 批准时刻 | 派单落盘＝2026-09-26 19:35（`ls` 现量的 mtime）；本程读到＝19:37 之后第一次读派单 |
| 批准的确切射程 | `internal/tools/bridge.go` 里**那一枚指针注释**（原 `:663-666` 那四行中的一句）**的措辞与所指路径**，一字一字算；不许动行为码／断言／阈值／golden；除这一处写面不变宽 |
| 「这不在票面原写面内」 | **明写：不在。** 票面 `:8` 的「地界」只给到 `internal/tools/**` 的**测试面**，这一枚注释在生产文件里 ⇒ 按原写面本程碰不到它 |
| 用的射程 | **用掉了，且只用这一处**（下方 commit `ef34118f`）。本程没有量到需要第二处 ⇒ 没有停手事件 |
| 提交形状 | 单独一枚 commit，标题含「注释面·编排者批准扩权」；`git commit -q -F <消息件> -- internal/tools/bridge.go`（显式 pathspec，消息件走 Write 工具，不经 shell 字符串） |

**改了什么**（`git diff --numstat`＝`7 3 internal/tools/bridge.go`，全仓唯一 dirty 的 `.go` 就是它）：
删掉那三行、换成七行。删的三行里过期的是两处——
① 只点 `clauses G1/G1b`，而 G5（r1 新增、射程正是**不过桥**那一族）没被点名；
② 落点写成 `docs/evidence/s1/154-host-id-never-closed-r1.md §2.3`，那一节装的是 **G1–G4 的文字**，
一条命令能跑的生成器其实是 `.scratch/wisp/probes/154/gate-clauses.sh`（G5 连那份文字都不在 §2.3 里）。
新句把两族各自管哪一条腿写清、并把**生成器路径**与**两份写-up**同时给出。

改后原文（现量 `sed -n '663,669p' internal/tools/bridge.go`）：

```
// Today no production code dispatches on the bridge except loop.go, so being
// such a caller means being the first one; the gates meant to ring when that
// happens are clauses of .scratch/wisp/probes/154/gate-clauses.sh - G1/G1b for
// this bridge leg, and G5 for a C25 taint scope opened with no paired close in
// the same file, i.e. the leg that never crosses this bridge at all (G5 is
// ticket 158's addition; §2 of 158-gate-scope-blind-spot-r1.md is its write-up,
// §2.3 of 154-host-id-never-closed-r1.md the write-up of G1-G4). Each is one
// `git grep` away, not a note in someone's head.
```

三件自证（票面 `:82` 要求的那三件，本程现量）：

```
$ git diff -U0 -- internal/tools/bridge.go | grep -E '^[+-][^+-]' | grep -cvE '^[+-][[:space:]]*//'
0                       # 非 // 开头的改动行＝0 枚（删 3 加 7 全是注释）
$ git diff --numstat -- internal/tools/bridge.go
7	3	internal/tools/bridge.go        # 删除列只落在这一枚文件
$ git status --porcelain -- '*.go'            # 改后、提交前
 M internal/tools/bridge.go                   # 全仓唯一脏的 .go
$ $(go env GOPATH)/bin/gofumpt.exe -l internal/tools/ cmd/wisp/    ; echo rc=$?
rc=0                                          # 空输出
$ go vet ./internal/tools/ ./cmd/wisp/        ; echo rc=$?
rc=0                                          # 空输出
```

**这一处改动对门本体的影响＝零枚读数**（这条不是推理，是现量，但它有一个**仪器坑**必须写下来）：
`gate-clauses.sh` 的锚点默认取 **`git rev-parse HEAD`**（脚本 `:10`），⇒ 注释还没提交时那把尺**看不见它**。
本程先在未提交状态下比了一次「改前／改后」，两遍只差时间戳一行——那是**假等价**。
提交后重跑（锚点 `ef34118f`）才量到真形状：判语一字未变（`UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)`／未成对＝1，
正控 0，负一负 5），名册里唯一变的是**行号** `internal/tools/bridge.go:691 → :695`（注释净加 4 行，把 `CloseScope` 那行推下去）。
复算尺＝`bash .scratch/wisp/probes/154/gate-clauses.sh` 两份读数原文
＝`.scratch/wisp/probes/158/r2b/gate-before-bb08d1f.txt` 与 `…/gate-after-comment-committed.txt`。

---

## D. 格 AC#5 —— 两包门禁：改前／改后四数 ＋ 名册两向 ＋ 三道工具

**档位：五子里四子已交、一子未裁（`sh scripts/d22scan.sh` 那一道，见 §D.5——它响在一枚**别人已提交**的前端码上，不在本程写面内）。**

### D.1 四数（逐包单跑；`internal/tools` 直接 `go test -v`，`cmd/wisp` 走 CI 同形入口）

尺＝仓内那把现成的（`tools/d22scan/runtests.sh:85-88`）：`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`，**顶层**口径。

| 跑 | 命令原文 | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | rc |
|---|---|---|---|---|---|---|
| 改前 · tools | `go test -count=1 -v ./internal/tools/` | **116** | **80** | **0** | **0** | 0（`ok … 13.937s`） |
| 改前 · cli | `bash scripts/wisp-cli-tests.sh` | **144** | **84** | **0** | **0** | 0 |
| 改后 · tools | `go test -count=1 -v ./internal/tools/`（tip `714f169`） | **116** | **80** | **0** | **0** | 0（`ok … 17.600s`） |
| 改后 · cli | `bash scripts/wisp-cli-tests.sh`（tip `714f169`） | **144** | **84** | **0** | **0** | 0 |
| 改后 · tools（仓内严格尺再跑一遍） | `sh tools/d22scan/runtests.sh ./internal/tools/ -count=1` | 116 | 80 | 0 | 0 | 0（`runtests.sh: OK`） |

改后 cli 那一跑打印的两行汇总（原文，含那枚 `-skip`）：

```
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$] top-level: PASS=84 FAIL=0 SKIP=0, === RUN=144, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=144  --- PASS=84  --- FAIL=0  --- SKIP=0
```

**分母为什么可以当分母用**：`internal/tools` 的 116/80 与 r1 §1.4 同值（那一格加了一枚用例后就是这个数）；
`cmd/wisp` 的 144/84 与票面 `Progress log [17:5x]` 那一条说的今天现量同值。两包的改前＝改后**逐枚同名册**，见 §D.2。

### D.2 名册两向 `comm` ＋ 差集逐枚归属

名册尺＝`^[[:space:]]*--- (PASS|FAIL|SKIP):` **全名（含子测试）**，去掉时长后 `LC_ALL=C sort`（与 r1 §1.2 同形）。
**并且这一发我差点读假，所以三条尺一起上**：`comm` 两向 **＋** `cmp -s`（字节级）＋ `md5sum`。

```
$ for pkg in tools cli; do cmp -s roster-$pkg-before.txt roster-$pkg-after.txt && echo "$pkg IDENTICAL"; LC_ALL=C comm -23 …; LC_ALL=C comm -13 …; done
tools: cmp => IDENTICAL   comm -23 => 空   comm -13 => 空   md5 两遍同值 31e18e28f1f03bd5ba387edf99c3ae17（各 116 行）
cli  : cmp => IDENTICAL   comm -23 => 空   comm -13 => 空   md5 两遍同值 c77310b70b17def6102f661ad543ca34（各 144 行）
```

⇒ **差集两向都是零枚**，归属因此是空的：本程唯一的非测试文件改动是 §B 那一处注释面，
它不改变任何一条读数的名字、状态或存在性。**这一格没有需要归名的差集，也没有需要解释的新红。**

（反面教材，本程自己撞的：第一次跑这个对比时我在 `LC_ALL=C sort` 之后**没带 `LC_ALL=C` 去跑 `comm`**，
MSYS 的 collation 与 C 排序不一致时 `comm` 会**安静地报「两向都空」**。同一形状在我量变异那一发上先露馅——
它把「一枚 PASS 翻成 FAIL」报成了零区别。三条尺一起上之后才敢说这是真等价。见 §G.2 的 I-2。）

**变异下的名册差集（用来证明上面那把尺吃得住这个形状，不是永远报空）**：
`.scratch/wisp/probes/158/r2b/roster-tools-before.txt` vs `…/roster-tools-mut.txt`（§E.1 那一发）——

```
只在改前（绿时）有：--- PASS: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
只在变异后有：      --- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
```
⇒ 同一条尺在真有一枚变脸时**点得出那一枚**，所以 §D.2 上面那个「两向皆空」不是尺瞎。

### D.3 `cmd/wisp` 的枚数账（把票面 `:42` 那一行判为**过期**，追加说明、不改写那一格）

票面 AC#5 第③条那句「`cmd/wisp` 声明 82 枚而本机只跑到 79（差的三枚在 `secret_dataroot_119b_test.go` 的 `//go:build !windows` 之下）」——
**那两个数（82／79）在本锚点已过期**，本程现量：

| 现量 | 命令 | 读数 |
|---|---|---|
| 源码里**声明**的顶层用例 | `grep -rn '^func Test' cmd/wisp/*_test.go \| sed 's/.*func \(Test[A-Za-z0-9_]*\).*/\1/' \| LC_ALL=C sort` | **87** 枚 |
| Windows 上**编译进去**的顶层用例 | `go test -list '.*' ./cmd/wisp/`（dll 已注入）`grep -c '^Test'` | **84** 枚 |
| CI 同形那一跑**真跑到**的顶层用例 | `^--- (PASS\|FAIL\|SKIP)` 计数 | **84** 枚 |
| 声明−编译 的差集（两向 `comm`） | `comm -23 declared-cli-all-src.txt list-…` | **3 枚**，全部在 `cmd/wisp/secret_dataroot_119b_test.go`，该文件 `:1` ＝ `//go:build !windows`（`TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119`／`TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119`／`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`）；反向差集（跑了但源码没声明）**空** |
| 那枚 `-skip` 在本包里挡住了几枚？ | 逐枚 `go test -list '^<名>$' ./cmd/wisp/` ＋ `grep -rl "func <名>(" cmd/wisp/` | **0 枚**：七枚名字在 `cmd/wisp` 源码里全部不存在（它们各自由 `portable-tests.sh` 的 ledger 在**自己的包**里作保，`internal/agent/approval`／`internal/memory`／`internal/proc`／`internal/audio`／`internal/models`／`internal/risk`），ledger 那枚 staleness 检（`scripts/portable-tests.sh:380-383`）按**行所命名的那个包**去 `go test -list` 验，不按当前 scope |

⇒ **那句的道理成立且本程复算到了**（结构性跑不到＝build tag 那三枚，别把「本机全绿」说成「全仓无影响」）；
**那句的数过期了**（今天 87／84／84，差仍是同一枚文件那 3 枚）。票面原句按派单 §2.3 一字不改，这一小段就是追加说明。

### D.4 「根本没跑到」的判据（派单 §2.1 那枚坑，本程自己量了一次正控与一发负控）

```
$ case "$PATH" in *sherpa-onnx*) echo YES;; *) echo NO;; esac
NO                                            # 这一跑的 PATH 里没有 dll 目录
$ go test -count=1 -v ./cmd/wisp/             # 不做 dll 注入
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.090s
FAIL                                          -> ^=== RUN 计数 = 0        rc=1
$ bash scripts/wisp-cli-tests.sh             # CI 同形（脚本自己 export PATH="$dll_dir:$PATH"）
^=== RUN 计数 = 144                            rc=0
```
⇒ 本程判「跑到了」只认 **`=== RUN` 枚数≠0** 这一件事，且两发都在本件里（原文＝`…/r2b/cli-negative-control-nodll.log` 与 `…/r2b/cli-after.log`）。
`0xc0000135` 本程不当结论用。dll 在位现量：`ls third_party/sherpa-onnx/` ＝ `onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`
（三枚，与 `deps.toml` 的 `[sherpa-onnx.dll.*]` 一致；脚本自己那行 `dll_dir=… (pinned: …)` 也在 `cli-after.log` 里，路径形＝shell 自己的 `/d/work/…`，正是脚本注释说的那个不能写成 `D:/…` 的坑）。

### D.5 三道工具

| 工具 | 命令原文 | 改前 | 改后（tip `714f169`） | 判定 |
|---|---|---|---|---|
| gofumpt | `$(go env GOPATH)/bin/gofumpt.exe --version`；`… -l . tools/d22scan tools/mockllm`；`… -l internal/tools/ cmd/wisp/` | 版本 `v0.12.0 (go1.27.1)`；全仓 **0 行** rc=0；窄 **0 行** rc=0 | 全仓 **0 行** rc=0；窄 **0 行** rc=0 | **绿**（版本自己 `--version` 现读并贴了） |
| go vet | `go vet ./internal/tools/ ./cmd/wisp/` | 空输出 rc=0 | 空输出 rc=0（字节数 0） | **绿** |
| D22 门 | `sh scripts/d22scan.sh` | **rc=1**（正控步红：`runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=28 FAIL=2 SKIP=0, === RUN=70`） | **rc=1**（同一枚红，两枚同名用例） | **未裁·受阻**，见下面整段 |

**`sh scripts/d22scan.sh` 这一子项：本程判「未裁」，不是「过」也不是「不过」。** 现量链条：

1. 红的两枚用例＝`TestScannerSelfScanOfRealRepoIsGreen`、`TestRealRepoLedgerIsHonest`，**都读运行时的工作树**。
   两枚的失败体里是**同一条** finding：
   ```
   scan_test.go:269: repo HEAD violates: frontend/src/components/harness/right-rail.tsx:108: [emoji] ban #8 glyph in scope frontend/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
   ```
2. 脚本 `set -eu` ⇒ 第 1 步（正控）红就退出，**第 2 步真扫没跑到**。本程因此把第 2 步单独跑了一遍拿读数：
   `cd tools/d22scan && go run . -root <仓根>` ⇒ **rc=1，1 finding，同一条**。
3. **这条 finding 不是本程造成的，也不是脏树造成的**：
   ① 那枚文件在 `git status --porcelain` 里**干净**（⇒ 第 108 行是 `HEAD` 里的内容），
   ② `git cat-file blob HEAD:frontend/src/components/harness/right-rail.tsx \| sed -n '108p' \| cat -A`
   ＝ `{ kind: "ok", text: "M-bM-^\M-^S built in 826ms" }`——`M-bM-^\M-^S`＝**U+2713**，在**字符串字面量**里（正是 AGENTS.md §1.2 那句「`✓`(U+2713) 不行」的形状），
   ③ 本程另外做了一件**独立口径**：`git archive HEAD \| tar -x -C $TMPDIR/…`（在**仓外**）再拿同一把尺扫那棵快照 ⇒ **同样 rc=1、同一条 finding**（`…/r2b/d22scan-headsnapshot.txt`）。
   ④ 来路：`git log -1 -- frontend/src/components/harness/right-rail.tsx` ＝ `5e23d99`（09-26 **17:46**，前端那一程的第 N 枚提交，早于本程 step-0 的 `bb08d1f` 19:36）。
4. ⇒ 本程**不许**修它：`frontend/**` 既在票面 AC#4 的冻结名单里、又「此刻正被别的会话写着」（派单 §3 明令不碰不还原不 commit），
   `tools/d22scan/**` 与 `allowlist.txt` 也全冻结；把它「豁免」掉正是票面禁的形状。**没有放宽任何断言、没有 `t.Skip`、没有改阈值。**

**同一把尺的另一半（`examined N` 非零）本程量到了，而且顺手把「别的会话的脏树」这件事量成了数**：

| scope | 工作树（含他人未提交改动） | `HEAD` 快照（仓外，纯净） |
|---|---|---|
| bans #1-5 `internal/` | 205 | 205 |
| bans #1-5 `cmd/` | 23 | 23 |
| ban #6 `frontend/` | 85 | 85 |
| ban #7 `internal/tools/` | 18 | 18 |
| ban #8 `design/` | **39** | **30** |
| ban #8 `frontend/` | 85 | 85 |
| ban #8 `internal/`（含注释与 `_test.go`） | 414 | 414 |
| ban #8 `cmd/` | 45 | 45 |

⇒ **八枚作用域 `examined N` 全部非零**（这一子项绿，且它才是「不是空仪器」的那半判据）；
唯一被他人未提交工作污染的是 `design/`（＋9，`git status --porcelain -- design` 现量 31 枚已跟踪改动 ＋ 61 枚未跟踪文件）。
**本程任何「零命中」宣称的口径**：`design/**` 与 `frontend/**` 一律**排除**在宣称之外（派单 §3），
本件里两族的数字只作为「仪器非空」的证据出现，不作为「谁也没写」的证据。
（快照那一跑额外打印了 `gitignore rules NOT APPLIED … not a git repository` ⇒ 快照口径的分母是**上界**，
它只会多扫不会少扫，所以「finding 相同」这一判语成立；这一句写在 `…/r2b/d22scan-headsnapshot.txt` 第 1 行。）

### D.6 AC#5 的**最小闭合集合**（只欠 §D.5 那一子项）

1. `frontend/**` 那一程（或任何拿到该写面的人）把 `right-rail.tsx:108` 字符串里的 U+2713 换成非禁字形（这是 D23／D29 的零 emoji 规矩，**不是阈值**）；
2. 然后重跑 `sh scripts/d22scan.sh`，要求 **rc=0** 且八枚作用域 `examined N` 仍全部非零（本件 §D.5 那张表就是改前基线）；
3. 若 owner 决定**不**改码而要走豁免，那必须给一条**可机读判据**并落在 `tools/d22scan/**`（本票与本程都写不到那儿）——那一支归 owner，不归 r2。
4. 这一子项闭合**不依赖**本程任何未交的东西：本程那两包的读数已经全绿，两件事没有耦合。

---

## E. 格 AC#6 —— 票面框 ↔ 本程格 双向对账

**档位：已交（两向都数过；未裁的两处明写，见 §E.4）。**

### E.0 对账尺本身（派单 §5：按节锚、不按行窗）

```
$ T=<票面路径>
$ wc -l < "$T"                                    -> 84 行（$ git hash-object -> sha 前缀 1cd6df70，工作树＝HEAD，干净）
$ grep -n '^## ' "$T"
11:## 这一格今天到底缺什么（三件，各自带出处）
22:## AC（判据一律要能答"这一发在未修码上响不响"）
45:## 本票**不**解决的事（登记，别顺手扩大）
52:## Acceptance（交件形状）
59:## Progress log (append-only, newest last)
$ grep -c '^- \[ \]' "$T"  -> 6      $ grep -c '^- \[x\]' "$T"  -> 0
$ grep -n '^- \[ \] \*\*AC#' "$T" -> 24:AC#1 / 28:AC#2 / 32:AC#3 / 34:AC#4 / 38:AC#5 / 43:AC#6
$ awk '/^## 本票\*\*不\*\*解决的事/{p=1;next} /^## /{p=0} p' "$T" | grep -c '^-'  -> 4 条
$ awk '/^## 这一格今天到底缺什么/{p=1;next} /^## /{p=0} p' "$T" | grep -cE '^[0-9]+\. '  -> 3 件
```
⇒ 本程**一枚 `sed -n 'A,Bp'` 都没用**；节与节之间是整节取的（另：本程开工第一件事是用 Read 工具读了票面全文 84 行，比任何行窗都全）。
分母：6 枚框、0 枚已勾（与 r1 §开工锚点那行说的「6 枚未勾／0 勾」同值，这一枚本程自己数过）。

### E.1 票面点名的三枚行号 ＋ 一句枚数，在**本锚点**全部复量

票面 `:57` 自己警告过那三枚行号是 14:2x 在 `5d286d5` 前后读的。现量（tip `8225cdc`）：

| 票面引的 | 本程现量命令 | 读数 | 漂没漂 |
|---|---|---|---|
| `bridge.go:559` | `sed -n '559p' internal/tools/bridge.go` | `\tb.OpenTask(dec.TaskID)` | **未漂**（本程那枚注释加在 `:663` 之后，不影响它） |
| `panel_assets.go:232` | `sed -n '232p' cmd/wisp/panel_assets.go` | `prov.OpenScope(taintSourceScopeID)` | **未漂**（今天仍是 G5 唯一点名的开方） |
| `task_scope_close_151_test.go:109/:113` | `grep -n 't\.Errorf\|t\.Fatal' cmd/wisp/task_scope_close_151_test.go` | `t.Errorf` 在 **`:62 :103 :109 :113`** 四枚 | **未漂**（票面只点后两枚，前两枚 r1 §0-P3 已自纠；两枚用例名＝`TestCompositionRootClosesTheLoopTasksTaintScope`／`TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit`） |
| AC#5 第③条「`cmd/wisp` 声明 82、本机跑到 79」 | 见 §D.3 的三行现量 | 今天 **87 声明／84 编译／84 跑到**，差集 3 枚仍全在 `secret_dataroot_119b_test.go` | **数过期、道理成立** ⇒ §D.3 已按派单 §2.3 追加说明，票面原句一字未改 |

### E.2 正向：一枚框 → 哪一程交的 → **这一发在未修码上响不响**（本程自己的读数）

| 框 | 交它的那一程 | 本程是否重做 | 「未修码上响不响」＝本程现量 | 判语 |
|---|---|---|---|---|
| **AC#1** 先复现变异再装本包守卫 | **r1**（上面 §1，编排者已独立复算） | **不重做**，只做「今天还响不响」的回归复量 | 变异**今日重造**（从当轮 `internal/tools/bridge.go` 拷一份、`sed -i '559d'`，`diff` 只差 `559d558 < b.OpenTask(dec.TaskID)` 一行）＋ `-overlay` 单跑（未与任何 `-cover` 同用）＝ **116／79／1／0、rc=1**，红名 `TestSensitiveReadAcrossTheBridgeOpensItsTaskScope`（原文 `probes/158/r2b/tools-mut-r2.log:55`）。**并先钉住「overlay 真被编译器吃」**：把 `:559` 改名成 `b.OpenTask_CANARY_R2` 的那发 `go build -overlay=…` 报 `.\.scratch\wisp\probes\158\r2b\bridge.canary-r2.go:559:4: … undefined`（rc=1）、不带 overlay 的 `go build ./internal/tools/` rc=0 ⇒ 那枚红**确实**是这发变异造出来的。**⇒ 响** | **已交·今日仍响** |
| **AC#2** 给门加 G5（射程＝开/合成对） | **r1**（上面 §2） | 不重做，只复量三腿 | 主尺 `UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)`、未成对＝**1** ⇒ **响**；正控未成对＝0（且 `git grep rc=0`、名册非空）；负一负＝5。锚点 `192ad56`（本程改前）与 `ef34118`（本程改后）两跑**判语一字未变**，唯一差别是名册里 `bridge.go:691 → :695` 那枚**行号位移**（§B 末解释过为什么） | **已交·今日仍响** |
| **AC#3** 「有界」要么给形状要么让它响（r1 选**支 B**） | **r1**（上面 §3） | 不重做、不改选边 | 支 B 的判据就是 G5 那一枚的读数，本程复量今天仍＝**1**（未熄火）⇒ 「响着并登记」这一支今天**仍在履约**。r1 给的两条形之路（`panel_assets.go` 拿到收尾／owner 给可机读谓词）本程**一条都没代做** | **已交·仍响＝登记仍生效** |
| **AC#4** 契约轴零字节 | **本程** | — | 这一枚框的形状不是「响不响」而是「**我这几枚 commit 有没有往冻结面上落一个字节**」⇒ 读数在 **§C**（名册尺原文都在那儿） | **本程交**（见 §C） |
| **AC#5** 两包门禁＋三道工具 | **本程** | — | 五子项四绿；**一子项未裁**＝`sh scripts/d22scan.sh` rc=1。它**正是**「在未修码上响」那一发：响在一个**本程没碰过**的、别人**已提交**的字形上（§D.5 三口径钉住了「不是脏树造成的」），本程不许它变绿 | **本程部分交**（＋未裁一处，见 §E.4） |
| **AC#6** 双向对账 | **本程** | — | 本节自身 | **本程交**（本节） |

### E.3 反向：本程动过的每一类东西 → 它属于哪枚框（孤儿工作清点）

| 本程动过／产出的东西 | 归到哪枚框 | 是否票面要求过 |
|---|---|---|
| `docs/evidence/s1/158-gate-scope-blind-spot-r1.md` 的 r2 追加节（§A–§G） | AC#4／AC#5／AC#6 ＋ Acceptance 节那四行交件形状 | 是（票面 `:52-56`） |
| `.scratch/wisp/probes/158/r2b/**`（37 枚读数件，**只建不删**） | 上面三格的原始输出落点 | 是（Acceptance 节要求「现量命令原文＋输出」） |
| `internal/tools/bridge.go` 那一枚**注释**（commit `ef34118f`） | **哪一枚框都不是**——它闭合的是 r1 §5-1「下一位读者会被误导」那一格 | **不在票面 6 枚框里**，来自派单 §4 的单次批准 ⇒ 请编排者在验收表里**单列一行**，别并进 AC#2 的「已交」（本程不替它安框） |
| 门禁／变异／普查的**跑**本身 | AC#5／AC#6 | 是 |
| **本程没有**新增 `internal/tools/**` 用例 | —— | 本程**不该**动 AC#1 的守卫（那是 r1 的格），也确实没动：`git log --no-walk --pretty=tformat: --name-only <本程号集合>` 里 `internal/tools/**` 只出现 `bridge.go` 一枚（注释面） |
| **本程没有**改门本体 | —— | `.scratch/wisp/probes/154/**` 在本程 commit 里**零枚**（G5 是 r1 加的，本程只跑不改） |

**票面「本票不解决的事」四条的逐项自查**（这一向也是反向对账的一部分：本程有没有顺手越界）：

| 那条不解决的事 | 本程有没有碰 | 机器核 |
|---|---|---|
| 不答 `Q-56` | **没碰** | 本程写的两件套（证据件 r2 追加节＋那枚注释）里新增行含 `Q-56` 的枚数＝**0／0**（`git show <号> --format='' -- <那两枚文件> \| grep -c '^+.*Q-56'`）。**口径要写清**：不拆文件时整棵 commit 命中 **5** 枚，逐枚全是**抄来的他程文本**——`bridge.go` 的两份变异副本各 1 枚，门本体那句节标题 `## G3 Q-56 那一支落地…` 被抄进 3 份读数件。不拆文件，这把尺会把"抄了什么"读成"答了什么" |
| 不裁 `panel_assets.go:232` 的后果链 | **没碰** | 本程只在 §D.5／§E.2 里把它当 G5 的**点名读数**引用；`cmd/wisp/**` 与 `internal/risk/**` 在本程 commit 里**零枚**（§C 的名册尺） |
| 不把门接进 CI | **没碰** | 本程 commit 里 `.github/**`、`tools/d22scan/**`、`scripts/**` 三类**零枚**（`git log --no-walk --pretty=tformat: --name-only <号集合> \| grep -E '^(\.github\|tools/d22scan\|scripts)/'` ＝ 空） |
| 不重开票 154 已结的六格 | **没碰** | 本程 commit 里任何 `154-*`／`151-*`／`156-*` 命名的证据件或票面**零枚**（同一条尺换 grep ＝ 空） |

### E.4 本程**未裁**的格子 ＋ 各自的最小闭合集合

1. **AC#5 的 `sh scripts/d22scan.sh` 子项＝未裁（受阻，不是不过）。**
   最小闭合集合＝§D.6 那三步，一句话版：`frontend/**` 的 owner 把 `right-rail.tsx:108` 字符串里那枚 U+2713 换成非禁字形 ⇒ 重跑该脚本要求 rc=0 且八枚作用域 `examined N` 仍全非零（本件 §D.5 那张表是它的改前基线）。
   **本程为什么自己闭不掉**：那条 finding 的两个可能落点（`frontend/**`、`tools/d22scan/**`＋`allowlist.txt`）全在票面 AC#4 的冻结名单里，且 `frontend/**` 此刻正被别的会话写着（派单 §3）。
   放宽断言／`t.Skip`／把字形抄进白名单都不在选项里（票面 AC#3 禁形状②③、派单 §6）。
2. **AC#1 的「发①：未修码＋未加守卫 ⇒ 零枚红」那一腿本程不重判。**
   今天要在未修码上复现它，必须**摘掉 r1 那枚已入库的守卫用例**（`internal/tools/bridge_scope_open_ticket158_test.go`）——那是删别人的证据，越界。
   本程交的是**相邻的那一发**：同一枚变异在今日码上只让**一枚**用例变脸（§D.2 末那两向差集），
   即「本包对这发变异的响应**只有** r1 那枚守卫」——那正是 AC#1 当初要的"唯一的守卫长在本包里"的反向确认。
   最小闭合集合（如果编排者要那一腿也重跑）：`git stash` 之外没有任何办法在共享树里临时删别人用例 ⇒ 只能由**非实现者验收程**在**快照**里跑（`git archive <r1 前一枚 commit>` ＋ 不带那枚文件的 `-overlay`），本程已把这条路留在 §F.6。
3. **票面 6 枚框本程一枚都没勾**（Acceptance 节末句：框由编排者按非实现者验收表定）。本件里出现的"已交"字样都是**程的自报档位**，不是勾。

---

## C. 格 AC#4 —— 契约轴零字节：按**本程 commit 号集合**算的名册普查

**档位：已交。** 本节**排在 §D／§E 之后**是故意的：AC#4 的名册必须覆盖那两格的 commit，
所以它只能最后量。编号按 AC 号，不按落盘顺序。

生成器＝`.scratch/wisp/probes/158/r2b/zero-byte-census.sh`（只读，不写仓里任何东西），
原始输出＝同目录 `zero-byte-census.txt`。本节所有数字都从那份文件抄，命令原文也在那份文件里。

### C.1 名册尺（主尺）＝按具体 commit 号集合，逐枚 `--no-walk`

派单 §3 的三条禁令本程都照做了，而且**把那条"没有 --no-walk 会怎样"自己量了一遍**：

```
$ bash .scratch/wisp/probes/158/r2b/zero-byte-census.sh        # 锚点 bb08d1f，跑时 tip=d041675（20:09:19，编排者的 ledger A309）
## 1. 本程 commit 号集合（按 subject 前缀选，不是按时间窗）
$ git log --format='%H' bb08d1f..HEAD --grep='^票 158 r2'
  ee41e63 20:03:49 票 158 r2 · 格 AC#6：票面框↔本程格 双向对账（…）
  8225cdc 19:59:38 票 158 r2 · 格 AC#5：两包门禁改前改后四数＋名册两向零差＋三道工具（…）
  ef34118 19:47:41 票 158 r2 · 注释面·编排者批准扩权：bridge.go 的门指路句补上 G5 与可跑的生成器路径
# 枚数=3
# 反向自证：同一区间里【不是】本程的提交（选规则没有多吞/漏吞）：
  d041675 20:09:19 ledger(A309): Cordis 交件——D21 那句"不依赖 cordis-rs"在真实上游里 0 命中；…
  2953203 20:07:xx ledger(A308)＋新增 Q-58: W1 把我们自己的工具面底盘量出来了——声明 36/实现 6/注册 5+1/可达 5
  e7f7741 20:04:07 ledger(A307): Pi 那枚 edit 的机制逐条验到原文＋PLAN.md D21 三处过期数并入 Q-45＋调研封批
  714f169 19:48:46 ledger(A306)＋停车点: owner 两句话到账＋再派四路（…）
  192ad56 19:40:36 ledger(A305): 外部 harness 十路对标交完＋我撤回两条已讲出去的话
# 每一枚的 --no-walk 自证（都是 commit，不是 tree/blob）：  commit ×3   rc=0

## 2. 名册尺（主尺）
$ git log --no-walk --pretty=tformat: --name-only <三枚号> | sed '/^$/d' | LC_ALL=C sort | uniq -c
# 名册去重后枚数=40
# 派单点名的那枚坑，本程自己量一遍（同三枚号，只差 --no-walk）：
  带 --no-walk   = 41 行
  不带 --no-walk = 4870 行          # ← 118 倍。少了这个旗标，整仓历史都会被算成"本程动过的"
```
⇒ 名册 40 枚去重路径 ＝ 38 枚 `.scratch/wisp/probes/158/r2b/**` ＋ 本证据件 ＋ `internal/tools/bridge.go`。**除这两类写面与那一处授权注释外，本程什么也没碰。**

### C.2 逐支账（把票面 AC#4 那张清单逐枚列，每支再套一条 pathspec）

28 支全零；三支"应有命中"的单列在下面。**"该面已跟踪文件数"那一列是仪器有没有落下去的证据**（非零＝尺扫得到这一支，不是空转）：

| 冻结面（票面 AC#4／派单 §3 点名的） | 越界枚数 | 该面已跟踪文件数 |
|---|---|---|
| `docs/PLAN.md` | **0** | 1 |
| `docs/specs/**` | **0** | 14 |
| `internal/risk/**`（`OpenScope`/`CloseScope` 的定义本体在这儿） | **0** | 37 |
| `internal/panel/**` | **0** | 20 |
| `internal/agent/**`（含 `approval/**`） | **0** | 58 |
| `internal/observe/**` | **0** | 25 |
| `thresholds.go`＝`internal/observe/thresholds.go` | **0** | 1 |
| golden：`internal/llm/golden/**` | **0** | 3 |
| golden：任何路径含 `golden`（`:(glob)**/*golden*`） | **0** | 6 |
| golden：任何 `**/testdata/golden/**`（`.sse` 那一大堆） | **0** | 52 |
| `allowlist.txt`＝`tools/d22scan/allowlist.txt` | **0** | 1 |
| `scripts/slo-check.ps1` | **0** | 1 |
| `tools/d22scan/**` | **0** | 6 |
| `frontend/**` | **0** | 85 |
| `design/**` | **0** | 30 |
| 票 151 证据件 ×2／票 154 证据件 ×2／票 156 证据件 ×3 | **0** | 各 1 |
| 票 151／154／156 的**票面** | **0** | 各 1 |
| **另加三支不在 AC#4 名单、但同样不是本程写面的**：本票票面（`:158-*`）／台账＋停车点（`docs/reports`）／门本体（`probes/154`，r1 写面） | **0／0／0** | 1／15／31 |

| 写面（**故意有命中**，藏起来反而可疑） | 越界枚数 | 归属 |
|---|---|---|
| `docs/evidence/s1/158-gate-scope-blind-spot-r1.md` | 1 | 本票证据件＝票面地界给到的写面 |
| `.scratch/wisp/probes/158/**` | 38 | 本程探针件（只建不删）。**注意口径**：同目录下还躺着 `probes/158/r2/`（被取消那一程的临时件，`git status` 里是 `??`），本程**没 stage、没读它内容、也没 commit 它**——那一支 38 枚全部在 `r2b/` 下 |
| `internal/tools/bridge.go` | 1 | §B 那一处**经批准的注释面**；`git show ef34118 --numstat` ＝ 加 7 删 3，改动行 100% 以 `//` 开头（§B 已量） |

### C.3 第二口径（字节级）＋ **那把尺的正控**（不然"零命中"跟"尺没通电"长一个样）

```
## 5. 把 numstat 的文件名列打上一串正则（与 §8 的正控共用同一串，不会两把尺漂移）
$ git log --no-walk --pretty=tformat: --numstat <三枚号> | cut -f3 | grep -E '<FIRE>'
  HIT internal/tools/bridge.go        # 只剩授权那一枚；其余冻结面全零
## 8. 同一串正则打在【已知动过冻结面】的三枚真提交上（必须响）
  12181923 命中=2   risk(票 141 具名解冻 ③): 审批卡 R7 文案那枚 `≥` 换成 ASCII…
      FIRE internal/risk/assessor_test.go
      FIRE internal/risk/rules_scale.go
  00bbb76f 命中=3   slo(66): read the D32 CPU row out of the measured tree (AC#3)
      FIRE internal/observe/observer_cost_test.go
      FIRE internal/observe/sampler.go
      FIRE internal/observe/thresholds.go        ← 连 thresholds.go 本尊都点得出
  cbbbdf85 命中=1   docs(agent): correct the spill fixture's scenario header (ticket 10 MINOR-8)
      FIRE internal/agent/testdata/golden/spill-tool.sse   ← golden 也点得出
```
⇒ 三枚全响 ⇒ §5 那个"只剩 `bridge.go`"是**码上真的没有**，不是一把产不出读数的装饰尺。
（这一发的来历也要登记：§5 的第一稿在正则里用了**字面 TAB**，而"TAB 有没有活着穿过 Write 工具"不该让下一位读者赌 ⇒ 改成 `cut -f3` 之后再打正则。
且**正控第一次就抓到我抄错的一枚号**（写成 `000bbb7f`，`git` 直接 `unknown revision`、那一行报 命中=0）：如果我当时忽略"三枚里有一枚是 0"，交出去的就是一枚半死的正控。）

### C.4 第三把尺：区间 diff —— 派单禁它当普查尺，本程把它当成**为什么禁**的证据

```
$ git diff --name-only bb08d1f HEAD | grep -v '^\.scratch/wisp/probes/158/r2b/'
  docs/evidence/s1/158-gate-scope-blind-spot-r1.md     ← 本程
  internal/tools/bridge.go                             ← 本程（授权注释面）
  docs/reports/HANDOVER.md                             ← 【不是本程】
  docs/reports/pending-and-issues.md                   ← 【不是本程】
```
逐枚归名（这五枚都不在 §1 的本程集合里）：`d041675`／`2953203`／`e7f7741`／`714f169`／`192ad56`，全是编排者的 ledger／停车点提交。
⇒ **谁今天造的假越界，这就是现场**：区间尺会把别人写成我的。普查必须按 commit 号集合算。

**另一枚"差点假越界"的形状，值得单独留一句**：`e7f7741` 的 subject 里写着「**PLAN.md** D21 三处过期数并入 Q-45」——
只扫 subject 的人会以为 `docs/PLAN.md` 被改过。现量两把尺都说没有：
§4 里 `docs/PLAN.md` 越界＝**0**（该面 1 枚已跟踪文件，尺扫得到），
且**连区间尺**在 `bb08d1f..HEAD` 全表里也没出现 `docs/PLAN.md`（上面那四支就是全部）。⇒ "过期数"是**并入台账 Q-45**，不是动了 PLAN.md 一字。

**`frontend/**`／`design/**` 的口径**（派单 §3 明令）：本程**没碰**（§4 两支越界＝0），
且本件里所有"零命中"宣称**都不把这两族算作证据**——它们的工作树此刻是别人的（`git status` 现量 `design/` 31 枚已跟踪改动＋61 枚未跟踪）。
它们在 §D.5 出现的唯一身份是"d22scan 仪器非空"的分母，见那里那张两列对照表（工作树 39 枚 vs `HEAD` 快照 30 枚，差＝别程未提交的工作）。

### C.5 本程 commit 之后的最终-tip 复跑（AC#4 那一枚 commit 自己不能量自己）

| 跑 | tip | 读数 |
|---|---|---|
| `go test -count=1 -v ./internal/tools/` | `d041675` | **116／80／0／0 rc=0**（`…/r2b/tools-finaltip.log`） |
| `bash scripts/wisp-cli-tests.sh` | `d041675` | **144／84／0／0 rc=0**（`…/r2b/cli-finaltip.log`，脚本自己那行 `portable-tests.sh: four numbers` 原文在内） |
| `$(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm` | `d041675` | **0 行 rc=0**（`…/r2b/gofumpt-full-finaltip.txt`；这一发还顺手把"新落库的两份 `bridge.*-r2.go` 探针副本会不会被格式化尺点中"一起量了：**没有**） |
| `go vet ./internal/tools/ ./cmd/wisp/` | `d041675` | **空 rc=0**（`…/r2b/vet-finaltip.txt`，0 字节） |
| `cd tools/d22scan && go run . -root <仓根>`（真扫那一步；`sh scripts/d22scan.sh` 仍被 §D.5 那枚正控步挡在 rc=1） | `d041675` | **rc=1、finding 枚数仍是 1、且还是同一枚 `right-rail.tsx:108`** ⇒ **本程的三格交件没有新增任何被禁字形**（`…/r2b/d22scan-finaltip.txt`） |

机器核过"别人同期提交没动到我这两包的输入"：`git diff --name-only 714f169 HEAD -- cmd/wisp internal/tools` ＝ **空**，
`-- '*.go'` 里除本程自己的探针副本外也**没有**第二枚。
⇒ §D 的改后读数在最终 tip 上仍成立；AC#4 自己那一枚 commit 只加证据件＋探针件（Go 字节 0 枚），
所以它**不可能**改变上面两个四数——这一句不是断言，是一条可复算的预测：**下一位把 `bash .scratch/wisp/probes/158/r2b/zero-byte-census.sh` 原样再跑一次**，
预期读数＝枚数从 **3 变 4**、名册去重从 **40 变多**但多出来的只可能是 `probes/158/r2b/**` 与本证据件两支、
`§5` 那一发仍然只剩 `bridge.go`、28 支冻结面仍全 0。**任何别的数字变化＝有人的提交混进了这个名册。**

### C.6 AC#4 判语

**已交：契约轴零字节成立。** 三条独立口径（逐枚名册尺／字节级 numstat 尺（带正控）／工作树侧）同判，
唯一命中在**经批准的注释面**那一枚文件、且那一枚的改动行 100% 是 `//`。
外加一条派单要求的自查：本程**没有**新增任何 `DEFERRED(D-xx)` 代码标记（新增行里的 2 处 `DEFERRED` 逐枚归名后全在 `bridge.go` 的两份**探针副本**里，
是从生产注释里整份抄来的），也**没有**新增或改动任何 `t.Skip`（新增行里的 2 处 `t.Skip` 是本件**说"没有 t.Skip"的那两句话自己**被尺扫到了——
同一条尺不拆文件就会把"否认"读成"做了"，与 §E.3 里 `Q-56` 那发是同一类口径事故）。凭据值零抄录，见 §G 末。

---

## F. 本程**没测／没做**的东西（按「漏了它谁会先被骗」排序）

1. **AC#5 的 `sh scripts/d22scan.sh` 子项＝未裁**（§D.5／§D.6）。⇒ 先被骗的是**任何把"r2 门禁交了"读成"门禁全绿"的下一程**：
   那一跑今天 rc=1，红在**别程已提交**的一枚字形上。本件的四绿一未裁就是这个意思，别并成五绿。
2. **本程只判了那枚字形"被禁"，没判它"该换成什么"**（`frontend/src/components/harness/right-rail.tsx:108` 的 U+2713）。
   ⇒ 会被骗的是**前端那一程的 owner**：本件给的是"扫描器说禁＋来路是哪枚 commit"，不是文案建议。
3. **`cmd/wisp` 的 84 枚只在「windows ＋ dll 在位」这一形量到**。`!windows` 那三枚（`secret_dataroot_119b_test.go`）本程跑不到、也没跑（§D.3）。
   ⇒ 会被骗的是把"cmd/wisp 全绿"读成"跨平台全绿"的人——票面 AC#5 第③条讲的正是这一形，本程只是把它的**数**换成今天的、把**归属**机器核到了那枚 tag。
4. **AC#1 那一腿「未修码＋未加守卫 ⇒ 零枚红」本程没重跑**（重跑得摘掉 r1 已入库的守卫用例＝删别人的证据）。
   ⇒ 会被骗的是以为"那一腿今天仍被本件守着"的人；本件今天重跑的是**另外那一发**（今日重造变异 ⇒ 116/79/1/0）。见 §E.4-2。
5. **`-cover*` 一枚没跑**（派单 §2.2／票面 AC#5②：与 `-overlay` 合用时 overlay 被静默忽略）。⇒ 本件**不做任何覆盖率宣称**；
   谁要覆盖率必须**另起一跑**，不能拿本件任何一发当覆盖率证据。
6. **只跑了 `internal/tools` 与 `cmd/wisp` 两包**：`go test ./...` 本程一枚没跑，全仓其余包在最终 tip 上**状态未知**（别的程在写）。
   ⇒ 会被骗的是把 AC#5 读成"全仓门禁"的人。
7. **G1–G4 只复算到"本程改动前后一字未变"**，没重判它们的射程该不该更宽（那是票 154 的账，票面"不解决的事"第④条）。
   同族里 `OpenTask/CloseTask`（G2）与 `DisposalScope`/`Defer` 那一族的成对性本程**没普查**（r1 §5-5 已登记，本程没补）。
8. **桥的 `b.scopes` 与 `prov.scopes` 两套账的一致性没裁**（票 154 的 U-5；r1 也未做）。本程只是把那枚指针注释修对了，
   ⇒ **注释说得对不对，本程没验**（"G5 守着不过桥那一腿"是 r1 §2 的读数＋脚本原文，不是本程新证的命题）。
9. **无凭据语料**：`%APPDATA%\wisp` 本程**未扫** ⇒ 既不能写"扫过很干净"，也没扫过。全程零抄录凭据值（只写变量名/文件名，见 §G 末）。
10. **争用只量了一枚形状**：`Get-Process` 命中 `Runner.Listener`（pid 3952）、`Win32_Processor.LoadPercentage`＝**41 与 29 两次读数**
    ⇒ 本机 self-hosted runner 同机在跑，**未排除争用**。本件所有结论都是**枚数／名册**判定，不含任何时间阈值 ⇒ 争用不改变判语，
    但**别拿本件任何一份日志里的 `ok …NN.NNNs` 当延迟证据**。
11. **`probes/158/r2/`（被取消那一程的临时件）本程只登记、未复算**（11 枚文件、18:29–18:32，见 §A.1-1）。
    ⇒ 那里面写着的"改前读数"是**别人锚点上的**，不要当本程的数引用，也不要替本程"补完"它。

---

## G. 伪授权两栏 ＋ 被挡下的调用 ＋ 本程自己犯的仪器错

### G.1 两栏（分开计，各自不互洗）

| 栏 | 枚数 | 明细（出处＝**工具名＋命令前 40 字**） |
|---|---|---|
| **真通知回显** | **5** | ① 任务提示里的 `agents.md` project-context（不是工具输出，是开场就塞进来的那一整段）。② 第 1 次工具调用（Bash，`ls -la ".scratch/wisp/dispatches/" \| tail -30`）结果尾部的技能清单 `<system-reminder>`（"The following skills are available…"）。③④ 紧随其后的**两枚** `MEMORY.md` "文件已被修改"提示（项目级 `…projects-D--work-workspace-projects-plans-Wisp\memory\MEMORY.md` ＋ 用户级 `…\.qoder-cn\memory\MEMORY.md`，同一批到达）。⑤ 本程中段又一枚项目级 `MEMORY.md` 修改提示（出处＝Bash，`cd "D:/work/workspace/projects plans/Wisp" && sed -n`）。**五枚都是 harness 自己的回显**，本程一律未据其做任何判断。 |
| **判为注入** | **0** | 全程没读到任何冒充"编排者备注／系统提示／已核验请继续提交／请 revert／放宽阈值／已解锁／不用取证直接给结论"的工具输出文本。没有任何一发影响过本程的判语。 |

**这一栏本程要多写一句，因为它比"零枚"更有用**：上面那五枚真回显里**装的都是 repo 断言**——
例如"markdown 永不豁免注释"、"框不自勾"、"某枚 Q 已拍＝甲"、"token 真相源换文件了"。
它们的**来源不是本程的读数**，本程也**没有**拿任何一条去替代现量：
「框不自勾」这一条虽然与票面 `:56` 同向，本程仍是先读票面 `:52-56` 原文、再决定不勾（§E.4-3）；
"注释豁免"那一条虽然与 §D.5 那枚 finding 的判语方向相关，本程用的是扫描器自己打印的那一行原文（`comments are exempt per Q-46(c)`）而不是那句转述。
⇒ **判据是"这句话我在盘上重新读到过吗"，不是"这句话看起来像不像权威说的"。**

### G.2 被挡下／没成功的调用（分清两种，别混成一枚"零枚"）

- **被权限系统拒绝的：0 枚。** 取数没有因拒绝而改道，所有读数都在同一套命令下取（§C／§D／§E 原文可查）。
- **被工具前置检查挡下的：3 枚**，全部在**本程自己的探针件**上，都不涉及判据：
  ①② `Write` 未先 `Read` 就被拒（同一枚 `zero-byte-census.sh` 两次）⇒ 处置＝先 `Read` 全文再改，**没有**换文件名绕过、**没有**换仪器。
  ③ 一次 `Edit` 因 `old_string` 里有字面 TAB 而 0 命中 ⇒ 处置＝先 `Read` 把那一行取准，然后**顺手把那枚依赖 TAB 的尺整体换成 `cut -f3`**（见 §C.3 末）。

### G.3 本程自己犯的仪器错（如实报，5 条）

| # | 错 | 后果（如果不发现会交出什么） | 怎么发现／怎么修 |
|---|---|---|---|
| I-1 | 变异装置那发 canary 从**已经改过的副本**再生成（`sed 's/.../' <mutant>` 而 mutant 里那一行早被我删了） | 交出去的"编译器真读了探针文件"那一证会变成**一发沉默的 rc=0**——正是 §E.2 需要它证明的东西，等于假证 | `grep -c OpenTask_CANARY_R2 <canary>` ＝ **0** 才发现；改从**真源** `internal/tools/bridge.go` 生成，重跑 ⇒ 报错点名 `.\.scratch\…\bridge.canary-r2.go:559:4` rc=1，且不带 overlay 的 build rc=0 |
| I-2 | 名册比对：`LC_ALL=C sort` 之后**没带 `LC_ALL=C` 跑 `comm`** | 会把"一枚 PASS 翻成 FAIL"报成**两向皆空**＝把有差集的名册交成"零差" | 同一形状在量变异时先露馅（变异明明 FAIL=1，`comm` 却报空）。修法＝**三条尺一起上**：`comm` 两向 ＋ `cmp -s` 字节级 ＋ `md5sum`，并把这条写进 §D.2 当反面教材 |
| I-3 | 注释改完**还没 commit** 就去跑门本体比"改前／改后" | 两遍只差时间戳一行 ⇒ 交出一条**假等价**（门默认取 `git rev-parse HEAD`，看不见未提交改动） | 从输出里那枚锚点 sha 与 `CloseScope` 行号都没动**反常地一致**起疑；提交后重跑 ⇒ 真读数：判语一字未变、名册里 `bridge.go:691 → :695` 位移（§B 末） |
| I-4 | 普查尺 §5 的正控里抄错一枚 commit 号（`000bbb7f`） | 正控三枚里会有一枚报 命中=0，整节变成"两响一不响"——看着像尺不稳 | 重跑时那行 `fatal: ambiguous argument '000bbb7f'` 直接响在输出里；换成 `00bbb76f` 后三枚全命中（2/3/1，含 `thresholds.go` 本尊与一枚 golden `.sse`） |
| I-5 | 普查尺 §7 第一稿对区间表用了 `head -60` | 带截断的输出当全表 ⇒ 正是编排者今天在这张票上栽过的那发（把"负一负 5 枚"读成"主尺 5 枚"） | 用之前先改成 `grep -v` 去噪＋**逐枚归名**（§C.4 那张五枚表），不保留任何截断读法 |

**凭据值零抄录**：本件与全部探针件里没有出现过任何凭据值——本程读过的凭据相关面只有
`deps.toml` 里 `[sherpa-onnx.dll.*]` 那几枚**文件名**（dll 名）与 `third_party/sherpa-onnx/` 的目录列表；
未打开 `wisp secret` 的任何存储面、未扫 `%APPDATA%\wisp`（§F.9）、未抄任何变量值。

---

## H. 本程交件形状（一眼账）＋ `next=`

| 格 | 档位 | 落在哪一节 |
|---|---|---|
| AC#4 契约轴零字节 | **已交**（三口径同判，唯一命中＝经批准的注释面；正控三枚全响） | §C |
| AC#5 两包门禁 | **四子项交 ＋ 一子项未裁**（`sh scripts/d22scan.sh` rc=1，响在别程已提交的一枚字形上） | §D（未裁的最小闭合集合＝§D.6 ＋ §E.4-1） |
| AC#6 双向对账 | **已交**（6 枚框两向都数过；两处未裁明写） | §E |
| 那一处注释面扩权 | **用掉且只用一处**；批准人／时刻／确切射程／"这不在票面原写面内"都登记了；**没量到第二处 ⇒ 没有停手事件** | §B |
| r1 的 AC#1／AC#2／AC#3 | **本程未重做、未改写**（上面第 0–7 节一字未动）；只做了"今日仍响"的回归复量 | §E.2 |

本程 commit 枚数＝**4**（一枚注释面 ＋ 三格各一枚），全部带显式 pathspec、全部**只 commit 未 push**；
`git diff --cached --name-only` 每枚提交前都看过，出现的都是 `docs/evidence/s1/158-…md` 与 `.scratch/wisp/probes/158/r2b/**` 两支，
**没有任何一支不属于本程写面**。别人的脏文件（`design/**` 那一堆、`probes/152/my152.py`、
`probes/158/r2/`、`docs/evidence/s1/152-…md` 那 3 行未提交自校）本程**未 stage、未 commit、未还原、未补完**。

`next=` **别信本件里任何一枚写死的 commit 号**（包括这一句旁边的那些——共享树会自己往前走，本程今天就撞到三次）。
接手时现量这四件：

```
$ date "+%Y-%m-%d %H:%M %z"
$ git rev-parse --short HEAD
$ git log --format='%h %ad %s' --date=format:'%H:%M' bb08d1f..HEAD --grep='^票 158 r2'   # 本程应交 4 枚
$ bash .scratch/wisp/probes/158/r2b/zero-byte-census.sh                                  # §C.5 末那条预测就是它的复算尺
```

从 **`bb08d1f`（本程 step-0 锚点，2026-09-26 19:37 现量）到 HEAD** 这一段里，**票 158 r2 的活已全部落库**；
剩下的三件事都不在实现方手里：
① §D.6／§E.4-1 那枚 d22scan 未裁子项要等 `frontend/**` 的 owner 换掉那枚字形；
② 票面 6 枚框的勾与 AC#1① 那一腿的重判要等**非实现者**验收程（§E.4-2 给了它可跑的形状）；
③ §B 那一处扩权请在验收表里**单列一行**，别并进 AC#2。

