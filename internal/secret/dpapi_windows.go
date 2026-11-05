//go:build windows

package secret

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// protect encrypts plaintext with DPAPI in CurrentUser scope (C28, SPEC-03
// §5.3): the blob can only be decrypted by the same Windows user on the same
// machine - that scope is the real access control, the blob file mode is
// defense in depth. CRYPTPROTECT_UI_FORBIDDEN keeps the call silent
// (headless/service safe).
func protect(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, errors.New("secret: dpapi protect: empty plaintext")
	}
	in := windows.DataBlob{Size: uint32(len(plaintext)), Data: &plaintext[0]}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, fmt.Errorf("secret: dpapi CryptProtectData: %w", err)
	}
	return takeBlob(&out)
}

// unprotect decrypts a blob produced by protect. Failing here is the P13
// path (different user/machine, corrupted blob): the caller decides whether
// that surfaces as the explicit portable-mode error or a plain failure.
func unprotect(blob []byte) ([]byte, error) {
	if len(blob) == 0 {
		return nil, errors.New("secret: dpapi unprotect: empty blob")
	}
	in := windows.DataBlob{Size: uint32(len(blob)), Data: &blob[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, fmt.Errorf("secret: dpapi CryptUnprotectData: %w", err)
	}
	return takeBlob(&out)
}

// dpapiProtector is C28's Windows mechanism behind the protector seam
// (protect.go). It is a pure delegation to the two functions above: the DPAPI
// calls, their flags, their buffer handling and their error strings are the
// ones that were here before the seam, so the bytes on disk and the messages
// callers see are unchanged. The seam adds a call, not a mechanism.
type dpapiProtector struct{}

func (dpapiProtector) protect(plaintext []byte) ([]byte, error) { return protect(plaintext) }

func (dpapiProtector) unprotect(blob []byte) ([]byte, error) { return unprotect(blob) }

// platformProtector is the Windows half of protect.go's seam. DPAPI's scope is
// the Windows user (CurrentUser), not the directory, so dir is deliberately
// unused: two stores under the same user still share one crypto scope, exactly
// as they did before this seam existed.
func platformProtector(string) protector { return dpapiProtector{} }

// takeBlob copies the DPAPI output buffer into Go-owned memory and frees the
// LocalAlloc'd original (the caller owns it per the DPAPI contract).
func takeBlob(out *windows.DataBlob) ([]byte, error) {
	if out.Size == 0 || out.Data == nil {
		return nil, errors.New("secret: dpapi returned an empty blob")
	}
	copied := make([]byte, out.Size)
	copy(copied, unsafe.Slice(out.Data, out.Size))
	if _, err := windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data))); err != nil {
		// The secret material is already copied; a failed free must not fail
		// the operation, but it is worth knowing about in diagnostics.
		return copied, nil
	}
	return copied, nil
}
