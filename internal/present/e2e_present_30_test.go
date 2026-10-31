package present

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 30 AC#2 and AC#6: four crafted results through the REAL pipeline -
// the real agent loop (real C5 seam, real turn assembly, real EvDone) with
// the router's channels observed per tier.
//
// The C5 seam is the sanctioned injection face (AGENTS 1.3). Three samples
// ride the real mockllm subprocess over the real openai-chat adapter: its
// non-golden mode echoes the last user text back, so the crafted sample IS
// what the loop's final text is made of - and with zero tools plus an empty
// scene the loop appends no volatile suffix message, so the echo is the
// sample, not the scene block. The >400-hanzi sample cannot ride the echo:
// mockllm caps the echo at 400 bytes (tools/mockllm/chat.go, the len(answer)
// clamp), and 402 hanzi is 1206. That one rides an in-process C5 provider
// instead - the same loop, the same router, the request count counted at the
// seam - and no golden fixture is added anywhere (testdata is outside this
// leg's write surface).
//
// The 口播摘要 line rides inside the produced response exactly the way a
// compliant model would produce it from the (6) style section's instruction,
// which is what proves the same-response summary rule end to end.

type noTools struct{}

func (noTools) Tools(context.Context) ([]agent.ToolInfo, error) { return nil, nil }

func (noTools) Execute(context.Context, agent.ToolRequest) (agent.ToolOutcome, error) {
	return agent.ToolOutcome{}, nil // the echo pipeline never calls a tool
}

// scriptedC5 is the in-process C5 provider: one scripted turn, a counted
// request. Only the >400-hanzi sample needs it (see the file comment).
type scriptedC5 struct {
	text  string
	calls int
}

func (p *scriptedC5) Stream(_ context.Context, _ *llm.Request, emit func(llm.StreamEvent) error) error {
	p.calls++
	if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: p.text}); err != nil {
		return err
	}
	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

func (p *scriptedC5) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Protocol: "openai-chat", Provider: "scripted",
		Model: "scripted-long", MaxContextWindow: 128000,
	}
}

