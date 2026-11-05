// Package command owns the user-facing command roster of ticket 213: one
// enumerable list of the commands a person can ring, one entry per command, with
// the four facts a menu needs - a slash-less name, one sentence of what it does,
// whether it costs a human approval, and whether it is runnable under the档 in
// force (with the reason when it is not).
//
// WHAT IT IS NOT. It is not an executor. Every entry NAMES the door that already
// executes that command, and the roster itself holds no store, no loop and no
// handler. That is the whole reason AC#5 ("命令的执行全部走现成审计") can be
// stated as a structural fact instead of a promise: there is no second path to
// audit, because there is no path here at all. The doors today are
//
//   - ChannelControlWord: D11(1)'s control words. internal/agent's MatchControl
//     is the ONE parser (display and execution share it), (*Loop).runControl runs
//     the verb, and the loop's default handler cancels its running task without a
//     provider round-trip (loop.go's run() checks MatchControl before any request
//     is built).
//   - ChannelPanelMode: C17's panel.mode.request, whose production handler is
//     panel.ModeWriteHandler in front of perm.Store.Set, attached by cmd/wisp's
//     two assembly points (run.go for the resident panel, panel_inbound.go for the
//     console leg). No method name is minted here: the entry carries bridge.go's
//     exported constant.
//
// WHY A COMMAND ENTERS ONLY WITH A DOOR (AC#2). The ticket's own failure mode is
// "菜单里有、点了没反应". So New refuses an entry that cannot resolve a door, and
// it refuses two entries claiming the same door - the second shape is how a
// renamed copy of an existing capability would otherwise sneak in. A control
// utterance is resolved through agent.MatchControl, a panel method through
// panel.ParseComposerRequest (which is the closed set in bridge.go, not a copy of
// it). Nothing in this file spells a route name of its own.
//
// WHY THE VERDICT DELEGATES (AC#3). "Does this request need a card" is permission
// policy, and this repository keeps exactly one definition of it:
// perm.JudgeCard/perm.RequiresCard, called by every transition point (perm.Store,
// panel.ModeWriteHandler and this roster). The roster calls the same function; it
// never re-spells `to == auto_approve` or `to > from`. A refusal here therefore
// always carries the sentence the door would give, and the roster cannot drift
// into a second policy. Since R9 the term that decides is the DIRECTION (every
// loosening costs a card, every tightening is free), so the verdict this function
// returns changed for the middle tier without a second rule appearing here.
//
// WHY THE ROSTER LIVES OUTSIDE internal/panel. The frozen grant-boundary
// instrument (internal/panel/l2_grant_boundary_test.go's poolJudgedByRealGuard)
// reports "a second switch/if chain ... exactly the shape a handler registry
// takes" for any route-shaped name in that package's sources that its guard
// answers and its own case list never named. A roster of command rows is exactly
// that shape, so ticket 213's ruling put it outside that package; nothing here
// adds a C17 method, a snapshot key or a top-level key, and the page's own type
// file stays untouched.
//
// AC#4, STATED SO NOBODY INFERS OTHERWISE. Before the roster reaches the page,
// ticket 213's terminal state is "名册可见（可枚举）＋命令行可执行，界面点不了".
// The enumerable half is this package; the executable half is the two doors'
// command-line spellings (`wisp run` carries a control word, `wisp panel-inbound`
// carries a panel.mode.request envelope); and the visible half is now a production
// caller - cmd/wisp's `wisp commands` leg prints this roster, so Production() is
// reached from a binary rather than from this package's tests (roster_test.go
// measures that off the disk, because "a test calls it" and "the product calls it"
// are different facts). What still does not exist is the PAGE's own path: no
// snapshot key, no C17 method, no row in ComposerState, so the menu's members are
// readable on the command line and unreadable in the chat box. Whoever adds that
// key owns it and the two rulers it moves (ticket 213's "交付 3"), and updates this
// paragraph and cmd/wisp/commands.go's footer, which says the same thing to the
// operator.
package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Channel names the existing door that executes one command.
type Channel string

