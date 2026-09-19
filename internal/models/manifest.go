package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manifest is the C29 signed model manifest (SPEC-04 §7.2, D33/F3). The JSON
// lives at repo path models/manifest.json and is signed by
// models/manifest.minisig; hashes NEVER come from a mirror (F3: same-source
// hash = zero auth value). Every entry carries the frozen contract fields
// {id, purpose, urls[], sha256, size_bytes, license, quant}.
//
// Field semantics for the two entry shapes:
//
//   - Single-artifact entry (Files == nil): urls[] lists the download
//     candidates for ONE artifact; sha256/size_bytes describe it; the
//     installed file name is the base name of the first URL.
//
//   - Multi-artifact entry (Files != nil): each Files[i] is one downloaded
//     artifact with its own urls[]/sha256/size_bytes; entry-level urls[] is
//     the flattened index of every artifact's official source (documentation
//     field, not the download list); entry-level sha256 = sha256 over the
//     byte-concatenation of all artifacts in Files order ("combined digest")
//     and size_bytes = their total. Authoring-time invariant, recomputed by
//     tools/signmodels and pinned by TestManifestCombinedDigest.
//
// An artifact may be an archive (Archive != nil): the archive bytes are
// verified by its own sha256, then extracted; Archive.Files pins the sha256
// of every extracted file that is installed, which is what integrity
// re-checks (cache hits, local_override) verify on disk.
type Manifest struct {
	ManifestVersion int          `json:"manifest_version"`
	Models          []ModelEntry `json:"models"`
}

// ModelEntry is one model of the C29 manifest.
type ModelEntry struct {
	ID        string     `json:"id"`
	Purpose   string     `json:"purpose"`
	Status    string     `json:"status,omitempty"` // "" or "ok" = shippable; "blocked-p3" = license block (P3), downloads refused
	URLs      []string   `json:"urls"`
	SHA256    string     `json:"sha256"`
	SizeBytes int64      `json:"size_bytes"`
	License   string     `json:"license"`
	Quant     string     `json:"quant"`
	Files     []FileSpec `json:"files,omitempty"`
}

// Entry statuses (P3: a non-commercial license finding blocks the model from
// shipping; the entry stays in the manifest with its hashes for provenance).
const (
	StatusOK        = "ok"
	StatusBlockedP3 = "blocked-p3"
)

// Shippable reports whether the model may be downloaded/installed.
func (e *ModelEntry) Shippable() bool {
	return e.Status == "" || e.Status == StatusOK
}

// FileSpec is one downloadable artifact of a model.
type FileSpec struct {
	Path      string       `json:"path"`   // install path relative to models/<id>/; for archives, the archive file name
	URLs      []string     `json:"urls"`   // candidates, tried in order (mirror failover)
	SHA256    string       `json:"sha256"` // of the artifact bytes as served
	SizeBytes int64        `json:"size_bytes"`
	Archive   *ArchiveSpec `json:"archive,omitempty"`
}

// ArchiveSpec describes an archived artifact: verified as bytes, then
// extracted. Format is "tar.bz2" (the only format the official sherpa-onnx
// releases use; Go stdlib decompresses it, zero new dependencies).
type ArchiveSpec struct {
	Format string          `json:"format"`
	Files  []ExtractedFile `json:"files"`
}

// ExtractedFile pins one file produced by extraction. Paths are flattened:
// if the archive root is a single directory it is stripped.
type ExtractedFile struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

// Purposes is the closed vocabulary of model purposes (aligns with the speech
// pipeline seams; ticket 15 consumes these).
var Purposes = []string{"kws", "vad", "asr-streaming", "asr-offline", "punctuation", "tts"}

// Quants is the closed vocabulary of quant labels. Values are factual, not
// aspirational: matcha-zh has no official int8 (S0 spike §4.7), so its entry
// says fp32.
var Quants = []string{"int8", "fp32"}

// ErrManifestInvalid is wrapped by every structural manifest rejection.
var ErrManifestInvalid = fmt.Errorf("models: manifest invalid")

func decodeManifestJSON(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}
	return &m, nil
}

