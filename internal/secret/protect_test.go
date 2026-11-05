package secret

import (
	"bytes"
	"testing"
)

// The protector seam (protect.go) is exercised on BOTH platforms by the same
// three properties, because they are the ones the store's contract actually
// depends on: a round trip returns the original bytes, the blob on disk is not
// the plaintext, and a blob that is not the one protect wrote is refused
// instead of half-opened.
//
// On Windows these run DPAPI; elsewhere they run protect_other.go's seal. Both
// must answer alike - a mechanism that leaked the plaintext or returned a value
// from a damaged blob would be a C28 defect regardless of platform.

const protectorTestSecret = "sk-FAKE-protector-2ruabwxj"

// TestProtectorRoundTripAndOpaqueBlob is the seam's acceptance criterion 1.
func TestProtectorRoundTripAndOpaqueBlob(t *testing.T) {
	p := platformProtector(t.TempDir())
	plain := []byte(protectorTestSecret)

	blob, err := p.protect(plain)
	if err != nil {
		t.Fatalf("protect: %v", err)
	}
	if len(blob) == 0 {
		t.Fatal("protect returned an empty blob")
	}
	if bytes.Contains(blob, plain) {
		t.Fatal("the blob carries the plaintext secret")
	}
	if bytes.Contains(bytes.ToLower(blob), []byte("fake-protector")) {
		t.Fatal("the blob carries a recognizable fragment of the secret")
	}

	got, err := p.unprotect(blob)
	if err != nil {
		t.Fatalf("unprotect: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("round trip returned %d bytes, want the %d that went in", len(got), len(plain))
	}
}

// TestProtectorRefusesEmptyInput pins parity between the mechanisms: the store
// already refuses an empty secret one layer up, and both protectors refuse it
// too rather than inventing an on-disk representation for "nothing".
func TestProtectorRefusesEmptyInput(t *testing.T) {
	p := platformProtector(t.TempDir())
	if _, err := p.protect(nil); err == nil {
		t.Fatal("protect(nil) must be an error")
	}
	if _, err := p.unprotect(nil); err == nil {
		t.Fatal("unprotect(nil) must be an error")
	}
}

// TestProtectorRejectsATamperedBlob is the "never a fallback" property at its
// smallest: a flipped byte must fail loudly. On Windows DPAPI's integrity check
// catches it; on POSIX the GCM tag does. Neither may return a value.
func TestProtectorRejectsATamperedBlob(t *testing.T) {
	p := platformProtector(t.TempDir())
	blob, err := p.protect([]byte(protectorTestSecret))
	if err != nil {
		t.Fatalf("protect: %v", err)
	}
	tampered := append([]byte(nil), blob...)
	tampered[len(tampered)-1] ^= 0xFF
	if got, err := p.unprotect(tampered); err == nil {
		t.Fatalf("a tampered blob opened and returned %d bytes", len(got))
	}
}

// TestProtectorRejectsForeignBytes is the "a file that is not a blob" case the
// command surface hits when a blob was overwritten by hand (ticket 63's
// corrupted-blob case). Truncation and arbitrary bytes both have to be errors,
// and the error must not echo the bytes back.
func TestProtectorRejectsForeignBytes(t *testing.T) {
	p := platformProtector(t.TempDir())
	for _, foreign := range [][]byte{
		[]byte("not a blob at all"),
		{0x00},
		[]byte("blob from another machine"),
	} {
		got, err := p.unprotect(foreign)
		if err == nil {
			t.Fatalf("unprotect(%d foreign bytes) returned %d bytes", len(foreign), len(got))
		}
		if bytes.Contains([]byte(err.Error()), foreign) {
			t.Errorf("the error echoes the blob bytes back: %q", err)
		}
	}
}
