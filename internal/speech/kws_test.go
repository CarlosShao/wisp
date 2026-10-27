package speech

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// errKwsTestBoom is the loader failure the slot test keys on.
var errKwsTestBoom = errors.New("kws loader boom")

// fakeKws41 records its Close into the slot test's recorder, the way
// fakeAsr26/fakeTts26 do for the other two families.
type fakeKws41 struct{ rec *recorder26 }

func (f *fakeKws41) Accept(samples []float32, out []KeywordHit) []KeywordHit {
	return out
}

func (f *fakeKws41) Close() error {
	f.rec.record("close:kws")
	return nil
}

// Ticket 41's receipts. The engine reads come from the staging recipe (the
// sherpa DLLs next to the test binary plus WISP_SPIKE_MODELS): a construction
// failure there is a red, never a skip -- this leg does not trade skips or
// loosened assertions for green. The model fixtures are the checkpoint's own
// test_wavs (0.wav..6.wav); no audio is synthesized for the positive hits.

// kwsFixtureDir returns the test_wavs directory of the pinned model.
func kwsFixtureDir(t *testing.T) string {
	t.Helper()
	dir, err := kwsModelDir()
	if err != nil {
		t.Fatalf("KWS model directory: %v", err)
	}
	return filepath.Join(dir, "test_wavs")
}

// kwsTableWords reads the display words of a keyword table: the text after
// the '@' on every non-empty line. The spotter reports keywords spelled the
// way the table spells them, so every hit the engine can ever produce must be
// one of these.
func kwsTableWords(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read keyword table %s: %v", path, err)
	}
	var words []string
	for _, line := range strings.Split(string(data), "\n") {
		if i := strings.IndexByte(line, '@'); i >= 0 {
			if w := strings.TrimSpace(line[i+1:]); w != "" {
				words = append(words, w)
			}
		}
	}
	if len(words) == 0 {
		t.Fatalf("keyword table %s parsed to zero display words", path)
	}
	return words
}

// kwsTableWordsIs shows whether w is one of the table's display words.
func kwsTableWordsIs(words []string, w string) bool {
	for _, want := range words {
		if w == want {
			return true
		}
	}
	return false
}

// kwsWavSamples loads a fixture wav and returns its samples at the seam scale.
// The fixtures are 16k/mono/PCM16 RIFF files with a LIST chunk before "data"
// (the sherpa release train wrote them), so the reader walks chunks until
// "data" and fails loud on anything but 16k/mono/PCM16. It deliberately does
// not drive audio.WavInjector: that source is a pump goroutine whose
// start/cancel pair is racy for a back-to-back drain (Stop can cancel the pump
// before its first push), and a deterministic reader needs no goroutine, no
// channel and no wall clock. The real C8 parser seam stays covered by audio's
// own tests and by the Consume path.
func kwsWavSamples(t *testing.T, path string) []float32 {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture wav: %v", err)
	}
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		t.Fatalf("%s is not a RIFF/WAVE file", path)
	}
	var (
		body []byte
		rate int
	)
	off := 12
	for off+8 <= len(data) {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		start := off + 8
		if start+size > len(data) {
			t.Fatalf("%s: chunk %q overruns the file", path, id)
		}
		switch id {
		case "fmt ":
			if size < 16 {
				t.Fatalf("%s: fmt chunk too short (%d)", path, size)
			}
			tag := binary.LittleEndian.Uint16(data[start : start+2])
			chans := binary.LittleEndian.Uint16(data[start+2 : start+4])
			rate = int(binary.LittleEndian.Uint32(data[start+4 : start+8]))
			bits := binary.LittleEndian.Uint16(data[start+14 : start+16])
			if tag != 1 || bits != 16 || chans != 1 || rate != audio.TargetRate {
				t.Fatalf("%s: want 16k/mono/PCM16, got tag=%d bits=%d chans=%d rate=%d", path, tag, bits, chans, rate)
			}
		case "data":
			body = data[start : start+size]
		}
		off = start + size + (size & 1) // chunks are word-aligned
	}
	if len(body) == 0 {
		t.Fatalf("%s: no data chunk", path)
	}
	n := len(body) / 2
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = float32(int16(binary.LittleEndian.Uint16(body[2*i:]))) / 32768
	}
	return out
}

