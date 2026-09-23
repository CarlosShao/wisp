# 票 133 AC#1 / AC#2 —— 第二把尺：分发跳本身（`cmd/wisp/leg_dispatch_gate_133_test.go`）

**执行方**：`worker-ticket133-ac1`（本文件唯一作者；实现方，只做 AC#1 与 AC#2 两格）
**开工锚定**：`git rev-parse --short HEAD` = **`c720494`**（分支 `dev`，`git status --porcelain` 内
`cmd/wisp/**` 为空 ⇒ 被验版本 = 这枚 sha，脏文件只有兄弟代理的 `docs/reports/*` 与新增票 136，一字未动）
**日期**：2026-09-23（开件 17:46 +08，逐节落盘，每裁一发 commit 一次）
**票面判据**：`.scratch/wisp/issues/133-no-instrument-sees-the-dispatch-hop-deleting-one-caller-keeps-76-tests-green.md`
的 AC#1（五发，含 09-23 17:3x 追加的第 5 发 X14 与其"第二拍也必须红"）与 AC#2（覆盖面主张）

## 0. 仪器是什么，以及它为什么不会与票 131 那扇门共用一根筋

新增一枚用例 `TestAC1AC2DispatchHopGate133`（`cmd/wisp/leg_dispatch_gate_133_test.go`，
自带一整套 AST 读取，**不调用** `leg_sink_gate_131_test.go` 里的任何函数、不读它的
`legNails131`、不读它的账本），判据五条：

| 判据 | 一句话 | 五发里谁落在它上面 |
| --- | --- | --- |
| (a) 能力向内读 | 凡够得着 `installLogSink` 或 `models.WireDownloading` 的符号，必须从 `func main` 可达 | **N-3** |
| (b) 装听众就得有钉 | 被分发的腿够得着听众 ⇒ 本尺自己的钉清单里必须有它 | **X4 X12 X14 第一拍** |
| (c) 分发目标必须是本包声明的函数 | 解析不出来的调用不是叶子，是"本尺瞎了"，红着说 | （X8 的合成腿顺带命中） |
| (d) 被分发的腿必须被覆盖 | 钉 / 某个 `Test*` 打到只属这条腿的符号 / 生产码里一行 `WISP-LEG-COVERAGE-RULING:` | **X14 第二拍** |
| (e) 腿清单与 usage 块双向差集为空 | 第二份清单读的是给人看的那段文字 | 四发顺带命中 |

**为什么 (a)(d) 是票 131 那扇门看不见的形状**：131 的门的每一条判据都以
"这条腿自己够不够得着 installLogSink" 为前提；把分发那一跳摘掉，那个前提自己就塌了
（`install=false` ⇒ 无义务 ⇒ 绿）。本尺的 (a) 反过来从听众往里读、(d) 干脆不看听众，
所以同一发变异在它这里没有"降级为无义务"这条路。

## 1. 基线（`c720494` 纯净快照，未加仪器；先证明"今天不红"）

```
$ git archive c720494 | tar -x -C /tmp/wisp133-ac1-a        # 再拷 third_party/sherpa-onnx/*.dll
$ PATH=<快照>/third_party/sherpa-onnx:$PATH go test -count=2 -v ./cmd/wisp/
RUN=200  --- PASS=106  --- FAIL=0  --- SKIP=0     ok 104.000s   rc=0
$ sh scripts/wisp-cli-tests.sh                # CI 逐字同形
portable-tests.sh: four numbers (all from -v output): === RUN=100  --- PASS=53  --- FAIL=0  --- SKIP=0
runtests.sh: OK - ... top-level: PASS=53 FAIL=0 SKIP=0, === RUN=100, '[no tests to run]'=0   rc=0
$ grep -c 'AC#4 RED' <baseline log>            → 0
```
时间戳：`date -u` = 2026-09-23 09:4x–09:52z（本地 17:5x–18:0x +08）。

## 2. AC#1 五发逐格

做法逐字按票面：①**不加仪器**、只做该变异（树 `c720494`）；②装上仪器、同一发变异，
**门开着**（整包）与**门关着**（`-run` 排除 `TestAC4EveryLegIsNailedOrRuled`，**不改它的文件**）各一份；
③还原（还原后与纯净树 `diff -r` 反查）。变异只落在 `/tmp` 快照树里。

**三枚树**（都从 sha 现取，DLL 一并拷进 `<树>/third_party/sherpa-onnx`，缺它就是 `0xc0000135`，
第一轮我就是这样把 8 枚真二进制用例读成红的——那是环境故障，不是发现，全部作废重跑）：

```
/tmp/wisp133c-base        = git archive c720494      （无本票仪器）
/tmp/wisp133c-ins         = git archive 3d43c3f      （有本票仪器，变异在此重放）
/tmp/wisp133c-ins-pristine= git archive 3d43c3f      （从不变异，还原反查的对照）
驱动脚本 /tmp/wisp133-run-v2.sh，变异器 /tmp/wisp133-edit.py（每处改动先 assert 锚点唯一、
改完 grep -n 落地行 + go build ./... rc=0 才读红名；还原 = 从该树自己的 sha 重取 main.go，
种进去的腿改名成 *.off，不删任何文件）
读数日志 /tmp/wisp133c-out/*.txt，逐行账 /tmp/wisp133c-out/summary.txt
```

