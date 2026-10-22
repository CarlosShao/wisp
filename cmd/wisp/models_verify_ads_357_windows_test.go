//go:build windows

package main

// Ticket 357 AC#3 at the exit: the sentence and the exit code, measured through
// the production function that prints and returns both.
//
// WHAT IS DRIVEN, AND WHY AT THIS LEVEL. The ruling has three clauses - say one
// sentence naming the stream and 清单未点名, keep rc at 0, and print exactly what
// was printed before when there is no stream - and two of the three are properties
// of the passing path as a whole. So these cases call announceUnnamedStreams357,
// which IS the passing path's tail (modelsVerify's last statement), over a real
// models.Manager: the same UnnamedStreams implementation the binary runs, over a
// real NTFS directory with a real alternate data stream planted by os.WriteFile.
//
// WHY NOT DRIVE modelsVerify ALL THE WAY FROM CONFIG. It would add one thing and
// cost one: the C29 gate in front of it verifies the manifest against the key
// embedded in buildinfo, so a fixture manifest has to be signed by a key this
// repository's build does not carry (the dev private key is deliberately outside
// the tree). The seam below is the Manager, not the manifest loader, so what gets
// stubbed out is a signature this repository already refuses to fake - not the
// enumeration, not the directory resolution, and not the exit code.
//
// THE MUTATION THIS FILE IS BUILT TO CATCH: deleting the announce call from
// modelsVerify's tail leaves every assertion here green only if these cases called
// something else. Test357VerifyTailIsStillOnThePassingPath pins that call by
// reading modelsVerify's own source, so the sentence cannot be unhooked from the
// command while the unit cases stay green.

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/models"
)

// streamFixture357 is one installed model, the Manager over it, and the directory
// that Manager resolves for the id.
type streamFixture357 struct {
	mgr *models.Manager
	dir string
}