func Test30FourResultsThroughTheRealPipeline(t *testing.T) {
	proc := adaptertest.StartMockllm(t) // skips under -short
	store := openStore(t)

	ep := llm.EndpointOptions{
		Provider: "mock", Model: "mock-small", Protocol: "openai-chat",
		BaseURL: proc.Base + "/v1", ContextWindow: 128000,
		HTTPClient: &http.Client{Transport: &http.Transport{Proxy: nil}},
	}
	mockProv, err := llm.NewProvider(ep)
	if err != nil {
		t.Fatalf("build provider: %v", err)
	}

	newLoop := func(t *testing.T, prov llm.LlmProvider, sink *agent.RecordSink) *agent.Loop {
		t.Helper()
		loop, err := agent.New(agent.Options{
			Provider: prov,
			Tools:    noTools{},
			Sink:     sink,
			Registry: observe.NewRegistry(),
			Config: agent.Config{
				Model:         "mock-small",
				ContextWindow: 128000,
				SceneFunc:     func() agent.Scene { return agent.Scene{} },
			},
		})
		if err != nil {
			t.Fatalf("agent.New: %v", err)
		}
		return loop
	}

	type sample struct {
		name   string
		prov   llm.LlmProvider
		prompt string
		check  func(t *testing.T, res agent.Result, rc Receipt, calls int)
	}
	samples := []sample{
		{
			name:   "tier1_table_with_marker_over_mockllm",
			prov:   mockProv,
			prompt: "生成结果如下：\n| 名称 | 数量 |\n|---|---|\n| 甲 | 1 |\n口播摘要：表格已存档，请查收",
			check: func(t *testing.T, res agent.Result, rc Receipt, calls int) {
				if calls != 1 {
					t.Fatalf("mockllm chat requests = %d, want 1: the summary must ride the same response, never a second completion", calls)
				}
				if !strings.Contains(res.Text, "口播摘要：表格已存档，请查收") {
					t.Fatalf("the response lost the same-response summary: %q", res.Text)
				}
				if rc.Tier != TierStructured {
					t.Fatalf("tier = %s (%s), want structured", rc.Tier, rc.Rule)
				}
				if !strings.HasPrefix(rc.Spoken, "表格已存档，请查收") {
					t.Errorf("spoken = %q, want the model's own summary", rc.Spoken)
				}
				if !strings.HasSuffix(rc.Spoken, rc.Artifact) {
					t.Errorf("spoken = %q, want the path announced", rc.Spoken)
				}
				if rc.Display != "echo: 生成结果如下：\n| 名称 | 数量 |\n|---|---|\n| 甲 | 1 |" {
					t.Errorf("display = %q, want the echoed body minus the marker line", rc.Display)
				}
				data, err := os.ReadFile(rc.Artifact)
				if err != nil {
					t.Fatalf("artifact: %v", err)
				}
				if string(data) != rc.Display {
					t.Errorf("artifact body = %q, want the display text", data)
				}
			},
		},
		{
			name:   "tier2_long_hanzi_with_marker_over_the_c5_seam",
			prov:   &scriptedC5{text: strings.Repeat("测", 402) + "\n口播摘要：长文已归档"},
			prompt: "把这份长文整理好",
			check: func(t *testing.T, res agent.Result, rc Receipt, calls int) {
				if calls != 1 {
					t.Fatalf("provider requests = %d, want 1", calls)
				}
				if !strings.Contains(res.Text, "口播摘要：长文已归档") {
					t.Fatalf("the response lost the summary: %q", truncRunesForTest(res.Text, 80))
				}
				if rc.Tier != TierLong {
					t.Fatalf("tier = %s (%s), want long", rc.Tier, rc.Rule)
				}
				if !strings.HasPrefix(rc.Spoken, "长文已归档") {
					t.Errorf("spoken = %q, want the marker summary", rc.Spoken)
				}
				if rc.Artifact == "" {
					t.Fatal("tier 2 wrote no artifact")
				}
				data, err := os.ReadFile(rc.Artifact)
				if err != nil {
					t.Fatalf("artifact: %v", err)
				}
				if string(data) != rc.Display {
					t.Errorf("artifact body = %q, want the display text minus the marker", truncRunesForTest(string(data), 80))
				}
			},
		},
		{
			name:   "tier3_medium_first_sentence_dual_channel_over_mockllm",
			prov:   mockProv,
			prompt: strings.Repeat("测", 61) + "。第二句在这里。",
			check: func(t *testing.T, res agent.Result, rc Receipt, calls int) {
				if calls != 1 {
					t.Fatalf("mockllm chat requests = %d, want 1", calls)
				}
				if rc.Tier != TierMedium {
					t.Fatalf("tier = %s (%s), want medium", rc.Tier, rc.Rule)
				}
			},
		},
		{
			name:   "tier4_short_clipboard_and_speak_over_mockllm",
			prov:   mockProv,
			prompt: "今天天气不错，适合出门散步。",
			check: func(t *testing.T, res agent.Result, rc Receipt, calls int) {
				if calls != 1 {
					t.Fatalf("mockllm chat requests = %d, want 1", calls)
				}
				if rc.Tier != TierShort {
					t.Fatalf("tier = %s (%s), want short", rc.Tier, rc.Rule)
				}
			},
		},
	}

	for _, s := range samples {
		t.Run(s.name, func(t *testing.T) {
			c, sp, nf, bg, tl, pn, cl := testChannels()
			c.Artifacts = store
			router := New(c)
			sink := &agent.RecordSink{}
			loop := newLoop(t, s.prov, sink)

			var calls int
			if sc, ok := s.prov.(*scriptedC5); ok {
				calls = sc.calls
			} else {
				proc.Reset(t)
			}
			res := loop.Run(context.Background(), s.prompt)
			if sc, ok := s.prov.(*scriptedC5); ok {
				calls = sc.calls
			} else {
				calls = proc.RouteCount(t, "chat")
			}
			if res.Status != agent.StatusCompleted {
				t.Fatalf("loop status = %s (%s), want completed", res.Status, res.Message)
			}
			if sink.LastOf(agent.EvDone) == nil {
				t.Fatal("the loop published no EvDone: the pipeline did not finish")
			}
			rc := router.Dispatch(context.Background(),
				Input{TaskID: res.TaskID, Text: res.Text})
			if !rc.OK {
				t.Fatalf("receipt not OK: %+v", rc.Notes)
			}
			s.check(t, res, rc, calls)

			// The per-tier channel observations that do not depend on the
			// sample, restated here so every subtest reads the same table.
			switch rc.Tier {
			case TierStructured, TierLong:
				if cl.calls != 0 {
					t.Errorf("tier 1/2 touched the clipboard %d times, want 0", cl.calls)
				}
				if len(pn.texts) != 1 || pn.arts[0] != rc.Artifact {
					t.Errorf("panel = %q/%q, want display plus artifact path", truncRunesForTest(pn.texts[0], 40), pn.arts)
				}
			case TierMedium:
				if nf.calls != 1 {
					t.Errorf("notify calls = %d, want 1", nf.calls)
				}
				if len(bg.texts) != 1 || len(tl.texts) != 1 {
					t.Errorf("badge=%d task=%d, want the unconditional dual-channel updates",
						len(bg.texts), len(tl.texts))
				}
				if len(pn.texts) != 1 || pn.arts[0] != "" {
					t.Errorf("panel = %q/%q, want text only", truncRunesForTest(pn.texts[0], 40), pn.arts)
				}
				if runes := len([]rune(rc.Spoken)); runes == 0 || runes > MaxSpokenRunes {
					t.Errorf("spoken = %q (%d runes), want a first sentence within the ceiling",
						truncRunesForTest(rc.Spoken, 40), runes)
				}
				if cl.calls != 0 {
					t.Errorf("tier 3 touched the clipboard %d times, want 0", cl.calls)
				}
			case TierShort:
				if cl.calls != 1 || cl.texts[0] != rc.Display {
					t.Errorf("clipboard = %d calls, want exactly the display text", cl.calls)
				}
				if got := sp.spoken(); len(got) != 1 || got[0] != rc.Display {
					t.Errorf("spoken = %q, want the whole short result", got)
				}
				if nf.calls != 0 || len(bg.texts) != 0 || len(tl.texts) != 0 || len(pn.texts) != 0 {
					t.Errorf("tier 4 pulled another channel: notify=%d badge=%d task=%d panel=%d",
						nf.calls, len(bg.texts), len(tl.texts), len(pn.texts))
				}
			}
		})
	}
}

func truncRunesForTest(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
