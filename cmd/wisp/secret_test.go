package main

// Tests for `wisp secret` (ticket 63). The point of this ticket is not that
// the command works, it is that the command cannot leak: these tests assert on
// the log surface, the two output streams, the error strings, the flag grammar
// and (in secret_argv_windows_test.go) the command line the OS can see.
//
// Every key value used here is an obviously fake placeholder (D36#5, R7): no
// real credential exists anywhere in this file, so a leak assertion failing is
// always a code problem, never a leaked secret.

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/secret"
)

// fakeKey is the fixed probe value. If any assertion below fails, this string
// is what escaped; it is not a credential.
const fakeKey = "sk-FAKE-placeholder-not-a-real-key-0f9e8d7c"

// captureLogs installs an in-memory slog handler for the duration of the test
// and returns the buffer every record lands in. JSON serialization is the same
// shape the production pipeline (observe.InitLog) writes, so an attribute
// cannot hide from the assertion by being a non-string value.
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return &buf
}

// hiddenRecorder stands in for golang.org/x/term ReadPassword: it records the
// prompts it was given (so a prompt bug that prints the value is caught) and
// returns queued answers.
type hiddenRecorder struct {
	answers []string
	prompts []string
	err     error
	i       int
}

func (h *hiddenRecorder) read(prompt string) (string, error) {
	h.prompts = append(h.prompts, prompt)
	if h.err != nil {
		return "", h.err
	}
	if h.i >= len(h.answers) {
		return "", fmt.Errorf("test: hidden input exhausted after %d reads", h.i)
	}
	v := h.answers[h.i]
	h.i++
	return v, nil
}

// probe is one `wisp secret` run wired to in-memory streams.
type probe struct {
	cmd  *secretCmd
	out  *bytes.Buffer
	errb *bytes.Buffer
	logs *bytes.Buffer
	hid  *hiddenRecorder
}

// newProbe builds a secretCmd bound to dataDir with every process surface
// injected. readHidden is nil unless answers are given, which is how the
// "no console" failure path is reproduced without a terminal.
func newProbe(t *testing.T, env buildinfo.Env, dataDir string, portable bool, stdin string, answers ...string) *probe {
	t.Helper()
	p := &probe{
		out:  &bytes.Buffer{},
		errb: &bytes.Buffer{},
		logs: captureLogs(t),
	}
	if len(answers) > 0 {
		p.hid = &hiddenRecorder{answers: answers}
	}
	p.cmd = &secretCmd{
		env:      env,
		dataDir:  dataDir,
		portable: portable,
		sio: secretIO{
			stdin:      strings.NewReader(stdin),
			stdout:     p.out,
			stderr:     p.errb,
			readHidden: p.hiddenReader(),
		},
	}
	return p
}

func (p *probe) hiddenReader() func(string) (string, error) {
	if p.hid == nil {
		return nil
	}
	return p.hid.read
}

// all surfaces a value could have leaked into: stdout, stderr, and the log.
func (p *probe) surfaces() map[string]string {
	return map[string]string{
		"stdout": p.out.String(),
		"stderr": p.errb.String(),
		"logs":   p.logs.String(),
	}
}

// assertNoPlaintext fails the test if want absent appears in any output
// surface, reporting which one and the offending line.
func (p *probe) assertNoPlaintext(t *testing.T, wantAbsent, ctx string) {
	t.Helper()
	for name, body := range p.surfaces() {
		for _, line := range strings.Split(body, "\n") {
			if strings.Contains(line, wantAbsent) {
				t.Errorf("%s: %q appears in %s: %s", ctx, wantAbsent, name, line)
			}
		}
	}
}

// blobPath is the on-disk location of a name in this probe's store.
func (p *probe) blobPath(name string) string {
	return filepath.Join(p.cmd.dataDir, "secrets", name)
}

func (p *probe) blobExists(t *testing.T, name string) bool {
	t.Helper()
	_, err := os.Stat(p.blobPath(name))
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		t.Fatalf("stat blob %s: %v", name, err)
	}
	return true
}

// ---------------------------------------------------------------------------
// AC1: set stores a DPAPI blob, get shows it masked
// ---------------------------------------------------------------------------

