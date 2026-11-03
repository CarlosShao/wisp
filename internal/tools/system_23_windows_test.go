//go:build windows

package tools

// Ticket 23's platform half, at the tool layer: the composition root in
// cmd/wisp/run.go writes &tools.SystemBridge{} with a nil Backend, which
// system_bridge.go documents as "selected by the platform build". This file pins
// that the zero value is USABLE rather than a nil-dereference trap, and that a
// tool assembled with no deps at all still resolves its action and reaches the
// bridge.
//
// Scope note, deliberately narrow: the claims here are about the ASSEMBLY (does
// the zero value route to a backend, and does a call survive it), not about
// reading this machine's desktop - so the one call made is one whose failure is
// also a pass: the assertion is "reached the backend and came back with either a
// report or a model-actable refusal", never "the volume is 40". A runner with no
// audio device, no window or no clipboard must not turn this leg red. The one
// outcome that IS red is ErrNoSystemBackend, because on this platform it would
// mean the zero value stopped reaching a backend.
//
// What is NOT here: the Win32 entry surface itself. That is
// system_backend_windows.go's own business, judged through the per-surface fakes
// in system_ops_windows.go; this ticket's tool layer runs every leg against the
// recording fake in system_23_test.go so no leg can touch the real machine.

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// Test23WindowsZeroBridgeReachesABackend is the positive control the
// platform-neutral leg cannot make: nil Backend is not "no backend", it is "the
// platform's backend", so a call must NOT come back as ErrNoSystemBackend.
func Test23WindowsZeroBridgeReachesABackend(t *testing.T) {
	var bridge SystemBridge // exactly what cmd/wisp/run.go writes
	entries := BuiltinSystemEntries(SystemDeps{Bridge: &bridge})
	if len(entries) != 11 {
		t.Fatalf("零值桥应出 11 枚，得到 %d", len(entries))
	}
	var got Entry
	for _, e := range entries {
		if e.Tool.Name() == "system.get" {
			got = e
			break
		}
	}
	if got.Tool == nil {
		t.Fatalf("system.get 不在名册里")
	}
	res, err := got.Tool.Execute(context.Background(), json.RawMessage(`{}`), nil)
	if err != nil {
		t.Fatalf("Windows 上零值桥的 system.get 返回宿主错误 %v：零值应当可用", err)
	}
	if res.IsError && strings.Contains(res.Text, ErrNoSystemBackend.Error()) {
		t.Fatalf("Windows 上零值桥报「没有后端」（%q）：nil Backend 应当是平台默认后端", res.Text)
	}
	t.Logf("零值桥读数：IsError=%v Text=%s", res.IsError, res.Text)
}
