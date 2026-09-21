package risk

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
//	client environment variables (OneDrive publishes its own root) -> "env"
//	registry (Windows, per-client account keys)                     -> "registry"
//	client config files                                             -> "config"
//	known default locations (exist on disk)                         -> "default"
//	fixture file (SyncFixturePath / injected)                       -> "fixture"
//
// and if NOTHING *confirmed* (env/registry/config grade, C26-canonical) can be
// established we fall back to the safe default demanded by the ticket: paths
// under the user profile are treated as sync-suspect, i.e. fs.write there is an
// exfil channel whenever content matches. Weaker evidence (a directory that
// merely exists under the profile, a fixture entry) adds roots but never
// switches that net off — otherwise one decoy `~\OneDrive` folder would
// disarm it (adversarial report M-3). The overall conclusion (including the
// residual gaps) is recorded in docs/PRECHECK.md §P12.
//
// Membership comparison always runs on C26-resolved canonical paths
// (Resolve in this package): junction/8.3/UNC spellings of a sync root
// collapse onto the same root. Because fs.write targets usually DO NOT exist
// yet (where C26's handle query cannot run), an unresolved candidate is
// re-anchored on its nearest existing ancestor and only then compared;
// something the OS still cannot verify fail-closes to "sync-suspect"
// (unverifiable = stricter side). A spelling that carries a `..` component is
// never classified at all — cleaning erases the segments the reparse audit
// depends on, so it fails closed the same way (N-10).

// SyncRoot names one detected sync location and how it was found.
type SyncRoot struct {
	// Provider is the client ("OneDrive", "Dropbox", "Nutstore",
	// "GoogleDrive"); suspect fallback uses "sync-suspect".
	Provider string `json:"provider"`
	// Path is the root as configured/discovered (raw spelling).
	Path string `json:"path"`
	// Source is the probe that produced it: "env" | "registry" | "config" |
	// "default" | "fixture" | "options" | "suspect-fallback".
	Source string `json:"source"`
}

// gradeConfirmed reports whether a probe result is strong enough to switch the
// under-profile sync-suspect net OFF. Only client-published configuration
// counts (env var the sync client sets, registry record, client config file);
// a directory that merely exists at a default location, a fixture line or an
// [fs] option is not proof that syncing happens there — P12 grades inaccurate
// detection as "this channel is open", so the weaker grades only ADD roots.
func gradeConfirmed(src string) bool {
	switch src {
	case "env", "registry", "config":
		return true
	}
	return false
}

// syncSet is the immutable, canonicalized view of the sync roots.
type syncSet struct {
	mu       sync.RWMutex
	roots    []syncRootC
	complete bool // >=1 C26-canonical root from a confirmed-grade probe
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
		Home:               orDefault(o.HomeDir, userHomeDir()),
		LocalAppData:       orDefault(o.LocalAppData, envOr("LOCALAPPDATA", "")),
		AppData:            orDefault(o.AppData, envOr("APPDATA", "")),
		OneDrive:           envOr("OneDrive", ""),
		OneDriveConsumer:   envOr("OneDriveConsumer", ""),
		OneDrivePersonal:   envOr("OneDrivePersonal", ""),
		OneDriveCommercial: envOr("OneDriveCommercial", ""),
		OneDrivePublic:     envOr("OneDrivePublic", ""),
	}

	var found []SyncRoot
	found = append(found, envConfiguredRoots(env)...) // client-published, first
	found = append(found, registryProbe(env)...)      // Windows HKCU keys; empty elsewhere
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
	s.finalize()
	return s
}

// finalize derives the fallback switch from the collected roots: only a root
// that is BOTH confirmed-grade and C26-canonical disarms the sync-suspect net
// (adversarial report M-3 + MINOR N-2: `canonical` used to be set, never read).
func (s *syncSet) finalize() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.complete = false
	for _, r := range s.roots {
		if r.canonical && gradeConfirmed(r.root.Source) {
			s.complete = true
			break
		}
	}
	if !s.complete {
		logf("risk/C25: sync-dir detection has no confirmed (env/registry/config) location — paths under the user profile fall back to sync-suspect (P12 safe default)")
	}
}