**门关着的做法**：`go test -list '.*'` 取全部顶层用例名，剔掉
`TestAC4EveryLegIsNailedOrRuled`，拼成 `-run '^(其余全部)$'`——一票 131 的门不注册即等于不存在，
它的文件一个字没动。

### 2.0 S0：两棵树各自的未变异读数（先证明仪器自己不是恒红，也证明门无关）

```
LABEL=S0-base-pristine     TREE=/tmp/wisp133c-base MODE=full        rc=0  RUN=100 TOPPASS=53 TOPFAIL=0 TOPSKIP=0
LABEL=S0-ins-pristine      TREE=/tmp/wisp133c-ins  MODE=full        rc=0  RUN=101 TOPPASS=54 TOPFAIL=0 TOPSKIP=0
LABEL=S0-ins-doorclosed    TREE=/tmp/wisp133c-ins  MODE=doorclosed  rc=0  RUN=100 TOPPASS=53 TOPFAIL=0 TOPSKIP=0
   started=2026-09-23T10:45:30Z / 10:46:31Z / 10:47:32Z（本地 18:45 / 18:46 / 18:47 +08）
```
⇒ 未变异的树上，仪器绿；把 131 的门关掉，仪器也绿（它不依赖那扇门才拿得到的事实）。
`RUN 101 / TOPPASS 54` 与 `100 / 53` 的差就是本票新增的那一枚用例，逐名账见 §4。

### 2.0.1 ⚠ 票面那句"今天一条不红"在 `c720494` 已经不成立（读数，不是感觉）

五发的第①栏（**不加仪器**、只做变异）没有一发是"全绿"：

```
LABEL=N3-1-noinstrument rc=1  TOPFAIL=1  reds: TestAC4EveryLegIsNailedOrRuled   131red=2
LABEL=X4-1-noinstrument rc=1  TOPFAIL=1  reds: TestAC4EveryLegIsNailedOrRuled   131red=2
LABEL=X8-1-noinstrument rc=1  TOPFAIL=1  reds: TestAC4EveryLegIsNailedOrRuled   131red=1
```
原因写在票 131 的续单里：它给每枚钉加了 **entry 声明**（R-131-2），于是 N-3 摘掉分发那一跳时
它自己的门就红了；X4/X8 是它 09-23 补的三形判据（越界读 argv、非字面量标签）。
**这不是本票的失败**：票面 AC#1 要的是"本票仪器能把这些形做成红"，而归属硬线要的是
"把 131 的门关掉仍然红、红名仍是本票的仪器"——§2.1..§2.5 的③栏逐发成立
（三门全绿的时代过去了，今天反而更要紧的是**第二把尺独立判红**，读数就在下面）。
本格照实登记"①栏不为全绿"，不拿它当本票的减分，也不改读数。

### 2.1 N-3 摘掉 `cmdModels` 的那一行调用者

```
① c720494 + 变异，无仪器        rc=1  100/52/1/0  reds: TestAC4EveryLegIsNailedOrRuled（131 的门，非本票）
② 3d43c3f + 变异，门开着        rc=1  101/52/2/0  reds: TestAC1AC2DispatchHopGate133, TestAC4EveryLegIsNailedOrRuled
③ 3d43c3f + 变异，门关着        rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133      ← 本票仪器独立判红
④ 还原                          rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
build：N3-1-build-base rc=0、N3-build-ins rc=0、N3-restore-build rc=0
```
③栏原文（`/tmp/wisp133c-out/N3-3-doorclosed.txt`）：

```
AC#1 RED: cmdModels (models.go:104) reaches installLogSink (or is a production entry of this package's dispatch), and no chain from func main reaches it any more.
AC#1 RED: modelStore.handOffModel (models.go:302) reaches installLogSink ... no chain from func main reaches it any more.
AC#1 RED: modelsEnsure (models.go:262) reaches installLogSink ... no chain from func main reaches it any more.
AC#1 RED: leg "models" (main.go:93) is claimed by nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" on entry "cmdModels", and the dispatch no longer reaches that symbol. This is the dispatch hop itself: the case still compiles and still reads a disk, but the command it names is not wired to the code it claims to cover.
```

### 2.2 X4 早退 `if` 分发的腿（`--diag`）

```
① c720494 + 变异，无仪器  rc=1  100/52/1/0  reds: TestAC4EveryLegIsNailedOrRuled
② 门开着                  rc=1  101/52/2/0  reds: TestAC1AC2DispatchHopGate133, TestAC4EveryLegIsNailedOrRuled
③ 门关着                  rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133      ← 红名点到 --diag
④ 还原                    rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
```
③栏原文：

