package agent

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// D15(3) long-output spill (SPEC-05 §4.2): a single tool result over the
// scaled spill threshold is written in full to
// <artifacts>\tool-output-<id>.txt and the context keeps only the head, the
// tail, the total length and the path; the model can re-read it with fs.read
// on demand. Raw output over the scaled hard cap is truncated first and
// marked truncated=true (the same 1MB ceiling D15 gives shell.exec stdout).
//
// This is a HOST-INTERNAL artifact write (D34 note(2)): it is deliberately not
// a gated tool, and it lands in the memory store's artifacts directory so the
// 500MB LRU quota job of ticket 04 owns its lifetime.

// Spiller caps and spills tool output against a scaled budget set.
type Spiller struct {
	dir string
	b   Budgets

	mu       sync.Mutex
	sequence int
}

// NewSpiller builds a spiller writing into dir (the store's artifacts dir).
func NewSpiller(dir string, b Budgets) *Spiller { return &Spiller{dir: dir, b: b} }

// Spill is the outcome of preparing one tool result for the context.
type Spill struct {
	// Text is what goes into the context: either the original text or the
	// head/tail stub.
	Text string
	// Spilled is true when the full text went to a file and Text is a stub.
	Spilled bool
	// Path is the artifact path ("" when not spilled).
	Path string
	// Name is the artifact file name ("" when not spilled).
	Name string
	// TotalBytes / TotalTokens describe the ORIGINAL output.
	TotalBytes  int
	TotalTokens int
	// TruncatedRaw is true when the raw hard cap cut the bytes before spill.
	TruncatedRaw bool
	// KeptHead / KeptTail report the stub sizes in tokens.
	KeptHead int
	KeptTail int
}

// CapRaw applies the D15(3) raw hard ceiling (scaled RawOutputCapBytes) and
// reports the cut text plus whether truncation happened.
func (s *Spiller) CapRaw(text string) (string, bool) {
	ceiling := s.b.RawOutputCapBytes
	if ceiling <= 0 || len(text) <= ceiling {
		return text, false
	}
	cut := ceiling
	for cut > 0 && !utf8Start(text[cut]) {
		cut-- // never split a rune
	}
	return text[:cut], true
}

// Prepare caps the raw output, then spills it when it exceeds the scaled spill
// threshold. callID identifies the tool call; it is model-supplied, so it is
// sanitized before it can reach a file name.
func (s *Spiller) Prepare(callID, text string) (Spill, error) {
	capped, rawTrunc := s.CapRaw(text)
	out := Spill{
		Text:         capped,
		TotalBytes:   len(capped),
		TotalTokens:  ApproxTokens(capped),
		TruncatedRaw: rawTrunc,
	}
	if rawTrunc {
		out.Text = capped + fmt.Sprintf(
			"\n[truncated=true: 原始输出超过 %d 字节硬上限，已截断]", s.b.RawOutputCapBytes)
	}
	if ApproxTokens(out.Text) <= s.b.SpillTokens {
		return out, nil
	}

	name := artifactName(callID, s.nextSequence())
	path := filepath.Join(s.dir, name)
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Spill{}, observe.Wrap(observe.ClassResource, err, "agent: create artifacts dir")
	}
	if err := writeFileExclusive(path, []byte(text)); err != nil {
		return Spill{}, observe.Wrap(observe.ClassResource, err, "agent: write spill artifact")
	}

	head := takeTokens(out.Text, s.b.SpillHeadTokens)
	tail := takeTokensLast(out.Text, s.b.SpillTailTokens)
	out.KeptHead = ApproxTokens(head)
	out.KeptTail = ApproxTokens(tail)
	out.Text = fmt.Sprintf(
		"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s",
		head, len(out.Text)-len(head)-len(tail), out.TotalBytes, out.TotalTokens, path, tail)
	out.Spilled = true
	out.Path = path
	out.Name = name
	return out, nil
}

func (s *Spiller) nextSequence() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sequence++
	return s.sequence
}

// artifactName renders the D15(3) file name tool-output-<id>.txt. Model-supplied
// ids are restricted to [A-Za-z0-9_-]; anything else falls back to a
// sequence-based name so no path shape can be injected through a tool-call id.
func artifactName(callID string, seq int) string {
	var b strings.Builder
	for i := 0; i < len(callID); i++ {
		c := callID[i]
		if c == '_' || c == '-' || (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			b.WriteByte(c)
		}
	}
	id := b.String()
	if id == "" {
		id = fmt.Sprintf("seq%d", seq)
	}
	return "tool-output-" + id + ".txt"
}

// writeFileExclusive writes the artifact, overwriting an existing file of the
// same call id (a retried call with the same id re-lands its own output).
func writeFileExclusive(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o600); err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return err
		}
		return err
	}
	return nil
}

// takeTokens returns the leading budget tokens of s (4 bytes per token under
// the shared heuristic, cut on a rune boundary).
func takeTokens(s string, budget int) string {
	if budget <= 0 {
		return ""
	}
	n := budget * 4
	if n >= len(s) {
		return s
	}
	for n > 0 && !utf8Start(s[n]) {
		n--
	}
	return s[:n]
}

// takeTokensLast returns the trailing budget tokens of s.
func takeTokensLast(s string, budget int) string {
	if budget <= 0 {
		return ""
	}
	n := budget * 4
	if n >= len(s) {
		return s
	}
	start := len(s) - n
	for start < len(s) && !utf8Start(s[start]) {
		start++
	}
	return s[start:]
}
