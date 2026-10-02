package main

// Ticket 248 - the settings path driven for real, and the leak ruler that has to
// stay quiet.
//
// Everything here runs the production assembly: the one config.Manager this process
// would build, the one secret.NewStore DPAPI behind it, the same
// panel.ComposerDispatch the WebView2 host and `wisp panel-inbound` both call. The
// only seam is nothing - there is no fake in this file, which is what makes AC#2's
// zero-hit reading worth anything at all.
//
// The canary below is a made-up string invented for this test. It is not, and must
// not become, any real credential; nothing here prints it on a failure path either,
// because a t.Errorf that quotes the artifact would put it in the test log - the
// very surface the ruler is checking.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// canary248 is the sentinel fed through the credential write path. It is shaped
// like a vendor key on purpose, because the ruler that has to stay quiet is a
// shape search, not a name search.
const canary248 = "sk-canary248notarealkey0f2a9b7c"

// secretShape248 is the candidate-key ruler: the literal, plus the spelling a real
// provider key would have.
var secretShape248 = []*regexp.Regexp{
	regexp.MustCompile(`sk-[A-Za-z0-9]{12,}`),
	regexp.MustCompile(canary248),
}

// hits248 runs the ruler over one artifact and returns what it found. It is a
// function rather than a pile of ifs so the positive control and the real case use
// the SAME code path - a ruler that is only ever pointed at clean trees proves
// nothing about its own eyesight.
func hits248(label string, blobs ...[]byte) []string {
	var out []string
	for _, b := range blobs {
		for _, re := range secretShape248 {
			if loc := re.Find(b); loc != nil {
				out = append(out, label+" matches "+re.String())
			}
		}
	}
	return out
}

// settingsSeed is a config the settings fields can address: one preset provider,
// one catalog entry, and the strictest档 so nothing here depends on a mode write.
const settingsSeed = `[risk]
permission_mode = "ask_every_step"

[llm]
text_chain = ["deepseek/deepseek-chat"]

[llm.roles.chat]
provider = "deepseek"
model = "deepseek-chat"

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
context_window = 64000
`

func settingsDir248(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(settingsSeed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	return dir
}

// readAllFiles248 returns every byte under the data root, one blob per file. The
// DPAPI ciphertext belongs in this surface: a store that wrote the value in plain
// text would be caught here, which is the point AC#3 makes about "one store".
func readAllFiles248(t *testing.T, dir string) [][]byte {
	t.Helper()
	var out [][]byte
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		out = append(out, b)
		// File names are part of the surface too: a blob named after its owner's
		// key would leak the same way.
		out = append(out, []byte(path))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", dir, err)
	}
	return out
}

// newRealSettingsLeg builds the leg out of the two production objects, the way
// newComposerDispatchChain does, and returns the dispatch in front of it.
func newRealSettingsLeg(t *testing.T, dir string, audit *[]string, logbuf *bytes.Buffer) (*panel.ComposerDispatch, *configStore) {
	t.Helper()
	mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	secrets, err := secret.NewStore(dir)
	if err != nil {
		t.Fatalf("secret.NewStore: %v", err)
	}
	auditf := func(format string, args ...any) {
		*audit = append(*audit, "[audit] "+fmt.Sprintf(format, args...))
	}
	if logbuf != nil {
		old := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(logbuf, &slog.HandlerOptions{Level: slog.LevelDebug})))
		t.Cleanup(func() { slog.SetDefault(old) })
	}
	leg := newConfigStore(mgr, secrets, auditf, "test-248")
	disp := &panel.ComposerDispatch{
		Config: &panel.ConfigWriteHandler{Store: leg, Audit: auditf, Actor: "test-248"},
		Audit:  auditf,
	}
	return disp, leg
}