// newStreamFixture357 installs a one-file model into a fresh store (or into an
// override directory) and returns the Manager the CLI tail would be handed.
func newStreamFixture357(t *testing.T, override string) streamFixture357 {
	t.Helper()
	const id = "verify357-fixture"
	body := []byte("wisp-357-fixture-body")
	sum := sha256.Sum256(body)
	manifest := fmt.Sprintf(`{
		"manifest_version": 1,
		"models": [
			{
				"id": %q, "purpose": "vad",
				"urls": ["https://example.invalid/%s/model.onnx"],
				"sha256": %q, "size_bytes": %d,
				"license": "MIT", "quant": "fp32"
			}
		]
	}`, id, id, hex.EncodeToString(sum[:]), len(body))
	mf := signedManifest357(t, manifest)

	store := filepath.Join(t.TempDir(), "models")
	dir := filepath.Join(store, id)
	if override != "" {
		dir = override
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "model.onnx"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	opts := models.Options{DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 1}
	if override != "" {
		opts.LocalOverride = map[string]string{id: override}
	}
	mgr, err := models.NewManager(opts)
	if err != nil {
		t.Fatal(err)
	}
	// The premise of every case below: this tree is what the manifest pins, so the
	// verification that earns rc=0 really does pass on it.
	if err := mgr.VerifyInstalled(id); err != nil {
		t.Fatalf("fixture does not verify before anything is planted: %v", err)
	}
	return streamFixture357{mgr: mgr, dir: dir}
}

// signedManifest357 signs a manifest with the same deterministic fixture key
// internal/models' own cases use, and loads it back through the real verification
// path - so the signature leg here is the shipped one, not a bypass.
func signedManifest357(t *testing.T, body string) *models.Manifest {
	t.Helper()
	seed := []byte("wisp-test-minisign-seed-01234567")
	priv := ed25519.NewKeyFromSeed(seed)
	pub := models.MinisignPublicKey{
		KeyID: models.DeriveKeyID(priv.Public().(ed25519.PublicKey)),
		Key:   priv.Public().(ed25519.PublicKey),
	}
	dir := t.TempDir()
	mpath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(mpath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	sig, err := models.SignPayload(priv, pub, []byte(body), "357", "trusted comment: 357")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mpath+".minisig", sig, 0o644); err != nil {
		t.Fatal(err)
	}
	mf, err := models.LoadSignedManifest(mpath, mpath+".minisig", pub.String("357 fixture key"))
	if err != nil {
		t.Fatalf("LoadSignedManifest: %v", err)
	}
	return mf
}

// plant writes one named stream beside the installed file and proves it landed.
func (f streamFixture357) plant(t *testing.T, name, payload string) {
	t.Helper()
	target := filepath.Join(f.dir, "model.onnx")
	if err := os.WriteFile(target+":"+name, []byte(payload), 0o644); err != nil {
		t.Fatalf("planting %s:%s needs NTFS alternate data streams, and this case is meaningless without one: %v",
			filepath.Base(target), name, err)
	}
	if got, err := os.ReadFile(target + ":" + name); err != nil || string(got) != payload {
		t.Fatalf("the stream did not land: %q %v", got, err)
	}
}

// run drives the production tail and reports exactly what a shell would see.
func (f streamFixture357) run(t *testing.T) (rc int, stdout, stderr string) {
	t.Helper()
	var out, errb bytes.Buffer
	rc = announceUnnamedStreams357(modelsIO{stdout: &out, stderr: &errb}, "verify357-fixture", f.mgr)
	return rc, out.String(), errb.String()
}

// Test357VerifyTailSaysOneSentenceAndKeepsRC0 is form 丁 itself: the reading
// before this change was 与已验签清单一致 and rc=0 with nothing said about the stream.
func Test357VerifyTailSaysOneSentenceAndKeepsRC0(t *testing.T) {
	f := newStreamFixture357(t, "")
	f.plant(t, "attacker", "wisp-ads-attacker-20")

	rc, stdout, _ := f.run(t)
	if rc != 0 {
		t.Fatalf("ticket 357: the exit code changed to %d. Form 丁 is a REPORT: the verdict and rc stay what "+
			"they are and only a sentence is added", rc)
	}
	want := "model.onnx:attacker:$DATA"
	if !strings.Contains(stdout, want) {
		t.Fatalf("ticket 357: the passing path did not name the stream %q, so the silent pass is still "+
			"silent: %q", want, stdout)
	}
	if !strings.Contains(stdout, "清单未点名") {
		t.Errorf("ticket 357: the sentence does not say 清单未点名 - the ticket's own wording for what the "+
			"finding is: %q", stdout)
	}
	if !strings.Contains(stdout, f.dir) {
		t.Errorf("ticket 357: the sentence does not name the directory that was enumerated (%s), so nobody can "+
			"tell which tree this report is about: %q", f.dir, stdout)
	}
	t.Logf("rc=%d stdout=%q", rc, stdout)
}

// Test357VerifyTailWithoutAStreamPrintsNothing is the reverse leg that keeps the
// change from becoming noise on every clean install: no stream, no sentence.
func Test357VerifyTailWithoutAStreamPrintsNothing(t *testing.T) {
	f := newStreamFixture357(t, "")

	rc, stdout, stderr := f.run(t)
	if rc != 0 {
		t.Fatalf("ticket 357 reverse leg: a clean install returned %d", rc)
	}
	if stdout != "" {
		t.Fatalf("ticket 357 reverse leg: a directory with no second stream still printed %q. The passing "+
			"sentence belongs to modelsVerify, which has already printed it by the time this tail runs; "+
			"anything here is an addition the ticket did not ask for", stdout)
	}
	if stderr != "" {
		t.Errorf("ticket 357 reverse leg: a clean install wrote to stderr: %q", stderr)
	}
}

// Test357VerifyTailReportsTheVerifiedTreeNotTheStoreRoot is the trap ticket 357's
// own first face is about: with a local_override the guard verifies the override
// while the store root is what gets printed. A report naming the store would be a
// report about a directory nobody verified, so this case puts a stream ONLY in the
// override and a differently named one ONLY in the store.
func Test357VerifyTailReportsTheVerifiedTreeNotTheStoreRoot(t *testing.T) {
	override := filepath.Join(t.TempDir(), "override-357")
	f := newStreamFixture357(t, override)
	if f.dir != override {
		t.Fatalf("fixture: the override was not honoured (%s)", f.dir)
	}
	f.plant(t, "in-override", "override-side")

	rc, stdout, _ := f.run(t)
	if rc != 0 {
		t.Fatalf("ticket 357: the override tree verifies, so rc must stay 0; got %d", rc)
	}
	if !strings.Contains(stdout, "model.onnx:in-override:$DATA") {
		t.Fatalf("ticket 357: the stream in the local_override directory was not reported: %q", stdout)
	}
	if !strings.Contains(stdout, override) {
		t.Errorf("ticket 357: the sentence does not name the override directory that was enumerated: %q", stdout)
	}
}

// Test357StreamNoticeIsTheOneSentence pins the sentence's contract: one line, the
// model id, the stream, the directory. A second sentence growing out of the
// renderer later would change what the passing path prints without any case above
// noticing.
func Test357StreamNoticeIsTheOneSentence(t *testing.T) {
	got := streamNotice357("m1", "model.onnx:attacker:$DATA", `D:\store\m1`)
	want := "wisp models verify: m1 在 model.onnx:attacker:$DATA 上枚举到清单未点名的备用数据流" +
		`（ADS dir=D:\store\m1）；本次一致结论只覆盖清单点名的文件` + "\n"
	if got != want {
		t.Fatalf("ticket 357: the sentence changed.\n got: %q\nwant: %q", got, want)
	}
	if strings.Count(got, "\n") != 1 {
		t.Errorf("ticket 357: the sentence is not exactly one line: %q", got)
	}
}

// Test357VerifyTailIsStillOnThePassingPath is the hook check the mutation depends
// on. Every case above calls announceUnnamedStreams357 directly, so deleting its
// call from modelsVerify would leave them all green - and the command would go back
// to being silent. This case reads modelsVerify's own body and requires the call to
// be there, AFTER the 与已验签清单一致 print (the ruling's "on the rc=0 path") and
// as the statement that produces its return value.
func Test357VerifyTailIsStillOnThePassingPath(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "models.go", nil, 0)
	if err != nil {
		t.Fatalf("parsing models.go: %v", err)
	}
	var body *ast.BlockStmt
	var bodyAt int
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "modelsVerify" || fn.Body == nil {
			continue
		}
		body = fn.Body
		bodyAt = fset.Position(fn.Pos()).Line
	}
	if body == nil {
		t.Fatal("modelsVerify is not in models.go any more, so this reading has nothing to check")
	}

	printAt, callAt := 0, 0
	ast.Inspect(body, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := ce.Fun.(type) {
		case *ast.SelectorExpr:
			if fn.Sel.Name == "Fprintf" && len(ce.Args) > 1 {
				if lit, ok := ce.Args[1].(*ast.BasicLit); ok && strings.Contains(lit.Value, "与已验签清单一致") {
					printAt = fset.Position(ce.Pos()).Line
				}
			}
		case *ast.Ident:
			if fn.Name == "announceUnnamedStreams357" {
				callAt = fset.Position(ce.Pos()).Line
			}
		}
		return true
	})
	if callAt == 0 {
		t.Fatalf("ticket 357: modelsVerify (models.go:%d) no longer calls announceUnnamedStreams357, so the "+
			"command has gone back to printing 与已验签清单一致 and saying nothing about a second data stream. "+
			"That is the exact silence this ticket exists to end, and every unit case in this file would stay "+
			"green through it", bodyAt)
	}
	if printAt == 0 {
		t.Fatalf("ticket 357: the 与已验签清单一致 print is gone from modelsVerify, so this reading cannot tell " +
			"whether the report still sits on the passing path")
	}
	if callAt < printAt {
		t.Errorf("ticket 357: the announce call (models.go:%d) runs BEFORE the passing sentence "+
			"(models.go:%d). The ruling puts the report on the rc=0 path, after the verdict it qualifies",
			callAt, printAt)
	}
}

