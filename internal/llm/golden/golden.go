// Package golden owns the Wisp golden-SSE format: ONE file format, TWO
// runners (ticket 09). The unit-test replayer (this package) serves recorded
// response bytes to a real adapter inside go test; tools/mockllm (separate
// module, stdlib-only) serves the same files over HTTP for compose and
// integration runs. Both parse with the grammar below; the format contract
// is mirrored verbatim in tools/mockllm/goldenfmt.go and its tests, so the
// two implementations cannot drift silently.
//
// Grammar (wisp golden sse v1):
//
//   - The file is UTF-8 text, LF line endings (CRLF is normalized on parse).
//
//   - Header zone: every line BEFORE the first "# @" directive is a comment
//     and ignored. Convention: first line "# wisp golden sse v1".
//
//   - Directive lines begin exactly with "# @":
//
//     # @scenario <free text>                (documentation, ignored)
//     # @response <status> [retry-after:<n>] [content-type:<mime>]
//     # @latency <ms>                        (per-chunk pacing hint)
//
//   - The body of a response is every line after its "# @response" until the
//     next "# @" directive, verbatim, joined with "\n" and terminated by one
//     trailing "\n".
//
//   - A file with no directives is a single response: status 200,
//     content-type text/event-stream, body = all non-header lines.
//
//   - A body line that must start with "# @" is escaped as "# @@".
//
// Directives are never part of a body. Retry ladders are expressed as
// multiple "# @response" sections: each successive request consumes the next
// section (that is how the recorded backoff-ladder tests work).
package golden

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// FormatHeader is the conventional first line of a golden file.
const FormatHeader = "# wisp golden sse v1"

// Response is one recorded HTTP response.
type Response struct {
	Status    int
	Header    http.Header
	Body      []byte
	LatencyMS int // per-chunk pacing hint (0 = none)
}

// Parse decodes golden file bytes into responses.
func Parse(src []byte) ([]Response, error) {
	text := strings.ReplaceAll(string(src), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var out []Response
	var cur *Response
	var bodyLines []string
	defaultLatency := 0 // a @latency seen before the first @response

	flush := func() {
		if cur == nil {
			return
		}
		if len(bodyLines) > 0 {
			cur.Body = []byte(strings.Join(bodyLines, "\n") + "\n")
		}
		out = append(out, *cur)
		cur, bodyLines = nil, nil
	}
	// ensure opens the implicit single response for directive-less files.
	ensure := func() {
		if cur == nil {
			cur = &Response{Status: 200, LatencyMS: defaultLatency}
			defaultLatency = 0
			cur.Header = http.Header{}
			cur.Header.Set("Content-Type", "text/event-stream")
		}
	}

	inHeaderZone := true
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if i == len(lines)-1 && line == "" {
			break // trailing newline artifact
		}
		if isDirective(line) {
			inHeaderZone = false
			verb, args := parseDirective(line)
			switch verb {
			case "response":
				flush()
				status, hdr, err := parseResponseArgs(args)
				if err != nil {
					return nil, fmt.Errorf("golden: line %d: %w", i+1, err)
				}
				cur = &Response{Status: status, Header: hdr, LatencyMS: defaultLatency}
				defaultLatency = 0
			case "latency":
				ms, err := strconv.Atoi(strings.TrimSpace(args))
				if err != nil || ms < 0 {
					return nil, fmt.Errorf("golden: line %d: bad @latency %q", i+1, args)
				}
				if cur != nil {
					cur.LatencyMS = ms
				} else {
					defaultLatency = ms
				}
			case "scenario":
				// documentation only
			default:
				return nil, fmt.Errorf("golden: line %d: unknown directive @%s", i+1, verb)
			}
			continue
		}
		if inHeaderZone {
			// Header zone: "#"-comment lines stay documentation; the first
			// non-comment line starts the implicit body of a directive-less
			// file.
			if strings.HasPrefix(line, "#") {
				continue
			}
			inHeaderZone = false
		}
		ensure()
		// Unescape body lines that legitimately start with "# @".
		if strings.HasPrefix(line, "# @@") {
			line = "# @" + strings.TrimPrefix(line, "# @@") // unescape: "# @@x" -> "# @x"
		}
		bodyLines = append(bodyLines, line)
	}
	flush()
	return out, nil
}

// LoadFile reads and parses a golden file from disk.
func LoadFile(path string) ([]Response, error) {
	src, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("golden: %w", err)
	}
	return Parse(src)
}

// Write encodes responses back into the golden format (recorder CLI, fixture
// generators). The output parses to an equal structure.
func Write(sb *strings.Builder, scenario string, rs []Response) {
	sb.WriteString(FormatHeader + "\n")
	if scenario != "" {
		sb.WriteString("# @scenario " + scenario + "\n")
	}
	for _, r := range rs {
		var d strings.Builder
		d.WriteString("# @response " + strconv.Itoa(r.Status))
		if ra := r.Header.Get("Retry-After"); ra != "" {
			d.WriteString(" retry-after:" + ra)
		}
		if ct := r.Header.Get("Content-Type"); ct != "" {
			d.WriteString(" content-type:" + ct)
		}
		sb.WriteString(d.String() + "\n")
		if r.LatencyMS > 0 {
			sb.WriteString("# @latency " + strconv.Itoa(r.LatencyMS) + "\n")
		}
		if len(r.Body) > 0 {
			body := strings.TrimSuffix(string(r.Body), "\n")
			for _, line := range strings.Split(body, "\n") {
				if strings.HasPrefix(line, "# @") {
					line = "# @" + line
				}
				sb.WriteString(line + "\n")
			}
		}
	}
}

// isDirective reports whether line begins a directive.
func isDirective(line string) bool {
	return strings.HasPrefix(line, "# @") && !strings.HasPrefix(line, "# @@")
}

// parseDirective splits "# @verb args..." (args may be empty).
func parseDirective(line string) (verb, args string) {
	rest := strings.TrimPrefix(line, "# @")
	rest = strings.TrimSpace(rest)
	if i := strings.IndexByte(rest, ' '); i >= 0 {
		return rest[:i], strings.TrimSpace(rest[i+1:])
	}
	return rest, ""
}

// parseResponseArgs parses "<status> [retry-after:<n>] [content-type:<mime>]".
func parseResponseArgs(args string) (int, http.Header, error) {
	hdr := http.Header{}
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("@response needs a status")
	}
	status, err := strconv.Atoi(fields[0])
	if err != nil || status < 100 || status > 599 {
		return 0, nil, fmt.Errorf("@response status %q out of range", fields[0])
	}
	for _, f := range fields[1:] {
		switch {
		case strings.HasPrefix(f, "retry-after:"):
			hdr.Set("Retry-After", strings.TrimPrefix(f, "retry-after:"))
		case strings.HasPrefix(f, "content-type:"):
			hdr.Set("Content-Type", strings.TrimPrefix(f, "content-type:"))
		default:
			return 0, nil, fmt.Errorf("@response: unknown option %q", f)
		}
	}
	if hdr.Get("Content-Type") == "" {
		if status == 200 {
			hdr.Set("Content-Type", "text/event-stream")
		} else {
			hdr.Set("Content-Type", "application/json")
		}
	}
	return status, hdr, nil
}
