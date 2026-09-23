# 票 133 AC#1 / AC#2 —— 非实现者验收方的裁决表

**执行方**：`acceptor-ticket133-r1`（**非实现者**，只读被验代码；不复用实现方的任何判据）
**被验版本锚定**：`git rev-parse --short HEAD` = **`52cf311`**（本验收方开工时刻的工作树 HEAD，
仅供对照）；**被验 sha = `c5f140c`**（= `c5f140c7a1c2be66ea16f1fcb09ae56737ea1ccf`）。
**纯净树取法**（不在仓库内建 worktree、不 checkout）：

```
$ mkdir -p /d/tmp/wisp133-acc-r1
$ git archive c5f140c | tar -x -C /d/tmp/wisp133-acc-r1        # 有仪器的被验树
$ mkdir -p /d/tmp/wisp133-acc-r1-noins
$ git archive c720494 | tar -x -C /d/tmp/wisp133-acc-r1-noins  # 无仪器的"①栏"树（前任的基线树）
$ cp third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/   # 三枚 DLL，缺即 0xc0000135
```

`git diff --stat c720494..c5f140c -- cmd/wisp/` ⇒ 只有 `leg_dispatch_gate_133_test.go` +1377 与
`main.go`/`panel_assets.go`/`slo_windows.go` 的 +19/+8/+6 行（裁决注释），与前任 §4.5 同形。

开工时刻 `date -u` = 2026-09-23 12:11:01z（本地 20:11 +08）。

## 0. 裁决表（逐节填毕，每裁一格 commit 一次）

| 格 | 判据出处 | 判定 |
| --- | --- | --- |
| AC#1 五发（N-3／X4／X8／X12／X14 两拍，门关着仍红） | 票面 AC#1 + 归属硬线 | **PASS**（§2：六栏三态独立复现，③栏红名逐发点到 `TestAC1AC2DispatchHopGate133`；①栏"今天不红"这一步的语义已失效，裁定见 §5，不构成本格的减分） |
| AC#2 覆盖面主张（清单式 vs 图式＋"删掉它哪条用例会红"） | 票面 AC#2 | **退回**（§3：覆盖判据 (d) 今天可被一枚**永不运行的同名方法**满足，红被洗掉；另两处同族洞。§6：§3 那句"摘掉本尺五发零红"与它自己 §2.0.1 矛盾） |

## 1. 独立基线（未变异，四枚读数）

驱动器 `/d/tmp/wisp133-acc-r1-run.sh`：`go test -count=1 -v ./cmd/wisp/`，`PATH` 前置该树自己的
`third_party/sherpa-onnx`（三枚 DLL）。四数只从 `-v` 输出量（`^=== RUN` / `^--- PASS` /
`^--- FAIL` / `^--- SKIP`）。日志 `/d/tmp/wisp133-acc-r1-out/A-base-*.log`。
测量 `date -u` 12:16:32z–12:20:1xZ（本地 20:16–20:20 +08）。

| 读数 | 树 | 附加旗标 | rc | RUN | PASS | FAIL | SKIP | 对前任 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `A-base-ins-full` | `c5f140c`（有仪器） | 无 | 0 | 101 | 54 | 0 | 0 | §2.0 `S0-ins-pristine` 101/54/0/0 逐字同 |
| `A-base-noins-full` | `c720494`（无仪器） | 无 | 0 | 100 | 53 | 0 | 0 | §1 基线 100/53/0/0、§2.0 `S0-base-pristine` 逐字同 |
| `A-base-ins-doorclosed` | `c5f140c` | `-skip '^TestAC4EveryLegIsNailedOrRuled$'` | 0 | 100 | 53 | 0 | 0 | §2.0 `S0-ins-doorclosed` 100/53/0/0 逐字同 |
| `A-base-ins-skipself` | `c5f140c` | `-skip '^TestAC1AC2DispatchHopGate133$'` | 0 | 100 | 53 | 0 | 0 | §4.3 基线独立复现 100/53/0/0 逐字同 |

**"变绿"与"被跳过"分清**（关门的读数是"少跑一枚"而不是"跑成一枚 SKIP"）：
`-skip` 后 RUN 101→100、PASS 54→53、**SKIP 仍 0**（`--- SKIP` 行不产出），
且被验尺自己在那份日志里仍有 `=== RUN   TestAC1AC2DispatchHopGate133` 一行
（`grep -c '^=== RUN   TestAC1AC2DispatchHopGate133$' A-base-ins-doorclosed.log` = 1）。
⇒ 门关着的三枚读数里，本尺是**真跑了再判**的，不是被一起跳掉才绿。

**门禁同形**（本验收方自己量，见 §7）：`c5f140c` 树上 `go build ./...` rc=0、`go vet ./cmd/wisp/` rc=0。

## 2. AC#1 五发独立复跑（①无仪器／②门开着／③门关着）

本验收方自己写的变异台 `/d/tmp/wisp133-acc-r1-mutate.py`（每处改动先 `assert` 锚点唯一、改完打印
落地行 + `go build ./...` rc，然后才读数；还原＝从 `.acc-r1.bak` 拷回 `main.go`、种件改名 `*.off` 不删），
驱动器 `/d/tmp/wisp133-acc-r1-run.sh`，日志 `/d/tmp/wisp133-acc-r1-out/*.log`。
门关着的做法：`-skip '^TestAC4EveryLegIsNailedOrRuled$'`（131 那扇门不跑＝不存在，它的文件一个字没动）。

