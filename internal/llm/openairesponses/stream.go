package openairesponses

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Decode side: the OpenAI Responses SSE dialect mapped onto C6.

// wireFrame is the union every Responses stream frame decodes into.
type wireFrame struct {
	Type string `json:"type"`

	// Frames carry their payload either at the top level or under `response`.
	Response *wireResponse `json:"response"`
	Item     *wireItem     `json:"item"`
	Delta    string        `json:"delta"`
	Text     string        `json:"text"`
	Reason   string        `json:"reason"`
	Part     *wirePart     `json:"part"`

	ItemID    string `json:"item_id"`
	OutputIdx int    `json:"output_index"`

	Error *wireError `json:"error"`
}

type wireResponse struct {
	ID                string     `json:"id"`
	Status            string     `json:"status"`
	Model             string     `json:"model"`
	Usage             *wireUsage `json:"usage"`
	IncompleteDetails *struct {
		Reason string `json:"reason"`
	} `json:"incomplete_details"`
	Error *wireError `json:"error"`
}

// wireItem is an output item (message / reasoning / function_call).
type wireItem struct {
	ID        string          `json:"id"`
	CallID    string          `json:"call_id"`
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Arguments string          `json:"arguments"`
	Content   json.RawMessage `json:"content"`
	Summary   json.RawMessage `json:"summary"`
}

