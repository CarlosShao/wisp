package panel

// C21 design-token four-way reconciliation (ticket 77 AC#2, rebuilt for the
// FOURTH generation by ticket 380).
//
// The four parties, re-pointed at the generation the frontend actually ships:
//
//	P1  frontend/scripts/gen-tokens.mjs          the fourth generation's truth source
//	P2  docs/evidence/s1/c21-native-tokens.md    the hand-maintained C21 table
//	P3  internal/ball/tokens.go                  the native (Direct2D) copy
//	P4  frontend/src/styles/tokens.generated.css the generated frontend theme
//
// P1 used to be design/assets/tokens.css, compared against P4. That file was
// deleted in 72b65c614; only a copy survives as the ball package's own CSS
// fixture, and it belongs to that package's test, not to this one. The fourth
// generation's reference implementation is the TOKENS table inside the
// generator - the committed theme is emitted from it verbatim - which is what
// lets this test compare one value at a time without running node.
//
// POLARITY. The first generation put its default (dark) palette in :root and
// the light one in [data-theme="light"]; the fourth generation inverted that -
// :root is LIGHT (--page: #fafafb) and [data-theme="dark"] is dark. Until
// ticket 380 this file read :root as dark, so P4's dark palette was invisible
// and the light leg parsed nothing at all. parseTokenBlocks now asserts both
// blocks are non-empty: put the old polarity back and the light leg is the
// fatal, exactly as it was.
//
// The four checks (ticket 380 AC#2). They replace "checked < 70", which had no
// object left to measure once the table stopped describing the shipped design:
//
//	1  every name the table cites as a fourth-generation slot is declared in
//	   BOTH P4 blocks - a slot declared in one theme only is half a decision.
//	2  every coloured token of P4 is accounted for exactly once, by the table's
//	   mapping or by the waiver list pinned in this file. The fourth generation
//	   has 34 of them: 11 mapped + 23 waived with a written reason in the table.
//	3  P4's two blocks declare the same name set, and every mapped name owns
//	   exactly one colour row per theme in P2 - so the row count is 2 x names.
//	4  at least c21FourthGenMappedPairs (22 = 11 x 2) of those pairs complete
//	   every leg. 22 is a ceiling, not a lowered bar: --accent-line,
//	   --danger-line and --success-line have no declaration in P4 at all.
//
// P2 -> P3 is re-derived here instead of borrowed from the ball test because
// internal/ball is not this ticket's tree to edit, so the fourth party cannot
// be bolted onto the ball's three-way case. Both tests cover that edge; if
// either starts failing the drift is real, and the day the ball's owner folds
// them together nothing is lost.
//
// Comparison is numeric, not textual: "#86C2B9", "rgba(134,194,185,1)" and
// "hex(0x86C2B9,1)" are three spellings of one decision, and a test that
// compared strings would fail on the spelling while passing on the value.
//
// The AC#2 judgement - "hand-edit one colour in the frontend theme and the test
// must go red" - and the orthogonal shape-vs-value mutation pair are recorded
// in docs/evidence/s1/380r2-impl.md with their readings.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const (
	// P1: the fourth generation's truth source. design/assets/tokens.css - the
	// first generation's reference implementation - was deleted in 72b65c614;
	// the committed theme is emitted from this table verbatim.
	generatorPath     = "frontend/scripts/gen-tokens.mjs"
	c21TablePath      = "docs/evidence/s1/c21-native-tokens.md"
	nativeTokensPath  = "internal/ball/tokens.go"
	frontendThemePath = "frontend/src/styles/tokens.generated.css"
)

// c21FourthGenMappedPairs is check 4's bound: 11 mapped names, one colour row
// per theme, so 22 (name, theme) pairs. It is a ceiling as much as a floor - of
// the 14 native-generation names that have a fourth-generation counterpart,
// three (--accent-line, --danger-line, --success-line) have no declaration in
// the generated theme at all (it carries --warn-line only), so no row for them
// can be reconciled. Deleting a row drops the count below 22 and the fatal
// fires; lowering 22 without changing the contract is the "relax the
// threshold" move ticket 380 exists to prevent.
const c21FourthGenMappedPairs = 22

