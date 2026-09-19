# SPEC-02 · 数据存储

> 追溯：D20/D26/D35/§14.4/§14.8/§14.11/§14.7、C13/C23/C28；环境隔离的数据目录分叉见 SPEC-03 §5。

## 1. 背景与问题

D35 给出了八张表的骨架，但缺字段类型、迁移机制、清理实现与并发细节；这些不定，agent 会自己发明
表结构。存储是隐私的核心载体（画像/记忆/任务日志/取证），一切保留、可见、可删要求必须落到 schema。

## 2. 存储选型与连接策略

- 【PLAN】SQLite，`%APPDATA%\wisp\wisp.db`；连接按需开关，不常驻（D20）。
- 【SPEC】驱动 `modernc.org/sqlite`（纯 Go，无 cgo——cgo 预算全部留给 sherpa-onnx，交叉编译也随之简化）。
- 【PLAN】WAL 与并发（D35 规则 3）：`PRAGMA journal_mode=WAL` · `synchronous=NORMAL` ·
  `busy_timeout=3000` · `wal_autocheckpoint=1000` · **启动时 `wal_checkpoint(TRUNCATE)`**
  （防 WAL 无界增长）。
- 【PLAN】**单写者模式**：所有写经 `db-writer` goroutine 串行化（命令 channel），避免
  `SQLITE_BUSY` 重试逻辑散落各处；读可并发（WAL 读不阻塞）。
- 【PLAN】连接按需开关与 WAL 的配合需验证——WAL 依赖持久 `-wal`/`-shm` 文件，测试必须覆盖
  「进程正常退出 / kill / 断电三种退出后重开」。

## 3. Schema（契约级，不得自行增删字段；【SPEC 细化】给出类型）

```sql
-- schema 版本与迁移水位（§14.5 的地基）
CREATE TABLE schema_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
-- 内置行：('schema_version', '1')

-- L1 用户画像（D20）：slot 有限枚举，≤20 行由应用层强制
CREATE TABLE profile (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  slot       TEXT NOT NULL UNIQUE,      -- 'pref.language'|'pref.tone'|'habit.work_hours'|'fact.family'|…
  content    TEXT NOT NULL,             -- 一句人话
  source     TEXT NOT NULL,             -- 'extracted' | 'manual'
  updated_at INTEGER NOT NULL           -- unix 秒（落盘时间戳一律 wall clock，D42#9）
);
CREATE INDEX idx_profile_updated ON profile(updated_at);  -- LRU 淘汰依据

-- L2 显式记忆（D20）
CREATE TABLE memory (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  kind        TEXT NOT NULL,             -- 'explicit'（预留 kind 列位，不新增实现）
  content     TEXT NOT NULL,
  keywords    TEXT NOT NULL,             -- 空格分隔、小写规范化的检索词
  created_at  INTEGER NOT NULL,
  last_hit_at INTEGER NOT NULL DEFAULT 0,
  hit_count   INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_memory_keywords ON memory(keywords);
CREATE INDEX idx_memory_last_hit ON memory(last_hit_at);

-- L3 任务日志（D20；query_text 按 §14.4 脱敏）
CREATE TABLE task_log (
  id             TEXT PRIMARY KEY,       -- 任务 UUID
  started_at     INTEGER NOT NULL,
  ended_at       INTEGER,
  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
  query_text     TEXT NOT NULL,          -- 已脱敏；模式匹配掩码是尽力而为，须标注（§14.4）
  summary_text   TEXT,
  cost_tokens_in  INTEGER NOT NULL DEFAULT 0,
  cost_tokens_out INTEGER NOT NULL DEFAULT 0,
  cost_amount_micro INTEGER NOT NULL DEFAULT 0,  -- 【SPEC】1e-6 货币单位整数，避免浮点
  currency       TEXT NOT NULL DEFAULT 'CNY',    -- 【SPEC】
  error_class    TEXT                             -- D37 枚举，可空
);
CREATE INDEX idx_task_log_started ON task_log(started_at);

-- 工具调用与审批取证（owner 隐私敏感：「它到底对我的机器做了什么」必须可查）
CREATE TABLE tool_call (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id        TEXT NOT NULL REFERENCES task_log(id),
  seq            INTEGER NOT NULL,
  tool           TEXT NOT NULL,
  args_json      TEXT NOT NULL,          -- 长字符串截断（§5.1 脱敏规则）
  risk_level     TEXT NOT NULL,          -- 'L0'|'L1'|'L2'
  decision       TEXT,                   -- 'allow'|'allow_session_grant'|'reject'|'timeout'|'batch_aggregated'
  decided_at     INTEGER,
  started_at     INTEGER,
  ended_at       INTEGER,
  outcome        TEXT,                   -- 'success'|'error'|'cancelled'|'truncated'
  error_class    TEXT,
  correlation_id TEXT NOT NULL,          -- C18 路由与取证
  grant_id       INTEGER                 -- 关联 approval_grant（D45）
);
CREATE INDEX idx_tool_call_task ON tool_call(task_id, seq);
CREATE INDEX idx_tool_call_corr ON tool_call(correlation_id);

-- D45 作用域会话授权
CREATE TABLE approval_grant (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  scope      TEXT NOT NULL,              -- 'session'
  tool       TEXT NOT NULL,
  pattern    TEXT NOT NULL,              -- 路径模式 / 目标进程名（input.type）
  session_id TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,           -- 会话结束时间（失效后仅作审计）
  revoked_at INTEGER
);
CREATE INDEX idx_grant_session ON approval_grant(session_id);

-- C23 日聚合（避免每次全表扫 task_log）
CREATE TABLE cost_daily (
  day        TEXT PRIMARY KEY,           -- 'YYYY-MM-DD'
  tokens_in  INTEGER NOT NULL DEFAULT 0,
  tokens_out INTEGER NOT NULL DEFAULT 0,
  cost_micros INTEGER NOT NULL DEFAULT 0,
  tasks      INTEGER NOT NULL DEFAULT 0
);

-- 插件安装状态 + 完整性校验（§14.7）
CREATE TABLE plugin_state (
  id                 TEXT PRIMARY KEY,   -- 插件 id
  version            TEXT NOT NULL,
  enabled            INTEGER NOT NULL DEFAULT 1,
  capabilities_json  TEXT NOT NULL,
  net_allowlist_json TEXT,
  installed_at       INTEGER NOT NULL,
  hash               TEXT NOT NULL,      -- manifest 哈希，加载时不符 → 拒绝加载并告警
  exe_hash           TEXT                -- 【SPEC】D46 command 插件的目标 exe sha256（安装时钉死）
);
```

