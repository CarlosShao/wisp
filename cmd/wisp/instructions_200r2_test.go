package main

// Ticket 200 AC#7's last hop, measured from the outside (200-r2 §2(b) 甲).
//
// At r1 the carrier was a struct field with no producer: the pump's assembly
// root had a reader for every other section of the packet and none for this
// one, so the JSON key could never arrive in a run anybody actually executed.
// That is the wire-side twin of what ledger A408 calls 文案在、控件不在, and the
// panel package's own test (TestAPumpWithoutAnInstructionsReaderSendsNoKey)
// cannot see it: only a real assembly has a pump worth reading.
//
// So nothing below constructs a panel.Snapshot or a section. Every packet these
// cases assert on came out of the run's own pump and its own exit, through the
// same rtHook shape panel_pump_test.go established, and the states are read
// from the BYTES (the wire), not from a Go struct that could hold a field the
// encoder drops.
//
// AC#8's switch is re-measured here from the panel side, because 200-r2's whole
// claim is that "off" and "nothing there" are now two different readings on the
// wire. The paired control in each case is the same fixture with the switch on:
// without it, a status string nobody flipped would pass by being wrong twice.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// instrWire is the packet's instructions section as the RENDERER sees it: three
// keys, no Go-side help.
type instrWire struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
	Files  []struct {
		Path   string `json:"path"`
		Tier   string `json:"tier"`
		Depth  int    `json:"depth"`
		Bytes  int    `json:"bytes"`
		Source string `json:"source"`
	} `json:"files"`
}

// packetInstructions reads the section back out of a run's own packet bytes and
// fails if the key is not there. The absence is the finding 200-r2 was written
// to remove, so it is a Fatal, not a skip.
func packetInstructions(t *testing.T, data []byte) instrWire {
	t.Helper()
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("the run's packet is not JSON: %v", err)
	}
	raw, ok := wire["instructions"]
	if !ok {
		t.Fatalf("the packet carries no instructions key at all: %s", data)
	}
	var sec instrWire
	if err := json.Unmarshal(raw, &sec); err != nil {
		t.Fatalf("instructions section unreadable: %v (%s)", err, raw)
	}
	if sec.Status == "" {
		t.Errorf("the section arrived with no status, which is the shape this cell "+
			"exists to close: %s", raw)
	}
	// The untrusted body never rides the manifest (ticket 200 AC#4's layering,
	// applied to the panel wire by AC#7's carrier).
	if strings.Contains(string(raw), "E2E200R2BODY") {
		t.Errorf("instruction body leaked onto the wire: %s", raw)
	}
	return sec
}

