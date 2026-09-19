// Package buildinfo carries values injected at build time (SPEC-11 §2.2) plus
// the hardcoded C29 minisign public key (SPEC-11 §7.3).
//
// scripts/build.ps1 overrides these via
// -ldflags "-X github.com/CarlosShao/wisp/internal/buildinfo.Name=value".
// Native dependency versions are parsed from deps.toml at build time so that
// the binary itself is the runtime self-check comparison source (SPEC-11 §7.2).
package buildinfo

import (
	"os"
	"runtime/debug"
)

// Build-time injected values. Defaults are for "go run" / non-script builds.
var (
	Version            = "0.0.0-dev" // product version; scheme decided at release time
	Commit             = "unknown"   // git short hash of the built tree
	BuildDate          = "unknown"   // UTC build timestamp
	DefaultEnv         = "dev"       // WISP_ENV when the env var is unset (SPEC-03 §5.1)
	SherpaOnnxVersion  = "unknown"   // pinned in deps.toml
	OnnxRuntimeVersion = "unknown"   // pinned in deps.toml
)

// MinisignPublicKey is the C29 update/model signing public key embedded in the
// binary (SPEC-11 §7.3), in minisign two-line format (keyid a3c8794f3fd94fc5).
// It is the DEV keypair generated for ticket 14: the secret key lives
// out-of-repo at E:\work\base\wisp-minisign\wisp-models.key (never committed;
// see docs/reports/2026-09-19-t14-dev-minisign-key.md). The PRODUCTION key
// ceremony is ticket S8 and MUST rotate this value.
const MinisignPublicKey = `untrusted comment: wisp models signing key (dev)
RWSjyHlPP9lPxdEQRvWj3zFLMbc1tTEkKMwTDuVXXQDxsWRpA/m5jk9j`

// EnvString resolves the effective WISP_ENV as a raw string: the environment
// variable wins over the build-time default (SPEC-03 §5.1). Prefer ResolveEnv
// for the typed, strictly-validated form (env.go).
func EnvString() string {
	if e := os.Getenv("WISP_ENV"); e != "" {
		return e
	}
	return DefaultEnv
}

// DependencyVersion returns the version of a Go module dependency linked into
// this binary, as recorded by the Go toolchain.
func DependencyVersion(modulePath string) (string, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, dep := range bi.Deps {
		if dep.Path == modulePath {
			return dep.Version, true
		}
	}
	return "", false
}
