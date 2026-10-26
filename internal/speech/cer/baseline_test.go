package cer

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// This file is where the calculator meets the committed corpus. Two separate
// jobs, deliberately in two tests:
//
//   - TestBaselineCorpusIsCommittedAndTheGateCanRun is the real thing: it runs
//     the gate over testdata/asr-baseline/ when the corpus is there, and SKIPS
//     BY NAME when it is not (today: not there, it is a separate leg).
//   - TestAbsentCorpusIsANamedErrorNotAnEmptyPass runs unconditionally and is
//     the reason this package cannot become permanently green: it asserts that
//     absence is a named error, and that nothing without a denominator can
//     produce Pass=true.

// repoRoot walks up from this package's working directory to the module root, so
// the corpus is found whether the tests run in the package directory or from a
// staging directory that carries the sherpa DLLs.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("module root not found above %s", dir)
	return ""
}

// baselineRoot is where SPEC-01:86-87 puts the two committed sets, repeated
// verbatim at SPEC-04:125.
func baselineRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "testdata", "asr-baseline")
}

// checkCommittedSet applies every requirement the real corpus must satisfy: both
// conditions load, digests match, the two conditions share no speaker
// (cer-a1 section 4: taken from different speakers so noisy_far cannot leak the
// clean sentences), a perfect transcript passes and an empty one fails.
//
// It is exercised today against a synthetic set (see the test below), so the day
// the corpus lands this is not first-run code.
func checkCommittedSet(t *testing.T, root string) {
	t.Helper()
	clean, err := LoadSet(root, NearClean)
	if err != nil {
		t.Fatalf("LoadSet near_clean: %v", err)
	}
	far, err := LoadSet(root, NoisyFar)
	if err != nil {
		t.Fatalf("LoadSet noisy_far: %v", err)
	}
	if len(clean.Utterances) == 0 || len(far.Utterances) == 0 {
		t.Fatalf("loaded %d near_clean and %d noisy_far utterances, want both non-empty",
			len(clean.Utterances), len(far.Utterances))
	}
	if err := VerifyManifest(root); err != nil {
		t.Fatalf("VerifyManifest: %v", err)
	}

	srcs, err := ReadManifest(root)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	bySpeaker := map[string]Condition{}
	for _, s := range srcs {
		if prev, ok := bySpeaker[s.Speaker]; ok && prev != s.Condition {
			t.Fatalf("speaker %s appears in both %q and %q; the sets must not share speakers",
				s.Speaker, prev, s.Condition)
		}
		bySpeaker[s.Speaker] = s.Condition
	}

	for _, set := range []Set{clean, far} {
		perfect := make(map[string]string, len(set.Utterances))
		empty := make(map[string]string, len(set.Utterances))
		for _, u := range set.Utterances {
			perfect[u.ID] = u.Ref
			empty[u.ID] = ""
		}
		rep, err := set.Score(perfect)
		if err != nil {
			t.Fatalf("%s: scoring a perfect transcript: %v", set.Condition, err)
		}
		v, err := rep.Verdict()
		if err != nil {
			t.Fatalf("%s: Verdict: %v", set.Condition, err)
		}
		if !v.Pass || v.CER != 0 {
			t.Fatalf("%s: a perfect transcript must pass with CER 0, got %s", set.Condition, v)
		}
		// The gate has teeth on real data too: nothing heard is 100%, and 100%
		// fails both gates.
		rep, err = set.Score(empty)
		if err != nil {
			t.Fatalf("%s: scoring an empty transcript: %v", set.Condition, err)
		}
		v, err = rep.Verdict()
		if err != nil {
			t.Fatalf("%s: Verdict: %v", set.Condition, err)
		}
		if v.Pass || v.CER != 1.0 {
			t.Fatalf("%s: an empty transcript must fail with CER 100%%, got %s", set.Condition, v)
		}
		t.Logf("%s: %s", set.Condition, v)
	}
}

