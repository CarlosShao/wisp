package agent

// ACCEPTANCE PROBE for ticket 153 (accept-r1) - lives ONLY in the out-of-repo
// snapshot, never in the delivered tree. It answers the one reading the impl
// listed as "not measured" (its s7 item 1): what does the compression trace
// look like once it has gone through the real redacting JSON pipeline onto a
// file on disk, tagged and untagged.

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

func TestAccept153TraceOnDiskThroughRedactor(t *testing.T) {
	dir := t.TempDir()
	p, err := observe.InitLogWithRegistry(observe.LogConfig{Dir: dir}, observe.NewRegistry())
	if err != nil {
		t.Fatalf("observe init: %v", err)
	}
	defer func() { _ = p.Close() }()

	c := NewCompressor(BudgetsFor(4096), nil, WithLogger(slog.New(p.Handler())))
	hist := buildRoundHistory(6, 400)

	if _, _, err := c.Compress(context.Background(), hist); err != nil {
		t.Fatalf("untagged Compress: %v", err)
	}
	id := newTaskID()
	if _, _, err := c.Compress(withTraceTask(context.Background(), id), hist); err != nil {
		t.Fatalf("tagged Compress: %v", err)
	}
	if err := p.FlushNow(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	var lines []string
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		b, rerr := os.ReadFile(filepath.Join(dir, f.Name()))
		if rerr != nil {
			t.Fatalf("read %s: %v", f.Name(), rerr)
		}
		for _, l := range strings.Split(string(b), "\n") {
			if strings.Contains(l, "history compressed") {
				lines = append(lines, l)
			}
		}
	}

	if len(lines) != 2 {
		t.Fatalf("on-disk trace records = %d, want 2 (%v)", len(lines), lines)
	}
	for _, l := range lines {
		t.Logf("PROBE-153 DISK: %s", l)
	}
	untagged, tagged := lines[0], lines[1]
	if strings.Contains(untagged, `"task"`) {
		t.Errorf("untagged record carries a task key on disk: %s", untagged)
	}
	if strings.Contains(untagged, `task=`) || strings.Contains(untagged, `"task":""`) {
		t.Errorf("untagged record carries a placeholder task on disk: %s", untagged)
	}
	if !strings.Contains(tagged, `"task":"`+id+`"`) {
		t.Errorf("tagged record lost its task on disk: want %q in %s", id, tagged)
	}
}
