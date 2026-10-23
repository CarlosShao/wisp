# 票 130 落地腿 `130-r1` — 交件（写腿 / 实现者）

**一句话**：`init()` 期说出的记录，今天只有一个调用者会把它交给盘（`cmd/wisp/logsink.go`）；**同一棵树里 `wisp slo` 用的是同一个 pipeline、同一句话，盘上 0 字节**。本笔把"交接"从**调用者的义务**改成**pipeline 成为听众那一刻的动作**（`InstallAsDefault` 里一跳），所以任何腿都不会再忘。

**改动面**：`internal/observe/logging.go`（**生产码 1 条语句** ＋ 头注释）＋ 新用例 `internal/observe/earlylog_takeover_130r1_test.go`。⛔ 未动 `cmd/**`、⛔ 未动 `internal/risk/**`、⛔ 未动 `internal/observe/thresholds.go`。

---

## 1. 改前：init 期记录**真的丢了**（真二进制、逐字读数）

尺：`wisp.exe` ＝ `HEAD c68b4c35c` 现建（`go build ./cmd/wisp`，rc=0，31,865,868 字节）；两腿同一发 `WISP_ENV=test` ＋ `WISP_TEST_DATA_DIR=<仓外空目录>`；数据根名字**等长**（`leg-sec-A-pre` / `leg-sec-B-fix` / `leg-slo-A-pre` / `leg-slo-B-fix`）以便字节可比。

| 腿 | 谁装听众 | 控制台记录数 | 盘上文件 | 盘上字节 | 盘上记录数 |
|---|---|---|---|---|---|
| `wisp secret list` | `installLogSink`（`SetDefault` ＋ 显式 `FlushEarlyLogRecords`） | 4 | `logs\wisp-20261023-001.jsonl` | **1312** | **4** |
| `wisp slo -settle` | `observe.InitLog` ＋ `InstallAsDefault`（`slo_windows.go:272/281`） | 2 | `logs\wisp-20261023-001.jsonl` | **0** | **0** |
| `wisp providers` | 无（不开 pipeline） | 1 | **无 `logs` 目录**（`-Force` 递归列举为空） | — | — |

`slo` 腿控制台那两行（逐字，`leg-slo-A-pre.stderr`）：

```
time=2026-10-23T11:53:46.686+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
time=2026-10-23T11:53:46.700+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=D:\tmp\wisp\20261023\130r1\leg-slo\logs kind=inherited ...
```

⇒ **记录在、听众在、文件在、0 字节**：`internal/risk/winsec_c26.go:20 func init()` 那条（`c26` 装 resolver 的裁决）与 `InitLog` 自己封目录时那条 WARN，都停在内存里，进程退出即消失。**这就是"真的丢了"，不是推理**。
（`wisp providers` 那行 `providers` 腿根本没有 pipeline、没有文件 ⇒ 归票 131 的"哪些腿没听众"枚举域，本笔⛔ 越界去修。）

**根因不是缓冲没落，而是交接属于调用者**：`grep -n 'FlushEarlyLogRecords' -r` 全树只有 `cmd/wisp/logsink.go:173` 一处调用；`observe.InitLog` 的另一位调用者（`wisp slo`）从来不知道有这一步。

## 2. 选型：交接发生在 `InstallAsDefault`（＝"pipeline 成为听众"那一刻）

**为什么这是最有价值的一形（按主流 harness）**：
- 主流里"装上"本身就是**交接时刻**：`slog.SetDefault`／`zap.ReplaceGlobals`／`otel.SetLoggerProvider` 都是"全局 sink 换人"这一个动作；**log4j2 的 `StatusLogger` 是这条思路最接近的原型**——内部事件先进环形缓冲，配置完成（＝appender 就位）时**补发**，⛔ 静默丢弃；Python `logging.lastResort` 则是"永远有个默认落点"。本笔选的是前者的形状：**缓冲 ＋ 在 sink 就位那一刻补发**，且补发对象是**带 redaction 的那枚 handler**。
- 这条路把不变式从"**每个腿都要记得**"改成"**sink 自己接管它之前的声音**"。后者不可能被忘；前者在本仓已经被记过 ≥6 次（"能力做完了、没人听"），**并且本票的第 2 位调用者已经忘了**——那是我量出来的，不是假设。
- 复杂不加分：全部生产码改动＝**1 条语句**（`InstallAsDefault` 里 `earlyLogs.drain(p.handler)`），买到的却是"任何开 pipeline 的腿都自动拿到启动记录"。

