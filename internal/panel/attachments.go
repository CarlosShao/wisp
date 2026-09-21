package panel

// Attachment ingestion for the composer (ticket 92 AC#2).
//
// The shape of this file is decided by one rule: an attachment is a byte source
// the panel controls, and the panel is the least trusted part of the process
// (a compromised renderer is the threat model SPEC-06 §1 and ban #6 are written
// against). So nothing here trusts anything the caller asserts:
//
//   - the FILE NAME is rejected with the very guard the artifacts store applies
//     at write time (memory.ValidArtifactName, ticket 76). A name that spells a
//     path - "\" or "/" or ".." - is refused before a byte is read, so the
//     reject also means "nothing landed on disk";
//   - the DECLARED mime type is only ever a claim. The type that goes into the
//     message is what the leading bytes say (sniffAttachment), and a claim that
//     does not match the bytes is a refusal, not a re-label. An .exe renamed to
//     .png therefore produces "this is not a PNG", never a PNG;
//   - the SIZE is read from the source metadata before the source is opened, so
//     a 2 GB pick is refused without streaming it into memory. The read is then
//     bounded by a LimitReader anyway, because a reported size is also a claim;
//   - every refusal is an error with a reason the user can read. A silent drop
//     is the "lying config key" defect family (ticket 83): the user pressed a
//     button and the product swallowed their intent.
//
// Storage reuses the pinned artifacts face and nothing else: memory.Store's
// PutArtifact (injective name + real O_EXCL + recursive quota, tickets 76/79).
// There is no "user files" directory in this package and adding one is out of
// scope for the project.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

// MaxAttachmentBytes is the per-attachment ceiling the composer advertises and
// enforces. It is far below memory.ArtifactsQuotaBytes (500 MB) on purpose: one
// user pick must not be able to evict a session's worth of tool evidence under
// the LRU.
const MaxAttachmentBytes int64 = 64 << 20 // 64 MiB

// ErrAttachmentRejected is the sentinel behind every composer attachment
// refusal. It exists so a caller - or a test - can tell "the gate said no,
// nothing was stored" apart from "the store failed", without matching on
// message text.
var ErrAttachmentRejected = errors.New("panel: attachment rejected")

// ArtifactSink is the storage surface an attachment needs, and nothing more.
// *memory.Store satisfies it in production (ArtifactsDir already exists there,
// PutArtifact was added by this ticket so the artifacts name guard has exactly
// one owner).
type ArtifactSink interface {
	PutArtifact(ctx context.Context, name string, data []byte) error
	ArtifactsDir() string
}

// AttachmentSource is one file or clipboard blob as the panel presents it.
// Open is lazy for a reason: a source that the guards already refuse must never
// have its bytes read, which is also what lets a 2 GB pick be refused in
// constant time.
type AttachmentSource struct {
	// DisplayName is what the user sees; it is never used to build a path.
	DisplayName string
	// DeclaredMIME is the picker's claim (possibly ""). It is checked against
	// the sniffed type and never believed on its own.
	DeclaredMIME string
	// SizeBytes is the picker's claim; <= 0 means "unknown", which is read with
	// a bounded reader rather than believed.
	SizeBytes int64
	// Open returns a reader over the content. nil means "no bytes reachable".
	Open func() (io.ReadCloser, error)
}

// AttachmentRef is the reference that travels in the message. The bytes stay on
// the native side; this is what the panel and the agent both read.
type AttachmentRef struct {
	ID string `json:"id"`
	// Name is the sanitized display name (for humans only).
	Name string `json:"name"`
	// MIME is the SNIFFED type, never the declared one.
	MIME string `json:"mime"`
	// Kind is "image" or "video": what the agent can do with it.
	Kind string `json:"kind"`
	// SizeBytes is the number of bytes actually stored.
	SizeBytes int64 `json:"sizeBytes"`
	// Artifact is the bare file name inside the artifacts directory.
	Artifact string `json:"artifact"`
	// Stored is true only when the bytes are provably in the artifacts dir.
	Stored bool `json:"stored"`
	// Deduplicated is true when identical bytes were already stored, so this
	// pick added no new file (an honest "nothing happened", not a silent one).
	Deduplicated bool `json:"deduplicated"`
	// Reason is the refusal text shown to the user; empty when Stored.
	Reason string `json:"reason"`
}

