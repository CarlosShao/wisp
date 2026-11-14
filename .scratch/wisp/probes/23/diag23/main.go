//go:build windows

// diag23 is the read-only diagnosis for the elevated-target predicate shipped in
// internal/tools/system_backend_windows.go. It does not touch the product: it
// replicates the SHIPPED shape (integrityOf over a process handle) next to the
// corrected shape (OpenProcessToken first), on the same pids, in the same run.
//
// Usage: diag23.exe <pid> [pid...]
package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

const processQueryLimitedInfo = 0x1000

// shippedValue is the product's integrityOf applied to what realProcessElevated
// actually hands it: windows.Token(handle) where handle came from OpenProcess.
func shippedValue(process windows.Handle) (uint32, error) {
	var n uint32
	err := windows.GetTokenInformation(windows.Token(process), windows.TokenIntegrityLevel, nil, 0, &n)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmt.Errorf("size 0")
	}
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(windows.Token(process), windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0, err
	}
	label := (*windows.Tokenmandatorylabel)(unsafe.Pointer(&buf[0]))
	count := label.Label.Sid.SubAuthorityCount()
	if count == 0 {
		return 0, fmt.Errorf("no sub-authorities")
	}
	return label.Label.Sid.SubAuthority(uint32(count - 1)), nil
}

// correctedValue opens a real token handle first.
func correctedValue(pid uint32) (uint32, string, error) {
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return 0, "", err
	}
	defer windows.CloseHandle(h)
	var tok windows.Token
	if err := windows.OpenProcessToken(h, windows.TOKEN_QUERY, &tok); err != nil {
		return 0, "", err
	}
	defer tok.Close()
	var n uint32
	_ = windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, nil, 0, &n)
	if n == 0 {
		return 0, "", fmt.Errorf("sizing call returned no size")
	}
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0, "", err
	}
	label := (*windows.Tokenmandatorylabel)(unsafe.Pointer(&buf[0]))
	sidStr := label.Label.Sid.String()
	count := label.Label.Sid.SubAuthorityCount()
	return label.Label.Sid.SubAuthority(uint32(count - 1)), sidStr, nil
}

func selfIntegrity() uint32 {
	var tok windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &tok); err != nil {
		return 0
	}
	defer tok.Close()
	var n uint32
	_ = windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, nil, 0, &n)
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0
	}
	label := (*windows.Tokenmandatorylabel)(unsafe.Pointer(&buf[0]))
	return label.Label.Sid.SubAuthority(uint32(label.Label.Sid.SubAuthorityCount() - 1))
}

func procName(pid uint32) string {
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)
	var buf [windows.MAX_PATH + 1]uint16
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return ""
	}
	return windows.UTF16ToString(buf[:size])
}

func main() {
	self := selfIntegrity()
	fmt.Printf("self pid=%d integrity=0x%x\n", os.Getpid(), self)
	for _, a := range os.Args[1:] {
		pid64, err := strconv.ParseUint(a, 10, 32)
		if err != nil {
			fmt.Printf("bad pid %q\n", a)
			continue
		}
		pid := uint32(pid64)
		proc, perr := windows.OpenProcess(processQueryLimitedInfo, false, pid)
		var shipped uint32
		var serr error
		if perr == nil {
			shipped, serr = shippedValue(proc)
			windows.CloseHandle(proc)
		} else {
			serr = perr
		}
		corr, sidStr, cerr := correctedValue(pid)
		fmt.Printf("pid=%d name=%s shipped_integrity=%v shipped_err=%v corrected_integrity=0x%x corrected_sid=%s corrected_err=%v shipped_says_elevated=%v corrected_says_elevated=%v\n",
			pid, procName(pid), shipped, serr, corr, sidStr, cerr, shipped > self, corr > self)
	}
}