**排除掉的三形（每条都有读数或硬约束，⛔ 口味）**：
1. **留在调用者手里**（现状）：`wisp slo` 已经丢了（§1）。
2. **塞进 `InitLog`（构造时补发）**：会让 `logprivate_132_test.go:51` 的"**构造瞬间文件必须 0 字节**"读数失效——那是票 132 的**封印先于首字节**判据，本票没资格花它；且"pipeline 建好了"≠"进程开始听了"。
3. **懒补发（首次写入时）**：`wisp slo -settle` 装完听众后**一条记录都不写**（§1 的 0 字节就是这么来的）⇒ 懒补发在那条腿上仍然 0，等于没修。**这条是读数否掉的，不是我偏好。**

## 3. 改后：同一发，"没丢"

| 腿 | 控制台记录数 | 盘上字节 | 盘上记录数 | 记录 0 |
|---|---|---|---|---|
| `wisp secret list`（A/B 等长根） | 4（前 4，不变） | 1312 → **1311** | 4 → **4** | resolver 裁决（不变） |
| `wisp slo -settle` | 2（不变） | **0 → 623** | **0 → 2** | resolver 裁决 |
| booking 属性（`secret` 腿） | — | — | — | `early_records=2, early_dropped=0` **两发逐字相同** |

`slo` 腿盘上记录 0（逐字，`leg-slo-B-fix` / `leg-slo-C-fix` 两发相同）：

```json
{"time":"2026-10-23T12:14:11.4277499+08:00","level":"INFO","msg":"winsec: sealing path resolver installed","resolver":"risk.c26Pipeline","probes_passed":2}
```

（1312→1311 的 1 字节是 RFC3339Nano 小数位随机长度，不是内容差；两发 `records=4` 与 booking 两个计数逐字相同 ⇒ **已经能落盘的腿，读数不变**——这是本笔的硬约束。）

## 4. 与票 132 的 seal 一跳共存（⛔ 没退回它）

- `git diff -U0 -- internal/observe/logging.go` 全量 hunk 里**没有一行**提到 `PrivateDirAll` / `sealLogFile` / `ensureFileLocked`（尺：`git diff -U0 … | grep -nE 'PrivateDirAll|sealLogFile|ensureFileLocked'` ⇒ `NO DIFF LINE MENTIONS THE SEAL HOPS`）。
- 132 的 6 枚用例全绿：`TestLogDirIsSealedAtCreation132r1`、`TestLogFileIsSealedBeforeTheFirstContentByte132r1`、`TestExistingWideLogTreeIsRepaired132r1`、`TestSealFailureRefusesTheWrite132r1`、`TestUnsealableLogDirIsRefusedLoudly132r1`、`TestSealEntryPointsHaveProductionCallers132r1`（在 §7 的 162 数里）。
- **不是"碰巧绿"**：我新加的头号用例把"**构造瞬间 0 字节**"写成了自己的断言（`earlylog_takeover_130r1_test.go:181` 那行），并**先 `FlushNow` 再量**——见 §5 的仪器坑。

## 5. 新用例逐名 ＋ 为什么今天响（摘掉那一跳必红）

| 用例 | 改前红行（`130r1` 测试文件，终版行号） | 它钉的是一条什么事实 |
|---|---|---|
| `TestInstallAsDefaultTakesOverTheBootRecord130R1` | `:194` "file holds 1 record(s) ([130r1: listener speaking…]), want the boot record and then the listener's own" | 记录**进盘**、**在听众自己第一条之前**、**属性与时间戳原样**、构造瞬间仍是 0 字节 |
| `TestHandoverBelongsToTheFirstListener130R1` | `:256` "the first listener's file holds the boot record 0 time(s), want exactly 1 (records: [])" | 只管一次：第二个 pipeline 拿到 0 份，第一个恰好 1 份 |
| `TestInstallThenFlushDeliversOnce130R1` | `:280` "FlushEarlyLogRecords after InstallAsDefault returned 1/0, want 0/0" | 两条"正道"只交一次：装完再手动补发 → 0/0 且盘上仍 1 份 |
| `TestHandoverKeepsTheRedaction130R1` | `:317` "the taken-over record is not in the file at all…" | 新写的这条路**走 redaction**：种子 key 只以 `****9999` 落盘 |