// feedKws pushes samples through the engine in production-shaped windows and
// returns every hit. Returns the hit count and the engine's total fed count.
func feedKws(k *KwsEngine, samples []float32, out []KeywordHit) ([]KeywordHit, int) {
	for i := 0; i < len(samples); i += WindowSamples {
		end := i + WindowSamples
		if end > len(samples) {
			end = len(samples)
		}
		out = k.Accept(samples[i:end], out)
	}
	return out, k.fed
}

// newTestKws opens the engine, or fails the test: on the staging host every
// ingredient is present, so a failure is a finding, not an environment note.
func newTestKws(t *testing.T) *KwsEngine {
	t.Helper()
	cfg, err := DefaultKwsConfig()
	if err != nil {
		t.Fatalf("DefaultKwsConfig: %v", err)
	}
	k, err := NewKwsEngine(cfg)
	if err != nil {
		t.Fatalf("NewKwsEngine: %v", err)
	}
	t.Cleanup(func() { _ = k.Close() })
	return k
}

// newTestKwsWithTable opens the engine pointed at a specific keyword table
// (the demo pairing test below needs the fixture's own demo table).
func newTestKwsWithTable(t *testing.T, keywordsPath string) *KwsEngine {
	t.Helper()
	cfg, err := DefaultKwsConfig()
	if err != nil {
		t.Fatalf("DefaultKwsConfig: %v", err)
	}
	cfg.KeywordsPath = keywordsPath
	k, err := NewKwsEngine(cfg)
	if err != nil {
		t.Fatalf("NewKwsEngine: %v", err)
	}
	t.Cleanup(func() { _ = k.Close() })
	return k
}

// TestKwsSpotterConfigPinsTheShape is the wiring receipt, asserted without a
// model file: one intra-op thread, the feature shape the model is trained on,
// the transducer triple wired where sherpa's KWS config keeps it, and the
// three search knobs exactly at the package constants. The literals are
// asserted on purpose -- asserting against the constants would let any future
// value pass (15-b's shape for this receipt).
func TestKwsSpotterConfigPinsTheShape(t *testing.T) {
	cfg := KwsConfig{
		EncoderPath:  "enc.onnx",
		DecoderPath:  "dec.onnx",
		JoinerPath:   "join.onnx",
		TokensPath:   "tokens.txt",
		KeywordsPath: "keywords.txt",
	}
	got := kwsSpotterConfig(cfg)

	if n := got.ModelConfig.NumThreads; n != 1 {
		t.Fatalf("ModelConfig.NumThreads = %d, want the literal 1: KWS inference must be "+
			"single-core pinned like every engine in this package (AC#6, D32)", n)
	}
	if sr := got.FeatConfig.SampleRate; sr != 16000 {
		t.Fatalf("FeatConfig.SampleRate = %d, want 16000", sr)
	}
	if fd := got.FeatConfig.FeatureDim; fd != 80 {
		t.Fatalf("FeatConfig.FeatureDim = %d, want 80 (the KWS zipformer's mel width)", fd)
	}
	tr := got.ModelConfig.Transducer
	if tr.Encoder != cfg.EncoderPath || tr.Decoder != cfg.DecoderPath || tr.Joiner != cfg.JoinerPath {
		t.Fatalf("transducer triple = (%q, %q, %q), want (%q, %q, %q)", tr.Encoder, tr.Decoder, tr.Joiner,
			cfg.EncoderPath, cfg.DecoderPath, cfg.JoinerPath)
	}
	if tok := got.ModelConfig.Tokens; tok != cfg.TokensPath {
		t.Fatalf("tokens = %q, want %q", tok, cfg.TokensPath)
	}
	if kf := got.KeywordsFile; kf != cfg.KeywordsPath {
		t.Fatalf("keywords file = %q, want %q", kf, cfg.KeywordsPath)
	}
	if p := got.MaxActivePaths; p != DefaultKwsMaxActivePaths {
		t.Fatalf("MaxActivePaths = %d, want %d", p, DefaultKwsMaxActivePaths)
	}
	if s := got.KeywordsScore; s != DefaultKwsScore {
		t.Fatalf("KeywordsScore = %v, want %v", s, DefaultKwsScore)
	}
	if th := got.KeywordsThreshold; th != DefaultKwsThreshold {
		t.Fatalf("KeywordsThreshold = %v with a zero override, want the default %v", th, DefaultKwsThreshold)
	}

	cfg.Threshold = 0.9
	if th := kwsSpotterConfig(cfg).KeywordsThreshold; th != 0.9 {
		t.Fatalf("KeywordsThreshold = %v with override 0.9, want the override", th)
	}
}

