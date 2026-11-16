package main

// The local audit档 behind ticket 219 AC#5's 可找回指针.
//
// WHAT THIS FILE IS FOR. The 2026-11-15 receipt (ticket 219) measured the last
// hole in AC#5: washReplyText cut a 200-rune reason to 160, the mark said
// 「原文共 200 字」, and the dropped 40 runes had NO second copy anywhere -
// dropped was only a count, so the words were gone. AC#5 requires a truncation
// to carry a findable pointer (票 147's same-family discipline). This file is
// that pointer's target: an append-only JSONL file under the data root that
// holds one line per truncated reason, keyed by a line number the mark itself
// carries.
//
// WHAT IT IS NOT, and the two rules that keep it honest:
//
//   - It is NOT the product log. The product log is the slog sink under
//     <dataDir>\logs (logsink.go) - it goes through observe's redactor, it is
//     what diagnostics export bundles, and AC#5's hard constraint says the
//     reason's second copy must not land there. This档 lives outside logs/,
//     is written directly with os (never through slog), is never read back by
//     the product, and is deliberately NOT bundled by diagnostics export.
//   - MaxLoggedString does NOT apply here, and that is the point of the 口径
//     (the dispatch's 「⛔ 假称全文」 clause): internal/observe/redact.go:37's
//     512 bound governs strings that reach the SLOG sink, not a file this
//     package writes itself. So the line below carries the COMPLETE washed
//     original, and the mark's 「全文 N 字」 claim is checkable against it -
//     if this file ever gets rerouted through slog, that claim becomes false
//     and the mark must say 「原文另存已到 512 上限」 instead.
//
// The seq returned by append is the 1-based line number of the line just
// written (recounted from the file on every append, so a pointer minted by a
// previous run of the same process or a different process still resolves).

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// replyReasonArchiveDirName and replyReasonArchiveFileName place the档: a
// sibling of logs/, secrets/ and models/ under the data root, outside the
// product log tree so no diagnostics bundle and no CI log can pick it up.
const (
	replyReasonArchiveDirName  = "audit"
	replyReasonArchiveFileName = "approval-reply-reasons.jsonl"
)

// replyReasonArchiveRelPath is the档's data-root-relative spelling, the one the
// pointer clause in the mark uses (a relative locator, not an OS path: it
// resolves against whatever data root the reading process booted with).
func replyReasonArchiveRelPath() string {
	return replyReasonArchiveDirName + "/" + replyReasonArchiveFileName
}

// replyReasonArchivePath is the absolute spelling, used by newReplyReasonArchive
// and by the teeth that read the档 back.
func replyReasonArchivePath(dataDir string) string {
	return filepath.Join(dataDir, replyReasonArchiveDirName, replyReasonArchiveFileName)
}

// replyReasonArchiveLine is one record of the档. Original is the COMPLETE washed
// original (the same text the mark's 原文共 N 字 counts - the wash is shape-only
// per 票 216, so washed is what "原文" means throughout this family);
// OriginalRunes lets an auditor check that count without counting runes by hand.
type replyReasonArchiveLine struct {
	Seq           int    `json:"seq"`
	When          string `json:"when"`
	Corr          string `json:"corr"`
	Action        string `json:"action"`
	Route         string `json:"route"`
	Dropped       int    `json:"dropped"`
	OriginalRunes int    `json:"original_runes"`
	Original      string `json:"original"`
}

// replyReasonArchive is the档's writer. One process owns one data root, so the
// mutex covers the writers that exist; the recount-on-every-append below is
// what keeps the line numbers honest across runs.
type replyReasonArchive struct {
	mu   sync.Mutex
	path string
}

// newReplyReasonArchive wires the档 to a data root. An empty dataDir is not an
// error - it means this assembly has no root to write in, and the surface then
// says so out loud in its pointer clause instead of minting a dead pointer.
func newReplyReasonArchive(dataDir string) *replyReasonArchive {
	if dataDir == "" {
		return nil
	}
	return &replyReasonArchive{path: replyReasonArchivePath(dataDir)}
}

// append writes one line and returns its 1-based line number: the pointer the
// mark carries. r.full must be the complete washed original; an empty full is
// a caller bug and is refused rather than written as a lying record.
func (a *replyReasonArchive) append(corr, action, route string, r replyReason) (int, error) {
	if a == nil || a.path == "" {
		return 0, errors.New("答复理由档未接入")
	}
	if r.full == "" {
		return 0, errors.New("答复理由档拒绝空原文：没有截断就没有第二份，这一行会成为假档")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(a.path), 0o700); err != nil {
		return 0, err
	}
	seq, err := replyReasonArchiveLineCount(a.path)
	if err != nil {
		return 0, err
	}
	seq++
	data, err := json.Marshal(replyReasonArchiveLine{
		Seq:           seq,
		When:          time.Now().Format(time.RFC3339),
		Corr:          corr,
		Action:        action,
		Route:         route,
		Dropped:       r.dropped,
		OriginalRunes: len([]rune(r.full)),
		Original:      r.full,
	})
	if err != nil {
		return 0, err
	}
	// One Write per line under O_APPEND: the档's only writer in a process is
	// this mutex, and the recount above is what heals the number across
	// process restarts sharing the same data root.
	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return 0, err
	}
	return seq, nil
}

// replyReasonArchiveLineCount counts the complete lines already in the档 (0
// when it does not exist yet). Truncations are operator-rare, so the file is
// small and the recount costs nothing worth caching.
func replyReasonArchiveLineCount(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return bytes.Count(data, []byte("\n")), nil
}

// readReplyReasonArchive is the teeth's read face: the档 parsed line by line.
// It is deliberately in the production file (not the test) so the format the
// teeth assert against is the one constant this file defines.
func readReplyReasonArchive(dataDir string) ([]replyReasonArchiveLine, error) {
	data, err := os.ReadFile(replyReasonArchivePath(dataDir))
	if err != nil {
		return nil, err
	}
	var lines []replyReasonArchiveLine
	for _, raw := range bytes.Split(bytes.TrimRight(data, "\n"), []byte("\n")) {
		if len(bytes.TrimSpace(raw)) == 0 {
			continue
		}
		var l replyReasonArchiveLine
		if err := json.Unmarshal(raw, &l); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

// MaxLoggedString is referenced here so the 口径 above is anchored to the bound
// it talks about: the redactor's limit binds the LOG, and this档 is not the log.
var _ = observe.MaxLoggedString
