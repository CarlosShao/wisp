package cer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The manifest is the contract between this calculator and the leg that brings
// the audio (AISHELL-1, cer-a1). Every way it can be absent or wrong gets its
// own named error here, because the acting caller is a CI gate: "the set has not
// landed" and "the set is broken" must not be the same reading.

const testOrigin = "AISHELL-1:test:S0002/BAC009S0002W0122"

func hashOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func writeFile(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", name, err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func manifestRow(file string, audio []byte, speaker string, cond Condition) string {
	return strings.Join([]string{file, hashOf(audio), speaker, string(cond), testOrigin}, "\t")
}

func TestParseManifestHappyPath(t *testing.T) {
	audio := []byte("not real audio, just bytes")
	body := "# provenance: AISHELL-1 (OpenSLR SLR33), Apache-2.0\n" +
		manifestHeader + "\n" +
		manifestRow("near_clean/0001.wav", audio, "S0002", NearClean) + "\n" +
		manifestRow("noisy_far/0002.wav", audio, "S0999", NoisyFar) + "\n"
	srcs, err := ParseManifest(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if len(srcs) != 2 {
		t.Fatalf("parsed %d rows, want 2", len(srcs))
	}
	want := Source{
		File:      "near_clean/0001.wav",
		SHA256:    hashOf(audio),
		Speaker:   "S0002",
		Condition: NearClean,
		Origin:    testOrigin,
	}
	if srcs[0] != want {
		t.Fatalf("row 0 = %+v, want %+v", srcs[0], want)
	}
	if srcs[1].Condition != NoisyFar || srcs[1].Speaker != "S0999" {
		t.Fatalf("row order is not manifest order: %+v", srcs[1])
	}
}

func TestParseManifestAcceptsCRLFAndBlankLines(t *testing.T) {
	audio := []byte("bytes")
	body := manifestHeader + "\r\n\r\n" +
		manifestRow("near_clean/0001.wav", audio, "S0002", NearClean) + "\r\n"
	srcs, err := ParseManifest(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseManifest with CRLF: %v", err)
	}
	if len(srcs) != 1 || srcs[0].File != "near_clean/0001.wav" {
		t.Fatalf("CRLF manifest parsed as %+v", srcs)
	}
}

func TestManifestFailuresAreNamed(t *testing.T) {
	audio := []byte("bytes")
	good := manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)
	for _, tc := range []struct {
		name    string
		body    string
		wantErr error
		wantMsg string
	}{
		{"empty file", "", ErrManifestEmpty, ""},
		{"comments only", "# nothing here\n", ErrManifestEmpty, ""},
		{"header only", manifestHeader + "\n", ErrManifestEmpty, ""},
		// %q spells the tabs of the wanted header as \t, which is what makes the
		// message safe to paste into a terminal.
		{"wrong header", "file\tspeaker\tsha256\n" + good + "\n", ErrManifestHeader, `file\tsha256\tspeaker\tcondition\torigin`},
		{"four columns", manifestHeader + "\nnear_clean/0001.wav\tabc\tS0002\tnear_clean\n", ErrManifestRow, "4 columns"},
		{"uppercase digest", manifestHeader + "\n" + strings.Replace(good, hashOf(audio), strings.ToUpper(hashOf(audio)), 1) + "\n", ErrManifestRow, "lowercase hex"},
		{"short digest", manifestHeader + "\n" + strings.Replace(good, hashOf(audio), hashOf(audio)[:63], 1) + "\n", ErrManifestRow, "lowercase hex"},
		{"wrong directory", manifestHeader + "\nnear_clean/0001.wav\t" + hashOf(audio) + "\tS0002\tnoisy_far\t" + testOrigin + "\n", ErrManifestRow, "not under"},
		{"backslash", manifestHeader + "\nnear_clean\\0001.wav\t" + hashOf(audio) + "\tS0002\tnear_clean\t" + testOrigin + "\n", ErrManifestRow, "backslash"},
		{"dot-dot", manifestHeader + "\nnear_clean/../0001.wav\t" + hashOf(audio) + "\tS0002\tnear_clean\t" + testOrigin + "\n", ErrManifestRow, "clean form"},
		{"absolute", manifestHeader + "\n/near_clean/0001.wav\t" + hashOf(audio) + "\tS0002\tnear_clean\t" + testOrigin + "\n", ErrManifestRow, "absolute"},
		{"not wav", manifestHeader + "\nnear_clean/0001.txt\t" + hashOf(audio) + "\tS0002\tnear_clean\t" + testOrigin + "\n", ErrManifestRow, ".wav"},
		{"unknown condition", manifestHeader + "\nnear_clean/0001.wav\t" + hashOf(audio) + "\tS0002\tdev\t" + testOrigin + "\n", ErrManifestRow, "neither"},
		{"no speaker", manifestHeader + "\nnear_clean/0001.wav\t" + hashOf(audio) + "\t\tnear_clean\t" + testOrigin + "\n", ErrManifestRow, "speaker"},
		{"no origin", manifestHeader + "\nnear_clean/0001.wav\t" + hashOf(audio) + "\tS0002\tnear_clean\t\n", ErrManifestRow, "origin"},
		{"duplicate", manifestHeader + "\n" + good + "\n" + good + "\n", ErrManifestDuplicate, "near_clean/0001.wav"},
	} {
		_, err := ParseManifest(strings.NewReader(tc.body))
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("%s: got %v, want %v", tc.name, err, tc.wantErr)
		}
		if tc.wantMsg != "" && !strings.Contains(err.Error(), tc.wantMsg) {
			t.Fatalf("%s: message %q does not name %q", tc.name, err, tc.wantMsg)
		}
	}
}

