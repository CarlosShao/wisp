package risk

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// C25 Provenance — taint marking, coarse fragment matching, and the six
// exfiltration channels (SPEC-06 §5, D33/F4, D30①, 16.9#1).
//
// Model: every sensitive-source read is Mark()ed under a task/session scope
// (the DisposalScope lifetime — scopes are opened by the composition root and
// Closed on Dispose, so a new session never inherits old taints). Before any
// call executes, Inspect() checks the outgoing parameters against every
// tainted fragment in that scope; a >=8-char normalized contiguous fragment
// hit is the R4 upgrade to L2, and the verdict names the source.
//
// Exfil channels (all six of D33/F4 — never only the HTTP body):
//
//	1  web.search query string
//	2  notify text + URL
//	3  TTS announce text (physical exfil: audible to the room)
//	4  clipboard.write
//	5  fs.write INTO A SYNC DIR (P12 detection; see syncdirs.go)
//	6  HTTP POST body / URL (net)
//
// Tools with a frozen D34 name (web.search, notify, clipboard.write, fs.write)
// get their table keys LABELED with the contract channel; every other
// parameter value (unknown key names, nested objects, arrays, byte slices) is
// still scanned fail-closed under ChUnknown, so a renamed or wrapped parameter
// cannot narrow the scan. TTS announce and the net body do not have a frozen
// tool name yet (tickets 22/26 name them), so they call
// CheckText(scope, ChTTS|ChHTTP, text) directly; and any call whose tool is
// NOT in the channel table fail-closes to a conservative scan of every string
// parameter.
//
// The ONE landing-site conditional channel is ⑤ (fs.write into a sync dir);
// channel ⑥ (HTTP POST body/URL) has no landing condition. Because a call
// without a frozen tool name carries no authority in its key names, the write
// gate is therefore decided by the SHAPE of the call (see writeGate), never by
// what the caller named its payload — the previous "skip the keys called
// content/data" rule let `{url, path, data}` through while `{url, path, body}`
// was caught (adversarial re-verification M-7).
//
// Wiring to the assessor: the frozen C19 seam is TaintDetector.TaintHit(params)
// — it carries no tool name — so callers bind it per scope:
// assessor.WithTaintDetector(prov.Detector(taskScopeID)) yields the
// conservative shape-based verdict; loop code that knows the tool name should
// prefer Inspect()/CheckText() directly (wiring lands with tickets 20/21/22/26).
//
// DEFERRED(C25-loop-wiring): the agent loop / tool providers (tickets 10, 20,
// 21, 22, 26) must (1) OpenScope at task start and Defer(CloseScope) on the
// task's DisposalScope, (2) call Mark(...) on every SPEC-06 §5 sensitive
// source output, (3) gate outgoing calls with Inspect(scope, tool, params)
// (or CheckText for the TTS/HTTP-body channels), and (4) record the R4 source
// name in the tool_call forensics row via the existing DAO — Hit.Fragment
// itself must NEVER be persisted (it is sensitive content; native card only).
//
// Deliberate limits (accepted residual risks, recorded in docs/PRECHECK.md):
//   - LLM paraphrase/translation of tainted content evades fragment matching
//     (16.9#1 rejected exact taint tracking; D30 layers 2/3/5 back-stop).
//   - screen.capture marks provenance but image bytes have no defined text
//     normalization; matching applies to textual renderings only.
//   - Sources shorter than the fragment floor (e.g. a 7-char secret) cannot
//     produce a contract-shaped match.

// Sensitive source tools whose output is marked (SPEC-06 §5). system.get is
// the sysinfo focused-window-title source; asr.transcript is the user
// transcript hook promised to ticket 15.
const (
	SrcFSRead        = "fs.read" // incl. out-of-allowlist attempts
	SrcSearchContent = "search.content"
	SrcClipboardRead = "clipboard.read"
	SrcSystemGet     = "system.get" // focused-window title (high-value leak)
	SrcWebFetch      = "web.fetch"
	SrcDocRead       = "doc.read"
	SrcScreenCapture = "screen.capture"
	SrcTranscript    = "asr.transcript"
)

