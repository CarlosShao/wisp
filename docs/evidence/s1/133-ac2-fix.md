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

### 0.05 共树未提交面（别人的，本程一个字没动）

`git status --short` 在本程中途起出现的兄弟未提交件：`.scratch/wisp/issues/124-*.md`、
`docs/evidence/s1/124-ac2b-4-conversion.md`（`worker-ticket124-ac2b-4` 在飞），
以及开工时就有的 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`。
本程每一枚 commit 都带**显式 pathspec**，`git diff --cached --name-only` 逐枚核过暂存清单
只含本程路径；未替兄弟 add、未替兄弟 commit、未 `git add -A`。

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

### 2.2 `R-133-3`（中）：裁决句要在 owning 这条腿的文件里才算覆盖

根因字节（验收方 §3.1 p2 那一发命中的）：`collectRulings133` 是**行扫描**，只取标记后第一个
token 当腿名、first-seen 即停，既不要求与这条腿的代码相邻、也不要求在同一枚文件里。

修法（**保留"第二把尺"这条线**：行扫描原样留着，没有改成从 AST 里读注释——派单硬约束 1
不许把这把尺改回只信 AST，票 135 面上同一形已按此裁过一次）：

- 扫描时把每一枚标记记成 `ruling133{leg, site, file}`（新增 `p.rulingAll`），
  `p.rulings[leg]` 仍是 first-seen 站点 ⇒ 两枚**不认领覆盖**的判据能见度不降级：
  "裁决＋钉同时在场"（`carries a coverage ruling … AND a nail claim`）与
  "裁决指向不在册的腿"（`which is not in the dispatch census`）照旧看得见任意文件里的标记。
- `rulingAdjacentTo133` ＋ `legOwnsFile133`：只有坐在"这条腿的 dispatch 所在文件 ∪ 它的 entry
  声明所在文件"里的标记才被认作覆盖证据（账本 `covered=ruling <site>` 打的是这一枚站点）。
- `misplacedRulingReds133`：一枚不落在任何 owning 文件里的标记**单独报红**，
  不让它因为"这条腿另有钉"就悄悄躺在树上。
- **没有采的形状**：验收方 §8 另半句"裁决句必须点到该腿的 entry 名"= 要求标记正文里出现
  `cmdSfx131` 这类符号名。今天树上那五枚裁决句（`main.go:65/70/74`、`panel_assets.go:16`、
  `slo_windows.go:186`）没有一枚写到 entry 符号名——加上这条判据就得改 `main.go` 等三枚生产文件
  的裁决注释才不红，而那三枚不在本程地界（本程只碰 `leg_dispatch_gate_133_test.go`），
  故**只做同文件核对**，把"点到 entry 名"登记进 §6 未验证面。

### 2.3 `R-133-4`（中）：一条 case 的每一枚命令标签各出一行腿

根因字节（原 `:965`）：`if key == "" { … key = v }` ⇒ 一条 clause 只有**第一条字面量标签**进清单，
第二条命令名整条从账本消失，也不与 usage 双向差撞。

修法 `caseLabelLegs133`：

- **命令词形态的标签逐枚成腿**：各自的覆盖义务（判据 (b)）、各自欠 `usage` 一行（双向差）。
- **旗标拼写折进主行的 `aliases=` 列**，条件是它的字母部分是主命令词的**前缀**
  （`-v`/`--version` 是 `version` 的拼写；`-h`/`--help` 是 `help` 的拼写）——
  今天 `main.go` 那两枚 clause 就是这个形状，折起来才不至于把"CLI 惯例的别名"当成第二条命令；
  不是前缀的旗标（`-x133`）不折叠，仍是一条自己的腿，既欠钉也欠文档。
- 非字面量标签仍**各自**成 `unparsed-label@<site>:<expr>` 腿并红（X8 的红名文案逐字未动；
  多条非字面量标签也不再互相吞掉）。
- `leg133` 新增 `aliases []string`，账本行尾部多一列 `aliases=…` ⇒ 别名从"零行"变成"在行上看得见"。

### 2.4 `R-133-6`（中，证据面）：把那句全称量词按实际读数更正

`docs/evidence/s1/133-ac1-ac2-instrument.md` §3 那句"把本尺摘掉，这五发的形状今天仍然
没有任何一枚用例会红"——**原文留在原地未抹**，紧跟其后插入一段带 `>` 引块的更正
（`git diff --numstat` 那枚文件删除列 = 0，只有 20 行插入）。更正段给的读数是：
摘掉本尺后**只有 X14 两拍零红**，N-3／X4／X8／X12 四发各红一枚、红的是票 131 的门
`TestAC4EveryLegIsNailedOrRuled`（验收方 §2.1 ①栏四行 100/52/1/0），并与同一份文件 §2.0.1
的自述对齐。本程自己的六发复算（在同一棵有本尺的树上 `-skip` 掉本尺）见 §3.4。

### 2.5 `R-133-5`（低，验收方未造 ⇒ 本程自造之后才算已证）

验收方 §8 那一格只登记了字节（`classifyCond133` 字面量分支 last-match、no-args 分支 first-match），
**没有造出探针**，并写明"下一位请补"。本程补了（探针 `p6`，见 §3.5）：造出来确实能绿 ⇒
按派单那句"若你造出来了再一并修并附三态"一并修掉，并从这一枚起它**是已证**、不再是读码断言。

修法：`classifyCond133` 改成返回 `(keys, aliases, kind)`，条件里**每一处** `argv 槽位 == 字面量`
比较都进 `keys`；`len(argv) == 0` 与字面量比较同时在场时两条腿都留（旧代码里字面量会吃掉 no-args）。
标签折叠规则与 `R-133-4` 合并成同一处 `splitCommandLabels133`（旗标拼写须是命令词的前缀才折进
`aliases`，否则自成一腿），两条判据不再各写一份、也就不会各自漂移。

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

### 3.2 `p2`（不相邻裁决句）与 `p5`（同 case 多标签）

时刻 `date -u` 13:59:3xZ–14:06Z（本地 21:59–22:06 +08）。
**先证一件事**：第一次跑 `p2` 用的树 `T2` 里，`probe.py` 的 `clean()` 把更早一轮留下的
`leg_dispatch_gate_133_test.go.p1.bak`（**只含第 1 枚修法**的版本）盖回了我刚拷进去的文件，
那次读到的"仍绿"是**改前的读数**，作废（日志 `POST-p2.log` 留盘不删，但不当证据）；
下面两发都在新建的树上看 `sha1sum` 与仓内文件逐字相同之后才取。

`p2`，树 `/d/tmp/wisp133-ac2fix-T3`（基树 `54123e0` ＋ 修法 1＋3）：

```
LANDED-MAIN x14b2c:
  main.go:80: 		case "sfx131":
  sfx131x.go:12: var holder131 = sinkHolder131{}
