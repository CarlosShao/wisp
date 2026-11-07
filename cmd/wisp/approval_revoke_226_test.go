package main

// Ticket 226 AC#5 - the reply grammar can trigger a revocation, and the
// revocation really removes the rule the next start would have honoured.
//
// AC#5's sentence is 「撤销一条长期规则必须既能落库、又能被答复语法触发」. The
// 落库 half has been Manager.SetAllowedDirs since the write-guard landing (two
// package-level cases pin it); what was missing - and what this file pins - is
// the entrance: an operator line on the assembled `wisp run` reply stream that
// ends with the rule gone from config.toml.
//
// The two E2E cases reuse the always-branch rig (widenRun / newWidenRun) because
// a rule can only be revoked after one was stored, and storing one is the 长期
// flow's job. The choreography matters: the revoke line is typed WHILE the
// original card is still pending, so the reply listener is provably still
// reading - there is no second listener and no direct call into replySurface
// from a test. The assertions are the file, read back the way the next start
// reads it (config.LoadFile), plus the audit's booking; a nil error from the
// write path is never the evidence.
//
//	TestReplyRevokeRemovesTheStoredRuleForTheNextStart   store via always + a
//	    settled widening card, then `revoke` the exact rule text the flow
//	    printed: config.toml loses that one line and keeps the operator's own,
//	    and LoadFile - the next start's own loader - reads the list without it.
//	TestReplyRevokeRefusesARuleTheProcessDoesNotHold     the fail-closed
//	    control: a rule the process does not hold is refused out loud and the
//	    file stays byte-for-byte what the widening left.
//	TestRevokeVerbRefusalsDoNotTouchTheConfigWriter      the two refusal shapes
//	    with no writer behind them at all: no argument, and no assembled config
//	    write surface - neither may reach for a Manager.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

// driveStoreThenRevoke runs the 长期 flow to a stored rule, then sends the
// revoke line while the original card is still pending, waits for the verb's
// own sentence, and only then answers the original card. storedOut, when
// non-nil, receives the config.toml bytes as they stood after the widening -
// the proof that the rule was really stored before the revocation touched
// anything.
func driveStoreThenRevoke(t *testing.T, w *widenRun, pw *io.PipeWriter, revokeLine, wantSentence string, answerOriginal string, storedOut *[]byte) {
	t.Helper()
	defer func() { _ = pw.Close() }()
	defer close(w.reply)

	waitForSentence187(t, w.h, "卡片编号：", "the first card must show the address the operator answers against")
	if _, err := fmt.Fprintln(pw, "always "+w.corr); err != nil {
		w.setFail(err)
		return
	}
	waitForSentence187(t, w.h, "要存的规则：", "the 长期 answer must print the rule it stores")
	corr := waitWidenCard187(t, w, 15*time.Second)
	if corr == "" {
		w.setFail(errors.New("the second, L2-level re-confirmation card never appeared"))
		return
	}
	if _, err := fmt.Fprintln(pw, "yes "+corr); err != nil {
		w.setFail(err)
		return
	}
	if !waitForAudit187(w.h, "approval: REPLY WIDEN-APPLIED", 25*time.Second) {
		w.setFail(errors.New("the widening never booked WIDEN-APPLIED, so there was no rule to revoke"))
		return
	}
	if storedOut != nil {
		body, err := os.ReadFile(filepath.Join(w.h.dir, configFileName))
		if err != nil {
			w.setFail(err)
			return
		}
		*storedOut = body
	}
	if _, err := fmt.Fprintln(pw, revokeLine); err != nil {
		w.setFail(err)
		return
	}
	waitForSentence187(t, w.h, wantSentence, "the revoke verb must answer on the console")
	if _, err := fmt.Fprintln(pw, answerOriginal+" "+w.corr); err != nil {
		w.setFail(err)
		return
	}
}

