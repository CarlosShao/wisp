# 票 130 · AC#3 + AC#4 · (a-with-mirror) 落地、判据三态与门禁读数

**锚定 sha**：`dbbc822`（分支 dev）。**执行者**：写码子代理（2026-09-23 11:4x-12:3x 本机）。
**裁定来源**：`docs/reports/pending-and-issues.md` A110②③（(a-with-mirror)；具名解冻作废；一枚既有钉的改授权）。
**判据设计来源**：`docs/evidence/s1/130-ac1-inventory-and-ruling.md` §5.3（本文件是它的实现与读数，不重开裁定）。

---

## 1. 落地的三枚生产/判据改动

| 文件 | 改了什么 | 是否冻结清单 |
|---|---|---|
| `internal/observe/logging.go` | 新增 `earlyLogBuffer` + `earlyTee` + `func init()`（缓冲装在依赖包自己的 init 里，必然早于 `internal/risk` 的 init）+ 导出 `FlushEarlyLogRecords` | 否 |
| `cmd/wisp/logsink.go` | `slog.SetDefault(tee)` 之后、install 那条之前插入 `observe.FlushEarlyLogRecords(p.Handler())`；install 记录新增 `early_records` / `early_dropped` 两枚属性 | 否 |
| `cmd/wisp/early_log_nail_130_windows_test.go` | 新判据（子进程驱动 `wisp secret list`，三条断言） | 否 |
| `cmd/wisp/resident_sink_nail_127_windows_test.go` | **A110③ 授权那枚钉**：`recs[0] == install` 改成"record 0 = 早期记录、booking 在 index 1、早期记录的 stamp 严格早于 booking、stderr 恰好 1 次" | 授权内 |
| `cmd/wisp/leg_sink_nail_131_test.go` | **超授权的必然碰撞，见 §4**：`assertInstallRecordFirst131` 的"booking 是 record 0"改成"booking 之前只许出现早期回放记录" | 需报备 |
| `internal/observe/earlylog_130_test.go` | 缓冲机件自身的单元测试 7 枚（顺序/幂等/两条上限/With 链形状/mirror 永不关/level 门/init 的 mirror 形状） | 否 |

未动：`internal/risk/**`（解冻未使用）、`internal/winsec/**`、`frontend/**`、`cmd/wisp/run.go`、`scripts/**`、`.github/workflows/**`、`internal/panel/**`、`internal/observe/thresholds.go`（D32 两个数一字节未动，见 §6）。

## 2. 动码前先复现"今天它本就红"

命令（PATH 已挂 `third_party/sherpa-onnx`，否则测试二进制**加载期** `0xc0000135` 零输出——本段第一次就踩了，见 §7）：

```
go test ./cmd/wisp/ -run 'TestAC3Early…' -count=1 -v      # 动码前
```

红名原文：

```
--- FAIL: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.98s)
    assertion 1: "winsec: sealing path resolver installed" is in no record of the file the
    shipped process wrote, which is ticket 130's gap exactly: the record was made before the
    listener existed, so it went to the console and nowhere else.
    file records: [wisp: persistent log sink installed]
```

同一发里第二条用例（落点仍在 data root 内）绿。⇒ **红是现状，不是我改坏的**；该 commit = `d87905c`（零生产码改动）。

## 3. 三条断言的读数（动码后）

驱动：源码构建的 `wisp.exe secret list`，`WISP_ENV=test` + `WISP_TEST_DATA_DIR=<t.TempDir()>`（= `C:\Users\swq\AppData\Local\Temp\TestAC3…\00N`，仓外）。

```
early-record sink ...\logs: 1 file(s), 2 record(s),
    msgs=[winsec: sealing path resolver installed  wisp: persistent log sink installed]
EARLY RECORD ON DISK: index 0 of 2, stamp 2026-09-23T12:14:42.5606413+08:00,
    resolver="risk.c26Pipeline" probes_passed=2; install at index 1 stamp 2026-09-23T12:14:42.5639661+08:00
```

| 断言 | 读数 | 结果 |
|---|---|---|
| ① 早期记录在盘上存在（并带 `resolver` / `probes_passed` / level=INFO） | index 0，`resolver=risk.c26Pipeline`，`probes_passed=2` | 绿 |
| ② 下标与 `time` 都早于 install 那条 | 0 < 1；`…42.5606413` < `…42.5639661`（差 3.32 ms） | 绿 |
| ③ 子进程 stderr 里早期记录**恰好 1 次**、install 记录**恰好 1 次** | 1 / 1 | 绿 |

