# 票 131 AC#4 复验（第三任）—— 裁的就是上一任判退回那一格

**验收方**：`acceptor-ticket131-r3`（本文件唯一作者；只读验收 + 自造变异，本文件是唯一写件）
**日期**：2026-09-23（开件 15:5x，逐节落盘，每节一枚 commit）
**被验交件**：票 131 的「续单」三件（续#1/续#2/续#3），实现方 `worker-ticket131-followup` + `worker-ticket131-followup-r2`
**被验 sha**：**`bcb03aa`**（= `717d822` 第 0 件 + `6f702ae` 三形判据 + `c49327c` 追认取证 + `bcb03aa` 票面翻格）
**判据原文**：`docs/evidence/s1/131-adversarial-acceptance.md` §4（X1..X12）、§6（`R-131-1..4`）、§8.3（结案三件事），
以及票面「续单 → 结案判据」那三句。

---

## 0. 锚定与被验面的完整性（先证明我读的是被验版本）

- 开工 `git rev-parse HEAD` = **`ab1462a8380a5a40d78a3d8b2e6c5b6598006a16`**，工作树 `git status --short` **为空** ⇒ 三任代理都入库了。
- HEAD 已前进一格（`bcb03aa..ab1462a` = `ab1462a` 本身）。按简报要求逐枚核被验面：

```
$ git diff --numstat bcb03aa..HEAD -- cmd/wisp
(零行)
$ git diff --numstat bcb03aa..HEAD
33  0   .scratch/wisp/issues/124-posix-harness-hands-the-sealing-floor-an-unresolved-tempdir-79-to-83-reds.md
24  0   docs/reports/pending-and-issues.md
```

⇒ **`cmd/wisp/**` 在 `bcb03aa..HEAD` 之间零字节改动**，唯一增量是两枚文档 ⇒ 被验面未被移动，本文全部读数同时适用于 `ab1462a` 的 `cmd/wisp`。
- **快照**：`git archive bcb03aa | tar -x -C /tmp/wisp131-acc-r3`，随后
  `cp third_party/sherpa-onnx/*.dll /tmp/wisp131-acc-r3/third_party/sherpa-onnx/`（三枚：`onnxruntime.dll`、
  `sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`）。**没拷就是 `0xc0000135`、四数全零**，简报预告过，我按它做了。
- **两记反查**：① `diff -r --exclude=third_party /tmp/wisp131-acc-r3 /tmp/wisp131-verify-acc-r3`（另一枚现取的
  `git archive bcb03aa`）⇒ **rc=0，零差异**；② `diff -r "D:/work/workspace/projects plans/Wisp/cmd/wisp" /tmp/wisp131-acc-r3/cmd/wisp`
  ⇒ **rc=0**，快照的 `cmd/wisp` 与仓库工作树逐字节相同。
- **未在仓内建 worktree / 未 checkout**；全程 `PATH=<快照>/third_party/sherpa-onnx`。
- ⚠ **一轮在飞观察**（不是裁决）：16:0x 我 `git status --short` 第二次量到
  ` M docs/reports/injection-timeline.md`（+21/-0，`git status --short -- cmd/wisp` 仍为空）。
  ⇒ 有兄弟代理正在写那枚文档，**被验面未被碰**，我一格未动它。
- **测前查在飞 run**（本机即 self-hosted runner）：`gh run list --repo CarlosShao/wisp --limit 5` ⇒ 五枚全 `completed`
  （最新 `35832874239` 07:39Z），无 in_progress ⇒ 本轮是安静窗口。整轮读数里 `0xc000013a` 早死 **零命中**（各发日志的
  `exit status` 计数见 §5 附）。墙钟只作旁证。

---

## 1. 基线复算（`bcb03aa` 纯净快照，未变异）—— 新基线是不是真的

简报警告"上一任的 82/82 已被三次移动过，谁拿它当分母谁就错"，也警告"重新基线必须是逐数解释过的"。两边我都量了。

| 仪器 | 续单方自报（`b723978`） | **我复算（`bcb03aa`）** | 判定 |
| --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | `=== RUN` 200 / 顶层 106 / 0 / 0，rc=0 | **RUN 200 / 顶层 PASS 106 / FAIL 0 / SKIP 0**，rc=0，`ok` 112.509s | 同数 |
| `sh scripts/wisp-cli-tests.sh`（CI 形状） | 100 / 53 / 0 / 0，rc=0 | **RUN 100 / 53 / 0 / 0**，rc=0，`ok` 54.953s | 同数 |
| `sh scripts/d22scan.sh` 八 scope | 203/22/40/18/16/40/390/37，rc=0 | **203 / 22 / 40 / 18 / 16 / 40 / 390 / 37**，rc=0 | 同数 |
| `gofmt -l cmd/wisp/` | 空 | **空**，rc=0 | 同 |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/`（v0.7.0，实存） | 空 | **空**，rc=0 | 同 |
| `go vet ./cmd/wisp/`（windows） | rc=0 | **rc=0** | 同 |
| `GOOS=linux go vet ./cmd/wisp/`（交叉） | rc=1，唯一错误行是 sherpa | **rc=1**，逐错误行归因见下 | 同 |
| Linux 原生 `go vet ./cmd/wisp/`（Docker） | rc=0 | **rc=0**，见 §3 | 同 |

基线整份 `-count=2 -v` 里 **`AC#4 RED` 命中数 = 0**（`grep -c 'AC#4 RED' /tmp/wisp131-acc-r3-base-count2.log`）。
⇒ 第二把尺在未变异树上零误伤，这一条我独立量到了（续单方在 `b723978` 报同数）。

### 1.1 "+18 RUN / +7 顶层 PASS 全归 129/130" 这句我逐枚对过

上一任锚 `56d8026` 的 CI 形状是 `82 / 46`，新基线是 `100 / 53` ⇒ 差 **+18 RUN / +7 顶层 PASS**。
逐枚查谁在这段区间往 `cmd/wisp` 加了用例：
我把两枚锚点各跑了一遍同一仪器，然后逐枚对名字集合：