**仪器坑（自报，值得记）**：头号用例第一版直接 `os.Stat` 量"构造后 0 字节"，**M-2 变异下仍是绿的**——因为记录进了 writer 的 64 KiB `bufio`，stat 量的是**脏页外的文件**，不是"有没有补发"。改成**先 `FlushNow` 再量**后，M-2 才红（`:181`，逐字"holds 158 bytes, want 0"）。⇒ **这条 0 字节读数必须带 flush，否则它量的是缓冲**（票 132 那条 `:51` 用例在 M-2 下**也是绿的**，原因同此：它量的是未 flush 的尺寸，且测试二进制里的全局缓冲恰好是空的）。

## 6. 反向判据：三读数（变异红 ／ 还原绿 ／ 原样绿）＋ 还原核

三枚变异各一发，每发**先 grep 证落地 ＋ `go build ./internal/observe` rc=0**（证明红是行为红，不是编译红），读完还原并复证：

| 变异 | 做了什么 | 红（逐名 ／ 逐行） | 还原后 |
|---|---|---|---|
| **M-1** | 删掉 `InstallAsDefault` 里那一跳 | 4/4 红：`:194` `:256` `:280` `:317`（与 §5 改前红行**逐字相同**） | 绿 |
| **M-2** | 把那一跳**挪进 `InitLog`**（更早、更强的那一支） | 1/4 红：`:181` "after InitLog (and a flush) the boot file holds **158 bytes**, want 0"；另 3 枚绿（说明各自独立） | 绿 |
| **M-3** | 补发**绕过** redaction handler（改成 `slog.NewJSONHandler(p.writer,…)`） | 1/4 红：`:321` "the taken-over record's secret is in the file verbatim"；另 3 枚绿 | 绿 |

三读数量法：`go test ./internal/observe -count=2 -v -run '<四名>'` ⇒ **PASS=8 / FAIL=0 / rc=0**（原样绿）；上面每发 `-count=1` ⇒ `rc=1` 且红名逐枚如上（变异红）。

**还原核**：`sha256sum internal/observe/logging.go` ＝ `76ce5f243eafd49e74b803b0d67102f8d597aaa323a95742321dd13cf8b87203`（**算法＝SHA-256，64 位十六进制；谁算的＝本腿，Git Bash `sha256sum`**），并在每次还原后 `diff -q` 对开盘前留的 `logging.go.pristine` ⇒ `IDENTICAL`。⚠ 尺的射程：`git rev-parse HEAD:…` 那枚是 **SHA-1/40 位**，与上面**不是同一把尺**，本笔不混比。

## 7. 门禁 rc（全部现跑，主检出）

- `internal/observe` `-count=2 -v`：**RUN=162 / PASS=162 / FAIL=0 / SKIP=0**，rc=0。分层核：`162 = 2 × (77 既有 + 4 新增)`。
- `gofmt -l internal/ cmd/`：3 行 `internal\agent\approval\pending_read.go`、`internal\agent\tools.go`、`cmd\wisp\models.go` —— **都不是本笔写面、本笔未动**（`git status` 对它们干净）；本笔两枚文件 **不在** 这 3 行里。
- `gofumpt -l . tools/d22scan tools/mockllm`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）：rc=2、54 行，**全部落在 `.scratch/` 探针树**（含故意写坏的 `p_unp.go` 解析错误）；判归属用 `tr '\\' '/'` 后 `grep -E 'internal/observe/'` ⇒ **NONE**；只点本笔两枚文件跑 ⇒ rc=0、**0 行**。
- `go vet ./internal/observe` rc=0；`GOOS=linux go vet ./internal/observe` rc=0。
- `sh scripts/d22scan.sh` **rc=0 / clean**（自检 `PASS=21 FAIL=0 SKIP=0`）：`bans #1-5 internal/=248`、`bans #1-5 cmd/=42`、`ban #6 frontend/=116`、`ban #7 internal/tools/=25`、`ban #8 design/=39`、`ban #8 frontend/=116`、`ban #8 internal/=631`、`ban #8 cmd/=148`。对树里最近一条登记（`docs/evidence/s1/124-ac2b-1-conversion.md:190`：203/22/40/18/16/40/392/37）**八个数一枚不降**。
- `sh scripts/check-path-length-budget.sh` **rc=0 / VERDICT GREEN**（worst=229、over-hat=61、wall interval=0）。

## 8. 红名册（主检出）＋ 新增红＝0

