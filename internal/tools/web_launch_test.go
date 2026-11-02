package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// web_launch_test.go — AC row 4 (the cross-platform shape half):
// "app.launch launches ... AUMID + name + path paths work." The real
// Job/notepad execution is pinned in app_launch_windows_test.go (the
// product platform per D7); this half proves the argv-vector rule (D33),
// the resolver seam (search.apps, ticket 23), the posture honesty, and
// the default-deny failure texts.

type fakeJobs struct {
	mu      sync.Mutex
	cmds    []*exec.Cmd
	started bool
	fail    error
}

func (f *fakeJobs) StartInJob(cmd *exec.Cmd) (*os.Process, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return nil, f.fail
	}
	f.cmds = append(f.cmds, cmd)
	f.started = true
	return nil, nil // tools discard the handle; no real start happens here
}

func (f *fakeJobs) last() *exec.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.cmds) == 0 {
		return nil
	}
	return f.cmds[len(f.cmds)-1]
}

type mapResolver map[string]string

func (m mapResolver) ResolveApp(name string) (string, bool) {
	p, ok := m[name]
	return p, ok
}

func Test22AppLaunchNameResolvesThroughSeam(t *testing.T) {
	fj := &fakeJobs{}
	w := &webShared{
		deps: WebDeps{
			AppResolver: mapResolver{"testapp": "/usr/bin/testapp22"},
			Jobs:        fj,
		},
		logf: func(string, ...any) {},
	}
	tool := appLaunch{w: w}
	res := execCall(t, tool, map[string]any{"app": "testapp", "args": []string{"--flag", "hello world"}})
	if res.IsError {
		t.Fatalf("name path: %s", res.Text)
	}
	cmd := fj.last()
	if cmd == nil {
		t.Fatal("no start recorded")
	}
	if cmd.Path != "/usr/bin/testapp22" {
		t.Errorf("cmd.Path = %q", cmd.Path)
	}
	// D33 forced-argv: the vector travels as Args, never a joined string.
	want := []string{"/usr/bin/testapp22", "--flag", "hello world"}
	if strings.Join(cmd.Args, "\x00") != strings.Join(want, "\x00") {
		t.Errorf("cmd.Args = %v want %v", cmd.Args, want)
	}
	if !strings.Contains(res.Text, "Job Object") {
		t.Errorf("posture text must report the Job membership that happened: %s", res.Text)
	}
}

func Test22AppLaunchPathAndMissingShapes(t *testing.T) {
	fj := &fakeJobs{}
	dir := t.TempDir()
	exe := filepath.Join(dir, "mytool.exe")
	if err := os.WriteFile(exe, []byte("MZ"), 0o700); err != nil {
		t.Fatal(err)
	}
	w := &webShared{deps: WebDeps{Jobs: fj}, logf: func(string, ...any) {}}
	res := execCall(t, appLaunch{w: w}, map[string]any{"app": exe})
	if res.IsError {
		t.Fatalf("direct path: %s", res.Text)
	}
	if fj.last() == nil || fj.last().Path != exe {
		t.Errorf("path target not started: %+v", fj.last())
	}
	res = execCall(t, appLaunch{w: w}, map[string]any{"app": filepath.Join(dir, "nope.exe")})
	if !res.IsError {
		t.Errorf("missing path target must refuse")
	}
	res = execCall(t, appLaunch{w: w}, map[string]any{"app": "definitely-not-an-app-22-zz"})
	if !res.IsError {
		t.Errorf("unresolvable name must refuse: %s", res.Text)
	}
	res = execCall(t, appLaunch{w: w}, map[string]any{})
	if !res.IsError {
		t.Errorf("missing app param must refuse")
	}
}

func Test22AppLaunchJobFailureSurfaces(t *testing.T) {
	fj := &fakeJobs{fail: os.ErrPermission}
	w := &webShared{deps: WebDeps{
		AppResolver: mapResolver{"x": "/bin/x"},
		Jobs:        fj,
	}, logf: func(string, ...any) {}}
	res := execCall(t, appLaunch{w: w}, map[string]any{"app": "x"})
	if !res.IsError || !strings.Contains(res.Text, "Job") {
		t.Errorf("StartInJob failure must surface as refusal, not half-truth: %s", res.Text)
	}
}

func Test22AppLaunchNoJobsStillHonest(t *testing.T) {
	dir := t.TempDir()
	// a script that exits immediately on every platform's shell reality:
	// use the go binary itself with "version" (instant, harmless).
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go binary not on PATH for the no-Job honesty fixture")
	}
	_ = dir
	w := &webShared{deps: WebDeps{AppResolver: mapResolver{}}, logf: func(string, ...any) {}}
	res := execCall(t, appLaunch{w: w}, map[string]any{"app": goBin, "args": []string{"version"}})
	if res.IsError {
		t.Fatalf("bare start: %s", res.Text)
	}
	if strings.Contains(res.Text, "已纳入 Job Object") {
		t.Errorf("must not claim a Job it did not join: %s", res.Text)
	}
	if !strings.Contains(res.Text, "未纳入") {
		t.Errorf("posture line missing the honest weaker statement: %s", res.Text)
	}
}

func Test22AUMIDIsWindowsOnlyHere(t *testing.T) {
	w := &webShared{deps: WebDeps{AppResolver: mapResolver{}}, logf: func(string, ...any) {}}
	res := execCall(t, appLaunch{w: w}, map[string]any{"app": "Microsoft.WindowsCalculator_8wekyb3d8bbwe!App"})
	if runtime.GOOS == "windows" {
		// classification into the aumid arm is pinned in the windows file;
		// here only assert the cross-platform contract: windows may run it.
		if res.IsError && strings.Contains(res.Text, "本平台不支持") {
			t.Errorf("windows must not answer the non-platform refusal")
		}
		return
	}
	if !res.IsError || !strings.Contains(res.Text, "不支持") {
		t.Errorf("non-windows must answer the platform refusal: %s", res.Text)
	}
}
