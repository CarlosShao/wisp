package main

// Ticket 329 AC#3's ruler: the host receives the pump's marshal bytes BIT FOR BIT.
//
// This case exists to answer one question and nothing softer: when a publish
// really books a packet, are the bytes the stub host receives exactly the bytes
// the pump marshaled - not "was Eval called once", not "is the payload similar",
// not "is the length right". Three separate rulers run here, and each is the
// removal target of one mandated mutation in the evidence file:
//
//	marshal-side  booked bytes   ==  this file's own json.Marshal(booked snapshot)
//	push-side     host JS        ==  panelSnapshotEvalJS(those marshal bytes)
//	state-pin     the sampled packet reports the one card the queue really held
//
// The closing packet is pinned the same way, plus the durable ledger record's own
// bytes/sha256, so the whole chain push-bytes == booked-bytes == marshal-bytes ==
// ledger-digested-bytes is one asserted equation, not three vibes.
//
// The sample is taken on a deliberately quiescent window: the hook runs before
// the task loop exists, the only concurrent publisher is the gate's own Prompt,
// and the snapshotCount delta around this case's publish must be exactly 1 - a
// noisy window fails loudly instead of producing a flaky byte compare.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
)

// capturePanelHost329 is the stub host of this case: it records every JS sentence
// evalPanelJS was handed, in order. It decides nothing and rewrites nothing - the
// recorded strings are compared byte-for-byte by the rulers above.
type capturePanelHost329 struct {
	mu  sync.Mutex
	got []string
}

func (h *capturePanelHost329) evalPanelJS(js string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.got = append(h.got, js)
}

// snapshot returns the recorded pushes oldest-first, plus the most recent one
// ("" when nothing arrived - an empty host roster is itself a finding, not a pass).
func (h *capturePanelHost329) snapshot() (all []string, last string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	all = append([]string(nil), h.got...)
	if len(all) > 0 {
		last = all[len(all)-1]
	}
	return all, last
}

// marshalPacket329 is this file's OWN encode of a booked snapshot. It is
// intentionally not panel.Snapshot's bytes handed over by the transport: the
// point of the two byte rulers above is that the pump's encode, the booking, and
// the push agree with an encode this test performs itself.
func marshalPacket329(t *testing.T, snap panel.Snapshot) []byte {
	t.Helper()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("this case's own marshal failed, which means the snapshot value handed back is not encodable: %v", err)
	}
	return data
}

func TestBookedSnapshotReachesHostBitForBit329(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	host := &capturePanelHost329{}
	var (
		driven     *agentRuntime
		sampleSnap panel.Snapshot
		sampleData []byte
		sampleJS   string
		quietDelta int
	)
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		rt.setPanelSnapshotHost(host)
		revoke := rt.gate.AdmitTextTask("deliver-task")
		defer revoke()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan struct{})
		go func() {
			defer close(done)
			// An out-of-allowlist read is the same L2-into-the-queue driver the
			// existing pump cases use: queue non-empty is what makes the sample
			// packet's pending section real, which is what the state-pin reads.
			_, _ = rt.bridge.Execute(ctx, agent.ToolRequest{
				TaskID: "deliver-task", CorrelationID: "deliver-corr-329", CallID: "d1",
				Name: "fs.read", Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, "C:/Windows/win.ini")),
			})
		}()
		for i := 0; i < 300 && rt.ui.shown() == 0; i++ {
			time.Sleep(10 * time.Millisecond)
		}
		if rt.ui.shown() == 0 {
			t.Fatal("the gate never displayed a card, so the queue never held anything to report")
		}
		// Let the Prompt's own publish finish before sampling (see the header).
		time.Sleep(50 * time.Millisecond)
		before := rt.snapshotCount()
		rt.publishPanelSnapshot()
		quietDelta = rt.snapshotCount() - before
		sampleSnap, sampleData, _ = rt.lastPanelSnapshot()
		_, sampleJS = host.snapshot()
		cancel()
		<-done
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if driven == nil {
		t.Fatal("the hook never saw the assembled runtime")
	}
	if quietDelta != 1 {
		t.Fatalf("this case's publish moved the pump count by %d, want exactly 1: the sampling window was not quiescent, and a byte compare taken through a concurrent publisher proves nothing", quietDelta)
	}

	// Ruler 1 (marshal-side): booked bytes == this file's own encode of the booked
	// snapshot. Catches any drift between what the pump marshals and what it books.
	wantSample := marshalPacket329(t, sampleSnap)
	if !bytes.Equal(sampleData, wantSample) {
		t.Errorf("marshal-side ruler: the bytes the run booked are not bit-equal to json.Marshal of the snapshot it booked: booked=%q remarshal=%q", sampleData, wantSample)
	}
	// Ruler 2 (push-side): the stub host received exactly the named global call
	// wrapping THOSE bytes. This is the equation ticket 329 is about - the de-fang
	// mutation in the evidence file (push replaced by an empty call, or by the
	// ledger summary line) is red HERE and nowhere else.
	if wantJS := panelSnapshotEvalJS(wantSample); sampleJS != wantJS {
		t.Errorf("push-side ruler: the JS the stub host received for the sampled packet is not bit-equal to %s(the pump's marshal bytes): host=%q want=%q",
			panelSnapshotGlobal, sampleJS, wantJS)
	}
	// Ruler 3 (state-pin): the sampled packet really carries the one card the
	// queue held at that instant. Catches a marshal side that trims live state
	// consistently (the ticket's "少一条待批" example), where rulers 1-2 alone
	// would still see a self-consistent pair.
	if len(sampleSnap.Pending) != 1 || sampleSnap.Pending[0].CorrelationID != "deliver-corr-329" {
		t.Errorf("state-pin: sampled pending = %+v, want the one card this run's queue held (deliver-corr-329)", sampleSnap.Pending)
	}

	// The closing packet, same two byte rulers, plus the durable ledger record's
	// own digest: push bytes == booked bytes == marshal bytes == digested bytes.
	finalSnap, finalData, seen := driven.lastPanelSnapshot()
	if !seen {
		t.Fatal("this run published no packet at all")
	}
	wantFinal := marshalPacket329(t, finalSnap)
	if !bytes.Equal(finalData, wantFinal) {
		t.Errorf("closing packet: booked bytes are not bit-equal to json.Marshal of the booked snapshot: booked=%q remarshal=%q", finalData, wantFinal)
	}
	pushes, finalJS := host.snapshot()
	if wantJS := panelSnapshotEvalJS(wantFinal); finalJS != wantJS {
		t.Errorf("closing packet: the last JS the stub host received is not bit-equal to the marshal bytes of the last booked packet: host=%q want=%q", finalJS, wantJS)
	}
	if len(pushes) < 2 {
		t.Errorf("the stub host received %d pushes, want at least the run's own Prompt publish plus this case's sample", len(pushes))
	}
	recs := ledgerSummaries145(t, f.dir)
	rec := recs[len(recs)-1]
	if got, want := fmt.Sprint(len(finalData)), rec["bytes"]; got != want {
		t.Errorf("the last ledger record says %s bytes, the host-received closing packet is %s bytes long", want, got)
	}
	if rec["sha256"] != sha256Short(finalData) {
		t.Errorf("the last ledger record digests %q, the host-received closing packet digests %q", rec["sha256"], sha256Short(finalData))
	}
}
