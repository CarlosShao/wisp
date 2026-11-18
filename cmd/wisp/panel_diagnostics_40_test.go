package main

// Readface wave 5 (2026-11-18), ticket 40 AC#5's teeth: the two diagnostics
// doors are wired to the REAL leg - the SAME diagnosticsBundleOptions180P
// assembler and the SAME Preview/Build pair the `wisp diagnostics` CLI uses,
// under the SAME [privacy] diagnostics_opt_in total switch - and the AC's
// three claims are each measured, not asserted:
//
//	Test40DiagnosticsGateIsTheCliTotalSwitch     the switch-off refusal names
//	                                             diagnostics_opt_in on BOTH
//	                                             verbs, and nothing is written
//	Test40DiagnosticsPreviewEntriesEqualTheZip   牙①: the preview's entry
//	                                             set EQUALS the exported
//	                                             zip's entry set, through
//	                                             the production chain
//	Test40DiagnosticsProbesDoNotComeOutTheOtherSide
//	                                             牙②: probe plaintext fed
//	                                             into the source dirs is
//	                                             grepped out of the unzipped
//	                                             bundle byte for byte - 0
//	                                             hits for audio / keys /
//	                                             transcript-full, with a
//	                                             positive control per class
//	                                             proving the grep can hit
//	Test40DiagnosticsOneEnumerateSourceNoSecondFileList
//	                                             source-level: the handler
//	                                             owns no enumerate of its
//	                                             own (the ⛔ preview-writes-
//	                                             one-list-export-another
//	                                             shape)
//	Test40DiagnosticsWiringNailInProductionChain 反向: the legs are on the
//	                                             assembled chain, not a nil
//	                                             socket
//
// THE GREP IS BYTE-LEVEL. Every negative reading below opens the exported
// zip, decompresses EVERY entry, and searches the bytes for the probe
// plaintext - never a whitelist, never a manifest entry name. Each probe
// class also gets a positive control: the same search over a byte that IS in
// the bundle (a marker in a log line, a marker in the config, the transcript
// probe itself via an explicit opt-in build) must hit, proving the search
// really can find that class of bytes if the exclusion ever broke. 探针全部
// 是假串（WISP40* 标记与 {{API_KEY_...}} 形状），没有任何真凭据进仓或进盘。

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// diagnosticsEnvelope40 builds one composer envelope for either verb.
func diagnosticsEnvelope40(requestID, method string) string {
	return `{"method":"` + method + `","requestId":"` + requestID + `","source":"` + panel.ComposerRequestSource + `"}`
}

