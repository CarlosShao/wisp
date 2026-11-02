package plugin

import (
	"errors"
	"strings"
	"testing"
)

// AC1 — C14 manifest parsing + per-field schema validation (ticket 50).
// Contract authority: PLAN.md C14 / SPEC-07 §4 (fields), C1 (tool name shape),
// C2 (risk levels), C3 (capability set).

const validManifest = `
schema_version = 1
id = "com.example.lark"
host_api = "^1.0"
name = "Lark CLI"
version = "0.1.0"

[[tools]]
name = "lark.calendar_agenda"
description = "查询飞书日程"
risk = "L0"
capabilities = ["net"]
net_allowlist = ["open.feishu.cn"]
parameters = { type = "object", properties = { to = { type = "string" } }, required = ["to"] }
`

func TestParseManifest_GoldenPositive(t *testing.T) {
	m, err := ParseManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("valid manifest rejected: %v", err)
	}
	if m.SchemaVersion != CurrentSchemaVersion || m.ID != "com.example.lark" ||
		m.HostAPI != "^1.0" || m.Name != "Lark CLI" || m.Version != "0.1.0" {
		t.Fatalf("top-level fields decoded wrong: %+v", m)
	}
	if len(m.Tools) != 1 {
		t.Fatalf("want 1 tool, got %d", len(m.Tools))
	}
	tool := m.Tools[0]
	if tool.Name != "lark.calendar_agenda" || tool.Risk != RiskL0 {
		t.Fatalf("tool fields decoded wrong: %+v", tool)
	}
	if !CapabilityDeclared(tool.Capabilities, "net") {
		t.Fatalf("capabilities lost: %+v", tool.Capabilities)
	}
	if len(tool.NetAllowlist) != 1 || tool.NetAllowlist[0] != "open.feishu.cn" {
		t.Fatalf("net_allowlist lost: %+v", tool.NetAllowlist)
	}
	if tool.Parameters["type"] != "object" {
		t.Fatalf("parameters (JSON Schema) not carried through: %+v", tool.Parameters)
	}
}