// sensitiveSourceTools is the SPEC-06 §5 set (marker-side only; a Mark() for
// any other tool is still honored fail-closed, just logged).
var sensitiveSourceTools = []string{
	SrcFSRead, SrcSearchContent, SrcClipboardRead, SrcSystemGet,
	SrcWebFetch, SrcDocRead, SrcScreenCapture, SrcTranscript,
}

// IsSensitiveSource reports whether a tool is in the contract marking set.
func IsSensitiveSource(tool string) bool {
	for _, t := range sensitiveSourceTools {
		if t == tool {
			return true
		}
	}
	return false
}

// Channel names the six exfil channels (identity strings land in Hit and in
// the tool_call forensics trail; they are display/audit data, not new rules).
type Channel string

const (
	ChWebSearch Channel = "web.search.query"
	ChNotify    Channel = "notify"
	ChTTS       Channel = "tts.announce"
	ChClipboard Channel = "clipboard.write"
	ChSyncWrite Channel = "fs.write.syncdir"
	ChHTTP      Channel = "net.http.body"
	ChUnknown   Channel = "unclassified.exfil-scan" // fail-closed generic scan
)

// channelTable maps the frozen D34 tool names to the parameter keys that
// carry exfil-able text (SPEC-06 §5 channel wording). Composition roots may
// extend/override via ProvOptions.ChannelKeys; entries must never be removed
// (removal = weakening a contract channel).
var channelTable = map[string]struct {
	ch   Channel
	keys []string
}{
	"web.search":      {ChWebSearch, []string{"query"}},
	"notify":          {ChNotify, []string{"text", "url"}},
	"clipboard.write": {ChClipboard, []string{"text", "content"}},
	"fs.write":        {ChSyncWrite, nil}, // special: content only when path is a sync dir
}

// writeChannelKeys are the payload parameters of the FROZEN fs.write contract,
// scanned ONLY when the call's path resolves into a sync dir (fs.write to a
// normal dir stays a plain L1 write; SPEC-06 §5 makes only the sync-dir
// landing an exfil channel). They are consulted for the named fs.write channel
// only — the contract gives those two names a meaning; for a call without a
// frozen name no key name grants an exemption (M-7), the exemption is decided
// by writeGate from the shape of the call.
var writeChannelKeys = map[string]bool{"content": true, "data": true}

// pathKeys are the parameters naming the write target for the sync-dir test.
var pathKeys = []string{"path", "file", "filepath", "dest", "destination"}

// remoteSinkKeys name a parameter that addresses a REMOTE destination (channel
// ⑥'s own wording: HTTP POST body/URL). A call carrying one is never exempted
// from the payload scan, whichever key the body arrived in.
var remoteSinkKeys = map[string]bool{
	"url": true, "uri": true, "href": true, "link": true, "endpoint": true,
	"webhook": true, "callback": true, "redirect": true, "to": true, "recipient": true,
}

// remoteSinkSuffixes catch the *_url / *_endpoint spellings of the same thing.
var remoteSinkSuffixes = []string{"url", "uri", "endpoint", "webhook", "href"}

// remoteSinkSchemes identify a remote sink by VALUE, so renaming the key cannot
// close the gate either: a parameter whose text is an address means content in
// the same call can leave the machine.
var remoteSinkSchemes = []string{"http://", "https://", "ftp://", "ftps://", "sftp://",
	"ws://", "wss://", "s3://", "webdav://", "dav://", "file://", "mailto:"}

// contract defaults for the scanning budgets (overridable via ProvOptions;
// zero/negative means "use the default", never "scan nothing").
const (
	// defaultMaxParamDepth bounds the generic parameter walk. Nesting deeper
	// than this fails closed as an unscanned-nesting hit instead of silently
	// skipping the tail (M-4).
	defaultMaxParamDepth = 8
	// defaultMaxScopeSources is the D32 per-scope index budget.
	defaultMaxScopeSources = 64
	// defaultMaxSourceRunes caps one stored tainted source.
	defaultMaxSourceRunes = 262144
	// defaultMaxScanRunes caps one scanned parameter value.
	defaultMaxScanRunes = 2097152
)

