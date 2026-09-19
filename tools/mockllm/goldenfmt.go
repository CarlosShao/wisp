package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Golden format parser - the SECOND runner of the wisp golden sse v1 format
// (the first is internal/llm/golden). The grammar is mirrored verbatim from
// that package; any change must land in both plus their tests (the
// byte-identical integration test is the tripwire).
//
// Grammar (wisp golden sse v1):
//   - UTF-8, LF; CRLF normalized.
//   - Lines before the first "# @" directive are comments.
//   - "# @response <status> [retry-after:<n>] [content-type:<mime>]" starts
//     a response; body = following lines verbatim (joined + "\n").
//   - "# @latency <ms>" paces the CURRENT or NEXT response.
//   - "# @scenario <text>" is documentation.
//   - No directives => one implicit 200 text/event-stream response.
//   - Body lines that must start with "# @" are escaped as "# @@".

type goldenResponse struct {
	status     int
	retryAfter string
	body       string // LF-joined, trailing newline
	latencyMS  int
}

func isDirective(line string) bool {
	return strings.HasPrefix(line, "# @") && !strings.HasPrefix(line, "# @@")
}

func parseGolden(src []byte) ([]goldenResponse, error) {
	text := strings.ReplaceAll(string(src), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var out []goldenResponse
	var cur *goldenResponse
	var body []string
	defaultLatency := 0

	flush := func() {
		if cur == nil {
			return
		}
		if len(body) > 0 {
			cur.body = strings.Join(body, "\n") + "\n"
		}
		out = append(out, *cur)
		cur, body = nil, nil
	}
	ensure := func() {
		if cur == nil {
			cur = &goldenResponse{status: 200, latencyMS: defaultLatency}
			defaultLatency = 0
		}
	}

	inHeader := true
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if i == len(lines)-1 && line == "" {
			break
		}
		if isDirective(line) {
			inHeader = false
			verb, args := splitDirective(line)
			switch verb {
			case "response":
				flush()
				st, ra, err := parseResponseArgs(args)
				if err != nil {
					return nil, fmt.Errorf("golden: line %d: %w", i+1, err)
				}
				cur = &goldenResponse{status: st, retryAfter: ra, latencyMS: defaultLatency}
				defaultLatency = 0
			case "latency":
				ms, err := strconv.Atoi(strings.TrimSpace(args))
				if err != nil || ms < 0 {
					return nil, fmt.Errorf("golden: line %d: bad @latency %q", i+1, args)
				}
				if cur != nil {
					cur.latencyMS = ms
				} else {
					defaultLatency = ms
				}
			case "scenario":
			default:
				return nil, fmt.Errorf("golden: line %d: unknown directive @%s", i+1, verb)
			}
			continue
		}
		if inHeader {
			if strings.HasPrefix(line, "#") {
				continue
			}
			inHeader = false
		}
		ensure()
		if strings.HasPrefix(line, "# @@") {
			line = "# @" + strings.TrimPrefix(line, "# @@") // unescape: "# @@x" -> "# @x"
		}
		body = append(body, line)
	}
	flush()
	return out, nil
}

func splitDirective(line string) (verb, args string) {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "# @"))
	if i := strings.IndexByte(rest, ' '); i >= 0 {
		return rest[:i], strings.TrimSpace(rest[i+1:])
	}
	return rest, ""
}

func parseResponseArgs(args string) (int, string, error) {
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return 0, "", fmt.Errorf("@response needs a status")
	}
	status, err := strconv.Atoi(fields[0])
	if err != nil || status < 100 || status > 599 {
		return 0, "", fmt.Errorf("@response status %q out of range", fields[0])
	}
	ra := ""
	for _, f := range fields[1:] {
		if v, ok := strings.CutPrefix(f, "retry-after:"); ok {
			ra = v
			continue
		}
		if _, ok := strings.CutPrefix(f, "content-type:"); ok {
			continue // accepted for format parity; mockllm sets headers itself
		}
		return 0, "", fmt.Errorf("@response: unknown option %q", f)
	}
	return status, ra, nil
}