// seedDiagnosticsDataDir40 lays out a data root the chain can assemble on and
// the bundle can read: a config with the switch ON, probe comments in the
// config, a log line carrying the log-class positive control, audio probe
// bytes planted where a careless collector would look, stray key material in
// the data root, and a transcript dir nobody opted into.
func seedDiagnosticsDataDir40(t *testing.T, optIn bool) string {
	t.Helper()
	dir := sealableTempDir124(t)
	cfg := "[risk]\npermission_mode = \"ask_every_step\"\n\n[privacy]\ndiagnostics_opt_in = " +
		map[bool]string{true: "true", false: "false"}[optIn] + "\n"
	if optIn {
		// The keys-class negative probe (a plaintext sk- key shape, fake) and
		// the keys-class POSITIVE control (an ordinary marker the redactor
		// passes) both live in the one file the bundle carries. If redaction
		// ever breaks, the probe plaintext comes out the other side and tooth
		// ② reddens; the control proves the config bytes really land in the
		// zip and the grep really can find them.
		cfg += "# key-probe-ctrl: WISP40CTRL-CONFIG-9f3\n"
		cfg += "# sk-wisp40probe9f3\n"
	}
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(cfg), 0o600); err != nil {
		t.Fatalf("cannot seed %s: %v", path, err)
	}
	if !optIn {
		return dir
	}
	logDir := filepath.Join(dir, logDirName)
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		t.Fatalf("cannot create the log dir: %v", err)
	}
	// The audio-class POSITIVE control rides the one log file the bundle
	// carries: a plain marker, in the zip by construction.
	logLine := "{\"level\":\"INFO\",\"msg\":\"wisp40 ctrl WISP40CTRL-LOG-9f3\"}\n"
	if err := os.WriteFile(filepath.Join(logDir, "wisp-20261118-001.jsonl"), []byte(logLine), 0o600); err != nil {
		t.Fatalf("cannot seed the log: %v", err)
	}
	// The audio-class NEGATIVE probes: a .wav with probe bytes INSIDE the log
	// dir (the one directory the builder walks - the prefix/ext filter must
	// keep it out) and a .opus at the data root (a directory the builder
	// never touches). If either shape ever enters the bundle, the probe bytes
	// come out and tooth ② reddens.
	wav := []byte("RIFF....WAVEfmt WISP40AUDIOPROBE-9f3")
	if err := os.WriteFile(filepath.Join(logDir, "wisp-audioprobe-9f3.wav"), wav, 0o600); err != nil {
		t.Fatalf("cannot seed the wav probe: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "audio-capture-probe.opus"), []byte("OggS WISP40AUDIOPROBE-9f3"), 0o600); err != nil {
		t.Fatalf("cannot seed the opus probe: %v", err)
	}
	// The keys-class NEGATIVE probe #2: key material sitting in the data root
	// (the secret store's neighbourhood). The builder reads exactly two
	// inputs - the log dir and config.toml - so this must never ship.
	secretsDir := filepath.Join(dir, "secrets")
	if err := os.MkdirAll(secretsDir, 0o700); err != nil {
		t.Fatalf("cannot create the secrets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(secretsDir, "probe-note.txt"), []byte("WISP40SECRET-PROBE-9f3"), 0o600); err != nil {
		t.Fatalf("cannot seed the secret probe: %v", err)
	}
	// The transcript-class NEGATIVE probe: a full transcript nobody opted
	// into. The builder reads NOTHING here by default.
	tDir := filepath.Join(dir, "transcripts")
	if err := os.MkdirAll(tDir, 0o700); err != nil {
		t.Fatalf("cannot create the transcripts dir: %v", err)
	}
	full := "WISP40TRANSFULL-9f3 用户在会话里说的完整原话，逐句记录，不该进诊断包。\n"
	if err := os.WriteFile(filepath.Join(tDir, "session-full.txt"), []byte(full), 0o600); err != nil {
		t.Fatalf("cannot seed the transcript probe: %v", err)
	}
	return dir
}

// unzipAll40 decompresses EVERY entry of the exported zip and returns the
// concatenation plus the entry names. The probes are searched over these
// bytes - never over a manifest, never over a whitelist.
func unzipAll40(t *testing.T, path string) (bytes string, names []string) {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open bundle zip %s: %v", path, err)
	}
	defer zr.Close()
	var b strings.Builder
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open zip entry %s: %v", f.Name, err)
		}
		data, _ := io.ReadAll(rc)
		_ = rc.Close()
		b.Write(data)
		names = append(names, f.Name)
	}
	return b.String(), names
}

// Test40DiagnosticsGateIsTheCliTotalSwitch: the switch-off refusal names
// diagnostics_opt_in on BOTH verbs, and a refused export writes nothing.
func Test40DiagnosticsGateIsTheCliTotalSwitch(t *testing.T) {
	dir := seedDiagnosticsDataDir40(t, false)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	ctx := context.Background()
	for _, method := range []string{panel.MethodDiagnosticsPreview, panel.MethodDiagnosticsExport} {
		_, err := disp.Handle(ctx, diagnosticsEnvelope40("rid-d-gate", method))
		if err == nil {
			t.Fatalf("%s ran with the total switch OFF: the gate is gone", method)
		}
		if !strings.Contains(err.Error(), "diagnostics_opt_in") {
			t.Errorf("%s switch-off refusal = %v, want it to name diagnostics_opt_in (the CLI's own switch)", method, err)
		}
	}
	// A refused export leaves no zip: the refusal cost zero bytes on disk.
	matches, _ := filepath.Glob(filepath.Join(dir, "wisp-diagnostics-*.zip"))
	if len(matches) != 0 {
		t.Errorf("a switch-off export left %d zip(s) behind (%v), want none", len(matches), matches)
	}
}

