package main

// Ticket 43, crash recovery: a start-after-crash marker and the one-shot
// interrupted-task sweep.
//
// The shape the ticket names is three facts, each of which lives in exactly one
// place:
//
//   - THE MARKER (writeRunMarker / clearRunMarker): a small JSON file in the
//     data root that exists if and only if a wisp process believes itself to be
//     running. It is written at resident boot and removed by the boot's own
//     deferred chain - which a kill -9 never runs. Finding it at the next boot
//     is therefore the definition of "the previous process died without its own
//     shutdown", and nothing else is consulted: not the pid (recycled), not a
//     wall-clock heuristic (D42#9).
//
//   - THE SWEEP (recoverCrashedTasks): runs once per boot, right after the task
//     store exists. Under a live marker, every task_log row still in the
//     'running' state is a task the previous process never closed, and it is
//     finished as interrupted - the same word memory's taskstate roster spells
//     for "cut off" (TaskStateInterrupted, ticket 196's roster, ticket 365's
//     ladder). The original summary, cost columns and error class are passed
//     back unchanged, so the sweep only moves the row to a terminal state; it
//     does not author a verdict the dead process never reached.
//
//   - NO RERUN, by construction rather than by a guard: nothing in this
//     repository auto-executes a task at boot - tasks start from an operator
//     gesture (console, voice, panel) and this file adds no second trigger. The
//     sweep's only writes are the terminal-state transitions above, so the
//     listing a panel or a CLI reads after the sweep shows 中断 exactly once.
//
// The marker is cross-platform on purpose: it is written and cleared from
// runResident, and the sweep runs from assembleRuntime, which the CLI also
// walks - a `wisp run` that follows a killed resident gets the same marking as
// the next resident boot.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// nowWallUnix is the marker's one persisted timestamp (D42#9: wall clock for
// records only, never for interval math).
func nowWallUnix() int64 { return observe.NowWallUTC().Unix() }

// runMarkerFile is the marker's file name inside the data root. It lives next
// to config.toml and wisp.db: the same directory the whole lifecycle of this
// installation is already keyed by.
const runMarkerFile = "run-marker.json"

// runMarkerBody is what the marker carries. The pid is diagnostics for the
// human reading the file, never a liveness test: pids are recycled, so the
// file's EXISTENCE is the only fact the sweep acts on.
type runMarkerBody struct {
	PID       int   `json:"pid"`
	StartedAt int64 `json:"started_at_wall_utc"`
}

// runMarkerPath is the marker's full path. It is a plain join of two values the
// caller already resolved; the risk.PathResolver rule governs PATH DECISIONS on
// behalf of file tools, not the data-root layout this package owns.
func runMarkerPath(dataDir string) string {
	return filepath.Join(dataDir, runMarkerFile)
}

// writeRunMarker records that this process is alive. A previous marker that
// somehow survived (a boot that crashed between the sweep and its own defer) is
// overwritten, not merged: the newest boot is the one whose death the file
// should describe.
func writeRunMarker(dataDir string) error {
	body, err := json.Marshal(runMarkerBody{PID: os.Getpid(), StartedAt: nowWallUnix()})
	if err != nil {
		return fmt.Errorf("crash-recovery: marker 序列化失败: %w", err)
	}
	p := runMarkerPath(dataDir)
	if err := os.WriteFile(p, body, 0o600); err != nil {
		return fmt.Errorf("crash-recovery: 写启动标记 %s 失败: %w", p, err)
	}
	return nil
}

// clearRunMarker removes the marker on the way out. It is registered as a defer
// by the resident boot, so every exit that walks that chain - signal, tray,
// WM_QUERYENDSESSION fast path - clears it; a kill -9 (or a power loss) is the
// remaining way the file survives, which is exactly the meaning it carries.
func clearRunMarker(dataDir string) error {
	err := os.Remove(runMarkerPath(dataDir))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("crash-recovery: 清除启动标记失败: %w", err)
	}
	return nil
}

