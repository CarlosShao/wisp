package memory

// Ticket 29 AC#2's replay half: the pinned golden-SSE recordings under
// internal/memory/testdata/golden_sse, served to the REAL C5 extraction chain.
//
// WHY BYTES ON DISK INSTEAD OF AN IN-TEST STRING. The AC#2 test must drive the
// production chain, not a rule-based double. The chain under test is
//
//	pinned .sse bytes -> httptest HTTP -> openaichat adapter (SSE parse)
//	    -> llm.TurnCollector -> ParseExtractionReply -> normalizeOutcome
//	    -> CompareExtraction against testdata/golden_extract.json
//
// so the fixture and the recordings are two independent pinned artefacts: the
// recording is what the assistant answered for that scripted conversation, the
// fixture is what the profile table must end up holding after the batch lands.
// A hand-written fake provider would collapse the two into one and skip the
// adapter's HTTP/SSE path entirely (the seam SPEC-05's golden-SSE harness
// exists for, cf. internal/llm/golden and the *.sse precedent under
// internal/agent/testdata/golden).
//
// BATCH-COMMIT IS ASSERTED AT THIS SEAM: each recording holds exactly ONE
// response and its replayer is consumed by the first request, so one scripted
// conversation costs exactly ONE C5 call - one prompt-cache invalidation per
// task. A second call for the same conversation reaches an exhausted replayer
// (599) and fails the test loudly instead of quietly reusing bytes. The router
// also refuses a request whose system prompt drifted from the constant prefix:
// if two tasks sent different prefixes, "one invalidation per task" would be
// false even with one call each.
//
// The recordings are embedded for the same reason the fixture is: -overlay
// mutations of testdata must reach the test binary, which is how the reverse
// criterion corrupts an expectation without editing a reviewed file.

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/llm/openaichat"
)

const (
	// goldenSSEDir holds one recording per golden case, named after its id.
	goldenSSEDir = "testdata/golden_sse"
	// goldenSSEModel is the model id the replayed chain is assembled with.
	goldenSSEModel = "golden-extract-replay"
	// goldenSSEProvider is the provider name carried in llm.ProviderInfo.
	goldenSSEProvider = "golden-sse"
)

//go:embed testdata/golden_sse/*.sse
var goldenSSEFS embed.FS

// newGoldenSSEChain assembles the real extraction chain over the pinned
// recordings: a real openai-chat adapter talking real HTTP to an in-process
// server that replays the case's recorded SSE.
func newGoldenSSEChain(t *testing.T, set GoldenExtractSet) (Extractor, *goldenSSERouter) {
	t.Helper()
	router := newGoldenSSERouter(t, set)
	provider := openaichat.New(llm.EndpointOptions{
		Provider: goldenSSEProvider,
		Model:    goldenSSEModel,
		Protocol: "openai-chat",
		BaseURL:  goldenSSEServerURL(t, router),
	})
	ex, err := NewL1Extractor(provider, goldenSSEModel)
	if err != nil {
		t.Fatalf("NewL1Extractor over the golden replay provider: %v", err)
	}
	return ex, router
}

// goldenSSEServerURL starts the replay handler on a real loopback HTTP server,
// which is what keeps the adapter's HTTP + SSE path genuine.
func goldenSSEServerURL(t *testing.T, h http.Handler) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}

// goldenSSERouter is the deterministic replay provider's wire half: an HTTP
// handler that answers the adapter's request with the ONE recorded response
// pinned for that conversation script.
type goldenSSERouter struct {
	mu        sync.Mutex
	replayers map[string]*golden.Replayer // case id -> single-response replayer
	byScript  map[string]string           // pinned transcript -> case id
	served    map[string]int              // case id -> responses served
	systems   map[string]int              // distinct system prompts seen
	requests  int                         // total extraction requests
}

// newGoldenSSERouter loads every recording and validates the file set: one
// response per case, one case per file, no orphan and no duplicate script.
// These are the AC#2 bookkeeping claims (10 conversations, one task each), so
// they fail HERE rather than as a confusing 599 mid-test.
func newGoldenSSERouter(t *testing.T, set GoldenExtractSet) *goldenSSERouter {
	t.Helper()
	router := &goldenSSERouter{
		replayers: make(map[string]*golden.Replayer, len(set.Cases)),
		byScript:  make(map[string]string, len(set.Cases)),
		served:    make(map[string]int, len(set.Cases)),
		systems:   make(map[string]int),
	}
	known := make(map[string]bool, len(set.Cases))
	for i := range set.Cases {
		c := set.Cases[i]
		known[c.ID] = true
		script, err := goldenScript(c)
		if err != nil {
			t.Fatalf("golden replay key: %v", err)
		}
		if other, dup := router.byScript[script]; dup {
			t.Fatalf("cases %s and %s carry the same transcript; the replay key would be ambiguous", other, c.ID)
		}
		router.byScript[script] = c.ID
	}

	entries, err := goldenSSEFS.ReadDir(goldenSSEDir)
	if err != nil {
		t.Fatalf("read embedded %s: %v", goldenSSEDir, err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".sse") {
			continue
		}
		id := strings.TrimSuffix(name, ".sse")
		if !known[id] {
			t.Fatalf("%s is an orphan recording: no golden case has id %q "+
				"(the fixture and the recording set must move together)", name, id)
		}
		src, err := goldenSSEFS.ReadFile(path.Join(goldenSSEDir, name))
		if err != nil {
			t.Fatalf("read embedded %s: %v", name, err)
		}
		rs, err := golden.Parse(src)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		if len(rs) != 1 {
			t.Fatalf("%s holds %d recorded responses, want exactly 1 "+
				"(one C5 call per task: a second response would mean a second prompt-cache invalidation)", name, len(rs))
		}
		if rs[0].Status != http.StatusOK {
			t.Fatalf("%s records status %d, want 200", name, rs[0].Status)
		}
		router.replayers[id] = golden.NewReplayer(rs)
	}
	for i := range set.Cases {
		if router.replayers[set.Cases[i].ID] == nil {
			t.Fatalf("case %s has no pinned recording under %s/%s.sse "+
				"(every scripted conversation needs one assistant-text SSE)", set.Cases[i].ID, goldenSSEDir, set.Cases[i].ID)
		}
	}
	return router
}

