package openaichat

import (
	"encoding/json"
	"strconv"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Wire types of the OpenAI Chat Completions stream (decode-side). Field
// names follow the protocol; provider quirks (reasoning_content vs
// reasoning, flat vs nested cached tokens, missing index) are tolerated here
// so the event emitter stays clean.

type wireChunk struct {
	Choices []wireChoice `json:"choices"`
	Usage   *wireUsage   `json:"usage"`
	Error   *wireError   `json:"error"`
}

type wireChoice struct {
	Index        int       `json:"index"`
	Delta        wireDelta `json:"delta"`
	FinishReason string    `json:"finish_reason"`
}

type wireDelta struct {
	Role             string         `json:"role"`
	Content          string         `json:"content"`
	Reasoning        string         `json:"reasoning"`
	ReasoningContent string         `json:"reasoning_content"`
	Refusal          string         `json:"refusal"`
	ToolCalls        []wireToolFrag `json:"tool_calls"`
}

func (d *wireDelta) text() string { return d.Content }

// reasoning joins the two reasoning spellings (deepseek/qwen style
// reasoning_content and the newer reasoning).
func (d *wireDelta) reasoning() string {
	if d.ReasoningContent != "" {
		return d.ReasoningContent
	}
	return d.Reasoning
}

// wireToolFrag is one streamed tool_calls fragment.
type wireToolFrag struct {
	// Index is a pointer: absence is a provider quirk the assembler treats
	// per compat mode (null index = malformed in strict mode).
	Index    *int   `json:"index"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type wireUsage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	CachedTokens        int `json:"cached_tokens"` // flat spelling (qwen/deepseek style)
	PromptTokensDetails *struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
	CompletionTokensDetails *struct {
		ReasoningTokens int `json:"reasoning_tokens"`
	} `json:"completion_tokens_details"`
}

func (u *wireUsage) isZero() bool {
	return u.PromptTokens == 0 && u.CompletionTokens == 0 &&
		u.CachedTokens == 0 && u.PromptTokensDetails == nil
}

// normalized converts the wire usage to the C6 Usage event payload.
func (u *wireUsage) normalized() llm.Usage {
	out := llm.Usage{InputTokens: u.PromptTokens, OutputTokens: u.CompletionTokens}
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > u.CachedTokens {
		out.CachedTokens = u.PromptTokensDetails.CachedTokens
	} else {
		out.CachedTokens = u.CachedTokens
	}
	return out
}

// wireError is the OpenAI error payload (HTTP error bodies and some
// mid-stream error chunks).
type wireError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

func parseWireError(body []byte) *wireError {
	var probe struct {
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(body, &probe); err != nil || probe.Error == nil {
		return nil
	}
	return probe.Error
}

// codeString renders the code field, which some providers send as a string
// and others as a number.
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

// classify maps the wire error to D37. Status is unknown for mid-stream
// errors: 500 keeps provider/retryable semantics while 4xx-style "type"
// strings (invalid_request_error / authentication_error / ...) refine it.
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
	case "overloaded_error":
		status = 529
	}
	return llm.NewHTTPError(llm.HTTPErrorDetail{
		Status: status, Code: e.codeString(), Type: e.Type, Message: e.Message,
	}, 0)
}
