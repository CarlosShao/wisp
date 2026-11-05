package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 230 AC#4, the cmd/wisp leg: the operator's STANDING B-tier exemption
// (config.toml's risk.blacklist_overrides) reaches risk.GateStanding through the
// bridge this composition root assembles, and stops reaching it when the entry
// is removed.
//
// WHY THE ALLOW IS READ OFF risk's OWN SECURITY LINE. The wiring's question is
// "does the configured set reach the gate", and the gate's only observable
// answer is the verdict it produces: an Allow writes exactly one
// "B-list OVERRIDE (single-file, source=config:risk.blacklist_overrides)" line
// naming the folded path and the rule, and a set with no exemption writes
// "B-list default DENY -> L2" instead (internal/risk/blacklist.go). Both lines
// come out of the production call site (internal/tools/mode.go's readBlacklist,
// reached by this bridge on a real Execute), so the reading is end-to-end:
// config.toml -> load pipeline -> Manager -> the armed snapshot here ->
// tools.Options.PersistedExemptions -> risk.GateStanding.
//
// WHAT THIS TEST DOES NOT CLAIM. The call itself is still refused as L2, and
// that is the shipped behaviour, not a shortfall of this leg: R3 escalates a
// B-tier path to L2 regardless of any override (internal/risk/rules_gateway.go),
// and an override record changes the gate's note, never the routing - pinned by
// internal/tools/ticket90_test.go and internal/tools/persisted_exemptions_230_test.go.
// The calls below are made WITHOUT admitting the task to the D47 gate, so the
// L2 route refuses immediately instead of parking a console card for the 300s
// deadline; nothing about the note or the security line depends on that refusal.
func Test230AC4AConfiguredExemptionReachesTheGate(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	cfgPath := filepath.Join(f.dir, configFileName)
	base, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read the fixture config: %v", err)
	}

	// A real B-tier file: ".env" is the tier's canonical example
	// (internal/risk/blacklist.go's ".env*" rule, class B).
	project := filepath.Join(f.dir, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	envFile := filepath.Join(project, ".env")
	if err := os.WriteFile(envFile, []byte("TOKEN=fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// The operator's line, spelled as a TOML literal so backslashes stay literal.
	withEntry := append(append([]byte{}, base...),
		[]byte("\n[risk]\nblacklist_overrides = ["+exemptionToml230(envFile)+"]\n")...)
	if err := os.WriteFile(cfgPath, withEntry, 0o600); err != nil {
		t.Fatal(err)
	}

	// The expectation is derived through the same resolver the load pipeline
	// uses, so the assertion is "the gate's key is the folded canonical path",
	// not "whatever this test thinks a path looks like".
	resolved, err := risk.Resolve(envFile, nil)
	if err != nil {
		t.Fatalf("premise: the fixture path cannot be resolved: %v", err)
	}
	key := risk.OverrideKey(resolved.Canonical)

	var (
		mu    sync.Mutex
		lines []string
	)
	// The gate arms risk.Logf at construction (internal/agent/approval/gate.go),
	// so the capture has to be installed after the assembly and before the
	// calls, inside the hook. Restored at the end of the test.
	prev := risk.Logf
	t.Cleanup(func() { risk.Logf = prev })
	take := func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := lines
		lines = nil
		return out
	}

	var armedAfterAssembly, armedAfterRemoval, armedAfterReAdd []string
	f.rtHook = func(rt *agentRuntime) {
		risk.Logf = func(format string, args ...any) {
			mu.Lock()
			lines = append(lines, fmt.Sprintf(format, args...))
			mu.Unlock()
		}

		// One real dispatch through the bridge this root assembled. The path is
		// the same one the config exempts; the task is deliberately NOT admitted
		// to the D47 gate (see the file comment above).
		call := func(tag string) agent.ToolOutcome {
			out, execErr := rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: tag, CorrelationID: tag + "-corr", CallID: tag,
				Name: "fs.read",
				Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(envFile))),
			})
			if execErr != nil {
				t.Fatalf("Execute(%s) returned a host fault: %v", tag, execErr)
			}
			return out
		}
		// The snapshot the closure hands the bridge, read directly as well: the
		// key must be the FOLD (OverrideKey output), and the set must carry
		// nothing else - a key spelled with the canonical string instead would
		// unlock nothing, which the security line below would then expose.
		snapshot := func() []string {
			set := rt.persistedExemptions()
			keys := make([]string, 0, len(set))
			for k, v := range set {
				if !v {
					t.Errorf("the armed set carries %q=false; a false entry is not an exemption", k)
				}
				keys = append(keys, k)
			}
			return keys
		}

		// (1) ARMED AT ASSEMBLY: the entry is in config.toml before this process
		// assembled, so the boot load must have armed it.
		armedAfterAssembly = snapshot()
		if len(armedAfterAssembly) != 1 || armedAfterAssembly[0] != key {
			t.Fatalf("the configured exemption was not armed: got %q, want exactly %q "+
				"(a key built from the canonical spelling instead of risk.OverrideKey would look like this)",
				armedAfterAssembly, key)
		}
		out := call("t230-armed")
		armedLines := take()
		var overrideLine string
		for _, l := range armedLines {
			if strings.Contains(l, "B-list OVERRIDE") {
				overrideLine = l
			}
			if strings.Contains(l, "l2-confirm") {
				t.Errorf("a config-sourced allow claims an L2 confirm: %q", l)
			}
		}
		if overrideLine == "" {
			t.Fatalf("a configured B-tier exemption did not reach risk.Gate (no override line): %q\n"+
				"the allow would have been read as 'ask first', which is what this leg wires",
				armedLines)
		}
		if !strings.Contains(overrideLine, "source=config:risk.blacklist_overrides") {
			t.Errorf("the allow line does not name the config source (the two sources were merged?): %q", overrideLine)
		}
		if !strings.Contains(overrideLine, "path="+key) {
			t.Errorf("the allow line names a different path than the folded exemption key %q: %q", key, overrideLine)
		}
		if !strings.Contains(overrideLine, "rule=") || !strings.Contains(overrideLine, ".env") {
			t.Errorf("the allow line does not carry the rule that granted it: %q", overrideLine)
		}
		// The routing is unchanged - R3 stays silence-immune, so the note is
		// evidence and never authority. This is the shipped contract, asserted
		// here so a later "wiring" change cannot silently turn the note into an
		// allow.
		if !out.IsError {
			t.Errorf("an unadmitted L2 call executed: the exemption is evidence, not an authority "+
				"over R3 (got %q)", out.Text)
		}

		// (2) REMOVED BY A RELOAD: the operator deletes the entry while the
		// process runs. Removing it is a tightening, so D36 asks for no card and
		// the reload commits by itself; the snapshot must lose the key.
		if err := os.WriteFile(cfgPath, base, 0o600); err != nil {
			t.Fatalf("rewrite config without the exemption: %v", err)
		}
		rt.reloadOnce()
		armedAfterRemoval = snapshot()
		if len(armedAfterRemoval) != 0 {
			t.Fatalf("the reload did not rebuild the exemption set: still armed %q", armedAfterRemoval)
		}
		call("t230-cleared")
		clearedLines := take()
		var denyLine string
		for _, l := range clearedLines {
			if strings.Contains(l, "B-list OVERRIDE") {
				t.Errorf("the removed exemption still allows the path: %q", l)
			}
			if strings.Contains(l, "B-list default DENY") && strings.Contains(l, "path="+key) {
				denyLine = l
			}
		}
		if denyLine == "" {
			t.Fatalf("with the entry removed the gate did not answer 'ask first' for %q: %q", key, clearedLines)
		}

		// (3) RE-ADDED BY A RELOAD: adding an exemption loosens a locked section,
		// so D36 puts the change behind one L2 card. The hook is replaced here by
		// a stand-in for the operator's click (the same device cmd/wisp's own
		// mode-switch cases use: a console run has no native channel, so the real
		// hook could only time out); what is measured is the half this leg owns -
		// that a reload which COMMITS a new entry re-arms the set.
		rt.mgr.ConfirmLocked = func(string, []string) bool { return true }
		if err := os.WriteFile(cfgPath, withEntry, 0o600); err != nil {
			t.Fatalf("re-add the exemption: %v", err)
		}
		rt.reloadOnce()
		armedAfterReAdd = snapshot()
		if len(armedAfterReAdd) != 1 || armedAfterReAdd[0] != key {
			t.Fatalf("the reload did not arm the newly added exemption: got %q, want %q", armedAfterReAdd, key)
		}
		call("t230-rearmed")
		rearmedLines := take()
		var readded string
		for _, l := range rearmedLines {
			if strings.Contains(l, "B-list OVERRIDE") && strings.Contains(l, "source=config:risk.blacklist_overrides") {
				readded = l
			}
		}
		if readded == "" {
			t.Fatalf("a reload that committed a new exemption did not reach the gate: %q", rearmedLines)
		}
	}

	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("run exit = %d, stderr:\n%s", code, f.err.String())
	}
}

