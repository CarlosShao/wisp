package agent

import (
	"crypto/sha256"
	"encoding/hex"
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
// <artifacts>\tool-output-<encoded id>.txt and the context keeps only the head,
// the tail, the total length and the path; the model can re-read it with fs.read
// on demand. Raw output over the scaled hard cap is truncated first and
// marked truncated=true (the same 1MB ceiling D15 gives shell.exec stdout).
//
// "<encoded id>", not "<id>": the name is an injective encoding of the
// model-supplied tool-call id, because a name two different ids share lets one
// tool call overwrite another's saved output (ticket 79, C25).
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
// percent-encoded (injectively, see artifactName) before it can reach a file
// name, and the artifact itself is created exclusively.
func (s *Spiller) Prepare(callID, text string) (Spill, error) {
	capped, rawTrunc := s.CapRaw(text)
	// The truncation notice is context text, never artifact bytes: the file has
	// to stay INSIDE the ceiling it announces (PLAN:432 truncates first and
	// lands the file afterwards), and a small window's tail budget would cut a
	// marker baked into the text right off the stub.
	notice := ""
	if rawTrunc {
		notice = fmt.Sprintf(
			"\n[truncated=true: 原始输出超过 %d 字节硬上限，已截断]", s.b.RawOutputCapBytes)
	}
	out := Spill{
		Text:         capped + notice,
		TotalBytes:   len(capped),
		TotalTokens:  ApproxTokens(capped),
		TruncatedRaw: rawTrunc,
	}
	if ApproxTokens(out.Text) <= s.b.SpillTokens {
		return out, nil
	}

	name := artifactName(callID, s.nextSequence())
	path := filepath.Join(s.dir, name)
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Spill{}, observe.Wrap(observe.ClassResource, err, "agent: create artifacts dir")
	}
	// The capped bytes, not the raw ones: TotalBytes below is what the stub
	// announces, so the artifact size and the announced length must agree.
	err := writeFileExclusive(path, []byte(capped))
	if errors.Is(err, fs.ErrExist) {
		// Because the name is injective in the logical id, an occupied path can
		// only be this same id's own earlier artifact: a retried tool call (or a
		// call whose id repeats across a restart). Documented behavior is
		// last-writer-wins, and D31 says the swap is a temp write plus rename, so
		// a model re-reading the path sees either the whole old output or the
		// whole new one, never a torn artifact.
		tmp := path + ".retry"
		if wErr := os.WriteFile(tmp, []byte(capped), 0o600); wErr != nil {
			return Spill{}, observe.Wrap(observe.ClassResource, wErr, "agent: write spill artifact retry")
		}
		if rErr := os.Rename(tmp, path); rErr != nil {
			_ = os.Remove(tmp)
			return Spill{}, observe.Wrap(observe.ClassResource, rErr, "agent: replace spill artifact")
		}
	} else if err != nil {
		return Spill{}, observe.Wrap(observe.ClassResource, err, "agent: write spill artifact")
	}

	head := takeTokens(capped, s.b.SpillHeadTokens)
	tail := takeTokensLast(capped, s.b.SpillTailTokens)
	out.KeptHead = ApproxTokens(head)
	out.KeptTail = ApproxTokens(tail)
	out.Text = fmt.Sprintf(
		"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s",
		head, len(capped)-len(head)-len(tail), out.TotalBytes, out.TotalTokens, path, tail) +
		notice
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

// artifactName renders the D15(3) file name tool-output-<id>.txt.
//
// The id is model-supplied, so it is ENCODED, never filtered. Ticket 79's
// defect was exactly the filtering: keeping only [A-Za-z0-9_-] folded the
// genuinely different ids `p/q`, `p\q` and `pq` onto the one name
// tool-output-pq.txt, and because the write below was not exclusive, a later
// call silently replaced an earlier call's bytes while the earlier Spill.Path
// (handed to the model, re-readable with fs.read) still named that file - a C25
// provenance break, not cosmetics.
//
// The encoding is percent-escaping (uppercase hex) of every byte outside
// [a-z0-9_-]. Escape, not hash: it is one-to-one without leaning on digest
// collision-resistance, so name(a) == name(b) => a == b holds outright, and the
// literal part keeps the artifacts dir debuggable. Upper-case letters escape too
// because the artifacts dir lives on NTFS, which FOLDS LETTER CASE - a name that
// kept "PQ" verbatim would still be the same file as "pq", which is the very
// defect this function used to have wearing a different hat. With [a-z0-9_-] the
// only literal set, `%XX` is the only upper-case a name can hold, a literal `%`
// itself becomes `%25`, and the mapping stays one-to-one even case-folded.
//
// An id whose encoding overflows maxArtifactIDBytes gets its readable part cut
// on an escape boundary and a digest of the WHOLE id appended: escaping triples
// bytes, and a name over Windows' 255-byte component limit fails the write
// outright, so the alternative was to lose the spill entirely. Such names are
// collision-resistant (sha256) rather than provably distinct - the only place
// this function is not information-theoretically injective, and it costs nothing
// to tell which branch produced a name: the digest tail is always there or not.
func artifactName(callID string, seq int) string {
	id := encodeArtifactID(callID)
	if id == "" {
		// The empty id is the one input that encodes to nothing, so it needs a
		// synthetic token - and a token made only of literals is reachable by a
		// literal id (id "seq7" would then share a name with the 7th empty-id
		// call, which is the same fold in a smaller hat). The dot in "seq.<n>" is
		// what closes that: the encoder escapes every dot, so no non-empty id can
		// produce a name holding a literal "." anywhere.
		id = fmt.Sprintf("seq.%d", seq)
	}
	if len(id) > maxArtifactIDBytes {
		id = cutEncodedID(id, maxArtifactIDBytes) + "-" + digestArtifactID(callID)
	}
	return "tool-output-" + id + ".txt"
}

const (
	// maxArtifactIDBytes bounds the encoded-id part of a name; with the fixed
	// digest tail this keeps the whole name well under a Windows path component.
	maxArtifactIDBytes = 96
	// artifactIDDigestBytes is how much sha256 goes into a long id's name.
	artifactIDDigestBytes = 8
)

func encodeArtifactID(callID string) string {
	const upperHex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(callID))
	for i := 0; i < len(callID); i++ {
		c := callID[i]
		if c == '_' || c == '-' || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperHex[c>>4])
		b.WriteByte(upperHex[c&0xF])
	}
	return b.String()
}

// cutEncodedID shortens an encoded id to at most limit bytes without splitting a
// %XX escape: half an escape would be a spelling no id encodes to, and the whole
// point of this function is that names are never ambiguous.
func cutEncodedID(enc string, limit int) string {
	cut := min(limit, len(enc))
	for cut > 0 && (enc[cut-1] == '%' || (cut >= 2 && enc[cut-2] == '%')) {
		cut--
	}
	return enc[:cut]
}

func digestArtifactID(callID string) string {
	sum := sha256.Sum256([]byte(callID))
	return hex.EncodeToString(sum[:artifactIDDigestBytes])
}

// writeFileExclusive creates path and writes data into it, FAILING with
// fs.ErrExist if anything already occupies the name - the O_EXCL the name has
// always promised (ticket 79: this was os.WriteFile, i.e. O_CREATE|O_TRUNC, so
// the second writer won silently and "Exclusive" was a lie). Deciding what a
// same-id retry means is the caller's job, not something O_TRUNC should
// half-decide by accident: see Prepare, which documents retry as last-writer-
// wins and performs the swap atomically.
func writeFileExclusive(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		// A half-written artifact would outlive this call: the name is exactly
		// the one a re-read of Spill.Path resolves to.
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	return f.Close()
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
