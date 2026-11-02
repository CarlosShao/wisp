package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// webopen.go lands D34's three "open doors" and the 16.5.3 split that solves
// the D23 <-> D30(3) contradiction (PLAN.md:2577):
//
//	web.open   ONLY http/https. Everything else is rejected and the rejection
//	           names file.open as the door that owns it (D30(3) strict side).
//	file.open  a SEPARATE, default-DENY protocol allowlist:
//	           file: (authorized dirs only, through C26), mailto:,
//	           ms-settings:, registered ms-* app protocols;
//	           explicitly denied: search-ms: (CLSID-hijack class), mshta:,
//	           javascript:, vbscript:, shell:, tel: — and every unknown
//	           scheme lands on the deny side because the list is an
//	           allowlist, not a blocklist ("新协议出现时自动落到拒绝侧").
//	app.launch a name / path / UWP AUMID resolved to a process, and the
//	           child ENTERS the C30 Job Object (D32/D38e): the launched
//	           process tree must not outlive wisp, and the guard against a
//	           prompt-injected "launch something" rides that same door.
//
// The OS hand-off itself is a seam (Launcher, below) so the matrix and the
// C26 path rule are testable on every platform while the real ShellExecute
// / Job assignment lives in the windows files.

// ---------------------------------------------------------------------------
// the 16.5.3 scheme matrix (data, not code paths)
// ---------------------------------------------------------------------------

// openMatrix is the 16.5.3 table, expressed as one classifier function and a
// data map — classifyOpenTarget below is the single decision point every
// file.open call passes, and web_matrix_test.go walks the whole allow/deny
// grid against it.

// explicitDenySchemes are named in 16.5.3 as hard denials for every web/open
// door. A scheme in this map is rejected even if the OS has registered a
// handler for it — registration is not authorization (search-ms: ships with
// Windows and is still on the deny list).
var explicitDenySchemes = map[string]string{
	"search-ms":  "search-ms: 被明确拒绝（可被用于 CLSID 劫持执行，16.5.3）",
	"mshta":      "mshta: 被明确拒绝（脚本宿主，16.5.3）",
	"javascript": "javascript: 被明确拒绝（脚本协议，16.5.3）",
	"vbscript":   "vbscript: 被明确拒绝（脚本协议，16.5.3）",
	"shell":      "shell: 被明确拒绝（shell 命名空间执行，16.5.3）",
	"tel":        "tel: 被明确拒绝（16.5.3）",
}

// classifyOpenTarget resolves one file.open target against the 16.5.3
// allowlist. raw is whatever the model passed (a bare path is the common
// case). Returns the resolved scheme ("" = a filesystem path), and a refusal
// when the scheme is not on the allow side.
func classifyOpenTarget(d WebDeps, raw string) (scheme string, target string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", errors.New("缺少 path 参数")
	}
	u, perr := url.Parse(raw)
	if perr != nil {
		u = &url.URL{}
	}
	scheme = strings.ToLower(u.Scheme)
	driveLike := len(scheme) == 1 // one-letter "scheme" is a Windows drive letter, not a protocol
	if driveLike {
		scheme = ""
	}
	if scheme == "" {
		if perr != nil && looksLikeURIScheme(raw) && !looksLikeWindowsDrivePath(raw) {
			// scheme-shaped but unparseable and not a drive path: fail closed
			// as an unknown protocol, never fall through to the path arm.
			return "", "", fmt.Errorf("协议 %q 不在 file.open 白名单，按默认拒绝处理（16.5.3 一切未知协议拒绝）", firstSchemeToken(raw))
		}
		// a filesystem path: it goes through C26 and the [fs] allowlist in
		// the tool's Execute, same rule as fs.read's out-of-scope shape.
		return "", raw, nil
	}
	if reason, deny := explicitDenySchemes[scheme]; deny {
		return "", "", errors.New(reason)
	}
	switch scheme {
	case "http", "https":
		return "", "", errors.New("web 协议请用 web.open（file.open 的白名单不含 http/https，16.5.3 拆分）")
	case "file":
		return "file", raw, nil
	case "mailto":
		return "mailto", raw, nil
	case "ms-settings":
		if !msPrefixRegistered(d, scheme) {
			return "", "", errors.New("ms-settings: 在本机未注册（默认拒绝侧）")
		}
		return "ms-settings", raw, nil
	default:
		if strings.HasPrefix(scheme, "ms-") {
			if !msPrefixRegistered(d, scheme) {
				return "", "", fmt.Errorf("协议 %q 未在本机注册为已知 ms-* 应用协议，默认拒绝（16.5.3）", scheme)
			}
			return scheme, raw, nil
		}
		return "", "", fmt.Errorf("协议 %q 不在白名单，默认拒绝（16.5.3）", scheme)
	}
}

