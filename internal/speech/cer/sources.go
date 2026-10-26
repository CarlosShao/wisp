package cer

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ManifestName is the one manifest of the committed baseline set, at
// testdata/asr-baseline/sources.tsv. One manifest covers both conditions; the
// file column carries the condition as its first path segment (SPEC-01:86-87).
const ManifestName = "sources.tsv"

// manifestHeader is the frozen column order. It is a header line, not a
// comment: the reader refuses a file whose first data line is not exactly this,
// so a column added or reordered by the corpus leg fails loudly instead of
// silently shifting every sha256 into the speaker column.
const manifestHeader = "file\tsha256\tspeaker\tcondition\torigin"

// sha256HexLen is the length of a lowercase hex SHA-256 digest.
const sha256HexLen = 64

// Source is one manifest row: which audio file, which digest, which speaker,
// which acoustic condition, and where the audio came from.
type Source struct {
	// File is relative to the baseline root, forward-slash separated, ending in
	// .wav; its first segment is the condition directory.
	File string
	// SHA256 is the lowercase hex digest of that file. OpenSLR publishes none,
	// so this manifest is the only integrity record the repo has
	// (cer-a1 section 1: "no official sha256").
	SHA256 string
	// Speaker is the corpus speaker id; the near_clean and noisy_far sets must
	// not share speakers (cer-a1 section 4).
	Speaker string
	// Condition is the acoustic condition, and must equal File's first segment.
	Condition Condition
	// Origin is the provenance string, e.g. "AISHELL-1:test:S0002/BAC009S0002W0122".
	Origin string
}

// The manifest's failure modes. Each one is named because the caller that has
// to act on it is a CI gate: "the set is gone" and "the set is wrong" must not
// look alike.
var (
	ErrManifestMissing   = errors.New("cer: baseline manifest is missing")
	ErrManifestHeader    = errors.New("cer: baseline manifest header is not the frozen one")
	ErrManifestRow       = errors.New("cer: baseline manifest row is malformed")
	ErrManifestEmpty     = errors.New("cer: baseline manifest has no data rows")
	ErrManifestDuplicate = errors.New("cer: baseline manifest lists one file twice")
	ErrAudioMissing      = errors.New("cer: baseline audio file is missing")
	ErrReferenceMissing  = errors.New("cer: baseline reference text is missing")
	ErrDigestMismatch    = errors.New("cer: baseline audio sha256 does not match the manifest")
)

// ParseManifest reads a manifest. Blank lines and lines starting with '#' are
// ignored; CRLF is accepted because the file may be edited on this platform.
// A header-only or comment-only file is ErrManifestEmpty, never an empty slice
// with a nil error.
func ParseManifest(r io.Reader) ([]Source, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	var (
		out     []Source
		line    int
		header  bool
		seen    = make(map[string]bool)
		dataRow int
	)
	for sc.Scan() {
		line++
		text := strings.TrimSuffix(sc.Text(), "\r")
		trimmed := strings.TrimSpace(text)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !header {
			if text != manifestHeader {
				return nil, fmt.Errorf("%w: line %d is %q, want %q",
					ErrManifestHeader, line, text, manifestHeader)
			}
			header = true
			continue
		}
		src, err := parseRow(text, line)
		if err != nil {
			return nil, err
		}
		if seen[src.File] {
			return nil, fmt.Errorf("%w: %s (line %d)", ErrManifestDuplicate, src.File, line)
		}
		seen[src.File] = true
		out = append(out, src)
		dataRow++
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("cer: read %s: %w", ManifestName, err)
	}
	if !header {
		// No header at all means the file is empty or comment-only: that is the
		// "the corpus is gone" case, not the "the corpus is wrong" case. A
		// present-but-different header is caught above.
		return nil, fmt.Errorf("%w: no header line %q and no data", ErrManifestEmpty, manifestHeader)
	}
	if dataRow == 0 {
		return nil, fmt.Errorf("%w: header %q then nothing", ErrManifestEmpty, manifestHeader)
	}
	return out, nil
}

// parseRow validates one tab-separated data row.
func parseRow(text string, line int) (Source, error) {
	f := strings.Split(text, "\t")
	if len(f) != 5 {
		return Source{}, fmt.Errorf("%w: line %d has %d columns, want 5 (%s)",
			ErrManifestRow, line, len(f), manifestHeader)
	}
	src := Source{
		File:      f[0],
		SHA256:    f[1],
		Speaker:   f[2],
		Condition: Condition(f[3]),
		Origin:    f[4],
	}
	if err := validateSource(src); err != nil {
		return Source{}, fmt.Errorf("%w: line %d: %v", ErrManifestRow, line, err)
	}
	return src, nil
}