func TestSecretSetGetListRoundTrip(t *testing.T) {
	dir := t.TempDir()
	logs := captureLogs(t)

	// --from-stdin: the value arrives on the standard input, in one read.
	p := newProbe(t, buildinfo.EnvDev, dir, false, fakeKey+"\n")
	if code := p.cmd.run([]string{"set", "deepseek", "--from-stdin"}); code != 0 {
		t.Fatalf("set --from-stdin exit = %d, want 0 (stderr: %s)", code, p.errb)
	}
	if !p.blobExists(t, "deepseek") {
		t.Fatal("set --from-stdin wrote no blob file")
	}
	// The blob on disk is a DPAPI blob, not the plaintext.
	raw, err := os.ReadFile(p.blobPath("deepseek"))
	if err != nil {
		t.Fatalf("read blob: %v", err)
	}
	if bytes.Contains(raw, []byte(fakeKey)) {
		t.Error("blob file contains the plaintext secret")
	}
	if bytes.Contains(bytes.ToLower(raw), []byte("placeholder")) {
		t.Error("blob file contains a fragment of the secret")
	}
	out := p.out.String()
	if !strings.Contains(out, `api_key_ref = "dpapi:deepseek"`) {
		t.Errorf("set must print the config reference, got:\n%s", out)
	}
	if !strings.Contains(out, "WISP_ENV=dev") {
		t.Errorf("set must echo the active env (dev key into the prod store is the failure mode), got:\n%s", out)
	}
	if !strings.Contains(out, secret.RedactSecret(fakeKey)) {
		t.Errorf("set must echo the masked form for confirmation, got:\n%s", out)
	}
	if strings.Contains(out, fakeKey) {
		t.Error("set printed the plaintext secret to stdout")
	}
	if strings.Contains(p.errb.String(), fakeKey) {
		t.Error("set printed the plaintext secret to stderr")
	}
	if strings.Contains(logs.String(), fakeKey) {
		t.Errorf("set logged the plaintext secret:\n%s", logs.String())
	}

	// Same name through the hidden-input path (no --from-stdin).
	p2 := newProbe(t, buildinfo.EnvDev, dir, false, "", fakeKey+"-b", fakeKey+"-b")
	if code := p2.cmd.run([]string{"set", "stepfun"}); code != 0 {
		t.Fatalf("set interactive exit = %d, want 0 (stderr: %s)", code, p2.errb)
	}
	for _, pr := range p2.hid.prompts {
		if strings.Contains(pr, fakeKey) {
			t.Errorf("hidden-input prompt carried the value: %q", pr)
		}
		if !strings.Contains(pr, "hidden") {
			t.Errorf("hidden-input prompt must say the input is hidden: %q", pr)
		}
	}
	if len(p2.hid.prompts) != 2 {
		t.Errorf("set read %d hidden inputs, want exactly 2 (entry + confirmation)", len(p2.hid.prompts))
	}

	// get is masked by default, plaintext only under --show.
	pg := newProbe(t, buildinfo.EnvDev, dir, false, "")
	if code := pg.cmd.run([]string{"get", "deepseek"}); code != 0 {
		t.Fatalf("get exit = %d, want 0 (stderr: %s)", code, pg.errb)
	}
	if got := pg.out.String(); strings.Contains(got, fakeKey) || !strings.Contains(got, secret.RedactSecret(fakeKey)) {
		t.Errorf("get must mask by default, got:\n%s", got)
	}
	ps := newProbe(t, buildinfo.EnvDev, dir, false, "")
	if code := ps.cmd.run([]string{"get", "deepseek", "--show"}); code != 0 {
		t.Fatalf("get --show exit = %d, want 0 (stderr: %s)", code, ps.errb)
	}
	if ps.out.String() != fakeKey {
		t.Errorf("get --show stdout = %q, want the key verbatim with no trailing newline (clip-friendly)", ps.out.String())
	}
	// --show is the one sanctioned plaintext output; it must still not log.
	if strings.Contains(ps.logs.String(), fakeKey) {
		t.Errorf("--show must not log the plaintext:\n%s", ps.logs.String())
	}

	// list shows ids and times only.
	pl := newProbe(t, buildinfo.EnvDev, dir, false, "")
	if code := pl.cmd.run([]string{"list"}); code != 0 {
		t.Fatalf("list exit = %d, want 0 (stderr: %s)", code, pl.errb)
	}
	listing := pl.out.String()
	for _, want := range []string{"deepseek", "stepfun", "dpapi:deepseek", "WISP_ENV=dev"} {
		if !strings.Contains(listing, want) {
			t.Errorf("list must show %q, got:\n%s", want, listing)
		}
	}
	if strings.Contains(listing, fakeKey) || strings.Contains(listing, "placeholder") {
		t.Errorf("list must never show secret content, got:\n%s", listing)
	}
	// Each listed row must carry a UTC creation timestamp (AC: id + 时间).
	rows := 0
	for _, line := range strings.Split(listing, "\n") {
		if !strings.Contains(line, "dpapi:") {
			continue
		}
		rows++
		if !strings.Contains(line, "T") || !strings.Contains(line, "Z") {
			t.Errorf("list row has no UTC timestamp: %q", line)
		}
		if strings.Contains(line, fakeKey) {
			t.Errorf("list row carries secret content: %q", line)
		}
	}
	if rows != 2 {
		t.Errorf("list showed %d blob rows, want 2:\n%s", rows, listing)
	}
}