// writeInstrFile puts one instruction file in dir and returns its path.
func writeInstrFile(t *testing.T, dir, body string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "AGENTS.md")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// disableProjectInstructions writes AC#8's switch into the fixture's config. It
// appends a section the fixture does not carry, the same way nonDefaultConfig145
// does for [risk]: the point is that the packet is only right if something READ
// the file.
func disableProjectInstructions(t *testing.T, dataDir string) {
	t.Helper()
	cfgPath := filepath.Join(dataDir, configFileName)
	old, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath,
		[]byte(string(old)+"\n[agent]\nproject_instructions_enabled = false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestRunPacketCarriesTheLoadedInstructionFiles is AC#7's end-to-end half: a
// real `wisp run`, a real workspace narrowing through C26, a real AGENTS.md in
// that tree, and the packet the run published on its own.
func TestRunPacketCarriesTheLoadedInstructionFiles(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	sub := filepath.Join(f.dir, "instr-ws")
	note := writeInstrFile(t, sub, "E2E200R2HEAD\n这一份是工作区档（depth=0）的说明文件。\n"+strings.Repeat("k", 200))

	var driven *agentRuntime
	var switchErr error
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		// The narrowing goes through the production handler, which is the only
		// producer of a WorkspaceView that carries the rewrite account. Passing
		// it is the point: the loader's directory comes from that view, and from
		// panel's read of the account, not from a path this file typed.
		_, switchErr = panel.RequestWorkspaceSwitch(rt.paths, filepath.ToSlash(sub), rt.auditf)
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if switchErr != nil {
		t.Fatalf("a workspace inside the allowlist must be accepted: %v", switchErr)
	}
	if driven == nil {
		t.Fatal("the hook never saw the assembled runtime")
	}
	snap, data, seen := driven.lastPanelSnapshot()
	if !seen {
		t.Fatal("this run published no packet, so the carrier has nothing to be read from")
	}
	sec := packetInstructions(t, data)

	if sec.Status != panel.InstructionsStatusLoaded {
		t.Errorf("status = %q, want %q (reason %q)", sec.Status, panel.InstructionsStatusLoaded, sec.Reason)
	}
	if snap.Instructions == nil {
		t.Fatal("the retained snapshot carries no section although its bytes do")
	}
	found := false
	for i := range sec.Files {
		if sameTreePath(sec.Files[i].Path, note) {
			found = true
			if sec.Files[i].Tier != "project" || sec.Files[i].Depth != 0 {
				t.Errorf("the workspace file arrived as tier=%q depth=%d, want project/0",
					sec.Files[i].Tier, sec.Files[i].Depth)
			}
			if sec.Files[i].Bytes <= 0 {
				t.Errorf("the workspace file reports %d bytes, want the bytes it contributed",
					sec.Files[i].Bytes)
			}
			if sec.Files[i].Source != string(risk.SrcFSRead) {
				t.Errorf("source = %q, want the existing C25 name %q: AC#6's stamp is on "+
					"the manifest too, and a name outside the roster is a new coinage",
					sec.Files[i].Source, risk.SrcFSRead)
			}
		}
	}
	if !found {
		t.Fatalf("the packet's instructions carry no entry for %s: %+v", note, sec.Files)
	}
	if len(snap.Instructions.Files) != len(sec.Files) {
		t.Errorf("the retained snapshot holds %d entries and its bytes %d",
			len(snap.Instructions.Files), len(sec.Files))
	}
	// And the same fact on the operator's surface: the per-turn print exists,
	// naming the file the packet carries.
	if !strings.Contains(f.err.String(), "projctx:") ||
		!strings.Contains(f.err.String(), filepath.Base(note)) {
		t.Errorf("no projctx: manifest line for the loaded file on stderr:\n%s", f.err.String())
	}
	t.Logf("packet instructions section: %s", sec.asLine())
}

// TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound is AC#8 read off the
// wire, and the reverse half of §2(b): the off switch and an empty tree used to
// marshal identically (no key). Each run carries the SAME global-tier file, so
// the only thing that can move the reading is the config line.
func TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound(t *testing.T) {
	body := "E2E200R2HEAD\n这一份是全局档（数据目录里）的说明文件。\n" + strings.Repeat("m", 200)

	off := newRunFixture(t, "openai-chat")
	writeInstrFile(t, off.dir, body)
	disableProjectInstructions(t, off.dir)
	var rtOff *agentRuntime
	off.rtHook = func(rt *agentRuntime) { rtOff = rt }
	if code := off.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, off.out.String(), off.err.String())
	}
	_, dataOff, seen := rtOff.lastPanelSnapshot()
	if !seen {
		t.Fatal("the disabled run published no packet")
	}
	secOff := packetInstructions(t, dataOff)
	if secOff.Status != panel.InstructionsStatusOff {
		t.Errorf("off-switch status = %q, want %q", secOff.Status, panel.InstructionsStatusOff)
	}
	if !strings.Contains(secOff.Reason, "已按你的配置跳过") {
		t.Errorf("the off state lost the sentence that names the switch: %q", secOff.Reason)
	}
	if len(secOff.Files) != 0 {
		t.Errorf("the off state still lists files: %+v", secOff.Files)
	}
	if !strings.Contains(off.err.String(), "已按你的配置跳过") {
		t.Errorf("the disabled run stayed silent on stderr:\n%s", off.err.String())
	}

	// Control: same fixture shape, switch left at its default (on), so the very
	// same directory now reads as loaded. Without this the off assertions above
	// would only prove that nothing was ever found.
	on := newRunFixture(t, "openai-chat")
	noteOn := writeInstrFile(t, on.dir, body)
	var rtOn *agentRuntime
	on.rtHook = func(rt *agentRuntime) { rtOn = rt }
	if code := on.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, on.out.String(), on.err.String())
	}
	_, dataOn, seenOn := rtOn.lastPanelSnapshot()
	if !seenOn {
		t.Fatal("the enabled run published no packet")
	}
	secOn := packetInstructions(t, dataOn)
	if secOn.Status != panel.InstructionsStatusLoaded {
		t.Fatalf("control status = %q, want %q (reason %q): the off reading above "+
			"then only says this fixture never had a file", secOn.Status,
			panel.InstructionsStatusLoaded, secOn.Reason)
	}
	var sawGlobal bool
	for _, f := range secOn.Files {
		if sameTreePath(f.Path, noteOn) {
			sawGlobal = true
			if f.Tier != "global" || f.Depth != -1 {
				t.Errorf("the data-dir file arrived as tier=%q depth=%d, want global/-1",
					f.Tier, f.Depth)
			}
		}
	}
	if !sawGlobal {
		t.Errorf("the control's packet carries no entry for %s: %+v", noteOn, secOn.Files)
	}
	if string(dataOff) == string(dataOn) {
		t.Error("off and on produced the same bytes, so one of the two readings is a constant")
	}
	t.Logf("off packet section: %s", secOff.asLine())
	t.Logf("on  packet section: %s", secOn.asLine())
}

// sameTreePath compares two spellings of one file the way this repository's own
// workspace report does (see panel_pump_test's fold-key note): the value native
// holds for a narrowed root is internal/tools' fold key, so the packet may carry
// a lowercased or slash-flipped spelling of the same tree. Case and separator
// are compared as a tree, not as a string; no path decision is made here, which
// is why nothing is Cleaned or Abs-ed - that belongs to C26, not to a test.
func sameTreePath(a, b string) bool {
	return strings.EqualFold(filepath.ToSlash(a), filepath.ToSlash(b))
}

// asLine keeps the artifact readable from a -v run: this is the first time a
// production packet carried this section, and the evidence file quotes it.
func (s instrWire) asLine() string {
	parts := make([]string, 0, len(s.Files))
	for _, f := range s.Files {
		parts = append(parts, filepath.Base(filepath.Dir(f.Path))+"/"+filepath.Base(f.Path)+
			" tier="+f.Tier+" depth="+strconv.Itoa(f.Depth)+
			" bytes="+strconv.Itoa(f.Bytes))
	}
	return "status=" + s.Status + " reason=" + s.Reason + " files=[" + strings.Join(parts, ", ") + "]"
}
