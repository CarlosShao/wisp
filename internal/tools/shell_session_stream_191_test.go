package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// Ticket 191 tool-side teeth: shell.session's edge delivery (SessionOptions.
// OnChunk). The frozen rule set this file pins:
//
//   - OnChunk nil => the Result is byte-for-byte the pre-191 one (tooth 1);
//     the whole 163 suite passing with nil callbacks is the broad regression,
//     and the paired-call comparison below is the pointed one.
//   - OnChunk set => the chunks, concatenated in delivery order, equal the
//     RAW stream the shell produced for this command - nothing dropped,
//     reordered or duplicated (tooth 2). len(chunks) >= 2 is the detector
//     arm: a reader mutated to deliver only once fails BOTH the count and
//     the concatenation ruler.
//   - The per-command cumulative chunk cap equals sessionPendingCap, the
//     SAME bound the search loop enforces on the buffer (tooth 4's ruler):
//     the consumer can never receive one byte more than the buffer bound,
//     and the cap-path refusal still fires.
//
// Mutation readouts for teeth 3 and 4 (reader-delivers-once; cap removed)
// are recorded in docs/evidence/s1/191-stream-tools-r1.md.

// piecewiseResponder answers like a real command's output: every line is
// written to the pipe in SEVERAL writes (so the session's reader performs
// several reads and the edge face must deliver several chunks), and every
// byte it wrote is recorded - the completeness ruler the concatenation
// assertion runs against. The marker line itself is split mid-token, which
// the accumulate-whole-stream search must survive exactly as before.
func piecewiseResponder(code int) (func(*fakeSessionShell, string), func() string) {
	var mu sync.Mutex
	var raw bytes.Buffer
	record := func(p []byte) {
		mu.Lock()
		raw.Write(p)
		mu.Unlock()
	}
	respond := func(f *fakeSessionShell, line string) {
		if strings.HasPrefix(line, "echo w163") {
			fields := strings.Fields(line)
			if len(fields) >= 3 && fields[2] == "%errorlevel%" {
				m := fmt.Sprintf("%s %d\r\n", fields[1], code)
				record([]byte(m[:4]))
				_, _ = f.pipeW.Write([]byte(m[:4]))
				record([]byte(m[4:8]))
				_, _ = f.pipeW.Write([]byte(m[4:8]))
				record([]byte(m[8:]))
				_, _ = f.pipeW.Write([]byte(m[8:]))
				return
			}
		}
		for _, piece := range []string{"ok: ", line, "\r\n"} {
			record([]byte(piece))
			_, _ = f.pipeW.Write([]byte(piece))
		}
	}
	return respond, func() string {
		mu.Lock()
		defer mu.Unlock()
		return raw.String()
	}
}

// streamCollector is a mutex-guarded chunk sink; Sum reads the delivery
// order back and Bytes the total handed out (the cap ruler).
type streamCollector struct {
	mu     sync.Mutex
	chunks []string
}

func (c *streamCollector) sink(chunk string) {
	c.mu.Lock()
	c.chunks = append(c.chunks, chunk)
	c.mu.Unlock()
}

func (c *streamCollector) all() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return strings.Join(c.chunks, "")
}

func (c *streamCollector) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.chunks)
}

func (c *streamCollector) bytes() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, ch := range c.chunks {
		n += len(ch)
	}
	return n
}

// sessionCallStream is sessionCall with the 191 streaming face: the third
// Execute argument is the onUpdate seam the bridge feeds (nil = no stream).
func sessionCallStream(t *testing.T, d ShellDeps, argv []string, onUpdate func(string)) Result {
	t.Helper()
	raw := fmt.Sprintf(`{"argv":%s}`, shellArgvJSON(argv))
	res, err := (shellSession{d: d}).Execute(mustTaskCtx(t), json.RawMessage(raw), onUpdate)
	if err != nil {
		t.Fatalf("Execute: unexpected provider fault: %v", err)
	}
	return res
}

