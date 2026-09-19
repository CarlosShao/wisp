package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sys/windows"
)

// check is one doctor result line. Status is "PASS", "FAIL" or "INFO".
type check struct {
	name     string
	status   string
	detail   string
	critical bool
}

// cmdDoctor runs the self-check (SPEC-11 §7.1 startup checks; ticket 01 scope:
// toolchain/DLL versions and the colocated-DLL rule, PASS/FAIL output).
// Returns true when no critical check failed.
func cmdDoctor() bool {
	var results []check
	pass := func(name, detail string) {
		results = append(results, check{name: name, status: "PASS", detail: detail, critical: true})
	}
	fail := func(name, detail string) {
		results = append(results, check{name: name, status: "FAIL", detail: detail, critical: true})
	}
	info := func(name, detail string) {
		results = append(results, check{name: name, status: "INFO", detail: detail, critical: false})
	}

	// Build identity.
	info("wisp build", fmt.Sprintf("version=%s commit=%s built=%s WISP_ENV=%s",
		buildinfo.Version, buildinfo.Commit, buildinfo.BuildDate, buildinfo.EnvString()))
	info("Go runtime", fmt.Sprintf("%s (toolchain pinned by go.mod)", runtime.Version()))
	info("C29 minisign public key", buildinfo.MinisignPublicKey+
		" (placeholder until C29 lands; hardcoded into buildinfo per SPEC-11 §7.3)")

	// gcc is a build-time-only dependency; report it for dev-machine diagnosis.
	if v, err := gccVersion(); err != nil {
		fail("gcc (build-time)", "not found: "+err.Error())
	} else {
		pass("gcc (build-time)", v)
	}

	// Native DLLs must be colocated with the exe (SPEC-11 §7.1). This binary
	// links the sherpa-onnx C API directly, so when the DLLs are missing the
	// OS loader refuses to start the process at all - any printed line here is
	// therefore already loader-level proof of the colocated-DLL rule.
	exeDir := executableDir()
	sherpaVer := sherpa.GetVersion()
	switch {
	case sherpaVer == "":
		fail("sherpa-onnx C API", "linked but returned an empty version string")
	case sherpaVer != buildinfo.SherpaOnnxVersion:
		fail("sherpa-onnx C API", fmt.Sprintf("runtime %s does not match build pin %s",
			sherpaVer, buildinfo.SherpaOnnxVersion))
	default:
		pass("sherpa-onnx C API", fmt.Sprintf("runtime %s matches build pin (exe dir %s)", sherpaVer, exeDir))
	}

	ortPath := filepath.Join(exeDir, "onnxruntime.dll")
	switch v, ok := dllFileVersion(ortPath); {
	case !ok:
		fail("onnxruntime.dll colocated", fmt.Sprintf("missing or unreadable next to exe: %s", ortPath))
	case normalizeVer(v) != normalizeVer(buildinfo.OnnxRuntimeVersion):
		fail("onnxruntime.dll version", fmt.Sprintf("file version %s does not match build pin %s",
			v, buildinfo.OnnxRuntimeVersion))
	default:
		pass("onnxruntime.dll colocated", fmt.Sprintf("file version %s matches build pin %s", v, buildinfo.OnnxRuntimeVersion))
	}
	info("sherpa-onnx built against onnxruntime", sherpa.GetOnnxruntimeVersion())

	for _, dll := range []string{"sherpa-onnx-c-api.dll", "sherpa-onnx-cxx-api.dll"} {
		p := filepath.Join(exeDir, dll)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			pass("DLL colocated: "+dll, p)
		} else {
			fail("DLL colocated: "+dll, fmt.Sprintf("not found next to exe: %s", p))
		}
	}

	// deps.toml cross-check: found when running from a repo layout; installed
	// copies have no deps.toml and rely on the build-time pins instead.
	if parsed, path, found := findDepsToml(exeDir); !found {
		info("deps.toml cross-check", "deps.toml not found near the exe (expected for installed copies; build-time pins already verified above)")
	} else {
		depsCrossCheck(pass, fail, *parsed, path)
	}

	// Data dir writable (SPEC-03 §5.2).
	env := buildinfo.EnvString()
	dir := resolveDataDir(env)
	if err := probeWritable(dir); err != nil {
		fail("data dir writable ("+env+")", dir+": "+err.Error())
	} else {
		pass("data dir writable ("+env+")", dir)
	}

	// Report.
	fmt.Println("wisp doctor - build chain self-check")
	fmt.Println(strings.Repeat("-", 72))
	anyFail := false
	for _, c := range results {
		fmt.Printf("%-6s %-34s %s\n", "["+c.status+"]", c.name, c.detail)
		if c.status == "FAIL" && c.critical {
			anyFail = true
		}
	}
	fmt.Println(strings.Repeat("-", 72))
	if anyFail {
		fmt.Println("wisp doctor: FAIL")
		return false
	}
	fmt.Println("wisp doctor: PASS")
	return true
}

