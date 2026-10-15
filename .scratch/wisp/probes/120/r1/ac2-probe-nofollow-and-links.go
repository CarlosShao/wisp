package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	base := "/tmp/exp120d"
	os.RemoveAll(base)
	os.MkdirAll(base+"/root/sib", 0o700)
	os.MkdirAll(base+"/foreign", 0o700)
	os.WriteFile(base+"/root/real.txt", []byte("y"), 0o666)
	os.WriteFile(base+"/foreign/keep.txt", []byte("x"), 0o666)
	os.WriteFile(base+"/root/sib/inner.txt", []byte("i"), 0o666)
	os.Symlink(base+"/foreign/keep.txt", base+"/root/absleaf")
	os.Symlink("sib/inner.txt", base+"/root/relleaf")
	os.Symlink("../foreign", base+"/root/updir")
	os.Symlink(base+"/foreign", base+"/root/absdir")
	os.Symlink("sib", base+"/root/reldir")

	root, _ := os.OpenRoot("/")
	defer root.Close()
	rel := "tmp/exp120d/root/"

	try := func(what, name string, flag int) {
		f, err := root.OpenFile(name, flag, 0)
		if err != nil {
			fmt.Printf("%-28s %-40s REFUSED: %v\n", what, name, err)
			return
		}
		st, _ := f.Stat()
		fmt.Printf("%-28s %-40s OPENED mode=%v size=%d\n", what, name, st.Mode().Perm(), st.Size())
		f.Close()
	}

	try("leaf abs-target, plain", rel+"absleaf", os.O_RDONLY)
	try("leaf abs-target, NOFOLLOW", rel+"absleaf", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("leaf rel-target, NOFOLLOW", rel+"relleaf", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("leaf rel-target, plain", rel+"relleaf", os.O_RDONLY)
	try("anc abs-target, NOFOLLOW", rel+"absdir/keep.txt", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("anc rel-target, NOFOLLOW", rel+"reldir/inner.txt", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("anc updir, NOFOLLOW", rel+"updir/keep.txt", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("real file, NOFOLLOW", rel+"real.txt", os.O_RDONLY|syscall.O_NOFOLLOW)
	try("real dir, NOFOLLOW", rel+"sib", os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK)

	// Root.OpenRoot per component
	for _, name := range []string{"sib", "reldir", "absdir", "updir"} {
		sub, err := root.OpenRoot(rel + name)
		if err != nil {
			fmt.Printf("OpenRoot(%-16s) REFUSED: %v\n", rel+name, err)
			continue
		}
		fmt.Printf("OpenRoot(%-16s) fd-anchored OK\n", rel+name)
		sub.Close()
	}
}
