package main

// Ticket 21 segment 2: the native card wiring. The console surface is where a
// `wisp run` operator reads the card, so the cases below run a REAL
// approval.Gate through it - a Decision in, the printed card out, a native
// allow with the minted grant to settle - and assert the two disclosures the
// ticket says the card itself owes: FULL params (not a summary) and the D45-1
// batch view.
//
// The pre-existing printers these cases guard are pinned elsewhere too
// (approval_card_window_376_test.go on the window number, panel_pump_test.go
// on the cross-surface header/rules shape); this file adds no new pin, it
// proves the NEW lines arrive through the production plumbing rather than only
// through a hand-built Prompt.

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// syncBuffer is a mutex-guarded io.Writer so the gate's goroutine and the
// test's poll can share it without a data race (approval's own tests run the
// whole suite under -race; a bare bytes.Buffer here would trip it).
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

// capturingUI wraps the real console printer so a case can read the Prompt the
// gate produced AND the bytes the operator saw - the pair the ticket's AC#2 is
// about (delivered data is not yet a shown card).
type capturingUI struct {
	mu   sync.Mutex
	last approval.Prompt
	out  *syncBuffer
}

func (c *capturingUI) Prompt(ctx context.Context, p approval.Prompt) error {
	c.mu.Lock()
	c.last = p
	c.mu.Unlock()
	return (&consoleApprovalUI{out: c.out}).Prompt(ctx, p)
}

func (c *capturingUI) Update(ctx context.Context, e approval.Event) error {
	return (&consoleApprovalUI{out: c.out}).Update(ctx, e)
}

func (c *capturingUI) prompt(t *testing.T) approval.Prompt {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.last.CorrelationID == "" {
		t.Fatal("the gate produced no Prompt")
	}
	return c.last
}

func (c *capturingUI) text() string { return c.out.String() }

func (c *capturingUI) waitForCard(t *testing.T) approval.Prompt {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		got := c.last.CorrelationID != ""
		c.mu.Unlock()
		if got {
			return c.prompt(t)
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("no card reached the console surface")
	return approval.Prompt{}
}

// TestTicket21S2ConsoleL2CardShowsFullParamsAndDepth is AC#2's rendering half
// on the console face: an fs.write L2 raised by the real route must print the
// whole argument object and the queue depth, and a native allow with the
// card's grant must still be the only thing that releases it.
func TestTicket21S2ConsoleL2CardShowsFullParamsAndDepth(t *testing.T) {
	ui := &capturingUI{out: &syncBuffer{}}
	g := approval.New(approval.Options{UI: ui})
	revoke := g.AdmitTextTask("task-21s2-console")
	defer revoke()

	d := tools.Decision{
		Tool:          "fs.write",
		Params:        map[string]any{"content": "0123456789", "path": "C:/work/deep/dir/target.txt"},
		Args:          json.RawMessage(`{"content":"0123456789","path":"C:/work/deep/dir/target.txt"}`),
		Level:         risk.L2,
		RulesHit:      []risk.RuleID{risk.R2},
		Reason:        "覆盖已有文件越出 [fs] 授权目录",
		Paths:         []string{"C:/work/deep/dir/target.txt"},
		CorrelationID: "21s2-console-l2",
		TaskID:        "task-21s2-console",
	}

	done := make(chan tools.Answer, 1)
	go func() {
		a, _ := g.PendingApproval(context.Background(), d)
		done <- a
	}()

	p := ui.waitForCard(t)
	if !strings.Contains(p.Grant, approvalGrantPrefix21S2) {
		t.Fatalf("premise broke: the card carried no native grant (%q)", p.Grant)
	}
	text := ui.text()
	for _, want := range []string{
		"[确认 L2 fs.write]",
		"完整参数（非摘要）：",
		`"content": "0123456789"`, // json.Indent spacing
		`"path": "C:/work/deep/dir/target.txt"`,
		"规则 R2",
		"影响路径：C:/work/deep/dir/target.txt",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("printed card missing %q\nfull output:\n%s", want, text)
		}
	}

	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("native allow failed: %v", err)
	}
	select {
	case a := <-done:
		if a != tools.AnswerAllow {
			t.Errorf("a native allow answered %q, want allow", a)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the route did not return after a valid native allow")
	}
}

const approvalGrantPrefix21S2 = "grant_"

// TestTicket21S2ConsoleL1BatchCardIsAggregatedView proves the D45-1 lines
// reach the operator for the shape Aggregate exists for: >=3 homogeneous L1
// paths in one call print ONE confirm whose body names the total, the affected
// dirs, the first-5 preview and the expandable remainder.
func TestTicket21S2ConsoleL1BatchCardIsAggregatedView(t *testing.T) {
	paths := make([]string, 0, 7)
	for i := range 7 {
		paths = append(paths, "C:/batch/"+string(rune('a'+i))+"/f.txt")
	}
	p := approval.Prompt{
		CorrelationID: "21s2-batch",
		TaskID:        "task-21s2-batch",
		Tool:          "fs.write",
		Level:         "L1",
		Args:          json.RawMessage(`{"paths":["a","b","c","d","e","f","g"]}`),
		RulesHit:      []risk.RuleID{risk.R1},
		Reason:        "批量可逆写",
		Paths:         paths,
		Batch:         approval.Aggregate(l1Decision21S2(paths...)),
		Window:        3 * time.Second,
		Channels:      approval.NewChannels(approval.ChannelBall, approval.ChannelEsc).Statuses(),
	}
	if p.Batch == nil {
		t.Fatal("premise broke: 7 homogeneous L1 paths must aggregate")
	}
	out := &bytes.Buffer{}
	if err := (&consoleApprovalUI{out: out}).Prompt(context.Background(), p); err != nil {
		t.Fatalf("console Prompt: %v", err)
	}
	text := out.String()
	for _, want := range []string{
		"批量确认：共 7 项",
		"涉及目录：",
		"明细（前 5 条）：",
		"其余 2 条一并由本次确认覆盖",
		"取消方式：单击悬浮球（已加载）",
		"窗口 3.0s",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("batch card missing %q\nfull output:\n%s", want, text)
		}
	}
}

// l1Decision21S2 builds the Decision Aggregate reads for a homogeneous L1 batch.
func l1Decision21S2(paths ...string) tools.Decision {
	return tools.Decision{
		Tool:     "fs.write",
		Provider: tools.KindBuiltin,
		Params:   map[string]any{"paths": len(paths)},
		Level:    risk.L1,
		RulesHit: []risk.RuleID{risk.R1},
		Reason:   "声明 L1：可逆写",
		Paths:    paths,
	}
}
