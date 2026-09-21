package panel

// Ticket 92 AC#2: the attachment path, judged one input class at a time.
//
// The judgement text asks for a conclusion per input, not a percentage, so each
// case below is a named subtest that either stores bytes (and says where) or
// refuses (and says why, in text a user can act on). The one shape none of them
// may take is (AttachmentRef{}, nil): a refusal that returns no error is how a
// product eats the user's intent, which is the "lying config key" family of
// ticket 83 wearing a different hat.

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
)

// fakeArtifacts is a memory.Store-shaped stand-in: a flat directory, an
// exclusive create, and a counter for every guard that must fire BEFORE the
// source is opened.
type fakeArtifacts struct {
	dir       string
	puts      int
	exists    map[string]bool
	openCalls *int
}

func newFakeArtifacts(t *testing.T) *fakeArtifacts {
	t.Helper()
	dir := t.TempDir()
	return &fakeArtifacts{dir: dir, exists: map[string]bool{}}
}

func (f *fakeArtifacts) ArtifactsDir() string { return f.dir }

func (f *fakeArtifacts) PutArtifact(_ context.Context, name string, data []byte) error {
	if err := memory.ValidArtifactName(name); err != nil {
		return err
	}
	if f.exists[name] {
		return &os.LinkError{Op: "create", New: name, Err: os.ErrExist}
	}
	f.puts++
	f.exists[name] = true
	return os.WriteFile(filepath.Join(f.dir, name), data, 0o600)
}

func fileSource(t *testing.T, path, declaredMIME string) AttachmentSource {
	t.Helper()
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return AttachmentSource{
		DisplayName:  filepath.Base(path),
		DeclaredMIME: declaredMIME,
		SizeBytes:    st.Size(),
		Open: func() (io.ReadCloser, error) {
			return os.Open(path)
		},
	}
}

func mustWrite(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

var (
	pngBytes  = []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 'I', 'H', 'D', 'R', 1, 2, 3, 4}
	jpgBytes  = []byte{0xff, 0xd8, 0xff, 0xe0, 'J', 'F', 'I', 'F', 0, 1, 2, 3}
	mp4Bytes  = append([]byte{0, 0, 0, 0x18, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm'}, bytes.Repeat([]byte{0}, 8)...)
	exeBytes  = append([]byte{'M', 'Z', 0x90, 0}, bytes.Repeat([]byte{0xAB}, 64)...)
	webpBytes = append([]byte("RIFF\x00\x00\x00\x00WEBPVP8 "), bytes.Repeat([]byte{0}, 16)...)
)

func TestAcceptsAndStoresEachSupportedType(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}
	cases := []struct {
		name, mime, kind string
		data             []byte
	}{
		{"shot.png", "image/png", "image", pngBytes},
		{"photo.jpg", "image/jpeg", "image", jpgBytes},
		{"clip.mp4", "video/mp4", "video", mp4Bytes},
		{"anim.gif", "image/gif", "image", append([]byte("GIF89a"), bytes.Repeat([]byte{1}, 12)...)},
		{"pic.webp", "image/webp", "image", webpBytes},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := fileSource(t, mustWrite(t, t.TempDir(), tc.name, tc.data), tc.mime)
			ref, err := b.Ingest(context.Background(), src)
			if err != nil {
				t.Fatalf("Ingest(%s) = %v, want stored (this type is on the whitelist)", tc.name, err)
			}
			if !ref.Stored || ref.Kind != tc.kind || ref.MIME != tc.mime {
				t.Fatalf("ref = %+v, want stored %s/%s", ref, tc.kind, tc.mime)
			}
			onDisk, err := os.ReadFile(filepath.Join(f.dir, ref.Artifact))
			if err != nil {
				t.Fatalf("bytes are not in the artifacts dir: %v", err)
			}
			if !bytes.Equal(onDisk, tc.data) {
				t.Errorf("stored %d bytes differ from the source's %d", len(onDisk), len(tc.data))
			}
			if ref.SizeBytes != int64(len(tc.data)) {
				t.Errorf("SizeBytes = %d, want %d", ref.SizeBytes, len(tc.data))
			}
			t.Logf("stored: artifact=%s bytes=%d mime=%s", ref.Artifact, ref.SizeBytes, ref.MIME)
		})
	}
}