**三轮种法的账（先说清，因为前两轮的读数与前任不符，原因在本验收方的种法，不在被验的树）**：

| 轮 | 种法 | 后果 |
| --- | --- | --- |
| v1 | 四条新腿都用 `if` 早退分发，且种件里调 `resolveDataDir` | 除 131 的门以外，还多点红出票 128 的门（`cmdXxx131 解析了数据根却没有腿驱动它`），且 128 的门 `t.Fatalf` 掉自己 5 枚子用例 ⇒ RUN 分母从 100 缩到 95 |
| v2 | X12／X14 改用票面明写的 `case "…":` 进 `switch`；X4 保留早退 `if`（那就是 X4 的形） | 四数与前任逐字对上；X14 两拍仍多一枚红：种件里那行 `slog.Info` 走的是**进程级听众**，131 的 R-117-1 判据（记了账却没装听众、又没裁决句）据此判红 |
| v3 | 在 v2 之上把种件里的 `slog.Info` 去掉（复验方当年那条账本原文就是 `install=false records=false`） | 六栏读数与前任 §2.1–§2.5 逐字对上 |

**三轮都留盘**（`shots.txt`＝v1，`shots-v2.txt`，`v3-probes.txt`＝v3＋仪器探针），因为前两轮的
红名差本身就是读数：**同一形在不同种法下会多点红不同门**，这对"归属硬线"的判断有影响（见 §5）。

### 2.1 五发三态总表（采用与票面形状同形的 v2/v3 轮；N-3／X8 只有 v1＝v2 同种法）

| 发 | 轮 | ①无仪器 `c720494` | ②门开着 `c5f140c` | ③门关着（`-skip` 131 的门） | ④还原 |
| --- | --- | --- | --- | --- | --- |
| N-3 摘掉 `cmdModels` 的调用者 | v1 | rc=1 100/52/1/0 红=`TestAC4EveryLegIsNailedOrRuled` | rc=1 101/52/2/0 红=本尺＋131 的门 | rc=1 100/52/1/0 红=**只有 `TestAC1AC2DispatchHopGate133`** | difflines=0 |
| X4 早退 `if` 的 `--diag` 腿 | v2 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 红=本尺＋131 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X8 `case "slo":` 换命名常量 | v1 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X12 install 经包级函数值别名 | v2 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X14 一拍 字段初值装听众 | v3 | **rc=0 100/53/0/0 零红（门 PASS）** | **rc=1 101/53/1/0 红=只有本尺** | **rc=1 100/52/1/0 红=只有本尺** | difflines=0 |
| X14 二拍 拆掉 install | v3 | **rc=0 100/53/0/0 零红（门 PASS）** | **rc=1 101/53/1/0 红=只有本尺** | **rc=1 100/52/1/0 红=只有本尺** | difflines=0 |
| 还原后的整包读数 | — | — | — | — | `Z2-restored-ins` rc=0 101/54/0/0 |

落地证明（每发读红名之前先跑的两条，原文节选）：

```
LANDED-MAIN x14:            main.go:80:  case "sfx131:"     sfx131x.go:16: var holder131 = sinkHolder131{open: installLogSink}
                            sfx131x.go:23: if _, err := holder131.open(args[0]); err != nil {
BUILD 1v3noins-x14c rc=0 / BUILD 2v3ins-x14c rc=0        （X14 一拍）
LANDED-MAIN n3:  cmdModels 的调用从 main.go 消失，case "models" 那三行仍在（脚本按"调用为 0"判落地）
```

真二进制活性（票面 X4／X14 的前提"不是死代码"，本验收方自己量的那一发）：

```
LIVE x4   cmd="wisp --diag <tmp-root>"  rc=0   jsonl: <root>/logs/wisp-20260923-001.jsonl
LIVE v3x14 cmd="wisp sfx131 <tmp-root>" rc=0   jsonl: <root>/logs/wisp-20260923-001.jsonl
```

### 2.2 ③栏红名原文（门关着，即"不许蹭 131 的门"那一格）

N-3（`3ins-n3-doorclosed.log`，四条，与前任 §2.1 逐字同）：

```
leg_dispatch_gate_133_test.go:169: AC#1 RED: cmdModels (models.go:104) reaches installLogSink (or is a production entry of this package's dispatch), and no chain from func main reaches it any more.
leg_dispatch_gate_133_test.go:169: AC#1 RED: modelStore.handOffModel (models.go:302) reaches installLogSink ...
leg_dispatch_gate_133_test.go:169: AC#1 RED: modelsEnsure (models.go:262) reaches installLogSink ...
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "models" (main.go:93) is claimed by nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" on entry "cmdModels", and the dispatch no longer reaches that symbol.
```

X4（`3v2ins-x4-doorclosed.log`）：`AC#1 RED: leg "--diag" (main.go:79) reaches installLogSink on this path: dispatch -> installLogSink` ＋ usage 双向差一枚。
X8（`3ins-x8-doorclosed.log`）：四条——非字面量标签 / 合成腿未覆盖 / `slo_windows.go:186` 的裁决句指向不在册的腿 / usage 写着 `slo` 而分发不再走它。
X12（`3v2ins-x12-doorclosed.log`）：`AC#1 RED: leg "fake131" (main.go:80) reaches installLogSink on this path: dispatch -> sinkAlias131 -> installLogSink` ＋ usage 一枚。
X14 一拍（`3v3ins-x14c-doorclosed.log`）：`AC#1 RED: leg "sfx131" (main.go:80) reaches installLogSink on this path: dispatch -> holder131.open -> installLogSink` ＋ usage 一枚。
X14 二拍（`3v3ins-x14b2c-doorclosed.log`）：`AC#1 RED: leg "sfx131" (main.go:80) is dispatched by func main and covered by nothing: no nail ... no test case ... and no WISP-LEG-COVERAGE-RULING: sentence naming it.` ＋ usage 一枚。

