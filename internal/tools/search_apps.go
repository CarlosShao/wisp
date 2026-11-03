package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// D34 `search.apps` (docs/PLAN.md:2551 / SPEC-07 §3 row: "检索已安装应用", L0,
// capability sysinfo, "只检索不启动（启动走 app.launch）"; ticket 24:23-24
// "installed-app enumeration feeding app.launch; enumeration only, never
// launches", AC#4 "search.apps enumerates fixture-installed app; feeding result
// to app.launch works (22)").
//
// The rule that makes AC#4 more than a hope: a candidate is only LISTED after
// it has been run through `resolveApp` - the very resolver `app.launch` uses
// (launcher_windows.go / launcher_other.go, untouched). Enumeration therefore
// emits either an absolute executable path or a validated AUMID, which is
// exactly the target shape the launch arm consumes. A candidate no resolver can
// turn into a launch target is counted and dropped, so the listing never
// promises an app the launch tool could not start.
//
// "Enumeration only, never launches" is structural, not a promise in a comment:
// this file calls the resolver, never `startApp`, and it declares capability
// `sysinfo` only - not `shell` - so if it ever grew a launch call the bridge
// would hard-reject it for needing a capability it never declared (SPEC-07 §2).
// The result line repeats the fact for the reader.
//
// Sources are an injectable seam (AppSource) so a test fixture can be the
// "installed app" AC#4 names, without this tool owning a registry or a process.
// Default sources: the App Paths registry hive on Windows (registered from
// search_apps_windows.go, the same hive windowsAppPathsResolver reads) plus the
// executable files in the PATH directories on every platform.
//
// What is NOT here, and why it is reported rather than invented: a UWP/AUMID
// enumeration needs either the Appx packaging API (a new dependency, whitelist
// change = human approval) or a spawned shell command (a launch-shaped side
// effect this L0 sysinfo tool must not take, and it would need `shell`). So no
// AUMID is EMITTED by the default sources today; the seam and the emit shape
// already carry `kind=aumid`, and `app.launch` accepts a user-supplied AUMID
// through its own shape gate. See docs/evidence/s1/24r3-impl.md §4.

// InstalledAppCandidate is one name a source says is installed.
type InstalledAppCandidate struct {
	Name   string
	AUMID  string // non-empty = a UWP activation id, handed over verbatim
	Source string
}

// AppSource is one enumeration source. It lists; it never launches.
type AppSource interface {
	// Label names the source for the result listing and the logs.
	Label() string
	// Apps returns the candidate apps visible to the current user. The ctx is
	// the call's budget: a source that cannot finish checks it.
	Apps(ctx context.Context) ([]InstalledAppCandidate, error)
}

// SearchAppsDeps is what search.apps needs from the host.
type SearchAppsDeps struct {
	// MaxResults caps the apps one call lists (default
	// defaultMaxSearchApps). Hitting it sets Truncated, the spill signal the
	// loop reads, exactly like the fs.read siblings.
	MaxResults int
	// MaxCandidates caps how many names a single source may hand over before
	// resolution, so one pathological PATH directory cannot turn this into an
	// unbounded resolver loop.
	MaxCandidates int
	// Sources is the enumeration seam; nil selects DefaultAppSources(). A test
	// fixture lands here (AC#4's "fixture-installed app").
	Sources []AppSource
	// Resolver is the C1 name-to-executable seam shared with app.launch
	// (WebDeps.AppResolver); nil means the platform default, which is what
	// app.launch itself would use.
	Resolver AppResolver
	// Logf receives host-side notes; nil means discard.
	Logf func(string, ...any)
}

const (
	defaultMaxSearchApps    = 50
	defaultMaxAppCandidates = 2000
)

func (d SearchAppsDeps) resultsCap() int {
	if d.MaxResults > 0 {
		return d.MaxResults
	}
	return defaultMaxSearchApps
}

func (d SearchAppsDeps) candidatesCap() int {
	if d.MaxCandidates > 0 {
		return d.MaxCandidates
	}
	return defaultMaxAppCandidates
}

func (d SearchAppsDeps) logf() func(string, ...any) {
	if d.Logf != nil {
		return d.Logf
	}
	return func(string, ...any) {}
}

// platformAppSources are the OS-specific enumeration sources. nil on platforms
// without one (the PATH source below is the cross-platform floor); Windows
// registers the App Paths registry hive from search_apps_windows.go.
var platformAppSources []AppSource

// DefaultAppSources is the enumeration set when the host injects none.
func DefaultAppSources() []AppSource {
	out := make([]AppSource, 0, len(platformAppSources)+1)
	out = append(out, platformAppSources...)
	return append(out, pathDirsAppSource{})
}

// pathDirsAppSource lists executable files found in the PATH directories. Every
// entry it emits is later handed to the same resolver app.launch uses, so a
// stale or non-executable hit drops out instead of being advertised.
type pathDirsAppSource struct{}

// Label implements AppSource.
func (pathDirsAppSource) Label() string { return "PATH 目录" }

// Apps implements AppSource.
func (pathDirsAppSource) Apps(ctx context.Context) ([]InstalledAppCandidate, error) {
	var out []InstalledAppCandidate
	seen := make(map[string]bool)
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		if strings.TrimSpace(dir) == "" {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // a PATH entry that is gone or unreadable is not a failure
		}
		for _, e := range entries {
			if len(out) >= defaultMaxAppCandidates {
				return out, nil
			}
			if e.IsDir() || !looksExecutable(e) {
				continue
			}
			key := strings.ToLower(e.Name())
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, InstalledAppCandidate{Name: e.Name(), Source: "PATH 目录"})
		}
	}
	return out, nil
}

