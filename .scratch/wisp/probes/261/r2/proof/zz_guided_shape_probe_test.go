package llm_test

// Ticket 261 leg r2 - the never-vacuous self-proof for fix A (firstrun
// guidance). It exists only through -overlay (proof/overlay.json maps it onto
// internal/llm/zz_guided_shape_probe_test.go); the tracked tree never carries
// this file.
//
// What it proves: the entry shape the firstrun receipt now TEACHES - a
// hand-written [llm.providers.<名>].models.<id> entry carrying an explicit
// `enabled = true` - really passes the enabled gate, and the SAME shape minus
// that one line is refused. Without the positive half, the guidance would be
// teaching a shape that still dies at the gate; without the negative half, the
// positive half could pass for reasons unrelated to the flag.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
)

const guidedShapeTOML = `schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.providers.acme]
protocol = "openai-chat"
base_url = "http://127.0.0.1:9/v1"

[llm.providers.acme.models.m1]
GUIDED_LINE
context_window = 128000
`

func loadGuidedShape(t *testing.T, guidedLine string) (*config.Config, error) {
	t.Helper()
	body := strings.Replace(guidedShapeTOML, "GUIDED_LINE", guidedLine, 1)
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, _, err := config.LoadFile(path, nil)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// TestZZ261R2GuidedShapePassesTheGate is the hard requirement from the
// dispatch: after fix A, the shape the guidance teaches must resolve through
// the gate, and the shape without the taught line must not.
func TestZZ261R2GuidedShapePassesTheGate(t *testing.T) {
	// The guided shape: exactly the line the receipt tells the user to write.
	cfgOn, err := loadGuidedShape(t, "enabled = true")
	if err != nil {
		t.Fatalf("LoadFile refused the guided shape: %v", err)
	}
	eps, err := llm.NewResolver(cfgOn, nil).ResolveChain()
	if err != nil {
		t.Fatalf("the GUIDED shape (enabled = true) did not pass the gate: %v", err)
	}
	if len(eps) != 1 || eps[0].Model != "m1" {
		t.Fatalf("guided shape resolved %+v, want one endpoint for m1", eps)
	}

	// The no-key shape: the same file minus the taught line. It decodes to
	// false (defaults never enter maps) and must be refused - by name, with
	// the guidance sentence fix B added.
	cfgOff, err := loadGuidedShape(t, "")
	if err != nil {
		t.Fatalf("LoadFile refused the no-key shape: %v", err)
	}
	_, err = llm.NewResolver(cfgOff, nil).ResolveChain()
	if err == nil {
		t.Fatal("the no-key shape passed the gate; the missing-key guidance teaches a lie")
	}
	for _, want := range []string{"disabled", "缺 enabled 键的条目视为关闭", "enabled = true"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("no-key refusal = %q, want it to carry %q", err.Error(), want)
		}
	}
}
