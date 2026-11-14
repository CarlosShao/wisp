package llm

// Ticket 216 AC#5's instrument for the seam's error classification. The
// criterion is "the class - and therefore the sentence internal/panel shows -
// is selected by a CODE, never by the upstream's English prose". The negative
// example is openchamber ChatInput.tsx:2205-2206 (three substring ORs on the
// error text, one of them "413"), transcribed verbatim in
// .scratch/wisp/probes/216/ac5/control/openchamber-chatinput-2205.txt.
//
// Teeth: the prose table below carries exactly the words the old implementation
// matched on ("model", "context length"). Put `strings.Contains(strings.ToLower
// (d.Message), "model")` back into ClassifyHTTPStatus and the 400/404 rows go
// red - the class flips to model for a body the CODE says nothing about. The
// second test is the positive control: the code token still selects, so the
// first test cannot be satisfied by collapsing every class to provider.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// messageProse216 is prose an upstream provider might send. None of it is a
// code: it must not change the class.
var messageProse216 = []string{
	"upstream 413 payload too large",
	"model not found: gpt-9",
	"the model `gpt-9` does not exist",
	"context length exceeded: 200001 > 128000",
	"insufficient_quota: your account is out of credits",
}

func TestTheUpstreamMessageTextNeverPicksTheClass216(t *testing.T) {
	statuses := []int{400, 404, 413, 422, 429}
	for _, status := range statuses {
		// The body carries no code: the class is the status's own answer.
		base, baseRetry := ClassifyHTTPStatus(HTTPErrorDetail{Status: status})
		for _, msg := range messageProse216 {
			got, retry := ClassifyHTTPStatus(HTTPErrorDetail{Status: status, Message: msg})
			if got != base {
				t.Errorf("status %d, message %q: class = %s, want %s - the prose picked the class (openchamber ChatInput.tsx:2205 shape)",
					status, msg, got, base)
			}
			if retry != baseRetry {
				t.Errorf("status %d, message %q: retryable = %v, want %v - the prose moved the retry decision too",
					status, msg, retry, baseRetry)
			}
		}
	}
}

// TestTheUpstreamCodeStillPicksTheClass216 is the positive control: a code is
// not prose, and the refinement the seam promises must still happen through
// the code token the adapters fill (openaichat/wire.go codeString, anthropic
// stream.go's code-or-type).
func TestTheUpstreamCodeStillPicksTheClass216(t *testing.T) {
	cases := []struct {
		status int
		code   string
		want   observe.ErrorClass
	}{
		{400, "model_not_found", observe.ClassModel},
		{422, "context_length_exceeded", observe.ClassModel},
		{404, "model_not_found", observe.ClassModel},
		{413, "context_length_exceeded", observe.ClassModel},
		{429, "insufficient_quota", observe.ClassBudget},
	}
	for _, c := range cases {
		// A contradicting message must not be able to veto the code.
		got, _ := ClassifyHTTPStatus(HTTPErrorDetail{
			Status: c.status, Code: c.code, Message: "unrelated prose",
		})
		if got != c.want {
			t.Errorf("status %d code %q: class = %s, want %s - the code branch lost to the prose",
				c.status, c.code, got, c.want)
		}
	}
}
