package risk

// Ticket 160 AC#4 — "what this machine does NOT treat", pinned as readings so
// the sentence in docs/evidence/s1/160-r1-impl.md cannot quietly grow into "it
// treats everything". Every reading below is taken from OUTSIDE the function
// under test: a registry entry, a taint count, a verdict, a returned error, an
// injected log sink.
//
//   - a scope whose handle is thrown away with no Close is never reclaimed and
//     nothing says so: "the leak that is never touched again" still leaks;
//   - the Detector(scopeID) borrow keeps working only while the scope is open,
//     and it owns nothing: a close added at the borrow point would take the
//     task's taints away from the handle that owns them;
//   - a registration that Mark created with no owner has no close route at all
//     except re-opening that exact id;
//   - the borrowed face is one read method wide, so "close it here" cannot be
//     written without widening that seam first.
//
// Why the last one is a shape pin and not a behavior pin: AC#2 freezes the
// RiskLevel verdicts, so no hop added for AC#4 may change one. What is left to
// pin is the shape that makes the wrong edit impossible to write silently.

import (
	"fmt"
	"reflect"
	"testing"
)

// captureLogf records the lines the engine logs while the test runs and
// restores the previous sink afterwards. Same shape as the two captures already
// in provenance_test.go, so a nil sink at entry stays nil at exit.
func captureLogf(t *testing.T) *[]string {
	t.Helper()
	old := Logf
	lines := &[]string{}
	Logf = func(f string, a ...any) { *lines = append(*lines, fmt.Sprintf(f, a...)) }
	t.Cleanup(func() { Logf = old })
	return lines
}

// TestScopeBoundary160UntouchedLeakIsSilent is the first AC#4 sentence: a scope
// whose handle nobody closed — the production shape in
// cmd/wisp/panel_assets.go's (*taintSourceFlag).detector, which does
// `_ = prov.OpenScope(...)` — is never reclaimed. No goroutine, no timer, no
// finalizer watches it (AGENTS §1.2 ban #1, D38), so the registration and its
// marks stay for the life of the process and NOTHING is logged. That is the
// boundary: this machine treats "closed by somebody who does not own it", not
// "never closed at all".
func TestScopeBoundary160UntouchedLeakIsSilent(t *testing.T) {
	p := NewProvenance(baseOptions(t))
	logged := captureLogf(t)

	_ = p.OpenScope("panel-assets-l2") // handle thrown away on purpose: this is the leak
	p.Mark("panel-assets-l2", SrcFSRead, "/legacy/secret", marker)

	reg, known := p.scopes["panel-assets-l2"]
	if !known || reg == nil {
		t.Fatal("the thrown-away handle's registration must still be there: nothing reclaims it")
	}
	taints := len(p.ScopeTaints("panel-assets-l2"))
	t.Logf("READING leak.silent: registered=%v owner_was_issued=%v taints=%d logged_lines=%d",
		known, reg.owner != nil, taints, len(*logged))
	if taints != 1 {
		t.Fatalf("ScopeTaints=%d, want 1: the leaked scope keeps its marks", taints)
	}
	// The leaked scope is still fully live. That is the leak, not a defect of
	// this test: a victim's content stays matchable for the process lifetime.
	if _, hit := p.Detector("panel-assets-l2").TaintHit(map[string]any{"text": marker}); !hit {
		t.Fatal("a leaked scope must still gate like any other: its marks were never dropped")
	}
	if len(*logged) != 0 {
		t.Fatalf("the leak must be silent here (the G5 census names it, not this engine): %v", *logged)
	}
}

