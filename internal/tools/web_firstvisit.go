package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// web_firstvisit.go is D30 layer 5 (ticket 22, SPEC-06 §6.5): a builtin web
// tool visiting a domain for the first time asks the user once, and the
// answer is per-domain rememberable. The prompt itself is a host seam —
// the production wiring hands an Ask that renders the native one-time card
// (same UI surface ticket 21/37 own); until one is wired this layer fails
// CLOSED: an unwired prompt channel is an unanswered question, and the
// answer to an unanswered first visit is "do not connect". That is the same
// posture NoGate takes for L1/L2 (gate.go), not a policy invention.
//
// Storage: the remembered set is a small JSON file next to the user's wisp
// data dir (WebDeps.DomainStorePath, composition-root-fed; D30⑤ says
// "stored in config/db", and neither the frozen [net] config schema nor the
// frozen SPEC-02 DDL has a home for it — adding a column or a locked key is
// a contract change this ticket must not make silently). The file is
// loaded lazily, kept in memory, and written back atomically (temp +
// rename, the D31 rule for anything that replaces a live file).

// FirstVisitAsk is the host's one-time prompt. allow=false refuses the
// visit; remember=true keeps the domain for later visits of any builtin
// web tool in this user profile.
type FirstVisitAsk func(ctx context.Context, tool, domain string) (allow, remember bool)

// domainMemory is the per-domain remembered set. Safe for concurrent use.
type domainMemory struct {
	mu    sync.Mutex
	path  string // "" = in-memory only (tests, or a composition that has no data dir yet)
	set   map[string]bool
	ready bool
}

func newDomainMemory(path string) *domainMemory {
	return &domainMemory{path: path, set: map[string]bool{}}
}

func (m *domainMemory) known(domain string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadLocked()
	return m.set[domain]
}

func (m *domainMemory) remember(domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadLocked()
	if m.set[domain] {
		return
	}
	m.set[domain] = true
	m.persistLocked()
}

// loadLocked reads the store file once. A broken file is treated as empty
// (and the next successful remember rewrites it): a corrupt domain memory
// may not silently ALLOW anything — the memory only ever SUPPRESSES prompts,
// and losing it means more prompts, which is the fail-closed direction.
func (m *domainMemory) loadLocked() {
	if m.ready || m.path == "" {
		m.ready = true
		return
	}
	m.ready = true
	b, err := os.ReadFile(m.path)
	if err != nil {
		return
	}
	var doc struct {
		Version int      `json:"version"`
		Domains []string `json:"domains"`
		Updated string   `json:"updated"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return
	}
	for _, d := range doc.Domains {
		if d = normalizeDomain(d); d != "" {
			m.set[d] = true
		}
	}
}

func (m *domainMemory) persistLocked() {
	if m.path == "" {
		return
	}
	doms := make([]string, 0, len(m.set))
	for d := range m.set {
		doms = append(doms, d)
	}
	// stable order so the file diff is honest
	for i := 1; i < len(doms); i++ {
		for j := i; j > 0 && doms[j] < doms[j-1]; j-- {
			doms[j], doms[j-1] = doms[j-1], doms[j]
		}
	}
	doc := struct {
		Version int      `json:"version"`
		Domains []string `json:"domains"`
		Updated string   `json:"updated"`
	}{Version: 1, Domains: doms, Updated: time.Now().UTC().Format(time.RFC3339)}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return
	}
	dir := filepath.Dir(m.path)
	tmp, err := os.CreateTemp(dir, ".wisp-domains-*")
	if err != nil {
		return
	}
	name := tmp.Name()
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return
	}
	if err := os.Rename(name, m.path); err != nil {
		_ = os.Remove(name)
	}
}

// normalizeDomain lowercases and strips a trailing dot ("EXAMPLE.COM." ->
// "example.com"), punycode in / out without conversion (the guard never
// re-encodes a host).
func normalizeDomain(host string) string {
	h := strings.ToLower(strings.TrimSpace(host))
	h = strings.TrimSuffix(h, ".")
	if strings.HasPrefix(h, "[") { // bracketed IPv6 literal
		h = strings.TrimSuffix(strings.TrimPrefix(h, "["), "]")
	}
	return h
}

// firstVisitGate wraps the seam + memory into the single question the web
// tools ask before any connect/open.
type firstVisitGate struct {
	mem *domainMemory
	ask FirstVisitAsk
	log func(string, ...any)
}

// ensure returns nil when the visit may proceed. errRefused marks a user
// "no"; a missing seam (ask == nil) refuses fail-closed.
var errFirstVisitRefused = errors.New("首次访问该域名未获用户允许")

func (f *firstVisitGate) ensure(ctx context.Context, tool, rawHost string) error {
	domain := normalizeDomain(rawHost)
	if domain == "" {
		return errors.New("URL 缺少主机名")
	}
	// IP literals are not "new domains" — the SSRF guard already judges
	// addresses, and an address memory would grow without bound per scan.
	// The prompt is for DNS domains (SPEC-06 §6.5 wording "新域名").
	if _, encoded := decodeNumericHost(domain); encoded {
		return nil
	}
	if f.mem.known(domain) {
		return nil
	}
	if f.ask == nil {
		if f.log != nil {
			f.log("tools: D30-5 first-visit prompt channel unwired, refusing %s -> %s (fail-closed)", tool, domain)
		}
		return fmt.Errorf("%w：首次访问提示通道尚未接入（fail-closed）", errFirstVisitRefused)
	}
	allow, remember := f.ask(ctx, tool, domain)
	if !allow {
		return fmt.Errorf("%w：%s", errFirstVisitRefused, domain)
	}
	if remember {
		f.mem.remember(domain)
	}
	return nil
}
