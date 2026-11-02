package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// web_matrix_test.go — AC row 3: "Protocol matrix: web.open/file.open
// allow/deny table per 16.5.3 — exhaustive test." The rows below are the
// ticket's own list, and the SHAPE of the list (allowlist + named explicit
// denies + default-deny tail) is what each direction pins.

// msRoster is the injected "known ms-*" registry face: what the production
// reader (HKCR) would answer on a normal Windows box.
func msRoster() func(string) bool {
	return func(scheme string) bool {
		switch scheme {
		case "ms-settings", "ms-calculator", "ms-paint", "ms-clock",
			"search-ms", "mshta": // last two: REGISTERED on Windows and still
			// denied — registration is not authorization (16.5.3).
			return true
		}
		return false
	}
}

type matrixRow struct {
	name           string
	target         string
	wantOpen       bool   // allowed by web.open
	wantFile       bool   // allowed by file.open
	wantFileLaunch string // when allowed: expected launch kind
}

var matrixRows = []matrixRow{
	// --- web.open: ONLY http/https (D30(3) strict side)
	{"https", "https://ok.test/a", true, false, ""},
	{"http", "http://ok.test/a", true, false, ""},
	{"file via web.open", "file:///C:/x", false, false, ""},
	{"mailto via web.open", "mailto:a@b.test", false, true, "uri"},
	{"ms-settings via web.open", "ms-settings:display", false, true, "uri"},
	{"search-ms registered-but-denied", "search-ms:query=x", false, false, ""},
	{"mshta registered-but-denied", "mshta:javascript.run", false, false, ""},
	{"javascript", "javascript:alert(1)", false, false, ""},
	{"vbscript", "vbscript:msgbox(1)", false, false, ""},
	{"shell", "shell:startup", false, false, ""},
	{"tel", "tel:5551234", false, false, ""},
	{"ftp", "ftp://host.test/x", false, false, ""},
	{"data", "data:text/html,x", false, false, ""},
	{"unknown scheme", "whatev:thing", false, false, ""},
	{"no scheme", "plain text no colon", false, false, ""},
	// --- file.open allow side
	{"registered ms- app protocol", "ms-calculator:", false, true, "uri"},
	{"unregistered ms- protocol", "ms-notreal:x", false, false, ""},
}

func Test22ProtocolMatrix(t *testing.T) {
	allowed := t.TempDir()
	rs := newRecShared(t, func(x *WebDeps) {
		x.IsKnownMsProtocol = msRoster()
		x.Paths = NewPathCanonicalizer([]string{allowed}, nil)
	})

	for _, tc := range matrixRows {
		t.Run(tc.name, func(t *testing.T) {
			rs.launches = nil

			// web.open door:
			res := execCall(t, rs.openTool(), map[string]any{"url": tc.target})
			if got := !res.IsError; got != tc.wantOpen {
				t.Errorf("web.open(%q): allowed=%v want=%v (text: %s)", tc.target, got, tc.wantOpen, res.Text)
			}
			if tc.wantOpen {
				if len(rs.launches) != 1 || rs.launches[0] != "url|"+tc.target {
					t.Errorf("web.open launches = %v, want [url|%s]", rs.launches, tc.target)
				}
			} else if len(rs.launches) != 0 {
				t.Errorf("web.open refused %q yet still launched: %v", tc.target, rs.launches)
			}

			// file.open door (fresh recorder each row so counts are clean):
			rs.launches = nil
			res = execCall(t, rs.fileTool(), map[string]any{"path": tc.target})
			got := !res.IsError
			if got != tc.wantFile {
				t.Errorf("file.open(%q): allowed=%v want=%v (text: %s)", tc.target, got, tc.wantFile, res.Text)
			}
			if tc.wantFile {
				if len(rs.launches) != 1 {
					t.Fatalf("file.open(%q) launches = %v, want exactly 1", tc.target, rs.launches)
				}
				if !strings.HasPrefix(rs.launches[0], tc.wantFileLaunch+"|") {
					t.Errorf("file.open(%q) launch kind = %q want %q", tc.target, rs.launches[0], tc.wantFileLaunch)
				}
			} else if len(rs.launches) != 0 {
				t.Errorf("file.open(%q) DENIED target still launched: %v", tc.target, rs.launches)
			}
		})
	}
}

