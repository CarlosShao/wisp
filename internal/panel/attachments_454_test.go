package panel

// Ticket 454: mp4 and QuickTime share one ISO base media container, and the
// composer's sniff table used to give BOTH whitelisted types the same predicate.
// The table's order then decided the verdict, so an honest .mov (major brand
// "qt  ") declared video/quicktime was judged as video/mp4 and refused for
// disagreeing with a declaration that was right.
//
// The tests below judge the two directions that keep the repair honest:
//
//   - AC#1: the honest QuickTime file is STORED, under video/quicktime, with the
//     .mov extension - the verdict follows the brand, not the table's order;
//   - AC#3: a file declared video/mp4 whose bytes are a different type is still
//     refused. A repair that widens "declared vs sniffed must agree" into "any
//     video claim is fine" turns this green-when-it-must-be-red.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
)

// quickTimeBytes is a real QuickTime file header: a 20-byte "ftyp" box whose
// MAJOR BRAND is "qt  " (the four bytes a .mov written by QuickTime carries),
// a minor version, and the compatible brand a real file repeats.
var quickTimeBytes = []byte{
	0x00, 0x00, 0x00, 0x14, // box size: 20 bytes
	'f', 't', 'y', 'p',
	'q', 't', ' ', ' ', // major brand: QuickTime
	0x00, 0x00, 0x02, 0x00, // minor version
	'q', 't', ' ', ' ', // compatible brand
}

// Test454HonestQuickTimeFileIsStored is AC#1: a genuinely honest .mov must be
// accepted. Before the fix it failed with "... 声明为 video/quicktime，但文件头
// 表明它是 video/mp4", because both rows shared matchesISOBaseMedia and the mp4
// row came first.
func Test454HonestQuickTimeFileIsStored(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}
	src := fileSource(t, mustWrite(t, t.TempDir(), "holiday.mov", quickTimeBytes), "video/quicktime")
	ref, err := b.Ingest(context.Background(), src)
	if err != nil {
		t.Fatalf("AC#1 一枚诚实的 QuickTime 文件（主品牌 \"qt  \"）声明 video/quicktime 被拒了：%v", err)
	}
	if !ref.Stored || ref.MIME != "video/quicktime" || ref.Kind != "video" {
		t.Fatalf("ref = %+v, want a stored video/quicktime/video reference", ref)
	}
	if got := filepath.Ext(ref.Artifact); got != ".mov" {
		t.Errorf("artifact %q carries extension %q, want the QuickTime .mov the bytes actually are", ref.Artifact, got)
	}
	onDisk, err := os.ReadFile(filepath.Join(f.dir, ref.Artifact))
	if err != nil {
		t.Fatalf("the accepted .mov's bytes are not in the artifacts dir: %v", err)
	}
	if !bytes.Equal(onDisk, quickTimeBytes) {
		t.Errorf("stored bytes differ from the QuickTime source: %x != %x", onDisk, quickTimeBytes)
	}
	if f.puts != 1 {
		t.Errorf("one honest .mov produced %d writes, want 1", f.puts)
	}
}

// Test454DeclaredMP4WithForeignBytesIsStillRefused is AC#3: the repair must not
// buy AC#1 by loosening the type gate. Every case below declares video/mp4 (or
// video/quicktime) over bytes that are NOT that type, and every one must stay a
// loud refusal - this is the test that goes red if the sniff predicate or the
// declared-vs-sniffed step is widened to "any video is fine".
func Test454DeclaredMP4WithForeignBytesIsStillRefused(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}

	t.Run("declared video/mp4 over a PNG is refused", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "actually-a-png.mp4", pngBytes), "video/mp4")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "声明为 video/mp4，但文件头表明它是 image/png")
	})

	t.Run("declared video/mp4 over an executable is refused as unsupported", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "payload.mp4", exeBytes), "video/mp4")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "不受支持")
	})

	t.Run("declared video/quicktime over a PNG is refused", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "actually-a-png.mov", pngBytes), "video/quicktime")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "声明为 video/quicktime，但文件头表明它是 image/png")
	})

	t.Run("declared video/mp4 over QuickTime bytes is refused, not re-labelled a mov", func(t *testing.T) {
		// The brand decides: "qt  " is QuickTime. Accepting it under an mp4
		// declaration would be the mirror-image lie of ticket 454 - the panel
		// would store QuickTime bytes as .mp4 and tell the agent they are mp4.
		src := fileSource(t, mustWrite(t, t.TempDir(), "actually-a-mov.mp4", quickTimeBytes), "video/mp4")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "声明为 video/mp4，但文件头表明它是 video/quicktime")
	})

	if f.puts != 0 {
		t.Errorf("refusing foreign bytes wrote %d artifact(s)", f.puts)
	}
	ents, _ := os.ReadDir(f.dir)
	if len(ents) != 0 {
		t.Errorf("the artifacts dir holds %d entr(ies) after four refusals: %v", len(ents), ents)
	}
}
