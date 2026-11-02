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
//
// Deferred fields are NOT accepted by this slice; the key auditor names the
// owning ticket instead of pretending they never existed:
//   - entry                                    -> ticket 51 (Tier-2 goja)
//   - kind/exe/exe_hash/argv_template/extract/
//     timeout_ms/env_allowlist                 -> ticket 52 (D46 command tools)
//
// CurrentSchemaVersion is the only schema_version this loader understands
// (SPEC-07 §4 example: schema_version = 1).
const CurrentSchemaVersion = 1

// Manifest is a parsed, C14-validated plugin manifest.
type Manifest struct {
	SchemaVersion int       `toml:"schema_version"`
	ID            string    `toml:"id"`
	HostAPI       string    `toml:"host_api"`
	Name          string    `toml:"name"`
	Version       string    `toml:"version"`
	Tools         []ToolDef `toml:"tools"`
}

// ToolDef is one declarative Tier-1 tool: parameters (a JSON Schema) mapped
// onto declared host capabilities.
type ToolDef struct {
	Name         string         `toml:"name"`
	Description  string         `toml:"description"`
	Risk         string         `toml:"risk"`
	Capabilities []string       `toml:"capabilities"`
	NetAllowlist []string       `toml:"net_allowlist"`
	Parameters   map[string]any `toml:"parameters"`
}

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
	toolKeys     = keySet("name", "description", "risk", "capabilities", "net_allowlist", "parameters")

	// deferredKeys maps a contract key that exists in C14 to the ticket that
	// owns its implementation. Present but unimplemented => refused load with
	// a named owner (never an "unknown key" lie, never a silent accept).
	deferredKeys = map[string]string{
		"entry":         "ticket 51 (Tier-2 goja entry)",
		"kind":          "ticket 52 (D46 command tools)",
		"exe":           "ticket 52 (D46 command tools)",
		"exe_hash":      "ticket 52 (D46 command tools)",
		"argv_template": "ticket 52 (D46 command tools)",
		"extract":       "ticket 52 (D46 command tools)",
		"timeout_ms":    "ticket 52 (D46 command tools)",
		"env_allowlist": "ticket 52 (D46 command tools)",
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
// in the C14 Tier-1 field table nor implemented by this slice.
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
		for _, e := range auditTable(table, prefix, toolKeys) {
			errs = append(errs, e)
		}
	}
	return errs
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
