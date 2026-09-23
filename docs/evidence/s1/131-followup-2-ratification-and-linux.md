# 票 131 续单 第 2、3 件：第 0 件的 Linux 复核、票 130 越权改动的追认取证、`R-131-4` 措辞

会话：`worker-ticket131-followup-r2`。锚点：开工 `git rev-parse HEAD` = `b723978`
（跑到 14:3x 时共享树被兄弟代理推进到 `7b4c36a`、再到 `c66d683`；新落的那几枚没有碰 `cmd/wisp/**`，
所以本档所有读数仍对得上 `b723978` 那一版分母，见 `131-followup-1-three-shapes.md` §1）。

## 1. 复核第 0 件（`717d822` 把 131 的钉挪进 `_windows_test.go`）：Linux 侧两条读数

仪器：Docker 原生 `golang:1.27` / `linux/amd64` / `CGO_ENABLED=1`，仓内工作树挂到 `/src`
（`MSYS_NO_PATHCONV=1 docker run -v "/d/...:/src"`，进去先 `ls -l /src/go.mod` 证挂载不是空挂）：

```
-rwxrwxrwx 1 root root 855 Sep 21 01:17 /src/go.mod
linux
amd64
1
--- go vet ./cmd/wisp/ ---
VET_RC=0
```

⇒ **`undefined: sinkInstallRecord` 在 Linux 上不再出现**（第 0 件成立；那句曾在 CI run `35817761098` 里带病 4 小时）。
顺带：这一发容器里 `go vet` 之前先真下载了 `sherpa-onnx-go-linux v1.13.8` 等模块 ⇒ 不是"包没加载所以零检查"那种假 rc=0。

**第二问：那扇枚举门在 Linux 上是不是恒绿？不是。**`go test -count=1 -v -run TestAC4EveryLegIsNailedOrRuled ./cmd/wisp/` ⇒ **FAIL**，
零分母那一发如实报名 GOOS、如实点出四枚读不见的腿，并把"看不见"与"没得抱怨"分开写死：

```
    leg_sink_gate_131_test.go:397: AC#4 RED (the instrument, not the code): zero nails registered in this test binary, so the gate has nothing to reconcile and would pass on an empty list.
        GOOS=linux compiled no nail file: every registerLegNail131 call lives in a *_windows_test.go file, so the 4 leg(s) in the ledger that reach installLogSink (models, no-args, run, secret) are UNREAD on this platform, not unnailed - and no row of this ledger is coverage here, in either direction.
        Fix: run this package where the nails compile (scripts/wisp-cli-tests.sh, the windows leg), or give this GOOS its own nail file and register it. This reading stays red on purpose: a gate that cannot see is not a gate that has seen nothing to complain about.
```

同一发的账本（Linux 上清单仍从磁盘源码现算，15 行一字不少，只有 `nails=` 那列空）：

```
  leg models  main.go:74  install=true  records=true  ruled=false nails=-  -> RED listener installed, no nail
  leg no-args main.go:51  install=true  records=true  ruled=false nails=-  -> RED listener installed, no nail
  leg run     main.go:61  install=true  records=true  ruled=false nails=-  -> RED listener installed, no nail
  leg secret  main.go:72  install=true  records=true  ruled=false nails=-  -> RED listener installed, no nail
  leg slo     main.go:82  install=false records=true  ruled=true  nails=-  -> ruled
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.14s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.153s
```

⚠ 仪器瑕疵登记：容器里那句 `GATE_RC=${PIPESTATUS[0]}` 在 `sh`(dash) 下报 `Bad substitution`，
所以 `GATE_RC` 没打出来；红是 `--- FAIL` + `FAIL ... 0.153s` + `FAIL` 三行原文给的。

Windows 侧同码对照（`-count=1 -v`，见 followup-1 §10）：门 PASS、四行 `nailed` 逐字带钉名 ⇒
第 0 件那句"两平台拿到同一张腿表、只有 nails 列会空"成立，Linux 那一发的红**不是**新引入的假红，
它正是"零分母不许当绿"的兑现。

