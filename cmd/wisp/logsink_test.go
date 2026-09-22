package main

// Ticket 117 AC#3's placement invariant, stated as a test rather than as a
// comment: the persistent sink that carries the security notices may only ever
// write under the resolved data root of the running environment.
//
// Why this is the shape AC#3 takes while owner's question ("is a log file
// private data?", the one open item of ticket 95's seven) is still unanswered:
// sealing the log file is a decision this ticket must not make, and putting the
// notices somewhere loose "just to get them written" would turn a security
// outlet into a disclosure surface. What is NOT a decision, and what is pinned
// here, is that the file lives in the same tree ticket 95 already put under
// private-data discipline (config.toml, the DPAPI blob dir, memory.db and the
// backups all land sealed inside it), and that a missing data root is a refusal
// instead of a fallback to some world-shared temp default.
//
// PLATFORM LEG, stated because ticket 117's §六 overstated it (R-117-C, fixed by
// ticket 127 AC#2): this file carries no build tag, and being untagged is *not*
// evidence that it compiles under GOOS=linux. This package does not build for
// linux at all - `GOOS=linux go vet ./cmd/wisp/` is rc=1 with the verbatim line
//
//	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files
//
// measured identical at 6a39820 and at 8663a39, so it is this package's cgo
// shape and not anything these cases introduced - and the only CI leg that runs
// ./cmd/wisp/ is the "cmd/wisp CLI tests" step on windows-latest
// (.github/workflows/ci.yml; the ubuntu job's --scope=core never carries this
// package). So the honest reading of these three untagged cases is: they run on
// every platform this package builds on, which today is Windows alone.
// `GOOS=linux go vet ./internal/observe/` rc=0 proves observe's own Linux
// buildability and nothing more: cmd/wisp appears 0 times in that command's
// dependency closure (`GOOS=linux go list -deps ./internal/observe/ | grep -c
// CarlosShao/wisp/cmd/wisp` = 0).

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC3LogSinkLandsInsideTheEnvDataRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "wisp-data")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	dir := logSinkDir(root)
	// "logs" is spelled literally here on purpose: pinning it against
	// logDirName would let a rename of the constant move the goalpost, and the
	// name is what an operator (and `wisp slo`'s retention sweep) looks for.
	if got := filepath.Base(dir); got != "logs" {
		t.Errorf("sink dir base = %q, want the literal \"logs\"", got)
	}
	if got := filepath.Dir(dir); got != root {
		t.Errorf("sink dir parent = %q, want the data root %q", got, root)
	}
	// The one sentence that makes this a placement rule and not a string
	// comparison: the sink directory is a direct child of the data root, never
	// the root itself and never above it.
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		t.Fatal(err)
	}
	if rel != "logs" {
		t.Errorf("Rel(data root, sink dir) = %q, want %q", rel, "logs")
	}
}

// TestAC3EmptyDataRootIsARefusalNotAFallback is the fail-closed half: with no
// resolved data root the install refuses, so nothing lands in the process's
// current working directory (which for a test run is the package source tree,
// and for a stray caller is wherever the binary happened to be standing).
func TestAC3EmptyDataRootIsARefusalNotAFallback(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadDir(cwd)
	if err != nil {
		t.Fatal(err)
	}
	sink, err := installLogSink("")
	if err == nil {
		if sink != nil {
			sink.close()
		}
		t.Fatal("installLogSink(\"\") returned no error: an unresolved data root must be a refusal")
	}
	if sink != nil {
		sink.close()
		t.Errorf("installLogSink(\"\") returned a sink alongside its error: %+v", sink)
	}
	after, err := os.ReadDir(cwd)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Errorf("the refused install changed the working directory: %d entries before, %d after", len(before), len(after))
	}
}

// TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole is the other half of the
// same rule: making the notice persistent may not be what makes it disappear
// from the terminal. Before this ticket the process-default slog logger WAS the
// console, so an install that only swapped the default for a file handler would
// answer R-104-5 by producing R-104-5 one screen over - a WARN nobody reads
// live any more. installLogSink therefore fans out, and this case names both
// halves: the record is in the file, and the same record is still on stderr.
func TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole(t *testing.T) {
	dataDir := t.TempDir()
	// installLogSink reads os.Stderr as it runs, so handing it a pipe first is
	// what lets a test see the console side of the fan-out.
	oldStderr, oldDefault := os.Stderr, slog.Default()
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = pipeW
	slog.SetDefault(oldDefault)
	t.Cleanup(func() {
		os.Stderr = oldStderr
		slog.SetDefault(oldDefault)
	})

	sink, err := installLogSink(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	defer sink.close()

	const probe = "wisp117 probe: a warning the operator must still see live"
	slog.Warn(probe)
	if err := sink.pipeline.FlushNow(); err != nil {
		t.Fatal(err)
	}
	// Close the writer before draining: an open pipe would keep ReadAll blocked.
	os.Stderr = oldStderr
	if err := pipeW.Close(); err != nil {
		t.Fatal(err)
	}
	console, err := io.ReadAll(pipeR)
	if err != nil {
		t.Fatal(err)
	}
	if err := pipeR.Close(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(console), probe) {
		t.Errorf("the console saw nothing of the WARN; captured:\n%s", string(console))
	}
	if !strings.Contains(string(console), "level=WARN") {
		t.Errorf("the console line lost the level the stock default logger used to print; captured:\n%s", string(console))
	}

	lines := sinkLines117(t, sink.logDir())
	var inFile []string
	for _, l := range lines {
		var rec struct {
			Level string `json:"level"`
			Msg   string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(l), &rec); err != nil {
			t.Fatalf("line %q is not the JSON object the pipeline writes: %v", l, err)
		}
		if rec.Msg == probe {
			if rec.Level != "WARN" {
				t.Errorf("file record level = %q, want WARN", rec.Level)
			}
			inFile = append(inFile, rec.Msg)
		}
	}
	if len(inFile) != 1 {
		t.Errorf("the probe reached the file %d times, want exactly once (a fan-out, not a duplication)", len(inFile))
	}
}

// sinkLines117 returns the non-empty lines of the pipeline's own files under
// dir, refusing an empty read: "the file exists and says nothing" is not a pass.
func sinkLines117(t *testing.T, dir string) []string {
	t.Helper()
	names, err := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatalf("no .jsonl under %s", dir)
	}
	var out []string
	for _, n := range names {
		b, err := os.ReadFile(n)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(l) != "" {
				out = append(out, l)
			}
		}
	}
	return out
}