// ProvOptions configures the engine. Everything is optional; zero values take
// the contract defaults (never a weaker behavior — see field docs).
type ProvOptions struct {
	// MinFragmentChars sets the matching window. The contract floor is 8;
	// values >=8 or <=0 are clamped to 8 (weakening is impossible), values
	// in (0,8) are accepted as STRICTER (用户只能调严, SPEC-06 §2).
	MinFragmentChars int
	// MaxSourceRunes caps one stored source (default 262144). Truncation is
	// logged; fragments past the cap can be missed (documented residual).
	MaxSourceRunes int
	// MaxScanRunes bounds how much of one parameter value is scanned
	// (default 2097152 runes); beyond is logged-not-scanned (documented).
	MaxScanRunes int
	// MaxParamDepth bounds the nesting the generic parameter walk descends
	// (default defaultMaxParamDepth). Reaching it is NOT a silent miss: the
	// call fails closed as an unscanned-nesting hit (adversarial report M-4).
	// Values below the default are accepted as looser scanning budgets only
	// if the caller also accepts the fail-closed escalation they trigger.
	MaxParamDepth int
	// MaxScopeSources is the D32 memory budget for one scope's indexed
	// sources (default 64). Exceeding it is logged and every taint is still
	// kept — dropping marks would open the gate (fail-open), which the budget
	// exists to bound, not to cause.
	MaxScopeSources int
	// ReparseExceptions mirrors [fs] reparse_point_exceptions for C26
	// Resolve (sync-root membership is decided by Resolve, never by hand).
	ReparseExceptions []string
	// SyncRoots injects configured sync roots (user escape hatch for
	// relocated clients); SyncFixturePath loads the JSON fixture fallback.
	SyncRoots       []SyncRoot
	SyncFixturePath string
	// ChannelKeys overrides/extends the exfil parameter mapping per tool.
	// Additive only: merged over channelTable, keys appended.
	ChannelKeys map[string][]string
	// HomeDir/LocalAppData/AppData relocate the P12 probes (tests); empty
	// falls back to the live environment.
	HomeDir      string
	LocalAppData string
	AppData      string
	// NoProbe skips registry/config/default probing entirely (pure unit
	// tests); fixture/injected roots still apply.
	NoProbe bool
}

// taintMark is one immutable marked source inside a scope.
type taintMark struct {
	tool   string
	origin string
	idx    *fragmentIndex
	at     int64
}

// Provenance is the C25 engine. Construct with NewProvenance. It is safe for
// concurrent use.
type Provenance struct {
	mu       sync.RWMutex
	scopes   map[string][]*taintMark
	minChars int
	maxSrc   int
	maxScan  int
	maxDepth int
	maxMarks int
	channel  map[string][]string
	sync     *syncSet
}

// Hit describes one R4 taint finding.
type Hit struct {
	ScopeID  string
	Channel  Channel
	SrcTool  string // the sensitive SOURCE tool (e.g. web.fetch)
	Origin   string // its origin (path / URL / window title...)
	Fragment string // the matched fragment (native confirmation card only —
	// NEVER write this into the DB/args_json: the fragment IS sensitive text)
}

// Source renders the "<源>" the confirmation card must name
// ("本次操作包含来自 `<源>` 的内容", SPEC-06 §5 / R4 reason).
func (h Hit) Source() string {
	if strings.TrimSpace(h.Origin) != "" {
		return h.SrcTool + " " + h.Origin
	}
	return h.SrcTool
}

