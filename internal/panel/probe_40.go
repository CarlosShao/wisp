package panel

// Ticket 40 AC#3 (first two halves) - the model-management probe door's shared
// half: the roster name, the machine payload marker and the receipt renderer.
//
// WHAT THIS FILE OWNS AND WHAT IT DOES NOT. The route (config.model.probe) is
// declared in bridge.go and routed in composer_dispatch.go like every other
// door; the LEG that answers it lives in cmd/wisp, because only the assembly
// root can reach the three objects a real probe needs (the config Manager for
// the catalog, the secret store for the key, and the llm probe machinery that
// writes provider_health). This file is the part both sides share, so the page
// and the tests read ONE symbol per spelling - the same rule every payload
// family in this package follows (SettingsPayloadMarker, ModelsListPayloadMarker,
// ConfigReloadEventMarker).
//
// THE ONE RULE THE RECEIPT ENFORCES: a probe MEASURES, it never saves. Nothing
// it says may claim a write happened or that anything takes effect - the CRUD
// writes next to it are restart tier and say 需要重启 through tierSentence, and
// this receipt deliberately carries no tier sentence at all, because a receipt
// that said "需要重启" for a read-only measurement would be the same wrong-tier
// lie AC#8's wording rule forbids in the other direction.

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ModelProbePayloadMarker starts the machine half of a probe answer: one line
// the page finds by the marker and parses as JSON. Same two properties as the
// settings markers: no `panel.` prefix (the drift instruments key on those
// literals) and no dot (it is neither a method name nor a route).
const ModelProbePayloadMarker = "wisp-model-probe:"

// ModelProbeOutcome is one capability's verdict as the wire carries it:
// declared is the catalog's claim, ok is what the wire demonstrated, detail is
// the probe's own failure text VERBATIM (an empty detail with ok=true is the
// normal answer; a detail with ok=false is the error原文 the ticket names -
// the page shows it as-is, never paraphrased).
type ModelProbeOutcome struct {
	Capability string `json:"capability"`
	Declared   bool   `json:"declared"`
	OK         bool   `json:"ok"`
	// Attempted is false when no request could be sent at all (audio today);
	// then OK carries no verdict and the page must not read it as a failure.
	Attempted bool   `json:"attempted"`
	Detail    string `json:"detail,omitempty"`
}

// ModelProbeMismatch is one 「声明可用，实测不可用」 fact. Label is the frozen
// copy (llm.ProbeMismatch.Label's own sentence); Detail is the probe's own
// failure text verbatim. Both are shown, never re-worded.
type ModelProbeMismatch struct {
	Capability string `json:"capability"`
	Detail     string `json:"detail"`
	Label      string `json:"label"`
}

// ModelProbeResult is the whole probe answer as the leg reports it. Error is
// the REAL failure text when the run could not complete (an endpoint that
// never answered, a key that would not resolve) - an empty Error with
// Outcomes is the only shape that may read as a success, so a stubbed "OK"
// with nothing behind it has no wire shape it can wear.
type ModelProbeResult struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	LatencyMS int64  `json:"latencyMS"`
	AllOK     bool   `json:"allOK"`
	// Written says the run's verdicts reached provider_health (ticket 09's
	// sink). false = the measurement happened but nothing was recorded (or
	// the run failed outright) - a state the page must show, never fold into
	// "measured fine".
	Written    bool                 `json:"written"`
	Outcomes   []ModelProbeOutcome  `json:"outcomes"`
	Mismatches []ModelProbeMismatch `json:"mismatches"`
	Error      string               `json:"error,omitempty"`
}

// RenderModelProbePayload renders the machine half of a probe answer: one line
// beginning with ModelProbePayloadMarker. Outcomes/Mismatches are always
// arrays on the wire, never null.
func RenderModelProbePayload(r ModelProbeResult) string {
	if r.Outcomes == nil {
		r.Outcomes = []ModelProbeOutcome{}
	}
	if r.Mismatches == nil {
		r.Mismatches = []ModelProbeMismatch{}
	}
	data, err := json.Marshal(r)
	if err != nil {
		// A struct of strings, bools and ints cannot fail to marshal; if a
		// later addition ever makes that untrue, the page is told rather than
		// handed a silently empty line.
		return fmt.Sprintf("%s{\"payloadError\":%q}", ModelProbePayloadMarker, err.Error())
	}
	return ModelProbePayloadMarker + string(data)
}

// renderModelProbeWord is one capability's verdict word. ASCII/Chinese words
// only: the console PASS/FAIL precedent (ticket 67 AC#3) and the D22 emoji ban
// both forbid check/cross glyphs, so the verdict is a WORD.
func renderModelProbeWord(o ModelProbeOutcome) string {
	if !o.Attempted {
		return "未实测"
	}
	if o.OK {
		return "实测可用"
	}
	return "实测不可用"
}

// RenderModelProbeReceipt is the page-visible sentence for one probe run: what
// was measured (provider/model verbatim), the measured p50 latency, one verdict
// word per attempted capability with the declared/实测 pairing, every mismatch's
// frozen label and its detail原文, and the provider_health half. It says NOTHING
// about saving or taking effect: a probe writes no config key, and the CRUD
// writes beside it own the tier sentences.
func RenderModelProbeReceipt(r ModelProbeResult) string {
	var b strings.Builder
	if r.Error != "" {
		// The failure path: the error原文 IS the answer. The prose names what
		// was attempted so the user's own click is traceable, and the raw text
		// rides verbatim - a paraphrase would be the second judge the
		// RetryNotice rule forbids.
		fmt.Fprintf(&b, "探测未完成：%s/%s —— %s", r.Provider, r.Model, r.Error)
		return b.String()
	}
	fmt.Fprintf(&b, "探测完成：%s/%s 实测 p50 %dms", r.Provider, r.Model, r.LatencyMS)
	attempted := 0
	for _, o := range r.Outcomes {
		if !o.Attempted {
			continue
		}
		attempted++
		declared := "未声明"
		if o.Declared {
			declared = "声明可用"
		}
		fmt.Fprintf(&b, "；%s %s %s", o.Capability, declared, renderModelProbeWord(o))
		if o.Detail != "" && !o.OK {
			fmt.Fprintf(&b, "（细节：%s）", oneLine40(o.Detail))
		}
	}
	if attempted == 0 {
		b.WriteString("；没有任何能力发出实测请求")
	}
	for _, m := range r.Mismatches {
		fmt.Fprintf(&b, "；%s", m.Label)
		if m.Detail != "" {
			fmt.Fprintf(&b, "（细节：%s）", oneLine40(m.Detail))
		}
	}
	if r.Written {
		b.WriteString("。判定已写入 provider_health。")
	} else {
		b.WriteString("。判定没有写入 provider_health（本次只测未录）。")
	}
	return b.String()
}

// oneLine40 keeps a probe detail on a single receipt line without hiding it -
// the same join the CLI's oneLine does (providers.go), because a multi-line
// provider error would otherwise break the receipt's one-sentence shape.
func oneLine40(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
