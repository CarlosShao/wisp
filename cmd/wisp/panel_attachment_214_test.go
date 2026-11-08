package main

// Ticket 214's teeth: the attachment ingest door is wired, and the tests below
// prove the bytes actually travel - envelope in, broker verdict, disk, draft
// ledger, snapshot wire, agent text. None of them stops at "the handler was
// called": the acceptance criterion is that an attachment REACHES the bottom,
// and that a refusal reaches it too, with its reason, where the user and the
// agent can read it.
//
// What each case pins:
//
//	Test214AttachmentAddEndToEndStoresBytesInArtifacts      AC#1+AC#2 正控: real chain,
//	                                                        real store, bytes on disk.
//	Test214AttachmentRefusalReachesTheSnapshotWithItsReason AC#2: 拒 + 为什么 into the
//	                                                        snapshot, nothing landed.
//	Test214DisguisedExecutableIsRefusedByItsBytes           AC#3: sniff, not extension.
//	Test214OversizeAttachmentIsRefusedBeforeAnyRead         AC#2: ceiling, constant time.
//	Test214PathSpellingDisplayNameIsRefusedAndNothingLands  AC#4 family: bare-name guard.
//	Test214IngestRefRenderingCarriesTheArtifactPath...      AC#5 half: ForAgent's
//	                                                        untrusted-external-text shape.
//	Test214PumpAssemblyCarriesTheAttachmentReaders          the wiring nail (run.go literal).
//	Test214DraftLedgerDropsOldestPastItsCap                 the ledger's own bound.

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

func auditQuiet214(string, ...any) {}

func newChain214(t *testing.T, dir string, drafts *attachmentDraftLedger) *panel.ComposerDispatch {
	t.Helper()
	disp, _, err := newComposerDispatchChain(dir, auditQuiet214, panelInboundActor, nil, nil, drafts)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	return disp
}

func attachmentEnvelope214(requestID, name, mime string, sizeBytes int, data []byte) string {
	return `{"method":"` + panel.MethodAttachmentAdd + `","requestId":"` + requestID +
		`","source":"` + panel.ComposerRequestSource + `","name":"` + name +
		`","declaredMime":"` + mime + `","sizeBytes":` +
		strconv.Itoa(sizeBytes) + `,"dataBase64":"` + base64.StdEncoding.EncodeToString(data) + `"}`
}

// pngBytes214 is a minimal PNG: the 8-byte signature the sniff rule matches,
// plus three bytes so the payload is more than the signature itself.
func pngBytes214() []byte {
	return []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00}
}

func artifactsEntries214(t *testing.T, dir string) []os.DirEntry {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(dir, "artifacts"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		t.Fatalf("read artifacts dir: %v", err)
	}
	return entries
}

func Test214AttachmentAddEndToEndStoresBytesInArtifacts(t *testing.T) {
	// The chain reads config.toml on the data root and never creates it (the
	// 128 discipline), so the root is seeded the way every other chain test
	// seeds it (newInboundDataDir33).
	dir := newInboundDataDir33(t, "ask_every_step")
	drafts := newAttachmentDraftLedger()
	disp := newChain214(t, dir, drafts)
	defer func() { _ = drafts.closeStore() }()

	png := pngBytes214()
	raw, herr := disp.Handle(context.Background(), attachmentEnvelope214("r-214-accept", "photo.png", "image/png", len(png), png))
	if herr != nil {
		t.Fatalf("an accepted attachment came back refused: %v", herr)
	}
	if raw != "" {
		t.Fatalf("an accepted attachment returned a receipt: %q", raw)
	}

	// The broker's first production caller stored REAL bytes in the REAL tree.
	entries := artifactsEntries214(t, dir)
	if len(entries) != 1 {
		t.Fatalf("artifacts dir holds %d entries, want exactly the one attachment: %v", len(entries), entries)
	}
	if !regexp.MustCompile(`^attachment-[0-9a-f]{16}\.png$`).MatchString(entries[0].Name()) {
		t.Errorf("artifact name %q is not the content-hash shape attachments.go mints", entries[0].Name())
	}
	stored, err := os.ReadFile(filepath.Join(dir, "artifacts", entries[0].Name()))
	if err != nil {
		t.Fatalf("read the stored artifact back: %v", err)
	}
	if string(stored) != string(png) {
		t.Errorf("the bytes on disk are not the bytes that were sent: %d stored vs %d sent", len(stored), len(png))
	}

	// The draft row says so too, and the snapshot carries it.
	rows := drafts.rows()
	if len(rows) != 1 || !rows[0].Stored || rows[0].MIME != "image/png" ||
		rows[0].SizeBytes != int64(len(png)) || rows[0].Artifact != entries[0].Name() {
		t.Fatalf("the ledger row does not describe the stored artifact: %+v", rows)
	}
	if drafts.reason() != "" {
		t.Errorf("an accepted ingest left a refusal sentence behind: %q", drafts.reason())
	}
	snap := panel.NewSnapshotPump(panel.PumpSources{
		Attachments:      drafts.rows,
		AttachmentReason: drafts.reason,
	}).Snapshot()
	if len(snap.Composer.Attachments) != 1 || !snap.Composer.Attachments[0].Stored {
		t.Fatalf("the snapshot lost the stored row: %+v", snap.Composer.Attachments)
	}
}