// looksLikeWindowsDrivePath reports "c:\..." / "c:/..." shapes (checked
// before the generic scheme test, since Go url-parses those inconsistently).
func looksLikeWindowsDrivePath(raw string) bool {
	if len(raw) < 3 {
		return false
	}
	c0 := raw[0]
	if !((c0 >= 'a' && c0 <= 'z') || (c0 >= 'A' && c0 <= 'Z')) || raw[1] != ':' {
		return false
	}
	return raw[2] == '\\' || raw[2] == '/'
}

// looksLikeURIScheme reports whether "x:y..." is scheme-shaped (RFC3986 head:
// ALPHA *( ALPHA / DIGIT / "+" / "-" / "." ) ":"). A remainder starting with
// '/' or '\' is a path (or a drive path), never a bare opaque scheme target.
func looksLikeURIScheme(raw string) bool {
	i := strings.IndexByte(raw, ':')
	if i <= 0 {
		return false
	}
	for j := 0; j < i; j++ {
		c := raw[j]
		ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' ||
			j > 0 && (c >= '0' && c <= '9' || c == '+' || c == '-' || c == '.')
		if !ok {
			return false
		}
	}
	rest := raw[i+1:]
	if len(rest) > 0 && (rest[0] == '/' || rest[0] == '\\') {
		return false
	}
	return true
}

// firstSchemeToken extracts "foo:" from "foo:bar" for the denial message.
func firstSchemeToken(raw string) string {
	if i := strings.IndexByte(raw, ':'); i > 0 {
		return strings.ToLower(raw[:i+1])
	}
	return raw
}

// msPrefixRegistered answers "known ms-* protocol" for 16.5.3. The production
// reader is the per-machine protocol handler registration (launcher_windows.go);
// the injected seam keeps the matrix testable on any runner without touching
// a real registry.
func msPrefixRegistered(d WebDeps, scheme string) bool {
	if d.IsKnownMsProtocol != nil {
		return d.IsKnownMsProtocol(scheme)
	}
	return isKnownMsProtocolDefault(scheme)
}

// ---------------------------------------------------------------------------
// web.open
// ---------------------------------------------------------------------------

type webOpen struct{ w *webShared }

type webOpenArgs struct {
	URL string `json:"url"`
}

var webOpenSchema = JSONSchema(`{"type":"object","properties":{` +
	`"url":{"type":"string","description":"要交给默认浏览器打开的 URL（仅 http/https）"}}` +
	`,"required":["url"],"additionalProperties":false}`)

func (webOpen) Name() string { return "web.open" }

func (webOpen) Description() string {
	return "在用户的默认浏览器里打开一个 http/https 链接（本地文件与其他协议一律拒绝，改走 file.open；16.5.3）"
}

func (webOpen) Parameters() JSONSchema { return webOpenSchema }

// Execute implements Tool.
func (t webOpen) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a webOpenArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	raw := strings.TrimSpace(a.URL)
	if raw == "" {
		return Result{Text: "缺少 url 参数", IsError: true}, nil
	}
	if len(raw) > risk.MaxURLLength {
		return Result{Text: fmt.Sprintf("web.open: URL 超过 %d 字符上限（D30③）", risk.MaxURLLength), IsError: true}, nil
	}
	u, err := normalizeGuardedURL(raw)
	if err != nil {
		return Result{Text: "web.open: " + err.Error(), IsError: true}, nil
	}
	if host := u.Hostname(); host != "" {
		if err := t.w.firstVisit.ensure(ctx, "web.open", host); err != nil {
			return Result{Text: "web.open: " + err.Error(), IsError: true}, nil
		}
	}
	if err := t.w.launch(ctx, "url", u.String()); err != nil {
		return Result{Text: "web.open: " + err.Error(), IsError: true}, nil
	}
	return Result{Text: "已在默认浏览器打开：" + u.String(), Origin: u.String()}, nil
}

