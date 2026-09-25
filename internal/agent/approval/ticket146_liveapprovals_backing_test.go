package approval

// Ticket 146's guard for Queue.LiveApprovals (pending_read.go).
//
// WHY THIS CASE EXISTS. pending_read.go's header promised "copies of verdicts"
// while the implementation copied tools.Decision BY VALUE - and Decision itself
// carries reference-typed fields, so the returned row shared one map and seven
// backing arrays with the queue's own stored record. Writing a slot of a
// returned row therefore rewrote an admitted verdict. Nothing in the tree does
// that today (ticket 146 AC#1 measured zero such callers), which is exactly why
// the case had to be written: a discipline with no instrument is a comment, and
// the symptom of forgetting it once is a silently rewritten entry in the
// approval queue - the record a native L2 card is rendered from.
//
// THREE MECHANISMS, EACH REQUIRED TO GO RED FOR A DIFFERENT BREAKING SHAPE:
//
//	probe 1  names every reference-typed slot of the returned Decision whose
//	         backing pointer is still the queue's. It walks by DECLARED TYPE
//	         (reflect.TypeOf), not by value, so a ninth reference field on
//	         tools.Decision is counted even when the fixture leaves it nil, and
//	         a slot of a kind the census has never seen (pointer / interface /
//	         chan / func / unexported) FAILS CLOSED instead of being skipped
//	         silently. That is what makes "a future reference field that
//	         cloneDecision forgets is named by this case rather than discovered
//	         by an exploit" an instrument and not a promise. The slot list it
//	         expects is asserted as a set, so growing or shrinking Decision
//	         cannot pass silently either.
//	probe 2  is ticket 146 AC#2's literal assertion: write IN PLACE through every
//	         one of those slots, then read the record back out of the queue
//	         (its stored qitem.Dec, its own view()/head() projection, and a
//	         second LiveApprovals() call) and require it unchanged. BEFORE the
//	         write it also pins VALUE FIDELITY with reflect.DeepEqual: the
//	         returned Decision must equal the queue's stored record field for
//	         field. The backing-pointer check only proves the two do not SHARE
//	         memory; DeepEqual is what proves the copy did not drop the values
//	         on the way out (a cloneParamsMap that keeps keys and nils the
//	         values, or a cloneBacking that allocates an empty slice, shares
//	         nothing and so sails through probe 1 untouched).
//
// What is NOT claimed here: the copy is ONE level deep (see cloneDecision).
// The containers reachable *through* Params' values - the []any / map[string]any
// a JSON decode leaves inside it - are still shared, and this case does not
// probe them. That limit is stated in pending_read.go's header, not hidden.
//
// This file is `package approval` like pending_read_test.go (white box) because
// probe 2 has to read the queue's own stored record, and no public path exposes
// it - the whole point of the ticket is that the public path hands out an alias.

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// sharedDecision is one admitted verdict with EVERY reference-typed slot of
// tools.Decision filled, including the three that live inside Blacklist. The
// TYPE census (declaredRefSlots) no longer depends on the fixture to notice a
// slot, but probe 1's backing-pointer comparison still does: a slot left nil or
// empty here would make that comparison blind to it, so probe 1 asserts this
// fixture fills every slot the type census names. The fixture is not the census
// any more - it is the test of the census.
func sharedDecision() tools.Decision {
	return tools.Decision{
		CorrelationID:          "corr-backing",
		TaskID:                 "task-backing",
		Tool:                   "shell.run",
		Provider:               tools.KindBuiltin,
		Params:                 map[string]any{"argv": []any{"echo", "hi"}, "command": "echo hi"},
		Args:                   json.RawMessage(`{"argv":["echo","hi"],"command":"echo hi"}`),
		Level:                  risk.L2,
		RulesHit:               []risk.RuleID{risk.R2, risk.R3},
		Reason:                 "reason-backing",
		Paths:                  []string{"C:/dir/corr-backing"},
		Capabilities:           []tools.Capability{tools.CapShell},
		SessionOverrideBlocked: true,
		Mode:                   risk.ModeAskEveryStep,
		Blacklist: tools.BlacklistNote{
			Class:           risk.ClassB,
			Absolute:        []string{"C:/Windows/System32/config/SAM"},
			Unlockable:      []string{"C:/proj/.env"},
			AlreadyUnlocked: []string{"C:/proj/readme.md"},
			Reason:          "blacklist-reason-backing",
			Unresolved:      0,
		},
	}
}

