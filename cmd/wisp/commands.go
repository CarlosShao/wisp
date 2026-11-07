package main

// wisp commands - ticket 213 AC#4: the command roster's first production caller.
//
// WHAT THIS IS. internal/command owns the roster (ticket 213 leg r1: the rows, the
// AC#2 door gate, the AC#3 verdict), and until this file it had zero production
// callers - `grep -rn --include=*.go '"github.com/CarlosShao/wisp/internal/command"' cmd/
// internal/ | grep -v '^internal/command/'` answered 0. A roster nothing reads is a
// menu with no surface: AC#4's terminal state is "名册可见 + 命令行可执行，界面点不了",
// and the first half of that sentence needs a leg an operator can actually run.
//
// WHY THE CONSOLE AND NOT THE PANEL. The page's own path to this list is a snapshot
// key or a C17 method, and both are other tickets' faces: the snapshot's top-level
// key set is pinned by a contract ruler that moves on both sides at once
// (pump_test.go's four keys, composer_test.go's interface match), and a C17 method
// name is a frozen closed set that ticket 213's禁区 forbids adding from here. The
// console is the surface that already exists, needs no new contract, and prints the
// same rows the page will one day read - so the roster becomes measurable in the
// product without a second copy of it.
//
// IT PRINTS THE ROSTER, IT DOES NOT KEEP ONE. The rows below come from
// command.Production() itself, never from a table written here: a private copy
// would keep printing after the roster refused a row, and the whole point of AC#2
// is that a command with no door cannot be listed. The line this leg prints for a
// row is derived from the row (name, summary, args, approval, door), and the
// command-line spelling of each door is derived from the roster's own entries - the
// control words in the footer are the utterances the roster carries, and the method
// in the envelope example is the roster's panel method. There is no second spelling
// of a frozen name in this file.
//
// NO EXECUTOR HERE EITHER. This leg reads and prints. It opens no store, resolves no
// data root (ticket 128 AC#2's table of legs that resolve one does not grow because
// of it), and calls none of the doors: the rows tell the operator which door to
// ring and where that door lives, and the ring happens in the two legs that own it.
// That is also why it books no records and installs no log sink - a listing that
// dies with the terminal is not a missing audit trail, it is a listing.

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/command"
	"github.com/CarlosShao/wisp/internal/panel"
)

const commandsUsage = `wisp commands - the command roster this machine can execute (ticket 213)

Usage:
  wisp commands
        print every command in the roster: its name (no slash), one sentence of
        what it does, the argument shape, whether a human approval is involved,
        and which existing door executes it.

The rows come from internal/command's production roster. A command enters that
table only if a production door already executes it, so a name listed here is a
name that works - and a name that is missing here is a capability this machine
does not have yet (/compact, /goal, /plan, /model, /export, /feedback: see the
ticket's census for why each one is still out).

The roster also rides the snapshot now (panel.Snapshot's commands key, ticket
213 交付 3's Go half, wired at the assembly root). The page-side TypeScript
declaration is still unwritten (Q-51), and ringing these from the page stays
where it was: mode's door is the one C17 route the page could already reach,
and the control words have no page entry at all today.
`

// commandsIO is this leg's whole world: two streams, no store. It is injectable for
// the same reason sessionsIO is - the case that pins the output drives this function
// in-process and reads what it printed, rather than shelling out.
type commandsIO struct {
	stdout io.Writer
	stderr io.Writer
}

func (c commandsIO) out() io.Writer {
	if c.stdout != nil {
		return c.stdout
	}
	return os.Stdout
}

func (c commandsIO) err() io.Writer {
	if c.stderr != nil {
		return c.stderr
	}
	return os.Stderr
}