func depsCrossCheck(pass, fail func(string, string), parsed depsToml, path string) {
	if parsed.SherpaOnnx.Version == "" {
		fail("deps.toml sherpa-onnx pin", path+": version missing")
	} else if parsed.SherpaOnnx.Version == buildinfo.SherpaOnnxVersion {
		pass("deps.toml sherpa-onnx pin", parsed.SherpaOnnx.Version)
	} else {
		fail("deps.toml sherpa-onnx pin", fmt.Sprintf("deps.toml says %s, binary pinned %s",
			parsed.SherpaOnnx.Version, buildinfo.SherpaOnnxVersion))
	}

	if parsed.OnnxRuntime.Version == "" {
		fail("deps.toml onnxruntime pin", path+": version missing")
	} else if normalizeVer(parsed.OnnxRuntime.Version) == normalizeVer(buildinfo.OnnxRuntimeVersion) {
		pass("deps.toml onnxruntime pin", parsed.OnnxRuntime.Version)
	} else {
		fail("deps.toml onnxruntime pin", fmt.Sprintf("deps.toml says %s, binary pinned %s",
			parsed.OnnxRuntime.Version, buildinfo.OnnxRuntimeVersion))
	}

	if parsed.GoBinding.Version == "" {
		fail("deps.toml sherpa-onnx-go pin", path+": version missing")
	} else if v, ok := buildinfo.DependencyVersion(parsed.GoBinding.Module); !ok {
		fail("deps.toml sherpa-onnx-go pin", fmt.Sprintf("module %s not found in build info", parsed.GoBinding.Module))
	} else if v == parsed.GoBinding.Version {
		pass("deps.toml sherpa-onnx-go pin", v)
	} else {
		fail("deps.toml sherpa-onnx-go pin", fmt.Sprintf("deps.toml says %s, binary linked %s",
			parsed.GoBinding.Version, v))
	}
}

// depsToml mirrors the deps.toml sections doctor consumes.
type depsToml struct {
	SherpaOnnx struct {
		Version string            `toml:"version"`
		Source  string            `toml:"source"`
		Dll     map[string]dllPin `toml:"dll"`
	} `toml:"sherpa-onnx"`
	OnnxRuntime struct {
		Version string `toml:"version"`
	} `toml:"onnxruntime"`
	GoBinding struct {
		Module  string `toml:"module"`
		Version string `toml:"version"`
	} `toml:"go-binding-sherpa-onnx"`
}

type dllPin struct {
	ArchivePath string `toml:"archive_path"`
	SHA256      string `toml:"sha256"`
}

// findDepsToml looks next to the exe, then up the directory tree, then in the
// working directory; the first deps.toml wins.
func findDepsToml(exeDir string) (*depsToml, string, bool) {
	var dirs []string
	if exeDir != "" {
		dirs = append(dirs, exeDir)
		for d := filepath.Dir(exeDir); ; d = filepath.Dir(d) {
			dirs = append(dirs, d)
			if parent := filepath.Dir(d); parent == d {
				break
			}
		}
	}
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, wd)
	}
	seen := map[string]bool{}
	for _, d := range dirs {
		if seen[d] {
			continue
		}
		seen[d] = true
		p := filepath.Join(d, "deps.toml")
		fi, err := os.Stat(p)
		if err != nil || fi.IsDir() {
			continue
		}
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var parsed depsToml
		if err := toml.Unmarshal(b, &parsed); err != nil {
			continue
		}
		return &parsed, p, true
	}
	return nil, "", false
}

