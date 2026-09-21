package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Phase enumerates the visible stages of one Ensure run (ball Downloading
// ring: connecting -> downloading -> verifying -> extracting -> installing).
type Phase string

const (
	PhaseConnecting  Phase = "connecting"
	PhaseDownloading Phase = "downloading"
	PhaseVerifying   Phase = "verifying"
	PhaseExtracting  Phase = "extracting"
	PhaseInstalling  Phase = "installing"
	PhaseDone        Phase = "done"
)

// ProgressEvent is one tick of download progress (SPEC-04 §7.1: progress ring
// + percent). Percent is the OVERALL entry progress across artifacts.
type ProgressEvent struct {
	ModelID       string  `json:"model_id"`
	Phase         Phase   `json:"phase"`
	URL           string  `json:"url,omitempty"`  // candidate currently in use
	Attempt       int     `json:"attempt"`        // 1-based attempt on the current URL
	ArtifactIndex int     `json:"artifact_index"` // 0-based
	ArtifactCount int     `json:"artifact_count"`
	BytesDone     int64   `json:"bytes_done"`
	TotalBytes    int64   `json:"total_bytes"`
	Percent       float64 `json:"percent"`
}

// Options configures the Manager. See NewManager for invariants.
type Options struct {
	// DataDir is the model store root: installs land in <DataDir>/<id>/,
	// transient bytes in <DataDir>/staging/<id>/.
	DataDir string
	// Manifest is a signature-verified manifest (LoadSignedManifest). Never
	// construct one by hand for runtime use.
	Manifest *Manifest
	// Mirrors are the [models] mirror base URLs, tried before the manifest's
	// own URLs. A mirror serves bytes at <base>/<model-id>/<file name>;
	// hashes come only from the manifest (F3).
	Mirrors []string
	// VerifySignature is hard-required true (C29). This field exists only so
	// a `false` from config is a runtime error too (double-check behind the
	// config layer's hard rejection); there is no code path that disables
	// verification.
	VerifySignature bool
	// LocalOverride maps model-id -> local directory (sideload escape hatch,
	// SPEC-04 §7.1). Overridden models never touch the network; integrity is
	// still enforced (sha256 against the signed manifest).
	LocalOverride map[string]string
	// HTTPClient is injectable for tests; nil = a hardened default.
	HTTPClient *http.Client
	// Progress receives progress ticks; may be nil.
	Progress func(ProgressEvent)
	// Attempts is the per-URL retry budget (default 3, D37 backoff rows).
	Attempts int
	// BackoffBase is the exponential backoff base between attempts
	// (default 500ms; tests shrink it).
	BackoffBase time.Duration
}

// Manager downloads, verifies and installs models from the signed manifest.
// Safe for concurrent use; per-model Ensures are serialized by id.
type Manager struct {
	opts   Options
	client *http.Client

	mu      sync.Mutex
	cancels map[string]context.CancelFunc
	busy    map[string]bool
}

// stagingDirName is kept inside DataDir so partial bytes can never be mistaken
// for installed models and cancel cleanup has exactly one place to look.
const stagingDirName = "staging"

func NewManager(opts Options) (*Manager, error) {
	if opts.DataDir == "" {
		return nil, observeNew(ClassConfig, "models: DataDir required")
	}
	if opts.Manifest == nil {
		return nil, observeNew(ClassConfig, "models: signed Manifest required (C29: no manifest, no downloads)")
	}
	// C29 runtime double-check: config already hard-rejects
	// models.verify_signature=false; this is the second gate.
	if !opts.VerifySignature {
		return nil, observeNew(ClassConfig,
			"C29: signature verification is NOT disableable (models.verify_signature is read-only true)")
	}
	if opts.Attempts <= 0 {
		opts.Attempts = 3
	}
	if opts.BackoffBase <= 0 {
		opts.BackoffBase = 500 * time.Millisecond
	}
	if opts.HTTPClient == nil {
		opts.HTTPClient = defaultHTTPClient()
	}
	return &Manager{opts: opts, client: opts.HTTPClient, cancels: map[string]context.CancelFunc{}, busy: map[string]bool{}}, nil
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ResponseHeaderTimeout: 20 * time.Second,
		},
	}
}

