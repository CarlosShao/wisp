package panel

// The snapshot pump: the first place production BUILDS the packet the panel is
// rendered from (ticket 35's data half).
//
// WHAT WAS TRUE BEFORE THIS FILE. Snapshot (composer.go:44) has four fields and
// NewSnapshot (composer.go:63) had zero non-test callers, so "the panel displays
// authority and requests it" (composer.go:6-12) described a packet nobody sent.
// ticket 145's census measured the same thing from the other side and named the
// disease in composer.go:41-43: "a view model that only a test can build is a
// view model production never sends". Moving Snapshot into production code did
// not move its callers, so this file's job is exactly that move: constructors
// called from a running process, over values read out of live objects.
//
// WHAT THIS FILE IS NOT: the transport. There is no Go -> page channel in this
// tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer,
// no local HTTP/SSE/websocket server, and assets.go:8 quotes D29 saying there
// must not be one. docs/evidence/s1/35-panel-snapshot-pump-r1.md §1.2 measures
// all four. So this pump stops at the byte boundary: Marshal hands out the JSON,
// Publish hands it to src.Out, and src.Out is the one line a later ticket
// re-points at the host. The last mile is NOT here, and it was not invented
// here either - a file the renderer polls would be a second channel, which
// composer.go:14 ("no second channel") and D29 both refuse.
//
// NO NEW KEYS. Snapshot keeps its four JSON keys and frontend/src/lib/panel.ts
// is untouched, so TestComposerContractTypesMatchFrontend (composer_test.go:48)
// stays the two-way ruler it is. Adding the fifth key - the "which screen" one -
// is Q-51 territory and ticket 145 AC#3 already answered it: not in this slice.
// Nor does this file reach for the six fields ticket 145 marks as having no
// vocabulary in the repository (thinkingMs / reasoningMs / durationMs / humanText
// / fragment / IconClass): a pump that fills a field with a constant to look
// complete is the same failure this file exists to end.

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// NativeVerdict is one live approval as the host reads it off its own queue.
// Every field is a copy of a verdict internal/risk already reached and
// internal/agent/approval already admitted; nothing here re-assesses, because a
// second assessment is a second judge and the second one would be this file.
//
// Args is the argv-shaped view of the call's arguments (ApprovalCardView.args),
// and CallChain has no producer in this repository yet - see the field comment.
type NativeVerdict struct {
	CorrelationID string
	Tool          string
	Args          []string
	// CallChain is what the host actually recorded about how this call was
	// reached. The run path records none (tools.Decision has no chain field), so
	// production hands in nil and the card renders an empty list. That is an
	// honest empty, not a filler: inventing ["session","agent-loop",tool] here
	// would put a string no code produced onto a security card.
	CallChain              []string
	Level                  risk.Level
	RulesHit               []risk.RuleID
	Reason                 string
	SessionOverrideBlocked bool
}

// CardView renders one live verdict as the card the panel shows, through the
// same pure half NewApprovalCardView uses (approval.go:72) - which is the point:
// the card in a pushed snapshot and the card the CLI diagnostic prints are the
// same function over the same verdict, so they cannot disagree about what L2
// meant.
func (v NativeVerdict) CardView() ApprovalCardView {
	return CardViewFromDecision(
		ApprovalSubject{
			CorrelationID: v.CorrelationID,
			Tool:          v.Tool,
			Args:          v.Args,
			CallChain:     v.CallChain,
		},
		risk.Decision{
			Level:                  v.Level,
			RulesHit:               v.RulesHit,
			Reason:                 v.Reason,
			SessionOverrideBlocked: v.SessionOverrideBlocked,
		},
	)
}