```
AC#1 RED: leg "--diag" (main.go:79) reaches installLogSink on this path: dispatch -> installLogSink
A listener installed on a dispatched leg has to have a nail in this gate's registry, or the block can be deleted in silence, which is the reading ticket 133 was filed for. Fix: write the case, then add {leg: "--diag", test: TestYourCase, entry: "cmdDiag131"} to legCovers133.
AC#1/#2 RED: leg "--diag" is dispatched by func main and documented in no line of the usage block: an operator cannot find it, and this gate has no second reading of the census to check it against.
```
⇒ 与 131 的做法不同：131 把这条 `if` 读成"越界读 argv"判红；本尺把它读成**一条腿**，
于是它当场欠一枚钉。同一发，两种前提，本尺不需要 131 那条判据也红。

### 2.3 X8 `case "slo":` 换成命名常量

```
① c720494 + 变异，无仪器  rc=1  100/52/1/0  reds: TestAC4EveryLegIsNailedOrRuled
② 门开着                  rc=1  101/52/2/0  reds: TestAC1AC2DispatchHopGate133, TestAC4EveryLegIsNailedOrRuled
③ 门关着                  rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133      ← 四条红名
④ 还原                    rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
```
③栏原文：

```
AC#1/#2 RED: main.go:103: the case label sloCmdName131 is not a string literal, which resolves to the literal "slo". A leg is identified by a literal command name; a label that has to be resolved through a constant is a leg that can leave the census without saying so, which is ticket 131's second shape. This row carries it as "unparsed-label@main.go:103:sloCmdName131" instead of folding it into default.
AC#1 RED: leg "unparsed-label@main.go:103:sloCmdName131" (main.go:103) is dispatched by func main and covered by nothing: no nail in this gate's registry, no test case in this directory that drives a symbol belonging to this leg alone, and no WISP-LEG-COVERAGE-RULING: sentence naming it.
AC#1 RED: slo_windows.go:186 carries a coverage ruling for leg "slo", which is not in the dispatch census: a ruling about a command nobody dispatches is prose, not a decision.
AC#1/#2 RED: the usage block documents command "slo", which func main does not dispatch: the promise and the dispatch disagree, which is how a leg goes missing from this census without anything noticing.
```
⇒ 票面那句"不许把它修成整包坏掉才响"这里成立：`go build ./...` rc=0（X8-build-ins rc=0），
整包 100 条里只红这一枚，红的四条都是**这条腿自己出账**的形状。

### 2.4 X12 `var 别名 = installLogSink` + 经别名调用

```
① c720494 + 变异，无仪器  rc=1  100/52/1/0  reds: TestAC4EveryLegIsNailedOrRuled（131 补过的形 c 判据）
② 门开着                  rc=1  101/52/2/0  reds: TestAC1AC2DispatchHopGate133, TestAC4EveryLegIsNailedOrRuled
③ 门关着                  rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133      ← 红名点到 fake131
④ 还原                    rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
build：X12-1-build-base rc=0、X12-build-ins rc=0、X12-restore-build rc=0
```
③栏原文：

```
AC#1 RED: leg "fake131" (main.go:80) reaches installLogSink on this path: dispatch -> sinkAlias131 -> installLogSink
A listener installed on a dispatched leg has to have a nail in this gate's registry, or the block can be deleted in silence, ...
AC#1/#2 RED: leg "fake131" is dispatched by func main and documented in no line of the usage block: an operator cannot find it, and this gate has no second reading of the census to check it against.
```
⇒ 别名这一支本尺**自己建边**（`sinkAlias131 -> installLogSink` 印在红名里），
不是"因为整包坏了才响"；与 X4 同一枚判据 (b)，两处独立读数。

### 2.5 X14 install 藏在结构体字段的初值里（两拍都必须红）

**这一发是五发里唯一①栏真全绿的一发**——它也是票面上写死"第二拍也必须红"的那一发。

```
第一拍（装了听众、不给钉，不碰任何测试文件）
① c720494 + 变异，无仪器  rc=0  100/53/0/0  reds: 无   gate131: 1 hits PASS, 131red=0   ← 今天全绿，门也绿
② 门开着                  rc=1  101/53/1/0  reds: TestAC1AC2DispatchHopGate133          ← 只有本票仪器红
③ 门关着                  rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133          ← 门本来就没红，关不关同一份红名
④ 还原                    rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
build：X14-1-build-base rc=0、X14-build-ins rc=0、X14-restore-build rc=0
```
②栏原文（门关着那份与它同一批红名，见③）：

```
AC#1 RED: leg "sfx131" (main.go:80) reaches installLogSink on this path: dispatch -> holder131.open -> installLogSink
A listener installed on a dispatched leg has to have a nail in this gate's registry, or the block can be deleted in silence, which is the reading ticket 133 was filed for. Fix: write the case, then add {leg: "sfx131", test: TestYourCase, entry: "cmdSfx131"} to legCovers133.
AC#1/#2 RED: leg "sfx131" is dispatched by func main and documented in no line of the usage block: an operator cannot find it, and this gate has no second reading of the census to check it against.
```
⇒ 红名点到的正是票面要的那句话：**这条腿装了听众却没钉**，且路径写的是
`dispatch -> holder131.open -> installLogSink`（字段初值那条边被本尺接住了）。
同一发变异在**门开着**时 `TestAC4EveryLegIsNailedOrRuled` 仍 `PASS`、`AC#4 RED` 命中 **0**
⇒ 这一形今天只有本票这把尺看得见，不是蹭 131 的门。

