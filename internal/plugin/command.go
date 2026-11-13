package plugin

// D46 command plugins — ticket 52 (SPEC-07 §4 / §16.5.6, C14 D46 fields).
//
// This file owns everything about a kind = "command" tool that CAN be decided
// without spawning a process, plus the one decision that must happen at pin
// time:
//
//  1. INSTALL-TIME PINNING (PinCommandExe): the wrapped executable is resolved
//     through C26, refused unless it is an absolute, stable, fixed local path,
//     and its bytes are hashed. The resulting Pin is what plugin_state stores
//     as exe_hash. A manifest that pins a hash the bytes do not match is
//     refused here as a local swap.
//  2. EVERY-EXECUTION RECHECK (Pin.Verify): the same bytes are re-hashed before
//     a call and a mismatch is an alert + refusal, never a warning. The runtime
//     in internal/plugin/manifestslot calls this on every call; the test leg
//     in this package proves the alert fires.
//
// The execution half - argv building, environment allowlist, child process,
// timeouts, output caps and extraction - lives in the manifestslot runtime
// (command.go there), because it needs tools.Result and the C11 scope and
// internal/plugin must not import internal/tools (import cycle, see
// manifestslot.go). TOGETHER the two halves are ticket 52; neither alone is.
//
// What this file deliberately does NOT do: it never builds a command STRING.
// The argv template is a vector ([]string) from the manifest to the process
// creation call, and the only text ever substituted into it comes from
// parameters that passed both the JSON Schema and the metacharacter screen.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// D46 execution budget constants (SPEC-07 §4, one source of truth shared by
// the validator, the pin and the runtime).
const (
	// DefaultCommandTimeoutMS is the manifest-overridable default (D46:
	// "默认 timeoutMs = 15000").
	DefaultCommandTimeoutMS = 15000
	// MinCommandTimeoutMS and MaxCommandTimeoutMS bound an explicit value. A
	// sub-50ms budget cannot start a process on Windows; an hour-long one is a
	// hung task wearing a deadline.
	MinCommandTimeoutMS = 50
	MaxCommandTimeoutMS = 10 * 60 * 1000

	// CommandStdoutCapBytes / CommandStderrCapBytes are the per-stream caps
	// ("stdout/stderr 各 <=1MB 超出即 kill 并标 truncated").
	CommandStdoutCapBytes = 1 << 20
	CommandStderrCapBytes = 1 << 20
	// CommandExtractCapBytes is what may enter the model context from ONE
	// call ("输出提取 <=32KB 进上下文"). Raw stdout never gets there.
	CommandExtractCapBytes = 32 * 1024
)

// Command execution refusals. Callers branch on these; none of them is a
// warning.
var (
	// ErrExeUnstable is the install-time refusal for an executable that is not
	// an absolute, fixed, local path (relative, %VAR%, ~, TEMP under any of its
	// spellings, or not local at all).
	ErrExeUnstable = errors.New("plugin: command exe is not a fixed absolute local path (D46: a movable executable cannot be hash-pinned)")
	// ErrExeSwap is the runtime refusal: the bytes on disk are not the bytes
	// pinned at install.
	ErrExeSwap = errors.New("plugin: command exe bytes do not match the pinned exe_hash (possible in-place swap)")
	// ErrNoPin is the runtime refusal when nothing was pinned at all. An
	// unpinned exe is exactly the supply-chain hole D46 closes.
	ErrNoPin = errors.New("plugin: no pinned exe_hash for this command tool; refusing execution")
)

// ExePin is the pinned identity of one wrapped executable: the C26-canonical
// path plus the sha256 of its bytes ("sha256:<64 hex>", same shape as
// plugin_state.hash and C14 exe_hash).
type ExePin struct {
	// Path is the canonical (C26) absolute path the bytes were read from.
	Path string
	// Hash is "sha256:<hex>".
	Hash string
	// Size is the byte length the hash was taken over (diagnostics).
	Size int64
}

func (p ExePin) String() string { return p.Path + " @ " + p.Hash }

// commandExeExts are the extensions a pinned local CLI may carry on Windows. A
// batch or PowerShell script is never a D46 exe: cmd.exe/sh would have to
// interpret it, and "never through a shell" is the whole point of the channel.
var commandExeExts = []string{".exe", ".com"}

// CommandExeAllowed reports whether name is a form D46 can execute directly
// (a real executable image) as opposed to something a shell would have to
// interpret (.bat, .cmd, .ps1, .sh, or nothing at all).
func CommandExeAllowed(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	for _, e := range commandExeExts {
		if ext == e {
			return true
		}
	}
	return false
}

