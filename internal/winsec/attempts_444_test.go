package winsec

// Ticket 444: the one knob the seven instruments in this file's neighbourhood share.
//
// WHY A KNOB AND NOT A SMALLER CONSTANT. AC#1 asks each instrument to be run with the
// window NOT reached, so that the sentence it prints when the window is absent is a
// READING rather than an argument. The honest way to do that is to shrink the
// denominator - the same move ticket 441 used (`WISP_120_ATTEMPTS`) - and not to edit
// the committed denominator, because the committed number is what the census and the CI
// readings are comparable with.
//
// The default is the caller's own constant, verbatim, so a run without the variable is
// the instrument that was there before this ticket, byte-for-byte in behaviour. Only the
// probe scripts under `.scratch/wisp/probes/444/` set it.

import (
	"os"
	"strconv"
)

// attemptsFor444 scales one instrument's denominator for a probe run.
func attemptsFor444(def int) int {
	if v := os.Getenv("WISP_444_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return def
}