// ServeHTTP answers one extraction request with its pinned recording. The
// request is decoded only to find which script it carries; the response bytes
// are the adapter's real input.
func (r *goldenSSERouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		http.Error(w, "golden replay: reading the adapter request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	call, err := decodeGoldenReplayCall(body)
	if err != nil {
		http.Error(w, "golden replay: "+err.Error(), http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests++
	r.systems[call.System]++
	if len(r.systems) > 1 {
		http.Error(w, "golden replay: the extraction system prompt is not byte-identical across tasks; "+
			"the cache-stable prefix drifted, so \"one invalidation per task\" is no longer true", 599)
		return
	}
	id, ok := r.byScript[call.Transcript]
	if !ok {
		http.Error(w, "golden replay: the request transcript matches no pinned case "+
			"(the extractor changed what it sends, or a case's turns changed)", 599)
		return
	}
	r.served[id]++
	// One response only: a second call for this case is an exhausted replayer.
	r.replayers[id].ServeHTTP(w, req)
}

// goldenReplayCall is what the router needs from one adapter request: the
// constant system prompt (cache-prefix check) and the rendered transcript
// (which pinned case this is).
type goldenReplayCall struct {
	System     string
	Transcript string
}

// decodeGoldenReplayCall decodes the adapter's wire request. L1Extractor sends
// one system message plus one user message carrying the transcript; anything
// else is a shape the recordings cannot answer honestly, so it is refused.
func decodeGoldenReplayCall(body []byte) (goldenReplayCall, error) {
	var wire struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		return goldenReplayCall{}, fmt.Errorf("decoding the adapter request: %w", err)
	}
	var call goldenReplayCall
	users := 0
	for _, m := range wire.Messages {
		// The extraction request is text-only, so the adapter encodes content
		// as a plain JSON string (the parts array is the image path).
		var text string
		if err := json.Unmarshal(m.Content, &text); err != nil {
			return goldenReplayCall{}, fmt.Errorf("%s message content is not plain text: %w", m.Role, err)
		}
		switch m.Role {
		case "system":
			if call.System != "" {
				return goldenReplayCall{}, fmt.Errorf("request carried more than one system message")
			}
			call.System = text
		case "user":
			users++
			call.Transcript = text
		default:
			return goldenReplayCall{}, fmt.Errorf("request carried a %q message; extraction is text-only", m.Role)
		}
	}
	if call.System == "" {
		return goldenReplayCall{}, fmt.Errorf("request carried no system prompt")
	}
	if users != 1 {
		return goldenReplayCall{}, fmt.Errorf("request carried %d user messages, want 1 (one task = one prompt)", users)
	}
	return call, nil
}

// goldenScript renders the request body the extractor must send for c. The
// replay key is the credential the recording answers to, so the fixture's own
// turns are checked against the extractor's window rules first: a turn the
// window would drop or truncate would make the key disagree with the request.
func goldenScript(c GoldenCase) (string, error) {
	turns := make([]GoldenTurn, 0, len(c.Turns))
	for _, tr := range c.Turns {
		role := strings.ToLower(strings.TrimSpace(tr.Role))
		if role != "user" && role != "assistant" {
			return "", fmt.Errorf("case %s turn role %q is outside user|assistant", c.ID, tr.Role)
		}
		if strings.TrimSpace(tr.Text) == "" {
			return "", fmt.Errorf("case %s has an empty turn; the extractor drops it and the replay key would drift", c.ID)
		}
		if n := utf8.RuneCountInString(tr.Text); n > l1DefaultMaxTurnRune {
			return "", fmt.Errorf("case %s turn is %d runes, past the extractor's %d-rune window", c.ID, n, l1DefaultMaxTurnRune)
		}
		turns = append(turns, GoldenTurn{Role: role, Text: tr.Text})
	}
	if len(turns) > l1DefaultMaxTurns {
		return "", fmt.Errorf("case %s has %d turns, past the extractor's %d-turn window", c.ID, len(turns), l1DefaultMaxTurns)
	}
	return renderTranscript(turns), nil
}

// callsFor reports how many recorded responses this case has consumed.
func (r *goldenSSERouter) callsFor(id string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.served[id]
}

// totalCalls reports how many extraction requests reached the seam.
func (r *goldenSSERouter) totalCalls() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.requests
}

// prefixVariants reports how many distinct extraction system prompts were sent.
// One variant across all tasks is what makes the prefix cache-stable.
func (r *goldenSSERouter) prefixVariants() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.systems)
}

// exhausted reports whether every recording was consumed.
func (r *goldenSSERouter) exhausted() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rep := range r.replayers {
		if rep.Remaining() != 0 {
			return false
		}
	}
	return true
}
