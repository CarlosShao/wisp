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
