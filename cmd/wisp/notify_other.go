//go:build !windows

package main

import "fmt"

// postSystemNotification has no shell to talk to outside Windows (C21: Wisp is
// a Windows product; the non-Windows build exists for CI type-checking only).
// It fails LOUDLY rather than pretending the user saw something, because every
// caller treats a notification error as a run failure.
func postSystemNotification(title, body string) error {
	return fmt.Errorf("notify: 当前平台没有系统通知实现（Wisp 的通知是 Windows shell 能力）：%s", title)
}
