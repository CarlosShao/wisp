package secret

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExistsAndDelete pins the two read/remove primitives `wisp secret
// unset`/`set` are built on (ticket 63 adds no storage mechanism, only these
// command-shaped wrappers over the existing one-file-per-ref layout).
func TestExistsAndDelete(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	ref := NewRef()

	if ok, err := st.Exists(ref); err != nil || ok {
		t.Fatalf("Exists(before) = (%v, %v), want (false, nil)", ok, err)
	}
	if err := st.Store(ref, testSecret); err != nil {
		t.Fatalf("Store: %v", err)
	}
	if ok, err := st.Exists(ref); err != nil || !ok {
		t.Fatalf("Exists(after) = (%v, %v), want (true, nil)", ok, err)
	}
	if err := st.Delete(ref); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if ok, err := st.Exists(ref); err != nil || ok {
		t.Fatalf("Exists(after delete) = (%v, %v), want (false, nil)", ok, err)
	}
	if _, err := os.Stat(filepath.Join(st.Dir(), strings.TrimPrefix(ref, RefPrefixDPAPI))); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("blob file survived Delete (stat err = %v)", err)
	}
	if _, err := st.Resolve(ref); err == nil {
		t.Fatal("Resolve after Delete must fail")
	}

	// Deleting twice reports "no blob file" as fs.ErrNotExist (callers must be
	// able to tell already-gone from I/O failure).
	err = st.Delete(ref)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Delete(missing) err = %v, want a wrapped fs.ErrNotExist", err)
	}
	if strings.Contains(err.Error(), testSecret) {
		t.Error("delete error must not carry the secret")
	}

	// env: refs name process state, not a file.
	if err := st.Delete("env:WISP_SECRETSTORE_UNSET"); err == nil {
		t.Error("Delete(env ref) must fail")
	}
	if _, err := st.Exists("env:WISP_SECRETSTORE_UNSET"); err == nil {
		t.Error("Exists(env ref) must fail")
	}
	if _, err := st.Exists("sk-plaintext"); !errors.Is(err, ErrInvalidRef) {
		t.Errorf("Exists(invalid ref) err = %v, want ErrInvalidRef", err)
	}
}

// TestBlobPathRejectsTraversalOutsideCharset is the guard Exists/Delete rely
// on: an id that could escape the secrets dir never reaches the filesystem.
func TestBlobPathRejectsTraversalOutsideCharset(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	for _, bad := range []string{"../evil", "..", "a/b", `a\b`, "sp ace", ""} {
		ref := RefPrefixDPAPI + bad
		if err := st.Store(ref, "x"); !errors.Is(err, ErrInvalidRef) {
			t.Errorf("Store(%q) err = %v, want ErrInvalidRef", ref, err)
		}
		if err := st.Delete(ref); !errors.Is(err, ErrInvalidRef) {
			t.Errorf("Delete(%q) err = %v, want ErrInvalidRef", ref, err)
		}
		if _, err := st.Exists(ref); !errors.Is(err, ErrInvalidRef) {
			t.Errorf("Exists(%q) err = %v, want ErrInvalidRef", ref, err)
		}
	}
	entries, err := os.ReadDir(st.Dir())
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("rejected refs left %d files", len(entries))
	}
}

// TestBlobsListsMetadataOnly is the `wisp secret list` contract (ticket 63):
// ids and times, never content. The assertion is made on the serialized
// result, so a future field that starts carrying secret material fails here.
func TestBlobsListsMetadataOnly(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	const secretValue = "sk-bloblisting-should-never-appear-1234"
	// Names chosen so the sorted order is verifiable, and one stray non-ref
	// file that must be skipped rather than listed.
	for _, id := range []string{"zeta", "alpha", "10-middle"} {
		if err := st.Store(RefPrefixDPAPI+id, secretValue); err != nil {
			t.Fatalf("Store(%s): %v", id, err)
		}
	}
	if err := os.WriteFile(filepath.Join(st.Dir(), "not.a.ref!"), []byte("stray"), 0o600); err != nil {
		t.Fatalf("write stray file: %v", err)
	}

	blobs, err := st.Blobs()
	if err != nil {
		t.Fatalf("Blobs: %v", err)
	}
	if len(blobs) != 3 {
		t.Fatalf("Blobs returned %d entries, want the 3 valid blob ids (stray file skipped): %+v", len(blobs), blobs)
	}
	want := []string{"10-middle", "alpha", "zeta"}
	var got []string
	for i, b := range blobs {
		got = append(got, b.ID)
		if b.ID != want[i] {
			t.Errorf("blobs[%d].ID = %q, want %q (sorted by id)", i, b.ID, want[i])
		}
		if b.Ref != RefPrefixDPAPI+want[i] {
			t.Errorf("blobs[%d].Ref = %q, want %q", i, b.Ref, RefPrefixDPAPI+want[i])
		}
		if b.Created.IsZero() || b.Modified.IsZero() {
			t.Errorf("blobs[%d]: created=%v modified=%v, want real timestamps", i, b.Created, b.Modified)
		}
		if b.Size == 0 {
			t.Errorf("blobs[%d].Size = 0, want the encrypted blob size", i)
		}
	}

	dump, err := json.Marshal(blobs)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(dump), secretValue) {
		t.Errorf("Blobs metadata leaked the secret: %s", dump)
	}
	if strings.Contains(string(dump), "should-never-appear") {
		t.Errorf("Blobs metadata leaked a secret fragment: %s", dump)
	}
}