// AC#2's standing nail: the sentinel goes in through the credential write, and the
// whole artifact surface stays clean.
func TestAC2CredentialSentinelAppearsInNoArtifact(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	logbuf := &bytes.Buffer{}
	disp, leg := newRealSettingsLeg(t, dir, &audit, logbuf)

	raw := `{"method":"config.set","requestId":"r-248-cred","source":"panel-composer",` +
		`"configField":"provider_credential","configProvider":"deepseek","` + credentialValueKey + `":"` + canary248 + `"}`
	reply, err := disp.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("the credential write was refused: %v", err)
	}

	snap := panel.NewSnapshotPump(panel.PumpSources{Credential: leg.credentialStatus})
	snapshotBytes, err := snap.Marshal()
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if !bytes.Contains(snapshotBytes, []byte(`"credentialState"`)) {
		t.Errorf("the snapshot carries no credential dimension at all: %s", snapshotBytes)
	}

	auditBytes := []byte(strings.Join(audit, "\n"))
	logBytes := logbuf.Bytes()
	fileBytes := readAllFiles248(t, dir)

	var hits []string
	hits = append(hits, hits248("receipt", []byte(reply))...)
	hits = append(hits, hits248("audit", auditBytes)...)
	hits = append(hits, hits248("slog-sink", logBytes)...)
	hits = append(hits, hits248("snapshot", snapshotBytes)...)
	hits = append(hits, hits248("data-root", fileBytes...)...)
	if len(hits) > 0 {
		t.Errorf("the credential value reached an artifact (%d hits); the artifacts are not printed here on purpose",
			len(hits))
		for _, h := range hits {
			// Name the ruler and the surface, never the matched text.
			t.Errorf("leak surface: %s", strings.ReplaceAll(h, canary248, "<canary-redacted>"))
		}
	}

	// And the write itself happened, or the clean reading above is just a silence.
	if !strings.Contains(reply, "dpapi:") {
		t.Errorf("no reference came back, so nothing was really recorded: %q", reply)
	}
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(body), "api_key_ref = 'dpapi:") {
		t.Errorf("the reference never landed in config.toml:\n%s", body)
	}
}

// The positive control the rule above needs: the same ruler, pointed at an artifact
// that does carry the string, must fire. Without this the zero hits could mean the
// ruler is blind.
func TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	_, leg := newRealSettingsLeg(t, dir, &audit, nil)

	// The results channel is a text channel and is allowed to carry anything the
	// model said; parking the sentinel there is the shape the ruler exists for.
	pump := panel.NewSnapshotPump(panel.PumpSources{
		Credential: leg.credentialStatus,
		Results: func() []panel.ResultChunk {
			return []panel.ResultChunk{{CorrelationID: "ctl", Text: "here is the key: " + canary248}}
		},
	})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := hits248("planted-snapshot", data); len(got) == 0 {
		t.Fatal("the credential ruler did not fire on a snapshot that carries the sentinel: the zero-hit case above proves nothing")
	}
	if got := hits248("planted-audit", []byte("audit line: "+canary248)); len(got) == 0 {
		t.Fatal("the credential ruler did not fire on an audit line that carries the sentinel")
	}
}

