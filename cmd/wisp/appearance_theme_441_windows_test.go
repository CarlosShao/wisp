//go:build windows

package main

// 票 441 的球侧主题牙（AC#2）与互不串证明（AC#4）。
//
// 三层各钉一次：
//   1. 种子映射：config.toml 的 [app].theme → ball.Theme（light/dark 真映射，
//      auto 与读不了的文件 ok=false 不冒充）；
//   2. 装配点静钉：resident_windows.go 里 SetTheme 的生产调用者在
//      startResidentBall（ball.New）之前——时序就是"第一帧画新 palette"的
//      依据，钉子照 panel_appearance_wave_test.go 的 parse 款式做；
//   3. AC#4 互不串：只动 font_size 的两份文件 ⇒ 种子的 theme 答案逐位相同；
//      只动 theme 的两份文件 ⇒ PanelFontSizeOptions().FontSize 逐位相同；
//      以及真写半：走真门 SetPanelFontSize(15) 落盘后再跑种子 ⇒ theme 不动。
//      两枚牙各喂自己的键、互读不到对方的键——没有同一枚假阳性牙充数。

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
)

// theme441Plant writes a minimal valid config.toml carrying [app].theme and
// [panel].font_size - the same section set wavePlantDir keeps so validation
// sees a loadable file. The two knobs are the two teeth's ONLY dials.
func theme441Plant(t *testing.T, theme string, fontSize int) string {
	t.Helper()
	dir := t.TempDir()
	body := "schema_version = 2\n\n[ball]\nsize = 56\n" +
		"\n[app]\ntheme = \"" + theme + "\"\nlanguage = \"zh-CN\"\n" +
		"\n[voice]\nenabled = true\n\n[voice.asr]\nprovider = \"local-sherpa\"\nmodel = \"parakeet-zh\"\n" +
		"\n[panel]\nwidth = 640\nfont_size = " + strconv.Itoa(fontSize) + "\n" +
		"\n[fs]\nallowed_dirs = []\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

// Test441BallThemeSeedMapsTheConfigKey: the seed is the file's own answer for
// the two resolvable words, and it REFUSES to fake an answer for the rest -
// auto resolves nothing here (the Warn names it) and an unreadable file keeps
// the renderer's dark default standing. The mutation tooth is ok: delete the
// c.App.Theme read and the light case reds on ok=false.
func Test441BallThemeSeedMapsTheConfigKey(t *testing.T) {
	light, okLight := ballThemeSeed441(theme441Plant(t, "light", 12))
	if !okLight || light != ball.ThemeLight {
		t.Fatalf("seed(app.theme=light) = (%v, %v), want (ThemeLight, true)", light, okLight)
	}
	dark, okDark := ballThemeSeed441(theme441Plant(t, "dark", 12))
	if !okDark || dark != ball.ThemeDark {
		t.Fatalf("seed(app.theme=dark) = (%v, %v), want (ThemeDark, true)", dark, okDark)
	}
	if light == dark {
		t.Fatalf("the two themes mapped to the same Theme value (%v): the tooth cannot tell them apart", light)
	}
	if got, ok := ballThemeSeed441(theme441Plant(t, "auto", 12)); ok {
		t.Fatalf("seed(app.theme=auto) = ok=true (%v): this boot resolves no system answer, faking one is the lie", got)
	}
	if _, ok := ballThemeSeed441(t.TempDir()); ok {
		t.Fatalf("seed over a config-less dir = ok=true: an unreadable file must keep the named fallback, not invent an answer")
	}
}

// Test441ThemeAssemblySitePin: the production caller exists at the assembly
// root, reads the seed, and sits BEFORE startResidentBall - ball.New is where
// the first frame paints, so a call after it would switch a palette the
// renderer never reads without a repaint (the order is the claim).
func Test441ThemeAssemblySitePin(t *testing.T) {
	src, err := os.ReadFile("resident_windows.go")
	if err != nil {
		t.Fatalf("resident_windows.go: %v", err)
	}
	s := string(src)
	seedIdx := strings.Index(s, "ballThemeSeed441(rt.Layout.DataDir)")
	callIdx := strings.Index(s, "ball.SetTheme(ballTheme441)")
	ballIdx := strings.Index(s, "startResidentBall(rt.Registry")
	if seedIdx < 0 || callIdx < 0 || ballIdx < 0 {
		t.Fatalf("the assembly site moved: seedIdx=%d callIdx=%d ballIdx=%d - SetTheme's production caller is gone or renamed", seedIdx, callIdx, ballIdx)
	}
	if !(seedIdx < callIdx && callIdx < ballIdx) {
		t.Fatalf("ordering broken: seed(%d) -> SetTheme(%d) -> startResidentBall(%d) must be strictly ascending - the call must land before ball.New, not after", seedIdx, callIdx, ballIdx)
	}
}

// Test441TeethDoNotCross (AC#4): the font tooth and the theme tooth each read
// ONLY their own key, in both directions, plus one hot-write leg through the
// real guarded setter.
func Test441TeethDoNotCross(t *testing.T) {
	// Forward: two files differing ONLY in font_size - the theme seed's
	// answer must be bit-identical (the font tooth cannot move the theme).
	tDark12, okA := ballThemeSeed441(theme441Plant(t, "dark", 12))
	tDark15, okB := ballThemeSeed441(theme441Plant(t, "dark", 15))
	if !okA || !okB || tDark12 != tDark15 || tDark12 != ball.ThemeDark {
		t.Fatalf("changing panel.font_size moved the theme seed: (%v,%v) vs (%v,%v)", tDark12, okA, tDark15, okB)
	}
	// Reverse: two files differing ONLY in theme - the font read face must
	// answer the same planted px (the theme tooth cannot move the font).
	mgrDark, err := config.NewManager(filepath.Join(theme441Plant(t, "dark", 14), configFileName), nil)
	if err != nil {
		t.Fatalf("manager(dark): %v", err)
	}
	mgrLight, err := config.NewManager(filepath.Join(theme441Plant(t, "light", 14), configFileName), nil)
	if err != nil {
		t.Fatalf("manager(light): %v", err)
	}
	fontDark := panelFontSizeSectionFromConfig(mgrDark)
	fontLight := panelFontSizeSectionFromConfig(mgrLight)
	if fontDark.FontSize != 14 || fontLight.FontSize != 14 {
		t.Fatalf("changing app.theme moved the font read face: dark=%+v light=%+v, want both font_size=14", fontDark, fontLight)
	}
	// Hot-write leg: a real SetPanelFontSize lands on disk, then the theme
	// seed re-reads the SAME file - the theme must not have moved.
	dir := theme441Plant(t, "light", 12)
	before, okBefore := ballThemeSeed441(dir)
	if !okBefore || before != ball.ThemeLight {
		t.Fatalf("fixture broke: seed before the write = (%v, %v)", before, okBefore)
	}
	mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	written, err := mgr.SetPanelFontSize(15)
	if err != nil || len(written) != 1 || written[0] != "panel.font_size" {
		t.Fatalf("SetPanelFontSize(15) = (%v, %v), want ([panel.font_size], nil)", written, err)
	}
	if v := mgr.PanelFontSizeOptions().FontSize; v != 15 {
		t.Fatalf("PanelFontSizeOptions().FontSize = %d after the write, want 15 (the positive control)", v)
	}
	after, okAfter := ballThemeSeed441(dir)
	if !okAfter || after != ball.ThemeLight {
		t.Fatalf("the font write moved the theme seed: (%v, %v) after SetPanelFontSize(15), want (ThemeLight, true)", after, okAfter)
	}
}