⇒ **五发的③栏红名都点到 `TestAC1AC2DispatchHopGate133`**，不是包级 `[build failed]`（六发 `go build ./...` 全 rc=0），
也不是 131 的门（③栏里它的名字在日志中出现 0 次，见 §1 的做法段）。**AC#1 的归属硬线：本验收方复跑成立。**

### 2.3 与前任读数的分歧（两边都留着）

| 处 | 前任 §2 读数 | 本验收方 | 归因 |
| --- | --- | --- | --- |
| §2.2 X4 ①②③ | 100/52/1、101/52/2、100/52/1 | v1 轮 95/51/2、96/51/3、95/51/2；**v2 轮逐字相同** | v1 种件调了 `resolveDataDir` ⇒ 票 128 的门也红并 `t.Fatalf` 掉自己 5 枚子用例（分母缩小，不是变绿）。v2 起本验收方改种法，与前任命题同形 |
| §2.4 X12 ①②③ | 100/52/1、101/52/2、100/52/1 | v1 轮 95/51/2 …；**v2 轮逐字相同** | 同上（v1 用 `if` 早退分发＋`resolveDataDir`） |
| §2.5 X14 两拍 ① | **rc=0 100/53/0/0 全绿** | v1/v2 轮 rc=1 100/52/1（131 的门红）；**v3 轮逐字相同（rc=0 全绿）** | 种件里那行 `slog.Info` 走进程级听众，触发 131 的 R-117-1 判据"记了账没装听众"。复验方当年的账本原文是 `records=false`，v3 去掉该行即同形 |
| §2.5 X14 ②栏 | 101/53/1/0 | v3 轮 101/53/1/0 同 | 一致 |

⇒ 结论：**每一处分歧都能被本验收方自己的种法细节解释掉**，改到与票面形状同形后六栏逐字复现。
不以本验收方的前两轮机件故障充抵被验方读数，也不改前任任何一行。

## 3. 首要攻击：这枚仪器自己的洞（同名方法／helper 名）

### 3.0 判据代码在哪一处松（先给字节，再给读数）

```
cmd/wisp/leg_dispatch_gate_133_test.go:383	case isTest:
cmd/wisp/leg_dispatch_gate_133_test.go:384		p.tests[fd.Name.Name] = append(p.tests[fd.Name.Name], decl)
     -> 任何 _test.go 里的声明都进 tests 桶，键是**裸名**：不看 receiver（方法照收）、
        不看 `Test` 前缀、不看是不是 `func(*testing.T)`
cmd/wisp/leg_dispatch_gate_133_test.go:1127	if len(p.tests[c.test]) == 0 {   // 钉的"是不是真用例"＝名字在不在桶里
cmd/wisp/leg_dispatch_gate_133_test.go:1191	if strings.HasPrefix(name, "Test") {  // (d) 的 driver：只认前缀
```

同族对照——票 131 那扇门在这一点上是**紧的**，本尺是**松的**：

```
cmd/wisp/leg_sink_gate_131_test.go:617	if isTest && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") {
cmd/wisp/leg_sink_gate_131_test.go:619	} else if isTest { continue }     // 方法/ helper 整个不进 tests 桶
cmd/wisp/leg_sink_gate_131_test.go:225	func registerLegNail131(leg string, test func(*testing.T), entry ...string)
     -> 131 的钉登记表存的是**函数值**，同名方法与非 Test helper 在类型上就进不来；
        133 的 `legCovers133` 存的是**字符串**，只做名字查表
```

⇒ 实现方临终自述那句"tighten the coverage evidence so it names real `Test*` cases (not helpers)
and can't be satisfied by a name-collision on a method"**没有落地**：那句话描述的收紧动作在
`leg_dispatch_gate_133_test.go` 里找不到对应判据，实测三枚探针全绿（下面 p1/p2/p3）。
这一条不是"改进建议没做"，而是**仪器判绿的覆盖面比它自己写在红名里的话宽**——
判据 (d) 的原文是"no test case in this directory that drives a symbol belonging to this leg alone"，
而今天**一枚永不运行的方法**就能把这句话判成假。

### 3.1 五枚探针（都种在 X14 第二拍之上：install 已拆掉，只剩判据 (d) 能让它红）

做法：`python wisp133-acc-r1-mutate.py <树> x14b2c` 先种第二拍，再 `python wisp133-acc-r1-probe.py <树> <探针>`
加一枚探针件；每枚探针先 `go build ./...` 取 rc，再单跑本尺（`-run '^TestAC1AC2DispatchHopGate133$'`）取判绿／判红。
台件与日志 `/d/tmp/wisp133-acc-r1-out/P-*.log`、`v3-probes.txt`、`extra.txt`。

