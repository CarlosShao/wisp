package main

// Ticket 197's 载体层 (leg C), measured on the production path.
//
// WHY THESE CASES SIT IN cmd/wisp AND NOT ONLY IN internal/panel. The pump's own
// package can build a roster state by hand and show it marshals correctly; that
// proves the mapping and nothing else. What ticket 197 asks - and what 200-r2
// established the shape of at run.go's PumpSources literal - is a reading out of a
// REAL assembly: spawn a subagent through the bridge this run built, and show the
// packet this run published carries it with its task id, its D43 state and its
// parent. Ledger A408's disease (文案在、控件不在) reached the wire side in 200-r1 as
// a carrier with no producer; a carrier whose only producer is a test fixture is the
// same disease one hop away, so nothing below constructs a panel.Snapshot.
//
// HOW a spawn is reached at all, and why the shape below looks like it does:
// task.spawn builds its child from rt.loopOpt, which execute() only publishes at its
// own top (cmd/wisp/run.go), and s.onRuntime runs BEFORE execute() (run.go:201). So
// nothing can spawn synchronously from the hook without failing closed, and nothing
// may call the spawner by hand either - the only honest route is the assembled
// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
// NOT the loop's: since bd124b2a callCorr mints one corr per call
// (<taskID>#<callID>, internal/agent/loop.go:603, used :676). The spawn waits
// for the ROOT's own roster row - proof it is live - then uses rt.bridge.
//
// The /__control/latency switch widens the distance between the root's own publishes
// and the child's join, so that "the packet the run published" describes a state of
// the run rather than a coin toss. Case 1 never calls publishPanelSnapshot: every
// byte it asserts on was built by the run and booked by the run's own exit, and the
// sha256 tie to the durable ledger record is checked below. Case 3 does sample, and
// says so where it does.

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// carrier197RootMarker names the root's own prompt, so a child's page can be checked
// for NOT carrying it - the merge this ticket forbids on the display side.
const carrier197RootMarker = "rootprompt197r3"

// carrier197Label goes into task.spawn's description field, so the roster row's label
// is compared against the input, not against a guess.
const carrier197Label = "载体层正控：把一句话原样说出来"

// carrier197RootTask is long enough that mockllm's echo comes back in several text
// chunks, which with the latency switch puts the run's own publishes well after the
// child's row exists.
var carrier197RootTask = "总结一下 " + carrier197RootMarker + " " + strings.Repeat("a", 360)

// carrier197ChildPrompt is the child's task book. What the child answers is
// mockllm's echo of the message the loop sent, and the loop's scene block comes
// first, so this case never asserts on the reply's wording - only that SOME text
// arrived on the child's own row and none of the root's text did.
var carrier197ChildPrompt = "把这句话原样说出来：wisp197r3"

// carrier197Row is one roster row exactly as the page would read it.
type carrier197Row struct {
	TaskID            string `json:"taskId"`
	Label             string `json:"label"`
	Kind              string `json:"kind"`
	ParentTaskID      string `json:"parentTaskId"`
	Status            string `json:"status"`
	StatusKnown       bool   `json:"statusKnown"`
	StatusReason      string `json:"statusReason"`
	StreamKey         string `json:"streamKey"`
	BlockedOnApproval bool   `json:"blockedOnApproval"`
	StreamTruncated   bool   `json:"streamTruncated"`
	StreamElidedRunes int    `json:"streamElidedRunes"`
	StreamDropped     bool   `json:"streamDropped"`
}

// carrier197Wire is the packet's tasks section as the renderer reads it: JSON keys
// only, no Go-side help.
type carrier197Wire struct {
	Rows              []carrier197Row `json:"rows"`
	InFlightSlots     int             `json:"inFlightSlots"`
	PoolCap           int             `json:"poolCap"`
	StreamTruncated   bool            `json:"streamTruncated"`
	StreamElidedRunes int             `json:"streamElidedRunes"`
	DroppedStreamKeys []string        `json:"droppedStreamKeys"`
}