// PumpSources is the live native state a snapshot is assembled from. Each field
// is a READER, not a value, because the only version of this that is worth
// anything reads what the running process is actually doing. A nil reader is a
// statement about the assembly, and each section answers it in its own documented
// direction (see Snapshot): the mode and the workspace refuse to invent, the
// queue and the stream report what they hold.
type PumpSources struct {
	// Verdicts reads the approvals pending right now, head first.
	Verdicts func() []NativeVerdict
	// Mode reads the permission mode in force (perm.Store.PermissionMode).
	Mode func() risk.Mode
	// Workspace reads the narrowed root as native resolved it.
	Workspace func() WorkspaceView
	// Results reads the streamed assistant text so far.
	Results func() []ResultChunk
	// AttachmentMax is the composer's per-attachment ceiling. Zero means
	// MaxAttachmentBytes, the constant the broker itself defaults to
	// (attachments.go:186), so the packet never advertises a ceiling that the
	// store would not apply.
	AttachmentMax int64
	// Now is the clock for the generatedAt stamp. nil is time.Now. This is a
	// display stamp and nothing else (panel.ts:132-134: "display only, never
	// used to derive state"), so it is not a deadline and not a timeout.
	Now func() time.Time
	// Out is the far end of the exit. Today the only production assembly attaches
	// the run's persistent ledger; a later ticket re-points it at the host. nil
	// means "assembled but nowhere to send", which Publish reports as an error
	// instead of pretending the bytes went somewhere.
	Out func(snapshot Snapshot, data []byte) error
}

// SnapshotPump assembles and emits the panel's packet.
type SnapshotPump struct {
	src PumpSources

	mu        sync.Mutex
	publishes int
}

// NewSnapshotPump wires a pump to its sources. The result is never nil; a zero
// PumpSources is legal and produces a snapshot whose composer section says
// "unknown" and whose lists are empty - a packet that states what it does not
// know, which is the AC#4 rule this file inherits.
func NewSnapshotPump(src PumpSources) *SnapshotPump {
	return &SnapshotPump{src: src}
}

// Snapshot builds the whole truth as of this instant.
func (p *SnapshotPump) Snapshot() Snapshot {
	if p == nil {
		return NewSnapshot(nil, nil, ComposerState{}, time.Time{})
	}
	var pending []NativeVerdict
	if p.src.Verdicts != nil {
		pending = p.src.Verdicts()
	}
	cards := make([]ApprovalCardView, 0, len(pending))
	for _, v := range pending {
		cards = append(cards, v.CardView())
	}

	var results []ResultChunk
	if p.src.Results != nil {
		results = p.src.Results()
	}

	// The composer section is built by NewComposerState, which is the constructor
	// ticket 92 made mandatory and nobody ever called: the accepted MIME list and
	// the attachment ceiling the renderer reads come from there, so the packet
	// cannot advertise one ceiling while the broker applies another.
	mode, modeReadable := risk.Mode(0), false
	if p.src.Mode != nil {
		mode = p.src.Mode()
		modeReadable = mode.Valid()
	}
	workspace := UnsetWorkspaceView()
	if p.src.Workspace != nil {
		workspace = p.src.Workspace()
	}
	maxAttachment := p.src.AttachmentMax
	if maxAttachment == 0 {
		maxAttachment = MaxAttachmentBytes
	}
	composer := NewComposerState(mode, workspace, nil, maxAttachment)
	if !modeReadable {
		// AC#4 through the pump: a mode nobody read is reported as unknown, and
		// never as the strictest-sounding or the safest-sounding default.
		composer.Mode = ModeUnknownView()
	}

	now := time.Now
	if p.src.Now != nil {
		now = p.src.Now
	}
	return NewSnapshot(cards, results, composer, now())
}

// Marshal is the exit function: the bytes the panel would be rendered from.
// It is a separate step from Snapshot so the last mile has something to hand to
// a transport without re-reading the queue (and racing whatever answered in
// between). It owns no state, so two pumps - or two goroutines on one pump -
// cannot see each other's packet.
func (p *SnapshotPump) Marshal() ([]byte, error) {
	data, err := json.Marshal(p.Snapshot())
	if err != nil {
		return nil, fmt.Errorf("panel: snapshot not encodable, nothing was sent: %w", err)
	}
	return data, nil
}

