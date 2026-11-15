package main

// Ticket 451 - the L1 profile injection's acceptance box.
//
// THE DEFECT THIS CASE WAS WRITTEN AGAINST (measured before the wiring existed):
// `agent.Config.Profiles` had ZERO assignments in production code. The profile
// rows were really written by the L1 extractor (`internal/memory`'s
// UpsertProfile), the prompt assembler really had a (3) profile section with a
// D15 budget (`internal/agent/prompt.go`'s profileText), the loop really
// forwarded `Config.Profiles` into `PromptInput.Profiles` (loop.go's
// buildRequest) - and nobody ever put the rows on that field. So the promise
// "remembers a fact from round 1" held in the database and nowhere else: the
// model never saw the fact, and the user saw "it wrote it down, but it does not
// mention it".
//
// WHY THIS CASE IS THE ONE THAT CAN GO RED: it drives the REAL composition root
// (`runTextTask`, the same function `wisp run` calls: config.toml -> secret
// store -> provider -> approval gate -> tools bridge -> agent loop) against the
// REAL mockllm, and asserts on the request body mockllm CAPTURED - not on a
// struct this test built. The profile is seeded through the product's own DAO
// into the same data dir the run opens, so the only path from the row to the
// request is the production wiring this ticket is about. Delete that wiring and
// the captured prompt loses the section: that is this ticket's mutation M1.
//
// WHAT THE ASSERTIONS READ, AND WHY THEY ARE ASCII: the fact and the section
// marker are pure ASCII ("mother birthday 05-06", "## profile"), so a JSON
// transport that escapes or re-encodes non-ASCII cannot make the reading depend
// on the encoder. The section marker is the context message's own "## profile"
// heading (prompt.go's ContextMessage), which only exists when the profile
// section is non-empty.

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// t451ProfileFact is the L1 row content this case seeds. It is a plain fact
// ("第 1 轮说过的一个事实" in the product's words), distinct enough that finding
// it in the request body cannot be an accident of the user's own utterance -
// the task text below never contains it.
const t451ProfileFact = "mother birthday 05-06"

// Test451L1ProfileReachesTheProviderPrompt is AC#1: with an L1 profile row on
// disk, one real conversation round must send a prompt that carries that fact.
func Test451L1ProfileReachesTheProviderPrompt(t *testing.T) {
	f := newRunFixture(t, "openai-chat")

	// Seed the row BEFORE the run, through the same DAO production uses, into
	// the same data dir assembleRuntime opens (memory.Open(s.dataDir)).
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer seedCancel()
	seed, err := memory.Open(f.dir)
	if err != nil {
		t.Fatalf("seed store: %v", err)
	}
	if _, err := seed.UpsertProfile(seedCtx, "fact.family", t451ProfileFact, memory.ProfileManual); err != nil {
		seed.Close()
		t.Fatalf("seed L1 profile row: %v", err)
	}
	if err := seed.Close(); err != nil {
		t.Fatalf("close seed store: %v", err)
	}

	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}

	// The request the provider ACTUALLY received, captured by mockllm.
	body, _ := f.srv.LastRequest(t, "chat")
	marker := strings.Index(body, "## profile")
	if marker < 0 {
		t.Fatalf("送给 provider 的 prompt 里没有画像段（找不到 \"## profile\"）：库里那枚 L1 画像 %q 没有被递进请求，模型看不到这条事实。mockllm 捕获的请求体：\n%s",
			t451ProfileFact, body)
	}
	if !strings.Contains(body[marker:], t451ProfileFact) {
		t.Fatalf("prompt 里有画像段，但段内没有库里那枚 L1 画像事实 %q：画像写进了库，却没有进 prompt。mockllm 捕获的请求体：\n%s",
			t451ProfileFact, body)
	}
}

// Test451L1ProfileInjectionStaysInsideTheD15Budget is AC#3's boundary reading:
// the same production path must not bypass the D15 (3) section budget. The
// store's cap is 20 rows; each seeded row is long enough that all twenty cannot
// fit the section budget, so a bare concatenation would put every row in the
// request while the budgeted path drops the tail. The assertion is the shape,
// not a byte count: the section is present (so the wiring really fed it) AND at
// least one seeded row is absent (so the budget really bit).
func Test451L1ProfileInjectionStaysInsideTheD15Budget(t *testing.T) {
	f := newRunFixture(t, "openai-chat")

	seedCtx, seedCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer seedCancel()
	seed, err := memory.Open(f.dir)
	if err != nil {
		t.Fatalf("seed store: %v", err)
	}
	// 20 rows is the DAO cap. Each line is ~140 ASCII bytes, so the 20 lines
	// plus the header are ~2800 bytes (~700 tokens of the shared heuristic) -
	// well past the (3) profile budget's 400-token grant at the reference
	// window, and a bare concatenation would therefore be visible as 20 rows.
	const budgetRowPrefix = "t451-budget-row-"
	padding := strings.Repeat("x", 100)
	for i := 0; i < memory.L1MaxProfiles; i++ {
		content := fmt.Sprintf("%s%02d-%s", budgetRowPrefix, i, padding)
		slot := fmt.Sprintf("fact.budget%02d", i)
		if _, err := seed.UpsertProfile(seedCtx, slot, content, memory.ProfileManual); err != nil {
			seed.Close()
			t.Fatalf("seed L1 row %d: %v", i, err)
		}
	}
	if err := seed.Close(); err != nil {
		t.Fatalf("close seed store: %v", err)
	}

	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	body, _ := f.srv.LastRequest(t, "chat")
	marker := strings.Index(body, "## profile")
	if marker < 0 {
		t.Fatalf("送给 provider 的 prompt 里没有画像段：画像根本没有接线，预算无从谈起。请求体：\n%s", body)
	}
	section := body[marker:]
	kept := strings.Count(section, budgetRowPrefix)
	t.Logf("预算路径读数：画像段保留 %d/%d 条（D15 (3) 段预算把它裁到预算内；裸拼接会是 20 条）",
		kept, memory.L1MaxProfiles)
	if kept == 0 {
		t.Fatalf("画像段里一条预算内的画像都没有：预算路径把整段都丢光了。请求体：\n%s", body)
	}
	if kept >= memory.L1MaxProfiles {
		t.Fatalf("画像段把 %d 条全部塞进了 prompt（库上限 %d）：这等于裸拼接、绕过了 D15 预算。请求体：\n%s",
			kept, memory.L1MaxProfiles, body)
	}
}