```
git archive 56d8026 | tar -x -C /tmp/wisp131-r3-old        (+ 三枚 dll)
go test -count=1 -v ./cmd/wisp/   -> RUN 82 / TOP PASS 46 / FAIL 0 / SKIP 0, rc=0   # 与上一任 §2 的 CI 形状逐字相同
go test -count=1 -v ./cmd/wisp/   -> RUN 100 / TOP PASS 53 / FAIL 0 / SKIP 0, rc=0  # bcb03aa（§1 基线）
comm -23 <旧名字集> <新名字集>    -> 空    # 一枚也没丢
comm -13 <旧名字集> <新名字集>    -> 恰 18 枚
```

那 18 枚逐字是：**票 128** 的 `4e5d240`（`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` +5 子、
`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` +4 子、`TestAC2RefusalMarkersAreNotAShortenableList128`、
`TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128` +2 子、`TestAC2TestDataDirBranchStillResolves128`
⇒ **16 枚 = 5 顶层 + 11 子用例**）与**票 130** 的 `d87905c`
（`TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord`、`TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot` ⇒ **2 枚顶层**）。
逐枚静态反查同数：`func Test` 枚数 `56d8026`=49 → `4e5d240`=54（+5，128）→ `d87905c`=56（+2，130）→ `bcb03aa`=56。
CI 的 `-skip` 名单与这 18 枚零交集（逐名 grep）。

⇒ **数字结论成立、归因写错了名字。** `+18 RUN / +7 顶层 PASS` 这一笔我逐枚对上，**零枚被洗掉、零枚丢失**；
但它是**票 128（16）＋票 130（2）**，不是续单证据与票面写的"全归 129/130"。
`git log 56d8026..b723978 -- cmd/wisp` 里**根本没有 129**（129 的改动面在 `internal/**`，一枚 `cmd/wisp` 用例也没加），
而 128 被漏名。⇒ 登记为 `R-131r3-4`（措辞/归因级，不改判"四数不降"）。
另两枚在该区间动了 `cmd/wisp` 的 commit（`11f3927`＝114 只改 `run.go`、`9b5d64d`＝130 只改既有钉文件、
`717d822`＝本票第 0 件改名）对分母零影响，我也逐枚 `--name-only` 看过。

### 1.2 被驱腿自己的行号也下移了（续单那句话不完整）

续单/票面写"六发红名逐字不变，**只有门与钉自己的 `file:line` 前缀下移**"。前半成立，后半不完整：

```
bcb03aa vs 56d8026 的 install 点：logsink.go:139->144、models.go:276->284、secret.go:219->226、run.go:164->178、resident_windows.go:57 未动
X2 红名里的 emit-site：slog.Info@secret.go:321 -> :328, slog.Info@:461 -> :468, slog.Warn@:461 -> :468
```

⇒ 红名的**文本结构与所点的那条腿一字未变**（我逐字比过，见 §2），但**位移的不止仪器自己的前缀**——
被驱腿文件的行号也被票 128（`secret.go` +11/-4）与票 130（`models.go` +9/-1、`logsink.go` +22/-2）推着走了 7–14 行。
这句话要更正，否则下一个人拿"只有前缀动"去核 X2 会以为读数错了一位。并入 `R-131r3-4`。

---

## 2. 结案判据 1：六发红名逐字不变（X1/X2/X3/X6/X9/X10）—— **PASS**

仪器与上一任同形：`git archive bcb03aa` 纯净快照 + 三枚 dll + `go build ./cmd/wisp/` rc=0 之后才读 `go test -count=1 -v ./cmd/wisp/`。
变异我**按上一任的字节形状重造**（X3/X6/X8/X12 的 `main.go` 与被种的 `fake131.go`/`diag131.go` 直接取上一任快照里那枚文件，
`diff` 证同形；X1/X2/X3 的删行区间逐字同 `284,292d283` / `226,230d225` / `74,81d73`）。

| 发 | 造法（同上一任） | 上一任四数 | **我这次四数**（RUN/顶层PASS/FAIL/SKIP） | 原来那枚红 | 名字换没换 |
| --- | --- | --- | --- | --- | --- |
| **X1** 删 `models.go` 那 9 行 | 删 284-292，`grep -c installLogSink` = 0，334→325 行 | 82/78/**4**/0（3 顶层+1 子） | **100/50/3/0 + 1 子用例红**，rc=1 | `AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "models", which no longer installs the persistent sink (main.go:74): the nail is now proving something else, or nothing` | **未换**，仅 `214→318`、钉 `307→358`、`455→506` |
| **X2** 删 `secret.go` 那 5 行 | 删 226-230 | 82/78/**4**/0 | **100/48/5/0 + 1 子用例红**，rc=1 | 门给**两条独立读数**（`:298` 要裁决、`:318` 要钉）文本逐字同 | **未换**；多出的 2 枚红点名归票 130 的两枚 early-record 用例（`d87905c`），方向是查得更多 |
| **X3** 删掉 `case "models"` 整支、`models.go` 一字不动 | 删 74-81，生产调用者只剩声明 | 82/81/**1**/0 | **100/52/1/0**，rc=1，唯一一枚红在门上 | `... which main.go does not dispatch to any more` 逐字同 | **未换** |
| **X6** 新腿直呼 install、不给钉 | 取上一任那枚 `fake131.go` + `81a82,83` | 82/81/**1**/0 | **100/52/1/0**，rc=1 | `leg "fake131" (main.go:82) reaches installLogSink on 1 line(s) of this package and no registered nail names it.` 逐字同，行号也是 `main.go:82` | **未换** |
| **X9** 两处 install 挪去 `<根>-elsewhere` | 各改 1 行 | 门 PASS + 两枚钉红 | **100/48/5/0 + 2 子用例红**；门 `--- PASS`、`AC#4 RED` 命中 **0** | 两枚钉各红一次，红名正文带 `dir=...\001-elsewhere\logs` 原文 | **未换**；+2 枚同样是 130 那两枚 |
| **X10** 注册一枚不读盘的用例当钉 | `registerLegNail131("models", TestSecretFlagsAreBoolOnly)` | 82/81/**1**/0 | **100/52/1/0**，rc=1 | `nail "TestSecretFlagsAreBoolOnly" for leg "models" calls none of the shared sink readers (readLegSink131, readResidentSink, readSink) ... (leg site main.go:74)` 逐字同 | **未换**；门内**另多一枚**红＝本轮为 `R-131-2` 加的"没带入口符号"那条（同一枚用例内，顶层枚数不变） |