// ---------------------------------------------------------------------------
// file.open
// ---------------------------------------------------------------------------

type fileOpen struct{ w *webShared }

type fileOpenArgs struct {
	Path string `json:"path"`
}

var fileOpenSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"本地文件路径（必须落在 [fs] 授权目录内，经 C26 规范化）、file: URL、mailto: 或已知 ms-* 协议 URI"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (fileOpen) Name() string { return "file.open" }

func (fileOpen) Description() string {
	return "用系统默认程序打开本地文件（仅 [fs] 授权目录内）或 mailto:/ms-settings:/已知 ms-* 协议 URI；未知协议与危险协议默认拒绝（16.5.3）"
}

func (fileOpen) Parameters() JSONSchema { return fileOpenSchema }

// Execute implements Tool.
func (t fileOpen) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fileOpenArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	scheme, target, err := classifyOpenTarget(t.w.deps, a.Path)
	if err != nil {
		return Result{Text: "file.open: " + err.Error(), IsError: true}, nil
	}
	switch scheme {
	case "":
		// plain path: same C26 + allowlist rule as every fs tool —
		// "file: 仅授权目录内 + 经 C26 规范化".
		if t.w.deps.Paths == nil {
			return Result{Text: "file.open: 未配置 C26 解析器，fail-closed 拒绝", IsError: true}, nil
		}
		c, err := t.w.deps.Paths.Canonicalize(target)
		if err != nil {
			return Result{Text: "file.open: 路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
		}
		if !t.w.deps.Paths.InAllowlist(c) {
			return Result{Text: "file.open: 路径不在 [fs] allowed_dirs 授权目录内，拒绝（file: 仅授权目录内，16.5.3）", IsError: true}, nil
		}
		if err := t.w.launch(ctx, "file", c); err != nil {
			return Result{Text: "file.open: " + err.Error(), IsError: true}, nil
		}
		return Result{Text: "已交给系统默认程序打开：" + c, Origin: c}, nil
	case "file":
		c, err := t.w.fileURLToAuthorizedPath(target)
		if err != nil {
			return Result{Text: "file.open: " + err.Error(), IsError: true}, nil
		}
		if err := t.w.launch(ctx, "file", c); err != nil {
			return Result{Text: "file.open: " + err.Error(), IsError: true}, nil
		}
		return Result{Text: "已交给系统默认程序打开：" + c, Origin: c}, nil
	default:
		// mailto: / ms-settings: / ms-* : the URI rides untouched to the
		// registered handler; explicit-deny shapes never reach here.
		if err := t.w.launch(ctx, "uri", target); err != nil {
			return Result{Text: "file.open: " + err.Error(), IsError: true}, nil
		}
		return Result{Text: fmt.Sprintf("已按 %s: 协议交给系统处理（默认拒绝白名单之外不会走到这里）", scheme)}, nil
	}
}

// fileURLToAuthorizedPath turns a file: URL into a C26-canonical, allowlist-
// checked path. Remote file://hosts are refused outright (there is no
// authorized remote reader in this product).
func (w *webShared) fileURLToAuthorizedPath(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("file: URL 无法解析：%w", err)
	}
	if h := u.Hostname(); h != "" && !strings.EqualFold(h, "localhost") {
		return "", fmt.Errorf("file: 仅允许本机路径（主机 %q 被拒绝）", h)
	}
	p := u.Path
	// file:///C:/x on windows spells the path "/C:/x"; the C26 resolver
	// wants the drive-letter form. Only strip when the shape IS a drive
	// path (one letter + ':'), never for ordinary POSIX roots.
	if len(p) >= 3 && p[0] == '/' && isLetter(p[1]) && p[2] == ':' {
		p = p[1:]
	}
	if w.deps.Paths == nil {
		return "", errors.New("未配置 C26 解析器，fail-closed 拒绝")
	}
	c, err := w.deps.Paths.Canonicalize(p)
	if err != nil {
		return "", err
	}
	if !w.deps.Paths.InAllowlist(c) {
		return "", errors.New("file: 目标不在 [fs] allowed_dirs 授权目录内（16.5.3：file: 仅授权目录内）")
	}
	return c, nil
}