// declaredRefSlots lists, by dotted path, every slot of type t whose DECLARED
// kind carries memory shared with the source: map, slice, pointer, chan, func,
// interface. Structs and arrays are descended into (that is how
// Blacklist.Absolute is reached, and how a reference hiding behind a struct or
// a fixed array is still caught); map and slice ELEMENTS are not, because
// one-level copying is the promise pending_read.go makes.
//
// It walks by TYPE, not by value, and it FAILS CLOSED: any kind outside the
// explicit value-safe list is either named or a t.Fatalf, never skipped. That
// is deliberate. The old value walker skipped pointer / interface / chan / func
// and unexported fields with no branch at all, and skipped any reference slot a
// fixture left nil - so a ninth reference field on tools.Decision could arrive
// quietly. Here it cannot: a newly declared map/slice/pointer/... is counted
// (census set grows -> mismatch -> red), and a newly declared kind nobody has
// reasoned about stops the test with an error naming the kind.
//
// Unexported struct fields are descended too, even though reflect can't read
// their values here: cloneDecision lives in package approval and cannot assign
// a tools.Decision private field anyway, so a private reference slot is exactly
// the shape that must be surfaced rather than walked past.
func declaredRefSlots(t *testing.T, prefix string, ty reflect.Type, into *[]string) {
	t.Helper()
	switch ty.Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface:
		*into = append(*into, prefix)
	case reflect.Struct:
		for i := 0; i < ty.NumField(); i++ {
			f := ty.Field(i)
			child := f.Name
			if prefix != "" {
				child = prefix + "." + f.Name
			}
			declaredRefSlots(t, child, f.Type, into)
		}
	case reflect.Array:
		declaredRefSlots(t, prefix+"[]", ty.Elem(), into)
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		// Copied by value into the returned struct; nothing shared, nothing to name.
	default:
		t.Fatalf("引用槽位普查(%s): 未预期的 reflect.Kind %s —— 普查必须显式命名或下降它，"+
			"不许静默跳过；这是 fail-closed，新增字段类型要连这把尺一起改。", prefix, ty.Kind())
	}
}

// typeCensus runs the walker over a TYPE and hands back the sorted census of
// every declared reference-typed slot. It is the census probe 1 asserts against
// the hand-maintained set below.
func typeCensus(t *testing.T, ty reflect.Type) []string {
	t.Helper()
	var found []string
	declaredRefSlots(t, "", ty, &found)
	return sorted(found)
}

// backingPointer is the identity of one slot's shared memory: the data pointer
// of a slice, the runtime pointer of a map. Zero means "not a slot this case
// measures".
func backingPointer(v reflect.Value) uintptr {
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		return v.Pointer()
	}
	return 0
}

func slotValue(t *testing.T, root reflect.Value, path string) reflect.Value {
	t.Helper()
	cur := root
	for _, p := range strings.Split(path, ".") {
		if cur.Kind() != reflect.Struct {
			t.Fatalf("slot path %q: %q is not a struct, cannot take field %q", path, cur.Kind(), p)
		}
		next := cur.FieldByName(p)
		if !next.IsValid() {
			t.Fatalf("slot path %q: no field %q", path, p)
		}
		cur = next
	}
	if k := cur.Kind(); k != reflect.Map && k != reflect.Slice {
		t.Fatalf("slot path %q ends on a %s, not a map or slice", path, k)
	}
	return cur
}

// queueSlotPaths is the hand-maintained census of tools.Decision's reference
// slots (ticket 146 §1.2). probe 1 compares it, as a SET, against the census
// the TYPE walker (declaredRefSlots) derives live from reflect.TypeOf. The
// whole point of deriving it by type is that neither of the two ways this used
// to be an always-green check survives: a field ADDED to Decision now grows the
// derived set and trips the mismatch (the old value walker missed it whenever
// the fixture left it nil), and a kind the walker has no branch for stops the
// test with an error instead of being descended past. When Decision genuinely
// gains a reference slot, this list and cloneDecision must change together, in
// the same review - which is the review the mismatch forces.
var queueSlotPaths = []string{
	"Args",
	"Blacklist.Absolute",
	"Blacklist.AlreadyUnlocked",
	"Blacklist.Unlockable",
	"Capabilities",
	"Params",
	"Paths",
	"RulesHit",
}