⇒ **判 PASS。** 六发原来那枚红**逐枚还在、名字一个没换**；新增的红全部点名归到票 130 的用例与本轮 `R-131-2` 的新判据上，
方向是**查得更多**。⚠ 两处措辞要更正（不影响判定）：X2/X9 各多的是**两枚**（`100/48/5` 对上一任 `82/78/4`），
续单证据与 A122① 一处写"两枚"一处写"一枚"，实测为**两枚**；X10 多的那一枚在**同一枚用例内部**，顶层 FAIL 枚数没变。

---

## 3. 结案判据 2：三发从绿变红（X4／X8／X12）—— **PASS**

同一仪器、同一台机器、`go build` rc=0 之后才读数。三发修复前都是**门 `PASS`／整包 82 全绿**（上一任 §4 原文）。

### 3.1 X4 ＝ 早退 `if` 分发一条装了听众的腿：红，且红名逐字点到 `--diag`

造法（照上一任）：`main()` 里 `args := os.Args[1:]` 之后插
`if len(args) > 0 && args[0] == "--diag" { os.Exit(cmdDiag131(args[1:])) }` ＋ 新生产文件 `cmd/wisp/diag131.go`
（内含 `installLogSink(root)` + `defer sink.close()`，`grep -c installLogSink` = 1），**不给它钉、不碰任何测试文件**。
`go build` rc=0、`go vet` rc=0。

`go test -count=1 -v ./cmd/wisp/` ⇒ rc=**1**，`RUN 100 / 顶层 52 / 1 / 0`，**唯一一枚顶层红就是门本人**，红名原文两枚：

```
leg_sink_gate_131_test.go:348: AC#4 RED: main.go:51 reads args (the local os.Args was bound to) in func main outside every branch
    this gate classifies (the argv binding, the `if len(args) == 0` leg, and the `switch args[...]`). Literals in that statement: "--diag"
    Statement: if len(args) > 0 && args[0] == "--diag" { os.Exit(cmdDiag131(args[1:])) }
    A dispatch that lives beside the switch is invisible to the leg list, so the leg it selects can install a listener with no nail and no row. ...
leg_sink_gate_131_test.go:348: AC#4 RED: main.go:52 reads args ... The statement names no literal command, so this reader cannot say which leg it dispatches.
    Statement: os.Exit(cmdDiag131(args[1:]))
```

⇒ 判据 2 那句"红名须点到 `--diag`"**逐字成立**（第一枚的 `Literals in that statement: "--diag"`；整份日志里 `--diag` 命中 2 次，全在红名内）。
账本仍是 **15 行、`--diag` 零命中** ⇒ 拦下这形的不是"多了一条腿"，是"**越权读 argv 就红**"——这条比"数腿"结实，
因为它不依赖门认得这条腿的形状。
零误伤：未变异基线整份 `-count=2` 里 `AC#4 RED` 命中 **0**（§1）。还原 ⇒ `diff -r` 只差一枚 `go build` 落的 `wisp.exe`、门 `--- PASS`。

### 3.2 X8 ＝ `case "slo":` 写成命名常量：红，且红来自门本人而不是"整包坏掉"

造法：`case "slo":` ⇒ `case sloCmdName131:` ＋ `const sloCmdName131 = "slo"`（**直接取上一任快照里那枚 `main.go`，逐字节同形**）。
`go build`/`go vet` rc=0。

- 整包：rc=1，`RUN 100 / 52 / 1 / 0` ⇒ **一枚顶层红，且它就是 `TestAC4EveryLegIsNailedOrRuled`**。
  这一条正是简报点名的"不许靠整包坏掉才响"：红的不是别的用例连带炸，是门自己判红。
- 红名点到常量与行号，并给出解析到的值：
  `the case label at main.go:84 is sloCmdName131, which is not a string literal. It is package-level const sloCmdName131 = "slo", so the leg this clause dispatches is "slo" ...`
- 账本 **15 行一字不少**，那条腿以 `unparsed-label@main.go:84:sloCmdName131` 自己的行在账
  （`install=false records=true ruled=true nails=- -> ruled`），`leg default main.go:95 -> no records` **不受影响**
  ⇒ 上一任点名的两个后果（静默出账 / 与别人撞键后静默）**各自都有了对面的红**。
  ⚠ 关键差别要说清：修复前是**两行都叫 `default`**、`leg slo` 这个键零命中；现在那个键改成了 `unparsed-label@...`，
  但那条腿的三列读数与裁决仍在账、且门为它红一次 ⇒ 裁决标记的承重没有被改名洗掉。
- 还原 ⇒ 只差 `wisp.exe`、门 `--- PASS`。

### 3.3 X12 ＝ 新腿 + 一枚 `var 别名 = installLogSink`：红，且与 X6 同一枚红名

造法：取上一任 `wisp131-acc131r2-x12-aliasnewleg/cmd/wisp/fake131.go` **原文件**（内含 `var sinkAlias131 = installLogSink` +
`sinkAlias131(root)` + `defer sink.close()`）＋它那枚 `main.go`（`81a82,83` 插 `case "fake131":`）⇒ 与上一任逐字节同形。
`go build`/`go vet` rc=0。

