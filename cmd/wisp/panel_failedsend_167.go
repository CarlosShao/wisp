package main

// Ticket 167 AC#5 - the failed-send slot: where the words of a send that did
// NOT go out wait for their second chance.
//
// WHAT THIS IS. The panel's message door (panel_message_453.go) is the one
// place a send can fail after the page handed the native side the user's own
// text: the relay is not bound (no task pipeline), or the task slot refused the
// submission. Until this ledger existed those words were gone - the page had
// already handed them over, the door answered an error sentence, and the only
// copy of the original was the page's own input row, which is exactly the copy
// the failure is most likely to destroy (the census row ④ in
// .scratch/wisp/issues/167-*.md §9: "页面把字发得进来，但发送失败后拿不回原文").
// So the door now records the verbatim original HERE, on the two branches where
// a send actually failed, and the snapshot pump carries it out under the
// composer section's own keys (failedSendText / failedSendKnown, registered as
// a judged page-side gap in internal/panel's composerTSGapLedger).
//
// THE ANTI-OVERWRITE RULE, AND WHERE ITS HALVES LIVE. The ticket's sentence is
// "用户后来新写的字不许被它覆盖". The Go half is the separation: the failed
// original lives under its OWN keys and the failure path touches no other
// composer field - there is no "replace the composer content" instruction on
// this wire, so a page that renders the slot as a recoverable chip cannot
// clobber the user's newer typing by accident of the data shape. The page half
// (how the restore is offered and consumed) is the UI leg's; the slot stays
// filled until a page-side consume route exists - a later success does NOT
// clear it (the positive control the dispatch names: 不失败时草稿不被抹), and a
// later failure REPLACES it with that failure's own text, never a merge.
//
// WHY NO SECOND CAP. The ledger keeps ONE string, verbatim, bounded by the
// bound the transport already applied to the request that carried it (the
// inbound envelope ceiling); truncating the original here would make the
// "recover your words" promise a lie in exactly the case it exists for. One
// slot, one string, the wire's own bound.

import (
	"sync"

	"github.com/CarlosShao/wisp/internal/panel"
)

// failedSendLedger is the one failed-send state a send face holds: the verbatim
// original of the most recent failed send, the request id it failed under (so a
// reader can tell WHICH attempt the words belong to), and how many failures
// this face has seen. The door writes it; the snapshot pump reads it through
// reading(). Nothing else touches it.
type failedSendLedger struct {
	mu        sync.Mutex
	text      string
	requestID string
	failures  int
}

// record books one failed send's original text, verbatim, replacing whatever
// the slot held (the newest failure is the current truth; the older one is not
// merged and not deleted into - there is one slot, on purpose).
func (l *failedSendLedger) record(requestID, text string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.text = text
	l.requestID = requestID
	l.failures++
}

// reading is the snapshot's view: the slot's text verbatim and Known=true,
// always - a ledger exists only on a door that has a send face, so "the face
// exists" is the only sentence this type can answer with. nil answers the zero
// reading (Known=false), which keeps both keys off the wire.
func (l *failedSendLedger) reading() panel.FailedSendReading {
	if l == nil {
		return panel.FailedSendReading{}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return panel.FailedSendReading{Text: l.text, Known: true}
}

// snapshot copies the slot's fields under its lock, for tests and forensics.
func (l *failedSendLedger) snapshot() (text, requestID string, failures int) {
	if l == nil {
		return "", "", 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.text, l.requestID, l.failures
}

// bindFailedSendLedger is the resident boot's one write of the runtime's
// failed-send field: the message door owns the ledger (it is where sends
// fail), the PIPELINE pump reads through the runtime, and this call crosses
// that ordering the same way bindAttachmentDrafts does - the door's ledger
// bound onto rt after both objects exist, so the pump's reader sees the slot
// the door writes, never a copy and never the nil the field held at assembly
// time. nil-safe everywhere; nil ledger leaves the field unbound.
func (rt *agentRuntime) bindFailedSendLedger(l *failedSendLedger) {
	if rt == nil || l == nil {
		return
	}
	rt.fsMu.Lock()
	rt.failedSend = l
	rt.fsMu.Unlock()
}

// failedSendReading is the pipeline pump's FailedSend reader (run.go's
// PumpSources literal): the bound door ledger's reading, or the zero reading
// when nothing is bound - a plain `wisp run` has no send door, so "this
// dimension does not exist here" is the honest answer and both composer keys
// stay off the wire (the omitempty rule).
func (rt *agentRuntime) failedSendReading() panel.FailedSendReading {
	if rt == nil {
		return panel.FailedSendReading{}
	}
	rt.fsMu.Lock()
	l := rt.failedSend
	rt.fsMu.Unlock()
	return l.reading()
}