// NewProvenance resolves options (clamping toward the stricter side) and runs
// the P12 sync-directory detection.
func NewProvenance(o ProvOptions) *Provenance {
	minChars := o.MinFragmentChars
	if minChars <= 0 || minChars >= contractMinFragmentChars {
		minChars = contractMinFragmentChars
	}
	if minChars < 1 {
		minChars = contractMinFragmentChars
	}
	maxSrc := o.MaxSourceRunes
	if maxSrc <= 0 {
		maxSrc = defaultMaxSourceRunes
	}
	maxScan := o.MaxScanRunes
	if maxScan <= 0 {
		maxScan = defaultMaxScanRunes
	}
	maxDepth := o.MaxParamDepth
	if maxDepth <= 0 {
		maxDepth = defaultMaxParamDepth
	}
	maxMarks := o.MaxScopeSources
	if maxMarks <= 0 {
		maxMarks = defaultMaxScopeSources
	}
	ch := map[string][]string{}
	for t, d := range channelTable {
		ch[t] = d.keys
	}
	for t, keys := range o.ChannelKeys {
		existing := map[string]bool{}
		var merged []string
		for _, k := range ch[t] {
			existing[k] = true
			merged = append(merged, k)
		}
		for _, k := range keys {
			if !existing[k] {
				merged = append(merged, k)
			}
		}
		ch[t] = merged
	}
	p := &Provenance{
		scopes:   map[string][]*taintMark{},
		minChars: minChars,
		maxSrc:   maxSrc,
		maxScan:  maxScan,
		maxDepth: maxDepth,
		maxMarks: maxMarks,
		channel:  ch,
	}
	if o.NoProbe {
		s := &syncSet{home: normPath(orDefault(o.HomeDir, userHomeDir())), excs: o.ReparseExceptions}
		for _, r := range o.SyncRoots {
			s.add(r)
		}
		if o.SyncFixturePath != "" {
			for _, r := range fixtureRoots(o.SyncFixturePath) {
				s.add(r)
			}
		}
		s.finalize()
		p.sync = s
	} else {
		p.sync = detectSyncSet(o)
	}
	return p
}

// --- scope lifetime (DisposalScope binding) ----------------------------------

// OpenScope registers a task/session scope (composition root calls this at
// task start; idempotent).
func (p *Provenance) OpenScope(scopeID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.scopes[scopeID]; !ok {
		p.scopes[scopeID] = nil
	}
}

// CloseScope drops every taint bound to a scope (wire via
// disposalScope.Defer(p.CloseScope) — disposal MUST not leak taints past the
// session, AC: new sessions never inherit).
func (p *Provenance) CloseScope(scopeID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.scopes, scopeID)
}

// Mark records one sensitive-source read in a scope. Provenance identity is
// {tool, origin, sensitive} (C25 contract); content is normalized + indexed
// for >=8-char contiguous fragment matching. Marking a closed/never-opened
// scope still records the taint (fail-closed on the marking side) and logs.
// Returns false only when the content has no matchable characters after
// normalization.
func (p *Provenance) Mark(scopeID, tool, origin, content string) bool {
	norm := normalizeTaint(content)
	if norm == "" {
		logf("risk/C25: mark %s origin=%q skipped: empty after normalization", tool, origin)
		return false
	}
	rp := []rune(norm)
	if len(rp) < p.minChars {
		// Shorter than the contract floor: no >=8 fragment can exist.
		logf("risk/C25: mark %s origin=%q has only %d normalized chars (< floor %d): indexed anyway as shorter fragments are impossible (documented residual)", tool, origin, len(rp), p.minChars)
	}
	truncated := false
	if len(rp) > p.maxSrc {
		rp, truncated = rp[:p.maxSrc], true
	}
	if !IsSensitiveSource(tool) {
		logf("risk/C25: Mark tool %q is not a SPEC-06 §5 source; recorded fail-closed as sensitive anyway (origin=%q)", tool, origin)
	}
	m := &taintMark{
		tool: tool, origin: origin,
		idx: newFragmentIndex(string(rp), p.minChars),
		at:  observe.NowWallUTC().Unix(),
	}
	p.mu.Lock()
	if _, ok := p.scopes[scopeID]; !ok {
		logf("risk/C25: scope %q not open (CloseScope raced or composition gap); taint stored — it cannot leak into other scopes", scopeID)
	}
	if n := len(p.scopes[scopeID]); n >= p.maxMarks {
		// D32 memory budget exceeded: keep the taint (dropping it would be a
		// fail-open) but make the budget miss visible for the SLO counters.
		logf("risk/C25: scope %q holds %d tainted sources (MaxScopeSources=%d, D32 budget); keeping every taint — index memory grows past the budget, see docs/PRECHECK.md", scopeID, n, p.maxMarks)
	}
	p.scopes[scopeID] = append(p.scopes[scopeID], m)
	p.mu.Unlock()
	if truncated {
		logf("risk/C25: source %s origin=%q truncated to MaxSourceRunes=%d; fragments past the cap can be missed (documented residual)", tool, origin, p.maxSrc)
	}
	return true
}