- **`provider_health`（schema v2 追加，2026-09-19 用户批准——LLM 接入补充，SPEC-03 §3.1）**：
  `provider` PK · `model` PK（联合主键）· `probe{text,vision,audio_in,audio_out,thinking,fc}`（JSON）
  · `last_probe_at` · `last_error` · `last_error_at` · `latency_ms_p50` · `quota_state`。
  存**运行观测态**（探测结果/健康/最近错误），与 config.toml 的目录配置（真相源）严格分界；
  由票 09 定义原语、票 11 实测写入、票 40 面板展示。
  schema v2 迁移 DDL（与 `internal/memory/schema.go` 的 `ddlV2` 逐字节一致；票 09 已落地）：

```sql
-- provider 健康与实测能力（schema v2：LLM 接入补充，SPEC-03 §3.1；票 09 定义原语，票 11 实测写入）
-- 存「运行观测态」（探测/健康/最近错误），与 config.toml 的目录真相源严格分界
CREATE TABLE provider_health (
  provider       TEXT NOT NULL,            -- 目录 provider 名
  model          TEXT NOT NULL,            -- 目录 model id
  probe_json     TEXT NOT NULL DEFAULT '{}', -- 实测能力位 JSON：{"text":..,"vision":..,"audio_in":..,"audio_out":..,"thinking":..,"fc":..}
  last_probe_at  INTEGER,                  -- unix 秒；NULL = 尚未实测
  last_error     TEXT,                     -- 最近一次失败（脱敏后细节）
  last_error_at  INTEGER,
  latency_ms_p50 INTEGER NOT NULL DEFAULT 0,
  quota_state    TEXT NOT NULL DEFAULT 'ok', -- 'ok' | 'throttled' | 'exhausted' | 'unknown'
  PRIMARY KEY (provider, model)
);
```
- L2 记忆检索：【SPEC】显式记忆量级小（几十条），用 `keywords` LIKE + 最近优先排序即可；
  超过 500 条再引入 FTS5（届时登记 schema 迁移 v2）。
- 时间戳规则：**落盘时间戳用 wall clock；一切超时/保留期计算用单调时钟**（D42#9）。

## 4. 保留期与清理（有归属才有效）

