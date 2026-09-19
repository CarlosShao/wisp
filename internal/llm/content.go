package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// C7 internal message representation (SPEC-05 sec 3.3). The agent core only
// ever sees these types - protocol shapes (OpenAI content parts, Anthropic
// blocks) never cross the adapter boundary.
//
// Privacy rule (SPEC-05 sec 3.3): ImagePart.BytesRef points at memory or a
// temp file. It is NEVER persisted to the DB, NEVER written to a log and
// NEVER included in a diagnostics bundle. RedactContent enforces the log half
// of that rule; adapters must resolve the ref into bytes only at wire-encode
// time (Request.ResolveBytes) and must not log the encoded payload.

// Role of a normalized Message.
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
	RoleTool      MessageRole = "tool" // carries ToolResultPart content
)

// Content is the sealed C7 content-part interface. The only implementors are
// the four structs below (pinned by TestContentSealed).
type Content interface {
	contentPart()
}

// TextPart is a plain text fragment.
type TextPart struct {
	Text string
}

// ImagePart is an image reference. BytesRef points at in-memory or temporary
// bytes (never a DB row, never a log line); MimeType is the RFC 2046 type;
// Alt is the user-visible description.
type ImagePart struct {
	MimeType string
	BytesRef string
	Alt      string
}

// ToolUsePart is a model-issued tool call. Input is the raw JSON arguments.
type ToolUsePart struct {
	ID    string
	Name  string
	Input json.RawMessage
}

// ToolResultPart is the host's answer to a ToolUsePart. Content carries the
// result parts (text today; richer media when the protocol allows).
type ToolResultPart struct {
	ID      string
	Content []Content
	IsError bool
}

func (TextPart) contentPart()       {}
func (ImagePart) contentPart()      {}
func (ToolUsePart) contentPart()    {}
func (ToolResultPart) contentPart() {}

// contentPartNames is used by the sealing test and by RedactContent.
var contentPartNames = []string{"TextPart", "ImagePart", "ToolUsePart", "ToolResultPart"}

// Message is one normalized conversation turn.
type Message struct {
	Role    MessageRole
	Content []Content
}

// ToolDef describes one tool the model may call. Parameters is the JSON
// Schema object (passed through verbatim to the provider).
type ToolDef struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

// ToolChoiceMode selects how the model may use the declared tools.
type ToolChoiceMode string

const (
	ToolChoiceAuto     ToolChoiceMode = "auto"
	ToolChoiceNone     ToolChoiceMode = "none"
	ToolChoiceRequired ToolChoiceMode = "required"
	ToolChoiceTool     ToolChoiceMode = "tool" // force ToolName
)

// ToolChoice is the normalized tool-selection directive.
type ToolChoice struct {
	Mode     ToolChoiceMode
	ToolName string // used when Mode == ToolChoiceTool
}

// RedactContent renders content parts for LOGGING ONLY: image bytes refs are
// replaced by a stable placeholder so no path or ref can leak into a log,
// diagnostics bundle or error message (SPEC-05 sec 3.3 privacy rule). Tool
// arguments are truncated to keep log lines bounded.
func RedactContent(parts []Content, maxArgRunes int) string {
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			b.WriteString(" | ")
		}
		switch c := p.(type) {
		case TextPart:
			b.WriteString(truncateRunes(c.Text, maxArgRunes))
		case ImagePart:
			fmt.Fprintf(&b, "<image mime=%s bytes-redacted alt=%s>", c.MimeType, truncateRunes(c.Alt, 60))
		case ToolUsePart:
			fmt.Fprintf(&b, "<tool_use id=%s name=%s args=%s>", c.ID, c.Name, truncateRunes(string(c.Input), maxArgRunes))
		case ToolResultPart:
			fmt.Fprintf(&b, "<tool_result id=%s error=%v parts=%d>", c.ID, c.IsError, len(c.Content))
		default:
			fmt.Fprintf(&b, "<unknown part %T>", p)
		}
	}
	return b.String()
}

func truncateRunes(s string, n int) string {
	if n <= 0 || len(s) == 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