// TaintInfo is the introspection view for the panel Security page (ticket 40).
type TaintInfo struct {
	ScopeID  string
	Tool     string
	Origin   string
	RuneLen  int
	MarkedAt int64
}

// ScopeTaints lists a scope's marked sources (never the content itself).
func (p *Provenance) ScopeTaints(scopeID string) []TaintInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var out []TaintInfo
	for _, m := range p.scopes[scopeID] {
		out = append(out, TaintInfo{ScopeID: scopeID, Tool: m.tool, Origin: m.origin,
			RuneLen: m.idx.nrunes, MarkedAt: m.at})
	}
	return out
}

// --- detection ---------------------------------------------------------------

// Fail-closed source names for hits that are NOT content matches. They are
// deliberately not tool names: the confirmation card then names the missing
// information ("包含来自 unbound-scope …"), and the forensics row shows why
// the extra confirm happened. Repo rule: missing information produces a DENY.
const (
	SrcUnboundScope     = "unbound-scope"     // Inspect on a scope that was never opened / already closed
	SrcUnscannedNesting = "unscanned-nesting" // parameter nesting deeper than MaxParamDepth
)

// scopeMarks snapshots a scope's taint marks. The second result reports the
// fail-closed condition: the scope id is not registered at all (typo, missing
// OpenScope, or Inspect after CloseScope) while the engine holds taints in
// some scope — exactly the shape of the "wiring got the scope id wrong" bug,
// which must upgrade rather than pass (adversarial report M-1).
func (p *Provenance) scopeMarks(scopeID string) ([]*taintMark, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	marks := p.scopes[scopeID]
	if len(marks) > 0 {
		return marks[:len(marks):len(marks)], false // snapshot (append never mutates the shared prefix)
	}
	if _, known := p.scopes[scopeID]; known {
		return nil, false // opened and empty: this session read nothing tainted
	}
	for _, m := range p.scopes {
		if len(m) > 0 {
			logf("risk/C25: Inspect/CheckText on unregistered scope %q while other scopes hold taints: fail-closed R4 (ensure OpenScope at task start)", scopeID)
			return nil, true
		}
	}
	return nil, false
}

