# 票 133 AC#2 退回件的回修 —— `R-133-1`／`-2`／`-3`／`-4`／`-6` 的修法与复跑读数

**执行方**：`worker-ticket133-ac2-fix`（写码位，非验收方）
**裁决表来源**：`docs/evidence/s1/133-adversarial-acceptance.md`（`acceptor-ticket133-r1`，被验 sha `c5f140c`）
**票面**：`.scratch/wisp/issues/133-*.md`「AC#2 判退回 ＋ R-133-x 归单」那节（翻格条件＝`R-133-1`＋`R-133-6` 必办）

## 0. 锚点与树

- 开工 `git rev-parse --short HEAD` = **`54123e0`**（分支 `dev`）。开工时工作树未提交面：
  `.scratch/wisp/issues/133-*.md`／`135-*.md` 两枚 M、`137-*.md`、`docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`、
  `internal/observe/sampler_zerosample_136_test.go` 三枚 ?? ——**都不是本程的路径，一个字节没动**。
- 纯净树（**不在仓内建 worktree、不 checkout**）：
  `mkdir -p /d/tmp/wisp133-ac2fix && git archive 54123e0 | tar -x -C /d/tmp/wisp133-ac2fix`
  该树 `go.mod` 883 字节，`third_party/sherpa-onnx/` 不在归档里（未入库），按验收方同法补三枚 DLL：
  `cp <repo>/third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/`（`onnxruntime.dll` 17799168／
  `sherpa-onnx-c-api.dll` 4605952／`sherpa-onnx-cxx-api.dll` 259584 字节），缺即 0xc0000135。
- 复跑台：`/d/tmp/wisp133-ac2fix-mutate.py`、`/d/tmp/wisp133-ac2fix-probe.py`
  = 验收方 `/d/tmp/wisp133-acc-r1-{mutate,probe}.py` 的**逐字节副本**（`sha1sum` 相同），
  所以 §1／§3 的变异与探针是**同一发件**，不是本程另写一枚形状。驱动器 `/d/tmp/wisp133-ac2fix-run.sh`。
  日志与逐发读数全在 `/d/tmp/wisp133-ac2fix-out/`。
- 测量树：`/d/tmp/wisp133-ac2fix-pristine`（`54123e0` 的副本，还原核对基准）、
  `-T1`（改前复跑台）、`-T2`（改后复跑台，只把 `cmd/wisp/leg_dispatch_gate_133_test.go` 换成本程版本）。
  **每发读数前都核过 `diff -r -q` 到 `pristine` 的 difflines=0**（还原检查），并且这三棵树里
  **没有**兄弟代理在途中落进 `cmd/wisp/` 的新文件（见 §0.1）。
- **共树漂移（按名甩清）**：本程开工锚 `54123e0`；第 1 枚修法 commit 前 `git diff --name-only 54123e0..HEAD -- cmd/wisp/`
  量到 `cmd/wisp/providers_test.go`、`cmd/wisp/secret_test.go`、`cmd/wisp/tempdir_resolved_124_test.go` 三枚
  （兄弟 `worker-ticket124-ac2b-4` 的 AC#2b 批次）⇒ §1／§3 的对照读数都在 `54123e0` 基树上取，
  §4 的整包四数**另在含那三枚文件的合并态快照**上重取一遍，两批读数分开标。本程不碰那三枚文件。
- 本程只动 `cmd/wisp/leg_dispatch_gate_133_test.go` ＋ 本证据文件 ＋ 票面 append 一节 ＋
  前任证据文件 `133-ac1-ac2-instrument.md` 的 §3 更正段（只插入、不抹）。

### 0.1 基线（未变异，`54123e0` 纯净树）

```
$ sh /d/tmp/wisp133-ac2fix-run.sh /d/tmp/wisp133-ac2fix B0-base-pristine
START 2026-09-23T13:44:46Z label=B0-base-pristine tree=/d/tmp/wisp133-ac2fix args=
READ label=B0-base-pristine rc=0 RUN=101 TOPPASS=54 TOPFAIL=0 TOPSKIP=0
```

与验收方 §1 的 `A-base-ins-full` 101/54/0/0 逐字同（今天基线：无本尺 100/53/0/0、有本尺 101/54/0/0）。

## 1. 改前：验收方那四发探针在原树上的红绿原文

时刻 `date -u` 13:46:48Z–13:47:5xZ（本地 21:46–21:47 +08）。做法与验收方 §3.1 同形：
先 `python …-mutate.py <树> x14b2c` 种第五发第二拍（install 已拆掉，只剩判据 (d) 能让它红），
再 `python …-probe.py <树> <探针>` 加探针件；每枚先 `go build ./...` 取 rc 证明落地，
再单跑本尺 `-run '^TestAC1AC2DispatchHopGate133$'` 取判绿／判红。p5 按验收方 `extra.sh` 那次的形状在
**pristine `main.go`** 上种（`probe.py` 的 `clean()` 会还原一份陈掉的 `.p1.bak`，
会把 x14b2c 那行 case 带回 `main.go` 而让 `go build` 停在 `cmd\wisp\main.go:81:11: undefined: cmdSfx131`
——这一形本程踩到一次，改法与验收方同：先 cp pristine 再手种那行标签；两份读数都留着）。