// add canonicalizes one raw root through C26 and appends it.
func (s *syncSet) add(r SyncRoot) {
	res, err := Resolve(r.Path, s.excs)
	c := syncRootC{root: r}
	if err == nil && res.Resolved && res.Canonical != "" {
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

// errTargetUnverified marks a write target whose spelling C26 could not be
// anchored on anything the OS confirms; the caller fails closed.
var errTargetUnverified = errors.New("no existing ancestor of the target could be verified through C26")

// IsSyncPath is the exported P12 check used by the fs.write channel (and the
// panel Security page, ticket 40). rawPath is canonicalized through C26,
// re-anchored on its nearest existing ancestor when the target does not exist
// yet, and only then compared.
// Fail-closed: an unresolvable path, or any under-profile path while no
// confirmed sync location exists, counts as sync.
func (p *Provenance) IsSyncPath(rawPath string) SyncStatus { return p.sync.match(rawPath) }

// resolveTarget produces the canonical spelling one write candidate must be
// compared with. C26's handle query (the step that expands 8.3 short names and
// strips the \\?\ prefix) only succeeds on paths that already exist, and
// fs.write targets usually do not — taking the lexical fallback at face value
// let `\\?\C:\Users\x\OneDrive\...` and `ONEDRI~1\...` spellings of a sync root
// read as "not a sync dir" while the bytes landed inside it (adversarial
// report B-1). Fix: resolve the nearest EXISTING ancestor through the very
// same Resolve entry point and re-attach the lexical remainder.
func (s *syncSet) resolveTarget(raw string) (string, error) {
	res, err := Resolve(raw, s.excs)
	if err != nil {
		return "", err // non-exempted reparse traversal: fail closed upstream
	}
	if res.Canonical == "" {
		return "", errTargetUnverified
	}
	if res.Resolved {
		return normPath(res.Canonical), nil
	}
	anc, rest := deepestExistingAncestor(res.Canonical)
	if anc == "" {
		return "", errTargetUnverified // nothing of this chain exists on disk
	}
	ares, err := Resolve(anc, s.excs)
	if err != nil {
		return "", err // the existing part of the chain traverses a reparse point
	}
	if ares.Canonical == "" {
		return "", errTargetUnverified
	}
	if !ares.Resolved {
		if pathHandleVerify {
			// Windows: an existing directory C26 still cannot verify (ACL,
			// offline volume) -> the spelling is unproven.
			return "", errTargetUnverified
		}
		// Platforms without handle-based resolution (see the DEFERRED marker
		// in pathresolver_other.go): the lexical form is all C26 can offer
		// there; the ancestor chain at least is known to exist.
		//
		// NOTE(adversarial re-verification N-9): on POSIX this branch is
		// currently UNREACHABLE — deepestExistingAncestor Lstats `\`-joined
		// components, which never exist there, so the walk returns anc=="" and
		// the caller lands in errTargetUnverified: every write is judged
		// sync-suspect (the strict side, not a hole, but it would flag every
		// plain local write once the engine is wired on macOS). Closing that —
		// native separators on the walk, plus a real CloudStorage probe — is
		// ticket 55's AC line, not this file's; compiling for darwin does NOT
		// mean this path behaves.
	}
	return normPath(ares.Canonical + `\` + strings.Join(rest, `\`)), nil
}

// deepestExistingAncestor splits an absolute, lexically cleaned path into the
// longest prefix whose components all exist on disk and the remaining (not yet
// existing) components. A plain file stops the walk because nothing can live
// under it.
//
// Symlinks/junctions are stepped over, and NOT for the reason one might
// guess (adversarial re-verification N-9): every prefix here is Lstat-ed with a
// trailing separator, and an Lstat with a trailing `\` follows a mount point,
// so a junction stats as its TARGET directory rather than as a reparse point
// (Go's Lstat on the same junction without the separator reports
// isdir=false/symlink=false plus the reparse attribute bit). This walk
// therefore never detects a reparse point at all; detection is the caller's
// `Resolve(raw)` over the whole cleaned chain, which runs first and denies a
// non-exempted traversal before this function is ever reached. For an
// explicitly EXEMPTED junction stepping through is exactly what is wanted:
// Resolve(anc) then normalizes through the mount point onto the real root.
func deepestExistingAncestor(canon string) (string, []string) {
	root, comps := splitPathComponents(canon)
	if root == "" || len(comps) == 0 {
		return "", nil
	}
	base, n := root, 0
	for i, c := range comps {
		next := base + c
		if !strings.HasSuffix(next, `\`) {
			next += `\`
		}
		st, err := os.Lstat(next)
		if err != nil {
			break // this component and nothing below it exists
		}
		if !st.IsDir() && st.Mode()&os.ModeSymlink == 0 {
			break // a plain file cannot be traversed
		}
		base, n = next, i+1
	}
	if n == len(comps) {
		return "", nil // full path exists yet was not handle-resolvable
	}
	return strings.TrimSuffix(base, `\`), comps[n:]
}

// hasFoldedDotDot reports whether the raw spelling carries a `..` component,
// i.e. a segment that lexical cleaning (C26's lexCanonical, and Windows' own
// pre-parse normalization) folds away together with whatever stood before it.
// Such a spelling cannot be classified from its cleaned form, so the caller
// fails closed on it (adversarial re-verification N-10); ordinary fs.write
// targets contain no `..`, so the false-positive surface is ~zero.
func hasFoldedDotDot(raw string) bool {
	for _, c := range strings.Split(strings.ReplaceAll(raw, "/", `\`), `\`) {
		if c == ".." {
			return true
		}
	}
	return false
}

// splitPathComponents separates the volume/UNC root of an absolute path
// (C:\, \\?\C:\, \\server\share\, \\?\UNC\server\share\, or the POSIX root
// mapped onto `\`) from its components. An empty root means "unclassifiable".
func splitPathComponents(p string) (string, []string) {
	u := strings.ReplaceAll(p, "/", `\`)
	var root, rest string
	switch {
	case strings.HasPrefix(u, `\\?\UNC\`):
		body := u[len(`\\?\UNC\`):]
		server, share, ok := twoComponents(body)
		if !ok {
			return "", nil
		}
		root, rest = `\\?\UNC\`+server+`\`+share+`\`, body[len(server)+1+len(share)+1:]
	case strings.HasPrefix(u, `\\?\`):
		if len(u) < 6 || u[5] != ':' {
			return "", nil // \\?\PhysicalDrive0&... and friends: not classifiable
		}
		root, rest = u[:len(`\\?\C:`)]+`\`, u[6:]
	case strings.HasPrefix(u, `\\`):
		body := u[2:]
		server, share, ok := twoComponents(body)
		if !ok {
			return "", nil
		}
		root, rest = `\\`+server+`\`+share+`\`, body[len(server)+1+len(share)+1:]
	case len(u) >= 2 && u[1] == ':':
		root, rest = u[:2]+`\`, u[2:]
	case strings.HasPrefix(u, `\`):
		root, rest = `\`, u[1:]
	default:
		return "", nil
	}
	var comps []string
	for _, c := range strings.Split(rest, `\`) {
		if c != "" {
			comps = append(comps, c)
		}
	}
	return root, comps
}

func twoComponents(body string) (string, string, bool) {
	i := strings.Index(body, `\`)
	if i <= 0 || i == len(body)-1 {
		return "", "", false
	}
	server := body[:i]
	rest := body[i+1:]
	j := strings.Index(rest, `\`)
	if j < 0 {
		return server, rest, true
	}
	if j == 0 {
		return "", "", false
	}
	return server, rest[:j], true
}

func (s *syncSet) match(rawPath string) SyncStatus {
	if s == nil {
		return SyncStatus{
			Sync: true, Root: SyncRoot{Provider: "sync-suspect", Source: "no-detector"},
			Why: "no sync detector configured; fail-closed: write target unverifiable",
		}
	}
	if strings.TrimSpace(rawPath) == "" {
		return SyncStatus{
			Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Source: "suspect-fallback"},
			Why:  "empty/unparseable write target; fail-closed as sync-suspect",
		}
	}
	if hasFoldedDotDot(rawPath) {
		// N-10 hardening (adversarial re-verification R2/R5): a `..` component
		// is removed by lexical cleaning BEFORE anything is audited, so the
		// string the reparse-point check sees is not the string that was
		// classified. The landing probes showed Windows folds `.`/`..` the same
		// way before it opens the file, so this is safe *today* — but that is
		// OS semantics, not a check in here, so make it an invariant instead.
		return SyncStatus{
			Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Source: "suspect-fallback"},
			Why:  "write target folds a '..' component: cleaning erases the segments the reparse check audits; fail-closed as sync-suspect (N-10)",
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	cand, err := s.resolveTarget(rawPath)
	if err != nil || cand == "" {
		why := "write target spelling could not be anchored on a C26-verifiable ancestor (8.3/\\\\?\\/missing chain); fail-closed as sync-suspect"
		if errors.Is(err, ErrReparseDenied) {
			why = "write target traverses a non-exempted reparse point; fail-closed as sync-suspect"
		}
		if err != nil && !errors.Is(err, ErrReparseDenied) && !errors.Is(err, errTargetUnverified) {
			why += ": " + err.Error()
		}
		return SyncStatus{
			Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Source: "suspect-fallback"}, Why: why,
		}
	}
	for _, e := range s.roots {
		if isUnder(cand, e.canon) {
			return SyncStatus{Sync: true, Root: e.root, Why: "write target is under a detected sync root"}
		}
	}
	// Component-bounded under-profile test (isUnder, not a raw prefix: a
	// sibling `profileevil` must not be swept in — MINOR N-5).
	if !s.complete && s.home != "" && isUnder(cand, s.home) {
		return SyncStatus{
			Sync: true,
			Root: SyncRoot{Provider: "sync-suspect", Path: s.home, Source: "suspect-fallback"},
			Why:  "no confirmed sync location (P12); path under the user profile treated as sync-suspect (safe default)",
		}
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

// SyncDetectionComplete reports whether a CONFIRMED sync location (env var,
// registry or client-config grade, and C26-canonical) was established.
// false -> under-profile writes stay sync-suspect. Weak-grade roots (default
// location that merely exists, fixture, injected options) are listed by
// SyncRoots() but do not clear the net (adversarial report M-3).
func (p *Provenance) SyncDetectionComplete() bool {
	return p.sync != nil && p.sync.complete
}

// --- probes -----------------------------------------------------------------

// probeEnv locates client configs for a (possibly relocated) profile; tests
// point it at a fixture tree.
type probeEnv struct {
	Home, LocalAppData, AppData string
	// OneDrive publishes its configured root(s) as environment variables of
	// the user session; on machines where the HKCU Accounts key carries no
	// UserFolder this is the only configured-location signal in place
	// (adversarial report M-2).
	OneDrive           string
	OneDriveConsumer   string
	OneDrivePersonal   string
	OneDriveCommercial string
	OneDrivePublic     string
}

// registryValueSource abstracts the HKCU reads the P12 registry probe needs.
// Injecting it keeps the probe logic (which subkey names carry a UserFolder,
// which WOW64 view to trust, what counts as a root) table-testable without
// touching the live registry — the gap adversarial report M-2 flagged.
type registryValueSource interface {
	// subKeys lists the value names directly under key (nil on any error).
	subKeys(key string) []string
	// string reads one REG_SZ value ("" on any error).
	string(key, name string) string
}

// registryProbeFor is the portable half of the Windows registry probe: the
// key paths and value names below are the documented per-client account
// records that survive a relocated client.
func registryProbeFor(src registryValueSource) []SyncRoot {
	var out []SyncRoot
	if src == nil {
		return nil
	}
	// OneDrive: HKCU\Software\Microsoft\OneDrive\Accounts\<Personal|Business#>\UserFolder
	const accounts = `Software\Microsoft\OneDrive\Accounts`
	seen := map[string]bool{}
	for _, sk := range src.subKeys(accounts) {
		if !strings.EqualFold(sk, "Personal") && !strings.HasPrefix(sk, "Business") {
			continue // unknown nodes (e.g. "ConsumingAccounts") carry no UserFolder
		}
		v := src.string(accounts+`\`+sk, "UserFolder")
		if strings.TrimSpace(v) == "" || seen[strings.ToLower(v)] {
			continue
		}
		seen[strings.ToLower(v)] = true
		out = append(out, SyncRoot{Provider: "OneDrive", Path: v, Source: "registry"})
	}
	// Dropbox installer record: HKCU\Software\Dropbox\Dropbox "Path"
	if v := src.string(`Software\Dropbox\Dropbox`, "Path"); strings.TrimSpace(v) != "" {
		out = append(out, SyncRoot{Provider: "Dropbox", Path: v, Source: "registry"})
	}
	// Nutstore (坚果云): InstallPath is the client install dir, not the sync
	// root; sync roots live in the (undocumented, version-specific) client
	// config -> P12 residual, covered by the default-location probe and the
	// suspect fallback. See docs/PRECHECK.md.
	return out
}

// envConfiguredRoots reads the client-published environment variables. These
// are confirmed-grade evidence: the sync client itself writes them at sign-in
// with the root it actually syncs, which is exactly what P12 asks for when it
// rejects path-string guessing.
func envConfiguredRoots(e probeEnv) []SyncRoot {
	var out []SyncRoot
	for _, c := range []struct{ provider, dir string }{
		{"OneDrive", e.OneDrive},
		{"OneDrive", e.OneDriveConsumer},
		{"OneDrive", e.OneDrivePersonal},
		{"OneDrive", e.OneDriveCommercial},
		{"OneDrive", e.OneDrivePublic},
	} {
		if strings.TrimSpace(c.dir) == "" {
			continue
		}
		out = append(out, SyncRoot{Provider: c.provider, Path: c.dir, Source: "env"})
	}
	return out
}

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