// Ensure makes the model available and returns its install directory.
//
// Order of operations per call:
//  1. local_override: verify the given directory against the signed manifest
//     hashes and return it; ZERO network activity in this branch.
//  2. cache: if <DataDir>/<id> exists and every installed file hashes true,
//     reuse it (tamper surveillance on every hit).
//  3. download: mirror chain -> HTTP Range resume -> retry/backoff ->
//     sha256 verify -> archive extract -> atomic move into <DataDir>/<id>.
func (m *Manager) Ensure(ctx context.Context, id string, onProgress ...func(ProgressEvent)) (string, error) {
	entry, err := m.opts.Manifest.FindModel(id)
	if err != nil {
		return "", err
	}
	// P3 gate: a license-blocked model is never downloaded or installed.
	if !entry.Shippable() {
		return "", observeNew(ClassModel,
			fmt.Sprintf("model %s is P3-BLOCKED (license: %s); download refused, replacement decision pending", id, entry.License))
	}
	progress := func(ev ProgressEvent) {
		if m.opts.Progress != nil {
			m.opts.Progress(ev)
		}
		for _, cb := range onProgress {
			if cb != nil {
				cb(ev)
			}
		}
	}

	// 1. local_override: sideload with integrity check, no network.
	if dir, ok := m.opts.LocalOverride[id]; ok {
		if err := m.VerifyDir(entry, dir); err != nil {
			return "", observeNew(ClassModel,
				fmt.Sprintf("local_override %q failed integrity check: %v", id, err))
		}
		progress(ProgressEvent{ModelID: id, Phase: PhaseDone, ArtifactCount: entry.ArtifactCount(), Percent: 100})
		return dir, nil
	}

	// 2. cache hit with re-verification.
	installDir := filepath.Join(m.opts.DataDir, id)
	if err := m.VerifyDir(entry, installDir); err == nil {
		progress(ProgressEvent{ModelID: id, Phase: PhaseDone, ArtifactCount: entry.ArtifactCount(), Percent: 100})
		return installDir, nil
	}

	// 3. download. Serialize per id; Cancel targets the active Ensure.
	m.mu.Lock()
	if m.busy[id] {
		m.mu.Unlock()
		return "", observeNew(ClassModel, fmt.Sprintf("model %s already being ensured", id))
	}
	m.busy[id] = true
	cctx, cancel := context.WithCancel(ctx)
	m.cancels[id] = cancel
	m.mu.Unlock()
	defer func() {
		m.mu.Lock()
		delete(m.cancels, id)
		delete(m.busy, id)
		m.mu.Unlock()
		cancel()
	}()

	if err := m.downloadEntry(cctx, entry, progress); err != nil {
		if cctx.Err() != nil {
			m.removeStaging(id)
			return "", observeNew(ClassCancelled, fmt.Sprintf("download of %s cancelled; staging cleaned", id))
		}
		return "", err
	}
	if err := m.install(entry, installDir, progress); err != nil {
		return "", err
	}
	m.removeStaging(id)
	progress(ProgressEvent{ModelID: id, Phase: PhaseDone, ArtifactCount: entry.ArtifactCount(), Percent: 100})
	return installDir, nil
}

// VerifyInstalled re-checks an installed model against the signed manifest at
// the moment the model is handed over, not only when Ensure decided to return
// it (ticket 109, AC95-R1).
//
// The window this closes is a clock window, not a permission one: Ensure
// verifies per file *before* it returns, and the reader opens the files
// *after*, and this package deliberately does not seal them - the bytes are
// public and the wide directory is what lets a second instance under another
// account reuse a multi-gigabyte cache (ticket 95's ruling, which this does not
// touch). So anything that can write the install directory owns the span
// between those two events. The closure is therefore the second verification,
// and it is reachable without a *ModelEntry because that is exactly what a
// hand-off site does not have.
//
// It re-reads and re-hashes; it does not trust anything cached about a previous
// pass. Cost is one more full read of the installed files at hand-off, the same
// cost Ensure's cache-hit branch already pays per call.
func (m *Manager) VerifyInstalled(id string) error {
	entry, err := m.opts.Manifest.FindModel(id)
	if err != nil {
		return err
	}
	dir := filepath.Join(m.opts.DataDir, id)
	if override, ok := m.opts.LocalOverride[id]; ok {
		dir = override
	}
	if err := m.VerifyDir(entry, dir); err != nil {
		return observeWrap(ClassModel, "hand-off re-verification",
			fmt.Errorf("model %s installed files no longer match the signed manifest in %s: %w", id, dir, err))
	}
	return nil
}

// Cancel aborts an in-flight Ensure for id. The aborted run cleans its
// staging directory before returning (no partials outside staging, ever).
func (m *Manager) Cancel(id string) {
	m.mu.Lock()
	cancel, ok := m.cancels[id]
	m.mu.Unlock()
	if ok {
		cancel()
	}
}