- 整包 rc=1，`100 / 52 / 1 / 0`，红名与 X6 **同一枚**：
  `AC#4 RED: leg "fake131" (main.go:82) reaches installLogSink on 1 line(s) of this package and no registered nail names it.`
- 门另外把它新加的边打进日志（这一条上一任没有、本轮才有；它把"我看不见这条边"变成"我接上了并且说出来"）：
  `leg_sink_gate_131_test.go:371: call edges added through package-level function values (R-131-1 形 c): cmdFake131:sinkAlias131->installLogSink`
- 账本 **16 行**，`leg fake131 ... install=true records=true ... -> RED listener installed, no nail`
  ⇒ 上一任最刺眼的那一点（活腿被写成 `install=false records=false -> no records`）已经翻成**真话 + 红**。
- 还原 ⇒ 只差 `wisp.exe`、门 `--- PASS`。

**X11 侧（同一形打在已有钉的 models 腿上）**：别名被接成边之后账本恢复 `install=true`，门 `PASS` 并打印
`modelsEnsure:sinkViaVar131->installLogSink`（这一枚是续单方自报的读数；我在 §4 的 y3 那一发从另一侧独立撞过同一条路）。

⇒ **判 PASS。** 三形全部从"绿/PASS"翻成"门本人红"：X4 点到 `--diag`、X8 不靠整包坏掉、X12 与 X6 同名。

---

## 4. 第五形（Y3 那一支）独立重走 —— 三处都办了，**其中一处测出真洞**

简报要我：① 复造指名豁免、② 核那把尺的两枚误伤是**在未变异树上量到并修掉**的、③ 它自报"三枚恒真变体读数取自补尺**之前**那版码"
⇒ **这一格由我补跑**。三处逐一，外加两枚我自己造的形状。

### 4.1 ① 指名豁免那一支：现在**红**（它自报"补判据前实测静默 rc=0"，我在最终码上复造）

造法（逐字照它的描述）：门 `leg_sink_gate_131_test.go:956`
`if strings.Contains(fi.doc, legSinkRulingMarker) {` ⇒ `... || fi.name == "cmdSLO" {`，**同时**把
`slo_windows.go:171` 的标记句改名成 `LEG-SINK-DECISION-REMOVED:`（即"账本说裁决过、文件里已经没有裁决句"）。`go build` rc=0。

- 只跑门：`GATE_RC=1`，红名原文：
  `leg_sink_gate_131_test.go:379: AC#4 RED: leg "slo" is booked as ruled at [slo_other.go:34 | slo_windows.go:185], but the marker sentence "WISP-LEG-SINK-RULING:" is in none of those files: the reading came from the instrument's own predicate, not from a ruling anybody wrote down.`
- 整包：`RUN 100 / 52 / 1 / 0` —— **唯一一枚顶层红就是门**，不是"整包坏了才响"。
⇒ **它自报的那一格（旧码 rc=0、账本仍写 `ruled=true` 而标记句已删）在 `bcb03aa` 上不再静默。** 〔独立复现〕

### 4.2 ③ 三枚恒真变体：我在**最终码**上逐枚补跑（它自报这三枚取自补尺之前那版）

| 枚 | 门里被改掉的那一行（＋同时删标记句／叠一发旧红） | 续单方自报（补尺前） | **我在 `bcb03aa` 的读数** |
| --- | --- | --- | --- |
| **A1** | `:956` ⇒ `if !fi.calls[installsFuncName] {`（不装听众就算裁决过） | "不绿：4 枚新红" | **不绿**：`AC#4 RED` 命中 **19 枚**，含 4 枚 `carries a WISP-LEG-SINK-RULING: sentence (...) AND reaches installLogSink`（models/no-args/run/secret），**另加 15 枚第二把尺的红**（`booked as ruled at ... is in none of those files`）；整包 `100/52/1/0`，顶层红仍是门本人 |
| **A2** | `:946` install 谓词 ⇒ `if len(fi.slogSel) > 0 {` ＋ 叠 X1（删 models 那 9 行） | "不绿：3 枚红，X1 原有那枚仍在" | **不绿**：3 枚红，其中 **`AC#4 RED: nail "TestAC2ModelsLeg..." claims leg "models", which no longer installs the persistent sink (main.go:74)` 逐字仍在**（⇒ 旧红没被新谓词冲掉）；另两枚是 `leg "slo" ... reaches installLogSink ... no registered nail` 与 `leg "slo" carries a ... AND reaches installLogSink`；整包 `100/50/3/0` ＋ 1 枚子用例红 |
| **A3** | `:950` records 谓词 ⇒ `if len(fi.slogSel) >= 0 {`（恒真）＋ 叠 X2 | "不绿：每条不发记录的腿各红一枚" | **不绿**：12 枚 `AC#4 RED: leg "--help" ("-v"/"default"/...) books records () through the process logger, installs no listener, and carries no WISP-LEG-SINK-RULING: sentence.`（每条无记录腿一枚）；整包 `100/48/5/0` ＋ 1 枚子用例红 |
| **A0**（对照，仪器一字未动） | 只改 `slo_windows.go` 的标记词 | "门红并点名 `slo`" | **红**，点名 `leg "slo" ... carries no WISP-LEG-SINK-RULING: sentence`（⇒ 那枚标记确实承重，删了必红） |

⇒ 三枚恒真变体在**最终码**上各自撞出新红，**没有一枚静默**；A2 还额外证明"旧红 X1 不被新谓词吞掉"。**这一格由我跑完了。** 〔独立复现〕

### 4.3 ② 那把尺的两枚误伤：在**未变异树**上确已修掉，而且修完仍会红（我用两枚新形状顶它）

