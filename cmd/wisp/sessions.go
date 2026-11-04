package main

// wisp sessions - the cross-restart conversation list (ticket 364-r1a).
//
// internal/session.Manager has owned the four conversation outlets since ticket
// 166 (list / search / rename / archive with undo), but no operator-facing
// command had ever wired even the first one up: "what did I talk about last
// time" could only be answered by opening the panel. This file is the smallest
// honest wiring of outlets #1 and #2:
//
//	wisp sessions list          the newest sessions, first: id, updated_at, title
//	wisp sessions search <q>    the sessions whose title or a turn mentions q
//
// They call session.Manager.List and session.Manager.Search - the methods that
// already exist - over the same memory.Store `wisp run` opens, and print the
// rows they get back through one shared row printer. The ordering is not decided
// here either: internal/memory's ListSessions and SearchSessions both order by
// updated_at DESC, id DESC, and the Manager passes the options straight through,
// so the zero SessionListOptions is exactly "every active session, newest first".
// Archived rows stay hidden for the same reason (that is the DAO's default).
//
// Ticket 364-r1c wires outlets #3 and #4 - the ones that change a row instead of
// reading rows:
//
//	wisp sessions rename <id> <title>   rewrite one session's title column
//	wisp sessions archive <id>          hide it from list / search
//	wisp sessions unarchive <id>        undo that archive
//
// They call session.Manager.Rename / Archive / Unarchive - the methods
// TestSessionRenameAndArchiveUndo already pins against real rows - so this file
// decides nothing about what a rename or an archive means; the console supplies
// the id (and, for rename, the title) and prints one confirmation line. All
// three go through sessionsMutate, so archive and unarchive cannot drift apart
// in the store they open or the exit code they return. A wrong argument shape is
// a usage error; a blank title is left to Manager.Rename's own refusal, because
// that is the layer that owns the "a session needs a title" rule.
//
// WISP-LEG-COVERAGE-RULING: sessions is a console command with no listener and
// no record of its own: it resolves the data root, opens memory.db, asks
// session.Manager and prints what came back. The behaviour that can regress
// behind it is pinned against real SQLite rows in internal/session/manage_test.go
// (TestSessionListByTime for the newest-first order, TestSessionListByWorkspace
// for the workspace filter, TestSessionRenameAndArchiveUndo for the archived
// default and for the r1c rename / archive / unarchive round trip) and, for the
// search outlet, TestSessionSearchTitleAndBody (title and body both match) plus
// TestAC2SearchHonoursArchiveMaskAndWorkspace (the archived/workspace mask) -
// this directory would only re-assert those same rows through a second manager,
// and this leg's own wiring has one failure shape an operator reads (the refusal
// lines below). What this directory does drive is the no-data-root refusal,
// through the dataroot_128 leg table. If the leg grows a second decision, fix
// that by writing the case here and registering it in the dispatch gate rather
// than widening this ruling.

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/session"
)

const sessionsUsage = `wisp sessions - the cross-restart conversation list

Usage:
  wisp sessions list
        print the most recent sessions, newest first, one line each:
        <id>  <updated_at>  <title>. The 20 newest are shown; archived sessions
        stay hidden, exactly as session.Manager.List returns them.

  wisp sessions search <query>
        print the sessions whose title or any turn text contains <query>, in the
        same <id>  <updated_at>  <title> rows, capped at the same 20. Quote the
        query if it has spaces.

  wisp sessions rename <id> <title>
        give session <id> a new title. Quote the title if it has spaces (this
        takes exactly two arguments). Only the title column is rewritten: the
        id, the turns and the created_at of the history stay as they were.

  wisp sessions archive <id>
        hide session <id> from list and search. Nothing is deleted - the row
        and its turns are still on disk.

  wisp sessions unarchive <id>
        undo archive: bring session <id> back into list and search.

The rows are the ones already on disk (memory.db, table session). This command
adds no query of its own: it opens the store the text path opens and calls the
existing Manager.List / Manager.Search / Manager.Rename / Manager.Archive /
Manager.Unarchive.
`

// sessionsDefaultLimit is how many rows `wisp sessions list` prints. It is a
// display cap only - the read itself is the Manager's untouched List.
const sessionsDefaultLimit = 20

// sessionsIO is the process surface of one `wisp sessions` run, injectable for
// the same reason providersIO and modelsIO are: the dataroot refusal leg drives
// this command in-process and reads what it printed.
type sessionsIO struct {
	stdout  io.Writer
	stderr  io.Writer
	dataDir string
}