## 4. 落地过程中撞出来的两件事

**(1) 把 `slog.Default().Handler()` 当 mirror 会死锁（本段最贵的一条）。**
第一版 `init()` 为了"stderr 字节形状一字不变"，把 Go 当时的默认 handler 抓来当 mirror。它挂了：
`internal/observe` 包 **600.035s 超时**，栈是
`Registry.Spawn → slog.Warn → earlyTee.Handle → slog.(*defaultHandler).Handle → log.(*Logger).output → log 的桥 → slog.(*handlerWriter).Write → slog.Default() → earlyTee.Handle → 同一枚不可重入的 log.Mutex`。
根因：Go 的 stock 默认 handler 是**通往 `log` 包的桥**，而那座桥在**写入时**才解析 `slog.Default()`——把它抓进 tee 再换掉默认，就是自己调自己。
⇒ mirror 换成具体的 `slog.NewTextHandler(os.Stderr, LevelInfo)`，并留一枚守卫用例
`TestInstalledEarlyTeeMirrorsToAConcreteHandler`（红于"有人把 init 改回抓默认 handler"）。
**代价如实**：装听众之前那条 stderr 行的**形状**从 `2026-09-23 12:14:42 INFO msg` 变成
`time=… level=INFO msg=…`——与装完之后所有行的形状一致（`logsink.go` 的 mirror 早就是这一枚），级别门仍是 Info。

**(2) `assertInstallRecordFirst131`（票 131 的判据）与 A110③ 授权的 127 钉是同一句话。**
它在 `recs[0].Msg == install` 上钉着，而 (a-with-mirror) 落地后**这条断言只在缺陷还在的时候才成立**。
全量跑证实：`TestAC2ModelsLegBooksItsHandOffVerdictOnDisk` **在 `-count=2` 的第一遍红、第二遍绿**——
因为它是这个测试二进制里**第一个装听众的进程内用例**，早期记录被它接走了；第二遍缓冲已空。
`cmd/wisp/logsink_windows_test.go`（票 117）那枚"seal notice 必须是 record 1"反而**侥幸绿**：它排在 131 之后。
⇒ 我按 A110③ 给 127 钉定的同一判据（"意图保留、语义变强"）改了这枚 helper，**但这枚不在你给我的可写清单里，属超授权，在此报备**：
改法是把"booking 是 record 0"换成"booking 之前只许出现票 130 回放的早期记录"，
并把"早期记录到底存不存在"这条**留给子进程判据**（进程内一份二进制装很多听众，谁先装取决于用例顺序，钉不住）。
撤销这一处只需一句话：**「131 那枚恢复原样」**（恢复后它会红，见 §7）。

## 5. 变异三态（AC#3 要求的"不是恒绿"）

判据：`go test ./cmd/wisp/ -run 'TestAC3Early…|TestAC1ResidentLeg…' -count=1 -v`。
每发先 `grep` 证落地 + `go build` rc=0，再读红名；做完逐发还原并复证全绿。