// Inspect checks one outgoing call against the scope's taints. tool is the
// D34 name ("" = unknown, forces the conservative generic scan). Returns the
// first hit in deterministic order (channel keys in table order, remaining
// params sorted, marks oldest-first).
//
// The channel table only LABELS parameters (adversarial report B-2): every
// remaining parameter value — unknown key names, nested objects, arrays, byte
// payloads — is additionally scanned, so renaming `query` to `q` or wrapping
// the payload cannot narrow the scan. The one exception is the sync-dir gate
// (SPEC-06 §5 channel ⑤, the only landing-site conditional channel), which
// writeGate derives from the shape of the call: no payload key name can open or
// close it (adversarial re-verification M-7).
func (p *Provenance) Inspect(scopeID, tool string, params map[string]any) (Hit, bool) {
	marks, unbound := p.scopeMarks(scopeID)
	if unbound {
		return Hit{ScopeID: scopeID, Channel: ChUnknown, SrcTool: SrcUnboundScope,
			Origin: "scope is not open (OpenScope missing or already closed)"}, true
	}
	if len(marks) == 0 {
		return Hit{}, false
	}

	// Build the candidate list: (channel, text) pairs derived from params.
	type cand struct {
		ch   Channel
		text string
	}
	var cands []cand
	truncated := false
	appendCands := func(key string, ch Channel) {
		ss, cut := collectStrings(params[key], p.maxDepth)
		truncated = truncated || cut
		for _, s := range ss {
			cands = append(cands, cand{ch, s})
		}
	}

	// The single conditional channel of SPEC-06 §5 (⑤ fs.write INTO A SYNC DIR)
	// is decided here, from the shape of the call — never from the name the
	// caller gave its payload (adversarial re-verification M-7).
	gateOpen, gateCh := p.writeGate(tool, params)
	if def, ok := channelTable[tool]; ok && def.ch == ChSyncWrite {
		// Named fs.write: the D34 contract gives `content`/`data` the meaning
		// "the bytes going into the file", so those two keys are what the
		// sync-dir rule exempts. Every other key has no defined meaning in
		// fs.write and is scanned anyway — a taint can ride a file name or an
		// extra field (adversarial report B-2). That asymmetry is fail-CLOSED
		// by key name and never fail-open; it is recorded as N-8 in
		// docs/PRECHECK.md §P12.
		for _, k := range sortedKeys(params) {
			if !gateOpen && writeChannelKeys[k] {
				continue // confirmed plain local write: the body is no channel
			}
			ch := ChUnknown
			if writeChannelKeys[k] {
				ch = gateCh
			}
			appendCands(k, ch)
		}
	} else {
		labelled := map[string]bool{}
		if ok {
			for _, k := range p.channel[tool] {
				labelled[k] = true
				if _, has := params[k]; !has {
					continue
				}
				appendCands(k, def.ch)
			}
		}
		for _, k := range sortedKeys(params) {
			if labelled[k] {
				continue
			}
			if !gateOpen && !isPathKey(k) {
				// No frozen name, and the shape of the call is a plain local
				// write: nothing in it can carry content off the machine, so
				// the payload is exempt WHICHEVER key carries it (M-7). The
				// write target itself is still scanned.
				continue
			}
			ch := ChUnknown
			if gateOpen && !isPathKey(k) {
				ch = gateCh
			}
			appendCands(k, ch)
		}
	}
	for _, c := range cands {
		if h, ok := p.matchText(scopeID, marks, c.ch, c.text); ok {
			return h, true
		}
	}
	if truncated {
		// M-4: the scan was cut off by the nesting budget. Model-controlled
		// JSON nests for free, so an unscanned tail is treated as tainted
		// rather than silently passed.
		logf("risk/C25: params of tool %q nest deeper than MaxParamDepth=%d; tail unscanned, fail-closed R4", tool, p.maxDepth)
		return Hit{ScopeID: scopeID, Channel: ChUnknown, SrcTool: SrcUnscannedNesting,
			Origin: fmt.Sprintf("parameters nested deeper than %d levels", p.maxDepth)}, true
	}
	return Hit{}, false
}

// CheckText matches one outgoing string against a scope's taints under an
// explicit channel label — the entry point for channels without a frozen D34
// tool name yet: TTS announce text (ticket 26) and HTTP POST body/url
// (ticket 22 / net layer). An unbound scope fail-closes exactly as in
// Inspect (M-1).
func (p *Provenance) CheckText(scopeID string, ch Channel, text string) (Hit, bool) {
	marks, unbound := p.scopeMarks(scopeID)
	if unbound {
		return Hit{ScopeID: scopeID, Channel: ch, SrcTool: SrcUnboundScope,
			Origin: "scope is not open (OpenScope missing or already closed)"}, true
	}
	if len(marks) == 0 {
		return Hit{}, false
	}
	return p.matchText(scopeID, marks, ch, text)
}

