package llm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Error classification at the seam (SPEC-05 sec 3.4, D42#4).
//
// Every failure that crosses the C5 boundary is normalized into a D37
// observe.Error. The provider_code field carries the discriminating raw
// code (HTTP status, upstream error code, transport detail) so statistics
// keep the 17-class enum stable while still distinguishing, e.g., an
// untrusted certificate from an unreachable host:
//
//	network (tls_untrusted_cert) - corporate TLS interception (D42#4)
//	network (connect_failed)     - DNS / refused / unreachable
//	network (timeout)            - deadline or socket timeout
//	network (stream_disconnected)- mid-stream EOF (partial content kept)
//	auth                         - 401/403: no retry, Unconfigured semantics
//	budget                       - 402 / 429 insufficient_quota: account problem
//	rate_limit                   - 429: honors retry-after
//	provider                     - 5xx / malformed payloads (retryable)
//	model                        - unknown model / context length exceeded
//	cancelled                    - caller cancelled the context

// Provider codes used by the seam and its adapters. They are log-safe and
// belong in task_log.provider_code / StreamEvent.Error.
const (
	CodeTLSUntrusted     = "tls_untrusted_cert"
	CodeConnectFailed    = "connect_failed"
	CodeTimeout          = "timeout"
	CodeStreamDisconnect = "stream_disconnected"
	CodeRetryAfter       = "retry_after"
)

// retryAfterFallback is used when a 429 arrives without a retry-after hint:
// D37's rate_limit row is RetryAfterHeader (a retry must honor a hint), so
// the seam substitutes this fallback rather than making the error
// non-retryable. Present headers always win.
const retryAfterFallback = time.Second

// transportErr assembles a classified transport error: class + discriminating
// provider code + retryability in one place.
func transportErr(err error, class observe.ErrorClass, code, detail string, retryable bool) *observe.Error {
	e := observe.Wrap(class, err, detail)
	e.ProviderCode = code
	if !retryable {
		// observe derives retryability from the class; for network-class
		// failures that must not be retried (TLS interception) we keep the
		// class but mark the wrapped error chain with a sentinel the retry
		// wrapper checks. See retryPolicyFor.
		e.Detail += " [no-retry]"
		e.Err = &noRetrySentinel{cause: e.Err}
	}
	return e
}

// noRetrySentinel marks a classified error whose class policy would allow a
// retry but whose concrete cause makes retrying pointless (persistent TLS
// interception, ...). errors.Is-aware via its own Error/Unwrap.
type noRetrySentinel struct{ cause error }

func (s *noRetrySentinel) Error() string { return "no retry possible" }
func (s *noRetrySentinel) Unwrap() error { return s.cause }

func isNoRetry(err error) bool {
	var s *noRetrySentinel
	return errors.As(err, &s)
}

// ClassifyTransportError maps a transport-layer error to its D37 class and
// discriminating provider code (D42#4: "cert untrusted" must be
// distinguishable from "cannot connect").
func ClassifyTransportError(err error) *observe.Error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return observe.Wrap(observe.ClassCancelled, err, "request cancelled")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return transportErr(err, observe.ClassNetwork, CodeTimeout, "deadline exceeded", true)
	}

	// TLS: certificate problems are machine-configuration problems - a retry
	// against the same endpoint cannot succeed, so they are NOT retryable
	// (the text_chain moves on to the next provider instead).
	var x509Err x509.UnknownAuthorityError
	if errors.As(err, &x509Err) {
		return transportErr(err, observe.ClassNetwork, CodeTLSUntrusted,
			"TLS certificate signed by unknown authority (intercepting proxy?)", false)
	}
	var certErr x509.CertificateInvalidError
	if errors.As(err, &certErr) {
		return transportErr(err, observe.ClassNetwork, CodeTLSUntrusted,
			"TLS certificate invalid: "+x509ReasonString(certErr.Reason), false)
	}
	var hostErr x509.HostnameError
	if errors.As(err, &hostErr) {
		return transportErr(err, observe.ClassNetwork, CodeTLSUntrusted,
			"TLS certificate does not match host", false)
	}
	var verifyErr *tls.CertificateVerificationError
	if errors.As(err, &verifyErr) {
		return transportErr(err, observe.ClassNetwork, CodeTLSUntrusted,
			"TLS certificate verification failed", false)
	}
	var recordErr *tls.RecordHeaderError
	if errors.As(err, &recordErr) {
		// e.g. a plaintext proxy answered on a TLS port.
		return transportErr(err, observe.ClassNetwork, CodeTLSUntrusted,
			"TLS record header error", false)
	}

	// DNS failure: authoritative host-not-found is not fixed by retrying, but
	// temporary resolver failures are; keep it retryable (backoff ladder,
	// chain moves on after exhaustion).
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return transportErr(err, observe.ClassNetwork, CodeConnectFailed, "DNS: "+dnsErr.Err, true)
	}

	// Connection refused / unreachable.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return transportErr(err, observe.ClassNetwork, CodeConnectFailed, "connection failed", true)
	}

	// Socket-level timeouts (net/http wraps them as net.Error).
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return transportErr(err, observe.ClassNetwork, CodeTimeout, "i/o timeout", true)
	}

	// EOF before a complete response body.
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return transportErr(err, observe.ClassNetwork, CodeStreamDisconnect, "stream disconnected", true)
	}

	return observe.Wrap(observe.ClassNetwork, err, "transport error")
}