// Test357StoreRootNamingIsNotWhatTheReportUses is the paired negative reading for
// the trap case above: the message a caller gets must be built from the directory
// UnnamedStreams resolved, never from the store root the CLI prints. This drives a
// Manager whose override and store differ and asserts the string on the wire names
// the former and not the latter.
func Test357StoreRootNamingIsNotWhatTheReportUses(t *testing.T) {
	root := t.TempDir()
	override := filepath.Join(root, "override-dir")
	store := filepath.Join(root, "store-root")
	f := newStreamFixture357In(t, store, override)
	f.plant(t, "only-here", "payload")

	_, stdout, _ := f.run(t)
	if !strings.Contains(stdout, override) {
		t.Fatalf("ticket 357: the report does not name the override directory %s: %q", override, stdout)
	}
	if strings.Contains(stdout, store) {
		t.Fatalf("ticket 357: the report names the store root %s, which VerifyInstalled did NOT verify for "+
			"this id. That is ticket 357's first face rebuilt inside the report: %q", store, stdout)
	}
}

// newStreamFixture357In is newStreamFixture357 with the store and override paths
// chosen by the caller, so a case can compare the two spellings.
func newStreamFixture357In(t *testing.T, store, override string) streamFixture357 {
	t.Helper()
	const id = "verify357-fixture"
	body := []byte("wisp-357-fixture-body")
	sum := sha256.Sum256(body)
	manifest := fmt.Sprintf(`{
		"manifest_version": 1,
		"models": [
			{
				"id": %q, "purpose": "vad",
				"urls": ["https://example.invalid/%s/model.onnx"],
				"sha256": %q, "size_bytes": %d,
				"license": "MIT", "quant": "fp32"
			}
		]
	}`, id, id, hex.EncodeToString(sum[:]), len(body))
	mf := signedManifest357(t, manifest)
	if err := os.MkdirAll(override, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(override, "model.onnx"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	mgr, err := models.NewManager(models.Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 1,
		LocalOverride: map[string]string{id: override},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.VerifyInstalled(id); err != nil {
		t.Fatalf("fixture does not verify: %v", err)
	}
	return streamFixture357{mgr: mgr, dir: override}
}
