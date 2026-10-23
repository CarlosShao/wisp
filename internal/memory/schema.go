package memory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Schema v1 (D35 / SPEC-02 §3) + schema v2 (provider_health, SPEC-02 §3
// note). ddlV1/ddlV2/ddlV3 are CONTRACT-LEVEL: they must stay byte-identical to
// the ```sql blocks of docs/specs/SPEC-02-data-storage.md §3 (including
// comments). No added or removed columns, tables, indexes or constraints —
// the acceptance test introspects sqlite_master against the same contract.
// Enum-shaped columns (task_log.error_class, tool_call
// risk_level/decision/outcome) are validated at the DAO layer
// (observe.ValidateErrorClass and friends), NOT as SQL CHECKs, precisely so
// this text never forks from the spec.
//
// ⚠ ONE DELIBERATE EXCEPTION, since schema v4 (ticket 196): task_log.state now
// carries a SQL CHECK. It is added by ddlV4 — a REBUILD step, not an edit of
// ddlV1 — so ddlV1's text and statement count stay exactly what SPEC-02 §3
// declares, and a fresh database reaches the constrained shape through the same
// 3->4 step an existing one does (the rule ddlV3 states below). The owner
// approved touching D35 for this (ticket 196's shape ruling, owner's words
// quoted in docs/evidence/s1/196-r1-impl.md), and SPEC-02 §3's text still has
// no CHECK: that is a spec-note debt recorded for the owner, not a drift this
// file is allowed to paper over by editing the spec.
//
// provider_health (SPEC-02 §3 note) is schema v2 and landed with ticket 09:
// the DAO primitive lives in dao_providerhealth.go; ticket 11 writes real
// probe results.
const ddlV1 = `-- schema 版本与迁移水位（§14.5 的地基）
CREATE TABLE schema_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
-- 内置行：('schema_version', '1')

-- L1 用户画像（D20）：slot 有限枚举，<=20 行由应用层强制
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
  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）
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
  hash               TEXT NOT NULL,      -- manifest 哈希，加载时不符 -> 拒绝加载并告警
  exe_hash           TEXT                -- 【SPEC】D46 command 插件的目标 exe sha256（安装时钉死）
);`

// SchemaVersionTarget is the schema version this binary produces (SPEC-02 §3:
// builtin row ('schema_version', '<target>')). History:
//
//	v1 - D35 core tables (ticket 04).
//	v2 - provider_health (2026-09-19 LLM access supplement, SPEC-03 §3.1 /
//	     SPEC-02 §3 note; ticket 09): per provider+model runtime observation
//	     state (probe results, last error, latency, quota state).
//	v3 - session management (ticket 166): the cross-restart conversation
//	     entity (`session`) plus its turns (`session_turn`). Added as its own
//	     migration step so NO text of ddlV1/ddlV2 changes: new databases reach
//	     v3 through the same 2->3 step old databases do, so the two paths
//	     cannot fork on which side creates which table (ticket 166 AC#3).
//	v4 - task_log.state gets an explicit CHECK (ticket 196). Same rule again:
//	     ddlV1 is untouched and every database, fresh or existing, reaches the
//	     constrained shape through the 3->4 rebuild in ddlV4. The accepted
//	     values come from taskStateRoster, so the SQL and the DAO cannot fork
//	     either.
//
// ⚠ D35 is a frozen contract (AGENTS.md §1.1) and SPEC-02 §3 still lists nine
// tables; this v3 roster extension carries a spec-note debt that
// ticket 166's progress note names for the owner (per AGENTS.md §0 "改契约＝
// 人工批准" and probes/166/a2 §H). v4 adds a second debt of the same kind:
// SPEC-02 §3's task_log block carries no CHECK, and the owner's approval to
// touch D35 does not edit the spec text. Both are named in
// docs/evidence/s1/196-r1-impl.md rather than fixed by an agent.
const SchemaVersionTarget = 4

// ddlV2 appends the provider_health table (SPEC-02 §3 note, schema v2). The
// DDL below mirrors the ```sql block in docs/specs/SPEC-02-data-storage.md
// §3 byte-for-byte (same contract rule as ddlV1).
const ddlV2 = `-- provider 健康与实测能力（schema v2：LLM 接入补充，SPEC-03 §3.1；票 09 定义原语，票 11 实测写入）
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
);`

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
// QueryContext joined the set for ticket 196's v4 step, which has to read the
// offending values out of task_log.state BEFORE it rebuilds the table: a raw
// CHECK failure names the constraint but never the row, and an operator staring
// at "CHECK constraint failed" with no value in it has nothing to fix.
type txExec interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// migrationChain is the ordered version chain. The base step creates the
// D35 schema from an empty database; later steps append (v1->v2 ...).
var migrationChain = []migrationStep{
	{from: 0, to: 1, apply: applyV1},
	{from: 1, to: 2, apply: applyV2},
	{from: 2, to: 3, apply: applyV3},
	{from: 3, to: 4, apply: applyV4},
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

// applyV2 creates the provider_health table (ticket 09 primitive).
func applyV2(tx txExec) error {
	for _, stmt := range splitSQLStatements(ddlV2) {
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return err
		}
	}
	return nil
}

// ddlV3 adds the session-management tables (ticket 166): the persisted
// conversation the ticket makes listable / searchable / renameable /
// archive-able / forkable, plus its turns (the "正文" the search reads).
//
// Why its own const and not an edit of ddlV1: a new database reaches v3 through
// the 2->3 step exactly like an existing one, so there is never a shape where a
// fresh table comes from CREATE and an upgraded one from ALTER - the divergence
// ticket 166 AC#3's second half exists to prevent.
//
// Naming: the entity carries SPEC-05 §5's word "会话", which there means the
// WITHIN-BOOT voice session (唤醒 → 90s 空闲). A cross-restart conversation
// sharing that word is an open naming question (probes/166/a2 §H1), so this DDL
// uses the table name `session` (the name the ticket and the a2 census both
// use) and leaves the vocabulary ruling to the owner rather than inventing a
// second term here.
const ddlV3 = `-- 会话管理（schema v3：票 166，跨重启的对话）
CREATE TABLE session (
  id          TEXT PRIMARY KEY,          -- 会话 id（crypto/rand 铸，非派生；与 C31 的 session.ID 同形）
  title       TEXT NOT NULL,             -- 用户可见标题；改名写这里，不动任何 id
  workspace   TEXT NOT NULL DEFAULT '',  -- 工作区标签；'' = 未归属
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL,
  archived_at INTEGER,                   -- NULL = 活动；归档 = 写时间戳；撤销 = 置回 NULL
  forked_from TEXT                       -- NULL = 原件；分叉出的新会话指回原件 id
);
CREATE INDEX idx_session_updated ON session(updated_at);
CREATE INDEX idx_session_workspace ON session(workspace, updated_at);

-- 会话轮次（"正文"）：搜索的正文面。改名/归档/分叉都不动这些行
-- ⇒ 原会话历史逐字仍在（票 166 AC#4）。
CREATE TABLE session_turn (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id TEXT NOT NULL REFERENCES session(id),
  seq        INTEGER NOT NULL,
  role       TEXT NOT NULL,              -- 'user' | 'assistant' | …
  content    TEXT NOT NULL,
  created_at INTEGER NOT NULL
);
CREATE INDEX idx_session_turn_session ON session_turn(session_id, seq);`

// applyV3 creates the session tables (ticket 166).
func applyV3(tx txExec) error {
	for _, stmt := range splitSQLStatements(ddlV3) {
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return err
		}
	}
	return nil
}

