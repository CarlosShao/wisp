package ball

// Ticket 69 (registry A24-D4) — machine checks for the C21 design-token table.
//
// docs/evidence/s1/c21-native-tokens.md is the hand-maintained C21 table. Since
// ticket 12 AC#7 it carries 61 rows nobody tested (25 geometry/motion rows +
// 36 liquid-glass look colours): TestTokenGoldenValues pins 20 palette colours
// and TestNoHardcodedColorsInBallPackage only forbids literals in the drawing
// code. Table and code could therefore drift in opposite directions with CI
// green — and one row (tokens.go CountdownFontPx) had already drifted to zero
// consumers without anything being able to see it.
//
// These tests close that gap in BOTH directions, inside the table's own
// declared scope (the "范围" block ticket 12 wrote, plus A24-D6):
//
//	table -> code   every row names a token that exists, with the value stated
//	code -> table   every Palette field, every look colour and every exported
//	                tokens.go geometry constant is declared by exactly one row;
//	                since ticket 74 the geometry rows may also name the constants
//	                that live beside them (hit.go, liquid.go), and those values
//	                are pinned by c21MotionGolden
//
// What is deliberately NOT asserted here, and why:
//   - the table's `tokens.css 值` column against design/assets/tokens.css. That
//     is ticket 12's manual three-way cross-check (A24-D6 leaves the 80
//     panel-only CSS declarations out of scope); converting it into a machine
//     check is a separate call, not this ticket's.
//   - any *default* ball value. Ticket 68 AC#2 owns flipping the prototype
//     default, and SPEC-08 §2 is frozen. Where table and code disagree on a
//     value, the row is reported for the owner (A24-D1..D6 precedent), never
//     "fixed" by one side.
//   - whether an unconsumed token is legal. It is not this test's call: the
//     zero-consumer report prints the list (AC: report first, not fatal) and
//     the ticket carries it to the owner.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const c21TablePath = "docs/evidence/s1/c21-native-tokens.md"

// ---------------------------------------------------------------- table model

// c21ColourRow is one row of the three colour tables (dark, light, look).
type c21ColourRow struct {
	where   string // "docs/evidence/s1/c21-native-tokens.md:41"
	section string // "dark" | "light" | "look"
	cssVar  string // first column, verbatim (--accent, or 无 for the look rows)
	goRef   string // "Palette.Accent" / "looks[aurora].BlobA" (backticks off)
	claim   string // "hex(0x86C2B9,1)" - what the table says the code holds
	line    int
}

// c21GeomRow is one row of the geometry/motion table.
type c21GeomRow struct {
	where  string
	rule   string        // first column, for the log lines
	names  []string      // the row's "Go 常量" column
	claims []c21NumClaim // numbers stated in its "值" and "用途" columns
	body   string        // its "值" and "用途" columns, for string-token assertions
	line   int
}

// c21NumClaim is one number written into a geometry row, with the unit that
// stood next to it ("ms", "s", "px", "fps", "dpi", "%" or "" when bare).
type c21NumClaim struct {
	val    float64
	unit   string
	tagged bool
}

var (
	c21IdentRe  = regexp.MustCompile("`([A-Za-z_][A-Za-z0-9_]*)`")
	c21ExportRe = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	c21NumRe    = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*(ms|fps|dpi|px|s|%)?`)
	c21HexLitRe = regexp.MustCompile(`0[xX][0-9A-Fa-f]+`)
	c21RGBARe   = regexp.MustCompile(`^rgba\(\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9.]+)\s*\)$`)
	c21HexRe    = regexp.MustCompile(`^hex\(\s*(0[xX][0-9A-Fa-f]{6})\s*,\s*([0-9.]+)\s*\)$`)
	c21LookRef  = regexp.MustCompile(`^looks\[([^\]]+)\]\.([A-Za-z0-9]+)$`)
)

// c21RepoRoot locates the checkout from this file's own path. A missing table is
// fatal on purpose: an unreadable table must never be able to pass as "checked".
func c21RepoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(c21TablePath))); err != nil {
		t.Fatalf("%s is not reachable from %s: %v - this check must never skip", c21TablePath, root, err)
	}
	return root
}

func c21PkgDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Dir(thisFile)
}

// c21ExemptRow is one row of the table's 未接线豁免 section: the Go references it
// attributes and the one-line reason ticket 74 family 3 demands per row. The
// reasons live HERE (in the document the check reads) rather than in a Go map so
// the exemption list cannot become a second truth that drifts from the table.
type c21ExemptRow struct {
	where  string
	tokens []string // "Palette.FgPrimary" / "looks[*].Deep" / "CountdownFontPx"
	reason string
	line   int
}

// c21ExemptRefRe picks the backticked token references out of an exemption row's
// first column. Dot- and bracket-qualified forms are allowed because that is how
// the two colour tables name their fields.
var c21ExemptRefRe = regexp.MustCompile("`(looks\\[\\*\\]\\.[A-Za-z0-9]+|[A-Z][A-Za-z0-9]*(\\.[A-Za-z0-9]+)?)`")

// c21ParseTable reads the markdown rows of the four token tables. Rows are only
// collected under the four known section headings, so the header block's own
// comparison table and the 审计 section cannot leak in.
func c21ParseTable(t *testing.T, root string) ([]c21ColourRow, []c21GeomRow, []c21ExemptRow) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(c21TablePath)))
	if err != nil {
		t.Fatalf("read %s: %v", c21TablePath, err)
	}
	var (
		colour  []c21ColourRow
		geom    []c21GeomRow
		exempt  []c21ExemptRow
		section string
	)
	for i, raw := range strings.Split(string(data), "\n") {
		line, lineNo := strings.TrimRight(raw, "\r"), i+1
		if strings.HasPrefix(line, "#") {
			switch {
			case strings.Contains(line, "明色系"):
				section = "dark"
			case strings.Contains(line, "亮色系"):
				section = "light"
			case strings.Contains(line, "色板"):
				section = "look"
			case strings.Contains(line, "几何"):
				section = "geom"
			case strings.Contains(line, "豁免"):
				section = "unconsumed"
			default:
				section = ""
			}
			continue
		}
		if section == "" || !strings.HasPrefix(line, "|") {
			continue
		}
		cells := c21Cells(line)
		if c21IsSeparator(cells) {
			continue
		}
		where := fmt.Sprintf("%s:%d", c21TablePath, lineNo)
		switch section {
		case "dark", "light", "look":
			if len(cells) != 4 || strings.HasPrefix(cells[0], "CSS 变量") {
				continue
			}
			colour = append(colour, c21ColourRow{
				where: where, section: section, cssVar: cells[0], line: lineNo,
				goRef: c21TrimTicks(cells[2]), claim: c21TrimTicks(cells[3]),
			})
		case "geom":
			if len(cells) != 4 || strings.HasPrefix(cells[0], "CSS 变量") {
				continue
			}
			row := c21GeomRow{where: where, line: lineNo, rule: cells[0], body: cells[1] + " " + cells[3]}
			for _, m := range c21IdentRe.FindAllStringSubmatch(cells[2], -1) {
				if c21ExportRe.MatchString(m[1]) {
					row.names = append(row.names, m[1])
				}
			}
			row.claims = c21NumClaims(row.body)
			geom = append(geom, row)
		case "unconsumed":
			if len(cells) != 2 || strings.HasPrefix(cells[0], "Go 引用") {
				continue
			}
			row := c21ExemptRow{where: where, line: lineNo, reason: cells[1]}
			for _, m := range c21ExemptRefRe.FindAllStringSubmatch(cells[0], -1) {
				row.tokens = append(row.tokens, m[1])
			}
			if len(row.tokens) == 0 {
				t.Errorf("%s: an exemption row whose first column names no Go reference attributes nothing: %q", where, cells[0])
				continue
			}
			exempt = append(exempt, row)
		}
	}
	if len(colour) == 0 || len(geom) == 0 {
		t.Fatalf("%s parsed to %d colour rows and %d geometry rows: the document shape changed and every check below would pass vacuously",
			c21TablePath, len(colour), len(geom))
	}
	return colour, geom, exempt
}

// c21Cells splits a markdown table row into its trimmed cells.
func c21Cells(line string) []string {
	parts := strings.Split(line, "|")
	if len(parts) > 0 && strings.TrimSpace(parts[0]) == "" {
		parts = parts[1:]
	}
	if len(parts) > 0 && strings.TrimSpace(parts[len(parts)-1]) == "" {
		parts = parts[:len(parts)-1]
	}
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.TrimSpace(p)
	}
	return out
}

func c21IsSeparator(cells []string) bool {
	for _, c := range cells {
		if c == "" || strings.Trim(c, "-: ") != "" {
			return false
		}
	}
	return true
}

func c21TrimTicks(s string) string {
	return strings.Trim(strings.TrimSpace(s), "`")
}