func TestSecretSetRejectsBadNamesAndEmptyInput(t *testing.T) {
	dir := t.TempDir()

	cases := []struct {
		name  string
		args  []string
		stdin string
	}{
		{"traversal", []string{"set", "../escape", "--from-stdin"}, fakeKey},
		{"dotdot", []string{"set", "..", "--from-stdin"}, fakeKey},
		{"separator", []string{"set", `a\b`, "--from-stdin"}, fakeKey},
		{"space", []string{"set", "a b", "--from-stdin"}, fakeKey},
		{"empty", []string{"set", "", "--from-stdin"}, fakeKey},
		{"empty stdin", []string{"set", "ok", "--from-stdin"}, "\r\n"},
		{"multiline stdin", []string{"set", "ok", "--from-stdin"}, fakeKey + "\nsecond-line\n"},
		{"no name", []string{"set", "--from-stdin"}, fakeKey},
		{"two names", []string{"set", "a", "b", "--from-stdin"}, fakeKey},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logs := captureLogs(t)
			p := newProbe(t, buildinfo.EnvTest, dir, false, tc.stdin)
			code := p.cmd.run(tc.args)
			if code == 0 {
				t.Fatalf("run(%v) exit = 0, want a refusal", tc.args)
			}
			if strings.Contains(p.out.String(), fakeKey) || strings.Contains(p.errb.String(), fakeKey) {
				t.Errorf("refusal must not echo the value:\nstdout=%s\nstderr=%s", p.out, p.errb)
			}
			if strings.Contains(logs.String(), fakeKey) {
				t.Errorf("refusal must not log the value:\n%s", logs.String())
			}
			for _, stray := range []string{"escape", "a b", "ok"} {
				if p.blobExists(t, stray) {
					t.Errorf("rejected set left a blob for %q", stray)
				}
			}
		})
	}
	// The traversal attempts must not have created anything outside secrets\.
	if _, err := os.Stat(filepath.Join(filepath.Dir(dir), "escape")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a name traversal escaped the secrets dir (stat err = %v)", err)
	}
}

func TestSecretSetConfirmationMismatchStoresNothing(t *testing.T) {
	dir := t.TempDir()
	logs := captureLogs(t)
	p := newProbe(t, buildinfo.EnvTest, dir, false, "", fakeKey, fakeKey+"-typo")
	if code := p.cmd.run([]string{"set", "mismatch"}); code == 0 {
		t.Fatal("mismatched confirmation must fail")
	}
	if p.blobExists(t, "mismatch") {
		t.Error("mismatched confirmation wrote a blob")
	}
	msg := p.errb.String()
	if !strings.Contains(msg, "did not match") {
		t.Errorf("error must say the entries did not match, got: %s", msg)
	}
	p.assertNoPlaintext(t, fakeKey, "confirmation mismatch")
	if strings.Contains(logs.String(), fakeKey[:8]) {
		t.Errorf("no part of either entry may be logged:\n%s", logs.String())
	}
}

func TestSecretSetWithoutConsolePointsAtFromStdin(t *testing.T) {
	dir := t.TempDir()
	p := newProbe(t, buildinfo.EnvTest, dir, false, "") // readHidden == nil
	if code := p.cmd.run([]string{"set", "nokey"}); code == 0 {
		t.Fatal("set with no console and no --from-stdin must fail")
	}
	msg := p.errb.String()
	if !strings.Contains(msg, "--from-stdin") {
		t.Errorf("the failure must name the usable channel, got: %s", msg)
	}
	if p.blobExists(t, "nokey") {
		t.Error("failed set wrote a blob")
	}
	p.assertNoPlaintext(t, fakeKey, "no console")
}

// ---------------------------------------------------------------------------
// No argv slot can carry a secret: the flag grammar of this command is bool
// only. Parsed from the source so a future -value=<key> flag fails the gate.
// ---------------------------------------------------------------------------

func TestSecretFlagsAreBoolOnly(t *testing.T) {
	src, err := os.ReadFile("secret.go")
	if err != nil {
		t.Fatalf("read secret.go: %v", err)
	}
	for _, banned := range []string{
		"fs.String(", "fs.StringVar(", "fs.Int(", "fs.Int64(", "fs.Uint(",
		"fs.Float64(", "fs.Duration(", "fs.Func(", "fs.Var(", "fs.TextVar(",
		"flag.String(", "flag.Var(", "flag.Func(",
	} {
		if i := bytes.Index(src, []byte(banned)); i >= 0 {
			line := 1 + bytes.Count(src[:i], []byte("\n"))
			t.Errorf("secret.go:%d defines the value-taking flag form %q: `wisp secret` must expose "+
				"boolean flags only, so no credential can ever arrive as an argv value", line, banned)
		}
	}
	// The three flags this ticket does expose.
	for _, want := range []string{`fs.Bool("from-stdin"`, `fs.Bool("show"`, `fs.Bool("force"`} {
		if !bytes.Contains(src, []byte(want)) {
			t.Errorf("secret.go must define %s", want)
		}
	}
}

