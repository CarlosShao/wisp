package agent

import (
	"strings"
)

// D11(1) session-control layer: the five control words are executed locally,
// WITHOUT an LLM round-trip (SPEC-05 §2: "命中控制词（停/取消/重说/大声点/
// 确认）→ 直接执行控制语义，不经 LLM（几十毫秒级）"). The word set is exactly
// the D11 set; anything else is business intent and goes to the loop.
//
// D11(2) removed every pre-classification layer: business intent is decided by
// the model itself through tool-calling. D11(3) (force_tool/force_chat TOML
// rule table) is a user-preference injection point and is not implemented by
// this ticket (DEFERRED(D11-3), see the ticket report).

// ControlVerb is one D11(1) control semantic.
type ControlVerb string

const (
	// ControlStop halts the current output/task ("停").
	ControlStop ControlVerb = "stop"
	// ControlCancel abandons the task and its effects ("取消").
	ControlCancel ControlVerb = "cancel"
	// ControlRepeat re-presents the last result ("重说").
	ControlRepeat ControlVerb = "repeat"
	// ControlLouder raises the spoken output ("大声点").
	ControlLouder ControlVerb = "louder"
	// ControlConfirm answers a pending confirmation with yes ("确认").
	ControlConfirm ControlVerb = "confirm"
)

// controlWords maps each verb to its utterance. The mapping is one-to-one:
// every entry is a control word of D11(1), no synonyms are invented here
// (expanding the vocabulary is an i18n/UX decision, ticket 58).
var controlWords = map[string]ControlVerb{
	"停":   ControlStop,
	"取消":  ControlCancel,
	"重说":  ControlRepeat,
	"大声点": ControlLouder,
	"确认":  ControlConfirm,
}

// controlTrim strips the punctuation an ASR transcript may carry around a
// bare control word. It removes only leading/trailing separators, never
// interior text, so "取消 这个任务" stays business intent.
const controlTrim = " \t\r\n。，,、.!?！？；;：:"

// MatchControl reports whether text is a bare control utterance and returns
// the verb it maps to.
func MatchControl(text string) (ControlVerb, bool) {
	t := strings.Trim(text, controlTrim)
	if t == "" {
		return "", false
	}
	v, ok := controlWords[t]
	return v, ok
}

// ControlOutcome is what a ControlHandler reports back for one verb.
type ControlOutcome struct {
	// Handled is false when the consumer could not act on the verb (no task
	// running, nothing pending): the caller then reports Text to the user.
	Handled bool
	// Text is the user-visible acknowledgement (may be empty).
	Text string
}

// ControlHandler executes a control verb in the host (cancel a running task,
// re-speak, raise volume, answer a pending confirmation). Wiring to the ball /
// TTS / approval queue belongs to those tickets; the loop only guarantees it
// never spends an LLM round-trip for a control utterance.
type ControlHandler func(verb ControlVerb, utterance string) ControlOutcome
