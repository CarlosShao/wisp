package memory

// Ticket 29 (29-r2): the L1 extraction goroutine must keep the frozen D38b
// roster name, and that name must be a temporary-roster member.
//
// WHY THIS TEST EXISTS. internal/memory/l1_warm.go:221 calls
// reg.Spawn(l1ExtractGoroutine, ...) - a package const, not a string literal -
// so the repo-wide 229 census (internal/tools/spawn_roster_census_229_test.go)
// cannot statically classify it and books the site by hand in
// spawn229ReviewedSites, recording "memory-extract => temporary" as prose.
// Prose in another package's ledger cannot notice a rename: changing the const
// value would leave the census green while the ledger's sentence became a lie,
// and every spawn would start crying "goroutine outside the D38 roster".
// This test is the other end of that handshake, on the package that owns the
// const: rename the const to anything off-roster and this goes red.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

func TestTicket29_L1ExtractGoroutineIsARosterTemporaryName(t *testing.T) {
	if l1ExtractGoroutine != "memory-extract" {
		t.Fatalf("L1 extraction goroutine name = %q；D38b 临时族里的名字是 memory-extract，"+
			"改它就得同时改 spawn229ReviewedSites 的登记（票 29 29-r2）", l1ExtractGoroutine)
	}
	if cat := observe.ClassifyGoroutine(l1ExtractGoroutine); cat != observe.CategoryTemporary {
		t.Fatalf("ClassifyGoroutine(%q) = %s，不是 temporary：名册外名字每次 spawn 都会让"+
			"「goroutine outside the D38 roster (leak symptom)」响一次", l1ExtractGoroutine, cat)
	}
	// 正控：分类器没有变成"什么名字都放行"——真名册外的名字仍是 unknown，
	// 那才是泄漏报警唯一的那一发信号。
	if cat := observe.ClassifyGoroutine("29r2-not-on-any-roster-table"); cat != observe.CategoryUnknown {
		t.Fatalf("正控失败：名册外名字被判成 %s，泄漏报警的判据没了", cat)
	}
}
