//go:build !windows

package secret

import (
	"bytes"
	"testing"
)

// These pin the non-Windows seal's own contract (protect_other.go). They are
// tagged to the platform whose mechanism they describe - unlike the case files
// ticket 70 moved out of the ubuntu denominator, they are not store tests that
// were made unrunnable; nothing here needs DPAPI, and nothing here duplicates
// what protect_test.go already asks both mechanisms.

// TestPOSIXSealIsBoundToItsStore pins the one scope property the seal claims:
// a blob carried out of the store that wrote it stops opening (the DPAPI
// analogue is a blob opened under another user or machine, P13). Without this,
// the "key is derived from machine + directory" sentence would be a comment
// nobody measures.
func TestPOSIXSealIsBoundToItsStore(t *testing.T) {
	writer := platformProtector(t.TempDir())
	other := platformProtector(t.TempDir())

	blob, err := writer.protect([]byte(protectorTestSecret))
	if err != nil {
		t.Fatalf("protect: %v", err)
	}
	if got, err := other.unprotect(blob); err == nil {
		t.Fatalf("a blob opened under a different store and returned %d bytes", len(got))
	}
	// Positive control: the same blob under its own store still opens, so the
	// refusal above is about the store and not about a seal that never works.
	if got, err := writer.unprotect(blob); err != nil || string(got) != protectorTestSecret {
		t.Fatalf("the writing store could not open its own blob (err=%v)", err)
	}
}

// TestPOSIXSealUsesAFreshNonce pins that two seals of one plaintext differ.
// The command surface depends on this: ticket 63 AC#4 reads three envs' blobs
// back and requires three DISTINCT ciphertexts, which a deterministic seal
// would satisfy only by accident of different keys.
func TestPOSIXSealUsesAFreshNonce(t *testing.T) {
	p := platformProtector(t.TempDir())
	first, err := p.protect([]byte(protectorTestSecret))
	if err != nil {
		t.Fatalf("first protect: %v", err)
	}
	second, err := p.protect([]byte(protectorTestSecret))
	if err != nil {
		t.Fatalf("second protect: %v", err)
	}
	if bytes.Equal(first, second) {
		t.Fatal("two seals of the same plaintext produced identical blobs")
	}
	for _, blob := range [][]byte{first, second} {
		if got, err := p.unprotect(blob); err != nil || string(got) != protectorTestSecret {
			t.Fatalf("a blob sealed by this protector did not open (err=%v)", err)
		}
	}
}

// TestPOSIXSealNamesTheShapeOfForeignBytes pins the error text the store wraps:
// garbage must be refused by shape (version/size), before any AEAD work, and
// the message must not carry the bytes. The corrupted-blob command case
// (ticket 63) reads this message through `wisp secret get`, so it has to be
// stable and free of the file's contents.
func TestPOSIXSealNamesTheShapeOfForeignBytes(t *testing.T) {
	p := platformProtector(t.TempDir())
	foreign := []byte("not a wisp seal at all")
	_, err := p.unprotect(foreign)
	if err == nil {
		t.Fatal("garbage opened")
	}
	if bytes.Contains([]byte(err.Error()), foreign) {
		t.Errorf("the error echoes the bytes back: %q", err)
	}
}
