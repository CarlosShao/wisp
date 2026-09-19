package openaichat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// C7 -> wire encoding (request side).

type wireMessage struct {
	Role       string          `json:"role"`
	Content    json.RawMessage `json:"content"`
	ToolCalls  []wireToolUse   `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
	Name       string          `json:"name,omitempty"`
}

type wireToolUse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type wireToolDef struct {
	Type     string           `json:"type"`
	Function wireToolDefInner `json:"function"`
}

type wireToolDefInner struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type wireToolChoice struct {
	// Mode variants: "auto" | "none" | "required" as a bare string, or the
	// forced-function form below.
	Type     string            `json:"type,omitempty"`
	Function *wireToolChoiceFn `json:"function,omitempty"`
}

type wireToolChoiceFn struct {
	Name string `json:"name"`
}

type wireRequest struct {
	Model         string          `json:"model"`
	Messages      []wireMessage   `json:"messages"`
	Tools         []wireToolDef   `json:"tools,omitempty"`
	ToolChoice    json.RawMessage `json:"tool_choice,omitempty"`
	Temperature   *float64        `json:"temperature,omitempty"`
	MaxTokens     int             `json:"max_tokens,omitempty"`
	Stop          []string        `json:"stop,omitempty"`
	Stream        bool            `json:"stream"`
	StreamOptions *struct {
		IncludeUsage bool `json:"include_usage"`
	} `json:"stream_options,omitempty"`
}

// buildWire encodes the normalized request. Text-only messages use the plain
// string content form (broadest compat); part-bearing messages use the
// content-parts array.
func (a *Adapter) buildWire(req *llm.Request) ([]byte, error) {
	w := wireRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxOutputTokens,
		Stop:        req.StopSequences,
		Stream:      true,
	}
	if !a.loose() {
		w.StreamOptions = &struct {
			IncludeUsage bool `json:"include_usage"`
		}{IncludeUsage: true}
	}

	// System prompt first (cache-prefix order lives with the caller, sec 4.1).
	if len(req.System) > 0 {
		var sb strings.Builder
		for _, p := range req.System {
			t, ok := p.(llm.TextPart)
			if !ok {
				return nil, observe.New(observe.ClassInternal,
					"openai-chat: system prompt supports TextPart only")
			}
			sb.WriteString(t.Text)
		}
		w.Messages = append(w.Messages, wireMessage{
			Role: "system", Content: jsonString(sb.String()),
		})
	}

	for _, m := range req.Messages {
		wm, err := a.encodeMessage(req, m)
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
		w.Tools = append(w.Tools, wireToolDef{
			Type: "function",
			Function: wireToolDefInner{
				Name: t.Name, Description: t.Description, Parameters: params,
			},
		})
	}

	if req.ToolChoice != nil {
		switch req.ToolChoice.Mode {
		case llm.ToolChoiceAuto:
			w.ToolChoice = json.RawMessage(`"auto"`)
		case llm.ToolChoiceNone:
			w.ToolChoice = json.RawMessage(`"none"`)
		case llm.ToolChoiceRequired:
			w.ToolChoice = json.RawMessage(`"required"`)
		case llm.ToolChoiceTool:
			b, _ := json.Marshal(wireToolChoice{
				Type:     "function",
				Function: &wireToolChoiceFn{Name: req.ToolChoice.ToolName},
			})
			w.ToolChoice = b
		default:
			return nil, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-chat: unknown tool_choice mode %q", req.ToolChoice.Mode))
		}
	}

	b, err := json.Marshal(w)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "openai-chat: encode request")
	}
	return b, nil
}

// encodeMessage maps one normalized message.
func (a *Adapter) encodeMessage(req *llm.Request, m llm.Message) (wireMessage, error) {
	switch m.Role {
	case llm.RoleUser:
		return a.encodeUser(req, m)
	case llm.RoleAssistant:
		return a.encodeAssistant(m)
	case llm.RoleTool:
		return a.encodeToolResult(m)
	default:
		return wireMessage{}, observe.New(observe.ClassInternal,
			fmt.Sprintf("openai-chat: unsupported message role %q", m.Role))
	}
}

// encodeUser emits plain-string content for pure text turns and the parts
// array when images are present.
func (a *Adapter) encodeUser(req *llm.Request, m llm.Message) (wireMessage, error) {
	texts := 0
	images := 0
	for _, p := range m.Content {
		switch p.(type) {
		case llm.TextPart:
			texts++
		case llm.ImagePart:
			images++
		default:
			return wireMessage{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-chat: user message cannot carry %T", p))
		}
	}
	if images == 0 {
		var sb strings.Builder
		for i, p := range m.Content {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(p.(llm.TextPart).Text)
		}
		return wireMessage{Role: "user", Content: jsonString(sb.String())}, nil
	}

	parts := make([]json.RawMessage, 0, len(m.Content))
	for _, p := range m.Content {
		switch c := p.(type) {
		case llm.TextPart:
			b, _ := json.Marshal(map[string]string{"type": "text", "text": c.Text})
			parts = append(parts, b)
		case llm.ImagePart:
			dataURL, err := a.imageDataURL(req, c)
			if err != nil {
				return wireMessage{}, err
			}
			b, _ := json.Marshal(map[string]any{
				"type":      "image_url",
				"image_url": map[string]string{"url": dataURL},
			})
			parts = append(parts, b)
		}
	}
	arr, err := json.Marshal(parts)
	if err != nil {
		return wireMessage{}, observe.Wrap(observe.ClassInternal, err, "openai-chat: encode user parts")
	}
	return wireMessage{Role: "user", Content: arr}, nil
}

// imageDataURL resolves the ImagePart.BytesRef into a data: URL. Resolution
// errors are ClassInternal (the caller owns the bytes store); the encoded
// bytes NEVER reach a log (wire payload only).
func (a *Adapter) imageDataURL(req *llm.Request, img llm.ImagePart) (string, error) {
	if req.ResolveBytes == nil {
		return "", observe.New(observe.ClassInternal,
			"openai-chat: request carries an image but ResolveBytes is not wired")
	}
	mime, data, err := req.ResolveBytes(img.BytesRef)
	if err != nil {
		return "", observe.Wrap(observe.ClassInternal, err,
			"openai-chat: resolve image bytes ref")
	}
	if mime == "" {
		mime = img.MimeType
	}
	if mime == "" {
		mime = "application/octet-stream"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

// encodeAssistant emits text and/or tool_use parts (OpenAI merges them into
// one assistant message: tool_calls beside content).
func (a *Adapter) encodeAssistant(m llm.Message) (wireMessage, error) {
	wm := wireMessage{Role: "assistant"}
	var texts []string
	for _, p := range m.Content {
		switch c := p.(type) {
		case llm.TextPart:
			texts = append(texts, c.Text)
		case llm.ToolUsePart:
			args := string(c.Input)
			if args == "" {
				args = "{}"
			}
			var wt wireToolUse
			wt.ID = c.ID
			wt.Type = "function"
			wt.Function.Name = c.Name
			wt.Function.Arguments = args
			wm.ToolCalls = append(wm.ToolCalls, wt)
		default:
			return wireMessage{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-chat: assistant message cannot carry %T", p))
		}
	}
	if len(texts) > 0 {
		wm.Content = jsonString(strings.Join(texts, "\n"))
	} else if wm.ToolCalls != nil {
		wm.Content = json.RawMessage(`null`) // tool-only turn: content null
	} else {
		wm.Content = json.RawMessage(`""`)
	}
	return wm, nil
}

// encodeToolResult flattens ToolResultPart content into the role:"tool"
// message OpenAI expects (one message per tool call id).
func (a *Adapter) encodeToolResult(m llm.Message) (wireMessage, error) {
	if len(m.Content) != 1 {
		return wireMessage{}, observe.New(observe.ClassInternal,
			"openai-chat: a tool message must carry exactly one ToolResultPart")
	}
	res, ok := m.Content[0].(llm.ToolResultPart)
	if !ok {
		return wireMessage{}, observe.New(observe.ClassInternal,
			fmt.Sprintf("openai-chat: tool message carries %T, want ToolResultPart", m.Content[0]))
	}
	var sb strings.Builder
	for i, p := range res.Content {
		t, ok := p.(llm.TextPart)
		if !ok {
			return wireMessage{}, observe.New(observe.ClassInternal,
				fmt.Sprintf("openai-chat: tool result part %d is %T (only TextPart)", i, p))
		}
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(t.Text)
	}
	if res.IsError {
		// Keep the error visible to the model without a second vocabulary.
		out := "tool error: " + sb.String()
		return wireMessage{Role: "tool", ToolCallID: res.ID, Content: jsonString(out)}, nil
	}
	return wireMessage{Role: "tool", ToolCallID: res.ID, Content: jsonString(sb.String())}, nil
}

// jsonString marshals a Go string as a JSON string value.
func jsonString(s string) json.RawMessage {
	b, _ := json.Marshal(s)
	return b
}
