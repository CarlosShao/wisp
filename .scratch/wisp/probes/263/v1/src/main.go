package main

// 263-v1 台件假 wisp：只实现 scripts/slo-check.ps1 用到的那一小组命令行形状
// （slo -state/-settle/-leak -seconds -interval-ms -out），判定线全由 V1_VERDICT 控制。
// ⛔ 它不测任何真实性能：本件只证"门有没有真的在等、有没有真的把码取回来"。
// 用法：probe263v1 slo -state Sleeping -seconds 2 -interval-ms 250 -out F.json
// 环境：V1_VERDICT=pass（默认，档位达标、exit 0）/ fail（真超标形状：pass:false + exit 1）
//       V1_SLEEP_SCALE 存在时忽略 -seconds、固定睡 2s（GUI/CUI 对照件用）

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

func argVal(args []string, name string) string {
	for i, a := range args {
		if a == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func main() {
	args := os.Args[1:]
	verdict := os.Getenv("V1_VERDICT")
	if verdict == "" {
		verdict = "pass"
	}
	isLeak := false
	for _, a := range args {
		if a == "-leak" {
			isLeak = true
		}
	}
	secs := 1.0
	if s := argVal(args, "-seconds"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			secs = f
		}
	}
	if os.Getenv("V1_SLEEP_SCALE") != "" {
		secs = 2.0
	}
	out := argVal(args, "-out")
	name := "probe263v1"
	if len(os.Args) > 0 && os.Args[0] != "" {
		if len(os.Args[0]) > 4 {
			name = os.Args[0]
		}
	}
	fmt.Printf("%s: started pid=%d args=%v\n", "probe263v1", os.Getpid(), args)
	time.Sleep(time.Duration(secs * float64(time.Second)))
	fmt.Printf("%s: exiting\n", "probe263v1")

	pass := verdict == "pass"
	code := 0
	if !pass {
		code = 1
	}
	if isLeak {
		// leak 档：门要求它必须翻红（exit 1），假件按契约给出 pass:false + exit 1
		pass = false
		code = 1
	}
	report := map[string]any{
		"pass":     pass,
		"state":    argVal(args, "-state"),
		"settle":   map[string]any{"free_os_memory_count": 1},
		"seconds":  secs,
		"selfname": name,
	}
	if out != "" {
		b, _ := json.MarshalIndent(report, "", "  ")
		if err := os.WriteFile(out, b, 0644); err != nil {
			fmt.Printf("probe263v1: write %s failed: %v\n", out, err)
		}
	}
	os.Exit(code)
}
