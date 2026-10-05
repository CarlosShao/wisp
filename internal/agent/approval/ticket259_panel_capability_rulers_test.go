package approval

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// Ticket 259 AC#3 (leg 259-r1): the three capability rulers ticket 242 AC#2 was
// judged NOT MET over. 242 shipped a name-level fence - a field allowlist the
// panel-facing struct declares and its own test maintains - and the acceptance
// leg measured that a renamed field walks straight through it (M-D2: a
// "Permitted bool" plus its own name added to the list, 66 PASS), and that a
// real allow verb added to PanelAPI walks through the method side too (M-E: an
// Allow that actually spends a grant, 66 PASS).
//
// So these three rulers do not look at spellings at all. They ask the capability
// question three different ways:
//
//	1. the carrier's VERB set - what a panel host can make the process DO;
//	2. the payload's KIND and VALUE set - what shape can cross at all, and
//	   whether anything that crosses is a live proof or a state read-out;
//	3. the LAUNDERING direction - whether anything read out of the outward read
//	   face can be handed back as an answer.
//
// A ruler that only matches strings on names is not evidence here, and the two
// 242 rulers stay exactly as they were (this leg neither widens nor edits them):
// M-D2 is what proves the difference.
//
// Scope: rulers only. No frozen file is touched, no C17 inbound method name is
// added, and nothing here unfreezes a contract face - adding a name to the panel
// surface is a contract question (Q-76 is still with the owner) and the correct
// reaction of these rulers is RED, not a quieter gate.

// ---------------------------------------------------------------------------
// Ruler 1: the panel carrier's method set is closed, by capability
// ---------------------------------------------------------------------------

// panelAPIAllowedMethods is the whole of what a host holding Gate.Panel() may
// call. Reject is a refusal (always safe, SPEC-06 §9); Head and View are reads.
// Anything else is a new capability arriving on the untrusted surface and has to
// be argued for in the open, not slipped in as a fourth name.
var panelAPIAllowedMethods = []string{"Reject", "Head", "View"}

func r1MethodNames(typ reflect.Type) []string {
	out := make([]string, 0, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		out = append(out, typ.Method(i).Name)
	}
	return out
}

func r1MethodNameSet(names []string) map[string]bool {
	set := map[string]bool{}
	for _, n := range names {
		set[n] = true
	}
	return set
}

// TestTicket259R1PanelAPIMethodSetIsClosed is the ruler 242 M-E was planted
// against: a new method name on PanelAPI - any name, whatever it is called - goes
// red here, because the closed set is the assertion and the set has no room for
// a state-changing verb.
func TestTicket259R1PanelAPIMethodSetIsClosed(t *testing.T) {
	typ := reflect.TypeOf((*PanelAPI)(nil)).Elem()
	have := r1MethodNameSet(r1MethodNames(typ))
	allowed := r1MethodNameSet(panelAPIAllowedMethods)

	for _, n := range r1MethodNames(typ) {
		if !allowed[n] {
			t.Errorf("AC#3 RED: PanelAPI gained method %q; the panel surface is closed at Reject/Head/View, and a fourth name is a new capability on the untrusted face (ticket 242 M-E shape: an allow verb that spends a real grant). Refusing a card is safe; settling one as allowed is not, and no spelling of a new verb makes it the same act", n)
		}
	}
	for _, n := range panelAPIAllowedMethods {
		if !have[n] {
			t.Errorf("AC#3 RED: PanelAPI no longer carries %q; the read/refusal roster and its ruler drifted apart, and the fence below can no longer see what arrived", n)
		}
	}
	if typ.NumMethod() != len(panelAPIAllowedMethods) {
		t.Errorf("AC#3 RED: PanelAPI has %d methods, the closed set has %d; a count match is part of the fence, otherwise an alias for an existing verb passes silently",
			typ.NumMethod(), len(panelAPIAllowedMethods))
	}
}