const (
	// ChannelControlWord is D11(1): the utterance is resolved by
	// agent.MatchControl and run by (*Loop).runControl with zero provider
	// round-trips. The entry must carry the utterance, never a synonym.
	ChannelControlWord Channel = "control-word"
	// ChannelPanelMode is C17's panel.mode.request (panel.MethodModeRequest),
	// whose production handler is panel.ModeWriteHandler -> perm.Store.Set.
	ChannelPanelMode Channel = "panel-mode"
)

// Approval says whether ringing one command costs a human approval.
type Approval uint8

const (
	// ApprovalNever: the door never asks. Stopping is not a privilege.
	ApprovalNever Approval = iota
	// ApprovalPerRequest: whether this particular request asks depends on its
	// argument, and the answer comes from the door's own gate (for the mode
	// command: perm.JudgeCard on the from -> to direction). The roster answers it
	// per request in Judge rather than pretending one boolean describes every
	// target.
	ApprovalPerRequest
)

// String renders the approval the way a menu row or a console line shows it. An
// unknown value says so instead of borrowing the nearest sentence: a row that
// cannot classify its own approval is the "clicked and nothing happened" family.
func (a Approval) String() string {
	switch a {
	case ApprovalNever:
		return "不需要"
	case ApprovalPerRequest:
		return "看参数（要不要卡由门自己的闸答）"
	}
	return fmt.Sprintf("未知(%d)", uint8(a))
}

// Verdict is what the roster says about one concrete request under one档.
type Verdict uint8

const (
	// Allowed: the door runs it without asking anybody.
	Allowed Verdict = iota
	// NeedsApproval: the door asks for one L2 strong confirmation first (and an
	// unanswered card is a refusal, not an allow).
	NeedsApproval
	// Refused: the door refuses this request, and Reason says why. Never silent.
	Refused
)

// Availability is the verdict plus the sentence a menu or a log can show.
type Availability struct {
	Verdict Verdict
	Reason  string
}

// Command is one roster row. Name carries no slash (the menu draws it); Utterance
// and Method are the door-specific identifiers, and exactly one of them is set,
// decided by Channel.
type Command struct {
	// Name is the slash-less command name: lowercase ascii letters, digits and
	// hyphens, starting with a letter. No dots - a dotted spelling is a route
	// name, and this is not a route.
	Name string
	// Summary is the one sentence a menu row shows. Empty is refused: a row that
	// cannot say what it does is the "clicked and nothing happened" shape.
	Summary string
	// Args describes the argument shape. Empty means the command takes none.
	Args string
	// Approval says whether a human approval is ever involved.
	Approval Approval
	// Channel names the door.
	Channel Channel
	// Utterance is the D11 control word (ChannelControlWord only).
	Utterance string
	// Method is C17's method constant (ChannelPanelMode only).
	Method string
}

// Door renders the door this row points at in this package's own words: the frozen
// contract, and the identifier that door parses. Nothing here is a synonym - the
// control word is the utterance agent.MatchControl takes, and the method is C17's
// own constant - so a printed row and an executed request name the same thing.
func (c Command) Door() string {
	switch c.Channel {
	case ChannelControlWord:
		return fmt.Sprintf("D11(1) 控制词「%s」（agent.MatchControl 解析、(*Loop).runControl 执行）", c.Utterance)
	case ChannelPanelMode:
		return fmt.Sprintf("C17 %s（panel.ModeWriteHandler -> perm.Store.Set）", c.Method)
	}
	return fmt.Sprintf("未知通道 %q（fail-closed：这条门没有执行者）", c.Channel)
}