// AC#3: the value lives in the one DPAPI store and nowhere else; the file keeps a
// reference; and the store really can read it back.
func TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	raw := `{"method":"config.set","requestId":"r-248-store","source":"panel-composer",` +
		`"configField":"provider_credential","configProvider":"deepseek","` + credentialValueKey + `":"` + canary248 + `"}`
	reply, err := disp.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("credential write: %v", err)
	}
	ref := refFromReceipt248(t, reply)

	secrets, err := secret.NewStore(dir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	got, err := secrets.Resolve(ref)
	if err != nil {
		t.Fatalf("the reference the page was told about does not resolve: %v", err)
	}
	if got != canary248 {
		// Deliberately no values in this message.
		t.Error("the blob behind the reference is not the value that was written")
	}

	blobs, err := filepath.Glob(filepath.Join(dir, "secrets", "*"))
	if err != nil {
		t.Fatalf("glob secrets dir: %v", err)
	}
	if len(blobs) != 1 {
		t.Errorf("secrets dir holds %d files, want exactly this one write", len(blobs))
	}
	// No second store appeared anywhere under the data root.
	extra, err := filepath.Glob(filepath.Join(dir, "*.key"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(extra) != 0 {
		t.Errorf("a second place for keys appeared under the data root: %v", extra)
	}
	// The reference is what config.toml carries, and the reference form is legal.
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(body), "api_key_ref = '"+ref+"'") {
		t.Errorf("config.toml does not carry the reference this write minted:\n%s", body)
	}
}

func refFromReceipt248(t *testing.T, reply string) string {
	t.Helper()
	for _, tok := range strings.Fields(strings.ReplaceAll(reply, "。", " ")) {
		if strings.HasPrefix(tok, "dpapi:") {
			return strings.Trim(tok, "，,")
		}
	}
	t.Fatalf("no reference appears in the receipt: %q", reply)
	return ""
}

// AC#2's structural half, and the reason the ruler above is not a coincidence: the
// shared envelope must not declare the key the credential value is bound from. A
// planted envelope that DOES declare it has to be caught by the same comparison,
// which is what makes this a gate rather than a note.
func TestAC2SharedEnvelopeCannotCarryTheCredentialValue(t *testing.T) {
	shared := jsonKeys248(reflect.TypeOf(panel.ComposerRequest{}))
	writeOnly := jsonKeys248(reflect.TypeOf(credentialWriteRequest{}))

	overlap := map[string]bool{}
	for k := range writeOnly {
		if shared[k] {
			overlap[k] = true
		}
	}
	for k := range overlap {
		if k == credentialValueKey {
			t.Errorf("panel.ComposerRequest declares %q: the credential value would be decoded into the envelope every answered method shares (J4 forbids it)", k)
		}
	}
	if !writeOnly[credentialValueKey] {
		t.Errorf("the write-only envelope lost its value key (%q): the credential route cannot be answered", credentialValueKey)
	}
	if shared[credentialValueKey] {
		t.Error("the shared envelope grew the credential key")
	}

	// The control: a type that binds a key the shared envelope also declares must
	// report the overlap, or the comparison above is decoration.
	type plantedCarrier struct {
		Method string `json:"method"`
		Text   string `json:"text"`
	}
	carrier := jsonKeys248(reflect.TypeOf(plantedCarrier{}))
	found := false
	for k := range carrier {
		if shared[k] && k != credentialValueKey {
			found = true
		}
	}
	if !found {
		t.Error("the key-overlap comparison found no overlap with a type that plainly shares method/text with ComposerRequest")
	}
}

// configField / configProvider / configValue are the three keys the settings route
// is allowed to share with the rest of the panel, and nothing else: an assertion
// about the shape of the envelope rather than about a comment describing it.
func TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope(t *testing.T) {
	shared := jsonKeys248(reflect.TypeOf(panel.ComposerRequest{}))
	for _, k := range []string{"configField", "configProvider", "configModel", "configValue"} {
		if !shared[k] {
			t.Errorf("the shared envelope lost the settings key %q", k)
		}
	}
	if shared[credentialValueKey] {
		t.Errorf("the shared envelope declares %q: settings writes would put a credential value on the envelope the message route also decodes", credentialValueKey)
	}
}

func jsonKeys248(ty reflect.Type) map[string]bool {
	out := map[string]bool{}
	if ty.Kind() != reflect.Struct {
		return out
	}
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			if f.Type.Kind() == reflect.Struct && !f.Anonymous {
				continue
			}
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				for k := range jsonKeys248(f.Type) {
					out[k] = true
				}
				continue
			}
			continue
		}
		out[name] = true
	}
	return out
}

// AC#1/AC#8 through the leg as main.go drives it: a settings write arrives on
// stdin, lands on disk, and the receipt the caller sees names both the key path and
// the restart tier. The unlisted name on the same run must leave the exit code at 1.
func TestAC1InboundLegAnswersSettingsRouteEndToEnd(t *testing.T) {
	dir := settingsDir248(t)

	set := `{"method":"config.set","requestId":"r-248-leg","source":"panel-composer",` +
		`"configField":"model_context_window","configProvider":"deepseek",` +
		`"configModel":"deepseek-chat","configValue":"200000"}` + "\n"
	get := `{"method":"config.get","requestId":"r-248-get","source":"panel-composer"}` + "\n"
	bad := `{"method":"panel.review.allow","requestId":"r-248-bad","source":"panel-composer"}` + "\n"

	code, out, errLog := runInboundLeg33(t, set+get+bad, dir)

	if code != 1 {
		t.Errorf("exit = %d, want 1 (one accepted, one refused): stdout=%q stderr=%q", code, out, errLog)
	}
	if !strings.Contains(out, "llm.providers.deepseek.models.deepseek-chat.context_window") {
		t.Errorf("the accepted write's receipt lost its key paths: stdout=%q", out)
	}
	if !strings.Contains(out, "重启") {
		t.Errorf("the receipt does not tell the user a restart is needed: stdout=%q", out)
	}
	if strings.Contains(out, "保存即生效") || strings.Contains(out, "立即生效") {
		t.Errorf("the leg claims a live save: stdout=%q", out)
	}
	if !strings.Contains(out, "凭据状态") {
		t.Errorf("config.get answered nothing readable: stdout=%q", out)
	}
	if !strings.Contains(out, "被拒绝") || !strings.Contains(out, "r-248-bad") {
		t.Errorf("the unlisted name was not refused by name: stdout=%q", out)
	}
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(body), "context_window = 200000") {
		t.Errorf("the settings write never reached disk:\n%s", body)
	}
	if !strings.Contains(string(errLog), "INBOUND-DISPATCH") {
		t.Errorf("the refusal left no audit line: stderr=%q", errLog)
	}
}

