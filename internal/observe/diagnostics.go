package observe

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Diagnostics bundle (ticket 08 API stub; ticket 45 completes the UX wiring,
// ticket 40 wires the panel entry point).
//
// The bundle collects, fail-closed:
//   - the rolling logs (already redacted at write time by the non-disableable
//     redactor - the bundle never writes a second, weaker format),
//   - a config snapshot re-passed through the redactor (refs only by design;
//     a snapshot in which an inline secret shape survives ABORTS the bundle),
//   - the SLO sampler snapshot JSON to attach (per acceptance: "sampler data
//     attachable to a bundle"),
//   - version info and the DEFERRED register items supplied by the caller.
//
// Anything optional that is missing is recorded in manifest.json with its
// reason - the manifest never lies about what is inside.

// DeferredItem is one DEFERRED register entry for deferred.json.
type DeferredItem struct {
	ID      string `json:"id"`     // e.g. "D-xx" anchor
	Ticket  string `json:"ticket"` // owning ticket
	Summary string `json:"summary"`
}

// BundleOptions configures BuildDiagnosticsBundle.
type BundleOptions struct {
	// OutPath is the target zip (created/overwritten).
	OutPath string
	// LogDir is the rolling log directory (wisp-*.jsonl files).
	LogDir string
	// ConfigPath is the config.toml (refs only, re-redacted defensively).
	ConfigPath string
	// RedactPaths mirrors [privacy] redact_paths for the snapshot copy.
	RedactPaths bool
	// SLOSnapshot is the sampler report JSON to attach (may be nil).
	SLOSnapshot []byte
	// SLOSnapshotName overrides the default entry name (e.g. per-state).
	SLOSnapshotName string
	// Version is the buildinfo block from the caller (cmd wiring owns
	// internal/buildinfo; this package stays dependency-clean).
	Version map[string]string
	// Deferred is the DEFERRED register; nil/empty records the register as
	// not-yet-embedded (ticket 45), never as "none exist".
	Deferred []DeferredItem
	// MaxLogBytes caps each copied log file (skip larger, record it).
	MaxLogBytes int64
	// Now is injectable for tests.
	Now time.Time
}

// BuildDiagnosticsBundle writes the zip and returns its path.
func BuildDiagnosticsBundle(o BundleOptions) (string, error) {
	if o.OutPath == "" {
		return "", New(ClassConfig, "observe: bundle output path is empty")
	}
	if o.MaxLogBytes <= 0 {
		o.MaxLogBytes = 10 << 20
	}
	if o.Now.IsZero() {
		o.Now = NowWallUTC()
	}

	f, err := os.Create(o.OutPath)
	if err != nil {
		return "", Wrap(ClassResource, err, "observe: create bundle zip")
	}
	complete := false
	defer func() {
		if !complete {
			// Fail-closed: an aborted bundle leaves no zip behind.
			f.Close()
			_ = os.Remove(o.OutPath)
		}
	}()
	zw := zip.NewWriter(f)

	var manifest strings.Builder
	manifest.WriteString("{\n  \"created_at\": \"" + WallTimestampUTC(o.Now) + "\",\n  \"contents\": [\n")
	first := true
	addManifest := func(entry, detail string) {
		if !first {
			manifest.WriteString(",\n")
		}
		first = false
		fmt.Fprintf(&manifest, "    {\"entry\": %q, \"detail\": %q}", entry, detail)
	}

	// 1. Logs (redacted at rest; re-redact lines defensively anyway).
	if o.LogDir != "" {
		entries, lerr := os.ReadDir(o.LogDir)
		if lerr != nil {
			return "", Wrap(ClassResource, lerr, "observe: read log dir for bundle")
		}
		copied, skipped := 0, 0
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasPrefix(name, logFilePrefix) || !strings.HasSuffix(name, logFileExt) {
				continue
			}
			full := filepath.Join(o.LogDir, name)
			data, rerr := os.ReadFile(full)
			if rerr != nil {
				skipped++
				continue
			}
			if int64(len(data)) > o.MaxLogBytes {
				skipped++
				continue
			}
			red := Redactor{RedactPaths: o.RedactPaths}
			clean := red.String(string(data))
			if err := writeZipEntry(zw, "logs/"+name, []byte(clean)); err != nil {
				return "", err
			}
			copied++
		}
		addManifest("logs/", fmt.Sprintf("%d file(s) copied, %d skipped (missing/oversized)", copied, skipped))
	} else {
		addManifest("logs/", "not collected: log dir not provided")
	}

	// 2. Config snapshot (fail-closed on surviving secret shapes).
	if o.ConfigPath != "" {
		raw, rerr := os.ReadFile(o.ConfigPath)
		if rerr != nil {
			return "", Wrap(ClassResource, rerr, "observe: read config for bundle")
		}
		red := Redactor{RedactPaths: o.RedactPaths}
		clean, ok := red.ConfigSnapshot(string(raw))
		if !ok {
			return "", New(ClassInternal,
				"observe: config snapshot still contains key-shaped material after redaction; bundle aborted (fail-closed)")
		}
		if err := writeZipEntry(zw, "config.redacted.toml", []byte(clean)); err != nil {
			return "", err
		}
		addManifest("config.redacted.toml", "config.toml re-passed through the non-disableable redactor")
	} else {
		addManifest("config.redacted.toml", "not collected: config path not provided")
	}

	// 3. SLO snapshot (the attachable sampler data).
	if len(o.SLOSnapshot) > 0 {
		name := o.SLOSnapshotName
		if name == "" {
			name = "slo-snapshot.json"
		}
		if !strings.HasSuffix(name, ".json") {
			return "", New(ClassConfig, "observe: slo snapshot entry must end in .json")
		}
		if err := writeZipEntry(zw, name, o.SLOSnapshot); err != nil {
			return "", err
		}
		addManifest(name, "sampler report attached")
	} else {
		addManifest("slo-snapshot.json", "not attached: no sampler data supplied")
	}

	// 4. Version info (caller-supplied buildinfo block).
	ver, _ := json.MarshalIndent(o.Version, "", "  ")
	if o.Version == nil {
		ver = []byte("{\"note\": \"version block not supplied by caller\"}")
	}
	if err := writeZipEntry(zw, "version.json", ver); err != nil {
		return "", err
	}
	addManifest("version.json", "build/version block")

	// 5. DEFERRED register (caller-supplied; absence is recorded honestly).
	var deferred any
	detail := "register embedded"
	if len(o.Deferred) == 0 {
		deferred = map[string]string{
			"registered": "false",
			"reason":     "deferred register not embedded in this build (ticket 45)",
		}
		detail = "register not embedded (ticket 45)"
	} else {
		deferred = o.Deferred
	}
	def, _ := json.MarshalIndent(deferred, "", "  ")
	if err := writeZipEntry(zw, "deferred.json", def); err != nil {
		return "", err
	}
	addManifest("deferred.json", detail)

	manifest.WriteString("\n  ]\n}\n")
	if err := writeZipEntry(zw, "manifest.json", []byte(manifest.String())); err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", Wrap(ClassResource, err, "observe: close bundle zip")
	}
	if err := f.Sync(); err != nil {
		return "", Wrap(ClassResource, err, "observe: sync bundle zip")
	}
	if err := f.Close(); err != nil {
		return "", Wrap(ClassResource, err, "observe: close bundle file")
	}
	complete = true
	return o.OutPath, nil
}

func writeZipEntry(zw *zip.Writer, name string, data []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return Wrap(ClassResource, err, "observe: zip entry "+name)
	}
	if _, err := w.Write(data); err != nil {
		return Wrap(ClassResource, err, "observe: zip write "+name)
	}
	return nil
}
