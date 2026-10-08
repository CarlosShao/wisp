package main

// Ticket 33 slice B (票 33 AC#9) - the first PRODUCTION listener for the inbound
// hop's router.
//
// WHAT THIS IS. composer_dispatch.go (slice A) is the hop that takes one
// postMessage string to the handler that owns it, and it shipped with zero
// production callers - `grep -rn "ComposerDispatch" internal/ cmd/ tools/ |
// grep -v _test.go` answered only its own definition and comments. This file is
// that missing caller: `wisp panel-inbound` reads one raw envelope per stdin
// line and hands the bytes to (*panel.ComposerDispatch).Handle. A second
// production caller now exists - cmd/wisp/panel_host_windows.go's dispatchRaw,
// the C27 host's dispatch door - and the judgement in panel_inbound_33_test.go
// reads the cmd/wisp call sites off the disk rather than off this paragraph,
// so deleting the line below reddens a case instead of passing quietly.
//
// WHY THIS SHAPE AND NOT THE WINDOW. The 33-a1 census (docs/evidence/s1/
// 33-inbound-hop-design-a1.md §4) put the only common hole across H4-H9 in "raw
// 交进来并按方法分流", and put WebView2 in H2/H3/H10 only. H4/H5 are what ticket
// 186's Go half, ticket 187's skeleton and ticket 114's AC#5 all wait on, so this
// leg attaches the router to the channels that already exist on this machine -
// stdin and the assembled write leg - rather than to a window that does not.
// dispatch §1 命名的"未接线那一段"因此不在本文件里，而在下面的 WHAT THIS IS NOT.
//
// THE ASSEMBLY IS THE PRODUCTION ONE. config.NewManager -> perm.New ->
// panel.ModeWriteHandler is the same chain cmd/wisp/run.go:417 builds for
// `wisp run`, with the same constructors and the same object types; this leg
// reads and writes the same config.toml [risk] permission_mode that the CLI
// write leg owns. It is not a test host: nothing here is reachable only from a
// *_test.go file, and the observable consequence of an accepted request is a
// changed byte on disk.
//
// CONFIRM IS NIL, ON PURPOSE. run.go hands the Store the console L2 leg
// (rt.confirmModeSwitch, which raises a C18 card). This leg has no card to raise
// - it is a pipe, not a window - and perm.Options says "nil = auto_approve is
// unreachable (fail closed)", which is the same posture ModeWriteHandler takes
// for a widening request (composer_handlers.go:126 ErrNoL2Confirm). Attaching a
// fake confirm here would be the widest thing this file could do, so it attaches
// none: a request that widens the档 is refused, audited, and says which leg is
// missing.
//
// WHAT THIS IS NOT, stated so nobody infers it from the file's existence:
//   - it is not the WebView2 receiver (H2/H3). The page's postMessage still has
//     no listener in this tree, and no go.mod dependency moved for this slice.
//   - it is not the reply path (H10). Handle returns the refusal sentence and,
//     since ticket 248, the settings receipt; this leg prints both on the console.
//     A page still gets nothing back from THIS process - the WebView2 host hands
//     Handle's return value to the awaiting binding (panel_host_windows.go's bind
//     closure), and that host is ticket 33's, not this leg's.
//   - it installs no persistent log sink (installLogSink). That call would make
//     this leg a records-booking leg under ticket 131's AC#4 and owe that gate a
//     registered nail; the honest, smaller step is to keep the audit on stderr
//     here (the "[audit]" family, same prefix rt.auditf uses) and let the sink
//     decision be made when this leg grows into the resident process.
//   - it answers two of the six whitelisted methods. The mode write is ticket
//     114's leg and the settings read/write is ticket 248's; workspace/attachment/
//     message handlers are ticket 186 / ticket 92 / ticket 35's land, so those
//     three sockets stay nil and the router refuses them by name (ErrNoHandlerAttached)
//     rather than pretending.
//
// EXIT CODES follow the same three-way split `wisp run` uses (SPEC-05 §3.4): 2 =
// this leg could not assemble (no data root, unreadable config), 1 = at least one
// request was refused, 0 = every line that arrived was handled.

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/secret"
)

// panelInboundActor is how perm's switch record names whoever rang the doorbell
// on this leg. It is a fact about this process, not a user identity: an
// operator-driven pipe is not the panel, and pretending otherwise would blur the
// two origins the audit exists to tell apart.
const panelInboundActor = "cli-panel-inbound"

// panelInboundMaxLine is the ceiling for one envelope line. The attachment
// envelope carries a payload reference, not bytes, so this is generous already;
// it is here because a scanner with no ceiling turns one runaway line into an
// allocation the pipe does not have.
const panelInboundMaxLine = 1 << 20