// validateSource enforces the schema the corpus leg has to satisfy.
func validateSource(s Source) error {
	if s.File == "" {
		return errors.New("file is empty")
	}
	if strings.Contains(s.File, `\`) {
		return fmt.Errorf("file %q uses a backslash; the manifest is forward-slash only", s.File)
	}
	if path.IsAbs(s.File) || strings.HasPrefix(s.File, "/") {
		return fmt.Errorf("file %q is absolute; it is relative to the baseline root", s.File)
	}
	if cleaned := path.Clean(s.File); cleaned != s.File {
		return fmt.Errorf("file %q is not in clean form (path.Clean gives %q)", s.File, cleaned)
	}
	if !strings.HasSuffix(s.File, ".wav") {
		return fmt.Errorf("file %q does not end in .wav", s.File)
	}
	if s.Condition != NearClean && s.Condition != NoisyFar {
		return fmt.Errorf("condition %q is neither %q nor %q", string(s.Condition), NearClean, NoisyFar)
	}
	if dir := strings.SplitN(s.File, "/", 2)[0]; dir != string(s.Condition) {
		return fmt.Errorf("file %q is not under the %q directory; SPEC-01:86-87 freezes those names",
			s.File, string(s.Condition))
	}
	if !validSHA256(s.SHA256) {
		return fmt.Errorf("sha256 %q is not 64 lowercase hex characters", s.SHA256)
	}
	if s.Speaker == "" {
		return errors.New("speaker is empty")
	}
	if s.Origin == "" {
		return errors.New("origin is empty (the corpus and utterance id belong here)")
	}
	return nil
}

// validSHA256 accepts lowercase hex only: an uppercase digest is a different
// string, and a manifest that half-matches would produce digest mismatches that
// look like corrupt audio.
func validSHA256(s string) bool {
	if len(s) != sha256HexLen {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// ReadManifest opens <root>/sources.tsv. A missing file wraps BOTH
// ErrManifestMissing and fs.ErrNotExist, so a caller can tell "the corpus has
// not landed yet" (errors.Is(err, fs.ErrNotExist)) from "the corpus is broken"
// without string matching.
func ReadManifest(root string) ([]Source, error) {
	p := filepath.Join(root, ManifestName)
	f, err := os.Open(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %w (the committed baseline set is a separate leg; "+
				"an absent set must never be read as an empty, passing one)", ErrManifestMissing, err)
		}
		return nil, fmt.Errorf("cer: open %s: %w", p, err)
	}
	defer f.Close()
	srcs, err := ParseManifest(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	return srcs, nil
}

// LoadSet loads one acoustic condition of the committed baseline: the manifest,
// each listed audio file's presence, and each reference transcription. It does
// not decode audio and does not hash it - hypotheses come from the ASR runner,
// which is a different slice of ticket 15.
//
// Every failure is named. In particular a condition with no rows, a missing
// wav and a missing/empty .txt are three different errors, because a silently
// shortened set is indistinguishable from a better recogniser.
func LoadSet(root string, cond Condition) (Set, error) {
	if _, err := cond.Gate(); err != nil {
		return Set{}, err
	}
	srcs, err := ReadManifest(root)
	if err != nil {
		return Set{}, err
	}
	set := Set{Condition: cond}
	for _, src := range srcs {
		if src.Condition != cond {
			continue
		}
		wav := filepath.Join(root, filepath.FromSlash(src.File))
		if _, err := os.Stat(wav); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return Set{}, fmt.Errorf("%w: %s (manifest lists %s)", ErrAudioMissing, wav, src.File)
			}
			return Set{}, fmt.Errorf("cer: stat %s: %w", wav, err)
		}
		txt := refPath(wav)
		raw, err := os.ReadFile(txt)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return Set{}, fmt.Errorf("%w: %s (the reference is the sibling .txt of %s)",
					ErrReferenceMissing, txt, src.File)
			}
			return Set{}, fmt.Errorf("cer: read %s: %w", txt, err)
		}
		ref := refText(raw)
		if ref == "" {
			return Set{}, fmt.Errorf("%w: %s", ErrEmptyReference, txt)
		}
		set.Utterances = append(set.Utterances, Utterance{ID: src.File, Ref: ref})
	}
	if len(set.Utterances) == 0 {
		return Set{}, fmt.Errorf("%w: %s lists no %q rows",
			ErrNoUtterances, filepath.Join(root, ManifestName), string(cond))
	}
	return set, nil
}

// refPath is the reference-transcript path for a wav: the same stem, .txt.
func refPath(wav string) string { return strings.TrimSuffix(wav, ".wav") + ".txt" }

// refText strips a UTF-8 BOM and surrounding whitespace. Everything else is the
// normaliser's business; this only guards against an editor's newline or BOM
// becoming a reference character.
func refText(b []byte) string {
	return strings.TrimSpace(strings.TrimPrefix(string(b), "\ufeff"))
}

// VerifyManifest checks every listed file's digest. It is separate from LoadSet
// because it is the one step that must touch all the bytes, and a gate run
// should not pay for it on every invocation; the corpus leg pins the digests
// once, and CI can run this as its own step.
func VerifyManifest(root string) error {
	srcs, err := ReadManifest(root)
	if err != nil {
		return err
	}
	for _, src := range srcs {
		if err := src.Verify(root); err != nil {
			return err
		}
	}
	return nil
}

// Verify checks one source's digest against the file on disk.
func (s Source) Verify(root string) error {
	p := filepath.Join(root, filepath.FromSlash(s.File))
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("cer: verify %s: %w", p, err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("cer: hash %s: %w", p, err)
	}
	got := hex.EncodeToString(h.Sum(nil))
	if got != s.SHA256 {
		return fmt.Errorf("%w: %s hashes to %s, manifest says %s",
			ErrDigestMismatch, s.File, got, s.SHA256)
	}
	return nil
}
