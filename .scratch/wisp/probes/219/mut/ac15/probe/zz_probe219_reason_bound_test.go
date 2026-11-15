package main

// PROBE (receipt 219-ac1-ac5, leg A). Not a test of the product: a reading.
//
// It is injected through `go test -overlay` as a file that does not exist on the
// tree, so nothing in cmd/wisp is written. What it measures is the ONE thing the
// receipt's AC#5 cell could otherwise only infer from source:
//
//	does sanitizeReplyText truncate an over-long typed reason SILENTLY?
//
// AC#5 says 「不许静默截断，票 147 那一族「截断要带可找回指针」同纪律」. The
// function's body breaks its loop at replyReasonMax runes and returns the prefix
// with no marker, no count and no pointer; this probe prints the numbers so the
// receipt can quote a measurement instead of an opinion.
//
// It also prints the empty / whitespace-only answers, because AC#5 owes a defined
// behaviour for those two as well and there is no case anywhere that reads them
// (grep for 未填理由 finds cmd/wisp/approval_reply.go:341 and nothing else).

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestProbe219ReasonBound(t *testing.T) {
	long := strings.Repeat("字", 200)
	got := sanitizeReplyText(long)
	fmt.Printf("PROBE219 reason_bound: input_runes=%d output_runes=%d replyReasonMax=%d\n",
		utf8.RuneCountInString(long), utf8.RuneCountInString(got), replyReasonMax)
	markers := []string{"…", "...", "截断", "省略", "more", "+", "truncat"}
	found := ""
	for _, m := range markers {
		if strings.Contains(got, m) {
			found = m
			break
		}
	}
	fmt.Printf("PROBE219 reason_bound: truncation_marker_in_output=%q (empty means SILENT)\n", found)
	fmt.Printf("PROBE219 reason_bound: output_tail=%q\n", tailRunes(got, 3))

	empty := sanitizeReplyText("")
	blank := sanitizeReplyText("   ")
	fmt.Printf("PROBE219 reason_bound: empty=%q blank=%q\n", empty, blank)

	// The 160/161 boundary, so the receipt can say the cut is exactly one rune past
	// the constant and not somewhere else.
	for _, n := range []int{159, 160, 161} {
		fmt.Printf("PROBE219 reason_bound: n=%d -> %d runes\n", n, utf8.RuneCountInString(sanitizeReplyText(strings.Repeat("a", n))))
	}
}

func tailRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[len(r)-n:])
}