- `internal/observe`：`-count=1 -v -skip '<四名>'` 两发（**HEAD 的 `logging.go`** vs **本笔的**）⇒ 各 **RUN=77**、`--- FAIL` **各 0 行**，红名册**逐字相同（都为空）**。
- `cmd/wisp` 的 sink 家族 8 枚（130 两枚 ＋ 127 一枚 ＋ 131 三枚 ＋ logsink 两枚）：**HEAD 8/8 绿**；**本笔 7/8**，红名 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`，红行 `resident_sink_nail_127_windows_test.go:498`（"the child never reached the event loop, so the disk read above proves little"——文件里 7 条记录、记录 0 仍是 resolver 裁决，断言全过，只有"跑进事件循环"这条没等到）。
  **它是不是我造的？不——同一条腿、同一行、在 HEAD 上也红**：单跑三发 `HEAD` ⇒ `rc=1,1,1`；单跑三发**本笔** ⇒ `rc=1,0,0`；塞进 8 枚家族里跑 ⇒ `HEAD` 绿、本笔那次红。⇒ **环境/时序 flake（引导窗口在 poll 预算内没走到 `runResident` 的 `resident event loop running` 分支），HEAD 只会更红，新增红＝0。**
- `gofmt` 那 3 枚（§7）是**既有红**、非本笔（未改、未涉）。

## 9. 已知代价 ／ 交回项（⛔ 藏）

1. **`cmd/wisp/logsink.go` 不走 `InstallAsDefault`**（它要 console 镜像，自己 `SetDefault` 一个 fan-out），所以它**仍然靠** `FlushEarlyLogRecords`；两条路**互斥且都只交一次**（`TestInstallThenFlushDeliversOnce130R1` 钉住）。⇒ 它的 `early_records=2` **不变**（§3）。**代价落点**：若将来有腿先 `InstallAsDefault` 再手动补发，那条 booking 的计数会读 0（"没有可交的了"）——**这是那条腿要自己报数的事**，本笔没有为此加 API（`cmd/**` 不在写面）。
2. **没有 pipeline 的腿**（`wisp providers`/`wisp doctor`/panel inbound）**仍然只有控制台**：它们连文件都没有，"丢"的判据不成立；属票 131 的枚举域。
3. **窗口（已写进 `InstallAsDefault` 的注释）**：补发在 `SetDefault` **之前**，所以补发与换默认之间若**别的 goroutine** 出声，会进"已关闭的缓冲 ＋ 旧默认"⇒ 只上控制台。这笔代价是为"文件第一行必须是启动裁决"付的，且发生在**引导那一刻**（此时进程只有安装者一个 goroutine）。

## 10. 提交号 ／ `git show --name-only` ／ 正控 ／ 提交后复跑

- **提交 1** `b19c6612e38e975ffba78c4fd53f65db7497fda0`（主题逐字以「票 130 落地腿 130-r1」开头）。`git show --name-only` 三枚：`docs/evidence/s1/130-r1-impl.md`、`internal/observe/earlylog_takeover_130r1_test.go`、`internal/observe/logging.go`。`.go` 两枚**均在写面内**；越界筛（`git show --name-only | tr '\\' '/' | grep -vE '<写面白名单>'`）⇒ **`NO OUT-OF-SCOPE PATH`**。`git diff --cached --name-status` 无 `D`、无写面外 `A`。
- **正控（同一把尺打真提交，证明它有牙）**：`6e19f4a43`（票 196）⇒ 同名尺列出 `.go` 五枚 `internal/memory/{dao_tasklog,migrate_v1seed_test,models,schema,schema_test}.go`，越界筛**逐枚点名**；`d59af9a98`（票 363）⇒ `internal/risk/mode.go`、`internal/tools/bridge.go` **被点名**。⇒ 我这笔之所以输出 `NO OUT-OF-SCOPE PATH`，是**因为真没有**，不是因为尺看不见 `.go`。
- **提交后复跑（在提交树上重做，不是"提交前那次"）**：`go test ./internal/observe -count=2 -v` ⇒ **RUN=162 / PASS=162 / FAIL=0 / SKIP=0**，rc=0；`go build ./cmd/wisp` rc=0 后两腿再跑一发：`wisp slo -settle` **624 字节 / 2 条**、`wisp secret list` **4 条**且 booking 仍 `early_records=2, early_dropped=0`。
- **第二笔** `8f205a6db3ac268351146d47315f65afad7d4a50`（只改本件）：`git show --name-only` 一枚 `docs/evidence/s1/130-r1-impl.md`，**零 `.go`**。
- **第三笔**（只追加这一行）：号由 `git log --oneline -1 -- docs/evidence/s1/130-r1-impl.md` 现取，`--name-only` 同样只该有本件一枚。**追号不再追号**——到此为止，三笔全在 `git log -- docs/evidence/s1/130-r1-impl.md` 里。