- **误伤 1（本文件那枚标记常量被当成裁决句）**：修复手段＝测试文件不进字节扫描。
  我在**未变异树**上从另一侧独立验到它：把标记句只写进一枚 `_test.go`（`logsink_test.go` 末尾追加一句
  `// WISP-LEG-SINK-RULING: planted in a test file`）⇒ 门 `--- PASS`、整包 **`100/53/0/0` rc=0**。
  这条**静默是有道理的、且不可被利用**：AST 那条路同样不收测试文件里的函数（`loadMainPackage131` 对 `isTest` 的非 `Test*` 函数
  直接 `continue`），所以"往测试文件写标记"**伪造不出 `ruled=true`**，也就没有"拿测试文件洗掉一条真裁决"这条路。
  ⇒ 不算缺陷，登记为口径条 `R-131r3-3`。
- **误伤 2（`_windows`/`_other` 被合并成一个 site）**：修复手段＝按 `" | "` 逐段查。
  我先复现它**修完仍然会红**（§4.1 那一发：slo 的 site 正是 `slo_other.go:34 | slo_windows.go:185` 这种合并串，两半都没有标记 ⇒ 红），
  ⇒ 把尺调哑并不是这两枚误伤的唯一出路，这一点我证到了。

**再拿两枚"不是 slo、不是指名豁免"的新形状顶它**（这一步就是简报要的"证明它修完仍会红"）：

- **第二枚恒真豁免（换一形，不针对 slo）**：`:956` ⇒ `... || len(fi.slogSel) > 0 {`（"发过记录就算裁决过"）＋删 slo 标记。
  ⇒ 门红 **9 枚**：既有老那条 `carries a ... AND reaches installLogSink`（4 枚，models/no-args/run/secret），
  也有第二把尺在**非 slo 的腿**上咬下来的
  `leg "models" is booked as ruled at [logsink.go:144], but the marker sentence ... is in none of those files`
  ⇒ 这把尺不是只认 `slo` 的定制补丁。整包 `100/52/1/0`，顶层红仍是门。
- **反向那一支在全新的文件上**：把 `// WISP-LEG-SINK-RULING: prose nobody dispatches to` 追加到 `doctor.go` 末尾
  （自由浮动注释，不进任何函数 doc）⇒ 门红一枚，原文
  `AC#4 RED: doctor.go carries a "WISP-LEG-SINK-RULING:" sentence that no ruled ledger row claims: a ruling nobody's dispatch reaches is prose, not a decision.`
  ⇒ 反向读数也不只认 `slo_windows.go`。整包 `100/52/1/0`。

### 4.4 但是我测出一枚**真的洞**：第二把尺只到"文件粒度"，一枚裸标记就能洗钱

造法（`/tmp/wisp131-r3-y6-filegran`）：在上面那枚"发过记录就算裁决过"的豁免之上，把
`// WISP-LEG-SINK-RULING: a sentence parked in this file by the probe, attached to no function.` **追加到 `slo_other.go` 末尾**，
而 `slo_windows.go` 里那枚真裁决句仍然改走。`go build` rc=0。读数：

- `leg "slo" is booked as ruled ... in none of those files` 那一枚红**消失了**。同一份码、只差这枚裸标记：
  不带裸标记那一发该句命中 **1** 次、带裸标记这一发命中 **0** 次（逐枚 grep 对过，快照 `y5-waiver2` vs `y6-filegran`）。
- 整包仍然红（那枚豁免让另外四条腿各自撞红），**所以今天它不静默**；但
  **"放宽 `ruled` 谓词 ＋ 目标腿的 site 是合并串"这一形**下，只要往合并串里任意一枚文件塞一行裸标记，
  `ruled=true` 就能在没有裁决句的那个文件上拿到背书。

⇒ 这正是续单自己 `next=` ② 写的那格（"两个文件各写半句这种拼接仍未测"）——**我今天把它测出来了，方向确凿（红 → 不红）**。
严重度我判**中低**：它需要有人**同时**放宽 `ruled` 谓词并塞一枚裸标记，而被放宽的谓词本身会先把另外四条腿弄红（不是免检），
所以它**不是** AC#4 声称要防的那个结局（"新增一条腿不给钉、门静默"）⇒ **不构成退回触发线**。登记为 `R-131r3-1`。

---

## 5. Linux 侧"如实红"复走（原生容器，不是交叉编译）—— **PASS**

### 5.1 仪器与挂载证明

`git archive bcb03aa` 快照拷到 `D:\work\build-tmp\wisp131-acc-r3`，用
`MSYS_NO_PATHCONV=1 docker run -v /d/work/build-tmp/wisp131-acc-r3:/src ... golang:1.27`（**真 `/d/...` 形，避开 `docker -v C:\…` 静默空挂那一枚假绿**），
进容器**第一件事就是证挂载非空**，退出码一律 `cmd > log 2>&1; echo RC=$?`（**不碰 `${PIPESTATUS[0]}`，dash 下会 Bad substitution**）：

```
--- MOUNT PROOF ---
-rwxrwxrwx 1 root root 883 Sep 23 08:08 /src/go.mod          # 非空挂
--- PLATFORM ---  Linux / x86_64 / GOOS=linux GOARCH=amd64 CGO_ENABLED=1
--- GO VET cmd/wisp ---  VET_RC=0
go: downloading github.com/k2-fsa/sherpa-onnx-go v1.13.8
go: downloading github.com/k2-fsa/sherpa-onnx-go-linux v1.13.8
go: downloading modernc.org/sqlite v1.59.0  ...（另有 7 枚真下载）
```

⇒ **第 0 件成立**：`vet: cmd/wisp/leg_sink_nail_131_test.go:127:12: undefined: sinkInstallRecord` 这一行在原生 Linux 上**不再出现**，
且 rc=0 **不是**"包没加载所以零检查"那种假绿（它先真下载了 sherpa-onnx-go-linux 等模块才报 0）。

### 5.2 那扇门在 Linux 上不是恒绿，是**如实红**（原文）