// TestDefaultKwsConfigResolution pins the WISP_SPIKE_MODELS contract: both
// spellings (the spike-models root and the KWS directory itself) resolve, and
// a directory without the keyword table is a named error, never a guess.
func TestDefaultKwsConfigResolution(t *testing.T) {
	real, err := kwsModelDir()
	if err != nil {
		t.Fatalf("this host has no KWS model directory at all: %v", err)
	}

	t.Setenv("WISP_SPIKE_MODELS", real)
	if got, err := kwsModelDir(); err != nil || got != real {
		t.Fatalf("WISP_SPIKE_MODELS at the KWS directory resolved to (%q, %v)", got, err)
	}

	root := filepath.Dir(real)
	t.Setenv("WISP_SPIKE_MODELS", root)
	if got, err := kwsModelDir(); err != nil || got != real {
		t.Fatalf("WISP_SPIKE_MODELS at the spike-models root resolved to (%q, %v), want %q", got, err, real)
	}

	t.Setenv("WISP_SPIKE_MODELS", t.TempDir())
	_, err = kwsModelDir()
	if err == nil {
		t.Fatalf("an empty WISP_SPIKE_MODELS resolved; it must be a named error")
	}
	if !strings.Contains(err.Error(), "WISP_SPIKE_MODELS") {
		t.Fatalf("error %v does not name the variable a reader would fix", err)
	}
	t.Setenv("WISP_SPIKE_MODELS", "")
	// The walk-up is proven from the package's own source directory, which
	// runtime.Caller resolves to its path inside the repository no matter where
	// the test binary runs from (under staging the working directory is the
	// staging bin, which has no repository above it).
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller could not resolve this test's source path")
	}
	pkgDir := filepath.Dir(thisFile)
	got, err := kwsWalkUpFrom(pkgDir)
	if err != nil {
		t.Fatalf("walk-up from the package directory (%s) failed: %v", pkgDir, err)
	}
	if got != real {
		t.Fatalf("walk-up resolved %q, want %q", got, real)
	}
}

// TestKwsEngineRosterIsTheC9Roster: the engine's exported method set is
// exactly the frozen [Accept Close] roster -- no more, no less. Removing a
// method already breaks the compile-time assertion in kws.go; ADDING one is
// invisible to structural typing, which is what this pin is for.
func TestKwsEngineRosterIsTheC9Roster(t *testing.T) {
	typ := reflect.TypeOf((*KwsEngine)(nil))
	var got []string
	for i := 0; i < typ.NumMethod(); i++ {
		got = append(got, typ.Method(i).Name)
	}
	want := []string{"Accept", "Close"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("*KwsEngine exported methods = %v, want exactly %v (the C9 roster of "+
			"WakeWordEngine): an added method has to be argued into the seam first", got, want)
	}
}

// TestKwsEngineRejectsIncompleteConfig: a half-specified config must be an
// error, not a session that fails later at the cgo boundary.
func TestKwsEngineRejectsIncompleteConfig(t *testing.T) {
	for _, cfg := range []KwsConfig{
		{},
		{EncoderPath: "enc.onnx"},
		{EncoderPath: "enc.onnx", DecoderPath: "dec.onnx"},
		{EncoderPath: "enc.onnx", DecoderPath: "dec.onnx", JoinerPath: "join.onnx"},
		{EncoderPath: "enc.onnx", DecoderPath: "dec.onnx", JoinerPath: "join.onnx", TokensPath: "tokens.txt"},
	} {
		if _, err := NewKwsEngine(cfg); err == nil {
			t.Fatalf("NewKwsEngine(%+v) returned an engine, want an error", cfg)
		}
	}
}