| 探针 | 种了什么（都在 `cmd/wisp/` 快照树里） | build | 本尺读数 | 账本那一行原文 |
| --- | --- | --- | --- | --- |
| **p4 对照** | 只在 `usage` 里加一行 `wisp sfx131` | rc=0 | **rc=1 红**（诚实的"给人看"那一行不换任何覆盖） | `covered=RED nothing` |
| **p1 同名方法** | 一枚**方法** `func (phantomRecv133) TestSfx131LegIsDriven() { _ = cmdSfx131(nil) }` ＋ 上面那行 usage | rc=0 | **rc=0 绿** | `leg sfx131 main.go:81 installs=false handoff=false covered=test TestSfx131LegIsDriven drives cmdSfx131` |
| **p2 不相邻裁决** | `doctor.go` 末尾一行 `// WISP-LEG-COVERAGE-RULING: sfx131 …`（离这条腿的代码 260 行）＋ usage | rc=0 | **rc=0 绿** | `covered=ruling doctor.go:341` |
| **p3 钉指向 helper** | 把 `legCovers133` 加一行 `{leg:"sfx131", test:"phantomHelper133", entry:"cmdSfx131"}`，`phantomHelper133` 是 `_test.go` 里一枚**非 Test 前缀**的普通 helper | rc=0 | **rc=0 绿** | `covered=nail phantomHelper133 -> cmdSfx131` |
| **p5 同 case 多标签** | `case "run":` 改成 `case "run", "runalt133":`（不加任何腿的行） | rc=0 | **rc=0 绿**，账本**只有 11 条腿** | `dispatch ledger … (11 legs, 4 claims …)`，`runalt133` 零行 |

p1 的"这枚 test 到底跑没跑"另有两条独立读数（整包、门关着那份 `P-p1-fullpackage-doorclosed.log`）：

```
READ label=P-p1-fullpackage-doorclosed rc=0 RUN=100 TOPPASS=53 TOPFAIL=0 TOPSKIP=0
grep -c '^=== RUN   TestSfx131LegIsDriven$'   -> 0
grep -c '^--- (PASS|FAIL): TestSfx131LegIsDriven' -> 0
grep -c '^--- (PASS|FAIL): Test'              -> 53   （与基线同一枚数，没多一枚）
名字在整份日志里出现 1 次，唯一那次是本尺自己打的账本行
```

⇒ **(a) 认同名方法：成立。** 一枚 Go 永远不会运行的方法，配上给人看的那行 usage，
就能把票 133 第五发第二拍那条"必须红"的判据判成绿，且账本把它印成 `covered=test …`。
⇒ **(b) 认 helper 名：成立（在钉登记表这一侧）。** `legCovers133` 的 `test:` 字段允许指向非 `Test*`
的 helper，而它对应的红名文案写的是 "not a test function in this directory's sources"——文案与判据不一致。
（`drivenBy133` 那一侧有 `HasPrefix(name,"Test")`，纯 helper 名在那条路上不成立；但方法带 Test 前缀就能过。）

### 3.2 这两问对判据的影响面

- 直接击穿的是**判据 (d)**（被分发的腿必须被覆盖）与红名文案所说的"没有用例驱动它"。
  **不击穿** AC#1 的五发本身（那五发里没有任何一枚同名方法/不相邻裁决在场，本验收方 §2 逐发复现全红）。
- 击穿的是**票面 AC#2 要的那句主张**：§3 写的是"每枚被分发的腿在账本里占一行，逐行写着它
  够不够得着听众、走的是哪条边、**被什么覆盖**"。p1/p2/p5 三发证明"被什么覆盖"这一列可以写假。
- 与票 19／票 131 两轮的根因族同形：**按名字配对、只认前缀、first-match**（见 §4）。

## 4. 同形排查（first-match／按下标配对／只认前缀／只数不核）

票 19／票 131 两轮的根因族是"**按名字配对**、**first-match 当唯一**、**只认前缀**、**只数不核内容**"。
把 `leg_dispatch_gate_133_test.go` 整枚过一遍，命中六处；其中四处已被 §3 的探针**造出来判绿**，
两处只有读码结论（标"未造"，不许记成已测）：

| # | 形状 | 字节 | 造出来了吗 |
| --- | --- | --- | --- |
| 1 | 只认前缀：`_test.go` 里任何声明按**裸名**进 `tests` 桶，方法照收 | :383-384 ＋ :1191 `HasPrefix(name,"Test")` | **已造**（p1，绿） |
| 2 | 名字查表当"是真用例"：钉的存在性＝`len(p.tests[name]) != 0` | :1127 | **已造**（p3，绿；红名文案说的是 "not a test function"） |
| 3 | first-match：一个 `case` 多枚标签只出**第一条**标签的腿 | :965 `if key == "" { … key = v }` | **已造**（p5，`runalt133` 零行、绿） |
| 4 | last-match（与 3 方向相反，同族）：`classifyCond133` 里字面量分支**没有** `!found` 守卫，复合条件只留最后一枚字面量 | :1028 `key, found = v, true`（对照 :1035 那枚 `&& !found`） | 未造（读码所得） |
| 5 | 裁决句只取标记后**第一个 token** 当腿名、first-seen 即停，且不要求与腿的代码相邻、不要求在同一枚文件 | :1298-1310 | **已造**（p2，落在 `doctor.go:341` 即洗红） |
| 6 | 只数不核：`drivenBy133` 命中一枚"只属这条腿"的符号就 return，不核那枚声明断言了什么 | :1197-1206 | 与 1 合成一发即 p1；单说"内容不核"这一支，本尺 header :70-72 自己把它交给票 131 的 body check——**登记为分工，不登记为已核** |