// Test40DiagnosticsPreviewEntriesEqualTheZip is 牙①: the preview's entry set
// EQUALS the exported zip's entry set, both read through the production
// chain. The set comparison is over the payload's entries field and the real
// archive's names - not over a manifest, not over prose.
func Test40DiagnosticsPreviewEntriesEqualTheZip(t *testing.T) {
	dir := seedDiagnosticsDataDir40(t, true)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	ctx := context.Background()
	previewReply, err := disp.Handle(ctx, diagnosticsEnvelope40("rid-d-prev", panel.MethodDiagnosticsPreview))
	if err != nil {
		t.Fatalf("diagnostics.preview refused on a wired chain with the switch on: %v\nreply=%q", err, previewReply)
	}
	prevParsed := parsePreviewEntries40(t, previewReply, DiagnosticsPreviewPayloadMarker)
	exportReply, err := disp.Handle(ctx, diagnosticsEnvelope40("rid-d-exp", panel.MethodDiagnosticsExport))
	if err != nil {
		t.Fatalf("diagnostics.export refused on a wired chain with the switch on: %v\nreply=%q", err, exportReply)
	}
	expParsed := parsePreviewEntries40(t, exportReply, DiagnosticsExportPayloadMarker)
	var outPath string
	if i := strings.Index(exportReply, "已导出诊断包："); i >= 0 {
		rest := exportReply[i+len("已导出诊断包："):]
		if j := strings.Index(rest, "（"); j >= 0 {
			outPath = rest[:j]
		}
	}
	if outPath == "" {
		t.Fatalf("the export receipt names no written path: %q", exportReply)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("the named bundle does not exist on disk (%s): %v", outPath, err)
	}
	_, zipNames := unzipAll40(t, outPath)
	zipSet := map[string]bool{}
	for _, n := range zipNames {
		zipSet[n] = true
	}
	prevSet := map[string]bool{}
	for _, n := range prevParsed.entries {
		prevSet[n] = true
	}
	if len(prevSet) != len(prevParsed.entries) || len(zipSet) != len(zipNames) {
		t.Fatalf("duplicate entry names: preview %d/%d, zip %d/%d",
			len(prevSet), len(prevParsed.entries), len(zipSet), len(zipNames))
	}
	var missing, extra []string
	for _, n := range prevParsed.entries {
		if !zipSet[n] {
			missing = append(missing, n)
		}
	}
	for _, n := range zipNames {
		if !prevSet[n] {
			extra = append(extra, n)
		}
	}
	if len(missing) != 0 || len(extra) != 0 {
		t.Fatalf("牙① 红: the preview's entry set does not equal the zip's entry set\n  preview-only: %v\n  zip-only:     %v", missing, extra)
	}
	// The export payload's own list must agree too: one enumerate source
	// means preview payload, export payload and zip all say the same thing.
	for _, n := range expParsed.entries {
		if !zipSet[n] {
			t.Errorf("the export payload lists %q but the zip does not carry it", n)
		}
	}
	if len(expParsed.entries) != len(zipNames) {
		t.Errorf("the export payload lists %d entries, the zip carries %d", len(expParsed.entries), len(zipNames))
	}
}