// TestKwsEngineHitsTheFixtureWavs is the ticket's headline reading. The
// checkpoint's own test_wavs are the demo recordings for the demo table
// (test_wavs/test_keywords.txt) -- that pairing is what the sherpa c-api
// example documents for this checkpoint family, and what this host measured:
// 0.wav..2.wav are the set's official negatives, 3.wav..6.wav carry the eight
// demo keywords (蒋友伯+女儿，周望军+落实，朱丽楠+见面会，文森特卡索). The
// production table (keywords.txt, the eight wake words) scores zero hits on
// all seven fixtures because nobody said those words -- the engine correctly
// rejecting audio that does not contain them, which is its own control test
// below. Every hit's Word must be one of the demo table's own display words,
// and the per-file readings (Word, At) are logged for the verdict table.
func TestKwsEngineHitsTheFixtureWavs(t *testing.T) {
	fixtures := kwsFixtureDir(t)
	demoTable := filepath.Join(fixtures, "test_keywords.txt")
	words := kwsTableWords(t, demoTable)
	k := newTestKwsWithTable(t, demoTable)

	for n := 0; n <= 6; n++ {
		name := filepath.Join(fixtures, string(rune('0'+n))+".wav")
		samples := kwsWavSamples(t, name)
		if len(samples) == 0 {
			t.Fatalf("%s parsed to zero samples", name)
		}
		hits, fed := feedKws(k, samples, nil)
		// The set's own shape, measured and pinned word-for-word: 0.wav..2.wav
		// are the official negatives, 3.wav..6.wav carry the eight demo
		// keywords. The sequences are exact because the search is deterministic
		// (dither=0, one thread, fixed model and audio) and because the
		// sequences are where the sherpa Reset clause shows: without the Reset
		// after every hit, 3.wav stops after its first keyword (the binding's
		// "You MUST call it right after detecting a keyword" is what lets the
		// next keyword in the same stream fire).
		want := []string(nil)
		switch n {
		case 3:
			want = []string{"文森特卡索", "法国"}
		case 4:
			want = []string{"蒋友伯", "女儿"}
		case 5:
			want = []string{"周望军", "落实"}
		case 6:
			want = []string{"朱丽楠", "见面会"}
		}
		var got []string
		for _, h := range hits {
			got = append(got, h.Word)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s hits = %v, want exactly %v over %d samples: a wrong count is a search "+
				"or parameter drift, and a missing second keyword is the sherpa Reset clause broken",
				filepath.Base(name), got, want, fed)
		}
		for _, hit := range hits {
			if !kwsTableWordsIs(words, hit.Word) {
				t.Fatalf("%s hit word %q is not one of the demo table's display words %v", name, hit.Word, words)
			}
			if hit.At <= 0 || hit.At > fed {
				t.Fatalf("%s hit At=%d is not on the audio clock (fed=%d): a hit position must name "+
					"the samples that produced it (D42#9)", name, hit.At, fed)
			}
			if hit.Veto {
				t.Fatalf("%s hit word %q claims Veto: the demo table has no veto word in it", name, hit.Word)
			}
		}
		var spell []string
		for _, h := range hits {
			spell = append(spell, h.Word)
		}
		t.Logf("fixture %s: %d samples -> Word=%v (fed=%d)", filepath.Base(name), len(samples), spell, fed)
	}
}

// TestKwsEngineProductionTableRejectsTheDemoWavs is the control that makes the
// headline reading mean something: the production table's eight wake words
// were never said in the demo recordings, so the engine over the production
// table must score zero hits on all seven fixtures. This is the production
// shape of "a hit requires the word in the audio".
func TestKwsEngineProductionTableRejectsTheDemoWavs(t *testing.T) {
	k := newTestKws(t)
	fixtures := kwsFixtureDir(t)
	for n := 0; n <= 6; n++ {
		name := filepath.Join(fixtures, string(rune('0'+n))+".wav")
		samples := kwsWavSamples(t, name)
		hits, fed := feedKws(k, samples, nil)
		if len(hits) != 0 {
			t.Fatalf("%s produced %d hit(s) (%v) against the production table over %d samples: "+
				"nobody said a production wake word in this recording", name, len(hits), hits, fed)
		}
	}
	t.Logf("production table keywords.txt: 7 fixtures, 0 hits everywhere")
}

