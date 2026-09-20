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
// are matched key-by-key by Inspect. TTS announce and the net body do not have
// a frozen tool name yet (tickets 22/26 name them), so they call
// CheckText(scope, ChTTS|ChHTTP, text) directly; and any call whose tool is
// NOT in the channel table fail-closes to a conservative scan of every string
// parameter (write payloads excepted: content/data params are only scanned
// when the call's path lands in a sync dir, which is the spec's own rule for
// fs.write, not an invention).
//
// Wiring to the assessor: the frozen C19 seam is TaintDetector.TaintHit(params)
// — it carries no tool name — so callers bind it per scope:
// assessor.WithTaintDetector(prov.Detector(taskScopeID)) yields the
// conservative shape-based verdict; loop code that knows the tool name should
// prefer Inspect()/CheckText() directly (wiring lands with tickets 20/21/22/26).
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

// writeChannelKeys are payload parameters scanned ONLY when the call's path
// resolves into a sync dir (fs.write to a normal dir stays a plain L1 write;
// SPEC-06 §5 makes only the sync-dir landing an exfil channel).
var writeChannelKeys = map[string]bool{"content": true, "data": true}

// pathKeys are the parameters naming the write target for the sync-dir test.
var pathKeys = []string{"path", "file", "filepath", "dest", "destination"}

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
	// AllowedDirs / ReparseExceptions mirror [fs] config (the suspect
	// fallback and C26 Resolve need them; membership is decided by Resolve).
	AllowedDirs       []string
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
		maxSrc = 262144
	}
	maxScan := o.MaxScanRunes
	if maxScan <= 0 {
		maxScan = 2097152
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
		s.complete = len(s.roots) > 0
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
			RuneLen: len(m.idx.norm), MarkedAt: m.at})
	}
	return out
}

// --- detection ---------------------------------------------------------------

// Inspect checks one outgoing call against the scope's taints. tool is the
// D34 name ("" = unknown, forces the conservative generic scan). Returns the
// first hit in deterministic order (channel keys in table order, remaining
// params sorted, marks oldest-first).
func (p *Provenance) Inspect(scopeID, tool string, params map[string]any) (Hit, bool) {
	p.mu.RLock()
	marks, known := p.scopes[scopeID]
	if len(marks) == 0 {
		p.mu.RUnlock()
		if !known {
			logf("risk/C25: Inspect on unknown scope %q: treated as untainted (empty store); ensure OpenScope at task start", scopeID)
		}
		return Hit{}, false
	}
	marks = marks[:len(marks):len(marks)] // snapshot (append never mutates the shared prefix)
	p.mu.RUnlock()

	// Build the candidate list: (channel, text) pairs derived from params.
	type cand struct {
		ch   Channel
		text string
	}
	var cands []cand
	if def, ok := channelTable[tool]; ok {
		switch def.ch {
		case ChSyncWrite:
			path, hasPath := firstStringParam(params, pathKeys)
			// A malformed fs.write without a usable path is unverifiable and
			// fail-closes to the sync side (strict).
			if !hasPath || p.IsSyncPath(path).Sync {
				for _, k := range sortedKeys(params) {
					if writeChannelKeys[k] {
						if v, ok := params[k].(string); ok {
							cands = append(cands, cand{ChSyncWrite, v})
						}
					}
				}
			}
			// Non-sync fs.write: content is NOT an exfil channel (SPEC-06 §5);
			// fall through with no candidates.
		default:
			for _, k := range p.channel[tool] {
				if v, ok := params[k].(string); ok {
					cands = append(cands, cand{def.ch, v})
				}
			}
		}
	} else {
		// Unknown tool (or the TaintDetector adapter): fail-closed generic
		// scan of every string parameter, write payloads gated on sync path.
		path, hasPath := firstStringParam(params, pathKeys)
		syncGateOpen := !hasPath || p.IsSyncPath(path).Sync // no path = unverifiable = fail-closed open
		for _, k := range sortedKeys(params) {
			v := params[k]
			if writeChannelKeys[k] && !syncGateOpen {
				continue
			}
			for _, s := range collectStrings(v, 0) {
				cands = append(cands, cand{ChUnknown, s})
			}
		}
	}
	for _, c := range cands {
		if h, ok := p.matchText(scopeID, marks, c.ch, c.text); ok {
			return h, true
		}
	}
	return Hit{}, false
}

// CheckText matches one outgoing string against a scope's taints under an
// explicit channel label — the entry point for channels without a frozen D34
// tool name yet: TTS announce text (ticket 26) and HTTP POST body/url
// (ticket 22 / net layer).
func (p *Provenance) CheckText(scopeID string, ch Channel, text string) (Hit, bool) {
	p.mu.RLock()
	marks := p.scopes[scopeID]
	marks = marks[:len(marks):len(marks)]
	p.mu.RUnlock()
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

// collectStrings flattens nested map/slice string values (bounded depth).
func collectStrings(v any, depth int) []string {
	if depth > 4 {
		return nil
	}
	switch t := v.(type) {
	case string:
		if t != "" {
			return []string{t}
		}
	case map[string]any:
		var out []string
		for _, k := range sortedKeys(t) {
			out = append(out, collectStrings(t[k], depth+1)...)
		}
		return out
	case []any:
		var out []string
		for _, e := range t {
			out = append(out, collectStrings(e, depth+1)...)
		}
		return out
	case []string:
		return t
	}
	return nil
}

// String for diagnostics without leaking content.
func (h Hit) String() string {
	return fmt.Sprintf("taint-hit channel=%s source=%s fragment_len=%d", h.Channel, h.Source(), len([]rune(h.Fragment)))
}
