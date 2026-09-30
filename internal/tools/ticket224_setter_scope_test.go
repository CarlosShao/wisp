package tools

// Ticket 224-r2 cell N#6 - widening the reflection pin's setter half.
//
// What 224-v1 left on the record (§1, and §6-c row 6): the existing
// TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority checks Bridge fields for
// names starting with "grant" AND `Kind() == reflect.Func`. `Bridge.grants` is
// interface-typed, so that loop never looked at it, and a second grant-shaped
// channel of ANY other kind - interface, struct, map, pointer - would have passed
// it. 224-v1 judged the hole real, "轻", and today unfilled; its own note said the
// sanctioned seam is safe "today" by way of a grep, not by way of a test.
//
// This case widens the reading without touching that one (the original stays on the
// books, still passing, so nobody can read its change as a relaxation):
//
//	ONE channel  - exactly one grant-shaped field on Bridge, and it must be the
//		       sanctioned read seam of type GrantSource. A second one of any
//		       Kind is the shape the Func-only loop let through.
//	NO setter    - no method on *Bridge may name a grant and an act
//		       (set/record/revoke/insert/add/delete/swap), and no exported
//		       Set* method at all: re-wiring an authorization source at runtime
//		       is the same M-7/C-3 family as a caller-supplied session field,
//		       just reached after construction instead of during it.
//
// Why the sanctioned field is allowed to exist and be interface-typed:
// internal/tools/grant.go's header states the rule - the identity lives on the
// injected source, which the composition root minted, so "a caller that wanted a
// different session would have to swap the bridge's GrantSource, which is a
// construction-time act visible in the audit log, not an argument". This case is
// what keeps that sentence a reading of the code rather than a promise.

import (
	"reflect"
	"strings"
	"testing"
)

// grantActs are the verbs that would turn the read seam into a write or a swap.
var grantActs = []string{"set", "record", "revoke", "insert", "add", "delete", "swap", "replace"}

func TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter(t *testing.T) {
	bTyp := reflect.TypeOf(Bridge{})

	// HALF 1 - exactly one grant-shaped field, and it is the seam itself.
	var grantFields []reflect.StructField
	for i := 0; i < bTyp.NumField(); i++ {
		f := bTyp.Field(i)
		if !strings.Contains(strings.ToLower(f.Name), "grant") {
			continue
		}
		grantFields = append(grantFields, f)
	}
	if len(grantFields) != 1 {
		t.Fatalf("Bridge carries %d grant-shaped fields, want exactly 1: %v. The pre-224 bridge "+
			"carried none and ticket 224 added one read seam; a second channel is a second place "+
			"an authorization source can come from, and only one of them is in the audit log",
			len(grantFields), fieldNames(grantFields))
	}
	f := grantFields[0]
	want := reflect.TypeOf((*GrantSource)(nil)).Elem()
	if f.Type != want {
		t.Errorf("Bridge.%s has type %v, want the tools.GrantSource interface: a differently shaped "+
			"grant field is how a swappable or writable channel reappears while still reading as "+
			"a seam to the old Func-only loop", f.Name, f.Type)
	}
	// The old loop's predicate, restated as a positive so the widening is auditable:
	// whatever the grant field is, it is NOT a function-typed channel either.
	if f.Type.Kind() == reflect.Func {
		t.Errorf("Bridge.%s is a function-typed grant channel, which the original pin already refuses "+
			"two ways at once", f.Name)
	}
	for _, k := range []reflect.Kind{reflect.Ptr, reflect.Map, reflect.Slice, reflect.Struct} {
		if f.Name == "grants" && f.Type.Kind() == k {
			t.Errorf("the sanctioned seam migrated to Kind %v: it is a one-method interface on "+
				"purpose, and a concrete type here means someone gave it methods", k)
		}
	}

	// HALF 2 - no setter, on either side of the naming rule.
	pTyp := reflect.PointerTo(bTyp)
	for i := 0; i < pTyp.NumMethod(); i++ {
		m := pTyp.Method(i)
		low := strings.ToLower(m.Name)
		if strings.Contains(low, "grant") {
			for _, act := range grantActs {
				if strings.Contains(low, act) {
					t.Errorf("(*Bridge).%s names an act on a grant: the bridge may read one seam and "+
						"must not be able to install, remove or trade it at runtime (that is a "+
						"caller-side authority reached after construction)", m.Name)
				}
			}
		}
		if strings.HasPrefix(m.Name, "Set") {
			t.Errorf("(*Bridge).%s is an exported setter on the enforcement layer: options are chosen "+
				"at New() and appear in the assembly log; a runtime re-wire does not", m.Name)
		}
	}

	// HALF 3 - the seam itself is still read-only, which is what makes the single
	// allowed field safe. Restated here (the original case asserts it too) so this
	// file cannot pass while the interface grows a write method.
	src := reflect.TypeOf((*GrantSource)(nil)).Elem()
	if src.NumMethod() != 1 {
		t.Fatalf("tools.GrantSource has %d methods, want exactly 1 (read-only, like ModeSource)",
			src.NumMethod())
	}
	for _, act := range grantActs {
		if strings.Contains(strings.ToLower(src.Method(0).Name), act) {
			t.Errorf("tools.GrantSource's only method %q names the act %q: the read side of an "+
				"authorization must not create, remove or swap one", src.Method(0).Name, act)
		}
	}
}

func fieldNames(fs []reflect.StructField) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Name+"("+f.Type.String()+", kind "+f.Type.Kind().String()+")")
	}
	return out
}