// c21NumClaims extracts every number with its unit from row text. Hex literals
// are removed first so their digits cannot be read as values, and "dpi" is kept
// as its own unit so a screen resolution can never satisfy a px token.
func c21NumClaims(text string) []c21NumClaim {
	low := c21HexLitRe.ReplaceAllString(strings.ToLower(text), "")
	var out []c21NumClaim
	for _, m := range c21NumRe.FindAllStringSubmatch(low, -1) {
		v, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			continue
		}
		out = append(out, c21NumClaim{val: v, unit: m[2], tagged: m[2] != ""})
	}
	return out
}

func c21ClaimsText(claims []c21NumClaim) string {
	if len(claims) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(claims))
	for _, c := range claims {
		parts = append(parts, fmt.Sprintf("%g%s", c.val, c.unit))
	}
	return strings.Join(parts, ", ")
}

// c21TableNames is every Go name the table declares, keyed by the reference form
// used in each table's Go column.
func c21TableNames(colour []c21ColourRow, geom []c21GeomRow) map[string]string {
	out := map[string]string{}
	for _, r := range colour {
		out[r.goRef] = r.where
	}
	for _, r := range geom {
		for _, n := range r.names {
			out[n] = r.where
		}
	}
	return out
}

// ------------------------------------------------------------- the code side

// c21ConstInfo is one exported package-level constant read off tokens.go.
type c21ConstInfo struct {
	name string
	kind string // "num" | "string" | "derived" (constant expression) | "other"
	file string // which file of the package declares it
}

// c21TokensGoConsts enumerates the exported constants tokens.go actually declares
// straight from its syntax tree, so a constant added without a table row cannot
// hide.
func c21TokensGoConsts(t *testing.T, dir string) map[string]c21ConstInfo {
	t.Helper()
	out := map[string]c21ConstInfo{}
	for _, c := range c21ExportedConsts(t, filepath.Join(dir, "tokens.go")) {
		if c.kind == "enum" {
			continue // an iota enumeration member is not a design token value
		}
		out[c.name] = c
	}
	if len(out) < 50 {
		t.Fatalf("tokens.go enumerated only %d exported constants: the enumeration is broken, not the file", len(out))
	}
	return out
}

// c21PackageConsts enumerates the same surface across the WHOLE ball package,
// i.e. every non-test file next to tokens.go. The geometry table names constants
// that live beside the code they drive (hit.go's click margins, liquid.go's spin
// group) since ticket 74, and a table row naming a constant nobody declares must
// fail even when that constant is not in tokens.go. tokens.go entries win, so the
// stricter inventory (c21GeometryGolden) keeps its meaning.
func c21PackageConsts(t *testing.T, dir string) map[string]c21ConstInfo {
	t.Helper()
	out := map[string]c21ConstInfo{}
	for _, c := range c21TokensGoConsts(t, dir) {
		out[c.name] = c
	}
	for _, f := range c21DrawingFiles(t, dir) {
		for _, c := range c21ExportedConsts(t, f) {
			if c.kind == "enum" || c.kind == "other" {
				continue // an enumeration member is not a design token value
			}
			if _, taken := out[c.name]; taken {
				continue // tokens.go owns it
			}
			out[c.name] = c
		}
	}
	return out
}

// c21DrawingFiles is the package's non-test, non-tokens.go Go source.
func c21DrawingFiles(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil || len(files) == 0 {
		t.Fatalf("glob the ball package: %v (%d files)", err, len(files))
	}
	var out []string
	for _, f := range files {
		base := filepath.Base(f)
		if base == "tokens.go" || strings.HasSuffix(base, "_test.go") {
			continue
		}
		out = append(out, f)
	}
	if len(out) == 0 {
		t.Fatal("the ball package has no non-test source files: every check below would be vacuous")
	}
	return out
}

