package observe

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// --- redaction rules (SPEC-10 §5.1, ticket acceptance: five seeded classes
// must never appear in the rolling logs) ---

func newTestHandler(t *testing.T, red Redactor) (*bytes.Buffer, slog.Handler) {
	t.Helper()
	var buf bytes.Buffer
	h := &redactHandler{
		red:  red,
		next: slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}
	return &buf, h
}

func TestRedactSecretAttrKeepsLast4Only(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	slog.New(h).Info("provider call", "api_key", "sk-1234567890abcdefgh")
	out := buf.String()
	if strings.Contains(out, "sk-1234567890abcdefgh") {
		t.Fatalf("full key leaked into log: %s", out)
	}
	want := secretRedactForTest("sk-1234567890abcdefgh")
	if !strings.Contains(out, want) {
		t.Fatalf("masked key %q not found in log: %s", want, out)
	}
	// Zero-value redactor already redacts: there is no off switch.
}

// secretRedactForTest pins the C28 shape (last 4) without importing internal/secret.
func secretRedactForTest(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return "****" + s[len(s)-4:]
}

func TestRedactInlineKeyShapes(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	slog.New(h).Info("calling upstream with sk-abcdefghij1234567890 and Bearer Zx1y2z3a4b5c6d7e8f9g0h plus api_key: NTP1234567890abcd")
	out := buf.String()
	for _, leak := range []string{"sk-abcdefghij1234567890", "Zx1y2z3a4b5c6d7e8f9g0h", "NTP1234567890abcd"} {
		if strings.Contains(out, leak) {
			t.Fatalf("inline secret %q leaked into log: %s", leak, out)
		}
	}
}

func TestRedactAudioBufferNeverLogged(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	audio := make([]byte, 4096)
	for i := range audio {
		audio[i] = byte(i)
	}
	slog.New(h).Info("capture", "audio", audio, "pcm_samples", audio)
	out := buf.String()
	if bytes.Contains([]byte(out), audio[:64]) {
		t.Fatal("audio bytes leaked into log")
	}
	if !strings.Contains(out, "[audio buffer redacted: 4096 bytes]") {
		t.Fatalf("audio placeholder missing: %s", out)
	}
}

func TestRedactFetchBodyNeverLogged(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	body := "<html><body>secret-page-content-marker-9182</body></html>"
	slog.New(h).Info("web.fetch done", "body", body)
	out := buf.String()
	if strings.Contains(out, "secret-page-content-marker-9182") {
		t.Fatalf("fetch body leaked into log: %s", out)
	}
	if !strings.Contains(out, "[content redacted:") {
		t.Fatalf("content placeholder missing: %s", out)
	}
}

func TestRedactLongArgTruncated(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	long := strings.Repeat("x", 4096) + "TAIL-MARKER-777"
	slog.New(h).Info("tool call", "args", long)
	out := buf.String()
	if strings.Contains(out, "TAIL-MARKER-777") {
		t.Fatal("long arg tail leaked into log")
	}
	if !strings.Contains(out, "(truncated, 4111 chars total)") {
		t.Fatalf("truncation marker missing: %s", out)
	}
}

func TestRedactPathsOptIn(t *testing.T) {
	off := Redactor{}
	if got := off.String("opened C:\\Users\\swq\\notes.txt today"); !strings.Contains(got, "C:\\Users\\swq\\notes.txt") {
		t.Fatalf("path masking must be opt-in, got %q", got)
	}
	on := Redactor{RedactPaths: true}
	got := on.String("opened C:\\Users\\swq\\notes.txt today")
	if strings.Contains(got, "swq") {
		t.Fatalf("path leaked with redact_paths on: %q", got)
	}
	if !strings.Contains(got, "<path>") {
		t.Fatalf("path placeholder missing: %q", got)
	}
}

func TestRedactBakedInAttrsRedacted(t *testing.T) {
	buf, h := newTestHandler(t, Redactor{})
	log := slog.New(h.WithAttrs([]slog.Attr{{Key: "token", Value: slog.StringValue("ghp_abcdefghijklmnopqrstuvwxyz")}}))
	log.Info("auth ok")
	if strings.Contains(buf.String(), "ghp_abcdefghijklmnopqrstuvwxyz") {
		t.Fatalf("WithAttrs-baked token leaked: %s", buf.String())
	}
}

// --- rolling writer ---

