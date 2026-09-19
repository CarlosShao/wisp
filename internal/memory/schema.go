package memory

import (
	"context"
	"database/sql"
	"strings"
)

// Schema v1 (D35 / SPEC-02 §3). The DDL below is CONTRACT-LEVEL: it must stay
// byte-identical to the ```sql block of docs/specs/SPEC-02-data-storage.md §3
// (including comments). No added or removed columns, tables, indexes or
// constraints — the acceptance test introspects sqlite_master against the
// same contract. Enum-shaped columns (task_log.error_class, tool_call
// risk_level/decision/outcome) are validated at the DAO layer
// (observe.ValidateErrorClass and friends), NOT as SQL CHECKs, precisely so
// this text never forks from the spec.
//
// provider_health (SPEC-02 §3 note) is schema v2 and belongs to ticket 09 —
// deliberately absent here.
const ddlV1 = `-- schema 版本与迁移水位（§14.5 的地基）
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
);`

// SchemaVersionTarget is the schema version this binary produces (SPEC-02 §3:
// builtin row ('schema_version', '1')).
const SchemaVersionTarget = 1

// schemaMetaVersionKey is the schema_meta row holding the migration watermark.
const schemaMetaVersionKey = "schema_version"

// migrationStep is one link of the hand-written version chain (SPEC-02 §7:
// no migration framework dependency). Every step runs inside a SINGLE
// transaction, after a pre-backup backup\wisp.db.bak-<from>-<to> has been
// written; the schema_version row is updated inside the same transaction.
type migrationStep struct {
	from  int
	to    int
	apply func(tx txExec) error
}

// txExec is the subset of *sql.Tx the steps need (keeps the chain testable).
type txExec interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// migrationChain is the ordered version chain. The base step creates the
// D35 schema from an empty database; later steps append (v1->v2 ...).
var migrationChain = []migrationStep{
	{from: 0, to: 1, apply: applyV1},
}

// applyV1 creates the eight D35 tables + indexes and seeds schema_version=1.
func applyV1(tx txExec) error {
	for _, stmt := range splitSQLStatements(ddlV1) {
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return err
		}
	}
	return nil
}

// splitSQLStatements splits a script on ';' terminators. Safe for the D35 DDL:
// no string literal or comment in it contains a semicolon (asserted by
// TestDDLContractIntrospection, which round-trips the script through SQLite).
func splitSQLStatements(script string) []string {
	parts := strings.Split(script, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		// Drop pure-comment chunks (statement-leading comments are kept and
		// understood by SQLite itself).
		body := strings.Map(func(r rune) rune {
			switch r {
			case ' ', '\t', '\r', '\n':
				return -1
			}
			return r
		}, strings.TrimSpace(strings.TrimPrefix(p, "--")))
		if body == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}