// TestBaselineCorpusIsCommittedAndTheGateCanRun is the only test that touches
// testdata/asr-baseline/. Today the directory does not exist, so it SKIPS with a
// named reason instead of passing with zero utterances.
func TestBaselineCorpusIsCommittedAndTheGateCanRun(t *testing.T) {
	root := baselineRoot(t)
	manifest := filepath.Join(root, ManifestName)
	if _, err := os.Stat(manifest); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("stat %s: %v", manifest, err)
		}
		// Assert the named error BEFORE skipping: this line is what makes the
		// skip a report about missing data rather than a silent pass.
		if _, err := LoadSet(root, NearClean); !errors.Is(err, ErrManifestMissing) {
			t.Fatalf("absent baseline set: LoadSet returned %v, want ErrManifestMissing", err)
		}
		t.Skipf("baseline corpus absent (%s): the CER gate is NOT verified by this run. "+
			"The committed sets are a separate leg - AISHELL-1, 60 near_clean + 40 noisy_far, "+
			"see docs/evidence/s1/cer-a1-research.md. SKIP, not pass: a green over 0 utterances "+
			"would be a lie.", manifest)
	}
	checkCommittedSet(t, root)
}

// TestCommittedSetChecksRunTodayAgainstASyntheticSet exercises the same
// checkCommittedSet the corpus test uses, on a manifest and files this test
// builds in t.TempDir(). The bytes are not audio and are never decoded: this
// package hashes them and reads the sibling .txt, which is the whole contract.
func TestCommittedSetChecksRunTodayAgainstASyntheticSet(t *testing.T) {
	root := t.TempDir()
	near := []byte("not audio, near_clean bytes")
	far := []byte("not audio, noisy_far bytes")
	writeFile(t, root, "near_clean/0001.wav", near)
	writeFile(t, root, "near_clean/0001.txt", []byte("今天天气"))
	writeFile(t, root, "noisy_far/0001.wav", far)
	writeFile(t, root, "noisy_far/0001.txt", []byte("今天天气很好"))
	body := manifestHeader + "\n" +
		manifestRow("near_clean/0001.wav", near, "S0002", NearClean) + "\n" +
		manifestRow("noisy_far/0001.wav", far, "S0999", NoisyFar) + "\n"
	writeFile(t, root, "sources.tsv", []byte(body))

	checkCommittedSet(t, root)
}

// TestAbsentCorpusIsANamedErrorNotAnEmptyPass runs whether or not the corpus
// exists, and is the guard against "no data means green": load an absent tree,
// then demand that every route to a verdict without a denominator is refused.
func TestAbsentCorpusIsANamedErrorNotAnEmptyPass(t *testing.T) {
	// A tree that has never existed, not even testdata/.
	absent := filepath.Join(t.TempDir(), "testdata", "asr-baseline")
	set, err := LoadSet(absent, NearClean)
	if !errors.Is(err, ErrManifestMissing) {
		t.Fatalf("absent corpus: LoadSet returned %d utterances and %v, want ErrManifestMissing",
			len(set.Utterances), err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("absent corpus must stay distinguishable from a broken one, got %v", err)
	}

	// The repository's own tree, today absent. Once the corpus lands this branch
	// stops running and the test above takes over - it must not fail then.
	realRoot := baselineRoot(t)
	if _, err := os.Stat(filepath.Join(realRoot, ManifestName)); errors.Is(err, fs.ErrNotExist) {
		if set, err := LoadSet(realRoot, NearClean); !errors.Is(err, ErrManifestMissing) {
			t.Fatalf("the repository's absent baseline set scored as %d utterances with %v",
				len(set.Utterances), err)
		}
	}

	// No denominator, no verdict: every route is an error, and none of the
	// zero-value reports is marked as passing.
	for _, rep := range []Report{
		{},
		{Condition: NearClean},
		{Condition: NoisyFar, Aggregate: Aggregate{}},
	} {
		if v, err := rep.Verdict(); err == nil || v.Pass {
			t.Fatalf("report %+v produced verdict %s / %v, want a refusal", rep, v, err)
		}
	}
	if v, err := (Set{Condition: NoisyFar}).Score(map[string]string{}); !errors.Is(err, ErrNoUtterances) {
		t.Fatalf("empty set scored as %v / %v, want ErrNoUtterances", v, err)
	}
}
