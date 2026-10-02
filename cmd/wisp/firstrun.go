package main

// Ticket 198: the first-run creation of config.toml, hooked to ONE entry.
//
// WHAT WAS MISSING (普查件 198-a1 §1.0, 推翻票面核心事实后成立的那一半): the
// capability was never absent. internal/config exports the default table
// (NewDefaults, defaults.go:58 - the `default:"..."` tags in schema.go are the
// single source, D36 rule 3) and the atomic sealed writer (SaveFile,
// loader.go:238), and writeguard.go:125-131 already proves a missing file is
// creatable through them ("writing it creates what first-run did not"). What
// no production path ever did was CALL that pair once, at the moment a user
// first asks `wisp run` for work. On a fresh machine the run leg therefore
// died at run.go's 配置未就绪 branch with exit 2, quoting a missing file as
// the user's only experience of the product.
//
// WHY THIS FILE AND NOT assembleRuntime (编排者裁定 J1, 账 A514): assembleRuntime
// is the shared assembly root - the resident leg calls it directly
// (resident_task_source_windows.go:265). Two existing pins forbid the resident
// leg from ever producing a config.toml nobody asked it for:
// resident_task_source_246_windows_test.go:389-391 (red sentence "the leg
// created a config.toml it was never asked for") and logsink_windows_test.go
// :159-164 (a run-leg boot that stops at Unconfigured because no config.toml
// exists is "the point of the case rather than an accident"). Creation inside
// assembleRuntime would falsify three assertions there today; loosening those
// pins to land this ticket is on the forbidden list (AGENTS.md §1.1), so
// creation hangs ONLY on `wisp run`, the entry the user explicitly starts, via
// runTextTask - the single production caller chain main.go cmdRun ->
// runTextTask. The J1 shape itself is pinned by
// TestTicket198FirstRunCallerIsTheRunEntryOnly below.
//
// WHAT THIS DOES NOT DO: no new defaults table, no template string, no second
// serializer (198-a1 §1.5 - that would mint the unguarded second truth source
// this repository's priciest failures are named for). No relaxation of the
// exit-2 path: the created file carries NO model and NO api_key_ref, so
// run.go's ResolveRole branch still fails loudly with 2 afterwards (R13).
// "没有配置就当默认跑" stays refused.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/secret"
)

// ensureFirstRunConfig creates <dataDir>\config.toml from the exported
// defaults and returns true, but ONLY when the file is really absent.
//
// The absence test is os.Stat + errors.Is(err, fs.ErrNotExist) - the exported
// equivalent of internal/config's unexported fileMissing (parse.go:232-235):
//   - a DIRECTORY standing under the config name (ticket 101's 权限读不到
//     stand-in, run_mode101_test.go:621-633) stats clean, so it is not missing
//     and this function never touches it. A "stat failed -> create" test would
//     try to rename over it (atomicWrite, parse.go:226) and turn an honest
//     read failure into a write accident.
//   - any OTHER stat failure (permission, reparse, ...) also declines to
//     create: it is not evidence of absence, and the load path classifies it
//     through describeReloadFailure's own branches instead of this one.
//
// The directory side needs no new code: the data root exists as a sealed tree
// only through the same decider production already uses (secret.NewStore ->
// winsec.PrivateDirAll, store.go:49), so this function calls that and never
// hand-decides a path beyond joining the resolved root with the existing
// configFileName constant (the same join run.go:391 already performs inside
// assembleRuntime; the
// d22scan prohibition is on filepath.Clean/Abs decisions outside
// risk.PathResolver, PLAN.md:1288, not on Join over a resolved root).
func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error) {
	cfgPath := filepath.Join(dataDir, configFileName)
	if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {
		// Exists (file or directory), or the stat answered something that is
		// not absence: not ours to create, not ours to report.
		return false, nil
	}
	if _, err := secret.NewStore(dataDir); err != nil {
		return false, fmt.Errorf("数据根不可建（%s）：%w", dataDir, err)
	}
	if err := config.SaveFile(cfgPath, config.NewDefaults()); err != nil {
		// atomicWrite is temp+rename (parse.go:196), so a failure here leaves
		// no half file at cfgPath; the absence stands and the load path below
		// still answers the honest missing-config verdict, exit 2.
		return false, fmt.Errorf("config.toml 首建未写成（%s）：%w", cfgPath, err)
	}
	// 说实话的回执（票面"要建什么"2 的前半：建了没建成 + 绝对路径 + 值从哪来
	// + 本轮仍缺什么）。它刻意 NOT 复用 cause=missing 那一句（config_reload.go
	// :317-320 是热加载"缺文件"分支的台词，说的是"继续用内存里的旧配置"，而首
	// 建场景根本没有旧配置；票 223 的名册钉着四句归因互斥，198-a1 §4-N9）。
	fmt.Fprintf(stderr,
		"wisp run: 已在 %s 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），"+
			"未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码\n",
		cfgPath)
	// AC#4（票面"去哪儿补"，198-r2）：下面两句只写用户今天真走得通的那条路，
	// 逐枚现读得来，不写票面原话那句"再补一行配置"：
	//   - key 的录入者是 cmdSecret 的 set 支（secret.go:333-387 runSet）：明文只
	//     经隐藏输入或 --from-stdin 进 internal/secret 的 DPAPI 存储
	//     （store.go:64-83 Store），命令回显一行 api_key_ref = "dpapi:<blob 名>"
	//     （secret.go:380）。config.toml 里根本没有承载明文的字段
	//     （schema.go:384-387，D36 rule 5：只许 dpapi:/env: 引用），所以手改文件
	//     补的是那个 *名字*、不是 key 本身；env: 那一支的值由系统环境提供
	//     （store.go:95-103 Resolve，没有明文兜底）。
	//   - 模型目录没有任何 CLI 写入者：cmdProviders 的 discover/probe 先 LoadFile
	//     才动作（providers.go:99-103、:118-133 目录里没有该 provider 就退），
	//     wisp models 是 C29 的本地签名模型库（models.go:60-74），不是 LLM 目录。
	//     所以模型这一格今天只有改这份文件一条路，文案就照实说。
	// 两句里只许出现命令名、字段名、占位名，绝不出现任何凭据值（C28；票 63 的
	// argv/日志两把尺对同一个边界已在别处钉着）。
	fmt.Fprint(stderr,
		"wisp run: 缺的两样各有各的入口。key：先跑 wisp secret set <blob 名>（隐藏输入，"+
			"不进 argv 也不进日志，也可以 --from-stdin 从管道喂），它把明文交给 DPAPI 存储，"+
			"并回显一行 api_key_ref = \"dpapi:<blob 名>\"；这份文件里没有写明文 key 的字段，"+
			"要补的是那个名字。不想用 DPAPI 就把 api_key_ref 写成 env:<环境变量名>，值由系统环境提供。\n"+
			"wisp run: 模型：在同一份文件的 [llm.providers.<名>] 里补 api_key_ref、base_url 与 models.<id>"+
			"（名字对上内置预设的，protocol 与 base_url 可以留空），再在 [llm] 的 text_chain 或 roles.chat "+
			"里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。\n")
	return true, nil
}