// TestTicket259R1PanelCarrierCarriesNoAnswerVerb checks the half a closed
// interface list cannot see: the CONCRETE value handed out by Gate.Panel() can
// hold more methods than the interface declares, and any host that type-asserts
// (or is handed the concrete type by a future refactor) gets them. The probe is
// behavioural - "can this value answer a card" - not a name lookup, so an
// Allow-shaped verb earns its red whatever it is called.
func TestTicket259R1PanelCarrierCarriesNoAnswerVerb(t *testing.T) {
	g := New(Options{})
	carrier := g.Panel()
	typ := reflect.TypeOf(carrier)

	allowed := r1MethodNameSet(panelAPIAllowedMethods)
	for _, n := range r1MethodNames(typ) {
		if !allowed[n] {
			t.Errorf("AC#3 RED: the value Gate.Panel() returns (%s) carries method %q, which the PanelAPI list does not name; a panel host that asserts a wider interface calls it", typ, n)
		}
	}

	// The same question asked through Go's own type system, so a verb that only
	// exists on the concrete carrier is caught even if the roster above were
	// edited to match it.
	if _, ok := carrier.(interface {
		Allow(corr, grant string) error
	}); ok {
		t.Error("AC#3 RED: the panel carrier satisfies an Allow-shaped verb; a host holding Panel() can settle a card as allowed (ticket 242 M-E)")
	}
	if _, ok := carrier.(interface {
		AllowSession(corr, grant string) error
	}); ok {
		t.Error("AC#3 RED: the panel carrier satisfies a session-Allow-shaped verb; that verb lives on NativeAPI only")
	}
	if _, ok := carrier.(interface {
		Approve(corr string) error
	}); ok {
		t.Error("AC#3 RED: the panel carrier satisfies an Approve-shaped verb under a different spelling")
	}
	if _, ok := carrier.(interface {
		Settle(corr string, allow bool) error
	}); ok {
		t.Error("AC#3 RED: the panel carrier satisfies a Settle-shaped verb that carries the allow bit as an argument")
	}
}

// ---------------------------------------------------------------------------
// Ruler 2: the panel payload's fence is by capability, not by spelling
// ---------------------------------------------------------------------------

// r1PanelFieldIsDisplayData answers "can this field, by its KIND, carry a
// verdict, a verb, a handle into mutable state, or raw credential bytes". It is
// deliberately name-blind: M-D2 renamed nothing about the type, and a fence that
// only reads names is what that mutation measured as absent.
func r1PanelFieldIsDisplayData(f reflect.StructField) (string, bool) {
	switch f.Type.Kind() {
	case reflect.Bool:
		return "a bool on the card a panel reads IS a verdict bit (allowed / not allowed); the panel read face is plain display data and has no decisions to report", false
	case reflect.Func:
		return "a func field is a verb the host can call; changing a card's state through the read face is exactly what F2 layer 3 removes", false
	case reflect.Chan, reflect.Ptr, reflect.UnsafePointer:
		return "a live handle or a pointer aliases mutable state on the other side of the boundary; the read face must hand out values, not access", false
	case reflect.Interface:
		return "an interface field can hold anything, including a carrier with verbs on it; the fence cannot see past the dynamic type it was given", false
	case reflect.Map:
		return "a map field is an open-ended payload; the declared read face is a fixed set of display scalars, and an unbounded slot is where a grant gets smuggled", false
	case reflect.Array, reflect.Struct:
		return "a nested aggregate widens the read face by value rather than by name, so the fence has to look inside it; the panel projection is flat display data", false
	case reflect.Slice:
		if f.Type.Elem().Kind() != reflect.String {
			return "a slice of " + f.Type.Elem().Kind().String() + " carries raw bytes or aggregates, which is credential-shaped material (a nonce is bytes); the read face allows only a list of display strings", false
		}
	case reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Scalars the card renders (the depth badge and the strings it shows).
	default:
		return "field kind " + f.Type.Kind().String() + " is outside the display-scalar alphabet the panel read face is fenced to", false
	}
	return "", true
}