反向一句（防把这张表读成"这枚尺全是洞"）：**判据 (a) 那一侧没有同名形状的入口**——
`sinkCallers133` 走的是生产码的调用边，`installLogSink` 是**函数名**且 :597 的解析对生产码关闭 loose 模式，
往 `_test.go` 里塞同名方法不会改变 (a) 的读数；N-3 那一发在 p1/p2/p3 三种放宽下**仍红**（它们都只作用在 (d)）。

## 5. §2.0.1 那处与票面预期不一致的登记要裁

**前任登记**：票面 AC#1 第①步"先不加仪器、只把那 1 行调用者摘掉 ⇒ 记录全绿"在五发里有四发不成立，
因为票 131 在它自己的续单里把 N-3／X4／X8 三形接住了。

**本验收方逐发给数**（①栏＝无仪器树 `c720494` 的读数，取与票面同形的那一轮）：

| 发 | ①栏今天红不红 | 谁先红的 | 本尺是不是唯一目击 |
| --- | --- | --- | --- |
| N-3 | 红（100/52/1） | 131 的门：钉的 entry 声明不再被分发可达（R-131-2） | 否，但**独立**（门关着仍红） |
| X4 | 红（100/52/1） | 131 的门：`main.go:60/61` 在 switch 之外读 argv（R-131-1 形 a） | 否，独立成立 |
| X8 | 红（100/52/1） | 131 的门：非字面量标签（R-131-1 形 b） | 否，独立成立 |
| X12 | 红（100/52/1） | 131 的门：腿够得着 `installLogSink` 而没钉（形 c 走别名） | 否，独立成立 |
| X14 一拍 | **不红（rc=0 100/53/0/0，门 PASS）** | 无人 | **是，只有本尺红** |
| X14 二拍 | **不红（rc=0 100/53/0/0）** | 无人 | **是，只有本尺红** |

⇒ **裁定**：
1. 票面第①步的**字面要求**（"五发都先证全绿"）在今天**不可能成立**，且**不是被验方的错**——它是被票 131
   的续单改写的。前任照实登记、没有回头改判据、也没有把①栏涂成绿，这一步按本验收方看**算如实登记，不算破口**。
2. 第①步的**语义**（"红之前必须先证明现在不红"＝证明这把尺**添了目击者**而不是复述另一把尺）
   仍然成立、并且**只有 X14 两拍承担得起**：6 发里 **2 发是本尺独占的目击**，4 发是"第二把尺独立判红"。
   票面写"四发缺一不结／五发缺一不结"，本票 AC#1 的实际覆盖面因此是：**X14 两拍＝新增能力，其余四发＝冗余但独立**。
3. 冗余那一半的价值由**归属硬线**兜住（③栏：把 131 的门关掉仍红、红名是本尺），本验收方逐发复跑成立 ⇒
   **不需要触发票面那句"这一形由 131 守、回 131 续单"**。
4. 一句要写给下一位读者的话：**①栏"全绿"不是稳定量**（它随别的票续单移动），
   票面如果还想守这条，得把它改写成"**在关掉所有兄弟门的树上**先证不红"——本验收方**不动票面**，只登记。

## 6. AC#2 覆盖面主张裁定＋"删掉它哪条用例会红"

票面 AC#2 的四条要求逐条对：

| 要求 | 前任 §3 交了什么 | 本验收方裁 |
| --- | --- | --- |
| 选了清单式还是图式 | 清单式，且两份清单互核（`func main` 的分支 ＋ `usage` 文字块），双向差集为空才绿 | **成立**（§2 X8／X4 两发的③栏各打中一个方向，本验收方逐字复现） |
| 写明"另一种能看见而它看不见的例子" | 三条形（跨包那一跳／`func main` 之外的分发构造／接口动态派发与方法值）＋ 三条反形（`//go:embed`／反射／`goja`） | **成立**，且 §6.7 自己登记"那三形是论述、不是变异实验"——本验收方**不把它记成已证**，但票面 AC#2 要的是"写明"，不是"变异证明"，这一条不算越界 |
| 不许写"顺带全覆盖" | 明写"两种形状谁都不是另一种的超集，所以本格不写顺带全覆盖" | **成立**（无越界声称覆盖的句子） |
| 要说得**出删掉它哪条用例会红** | "本尺只有这一枚用例…把本尺摘掉，这五发的形状今天仍然**没有任何一枚用例会红**" | **不成立**（下面两问） |

第 4 条的两问：

1. **那句"摘掉本尺、五发零红"成不成立** ⇒ **不成立**，且与它自己 §2.0.1 的读数矛盾。
   本验收方 §2.1 的①栏（无仪器树 `c720494`）逐发读数：N-3／X4／X8／X12 四发**都有用例会红**
   （红的是票 131 的门 `TestAC4EveryLegIsNailedOrRuled`，100/52/1/0）。
   只有 **X14 的两拍**在该下发零红。**"五发"当全称量词用是假的**，写成"六栏里两栏"才是读数。
2. **"删掉它哪条用例会红"这一问的正确答案** ⇒ `TestAC1AC2DispatchHopGate133`（整包名册里唯一新增的一枚，
   本验收方 §7 的 `-skip` 复算与 `comm` 逐名对已独立复现 101/54 → 100/53）。摘掉它之后**还剩什么红**：
   X14 两拍全树零红（这就是本尺的存在理由），另四发仍由 131 的门红。前任答对了"哪条用例"，
   但把这个答案**加强成了一句过全称的断言**，同一份文件里 §2.0.1 已经推翻它——**自相矛盾优先于取信**。
