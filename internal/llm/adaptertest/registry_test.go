package adaptertest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/anthropic"
	"github.com/CarlosShao/wisp/internal/llm/openaichat"
	"github.com/CarlosShao/wisp/internal/llm/openairesponses"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Cross-adapter proofs (ticket 11 AC#1). These live in the harness package so
// they see ALL THREE adapters at once, which no single adapter package can.

// threeUnits builds one Unit per registered protocol.
func threeUnits() []Unit {
	return []Unit{
		{
			Name: "openai-chat", Prefix: "chat-", EndpointPath: "/chat/completions",
			New: func(t *testing.T, base string, o Options) llm.LlmProvider {
				return openaichat.New(chatOptions(base, o))
			},
			Fixture: func(sc ScenarioID) string {
				switch sc {
				case ScNoTerminal:
					return "missing-finish"
				case ScCancel:
					return "long-text"
				case ScDisconnect, ScMaxTokens, ScToolCall, ScUsageMultichunk,
					ScBackoff429, ScProvider500, ScUnauthorized, ScMissingUsage:
					return string(sc) // ticket 09 canonical names
				default:
					return "chat-" + string(sc)
				}
			},
		},
		{
			Name: "anthropic", Prefix: "anthropic-", EndpointPath: "/messages",
			New: func(t *testing.T, base string, o Options) llm.LlmProvider {
				return anthropic.New(anthropicOptions(base, o))
			},
			Fixture: dialectFixture("anthropic-"),
		},
		{
			Name: "openai-responses", Prefix: "responses-", EndpointPath: "/responses",
			New: func(t *testing.T, base string, o Options) llm.LlmProvider {
				return openairesponses.New(responsesOptions(base, o))
			},
			Fixture: dialectFixture("responses-"),
		},
	}
}

func dialectFixture(prefix string) func(ScenarioID) string {
	return func(sc ScenarioID) string {
		if sc == ScCancel {
			return prefix + "long-text"
		}
		return prefix + string(sc)
	}
}

// TestAllThreeProtocolsRegistered pins the C5 registry state ticket 11 had to
// reach (spec: three protocols behind one seam).
func TestAllThreeProtocolsRegistered(t *testing.T) {
	got := llm.Protocols()
	want := []string{"anthropic", "openai-chat", "openai-responses"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("registered protocols = %v, want %v", got, want)
	}
}

// TestSharedTableCoversEveryAdapter asserts the anti-shortcut property itself:
// each adapter package runs the SAME table, and each row has a fixture for
// each of the three dialects.
func TestSharedTableCoversEveryAdapter(t *testing.T) {
	for _, u := range threeUnits() {
		AssertCoverage(t, u)
	}
}

// TestC6EventsIdenticalAcrossAdapters is the load-bearing AC#1 proof: for the
// scenarios whose wire bytes differ only in dialect, the emitted C6 event
// sequence (kinds and payloads, ignoring EvUsage position, which the dialects
// legitimately report at different points) must be byte-identical across all
// three adapters.
func TestC6EventsIdenticalAcrossAdapters(t *testing.T) {
	cases := map[ScenarioID]string{
		ScToolCall:        "tool calls + text + stop",
		ScMaxTokens:       "truncation",
		ScThinking:        "reasoning then text",
		ScDisconnect:      "mid-stream disconnect",
		ScMidStreamError:  "error inside a 200 stream",
		ScUsageMultichunk: "usage aggregation",
	}
	for _, sc := range AllScenarios() {
		desc, ok := cases[sc]
		if !ok {
			continue
		}
		t.Run(string(sc), func(t *testing.T) {
			var refName string
			var ref []string
			for _, u := range threeUnits() {
				u := u
				base, _ := Serve(t, u, sc, false)
				strict := u.New(t, base, Options{})
				looseBase, _ := Serve(t, u, sc, false)
				loose := u.New(t, looseBase, Options{Loose: true})

				events := runBoth(t, strict, loose, sc)
				if refName == "" {
					refName, ref = u.Name, events
					continue
				}
				if !reflect.DeepEqual(ref, events) {
					t.Fatalf("%s (%s): event sequences differ\n  %s: %v\n  %s: %v",
						sc, desc, refName, ref, u.Name, events)
				}
			}
		})
	}
}

// trace renders one adapter's event stream into a canonical, dialect-free
// sequence. Two things are normalized on purpose, because they are genuinely
// wire-dialect detail and NOT part of the C6 contract:
//   - a tool call's argument fragments merge into one entry at the position of
//     its ToolCallStart (dialects differ in where ToolCallEnd fires: chat
//     closes all calls at finish_reason, Messages closes each block, Responses
//     closes each item); the entry still carries the assembled arguments and
//     whether an End was seen,
//   - an error contributes only its D37 class (the raw upstream code differs
//     per protocol; TestFaultClassIsAdapterIndependent pins the class).
func trace(events []llm.StreamEvent) []string {
	var out []string
	type call struct {
		id, name string
		args     []byte
		ended    bool
	}
	var calls []*call
	index := map[string]*call{}
	flush := func() {
		for _, c := range calls {
			out = append(out, "tool_args:"+c.id+":"+string(c.args)+
				":ended="+boolStr(c.ended))
		}
		calls = nil
		index = map[string]*call{}
	}
	for _, ev := range events {
		switch ev.Type {
		case llm.EvUsage:
			// position differs per dialect; values are asserted by the table
		case llm.EvTextDelta:
			flush()
			out = append(out, "text:"+ev.Text)
		case llm.EvReasoningDelta:
			flush()
			out = append(out, "reasoning:"+ev.Text)
		case llm.EvToolCallStart:
			c := &call{id: ev.ToolCallID, name: ev.ToolName}
			calls = append(calls, c)
			index[ev.ToolCallID] = c
			out = append(out, "tool_start:"+ev.ToolCallID+":"+ev.ToolName)
		case llm.EvToolCallArgsDelta:
			if c, ok := index[ev.ToolCallID]; ok {
				c.args = append(c.args, ev.ArgsDelta...)
			}
		case llm.EvToolCallEnd:
			if c, ok := index[ev.ToolCallID]; ok {
				c.ended = true
			}
		case llm.EvStop:
			flush()
			out = append(out, "stop:"+string(ev.Stop))
		case llm.EvError:
			flush()
			cls := "unclassified"
			if ev.Err != nil {
				cls = string(ev.Err.Class)
			}
			out = append(out, "error:"+cls)
		case llm.EvDone:
			flush()
			out = append(out, "done")
		}
	}
	flush()
	return out
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// runBoth streams with the strict and the loose adapter (failure scenarios
// are strict-only for their error class; the loose leg proves tolerance) and
// returns the canonical trace of the strict run.
func runBoth(t *testing.T, strict, loose llm.LlmProvider, sc ScenarioID) []string {
	t.Helper()
	events, _, _ := Drain(t, strict, BaseRequest())
	kinds := trace(events)
	// The loose run must at least terminate cleanly (the compat switch is part
	// of the shared contract, so it is exercised for every dialect here).
	looseEvents, _, _ := Drain(t, loose, BaseRequest())
	if len(looseEvents) == 0 {
		t.Fatalf("scenario %s: loose adapter emitted nothing", sc)
	}
	return kinds
}

// evSignature renders one event compactly and deterministically.
func evSignature(ev llm.StreamEvent) string {
	b, _ := json.Marshal(ev)
	return string(b)
}

// TestFaultClassIsAdapterIndependent answers the "did every adapter classify
// the SAME upstream failure the SAME way" question directly: one injected
// fault body per row, replayed through all three protocols against the live
// mockllm, must yield one identical D37 class each.
func TestFaultClassIsAdapterIndependent(t *testing.T) {
	tests := []struct {
		name string
		body string
		want observe.ErrorClass
	}{
		{"5xx", `{"type":"internal_server_error","error":{"type":"server_error","message":"upstream exploded","code":"internal_error"}}`, observe.ClassProvider},
		{"429", `{"type":"rate_limit_error","error":{"type":"rate_limit_error","message":"slow down","code":"rate_limit_exceeded"}}`, observe.ClassRateLimit},
		{"401", `{"type":"authentication_error","error":{"type":"authentication_error","message":"bad key","code":"invalid_api_key"}}`, observe.ClassAuth},
		{"402-quota", `{"type":"insufficient_quota","error":{"type":"insufficient_quota","message":"no credits","code":"insufficient_quota"}}`, observe.ClassBudget},
	}
	proc := StartMockllm(t)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var refName string
			var ref observe.ErrorClass
			for _, u := range threeUnits() {
				proc.Control(t, "/__control/reset", "{}")
				proc.QueueFault(t, Fault{Status: statusOf(tc.name), Times: 1, Body: tc.body})
				p := u.New(t, proc.Base+"/v1", Options{})
				_, _, err := Drain(t, p, BaseRequest())
				if err == nil {
					t.Fatalf("%s: expected a failure for %s", u.Name, tc.name)
				}
				got := ClassOf(err)
				if refName == "" {
					refName, ref = u.Name, got
				} else if got != ref {
					t.Fatalf("%s classified %s body as %s but %s as %s (same fault, two classes: %v)",
						u.Name, tc.name, got, refName, ref, err)
				}
				if got != tc.want {
					t.Errorf("%s/%s: class = %s, want %s (%v)", u.Name, tc.name, got, tc.want, err)
				}
				if got.Retryable() != tc.want.Retryable() {
					t.Errorf("%s/%s: retryable = %v, want %v (retry policy is class-derived)",
						u.Name, tc.name, got.Retryable(), tc.want.Retryable())
				}
			}
			proc.Control(t, "/__control/reset", "{}")
		})
	}
}

