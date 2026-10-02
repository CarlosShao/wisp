package approval

import (
	"reflect"
	"strings"
	"testing"
)

// Ticket 242 AC#2 (form A, per 242-precheck-1): the panel-bound read face
// (approval.PanelItem, ui.go:50-57) had zero instruments - no ruler counts
// its fields, so a new field could silently join what the panel is allowed
// to see. This ruler walks PanelItem by reflection and requires every leaf
// to be named in panelItemReadFace - the allowlist IS the read face; adding
// a field without declaring it goes red here.
//
// ⛔ Deliberately NOT panel.ApprovalCardView: that type's two-sided gate
// (TestApprovalCardViewJSONKeysMatchFrontendTypes) is already red by known
// pre-existing reasons; piling this ruler onto it would drown the mutation
// readings. PanelItem is the panel's read face in THIS package (242-precheck
// §2: it never enters Snapshot - the card view travels a parallel path).
var panelItemReadFace = []string{
	"CorrelationID",
	"Tool",
	"Level",
	"Reason",
	"Paths",
	"Depth",
}

func TestTicket242PanelItemReadFaceIsFullyDeclared(t *testing.T) {
	typ := reflect.TypeOf(PanelItem{})
	declared := map[string]bool{}
	for _, n := range panelItemReadFace {
		declared[n] = true
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !declared[f.Name] {
			t.Errorf("AC#2 RED: PanelItem.%s (%s) is visible to the panel but not declared in panelItemReadFace (ticket 242: a new field on the panel read face must be declared, not slipped in)", f.Name, f.Type)
		}
	}
	for n := range declared {
		if _, ok := typ.FieldByName(n); !ok {
			t.Errorf("AC#2 RED: panelItemReadFace declares %s but PanelItem has no such field; the read face and its declaration drifted", n)
		}
	}
}

// The read face must keep its strongest property: no field can express
// "allow" (SPEC-06 §9 layer 3 - a struct that cannot say allow is a stronger
// guarantee than a check that hopes the caller meant it). Any new field whose
// name or type smells like a grant/allow channel goes red here.
func TestTicket242PanelItemStaysGrantFree(t *testing.T) {
	typ := reflect.TypeOf(PanelItem{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name := strings.ToLower(f.Name)
		typName := strings.ToLower(f.Type.String())
		for _, banned := range []string{"grant", "allow", "approve", "nonce", "token"} {
			if strings.Contains(name, banned) || strings.Contains(typName, banned) {
				t.Errorf("AC#2 RED: PanelItem.%s (%s) smells like an allow channel (%q); the panel read face must stay grant-free (SPEC-06 §9 layer 3, ticket 242)", f.Name, f.Type, banned)
			}
		}
	}
}
