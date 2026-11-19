package main

// Ticket 167 AC#6 - the crash self-rescue exit: the scene a human can copy, and
// ONE instruction that says what to do with it.
//
// WHAT ALREADY EXISTS (measured, and reused here - no second recording path was
// built). A panic in any spawned goroutine is recovered, converted to a
// structured record and written through the process's panic sink
// (internal/observe/goroutine.go's defaultPanicSink: name, owner, root, class,
// recovered, stack). The resident boots install the persistent JSONL pipeline
// (cmd/wisp/logsink.go installLogSink), so the record IS on disk, under
// <data root>\logs\wisp-*.jsonl. What was missing is the LAST MILE: doctor's
// checks never read logs\ (the 167-a2 census counted the check items that do:
// zero), so the operator - who on a double-click boot has NO console and
// therefore no stderr copy of the stack - had no way to find, let alone copy,
// the scene. This file is that last mile, and only that: it READS what the
// existing sink already wrote and tells the user where it is and what to do.
//
// WHAT THIS CELL LEAVES UNDECIDED, STATED (the dispatch's 明写未裁): the on-disk
// stack is NOT the full scene. The log pipeline bounds any single string at
// observe.MaxLoggedString before it reaches the file (internal/observe/redact.go,
// rule 4), so the file keeps the stack's head only - the full stack existed only
// on stderr, which a GUI boot does not have. Widening that bound (form 甲),
// writing a second crash file that bypasses the redaction pipeline (form 乙,
// which first has to answer whether a raw stack can carry a secret), and the
// "disable plugins and restart" recovery branch (waiting for the plugin engine,
// tickets 50/51) are ALL unruled - none of them is implemented here, none of
// them is promised in the instruction below, and this cell's closed set is
// exactly: existing sink + a doctor line that makes the existing record
// findable and copyable + one honest instruction. The instruction therefore
// never claims a plugin-disable restart works.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// crashSceneMaxFiles bounds the scan: doctor runs once per invocation and the
// log directory rolls, so the newest files are where the newest panic lives.
// The newest file that carries a record wins; older files are read only until
// one is found.
const crashSceneMaxFiles = 64

// crashSceneMsg is the msg value defaultPanicSink writes ("goroutine panic
// recovered", goroutine.go:168) - the one string the scan keys on.
const crashSceneMsg = "goroutine panic recovered"

// crashSceneRecord is one panic record as the JSONL sink wrote it: the slog
// JSON handler's own keys plus defaultPanicSink's attributes. Parsed tolerantly
// - a line that is not this shape is skipped, never guessed.
type crashSceneRecord struct {
	Time      string `json:"time"`
	Level     string `json:"level"`
	Msg       string `json:"msg"`
	Goroutine string `json:"goroutine"`
	Owner     string `json:"owner"`
	Root      string `json:"root"`
	Recovered string `json:"recovered"`
	Stack     string `json:"stack"`
	At        string `json:"at"`
}

// crashSceneHit is the scan's answer: which file holds the newest panic
// record, and the record itself verbatim off the disk.
type crashSceneHit struct {
	File   string
	Record crashSceneRecord
}

// doctorCrashScenes scans the data root's log directory for the newest goroutine
// panic record the existing sink already wrote. Files are read newest-first
// (by mod time), line by line, and the LAST matching line in the first matching
// file wins - the newest record is the current scene. found=false means no
// record exists in the scanned window, which is the honest "no crash on file",
// never an error.
func doctorCrashScenes(dataDir string) (crashSceneHit, bool) {
	dir := logSinkDir(dataDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return crashSceneHit{}, false
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	// Newest first: the scene the user needs is the most recent one.
	sort.Slice(files, func(i, j int) bool {
		fi, fj := fileModTime(files[i]), fileModTime(files[j])
		return fi.After(fj)
	})
	scanned := 0
	for _, f := range files {
		if scanned >= crashSceneMaxFiles {
			break
		}
		scanned++
		if rec, ok := lastPanicRecordIn(f); ok {
			return crashSceneHit{File: f, Record: rec}, true
		}
	}
	return crashSceneHit{}, false
}

// lastPanicRecordIn reads one JSONL file and returns its last panic record.
// A read error mid-file still answers whatever matched before it: the scene on
// disk is the scene, even when the tail of a torn file will not parse.
func lastPanicRecordIn(path string) (crashSceneRecord, bool) {
	f, err := os.Open(path)
	if err != nil {
		return crashSceneRecord{}, false
	}
	defer func() { _ = f.Close() }()
	var last crashSceneRecord
	found := false
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var rec crashSceneRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}
		if rec.Msg != crashSceneMsg {
			continue
		}
		last, found = rec, true
	}
	return last, found
}

// fileModTime is the scan's ordering read; a file that cannot be stat'ed sorts
// oldest (zero time) rather than aborting the scan.
func fileModTime(path string) time.Time {
	if fi, err := os.Stat(path); err == nil {
		return fi.ModTime()
	}
	return time.Time{}
}

// crashSelfRescueLine is THE sentence this cell owes the user: the copyable
// scene (file, the record's own time and cause, and the honest truncation
// note) plus one instruction. The instruction names only what exists today -
// copying the line, and the diagnostics export the diagnostics doors already
// answer - and never claims the plugin-disable restart that tickets 50/51 own.
func crashSelfRescueLine(hit crashSceneHit) string {
	rec := hit.Record
	return fmt.Sprintf(
		"崩溃现场可复制：%s（panic at %s，恢复于 %s，goroutine=%s，cause=%q；"+
			"盘上栈最多 %d 字（日志管线截断规则，完整现场的落盘口径未裁））。"+
			"自救指令：把该文件里这一行 goroutine panic recovered 整行复制出来交给维护者；"+
			"更完整的现场用 wisp diagnostics export 导出诊断包"+
			"（config.toml [privacy] diagnostics_opt_in 关着时先改为 true 再跑）。",
		hit.File, rec.At, rec.Time, rec.Goroutine, rec.Recovered,
		observe.MaxLoggedString)
}

// crashSceneQuietLine is the no-crash reading: the doctor line must say WHERE it
// looked, so "nothing crashed" and "nobody looked" never render alike.
func crashSceneQuietLine(dataDir string) string {
	return "没有读到崩溃记录（查过 " + logSinkDir(dataDir) + "）"
}
