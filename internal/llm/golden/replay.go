package golden

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"
)

// The unit-test replayer: the in-process runner of the golden format. It
// serves recorded Response bytes through a real httptest.Server so adapters
// run their genuine HTTP + SSE path (no function mocks - ticket 09).
// tools/mockllm is the second runner of the same format (compose / manual
// runs) and must produce byte-identical event streams; that equivalence is
// pinned by the integration test.

// Replayer serves a parsed golden sequence: each HTTP request consumes the
// next Response; an exhausted sequence fails loudly (599) instead of
// silently repeating.
type Replayer struct {
	rs     []Response
	cursor int
	// Pace applies the per-response LatencyMS pacing when true (used by
	// cancellation/timing tests). False by default: most tests want bytes as
	// fast as possible.
	Pace bool

	// Requests records every request line seen (for attempt-count asserts).
	Requests []string
}

// NewReplayer builds a replayer over rs.
func NewReplayer(rs []Response) *Replayer { return &Replayer{rs: rs} }

// Remaining reports how many responses are left.
func (r *Replayer) Remaining() int { return len(r.rs) - r.cursor }

// ServeHTTP implements the sequential replay.
func (r *Replayer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Requests = append(r.Requests, req.Method+" "+req.URL.Path)
	if r.cursor >= len(r.rs) {
		http.Error(w, "golden: sequence exhausted ("+strconv.Itoa(len(r.rs))+" responses served)", 599)
		return
	}
	res := r.rs[r.cursor]
	r.cursor++
	for k, vs := range res.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	if res.Status != 200 {
		w.WriteHeader(res.Status)
		_, _ = w.Write(res.Body)
		return
	}
	// 200 responses stream: flush per chunk so the adapter sees genuine
	// incremental SSE.
	w.WriteHeader(200)
	if f, ok := w.(http.Flusher); ok {
		_ = f
	}
	writeBodyStreaming(w, res.Body, r.paceDelay(res))
}

// paceDelay is the per-chunk delay for the streaming write.
func (r *Replayer) paceDelay(res Response) time.Duration {
	if !r.Pace || res.LatencyMS <= 0 {
		return 0
	}
	return time.Duration(res.LatencyMS) * time.Millisecond
}

// writeBodyStreaming writes body bytes SSE-chunk-wise: every complete
// "data:"-line (including its newline) is one chunk. A trailing partial line
// (mid-JSON cut, as in disconnect fixtures) is written as-is and the
// connection closes right after - that IS the mid-stream disconnect.
func writeBodyStreaming(w io.Writer, body []byte, perChunk time.Duration) {
	if perChunk <= 0 {
		_, _ = w.Write(body)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	for _, chunk := range sseChunks(body) {
		_, _ = w.Write([]byte(chunk))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(perChunk)
	}
}

// sseChunks splits body into one chunk per SSE line group: a chunk is every
// line up to and including the next blank line, or a single "data:" line
// when blanks are rare. Simpler and sufficient for pacing: one chunk per
// non-empty line plus the blank dispatcher lines attached to the preceding
// line.
func sseChunks(body []byte) []string {
	lines := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
	var chunks []string
	var cur strings.Builder
	for _, ln := range lines {
		cur.WriteString(ln + "\n")
		if ln == "" {
			chunks = append(chunks, cur.String())
			cur.Reset()
		}
	}
	if cur.Len() > 0 {
		chunks = append(chunks, cur.String())
	}
	return chunks
}

// Server starts an httptest.Server backed by the replayer. The caller closes
// it. Panic-mode: false.
func (r *Replayer) Server() *httptest.Server {
	srv := httptest.NewServer(r)
	return srv
}

// ServerFromBytes parses src and starts the server (test convenience).
func ServerFromBytes(src []byte) (*httptest.Server, *Replayer, error) {
	rs, err := Parse(src)
	if err != nil {
		return nil, nil, err
	}
	rep := NewReplayer(rs)
	return rep.Server(), rep, nil
}

// ServerFromFile loads a golden file and starts the server.
func ServerFromFile(path string) (*httptest.Server, *Replayer, error) {
	rs, err := LoadFile(path)
	if err != nil {
		return nil, nil, err
	}
	rep := NewReplayer(rs)
	return rep.Server(), rep, nil
}

// readFile is split out so tests can stub the FS if ever needed.
func readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return b, nil
}
