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
// TWO PROBES, BOTH REQUIRED TO GO RED IF THE COPY IS REMOVED:
//
//	probe 1  names every reference-typed slot of the returned Decision whose
//	         backing pointer is still the queue's. It walks by reflection, so a
//	         future reference field on tools.Decision that cloneDecision forgets
//	         is named by this case rather than discovered by an exploit; the
//	         slot list it expects is asserted as a set, so growing or shrinking
//	         Decision cannot pass silently either.
//	probe 2  is ticket 146 AC#2's literal assertion: write IN PLACE through every
//	         one of those slots, then read the record back out of the queue
//	         (its stored qitem.Dec, its own view()/head() projection, and a
//	         second LiveApprovals() call) and require it unchanged.
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
// tools.Decision filled, including the three that live inside Blacklist. A slot
// left nil here would make probe 1 blind to it, so the fixture is the census.
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

// refSlots lists every reference-typed slot of v that is non-nil, by dotted
// path. Structs are descended into (that is how Blacklist.Absolute is reached);
// map and slice ELEMENTS are not, because one-level copying is the promise
// pending_read.go makes.
func refSlots(prefix string, v reflect.Value, into *[]string) {
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		if !v.IsNil() {
			*into = append(*into, prefix)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if !f.IsExported() {
				continue
			}
			child := prefix + "." + f.Name
			if prefix == "" {
				child = f.Name
			}
			refSlots(child, v.Field(i), into)
		}
	}
}

// slotPaths runs the walker and hands back the sorted census of one value.
func slotPaths(v reflect.Value) []string {
	var found []string
	refSlots("", v, &found)
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

// queueSlotPaths is what the census in ticket 146 §1.2 measured on
// tools.Decision. It is asserted as a SET on both sides of the probe so that
// neither a field added to Decision nor a walker that quietly stopped
// descending can turn probe 1 into an always-green check.
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

	live := slotPaths(reflect.ValueOf(got[0].Decision))
	found := live
	expect := sorted(queueSlotPaths)
	if len(found) != len(expect) {
		t.Fatalf("返回值里的引用槽位数 %d (%v)，期望 %d (%v)。\n"+
			"数目不符说明 fixture 或 refSlots 自己失效了 —— 这条用例就会变成恒真判据，"+
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
