package risk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// C26 PathResolver (SPEC-06 §4): the ONLY path normalization entry point in
// the codebase. Pipeline order is fixed and must not be reordered:
//
//	expand (env / ~) -> absolutize -> Clean -> handle-based real path
//	(GetFinalPathNameByHandle VOLUME_NAME_DOS) -> FILE_ATTRIBUTE_REPARSE_POINT
//	detection with default DENY -> 8.3 short-name expansion -> UNC normalization
//
// "Parse-then-continue" on reparse points is forbidden (SPEC-06 §4): a
// junction whose source sits outside and target inside (or vice versa) is
// only caught by refusing reparse traversal itself, so the resolver denies
// any path that traverses a reparse point unless that exact component is
// listed in the user's [fs] reparse_point_exceptions.
//
// D22: filepath.Clean / filepath.Abs are sanctioned ONLY inside this file.
// Any other use for fs decisions is a CI-failing violation (see
// scripts/check-pathclean-ban.sh).

// ErrReparseDenied is returned when a path traverses a junction/symlink that
// is not covered by reparse_point_exceptions. Fail-closed by design.
var ErrReparseDenied = errors.New("risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions")

// Result is the outcome of resolving one path through the C26 pipeline.
type Result struct {
	// Canonical is the real path after handle resolution (8.3 expanded,
	// reparse targets resolved, UNC and \\?\ prefixes normalized). When the
	// path does not exist yet, it is the lexically cleaned spelling —
	// classification still applies so nonexistent sensitive targets cannot
	// be probed into existence.
	Canonical string
	// Reparse reports whether the INPUT traversed any reparse point that
	// was covered by an explicit exception (non-excepted traversal returns
	// ErrReparseDenied instead).
	Reparse bool
	// Resolved reports whether handle-based resolution succeeded (the path
	// exists and the OS gave us its final form).
	Resolved bool
}

// Resolve runs the fixed C26 pipeline on input. exceptions carries the
// user-configured [fs] reparse_point_exceptions (explicit paths only, never
// prefixes). Nonexistent paths resolve lexically and still classify; only a
// non-exempted reparse traversal is an error.
func Resolve(input string, exceptions []string) (Result, error) {
	p := expandInput(input)
	p = lexCanonical(p)      // the one sanctioned lexical step
	p = normalizeLocalUNC(p) // \\localhost\c$\... and \\?\UNC\localhost\c$\... -> c:\...
	p = lexCanonical(p)      // re-clean after prefix rewriting
	exems := make(map[string]bool, len(exceptions))
	for _, e := range exceptions {
		ee := normalizeLocalUNC(lexCanonical(expandInput(e)))
		exems[strings.ToLower(ee)] = true
	}

	res := Result{}
	if comps := reparseComponents(p); len(comps) > 0 {
		for _, c := range comps {
			if exems[strings.ToLower(c)] {
				res.Reparse = true
				continue
			}
			return res, ErrReparseDenied
		}
	}

	if final, ok := resolveHandle(p); ok {
		final = stripExtendedPrefix(final)
		final = normalizeLocalUNC(final)
		if final != "" {
			res.Canonical = final
			res.Resolved = true
			return res, nil
		}
	}
	// Path does not exist (or handle query unsupported): lexical fallback,
	// classification still enforced — fail-closed.
	res.Canonical = p
	return res, nil
}

// expandInput expands environment variables (%VAR% and $VAR) and a leading
// ~ (the user home). Unknown constructs pass through untouched.
func expandInput(input string) string {
	p := strings.TrimSpace(input)
	if p == "" {
		return p
	}
	// Windows-style %VAR%.
	for start := strings.Index(p, "%"); start >= 0; start = strings.Index(p[start+1:], "%") + start + 1 {
		end := strings.Index(p[start+1:], "%")
		if end < 0 {
			break
		}
		name := p[start+1 : start+1+end]
		if v := os.Getenv(name); v != "" && !strings.Contains(name, "%") {
			p = p[:start] + v + p[start+1+end+1:]
		}
	}
	p = os.ExpandEnv(p)
	// Leading ~ -> home. Only bare ~, ~/ or ~\ (never ~user).
	if p == "~" || strings.HasPrefix(p, `~\`) || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			p = home + p[1:]
		}
	}
	return p
}

// lexCanonical is the single sanctioned filepath.Abs + filepath.Clean call
// site in the codebase (SPEC-06 §4 pipeline step 2-3).
func lexCanonical(p string) string {
	if p == "" {
		return p
	}
	if a, err := filepath.Abs(p); err == nil {
		p = a
	}
	return filepath.Clean(p)
}

// normalizeLocalUNC maps UNC spellings of local drives back to drive paths:
// \\localhost\c$\x\y, \\127.0.0.1\c$\x\y and \\?\UNC\localhost\c$\x\y
// become c:\x\y. Other UNC paths are left as canonical \\server\share form.
func normalizeLocalUNC(p string) string {
	u := strings.ReplaceAll(p, "/", `\`)
	for _, server := range []string{"localhost", "127.0.0.1"} {
		for _, prefix := range []string{`\\?\UNC\` + server + `\`, `\\\?\UNC\` + server + `\`, `\\` + server + `\`} {
			if len(u) > len(prefix)+2 && strings.HasPrefix(strings.ToLower(u), strings.ToLower(prefix)) {
				share := u[len(prefix):]
				i := strings.Index(share, `\`)
				if i < 0 {
					continue
				}
				drive, rest := share[:i], share[i+1:]
				if len(drive) == 2 && drive[1] == '$' && drive[0] >= 'a' && drive[0] <= 'z' || len(drive) == 2 && drive[1] == '$' && drive[0] >= 'A' && drive[0] <= 'Z' {
					return strings.ToUpper(drive[:1]) + `:\` + rest
				}
			}
		}
	}
	return u
}

// stripExtendedPrefix removes the \\?\ (and \\?\UNC\) prefix from a
// GetFinalPathNameByHandle result, yielding the canonical user-visible form.
func stripExtendedPrefix(p string) string {
	if strings.HasPrefix(p, `\\?\UNC\`) {
		return `\\` + p[len(`\\?\UNC\`):]
	}
	if strings.HasPrefix(p, `\\?\`) {
		return p[len(`\\?\`):]
	}
	return p
}
