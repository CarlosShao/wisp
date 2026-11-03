package plugin

import (
	"fmt"
	"regexp"
	"strings"
)

// envNamePattern is the shape of an environment variable NAME (no '=',
// no path, nothing a shell would interpret).
var envNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError is one C14 schema finding: a dotted field path plus the
// human-readable reason. Reasons name the contract clause they come from so
// a plugin author can look it up.
type ValidationError struct {
	Field  string // dot-path; "(document)" for parse-level errors
	Reason string
}

// ValidationErrors is the full finding list. It implements error, and the
// loader surfaces every finding at once (mainstream harness behaviour: fix
// the file in one pass instead of one error per rebuild).
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "plugin: no validation errors"
	}
	var b strings.Builder
	b.WriteString("plugin: manifest failed C14 validation:")
	for _, e := range v {
		b.WriteString("\n  - ")
		b.WriteString(e.Field)
		b.WriteString(": ")
		b.WriteString(e.Reason)
	}
	return b.String()
}

// Fields returns the offending field paths (stable order, as reported).
func (v ValidationErrors) Fields() []string {
	out := make([]string, 0, len(v))
	for _, e := range v {
		out = append(out, e.Field)
	}
	return out
}

// Validate runs the per-field C14 rules over a decoded manifest. It is the
// second half of AC1 (ParseManifest does the closed-schema key audit first);
// it returns nil when the manifest is loadable.
func Validate(m *Manifest) ValidationErrors {
	var errs ValidationErrors

	if m.SchemaVersion != CurrentSchemaVersion {
		errs = append(errs, ValidationError{
			Field:  "schema_version",
			Reason: fmt.Sprintf("unsupported schema_version %d (this loader implements %d)", m.SchemaVersion, CurrentSchemaVersion),
		})
	}

	// id becomes the on-disk directory name plugins\<id>\ (SPEC-07 §6), so the
	// charset is restricted to a single safe path segment.
	switch {
	case m.ID == "":
		errs = append(errs, ValidationError{Field: "id", Reason: "required (SPEC-07 §4)"})
	case m.ID == "." || m.ID == "..":
		errs = append(errs, ValidationError{Field: "id", Reason: "reserved path segment"})
	case !idPattern.MatchString(m.ID):
		errs = append(errs, ValidationError{
			Field:  "id",
			Reason: "must match [A-Za-z0-9._-]+ (it is used as the plugins\\<id> directory name; separators and '.'/'..' are rejected)",
		})
	}

	// Presence/cleanliness is this cell (AC1); the range grammar is checked
	// here too so an unreadable host_api never reaches the C16 gate, whose
	// verdict (load vs refuse) is ticket 50 AC2 — hostapi.go.
	switch {
	case m.HostAPI == "":
		errs = append(errs, ValidationError{Field: "host_api", Reason: "required (C16 range, e.g. \"^1.0\")"})
	case hasControl(m.HostAPI):
		errs = append(errs, ValidationError{Field: "host_api", Reason: "contains control characters"})
	case func() bool { _, err := ParseCaretRange(m.HostAPI); return err != nil }():
		_, err := ParseCaretRange(m.HostAPI)
		errs = append(errs, ValidationError{Field: "host_api", Reason: fmt.Sprintf("invalid C16 range: %v", err)})
	}

	switch {
	case m.Name == "":
		errs = append(errs, ValidationError{Field: "name", Reason: "required"})
	case hasControl(m.Name):
		errs = append(errs, ValidationError{Field: "name", Reason: "contains control characters"})
	}

	switch {
	case m.Version == "":
		errs = append(errs, ValidationError{Field: "version", Reason: "required"})
	case hasControl(m.Version):
		errs = append(errs, ValidationError{Field: "version", Reason: "contains control characters"})
	}

	for i := range m.Tools {
		errs = append(errs, validateTool(&m.Tools[i], i)...)
	}
	errs = append(errs, validateDuplicateToolNames(m.Tools)...)
	return errs
}