// c21ExportedConsts reads one file's exported package-level constants. A
// parenthesized block inherits the previous spec's type, which is how
// `const ( ThemeDark Theme = iota; ThemeLight )` makes ThemeLight a Theme: that
// type is not a plain numeric/string type, so both members are left out.
func c21ExportedConsts(t *testing.T, path string) []c21ConstInfo {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", filepath.Base(path), err)
	}
	var out []c21ConstInfo
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		var inherited ast.Expr
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			typ := vs.Type
			if typ == nil {
				typ = inherited
			} else {
				inherited = vs.Type
			}
			if typ != nil && !c21PlainValueType(typ) {
				continue
			}
			for i, id := range vs.Names {
				if !id.IsExported() {
					continue
				}
				var expr ast.Expr
				if i < len(vs.Values) {
					expr = vs.Values[i]
				}
				out = append(out, c21ConstInfo{name: id.Name, kind: c21ExprKind(expr), file: filepath.Base(path)})
			}
		}
	}
	return out
}

func c21PlainValueType(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	switch id.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32",
		"float32", "float64", "string", "rune", "byte":
		return true
	}
	return false
}

func c21ExprKind(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT, token.FLOAT:
			return "num"
		case token.STRING:
			return "string"
		}
	case *ast.BinaryExpr, *ast.ParenExpr, *ast.UnaryExpr:
		return "derived"
	case *ast.Ident:
		if e.Name == "iota" {
			return "enum"
		}
	}
	return "other"
}

func c21SortedNames(m map[string]c21ConstInfo) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// c21SortedKeys is the same for the float64 value inventories.
func c21SortedKeys(m map[string]float64) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// c21GeometryGolden is the test's inventory of the numeric token surface. Every
// value here is a reference to the constant itself, so it cannot drift away from
// tokens.go; what can drift is the NAME list, and
// TestC21GeometryRowsMatchCodeConstants fails if it stops matching the AST
// enumeration of tokens.go in either direction.
var c21GeometryGolden = map[string]float64{
	// Liquid blob placement + glass edge weights (ticket 62).
	"SleepRestRatio": SleepRestRatio,
	"LiquidRadiusA":  LiquidRadiusA,
	"LiquidRadiusB":  LiquidRadiusB,
	"LiquidRadiusC":  LiquidRadiusC,
	"LiquidOffsetA":  LiquidOffsetA,
	"LiquidOffsetB":  LiquidOffsetB,
	"LiquidOffsetC":  LiquidOffsetC,
	"GlassRimPx":     GlassRimPx,
	"GlassLipPx":     GlassLipPx,
	"GlassCaustic":   GlassCaustic,
	"BorderRingPx":   BorderRingPx,
	// SwimLevelGain / SpinLevelGain are gone (ticket 74 family 2): declared,
	// tabled and read by nothing, while the gains that move the liquid lived in
	// renderer_windows.go as bare literals. The two below are those literals
	// lifted into tokens, and liquid.go's spin group is tabled next to them.
	"LiquidGatherPerLevel": LiquidGatherPerLevel,
	"SummonFlowSpread":     SummonFlowSpread,
	"BorderOpenMs":         BorderOpenMs,
	"BorderCloseMs":        BorderCloseMs,
	"SummonFlowMs":         SummonFlowMs,
	"DockAnimMs":           DockAnimMs,
	"DockOverlapFrac":      DockOverlapFrac,
	"DockTriggerPx":        DockTriggerPx,

	// Geometry + motion tokens (SPEC-08 §2 / tokens.css).
	"BallSizeDefaultPx":    BallSizeDefaultPx,
	"BallSizeSmallPx":      BallSizeSmallPx,
	"BallSizeLargePx":      BallSizeLargePx,
	"BallSizeMinPx":        BallSizeMinPx,
	"BallSizeMaxPx":        BallSizeMaxPx,
	"SleepingDotPx":        SleepingDotPx,
	"SleepOpacity":         SleepOpacity,
	"SleepingRestMinPx":    SleepingRestMinPx,
	"RestSettledOpacity":   RestSettledOpacity,
	"DurFastMs":            DurFastMs,
	"DurBaseMs":            DurBaseMs,
	"DurSlowMs":            DurSlowMs,
	"WarmBreathPeriodMs":   WarmBreathPeriodMs,
	"SpeakingBreathMs":     SpeakingBreathMs,
	"ListeningWaveMs":      ListeningWaveMs,
	"ThinkingSweepMs":      ThinkingSweepMs,
	"ConfirmingPulseMs":    ConfirmingPulseMs,
	"FirstRunGuidePulseMs": FirstRunGuidePulseMs,
	"SettlingFadeMs":       SettlingFadeMs,
	"WarmOpacityLow":       WarmOpacityLow,
	"WarmOpacityHigh":      WarmOpacityHigh,
	"MaxAnimFPS":           MaxAnimFPS,
	"MinFrameMs":           MinFrameMs,
	"WarmBreathFPS":        WarmBreathFPS,
	"SlashStrokePx":        SlashStrokePx,
	"IconStrokePx":         IconStrokePx,
	"RingStrokePx":         RingStrokePx,
	"ConvRingStrokePx":     ConvRingStrokePx,
	"ThinkingBandPx":       ThinkingBandPx,
	"BadgeDiameterPx":      BadgeDiameterPx,
	"BadgeFontPx":          BadgeFontPx,
	"BadgeBorderPx":        BadgeBorderPx,
	"QueueDotPx":           QueueDotPx,
	"FontSizeMonoPx":       FontSizeMonoPx,
	"FontSizeMicroPx":      FontSizeMicroPx,
	"CountdownFontPx":      CountdownFontPx,
}