`go test -count=1 -v -run TestAC4EveryLegIsNailedOrRuled ./cmd/wisp/` ⇒ `GATE_RC=1`，`--- FAIL`，红名原文（逐字）：

```
    leg_sink_gate_131_test.go:275: AC#4 RED: leg "models" (main.go:74) reaches installLogSink on 1 line(s) of this package and no registered nail names it.
    leg_sink_gate_131_test.go:275: AC#4 RED: leg "no-args" (main.go:51) ...
    leg_sink_gate_131_test.go:275: AC#4 RED: leg "run" (main.go:61) ...
    leg_sink_gate_131_test.go:275: AC#4 RED: leg "secret" (main.go:72) ...
    leg_sink_gate_131_test.go:397: AC#4 RED (the instrument, not the code): zero nails registered in this test binary, so the gate
        has nothing to reconcile and would pass on an empty list.
        GOOS=linux compiled no nail file: every registerLegNail131 call lives in a *_windows_test.go file, so the 4 leg(s) in the
        ledger that reach installLogSink (models, no-args, run, secret) are UNREAD on this platform, not unnailed - and no row of
        this ledger is coverage here, in either direction.
        Fix: run this package where the nails compile (scripts/wisp-cli-tests.sh, the windows leg), or give this GOOS its own nail
        file and register it. This reading stays red on purpose: a gate that cannot see is not a gate that has seen nothing to complain about.
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.18s)
```

同一发的账本：清单仍是**磁盘源码现算的 15 行**（`leg slo ... ruled=true -> ruled` 也在），只有 `nails=` 那一列空
⇒ 简报担心的那一格（"把钉文件挪进 `_windows_test.go` ⇒ 门在 Linux 变成恒绿"）**没有发生**；
而且它把"看不见"与"没得抱怨"分开写死了，那四行 `no registered nail` 之上还压了一条自报仪器的红。〔独立复现〕

### 5.3 交叉编译那枚 rc=1 我按"逐错误行"归因，没有整体归因

Windows 上 `GOOS=linux go vet ./cmd/wisp/` ⇒ rc=**1**，整份输出 **3 行**，逐行归因：

```
1  package github.com/CarlosShao/wisp/cmd/wisp                      <- 包名头（导入链上下文，不是诊断行）
2      imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx          <- 导入链上下文
3      imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in ...sherpa-onnx-go-linux@v1.13.8   <- 唯一一枚诊断，落在第三方模块目录
```

`grep -nE 'cmd/wisp/[^ ]*\.go'` 对这份输出 ⇒ **rc=1（零命中）**，即**没有任何 `cmd/wisp/<file>.go:line:col` 形诊断**。
⚠ 但"唯一错误行只指向 sherpa 模块目录"这句要加一个限定才严谨：`cmd/wisp` **确实出现在第 1 行**，那是**包名头**而不是诊断行。
⇒ 归因结论不变（工具链约束、非本票引入、Linux 原生 rc=0 为反证），**但引用那句话的人别把它读成"输出里找不到 cmd/wisp 字样"**。
这一格我按上一任踩过的坑办：不整体归因，逐行点名。

---

## 6. 票 130 越权改 131 的 helper：追认成不成立（我有一票否决）

三发我都**在最终码 `bcb03aa` 上独立复跑**（续单方自报 P1/P2 的整包四数取自补第二把尺之前那版，那一格它登记为"待核"，现在我跑完了）。

| 探针 | 造法 | 只跑目标用例 | **整包 `-count=1 -v` 四数与红名** | 与续单自报是否同账 |
| --- | --- | --- | --- | --- |
| **P1** 删掉 130 的冲刷 | `logsink.go:173` ⇒ `var earlyFlushed, earlyDropped int`（`grep -c FlushEarlyLogRecords` = 0），build/vet rc=0 | models 钉 **`--- PASS`**（rc=0） | rc=1，`100/50/3/0`：`TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord`（130 断言 1）、`TestAC1ResidentLegInstallsItsLogListenerOnDisk`、`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`（127 两枚钉） | **同账，逐字** |
| **P2** 原始疾病：某腿的事件早于它自己的 install | `models.go` 加包级 `func init() { slog.Info("models: P2 planted event ...") }` + `log/slog` import，build/vet rc=0 | models 钉 **`--- FAIL`**：`leg_sink_nail_131_windows_test.go:366: record 1 sits before the listener's booking and is "models: P2 planted event booked before this leg had a listener", want "winsec: sealing path resolver installed" or nothing ...（记录序列里 `wisp: persistent log sink installed` 排在第 3） | rc=1，`100/50/3/0`：models 钉 + 127 那两枚 | **同账**（这一发的四数取自**最终码**，续单那版是补尺前） |
| **P3** 把 helper 第一条判据换回 130 之前的形状、打在**健康树**上 | `for i, r := range s.records[:installIdx] {...}` ⇒ `if installIdx != 0 { t.Errorf("install record is at index %d, want 0: ...") }` | — | rc=1，`100/52/1/0`，唯一一枚红＝`TestAC2ModelsLegBooksItsHandOffVerdictOnDisk`（`leg_sink_nail_131_windows_test.go:362: install record is at index 1, want 0`）；**`TestAC3SecretLeg...` 与门都 `--- PASS`** | **同账**（红谁取决于顺序，不是取决于秩序） |

### 6.1 正面回答：**"追认成立"在 P1 前提失效之后仍然站得住，但依据必须换成 P2 单发**

先把话说清：简报里"P1 ⇒ helper 仍须红"这条**前提**我复跑证否了（P1 下 131 的钉全绿）。但**前提错 ≠ 结论错**，
因为"P1 这一发要回答的问题"从来不是"helper 还红不红"，而是**"130 有没有把它自己的修法偷渡成 131 的必要条件"**。
按这两问分别看：

1. **有没有把 131 的侦测力改弱？** 判据是 **P2**：原始疾病（本腿的事件排在 booking 之前）在改后的 helper 下**照样红**，
   红名点到 `record 1 sits before the listener's booking`。⇒ 这一改保住了 helper 存在的理由。**成立。**
