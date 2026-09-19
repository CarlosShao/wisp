package buildinfo

import "testing"

func TestParseEnv(t *testing.T) {
	for _, want := range []Env{EnvProd, EnvDev, EnvTest} {
		got, err := ParseEnv(string(want))
		if err != nil || got != want {
			t.Errorf("ParseEnv(%q) = %q, %v; want %q, nil", want, got, err, want)
		}
	}
	for _, bad := range []string{"", "Prod", "production", "dev ", " prod", "staging"} {
		if _, err := ParseEnv(bad); err == nil {
			t.Errorf("ParseEnv(%q) = nil error, want error", bad)
		}
	}
}

// TestResolveEnv covers the SPEC-03 §5.1 order: WISP_ENV variable wins, then
// the build-time default (DefaultEnv is ldflags-injected, so it is only
// asserted relative to itself here).
func TestResolveEnv(t *testing.T) {
	t.Setenv("WISP_ENV", "test")
	if got, err := ResolveEnv(); err != nil || got != EnvTest {
		t.Fatalf("ResolveEnv with WISP_ENV=test = %q, %v", got, err)
	}

	t.Setenv("WISP_ENV", "nonsense")
	if _, err := ResolveEnv(); err == nil {
		t.Fatal("ResolveEnv with invalid WISP_ENV must fail")
	}

	t.Setenv("WISP_ENV", "")
	want, err := ParseEnv(DefaultEnv)
	if err != nil {
		t.Fatalf("build default %q unparseable: %v", DefaultEnv, err)
	}
	if got, err := ResolveEnv(); err != nil || got != want {
		t.Fatalf("ResolveEnv fallback = %q, %v; want %q", got, err, want)
	}
}
