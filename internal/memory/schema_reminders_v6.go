package memory

import (
	"context"
	"strings"
)

// ddlV6 adds the reminders table (ticket 31, leg 31-r2): the one-shot
// in-process reminder row. Same rule as ddlV3/ddlV5: a NEW table in its own
// step, no edit of ddlV1..ddlV5 text, so a fresh database reaches the table
// through the same 5->6 step an existing one does — and existing rows are
// never touched, so 旧数据不丢.
//
// IF NOT EXISTS everywhere: ticket 31 asks the step to be idempotent. The
// chain already gates it on from==5, but a retried transaction re-running the
// step must not error and must not restate any row.
//
// The status CHECK is GENERATED from reminderStatusRoster (dao_reminders.go),
// never spelled a second time — ticket 196's rule applied to a new enum for
// the first time: the Go roster is the authority, the CHECK is one of its
// mouths, the DAO validator is the other.
//
// fire_at_wallclock_utc is a WALL-CLOCK instant (「下午三点」是墙钟概念); the
// WAIT to it is computed on the monotonic clock by the caller
// (internal/tools + internal/agent), so a system-clock adjustment never moves
// a scheduled fire.
//
// 禁止 这张表刻意表达不了的（票 31 拒绝清单；同样写进工具注册的 note）：
//   - 无 recurrence/cron 列：不做循环提醒；
//   - 无 action/payload 列：fire 只把 note 推给 Notifier 端口，绝不触发任何
//     工具执行；
//   - 不是 OS 计划任务、重启不补发：boot 只对 status='armed' 且已到点的行回
//     一个计数（CountMissedArmed），fire 路径只活在进程内。
//
// SPEC-02 §3 still lists nine tables; this is the fourth roster extension of
// the same kind as v3/v4/v5 (AGENTS.md §0「改契约＝人工批准」), named in
// docs/evidence/s1/31r2-impl.md rather than fixed by an agent.
var ddlV6 = buildDDLV6()

func buildDDLV6() string {
	quoted := make([]string, len(reminderStatusRoster))
	for i, st := range reminderStatusRoster {
		quoted[i] = "'" + st + "'"
	}
	return `-- 一次性进程内提醒（schema v6：票 31 / 31-r2；SPEC-02 §3 roster 的第四枚扩展，spec-note 债摆机主批）
-- 禁止 无 recurrence/cron 列 · 禁止 无 action/payload 列（fire 只推通知，绝不触发工具）
-- 禁止 不是 OS 计划任务：进程退出即失效，重启只计数不补发
CREATE TABLE IF NOT EXISTS reminders (
  id                    TEXT PRIMARY KEY,            -- 'rem_' + crypto/rand hex（C31 同形铸造）
  fire_at_wallclock_utc INTEGER NOT NULL,            -- unix 秒（墙钟目标时刻）
  note                  TEXT NOT NULL,               -- 给用户的一句话，不是指令
  status                TEXT NOT NULL CHECK(status IN (` + strings.Join(quoted, ",") + `)),
  created_utc           INTEGER NOT NULL,            -- unix 秒
  fired_utc             INTEGER                      -- NULL = 尚未 fire（armed/cancelled 均 NULL）
);
CREATE INDEX IF NOT EXISTS idx_reminders_due ON reminders(status, fire_at_wallclock_utc);`
}

// applyV6 creates the reminders table (ticket 31). Runs inside the chain's
// single transaction for the step, after the pre-backup.
func applyV6(tx txExec) error {
	for _, stmt := range splitSQLStatements(ddlV6) {
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return err
		}
	}
	return nil
}