// Test230AC4AnATierEntryCannotBeArmedFromConfig is the other half of the same
// wire: the A tier is absolutely denied (SPEC-06 4.1), so the composition root
// must never get the chance to arm an A-tier path - the load refuses the entry
// and this process does not assemble at all (exit 2, Unconfigured). Accepting it
// here and dropping it later would leave a config line reading as a standing
// exemption that the gate will never honour.
//
// The A-tier spelling is built under a fake profile (the same technique
// internal/config/overrides_230_test.go uses), so the real profile is untouched.
func Test230AC4AnATierEntryCannotBeArmedFromConfig(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	cfgPath := filepath.Join(f.dir, configFileName)
	base, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read the fixture config: %v", err)
	}

	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("USERPROFILE", root)
	t.Setenv("APPDATA", filepath.Join(root, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(root, "AppData", "Local"))
	sshDir := filepath.Join(root, ".ssh")
	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		t.Fatal(err)
	}
	idRsa := filepath.Join(sshDir, "id_rsa")
	if err := os.WriteFile(idRsa, []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Premise: this spelling really is the absolutely denied tier, measured with
	// the frozen table rather than assumed.
	if class := risk.Classify(idRsa); class != risk.ClassA {
		t.Fatalf("premise broken: %q classifies %v, want A", idRsa, class)
	}

	withEntry := append(append([]byte{}, base...),
		[]byte("\n[risk]\nblacklist_overrides = ["+exemptionToml230(idRsa)+"]\n")...)
	if err := os.WriteFile(cfgPath, withEntry, 0o600); err != nil {
		t.Fatal(err)
	}

	code := f.run("总结一下 这份笔记")
	if code != 2 {
		t.Fatalf("a config naming an A-tier exemption assembled anyway (exit %d); stderr:\n%s",
			code, f.err.String())
	}
	msg := f.err.String()
	if !strings.Contains(msg, "blacklist_overrides[0]") {
		t.Errorf("the refusal does not name the entry, so the operator cannot find it:\n%s", msg)
	}
	if !strings.Contains(msg, "A-tier sensitive path") {
		t.Errorf("the refusal is not the TIER verdict:\n%s", msg)
	}
}

// exemptionToml230 quotes one path as a TOML literal string: Windows paths are
// full of backslashes, and a basic string would read them as escapes.
func exemptionToml230(p string) string { return "'" + p + "'" }