// ddlV4 rebuilds task_log with an explicit CHECK on state (ticket 196).
//
// Why a rebuild and not an edit of ddlV1: SQLite has no ALTER TABLE ADD
// CONSTRAINT, and ddlV1 is a frozen contract whose text and statement count the
// introspection test pins. The rebuild follows SQLite's documented 12-step
// procedure (create new -> copy rows -> drop old -> rename -> recreate indexes),
// all inside the one transaction the migration chain already wraps each step in.
//
// The accepted values are GENERATED from taskStateRoster, never spelled a second
// time here. That is the whole point of ticket 196's shape: the Go roster is the
// authority, and the CHECK is one of its two mouths (the DAO's validator is the
// other). Hand-copying the list into this string is exactly the fork the ticket
// exists to close.
var ddlV4 = buildDDLV4()

func buildDDLV4() string {
	quoted := make([]string, len(taskStateRoster))
	for i, ts := range taskStateRoster {
		quoted[i] = "'" + string(ts) + "'"
	}
	return `-- task_log 重建（schema v4：票 196，state 列加显式 CHECK）
-- 接受的值＝internal/memory/taskstate.go 的 taskStateRoster（Go 侧唯一权威），
-- 由 buildDDLV4 生成；这里不是第二份手抄名册。
-- 重建走 SQLite 官方 12 步 ALTER 过程：建新表 -> 搬行 -> 删旧表 -> 改名 -> 重建索引。
CREATE TABLE task_log_v4 (
  id             TEXT PRIMARY KEY,
  started_at     INTEGER NOT NULL,
  ended_at       INTEGER,
  state          TEXT NOT NULL CHECK (state IN (` + strings.Join(quoted, ", ") + `)),
  query_text     TEXT NOT NULL,
  summary_text   TEXT,
  cost_tokens_in  INTEGER NOT NULL DEFAULT 0,
  cost_tokens_out INTEGER NOT NULL DEFAULT 0,
  cost_amount_micro INTEGER NOT NULL DEFAULT 0,
  currency       TEXT NOT NULL DEFAULT 'CNY',
  error_class    TEXT
);
INSERT INTO task_log_v4 (id, started_at, ended_at, state, query_text, summary_text,
    cost_tokens_in, cost_tokens_out, cost_amount_micro, currency, error_class)
  SELECT id, started_at, ended_at, state, query_text, summary_text,
    cost_tokens_in, cost_tokens_out, cost_amount_micro, currency, error_class
  FROM task_log;
DROP TABLE task_log;
ALTER TABLE task_log_v4 RENAME TO task_log;
CREATE INDEX idx_task_log_started ON task_log(started_at);`
}