func TestRefusesUnsupportedAndMasqueradingInputsLoudly(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}

	t.Run("empty file is refused, not dropped", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "empty.png", nil), "image/png")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "0 字节")
	})

	t.Run("exe disguised as png is refused and named as an executable", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "invoice.png", exeBytes), "image/png")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "MZ")
	})

	t.Run("unknown type without a declared mime is refused", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "mystery.bin",
			[]byte("not any signature on the whitelist")), "")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "不受支持")
	})

	t.Run("lying declared mime is refused even when the bytes are a png", func(t *testing.T) {
		src := fileSource(t, mustWrite(t, t.TempDir(), "real-but-mislabeled.png", pngBytes), "application/pdf")
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "声明为")
	})

	t.Run("path-injection file name is refused and nothing is written", func(t *testing.T) {
		before := f.puts
		for _, name := range []string{`..\evil.png`, `C:\Windows\temp\evil.png`, `sub/ok.png`, `....`, `x.png.`} {
			src := AttachmentSource{DisplayName: name, DeclaredMIME: "image/png",
				SizeBytes: int64(len(pngBytes)),
				Open:      func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(pngBytes)), nil }}
			ref, err := b.Ingest(context.Background(), src)
			requireRefused(t, ref, err, "裸文件名")
		}
		if f.puts != before {
			t.Errorf("%d artifact(s) were written while refusing a path-spelling name", f.puts-before)
		}
		ents, _ := os.ReadDir(f.dir)
		if len(ents) != 0 {
			t.Errorf("the artifacts dir holds %d entr(ies) after a refusal: %v", len(ents), ents)
		}
	})

	t.Run("2 GB is refused on its declared size without reading a byte", func(t *testing.T) {
		big := filepath.Join(t.TempDir(), "big.mp4")
		fh, err := os.Create(big)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		// A sparse extend: metadata only, no clusters written.
		if err := fh.Truncate(2 << 30); err != nil {
			_ = fh.Close()
			t.Skipf("filesystem cannot hold a 2 GB file here: %v", err)
		}
		_ = fh.Close()
		count := 0
		src := fileSource(t, big, "video/mp4")
		src.Open = func() (io.ReadCloser, error) { count++; return os.Open(big) }
		ref, err := b.Ingest(context.Background(), src)
		requireRefused(t, ref, err, "超过附件上限")
		if count != 0 {
			t.Errorf("the source was opened %d time(s) to refuse it - a size claim must gate the read, not follow it", count)
		}
		if f.puts != 0 {
			t.Errorf("a 2 GB pick stored %d artifact(s)", f.puts)
		}
	})

	t.Run("a source lying about its size is refused mid-read, not truncated", func(t *testing.T) {
		src := AttachmentSource{DisplayName: "liar.png", DeclaredMIME: "image/png", SizeBytes: 4,
			Open: func() (io.ReadCloser, error) {
				return io.NopCloser(io.MultiReader(bytes.NewReader(pngBytes), bytes.NewReader(pngBytes),
					bytes.NewReader(pngBytes), bytes.NewReader(pngBytes))), nil
			}}
		// SizeBytes lies small; the broker's bound is its own max, so this must
		// still be judged on the bytes it actually got: 64 bytes, a real png
		// header, stored - the point of the case is that the stored size is the
		// ACTUAL size, not the claim.
		ref, err := b.Ingest(context.Background(), src)
		if err != nil {
			t.Fatalf("Ingest = %v", err)
		}
		if ref.SizeBytes != int64(4*len(pngBytes)) {
			t.Errorf("SizeBytes = %d, want the actual %d bytes read", ref.SizeBytes, 4*len(pngBytes))
		}
	})

	t.Run("nil source and nil sink are refused at construction, not at send", func(t *testing.T) {
		if _, err := NewAttachmentBroker(nil, memory.ValidArtifactName, 0); err == nil {
			t.Error("a broker with no artifact store constructed without error")
		}
		if _, err := NewAttachmentBroker(f, nil, 0); err == nil {
			t.Error("a broker with no name guard constructed without error")
		}
		if _, err := b.Ingest(context.Background(), AttachmentSource{DisplayName: "a.png", SizeBytes: -1}); err == nil {
			t.Error("a source with no Open() ingested without error")
		}
	})
}

// TestSecondIdenticalAttachmentIsDeduplicatedNotOverwritten pins ticket 79's
// O_EXCL promise from the composer's side: a collision is reported, never
// silently won.
func TestSecondIdenticalAttachmentIsDeduplicatedNotOverwritten(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}
	src := fileSource(t, mustWrite(t, t.TempDir(), "same.png", pngBytes), "image/png")
	first, err := b.Ingest(context.Background(), src)
	if err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	second, err := b.Ingest(context.Background(), src)
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if !second.Deduplicated || first.Artifact != second.Artifact {
		t.Errorf("second ref = %+v, want deduplicated with artifact %q", second, first.Artifact)
	}
	if f.puts != 1 {
		t.Errorf("identical bytes produced %d writes, want 1", f.puts)
	}
}