func (s sessionsIO) out() io.Writer {
	if s.stdout != nil {
		return s.stdout
	}
	return os.Stdout
}

func (s sessionsIO) err() io.Writer {
	if s.stderr != nil {
		return s.stderr
	}
	return os.Stderr
}

// cmdSessions runs `wisp sessions ...` and returns the exit code. The codes are
// the ones the other console legs use: 0 listed, matched or changed, 2
// Unconfigured (no data root, storage unusable, wrong arguments), 1 the store
// refused the read or the write.
func cmdSessions(argv []string, io_ sessionsIO) int {
	if len(argv) == 0 {
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	if io_.dataDir == "" {
		dir, dirErr := resolveDataDir(buildinfo.EnvString())
		if dirErr != nil {
			// Ticket 128 AC#2: the same refusal `wisp run`, `wisp models` and
			// `wisp providers` make. Listing sessions out of the start-up
			// directory would answer with a different memory.db than the one
			// the agent writes - a silent empty list, not a list.
			fmt.Fprintf(io_.err(), "wisp sessions: %v\n", dirErr)
			return 2
		}
		io_.dataDir = dir
	}
	switch argv[0] {
	case "list":
		return sessionsList(io_, argv[1:])
	case "search":
		return sessionsSearch(io_, argv[1:])
	case "rename":
		return sessionsRename(io_, argv[1:])
	case "archive":
		return sessionsArchive(io_, argv[1:])
	case "unarchive":
		return sessionsUnarchive(io_, argv[1:])
	default:
		fmt.Fprintf(io_.err(), "wisp sessions: 未知子命令 %q\n\n", argv[0])
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
}

// sessionsList is outlet #1 through the console: Manager.List, then print.
func sessionsList(io_ sessionsIO, rest []string) int {
	if len(rest) != 0 {
		fmt.Fprintf(io_.err(), "wisp sessions: list 不接受参数，收到 %q\n\n", strings.Join(rest, " "))
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	store, err := memory.Open(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: 存储不可用：%v\n", err)
		return 2
	}
	defer func() { _ = store.Close() }()

	mgr, err := session.NewManager(session.ManagerOptions{Store: store})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: %v\n", err)
		return 2
	}
	// No deadline on this read: it is one local query, not a task, and D42#9's
	// rule is against inventing one out of a wall-clock difference.
	rows, err := mgr.List(context.Background(), memory.SessionListOptions{})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: 读取会话列表失败：%v\n", err)
		return 1
	}
	sessionsPrint(io_,
		"wisp sessions: 最近 %d 个会话（共 %d 个，按更新时间降序）\n",
		"wisp sessions: 还没有会话",
		rows)
	return 0
}

// sessionsSearch is outlet #2 through the console: Manager.Search, printed
// through the same row formatter list uses. It needs exactly one query - no
// argument and a blank one are the same operator mistake ("what should I look
// for?"), so both are usage errors rather than an empty listing.
func sessionsSearch(io_ sessionsIO, rest []string) int {
	if len(rest) != 1 || strings.TrimSpace(rest[0]) == "" {
		fmt.Fprintf(io_.err(), "wisp sessions: search 需要恰好一个查询词，收到 %q\n\n", strings.Join(rest, " "))
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	query := rest[0]

	store, err := memory.Open(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: 存储不可用：%v\n", err)
		return 2
	}
	defer func() { _ = store.Close() }()

	mgr, err := session.NewManager(session.ManagerOptions{Store: store})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: %v\n", err)
		return 2
	}
	// Same reasoning as list: one local query, no invented deadline.
	rows, err := mgr.Search(context.Background(), query, memory.SessionListOptions{})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: 搜索会话失败：%v\n", err)
		return 1
	}
	sessionsPrint(io_,
		"wisp sessions: 匹配 %d 个会话（共 %d 个，按更新时间降序）\n",
		"wisp sessions: 没有匹配的会话",
		rows)
	return 0
}