func (p *Provenance) matchText(scopeID string, marks []*taintMark, ch Channel, text string) (Hit, bool) {
	norm := normalizeTaint(text)
	rp := []rune(norm)
	if len(rp) > p.maxScan {
		logf("risk/C25: exfil candidate truncated to MaxScanRunes=%d; the tail is not scanned (documented residual)", p.maxScan)
		rp = rp[:p.maxScan]
		norm = string(rp)
	}
	for _, m := range marks {
		if frag, ok := m.idx.contains(norm); ok {
			return Hit{ScopeID: scopeID, Channel: ch, SrcTool: m.tool, Origin: m.origin, Fragment: frag}, true
		}
	}
	return Hit{}, false
}

// Detector returns the frozen C19 TaintDetector seam bound to one scope.
// Without a tool name it applies the conservative generic scan (see Inspect).
func (p *Provenance) Detector(scopeID string) TaintDetector {
	return boundDetector{p, scopeID}
}

type boundDetector struct {
	p     *Provenance
	scope string
}

func (d boundDetector) TaintHit(params map[string]any) (string, bool) {
	h, ok := d.p.Inspect(d.scope, "", params)
	if !ok {
		return "", false
	}
	return h.Source(), true
}

// --- the write / sync gate (SPEC-06 §5 channel ⑤) ----------------------------

// writeGate decides whether payload parameters must be scanned at all, and
// under which channel label the scan result is reported.
//
// SPEC-06 §5 makes exactly ONE of the six exfil channels conditional on a
// landing site: ⑤ fs.write INTO A SYNC DIR. Channel ⑥ (HTTP POST body/URL)
// has no landing condition, so an exemption must never be able to travel with
// a parameter name: the previous rule ("skip the keys called content/data
// whenever the call has a path outside a sync dir") let
// `{url, path:<non-sync>, data:<tainted>}` through while the same call with
// `body` was caught (adversarial re-verification M-7 — a key-name-selected
// fail-open).
//
// The gate therefore closes (payload exempt) ONLY for a call that is a plain
// local write and nothing else:
//
//   - the tool is not a named channel with its own sink (web.search/notify/
//     clipboard.write: a `path` key there is just another parameter);
//   - a write target is present and CONFIRMED to sit outside every sync root —
//     a missing, unresolvable or reparse-traversing target keeps the gate open
//     (missing information produces a DENY, not a pass);
//   - nothing in the call addresses a remote destination, neither by key name
//     nor by a URL-shaped value, and the nesting budget did not cut that check
//     short (unproven = open).
//
// Which key carries the payload is not an input to that decision, so the
// exemption cannot be dodged — nor claimed — by naming.
func (p *Provenance) writeGate(tool string, params map[string]any) (open bool, ch Channel) {
	def, named := channelTable[tool]
	if named && def.ch != ChSyncWrite {
		return true, ChUnknown // named non-write channel: no local-write exemption
	}
	path, hasPath := firstStringParam(params, pathKeys)
	if !hasPath {
		if named {
			// fs.write whose target cannot be verified is the shape the gate
			// exists to protect: fail closed toward scanning.
			return true, ChSyncWrite
		}
		return true, ChUnknown // nothing write-shaped: the generic scan
	}
	remote, cut := hasRemoteSink(params, p.maxDepth)
	st := p.IsSyncPath(path)
	switch {
	case st.Sync:
		return true, ChSyncWrite
	case remote || cut:
		return true, ChHTTP // a sink in the same call: the payload can leave
	}
	return false, ChSyncWrite // confirmed plain local write: not an exfil channel
}

// isPathKey reports whether a parameter names the write target.
func isPathKey(k string) bool {
	lk := strings.ToLower(strings.TrimSpace(k))
	for _, pk := range pathKeys {
		if lk == pk {
			return true
		}
	}
	return false
}

