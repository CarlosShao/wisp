package main

import (
	"fmt"
	"errors"
	"syscall"
	"golang.org/x/sys/windows"
)

func main() {
	fmt.Printf("go=%s\n", "1.27.1")
	fmt.Printf("syscall.Errno(0).Error() = %q\n", syscall.Errno(0).Error())
	fmt.Printf("windows.Errno(0).Error() = %q\n", windows.Errno(0).Error())
	fmt.Printf("bytes=%x\n", []byte(syscall.Errno(0).Error()))
	fmt.Printf("syscallErrno0_Is_errorsIs = %v\n", errors.Is(syscall.Errno(0), syscall.Errno(0)))
	fmt.Printf("wrapped_Is = %v\n", errors.Is(fmt.Errorf("ctx: %w", syscall.Errno(0)), syscall.Errno(0)))
	fmt.Printf("wrapped_Error = %q\n", fmt.Errorf("ctx: %w", syscall.Errno(0)).Error())
	fmt.Printf("errno5_Error = %q\n", syscall.Errno(5).Error())
	fmt.Printf("errno5_Is5 = %v\n", errors.Is(syscall.Errno(5), syscall.Errno(5)))
	fmt.Printf("errno0_Is5 = %v\n", errors.Is(syscall.Errno(0), syscall.Errno(5)))
	fmt.Printf("errno87_Error = %q\n", syscall.Errno(87).Error())
	fmt.Printf("errno1400_Error = %q\n", syscall.Errno(1400).Error())
}
