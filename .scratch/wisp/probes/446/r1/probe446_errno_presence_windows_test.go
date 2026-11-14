//go:build windows

// Probe 446 AC#3 -- the tooth that must be able to go RED.
//
// It is NOT a repo test: it lives under .scratch/wisp/probes/446/ and is pulled
// into package `tools` through `go test -overlay` (see overlay-446.json), so the
// tree keeps no test file it did not already have. It exercises the single
// decision point both call sites now share (errnoPresent), and it walks the
// errno=0 path in three shapes: the bare value Proc.Call hands back after a
// success, that value behind one wrapper, and the same behind two.
//
// Teeth: put the old message-text comparison back inside errnoPresent and the
// wrapped rows go RED, because that text is the prefix ("ctx: ") plus the
// sentence, and a sentence compare only ever sees the outermost text.
package tools

import (
	"errors"
	"fmt"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

type probe446Case struct {
	name string
	err  error
	// want = "this error carries a real errno" (i.e. something to report).
	want bool
}

// localizedZeroErrno is the same errno 0 the ticket is about, on the Windows
// that broke the old code: the code is 0, only the TEXT is not the English
// sentence. A real box reaches this state without any fake type (Go asks
// FormatMessage for en-US, falls back to the system default when it cannot,
// syscall_windows.go:169-171); the fake exists only because this bench has an
// English UI language pack. Is() is what errors.Is consults after unwrapping,
// so the judgement must answer "no errno" here too.
type localizedZeroErrno struct{}

func (localizedZeroErrno) Error() string        { return "操作成功完成。" }
func (localizedZeroErrno) Is(target error) bool { return target == syscall.Errno(0) }

func TestProbe446ErrnoPresenceIsDecidedByCode(t *testing.T) {
	wrappedZero := fmt.Errorf("SendInput wrapper: %w", syscall.Errno(0))
	cases := []probe446Case{
		// The errno=0 path, the one the old comparison claimed to recognise.
		{name: "nil", err: nil, want: false},
		{name: "errno-zero-raw", err: syscall.Errno(0), want: false},
		{name: "errno-zero-via-windows-alias", err: windows.Errno(0), want: false},
		{name: "errno-zero-localized-text", err: localizedZeroErrno{}, want: false},
		{name: "errno-zero-wrapped", err: wrappedZero, want: false},
		{name: "errno-zero-wrapped-twice", err: fmt.Errorf("outer: %w", wrappedZero), want: false},
		// Real failures must still be reported.
		{name: "errno-5-access-denied", err: syscall.Errno(5), want: true},
		{name: "errno-87-bad-parameter", err: syscall.Errno(87), want: true},
		{name: "errno-87-wrapped", err: fmt.Errorf("ctx: %w", syscall.Errno(87)), want: true},
		// The sentence is not an identity: an error that merely READS like the
		// success text is still an error, and the old comparison swallowed it
		// as "no errno". This row is why the literal below is here at all.
		{name: "lookalike-text-is-not-errno-zero",
			err:  errors.New("The operation completed successfully."),
			want: true},
	}

	zeroCase := 0
	for _, c := range cases {
		got := errnoPresent(c.err)
		if c.err != nil && !c.want {
			zeroCase++
		}
		t.Logf("probe446 case %-32s errnoPresent=%-5v want=%-5v err=%q",
			c.name, got, c.want, probe446Text(c.err))
		if got != c.want {
			t.Fatalf("PROBE 446 RED at case %q: errnoPresent(%q) = %v, want %v -- the decision is reading the message text again instead of the code",
				c.name, probe446Text(c.err), got, c.want)
		}
	}
	t.Logf("probe446 GREEN: %d cases, %d of them on the errno=0 path", len(cases), zeroCase)
}

func probe446Text(e error) string {
	if e == nil {
		return "<nil>"
	}
	return e.Error()
}
