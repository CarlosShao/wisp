//go:build !windows

package tools

import (
	"errors"
	"os"
)

// Non-Windows fallbacks for the staging sweeper. Wisp's supported platform is
// Windows (SPEC-01), and this package already refuses the platform-specific
// halves there (see platform_other.go: fs.trash has no bin, so it refuses
// instead of unlinking). The sweeper follows the same shape: it does not guess
// at liveness on a platform whose process semantics this ticket never measured.
//
// Reporting "undecidable" makes stagingCreatorAlive fail closed, i.e. every
// attributable orphan is SKIPPED rather than deleted. So the non-Windows build
// can never lose a live writer's file - and never reclaims one either, which is
// the honest asymmetry: the proof of this behavior lives in a Windows-only test.

func stagingCreatorAlive(string) (bool, error) {
	return false, errors.New("tools: 本平台无法判定进程是否存活，暂存文件一律不清扫")
}

// stagingIsReparse has no reparse-point concept off Windows; Lstat's mode check
// in stagingDeletable still blocks anything that is not a plain regular file.
func stagingIsReparse(os.FileInfo) bool { return false }