// taskLogStatesOutsideRoster lists the distinct task_log.state values the v4
// roster would reject. An empty result means every existing row survives the
// rebuild, which is what makes "an old database still opens" a measured fact
// rather than a hope.
func taskLogStatesOutsideRoster(ctx context.Context, tx txExec) ([]string, error) {
	names := TaskStateRoster()
	ph := make([]string, len(names))
	args := make([]any, len(names))
	for i, n := range names {
		ph[i] = "?"
		args[i] = n
	}
	rows, err := tx.QueryContext(ctx,
		`SELECT DISTINCT state FROM task_log WHERE state NOT IN (`+strings.Join(ph, ",")+`) ORDER BY state`,
		args...)
	if err != nil {
		return nil, fmt.Errorf("memory: census task_log.state for v4: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("memory: census task_log.state for v4: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("memory: census task_log.state for v4: %w", err)
	}
	return out, nil
}

// pragmaForeignKeysQuery reads the connection's foreign-key enforcement switch.
//
// It is READ ONLY on purpose. SQLite makes "PRAGMA foreign_keys=..." a no-op
// while a transaction is open, and the migration chain already runs every step
// inside one (open.go's applyMigrationStep). Measured on this driver: with the
// switch off, setting it ON inside the step's transaction leaves it at 0, while
// the same statement on the same connection with no transaction open flips it
// to 1. So applyV4 cannot disable enforcement for itself - it can only refuse
// to run under it, which is what requireForeignKeysOff does.
const pragmaForeignKeysQuery = "PRAGMA foreign_keys"

// requireForeignKeysOff is the guard ticket 196's r2 leg puts around the rebuild.
//
// The v4 rebuild drops task_log and recreates it, while tool_call.task_id
// carries a REFERENCES task_log(id) clause (ddlV1). Under enforcement ON,
// DROP TABLE runs an implicit DELETE FROM first, so with any tool_call row
// present it aborts on SQLite's "FOREIGN KEY constraint failed" - a message
// that names neither the table nor a row, and that arrives three statements
// into a destructive sequence, after the CREATE and the INSERT..SELECT have
// already run inside the same transaction. That is the failure this guard
// replaces.
//
// Refusing here moves the failure to the first statement of the step, before
// anything is written, and the message names the pragma, the child table and
// the fix. It does NOT make the rebuild work under enforcement ON - it cannot,
// see pragmaForeignKeysQuery. SQLite's documented 12-step ALTER procedure
// disables enforcement for the duration of a rebuild (its step 1); the migration chain
// takes that step for itself in migrationDSN (open.go), which drops every foreign_keys
// entry from the shared list and appends exactly one OFF. This guard is therefore a
// backstop rather than the load-bearing precondition: it fires if that declaration is
// ever lost.
//
// Cost, stated so nobody has to rediscover it: a database that would have
// survived by luck (zero tool_call rows) is refused too. That is deliberate - a
// rebuild whose success depends on how many child rows happen to exist is not a
// rebuild worth running.
//
// A driver that answers the PRAGMA with no row leaves the pre-r2 behaviour in
// place rather than refusing to open every database on earth; the read itself
// failing is reported, because that is a real defect rather than a dialect gap.
func requireForeignKeysOff(tx txExec) error {
	rows, err := tx.QueryContext(context.Background(), pragmaForeignKeysQuery)
	if err != nil {
		return fmt.Errorf("memory: read PRAGMA foreign_keys before the task_log rebuild: %w", err)
	}
	defer rows.Close()
	on := 0
	if rows.Next() {
		if err := rows.Scan(&on); err != nil {
			return fmt.Errorf("memory: read PRAGMA foreign_keys before the task_log rebuild: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("memory: read PRAGMA foreign_keys before the task_log rebuild: %w", err)
	}
	if on == 0 {
		return nil
	}
	return fmt.Errorf("memory: refusing to rebuild task_log with PRAGMA foreign_keys=ON: " +
		"tool_call.task_id REFERENCES task_log(id), so the DROP TABLE the v4 rebuild needs would " +
		`abort on SQLite's "FOREIGN KEY constraint failed" and name neither a table nor a row. ` +
		"SQLite's 12-step ALTER procedure disables enforcement for a rebuild (its step 1), and " +
		"the migration connection is supposed to declare that itself (migrationDSN in open.go): " +
		"something removed that declaration. Nothing was modified.")
}

// applyV4 rebuilds task_log with the state CHECK (ticket 196).
//
// It refuses BY NAME before touching anything: if a row already holds a word the
// roster does not accept, the step stops and says which words those are. The
// alternative - letting the INSERT..SELECT hit the CHECK - aborts the migration
// with "CHECK constraint failed" and no row identity, which is a dead end for
// whoever has to repair the database. Refusing here keeps the promise the
// ticket's second hard boundary makes: a database that satisfies today's
// vocabulary must open, and a database that does not must say why.
//
// Order matters: the foreign-key precondition is checked FIRST, because it is
// about the connection rather than about the rows, and because SQLite's own
// procedure puts it at step 1. A run that is wrong in both ways therefore hears
// about the environment problem before the data problem.
func applyV4(tx txExec) error {
	ctx := context.Background()
	if err := requireForeignKeysOff(tx); err != nil {
		return err
	}
	off, err := taskLogStatesOutsideRoster(ctx, tx)
	if err != nil {
		return err
	}
	if len(off) > 0 {
		return fmt.Errorf("task_log.state holds %d value(s) outside the v4 roster (%s); "+
			"refusing to rebuild rather than abort on a CHECK that cannot name a row. "+
			"Accepted values: %s",
			len(off), strings.Join(off, ", "), strings.Join(TaskStateRoster(), "|"))
	}
	for _, stmt := range splitSQLStatements(ddlV4) {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
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
