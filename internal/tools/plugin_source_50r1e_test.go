package tools

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ac6Provider is a minimal C4 Manifest slot over already-built entries, so this
// file tests the registry + bridge half and leaves the manifest reading itself
// to internal/plugin's provider test.
type ac6Provider struct {
	kind    Kind
	entries []Entry
}

func (p *ac6Provider) Kind() Kind                  { return p.kind }
func (p *ac6Provider) List() []Entry               { return p.entries }
func (p *ac6Provider) Lookup(string) (Entry, bool) { return Entry{}, false }

func ac6Entry(name, pluginID string, provider Kind) Entry {
	return Entry{
		Tool: &fixtureTool{name: name},
		Decl: Decl{
			Capabilities: []Capability{CapNotify}, Needs: []Capability{CapNotify},
			Declared: risk.L0, Provider: provider, PluginID: pluginID,
		},
	}
}

// TestAC6BridgeExecutionRecordsPluginID is the AC6 execution half: a tool that
// came from a plugin manifest runs behind the same gate as a builtin, and its
// execution record names the plugin that ran it. Two surfaces carry the id (the
// Decision the host sees, and the audit line), and the builtin keeps both clean.
func TestAC6BridgeExecutionRecordsPluginID(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterProvider(NewBuiltinProvider(
		ac6Entry("builtin.ping", "", KindBuiltin))); err != nil {
		t.Fatal(err)
	}
	if err := reg.RegisterProvider(&ac6Provider{kind: KindManifest, entries: []Entry{
		ac6Entry("lark.search", "com.example.lark", KindManifest),
	}}); err != nil {
		t.Fatalf("the Manifest slot is landed by ticket 50: %v", err)
	}

	var lines []string
	var seen []Decision
	b := New(Options{
		Registry:   reg,
		Logf:       func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) },
		OnDecision: func(d Decision) { seen = append(seen, d) },
	})

	// 1) The plugin call: the audit trail must attribute it to its plugin.
	if _, err := b.Execute(t.Context(), req("lark.search", `{}`)); err != nil {
		t.Fatalf("a plugin call that is not refused must not return an error: %v", err)
	}
	var src []string
	for _, l := range lines {
		if strings.HasPrefix(l, "tools: PLUGIN-SOURCE ") {
			src = append(src, l)
		}
	}
	if len(src) != 1 {
		t.Fatalf("PLUGIN-SOURCE lines = %d (%v), want exactly 1 for the plugin call", len(src), lines)
	}
	for _, want := range []string{"plugin_id=com.example.lark", "tool=lark.search", "kind=manifest", "task=task-1"} {
		if !strings.Contains(src[0], want) {
			t.Errorf("PLUGIN-SOURCE line %q lacks %q", src[0], want)
		}
	}
	if len(seen) == 0 {
		t.Fatal("no Decision emitted")
	}
	if got := seen[len(seen)-1]; got.PluginID != "com.example.lark" || got.Provider != KindManifest {
		t.Errorf("Decision out of the plugin call = %+v, want PluginID com.example.lark / provider manifest", got)
	}

	// 2) The builtin call: same execution path, nothing to attribute.
	lines, seen = nil, nil
	if _, err := b.Execute(t.Context(), req("builtin.ping", `{}`)); err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "tools: PLUGIN-SOURCE ") {
			t.Errorf("a builtin execution must not claim a plugin: %q", l)
		}
	}
	if len(seen) == 0 || seen[0].PluginID != "" {
		t.Errorf("builtin Decision PluginID = %q, want empty", seen[0].PluginID)
	}

	// 3) A refused plugin call leaves no execution attribution: nothing executed.
	lines = nil
	reg2 := NewRegistry()
	if err := reg2.RegisterProvider(&ac6Provider{kind: KindManifest, entries: []Entry{
		ac6Entry("lark.write", "com.example.lark", KindManifest),
	}}); err != nil {
		t.Fatal(err)
	}
	b2 := New(Options{
		Registry: reg2,
		Logf:     func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) },
		// The host's view: this tool never declared notify, so C3 hard-rejects it.
		DeclaredCaps: func(string) []Capability { return []Capability{CapFSRead} },
	})
	out, err := b2.Execute(t.Context(), req("lark.write", `{}`))
	if err != nil || !out.IsError {
		t.Fatalf("an undeclared capability must be a refusal, not a fault: err=%v out=%+v", err, out)
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "tools: PLUGIN-SOURCE ") {
			t.Errorf("a call that never executed must not be attributed to a plugin: %q", l)
		}
	}
}

// TestAC6ManifestEntryWithoutAnOwnerIsRefused is the registration-side reverse
// check: the provenance gap cannot be created by a provider that forgets the id.
func TestAC6ManifestEntryWithoutAnOwnerIsRefused(t *testing.T) {
	reg := NewRegistry()
	err := reg.RegisterProvider(&ac6Provider{kind: KindManifest, entries: []Entry{
		ac6Entry("lark.search", "", KindManifest),
	}})
	if err == nil {
		t.Fatal("RegisterProvider accepted a manifest entry with no plugin id")
	}
	if !strings.Contains(err.Error(), "plugin id") || !strings.Contains(err.Error(), "lark.search") {
		t.Fatalf("err = %v, want it to name the entry and the missing plugin id", err)
	}
	if _, ok := reg.Lookup("lark.search"); ok {
		t.Error("the refused entry leaked into the directory")
	}
}

// TestAC6ManifestIsNoLongerAnUnlandedSlot pins the other half of the flip: the
// goja and MCP slots still answer the checkable placeholder error, and an
// unknown slot name is still refused rather than quietly accepted.
func TestAC6ManifestIsNoLongerAnUnlandedSlot(t *testing.T) {
	reg := NewRegistry()
	for _, kind := range []Kind{KindGoja, KindMCP} {
		err := reg.RegisterProvider(&ac6Provider{kind: kind})
		if !errors.Is(err, ErrSlotNotLanded) {
			t.Fatalf("slot %s: err = %v, want ErrSlotNotLanded", kind, err)
		}
	}
	if err := reg.RegisterProvider(&ac6Provider{kind: Kind("sidestep")}); err == nil {
		t.Fatal("an unknown slot was accepted")
	}
	if got := len(UnlandedSlots()); got != 2 {
		t.Fatalf("UnlandedSlots() = %d, want 2 (goja/mcp)", got)
	}
	// A manifest provider with nothing behind it is still a legal registration
	// (a plugin that declares no tools), and its entries default to the slot kind.
	if err := reg.RegisterProvider(&ac6Provider{kind: KindManifest}); err != nil {
		t.Fatalf("empty manifest slot: %v", err)
	}
}