func TestSecretValueCarryingFlagsAreRefusedAndUnechoed(t *testing.T) {
	dir := t.TempDir()
	// Every plausible way an operator (or a wrapper script) might try to hand
	// the value over on the command line must be rejected before the secret is
	// read, and must not be reflected in the error text.
	forms := []string{"--value", "--secret", "--key", "-k", "--api-key", "--from-file", "--file", "--stdin"}
	for _, form := range forms {
		t.Run(form, func(t *testing.T) {
			logs := captureLogs(t)
			p := newProbe(t, buildinfo.EnvTest, dir, false, "", fakeKey, fakeKey)
			code := p.cmd.run([]string{"set", "probe", form + "=" + fakeKey})
			if code == 0 {
				t.Fatalf("%s=... was accepted; a value-taking flag is an argv leak", form)
			}
			if p.blobExists(t, "probe") {
				t.Errorf("%s=... stored a blob: argv would carry the secret", form)
			}
			for name, body := range map[string]string{"stderr": p.errb.String(), "stdout": p.out.String(), "logs": logs.String()} {
				if strings.Contains(body, fakeKey) {
					t.Errorf("%s=... leaked into %s (flag errors must name the flag, not its value)", form, name)
				}
			}
			if !strings.Contains(p.errb.String(), form) && !strings.Contains(p.errb.String(), strings.TrimPrefix(form, "-")) {
				t.Errorf("error must name the rejected flag %q, got: %s", form, p.errb)
			}
		})
	}
	// Positional overflow is refused the same way: the extra argument is a
	// name slot, not a value slot, and must not be stored or echoed as a key.
	t.Run("positional", func(t *testing.T) {
		p := newProbe(t, buildinfo.EnvTest, dir, false, "", fakeKey, fakeKey)
		if code := p.cmd.run([]string{"set", "probe", fakeKey}); code == 0 {
			t.Fatal("an extra positional was accepted")
		} else if strings.Contains(p.errb.String(), fakeKey) {
			t.Error("the positional-arity error echoed the extra argument")
		}
	})
}

// ---------------------------------------------------------------------------
// AC2: no plaintext in any log line or error string, including failures
// ---------------------------------------------------------------------------

func TestSecretFailurePathsLogAndPrintNoPlaintext(t *testing.T) {
	const goodKey = "sk-FAKE-intact-before-failure-abcdef1234"

	t.Run("bad dpapi decrypt (corrupted blob, non-portable)", func(t *testing.T) {
		dir := t.TempDir()
		p := newProbe(t, buildinfo.EnvTest, dir, false, "", goodKey, goodKey)
		if code := p.cmd.run([]string{"set", "corrupt"}); code != 0 {
			t.Fatalf("set: exit %d (%s)", code, p.errb)
		}
		if err := os.WriteFile(p.blobPath("corrupt"), []byte("not a DPAPI blob at all"), 0o600); err != nil {
			t.Fatalf("corrupt: %v", err)
		}
		p2 := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := p2.cmd.run([]string{"get", "corrupt"}); code == 0 {
			t.Fatal("get on a corrupted blob must fail")
		}
		p2.assertNoPlaintext(t, goodKey, "corrupted blob get")
		if !strings.Contains(p2.errb.String(), "dpapi:corrupt") {
			t.Errorf("the failure must identify the ref, got stderr=%s logs=%s", p2.errb, p2.logs)
		}
	})

	t.Run("bad dpapi decrypt (portable, P13)", func(t *testing.T) {
		dir := t.TempDir()
		p := newProbe(t, buildinfo.EnvTest, dir, true, "", goodKey, goodKey)
		if code := p.cmd.run([]string{"set", "portable"}); code != 0 {
			t.Fatalf("set in portable mode: exit %d (%s)", code, p.errb)
		}
		if err := os.WriteFile(p.blobPath("portable"), []byte("blob from another machine"), 0o600); err != nil {
			t.Fatalf("corrupt: %v", err)
		}
		p2 := newProbe(t, buildinfo.EnvTest, dir, true, "")
		if code := p2.cmd.run([]string{"get", "portable", "--show"}); code == 0 {
			t.Fatal("portable decrypt failure must surface")
		}
		p2.assertNoPlaintext(t, goodKey, "portable decrypt failure")
		if !strings.Contains(p2.errb.String(), "env:") {
			t.Errorf("P13 must guide to env: refs, got: %s", p2.errb)
		}
		if strings.Contains(p2.out.String(), goodKey) {
			t.Error("--show must print nothing when the decrypt failed")
		}
	})

	t.Run("unwritable store dir", func(t *testing.T) {
		dir := t.TempDir()
		// A regular file where secrets\ must go: NewStore's MkdirAll fails, so
		// the whole command must fail without touching the value.
		if err := os.WriteFile(filepath.Join(dir, "secrets"), []byte("blocker"), 0o600); err != nil {
			t.Fatalf("block secrets dir: %v", err)
		}
		p := newProbe(t, buildinfo.EnvTest, dir, false, "", goodKey, goodKey)
		if code := p.cmd.run([]string{"set", "unwritable"}); code == 0 {
			t.Fatal("set into an unwritable store must fail")
		}
		p.assertNoPlaintext(t, goodKey, "unwritable dir")
		for _, sub := range []string{"get", "list", "unset"} {
			q := newProbe(t, buildinfo.EnvTest, dir, false, "")
			if code := q.cmd.run([]string{sub, "unwritable"}); code == 0 {
				t.Fatalf("%s against an unwritable store must fail", sub)
			}
			q.assertNoPlaintext(t, goodKey, "unwritable dir: "+sub)
		}
	})

	t.Run("unset name (no such blob)", func(t *testing.T) {
		dir := t.TempDir()
		p := newProbe(t, buildinfo.EnvTest, dir, false, "", goodKey, goodKey)
		if code := p.cmd.run([]string{"get", "never-stored"}); code == 0 {
			t.Fatal("get of an unset name must fail")
		}
		p.assertNoPlaintext(t, goodKey, "unset name get")
		q := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q.cmd.run([]string{"unset", "never-stored"}); code == 0 {
			t.Fatal("unset of a name that was never stored must fail")
		}
		q.assertNoPlaintext(t, goodKey, "unset name unset")
		if !strings.Contains(q.errb.String(), "never-stored") {
			t.Errorf("the failure must name the ref, got: %s", q.errb)
		}
	})

	t.Run("store write failure surfaces the ref only", func(t *testing.T) {
		dir := t.TempDir()
		// Force Store to fail after Exists succeeded: the blob directory is
		// replaced by a file between the two calls.
		p := newProbe(t, buildinfo.EnvTest, dir, false, "", goodKey, goodKey)
		if code := p.cmd.run([]string{"set", "late"}); code != 0 {
			t.Fatalf("set: exit %d (%s)", code, p.errb)
		}
		blob := p.blobPath("late")
		if err := os.Remove(blob); err != nil {
			t.Fatalf("remove blob: %v", err)
		}
		if err := os.Mkdir(blob, 0o700); err != nil {
			t.Fatalf("occupy blob path with a dir: %v", err)
		}
		q := newProbe(t, buildinfo.EnvTest, dir, false, fakeKey+"\n")
		if code := q.cmd.run([]string{"set", "late", "--from-stdin"}); code == 0 {
			t.Fatal("Store over a directory must fail")
		}
		q.assertNoPlaintext(t, fakeKey, "store write failure")
		if !strings.Contains(q.errb.String(), "dpapi:late") {
			t.Errorf("the write failure must name the ref, got: %s", q.errb)
		}
	})
}