// Publish assembles, encodes, and hands the bytes to src.Out. The returned
// snapshot is the one those bytes describe, so a caller can audit the depth it
// just booked without asking the queue a second time.
func (p *SnapshotPump) Publish() (Snapshot, []byte, error) {
	if p == nil {
		return Snapshot{}, nil, fmt.Errorf("panel: no pump attached, the snapshot was not built")
	}
	snap := p.Snapshot()
	data, err := json.Marshal(snap)
	if err != nil {
		p.count()
		return snap, nil, fmt.Errorf("panel: snapshot not encodable, nothing was sent: %w", err)
	}
	if p.src.Out == nil {
		p.count()
		return snap, data, fmt.Errorf("panel: 快照已构造但没有出口（最后一公里在票 33/35 名下），本次未送达")
	}
	if err := p.src.Out(snap, data); err != nil {
		p.count()
		return snap, data, fmt.Errorf("panel: 快照出口拒收：%w", err)
	}
	p.count()
	return snap, data, nil
}

// Publishes reports how many snapshots this pump has put through its exit,
// including the ones the exit refused. A count that never moves is how a caller
// notices the pump is assembled but never driven.
func (p *SnapshotPump) Publishes() int {
	if p == nil {
		return 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.publishes
}

func (p *SnapshotPump) count() {
	p.mu.Lock()
	p.publishes++
	p.mu.Unlock()
}

// StreamLog produces the snapshot's results section from the deltas a run
// streams, and it is the reason "results" has a source at all.
//
// One chunk per correlation key, deltas appended to it, Close marking that key
// done. The key the loop actually has is the task id (agent.Event carries
// TaskID, and Event has no per-chunk id of its own) - naming that here is the
// honest version of a field the pump must fill; a made-up id would be the packet
// lying about its own provenance.
//
// Overflow MERGES and never drops: at maxKeys the two oldest chunks are
// concatenated, which is what a stream looks like when you stop splitting it.
// That is the small end of ticket 35 AC#5's backpressure rule (the queue-level
// merge belongs to the transport), and it is here because a pump that grows
// without bound is not a pump.
type StreamLog struct {
	mu      sync.Mutex
	order   []string
	chunks  map[string]ResultChunk
	maxKeys int
}

// DefaultStreamKeys is how many concurrent result streams a pump tracks before
// it starts merging. 32 is generous against the only producer today (one task
// per `wisp run`) and small enough to bound a long-lived resident process.
const DefaultStreamKeys = 32

// NewStreamLog builds an empty results log. maxKeys <= 0 means DefaultStreamKeys.
func NewStreamLog(maxKeys int) *StreamLog {
	if maxKeys <= 0 {
		maxKeys = DefaultStreamKeys
	}
	return &StreamLog{chunks: map[string]ResultChunk{}, maxKeys: maxKeys}
}

// Append adds one delta to its key's chunk.
func (s *StreamLog) Append(key, text string) {
	if s == nil || text == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.chunks[key]
	if !ok {
		s.chunks[key] = ResultChunk{CorrelationID: key, Text: text}
		s.order = append(s.order, key)
		s.mergeOverflowLocked()
		return
	}
	cur.Text += text
	s.chunks[key] = cur
}

// Close marks one key's chunk as the end of its stream. An unknown key still
// records a done chunk, because "the stream ended and nobody had noticed it
// start" is a fact the panel should see rather than lose.
func (s *StreamLog) Close(key string) {
	if s == nil || key == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.chunks[key]
	if !ok {
		s.chunks[key] = ResultChunk{CorrelationID: key, Done: true}
		s.order = append(s.order, key)
		s.mergeOverflowLocked()
		return
	}
	cur.Done = true
	s.chunks[key] = cur
}

// Chunks returns the accumulated result chunks in arrival order.
func (s *StreamLog) Chunks() []ResultChunk {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ResultChunk, 0, len(s.order))
	for _, k := range s.order {
		out = append(out, s.chunks[k])
	}
	return out
}

// mergeOverflowLocked folds the oldest chunk into the one behind it until the
// log is back inside its bound. The fold keeps BOTH texts, so what the panel
// loses is a split, never a sentence. Caller holds s.mu.
func (s *StreamLog) mergeOverflowLocked() {
	for len(s.order) > s.maxKeys && len(s.order) >= 2 {
		oldest, next := s.order[0], s.order[1]
		merged := s.chunks[next]
		merged.Text = s.chunks[oldest].Text + merged.Text
		merged.Done = merged.Done || s.chunks[oldest].Done
		s.chunks[next] = merged
		delete(s.chunks, oldest)
		copy(s.order, s.order[1:])
		s.order = s.order[:len(s.order)-1]
	}
}
