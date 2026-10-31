//go:build windows

package main

// 票 12 AC#9（registry A25）的 3s 回落实测：≥5 个样本、单调钟、逐样本毫秒全报、
// 坏尾部照录不许平均抹掉。判据原文（票面 :150-159）：
//
//	①从「最后一个 token 送达」起算到进入 Settling ≤3000ms，必须用单调钟；
//	②≥5 个样本全报；③机器可判形态优先，接缝＝产码。
//
// 样本走的是产线链，只有三处按本仓注入面纪律替身：
//   - C5 LlmProvider＝mockllm（newRunFixture，run_test.go 同一套）——流式回复的
//     最后一条 delta 在产线 consoleSink 上落 T0；
//   - C9 TTS 引擎＝三方法桩（settle12r2Tts）——让被测的腿只量回落、不量合成，
//     输出设备沿 wire-r1 的 voiceWireSpeaker 桩；
//   - 机器桥＝测试侧按冻结表自己的行走位（#4、#11，26r3/381r2 同款先例）——
//     #15/#26/EvWarmIdle 半边的派发仍是产线 speak()/scope 自己的。
//
// 暖窗经 voicePlan.WarmTimeout（= [session].warm_timeout_sec 的读数位）设 1.2s，
// 让一个样本在秒级完成；3s 静置窗不动（冻结默认 SettleStaticDefault，被测的就是它）。
// 产线默认 90s 暖窗的形状在证据表里以分析记录，不进样本。

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// settle12r2Tts is the C9 TTS seam stub: one chunk, no synthesis cost, so the
// measured legs carry the walk and nothing else.
type settle12r2Tts struct{}

func (s *settle12r2Tts) Speak(text string, emit func([]float32) error) error {
	return emit(make([]float32, 1600)) // 100ms of silence at 16kHz: Drain has work
}

func (s *settle12r2Tts) SampleRate() int { return 16000 }

func (s *settle12r2Tts) Close() error { return nil }

// Test12r2SettleWalkFiveSamplesEndToEnd rides the production fall-back walk five
// times and reports every sample in milliseconds: the reply streams through the
// real sink (T0), the scope settles on its own timer (T1, the SCOPE-SETTLING
// announcement), the frozen 3s static expires and the dispose cascade stamps T2.
func Test12r2SettleWalkFiveSamplesEndToEnd(t *testing.T) {
	f := newRunFixture(t, "openai-chat")

	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "12-r2 settle harness",
		Tts:          voiceWireTtsFiles(t),
		WarmTimeout:  1200 * time.Millisecond,
	}
	rv, _ := voiceWireStart(t, plan, &residentBall{}, &voice15eSink{}, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return &voiceWireSpeaker{rate: rate}, nil }
		rv.residency.ttsLoader = func() (speech.TtsEngine, error) { return &settle12r2Tts{}, nil }
	})
	trace := rv.residency.settle
	// The same wiring the resident assembly root does (resident_windows.go's
	// sibling segment): the run's sink stamps into the voice leg's trace.
	f.rtHook = func(rt *agentRuntime) {
		rt.voiceAnnounce = rv.announce
		rt.settle = trace
	}

	const wantSamples = 5
	for i := 0; i < wantSamples; i++ {
		// The previous sample ended with the scope disposed; closeEpoch clears
		// the epoch so this sample's announce opens a fresh session scope.
		if err := rv.residency.closeEpoch("12r2 sample boundary"); err != nil {
			t.Fatalf("sample %d: closeEpoch: %v", i+1, err)
		}
		// The machine bridge, on the table's own rows: #4 (Sleeping/Warm -> Listening)
		// then #11 (Listening -> Thinking). From Thinking the announce's own #15 and
		// #26 are legal, which is the production walk into Warm.
		rv.dispatch(statemachine.EvSummon, nil, "test: D43#4/#29 让球到 Listening（12r2 样本桥）")
		if got := rv.machine.State(); got != statemachine.StateListening {
			t.Fatalf("sample %d: the bridge walk landed in %s, want Listening", i+1, string(got))
		}
		rv.dispatch(statemachine.EvVadStop, &statemachine.Facts{SpeechMS: 400}, "test: D43#11 让球到 Thinking（12r2 样本桥）")
		if got := rv.machine.State(); got != statemachine.StateThinking {
			t.Fatalf("sample %d: the bridge walk landed in %s, want Thinking", i+1, string(got))
		}

		code := f.run("总结一下 12r2 的 settle 样本")
		if code != 0 {
			t.Fatalf("sample %d: exit %d\nstdout:\n%s\nstderr:\n%s", i+1, code, f.out.String(), f.err.String())
		}

		sample := i
		voice15eWaitFor(t, "the reply's fall-back to land in the trace (T2)", func() bool {
			return len(trace.samplesSnapshot()) > sample
		})
		s := trace.samplesSnapshot()[sample]
		t.Logf("12R2-SETTLE sample=%d token_to_settling=%dms settling_to_fallback=%dms token_to_fallback=%dms",
			i+1, s.LastTokenToSettling.Milliseconds(), s.SettlingToFallback.Milliseconds(),
			s.LastTokenToFallback.Milliseconds())

		// AC#9 ①: the bound is per-sample and a bad tail is recorded, never
		// averaged away - hence Errorf per sample, not a summary assert.
		if s.LastTokenToSettling <= 0 || s.LastTokenToSettling > 3000*time.Millisecond {
			t.Errorf("AC#9 sample %d: 最后 token 到进 Settling %dms，超出 3000ms 判据线（照录 FAIL，不挑样本）",
				i+1, s.LastTokenToSettling.Milliseconds())
		}
		// The settle window itself: the frozen 3s static, never fires early, and
		// reaches the stamp through the owner queue. A reading outside 2.9-3.6s
		// means the walk measured something that is not the D43#32 window.
		if s.SettlingToFallback < 2900*time.Millisecond || s.SettlingToFallback > 3600*time.Millisecond {
			t.Errorf("sample %d: Settling 到回落 %dms，不是冻结 3s 静置窗（D43#32）的形状",
				i+1, s.SettlingToFallback.Milliseconds())
		}
	}
	if got := len(trace.samplesSnapshot()); got < wantSamples {
		t.Fatalf("the walk produced %d samples, want %d (AC#9 ②: every sample reported)", got, wantSamples)
	}
}

