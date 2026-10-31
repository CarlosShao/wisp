package present

import "strings"

// The same-response spoken-summary protocol. The loop's (6) style section
// instructs the model to append a spoken summary of at most 60 characters to
// long results; this file is where that summary is lifted back out. The
// protocol shape is the dispatcher's: a LAST line that reads
// 口播摘要：<summary> (full-width or half-width colon). Only a trailing
// marker counts so an ordinary sentence that merely contains the phrase
// stays content, and the marker never reaches the panel or the spoken line
// twice.

const (
	spokenMarker = "口播摘要"
	// spokenSuffix is the tier 1/2 tail of the spoken line: D10 says the
	// announcement carries the summary AND the path.
	spokenHaveFile = "。全文已存："
	spokenNoFile   = "。全文未能保存"
)

// ExtractSummary splits one final LLM response into the display text and the
// model's own spoken summary. summary is "" when the response carried no
// trailing marker line - the caller then falls back to FirstSentence, which
// is also the tier-3 rule. The display never contains the marker line.
func ExtractSummary(text string) (display, summary string) {
	lines := strings.Split(text, "\n")
	last := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			last = i
			break
		}
	}
	if last >= 0 {
		if s, ok := markerOf(lines[last]); ok {
			lines[last] = ""
			return strings.TrimRight(strings.Join(lines, "\n"), "\n"), s
		}
	}
	return text, ""
}

// markerOf strips one 口播摘要 marker line into its payload.
func markerOf(line string) (string, bool) {
	t := strings.TrimSpace(line)
	rest, ok := strings.CutPrefix(t, spokenMarker+"：")
	if !ok {
		rest, ok = strings.CutPrefix(t, spokenMarker+":")
	}
	if !ok {
		return "", false
	}
	return strings.TrimSpace(rest), true
}

// FirstSentence returns the first sentence of s for the spoken line, cut at
// the first sentence punctuation or newline and clamped to MaxSpokenRunes.
// It is the tier-3 rule verbatim (D10: 只播报首句) and the fallback summary
// when a tier 1/2 response carried no marker.
func FirstSentence(s string, limit int) string {
	if limit <= 0 {
		limit = MaxSpokenRunes
	}
	s = strings.TrimSpace(s)
	if cut := strings.IndexAny(s, "。！？；\n"); cut >= 0 {
		s = strings.TrimSpace(s[:cut])
	}
	return clampRunes(s, limit)
}

// clampRunes cuts s to at most n runes. A hard cut can split a word; the
// alternative - speaking over the ceiling - is the one thing D10 forbids, so
// the cut wins.
func clampRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// spokenWithFile composes the tier 1/2 announcement: the summary (or the
// first-sentence fallback) plus the artifact path. The D10 ceiling ("摘要
// <=60 字") binds the summary, not the announcement: the path is appended
// whole, because a truncated path is a path the user cannot open. The
// not-saved tail stays honest when the artifacts write failed.
func spokenWithFile(summary, display, path string, saved bool) string {
	head := summary
	if head == "" {
		head = FirstSentence(display, MaxSpokenRunes)
	}
	tail := spokenNoFile
	if saved {
		tail = spokenHaveFile + path
	}
	// The model was told 60, but the ceiling is enforced here, not trusted:
	// a marker summary over the ceiling is clamped like the fallback is.
	return clampRunes(head, MaxSpokenRunes) + tail
}