**第二拍（把字段初值与那次调用删掉，该文件代码里 `installLogSink` 引用为 0）**

```
① c720494 + 变异，无仪器  rc=0  100/53/0/0  reds: 无   gate131 PASS, 131red=0    ← 复验方实测同一形"整包仍是 100/53/0/0"，逐字复现
② 门开着                  rc=1  101/53/1/0  reds: TestAC1AC2DispatchHopGate133   ← 门的账本仍写着 install=false，只有本尺红
③ 门关着                  rc=1  100/52/1/0  reds: TestAC1AC2DispatchHopGate133   ← 本票仪器仍红（票面写死的那一条）
④ 还原                    rc=0  101/54/0/0  RESTORE-DIFF rc=0 lines=0
build：X14B2-build-base rc=0、X14B2-build-ins rc=0、X14B2-restore-build rc=0
```
③栏原文，连同账本那一行：

```
AC#1 RED: leg "sfx131" (main.go:80) is dispatched by func main and covered by nothing: no nail in this gate's registry, no test case in this directory that drives a symbol belonging to this leg alone, and no WISP-LEG-COVERAGE-RULING: sentence naming it.
That is the second beat of ticket 133's fifth shot: a leg whose install block was deleted has no install-based obligation left, and "nobody anywhere verifies this command" is the fact that survives it. Fix: drive it from a case, or write the ruling next to the code that owns the leg.
AC#1/#2 RED: leg "sfx131" is dispatched by func main and documented in no line of the usage block: an operator cannot find it, and this gate has no second reading of the census to check it against.
  leg sfx131  main.go:80  installs=false handoff=false covered=RED nothing  entries=cmdSfx131
```
⇒ 红名换了判据但没换仪器：第一拍红在 (b)"装了听众没钉"，第二拍 install 已经拆掉，
红在 (d)"这条被分发的腿什么都没覆盖"。**只把第一拍弄红等于放过这一形**，所以两拍的③栏都在上面。

**两处如实登记的口径差**（都不影响判定，写出来防下一个人复核时对不上）：
1. 复验方那句"`grep -c installLogSink` = 0"是连注释一起数的；本尺种的 beat-2 文件在**注释**里
   引用了一次那个名字（`sfx131.go.x14b2.off:8`），代码里 0 处。判定材料是 AST 的调用边
   （账本 `installs=false`），与 grep 无关。
2. 快照里退役的种件改名成 `*.off` 保留着（"临时件只建不删"），所以 `diff -r` 反查带
   `--exclude='*.off'`；`RESTORE-DIFF lines=0` 指的是除此之外该树与纯净树逐字节相同。

### 2.6 AC#1 逐发红名一览（判据：五发都红，且关门仍红、红名是本票仪器）

| 发 | ①无仪器 | ②门开着 | ③门关着：红名 | 判据 |
| --- | --- | --- | --- | --- |
| N-3 | rc=1（只 131 的门红） | 本尺 + 131 的门 | `TestAC1AC2DispatchHopGate133` 唯一 | (a) 三枚孤儿符号 + 钉的入口不再可达 |
| X4 | rc=1（只 131 的门红） | 本尺 + 131 的门 | 同上 唯一 | (b) `--diag` 装了听众没钉（本尺把它读成一条腿） |
| X8 | rc=1（只 131 的门红） | 本尺 + 131 的门 | 同上 唯一 | 标签非字面量 + 合成腿未覆盖 + 裁决/usage 双向差 |
| X12 | rc=1（只 131 的门红） | 本尺 + 131 的门 | 同上 唯一 | (b) 别名边 `sinkAlias131 -> installLogSink` |
| X14 一拍 | **rc=0 全绿、门 PASS** | **只有本尺红** | 同上 唯一 | (b) 字段边 `holder131.open -> installLogSink` |
| X14 二拍 | **rc=0 全绿、门 PASS** | **只有本尺红** | 同上 唯一 | (d) 被分发的腿未覆盖（install 已拆） |

⇒ **AC#1 判定：PASS**。五发（含 X14 两拍）逐一红、红名点到本票仪器、把 131 的门 `-run` 排除后
仍红；④栏每发还原复绿且 `diff -r` 零差异。没有哪一发"只有 131 的门能红"，因此不必触发
"这一形由 131 守、回 131 续单"那一条。唯一与票面预期不同的是 §2.0.1：①栏今天不再全绿，
因为 131 在它自己的续单里把 N-3/X4/X8 三形接住了——照实登记，不改读数。

## 3. AC#2 覆盖面主张

**选的是清单式**，且是两份互相核的清单：一份从 `func main` 自己的分支现读
（`switch` 的字面量 case + `if len(args)==0` + 任何把 argv 槽位与字面量比较的早退 `if`），
一份从给人看的 `usage` 文字块现读，两者**双向差集为空**才绿。每枚被分发的腿在账本里占一行，
逐行写着它够不够得着听众、走的是哪条边、被什么覆盖。生产入口 → 分发调用者的对应关系由
判据 (a)(c) 钉住：凡够得着 `installLogSink` / `models.WireDownloading` 的符号，必须从
`func main` 走得到，且分发目标必须是本目录声明的函数。

