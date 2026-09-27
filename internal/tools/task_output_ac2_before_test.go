package tools

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
)

// Ticket 164 AC#2 - the BEFORE reading, at CALL level.
//
// The question the ticket asks is not "is the string task.output in a registry"
// (that is the ruler ticket 164 AC#1 already ran, and a judgement that reuses
// it can never be falsified by an implementation). It is: "start something in
// the background, then read what it printed" - can that call be made to land
// today? So this probe sends one real request through the real C1/C3/C19/C25
// bridge (internal/tools/bridge.go Execute, the same choke point the loop
// calls) against the roster production registers today, which is
// tools.BuiltinFSEntries (cmd/wisp/run.go).
//
// What "before" looks like is asserted, not narrated: the call must fall into
// the Lookup failure branch at bridge.go:247-253 - an IsError outcome whose
// text names the tool as unknown, error class "tool" (D37: self-correctable,
// so the task keeps running), and NOT a Go error (a rejection is not a fault).
//
// This file deliberately references no task.output implementation, so it
// compiles and runs against the pre-implementation tree: the reading below is
// reproducible by
//
//	go test -count=1 -v -run TestTaskOutputAC2 ./internal/tools/
//
// The AFTER half of the comparison is TestTaskOutputAC2AfterLegIsReachable in
// task_output_leg_test.go, which issues the identical request against the
// roster that carries the tool.

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

	t.Logf("AC#2 BEFORE verbatim: IsError=%v ErrorClass=%q Truncated=%v Text=%q",
		out.IsError, out.ErrorClass, out.Truncated, out.Text)

	if !out.IsError {
		t.Fatalf("expected the pre-implementation call to be refused, got a success: %+v", out)
	}
	if !strings.Contains(out.Text, "未知工具") || !strings.Contains(out.Text, "task.output") {
		t.Fatalf("expected the unknown-tool branch (bridge.go:247), got: %q", out.Text)
	}
	if out.ErrorClass != "tool" {
		t.Errorf("error_class = %q, want tool (an unknown tool is self-correctable, D37)", out.ErrorClass)
	}
	if out.Text != "" && strings.Contains(out.Text, "bg-1") {
		t.Errorf("the refusal leaked the task id it never looked up: %q", out.Text)
	}
	// And the mirror of AC#1's claim, at call level: the directory the model is
	// shown today holds no task.* entry, so the model cannot even learn to ask.
	dir1 := tempCanonical(t)
	b2, _ := fsBridgeWith(t, nil, NoGate{}, dir1)
	tools, err := b2.Tools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	for _, ti := range tools {
		if strings.HasPrefix(ti.Name, "task.") {
			t.Fatalf("a task.* tool is registered: %s", ti.Name)
		}
	}
}