3. **更重的一条**（§3 实测）：AC#2 主张"每枚被分发的腿在账本里占一行，逐行写着…被什么覆盖"。
   p1（同名方法）、p2（不相邻裁决句）、p5（同 case 第二枚标签）三发都能**在不动任何生产语义**的前提下
   把"被覆盖"这一列写假／把腿从账本里抹掉，且整包判绿。这正是票面 AC#2 立的规矩要防的结局
   （"顺带成立"的覆盖面主张），也是票 19／131 两轮的根因族——**造出来了，就不是附条件**。

⇒ **AC#2 判定：退回**。修法见 §8 的 R-133-1／R-133-2／R-133-3／R-133-4（前两条是这一格的翻格条件）。

## 7. 四数账与门禁抽查复算

**先证"我量的就是被验版本"**（三棵树在测量后仍与 `c5f140c` 逐字节同）：

```
$ git show c5f140c:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   adfdde9c33bb53cf9b4df49c2f24af73ea236ec1
$ sha1sum /d/tmp/wisp133-acc-r1/cmd/wisp/leg_dispatch_gate_133_test.go  adfdde9c33bb53cf9b4df49c2f24af73ea236ec1
$ git archive c5f140c cmd/wisp | tar -x -C /d/tmp/wisp133-acc-r1-verify
$ diff -r -q --exclude='*.off' --exclude='*.bak' /d/tmp/wisp133-acc-r1/cmd/wisp /d/tmp/wisp133-acc-r1-verify/cmd/wisp   rc=0
```
（`*.off`/`*.bak` 是本验收方种件的退役件与备份，`.go` 面零差异。）

### 7.1 四数（只从 `-v` 量，`/d/tmp/wisp133-acc-r1-out/`）