// ---- download pipeline -----------------------------------------------------

func (m *Manager) downloadEntry(ctx context.Context, entry *ModelEntry, progress func(ProgressEvent)) error {
	staging := m.stagingDir(entry.ID)
	if err := os.MkdirAll(staging, 0o755); err != nil {
		return observeWrap(ClassModel, "create staging dir", err)
	}
	count := entry.ArtifactCount()
	var doneAll, totalAll int64
	var archived []int
	for i := 0; i < count; i++ {
		a := entry.Artifact(i)
		totalAll += a.SizeBytes
		done, err := m.downloadArtifact(ctx, entry, i, a, staging, progress, doneAll, totalAll)
		if err != nil {
			// Integrity failures poison everything staged: reject+delete.
			// Network failures keep the partial for resume.
			if cls, ok := observe.ClassOf(err); ok && cls == ClassModel {
				_ = os.RemoveAll(staging)
			}
			return err
		}
		doneAll += done
		if a.Archive != nil {
			archived = append(archived, i)
		}
	}

	// Extract archives (bytes were hash-verified inside downloadArtifact).
	for _, i := range archived {
		a := entry.Artifact(i)
		progress(ProgressEvent{
			ModelID: entry.ID, Phase: PhaseExtracting, ArtifactIndex: i, ArtifactCount: count,
			BytesDone: doneAll, TotalBytes: totalAll, Percent: percentOf(doneAll, totalAll),
		})
		archivePath := filepath.Join(staging, filepath.FromSlash(a.Path))
		if err := ExtractTarBz2(archivePath, staging); err != nil {
			_ = os.RemoveAll(staging)
			return observeNew(ClassModel, fmt.Sprintf("model %s: extraction failed: %v", entry.ID, err))
		}
		_ = os.Remove(archivePath) // archive bytes verified; keep only installed files
	}

	// Final gate before anything leaves staging: verify every installed file
	// hash inside staging (F3: mirror bytes are never trusted).
	for _, want := range entry.InstalledFiles() {
		progress(ProgressEvent{
			ModelID: entry.ID, Phase: PhaseVerifying, ArtifactCount: count,
			BytesDone: doneAll, TotalBytes: totalAll, Percent: percentOf(doneAll, totalAll),
		})
		p := filepath.Join(staging, filepath.FromSlash(want.Path))
		if err := verifyFileHash(p, want.SHA256, want.SizeBytes); err != nil {
			// A failed verification means staged bytes are wrong: reject+delete.
			_ = os.RemoveAll(staging)
			return observeNew(ClassModel, fmt.Sprintf("model %s: staged %s failed sha256: %v", entry.ID, want.Path, err))
		}
	}
	return nil
}

// downloadArtifact fetches one artifact into staging with the full candidate
// chain, Range resume and retry/backoff. Returns bytes counted for this
// artifact (its full size on success).
func (m *Manager) downloadArtifact(ctx context.Context, entry *ModelEntry, index int, a FileSpec,
	staging string, progress func(ProgressEvent), doneBefore, totalAll int64,
) (int64, error) {
	dst := filepath.Join(staging, filepath.FromSlash(a.Path))
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, observeWrap(ClassModel, "create staging subdir", err)
	}
	offset := int64(0)
	if st, err := os.Stat(dst); err == nil && st.Mode().IsRegular() {
		offset = st.Size()
	}
	// A staged file that already claims completion gets hash-checked right
	// here: pass = resume skip, fail = reject+delete and refetch.
	if offset == a.SizeBytes && offset > 0 {
		if err := verifyFileHash(dst, a.SHA256, a.SizeBytes); err == nil {
			progress(ProgressEvent{
				ModelID: entry.ID, Phase: PhaseDownloading, URL: currentCandidateLabel(entry, a),
				ArtifactIndex: index, ArtifactCount: entry.ArtifactCount(),
				BytesDone: doneBefore + offset, TotalBytes: totalAll, Percent: percentOf(doneBefore+offset, totalAll),
			})
			return offset, nil
		}
		_ = os.Remove(dst)
		offset = 0
	}

	candidates := BuildCandidates(m.opts.Mirrors, entry.ID, a.URLs)
	var lastErr error
	for _, cand := range candidates {
		for attempt := 1; attempt <= m.opts.Attempts; attempt++ {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			progress(ProgressEvent{
				ModelID: entry.ID, Phase: PhaseConnecting, URL: cand, Attempt: attempt,
				ArtifactIndex: index, ArtifactCount: entry.ArtifactCount(),
				BytesDone: doneBefore + offset, TotalBytes: totalAll, Percent: percentOf(doneBefore+offset, totalAll),
			})

			_, err := m.fetchOne(ctx, cand, dst, &offset, a, entry, index, progress, doneBefore, totalAll)
			if err == nil {
				return a.SizeBytes, nil
			}
			lastErr = err
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			if cls, ok := observe.ClassOf(err); ok && cls == ClassModel {
				// Integrity failure: bytes wrong. Staged file deleted inside
				// fetchOne; mirror bytes cannot be trusted -> next candidate.
				break
			}
			if attempt < m.opts.Attempts {
				time.Sleep(m.opts.BackoffBase << (attempt - 1)) // 1x, 2x, 4x...
			}
		}
	}
	return 0, lastErr
}

