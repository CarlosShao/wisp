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
// WISP-LEG-COVERAGE-RULING: sessions is a console read with no listener and no
// record of its own: it resolves the data root, opens memory.db, asks
// session.Manager.List and prints what came back. The behaviour that can regress
// behind it is pinned against real SQLite rows in internal/session/manage_test.go
// (TestSessionListByTime for the newest-first order, TestSessionListByWorkspace
// for the workspace filter, TestSessionRenameAndArchiveUndo for the archived
// default) and, for the search outlet, TestSessionSearchTitleAndBody (title and
// body both match) plus TestAC2SearchHonoursArchiveMaskAndWorkspace (the
// archived/workspace mask) - this directory would only re-assert those same rows
// through a second manager, and this leg's own wiring has one failure shape an operator
// reads (the refusal lines below). What this directory does drive is the
// no-data-root refusal, through the dataroot_128 leg table. If the leg grows a
// second decision, fix that by writing the case here and registering it in the
// dispatch gate rather than widening this ruling.

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

The rows are the ones already on disk (memory.db, table session). This command
adds no query of its own: it opens the store the text path opens and calls the
existing Manager.List / Manager.Search.
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
// the ones the other console legs use: 0 listed or matched, 2 Unconfigured (no
// data root, storage unusable, wrong arguments), 1 the store refused the read.
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