// Test22FileOpenPathArms pins the file:-inside-allowlist rule through the
// REAL C26 seam (not a fake): inside path opens, outside path refuses,
// remote file://host refuses, and the denial never reaches the launcher.
func Test22FileOpenPathArms(t *testing.T) {
	allowed := t.TempDir()
	inside := filepath.Join(allowed, "note.txt")
	if err := os.WriteFile(inside, []byte("hi"), 0o600); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(outside, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	rs := newRecShared(t, func(x *WebDeps) {
		x.Paths = NewPathCanonicalizer([]string{allowed}, nil)
	})
	f := rs.fileTool()

	rs.launches = nil
	res := execCall(t, f, map[string]any{"path": inside})
	if res.IsError {
		t.Errorf("in-allowlist path refused: %s", res.Text)
	}
	if len(rs.launches) != 1 || !strings.HasPrefix(rs.launches[0], "file|") {
		t.Errorf("launches = %v", rs.launches)
	}

	rs.launches = nil
	res = execCall(t, f, map[string]any{"path": outside})
	if !res.IsError || !strings.Contains(res.Text, "allowed_dirs") {
		t.Errorf("out-of-allowlist path must refuse, got: %s", res.Text)
	}
	if len(rs.launches) != 0 {
		t.Errorf("refused path still launched: %v", rs.launches)
	}

	rs.launches = nil
	fileURL := "file:///" + strings.ReplaceAll(filepath.ToSlash(inside), "//", "/")
	res = execCall(t, f, map[string]any{"path": fileURL})
	if res.IsError {
		t.Errorf("file: URL inside allowlist refused: %s (%s)", res.Text, fileURL)
	}

	rs.launches = nil
	res = execCall(t, f, map[string]any{"path": "file://otherhost/share/x"})
	if !res.IsError || !strings.Contains(res.Text, "本机") {
		t.Errorf("remote file:// host must refuse: %s", res.Text)
	}
	if len(rs.launches) != 0 {
		t.Errorf("refused file: still launched: %v", rs.launches)
	}
}

// Test22DenyBeforeAllowOrder is the mutation target. The six schemes of
// 16.5.3 are DENIED-BY-NAME: the refusal must quote the plan's own reason
// (CLSID 劫持 / 脚本宿主 / 脚本协议 / shell 命名空间), not just any
// "拒绝" — a default-deny tail that swallows them into "不在白名单" is
// weaker than the contract (registration-vs-deny order unproven, and the
// user card loses the WHY). If the code ever checks registration BEFORE
// the explicit-deny list, search-ms:/mshta: (both genuinely registered on
// Windows) flip to allow and this goes red.
func Test22DenyBeforeAllowOrder(t *testing.T) {
	rs := newRecShared(t, func(x *WebDeps) { x.IsKnownMsProtocol = msRoster() })
	named := []struct{ target, wantReason string }{
		{"search-ms:query=x", "CLSID"},
		{"mshta:x", "脚本宿主"},
		{"javascript:alert(1)", "脚本协议"},
		{"vbscript:x", "脚本协议"},
		{"shell:x", "命名空间"},
		{"tel:1", "明确拒绝"},
	}
	for _, tc := range named {
		res := execCall(t, rs.fileTool(), map[string]any{"path": tc.target})
		if !res.IsError || !strings.Contains(res.Text, tc.wantReason) {
			t.Errorf("%s must land on the NAMED denial (reason %q) even when OS-registered: %s", tc.target, tc.wantReason, res.Text)
		}
	}
	// web.open side of the same order (R5's own map lives in risk; the tool
	// door re-checks shape):
	for _, s := range []string{"javascript:alert(1)", "vbscript:x", "search-ms:q", "shell:x", "tel:1", "mshta:x"} {
		res := execCall(t, rs.openTool(), map[string]any{"url": s})
		if !res.IsError {
			t.Errorf("web.open allowed %s", s)
		}
	}
}

// Test22UnknownSchemeDefaultsToDeny is the "默认拒绝 + 白名单放行，不是
// 黑名单" sentence from 16.5.3 as a test: every scheme NOT on the allow
// side — however novel — is refused. The exhaustive list is impossible, so
// the pin is structural: the classifier has no tail fallthrough that
// allows, asserted here across a spread of invented shapes.
func Test22UnknownSchemeDefaultsToDeny(t *testing.T) {
	rs := newRecShared(t, func(x *WebDeps) { x.IsKnownMsProtocol = func(string) bool { return false } })
	inventions := []string{
		"zeta:thing", "x-ms-evil:y", "mS-Settings-notreal:z", "http2://a.test/",
		"ms-", "ledger:x", "filex:///C:/x", "maileto:a@b", "mssetting:x",
	}
	for _, s := range inventions {
		res := execCall(t, rs.fileTool(), map[string]any{"path": s})
		if !res.IsError {
			t.Errorf("invented scheme %q was ALLOWED by file.open: %s", s, res.Text)
		}
	}
}

// Test22FileOpenNoFirstVisitPrompt: D30(5) is about NEW DOMAINS via web
// tools; local opens must not spam domain prompts.
func Test22FileOpenNoFirstVisitPrompt(t *testing.T) {
	allowed := t.TempDir()
	p := filepath.Join(allowed, "a.txt")
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	rs := newRecShared(t, func(x *WebDeps) { x.Paths = NewPathCanonicalizer([]string{allowed}, nil) })
	execCall(t, rs.fileTool(), map[string]any{"path": p})
	if len(rs.asks) != 0 {
		t.Errorf("file.open asked a domain prompt: %v", rs.asks)
	}
}
