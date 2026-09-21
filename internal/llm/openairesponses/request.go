package openairesponses

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// C7 -> wire encoding (request side) for /v1/responses.

type wireInputItem struct {
	Type      string          `json:"type"`
	Role      string          `json:"role,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	CallID    string          `json:"call_id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Arguments string          `json:"arguments,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
	ID        string          `json:"id,omitempty"`
	Summary   json.RawMessage `json:"summary,omitempty"`
}

type wireContentPart struct {
	Type       string `json:"type"`
	Text       string `json:"text,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	Detail     string `json:"detail,omitempty"`
	ToolCallID string `json:"tool_use_id,omitempty"`
}

type wireToolDef struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	Strict      *bool           `json:"strict,omitempty"`
}

type wireReasoning struct {
	Effort  string `json:"effort,omitempty"`
	Summary string `json:"summary,omitempty"`
}

type wireRequest struct {
	Model           string          `json:"model"`
	Instructions    string          `json:"instructions,omitempty"`
	Input           []wireInputItem `json:"input"`
	Tools           []wireToolDef   `json:"tools,omitempty"`
	ToolChoice      json.RawMessage `json:"tool_choice,omitempty"`
	Temperature     *float64        `json:"temperature,omitempty"`
	MaxOutputTokens int             `json:"max_output_tokens,omitempty"`
	Reasoning       *wireReasoning  `json:"reasoning,omitempty"`
	Store           *bool           `json:"store,omitempty"`
	Stream          bool            `json:"stream"`
	Include         []string        `json:"include,omitempty"`
}

// effortLevels maps the config thinking enum onto reasoning.effort.
var effortLevels = map[string]string{"low": "low", "medium": "medium", "high": "high"}

// buildWire encodes the normalized request.
func (a *Adapter) buildWire(req *llm.Request) ([]byte, error) {
	if len(req.StopSequences) > 0 {
		if !a.loose() {
			return nil, observe.New(observe.ClassInternal,
				"openai-responses: the protocol has no stop-sequence field; drop StopSequences or set compat.loose")
		}
		// compat.loose: tolerate a provider that cannot honor them.
	}

	w := wireRequest{
		Model:           req.Model,
		Temperature:     req.Temperature,
		MaxOutputTokens: req.MaxOutputTokens,
		Stream:          true,
	}

	// System prompt: the protocol's own instructions field (prefix-stable, so
	// implicit provider caching still applies).
	for _, p := range req.System {
		t, ok := p.(llm.TextPart)
		if !ok {
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-responses: system prompt supports TextPart only, got %T", p))
		}
		if w.Instructions != "" {
			w.Instructions += "\n"
		}
		w.Instructions += t.Text
	}

	for _, m := range req.Messages {
		items, err := a.encodeMessage(req, m)
		if err != nil {
			return nil, err
		}
		w.Input = append(w.Input, items...)
	}

	for _, t := range req.Tools {
		params := t.Parameters
		if len(params) == 0 {
			params = json.RawMessage(`{"type":"object"}`)
		}
		if !json.Valid(params) {
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-responses: tool %q has invalid JSON schema", t.Name))
		}
		w.Tools = append(w.Tools, wireToolDef{
			Type: "function", Name: t.Name,
			Description: t.Description, Parameters: params,
		})
	}

	if req.ToolChoice != nil {
		tc, err := encodeToolChoice(req.ToolChoice)
		if err != nil {
			return nil, err
		}
		w.ToolChoice = tc
	}

	if intensity := strings.ToLower(strings.TrimSpace(req.ThinkingIntensity)); intensity != "" && intensity != "off" {
		effort, ok := effortLevels[intensity]
		if !ok {
			return nil, observe.New(observe.ClassConfig,
				fmt.Sprintf("openai-responses: unknown thinking_intensity %q (want off|low|medium|high)", req.ThinkingIntensity))
		}
		// summary="auto" is what makes reasoning visible as a stream of
		// ReasoningDelta events instead of a silent item.
		w.Reasoning = &wireReasoning{Effort: effort, Summary: "auto"}
	}

	// Reasoning items must come back in the response for the reasoning
	// deltas to stream; this is the documented include path.
	if w.Reasoning != nil {
		w.Include = []string{"reasoning.encrypted_content"}
	}

	b, err := json.Marshal(w)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "openai-responses: encode request")
	}
	return b, nil
}

// encodeMessage maps one normalized turn onto input items (a tool result is
// its own item, so a message can produce more than one).
func (a *Adapter) encodeMessage(req *llm.Request, m llm.Message) ([]wireInputItem, error) {
	switch m.Role {
	case llm.RoleUser:
		parts := make([]wireContentPart, 0, len(m.Content))
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.TextPart:
				parts = append(parts, wireContentPart{Type: "input_text", Text: c.Text})
			case llm.ImagePart:
				url, err := a.imageDataURL(req, c)
				if err != nil {
					return nil, err
				}
				parts = append(parts, wireContentPart{Type: "input_image", ImageURL: url})
			default:
				return nil, observe.New(observe.ClassInternal,
					fmt.Sprintf("openai-responses: user message cannot carry %T", p))
			}
		}
		return []wireInputItem{{
			Type: "message", Role: "user",
			Content: mustJSON(parts),
		}}, nil

	case llm.RoleAssistant:
		var out []wireInputItem
		var parts []wireContentPart
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.TextPart:
				if c.Text != "" {
					parts = append(parts, wireContentPart{Type: "output_text", Text: c.Text})
				}
			case llm.ToolUsePart:
				args := string(c.Input)
				if strings.TrimSpace(args) == "" {
					args = "{}"
				}
				if !json.Valid([]byte(args)) {
					return nil, observe.New(observe.ClassInternal,
						fmt.Sprintf("openai-responses: tool_use %q arguments are not valid JSON", c.ID))
				}
				out = append(out, wireInputItem{
					Type:   "function_call",
					CallID: c.ID, Name: c.Name, Arguments: args,
				})
			default:
				return nil, observe.New(observe.ClassInternal,
					fmt.Sprintf("openai-responses: assistant message cannot carry %T", p))
			}
		}
		if len(parts) > 0 {
			out = append([]wireInputItem{{
				Type: "message", Role: "assistant",
				Content: mustJSON(parts),
			}}, out...)
		}
		if len(out) == 0 {
			out = append(out, wireInputItem{
				Type: "message", Role: "assistant",
				Content: mustJSON([]wireContentPart{{Type: "output_text", Text: ""}}),
			})
		}
		return out, nil

	case llm.RoleTool:
		if len(m.Content) != 1 {
			return nil, observe.New(observe.ClassInternal,
				"openai-responses: a tool message must carry exactly one ToolResultPart")
		}
		res, ok := m.Content[0].(llm.ToolResultPart)
		if !ok {
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-responses: tool message carries %T, want ToolResultPart", m.Content[0]))
		}
		var sb strings.Builder
		for i, p := range res.Content {
			t, ok := p.(llm.TextPart)
			if !ok {
				return nil, observe.New(observe.ClassInternal,
					fmt.Sprintf("openai-responses: tool result part %d is %T (only TextPart)", i, p))
			}
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(t.Text)
		}
		out := sb.String()
		if res.IsError {
			out = "tool error: " + out
		}
		return []wireInputItem{{
			Type: "function_call_output", CallID: res.ID,
			Output: mustJSON([]wireContentPart{{Type: "input_text", Text: out}}),
		}}, nil

	case llm.RoleSystem:
		return nil, observe.New(observe.ClassInternal,
			"openai-responses: a mid-conversation system message has no input-item form (put it in Request.System)")
	default:
		return nil, observe.New(observe.ClassInternal,
			fmt.Sprintf("openai-responses: unsupported message role %q", m.Role))
	}
}

// imageDataURL resolves the image ref into a data: URL (bytes never logged).
func (a *Adapter) imageDataURL(req *llm.Request, img llm.ImagePart) (string, error) {
	if req.ResolveBytes == nil {
		return "", observe.New(observe.ClassInternal,
			"openai-responses: request carries an image but ResolveBytes is not wired")
	}
	mime, data, err := req.ResolveBytes(img.BytesRef)
	if err != nil {
		return "", observe.Wrap(observe.ClassInternal, err, "openai-responses: resolve image bytes ref")
	}
	if mime == "" {
		mime = img.MimeType
	}
	if mime == "" {
		mime = "application/octet-stream"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

// encodeToolChoice maps the normalized choice onto the Responses vocabulary.
func encodeToolChoice(tc *llm.ToolChoice) (json.RawMessage, error) {
	switch tc.Mode {
	case llm.ToolChoiceAuto:
		return mustJSON("auto"), nil
	case llm.ToolChoiceNone:
		return mustJSON("none"), nil
	case llm.ToolChoiceRequired:
		return mustJSON("required"), nil
	case llm.ToolChoiceTool:
		if tc.ToolName == "" {
			return nil, observe.New(observe.ClassInternal,
				"openai-responses: tool_choice=tool requires ToolName")
		}
		return mustJSON(map[string]string{"type": "function", "name": tc.ToolName}), nil
	default:
		return nil, observe.New(observe.ClassInternal,
			fmt.Sprintf("openai-responses: unknown tool_choice mode %q", tc.Mode))
	}
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