// attachmentType is one whitelisted content type: the sniff rule that recognises
// it, the kind it renders as, and the extension its artifact name carries.
type attachmentType struct {
	mime  string
	kind  string
	ext   string
	match func(head []byte) bool
}

// attachmentTypes is the whitelist. It is a table, not a switch, so AC#2's
// "unsupported types must fail loudly" has one list to point at. Video is
// accepted as BYTES ONLY: semantic understanding of video is Q-28 and is not
// this ticket - the message therefore says what was stored, not what it means.
var attachmentTypes = []attachmentType{
	{mime: "image/png", kind: "image", ext: "png", match: func(h []byte) bool {
		return len(h) >= 8 && bytes.Equal(h[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	}},
	{mime: "image/jpeg", kind: "image", ext: "jpg", match: func(h []byte) bool {
		return len(h) >= 3 && h[0] == 0xff && h[1] == 0xd8 && h[2] == 0xff
	}},
	{mime: "image/gif", kind: "image", ext: "gif", match: func(h []byte) bool {
		return len(h) >= 6 && (bytes.Equal(h[:6], []byte("GIF87a")) || bytes.Equal(h[:6], []byte("GIF89a")))
	}},
	{mime: "image/webp", kind: "image", ext: "webp", match: func(h []byte) bool {
		return len(h) >= 12 && bytes.Equal(h[:4], []byte("RIFF")) && bytes.Equal(h[8:12], []byte("WEBP"))
	}},
	{mime: "video/mp4", kind: "video", ext: "mp4", match: matchesISOBaseMedia},
	{mime: "video/quicktime", kind: "video", ext: "mov", match: matchesISOBaseMedia},
}

// matchesISOBaseMedia recognises the ISO base media container (mp4/mov): a
// 4-byte size, the literal "ftyp", then a compatible brand from the small set
// the composer claims to read. A brand outside the set is an unsupported type,
// which is a refusal - not a guess.
func matchesISOBaseMedia(h []byte) bool {
	if len(h) < 12 || !bytes.Equal(h[4:8], []byte("ftyp")) {
		return false
	}
	switch string(h[8:12]) {
	case "isom", "avc1", "mp42", "M4V ", "qt  ", "mmp4":
		return true
	}
	return false
}

// AcceptedMIMETypes lists the whitelist for the composer UI, so the panel can
// tell the user what it can take instead of letting them discover it.
func AcceptedMIMETypes() []string {
	out := make([]string, 0, len(attachmentTypes))
	for _, t := range attachmentTypes {
		out = append(out, t.mime)
	}
	return out
}

// NameGuard is the artifacts package's bare-name rule, injected rather than
// imported: memory.ValidArtifactName is the same function the store applies at
// write time (ticket 76), so this package refuses a path-spelling name with the
// store's own verdict instead of a second copy of the rules that can drift.
type NameGuard func(name string) error

// AttachmentBroker turns one source into either a stored reference or a loud
// error. It is safe for concurrent use: it holds only immutable config.
type AttachmentBroker struct {
	sink  ArtifactSink
	guard NameGuard
	max   int64
}

// NewAttachmentBroker wires the broker to the artifacts face. A nil sink or a
// nil guard is a construction error, not a no-op broker: a broker that quietly
// accepted bytes with nowhere to put them, or with no name rule, is exactly the
// silent drop this file forbids.
func NewAttachmentBroker(sink ArtifactSink, guard NameGuard, max int64) (*AttachmentBroker, error) {
	if sink == nil {
		return nil, fmt.Errorf("%w: no artifact store attached", ErrAttachmentRejected)
	}
	if guard == nil {
		return nil, fmt.Errorf("%w: no artifact name guard attached", ErrAttachmentRejected)
	}
	if max <= 0 {
		max = MaxAttachmentBytes
	}
	return &AttachmentBroker{sink: sink, guard: guard, max: max}, nil
}

// MaxBytes reports the enforced per-attachment ceiling.
func (b *AttachmentBroker) MaxBytes() int64 { return b.max }

// Ingest validates and stores one attachment. The contract is total: either the
// bytes are in the artifacts directory and the returned reference says so, or
// an error wrapping ErrAttachmentRejected (or a store error) comes back with a
// reason and Stored=false.
func (b *AttachmentBroker) Ingest(ctx context.Context, src AttachmentSource) (AttachmentRef, error) {
	ref := AttachmentRef{
		ID:   "att-" + shortHash([]byte(src.DisplayName+"|"+src.DeclaredMIME)),
		Name: src.DisplayName,
		MIME: src.DeclaredMIME,
	}
	// Guard 1: a file name that spells a path is refused, and refused before
	// the source is opened, so the reject is also "nothing was written".
	if err := checkDisplayName(src.DisplayName, b.guard); err != nil {
		return b.refuse(ref, fmt.Sprintf("文件名 %q 不是一个裸文件名：%v", src.DisplayName, err))
	}
	// Guard 2/3: the declared size is read before the bytes are.
	if src.SizeBytes == 0 {
		return b.refuse(ref, fmt.Sprintf("%q 是 0 字节的空文件，没有内容可发送", src.DisplayName))
	}
	if src.SizeBytes > b.max {
		return b.refuse(ref, fmt.Sprintf("%q 有 %d 字节，超过附件上限 %d 字节；未读取任何内容",
			src.DisplayName, src.SizeBytes, b.max))
	}
	if src.Open == nil {
		return b.refuse(ref, fmt.Sprintf("%q 没有可读的内容来源", src.DisplayName))
	}
	rc, err := src.Open()
	if err != nil {
		return b.refuse(ref, fmt.Sprintf("%q 打不开：%v", src.DisplayName, err))
	}
	defer func() { _ = rc.Close() }()
	// The limit is one byte past the cap: reaching it proves the source lied
	// about its size, and the excess is never buffered.
	data, err := io.ReadAll(io.LimitReader(rc, b.max+1))
	if err != nil {
		return b.refuse(ref, fmt.Sprintf("%q 读取失败：%v", src.DisplayName, err))
	}
	if int64(len(data)) > b.max {
		return b.refuse(ref, fmt.Sprintf("%q 实际内容超过附件上限 %d 字节；只读取了前 %d 字节用于判定",
			src.DisplayName, b.max, b.max+1))
	}
	if len(data) == 0 {
		return b.refuse(ref, fmt.Sprintf("%q 声称有内容，实际是 0 字节", src.DisplayName))
	}
	// Guard 4: content sniffing, and the declared-vs-sniffed agreement.
	found := false
	var accepted attachmentType
	for _, t := range attachmentTypes {
		if t.match(data) {
			accepted, found = t, true
			break
		}
	}
	if !found {
		return b.refuse(ref, fmt.Sprintf("%q 的类型（%s）不受支持；本产品的附件只接受 %s",
			src.DisplayName, describeHead(data), strings.Join(AcceptedMIMETypes(), " / ")))
	}
	if src.DeclaredMIME != "" && !mimeAgrees(src.DeclaredMIME, accepted.mime) {
		return b.refuse(ref, fmt.Sprintf("%q 声明为 %s，但文件头表明它是 %s",
			src.DisplayName, src.DeclaredMIME, accepted.mime))
	}
	ref.MIME, ref.Kind, ref.Name = accepted.mime, accepted.kind, sanitizeDisplayName(src.DisplayName)

	// Guard 5: the name is a content hash, so it is injective in the bytes and
	// carries none of the caller's spelling.
	name := "attachment-" + shortHash(data) + "." + accepted.ext
	if err := b.guard(name); err != nil {
		return b.refuse(ref, fmt.Sprintf("生成的附件名未通过存储守卫：%v", err))
	}
	err = b.sink.PutArtifact(ctx, name, data)
	switch {
	case err == nil:
	case errors.Is(err, fs.ErrExist):
		ref.Deduplicated = true
	default:
		return ref, fmt.Errorf("panel: store attachment %q: %w", name, err)
	}
	ref.Artifact = name
	ref.SizeBytes = int64(len(data))
	ref.Stored = true
	ref.Reason = ""
	return ref, nil
}

// refuse is the single exit for a rejection: it fills the reason, clears
// everything a partial verdict might have leaked, and returns an error that
// wraps ErrAttachmentRejected so the caller can classify without text matching.
func (b *AttachmentBroker) refuse(ref AttachmentRef, reason string) (AttachmentRef, error) {
	ref.Stored = false
	ref.Deduplicated = false
	ref.Artifact = ""
	ref.SizeBytes = 0
	ref.Reason = reason
	return ref, fmt.Errorf("%w: %s", ErrAttachmentRejected, reason)
}

// checkDisplayName applies the store's name guard to a display name, after
// rejecting the two shapes a bare-name rule alone does not cover but the
// composer must: an empty name and control characters (a NUL inside a name is a
// truncation attack on every consumer that C-strings it later).
func checkDisplayName(name string, guard NameGuard) error {
	if name == "" {
		return errors.New("required (bare file names only)")
	}
	if strings.ContainsAny(name, "\x00\r\n\t") {
		return errors.New("contains control characters")
	}
	return guard(name)
}

// sanitizeDisplayName keeps a display name displayable: the store already
// refused anything spelling a path, so this only trims what a human would not
// want echoed back verbatim into a message.
func sanitizeDisplayName(name string) string {
	return filepathBaseLike(strings.TrimSpace(name))
}

// filepathBaseLike keeps the last path segment of a display name using EITHER
// separator, so a name that arrived as "C:\dir\photo.png" from a picker that
// lies about its separator still shows "photo.png".
func filepathBaseLike(s string) string {
	if i := strings.LastIndexAny(s, `/\`); i >= 0 {
		return s[i+1:]
	}
	return s
}

// describeHead names the bytes that were actually found, for the refusal text.
// It is deliberately blunt about an executable: "伪装成 .png 的 .exe" must read
// back as an executable, not as "unknown type".
func describeHead(data []byte) string {
	n := len(data)
	if n > 8 {
		n = 8
	}
	head := data[:n]
	switch {
	case bytes.HasPrefix(data, []byte("MZ")):
		return "可执行文件（MS-DOS/PE 头 \"MZ\"）"
	case bytes.HasPrefix(data, []byte{0x7f, 'E', 'L', 'F'}):
		return "可执行文件（ELF 头）"
	case bytes.HasPrefix(data, []byte("!<arch>")):
		return "归档/可执行文件（ar 头）"
	default:
		return fmt.Sprintf("未知类型，文件头 0x%x", head)
	}
}

// mimeAgrees compares two mime spellsings tolerantly (parameters, case).
func mimeAgrees(declared, actual string) bool {
	d := strings.ToLower(strings.TrimSpace(strings.Split(declared, ";")[0]))
	a := strings.ToLower(strings.TrimSpace(actual))
	if d == a {
		return true
	}
	// image/jpg is the one alias users type for image/jpeg.
	return d == "image/jpg" && a == "image/jpeg"
}

// shortHash is the first 16 hex chars of sha256 - long enough that the artifact
// name stays injective for practical purposes, short enough to read in a log.
func shortHash(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}
