package plugin

import (
	"fmt"
	"regexp"
	"sort"

	toml "github.com/pelletier/go-toml/v2"
)

// C14 plugin manifest — Tier-1 declarative base (SPEC-07 §4, PLAN.md contract
// table C14, ticket 50 AC1).
//
// Field table authority:
//   - top level: schema_version, id, host_api, name, version, tools[]
//   - tools[]:   name, description, risk, capabilities, net_allowlist, parameters
//   - tools[] where kind = "command" (D46, ticket 52): additionally
//     exe, exe_hash, argv_template, extract, timeout_ms, env_allowlist
//     and the two R8 declaration knobs irreversible / taint_output
//
// Deferred fields are NOT accepted by this slice; the key auditor names the
// owning ticket instead of pretending they never existed:
//   - entry                                    -> ticket 51 (Tier-2 goja)
//
// "Deferred" is a per-field fact, not a per-ticket one: the D46 fields are
// accepted on a kind = "command" tool only. On a declarative tool the auditor
// names the ticket that owns them instead of reading them and doing nothing.
//
// CurrentSchemaVersion is the only schema_version this loader understands
// (SPEC-07 §4 example: schema_version = 1).
const CurrentSchemaVersion = 1

// ToolKindDeclarative is the default tool kind: a Tier-1 tool whose behaviour
// is declared (host capabilities today, the goja entry of ticket 51 later).
const ToolKindDeclarative = "declarative"

// ToolKindCommand is the D46 kind: a wrapper around ONE pinned local CLI
// (SPEC-07 §4 "D46 command 型", PLAN.md §16.5.6).
const ToolKindCommand = "command"

// Manifest is a parsed, C14-validated plugin manifest.
type Manifest struct {
	SchemaVersion int       `toml:"schema_version"`
	ID            string    `toml:"id"`
	HostAPI       string    `toml:"host_api"`
	Name          string    `toml:"name"`
	Version       string    `toml:"version"`
	Tools         []ToolDef `toml:"tools"`
}

// ToolDef is one Tier-1 tool.
//
// The declaration face (name/description/risk/capabilities/net_allowlist/
// parameters) is common to both kinds. The D46 fields are only legal on a
// kind = "command" tool; the validator (validator.go) refuses them anywhere
// else rather than letting them sit unread, because a field the host ignores
// is a promise the host is not keeping.
type ToolDef struct {
	Name         string         `toml:"name"`
	Description  string         `toml:"description"`
	Risk         string         `toml:"risk"`
	Capabilities []string       `toml:"capabilities"`
	NetAllowlist []string       `toml:"net_allowlist"`
	Parameters   map[string]any `toml:"parameters"`
	Kind         string         `toml:"kind"`

	// Exe is the absolute, stable path of the wrapped CLI (D46). A manifest
	// that names it relative, through %VARS%/~, or under a TEMP directory is
	// refused at install: a movable executable cannot be hash-pinned.
	Exe string `toml:"exe"`
	// ExeHash is the "sha256:<hex>" pin. It is optional IN THE MANIFEST: the
	// install leg computes it from the bytes on disk when it is absent (the
	// ticket's "install computes hash"), and a value that is present must
	// match those bytes or the install is refused as a local swap.
	ExeHash string `toml:"exe_hash"`
	// ArgvTemplate is the argv VECTOR handed to process creation. Element 0
	// must be the exe itself; "{{name}}" placeholders are substituted from the
	// schema-validated parameters and never through a shell.
	ArgvTemplate []string `toml:"argv_template"`
	// Extract is the declarative output extractor (json path or regex). It is
	// required: raw stdout never enters the model context (D46).
	Extract map[string]any `toml:"extract"`
	// TimeoutMs is the per-call budget; 0 means DefaultCommandTimeoutMS.
	TimeoutMs int `toml:"timeout_ms"`
	// EnvAllowlist names the environment variables the child may see. The
	// child NEVER inherits the host environment, so the contract's base set
	// (PATH trimmed, HOME, USERPROFILE, LANG, SystemRoot) is always present
	// even when this list is empty; entries here are additions.
	EnvAllowlist []string `toml:"env_allowlist"`

	// Irreversible declares the R8 classes a call of this tool amounts to
	// ("send", "delete", "overwrite", "power", "close-window"). It is a
	// DECLARATION, never a verdict: C19/R8 fuses it with everything else and
	// takes the max, so a declaration can only raise the level. It exists
	// because C19 cannot read intent out of a CLI's argv template - the host
	// knows that "write" is a write, but only the plugin author knows that
	// "lark-cli im send" sends a message to another human.
	Irreversible []string `toml:"irreversible"`
	// TaintOutput marks the command's stdout as an untrusted-content source
	// for C25 (the equivalent of the builtins whose result is a sensitive
	// source), so a fragment of what the CLI printed cannot later leave
	// through an exfil channel unnoticed. The extracted output is always
	// treated as data, never as instructions (D30 layer 4).
	TaintOutput bool `toml:"taint_output"`
}

// IsCommand reports whether this tool is a D46 command wrapper.
func (t ToolDef) IsCommand() bool { return t.Kind == ToolKindCommand }

// IrreversibleClasses is the frozen R8 class set (SPEC-06 §3). Declaring a
// class outside it is a contract violation, so C14 refuses the manifest
// rather than letting an unreadable class sit in a declaration.
var IrreversibleClasses = []string{"delete", "overwrite", "power", "close-window", "send"}

// Risk levels (C2). A manifest-declared level is only the R1 lower bound: the
// central C19 inference may raise it, never the loader.
const (
	RiskL0 = "L0"
	RiskL1 = "L1"
	RiskL2 = "L2"
)