// TestTicket259R1PanelItemFieldCapabilityFence is the ruler M-D2 must trip: a
// field whose KIND can express a state verdict is refused whatever it is called,
// so adding it plus its own name to the string allowlist buys nothing.
func TestTicket259R1PanelItemFieldCapabilityFence(t *testing.T) {
	typ := reflect.TypeOf(PanelItem{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if why, ok := r1PanelFieldIsDisplayData(f); !ok {
			t.Errorf("AC#3 RED: PanelItem.%s has type %s - %s. This is the M-D2 shape: the field name is free, the list is updatable by the same hand that adds the field, so the fence has to be about what the field CAN carry (ticket 259 AC#3)", f.Name, f.Type, why)
		}
	}
}

// TestTicket259R1PanelItemCarriesNoCredentialValue asks the credential question
// about VALUES instead of names, on a card that holds a live grant: nothing the
// panel projection renders may be, or contain, the nonce or the binding digest.
// A field can keep an approved name and an approved type and still leak; only
// this kind of check sees that.
func TestTicket259R1PanelItemCarriesNoCredentialValue(t *testing.T) {
	g, _, it, nonce, _ := r1Card(t, "task-259-value", "corr-259-value", "shell.run")
	digest := it.bind

	for name, item := range map[string]PanelItem{"Head": mustPanelHead(t, g), "View": mustPanelView(t, g, it.Corr)} {
		typ := reflect.TypeOf(item)
		val := reflect.ValueOf(item)
		for i := 0; i < typ.NumField(); i++ {
			for _, v := range r1FieldStrings(val.Field(i)) {
				if v == "" {
					continue
				}
				// Only the field is named in the message: a grant value must not
				// land in a test log either.
				if v == nonce {
					t.Errorf("AC#3 RED: %s renders a live native grant in PanelItem.%s; the value a panel can read is the value that opens the card", name, typ.Field(i).Name)
				}
				if v == digest {
					t.Errorf("AC#3 RED: %s renders the item's binding digest in PanelItem.%s", name, typ.Field(i).Name)
				}
				if v != nonce && v != digest && strings.Contains(v, nonce) {
					t.Errorf("AC#3 RED: %s embeds a live native grant inside PanelItem.%s", name, typ.Field(i).Name)
				}
			}
		}
	}
}

// TestTicket259R1PanelItemProjectionIgnoresGrantState: the read face must render
// the same item whether or not its grant is live. A field that changes when the
// grant changes is reading decision state onto a surface that is only allowed to
// display it - which is the capability M-D2's field was standing in for.
func TestTicket259R1PanelItemProjectionIgnoresGrantState(t *testing.T) {
	g, q, it, _, _ := r1Card(t, "task-259-invariance", "corr-259-invariance", "shell.run")
	before := mustPanelView(t, g, it.Corr)
	q.revokeGrants(it.Corr)
	after := mustPanelView(t, g, it.Corr)

	if reflect.DeepEqual(before, after) {
		return
	}
	bt, av := reflect.TypeOf(before), reflect.ValueOf(after)
	bv := reflect.ValueOf(before)
	for i := 0; i < bt.NumField(); i++ {
		if !reflect.DeepEqual(bv.Field(i).Interface(), av.Field(i).Interface()) {
			t.Errorf("AC#3 RED: PanelItem.%s changed when the card's grants were revoked; the panel read face must not reflect grant or decision state (name-blind: this is the capability half of M-D2)", bt.Field(i).Name)
		}
	}
}

// r1FieldStrings flattens one field's observable strings (a plain string, or the
// elements of a []string) so the value checks do not hard-code field names.
func r1FieldStrings(v reflect.Value) []string {
	switch v.Kind() {
	case reflect.String:
		return []string{v.String()}
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.String {
			return nil
		}
		out := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			out = append(out, v.Index(i).String())
		}
		return out
	default:
		return nil
	}
}