// wirePart is an added content part (response.content_part.added).
type wirePart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type wireUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	// Both nested spellings appear in the wild for cache reads.
	InputTokensDetails *struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"input_tokens_details"`
	CacheReadInputTokens int `json:"cache_read_input_tokens"`
	CachedTokens         int `json:"cached_tokens"` // flat spelling (compat servers)
}

func (u *wireUsage) isZero() bool {
	return u == nil || (u.InputTokens == 0 && u.OutputTokens == 0 &&
		u.CachedTokens == 0 && u.CacheReadInputTokens == 0 &&
		u.InputTokensDetails == nil)
}

func (u *wireUsage) normalized() llm.Usage {
	if u == nil {
		return llm.Usage{}
	}
	out := llm.Usage{InputTokens: u.InputTokens, OutputTokens: u.OutputTokens}
	if u.InputTokensDetails != nil && u.InputTokensDetails.CachedTokens > out.CachedTokens {
		out.CachedTokens = u.InputTokensDetails.CachedTokens
	}
	if u.CacheReadInputTokens > out.CachedTokens {
		out.CachedTokens = u.CacheReadInputTokens
	}
	if u.CachedTokens > out.CachedTokens {
		out.CachedTokens = u.CachedTokens
	}
	return out
}

// wireError is the OpenAI error payload (HTTP bodies and terminal frames).
type wireError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

func (e *wireError) codeString() string {
	switch v := e.Code.(type) {
	case string:
		return v
	case float64:
		return strconv.Itoa(int(v))
	default:
		if v == nil {
			return ""
		}
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// classify maps an error payload to D37 through the shared HTTP table (a
// mid-stream error carries no status, so the type stands in, exactly like the
// openai-chat adapter does).
func (e *wireError) classify() *observe.Error {
	status := 500
	switch e.Type {
	case "invalid_request_error":
		status = 400
	case "authentication_error":
		status = 401
	case "permission_error":
		status = 403
	case "not_found_error":
		status = 404
	case "rate_limit_error":
		status = 429
	case "insufficient_quota", "billing_error":
		status = 402
	case "server_error", "overloaded_error":
		status = 500
	}
	code := e.codeString()
	if code == "" {
		code = e.Type
	}
	return llm.NewHTTPError(llm.HTTPErrorDetail{
		Status: status, Code: code, Type: e.Type, Message: e.Message,
	}, 0)
}

// parseWireError extracts an error object from a JSON body.
func parseWireError(body []byte) *wireError {
	var probe struct {
		Type  string     `json:"type"`
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return nil
	}
	if probe.Error != nil && (probe.Error.Message != "" || probe.Error.Type != "") {
		return probe.Error
	}
	if probe.Type == "error" {
		var flat wireError
		if err := json.Unmarshal(body, &flat); err == nil && (flat.Message != "" || flat.Type != "") {
			return &flat
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Stream state
// ---------------------------------------------------------------------------

// itemState tracks one streamed output item (deltas arrive keyed by item id).
type itemState struct {
	kind       string // message | reasoning | function_call | refusal
	callID     string
	name       string
	streamText bool // any text/reasoning delta was seen for it
	streamArgs bool // any argument delta was seen for it
	open       bool
}

// streamState accumulates one Responses stream.
type streamState struct {
	items map[string]*itemState // by item id
	// function_call items keyed by their event item id -> tool call id, so
	// argument deltas that only carry item_id resolve to the caller's id.
	calledFunction bool
	sawRefusal     bool
	sawUsage       bool
	sawDone        bool
	sawTerminal    bool
	status         string
	// incompleteReason is response.incomplete_details.reason.
	incompleteReason string
	usage            llm.Usage
}

func newStreamState() *streamState {
	return &streamState{items: map[string]*itemState{}}
}

// handle applies one frame; the bool reports an early terminal (a failed
// response ends the stream).
func (s *streamState) handle(data string, emit func(llm.StreamEvent) error) (bool, error) {
	var f wireFrame
	if err := json.Unmarshal([]byte(data), &f); err != nil {
		return false, observe.Wrap(observe.ClassProvider, err,
			"openai-responses: stream frame is not valid JSON")
	}
	switch f.Type {
	case "", "response.created", "response.in_progress",
		"response.output_text.done", "response.content_part.added",
		"response.content_part.done", "response.reasoning_summary_part.added",
		"response.reasoning_summary_text.done", "response.output_text.annotation.added":
		return false, nil

	case "response.output_item.added":
		return false, s.itemAdded(f, emit)

	case "response.output_item.done":
		return false, s.itemDone(f, emit)

	case "response.output_text.delta", "response.refusal.delta":
		if f.Type == "response.refusal.delta" {
			s.sawRefusal = true
		}
		if f.Delta != "" {
			if st := s.item(f.ItemID); st != nil {
				st.streamText = true
			}
			if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: f.Delta}); err != nil {
				return false, err
			}
		}
		return false, nil

	case "response.reasoning_summary_text.delta", "response.reasoning_text.delta":
		// Reasoning items -> ReasoningDelta (AC#3).
		if f.Delta != "" {
			if st := s.item(f.ItemID); st != nil {
				st.streamText = true
			}
			if err := emit(llm.StreamEvent{Type: llm.EvReasoningDelta, Text: f.Delta}); err != nil {
				return false, err
			}
		}
		return false, nil

	case "response.function_call_arguments.delta":
		st := s.item(f.ItemID)
		if st == nil || st.kind != "function_call" {
			return false, observe.New(observe.ClassProvider,
				"openai-responses: argument delta for unknown function_call item "+f.ItemID)
		}
		if f.Delta == "" {
			return false, nil
		}
		st.streamArgs = true
		if err := emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta,
			ToolCallID: st.callID, ArgsDelta: f.Delta}); err != nil {
			return false, err
		}
		return false, nil

	case "error":
		if f.Error == nil {
			var flat wireError
			if err := json.Unmarshal([]byte(data), &flat); err != nil ||
				(flat.Message == "" && flat.Type == "") {
				return false, observe.New(observe.ClassProvider,
					"openai-responses: error frame without an error payload")
			}
			f.Error = &flat
		}
		return true, f.Error.classify()

	case "response.failed", "response.incomplete", "response.completed",
		"response.usage":
		return s.terminal(f, emit)

	default:
		// Unknown frame kinds (custom tools, web search, ...) are ignored: the
		// C6 event set is contract-exact and no other event fits them.
		return false, nil
	}
}

// itemAdded registers a new output item and opens its tool call.
func (s *streamState) itemAdded(f wireFrame, emit func(llm.StreamEvent) error) error {
	it := s.itemOf(f)
	if it == nil {
		return nil
	}
	st := &itemState{kind: it.Type, callID: it.callID(), name: it.Name, open: true}
	s.items[it.id()] = st
	switch st.kind {
	case "function_call":
		s.calledFunction = true
		return emit(llm.StreamEvent{Type: llm.EvToolCallStart,
			ToolCallID: st.callID, ToolName: st.name})
	case "refusal":
		s.sawRefusal = true
	}
	return nil
}

// itemDone closes an item, back-filling payload a non-incremental server only
// delivered in the done frame.
func (s *streamState) itemDone(f wireFrame, emit func(llm.StreamEvent) error) error {
	it := s.itemOf(f)
	if it == nil {
		return nil
	}
	id := it.id()
	st := s.items[id]
	if st == nil {
		// A server that skipped output_item.added: adopt the item so its
		// payload is not lost (a compat-provider quirk, not a contract gap).
		st = &itemState{kind: it.Type, callID: it.callID(), name: it.Name, open: true}
		s.items[id] = st
		if st.kind == "function_call" {
			s.calledFunction = true
			if err := emit(llm.StreamEvent{Type: llm.EvToolCallStart,
				ToolCallID: st.callID, ToolName: st.name}); err != nil {
				return err
			}
		}
	}
	switch st.kind {
	case "function_call":
		if !st.streamArgs && strings.TrimSpace(it.Arguments) != "" {
			if err := emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta,
				ToolCallID: st.callID, ArgsDelta: it.Arguments}); err != nil {
				return err
			}
		}
		delete(s.items, id)
		return emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: st.callID})
	case "message":
		if !st.streamText {
			for _, t := range it.outputTexts() {
				if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: t}); err != nil {
					return err
				}
			}
		}
		delete(s.items, id)
		return nil
	default:
		// reasoning / unknown: nothing to close.
		delete(s.items, id)
		return nil
	}
}

// terminal absorbs a status/usage-bearing frame.
func (s *streamState) terminal(f wireFrame, emit func(llm.StreamEvent) error) (bool, error) {
	r := f.Response
	if f.Type == "response.usage" {
		// A gateway-only extra usage frame: merge the counters, keep the
		// terminal state untouched.
		if !r.usage().isZero() {
			s.sawUsage = true
			s.usage = s.usage.Add(r.usage().normalized())
			if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: s.usage}); err != nil {
				return false, err
			}
		}
		return false, nil
	}
	s.sawTerminal = true
	// The frame TYPE is the authoritative status (a terminal frame without a
	// status field is still that terminal); an explicit status wins.
	if statusFromType(f.Type) != "" {
		s.status = statusFromType(f.Type)
	}
	if r != nil {
		if r.Status != "" {
			s.status = r.Status
		}
		if r.IncompleteDetails != nil {
			s.incompleteReason = r.IncompleteDetails.Reason
		}
	}
	if !r.usage().isZero() {
		s.sawUsage = true
		s.usage = s.usage.Add(r.usage().normalized())
		if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: s.usage}); err != nil {
			return false, err
		}
	}
	if s.status == "failed" {
		e := r.failureError(f)
		return true, e
	}
	return false, nil
}

// statusFromType maps a terminal frame type to the response status.
func statusFromType(t string) string {
	switch t {
	case "response.completed":
		return "completed"
	case "response.incomplete":
		return "incomplete"
	case "response.failed":
		return "failed"
	}
	return ""
}

// item resolves an item state by either id spelling.
func (s *streamState) item(id string) *itemState {
	if st, ok := s.items[id]; ok {
		return st
	}
	return nil
}

// itemOf returns the frame's item, falling back to the response-carried item
// shape (some frames nest it).
func (s *streamState) itemOf(f wireFrame) *wireItem {
	if f.Item != nil && (f.Item.Type != "" || f.Item.ID != "" || f.Item.CallID != "") {
		return f.Item
	}
	return nil
}

// id prefers the item id (the key deltas use).
func (it *wireItem) id() string {
	if it.ID != "" {
		return it.ID
	}
	return it.CallID
}

// callID prefers call_id (the id the caller echoes back) and falls back to the
// item id for servers that reuse one field.
func (it *wireItem) callID() string {
	if it.CallID != "" {
		return it.CallID
	}
	return it.ID
}

// outputTexts pulls the text parts out of a message item's content.
func (it *wireItem) outputTexts() []string {
	if len(it.Content) == 0 {
		return nil
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(it.Content, &parts); err != nil {
		return nil
	}
	var out []string
	for _, p := range parts {
		if (p.Type == "output_text" || p.Type == "refusal") && p.Text != "" {
			out = append(out, p.Text)
		}
	}
	return out
}

// usage tolerates the nil response pointer on malformed frames.
func (r *wireResponse) usage() *wireUsage {
	if r == nil {
		return nil
	}
	return r.Usage
}

// failureError picks the error payload of a failed terminal frame.
func (r *wireResponse) failureError(f wireFrame) *observe.Error {
	if r != nil && r.Error != nil && (r.Error.Message != "" || r.Error.Type != "") {
		return r.Error.classify()
	}
	if f.Error != nil {
		return f.Error.classify()
	}
	return observe.New(observe.ClassProvider, "openai-responses: response failed without an error payload")
}

var _ = strconv.Itoa
