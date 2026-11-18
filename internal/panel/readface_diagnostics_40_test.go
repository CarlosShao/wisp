package panel

// Readface wave 5 (2026-11-18) - the two diagnostics doors' internal/panel
// teeth.
//
// Ticket 40 AC#5's diagnostics.preview / diagnostics.export join the C17
// whitelist this wave. The teeth here are the package's own: the envelope
// parses both names, the router reaches the attached handler and refuses an
// unattached one by name, and the pair's split is measured (preview writes
// nothing, export builds - two methods, not one verb with an argument). The
// legs themselves (the one BundleOptions assembler, the gate, the zip) are
// judged on the cmd/wisp side; the closed-roster and page-literal rulers
// moved in the same change (inbound_roster_253_test.go's 18 -> 20, panel.ts's
// two request literals).

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestReadfaceWave5DoorsAreOnTheWhitelistAndRouted drives both doors through
// the production envelope parser and the production router. The pair is
// method-per-verb on purpose: the router's case labels ARE the proof that a
// preview envelope cannot reach the build leg, and an export envelope cannot
// reach a preview receipt.
func TestReadfaceWave5DoorsAreOnTheWhitelistAndRouted(t *testing.T) {
	diag := &recordingDiagnosticsHandlerWave5{}
	d := &ComposerDispatch{Diagnostics: diag}

	// diagnostics.preview with no selectors: parsed, routed, receipt
	// passthrough. The door takes no argument - the bundle's contents are the
	// machine's own, not a request parameter.
	reply, err := d.Handle(context.Background(),
		`{"method":"diagnostics.preview","requestId":"rw5-p1","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("diagnostics.preview was refused by the production hop: %v", err)
	}
	if reply != diag.previewReceipt {
		t.Errorf("diagnostics.preview reply = %q, want the handler's own receipt %q", reply, diag.previewReceipt)
	}
	if !diag.calledPreview {
		t.Error("diagnostics.preview never reached the attached handler")
	}

	// diagnostics.export: the second verb, its own method and its own slot.
	reply, err = d.Handle(context.Background(),
		`{"method":"diagnostics.export","requestId":"rw5-e1","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("diagnostics.export was refused by the production hop: %v", err)
	}
	if reply != diag.exportReceipt {
		t.Errorf("diagnostics.export reply = %q, want the handler's own receipt %q", reply, diag.exportReceipt)
	}
	if !diag.calledExport {
		t.Error("diagnostics.export never reached the attached handler")
	}
}

// TestReadfaceWave5UnattachedDoorsRefuseByName is the nil-socket half: a
// whitelisted name on a machine with no leg attached refuses by name - never
// a silence, and for an EXPORT never a fake success.
func TestReadfaceWave5UnattachedDoorsRefuseByName(t *testing.T) {
	d := &ComposerDispatch{}
	for _, method := range []string{MethodDiagnosticsPreview, MethodDiagnosticsExport} {
		raw := `{"method":"` + method + `","requestId":"rw5-nil","source":"panel-composer"}`
		_, err := d.Handle(context.Background(), raw)
		if !errors.Is(err, ErrNoHandlerAttached) {
			t.Errorf("%s with no handler attached returned %v, want ErrNoHandlerAttached", method, err)
		}
	}
}

// TestReadfaceWave5NamesStayClearOfTheBannedVocabularies pins the naming
// rulings measured, not asserted in prose: no "panel." prefix (ticket 181's
// extractor counts those literals), no grant-boundary verdict word
// (grantRouteWords reddens approved-shaped routes), and no task-axis word
// (internal/tools' ticket-188 scan reddens task/state in panel Method*
// values). The spellings land clear of all three BY VALUE; this case makes a
// later rename re-read the instruments instead of rediscovering them red.
func TestReadfaceWave5NamesStayClearOfTheBannedVocabularies(t *testing.T) {
	for _, c := range []struct {
		name  string
		value string
	}{
		{"MethodDiagnosticsPreview", MethodDiagnosticsPreview},
		{"MethodDiagnosticsExport", MethodDiagnosticsExport},
	} {
		lower := strings.ToLower(c.value)
		if strings.HasPrefix(lower, "panel.") {
			t.Errorf("%s = %q: the value claims the panel. prefix, and ticket 181's extractor counts those literals - spell the face's own verb instead", c.name, c.value)
		}
		for _, banned := range []string{"task", "state", "approve", "grant", "allow", "permit", "ratify", "authorize", "decide", "decision", "verdict"} {
			if strings.Contains(lower, banned) {
				t.Errorf("%s = %q: the value carries the banned word %q (the grant-boundary / task-axis instruments redden exactly that shape) - rename the door or re-read the instruments before landing", c.name, c.value, banned)
			}
		}
	}
}

type recordingDiagnosticsHandlerWave5 struct {
	calledPreview  bool
	calledExport   bool
	previewReceipt string
	exportReceipt  string
}

func (h *recordingDiagnosticsHandlerWave5) HandleDiagnosticsPreview(ctx context.Context, req ComposerRequest) (string, error) {
	h.calledPreview = true
	h.previewReceipt = "诊断包预览已生成。" // 形状占位由 cmd/wisp 的真腿判；这里只证路由与封套
	return h.previewReceipt, nil
}

func (h *recordingDiagnosticsHandlerWave5) HandleDiagnosticsExport(ctx context.Context, req ComposerRequest) (string, error) {
	h.calledExport = true
	h.exportReceipt = "已导出诊断包。"
	return h.exportReceipt, nil
}
