package speech

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"
)

// This file is the receipt for ticket 27's punctuation stage. It is split in
// two halves on purpose:
//
//   - the shape half (config pin, roster pin, constructor refusals, toggle) runs
//     everywhere, including a host with a shadowed onnxruntime, because none of
//     it needs the model;
//   - the model half needs the real CT-Transformer and a staged onnxruntime, and
//     it is where the actual readings come from. A missing model is a named
//     skip, never a silent pass.

// puncModelFile is the file name the manifest pins inside the archive.
const puncModelFile = "model.int8.onnx"

// puncModelPath locates the CT-Transformer model for tests.
//
// It honours WISP_PUNC_MODEL, the same staging seam WISP_SILERO_MODEL and
// WISP_SPIKE_MODELS are, because the staging directory holds the DLLs and not
// the models. Otherwise it looks in third_party/spike-models, where the other
// pinned spike models live.
//
// It deliberately does NOT fall back to third_party/model-fetch: the .onnx stubs
// in there are 29-byte failed-download payloads, and a helper that hands one to
// sherpa would turn "the model was never fetched" into a confusing load error.
// The intact archive next to those stubs is the real source (its contents match
// build/models/manifest.json byte for byte); extract it and point WISP_PUNC_MODEL
// at model.int8.onnx.
func puncModelPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_PUNC_MODEL"); p != "" {
		return p
	}
	rel := filepath.Join("third_party", "spike-models", puncModelFile)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, rel)); err == nil {
			return filepath.Join(dir, rel)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skipf("punctuation model not available: set WISP_PUNC_MODEL to the CT-Transformer "+
		"model.int8.onnx (extract third_party/model-fetch/punc-int8.tar.bz2; its %s is "+
		"sha256 65a3fb9f5ad7bfb96bf69e0dc4481df97f6ee60513c1d94ce981ba6effd524b1), or "+
		"place it at %s - looked up from %s", puncModelFile, rel, dir)
	return ""
}