// statusOf maps a fault label onto the HTTP status to inject.
func statusOf(label string) int {
	switch label {
	case "5xx":
		return 500
	case "429":
		return 429
	case "401":
		return 401
	default: // 402-quota
		return 402
	}
}

// ---------------------------------------------------------------------------
// Endpoint options (fake keys only; this test never reads a real credential).
// ---------------------------------------------------------------------------

func optionsFor(base string, protocol string, o Options) llm.EndpointOptions {
	var ep llm.EndpointOptions
	ep.Provider = "mock"
	ep.Model = "mock-small"
	ep.Protocol = protocol
	ep.BaseURL = base
	ep.APIKey = "sk-test-not-a-real-key"
	ep.ContextWindow = 8192
	ep.Compat.Loose = o.Loose
	ep.Compat.AllowMissingUsage = o.AllowMissingUsage
	ep.HTTPClient = noProxyClient()
	return ep
}

func chatOptions(base string, o Options) llm.EndpointOptions {
	return optionsFor(base, "openai-chat", o)
}

func anthropicOptions(base string, o Options) llm.EndpointOptions {
	return optionsFor(base, "anthropic", o)
}

func responsesOptions(base string, o Options) llm.EndpointOptions {
	return optionsFor(base, "openai-responses", o)
}

// noProxyClient keeps machine proxy settings out of the harness.
func noProxyClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: nil}}
}