PROBE p2 applied
  doctor.go:341: // WISP-LEG-COVERAGE-RULING: sfx131 acceptor probe: planted in an unrelated production file, far from any code that owns this leg.
BUILD POST3-p2 rc=0
GATE POST3-p2 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:175: AC#1 RED: doctor.go:341 carries WISP-LEG-COVERAGE-RULING: for leg "sfx131", but doctor.go owns
  neither this leg's dispatch (main.go:81) nor any symbol it calls (main.go, sfx131x.go). A coverage ruling has to sit next to the
  code that owns the leg: as placed, one line at the end of any production file in this directory would back the claim that
  somebody verifies "sfx131", which is the reading ticket 133's own acceptor registered as R-133-3.
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing: …
  leg sfx131 main.go:81 installs=false handoff=false covered=RED nothing entries=cmdSfx131
```

⇒ 由绿转红，两枚红都点到 `TestAC1AC2DispatchHopGate133`，`go build ./...` rc=0（不是包级坏掉）。
**正向控制一枚**（收紧判据不许把"写在该写地方的裁决句"一起打死）：同一发 x14b2c ＋ 同一行 usage，
把裁决句从 `doctor.go` 末尾改到**这条腿自己的文件** `sfx131x.go:22`：

```
BUILD POST-p2b rc=0
GATE POST-p2b rc=0 red=0 pass=1
  leg sfx131 main.go:81 installs=false handoff=false covered=ruling sfx131x.go:22 entries=cmdSfx131 aliases=