// TestFieldMatrix is the "per-field positive/negative matrix" AC1 asks for:
// one mutation per field, each with its own named case.
func TestFieldMatrix(t *testing.T) {
	cases := []struct {
		name    string
		snippet string
		wantErr bool
		contain string
	}{
		// schema_version
		{"schema_version/ok", validManifest, false, ""},
		{"schema_version/too_new", strings.Replace(validManifest, "schema_version = 1", "schema_version = 2", 1), true, "schema_version"},
		{"schema_version/missing", "\n" + strings.Replace(validManifest, "schema_version = 1\n", "", 1), true, "schema_version"},

		// id
		{"id/ok", validManifest, false, ""},
		{"id/empty", strings.Replace(validManifest, `id = "com.example.lark"`, `id = ""`, 1), true, "id"},
		{"id/path_separator", strings.Replace(validManifest, `id = "com.example.lark"`, `id = "com.example/lark"`, 1), true, "id"},
		{"id/backslash", strings.Replace(validManifest, `id = "com.example.lark"`, `id = 'com\\example\\lark'`, 1), true, "id"},
		{"id/dotdot", strings.Replace(validManifest, `id = "com.example.lark"`, `id = ".."`, 1), true, "id"},

		// host_api
		{"host_api/ok", validManifest, false, ""},
		{"host_api/empty", strings.Replace(validManifest, `host_api = "^1.0"`, `host_api = ""`, 1), true, "host_api"},

		// name / version
		{"name/empty", strings.Replace(validManifest, `name = "Lark CLI"`, `name = ""`, 1), true, "name"},
		{"version/empty", strings.Replace(validManifest, `version = "0.1.0"`, `version = ""`, 1), true, "version"},

		// risk (C2)
		{"risk/L0", validManifest, false, ""},
		{"risk/L1", strings.Replace(validManifest, `risk = "L0"`, `risk = "L1"`, 1), false, ""},
		{"risk/L2", strings.Replace(validManifest, `risk = "L0"`, `risk = "L2"`, 1), false, ""},
		{"risk/self_downgrade_word", strings.Replace(validManifest, `risk = "L0"`, `risk = "read"`, 1), true, "risk"},
		{"risk/lowercase_l0", strings.Replace(validManifest, `risk = "L0"`, `risk = "l0"`, 1), true, "risk"},
		{"risk/empty", strings.Replace(validManifest, `risk = "L0"`, `risk = ""`, 1), true, "risk"},

		// capabilities (C3, 11 items)
		{"capabilities/fs.read", strings.Replace(validManifest, `capabilities = ["net"]`, `capabilities = ["fs.read"]`, 1), false, ""},
		{"capabilities/memory", strings.Replace(validManifest, `capabilities = ["net"]`, `capabilities = ["memory"]`, 1), false, ""},
		{"capabilities/not_in_C3", strings.Replace(validManifest, `capabilities = ["net"]`, `capabilities = ["filesystem"]`, 1), true, "capabilities"},
		{"capabilities/duplicate", strings.Replace(validManifest, `capabilities = ["net"]`, `capabilities = ["net", "net"]`, 1), true, "capabilities"},
		{"capabilities/empty_string", strings.Replace(validManifest, `capabilities = ["net"]`, `capabilities = [""]`, 1), true, "capabilities"},

		// net_allowlist (SPEC-07 §4: net 能力必须带目标域名白名单)
		{"net_allowlist/ok", validManifest, false, ""},
		{"net_allowlist/net_without_list", strings.Replace(validManifest, `net_allowlist = ["open.feishu.cn"]`, `net_allowlist = []`, 1), true, "net_allowlist"},
		{"net_allowlist/wildcard", strings.Replace(validManifest, `net_allowlist = ["open.feishu.cn"]`, `net_allowlist = ["*.feishu.cn"]`, 1), true, "net_allowlist"},
		{"net_allowlist/with_scheme_and_path", strings.Replace(validManifest, `net_allowlist = ["open.feishu.cn"]`, `net_allowlist = ["https://open.feishu.cn/api"]`, 1), true, "net_allowlist"},
		{"net_allowlist/with_port", strings.Replace(validManifest, `net_allowlist = ["open.feishu.cn"]`, `net_allowlist = ["open.feishu.cn:443"]`, 1), true, "net_allowlist"},
		{"net_allowlist/empty_entry", strings.Replace(validManifest, `net_allowlist = ["open.feishu.cn"]`, `net_allowlist = [""]`, 1), true, "net_allowlist"},

		// tool name (C1: namespace.action)
		{"name/no_namespace", strings.Replace(validManifest, `name = "lark.calendar_agenda"`, `name = "calendar_agenda"`, 1), true, "tools[0].name"},
		{"name/uppercase", strings.Replace(validManifest, `name = "lark.calendar_agenda"`, `name = "Lark.Calendar_Agenda"`, 1), true, "tools[0].name"},
		{"name/path_traversal", strings.Replace(validManifest, `name = "lark.calendar_agenda"`, `name = "lark/../../etc"`, 1), true, "tools[0].name"},
		{"description/empty", strings.Replace(validManifest, `description = "查询飞书日程"`, `description = ""`, 1), true, "tools[0].description"},

		// parameters (JSON Schema)
		{"parameters/missing", strings.Replace(validManifest, "parameters = { type = \"object\", properties = { to = { type = \"string\" } }, required = [\"to\"] }\n", "", 1), true, "tools[0].parameters"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseManifest([]byte(tc.snippet))
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("expected accept, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected rejection containing %q, got accept", tc.contain)
			}
			if !strings.Contains(err.Error(), tc.contain) {
				t.Fatalf("expected error containing %q, got: %v", tc.contain, err)
			}
		})
	}
}