2. **有没有把 130 的回放偷渡成 131 的前提？** 判据是 **P1 的实测结果本身**：把 130 的冲刷删掉，131 的 helper 与两枚钉**全 PASS**，
   红的是 130 自己的断言与 127 的钉。⇒ 131 的判据**不依赖** 130 的回放在位（helper 的注释也逐字声明"它不主张 boot 记录存在"，
   并把那一形的证人指到 `early_log_nail_130_windows_test.go` 与 127 的文件）。**没有偷渡。成立。**
   ⚠ 这一条**恰恰是被"P1 应红"那个错前提挡住的**：若真按简报预期红起来，反而说明偷渡发生了。
3. **旧形状能不能就这么留着？** 判据是 **P3**：换回"下标必须是 0"之后，**健康的树**上 models 钉红、secret 钉绿
   ⇒ 旧判据红不红取决于测试顺序而不是秩序，它不是"更强的同一件事"，是"另一件会误报的事"。⇒ 130 动它**有正面必要**。

⇒ **我的裁决：追认成立（不行使否决票），但成立的全部重量在 P2 与 P3 上，P1 只是"未偷渡"的反证，不是"仍侦测"的正证。**
续单证据那段结论的**判定**与我一致，它把 P1 列作依据之一（"它没有把 130 的回放当成新的必要条件"）其实已经是正确用法，
只是与简报的错前提并排写会让下一个人误读成"两发都支持侦测力"。⇒ 登记 `R-131r3-2`（归因/表述级）。
**撤销令不在我这一程**：我没有回滚任何一枚文件，`9b5d64d` 那处改动原地未动（我只在 `/tmp` 的一次性副本里做过 P1/P2/P3）。

---

## 7. 结案判据 3：四数不降 —— **PASS**（含一处归因更正）

⚠ 简报说对了要害：**上一任的 `82/82` 已被票 129/130/续单三次移动，谁拿它当分母谁就错**；
但"重新基线"必须**逐数解释过**。我 therefore 把**两枚锚点各自都自己跑了一遍**，不引用任何一方的自述：

| 仪器 | 上一任锚 `56d8026`（**我在同一台机器复测**） | 被验版本 `bcb03aa`（**§1，我复测**） | 逐数判定 |
| --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | RUN **164** / 顶层 PASS **92** / FAIL 0 / SKIP 0，rc=0 | RUN **200** / 顶层 PASS **106** / FAIL 0 / SKIP 0，rc=0 | RUN +36、PASS +14、FAIL/SKIP 原地 ⇒ **不降** |
| `go test -count=1 -v ./cmd/wisp/` | RUN **82** / 顶层 46 / 0 / 0 | RUN **100** / 顶层 53 / 0 / 0 | +18 / +7 ⇒ **不降**，且 18 枚逐名可点（§1.1） |
| `sh scripts/wisp-cli-tests.sh`（CI 形状） | 同上（脚本口径与 `-count=1 -v` 同数） | RUN 100 / 53 / 0 / 0，rc=0 | **不降** |
| `sh scripts/d22scan.sh` 八 scope | **202 / 22 / 40 / 18 / 16 / 40 / 387 / 34**，rc=0 | **203 / 22 / 40 / 18 / 16 / 40 / 390 / 37**，rc=0 | 八数**逐枚 ≥**，一枚没降 |
| `gofmt -l` / `gofumpt -l cmd/wisp/` | — | 空 / 空，rc=0 | 同 |
| `go vet ./cmd/wisp/`（windows）/ `go build ./cmd/wisp/` | — | rc=0 / rc=0 | 同 |

⚠ **一枚口径差**（写给下一位引用的人）：上一任 §2 写的是"`RUN 164 / PASS 164`"，我在同一枚树同一条命令量到
"RUN **164** / **顶层** `--- PASS` **92** / 子用例 72" ⇒ 92+72=164。两个数都是真的，**差在 PASS 那列有没有把子用例算进去**。
引用"四数"不写口径就已经腐坏过一次（`79→83` 那一族同病）。本文一律用**顶层**口径，子用例另列。
`cmd/` 那一格 34→37 的三枚增量我逐枚对过＝`dataroot_128_test.go`、`dataroot_128_windows_test.go`、
`early_log_nail_130_windows_test.go`（票 128 两枚 + 票 130 一枚），**本票三形判据一枚 `cmd/` 新文件也没加**（它只改在已有文件里）。

⇒ **判 PASS。** 基线移动被逐数、逐名解释完，零枚被洗小；归因名字要更正（`R-131r3-4`）。

---

## 8. `R-131r3-x` 清单（严重度 / 能否复现 / 修法 / **归谁**）