| 探针 | 种了什么 | build | 本尺改前读数 | 账本那一行（原文，已裁到 190 列） |
| --- | --- | --- | --- | --- |
| `p4` 对照 | 只在 `usage` 里加一行 `wisp sfx131` | rc=0 | **rc=1 红**（`red=1 pass=0`） | `leg sfx131 main.go:81 installs=false handoff=false covered=RED nothing entries=cmdSfx131` |
| `p1` 同名方法 | 一枚**方法** `func (phantomRecv133) TestSfx131LegIsDriven(){ _ = cmdSfx131(nil) }` ＋上面那行 usage | rc=0 | **rc=0 绿**（`red=0 pass=1`） | `leg sfx131 … covered=test TestSfx131LegIsDriven drives cmdSfx131` |
| `p2` 不相邻裁决 | `doctor.go` 末尾一行 `// WISP-LEG-COVERAGE-RULING: sfx131 …`（与这条腿的 entry `cmdSfx131`（`sfx131x.go`）不同文件）＋ usage | rc=0 | **rc=0 绿** | `leg sfx131 … covered=ruling doctor.go:341` |
| `p3` 钉指向 helper | 钉表加一行 `{leg:"sfx131", test:"phantomHelper133", entry:"cmdSfx131"}`，`phantomHelper133` 是 `_test.go` 里一枚非 `Test` 前缀的普通 helper | rc=0 | **rc=0 绿** | `leg sfx131 … covered=nail phantomHelper133 -> cmdSfx131` |
| `p5` 同 case 多标签 | `case "run":` 改成 `case "run", "runalt133":`（pristine `main.go`，不加任何腿） | rc=0 | **rc=0 绿** | `dispatch ledger … (**11 legs**, 4 claims …)`，`runalt133` 在整份日志里**零行** |

`p1` 那枚"用例"到底跑没跑（整包、门关着 `-skip '^TestAC4EveryLegIsNailedOrRuled$'`）：

```
READ label=PRE-P-p1-fullpackage-doorclosed rc=0 RUN=100 TOPPASS=53 TOPFAIL=0 TOPSKIP=0
roster-has-phantom=1        （名字在整份日志里出现 1 次，唯一那次是本尺自己打的账本行）
phantom-RUN=0               （grep -c '^=== RUN   TestSfx131LegIsDriven$'）
phantom-PASSFAIL=0          （grep -c '^--- (PASS|FAIL): TestSfx131LegIsDriven'）
top-level-test-count=53     （^--- (PASS|FAIL): Test 共 53 枚，与基线同一枚数，没多一枚）
```

⇒ 复现验收方 §3.1 的四形：**(a) 认同名方法＝成立（p1 绿）、(b) 认 helper 名＝成立（p3 绿）、
裁决句可在另一枚文件末尾拿到背书＝成立（p2 绿）、多标签第二条命令整条不出账本＝成立（p5 绿）**，
而诚实的那一发（p4，只加给人看的那行 usage）红。

## 2. 每条修法

### 2.1 `R-133-1`（高）＋ `R-133-2`（中）：覆盖桶只收 Go 真会跑的那一形

根因两处（验收方 §3.0 的字节）：`collectTopLevel133` 把 `_test.go` 里**任何**声明按**裸名**进 `tests` 桶
（不看 receiver／签名），`drivenBy133` 只认 `HasPrefix("Test")` ⇒ 一枚永不运行的同名方法就是"覆盖"。

修法照 131 那扇门收紧（`leg_sink_gate_131_test.go:617` 明查 `fn.Recv == nil && HasPrefix(fn.Name.Name,"Test")`），
并比它多核一步签名：

- 新增 `isRunnableCase133(fd, testingPkg)`：`fd.Recv == nil` ＋ `Test` 前缀 ＋ 参数表恰一枚、类型
  `*<testing>.T`（`testingPkgName133` 按该文件的 import 别名解析，`import t "testing"` 这一形不放过）＋ 无返回值。
  只有这种声明进 `p.tests`。
- 其余 `_test.go` 里的声明（helper、方法、Benchmark）落进新增的 `p.helpers` 桶。
  **这是为了不退化**：`testClosure133`/`lookup133`/`resolve133` 的 loose 走图现在两桶都看得见，
  所以"一枚真用例经 helper 驱动这条腿"仍是覆盖证据（前任账本里 `leg doctor`／`leg providers` 那两行
  `covered=test TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 drives cmdDoctor|cmdProviders` 逐字未动，
  见 §2.1 末的 ledger diff）。
- `drivenBy133` 的 `HasPrefix` 保留，注释改成写明它是**同一枚不变式的第二遍读**，真正的
  receiver／签名判定在 `isRunnableCase133`。