// TestKwsEngineOnSilenceStaysSilent is the negative control, recognizer-level:
// three seconds of digital silence must produce zero hits.
func TestKwsEngineOnSilenceStaysSilent(t *testing.T) {
	k := newTestKws(t)
	samples := pcmToFloat(silence(msSamples(3000)))
	hits, fed := feedKws(k, samples, nil)
	if len(hits) != 0 {
		t.Fatalf("3s of digital silence produced %d hit(s) (%v), want 0: the engine is "+
			"firing on nothing", len(hits), hits)
	}
	t.Logf("silence: %d samples -> 0 hits", fed)
}

// TestKwsEngineOnNonKeywordSpeechStaysSilent is the second negative control:
// the package's synthetic speech surrogate -- real-speech-shaped, and
// semantically far from every table word -- must produce zero hits.
func TestKwsEngineOnNonKeywordSpeechStaysSilent(t *testing.T) {
	k := newTestKws(t)
	samples := pcmToFloat(synthSpeech(msSamples(5000)))
	hits, fed := feedKws(k, samples, nil)
	if len(hits) != 0 {
		t.Fatalf("5s of non-keyword speech produced %d hit(s) (%v), want 0: a false trigger "+
			"on generic speech-shaped audio", len(hits), hits)
	}
	t.Logf("non-keyword speech: %d samples -> 0 hits", fed)
}

// TestKwsEngineResetAfterHitLetsTheNextUtteranceHitAgain is the receipt for
// the binding's verbatim clause (sherpa_onnx.go:2466, "You MUST call it right
// after detecting a keyword"): with the engine's own Reset in place, the same
// keyword audio fed a second time hits again (the demo recordings carry
// several keywords each, so a fresh spotter still has work to do after the
// first hit). The reverse criterion for this test is the mutation that deletes
// the Reset call in kws.go: the second pass then goes silent or wrong, and
// this test goes red.
func TestKwsEngineResetAfterHitLetsTheNextUtteranceHitAgain(t *testing.T) {
	fixtures := kwsFixtureDir(t)
	demoTable := filepath.Join(fixtures, "test_keywords.txt")
	k := newTestKwsWithTable(t, demoTable)

	name := filepath.Join(fixtures, "4.wav")
	samples := kwsWavSamples(t, name)

	first, fed1 := feedKws(k, samples, nil)
	if len(first) == 0 {
		t.Fatalf("first pass produced 0 hit(s), want at least one")
	}
	second, fed2 := feedKws(k, samples, nil)
	if len(second) == 0 {
		t.Fatalf("second pass of the same audio produced 0 hit(s) over %d samples: without a "+
			"search reset after the first hit the spotter cannot fire again", fed2)
	}
	if second[0].At <= first[0].At {
		t.Fatalf("second pass At=%d is not after first pass At=%d: the clock must keep running forward", second[0].At, first[0].At)
	}
	var f1, f2 []string
	for _, h := range first {
		f1 = append(f1, h.Word)
	}
	for _, h := range second {
		f2 = append(f2, h.Word)
	}
	t.Logf("two passes over %s: first=%v (at %d of %d) second=%v (at %d of %d)", filepath.Base(name),
		f1, first[0].At, fed1, f2, second[0].At, fed2)
}

