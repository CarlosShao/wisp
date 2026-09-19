package observe

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"

	"github.com/CarlosShao/wisp/internal/secret"
)

// Log redaction (ticket 08; D33, SPEC-10 §5.1) - HARD-CODED, NON-DISABLEABLE.
//
// The Redactor sits between slog and the JSONL sink and applies five frozen
// rule classes. There is deliberately NO config switch that turns redaction
// off (D22: an off switch would be the "loosen the gate instead of meet it"
// failure mode); the only knob is [privacy] redact_paths, which the frozen
// config model already defines as an opt-in for PATH masking, and which can
// only ever redact MORE, never less.
//
//   1. Secret-shaped attribute values (attribute names whose word parts
//      include key/token/secret/password/...): last 4 characters at most
//      (C28, via secret.RedactSecret).
//   2. Audio buffers never reach the log: audio/pcm/sample-shaped attribute
//      names are replaced by a byte-count placeholder.
//   3. web.fetch / file bodies never reach the log: body/content-shaped
//      attribute names are replaced by a byte-count placeholder.
//   4. Long argument strings are truncated (bounded log lines).
//   5. Inline secret material (sk-... keys, Bearer tokens, api_key=... in
//      free text) is masked wherever it appears in a string value or message.
//   6. Path masking (rule 6, [privacy] redact_paths) masks filesystem paths
//      in string values when enabled.

const (
	// MaxLoggedString bounds any single string that reaches the log (rule 4).
	MaxLoggedString = 512
	// maxScanBytes bounds the inline-secret scan; longer values are truncated
	// first (scanning megabyte bodies would itself be a DoS).
	maxScanBytes = 64 << 10
)

// Redactor applies the frozen redaction rules. The zero value is valid and
// already redacts (only path masking is off by default).
type Redactor struct {
	// RedactPaths mirrors [privacy] redact_paths: mask filesystem paths in
	// string values. Opt-in, and it can only redact more, never less.
	RedactPaths bool
}

// secretWords are the attribute-name word parts that mark a value as secret
// material. Word parts come from splitting the key on non-alphanumerics, so
// "openai.api_key" and "provider_api_key" match via their "key" part while
// "hotkey" (one word) does not.
var secretWords = map[string]bool{
	"key": true, "token": true, "secret": true, "password": true,
	"passwd": true, "pwd": true, "auth": true, "authorization": true,
	"credential": true, "credentials": true, "bearer": true, "passphrase": true,
}

// audioWords mark audio buffer payloads (rule 2: never logged).
var audioWords = map[string]bool{
	"audio": true, "pcm": true, "samples": true, "waveform": true, "wav": true,
}

// contentWords mark fetched/file payload bodies (rule 3: never logged).
var contentWords = map[string]bool{
	"body": true, "html": true, "content": true, "document": true, "payload": true,
}

// Inline secret shapes (rule 5): OpenAI-style keys, Bearer headers, and
// api_key=... style assignments in free text. Only key-shaped runs are
// masked; ordinary words pass untouched. providerPrefixRe adds the common
// token prefixes (GitHub, Anthropic, Google, Slack) to both maskers and the
// bundle's post-mask detector.
var (
	inlineKeyRe = regexp.MustCompile(`sk-[A-Za-z0-9_-]{8,}`)
	bearerRe    = regexp.MustCompile(`(?i)\bbearer[ \t_-]+[A-Za-z0-9._~-]{16,}`)
	assignRe    = regexp.MustCompile(`(?i)\b(api[_-]?key|access[_-]?token|refresh[_-]?token|secret[_-]?key)(\s*[:=]\s*)([A-Za-z0-9._~-]{12,})`)
	providerRe  = regexp.MustCompile(`(?:\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{8,}|\bgithub_pat_[A-Za-z0-9_]{8,}|\bsk-ant-[A-Za-z0-9_-]{8,}|\bAIza[A-Za-z0-9_-]{8,}|\bxox[baprs]-[A-Za-z0-9-]{8,})`)
	pathRe      = regexp.MustCompile(`(?i)(?:\\{1,2}\?\\)?(?:[a-z]:\\|\\\\)[^\s"',;:)>\]]+|(?:(?:/home|/users|/root)/[^\s"',;:)>\]]+)`)

	// snapshotDetectRe is deliberately >= the mask set (mask + provider
	// prefixes): if the detector still fires AFTER masking, the bundle
	// builder aborts fail-closed instead of shipping a leaking snapshot.
	snapshotDetectRe = regexp.MustCompile(`sk-[A-Za-z0-9_-]{8,}|(?i)\bbearer[ \t_-]+[A-Za-z0-9._~-]{16,}|(?i)\b(api[_-]?key|access[_-]?token|refresh[_-]?token|secret[_-]?key)(\s*[:=]\s*)[A-Za-z0-9._~-]{12,}|` + providerRe.String())
)

