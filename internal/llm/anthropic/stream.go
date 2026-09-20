package anthropic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Decode side: the Anthropic Messages SSE dialect (event:/data: pairs whose
// payload carries a `type` discriminator) mapped onto C6.

// wireEvent is the envelope every Messages SSE data payload shares.
type wireEvent struct {
	Type string `json:"type"`

	// message_start
	Message *wireMessage `json:"message"`

	// content_block_start / content_block_delta / content_block_stop
	Index        int        `json:"index"`
	ContentBlock *wireBlock `json:"content_block"`
	Delta        *wireDelta `json:"delta"`
	Usage        *wireUsage `json:"usage"`

	// message_delta
	StopReason string `json:"stop_reason"`

	// error
	Error *wireError `json:"error"`
}

type wireMessage struct {
	ID      string      `json:"id"`
	Role    string      `json:"role"`
	Model   string      `json:"model"`
	Usage   *wireUsage  `json:"usage"`
	Content []wireBlock `json:"content"`
}

// wireBlock is a content block at start time (text / tool_use / thinking).
type wireBlock struct {
	Type  string          `json:"type"`
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
	Text  string          `json:"text"`
}

// wireDelta is a content_block_delta / message_delta payload.
type wireDelta struct {
	Type        string `json:"type"` // text_delta|input_json_delta|thinking_delta|signature_delta|stop_reason
	Text        string `json:"text"`
	Thinking    string `json:"thinking"`
	Signature   string `json:"signature"`
	PartialJSON string `json:"partial_json"`
	StopReason  string `json:"stop_reason"`
}

// wireUsage is the Anthropic usage object (both message_start and
// message_delta carry a partial version of it).
type wireUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheCreation            *struct {
		Ephemeral5mInputTokens int `json:"ephemeral_5m_input_tokens"`
		Ephemeral1hInputTokens int `json:"ephemeral_1h_input_tokens"`
	} `json:"cache_creation"`
}

func (u *wireUsage) isZero() bool {
	return u == nil || (u.InputTokens == 0 && u.OutputTokens == 0 &&
		u.CacheReadInputTokens == 0 && u.CacheCreationInputTokens == 0)
}

// normalized folds the wire usage into the C6 payload. Cache READS are the
// CachedTokens; cache CREATION is deliberately not folded in (see the package
// usage-mapping note).
func (u *wireUsage) normalized() llm.Usage {
	if u == nil {
		return llm.Usage{}
	}
	return llm.Usage{
		InputTokens:  u.InputTokens,
		OutputTokens: u.OutputTokens,
		CachedTokens: u.CacheReadInputTokens,
	}
}

// wireError is the {"type":"error"} payload and the non-2xx body shape.
type wireError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    any    `json:"code"`
}

// parseWireError extracts the error object from a JSON body ({type:error} or
// {error:{...}}), returning nil when the payload is not an error.
func parseWireError(body []byte) *wireError {
	var probe struct {
		Type  string     `json:"type"`
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return nil
	}
	if probe.Error != nil && (probe.Error.Message != "" || probe.Error.Type != "") {
		if probe.Error.Type == "" {
			probe.Error.Type = "api_error"
		}
		return probe.Error
	}
	if probe.Type == "error" {
		// Some gateways flatten the payload: {"type":"error","message":...}.
		var flat wireError
		if err := json.Unmarshal(body, &flat); err == nil && (flat.Message != "" || flat.Type != "") {
			if flat.Type == "" {
				flat.Type = "api_error"
			}
			return &flat
		}
	}
	// OpenAI-shaped error body on an Anthropic endpoint (proxies do this):
	// {"error":{"message":...}} is handled above, so anything left is not an
	// error object.
	return nil
}

// codeString renders the optional code field (string or number).
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