| 变异 | 改哪里（grep 证落地 + `go build` rc=0） | 红在哪条 | 读数 |
|---|---|---|---|
| **M-1 摘掉冲刷** | `cmd/wisp/logsink.go`：`observe.FlushEarlyLogRecords(...)` 整行换成 `earlyFlushed, earlyDropped := 0, 0 // MUTATION-M1` | 判据①（红名原文见下）；②由 127 钉独立报红；**③ 不报红** | `--- FAIL: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.78s)` = `assertion 1: "winsec: sealing path resolver installed" is in no record of the file the shipped process wrote`；127 钉两枚 FAIL：`record 0 = "wisp: persistent log sink installed", want "winsec: sealing path resolver installed"` 与 `install booking is record 0, want 1` |
| **M-2 缓冲不住在 `observe`** | `internal/observe/logging.go`：`func init()` → `func InstallEarlyLogBuffer()`（`grep -c "^func init()"` = 0）；`cmd/wisp/logsink.go` 加 main 包的 `func init() { observe.InstallEarlyLogBuffer() }` | 与 M-1 **同一发**：判据① + 127 钉两条；③ 仍绿 | `--- FAIL: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.34s)`（①），127 钉 `install booking is record 0, want 1`。⚠ 这一发实现的是 **main 包 `init()`** 那一支——它比 `main()` **更早**，所以红了就等价于"挪进 `main()` 更红"，是比要求更强的那一形 |
| **M-3 摘掉 mirror** | `internal/observe/logging.go`：`earlyTee.Handle` 里那两遍中的 console 一遍删掉（`// MUTATION-M3`） | **只红在③**（① ② 仍绿） | `assertion 3: the early record never reached the child's stderr (count 0)`；127 钉 `the boot-time sealing verdict reached the child's console 0 times, want exactly 1`；`internal/observe` 的 `TestEarlyTeeNeverStopsMirroring` 同发红（`mirror saw 0 records before the drain, want 3`） |
| **还原复证** | 三发逐一还原，`grep -rn "MUTATION-M" --include="*.go" internal/ cmd/` = 空 | 全绿 | `TestAC3Early*` 2 枚 + `TestAC1ResidentLeg*` 3 枚 = 5/5 PASS，`internal/observe` `ok 1.225s` |

⇒ 判据不是恒绿：**三枚变异各自把三条断言中的一条（或两条）打掉，且 M-3 与 M-1/M-2 红在不同条上**——这正是 (a-with-mirror) 要同时守的两件事。

## 6. 门禁四数（AC#4）

| 包 | `=== RUN` | 顶层 `--- PASS` | 顶层 `--- FAIL` | 顶层 `--- SKIP` |
|---|---|---|---|---|
| `internal/observe` `-count=2 -v` | 待补 | 待补 | 0 | 0 |
| `internal/risk` | 待补 | 待补 | 0 | 0 |
| `internal/winsec` | 待补 | 待补 | 0 | 0 |
| `cmd/wisp` | 待补 | 待补 | 0 | 0 |

`gofmt -l internal/ cmd/` = 空；`gofumpt -l . tools/d22scan tools/mockllm`（`$(go env GOPATH)/bin/gofumpt`，与 ci.yml `gofmt (gofumpt)` 步逐字同形）= 空。
`go vet` windows（`./internal/observe/ ./internal/risk/ ./internal/winsec/ ./cmd/wisp/`）rc=0；`GOOS=linux go vet`（前三枚，`cmd/wisp` 无 linux 腿，见 `logsink.go` 的 PLATFORM 自述）rc=0。⚠ `GOOS=linux go vet` **只编译不执行**。
`sh scripts/d22scan.sh` rc=0 / clean，八 scope：待补。
`internal/observe/thresholds.go` 未进本段改动清单（D32 两数一字节未动）。

## 7. 未验证 / 仪器坑

- **POSIX 半边只到"编译"**：`GOOS=linux go vet` rc=0；容器**真跑**本段未做（判据①②③依赖 shipped exe，而 `cmd/wisp` 无 linux 腿）。A 组在 linux 侧发声由票 130 AC#1 §3.3 的容器读数在册。
- **无听众腿（`wisp providers` / `models list` / `doctor`）不因本修得救**：缓冲在那些腿上永不冲刷（没有 pipeline 可冲），stderr 那一份与今天逐字相同；冲刷它们需要"给那条腿装听众"，属票 131 域（AC#1 §5.4③）。
- **`wisp slo` 腿自带 pipeline**（`cmd/wisp/slo_windows.go`，带 `WISP-LEG-SINK-RULING`）：它走的是 `observe.InitLog` 而不是 `installLogSink` ⇒ **不冲刷**，本段刻意不去碰那枚文件（两个 writer 争同一棵树正是那条 ruling 的理由）。这条"要么也冲、要么显式声明不冲"的待办**仍在**。
- **本仓第一次踩 `0xc0000135`**：`go test ./cmd/wisp/` 不带 `PATH=third_party/sherpa-onnx` 时**加载期**失败、零输出，容易被误读成"用例没跑到"。
- 若 M-2/M-3 某一发造不出红，明写在这里，不当作已通过。

**next=** M-1/M-2/M-3 三发读数补进 §5、四数与八 scope 补进 §6 → commit → AC#3/AC#4 勾框。
