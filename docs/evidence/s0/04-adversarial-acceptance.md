# T04 对抗验收报告 — 04-sqlite-core

- 验收角色：T04-adv（对抗验收，独立上下文；只读仓库 + 临时目录实验，唯一仓库写入为本报告）
- 验收时间：2026-09-19T09:00Z ~ 09:45Z
- 被验收对象：T04-impl 的 6 个 commit（`83acde2` skeleton → `a309579` DAO → `2cede0e` retention →
  `f27869e` privacy → `41d707c` concurrency/crash tests → `6d506a1` handoff；均为当时 HEAD `8c5a72c`
  的祖先链，验收期间 orchestrator 新增的票 05 认领 commit 不涉及票 04 代码）
- 验收环境：Windows 11 x64，Git Bash；Go `go1.27.1`（`D:\work\base\go\bin`）；
  GCC `(Rev3, Built by MSYS2 project) 16.2.0`（`E:\work\base\msys64\mingw64\bin`，与 BUILD.md §1 一致）；
  `GOPROXY=https://goproxy.cn,direct`
- 边界遵守：T02-impl 在途文件（`scripts/spike/`、`docs/SLO.md`、`docs/PRECHECK.md`、`docs/evidence/s0/json/`）
  **未触碰、未评判**；全程显式路径，未使用 `git add -A`；全部实验在 `%TEMP%\t04adv\` 完成
  （独立探针程序 + `go test -overlay` 注入测试二进制 + 独立验证器，仓库磁盘零改动）。

## 逐项裁决

### 1. 全量复跑（build / vet / test / race）— PASS（`-race` 一项 FAIL，见 MAJOR-2）

```
CGO_ENABLED=1 go build ./...           → BUILD_OK（整仓，含 cmd/wisp sherpa 链接）
CGO_ENABLED=1 go vet ./...             → VET_OK
CGO_ENABLED=0 go build ./internal/memory/... ./internal/observe/... ./internal/plugin/...
                                       → OK（memory 链纯 Go 无 cgo）
CGO_ENABLED=0 go test ./...            → buildinfo/observe/plugin/proc/memory 全 ok；
                                         唯一 FAIL=cmd/wisp [setup failed]（CGO=0 下 sherpa
                                         build-constraint 排除 —— 环境固有约束，非本票缺陷，
                                         CGO=1 下整仓 build 绿已补证）
CGO_ENABLED=0 go test -count=1 -v ./internal/memory/...
                                       → 21 个顶层测试函数 = 20 PASS + 1 SKIP
                                         （TestSubprocessCrashWriter 为崩溃测试子进程角色，
                                         SKIP 属预期），suite 11.3s —— 与实现声称一致
CGO_ENABLED=1 go test -race -count=1 ./internal/memory/...
                                       → **FAIL：8 × "WARNING: DATA RACE"，--- FAIL:
                                         TestRetentionSchedule**（复跑 2/2 稳定复现）
```

- 测试数对账：实现日志称「27 tests」；实测 21 个顶层函数 + 16 个 introspection 子测试，
  顶层口径 20 绿 + 1 SKIP。计数口径不同但「全绿（除子进程角色 SKIP）」成立。
- **MAJOR-2**：`go test -race ./internal/memory/...` 并非全绿，见下文专节。

### 2. DDL 对账独立复核 — PASS（字节级一致属实）

用 Temp 目录独立程序（自写提取器 + modernc.org/sqlite v1.59.0 直驱，不经仓库代码）三方对账：

1. **文本层**：从 `docs/specs/SPEC-02-data-storage.md` §3 提取 ```sql 块，从 `internal/memory/schema.go`
   提取 `ddlV1` 字面量 —— **去首尾空白后 byte-identical = true**（含注释、含中文注释）。
