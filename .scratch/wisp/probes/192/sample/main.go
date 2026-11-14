// Ticket 192 AC#1 Q4 -- "what does the user actually see when a preview is
// truncated". This program drives the real resolver over a real memory.Store
// (a real artifacts directory, the real name guard) and prints the first view
// verbatim plus the walk that proves the pointer finds the whole file.
//
// It is a probe, not a test: the criteria in internal/panel pin the same
// numbers. Run:
//
//	go run ./.scratch/wisp/probes/192/sample
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

func main() {
	dir, err := os.MkdirTemp("", "wisp-192-preview")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	store, err := memory.Open(dir, memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = store.Close() }()

	var full bytes.Buffer
	for i := 0; full.Len() < 300000; i++ {
		fmt.Fprintf(&full, "第%06d行：这是一段用于预览截断与续读的内容，带多字节字符。\n", i)
	}
	data := full.Bytes()
	ctx := context.Background()
	if err := store.PutArtifact(ctx, "big-log.txt", data); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	resolver, err := panel.NewPreviewResolver(store, memory.ValidArtifactName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	doc := resolver.Lookup(ctx, panel.PreviewURL("big-log.txt"))
	fmt.Printf("artifact bytes=%d kind=%s status=%d complete=%t next=%d\n",
		len(data), doc.Kind, doc.Status, doc.Complete, doc.NextOffset)
	fmt.Printf("body bytes=%d\n", len(doc.Body))
	lines := strings.SplitN(string(doc.Body), "\n", 4)
	fmt.Printf("line 1 (notice): %s\n", lines[0])
	fmt.Printf("line 2 (pointer): %s\n", lines[1])
	fmt.Printf("line 3 (head marker): %s\n", lines[2])
	// The elision line sits after the head window; find it by its constant.
	for _, line := range strings.Split(string(doc.Body), "\n") {
		if strings.HasPrefix(line, "[wisp 省略]") || strings.HasPrefix(line, "[wisp 正文] 结尾") {
			fmt.Printf("notice: %s\n", line)
		}
	}
	fmt.Printf("head bytes shown=%d first 20=%q\n", doc.NextOffset, firstN(string(doc.Body), 20))
	tailLine := "[wisp 正文] 结尾"
	for _, line := range strings.Split(string(doc.Body), "\n") {
		if strings.HasPrefix(line, tailLine) {
			fmt.Printf("notice: %s\n", line)
		}
	}

	// The walk: every byte, in windows, by the pointer alone.
	var got []byte
	offset := int64(0)
	windows := 0
	for {
		windows++
		if windows > 64 {
			fmt.Println("walk did not terminate")
			os.Exit(1)
		}
		w := resolver.Lookup(ctx, panel.PreviewURLAt("big-log.txt", offset, 65536))
		if w.Status != 200 {
			fmt.Fprintf(os.Stderr, "window at %d: %d %s\n", offset, w.Status, w.Reason)
			os.Exit(1)
		}
		got = append(got, w.Body...)
		if w.Complete {
			break
		}
		offset = w.NextOffset
	}
	fmt.Printf("walk windows=%d recovered=%d/%d bytes equal=%t sha256=%s\n",
		windows, len(got), len(data), bytes.Equal(got, data), shortSHA(got))
	fmt.Printf("artifact sha256=%s\n", shortSHA(data))
}

func firstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func shortSHA(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}