**"删掉它哪条用例会红"的答案**：本尺只有这一枚用例 `TestAC1AC2DispatchHopGate133`，
而"没有它的树"不是假设——§2 每一发的第①栏就是在**不含这枚用例的 `c720494` 快照**上做的同一发变异，
五发（含 X14 第二拍）**逐栏全绿、四数不降**（100/53/0/0）。也就是说：把本尺摘掉，这五发的形状
今天仍然没有任何一枚用例会红；装上它，五发都红，且红名点到的是本尺（`leg_dispatch_gate_133_test.go`）。
仪器自己的钉清单 `legCovers133` 里那四枚声明（run/models/secret/no-args）也是同一回事：
删掉其中一枚声明，对应那条腿立刻从"nailed"掉进"covered by nothing"，红的仍是本尺本人。

**清单式看不见、图式（从 `main` 出发的可达性）能看见的形状**，逐条形如"本尺今天判不出"而非"本尺判错了"：

1. **跨包的那一跳**。本尺的两枚能力根都在 `cmd/wisp` 这一个目录里按名字读；
   若哪天持久听众是 `internal/**` 里某条腿替 `cmd/wisp` 装的（票 131 的复验方 §10 登记过同一残窗），
   图式若以模块依赖为边就看得见，本尺的清单只到本目录，看不见。票面已把跨包那一支划出本票地界，
   本格不声称覆盖。
2. **`func main` 之外的分发构造**：`init()` 里注册的命令表、`flag` 包的子命令、脚本侧注册。
   这类边不在 `func main` 的分支上，本尺的清单收不到它；图式以 `init` 为额外根就能走到 handler。
   本尺对今天的树零误伤（实测 `func main` 的每一处调用都落在被归类的分支里），
   但那份"零"只覆盖 `func main` 这一枚函数。
3. **接口动态派发与方法值**（`var f = v.M`、`iface.Do()`）。图式接 CHA/VTA 能把接口的全部实现连上，
   本尺按名字并集过近似；放不下的边本尺**红着承认自己是瞎的**（判据 (c) 与 `pkg.blind`），
   但"承认瞎"不等于"看见了那条边"。票面把方法值那一支划走，本格照登。

**反方向也写清楚，因为图式不是本尺的上级**：`//go:embed`、反射、`goja` 注册这类**静态调用图上根本没有边**的
依赖，图式会被骗过——`panel-assets` 那条腿卖的字节是 `internal/panel` 里 `//go:embed` 打进二进制的，
从 `main` 出发的调用图看不见"这条腿依赖那批文件"这件事；本尺的清单至少把这条腿的**命令名**与
`usage` 里那行**给人看的字**绑成互相核对的两处，少一处就红。两种形状谁都不是另一种的超集，
所以本格不写"顺带全覆盖"，只写上面这六条（三条形 + 三条反形）。

## 4. 动了哪几枚数（四数账）

**接续方**：`worker-ticket133-ac1-r2`（补 §4/§5/§6；不改 §2/§2.6/§3 任何一行读数）。
**锚定**：开工 `git rev-parse --short HEAD` = **`e113b1a`**（分支 `dev`、`git status --porcelain` 空）。
测量时刻 `date -u`：11:49:28z–11:51:55z（本地 19:49–19:51 +08）跑 §4.1；11:53:06z–11:54:06z 跑 §4.2；
11:54:32z–11:55:46z 跑 §4.3；全部在**工作树本体**、非 `/tmp` 快照（§4 是"今天的树动了哪几枚数"的账，
不是变异复现，不需要前任那三棵树）。**PATH 就位**：`export PATH="$PWD/third_party/sherpa-onnx:$PATH"`
（`onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll` 三枚齐全，缺一即 `0xc0000135`
＝环境故障不是发现，见 §1）。

### 4.1 字面 AC#5 命令 `go test -count=2 -v ./cmd/wisp/`（四数只从 `-v` 量）

```
$ go test -count=2 -v ./cmd/wisp/          rc=0   141.780s   ok github.com/CarlosShao/wisp/cmd/wisp
grep -cE '^=== RUN'   → 202
grep -cE '^--- PASS'  → 108     (顶栏；indented subtest 另加 47 枚 = 202 枚 PASS 行)
grep -cE '^--- FAIL'  → 0
grep -cE '^--- SKIP'  → 0
grep -cE '\[no tests to run\]' → 0
```
⇒ `count=2` 原始读数 **202 / 108 / 0 / 0**，per-count 归一 **101 / 54 / 0 / 0**，
逐字复现 §2.0 `S0-ins-pristine`（`RUN=101 TOPPASS=54 TOPFAIL=0 TOPSKIP=0`）；
本包 54 枚顶层 + 47 枚子用例，`count=2` 下 `=== RUN` 202、顶栏 `--- PASS` 108 都是 per-count ×2。