| 号 | 现象 | 严重度 | 能否复现 | 修法 | **归谁** |
| --- | --- | --- | --- | --- | --- |
| **R-131r3-1** | **第二把尺（`rulingCrossChecks131`）只到"文件粒度"**：一枚**不绑任何函数**的裸标记句，只要落在被 `rulingSites` 合并串点到的**任意一枚**文件里，就能给一条**已经不存在的裁决**拿到背书。实测同一份放宽码（`ruled` 谓词 ⇒ `... || len(fi.slogSel) > 0`）＋ slo 真裁决句改走：不带裸标记 ⇒ `leg "slo" is booked as ruled ... in none of those files` 命中 **1**；把裸标记追加到 `slo_other.go` 末尾 ⇒ 命中 **0**（红被洗掉）。⇒ 这正是续单自己 `next=` ② 写"仍未测"那一格，**今天测出来了，方向确凿** | **中低**（今天不静默：放宽谓词本身会把另外四条腿各自弄红，所以它**不是** AC#4 声称要防的结局 ⇒ **不触发退回硬线**；但"两枚变异叠起来洗掉一枚红"这一族它已经占了一次） | **能**：`/tmp/wisp131-r3-y5-waiver2` 与 `/tmp/wisp131-r3-y6-filegran` 两棵树逐字对照；红名命中数各 grep 一次即得 | 两条任选：① `rulingFiles` 从"文件含标记词"收到"**标记句所在的注释块绑到哪一枚 func**"（与 `leg.rulingSites` 的 `file:line` 同段才算背书）；② 按票 135 的做法给这把尺装**一枚主动弄哑自己**的自证腿（把那枚裸标记做成常备用例 ⇒ 现在就红）。我倾向 ②，因为 ① 会把尺重新拉回"只信 AST"那一侧、又少一把独立读数 | **票 135**（排队中，`worker-ticket124-ac2a` 不在飞它）——它的 AC#1/AC#2 就是"两枚变异叠起来 16/16 全绿"与"自证腿不许是哑的"，同一根。⚠ 与 **票 133 AC#3**（`ready-for-agent`，不在简报给的排队名单里）也同族，**并到哪一枚由编排者裁，别让它两票都不接** |
| **R-131r3-2** | **追认那段结论的"依据列"会把下一个人带偏**：`131-followup-2` §2 结论把 P1 与 P2 并排列作依据，而 P1 的**实测结果**（131 的 helper 不红）在简报那句错前提下读起来像"侦测力掉了"。实际两问要分开：**侦测力＝P2（红）**，**没偷渡＝P1（不红才是好消息）**，**旧形状留不得＝P3（健康树上就红）** | 低（措辞／归因；**判定不变**，我这一程的正面回答见 §6.1：追认**成立**） | 能：§6 那三发逐发可重跑（快照 `/tmp/wisp131-r3-p1`、`-p2`、`-p3`） | 在 `131-followup-2` §2 结论段 **append 一句**（票面 append-only 同理，原文不抹）："P1 的角色是'未偷渡'的反证，不是'仍侦测'的正证；后者只有 P2" | **编排者**（改的是两份证据/票面里的一句话，不是码；`9b5d64d` 那处改动**原地不动**，撤销令只服从他的对话） |
| **R-131r3-3** | **口径条，不是缺陷**：标记词只出现在 `_test.go` 时**两把尺都不响**（第二把尺按 `!isTest` 排除；AST 那条路 `loadMainPackage131` 对测试文件里的非 `Test*` 函数直接 `continue`）。我实测：往 `logsink_test.go` 末尾追加一枚裸标记 ⇒ 门 `--- PASS`、整包 `100/53/0/0`。⇒ **不可被用来伪造 `ruled=true`**（伪造需要一条生产函数带这句），所以这一格今天没有牙也咬不到人 | 极低（**登记给下一位**，别当成已收的洞） | 能：`/tmp/wisp131-r3-y8-testfile` | 无需动作；若 R-131r3-1 收成"标记句必须绑到 func"，这一格会被同一条修法顺带关掉，**届时改一句注释即可** | **无人为动作**（与 R-131r3-1 同批收的话归 135） |
| **R-131r3-4** | **两处引用/归因腐坏**（都在"基线移动"那笔账上）：<br>① "+18 RUN / +7 顶层 PASS **全归票 129/130**"——实为**票 128（`4e5d240`，16 枚＝5 顶层＋11 子用例）＋票 130（`d87905c`，2 枚）**；`git log 56d8026..b723978 -- cmd/wisp` 里**没有 129**，129 在 `cmd/wisp` 零枚用例。<br>② "六发红名逐字不变，**只有门与钉自己的 `file:line` 前缀下移**"——不完整：被驱腿自己的行号也移了（`models.go:276→284`、`secret.go:219→226`、`run.go:164→178`、`logsink.go:139→144`；X2 红名里的 `slog.Info@secret.go:321→:328`）。红名文本与所点的腿一字未变，但位移不止仪器前缀。<br>③ 同一件事的枚数在票面写"各多两枚"、A122① 写"各多一枚"，**实测两枚**（X2 与 X9） | 低（**判定不受影响**：四数不降成立、六发红名成立；但"引用某区间几枚红必须连口径一起引"这一条老账又踩了一次） | 能：§1.1 的 `comm` 两段命令 + §1.2 的 `git diff --numstat`／`git grep -n` 即得 | 票面与 A122① 各 **append 一行更正**（不抹原文）；把"四数不降"的证据一律改成"逐名 `comm`"而不是"逐数差" | **编排者**（账目与票面）；`cmd/wisp` 侧**无需改码** |
| **R-131r3-5** | **"交叉 `GOOS=linux go vet` 的唯一错误行只指向 sherpa 模块目录"这句，字面读会翻车**：那份输出第 **1** 行就是 `package github.com/CarlosShao/wisp/cmd/wisp`（导入链的包名头，不是诊断行）。零枚 `cmd/wisp/*.go:line:col` 形诊断是真的，"输出里找不到 `cmd/wisp`"是假的 | 低（今天这三个代理的结论都对，但这一句被下一个人按字面复述就会造成一次误否/误归） | 能：§5.3 那三行 + `grep -nE 'cmd/wisp/[^ ]*\.go'` 一次即得 | 引用时改成"唯一**诊断行**落在第三方模块目录，另有两行导入链上下文"；把 §5.3 那张逐行表一起带走 | **编排者**；`internal/observe`／`logsink.go` 侧无动作 |

**归单一句话**：本轮**五枚全部不需要在飞的 `acceptor-ticket134-r1` 或 `worker-ticket124-ac2a` 停手**——
R-131r3-1 归**票 135**（唯一一枚改码项，`cmd/wisp/leg_sink_gate_131_test.go` 一枚函数或新增一枚自证用例），
R-131r3-2/4/5 归**编排者自己的账目与票面**（append 更正，零枚代码），R-131r3-3 无人需动作。
**我没有新开票、没有翻 AC#4 那一格的勾。**
