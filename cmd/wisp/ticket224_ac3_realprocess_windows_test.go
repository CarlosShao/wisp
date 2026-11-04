//go:build windows

package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/session"
)

// TestTicket224ProductionSessionChangesAcrossRealOSProcesses is ticket 224 AC#3's
// "真 OS 进程" half, which the t101host-based restart test could not carry: that
// one runs a SECOND assembly inside the SAME test process, and its own JUDGEMENT
// CEILING says so (ticket224_assembly_test.go:31-46). This test boots the ACTUAL
// built binary as two separate child processes over ONE data dir, reads the
// SESSION-MINT line each writes at boot (cmd/wisp/run.go:578), and — to make the
// "永久免审通行证" claim airtight at the OS-process level — seeds a D45 grant row
// keyed on process 1's real identity and proves process 2's identity cannot see it.
//
// Why this is the load-bearing test: every D45 grant row is keyed by session_id and
// lookup is ListGrantsBySession(id). "本次会话内允许" dying with the process (A435
// clause 2) therefore reduces to "a restart mints a fresh identity, so the old
// row's key never matches". The in-process test proves the DB half (Covering does
// not cross to a new session); this test proves the OS-process identity half and
// that a real process-2 cannot read a row a real process-1 left behind.
func TestTicket224ProductionSessionChangesAcrossRealOSProcesses(t *testing.T) {
	h := t101boot(t, "") // valid config + mockllm + secret store, all in h.dir
	exe := buildWispForTest(t)
	ctx := context.Background()

	// runOnce boots one real `wisp run` child over h.dir and returns the identity
	// that child minted at boot.
	mint := func() string {
		c, cancel := context.WithTimeout(ctx, 120*time.Second)
		defer cancel()
		cmd := exec.CommandContext(c, exe, "run", "hi")
		cmd.Dir = h.dir
		// WISP_TEST_DATA_DIR is the explicit data-root injection (proc/envfork.go:36),
		// so both child processes share h.dir and therefore the same wisp.db.
		cmd.Env = append(os.Environ(), "WISP_TEST_DATA_DIR="+h.dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		runErr := cmd.Run()
		line := strings.ToValidUTF8(out.String(), "?")
		if c.Err() == context.DeadlineExceeded {
			t.Fatalf("a real `wisp run` child never exited within the 120s guard:\n%s", line)
		}
		// The task itself may fail for reasons unrelated to this ticket (it is not
		// what we measure); what MUST be present is the boot-time SESSION-MINT line,
		// which run.go writes during assembly before the task runs.
		m := sessionMintRe.FindStringSubmatch(line)
		if m == nil {
			t.Fatalf("a real `wisp run` child wrote no SESSION-MINT line (runErr=%v):\n%s", runErr, line)
		}
		return m[1]
	}

	id1 := mint()
	if !session.ID(id1).Valid() {
		t.Fatalf("process 1 minted %q, which the shape lock (session.ID.Valid) rejects", id1)
	}

	// Seed a session grant under process 1's REAL identity — as if boot 1 had
	// answered a card for fs.write on this path (the "答复⇒落行" write is AC#2's
	// subject; here we need a live-looking row that boot 2 must NOT be able to use).
	st, err := memory.Open(h.dir)
	if err != nil {
		t.Fatalf("open store over the shared data dir: %v", err)
	}
	if _, err := st.InsertGrant(ctx, memory.ApprovalGrant{
		Scope: memory.GrantScopeSession, Tool: "fs.write", Pattern: "C:/work/x.txt",
		SessionID: id1, CreatedAt: time.Now().Unix(), ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}); err != nil {
		st.Close()
		t.Fatalf("seed a grant under process 1's identity %q: %v", id1, err)
	}
	rows1, err := st.ListGrantsBySession(ctx, id1)
	if err != nil || len(rows1) != 1 {
		st.Close()
		t.Fatalf("process 1's identity should hold exactly 1 seeded grant, got %d (err %v)", len(rows1), err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close store before boot 2: %v", err)
	}

	id2 := mint()
	if id2 == id1 {
		t.Fatalf("two separate `wisp run` OS processes minted the same identity %q: 「本次会话内允许」 would then be a permanent pass across restarts (A435 clause 2 forbids it)", id1)
	}
	if !session.ID(id2).Valid() {
		t.Fatalf("process 2 minted %q, which the shape lock rejects", id2)
	}

	// The whole point: a grant row that a real process-1 left on disk is invisible
	// to the real process-2 identity, so any same-class request in boot 2 finds no
	// covering row and is asked again. That is AC#3 at the OS-process level.
	st2, err := memory.Open(h.dir)
	if err != nil {
		t.Fatalf("reopen store to inspect boot 2's view: %v", err)
	}
	defer st2.Close()
	rows1After, err := st2.ListGrantsBySession(ctx, id1)
	if err != nil || len(rows1After) != 1 {
		t.Fatalf("process 1's grant should still be on disk (audit), got %d (err %v)", len(rows1After), err)
	}
	rows2, err := st2.ListGrantsBySession(ctx, id2)
	if err != nil {
		t.Fatalf("ListGrantsBySession(process2 id): %v", err)
	}
	if len(rows2) != 0 {
		t.Fatalf("process 2's identity sees %d of process 1's grants: a restart did NOT invalidate the session grant (永久免审通行证)", len(rows2))
	}
	t.Logf("real OS processes: boot1 minted %s (holds 1 grant), boot2 minted %s (sees 0) — restart invalidates", id1, id2)
}

var sessionMintRe = regexp.MustCompile(`SESSION-MINT id=(\S+)`)
