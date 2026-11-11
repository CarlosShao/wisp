package main

// The deferred-disposition 20261110 #5 wiring: privacy.diagnostics_opt_in, the
// total switch that turns ticket 45's diagnostics-bundle engine into a command
// an owner can actually run.
//
// THE SHAPE THE TICKET LEFT. internal/observe has carried the whole engine
// since ticket 45 - PreviewDiagnosticsBundle (review-before-send, entry names
// guaranteed to match the zip), BuildDiagnosticsBundle (fail-closed secret
// detection, default-excluded transcript and watchdog snapshot categories,
// per-item opt-in,
// an oversize skip the manifest records, a manifest that never lies) - and the
// all-repo census found zero production callers. The engine had no方向盘.
// PLAN.md's product line「诊断包一键导出」names the gap, so this is a欠账
// being paid, not a feature being invented.
//
// THE SWITCH IS THE PRODUCT MEANING, NOT A DECORATION. privacy
// .diagnostics_opt_in is a bool总闸 (schema.go:548-549, default false): false
// refuses the whole command - preview included - naming the key, the file and
// the line to flip, because a diagnostics surface that exists while the user
// said "no diagnostics" is not fail-closed, it is fail-rude. true runs
// Preview -> show -> Build. The per-item opt-in (transcripts, watchdog
// snapshots) stays library-side: those categories are default-excluded and the
// builder refuses to read anything the user did not name after a preview; the
// CLI names no items because this build has no file directory it could honestly
// name (transcripts live in the SQLite session store, watchdog snapshots ride
// alerts, neither is a directory on disk), and inventing one would be exactly
// the "a preview that promises an entry the zip will not hold" lie the ticket
// spent a sibling preview function refusing. NOTE (契约面): the usage copy
// below names the config key and the D36 tiers; both already exist in
// PLAN.md/SPEC-03 - this file changes no contract text, and the usage block is
// reconciled against func main's dispatch by censusVsUsage133 in both
// directions.
//
// THE LAYOUT FOLLOWS THE SECRET/CONFIG LEGS line for line (config_cmd_65.go):
// env resolved the one way, data root through resolveSecretLayout, exit codes
// 2 usage/environment, 1 operational refusal, 0 success, subcommands no flags
// to remember except export's two optional ones (-o for the zip path,
// -slo-snapshot to attach a sampler JSON the operator already produced with
// `wisp slo`). The bundle builder redacts; this file adds no redactor of its
// own and never copies a log line anywhere.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// diagnosticsCliUsage180P is the subcommand block printed by `wisp diagnostics`
// and by every usage refusal this leg answers with.
const diagnosticsCliUsage180P = `wisp diagnostics - export a diagnostics bundle (ticket 45)

Usage:
  wisp diagnostics preview
                            print what a bundle WOULD carry: every category,
                            every item, the exact zip entry list, nothing
                            written to disk
  wisp diagnostics export [-o <path.zip>] [-slo-snapshot <file.json>]
                            show the same preview, then write the bundle zip.
                            -o names the output (default: inside the data
                            dir); -slo-snapshot attaches one sampler JSON the
                            operator produced with ` + "`wisp slo`" + ` (optional).

总闸：[privacy] diagnostics_opt_in（config.toml，默认 false）。false 时本命令整体拒绝
（preview 也拒）：总闸关着就没有诊断面。true 才走 预览 -> 导出。
转写与 watchdog 快照两类默认排除、逐项同意（票 45）：本 CLI 不发明文件目录，不指名任何
逐项内容，zip 的 manifest.json 会如实记录「未收集、需逐项同意」。
Environment: the active WISP_ENV is resolved like ` + "`wisp secret`" + ` (SPEC-03 §5).
`

// diagnosticsIO is the process surface of this leg: the same injection seam
// commandsIO/configIO use, so the tests drive the real entry point.
type diagnosticsIO struct {
	stdout io.Writer
	stderr io.Writer
}

// cmdDiagnostics is the `wisp diagnostics` entry point (main dispatch). It
// resolves the env and the data root exactly like `wisp config`, then
// dispatches; the gate (the config key) is consulted inside the two runners so
// the refusal sentences can name the file they read.
func cmdDiagnostics(args []string, dio diagnosticsIO) int {
	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: %v\n", err)
		return 2
	}
	layout, err := resolveSecretLayout(env)
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: %v\n", err)
		return 2
	}
	if len(args) == 0 {
		fmt.Fprint(dio.stderr, diagnosticsCliUsage180P)
		return 2
	}
	switch args[0] {
	case "preview":
		return runDiagnosticsPreview180P(dio, layout.DataDir)
	case "export":
		return runDiagnosticsExport180P(dio, layout.DataDir, args[1:])
	default:
		fmt.Fprintf(dio.stderr, "wisp diagnostics: unknown subcommand %q\n\n", args[0])
		fmt.Fprint(dio.stderr, diagnosticsCliUsage180P)
		return 2
	}
}

