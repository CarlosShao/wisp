package risk

import (
	"fmt"
	"strings"
)

// shellMetachars is the frozen metacharacter set for R6 (SPEC-06 §3 R6):
// pipe, command separators, redirection, command/command substitution, the
// cmd.exe escape caret, and control characters. Presence of ANY of these in
// any argv element -> L2.
const shellMetachars = "|;&><`$()^\n\r\x00"

// ruleShell implements R6 on the forced argv vector (SPEC-06 §10: shell.exec
// accepts argv only; string mode is [risk] allow_shell_string=true and is
// always L2). Verdicts:
//   - string mode -> L2
//   - any argv element contains a metacharacter -> L2
//   - clean argv, program in the allowlist -> L1
//   - clean argv, program not in the allowlist -> L2
//
// A nil/empty argv vector means "not a shell call" and does not fire.
func ruleShell(ctx *assessCtx) *contribution {
	if ctx.facts.ShellString {
		return &contribution{
			rules:  []RuleID{R6},
			level:  L2,
			reason: "R6: shell 字符串模式（allow_shell_string），一律 L2",
		}
	}
	if len(ctx.facts.ShellArgv) == 0 {
		return nil
	}
	if metas := findMetachars(ctx.facts.ShellArgv); metas != "" {
		return &contribution{
			rules:  []RuleID{R6},
			level:  L2,
			reason: fmt.Sprintf("R6: argv 含元字符（%s）", metas),
		}
	}
	if shellAllowlisted(ctx.facts.ShellArgv[0], ctx.facts.ShellAllowlist) {
		return &contribution{
			rules:  []RuleID{R6},
			level:  L1,
			reason: "R6: 程序命中 shell 白名单",
		}
	}
	return &contribution{
		rules:  []RuleID{R6},
		level:  L2,
		reason: fmt.Sprintf("R6: 程序不在 shell 白名单: %s", ctx.facts.ShellArgv[0]),
	}
}

// findMetachars returns the distinct metacharacters found across argv, in the
// frozen shellMetachars order, with control characters escaped for display.
func findMetachars(argv []string) string {
	found := make(map[rune]bool)
	for _, arg := range argv {
		for _, r := range arg {
			if strings.ContainsRune(shellMetachars, r) {
				found[r] = true
			}
		}
	}
	if len(found) == 0 {
		return ""
	}
	var b strings.Builder
	for _, r := range shellMetachars {
		if !found[r] {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		switch r {
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case 0:
			b.WriteString("\\0")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// shellAllowlisted matches argv[0] against the allowlist by full spelling or
// base name (last / or \ segment), case-insensitively.
func shellAllowlisted(argv0 string, allowlist []string) bool {
	base := argv0
	if i := strings.LastIndexAny(argv0, `/\`); i >= 0 {
		base = argv0[i+1:]
	}
	for _, entry := range allowlist {
		if strings.EqualFold(argv0, entry) || strings.EqualFold(base, entry) {
			return true
		}
	}
	return false
}
