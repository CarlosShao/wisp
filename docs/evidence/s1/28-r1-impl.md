# 票 28 落地腿 `28-r1` —— SessionScope 落地件（写腿）

> **尺**：`HEAD=ba9ca8fa7`（`dev`），现跑 2026-10-26 18:31–19:05 +08。★**读数逐条现跑**，⛔ 转述。
> **写面**：`internal/session/**` —— **新 4 枚**（`scope.go` 752 行／`scope_test.go` 773 行／`scope_cascade_test.go` 196 行／`scope_engines_windows_test.go` 375 行）＋本文件（**95 行**，未过 170 帽 ⇒ 不拆 `-partN`）。
> **零字节**：`cmd/wisp/**`（`git status --porcelain -- cmd/wisp` ＝ **0 行**）、`internal/speech/**`（0 行，见 §7 的归属）、`internal/audio/**`（0 行）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\28r1\{bin,bin15e,speechbin,cases,cases15e,logs,mut,overlay}`。★**变异走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**。

## 0. 一句话

`internal/session` 从「只有 doc.go 的承诺」变成**真的 keep-alive 所有者**：Warm 里引擎热着、麦克风先关、唤醒**零装载**；90s＝Warm→Settling→dispose（3s 静态期仍热）、30s Conversation 落 Warm、断电器 30s/15s；Conversation 首开走 L2 卡片并落审计；dispose 按 D43#31 四步定序；100 轮后 refs/goroutines 归零。★**真实引擎读数**：冷启第一请求 **3.409s** → Warm **171ms**（**19.9x**，四发 17.7–33.2x），**Warm 装载计数 = 0**。★**Path C 那条按票面 blocked by 59 ⇒ 登记残余、⛔ 硬做**。

## 1. ★★★ 票面 6 格 AC 逐条（原文 ＋ 判定 ＋ 读数）

| # | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| 1 | `Cold vs Warm wake: second-round first-token latency measured; Warm wake shows zero engine loads (instrumentation counter) — B4 evidence.` | **可翻** | 真引擎：冷 `acquire=3.2586296s`／首文本 `150.0ms`／**第一请求 3.4086s**；Warm-1 `acquire=0ns`／`171.35ms`／**171.35ms** ⇒ **19.9x**；`Stats().LoadsAfterWarm=0`（portable 5 轮 ＋ 真引擎 3 轮）；Loader 全程只被调 **2 次** |
| 2 | `Timers: 90s Warm → Settling → dispose; 30s Conversation → Warm with mic closed FIRST (order asserted); battery-saver values applied.` | **可翻** | 89s 仍 Warm／90s＝**Settling**（此刻 `EnginesHot=2`、`closes=0`）／92.999s 仍 Settling／93s＝Disposed 且两枚引擎各 `closes=1`、`Pending()=0`；Conversation 30.000s 落 Warm，`mic.close` 下标 **＜** `SCOPE-WARM` 日志下标，且 MicClose 里读到的 `Stage()=Awake`；断电器 14.999s 仍 Awake／15s Warm／29.999s 仍 Warm／30s Settling |
| 3 | `Path C full duplex (D47, blocked by 59): during Conversation playback mic stays open (AEC), barge-in works (≤400ms), Conversation memory measured and budget row backfilled (or degradation mode verified if over budget).` | **⛔ 残余（阻滞未清）** | ★`59-p15-aec-spike.md` 现读 `Status: ready-for-agent`、**无 `-done` 后缀**、7 格 AC **0 勾** ⇒ **本条的前置没清**。按派单：**登记残余、⛔ 硬做**。本腿不碰 AEC、不碰播放期开麦、不出 400ms 数、不填预算行 |
| 4 | `Conversation first-enable → privacy confirm + log; subsequent enables no prompt.` | **可翻** | 两次 enable ⇒ `prompts=1`、`CONVERSATION-PRIVACY-CONFIRM granted` **恰 1 行**；拒绝 ⇒ `(false,nil)`、仍 OFF、`首次开启被拒` 落日志；`PrivacyConfirm=nil` ⇒ **具名错误**（fail-closed，⛔ 无卡片就开麦） |
| 5 | `Dispose cascade order test (models → panel-hide → FreeOSMemory → grants expiry) and RSS return ≤10s (sampler).` | **可翻** | 序断言：`engine.close:*` **＜** `panel.hide` **＜** `free-os-memory` **＜** `grants.expire`（显式结束与 3s 过期两条路都过）；采样器 `base=474016 during=50806312 after=479112`、**samples=0（8.58ms）** ⇒ 远在 10s 内（48MiB 载荷） |
| 6 | `Refcount leak test: 100 wake/sleep cycles → zero residual refs/goroutines.` | **可翻** | 100 轮（真 `time.Timer`）：`loader calls=200`、每类 `closes=100`、每轮 `Refs=0 EnginesHot=0 LoadsAfterWarm=0 WarmResumes=1 Disposals=1`；`goroutines base=2 after=2` |

★**AC3 之外的票面 Key constraints 落点**：Session=wake→显式结束或 90s idle（✓）；Warm＝模型热＋麦关＋零装载（✓）；Conversation 默认 OFF（✓ `Scope` 初值 false，只能 `EnableConversation()` 开）；refcount 只走 SessionScope（✓ 唯一持有 `Engine` 的地方，`Loader` 注入，别处拿不到句柄）；Warm 窗口是 29/10 的异步入口（✓ `Hooks.WarmWindow`，每进一次 Warm 调一次且已在 Warm 之后）；Settling 里新唤醒取消（✓ D43#33 零装载）；显式结束立即 dispose（✓）；**KWS 断电器常驻 OFF** ⇒ 见 §9 残余。

## 2. ★★★ 「session scope warm」核心行为：warming 后第一个请求真的快了多少

| 轮 | acquire（装载/借用） | 首文本 | **第一请求 = acquire＋首文本** | 终稿 |
|---|---|---|---|---|
| **冷启**（新 SessionScope，两枚引擎都没热） | **3.2586296s**（3258629600ns） | 150.0ms | **3.4086319s** | `嗯嗯嗯嗯嗯嗯哦哦` |
| **Warm-1** | **0ns** | 171.35ms | **171.3518ms** | 同上 |
| Warm-2／3 | 0ns／0ns | 189.39ms／161.19ms | 189.39ms／161.19ms | 同上 |
| **倍数** | **冷/暖 acquire ＝ 无穷**（0ns 是量测粒度下的真 0） | —— | **19.9x**（四发区间 **17.7x–33.2x**） | —— |

★引擎是真的：`speech.NewDetector`（silero）＋ `speech.NewStreamingAsr`（流式 Paraformer，encoder 165MB）；`dispose` 的级联关掉真会话 `cascade=67.6ms`。★**暖的不是"算得快"，是"不用装"**：首文本那段（171ms vs 150ms）两轮同量级，差的全在 acquire。

## 3. ★★ 反向判据：把 warm 摘掉 ⇒ 红/变慢（三读数 ＋ 还原核，时机在最后一发变异之后）

变异均在仓外 `mut\` ＋ `-overlay`（`go test -overlay`），工作树零突变窗口。

| 发 | 变异 | 面 | rc | 红的断言（逐字节选） |
|---|---|---|---|---|
| **A** | `EnterWarm` 里把引擎全关并清槽（＝**根本没有 warm 池**，"回合结束就释放"） | portable | **1** | **5 红/15**：`warm stats = {… EnginesHot:0 …}, want EnginesHot=2`；`Settling unloaded engines: EnginesHot = 0`；`EnginesHot after the fallback = 0, want 1` |
| A | 同上 | 真引擎 | **1** | `warm stats = {… EnginesHot:0 …}, want 2 hot engines and 0 loads after warm` |
| **B** | `Resume` 里把热槽丢掉（Warm 时热、**唤醒时重新装**） | portable | **1** | **2 红/15**：`LoadsAfterWarm = 10 after 5 warm wakes, want 0 (B4 …)` |
| B | 同上 | 真引擎 | **1** | `warm round 0 acquire = 2.7351729s (2735172900ns), want at most a tenth of the cold load 4.3708357s (4370835700ns): the wake paid for a load` ⇒ **变慢被具名** |
| **还原** | 无 overlay，重编重跑 | portable | **0** | **42 PASS／0 FAIL／0 SKIP** |
| 还原 | 同上 | 真引擎 | **0** | 冷 `3.2324579s`／Warm `121.9134ms` ⇒ **26.5x**；`loads after warm=0`；正控 `2.3510857s` |
| **原样** | 同一发再跑一次（无 overlay） | portable | **0** | **42 PASS／0 FAIL／0 SKIP** |
| 原样 | 同上 | 真引擎 | **0** | 冷 `3.2586296s`／Warm `171.3518ms` ⇒ **19.9x**；`loads after warm=0` |

★★**还原核（时机＝最后一发变异之后）**：①`internal/session/scope.go` sha256 **T1＝T2＝`9EDAAAC0C8DE12A5C7E889A603C7BE2E688A7A459ECDEE0F332718B835FE6E0E`**；②树内 `MUTANT` 标记 **0 命中**（`internal/session/*.go`）；③两枚变异件对工作树 `git diff --no-index --stat` ＝ **各 8 insertions、0 deletions**（差的只有变异块本身，不是另一个版本）；④`git status` 只列我这 4 枚新文件。⇒ **变异只存在于 `D:\tmp\…\mut\`，`-overlay` 只换编译输入**。

## 4. ★★ 正控：冷启动路径仍然对（⛔ warm 了就把冷路径弄坏）

- **发内正控**：冷那一轮**真的转写出文本** `嗯嗯嗯嗯嗯嗯哦哦`（8 runes）——不是"没报错"，是模型输出。
- **发后正控**：Warm 那三轮跑完、`Dispose` 之后**再开一个新 SessionScope** ⇒ `acquire=2.4495831s`（**又付一次冷装载**）＋终稿仍 `嗯嗯嗯嗯嗯嗯哦哦`；新 scope `ColdLoads=2 LoadsAfterWarm=0`。⇒ ①warm 没把冷路径弄坏；②"热"是**每会话**的，不是全局泄漏的缓存（否则跨会话的池＝泄漏）。
- **反面对照**：变异 A 下同一句断言红（`EnginesHot=0`）；变异 B 下新会话仍冷（说明 B 只破坏唤醒路径）。

## 5. ★★ `cmd/wisp` 零回归：15e staging **6/6** ＋ 延迟三数 `262/640/902` **逐位同值**

| 尺 | 命令 | rc | 读数 |
|---|---|---|---|
| 15e staging（全 6 案） | §6 的三步，`-test.run=Test15e` | **0** | **6 PASS／0 FAIL** |
| **延迟三数** | 同上 | —— | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`、`seamWaitToFirstText=5.51s finalLatency=4.096s speechMS=3392 closedBy=end` ⇒ **262/640/902 逐位同值** |
| 关键计数 | 同上 | —— | `drained=181/181 partials=5 finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]`；`handoff_exit=0` |
| 定向（带 `PATH`） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS … 0.01s` |
| 归属 | `git status --porcelain -- cmd/wisp` | —— | **0 行** ⇒ 本腿对 `cmd/wisp` **零字节**，15e 三处落点语义（`handOffModel` 接线／C8 帧消费／唤醒暂停）**一行未动** |

## 6. ★★ staging 跑法（完整命令）＋ 为什么 `PATH` 不够

```
# 真引擎那条（internal/session 的带 tag 文件）
cd D:\work\workspace\projects plans\Wisp
go test -c -tags "wisp_voice_engines" -o D:\tmp\wisp\20261026\28r1\bin\session28.test.exe ./internal/session
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\28r1\bin
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
D:\tmp\wisp\20261026\28r1\bin\session28.test.exe -test.run=TestTicket28ColdVsWarm -test.v   =>  PASS, rc=0

# 15e（§5）
go test -c -o D:\tmp\wisp\20261026\28r1\bin15e\wisp.test.exe ./cmd/wisp
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\28r1\bin15e
set WISP_15E_TMP=D:\tmp\wisp\20261026\28r1\cases15e
D:\tmp\...\bin15e\wisp.test.exe -test.run=Test15e -test.v                              =>  6/6, rc=0
```

★**为什么 `PATH` 不够（15e §8 已测、本腿按同一 staging 复现）**：Windows 对**依赖 DLL** 的搜索序是「**主模块目录** → System32 → … → `PATH`」。`sherpa-onnx-c-api.dll` 能从 `PATH` 找到，但它依赖的 `onnxruntime.dll` 会被 **System32 的 1.17.1** 先接走，而 linked C API 要 **1.28**（`deps.toml` 钉 1.28.2）⇒ 三层口径：**什么都不给 ⇒ `0xc0000135`（进程起不来）／只给 `PATH` ⇒ 起得来但拿 1.17，`internal/speech` 具名拒绝（不崩）／staging（二进制与三枚 DLL 同目录）⇒ 才全绿**。★**本腿的带 tag 文件正是为此存在**：`internal/session` 在 portable core scope 上（ubuntu），不带 tag 时此文件**根本不编译**，`go test ./internal/session` 的 portable 分母（42 案）不需要 DLL 与模型。

## 7. ★ 四数 ＋ 三读数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `cmd/wisp` **定向＋`PATH`** | §5 | **0** | PASS；`ok … 0.095s` |
| `internal/speech` **改前**（18:37 时的工作树） | staging 全套（§6 同法） | **0** | **31 PASS／0 FAIL／0 SKIP** |
| `internal/speech` **改后**（19:0x） | 同上 | **0** | **49 PASS／0 FAIL／6 SKIP**；★**原 31 名逐名对比 missing=0**（新增 18 名、6 个 SKIP 全是 `TestPunctuat*`／`TestFinalTranscriptCrossesTheSeamPunctuated`） |
| `internal/audio`（零字节，独立跑） | `go test -count=1 ./internal/audio` | **0** | `ok … 16.446s`；`git status` **0 行** |
| `internal/session` **改前／改后** | `go test -count=1 -v ./internal/session` | **0／0** | **27 PASS／0／0** ⇒ **42 PASS／0 FAIL／0 SKIP**（+15 案） |
| `go vet` | `./internal/session`；`-tags wisp_voice_engines ./internal/session` | **0／0** | 空 |

★**归属（⛔ 冒领）**：`internal/speech` 现在有**另一条在飞腿**的文件（`M transcriber.go`、`?? punctuate*.go／rss_*.go／slot*.go／tts*.go`，共 10 项）。本腿对它 **零字节**；"31→49＋6 SKIP"全是那条腿的增量，31 个原名**一个没丢**。

## 8. ★ 格式甲乙两尺 ＋ 门禁 rc

| 尺 | 命令 | 输出 |
|---|---|---|
| 甲 `gofmt -l` | `internal cmd tools scripts frontend` | 恰 **3 行**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`（**已知 3 条既存假红**，CRLF/`autocrlf` 造的，⛔ 本腿责任） |
| 乙 `gofumpt -l`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | 同上 | **同上 3 行、逐字相同** |
| 归属判据 | 两尺对我 4 枚新文件单点跑 | **甲空、乙空** |
| 门禁 `d22scan` | `cd tools/d22scan && go run . -root <repo>` | **rc 0** `clean - no D22 ban violations`；`ban#8 internal/=679`（15-v1 记 666 ⇒ **+13**：我 4 ＋ 另一腿 9）、`bans#1-5 internal/=263`、`cmd/=43`、`ban#8 cmd/=155` |
| ★中途一发 rc 1 | 同上（18:4x 那次） | 唯一 finding＝**另一腿**未跟踪件 `internal/speech/punctuate.go:132` 的 phantom-citation（引 `docs/evidence/s1/27-r1-impl.md` 尚未落盘）；该腿落盘后复跑 **rc 0**。⛔ 我的 4 枚文件任何一发都是 0 finding |
| ⛔ 错尺 | —— | **没用 ban#8 扫 `docs/`**（它只扫 `internal/`＋`cmd/`＋`design/`＋`frontend/`） |

## 9. 未做／边界（⛔ 越界即作废，逐条自陈）

- **AC3 Path C 全双工＝具名残余（阻滞票 59 未 done）**：播放期开麦＋AEC、barge-in ≤400ms、Conversation 内存实测与预算行回填，全部**未做**。★按票面与 59 的 `Blocks:` 行，这条只能等 59 的 go/no-go；⛔ 没有"先写个假的"。
- **生产调用者 0**：`cmd/wisp` **没有**构造 `SessionScope`。原因不是省事——把 warm 池接进 `residentVoice` 就必须改 15-e 的落点语义（那条腿在**回合结束**就释放引擎，`stop()` 一停全关），正是派单第 1 条硬边界 ⇒ **停手上报、⛔ 动**。⇒ `internal/session` 今天仍是"库＋读数齐、装配等下一波"，与票 166 当时的形状同类（166-v1：`生产调用者 0/0/0/0/0` 不算接线）。
- **KWS 断电器常驻 OFF 未接线**：KWS 驻留发生在 `Armed`，那里**没有 session**；本腿只把 30s/15s 两个值钉死并测。KWS 本体是票 41。
- **Conversation 的 UI 半边**（红实心环不消退、Speaking-end 直接回 Listening）不在本包：D43#28 行早已冻结，本腿**一行未改** `internal/statemachine`。
- **Warm 窗口的异步活儿只有口子**：`Hooks.WarmWindow` 每进 Warm 调一次（已测），真正的 L1 抽取（29）／历史压缩（10）是它们的执行器。
- **`doc.go` 的 `DEFERRED(SessionScope)` 标记**现在半过期（lifecycle 已落、装配未落）；摘标记属 `225-deferred-registry-audit` 的写面。
- **配置键 0 新增**：`WarmTimeout` 是 Option（生产根传既有 `[session] warm_timeout_sec`，默认 90）；断电器 30s/15s 是 D42#6 的**包内常量**。
- 包内**零裸 `go func(`**（本包不起 goroutine：唯一异步源是 `Clock` 定时器，回调按 generation 判陈旧）、**生产代码零墙钟算术**（`time.Now` 只在证据测试里量装载耗时）——两句都是 `d22scan` rc 0 ＋ 现跑 grep，不是声明。
- **新依赖边 0**：`internal/session` 生产代码只 import 标准库（`manage.go`/`grants.go` 原有的 `internal/memory` 不变）；`internal/speech` 只在**带 tag 的测试文件**里 import。
