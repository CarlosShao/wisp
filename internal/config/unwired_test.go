package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Ticket 83 (ruling A53②, ticket 80 option (C)): a locked-section key that the
// direction auditor treats as a security change but that nothing reads must
// fail the load loudly instead of being swallowed.

// unwiredCases pairs each guarded key with the fragment that trips it and the
// two things the error must name: the key path and who lands the capability.
var unwiredCases = []struct {
	path     string
	fragment string
	wantTalk string // plan/owner sentence marker
}{
	{"risk.shell_enabled", "[risk]\nshell_enabled = true\n", "shell.exec"},
	{"risk.allow_shell_string", "[risk]\nallow_shell_string = true\n", "shell.exec"},
	{"risk.shell_allowlist", "[risk]\nshell_allowlist = [\"git\"]\n", "shell.exec"},
	{"risk.blacklist_overrides", "[risk]\nblacklist_overrides = [\"D:\\\\secret\"]\n", "ticket 21"},
	{"net.allowlist", "[net]\nallowlist = [\"api.example.com\"]\n", "ticket 22"},
	{"net.block_private_ranges", "[net]\nblock_private_ranges = false\n", "D22"},
}

// TestUnwiredSecurityKeysFailLoudly loads each guarded key through the real
// LoadFile pipeline: the load must error, and the error must name the key and
// say who implements the capability (AC#2).
func TestUnwiredSecurityKeysFailLoudly(t *testing.T) {
	for _, tc := range unwiredCases {
		t.Run(tc.path, func(t *testing.T) {
			path := writeConfigFile(t, "schema_version = 2\n\n"+tc.fragment)
			if _, _, err := LoadFile(path, nil); err == nil {
				t.Fatalf("%s must fail the load, got no error", tc.path)
			} else {
				msg := err.Error()
				if !strings.Contains(msg, tc.path) {
					t.Errorf("error %q must name the key %q", msg, tc.path)
				}
				if !strings.Contains(msg, tc.wantTalk) {
					t.Errorf("error %q must point at %q", msg, tc.wantTalk)
				}
			}
		})
	}
}

// TestUnwiredGuardLeavesHonestConfigsAlone is the other half of AC#2: a user
// who never writes these keys must not be affected. Both "absent" and "present
// at the default value" have to load, because SaveFile writes every field out.
func TestUnwiredGuardLeavesHonestConfigsAlone(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		path := writeConfigFile(t, "schema_version = 2\n")
		if _, _, err := LoadFile(path, nil); err != nil {
			t.Fatalf("bare config must load: %v", err)
		}
	})
	t.Run("guarded keys at defaults", func(t *testing.T) {
		path := writeConfigFile(t, `schema_version = 2

[risk]
confirm_timeout_sec = 300
l1_window_sec = 2
shell_enabled = false
allow_shell_string = false
shell_allowlist = []
blacklist_overrides = []

[net]
allowlist = []
block_private_ranges = true
`)
		c, _, err := LoadFile(path, nil)
		if err != nil {
			t.Fatalf("defaults written out explicitly must load: %v", err)
		}
		if c.Risk.ShellEnabled || c.Risk.AllowShellString || len(c.Risk.ShellAllowlist) > 0 ||
			len(c.Risk.BlacklistOverrides) > 0 || len(c.Net.Allowlist) > 0 || !c.Net.BlockPrivateRanges {
			t.Fatalf("defaults did not survive: %+v %+v", c.Risk, c.Net)
		}
	})
	t.Run("SaveFile output reloads", func(t *testing.T) {
		dir := sealableTempDir124(t)
		path := filepath.Join(dir, "config.toml")
		if err := SaveFile(path, NewDefaults()); err != nil {
			t.Fatal(err)
		}
		if _, _, err := LoadFile(path, nil); err != nil {
			t.Fatalf("a config the app itself wrote must load: %v", err)
		}
	})
	t.Run("tightening a guarded key past its default is still rejected", func(t *testing.T) {
		// shell_enabled=false is the default; the guard is about the value
		// nothing consumes, so the false case must stay loadable (checked by
		// the subtest above) while true fails.
		path := writeConfigFile(t, "schema_version = 2\n\n[risk]\nshell_enabled = true\n")
		if _, _, err := LoadFile(path, nil); err == nil {
			t.Fatal("shell_enabled = true must fail")
		}
	})
}