// sessionsRename is outlet #3 through the console: Manager.Rename, then one
// confirmation line. It takes exactly two arguments, the id and the title, so a
// title with spaces has to be quoted - the same bargain search makes with its
// query. A third word is an operator mistake, not a longer title: joining the
// tail here would make "rename <id> my new title" and a typo indistinguishable.
func sessionsRename(io_ sessionsIO, rest []string) int {
	if len(rest) != 2 {
		fmt.Fprintf(io_.err(), "wisp sessions: rename 需要 <id> 和 <title> 两个参数，收到 %q\n\n", strings.Join(rest, " "))
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	id, title := rest[0], rest[1]
	return sessionsMutate(io_, "改名", "已改名 "+id, func(ctx context.Context, mgr *session.Manager) error {
		return mgr.Rename(ctx, id, title)
	})
}

// sessionsArchive is outlet #4's forward half: Manager.Archive, then the same
// confirmation shape rename prints. It takes exactly one id.
func sessionsArchive(io_ sessionsIO, rest []string) int {
	if len(rest) != 1 {
		fmt.Fprintf(io_.err(), "wisp sessions: archive 需要恰好一个会话 id，收到 %q\n\n", strings.Join(rest, " "))
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	id := rest[0]
	return sessionsMutate(io_, "归档", "已归档 "+id, func(ctx context.Context, mgr *session.Manager) error {
		return mgr.Archive(ctx, id)
	})
}

// sessionsUnarchive is outlet #4's undo, addressed by id rather than by "undo
// the last archive" - a console has no session state of its own to undo, so the
// operator names the row. Manager.Unarchive is a first-class outlet for the same
// reason it exists there: an archive you cannot take back is a delete.
func sessionsUnarchive(io_ sessionsIO, rest []string) int {
	if len(rest) != 1 {
		fmt.Fprintf(io_.err(), "wisp sessions: unarchive 需要恰好一个会话 id，收到 %q\n\n", strings.Join(rest, " "))
		fmt.Fprint(io_.err(), sessionsUsage)
		return 2
	}
	id := rest[0]
	return sessionsMutate(io_, "取消归档", "已取消归档 "+id, func(ctx context.Context, mgr *session.Manager) error {
		return mgr.Unarchive(ctx, id)
	})
}

// sessionsMutate is the one opener/printer the three row-changing outlets share:
// open the same store the read outlets above open, compose a Manager, make the
// one call, print the one confirmation line. Sharing it is what keeps archive
// and unarchive from drifting in the store they touch or the code they return -
// rename goes through it too, so all three fail the same way.
//
// The split matches list and search: a storage that will not open is rc=2 (no
// data root, unusable file), while a call the layer below refused is rc=1 -
// that covers an unknown id (the DAO's RowsAffected==0 check, so a typo cannot
// print "已归档 <typo>" as success), a blank title (Manager.Rename's own rule)
// and any SQLite error. No deadline: one local statement, and D42#9's rule is
// against inventing one out of a wall-clock difference.
func sessionsMutate(io_ sessionsIO, op, confirm string, apply func(context.Context, *session.Manager) error) int {
	store, err := memory.Open(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: 存储不可用：%v\n", err)
		return 2
	}
	defer func() { _ = store.Close() }()

	mgr, err := session.NewManager(session.ManagerOptions{Store: store})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: %v\n", err)
		return 2
	}
	if err := apply(context.Background(), mgr); err != nil {
		fmt.Fprintf(io_.err(), "wisp sessions: %s失败：%v\n", op, err)
		return 1
	}
	fmt.Fprintf(io_.out(), "wisp sessions: %s\n", confirm)
	return 0
}

// sessionsPrint is the one place the console row shape lives: cap at
// sessionsDefaultLimit, then <id>  <updated_at>  <title> one per line, the
// title flattened by oneLine. heading is a fmt format that receives the shown
// and total counts; empty is printed when there is nothing to show. list and
// search both print through it so the two outlets cannot drift apart in the
// columns, the cap, or the stamp they use.
func sessionsPrint(io_ sessionsIO, heading, empty string, rows []memory.Session) {
	shown := rows
	if len(shown) > sessionsDefaultLimit {
		shown = shown[:sessionsDefaultLimit]
	}
	if len(shown) == 0 {
		fmt.Fprintln(io_.out(), empty)
		return
	}
	fmt.Fprintf(io_.out(), heading, len(shown), len(rows))
	for _, s := range shown {
		// updated_at is the stamp the list is ordered by (ticket 326 AC#1: the
		// conversation you keep talking in comes back to the top), so that is
		// the stamp printed - not created_at.
		fmt.Fprintf(io_.out(), "  %s  %s  %s\n", s.ID, sessionStamp(s.UpdatedAt), oneLine(s.Title))
	}
}

// sessionStamp renders a stored Unix second the way the rest of the console
// legs render audit stamps (secret.go's timeLayout): UTC, RFC3339, seconds.
func sessionStamp(unix int64) string {
	return time.Unix(unix, 0).UTC().Format(timeLayout)
}
