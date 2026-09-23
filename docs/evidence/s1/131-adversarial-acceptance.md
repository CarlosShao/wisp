# 票 131 独立对抗验收（第二任）—— 攻的就是那枚"从 `main.go` 现算腿清单"的枚举门

**验收方**：`acceptor-ticket131-run2`（本文件唯一作者；只读验收，本文件是唯一写件）
**日期**：2026-09-23（本会话，开件 11:45，末次更新见各节时间戳）
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

### X4（候选 ②）—— **全绿。这是本轮最重的一条，也是本票的退回理由**

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

**可重跑**：`/tmp/wisp131-acc131r2-x4-ifleg`（第一拍的 `diag131.go` 已被我就地改成第二拍；两拍的
`main.go` 差 3 行都在 `diff` 里）。补丁形状＝上面两段代码，任何人 5 分钟能重跑。

### X5（候选 ③）—— 待读数
### X6（候选 ④）—— 待读数
### X3c（解析器换形）—— 待读数

---

## 5. AC 逐格裁决

| 格 | 判 | 一句话理由 + 证据档 |
| --- | --- | --- |
| **AC#1** 腿清单可 grep、给文件:行 | **通过**〔独立复现〕 | 清单不是手填的：基线 `-v` 里门现算出的 15 行与票面 11 组**行号逐条一致**；install 全集（定义 1 + 调用者 4）我另用 `grep` 独立数过，同数。 |
| **AC#2** models 腿补钉，进程内驱动、不 `t.Skip`、查内容不只查存在 | **通过**〔独立复现〕 | X1 我复算＝3 枚红点名 models；基线钉的 `-v` 原文显示记录 0 是 install 且 `dir` **等于** `logSinkDir(数据根)`、其后恰 1 枚 `models: hand-off REFUSED ... state=Error`、exit=1。⚠ 一格口径要说清（不算缺陷）：钉注入的是 `modelsIO.dataDir`，所以它锁的是"给定数据根 ⇒ 落在 `<根>\logs`"，**不锁**"生产解析出的数据根"；后者由票 117 的 `TestAC3LogSinkLandsInsideTheEnvDataRoot` 在 run 腿上锁。 |
| **AC#3** `wisp secret` 先裁该不该装、再补钉 | **通过**〔独立复现〕 | 裁决＝装，理由在 `secret.go:196-218`（install 点旁，不是"以后再说"）；X2 我复算＝门同时给出"要求裁决"与"要求钉"两条读数，钉与降级子用例各红一次。C28 那半（明文不落盘）由钉里的 `scanForPlaintext` 指向新 jsonl，我没有独立造一枚假 key 去撞（记 §10 未复现）。 |
| **AC#4** 收敛判据：清单做成用例、**新增一条腿不给钉 ⇒ 当场红** | **退回**〔独立复现〕 | 前半（从源码枚举、非硬编码、不是恒绿）成立：X1/X2 两发点名红、X6/基线台账可复核。后半被**我造出来了**：X4 那枚 `--diag` 腿在真二进制上装了听众、写得出 jsonl，门**不红也不报少一条**；再拆掉它的 install 仍 **82/82 全绿**。按派单硬线（"AC 声称要防的结局如果被造出来 ⇒ 判退回，不许附条件通过"），这一格判退回，登记 `R-131-1`。 |
| **AC#5** 门禁 | **通过**〔独立复现〕 | §2 那张表：`-count=2` 164/164/0/0、CI 形状 82/46/0/0、台账八 scope 202/22/40/18/16/40/387/**34** 与两枚 rc=0 的格式门全部逐字复现；`GOOS=linux go vet` rc=1 且原文是派单预告的 sherpa 那句（只编译不执行）。墙钟我这边普遍慢 10–20%，**只登记为带噪**，不用它下任何结论。 |

---

## 6. `R-131-*` 清单