// unstableExeMarkers are the spellings that make a path movable. They are
// checked on the RAW manifest string before C26 expands anything: C26 expands
// %VAR% and ~ by contract (SPEC-06 §4 step 1), so a manifest that names
// `%TEMP%\x.exe` would otherwise be pinned to today's temp directory and called
// "stable". D46's red-team case is exactly this one (PLAN.md:2678).
//
// The tilde is deliberately NOT a substring marker: its home-relative meaning
// lives only in the LEADING position (expandInput expands a leading ~, never
// one inside a component), and GitHub's windows runners answer TEMP through an
// 8.3 alias (C:\Users\RUNNER~1\AppData\Local\Temp) whose mid-component tilde is
// a FIXED alias of the runner user's profile, not a movable spelling. Marking
// that tilde movable refused the runner's own TEMP tree with the wrong reason
// and buried the TEMP rule (run 43362781437's test-windows log); the
// home-relative half of the old check lives on in isHomeRelativeSpelling.
var unstableExeMarkers = []string{"%", "$"}

// isHomeRelativeSpelling is the tilde half of the old marker check: a LEADING
// ~ (bare, ~/ or ~\ - the spellings expandInput treats as home) resolves
// against whoever runs the manifest and cannot be hash-pinned. A tilde inside
// a component is an 8.3 alias artifact (RUNNER~1) and is not refused here.
func isHomeRelativeSpelling(p string) bool {
	return strings.HasPrefix(p, "~")
}

// tempDirEnvnames are the variables whose value would make an exe live in a
// directory the OS recycles. Checked as a component of the canonical path, so
// an unusual TEMP location is caught as well as the default one.
var tempDirEnvnames = []string{"TEMP", "TMP", "TMPDIR"}

// resolveExe runs the C26 pipeline over a manifest exe spelling and applies the
// D46 stability rules. It is the shared front half of install and runtime so a
// call can never agree with a manifest that install would have refused.
//
// Stability rules, in order:
//   - no %VAR% / $VAR in the spelling, and no LEADING ~ (that is what
//     "movable" means here; a mid-component ~ is an 8.3 alias, not a variable);
//   - the spelling is absolute as written (a relative exe depends on the cwd);
//   - C26's expansion step must NOT have moved it (Result.Actable);
//   - it is a local drive path, not a UNC share (a share name is not a pin);
//   - it does not live under a TEMP/TMP/TMPDIR root;
//   - it names an executable image (.exe/.com), not a script a shell would run.
func resolveExe(spelling string, exceptions []string) (string, error) {
	if strings.TrimSpace(spelling) == "" {
		return "", fmt.Errorf("%w: empty exe", ErrExeUnstable)
	}
	for _, marker := range unstableExeMarkers {
		if strings.Contains(spelling, marker) {
			return "", fmt.Errorf("%w: %q contains %q (variable/home-relative paths are refused before C26 expands them; a movable executable cannot be hash-pinned)",
				ErrExeUnstable, spelling, marker)
		}
	}
	if isHomeRelativeSpelling(spelling) {
		return "", fmt.Errorf("%w: %q is home-relative (a leading ~ resolves against whoever runs the manifest; a movable executable cannot be hash-pinned)",
			ErrExeUnstable, spelling)
	}
	if !isAbsSpelling(spelling) {
		return "", fmt.Errorf("%w: %q is not an absolute path", ErrExeUnstable, spelling)
	}
	res, err := risk.Resolve(spelling, exceptions)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrExeUnstable, err)
	}
	canonical, err := res.Actable()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrExeUnstable, err)
	}
	if strings.HasPrefix(canonical, `\\`) {
		return "", fmt.Errorf("%w: %q resolves to a UNC path %q (a share name is not a local pin)", ErrExeUnstable, spelling, canonical)
	}
	if vol := filepath.VolumeName(canonical); runtime.GOOS == "windows" && vol == "" {
		return "", fmt.Errorf("%w: %q resolves to a drive-less path %q (a volume-relative spelling is resolved against an unpredictable current directory)", ErrExeUnstable, spelling, canonical)
	}
	if env, root, inTemp := underTempRoot(canonical); inTemp {
		return "", fmt.Errorf("%w: %q resolves under %s (%q); a TEMP directory is recycled", ErrExeUnstable, spelling, env, root)
	}
	if !CommandExeAllowed(canonical) {
		return "", fmt.Errorf("%w: %q is not an executable image (%s); scripts would need a shell, which D46 never uses",
			ErrExeUnstable, canonical, strings.Join(commandExeExts, "/"))
	}
	return canonical, nil
}