// fetchOne performs one HTTP exchange for the artifact at the current byte
// offset. 206 appends (resume), 200 restarts (server ignored Range), 416 with
// a complete offset means the file was already fully staged.
func (m *Manager) fetchOne(ctx context.Context, cand, dst string, offset *int64, a FileSpec,
	entry *ModelEntry, index int, progress func(ProgressEvent), doneBefore, totalAll int64,
) (int64, error) {
	wantSize := a.SizeBytes

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand, nil)
	if err != nil {
		return 0, observeWrap(ClassNetwork, "build request", err)
	}
	if *offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", *offset))
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return 0, observeWrap(ClassNetwork, "request", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusPartialContent && *offset > 0:
		// resume path: append
	case resp.StatusCode == http.StatusOK:
		// server ignored Range (or offset==0): restart from zero
		*offset = 0
	case resp.StatusCode == http.StatusRequestedRangeNotSatisfiable && *offset == wantSize && wantSize > 0:
		// Already fully staged (previous process completed the body but died
		// before verification). Verify here so nothing unverified passes on.
		if err := verifyFileHash(dst, a.SHA256, wantSize); err != nil {
			_ = os.Remove(dst)
			*offset = 0
			return 0, observeNew(ClassModel, fmt.Sprintf("staged file failed sha256: %v", err))
		}
		return *offset, nil
	default:
		return 0, observeNew(ClassNetwork, fmt.Sprintf("status %d for %s", resp.StatusCode, cand))
	}

	openFlag := os.O_WRONLY
	if *offset > 0 {
		openFlag |= os.O_APPEND
	} else {
		openFlag |= os.O_CREATE | os.O_TRUNC
	}
	f, err := os.OpenFile(dst, openFlag, 0o644)
	if err != nil {
		return 0, observeWrap(ClassModel, "open staging file", err)
	}
	defer f.Close()

	digester := sha256.New()
	// Hash only the FULL payload: on resume, prefix bytes were hashed in a
	// previous process, so re-feed them here to keep the digest whole.
	if *offset > 0 {
		if err := hashPrefix(digester, dst, *offset); err != nil {
			return 0, observeWrap(ClassModel, "hash staged prefix", err)
		}
	}

	counted := *offset
	buf := make([]byte, 32*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			wn, werr := f.Write(buf[:n])
			counted += int64(wn)
			digester.Write(buf[:wn])
			progress(ProgressEvent{
				ModelID: entry.ID, Phase: PhaseDownloading, URL: cand,
				ArtifactIndex: index, ArtifactCount: entry.ArtifactCount(),
				BytesDone: doneBefore + counted, TotalBytes: totalAll, Percent: percentOf(doneBefore+counted, totalAll),
			})
			if werr != nil {
				return 0, observeWrap(ClassModel, "write staging", werr)
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			// Keep the partial file: this is exactly the resume path. The
			// next attempt (or the next process) continues from the offset.
			*offset = counted
			return 0, observeWrap(ClassNetwork, "transfer interrupted (partial kept for resume)", rerr)
		}
	}

	// Completed body. Size then hash decide everything.
	if counted != wantSize {
		*offset = counted
		_ = os.Remove(dst) // wrong size cannot be resumed from: reject+delete
		*offset = 0
		return 0, observeNew(ClassModel, fmt.Sprintf("size mismatch for %s: got %d bytes, manifest pins %d", cand, counted, wantSize))
	}
	sum := hex.EncodeToString(digester.Sum(nil))
	if sum != a.SHA256 {
		_ = os.Remove(dst) // reject+delete tampered bytes
		*offset = 0
		return 0, observeNew(ClassModel, fmt.Sprintf("sha256 MISMATCH for %s (tampered mirror?): manifest pins %s, got %s; staged file deleted", cand, a.SHA256[:12], sum[:12]))
	}
	*offset = counted
	return counted, nil
}