// packetTasks197 reads the roster section back out of a run's own bytes. A missing
// key is Fatal, not a skip: "the carrier exists but never arrived" is the finding this
// leg was written to close.
func packetTasks197(t *testing.T, data []byte) carrier197Wire {
	t.Helper()
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("the run's packet is not JSON: %v", err)
	}
	raw, ok := wire["tasks"]
	if !ok {
		t.Fatalf("the packet carries no tasks key at all: %s", data)
	}
	var sect carrier197Wire
	if err := json.Unmarshal(raw, &sect); err != nil {
		t.Fatalf("tasks section unreadable: %v (%s)", err, raw)
	}
	// The byte-level reading of the positive wire shape, so §② of the evidence
	// file quotes bytes this run published rather than a Go struct it re-marshalled.
	t.Logf("tasks wire bytes: %s", raw)
	return sect
}

// carrierRow returns the row for one task id, or fails naming every row present: an
// absent row is the lie this ticket is about, so it is never a soft reading.
func carrierRow(t *testing.T, sect carrier197Wire, taskID string) carrier197Row {
	t.Helper()
	for _, r := range sect.Rows {
		if r.TaskID == taskID {
			return r
		}
	}
	t.Fatalf("the packet carries no roster row for %q: %+v", taskID, sect.Rows)
	return carrier197Row{}
}

// spawn197 records what one or more spawns dispatched from a goroutine learned, so no
// t.Fatal ever runs on a non-test goroutine (which would abort the wrong thing).
type spawn197 struct {
	mu       sync.Mutex
	rootID   string
	childIDs []string
	err      error
	done     chan struct{}
}

func newSpawn197() *spawn197 { return &spawn197{done: make(chan struct{})} }

// launch waits for the root's roster row, then dispatches `repeats` task.spawn calls
// through the bridge this runtime assembled, each with prompt as the child's task
// book. Sequential on purpose: the pool cap (4) is ticket 211's business, and case 3
// needs stream KEYS, which is a different bound.
func (s *spawn197) launch(rt *agentRuntime, prompt string, repeats int) {
	defer close(s.done)
	rootID, err := waitRootRow197(rt)
	if err != nil {
		s.setErr(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()
	for i := 0; i < repeats; i++ {
		desc := carrier197Label
		if repeats > 1 {
			// A long, distinct label per child: it is what the row reads as, and it
			// is the text a merge would leak through.
			desc = fmt.Sprintf("溢出正控 %02d %s", i, strings.Repeat("l", 90))
		}
		args, merr := json.Marshal(map[string]string{"description": desc, "prompt": prompt})
		if merr != nil {
			s.setErr(merr)
			return
		}
		out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),
			Name: "task.spawn", Args: args,
		})
		if eerr != nil {
			s.setErr(eerr)
			return
		}
		if out.IsError {
			s.setErr(fmt.Errorf("task.spawn 在被派发的桥上回了错误：%s", out.Text))
			return
		}
		s.mu.Lock()
		s.rootID = rootID
		s.childIDs = append(s.childIDs, childRowAfterJoin197(rt, rootID, i+1))
		s.mu.Unlock()
	}
}

func (s *spawn197) setErr(err error) {
	s.mu.Lock()
	s.err = err
	s.mu.Unlock()
}