// c21FourthGenUncovered is the pinned allow-list for check 2: fourth-generation
// colour tokens the native generation deliberately does not cover. The reasons
// live in the table's 第四代未覆盖 section; the NAMES are pinned here so that
// deleting a row and pasting a waiver into the table cannot make the check
// agree with itself - widening the waiver takes a code change, which shows up
// in review. 11 mapped + 23 waived = the fourth generation's 34 coloured
// tokens (43 declarations minus six --shadow-*, --panel-alpha and the two font
// stacks).
var c21FourthGenUncovered = []string{
	"--accent-ink",
	"--accent-tint",
	"--canvas",
	"--field",
	"--green",
	"--green-tint",
	"--hover",
	"--hover-2",
	"--line-strong",
	"--orange",
	"--orange-tint",
	"--page",
	"--red",
	"--red-tint",
	"--stop",
	"--stripe",
	"--stripe-bg",
	"--success-soft",
	"--tooltip-bg",
	"--tooltip-border",
	"--tooltip-fg",
	"--tooltip-muted",
	"--warn-soft",
}

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
	// :root is the LIGHT block in the fourth generation and [data-theme="dark"]
	// the dark one - the inverse of the first generation's convention, which is
	// why reading :root as dark left the light leg empty for a whole generation.
	blockRe    = regexp.MustCompile(`(?s)(:root|\[data-theme="dark"\])\s*\{([^}]*)\}`)
	commentRe  = regexp.MustCompile(`(?s)/\*.*?\*/`)
	declRe     = regexp.MustCompile(`^--([A-Za-z0-9-]+)\s*:\s*(.*)$`)
	plainColRe = regexp.MustCompile(`^(#[0-9A-Fa-f]{3,8}|rgba?\([^)]*\))$`)
	colourInRe = regexp.MustCompile(`rgba?\([^)]*\)|#[0-9A-Fa-f]{3,8}`)
)