// validateDuplicateToolNames enforces C1's uniqueness inside one manifest. Two
// [[tools]] entries with the same name would silently collapse in the slot's
// byName map, and the one that lost is the one nobody would notice missing.
func validateDuplicateToolNames(tools []ToolDef) ValidationErrors {
	var errs ValidationErrors
	seen := make(map[string]int, len(tools))
	for i := range tools {
		name := tools[i].Name
		if name == "" {
			continue
		}
		if first, dup := seen[name]; dup {
			errs = append(errs, ValidationError{
				Field:  fmt.Sprintf("tools[%d].name", i),
				Reason: fmt.Sprintf("duplicate tool name %q (already declared at tools[%d]; C1 requires global uniqueness)", name, first),
			})
			continue
		}
		seen[name] = i
	}
	return errs
}

func validateTool(t *ToolDef, i int) ValidationErrors {
	path := fmt.Sprintf("tools[%d]", i)
	var errs ValidationErrors

	// C1: tool names are globally unique and shaped 'namespace.action'.
	switch {
	case t.Name == "":
		errs = append(errs, ValidationError{Field: path + ".name", Reason: "required"})
	case hasControl(t.Name):
		errs = append(errs, ValidationError{Field: path + ".name", Reason: "contains control characters"})
	case !toolNamePair.MatchString(t.Name):
		errs = append(errs, ValidationError{
			Field:  path + ".name",
			Reason: "must be 'namespace.action' (C1): lowercase letters/digits/'_'/'-', dot-separated, at least two segments",
		})
	}

	switch {
	case t.Description == "":
		errs = append(errs, ValidationError{Field: path + ".description", Reason: "required"})
	case hasControl(t.Description):
		errs = append(errs, ValidationError{Field: path + ".description", Reason: "contains control characters"})
	}

	// C2: the declaration is only the R1 lower bound; anything that is not a
	// known level is rejected instead of defaulting to the safest-looking one.
	switch t.Risk {
	case RiskL0, RiskL1, RiskL2:
	default:
		errs = append(errs, ValidationError{
			Field:  path + ".risk",
			Reason: fmt.Sprintf("must be one of %s/%s/%s (C2); got %q", RiskL0, RiskL1, RiskL2, t.Risk),
		})
	}

	errs = append(errs, validateCapabilities(t.Capabilities, path)...)
	errs = append(errs, validateNetAllowlist(t, path)...)

	// parameters is a JSON Schema object (C14). Its inner keys are another
	// grammar, so the manifest key auditor does not descend into it.
	switch {
	case t.Parameters == nil:
		errs = append(errs, ValidationError{
			Field:  path + ".parameters",
			Reason: "required: a JSON Schema table mapping call arguments to the tool (C14)",
		})
	case len(t.Parameters) == 0:
		errs = append(errs, ValidationError{
			Field:  path + ".parameters",
			Reason: "must not be an empty table (declare the schema, e.g. type/properties)",
		})
	}
	if typ, ok := t.Parameters["type"]; ok {
		if s, ok := typ.(string); ok && s == "" {
			errs = append(errs, ValidationError{
				Field:  path + ".parameters.type",
				Reason: "must not be an empty string",
			})
		}
	}

	// kind selects which of the two shapes this entry is. An unrecognized
	// spelling is refused: defaulting it to "declarative" would leave the D46
	// fields of a mis-typed command tool unread, and defaulting it to
	// "command" would demand an exe nobody wrote.
	var command bool
	switch t.Kind {
	case "", ToolKindDeclarative:
	case ToolKindCommand:
		command = true
	default:
		errs = append(errs, ValidationError{
			Field:  path + ".kind",
			Reason: fmt.Sprintf("must be %q or %q (C14); got %q", ToolKindDeclarative, ToolKindCommand, t.Kind),
		})
	}

	errs = append(errs, validateIrreversible(t, path)...)
	if command {
		errs = append(errs, validateCommand(t, path)...)
	}
	return errs
}