### 4.2 CI 逐字同形 `sh scripts/wisp-cli-tests.sh`（`-v -count=1`，与基线 100/53/0/0 直接同形）

```
$ sh scripts/wisp-cli-tests.sh                rc=0   57.707s
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|
  TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|
  TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$]
  top-level: PASS=54 FAIL=0 SKIP=0, === RUN=101, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=101  --- PASS=54  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
```
⇒ **CI 形状读数 101 / 54 / 0 / 0**（顶栏 PASS=54、RUN=101，SKIP=0，`runtests.sh` 判 OK）。
CI 的 `-skip` 名单是 `portable-tests.sh` 的账本，**一枚 cmd/wisp 用例都不在里面**（那些 skip 项都在
`internal/**`），所以新用例在 CI 里真跑、真计。

### 4.3 复现基线：`-count=1 -v -skip '^TestAC1AC2DispatchHopGate133$'`

把本票新增那一枚顶层用例用 `-skip` 摘掉，**工作树未动、任何文件未改**（`-skip` 使 Go 在跑前过滤掉
它，`--- SKIP` 行不产出，所以是"少跑一枚"而不是"跑成一枚 SKIP"）：

```
$ go test -count=1 -v -skip '^TestAC1AC2DispatchHopGate133$' ./cmd/wisp/
  rc=0   57.691s
grep -cE '^=== RUN'   → 100
grep -cE '^--- PASS'  → 53
grep -cE '^--- FAIL'  → 0
grep -cE '^--- SKIP'  → 0
```
⇒ **100 / 53 / 0 / 0**——与票 131 第三任在 `bcb03aa` 量的基线逐字同形。
这既是基线的**独立复现**（工作树没退步到 `bcb03aa` 那时代码的形状；差异就是新增那一枚用例），
也是 §4.4 名册 diff 的来源。

### 4.4 名册逐名对（comm，不只给差值）

从 `-v` 顶栏 `--- (PASS|FAIL|SKIP): <name>` 抽名（sed 掉状态与时延）→ `sort -u`：
- `roster-head.txt` = §4.1 的 54 枚（HEAD 上有仪器）
- `roster-skip133.txt` = §4.3 的 53 枚（摘掉本票那一枚）
两份都在 `/d/tmp/wisp133-ac1-r2-sec4/`；命令与逐名输出：

```
$ comm -23 roster-head.txt roster-skip133.txt     # 只在 head 的
TestAC1AC2DispatchHopGate133

$ comm -13 roster-head.txt roster-skip133.txt     # 只在 skip133 的（应为空 ⇒ 无回归）
（空）

$ comm -12 roster-head.txt roster-skip133.txt | wc -l
53                                                 # 共有 53 枚，一字未变
```

**54 枚名册全列（`TestAC1AC2DispatchHopGate133` 是本票新增，其余 53 枚与基线同）**：

```
* TestAC1AC2DispatchHopGate133                      <- 本票新增（ed18727）
  TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
  TestAC1ResidentLegInstallsItsLogListenerOnDisk
  TestAC1ResidentLegOutlivesItsOwnLogFailure
  TestAC2AC3DegradedLegsStillDeliverTheirVerdict
  TestAC2AuditTrailLandsInTheRunLegLogFile
  TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128
  TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
  TestAC2RealProcessRefusesOnEveryLegWithoutAppData128
  TestAC2RefusalMarkersAreNotAShortenableList128
  TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128
  TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
  TestAC2SealNoticeLandsInTheRunLegLogFile
  TestAC2TestDataDirBranchStillResolves128
  TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
  TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
  TestAC3EmptyDataRootIsARefusalNotAFallback
  TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole
  TestAC3LogSinkLandsInsideTheEnvDataRoot
  TestAC3SecretLegBooksItsAuditRecordsOnDisk
  TestAC4EveryLegIsNailedOrRuled
  TestComposedGateBlocksAWriteForTwoSeconds
  TestHostDispatchThroughTheAssembledBridge
  TestMissingBlobFailsUnconfiguredNeverSilently
  TestProcessCommandLineProbeDetectsAPlantedValue
  TestProcessCommandLineProbeHelperProcess
  TestProvidersDiscoverListsWhatTheServerServes
  TestProvidersProbeRecordsMeasuredThinkingFalse
  TestProvidersProbeRecordsMeasuredThinkingTrue
  TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
  TestRunTextTaskFailNextIsClassified
  TestRunTextTaskKeyResolvesInTheStore
  TestRunTextTaskTextPathEndToEnd
  TestSecretArgvCarriesNoSecret
  TestSecretEndToEndConfigRefResolvesAtRequestTime
  TestSecretFailurePathsLogAndPrintNoPlaintext
  TestSecretFlagsAreBoolOnly
  TestSecretFromStdinWritesNoIntermediateFile
  TestSecretOverwriteIsAnnounced
  TestSecretPortableModeUsesTicket06Seam
  TestSecretRealBinaryRefusesValueFlag
  TestSecretSameNameUnderThreeEnvsIsThreeBlobs
  TestSecretSetConfirmationMismatchStoresNothing
  TestSecretSetGetListRoundTrip
  TestSecretSetRejectsBadNamesAndEmptyInput
  TestSecretSetWithoutConsolePointsAtFromStdin
  TestSecretUnsetRefusesWhileReferenced
  TestSecretUsageAndUnknownSubcommand
  TestSecretValueCarryingFlagsAreRefusedAndUnechoed
  TestTicket101ManualSwitchSurvivesRestart
  TestTicket101ModeSwitchUsesTheRealL2Gate
  TestTicket101SessionGrantDoesNotCrossRestart
  TestTicket101UnreadableModeFailsLoudlyAndStrict
  TestTicket101UntouchedConfigRestartsAtDefault
```

