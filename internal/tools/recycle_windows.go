//go:build windows

package tools

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The real shell recycle-bin API (D34 note②: trash is L1 BECAUSE the shell can
// put the item back). fs.trash goes through SHFileOperationW with
// FOF_ALLOWUNDO - the same call Explorer's Delete key makes - and the entry is
// then verified by the drive's item count via SHQueryRecycleBinW.
//
// os.Remove is deliberately never reached from this file: an unlink is not a
// trash, and a tool that unlinks while declaring L1 would be the exact lie the
// ticket forbids. Where the shell refuses to take the item, this returns an
// error and fs.trash fails without deleting anything either.

const (
	foDelete = 0x0000

	fofSILENT          = 0x0004
	fofNoCONFIRMATION  = 0x0010
	fofALLOWUNDO       = 0x0040
	fofNOERRORUI       = 0x0400
	fofWANTNUKEWARNING = 0x4000
)

// shell32 is loaded once; the procs are resolved lazily on first use and a
// missing export surfaces as a trash refusal, never as a fallback delete.
var (
	procSHFileOperationW   = windows.NewLazySystemDLL("shell32.dll").NewProc("SHFileOperationW")
	procSHQueryRecycleBinW = windows.NewLazySystemDLL("shell32.dll").NewProc("SHQueryRecycleBinW")
)

// shfileopstructW mirrors the Win32 SHFILEOPSTRUCTW exactly. FILEOP_FLAGS is a
// WORD, not a DWORD, and Go's automatic field padding reproduces the C layout
// (checked against unsafe.Sizeof by TestRecycleBinStructLayoutMatchesWin32).
type shfileopstructW struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

// shqueryrbinfoW mirrors SHQUERYRBINFO (DWORD cbSize + two __int64 counters).
type shqueryrbinfoW struct {
	cbSize      uint32
	i64Size     int64
	i64NumItems int64
}

// trashDetail is what the shell backend reports about one successful call, so
// the tool's applied-steps line can carry the evidence rather than an
// unqualified "done".
type trashDetail struct {
	// ItemsBefore/ItemsAfter are that drive's recycle-bin entry counts around
	// the call. ItemsAfter >= ItemsBefore+1 is the proof the item landed.
	ItemsBefore int64
	ItemsAfter  int64
	// API names the shell entry point that was used, for the audit line.
	API string
}

// recycleBinSupported is the platform answer fs.trash gives before it does
// anything else. On Windows the API exists, so refusing here would be a lie.
func recycleBinSupported() bool {
	// LazyProc.Call PANICS when the export cannot be resolved, and a panicking
	// tool would take the task down with it. Resolve both procs up front so a
	// missing export degrades into "no recycle bin here" instead.
	if err := procSHFileOperationW.Find(); err != nil {
		return false
	}
	return procSHQueryRecycleBinW.Find() == nil
}

// shellTrash moves one canonical path into the Recycle Bin. A returned error
// means the item was NOT moved (and nothing was deleted either).
func shellTrash(canonical string) (trashDetail, error) {
	d := trashDetail{API: "SHFileOperationW(FOF_ALLOWUNDO)"}
	if strings.TrimSpace(canonical) == "" {
		return d, errors.New("回收站：空路径")
	}
	if !recycleBinSupported() {
		return d, errors.New("shell32 的回收站入口不可用，已拒绝执行（绝不退化成删除）")
	}
	root := volumeRootOf(canonical)
	before, err := recycleBinCount(root)
	if err != nil {
		return d, err
	}
	d.ItemsBefore = before

	from := utf16DoubleZ(canonical)
	op := &shfileopstructW{
		wFunc: foDelete,
		pFrom: &from[0],
		// SILENT + NOCONFIRMATION + NOERRORUI: the gate already asked the
		// human, so a second modal shell dialog would duplicate a decision the
		// UI does not own. WANTNUKEWARNING stays ON so an item the bin cannot
		// hold is reported as the refusal it is instead of being destroyed.
		fFlags: fofALLOWUNDO | fofSILENT | fofNoCONFIRMATION | fofNOERRORUI | fofWANTNUKEWARNING,
	}
	r, _, sc := procSHFileOperationW.Call(uintptr(unsafe.Pointer(op)))
	runtime.KeepAlive(from)
	runtime.KeepAlive(op)
	if r != 0 {
		return d, fmt.Errorf("SHFileOperationW 返回 %d（%v）：项目未进入回收站", int32(r), sc)
	}
	if op.fAnyOperationsAborted != 0 {
		return d, errors.New("Shell 操作被中止，项目未进入回收站")
	}
	after, err := recycleBinCount(root)
	if err != nil {
		return d, err
	}
	d.ItemsAfter = after
	if after < before+1 {
		// The count did not grow: the shell did not take the item (a drive with
		// its recycle bin disabled deletes instead). Reporting success here
		// would be the lie this function exists to prevent, so the tool fails
		// and says so.
		return d, fmt.Errorf("回收站条目数未增加（%d -> %d）：拒绝按已回收处理", before, after)
	}
	return d, nil
}

// recycleBinCount asks the shell how many entries the bin of one volume root
// holds. It is a read-only query, so it is safe to run before and after a trash.
func recycleBinCount(root string) (int64, error) {
	var info shqueryrbinfoW
	info.cbSize = uint32(unsafe.Sizeof(info))
	var p *uint16
	if root != "" {
		b := utf16DoubleZ(root)
		p = &b[0]
		defer runtime.KeepAlive(b)
	}
	r, _, sc := procSHQueryRecycleBinW.Call(uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&info)))
	if r != 0 {
		return 0, fmt.Errorf("SHQueryRecycleBinW(%q) 返回 0x%x（%v）", root, uint32(r), sc)
	}
	return info.i64NumItems, nil
}

// utf16DoubleZ encodes a shell file-list entry: the path, then TWO terminators,
// because SHFileOperationW reads pFrom as a list.
func utf16DoubleZ(p string) []uint16 {
	return windows.StringToUTF16(p + "\x00")
}