// diagnosticsGate180P reads the total switch. It answers the config the gate
// read (nil when the switch is off) and prints the refusal itself, because the
// three refusal shapes - no file, unreadable file, switch off - are three
// different fixes an operator can act on (票 216's discipline: collapsed
// causes are lies by omission).
func diagnosticsGate180P(dio diagnosticsIO, dataDir string) (*config.Config, int) {
	cfgPath := filepath.Join(dataDir, configFileName)
	if _, statErr := os.Stat(cfgPath); statErr != nil {
		fmt.Fprintf(dio.stderr,
			"wisp diagnostics: 拒绝：%s 尚不存在，总闸 [privacy] diagnostics_opt_in 取 schema 默认值 false（先创建配置并把这一行设为 true）\n",
			cfgPath)
		return nil, 1
	}
	c, _, lerr := config.LoadFile(cfgPath, nil)
	if lerr != nil || c == nil {
		fmt.Fprintf(dio.stderr,
			"wisp diagnostics: 拒绝：%s 读不进来（%v）——总闸读不到就是关着，fail-closed；先修好配置文件\n",
			cfgPath, lerr)
		return nil, 1
	}
	if !c.Privacy.DiagnosticsOptIn {
		fmt.Fprintf(dio.stderr,
			"wisp diagnostics: 拒绝：总闸 [privacy] diagnostics_opt_in = false（来源 %s）。"+
				"要导出诊断包，请把 config.toml 里 [privacy] 段的 diagnostics_opt_in 改为 true 后重试；"+
				"总闸关着时本命令不做任何预览、不读任何日志。\n", cfgPath)
		return nil, 1
	}
	return c, 0
}

// diagnosticsBundleOptions180P assembles the one BundleOptions this CLI
// produces: the env's real log dir, the real config.toml, the redact_paths
// value the file itself carries, the caller-supplied extras. The fields left
// zero are left zero ON PURPOSE and the manifest says so - this file never
// fills a slot with something it did not read.
func diagnosticsBundleOptions180P(dataDir string, c *config.Config, outPath, sloSnapshotPath string) (observe.BundleOptions, error) {
	o := observe.BundleOptions{
		OutPath:     outPath,
		LogDir:      filepath.Join(dataDir, logDirName),
		ConfigPath:  filepath.Join(dataDir, configFileName),
		RedactPaths: c.Privacy.RedactPaths,
		Version: map[string]string{
			"wisp_version": buildinfo.Version,
			"wisp_commit":  buildinfo.Commit,
			"wisp_built":   buildinfo.BuildDate,
			"wisp_env":     buildinfo.EnvString(),
			"sherpa_onnx":  sherpa.GetVersion(),
			"collected_by": "wisp diagnostics export (ticket 45 engine, deferred disposition 20261110 #5)",
			"opt_in_gate":  "privacy.diagnostics_opt_in=true",
			"redact_paths": fmt.Sprintf("%t", c.Privacy.RedactPaths),
		},
		// SLOSnapshot: attached only when the operator names a file; otherwise
		// the manifest records "not attached: no sampler data supplied".
		// TranscriptDir/SnapshotDir/OptIn/Deferred/DiskStatus: all zero - no
		// file transcript dir exists in this build, no watchdog snapshot dir
		// exists, no per-item consent was taken, the DEFERRED register is not
		// embedded, and no disk-guard reading is available from a CLI process
		// that never installed one. Every one of these zeros is a manifest
		// line, which is the ticket-45 contract this leg leans on.
	}
	if sloSnapshotPath != "" {
		data, err := os.ReadFile(sloSnapshotPath)
		if err != nil {
			return o, fmt.Errorf("读 -slo-snapshot 失败：%w", err)
		}
		o.SLOSnapshot = data
	}
	return o, nil
}

