# 236-r1 交付件 — AC#5：红名册这条尺的两发（整包名册红 ＋ 单包隔离复量）

- 腿：`236-r1`（实现者／写腿）。票 236，只剩 `AC#5` 一格。
- 起手锚 `53b9be53`（2026-10-22 12:47 +08）。
- **仪器说明（AC#5 要的那份，具名）**：`docs/evidence/s1/236-r1-instrument-note.md`。
  规矩本体全在那里（sec 1..9）；本件只放读数、判语与门禁。
- 写面：只 `236-r1-impl.md`、`236-r1-impl-part2.md`、`236-r1-instrument-note.md` 三枚。
  零 `*.go`、零 `scripts/**`、零 `.github/**`、零 `cmd/wisp` 写、零 push、零翻框、零改票名／`Status:`。
- 原始日志：`D:\tmp\wisp\20261022\236r1\logs\`（只建不删，清单见 part2 sec 6 尾）。
- **本件过帽（>160 行）拆两半**：本件＝sec 1..3（121 行）；
  **`236-r1-impl-part2.md` ＝ sec 4..6**（两个口径与互咬算法／反向判据／门禁 rc，93 行）。

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

