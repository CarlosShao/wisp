# 票 131 独立对抗验收（第二任）—— 攻的就是那枚"从 `main.go` 现算腿清单"的枚举门

**验收方**：`acceptor-ticket131-run2`（本文件唯一作者；只读验收，本文件是唯一写件）
**日期**：2026-09-23（本会话，开件 11:45，末次读数 12:2x；每裁一格落一次盘，见本文件的两枚 commit）
**被验交件**：票 131，实现方 `agent-ticket131`
**被验 sha**：`56d8026`（代码态终点；四枚 `7e60d31`→`627aa66`→`013e703`→`56d8026`）
**控制组**：`db9fafc`（票 131 开工前的树，实现方自选，我沿用其口径）

---

## 0. 接管声明（我不继承前任骨架的一个字）

本文件在原路径上被前任验收代理（`Accept ticket 131 enumeration gate`）建过一枚 34 行骨架：
只有"一～九"的标题清单、九个框全部未勾、零裁决。它的转录最后一次写入是 09:16，距本轮开件约 2.5 小时，
无交件通知。**本轮把那 9 个空标题整体删除、换成下面这套我自己的结构**，骨架里没有一个字被继承为结论。
前任的标题里有三格（"六、`627aa66` 单判"、"七、`56d8026` 单判"、"四 (c) 整段拆 install"）与本轮的
关注点重合，我在下面用**我自己的编号与判法**重做，不是填它的空。

本轮没有重做前任的任何事（它什么都没做），也没有相信实现方的任何一句话：
下面每一条读数都是我在这台机器上、在被验 sha 的纯净快照里自己跑出来的。

---

## 1. 仪器与被验版本的锚定（先证明我读的不是脏树）

```
git archive 56d8026 | tar -x -C /tmp/wisp131-acc131r2          # 目录带本会话后缀
cp third_party/sherpa-onnx/*.dll <快照>/third_party/sherpa-onnx/ # git archive 不含未入库件
PATH=<快照>/third_party/sherpa-onnx <所有 go 命令>
```

- **完整性反查（12:0x 自己做的）**：另起一枚 `/tmp/wisp131-verify-acc131r2` 从 `git archive 56d8026` 现取，
  `diff -r --exclude=third_party` 两棵树 ⇒ **唯一差异是我 `go build` 落下的 `wisp.exe` 产物**（已删）。
  基线快照的全部源文件与 `56d8026` 逐字节相同。
- **仓库工作树**：全程未在仓内建 worktree、未改任何生产文件（`git status --short` 只有本文件；见 §11）。
- **A103④ 三连检查（取任何读数之前，11:50）**：
  ① `gh run list --branch dev --status in_progress` ⇒ 有 1 枚在飞（`35815467601`，dev push）；
  ② `E:/work/base/actions-runner/_diag/Worker_*.log` 最新 mtime = **11:47:02**（开测前 3.5 分钟）；
  ③ `Get-CimInstance Win32_Process | Where ExecutablePath -like '*actions-runner\_work*'` ⇒ **零枚进程**；
  ④ 追一步：`gh run view 35815467601 --json jobs` ⇒ **`slo-full`（本机 self-hosted、采 D32 那两个数的那一步）已 `completed success`**，
  在飞的只剩 `test-windows`（GitHub 托管，不落本机 CPU）。
  ⇒ 判定：**时序采样负载已散，本轮可以开测**；但仍**不采 RSS/时序结论**，本文所有墙钟只作旁证并标"带噪"。
- ⚠ **一枚仪器坑（本轮自己踩的，写给下一位）**：我在中途用**脏工作树**（`dbbc822`，票 130 已改过 `cmd/wisp/logsink.go`）
  `grep -n installLogSink cmd/wisp/*.go` 取行号，读到的 `secret.go:226` 与快照里的 `219` 差 7 行。
  **票面与快照的行号一致、脏树不一致** ⇒ 凡引用行号必须从快照取；脏树连行号都不能信。

---

## 2. 基线（`56d8026` 纯净快照，未变异）