// hasRunMarker answers whether the previous process died without its own
// shutdown. A marker that fails to PARSE still counts: a half-written file from
// a power loss mid-write is as much evidence of an unclean exit as an intact
// one, and refusing to act on it would trade one unclean exit for a silent
// unmarked one.
func hasRunMarker(dataDir string) bool {
	_, err := os.Stat(runMarkerPath(dataDir))
	return err == nil
}

// readRunMarker returns the parsed marker, or the parse failure. It exists for
// the boot report's diagnostics; the sweep only asks hasRunMarker.
func readRunMarker(dataDir string) (runMarkerBody, error) {
	body, err := os.ReadFile(runMarkerPath(dataDir))
	if err != nil {
		return runMarkerBody{}, err
	}
	var m runMarkerBody
	if err := json.Unmarshal(body, &m); err != nil {
		return runMarkerBody{}, fmt.Errorf("crash-recovery: 启动标记解析失败: %w", err)
	}
	return m, nil
}

// recoverCrashedTasks is the sweep. Under a live marker it finishes every still
// 'running' task_log row as interrupted and reports how many it moved; with no
// marker it is a no-op that answers 0, nil, so the call site needs no branch.
//
// The write goes through the store's own FinishTaskLog (validated terminal
// state, one UPDATE, ended_at set) with the row's own cost columns handed back,
// and the summary slot passed empty: FinishTaskLog keeps the existing summary
// when the argument is empty (COALESCE in dao_tasklog.go), so the sweep does not
// overwrite what the dead process had already written. The diagnostics hint the
// ticket asks for is the log line this function and the call site print - it
// names the count and the no-rerun posture where the operator reads it.
func recoverCrashedTasks(ctx context.Context, store *memory.Store, dataDir string) (int, error) {
	if store == nil {
		return 0, nil
	}
	if !hasRunMarker(dataDir) {
		return 0, nil
	}
	marker, merr := readRunMarker(dataDir)
	rows, err := store.ListTaskLogs(ctx)
	if err != nil {
		return 0, fmt.Errorf("crash-recovery: 读取任务列表失败: %w", err)
	}
	n := 0
	for _, tl := range rows {
		if tl.State != string(memory.TaskStateRunning) {
			continue
		}
		errClass := tl.ErrorClass
		if err := store.FinishTaskLog(ctx, tl.ID, string(memory.TaskStateInterrupted),
			"", tl.CostTokensIn, tl.CostTokensOut, tl.CostAmountMicro, errClass); err != nil {
			return n, fmt.Errorf("crash-recovery: 任务 %s 标记中断失败: %w", tl.ID, err)
		}
		n++
	}
	if n > 0 {
		who := "上一次进程"
		if merr == nil {
			who = fmt.Sprintf("上一次进程（pid=%d）", marker.PID)
		}
		slog.Warn("crash-recovery: 检测到上次进程异常退出，未收口任务已标记为中断",
			"interrupted", n, "previous", who, "note", "不会自动重跑（任务只能由操作者再次发起）")
	}
	return n, nil
}

// updatePendingFile is the apply-handoff cue ticket 56's updater leaves behind.
// Ticket 43's half is the hook alone: on boot, if the file is present, book the
// event the updater consumes. The file is NOT consumed here (renamed or
// deleted): the updater owns that protocol, and a boot that ate the cue would
// make ticket 56's rollback counters start with no one having applied anything.
const updatePendingFile = "update-pending.json"

// markUpdatePendingHandoff books the boot-time apply-handoff event when the
// cue exists. Today's only consumer is the audit log (ticket 56 owns the
// updater binary and the N-1 rollback counters this hook starts); the event is
// booked by name so the consumer can grep for it rather than infer from
// silence.
func markUpdatePendingHandoff(dataDir string) bool {
	p := filepath.Join(dataDir, updatePendingFile)
	if _, err := os.Stat(p); err != nil {
		return false
	}
	slog.Info("update-pending: 检测到待应用的更新交接文件，登记 apply-handoff 事件",
		"event", "UPDATE-APPLY-HANDOFF", "path", p,
		"note", "更新器本体是票 56；本事件由该票消费，此处不移动也不删除文件")
	return true
}
