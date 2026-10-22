package tools

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
)

// Guard, not a historical reading: the shape of the unknown-tool refusal and
// its D37 class.
//
// This case is NOT evidence for ticket 164 AC#2 ("start something in the
// background, then read what it printed - can it be done today"). That reading
// was taken on the real tree at anchor f206e9f and lives in the ticket's AC#2
// note; this file cannot speak for it, because the bridge below is composed
// from BuiltinFSEntries, which is the fs family only, while the production
// composition root appends the task family separately (cmd/wisp/run.go,
// BuiltinTaskEntries). What this file pins is the bridge's own contract at the
// choke point every tool call crosses (Bridge.Execute):
//
//  1. REFUSAL SHAPE - an unknown name falls into the Lookup-failure branch and
//     comes back as DATA, not as a Go error: err == nil, because D37 says a
//     rejection is not a fault (SPEC-07 2 states the C3 rule as a refusal, not
//     as an error). The text names the requested tool as unknown and points the
//     model at list_tools, and it does not echo anything it never looked up.
//  2. D37 CLASS - that refusal carries class "tool" (observe.ClassTool:
//     self-correctable, so the task keeps running), never "internal".
//  3. FS ENTRY POINT - BuiltinFSEntries does not smuggle a task.* name into the
//     directory: the fs family stays single-purpose. This clause is about the
//     fs registration path only. It is NOT a claim about the production
//     registry, which does carry task.output today.
//
// The teeth were measured in both directions on 2026-10-22 (tree ce84d53d):
// rewording the refusal opener in bridge.go makes this red on the shape
// assertion; switching that branch's class to internal makes it red on the
// class assertion; restoring bridge.go (same sha256) makes it green. It pins
// the live behaviour of the fs-only bridge, so it is green on a tree where the
// feature has already landed (the measurement above was taken on one).
//
// Reproduce:
//
//	go test -count=1 -v -run '^TestTaskOutputAC2BeforeLegIsUnreachable$' ./internal/tools/
//
// The file name and the test name are historical keys - ticket 164 AC#2b, the
// archived anchor instrument .scratch/wisp/probes/164/v1/ac2-anchor-reading.sh
// and the CI logs address this case by them - so they are kept unchanged even
// though the case no longer signs itself as a before/after reading.

func TestTaskOutputAC2BeforeLegIsUnreachable(t *testing.T) {
	dir := tempCanonical(t)
	b, _ := fsBridgeWith(t, nil, NoGate{}, dir)

	req := agent.ToolRequest{
		TaskID: "164-ac2-before",
		CallID: "call-164-ac2-before",
		Name:   "task.output",
		Args:   json.RawMessage(`{"task_id":"bg-1"}`),
	}
	out, err := b.Execute(t.Context(), req)
	if err != nil {
		t.Fatalf("Execute must answer a rejection as data, not as a Go error (D37): %v", err)
	}

	t.Logf("unknown-tool refusal: IsError=%v ErrorClass=%q Truncated=%v Text=%q",
		out.IsError, out.ErrorClass, out.Truncated, out.Text)

	if !out.IsError {
		t.Fatalf("expected an unknown-tool call to be refused, got a success: %+v", out)
	}
	if !strings.Contains(out.Text, "未知工具") || !strings.Contains(out.Text, "task.output") {
		t.Fatalf("expected the unknown-tool refusal (Bridge.Execute, Lookup-failure branch), got: %q", out.Text)
	}
	if out.ErrorClass != "tool" {
		t.Errorf("error_class = %q, want tool (an unknown tool is self-correctable, D37)", out.ErrorClass)
	}
	if out.Text != "" && strings.Contains(out.Text, "bg-1") {
		t.Errorf("the refusal leaked the task id it never looked up: %q", out.Text)
	}
	// The fs entry point itself: the directory BuiltinFSEntries builds must not
	// carry a task.* name. Scope note: this is the fs registration path only;
	// the production composition root registers the task family on top of it.
	dir1 := tempCanonical(t)
	b2, _ := fsBridgeWith(t, nil, NoGate{}, dir1)
	tools, err := b2.Tools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	for _, ti := range tools {
		if strings.HasPrefix(ti.Name, "task.") {
			t.Fatalf("BuiltinFSEntries carried a task.* name: %s", ti.Name)
		}
	}
}
