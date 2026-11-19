package main

// Ticket 167 AC#6's teeth: the crash self-rescue exit reads the log directory
// the panic sink already writes, hands back a COPYABLE scene, and the one
// instruction names only what exists today.
//
// The census this closes (167-a2 §3/§9 rows 5-7): the panic record IS on disk
// (existing sink, no second recording path), but doctor's checks read zero of
// logs\ - "话说出口没接线". These cases pin the wire: the scan finds the record,
// the line carries the file path and the copy instruction, the quiet line says
// WHERE it looked, and the instruction never promises the plugin-disable
// restart that tickets 50/51 own (这一格明写未裁：盘上栈的截断口径、不走脱敏
// 的现场文件、禁用插件重启，都未裁).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// plantPanicRecord167 writes one sink-shaped JSONL file: an install line, a
// decoy error, one non-JSON torn line, and the panic record itself - the scan
// must key on the msg field and skip everything else.
func plantPanicRecord167(t *testing.T, dataDir, file string, recovered string) string {
	t.Helper()
	logDir := logSinkDir(dataDir)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatal(err)
	}
	line := `{"time":"2026-11-19T03:00:01Z","level":"ERROR","msg":"` + crashSceneMsg + `",` +
		`"goroutine":"agent.stream","owner":"agent","root":"t-167",` +
		`"error_class":"internal","recovered":` + jsonQuote167(recovered) + `,` +
		`"stack":"goroutine 1 [running]: wisp.(*agentRuntime).execute(...)","at":"2026-11-19T03:00:00Z"}`
	body := strings.Join([]string{
		`{"time":"2026-11-19T02:59:59Z","level":"INFO","msg":"wisp: persistent log sink installed","dir":"` + logDir + `"}`,
		`{"time":"2026-11-19T03:00:00Z","level":"WARN","msg":"some other error","err":"not a panic"}`,
		`torn line, not JSON at all`,
		line,
	}, "\n")
	p := filepath.Join(logDir, file)
	if err := os.WriteFile(p, []byte(body+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// jsonQuote167 quotes one string as a JSON value without pulling a second
// encoder into the assertion.
func jsonQuote167(s string) string {
	b := make([]byte, 0, len(s)+2)
	b = append(b, '"')
	for _, r := range s {
		switch r {
		case '"', '\\':
			b = append(b, '\\', byte(r))
		default:
			b = append(b, string(r)...)
		}
	}
	return string(append(b, '"'))
}

func Test167AC6DoctorFindsTheCrashSceneAndSaysWhatToDo(t *testing.T) {
	dir := t.TempDir()
	path := plantPanicRecord167(t, dir, "wisp-20261119-001.jsonl", "boom: 空指针在读快照时")

	hit, found := doctorCrashScenes(dir)
	if !found {
		t.Fatalf("the planted panic record at %s was not found: the scene the sink already wrote stays unreachable", path)
	}
	if hit.File != path {
		t.Errorf("hit.File = %q, want the planted file %q", hit.File, path)
	}
	if hit.Record.Recovered != "boom: 空指针在读快照时" {
		t.Errorf("recovered = %q, want the record's own cause verbatim", hit.Record.Recovered)
	}
	if hit.Record.Goroutine != "agent.stream" {
		t.Errorf("goroutine = %q, want the record's own name", hit.Record.Goroutine)
	}

	line := crashSelfRescueLine(hit)
	for _, want := range []string{
		path,                      // the copyable scene's address
		crashSceneMsg,             // the record the user must copy
		"复制",                      // the instruction is a copy instruction
		"wisp diagnostics export", // the second exit the doors already answer
		"diagnostics_opt_in",      // named with its gate, not assumed open
	} {
		if !strings.Contains(line, want) {
			t.Errorf("self-rescue line misses %q: %s", want, line)
		}
	}
	// The unruled part, pinned by absence: the instruction must NOT present the
	// plugin-disable restart as an available action (tickets 50/51 own that
	// branch; this cell 明写未裁).
	if strings.Contains(line, "禁用插件") {
		t.Errorf("self-rescue line promises the plugin-disable restart: %s", line)
	}
}

func Test167AC6QuietWhenNoCrashOnFile(t *testing.T) {
	dir := t.TempDir()
	// A logs directory with no panic record - the honest quiet, and it must say
	// where it looked, so "nothing crashed" and "nobody looked" never match.
	if err := os.MkdirAll(logSinkDir(dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, found := doctorCrashScenes(dir); found {
		t.Fatal("an empty log directory answered found=true: the scan invents a scene")
	}
	line := crashSceneQuietLine(dir)
	if !strings.Contains(line, logSinkDir(dir)) {
		t.Errorf("quiet line = %q, want it to name the directory it scanned", line)
	}
	// No logs directory at all reads the same way - a boot that never crashed
	// never wrote one, and that is not an error.
	if _, found := doctorCrashScenes(filepath.Join(dir, "absent")); found {
		t.Fatal("a missing log directory answered found=true")
	}
}

func Test167AC6NewestFileWins(t *testing.T) {
	dir := t.TempDir()
	oldPath := plantPanicRecord167(t, dir, "wisp-20261118-001.jsonl", "旧的那次崩溃")
	newPath := plantPanicRecord167(t, dir, "wisp-20261119-001.jsonl", "最新那次崩溃")
	older := time.Date(2026, 11, 18, 9, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 11, 19, 9, 0, 0, 0, time.UTC)
	if err := os.Chtimes(oldPath, older, older); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(newPath, newer, newer); err != nil {
		t.Fatal(err)
	}
	hit, found := doctorCrashScenes(dir)
	if !found {
		t.Fatal("the newest panic record was not found")
	}
	if hit.File != newPath {
		t.Errorf("hit.File = %q, want the NEWEST file %q (the current scene, not history)", hit.File, newPath)
	}
	if hit.Record.Recovered != "最新那次崩溃" {
		t.Errorf("recovered = %q, want the newest record's own cause", hit.Record.Recovered)
	}
}

// Test167AC6DoctorCarriesTheCrashSceneCheck is the wiring nail for the census's
// "话说出口没接线" finding: the instruction only exists if doctor actually calls
// the check. Unplugging the call reddens this line.
func Test167AC6DoctorCarriesTheCrashSceneCheck(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("doctor.go"))
	if err != nil {
		t.Fatalf("read doctor.go: %v", err)
	}
	for _, want := range []string{"doctorCrashScenes(dir)", "crashSelfRescueLine(hit)", "crash scene (logs)"} {
		if !strings.Contains(string(src), want) {
			t.Errorf("doctor.go no longer carries the crash-scene check (%q missing): the self-rescue instruction lost its only production caller", want)
		}
	}
}