### 4.5 逐名账对上游两票的意义（票 124 批次与票 131 翻格都靠这两个数）

- **基线 100/53/0/0 → 101/54/0/0**（per-count）：**只动 RUN、PASS 两枚数**（各 +1），
  FAIL/SKIP **零动**（0/0）。
- **逐名对上**：新增一枚 `TestAC1AC2DispatchHopGate133`（`cmd/wisp/leg_dispatch_gate_133_test.go:144`，
  `ed18727` 引入），其余 53 枚逐字未变（`comm -12` = 53、`comm -13` = 空）。
- **`ed18727` 里 `main.go`/`panel_assets.go`/`slo_windows.go` 的 +19/+8/+6 行都是 `WISP-LEG-COVERAGE-RULING:`
  注释**（给 version/help/default/panel-assets/slo 五枚腿的裁决句），**不新增用例、不改变四数**；
  工作树 `cmd/wisp` 从 `c720494`（前任的基线树）到 `e113b1a` 的全部差异文件与行数：
  ```
  $ git log --oneline c720494..HEAD -- cmd/wisp/
  ed18727 feat(133 AC#1/AC#2): 第二把尺——分发跳本身...
  $ git diff --stat c720494..HEAD -- cmd/wisp/  (只 ed18727 一枚)
   cmd/wisp/leg_dispatch_gate_133_test.go | 1378 +++++++++++++++++
   cmd/wisp/main.go                       |   19 +
   cmd/wisp/panel_assets.go               |    8 +
   cmd/wisp/slo_windows.go                |    6 +
  ```
- **票 124 的批次**：其 §7/§9 的兄弟代理 (`worker-ticket124-ac2b-2-r2`) 会引用 124 自己的 RUN/SKIP 差；
  **本票 §4 的四数不动那两枚数**——本票的读数域是 `cmd/wisp/`，124 是 `tools/llm`／`internal/*`。
- **票 131 的翻格**：`TestAC4EveryLegIsNailedOrRuled`（131 的门）仍在这 53 枚共有名里、仍 `--- PASS`；
  本票仪器**不共用**它的钉清单也不读它的账本（§0 明写、`leg_dispatch_gate_133_test.go:126-133` 注释重申），
  所以 131 那扇门的名字与判定都不因本票而变——131 的门 `-run` 关掉后本尺的读数见 §2 每一发③栏。

## 5. 门禁与全仓仪器

接续方补记；工作树本体、非快照、**HEAD 未动**（§4 的 `00d3f34` 之后继续，`git status --porcelain`
内 `cmd/wisp/**` 始终为空，测量与被验版本同形）。

### 5.1 `gofmt -l cmd/wisp/`（整包）

```
$ gofmt -l cmd/wisp/          rc=0   （无输出）
```
⇒ 整包**零文件待格式化**（含新落的 `leg_dispatch_gate_133_test.go`、`ed18727` 改过的
`main.go`/`panel_assets.go`/`slo_windows.go`）。

### 5.2 `go vet ./cmd/wisp/`（windows native）

```
$ go vet ./cmd/wisp/          rc=0   （无诊断）
```
⇒ 本包类型检查／vet 分析**零告警**，与前任 §2.5 三处 build 记 `rc=0`（X14-restore-build 那一栏）
同调。

### 5.3 `GOOS=linux go vet ./cmd/wisp/`（交叉 vet，工具链假象，**既不算破口也不算清白**）

```
$ GOOS=linux go vet ./cmd/wisp/          rc=1   2026-09-23 11:57:56z
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
$ GOOS=linux CGO_ENABLED=0 go vet ./cmd/wisp/          rc=1
（同上一条诊断，逐字）
```
⇒ **rc=1 死在包加载**（`build constraints exclude all Go files`），**够不到 `cmd/wisp` 的类型检查**：
交叉编译下 cgo 默认关，linux sherpa 绑定包内所有 `.go` 都带 `//go:build cgo && linux`，
被 `exclude all`，`go vet` 在解析 import 阶段就退出。诊断文本**不含任何 `cmd/wisp/*.go:line:col`**、
只有第三方包路径——这是工具链假象（brief 明写"既不算破口也不算清白"）；**票 131 的 Linux 编译破口
曾被"整树 rc=1＝工具链假象"放过、带病 4h17m**（`feedback-verification-blind-spots`），所以这里
显式逐错误行归因到 `sherpa-onnx-go-linux@v1.13.8`（**不是** `cmd/wisp/*.go`），并按下面 §5.4 单独
补一枚真类型读数。