// TestUnwiredGuardFiresBeforeAnyFileWrite pins AC#2's "no side effects" half.
// Migration is the only load step that writes to disk (backup + atomic
// rewrite), and it dry-checks the migrated bytes before touching the file. A
// v1 config carrying a lying key must therefore leave the file byte-identical
// and create no backup.
func TestUnwiredGuardFiresBeforeAnyFileWrite(t *testing.T) {
	const v1WithLyingKey = `schema_version = 1

[app]
theme = "dark"

[risk]
shell_allowlist = ["git"]

[llm]
default_provider = "openai"
timeout = 45000

[llm.providers.openai]
base_url = "https://api.openai.com/v1"
model = "gpt-4o-mini"
`
	path := writeConfigFile(t, v1WithLyingKey)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadFile(path, nil); err == nil {
		t.Fatal("a v1 file carrying risk.shell_allowlist must not load")
	} else if !strings.Contains(err.Error(), "risk.shell_allowlist") {
		t.Fatalf("error must name the key, got %q", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("the rejected file changed on disk:\nbefore=%q\nafter =%q", before, after)
	}
	if _, err := os.Stat(path + ".bak-1"); !os.IsNotExist(err) {
		t.Fatalf("no backup may be written before validation passes (stat err=%v)", err)
	}
	tmp, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".wisp-config-*.tmp"))
	if len(tmp) > 0 {
		t.Fatalf("no temp write may happen before validation passes, found %v", tmp)
	}
}

// TestUnwiredGuardIsScopedNotABigStick is the AC#3(ii) anchor: the guard must
// be a named list of keys, not a blanket "anything we do not consume is an
// error". Every zero-consumer key outside the guarded list has to keep loading
// (their owning subsystems are scheduled work, see ticket 83 AC#1's table), and
// an unrelated unknown key must still be reported as an *unknown key* by the
// decoder, not by this guard.
func TestUnwiredGuardIsScopedNotABigStick(t *testing.T) {
	path := writeConfigFile(t, `schema_version = 2

[app]
language = "en-US"
autostart = true

[ball]
size = 60
click_through = false

[session]
warm_timeout_sec = 120

[voice]
enabled = false
cloud_asr_chain = ["openai/gpt-4o-mini"]

[audio]
input_device = "mic-1"
sample_rate = 48000

[agent]
max_rounds = 10
token_budget = 1000

[risk]
confirm_timeout_sec = 60
l1_window_sec = 5

[fs]
allowed_dirs = ["D:\\work"]
reparse_point_exceptions = ["C:\\link"]
delete_enabled = true

[net]
proxy.mode = "manual"
proxy.url = "http://proxy:8080"

[llm]
text_chain = ["openai/gpt-4o-mini"]

[llm.providers.openai]
models."gpt-4o-mini".quota_daily_micro = 1000000
models."gpt-4o-mini".thinking_levels = ["off", "low"]
models."gpt-4o-mini".context_window = 128000

[privacy]
diagnostics_opt_in = true
retention_days = 7

[memory]
l1_max = 5
extract_model = "openai/gpt-4o-mini"

[panel]
width = 800
height = 600
scale = 1.5

[cost]
daily_budget = 1000000
alert_threshold = 0.5

[models]
dir = "models"
mirror = ["http://mirror.example"]
local_override = { "gpt-4o-mini" = "D:\\local\\m" }

[observe]
level = "warn"
slo_sample_interval_sec = 2

[plugins]
tier2_enabled = true

[plugins.clip]
enabled = true
capabilities = ["clipboard"]
net_allowlist = ["a.example"]
host_api = ["fs"]
`)
	if _, _, err := LoadFile(path, nil); err != nil {
		t.Fatalf("only the six named keys are guarded, this config must load: %v", err)
	}

	// A key that is genuinely unknown is still the decoder's job, and its
	// message shape is different: it says "unknown key", not "does nothing".
	bad := writeConfigFile(t, "schema_version = 2\n\n[risk]\nnot_a_real_key = 1\n")
	_, _, err := LoadFile(bad, nil)
	if err == nil {
		t.Fatal("unknown keys must still be rejected")
	}
	if !strings.Contains(err.Error(), "unknown key") {
		t.Fatalf("unknown key must be reported by the strict decoder, got %q", err)
	}
	if strings.Contains(err.Error(), "does nothing") {
		t.Fatalf("the unwired guard must not be the thing that rejects unknown keys, got %q", err)
	}
}

