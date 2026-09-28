package panel

// Ticket 200-r2's instruments, in the package that owns the carrier.
//
// Two accounts are measured here and neither is a word-game:
//
//   - §2(a), ticket 102's rewrite account: hand the loader's decision-maker a
//     WorkspaceView whose spelling expansion MOVED onto another tree, and prove
//     by RUNNING a real loader over real files that it opens none of them -
//     workspace tier and global tier both. The paired control (the same two
//     directories, Rewritten=false, files actually appear) is what stops that
//     from being an assertion about an empty fixture.
//   - §2(b), the carrier's states: "off", "on but nothing there", "on and
//     loaded", "refused" and "no reading yet" must not collapse into one absent
//     key, and the key's presence must be caused by a pump that has a reader -
//     which the no-reader pump below pins from the other side.
//
// Nothing here loosens anything that was already pinned: pump_test's two
// byte-level key nails still read exactly four keys, because their fixtures
// assemble a pump with no instructions reader, and TestAPumpWithoutAnInstru…
// is the case that proves that reading is the reader's absence and not a rule
// this file can switch off.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/projctx"
)

// instrClock is a fixed instant: the packet's generatedAt must never be what a
// status assertion is accidentally reading.
var instrClock = time.Unix(1700000000, 0).UTC()

// mkInstr writes one instruction file and returns its path.
func mkInstr(t *testing.T, dir, name, body string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// load runs one turn with the wiring the composition root uses. The budget is
// generous on purpose: these cases are about WHETHER a file was read, never
// about how much of it fit (projctx_test owns the budget's own readings).
func load(o projctx.Options) *projctx.Bundle {
	o.BudgetTokens = 100000
	return projctx.New(o).Turn()
}

func TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll(t *testing.T) {
	wsDir := filepath.Join(t.TempDir(), "moved")
	globalDir := filepath.Join(t.TempDir(), "wispdata")
	wsFile := mkInstr(t, wsDir, "AGENTS.md", "WS200R2WORKSPACE: 改写后的树不许把这份字送进请求。")
	globalFile := mkInstr(t, globalDir, "AGENTS.md", "WS200R2GLOBAL: 被拒绝的那一轮连全局级也不许读。")

	rewritten := WorkspaceView{
		Set:       true,
		Spelling:  filepath.Join("%WS200R2VAR%", "moved"),
		Canonical: wsDir,
		Rewritten: true,
		Reason:    "展开改写了这条路径（env），已拒绝授权改写后的树",
	}

	req := ProjectInstructionLoadRequestFor(rewritten, globalDir)
	if req.Refused == "" {
		t.Fatalf("refusal lost for a rewritten view: %+v", req)
	}
	if req.WorkspaceDir != "" || req.GlobalDir != "" {
		t.Errorf("a rewritten view still handed out directories: %+v - ticket 102's hole "+
			"is exactly that the expanded tree gets consumed anyway", req)
	}
	for _, want := range []string{"rewritten=true", wsDir, rewritten.Spelling} {
		if !strings.Contains(req.Refused, want) {
			t.Errorf("the refusal does not say %q, so nobody reading the log can tell "+
				"which tree was refused: %q", want, req.Refused)
		}
	}

	// The behavioural half: a REAL loader, handed what the request says, over
	// directories that really do hold files. This is the positive control the
	// dispatch named, not a grep for a field name.
	var logged []string
	b := load(projctx.Options{
		WorkspaceDir: req.WorkspaceDir,
		GlobalDir:    req.GlobalDir,
		Enabled:      true,
		Refused:      req.Refused,
		Log:          func(s string) { logged = append(logged, s) },
	})
	if b.Block != "" {
		t.Errorf("the refused loader still built a block: %q", b.Block)
	}
	if len(b.Files) != 0 {
		t.Errorf("the refused loader still read files: %+v", b.Files)
	}
	if b.SkipReason != projctx.SkipRewriteRefused {
		t.Errorf("SkipReason = %q, want %q: the panel cannot tell a refusal from a "+
			"config-off without it", b.SkipReason, projctx.SkipRewriteRefused)
	}
	joined := strings.Join(logged, "\n")
	if !strings.Contains(joined, "一份都没有读") || !strings.Contains(joined, "rewritten=true") {
		t.Errorf("the refusal stayed silent instead of saying why: %q", joined)
	}

	// Control: the same two directories, the same wiring, one field different.
	// Without it the block assertion above would only be an empty fixture.
	plain := rewritten
	plain.Rewritten = false
	req2 := ProjectInstructionLoadRequestFor(plain, globalDir)
	if req2.Refused != "" || req2.WorkspaceDir != wsDir || req2.GlobalDir != globalDir {
		t.Fatalf("an unrewritten view must hand out both trees, got %+v", req2)
	}
	b2 := load(projctx.Options{
		WorkspaceDir: req2.WorkspaceDir,
		GlobalDir:    req2.GlobalDir,
		Enabled:      true,
	})
	if !strings.Contains(b2.Block, "WS200R2WORKSPACE") || !strings.Contains(b2.Block, "WS200R2GLOBAL") {
		t.Fatalf("the control loaded neither file (%+v): the refusal case above proves nothing", b2.Files)
	}
	for _, p := range []string{wsFile, globalFile} {
		if !strings.Contains(b2.Block, p) {
			t.Errorf("the control block lost %s", p)
		}
	}
	if len(b2.Files) < 2 {
		t.Errorf("the control read %d files, want at least the workspace tier and the global tier",
			len(b2.Files))
	}
}

func TestTheInstructionStatesDoNotShareOneWireShape(t *testing.T) {
	dir := t.TempDir()
	mkInstr(t, dir, "AGENTS.md", "STATE200R2HEAD\n"+strings.Repeat("x", 200))
	empty := filepath.Join(t.TempDir(), "noinstructions")

	cases := []struct {
		name    string
		opts    projctx.Options
		want    string
		hasFile bool
	}{
		{
			"off",
			projctx.Options{WorkspaceDir: empty, GlobalDir: empty, Enabled: false},
			InstructionsStatusOff, false,
		},
		{
			"none_found",
			projctx.Options{WorkspaceDir: empty, GlobalDir: empty, Enabled: true},
			InstructionsStatusNoneFound, false,
		},
		{
			"loaded",
			projctx.Options{WorkspaceDir: dir, GlobalDir: empty, Enabled: true},
			InstructionsStatusLoaded, true,
		},
		{"refused", projctx.Options{
			WorkspaceDir: "", GlobalDir: "", Enabled: true,
			Refused: ProjectInstructionLoadRequestFor(WorkspaceView{
				Set: true, Canonical: dir, Rewritten: true, Spelling: filepath.Join("%V%", "dir"),
			}, empty).Refused,
		}, InstructionsStatusRefused, false},
	}

	wire := map[string]string{}
	for _, c := range cases {
		sec := InstructionsSectionFromBundle(load(c.opts))
		if sec.Status != c.want {
			t.Errorf("%s: status = %q, want %q (reason %q)", c.name, sec.Status, c.want, sec.Reason)
		}
		if sec.Files == nil {
			t.Errorf("%s: Files is nil, which marshals as null and hands the renderer a nil check", c.name)
		}
		if strings.Contains(sec.Status, " ") || sec.Status == "" {
			t.Errorf("%s: status %q is not a wire name", c.name, sec.Status)
		}
		snap := NewSnapshot(nil, nil, ComposerState{}, instrClock)
		snap.Instructions = sec
		data, err := json.Marshal(snap)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"instructions"`) {
			t.Errorf("%s: the carrier lost its key, the exact shape 200-r2 exists to remove: %s",
				c.name, data)
		}
		if c.hasFile != (len(sec.Files) > 0) {
			t.Errorf("%s: files = %+v, want any=%v", c.name, sec.Files, c.hasFile)
		}
		if c.want != InstructionsStatusLoaded && strings.Contains(string(data), "STATE200R2HEAD") {
			t.Errorf("%s: an empty state still carried a file entry: %s", c.name, data)
		}
		wire[c.name] = string(data)
	}

	// The complaint this cell closes: at r1 off / none_found / refused all
	// marshalled as the SAME absent key. Every pair must now read apart.
	for _, pair := range [][2]string{
		{"off", "none_found"},
		{"off", "loaded"},
		{"none_found", "loaded"},
		{"refused", "off"},
		{"refused", "none_found"},
	} {
		if wire[pair[0]] == wire[pair[1]] {
			t.Errorf("%s and %s marshal to the same bytes: %s", pair[0], pair[1], wire[pair[0]])
		}
	}
	if !strings.Contains(wire["off"], "已按你的配置跳过") {
		t.Errorf("the off state lost the sentence that names the switch: %s", wire["off"])
	}
	if strings.Contains(wire["none_found"], "已按你的配置跳过") {
		t.Errorf("the nothing-found state claims a switch nobody flipped: %s", wire["none_found"])
	}
	if !strings.Contains(wire["refused"], "rewritten=true") {
		t.Errorf("the refused state lost the account that refused it: %s", wire["refused"])
	}
	// The body never travels on the wire (AC#4's layering, applied to the panel).
	if strings.Contains(wire["loaded"], "STATE200R2HEAD") {
		t.Errorf("instruction BODY leaked into the snapshot wire: %s", wire["loaded"])
	}
}

func TestAPumpWithoutAnInstructionsReaderSendsNoKey(t *testing.T) {
	// §3③'s mutation control: with no reader on the pump the key is absent from
	// the bytes a real pump marshals, so its presence below is caused by the
	// reader and not by a field that is always on the wire with nothing behind it.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"instructions"`) {
		t.Errorf("a pump with no instructions reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d (%v), want the four pump_test pins for a reader-less pump",
			len(generic), generic)
	}

	// A reader that exists always sends the key, and the two readings it can
	// give (no turn yet / one file loaded) are two statuses, not one status and
	// a missing key.
	dir := t.TempDir()
	name := mkInstr(t, dir, "AGENTS.md", "READER200R2HEAD\n"+strings.Repeat("y", 200))
	var b *projctx.Bundle
	pump := NewSnapshotPump(PumpSources{
		Instructions: func() *projctx.Bundle { return b },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"instructions"`) ||
		!strings.Contains(string(data), `"status":"`+InstructionsStatusNotRun+`"`) {
		t.Errorf("a wired reader with no bundle yet lost the not_run state: %s", data)
	}

	b = load(projctx.Options{WorkspaceDir: dir, GlobalDir: dir, Enabled: true})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"status":"`+InstructionsStatusLoaded+`"`) ||
		!strings.Contains(string(data), filepath.Base(name)) {
		t.Errorf("the reader's bundle did not reach the wire: %s", data)
	}
	if strings.Contains(string(data), "READER200R2HEAD") {
		t.Errorf("instruction body leaked onto the wire: %s", data)
	}
}

func TestWithInstructionsKeepsTheR1CarrierContract(t *testing.T) {
	// projctx_test's AC#7 case (r1) builds the carrier through this function and
	// pins the file entry's keys. It is the same promise, so it is pinned here
	// too rather than only from the other package.
	dir := t.TempDir()
	p := mkInstr(t, dir, "AGENTS.md", "CARRIER200R2: 面板要能拿到这一份的清单。")
	sec := InstructionsSectionFromBundle(load(projctx.Options{WorkspaceDir: dir, Enabled: true}))
	snap := WithInstructions(NewSnapshot(nil, nil, ComposerState{}, instrClock), sec.Files)
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"instructions"`, `"tier":"project"`, `"depth":0`, `"bytes":`,
		`"source":"fs.read"`, filepath.Base(p),
	} {
		if !strings.Contains(string(data), want) {
			t.Errorf("carrier lost %q: %s", want, data)
		}
	}
	if strings.Contains(string(data), "CARRIER200R2") {
		t.Errorf("carrier carried the body, not the manifest: %s", data)
	}
	// And the copy rule: the snapshot must not alias the bundle's live slice.
	if len(sec.Files) == 0 {
		t.Fatal("premise: the carrier is empty, so the aliasing check would prove nothing")
	}
	sec.Files[0].Path = "MUTATED"
	if strings.Contains(string(mustMarshalInstr(t, snap)), "MUTATED") {
		t.Error("WithInstructions aliased the caller's slice: a later turn would rewrite an earlier packet")
	}
}

func mustMarshalInstr(t *testing.T, s Snapshot) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
