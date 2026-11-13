//go:build windows

package main

// The appearance wave's wiring teeth, the 65/413/422 shape: the boot seed
// reads the two [ball] keys the config screen left locked (ball.opacity_idle,
// ball.click_through), the settings write route lands all four fields through
// the REAL guarded setters, the two [ball] writes answer their measured tier
// (now when an applier is installed - the real product functions, not stubs -
// restart when this process hosts no ball), the two [panel] writes carry the
// honest Note about the per-create reader, and every refusal is the
// single-sourced sentence.
//
// The reverse criteria, per case: deleting the seed read
// (ball_idle_appearance_windows.go's c.Ball.* line) reddens the seed test on
// ok=false; deleting the withBallIdleAppearance hand-over or the Options
// literal reddens the reflection test on nil; deleting an ApplySetting case
// reddens the route test on "没有对应的写入腿"; deleting the tier/Note
// stamps reddens the receipt cases.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// wavePlantDir writes a minimal valid config.toml whose [ball] section
// carries a NON-DEFAULT opacity_idle and click_through (the values the seed
// tests read back), plus the [panel] rows the store tests need - the same
// section set ballStyleTestConfig65 keeps so validation sees a loadable
// file. 0.7 and false both differ from the schema defaults (0.35 / true), so
// a default leaking into any reading cannot fake the answer.
func wavePlantDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	body := "schema_version = 2\n\n[ball]\nsize = 56\nopacity_idle = 0.7\nclick_through = false\n" +
		"\n[app]\ntheme = \"dark\"\nlanguage = \"zh-CN\"\n" +
		"\n[voice]\nenabled = true\n\n[voice.asr]\nprovider = \"local-sherpa\"\nmodel = \"parakeet-zh\"\n" +
		"\n[panel]\nwidth = 700\nfont_size = 14\n" +
		"\n[fs]\nallowed_dirs = []\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

// ballOptsWave reads the two appearance-wave fields off the ball's unexported
// Options through reflect - the 413 ruling's instrument (reflect reads, never
// writes, and production carries no reader of these options back out of the
// ball), so the test cannot have fed the answer it reports.
func ballOptsWave(t *testing.T, b *ball.Ball) (*float64, bool) {
	t.Helper()
	if b == nil {
		t.Fatal("no ball handle: the reflection ruler has no subject")
	}
	f := reflect.ValueOf(b).Elem().FieldByName("opts").FieldByName("OpacityIdle")
	if !f.IsValid() || f.Kind() != reflect.Ptr {
		t.Fatalf("ball.Ball's opts.OpacityIdle is not a *float64 field (valid=%v kind=%v): the option the feed rides has moved", f.IsValid(), f.Kind())
	}
	var op *float64
	if !f.IsNil() {
		// Read through the pointer without Interface(): RO-flagged values
		// (read off an unexported field) allow Float() and forbid only
		// Interface/Set - the read-only rule the 413 instrument lives under.
		v := f.Elem().Float()
		op = &v
	}
	ct := reflect.ValueOf(b).Elem().FieldByName("opts").FieldByName("ClickThrough")
	if !ct.IsValid() || ct.Kind() != reflect.Bool {
		t.Fatalf("ball.Ball's opts.ClickThrough is not a bool field: the option the feed rides has moved")
	}
	return op, ct.Bool()
}

// TestWaveSeedReadsTheFile: the seed is the file's own answer for both keys
// - the mutation tooth is ok: delete the c.Ball.* read and ok=false reds the
// case, because no local default exists to fake it. The unreadable-file half
// names the carrier: nil opacity (the ball package's own idle rendering
// stands; 0 is a real configured value, so a sentinel float would be a lying
// option) and the pre-wave clickable idle ball.
func TestWaveSeedReadsTheFile(t *testing.T) {
	dir := wavePlantDir(t)
	op, ct, ok := ballIdleAppearanceSeed(dir)
	if !ok {
		t.Fatalf("ballIdleAppearanceSeed = ok=false for a loadable config - the read is load-bearing")
	}
	if op == nil || *op != 0.7 {
		t.Fatalf("seed opacity = %v, want the planted 0.7", op)
	}
	if ct {
		t.Fatalf("seed clickThrough = true, want the planted false")
	}
	// The seed adds nothing of its own: a fresh LoadFile of the same file
	// answers the same pair.
	c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil || c == nil {
		t.Fatalf("control LoadFile failed: %v", err)
	}
	if op == nil || *op != c.Ball.OpacityIdle || ct != c.Ball.ClickThrough {
		t.Fatalf("seed (%v, %v) disagrees with cfg.Ball (%v, %v): the seed answered something the file did not say",
			op, ct, c.Ball.OpacityIdle, c.Ball.ClickThrough)
	}
	// The unreadable half: nil carrier, false switch, ok=false - said by the
	// seed's Warn, not silently swallowed.
	empty := t.TempDir()
	op, ct, ok = ballIdleAppearanceSeed(empty)
	if ok || op != nil || ct {
		t.Fatalf("seed over a config-less dir = (%v, %v, %v), want (nil, false, false)", op, ct, ok)
	}
}