// allLogs is the whole log surface of the probe.
func (p *probe) allLogs() string { return p.logs.String() }

// ---------------------------------------------------------------------------
// AC3: unset refuses while config references the blob; --force audits
// ---------------------------------------------------------------------------

func TestSecretUnsetRefusesWhileReferenced(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	writeConfig := func(body string) {
		t.Helper()
		if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}
	}
	writeConfig(`schema_version = 2

[llm.providers.deepseek]
api_key_ref = "dpapi:deepseek"

[llm.providers.stepfun]
api_key_ref = "dpapi:stepfun"
`)

	p := newProbe(t, buildinfo.EnvTest, dir, false, fakeKey+"\n")
	if code := p.cmd.run([]string{"set", "deepseek", "--from-stdin"}); code != 0 {
		t.Fatalf("set: exit %d (%s)", code, p.errb)
	}

	t.Run("refused while referenced", func(t *testing.T) {
		q := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q.cmd.run([]string{"unset", "deepseek"}); code == 0 {
			t.Fatal("unset of a referenced blob must be refused")
		}
		if !q.blobExists(t, "deepseek") {
			t.Error("a refused unset must leave the blob in place")
		}
		msg := q.errb.String()
		if !strings.Contains(msg, "llm.providers.deepseek.api_key_ref") {
			t.Errorf("the refusal must name the referencing field, got:\n%s", msg)
		}
		if strings.Contains(msg, "llm.providers.stepfun.api_key_ref") {
			t.Errorf("the refusal must name only the fields that reference this blob, got:\n%s", msg)
		}
		if !strings.Contains(msg, "--force") {
			t.Errorf("the refusal must point at --force, got:\n%s", msg)
		}
	})

	t.Run("refused for a voice.realtime reference too", func(t *testing.T) {
		writeConfig(`schema_version = 2

[voice.realtime]
enabled = false
api_key_ref = "dpapi:deepseek"
`)
		q := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q.cmd.run([]string{"unset", "deepseek"}); code == 0 {
			t.Fatal("unset must also respect voice.realtime.api_key_ref")
		}
		if !strings.Contains(q.errb.String(), "voice.realtime.api_key_ref") {
			t.Errorf("refusal must name voice.realtime.api_key_ref, got:\n%s", q.errb)
		}
	})

	t.Run("an unreadable config fails closed", func(t *testing.T) {
		writeConfig("this is = = not toml")
		q := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q.cmd.run([]string{"unset", "deepseek"}); code == 0 {
			t.Fatal("unset must fail closed when the reference index is unreadable")
		}
		if !q.blobExists(t, "deepseek") {
			t.Error("a failed reference check must not delete the blob")
		}
	})

	t.Run("force deletes and writes an audit line", func(t *testing.T) {
		writeConfig(`schema_version = 2

[llm.providers.deepseek]
api_key_ref = "dpapi:deepseek"
`)
		q := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q.cmd.run([]string{"unset", "deepseek", "--force"}); code != 0 {
			t.Fatalf("unset --force: exit %d (%s)", code, q.errb)
		}
		if q.blobExists(t, "deepseek") {
			t.Error("--force must delete the blob")
		}
		audit := q.allLogs()
		if !strings.Contains(audit, "dpapi:deepseek") || !strings.Contains(audit, "llm.providers.deepseek.api_key_ref") {
			t.Errorf("the audit line must record the ref and the dangling field:\n%s", audit)
		}
		if !strings.Contains(audit, "forced=true") {
			t.Errorf("the audit line must record that --force was used:\n%s", audit)
		}
		if !strings.Contains(audit, "WARN") {
			t.Errorf("a forced delete of a referenced blob must be a Warn record:\n%s", audit)
		}
		if !strings.Contains(q.out.String(), "audit:") {
			t.Errorf("the audit line must also be visible on stdout:\n%s", q.out)
		}
		assertNoPlaintextIn(t, audit, fakeKey, "force-unset audit")
	})

	t.Run("unreferenced delete is audited at info", func(t *testing.T) {
		writeConfig(`schema_version = 2
`)
		q := newProbe(t, buildinfo.EnvTest, dir, false, fakeKey+"\n")
		if code := q.cmd.run([]string{"set", "spare", "--from-stdin"}); code != 0 {
			t.Fatalf("set spare: exit %d (%s)", code, q.errb)
		}
		q2 := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := q2.cmd.run([]string{"unset", "spare"}); code != 0 {
			t.Fatalf("unset spare: exit %d (%s)", code, q2.errb)
		}
		if q2.blobExists(t, "spare") {
			t.Error("unset of an unreferenced blob must delete it")
		}
		if !strings.Contains(q2.allLogs(), "dpapi:spare") {
			t.Errorf("unreferenced deletes are audited too:\n%s", q2.allLogs())
		}
		if strings.Contains(q2.allLogs(), "forced=true") {
			t.Errorf("forced=true only belongs to a forced delete:\n%s", q2.allLogs())
		}
	})

	t.Run("no config at all means no references", func(t *testing.T) {
		fresh := t.TempDir()
		q := newProbe(t, buildinfo.EnvTest, fresh, false, fakeKey+"\n")
		if code := q.cmd.run([]string{"set", "orphan", "--from-stdin"}); code != 0 {
			t.Fatalf("set: exit %d (%s)", code, q.errb)
		}
		q2 := newProbe(t, buildinfo.EnvTest, fresh, false, "")
		if code := q2.cmd.run([]string{"unset", "orphan"}); code != 0 {
			t.Fatalf("unset without a config: exit %d (%s)", code, q2.errb)
		}
		if q2.blobExists(t, "orphan") {
			t.Error("blob survived unset")
		}
	})
}