func mustPanelHead(t *testing.T, g *Gate) PanelItem {
	t.Helper()
	item, ok := g.Panel().Head()
	if !ok {
		t.Fatal("AC#3 fixture: Panel().Head() found no card while one was pending")
	}
	return item
}

func mustPanelView(t *testing.T, g *Gate, corr string) PanelItem {
	t.Helper()
	item, ok := g.Panel().View(corr)
	if !ok {
		t.Fatalf("AC#3 fixture: Panel().View(%s) found no card while it was pending", corr)
	}
	return item
}

// ---------------------------------------------------------------------------
// Ruler 3: what the read face hands out cannot be handed back as an answer
// ---------------------------------------------------------------------------

// TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer is the semantic half
// the name fences cannot state at all: take EVERYTHING a panel host can read
// through Gate.Panel() - Head and View, every field, every element - and try to
// answer a live card with it, from both inbound routes. Nothing may land. The
// control that makes the reading mean something runs last: the card the read
// face was refused on must still be openable by the native grant, so a pass here
// cannot be the card having quietly died.
func TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer(t *testing.T) {
	g, q, x, nonceX, _ := r1Card(t, "task-259-launder-x", "corr-259-launder-x", "shell.run")
	y, _ := r1CardIn(t, q, "task-259-launder-y", "corr-259-launder-y", "fs.write")
	ctx := context.Background()

	for _, corr := range []string{x.Corr, y.Corr} {
		if _, ok := g.Panel().View(corr); !ok {
			t.Fatalf("AC#3 fixture: card %s is not readable while pending", corr)
		}
	}
	readOut := r1EverythingThePanelCanRead(g, x.Corr, y.Corr)
	if len(readOut) == 0 {
		t.Fatal("AC#3 RED: the read face yielded nothing to try, so this case asserted nothing")
	}

	// Native route: every panel-observable value, offered as the grant.
	for _, v := range readOut {
		if err := g.Native().Allow(ctx, x.Corr, v); err == nil {
			t.Errorf("AC#3 RED: a value read out of the panel face spent as a native grant and allowed the card; the outward read face and the answer face are not the same channel, and this is the case that says so")
			break
		}
	}
	// The control, and it has to be the last native call: the card is still live
	// and the minted grant still opens it.
	if err := g.Native().Allow(ctx, x.Corr, nonceX); err != nil {
		t.Fatalf("AC#3 RED: the panel-observable values were refused, but so was the live native grant (%v) - the card stopped being answerable before the control ran, which makes the pass above meaningless", err)
	}

	// Panel route: the same values, plus the real grant shape, offered where an
	// allow is refused on the route alone. Y stays pending through this: a grant
	// surfaced on the untrusted route burns itself, and answering is still a
	// refusal.
	for _, v := range readOut {
		if err := g.DecideFromPanel(ctx, Request{
			CorrelationID: y.Corr, Allow: true, Grant: v, Source: "native",
		}); err == nil {
			t.Errorf("AC#3 RED: DecideFromPanel accepted an allow built from panel-read material; the route is supposed to refuse before it looks at a grant")
			break
		}
	}
	if q.pendingCount() == 0 {
		t.Fatal("AC#3 RED: no card is pending after the laundering attempts, so an answer landed somewhere it should not have")
	}
}

// r1EverythingThePanelCanRead walks the entire outbound surface (Head plus View
// for each readable card) and returns every string a host can observe in it,
// including each element of every string slice. Field names are never consulted.
func r1EverythingThePanelCanRead(g *Gate, corrs ...string) []string {
	var out []string
	add := func(item PanelItem) {
		typ := reflect.TypeOf(item)
		val := reflect.ValueOf(item)
		for i := 0; i < typ.NumField(); i++ {
			out = append(out, r1FieldStrings(val.Field(i))...)
		}
	}
	if head, ok := g.Panel().Head(); ok {
		add(head)
	}
	for _, c := range corrs {
		if item, ok := g.Panel().View(c); ok {
			add(item)
		}
	}
	return out
}
