package tools

import (
	"context"
	"net/netip"
	"os"
	"os/exec"
	"time"
)

// web.go is the composition face for the D30 web/open family (ticket 22):
// web.search (L0, R4-guarded query), web.fetch (L0, SSRF-hard-denied,
// Markdown-only), web.open (L1, http/https strict), file.open (L1, the
// separate 16.5.3 protocol allowlist), app.launch (L1, C30 Job).
//
// Everything the family needs from the host arrives through WebDeps; every
// zero value resolves to the FAIL-CLOSED shape (guarded client, provider
// without a base URL answers "not configured", no prompt channel = refusals
// on new domains), so a half-wired composition cannot quietly run the family
// with the D30 layers off. The wiring ticket (12, per R12) fills Deps from
// the frozen config faces it already has: [net].allowlist -> NetAllowlist,
// the DPAPI/env key source -> SearchKey, the data dir -> DomainStorePath,
// the process-wide JobScope -> Jobs.

// WebDeps configures the whole web/open family.
type WebDeps struct {
	// Paths is the C26 canonicalizer + [fs] allowlist judge used by the
	// file.open path arms (16.5.3: "file: 仅授权目录内 + 经 C26 规范化").
	Paths *PathCanonicalizer

	// NetAllowlist mirrors [net] allowlist[] and is fed to R5 through the
	// Facts hooks of every model-supplied target. Empty = domain gate off
	// (the config owns that choice; see unwired.go on why the composition
	// root, not this package, decides when the list goes live).
	NetAllowlist []string

	// AskFirstVisit is D30 layer 5's prompt seam. nil = the channel is not
	// wired and every NEW domain visit is refused fail-closed (same posture
	// as NoGate for L1/L2). Remembered domains keep working without it.
	AskFirstVisit FirstVisitAsk
	// DomainStorePath is where the remembered-domain file lives ("" =
	// session memory only: the ask still runs once per domain per process,
	// nothing persists; the composition root passes the wisp data-dir path).
	DomainStorePath string

	// SearchProvider is the R4-ruled implementation face (anysearch
	// compatible). nil = the default provider is an AnySearchProvider over
	// SearchBaseURL/SearchKey below.
	SearchProvider SearchProvider
	// SearchBaseURL is the provider egress for the default face (production:
	// https://api.anysearch.com/v1, ledger R9). Empty = web.search answers
	// "未配置 baseUrl" and never dials.
	SearchBaseURL string
	// SearchKey returns the API key at call time (from internal/secret or
	// the user's environment reference — R7: never config plaintext, never
	// argv, never logged). Empty string = anonymous (rate-limited) call.
	SearchKey func() string

	// LookupIP overrides DNS resolution inside the guard (nil =
	// net.DefaultResolver). The DNS-rebinding fixture flips answers
	// through exactly this seam.
	LookupIP func(ctx context.Context, host string) ([]netip.Addr, error)

	// MaxResponseBytes caps one web.fetch body (default: the spill raw hard
	// ceiling at reference scale, 1 MiB — internal/agent/budgets.go
	// refRawOutputBytes; the loop's Spiller owns everything above the
	// context budget after that).
	MaxResponseBytes int
	// FetchTimeout bounds one web.fetch end to end (C22 Decl.Timeout is the
	// outer belt; this is the tool's own). Zero = defaultFetchTimeout.
	FetchTimeout time.Duration
	// DialTimeout bounds resolution + connect per candidate address.
	DialTimeout time.Duration
	// MaxRedirects caps the redirect chain (zero = defaultMaxRedirects).
	MaxRedirects int

	// IsKnownMsProtocol answers "known ms-* application protocol" for the
	// 16.5.3 allowlist (registry-fed on Windows; tests inject a roster).
	IsKnownMsProtocol func(scheme string) bool

	// AppResolver maps an app NAME to a launch target (the seam ticket 23's
	// search.apps will share). nil = the platform default (App Paths
	// registry on Windows, PATH lookup elsewhere). Absolute paths and AUMID
	// shapes bypass the resolver before it is asked.
	AppResolver AppResolver

	// Jobs is the C30 Job Object (internal/proc.JobScope satisfies it on
	// Windows). When set, every app.launch child enters it via StartInJob.
	// nil = children start outside any wisp job — allowed on platforms
	// without a Job implementation, and RECORDED in the launch result text
	// so the model/user sees the weaker posture.
	Jobs JobStarter

	// Launch overrides the OS hand-off for web.open/file.open ("url",
	// "file", "uri" deliveries). Production leaves it nil (ShellExecute on
	// Windows, the honest "no door here" refusal elsewhere); the protocol
	// matrix suite injects a recorder so 16.5.3 can be pinned without a
	// browser ever starting.
	Launch func(ctx context.Context, kind, target string) error

	// Logf is the audit sink. It receives URLs, hosts, statuses — never
	// response bodies, never keys.
	Logf func(format string, args ...any)
}

