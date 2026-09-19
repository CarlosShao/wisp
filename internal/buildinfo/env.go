package buildinfo

import (
	"fmt"
	"os"
)

// Env is the typed WISP_ENV value (SPEC-03 §5.1): prod | dev | test.
type Env string

// The three environments.
const (
	EnvProd Env = "prod"
	EnvDev  Env = "dev"
	EnvTest Env = "test"
)

// ParseEnv parses a WISP_ENV value strictly (no defaults applied). Empty or
// unknown values are rejected; use ResolveEnv for the defaulting order.
func ParseEnv(s string) (Env, error) {
	switch Env(s) {
	case EnvProd, EnvDev, EnvTest:
		return Env(s), nil
	case "":
		return "", fmt.Errorf("buildinfo: WISP_ENV is empty, want one of prod|dev|test")
	default:
		return "", fmt.Errorf("buildinfo: invalid WISP_ENV %q, want one of prod|dev|test", s)
	}
}

// ResolveEnv resolves the effective environment: the WISP_ENV environment
// variable wins over the build-time default (SPEC-03 §5.1). DefaultEnv is a
// string because scripts/build.ps1 injects it via -ldflags -X, which only
// supports string variables.
func ResolveEnv() (Env, error) {
	if e := os.Getenv("WISP_ENV"); e != "" {
		return ParseEnv(e)
	}
	return ParseEnv(DefaultEnv)
}