// parseTokenBlocks reads the :root (light) and [data-theme="dark"] (dark)
// blocks of a CSS file into per-theme declaration lists.
//
// It reports the declarations written in the file and nothing else. An earlier
// version also synthesised a <name>-color entry for every composite value; that
// existed to line P4 up with tokens.css, which derived the same names. The
// generator's table carries no derived entries, and a name P4 declares that its
// generator does not is precisely what the value leg forbids - so the
// derivation would now manufacture a mismatch. Composite values still compare
// correctly because valuesAgree falls back to the colour slice inside them,
// which is the only part that can reach a utility class.
func parseTokenBlocks(t *testing.T, text, from string) map[string][]tokenDecl {
	t.Helper()
	body := commentRe.ReplaceAllString(text, "")
	out := map[string][]tokenDecl{}
	for _, m := range blockRe.FindAllStringSubmatch(body, -1) {
		theme := "light"
		if strings.HasPrefix(m[1], "[data-theme") {
			theme = "dark"
		}
		var decls []tokenDecl
		for _, chunk := range strings.Split(m[2], ";") {
			// Flattened BEFORE matching: the long values (the --sheen gradient,
			// the --shadow-* stacks, the font stacks) are written across several
			// lines, and a declaration regex anchored to one line would silently
			// drop exactly those tokens.
			chunk = flattenWS(chunk)
			if chunk == "" {
				continue
			}
			dm := declRe.FindStringSubmatch(chunk)
			if dm == nil {
				continue
			}
			decls = append(decls, tokenDecl{name: "--" + dm[1], value: flattenWS(dm[2])})
		}
		out[theme] = decls
	}
	// Both blocks must be non-empty, and that is where the polarity bug lands:
	// with :root read as dark, the light leg had nothing to compare and every
	// row below went unchecked while the test still looked alive.
	if len(out["light"]) == 0 {
		t.Fatalf("%s: no :root block parsed - :root is the LIGHT block in the fourth generation, and an empty light leg makes every check below vacuous", from)
	}
	if len(out["dark"]) == 0 {
		t.Fatalf("%s: no [data-theme=\"dark\"] block parsed - the dark palette would be invisible to every check below", from)
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
	if len(rows) == 0 {
		t.Fatalf("parsed no colour rows from %s - the table shape changed and every check below would be vacuous", c21TablePath)
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

// -------------------------------------------------- fourth generation model

// genToken is one entry of the generator's TOKENS table: the name and the two
// values it emits, one per block.
type genToken struct {
	name  string
	light string
	dark  string
}

// fourthGenSlot is one row of the table's 第四代槽位对照 section: the
// native-generation name on the left, the fourth-generation slot it lands on.
type fourthGenSlot struct {
	where  string
	native string
	slot   string
}

// fourthGenWaiver is one row of the table's 第四代未覆盖 section: a
// fourth-generation token the native generation does not cover, and why.
type fourthGenWaiver struct {
	where  string
	name   string
	reason string
}

var (
	// One entry per token, single line, double-quoted: { name: "--page",
	// light: "#fafafb", dark: "#17181a" } with an optional note: after it.
	genTokenRe = regexp.MustCompile(`\{\s*name:\s*"(--[A-Za-z0-9-]+)",\s*light:\s*"([^"]*)",\s*dark:\s*"([^"]*)"`)
	// A token name in the first cell, backticked or not. "--" alone (a
	// separator row) must not match, hence the [A-Za-z0-9] after the dashes.
	tokenCellRe = regexp.MustCompile(`^--[A-Za-z0-9][A-Za-z0-9-]*$`)
)

// parseGeneratorTokens reads the fourth generation's truth source: the TOKENS
// table the committed theme is emitted from. Parsing it is what lets this test
// check values without running node (the generator's own --check needs node and
// is not on the test path). A shape change that hides an entry shows up as a
// name the generated theme declares and the table does not.
func parseGeneratorTokens(t *testing.T, text, from string) []genToken {
	t.Helper()
	var out []genToken
	for _, m := range genTokenRe.FindAllStringSubmatch(text, -1) {
		out = append(out, genToken{name: m[1], light: m[2], dark: m[3]})
	}
	if len(out) == 0 {
		t.Fatalf("%s: parsed no token entries - the generator's table changed shape and every value check below would be vacuous", from)
	}
	return out
}

// parseC21FourthGenTables walks the two 2-column sections the table carries for
// the fourth generation. They are deliberately two columns wide: the ball
// package's parser skips any row that does not have exactly four cells
// (internal/ball/tokens_table_test.go:181), so these sections are invisible to
// it - which is the reason the 78 native-generation rows can stay untouched.
func parseC21FourthGenTables(t *testing.T, text string) ([]fourthGenSlot, []fourthGenWaiver) {
	t.Helper()
	const (
		sectionSlots   = "slots"
		sectionWaivers = "waivers"
	)
	var (
		slots   []fourthGenSlot
		waivers []fourthGenWaiver
	)
	section := ""
	for i, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			switch {
			case strings.HasPrefix(line, "## 第四代槽位对照"):
				section = sectionSlots
			case strings.HasPrefix(line, "## 第四代未覆盖"):
				section = sectionWaivers
			default:
				section = ""
			}
			continue
		}
		if section == "" || !strings.HasPrefix(line, "|") {
			continue
		}
		cells := splitMarkdownRow(line)
		if len(cells) != 2 {
			continue
		}
		name := strings.Trim(cells[0], "` ")
		if !tokenCellRe.MatchString(name) {
			continue
		}
		where := fmt.Sprintf("%s:%d", c21TablePath, i+1)
		switch section {
		case sectionSlots:
			slots = append(slots, fourthGenSlot{where: where, native: name, slot: strings.Trim(cells[1], "` ")})
		case sectionWaivers:
			waivers = append(waivers, fourthGenWaiver{where: where, name: name, reason: cells[1]})
		}
	}
	if len(slots) == 0 || len(waivers) == 0 {
		t.Fatalf("%s: the fourth-generation sections parsed %d mapped rows and %d waiver rows - both sections must exist and be non-empty", c21TablePath, len(slots), len(waivers))
	}
	return slots, waivers
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
	gen := parseGeneratorTokens(t, readRepoFile(t, root, generatorPath), generatorPath)
	front := parseTokenBlocks(t, readRepoFile(t, root, frontendThemePath), frontendThemePath)
	rows := parseC21ColourTable(t, readRepoFile(t, root, c21TablePath))
	slots, waivers := parseC21FourthGenTables(t, readRepoFile(t, root, c21TablePath))
	native := parseNativePalettes(t, readRepoFile(t, root, nativeTokensPath))

	generated := map[string]map[string]string{"light": {}, "dark": {}}
	for _, g := range gen {
		if prev, dup := generated["light"][g.name]; dup {
			t.Fatalf("%s carries %s twice (light %q / %q) - the table has no single answer", generatorPath, g.name, prev, g.light)
		}
		generated["light"][g.name] = g.light
		generated["dark"][g.name] = g.dark
	}
	declared := map[string]map[string]string{
		"light": indexDecls(t, front["light"], frontendThemePath),
		"dark":  indexDecls(t, front["dark"], frontendThemePath),
	}

	t.Logf("parties: %s=%d tokens, %s=%d colour rows, %s=%d/%d palette fields, %s=%d light + %d dark declarations, %d mapped + %d waived",
		generatorPath, len(gen), c21TablePath, len(rows), nativeTokensPath,
		len(native["dark"]), len(native["light"]), frontendThemePath,
		len(declared["light"]), len(declared["dark"]), len(slots), len(waivers))

	// ------------------------------------------------------------------ value leg
	// P1 -> P4, per theme, one value at a time: the committed theme must be
	// exactly what the generator's table says, nothing more and nothing less.
	// This is the only leg that can see a hand-edited value (--inset dark
	// #1f2022 -> #1f2023): the name set does not move, so no name-level check
	// has anything to complain about.
	reconciled := map[string]int{}
	for _, theme := range []string{"light", "dark"} {
		want := generated[theme]
		got := declared[theme]
		if len(want) != len(got) {
			t.Errorf("%s declares %d tokens in the %s block but %s carries %d entries - one side gained or lost a name",
				frontendThemePath, len(got), theme, generatorPath, len(want))
		}
		for _, name := range sortedNames(want) {
			gv, ok := got[name]
			if !ok {
				t.Errorf("%s: the %s block does not declare %s = %q that %s carries - a token never reached the panel",
					frontendThemePath, theme, name, want[name], generatorPath)
				continue
			}
			same, err := valuesAgree(want[name], gv)
			if err != nil {
				t.Errorf("%s: cannot compare %s in the %s block: %v", frontendThemePath, name, theme, err)
				continue
			}
			if !same {
				t.Errorf("%s: %s declares %s = %q in the %s block, but %s says %q - a hand edit, not a decision",
					frontendThemePath, theme, name, gv, theme, generatorPath, want[name])
				continue
			}
			reconciled[theme]++
		}
		for _, name := range sortedNames(got) {
			if _, ok := want[name]; !ok {
				t.Errorf("%s declares %s = %q in the %s block, which %s does not - a hand-written token reached the panel",
					frontendThemePath, name, got[name], theme, generatorPath)
			}
		}
		t.Logf("value leg theme %s: %d tokens reconciled against %s", theme, reconciled[theme], generatorPath)
	}

	// ---------------------------------------------------------- check 3 (shape)
	// The two blocks must declare the same names: the generator emits every
	// entry into both blocks, so a name in one block only is a truncation, and
	// it would silently halve the dark (or light) side of every check below.
	lightNames, darkNames := nameSet(declared["light"]), nameSet(declared["dark"])
	for _, name := range sortedTrue(lightNames) {
		if !darkNames[name] {
			t.Errorf("check 3: %s declares %s in :root (light) but not in [data-theme=\"dark\"] - the two blocks are one decision", frontendThemePath, name)
		}
	}
	for _, name := range sortedTrue(darkNames) {
		if !lightNames[name] {
			t.Errorf("check 3: %s declares %s in [data-theme=\"dark\"] but not in :root (light) - the two blocks are one decision", frontendThemePath, name)
		}
	}

	// ------------------------------------------- P2 -> P3 (all 78 colour rows)
	// Every colour row of the native-generation table still states the literal
	// tokens.go holds, in the palette constructor that matches its theme. The
	// table's CSS column is no longer compared anywhere in this package: it
	// documents the retired first-generation file, and the ball package owns
	// that edge against its own fixture.
	type rowKey struct{ where, theme string }
	agreed := map[rowKey]bool{}
	rowBy := map[string]map[string]tableColourRow{}
	for _, r := range rows {
		if rowBy[r.cssVar] == nil {
			rowBy[r.cssVar] = map[string]tableColourRow{}
		}
		if prev, dup := rowBy[r.cssVar][r.theme]; dup {
			t.Errorf("%s and %s both claim to be the %s colour row of %s - a name gets one row per theme",
				prev.where, r.where, r.theme, r.cssVar)
		}
		rowBy[r.cssVar][r.theme] = r

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
		same, err := valuesAgree(claim, held)
		if err != nil {
			t.Errorf("%s: %v", r.where, err)
			continue
		}
		if !same {
			t.Errorf("%s: the table says %s holds %q, but %s holds %q",
				r.where, r.goRef, claim, nativeTokensPath, held)
			continue
		}
		agreed[rowKey{r.where, r.theme}] = true
	}

	// --------------------------------- check 1 + check 3, per mapped (name,theme)
	// The table's mapping is the contract between the two generations: each
	// native name lands on one fourth-generation slot, the slot is declared in
	// the block of that theme, and the name owns exactly one colour row per
	// theme (so the row count is 2 x the name count).
	slotOf := map[string]string{}
	for _, s := range slots {
		if prev, dup := slotOf[s.native]; dup {
			t.Errorf("%s maps %s to %s, but an earlier row already mapped it to %s - one name, one slot",
				s.where, s.native, s.slot, prev)
		}
		slotOf[s.native] = s.slot
	}
	checked := 0
	for _, s := range slots {
		for _, theme := range []string{"dark", "light"} {
			r, ok := rowBy[s.native][theme]
			if !ok {
				t.Errorf("check 3: %s maps %s to %s, but the table has no %s colour row for %s - a mapped name gets one row per theme",
					s.where, s.native, s.slot, theme, s.native)
				continue
			}
			if _, ok := declared[theme][s.slot]; !ok {
				t.Errorf("check 1: %s cites %s as the fourth-generation slot of %s, but the %s block of %s does not declare it - half a decision",
					r.where, s.slot, s.native, theme, frontendThemePath)
				continue
			}
			if !agreed[rowKey{r.where, theme}] {
				continue // already reported by the P2 -> P3 pass above
			}
			checked++
		}
	}
	if want := 2 * len(slots); len(slots) != c21FourthGenMappedPairs/2 {
		t.Errorf("check 3: the table maps %d names but the bound %d is 2 x 11 - a mapping change is a contract change, update the bound deliberately (the table would need %d mapped rows)",
			len(slots), c21FourthGenMappedPairs, want)
	} else if rowsMapped := countMappedRows(rowBy, slots); rowsMapped != want {
		t.Errorf("check 3: %d mapped colour rows for %d mapped names, want %d (one per theme each)", rowsMapped, len(slots), want)
	}

	// ---------------------------------------------------------- check 2 (cover)
	// Every coloured token of the fourth generation is accounted for exactly
	// once: mapped by the table, or waived by the list pinned in this file.
	//
	// "Mapped" is a join, not a claim: the mapping row must still have a colour
	// row behind it. Reading the mapping alone would let a deleted row keep
	// covering its slot - the name would stay in the 2-column table while the
	// record it points at is gone, which is the exact failure this ticket is
	// about.
	colourFourth := colouredTokenNames(declared["light"])
	claimed := map[string]string{}
	for _, s := range slots {
		claimed[s.slot] = s.where
	}
	covered := map[string]string{}
	for _, s := range slots {
		for _, theme := range []string{"dark", "light"} {
			if r, ok := rowBy[s.native][theme]; ok {
				covered[s.slot] = r.where
				break
			}
		}
	}
	waived := map[string]string{}
	for _, w := range waivers {
		if prev, dup := waived[w.name]; dup {
			t.Errorf("%s waives %s, already waived at %s - one reason, one row", w.where, w.name, prev)
		}
		waived[w.name] = w.where
		if len([]rune(strings.TrimSpace(w.reason))) < 12 {
			t.Errorf("%s waives %s with a %d-character reason - say which side owns the token instead",
				w.where, w.name, len([]rune(strings.TrimSpace(w.reason))))
		}
	}
	pinned := map[string]bool{}
	for _, name := range c21FourthGenUncovered {
		pinned[name] = true
	}
	for _, name := range sortedTrue(pinned) {
		if _, ok := waived[name]; !ok {
			t.Errorf("check 2: %s is pinned as uncovered but the 第四代未覆盖 section has no row for it - the table and the code list disagree", name)
		}
	}
	for _, name := range colourFourth {
		_, inTable := covered[name]
		_, inWaivers := waived[name]
		switch {
		case inTable && inWaivers:
			t.Errorf("check 2: %s is both mapped (%s) and waived (%s) - a token is covered one way", name, covered[name], waived[name])
		case !inTable && !inWaivers:
			t.Errorf("check 2: %s is a coloured token of %s but neither 第四代槽位对照 maps it nor 第四代未覆盖 waives it - it would go unchecked",
				name, frontendThemePath)
		}
	}
	for _, name := range sortedNames(claimed) {
		if !containsString(colourFourth, name) {
			t.Errorf("check 2: %s is mapped to %s, but %s does not declare a colour there - the mapping points at nothing",
				claimed[name], name, frontendThemePath)
		}
	}
	for _, name := range sortedNames(waived) {
		if !containsString(colourFourth, name) {
			t.Errorf("check 2: %s waives %s, which is not a coloured token of %s - a stale waiver hides a real gap",
				waived[name], name, frontendThemePath)
		}
		if !pinned[name] {
			t.Errorf("check 2: %s waives %s, which is not in the pinned c21FourthGenUncovered list - widening a waiver is a contract change, not a table edit",
				waived[name], name)
		}
	}
	if got, want := len(covered)+len(waived), len(colourFourth); got != want {
		t.Errorf("check 2: %d mapped + %d waived = %d, but %s declares %d coloured tokens - the coverage does not add up",
			len(covered), len(waived), got, frontendThemePath, want)
	}
	t.Logf("check 2: %s declares %d coloured tokens = %d mapped + %d waived",
		frontendThemePath, len(colourFourth), len(covered), len(waived))

	// ------------------------------------------------------------------ check 4
	// 22 is the ceiling as well as the floor, so this is not a lowered bar that
	// any smaller table clears: it is "all 11 mapped names, in both themes".
	if checked < c21FourthGenMappedPairs {
		t.Fatalf("check 4: only %d of the %d mapped (name, theme) pairs completed every leg - the check is not doing its job", checked, c21FourthGenMappedPairs)
	}
	t.Logf("checks 1-4 reconciled %d of %d mapped pairs (%d colour rows read, %d agreed with %s)",
		checked, c21FourthGenMappedPairs, len(rows), len(agreed), nativeTokensPath)
}

// sortedNames returns the keys of m in a stable order, so failure output does
// not shuffle between runs.
func sortedNames(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func nameSet(decls map[string]string) map[string]bool {
	out := make(map[string]bool, len(decls))
	for name := range decls {
		out[name] = true
	}
	return out
}

// sortedTrue returns the set's members in a stable order.
func sortedTrue(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// countMappedRows counts the colour rows that belong to a mapped name, one per
// theme each - the "row count == 2 x name count" half of check 3.
func countMappedRows(rowBy map[string]map[string]tableColourRow, slots []fourthGenSlot) int {
	n := 0
	for _, s := range slots {
		for _, theme := range []string{"dark", "light"} {
			if _, ok := rowBy[s.native][theme]; ok {
				n++
			}
		}
	}
	return n
}

// colouredTokenNames returns the tokens of a block whose value is a colour -
// directly, through a var() alias, or as a color-mix(). Everything else (the
// --shadow-* stacks, --panel-alpha, the two font stacks) is structure, and the
// native generation has no colour to reconcile against it.
func colouredTokenNames(decls map[string]string) []string {
	var out []string
	for name := range decls {
		if valueIsColour(decls, name, 0) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func valueIsColour(decls map[string]string, name string, depth int) bool {
	if depth > 8 {
		return false // an alias cycle; the generator's verify() forbids it anyway
	}
	value, ok := decls[name]
	if !ok {
		return false
	}
	value = flattenWS(value)
	switch {
	case plainColRe.MatchString(value):
		return true
	case strings.HasPrefix(value, "var("):
		target := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "var("), ")"))
		return valueIsColour(decls, target, depth+1)
	case strings.HasPrefix(value, "color-mix("):
		return true
	}
	return false
}
