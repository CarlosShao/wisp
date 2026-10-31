package present

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
)

// The channel behavior per tier, on REAL ports where the repo has them (the
// memory artifacts store) and recording fakes where the host owns the real
// surface (TTS queue, toast, ball, panel, clipboard).

type fakeSpeaker struct {
	mu    sync.Mutex
	lines []string
	taken bool
}

func (f *fakeSpeaker) Speak(text string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lines = append(f.lines, text)
	return f.taken
}

func (f *fakeSpeaker) spoken() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.lines...)
}

type fakeNotifier struct {
	mu     sync.Mutex
	calls  int
	bodies []string
	err    error
}

func (f *fakeNotifier) Notify(title, body string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.bodies = append(f.bodies, body)
	return f.err
}

type fakeBadger struct {
	mu    sync.Mutex
	texts []string
}

func (f *fakeBadger) Badge(text string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.texts = append(f.texts, text)
}

type fakeTaskLister struct {
	mu    sync.Mutex
	texts []string
}

func (f *fakeTaskLister) TaskEntry(text string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.texts = append(f.texts, text)
}

type fakePanel struct {
	mu    sync.Mutex
	texts []string
	arts  []string
}

func (f *fakePanel) ShowResult(text, artifact string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.texts = append(f.texts, text)
	f.arts = append(f.arts, artifact)
}

type fakeClipboard struct {
	mu    sync.Mutex
	calls int
	texts []string
	err   error
}

func (f *fakeClipboard) WriteClipboard(text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.texts = append(f.texts, text)
	return f.err
}

// quotaArtifacts is the artifacts face with an injected failure, so a quota
// refusal on the write path can be observed without a 500MB fixture: the
// real quota enforcement is the store's (ticket 04, on the write path), and
// what THIS layer must prove is that the refusal is named, not swallowed.
type quotaArtifacts struct {
	dir string
	err error
}

func (f quotaArtifacts) PutArtifact(context.Context, string, []byte) error { return f.err }
func (f quotaArtifacts) ArtifactsDir() string                              { return f.dir }

func testChannels() (c Channels, sp *fakeSpeaker, nf *fakeNotifier, bg *fakeBadger,
	tl *fakeTaskLister, pn *fakePanel, cl *fakeClipboard,
) {
	sp = &fakeSpeaker{taken: true}
	nf = &fakeNotifier{}
	bg = &fakeBadger{}
	tl = &fakeTaskLister{}
	pn = &fakePanel{}
	cl = &fakeClipboard{}
	c = Channels{
		Speak: sp, Notify: nf, Badge: bg, TaskList: tl, Panel: pn,
		Clipboard: cl,
	}
	return c, sp, nf, bg, tl, pn, cl
}

func openStore(t *testing.T) *memory.Store {
	t.Helper()
	store, err := memory.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

const tier1Text = "结果如下：\n| 名称 | 数量 |\n|---|---|\n| 甲 | 1 |\n口播摘要：表格已整理完毕"

func TestTierOneWritesArtifactShowsPanelSpeaksSummaryWithPathAndNeverClips(t *testing.T) {
	store := openStore(t)
	c, sp, _, _, _, pn, cl := testChannels()
	c.Artifacts = store
	var buf bytes.Buffer
	c.Logger = slog.New(slog.NewTextHandler(&buf, nil))
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "task-abc-1", Text: tier1Text})

	if rc.Tier != TierStructured {
		t.Fatalf("tier = %s, want structured", rc.Tier)
	}
	if rc.Artifact == "" {
		t.Fatal("receipt carries no artifact path")
	}
	if !strings.HasPrefix(rc.Artifact, store.ArtifactsDir()+string(filepath.Separator)) {
		t.Errorf("artifact path %q is not inside the store's artifacts dir %q", rc.Artifact, store.ArtifactsDir())
	}
	data, err := os.ReadFile(rc.Artifact)
	if err != nil {
		t.Fatalf("artifact file not written: %v", err)
	}
	if string(data) != rc.Display {
		t.Errorf("artifact body = %q, want the display text", data)
	}
	if !strings.Contains(buf.String(), rc.Artifact) {
		t.Errorf("log does not carry the artifact path: %s", buf.String())
	}
	if len(pn.texts) != 1 || pn.texts[0] != rc.Display || pn.arts[0] != rc.Artifact {
		t.Errorf("panel got (%q, %q), want display plus artifact path", pn.texts, pn.arts)
	}
	wantSpoken := "表格已整理完毕" + spokenHaveFile + rc.Artifact
	if got := sp.spoken(); len(got) != 1 || got[0] != wantSpoken {
		t.Errorf("spoken = %q, want %q", got, wantSpoken)
	}
	if cl.calls != 0 {
		t.Errorf("clipboard called %d times on tier 1, want 0: long results never take the clipboard", cl.calls)
	}
	if !rc.OK {
		t.Errorf("receipt not OK: %+v", rc.Notes)
	}
}

