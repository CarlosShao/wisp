# 票 130 刷新腿 `130-r2` — 交件（写腿 / 实现者）

**两件**：①`AC#1` 的"哪些记录会掉"清单按 HEAD 重刷（**零生产码**）；②两处 `.go` 注释里那句"塞进 `InitLog` 会花掉票 132 的封印判据"换成**真代价**（**只改注释**）。
**锚定** `dbe861a02`（起手 HEAD）｜**写面** `internal/observe/logging.go`（`:486-495` 注释）＋`internal/observe/earlylog_takeover_130r1_test.go`（`:16-22` 注释）＋本件＋`130-r2-impl-part2.md`。
**零行为改动的尺**（§6 现跑）：`git diff -U0 -- internal/observe/logging.go internal/observe/earlylog_takeover_130r1_test.go | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vE '^[+-][[:space:]]*//'` ⇒ **0 行**。

## 0. 三句结论

1. **分母今天＝129 行 / 27 包**（件里 53/22；判语 129/27 ＝ **对**）。`internal/audio`／`internal/ball` **今天在闭包内**，主式子**闭包外 0 行** ⇒ 件里"闭包内 28 / 闭包外 25"那套分法今天不成立。
2. **"其余 24 枚"既不是 24、也不是 19：今天 118 枚**，逐枚在 `130-r2-impl-part2.md`（129 行全表、每行带 `grep -nF` 锚）。恒等式 **`129 = 5(A) + 4(B) + 1(C') + 1(F) + 118`**；另有 **1 枚 `Debug` 形（A6）在 129 之外**。
3. **`internal/winsec/winsec_windows.go:143` 归 B 窗口**（判语对、件漏了 B 这一半）：它在 `observe.InitLog` 内部就发（pipeline 自己封 `logs\` 那次）⇒ 判语 §4 的 post 发里它正是**记录 1**。另**新补一枚 C'**：second launch 腿**根本没有听众**（`internal/proc/singleinstance_windows.go:171`）。

## 1. 两个分母（现跑，⛔ 抄）

```bash
grep -rn "slog\.\(Error\|Warn\|Info\)(" --include="*.go" internal/ cmd/ | grep -v _test | wc -l   # 129
go list -deps ./cmd/wisp/ | grep -c 'CarlosShao'                                                  # 27
go list -deps ./cmd/wisp/ | grep -c 'internal/audio'   # 1（在闭包内）
go list -deps ./cmd/wisp/ | grep -c 'internal/ball'    # 1（在闭包内）
grep -rn "slog\.Debug(" --include="*.go" internal/ cmd/ | grep -v _test | wc -l                  # 1（A6，尺①不含）
grep -rnE "\b[a-z]+\.log(ger|\(\))\.(Info|Warn|Error|Debug)\(" --include="*.go" internal/ cmd/ | grep -v _test | wc -l  # 25（捕获式，不在 129 里）
```

- **闭包外只剩两包**：`internal/speech`、`internal/watchdog`（`comm -23 <(ls -d internal/*/) <(closure)`），两者主式子 **0 行** ⇒ **129/129 在闭包内**。
- 件里那 53 行的"闭包内 28 / 闭包外 25"**今天读法相反**：audio(7)+ball(23) 今天都在闭包内（resident 腿已接球与采集）。
- **捕获式的分母也漂了**：件写 21 枚，今天 **25 枚**（memory 16：`artifacts.go:102/271/321`、`open.go:236/240/343/507/522`、`retention.go:109/120/141/167/174/180`、**`writer.go:148/184` 新增**；agent 9：`loop.go:391/427/650/723/764/804/1016`＋**`compress.go:487/516` 新增**）⇒ 判语 §2(c) 的 agent 七枚/构造点 215 对，但 memory 14 枚漏了两枚。

## 2. 「其余 118 枚」怎么来的（⛔ 24、⛔ 19）

- 件说"其余 24 枚"、逐条只列 19；判语点出 `28 = 6 + 2 + 1 + 19` 才对。**两个数今天都不成立**：分母 129、且组变了（B 从 2 涨到 4；C' 是新的一枚）。
- 今天：**129 = 5(A) + 4(B) + 1(C') + 1(F) + 118(其余)**，逐条见 `130-r2-impl-part2.md` §4。该表每行 = `# | 站点 | 组 | 级 | 逐字锚 | 判 | 据`，**118 枚一枚不少**（判：113 `产不可达` / 12 `造不出` / 3 `可造` / 1 `非缺口`；级：43 Error / 43 Warn / 43 Info）。
- **判语 §2(c) 的 16 处行号复核**：四腿装点 `109/222/226/339` ✅、`logsink 149/153/184`（补发 173）✅、agent 七枚＋构造点 `215` ✅、B2 `goroutine.go:271` ✅、C1 `winsec_windows.go:143` ✅、F1 `secret/migrate.go:198` ✅、A 组 6 枚 ✅；**两处要修**：B1 = `logging.go:257`（**`379` 是 `flushLoop` 的 tick，不是 B 窗口**）；memory/agent 捕获式各多两枚（见 §1）。对照表在 part2 §5.1。

## 3. 具名 12 枚（A 6 ＋ B 4 ＋ C' 1 ＋ F 1）与"能不能造"

| 组 | 站点 | 级 | 消息（锚，逐字） | 判 | 据 |
|---|---|---|---|---|---|
| A1 | `internal/winsec/resolve.go:182` | Info | `winsec: sealing path resolver installed` | **可造** | 发声半**本笔现跑**：`go test ./internal/risk/ -run 'TestNothingMatches130Audit' -v` ⇒ stderr 首行 `level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2`；落盘半＝判语 §4（slo post 记录 0） |
| A2 | `internal/winsec/resolve.go:158` | Error | `winsec: refusing to install a path resolver into the sealing seam` | **本机造不出** | 本机 conformance probe 全过（`probes_passed=2`，上面那行就是）⇒ 要一台让 probe 失败的主机；⛔ 写成"不会掉" |
| A3 | `internal/winsec/resolve.go:166` | Error | `winsec: refusing to replace the already installed sealing path resolver` | **生产不可达** | 需第二次、不同实例的安装；`SetPathResolver(` 非测试调用者 **1**（`internal/risk/winsec_c26.go:21`） |
| A4 | `internal/winsec/resolve.go:151` | Error | `winsec: refusing to release the sealing path resolver` | **生产不可达** | 需 `SetPathResolver(nil)`，非测试调用者 **0** |
| A5 | `internal/winsec/resolve.go:147` | Info | `winsec: sealing path resolver left at the built-in floor` | **生产不可达** | 同 A4 的前置 |
| A6 | `internal/winsec/resolve.go:176` | Debug | `winsec: sealing path resolver installed again identically` | **生产不可达** | 同 A3，且 **Debug 被默认 Info 门挡** ⇒ 今天连 stderr 都不到 |
| B1 | `internal/observe/logging.go:257` | Warn | `observe: log retention sweep failed` | **本机造不出** | boot sweep 的 `os.ReadDir(<data>\logs)` 要在 `MkdirAll` 刚成功后失败；本机未造出（`379` 是同一条消息的 tick 版，**不在窗口内**） |
| B2 | `internal/observe/goroutine.go:271` | Warn | `goroutine outside the D38 roster (leak symptom)` | **可造** | 变异：把 `log-flusher` 从 D38 名册摘掉即出声（判语 §3 用的同形） |
| B3 | `internal/observe/goroutine.go:168` | Error | `goroutine panic recovered` | **本机造不出** | 需窗口内已 spawn 的 goroutine（`logging.go:123` 起的 `log-flusher`）panic |
| B4 | `internal/winsec/winsec_windows.go:143` | Warn | `winsec: seal cleared principals that stood on this object` | **可造** | 两向实测（判语 §4）：`slo -settle` post = **记录 1**（盘上，`path=…\logs`，时间戳早于任何 listener 记录）；`providers` 腿 = stderr 且 `find` 零文件。**它就是 B 窗口**（`InitLog` 内部封 `logs\` 那次）**且同时是 C 组**（`models ensure` 腿 `:323`（`loadModelConfig`）/:328（`openModelStore`）早于装点 `:339`） |
| C' | `internal/proc/singleinstance_windows.go:171` | Warn | `proc: AllowSetForegroundWindow(ASFW_ANY) was refused; …` | **本机造不出** | 需 `ASFW_ANY` 被拒（前台锁）。**且这条腿没有听众**：second launch 在 `cmd/wisp/resident_windows.go:88` 就 `return`，永不到 `:109` ⇒ 只到 stderr。**件与判语都没有这一枚** |
| F1 | `internal/secret/migrate.go:198` | Warn | `plaintext api_key migrated to DPAPI-protected references (D33)` | **生产不可达** | `MigratePlaintext` 非测试调用者 **0** |

## 4. 今天的"掉"面已经换了形状（判语 §4 之后）

- **装得上听众的腿**（`resident`/`run`/`secret`/`models ensure`/`slo`）：A、B 组的记录**今天落盘**（判语 §4：slo post 记录 0 = A1、记录 1 = B4）。
- **掉只剩三处**：①**无听众腿** —— `providers`/`doctor`/`panel-assets`/`models list`+`verify`/`version`/`help`/**`panel-inbound`**/ **second launch** / `cmd/balldebug`；这些腿上 `FlushEarlyLogRecords`（生产调用者 1＝`logsink.go:173`）与 `InstallAsDefault`（1＝`slo_windows.go:281`）都没有调用者 ⇒ 缓冲**永不冲刷**；②**缓冲溢出**（`earlyLogMaxRecords` 64 / 64 KiB，`logging.go:514-515`；溢出时补 `observe: early log buffer overflow` 带 `dropped`/`first_kept`/`capacity`）；③ 装听众**失败**的腿（`installLogSink` 报错 → stderr 兜底）。
- ⇒ 清单的"判"按**腿**判，不按时刻判：同一条站点在 `run` 腿落盘、在 `panel-inbound` 腿掉。

## 5. 两处注释：改前 → 改后（逐字）

**① `internal/observe/logging.go`（改前 `:486-491`，改后 `:486-495`）**

```go
// 改前（6 行）
// Why the handover is NOT in InitLog: ticket 132 reads the file seal at exactly
// that instant and needs the boot file to still hold zero bytes
// (logprivate_132_test.go's TestLogFileIsSealedBeforeTheFirstContentByte132r1).
// Moving the replay into construction would spend a seal reading this ticket has
// no business spending - and "the pipeline exists" is not yet "the process
// listens", which is the distinction the two hooks above make.
// 改后（10 行）
// Why the handover is NOT in InitLog: the reading a replay would spend is this
// ticket's own, not ticket 132's. Moving the replay into construction turns
// earlylog_takeover_130r1_test.go's "after InitLog (and a flush) the boot file
// holds 0 bytes" red - that case flushes before it measures, so it is the one
// reading that separates a replay from a buffer. Ticket 132's six cases stay
// green under the same move (6 PASS / 0 FAIL, measured in
// docs/evidence/s1/130-v1-verdict.md §6): they stat the file without flushing,
// so they read the writer's bufio, not the handover.
// What stands on its own is the distinction the two hooks above make:
// "the pipeline exists" is not yet "the process listens", and a pipeline built
// but not installed has no business writing the boot verdict into its own file.
```

**② `internal/observe/earlylog_takeover_130r1_test.go`（改前 `:13-20`，改后 `:16-22`）**（第 2 条 `InstallAsDefault` 的理由未动）

```go
// 改前（第 1 条，5 行）
//  1. ticket 132 reads the seal at the construction instant
//     (logprivate_132_test.go's TestLogFileIsSealedBeforeTheFirstContentByte132r1
//     requires the boot file to hold 0 bytes right after InitLog). A construction
//     that writes the boot records would make that reading unreadable, and the
//     seal question it answers is not this ticket's to spend.
// 改后（第 1 条，7 行）
//  1. the reading a construction-time replay would spend is this file's own, not
//     ticket 132's: moving the replay into InitLog turns
//     TestInstallAsDefaultTakesOverTheBootRecord130R1 red on its "after InitLog
//     (and a flush) the boot file holds 0 bytes" assertion. Ticket 132's six
//     cases stay green under the same move (6 PASS / 0 FAIL, measured in
//     docs/evidence/s1/130-v1-verdict.md §6): they stat the file without
//     flushing, so they read the writer's bufio.
```

**为什么这句必须换**：判语造 `M-2`（把补发挪进 `InitLog`）后，票 132 六枚 **6 PASS / 0 FAIL（rc=0）** ⇒ "会花掉 132 的封印判据"是**假话**。真代价是本笔自己的 `earlylog_takeover_130r1_test.go:182`（**带 `FlushNow`** 的"构造瞬间 0 字节"）变红。机制可复跑：`logprivate_132_test.go:68` 的 `info.Size()` **不带 flush**（量的是 `bufio`），本笔 `:179` 的 `FlushNow()` **带**。

## 6. 门禁（全现跑，主检出）

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 0 行输出 |
| `go test ./internal/observe/ -run '<130r1 四名>' -count=2 -v` | **0** | **RUN 8 / PASS 8 / FAIL 0** |
| `go test ./internal/observe/ -run '132r1' -count=2 -v` | **0** | **RUN 12 / PASS 12 / FAIL 0**（＝六枚 × 2） |
| `go vet ./internal/observe/` | **0** | — |
| `gofmt -l` 两枚本笔文件 | **0** | 空；`gofmt -l internal/ cmd/` 仍是既有的 3 枚 eol 伪红（`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`），⛔ 本笔未动 |
| `gofumpt -l`（`/d/work/base/gopath/bin/gofumpt.exe`）两枚 | **0** | 空；全树 54 行＝51 `.scratch/` + 上面 3 枚 eol 伪红（`tr '\\' '/'` 后 `grep -E 'internal/observe/'` ⇒ `NONE`） |
| `sh scripts/d22scan.sh` | **0** | clean：`248/42/116/25/39/116/631/148`（八 scope 与判语 §7 乙列逐字相同，**不降**） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（worst=229、hat=165） |
| CR 字节（`tr -cd "\r" \| wc -c`） | — | 两枚均 **0** |

## 7. 一条命令重验（§1/§2 全在里面）

```bash
cd '/d/work/workspace/projects plans/Wisp' && git rev-parse --short HEAD && \
grep -rn "slog\.\(Error\|Warn\|Info\)(" --include="*.go" internal/ cmd/ | grep -v _test | cut -d: -f1,2 | sort > /tmp/130r2-now.txt && \
grep -E '^\| [0-9]+ \| `[^`]+\.go:[0-9]+` \|' docs/evidence/s1/130-r2-impl-part2.md | grep -oE '`[^`]+\.go:[0-9]+`' | tr -d '`' | sort > /tmp/130r2-want.txt && \
printf 'census now=%s want=%s\n' "$(wc -l < /tmp/130r2-now.txt)" "$(wc -l < /tmp/130r2-want.txt)" && \
diff /tmp/130r2-now.txt /tmp/130r2-want.txt && echo CENSUS-OK && \
printf 'closure=%s audio=%s ball=%s\n' "$(go list -deps ./cmd/wisp/ | grep -c CarlosShao)" "$(go list -deps ./cmd/wisp/ | grep -c internal/audio)" "$(go list -deps ./cmd/wisp/ | grep -c internal/ball)"
```

**本笔现跑的输出（逐字）**：`census now=129 want=129` / `CENSUS-OK` / `closure=27 audio=1 ball=1`。

## 8. 交回项（不属本笔写面，⛔ 藏）

1. **同一句假话还活着**：`earlylog_takeover_130r1_test.go:183-185` 的 `Fatalf` 文案仍写"ticket 132 reads the file seal exactly there, and a replay this early makes that reading unreadable"——那是**字符串**（在 `:13-20` 注释写面之外），本笔⛔ 动。要换请另开一格。
2. **票面前提已腐坏 3 处**（本笔逐条现跑）：`:20` 引 `internal/winsec/resolve.go:157`（今天错一句 ⇒ 是 `:158`）；`:22` 引 `internal/observe/logging.go:204`（今天是 `:257`）；`:65/:78` 引 `cmd/wisp/leg_sink_nail_131_test.go`（今天叫 `cmd/wisp/leg_sink_nail_131_windows_test.go`）。
3. **`panel-inbound` 腿不在任何一份"无听众腿"名单里**（件 E 组只列 providers/doctor/panel-assets/models list）：它自述不装 sink（`cmd/wisp/panel_inbound.go:50`）却会走到 `secret.NewStore`（`:274`）与 `config.NewManager`（`:242`）⇒ 票 131 的枚举域漏了它。
4. **`AC#1` 的框本笔⛔ 勾**（翻框是验收方的事）；本笔只交读数与清单。
5. **写面外零字节**：`cmd/**`、`internal/{risk,winsec,tools,config,agent,session,memory,panel,ball,proc,audio,watchdog}/**`、`scripts/**`、`.github/**`、`tools/d22scan/**`、`go.mod`/`go.sum`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、`thresholds.go`／golden／`allowlist.txt`／三枚冻结件 —— 全部 `git status --porcelain` 干净（开工前核过，`git diff --cached --name-status` 空、无 `D`/`A`）。