// TestWaveFeedReachesBallOptions is the positive criterion: the full
// production wiring (startResidentBall over a live ball, both values handed
// in with withBallIdleAppearance the way resident_windows.go hands them) puts
// the configured pair into ball.Options, and a reflection read of the ball's
// own opts answers cfg.Ball's values.
func TestWaveFeedReachesBallOptions(t *testing.T) {
	dir := wavePlantDir(t)
	op, ct, ok := ballIdleAppearanceSeed(dir)
	if !ok || op == nil || *op != 0.7 || ct {
		t.Fatalf("seed = (%v, %v, %v), want the planted (0.7, false, true): the fixture broke, not the feed", op, ct, ok)
	}
	rb := startResidentBallBounded(t, observe.NewRegistry(), nil, nil, nil, withBallIdleAppearance(op, ct))
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the feed has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	gotOp, gotCt := ballOptsWave(t, rb.b)
	if gotOp == nil || *gotOp != 0.7 {
		t.Errorf("ball.Options.OpacityIdle = %v, want the file's 0.7: the configured idle opacity did not reach the option it names (the pre-wave gap is back)", gotOp)
	}
	if gotCt {
		t.Errorf("ball.Options.ClickThrough = true, want the file's false: the configured click-through did not reach the option")
	}
}

// TestWaveUnfedBallKeepsTheStanding is the reverse criterion: with no
// withBallIdleAppearance hook (the pre-wave wiring), the same construction
// answers nil opacity (the package's own idle rendering stands - NOT 0.35,
// which would be a swallowed zero-value) and the pre-wave clickable idle
// ball - proving the DIFFERENTIATION the two live boots make.
func TestWaveUnfedBallKeepsTheStanding(t *testing.T) {
	rb := startResidentBallBounded(t, observe.NewRegistry(), nil, nil, nil)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the unfed construction has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	gotOp, gotCt := ballOptsWave(t, rb.b)
	if gotOp != nil {
		t.Errorf("unfed ball.Options.OpacityIdle = %v, want nil (the package's own idle rendering stands)", *gotOp)
	}
	if gotCt {
		t.Errorf("unfed ball.Options.ClickThrough = true, want false (the pre-wave clickable idle ball: no bit is ever set from a zero-value option)")
	}
}