// TestScopeBoundary160BorrowedDetectorOwnsNothing is the second AC#4 sentence,
// and the reason Detector(scopeID) must NOT be given a close: it is a BORROW.
// Whoever opened the scope owns the close; the borrower only reads. The two
// readings below are exactly how a close added at the borrow point would show
// up — the second hit on the same borrowed detector goes quiet, and the owner's
// own Close stops being the owner's.
func TestScopeBoundary160BorrowedDetectorOwnsNothing(t *testing.T) {
	p := NewProvenance(baseOptions(t))
	owner := p.OpenScope("panel-assets-l2") // the one that owns the close
	p.Mark("panel-assets-l2", SrcFSRead, "/legacy/secret", marker)
	borrowed := p.Detector("panel-assets-l2") // same shape as the panel's detector()
	params := map[string]any{"text": marker}

	for i := 1; i <= 2; i++ {
		src, hit := borrowed.TaintHit(params)
		t.Logf("READING borrow.hit#%d: hit=%v src=%q taints=%d",
			i, hit, src, len(p.ScopeTaints("panel-assets-l2")))
		if !hit {
			t.Fatalf("hit #%d: the borrowed detector must stay live across calls "+
				"(a close at the borrow point is what kills it)", i)
		}
	}
	if got := len(p.ScopeTaints("panel-assets-l2")); got != 1 {
		t.Fatalf("a borrow spent the scope: ScopeTaints=%d want 1", got)
	}
	// The owner can still close what it opened: a borrow must not spend that.
	if err := owner.Close(); err != nil {
		t.Fatalf("owner's Close after two borrowed hits: %v "+
			"(ErrScopeAlreadyClosed/ErrScopeNotOpen here = a close was added at the borrow point)", err)
	}
	if _, hit := borrowed.TaintHit(params); hit {
		t.Fatal("after the owner closed the scope the borrowed detector must be silent")
	}
}

// TestScopeBoundary160OwnerlessRegistrationHasNoCloser pins the third residual:
// Mark on an id nobody opened registers the taint with owner nil (Mark's head,
// the fail-closed marking path), and Scope.Close only drops a registration it
// owns. With the by-id CloseScope route gone (ticket 160) such a registration
// has no closer at all: only a later OpenScope on that EXACT id hands ownership
// to somebody who can close it, and nothing in the engine enumerates it.
func TestScopeBoundary160OwnerlessRegistrationHasNoCloser(t *testing.T) {
	p := NewProvenance(baseOptions(t))
	p.Mark("typo-scope", SrcWebFetch, "https://v", marker)

	reg := p.scopes["typo-scope"]
	if reg == nil || reg.owner != nil {
		t.Fatalf("Mark on a never-opened id must register it with no owner, got %+v", reg)
	}
	other := p.OpenScope("another-task")
	mustClose(t, other) // a foreign handle cannot reach it
	t.Logf("READING ownerless: survives_foreign_close=%v taints=%d",
		p.scopes["typo-scope"] != nil, len(p.ScopeTaints("typo-scope")))
	if _, still := p.scopes["typo-scope"]; !still {
		t.Fatal("a foreign handle dropped a registration it never owned")
	}
	// The relief route, pinned so the sentence above stays checkable: re-open
	// that exact id (ownership moves to the new handle) and close it.
	late := p.OpenScope("typo-scope")
	mustClose(t, late)
	if _, still := p.scopes["typo-scope"]; still {
		t.Fatal("re-opening then closing the same id must still be able to drop it")
	}
}

// TestScopeBoundary160BorrowFaceCarriesNoCloser pins the shape the other three
// readings lean on: the borrowed detector holds a one-method seam, not the
// engine. Putting *Provenance back into that field is what would let a close be
// written at the borrow point without touching any declared contract, and this
// is the reading that notices — it is ticket 160 AC#4's reverse criterion.
func TestScopeBoundary160BorrowFaceCarriesNoCloser(t *testing.T) {
	field, ok := reflect.TypeOf(boundDetector{}).FieldByName("p")
	if !ok {
		t.Fatal("boundDetector has no p field: the borrow face moved, re-derive this pin")
	}
	face := field.Type
	names := make([]string, 0, face.NumMethod())
	for i := 0; i < face.NumMethod(); i++ {
		names = append(names, face.Method(i).Name)
	}
	t.Logf("READING borrow.face: type=%v kind=%v methods=%v", face, face.Kind(), names)
	if face.Kind() != reflect.Interface {
		t.Fatalf("the borrowed detector holds %v itself, so a close at the borrow point "+
			"becomes writable again: Detector must hand out detectorSeam", face)
	}
	if len(names) != 1 || names[0] != "Inspect" {
		t.Fatalf("the borrowed face must stay one read method (Inspect), got %v", names)
	}
}
