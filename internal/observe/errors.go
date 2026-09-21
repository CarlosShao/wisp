package observe

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// ErrorClass is the finite error taxonomy of D37(a) / SPEC-01 §5.1. It is an
// enum by contract: logs, StreamEvent.Error (C6) and the task_log.error_class
// column (ticket 04) accept these 17 values and nothing else. Free-form
// strings are banned (D37: they silently break statistics and alerting).
type ErrorClass string

// The 17 classes, in D37 table order.
const (
	ClassConfig           ErrorClass = "config"
	ClassAuth             ErrorClass = "auth"
	ClassNetwork          ErrorClass = "network"
	ClassRateLimit        ErrorClass = "rate_limit"
	ClassProvider         ErrorClass = "provider"
	ClassModel            ErrorClass = "model"
	ClassAudioDevice      ErrorClass = "audio_device"
	ClassASR              ErrorClass = "asr"
	ClassTool             ErrorClass = "tool"
	ClassPermissionDenied ErrorClass = "permission_denied"
	ClassUserRejected     ErrorClass = "user_rejected"
	ClassCancelled        ErrorClass = "cancelled"
	ClassBudget           ErrorClass = "budget"
	ClassLoop             ErrorClass = "loop"
	ClassInjection        ErrorClass = "injection"
	ClassInternal         ErrorClass = "internal"
	ClassResource         ErrorClass = "resource"
)

// allClasses is the authoritative set, in D37 table order.
var allClasses = []ErrorClass{
	ClassConfig, ClassAuth, ClassNetwork, ClassRateLimit, ClassProvider,
	ClassModel, ClassAudioDevice, ClassASR, ClassTool, ClassPermissionDenied,
	ClassUserRejected, ClassCancelled, ClassBudget, ClassLoop, ClassInjection,
	ClassInternal, ClassResource,
}

// AllClasses returns the 17 error classes in D37 table order. The returned
// slice is a copy and may be freely mutated by callers.
func AllClasses() []ErrorClass {
	return slices.Clone(allClasses)
}

// Valid reports whether c is one of the 17 enum values.
func (c ErrorClass) Valid() bool {
	return slices.Contains(allClasses, c)
}

// ValidateErrorClass validates a string destined for a log error_class field
// or the task_log.error_class column. The DAO-level CHECK constraint for
// task_log is reserved in internal/memory and lands with the schema (ticket
// 04); this function is the shared validator until then.
func ValidateErrorClass(s string) error {
	if ErrorClass(s).Valid() {
		return nil
	}
	return fmt.Errorf("invalid error_class %q: must be one of the 17 D37 classes (%s)",
		s, joinClasses(allClasses))
}

func joinClasses(cs []ErrorClass) string {
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = string(c)
	}
	return strings.Join(parts, "|")
}

// RetryPolicy is the D37 "retryable" column made executable. Callers own the
// actual loop; the policy only decides whether and how a retry is allowed.
type RetryPolicy uint8

const (
	// RetryNever: a retry cannot succeed without user or caller intervention
	// (config fix, new key, user consent, budget raise, ...).
	RetryNever RetryPolicy = iota
	// RetryBackoff: exponential backoff, at most MaxBackoffAttempts attempts
	// (D37 network/provider rows).
	RetryBackoff
	// RetryAfterHeader: retry only after the upstream retry-after hint
	// (Error.RetryAfter); immediate retry is forbidden.
	RetryAfterHeader
	// RetryOnce: exactly one immediate retry (D37 audio_device "re-enumerate
	// once" and asr "1 time" rows).
	RetryOnce
	// RetryAgentDecided: the outcome goes back to the LLM as a tool error and
	// the model decides whether to try again (D37 tool row).
	RetryAgentDecided
)

// MaxBackoffAttempts is the D37 ceiling for RetryBackoff classes.
const MaxBackoffAttempts = 3

// RetryPolicy maps the class to its D37 retry column.
func (c ErrorClass) RetryPolicy() RetryPolicy {
	switch c {
	case ClassNetwork, ClassProvider:
		return RetryBackoff
	case ClassRateLimit:
		return RetryAfterHeader
	case ClassAudioDevice, ClassASR:
		return RetryOnce
	case ClassTool:
		return RetryAgentDecided
	default:
		// config, auth, model, permission_denied, user_rejected, cancelled,
		// budget, loop, injection, internal, resource
		return RetryNever
	}
}

