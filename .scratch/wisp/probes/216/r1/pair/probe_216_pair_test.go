package panel

// 216-r1's PAIRED PROBE - not committed into the repository's own test set, and
// deliberately so. It is copied verbatim into two OUT-OF-REPO export trees built
// with `git archive` (before = 5326f65e3819fd0165e51ba45ab8c172d029ecb2, after =
// 92e8cc575efe29387e5b277af4859a9701ad1e25) and run with the same ruler on the same
// platform. It names nothing ticket 216 added, so it compiles in both trees:
//   - in the before-tree it goes RED, and the red sentence is the "穿过" reading:
//     the menu row reaches the page carrying a live U+001B;
//   - in the after-tree it goes GREEN, which is the "洗掉" reading of the same
//     assertion, on the same assembly, with nothing edited in between.
//
// Run: go test ./internal/panel/ -run '^TestProbe216HostileDisplayNameAtTheDisplayBoundary$' -v

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

func probe216PngSource(name string) AttachmentSource {
	return AttachmentSource{
		DisplayName:  name,
		DeclaredMIME: "image/png",
		SizeBytes:    int64(len(pngBytes)),
		Open: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(pngBytes)), nil
		},
	}
}

func TestProbe216HostileDisplayNameAtTheDisplayBoundary(t *testing.T) {
	store, err := memory.Open(t.TempDir())
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	broker, err := NewAttachmentBroker(store, memory.ValidArtifactName, 0)
	if err != nil {
		t.Fatalf("NewAttachmentBroker: %v", err)
	}
	const hostile = "photo\x1b已允许.png"

	ref, err := broker.Ingest(context.Background(), probe216PngSource(hostile))
	if err != nil {
		t.Fatalf("the real broker refused the hostile name instead of accepting it: %v", err)
	}
	if !ref.Stored {
		t.Fatalf("Stored = false with a nil error: %s", ref.Reason)
	}
	// The premise the pair measures, asserted on BOTH trees so a red after-tree
	// reading cannot be mistaken for a fix that moved the ingress guard instead:
	// the entry guard still only refuses NUL/CR/LF/TAB, so an ESC name arrives.
	if !strings.ContainsRune(ref.Name, 0x1b) {
		t.Fatalf("premise broken in this tree: ref.Name = %q no longer carries the raw ESC", ref.Name)
	}

	snap := NewSnapshot(nil, nil, NewComposerState(risk.ModeAskEveryStep,
		UnsetWorkspaceView(), []AttachmentRef{ref}, MaxAttachmentBytes), time.Now())
	row := snap.Composer.Attachments[0].Name
	// And the page's own view of it: through JSON and back.
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back Snapshot
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	wireRow := back.Composer.Attachments[0].Name
	t.Logf("PROBE go-row=%q raw-ESC-in-go-row=%v wire-row=%q raw-ESC-in-wire-row=%v",
		row, strings.ContainsRune(row, 0x1b), wireRow, strings.ContainsRune(wireRow, 0x1b))

	for _, s := range []string{row, wireRow} {
		if i := strings.IndexFunc(s, unicode.IsControl); i >= 0 {
			t.Errorf("the menu row reached the page with a live control character U+%04X at offset %d: %q "+
				"(this is the AC#2 穿过 reading; the display boundary washed nothing in this tree)",
				s[i], i, s)
		}
	}
}
