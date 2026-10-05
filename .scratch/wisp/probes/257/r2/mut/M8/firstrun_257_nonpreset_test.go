package main

// Ticket 257, shape (c), the half of the receipt that the seven-field walk does
// NOT cover: the sentence about a provider name that is not one of the built-in
// presets. firstrun.go tells the owner two opposite things in the same breath -
// a preset name may leave protocol and base_url empty, a non-preset name has to
// state protocol itself "否则这份文件加载不过". That clause is the difference
// between guidance and a trap: an owner who invents a name (the common case for
// a self-hosted gateway) and follows the rest of the sentence gets a config.toml
// that no longer loads at all, which is the same failure mode A543 recorded for
// the [providers.x] spelling - just one section later.
//
// So this file does two things and each one has its own tooth (evidence
// .scratch/wisp/probes/257/r2/evidence.md §5, M7 and M8):
//
//   1. the receipt must say the clause out loud (M7: delete it, this goes red);
//   2. the clause must be TRUE of the real gate - a non-preset row with no
//      protocol fails to load, and with protocol it loads AND its three
//      provider-level fields become writable (M8: point the same code at a
//      preset name, the "must fail" step goes red, which is what proves the
//      behavioural half is not vacuous).
//
// Nothing here touches internal/config: the gate is exercised through the
// exported loader the production process itself uses (config.NewManager, the
// res=nil shape every production call site passes, ledger A560).

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
)

// t257npGhost is deliberately NOT a key of internal/config's providerPresets
// (defaults.go's map, 12 names as of this writing). Asserted below rather than
// trusted, because a future preset with this exact name would silently turn
// half of this case into the preset path.
const t257npGhost = "t257-r2-ghost-gateway"

// t257npCleanMachine is this file's own clean machine: empty temp root, premise
// asserted, ONE pass of the user-started run entry, and the file production
// wrote. Kept self-contained on purpose - it shares no helper with the sibling
// 257 file so neither half can be silenced by the other being renamed.
func t257npCleanMachine(t *testing.T) (dir, cfgPath, receipt string) {
	t.Helper()
	dir = t.TempDir()
	cfgPath = filepath.Join(dir, configFileName)
	if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("AC#1 premise broke: the clean machine already holds %q", cfgPath)
	}
	code, _, stderrText := run198(t, dir)
	if code != 2 {
		t.Fatalf("first run over a fresh root: exit %d, want 2, stderr:\n%s", code, stderrText)
	}
	if !strings.Contains(stderrText, "新建默认配置") {
		t.Fatalf("no creation receipt to read:\n%s", stderrText)
	}
	return dir, cfgPath, stderrText
}

// TestTicket257R2AC1ReceiptStatesTheNonPresetCondition is the prose half: the
// receipt must carry both sides of the preset rule, in the same sentence it
// prints once on a machine that never had a config.toml.
func TestTicket257R2AC1ReceiptStatesTheNonPresetCondition(t *testing.T) {
	_, _, receipt := t257npCleanMachine(t)

	for _, mustSay := range []string{
		"名字对上内置预设的，protocol 与 base_url 可以留空",
		"非预设名必须自己写 protocol",
		"否则这份文件加载不过",
	} {
		if !strings.Contains(receipt, mustSay) {
			t.Errorf("AC#1 RED: the receipt never states %q, so the guidance is only honest about preset names:\n%s",
				mustSay, receipt)
		}
	}
	// The clause is a condition, not an invitation to leave it out: a receipt that
	// folded both sides into "protocol 可以留空" would send a self-hosted gateway
	// straight into a file that cannot load.
	if strings.Contains(receipt, "非预设名也可以留空") {
		t.Errorf("AC#1 RED: the receipt now promises protocol may stay empty for a non-preset name:\n%s", receipt)
	}
	if names := config.PresetNames(); !slicesContainsLocal(names, t257npGhost) {
		t.Logf("premise holds: %q is absent from the %d built-in presets %v", t257npGhost, len(names), names)
	} else {
		t.Fatalf("premise broke: %q became a built-in preset, so the non-preset half of the receipt has nothing to warn about here", t257npGhost)
	}
}

// TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads is the behavioural
// half: it runs the two shapes the sentence describes against the real loader,
// on the file the production first run created.
func TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads(t *testing.T) {
	_, cfgPath, receipt := t257npCleanMachine(t)
	base := t257npRead(t, cfgPath)

	// Shape 1: what an owner gets if they invent a name and follow only the
	// "留空" half. The file must be refused, and refused for the reason the
	// receipt names - not by some unrelated gate.
	noProtocol := base + "\n[llm.providers." + t257Provider + "]\napi_key_ref = 'env:T257_R2_GHOST_NEVER_SET'\n"
	if err := os.WriteFile(cfgPath, []byte(noProtocol), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := config.NewManager(cfgPath, nil)
	if err == nil {
		t.Fatal("AC#1 RED: a non-preset provider row with no protocol loaded, so the receipt's warning is describing a gate that is not there")
	}
	for _, mustBlame := range []string{"protocol", t257npGhost} {
		if !strings.Contains(err.Error(), mustBlame) {
			t.Errorf("AC#1 RED: the refusal does not name %q, so the receipt's clause and the real gate are not the same rule: %v", mustBlame, err)
		}
	}

	// Shape 2: the same invented name WITH the protocol the receipt says to
	// write. It must load, and the provider-level fields must then be writable -
	// this is what makes the clause guidance rather than a warning that the
	// non-preset route is closed.
	withProtocol := noProtocol + "\n[llm.providers." + t257npGhost + ".models.ghost-chat]\nenabled = true\n"
	withProtocol = strings.Replace(withProtocol,
		"[llm.providers."+t257npGhost+"]\napi_key_ref = 'env:T257_R2_GHOST_NEVER_SET'\n",
		"[llm.providers."+t257npGhost+"]\nprotocol = \""+config.ProtocolOpenAIChat+"\"\napi_key_ref = 'env:T257_R2_GHOST_NEVER_SET'\n", 1)
	if err := os.WriteFile(cfgPath, []byte(withProtocol), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("AC#1 RED: the shape the receipt teaches for a non-preset name does not load: %v", err)
	}
	if p := m.Config().LLM.Providers[t257npGhost]; p.Protocol != config.ProtocolOpenAIChat {
		t.Fatalf("the protocol the owner wrote did not survive the load: %#v", p)
	}
	if _, err := m.SetProviderBaseURL(t257npGhost, "https://ghost-gateway.example.invalid/v1"); err != nil {
		t.Errorf("AC#1 RED: a preset-named row is writable but an invented one is not, which the receipt never says: %v", err)
	}
	if _, err := m.SetProviderAPIKeyRef(t257npGhost, "dpapi:t257r2ghostcanaryref"); err != nil {
		t.Errorf("AC#1 RED: the credential-reference field refused an invented provider: %v", err)
	}

	// AC#3's line stays where it is: both shapes above carried a reference and
	// never a value, and neither run put a plaintext key on disk.
	after := t257npRead(t, cfgPath)
	if strings.Contains(after, "\napi_key =") || strings.Contains(after, "sk-") {
		t.Errorf("AC#3 RED: the non-preset walk put a value-shaped credential on disk:\n%s", after)
	}
	// The receipt the operator actually read still promises the reference form,
	// so this path did not silently become a second credential door.
	if !strings.Contains(receipt, "dpapi:<blob 名>") || !strings.Contains(receipt, "env:<环境变量名>") {
		t.Errorf("AC#3: the receipt stopped naming the two reference forms as placeholders:\n%s", receipt)
	}
}

// --- helpers ---------------------------------------------------------------

func t257npRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}

// slicesContainsLocal keeps the premise check free of an extra import in a file
// whose whole job is two sentences of operator-facing text.
func slicesContainsLocal(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}
