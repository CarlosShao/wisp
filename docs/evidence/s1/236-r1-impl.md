# 236-r1 交付件 — AC#5：红名册这条尺的两发（整包名册红 ＋ 单包隔离复量）

- 腿：`236-r1`（实现者／写腿）。票 236，只剩 `AC#5` 一格。
- 起手锚 `53b9be53`（2026-10-22 12:47 +08）。
- **仪器说明（AC#5 要的那份，具名）**：`docs/evidence/s1/236-r1-instrument-note.md`。
  规矩本体全在那里（sec 1..9）；本件只放读数、判语与门禁。
- 写面：只 `docs/evidence/s1/236-r1-impl.md` 与 `docs/evidence/s1/236-r1-instrument-note.md`。
  零 `*.go`、零 `scripts/**`、零 `.github/**`、零 `cmd/wisp` 写、零 push、零翻框、零改票名／`Status:`。
- 原始日志：`D:\tmp\wisp\20261022\236r1\logs\`（只建不删，清单见 sec 6 尾）。

---

## sec 1 我自己跑的两形读数（逐字；NOT 抄票面的 5.15／5.055）

环境：`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`、`PATH` 前置 `third_party/sherpa-onnx` 与 `build`、
`CGO_ENABLED` 未设、`-timeout 30m`、一次一发（起跑 `go` 进程 0）、本机 CPU 73%。

### 发二：单包隔离复量，5 发（`-run '^TestTicket223RefusedLooseningKeepsOldValues$' ./cmd/wisp`）

```
shot1 rc=0 :: --- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.24s) | ok ... cmd/wisp 7.474s
shot2 rc=0 :: --- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.62s) | ok ... cmd/wisp 7.738s
shot3 rc=0 :: --- PASS: TestTicket223RefusedLooseningKeepsOldValues (8.31s) | ok ... cmd/wisp 8.568s
shot4 rc=0 :: --- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.44s) | ok ... cmd/wisp 7.638s
shot5 rc=0 :: --- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.70s) | ok ... cmd/wisp 7.909s
```

5/5 绿。票面给的是 `3/3` 绿（5.055／5.019／5.034 s）：**判决方向一致**，秒数本机高约 1.5 倍（负载 73%）。

### 发一：整包名册红（`go test -count=1 -v -timeout 30m ./cmd/wisp`），三发

| 发 | rc | 耗时 | 顶层 RUN | 顶层判决 | 子测试 RUN | 子测试判决 | 名级红 |
|---|---|---|---|---|---|---|---|
| A | 1 | 353.54 s | 162 | 161 | 58 | 58 | 2 |
| B | 1 | 361.83 s | 167 | 166 | 58 | 58 | 1 |
| C | 1 | 379.01 s | 162 | 161 | 58 | 58 | **0** |

三发的红名册逐字：

- A：`--- FAIL: TestPanelHostRealWindowHopAndLifecycle (3.21s)` ＋ `--- FAIL: TestPanelHostLatencyPercentilesAC2 (0.00s)`
- B：`--- FAIL: TestAC14GoSideEvalPushReachesThePage (0.59s)`
- C：（空）

三发里那枚 223 用例的判决行逐字：

- A `--- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.33s)`
- B `--- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.18s)`
- C `--- PASS: TestTicket223RefusedLooseningKeepsOldValues (7.89s)`

### 与票面的差：这一条必须单列（派单条款 2）

票面 `AC#5` 写「整包 5.15s **红**、隔离 3/3 绿」。**我整包三发全绿**。
=> 不是"秒数漂了"，是**判决翻转**（红变绿）。旁证票面确已过期：票面引的失败行
`config_reload_223_test.go:379`，那枚 `t.Errorf` 今天在 **`:512`**（文件已挪约 133 行）。
=> 本件口径：**"票面那一发未复现"**，NOT "它不是计时红"。两者是不同的断言。

---