// Test12r2SeamCallSitesAreTheMeasurement is the always-armed half of the reverse
// criterion: a stamp call site removed from the production files reds here
// without any mutation run. The behavioural half (remove the stamps, the
// measurement test reds) is the -overlay pair recorded in the evidence file.
func Test12r2SeamCallSitesAreTheMeasurement(t *testing.T) {
	fset := token.NewFileSet()
	runAst, err := parser.ParseFile(fset, "run.go", nil, 0)
	if err != nil {
		t.Fatalf("parse run.go: %v", err)
	}
	enginesAst, err := parser.ParseFile(fset, "resident_voice_engines_windows.go", nil, 0)
	if err != nil {
		t.Fatalf("parse resident_voice_engines_windows.go: %v", err)
	}

	// T0: consoleSink.Publish's delta branch must stamp the last token.
	t0 := false
	ast.Inspect(runAst, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Publish" || fn.Recv == nil {
			return true
		}
		ast.Inspect(fn.Body, func(m ast.Node) bool {
			ce, ok := m.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := ce.Fun.(*ast.SelectorExpr)
			if ok && sel.Sel.Name == "markLastToken" {
				t0 = true
			}
			return true
		})
		return false
	})
	if !t0 {
		t.Fatal("12r2: run.go's consoleSink.Publish no longer stamps the last token (markLastToken), " +
			"so the AC#9 measurement has no T0 and the 3s settle record is unmeasurable")
	}

	// T1/T2: openScope must stamp on the scope's own settling announcement and
	// on the dispose cascade's memory-return step, anchored on the scope's own
	// audit prefix.
	t1, t2, tag := false, false, false
	ast.Inspect(enginesAst, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "openScope" || fn.Body == nil {
			return true
		}
		ast.Inspect(fn.Body, func(m ast.Node) bool {
			switch x := m.(type) {
			case *ast.CallExpr:
				if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "markSettling" {
						t1 = true
					}
					if sel.Sel.Name == "markFellBack" {
						t2 = true
					}
				}
			case *ast.BasicLit:
				if x.Value == `"session: SCOPE-SETTLING"` {
					tag = true
				}
			}
			return true
		})
		return false
	})
	if !t1 || !tag {
		t.Fatal("12r2: openScope no longer stamps T1 on the scope's SCOPE-SETTLING announcement, " +
			"so the settle entry leg of AC#9 item 1 is unmeasurable")
	}
	if !t2 {
		t.Fatal("12r2: openScope no longer stamps T2 on the dispose cascade (FreeOSMemory step), " +
			"so the fall-back leg (Settling->dispose) is unmeasurable")
	}
}