func TestRollingWriterSizeRoll(t *testing.T) {
	dir := t.TempDir()
	w, err := newRollingWriter(dir, 0, 7)
	if err != nil {
		t.Fatal(err)
	}
	// Force a tiny threshold for the test.
	w.maxSize = 1024
	payload := bytes.Repeat([]byte("0123456789abcdef\n"), 128) // 2 KiB
	for i := 0; i < 3; i++ {
		if _, err := w.Write(payload); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	n, err := CountLogFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n < 3 {
		t.Fatalf("size roll did not create >=3 files, got %d", n)
	}
}

func TestRollingWriterDayRoll(t *testing.T) {
	dir := t.TempDir()
	day1 := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	day2 := day1.Add(25 * time.Hour)
	// The clock goes in through the constructor: the eager first open must
	// already speak the fixture's calendar, or the writer also creates a file
	// for the REAL today and the "one file per day" count is off by one on
	// every date outside {day1, day2}.
	w, err := newRollingWriterClock(dir, 0, 7, func() time.Time { return day1 })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("day one\n")); err != nil {
		t.Fatal(err)
	}
	w.now = func() time.Time { return day2 }
	if _, err := w.Write([]byte("day two\n")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	names := logFileNames(t, dir)
	if len(names) != 2 {
		t.Fatalf("want 2 files (one per day), got %v", names)
	}
	if !strings.Contains(names[0], "20260919") || !strings.Contains(names[1], "20260920") {
		t.Fatalf("day roll files wrong: %v", names)
	}
}

func TestRollingWriterRetentionSweep(t *testing.T) {
	dir := t.TempDir()
	// Both fixture days are DERIVED from one injected instant, never written as
	// a wall-clock-adjacent literal. A literal ages out of the 7-day window and
	// says nothing while it does: this case hardcoded its "fresh" file as
	// 20260918 and flipped at 2026-09-25 00:00 UTC, deterministically red on
	// every machine and every runner afterwards
	// (docs/reports/pending-and-issues.md A212 ③). TestRollingWriterDayRoll had
	// the same disease earlier, which is the only reason newRollingWriterClock
	// (logging.go:181) takes a clock at all. fixedNow is an arbitrary pinned
	// instant; nothing on this path may consult the real wall clock.
	fixedNow := time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC)
	// The window is 7 days, so the two fixtures sit 5 days clear of the edge on
	// their own sides (cutoff = fixedNow-7d = 2026-09-11 12:00 UTC). Do not
	// tighten these offsets toward 7: the margins are what keep both halves of
	// the assertion true regardless of where the real clock has walked to.
	oldDay := fixedNow.Add(-12 * 24 * time.Hour) // outside: must be pruned
	old := filepath.Join(dir, logFilePrefix+oldDay.Format(logDayLayout)+"-001"+logFileExt)
	if err := os.WriteFile(old, []byte("stale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The mtime agrees with the name: sweepLocked parses the day from the name
	// and falls back to mtime only when that fails, so a fixture that
	// contradicts itself pins neither signal.
	if err := os.Chtimes(old, oldDay, oldDay); err != nil {
		t.Fatal(err)
	}
	freshDay := fixedNow.Add(-2 * 24 * time.Hour) // inside: must survive
	fresh := filepath.Join(dir, logFilePrefix+freshDay.Format(logDayLayout)+"-001"+logFileExt)
	if err := os.WriteFile(fresh, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := newRollingWriterClock(dir, 0, 7, func() time.Time { return fixedNow })
	if err != nil {
		t.Fatal(err)
	}
	if err := w.sweepLocked(); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("stale file not pruned: %v", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Fatalf("fresh file wrongly pruned: %v", err)
	}
}

// --- pipeline integration ---

func TestLogPipelineEndToEnd(t *testing.T) {
	dir := t.TempDir()
	reg := NewRegistry()
	p, err := InitLogWithRegistry(LogConfig{Dir: dir, Level: "debug", RedactPaths: true}, reg)
	if err != nil {
		t.Fatal(err)
	}
	p.InstallAsDefault()
	t.Cleanup(func() { slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil))) })

	slog.Info("pipeline check", "api_key", "sk-this-is-a-long-seeded-key-9999")
	if err := p.FlushNow(); err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	files := logFileNames(t, dir)
	if len(files) == 0 {
		t.Fatal("no log file written")
	}
	data, err := os.ReadFile(filepath.Join(dir, files[0]))
	if err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])
	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		t.Fatalf("log line is not JSON: %q: %v", line, err)
	}
	if strings.Contains(string(data), "sk-this-is-a-long-seeded-key-9999") {
		t.Fatal("seeded key leaked through the pipeline")
	}
	// log-flusher is a named roster goroutine spawned through the registry.
	if reg.CountByName("log-flusher") != 0 {
		t.Fatal("log-flusher still live after Close")
	}
	if reg.TotalStarted() == 0 {
		t.Fatal("flusher was never spawned through the registry (D38b)")
	}
}

func TestLogPipelineConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	reg := NewRegistry()
	p, err := InitLogWithRegistry(LogConfig{Dir: dir}, reg)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		h := p.Handler()
		// Direct handler writes exercise rollingWriter locking without the
		// process-wide default (which InitLogWithRegistry deliberately does
		// not touch - tests must stay isolated).
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = h.Handle(context.Background(),
					slog.NewRecord(time.Now(), slog.LevelInfo, "concurrent", 0))
			}
		}(i)
	}
	wg.Wait()
	_ = p.FlushNow()
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
}

func logFileNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), logFilePrefix) && strings.HasSuffix(e.Name(), logFileExt) {
			names = append(names, e.Name())
		}
	}
	return names
}