// ParseManifest decodes and validates manifest JSON structure. Signature
// verification happens in LoadSignedManifest BEFORE any network use (C29:
// verification is offline).
func ParseManifest(data []byte) (*Manifest, error) {
	m, err := decodeManifestJSON(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrManifestInvalid, err)
	}
	if m.ManifestVersion != 1 {
		return nil, fmt.Errorf("%w: manifest_version %d, want 1", ErrManifestInvalid, m.ManifestVersion)
	}
	if len(m.Models) == 0 {
		return nil, fmt.Errorf("%w: zero models", ErrManifestInvalid)
	}
	seen := map[string]bool{}
	for i := range m.Models {
		e := &m.Models[i]
		if e.ID == "" {
			return nil, fmt.Errorf("%w: model[%d] empty id", ErrManifestInvalid, i)
		}
		if seen[e.ID] {
			return nil, fmt.Errorf("%w: duplicate model id %q", ErrManifestInvalid, e.ID)
		}
		seen[e.ID] = true
		if !contains(Purposes, e.Purpose) {
			return nil, fmt.Errorf("%w: %s purpose %q not in %v", ErrManifestInvalid, e.ID, e.Purpose, Purposes)
		}
		if !contains(Quants, e.Quant) {
			return nil, fmt.Errorf("%w: %s quant %q not in %v", ErrManifestInvalid, e.ID, e.Quant, Quants)
		}
		if e.License == "" {
			return nil, fmt.Errorf("%w: %s license empty (P3: every entry records its license finding)", ErrManifestInvalid, e.ID)
		}
		if !isHex64(e.SHA256) {
			return nil, fmt.Errorf("%w: %s sha256 %q not 64 hex chars", ErrManifestInvalid, e.ID, e.SHA256)
		}
		if e.Status != "" && e.Status != StatusOK && e.Status != StatusBlockedP3 {
			return nil, fmt.Errorf("%w: %s status %q not in {ok, blocked-p3}", ErrManifestInvalid, e.ID, e.Status)
		}
		if e.SizeBytes < 0 {
			return nil, fmt.Errorf("%w: %s negative size", ErrManifestInvalid, e.ID)
		}
		// Files consistency (multi-artifact entries).
		if len(e.Files) > 0 {
			for j := range e.Files {
				f := &e.Files[j]
				if err := validRelPath(f.Path); err != nil {
					return nil, fmt.Errorf("%w: %s files[%d].path: %v", ErrManifestInvalid, e.ID, j, err)
				}
				if !isHex64(f.SHA256) {
					return nil, fmt.Errorf("%w: %s files[%d] sha256 not 64 hex", ErrManifestInvalid, e.ID, j)
				}
				if f.SizeBytes < 0 {
					return nil, fmt.Errorf("%w: %s files[%d] negative size", ErrManifestInvalid, e.ID, j)
				}
				if f.Archive != nil {
					if f.Archive.Format != "tar.bz2" {
						return nil, fmt.Errorf("%w: %s files[%d] archive format %q", ErrManifestInvalid, e.ID, j, f.Archive.Format)
					}
					if len(f.Archive.Files) == 0 {
						return nil, fmt.Errorf("%w: %s files[%d] archive lists no files", ErrManifestInvalid, e.ID, j)
					}
					for k := range f.Archive.Files {
						x := &f.Archive.Files[k]
						if err := validRelPath(x.Path); err != nil {
							return nil, fmt.Errorf("%w: %s files[%d] archive file[%d]: %v", ErrManifestInvalid, e.ID, j, k, err)
						}
						if !isHex64(x.SHA256) {
							return nil, fmt.Errorf("%w: %s files[%d] archive file[%d] sha256 not 64 hex", ErrManifestInvalid, e.ID, j, k)
						}
					}
				}
			}
		}
	}
	return m, nil
}

// ArtifactCount returns the number of downloadable artifacts for the entry
// (1 for single-artifact entries).
func (e *ModelEntry) ArtifactCount() int {
	if len(e.Files) == 0 {
		return 1
	}
	return len(e.Files)
}

// ImplicitFile renders a single-artifact entry as its FileSpec, so download
// and verification logic has exactly one code path. The install file name is
// the base name of the first URL.
func (e *ModelEntry) ImplicitFile() FileSpec {
	name := "model.bin"
	if len(e.URLs) > 0 {
		name = filepath.Base(strings.SplitN(e.URLs[0], "?", 2)[0])
	}
	return FileSpec{Path: name, URLs: e.URLs, SHA256: e.SHA256, SizeBytes: e.SizeBytes}
}

