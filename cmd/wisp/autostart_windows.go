//go:build windows

package main

// Ticket 43, autostart: the [app].autostart opt-in as a registry Run key.
//
// The contract pieces this file assembles, each already frozen elsewhere:
//
//   - the KEY is [app].autostart (internal/config/schema.go, default false -
//     the OFF default the ticket's fresh-config cell asks for is the schema's
//     own zero value, not a second default table here);
//   - the TIER is restart (internal/config/tiers.go, TierRegistry
//     "app.autostart" -> TierRestart, ticket 233): the OS registration is a
//     NEXT-BOOT property, so this file applies the config's wish at BOOT and
//     nowhere else. A toggle that wrote the key immediately would race the
//     Manager's restart-pending receipt and make "需要重启进程" a lie.
//   - the VALUE is HKCU\...\Run: HKCU needs no elevation (the resident process
//     never elevates), and the entry points at the running executable quoted,
//     because the install root may contain spaces.
//
// What this file does NOT do, said out loud: the tray-menu half of the ticket's
// "toggle in tray menu + config" lives in internal/ball's frozen menu table,
// which is not this leg's write surface. The config half is done (any config
// writer - CLI, panel settings, a hand edit - moves the key); the tray item is
// named as a residual in the evidence table rather than half-wired here.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
	"golang.org/x/sys/windows/registry"
)

// autostartRunKeyPath is the per-user Run key. HKCU on purpose: autostart is a
// per-session opt-in of a per-session process (single instance is per-session
// too, D42#7), and HKLM would need an elevation this process never performs.
const autostartRunKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`

// autostartRunValueName is this installation's value inside the Run key. It is
// a parameter of autostartSync (the assembly root passes autostartRunValueName;
// tests pass their own name) so a verification run can write, read and remove a
// value of its own without touching the owner's real entry.
const autostartRunValueName = "Wisp"

// autostartSync reads config.toml once and makes the Run key agree with
// [app].autostart. The four answers are distinct on purpose - "wrote it",
// "removed it", "already agreed (on)", and "agreed by absence (off)" - because
// the boot report prints the answer, and a boot that silently did nothing when
// the user asked for autostart is exactly the quiet mismatch ticket 296 ruled
// out for the hotkey chain.
//
// A config that cannot be read is not an error here: it is the same fallback
// every other boot leg takes (the compiled defaults, whose [app].autostart is
// false), and the fallback is named in the returned receipt. Boot never stops
// for this - a Run key nobody can read is not a reason to refuse the ball.
func autostartSync(cfgPath string, valueName string) string {
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		c = config.NewDefaults()
		return applyAutostartValue(valueName, c.App.Autostart) +
			fmt.Sprintf("（config.toml 不可读 %v，按编译默认表 autostart=false 处理）", err)
	}
	return applyAutostartValue(valueName, c.App.Autostart)
}

// applyAutostartValue is the registry half: make one Run value agree with want.
func applyAutostartValue(valueName string, want bool) string {
	k, err := registry.OpenKey(registry.CURRENT_USER, autostartRunKeyPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		// The Run key itself can be absent on a bare profile; that absence is
		// exactly what the OFF state needs, so only the "wanted ON" branch
		// treats it as a failure.
		if want {
			return fmt.Sprintf("autostart=on 但注册表不可写（%v）：开机自启未生效，需要重启后重试", err)
		}
		return "autostart=off：Run 键不存在（本就无自启项），无动作"
	}
	defer k.Close()

	if !want {
		if _, _, err := k.GetStringValue(valueName); err != nil {
			return "autostart=off：Run 键无本程序条目，无动作"
		}
		if err := k.DeleteValue(valueName); err != nil {
			return fmt.Sprintf("autostart=off 但删除 Run 条目失败（%v）：开机自启仍在，需要重启后重试", err)
		}
		return "autostart=off：已删除 Run 键条目（开机自启已关闭）"
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("autostart=on 但可执行文件路径解析失败（%v）：开机自启未写入", err)
	}
	// Quoted, because the install root may contain spaces ("projects plans" does,
	// in the dev checkout): an unquoted path with spaces runs the wrong binary.
	val := `"` + exe + `"`
	if prev, _, err := k.GetStringValue(valueName); err == nil && prev == val {
		return "autostart=on：Run 键条目已存在且一致，无动作"
	}
	if err := k.SetStringValue(valueName, val); err != nil {
		return fmt.Sprintf("autostart=on 但写入 Run 键失败（%v）：开机自启未生效，需要重启后重试", err)
	}
	return fmt.Sprintf("autostart=on：已写入 Run 键条目（%s -> %s）", valueName, filepath.Base(exe))
}
