//go:build windows

package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"
)

// The trash-is-honest proof, in the only place it can be made: fs.trash must
// leave an entry the Recycle Bin can see, and it must be the SHELL that wrote
// it. An unlink - the tempting implementation - cannot pass these.

func TestRecycleBinStructLayoutMatchesWin32(t *testing.T) {
	// SHFILEOPSTRUCTW on AMD64: HWND(8) UINT(8) LPCWSTR(8) LPCWSTR(8)
	// FILEOP_FLAGS(2+pad) BOOL(4) LPVOID(8) LPCWSTR(8) = 56. A wrong layout
	// means shell32 reads garbage out of our memory.
	var op shfileopstructW
	if got := unsafe.Sizeof(op); got != 56 {
		t.Fatalf("sizeof(SHFILEOPSTRUCTW) = %d, want 56", got)
	}
	if off := unsafe.Offsetof(op.fAnyOperationsAborted); off != 36 {
		t.Errorf("fAnyOperationsAborted sits at %d, want 36 (the WORD FILEOP_FLAGS must keep its 2-byte pad)", off)
	}
	if op.fFlags != 0 {
		t.Fatal("zero value must have no flags set")
	}
}

// TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate calls the backend
// directly (no gate, no bridge) so the assertion is about the shell API alone:
// the file is gone from its path, and the volume's bin holds a restore record
// naming it.
func TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate(t *testing.T) {
	if !recycleBinSupported() {
		t.Skip("shell32 recycle-bin entry point not resolvable on this machine")
	}
	root := tempRaw(t)
	target := filepath.Join(root, "evidence.txt")
	if err := os.WriteFile(target, []byte("recycle-bin evidence"), 0o600); err != nil {
		t.Fatal(err)
	}
	canon := mustCanonical(t, target)
	bin := recycleBinDir(canon)
	if bin == "" {
		t.Fatalf("no bin dir derived for %s", canon)
	}
	before, err := recycleRecordNames(bin)
	if err != nil {
		t.Fatal(err)
	}
	detail, err := shellTrash(canon)
	if err != nil {
		t.Fatalf("shellTrash: %v", err)
	}
	if existsFile(target) {
		t.Fatal("the item is still on disk after a successful trash")
	}
	if !strings.HasPrefix(filepath.Base(detail.Record), "$I") {
		t.Errorf("restore record %q is not an $I file", detail.Record)
	}

	// Re-read the record the tool claimed, straight off the volume.
	after, err := recycleRecordNames(bin)
	if err != nil {
		t.Fatal(err)
	}
	var found string
	for rel := range after {
		if before[rel] {
			continue
		}
		body, err := os.ReadFile(filepath.Join(bin, rel))
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(utf16BytesToString(body)), strings.ToLower(baseOf(canon))) {
			found = rel
			break
		}
	}
	if found == "" {
		t.Fatal("no new $I restore record on the volume names the trashed file: this is not a recycle-bin operation")
	}
}

// TestShellTrashRefusesAMissingPathWithoutDeletingAnything guards the "no
// fallback to unlink" rule from the other side: the call must fail, and it must
// fail before any destructive step.
func TestShellTrashRefusesAMissingPathWithoutDeletingAnything(t *testing.T) {
	root := tempRaw(t)
	missing := mustCanonical(t, filepath.Join(root, "ghost.txt"))
	if _, err := shellTrash(missing); err == nil {
		t.Fatal("trashing a path that is not there must not report success")
	}
	if existsFile(missing) {
		t.Fatal("the path exists after all")
	}
}