// c21MotionGolden pins the values of the geometry/motion constants that the C21
// table declares from OUTSIDE tokens.go (hit.go's click margins, liquid.go's
// spin and envelope group). Ticket 74 family 1: the table's scope stopped at
// tokens.go, so 11 constants that move pixels every frame were neither tabled nor
// value-checked - they were only "exempt", which is a note, not a contract. Same
// rule as the map above: every entry is a reference to the constant itself, so
// the value cannot drift away from the code; what can drift is the NAME list, and
// TestC21GeometryRowsMatchCodeConstants fails if it stops matching the table in
// either direction.
var c21MotionGolden = map[string]float64{
	// hit.go: window sizing + the clickable circle.
	"RingMarginPx":     RingMarginPx,
	"ClickTolerancePx": ClickTolerancePx,

	// liquid.go: what the audio envelope does to the liquid.
	"SpinBaseRadPerS":   SpinBaseRadPerS,
	"SpinLevelRadPerS":  SpinLevelRadPerS,
	"SpinSummonRadPerS": SpinSummonRadPerS,
	"SilenceLevelGate":  SilenceLevelGate,
	"SilenceHoldMs":     SilenceHoldMs,
	"LevelAttackTauMs":  LevelAttackTauMs,
	"LevelReleaseTauMs": LevelReleaseTauMs,
	"MotionEpsilon":     MotionEpsilon,
	"FrameIntervalMs":   FrameIntervalMs,
}

// c21TabledValue looks a tabled constant's value up in either inventory.
func c21TabledValue(name string) (float64, bool) {
	if v, ok := c21GeometryGolden[name]; ok {
		return v, true
	}
	v, ok := c21MotionGolden[name]
	return v, ok
}

// c21StringTokens are the non-numeric tokens of the same surface. The table names
// FontFamily's CSS source (--font-sans) but states no value for it, so there is
// nothing numeric to compare: the name is still checked in both directions, and
// this list is what keeps that one gap visible instead of silent.
var c21StringTokens = map[string]string{
	"FontFamily": FontFamily,
}

// c21OutTableExempt lists the exported constants the ball package keeps OUTSIDE
// tokens.go, which therefore sit outside the table's declared scope ("本表只覆盖
// internal/ball/tokens.go 的 40 个 Palette 字段 + 全部导出的几何/动效常量", A24-D6).
// Each entry carries the reason it is not a C21 token. A new exported constant in
// one of these files that is neither tabled nor listed here fails the test, which
// is the ticket's "要么进表要么显式豁免并写明理由".
//
// Ticket 74 emptied the two GLOW families out of this map: the 9 liquid.go motion
// constants (D8) and the 2 hit.go boundary constants (D7) are now rows of the
// table with their real units and pinned values, which is what A24-D4 asked for.
// What is left genuinely is not a design token.
var c21OutTableExempt = map[string]string{
	// hit.go: hit-test result codes. They are Win32 protocol, not ball geometry:
	// the numbers come from the WM_NCHITTEST contract, and changing them changes
	// no pixel.
	"HTClient":      "Win32 hit-test return value, not a design token",
	"HTTransparent": "Win32 hit-test return value, not a design token",
	"HTNowhere":     "Win32 hit-test return value, not a design token",

	// hotkey_windows.go: a keybinding string, not geometry or colour.
	"AltSummonSpace": "Default hotkey text (SPEC-04 surface), not a C21 token",
}

// c21StructColour reads one Color field by name.
func c21StructColour(holder any, field string) (Color, bool) {
	v := reflect.ValueOf(holder).FieldByName(field)
	if !v.IsValid() || v.Type() != reflect.TypeOf(Color{}) {
		return Color{}, false
	}
	c, ok := v.Interface().(Color)
	return c, ok
}

func c21PaletteOf(section string) Palette {
	if section == "light" {
		return LightPalette()
	}
	return DarkPalette()
}

func c21LookByNameExact(name string) (LiquidLook, bool) {
	for _, l := range looks {
		if l.Name == name {
			return l, true
		}
	}
	return LiquidLook{}, false
}

func c21ColourFieldNames(t *testing.T, typ reflect.Type) []string {
	t.Helper()
	var out []string
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type != reflect.TypeOf(Color{}) {
			continue
		}
		out = append(out, f.Name)
	}
	return out
}

// ----------------------------------------------------- check 1: colour rows