func sorted(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// TestLiveApprovalsSharesNoReferenceSlotWithTheQueue is probe 1.
func TestLiveApprovalsSharesNoReferenceSlotWithTheQueue(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	want := sharedDecision()
	mustPush(t, q, want)

	got := q.LiveApprovals()
	if len(got) != 1 {
		t.Fatalf("期望 1 行，得 %d：%+v", len(got), got)
	}

	live := typeCensus(t, reflect.TypeOf(got[0].Decision))
	found := live
	expect := sorted(queueSlotPaths)
	if len(found) != len(expect) {
		t.Fatalf("返回值里的引用槽位数 %d (%v)，期望 %d (%v)。\n"+
			"数目不符说明 tools.Decision 增删了引用字段而 queueSlotPaths / cloneDecision 没跟上 —— 这条用例正是为此而红，"+
			"必须和 tools.Decision 的字段表一起改，不许放宽。",
			len(found), found, len(expect), expect)
	}
	for i := range expect {
		if found[i] != expect[i] {
			t.Fatalf("引用槽位名册不符：第 %d 枚 got=%q want=%q（全量 got=%v）", i, found[i], expect[i], found)
		}
	}

	stored := reflect.ValueOf(mustFindLive(q, "corr-backing").Dec)
	var shared []string
	for _, path := range found {
		a := slotValue(t, reflect.ValueOf(got[0].Decision), path)
		b := slotValue(t, stored, path)
		// The backing-pointer comparison below only has teeth on a non-empty
		// slot: a nil or zero-length slice/map has a zero (or shared-empty)
		// data pointer, so "pa != 0 && pa == pb" would report no alias even if
		// the slot were handed back straight from the queue. Refuse to run the
		// comparison blind - if the fixture stops filling a declared slot, this
		// goes red instead of quietly weakening probe 1.
		if b.IsNil() || b.Len() == 0 {
			t.Fatalf("fixture 未填引用槽 %q（nil 或空）：backing 指针对它是失明的，"+
				"probe 1 会退化成恒真。sharedDecision 必须填满 typeCensus 点名的每一枚槽位。", path)
		}
		if pa, pb := backingPointer(a), backingPointer(b); pa != 0 && pa == pb {
			shared = append(shared, path)
		}
	}
	if len(shared) != 0 {
		t.Errorf("AC#2 RED: LiveApprovals() 返回行里这些槽位与队列存储的那条记录共用底层：%v。\n"+
			"pending_read.go 的注释说 copies，共用底层就是就地写会改掉已决断记录。\n"+
			"修法在 cloneDecision：给每一枚点名的槽位单独分配。", shared)
	}
}

// TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord is probe 2, AC#2's
// literal sentence: write in place through every reference slot of the returned
// row, then read the record back out of the queue three ways and require it
// unchanged. Under the pre-ticket-146 implementation all eight writes land in
// the queue's own memory and every assertion below goes red.
func TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	want := sharedDecision()
	mustPush(t, q, want)

	got := q.LiveApprovals()
	if len(got) != 1 {
		t.Fatalf("期望 1 行，得 %d：%+v", len(got), got)
	}
	row := got[0]

	// VALUE FIDELITY, before anything is clobbered. "copies" is a two-sided
	// promise: not only must the returned row not SHARE the queue's memory
	// (probe 1's backing-pointer check), the copy must also CARRY the same
	// values. A cloneParamsMap that keeps every key and nils every value, or a
	// cloneBacking that allocates a fresh-but-empty slice, shares nothing - so
	// probe 1 stays green - yet hands the panel a deformed verdict. DeepEqual on
	// the returned Decision vs the queue's stored record is the one ruler that
	// goes red for that shape, and it is checked here, before the in-place write
	// below, so a genuine clobber is never mistaken for a lost value.
	if !reflect.DeepEqual(row.Decision, mustFindLive(q, "corr-backing").Dec) {
		t.Fatalf("AC#2 RED (值保真): LiveApprovals() 交出的 Decision 与队列存储的那条不逐字段相等。\n"+
			"返回值:  %+v\n存储项: %+v\n"+
			"这就是注释里 \"copies\" 的第二半没人钉的形状：不共用底层，却把值拷丢了。\n"+
			"修法在 cloneDecision／cloneBacking／cloneParamsMap：分配的每一份都得把原值带上。",
			row.Decision, mustFindLive(q, "corr-backing").Dec)
	}

	// Clobber each slot IN PLACE, through the returned value only. No
	// re-assignment of whole fields: that would prove nothing about aliasing.
	row.Decision.Params["argv"] = "clobbered"
	row.Decision.Params["injected-by-caller"] = true
	copy(row.Decision.Args, []byte("XXXXXXXX"))
	row.Decision.RulesHit[0] = risk.RuleID("R9")
	row.Decision.Paths[0] = "C:/clobbered"
	row.Decision.Capabilities[0] = tools.CapNotify
	row.Decision.Blacklist.Absolute[0] = "C:/clobbered"
	row.Decision.Blacklist.Unlockable[0] = "C:/clobbered"
	row.Decision.Blacklist.AlreadyUnlocked[0] = "C:/clobbered"

	// (1) The queue's own stored record, read white-box.
	have := mustFindLive(q, "corr-backing").Dec
	if _, injected := have.Params["injected-by-caller"]; injected {
		t.Errorf("AC#2 RED: 队列存的 Params 里出现了调用方塞进的键 %q，map 是共用的（键 argv 现值 %v）",
			"injected-by-caller", have.Params["argv"])
	}
	if argv, ok := have.Params["argv"].(string); ok && argv == "clobbered" {
		t.Errorf("AC#2 RED: 队列存的 Params[\"argv\"] 被就地改成 %q，原本应是解出来的参数列表", argv)
	}
	if string(have.Args[:8]) == "XXXXXXXX" {
		t.Errorf("AC#2 RED: 队列存的 Args 前 8 字节被就地改成 %q", string(have.Args[:8]))
	}
	if have.RulesHit[0] != risk.R2 {
		t.Errorf("AC#2 RED: 队列存的 RulesHit[0]=%q，期望仍是 %q —— L2 卡逐字念这枚（票 17 冻结注）",
			have.RulesHit[0], risk.R2)
	}
	if have.Paths[0] != "C:/dir/corr-backing" {
		t.Errorf("AC#2 RED: 队列存的 Paths[0]=%q，期望原样", have.Paths[0])
	}
	if have.Capabilities[0] != tools.CapShell {
		t.Errorf("AC#2 RED: 队列存的 Capabilities[0]=%q，期望原样 %q —— 能力位被改写就是权限面被改写",
			have.Capabilities[0], tools.CapShell)
	}
	for name, slot := range map[string]string{
		"Blacklist.Absolute":        have.Blacklist.Absolute[0],
		"Blacklist.Unlockable":      have.Blacklist.Unlockable[0],
		"Blacklist.AlreadyUnlocked": have.Blacklist.AlreadyUnlocked[0],
	} {
		if slot == "C:/clobbered" {
			t.Errorf("AC#2 RED: 队列存的 %s 被就地改成 %q —— A/B 档证据被改，卡片据此说话", name, slot)
		}
	}

	// (2) The queue's own untrusted projection, through its public read path.
	// PanelItem already copies Paths (queue.go's viewLocked), so this is the
	// second ruler: if a future edit drops that copy, the card shows the
	// caller's rewrite.
	item, ok := q.view("corr-backing")
	if !ok {
		t.Fatalf("corr-backing 不在队列里了：一次只读把队列条目弄丢了")
	}
	if len(item.Paths) != 1 || item.Paths[0] != "C:/dir/corr-backing" {
		t.Errorf("AC#2 RED: PanelItem.Paths=%v，期望队列投影响原样（面板卡片的越界提示读这一枚）", item.Paths)
	}
	if item.Reason != "reason-backing" || item.Tool != "shell.run" {
		t.Errorf("AC#2 RED: PanelItem 非引用字段也被牵动：tool=%q reason=%q", item.Tool, item.Reason)
	}
	head, ok := q.head()
	if !ok {
		t.Fatal("队头读不出来了")
	}
	if head.Depth != item.Depth {
		t.Errorf("view 与 head 对同一枚项报出不同位次 %d vs %d", item.Depth, head.Depth)
	}

	// (3) A second read: the pump re-reads on every state change, so the row it
	// gets next must describe the queue, not the earlier caller's write.
	again := q.LiveApprovals()
	if len(again) != 1 {
		t.Fatalf("第二次读期望 1 行，得 %d", len(again))
	}
	if again[0].Decision.RulesHit[0] != risk.R2 {
		t.Errorf("AC#2 RED: 第二次读到 RulesHit[0]=%q，队列记录已被上一枚读者改掉", again[0].Decision.RulesHit[0])
	}
	if again[0].Decision.Capabilities[0] != tools.CapShell {
		t.Errorf("AC#2 RED: 第二次读到 Capabilities[0]=%q", again[0].Decision.Capabilities[0])
	}
	if _, injected := again[0].Decision.Params["injected-by-caller"]; injected {
		t.Errorf("AC#2 RED: 第二次读到的 Params 里带着别人塞的键")
	}

	// And the mutation must not have cost the queue an item or a state change.
	if q.Depth() != 1 {
		t.Errorf("Depth()=%d，期望 1：一次读加就地写不应动队列结构", q.Depth())
	}
}