### 5.4 Linux 容器**真类型读数**（`golang:1.27`，挂载已证非空）

```
$ date -u → 2026-09-23 12:00:43z
$ MSYS_NO_PATHCONV=1 docker run --rm \
    -v "/d/work/workspace/projects plans/Wisp:/src:ro" \
    -v "/d/work/base/gopath/pkg/mod:/gomod" \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e GOSUMDB=off -e GOMODCACHE=/gomod -e CGO_ENABLED=1 \
    golang:1.27 sh -c 'set -e; echo "--- mount proof ---"; ls -l /src/go.mod; cd /src; \
      go env GOOS GOARCH CGO_ENABLED GOMODCACHE; go vet ./cmd/wisp/; echo VET_RC=$?'
--- mount proof ---
-rwxrwxrwx 1 root root 855 Sep 21 01:17 /src/go.mod
--- go env ---
linux
amd64
1
/gomod
--- go vet ./cmd/wisp/ ---
VET_RC=0
```
⇒ **rc=0、零诊断**：这是 §5.3 那一发**唯一能替代的真类型读数**——linux 容器里 CGO_ENABLED=1、
`//go:build cgo && linux` 满足，包加载成功，`go vet` 真的跑完了分析器。挂载**先证非空**（`ls -l /src/go.mod`
出 stat 行、`ls /gomod` 出 `github.com/...` 目录、`ls /src/cmd/wisp/*.go` 出被测文件），
不走"Git Bash 下 `-v "C:\…"` 静默挂空且 rc=0"那条假绿（`MSYS_NO_PATHCONV=1` + `/d/...` 路径）。
`GOPROXY=off` 强制走宿主 module cache、**未联网**、**未新下模块**。仓库 `:ro` 挂载，容器**不会**
往工作树里写。**不覆盖 §2 任何读数**（那是前任变异读数，本节是本票门禁，两者是两回事）。

### 5.5 全仓仪器 `sh scripts/d22scan.sh`（作用域是整个 `internal/`，按包门禁看不见它）

调用形状（**wrapper 从仓根之外也不能省**——`scripts/d22scan.sh` 从自身位置派生仓库根、
先 `cd tools/d22scan` 再 `go run . -root "$root"`；直接 `go run ./tools/d22scan -root .` 会撞
"main module does not contain package"，`cd tools/d22scan && go run .` 会扫到空树——两种都曾是假绿）：

```
$ sh scripts/d22scan.sh                 rc=0   2026-09-23 12:01:58z → 12:02:03z (5s)
```

**Step 1（正向控制，runtests.sh `-count=1 -v`，证明门"能红"）**：
```
runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0
  关键子用例（都 PASS）:
    TestScanDetectsAllSeededViolations              种子违规被检出（否则后续"绿"=仪器瞎）
    TestScanCleanRepoIsGreen                        纯净仓库判绿
    TestAllowlistSuppressesOnlyListedPaths          allowlist 只压列出的路径
    TestCheckRootRejectsBlindRoots/{empty_dir, go.mod_only,
      the_d22scan_module_itself, repo_skeleton_with_two_go_files}
                                                    拒绝四种"看起来像根其实扫不到东西"的根
    TestBuiltBinaryGoesRedEndToEnd/{seeded_violation_exits_1, empty_live_scope_exits_2,
      armed_ban_6_goes_red_on_the_panel_violation_exits_1,
      frontend_tree_gone_while_declared_live_exits_2,
      ban_7_tree_gone_while_declared_live_exits_2,
      fully_live_fixture_exits_0}                   真二进制端到端 6 形
```
⇒ 正向控制 **21 枚顶层 PASS / 0 FAIL / 0 SKIP**——仪器**能红**，随后的"clean"不是"仪器瞎"。

**Step 2（真扫工作树）**：
```
d22scan: examined 225 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/         examined  43 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  43 text files
d22scan: scope ban #8 internal/         examined 397 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  38 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
```
⇒ **rc=0，八个 scope 全部非零命中、无违规**；`internal/` 203 / `cmd/` 22 / `internal/tools/` 18 枚
生产 `.go` 全在扫描面内（票 133 的新文件不在 `internal/` 也不在 `tools/d22scan/` 的禁改面上，
被 `ban #8 internal/` 与 `ban #8 cmd/` 的 emoji 扫各覆盖到）。

### 5.6 §5 总判

- `gofmt -l cmd/wisp/` **整包 rc=0**、无输出。
- `go vet ./cmd/wisp/` windows native **rc=0**、零诊断。
- `GOOS=linux go vet ./cmd/wisp/` **rc=1** 停在 cgo 包加载（诊断只有第三方包路径，**没有** `cmd/wisp/*.go`
  的任何 file:line），**既不算破口也不算清白**。
- Linux 容器真类型读数（golang:1.27，挂载证非空、CGO_ENABLED=1） **`go vet ./cmd/wisp/` rc=0、零诊断**。
- 全仓仪器 `sh scripts/d22scan.sh` **rc=0**、正向控制 21/0/0/0 与真扫 8 scope 全绿。

## 6. 没做到的

（待填）