// TestC21TableColourRowsMatchCode walks the C21 table's three colour sections and
// asserts, row by row, that the token exists and holds exactly the value the
// table writes - then walks the other way and asserts that no Palette field or
// look colour exists without a row. That covers all 114 colour rows (40 dark +
// 38 light + 36 look), against the 20 TestTokenGoldenValues kept for the CSS
// spot check.
func TestC21TableColourRowsMatchCode(t *testing.T) {
	root := c21RepoRoot(t)
	colour, geom, _ := c21ParseTable(t, root)
	if len(geom) == 0 {
		t.Fatal("geometry table empty")
	}

	var dark, light, look []c21ColourRow
	for _, r := range colour {
		switch r.section {
		case "dark":
			dark = append(dark, r)
		case "light":
			light = append(light, r)
		case "look":
			look = append(look, r)
		}
	}
	// These are the counts the table itself declares (40 dark + 38 light rows,
	// and 4 looks x 9 colour fields). Deleting a row narrows the contract, so it
	// stops here instead of quietly checking less; adding rows needs no edit.
	for _, c := range []struct {
		name string
		got  int
		want int
	}{
		{"dark (:root)", len(dark), 40},
		{"light ([data-theme=light])", len(light), 38},
		{"look", len(look), len(looks) * 9},
	} {
		if c.got < c.want {
			t.Errorf("the %s colour table declares %d rows, want at least %d: a row that vanished takes its machine check with it (%s)",
				c.name, c.got, c.want, c21TablePath)
		}
	}

	// (a) table -> code: same name, same value.
	seen := make(map[string]c21ColourRow, len(colour))
	verified := 0
	for _, r := range colour {
		key := r.section + "|" + r.goRef
		if prev, dup := seen[key]; dup {
			t.Errorf("%s: %s is declared twice in the %s table (first at %s): one token needs exactly one row, or the check cannot say which value wins",
				r.where, r.goRef, r.section, prev.where)
			continue
		}
		seen[key] = r

		var got Color
		switch r.section {
		case "dark", "light":
			if !strings.HasPrefix(r.goRef, "Palette.") {
				t.Errorf("%s: %q is not a Palette.<Field> reference - this row no longer matches the table's documented form", r.where, r.goRef)
				continue
			}
			field := strings.TrimPrefix(r.goRef, "Palette.")
			var ok bool
			if got, ok = c21StructColour(c21PaletteOf(r.section), field); !ok {
				t.Errorf("%s: the table declares Palette.%s, but internal/ball/tokens.go has no Color field of that name", r.where, field)
				continue
			}
		case "look":
			m := c21LookRef.FindStringSubmatch(r.goRef)
			if m == nil {
				t.Errorf("%s: %q is not a looks[<name>].<Field> reference - this row no longer matches the table's documented form", r.where, r.goRef)
				continue
			}
			l, ok := c21LookByNameExact(m[1])
			if !ok {
				t.Errorf("%s: the table declares looks[%s], but tokens.go ships no look of that name (looks are %s)",
					r.where, m[1], strings.Join(LookNames(), ", "))
				continue
			}
			if got, ok = c21StructColour(l, m[2]); !ok {
				t.Errorf("%s: LiquidLook has no Color field %q, so this row checks nothing", r.where, m[2])
				continue
			}
		default:
			t.Errorf("%s: colour row in unknown section %q", r.where, r.section)
			continue
		}

		want, err := c21EvalClaim(r.claim)
		if err != nil {
			t.Errorf("%s: cannot read the table's Go value %q: %v", r.where, r.claim, err)
			continue
		}
		if !colorEq(got, want) {
			t.Errorf("%s: %s drifted - the table states %s (%s) and internal/ball/tokens.go holds %s",
				r.where, r.goRef, r.claim, c21ColourText(want), c21ColourText(got))
			continue
		}
		verified++
	}

	// (b) code -> table: nothing in the colour surface may be undeclared.
	darkIdx := map[string]c21ColourRow{}
	for _, r := range dark {
		darkIdx[strings.TrimPrefix(r.goRef, "Palette.")] = r
	}
	lightIdx := map[string]bool{}
	for _, r := range light {
		lightIdx[strings.TrimPrefix(r.goRef, "Palette.")] = true
	}
	paletteType := reflect.TypeOf(Palette{})
	darkPal, lightPal := reflect.ValueOf(DarkPalette()), reflect.ValueOf(LightPalette())
	var cascade []string
	for i := 0; i < paletteType.NumField(); i++ {
		f := paletteType.Field(i)
		if f.Type != reflect.TypeOf(Color{}) {
			t.Errorf("Palette.%s is a %s, not a Color: the colour tables only cover Color fields, so this one would go unchecked - extend these tests deliberately", f.Name, f.Type)
			continue
		}
		if _, declared := darkIdx[f.Name]; !declared {
			t.Errorf("Palette.%s exists in tokens.go but the dark (:root) table declares no row for it (code -> table direction, A24-D4)", f.Name)
		}
		if lightIdx[f.Name] {
			continue
		}
		// Not in the light table: the table's prose promises the light block does
		// not redefine it, i.e. CSS cascade carries the :root value. That promise
		// is only true if the two palettes really agree.
		dv, _ := c21StructColour(darkPal.Field(i).Interface(), f.Name)
		lv, _ := c21StructColour(lightPal.Field(i).Interface(), f.Name)
		if !colorEq(dv, lv) {
			t.Errorf("Palette.%s has no light-table row, so cascade says it carries the :root value, but dark=%s light=%s",
				f.Name, c21ColourText(dv), c21ColourText(lv))
			continue
		}
		cascade = append(cascade, f.Name)
	}
	if len(cascade) > 0 {
		sort.Strings(cascade)
		t.Logf("light table reuses the :root value by cascade for %d field(s): %s (the table states this in prose above its light rows)",
			len(cascade), strings.Join(cascade, ", "))
	}

	lookIdx := map[string]bool{}
	for _, r := range look {
		if m := c21LookRef.FindStringSubmatch(r.goRef); m != nil {
			lookIdx[m[1]+"/"+m[2]] = true
		}
	}
	lookType := reflect.TypeOf(LiquidLook{})
	lookChecked := 0
	for _, l := range looks {
		for i := 0; i < lookType.NumField(); i++ {
			f := lookType.Field(i)
			if f.Type != reflect.TypeOf(Color{}) {
				continue
			}
			if !lookIdx[l.Name+"/"+f.Name] {
				t.Errorf("looks[%s].%s exists in tokens.go but the look table declares no row for it (code -> table direction, A24-D4)", l.Name, f.Name)
				continue
			}
			lookChecked++
		}
	}
	t.Logf("verified %d/%d colour rows against tokens.go (%d palette fields checked in code -> table, %d look colours, %d looks)",
		verified, len(colour), paletteType.NumField(), lookChecked, len(looks))
}

// c21EvalClaim evaluates the Go construction the table writes for a colour
// (rgba(0..255,0..255,0..255,alpha) or hex(0xRRGGBB,alpha)) with the very
// helpers tokens.go uses, so "equal" means equal channels.
func c21EvalClaim(claim string) (Color, error) {
	if m := c21RGBARe.FindStringSubmatch(claim); m != nil {
		var ch [3]uint8
		for i := 0; i < 3; i++ {
			v, err := strconv.ParseUint(m[i+1], 10, 8)
			if err != nil {
				return Color{}, fmt.Errorf("channel %q is not a 0..255 integer", m[i+1])
			}
			ch[i] = uint8(v)
		}
		a, err := strconv.ParseFloat(m[4], 32)
		if err != nil {
			return Color{}, fmt.Errorf("alpha %q is unreadable", m[4])
		}
		return rgba(ch[0], ch[1], ch[2], float32(a)), nil
	}
	if m := c21HexRe.FindStringSubmatch(claim); m != nil {
		v, err := strconv.ParseUint(strings.TrimPrefix(strings.TrimPrefix(m[1], "0x"), "0X"), 16, 32)
		if err != nil {
			return Color{}, fmt.Errorf("hex %q is unreadable", m[1])
		}
		a, err := strconv.ParseFloat(m[2], 32)
		if err != nil {
			return Color{}, fmt.Errorf("alpha %q is unreadable", m[2])
		}
		return hex(uint32(v), float32(a)), nil
	}
	return Color{}, fmt.Errorf("neither rgba(...) nor hex(0xRRGGBB, alpha)")
}

func c21ColourText(c Color) string {
	return fmt.Sprintf("rgba(%.0f,%.0f,%.0f,%.6f)", c.R*255, c.G*255, c.B*255, c.A)
}

// ------------------------------------------------- check 2: geometry rows

