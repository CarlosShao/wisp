package present

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// D10's frozen tier ladder (PLAN.md D10 / SPEC-05 §7, fourth-round final
// values). These are contract constants, not tuning knobs: there is no
// configuration key for them and none may be added (the dispatch thresholds
// are part of the frozen contract text).
const (
	// shortLimit is tier 3's floor: at most this many hanzi (and no
	// structured marker) is a TTS-only short result.
	shortLimit = 60
	// longHanziLimit is tier 2's hanzi condition: more than this many hanzi
	// is a file-drop result.
	longHanziLimit = 400
	// longCharLimit is tier 2's character backstop: more than this many
	// non-whitespace characters is a file-drop result even when the hanzi
	// count stays low (a long English reply must not be read aloud whole).
	longCharLimit = 800
	// MaxSpokenRunes caps every spoken line at the D10 "60 字" ceiling,
	// counted in runes.
	MaxSpokenRunes = 60
)

// Tier names the D10 rule that matched. TierNone means there was nothing to
// present (an empty result routes nowhere, on purpose).
type Tier string

const (
	// TierNone: empty or whitespace-only result. No channel fires.
	TierNone Tier = "none"
	// TierStructured is D10 rule 1: code block / table / list of >=3 items /
	// >=2 links / file-path list. File-drop + panel + spoken summary + path.
	TierStructured Tier = "structured"
	// TierLong is D10 rule 2: more than 400 hanzi or more than 800
	// characters. Same channels as rule 1.
	TierLong Tier = "long"
	// TierMedium is D10 rule 3: 61-400 hanzi. Notify + panel, first
	// sentence spoken, badge and task list updated unconditionally (D42#8).
	TierMedium Tier = "medium"
	// TierShort is D10 rule 4: at most 60 hanzi and unstructured. TTS +
	// ball + clipboard write; no panel.
	TierShort Tier = "short"
)

// HanziCount counts the Han-script runes of s - the unit D10's thresholds
// are written in.
func HanziCount(s string) int {
	n := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			n++
		}
	}
	return n
}

// CharCount counts the non-whitespace runes of s - the "字符" backstop of
// D10 rule 2. Whitespace is excluded so padding cannot push a result into
// file-drop.
func CharCount(s string) int {
	n := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			n++
		}
	}
	return n
}

// Classify applies the D10 ladder to one final text: first match wins and
// structured markers beat length, so a 40-character table is rule 1 no
// matter how short it is. The returned rule string names the matched
// condition in receipt-visible words. An empty (or whitespace-only) text is
// TierNone: there is nothing to present on any channel.
func Classify(text string) (Tier, string) {
	if strings.TrimSpace(text) == "" {
		return TierNone, "empty result: nothing to present"
	}
	if reason, ok := structuredReason(text); ok {
		return TierStructured, reason
	}
	h, c := HanziCount(text), CharCount(text)
	switch {
	case h > longHanziLimit:
		return TierLong, fmt.Sprintf("length: %d hanzi over the %d limit", h, longHanziLimit)
	case c > longCharLimit:
		return TierLong, fmt.Sprintf("length: %d chars over the %d limit", c, longCharLimit)
	case h > shortLimit:
		return TierMedium, fmt.Sprintf("length: %d hanzi in the 61-%d band", h, longHanziLimit)
	default:
		return TierShort, fmt.Sprintf("length: %d hanzi at or under %d, unstructured", h, shortLimit)
	}
}

// ---------------------------------------------------------------------------
// The rule-1 structured-marker detector. Each rule is deliberately plain:
// the detector must stay explainable in a receipt, and the "40-character
// table" case of SPEC-05 §9 must hit it.

var (
	mdLinkRe   = regexp.MustCompile(`\[[^\]]*\]\([^)]+\)`)
	bulletRe   = regexp.MustCompile(`^\s*(?:[-*+]|\d{1,3}[.)])\s+\S`)
	driveRe    = regexp.MustCompile(`^[A-Za-z]:[\\/]`)
	uncRe      = regexp.MustCompile(`^\\\\`)
	extRe      = regexp.MustCompile(`\.[A-Za-z0-9]{1,8}$`)
	schemeRe   = regexp.MustCompile(`https?://`)
	leadBullet = regexp.MustCompile(`^\s*(?:[-*+]|\d{1,3}[.)])\s+`)
)

// structuredReason reports which rule-1 marker the text carries, if any.
// A code fence or a markdown table anywhere, a list of at least three items,
// at least two links, or a list of at least two file-path lines makes the
// result rule 1. A single path inside prose stays unstructured on purpose:
// D10 says a short spoken answer must not be forced into a file because it
// mentions where a file went.
func structuredReason(text string) (string, bool) {
	lines := strings.Split(text, "\n")
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			return "structured: code fence", true
		}
	}
	for i := 0; i+1 < len(lines); i++ {
		if strings.Count(lines[i], "|") >= 2 && isTableDelimiter(lines[i+1]) {
			return "structured: markdown table", true
		}
	}
	items := 0
	for _, ln := range lines {
		if bulletRe.MatchString(ln) {
			items++
		}
	}
	if items >= 3 {
		return fmt.Sprintf("structured: list of %d items", items), true
	}
	if n := linkCount(text); n >= 2 {
		return fmt.Sprintf("structured: %d links", n), true
	}
	paths := 0
	for _, ln := range lines {
		if isPathLine(ln) {
			paths++
		}
	}
	if paths >= 2 {
		return fmt.Sprintf("structured: %d path lines", paths), true
	}
	return "", false
}

// isTableDelimiter matches the |---|:-:| separator row under a markdown
// table header: only pipes, dashes, colons and blanks, with at least one
// dash run so a plain cell row never passes.
func isTableDelimiter(line string) bool {
	t := strings.TrimSpace(line)
	if t == "" || !strings.Contains(t, "-") {
		return false
	}
	for _, r := range t {
		if r != '|' && r != '-' && r != ':' && r != ' ' && r != '\t' {
			return false
		}
	}
	return true
}

// linkCount counts markdown links plus bare http(s) URLs, with the markdown
// spans removed first so a link's own href is never counted twice.
func linkCount(text string) int {
	n := len(mdLinkRe.FindAllString(text, -1))
	rest := mdLinkRe.ReplaceAllString(text, "")
	return n + len(schemeRe.FindAllString(rest, -1))
}

// isPathLine reports whether one line reads as a file path: a drive letter,
// a UNC spelling, a POSIX absolute path, or any line that carries a path
// separator and ends in a short extension.
func isPathLine(line string) bool {
	t := leadBullet.ReplaceAllString(line, "")
	t = strings.TrimSpace(t)
	if t == "" || mdLinkRe.MatchString(t) {
		return false
	}
	if driveRe.MatchString(t) || uncRe.MatchString(t) || strings.HasPrefix(t, "/") {
		return true
	}
	return (strings.ContainsRune(t, '/') || strings.ContainsRune(t, '\\')) &&
		extRe.MatchString(t)
}