// TestMissingManifestIsNamedAndNotAnEmptySet is reading 3 of the slice: absence
// must be a named error that a caller can branch on, never a nil error with zero
// rows. It also wraps fs.ErrNotExist so "not landed yet" is distinguishable from
// "broken" without parsing the message.
func TestMissingManifestIsNamedAndNotAnEmptySet(t *testing.T) {
	root := t.TempDir()
	srcs, err := ReadManifest(root)
	if !errors.Is(err, ErrManifestMissing) {
		t.Fatalf("ReadManifest on an empty directory: got %v, want ErrManifestMissing", err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing manifest must also wrap fs.ErrNotExist, got %v", err)
	}
	if srcs != nil {
		t.Fatalf("missing manifest returned %d rows, want nil", len(srcs))
	}
	if !strings.Contains(err.Error(), ManifestName) {
		t.Fatalf("message %q does not name %s", err, ManifestName)
	}
	if _, err := LoadSet(root, NearClean); !errors.Is(err, ErrManifestMissing) {
		t.Fatalf("LoadSet on an empty directory: got %v, want ErrManifestMissing", err)
	}
}

func TestLoadSetHappyPath(t *testing.T) {
	root := t.TempDir()
	audio := []byte("RIFF....WAVE")
	writeFile(t, root, "near_clean/0001.wav", audio)
	// Editors on this platform add a BOM and a trailing newline; neither is a
	// reference character.
	writeFile(t, root, "near_clean/0001.txt", []byte("\ufeff今天天气。\r\n"))
	writeFile(t, root, "sources.tsv", []byte(manifestHeader+"\n"+
		manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)+"\n"))

	set, err := LoadSet(root, NearClean)
	if err != nil {
		t.Fatalf("LoadSet: %v", err)
	}
	if set.Condition != NearClean {
		t.Fatalf("condition %q, want near_clean", set.Condition)
	}
	if len(set.Utterances) != 1 {
		t.Fatalf("loaded %d utterances, want 1", len(set.Utterances))
	}
	got := set.Utterances[0]
	if got.ID != "near_clean/0001.wav" {
		t.Fatalf("utterance id %q, want the manifest path", got.ID)
	}
	if got.Ref != "今天天气。" {
		t.Fatalf("reference %q, want the trimmed text", got.Ref)
	}
	// The reference still carries its full stop; normalisation is the scorer's
	// job, so the loaded set keeps the corpus bytes.
	if res, err := Score(got.Ref, "今天天气"); err != nil || res.CER != 0 {
		t.Fatalf("scoring the loaded reference: %v / %v", res, err)
	}
}