// ---------------------------------------------------------------------------
// app.launch
// ---------------------------------------------------------------------------

type appLaunch struct{ w *webShared }

type appLaunchArgs struct {
	App  string   `json:"app"`
	Args []string `json:"args"`
}

var appLaunchSchema = JSONSchema(`{"type":"object","properties":{` +
	`"app":{"type":"string","description":"应用名（search.apps 可检索到的名字）、绝对路径，或 UWP AUMID（含 ! 的包名）"},` +
	`"args":{"type":"array","items":{"type":"string"},"description":"传给应用的参数（argv 向量，不拼接字符串，D33）"}}` +
	`,"required":["app"],"additionalProperties":false}`)

func (appLaunch) Name() string { return "app.launch" }

func (appLaunch) Description() string {
	return "启动本机应用（名称 / 绝对路径 / UWP AUMID）。启动的进程进入 Wisp 的 Job Object（C30），随宿主一起终止"
}

func (appLaunch) Parameters() JSONSchema { return appLaunchSchema }

// Execute implements Tool.
func (t appLaunch) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a appLaunchArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	name := strings.TrimSpace(a.App)
	if name == "" {
		return Result{Text: "缺少 app 参数", IsError: true}, nil
	}
	for _, s := range a.Args {
		if len(name)+len(s) > 4096 {
			return Result{Text: "app.launch: 参数总长超出预算，拒绝", IsError: true}, nil
		}
	}
	target, err := t.w.resolveApp(ctx, name)
	if err != nil {
		return Result{Text: "app.launch: " + err.Error(), IsError: true}, nil
	}
	posture, err := t.w.startApp(ctx, target, a.Args)
	if err != nil {
		return Result{Text: "app.launch: " + err.Error(), IsError: true}, nil
	}
	return Result{Text: fmt.Sprintf("已启动 %s（%s 类目标，%s）", name, target.Kind, posture)}, nil
}

// launchTarget is one resolved app.launch shape.
type launchTarget struct {
	Kind  string // "path" | "aumid"
	Path  string // for path targets
	AUMID string // for UWP targets
}

// ---------------------------------------------------------------------------
// declarations
// ---------------------------------------------------------------------------

// WebOpenDecl declares web.open (D34: L1, capability net; R5 judges the
// target through the same url Facts hook as web.fetch).
func WebOpenDecl(d WebDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapNet},
		Needs:        []Capability{CapNet},
		Declared:     risk.L1,
		Resident:     true,
		Provider:     KindBuiltin,
		Facts:        netURLFacts(d.NetAllowlist),
	}
}

// FileOpenDecl declares file.open (D34: L1, capability shell; the path arm
// rides C26/R2/R3 like every fs target, the protocol arms are owned by the
// 16.5.3 matrix inside Execute).
func FileOpenDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapShell},
		Needs:        []Capability{CapShell},
		Declared:     risk.L1,
		Resident:     false,
		Provider:     KindBuiltin,
		PathParams:   []string{"path"},
	}
}

// AppLaunchDecl declares app.launch (D34: L1, capability shell, unified L1
// per "从 search 里独立出来成为一等公民工具；统一 L1").
func AppLaunchDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapShell},
		Needs:        []Capability{CapShell},
		Declared:     risk.L1,
		Resident:     false,
		Provider:     KindBuiltin,
	}
}

// isLetter reports ASCII letters (drive-letter head of file:///C:/...).
func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// absHint is a tiny path helper for the launch resolvers: it only decides
// "is this spelled like a direct file target", never a filesystem DECISION
// (risk.PathResolver owns those; D22 bans making decisions with
// filepath.Clean/Abs outside it — the allowlist verdict stays InAllowlist).
func looksLikeDirectPath(raw string) bool {
	if len(raw) >= 2 && raw[1] == ':' && (raw[0] >= 'A' && raw[0] <= 'Z' || raw[0] >= 'a' && raw[0] <= 'z') {
		return true
	}
	if filepath.IsAbs(raw) || strings.HasPrefix(raw, `\\`) || strings.HasPrefix(raw, "./") || strings.HasPrefix(raw, "../") {
		return true
	}
	return false
}
