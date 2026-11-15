// Probe 214/AC#3: read the broker's sniff semantics off the REAL production
// broker (panel.NewAttachmentBroker + panel.AttachmentBroker.Ingest), with a
// throwaway sink so nothing touches disk. No production code is imported
// beyond internal/panel, and nothing is written anywhere.
//
// What this measures, case by case:
//   - whether the DECLARED mime is believed (it must not be),
//   - whether the file NAME's extension is ever consulted (it must not be),
//   - what the shared ISO-base-media predicate resolves "qt  " to (table
//     order: video/mp4 is the first entry whose match fires).
package main

import (
	"context"
	"fmt"
	"io"

	"github.com/CarlosShao/wisp/internal/panel"
)

type nopSink struct{}

func (nopSink) PutArtifact(context.Context, string, []byte) error { return nil }
func (nopSink) ArtifactsDir() string                             { return "" }

type caseT struct {
	label    string
	name     string
	declared string
	data     []byte
}

func png() []byte {
	return []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00}
}

func jpeg() []byte { return []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46} }

// isoBase builds a size+ftyp+brand box: bytes 0..3 a fake size, 4..7 "ftyp",
// 8..11 the brand, then padding.
func isoBase(brand string) []byte {
	b := []byte{0, 0, 0, 0x18, 'f', 't', 'y', 'p'}
	b = append(b, []byte(brand)...)
	return append(b, make([]byte, 8)...)
}

func main() {
	broker, err := panel.NewAttachmentBroker(nopSink{}, func(string) error { return nil }, 0)
	if err != nil {
		fmt.Println("broker construction failed:", err)
		return
	}
	cases := []caseT{
		{"png bytes / photo.png / declared image/png", "photo.png", "image/png", png()},
		{"MZ bytes / photo.png / declared image/png", "photo.png", "image/png",
			[]byte("MZ\x90\x00this is not a png at all")},
		{"png bytes / photo.exe / declared empty (extension lies LOW)", "photo.exe", "", png()},
		{"png bytes / photo.jpg / declared image/jpeg (declared lies)", "photo.jpg", "image/jpeg", png()},
		{"jpeg bytes / photo.png / declared image/png", "photo.png", "image/png", jpeg()},
		{"ISO 'qt  ' brand / clip.mov / declared video/quicktime", "clip.mov", "video/quicktime", isoBase("qt  ")},
		{"ISO 'qt  ' brand / clip.mov / declared video/mp4", "clip.mov", "video/mp4", isoBase("qt  ")},
		{"ISO 'mmp4' brand / clip.mp4 / declared empty", "clip.mp4", "", isoBase("mmp4")},
		{"ISO unknown brand 'heic' / clip.mp4 / declared video/mp4", "clip.mp4", "video/mp4", isoBase("heic")},
		{"ISO bytes shorter than 12 / clip.mp4 / declared video/mp4", "clip.mp4", "video/mp4",
			[]byte{0, 0, 0, 0x18, 'f', 't', 'y', 'p'}},
	}
	for i, c := range cases {
		data := c.data
		src := panel.AttachmentSource{
			DisplayName:  c.name,
			DeclaredMIME: c.declared,
			SizeBytes:    int64(len(data)),
			Open: func() (io.ReadCloser, error) {
				return io.NopCloser(bytesReader(data)), nil
			},
		}
		ref, ierr := broker.Ingest(context.Background(), src)
		fmt.Printf("[%02d] %s\n", i+1, c.label)
		fmt.Printf("     err=%v\n", ierr)
		fmt.Printf("     stored=%v mime=%q kind=%q artifact=%q reason=%q\n",
			ref.Stored, ref.MIME, ref.Kind, ref.Artifact, ref.Reason)
	}
}

type sliceReader struct {
	b []byte
	i int
}

func bytesReader(b []byte) *sliceReader { return &sliceReader{b: b} }

func (r *sliceReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