// isAbsSpelling reports whether p is absolute on the native filesystem
// without any C26 step (the manifest string is judged as written).
//
// A drive-qualified path is absolute whether or not the next character is a
// separator: on Windows both `C:\x\cli.exe` and `C:/x/cli.exe` are absolute,
// and a drive-relative spelling (`C:x`, no separator at all) is caught below
// by filepath.IsAbs, which is the platform's own verdict. The check is
// deliberately not written as "VolumeName != \"\" && root separator", the
// shape that once refused every forward-slash Windows path D46's own spec
// example uses (SPEC-07 §4: `exe = "C:/Users/me/bin/lark-cli.exe"`).
func isAbsSpelling(p string) bool {
	if filepath.VolumeName(p) != "" {
		return filepath.IsAbs(p)
	}
	// Posix absolute spelling (`/tmp/cli.exe`): absolute where POSIX paths are
	// the native spelling, so D46 can pin the exact bytes of a real file there.
	// On a WINDOWS build the same spelling is drive-root-relative (the volume
	// letter would come from the current directory), so it stays refused - the
	// pre-triage verdict on windows is byte-for-byte unchanged.
	if runtime.GOOS == "windows" {
		return false
	}
	return strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "//")
}

// underTempRoot reports whether p sits inside a TEMP/TMP/TMPDIR root, and which
// variable named that root. Both the validator (C14, on the spelling as
// written) and resolveExe (install, on the C26-canonical path) call it, so "a
// TEMP path is refused" is one rule expressed once rather than two guesses that
// can drift apart. The root spelling is the env value AS ANSWERED - on GitHub's
// windows runners TEMP is the 8.3 short form (...\RUNNER~1\...), and because
// the check is a prefix comparison against a root this process got from the
// same environment as the exe spelling, the short form matches a short
// spelling; an 8.3-alias root would only evade this check on a host that
// spells the exe long and TEMP short, which the C26-resolved canonical side
// absorbs by construction (resolveExe compares one resolution against the
// same host's own answer).
//
// On POSIX with NO TEMP/TMP/TMPDIR answer at all, the recycled-tree rule still
// has a subject: os.TempDir()'s own default (/tmp) is the platform's recycled
// scratch root, and an exe pinned there is exactly the red-team case the rule
// exists for (t.TempDir() itself lands there when TMPDIR is unset). The /tmp
// default is consulted only when all three variables are empty, so a host that
// points TMPDIR elsewhere keeps answering through its own variable.
func underTempRoot(p string) (envname, root string, ok bool) {
	for _, env := range tempDirEnvnames {
		r := strings.TrimSpace(os.Getenv(env))
		if r == "" {
			continue
		}
		if isWithinDir(p, r) {
			return env, r, true
		}
	}
	if runtime.GOOS != "windows" {
		for _, env := range tempDirEnvnames {
			if strings.TrimSpace(os.Getenv(env)) != "" {
				return "", "", false // a POSIX host answered; its answer is the rule's subject, /tmp is not
			}
		}
		if isWithinDir(p, "/tmp") {
			return "os.TempDir default", "/tmp", true
		}
	}
	return "", "", false
}

