package main

// Readface wave 1 (2026-11-17), cmd/wisp teeth for tickets 37 and 39:
//
//	Test37VerdictBatchFieldsReachThePacket   the pump's verdict reader carries
//	                                         approval.CardSummary's batch fields
//	                                         for a 10-path verdict and nothing
//	                                         for a 2-path one.
//	Test39ConfigStoreReadFaceSectionsErrors  the settings leg answers the full
//	                                         section tree and, over a broken
//	                                         file on disk, line-numbered errors.

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Test37VerdictBatchFieldsReachThePacket is AC#5's Go wiring tooth: an L2
// verdict carrying 10 paths (the R7-escalated shape arrives here the same
// way) renders into the packet's card with the aggregate fields, and a
// 2-path verdict renders none - the card JSON stays free of batch keys.
func Test37VerdictBatchFieldsReachThePacket(t *testing.T) {
	block := make(chan struct{})
	ui := approval.UIFuncs{
		PromptFn: func(context.Context, approval.Prompt) error {
			<-block
			return context.Canceled
		},
	}
	g := approval.New(approval.Options{UI: ui})
	revoke := g.AdmitTextTask("task-37-batch")
	defer revoke()

	big := tools.Decision{
		Tool: "fs.write", Level: risk.L2, RulesHit: []risk.RuleID{risk.R7},
		Reason: "batch",
		Paths: []string{
			"C:/d/x/a.txt", "C:/d/x/b.txt", "C:/d/x/c.txt", "C:/d/x/d.txt", "C:/d/x/e.txt",
			"C:/d/y/f.txt", "C:/d/y/g.txt", "C:/d/y/h.txt", "C:/d/y/i.txt", "C:/d/y/j.txt",
		},
		CorrelationID: "corr-37-big", TaskID: "task-37-batch",
	}
	done := make(chan tools.Answer, 1)
	go func() {
		a, _ := g.PendingApproval(context.Background(), big)
		done <- a
	}()
	waitForQueueDepth37(t, g, 1)

	verdicts := nativeVerdictsFromGate(g)
	var bigCard *panel.ApprovalCardView
	for _, v := range verdicts {
		if v.CorrelationID == "corr-37-big" {
			card := v.CardView()
			bigCard = &card
		}
	}
	if bigCard == nil {
		t.Fatal("the 10-path verdict is not in the packet")
	}
	if bigCard.BatchTotal != 10 || bigCard.BatchDeferred != 5 || len(bigCard.BatchDirs) != 2 {
		t.Errorf("card batch fields = total %d deferred %d dirs %v, want 10/5/2 dirs (the summary's own numbers)",
			bigCard.BatchTotal, bigCard.BatchDeferred, bigCard.BatchDirs)
	}
	if len(bigCard.BatchFirst) != approval.BatchPreviewCount {
		t.Errorf("card batchFirst = %d rows, want %d", len(bigCard.BatchFirst), approval.BatchPreviewCount)
	}
	// And the wire shape: the batch keys are on the card's JSON.
	data, err := json.Marshal(bigCard)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"batchTotal":10`) {
		t.Errorf("card wire carries no batchTotal 10: %s", data)
	}

	close(block)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the pending approval did not settle after the prompt returned")
	}

	// The mirror: a 2-path verdict carries no batch half at all.
	g2 := approval.New(approval.Options{UI: approval.UIFuncs{PromptFn: func(context.Context, approval.Prompt) error {
		return context.Canceled
	}}})
	g2.AdmitTextTask("task-37-small")
	_, _ = g2.PendingApproval(context.Background(), tools.Decision{
		Tool: "fs.write", Level: risk.L1, Paths: []string{"C:/d/one.txt", "C:/d/two.txt"},
		CorrelationID: "corr-37-small", TaskID: "task-37-small",
	})
	for _, v := range nativeVerdictsFromGate(g2) {
		if v.Batch != nil {
			t.Errorf("a 2-path verdict carries a batch half (%+v): the aggregate fields are for real batches only", v.Batch)
		}
		card := v.CardView()
		if card.HasBatch() {
			t.Error("a 2-path verdict's card claims a batch")
		}
	}
}

func waitForQueueDepth37(t *testing.T, g *approval.Gate, want int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if g.Queue().Depth() == want {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("queue depth never reached %d (now %d)", want, g.Queue().Depth())
}

// Test39ConfigStoreReadFaceSectionsErrors drives the real settings leg over a
// real config Manager: the section tree is the full D36 walk with the
// registry's tier words, a clean file answers no errors, and a broken file on
// disk answers line-numbered errors without moving the running config.
func Test39ConfigStoreReadFaceSectionsErrors(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.toml"
	body := "schema_version = 2\n\n[app]\ntheme = \"dark\"\n\n[ball]\nsize = 90\n"
	if err := writeFileSync39(path, body); err != nil {
		t.Fatalf("write config: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	leg := newConfigStore(mgr, nil, nil, "test-39")
	view, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	if len(view.Sections) == 0 {
		t.Fatal("ReadSettings answered no section rows - the D36 tree read face is dead")
	}
	var ballRow, riskRow bool
	for _, r := range view.Sections {
		switch r.Key {
		case "ball.size":
			if r.Value != "90" || r.Tier != "hot" {
				t.Errorf("ball.size row = %+v, want value 90 tier hot", r)
			}
			ballRow = true
		case "risk.permission_mode":
			if r.Tier != "locked" {
				t.Errorf("risk.permission_mode tier = %q, want locked", r.Tier)
			}
			riskRow = true
		}
	}
	if !ballRow || !riskRow {
		t.Errorf("section tree missing rows (ball.size=%v risk.permission_mode=%v)", ballRow, riskRow)
	}
	if len(view.LoadErrors) != 0 {
		t.Errorf("a clean file answered %d load errors, want 0", len(view.LoadErrors))
	}

	// The broken-file half: the Manager keeps running on the good config, and
	// the read face reports the disk file's line-numbered breakage.
	broken := "schema_version = 2\n\n[app]\ntheme = \"dark\"\nbogus_key = 1\n"
	if err := writeFileSync39(path, broken); err != nil {
		t.Fatalf("write broken config: %v", err)
	}
	view2, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings over broken file: %v", err)
	}
	if len(view2.LoadErrors) == 0 {
		t.Fatal("the broken file answered no load errors")
	}
	found := false
	for _, e := range view2.LoadErrors {
		if strings.Contains(e.Message, "bogus_key") && e.Line == 5 {
			found = true
		}
	}
	if !found {
		t.Errorf("no error row names bogus_key at line 5: %+v", view2.LoadErrors)
	}
	// The tree itself still answers from the RUNNING config (untouched by the
	// broken file) - that is the Manager's own contract, now visible to the page.
	if len(view2.Sections) == 0 {
		t.Error("the section tree vanished while the disk file was broken")
	}
}

func writeFileSync39(path, body string) error {
	return os.WriteFile(path, []byte(body), 0o600)
}
