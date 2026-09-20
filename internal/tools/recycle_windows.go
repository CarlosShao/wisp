//go:build windows

package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The real shell recycle-bin API (D34 note②: trash is L1 BECAUSE the shell can
// give the item back). fs.trash calls SHFileOperationW with FOF_ALLOWUNDO - the
// same entry point Explorer's Delete key uses - and then PROVES the item landed
// by reading the restore record the shell writes into the volume's
// `$Recycle.Bin\<SID>\$I*` metadata.
//
// os.Remove is never reached from this file: an unlink is not a trash, and a
// tool that unlinks while declaring L1 is the exact lie the ticket forbids. If
// the shell refuses, or if the restore record cannot be found afterwards, this
// returns an error and fs.trash fails - and it has already deleted nothing,
// because there is no code path here that deletes without the bin.
//
// SHQueryRecycleBinW is deliberately NOT used: from a Go thread it blocks
// forever (measured - the test binary sat in that syscall until the go-test
// timeout killed it), and a background tool call must not be able to hang a
// task. The metadata scan below is a plain directory read: same proof, no COM
// round trip.

const (
	foDelete = 0x0003 // FO_DELETE

	fofSILENT         = 0x0004
	fofNoCONFIRMATION = 0x0010
	fofALLOWUNDO      = 0x0040
	fofNOCANCELLATION = 0x0080
	fofNOERRORUI      = 0x0400
	// FOF_WANTNUKEWARNING is NOT set: it asks shell32 for a MODAL "this cannot
	// be recovered" dialog, which turns a background tool call into a window
	// nobody told the user about. The restore-record check answers the same
	// question with data.

	// cocinitAPARTMENTTHREADED / cocinitDisableOLE1DDE / sFalse: shell32's file
	// operations marshal through COM, which wants an apartment on the CALLING
	// THREAD; S_FALSE means this thread already had one and must still uninit.
	cocinitAPARTMENTTHREADED = 0x00000002
	cocinitDisableOLE1DDE    = 0x00000004
	sFalse                   = 1
)

var (
	procSHFileOperationW = windows.NewLazySystemDLL("shell32.dll").NewProc("SHFileOperationW")
	procCoInitializeEx   = windows.NewLazySystemDLL("ole32.dll").NewProc("CoInitializeEx")
	procCoUninitialize   = windows.NewLazySystemDLL("ole32.dll").NewProc("CoUninitialize")
)

// shfileopstructW mirrors the Win32 SHFILEOPSTRUCTW. FILEOP_FLAGS is a WORD, not
// a DWORD, and the field order/widths below are what the OS reads; the layout is
// pinned by TestRecycleBinStructLayoutMatchesWin32.
type shfileopstructW struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	_                     [2]byte // C alignment padding after the WORD
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

// trashDetail is what the shell backend reports about one successful call, so
// the applied-steps line carries evidence instead of an unqualified "done".
type trashDetail struct {
	// API names the shell entry point that was used.
	API string
	// Record is the restore record ($I file) the shell wrote for this item.
	Record string
	// Bin is the volume's recycle-bin folder the record landed in.
	Bin string
}

// recycleBinSupported is the answer fs.trash gives before it touches anything:
// no shell entry point means no trash, and the tool refuses instead of deleting.
func recycleBinSupported() bool {
	// LazyProc.Call PANICS when the export cannot be resolved, and a panicking
	// tool would take the task down with it; resolve up front so a missing
	// export degrades into "no recycle bin here" instead.
	if err := procSHFileOperationW.Find(); err != nil {
		return false
	}
	return procCoInitializeEx.Find() == nil
}

// shellTrash moves one canonical path into the Recycle Bin. A returned error
// means the item was NOT moved - and nothing was deleted either.
func shellTrash(canonical string) (trashDetail, error) {
	d := trashDetail{API: "SHFileOperationW(FOF_ALLOWUNDO)"}
	if strings.TrimSpace(canonical) == "" {
		return d, errors.New("回收站：空路径")
	}
	if !recycleBinSupported() {
		return d, errors.New("shell32/ole32 的回收站入口不可用，已拒绝执行（绝不退化成删除）")
	}
	if !existsViaLstatQuiet(canonical) {
		return d, fmt.Errorf("回收站：目标不存在，未做任何操作: %s", canonical)
	}
	bin := recycleBinDir(canonical)
	if bin == "" {
		return d, errors.New("回收站：无法定位该卷的 $Recycle.Bin，拒绝按已回收处理")
	}
	before, err := recycleRecordNames(bin)
	if err != nil {
		return d, err
	}

	// COM's apartment is per OS THREAD, and shell32 marshals through it; a Go
	// goroutine can be rescheduled off its thread mid-call, so the thread is
	// pinned for the whole operation. An uninitialized apartment is what made
	// the earlier shell query hang forever.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	coInit, _, _ := procCoInitializeEx.Call(0, uintptr(cocinitAPARTMENTTHREADED|cocinitDisableOLE1DDE))
	if coInit == 0 || coInit == uintptr(sFalse) {
		defer procCoUninitialize.Call()
	}

	from := utf16DoubleZ(canonical)
	op := &shfileopstructW{
		wFunc: foDelete,
		pFrom: &from[0],
		// SILENT + NOCONFIRMATION + NOERRORUI + NOCANCELLATION: the gate already
		// asked the human, so a second shell dialog would duplicate a decision
		// the shell does not own - and a modal dialog from a background call is a
		// task that hangs in silence.
		fFlags: fofALLOWUNDO | fofSILENT | fofNoCONFIRMATION | fofNOERRORUI | fofNOCANCELLATION,
	}
	r, _, sc := procSHFileOperationW.Call(uintptr(unsafe.Pointer(op)))
	runtime.KeepAlive(from)
	runtime.KeepAlive(op)
	if r != 0 {
		return d, fmt.Errorf("SHFileOperationW 返回 %d（%v）：项目未进入回收站，也未被删除", int32(r), sc)
	}
	if op.fAnyOperationsAborted != 0 {
		return d, errors.New("Shell 操作被中止：项目未进入回收站，也未被删除")
	}

	// THE PROOF. The shell writes one `$I` restore record per recycled item,
	// and its body carries the original path; a new record that names this path
	// is direct evidence the item can be given back. An unlink cannot produce
	// this, so a "trash" implemented as a delete fails right here.
	record, err := findRecycleRecord(bin, before, canonical)
	if err != nil {
		return d, err
	}
	d.Record = filepath.Base(record)
	d.Bin = bin
	return d, nil
}

// existsViaLstatQuiet is the pre-flight existence check for the shell call,
// where a probe error counts as absent (the shell would fail on it anyway, and
// the answer fs.trash gives must not be "we deleted something we could not see").
func existsViaLstatQuiet(p string) bool {
	ok, err := existsViaLstat(p)
	return err == nil && ok
}

// recycleBinDir is `<drive>\$Recycle.Bin` for one canonical path. A UNC target
// has no local bin, which shellTrash turns into a refusal.
func recycleBinDir(canonical string) string {
	root := volumeRootOf(canonical)
	if len(root) < 2 || root[1] != ':' {
		return ""
	}
	return root + `$Recycle.Bin`
}

// recycleRecordNames lists every `$I*` restore record currently on this volume,
// across every per-user bin folder we can read. It is a name-only enumeration
// (one directory pass, no stat), which is what makes it safe to call twice
// around a trash even on a bin with tens of thousands of entries.
func recycleRecordNames(bin string) (map[string]bool, error) {
	out := map[string]bool{}
	sids, err := os.ReadDir(bin)
	if err != nil {
		return nil, fmt.Errorf("无法读取回收站目录 %s：%w", bin, err)
	}
	for _, sid := range sids {
		if !sid.IsDir() {
			continue
		}
		entries, err := os.ReadDir(filepath.Join(bin, sid.Name()))
		if err != nil {
			// Another user's bin (access denied) or a folder that vanished mid
			// scan: skipping it is safe, because the record our own delete
			// writes always lands in OUR sid folder, which is readable.
			if errors.Is(err, os.ErrPermission) {
				continue
			}
			return nil, fmt.Errorf("无法读取回收站子目录：%w", err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasPrefix(e.Name(), "$I") {
				out[filepath.Join(sid.Name(), e.Name())] = true
			}
		}
	}
	return out, nil
}

// findRecycleRecord returns the path of a NEW `$I` restore record whose body
// names the trashed item. "New" is taken against the before-snapshot so the
// scan does not have to trust the clock.
func findRecycleRecord(bin string, before map[string]bool, canonical string) (string, error) {
	now, err := recycleRecordNames(bin)
	if err != nil {
		return "", err
	}
	base := strings.ToLower(baseOf(canonical))
	parent := strings.ToLower(dirOf(canonical))
	var fresh int
	for rel := range now {
		if before[rel] {
			continue
		}
		fresh++
		full := filepath.Join(bin, rel)
		body, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		// The record stores the original path as UTF-16LE; decoding the whole
		// blob and matching folded substrings is enough to tell "this item"
		// from "some other item deleted at the same moment".
		txt := strings.ToLower(utf16BytesToString(body))
		if strings.Contains(txt, base) && (parent == "" || strings.Contains(txt, parent)) {
			return full, nil
		}
	}
	if fresh == 0 {
		return "", fmt.Errorf("回收站中没有出现任何新的还原记录（%s）：拒绝按已回收处理", bin)
	}
	return "", fmt.Errorf("回收站新增了 %d 条还原记录，但没有一条指向 %s：拒绝按已回收处理", fresh, canonical)
}

// utf16BytesToString decodes a UTF-16LE blob, dropping the binary header
// (which decodes to unprintable runes, harmless for a substring search).
func utf16BytesToString(b []byte) string {
	u := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		u = append(u, uint16(b[i])|uint16(b[i+1])<<8)
	}
	return string(utf16.Decode(u))
}

// utf16DoubleZ encodes a shell file-list entry: the path, then TWO NUL units,
// because SHFileOperationW reads pFrom as a double-terminated list. It is NOT
// windows.StringToUTF16, which panics on an embedded NUL.
func utf16DoubleZ(p string) []uint16 {
	return append(utf16.Encode([]rune(p)), 0, 0)
}
