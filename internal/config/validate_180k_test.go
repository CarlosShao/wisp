package config

// Ticket 180-census follow-up (the 19 orphan config keys leg, 423's second
// wave): two [audio]/[voice] keys whose file values the shipped code could
// never honour get the load-time guard instead of a silent no-op. The census
// measured audio.sample_rate as "12 prod hits, all the audio engine's own
// same-named fields - the config value never reaches the chain", and
// voice.aec.echo_ref as zero mentions anywhere. The chain itself is a 16k
// chain end to end and the AEC path implements exactly one echo-reference
// source, so the honest reader for both keys is the guard the half_duplex
// precedent already keeps: a value the code cannot act on is a load error
// naming the key, never a setting that does nothing.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// guard180Config is a minimal valid file with one section appended; the
// decoder merges over the compiled defaults, so only the planted key differs.
func guard180Write(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func Test180KSampleRateChainValueLoads(t *testing.T) {
	c, _, err := LoadFile(guard180Write(t, "[audio]\nsample_rate = 16000\n"), nil)
	if err != nil {
		t.Fatalf("sample_rate = 16000 (the chain value) must load: %v", err)
	}
	if c.Audio.SampleRate != 16000 {
		t.Fatalf("sample_rate carried %d, want 16000", c.Audio.SampleRate)
	}
}

func Test180KSampleRateOffChainValueIsALoadError(t *testing.T) {
	_, _, err := LoadFile(guard180Write(t, "[audio]\nsample_rate = 48000\n"), nil)
	if err == nil {
		t.Fatal("sample_rate = 48000 must be a load error: the shipped chain honours no other rate, so accepting it would be a key that lies")
	}
	if !strings.Contains(err.Error(), "audio.sample_rate") {
		t.Fatalf("the refusal must name the key, got: %v", err)
	}
}

func Test180KSampleRateGuardRulerStillNamesTheRealChainRate(t *testing.T) {
	// The guard's constant must equal the chain's own resample target. This
	// package must not import internal/audio (the dependency graph runs one
	// direction), so the number is written twice on purpose and pinned here
	// by reading internal/audio's source: a drift reddens this test, not the
	// user's boot.
	src, err := os.ReadFile(filepath.Join("..", "audio", "audio.go"))
	if err != nil {
		t.Fatalf("read internal/audio/audio.go: %v", err)
	}
	if !strings.Contains(string(src), "TargetRate = 16000") {
		t.Fatal("internal/audio TargetRate moved away from 16000; internal/config's audioSampleRateChain and the guard sentence must move with it (and the whole 16k chain contract with both)")
	}
}

func Test180KAECEchoRefVocabulary(t *testing.T) {
	// The default (and empty merge) loads; the one implemented source loads;
	// anything else is a load error naming the vocabulary.
	if _, _, err := LoadFile(guard180Write(t, "[voice.aec]\necho_ref = \"self-render\"\n"), nil); err != nil {
		t.Fatalf("echo_ref = self-render must load: %v", err)
	}
	_, _, err := LoadFile(guard180Write(t, "[voice.aec]\necho_ref = \"loopback\"\n"), nil)
	if err == nil {
		t.Fatal("echo_ref = loopback must be a load error: no second echo-reference source exists, so accepting it would be a silent fallback")
	}
	if !strings.Contains(err.Error(), "voice.aec.echo_ref") || !strings.Contains(err.Error(), "self-render") {
		t.Fatalf("the refusal must name the key and its vocabulary, got: %v", err)
	}
}