// classify maps the wire error type to D37 through the HTTP-status table: the
// SSE error event carries no status, so the error type stands in for it
// (overloaded_error -> 529 provider, authentication_error -> 401 auth, ...).
func (e *wireError) classify() *observe.Error {
	status := 500
	switch e.Type {
	case "invalid_request_error":
		status = 400
	case "authentication_error", "unauthorized":
		status = 401
	case "permission_error", "forbidden":
		status = 403
	case "not_found_error":
		status = 404
	case "request_too_large":
		status = 413
	case "rate_limit_error":
		status = 429
	case "billing_error", "payment_required":
		status = 402
	case "overloaded_error":
		status = 529
	case "api_error":
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

// ---------------------------------------------------------------------------
// Stream state machine
// ---------------------------------------------------------------------------

// blockState tracks one open content block (Anthropic interleaves blocks by
// index; tool argument fragments belong to the block that opened them).
type blockState struct {
	kind string // text | tool_use | thinking | redacted_thinking
	id   string
	name string
}

// streamState accumulates one Messages stream.
type streamState struct {
	blocks     map[int]*blockState
	usage      llm.Usage
	sawUsage   bool
	stopReason string
	// sawMessageStop is the terminal sentinel (the [DONE] analogue).
	sawMessageStop bool
}

func newStreamState() *streamState {
	return &streamState{blocks: map[int]*blockState{}}
}

// openBlocks reports a content block that never received content_block_stop.
func (s *streamState) openBlocks() bool { return len(s.blocks) > 0 }

// handle applies one decoded SSE event, emitting C6 events. The bool reports
// whether the terminal sentinel was seen; err is a classified failure that
// caller must surface (the caller owns the terminal triple).
func (s *streamState) handle(_ context.Context, eventName string, data []byte,
	emit func(llm.StreamEvent) error) (bool, error) {

	var ev wireEvent
	if err := json.Unmarshal(data, &ev); err != nil {
		return false, observe.Wrap(observe.ClassProvider, err,
			"anthropic: stream payload is not valid JSON")
	}
	if ev.Type == "" {
		ev.Type = eventName // tolerate gateways that only set `event:`
	}

	switch ev.Type {
	case "ping", "":
		return false, nil

	case "error":
		if ev.Error == nil {
			var flat wireError
			if err := json.Unmarshal(data, &flat); err != nil || (flat.Message == "" && flat.Type == "") {
				return false, observe.New(observe.ClassProvider,
					"anthropic: error event without an error payload")
			}
			ev.Error = &flat
		}
		return false, ev.Error.classify()

	case "message_start":
		if ev.Message != nil && !ev.Message.Usage.isZero() {
			s.sawUsage = true
			s.usage = s.usage.Add(ev.Message.Usage.normalized())
			if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: s.usage}); err != nil {
				return false, err
			}
		}
		return false, nil

	case "content_block_start":
		blk := ev.ContentBlock
		if blk == nil {
			return false, nil
		}
		st := &blockState{kind: blk.Type, id: blk.ID, name: blk.Name}
		s.blocks[ev.Index] = st
		if st.kind == "tool_use" {
			if err := emit(llm.StreamEvent{Type: llm.EvToolCallStart,
				ToolCallID: st.id, ToolName: st.name}); err != nil {
				return false, err
			}
			// A non-empty input at start time (rare, but real for
			// pre-filled turns) is the args payload: replay it as one delta
			// so the assembler never loses it.
			if args := strings.TrimSpace(string(blk.Input)); args != "" && args != "{}" && args != "null" {
				if err := emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta,
					ToolCallID: st.id, ArgsDelta: args}); err != nil {
					return false, err
				}
			}
		}
		return false, nil

	case "content_block_delta":
		d := ev.Delta
		if d == nil {
			return false, nil
		}
		switch d.Type {
		case "text_delta":
			if d.Text != "" {
				if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: d.Text}); err != nil {
					return false, err
				}
			}
		case "thinking_delta":
			// Extended thinking -> ReasoningDelta (AC#3).
			if d.Thinking != "" {
				if err := emit(llm.StreamEvent{Type: llm.EvReasoningDelta, Text: d.Thinking}); err != nil {
					return false, err
				}
			}
		case "signature_delta":
			// Thinking signature verification material: no C6 event exists
			// for it and the event set is contract-exact, so it is dropped at
			// the seam (the agent core never needs it).
		case "input_json_delta":
			st := s.blocks[ev.Index]
			if st == nil || st.kind != "tool_use" {
				if st != nil {
					// A json delta on a text/thinking block is provider
					// noise: drop it rather than corrupt a tool call.
					return false, nil
				}
				return false, observe.New(observe.ClassProvider,
					"anthropic: input_json_delta for unknown content block index")
			}
			if d.PartialJSON == "" {
				return false, nil
			}
			if err := emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta,
				ToolCallID: st.id, ArgsDelta: d.PartialJSON}); err != nil {
				return false, err
			}
		default:
			// Unknown delta kind: strict-but-quiet - the event set is exact
			// and nothing here invalidates the turn.
		}
		return false, nil

	case "content_block_stop":
		st := s.blocks[ev.Index]
		if st == nil {
			return false, nil
		}
		if st.kind == "tool_use" {
			if err := emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: st.id}); err != nil {
				return false, err
			}
		}
		delete(s.blocks, ev.Index)
		return false, nil

	case "message_delta":
		if r := deltaStopReason(ev); r != "" {
			s.stopReason = r
		}
		if !ev.Usage.isZero() {
			s.sawUsage = true
			s.usage = s.usage.Add(ev.Usage.normalized())
			if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: s.usage}); err != nil {
				return false, err
			}
		}
		return false, nil

	case "message_stop":
		s.sawMessageStop = true
		return true, nil

	default:
		// session/agent/stream-style events from gateways: ignored.
		return false, nil
	}
}

// deltaStopReason reads the stop_reason from either place the dialect puts it
// (delta.stop_reason per spec; some gateways hoist it to the top level).
func deltaStopReason(ev wireEvent) string {
	if ev.Delta != nil && ev.Delta.StopReason != "" {
		return ev.Delta.StopReason
	}
	return ev.StopReason
}