func keyWords(name string) []string {
	return strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func hasWord(name string, set map[string]bool) bool {
	for _, w := range keyWords(name) {
		if set[w] {
			return true
		}
	}
	return false
}

// Attr redacts one log attribute. It is the single choke point every record
// passes through (redactHandler), including attributes baked in via
// slog.Handler.WithAttrs.
func (r Redactor) Attr(key string, v slog.Value) slog.Attr {
	switch {
	case hasWord(key, secretWords):
		// C28: named secret material is reduced to the last 4 characters by
		// the one sanctioned helper, never by a local copy.
		return slog.String(key, secret.RedactSecret(v.String()))
	case hasWord(key, audioWords):
		return slog.String(key, fmt.Sprintf("[audio buffer redacted: %d bytes]", valueLen(v)))
	case hasWord(key, contentWords):
		return slog.String(key, fmt.Sprintf("[content redacted: %d bytes]", valueLen(v)))
	}
	switch v.Kind() {
	case slog.KindString:
		return slog.String(key, r.String(v.String()))
	case slog.KindAny:
		if b, ok := v.Any().([]byte); ok {
			return slog.String(key, r.String(string(b)))
		}
		return slog.Attr{Key: key, Value: v}
	default:
		return slog.Attr{Key: key, Value: v}
	}
}

// String applies the frozen rules to a free-form string (attribute values,
// messages, config snapshot lines).
func (r Redactor) String(s string) string {
	if len(s) > maxScanBytes {
		s = s[:maxScanBytes]
	}
	s = maskInlineSecrets(s)
	if r.RedactPaths {
		s = pathRe.ReplaceAllString(s, "<path>")
	}
	return truncate(s, MaxLoggedString)
}

// Message redacts a log message body (same rules as String).
func (r Redactor) Message(msg string) string { return r.String(msg) }

// maskInlineSecrets masks sk-/Bearer/api_key=... shaped material, keeping the
// last 4 characters (C28).
func maskInlineSecrets(s string) string {
	s = inlineKeyRe.ReplaceAllStringFunc(s, func(m string) string {
		return "sk-" + last4(m[len("sk-"):])
	})
	s = bearerRe.ReplaceAllStringFunc(s, func(m string) string {
		i := strings.LastIndexAny(m, " \t_-")
		return m[:i+1] + "****" + last4(strings.TrimSpace(m[i+1:]))
	})
	s = maskProviderTokens(s)
	return assignRe.ReplaceAllStringFunc(s, func(m string) string {
		sub := assignRe.FindStringSubmatch(m)
		if len(sub) < 4 {
			return m
		}
		return sub[1] + sub[2] + "****" + last4(sub[3])
	})
}

// maskProviderTokens masks the common provider token prefixes.
func maskProviderTokens(s string) string {
	return providerRe.ReplaceAllStringFunc(s, func(m string) string {
		i := strings.LastIndexByte(m, '_')
		if i < 0 || i == len(m)-1 {
			i = len(m) - 5
			if i < 0 {
				return "****"
			}
			return "****" + last4(m[i:])
		}
		return m[:i+1] + "****" + last4(m[i+1:])
	})
}

func last4(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return "****"
	}
	if len(r) <= 4 {
		return string(r)
	}
	return string(r[len(r)-4:])
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + fmt.Sprintf("...(truncated, %d chars total)", len(runes))
}

func valueLen(v slog.Value) int {
	switch v.Kind() {
	case slog.KindString:
		return len(v.String())
	case slog.KindAny:
		switch b := v.Any().(type) {
		case []byte:
			return len(b)
		case string:
			return len(b)
		}
	}
	return 0
}

// ConfigSnapshot runs a whole config/log file text through the redactor
// (diagnostics bundle path). It returns the redacted text and reports whether
// any key-shaped material survived, so the bundle builder can refuse to ship
// a snapshot it could not fully clean (fail-closed; the detector is
// deliberately >= the mask set).
func (r Redactor) ConfigSnapshot(text string) (string, bool) {
	out := r.String(text)
	return out, !snapshotDetectRe.MatchString(out)
}