// isWithinDir reports whether canonical sits inside dir (or is dir). Both
// sides are C26 output or an OS-provided root; the comparison is
// case-insensitive because Windows filesystems are, and component-bounded so
// "C:\Temp2" is not inside "C:\Temp".
//
// Each side is compared in BOTH of the spellings C26 distinguishes (the same
// two forms risk's pathForms separates): the spelling AS ANSWERED (raw) and
// the handle-resolved real path when the OS can answer one. A hosted windows
// runner answers TEMP as the 8.3 short form (C:\Users\RUNNER~1\...) while a
// resolved path comes back long (C:\Users\runneradmin\...), so a single-form
// comparison misses exactly there (risk/pathresolver.go's anchor-spelling
// note); a containment hit on EITHER paired form is enough, because both
// forms of one side name the same tree on the same host. This predicate feeds
// refusal rules only, so the extra hits are strictness, never a pass.
func isWithinDir(canonical, dir string) bool {
	dirRes, err := risk.Resolve(dir, nil)
	if err != nil {
		return false
	}
	cRes, err := risk.Resolve(canonical, nil)
	if err != nil {
		return false
	}
	roots := []string{dirRes.Canonical, dir}
	cands := []string{cRes.Canonical, canonical}
	for _, c := range cands {
		c = strings.ToLower(c)
		if c == "" {
			continue
		}
		for _, r := range roots {
			r = strings.ToLower(strings.TrimRight(r, `\/`))
			if r == "" {
				continue
			}
			if c == r || strings.HasPrefix(c, r+`\`) || strings.HasPrefix(c, r+"/") {
				return true
			}
		}
	}
	return false
}

// ResolveExeForRun re-runs the install-time C26 stability rules over a
// manifest's exe at execution time. The runtime uses it for a pin record that
// carries only a hash (the single-tool legacy form); a pin with a path uses the
// path, because that is the file whose bytes were hashed.
func ResolveExeForRun(t ToolDef) (string, error) { return resolveExe(t.Exe, nil) }

// HashFileBytes returns the "sha256:<hex>" pin of one file's bytes, plus the
// byte count. It is the ONE hashing implementation behind both exe_hash values
// and the runtime recheck.
func HashFileBytes(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, fmt.Errorf("plugin: open %s: %w", path, err)
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, fmt.Errorf("plugin: read %s: %w", path, err)
	}
	return HashAlgorithm + hex.EncodeToString(h.Sum(nil)), n, nil
}

// PinCommandExe is the install-time step for one command tool: resolve the exe
// through C26, refuse anything unstable, hash the bytes, and compare with the
// manifest's own exe_hash when it declared one.
//
// The manifest's exe_hash is optional at PARSE time (the loader cannot hash a
// file it must not open) and mandatory at PIN time: this function produces the
// value that plugin_state stores and every later call rechecks.
func PinCommandExe(t ToolDef, exceptions []string) (ExePin, error) {
	if !t.IsCommand() {
		return ExePin{}, fmt.Errorf("plugin: tool %q is not a command tool", t.Name)
	}
	canonical, err := resolveExe(t.Exe, exceptions)
	if err != nil {
		return ExePin{}, fmt.Errorf("plugin: tool %q: %w", t.Name, err)
	}
	hash, size, err := HashFileBytes(canonical)
	if err != nil {
		return ExePin{}, fmt.Errorf("plugin: tool %q: %w", t.Name, err)
	}
	if t.ExeHash != "" {
		want, err := normalizePinnedHash(t.ExeHash)
		if err != nil {
			return ExePin{}, fmt.Errorf("plugin: tool %q: %w", t.Name, err)
		}
		if !strings.EqualFold(want, hash) {
			return ExePin{}, fmt.Errorf("%w: tool %q declared %s, bytes hash to %s",
				ErrExeSwap, t.Name, want, hash)
		}
	}
	return ExePin{Path: canonical, Hash: hash, Size: size}, nil
}

// PinManifestExe is the plugin-level install step: one Pin per command tool,
// keyed by tool name. It returns nil for a manifest without command tools, so
// the installer's Tier-1 path is unchanged.
//
// exe_hash is stored per TOOL, not per plugin, because one plugin may wrap more
// than one CLI; pinning only the first would leave the rest unverified.
func PinManifestExe(m Manifest, exceptions []string) (map[string]ExePin, error) {
	var pins map[string]ExePin
	for _, t := range m.Tools {
		if !t.IsCommand() {
			continue
		}
		pin, err := PinCommandExe(t, exceptions)
		if err != nil {
			return nil, err
		}
		if pins == nil {
			pins = make(map[string]ExePin, len(m.Tools))
		}
		pins[t.Name] = pin
	}
	return pins, nil
}

// ExePinEnvelopePrefix marks the serialized multi-tool pin form written to
// plugin_state. A single-tool plugin stores a bare "sha256:<hex>" (the shape
// ticket 50 reserved); several command tools need their names attached, so the
// value becomes this prefix plus "tool=path@hash" lines. Both forms are read
// back by ParseExeHashField, and neither is a manifest hash.
const ExePinEnvelopePrefix = "d46-exe-v1"

// EncodeExePins serializes per-tool pins into the plugin_state value. One
// command tool with no newline or ';' in its name collapses to the bare hash,
// which keeps the common case readable in the database.
func EncodeExePins(pins map[string]ExePin) string {
	if len(pins) == 0 {
		return ""
	}
	if len(pins) == 1 {
		for _, p := range pins {
			return p.Hash
		}
	}
	names := sortedKeys(pins)
	var b strings.Builder
	b.WriteString(ExePinEnvelopePrefix)
	for _, name := range names {
		p := pins[name]
		b.WriteString("\n")
		b.WriteString(name + "=" + p.Path + "@" + p.Hash)
	}
	return b.String()
}

// ParseExeHashField reads back what EncodeExePins wrote (or the legacy bare
// hash) into per-tool pins. An empty value yields no pins, which the runtime
// refuses rather than treating as "nothing to verify".
func ParseExeHashField(raw string) (map[string]ExePin, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if !strings.HasPrefix(raw, ExePinEnvelopePrefix) {
		if _, err := normalizePinnedHash(raw); err != nil {
			return nil, fmt.Errorf("plugin: exe_hash field: %w", err)
		}
		return map[string]ExePin{"": {Hash: raw}}, nil
	}
	lines := strings.Split(raw, "\n")
	pins := make(map[string]ExePin, len(lines)-1)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, rest, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("plugin: exe_hash field: malformed entry %q", line)
		}
		at := strings.LastIndex(rest, "@")
		if at < 0 {
			return nil, fmt.Errorf("plugin: exe_hash field: malformed entry %q", line)
		}
		pins[name] = ExePin{Path: rest[:at], Hash: rest[at+1:]}
	}
	if len(pins) == 0 {
		return nil, fmt.Errorf("plugin: exe_hash field: %q carries no pins", ExePinEnvelopePrefix)
	}
	return pins, nil
}

// Verify re-hashes the pinned executable and compares it with the pin. It is
// the every-execution gate: mismatch or a missing pin is the error, and the
// caller alerts before it returns the refusal.
func (p ExePin) Verify() error {
	if p.Hash == "" {
		return ErrNoPin
	}
	want, err := normalizePinnedHash(p.Hash)
	if err != nil {
		return err
	}
	hash, _, err := HashFileBytes(p.Path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExeSwap, err)
	}
	if !strings.EqualFold(want, hash) {
		return fmt.Errorf("%w: pinned %s, computed %s (%s)", ErrExeSwap, want, hash, p.Path)
	}
	return nil
}

// VerifyExeAlert is Pin.Verify plus the alert contract: a mismatch or missing
// pin is reported to the alerter BEFORE the refusal leaves, so a caller that
// drops the error cannot drop the evidence (same rule as the manifest tamper
// gate in tamper.go). It returns the current hash for the alert record.
func VerifyExeAlert(p ExePin, pluginID, toolName string, alerter Alerter) error {
	if alerter == nil {
		alerter = LogAlerter{}
	}
	err := p.Verify()
	if err == nil {
		return nil
	}
	actual := ""
	if h, _, herr := HashFileBytes(p.Path); herr == nil {
		actual = h
	}
	alerter.AlertTamper(TamperAlert{
		Kind:     "exe",
		PluginID: pluginID,
		Path:     p.Path,
		Expected: p.Hash,
		Actual:   actual,
		Reason:   fmt.Sprintf("tool %s: %v", toolName, err),
	})
	return err
}

// placeholderPattern finds every {{name}} reference in one argv template
// element. Only the double-brace form is a placeholder: a single "{name}" is
// passed through as literal text, because a CLI is free to use braces itself.
var placeholderPattern = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

// ArgvPlaceholders returns the placeholder names inside one template element.
func ArgvPlaceholders(element string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(element, -1)
	if len(matches) == 0 {
		return nil
	}
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m[1])
	}
	return out
}

// substitutedPlaceholders is the union over a whole template, in first-seen
// order (stable, so error text and the "declared but unused" check agree).
func substitutedPlaceholders(argvTemplate []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, el := range argvTemplate {
		for _, name := range ArgvPlaceholders(el) {
			if seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

// validateArgvTemplate checks the template as a VECTOR: one empty string is the
// classic string-concatenation slip, and element 0 is the exe, never text.
func validateArgvTemplate(argv []string, path string) ValidationErrors {
	var errs ValidationErrors
	if len(argv) == 0 {
		return ValidationErrors{{
			Field:  path,
			Reason: "required for a command tool: the argv VECTOR handed to process creation (D46 never concatenates a command string)",
		}}
	}
	for i, el := range argv {
		at := fmt.Sprintf("%s[%d]", path, i)
		if el == "" {
			errs = append(errs, ValidationError{Field: at, Reason: "must not be an empty string (an empty argv slot is how a string command gets smuggled in)"})
		}
		if hasControl(el) {
			errs = append(errs, ValidationError{Field: at, Reason: "contains control characters (newlines included): argv elements are single arguments"})
		}
	}
	if i := strings.Index(argv[0], "{{"); i >= 0 {
		errs = append(errs, ValidationError{
			Field:  path + "[0]",
			Reason: fmt.Sprintf("argv[0] names the pinned exe and must be a literal, found a placeholder at offset %d (D46: a command tool never runs an exe that is not fixed in the manifest)", i),
		})
	}
	return errs
}

// validatePlaceholderSchemas is the "parameters must be a JSON Schema per
// placeholder" rule. It reads the properties table written in the manifest and
// requires every placeholder to have an entry there.
func validatePlaceholderSchemas(t *ToolDef, path string) ValidationErrors {
	props, propErr := schemaProperties(t.Parameters)
	if propErr != nil {
		return ValidationErrors{{Field: path + ".parameters", Reason: propErr.Error()}}
	}
	var errs ValidationErrors
	for _, name := range substitutedPlaceholders(t.ArgvTemplate) {
		if _, ok := props[name]; !ok {
			errs = append(errs, ValidationError{
				Field: fmt.Sprintf("%s.parameters.properties.%s", path, name),
				Reason: fmt.Sprintf("argv_template references {{%s}} but no schema is declared for it: every placeholder needs its own type/enum/range/maxLength entry (D46), and an undeclared parameter would reach argv unjudged",
					name),
			})
		}
	}
	return errs
}

// validateExtract checks the declarative extractor. Raw stdout never enters the
// context, so a command tool without a usable extractor cannot answer at all.
func validateExtract(extract map[string]any, path string) ValidationErrors {
	if len(extract) == 0 {
		return ValidationErrors{{
			Field:  path,
			Reason: "required: declare how the output is extracted (kind = \"json\" with a path, or kind = \"regex\" with a pattern); raw stdout never enters the context (D46)",
		}}
	}
	ext, err := newExtractor(extract)
	if err != nil {
		return ValidationErrors{{Field: path, Reason: err.Error()}}
	}
	_ = ext
	return nil
}

// ---------------------------------------------------------------------------
// Parameter schema access (shared by the validator, the slot's C1 Parameters()
// and the runtime's pre-argv screen, so the three can never disagree).
// ---------------------------------------------------------------------------

// schemaProperties reads parameters.properties as a name -> schema table.
func schemaProperties(parameters map[string]any) (map[string]map[string]any, error) {
	out := map[string]map[string]any{}
	if parameters == nil {
		return out, nil
	}
	raw, ok := parameters["properties"]
	if !ok || raw == nil {
		return out, nil
	}
	table, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("properties must be a table of JSON Schemas (per-parameter)")
	}
	for name, v := range table {
		sub, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("properties.%s must be a JSON Schema table", name)
		}
		out[name] = sub
	}
	return out, nil
}

// requiredParameters reads parameters.required. A parameter that is required
// by the schema AND used by the template must be present on every call; the
// runtime enforces it before it builds argv.
func requiredParameters(parameters map[string]any) []string {
	raw, ok := parameters["required"]
	if !ok {
		return nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// ParameterSchema returns the JSON Schema a command tool advertises to the
// model: the manifest's parameters table, restricted to the placeholders the
// argv template actually substitutes.
//
// Why restricted: the parenthetical form the template uses means the
// enforcement lives here, not in the model-facing schema. A model that sees a
// property the template never substitutes would believe it changed the call,
// and would be wrong silently.
func ParameterSchema(t ToolDef) map[string]any {
	if !t.IsCommand() {
		return t.Parameters
	}
	props, err := schemaProperties(t.Parameters)
	if err != nil {
		return t.Parameters
	}
	out := make(map[string]any, len(t.Parameters))
	for k, v := range t.Parameters {
		out[k] = v
	}
	used := substitutedPlaceholders(t.ArgvTemplate)
	usedSet := make(map[string]bool, len(used))
	for _, name := range used {
		usedSet[name] = true
	}
	narrowed := make(map[string]any, len(used))
	for _, name := range used {
		if sub, ok := props[name]; ok {
			narrowed[name] = sub
		}
	}
	out["properties"] = narrowed
	if req := requiredParameters(t.Parameters); len(req) > 0 {
		kept := make([]any, 0, len(req))
		for _, name := range req {
			if usedSet[name] {
				kept = append(kept, name)
			}
		}
		out["required"] = kept
	}
	return out
}