// cmdCommands runs `wisp commands` and returns the exit code: 0 printed, 1 the
// roster refused to assemble (which is a real failure - AC#2's gate just told the
// product that a row has no door), 2 an operator mistake (this leg takes no
// arguments).
func cmdCommands(argv []string, s commandsIO) int {
	if len(argv) != 0 {
		fmt.Fprintf(s.err(), "wisp commands: 不接受参数，收到 %q\n\n%s", strings.Join(argv, " "), commandsUsage)
		return 2
	}
	roster, err := command.Production()
	if err != nil {
		// The roster's own sentence, not a paraphrase: what the operator reads and
		// what the gate said are the same string, which is the property a refusal
		// reader actually wants.
		fmt.Fprintf(s.err(), "wisp commands: 名册装配失败（AC#2 的闸拒了生产表本身，这是缺陷不是用法）：%v\n", err)
		return 1
	}
	entries := roster.Entries()
	fmt.Fprintf(s.out(), "wisp commands: %d 条命令（票 213 名册；每一条都落到一扇已有的生产门上）\n\n", len(entries))
	for _, e := range entries {
		fmt.Fprintf(s.out(), "  %s  %s\n", e.Name, e.Summary)
		fmt.Fprintf(s.out(), "    参数：%s  ·  人批：%s\n", commandArgs(e.Args), e.Approval)
		fmt.Fprintf(s.out(), "    门：%s\n", e.Door())
	}
	fmt.Fprint(s.out(), commandsFooter(entries))
	return 0
}

// commandArgs renders the argument shape of one row: "无" when the command takes
// none, so a menu row never shows a blank field an operator has to interpret.
func commandArgs(args string) string {
	if strings.TrimSpace(args) == "" {
		return "无"
	}
	return args
}

// commandsFooter is the second half of AC#4's terminal state, spelled from the
// roster's own rows: how to ring each door from this console, and the named fact
// that the page cannot. The control words and the method are read out of entries,
// so they cannot drift from the table the rows above were printed from.
func commandsFooter(entries []command.Command) string {
	var controlWords []string
	method := ""
	for _, e := range entries {
		switch e.Channel {
		case command.ChannelControlWord:
			controlWords = append(controlWords, fmt.Sprintf("wisp run %q", e.Utterance))
		case command.ChannelPanelMode:
			method = e.Method
		}
	}
	var b strings.Builder
	b.WriteString("\n怎么从命令行执行（票 213 AC#4 的终态：名册可见 + 命令行可执行）\n")
	if len(controlWords) > 0 {
		fmt.Fprintf(&b, "  控制词门（D11(1)，不经 LLM、不要人批）：%s\n", strings.Join(controlWords, "  /  "))
		b.WriteString("      agent.MatchControl 在请求构建之前解析，(*Loop).runControl 停当前任务。\n")
	}
	if method != "" {
		fmt.Fprintf(&b, "  面板门（C17 既有方法，不新增方法名）：wisp panel-inbound -data <数据根>\n")
		fmt.Fprintf(&b, "      一行一封套：%s\n", commandsEnvelope(method))
		b.WriteString("      没有 -data 就按票 128 AC#2 拒绝启动，一个字节都不写。\n")
	}
	b.WriteString("\n名册在快照里：panel.Snapshot 现在携带 commands 键（票 213 交付 3 的 Go 侧一半，装配根\n" +
		"接了读口），页面侧的 TypeScript 声明（frontend/src/lib/panel.ts）还没人写（Q-51）：两侧\n" +
		"同批的两枚契约尺红着等它，那是界面那支的落笔。\n" +
		"执行照旧：mode 的门本就是页面可达的 C17 既有路由；控制词（取消/停）今天只有命令行——\n" +
		"页面文本到不了 agent.Loop（message 门未接），D11(1) 的入口仍是本控制台。\n")
	return b.String()
}

// commandsEnvelope is the smallest envelope the panel door accepts, built from the
// roster's method and C17's own source constant - the same two fields
// panel.ParseComposerRequest checks before any handler runs.
func commandsEnvelope(method string) string {
	return fmt.Sprintf(`{"method":%q,"requestId":"<请求 id>","source":%q,"to":"<目标档位>"}`,
		method, panel.ComposerRequestSource)
}