func Test214AttachmentRefusalReachesTheSnapshotWithItsReason(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	drafts := newAttachmentDraftLedger()
	disp := newChain214(t, dir, drafts)
	defer func() { _ = drafts.closeStore() }()

	notAnImage := []byte("this is definitely not an image")
	_, err := disp.Handle(context.Background(), attachmentEnvelope214("r-214-refuse", "notes.png", "image/png", len(notAnImage), notAnImage))
	if !errors.Is(err, panel.ErrAttachmentRejected) {
		t.Fatalf("a non-image payload was not refused with the sentinel: %v", err)
	}

	// Nothing landed: the refusal happened before any byte was written.
	if entries := artifactsEntries214(t, dir); len(entries) != 0 {
		t.Errorf("a refused attachment left %d artifacts behind: %v", len(entries), entries)
	}

	// The ledger row carries the human reason, and the snapshot shows 被拒＋为什么.
	rows := drafts.rows()
	if len(rows) != 1 || rows[0].Stored || rows[0].Reason == "" {
		t.Fatalf("the refusal did not reach the ledger as a row with a reason: %+v", rows)
	}
	if !strings.Contains(rows[0].Reason, "不受支持") {
		t.Errorf("the refusal reason does not say what the product supports: %q", rows[0].Reason)
	}
	snap := panel.NewSnapshotPump(panel.PumpSources{
		Attachments:      drafts.rows,
		AttachmentReason: drafts.reason,
	}).Snapshot()
	if snap.Composer.AttachmentReason != rows[0].Reason {
		t.Errorf("composer.attachmentError = %q, want the row's own reason", snap.Composer.AttachmentReason)
	}
	data, mErr := panel.NewSnapshotPump(panel.PumpSources{
		Attachments:      drafts.rows,
		AttachmentReason: drafts.reason,
	}).Marshal()
	if mErr != nil {
		t.Fatal(mErr)
	}
	if !strings.Contains(string(data), "不受支持") {
		t.Errorf("the wire lost the refusal reason:\n%s", data)
	}
}

func Test214DisguisedExecutableIsRefusedByItsBytes(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	drafts := newAttachmentDraftLedger()
	disp := newChain214(t, dir, drafts)
	defer func() { _ = drafts.closeStore() }()

	exe := []byte("MZ\x90\x00this is not a png at all, but the extension claims otherwise")
	_, err := disp.Handle(context.Background(), attachmentEnvelope214("r-214-disguise", "photo.png", "image/png", len(exe), exe))
	if !errors.Is(err, panel.ErrAttachmentRejected) {
		t.Fatalf("an executable named .png was not refused: %v", err)
	}
	if rows := drafts.rows(); len(rows) != 1 || !strings.Contains(rows[0].Reason, "可执行文件") {
		t.Fatalf("the refusal does not name the executable for what it is: %+v", drafts.rows())
	}
	if entries := artifactsEntries214(t, dir); len(entries) != 0 {
		t.Errorf("the disguised executable landed anyway: %v", entries)
	}
}

func Test214OversizeAttachmentIsRefusedBeforeAnyRead(t *testing.T) {
	dir := t.TempDir()
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer func() { _ = store.Close() }()
	broker, err := panel.NewAttachmentBroker(store, memory.ValidArtifactName, 0)
	if err != nil {
		t.Fatalf("NewAttachmentBroker: %v", err)
	}
	drafts := newAttachmentDraftLedger()

	// Open stays nil on purpose: the ceiling is judged from the declared size
	// before the source is opened, so a 2 GB pick costs no read at all.
	src := panel.AttachmentSource{DisplayName: "big.png", DeclaredMIME: "image/png", SizeBytes: panel.MaxAttachmentBytes + 1}
	ref, err := broker.Ingest(context.Background(), src)
	if !errors.Is(err, panel.ErrAttachmentRejected) {
		t.Fatalf("an oversize pick was not refused with the sentinel: %v", err)
	}
	if !strings.Contains(err.Error(), "超过附件上限") {
		t.Errorf("the refusal does not name the ceiling: %v", err)
	}
	if ref.Stored || ref.Artifact != "" {
		t.Errorf("a refused pick produced a stored reference: %+v", ref)
	}

	drafts.record(ref)
	wire, wErr := panel.NewSnapshotPump(panel.PumpSources{
		Attachments:      drafts.rows,
		AttachmentReason: drafts.reason,
	}).Marshal()
	if wErr != nil {
		t.Fatal(wErr)
	}
	if !strings.Contains(string(wire), "超过附件上限") {
		t.Errorf("the wire does not show the oversize refusal:\n%s", wire)
	}
	if entries := artifactsEntries214(t, dir); len(entries) != 0 {
		t.Errorf("the oversize pick landed anyway: %v", entries)
	}
}