## sec 2 写成的规矩（落点具名）

**落点＝`docs/evidence/s1/236-r1-instrument-note.md`**（与主交付同目录；派单给的最保守合法落点）。
未落 `scripts/**`、未落 `docs/reports/**`（那两块本票不授权）。

规矩的骨架（全文见仪器说明）：

1. **发一 整包名册红**：`go test -count=1 -v -timeout 30m ./<包>`，只取"哪些名字红了"。
2. **发二 单包隔离复量**：`-run '^(名1|名2)$'` 锚定全名，只对发一红了的名字跑，至少 3 发。
3. **记账判据是发二的 rc 与名，NOT 发一的秒数**：
   - 发二 `rc=0` 且 PASS => **计时红，不记账**（不进"新增红"，也不许进"我修好的红"）。
   - 发二 `rc=1` 且同名 FAIL => **真红，记账**并归因到具体一行。
   - 发二 `rc=1` 但红的是别的名 => **换位，记账并单独报**（说明两发之间字节或负载变了）。
4. **发一必须先过两道恒等式**（否则那份日志根本不是名册）：
   - 分层对账 `RUN_top - VERD_top = 0`；
   - 声明对账 `go test -list` 的枚数 减 `RUN_top` 等于"整发没跑到"的枚数，必须写出来。
5. **整包这把尺会整族漏**（`A1124(5)`）：`panel_*`／`plan_*` 之后的族整包恒 `RUN=0`，
   整包的绿对它们**恒无信息**，只能用显式 `-run` 另取。报告必须写"这一发跑到哪儿为止"。

---

## sec 3 "两枚计时红换位"的重钉（我自己那一发，NOT 复述票面）

票面那句是"在册常红名册今晚与实现腿那发不一致（两枚计时红换位）"。我用**同一台机、同一份字节、连续两发**重钉，
并把它推硬一层：**换位不需要跨机、跨 commit，同机连发就换。**

### 三发之间到底变没变字节（先证明可比性）

- `cmd/wisp` porcelain 起手 **0 行**；`git rev-parse --short HEAD` 全程 `53b9be53`。
- 发 A 结束时（12:56:09）工作树干净。
- **发 B 与发 C 期间，`cmd/wisp` 被别的腿落了四枚未跟踪探针件**（`zz227r1_probe{,2,3,4}_test.go`，
  mtime 12:57:28／12:58:20／12:59:27／13:11:15），而**未跟踪的 `_test.go` 会进同一测试二进制**。
- 我在发 C 前后对 `cmd/wisp` 全部 145 枚 `*.go` 逐个取 sha256：**发 C 期间文件集变了**
  （多出 `zz227r1_probe4_test.go`）=> 发 B 与发 C 之间不可比；发 A 与发 C 之间也不可比（探针 1-3 已入树）。
- => **三发没有一对是"同字节"的**，所以我**不**用它们互相当"同字节换位"的证据。
  这三发证明的是另一件更硬的事，见下。

### 三发证明的事：红名册在**同一 HEAD** 上逐发不同，且**与有没有红无关**

- A 红 2 枚、B 红 1 枚、C 红 **0** 枚，而**三发的 HEAD 都是 `53b9be53`**。
- **发 C 决定性**：`rc=1` 而**名级红＝0**。=> `rc` 记的是"这个进程没走完"，不是"哪一枚红了"。
  任何拿整包 `rc` 当"有红"的读法，在发 C 上直接失效。
- 三发的截断点也不一样：A／C 停在 `TestAC13BringUpSurvivesAReusedThreadQuit`（起了没判），
  B 停在 `TestPanelThreadIsSTAAndExitsCleanly`（起了没判）。

### 两枚计时红的现钉（这才是"隔离复量"要咬的东西）

对发 A 那两枚红做隔离复量（`-run '^(TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2|TestAC14GoSideEvalPushReachesThePage)$' ./cmd/wisp`）：