// Test40DiagnosticsProbesDoNotComeOutTheOtherSide is 牙②: the probe
// plaintext fed into the source dirs is grepped out of the unzipped bundle
// byte for byte - 0 hits for audio / keys / transcript-full - and each class
// carries a positive control proving the grep really can hit that class of
// bytes when a file IS included.
func Test40DiagnosticsProbesDoNotComeOutTheOtherSide(t *testing.T) {
	dir := seedDiagnosticsDataDir40(t, true)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	exportReply, err := disp.Handle(context.Background(), diagnosticsEnvelope40("rid-d-probe", panel.MethodDiagnosticsExport))
	if err != nil {
		t.Fatalf("diagnostics.export refused on a wired chain with the switch on: %v\nreply=%q", err, exportReply)
	}
	var outPath string
	if i := strings.Index(exportReply, "已导出诊断包："); i >= 0 {
		rest := exportReply[i+len("已导出诊断包："):]
		if j := strings.Index(rest, "（"); j >= 0 {
			outPath = rest[:j]
		}
	}
	if outPath == "" {
		t.Fatalf("the export receipt names no written path: %q", exportReply)
	}
	bundleBytes, zipNames := unzipAll40(t, outPath)
	if len(zipNames) == 0 {
		t.Fatal("the exported zip carries no entries: there is nothing the exclusion could be proven over")
	}

	// The three NEGATIVE readings, byte-level: probe plaintext, 0 hits.
	for _, probe := range []struct {
		class string
		probe string
	}{
		{"audio", "WISP40AUDIOPROBE-9f3"},
		{"keys (secret material planted in the data root)", "WISP40SECRET-PROBE-9f3"},
		{"keys (plaintext sk- key shape in config)", "sk-wisp40probe9f3"},
		{"transcript-full", "WISP40TRANSFULL-9f3"},
	} {
		if strings.Contains(bundleBytes, probe.probe) {
			t.Errorf("牙② 红 (exclusion broken): the %s probe plaintext %q came out the other side - it IS in the bundle bytes", probe.class, probe.probe)
		}
	}

	// The POSITIVE CONTROLS, per class: the same search over bytes that ARE
	// in the bundle must hit, so a green negative above is a proven reading,
	// not a broken grep.
	//   audio class: the marker rides the one log file the bundle carries.
	if !strings.Contains(bundleBytes, "WISP40CTRL-LOG-9f3") {
		t.Error("牙② 正控 failed (audio class): the log-class marker never hit the bundle bytes - the grep cannot prove the audio exclusion")
	}
	//   keys class: the marker rides config.toml, which ships as
	//   config.redacted.toml.
	if !strings.Contains(bundleBytes, "WISP40CTRL-CONFIG-9f3") {
		t.Error("牙② 正控 failed (keys class): the config-class marker never hit the bundle bytes - the grep cannot prove the keys exclusion")
	}
	//   transcript-full class: the SAME probe, included via an explicit
	//   per-item opt-in build over the same transcript dir, MUST hit - the
	//   one proof that if the default exclusion ever broke, this grep would
	//   catch the transcript's own bytes.
	o := observe.BundleOptions{
		OutPath:       filepath.Join(dir, "wisp-diagnostics-optin-posctrl.zip"),
		LogDir:        filepath.Join(dir, logDirName),
		ConfigPath:    filepath.Join(dir, configFileName),
		TranscriptDir: filepath.Join(dir, "transcripts"),
		OptIn:         observe.BundleOptIn{Transcripts: []string{"session-full.txt"}},
	}
	if _, err := observe.BuildDiagnosticsBundle(o); err != nil {
		t.Fatalf("the positive-control opt-in build failed: %v", err)
	}
	ctrlBytes, _ := unzipAll40(t, o.OutPath)
	if !strings.Contains(ctrlBytes, "WISP40TRANSFULL-9f3") {
		t.Error("牙② 正控 failed (transcript-full class): the opted-in transcript's probe never hit - the grep cannot prove the transcript exclusion")
	}
	// And the transcript class's second NEGATIVE reading, engine-level: the
	// SAME dir planted but NOTHING opted in - the default exclusion holds
	// even when the dir exists (the shape a mutation of the opt-in gate
	// would break, with a genuine probe hit to redden on).
	noOptIn := o
	noOptIn.OptIn = observe.BundleOptIn{}
	noOptIn.OutPath = filepath.Join(dir, "wisp-diagnostics-no-optin.zip")
	if _, err := observe.BuildDiagnosticsBundle(noOptIn); err != nil {
		t.Fatalf("the no-opt-in build failed: %v", err)
	}
	noOptBytes, _ := unzipAll40(t, noOptIn.OutPath)
	if strings.Contains(noOptBytes, "WISP40TRANSFULL-9f3") {
		t.Error("牙② 红 (exclusion broken): the un-opted-in transcript's probe came out the other side - the default exclusion is gone")
	}
}