| 仪器 | 实现方自报 | **我复算** |
| --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | RUN 164 / PASS 164 / FAIL 0 / SKIP 0，`ok` 86.648s | **RUN 164 / PASS 164 / FAIL 0 / SKIP 0，rc=0**，`ok` 95.960s（带噪） |
| `sh scripts/wisp-cli-tests.sh`（CI 形状） | RUN 82 / PASS 46 / FAIL 0 / SKIP 0，rc=0 | **RUN 82 / PASS 46 / FAIL 0 / SKIP 0，rc=0**，`ok` 44.098s（带噪） |
| `sh scripts/d22scan.sh` | rc=0，八 scope 202/22/40/18/16/40/387/**34** | **rc=0，同八数一字不差**（`ban #8 cmd/=34`） |
| `gofmt -l cmd/wisp/` | 空 | **空**，rc=0 |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/` | 空 | **空**，rc=0 |
| `go vet ./cmd/wisp/` / `go build ./...` | rc=0 | **rc=0 / rc=0** |
| `GOOS=linux go vet ./cmd/wisp/` | rc=1，原文是 sherpa 那句 | **rc=1**，原文逐字：`imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in .../sherpa-onnx-go-linux@v1.13.8` |

⇒ 四数与台账八 scope **逐字复现**，两枚 FAIL=0 ⇒ 票 123 那四枚红确实不在本包分母里（与实现方说法一致）。〔独立复现〕

**控制组我也自己取了一枚**（AC#5 那句"与同 sha 控制组逐数不降"要有我的数，不是抄它的）：
`git archive db9fafc | tar -x -C /tmp/wisp131-ctrl-acc131r2` + `sh scripts/d22scan.sh` ⇒ rc=**0**，
八 scope = **202/22/40/18/16/40/387/32**。对照本表候选组 `…/34` ⇒ **七数一字不降、`cmd/` +2**，
那两枚增量正好是被验的两枚新文件（`leg_sink_gate_131_test.go`、`leg_sink_nail_131_test.go`）。
⇒ 实现方那张对照表的**两侧我都独立取过**，同数。〔独立复现〕

**门每次运行现算出的腿清单**（基线 `-v` 原文，两趟逐字相同；这就是 AC#1 那张表的可执行形式）：

```
leg --help       main.go:90   install=false records=false ruled=false nails=-            -> no records
leg --version    main.go:87   install=false records=false ruled=false nails=-            -> no records
leg -h           main.go:90   install=false records=false ruled=false nails=-            -> no records
leg -v           main.go:87   install=false records=false ruled=false nails=-            -> no records
leg default      main.go:93   install=false records=false ruled=false nails=-            -> no records
leg doctor       main.go:67   install=false records=false ruled=false nails=-            -> no records
leg help         main.go:90   install=false records=false ruled=false nails=-            -> no records
leg models       main.go:74   install=true  records=true  ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk -> nailed
leg no-args      main.go:51   install=true  records=true  ruled=false nails=TestAC1ResidentLegInstallsItsLogListenerOnDisk -> nailed
leg panel-assets main.go:84   install=false records=false ruled=false nails=-            -> no records
leg providers    main.go:64   install=false records=false ruled=false nails=-            -> no records
leg run          main.go:61   install=true  records=true  ruled=false nails=TestAC2SealNoticeLandsInTheRunLegLogFile -> nailed
leg secret       main.go:72   install=true  records=true  ruled=false nails=TestAC3SecretLegBooksItsAuditRecordsOnDisk -> nailed
leg slo          main.go:82   install=false records=true  ruled=true  nails=-            -> ruled
leg version      main.go:87   install=false records=false ruled=false nails=-            -> no records
```

15 行 / 11 组 = 票面 AC#1 那张表，**行号逐条对得上**（含 `no-args main.go:51`、`slo main.go:82` 那两枚特殊形）。
install 全集我也单独 grep 过（快照、非测试）：定义在 `logsink.go:139`，调用者 4 枚 =
`run.go:164`、`resident_windows.go:57`、`models.go:276`、`secret.go:219` —— 与票面"调用者 4 枚"一致。

---

## 3. 复现实现方的两发变异（证明门不是恒绿）

每发固定三步：先 `diff` + `grep` 证那几字节真没了 ⇒ `go build ./cmd/wisp/` rc=0 ⇒ 才读红名。
变异全在 `/tmp/wisp131-acc131r2-<编号>` 的一次性副本里做，基线快照与工作树都没被污染（§1 的反查）。

### X1 ＝ 实现方的 M1（拆掉 `cmd/wisp/models.go:276-284` 那 9 行）

- 落地：`diff` 只有 `276,284d275`，被删的正是 `if sink, sinkErr := installLogSink(io_.dataDir); ...` 那 9 行；
  `grep -n "installLogSink\|持久日志未启用\|defer sink.close" cmd/wisp/models.go` ⇒ **零命中，rc=1**；`wc -l` 326→317。
- `go build` **rc=0**、`go vet` **rc=0**（编译器仍然不知道发生了任何事）。
- `go test -count=1 -v ./cmd/wisp/` ⇒ rc=**1**，`RUN 82 / PASS 78 / FAIL 4 / SKIP 0`：
  **3 枚顶层红 + 1 枚子用例红**，与实现方 M1 的枚数、名字逐字相同：

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.02s)
    leg_sink_gate_131_test.go:214: AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "models",
        which no longer installs the persistent sink (main.go:74): the nail is now proving something else, or nothing
--- FAIL: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
    leg_sink_nail_131_test.go:307: no wisp-<day>-<seq>.jsonl in <数据根>\logs: the leg ran to completion and
        its listener wrote nothing there.
--- FAIL: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models (0.01s)
    leg_sink_nail_131_test.go:455: a log directory that cannot be opened produced no named refusal
        ("wisp models: 持久日志未启用" on stderr). ...
```

⇒ **本票立项的那一发（"删那 9 行 76 条一条不红"）今天会红，且红名点到 `models` 那条腿。** `R-127-6` 的正面判据成立。〔独立复现〕

### X2 ＝ 实现方的 M5（拆掉 `cmd/wisp/secret.go:219-223` 那 5 行，本票新装的）

- 落地：`diff` 只有 `219,223d218`；`grep -c installLogSink cmd/wisp/secret.go` ⇒ **0**。
- `go build` rc=0、`go vet` rc=0；测试 rc=1，`RUN 82 / PASS 78 / FAIL 4 / SKIP 0`（3 顶层 + 1 子用例）。
- **枚数比实现方自报的"三枚红"多一枚子用例**，方向是"更多枚在查"，不是判据变松；红名逐字：

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.02s)
    leg_sink_gate_131_test.go:194: AC#4 RED: leg "secret" (main.go:72) books records
        (slog.Info@secret.go:321, slog.Info@secret.go:461, slog.Warn@secret.go:461) through the process logger,
        installs no listener, and carries no WISP-LEG-SINK-RULING: sentence.
    leg_sink_gate_131_test.go:214: AC#4 RED: nail "TestAC3SecretLegBooksItsAuditRecordsOnDisk" claims leg "secret",
        which no longer installs the persistent sink (main.go:72) ...
--- FAIL: TestAC3SecretLegBooksItsAuditRecordsOnDisk
--- FAIL: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret
```

⇒ 门对 `secret` 给出的是**两条独立读数**（要求裁决 + 要求钉），这正是 AC#3 想要的形状。〔独立复现〕

**两发都复现成功 ⇒ "门不是恒绿"这一条我背书。**

---

## 4. 我自己补的变异（这一节才是本验收的正文）

候选方向按派单给的四条：① 保留 import 只拆分发那一跳；② 让解析器吃不到某条腿 ⇒ 门**静默少一条**；
③ 钉装到错误的腿上；④ 新增一条腿但不登记。

### X3（候选 ①，同时是实现方自认盲区 `M-B1` 的独立重走）—— **红，只红在门上**

做法逐字按 `R-121-1` 的 N-3 形：`main.go` 里删掉 `case "models":` 那**整支 8 行**（含 `os.Exit(cmdModels(...))`
那一跳），`cmd/wisp/models.go` 与它的 import 一个字节不动。

- 落地：`diff` 只有 `74,81d73`；`grep -rn cmdModels cmd/wisp/*.go | grep -v _test.go` ⇒ 只剩
  `models.go:98`（注释）与 `models.go:104`（声明）⇒ **生产调用者 0 枚**；`models.go` 与快照逐字节相同。
- `go build` **rc=0**、`go vet` **rc=0**。
- `go test -count=1 -v ./cmd/wisp/` ⇒ rc=1，`RUN 82 / PASS 81 / FAIL 1 / SKIP 0`，唯一一枚红：

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.02s)
    leg_sink_gate_131_test.go:214: AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "models",
        which main.go does not dispatch to any more
```

- **同时**：`TestAC2ModelsLegBooksItsHandOffVerdictOnDisk` **PASS**（`-v` 原文打出 2 条记录：
  `INFO/wisp: persistent log sink installed` + `INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error`，exit=1）。
  ⇒ 那枚钉是进程内驱动 `cmdModels` 的，**它自己看不见分发有没有断**；拦住这一形的只有新门。

**判**：`R-121-1` 的 N-3 形 / 实现方自认的 `M-B1`（"import 保留 + 分发拆掉 ⇒ 绿"）**在本票的新门下不再成立**——
这一形今天会红，且红名点到"main.go 不再分发到 models"。
它说"触发不了"我没接受，我按字面重走了：结论是**它的自认盲区被它自己的门补掉了**（对**有钉的腿**而言）。〔独立复现〕
⚠ 但这一形的**成立条件要写清**：门红是因为"有一枚钉登记在 `models` 这个名字上"。
对**没有钉、也不 install、也不发 `slog` 记录**的腿（`providers`/`doctor`/`panel-assets`/`version`/`help`/`default`），
把它们的分发那一跳摘掉，本包**零红**——票 133 要补的那把"入口级可达性"尺**仍然空着**，本票没声称补它，我也不会替它记成补了。

### X4（候选 ②，形 a）—— **全绿。本票的退回理由第一条**（最省的那形是 X12，只要一行）

做法（一发两拍，都在快照副本里）：在 `main()` 里、`switch args[0]` **之前**插一枚早退分支，
并新增一枚生产文件 `cmd/wisp/diag131.go` 装听众：

```go
	if len(args) > 0 && args[0] == "--diag" {
		os.Exit(cmdDiag131(args[1:]))
	}
```
`cmdDiag131` 调 `installLogSink(root)` + `defer sink.close()`（其余两枚钉的形状照抄）。**不给它任何钉、不碰任何测试文件。**

**先证这条腿是活的（不是我造的一枚死代码）**——`go build` 出的**真二进制**：

```
$ /tmp/wisp131-x4.exe --diag C:\...\x4data
time=... level=INFO msg="wisp: persistent log sink installed" dir=C:\...\x4data\logs min_level=info
diag ok                                                                          # rc=0
$ ls C:\...\x4data\logs\           →  wisp-20260923-001.jsonl
$ cat 那枚 jsonl                   →  {"...","msg":"wisp: persistent log sink installed","dir":"...x4data\\logs",...}
```

**第一拍（装了听众、没给钉）读数**：`go vet` rc=0；`go test -count=1 -v ./cmd/wisp/` ⇒
**rc=0，`RUN 82 / PASS 82 / FAIL 0 / SKIP 0`，`ok` 41.8s**，门的账本仍是 **15 行**，
`grep -c diag` 在整份 `-v` 输出里 = **0** ⇒ **枚举门连"少了一条腿"这件事都不报**，
它只报它看见的那 15 行，看见与否取决于分支写成什么形状。

**第二拍（把这条腿的 install 拆掉，9 行级）**：`grep -c installLogSink cmd/wisp/diag131.go` = **0**、
`go build` rc=0、`go vet` rc=0；真二进制 `--diag` 仍 **rc=0** 但 `<root>\logs` 里**一个文件都不落**；
`go test -count=1 -v ./cmd/wisp/` ⇒ **rc=0，`RUN 82 / PASS 82 / FAIL 0 / SKIP 0`，全绿**。

⇒ **票 131 的立项理由逐字复现了第九次**："装了一条腿的听众、拆掉它、`cmd/wisp` 一条用例都不红"——
只差这条腿写在 `switch` 而不是写在 `if` 里。而票面第 18 行要求门回答的那句
"**新增一条 CLI 腿时，用例会不会自动要求我给它一枚钉**"，我量到的答案是：**不会**（形状不对就免检）。
这条**不在**实现方"残窗"那一段的自陈范围内——它自陈的是"records 谓词只认 `slog.*`/`observe.InitLog`，
toast-only / 纯 stderr 的腿判无义务"，那讲的是**账本里那条腿的判定材料**；我这一发讲的是
**腿根本没进账本**，是另一格。〔独立复现〕
⚠ 我还顺手量了"这条修法今天会不会误伤"：`grep -rn os.Args cmd/wisp/*.go`（非测试）在 `56d8026` 上**只有 `main.go:50` 一处命中**
⇒ "除已归类分支外不许再读 `os.Args`"这条断言**今天零误伤**，是一行量级就能补的洞，不是架构问题。

**可重跑**：`/tmp/wisp131-acc131r2-x4-ifleg`（第一拍的 `diag131.go` 已被我就地改成第二拍；两拍的
`main.go` 差 3 行都在 `diff` 里）。补丁形状＝上面两段代码，任何人 5 分钟能重跑。

### X5（候选 ③：钉装到错误的腿上，数量对、归属错）—— **全绿**

做法：只改 `leg_sink_nail_131_test.go:573-574` 两行的**腿键**，把 models 的钉登记到 `secret`、
secret 的钉登记到 `models`（`diff -r` 只有那 2 行 2 列）。用例本体、install 点、门一个字没动。

- `go build` rc=0、`go vet` rc=0；`go test -count=1 -v ./cmd/wisp/` ⇒ **rc=0，`RUN 82 / PASS 82 / FAIL 0 / SKIP 0`，全绿**。
- 门的账本原文（两行归属已经互换，仍各报 `nailed`）：

```
leg models  main.go:74  install=true records=true ruled=false nails=TestAC3SecretLegBooksItsAuditRecordsOnDisk  -> nailed
leg secret  main.go:72  install=true records=true ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk -> nailed
```

⇒ **门只核"有没有一枚读了盘的钉挂在这个名字上"，不核"那枚钉驱的是不是这条腿"**。
损害边界我也量清了（这条不是最重的）：即使归属错乱，**拆掉任一 install 仍会红**——因为 `staleRegistrations131`
看的是"腿还在不在 install"，与那枚钉驱的是谁无关。丢的是"这条腿其实没有一枚真正测它的用例"这一格：
账面 `nailed`，实际那条腿的语义（记录内容、落点、退出码）无人守。⇒ `R-131-2`。〔独立复现〕

### X6（候选 ④：新增一条腿但不登记）＝ 实现方 M3 的独立复现 —— **红，只红门**

做法：新增生产文件 `cmd/wisp/fake131.go`（`cmdFake131` 调 `installLogSink` + `defer sink.close()`），
并在 `main.go` 的 `switch` 里插 `case "fake131": os.Exit(cmdFake131(args[1:]))`（`diff` 只有 `81a82,83` 两行）。
**不写任何测试、不碰任何注册。**

- `go build` rc=0、`gofmt -l` 空、`go vet` rc=0；测试 rc=1，`RUN 82 / PASS 81 / FAIL 1 / SKIP 0`：

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.02s)
    leg_sink_gate_131_test.go:171: AC#4 RED: leg "fake131" (main.go:82) reaches installLogSink on 1 line(s)
        of this package and no registered nail names it.
        Delete that install block and this package stays green - the reading ticket 127 measured for models, and ticket 117 for resident.
        Fix: write the nail, and claim it with registerLegNail131("fake131", TestYourCase) in the file that owns it.
      leg fake131  main.go:82  install=true records=true ruled=false nails=-  -> RED listener installed, no nail
```

⇒ 与实现方 M3 逐字同形、同枚数。**"往 `switch` 里加一条装了听众的腿"这一形，门的响应是对的**：
点名、给行号、给修法，且没有顺手把别的用例拖红。这条我背书，实现方不是在编。〔独立复现〕
⚠ 但 X4 与 X8 就是这条判据的**形状前提**：腿得写成 `switch` 的字面量 case 才进得了账本。

### X3c（候选 ②的前半：把整个分发换个写法，门会不会静默交出一张短表）—— **响亮地红**

做法：`switch args[0] {` ⇒ `switch strings.ToLower(args[0]) {`（外加一行 `strings` import）。这是"分发还在、
只是标签不再是 `args[...]` 本身"的最常见重构形（顺手大小写归一）。

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.02s)
    leg_sink_gate_131_test.go:152: AC#4 RED (the instrument, not the code): func main has no `switch args[...]`:
        either the dispatch moved, in which case this reader has to be pointed at where it went, or it is gone.
        Neither is a green
```

⇒ 整枚 switch 换形时它**不交空表、不装绿**，而是自报"仪器坏了"。这一格实现方说的是对的（我复算）。〔独立复现〕
⇒ 正因为这样，X4/X8 那两发的价值才被放大：**整体坏掉它知道，逐条漏掉它不知道。**

### X8（候选 ②的第二个独立形状：case 标签写成命名常量）—— **全绿，且那条腿静默出账**

做法：`case "slo":` ⇒ `case sloCmdName131:`，并在 `main.go` 加 `const sloCmdName131 = "slo"`。
行为不变（`go build` rc=0；这一发改的是**分发标签的写法**，不是分发的目标）。

- `go test -count=1 -v -run TestAC4EveryLegIsNailedOrRuled` ⇒ **rc=0、`--- PASS`**。
- 账本原文（`leg slo` 那一行**没了**，取而代之的是两行都叫 `default`）：

```
leg default  main.go:95  install=false records=false ruled=false nails=-  -> no records
leg default  main.go:84  install=false records=true  ruled=true  nails=-  -> ruled
```

⇒ 两个后果，都在这枚门的正中：① `caseLabels131` 只吃 `*ast.BasicLit`，
**标签写成常量/标识符的腿根本不进账本**，门既不红也不报"少了一条"；② 落空的正是**中间那一行**——
`slo` 是唯一真实占用 `WISP-LEG-SINK-RULING:` 的腿，也就是 **R-117-1 的守卫**。
将来一条"发 `slog` 记录但不装听众"的新腿，只要把标签写成常量，就**不需要裁决**、门全绿。
⇒ 并入 `R-131-1` 作为形 b（与 X4 的形 a 同一根：清单来自形状而不是来自事实）。〔独立复现〕

### X7（对 AC#2 那句"`dir` 等于该 env 的数据根"定向拆台）—— **红，但红的不是本票的钉**

做法：`cmd/wisp/logsink.go:71` 的 `const logDirName = "logs"` ⇒ `"logz"`（落点目录改名，其余一字不动）。

```
--- PASS: TestAC4EveryLegIsNailedOrRuled
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk
--- FAIL: TestAC3LogSinkLandsInsideTheEnvDataRoot      ← 票 117 的那枚，按字面 "logs" 钉
```

⇒ 两枚新钉比的是 `logSinkDir(dataDir)`（生产函数本身），**实现搬走、判据跟着搬**，所以它们对这一发是哑的；
真正拦住它的是 `logsink_test.go:51` 那枚**故意拼字面量**的旧钉（它的注释逐字写着
"pinning it against `logDirName` would let a rename of the constant move the goalpost"）。
**判：这不是洞**（有独立一把尺、且理由写在码里），但**归因必须写清**：AC#2 那句"数据根"级判据的分母
不在本票两枚新文件里，动 `logDirName` 时红的是票 117 的用例。登记为口径条 `R-131-3`，不算缺陷。〔独立复现〕

### X11 / X12（候选 ②的第三种形状：那条 install 改经一枚函数值走）—— **X11 红在假话上、X12 全绿**

做法（都只有 1–3 行）：加 `var <别名> = installLogSink`，把调用点从直呼改成走别名。
名法游走（`pkg.funcs[名字]`）**吃不到函数值/变量别名** ⇒ 那条边在账本上消失。

**X11**（用在**已有钉**的 models 腿上）：`models.go` 加 `var sinkViaVar131 = installLogSink`，
`installLogSink(io_.dataDir)` ⇒ `sinkViaVar131(io_.dataDir)`。听众**照装不误**（那枚钉这次是 PASS 的），

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.03s)
    leg_sink_gate_131_test.go:214: AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "models",
        which no longer installs the persistent sink (main.go:74): the nail is now proving something else, or nothing
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)     ← 它当然会 PASS：听众还在装
```
⇒ 这次红了，但红的原因是"有钉 claim 它"、而且**红名说的是假话**（"no longer installs the persistent sink"——
它 install 得好好的）。账本那一行的原文：

```
leg models  main.go:74  install=false records=false ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk  -> no records
```
**门的判据是语法可达性，不是"有没有人装听众"**，这一点它文件头没说（文件头只说了"忽略 build tag ⇒ 只会多要钉不会少要钉"，
没说"直呼之外的边会静默消失"）。

**X12**（同一形，用在**没有钉的新腿**上；这才是真正的一发）：新增 `cmd/wisp/fake131.go`
（`var sinkAlias131 = installLogSink` + `sinkAlias131(root)` + `defer sink.close()`），
并在 `switch` 里加 `case "fake131": os.Exit(cmdFake131(args[1:]))`——**与 X6 逐字同形，只多一枚别名**。

- `go build` rc=0、`go vet` rc=0；`go test -count=1 -v ./cmd/wisp/`（**整包，不是只跑门**）⇒ **rc=0，
  `RUN 82 / PASS 82 / FAIL 0 / SKIP 0`，`ok` 47.283s（带噪）**；
- 账本共 **16 行**（基线 15 行 + 这条新腿那一行），而对这条腿的判定材料是假的：

```
leg fake131  main.go:82  install=false records=false ruled=false nails=-  -> no records
```

- **这条腿在真二进制上是活的**（我自己编了自己跑的，和 X4a 同一记仪器）：

```
$ /tmp/wisp131-x12.exe fake131 C:\...\x12data
time=... level=INFO msg="wisp: persistent log sink installed" dir=C:\...\x12data\logs min_level=info
fake131 ok                                                                        # rc=0
$ cat C:\...\x12data\logs\wisp-20260923-001.jsonl
{"...","msg":"wisp: persistent log sink installed","dir":"C:\\...\\x12data\\logs","min_level":"info"}
```

⇒ **门把一条真的装了持久听众、真写得出 jsonl 的 CLI 腿判成"无义务"**，全绿。
比 X4 更省：X4 要改分发形状，这一发只要**给函数起个别名**——`grep 'installLogSink'` 的人看不见它、
枚举门也看不见它，两侧同时瞎。⇒ 与 X4/X8 一起并入 `R-131-1`（形 a/b/c），修法③就是冲它写的。〔独立复现〕

### 我这轮的还原账（逐发还原 + 整树反查）

X1…X12 每发都在**独立副本目录**（`/tmp/wisp131-acc131r2-x*`）里就地改，基线快照本身没动过；
末了在基线快照上重跑 §1 那记 `diff -r --exclude=third_party` 对 `git archive 56d8026` 的反查 ⇒
**源码逐字节相同**（多出的一枚 `wisp.exe` 是 `go build` 的产物，已 `rm`）。
仓库工作树：`git status --short` 与 `git diff --numstat` 见 §11 原文，生产件零改动。

⚠ 还有一发 **X10**（注册一枚不读盘的用例当钉 ⇒ 门红）我把它记在 **§10 开头**，
因为它同时是实现方 **M6** 的复现，放在证据台账里好对照；它也是一发独立变异，读数同档。

---

## 5. AC 逐格裁决

| 格 | 判 | 一句话理由 + 证据档 |
| --- | --- | --- |
| **AC#1** 腿清单可 grep、给文件:行 | **通过**〔独立复现〕 | 清单不是手填的：基线 `-v` 里门现算出的 15 行与票面 11 组**行号逐条一致**；install 全集（定义 1 + 调用者 4）我另用 `grep` 独立数过，同数。 |
| **AC#2** models 腿补钉，进程内驱动、不 `t.Skip`、查内容不只查存在 | **通过**〔独立复现〕 | X1 我复算＝3 枚红点名 models；基线钉的 `-v` 原文显示记录 0 是 install 且 `dir` **等于** `logSinkDir(数据根)`、其后恰 1 枚 `models: hand-off REFUSED ... state=Error`、exit=1；**X9** 再把 install 挪到 `<数据根>-elsewhere` ⇒ 这枚钉立刻红（⇒ 它不是"有个 jsonl 就行"）；`grep 't\.Skip('` 在两枚新文件里 **0 命中**（唯一含 "Skip" 的是 `leg_sink_nail_131_test.go:52` 那句"NO t.Skip"注释——⚠ 仪器坑照旧：`grep -c 't.Skip'` 会把注释算成命中，这次直接用带左括号的形）。⚠ 一格口径要说清（不算缺陷）：钉注入的是 `modelsIO.dataDir`，所以它锁的是"给定数据根 ⇒ 落在 `<根>\logs`"，**不锁**"生产解析出的数据根"（后者由票 117 的 `TestAC3LogSinkLandsInsideTheEnvDataRoot` 守着，见 `R-131-3`）。 |
| **AC#3** `wisp secret` 先裁该不该装、再补钉 | **通过**〔独立复现〕 | 裁决＝装，理由在 `secret.go:196-218`（install 点旁，不是"以后再说"）；X2 我复算＝门同时给出"要求裁决"与"要求钉"两条读数，钉与降级子用例各红一次；**X9** 同形下 secret 那枚钉也红（落点搬走就红）。C28 那半（明文不落盘）由钉里的 `scanForPlaintext` 指向新 jsonl，我没有独立植一枚真形状密钥去撞（记 §10 未复现②）。 |
| **AC#4** 收敛判据：清单做成用例、**新增一条腿不给钉 ⇒ 当场红** | **退回**〔独立复现〕 | 前半（从源码枚举、非硬编码、不是恒绿）成立：X1/X2 两发点名红、X6 复现"新腿不给钉当场红"、X10 复现"注册一枚不读盘的用例 ⇒ 红"、X3c 证明整枚 switch 换形时它自报仪器坏了。后半被**我造出来三次**：X4 那枚 `--diag` 腿（早退 `if` 分发）在真二进制上装了听众、写得出 jsonl，门**不红也不报少一条**，再拆掉它的 install 仍 82/82 全绿；X8 把 `case "slo":` 的标签写成命名常量 ⇒ 那条腿**静默出账**、门 `PASS`；X12 与 X6 逐字同形、只多一枚 `var 别名 = installLogSink` ⇒ 账本把这条活腿判成"无义务"、门 `PASS`。三发同一根：**清单来自形状，不是来自事实**。按派单硬线（"AC 声称要防的结局如果被造出来 ⇒ 判退回，不许附条件通过"），这一格判退回，登记 `R-131-1`。 |
| **AC#5** 门禁 | **通过**〔独立复现〕 | §2 那张表：`-count=2` 164/164/0/0、CI 形状 82/46/0/0、台账八 scope 202/22/40/18/16/40/387/**34** 与两枚 rc=0 的格式门全部逐字复现；**控制组 `db9fafc` 的八数我也自己取了一枚**（`…/32` ⇒ 七数一字不降、`cmd/` +2 正好是那两枚新文件）；`GOOS=linux go vet` rc=1 且原文是派单预告的 sherpa 那句（只编译不执行）。墙钟我这边普遍慢 10–20%，**只登记为带噪**，不用它下任何结论。 |

---

## 6. `R-131-*` 清单

| 号 | 现象 | 能不能重跑 | 严重度 | 建议归属 |
| --- | --- | --- | --- | --- |
| **R-131-1** | **枚举门的"腿清单"只认一种形状，形状之外的腿免检**——`main()` 里 `switch args[0]` 的**字面量** case 加 `if len(args)==0`，且 install 的边必须是**被直呼的函数名**。三个独立复现：<br>**形 a（X4）**：一条经早退 `if` 分发的 CLI 腿（`--diag`）**根本不进账本**，门既不红也不报"少了一条腿"，整包 82/82 绿；把它的 install 再拆掉仍 82/82 绿。真二进制上它 `rc=0` 且落得出 `wisp-<day>-<seq>.jsonl` ⇒ 不是我种的死代码。<br>**形 b（X8）**：`case "slo":` 的标签写成命名常量 ⇒ 那条腿**静默出账**（账本里 `leg slo` 零命中、换成两行都叫 `default`），门 `PASS`；落空的正是"记了日志却没装听众 ⇒ 必须有裁决"那一行，也就是 `R-117-1` 的守卫。<br>**形 c（X11 有钉侧 / X12 无钉侧）**：`installLogSink` 的调用改经一枚函数值别名（`var 别名 = installLogSink`）。X11 在**已有钉**的 models 腿上做 ⇒ 账本报 `install=false`、门确实红，**但红名是假话**（"which no longer installs the persistent sink"，而它装得好好的、那枚钉 PASS）；X12 在**没有钉的新腿**上做（与 X6 逐字同形、只多一行别名）⇒ **门 `PASS`、账本写 `install=false records=false -> no records`**，而 `wisp131-x12.exe fake131 <根>` 真二进制 rc=0 并落得出 `wisp-<day>-<seq>.jsonl`。<br>⇒ 票 131 的立项理由（"装听众、拆掉、没人知道"）在这条门外**第九次逐字复现**，且 X12 的绕过成本只有**一行**。 | **能**：补丁都是分钟级——形 a＝`main.go` +3 行、新增 `cmd/wisp/diag131.go` 一枚函数；形 b＝`main.go` 换 1 个标签 + 1 行 const；形 c＝新增 1 行 `var 别名 = installLogSink` + 改 1 处调用名。快照目录 `/tmp/wisp131-acc131r2-x4-ifleg`、`-x8-constlabel`、`-x11-funcvalue`、`-x12-aliasnewleg`；读数在 §4 X4/X8/X11/X12 | **高**（本票唯一硬判据那一格被造出来了；今天无产品后果，因为这三形都是我种的，但它正是"下一次同族"的入口形状，而且**修法与病同量级：一行**） | **票 131 续单**，三行量级：① 断言"`main()` 里除已归类分支之外不再出现 `os.Args` 的读点"（读得到就红、点名那一行；我量过今天只有 `main.go:50` 一处命中 ⇒ 零误伤）；② `caseLabels131` 吃不下某个标签（`Ident`/常量）时**必须红**，而不是当它 `default`、或与另一行 `default` 撞键后静默；③ 走到 `pkg.funcs` 查不到的调用名时，若它是**函数值 var**（这一形可以在 `loadMainPackage131` 里顺手把包级 `var x = <funcName>` 解析成边）⇒ 建边或红着说"我看不见这条边"，**不许当叶子略过**。⚠ 另附一句给读的人：门的红名是"可达性"陈述不是"有没有人装听众"的事实陈述，X11 那一发它说了假话。⚠ **与 `R-121-1`／票 133 不是同一件事**：133 要的是"腿在 `switch` 里、但没人分发到 handler"（**下游**那一跳的可达性），我这条是"**腿清单本身**来自形状而不是来自事实"（**上游**少了一条）。两者同族、互补，且 **133 的 AC#1 硬判据（"必须能红 N-3 那一形"）今天并不覆盖 X4/X8/X12 这三形** ⇒ 建议把这三发作为**票 133 的第二～第四发判据**并进去，别再攒第三张同族票。 |
| **R-131-2** | **钉的"归属"列不可核**：门只核"有没有一枚读过盘的钉挂在这个腿名上"，不核"那枚钉驱的是不是这条腿"。X5 把 models 与 secret 两枚钉的腿名对调（用例本体、install 点、门全没动）⇒ `go build`/`go vet` rc=0、`go test -count=1 -v ./cmd/wisp/` **82/82 全绿**，账本两行仍各报 `nailed`，只是 `leg models` 的钉子写着 `TestAC3SecretLeg...`。 | **能**：一处 2 行的字符串对调；快照 `/tmp/wisp131-acc131r2-x5-wrongleg`，读数 §4 X5 | **中低**（损害边界我量清了：拆 install **仍会红**——`staleRegistrations131` 看的是腿还装不装听众，与钉驱谁无关；丢的是"账面 `nailed`、实际那条腿没人真测"这一格，也就是语义回归无人守，只保住"install 存在性"） | **票 131 续单**（与 R-131-1 同一枚文件、同批做最省）：最低形式是让注册带上"被驱入口的符号名"，门核一次"那枚钉的调用闭包里含不含这条腿的入口函数"（信息它已经全有：`pkg.tests` 里有每枚钉的 calls 集）。做不到就把这句话写进门的文件头，别让它当"每条腿都有人测"的证据。 |
| **R-131-3** | **口径归因，不是缺陷**：AC#2 那句"`dir` 等于该 env 的数据根"在 models/secret 两枚钉里是**相对**比较（与生产函数 `logSinkDir(dataDir)` 比）。X7 把 `logsink.go:71` 的 `const logDirName = "logs"` 改成 `"logz"` ⇒ 门与两枚新钉**全绿**，红的只有票 117 的 `TestAC3LogSinkLandsInsideTheEnvDataRoot`（它按字面量 `"logs"` 钉，注释逐字写了为什么不用常量）。 | **能**：一行常量改名；快照 `-x7-landingdir`，读数 §4 X7 | 低（**这一形今天被拦着**，拦它的尺是票 117 那枚；不改判 AC#2） | **无人为动作**，只登记一句给下一位：动落点名字时红的是 `logsink_test.go`，不要以为本票那两枚钉会响；`docs/reports/**` 由编排者记。 |
| **R-131-4** | **裁决注释里有一句描述与仪器输出不符**：`cmd/wisp/slo_windows.go:179` 的 `WISP-LEG-SINK-RULING:` 段写着"enumeration gate ... reports this row as `pipeline` rather than `install`"，而门的账本列是 `install=... records=... ruled=...` + 状态词 `ruled`，**基线整份 `-v` 输出里 `pipeline` 出现 0 次**（账本原文见 §2）。 | 能：§2 那张基线账本 + `grep -c pipeline` | 极低（措辞；但它是"注释描述仪器"的那一类，下一个人会照着找不存在的列） | **票 131 顺手**：把那半句改成实际列名（`ruled=true`）。别人不必动。 |

---

## 7. X9 ＋ 两枚后续 commit 的单判

### X9（"装对了听众、装错了树"——AC#2/AC#3 那句判据的真正分量）

做法：`models.go:276` 与 `secret.go:219` 的两处 `installLogSink(<数据根>)` 各加一个 `-elsewhere` 后缀
（`diff` 各只有 1 行；听众仍装、仍写盘、只是不写在这条腿的数据根里）。

```
--- PASS: TestAC4EveryLegIsNailedOrRuled (0.03s)          ← 门不看落点，这是分工
--- FAIL: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
--- FAIL: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.02s)
```
红名正文里带着被驱腿自己 stderr 的那一行原文（`... dir=C:\...\TestAC2ModelsLeg...551944006\001-elsewhere\logs`），
即读数是"它装到了哪"而不是"我猜它没装"。
⇒ **两枚新钉不是"查文件存在"那种假钉**：换树就红。这与票 127 的 `dir` 那一发同强度，我自己复算过。〔独立复现〕
⚠ 同时它也把门的分工说清了：**门管"这条腿 reach 得到 install 吗"，钉管"落点对不对、记录说了什么"**——
所以 R-131-1 那条（门看不见形状外的腿）**不能**拿"钉会红"来抵：X4 第二拍我连钉都没给那条腿，全绿。

### `627aa66` 单判（把"目录不存在"也读成"这条腿没有听众"）—— **通过，方向是对的**

`git show --stat 627aa66` ⇒ 只动 `leg_sink_nail_131_test.go`（+15/-4）。我在 `56d8026` 上读实现并复算：
`readLegSink131` 现在把 (a) `CountLogFiles` 报 `fs.ErrNotExist`、(b) 目录在而管道文件为零 合并成同一句主张
（`no wisp-<day>-<seq>.jsonl in <dir>: the leg ran to completion and its listener wrote nothing there.` + 被驱腿的 stdout/stderr），
另留一条 `err != nil && !errors.Is(err, fs.ErrNotExist)` 的响亮 `Fatalf`，和一条自称不可达的 `err == nil && n != 0`。
X1 的原文就是这句（不是实现方第一版那句 `counting pipeline files ... The system cannot find the file specified`）。
**判**：合并是对的——两种形状下"这条腿没有听众"都成立；被抹掉的只是 syscall 级细节，而它本来就在
`dir` 与整份 console 都被打进正文之后。**"路径拼错 ⇒ 红不红"这一问我的答案是 X9 + X7 两发**：
装错数据根 ⇒ 两枚钉红；落点目录本身改名 ⇒ 钉与门跟着搬、红的是票 117 那枚字面钉（⇒ `R-131-3` 归因条）。

### `56d8026` 单判（"两枚新文件为什么不打 windows tag"写成决定）—— **通过，且它自证了这条决定的可破性**

`git show --stat 56d8026` ⇒ 只动同一枚测试文件（+12，纯注释）。它写的决定是：
`package main` 没有 Linux 腿（`GOOS=linux go vet ./cmd/wisp/` rc=1，我在 §2 复算到同一句 sherpa 原文），
所以未打 tag 的文件"要么两边都编译得过、要么什么都没编译"；真出现 Linux 腿那天，
一枚编译不过的未打标文件比一枚静默掉分母的文件更响。
**判**：与实测一致（我这一轮所有 `GOOS=linux go vet` 读数都是 rc=1，从未跑到过执行）。
它**没有**为这句话新增用例，我也**不要求**——这一格的失效形（给门打上 `windows` tag 后 linux 静默）
被 `registerLegNail131` 收 func 值这件事堵在编译期：门文件与钉文件都未打标，
把任一枚挪到 tag 之外 ⇒ 另一枚引用不到 ⇒ 编译不过，而不是绿着没人发现。这条我核过形状，未实跑（Linux 腿不存在，跑不了）。

---

## 8. 总判

**总判：票 131 —— 退回一格（AC#4），其余四格通过。** 按派单硬线执行：
AC#4 声称要防的结局（"新增一条 CLI 腿、装了听众、不给钉、没人知道"）**被我真实造出来三形**
（X4 早退分支形、X8 命名常量标签形、X12 一行别名形），三形都伴随"整包 82/82 全绿或门 `PASS`，
且账本不报少一条腿"⇒ 不许写"附条件通过"。

一句话理由：**这枚门对它看得见的那 15 行是真判据（X1/X2/X6/X9 我都复算为红，且红名点到腿），
但它把"腿"定义成"`switch args[0]` 的字面量 case"，于是清单本身成了一条无人守护的假设——
这一族第九次不是"没人给钉"，是"没人核清单全不全"。**

### 8.1 它相对票 121/127 确实前进的地方（公道要替实现方说，逐条我复算过）

| 主张 | 我这轮的读数 |
| --- | --- |
| 拆 `models` 那 9 行今天会红、且点名 models | X1 ⇒ 3 顶层 + 1 子用例红，`R-127-6` 正面结清 |
| 拆 `secret` 那 5 行 ⇒ 门同时要求裁决与钉 | X2 ⇒ 两条独立读数（`leg_sink_gate_131_test.go:194` + `:214`） |
| 新增一条 switch 腿不给钉 ⇒ 当场红 | X6 ⇒ 唯一一枚红在门上，带行号与修法（＝实现方 M3，我复现） |
| 整枚 switch 换形 ⇒ 不交空表、自报仪器坏了 | X3c ⇒ `AC#4 RED (the instrument, not the code)` |
| `R-121-1`/`M-B1` 那一形（留 import、断分发） | X3 ⇒ **红**，且只红门（钉自己仍绿）——实现方自认的盲区被它自己的门补上了 |
| 钉不是"查文件存在" | X9 ⇒ 装到隔壁数据根就两枚齐红 |
| 注册一枚不读盘的用例当钉 ⇒ 红（防恒真钉那一半不是注释） | X10 ⇒ `AC#4 RED: nail ... calls none of the shared sink readers`（＝实现方 M6，我复现） |

### 8.2 能不能挡住第九次 —— **挡不住，而且绕它只要一行**

派单最后那一问的答案：它挡住了"**第八起逐字重演**"（拆掉某条在册腿的 install：X1/X2 都红），
挡不住"**换一条不在册的腿**"。三发独立绕过，成本从 12 行到 **1 行**：
X4（早退 `if` 分发，12 行）⇒ 全绿；X8（`case` 标签写成命名常量，2 行）⇒ 门 `PASS`、那条腿静默出账；
**X12（在 `switch` 里、与 X6 逐字同形，只把 `installLogSink` 换成一枚 `var 别名` —— 1 行）⇒ 门 `PASS`、
账本写 `no records`，而真二进制上这条腿确实装了听众并落出了 jsonl**。
⇒ 最省的那一发意味着：**这不是"有人恶意绕过"，而是"下一次顺手重构就会踩到"**——
`var newSink = installLogSink` 这种写法在重构里太常见了（换实现、按平台分派、注入替身）。
修法我写在 `R-131-1` 里，三条各一行量级（未归类的 `os.Args` 读点 ⇒ 红；标签吃不下 ⇒ 红；
调用名解析不到但是函数值 var ⇒ 建边或红），**不需要**新架构、不需要动别人地界
⇒ 所以我判"退回续单"而不是"整票推翻"：AC#1/#2/#3/#5 的交付是真的、可重跑的；只有 AC#4 那条"收敛"主张没到底。

### 8.3 票 131 能不能结案 —— **不能，按本表续一格**

- 阻断项有**一条、三形**：`R-131-1`（AC#4 那半格）。它需要**三发**新变异自证，缺一不可：
  ① 把 X4 那枚 `--diag` 腿种进**已补判据**的门 ⇒ 红、红名点到 `--diag`；
  ② X8 那枚命名常量标签 ⇒ 红（今天 `PASS`）；
  ③ X12 那枚"新腿 + 一枚 `var 别名 = installLogSink`" ⇒ 红（今天 `PASS`，且账本写着 `no records`）。
  三发各自还原 ⇒ 复绿。⇒ 判"退回"就是判"这三形的读数要从绿翻红再翻回绿"，不是重写这枚门。
- `R-131-2`（归属列）建议同批做（同一枚文件）；`R-131-4`（一句措辞）顺手。
- 结案判我写死成可复算的形式：**X1/X2/X3/X6/X9/X10 六发保持红名不变 ＋ X4/X8/X12 三发从绿变红 ＋ §2 那张表四数不降**，
  缺一条就不能翻 `-done`。别人的框我不翻，`Status` 行交回编排者。

---

## 9. 与同族账目的接口（归谁，逐条点名）

| 本表的号 | 与哪条旧账同一根 | 是不是同一件事 |
| --- | --- | --- |
| `R-131-1` | `R-121-1` / 票 133 AC#1、`A105④`、`A110④`（"能力装好没人调"第九次半例） | **不是同一件事，是同一族的上下游**。133 核"在册的腿还有没有人分发到 handler"；我这条核"清单收没收到这条腿、收进来的那行是不是真话"。**133 现在的唯一硬判据（"必须能红 N-3 那一形"）不覆盖 X4／X8／X12 这三形** ⇒ 建议把这三发直接写进 133 当第 2～4 发判据；若编排者愿意让 131 续单自修那三行，也可以，但**别让它变成第三张同族票**。 |
| `R-131-2` | 票 121 `R-121-2`（"策略清单加长被查、缩短静默"） | **同一种病的另一枚尺**：都是"账本的某一列由实现自己填，没人核"。修法也同为"几行断言"级。 |
| `R-131-3` | `A110③` 那条"允许改 127 的钉，但判据是意图保留、语义变强" | 无冲突，只是归因：落点这一形由票 117 的字面钉守着，不是 131 的两枚新钉。 |
| `R-131-4` | —— | 纯措辞，与任何旧账无关。 |

---

## 10. X10（实现方 M6 的独立复现）＋ 三档证据台账 ＋ 我没能复现的自述

### X10 —— 注册一枚"读了别处、不读盘"的用例当钉：**红**

做法：只改注册那一行的第二个实参，`registerLegNail131("models", TestAC2ModelsLeg...)` ⇒
`registerLegNail131("models", TestSecretFlagsAreBoolOnly)`（`secret_test.go:372` 的一枚真实用例，不碰任何读盘仪器）。

```
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.03s)
    leg_sink_gate_131_test.go:188: AC#4 RED: nail "TestSecretFlagsAreBoolOnly" for leg "models" calls none of the
        shared sink readers (readLegSink131, readResidentSink, readSink), so it cannot be the case that goes red
        when the install block is deleted (leg site main.go:74).
```
⇒ 实现方"防恒真钉"那一半**不是注释，是仪器**，我复算为真。〔独立复现〕
⚠ 但请注意它与 X5 的合起来读法：**它核"这枚钉读不读盘"，不核"它读的是不是这条腿的盘"**（X5 已证）。
⇒ 门的"非恒真"保证是**包级**的，不是**逐腿**的。这句话写进 `R-131-2` 的修法里。

### 三档证据台账

- **〔独立复现〕**（我自己在 `56d8026` 纯净快照上跑出来的）：§1 两记完整性反查、§2 基线全部读数
  （四数 ×2 种仪器、台账八 scope、gofmt/gofumpt/vet/`GOOS=linux vet`、门的 15 行账本原文）、
  候选①～④ 共 **13 发编号变异**（X4 分两拍 ⇒ 14 次读数）：X1(=M1)、X2(=M5)、X3(=M-B1/`R-121-1` N-3 形)、X3c、X4 两拍、X5、X6(=M3)、X7、X8、X9、X10(=M6)、X11、X12、
  真二进制 `--diag` 落得出 jsonl 这一条、install 全集与 `os.Args` 全集的 grep 计数、
  `R-131-1..4` 每条的现象列。
- **〔日志＋归档，我抽验〕**：**无**。本轮我没有把任何别人日志里的读数当结论用（CI 我只在 §1 用它判"本机忙不忙"）。
- **〔仅自述，我没能复现，不背书〕**：
  ① 实现方 **M4 的字面形**（种一条假腿**再**给它一枚钉 ⇒ 绿）。等价主张我承认（基线 15 行里 10 行判"no records"仍 PASS ⇒ 门会闭嘴），
     但我没照它那样种第二枚假腿再配钉 ⇒ 这一发的原文我不背书。
  ② `secret` 钉里"`scanForPlaintext` 也指向新 jsonl、假 key 不在任何落盘字节里"——我读了码与它的调用点，
     **没有**独立植一枚真形状的密钥去撞那条断言。
  ③ 票面"顺手动到的一处"：`resident_sink_nail_127_windows_test.go` 那枚 ordering 用例补的 `len(recs)==0` guard。
     我核到它确实在 `7e60d31` 的 `--name-only` 里（该文件 +13 行）、基线里这枚用例 PASS，
     **没有**复现"编队负载下子进程以 `0xc000013a` 早死 ⇒ 空文件 ⇒ panic 吃掉后面的红名"那一发（那需要再造负载）。
  ④ 墙钟账（"本票新增 0.17s/两趟、单趟 +0.20%"）：我没逐条重量；本机这一轮同形包总时长在 41.9–96.0s 之间漂，
     噪声远大于该量级 ⇒ 按 §1 的判据我只登记为带噪，不引用来支持或否证。
  ⑤ CI 侧分母那格（`next=` 第 ③ 条）："枚举门只在 windows-leg 有分母"我**没**去开 CI run 日志核 run id+step
     ⇒ 按本仓的规矩，说不出"上次真跑过的 run id + step"就当那一步不存在；我也因此**不**把"门在 CI 上有分母"记成已证。
  ⑥ `git show --stat` 四枚 commit 的文件清单我逐枚点了（与票面 `--name-only` 一致，零越界），
     但"票 123 那四枚 CLI 红不在本包"我只在**本包四数 FAIL=0** 这一层复算，没去重跑那四枚红所在的包。

---

## 11. 地界自证（原文）

先替被验方核一遍**它**的地界（`git diff --name-status db9fafc..56d8026` 原文）：

```
M  .scratch/wisp/issues/131-the-models-leg-has-no-nail-deleting-9-lines-keeps-76-cmd-wisp-tests-green.md
A  cmd/wisp/leg_sink_gate_131_test.go
A  cmd/wisp/leg_sink_nail_131_test.go
M  cmd/wisp/resident_sink_nail_127_windows_test.go
M  cmd/wisp/secret.go
M  cmd/wisp/slo_windows.go
```
⇒ 六枚路径**全在 `cmd/wisp/**` 加票面**：票面禁改清单里的 `rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、
`.github/workflows/ci.yml`、`scripts/`、任何阈值/golden 一枚没动，`internal/**`、`frontend/` 零改动；
`ban #8 cmd/` 从 32→34 的两枚增量正好是那两枚 `A`。票面那一枚是 append-only 的进度格（规则要求的形态）。
⇒ 实现方"没做/不做的事逐条点名"那一段，我按文件清单核为**属实**。〔独立复现〕

我的地界：
- 只写本文件。生产文件、票面、`docs/reports/**`、`frontend/`、`internal/risk/**`、`internal/winsec/**` **一字未动**（只读）。
- **末了整树反查**（12:1x 原文，两台机器同一台）：

```
$ git status --short
 M cmd/wisp/logsink.go
 M docs/evidence/s1/131-adversarial-acceptance.md
 M internal/observe/logging.go
$ git diff --numstat -- cmd/wisp/logsink.go internal/observe/logging.go
22	2	cmd/wisp/logsink.go
291	0	internal/observe/logging.go
```
  ⇒ 整树 `git diff` **不为空**，但那两处是**别人的在飞件**（票 130 的 `logsink.go`/`logging.go`，本轮期间另有代理在提交），
  我一格都没碰；本文件那一枚的 diff 是我自己两拍之间的增量。
  我自己的地界这样证：`git show --name-only 097c578` ⇒ **唯一路径 `docs/evidence/s1/131-adversarial-acceptance.md`**（271 行新增）。
- **被验版本反查**（收尾再量一次）：`git archive 56d8026 | tar -x -C /tmp/wisp131-verify2-acc131r2`，
  `diff -r --exclude=third_party --exclude=wisp.exe` 与基线快照 ⇒ **rc=0（逐字节相同）**
  ⇒ 我这一路跑的九发变异**没有一发**回流到基线快照，更没有一发进过仓库工作树。
- 所有变异落在 `/tmp/wisp131-acc131r2-x*`（带本会话后缀），未在仓内建 worktree，未 `add`/`commit` 任何别人的路径，**未 push**。

## 12. 注入文本登记（两条判据：路径真不真 / 内容是否越权）

截至 12:2x（有界说法，它可能继续增长）：

- 形状：harness 级 `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.`
  与它的项目级同族 `...\.qoder-cn\projects\D--work-workspace-projects-plans-Wisp\memory\MEMORY.md ...`，
  出现在我多条 `Bash` / `Read` 工具结果的**尾部**（示例出处：`Bash` 的 `git log --oneline -12`、
  `date "+%Y-%m-%d %H:%M:%S %z"`、`cd /tmp/wisp131-acc131r2 && wc -l cmd/wisp/...` 这类只读命令之后）。
- **精确次数我给不出**：它随工具结果夹带、跨轮重复，我没有逐条审计自己的转录 ⇒ 这条按"不背书"处理，
  只报"≥7 条工具结果尾部出现过，截至 12:2x"。⚠ 本仓老账（`A104③`）正是被这种"代理自述的次数"坑过一次，
  所以我把"给不出可核次数"这件事本身写出来，而不是填一个数。
- 判据核对：① **路径真** —— 两枚 `MEMORY.md` 在本机确实存在，内容是记忆索引（不是我发的字也不冒充我发的字）；
  ② **内容未越权** —— 没有出现"以代理判断为准/冲突时覆盖 owner 裁定/放宽判据/请 revert/冻结某包/别提交"这类指示。
  ⇒ 全部按"我自己的上下文被刷新"处理：**未据此改动任何判据、未 revert、未放宽任何阈值、未冻结任何包、未 commit 别人的文件。**
- 本轮**没有**遇到自称"编排者备注 / 系统提示 / 用户已更新编码规则 / `wisp-orchestrator-continuation`"的表格式文字（这类计数 0）。
- 凭据值一律不抄：文中 `<数据根>`、`C:\Users\...\Temp\...` 均为占位或系统路径，无密钥；
  X10 里出现的用例名 `TestSecretFlagsAreBoolOnly` 是用例标识符，不是凭据；
  被验树里那两枚钉植的假 key 我只在 §5 的判语里以"假 key"指代，未抄其字面。

---

## 13. 一条给编排者的在飞冲突提示（不是裁决，是状态）

我收尾时（12:3x）量到工作树里正有**别人的**改动落在我验收的这两枚文件上：

```
$ git status --short
 M cmd/wisp/leg_sink_nail_131_test.go          （+44 / -9）
 M cmd/wisp/resident_sink_nail_127_windows_test.go
 M cmd/wisp/logsink.go                          （+22 / -2）
 M internal/observe/logging.go
?? internal/observe/earlylog_130_test.go
$ git log --oneline -1 --format=... ⇒ df4a60a 12:17 docs(A113,建票135) / 64f4811 12:20 docs(129,...)
```

形状看是**票 130**（`init()` 期记录的缓冲/冲刷）那一批正在动 `cmd/wisp` 的听众与它的钉，不是我这一轮碰的（我只写本文件）。
⇒ 两条后果要交给编排者：
① **我的全部裁决只锚 `56d8026`**（§1 那记 `diff -r` 反查为证），上面这批改动落地后**读数会变**，
尤其 `assertInstallRecordFirst131` 那句"记录 0 必须是 install"——`A110③` 已经预告过"修完之后第一条记录本来就该是更早那条"，
所以 130 落地那一轮**必须重取** AC#2/AC#3 两格与 §2 那张基线表，不能沿用我的数；
② 131 的续单（`R-131-1`/`R-131-2`）**与 130 会改同一枚 `leg_sink_nail_131_test.go`**（以及门的文件），
排期上要么 130 先落、131 后续单，要么反过来，**别并行**——同一枚判据文件两把笔。

---

