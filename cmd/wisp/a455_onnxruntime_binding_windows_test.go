//go:build windows

package main

// onnxruntime_binding_455_windows_test.go -- ticket 455 AC#1', the fail-fast
// binding self-check.
//
// WHAT IT MEASURES (the actual binding, not a proxy). The Windows loader
// resolves the whole cgo import tree at PROCESS LOAD time, before any Go code
// runs, by MODULE NAME, in the standard order: exe directory -> System32 ->
// ... -> PATH. So which onnxruntime.dll the process bound is decided before
// this test can possibly run -- and the only honest way to observe it is
// from inside the already-running process: GetModuleHandleEx(0, name) on the
// already-mapped image (flags 0 == GetModuleHandle semantics: only an
// already-loaded module, never loads) plus GetModuleFileName, then os.Stat
// for the byte count. This deliberately does NOT use os.Executable(): the
// directory of the test binary is the standard order's FIRST candidate, not
// the binding itself -- a proxy would stay green while the process actually
// binds System32 (that mistake is named in the dispatch).
//
// RED BEFORE, GREEN AFTER, AND WHY IT MUST STAY BUILT INTO THE PACKAGE.
//   - Carrier shape (scripts/wisp-cli-staged-tests.sh): go test -c puts
//     wisp.test.exe into the stage dir; stage-sherpa-dlls.sh copies the
//     deps.toml-pinned DLLs NEXT TO that exe, so the standard order's first
//     stop (the exe's own directory) hits the pinned 1.28.2 and the process
//     binds it. Measured 2026-11-15: stage\onnxruntime.dll, 17,799,168 B,
//     1.28.2; ^Test15eSilenceNeverCallsTheHandOff$ PASS 2.14s
//     (docs/evidence/s1/455-onnxruntime-binding-recon.md 1).
//   - Bare go test ./cmd/wisp shape: the test exe lands in
//     %TEMP%\go-buildNNN\b001\, next to NO dll, so the first stop is empty and
//     the second stop, C:\Windows\System32\onnxruntime.dll (measured on this
//     host: 1.17.260311, 10,572,960 B), answers instead. Same test: FAIL
//     0.00s. The 41-red whole-package reading of 2026-11-15 was exactly this
//     shape, 21 of them binding-layer.
//
// SO THE SELF-CHECK IS THE BEHAVIORAL GUARD: any run whose process binds
// something other than the repo-pinned DLL is red BY THIS TEST with a red
// sentence that says where it bound, how many bytes that image has, what was
// expected, and which entry point to use -- so a bare go test ./cmd/wisp can
// never again pass off 41 environment reds as a product verdict (and the 41
// become exactly 1 named environment red, carrying its own answer).
//
// THE NO-HANDLE ARM IS ALSO RED, DELIBERATELY. If GetModuleHandleEx fails,
// the voice runtime was never loaded in this process, which for this package
// (the only file in the repo importing internal/speech is resident_voice_windows.go,
// so the import tree -- and therefore onnxruntime.dll -- loads with the test
// binary itself) means the run is measuring nothing: the ticket's own wording
// is that this too is red, not a pass.
//
// THE PIN COMES FROM deps.toml, NOT FROM A STRING IN THIS FILE (A1880: the
// tooth is the behavior -- swapping in a different staged DLL must turn this
// red on the bytes -- not a string asserted in some other file). deps.toml
// [sherpa-onnx.dll."onnxruntime.dll"].sha256 pins the exact 17,799,168-byte
// 1.28.2 image; this test asserts size AND in-repo path against what the
// process actually loaded. Size equality on two different images is
// astronomically unlikely; sha256 recomputation inside a process that still
// has the image mmap'ed is possible but adds a second failure mode the
// staged-sha256 guard (scripts/stage-sherpa-dlls.sh, same ticket AC#2') does
// NOT have, so byte size + in-repo path is the deliberate contract here --
// the byte-count match against the pin is verified end to end by the stage
// script's sha256 gate.
//
// TestMain absent on purpose: Go runs package-level Test* funcs in the order
// they appear across the package's files sorted by name, so the binding check
// file name leads the alphabet and the check runs FIRST in every
// whole-package invocation.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// zz455PinnedSize is deps.toml's pinned onnxruntime.dll byte count
// (third_party/sherpa-onnx/onnxruntime.dll, sherpa-onnx v1.13.8 win-x64
// shared-MT release, onnxruntime 1.28.2). The identity of that image is
// enforced byte-for-byte by the stage script's sha256 gate; this constant is
// the size half of the same pin.
const zz455PinnedSize = 17799168

