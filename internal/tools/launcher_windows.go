//go:build windows

package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// launcher_windows.go is the OS hand-off side of the ticket-22 open family.
//
//   - web.open / file.open legs ride ShellExecuteW "open" to the DEFAULT
//     handler the USER registered. The scheme decision was already made
//     (16.5.3 matrix in web_open.go + R5 facts in the bridge); this file
//     cannot widen it, it only delivers an accepted target.
//   - app.launch path/name legs go through CreateProcess via exec.Cmd and
//     ENTER the C30 Job Object (WebDeps.Jobs) before returning. argv is a
//     vector (D33 forced-argv rule: nothing is string-joined into a shell
//     command line here either).
//   - AUMID legs activate through explorer's shell:AppsFolder hand-off.
//     Honest platform fact: explorer delegating to the AppX activation
//     server means the UWP's own process is not a child of this command —
//     it cannot be assigned to our Job at creation time. The explorer
//     carrier enters the Job; the activated app lives under the OS
//     activation host. This is recorded in the launch result text and in
//     the ticket-22 evidence, not papered over.

// windowsAUMIDPattern is the conservative shape gate for !Name!App ids
// (package family name + app id, no path characters).
var windowsAUMIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{1,255}![A-Za-z0-9._-]{1,255}$`)

// platformLaunch delivers one accepted web/file/uri open via ShellExecute.
// kind was already vetted by the tool ("url"|"file"|"uri"); the target is
// the canonical/normalized string the gate saw. WebDeps.Launch (dispatcher
// in web.go) wins first when a test or host wires its own hand-off.
func (w *webShared) platformLaunch(_ context.Context, kind, target string) error {
	_ = kind // the matrix already decided; this is delivery, not judgment
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	if err := windows.ShellExecute(0, verb, file, nil, nil, windows.SW_SHOWNORMAL); err != nil {
		return fmt.Errorf("ShellExecute open 失败：%w", err)
	}
	w.logf("tools: %s open delivered via ShellExecute (target scheme/host only in logs, not full query where avoidable)", kind)
	return nil
}

// isKnownMsProtocolDefault answers "registered ms-* application protocol"
// from the per-machine Classes Root, the same hive the shell itself reads.
// Failure to open = not registered = the deny side (default-deny, 16.5.3).
func isKnownMsProtocolDefault(scheme string) bool {
	if !strings.HasPrefix(scheme, "ms-") {
		return false
	}
	k, err := registry.OpenKey(registry.CLASSES_ROOT, scheme, registry.READ)
	if err != nil {
		return false
	}
	_ = k.Close()
	return true
}

// resolveApp maps one app.launch target into a launchTarget.
// Order: explicit path shape, AUMID shape, then the name resolver
// (AppResolver seam, else App Paths registry, else PATH lookup).
func (w *webShared) resolveApp(_ context.Context, name string) (launchTarget, error) {
	if looksLikeDirectPath(name) {
		if _, err := os.Stat(name); err != nil {
			return launchTarget{}, fmt.Errorf("路径目标不存在或不可达：%v", err)
		}
		return launchTarget{Kind: "path", Path: name}, nil
	}
	if windowsAUMIDPattern.MatchString(name) {
		return launchTarget{Kind: "aumid", AUMID: name}, nil
	}
	res := w.deps.AppResolver
	if res == nil {
		res = windowsAppPathsResolver{}
	}
	if p, ok := res.ResolveApp(name); ok {
		return launchTarget{Kind: "path", Path: p}, nil
	}
	if p, err := exec.LookPath(name); err == nil {
		return launchTarget{Kind: "path", Path: p}, nil
	}
	if p, err := exec.LookPath(name + ".exe"); err == nil {
		return launchTarget{Kind: "path", Path: p}, nil
	}
	return launchTarget{}, fmt.Errorf("找不到应用 %q（search.apps 检索面尚未接入，App Paths 与 PATH 均未命中）", name)
}

// startApp lands the child in the C30 Job when wired and reports the
// POSTURE THAT ACTUALLY HAPPENED in one honest line, so the tool result
// never claims a Job membership stronger than the platform delivers (the
// AUMID arm's carrier enters the Job; the activated UWP process itself is
// hosted by the OS activation server — file header, recorded not hidden).
func (w *webShared) startApp(_ context.Context, target launchTarget, args []string) (string, error) {
	switch target.Kind {
	case "path":
		cmd := exec.Command(target.Path, args...)
		return w.startInJobOrBare(cmd, "进程")
	case "aumid":
		cmd := exec.Command("explorer.exe", `shell:AppsFolder\`+target.AUMID)
		posture, err := w.startInJobOrBare(cmd, "启动载体")
		if err != nil {
			return posture, err
		}
		if w.deps.Jobs != nil {
			return posture + "；UWP 应用本体由系统激活宿主承载，不在 wisp Job 树内（平台事实）", nil
		}
		return posture, nil
	}
	return "", errors.New("未知的启动目标类型")
}

func (w *webShared) startInJobOrBare(cmd *exec.Cmd, what string) (string, error) {
	if w.deps.Jobs != nil {
		if _, err := w.deps.Jobs.StartInJob(cmd); err != nil {
			return "", fmt.Errorf("进入 Job Object 失败（C30）：%w", err)
		}
		return what + "已纳入 Job Object（C30）", nil
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	w.logf("tools: app.launch started OUTSIDE a wisp Job (WebDeps.Jobs unwired) — weaker C30 posture, recorded")
	return what + "未纳入 wisp Job（C30 未接线，弱姿态）", nil
}

// windowsAppPathsResolver reads the two App Paths keys the shell itself
// consults (HKLM then HKCU), trying "name" then "name.exe".
type windowsAppPathsResolver struct{}

// ResolveApp implements AppResolver.
func (windowsAppPathsResolver) ResolveApp(name string) (string, bool) {
	if name == "" || strings.ContainsAny(name, `/\\*?"<>|`) {
		return "", false
	}
	candidates := []string{name}
	if !strings.HasSuffix(strings.ToLower(name), ".exe") {
		candidates = append(candidates, name+".exe")
	}
	for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		for _, cand := range candidates {
			key := `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\` + cand
			k, err := registry.OpenKey(root, key, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			v, _, err := k.GetStringValue("")
			_ = k.Close()
			if err != nil || v == "" {
				continue
			}
			// The registry value may carry arguments after the path; the
			// shell splits at the first space outside quotes — keep the
			// conservative side: take the quoted value or the bare value.
			p := strings.Trim(v, `"`)
			if i := strings.LastIndex(p, `.exe"`); i >= 0 && strings.Contains(p, `"`) {
				p = strings.Trim(p[:i+4], `"`)
			}
			if _, err := os.Stat(p); err == nil {
				return p, true
			}
		}
	}
	return "", false
}
