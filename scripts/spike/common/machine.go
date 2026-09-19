package common

import (
	"os"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var procGMSE2 = prock32.NewProc("GlobalMemoryStatusEx")

// MachineInfo is recorded once per JSON output so every number in the spike
// report is traceable to hardware/OS (acceptance criterion of ticket 02).
type MachineInfo struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"osName"`
	OSVersion    string `json:"osVersion"`
	OSBuild      string `json:"osBuild"`
	CPUName      string `json:"cpuName"`
	CPULogical   int    `json:"cpuLogical"`
	RAMTotalMB   uint64 `json:"ramTotalMB"`
	RAMAvailMB   uint64 `json:"ramAvailMB"`
	ProcessOwner string `json:"processOwner"`
	RecordedAt   string `json:"recordedAtUtc"`
}

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	UllTotalPhys         uint64
	UllAvailPhys         uint64
	UllTotalPageFile     uint64
	UllAvailPageFile     uint64
	UllTotalVirtual      uint64
	UllAvailVirtual      uint64
	UllAvailExtendedVirtual uint64
}

func MachineInfo() MachineInfo {
	mi := MachineInfo{RecordedAt: time.Now().UTC().Format(time.RFC3339)}
	if hn, err := os.Hostname(); err == nil {
		mi.Hostname = hn
	}
	if v := windows.RtlGetVersion(); v != nil {
		mi.OSVersion = strconv.Itoa(int(v.Major)) + "." + strconv.Itoa(int(v.Minor))
		mi.OSBuild = strconv.Itoa(int(v.BuildNumber))
		name := regStr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductName")
		dv := regStr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "DisplayVersion")
		if dv != "" {
			name += " " + dv
		}
		mi.OSName = name
	}
	mi.CPUName = regStr(`HARDWARE\DESCRIPTION\System\CentralProcessor\0`, "ProcessorNameString")
	mi.CPULogical = runtime.NumCPU()
	if ms, err := globalMemoryStatusEx(); err == nil {
		mi.RAMTotalMB = ms.UllTotalPhys / (1 << 20)
		mi.RAMAvailMB = ms.UllAvailPhys / (1 << 20)
	}
	mi.ProcessOwner = os.Getenv("USERNAME")
	return mi
}

func globalMemoryStatusEx() (*memoryStatusEx, error) {
	var ms memoryStatusEx
	ms.Length = uint32(unsafe.Sizeof(ms))
	r, _, err := procGMSE2.Call(uintptr(unsafe.Pointer(&ms)))
	if r == 0 {
		return nil, err
	}
	return &ms, nil
}

func regStr(keyPath, name string) string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	s, _, err := k.GetStringValue(name)
	if err != nil {
		return ""
	}
	return s
}