// JobStarter is the C30 seam: it starts cmd and lands the child in the
// process-wide Job Object before it can do work (proc.JobScope.StartInJob is
// the Windows implementation; same signature by construction).
type JobStarter interface {
	StartInJob(cmd *exec.Cmd) (*os.Process, error)
}

// AppResolver turns a user-visible app name into an executable path
// (search.apps' resolution table, ticket 23, lands behind this same seam).
type AppResolver interface {
	ResolveApp(name string) (path string, ok bool)
}

// webShared is the resolved (defaults-applied) family state, built once per
// BuiltinWebEntries call and shared by the entries' Execute closures.
type webShared struct {
	deps             WebDeps
	guard            *netGuard
	firstVisit       *firstVisitGate
	provider         SearchProvider
	fetchTimeout     time.Duration
	searchTimeout    time.Duration
	maxResponseBytes int
	logf             func(string, ...any)
}

// newWebShared applies the fail-closed defaults. One netGuard, one domain
// memory, one provider binding — the family shares exactly the state D30
// describes as load-bearing.
func newWebShared(d WebDeps) *webShared {
	log := d.Logf
	if log == nil {
		log = func(string, ...any) {}
	}
	fetchTimeout := orDefaultDur(d.FetchTimeout, defaultFetchTimeout)
	w := &webShared{
		deps: d,
		guard: newNetGuard(guardOptions{
			lookup:       d.LookupIP,
			dialTimeout:  orDefaultDur(d.DialTimeout, defaultDialTimeout),
			fetchTimeout: fetchTimeout,
			maxRedirects: orDefaultInt(d.MaxRedirects, defaultMaxRedirects),
		}),
		fetchTimeout:     fetchTimeout,
		searchTimeout:    searchBudget,
		maxResponseBytes: orDefaultInt(d.MaxResponseBytes, defaultFetchMaxBytes),
		logf:             log,
		firstVisit: &firstVisitGate{
			mem: newDomainMemory(d.DomainStorePath),
			ask: d.AskFirstVisit,
			log: log,
		},
	}
	if d.SearchProvider != nil {
		w.provider = d.SearchProvider
	} else {
		// The default face IS the R4 ruling: anysearch-compatible, dialed
		// THROUGH the guard. With no base URL set it answers "not
		// configured" and never connects — the honest half-wired state.
		w.provider = NewAnySearchProvider(d.SearchBaseURL, d.SearchKey, nil)
	}
	return w
}

func orDefaultDur(v, def time.Duration) time.Duration {
	if v > 0 {
		return v
	}
	return def
}

// launch is the dispatcher: an injected hand-off (the matrix suite, or a
// host that routes opens through its own broker) wins first; otherwise the
// platform file delivers. The scheme VERDICT is never here — it is made in
// web_open.go + the bridge's R5 before anything reaches this line.
func (w *webShared) launch(ctx context.Context, kind, target string) error {
	if w.deps.Launch != nil {
		return w.deps.Launch(ctx, kind, target)
	}
	return w.platformLaunch(ctx, kind, target)
}

// fetchBudget is the C22 outer belt for web.fetch: a hair above the tool's
// own fetch timeout so the tool's structured text usually wins, and the
// bridge only catches a tool that ignores its context.
func (d WebDeps) fetchBudget() time.Duration {
	return orDefaultDur(d.FetchTimeout, defaultFetchTimeout) + 2*time.Second
}

// ---------------------------------------------------------------------------
// registration (D34 rows, one entry point like the fs family's)
// ---------------------------------------------------------------------------

// BuiltinWebEntries returns the whole ticket-22 family with defaults applied:
//
//	web.search   L0 (query taint -> L2 via R4, computed by the bridge)
//	web.fetch    L0 (SSRF hard-deny in the netGuard, before DNS/socket)
//	web.open     L1 (http/https only)
//	file.open    L1 (16.5.3 allowlist; file: arms through C26)
//	app.launch   L1 (C30 Job when WebDeps.Jobs is wired)
//
// Register them through the same Registry the bridge owns
// (Registry().Register); nothing here reaches the bridge by itself, same
// rule as the fs family.
func BuiltinWebEntries(d WebDeps) []Entry {
	w := newWebShared(d)
	return []Entry{
		{Tool: webSearch{w: w}, Decl: WebSearchDecl()},
		{Tool: webFetch{w: w}, Decl: WebFetchDecl(d)},
		{Tool: webOpen{w: w}, Decl: WebOpenDecl(d)},
		{Tool: fileOpen{w: w}, Decl: FileOpenDecl()},
		{Tool: appLaunch{w: w}, Decl: AppLaunchDecl()},
	}
}