// waitRootRow197 blocks until execute() has filed the root's own roster row, i.e.
// until there is a live loop to derive a child from. Error, never Fatal: it runs on
// the spawn goroutine.
func waitRootRow197(rt *agentRuntime) (string, error) {
	for i := 0; i < 2400; i++ {
		for _, row := range rt.taskRosterState().Rows {
			if row.Kind == tools.TaskKindRoot && row.TaskID != "" {
				return row.TaskID, nil
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	return "", fmt.Errorf("execute() 从未登记根任务行：这一程没有可派生的活环路")
}

// childRowAfterJoin197 picks the n-th subagent row hanging under rootID. Ordering is
// by task id (what panel.TaskRosterSectionFrom sorts on too), so "the first child"
// means the same thing here as it does on the wire.
func childRowAfterJoin197(rt *agentRuntime, rootID string, n int) string {
	var kids []string
	for _, row := range rt.taskRosterState().Rows {
		if row.Kind == tools.TaskKindSubagent && row.ParentTaskID == rootID {
			kids = append(kids, row.TaskID)
		}
	}
	if n <= 0 || n > len(kids) {
		return ""
	}
	return kids[n-1]
}

// TestRunPacketCarriesTheSubagentItsRosterRowFed is cell (a) plus the page half of
// cell (b): 派一枚子代理 ⇒ 这枚 run 自己发布的快照里有它，带 taskID、D43 状态、父任务和
// 它自己那一页的键；那一页上有它流出来的正文；它收口之后名册与流都改口，不是常量。
func TestRunPacketCarriesTheSubagentItsRosterRowFed(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	f.srv.Control(t, "/__control/latency", `{"ms":80}`)
	sp := newSpawn197()
	var driven *agentRuntime
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		go sp.launch(rt, carrier197ChildPrompt, 1)
	}
	if code := f.run(carrier197RootTask); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}

	// Nothing above called publishPanelSnapshot, so the retained packet is the last
	// one the RUN built. The digest check below ties it to the durable record the
	// run's own exit wrote, which is what makes this a production reading rather than
	// a snapshot of a test's own call (the standard 145's pump test holds itself to).
	snap, data, seen := driven.lastPanelSnapshot()
	if !seen {
		t.Fatal("no packet was retained")
	}
	if driven.snapshotCount() < 1 {
		t.Fatal("this run published nothing, so the packet below proves nothing")
	}
	want := sha256Short(data)
	var booked bool
	for _, rec := range ledgerSummaries145(t, f.dir) {
		if rec["sha256"] == want {
			booked = true
		}
	}
	if !booked {
		t.Errorf("the packet under test (%s) was never booked by a production publish", want)
	}

	<-sp.done
	sp.mu.Lock()
	err, rootID, childIDs := sp.err, sp.rootID, append([]string(nil), sp.childIDs...)
	sp.mu.Unlock()
	if err != nil {
		t.Fatalf("the spawn call failed: %v", err)
	}
	if rootID == "" || len(childIDs) != 1 || childIDs[0] == "" {
		t.Fatalf("the roster this run filled gave no root/child pair (root %q, children %v)", rootID, childIDs)
	}
	childID := childIDs[0]

	sect := packetTasks197(t, data)
	if len(sect.Rows) != 2 {
		t.Fatalf("rows = %d, want the root and its one child: %+v", len(sect.Rows), sect.Rows)
	}
	for i := 1; i < len(sect.Rows); i++ {
		if sect.Rows[i-1].TaskID > sect.Rows[i].TaskID {
			t.Errorf("rows are not ordered by taskId (%s before %s): the ledger books a sha256 of "+
				"these bytes, so an unordered section is a packet nobody can re-derive",
				sect.Rows[i-1].TaskID, sect.Rows[i].TaskID)
		}
	}

	child := carrierRow(t, sect, childID)
	if child.Kind != tools.TaskKindSubagent {
		t.Errorf("kind = %q, want %q as the spawner filed it", child.Kind, tools.TaskKindSubagent)
	}
	if child.ParentTaskID != rootID {
		t.Errorf("parentTaskId = %q, want the root this run ran (%q)", child.ParentTaskID, rootID)
	}
	if child.Label != carrier197Label {
		t.Errorf("label = %q, want the description handed to task.spawn (%q)", child.Label, carrier197Label)
	}
	if !child.StatusKnown {
		t.Errorf("statusKnown = false for a row the spawner itself filed (reason %q)", child.StatusReason)
	}
	if !statemachine.Valid(statemachine.State(child.Status)) {
		t.Errorf("status %q is not one of D43's 20 names: ticket 197 §2(c) forbids a second "+
			"vocabulary, and this is the wire where one would appear", child.Status)
	}
	if child.StreamKey != panel.SubagentStreamKey(childID) {
		t.Errorf("streamKey = %q, want %q - this is the field that answers 这一行是谁的",
			child.StreamKey, panel.SubagentStreamKey(childID))
	}
	if child.BlockedOnApproval {
		t.Error("a child with no card on the queue read as blocked: the join is by correlation " +
			"id, so a false positive here would dress an idle row up as a waiting one")
	}

	// The parent row has to be in the same packet or the page cannot say who derived
	// whom. Its own status dimension is filed by nobody (ruling A394 leaves that
	// rewiring to ticket 196), so the honest reading is empty WITH a reason.
	root := carrierRow(t, sect, rootID)
	if root.Kind != tools.TaskKindRoot {
		t.Errorf("root row kind = %q, want %q", root.Kind, tools.TaskKindRoot)
	}
	if root.StatusKnown || root.Status != "" || root.StatusReason == "" {
		t.Errorf("root status = %q/known=%v/reason=%q, want the empty stated WITH a reason",
			root.Status, root.StatusKnown, root.StatusReason)
	}

	// Ticket 211's account becomes numbers on the wire rather than a sentence in a
	// Go comment: the cap the refusal measures against is the bridge's ceiling.
	if sect.PoolCap != tools.MaxConcurrentSubagents {
		t.Errorf("poolCap = %d, want tools.MaxConcurrentSubagents (%d)", sect.PoolCap, tools.MaxConcurrentSubagents)
	}

	// The child's own page, found in the same packet's results section by the key the
	// roster row names. Before this leg the child loop had NO sink at all, so the only
	// text in here was the spawner's lifecycle lines (197-r2 §6.1 named it).
	page := carrierPage(t, snap.Results, child.StreamKey)
	if !strings.Contains(page.Text, "已派生") {
		t.Errorf("the child's page lost its lifecycle line: %q", page.Text)
	}
	if !strings.Contains(page.Text, "echo:") {
		t.Errorf("the child's own streamed reply never reached its own row: %q - this is the hop "+
			"the carrier leg exists to close", page.Text)
	}
	if strings.Contains(page.Text, carrier197RootMarker) {
		t.Error("the child's page carries the root's prompt: two tasks merged into one row")
	}

	// After the join, the same production reader must have CHANGED its answer about
	// that row - a status field that prints one word whatever happens is exactly the
	// 空转 ticket 181 AC#7 was filed over.
	after := driven.taskRosterState()
	joined := findRow197(t, after, childID)
	if joined.Status != string(statemachine.StateSettling) || !joined.StatusKnown {
		t.Errorf("a joined child reads %q/known=%v, want %s: the finish watcher's own write is "+
			"what this dimension is supposed to be reporting",
			joined.Status, joined.StatusKnown, string(statemachine.StateSettling))
	}
	if after.InFlightSlots != 0 {
		t.Errorf("inFlightSlots = %d after the join, want 0 (the slot is released before the row "+
			"is finalised, and a reader that sees the settled row must never see it held)",
			after.InFlightSlots)
	}
	// The packet's OWN two halves must not contradict each other. This reads one
	// instant twice - the row's status and the pool count came out of the same
	// taskRosterState() call that built these bytes - so the invariant is real:
	// the slot is released before the row is finalised, hence a packet that shows
	// its child as Settling can never report a held slot.
	//
	// What this is NOT: it is not a claim about which instant the run happened to
	// publish. 197-r3 wrote the opposite assertion here (`InFlightSlots != 0 &&
	// <a status read after the spawn call returned>`) and comparing two different
	// instants is not an invariant - it went red in a full-package run, where the
	// last packet the run booked was sampled mid-child (slots=1, status=Thinking,
	// the child joined afterwards), and green in an isolated -run of the same
	// bytes. Its message also named the reading its condition had just excluded.
	if child.Status == string(statemachine.StateSettling) && sect.InFlightSlots != 0 {
		t.Errorf("one packet reports %s's row as %s while the SAME packet reports the pool still "+
			"holding %d slot(s): the slot is released before the row is finalised, so a single "+
			"reading that shows both is two readers disagreeing about one instant",
			child.TaskID, child.Status, sect.InFlightSlots)
	}
	closed := carrierPage(t, driven.stream.Chunks(), child.StreamKey)
	if !closed.Done {
		t.Error("the child's stream row is not closed after it joined: it would stream forever " +
			"on the page, which is the failure consoleSink's EvDone branch exists to prevent")
	}
	t.Logf("packet tasks section: rows=%d poolCap=%d child=%s runStatus=%s afterStatus=%s key=%s bytes=%d",
		len(sect.Rows), sect.PoolCap, childID, child.Status, joined.Status, child.StreamKey, len(data))
}

// carrierPage finds the results row a roster row points at. Absent is Fatal: a row
// whose page nobody can find is the "会有流却没人读得到" failure the key pins exist for.
func carrierPage(t *testing.T, results []panel.ResultChunk, streamKey string) panel.ResultChunk {
	t.Helper()
	for _, chunk := range results {
		if chunk.CorrelationID == streamKey {
			return chunk
		}
	}
	t.Fatalf("no results row carries correlationId %q, so clicking that row shows nothing: %+v",
		streamKey, results)
	return panel.ResultChunk{}
}

func findRow197(t *testing.T, sect panel.TaskRosterState, taskID string) panel.TaskRow {
	t.Helper()
	for _, row := range sect.Rows {
		if row.TaskID == taskID {
			return row
		}
	}
	t.Fatalf("the reader lost row %q: %+v", taskID, sect.Rows)
	return panel.TaskRow{}
}

// TestSubagentStreamKeyHasOneMintSite is the tooth the 载体层 grew when it collapsed
// the duplicated key spelling. internal/tools and internal/panel each used to write
// the literal themselves, and the composition root could only compare the two
// constants afterwards - nothing could stop a third copy, or a fourth. Now the literal
// lives in internal/streamkey alone, and this scan is what keeps it that way: mint it
// anywhere else and this leg goes red on its own.
//
// The scan reads STRING LITERALS through the go/ast, not bytes: a comment may name the
// shape (several do) and must not be mistaken for a mint site, while a real literal in
// a non-test file is one. Test files are out of scope on purpose - a test naming the
// expected spelling is how the shape stays pinned at its source (internal/tools'
// Test197StreamKeyShapeIsLiteral and internal/panel's key-shape leg both do), and
// forbidding that would only push the literal into a computed string nobody can read.
func TestSubagentStreamKeyHasOneMintSite(t *testing.T) {
	root := repoRoot197(t)
	want := filepath.ToSlash(filepath.Join("internal", "streamkey", "streamkey.go"))
	found := stringLitSites197(t, root, "internal")
	found = append(found, stringLitSites197(t, root, "cmd")...)
	if len(found) != 1 {
		t.Fatalf("the stream key literal %q is minted in %d non-test sources %v, want exactly one: %s",
			panel.SubagentStreamKeyPrefix, len(found), found, want)
	}
	if found[0] != want {
		t.Fatalf("the one mint site is %s, want %s", found[0], want)
	}
	// The two legs still alias it, so the writer's key is the reader's key - compared
	// against each other, not against a spelling this file typed.
	if tools.SubagentStreamKeyPrefix != panel.SubagentStreamKeyPrefix {
		t.Errorf("tools = %q, panel = %q: the aliases drifted", tools.SubagentStreamKeyPrefix,
			panel.SubagentStreamKeyPrefix)
	}
	for _, id := range []string{"task-197-a", "带中文的id"} {
		if got, wantKey := tools.SubagentStreamKey(id), panel.SubagentStreamKey(id); got != wantKey {
			t.Errorf("id %q: tools minted %q, panel reads %q", id, got, wantKey)
		}
	}
}

// stringLitSites197 returns the files under root/sub holding a top-level directory
// whose Go string literals include exactly `value` - the literal itself, not a
// sentence about it.
func stringLitSites197(t *testing.T, root, sub string) []string {
	t.Helper()
	var hits []string
	err := filepath.WalkDir(filepath.Join(root, sub), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return err
		}
		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, p, nil, 0)
		if perr != nil {
			// An unparseable file is its own finding: a scanner that skips one silently
			// is how a re-fork walks in. Name it and fail.
			return fmt.Errorf("parse %s: %w", p, perr)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			v := strings.Trim(lit.Value, "`\"")
			if v == panel.SubagentStreamKeyPrefix {
				rel, relErr := filepath.Rel(root, p)
				if relErr != nil {
					t.Errorf("rel %s: %v", p, relErr)
					return true
				}
				hits = append(hits, filepath.ToSlash(rel))
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hits
}

func repoRoot197(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the working directory - cannot locate the repository root")
		}
		dir = parent
	}
}