// hasRemoteSink reports whether one of the call's parameters addresses a remote
// destination, by key name or by an URL-shaped value. The second result is the
// nesting budget cutting the value check short, which the caller must treat as
// "a sink may be hidden below the budget".
func hasRemoteSink(params map[string]any, maxDepth int) (present, truncated bool) {
	for _, k := range sortedKeys(params) {
		if !isRemoteSinkKey(k) {
			continue
		}
		ss, cut := collectStrings(params[k], maxDepth)
		truncated = truncated || cut
		for _, s := range ss {
			if strings.TrimSpace(s) != "" {
				present = true
			}
		}
	}
	if present {
		return true, truncated
	}
	ss, cut := collectStrings(params, maxDepth)
	truncated = truncated || cut
	for _, s := range ss {
		if looksRemoteSink(s) {
			return true, truncated
		}
	}
	return false, truncated
}

// isRemoteSinkKey matches the remote-sink key names (case/edge insensitive,
// plus the *_url / *_endpoint spellings).
func isRemoteSinkKey(k string) bool {
	lk := strings.ToLower(strings.TrimSpace(k))
	if remoteSinkKeys[lk] {
		return true
	}
	for _, suf := range remoteSinkSuffixes {
		if len(lk) > len(suf) && strings.HasSuffix(lk, suf) {
			return true
		}
	}
	return false
}

// looksRemoteSink reports whether a VALUE is itself the address of a sink, so
// renaming the key cannot close the gate either. Only an address-shaped value
// counts: a document body that merely quotes a link is not a sink (otherwise
// every plain local write of a markdown file with URLs in it would be scanned,
// which is the false-positive side this gate must not create).
func looksRemoteSink(s string) bool {
	t := strings.TrimSpace(s)
	if t == "" || strings.ContainsAny(t, " \t\r\n") {
		return false // not a single token -> a body, not an address
	}
	lower := strings.ToLower(t)
	if len(lower) > 128 {
		lower = lower[:128] // a scheme always sits at the head
	}
	for _, pre := range remoteSinkSchemes {
		if strings.HasPrefix(lower, pre) {
			return true
		}
	}
	return false
}

// --- param walking -----------------------------------------------------------

func sortedKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func firstStringParam(m map[string]any, keys []string) (string, bool) {
	for _, k := range keys {
		if s, ok := m[k].(string); ok && strings.TrimSpace(s) != "" {
			return s, true
		}
	}
	return "", false
}

// collectStrings flattens nested map/slice string values down to maxDepth.
// The second result reports that the walk was CUT OFF by the nesting budget:
// callers must fail closed on it (adversarial report M-4) — a silent nil here
// is how a six-level wrapper escaped the generic scan.
func collectStrings(v any, maxDepth int) ([]string, bool) {
	return walkStrings(v, 0, maxDepth)
}

func walkStrings(v any, depth, maxDepth int) ([]string, bool) {
	switch t := v.(type) {
	case string:
		if t == "" {
			return nil, false
		}
		return []string{t}, false
	case []byte: // byte payloads are text carriers too (B-2 type escape)
		if len(t) == 0 {
			return nil, false
		}
		return []string{string(t)}, false
	case []string:
		return t, false
	case map[string]any:
		if depth >= maxDepth {
			return nil, len(t) > 0
		}
		var out []string
		cut := false
		for _, k := range sortedKeys(t) {
			ss, c := walkStrings(t[k], depth+1, maxDepth)
			out = append(out, ss...)
			cut = cut || c
		}
		return out, cut
	case []any:
		if depth >= maxDepth {
			return nil, len(t) > 0
		}
		var out []string
		cut := false
		for _, e := range t {
			ss, c := walkStrings(e, depth+1, maxDepth)
			out = append(out, ss...)
			cut = cut || c
		}
		return out, cut
	}
	return nil, false
}

// String for diagnostics without leaking content.
func (h Hit) String() string {
	return fmt.Sprintf("taint-hit channel=%s source=%s fragment_len=%d", h.Channel, h.Source(), len([]rune(h.Fragment)))
}