// The C14 schema is closed: an unknown key is a refusal, never an ignore.
func TestUnknownKeysRejected(t *testing.T) {
	t.Run("top_level", func(t *testing.T) {
		src := strings.Replace(validManifest, "schema_version = 1\n", "schema_version = 1\nauto_approve = true\n", 1)
		_, err := ParseManifest([]byte(src))
		if err == nil {
			t.Fatal("unknown top-level key accepted")
		}
		if !strings.Contains(err.Error(), "unknown key") || !strings.Contains(err.Error(), "manifest.auto_approve") {
			t.Fatalf("want unknown-key refusal naming the key, got: %v", err)
		}
	})

	t.Run("tool_level", func(t *testing.T) {
		src := strings.Replace(validManifest, `risk = "L0"`, `risk = "L0"
silent_shell = true`, 1)
		_, err := ParseManifest([]byte(src))
		if err == nil {
			t.Fatal("unknown tool-level key accepted")
		}
		if !strings.Contains(err.Error(), "unknown key") || !strings.Contains(err.Error(), "tools[0].silent_shell") {
			t.Fatalf("want unknown-key refusal naming the key, got: %v", err)
		}
	})

	t.Run("near_miss_spelling", func(t *testing.T) {
		src := strings.Replace(validManifest, `net_allowlist =`, `net_allow_list =`, 1)
		if _, err := ParseManifest([]byte(src)); err == nil {
			t.Fatal("misspelled contract key accepted (silently dropped the whitelist)")
		}
	})
}

// Deferred C14 keys must name their owning ticket, not masquerade as unknown.
func TestDeferredKeysNameOwner(t *testing.T) {
	t.Run("tier2_entry", func(t *testing.T) {
		src := strings.Replace(validManifest, "schema_version = 1\n", "schema_version = 1\nentry = \"main.js\"\n", 1)
		_, err := ParseManifest([]byte(src))
		if err == nil || !strings.Contains(err.Error(), "ticket 51") {
			t.Fatalf("want refusal naming ticket 51, got: %v", err)
		}
	})
	t.Run("command_exe", func(t *testing.T) {
		src := strings.Replace(validManifest, `risk = "L0"`, `risk = "L0"
exe = "C:/bin/lark-cli.exe"`, 1)
		_, err := ParseManifest([]byte(src))
		if err == nil || !strings.Contains(err.Error(), "ticket 52") {
			t.Fatalf("want refusal naming ticket 52, got: %v", err)
		}
		if strings.Contains(err.Error(), "unknown key") {
			t.Fatalf("deferred key reported as unknown (wrong owner signal): %v", err)
		}
	})
}

// parameters is a JSON Schema: its inner keys belong to another grammar and
// must not be audited as manifest keys.
func TestParametersKeysAreOpaque(t *testing.T) {
	src := strings.Replace(validManifest,
		`parameters = { type = "object", properties = { to = { type = "string" } }, required = ["to"] }`,
		`parameters = { "$schema" = "https://json-schema.org/draft/2020-12/schema", type = "object", properties = { to = { type = "string", maxLength = 200, enum = ["a", "b"] } }, required = ["to"], additionalProperties = false }`, 1)
	if _, err := ParseManifest([]byte(src)); err != nil {
		t.Fatalf("JSON Schema inner keys audited as manifest keys: %v", err)
	}
}

func TestParseManifest_MalformedDocument(t *testing.T) {
	for name, src := range map[string]string{
		"not_toml":        "this is not toml at all = = =",
		"tools_not_array": "schema_version = 1\nid = \"a.b\"\nhost_api = \"^1.0\"\nname = \"n\"\nversion = \"1\"\ntools = \"lark.x\"\n",
		"tool_not_table":  "schema_version = 1\nid = \"a.b\"\nhost_api = \"^1.0\"\nname = \"n\"\nversion = \"1\"\ntools = [\"lark.x\"]\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseManifest([]byte(src)); err == nil {
				t.Fatal("malformed document accepted")
			}
		})
	}
}

// All findings surface in one pass so a plugin author fixes the file once.
func TestValidationReportsEveryError(t *testing.T) {
	src := strings.Replace(validManifest, `risk = "L0"`, `risk = "nope"`, 1)
	src = strings.Replace(src, `capabilities = ["net"]`, `capabilities = ["bogus"]`, 1)
	src = strings.Replace(src, `id = "com.example.lark"`, `id = ""`, 1)
	_, err := ParseManifest([]byte(src))
	if err == nil {
		t.Fatal("manifest with three defects accepted")
	}
	var errs ValidationErrors
	if !errors.As(err, &errs) {
		t.Fatalf("want ValidationErrors, got %T: %v", err, err)
	}
	joined := errs.Error()
	for _, want := range []string{"id", "risk", "bogus"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("report missing %q:\n%s", want, joined)
		}
	}
	if len(errs) < 3 {
		t.Fatalf("want >=3 findings, got %d: %s", len(errs), joined)
	}
}
