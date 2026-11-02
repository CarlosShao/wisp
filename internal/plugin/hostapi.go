package plugin

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// C16 host_api version gate (ticket 50 AC2).
//
// Contract sources, verbatim scope:
//   - PLAN.md C16:   "semver；插件声明 host_api = \"^1.0\"；主版本不兼容 → 拒绝加载并给明确
//     错误，不得静默降级"
//   - PLAN.md §14.7 / SPEC-07 §6: the same gate on update and on load; refusal is
//     always an explicit error, never a silent skip.
//   - SPEC-07 §4: host_api is a TOML string, caret form "^1.0" in every example.
//
// What the contract does NOT define, this file does not invent: the only range
// operator that appears anywhere in SPEC/PLAN for host_api is the caret form.
// Any other syntax (bare "1.0", "=1.2.3", ">=1 <2", "*") is refused with an
// explicit "unsupported range syntax" error naming the shape the loader
// accepts — not treated as a pass, not treated as unknown-key noise. The
// minor/patch positions may be omitted ("^1" == "^1.0" == "^1.0.0"), the
// mainstream caret equivalence, because it cannot widen the accepted host set.
// Adding other operators is a contract change (AGENTS.md §0: 改契约＝人工批准).
//
// Caret semantics implemented here are the standard semver ones (node-semver,
// the mainstream harness convention): ^1.2 means ">=1.2.0 <2.0.0" — the host
// major must equal the declared major and the host must be at least the
// declared minimum. For a 0.x major the minor is the compatibility position
// (^0.2 means ">=0.2.0 <0.3.0"). The contract text ("主版本不兼容 → 拒绝；
// minor 在范围内 → 装载") matches this for every declared major >= 1.

// DefaultHostVersion is the host_api contract version the loader gates against:
// buildinfo.Version is the product version stamped by scripts/build.ps1.
// The "go run" default ("0.0.0-dev") matches no caret range, so development
// builds refuse every plugin by design (fail closed) rather than passing one.
var DefaultHostVersion = buildinfo.Version

// HostVersion is a parsed semver version (core, optional pre-release, optional
// build metadata). Build metadata is recorded but never participates in
// precedence (semver §10).
type HostVersion struct {
	Major, Minor, Patch uint64
	Pre                 []string
	Build               []string
}

// CaretRange is a parsed "^<major>.<minor>[.<patch>][-<pre>]" constraint.
type CaretRange struct {
	Raw      string
	Min      HostVersion
	HasPatch bool
}

// ParseHostVersion parses a strict semver 2.0.0 version string.
func ParseHostVersion(s string) (HostVersion, error) {
	if s == "" {
		return HostVersion{}, fmt.Errorf("empty version")
	}
	if strings.ContainsAny(s, " \t") || hasControl(s) {
		return HostVersion{}, fmt.Errorf("version %q contains whitespace or control characters", s)
	}
	var v HostVersion
	rest := s
	if i := strings.IndexByte(rest, '+'); i >= 0 {
		if err := v.setBuild(rest[i+1:]); err != nil {
			return HostVersion{}, fmt.Errorf("version %q: %w", s, err)
		}
		rest = rest[:i]
	}
	if i := strings.IndexByte(rest, '-'); i >= 0 {
		if err := v.setPre(rest[i+1:]); err != nil {
			return HostVersion{}, fmt.Errorf("version %q: %w", s, err)
		}
		rest = rest[:i]
	}
	parts := strings.Split(rest, ".")
	if len(parts) != 3 {
		return HostVersion{}, fmt.Errorf("version %q: core must be <major>.<minor>.<patch>", s)
	}
	nums := []*uint64{&v.Major, &v.Minor, &v.Patch}
	for i, p := range parts {
		n, err := parseNumericIdentifier(p)
		if err != nil {
			return HostVersion{}, fmt.Errorf("version %q: %s part: %w", s, [...]string{"major", "minor", "patch"}[i], err)
		}
		*nums[i] = n
	}
	return v, nil
}