// validateCommand is the D46 half of C14 (§16.5.6): the fields that only a
// kind = "command" tool carries, each checked against the constraint that
// makes the execution channel safe rather than merely convenient.
func validateCommand(t *ToolDef, path string) ValidationErrors {
	var errs ValidationErrors
	errs = append(errs, validateExe(t.Exe, path+".exe")...)
	errs = append(errs, validateArgvTemplate(t.ArgvTemplate, path+".argv_template")...)
	if _, perr := schemaProperties(t.Parameters); perr == nil {
		errs = append(errs, validatePlaceholderSchemas(t, path)...)
	}
	errs = append(errs, validateExtract(t.Extract, path+".extract")...)

	// timeout_ms: 0 means "use the D46 default (15000)". A negative value is
	// nonsense, and an explicit value gets the same floor/ceiling the runtime
	// applies, so a manifest cannot declare a budget the executor would then
	// silently change.
	switch {
	case t.TimeoutMs == 0:
	case t.TimeoutMs < MinCommandTimeoutMS:
		errs = append(errs, ValidationError{
			Field:  path + ".timeout_ms",
			Reason: fmt.Sprintf("must be 0 (default %d) or at least %d ms; got %d", DefaultCommandTimeoutMS, MinCommandTimeoutMS, t.TimeoutMs),
		})
	case t.TimeoutMs > MaxCommandTimeoutMS:
		errs = append(errs, ValidationError{
			Field:  path + ".timeout_ms",
			Reason: fmt.Sprintf("must be at most %d ms; got %d", MaxCommandTimeoutMS, t.TimeoutMs),
		})
	}

	seen := make(map[string]bool, len(t.EnvAllowlist))
	for j, name := range t.EnvAllowlist {
		at := fmt.Sprintf("%s.env_allowlist[%d]", path, j)
		switch {
		case name == "":
			errs = append(errs, ValidationError{Field: at, Reason: "must not be an empty string"})
		case seen[name]:
			errs = append(errs, ValidationError{Field: at, Reason: fmt.Sprintf("duplicate entry %q", name)})
		case !envNamePattern.MatchString(name):
			errs = append(errs, ValidationError{
				Field:  at,
				Reason: fmt.Sprintf("%q is not an environment variable name ([A-Za-z_][A-Za-z0-9_]*)", name),
			})
		}
		seen[name] = true
	}
	return errs
}

func validateIrreversible(t *ToolDef, path string) ValidationErrors {
	var errs ValidationErrors
	seen := make(map[string]bool, len(t.Irreversible))
	for j, class := range t.Irreversible {
		at := fmt.Sprintf("%s.irreversible[%d]", path, j)
		switch {
		case class == "":
			errs = append(errs, ValidationError{Field: at, Reason: "must not be an empty string"})
		case seen[class]:
			errs = append(errs, ValidationError{Field: at, Reason: fmt.Sprintf("duplicate class %q", class)})
		case !containsString(IrreversibleClasses, class):
			errs = append(errs, ValidationError{
				Field:  at,
				Reason: fmt.Sprintf("%q is not an R8 class (%s)", class, strings.Join(IrreversibleClasses, "|")),
			})
		}
		seen[class] = true
	}
	return errs
}