func TestTierTwoFallsBackToTheFirstSentenceWhenNoMarkerCame(t *testing.T) {
	store := openStore(t)
	c, sp, _, _, _, pn, _ := testChannels()
	c.Artifacts = store
	text := "首句在这里。 " + strings.Repeat("测", 401)
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "t2", Text: text})

	if rc.Tier != TierLong {
		t.Fatalf("tier = %s, want long", rc.Tier)
	}
	if rc.Artifact == "" {
		t.Fatal("tier 2 must land a file")
	}
	if len(pn.texts) != 1 {
		t.Fatalf("panel calls = %d, want 1", len(pn.texts))
	}
	if got := sp.spoken(); len(got) != 1 || !strings.HasPrefix(got[0], "首句在这里") {
		t.Fatalf("spoken = %q, want the first-sentence fallback", got)
	}
	if !strings.Contains(sp.spoken()[0], rc.Artifact) {
		t.Errorf("spoken = %q, want the artifact path announced", sp.spoken()[0])
	}
}

func TestTierThreeDualChannelSurvivesAToastFailure(t *testing.T) {
	c, sp, nf, bg, tl, pn, cl := testChannels()
	nf.err = errors.New("toast swallowed by focus assist")
	text := strings.Repeat("测", 61) + "。第二句"
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "t3", Text: text})

	if rc.Tier != TierMedium {
		t.Fatalf("tier = %s, want medium", rc.Tier)
	}
	if nf.calls != 1 {
		t.Errorf("notify calls = %d, want 1", nf.calls)
	}
	// D42#8: the badge and the task list update even though the toast died.
	if len(bg.texts) != 1 {
		t.Errorf("badge updates = %d, want 1 (dual-channel is unconditional)", len(bg.texts))
	}
	if len(tl.texts) != 1 {
		t.Errorf("task list entries = %d, want 1 (dual-channel is unconditional)", len(tl.texts))
	}
	if len(pn.texts) != 1 || pn.arts[0] != "" {
		t.Errorf("panel = %q with artifact %q, want the text with no artifact", pn.texts, pn.arts)
	}
	if cl.calls != 0 {
		t.Errorf("clipboard called on tier 3, want 0")
	}
	spoken := sp.spoken()
	if len(spoken) != 1 {
		t.Fatalf("spoken = %q, want the first sentence", spoken)
	}
	if runes := len([]rune(spoken[0])); runes > MaxSpokenRunes {
		t.Errorf("spoken = %d runes, over the ceiling", runes)
	}
	if rc.OK {
		t.Errorf("receipt OK with a failed toast: %+v", rc.Notes)
	}
	var notified, badged bool
	for _, n := range rc.Notes {
		if n.Channel == "notify" && n.Status == "failed" {
			notified = true
		}
		if n.Channel == "badge" && n.Status == "sent" {
			badged = true
		}
	}
	if !notified || !badged {
		t.Errorf("trace = %+v, want a failed notify beside a sent badge", rc.Notes)
	}
}

func TestTierFourClipsBoardAndSpeaksAndPullsNothingElse(t *testing.T) {
	c, sp, nf, bg, tl, pn, cl := testChannels()
	text := "今天天气不错"
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "t4", Text: text})

	if rc.Tier != TierShort {
		t.Fatalf("tier = %s, want short", rc.Tier)
	}
	if cl.calls != 1 || cl.texts[0] != text {
		t.Errorf("clipboard = %d calls with %q, want exactly the display text", cl.calls, cl.texts)
	}
	if got := sp.spoken(); len(got) != 1 || got[0] != text {
		t.Errorf("spoken = %q, want the whole short result", got)
	}
	if nf.calls != 0 || len(bg.texts) != 0 || len(tl.texts) != 0 || len(pn.texts) != 0 {
		t.Errorf("tier 4 pulled another channel: notify=%d badge=%d task=%d panel=%d, want all zero",
			nf.calls, len(bg.texts), len(tl.texts), len(pn.texts))
	}
	if !rc.OK {
		t.Errorf("receipt not OK: %+v", rc.Notes)
	}
}