// The mode door and the settings door are answered by one assembly, and the three
// doors nobody owns still refuse by name (the shape ticket 33's slice B pinned,
// re-asserted here because this ticket added a socket to the same struct).
func TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	logbuf := &bytes.Buffer{}
	disp, _ := newRealSettingsLeg(t, dir, &audit, logbuf)
	// A dispatch without the older sockets, exactly as the production chain builds
	// them for workspace / attachment / message.
	for _, method := range []string{panel.MethodWorkspaceRequest, panel.MethodAttachmentAdd, panel.MethodMessageSend} {
		raw := `{"method":"` + method + `","requestId":"r-248-nil","source":"panel-composer"}`
		if _, err := disp.Handle(context.Background(), raw); !strings.Contains(err.Error(), "处理器未接入") {
			t.Errorf("%s with no socket: %v, want the unattached refusal", method, err)
		}
	}
}

// A write the schema would reject never reaches the file, and the leg says which
// field it refused - through the router, not around it.
func TestAC7InvalidSettingsValueIsRefusedBeforeTheFile(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	before, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	raw := `{"method":"config.set","requestId":"r-248-bad-value","source":"panel-composer",` +
		`"configField":"model_context_window","configProvider":"deepseek",` +
		`"configModel":"deepseek-chat","configValue":"not-a-number"}`
	reply, err := disp.Handle(context.Background(), raw)
	if err == nil {
		t.Fatalf("a non-numeric context window was accepted: %q", reply)
	}
	if !strings.Contains(err.Error(), "model_context_window") {
		t.Errorf("the refusal does not name the field: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config after: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Errorf("config.toml changed under a refused write:\n%s", after)
	}
	if len(audit) == 0 {
		t.Error("a refused settings write left no audit line")
	}
}

// The pump's credential dimension is a state, and the state must survive a real
// reference that has no blob behind it (the "配了没给" case the read side names).
func TestAC2SnapshotReportsRefWithoutBlobAsAPartState(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, leg := newRealSettingsLeg(t, dir, &audit, nil)
	// Point the provider at a reference nothing stored.
	raw := `{"method":"config.set","requestId":"r-248-dangling","source":"panel-composer",` +
		`"configField":"provider_api_key_ref","configProvider":"deepseek","configValue":"dpapi:none-of-these-blobs"}`
	if _, err := disp.Handle(context.Background(), raw); err != nil {
		t.Fatalf("the reference write was refused: %v", err)
	}
	if got := leg.credentialStatus(); got != panel.CredentialPartlyMissing {
		t.Errorf("credential status = %q, want %q", got, panel.CredentialPartlyMissing)
	}
	pump := panel.NewSnapshotPump(panel.PumpSources{Credential: leg.credentialStatus})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var snap map[string]any
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("snapshot is not JSON: %v", err)
	}
	composer, ok := snap["composer"].(map[string]any)
	if !ok {
		t.Fatalf("snapshot has no composer section: %s", data)
	}
	if composer["credentialState"] != string(panel.CredentialPartlyMissing) {
		t.Errorf("credentialState = %v, want partly_missing: %s", composer["credentialState"], data)
	}
	if composer["credentialKnown"] != true {
		t.Errorf("credentialKnown = %v, want true (a reader ran): %s", composer["credentialKnown"], data)
	}
	// A pump with no reader at all must say nobody asked, not "nothing recorded".
	noread, err := panel.NewSnapshotPump(panel.PumpSources{}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(noread, []byte(`"credentialKnown":false`)) ||
		!bytes.Contains(noread, []byte(`"credentialState":"unknown"`)) {
		t.Errorf("an unread credential dimension rendered as a known state: %s", noread)
	}
}