// c21WantedUnits is the unit a constant's own name promises. The table writes
// durations in seconds ("Speaking 呼吸 1.6s") while the constants are ms, so "s"
// counts as ms; a bare number is what a ratio, opacity or gain is stated with.
func c21WantedUnits(name string) []string {
	switch {
	case strings.HasSuffix(name, "Ms"):
		return []string{"ms", "s"}
	case strings.HasSuffix(name, "Px"):
		return []string{"px"}
	case strings.HasSuffix(name, "FPS"):
		return []string{"fps"}
	default:
		return []string{""}
	}
}

// c21TakeClaim is the geometry row's value assertion. Ticket 69 shipped the match
// as containment ("does ANY number in this row equal the constant?"); this is the
// same question with the half that was missing: the number it matches is CONSUMED,
// so a second constant in the same row cannot be satisfied by the same written
// digit. Containment alone let a row mask a wrong value - measured before this
// change: setting tokens.go's DockTriggerPx to 160 kept the whole package green,
// because that row already writes 160ms for DockAnimMs. A row that legitimately
// states the same number for two constants must now write it twice, which is also
// the more honest table. A bare number that lacks the unit the constant's name
// promises is still a match, just a weaker one, and is reported as such.
func c21TakeClaim(val float64, name string, claims []c21NumClaim, taken []bool) (found, exactUnit bool) {
	want := c21WantedUnits(name)
	for i, c := range claims {
		if taken[i] || c.unit == "dpi" {
			continue
		}
		v := c.val
		if c.unit == "s" {
			v *= 1000
		}
		if math.Abs(v-val) > 1e-9 {
			continue
		}
		for _, u := range want {
			if c.unit == u {
				taken[i] = true
				return true, true
			}
		}
	}
	for i, c := range claims {
		if taken[i] || c.unit == "dpi" {
			continue
		}
		v := c.val
		if c.unit == "s" {
			v *= 1000
		}
		if math.Abs(v-val) > 1e-9 {
			continue
		}
		taken[i] = true
		return true, false
	}
	return false, false
}

// TestC21GeometryRowsMatchCodeConstants walks the geometry/motion table row by
// row: every constant a row names must exist in the ball package, and the value
// the package holds must be a number the row actually states. The reverse pass
// requires every exported tokens.go constant to be named by exactly one row.
func TestC21GeometryRowsMatchCodeConstants(t *testing.T) {
	root := c21RepoRoot(t)
	dir := filepath.Join(root, "internal", "ball")
	_, geom, _ := c21ParseTable(t, root)
	code := c21TokensGoConsts(t, dir)
	declared := c21PackageConsts(t, dir) // tokens.go + hit.go + liquid.go + ...

	// (0) the test's own inventory must be exactly tokens.go's constants, or the
	// rows below would be checked against a stale list.
	for _, name := range c21SortedNames(code) {
		info := code[name]
		switch info.kind {
		case "num", "derived":
			if _, ok := c21GeometryGolden[name]; !ok {
				t.Errorf("tokens.go declares the constant %s (%s) but c21GeometryGolden does not list it: add the value, add the table row, or say why it is not a token", name, info.kind)
			}
		case "string":
			if _, ok := c21StringTokens[name]; !ok {
				t.Errorf("tokens.go declares the string constant %s but c21StringTokens does not list it", name)
			}
		default:
			t.Errorf("tokens.go declares %s in a form this test does not understand (%s): extend the test deliberately, do not skip it", name, info.kind)
		}
	}
	for name := range c21GeometryGolden {
		if _, ok := code[name]; !ok {
			t.Errorf("c21GeometryGolden lists %s, which internal/ball/tokens.go no longer declares - it was renamed or deleted and the table row must follow", name)
		}
	}
	for name := range c21StringTokens {
		if _, ok := code[name]; !ok {
			t.Errorf("c21StringTokens lists %s, which internal/ball/tokens.go no longer declares", name)
		}
	}

	// (a) table -> code: names exist, one row per name.
	rowOf := map[string]c21GeomRow{}
	named := 0
	for _, r := range geom {
		if len(r.names) == 0 {
			t.Logf("%s: this row states a rule (%s) without naming a Go constant, so there is no value to compare", r.where, r.rule)
			continue
		}
		for _, n := range r.names {
			if prev, dup := rowOf[n]; dup {
				t.Errorf("%s: %s is declared by two table rows (%s and %s): one token, one home", r.where, n, prev.where, r.where)
			}
			rowOf[n] = r
			if _, ok := declared[n]; !ok {
				t.Errorf("%s: the table names %s, but no file in internal/ball declares that exported constant (tokens.go and its neighbours were both enumerated)", r.where, n)
			}
			named++
		}
	}

	// (a2) a row that names a constant outside tokens.go must pin its value too,
	// or the row states a number nothing is compared against.
	for _, name := range c21SortedKeys(c21MotionGolden) {
		if _, ok := code[name]; ok {
			t.Errorf("c21MotionGolden lists %s, which tokens.go declares: move it to c21GeometryGolden, the table's own scope", name)
			continue
		}
		if _, ok := declared[name]; !ok {
			t.Errorf("c21MotionGolden lists %s, which internal/ball no longer declares: it was renamed or deleted and the table row must follow", name)
			continue
		}
		if _, ok := rowOf[name]; !ok {
			t.Errorf("c21MotionGolden pins %s but no geometry-table row declares it: the value is checked against nothing", name)
		}
	}
	for _, r := range geom {
		for _, n := range r.names {
			if _, ok := code[n]; ok {
				continue // tokens.go's own surface: c21GeometryGolden is checked in full above
			}
			if _, ok := c21MotionGolden[n]; !ok {
				home := "declared nowhere in the package"
				if info, known := declared[n]; known {
					home = "declared in " + info.file
				}
				t.Errorf("%s: the table names %s, which is %s and therefore outside tokens.go's inventory, but c21MotionGolden pins no value for it - this row compares nothing",
					r.where, n, home)
			}
		}
	}

	// (b) code -> table.
	for _, name := range c21SortedNames(code) {
		if _, ok := rowOf[name]; !ok {
			t.Errorf("tokens.go exports %s but no geometry-table row declares it (code -> table direction, A24-D4)", name)
		}
	}

	// (c) values.
	var weak, unclaimed []string
	for _, r := range geom {
		if len(r.names) == 0 {
			continue
		}
		claimed := make([]bool, len(r.claims))
		for _, n := range r.names {
			if info, ok := declared[n]; ok && info.kind == "string" {
				want, pinned := c21StringTokens[n]
				if !pinned {
					t.Errorf("%s: the table names the string constant %s but c21StringTokens pins no value: this row compares nothing - pin it or say in the row that it carries no value", r.where, n)
					continue
				}
				// The value assertion FontFamily never had: the row quoted
				// --font-sans as its source but never wrote the family the code
				// actually asks DirectWrite for, so any string could sit in
				// tokens.go and the table still "matched".
				if !strings.Contains(r.body, want) {
					t.Errorf("%s: %s = %q, but this row never writes that string - a token value the table does not state cannot be audited: %s",
						r.where, n, want, c21TablePath)
				}
				continue
			}
			val, ok := c21TabledValue(n)
			if !ok {
				continue // already reported by (0)/(a)/(a2)
			}
			found, exact := c21TakeClaim(val, n, r.claims, claimed)
			if !found {
				t.Errorf("%s: %s = %s, but this row states no UNCLAIMED number equal to it (its numbers are: %s) - table and code drifted, or two constants are sharing one written digit; state the number once per constant",
					r.where, n, c21NumText(val), c21ClaimsText(r.claims))
				continue
			}
			if !exact {
				weak = append(weak, fmt.Sprintf("%s=%s @ %s", n, c21NumText(val), r.where))
			}
		}
		for i, c := range r.claims {
			if !claimed[i] {
				unclaimed = append(unclaimed, fmt.Sprintf("%g%s @ %s", c.val, c.unit, r.where))
			}
		}
	}
	// What the check still cannot say, out loud: a number written in prose that
	// belongs to no named constant (a SPEC clause number, "96 DPI", a commit hash)
	// is listed below and NOT judged, because surjectivity would turn every
	// sentence in the 用途 column into a token. What ticket 74 did close is the
	// other half - one written number can no longer satisfy two constants.
	if len(weak) > 0 {
		sort.Strings(weak)
		t.Logf("MATCHED WITHOUT THE PROMISED UNIT (%d): %s", len(weak), strings.Join(weak, ", "))
	}
	if len(unclaimed) > 0 {
		sort.Strings(unclaimed)
		t.Logf("NUMBERS NO NAMED CONSTANT CLAIMS (%d): %s", len(unclaimed), strings.Join(unclaimed, ", "))
	}
	t.Logf("checked %d constant declarations across %d geometry rows against %d exported tokens.go constants + %d tabled from hit.go/liquid.go",
		named, len(geom), len(code), len(c21MotionGolden))
}