// Env is the host's condition for one judgement: the permission档 the door will
// measure the switch FROM, and whether an L2 confirmation leg is attached.
//
// WHICH档, EXACTLY. Since R9 the mode verdict is a judgement about a DIRECTION, so
// the baseline has to be the one the write leg uses - perm.Store.ConfiguredMode,
// the persisted value - because perm.Store.Set compares the persisted pair and the
// panel gate compares the same pair. The EFFECTIVE档 (perm.PermissionMode, what the
// page displays) is NOT the right baseline for that arithmetic: on a legless host
// it collapses to ask_every_step (ticket 372-r2), which would make a configured
// auto_approve -> ask_high_risk request look like a loosening and the roster would
// refuse a change that only tightens the file. A caller that has only the effective
// value should pass that; the verdict then answers for that baseline, which is what
// a caller with no write leg to inspect actually knows.
type Env struct {
	Mode        risk.Mode
	CardChannel bool
}

// Request is one concrete invocation to judge.
type Request struct {
	// Name is the command name the user rang, without the slash.
	Name string
	// Target is the argument as the user typed it ("" = none). For the mode
	// command it is the target档 spelled by risk.ModeNames().
	Target string
}

// Roster is the enumerable list. It is immutable once built.
type Roster struct {
	entries []Command
	byName  map[string]int
	// doors maps a door (channel + identifier) to the command that owns it, so a
	// second command cannot claim a door an existing capability already owns.
	doors map[string]string
}

// New builds a roster out of entries, refusing every row that cannot name a door
// this machine already has. The refusals are the AC#2 gate; they are loud on
// purpose (`internal/tools`' duplicate-name rule, and Pi's "report the clash"
// reading the ticket picked over openchamber's silent ignore).
func New(entries ...Command) (*Roster, error) {
	r := &Roster{byName: make(map[string]int, len(entries)), doors: make(map[string]string, len(entries))}
	for _, e := range entries {
		if err := verify(e); err != nil {
			return nil, err
		}
		if _, dup := r.byName[e.Name]; dup {
			return nil, fmt.Errorf("command: 重名 %q：名册的名字必须唯一，撞名要响亮报（票 213 第 0.3 节选了 Pi 那一形，不是静默忽略）", e.Name)
		}
		door := doorKey(e)
		if owner, taken := r.doors[door]; taken {
			return nil, fmt.Errorf("command: 门 %q 已被 %q 占用，%q 又来认领：一扇门只许一条命令，否则就是把同一份能力换个名字塞进名册", door, owner, e.Name)
		}
		r.byName[e.Name] = len(r.entries)
		r.doors[door] = e.Name
		r.entries = append(r.entries, e)
	}
	return r, nil
}

// Names returns the command names in stable order.
func (r *Roster) Names() []string {
	out := make([]string, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e.Name)
	}
	sort.Strings(out)
	return out
}

// Entries returns a copy of the rows in roster order.
func (r *Roster) Entries() []Command {
	out := make([]Command, len(r.entries))
	copy(out, r.entries)
	return out
}

// Lookup finds one command by its slash-less name.
func (r *Roster) Lookup(name string) (Command, bool) {
	i, ok := r.byName[name]
	if !ok {
		return Command{}, false
	}
	return r.entries[i], true
}

// Judge answers one concrete request under the given conditions. An unknown name,
// an unparseable target and a gate that cannot be satisfied all answer Refused
// with a sentence saying why - the roster never stays silent about a refusal
// (AC#3), and it never invents a verdict the door does not share.
func (r *Roster) Judge(env Env, req Request) Availability {
	i, ok := r.byName[req.Name]
	if !ok {
		return Availability{Verdict: Refused, Reason: fmt.Sprintf(
			"名册里没有 %q（今天只有 %s）：没有能力的东西不进名册，票 213 AC#2", req.Name, strings.Join(r.Names(), ", "))}
	}
	e := r.entries[i]
	switch e.Channel {
	case ChannelControlWord:
		return Availability{Verdict: Allowed, Reason: fmt.Sprintf(
			"命令 %q 走 D11(1) 控制词 %q：agent.MatchControl 解析、(*Loop).runControl 执行，本机直接停当前任务，不经 LLM，也不要人批", e.Name, e.Utterance)}
	case ChannelPanelMode:
		return judgeMode(env, e, req)
	}
	// A channel value New would have refused; fail closed rather than guess.
	return Availability{Verdict: Refused, Reason: fmt.Sprintf("命令 %q 的门 %q 不认识，按 fail-closed 拒绝", e.Name, e.Channel)}
}