```
iso rc=1
--- PASS: TestPanelHostRealWindowHopAndLifecycle (1.79s)
--- PASS: TestPanelHostLatencyPercentilesAC2 (0.00s)
--- FAIL: TestAC14GoSideEvalPushReachesThePage (1.62s)
FAIL ... cmd/wisp 3.727s
```

- 发 A 那两枚 => **隔离下全绿** => **计时红，不记账**。
- 发 B 那枚 `TestAC14GoSideEvalPushReachesThePage` => **隔离下仍 FAIL** => **真红，记账**
  （与台账 `pending-and-issues.md:21050` 记的既有真红同名）。
- => 同一把尺上"记账／不记账"两形**今天都咬得动**。这是 `AC#5` 要的"自带隔离复量那一半"的**活证据**，
  且**不依赖**票面那枚 223 用例是红是绿。

---

## sec 4 两个口径：名级 23 NOT 包级 3，以及 four numbers 互咬的算法

### 4.1 两个口径不是同一件事

- **包级**＝`grep -cP '^FAIL\t'`：有几个**包**判决 FAIL。一个包红 100 枚也算 **1**。
- **名级**＝`grep -cE '^[[:space:]]*--- FAIL'`：有几枚**用例**红了（含缩进子测试）。
- 两者**没有换算关系**。"3 个包红"可以是 23 枚红，也可以是 1 枚红。
- => 票面"名级 23 NOT 包级 3"成立，且**必须连口径三样一起写**：
  (1) 跑范围（本机包集 vs CI 全 job）；(2) 层级（顶层 vs 含缩进子测试）；(3) tier 集合（哪几个、是否互斥）。

### 4.2 four numbers 互咬：具体算法

`scripts/portable-tests.sh:851-856` 打的四个数，**内部不在同一层**：

```
ran=$(count '^=== RUN')   passed=$(count '^--- PASS')
failed=$(count '^--- FAIL') skipped=$(count '^--- SKIP')
```

- `^=== RUN` **数得到子测试起点**（Go 把 `=== RUN   TestX/sub` 打在**行首**）。
- `^--- PASS/FAIL/SKIP` **数不到子测试判决**（子测试判决是**缩进**的 `    --- PASS: TestX/sub`）。

本机实测（`cmd/wisp`，`53b9be53`）：`^=== RUN`＝**220**、`^--- PASS`＝**161**，差 **59**。按层拆开：

| 层 | RUN | 判决 | 差 |
|---|---|---|---|
| 顶层 | 162 | 161 | **1** |
| 子测试 | 58 | 58 | 0 |

=> 那 59 里 **58 是层差、1 才是真的没判决**。**拿 59 当"59 枚没判决"是读错层**（本件第一版就这么错过，见仪器说明 sec 8(e)）。

互咬算法（逐步、可照抄）：

1. **分层**：`RUN_top`＝`grep -cE '^=== RUN[[:space:]]+[^/]+$'`；`VERD_top`＝`grep -cE '^--- (PASS|FAIL|SKIP):'`。
2. **恒等式一**：`RUN_top - VERD_top` 必须＝0。非 0 => **这一发不是名册**，不许记红，重取。
3. **恒等式二**：`DECLARED`（`go test -list '.*' ./<包>` 的枚数）减 `RUN_top` ＝ 整发没跑到的枚数，
   **必须写进报告**（本机 `cmd/wisp`＝350 减 162／167 ＝ **188／183**）。
4. **恒等式三**：跨 tier 相加**先证互斥**。`portable-tests.sh` 的四个 pin 现跑去重：
   `core=27 / windows=11 / cli=1 / winsec=1`，**相加 40、去重并集 28**；windows 的 11 个包**全在 core 里**，
   winsec 那 1 个也在 core 里。=> **不互斥，相加会重复计数。**
   要用相加，必须给"包级去重后的并集"并声明 tier 集合；给不出就**不许相加**。