2. **规范化层**：剥离 `--` 注释、折叠空白后 token 序列 identical。
3. **语义层**：两份脚本各自灌入两个全新 SQLite 库 ——
   - `sqlite_master` 各 15 个对象（8 表 + 7 索引），逐对象规范化 SQL **完全一致**（能抓住索引列/表达式漂移）；
   - 8 张表 `PRAGMA table_info` 全部一致（cid 顺序、类型、notnull、pk）。
4. **introspection 契约测试核查**（schema_test.go）：确实断言列**序**（按 cid 逐位比对 name/type/
   notnull/pk）、索引名 + 数量 + 归属表、schema_version=1、splitSQLStatements 恰 15 语句。
   其「PK 列 notnull=0」的契约豁免与实测 ground truth 完全吻合（TEXT PK 与 INTEGER PK 均如此）。

### 3. PRAGMA 实证（实查非读 DSN）— PASS

以与 `open.go dsn()` 完全相同的 DSN 开库后**实查**：

```
PRAGMA journal_mode       = wal      （writer 连接与 reader 连接均如此）
PRAGMA synchronous        = 1        （NORMAL）
PRAGMA busy_timeout       = 3000
PRAGMA wal_autocheckpoint = 1000
PRAGMA page_size          = 4096
```

- **busy_timeout 是活的**：A 连接持写锁时 B 连接写入实测阻塞 **3.017s** 后才返回 SQLITE_BUSY。
- **`_txlock=immediate` 是活的（行为学证明）**：无 `_txlock` 连接在活动写者下 `Begin()` 立即返回；
  带 `_txlock=immediate` 的连接 `Begin()` 实测阻塞 **3.023s** 后报 SQLITE_BUSY —— 即 BEGIN IMMEDIATE
  真实在 Begin 时抢 RESERVED 锁。此实证必要：modernc 驱动对**未知** `_param` 静默忽略（实测
  `_definitely_unknown_param(1)` 打开不报错），仅凭 DSN 字符串无法证明。
- 写池 `MaxOpenConns=1`、读池 8 连接、两池 `ConnMaxIdleTime=30s`（D20 按需开关）均在 open.go 核实。

### 4. 单写者审计 — PASS（附 MAJOR-1 缺陷）

- **绕过点扫描**：`grep` 全包 `s.writer.` 直接使用仅 2 处 —— `startupCheckpoint()` 与 `Close()`
  的 `PRAGMA wal_checkpoint(TRUNCATE)`（PRAGMA 非数据写，且均在无并发窗口）；迁移 `boot` 连接
  在池创建前单线程使用。全部 8 表 DAO 写 + privacy 删除/清空 + retention 清理 + SearchMemory
  命中统计，共 20 余个写路径**全部经 `s.write → writeQueue`**。其他包只能触达公开 API
  （`s.writer` 未导出），无外部绕过面。
- **空闲即退协议审查**（writer.go）：`enqueue` 与 `loop` 的退出决策持同一把 `q.mu` ——
  入队要么发生在 writer 最后一次空查之前（被排空），要么看到 `running=false` 并同步 Spawn 新
  writer；`Spawn` 在持 `q.mu` 下调用但只短暂取 `observe` 内部锁、无反向依赖，无死锁；
  `cmd.done` 带 1 缓冲，caller 弃等也不阻塞 writer。**非 panic 路径无丢失写**，结论成立。
- **并发复跑**：2 写 + 4 读 × 10s = **9,707 写 / 14,266 读循环 / 0 busy**，WAL 4,210,672B
  ≤ 2×1000 页 × 4096B（autocheckpoint 有界）；16 并发突发写全成功后 queue 空闲、db-writer 归零
  （TestWriteQueueLazyLifecycle —— 按需启动 + 空闲即退实证，RosterReport ≤ 基线 6 由「至多 1 个
  db-writer 实例 + 空闲退出」构造性成立，花名册本身零改动）。
- **MAJOR-1（panic 楔死）**：见下文专节。

### 5. 崩溃恢复独立复现（自建实验，非复用实现跑法）— PASS