// TestBlobsEmptyAndMissingDir: an env with no secrets dir lists as empty, not
// as an error (a fresh WISP_ENV has nothing stored yet).
func TestBlobsEmptyAndMissingDir(t *testing.T) {
	root := t.TempDir()
	blobs, err := func() ([]BlobInfo, error) {
		st, err := NewStore(filepath.Join(root, "unused"))
		if err != nil {
			return nil, err
		}
		if err := os.Remove(filepath.Join(root, "unused", "secrets")); err != nil {
			return nil, err
		}
		return st.Blobs()
	}()
	if err != nil {
		t.Fatalf("Blobs on a missing dir: %v", err)
	}
	if len(blobs) != 0 {
		t.Fatalf("Blobs on a missing dir = %d entries, want 0", len(blobs))
	}
}

// TestConfigRefs is the reference-index view `wisp secret unset` uses to
// refuse deleting a referenced blob: dotted path -> ref for every
// api_key_ref, including nested tables, arrays of tables, and no config.
func TestConfigRefs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	const body = `schema_version = 2

[llm.providers.deepseek]
api_key_ref = "dpapi:deepseek"

[llm.providers.ollama]
base_url = "http://127.0.0.1:11434/v1"

[voice.realtime]
enabled = false
api_key_ref = "env:WISP_RT_KEY"

[[llm.providers.matrix.pool]]
api_key_ref = "dpapi:matrix-pool"
`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	refs, err := ConfigRefs(path)
	if err != nil {
		t.Fatalf("ConfigRefs: %v", err)
	}
	want := map[string]string{
		"llm.providers.deepseek.api_key_ref":       "dpapi:deepseek",
		"voice.realtime.api_key_ref":               "env:WISP_RT_KEY",
		"llm.providers.matrix.pool[0].api_key_ref": "dpapi:matrix-pool",
	}
	if len(refs) != len(want) {
		t.Fatalf("ConfigRefs returned %d entries, want %d: %+v", len(refs), len(want), refs)
	}
	for k, v := range want {
		if refs[k] != v {
			t.Errorf("ConfigRefs[%q] = %q, want %q", k, refs[k], v)
		}
	}

	names, err := RefFieldNames(path, "dpapi:deepseek")
	if err != nil {
		t.Fatalf("RefFieldNames: %v", err)
	}
	if len(names) != 1 || names[0] != "llm.providers.deepseek.api_key_ref" {
		t.Fatalf("RefFieldNames = %v, want the one referencing field", names)
	}
	if names, err := RefFieldNames(path, "dpapi:absent"); err != nil || len(names) != 0 {
		t.Fatalf("RefFieldNames(absent) = (%v, %v), want no hits", names, err)
	}

	// No config yet (fresh env): empty, not an error.
	missing, err := ConfigRefs(filepath.Join(dir, "nope.toml"))
	if err != nil || len(missing) != 0 {
		t.Fatalf("ConfigRefs(missing file) = (%v, %v), want empty and no error", missing, err)
	}

	// A corrupt config must surface, not read as "nothing references it".
	bad := filepath.Join(dir, "bad.toml")
	if err := os.WriteFile(bad, []byte("this is = = not toml"), 0o600); err != nil {
		t.Fatalf("write bad config: %v", err)
	}
	if _, err := ConfigRefs(bad); err == nil {
		t.Fatal("ConfigRefs on an unparsable config must fail (unset fails closed)")
	}
}