// TestWaveApplySettingRoutesAllFour is the write route's teeth: all four
// fields land through the real setters, the [ball] tiers are the measured
// ones (restart with no applier, now with the real product functions
// installed), the [panel] receipts carry the per-create-reader Note, and
// every refusal names its field.
func TestWaveApplySettingRoutesAllFour(t *testing.T) {
	t.Cleanup(func() {
		// The applier seams are package globals; tests here must not leave
		// them installed for a later leg.
		ballIdleOpacityApplier = nil
		ballClickThroughApplier = nil
	})
	mgr, err := config.NewManager(filepath.Join(wavePlantDir(t), configFileName), nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	s := newConfigStore(mgr, nil, nil, "test-wave")

	// With no appliers (a process hosting no ball): both [ball] writes land
	// on disk but answer restart, and the [panel] writes carry the Note.
	res, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallOpacityIdle, Value: "0.5", Actor: "test-wave"})
	if err != nil {
		t.Fatalf("ApplySetting(ball_opacity_idle): %v", err)
	}
	if res.Tier != panel.EffectiveRestart || len(res.Written) != 1 || res.Written[0] != "ball.opacity_idle" {
		t.Fatalf("no-applier receipt = (tier %v, written %v), want (restart, [ball.opacity_idle])", res.Tier, res.Written)
	}
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallClickThrough, Value: "true", Actor: "test-wave"})
	if err != nil {
		t.Fatalf("ApplySetting(ball_click_through): %v", err)
	}
	if res.Tier != panel.EffectiveRestart || len(res.Written) != 1 || res.Written[0] != "ball.click_through" {
		t.Fatalf("no-applier receipt = (tier %v, written %v), want (restart, [ball.click_through])", res.Tier, res.Written)
	}
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldPanelWidth, Value: "520", Actor: "test-wave"})
	if err != nil {
		t.Fatalf("ApplySetting(panel_width): %v", err)
	}
	if res.Tier != panel.EffectiveRestart || res.Note == "" || !strings.Contains(res.Note, "关窗再开") {
		t.Fatalf("panel_width receipt = (tier %v, note %q), want restart plus the per-create-reader Note", res.Tier, res.Note)
	}
	if v := mgr.PanelWidthOptions(); v.Width != 520 {
		t.Fatalf("PanelWidthOptions().Width = %d, want 520 - the write must be readable hot in this same process", v.Width)
	}
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldPanelFontSize, Value: "13", Actor: "test-wave"})
	if err != nil {
		t.Fatalf("ApplySetting(panel_font_size): %v", err)
	}
	if res.Tier != panel.EffectiveRestart || res.Note == "" || !strings.Contains(res.Note, "关窗再开") {
		t.Fatalf("panel_font_size receipt = (tier %v, note %q), want restart plus the Note", res.Tier, res.Note)
	}
	if v := mgr.PanelFontSizeOptions(); v.FontSize != 13 {
		t.Fatalf("PanelFontSizeOptions().FontSize = %d, want 13", v.FontSize)
	}

	// The real product functions as the appliers (Ball.SetIdleOpacity /
	// Ball.SetClickThrough are posted to an STA thread; a nil-ball handle
	// answers false, so the SEAM shape is proven here and the live-ball
	// behavior is the winlive suite's - the same split the 65 leg made).
	var seamGot []string
	ballIdleOpacityApplier = func(v float64) bool { seamGot = append(seamGot, "opacity"); return true }
	ballClickThroughApplier = func(on bool) bool { seamGot = append(seamGot, "click"); return true }
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallOpacityIdle, Value: "0.4", Actor: "test-wave"})
	if err != nil || res.Tier != panel.EffectiveNow {
		t.Fatalf("applier receipt = (%v, %v), want (now, nil)", res.Tier, err)
	}
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallClickThrough, Value: "false", Actor: "test-wave"})
	if err != nil || res.Tier != panel.EffectiveNow {
		t.Fatalf("applier receipt = (%v, %v), want (now, nil)", res.Tier, err)
	}
	if len(seamGot) != 2 || seamGot[0] != "opacity" || seamGot[1] != "click" {
		t.Fatalf("the appliers saw %v, want [opacity click] - a now tier must mean a real apply ran", seamGot)
	}

	// Refusals: each parse failure is refused before the setter, each value
	// the setter refuses comes back as the single-sourced sentence, and a
	// refusal moves no byte.
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallOpacityIdle, Value: "tall", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), `字段 "ball_opacity_idle" 要的是 0..1 的数字`) {
		t.Fatalf("a non-float opacity must be refused by the parse sentence, got: %v", err)
	}
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallOpacityIdle, Value: "1.7", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), "config.toml: ball.opacity_idle 1.7 out of range [0, 1]") {
		t.Fatalf("an out-of-band opacity must be refused by the validate sentence verbatim, got: %v", err)
	}
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallClickThrough, Value: "maybe", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), `字段 "ball_click_through" 要的是 true/false`) {
		t.Fatalf("a non-bool switch must be refused by the parse sentence, got: %v", err)
	}
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldPanelWidth, Value: "wide", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), `字段 "panel_width" 要的是正整数`) {
		t.Fatalf("a non-integer width must be refused by the parse sentence, got: %v", err)
	}
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldPanelWidth, Value: "0", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), "panel.width") {
		t.Fatalf("a zero width must be refused by the setter, got: %v", err)
	}
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldPanelFontSize, Value: "-2", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), "panel.font_size") {
		t.Fatalf("a negative font size must be refused by the setter, got: %v", err)
	}
	if v := mgr.PanelWidthOptions(); v.Width != 520 {
		t.Fatalf("after the refusals Width = %d, want the untouched 520 - a refusal moves no byte", v.Width)
	}
}

// TestWaveStoreRefusesAnUnlistedField is the door's own tooth for the wave:
// the field must be admitted BY the union (the roster), not by an absence of
// checking - a misspelled wave name stays refused.
func TestWaveStoreRefusesAnUnlistedField(t *testing.T) {
	mgr, err := config.NewManager(filepath.Join(wavePlantDir(t), configFileName), nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	s := newConfigStore(mgr, nil, nil, "test-wave")
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: "ball_opacity", Value: "0.5", Actor: "test-wave"}); err == nil ||
		!strings.Contains(err.Error(), "没有对应的写入腿") {
		t.Fatalf("a misspelled wave field must be refused by the leg's default branch (the door's roster check is one hop above), got: %v", err)
	}
}