func TestReplyRevokeRemovesTheStoredRuleForTheNextStart(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t226-revoke", "t226-revoke-corr")
	// The exact text the 长期 flow prints and the audit books - the spelling an
	// operator copies off the screen, spaces and all.
	rule := `[fs] allowed_dirs += "` + filepath.ToSlash(w.h.outside) + `"`
	var stored []byte
	go driveStoreThenRevoke(t, w, pw, "revoke "+rule, "已撤销长期规则", "yes", &stored)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	// The rule really was in the file between the two answers - otherwise the
	// assertions below could be green because nothing was ever stored.
	if !strings.Contains(string(stored), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("the widening never stored the rule this case revokes; config.toml was:\n%s", stored)
	}

	cfgPath := filepath.Join(w.h.dir, configFileName)
	text, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(text), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("the revoked rule is still in config.toml (AC#5: 撤销必须真落库):\n%s", text)
	}
	if !strings.Contains(string(text), filepath.ToSlash(w.h.dir)) {
		t.Fatalf("the revocation removed more than the one rule it named; config.toml:\n%s", text)
	}
	// The next start reads this file through LoadFile, so that is the reading
	// the tooth is measured on, not a substring.
	loaded, _, err := config.LoadFile(cfgPath, nil)
	if err != nil {
		t.Fatalf("the revoked config no longer loads: %v", err)
	}
	if len(loaded.FS.AllowedDirs) != 1 || loaded.FS.AllowedDirs[0] != filepath.ToSlash(w.h.dir) {
		t.Fatalf("the next start would load fs.allowed_dirs = %v, want exactly [%s]",
			loaded.FS.AllowedDirs, filepath.ToSlash(w.h.dir))
	}
	audit := w.h.err.String()
	for _, want := range []string{
		"approval: REPLY WIDEN-APPLIED",
		`approval: REPLY REVOKED corr="" tool=config.allow_dir route="native/revoke"`,
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing %q; got:\n%s", want, audit)
		}
	}
	for _, banned := range []string{"REVOKE-FAILED", "REVOKE-WRITE-FAILED", "WIDEN-REFUSED"} {
		if strings.Contains(audit, banned) {
			t.Errorf("audit booked %s on the happy path:\n%s", banned, audit)
		}
	}
}

func TestReplyRevokeRefusesARuleTheProcessDoesNotHold(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t226-revoke2", "t226-revoke2-corr")
	var stored []byte
	go driveStoreThenRevoke(t, w, pw, "revoke D:/nowhere/else", "本进程账上没有", "no", &stored)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(w.h.dir, configFileName)
	text, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != string(stored) {
		t.Fatalf("a refused revocation changed the file:\nbefore:\n%s\nafter:\n%s", stored, text)
	}
	if !strings.Contains(string(text), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("the stored rule vanished without a successful revoke:\n%s", text)
	}
	audit := w.h.err.String()
	if !strings.Contains(audit, `approval: REPLY REFUSED corr="" tool=config.allow_dir route="native/revoke"`) {
		t.Errorf("a refused revocation was not booked as a refusal:\n%s", audit)
	}
	if strings.Contains(audit, "approval: REPLY REVOKED") {
		t.Errorf("a refused revocation was booked as applied:\n%s", audit)
	}
}

// TestRevokeVerbRefusalsDoNotTouchTheConfigWriter pins the two shapes that must
// never reach the Manager: no rule named, and no config write surface
// assembled. A bare replySurface has no audit sink and no runtime, so these are
// also the shapes' spelling guarantee - the operator gets told what went wrong
// instead of a silent nothing.
func TestRevokeVerbRefusalsDoNotTouchTheConfigWriter(t *testing.T) {
	s := &replySurface{}
	if _, err := s.revokeRule(""); err == nil || !strings.Contains(err.Error(), "要撤销哪一条") {
		t.Errorf("a bare revoke must ask which rule, got: %v", err)
	}
	if _, err := s.revokeRule(`[fs] allowed_dirs += "D:/x"`); err == nil || !strings.Contains(err.Error(), "配置写面没有装配") {
		t.Errorf("a revoke with no assembled writer must say so, got: %v", err)
	}
}