// panelInboundIO is this leg's whole world: where raw arrives, where the user's
// sentence goes, and which data root it writes. An empty dataDir is NOT a
// fallback: see the refusal inside cmdPanelInbound for why this leg takes its
// root from -data and never asks the OS.
type panelInboundIO struct {
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	dataDir string
}

// cmdPanelInbound is the leg main.go dispatches for `wisp panel-inbound`.
func cmdPanelInbound(args []string, s panelInboundIO) int {
	if s.stdin == nil {
		s.stdin = os.Stdin
	}
	if s.stdout == nil {
		s.stdout = os.Stdout
	}
	if s.stderr == nil {
		s.stderr = os.Stderr
	}
	dir, bad := parsePanelInboundFlags(args, s.dataDir)
	if bad != "" {
		if strings.HasPrefix(bad, "help:") {
			fmt.Fprintln(s.stdout, strings.TrimPrefix(bad, "help:"))
			return 0
		}
		fmt.Fprintf(s.stderr, "wisp panel-inbound: %s\n%s\n", bad, panelInboundUsage)
		return 2
	}
	if dir == "" {
		// This leg does NOT resolve the data root. It is a deliberate limit and
		// not a shortcut: every function in this package that calls
		// resolveDataDir is registered in ticket 128 AC#2's refusal-leg table
		// (cmd/wisp/dataroot_128_test.go:113), and that file is named by another
		// ticket's acceptance criteria, so it is not this slice's to edit. The
		// shape AC#2 rules on is "no root, no write, and the reason out loud",
		// which is what an absent -data does here: refuse with 2 before any
		// config.toml is opened, so nothing can land in the start-up directory
		// either way. The host that owns the resolution is ticket 33's WebView2
		// leg (H2/H3), and it inherits the run/resident legs' resolution when it
		// arrives; 33-a1 §4 names that window as the part still missing, not the
		// router this file calls.
		fmt.Fprintln(s.stderr, "wisp panel-inbound: 未指定数据根（-data <目录>），按票 128 AC#2 的形状拒绝启动并写任何文件\n"+panelInboundUsage)
		return 2
	}

	// One audit sink for both halves of the hop: the router's refusal lines and
	// the handler's MODE-REFUSED lines come out of the same function, so the two
	// cannot drift into two formats that a reader has to reconcile by eye.
	auditf := func(format string, args ...any) {
		fmt.Fprintln(s.stderr, "[audit] "+fmt.Sprintf(format, args...))
	}

	disp, err := newPanelInboundDispatch(dir, auditf)
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp panel-inbound: 入向装配未完成（Unconfigured）：%v\n", err)
		return 2
	}

	ctx := context.Background()
	sc := bufio.NewScanner(s.stdin)
	sc.Buffer(make([]byte, 0, 64*1024), panelInboundMaxLine)

	var arrived, accepted, refused int
	for lineNo := 0; sc.Scan(); lineNo++ {
		raw := strings.TrimSpace(sc.Text())
		if raw == "" {
			continue
		}
		arrived++
		reply, err := disp.Handle(ctx, raw)
		if err != nil {
			refused++
			// The sentence is Handle's, not a paraphrase of it: what the user is
			// told and what the audit recorded are then the same string, which is
			// the property a reader of a refusal actually wants.
			fmt.Fprintf(s.stdout, "wisp panel-inbound: 第 %d 行被拒绝：%s\n", lineNo+1, reply)
			continue
		}
		accepted++
		// An accepted request returns no text on purpose for the four doors slice A
		// shipped (the router must not become the second voice in the room), so the
		// receipt line below is this leg's own and says exactly which hop happened.
		// A settings request does return a sentence, and ticket 248 AC#8 is why: which
		// key paths landed, and that a restart is required, must reach the surface the
		// user reads instead of only the log. Printing Handle's own words keeps this leg
		// a pipe rather than an author.
		if reply != "" {
			fmt.Fprintf(s.stdout, "wisp panel-inbound: 第 %d 行回执：%s\n", lineNo+1, reply)
			continue
		}
		fmt.Fprintf(s.stdout, "wisp panel-inbound: 第 %d 行已受理，处理器已被调用\n", lineNo+1)
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(s.stderr, "wisp panel-inbound: 输入流中断：%v\n", err)
		return 2
	}

	fmt.Fprintf(s.stderr,
		"wisp panel-inbound: 收到 %d 封 / 受理 %d 封 / 拒绝 %d 封（数据根 %s）\n",
		arrived, accepted, refused, dir)
	if refused > 0 {
		return 1
	}
	return 0
}