// judgeMode judges one panel.mode.request against the same predicate the two
// production gate layers call (perm.JudgeCard), so the roster cannot disagree with
// the door about which switch needs a card. The `from` it passes is Env.Mode, the
// EFFECTIVE档, because that is the档 the panel handler judges against; a legless
// store answers the strictest档 there even when the file says otherwise (ticket
// 372-r2), and R9's direction rule must be read on the same pair the door reads.
func judgeMode(env Env, e Command, req Request) Availability {
	target := strings.TrimSpace(req.Target)
	if target == "" {
		return Availability{Verdict: Refused, Reason: fmt.Sprintf(
			"命令 %q 要一个目标档位（%s）：没说切到哪里的档位请求在写腿那儿也只是解析错误", e.Name, strings.Join(risk.ModeNames(), " / "))}
	}
	to, err := risk.ParseMode(target)
	if err != nil {
		return Availability{Verdict: Refused, Reason: fmt.Sprintf(
			"命令 %q 的目标档位 %q 不可解析：risk.ParseMode 对未知拼写返回错误而不是回退（票 83 的 lying-key 规矩）", e.Name, target)}
	}
	if to == env.Mode {
		return Availability{Verdict: Allowed, Reason: fmt.Sprintf(
			"目标就是当前档 %q：没有变化，不花一张卡（perm.Store.Set 的 no-change 分支）", to)}
	}
	switch perm.JudgeCard(env.Mode, to, env.CardChannel) {
	case perm.GateNoCard:
		return Availability{Verdict: Allowed, Reason: fmt.Sprintf(
			"从 %q 切到 %q 是调严，不出卡（R9：只有调松要确认）；perm.JudgeCard(%q, %q, card=%v) 报 GateNoCard，写腿直接落盘并写审计",
			env.Mode, to, env.Mode, to, env.CardChannel)}
	case perm.GateNoChannel:
		return Availability{Verdict: Refused, Reason: fmt.Sprintf(
			"从 %q 切到 %q 是调松，而本机没有接 L2 确认腿：perm.JudgeCard 报 GateNoChannel，面板腿在写之前就拒（panel.ErrNoL2Confirm），perm.Store.Set 也拒（失败闭合）", env.Mode, to)}
	case perm.GateAsk:
		return Availability{Verdict: NeedsApproval, Reason: fmt.Sprintf(
			"从 %q 切到 %q 是调松：perm.JudgeCard 报 GateAsk，perm.Store.Set 要一张已应答的 L2 强确认卡；只连了腿不算，卡没被答就是拒（perm.ResultCardNotAnswered）", env.Mode, to)}
	}
	return Availability{Verdict: Refused, Reason: fmt.Sprintf("perm.JudgeCard 对 %q -> %q 给出未知判决，按 fail-closed 拒绝", env.Mode, to)}
}

// Production is the assembled roster: the commands this machine can actually
// execute today. Its production caller is cmd/wisp's `wisp commands` leg, which
// prints these rows and nothing else; a caller that kept its own copy of the table
// would make this function decorative, so roster_test.go re-reads the disk for the
// import and the call rather than trusting the sentence.
func Production() (*Roster, error) { return New(productionEntries()...) }

// productionEntries is the only table. Each row is proved executable by New (the
// door resolves through the same parser the door uses) and by roster_test.go (each
// door is driven for real).
func productionEntries() []Command {
	return []Command{
		{
			Name:      "cancel",
			Summary:   "取消当前正在跑的任务：停住它，并且放弃它还没做完的效果",
			Approval:  ApprovalNever,
			Channel:   ChannelControlWord,
			Utterance: "取消",
		},
		{
			Name:      "stop",
			Summary:   "停住当前正在跑的任务：立刻停止输出与后续工具调用",
			Approval:  ApprovalNever,
			Channel:   ChannelControlWord,
			Utterance: "停",
		},
		{
			Name:     "mode",
			Summary:  "切换权限预设：默认最严、高风险才问、全自动",
			Args:     "目标档位：" + strings.Join(risk.ModeNames(), " / "),
			Approval: ApprovalPerRequest,
			Channel:  ChannelPanelMode,
			Method:   panel.MethodModeRequest,
		},
	}
}