经 `go test -overlay` 注入自写测试（仓库磁盘零改动）：

- **自写写循环**：经真实 `Store.write` 队列持续提交双语句事务（`task_log`+`tool_call`，
  4KB 随机 payload，拉长事务窗口）；
- **杀手**：Windows 原生 `windows.TerminateProcess`（非 `os.Process.Kill` 复用），在
  第 N 次 commit 公告后 **+2ms / +7ms 随机化**打入下一事务中段（kill 时 WAL 实测 2.0~2.2MB，
  证明 kill 落在写中途）；
- **独立验证器**（Temp 独立程序，纯 modernc 直驱，**零仓库代码**）读盘：

| 轮次 | 公告 commit | task_log/tool_call | integrity_check | WAL checkpoint(TRUNCATE) 后 |
|---|---|---|---|---|
| 1（+2ms） | 50 | 50 / 50（=公告数，零丢失，双语句事务原子） | ok | 2,117,712B → **0B**，integrity 复检 ok |
| 2（+7ms） | 53 | 53 / 53 | ok | 2,241,336B → **0B**，ok |

- 全部行 payload 完整（`length(query_text)=4096`）、无孤儿 task 行。WAL/断电 reopen 语义
  与实现代理的子进程测试结论一致，本轮为独立方法学复核。

### 6. 迁移安全独立实验 — PASS

- **版本过新（v99）**：伪造 `schema_version=99` 并 checkpoint 落主文件、删 -wal/-shm 后 ——
  `Open` 失败，`*SchemaError{CurrentVersion:99, Target:1, "downgrading is not supported"}`，
  `errors.Is(err, ErrSchemaUnmigratable)=true`；**主文件 sha256 前后逐字节一致**
  （`dc9dd3e8…`）；backup 目录清单零新增（拒绝迁移不产生任何备份/写副作用）。
- **备份不覆盖**：构造失败迁移（事务内 NOT NULL 违例）→ 首次失败写 `bak-1-2`；
  **在两次失败之间向 v1 库加新数据**，二次失败后 `bak-1-2` sha256 逐字节不变
  （`8ec3fc2d…`），内容仍为第一快照（1 行 profile，未混入新数据）——
  保留最老快照的语义有日志佐证（`migration backup already exists, keeping it`）。
- 水位同行更新 + 单事务：`applyMigrationStep` DDL 与 `schema_meta` upsert 同一 `*sql.Tx`，
  TestMigrationFailedStepIsAtomic 佐证回滚。

### 7. 保留期与隐私 API 复跑 + LRU 日志断言 — PASS

- `TestRetentionBoundaries`（29/30/31 天、399/400/401 天、grant 失效 30 天审计窗、严格大于
  语义）PASS；`TestPrivacyOpsAllDomains`（五域 × list/delete-one/purge-all/export-JSON +
  未知域拒绝 + 空导出为 `[]`）PASS。
- **LRU 日志断言为真**：`TestProfileUpsertAndLRUEvictionLogs` 用 bufLogger 捕获**真实淘汰**
  （改造 updated_at 制造确定性 LRU 序），断言 `profile LRU eviction` + `evicted_slot=pref.k00`
  出现在捕获文本 —— 断的是 logger 输出，非打印常量；artifacts LRU 用 `os.Chtimes` 造龄真实
  删除最老文件后断言 `artifacts LRU cleanup`。均非假断言。
- 调度契约：`defaultFirstDelay=5m`、`defaultPeriod=24h`（retention.go 常量核实），单调
  Timer/Ticker；`StartRetentionJob` 无 scope 即 panic（强制 C11），TestRetentionSchedule 证明
  Dispose 后不再有 pass（该测试本身有 MAJOR-2 的数据竞争）。

### 8. D22 依赖白名单 / commit 卫生 / 越界扫描 — PASS（附 MINOR-1）

