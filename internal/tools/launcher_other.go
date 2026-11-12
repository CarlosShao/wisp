//go:build !windows

package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// launcher_other.go is the compile-portability face of the open family on
// the Linux CI runner (test-core). D7 scopes the PRODUCT to Windows, so
// these legs answer honestly "this platform does not have that door"
// instead of pretending: no ShellExecute equivalent, no Job Object, no
// ms-* protocol registry. The file.open scheme matrix (web_open.go) still
// runs — the default-DENY verdict is platform-independent and is what the
// protocol-matrix acceptance row pins.

func (w *webShared) platformLaunch(_ context.Context, kind, target string) error {
	// A permissive desktop opener where one exists; refuse elsewhere.
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", target)
		}
	case "darwin":
		cmd = exec.Command("open", target)
	}
	if cmd == nil {
		return errors.New("本平台（" + runtime.GOOS + "）没有接好的系统默认程序通道；本产品按 D7 以 Windows 为目标平台")
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	w.logf("tools: %s open delivered via %s (no wisp Job on this platform)", kind, cmd.Path)
	return nil
}

func isKnownMsProtocolDefault(scheme string) bool {
	// No Windows protocol registration exists here; the allowlist therefore
	// lands on its deny side for every ms-* target — the correct fail-closed
	// answer, not a missing check.
	_ = scheme
	return false
}

func (w *webShared) resolveApp(_ context.Context, name string) (launchTarget, error) {
	if looksLikeDirectPath(name) {
		if _, err := os.Stat(name); err != nil {
			return launchTarget{}, err
		}
		return launchTarget{Kind: "path", Path: name}, nil
	}
	// Same shape gate as launcher_windows.go (launcher_aumid.go): an AUMID
	// must classify as Kind "aumid" HERE, before the name resolvers, so the
	// startApp aumid arm answers the honest platform refusal instead of the
	// name resolvers answering "找不到应用" for a shape they cannot resolve.
	if windowsAUMIDPattern.MatchString(name) {
		return launchTarget{Kind: "aumid", AUMID: name}, nil
	}
	if res := w.deps.AppResolver; res != nil {
		if p, ok := res.ResolveApp(name); ok {
			return launchTarget{Kind: "path", Path: p}, nil
		}
	}
	if p, err := exec.LookPath(name); err == nil {
		return launchTarget{Kind: "path", Path: p}, nil
	}
	return launchTarget{}, errors.New("找不到应用：" + name)
}

func (w *webShared) startApp(_ context.Context, target launchTarget, args []string) (string, error) {
	switch target.Kind {
	case "path":
		cmd := exec.Command(target.Path, args...)
		if w.deps.Jobs != nil {
			if _, err := w.deps.Jobs.StartInJob(cmd); err != nil {
				// Same failure shape as the Windows side (startInJobOrBare):
				// the refusal must name the Job it failed to enter, never a
				// bare transport error that reads as a half-truth.
				return "", fmt.Errorf("进入 Job Object 失败（C30）：%w", err)
			}
			return "进程已纳入 Job Object（C30）", nil
		}
		if err := cmd.Start(); err != nil {
			return "", err
		}
		w.logf("tools: app.launch started OUTSIDE a wisp Job (no Job implementation on this platform)")
		return "进程未纳入 wisp Job（本平台无 Job 实现，弱姿态）", nil
	case "aumid":
		return "", errors.New("UWP/AUMID 启动是 Windows 平台能力（D7），本平台不支持")
	}
	return "", errors.New("未知的启动目标类型")
}
