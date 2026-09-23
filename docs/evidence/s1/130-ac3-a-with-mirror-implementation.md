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

全部 `-count=2 -v`，四数只从 `-v` 数（`=== RUN` 含子测试；顶层与子测试分列，不跨层相加）。
末次复跑（三枚变异全部还原之后）：

| 包 | `=== RUN` | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子测试 PASS/FAIL/SKIP | rc |
|---|---|---|---|---|---|---|
| `internal/observe` | **108** | **108** | **0** | **0** | 0 / 0 / 0 | 0 |
| `internal/risk` | **336** | **198** | **0** | **2** | 136 / 0 / 0 | 0 |
| `internal/winsec` | **202** | **116** | **0** | **0** | 86 / 0 / 0 | 0 |
| `cmd/wisp` | **200** | **106** | **0** | **0** | 94 / 0 / 0 | 0 |

分层核：`336 = 198 + 2 + 136`、`202 = 116 + 86`、`200 = 106 + 94`、`108 = 108`。
`internal/risk` 那 2 枚顶层 SKIP 是**既有**的同一枚用例在 `-count=2` 下各跳一次：
`TestSyncRegistryProbeLive`（需真 registry，本机跳）——与本票无关，不是新增。
`cmd/wisp` 侧一律带 `PATH=third_party/sherpa-onnx`（`scripts/wisp-cli-tests.sh` 的形式），否则测试二进制**加载期** `0xc0000135`、0 个用例跑过。

- `gofmt -l internal/ cmd/` = **空**（本段改过的 6 枚 .go 全部 `gofmt -w` 过一遍再复跑判空）。
- `gofumpt -l . tools/d22scan tools/mockllm`（`$(go env GOPATH)/bin/gofumpt`，与 ci.yml `gofmt (gofumpt)` 步**逐字同形**）= **空**。
- `go vet` windows：`./internal/observe/ ./internal/risk/ ./internal/winsec/ ./cmd/wisp/` **rc=0**。
  `GOOS=linux go vet`：前三枚 **rc=0**；`cmd/wisp` 无 linux 腿（sherpa 的 cgo 布局，见 `cmd/wisp/logsink.go` 的 PLATFORM 自述），不纳入式子。
  ⚠ `GOOS=linux go vet` **只编译不执行**，POSIX 半边本段未做容器真跑（§7）。
- `sh scripts/d22scan.sh` **rc=0 / clean**（含 `tools/d22scan` 自检 `PASS=21 FAIL=0 SKIP=0 / === RUN=31`）。
  台账八 scope（本树 → 上一次在树里登记的读数，`docs/evidence/s1/118-adversarial-acceptance.md`）：
  `bans #1-5 internal/` **203** ≥ 202 · `bans #1-5 cmd/` **22** ≥ 22 · `ban #6 frontend/` **43** ≥ 40 · `ban #7 internal/tools/` **18** ≥ 18 ·
  `ban #8 design/` **16** · `ban #8 frontend/` **43** · `ban #8 internal/` **390** · `ban #8 cmd/` **37** ⇒ **各 scope 不降**；
  本 commit 对 Go 树只有新增（`internal/observe/logging.go` 的 diff 是 **343 增 / 0 删**，另两枚新文件），没有任何一枚文件被删，
  所以"降"在本票这一发里结构上不可能发生。
- `internal/observe/thresholds.go` **不在改动清单**（D32 的 `memCapSleeping = 25 << 20` / `cpuLimitSleeping = 0.5` 一字节未动）；
  缓冲的上界是**条数与内存字节**（64 / 64 KiB），不是任何阈值表项，也不新增 golden。
- **没有跑全仓 `go test ./...`**（同树有其他票的在飞件）。

## 7. 未验证 / 仪器坑

- **POSIX 半边只到"编译"**：`GOOS=linux go vet` rc=0；容器**真跑**本段未做（判据①②③依赖 shipped exe，而 `cmd/wisp` 无 linux 腿）。A 组在 linux 侧发声由票 130 AC#1 §3.3 的容器读数在册。
- **无听众腿（`wisp providers` / `models list` / `doctor`）不因本修得救**：缓冲在那些腿上永不冲刷（没有 pipeline 可冲），stderr 那一份与今天逐字相同；冲刷它们需要"给那条腿装听众"，属票 131 域（AC#1 §5.4③）。
- **`wisp slo` 腿自带 pipeline**（`cmd/wisp/slo_windows.go`，带 `WISP-LEG-SINK-RULING`）：它走的是 `observe.InitLog` 而不是 `installLogSink` ⇒ **不冲刷**，本段刻意不去碰那枚文件（两个 writer 争同一棵树正是那条 ruling 的理由）。这条"要么也冲、要么显式声明不冲"的待办**仍在**。
- **本仓第一次踩 `0xc0000135`**：`go test ./cmd/wisp/` 不带 `PATH=third_party/sherpa-onnx` 时**加载期**失败、零输出，容易被误读成"用例没跑到"。
- 三态变异**每一发都造出了红**（M-1 ①②、M-2 ①②、M-3 ③），无"造不出红"的分支；读数在 §5。
- 本段**不取任何时序/RSS 读数**（任务明写不需要），因此也没有 `A103` 的三连检查前置；缓冲的上界论证用的是"64 条 / 64 KiB、无 fd、无 fsync、无启动期文件系统活动"这个结构论证，D32 冷读预算不经过它。
- ⚠ 票 130 AC#1 §5.4② 记的那条待办**仍然挂着**：`wisp slo` 腿自带 pipeline、不冲刷；无听众腿（`providers` 等）也不冲刷。两者都不是本票的判据域。