| 数据 | 保留期 | 清理者 |
|---|---|---|
| `task_log` / `tool_call` | 30 天 | `memory.RetentionJob` |
| `cost_daily` | 400 天 | 同上 |
| `approval_grant` 失效行 | 保留 30 天供审计 | 同上 |
| `profile` / `memory` | 常驻，用户逐条删 / L1 满 20 条按 `updated_at` LRU 淘汰**并写日志**（淘汰必须可见） | 应用层 |
| `artifacts\` 目录 | 500MB 配额，超则 LRU 清理 | 同上 |

`RetentionJob` 调度（D35 规则 1）：**进程启动后延迟 5 分钟跑一次 + 每 24h 一次**——不能只靠周期
定时器（常驻进程可能只开 10 分钟就关机）。挂在 `retention-job` 临时 goroutine。

## 5. 面板「数据与隐私」页（D35 规则 2，非可选）

对 `profile` / `memory` / `task_log` / `tool_call` / `artifacts` 全部提供：
**一键清空 / 导出（JSON）/ 逐条删除**。ASR 转写历史全量可见、可逐条删、可一键清空（§14.4）。

## 6. 数据目录布局（§14.8/D41(a)；环境分叉见 SPEC-03 §5）

```
%APPDATA%\wisp\
├── config.toml              # A 档黑名单（F1）：Agent 不可读写
├── secrets\                 # C28 DPAPI blob（per-ref 一个文件）
├── models\                  # 可经 [models] dir 单独指定（D26 逃生门）
│   └── <model-id>\…
├── plugins\                 # 已安装插件（Tier1 清单 + Tier2 js）
│   └── <plugin-id>\…
├── artifacts\               # 宿主内部工件（D10 落盘/D15 spill）：tool-output-<id>.txt 等
├── logs\                    # JSONL，按大小+按天滚动，默认保留 7 天
├── wisp.db  (+ -wal/-shm)
├── staging\  backup\<version>\  update-pending.json   # D41(b) 更新
└── portable.txt             # 存在 → 便携模式：以上全部落在 exe 同级 data\（§14.8）
```

- 模型目录单独可指定：模型几百 MB，用户可能想放非系统盘（D26）。
- **便携模式 × DPAPI（P13，【OPEN】）**：DPAPI blob 换机器/换用户无法解密。【SPEC 建议】便携模式下
  `secret` 仅支持 `env:` 引用（不支持 DPAPI 持久 blob），首次尝试 `dpapi:` 引用解密失败时给明确
  错误并引导改为 `env:`；**不得退回明文**。S1 定案。
- 卸载：默认不删 `models\` 与 `wisp.db`，必须明确告知占用大小并提供「一并删除」选项（D41(c)）。
- 备份与迁移：`config.toml` + `wisp.db` 一键导出/导入（换机）；导出**不含** `secrets\`
  （DPAPI blob 跨机无效，导出的 config 里 `api_key_ref` 保留引用名，导入后需重录）。

## 7. schema 迁移机制（§14.5）

- `schema_meta.schema_version` 单调递增；迁移代码按版本链组织，每步在**单事务**内完成并记录
  起止时间；迁移前自动 `backup\wisp.db.bak-<from>-<to>`。
- 无法迁移 → 启动失败进 `Unconfigured` 态并给出明确指引，**不得静默用默认值覆盖**
  （会静默重置目录白名单与权限设置——属安全问题）。
- `config.toml` 的 schema 迁移独立于 DB（见 SPEC-03 §4.4）。

## 8. 测试决策

- 每张表一组契约测试：字段/索引存在性（用 `sqlite_master` 自省）+ 保留期清理边界（29/30/31 天）。
- 并发测试：2 写者 + 4 读者并发 10s，断言零 `SQLITE_BUSY` 泄漏、WAL 文件 ≤ autocheckpoint 预期。
- 崩溃恢复：写一半 kill 进程 → 重开库校验完整性（`PRAGMA integrity_check`）与 checkpoint 生效。
- 隐私测试：写入含「密码/卡号」样式的转写文本后，`task_log.query_text` 已按 §14.4 掩码（尽力而为
  标注仍保留在 UI，不因测试而虚报安全性）。

## 9. 不做什么

- 不做 ORM / 迁移框架依赖（手写版本链，依赖白名单内不含有此能力者）。
- 不做跨进程共享 DB（单进程单实例语义由 per-session 互斥保证，SPEC-09 §5）。
- 不做加密整库（SQLCipher 类）；敏感项由 C28 逐项加密 + A 档黑名单 + 转写脱敏兜住，
  全库加密的成本与 D1 相左，且用户可用 BitLocker。
