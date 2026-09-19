package secret

// RedactSecret renders the log-safe form of a secret: everything but the
// last 4 runes collapses to '*'. Secrets of 4 runes or fewer render as a
// fixed "****" - short secrets reveal nothing (C28, SPEC-03 §5.3: logs and
// diagnostics carry at most the last 4 characters of a key).
func RedactSecret(secret string) string {
	r := []rune(secret)
	if len(r) <= 4 {
		return "****"
	}
	return "****" + string(r[len(r)-4:])
}