// TestKwsProbeSherpaDebug is the raw-sherpa control behind the demo pairing:
// without the Go engine, with the checkpoint's documented demo parameters
// (score 3.0, threshold 0.1, c-api.h example), the spotter fires on the demo
// fixtures. If the engine ever diverges from this, the diff is in kws.go, not
// in sherpa.
func TestKwsProbeSherpaDebug(t *testing.T) {
	cfg, err := DefaultKwsConfig()
	if err != nil {
		t.Fatalf("DefaultKwsConfig: %v", err)
	}
	sc := &sherpa.KeywordSpotterConfig{
		FeatConfig: sherpa.FeatureConfig{SampleRate: 16000, FeatureDim: 80},
		ModelConfig: sherpa.OnlineModelConfig{
			Transducer: sherpa.OnlineTransducerModelConfig{
				Encoder: cfg.EncoderPath, Decoder: cfg.DecoderPath, Joiner: cfg.JoinerPath,
			},
			Tokens: cfg.TokensPath, NumThreads: 1, Provider: "cpu",
		},
		MaxActivePaths:    4,
		KeywordsFile:      filepath.Join(kwsFixtureDir(t), "test_keywords.txt"),
		KeywordsScore:     3.0,
		KeywordsThreshold: 0.1,
	}
	spotter := sherpa.NewKeywordSpotter(sc)
	if spotter == nil {
		t.Fatal("NewKeywordSpotter returned nil")
	}
	defer sherpa.DeleteKeywordSpotter(spotter)
	stream := sherpa.NewKeywordStream(spotter)
	if stream == nil {
		t.Fatal("NewKeywordStream returned nil")
	}
	defer sherpa.DeleteOnlineStream(stream)

	total := 0
	for n := 3; n <= 6; n++ {
		samples := kwsWavSamples(t, filepath.Join(kwsFixtureDir(t), string(rune('0'+n))+".wav"))
		stream.AcceptWaveform(audio.TargetRate, samples)
		for spotter.IsReady(stream) {
			spotter.Decode(stream)
			if r := spotter.GetResult(stream); r.Keyword != "" {
				total++
				t.Logf("raw sherpa hit: %d.wav -> %q", n, r.Keyword)
			}
		}
	}
	if total == 0 {
		t.Fatal("raw sherpa with the demo pairing scored zero hits over 3.wav..6.wav")
	}
	t.Logf("raw sherpa demo pairing: %d hits over 3.wav..6.wav", total)
}

// TestKwsEngineCloseIsIdempotent: C11 -- the second Close is a no-op, and an
// Accept after Close degrades to silence instead of panicking.
func TestKwsEngineCloseIsIdempotent(t *testing.T) {
	k := newTestKws(t)
	if err := k.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := k.Close(); err != nil {
		t.Fatalf("second Close must be idempotent (C11): %v", err)
	}
	if got := k.Accept(pcmToFloat(silence(WindowSamples)), nil); len(got) != 0 {
		t.Fatalf("Accept after Close produced %d hit(s), want a silent no-op", len(got))
	}
}

// Test41TheSlotHoldsTheKwsFamily extends the slot receipt to the third family:
// same rules as ASR/TTS -- release before load, failed load leaves the slot
// empty, redundant load returns the resident engine.
func Test41TheSlotHoldsTheKwsFamily(t *testing.T) {
	rec := &recorder26{}
	slot := NewSlot()

	kws, err := slot.LoadKws(func() (WakeWordEngine, error) {
		rec.record("load:kws")
		return &fakeKws41{rec: rec}, nil
	})
	if err != nil {
		t.Fatalf("LoadKws: %v", err)
	}
	if slot.Resident() != FamilyKws || slot.Kws() != kws {
		t.Fatalf("resident = %q, kws match = %v, want the KWS family", slot.Resident(), slot.Kws() == kws)
	}

	if _, err := slot.LoadTts(func() (TtsEngine, error) {
		rec.record("load:tts")
		return &fakeTts26{name: "tts", rec: rec}, nil
	}); err != nil {
		t.Fatalf("LoadTts: %v", err)
	}
	if got, want := rec.String(), "load:kws,close:kws,load:tts"; got != want {
		t.Fatalf("slot lifecycle = %q, want %q (KWS must be released before TTS loads)", got, want)
	}
	if slot.Kws() != nil {
		t.Fatalf("the KWS family is still reachable while TTS is resident")
	}

	// A failed load after a release leaves the slot empty, not half-built.
	if _, err := slot.LoadKws(func() (WakeWordEngine, error) {
		rec.record("load:kws-fail")
		return nil, errKwsTestBoom
	}); err != errKwsTestBoom {
		t.Fatalf("LoadKws = %v, want the loader's error", err)
	}
	if slot.Resident() != FamilyNone || slot.Kws() != nil {
		t.Fatalf("failed load left residency %q behind", slot.Resident())
	}
	if slot.Loads() != 2 || slot.Closes() != 2 {
		t.Fatalf("counters loads=%d closes=%d, want 2 and 2", slot.Loads(), slot.Closes())
	}

	if _, err := slot.LoadKws(nil); err == nil {
		t.Fatalf("LoadKws(nil) returned nil error")
	}
	if err := slot.Close(); err != nil {
		t.Fatalf("Close on an empty slot: %v", err)
	}
}
