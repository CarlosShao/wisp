package main

// Ticket 226 AC#1, end to end: the operator edits config.toml while `wisp run`
// is alive, then answers one 「以后这类都别问」 - and the edit must still be in
// the file afterwards.
//
// This is the shape the ticket describes in the owner's own words, driven through
// the seam AGENTS.md §1.3 names for the CLI (the assembled run with a scripted
// reply stream). Nothing stands in for the gate, the second L2 card, the config
// writer or the file: the widening lands through Manager.AddAllowedDir, which is
// the code that used to hand the whole startup snapshot to SaveFile.
//
// Two assertions, both needed. The hand edit surviving alone would also be true
// if the write simply refused to happen, so the stored rule is read back from the
// same file in the same case - that is why AC#1 calls a green here worthless
// unless "the test changed the second key" and "the write wrote" both hold.
//
// The last assertion is the boundary this ticket does NOT cross: merging the
// file in must not put the operator's value into the running process. Reading a
// hand edit back is ticket 223's job (and a locked loosening there still goes
// through ConfirmLocked); a write path that quietly adopts is a second, wider
// behaviour change with no ticket of its own.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

func TestAC1AlwaysBranchDoesNotRevertAHandEditedKey(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t226-always", "t226-always-corr")
	cfgPath := filepath.Join(w.h.dir, configFileName)

	// The edit lands after assembly, while the process is running and nothing
	// polls: [app] theme, a key with nothing whatever to do with [fs].
	inner := w.h.rtHook
	w.h.rtHook = func(rt *agentRuntime) {
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			w.setFail(fmt.Errorf("read config.toml: %w", err))
			return
		}
		hand := string(data) + "\n[app]\ntheme = \"light\"\n"
		if err := os.WriteFile(cfgPath, []byte(hand), 0o600); err != nil {
			w.setFail(fmt.Errorf("hand edit config.toml: %w", err))
			return
		}
		inner(rt)
	}

	go driveAlways(t, w, pw, true)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	text, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	body := string(text)
	if !strings.Contains(body, filepath.ToSlash(w.h.outside)) {
		t.Errorf("the run did not store the rule it was answering for, so the surviving "+
			"hand edit proves nothing:\n%s", body)
	}
	// Read the file the way the next start would, rather than matching the
	// canonical quoting (MarshalCanonical emits single-quoted strings).
	got, _, err := config.LoadFile(cfgPath, nil)
	if err != nil {
		t.Fatalf("the file the 长期 write left behind does not load: %v\n%s", err, body)
	}
	if got.App.Theme != "light" {
		t.Errorf("the 长期 write reverted the operator's hand edit (ticket 226): app.theme = %q "+
			"on disk, want the light they wrote\n%s", got.App.Theme, body)
	}
	if !slices.Contains(got.FS.AllowedDirs, filepath.ToSlash(w.h.outside)) {
		t.Errorf("fs.allowed_dirs on disk = %v, want the stored rule %q",
			got.FS.AllowedDirs, filepath.ToSlash(w.h.outside))
	}
	if !strings.Contains(w.h.err.String(), "approval: REPLY WIDEN-APPLIED") {
		t.Errorf("the widening was not booked as applied:\n%s", w.h.err.String())
	}
	// The merge reads the file to write it; it does not read the file into memory.
	rt := w.runtime()
	if rt == nil {
		t.Fatal("no runtime was captured")
	}
	if got := rt.mgr.Config().App.Theme; got != "dark" {
		t.Errorf("the guarded write applied the hand-edited theme to the running process "+
			"(got %q, want the dark this run loaded): reading a hand edit back is ticket 223", got)
	}
}