func containsString(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// validateExe is the C14 shape check on `exe`, the cheap half of D46's
// "absolute + stable path" rule. The expensive half (C26 resolution, TEMP
// detection, the executable-image check and the hash) runs at install time in
// PinCommandExe, which refuses the same spellings; keeping the shape check here
// too means a manifest with `%TEMP%\x.exe` is refused by the loader as well, so
// the red-team case fails even for a leg that never reaches the pin.
func validateExe(exe, path string) ValidationErrors {
	switch {
	case strings.TrimSpace(exe) == "":
		return ValidationErrors{{Field: path, Reason: "required for a command tool (D46: the absolute path of the wrapped CLI)"}}
	}
	for _, marker := range unstableExeMarkers {
		if strings.Contains(exe, marker) {
			return ValidationErrors{{
				Field:  path,
				Reason: fmt.Sprintf("%q contains %q: D46 refuses a variable or home-relative exe because a movable executable cannot be hash-pinned (PLAN.md §16.5.6)", exe, marker),
			}}
		}
	}
	var errs ValidationErrors
	if !isAbsSpelling(exe) {
		errs = append(errs, ValidationError{
			Field:  path,
			Reason: fmt.Sprintf("%q is not absolute: a relative exe depends on the working directory, so its hash proves nothing about what runs", exe),
		})
	}
	// The dangerous ABSOLUTE spelling is refused too, and the order is about a
	// readable message rather than about severity: a path under
	// `C:/.../Temp\cli.exe` is absolute, so a test asserting the TEMP marker
	// would pass on the absoluteness refusal alone and never prove the
	// stability rule. Being inside a TEMP root is the stronger statement, so it
	// is named first.
	if env, root, inTemp := underTempRoot(exe); inTemp {
		errs = append(errs, ValidationError{
			Field:  path,
			Reason: fmt.Sprintf("%q is under %s (%q): a TEMP directory is recycled, so an executable there cannot be hash-pinned (PLAN.md §16.5.6 red-team case)", exe, env, root),
		})
	}
	if strings.HasPrefix(exe, `\\`) || strings.HasPrefix(exe, "//") {
		errs = append(errs, ValidationError{
			Field:  path,
			Reason: fmt.Sprintf("%q is a UNC share: a share name is not a local pin", exe),
		})
	}
	if !CommandExeAllowed(exe) {
		errs = append(errs, ValidationError{
			Field:  path,
			Reason: fmt.Sprintf("%q is not an executable image (%s): a script would need a shell, and D46 never uses one", exe, strings.Join(commandExeExts, "/")),
		})
	}
	return errs
}

func validateCapabilities(caps []string, path string) ValidationErrors {
	var errs ValidationErrors
	seen := make(map[string]bool, len(caps))
	for j, c := range caps {
		at := fmt.Sprintf("%s.capabilities[%d]", path, j)
		switch {
		case c == "":
			errs = append(errs, ValidationError{Field: at, Reason: "must not be an empty string"})
		case seen[c]:
			errs = append(errs, ValidationError{Field: at, Reason: fmt.Sprintf("duplicate capability %q", c)})
		}
		seen[c] = true
		if c != "" && !CapabilityDeclared(Capabilities, c) {
			errs = append(errs, ValidationError{
				Field:  at,
				Reason: fmt.Sprintf("capability %q is not in the frozen C3 set (%s)", c, strings.Join(Capabilities, "|")),
			})
		}
	}
	return errs
}

func validateNetAllowlist(t *ToolDef, path string) ValidationErrors {
	var errs ValidationErrors
	// Contract direction (SPEC-07 §4): a tool declaring "net" must ship a
	// non-empty target-domain whitelist.
	if CapabilityDeclared(t.Capabilities, "net") && len(t.NetAllowlist) == 0 {
		errs = append(errs, ValidationError{
			Field:  path + ".net_allowlist",
			Reason: "required and non-empty when the tool declares the \"net\" capability (SPEC-07 §4: net 能力必须带目标域名白名单)",
		})
	}
	seen := make(map[string]bool, len(t.NetAllowlist))
	for j, d := range t.NetAllowlist {
		at := fmt.Sprintf("%s.net_allowlist[%d]", path, j)
		switch {
		case d == "":
			errs = append(errs, ValidationError{Field: at, Reason: "must not be an empty string"})
		case seen[d]:
			errs = append(errs, ValidationError{Field: at, Reason: fmt.Sprintf("duplicate entry %q", d)})
		case !domainPattern.MatchString(d):
			errs = append(errs, ValidationError{
				Field:  at,
				Reason: fmt.Sprintf("%q is not a bare hostname: no scheme, port, path, credentials or wildcard", d),
			})
		}
		seen[d] = true
	}
	return errs
}
