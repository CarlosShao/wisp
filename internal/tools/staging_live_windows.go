//go:build windows

package tools

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

// stillActive is Win32's STILL_ACTIVE exit code (winsafer.h), the value
// GetExitCodeProcess reports for a process that has not terminated. x/sys
// exports the error constants but not this one.
const stillActive = 259

// stagingCreatorAlive answers the one question that separates an orphan from a
// concurrent writer's in-flight staging file: is the process named in the file
// name still running?
//
// PROCESS_QUERY_LIMITED_INFORMATION is the least-privilege probe: it can ask,
// and nothing else. "No such process" is the only answer that reports dead;
// every other failure (including ERROR_ACCESS_DENIED on another user's process)
// is returned as an error, which the sweeper treats as ALIVE so the fold is
// always toward deleting nothing.
func stagingCreatorAlive(pid string) (bool, error) {
	n, err := strconv.Atoi(pid)
	if err != nil || n <= 0 {
		return false, fmt.Errorf("tools: 暂存文件名里的 pid %q 不是进程号", pid)
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(n))
	if err != nil {
		if errors.Is(err, syscall.Errno(windows.ERROR_INVALID_PARAMETER)) {
			return false, nil // the creator is gone: this is an orphan
		}
		return false, err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return false, err
	}
	return code == stillActive, nil
}

// stagingIsReparse is C26's rule applied to the sweeper's own traversal: an
// entry carrying a reparse point is judged as ITSELF, never through whatever it
// reaches. Deleting what sits beyond a junction would turn a 4-byte cleanup into
// a delete primitive pointed outside the allowed dirs.
func stagingIsReparse(st os.FileInfo) bool {
	data, ok := st.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		// Unknown shape: be conservative in the direction that deletes nothing.
		return true
	}
	return data.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}