// ---------------------------------------------------------------------------
// AC4: WISP_ENV isolation through the command surface
// ---------------------------------------------------------------------------

func TestSecretSameNameUnderThreeEnvsIsThreeBlobs(t *testing.T) {
	configRoot := t.TempDir()
	testData := filepath.Join(configRoot, "injected-test")
	t.Setenv("WISP_TEST_DATA_DIR", testData)

	envs := []buildinfo.Env{buildinfo.EnvDev, buildinfo.EnvTest, buildinfo.EnvProd}
	dirs := map[buildinfo.Env]string{}
	for _, env := range envs {
		// Same resolution order the real command uses (ticket 06's fork),
		// with the config root injected so nothing is written under %APPDATA%.
		layout, err := sessionLayout(env, configRoot, "")
		if err != nil {
			t.Fatalf("sessionLayout(%s): %v", env, err)
		}
		dirs[env] = layout.DataDir
	}
	if dirs[buildinfo.EnvDev] == dirs[buildinfo.EnvProd] ||
		dirs[buildinfo.EnvDev] == dirs[buildinfo.EnvTest] ||
		dirs[buildinfo.EnvProd] == dirs[buildinfo.EnvTest] {
		t.Fatalf("the three envs share a data dir: %v", dirs)
	}

	const name = "shared-key"
	keys := map[buildinfo.Env]string{}
	for i, env := range envs {
		keys[env] = fmt.Sprintf("sk-FAKE-%s-%08d", env, 1000+i*7)
	}

	for _, env := range envs {
		p := newProbe(t, env, dirs[env], false, keys[env]+"\n")
		if code := p.cmd.run([]string{"set", name, "--from-stdin"}); code != 0 {
			t.Fatalf("set under %s: exit %d (%s)", env, code, p.errb)
		}
		if !strings.Contains(p.out.String(), "WISP_ENV="+string(env)) {
			t.Errorf("set under %s must echo the env name:\n%s", env, p.out)
		}
		if !p.blobExists(t, name) {
			t.Fatalf("no blob under %s", dirs[env])
		}
	}

	// Three distinct files with three distinct ciphertexts.
	seen := map[string]bool{}
	for _, env := range envs {
		raw, err := os.ReadFile(filepath.Join(dirs[env], "secrets", name))
		if err != nil {
			t.Fatalf("read blob for %s: %v", env, err)
		}
		if bytes.Contains(raw, []byte(keys[env])) {
			t.Errorf("blob under %s holds the plaintext", env)
		}
		seen[string(raw)] = true
	}
	if len(seen) != 3 {
		t.Fatalf("the three envs produced %d distinct blobs, want 3", len(seen))
	}

	// Each env resolves exactly its own key and no other env's.
	for _, env := range envs {
		q := newProbe(t, env, dirs[env], false, "")
		if code := q.cmd.run([]string{"get", name, "--show"}); code != 0 {
			t.Fatalf("get --show under %s: exit %d (%s)", env, code, q.errb)
		}
		got := q.out.String()
		if got != keys[env] {
			t.Errorf("get --show under %s returned the wrong env's key", env)
		}
		for _, other := range envs {
			if other != env && got == keys[other] {
				t.Errorf("env %s returned env %s's key", env, other)
			}
		}
	}

	// A listing is per-env too.
	prod := newProbe(t, buildinfo.EnvProd, dirs[buildinfo.EnvProd], false, "")
	if code := prod.cmd.run([]string{"list"}); code != 0 {
		t.Fatalf("list: %d (%s)", code, prod.errb)
	}
	if !strings.Contains(prod.out.String(), name) || !strings.Contains(prod.out.String(), "WISP_ENV=prod") {
		t.Errorf("prod listing wrong:\n%s", prod.out)
	}
	if strings.Contains(prod.out.String(), "wisp-dev") {
		t.Errorf("prod listing reached into the dev dir:\n%s", prod.out)
	}
}