// newPanelInboundDispatch is the assembly, and the only reason it is a function
// rather than four lines in the loop is so the leg's shape stays readable: five
// constructors, in the order the data flows.
//
// The chain is the production one (see the header): the config Manager is the
// single writer of the档 (ticket 101's rule that loading it twice is the split
// state SPEC-03 §3.1 prevents), perm.Store is the single owner of it,
// panel.ModeWriteHandler is the single handler for panel.mode.request, and
// panel.ComposerDispatch is the single router in front of that handler.
func newPanelInboundDispatch(dir string, auditf panel.AuditFunc) (*panel.ComposerDispatch, error) {
	// Ticket 223 AC#4's fourth sentence, and it is this leg's true state: the tick
	// that re-reads config.toml lives in `wisp run` (config_reload.go), not here.
	// panel-inbound owns no approval gate - its Confirm is nil below, on purpose -
	// so a loosening it read could not be confirmed even if it polled, and polling
	// without a card surface is how a silent deny gets mistaken for hot reload.
	// Rather than fake the tier, this line says which host did not take the job,
	// and says it once per start instead of leaving the operator to infer it from
	// an edit that never lands.
	auditf("%s", hotReloadDisabledPanelInbound)
	return newComposerDispatchChain(dir, auditf, panelInboundActor)
}

// newComposerDispatchChain is the five-constructor body behind
// newPanelInboundDispatch, with the audit actor as a parameter. The actor is not
// decoration: the MODE-SWITCH audit line names whoever asked, so a chain built for
// the resident panel must not sign itself as the CLI seam (ticket 33's resident
// wiring passes "resident-panel"). Every other property of the chain is unchanged,
// including the two nils that keep a widening request from reaching Store.Set.
func newComposerDispatchChain(dir string, auditf panel.AuditFunc, actor string) (*panel.ComposerDispatch, error) {
	cfgPath := filepath.Join(dir, configFileName)
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		return nil, fmt.Errorf("配置未就绪（%s）：%w", cfgPath, err)
	}
	// Ticket 223 AC#4's hot-reload sentence is the CLI leg's own and moved up into
	// newPanelInboundDispatch with it, so this shared body says nothing about a
	// host it does not know.
	store, err := perm.New(perm.Options{
		Manager: mgr,
		// nil: no L2 card on this leg, so auto_approve is unreachable and every
		// widening request is refused before Set is ever reached. Not an
		// oversight - see this file's header.
		Confirm: nil,
		Logf:    auditf,
	})
	if err != nil {
		return nil, fmt.Errorf("权限档位不可用：%w", err)
	}
	modeWrites := &panel.ModeWriteHandler{
		Modes: store,
		// The handler reads only this field's nil-ness to refuse a widening
		// request; the forward below is run.go's, and this leg has no card.
		Confirm: nil,
		Audit:   auditf,
		Actor:   actor,
	}
	// Ticket 248 AC#1/AC#3: the settings leg. It writes through the SAME Manager
	// object as the mode chain above (one truth, no second loader) and through the
	// one secret store this repository has (internal/secret, DPAPI), so this route
	// creates no second place a key could live. A store-directory failure is an
	// assembly failure here, not a degraded settings page that quietly accepts
	// writes it cannot persist.
	secrets, err := secret.NewStore(dir)
	if err != nil {
		return nil, fmt.Errorf("凭据存储不可用（%s）：%w", dir, err)
	}
	configWrites := &panel.ConfigWriteHandler{
		Store: newConfigStore(mgr, secrets, auditf, actor),
		Audit: auditf,
		Actor: actor,
	}
	return &panel.ComposerDispatch{
		Mode: modeWrites,
		// Untouched by this slice, and refused by name when they arrive:
		// workspace = ticket 186, attachment = ticket 92, message = ticket 35.
		Workspace:  nil,
		Attachment: nil,
		Message:    nil,
		// Ticket 248's two settings doors are answered by the leg above.
		Config: configWrites,
		Audit:  auditf,
	}, nil
}

const panelInboundUsage = `用法：wisp panel-inbound -data <目录>
  从 stdin 逐行读入面板 composer 的原始封套（每行一枚 JSON），交进
  panel.ComposerDispatch.Handle 的入向派发腿。-data 是指定的数据根：本腿不解析
  数据根，缺它就按票 128 AC#2 的形状拒绝（不写任何文件）。
  退码：2=装配未完成，1=有请求被拒，0=全部受理。`

// parsePanelInboundFlags takes the one flag this leg has. Manual because the
// package's other legs read argv the same way and a flag.FlagSet would write
// usage text this gate's census cannot find.
func parsePanelInboundFlags(args []string, preset string) (dir string, bad string) {
	dir = preset
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-data":
			if i+1 >= len(args) {
				return "", "-data 需要一个目录参数"
			}
			i++
			dir = args[i]
		case "-h", "-help", "--help":
			return "", "help:" + panelInboundUsage
		default:
			return "", fmt.Sprintf("未识别的参数 %q", args[i])
		}
	}
	return dir, ""
}
