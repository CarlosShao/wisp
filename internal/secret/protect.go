package secret

// protector is the platform half of the C28 store: the operation that turns a
// plaintext secret into the bytes of one blob file, and back again.
//
// C28 names one mechanism per SUPPORTED platform - Windows DPAPI
// (CryptProtectData, CurrentUser scope) and, at S8, the macOS Keychain - so the
// mechanism is a build-tagged choice, not a runtime option: dpapi_windows.go
// supplies the Windows protector and protect_other.go supplies the non-Windows
// one. D7 puts Linux outside the support matrix ("源码可自编、CI 可编译，但不进
// 官方支持矩阵"), so the non-Windows protector has no C28 sentence to satisfy;
// it exists so the store's platform-neutral invariants (one opaque file per ref,
// no plaintext on disk, an explicit error instead of a fallback) are executable
// on the POSIX leg rather than unreachable there.
//
// The interface is unexported on purpose. It is the store's internal seam, not
// a caller knob: keeping "DPAPI on Windows" a property of the BUILD (rather
// than of an option some caller passes) is what makes this a seam inside C28
// instead of a new C28 surface.
type protector interface {
	protect(plaintext []byte) ([]byte, error)
	unprotect(blob []byte) ([]byte, error)
}

// platformProtector returns the protector for a store whose blob directory is
// dir. dir is part of the seam because a platform protector may scope its key
// to the store - the non-Windows one does, so a blob carried out of its store
// stops opening, the way a DPAPI blob stops opening under another user. The
// Windows protector scopes to the Windows user and ignores dir, which is its
// behaviour before and after this seam (C28).
//
// The body lives in the build-tagged files (Go allows an external-body-only
// declaration for assembly, not for a second Go file).