// ---------------------------------------------------------------------------
// AC5: portable and non-portable, through ticket 06's seam (WithPortable)
// ---------------------------------------------------------------------------

func TestSecretPortableModeUsesTicket06Seam(t *testing.T) {
	configRoot := t.TempDir()
	exeDir := t.TempDir()
	// The marker is ticket 06's portable trigger (SPEC-02 §6).
	if err := os.WriteFile(filepath.Join(exeDir, "portable.txt"), nil, 0o600); err != nil {
		t.Fatalf("write portable.txt: %v", err)
	}
	layout, err := sessionLayout(buildinfo.EnvProd, configRoot, exeDir)
	if err != nil {
		t.Fatalf("sessionLayout with portable.txt: %v", err)
	}
	if !layout.Portable {
		t.Fatal("portable.txt next to the exe must mark the session portable")
	}
	if layout.DataDir != filepath.Join(exeDir, "data") {
		t.Errorf("portable data dir = %q, want the exe-relative dir", layout.DataDir)
	}
	// Without the marker the same env resolves to the fork dir, non-portable.
	forked, err := sessionLayout(buildinfo.EnvProd, configRoot, exeDir+"-no-marker")
	if err != nil {
		t.Fatalf("sessionLayout without a marker: %v", err)
	}
	if forked.Portable {
		t.Error("a non-existent portable.txt must not mark the session portable")
	}

	// Portable mode still stores and reads its own blobs (P13 only forbids
	// decrypting a blob made elsewhere, never a portable install's own keys).
	p := newProbe(t, buildinfo.EnvProd, layout.DataDir, true, fakeKey+"\n")
	if code := p.cmd.run([]string{"set", "portable-own", "--from-stdin"}); code != 0 {
		t.Fatalf("portable set: exit %d (%s)", code, p.errb)
	}
	if !strings.Contains(p.out.String(), "portable=true") {
		t.Errorf("portable runs must say so:\n%s", p.out)
	}
	q := newProbe(t, buildinfo.EnvProd, layout.DataDir, true, "")
	if code := q.cmd.run([]string{"get", "portable-own", "--show"}); code != 0 {
		t.Fatalf("portable get: exit %d (%s)", code, q.errb)
	}
	if q.out.String() != fakeKey {
		t.Error("portable get returned the wrong value")
	}
	// And never degrades to a plaintext file: the blob is opaque.
	raw, err := os.ReadFile(p.blobPath("portable-own"))
	if err != nil {
		t.Fatalf("read portable blob: %v", err)
	}
	if bytes.Contains(raw, []byte(fakeKey)) {
		t.Error("portable mode wrote the secret in plaintext")
	}
}

// ---------------------------------------------------------------------------
// AC6: end to end - a stored secret is resolved through config at request time
// ---------------------------------------------------------------------------

