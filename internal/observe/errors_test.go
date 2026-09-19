package observe

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestAllClassesHasExactly17 pins the D37(a) contract: exactly 17 classes,
// unique, in the documented table order, all valid.
func TestAllClassesHasExactly17(t *testing.T) {
	classes := AllClasses()
	if len(classes) != 17 {
		t.Fatalf("AllClasses() len = %d, want 17", len(classes))
	}
	seen := map[ErrorClass]bool{}
	for i, c := range classes {
		if !c.Valid() {
			t.Fatalf("class %d (%q) not valid", i, c)
		}
		if seen[c] {
			t.Fatalf("duplicate class %q", c)
		}
		seen[c] = true
	}
	want := []ErrorClass{
		ClassConfig, ClassAuth, ClassNetwork, ClassRateLimit, ClassProvider,
		ClassModel, ClassAudioDevice, ClassASR, ClassTool, ClassPermissionDenied,
		ClassUserRejected, ClassCancelled, ClassBudget, ClassLoop, ClassInjection,
		ClassInternal, ClassResource,
	}
	for i := range want {
		if classes[i] != want[i] {
			t.Fatalf("class %d = %q, want %q", i, classes[i], want[i])
		}
	}
	// Mutating the returned copy must not affect the package state.
	classes[0] = "tampered"
	if AllClasses()[0] != ClassConfig {
		t.Fatal("AllClasses() returned a slice aliasing internal state")
	}
}

// TestClassMappingAndRetry covers all 17 classes one row each, asserting the
// D37 mapped-state and retryability columns exactly. (Per SPEC-01 §7 the
// end-to-end classification tests with real error sources - bad config,
// 429 mock, broken model file - land with the tickets that own those sources;
// this test pins the model itself.)
func TestClassMappingAndRetry(t *testing.T) {
	cases := []struct {
		class  ErrorClass
		states []string
		policy RetryPolicy
	}{
		{ClassConfig, []string{"Unconfigured"}, RetryNever},
		{ClassAuth, []string{"Unconfigured", "Error"}, RetryNever},
		{ClassNetwork, []string{"NoNetwork"}, RetryBackoff},
		{ClassRateLimit, []string{"Error"}, RetryAfterHeader},
		{ClassProvider, []string{"Error"}, RetryBackoff},
		{ClassModel, []string{"Downloading", "Error"}, RetryNever},
		{ClassAudioDevice, []string{"Error"}, RetryOnce},
		{ClassASR, []string{"Listening", "Error"}, RetryOnce},
		{ClassTool, []string{"Acting"}, RetryAgentDecided},
		{ClassPermissionDenied, []string{"Acting"}, RetryNever},
		{ClassUserRejected, []string{"Settling"}, RetryNever},
		{ClassCancelled, []string{"Settling"}, RetryNever},
		{ClassBudget, []string{"Stuck"}, RetryNever},
		{ClassLoop, []string{"Stuck"}, RetryNever},
		{ClassInjection, []string{"Error"}, RetryNever},
		{ClassInternal, []string{"Error"}, RetryNever},
		{ClassResource, []string{"WatchdogAlert"}, RetryNever},
	}
	if len(cases) != 17 {
		t.Fatalf("test table has %d rows, want 17", len(cases))
	}
	for _, tc := range cases {
		if got := tc.class.MappedStates(); !equalStrings(got, tc.states) {
			t.Errorf("%s.MappedStates() = %v, want %v", tc.class, got, tc.states)
		}
		if got := tc.class.RetryPolicy(); got != tc.policy {
			t.Errorf("%s.RetryPolicy() = %d, want %d", tc.class, got, tc.policy)
		}
		wantRetryable := tc.policy != RetryNever
		if got := tc.class.Retryable(); got != wantRetryable {
			t.Errorf("%s.Retryable() = %v, want %v", tc.class, got, wantRetryable)
		}
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestMappedStatesAreValidBallStates pins observe's state strings to the
// statemachine vocabulary so the two frozen tables cannot drift.
func TestMappedStatesAreValidBallStates(t *testing.T) {
	for _, c := range AllClasses() {
		for _, s := range c.MappedStates() {
			if !statemachine.Valid(statemachine.State(s)) {
				t.Errorf("class %s maps to %q which is not one of the 20 D43 states", c, s)
			}
		}
	}
}

// TestErrorTypeAndClassOf covers construction, wrapping, errors.As/Is and the
// internal fallback.
func TestErrorTypeAndClassOf(t *testing.T) {
	cause := errors.New("dial tcp: connection refused")
	err := Wrap(ClassNetwork, cause, "provider unreachable")

	if got := err.Error(); got != "network: provider unreachable: dial tcp: connection refused" {
		t.Fatalf("Error() = %q", got)
	}
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is failed through the wrap chain")
	}
	var oe *Error
	if !errors.As(err, &oe) || oe.Class != ClassNetwork {
		t.Fatalf("errors.As failed: %+v", oe)
	}
	if c, ok := ClassOf(err); !ok || c != ClassNetwork {
		t.Fatalf("ClassOf = %q,%v", c, ok)
	}
	if got := ClassOfOrInternal(cause); got != ClassInternal {
		t.Fatalf("ClassOfOrInternal(unclassified) = %q, want internal", got)
	}
	if _, ok := ClassOf(cause); ok {
		t.Fatal("ClassOf(unclassified) reported ok")
	}

	// Retryable semantics.
	rl := &Error{Class: ClassRateLimit, RetryAfter: 30 * time.Second, ProviderCode: "429"}
	if !rl.Retryable() {
		t.Fatal("rate_limit with retry-after must be retryable")
	}
	rl0 := &Error{Class: ClassRateLimit}
	if rl0.Retryable() {
		t.Fatal("rate_limit without retry-after must not be retryable")
	}
	if (&Error{Class: ClassProvider}).Retryable() != true {
		t.Fatal("provider must be retryable")
	}
	if (&Error{Class: ClassAuth}).Retryable() {
		t.Fatal("auth must not be retryable")
	}
	if err.Retryable() != true {
		t.Fatal("network must be retryable")
	}
	if (*Error)(nil).Retryable() {
		t.Fatal("nil error must not be retryable")
	}
}

// TestValidateErrorClass covers the validator that the task_log DAO (ticket
// 04) will reuse for the error_class CHECK constraint.
func TestValidateErrorClass(t *testing.T) {
	for _, c := range AllClasses() {
		if err := ValidateErrorClass(string(c)); err != nil {
			t.Errorf("ValidateErrorClass(%q) = %v, want nil", c, err)
		}
	}
	for _, bad := range []string{"", "unknown", "Internal", "network ", "net-work"} {
		err := ValidateErrorClass(bad)
		if err == nil {
			t.Errorf("ValidateErrorClass(%q) = nil, want error", bad)
			continue
		}
		if got := fmt.Sprint(err); len(got) < 50 {
			t.Errorf("error text for %q should enumerate the 17 classes, got %q", bad, got)
		}
	}
}

// TestMessageKeys pins the user-visible copy hooks (the copy text itself is
// D23, deferred; only the keys are frozen here).
func TestMessageKeys(t *testing.T) {
	for _, c := range AllClasses() {
		want := "error." + string(c)
		if got := c.MessageKey(); got != want {
			t.Errorf("MessageKey(%s) = %q, want %q", c, got, want)
		}
	}
}
