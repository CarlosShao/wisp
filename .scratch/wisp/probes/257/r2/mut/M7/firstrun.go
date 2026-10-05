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
			"里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。"+
			// 票 261 Ⓐ（A571）：引导教的条目形状必须能过 llm 门的形状——手写条目缺
			// enabled 键＝视为关闭（defaults 不进 map，解码落 false），点名一条没写
			// enabled = true 的模型会在起动时被具名拒绝。只加这一句指路，不重写引导。
			"注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。\n")
	// ---------------------------------------------------------------------
	// 票 257 形 ⓒ（账 A543）：干净机器上设置页那七枚字段一枚都写不进去，因为
	// 首份默认配置里服务商注册表是 nil，而写侧只改已有的行。裁定的那一形不是
	// 替用户建行，是把"你自己该怎么加"在这份给人看的回执里说全。三段各答一票：
	//
	//   1. AC#1 的终态＝ ⓒ 真兑现：手加三样（provider 行／model 行／roles.chat
	//      点名），静态链是 257-a1 §ⓒ 量的（只加第一样＝3/7，三样齐＝7/7）。第二样
	//      的落点是"就地填已有那一节"而不是"再追加一节"——首份文件本来就带
	//      [llm.roles.chat]，追加同名节是 TOML 的 duplicate table，文件自己加载不
	//      过（257-r1b 现量，M3 那发红句就是它）。拼法一律 [llm.providers.<名>]：
	//      顶格写 [providers.x] 会撞 internal/config/parse.go 的 DisallowUnknownFields
	//      （现量 :72 与 :123），整条面板链起不来——票面早期那半句被普查腿顶回，
	//      照 A543 的更正办，这里也照抄给机主的那一份就不许再出现那个拼法。
	//   2. AC#2 的三因各一句：措辞与 internal/config 那三枚 tag（settings.go 的
	//      refusalFileMissing／refusalRowMissing／refusalInvalid）逐字对齐，拒写
	//      现场说原话、这里预先把那三种说清，⛔ 不合成一句"配置未生效"（票 223
	//      AC#4 那条定式的同一形状）。
	//   3. §8 边界①＋A560 的更正：什么时候算用上，今天有三种进程形状。旧的两格措辞
	//      （"手改＝热加载认／面板写＝要重启"）在"任务腿没过控制台闸的常驻形状"下是
	//      半谎——那个进程里可能根本没有重读盘的东西，所以三形状各说一句。凭据那段
	//      只写真接了的路：引用是新建端点时解一次（run.go:435-436 → internal/llm/
	//      resolver.go:141），配置层那条解引用通道今天没接线（A560：三处生产
	//      NewManager 第二参数全 nil），⛔ 不把它列成入口、也不回显任何值（AC#3）。
	//
	// 三段都只走 stderr 这一条"给人看的回执"通道：一行都不落进被生成的
	// config.toml（票 198 的钉逐字禁首建文件里出现 [llm.providers 字样，
	// firstrun_198_test.go:215 那一族），也不复用 cause=missing 那句（同文件
	// :237 的钉）。
	fmt.Fprint(stderr,
		"wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，"+
			"模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；"+
			"要在这一页配上模型，得在这份文件里手加三样，缺一不可：\n"+
			"wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，"+
			"protocol 与 base_url 可以留空。）。"+
			"第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。"+
			"第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名"+
			"（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。"+
			"界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，"+
			"只加第一样只解锁服务商那三枚。\n")
	fmt.Fprintf(stderr,
		"wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：\n"+
			"  第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，%s 现在是真的文件；"+
			"首启之前没有任何旧配置可言。\n"+
			"  第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，"+
			"加完才写得进；这一条说的不是你的值不对。\n"+
			"  第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验"+
			"（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，"+
			"文件一个字节都没动。\n",
		cfgPath)
	fmt.Fprint(stderr,
		"wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：\n"+
			"  在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，"+
			"[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，"+
			"热加载不会替它换脑，真正发请求还是按启动时那一份。\n"+
			"  没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，"+
			"手改与面板写在两个方向上都只能等下一次启动。\n"+
			"  设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；"+
			"页面上的读数在重启之前也不会跟着你手改的文件走。\n"+
			"wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，"+
			"把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；"+
			"这一页任何时候都不回显密钥的值，只说已录入还是没录入。\n")
	return true, nil
}
