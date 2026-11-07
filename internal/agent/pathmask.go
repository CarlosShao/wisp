package agent

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Ticket 264, form 甲 (orchestrator decision, A-账): tool failure sentences
// carry the machine's absolute paths - C26 Canonicalize's
// "%w: %s expands to %s (%s)" plus the *fs.PathError / *os.LinkError spellings
// - into Result.Text, and that text reaches the cloud provider inside history.
// The assembly sites (103 / 112 by the AC#0 recount) are spread across 24
// files; editing them one by one means every future tool re-opens the leak on
// day one. The mask therefore runs at the ONE chokepoint every model-visible
// tool answer crosses: toolResultParts (loop.go). The user-facing side is fed
// from log.Text BEFORE this function runs (the EvToolEnd event and
// res.ToolLog), so the user still sees exactly which file was meant.
//
// Mask rule (the one sentence to quote): a span that starts at a drive-letter
// absolute path ("X:\..." or "X:/...", not the tail of a longer token such as
// https://) and runs until a filename-illegal character (whitespace, ":",
// "：", "，", "。", "；", quote, <, >, |, *, ?) is replaced wholesale by
// "<path>"; within that span, a space-separated continuation is absorbed only
// when the token itself contains a separator, so a path with a space in a
// middle segment is fully masked while the error prose that follows
// ("…: 系统找不到指定的文件。") survives.
//
// Known boundaries (documented, accepted): a space inside the FINAL segment
// ("C:\a\new file.txt") leaves that tail visible - paths are not re-parsed out
// of prose, and the bias is toward not eating error text; POSIX-style "/usr/…"
// paths are not masked (D7: Windows-only build; masking them would eat URLs
// and ratios). The rule can only replace a path with a placeholder; it never
// invents one.
//
// The single exemption is ticket 174's pointer: the path immediately after
// "全文见 " is written by the host itself (internal/agent/spill.go and
// internal/tools/task.go print the same shape), and the model's only road back
// to a spilled output is fs.read on that exact string - masking it would be
// 只截不指 (PLAN.md:2564) and would break the 174r8 pin
// (cmd/wisp/task_output_pointer_notice_174r8_test.go).

// maskPathPlaceholder replaces every masked absolute path.
const maskPathPlaceholder = "<path>"

// pathPointerPreamble is the exact preamble both spill legs print immediately
// before the artifact path ("，全文见 %s…").
const pathPointerPreamble = "全文见 "

// pathStartRe finds where an absolute Windows path may begin. The byte before
// the drive letter decides: a letter/digit/'-'/'_' there means this is the
// tail of a longer token (https://, an id), not a path.
var pathStartRe = regexp.MustCompile(`[A-Za-z]:[\\/]`)

// maskOutboundPaths rewrites text so no machine directory structure reaches
// the model-visible copy of a tool answer. It is called only from
// toolResultParts - the one chokepoint - never on the user-facing copy.
func maskOutboundPaths(text string) string {
	if !strings.Contains(text, ":\\") && !strings.Contains(text, ":/") {
		return text
	}
	var b strings.Builder
	b.Grow(len(text))
	last := 0
	masked := false
	for _, loc := range pathStartRe.FindAllStringIndex(text, -1) {
		start := loc[0]
		if start < last {
			continue // inside a span this pass already masked
		}
		if start > 0 && isPathGlueByte(text[start-1]) {
			continue // tail of a longer token (https://), not a path
		}
		if strings.HasSuffix(text[:start], pathPointerPreamble) {
			continue // ticket 174 pointer: the model's road back to the artifact
		}
		end := pathSpanEnd(text, start)
		b.WriteString(text[last:start])
		b.WriteString(maskPathPlaceholder)
		last = end
		masked = true
	}
	if !masked {
		return text
	}
	b.WriteString(text[last:])
	return b.String()
}

// isPathGlueByte reports whether c, the byte immediately before a "X:/" or
// "X:\" candidate start, makes that candidate the tail of a longer token
// ("https://" and friends) rather than a drive path. Multi-byte UTF-8
// continuation bytes fall outside every range here and correctly read as
// "not glue" (a Chinese prefix before a path is prose, not a token).
func isPathGlueByte(c byte) bool {
	return c == '-' || c == '_' ||
		(c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// pathSpanStops reports whether r ends a path span. Every rune here is one a
// Windows filename cannot contain (plus the fullwidth prose punctuation the
// tool sentences put right after a path), so stopping never cuts a real path.
func pathSpanStops(r rune) bool {
	switch r {
	case ':', '：', '，', '。', '；', '、', '"', '\'', '<', '>', '|', '*', '?':
		return true
	}
	return unicode.IsSpace(r)
}

// isDrivePrefix reports whether text starts with an absolute Windows drive
// prefix: a letter, a colon, and a separator.
func isDrivePrefix(text string) bool {
	r, size := utf8.DecodeRuneInString(text)
	if !unicode.IsLetter(r) || size != 1 {
		return false
	}
	if text[1] != ':' {
		return false
	}
	return text[2] == '\\' || text[2] == '/'
}

// pathSpanEnd scans from the drive letter to the end of the path span. The
// first three bytes are the "X:/" prefix itself (its colon is not a stop);
// the initial run then stops at any span-stopping rune; a following
// continuation (spaces/tabs + token) is absorbed only while the token itself
// contains a separator, so "C:\Users\a b\file.txt" is masked whole while the
// error prose after the path ("…: 系统找不到指定的文件。") is never eaten.
func pathSpanEnd(text string, start int) int {
	i, n := start+3, len(text) // skip the "X:/" prefix: its colon is part of it
	for i < n {
		r, size := utf8.DecodeRuneInString(text[i:])
		if pathSpanStops(r) {
			break
		}
		i += size
	}
	for i < n {
		j := i
		for j < n {
			r, size := utf8.DecodeRuneInString(text[j:])
			if r != ' ' && r != '\t' {
				break
			}
			j += size
		}
		if j == i {
			break // no whitespace: no continuation
		}
		k, hasSep := j, false
		// A continuation may itself start with a drive prefix ("rename
		// D:\a D:\b"): its colon is part of the token, not a stop.
		if k+2 < n && isDrivePrefix(text[k:]) {
			k += 3
			hasSep = true
		}
		for k < n {
			r, size := utf8.DecodeRuneInString(text[k:])
			if pathSpanStops(r) {
				break
			}
			if r == '\\' || r == '/' {
				hasSep = true
			}
			k += size
		}
		if !hasSep || k == j {
			break // prose after the path: leave it alone
		}
		i = k
	}
	return i
}