// TestUnwiredKeysStillRoundTripAtTheByteLevel keeps the schema honest: these
// keys are rejected at load, they are not deleted from the schema. Their
// marshal/decode symmetry still has to hold, because the day ticket 21/22 lands
// the consumer the row disappears and the GUI double-source guard (D36 rule 4)
// must already cover them. validate() is deliberately bypassed here.
func TestUnwiredKeysStillRoundTripAtTheByteLevel(t *testing.T) {
	c := NewDefaults()
	c.SchemaVersion = SchemaVersionCurrent
	c.Risk.ShellEnabled = true
	c.Risk.AllowShellString = true
	c.Risk.ShellAllowlist = []string{"git", "go"}
	c.Risk.BlacklistOverrides = []string{"B-1"}
	c.Net.Allowlist = []string{"api.example.com"}
	c.Net.BlockPrivateRanges = false

	if err := validate(c); err == nil {
		t.Fatal("sanity: the guarded values must be the ones validate rejects")
	} else if !strings.Contains(err.Error(), "risk.shell_enabled") {
		t.Fatalf("unexpected first rejection: %v", err)
	}

	data, err := MarshalCanonical(c)
	if err != nil {
		t.Fatal(err)
	}
	back := NewDefaults()
	if err := decodeStrict(data, back); err != nil {
		t.Fatalf("decode: %v", err)
	}
	normalizeConfig(back)
	// Presets are part of the load pipeline only (applyPresets); the schema
	// round-trip must be the identity without them.
	if got, want := back.LLM.Providers, c.LLM.Providers; !reflect.DeepEqual(got, want) {
		t.Fatalf("providers drifted: %v vs %v", got, want)
	}
	if !back.Risk.ShellEnabled || !back.Risk.AllowShellString ||
		!equalStrings(back.Risk.ShellAllowlist, []string{"git", "go"}) ||
		!equalStrings(back.Risk.BlacklistOverrides, []string{"B-1"}) ||
		!equalStrings(back.Net.Allowlist, []string{"api.example.com"}) ||
		back.Net.BlockPrivateRanges != false {
		t.Fatalf("guarded keys lost their values in the round trip: %+v %+v", back.Risk, back.Net)
	}
}

// TestEveryLockedSectionKeyIsAccountedFor makes AC#1's completeness
// executable: every key of a locked section must appear in lockedKeyDisposition
// with a stated reason, and the "unwired:" entries must match the guard table
// exactly. A new locked key added without a verdict - the way the four [risk]
// keys were added - fails here.
func TestEveryLockedSectionKeyIsAccountedFor(t *testing.T) {
	paths := []string{}
	var collect func(prefix string, v reflect.Value)
	collect = func(prefix string, v reflect.Value) {
		typ := v.Type()
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag := f.Tag.Get("toml")
			if tag == "-" || tag == "" {
				continue
			}
			path := prefix + "." + tag
			fv := v.Field(i)
			switch fv.Kind() {
			case reflect.Struct:
				paths = append(paths, path)
				collect(path, fv)
			case reflect.Map:
				elem := fv.Type().Elem()
				for elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Interface {
					elem = elem.Elem()
				}
				paths = append(paths, path+".<id>")
				if elem.Kind() == reflect.Struct {
					collect(path+".<id>", reflect.New(elem).Elem())
				}
			default:
				paths = append(paths, path)
			}
		}
	}
	collect("risk", reflect.ValueOf(RiskSection{}))
	collect("fs", reflect.ValueOf(FSSection{}))
	collect("net", reflect.ValueOf(NetSection{}))
	collect("plugins", reflect.ValueOf(PluginsSection{}))

	guarded := map[string]bool{}
	for _, k := range unwiredKeys {
		guarded[k.path] = false
	}
	for _, p := range paths {
		why, ok := lockedKeyDisposition[p]
		if !ok {
			t.Errorf("locked key %q has no disposition in lockedKeyDisposition (ticket 83 AC#1: every locked key needs a verdict)", p)
			continue
		}
		if strings.HasPrefix(why, "unwired:") {
			key := strings.TrimPrefix(why, "unwired:")
			if key != p {
				t.Errorf("locked key %q claims disposition %q, the path must match", p, why)
			}
			if _, seen := guarded[key]; !seen {
				t.Errorf("locked key %q points at guard %q which unwiredKeys does not contain", p, key)
				continue
			}
			guarded[key] = true
		}
	}
	for k, hit := range guarded {
		if !hit {
			t.Errorf("unwiredKeys guards %q but no locked-section key points at it", k)
		}
	}
}