// Retryable reports whether the class allows any automatic retry at all.
func (c ErrorClass) Retryable() bool {
	return c.RetryPolicy() != RetryNever
}

// MappedStates returns the D43 BallState names this class maps to (D37
// column "mapped state"). Values are strings that must equal
// statemachine.State constants; a test pins that correspondence so the two
// packages cannot drift apart. Where a class maps to two states, the first
// entry is the pre-transition ("retry in place") state and the second the
// terminal one.
func (c ErrorClass) MappedStates() []string {
	switch c {
	case ClassConfig:
		return []string{"Unconfigured"}
	case ClassAuth:
		return []string{"Unconfigured", "Error"}
	case ClassNetwork:
		return []string{"NoNetwork"}
	case ClassRateLimit, ClassProvider, ClassInjection, ClassInternal:
		return []string{"Error"}
	case ClassModel:
		return []string{"Downloading", "Error"}
	case ClassAudioDevice:
		return []string{"Error"}
	case ClassASR:
		return []string{"Listening", "Error"}
	case ClassTool, ClassPermissionDenied:
		return []string{"Acting"}
	case ClassUserRejected, ClassCancelled:
		return []string{"Settling"}
	case ClassBudget, ClassLoop:
		return []string{"Stuck"}
	case ClassResource:
		return []string{"WatchdogAlert"}
	default:
		return nil
	}
}

// MessageKey returns the stable user-visible copy key for the class. The copy
// text itself is deferred (D23 error-copy work, S6); UI and logs embed the key
// so copy can be added later without re-touching error sites.
func (c ErrorClass) MessageKey() string {
	return "error." + string(c)
}

// Error is the canonical error type (D37; C6 StreamEvent.Error shape:
// class, provider_code, retryable, retry_after). Every cross-boundary
// failure (llm, speech cgo, tools, config, watchdog) is normalized into it.
type Error struct {
	Class        ErrorClass    `json:"class"`
	ProviderCode string        `json:"provider_code,omitempty"` // raw upstream code (HTTP status, sherpa return code, ...); "" if none
	RetryAfter   time.Duration `json:"retry_after,omitempty"`   // >0 only with RetryAfterHeader
	Detail       string        `json:"detail,omitempty"`        // log-safe detail; never carry secrets (D33 redaction still applies)
	Err          error         `json:"-"`                       // wrapped cause, optional
}

// New creates a classified error.
func New(class ErrorClass, detail string) *Error {
	return &Error{Class: class, Detail: detail}
}

// Wrap creates a classified error around a cause.
func Wrap(class ErrorClass, err error, detail string) *Error {
	return &Error{Class: class, Detail: detail, Err: err}
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s", e.Class)
	if e.ProviderCode != "" {
		fmt.Fprintf(&b, " (%s)", e.ProviderCode)
	}
	if e.Detail != "" {
		fmt.Fprintf(&b, ": %s", e.Detail)
	}
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	return b.String()
}

// Unwrap exposes the wrapped cause to errors.Is / errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Retryable reports whether THIS error may be retried (class policy, with the
// rate_limit refinement that a retry-after hint must be present and honored).
func (e *Error) Retryable() bool {
	if e == nil {
		return false
	}
	p := e.Class.RetryPolicy()
	if p == RetryAfterHeader {
		return e.RetryAfter > 0
	}
	return p != RetryNever
}

// ClassOf extracts the D37 class from an error chain via errors.As. ok is
// false for unclassified errors; unclassified errors that must be logged or
// persisted should be recorded as ClassInternal (see ClassOfOrInternal).
func ClassOf(err error) (ErrorClass, bool) {
	var e *Error
	if errors.As(err, &e) && e.Class.Valid() {
		return e.Class, true
	}
	return "", false
}

// ClassOfOrInternal is the logging/persistence convenience: unclassified
// errors collapse to internal (D37: "unexpected" is internal).
func ClassOfOrInternal(err error) ErrorClass {
	if c, ok := ClassOf(err); ok {
		return c
	}
	return ClassInternal
}
