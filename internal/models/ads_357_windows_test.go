//go:build windows

package models

// Ticket 357 AC#3, the Windows half: the enumeration itself, measured on real
// NTFS with no mock.
//
// The reading these cases have to make impossible is form 丁's false negative. An
// alternate data stream is invisible to os.ReadDir and to ticket 121's
// "nothing unnamed" walk, and readable as model.onnx:attacker, so the tree a
// guard calls consistent can carry bytes nobody signed - 357-a1 measured exactly
// that against the real binary (D1: 12-byte model.onnx + 19-byte
// model.onnx:attacker => 与已验签清单一致, rc=0). So:
//
//	the forward case  : a planted stream must be ENUMERATED, by name
//	the reverse case  : a tree with no stream must come back empty, otherwise the
//	                    report would fire on every clean install
//	the many case     : 200 short streams must all come back. One 4096-byte call
//	                    answers ERROR_MORE_DATA and an empty list at this count
//	                    (357-r1's measurement), which is precisely the shape that
//	                    would turn "the file worth reporting" into "no streams
//	                    here", so the retry is what this case is for
//	the ceiling case  : the same call with a ceiling it cannot meet must return a
//	                    NAMED error, never an empty list - a bounded retry that
//	                    silently answers "nothing" is the same false negative one
//	                    layer down
//	the override case : the streams reported must come from the local_override
//	                    directory, not from <DataDir>/<id>. That is ticket 357's
//	                    first face, and it is the easiest thing to get wrong here

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// adsFixture installs the kws fixture into a fresh store and returns the store
// root, a Manager over it, and the installed directory.
func adsFixture(t *testing.T) (store string, m *Manager, dir string) {
	t.Helper()
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	mf := manifestWithURLs(t, srv.URL())
	store = filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	installed, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("install the fixture: %v", err)
	}
	return store, mgr, installed
}

// plantStream writes one named data stream onto path. On NTFS this is the same
// os.WriteFile an installer would use; the ":name" suffix is the OS's own
// spelling, and the reopen below is what proves the bytes really landed on that
// stream rather than in a file whose name happens to contain a colon.
func plantStream(t *testing.T, path, name, body string) {
	t.Helper()
	if err := os.WriteFile(path+":"+name, []byte(body), 0o644); err != nil {
		t.Fatalf("planting the stream %s:%s needs a filesystem with alternate data streams "+
			"(NTFS); this case is worthless without one: %v", filepath.Base(path), name, err)
	}
	got, err := os.ReadFile(path + ":" + name)
	if err != nil || string(got) != body {
		t.Fatalf("the stream %s:%s was not planted: read back %q, %v", filepath.Base(path), name, got, err)
	}
}

func Test357AlternateStreamIsEnumeratedByName(t *testing.T) {
	_, m, dir := adsFixture(t)
	plantStream(t, filepath.Join(dir, "model.onnx"), "attacker", "wisp-ads-attacker")

	gotDir, streams, err := m.UnnamedStreams("kws-fixture")
	if err != nil {
		t.Fatalf("ticket 357: enumerating the streams of a tree that carries one failed: %v", err)
	}
	if gotDir != dir {
		t.Errorf("ticket 357: the enumeration reports on %s, want the directory VerifyInstalled verified (%s) - a report about another tree is ticket 357's first face all over again", gotDir, dir)
	}
	want := "model.onnx:attacker:$DATA"
	if len(streams) != 1 || streams[0] != want {
		t.Fatalf("ticket 357: a planted alternate data stream was not enumerated by name; got %v, want exactly [%q]. "+
			"An ADS is invisible to os.ReadDir and to the tree walk ticket 121 added, so a miss here is the "+
			"silent pass form 丁 exists to end", streams, want)
	}
	t.Logf("enumerated %v under %s", streams, gotDir)
}

func Test357TreeWithoutStreamsEnumeratesEmpty(t *testing.T) {
	_, m, dir := adsFixture(t)

	gotDir, streams, err := m.UnnamedStreams("kws-fixture")
	if err != nil {
		t.Fatalf("ticket 357 reverse leg: a clean install failed to enumerate: %v", err)
	}
	if gotDir != dir {
		t.Errorf("ticket 357: enumerated %s, want %s", gotDir, dir)
	}
	if len(streams) != 0 {
		t.Fatalf("ticket 357 reverse leg: a directory whose files carry no second stream reported %v. "+
			"Form 丁 says a tree with no such stream must print exactly what it printed before, so a "+
			"non-empty answer here would add noise to every clean install", streams)
	}
	// The same tree must still verify: the enumeration is a read, not a refusal.
	if err := m.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("ticket 357: enumerating streams changed the verdict of a clean install: %v", err)
	}
}

