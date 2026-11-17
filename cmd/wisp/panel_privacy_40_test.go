package main

// Readface wave 1 (2026-11-17), ticket 40 AC#1's teeth: the three privacy
// doors are wired to the REAL leg - internal/memory's privacy layer over a
// real SQLite wisp.db - and each tooth counts the DB before and after, the
// way the AC words it ("delete-one/purge/export e2e against live DB, counts
// match"). None of them stops at a nil check:
//
//	Test40PrivacyDeleteOneDeletesExactlyOneRow  正控+计数: one row out, the
//	                                            others untouched.
//	Test40PrivacyPurgeClearsTheDomainAndCounts  正控+计数: n rows out, the
//	                                            count in the receipt is the
//	                                            store's own.
//	Test40PrivacyExportCarriesTheDomainJSON     正控: the export payload is
//	                                            the store's own envelope.
//	Test40PrivacyWiringNailInProductionChain    反向: unplug the leg and the
//	                                            door refuses by name.
//	Test40PrivacyRefusalsNameTheDomainVocabulary 域词与空 id 的拒绝都有名字。

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

func privacyEnvelope40(requestID, method, domain, itemID string) string {
	env := `{"method":"` + method + `","requestId":"` + requestID + `","source":"` + panel.ComposerRequestSource + `","privacyDomain":"` + domain + `"`
	if itemID != "" {
		env += `,"privacyItemId":"` + itemID + `"`
	}
	return env + `}`
}

// seedMemories40 adds n rows to the memory domain and returns their ids, so
// every count below is the DB's own answer, not a test-side tally.
func seedMemories40(t *testing.T, dir string, n int) []string {
	t.Helper()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	defer func() { _ = st.Close() }()
	ctx := context.Background()
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		id, err := st.AddMemory(ctx, "privacy tooth row "+string(rune('a'+i)), "tooth")
		if err != nil {
			t.Fatalf("AddMemory: %v", err)
		}
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	return ids
}

func countMemories40(t *testing.T, dir string) int {
	t.Helper()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	defer func() { _ = st.Close() }()
	rows, err := st.ListMemories(context.Background())
	if err != nil {
		t.Fatalf("ListMemories: %v", err)
	}
	return len(rows)
}

func Test40PrivacyDeleteOneDeletesExactlyOneRow(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	ids := seedMemories40(t, dir, 3)
	if before := countMemories40(t, dir); before != 3 {
		t.Fatalf("seed count = %d, want 3", before)
	}
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(context.Background(),
		privacyEnvelope40("rid-p-del", panel.MethodPrivacyDeleteOne, "memory", ids[0]))
	if err != nil {
		t.Fatalf("privacy.delete-one refused on a wired chain: %v\nreply=%q", err, reply)
	}
	if after := countMemories40(t, dir); after != 2 {
		t.Errorf("DB count after delete-one = %d, want 2 (one row out, the rest untouched)", after)
	}
	if !strings.Contains(reply, ids[0]) {
		t.Errorf("receipt does not name the deleted id %q: %q", ids[0], reply)
	}
}

func Test40PrivacyPurgeClearsTheDomainAndCounts(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	seedMemories40(t, dir, 3)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(context.Background(),
		privacyEnvelope40("rid-p-purge", panel.MethodPrivacyPurge, "memory", ""))
	if err != nil {
		t.Fatalf("privacy.purge refused on a wired chain: %v\nreply=%q", err, reply)
	}
	if after := countMemories40(t, dir); after != 0 {
		t.Errorf("DB count after purge = %d, want 0", after)
	}
	if !strings.Contains(reply, "3") {
		t.Errorf("receipt does not name the store's own count (3): %q", reply)
	}
}

func Test40PrivacyExportCarriesTheDomainJSON(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	seedMemories40(t, dir, 2)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(context.Background(),
		privacyEnvelope40("rid-p-exp", panel.MethodPrivacyExport, "memory", ""))
	if err != nil {
		t.Fatalf("privacy.export refused on a wired chain: %v\nreply=%q", err, reply)
	}
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, PrivacyExportPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, PrivacyExportPayloadMarker)
		}
	}
	if payloadLine == "" {
		t.Fatalf("export reply carries no %s line: %q", PrivacyExportPayloadMarker, reply)
	}
	var envelope struct {
		Domain  string           `json:"domain"`
		Count   int              `json:"count"`
		Records []map[string]any `json:"records"`
	}
	if err := json.Unmarshal([]byte(payloadLine), &envelope); err != nil {
		t.Fatalf("export payload is not JSON: %v\n%s", err, payloadLine)
	}
	if envelope.Domain != "memory" || len(envelope.Records) != 2 {
		t.Errorf("export envelope = domain %q count %d records %d, want the store's own 2 rows",
			envelope.Domain, envelope.Count, len(envelope.Records))
	}
	if !strings.Contains(reply, "2 条记录") {
		t.Errorf("receipt does not state the record count read off the payload: %q", reply)
	}
}

func Test40PrivacyWiringNailInProductionChain(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	if disp.Privacy == nil {
		t.Fatal("the production chain carries no privacy leg: the three doors are back on the nil socket that refuses by name")
	}
	if _, ok := disp.Privacy.(*privacyHandler40); !ok {
		t.Fatalf("privacy leg is %T, want the real handler privacyHandler40 (a stub here would be the non-真 answer)", disp.Privacy)
	}
	// The three roster names are answered by the assembled chain, not only by
	// the guard: a door that parse accepts but nobody handles is the
	// silent-drop shape the dispatch's default branch exists to refuse.
	for _, method := range []string{panel.MethodPrivacyDeleteOne, panel.MethodPrivacyPurge, panel.MethodPrivacyExport} {
		env := privacyEnvelope40("rid-p-nail", method, "memory", "")
		if _, err := disp.Handle(context.Background(), env); err != nil &&
			strings.Contains(err.Error(), panel.ErrNoHandlerAttached.Error()) {
			t.Errorf("roster name %q is on the chain but its handler socket is nil: %v", method, err)
		}
	}
}

func Test40PrivacyRefusalsNameTheDomainVocabulary(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	ctx := context.Background()
	// Unknown domain: the refusal lists the five D35 domains instead of
	// guessing.
	_, err = disp.Handle(ctx, privacyEnvelope40("rid-p-dom", panel.MethodPrivacyPurge, "sessions", ""))
	if err == nil || !strings.Contains(err.Error(), "task_log") {
		t.Errorf("unknown-domain refusal = %v, want one that names the D35 vocabulary", err)
	}
	// Delete-one without an id: nothing is deleted and the refusal says so.
	seedMemories40(t, dir, 1)
	_, err = disp.Handle(ctx, privacyEnvelope40("rid-p-noid", panel.MethodPrivacyDeleteOne, "memory", ""))
	if err == nil || !strings.Contains(err.Error(), "privacyItemId") {
		t.Errorf("no-id refusal = %v, want one that names the missing id", err)
	}
	if n := countMemories40(t, dir); n != 1 {
		t.Errorf("a refused delete-one moved the DB (%d rows), want 1", n)
	}
}