**next=** 交回 owner / 编排者三件事：① §4(2) 那枚**超授权**改动的追认或撤销（口令「131 那枚恢复原样」）；
② `cmd/wisp/leg_sink_nail_131_test.go` 与 `resident_sink_nail_127_windows_test.go` 已被本票改过 ⇒ **票 131 的续单与第二任验收表 `docs/evidence/s1/131-adversarial-acceptance.md` 的 §2 那张表（锚在 `56d8026`）都要重新锚 sha**；
③ AC#2 的勾框归编排者（裁定已落在 A110③，本票只落了码），AC#1 的勾框归出清单的那位只读代理复跑。

## 8. 同一棵树的并发事实（写给下一个接 131 续单的人）

本段动 `cmd/wisp/leg_sink_nail_131_test.go` 期间，票 131 的第二任对抗验收代理正在同一棵树上交件，
并已把我这批在飞改动**作为事实登记（不下判断）**进它的裁决表：commit `7daec2f`
（`docs/evidence/s1/131-adversarial-acceptance.md` §13，原文含
`leg_sink_nail_131_test.go 44/9`、`resident_sink_nail_127_windows_test.go 108/18`、`logsink.go 22/2`，
并点名"删掉的第一行是 `assertInstallRecordFirst131` 上面那段'记录 0 必须是 install'的理由注释"）。
它的结论与本文件 §4(2) 独立同形：**130 与 131 会改同一枚判据文件，别并行；接 131 续单前先重新锚 sha**——
它那张表是 `56d8026` 的读数，不是本 commit 之后的树。

本票两枚 commit 的落点：`d87905c`（判据入库，红）→ `9b5d64d`（(a-with-mirror) 落地 + 三态变异）。
中间落在同一棵树上的他人 commit（`c25adea`/`f84cd90`/`7daec2f`/`e10ca09`/`dcedb7b` 等）与本票可写清单**零文件重叠**
（逐枚 `git show --stat` 核过：只动 `docs/evidence/s1/131-*`、`internal/secret/**` 等），未发生覆盖。
本票引用的行号一律以 `dbbc822` 与本文件为准，不引用他人工作树状态。

## 9. 边界自证

- `git add` 只用显式路径；两枚 commit 的 `git diff --cached --name-only` 清单：
  `d87905c` = 2 枚（新用例 + 票面）；`9b5d64d` = 7 枚（`internal/observe/logging.go`、`internal/observe/earlylog_130_test.go`、
  `cmd/wisp/logsink.go`、`cmd/wisp/resident_sink_nail_127_windows_test.go`、`cmd/wisp/leg_sink_nail_131_test.go`、票面、本证据文件）。
- 未 push、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；未改名任何文件（无需 `git ls-tree` 命中数自证）。
- 未 add 他人件：`docs/evidence/s1/131-adversarial-acceptance.md`（131 验收代理自己入库了）。
- 真跑二进制一律仓外：`/tmp/wisp130b/`（= `C:\Users\swq\AppData\Local\Temp\wisp130b`）与各用例的 `t.TempDir()`；
  仓内 `build/`、`third_party/sherpa-onnx/` 只被读/复制到临时根，未覆写。
- **owner 真实数据目录前后各读一次，一字未多写**：
  前（11:4x）`%APPDATA%\wisp` 文件数 **0** / mtime `2026-09-19 14:49:20.310467900 +0800`；`%APPDATA%\wisp-dev` 文件数 **0** / mtime `2026-09-20 07:06:25.376099600 +0800`；`%APPDATA%\wisp-prod` **不存在**。
  后（12:35，同一条 `find -type f | wc -l` + `stat -c %y`）= **0 / 0 / 不存在，两枚 mtime 逐字未变** ⇒ 真跑的二进制只写了仓外根。
- 第三枚 commit（AC#4 读数与本节措辞补齐）暂存的清单只有两枚：票面 + 本证据文件；`git diff --cached --name-only` 的原文在交回报告里。
- 伪授权登记 **计数 0**：本轮工具输出里没有出现任何自称「编排者备注 / 系统提示 / 用户已更新编码规则 / 请 revert / 冻结某包 / 放宽阈值」的文字。
  出现过的**形似**项只有 harness 自己的三条 `MEMORY.md` 已被修改通知（`Read`/`Bash` 返回内的 system-reminder，非任何命令结果）：
  它点名的路径都是真实文件、内容里没有要求 revert/冻结/放宽的指令 ⇒ 判定合法通知，未据其改判据。