// HTTPErrorDetail carries what the adapter parsed from a non-2xx response.
type HTTPErrorDetail struct {
	Status  int
	Code    string // upstream error code string (e.g. "insufficient_quota")
	Type    string // upstream error type (e.g. "invalid_request_error")
	Message string // upstream message (log-safe: redaction still applies upstream)
}

// ClassifyHTTPStatus maps an HTTP status (plus the upstream error code when
// the body carried one) to a D37 class. The bool reports whether the class
// allows a retry at all (mirrors observe.ErrorClass.Retryable, with the
// rate_limit refinement handled by RetryAfter parsing, not here).
func ClassifyHTTPStatus(d HTTPErrorDetail) (observe.ErrorClass, bool) {
	code := strings.ToLower(d.Code)
	switch d.Status {
	case 400, 422:
		// Bad request: our bug or an unsupported parameter. Model-flavored
		// rejections get their own class for actionable copy.
		if strings.Contains(code, "model") || strings.Contains(strings.ToLower(d.Message), "model") {
			return observe.ClassModel, false
		}
		if strings.Contains(code, "context_length") || strings.Contains(strings.ToLower(d.Message), "context length") {
			return observe.ClassModel, false
		}
		return observe.ClassProvider, false
	case 401, 403:
		// SPEC-05 sec 3.4: never retried; Unconfigured semantics (the class
		// maps to the Unconfigured ball state via observe.MappedStates).
		return observe.ClassAuth, false
	case 402:
		// Account / credit problem: distinguishable from network trouble.
		return observe.ClassBudget, false
	case 404:
		if strings.Contains(code, "model") || strings.Contains(strings.ToLower(d.Message), "model") {
			return observe.ClassModel, false
		}
		return observe.ClassProvider, false
	case 408:
		return observe.ClassNetwork, true
	case 409:
		return observe.ClassProvider, false
	case 413:
		if strings.Contains(code, "context") {
			return observe.ClassModel, false
		}
		return observe.ClassProvider, false
	case 429:
		// OpenAI-style quota exhaustion arrives as a 429 with
		// code=insufficient_quota: that is an ACCOUNT problem (no retry),
		// not a pace problem.
		if strings.Contains(code, "insufficient_quota") || strings.Contains(code, "quota") ||
			strings.Contains(code, "balance") || strings.Contains(code, "arrears") {
			return observe.ClassBudget, false
		}
		return observe.ClassRateLimit, true
	case 500, 501, 502, 503, 504:
		return observe.ClassProvider, true
	case 529:
		// Overloaded (Anthropic-style; some compat providers use it too).
		return observe.ClassProvider, true
	default:
		if d.Status >= 500 {
			return observe.ClassProvider, true
		}
		return observe.ClassProvider, false
	}
}