- **go.sum 增量全账**（3778c96..HEAD）：恰为 `modernc.org/sqlite v1.59.0` 闭包
  （libc/mathutil/memory + humanize/uuid/go-isatty/strftime/bigfft 传递）+
  **x/sys 0.42.0→0.47.0**。`go mod graph` 证实 `modernc.org/libc@v1.75.7` 与
  `modernc.org/sqlite@v1.59.0` 均要求 `x/sys@v0.47.0` —— 「连带升级」自述属实。
- **白名单符合**：SPEC-01 L96-101 D22 白名单明列 `golang.org/x/sys` 与 `SQLite 驱动`，且
  【SPEC】点名 `modernc.org/sqlite`；传递闭包属白名单内选型的必然组成，无白名单外新依赖。
- **6 个 commit 逐一 `git show --stat`**：全部显式路径，仅触碰 `internal/memory/**`、
  自身票据文件 `.scratch/wisp/issues/04-sqlite-core.md`、`go.mod/go.sum`；
  **零 spike 文件卷入**（spike 路径的全部变更均属 T02 的 8 个 commit，已逐条归属核对）。
- **零 emoji**：`internal/memory/*.go` + 票据文件扫描真实 emoji 码位（1F000-1FAFF/FE0F）零命中
  （注释中的 `→`/`≤` 为普通排版符号，与全仓风格一致）。
- **`provider_health` 零出现**：grep 全包仅 schema.go 头注释一句「schema v2 belongs to ticket 09 —
  deliberately absent here」—— 防蔓延正确，无表、无 DDL、无代码引用。

### 9. 票据对照 — PASS（附 MINOR-2）

| 票面验收标准 | 裁决 | 证据 |
|---|---|---|
| 每表 sqlite_master 契约测试 = SPEC-02 §3 | PASS | 清单 2（双库对账 + 列序/类型/notnull/pk 断言） |
| 2 写 + 4 读 × 10s 零 busy、WAL 有界 | PASS | 清单 4（9,707/14,266/0 busy，4.21MB ≤ 预算） |
| kill 写中 → reopen integrity_check + checkpoint 截断 | PASS | 清单 5（两轮独立 TerminateProcess 实验） |
| 保留期边界 29/30/31、399/400/401 | PASS | 清单 7 |
| 隐私 API purge/export/delete-one × 五域 | PASS | 清单 7（含未知域拒绝、ErrNotFound 语义） |
| profile LRU 淘汰写日志（测试断言日志行） | PASS | 清单 7（真实缓冲捕获断言） |

- 票面之外声称核实：WAL 四 PRAGMA + 启动/收尾 `wal_checkpoint(TRUNCATE)`（清单 3）、迁移链
  单事务/水位同行/事前备份/不可迁移字节不变（清单 6）、db-writer 按需 + 空闲退 + Roster ≤ 6
  （清单 4）、五域隐私 + artifacts 500MB LRU（清单 7）—— 均成立。
- Progress log：8 条记录覆盖 claim→handoff 全程，**内容完整**；但 `[08:54:30]`（并发+崩溃条目）
  列于 `[08:54:03]`（handoff）之前，违反「append-only, newest last」（MINOR-2）。
- Status：`in-progress` ✓（等待 orchestrator 处置本报告后变更，符合流程）。

## 问题清单（按严重度）

### MAJOR-1 · 写队列 panic 楔死（writer.go:94-109 `writeQueue.loop`）— 实验证实

队列命令的 fn（`cmd.fn`）panic 时：observe 边界按 D37b 恢复（进程存活、internal 错误落日志），
但 `loop` 是被 panic 打断的 —— `q.running` **永远停留 true**，且 panic 命令及其后所有已入队命令
的 `done` 永不投递。此后**任何** `Store.write` 入队后看到 `running=true` 不再 Spawn 新 writer，
永久阻塞至 caller ctx 超时。store 的写路径自此刻整体失效直至进程重启，错误表现为
「write abandoned by caller: context deadline exceeded」（误导性）。

