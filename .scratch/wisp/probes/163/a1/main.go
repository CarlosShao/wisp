// 票 163 a1（只读设计核台件，产码零字节改动）：三发现量。
//  1. 真注册表（照 cmd/wisp/run.go 的组合形状）里 shell.exec / shell.session 在不在；
//  2. 两发现实的独立 spawn：先 cd 再 pwd —— 量「第二条命令不知道第一条」；
//  3. 退出码是不是 Wait 状态取回的字段，并把「从输出文本里猜」那一味做成读数。
// logdir 由旗标传入，调用方（gotest.sh）取脚本自身目录，不继承 CWD。
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

func main() {
	logdir := flag.String("logdir", "", "REQUIRED: probe script directory (never inherited CWD)")
	flag.Parse()
	if *logdir == "" {
		fmt.Fprintln(os.Stderr, "logdir required")
		os.Exit(2)
	}
	if err := os.MkdirAll(*logdir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var b strings.Builder
	say := func(f string, a ...any) {
		line := fmt.Sprintf(f, a...)
		b.WriteString(line + "\n")
		fmt.Println(line)
	}

	say("# 163/a1 go-probe 现量（go run，真包真进程，无镜像/假进程）")

	// ---- 发 1：真注册表，形状照 cmd/wisp/run.go:333-366 ----
	wd, _ := os.Getwd()
	paths := tools.NewPathCanonicalizer([]string{wd}, nil)
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: tools.NewTaskRoster()}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	var names []string
	for _, e := range reg.List() {
		names = append(names, fmt.Sprintf("%s(%v)", e.Tool.Name(), e.Decl.Declared))
	}
	say("## 注册表实名册: %s", strings.Join(names, " "))
	for _, q := range []string{"shell.exec", "shell.session", "task.output"} {
		_, ok := reg.Lookup(q)
		say("## Lookup(%q) => registered=%v", q, ok)
	}
	say("## IsSensitiveSource(shell.exec)=%v IsSensitiveSource(shell.session)=%v IsSensitiveSource(fs.read)=%v (mark 侧名闸, bridge.go:552 用它)",
		risk.IsSensitiveSource("shell.exec"), risk.IsSensitiveSource("shell.session"), risk.IsSensitiveSource("fs.read"))

	// ---- 发 2：两发现实独立 spawn，先 cd 再 pwd ----
	c1 := exec.Command("cmd.exe", "/c", "cd /d C:\\Windows & cd")
	o1, _ := c1.Output()
	say("## 发1(进程A) cmd.exe /c 'cd /d C:\\Windows & cd' => %q rc=%v", strings.TrimSpace(string(o1)), c1.ProcessState.ExitCode())
	c2 := exec.Command("cmd.exe", "/c", "cd")
	o2, _ := c2.Output()
	say("## 发2(进程B,全新) cmd.exe /c 'cd' => %q rc=%v", strings.TrimSpace(string(o2)), c2.ProcessState.ExitCode())
	say("## 读数: 发2 %s发1的目录 ⇒ 两次生发之间没有任何存活载体（cwd 每发只能靠 exec.Cmd.Dir 单发指定）",
		map[bool]string{true: "≠", false: "=="}[strings.TrimSpace(string(o2)) != `C:\Windows`])

	// ---- 发 3：退出码从哪来 + 摘掉取回机制那一味 ----
	c3 := exec.Command("cmd.exe", "/c", "exit 7")
	err3 := c3.Run()
	var ee *exec.ExitError
	if errors.As(err3, &ee) {
		say("## exit 7: errors.As(*exec.ExitError) 命中, ExitCode()=%d —— Wait 状态【取回的字段】", ee.ExitCode())
	} else {
		say("## exit 7: 未命中 ExitError (err=%v)", err3)
	}
	c4 := exec.Command("cmd.exe", "/c", "echo exit 1")
	o4, err4 := c4.Output()
	say("## 真值=0 的进程打印了 %q：若『从输出文本里猜』会把 %d 读成 1；字段取回 rc=%v err=%v",
		strings.TrimSpace(string(o4)), 1, c4.ProcessState.ExitCode(), err4 == nil)
	say("# 台件结束")

	out := filepath.Join(*logdir, "go-probe.log")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("# log=" + out)
}
