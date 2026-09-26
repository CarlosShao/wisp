package agent

// ACCEPTANCE PROBE, ticket 155 adversarial accept pass. Lives ONLY in the
// out-of-repo snapshot (D:/tmp/wisp155accept/snap6de/internal/agent/); a copy
// is committed under .scratch/wisp/probes/155-accept/ for re-runs. Zero bytes
// of it are in the delivered tree.
//
// It answers AC#3 sentence 2 ("does any externally visible reading change when
// you remove it") with a DIFFERENT ruler than the impl pass used. The impl
// pass and ticket 155's pass both counted the substring `"task":"` in the
// persisted JSONL. This probe parses every persisted record with
// encoding/json and asks whether the KEY `task` is present and what its VALUE
// is. Three shapes a substring count cannot tell apart: an empty value, a
// non-string value, and an id that is not the one belonging to that run.
//
// Helpers used (newHarness / withLogger / withConfig / buildRoundHistory /
// BudgetsFor / observe.InitLogWithRegistry) are all pre-existing ones from the
// delivered test files and internal/observe - no new seam is invented here.

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

func TestAccept155ParsedRecordsOnDisk(t *testing.T) {
	dir := t.TempDir()
	p, err := observe.InitLogWithRegistry(observe.LogConfig{Dir: dir}, observe.NewRegistry())
	if err != nil {
		t.Fatalf("observe init: %v", err)
	}
	defer func() { _ = p.Close() }()

	c := NewCompressor(BudgetsFor(4096), nil, WithLogger(slog.New(p.Handler())))
	hist := buildRoundHistory(6, 400)

	// leg 1 - direct call, nothing tagged.
	if _, rep, err := c.Compress(context.Background(), hist); err != nil || !rep.Ran {
		t.Fatalf("leg1 untagged: err=%v ran=%v", err, rep.Ran)
	}
	// leg 2 - direct call, ctx tagged by hand (this leg is what the two-leg
	// accept-r1 probe of ticket 153 can see and nothing else).
	leg2ID := newTaskID()
	if _, rep, err := c.Compress(withTraceTask(context.Background(), leg2ID), hist); err != nil || !rep.Ran {
		t.Fatalf("leg2 tagged: err=%v ran=%v", err, rep.Ran)
	}
	// leg 3 - the production leg: a real Loop.Run; only loop.go hands the id over.
	h := newHarness(t, "text-reply",
		withConfig(func(cc *Config) { cc.ContextWindow = 4096 }),
		withLogger(slog.New(p.Handler())))
	for _, m := range buildRoundHistory(4, 400) {
		h.loop.append(m)
	}
	res := h.run("收个尾")
	if !res.Compression.Ran {
		t.Fatalf("leg3 did not fold -> its record would be a phantom (status=%s)", res.Status)
	}
	if err := p.FlushNow(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	var records []map[string]any
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		b, rerr := os.ReadFile(filepath.Join(dir, f.Name()))
		if rerr != nil {
			t.Fatalf("read %s: %v", f.Name(), rerr)
		}
		for _, l := range strings.Split(string(b), "\n") {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			var m map[string]any
			if jerr := json.Unmarshal([]byte(l), &m); jerr != nil {
				t.Logf("PROBE-155-ACCEPT UNPARSEABLE %s", l)
				continue
			}
			if m["msg"] != "agent: history compressed" {
				continue
			}
			records = append(records, m)
		}
	}
	if len(records) == 0 {
		t.Fatalf("PROBE-155-ACCEPT no persisted compression record at all - dead tree, not a zero reading")
	}
	present := 0
	for i, m := range records {
		v, has := m["task"]
		if has {
			present++
		}
		t.Logf("PROBE-155-ACCEPT RECORD %d taskkey=%v taskval=%#v keys=%d", i, has, v, len(m))
	}
	leg3 := ""
	if len(records) == 3 {
		if v, has := records[2]["task"]; has {
			leg3, _ = v.(string)
		}
	}
	t.Logf("PROBE-155-ACCEPT COUNT records=%d taskkey_present=%d leg2_id=%s leg3_taskid=%s leg3_match=%v",
		len(records), present, leg2ID, res.TaskID, leg3 == res.TaskID)
}