## 2. 票 130 越权改 131 的 `assertInstallRecordFirst131`：追认取证（只给读数，不下撤销令）

被追认的改动（`9b5d64d`，commit message 自己登记为"属超授权改动"、撤销口令「131 那枚恢复原样」）：
把"booking 必须是 record 0"换成"booking 之前只许是 130 回放的 boot 记录"。三条探针：

### P1＝删掉 130 的冲刷（130 自报的 M-1 形），helper 还红不红

造法（快照 `/tmp/wisp131-r2-p12`，`logsink.go:173`）：`earlyFlushed, earlyDropped := observe.FlushEarlyLogRecords(p.Handler())`
⇒ `var earlyFlushed, earlyDropped int`。`go build` rc=0、`go vet` rc=0。

- 只跑 131 那枚用 helper 的用例（`-run TestAC2ModelsLegBooksItsHandOffVerdictOnDisk`）⇒ **rc=0、`--- PASS`**。
- 整包（`-count=1 -v`）⇒ rc=**1**，`RUN 100 / 50 / 3 / 0`，三枚红逐字是：

```
--- FAIL: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.15s)     <- 票 130 自己的断言 1
--- FAIL: TestAC1ResidentLegInstallsItsLogListenerOnDisk (3.98s)          <- 票 127 的钉
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (3.12s)  <- 票 127 的钉
```

127 那两枚的红名里能看到它钉的是"回放必须在"：

```
    resident_sink_nail_127_windows_test.go:434: install booking is record 0, want 1: the replayed early records come first and nothing else may sit before the booking (records: [wisp: persistent log sink installed])
    resident_sink_nail_127_windows_test.go:587: record 0 = "wisp: persistent log sink installed", want "winsec: sealing path resolver installed": ... record 0 being the booking again means the replay of the boot-time sealing verdict is gone
```

⇒ **简报里"P1 ⇒ helper 仍须红"这一条断言不成立**（如实报，不硬改）：131 的 helper 在冲刷被删之后**不红**，
红的是 130 自己的断言 1 与 127 的两枚钉——与 `9b5d64d` commit message 自报的
"M-1 删冲刷红在断言1与127钉的两条"**逐字同账**。还原 ⇒ 整包 `100 / 53 / 0 / 0` rc=0。

### P2＝原始疾病：某腿的事件早于它自己的 install（要害那一发）

造法（`cmd/wisp/models.go` 加一枚包级 `init()`，用 `slog.Info` 记一条本腿自己的事件；
`log/slog` 进 import）：这样那条记录比任何 `installLogSink` 都早，会被 130 的回放冲在 booking **前面**。
`go build` rc=0。只跑 models 那枚钉（进程内第一枚 install，回放归它）⇒ rc=**1**，红名逐字是 helper 自己的那一发：

```
    leg_sink_nail_131_windows_test.go:366: record 1 sits before the listener's booking and is "models: P2 planted event booked before this leg had a listener", want "winsec: sealing path resolver installed" or nothing: a leg that logs something before installLogSink loses exactly the record it exists to keep (records: [INFO/winsec: sealing path resolver installed INFO/models: P2 planted event booked before this leg had a listener INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED ...])
```

整包同发 ⇒ rc=1（`100 / 50 / 3 / 0`：131 的 models 钉 + 127 的两枚钉）。还原 ⇒ `100 / 53 / 0 / 0`。

### P3＝把 helper 的第一条判据换回 130 之前的形状，打在**没有生病**的树上

造法：`for i, r := range s.records[:installIdx] { if r.Msg != residentEarlyResolverMsg {...} }`
⇒ `if installIdx != 0 { t.Errorf("install record is at index %d, want 0: ...") }`（与旧形状语义等价：
旧码是 `s.records[0].Msg != residentInstallMsg`，而 `installIdx` 就是第一条 booking 的下标）。
`go build` rc=0，树上没有任何疾病（130 的回放在位）。读数
（`-count=1 -v -run 'TestAC2ModelsLegBooksItsHandOffVerdictOnDisk|TestAC3SecretLegBooksItsAuditRecordsOnDisk|TestAC2AC3DegradedLegsStillDeliverTheirVerdict'`）：

