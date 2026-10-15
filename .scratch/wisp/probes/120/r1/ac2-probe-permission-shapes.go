package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	base := "/tmp/exp120b"
	os.RemoveAll(base)
	os.MkdirAll(base+"/root", 0o700)
	os.MkdirAll(base+"/foreign", 0o700)
	os.WriteFile(base+"/root/real.txt", []byte("y"), 0o644)
	os.Chmod(base+"/root/real.txt", 0o666)

	// FIFO
	syscall.Mkfifo(base+"/root/fifo", 0o600)

	root, err := os.OpenRoot("/")
	if err != nil {
		fmt.Println("OpenRoot:", err)
		return
	}
	defer root.Close()

	// 1. O_NONBLOCK on a FIFO through Root
	done := make(chan string, 1)
	go func() {
		f, err := root.OpenFile("tmp/exp120b/root/fifo", os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			done <- fmt.Sprintf("fifo open err: %v", err)
			return
		}
		mode := fmt.Sprint(f.Chmod(0o600))
		st, _ := f.Stat()
		done <- fmt.Sprintf("fifo open ok, fchmod=%s mode=%v", mode, st.Mode())
		f.Close()
	}()
	select {
	case r := <-done:
		fmt.Println("1", r)
	case <-time.After(3 * time.Second):
		fmt.Println("1 fifo open BLOCKED (no O_NONBLOCK effect)")
	}

	// 2. rel "." (path == "/")
	f2, err2 := root.OpenFile(".", os.O_RDONLY, 0)
	fmt.Println("2 open(\".\") err:", err2)
	if err2 == nil {
		st, _ := f2.Stat()
		fmt.Println("   mode:", st.Mode().Perm())
		f2.Close()
	}

	// 3. trailing slash
	f3, err3 := root.OpenFile("tmp/exp120b/root/", os.O_RDONLY, 0)
	fmt.Println("3 open(dir/) err:", err3)
	if err3 == nil {
		f3.Close()
	}
	f4, err4 := root.OpenFile("tmp/exp120b/root/real.txt/", os.O_RDONLY, 0)
	fmt.Println("4 open(file/) err:", err4)
	if err4 == nil {
		f4.Close()
	}

	// 5. empty name
	f5, err5 := root.OpenFile("", os.O_RDONLY, 0)
	fmt.Println("5 open(\"\") err:", err5)
	if err5 == nil {
		f5.Close()
	}

	// 6. non-root permission shapes: write-only file, 0000 file
	uid := os.Getuid()
	for _, tc := range []struct {
		name string
		mode os.FileMode
	}{
		{"wonly.txt", 0o200},
		{"zero.txt", 0o000},
		{"rdonly.txt", 0o400},
	} {
		p := base + "/root/" + tc.name
		os.WriteFile(p, []byte("x"), 0o600)
		os.Chmod(p, tc.mode)
		f, err := root.OpenFile("tmp/exp120b/root/"+tc.name, os.O_RDONLY, 0)
		if err != nil {
			fmt.Printf("6 uid=%d %s(mode %o) O_RDONLY err: %v\n", uid, tc.name, tc.mode, err)
			continue
		}
		if err := f.Chmod(0o600); err != nil {
			fmt.Printf("6 uid=%d %s(mode %o) open ok, fchmod err: %v\n", uid, tc.name, tc.mode, err)
		} else {
			fmt.Printf("6 uid=%d %s(mode %o) open ok, fchmod ok\n", uid, tc.name, tc.mode)
		}
		f.Close()
	}

	// 7. 0300 directory (no read)
	d := base + "/root/nordir"
	os.Mkdir(d, 0o700)
	os.Chmod(d, 0o300)
	f7, err7 := root.OpenFile("tmp/exp120b/root/nordir", os.O_RDONLY, 0)
	fmt.Printf("7 uid=%d dir 0300 O_RDONLY err: %v\n", uid, err7)
	if err7 == nil {
		fmt.Println("   fchmod err:", f7.Chmod(0o700))
		f7.Close()
	}
}