// resolveDataDir implements the SPEC-03 §5.2 data dir rules (ticket 01 scope:
// portable mode + the per-env appdata fork; mutex/log/endpoint forks are later
// tickets).
func resolveDataDir(env string) string {
	if exeDir := executableDir(); exeDir != "" {
		if _, err := os.Stat(filepath.Join(exeDir, "portable.txt")); err == nil {
			if env == "dev" {
				return filepath.Join(exeDir, "data-dev")
			}
			return filepath.Join(exeDir, "data")
		}
	}
	if env == "test" {
		if v := os.Getenv("WISP_TEST_DATA_DIR"); v != "" {
			return v
		}
		return filepath.Join(os.TempDir(), fmt.Sprintf("wisp-test-%d", os.Getpid()))
	}
	base, err := os.UserConfigDir() // %APPDATA%
	if err != nil {
		base = "."
	}
	if env == "dev" {
		return filepath.Join(base, "wisp-dev")
	}
	return filepath.Join(base, "wisp")
}

func probeWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	probe := filepath.Join(dir, "doctor-write-probe.tmp")
	if err := os.WriteFile(probe, []byte("probe"), 0o644); err != nil {
		return err
	}
	return os.Remove(probe)
}

func dataDirForDisplay() string {
	// Prefer the real resolution order (per-env fork + portable override,
	// ticket 06); fall back to the static fork table when env/layout
	// resolution fails.
	if env, err := buildinfo.ResolveEnv(); err == nil {
		if sum, err := proc.Summarize(env); err == nil && sum.DataDir != "" {
			return sum.DataDir
		}
	}
	return resolveDataDir(buildinfo.EnvString())
}

func gccVersion() (string, error) {
	cc := os.Getenv("CC")
	if cc == "" {
		cc = "gcc"
	}
	out, err := exec.Command(cc, "--version").Output()
	if err != nil {
		return "", err
	}
	line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
	return strings.TrimSpace(line), nil
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exe)
}

// normalizeVer trims trailing ".0" components so "1.28.2.0" equals "1.28.2".
func normalizeVer(v string) string {
	parts := strings.Split(strings.TrimSpace(v), ".")
	for len(parts) > 1 && parts[len(parts)-1] == "0" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, ".")
}

// dllFileVersion reads the VS_FIXEDFILEINFO FileVersion of a PE file through
// the Win32 version API.
var (
	modVersion              = windows.NewLazySystemDLL("version.dll")
	procGetFileVersionInfoW = modVersion.NewProc("GetFileVersionInfoW")
	procGetFileVersionSizeW = modVersion.NewProc("GetFileVersionInfoSizeW")
	procVerQueryValueW      = modVersion.NewProc("VerQueryValueW")
)

type vsFixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

func dllFileVersion(path string) (string, bool) {
	p16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", false
	}
	size, _, _ := procGetFileVersionSizeW.Call(uintptr(unsafe.Pointer(p16)), 0)
	if size == 0 {
		return "", false
	}
	buf := make([]byte, size)
	r1, _, _ := procGetFileVersionInfoW.Call(uintptr(unsafe.Pointer(p16)), 0, size, uintptr(unsafe.Pointer(&buf[0])))
	if r1 == 0 {
		return "", false
	}
	root, err := windows.UTF16PtrFromString("\\")
	if err != nil {
		return "", false
	}
	var ffi *vsFixedFileInfo
	var ffiLen uint32
	r1, _, _ = procVerQueryValueW.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(root)),
		uintptr(unsafe.Pointer(&ffi)), uintptr(unsafe.Pointer(&ffiLen)))
	if r1 == 0 || ffi == nil || ffi.Signature != 0xFEEF04BD {
		return "", false
	}
	return fmt.Sprintf("%d.%d.%d.%d",
		ffi.FileVersionMS>>16, ffi.FileVersionMS&0xffff,
		ffi.FileVersionLS>>16, ffi.FileVersionLS&0xffff), true
}