func TestEmptyResultRoutesNowhereAndPanicsNowhere(t *testing.T) {
	c, sp, nf, bg, tl, pn, cl := testChannels()
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "t0", Text: "   "})
	if rc.Tier != TierNone {
		t.Fatalf("tier = %s, want none", rc.Tier)
	}
	if len(sp.spoken()) != 0 || nf.calls != 0 || len(bg.texts) != 0 ||
		len(tl.texts) != 0 || len(pn.texts) != 0 || cl.calls != 0 {
		t.Errorf("an empty result reached a channel: %+v", rc.Notes)
	}
	if !rc.OK {
		t.Errorf("receipt not OK: %+v", rc.Notes)
	}
}

func TestUnwiredChannelsAreSkippedNotFatal(t *testing.T) {
	// No ports at all: the named behavior is per-channel skip notes and a
	// clean receipt, never a panic.
	rc := New(Channels{}).Dispatch(context.Background(), Input{TaskID: "tn", Text: tier1Text})
	if rc.Tier != TierStructured {
		t.Fatalf("tier = %s, want structured", rc.Tier)
	}
	if rc.Artifact != "" {
		t.Errorf("artifact = %q, want empty without an artifacts face", rc.Artifact)
	}
	for _, want := range []string{"panel", "speak", "artifact"} {
		found := false
		for _, n := range rc.Notes {
			if n.Channel == want && n.Status == "skipped" && n.Detail == "not wired" {
				found = true
			}
		}
		if !found {
			t.Errorf("channel %q has no skipped note: %+v", want, rc.Notes)
		}
	}
	if !rc.OK {
		t.Errorf("skips must not spoil the receipt: %+v", rc.Notes)
	}
}

func TestRefusedAnnounceIsReportedNotSwallowed(t *testing.T) {
	c, sp, _, _, _, _, _ := testChannels()
	sp.taken = false
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "tr", Text: "短结果"})
	if rc.Tier != TierShort {
		t.Fatalf("tier = %s, want short", rc.Tier)
	}
	var failed bool
	for _, n := range rc.Notes {
		if n.Channel == "speak" && n.Status == "failed" {
			failed = true
		}
	}
	if !failed {
		t.Errorf("trace = %+v, want a failed speak note", rc.Notes)
	}
	if rc.OK {
		t.Errorf("receipt OK with a refused announce")
	}
}

func TestArtifactFailureIsNamedAndPresentationSurvives(t *testing.T) {
	c, sp, _, _, _, pn, _ := testChannels()
	c.Artifacts = quotaArtifacts{
		dir: t.TempDir(),
		err: errors.New("memory: enforce artifacts quota after \"result-x.md\": over quota"),
	}
	rc := New(c).Dispatch(context.Background(), Input{TaskID: "tq", Text: tier1Text})

	if rc.Artifact != "" {
		t.Errorf("artifact = %q, want empty after a failed write", rc.Artifact)
	}
	if len(pn.texts) != 1 {
		t.Fatalf("panel calls = %d, want 1: the result is still shown", len(pn.texts))
	}
	got := sp.spoken()
	if len(got) != 1 || !strings.HasSuffix(got[0], spokenNoFile) {
		t.Errorf("spoken = %q, want the honest could-not-save tail", got)
	}
	if rc.OK {
		t.Errorf("receipt OK with a failed artifact write")
	}
}

func TestArtifactNamePassesTheStoreGuardAndIsContentInjective(t *testing.T) {
	n1 := ArtifactName("task-abc-1", "正文甲")
	n2 := ArtifactName("task-abc-1", "正文乙")
	n3 := ArtifactName("task-abc-1", "正文甲")
	if n1 != n3 {
		t.Errorf("same content produced %q and %q, want one name", n1, n3)
	}
	if n1 == n2 {
		t.Errorf("different content produced the same name %q", n1)
	}
	if err := memory.ValidArtifactName(n1); err != nil {
		t.Errorf("name %q refused by the store's own guard: %v", n1, err)
	}
	weird := ArtifactName(`..\evil:task`, "x")
	if err := memory.ValidArtifactName(weird); err != nil {
		t.Errorf("sanitized name %q refused: %v", weird, err)
	}
	if strings.ContainsAny(weird, `\/:`) {
		t.Errorf("sanitized name %q still carries a separator", weird)
	}
}