| 号 | 现象 | 能不能重跑 | 严重度 | 建议归属 |
| --- | --- | --- | --- | --- |
| **R-131-1** | **枚举门的"腿清单"只认一种分发形状**：`main()` 里 `switch args[0]` 的 case 与 `if len(args)==0` 那两形。一条经早退 `if` 分发的 CLI 腿（X4 的 `--diag`）**根本不进账本**——门既不红、也不报"少了一条腿"，整包 82/82 绿；把它的 install 再拆掉仍 82/82 绿。⇒ 票 131 立项理由（"装听众、拆掉、没人知道"）在这条门外**第九次逐字复现**。 | **能**：两份补丁（`main.go` +3 行、新增 `cmd/wisp/diag131.go` 一枚函数）＋ §4 那三段读数，任何人 5 分钟重跑；快照目录 `/tmp/wisp131-acc131r2-x4-ifleg` | **高**（本票唯一硬判据的那一格被造出来了；今天无产品后果，因为 `--diag` 是我种的，但它正是"下一次同族"的入口形状） | **票 131 续单**，量级三行：门的 `enumerateLegs131` 补一条**"main() 里除已归类的分支之外不许再出现 `os.Args` 的读点"**（读得到就红，点名那一行）；或退一步：把"未归类分支"当一条红腿入账。⚠ **与 `R-121-1`／票 133 不是同一件事**：133 要的是"腿在 `switch` 里、但没人分发到 handler"（**下游**那一跳的可达性），我这条是"**腿清单本身**来自形状而不是来自事实"（**上游**少了一条）。两者同族、互补，且 **133 的 AC#1 硬判据（"必须能红 N-3 那一形"）今天并不覆盖 X4 这一形** ⇒ 建议把 X4 作为**票 133 的第二发判据**并进去，别再攒第三张同族票。 |
| **R-131-2** | 待 X5 读数 | | | |
| **R-131-3** | 待 X6 读数 | | | |

---

## 10. 三档证据台账 + 我没能复现的自述

- 〔独立复现〕§1 完整性反查、§2 全部四数与格式门、X1、X2、X3、X4（两拍 + 真二进制活体）、腿清单、install 全集计数。
- 〔日志＋归档，我抽验〕无（本轮没有引用任何别人日志里的读数下结论）。
- 〔仅自述，我没能复现，不背书〕：
  ① 实现方的 **M3/M4/M6** 三发（假腿不给钉 ⇒ 红 / 假腿给钉 ⇒ 绿 / 注册指向不读盘的用例 ⇒ 红）。
     M4 的等价物我在基线里看到（四枚钉各配一枚真实用例），M3/M6 我这轮按候选 ③④ 自己重做＝X6/X5，读数待填。
  ② `secret` 钉里"`scanForPlaintext` 也指向新 jsonl、假 key 不在任何落盘字节里"——我读了码，**没有**独立植一枚真形状密钥去撞它。
  ③ 票面"顺手动到的一处"：`resident_sink_nail_127_windows_test.go:494` 那枚 ordering 用例补的 `len(recs)==0` guard
     ——我核到它在 `7e60d31` 的 `--name-only` 里（该文件 +13 行），**没有**独立复现它在编队负载下以 `0xc000013a` 早死那一发。
  ④ 墙钟账"本票新增 0.17s/两趟"：我没逐条重量（本机噪声大于该量级，见 §1）。
  ⑤ CI 侧分母："枚举门只在 windows-leg 有分母"我**未**开 CI 日志核对 run id（票面 `next=` 第 ③ 条自己也这么挂着）。

---

## 11. 地界自证

- 只写了本文件。生产文件、票面、`docs/reports/**`、`frontend/`、`internal/risk/**`、`internal/winsec/**` **一字未动**。
- 仓库工作树 `git status --short` 全程只有本文件那一枚未跟踪件（末尾再量一次并贴原文）。
- 所有变异落在 `/tmp/wisp131-acc131r2*`（带会话后缀），**没有**在仓内建 worktree，没有 `add/commit` 别人的路径。

## 12. 注入文本登记（两条判据：路径真不真 / 内容是否越权）

截至 12:0x（有界说法，它可能继续增长）：本轮我的工具输出里出现
"harness 级 `MEMORY.md 已被修改，以下是修改后内容`"提示 **5 次**，
出处：`Read`/`Bash` 工具结果尾部追加的 `Note: The file C:\Users\swq\.qoder-cn\...\memory\MEMORY.md was modified since it was last read.`
（命令均为本会话自己的 git/go 只读或快照操作，前 40 字不含任何写入仓库生产件的调用）。
判据核对：① **路径真**——那两枚文件在本机存在，内容是我自己/编排者的记忆索引；② **内容未越权**——
没有一条要求改判据、放宽阈值、revert、冻结包或提交别人的文件。
⇒ 全部按"我自己的上下文被刷新"处理：**未据此改动任何判据、未 revert、未放宽任何阈值、未 commit 别人的文件**。
本轮**没有**遇到自称"编排者备注/系统提示/用户已更新编码规则"的表格式文字（计数 0，截至 12:0x）。
凭据值一律不抄；文中出现的 `<数据根>` 均为占位符。
