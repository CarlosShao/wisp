package main

import (
	"fmt"
	"os"
)

func main() {
	base := "/tmp/exp120"
	os.RemoveAll(base)
	os.MkdirAll(base+"/root", 0o700)
	os.MkdirAll(base+"/foreign", 0o700)
	foreign := base + "/foreign/keep-me.txt"
	os.WriteFile(foreign, []byte("x"), 0o600)
	os.Chmod(foreign, 0o666)

	link := base + "/root/leaf"
	os.Symlink(foreign, link)
	os.Symlink(base+"/foreign", base+"/root/dirlink")
	realFile := base + "/root/real.txt"
	os.WriteFile(realFile, []byte("y"), 0o644)
	zero := base + "/root/zero.txt"
	os.WriteFile(zero, []byte("z"), 0o600)
	os.Chmod(zero, 0o000)

	fmt.Println("uid =", os.Getuid())

	root, err := os.OpenRoot("/")
	fmt.Println("OpenRoot(/):", root != nil, err)
	if err != nil {
		return
	}
	defer root.Close()

	f, err := root.OpenFile("tmp/exp120/root/leaf", os.O_RDONLY, 0)
	fmt.Println("A leaf-symlink open err:", err)
	if err == nil {
		st, _ := f.Stat()
		fmt.Println("   A OPENED A LINK, mode:", st.Mode())
		f.Close()
	}

	f2, err2 := root.OpenFile("tmp/exp120/root/dirlink/keep-me.txt", os.O_RDONLY, 0)
	fmt.Println("B mid-symlink open err:", err2)
	if err2 == nil {
		f2.Close()
	}

	f3, err3 := root.OpenFile("tmp/exp120/root/real.txt", os.O_RDONLY, 0)
	fmt.Println("C real-file open err:", err3)
	if err3 == nil {
		fmt.Println("   C fchmod err:", f3.Chmod(0o600))
		st, _ := f3.Stat()
		fmt.Println("   C mode after fchmod:", st.Mode().Perm())
		f3.Close()
	}

	f4, err4 := root.OpenFile("tmp/exp120/root", os.O_RDONLY, 0)
	fmt.Println("D dir open err:", err4)
	if err4 == nil {
		fmt.Println("   D fchmod err:", f4.Chmod(0o700))
		st, _ := f4.Stat()
		fmt.Println("   D dir mode:", st.Mode().Perm())
		f4.Close()
	}

	f5, err5 := root.OpenFile("tmp/exp120/root/zero.txt", os.O_RDONLY, 0)
	fmt.Println("E mode-000 open err (uid", os.Getuid(), "):", err5)
	if err5 == nil {
		fmt.Println("   E fchmod err:", f5.Chmod(0o600))
		f5.Close()
	}

	_, err6 := root.OpenFile("/tmp/exp120/root/real.txt", os.O_RDONLY, 0)
	fmt.Println("F absolute name inside root err:", err6)

	nf, err7 := root.OpenFile("tmp/exp120/root/real.txt", os.O_WRONLY, 0)
	fmt.Println("G O_WRONLY err:", err7)
	if err7 == nil {
		nf.Close()
	}

	st, _ := os.Stat(foreign)
	fmt.Println("H foreign mode after all:", st.Mode().Perm())
}
