package secret

import (
	"errors"
	"strings"
	"testing"
)

func TestParseRefValid(t *testing.T) {
	cases := []struct {
		ref      string
		wantKind string
		wantVal  string
	}{
		{"dpapi:abc123", RefKindDPAPI, "abc123"},
		{"dpapi:llm.providers.openai.api_key", RefKindDPAPI, "llm.providers.openai.api_key"},
		{"dpapi:a-b_c.d", RefKindDPAPI, "a-b_c.d"},
		{"env:WISP_TEST_LLM_KEY", RefKindEnv, "WISP_TEST_LLM_KEY"},
		{"env:a", RefKindEnv, "a"},
	}
	for _, tc := range cases {
		kind, val, err := ParseRef(tc.ref)
		if err != nil {
			t.Fatalf("ParseRef(%q): %v", tc.ref, err)
		}
		if kind != tc.wantKind || val != tc.wantVal {
			t.Errorf("ParseRef(%q) = (%q, %q), want (%q, %q)", tc.ref, kind, val, tc.wantKind, tc.wantVal)
		}
	}
}

func TestParseRefInvalid(t *testing.T) {
	cases := []string{
		"",
		"dpapi:",                            // empty blob id
		"dpapi:.",                           // dot
		"dpapi:..",                          // dotdot
		"dpapi:a/b",                         // separator
		`dpapi:a\b`,                         // separator
		"dpapi:..\\evil",                    // traversal
		"dpapi:has space",                   // charset
		"dpapi:" + strings.Repeat("x", 129), // too long
		"env:",                              // empty name
		"env:HAS SPACE",                     // whitespace
		"env:A=B",                           // '=' reserved
		"env:A:B",                           // ':' reserved
		"DPAPI:abc",                         // prefixes are case-sensitive
		"sk-1234567890",                     // plaintext, no prefix
		"file:/etc/passwd",                  // unknown scheme
	}
	for _, ref := range cases {
		kind, val, err := ParseRef(ref)
		if err == nil {
			t.Errorf("ParseRef(%q) = (%q, %q), want error", ref, kind, val)
			continue
		}
		if !errors.Is(err, ErrInvalidRef) {
			t.Errorf("ParseRef(%q) err = %v, want ErrInvalidRef", ref, err)
		}
	}
}

func TestNewRefUniqueAndParseable(t *testing.T) {
	seen := make(map[string]bool, 100)
	for i := 0; i < 100; i++ {
		ref := NewRef()
		if !strings.HasPrefix(ref, RefPrefixDPAPI) {
			t.Fatalf("NewRef() = %q, want %q prefix", ref, RefPrefixDPAPI)
		}
		if seen[ref] {
			t.Fatalf("NewRef() repeated %q", ref)
		}
		seen[ref] = true
		if _, _, err := ParseRef(ref); err != nil {
			t.Fatalf("ParseRef(NewRef()): %v", err)
		}
	}
}

func TestValidBlobIDCharset(t *testing.T) {
	for _, ok := range []string{"a", "A9", "a.b", "a-b", "a_b", strings.Repeat("x", 128)} {
		if !ValidBlobID(ok) {
			t.Errorf("ValidBlobID(%q) = false, want true", ok)
		}
	}
	for _, bad := range []string{"", ".", "..", "a/b", `a\b`, "a b", "a:b", "a\nb", strings.Repeat("x", 129), "ünicode"} {
		if ValidBlobID(bad) {
			t.Errorf("ValidBlobID(%q) = true, want false", bad)
		}
	}
}