```
RUN=5 TOPPASS=2 TOPFAIL=1 TOPSKIP=0
--- FAIL: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
    leg_sink_nail_131_windows_test.go:362: install record is at index 1, want 0: the pre-ticket-130 assertion, restored for P3
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.03s)
```

⇒ 旧的那条判据在**健康的树**上就红（谁先安装谁拿到回放 ⇒ 它红不红取决于测试顺序，不是取决于秩序），
这一发是"130 为什么必须动这枚 helper"的正面证据；同时它也说明旧判据红得不稳（同一枚 mutation 下 secret 那枚 PASS）。


### 结论（只写追认成立／不成立）

⚠ 仪器版本登记（不藏着）：P1/P2 的**整包四数**取自 13:5x 那一版三枚件（还没有 `rulingCrossChecks131`）。
这三发都不碰任何腿的 `install/records/ruled` 读数（标记词全仓只在 `slo_windows.go` 与门文件里，逐枚 grep 过），
而新判据在未变异树上零命中（`131-followup-1-three-shapes.md` §6 末与 §10 的 `AC#4 RED` 计数 0）
⇒ 判定不变；但"最终码下整包重跑 P1/P2"这一发**我没有单独复跑**，登记为待核，不算已验。
P3 那一发是最终码跑的（`-run` 只圈三枚钉，门不在分母里）。

**追认成立。** 依据：P2 那一发（原始疾病）在改后的 helper 下**照样红**、且红名点到"booking 之前坐了一条本腿的记录"；
helper 的第一条判据从"下标 0"换成"booking 之前只许是回放"之后，
它没有把 130 的回放当成新的必要条件（P1 下 131 的用例不红，红的是 130/127 自己那两枚），
也就是说这一改**没有把别人的修法偷渡进 131 的判据**。
登记一句给下一位：P1 那发同时说明 131 的 helper **看不见**"回放缺失"这一形——那一形的证人在
`early_log_nail_130_windows_test.go` 与 `resident_sink_nail_127_windows_test.go`，不在本票。
**撤销令不在本档**：130 那枚改动要不要恢复原样，只服从编排者那边的对话。

## 3. `R-131-4`（`cmd/wisp/slo_windows.go` 裁决段那句描述）

验收方 §6 的原判：那句写着门会把这行报成 `pipeline`，而门的账本列是 `install=... records=... ruled=...` + 状态词，
"基线整份 `-v` 里 `pipeline` 出现 0 次"。复核：

```
grep -c pipeline /tmp/wisp131-r2-base-c2.log   -> 0        # b723978 纯净快照, go test -count=2 -v ./cmd/wisp/
账本列实际形状                                    -> install=false records=false ruled=false
状态词实际集合                                    -> -> nailed / -> no records / -> ruled
```

落点（只改措辞、零语义改动）：`cmd/wisp/slo_windows.go:179-186` 那半段，旧文
"reports this row as \"pipeline\" rather than \"install\"" ⇒ 新文改成门真读的与被真报的：

```
// enumeration gate in leg_sink_gate_131_test.go reads this function's own doc
// comment for the sentence above and its body for the calls below, so the row it
// books for `wisp slo` in that gate's ledger is install=false with records=true
// and ruled=true, and the state word on that line is "ruled" rather than
// "nailed". This sentence is the reason that is a decision and not an
// omission.
```

与实测对齐：门的 slo 行原文 `leg slo main.go:82 install=false records=true ruled=true nails=- -> ruled`
（`-count=1 -v` 与 `-count=2 -v` 两趟逐字相同）。标记词位置也照实登记：`slo_windows.go:171`（验收方写 `:179`，
差 8 行是 130 在同一段里追加过句子所致，不是谁的读数错）。
`WISP-LEG-SINK-RULING:` 承不承重由 followup-1 §6 的 A0 那一发回答（删掉标记词 ⇒ 门红并点名 `slo`）。
