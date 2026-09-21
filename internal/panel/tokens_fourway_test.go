package panel

// C21 design-token four-way reconciliation (ticket 77 AC#2).
//
// Ticket 74 delivered a THREE-way check (TestC21TableColourRowsMatchTokensCSS,
// in internal/ball): tokens.css = the C21 table's CSS column = tokens.go. This
// file adds the fourth party the panel brings with it - the generated frontend
// theme - and pins all four edges in one test:
//
//	P1  design/assets/tokens.css                 the C21 reference implementation
//	P2  docs/evidence/s1/c21-native-tokens.md    the hand-maintained C21 table
//	P3  internal/ball/tokens.go                  the native (Direct2D) copy
//	P4  frontend/src/styles/tokens.generated.css the frontend theme
//
// P4 is produced by frontend/scripts/gen-tokens.mjs from P1 and committed, so
// building the binary needs no node. What this test refuses:
//
//	P1 -> P4  every token P1 declares (plus the derived <name>-color slice of a
//	          composite value) appears in P4 with the same colour, and P4
//	          declares nothing P1 does not.
//	P2 -> P1  every colour row of the table states the value P1 declares.
//	P2 -> P3  every colour row's Go claim is the literal tokens.go holds, in the
//	          palette constructor that matches the row's theme.
//	P3 -> P4  the value the native side holds is the value the panel renders.
//
// Why P2->P3 is re-derived here instead of borrowed from the ball test:
// internal/ball is not this ticket's tree to edit (the ticket says so), so the
// fourth party cannot be bolted onto the three-way case. Both tests now cover
// the shared edges; if either starts failing the drift is real, and the day the
// ball's owner folds them together nothing is lost.
//
// Comparison is numeric, not textual: "#86C2B9", "rgba(134,194,185,1)" and
// "hex(0x86C2B9,1)" are three spellings of one decision, and a test that
// compared strings would fail on the spelling while passing on the value.
//
// The AC#2 judgement - "hand-edit one colour in the frontend theme and the test
// must go red" - is a mutation run recorded in the ticket's Progress log.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const (
	tokensCSSPath     = "design/assets/tokens.css"
	c21TablePath      = "docs/evidence/s1/c21-native-tokens.md"
	nativeTokensPath  = "internal/ball/tokens.go"
	frontendThemePath = "frontend/src/styles/tokens.generated.css"
)

// panelRepoRoot walks up from the test's working directory to the checkout.
func panelRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the working directory - cannot locate the repository root")
		}
		dir = parent
	}
}