func c21NumText(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// -------------------------------------- check 3: zero consumers must be judged
//
// Ticket 69 shipped this as a report ("an unconsumed token is not this test's
// call"). Ticket 74 family 3 closes that: a report nobody reads is exactly how
// SwimLevelGain managed to be documented, machine-checked and wired to nothing at
// the same time. So now every zero-consumer token needs a row of the table's own
// 未接线豁免 section saying why it is legal, and the section is compared in BOTH
// directions (a missing row and a stale row are both red). The reasons live in
// the document, not in a Go map, so this cannot become a second drifting copy of
// the table it audits.

// c21MinExemptReason is the shortest reason this check accepts. It exists to keep
// the placeholder answers out: the shortest row in the table today is a full
// sentence naming a substitute path or an owning ticket, and "暂未使用" is not.
const c21MinExemptReason = 20

func TestC21TokenConsumerReport(t *testing.T) {
	dir := c21PkgDir(t)
	_, geom, exempt := c21ParseTable(t, c21RepoRoot(t))
	code := c21TokensGoConsts(t, dir)

	idents := map[string]int{} // plain identifiers: constants in use
	fields := map[string]int{} // selector names: Palette / LiquidLook fields in use
	scanned := 0
	for _, f := range c21DrawingFiles(t, dir) {
		n := c21CollectRefs(t, f, idents, fields)
		if n == 0 {
			t.Fatalf("%s contributed no identifiers: the scan is broken, not the file", filepath.Base(f))
		}
		scanned++
	}
	if scanned == 0 {
		t.Fatal("no drawing files were scanned: the report would be vacuous")
	}

	tokens := make([]string, 0, len(code)+len(c21MotionGolden)+49)
	for _, n := range c21SortedNames(code) {
		tokens = append(tokens, n)
	}
	// Ticket 74: the constants the table names from hit.go and liquid.go are held
	// to the same rule as tokens.go's - tabled means somebody reads it, or the
	// table says why not.
	for _, n := range c21SortedKeys(c21MotionGolden) {
		tokens = append(tokens, n)
	}
	for _, n := range c21ColourFieldNames(t, reflect.TypeOf(Palette{})) {
		tokens = append(tokens, "Palette."+n)
	}
	for _, n := range c21ColourFieldNames(t, reflect.TypeOf(LiquidLook{})) {
		tokens = append(tokens, "looks[*]."+n) // per-field: one look is live at a time
	}

	var zero, live []string
	for _, tok := range tokens {
		n := tok
		if i := strings.LastIndex(tok, "."); i >= 0 {
			n = tok[i+1:]
			if fields[n] == 0 {
				zero = append(zero, tok)
			} else {
				live = append(live, tok)
			}
			continue
		}
		if idents[n] == 0 {
			zero = append(zero, tok)
		} else {
			live = append(live, tok)
		}
	}
	// Anti-vacuity: if the identifier scan matched nothing, every token would land
	// in `zero` and the report would look like a full pass.
	if len(live) == 0 {
		t.Fatalf("none of the %d tokens matched a reference in the %d drawing files scanned - the scan is broken, not the code", len(tokens), scanned)
	}
	sort.Strings(zero)
	t.Logf("ZERO-CONSUMER TOKENS (%d of %d checked; %d are referenced by %d non-test files):",
		len(zero), len(tokens), len(live), scanned)
	for _, tok := range zero {
		line := "  UNCONSUMED " + tok
		if r, ok := rowOfTableToken(geom, tok); ok {
			line += " (declared at " + r + ")"
		}
		t.Log(line)
	}

	// (a) the exemption rows, attributed.
	attributed := map[string]string{}
	for _, r := range exempt {
		if len([]rune(r.reason)) < c21MinExemptReason {
			t.Errorf("%s: the exemption reason is %d runes (%q) - ticket 74 wants a sentence naming the substitute path or the owning ticket, not a placeholder",
				r.where, len([]rune(r.reason)), r.reason)
		}
		if !strings.ContainsAny(r.reason, "`0123456789") {
			t.Errorf("%s: the exemption reason names no code path and no number (%q) - point it at what draws this instead, or at the ticket that will", r.where, r.reason)
		}
		for _, tok := range r.tokens {
			if prev, dup := attributed[tok]; dup {
				t.Errorf("%s: %s is exempted twice (%s and %s): one token, one home", r.where, tok, prev, r.where)
				continue
			}
			attributed[tok] = r.where
		}
	}

	// (b) zero-consumer <-> exemption: both directions.
	zeroSet := make(map[string]bool, len(zero))
	for _, tok := range zero {
		zeroSet[tok] = true
		if _, ok := attributed[tok]; !ok {
			t.Errorf("%s is tabled but nothing in the ball package reads it, and %s's 未接线豁免 table has no row for it: wire it, delete it with its replacement in one commit, or write the one-line reason (ticket 74 family 3, A33)",
				tok, c21TablePath)
		}
	}
	for _, r := range exempt {
		for _, tok := range r.tokens {
			if zeroSet[tok] {
				continue
			}
			if _, inTable := rowOfTableToken(geom, tok); !inTable && !strings.Contains(tok, ".") {
				t.Errorf("%s: the exemption row names %s, which the geometry table declares nowhere - the row audits nothing, point it at a tabled token or delete it", r.where, tok)
			}
			t.Errorf("%s: the exemption row for %s is stale - something in the package reads it again, so it is no longer documented-but-unconsumed; delete the row", r.where, tok)
		}
	}
	t.Logf("未接线豁免: %d row(s) attribute all %d zero-consumer token(s); %d/%d checked tokens are live",
		len(exempt), len(zero), len(live), len(tokens))
}

// rowOfTableToken finds the geometry row that declares a token, so the report
// points at the table line an owner would have to rule on.
func rowOfTableToken(geom []c21GeomRow, tok string) (string, bool) {
	name := tok
	if i := strings.LastIndex(tok, "."); i >= 0 {
		name = tok[i+1:]
	}
	for _, r := range geom {
		for _, n := range r.names {
			if n == name {
				return r.where, true
			}
		}
	}
	return "", false
}

// c21CollectRefs records every identifier and every selector name one file USES.
// Comments are excluded for free because the AST does not carry them: a token
// mentioned only in prose is not a consumer (that is how SleepRestRatio's
// comment-only mentions must not read as usage).
//
// The identifier's OWN declaration is not a use either. tokens.go is skipped by
// the caller so the point is moot there, but since ticket 74 the scan also covers
// hit.go and liquid.go, where the declarations and the uses share a file: without
// this, a constant nobody ever read would still count as consumed by its own
// `const` line, which is precisely the empty check A33 is about.
func c21CollectRefs(t *testing.T, path string, idents, fields map[string]int) int {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", filepath.Base(path), err)
	}
	declared := map[*ast.Ident]bool{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			if vs, ok := spec.(*ast.ValueSpec); ok {
				for _, id := range vs.Names {
					declared[id] = true
				}
			}
		}
	}
	n := 0
	ast.Inspect(file, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.SelectorExpr:
			if x.Sel != nil {
				fields[x.Sel.Name]++
				n++
			}
		case *ast.Ident:
			if declared[x] {
				return true // this file declaring it is not it being used
			}
			idents[x.Name]++
			n++
		}
		return true
	})
	return n
}