// printDiagnosticsPreview180P renders a BundlePreview as the operator's rows.
// The entry list is printed verbatim: it is the exact zip entry list the
// builder will produce from the same options, which is the claim a preview is
// for.
func printDiagnosticsPreview180P(w io.Writer, p observe.BundlePreview) {
	fmt.Fprintf(w, "诊断包预览（created_at %s，输出 %s）\n", p.CreatedAt, p.OutPath)
	for _, cat := range p.Categories {
		state := "collect"
		if !cat.Collected {
			state = "skip"
		}
		fmt.Fprintf(w, "  [%s] %s 默认包含=%t：%s\n", cat.Category, state, cat.DefaultIncluded, cat.Detail)
		for _, it := range cat.Items {
			mark := " "
			if it.Included {
				mark = "x"
			}
			line := fmt.Sprintf("    [%s] %-28s %d byte(s)", mark, it.Name, it.Bytes)
			if it.Reason != "" {
				line += "  " + it.Reason
			}
			if it.Head != "" {
				line += "  head=" + it.Head
			}
			fmt.Fprintln(w, line)
		}
	}
	fmt.Fprintf(w, "  zip 条目（%d）：%s\n", len(p.Entries), strings.Join(p.Entries, ", "))
	fmt.Fprintf(w, "  合计 %d byte(s)\n", p.TotalBytes)
}

// runDiagnosticsPreview180P is the preview subcommand: gate, options, preview,
// rows. Nothing is written anywhere - the builder is not called.
func runDiagnosticsPreview180P(dio diagnosticsIO, dataDir string) int {
	c, rc := diagnosticsGate180P(dio, dataDir)
	if c == nil {
		return rc
	}
	o, err := diagnosticsBundleOptions180P(dataDir, c, "(preview only: nothing is written)", "")
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: %v\n", err)
		return 1
	}
	p, err := observe.PreviewDiagnosticsBundle(o)
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: 预览失败（预览失败 = 这个包真的导不出来，先修原因）：%v\n", err)
		return 1
	}
	printDiagnosticsPreview180P(dio.stdout, p)
	return 0
}

// runDiagnosticsExport180P is the export subcommand: same gate, same options,
// then the builder. The preview runs FIRST and is printed, because「先看到、
// 再导出」is the ticket-45 flow the engine was built for - an export whose
// operator never saw the contents list is just a zip nobody vetted.
func runDiagnosticsExport180P(dio diagnosticsIO, dataDir string, args []string) int {
	c, rc := diagnosticsGate180P(dio, dataDir)
	if c == nil {
		return rc
	}
	outPath, sloPath := "", ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			if i+1 >= len(args) {
				fmt.Fprint(dio.stderr, diagnosticsCliUsage180P)
				return 2
			}
			i++
			outPath = args[i]
		case "-slo-snapshot":
			if i+1 >= len(args) {
				fmt.Fprint(dio.stderr, diagnosticsCliUsage180P)
				return 2
			}
			i++
			sloPath = args[i]
		default:
			fmt.Fprintf(dio.stderr, "wisp diagnostics export: 不认识的参数 %q\n\n", args[i])
			fmt.Fprint(dio.stderr, diagnosticsCliUsage180P)
			return 2
		}
	}
	if outPath == "" {
		stamp := observe.NowWallUTC().Format("20060102-150405Z")
		outPath = filepath.Join(dataDir, "wisp-diagnostics-"+stamp+".zip")
	}
	o, err := diagnosticsBundleOptions180P(dataDir, c, outPath, sloPath)
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: %v\n", err)
		return 1
	}
	p, err := observe.PreviewDiagnosticsBundle(o)
	if err != nil {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: 预览失败（导出中止，什么都没写）：%v\n", err)
		return 1
	}
	printDiagnosticsPreview180P(dio.stdout, p)
	written, err := observe.BuildDiagnosticsBundle(o)
	if err != nil {
		// The builder is fail-closed: an aborted bundle leaves no zip behind.
		// The sentence names that, so a red export never reads as "the file
		// exists but is bad".
		fmt.Fprintf(dio.stderr, "wisp diagnostics: 导出失败（zip 已被构建器回收，盘上没有半成品）：%v\n", err)
		return 1
	}
	if written != outPath {
		fmt.Fprintf(dio.stderr, "wisp diagnostics: 构建器回答的路径与请求的不一致（%q != %q），请上报\n", written, outPath)
		return 1
	}
	fmt.Fprintf(dio.stdout, "已导出：%s（内容以上方预览为准；manifest.json 逐条记录了每个条目与每个未收集项的原因）\n", written)
	return 0
}