func readRepoFile(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

// ---------------------------------------------------------------- token model

type tokenDecl struct {
	name  string
	value string // whitespace-flattened, exactly as the generator writes it
}

var (
	blockRe    = regexp.MustCompile(`(?s)(:root|\[data-theme="light"\])\s*\{([^}]*)\}`)
	commentRe  = regexp.MustCompile(`(?s)/\*.*?\*/`)
	declRe     = regexp.MustCompile(`^--([A-Za-z0-9-]+)\s*:\s*(.*)$`)
	plainColRe = regexp.MustCompile(`^(#[0-9A-Fa-f]{3,8}|rgba?\([^)]*\))$`)
	colourInRe = regexp.MustCompile(`rgba?\([^)]*\)|#[0-9A-Fa-f]{3,8}`)
)

// parseTokenBlocks reads the :root and [data-theme="light"] blocks of a CSS
// file into per-theme declaration lists, applying the same two derivations
// scripts/gen-tokens.mjs documents (whitespace flattening, and the
// <name>-color slice taken out of a composite value). A generated file that
// carries the derived lines literally is therefore checked against what the
// source re-derives here, which is what makes a hand edit visible.
func parseTokenBlocks(t *testing.T, text, from string) map[string][]tokenDecl {
	t.Helper()
	body := commentRe.ReplaceAllString(text, "")
	out := map[string][]tokenDecl{}
	for _, m := range blockRe.FindAllStringSubmatch(body, -1) {
		theme := "dark"
		if strings.HasPrefix(m[1], "[data-theme") {
			theme = "light"
		}
		var decls []tokenDecl
		for _, chunk := range strings.Split(m[2], ";") {
			// Flattened BEFORE matching: tokens.css writes the long values
			// (the --sheen gradient, the --shadow-* stacks, the font stacks)
			// across several lines, and a declaration regex anchored to one
			// line would silently drop exactly those tokens.
			chunk = flattenWS(chunk)
			if chunk == "" {
				continue
			}
			dm := declRe.FindStringSubmatch(chunk)
			if dm == nil {
				continue
			}
			value := flattenWS(dm[2])
			decls = append(decls, tokenDecl{name: "--" + dm[1], value: value})
			if !plainColRe.MatchString(value) {
				if slice := lastColour(value); slice != "" {
					decls = append(decls, tokenDecl{name: "--" + dm[1] + "-color", value: slice})
				}
			}
		}
		out[theme] = decls
	}
	if len(out["dark"]) == 0 {
		t.Fatalf("%s: no :root block parsed", from)
	}
	return out
}

// lastColour returns the final colour literal inside a composite value - the
// slice the panel can put on a utility class when the whole value is a
// background stack (C21 layers --bg-raised as var(--sheen) plus a colour).
func lastColour(value string) string {
	found := colourInRe.FindAllString(value, -1)
	if len(found) == 0 {
		return ""
	}
	return found[len(found)-1]
}

func flattenWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ------------------------------------------------------------------ table model

type tableColourRow struct {
	where    string
	theme    string // "dark" | "light"
	cssVar   string // "--accent", verbatim
	cssValue string // what the table says tokens.css declares
	goRef    string // "Palette.Accent"
	goClaim  string // "hex(0x86C2B9,1)" / "rgba(6,8,11,0.38)"
}

var (
	tableSectionRe = regexp.MustCompile(`^## (明色系|亮色系)`)
	tableRowHeadRe = regexp.MustCompile("^`--")
)

// parseC21ColourTable walks the two colour tables of the C21 document. Rows
// whose first column is not a CSS variable (the native-only liquid looks) are
// skipped for the same reason ticket 74 skipped them: they cite no token.
func parseC21ColourTable(t *testing.T, text string) []tableColourRow {
	t.Helper()
	var rows []tableColourRow
	theme, inTable := "", false
	for i, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			inTable = tableSectionRe.MatchString(line)
			switch {
			case strings.HasPrefix(line, "## 明色系"):
				theme = "dark"
			case strings.HasPrefix(line, "## 亮色系"):
				theme = "light"
			default:
				theme = ""
			}
			continue
		}
		if !inTable || !strings.HasPrefix(line, "|") {
			continue
		}
		cells := splitMarkdownRow(line)
		if len(cells) < 4 || !tableRowHeadRe.MatchString(cells[0]) {
			continue
		}
		rows = append(rows, tableColourRow{
			where:    fmt.Sprintf("%s:%d", c21TablePath, i+1),
			theme:    theme,
			cssVar:   strings.Trim(cells[0], "` "),
			cssValue: flattenWS(strings.Trim(cells[1], "` ")),
			goRef:    strings.Trim(cells[2], "` "),
			goClaim:  flattenWS(strings.Trim(cells[3], "` ")),
		})
	}
	if len(rows) < 70 {
		t.Fatalf("parsed only %d colour rows from %s - the table shape changed and this check would be vacuous", len(rows), c21TablePath)
	}
	return rows
}