func (v *HostVersion) setPre(s string) error {
	if s == "" {
		return fmt.Errorf("empty pre-release after '-'")
	}
	for _, id := range strings.Split(s, ".") {
		if id == "" || !isLegalAlnumIdentifier(id) {
			return fmt.Errorf("illegal pre-release identifier %q", id)
		}
		v.Pre = append(v.Pre, id)
	}
	return nil
}

func (v *HostVersion) setBuild(s string) error {
	if s == "" {
		return fmt.Errorf("empty build metadata after '+'")
	}
	for _, id := range strings.Split(s, ".") {
		if id == "" || !isLegalAlnumIdentifier(id) {
			return fmt.Errorf("illegal build identifier %q", id)
		}
		v.Build = append(v.Build, id)
	}
	return nil
}

// String renders the parsed form back to semver text (for error messages).
func (v HostVersion) String() string {
	s := strconv.FormatUint(v.Major, 10) + "." +
		strconv.FormatUint(v.Minor, 10) + "." +
		strconv.FormatUint(v.Patch, 10)
	if len(v.Pre) > 0 {
		s += "-" + strings.Join(v.Pre, ".")
	}
	if len(v.Build) > 0 {
		s += "+" + strings.Join(v.Build, ".")
	}
	return s
}

// ParseCaretRange parses the only host_api grammar C16 defines: a caret with an
// explicit major. An omitted minor/patch is read as 0 (node-semver: "^1.0" ==
// "^1.0.0" == "^1"), which is what the SPEC-07 §4 example relies on.
func ParseCaretRange(s string) (CaretRange, error) {
	if s == "" {
		return CaretRange{}, fmt.Errorf("host_api is empty")
	}
	if strings.ContainsAny(s, " \t") || hasControl(s) {
		return CaretRange{}, fmt.Errorf("host_api %q contains whitespace or control characters", s)
	}
	if !strings.HasPrefix(s, "^") {
		return CaretRange{}, fmt.Errorf("host_api %q: unsupported range syntax; C16 defines only the caret form \"^<major>.<minor>\" (e.g. \"^1.0\")", s)
	}
	body := s[1:]
	if body == "" {
		return CaretRange{}, fmt.Errorf("host_api %q: caret needs a version after it", s)
	}
	var v HostVersion
	rest := body
	if i := strings.IndexByte(rest, '-'); i >= 0 {
		if err := v.setPre(rest[i+1:]); err != nil {
			return CaretRange{}, fmt.Errorf("host_api %q: %w", s, err)
		}
		rest = rest[:i]
	}
	parts := strings.Split(rest, ".")
	if len(parts) > 3 {
		return CaretRange{}, fmt.Errorf("host_api %q: at most <major>.<minor>.<patch>", s)
	}
	nums := []*uint64{&v.Major, &v.Minor, &v.Patch}
	for i, p := range parts {
		n, err := parseNumericIdentifier(p)
		if err != nil {
			return CaretRange{}, fmt.Errorf("host_api %q: %s part: %w", s, [...]string{"major", "minor", "patch"}[i], err)
		}
		*nums[i] = n
	}
	r := CaretRange{Raw: s, Min: v, HasPatch: len(parts) == 3}
	return r, nil
}

// String renders the range back to its caret text.
func (r CaretRange) String() string { return r.Raw }

// CheckHostAPI is the C16 load gate: it returns nil when the host satisfies the
// plugin's declared caret range, and otherwise a load refusal carrying the
// plugin id, both versions, and the violated condition. It never stays silent:
// a range that cannot even be parsed is also a refusal (the error says so).
func CheckHostAPI(pluginID, rangeStr, hostVersion string) error {
	r, err := ParseCaretRange(rangeStr)
	if err != nil {
		return fmt.Errorf("plugin %q: host_api %q is invalid: %w", pluginID, rangeStr, err)
	}
	h, err := ParseHostVersion(hostVersion)
	if err != nil {
		return fmt.Errorf("plugin %q: host version %q is invalid: %w", pluginID, hostVersion, err)
	}
	return r.Check(pluginID, h)
}

