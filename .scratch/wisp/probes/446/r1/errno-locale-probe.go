package main

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func msg(langID uint32) string {
	b := make([]uint16, 300)
	flags := uint32(windows.FORMAT_MESSAGE_FROM_SYSTEM | windows.FORMAT_MESSAGE_ARGUMENT_ARRAY | windows.FORMAT_MESSAGE_IGNORE_INSERTS)
	n, err := windows.FormatMessage(flags, 0, 0, langID, b, nil)
	if err != nil {
		return "<FormatMessage FAILED: " + err.Error() + ">"
	}
	for ; n > 0 && (b[n-1] == '\n' || b[n-1] == '\r'); n-- {
	}
	return windows.UTF16ToString(b[:n])
}

func main() {
	k32 := windows.NewLazySystemDLL("kernel32.dll")
	uiLang := k32.NewProc("GetUserDefaultUILanguage")
	sysLang := k32.NewProc("GetSystemDefaultLangID")
	a, _, _ := uiLang.Call()
	b, _, _ := sysLang.Call()
	fmt.Printf("GetUserDefaultUILanguage   = 0x%04X\n", uint16(a))
	fmt.Printf("GetSystemDefaultLangID     = 0x%04X\n", uint16(b))
	fmt.Printf("FormatMessage(errno=0, en-US 0x0409)      = %q\n", msg(0x0409))
	fmt.Printf("FormatMessage(errno=0, zh-CN 0x0804)      = %q\n", msg(0x0804))
	fmt.Printf("FormatMessage(errno=0, langid 0 = default)= %q\n", msg(0))
	fmt.Printf("syscall.Errno(0).Error()                  = %q\n", syscall.Errno(0).Error())
	fmt.Printf("syscall.Errno(5).Error()                  = %q\n", syscall.Errno(5).Error())
}