5. **名级总数**只能在 2、3、4 全过之后才报，且报的是**顶层**名级数；含子测试的另报一栏。

=> 票面那句"4／7／12 三发各自相加＝23 对得上"：**只有"三发的 tier 集合互斥"成立时才是名级总数**。今天四个 pin **不互斥**，所以这条相加**不能单独作为名级口径**，必须配第 4 步的并集声明。

---

## sec 5 反向判据：这份规矩在哪种情形下失效（文档形，附读数）

仪器说明是**文档**（`docs/evidence/s1/236-r1-instrument-note.md`），NOT 可执行件。
=> 按派单，反向判据＝**指出它会在哪种情形下失效并给那情形的读数**。四个面，逐条给读数：

**(1) 发一截断时，"完整名册"是错觉。**
- 读数：三发顶层 `RUN=162／167／162`、`VERD=161／166／161` => **每发恰 1 枚起了没判**；
  三发红名册**逐发不同**（2 枚 / 1 枚 / **0 枚**）。
- 挡它的尺：恒等式一。`RUN_top - VERD_top != 0` => 这一发不许当名册。

**(2) 发二只跑一发就判"计时红"。**
- 读数：那枚 223 用例隔离 **5 发全绿**，而它整包三发**也全绿** => 它**今天压根没红过**。
  => "发二绿"只证明"隔离下今天绿"，**不证明**"发一那枚是计时红"。发一绿＋发二绿 => 它没红过，
  **不许**被写进"计时红"那一栏凑数。

**(3) 拿 tier 相加当名级总数。**
- 读数：`core=27 / windows=11 / cli=1 / winsec=1` => 相加 **40**、去重并集 **28**；windows 11 个包全在 core 里。
- 挡它的尺：恒等式三。

**(4) 本规矩挡不住"两发之间字节变了"。**
- 读数：`cmd/wisp` 在 12:57:28／12:58:20／12:59:27／13:11:15 被别的腿落了四枚 `zz227r1_probe*_test.go`，
  **横跨我的发一与发二**。未跟踪 `_test.go` 会进同一测试二进制 => 两发不可比。
- 挡它的办法（规矩里写了，但**没有**自动手段）：两发都记 HEAD ＋ `*.go` 逐个 sha256 的文件集指纹，
  两次不同 => 重取。**这一条本机真的咬到过**（发 C 前后指纹不同）。

**附加：规矩自己第一版就是反例。**
- 第一版把恒等式写成 `grep -c '^=== RUN'` 减 `grep -c '^--- PASS'`，把 `cmd/wisp` 读成"**59 枚没判决**"。
  按层拆开后真值＝**1 枚**，那 58 是层差。=> 任何人报出一个非零"未判决数"，都要先问它是哪一层的差。

---

## sec 6 门禁 rc 与交付清单

| 项 | 命令 | rc |
|---|---|---|
| 写面之外零改动 | `git status --porcelain -- internal cmd scripts .github`（起手） | 0 行 |
| `docs/evidence/s1` 起手 | `git status --porcelain -- docs/evidence/s1` | 0 行（无别人未提交改动） |
| 本腿新增文件 | 只 `236-r1-impl.md`、`236-r1-instrument-note.md` | 2 枚 |
| 禁码位扫描（两份都扫） | 本件与仪器说明 | **0 命中** |

**具名没做**：
- 未跑 `cmd/wisp` 的 `-race`（本票射程外，且与 `227-r1` 抢同一枚重包）。
- 未复跑 CI（`gh`）读数：`236-a1/census.md:141`／`:142` 的 CI 侧数字是**转述**，我本腿只做本机读数。
- 未证票面那枚 223 用例"曾经是计时红"；只证"**今天三发整包全绿**"。
- 未动 `scripts/portable-tests.sh`（票 302 射程）、未动任何 `*.go`、未动 `cmd/wisp`（`227-r1` 在飞）。