- `R-133-2`：钉表 `legCovers133.test` 查的就是这枚收紧后的桶 ⇒ 文案 "is not a test function" 与判据一致；
  红名再加一句"它找到的是什么"（`declares probe133_test.go:6, which this file sorts with isRunnableCase133:
  a helper or a method, not a case Go's testing package runs`／或"根本没有这个名字"）。

**没有采的另一种修法（报备，不是漏做）**：验收方 §8 给的第一个选项是"登记表存**函数值**（同
`registerLegNail131`，131 `:225`）"。本程**装不了**这一形：`legCovers133` 现在那四枚声明指向的用例
逐枚在 `cmd/wisp/leg_sink_nail_131_windows_test.go:346,407`、`cmd/wisp/logsink_windows_test.go:137`、
`cmd/wisp/resident_sink_nail_127_windows_test.go:393`，文件名后缀 `_windows_test.go` ＝ 隐式 `GOOS=windows`；
而本尺 `leg_dispatch_gate_133_test.go` 没有构建后缀、两平台都编。把四枚声明换成函数值会叫
`GOOS=linux go vet ./cmd/wisp/`／`go test` 直接 `undefined: TestAC2SealNotice…` ——
把"关着门仍红"的那把独立尺换成"linux 上根本编不过"，正是派单要我**报回来不要硬修**的那一类。
131 那份钉表能存函数值，是因为它的注册点就在 windows 那枚文件里（`registerLegNail131` 由
`leg_sink_nail_131_windows_test.go` 自己调），本尺的注册点跨不过这个构建标签。
⇒ 本程采验收方给的第二个选项（"至少在查表时校验 receiver 为 nil ＋ 参数表含 `*testing.T`"），
把判据与文案对齐；函数值那一支留给"钉注册点也搬进 windows 文件"的后续格（AC#3），本程不顺手做半截。

（待填：2.2 R-133-3／2.3 R-133-4／2.4 R-133-6）

## 3. 改后：同一发复跑

### 3.1 `p1`／`p3`（判据 (d) 的名字那一侧）

时刻 `date -u` 13:55:35Z–13:56:0xZ（本地 21:55–21:56 +08），树 `/d/tmp/wisp133-ac2fix-T2`
（= `54123e0` 基树 ＋ 本程第 1 枚修法）。每发都先 `go build ./...` 证明种件落地。

```
########## POST PROBE=p4 2026-09-23T13:55:35Z      （对照，仍须红）
BUILD POST-p4 rc=0
GATE POST-p4 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing: …
  leg sfx131 main.go:81 installs=false handoff=false covered=RED nothing entries=cmdSfx131
########## POST PROBE=p1 2026-09-23T13:55:45Z      （同名方法：改前绿 ⇒ 改后红）
PROBE p1 applied
  probe133_test.go:10: func (phantomRecv133) TestSfx131LegIsDriven() { _ = cmdSfx131(nil) }
BUILD POST-p1 rc=0
GATE POST-p1 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing: …
  leg sfx131 main.go:81 installs=false handoff=false covered=RED nothing entries=cmdSfx131
（账本表头：dispatch ledger … (12 legs, 4 claims …) —— 12＝11＋这条被分发的 sfx131）
########## POST PROBE=p3 2026-09-23T13:55:53Z      （钉表收 helper 名：改前绿 ⇒ 改后红，两枚红）
PROBE p3 applied
  leg_dispatch_gate_133_test.go:140: 	{leg: "sfx131", test: "phantomHelper133", entry: "cmdSfx131"},
BUILD POST-p3 rc=0
GATE POST-p3 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:176: AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing: …
leg_dispatch_gate_133_test.go:176: AC#1 RED: this gate claims leg "sfx131" is nailed by "phantomHelper133", which is not
  a test function in this directory's sources: declares probe133_test.go:6, which this file sorts with isRunnableCase133:
  a helper or a method, not a case Go's testing package runs. …
```

⇒ 红名两发都点到本尺 `TestAC1AC2DispatchHopGate133`（`grep -c '^--- FAIL: TestAC1AC2DispatchHopGate133'` = 1，
`--- PASS` = 0），不是包级 `[build failed]`（`go build ./...` rc=0）。

**未变异树上的账本逐字节不变**（收紧判据不许改今天的读数）：

```
$ diff /d/tmp/wisp133-ac2fix-out/LEDGER-pre-fix.txt /d/tmp/wisp133-ac2fix-out/LEDGER-post-fix1.txt   rc=0（空输出）
      11 /d/tmp/wisp133-ac2fix-out/LEDGER-pre-fix.txt
      11 /d/tmp/wisp133-ac2fix-out/LEDGER-post-fix1.txt
（两枚都是 `54123e0` 未变异树上 `-run '^TestAC1AC2DispatchHopGate133$'` 的 11 行 `leg …` 账本）
$ gofmt -l cmd/wisp/   rc=0 空输出
$ go vet ./cmd/wisp/   rc=0
```

（待填：3.2 p2／p5 复跑、3.3 五发老账不退化）

## 4. 四数逐名账（待填）

## 5. 门禁（待填）

## 6. 未验证与 next=（待填）