```

`p5`，树 `/d/tmp/wisp133-ac2fix-T4`（基树 `54123e0` ＋ 修法 1＋3＋4）：

```
P5 applied（pristine main.go 上手种 `case "run", "runalt133":`，与验收方 extra.sh 同法）
  main.go:82: 		case "run", "runalt133":
BUILD FIN-p5 rc=0
GATE POST-p5 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "runalt133" (main.go:80) reaches installLogSink on this path: dispatch -> runTextTask -> installLogSink …
leg_dispatch_gate_133_test.go:185: AC#1/#2 RED: leg "runalt133" is dispatched by func main and documented in no line of the usage block …
dispatch ledger, read out of func main's own branches at run time (12 legs, 4 claims in this gate's registry):
  leg run        main.go:80  installs=true  handoff=false covered=nail TestAC2SealNoticeLandsInTheRunLegLogFile -> runTextTask entries=attachParentConsole|cmdRun aliases=
  leg runalt133  main.go:80  installs=true  handoff=false covered=RED sink with no nail                                   entries=attachParentConsole|cmdRun aliases=
```

⇒ 由绿转红：第二条标签现在**有自己的行**（12 legs），并且因为这条腿够得着听众又没钉、
usage 里也没那行字，两枚红都点到本尺。别名折叠的正向控制同树可见：
未变异的树上 `leg version … aliases=--version|-v`、`leg help … aliases=--help|-h`、
`11 legs` 与基线逐字同（`go vet ./cmd/wisp/` rc=0）。

### 3.5 `p6`（本程自造的 `R-133-5` 探针，三态齐全）

种法：pristine `main.go` 的 `switch args[0] {` 之前插
`if len(args) > 0 && (args[0] == "x133a" || args[0] == "version") { attachParentConsole(); printVersions(""); return }`，
**不动 usage 块**（这台件 `/d/tmp/wisp133-ac2fix-p6.sh`，两棵树只差第 5 枚修法）。

```
# 改前（树 /d/tmp/wisp133-ac2fix-T7 ＝ 54123e0 基树 ＋ 修法 1-4）
P6 applied (compound early if; x133a deliberately NOT in the usage block)
  main.go:79: 	if len(args) > 0 && (args[0] == "x133a" || args[0] == "version") {
BUILD PRE5-p6 rc=0
GATE PRE5-p6 rc=0 red=0 pass=1
  dispatch ledger … (12 legs, 4 claims in this gate's registry):
  leg version main.go:79  … covered=ruling main.go:65 entries=attachParentConsole|printVersions
  leg version main.go:111 … covered=ruling main.go:65 entries=attachParentConsole|printVersions aliases=--vers
X133A-ROWS=0 NAME-APPEARANCES=0        （"x133a" 在整份读数里出现 0 次）
# 改后（树 /d/tmp/wisp133-ac2fix-T8 ＝ 同一基树 ＋ 修法 1-5）
BUILD POST5-p6 rc=0
GATE POST5-p6 rc=1 red=1 pass=0
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "x133a" (main.go:79) is dispatched by func main and covered by nothing: …
leg_dispatch_gate_133_test.go:185: AC#1/#2 RED: leg "x133a" is dispatched by func main and documented in no line of the usage block: …
  dispatch ledger … (13 legs, 4 claims …):
  leg x133a   main.go:79  installs=false handoff=false covered=RED nothing entries=attachParentConsole|printVersions aliases=
X133A-ROWS=1 NAME-APPEARANCES=3
```

⇒ 一枚只写在复合条件**第一个**槽位上的命令名，改前整条不进账本（0 次出现）、判绿；
改后有自己的行、两条红，红名仍是本尺。

### 3.3 五发老账不退化（最终 sha `fbf420c`，门关着）

台件与时刻：树 `/d/tmp/wisp133-ac2fix-T9` = `git archive fbf420c | tar -x` ＋ 三枚 DLL
（`ARCHIVE-MODULE` 行：仓内与本尺文件 `sha1sum` 同为 `906201f4a3995d10b0a65910aa4cc68e1f745a9d`；
`git diff --name-only fbf420c HEAD -- cmd/wisp/` 空 ⇒ 那棵快照里没有兄弟的新文件）。
`date -u` 14:27:49Z–14:45:3xZ（本地 22:27–22:45 +08）。
门关着的做法与验收方同：`-skip '^TestAC4EveryLegIsNailedOrRuled$'`（131 的门不跑＝不存在，
它的两枚文件本程一个字没动）。每发先 `python …-mutate.py <树> <发>` 打印 LANDED 行＋`go build ./...` rc 证落地。

| 发 | build | 门关着读数 | 红名（`grep '^--- FAIL: '`） | 六栏机检 |
| --- | --- | --- | --- | --- |
| N-3 | rc=0 | rc=1 100/52/1/0 | `TestAC1AC2DispatchHopGate133` | 131 门提及 0／SKIP 行 0／本尺 RUN 1／本尺 FAIL 1／build failed 0／panic 0 |
| X4 | rc=0 | rc=1 100/52/1/0 | 同上 | 同上 |
| X8 | rc=0 | rc=1 100/52/1/0 | 同上 | 同上 |
| X12 | rc=0 | rc=1 100/52/1/0 | 同上 | 同上 |
| X14 一拍（`x14c`） | rc=0 | rc=1 100/52/1/0 | 同上 | 同上 |
| X14 二拍（`x14b2c`） | rc=0 | rc=1 100/52/1/0 | 同上 | 同上 |

⇒ **六发两拍零退化**：红名逐发仍点到本尺（不是包级 `[build failed]`，`go build ./...` 六发全 rc=0；
131 的门在③栏日志里出现 0 次）；`--- SKIP` 行全 0（关门是"少跑一枚"而不是"跑成 SKIP"，
RUN 从 101 到 100 就是那一枚被 `-skip` 掉的 131 门）。逐名与四枚修法之前那一次
（同法、树 `090bb3e`、前缀 `F-`，§3.1／§3.2 引用的那批日志）逐字同。

### 3.4 摘掉本尺的六发（`R-133-6` 那句更正的本程第一手读数）

同一棵树、同六发，把 `-skip` 换成**摘掉本尺**（`-skip '^TestAC1AC2DispatchHopGate133$'`，
131 的门照跑）——这一形问的就是前任 §3 那句话："把本尺摘掉，还有没有用例会红"：

| 发 | 摘掉本尺的读数 | 红的那一枚是谁 |
| --- | --- | --- |
| N-3 | rc=1 100/52/1/0 | `TestAC4EveryLegIsNailedOrRuled`（票 131 的门） |
| X4 | rc=1 100/52/1/0 | 同上 |
| X8 | rc=1 100/52/1/0 | 同上 |
| X12 | rc=1 100/52/1/0 | 同上 |
| X14 一拍 | **rc=0 100/53/0/0** | 无人（全树零红） |
| X14 二拍 | **rc=0 100/53/0/0** | 无人（全树零红） |

⇒ 本程自己量到的读数与验收方 §2.1 ①栏（无仪器树 `c720494`）同一命题：**"摘掉本尺五发零红"是假的**，
四发仍由 131 的门红、只有 X14 两拍零红。§2.4 那段更正按这一枚读数写。
（两形树基不同的那笔账在 §6 第 7 条，别把这两枚当成同一枚读数。）

### 3.6 四枚探针在最终树上的复跑（`fbf420c` ＋ 第 5 枚修法之后，见 §3.6 段日志 `FIN-*.log`）

（待填：B 脚本读数）

## 4. 四数逐名账（合并态快照，只从 `-v` 量）

**树**：`/d/tmp/wisp133-ac2fix-T5` = `git archive 090bb3e`（本程第 4 枚修法 commit，
落在兄弟 `worker-ticket124-ac2b-4` 的 AC#2b 三枚文件之上）`| tar -x` ＋ 三枚 DLL。
**先证量的是被交版本**：`sha1sum cmd/wisp/leg_dispatch_gate_133_test.go` 仓内与归档后
同为 `18b74b45c747487be419a1494bc17026d5e563cd`（脚本 `FINAL-A.log` 的 `ARCHIVE-MODULE` 行）。
时刻 `date -u` 14:10:08Z–14:13:5xZ（本地 22:10–22:13 +08）。

| 命令 | rc | RUN | PASS | FAIL | SKIP | 对基线 |
| --- | --- | --- | --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | 0 | 202 | 108 | 0 | 0 | 字面 AC#5 命令；与前任 §4.1、验收方 §7.1 的 202/108/0/0 逐字同；`[no tests to run]` 命中 0；本尺在 count=2 里 `--- (PASS\|FAIL): TestAC1AC2DispatchHopGate133` 出现 2 次 |
| `go test -count=1 -v ./cmd/wisp/` | 0 | 101 | 54 | 0 | 0 | ＝今天"有本尺"基线 101/54/0/0 |
| `… -count=1 -v -skip '^TestAC1AC2DispatchHopGate133$'` | 0 | 100 | 53 | 0 | 0 | ＝今天"无本尺"基线 100/53/0/0（`-skip` 是少跑一枚、不是跑成 SKIP：SKIP 列仍 0） |

**逐名 `comm`（不是差值）**，两份名册 = 上面两枚读数的顶层 `--- (PASS|FAIL|SKIP):` 行：

```
comm -23 F-roster-ins.txt F-roster-skip.txt  ->  TestAC1AC2DispatchHopGate133   （新增，仅此一枚）
comm -13 F-roster-ins.txt F-roster-skip.txt  ->  （空）                          （无一名消失）
comm -12 F-roster-ins.txt F-roster-skip.txt  ->  53 枚共有，一字未动
```

**四枚修法各自动了哪几枚数**：全部改在判据里（`isRunnableCase133`／
`rulingAdjacentTo133`＋`misplacedRulingReds133`／`caseLabelLegs133`），
**零新增用例、零删除用例、零改名** ⇒ 名册逐名可核：有本尺 54 枚、无本尺 53 枚，
差的仍只有本尺那一枚，与改前（验收方 §7.1 在同一枚尺上量的 54/53）逐名同。
本程**没有**另起 `*_133_test.go` 新件：四条的"探针"是**变异件**（装在快照树里、跑完留盘不进仓），
不是一枚新的进仓用例；AC#3 那句"给这条判据装一枚主动弄哑自己的自证腿"是**另一格**（见 §6）。

## 5. 门禁（待填）

## 6. 未验证、没做的、与要报回来的

1. **AC#2 的翻勾权不在本程**。派单写死"不许自己勾 AC#2"，票面那一节也只 append 不翻格。
   本程交的是四条（＋一条自造）修法与复跑读数，终裁按 `133-adversarial-acceptance.md`
   的 §3.1／§8 由非实现者验收方复跑判。
2. **钉表存函数值那一支没装**（`R-133-2` 的第一选项）。理由与逐枚文件名见 §2.1 末：
   四枚钉指向的用例全在 `_windows_test.go` 后缀的文件里，本尺是跨平台文件，
   换成函数值会让 linux 侧 `go vet`／`go test` 直接 `undefined:` ——
   把"关着门仍红"的独立读数换成"编不过"，属派单点名的"修了会以另一种方式变哑"，**报回不硬修**。
   代价登记在这里：`covered=nail …` 那一列仍是**按名字**对源码做的查表，
   build-tag 藏起来的用例今天仍会被认作钉（本尺 header 的 blindness disclosure 一行已自陈）。
3. **`R-133-3` 只做到"同文件核对"**，没做验收方那半句"裁决句必须点到该腿的 entry 名"。
   今天五枚裁决句（`main.go:65/70/74`、`panel_assets.go:16`、`slo_windows.go:186`）
   没有一枚写了 entry 符号名 ⇒ 装上就得改三枚生产文件的裁决注释才不红，那出了本程地界
   （本程只碰 `leg_dispatch_gate_133_test.go`）。**残留形状**：一枚坐在 `main.go` 里、
   点名 `version`、但不提任何符号名的裁决句仍然算覆盖。
4. **AC#3 那句"给这条判据装一枚主动弄哑自己的自证腿"未做**——裁决表 §8 的修法栏把它列进
   `R-133-1` 那一格，但 AC#3 本身"本表未裁"（不在派单地界）。本程四法的"自证"是**变异探针**
   （装在快照树里、跑完留盘不进仓），不是一枚常驻的弄哑自证腿。
5. **`R-133-7`（gofumpt 版本没钉在 CI／d22scan 各 scope 无可比历史基线）不在本程**：
   §5 只写明**本程这一枚** gofumpt 的版本，没有代 CI 钉版本、也没有立进仓的 scope 基线数；
   那一格仍归票 133 的 AC#5。
6. **linux 面**：`GOOS=linux go vet ./cmd/wisp/` 停在 cgo 包加载（诊断里没有一枚
   `cmd/wisp/*.go:line` ⇒ 既非破口亦非清白），真类型读数只有 `golang:1.27` 容器那一枚；
   linux 上 `go test ./cmd/wisp/` 的 19 枚 FAIL 本程**未复跑**（package owner 地界，
   且 `scripts/wisp-cli-tests.sh` 明写那是 windows 腿）。
7. **五发的①栏（无仪器树 `c720494`）本程未重取**：§3.4 用的是"同一棵有本尺的树上 `-skip` 掉本尺"
   这一形，与验收方①栏的树基不同（那棵树连 `main.go`/`panel_assets.go`/`slo_windows.go` 的
   19/8/6 行裁决注释都没有）。两形都在说同一句"摘掉本尺之后还有谁红"，但**别把它当成同一枚读数**。
8. **AC#6（仪器自己进 CI：说不出 run id＋job id＋step 名就当不存在）本程未动**：只 commit 未 push。
9. **容器挂载与 TMPDIR 两枚坑**的处理见 §5：`MSYS_NO_PATHCONV=1` ＋ `/d/...` ＋ 容器内
   `ls -l /src/go.mod` 先证挂上；软链 temp 那一形是**另一枚读数**，见 §5 末（分母是否缩小要分清）。

`next=` 交回编排者：请**非实现者验收方**在本程最终 sha 上原样复跑 `p1`／`p2`／`p3`／`p5`
（四枚已知探针）＋ `p6`（本程自造的 `R-133-5`）＋ 五发两拍，然后按裁决表 §3.1／§6／§8 判
AC#2 是否翻格；`R-133-2` 的函数值那一支与 `R-133-3` 的 entry-name 半句**留给 AC#3／AC#5 那两格**
（要动生产文件，别在续单里顺手做）。本程不改判据阈值、不动 131 的门与钉。