// newTestPunctuator opens the stage, or skips by name when the host cannot run
// it (a shadowed onnxruntime, or no model). A skip is not a pass: the staging
// recipe in docs/evidence/s1/27-r1-impl.md is what makes these run.
func newTestPunctuator(t *testing.T) *CTTransformerPunctuator {
	t.Helper()
	p, err := NewCTTransformerPunctuator(PunctuationConfig{ModelPath: puncModelPath(t)})
	if err != nil {
		t.Skipf("punctuation stage unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })
	return p
}

// punctStubAsr is an AsrEngine that spells out a fixed text, one rune per feed.
// It shares no code with StreamingAsr on purpose: what the seam test proves is
// that a final transcript reaches the punctuation stage whatever produced it,
// which is C9's interchangeability property applied to the stage above it.
type punctStubAsr struct {
	open bool
	fed  int
	text []rune
	want []rune
}

func (c *punctStubAsr) Begin() { c.open, c.fed, c.text = true, 0, nil }

func (c *punctStubAsr) Feed(samples []float32) {
	if !c.open {
		return
	}
	c.fed += len(samples)
	if n := len(c.text); n < len(c.want) {
		c.text = append(c.text, c.want[n])
	}
}

func (c *punctStubAsr) Partial() string {
	if !c.open {
		return ""
	}
	return string(c.text)
}

func (c *punctStubAsr) Finish() string {
	if !c.open {
		return ""
	}
	c.open = false
	return string(c.text)
}

func (c *punctStubAsr) InFlight() bool { return c.open }
func (c *punctStubAsr) Fed() int       { return c.fed }
func (c *punctStubAsr) Close() error   { c.open = false; return nil }

var _ AsrEngine = (*punctStubAsr)(nil)

// emptyPunct is the worst-case stage: it returns nothing for everything. It
// exists to prove the driver's guard, not to stand in for a model.
type emptyPunct struct{ calls int }

func (e *emptyPunct) AddPunct(string) string { e.calls++; return "" }
func (e *emptyPunct) Close() error           { return nil }

var _ Punctuator = (*emptyPunct)(nil)

// marksIn returns the punctuation marks of s, in order, and their kinds as a
// sorted set. Deduplicated by rune so a repeated comma counts once in the set.
func marksIn(s string) (marks []rune, kinds []string) {
	seen := map[rune]bool{}
	for _, r := range s {
		if !unicode.IsPunct(r) {
			continue
		}
		marks = append(marks, r)
		if !seen[r] {
			seen[r] = true
			kinds = append(kinds, string(r))
		}
	}
	sort.Strings(kinds)
	return marks, kinds
}

// contentRunes is s with punctuation and whitespace removed: the invariant that
// a punctuation stage inserts marks and never edits the words.
func contentRunes(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// TestPunctuationConfigPinsTheShape is the wiring receipt for the three
// decisions that must not drift on the punctuation stage.
//
// It asserts the literal 1 on purpose: asserting against punctuationThreads
// would make the test agree with any value that constant ever takes, which is
// exactly how a pin stops being a pin.
func TestPunctuationConfigPinsTheShape(t *testing.T) {
	got := punctuationConfig(PunctuationConfig{ModelPath: "punc.onnx"})

	if n := got.Model.NumThreads; n != 1 {
		t.Fatalf("Model.NumThreads = %d, want the literal 1: punctuation inference runs "+
			"inside the single-core voice chain (D32), and the pin is only real while "+
			"nothing can raise it", n)
	}
	if got.Model.Provider != "cpu" {
		t.Fatalf("Model.Provider = %q, want cpu: one provider for the whole chain (D5)",
			got.Model.Provider)
	}
	if got.Model.Debug != 0 {
		t.Fatalf("Model.Debug = %d, want 0: a stage that prints during inference is not "+
			"measurable", got.Model.Debug)
	}
	if got.Model.CtTransformer != "punc.onnx" {
		t.Fatalf("Model.CtTransformer = %q, want %q", got.Model.CtTransformer, "punc.onnx")
	}
}

// TestPunctuatorRosterIsPinned is the punctuation seam's roster nail: two
// methods, no more.
//
// The reason is engines_test.go's reason applied one stage over. Removing a
// method is already fatal, but ADDING one is invisible to Go's structural
// typing, and this stage sits directly on the path that carries the user's
// words to the agent: a method that rewrote the transcript instead of
// punctuating it would compile everywhere and stay green. Ticketing it here
// means such an edit has to be argued first.
func TestPunctuatorRosterIsPinned(t *testing.T) {
	iface := reflect.TypeOf((*Punctuator)(nil)).Elem()
	got := make([]string, 0, iface.NumMethod())
	for i := 0; i < iface.NumMethod(); i++ {
		got = append(got, iface.Method(i).Name)
	}
	want := []string{"AddPunct", "Close"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Punctuator roster = %v, want exactly %v: this seam carries the user's "+
			"words to the agent, so a method that edited them rather than punctuating "+
			"them must not be addable in silence", got, want)
	}
	t.Logf("Punctuator: %d methods, %v", len(got), got)
}

// TestNewCTTransformerPunctuatorRejectsEmptyPath: no model path is an error, not
// an empty stage that silently degrades and hides a broken download.
func TestNewCTTransformerPunctuatorRejectsEmptyPath(t *testing.T) {
	if _, err := NewCTTransformerPunctuator(PunctuationConfig{}); err == nil {
		t.Fatalf("NewCTTransformerPunctuator with no ModelPath returned a stage, want an error")
	}
}

// TestOpenPunctuationSwitchedOffIsANilStage: the toggle-off path hands back a
// nil Punctuator and no error, so the caller above cannot branch on it.
//
// The nil is the contract. A typed nil would be a non-nil interface and would
// defeat the arrangement, so this asserts on the interface value, not on the
// concrete type.
func TestOpenPunctuationSwitchedOffIsANilStage(t *testing.T) {
	var p Punctuator
	p, err := OpenPunctuation(PunctuationConfig{ModelPath: "ignored.onnx"}, false)
	if err != nil {
		t.Fatalf("OpenPunctuation(enabled=false) = %v, want no error", err)
	}
	if p != nil {
		t.Fatalf("OpenPunctuation(enabled=false) = %#v, want a nil Punctuator: the raw "+
			"transcript has to be the switched-off answer", p)
	}
}

// punctuationSamples is the unpunctuated dictation corpus this stage is read
// against: the scenario-3 shapes from the ticket (timers, schedules, questions,
// multi-clause commands). They are unpunctuated on purpose -- that is what the
// ASR hands over -- so every mark below was inserted by the model.
var punctuationSamples = []struct {
	name string
	in   string
}{
	{"timer", "今天下午三点提醒我开会"},
	{"schedule", "帮我把明天上午的日程整理一下发给我"},
	{"question", "这个月的电费是多少钱"},
	{"clause-chain", "先保存文件然后关闭窗口最后重启电脑"},
	{"two-questions", "明天天气怎么样需要带伞吗"},
	{"long-command", "打开收音机然后把音量调小一点再帮我记一笔买牛奶"},
	{"mixed-digits", "周四下午两点半和三点各有一个会"},
	{"statement", "我刚才说的那些内容你都记住了"},
}

// dictation300 is the ticket's AC#3 fixture: 300 characters of continuous
// speech, unpunctuated, in one breath -- the shape a dictation user actually
// produces when they hold the button and talk.
const dictation300 = "明天上午九点我要先去公司开一个小时的周会然后十点半跟产品团队过一下下个季度的需求评审" +
	"下午两点还要去见一个客户大概四点回到办公室记得提醒我把这周的报销单交上去另外帮我查一下周五去上海的" +
	"火车票最好是上午出发的还有晚上七点约了人在楼下咖啡厅谈合作这件事你别忘了提前半小时提醒我" +
	"对了下周三之前还要把上个月的账单整理好发给财务顺便问一下他们两个工作日之内能不能处理完" +
	"如果不行就约个时间当面谈一下具体的细节问题还有我妈让我周末回家吃饭这件事也帮我记上" +
	"别忘了提前跟我说一声我好买东西对了上次说的那个新同事的入职材料我还没准备你提醒我周四之前弄完" +
	"不然又要被催了最后把今天所有的安排按时间顺序列一遍给我看看有没有冲突的地方"

// TestPunctuationOnThe300CharDictationFixture is the AC#3 fixture: one
// 300-character utterance through the stage in a single call.
//
// It asserts the mechanical properties (length, content preserved, marks
// inserted, no leading mark) and logs the whole pair. The subjective 8/10 score
// the ticket asks for is a judgement about the logged output and belongs in
// docs/evidence/s1/27-r1-impl.md, not in an assertion -- a test that scored its
// own output would be marking its own homework.
//
// A single call also answers a question the short samples cannot: whether the
// CT-Transformer handles a whole dictation run at once or truncates it. The
// content-preserved assertion is what detects truncation, and it would fail
// loudly rather than quietly shortening the transcript.
func TestPunctuationOnThe300CharDictationFixture(t *testing.T) {
	p := newTestPunctuator(t)

	n := len([]rune(dictation300))
	if n < 300 {
		t.Fatalf("the AC#3 fixture is %d characters, want at least 300", n)
	}
	out := p.AddPunct(dictation300)
	marks, kinds := marksIn(out)
	t.Logf("AC#3 fixture: %d characters in, %d marks out, kinds=%v",
		n, len(marks), kinds)
	t.Logf("AC#3 input:  %q", dictation300)
	t.Logf("AC#3 output: %q", out)

	if got, want := contentRunes(out), contentRunes(dictation300); got != want {
		t.Fatalf("the 300-character run did not survive intact (%d content runes out, %d in): "+
			"the model truncated or rewrote it, so a dictation-length utterance needs "+
			"chunking before this stage can be wired", len([]rune(got)), len([]rune(want)))
	}
	if len(marks) < 3 {
		t.Fatalf("a 300-character run came back with %d mark(s) (%q): that is not punctuation "+
			"restoration", len(marks), string(marks))
	}
	if out != "" && unicode.IsPunct([]rune(out)[0]) {
		t.Fatalf("the output starts with punctuation %q", string([]rune(out)[0]))
	}
}

// TestPunctuationIsRestoredPerExample is the core reading of ticket 27: an
// unpunctuated Chinese sentence goes in, and marks the model chose come out.
//
// Three properties are asserted, and the third is what keeps this from being a
// test of "the function did not return an error":
//
//   - marks were inserted: the output carries punctuation the input did not;
//   - the words are untouched: stripping punctuation from both sides leaves the
//     same text, so the stage punctuated rather than rewrote;
//   - the marks are not one constant mark everywhere: at least two distinct
//     kinds appear across the corpus, which a stage that appended a full stop to
//     every input could not fake.
func TestPunctuationIsRestoredPerExample(t *testing.T) {
	p := newTestPunctuator(t)

	kindsAcross := map[string]bool{}
	for _, tc := range punctuationSamples {
		out := p.AddPunct(tc.in)
		added, kinds := marksIn(out)
		for _, k := range kinds {
			kindsAcross[k] = true
		}
		t.Logf("case %-12s in=%q out=%q inserted=%q kinds=%v",
			tc.name, tc.in, out, string(added), kinds)

		if len(added) == 0 {
			t.Errorf("case %s: %q came back with no punctuation at all -- the stage ran "+
				"and did nothing", tc.name, tc.in)
		}
		if got, want := contentRunes(out), contentRunes(tc.in); got != want {
			t.Errorf("case %s: the words changed, not just the punctuation: got %q want "+
				"%q (in %q, out %q)", tc.name, got, want, tc.in, out)
		}
		if out != "" && unicode.IsPunct([]rune(out)[0]) {
			t.Errorf("case %s: the output starts with punctuation %q, which is never a "+
				"correct restoration of %q", tc.name, string([]rune(out)[0]), tc.in)
		}
	}
	if len(kindsAcross) < 2 {
		t.Fatalf("across %d samples the stage only ever inserted %v: that is a constant, "+
			"not a restore", len(punctuationSamples), kindsAcross)
	}
	t.Logf("distinct marks inserted across the corpus: %d", len(kindsAcross))
}

// TestPunctuatorPositiveControls are the three inputs that must not become an
// excuse: nothing, pure punctuation, and pure English.
//
// The last one is the honest one. The shipped model's label set is Chinese
// punctuation (config.yaml punc_list), so it is not a fair expectation that it
// adds English marks; what it must not do is crash, swallow, or rewrite the
// text. Whatever it does add is logged rather than asserted away.
func TestPunctuatorPositiveControls(t *testing.T) {
	p := newTestPunctuator(t)

	for _, tc := range []struct{ name, in string }{
		{"empty", ""},
		{"spaces", "   "},
		{"pure-marks", "\u3002\u3002\uFF0C\uFF1F"},
		{"pure-english", "set a timer for three pm tomorrow"},
	} {
		out := p.AddPunct(tc.in)
		added, _ := marksIn(out)
		t.Logf("positive control %-12s in=%q out=%q inserted=%q", tc.name, tc.in, out, string(added))

		if got, want := contentRunes(out), contentRunes(tc.in); got != want {
			t.Errorf("control %s: content changed: got %q want %q (in %q, out %q)",
				tc.name, got, want, tc.in, out)
		}
		switch tc.name {
		case "empty", "spaces":
			if out != tc.in {
				t.Errorf("control %s: %q became %q; there is nothing there to punctuate, "+
					"so any change is a hallucination", tc.name, tc.in, out)
			}
		}
	}
}

// TestPunctuatorLatencyWithinBudget measures the stage against D32's 300ms and
// reports P50 and P95 over the corpus, repeated.
//
// The first call is a warm-up and is excluded: the model is loaded by the
// constructor, but onnxruntime allocates its inference arena on first Run, and
// billing that to the steady state would misreport the stage the agent waits
// on. The loop runs a fixed number of times -- there is no clock in the
// decision, only in the measurement (D42#9 forbids wall-clock timeouts, not
// wall-clock readings).
func TestPunctuatorLatencyWithinBudget(t *testing.T) {
	p := newTestPunctuator(t)

	const repeats = 10
	_ = p.AddPunct(punctuationSamples[0].in) // warm-up, not measured

	var all []time.Duration
	for i := 0; i < repeats; i++ {
		for _, tc := range punctuationSamples {
			start := time.Now()
			out := p.AddPunct(tc.in)
			d := time.Since(start)
			if out == "" {
				t.Fatalf("case %s returned nothing during the latency run", tc.name)
			}
			all = append(all, d)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	p50 := all[len(all)*50/100]
	p95 := all[len(all)*95/100]
	t.Logf("latency over %d calls (warm-up excluded): P50=%v P95=%v max=%v budget=%v",
		len(all), p50.Round(time.Microsecond), p95.Round(time.Microsecond),
		all[len(all)-1].Round(time.Microsecond), PunctuationLatencyBudget)

	if p95 > PunctuationLatencyBudget {
		t.Fatalf("P95 %v exceeds the %v the stage is held to (D32 16.3.4): the punctuation "+
			"step would blow the dictation budget on its own", p95, PunctuationLatencyBudget)
	}
}

// TestFinalTranscriptCrossesTheSeamPunctuated is the end-to-end reading inside
// this package: the real endpointer decides the utterance, a non-sherpa engine
// produces the words, and the real punctuation stage marks them.
//
// This is the furthest point the assertion can reach without touching cmd/wisp:
// the agent-side sink (voiceTranscriptSink -> residentTaskSource) belongs to
// 15-e and is outside this ticket's write face, so what is proven here is that
// the transcript arriving AT that sink is punctuated and says so.
func TestFinalTranscriptCrossesTheSeamPunctuated(t *testing.T) {
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })

	const want = "今天下午三点提醒我开会"
	eng := &punctStubAsr{want: []rune(want)}
	tr := NewTranscriber(det, eng)
	tr.SetPunctuator(newTestPunctuator(t))

	// The 15-b calibrated fixture: 1.2s of speech, then the 1500ms nominal
	// silence that actually closes an utterance.
	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	if len(out) == 0 {
		t.Fatalf("the seam produced no transcripts at all")
	}
	last := out[len(out)-1]
	if last.Kind != TranscriptFinal {
		t.Fatalf("last transcript is a %v, want the final", last.Kind)
	}
	marks, kinds := marksIn(last.Text)
	t.Logf("final reaching the agent sink: %q (punctuated=%v, inserted=%q, kinds=%v)",
		last.Text, last.Punctuated, string(marks), kinds)

	if !last.Punctuated {
		t.Fatalf("the final reports Punctuated=false with a stage attached: downstream "+
			"cannot tell dictation-grade text from a raw hypothesis (final %q)", last.Text)
	}
	if len(marks) == 0 {
		t.Fatalf("the final %q reached the seam with no punctuation, so the stage is not "+
			"on the path", last.Text)
	}
	if got, wantContent := contentRunes(last.Text), contentRunes(want); got != wantContent {
		t.Fatalf("the words changed on the way through: got %q want %q", got, wantContent)
	}

	// Partials must stay raw: they are display-only (ticket 27).
	for _, tr := range out {
		if tr.Kind != TranscriptPartial {
			continue
		}
		if tr.Punctuated {
			t.Fatalf("a partial reports Punctuated=true: partials are display-only and "+
				"must stay unpunctuated (text %q)", tr.Text)
		}
		if m, _ := marksIn(tr.Text); len(m) > 0 {
			t.Fatalf("a partial carries punctuation %q (text %q): the stage is running "+
				"per frame instead of per utterance", string(m), tr.Text)
		}
	}
}

// TestPunctuationToggleOffPassesTheTranscriptThroughUnchanged is the toggle-off
// path: the same seam, the same engine, no stage, and the raw transcript is
// what arrives -- flagged as raw.
func TestPunctuationToggleOffPassesTheTranscriptThroughUnchanged(t *testing.T) {
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })

	const want = "今天下午三点提醒我开会"
	tr := NewTranscriber(det, &punctStubAsr{want: []rune(want)})
	tr.SetPunctuator(nil) // the toggle-off path

	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	last := out[len(out)-1]
	t.Logf("toggle off: final %q punctuated=%v", last.Text, last.Punctuated)
	if last.Text != want {
		t.Fatalf("toggle off changed the transcript: got %q want %q", last.Text, want)
	}
	if last.Punctuated {
		t.Fatalf("toggle off reports Punctuated=true, so downstream cannot tell the text is raw")
	}
	if m, _ := marksIn(last.Text); len(m) > 0 {
		t.Fatalf("toggle off still produced punctuation %q", string(m))
	}
}