// Check runs the gate against an already parsed host version.
func (r CaretRange) Check(pluginID string, h HostVersion) error {
	if r.Min.Major != 0 {
		if h.Major != r.Min.Major {
			return majorMismatch(pluginID, r, h)
		}
		if cmpCorePre(h, r.Min) < 0 {
			return belowMin(pluginID, r, h)
		}
		return nil
	}
	// Major 0: minor is the compatibility position (^0.2 => >=0.2.0 <0.3.0).
	if h.Major != 0 {
		return majorMismatch(pluginID, r, h)
	}
	if r.Min.Minor != 0 {
		if h.Minor != r.Min.Minor {
			return majorMismatch(pluginID, r, h)
		}
		if cmpCorePre(h, r.Min) < 0 {
			return belowMin(pluginID, r, h)
		}
		return nil
	}
	// ^0.0: patch is the compatibility position (>=0.0.<patch> <0.1.0).
	if h.Minor != 0 || cmpCorePre(h, r.Min) < 0 {
		return belowMin(pluginID, r, h)
	}
	return nil
}

func majorMismatch(pluginID string, r CaretRange, h HostVersion) error {
	return fmt.Errorf("plugin %q refused: host_api %q requires host major %d, this host is %s — "+
		"C16: major-version incompatibility must stop the load with this error, never a silent skip or downgrade",
		pluginID, r.Raw, r.Min.Major, h)
}

func belowMin(pluginID string, r CaretRange, h HostVersion) error {
	return fmt.Errorf("plugin %q refused: host version %s is below the minimum %s declared by host_api %q — "+
		"C16: the host must provide at least the API level the plugin was written against",
		pluginID, h, r.Min, r.Raw)
}

// CompatibleWith applies the manifest's declared host_api range to a host
// version string. The loader calls this before registering a plugin.
func (m Manifest) CompatibleWith(hostVersion string) error {
	return CheckHostAPI(m.ID, m.HostAPI, hostVersion)
}

// cmpCorePre compares two versions by precedence (semver §11) ignoring build
// metadata. Within one range the upper bound is always free of pre-releases,
// so this fully covers the gate's decisions.
func cmpCorePre(a, b HostVersion) int {
	if c := cmpUint(a.Major, b.Major); c != 0 {
		return c
	}
	if c := cmpUint(a.Minor, b.Minor); c != 0 {
		return c
	}
	if c := cmpUint(a.Patch, b.Patch); c != 0 {
		return c
	}
	return cmpPre(a.Pre, b.Pre)
}

// cmpPre implements semver §11.4: a version without pre-release outranks one
// with it; identifiers compare numeric-as-integers, alnum-as-ASCII, numeric
// identifiers of equal value still rank below alnum ones.
func cmpPre(a, b []string) int {
	switch {
	case len(a) == 0 && len(b) == 0:
		return 0
	case len(a) == 0:
		return 1
	case len(b) == 0:
		return -1
	}
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		x, y := a[i], b[i]
		if x == y {
			continue
		}
		xn, xe := strconv.ParseUint(x, 10, 64)
		yn, ye := strconv.ParseUint(y, 10, 64)
		switch {
		case xe == nil && ye == nil:
			if c := cmpUint(xn, yn); c != 0 {
				return c
			}
		case xe == nil:
			return -1 // numeric identifiers rank below alnum
		case ye == nil:
			return 1
		default:
			if x < y {
				return -1
			}
			return 1
		}
	}
	return cmpUint(uint64(len(a)), uint64(len(b)))
}

func cmpUint(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func parseNumericIdentifier(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty numeric identifier")
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return 0, fmt.Errorf("numeric identifier %q has non-digit character %q", s, r)
		}
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("numeric identifier %q has a leading zero", s)
	}
	return strconv.ParseUint(s, 10, 64)
}

func isLegalAlnumIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		isAlnum := (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '-'
		if !isAlnum {
			return false
		}
	}
	return true
}