// Test357ManyStreamsAllComeBack is the ERROR_MORE_DATA leg. 200 one-byte streams
// overflow the first 4096-byte buffer, and the call answers with an empty list and
// ERROR_MORE_DATA - so without the grow-and-retry this case reports zero streams
// on the one file that carries 200.
func Test357ManyStreamsAllComeBack(t *testing.T) {
	_, m, dir := adsFixture(t)
	target := filepath.Join(dir, "model.onnx")
	const count = 200
	for i := range count {
		plantStream(t, target, fmt.Sprintf("s%03d", i), "1")
	}

	gotDir, streams, err := m.UnnamedStreams("kws-fixture")
	if err != nil {
		t.Fatalf("ticket 357: enumerating %d streams failed instead of growing the buffer: %v", count, err)
	}
	if gotDir != dir {
		t.Errorf("ticket 357: enumerated %s, want %s", gotDir, dir)
	}
	if len(streams) != count {
		t.Fatalf("ticket 357 ERROR_MORE_DATA leg: %d streams planted, %d enumerated. A single 4096-byte call "+
			"answers ERROR_MORE_DATA and an EMPTY list at this count (357-r1 measured it), so a number below %d "+
			"here means the retry is missing and a flooded file reads as a clean one", count, len(streams), count)
	}
	for i := range count {
		want := fmt.Sprintf("model.onnx:s%03d:$DATA", i)
		if streams[i] != want {
			t.Fatalf("ticket 357: stream %d enumerated as %q, want %q (the list is sorted, so a stable order is "+
				"part of the answer)", i, streams[i], want)
		}
	}
	t.Logf("all %d streams enumerated from the same file", len(streams))
}

// Test357BoundedRetryNamesTheCeilingInsteadOfAnsweringEmpty pins the upper bound
// of that retry. The requirement is two-sided: the loop must stop, and stopping
// must be a NAMED failure rather than a quiet empty answer - the second half is
// what keeps a truncated read from being read as "no streams here".
func Test357BoundedRetryNamesTheCeilingInsteadOfAnsweringEmpty(t *testing.T) {
	_, _, dir := adsFixture(t)
	target := filepath.Join(dir, "model.onnx")
	// Enough streams that the answer cannot fit in 1024 bytes: each record is a
	// 24-byte header plus a UTF-16 name, so 64 of them are ~4 KiB.
	for i := range 64 {
		plantStream(t, target, fmt.Sprintf("s%03d", i), "1")
	}

	streams, err := listNamedStreamsBounded(target, 1024, 1024)
	if err == nil {
		t.Fatalf("ticket 357: a ceiling the answer cannot meet returned %d stream(s) and no error. Form 丁 "+
			"forbids reading an incomplete enumeration as an empty one, so a bounded retry that hits its bound "+
			"has to be an error the caller can print", len(streams))
	}
	if !strings.Contains(err.Error(), "did not fit in 1024 bytes") {
		t.Errorf("ticket 357: the ceiling refusal does not name the bound it hit, so nobody can tell a truncated "+
			"read from a failed call: %v", err)
	}
	if streams != nil {
		t.Errorf("ticket 357: the refusal still handed back %v; a partial answer next to an error invites a caller "+
			"to report the partial list as the whole one", streams)
	}
	// The same file with the production bounds must come back complete, so this
	// case cannot pass by making every enumeration fail.
	full, ferr := listNamedStreams(target)
	if ferr != nil {
		t.Fatalf("ticket 357: the same file fails under the production bounds: %v", ferr)
	}
	if len(full) != 64 {
		t.Fatalf("ticket 357: under the production bounds the same file enumerated %d of 64 streams", len(full))
	}
}

// Test357OverrideDirectoryIsTheOneEnumerated is the trap case. The manifest's
// files are installed in the local_override directory and only THERE, while
// <DataDir>/<id> holds a planted stream. If the enumeration follows the store root
// - which is what the CLI prints on that same line - it reports on a directory
// nobody verified; the answer has to come from the override.
func Test357OverrideDirectoryIsTheOneEnumerated(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	mf := manifestWithURLs(t, srv.URL())
	store := filepath.Join(t.TempDir(), "models")
	seed, err := NewManager(Options{DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := seed.Ensure(context.Background(), "kws-fixture"); err != nil {
		t.Fatalf("seed install: %v", err)
	}
	// Move the installed tree aside and address it as an override, so the store
	// copy is a different directory with a different stream set.
	override := filepath.Join(t.TempDir(), "override-kws")
	if err := os.Rename(filepath.Join(store, "kws-fixture"), override); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(store, "kws-fixture"), 0o755); err != nil {
		t.Fatal(err)
	}
	plantStream(t, filepath.Join(store, "kws-fixture", "model.onnx"), "only-in-store", "store-side")
	plantStream(t, filepath.Join(override, "model.onnx"), "only-in-override", "override-side")

	mgr, err := NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0,
		LocalOverride: map[string]string{"kws-fixture": override},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Fatalf("the override tree should verify (it holds exactly the manifest's files): %v", err)
	}

	gotDir, streams, err := mgr.UnnamedStreams("kws-fixture")
	if err != nil {
		t.Fatalf("ticket 357 override leg: %v", err)
	}
	if gotDir != override {
		t.Fatalf("ticket 357: the enumeration reports on %s, want the local_override directory %s. The guard "+
			"verifies the override and the CLI prints store=%s, so following the store here would file a "+
			"report about a tree nobody verified", gotDir, override, store)
	}
	if len(streams) != 1 || streams[0] != "model.onnx:only-in-override:$DATA" {
		t.Fatalf("ticket 357 override leg: got %v, want exactly [model.onnx:only-in-override:$DATA]. Seeing the "+
			"store-side stream means the wrong directory was enumerated; seeing both means the store leaked in", streams)
	}
}
