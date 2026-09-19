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
// binary (SPEC-11 §7.3). The real keypair is generated offline when ticket C29
// lands; key material must never be committed to the repo, hence placeholder.
const MinisignPublicKey = "PLACEHOLDER-C29-MINISIGN-PUBLIC-KEY"

// Env resolves the effective WISP_ENV: the environment variable wins over the
// build-time default (SPEC-03 §5.1).
func Env() string {
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
