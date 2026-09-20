//go:build windows

package risk

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// Red-team suite for the C26 PathResolver (SPEC-06 §4, SPEC-10 §6).
// Every bypass case creates REAL OS artifacts (mklink /J junctions, real
// 8.3 short names) inside t.TempDir(); no string mocks.

// tmpHome builds a fake "home" whose layout mirrors the A-list anchors, so
// classification can be exercised without touching the real user profile.
func tmpHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(home, "AppData", "Local"))
	mustWrite(t, filepath.Join(home, ".git-credentials"), "secret")
	mustWrite(t, filepath.Join(home, ".ssh", "id_testkey"), "secret")
	mustWrite(t, filepath.Join(home, ".aws", "credentials"), "secret")
	mustWrite(t, filepath.Join(home, ".kube", "config"), "secret")
	mustWrite(t, filepath.Join(home, "plain.txt"), "ok")
	return home
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// mkJunction creates a REAL junction via mklink /J (no admin needed).
func mkJunction(t *testing.T, link, target string) {
	t.Helper()
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
	if err != nil {
		t.Fatalf("mklink /J %s %s failed: %v: %s", link, target, err, out)
	}
}

// getShortPath returns the real 8.3 short spelling via GetShortPathNameW.
// Returns "" when the volume does not generate short names.
func getShortPath(t *testing.T, long string) string {
	t.Helper()
	lp, err := syscall.UTF16PtrFromString(long)
	if err != nil {
		return ""
	}
	n, err := syscall.GetShortPathName(lp, nil, 0)
	if err != nil || n == 0 {
		return ""
	}
	buf := make([]uint16, n)
	n, err = syscall.GetShortPathName(lp, &buf[0], n)
	if err != nil || n == 0 {
		return ""
	}
	s := syscall.UTF16ToString(buf[:n])
	if s == long || !strings.Contains(strings.ToLower(filepath.Base(s)), "~") {
		return ""
	}
	return s
}

func gateFor(t *testing.T, canonical string) PathDecision {
	t.Helper()
	return Gate(canonical, nil)
}

// Case 1 (red-team: junction): a REAL junction pointing into an A-list tree
// must be rejected by the resolver (default reparse deny), and its resolved
// target must classify ClassA regardless.
func TestPathResolverJunctionWindows(t *testing.T) {
	home := tmpHome(t)
	link := filepath.Join(home, "jnlink")
	mkJunction(t, link, filepath.Join(home, ".ssh"))

	viaLink := filepath.Join(link, "id_testkey")
	_, err := Resolve(viaLink, nil)
	if !errors.Is(err, ErrReparseDenied) {
		t.Fatalf("junction to A-list tree: want ErrReparseDenied, got %v", err)
	}

	res, rerr := Resolve(viaLink, []string{link})
	if rerr != nil {
		t.Fatalf("exempted junction should pass reparse gate: %v", rerr)
	}
	if got := Classify(res.Canonical); got != ClassA {
		t.Fatalf("canonical %q classified %v, want ClassA (defense in depth)", res.Canonical, got)
	}
	if d := gateFor(t, res.Canonical); d.Allow {
		t.Fatalf("A-list via exempted junction granted: %+v", d)
	}
}

// Case 2 (red-team: 8.3 short name): the short spelling of an A-list file
// must resolve to the long path and be denied. Skips with a note when the
// temp volume has 8.3 generation disabled (fsutil setshortname needs admin).
func TestPathResolverShortNameAListDenied(t *testing.T) {
	home := tmpHome(t)
	long := filepath.Join(home, ".git-credentials")
	short := getShortPath(t, long)
	if short == "" {
		t.Skip("volume has 8.3 short-name generation disabled (GetShortPathName returned the long name); legal workaround per ticket 18: covered by junction/case cases")
	}
	res, err := Resolve(short, nil)
	if err != nil {
		t.Fatalf("short name is not a reparse point, want clean resolve: %v", err)
	}
	if !strings.EqualFold(res.Canonical, long) {
		t.Fatalf("short name did not expand to long path: got %q want %q", res.Canonical, long)
	}
	if got := Classify(res.Canonical); got != ClassA {
		t.Fatalf("8.3 spelling of A-list file classified %v, want ClassA", got)
	}
	if d := gateFor(t, res.Canonical); d.Allow {
		t.Fatalf("8.3 spelling of A-list file granted: %+v", d)
	}
}