// looksExecutable is the cheap shape filter for the PATH source: a Windows
// executable extension, or a user-executable mode bit elsewhere. The
// authoritative veto is the resolver that follows, not this predicate.
func looksExecutable(e os.DirEntry) bool {
	if runtime.GOOS == "windows" {
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".exe", ".com":
			return true
		}
		return false
	}
	info, err := e.Info()
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0o111 != 0
}

type searchApps struct{ d SearchAppsDeps }

type searchAppsArgs struct {
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
}

var searchAppsSchema = JSONSchema(`{"type":"object","properties":{` +
	`"query":{"type":"string","description":"应用名字面子集（大小写不敏感，非正则）；留空＝列出全部已检索到的应用"},` +
	`"max_results":{"type":"integer","minimum":1,"description":"本次最多列出多少个应用"}}` +
	`,"additionalProperties":false}`)

func (searchApps) Name() string { return "search.apps" }

// Description implements Tool.
func (searchApps) Description() string {
	return "检索当前用户可见的已安装应用（只检索不启动；结果可直接交给 app.launch）"
}

// Parameters implements C1 Tool.
func (searchApps) Parameters() JSONSchema { return searchAppsSchema }

// Execute implements C1 Tool.
func (t searchApps) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a searchAppsArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	needle := strings.ToLower(strings.TrimSpace(a.Query))

	sources := t.d.Sources
	if len(sources) == 0 {
		sources = DefaultAppSources()
	}
	limit := t.d.resultsCap()
	shared := &webShared{deps: WebDeps{AppResolver: t.d.Resolver}, logf: t.d.logf()}

	var (
		apps       []launchedAppShape
		seen       = make(map[string]bool)
		candidates int
		unresolved int
		errored    int
		aborted    string
		hitCap     bool
	)
	for _, src := range sources {
		list, err := src.Apps(ctx)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				aborted = ctxErr.Error()
				break
			}
			errored++
			t.d.logf()("tools: search.apps source %q failed: %v", src.Label(), err)
			continue
		}
		for _, cand := range list {
			if candidates >= t.d.candidatesCap() {
				hitCap = true
				break
			}
			name := strings.TrimSpace(cand.Name)
			if name == "" {
				continue
			}
			if needle != "" && !strings.Contains(strings.ToLower(name), needle) {
				continue
			}
			candidates++
			key := strings.ToLower(name)
			if seen[key] {
				continue
			}
			target, rerr := shared.resolveApp(ctx, name)
			if rerr != nil {
				// Never advertised: an app this resolver cannot turn into a
				// launch target is not something app.launch could start.
				unresolved++
				seen[key] = true
				continue
			}
			seen[key] = true
			shape := launchedAppShape{Name: name, Kind: target.Kind, Source: src.Label()}
			switch target.Kind {
			case kindAUMID:
				shape.Target = target.AUMID
			default:
				shape.Target = target.Path
			}
			apps = append(apps, shape)
			if len(apps) >= limit {
				hitCap = true
				break
			}
		}
		if hitCap || aborted != "" {
			break
		}
	}
	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Name < apps[j].Name {
			return true
		}
		return apps[i].Target < apps[j].Target
	})

	head := fmt.Sprintf("search.apps 检索到 %d 个应用（候选 %d，未能解析 %d，上限 %d",
		len(apps), candidates, unresolved, limit)
	if aborted != "" {
		head += fmt.Sprintf("，已中止：%s，truncated", aborted)
	} else if hitCap {
		head += "，达到上限，truncated"
	}
	if errored > 0 {
		head += fmt.Sprintf("，检索源失败 %d 个", errored)
	}
	head += "；只检索，未启动任何应用)\n"

	if len(apps) == 0 && errored == len(sources) {
		// Every source failed is an error, not an honest zero - reporting "no
		// apps installed" here would be a silent lie about the machine.
		return Result{Text: strings.TrimRight(head, "\n") + "\n检索失败：所有枚举源均报错", IsError: true}, nil
	}

	var b strings.Builder
	b.WriteString(head)
	for _, ap := range apps {
		fmt.Fprintf(&b, "%s\t%s\t%s\t(来源:%s)\n", ap.Name, ap.Kind, ap.Target, ap.Source)
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("search.apps 检索到 %d 个应用", len(apps)))
	}
	truncated := hitCap || aborted != ""
	return Result{Text: b.String(), Truncated: truncated}, nil
}

// launchedAppShape is one listed app: the name a user would say, plus the
// target string app.launch consumes and which arm of it it is.
type launchedAppShape struct {
	Name   string
	Kind   string
	Target string
	Source string
}

// kindAUMID is the launcher's UWP arm name (launcher_windows.go).
const kindAUMID = "aumid"

// SearchAppsDecl is the host-side declaration for search.apps: capability
// sysinfo (D34 row), R1 floor L0, and deliberately NO shell capability and NO
// PathParams - this tool starts nothing and reads no user path, so a launch
// attempt inside it could not be declared, and would be hard-rejected.
func SearchAppsDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapSysinfo},
		Needs:        []Capability{CapSysinfo},
		Declared:     risk.L0,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinSearchAppsEntries returns the search.apps entry for the Builtin slot.
// Like the search.files / search.content / doc.read entries it is not wired into
// bridge.go in this limb - see docs/evidence/s1/24r3-impl.md §5.
func BuiltinSearchAppsEntries(d SearchAppsDeps) []Entry {
	return []Entry{{Tool: searchApps{d: d}, Decl: SearchAppsDecl()}}
}
