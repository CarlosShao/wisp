// 票 163 AC#1b 现量台件（只读；复刻 cmd/wisp/run.go:624-1082 的今日装配形状）：
// 以真 internal/tools 包组成生产注册表，量 shell.exec / shell.session 两枚名字在不在册，
// 并打印整册实名与风险档。与 a1 的差别：a1 只装了 FS+Task 两族（当时生产也只有这两族）；
// 今日 HEAD 的 run.go 还装 System / Subagent / Plan 三族，本台件照今日形状补齐，
// 否则名册会漏枚、读数会被误读成"比 a1 时代更少"。
// logdir 由旗标传入，调用方（ac1b.sh）取脚本自身目录，不继承 CWD。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

func main() {
	logdir := flag.String("logdir", "", "REQUIRED: probe script directory (never inherited CWD)")
	deleteEnabled := flag.Bool("delete-enabled", false, "mirror cfg.FS.DeleteEnabled (default false)")
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

	say("# 163/ac1b go-probe 现量（go run，真包真注册表，形状照 cmd/wisp/run.go:624-1082 今日装配）")

	// ---- 真注册表，装配顺序与 run.go 完全一致 ----
	wd, _ := os.Getwd()
	paths := tools.NewPathCanonicalizer([]string{wd}, nil)
	reg := tools.NewRegistry()
	roster := tools.NewTaskRoster()

	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{
		Paths:         paths,
		DeleteEnabled: *deleteEnabled,
	}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster, Paths: paths}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	for _, e := range tools.BuiltinSystemEntries(tools.SystemDeps{Bridge: &tools.SystemBridge{}}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	for _, e := range tools.BuiltinSubagentEntries(tools.SubagentDeps{
		Roster:      roster,
		ParentTools: nil,
		Provenance:  nil,
		Stream:      nil,
		BaseOptions: nil,
	}) {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}
	for _, e := range tools.BuiltinPlanEntries() {
		if err := reg.Register(e); err != nil {
			say("register %s: %v", e.Tool.Name(), err)
		}
	}

	// ---- 整册实名册（名字＋声明档），按字典序 ----
	var names []string
	for _, e := range reg.List() {
		names = append(names, fmt.Sprintf("%s(%v)", e.Tool.Name(), e.Decl.Declared))
	}
	sort.Strings(names)
	say("## 注册表实名册（%d 枚）: %s", len(names), strings.Join(names, " "))

	// ---- 本票问的两枚 ----
	for _, q := range []string{"shell.exec", "shell.session"} {
		_, ok := reg.Lookup(q)
		say("## Lookup(%q) => registered=%v", q, ok)
	}
	// 对照组：已落名的 task.output（164 号票那枚）仍在册
	_, okTO := reg.Lookup("task.output")
	say("## Lookup(\"task.output\") => registered=%v（对照组）", okTO)

	// ---- mark 侧名闸（bridge 的 C25 标记面；shell 两枚的对照） ----
	say("## IsSensitiveSource(shell.exec)=%v IsSensitiveSource(shell.session)=%v IsSensitiveSource(fs.read)=%v IsSensitiveSource(task.output)=%v",
		risk.IsSensitiveSource("shell.exec"), risk.IsSensitiveSource("shell.session"),
		risk.IsSensitiveSource("fs.read"), risk.IsSensitiveSource("task.output"))

	say("# 台件结束")

	out := filepath.Join(*logdir, "go-probe-ac1b.log")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("# log=" + out)
}