// TestRunPacketReportsTheStreamLogPastItsBound is the truncation half of cell (b) on a
// real assembly: enough streams to cross StreamLog's key bound, then a reading of the
// packet that says the same thing the live log says.
//
// This case DOES sample (publishPanelSnapshot is the run's own publisher), because
// crossing the bound needs 33 spawns and the root's own turn is long over by then. The
// claim is correspondingly narrower than case 1's: it shows the section's numbers are
// READ from the log rather than shaped, and that crossing the bound is visible at all.
// The numeric half - a row actually losing its middle - is not reachable here and is
// said so in the list at the bottom of this file.
func TestRunPacketReportsTheStreamLogPastItsBound(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	sp := newSpawn197()
	var driven *agentRuntime
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		go sp.launch(rt, carrier197ChildPrompt, panel.DefaultStreamKeys+1)
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}
	<-sp.done
	sp.mu.Lock()
	err := sp.err
	sp.mu.Unlock()
	if err != nil {
		t.Fatalf("the fan-out failed: %v", err)
	}

	driven.publishPanelSnapshot()
	snap, data, seen := driven.lastPanelSnapshot()
	if !seen {
		t.Fatal("no packet")
	}
	sect := packetTasks197(t, data)
	if len(driven.stream.Chunks()) <= panel.DefaultStreamKeys {
		t.Fatalf("only %d streams were opened; the bound (%d) was never crossed, so the reading "+
			"below would be vacuous", len(driven.stream.Chunks()), panel.DefaultStreamKeys)
	}
	if !sect.StreamTruncated || !driven.stream.Truncated() {
		t.Errorf("packet says truncated=%v and the live log says %v, want both true after %d streams",
			sect.StreamTruncated, driven.stream.Truncated(), panel.DefaultStreamKeys+1)
	}
	if sect.StreamElidedRunes != driven.stream.ElidedRunes() {
		t.Errorf("section total %d, want the live log's own %d: the reader has to be the log",
			sect.StreamElidedRunes, driven.stream.ElidedRunes())
	}
	if strings.Join(sect.DroppedStreamKeys, ",") != strings.Join(driven.stream.DroppedKeys(), ",") {
		t.Errorf("packet names %v dropped, want the log's own %v", sect.DroppedStreamKeys,
			driven.stream.DroppedKeys())
	}

	// Every child keeps its OWN row, under its own key, with nobody else's label in
	// it - the carrier-side reading of the invariant 197-r2 made StreamLog keep.
	seenKey := map[string]bool{}
	for _, row := range sect.Rows {
		if row.Kind != tools.TaskKindSubagent {
			continue
		}
		if seenKey[row.StreamKey] {
			t.Errorf("two roster rows share stream key %q: one page for two agents", row.StreamKey)
		}
		seenKey[row.StreamKey] = true
		page := carrierPage(t, snap.Results, row.StreamKey)
		if strings.Contains(page.Text, carrier197RootMarker) {
			t.Errorf("row %s's page carries the root's prompt", row.TaskID)
		}
		if live := driven.stream.TruncationFor(row.StreamKey); row.StreamElidedRunes != live.ElidedRunes ||
			row.StreamTruncated != live.Truncated || row.StreamDropped != live.Dropped {
			t.Errorf("row %s reads truncated=%v/elided=%d/dropped=%v, want the log's own %+v",
				row.TaskID, row.StreamTruncated, row.StreamElidedRunes, row.StreamDropped, live)
		}
	}
	if len(seenKey) != panel.DefaultStreamKeys+1 {
		t.Errorf("roster carries %d subagent rows, want one per stream opened (%d)",
			len(seenKey), panel.DefaultStreamKeys+1)
	}
	t.Logf("bound crossed: rows=%d streams=%d truncated=%v elided=%d dropped=%v",
		len(sect.Rows), len(driven.stream.Chunks()), sect.StreamTruncated,
		sect.StreamElidedRunes, sect.DroppedStreamKeys)
}