// TestPunctuationLoadFailureDegradesInsteadOfFailingTheUtterance is the
// degraded-load path, and it is also the receipt for the pre-flight guard.
//
// Why the guard exists at all: handing sherpa a file that is not an ONNX model
// does not return nil and does not return an error -- onnxruntime throws a C++
// exception (0xe06d7363) and the process is gone. That was measured on this host
// with the first draft of this very test, which aborted the whole binary. Go
// cannot recover a C++ throw, so checkOnnxModelFile has to run before sherpa
// ever sees the path. The 29-byte payload below is the file shape
// third_party/model-fetch actually ships.
func TestPunctuationLoadFailureDegradesInsteadOfFailingTheUtterance(t *testing.T) {
	if err := checkOnnxruntime(); err != nil {
		t.Skipf("cannot distinguish a load failure from a shadowed runtime here (%v)", err)
	}

	dir := t.TempDir()
	emptyFile := filepath.Join(dir, "empty.onnx")
	if err := os.WriteFile(emptyFile, nil, 0o600); err != nil {
		t.Fatalf("write empty file: %v", err)
	}
	errorPage := filepath.Join(dir, "error-page.onnx")
	if err := os.WriteFile(errorPage, []byte("Invalid username or password."), 0o600); err != nil {
		t.Fatalf("write error page: %v", err)
	}

	for _, tc := range []struct{ name, path string }{
		{"empty", emptyFile},
		{"error-page", errorPage},
		{"missing", filepath.Join(dir, "never-downloaded.onnx")},
		{"directory", dir},
	} {
		p, err := OpenPunctuation(PunctuationConfig{ModelPath: tc.path}, true)
		if err == nil {
			if p != nil {
				_ = p.Close()
			}
			t.Fatalf("OpenPunctuation accepted the %s case (%s): a load failure that is "+
				"not refused here aborts the process inside the C API", tc.name, tc.path)
		}
		if p != nil {
			t.Fatalf("OpenPunctuation returned a stage alongside error %v (%s)", err, tc.name)
		}
		t.Logf("refused %-10s %s: %v", tc.name, filepath.Base(tc.path), err)
	}

	// The seam above it keeps working: no stage, raw transcript, utterance intact.
	det, derr := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if derr != nil {
		t.Skipf("VAD session unavailable on this host: %v", derr)
	}
	t.Cleanup(func() { _ = det.Close() })

	const want = "今天下午三点提醒我开会"
	p, err := OpenPunctuation(PunctuationConfig{ModelPath: errorPage}, true)
	tr := NewTranscriber(det, &punctStubAsr{want: []rune(want)})
	tr.SetPunctuator(p) // nil: the failed load, handed on as "no stage"

	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	if len(out) == 0 {
		t.Fatalf("a failed punctuation load took the whole utterance with it")
	}
	last := out[len(out)-1]
	t.Logf("degraded final: %q punctuated=%v (load error was: %v)", last.Text, last.Punctuated, err)
	if last.Kind != TranscriptFinal || last.Text != want {
		t.Fatalf("degraded final = %v %q, want the raw transcript %q", last.Kind, last.Text, want)
	}
	if last.Punctuated {
		t.Fatalf("a degraded final reports Punctuated=true")
	}
}