// Case 3 (red-team: UNC): UNC spellings of a local-drive A-list file
// (\\?\UNC\localhost\c$\ and \\localhost\c$\) must normalize to the drive
// path and be denied.
func TestPathResolverUNCAListDenied(t *testing.T) {
	home := tmpHome(t)
	long := filepath.Join(home, ".git-credentials")
	vol := filepath.VolumeName(long) // e.g. C:
	share := strings.ToLower(vol[:1]) + "$"
	rest := strings.TrimPrefix(long, vol)
	for _, spelling := range []string{
		`\\?\UNC\localhost\` + share + rest,
		`\\localhost\` + share + rest,
		`//localhost/` + share + rest,
	} {
		res, err := Resolve(spelling, nil)
		if err != nil {
			t.Fatalf("UNC spelling %q: unexpected reparse error: %v", spelling, err)
		}
		if !strings.EqualFold(res.Canonical, long) {
			t.Fatalf("UNC spelling %q normalized to %q, want %q", spelling, res.Canonical, long)
		}
		if d := gateFor(t, res.Canonical); d.Allow {
			t.Fatalf("UNC spelling of A-list file granted: %q -> %+v", spelling, d)
		}
	}
}

// Case 4 (red-team: \\?\ extended-length prefix): the A-list file spelled
// with the \\?\ prefix must be denied.
func TestPathResolverExtendedLengthPrefixAListDenied(t *testing.T) {
	home := tmpHome(t)
	long := filepath.Join(home, ".git-credentials")
	res, err := Resolve(`\\?\`+long, nil)
	if err != nil {
		t.Fatalf("\\\\?\\ prefix is not a reparse point: %v", err)
	}
	if !strings.EqualFold(res.Canonical, long) {
		t.Fatalf("\\\\?\\ spelling normalized to %q, want %q", res.Canonical, long)
	}
	if d := gateFor(t, res.Canonical); d.Allow {
		t.Fatalf("\\\\?\\ spelling of A-list file granted: %+v", d)
	}
}

// Case 5 (red-team: case and . / .. mixing): redundant spellings of an
// A-list path must all be denied.
func TestPathResolverCaseAndDotMixDenied(t *testing.T) {
	home := tmpHome(t)
	base := filepath.Join(home, ".git-credentials")
	cases := []string{
		strings.ToUpper(base),
		strings.ToLower(base),
		filepath.Join(home, ".", ".git-credentials"),
		filepath.Join(home, ".ssh", "..", ".git-credentials"),
		filepath.Join(home, "no-such-dir", "..", ".git-credentials"),
	}
	for _, c := range cases {
		res, err := Resolve(c, nil)
		if err != nil {
			t.Fatalf("case/dot spelling %q: unexpected error: %v", c, err)
		}
		if d := gateFor(t, res.Canonical); d.Allow {
			t.Fatalf("case/dot spelling of A-list file granted: %q -> %+v", c, d)
		}
	}
}

// Case 6 (reparse exceptions): an explicitly excepted junction passes the
// reparse gate for its OWN subtree only, and a sibling junction stays denied.
func TestReparseExceptionExplicitPathOnly(t *testing.T) {
	home := tmpHome(t)
	link := filepath.Join(home, "jnlink")
	mkJunction(t, link, filepath.Join(home, ".ssh"))
	other := filepath.Join(home, "jnother")
	mkJunction(t, other, filepath.Join(home, ".aws"))

	if _, err := Resolve(filepath.Join(link, "id_testkey"), []string{link}); err != nil {
		t.Fatalf("excepted junction must pass: %v", err)
	}
	if _, err := Resolve(filepath.Join(other, "credentials"), []string{link}); !errors.Is(err, ErrReparseDenied) {
		t.Fatalf("non-excepted sibling junction must stay denied, got %v", err)
	}
	// Exception must match the specific path, not a sibling: an evil-twin
	// REAL junction next to the excepted one stays denied.
	evil := filepath.Join(home, "jnlink-eviltwin")
	mkJunction(t, evil, filepath.Join(home, ".aws"))
	if _, err := Resolve(filepath.Join(evil, "credentials"), []string{link}); !errors.Is(err, ErrReparseDenied) {
		t.Fatalf("exception must not leak to sibling junctions, got %v", err)
	}
}

// Case 7 (A-list breadth): every A-list anchor is denied read+write, and no
// override map can unlock it.
func TestAListDenyAndUnoverridable(t *testing.T) {
	home := tmpHome(t)
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(home, "AppData", "Local"))
	ad := filepath.Join(home, "AppData", "Roaming")
	ld := filepath.Join(home, "AppData", "Local")
	cases := []string{
		filepath.Join(home, ".git-credentials"),
		filepath.Join(home, "repo", ".git", "config"),
		filepath.Join(ad, "wisp", "config.toml"),
		filepath.Join(home, ".ssh", "known_hosts"),
		filepath.Join(home, ".ssh", "deep", "id_rsa"),
		filepath.Join(home, ".aws", "credentials"),
		filepath.Join(home, ".kube", "config"),
		filepath.Join(ld, "Google", "Chrome", "User Data", "Default", "Login Data"),
		filepath.Join(ld, "Microsoft", "Edge", "User Data", "Default", "Cookies"),
		filepath.Join(ad, "Microsoft", "Protect", "S-1-5-21", "df9d23cd"),
		filepath.Join(ld, "Microsoft", "Credentials", "AB12CD34"),
	}
	for _, c := range cases {
		if got := Classify(c); got != ClassA {
			t.Fatalf("A-list miss: %q classified %v", c, got)
		}
		d := Gate(c, map[string]bool{strings.ToLower(c): true})
		if d.Allow || d.NeedL2 {
			t.Fatalf("A-list must be unoverridable: %q -> %+v", c, d)
		}
	}
}

// Case 8 (B-list): default deny -> L2; single-file override grants with a
// log entry; non-listed paths stay ClassNone.
func TestBListDefaultDenyAndOverride(t *testing.T) {
	home := tmpHome(t)
	names := []string{".env", ".env.local", "server.pem", "cert.p12", "cert.pfx",
		"id_rsa", "secrets.yaml", "prod-credentials.json"}
	for _, n := range names {
		p := filepath.Join(home, "proj", n)
		if got := Classify(p); got != ClassB {
			t.Fatalf("B-list miss: %q classified %v", p, got)
		}
		if d := Gate(p, nil); d.Allow || !d.NeedL2 {
			t.Fatalf("B-list default must be deny->L2: %q -> %+v", p, d)
		}
	}
	var logs []string
	Logf = func(f string, a ...any) { logs = append(logs, sprintf(f, a...)) }
	defer func() { Logf = nil }()
	p := filepath.Join(home, "proj", ".env.local")
	if d := Gate(p, map[string]bool{strings.ToLower(p): true}); !d.Allow {
		t.Fatalf("B-list single-file override must allow: %+v", d)
	}
	if len(logs) == 0 {
		t.Fatal("B-list override must emit a log entry")
	}
	if got := Classify(filepath.Join(home, "proj", "main.go")); got != ClassNone {
		t.Fatalf("non-listed file classified %v, want ClassNone", got)
	}
}

// Case 9: resolver idempotence and plain-path round trip.
func TestResolveIdempotent(t *testing.T) {
	home := tmpHome(t)
	p := filepath.Join(home, "plain.txt")
	r1, err := Resolve(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := Resolve(r1.Canonical, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Canonical != r2.Canonical {
		t.Fatalf("not idempotent: %q vs %q", r1.Canonical, r2.Canonical)
	}
	if !r1.Resolved {
		t.Fatal("existing file must resolve via handle (Resolved=true)")
	}
}

// Case 10: env-var and ~ expansion feed the pipeline.
func TestResolveExpansion(t *testing.T) {
	home := tmpHome(t)
	t.Setenv("WISP_TEST_HOME", home)
	for _, in := range []string{
		`%WISP_TEST_HOME%\.git-credentials`,
		strings.ReplaceAll(filepath.Join(home, ".git-credentials"), `\`, "/"),
	} {
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("expand %q: %v", in, err)
		}
		if d := gateFor(t, res.Canonical); d.Allow {
			t.Fatalf("expanded A-list path granted: %q -> %+v", in, d)
		}
	}
}

func sprintf(format string, a ...any) string { return fmt.Sprintf(format, a...) }