func splitMarkdownRow(line string) []string {
	trimmed := strings.Trim(strings.TrimSpace(line), "|")
	parts := strings.Split(trimmed, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// ------------------------------------------------------------- native parsing

var (
	nativePaletteRe = regexp.MustCompile(`^func (DarkPalette|LightPalette)\(\) Palette \{`)
	nativeAssignRe  = regexp.MustCompile(`^\s*p\.([A-Z][A-Za-z0-9]*)\s*=\s*(.+)$`)
	nativeFieldRe   = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	goHexClaimRe    = regexp.MustCompile(`^hex\(\s*0[xX][0-9A-Fa-f]{6}\s*,\s*[0-9.]+\s*\)`)
	goRGBAClaimRe   = regexp.MustCompile(`^rgba\(\s*[0-9]{1,3}\s*,\s*[0-9]{1,3}\s*,\s*[0-9]{1,3}\s*,\s*[0-9.]+\s*\)`)
)

// parseNativePalettes returns theme -> Go field -> colour literal text, read out
// of the two palette constructors of internal/ball/tokens.go.
func parseNativePalettes(t *testing.T, text string) map[string]map[string]string {
	t.Helper()
	out := map[string]map[string]string{}
	theme := ""
	for i, line := range strings.Split(text, "\n") {
		if m := nativePaletteRe.FindStringSubmatch(line); m != nil {
			theme = "dark"
			if m[1] == "LightPalette" {
				theme = "light"
			}
			out[theme] = map[string]string{}
			continue
		}
		if theme == "" {
			continue
		}
		if strings.HasPrefix(line, "func ") {
			theme = ""
			continue
		}
		var field, rest string
		if m := nativeAssignRe.FindStringSubmatch(line); m != nil {
			// LightPalette() mutates a copy of DarkPalette() ("p.BGInset =
			// rgba(...)"), so the light theme's literals arrive in a different
			// shape than the dark theme's composite literal. A field the light
			// block does not redefine keeps the dark value, which is what the
			// table's "not redefined in light" note records - and why the light
			// table carries 38 rows against dark's 40.
			field, rest = m[1], m[2]
		} else {
			field = strings.TrimSpace(strings.SplitN(line, ":", 2)[0])
			if !nativeFieldRe.MatchString(field) {
				continue
			}
			rest = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}
		lit, err := colourLiteralAt(rest)
		if err != nil {
			continue
		}
		if _, dup := out[theme][field]; dup {
			t.Fatalf("tokens.go:%d: field %s parsed twice in %s - the parser no longer matches the file",
				i+1, field, theme)
		}
		out[theme][field] = lit
	}
	if len(out["dark"]) == 0 || len(out["light"]) == 0 {
		t.Fatalf("native palettes not parsed (dark %d / light %d fields)", len(out["dark"]), len(out["light"]))
	}
	return out
}

// colourLiteralAt cuts the leading colour constructor out of a value, dropping a
// trailing comma and anything after a comment marker.
func colourLiteralAt(s string) (string, error) {
	s = strings.TrimSpace(strings.SplitN(s, "//", 2)[0])
	s = strings.TrimSuffix(strings.TrimSpace(s), ",")
	for _, re := range []*regexp.Regexp{goHexClaimRe, goRGBAClaimRe} {
		if m := re.FindString(s); m != "" {
			return flattenWS(m), nil
		}
	}
	return "", fmt.Errorf("no colour constructor at %q", s)
}

// ---------------------------------------------------------------- colour math

type argb struct{ r, g, b, a float64 }

// evalColour turns either side's spelling of a colour into numbers.
func evalColour(s string) (argb, error) {
	s = flattenWS(s)
	switch {
	case strings.HasPrefix(s, "#"):
		digits := s[1:]
		if len(digits) == 3 {
			digits = string(digits[0]) + string(digits[0]) +
				string(digits[1]) + string(digits[1]) +
				string(digits[2]) + string(digits[2])
		}
		if len(digits) != 6 {
			return argb{}, fmt.Errorf("unsupported hex colour %q", s)
		}
		v, err := strconv.ParseUint(digits, 16, 32)
		if err != nil {
			return argb{}, fmt.Errorf("bad hex colour %q: %w", s, err)
		}
		return argb{float64(v>>16&0xFF) / 255, float64(v>>8&0xFF) / 255, float64(v&0xFF) / 255, 1}, nil

	case strings.HasPrefix(s, "rgba("), strings.HasPrefix(s, "rgb("):
		inner := s[strings.Index(s, "(")+1 : strings.LastIndex(s, ")")]
		parts := strings.Split(inner, ",")
		if len(parts) != 3 && len(parts) != 4 {
			return argb{}, fmt.Errorf("bad rgb colour %q", s)
		}
		channels := []float64{0, 0, 0}
		for i := 0; i < 3; i++ {
			n, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
			if err != nil {
				return argb{}, fmt.Errorf("bad rgb channel in %q: %w", s, err)
			}
			channels[i] = n / 255
		}
		c := argb{r: channels[0], g: channels[1], b: channels[2], a: 1}
		if len(parts) == 4 {
			n, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
			if err != nil {
				return argb{}, fmt.Errorf("bad rgb alpha in %q: %w", s, err)
			}
			c.a = n
		}
		return c, nil

	case strings.HasPrefix(s, "hex("):
		inner := s[strings.Index(s, "(")+1 : strings.LastIndex(s, ")")]
		parts := strings.Split(inner, ",")
		if len(parts) != 2 {
			return argb{}, fmt.Errorf("bad hex() claim %q", s)
		}
		v, err := strconv.ParseUint(strings.TrimSpace(strings.TrimPrefix(parts[0], "0x")), 16, 32)
		if err != nil {
			return argb{}, fmt.Errorf("bad hex() value in %q: %w", s, err)
		}
		a, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return argb{}, fmt.Errorf("bad hex() alpha in %q: %w", s, err)
		}
		return argb{float64(v>>16&0xFF) / 255, float64(v>>8&0xFF) / 255, float64(v&0xFF) / 255, a}, nil
	}
	return argb{}, fmt.Errorf("not a colour: %q", s)
}

func (c argb) String() string {
	return fmt.Sprintf("rgba(%.0f,%.0f,%.0f,%.3f)", c.r*255, c.g*255, c.b*255, c.a)
}

func sameColour(a, b argb) bool {
	const eps = 0.5 / 255
	for _, pair := range [][2]float64{{a.r, b.r}, {a.g, b.g}, {a.b, b.b}, {a.a, b.a}} {
		if diff := pair[0] - pair[1]; diff > eps || diff < -eps {
			return false
		}
	}
	return true
}

// valuesAgree compares two spellings of one token value. Colour constructors
// are compared numerically, because the same decision is written as
// "#86C2B9" in CSS, "rgba(134, 194, 185, 1)" in the table and
// "hex(0x86C2B9, 1)" in Go - and the whitespace inside those spellings is not
// the contract. Composite values (background stacks, shadows) fall back to the
// colour slice inside them, the only part that can reach a utility class.
func valuesAgree(a, b string) (bool, error) {
	fa, fb := flattenWS(a), flattenWS(b)
	if fa == fb {
		return true, nil
	}
	ca, ea := evalColour(fa)
	cb, eb := evalColour(fb)
	if ea == nil && eb == nil {
		return sameColour(ca, cb), nil
	}
	sa, sb := lastColour(fa), lastColour(fb)
	if sa != "" && sb != "" {
		ca, ea = evalColour(sa)
		cb, eb = evalColour(sb)
		if ea == nil && eb == nil {
			return sameColour(ca, cb), nil
		}
	}
	// At least one side is not a colour anywhere (a duration, a radius, a font
	// stack): the flattened strings above were the whole comparison, and they
	// already differed.
	if ea == nil || eb == nil {
		return false, nil
	}
	return false, fmt.Errorf("neither %q nor %q is a readable colour", a, b)
}

func indexDecls(t *testing.T, decls []tokenDecl, from string) map[string]string {
	t.Helper()
	out := make(map[string]string, len(decls))
	for _, d := range decls {
		if prev, dup := out[d.name]; dup && prev != d.value {
			t.Fatalf("%s declares %s twice with different values (%q / %q) - the file has no single answer",
				from, d.name, prev, d.value)
		}
		out[d.name] = d.value
	}
	return out
}

func themeFunc(theme string) string {
	if theme == "light" {
		return "LightPalette()"
	}
	return "DarkPalette()"
}

// ---------------------------------------------------------------- the one test

func TestC21DesignTokensFourWayAgree(t *testing.T) {
	root := panelRepoRoot(t)
	css := parseTokenBlocks(t, readRepoFile(t, root, tokensCSSPath), tokensCSSPath)
	front := parseTokenBlocks(t, readRepoFile(t, root, frontendThemePath), frontendThemePath)
	rows := parseC21ColourTable(t, readRepoFile(t, root, c21TablePath))
	native := parseNativePalettes(t, readRepoFile(t, root, nativeTokensPath))

	t.Logf("parties: %s=%d dark decls, %s=%d colour rows, %s=%d/%d palette fields, %s=%d dark decls",
		tokensCSSPath, len(css["dark"]), c21TablePath, len(rows), nativeTokensPath,
		len(native["dark"]), len(native["light"]), frontendThemePath, len(front["dark"]))

	// Leg 1 (P1 -> P4 and back), per theme: the generated theme must be exactly
	// what tokens.css generates, nothing more and nothing less.
	for _, theme := range []string{"dark", "light"} {
		want := indexDecls(t, css[theme], tokensCSSPath)
		got := indexDecls(t, front[theme], frontendThemePath)
		if len(want) == 0 {
			t.Fatalf("theme %s: %s parsed no declarations, this check would be vacuous", theme, tokensCSSPath)
		}
		for name, value := range want {
			fv, ok := got[name]
			if !ok {
				t.Errorf("%s does not carry %s = %s that %s declares - a token never reached the panel",
					frontendThemePath, name, value, tokensCSSPath)
				continue
			}
			same, err := valuesAgree(value, fv)
			if err != nil {
				t.Errorf("%s: cannot compare %s: %v", theme, name, err)
				continue
			}
			if !same {
				t.Errorf("%s declares %s = %q but %s carries %q - two style sources disagree",
					tokensCSSPath, name, value, frontendThemePath, fv)
			}
		}
		for name, value := range got {
			if _, ok := want[name]; !ok {
				t.Errorf("%s declares %s = %q, which %s does not - a hand-written token reached the panel",
					frontendThemePath, name, value, tokensCSSPath)
			}
		}
		t.Logf("leg 1 theme %s: %d tokens reconciled", theme, len(want))
	}

	// Legs 2-4: the table against the CSS file, against tokens.go, and the
	// native value against what the panel renders.
	checked := 0
	for _, r := range rows {
		cssDecl := indexDecls(t, css[r.theme], tokensCSSPath)
		frontDecl := indexDecls(t, front[r.theme], frontendThemePath)
		decl, ok := cssDecl[r.cssVar]
		if !ok {
			t.Errorf("%s: the table cites %s from the %s theme, %s does not declare it",
				r.where, r.cssVar, r.theme, tokensCSSPath)
			continue
		}
		// leg 2: table -> tokens.css
		if same, err := valuesAgree(r.cssValue, decl); err != nil {
			t.Errorf("%s: %v", r.where, err)
			continue
		} else if !same {
			t.Errorf("%s: the table states %s = %q for theme %s, but %s declares %q",
				r.where, r.cssVar, r.cssValue, r.theme, tokensCSSPath, decl)
			continue
		}
		// leg 3: table -> tokens.go
		field := r.goRef
		if i := strings.LastIndex(field, "."); i >= 0 {
			field = field[i+1:]
		}
		held, ok := native[r.theme][field]
		if !ok {
			t.Errorf("%s: the table names %s, but %s holds no such colour literal in %s",
				r.where, r.goRef, nativeTokensPath, themeFunc(r.theme))
			continue
		}
		claim, err := colourLiteralAt(r.goClaim)
		if err != nil {
			t.Errorf("%s: %v", r.where, err)
			continue
		}
		if same, err := valuesAgree(claim, held); err != nil {
			t.Errorf("%s: %v", r.where, err)
			continue
		} else if !same {
			t.Errorf("%s: the table says %s holds %q, but %s holds %q",
				r.where, r.goRef, claim, nativeTokensPath, held)
			continue
		}
		// leg 4: tokens.go -> the frontend theme, through the token the row cites.
		fv, ok := frontDecl[r.cssVar]
		if !ok {
			t.Errorf("%s: %s declares %s but %s does not carry it",
				r.where, tokensCSSPath, r.cssVar, frontendThemePath)
			continue
		}
		if same, err := valuesAgree(claim, fv); err != nil {
			t.Errorf("%s: %v", r.where, err)
			continue
		} else if !same {
			t.Errorf("%s: native %s holds %q (%s) but the panel renders %s = %q",
				r.where, r.goRef, claim, themeFunc(r.theme), r.cssVar, fv)
			continue
		}
		checked++
	}
	if checked < 70 {
		t.Fatalf("only %d of %d colour rows completed all four legs - the check is not doing its job", checked, len(rows))
	}
	t.Logf("legs 2-4 reconciled for %d of %d colour rows", checked, len(rows))
}
