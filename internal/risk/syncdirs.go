package risk

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Sync-directory detection (C25 / P12, SPEC-06 §5, D33/F4).
//
// Writing into a cloud-sync folder (OneDrive / Dropbox / Nutstore(坚果云) /
// Google Drive) is automatic third-party upload — the most covert of the six
// exfil channels. P12 is explicit that path-string matching ALONE is not
// sufficient: clients can be relocated, so we probe, in order,
//
//	registry (Windows, per-client account keys)  -> configured locations
//	client config files                          -> configured locations
//	known default locations (exist on disk)      -> known locations
//	fixture file (SyncFixturePath / injected)    -> test + user escape hatch
//
// and if NOTHING can be confirmed (undetectable) we fall back to the safe
// default demanded by the ticket: [fs] allowed_dirs entries under the user
// profile are treated as sync-suspect, i.e. fs.write there is an exfil
// channel whenever content matches. The overall conclusion (including the
// residual gaps) is recorded in docs/PRECHECK.md §P12.
//
// Membership comparison always runs on C26-resolved canonical paths
// (Resolve in this package): junction/8.3/UNC spellings of a sync root
// collapse onto the same root. An UNRESOLVABLE path fail-closes to
// "sync-suspect" (unverifiable = stricter side).

// SyncRoot names one detected sync location and how it was found.
type SyncRoot struct {
	// Provider is the client ("OneDrive", "Dropbox", "Nutstore",
	// "GoogleDrive"); suspect fallback uses "sync-suspect".
	Provider string `json:"provider"`
	// Path is the root as configured/discovered (raw spelling).
	Path string `json:"path"`
	// Source is the probe that produced it: "registry" | "config" |
	// "default" | "fixture" | "options" | "suspect-fallback".
	Source string `json:"source"`
}

// syncSet is the immutable, canonicalized view of the sync roots.
type syncSet struct {
	mu       sync.RWMutex
	roots    []syncRootC
	complete bool // at least one real location confirmed (no suspect fallback)
	home     string
	excs     []string // [fs] reparse_point_exceptions for C26 Resolve
}

type syncRootC struct {
	root      SyncRoot
	canon     string // normPath(Resolve(Path))
	canonical bool   // Resolve succeeded (else only lexical — fail toward suspect)
}

// detectSyncSet runs the P12 probes against the environment described by o.
func detectSyncSet(o ProvOptions) *syncSet {
	s := &syncSet{home: normPath(orDefault(o.HomeDir, userHomeDir())), excs: o.ReparseExceptions}
	env := probeEnv{
		Home:         orDefault(o.HomeDir, userHomeDir()),
		LocalAppData: orDefault(o.LocalAppData, envOr("LOCALAPPDATA", "")),
		AppData:      orDefault(o.AppData, envOr("APPDATA", "")),
	}

	var found []SyncRoot
	found = append(found, registryProbe(env)...) // Windows HKCU keys; empty elsewhere
	found = append(found, clientConfigProbe(env)...)
	found = append(found, defaultLocationProbe(env)...)
	if o.SyncFixturePath != "" {
		found = append(found, fixtureRoots(o.SyncFixturePath)...)
	}
	found = append(found, o.SyncRoots...)

	for _, r := range found {
		if strings.TrimSpace(r.Path) == "" {
			continue
		}
		s.add(r)
	}
	s.complete = len(s.roots) > 0
	if !s.complete {
		logf("risk/C25: sync-dir detection found no location (P12 undetectable) — [fs] allowed_dirs under the user profile fall back to sync-suspect")
	}
	return s
}

// add canonicalizes one raw root through C26 and appends it.
func (s *syncSet) add(r SyncRoot) {
	res, err := Resolve(r.Path, s.excs)
	c := syncRootC{root: r}
	if err == nil && res.Canonical != "" {
		c.canon, c.canonical = normPath(res.Canonical), true
	} else {
		// Root itself unresolvable (e.g. not yet created): keep the lexical
		// form for comparison but remember it is not verified.
		c.canon, c.canonical = normPath(expandInput(r.Path)), false
	}
	for _, e := range s.roots {
		if e.canon == c.canon {
			return
		}
	}
	s.roots = append(s.roots, c)
}

// SyncStatus is the verdict for one target path.
type SyncStatus struct {
	Sync bool
	Root SyncRoot // the matched (or fail-closed suspect) root
	Why  string
}

// IsSyncPath is the exported P12 check used by the fs.write channel (and the
// panel Security page, ticket 40). rawPath is canonicalized through C26.
// Fail-closed: an unresolvable path, or any under-profile allowed-dir entry
// while detection is undetectable, counts as sync.
func (p *Provenance) IsSyncPath(rawPath string) SyncStatus { return p.sync.match(rawPath) }