func TestSecretEndToEndConfigRefResolvesAtRequestTime(t *testing.T) {
	dir := t.TempDir()
	logs := captureLogs(t)

	// 1. The owner enters a key with no chat UI, no argv, no temp file.
	p := newProbe(t, buildinfo.EnvTest, dir, false, fakeKey+"\n")
	if code := p.cmd.run([]string{"set", "mockllm", "--from-stdin"}); code != 0 {
		t.Fatalf("set: exit %d (%s)", code, p.errb)
	}
	ref := "dpapi:mockllm"

	// 2. config.toml carries only the reference (D36 rule 5).
	configPath := filepath.Join(dir, "config.toml")
	body := "schema_version = 2\n\n[llm.providers.acme]\nprotocol = \"openai-chat\"\nbase_url = \"http://127.0.0.1:1/v1\"\napi_key_ref = \"" + ref + "\"\n\n[llm.providers.acme.models.m1]\n"
	if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if strings.Contains(body, fakeKey) {
		t.Fatal("test setup error: the config body would carry the key")
	}

	// 3. A request-time load resolves the ref through ticket 06's store.
	st, err := secret.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	cfg, resolved, err := config.LoadFile(configPath, st)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if got := resolved.ProviderKeys["acme"]; got != fakeKey {
		// The mismatch is reported by length only: a test failure must not
		// itself become a plaintext leak in the test log.
		t.Errorf("resolved provider key does not match the stored secret (got %d chars, want %d)", len(got), len(fakeKey))
	}
	if ref2 := cfg.LLM.Providers["acme"].APIKeyRef; ref2 != ref {
		t.Errorf("config still holds the ref, got %q", ref2)
	}
	// The config file itself has no plaintext (and never will).
	onDisk, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("reread config: %v", err)
	}
	if bytes.Contains(onDisk, []byte(fakeKey)) {
		t.Error("config.toml contains the plaintext key")
	}

	// 4. The resolved key is what goes on the wire: a chat request against a
	// local stand-in for tools/mockllm (a separate main package, so this test
	// serves the protocol shape itself rather than importing it).
	var sawAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"chatcmpl-fake","object":"chat.completion","model":"m1","choices":[{"index":0,"message":{"role":"assistant","content":"pong"},"finish_reason":"stop"}]}`))
	}))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", strings.NewReader(`{"model":"m1","messages":[{"role":"user","content":"hi"}]}`))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+resolved.ProviderKeys["acme"])
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("mock completion status = %d, want 200", resp.StatusCode)
	}
	if sawAuth != "Bearer "+fakeKey {
		t.Errorf("the mock saw an Authorization header of %d chars, want the %d-char key stored by `wisp secret set`",
			len(sawAuth), len(fakeKey)+len("Bearer "))
	}
	if strings.Contains(logs.String(), fakeKey) {
		t.Errorf("the end-to-end run logged the plaintext:\n%s", logs.String())
	}
}

// ---------------------------------------------------------------------------
// dispatch hygiene
// ---------------------------------------------------------------------------

func TestSecretUsageAndUnknownSubcommand(t *testing.T) {
	dir := t.TempDir()
	for _, args := range [][]string{
		nil,
		{"nope"},
		{"set"},
		{"get"},
		{"list", "extra"},
		{"unset"},
	} {
		p := newProbe(t, buildinfo.EnvTest, dir, false, "")
		if code := p.cmd.run(args); code != 2 {
			t.Errorf("run(%v) exit = %d, want the usage code 2", args, code)
		}
		if !strings.Contains(p.errb.String(), "Usage") {
			t.Errorf("run(%v) must print usage, got: %s", args, p.errb)
		}
	}
	p := newProbe(t, buildinfo.EnvTest, dir, false, "")
	if code := p.cmd.run([]string{"help"}); code != 0 {
		t.Errorf("help exit = %d, want 0", code)
	}
	if !strings.Contains(p.errb.String(), "from-stdin") || !strings.Contains(p.errb.String(), "WISP_ENV=") {
		t.Errorf("help must document the channels and the active env:\n%s", p.errb)
	}
}

func TestSecretOverwriteIsAnnounced(t *testing.T) {
	dir := t.TempDir()
	a := "sk-FAKE-first-value-aaaaaaaaaaaa1111"
	b := "sk-FAKE-second-value-bbbbbbbbbbbb2222"
	if code := newProbe(t, buildinfo.EnvTest, dir, false, a+"\n").cmd.run([]string{"set", "dup", "--from-stdin"}); code != 0 {
		t.Fatalf("first set: %d", code)
	}
	p := newProbe(t, buildinfo.EnvTest, dir, false, b+"\n")
	if code := p.cmd.run([]string{"set", "dup", "--from-stdin"}); code != 0 {
		t.Fatalf("second set: %d (%s)", code, p.errb)
	}
	if !strings.Contains(p.out.String(), "overwritten") {
		t.Errorf("an overwrite must be announced, got:\n%s", p.out)
	}
	q := newProbe(t, buildinfo.EnvTest, dir, false, "")
	if code := q.cmd.run([]string{"get", "dup", "--show"}); code != 0 || q.out.String() != b {
		t.Fatalf("get after overwrite = (%q, %d), want the new value", q.out.String(), code)
	}
}

// assertNoPlaintextIn is the free-function form used where a probe is not in
// hand (e.g. a raw log buffer).
func assertNoPlaintextIn(t *testing.T, body, wantAbsent, ctx string) {
	t.Helper()
	if strings.Contains(body, wantAbsent) {
		t.Errorf("%s: plaintext found in the surface: %s", ctx, body)
	}
}