// Test40DiagnosticsOneEnumerateSourceNoSecondFileList is the source-level
// tooth for the AC's hardest clause: the panel door owns NO enumerate of its
// own. It assembles the CLI's options (diagnosticsBundleOptions180P) and
// hands them to the engine's Preview/Build pair; a second file walk, a
// second zip writer or a hand-spelled entry list in this file is the
// "preview writes one list, the zip carries another" shape, and it must
// never exist.
func Test40DiagnosticsOneEnumerateSourceNoSecondFileList(t *testing.T) {
	src, err := os.ReadFile("panel_diagnostics_40.go")
	if err != nil {
		t.Fatalf("read the door's source: %v", err)
	}
	s := string(src)
	for _, must := range []string{
		"diagnosticsBundleOptions180P(",
		"observe.PreviewDiagnosticsBundle(",
		"observe.BuildDiagnosticsBundle(",
	} {
		if !strings.Contains(s, must) {
			t.Errorf("the door's source no longer rides the ONE enumerate source: %q is gone - the panel preview and the CLI's cannot be proven to agree without it", must)
		}
	}
	for _, banned := range []string{"os.ReadDir", "zip.NewWriter", "zip.OpenReader", "filepath.Walk", "filepath.WalkDir"} {
		if strings.Contains(s, banned) {
			t.Errorf("the door's source carries %q: a second enumerate/zip path in the panel leg is exactly the preview-vs-zip drift the AC bans - delete it and ride the engine's pair", banned)
		}
	}
}

// Test40DiagnosticsWiringNailInProductionChain is the reverse half: the legs
// are on the assembled chain, not a nil socket. A refusal that names
// ErrNoHandlerAttached for either door means the field is unplugged.
func Test40DiagnosticsWiringNailInProductionChain(t *testing.T) {
	dir := seedDiagnosticsDataDir40(t, true)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	if disp.Diagnostics == nil {
		t.Fatal("the production chain carries no diagnostics leg: the two doors are back on the nil socket that refuses by name")
	}
	if _, ok := disp.Diagnostics.(*diagnosticsDoor40); !ok {
		t.Fatalf("diagnostics leg is %T, want the real handler diagnosticsDoor40 (a stub here would be the non-真 answer)", disp.Diagnostics)
	}
	for _, method := range []string{panel.MethodDiagnosticsPreview, panel.MethodDiagnosticsExport} {
		_, err := disp.Handle(context.Background(), diagnosticsEnvelope40("rid-d-nail", method))
		if err != nil && strings.Contains(err.Error(), panel.ErrNoHandlerAttached.Error()) {
			t.Errorf("roster name %q is on the chain but its handler socket is nil: %v", method, err)
		}
	}
}

// parsePreviewEntries40 pulls the entries list off one payload line.
func parsePreviewEntries40(t *testing.T, reply, marker string) struct{ entries []string } {
	t.Helper()
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, marker) {
			payloadLine = strings.TrimPrefix(line, marker)
		}
	}
	if payloadLine == "" {
		t.Fatalf("reply carries no %s line: %q", marker, reply)
	}
	var envelope struct {
		Entries []string `json:"entries"`
	}
	if err := json.Unmarshal([]byte(payloadLine), &envelope); err != nil {
		t.Fatalf("payload is not JSON: %v\n%s", err, payloadLine)
	}
	if len(envelope.Entries) == 0 {
		t.Fatalf("payload lists no entries: %s", payloadLine)
	}
	return struct{ entries []string }{entries: envelope.Entries}
}