// verify is AC#2's gate, row by row.
func verify(e Command) error {
	if err := verifyName(e.Name); err != nil {
		return err
	}
	if strings.TrimSpace(e.Summary) == "" {
		return fmt.Errorf("command: %q 没有一句话说明：菜单里一条说不出话的行就是\"点了没反应\"的同族，票 213 交付 1", e.Name)
	}
	switch e.Channel {
	case ChannelControlWord:
		if strings.TrimSpace(e.Utterance) == "" {
			return fmt.Errorf("command: %q 声称走控制词门，却没有给控制词（D11 的词表是唯一来源）", e.Name)
		}
		if e.Method != "" {
			return fmt.Errorf("command: %q 同时带了控制词和面板方法：一条命令只许一扇门", e.Name)
		}
		if _, ok := agent.MatchControl(e.Utterance); !ok {
			return fmt.Errorf("command: %q 的控制词 %q 不在 agent.MatchControl 的词表里：没有解析器就没有执行者，不许凭空空转（票 213 AC#2）", e.Name, e.Utterance)
		}
	case ChannelPanelMode:
		if strings.TrimSpace(e.Method) == "" {
			return fmt.Errorf("command: %q 声称走面板门，却没有给方法名", e.Name)
		}
		if e.Utterance != "" {
			return fmt.Errorf("command: %q 同时带了控制词和面板方法：一条命令只许一扇门", e.Name)
		}
		// ParseComposerRequest is the closed set in bridge.go and the source
		// check; asking it is what makes "this door exists" a measurement instead
		// of a claim, and it is why no new C17 method name can enter through this
		// table.
		if _, err := panel.ParseComposerRequest(verifyEnvelope(e.Method)); err != nil {
			return fmt.Errorf("command: %q 的方法 %q 不在面板 composer 名册里（panel.ParseComposerRequest 拒收：%v）：名册里的一条命令不再是新能力", e.Name, e.Method, err)
		}
	case "":
		return fmt.Errorf("command: %q 没有执行门（channel 为空）：从没有能力的那一条不进名册（票 213 AC#2、交付 1）", e.Name)
	default:
		return fmt.Errorf("command: %q 的门 %q 不认识（只认 %q / %q）", e.Name, e.Channel, ChannelControlWord, ChannelPanelMode)
	}
	return nil
}

// verifyEnvelope is the smallest envelope ParseComposerRequest accepts for a
// method: the wall checks method membership, source identity and a requestId, and
// nothing method-specific reaches the parse.
func verifyEnvelope(method string) string {
	return fmt.Sprintf(`{"method":%q,"requestId":"roster-verify","source":%q}`, method, panel.ComposerRequestSource)
}

// doorKey names one door so two commands cannot share it.
func doorKey(e Command) string {
	if e.Channel == ChannelControlWord {
		return string(e.Channel) + ":" + e.Utterance
	}
	return string(e.Channel) + ":" + e.Method
}

// verifyName enforces the slash-less spelling: lowercase ascii, digits and
// hyphens after the first letter, no dots. A dotted name is a route name, and
// letting one in would make the roster look like a second route table.
func verifyName(name string) error {
	bad := func() error {
		return fmt.Errorf("command: 名字 %q 不合法：不带斜杠、不带点、只许小写字母/数字/连字符，且以字母开头（形如 \"mode\"）", name)
	}
	if name == "" || len(name) > 32 {
		return bad()
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z':
		case i > 0 && c >= '0' && c <= '9':
		case i > 0 && c == '-':
		default:
			return bad()
		}
	}
	if name[len(name)-1] == '-' {
		return bad()
	}
	return nil
}