// Capabilities is the frozen C3 set (PLAN.md C3 / SPEC-07 §2, 11 items).
// A capability outside it is a contract violation, so the validator rejects
// the manifest rather than silently dropping the grant.
var Capabilities = []string{
	"fs.read", "fs.write", "net", "clipboard", "shell", "sysinfo",
	"screen", "input", "window", "notify", "memory",
}

var (
	manifestKeys = keySet("schema_version", "id", "host_api", "name", "version", "tools")

	// toolKeys is the declaration face (ticket 50). It is joined with
	// toolCommandKeys for a tool whose kind = "command"; see auditKeys.
	toolKeys = keySet("name", "description", "risk", "capabilities", "net_allowlist", "parameters",
		"kind", "irreversible", "taint_output")

	// toolCommandKeys is the D46 surface (C14, §16.5.6). These names are the
	// contract's own, so they stay listed in deferredKeys below: what changed
	// is that a kind = "command" tool now implements them.
	toolCommandKeys = []string{"exe", "exe_hash", "argv_template", "extract", "timeout_ms", "env_allowlist"}

	// deferredKeys maps a contract key that exists in C14 to the ticket that
	// owns its implementation. Present but unimplemented => refused load with
	// a named owner (never an "unknown key" lie, never a silent accept).
	deferredKeys = map[string]string{
		"entry":         "ticket 51 (Tier-2 goja entry)",
		"exe":           "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
		"exe_hash":      "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
		"argv_template": "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
		"extract":       "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
		"timeout_ms":    "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
		"env_allowlist": "ticket 52 (D46 command tools): legal only on a tool with kind = \"command\"",
	}

	// opaqueKeys are tables whose inner keys belong to another grammar (JSON
	// Schema for parameters), so the manifest key auditor does not descend.
	opaqueKeys = keySet("parameters")

	idPattern     = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	toolNamePair  = regexp.MustCompile(`^[a-z0-9_-]+(\.[a-z0-9_-]+)+$`)
	domainPattern = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?(\.[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?)*$`)
)

func keySet(keys ...string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

// ParseManifest decodes and C14-validates raw manifest bytes in one call:
// unknown/deferred key audit, typed decode, then per-field validation.
// It returns the full error list as a *ValidationErrors; err is nil only when
// the manifest is loadable.
func ParseManifest(raw []byte) (Manifest, error) {
	if errs := auditKeys(raw); len(errs) > 0 {
		return Manifest{}, errs
	}
	var m Manifest
	if err := toml.Unmarshal(raw, &m); err != nil {
		return Manifest{}, ValidationErrors{{
			Field:  "(document)",
			Reason: fmt.Sprintf("manifest is not valid TOML for the C14 shape: %v", err),
		}}
	}
	if errs := Validate(&m); len(errs) > 0 {
		return m, errs
	}
	return m, nil
}

// auditKeys decodes into a generic tree and reports every key that is neither
// in the C14 field table nor implemented by this slice.
func auditKeys(raw []byte) ValidationErrors {
	var generic map[string]any
	if err := toml.Unmarshal(raw, &generic); err != nil {
		return ValidationErrors{{
			Field:  "(document)",
			Reason: fmt.Sprintf("cannot parse TOML: %v", err),
		}}
	}
	var errs ValidationErrors
	errs = append(errs, auditTable(generic, "manifest", manifestKeys)...)

	toolsRaw, ok := generic["tools"]
	if !ok {
		return errs
	}
	list, ok := toolsRaw.([]any)
	if !ok {
		errs = append(errs, ValidationError{Field: "tools", Reason: "must be an array of tables ([[tools]])"})
		return errs
	}
	for i, t := range list {
		table, ok := t.(map[string]any)
		if !ok {
			errs = append(errs, ValidationError{
				Field:  fmt.Sprintf("tools[%d]", i),
				Reason: "each [[tools]] entry must be a table",
			})
			continue
		}
		prefix := fmt.Sprintf("tools[%d]", i)
		allowed := toolKeys
		if kind, _ := table["kind"].(string); kind == ToolKindCommand {
			allowed = commandToolKeys()
		}
		errs = append(errs, auditTable(table, prefix, allowed)...)
	}
	return errs
}

// commandToolKeys is the key table of a kind = "command" tool: the declaration
// face plus the D46 surface.
func commandToolKeys() map[string]bool {
	m := make(map[string]bool, len(toolKeys)+len(toolCommandKeys))
	for k := range toolKeys {
		m[k] = true
	}
	for _, k := range toolCommandKeys {
		m[k] = true
	}
	return m
}

// auditTable reports unknown and deferred keys of one table, without
// descending into opaque sub-tables (parameters).
func auditTable(table map[string]any, path string, allowed map[string]bool) ValidationErrors {
	var errs ValidationErrors
	keys := make([]string, 0, len(table))
	for k := range table {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if allowed[k] || opaqueKeys[k] {
			continue
		}
		if owner, ok := deferredKeys[k]; ok {
			errs = append(errs, ValidationError{
				Field:  path + "." + k,
				Reason: fmt.Sprintf("field is part of the C14 contract but is not implemented by the Tier-1 base loader; owner: %s", owner),
			})
			continue
		}
		errs = append(errs, ValidationError{
			Field:  path + "." + k,
			Reason: "unknown key: the C14 manifest schema is closed, unknown fields are rejected (never ignored)",
		})
	}
	return errs
}

// CapabilityDeclared reports whether caps contains c.
func CapabilityDeclared(caps []string, c string) bool {
	for _, x := range caps {
		if x == c {
			return true
		}
	}
	return false
}

func hasControl(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}