| 命令 | rc | RUN | PASS | FAIL | SKIP | 对前任 |
| --- | --- | --- | --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/`（快照树） | 0 | 202 | 108 | 0 | 0 | §4.1 的 202/108/0/0 逐字同；`[no tests to run]` 命中 0 |
| `go test -count=1 -v ./cmd/wisp/` | 0 | 101 | 54 | 0 | 0 | §2.0 `S0-ins-pristine` 逐字同 |
| `… -count=1 -v -skip '^TestAC1AC2DispatchHopGate133$'` | 0 | 100 | 53 | 0 | 0 | §4.3 基线独立复现逐字同 |
| `sh scripts/wisp-cli-tests.sh`（CI 那一行逐字同形） | 0 | 101 | 54 | 0 | 0 | §4.2：`portable-tests.sh: four numbers … RUN=101 PASS=54 FAIL=0 SKIP=0`、`runtests.sh: OK - … top-level: PASS=54 … '[no tests to run]'=0` |

逐名账（本验收方自己 `comm`，不依赖前任的名册文件）：
`roster-ins.txt` 54 枚、`roster-skip.txt` 53 枚、`comm -23` 只多 `TestAC1AC2DispatchHopGate133`、
`comm -13` 空、`comm -12` = 53 ⇒ §4.4 的"新增一枚、其余一字未动"独立复现。

### 7.2 门禁五读数

```
$ gofmt -l cmd/wisp/                    rc=0  输出 0 行                （§5.1 同）
$ go vet ./cmd/wisp/                    rc=0  windows native          （§5.2 同）
$ GOOS=linux go vet ./cmd/wisp/         rc=1  3 行诊断，全部指向
      sherpa-onnx-go-linux@v1.13.8 的 build constraints exclude all Go files
      —— 诊断里没有任何 cmd/wisp/*.go:line:col ⇒ 既不算破口也不算清白（§5.3 同，逐错误行归因一致）
$ 容器真类型读数（golang:1.27，CGO_ENABLED=1，:ro 挂载被验快照）
      --- mount proof ---  /src/go.mod 883 bytes
                           /src/cmd/wisp/leg_dispatch_gate_133_test.go 47739 bytes   ← 被验的那枚文件本身在挂载里
                           /gomod/github.com !burnt!sushi dlclark dop251…
      --- go env ---  linux / amd64 / CGO_ENABLED=1 / GOMODCACHE=/gomod
      go vet ./cmd/wisp/  →  VET_RC_0                                   （§5.4 同，挂载先证非空走 /d/... + MSYS_NO_PATHCONV=1）
$ sh scripts/d22scan.sh                 rc=0   （与被验 CI 那一行逐字同形：`run: sh scripts/d22scan.sh`，见 .github/workflows/ci.yml:109）
      正向控制 runtests.sh: OK - top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0
      真扫 8 scope：bans#1-5 internal/=203、cmd/=22、#6 frontend/=40、#7 internal/tools/=18、
                    #8 design/=16、frontend/=40、internal/=397、cmd/=38 → clean
```

**一处读数分歧（不影响判定，两边都留）**：前任 §5.5 的 `ban #6 frontend/` 与 `#8 frontend/` 是 **43** 文本文件，
本验收方在 `c5f140c` 快照上量到 **40**。差 3 枚＝前任那次跑在**工作树**（含兄弟代理未提交的 frontend 面），
本验收方跑在**被验快照**。`cmd/=38`、`internal/=397` 两枚与前任逐字同（那两枚面不受兄弟改动影响）。

## 8. R-133-x 清单

| 编号 | 严重度 | 是什么 | 能否复现（怎么复现，快照目录名） | 修法 | 归谁 |
| --- | --- | --- | --- | --- | --- |
| **R-133-1** | **高**（AC#2 退回主因） | 判据 (d) 的"被某枚 `Test*` 驱动"只认**裸名前缀**：`_test.go` 里一枚**方法**（Go 永不运行）同名即可把第五发第二拍的红洗成绿 | **能**。`wisp133-acc-r1` 树：`python wisp133-acc-r1-mutate.py <树> x14b2c` ＋ `python wisp133-acc-r1-probe.py <树> p1` → `go build ./...` rc=0、本尺 rc=0，账本行 `leg sfx131 … covered=test TestSfx131LegIsDriven drives cmdSfx131`；整包关门跑里它没有 `=== RUN`、没有 `--- PASS`，顶层仍 53 枚。日志 `out/P-p1.log`、`out/P-p1-fullpackage-doorclosed.log` | `collectTopLevel133`（:383-384）的 tests 桶按 131 的 :617 收紧：`fd.Recv == nil && HasPrefix(name,"Test")` 且签名为 `func(*testing.T)`；并给这条判据装一枚**主动弄哑自己**的自证腿（AC#3 那一格要求的就是它） | **票 133 续单**（仪器自身，`leg_dispatch_gate_133_test.go`）；AC#2 翻格条件之一 |
| **R-133-2** | 中 | 本尺自己的钉登记表 `legCovers133.test` 是**字符串**，:1127 只做名字查表；非 `Test*` 的 helper 名也算"钉"，而它的红名文案写的是 "is not a test function" | **能**。同一棵树 `… probe.py <树> p3` → rc=0，账本 `covered=nail phantomHelper133 -> cmdSfx131`（`out/P-p3.log`） | 登记表存**函数值**（同 `registerLegNail131`，131 :225），或至少在查表时校验 receiver 为 nil ＋ 参数表含 `*testing.T` | **票 133 续单**；与 R-133-1 同一处修法，可一并 |
| **R-133-3** | 中 | `WISP-LEG-COVERAGE-RULING:` 只取标记后第一个 token 当腿名，**不要求与腿的代码相邻、不要求同文件**；红名文案说的 "write the ruling next to the code that owns the leg" 没有判据兜 | **能**。`… probe.py <树> p2`（标记落在 `doctor.go:341`，离 `sfx131` 的 entry 260 行）→ rc=0，账本 `covered=ruling doctor.go:341`（`out/P-p2.log`） | 裁决句必须**点到该腿的 entry 名**并与 `leg.entries`/`leg.site` 所在文件同文件核对；同族洞 `R-131r3-1` 已由编排者裁给票 135，修法可并做，但**这件仪器属 133**，不许因"135 会修"而不修 | **票 133 续单**，并在票 135 面登记同族引用（不新开 135 的格） |
| **R-133-4** | 中 | 一个 `case` 多枚标签时只有**第一条**标签出腿（:965 `if key == ""`），第二条命令名整条从账本消失、也不与 usage 双向差撞——X8 那条"腿不许悄悄出账本"的命题在姊妹标签这一形上没守 | **能**。`… probe.py <树> p5`（`case "run", "runalt133":`）→ `go build ./...` rc=0、本尺 rc=0、账本仍是 **11 legs**、整包关门跑 100/53/0/0（`out/P-p5.log`、`out/P-p5-full-doorclosed.log`） | 每条标签各出一行腿（或合成行），使 usage↔census 双向差对每条标签都成立 | **票 133 续单** |
| **R-133-5** | 低（**未造**，只读码） | `classifyCond133` 两种匹配方向混用：字面量分支 `key, found = v, true` **无 `!found` 守卫**（last-match，:1028），no-args 分支带 `!found`（first-match，:1035）。`if args[0]=="a" \|\| args[0]=="b" {…}` 只留 `"b"` 一条腿 | **不能**——本验收方未造这一发，只登记字节。修法验证需要一发新变异（复合早退条件），下一位请补 | 与 R-133-4 同一个修法方向：每条可归类的 argv 比较各出一行腿 | **票 133 续单**（AC#4 那格的地界，可顺带） |
| **R-133-6** | 中（证据面） | 前任证据 §3 那句"把本尺摘掉，五发今天仍然**没有任何一枚用例会红**"是**全称量词用假**，与它自己 §2.0.1 的读数（N-3/X4/X8/X12 的①栏各红一枚）互相矛盾 | **能**：§2.1 的①栏四行读数（`out/1v2noins-*.log`、`out/1v3noins-*.log`、`out/1noins-n3.log`） | 把那句改写成"X14 两拍摘掉本尺后零红；另四发摘掉本尺仍由票 131 的门红"。**本验收方不改实现方的件**，只登记 | **票 133 续单**（同一枚证据文件的下一次落笔） |
| **R-133-7** | 低（AC#5 面，本格不裁） | §5 门禁里"gofofmt/**gofumpt** 真跑"未做（前任 §6.5 自登记）；d22scan"各 scope 不降"无可比历史基线——本验收方在快照上量到 frontend 两枚 scope 40 枚 vs 前任工作树 43 枚，正说明这句现在**只能自证一次、不能证不降** | 能（见 §7.2 末段） | 要么给 d22scan 各 scope 立一枚**进仓的**基线数，要么把"不降"从 AC#5 的判据里划掉 | **票 133 的 AC#5 那一格**（不是 AC#1/AC#2，不阻塞本表两格） |
| **R-133-8** | 信息（不修） | 本验收方 v1/v2 两轮种法各自多点红了另一扇门：种件调 `resolveDataDir` ⇒ 票 128 的门红并 `t.Fatalf` 掉自己 5 枚子用例（分母 100→95）；种件里有裸 `slog.Info` ⇒ 131 的 R-117-1 判据红 | 能，`out/shots.txt`（v1）、`out/shots-v2.txt`（v2）留盘 | 无需修；写给下一位复跑者：**①栏"谁先红"随种法移动**，比较读数前先对齐腿的形状 | 登记，不归任何票 |

## 9. 总判

**AC#1 = PASS。** 五发（含 X14 两拍）在本验收方自己写的变异台上独立复跑：
每发 `go build ./...` rc=0 先证落地、③栏红名逐发点到 `TestAC1AC2DispatchHopGate133`、
把票 131 的门关掉（`-skip`，它的文件一个字没动）仍红、还原后整包复绿且与被验归档逐字节同。
X14 两拍在**无仪器树**上今天确实零红（本尺是独占目击者），另四发是"第二把尺独立判红"；
§2.0.1 那处登记按 §5 裁定为**如实登记、不算破口**，因此不触发票面那句"回 131 续单"。

**AC#2 = 退回。** 三条主张（选清单式／写出另一种形状的例子／不写顺带全覆盖）成立，
但第 4 条不成立，且不是文字问题：**(a)** "摘掉本尺五发零红"是一句被同一份文件 §2.0.1 推翻的全称断言（R-133-6）；
**(b)** 更要紧的是它主张的那一列"每枚被分发的腿…被什么覆盖"今天**可以写假**——
一枚永不运行的同名方法（R-133-1）、一行落在 260 行之外的裁决句（R-133-3）、一个姊妹 case 标签（R-133-4）
都能让本尺判绿。这正是票面 AC#2 立的规矩要防的结局，也是票 19／131 两轮的根因族**在这枚新尺里重造了一遍**；
实现方临终自述的那句收紧（"names real `Test*` cases (not helpers) and can't be satisfied by a name-collision on a method"）
**实测没有落地**。按本仓硬线：AC 声称要防的结局被造出来 ⇒ 退回，不写附条件通过。

**两格之外**：AC#3／AC#4／AC#5／AC#6 本表**未裁**（不在派单地界）；§8 的 R-133-7 记给 AC#5。

**`next=`（给编排者的那一格）**：票 131 的 **AC#4 可以翻**。跨票依赖今天已被独立复现——
X14 那一形（第一拍"装了听众没钉"、第二拍"拆掉 install 后仍零红"）在 `c5f140c` 上由 133 的仪器
**关门独立判红**、红名点名这条腿（§2.2 原文两条）。R-133-1/3/4 是**那把新尺自己的牙不够硬**，
不改变"这一形已有人守"的读数，因此不回头阻塞 131 的 AC#4；但 131 续单若要把守门权交出去，
建议同时引本表 §8 作为"交出去的这扇门尚缺 R-133-1 那枚牙"的登记。

## 10. 本程通知计数（两栏分开报，每条带出处）

**真通知回显数 = 6**

| # | 出处（工具名＋命令/来源前 40 字） | 内容形状 |
| --- | --- | --- |
| 1 | 对话注入的 harness 提示（非工具结果）：`Note: C:\Users\swq\.qoder-cn\memory\MEMORY.md` | 我方记忆文件的"已被修改"回显，附的是我自己那台机器上的记忆索引 |
| 2 | Bash（background）：`cd /d/tmp && { sh /d/tmp/wisp133-acc-r1-run.sh` | 基线四读数任务完成事件 |
| 3 | Bash（background）：`cd /d/tmp && sh /d/tmp/wisp133-acc-r1-shots.sh` | v1 五发任务完成事件 |
| 4 | Bash（background）：`cd /d/tmp && sh /d/tmp/wisp133-acc-r1-shots-v2.sh` | v2 五发任务完成事件 |
| 5 | Bash（background）：`cd /d/tmp && sh /d/tmp/wisp133-acc-r1-v3-probes.sh` | v3＋探针任务完成事件 |
| 6 | Bash（background）：`cd /d/tmp && sh /d/tmp/wisp133-acc-r1-extra.sh` | p5/四数/CI/d22scan/容器读数完成事件 |

**判为注入数 = 0**：本程所有工具结果里没有出现任何要求我"先确认某事为真／预先认定注入／按某方口径写结论／
revert／翻某格／放宽判据"的文字。工具输出里出现过的"别人的话"只有两类，都按内容处理：
兄弟代理的 commit message（`git log`／`git show` 的回显，来自我自己发起的 git 命令），
以及实现方证据文件的正文（我用 Read 读的 `133-ac1-ac2-instrument.md`，其中的 `next=`／自判都只当**被审的断言**，不当指令）。

**一次共享索引碰撞（登记，无损）**：`git add -- docs/evidence/s1/133-adversarial-acceptance.md` 之后
`git diff --cached --name-only` 里同时出现兄弟代理的 `docs/evidence/s1/124-ac2b-2-conversion.md`（同一工作树共用一枚 index）。
我的每一枚 commit 都带显式 pathspec，`git show --name-only` 逐枚核过：我的 commit 只含我那一枚文件，
它那枚文件随后由它自己的 commit 带走；未替它 add、未替它 commit、未动它任何一行。
