package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// AC4 — tamper detection (SPEC-07 §6 "加载时校验", ticket 50).
//
// The install record pins the manifest hash (plugin_state.hash, already modelled
// by internal/memory.PluginState.Hash). At every load the bytes on disk are
// re-hashed and compared: mismatch => the load is refused AND an alert is
// emitted. Silent skip and "load anyway with a warning" are both forbidden by
// the contract, so refusal is the only non-alert-free path and there is no
// non-refusing path at all.
//
// Scope note: the command-plugin exe_hash check (D46, verified before every
// execution) belongs to ticket 52 and is not implemented here.

// HashAlgorithm is the only accepted manifest-hash prefix (C14 exe_hash and
// plugin_state.hash use the same "sha256:<hex>" shape).
const HashAlgorithm = "sha256:"

// Refusal sentinel errors. Callers branch on these instead of string matching.
var (
	// ErrTampered: the bytes do not hash to the pinned value.
	ErrTampered = errors.New("plugin: manifest tamper detected: hash mismatch vs install record")
	// ErrNoPinnedHash: nothing was pinned, so the bytes cannot be trusted.
	// An unverifiable manifest is refused, never loaded "as good enough".
	ErrNoPinnedHash = errors.New("plugin: no pinned manifest hash in the install record; refusing load")
	// ErrMalformedHash: the pinned value is not a sha256:<hex> string.
	ErrMalformedHash = errors.New("plugin: pinned manifest hash is malformed (expected sha256:<64 hex>)")
)

// ManifestHash returns the pinned-hash representation of raw manifest bytes:
// "sha256:<lowercase hex of the SHA-256 digest>".
func ManifestHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return HashAlgorithm + hex.EncodeToString(sum[:])
}

// VerifyBytes compares raw bytes against a pinned hash. It is the pure core of
// the load-time check (usable when the bytes came from somewhere other than a
// file, e.g. an already-downloaded install payload).
func VerifyBytes(raw []byte, pinnedHash string) error {
	want, err := normalizePinnedHash(pinnedHash)
	if err != nil {
		return err
	}
	if want == "" {
		return ErrNoPinnedHash
	}
	got := ManifestHash(raw)
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("%w: pinned %s, computed %s", ErrTampered, want, got)
	}
	return nil
}

func normalizePinnedHash(h string) (string, error) {
	h = strings.TrimSpace(h)
	if h == "" {
		return "", nil
	}
	if !strings.HasPrefix(strings.ToLower(h), HashAlgorithm) {
		return "", fmt.Errorf("%w: missing %q prefix in %q", ErrMalformedHash, HashAlgorithm, h)
	}
	digits := h[len(HashAlgorithm):]
	if len(digits) != sha256.Size*2 {
		return "", fmt.Errorf("%w: digest is %d hex chars, want %d", ErrMalformedHash, len(digits), sha256.Size*2)
	}
	if _, err := hex.DecodeString(digits); err != nil {
		return "", fmt.Errorf("%w: %v", ErrMalformedHash, err)
	}
	return HashAlgorithm + strings.ToLower(digits), nil
}

// TamperAlert describes one refused load. Kind distinguishes the two tamper
// surfaces so the panel/audit leg can render them differently.
type TamperAlert struct {
	Kind     string `json:"kind"` // "manifest" | "bytes"
	PluginID string `json:"plugin_id,omitempty"`
	Path     string `json:"path,omitempty"`
	Expected string `json:"expected,omitempty"` // pinned value as recorded
	Actual   string `json:"actual,omitempty"`   // recomputed from the bytes
	Reason   string `json:"reason"`
}

// Alerter receives tamper alerts. The load path calls Alert synchronously
// before returning the refusal, so an alert can never be lost by a caller that
// drops the error.
type Alerter interface {
	AlertTamper(a TamperAlert)
}

// LogAlerter is the default Alerter: one structured Error log line per alert
// with the stable event name "plugin_tamper".
type LogAlerter struct{}

// AlertTamper implements Alerter.
func (LogAlerter) AlertTamper(a TamperAlert) {
	slog.Error("plugin_tamper",
		"kind", a.Kind,
		"plugin_id", a.PluginID,
		"path", a.Path,
		"expected", a.Expected,
		"actual", a.Actual,
		"reason", a.Reason,
	)
}

// CaptureAlerter collects alerts in memory (tests, and the audit/panel leg
// until it wires a real sink). Safe for concurrent use.
type CaptureAlerter struct {
	mu     sync.Mutex
	alerts []TamperAlert
}

// AlertTamper implements Alerter.
func (c *CaptureAlerter) AlertTamper(a TamperAlert) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alerts = append(c.alerts, a)
}

// Alerts returns a copy of what was captured.
func (c *CaptureAlerter) Alerts() []TamperAlert {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]TamperAlert, len(c.alerts))
	copy(out, c.alerts)
	return out
}

// LoadOptions configures a LoadInstalled call.
type LoadOptions struct {
	// PinnedHash is plugin_state.hash for this plugin id. Empty => refused.
	PinnedHash string
	// PluginID labels alerts/logs only; it is not trusted for identity here
	// (the manifest's own id is checked against it when provided).
	PluginID string
	// Alerter receives tamper alerts; nil means LogAlerter.
	Alerter Alerter
}

// LoadInstalled reads an installed manifest from path, verifies it against the
// pinned hash, and only then parses/validates it (C14). On any tamper or
// unverifiable-hash condition it emits an alert and returns a refusal; the
// returned Manifest is the zero value in that case.
//
// The path is used exactly as given: no Clean/Abs here, path decisions belong
// to risk.PathResolver (AGENTS §1.2).
func LoadInstalled(path string, opt LoadOptions) (Manifest, error) {
	alerter := opt.Alerter
	if alerter == nil {
		alerter = LogAlerter{}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("plugin: read manifest %s: %w", path, err)
	}
	return VerifyAndParse(raw, path, opt)
}

// VerifyAndParse is LoadInstalled without the file read: tamper check first,
// then C14 parse. Split out so the install leg (ticket 50 AC3 / 51 / 52) can
// reuse the same gate on bytes it obtained elsewhere.
func VerifyAndParse(raw []byte, path string, opt LoadOptions) (Manifest, error) {
	alerter := opt.Alerter
	if alerter == nil {
		alerter = LogAlerter{}
	}

	pinned, err := normalizePinnedHash(opt.PinnedHash)
	if err != nil {
		alerter.AlertTamper(TamperAlert{
			Kind: "manifest", PluginID: opt.PluginID, Path: path,
			Expected: opt.PinnedHash, Actual: ManifestHash(raw),
			Reason: err.Error(),
		})
		return Manifest{}, err
	}
	if pinned == "" {
		alerter.AlertTamper(TamperAlert{
			Kind: "manifest", PluginID: opt.PluginID, Path: path,
			Actual: ManifestHash(raw),
			Reason: ErrNoPinnedHash.Error(),
		})
		return Manifest{}, ErrNoPinnedHash
	}

	got := ManifestHash(raw)
	if !strings.EqualFold(got, pinned) {
		alerter.AlertTamper(TamperAlert{
			Kind: "manifest", PluginID: opt.PluginID, Path: path,
			Expected: pinned, Actual: got,
			Reason: "installed manifest bytes do not match the hash pinned at install time (possible in-place swap)",
		})
		return Manifest{}, fmt.Errorf("%w: pinned %s, computed %s", ErrTampered, pinned, got)
	}

	m, err := ParseManifest(raw)
	if err != nil {
		return Manifest{}, fmt.Errorf("plugin: manifest %s: %w", path, err)
	}
	return m, nil
}