实证（overlay 注入测试 `TestAdvPanicWedge`）：注入 `panic("boom")` 的写命令 300ms 超时返回后，
2 秒超时的正常写**同样超时失败**（预期应为成功）——楔死复现。

现状无人能 panic（现有 DAO 全为参数校验 + SQL Exec），故非 BLOCKER；但 `Store.write` 是票 29+
agent 回路复用的公开写入口，未来任何一个闭包 bug 即把存储核心整体打瘫，且违背 writer.go 自述
「No lost writes」。**最小修复**：`loop` 内对 `q.exec` 包 recover（或 `defer` 中重置 `q.running`
并向当前命令投递 internal 错误），保证 panic 后队列可重建；补一条「panic 后队列恢复」回归测试。

### MAJOR-2 · `go test -race ./internal/memory/...` 红（retention_test.go 测试内竞争）

钉死工具链下 2/2 稳定复现：8 个 DATA RACE，唯一失败测试 `TestRetentionSchedule`。
根因在**测试挂具**而非生产代码：`bufLogger` 的共享 `*bytes.Buffer` 被测试主 goroutine 在轮询中
读取（retention_test.go:176 `strings.Count(buf.String(), …)`），同时 retention-job goroutine 经
slog 向同一 buffer 写日志 —— `bytes.Buffer` 非并发安全。生产路径（Store/queue/DAO/retention/
open）在本轮 -race 下**零竞争报告**。

影响：SPEC-10 的 race 门禁将永久红；实现代理 handoff 未直接声称 `-race`（票据日志只称
CGO_ENABLED=0-ok），但协调方背景陈述中的「-race 全绿」与事实不符。**最小修复**：bufLogger 的
writer 加 mutex（或换成带锁 Handler / 原子计数器），-race 即绿；预计 ≤10 行且不影响断言语义。

### MINOR

1. **go.mod 不 tidy**：`modernc.org/sqlite v1.59.0` 停在 indirect 块且带过期 `// indirect` 标记，
   而 internal/memory 直接 import（`go mod tidy -diff` 显示应上移直接依赖块）。纯卫生问题。
2. **Progress log 时序**：`[08:54:30]` 条目列于 `[08:54:03]`（handoff）之前，违反 append-only
   newest-last；handoff 时间戳早于其声称已完成的测试条目。
3. **`writeQueue.last` 死字段**：仅赋值（writer.go:81）从不读取。
4. **retention 边缘**：`tool_call` 行 `started_at`/`decided_at` 均为 NULL 时永不过期
   （deleteOlderThan 的 `IS NOT NULL` 守卫）—— 当前 DAO 允许插入双 NULL 的 pending 行。
5. **`parseVersion` 前缀解析**：`Sscanf("%d")` 使 `"1abc"` 被接受为版本 1 —— 被篡改/损坏水位
   未按「不可迁移」明确失败。建议 strconv.Atoi + 全串校验。

## 裁决

九项验收中八项 PASS 且证据独立复核（DDL 字节级、PRAGMA 行为学实证、双轮独立崩溃实验、
迁移字节不变/备份不覆盖、白名单与 commit 卫生全对账）；唯一未过的硬性项是清单 1 明确要求的
`go test -race ./internal/memory/...`（MAJOR-2），另证实一项写路径健壮性缺陷（MAJOR-1）。
两项 MAJOR 修复集均极小（合计约 20 行 + 2 个回归测试），不触及 DDL 契约与架构。

**VERDICT: FAIL (最小修复集)**
1. `writer.go`：loop 内 recover / 重置 `q.running`，panic 后队列可重建 + 回归测试。
2. `retention_test.go`：bufLogger 加锁，使 `-race` 绿 + 保留断言语义。
3. `go mod tidy`（modernc.org/sqlite 移入直接依赖块）。
4. （顺手）Progress log 时序勘误备注；`writeQueue.last` 删除或启用。