// Test191StreamChunksConcatenateToRawStream - tooth 2: with OnChunk set, the
// chunks in delivery order are the command's raw stream, complete and in
// order, including the echoed command line and the (host machinery) marker
// line; the final Result stays marker-stripped as before. The len>=2 arm is
// what a reader mutated to "deliver only once" fails first.
func Test191StreamChunksConcatenateToRawStream(t *testing.T) {
	respond, raw := piecewiseResponder(0)
	ff := &fakeShellFactory{respond: respond}
	d := shellDepsWith(ff)
	col := &streamCollector{}
	res := sessionCallStream(t, d, []string{"echo", "hello"}, col.sink)
	if res.IsError {
		t.Fatalf("streaming command must run: %s", res.Text)
	}
	if got, want := res.Text, "ok: echo hello 2>&1\n退出码：0（会话内埋标记回读，每发单独判级）"; got != want {
		t.Fatalf("Result text must stay the pre-191 form:\n got %q\nwant %q", got, want)
	}
	if col.count() < 2 {
		t.Fatalf("piecewise output must arrive as SEVERAL chunks, count=%d", col.count())
	}
	if got, want := col.all(), raw(); got != want {
		t.Fatalf("concatenated chunks must equal the raw stream:\n got %q\nwant %q", got, want)
	}
	if strings.Contains(res.Text, "w163") {
		t.Fatalf("marker must stay out of the Result even while streaming: %s", res.Text)
	}
}

// Test191NilCallbackByteIdentical - tooth 1: the same command run twice on
// fresh sessions, once with the callback nil and once with it set, produces
// the SAME Result bytes. Streaming is a second face, never a change to the
// answer.
func Test191NilCallbackByteIdentical(t *testing.T) {
	respondNil, _ := piecewiseResponder(0)
	resNil := sessionCallStream(t, shellDepsWith(&fakeShellFactory{respond: respondNil}),
		[]string{"echo", "hello"}, nil)
	respondCb, _ := piecewiseResponder(0)
	resCb := sessionCallStream(t, shellDepsWith(&fakeShellFactory{respond: respondCb}),
		[]string{"echo", "hello"}, func(string) {})
	if resNil.IsError != resCb.IsError || resNil.Text != resCb.Text {
		t.Fatalf("nil and set callbacks must answer byte-identically:\n nil %q %v\n cb  %q %v",
			resNil.Text, resNil.IsError, resCb.Text, resCb.IsError)
	}
}

// floodResponder answers the command with `total` bytes of output in 1MiB
// writes and NEVER writes the marker line: the flood is the whole command,
// so the only way out is the cap-path refusal.
func floodResponder(total int) func(*fakeSessionShell, string) {
	return func(f *fakeSessionShell, line string) {
		if strings.HasPrefix(line, "echo w163") {
			return // no boundary in the stream, ever
		}
		buf := bytes.Repeat([]byte{'a'}, 1<<20)
		for sent := 0; sent < total; sent += len(buf) {
			n := len(buf)
			if total-sent < n {
				n = total - sent
			}
			_, _ = f.pipeW.Write(buf[:n])
		}
	}
}

// Test191StreamChunkCapIsTheBufferCap - tooth 4's ruler: a command whose
// output blows past sessionPendingCap still takes the SAME refusal path
// (session reset, no exit code), and the streaming consumer received at
// most sessionPendingCap bytes - the incremental face never opens a second,
// uncapped buffer. If the cap is removed (either side), the refusal text
// disappears and the delivery sum blows through the bound: this test goes
// red on both mutations.
func Test191StreamChunkCapIsTheBufferCap(t *testing.T) {
	total := sessionPendingCap + (1 << 20)
	ff := &fakeShellFactory{respond: floodResponder(total)}
	d := shellDepsWith(ff)
	col := &streamCollector{}
	ctx, cancel := shellTaskCtx(60 * time.Second)
	defer cancel()
	rawJSON := fmt.Sprintf(`{"argv":%s}`, shellArgvJSON([]string{"type", "huge.log"}))
	res, err := (shellSession{d: d}).Execute(ctx, json.RawMessage(rawJSON), col.sink)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "单发上限") {
		t.Fatalf("over-cap command must take the cap refusal: %s", res.Text)
	}
	if col.bytes() > sessionPendingCap {
		t.Fatalf("streamed chunks must stay inside the buffer cap: delivered=%d cap=%d",
			col.bytes(), sessionPendingCap)
	}
	// The stream cannot outlive the command's own budget either: the ctx
	// was a real deadline, and the refusal arrived before it. A desync-
	// timeout result here would mean the cap path never fired (that is the
	// red the cap-removal mutation produces).
	if strings.Contains(res.Text, "失去同步") {
		t.Fatalf("the cap path, not the deadline, must end an over-cap command: %s", res.Text)
	}
}