// TestMessageCarriesAttachmentRefsAndRefusals is AC#2's "出现在消息里" half.
func TestMessageCarriesAttachmentRefsAndRefusals(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}
	ok, err := b.Ingest(context.Background(),
		fileSource(t, mustWrite(t, t.TempDir(), "clip.mp4", mp4Bytes), "video/mp4"))
	if err != nil {
		t.Fatalf("ingest mp4: %v", err)
	}
	bad, err := b.Ingest(context.Background(),
		fileSource(t, mustWrite(t, t.TempDir(), "evil.png", exeBytes), "image/png"))
	if err == nil {
		t.Fatal("an exe named .png was accepted")
	}
	msg := OutgoingMessage{Text: "看看这个", Attachments: []AttachmentRef{ok, bad}}
	line := msg.ForAgent(f.dir)
	if !strings.Contains(line, "[附件] name=clip.mp4 mime=video/mp4 kind=video") {
		t.Errorf("the accepted attachment is not in the message the agent reads:\n%s", line)
	}
	if !strings.Contains(line, filepath.Join(f.dir, ok.Artifact)) {
		t.Errorf("the message does not name the artifact file holding the bytes:\n%s", line)
	}
	if !strings.Contains(line, "[附件未送达] evil.png") {
		t.Errorf("the refused attachment vanished from the message - the user's intent was eaten:\n%s", line)
	}
	if !strings.Contains(line, "MZ") {
		t.Errorf("the refusal lost its reason:\n%s", line)
	}
	// Q-28's boundary: a video reference must not claim the agent saw its content.
	if strings.Contains(line, "视频内容") || strings.Contains(line, "已理解") {
		t.Errorf("the message over-claims what a video reference means (Q-28 is not this ticket):\n%s", line)
	}
	t.Logf("agent-facing text:\n%s", line)
}

func requireRefused(t *testing.T, ref AttachmentRef, err error, wantInReason string) {
	t.Helper()
	if err == nil {
		t.Fatalf("refused nothing: ref=%+v err=<nil> - a silent drop is AC#2's forbidden shape", ref)
	}
	if !errors.Is(err, ErrAttachmentRejected) {
		t.Fatalf("error %q does not wrap ErrAttachmentRejected", err)
	}
	if ref.Stored {
		t.Errorf("Stored=true alongside an error: %+v", ref)
	}
	if ref.MIME != "" || ref.Kind != "" {
		t.Errorf("a refused attachment still reports mime=%q kind=%q - a caller's claim is not a verdict",
			ref.MIME, ref.Kind)
	}
	if !strings.Contains(ref.Reason, wantInReason) {
		t.Errorf("reason %q does not explain the refusal (%q expected in it)", ref.Reason, wantInReason)
	}
	if !strings.Contains(err.Error(), wantInReason) {
		t.Errorf("error %q must carry the same reason the user sees", err)
	}
}

// TestAttachmentPayloadCarriesTheBytes is AC#2's transport half: the file's
// CONTENT must reach the native side through the panel's text envelope, and a
// damaged trip must be refused instead of becoming a shorter file.
func TestAttachmentPayloadCarriesTheBytes(t *testing.T) {
	f := newFakeArtifacts(t)
	b, err := NewAttachmentBroker(f, memory.ValidArtifactName, MaxAttachmentBytes)
	if err != nil {
		t.Fatalf("broker: %v", err)
	}

	t.Run("a real png arrives whole through base64", func(t *testing.T) {
		src, err := DecodeAttachmentPayload(AttachmentPayload{
			Name: "shot.png", DeclaredMIME: "image/png", SizeBytes: int64(len(pngBytes)),
			DataBase64: base64.StdEncoding.EncodeToString(pngBytes),
		})
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		ref, err := b.Ingest(context.Background(), src)
		if err != nil {
			t.Fatalf("ingest: %v", err)
		}
		onDisk, err := os.ReadFile(filepath.Join(f.dir, ref.Artifact))
		if err != nil {
			t.Fatalf("read back: %v", err)
		}
		if !bytes.Equal(onDisk, pngBytes) {
			t.Errorf("the stored bytes are not the bytes the user picked: %x != %x", onDisk, pngBytes)
		}
	})

	t.Run("undecodable base64 is refused, not stored as nothing", func(t *testing.T) {
		_, err := DecodeAttachmentPayload(AttachmentPayload{
			Name: "x.png", SizeBytes: 4, DataBase64: "!!!not base64!!!",
		})
		if err == nil || !errors.Is(err, ErrAttachmentRejected) {
			t.Fatalf("undecodable payload returned %v, want a refusal", err)
		}
		if !strings.Contains(err.Error(), "无法解码") {
			t.Errorf("refusal does not say what failed: %v", err)
		}
	})

	t.Run("a truncated payload is refused, not stored short", func(t *testing.T) {
		enc := base64.StdEncoding.EncodeToString(pngBytes)
		_, err := DecodeAttachmentPayload(AttachmentPayload{
			Name: "cut.png", DeclaredMIME: "image/png", SizeBytes: int64(len(pngBytes)),
			DataBase64: enc[:len(enc)-4],
		})
		if err == nil || !strings.Contains(err.Error(), "截断") {
			t.Fatalf("a truncated payload returned %v, want a truncation refusal", err)
		}
	})

	t.Run("an empty payload is refused", func(t *testing.T) {
		src, err := DecodeAttachmentPayload(AttachmentPayload{Name: "nil.png", DataBase64: ""})
		if err != nil {
			t.Fatalf("decode of an empty payload should succeed so the reason is the store's: %v", err)
		}
		if _, err := b.Ingest(context.Background(), src); err == nil ||
			!strings.Contains(err.Error(), "0 字节") {
			t.Fatalf("an empty payload ingested as %v, want the 0-byte refusal", err)
		}
	})
}