// hashPrefix re-feeds the first n bytes of path into h (resume correctness).
func hashPrefix(h hash.Hash, path string, n int64) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.CopyN(h, f, n); err != nil {
		return err
	}
	return nil
}

// ---- candidates (mirror chain) ---------------------------------------------

// BuildCandidates assembles the failover chain for one artifact:
// [models] mirrors first (bytes at <base>/<model-id>/<file name>), then the
// manifest's own URLs in order, then transport-level fallbacks appended AFTER
// each official URL (ghfast.top for github.com release bytes - the P5 finding
// that direct github TLS fails intermittently from CN networks). Order is
// preserved and duplicates dropped.
func BuildCandidates(mirrors []string, modelID string, urls []string) []string {
	var out []string
	push := func(u string) {
		if u != "" && !contains(out, u) {
			out = append(out, u)
		}
	}
	name := ""
	if len(urls) > 0 {
		name = filepath.Base(strings.SplitN(urls[0], "?", 2)[0])
	}
	for _, base := range mirrors {
		push(strings.TrimRight(base, "/") + "/" + modelID + "/" + name)
	}
	for _, u := range urls {
		push(u)
		if fb, ok := transportFallback(u); ok {
			push(fb)
		}
	}
	return out
}

func currentCandidateLabel(entry *ModelEntry, a FileSpec) string {
	c := BuildCandidates(nil, entry.ID, a.URLs)
	if len(c) > 0 {
		return c[0]
	}
	return ""
}

// transportFallback maps an official URL to its domestic transport fallback.
func transportFallback(u string) (string, bool) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", false
	}
	switch parsed.Host {
	case "github.com":
		return "https://ghfast.top/" + u, true
	}
	return "", false
}

// ---- verification & install ------------------------------------------------

// VerifyDir checks that dir contains every installed file of the entry with
// exact sha256 and size (cache hits and local_override both use this; cheap:
// pure streaming hashes, no network).
func (m *Manager) VerifyDir(entry *ModelEntry, dir string) error {
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return fmt.Errorf("model dir %s missing", dir)
	}
	for _, want := range entry.InstalledFiles() {
		p := filepath.Join(dir, filepath.FromSlash(want.Path))
		if err := verifyFileHash(p, want.SHA256, want.SizeBytes); err != nil {
			return fmt.Errorf("%s: %w", want.Path, err)
		}
	}
	return nil
}

func verifyFileHash(path, wantSHA string, wantSize int64) error {
	st, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}
	if !st.Mode().IsRegular() {
		return fmt.Errorf("not a regular file")
	}
	if st.Size() != wantSize {
		return fmt.Errorf("size %d, manifest pins %d", st.Size(), wantSize)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != wantSHA {
		return fmt.Errorf("sha256 mismatch (pins %s..., got %s...)", wantSHA[:12], got[:12])
	}
	return nil
}

// install moves verified staged files into <DataDir>/<id>. Same-volume
// rename with a copy fallback; extraction layout is preserved.
func (m *Manager) install(entry *ModelEntry, installDir string, progress func(ProgressEvent)) error {
	staging := m.stagingDir(entry.ID)
	if progress != nil {
		progress(ProgressEvent{ModelID: entry.ID, Phase: PhaseInstalling, ArtifactCount: entry.ArtifactCount()})
	}
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return observeWrap(ClassModel, "create install dir", err)
	}
	for _, want := range entry.InstalledFiles() {
		src := filepath.Join(staging, filepath.FromSlash(want.Path))
		dst := filepath.Join(installDir, filepath.FromSlash(want.Path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return observeWrap(ClassModel, "create install subdir", err)
		}
		if err := os.Rename(src, dst); err != nil {
			if cpErr := copyFile(src, dst); cpErr != nil {
				return observeWrap(ClassModel, fmt.Sprintf("install %s", want.Path), err)
			}
			_ = os.Remove(src)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func (m *Manager) stagingDir(id string) string {
	return filepath.Join(m.opts.DataDir, stagingDirName, id)
}

func (m *Manager) removeStaging(id string) {
	_ = os.RemoveAll(m.stagingDir(id))
}

func percentOf(done, total int64) float64 {
	if total <= 0 {
		return 0
	}
	p := float64(done) / float64(total) * 100
	if p > 100 {
		p = 100
	}
	return p
}