// NewHTTPError builds the classified observe.Error for a non-2xx response.
// retryAfterHint is the parsed Retry-After header (0 = absent): for 429 it is
// attached so observe.Error.Retryable() holds; a missing hint falls back to
// retryAfterFallback so pace-limiting still backs off (see retryAfterFallback).
func NewHTTPError(d HTTPErrorDetail, retryAfterHint time.Duration) *observe.Error {
	class, _ := ClassifyHTTPStatus(d)
	e := observe.New(class, httpErrorDetailString(d))
	e.ProviderCode = d.Code
	if e.ProviderCode == "" {
		e.ProviderCode = strconv.Itoa(d.Status)
	}
	if class == observe.ClassRateLimit {
		if retryAfterHint > 0 {
			e.RetryAfter = retryAfterHint
		} else {
			e.RetryAfter = retryAfterFallback
			e.ProviderCode = CodeRetryAfter
		}
	}
	return e
}

func httpErrorDetailString(d HTTPErrorDetail) string {
	var b strings.Builder
	fmt.Fprintf(&b, "HTTP %d", d.Status)
	if d.Type != "" {
		fmt.Fprintf(&b, " type=%s", d.Type)
	}
	if d.Code != "" {
		fmt.Fprintf(&b, " code=%s", d.Code)
	}
	if d.Message != "" {
		fmt.Fprintf(&b, ": %s", truncateRunes(d.Message, 300))
	}
	return b.String()
}

// x509ReasonString names an x509.InvalidReason without depending on the
// stdlib's fmt printing of the enum.
func x509ReasonString(r x509.InvalidReason) string {
	switch r {
	case x509.NotAuthorizedToSign:
		return "not_authorized_to_sign"
	case x509.Expired:
		return "expired"
	case x509.CANotAuthorizedForThisName:
		return "ca_not_authorized_for_this_name"
	case x509.TooManyIntermediates:
		return "too_many_intermediates"
	case x509.IncompatibleUsage:
		return "incompatible_usage"
	case x509.NameMismatch:
		return "name_mismatch"
	case x509.NameConstraintsWithoutSANs:
		return "name_constraints_without_sans"
	case x509.UnconstrainedName:
		return "unconstrained_name"
	case x509.TooManyConstraints:
		return "too_many_constraints"
	case x509.CANotAuthorizedForExtKeyUsage:
		return "ca_not_authorized_for_ext_key_usage"
	default:
		return "invalid"
	}
}

// ParseRetryAfter parses a Retry-After header value: delta-seconds (integer)
// or an HTTP-date. Absent/unparseable yields 0. Dates are evaluated against
// the monotonic-safe now parameter.
func ParseRetryAfter(h string, now time.Time) time.Duration {
	h = strings.TrimSpace(h)
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		d := t.Sub(now)
		if d > 0 {
			return d
		}
	}
	return 0
}

// ProxyFunc returns the proxy function the seam's default HTTP transport
// uses: HTTP(S)_PROXY environment (via http.ProxyFromEnvironment, which also
// honors NO_PROXY) layered with the Windows system (WinINET/WinHTTP default)
// proxy settings. mode is the config net.proxy.mode enum: "system" consults
// the OS settings, "manual" forces cfgURL, "none" disables proxying.
//
// On non-Windows platforms the system component is a no-op (env rules apply).
func ProxyFunc(mode, cfgURL string) func(*http.Request) (*url.URL, error) {
	switch mode {
	case "manual":
		u, err := url.Parse(cfgURL)
		return func(*http.Request) (*url.URL, error) {
			if err != nil {
				return nil, err
			}
			return u, nil
		}
	case "none":
		return func(*http.Request) (*url.URL, error) { return nil, nil }
	default: // "system" and anything else: env + OS settings
		envProxy := http.ProxyFromEnvironment
		sysProxy := systemProxyFunc()
		if sysProxy == nil {
			return envProxy
		}
		return func(req *http.Request) (*url.URL, error) {
			if u, err := sysProxy(req); u != nil || err != nil {
				return u, err
			}
			return envProxy(req)
		}
	}
}

// NewDefaultHTTPClient builds the transport the adapters use unless overridden:
// proxy rules per ProxyFunc, sane timeouts handled by contexts (no
// wall-clock math, D42#9) - only the TLS handshake and idle pool get bounds.
func NewDefaultHTTPClient(proxyMode, proxyURL string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 ProxyFunc(proxyMode, proxyURL),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          8,
			MaxIdleConnsPerHost:   4,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
	}
}
