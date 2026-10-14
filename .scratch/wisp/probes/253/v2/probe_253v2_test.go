// 253-v2 acceptance leg: my own Q&A of the RUNNING guard, compiled into package
// panel through `go test -overlay` only (this file never exists at that path on
// disk). It re-implements no judgement: it calls the real knownComposerMethod.
package panel

import "testing"

func TestV2RunningGuardQandA253v2(t *testing.T) {
	six := []string{
		"panel.mode.request", "panel.workspace.request", "panel.attachment.add",
		"panel.message.send", "config.get", "config.set",
	}
	for _, n := range six {
		if !knownComposerMethod(n) {
			t.Errorf("253-v2: the RUNNING guard refuses the roster name %q", n)
		}
	}
	off := []string{
		"config.peek", "settings.read253v2", "config.toml", "index.html",
		"panel.mode.request.extra", "panel-composer",
	}
	for _, n := range off {
		if knownComposerMethod(n) {
			t.Errorf("253-v2: the RUNNING guard answers the off-roster spelling %q", n)
		}
	}
	t.Logf("253-v2 running-guard Q&A: the six roster names answered verbatim; the six off-roster spellings (%q) refused", off)
}
