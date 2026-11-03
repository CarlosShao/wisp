package observe

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	// TranscriptDir is the ASR transcript directory. Nothing under it is read
	// unless its name is in OptIn.Transcripts: transcripts are default-excluded
	// and opted in item by item (ticket 45, SPEC-14 §14.4).
	TranscriptDir string
	// SnapshotDir is the watchdog snapshot directory, opted in the same way.
	SnapshotDir string
	// OptIn is the user's explicit per-item consent, taken after seeing a
	// PreviewDiagnosticsBundle. It is the zero value for a caller that never
	// previews, which is the default-excluded bundle.
	OptIn BundleOptIn
	// DiskStatus is the disk-full guard's reading to attach (nil records the
	// guard as not installed - the manifest never claims "no guard, no problem").
	DiskStatus *DiskStatus
}

// BundleOptIn is the explicit consent §14.4 requires before anything beyond the
// default set leaves the machine. Names are BASE names inside the matching
// directory: a name holding a separator or ".." is refused outright rather than
// resolved, because the consent the user gave is for the item they saw, not for
// a path that reaches elsewhere.
type BundleOptIn struct {
	// Transcripts names ASR transcript files to include, one item at a time.
	Transcripts []string
	// Snapshots names watchdog snapshot files to include.
	Snapshots []string
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

	// 6. The disk-full guard's reading (ticket 45, D42#5). Attached whether or
	// not the guard is tripped: a bundle that shows only problems is a bundle
	// that cannot tell "nothing happened" from "nobody was looking".
	if o.DiskStatus != nil {
		ds, derr := json.MarshalIndent(o.DiskStatus, "", "  ")
		if derr != nil {
			return "", Wrap(ClassInternal, derr, "observe: marshal disk status")
		}
		if err := writeZipEntry(zw, "disk-guard.json", ds); err != nil {
			return "", err
		}
		addManifest("disk-guard.json", fmt.Sprintf("tripped=%t, trips=%d", o.DiskStatus.Tripped, o.DiskStatus.Trips))
	} else {
		addManifest("disk-guard.json", "not attached: no disk-full guard reading supplied")
	}

	// 7. Opt-in material (SPEC-14 §14.4). Nothing here is read from disk unless
	// the user named it after seeing the preview, and audio is never read at
	// all - it is not collected, so it cannot leak.
	if len(o.OptIn.Transcripts) == 0 {
		addManifest("transcripts/", "not collected: default-excluded, needs an explicit per-item opt-in")
	} else {
		if o.TranscriptDir == "" {
			return "", New(ClassConfig, "observe: transcript opt-in supplied without a transcript dir")
		}
		red := Redactor{RedactPaths: o.RedactPaths}
		for _, name := range o.OptIn.Transcripts {
			if err := checkBundleItemName(name); err != nil {
				return "", err
			}
			data, rerr := os.ReadFile(filepath.Join(o.TranscriptDir, name))
			if rerr != nil {
				return "", Wrap(ClassResource, rerr, "observe: read opted-in transcript "+name)
			}
			// Transcripts are the user's own words, and words are where a key
			// shape hides: the same non-disableable redactor, no second format.
			if err := writeZipEntry(zw, "transcripts/"+name, []byte(red.String(string(data)))); err != nil {
				return "", err
			}
			addManifest("transcripts/"+name, fmt.Sprintf("opted in by the user, %d byte(s), re-passed through the redactor", len(data)))
		}
	}
	if len(o.OptIn.Snapshots) == 0 {
		addManifest("watchdog/", "not collected: default-excluded, needs an explicit per-item opt-in")
	} else {
		if o.SnapshotDir == "" {
			return "", New(ClassConfig, "observe: watchdog snapshot opt-in supplied without a snapshot dir")
		}
		for _, name := range o.OptIn.Snapshots {
			if err := checkBundleItemName(name); err != nil {
				return "", err
			}
			data, rerr := os.ReadFile(filepath.Join(o.SnapshotDir, name))
			if rerr != nil {
				return "", Wrap(ClassResource, rerr, "observe: read opted-in snapshot "+name)
			}
			if err := writeZipEntry(zw, "watchdog/"+name, data); err != nil {
				return "", err
			}
			addManifest("watchdog/"+name, fmt.Sprintf("opted in by the user, %d byte(s)", len(data)))
		}
	}
	addManifest("audio/", "never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)")

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

// ---------------------------------------------------------------------------
// The preview (ticket 45, SPEC-14 §14.4: "user explicitly opts in per-category
// with a preview of what's included BEFORE export").
//
// WHY IT IS A SIBLING OF THE BUILDER AND NOT A CONSTANT. "Bundle contents
// exactly as previewed" is a claim about two code paths agreeing, and it can
// only be tested by walking the same directories the builder walks with the same
// caps and the same name rules. A preview that enumerates a fixed category list
// would pass its own test and then lie about the zip. So PreviewDiagnosticsBundle
// reads the same inputs BundleOptions carries, applies the same skips, and hands
// back the entry names BuildDiagnosticsBundle would write - entry names, not a
// prose summary, because the test compares them against the real archive.
//
// The categories a bundle CANNOT contain are in the preview too, marked as such
// (audio): a review that only lists what is there cannot show that the thing you
// were worried about is genuinely not there.
// ---------------------------------------------------------------------------

// PreviewItem is one file the bundle would carry, with the head of its content
// so a person can recognise it.
type PreviewItem struct {
	Name  string `json:"name"`
	Bytes int64  `json:"bytes"`
	// Head is the per-item preview: the first line, redacted through the same
	// non-disableable redactor the builder uses, capped.
	Head     string `json:"head,omitempty"`
	Included bool   `json:"included"`
	Reason   string `json:"reason,omitempty"`
}

// PreviewCategory is one category of the bundle: its policy (default in or
// default out), what the current options actually do, and its items.
type PreviewCategory struct {
	Category string `json:"category"`
	// DefaultIncluded is the policy: transcripts and watchdog snapshots start
	// out, audio is never in.
	DefaultIncluded bool          `json:"default_included"`
	Collected       bool          `json:"collected"`
	Detail          string        `json:"detail"`
	Items           []PreviewItem `json:"items,omitempty"`
}

// BundlePreview is the review-before-send reading of a BundleOptions.
type BundlePreview struct {
	CreatedAt  string            `json:"created_at"`
	OutPath    string            `json:"out_path"`
	Categories []PreviewCategory `json:"categories"`
	// Entries are the exact zip entry names BuildDiagnosticsBundle(o) would
	// write, in the order it writes them, manifest.json last.
	Entries    []string `json:"entries"`
	TotalBytes int64    `json:"total_bytes"`
}

// PreviewDiagnosticsBundle reads what a bundle WOULD contain without writing
// anything. It fails on the same inputs the builder fails on (an unreadable log
// dir, a missing config), so a preview that succeeded is a preview of a bundle
// that could actually be built.
func PreviewDiagnosticsBundle(o BundleOptions) (BundlePreview, error) {
	if o.OutPath == "" {
		return BundlePreview{}, New(ClassConfig, "observe: bundle preview output path is empty")
	}
	if o.MaxLogBytes <= 0 {
		o.MaxLogBytes = 10 << 20
	}
	if o.Now.IsZero() {
		o.Now = NowWallUTC()
	}
	red := Redactor{RedactPaths: o.RedactPaths}
	p := BundlePreview{CreatedAt: WallTimestampUTC(o.Now), OutPath: o.OutPath}

	// 1. Logs.
	logCat := PreviewCategory{
		Category: "logs", DefaultIncluded: true,
		Detail: "rolling JSONL logs; redacted at rest and re-passed through the redactor",
	}
	if o.LogDir == "" {
		logCat.Detail = "not collected: log dir not provided"
	} else {
		entries, err := os.ReadDir(o.LogDir)
		if err != nil {
			return p, Wrap(ClassResource, err, "observe: read log dir for bundle preview")
		}
		copied, skipped := 0, 0
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasPrefix(name, logFilePrefix) || !strings.HasSuffix(name, logFileExt) {
				continue
			}
			full := filepath.Join(o.LogDir, name)
			head, size, herr := previewItemHead(full, red)
			if herr != nil {
				return p, Wrap(ClassResource, herr, "observe: read log file for bundle preview")
			}
			if size > o.MaxLogBytes {
				// The builder skips oversized files; a preview that listed them
				// would promise an entry the zip will not hold.
				skipped++
				logCat.Items = append(logCat.Items, PreviewItem{
					Name: "logs/" + name, Bytes: size, Included: false,
					Reason: fmt.Sprintf("skipped: over the %d byte per-file cap", o.MaxLogBytes),
				})
				continue
			}
			copied++
			p.addEntry("logs/" + name)
			logCat.Items = append(logCat.Items, PreviewItem{
				Name: "logs/" + name, Bytes: size, Head: head, Included: true,
			})
		}
		logCat.Collected = copied > 0
		logCat.Detail = fmt.Sprintf("%d file(s) would be copied, %d skipped (oversized); %s",
			copied, skipped, logCat.Detail)
	}
	p.Categories = append(p.Categories, logCat)

	// 2. Config snapshot.
	cfgCat := PreviewCategory{
		Category: "config", DefaultIncluded: true,
		Detail: "config.toml re-passed through the non-disableable redactor; a surviving key shape aborts the build",
	}
	if o.ConfigPath == "" {
		cfgCat.Detail = "not collected: config path not provided"
	} else {
		head, size, herr := previewItemHead(o.ConfigPath, red)
		if herr != nil {
			return p, Wrap(ClassResource, herr, "observe: read config for bundle preview")
		}
		cfgCat.Collected = true
		cfgCat.Items = append(cfgCat.Items, PreviewItem{
			Name: "config.redacted.toml", Bytes: size, Head: head, Included: true,
		})
		p.addEntry("config.redacted.toml")
	}
	p.Categories = append(p.Categories, cfgCat)

	// 3. SLO snapshot.
	sloName := o.SLOSnapshotName
	if sloName == "" {
		sloName = "slo-snapshot.json"
	}
	sloCat := PreviewCategory{
		Category: "slo", DefaultIncluded: true,
		Detail: "the sampler report the caller supplies",
	}
	if len(o.SLOSnapshot) > 0 {
		if !strings.HasSuffix(sloName, ".json") {
			return p, New(ClassConfig, "observe: slo snapshot entry must end in .json")
		}
		sloCat.Collected = true
		sloCat.Items = append(sloCat.Items, PreviewItem{
			Name: sloName, Bytes: int64(len(o.SLOSnapshot)), Included: true,
		})
		p.addEntry(sloName)
	} else {
		sloCat.Detail = "not attached: no sampler data supplied"
	}
	p.Categories = append(p.Categories, sloCat)

	// 4 + 5. Version and the DEFERRED register (always written by the builder).
	verCat := PreviewCategory{
		Category: "version", DefaultIncluded: true, Collected: true,
		Detail: "build/version block supplied by the caller",
	}
	verBytes := int64(len("{\"note\": \"version block not supplied by caller\"}"))
	if o.Version != nil {
		ver, _ := json.MarshalIndent(o.Version, "", "  ")
		verBytes = int64(len(ver))
		verCat.Detail = fmt.Sprintf("build/version block: %d key(s)", len(o.Version))
	}
	verCat.Items = append(verCat.Items, PreviewItem{Name: "version.json", Bytes: verBytes, Included: true})
	p.addEntry("version.json")
	p.Categories = append(p.Categories, verCat)

	defCat := PreviewCategory{
		Category: "deferred", DefaultIncluded: true, Collected: true,
		Detail: "the DEFERRED register, or an honest record that it is not embedded",
	}
	defCat.Items = append(defCat.Items, PreviewItem{
		Name: "deferred.json", Bytes: int64(len(o.Deferred)) * 64, Included: true,
		Reason: fmt.Sprintf("%d register entry/ies supplied", len(o.Deferred)),
	})
	p.addEntry("deferred.json")
	p.Categories = append(p.Categories, defCat)

	// 6. The disk-full guard reading.
	diskCat := PreviewCategory{
		Category: "disk-guard", DefaultIncluded: false,
		Detail: "the disk-full guard's live reading (ticket 45); attached when the caller supplies it",
	}
	if o.DiskStatus != nil {
		diskCat.Collected = true
		diskCat.Items = append(diskCat.Items, PreviewItem{
			Name: "disk-guard.json", Included: true,
			Reason: fmt.Sprintf("tripped=%t, trips=%d", o.DiskStatus.Tripped, o.DiskStatus.Trips),
		})
		p.addEntry("disk-guard.json")
	} else {
		diskCat.Detail = "not attached: no disk-full guard reading supplied"
	}
	p.Categories = append(p.Categories, diskCat)

	// 7. Transcripts: default-excluded, per-item consent, per-item preview.
	tCat := PreviewCategory{
		Category: "transcripts", DefaultIncluded: false,
		Detail: "ASR transcripts - excluded by default; each item must be named after it was previewed",
	}
	if len(o.OptIn.Transcripts) == 0 {
		tCat.Detail += " (none opted in: nothing under the transcript dir is read)"
		p.Categories = append(p.Categories, tCat)
	} else {
		if o.TranscriptDir == "" {
			return p, New(ClassConfig, "observe: transcript opt-in supplied without a transcript dir")
		}
		for _, name := range o.OptIn.Transcripts {
			if err := checkBundleItemName(name); err != nil {
				return p, err
			}
			head, size, herr := previewItemHead(filepath.Join(o.TranscriptDir, name), red)
			if herr != nil {
				return p, Wrap(ClassResource, herr, "observe: read opted-in transcript for bundle preview")
			}
			tCat.Collected = true
			tCat.Items = append(tCat.Items, PreviewItem{
				Name: "transcripts/" + name, Bytes: size, Head: head, Included: true,
			})
			p.addEntry("transcripts/" + name)
		}
		p.Categories = append(p.Categories, tCat)
	}

	// 8. Watchdog snapshots: default-excluded, same consent shape.
	wCat := PreviewCategory{
		Category: "watchdog", DefaultIncluded: false, Collected: len(o.OptIn.Snapshots) > 0,
		Detail: "watchdog snapshots - excluded by default; each item must be named after it was previewed",
	}
	for _, name := range o.OptIn.Snapshots {
		if err := checkBundleItemName(name); err != nil {
			return p, err
		}
		if o.SnapshotDir == "" {
			return p, New(ClassConfig, "observe: watchdog snapshot opt-in supplied without a snapshot dir")
		}
		_, size, herr := previewItemHead(filepath.Join(o.SnapshotDir, name), red)
		if herr != nil {
			return p, Wrap(ClassResource, herr, "observe: read opted-in snapshot for bundle preview")
		}
		wCat.Items = append(wCat.Items, PreviewItem{
			Name: "watchdog/" + name, Bytes: size, Included: true,
		})
		p.addEntry("watchdog/" + name)
	}
	p.Categories = append(p.Categories, wCat)

	// 9. Audio: stated out loud, because "absent" is the claim under review.
	p.Categories = append(p.Categories, PreviewCategory{
		Category: "audio", DefaultIncluded: false, Collected: false,
		Detail: "never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)",
	})

	p.addEntry("manifest.json")
	for _, c := range p.Categories {
		for _, it := range c.Items {
			if it.Included {
				p.TotalBytes += it.Bytes
			}
		}
	}
	return p, nil
}

// addEntry records one entry name in the order the builder writes it.
func (p *BundlePreview) addEntry(name string) {
	p.Entries = append(p.Entries, name)
}

// checkBundleItemName is the consent's shape rule: a name the user could have
// seen in a preview is a base name, so anything that could reach outside the
// directory is refused instead of resolved.
func checkBundleItemName(name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) || strings.HasPrefix(name, ".") {
		return New(ClassConfig, fmt.Sprintf("observe: bundle opt-in item %q is not a plain base name", name))
	}
	return nil
}

// previewItemHead reads the first line of a file for the per-item preview, up to
// 4KiB, and passes it through the caller's redactor: the preview never shows a
// shape the bundle itself would refuse to ship.
func previewItemHead(path string, red Redactor) (head string, size int64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	info, serr := f.Stat()
	if serr != nil {
		return "", 0, serr
	}
	buf := make([]byte, 4096)
	n, rerr := io.ReadFull(f, buf)
	if rerr != nil && !errors.Is(rerr, io.ErrUnexpectedEOF) && !errors.Is(rerr, io.EOF) {
		return "", info.Size(), rerr
	}
	line := string(buf[:n])
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	if len(line) > 200 {
		line = line[:200]
	}
	return red.String(line), info.Size(), nil
}