func Test214PathSpellingDisplayNameIsRefusedAndNothingLands(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	drafts := newAttachmentDraftLedger()
	disp := newChain214(t, dir, drafts)
	defer func() { _ = drafts.closeStore() }()

	png := pngBytes214()
	_, err := disp.Handle(context.Background(), attachmentEnvelope214("r-214-traverse", "../evil.png", "image/png", len(png), png))
	if !errors.Is(err, panel.ErrAttachmentRejected) {
		t.Fatalf("a path-spelling name was not refused: %v", err)
	}
	if rows := drafts.rows(); len(rows) != 1 || !strings.Contains(rows[0].Reason, "裸文件名") {
		t.Fatalf("the refusal does not name the bare-name rule: %+v", drafts.rows())
	}
	if entries := artifactsEntries214(t, dir); len(entries) != 0 {
		t.Errorf("the traversal name landed anyway: %v", entries)
	}
}

func Test214IngestRefRenderingCarriesTheArtifactPathAndRefusalsToTheAgent(t *testing.T) {
	// AC#5's available half: the reference the ingest produced, rendered into
	// agent-visible text, is framed as stored bytes at an artifact path plus a
	// verbatim refusal - the untrusted-external-text shape the message door
	// will hand the model when ticket 35's socket lands. The C25 provenance
	// marking of that text is that door's job and is named in the evidence
	// file as still absent, never pretended here.
	artDir := filepath.Join("data", "artifacts")
	msg := panel.OutgoingMessage{
		Text: "看这张图",
		Attachments: []panel.AttachmentRef{
			{
				ID: "a", Name: "photo.png", MIME: "image/png", Kind: "image",
				SizeBytes: 11, Artifact: "attachment-deadbeef00000000.png", Stored: true,
			},
			{ID: "b", Name: "notes.png", Stored: false, Reason: "类型不受支持"},
		},
	}
	out := msg.ForAgent(artDir)
	if !strings.Contains(out, "看这张图") {
		t.Errorf("the agent text lost the user's own sentence: %q", out)
	}
	wantPath := filepath.Join(artDir, "attachment-deadbeef00000000.png")
	if !strings.Contains(out, "name=photo.png mime=image/png kind=image bytes=11 path="+wantPath) {
		t.Errorf("the stored row did not render as bytes-at-an-artifact-path:\n%s", out)
	}
	if !strings.Contains(out, "[附件未送达] notes.png: 类型不受支持") {
		t.Errorf("the refused row is not repeated verbatim where the model reads:\n%s", out)
	}
}

func Test214PumpAssemblyCarriesTheAttachmentReaders(t *testing.T) {
	// The wiring nail, in the shape inbound_roster_253 and the 265 ruler use:
	// read the assembly root off the disk and pin the two reader lines plus the
	// resident bind. Deleting any of the three turns this red instead of
	// silently restoring the always-empty packet 145 measured.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	pkgDir := filepath.Dir(thisFile)
	read := func(name string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(pkgDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		return string(data)
	}
	run := read("run.go")
	if !regexp.MustCompile(`Attachments:\s+rt\.attachmentDraftRows,`).MatchString(run) {
		t.Error("run.go's PumpSources literal no longer carries the attachment draft reader")
	}
	if !regexp.MustCompile(`AttachmentReason:\s+rt\.attachmentDraftReason,`).MatchString(run) {
		t.Error("run.go's PumpSources literal no longer carries the refusal-sentence reader")
	}
	if !strings.Contains(read("resident_windows.go"), "bindAttachmentDrafts(attachmentDrafts)") {
		t.Error("the resident boot no longer binds the ingest ledger onto the runtime")
	}
}

func Test214DraftLedgerDropsOldestPastItsCap(t *testing.T) {
	l := newAttachmentDraftLedger()
	for i := 0; i < attachmentDraftMaxRows+8; i++ {
		l.record(panel.AttachmentRef{ID: strconv.Itoa(i), Stored: true})
	}
	rows := l.rows()
	if len(rows) != attachmentDraftMaxRows {
		t.Fatalf("the ledger holds %d rows, want the cap %d", len(rows), attachmentDraftMaxRows)
	}
	if rows[0].ID != strconv.Itoa(8) {
		t.Errorf("the oldest row survived past the cap: rows[0].ID = %q, want 8", rows[0].ID)
	}
	if rows[len(rows)-1].ID != strconv.Itoa(attachmentDraftMaxRows+7) {
		t.Errorf("the newest row was dropped: rows[last].ID = %q", rows[len(rows)-1].ID)
	}
}
