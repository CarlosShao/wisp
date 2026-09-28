package main

// internal/tools/subagent_197.go:43-47 promises exactly this file: the two legs
// of ticket 197 spell a subagent's stream key with the SAME literal
// ("subagent:<taskID>"), the writer being internal/tools (the spawner feeds it)
// and the reader being internal/panel (pump.go's StreamLog keys its rows by it).
// The two packages must not import each other, so neither one can pin the other's
// copy from inside itself - only the composition root, which imports both, can
// compare them. 197-r1 wrote that sentence and left the file out; this leg makes
// the sentence true (dispatch 2026-09-28 19:2x, §2(c) 丙).
//
// What this file deliberately does NOT do: merge the two literals into a single
// source of truth. Ruling A406 says the panel side is the real one and cmd/wisp
// injects it, and that consolidation is ticket 197's 载体层 r3, not this leg.
//
// Positive control taken at delivery: changing internal/tools' SubagentStreamKeyPrefix
// to another spelling turns the first case below red. Reading is in
// docs/evidence/s1/197-subagent-entity-r1b.md.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The prefix is one spelling, not two that happen to match today. A reader that
// looks up "subagentX:<id>" would find nothing while the writer keeps appending,
// and the panel would show a subagent that never produces a line.
func TestSubagentStreamKeyPrefixAgreesAcrossBothPackages(t *testing.T) {
	if tools.SubagentStreamKeyPrefix != panel.SubagentStreamKeyPrefix {
		t.Errorf("两包的前缀分家了：tools = %q, panel = %q（写侧键与读侧键从此对不上，"+
			"名册里的子代理会有流却没人读得到）",
			tools.SubagentStreamKeyPrefix, panel.SubagentStreamKeyPrefix)
	}
	// The literal itself is pinned too, so this leg cannot silently turn into
	// "both sides changed together": ticket 197 §0 froze this exact spelling.
	if tools.SubagentStreamKeyPrefix != "subagent:" {
		t.Errorf("前缀 = %q, want %q（§0 冻的键形状：全小写、冒号分隔）",
			tools.SubagentStreamKeyPrefix, "subagent:")
	}
}

// The two builders agree for the ids production actually mints. Compared against
// each other, not against a literal this file typed.
//
// Empty/whitespace ids are deliberately NOT compared: panel's SubagentStreamKey
// answers "" for them (Append("") would open a row nobody can attribute, its own
// comment at pump.go:382-389 says why) while the tools leg has no such branch.
// That divergence is documented, and closing it is 载体层 r3's call, not a
// composition-root test's.
func TestSubagentStreamKeyBuildersAgreeAcrossBothPackages(t *testing.T) {
	for _, id := range []string{
		"task-197-a",
		"00000000-0000-4000-8000-000000000000",
		"带中文的id",
	} {
		got, want := tools.SubagentStreamKey(id), panel.SubagentStreamKey(id)
		if got != want {
			t.Errorf("同一个 id %q 两侧拼出不同的键：tools = %q, panel = %q", id, got, want)
		}
	}
}