// Artifact returns the i-th download artifact.
func (e *ModelEntry) Artifact(i int) FileSpec {
	if len(e.Files) == 0 {
		return e.ImplicitFile()
	}
	return e.Files[i]
}

// InstalledFiles lists every file expected on disk after installation:
// archive members for archived artifacts, artifact files otherwise.
func (e *ModelEntry) InstalledFiles() []ExtractedFile {
	var out []ExtractedFile
	for i := range e.ArtifactCount() {
		a := e.Artifact(i)
		if a.Archive != nil {
			out = append(out, a.Archive.Files...)
			continue
		}
		out = append(out, ExtractedFile{Path: a.Path, SHA256: a.SHA256, SizeBytes: a.SizeBytes})
	}
	return out
}

// CombinedDigest computes the entry-level digest contract: sha256 over the
// byte-concatenation of all artifacts in Files order. The authoring tool
// asserts entry.SHA256 equals this; see Manifest doc comment.
func (e *ModelEntry) CombinedDigest(artifactBytes [][]byte) string {
	h := sha256.New()
	for _, b := range artifactBytes {
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// LoadSignedManifest reads manifest+signature from disk and verifies the
// minisign signature against pubKey BEFORE the manifest is parsed. This is
// the C29 gate: a manifest that does not verify never becomes model
// metadata, so no URL in it is ever contacted.
func LoadSignedManifest(manifestPath, sigPath string, pubKey string) (*Manifest, error) {
	pub, err := ParseMinisignPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("models: embedded C29 public key invalid: %w", err)
	}
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, observeWrap(ClassModel, fmt.Sprintf("read manifest %s", manifestPath), err)
	}
	sigBytes, err := os.ReadFile(sigPath)
	if err != nil {
		return nil, observeWrap(ClassModel, fmt.Sprintf("read manifest signature %s", sigPath), err)
	}
	if err := VerifyMinisignSignature(pub, manifestBytes, sigBytes); err != nil {
		return nil, observeNew(ClassModel, fmt.Sprintf("C29 manifest signature verification FAILED (%s): %v", manifestPath, err))
	}
	m, err := ParseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// FindModel returns the entry with the given id.
func (m *Manifest) FindModel(id string) (*ModelEntry, error) {
	for i := range m.Models {
		if m.Models[i].ID == id {
			return &m.Models[i], nil
		}
	}
	return nil, observeNew(ClassModel, fmt.Sprintf("model id %q not in signed manifest", id))
}

// validRelPath rejects absolute paths, escapes, and separators that could
// escape the model directory (zip-slip guard, C26 spirit applied to archives).
func validRelPath(p string) error {
	if p == "" {
		return fmt.Errorf("empty path")
	}
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") || strings.Contains(p, `\`) {
		return fmt.Errorf("not a clean relative path: %q", p)
	}
	if filepath.IsAbs(p) || filepath.VolumeName(p) != "" {
		return fmt.Errorf("volume path rejected: %q", p)
	}
	clean := filepath.ToSlash(filepath.Clean(p))
	if clean != p || strings.HasPrefix(clean, "../") || clean == ".." || strings.Contains(clean, "/../") {
		return fmt.Errorf("path escapes model dir: %q", p)
	}
	return nil
}

func isHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// ResolveManifestPath finds models/manifest.json (and .minisig) for the
// current process, in priority order: explicit override, WISP_MODELS_MANIFEST
// env, <exe dir>/models/, <cwd>/models/ (go run from the repo root). Used by
// the runtime wiring; tests pass explicit paths.
func ResolveManifestPath(explicit string) (manifest, sig string, err error) {
	if explicit != "" {
		return explicit, explicit + ".minisig", nil
	}
	if env := os.Getenv("WISP_MODELS_MANIFEST"); env != "" {
		return env, env + ".minisig", nil
	}
	var dirs []string
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(exe), "models"))
	}
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(wd, "models"))
	}
	for _, d := range dirs {
		p := filepath.Join(d, "manifest.json")
		if fileExists(p) {
			return p, p + ".minisig", nil
		}
	}
	return "", "", observeNew(ClassModel, "signed manifest not found (looked in WISP_MODELS_MANIFEST, <exe>/models, <cwd>/models)")
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}
