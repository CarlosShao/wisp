package anthropic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// C7 -> wire encoding (request side) plus the explicit prompt-caching
// breakpoint placement documented in the package comment.

// cacheControl is the Anthropic cache_control marker object.
var cacheControl = json.RawMessage(`{"type":"ephemeral"}`)

type wireTextBlock struct {
	Type         string          `json:"type"`
	Text         string          `json:"text"`
	CacheControl json.RawMessage `json:"cache_control,omitempty"`
}

type wireToolUseBlock struct {
	Type         string          `json:"type"`
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Input        json.RawMessage `json:"input"`
	CacheControl json.RawMessage `json:"cache_control,omitempty"`
}

type wireToolResultBlock struct {
	Type         string          `json:"type"`
	ToolUseID    string          `json:"tool_use_id"`
	Content      json.RawMessage `json:"content"`
	IsError      bool            `json:"is_error,omitempty"`
	CacheControl json.RawMessage `json:"cache_control,omitempty"`
}

type wireImageBlock struct {
	Type   string `json:"type"`
	Source struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source"`
}

type wireMessageIn struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type wireToolDef struct {
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	InputSchema  json.RawMessage `json:"input_schema"`
	CacheControl json.RawMessage `json:"cache_control,omitempty"`
}

type wireToolChoice struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type wireThinking struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens,omitempty"`
}

type wireRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`

	// System is the array form so a cache breakpoint can sit on it (the
	// string form cannot carry cache_control).
	System     json.RawMessage `json:"system,omitempty"`
	Messages   []wireMessageIn `json:"messages"`
	Tools      []wireToolDef   `json:"tools,omitempty"`
	ToolChoice json.RawMessage `json:"tool_choice,omitempty"`

	Temperature   *float64      `json:"temperature,omitempty"`
	TopP          *float64      `json:"top_p,omitempty"`
	StopSequences []string      `json:"stop_sequences,omitempty"`
	Thinking      *wireThinking `json:"thinking,omitempty"`
	Stream        bool          `json:"stream"`
}

// thinkingBudgets maps the config enum to budget_tokens (adapter-local
// numbers; see the package doc).
var thinkingBudgets = map[string]int{"low": 1024, "medium": 4096, "high": 8192}

// buildWire encodes the normalized request.
func (a *Adapter) buildWire(req *llm.Request) ([]byte, error) {
	breaks, err := a.planBreakpoints(req)
	if err != nil {
		return nil, err
	}

	w := wireRequest{
		Model:         req.Model,
		MaxTokens:     req.MaxOutputTokens,
		StopSequences: req.StopSequences,
		Stream:        true,
	}
	if w.MaxTokens <= 0 {
		w.MaxTokens = defaultMaxTokens // protocol requires the field
	}

	// Extended thinking: the intensity enum owns the budget parameter.
	if intensity := strings.ToLower(strings.TrimSpace(req.ThinkingIntensity)); intensity != "" && intensity != "off" {
		budget, ok := thinkingBudgets[intensity]
		if !ok {
			return nil, observe.New(observe.ClassConfig,
				fmt.Sprintf("anthropic: unknown thinking_intensity %q (want off|low|medium|high)", req.ThinkingIntensity))
		}
		if budget >= w.MaxTokens {
			budget = w.MaxTokens - 1 // protocol: max_tokens must exceed the budget
		}
		if budget < 0 {
			budget = 0
		}
		if budget > 0 {
			w.Thinking = &wireThinking{Type: "enabled", BudgetTokens: budget}
			// Anthropic rejects temperature != 1 while thinking is on.
			req = cloneWithoutTemperature(req)
		}
	} else if intensity == "off" {
		w.Thinking = &wireThinking{Type: "disabled"}
	}

	if len(req.System) > 0 {
		sys, err := a.encodeSystem(req, breaks.afterSystem)
		if err != nil {
			return nil, err
		}
		w.System = sys
	}

	for i, m := range req.Messages {
		wm, err := a.encodeMessage(req, m, breaks.forMessage[i])
		if err != nil {
			return nil, err
		}
		w.Messages = append(w.Messages, wm)
	}

	for _, t := range req.Tools {
		params := t.Parameters
		if len(params) == 0 {
			params = json.RawMessage(`{"type":"object"}`)
		}
		if !json.Valid(params) {
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: tool %q has invalid JSON schema", t.Name))
		}
		w.Tools = append(w.Tools, wireToolDef{Name: t.Name, Description: t.Description, InputSchema: params})
	}

	if req.ToolChoice != nil {
		tc, err := encodeToolChoice(req.ToolChoice)
		if err != nil {
			return nil, err
		}
		w.ToolChoice = tc
	}

	if req.Temperature != nil && w.Thinking == nil {
		v := *req.Temperature
		w.Temperature = &v
	}

	b, err := json.Marshal(w)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "anthropic: encode request")
	}
	return b, nil
}

// cloneWithoutTemperature returns a shallow copy with Temperature unset so the
// thinking parameter can be sent (see the package doc).
func cloneWithoutTemperature(req *llm.Request) *llm.Request {
	cp := *req
	cp.Temperature = nil
	return &cp
}

// breakpointPlan is the resolved cache_control placement: which system block
// and which messages carry a marker.
type breakpointPlan struct {
	afterSystem bool
	forMessage  []bool
}

func (p breakpointPlan) count() int {
	n := 0
	if p.afterSystem {
		n++
	}
	for _, b := range p.forMessage {
		if b {
			n++
		}
	}
	return n
}

// planBreakpoints resolves Request.CacheBreakpoints into wire placement.
// -1 marks the end of System; i >= 0 marks the end of Messages[i]. Anthropic
// allows 4 markers: more is our bug (strict) or a trim (loose).
func (a *Adapter) planBreakpoints(req *llm.Request) (breakpointPlan, error) {
	plan := breakpointPlan{forMessage: make([]bool, len(req.Messages))}
	for _, idx := range req.CacheBreakpoints {
		if idx == -1 {
			plan.afterSystem = true
			continue
		}
		if idx < 0 || idx >= len(req.Messages) {
			return plan, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: cache breakpoint index %d out of range for %d messages",
					idx, len(req.Messages)))
		}
		plan.forMessage[idx] = true
	}
	if plan.count() > maxCacheControl {
		if !a.loose() {
			return plan, observe.New(observe.ClassInternal, fmt.Sprintf(
				"anthropic: %d cache breakpoints exceed the protocol maximum of %d",
				plan.count(), maxCacheControl))
		}
		// loose: keep the declared order and drop the overflow.
		plan.trimTo(maxCacheControl)
	}
	return plan, nil
}

// trimTo drops trailing markers until n remain (system first, then message
// order - matching the declared priority).
func (p *breakpointPlan) trimTo(n int) {
	if p.count() <= n {
		return
	}
	for i := len(p.forMessage) - 1; i >= 0 && p.count() > n; i-- {
		p.forMessage[i] = false
	}
	if p.count() > n {
		p.afterSystem = false
	}
}

// encodeSystem renders the top-level system blocks, marking the last one when
// the -1 breakpoint is set.
func (a *Adapter) encodeSystem(req *llm.Request, markLast bool) (json.RawMessage, error) {
	blocks := make([]wireTextBlock, 0, len(req.System))
	for _, p := range req.System {
		t, ok := p.(llm.TextPart)
		if !ok {
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: system prompt supports TextPart only, got %T", p))
		}
		blocks = append(blocks, wireTextBlock{Type: "text", Text: t.Text})
	}
	if len(blocks) == 0 {
		return nil, nil
	}
	if markLast {
		blocks[len(blocks)-1].CacheControl = cacheControl
	}
	return json.Marshal(blocks)
}

// encodeMessage maps one normalized turn. Anthropic has no tool role: a tool
// result becomes a user message carrying a tool_result block.
func (a *Adapter) encodeMessage(req *llm.Request, m llm.Message, mark bool) (wireMessageIn, error) {
	switch m.Role {
	case llm.RoleUser:
		return a.encodeUser(req, m, mark)
	case llm.RoleAssistant:
		return a.encodeAssistant(m, mark)
	case llm.RoleTool:
		return a.encodeToolResult(m, mark)
	case llm.RoleSystem:
		return wireMessageIn{}, observe.New(observe.ClassInternal,
			"anthropic: a mid-conversation system message has no wire form (put it in Request.System)")
	default:
		return wireMessageIn{}, observe.New(observe.ClassInternal,
			fmt.Sprintf("anthropic: unsupported message role %q", m.Role))
	}
}

func (a *Adapter) encodeUser(req *llm.Request, m llm.Message, mark bool) (wireMessageIn, error) {
	blocks := make([]json.RawMessage, 0, len(m.Content))
	for _, p := range m.Content {
		switch c := p.(type) {
		case llm.TextPart:
			b, err := json.Marshal(wireTextBlock{Type: "text", Text: c.Text})
			if err != nil {
				return wireMessageIn{}, err
			}
			blocks = append(blocks, b)
		case llm.ImagePart:
			b, err := a.encodeImage(req, c)
			if err != nil {
				return wireMessageIn{}, err
			}
			blocks = append(blocks, b)
		default:
			return wireMessageIn{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: user message cannot carry %T", p))
		}
	}
	if len(blocks) == 0 {
		return wireMessageIn{}, observe.New(observe.ClassInternal,
			"anthropic: user message encodes to no content block")
	}
	if mark {
		markLastBlock(blocks)
	}
	return wireMessageIn{Role: "user", Content: mustArray(blocks)}, nil
}

func (a *Adapter) encodeAssistant(m llm.Message, mark bool) (wireMessageIn, error) {
	blocks := make([]json.RawMessage, 0, len(m.Content))
	for _, p := range m.Content {
		switch c := p.(type) {
		case llm.TextPart:
			if c.Text == "" {
				continue
			}
			b, err := json.Marshal(wireTextBlock{Type: "text", Text: c.Text})
			if err != nil {
				return wireMessageIn{}, err
			}
			blocks = append(blocks, b)
		case llm.ToolUsePart:
			input := c.Input
			if len(input) == 0 {
				input = json.RawMessage(`{}`)
			}
			if !json.Valid(input) {
				return wireMessageIn{}, observe.New(observe.ClassInternal,
					fmt.Sprintf("anthropic: tool_use %q input is not valid JSON", c.ID))
			}
			b, err := json.Marshal(wireToolUseBlock{Type: "tool_use", ID: c.ID,
				Name: c.Name, Input: input})
			if err != nil {
				return wireMessageIn{}, err
			}
			blocks = append(blocks, b)
		default:
			return wireMessageIn{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: assistant message cannot carry %T", p))
		}
	}
	if len(blocks) == 0 {
		// An empty assistant turn is not encodable; a zero-length text block
		// keeps the protocol happy (and history honest).
		b, _ := json.Marshal(wireTextBlock{Type: "text", Text: ""})
		blocks = append(blocks, b)
	}
	if mark {
		markLastBlock(blocks)
	}
	return wireMessageIn{Role: "assistant", Content: mustArray(blocks)}, nil
}

// encodeToolResult flattens a ToolResultPart into the user-message block
// Anthropic expects.
func (a *Adapter) encodeToolResult(m llm.Message, mark bool) (wireMessageIn, error) {
	if len(m.Content) != 1 {
		return wireMessageIn{}, observe.New(observe.ClassInternal,
			"anthropic: a tool message must carry exactly one ToolResultPart")
	}
	res, ok := m.Content[0].(llm.ToolResultPart)
	if !ok {
		return wireMessageIn{}, observe.New(observe.ClassInternal,
			fmt.Sprintf("anthropic: tool message carries %T, want ToolResultPart", m.Content[0]))
	}
	inner := make([]json.RawMessage, 0, len(res.Content))
	for i, p := range res.Content {
		t, ok := p.(llm.TextPart)
		if !ok {
			return wireMessageIn{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("anthropic: tool result part %d is %T (only TextPart)", i, p))
		}
		b, err := json.Marshal(wireTextBlock{Type: "text", Text: t.Text})
		if err != nil {
			return wireMessageIn{}, err
		}
		inner = append(inner, b)
	}
	if len(inner) == 0 {
		inner = append(inner, mustJSON(wireTextBlock{Type: "text", Text: " "}))
	}
	b, err := json.Marshal(wireToolResultBlock{Type: "tool_result", ToolUseID: res.ID,
		Content: mustArray(inner), IsError: res.IsError})
	if err != nil {
		return wireMessageIn{}, err
	}
	blocks := []json.RawMessage{b}
	if mark {
		markLastBlock(blocks)
	}
	return wireMessageIn{Role: "user", Content: mustArray(blocks)}, nil
}

// encodeImage resolves an ImagePart ref into a base64 source block. Resolved
// bytes reach the wire only, never a log.
func (a *Adapter) encodeImage(req *llm.Request, img llm.ImagePart) (json.RawMessage, error) {
	if req.ResolveBytes == nil {
		return nil, observe.New(observe.ClassInternal,
			"anthropic: request carries an image but ResolveBytes is not wired")
	}
	mime, data, err := req.ResolveBytes(img.BytesRef)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "anthropic: resolve image bytes ref")
	}
	if mime == "" {
		mime = img.MimeType
	}
	if mime == "" {
		mime = "image/png"
	}
	var blk wireImageBlock
	blk.Type = "image"
	blk.Source.Type = "base64"
	blk.Source.MediaType = mime
	blk.Source.Data = base64.StdEncoding.EncodeToString(data)
	return json.Marshal(blk)
}

// encodeToolChoice maps the normalized choice onto the Anthropic vocabulary
// (auto | any | tool).
func encodeToolChoice(tc *llm.ToolChoice) (json.RawMessage, error) {
	switch tc.Mode {
	case llm.ToolChoiceAuto:
		return mustJSON(wireToolChoice{Type: "auto"}), nil
	case llm.ToolChoiceNone:
		// Anthropic has no "none": omitting the tools is the equivalent, and
		// dropping the choice here keeps the declaration honest.
		return nil, nil
	case llm.ToolChoiceRequired:
		return mustJSON(wireToolChoice{Type: "any"}), nil
	case llm.ToolChoiceTool:
		if tc.ToolName == "" {
			return nil, observe.New(observe.ClassInternal,
				"anthropic: tool_choice=tool requires ToolName")
		}
		return mustJSON(wireToolChoice{Type: "tool", Name: tc.ToolName}), nil
	default:
		return nil, observe.New(observe.ClassInternal,
			fmt.Sprintf("anthropic: unknown tool_choice mode %q", tc.Mode))
	}
}

// markLastBlock attaches cache_control to the final block of an encoded
// message (the breakpoint position is AFTER the message).
func markLastBlock(blocks []json.RawMessage) {
	last := blocks[len(blocks)-1]
	var text wireTextBlock
	if err := json.Unmarshal(last, &text); err == nil && text.Type == "text" {
		blocks[len(blocks)-1] = mustJSON(wireTextBlock{Type: "text", Text: text.Text, CacheControl: cacheControl})
		return
	}
	var tu wireToolUseBlock
	if err := json.Unmarshal(last, &tu); err == nil && tu.Type == "tool_use" {
		blocks[len(blocks)-1] = mustJSON(wireToolUseBlock{Type: "tool_use", ID: tu.ID, Name: tu.Name,
			Input: tu.Input, CacheControl: cacheControl})
		return
	}
	var tr wireToolResultBlock
	if err := json.Unmarshal(last, &tr); err == nil && tr.Type == "tool_result" {
		blocks[len(blocks)-1] = mustJSON(wireToolResultBlock{Type: "tool_result", ToolUseID: tr.ToolUseID,
			Content: tr.Content, IsError: tr.IsError, CacheControl: cacheControl})
	}
	// Any other final block type cannot carry a marker in this protocol; the
	// breakpoint is then a no-op (prefix caching degrades, semantics do not).
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func mustArray(blocks []json.RawMessage) json.RawMessage {
	b, _ := json.Marshal(blocks)
	return b
}