// WHAT THIS FILE DID NOT MEASURE.
//
//  1. blockedOnApproval is never asserted on a production packet. The join has two
//     producers (subagent_roster_197.go:214, :219-220): a card's corr, never a task id
//     for the loop (loop.go:603), and L1 windows no run feeds; no tool_choice is sent, so mockllm
//     never answers with a tool call, so a spawned child never reaches the gate. The
//     field is therefore asserted at the pump against a real ApprovalCardView
//     (internal/panel's TestTheRosterReaderPutsSubagentsOnTheWire), and on this path it
//     is measured only as false.
//  2. A stream row that actually loses its middle is NOT reachable here. mockllm caps an
//     answer at 400 BYTES and the loop's scene block comes first in what it echoes, so
//     one child's whole row is ~150 runes - under the 512-rune retained window. The
//     bound-crossing case therefore reads truncated=true with elided=0, which is true
//     and honest, and the per-row arithmetic lives in internal/panel's
//     TestTruncationFactsRideThePacket instead (real StreamLog, real pump, real bytes).
//  3. Dropped stream keys are never named on a production packet: the hard ceiling is
//     DefaultStreamKeys x 2 = 64 keys, i.e. 64 real child loops. 33 was already the
//     slowest case in this suite.
//  4. The in-flight moment is not caught on the production path. A child's turn is
//     ~2-4 SSE chunks and the root's is ~9, so which one is observable first is a race,
//     and no instrument in this tree makes a child provably slower. Case 1 therefore
//     asserts the row was on the packet the run published (which is the claim ticket
//     197 asks for) and the settled value afterwards; "a running child reads Thinking"
//     stays internal/tools' Test197RowExistsBeforeFirstChildModelCall.
//  5. Cancellation is not re-measured here. 取消不级联 is internal/tools'
//     Test197CancelIsPerRowAndNeverCascades, which this leg neither rewrote nor worked
//     around, and the carrier adds no cancel path of its own.
//  6. No click. Every reading stops at the packet bytes: the transport is still ticket
//     33's, so nothing here proves a page rendered a row, only that a page could read
//     one. panel_pump.go names that last mile and this leg did not invent a substitute.
//  7. The root's own status dimension is asserted as the honest empty, because nothing
//     files a D43 name for a root row yet (ruling A394 leaves that with ticket 196). If
//     a later leg fills it, case 1's root assertions are the ones that must change -
//     they will not quietly start passing.
