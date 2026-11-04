package plugin

// Ticket 403: this test used to live in internal/risk/provenance_test.go, where
// importing internal/plugin closed the loop risk(test) -> plugin -> risk and the
// whole risk test binary failed to build ("import cycle not allowed in test").
// plugin -> risk is a real one-way production edge (command.go, netallow.go), so
// this is the side of the pair that can legally hold the integration test. It
// uses risk's EXPORTED API only: the C11 DisposalScope and the C25 taint store
// are wired together exactly as a session boundary wires them.
//
// The assertions are the ones ticket 160 left in place, unchanged: the deferred
// closer must be able to close the scope it was handed, Dispose must drop every
// taint, and a new session's store must start empty.

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// riskCycleMarker is the planted sensitive fragment. It is the same literal the
// risk package's own suite uses; it is repeated here because that constant is
// test-only and this file may not reach it.
const riskCycleMarker = "MARKER-QZX-98WVE7-TAINTED-SOURCE"

// riskCycleOptions mirrors internal/risk's baseOptions fixture: probing is
// skipped and the sync root is injected under a fake profile, so this test never
// depends on the machine's real OneDrive/Dropbox discovery. The sync root is not
// the subject here (notify is), but keeping the fixture identical means the
// C11 x C25 test runs on the same engine shape it always did.
func riskCycleOptions(t *testing.T) risk.ProvOptions {
	t.Helper()
	root := filepath.Join(t.TempDir(), "profile", "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	return risk.ProvOptions{
		NoProbe:   true,
		SyncRoots: []risk.SyncRoot{{Provider: "OneDrive", Path: root, Source: "options"}},
	}
}

// TestDisposalScopeClearsTaints is the ticket 403 relocation of the ticket 160
// AC#2 integration test: C11 disposal is what closes a C25 scope at session end,
// and closing it must drop every taint rather than leave the next session
// inheriting the old one.
func TestDisposalScopeClearsTaints(t *testing.T) {
	// Session end (C11 binding): closing the scope must drop every taint; the
	// new session's scope starts clean (AC #5).
	scope := NewDisposalScope("session-1", context.Background())
	p := risk.NewProvenance(riskCycleOptions(t))
	handle := p.OpenScope("session-1") // ticket 160: the handle is what closes it
	scope.Defer(func() {
		// Stricter than the old p.CloseScope("session-1") line, which could not
		// fail: a Defer-ed closer that cannot close its own scope is the bug.
		if err := handle.Close(); err != nil {
			t.Errorf("DisposalScope must be able to close the scope it was handed: %v", err)
		}
	})

	p.Mark("session-1", risk.SrcFSRead, "/legacy/secret", riskCycleMarker)
	if _, ok := p.Inspect("session-1", "notify", map[string]any{"text": riskCycleMarker}); !ok {
		t.Fatal("taint live within the session")
	}
	if res := scope.Dispose(); res.Incomplete {
		t.Fatalf("dispose incomplete: %+v", res)
	}
	if _, ok := p.Inspect("session-1", "notify", map[string]any{"text": riskCycleMarker}); ok {
		t.Fatal("disposed session must not retain taint")
	}
	p.OpenScope("session-2")
	if _, ok := p.Inspect("session-2", "notify", map[string]any{"text": riskCycleMarker}); ok {
		t.Fatal("new session does not inherit old taints")
	}
	if len(p.ScopeTaints("session-2")) != 0 {
		t.Fatal("new session store must start empty")
	}
}