// TestPunctuationNeverErasesAnUtterance is the guard that keeps a broken model
// from manufacturing the one signal the caller treats as meaningful.
//
// An empty final means "nothing sayable" and is what AC#5's retry-then-"did not
// catch that" policy keys on. A stage that returned nothing for a sentence has
// failed, not punctuated, so the raw transcript must survive -- otherwise a
// punctuation bug would read downstream as a microphone problem.
func TestPunctuationNeverErasesAnUtterance(t *testing.T) {
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })

	const want = "今天下午三点提醒我开会"
	eng := &punctStubAsr{want: []rune(want)}
	tr := NewTranscriber(det, eng)
	bad := &emptyPunct{}
	tr.SetPunctuator(bad)

	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	last := out[len(out)-1]
	t.Logf("erasing stage: %d call(s), final %q punctuated=%v", bad.calls, last.Text, last.Punctuated)
	if bad.calls == 0 {
		t.Fatalf("the erasing stage was never called, so this control proves nothing")
	}
	if last.Text != want {
		t.Fatalf("a stage that returned nothing ate the utterance: final %q, want %q",
			last.Text, want)
	}
	if last.Punctuated {
		t.Fatalf("a final kept against a failing stage reports Punctuated=true")
	}
}

// TestPunctuationOnTheRealAsrFinal is the honest reading: the real streaming
// recognizer over the 15-b synthetic fixture, punctuated by the real model.
//
// The surrogate has no words, so the recognizer emits a run of one vowel and
// the punctuator has almost nothing to work with. That is logged, not asserted
// around -- the per-example readings above are what show the stage restoring
// punctuation, and this one shows what the stage does when the text it is
// handed has no sentence structure at all.
func TestPunctuationOnTheRealAsrFinal(t *testing.T) {
	asr := newTestAsr(t)
	p := newTestPunctuator(t)

	asr.Begin()
	pcm := synthSpeech(msSamples(1200))
	for i := 0; i+WindowSamples <= len(pcm); i += WindowSamples {
		asr.Feed(pcmToFloat(pcm[i : i+WindowSamples]))
	}
	raw := asr.Finish()
	out := p.AddPunct(raw)
	marks, kinds := marksIn(out)
	t.Logf("real ASR final %q -> punctuated %q (inserted=%q, kinds=%v)",
		raw, out, string(marks), kinds)

	if out == "" && raw != "" {
		t.Fatalf("the stage ate a real ASR final: %q became empty", raw)
	}
	if got, want := contentRunes(out), contentRunes(raw); got != want {
		t.Fatalf("the stage rewrote real ASR output: got %q want %q", got, want)
	}
}