// -------------------------------------- check 4: no unlisted code-side tokens

// TestC21CodeTokensAreTabledOrExempt is the other half of the reverse direction:
// an exported constant anywhere in the ball package is either declared by the C21
// table or exempt with a written reason. tokens.go's own constants are held to the
// stricter rule by TestC21GeometryRowsMatchCodeConstants (the table's declared
// scope), so this test covers what lives next to it.
func TestC21CodeTokensAreTabledOrExempt(t *testing.T) {
	root := c21RepoRoot(t)
	dir := filepath.Join(root, "internal", "ball")
	colour, geom, _ := c21ParseTable(t, root)
	tabled := c21TableNames(colour, geom)

	seen := map[string]bool{}
	for _, info := range c21TokensGoConsts(t, dir) {
		seen[info.name] = true
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil || len(files) == 0 {
		t.Fatalf("glob the ball package: %v", err)
	}
	for _, f := range files {
		base := filepath.Base(f)
		if base == "tokens.go" || strings.HasSuffix(base, "_test.go") {
			continue
		}
		for _, c := range c21FileConsts(t, f) {
			if seen[c.name] {
				continue // tokens.go owns it; the geometry test judges it
			}
			if _, ok := tabled[c.name]; ok {
				continue // tabled even though it lives elsewhere
			}
			if reason, ok := c21OutTableExempt[c.name]; ok {
				t.Logf("EXEMPT %s (%s) is outside the table on purpose: %s", c.name, base, reason)
				continue
			}
			t.Errorf("%s exports %s (%s), which is neither declared by %s nor listed in c21OutTableExempt with a reason - add the table row or exempt it deliberately",
				base, c.name, c.kind, c21TablePath)
		}
	}
	// A stale exemption hides as much as a missing one.
	for name := range c21OutTableExempt {
		if _, ok := tabled[name]; ok {
			t.Errorf("c21OutTableExempt lists %s, which the table now declares: drop the exemption", name)
			continue
		}
		found := false
		for _, f := range files {
			if filepath.Base(f) == "tokens.go" || strings.HasSuffix(f, "_test.go") {
				continue
			}
			for _, c := range c21FileConsts(t, f) {
				if c.name == name {
					found = true
				}
			}
		}
		if !found {
			t.Errorf("c21OutTableExempt lists %s, which the package no longer declares: remove the stale exemption", name)
		}
	}
}

// c21FileConsts is c21ExportedConsts narrowed to the kinds that could carry a
// token value: enums (IconKind, AnimKind, HotkeyStatus, the Win32 hit-test codes)
// and iota members are left out for the same reason as in tokens.go.
func c21FileConsts(t *testing.T, path string) []c21ConstInfo {
	t.Helper()
	var out []c21ConstInfo
	for _, c := range c21ExportedConsts(t, path) {
		if c.kind == "num" || c.kind == "string" || c.kind == "derived" {
			out = append(out, c)
		}
	}
	return out
}