func TestLoadSetFailuresAreNamed(t *testing.T) {
	audio := []byte("RIFF....WAVE")
	for _, tc := range []struct {
		name     string
		cond     Condition
		manifest []string // rows, without the header
		wav      bool
		txt      *string
		wantErr  error
	}{
		{
			name:     "missing audio",
			cond:     NearClean,
			manifest: []string{manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)},
			wav:      false,
			txt:      strp("今天天气"),
			wantErr:  ErrAudioMissing,
		},
		{
			name:     "missing reference",
			cond:     NearClean,
			manifest: []string{manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)},
			wav:      true,
			txt:      nil,
			wantErr:  ErrReferenceMissing,
		},
		{
			name:     "empty reference",
			cond:     NearClean,
			manifest: []string{manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)},
			wav:      true,
			txt:      strp(" \r\n"),
			wantErr:  ErrEmptyReference,
		},
		{
			name:     "no rows for the condition",
			cond:     NearClean,
			manifest: []string{manifestRow("noisy_far/0001.wav", audio, "S0002", NoisyFar)},
			wav:      true,
			txt:      strp("今天天气"),
			wantErr:  ErrNoUtterances,
		},
	} {
		root := t.TempDir()
		if tc.wav {
			writeFile(t, root, "near_clean/0001.wav", audio)
			writeFile(t, root, "noisy_far/0001.wav", audio)
		}
		if tc.txt != nil {
			writeFile(t, root, "near_clean/0001.txt", []byte(*tc.txt))
			writeFile(t, root, "noisy_far/0001.txt", []byte(*tc.txt))
		}
		body := manifestHeader + "\n" + strings.Join(tc.manifest, "\n") + "\n"
		writeFile(t, root, "sources.tsv", []byte(body))

		set, err := LoadSet(root, tc.cond)
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("%s: got %v, want %v", tc.name, err, tc.wantErr)
		}
		if len(set.Utterances) != 0 {
			t.Fatalf("%s: failed load returned %d utterances", tc.name, len(set.Utterances))
		}
	}
	// A condition with no gate is refused before the manifest is even read.
	if _, err := LoadSet(t.TempDir(), "dev"); !errors.Is(err, ErrUnknownCondition) {
		t.Fatalf("unknown condition: got %v, want ErrUnknownCondition", err)
	}
}

func strp(s string) *string { return &s }

func TestVerifyManifest(t *testing.T) {
	root := t.TempDir()
	audio := []byte("RIFF....WAVE")
	wav := writeFile(t, root, "near_clean/0001.wav", audio)
	writeFile(t, root, "near_clean/0001.txt", []byte("今天天气"))
	writeFile(t, root, "sources.tsv", []byte(manifestHeader+"\n"+
		manifestRow("near_clean/0001.wav", audio, "S0002", NearClean)+"\n"))

	if err := VerifyManifest(root); err != nil {
		t.Fatalf("VerifyManifest on intact set: %v", err)
	}
	// One byte of drift is a mismatch, and the error names the file.
	if err := os.WriteFile(wav, []byte("RIFF....WAVe"), 0o600); err != nil {
		t.Fatalf("tamper: %v", err)
	}
	err := VerifyManifest(root)
	if !errors.Is(err, ErrDigestMismatch) {
		t.Fatalf("tampered set: got %v, want ErrDigestMismatch", err)
	}
	if !strings.Contains(err.Error(), "near_clean/0001.wav") {
		t.Fatalf("message %q does not name the file", err)
	}
}