func (s *syncSet) match(rawPath string) SyncStatus {
	if s == nil {
		return SyncStatus{Sync: true, Root: SyncRoot{Provider: "sync-suspect", Source: "no-detector"},
			Why: "no sync detector configured; fail-closed: write target unverifiable"}
	}
	res, err := Resolve(rawPath, s.excs)
	if err != nil {
		return SyncStatus{Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Source: "suspect-fallback"},
			Why:  "C26 could not resolve the write target; fail-closed as sync-suspect: " + err.Error()}
	}
	if res.Canonical == "" {
		return SyncStatus{Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Source: "suspect-fallback"},
			Why:  "empty/unparseable write target; fail-closed as sync-suspect"}
	}
	cand := normPath(res.Canonical)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.roots {
		if isUnder(cand, e.canon) {
			return SyncStatus{Sync: true, Root: e.root, Why: "write target is under a detected sync root"}
		}
	}
	if !s.complete && s.home != "" && strings.HasPrefix(cand, s.home) {
		return SyncStatus{Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Path: s.home, Source: "suspect-fallback"},
			Why:  "sync locations undetectable (P12); under-profile path treated as sync-suspect (safe default)"}
	}
	return SyncStatus{Sync: false, Why: "write target is not under any sync root"}
}

// Roots returns the detected roots (panel/debug view). The suspect fallback
// is reported via Complete()==false, not as a root.
func (p *Provenance) SyncRoots() []SyncRoot {
	if p.sync == nil {
		return nil
	}
	p.sync.mu.RLock()
	defer p.sync.mu.RUnlock()
	out := make([]SyncRoot, 0, len(p.sync.roots))
	for _, e := range p.sync.roots {
		out = append(out, e.root)
	}
	return out
}

// Complete reports whether at least one concrete sync location was detected.
// false = undetectable -> under-profile writes are sync-suspect.
func (p *Provenance) SyncDetectionComplete() bool {
	return p.sync != nil && p.sync.complete
}

// --- probes -----------------------------------------------------------------

// probeEnv locates client configs for a (possibly relocated) profile; tests
// point it at a fixture tree.
type probeEnv struct{ Home, LocalAppData, AppData string }

func defaultLocationProbe(e probeEnv) []SyncRoot {
	var out []SyncRoot
	for _, d := range []struct{ provider, dir string }{
		{"OneDrive", "OneDrive"},
		{"Dropbox", "Dropbox"},
		{"Nutstore", "Nutstore"},
		{"GoogleDrive", "Google Drive"},
	} {
		if e.Home == "" {
			continue
		}
		p := filepath.Join(e.Home, d.dir)
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			out = append(out, SyncRoot{Provider: d.provider, Path: p, Source: "default"})
		}
	}
	return out
}

// clientConfigProbe reads the on-disk client configurations that hold the
// configured root when the default location was moved.
func clientConfigProbe(e probeEnv) []SyncRoot {
	var out []SyncRoot
	// Dropbox: host.db — line 1 is a version dir, line 2 the base64 root.
	for _, base := range []string{e.LocalAppData, e.AppData} {
		if base == "" {
			continue
		}
		if root := dropboxHostDB(filepath.Join(base, "Dropbox", "host.db")); root != "" {
			out = append(out, SyncRoot{Provider: "Dropbox", Path: root, Source: "config"})
		}
	}
	// Google Drive: %LOCALAPPDATA%\Google\Drive\preference (JSON "root";
	// newer clients store "root_encrypted" -> not parseable, P12 gap: the
	// default-location probe + suspect fallback cover it).
	if e.LocalAppData != "" {
		if root := jsonKeyRoot(filepath.Join(e.LocalAppData, "Google", "Drive", "preference"), "root"); root != "" {
			out = append(out, SyncRoot{Provider: "GoogleDrive", Path: root, Source: "config"})
		}
	}
	return out
}

func dropboxHostDB(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	for _, ln := range lines[1:] {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if dec, err := base64.StdEncoding.DecodeString(ln); err == nil && len(dec) > 0 {
			return string(dec)
		}
	}
	return ""
}

func jsonKeyRoot(path, key string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// fixtureRoots loads a JSON array of {provider,path} — the acceptance-criteria
// fixture fallback (and the relocated-client user escape hatch fed by config).
func fixtureRoots(path string) []SyncRoot {
	b, err := os.ReadFile(path)
	if err != nil {
		logf("risk/C25: sync fixture %s unreadable (fail-closed path: fewer confirmed roots): %v", path, err)
		return nil
	}
	var rs []SyncRoot
	if err := json.Unmarshal(b, &rs); err != nil {
		logf("risk/C25: sync fixture %s malformed: %v", path, err)
		return nil
	}
	for i := range rs {
		rs[i].Source = "fixture"
	}
	return rs
}

func orDefault(v, def string) string {
	if v != "" {
		return v
	}
	return def
}