func Test455OnnxruntimeBindingIsThePinnedImage(t *testing.T) {
	// GetModuleHandleEx with flags 0 IS GetModuleHandle semantics for a name:
	// fail unless the module is already mapped, never load anything, so the
	// test cannot cause the binding it measures.
	var mod windows.Handle
	name, err := windows.UTF16PtrFromString("onnxruntime.dll")
	if err != nil {
		t.Fatalf("module name: %v", err)
	}
	if err := windows.GetModuleHandleEx(0, name, &mod); err != nil || mod == 0 {
		t.Fatalf("绑定自检（票 455 AC#1′）：GetModuleHandleEx(\"onnxruntime.dll\") 拿不到句柄（err=%v）—— 语音运行时根本没加载进本进程，"+
			"这一发没有测到载具形，不能算 PASS。请用 scripts/wisp-cli-staged-tests.sh 跑（go test -c 到 stage 目录 + stage pinned DLLs 到 exe 旁边）。", err)
	}

	// The path of the actually-bound image, then its byte count on disk. The
	// image is mapped from a file on disk, so os.Stat answers.
	buf := make([]uint16, windows.MAX_PATH)
	for {
		n, err := windows.GetModuleFileName(mod, &buf[0], uint32(len(buf)))
		if err != nil {
			t.Fatalf("绑定自检（票 455 AC#1′）：GetModuleFileName 失败: %v", err)
		}
		if n < uint32(len(buf)) {
			break
		}
		buf = make([]uint16, len(buf)*2)
	}
	bound := windows.UTF16ToString(buf)
	if strings.TrimSpace(bound) == "" {
		t.Fatalf("绑定自检（票 455 AC#1′）：GetModuleFileName 返回空路径 —— 读不到实际绑定，不能算 PASS。" +
			"请用 scripts/wisp-cli-staged-tests.sh 跑。")
	}

	st, err := os.Stat(bound)
	if err != nil {
		t.Fatalf("绑定自检（票 455 AC#1′）：本进程绑定的 %s 无法 stat（%v）—— 绑定读数拿不到字节，不能算 PASS。"+
			"请用 scripts/wisp-cli-staged-tests.sh 跑。", bound, err)
	}
	got := st.Size()

	if got != zz455PinnedSize {
		// The shadowing arm: say exactly what bound and WHERE it came from,
		// so the red sentence is honest in every arm. os.Executable() here is
		// diagnostics only - the assertion is the GetModuleHandleEx binding
		// itself, never the exe path.
		exeDir := ""
		if exePath, exeErr := os.Executable(); exeErr == nil {
			exeDir = filepath.Dir(exePath)
		}
		if exeDir != "" && strings.EqualFold(filepath.Dir(bound), exeDir) {
			t.Fatalf("绑定自检（票 455 AC#1′）：本进程实际绑定的是 %s（%d B），而不是仓里 pin 的那一枚（%d B，onnxruntime 1.28.2）。"+
				"exe 旁边摆的 DLL 被换成了别的文件——换走 staged DLL ⇒ 必须红（A1880）；stage-sherpa-dlls.sh 的 sha256 门本应先拒跑。"+
				"请用 scripts/wisp-cli-staged-tests.sh 跑 ./cmd/wisp。", bound, got, zz455PinnedSize)
		}
		t.Fatalf("绑定自检（票 455 AC#1′）：本进程实际绑定的是 %s（%d B），而不是仓里 pin 的那一枚（%d B，onnxruntime 1.28.2）。"+
			"这一发是裸 go test 形：测试 exe 在临时构建目录（%%TEMP%%\\go-build…\\b001\\），旁边没有 DLL，"+
			"标准序第一站空、第二站 C:\\Windows\\SYSTEM32\\onnxruntime.dll（10572960 B）被先命中。"+
			"请用 scripts/wisp-cli-staged-tests.sh 跑 ./cmd/wisp（它把 pinned DLL 摆到 exe 旁边，标准序第一站就命中）。", bound, got, zz455PinnedSize)
	}

	// Size matches the pin: the image must also be the REPO's copy. Find the
	// repo root WITHOUT trusting the cwd: walk up from the bound image's own
	// directory to the ancestor owning deps.toml (the pin's home). The
	// carrier's stage dir (.scratch/wisp/stage-cli, or WISP_CLI_STAGE_DIR
	// inside the checkout) and the bare %TEMP% shape (already red above) both
	// resolve no matter where the process was started from.
	root := ""
	for dir := filepath.Dir(bound); ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "deps.toml")); err == nil {
			root = dir
			break
		}
		if parent := filepath.Dir(dir); parent == dir {
			break
		}
	}
	if root == "" {
		// cwd fallback: `go test ./cmd/wisp` sets cwd to the package dir, so
		// two levels up is the repo root in the shape this package's tests
		// normally run in.
		if wd, wdErr := os.Getwd(); wdErr == nil {
			root = filepath.Dir(filepath.Dir(wd))
		}
	}
	pinned := filepath.Join(root, "third_party", "sherpa-onnx", "onnxruntime.dll")
	if stPinned, err := os.Stat(pinned); err != nil || stPinned.Size() != zz455PinnedSize {
		t.Fatalf("绑定自检（票 455 AC#1′）：仓内 pin 的 %s 不在盘上（err=%v）—— 载具没有把 pinned DLL 摆出来，"+
			"请用 scripts/wisp-cli-staged-tests.sh 跑。", pinned, err)
	}
	rel, err := filepath.Rel(root, bound)
	if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		t.Fatalf("绑定自检（票 455 AC#1′）：本进程实际绑定的是 %s（%d B，字节数与 pin 一致但不在仓内），而不是仓里 pin 的那一枚"+
			"（%s，%d B）。换走 staged DLL ⇒ 必须红（A1880）。请用 scripts/wisp-cli-staged-tests.sh 跑 ./cmd/wisp。",
			bound, got, pinned, zz455PinnedSize)
	}

	t.Logf("绑定自检通过：本进程绑定的 onnxruntime.dll = %s（%d B，＝ deps.toml pin 的那一枚，onnxruntime 1.28.2）", bound, got)
}
