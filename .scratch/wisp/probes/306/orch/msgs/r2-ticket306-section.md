---

## 编排者收 `306-a1`（2026-10-10 22:1x，三笔 `4f210227`→`5940f64a`→`408c0d06`，件 `probes/306/a1/**`；我自己的复验件＝`probes/306/orch/logs/r1-orch-verify-306a1-20261010-221240.txt`，脚本＝`probes/306/orch/r1-verify-306a1.sh`）⇒ **`AC#0` 由我翻**＋★**第三枚行内权威值 `:196` 入账成新格 `AC#2b`**＋★**"要不要升级成接线票"我裁＝⛔ 不升级（理由带盘上出处，见第 2 节）**

### 0. 台面（⛔ 与腿同台面＝同一枚 HEAD，⛔ 在共享工作树里做任何 checkout／改产码）

- 起手尺＝`git status --porcelain -- internal cmd docs scripts .github` 现量 **0 行**；锚＝`HEAD_sha=408c0d06`（＝腿的终笔之后同一枚台面）。
- `git show HEAD:internal/audio/wavinjector.go | cmp - internal/audio/wavinjector.go` ＝ **`rc_blob_vs_worktree_cmp=0`** ⇒ 本节引用的每一枚行号在工作树与 HEAD blob 上逐字相同（两把尺并排，是读数不是假设）。
- 终局再量一把：`final_porcelain_lines=0`＋`rc_final_blob_cmp=0` ⇒ **本程产码面一字未动**（那枚"改成引用常量"的决定性突变只落在仓外导出树 `C:/Users/swq/tmp/306o-refconst-*`）。
- 本程⛔ 动任何产码；腿的写面我逐笔核过＝三笔并集只含 `probes/306/a1/**`（新 `.go` **0** 枚、`.out` **0** 枚、票面框改动 **0**）。

### 1. 三张表里我**自己现跑**的读数（⛔ 抄腿的转述；每条带尺）

| 表 | 尺 | 我这发的读数 |
|---|---|---|
| ②覆盖尺（洞的**原因**） | `GOFLAGS=-mod=mod go test ./internal/audio/ -count=1 -coverprofile=…` | `rc_test=0`；末行逐字 `ok  	github.com/CarlosShao/wisp/internal/audio	15.967s	coverage: 59.3% of statements`；函数级逐字 `parseWav 70.0%`／`NewWavInjector 100.0%`／`parseWaveFormat 81.8%`／`Drain 0.0%`／`convertPacket 0.0%`／`total 59.3%` ⇒ **六数与腿逐字对上**（腿那发 15.960s，我 15.967s＝时长⛔ 分母） |
| ②块级（那两支零执行者） | 同一件块图，射程＝`wavinjector.go:188-226` | **零命中六块逐字**：`193.5,193.18 1 0`／`194.6,195.1 1 0`／`196.5,196.65 1 0`／`219.3,221.26 3 0`／`222.4,223.1 1 0`／`224.3,224.43 1 0`；正控三块有命中＝`188.4,192.24 5 1`（只走到 `if` 头）／`212.3,214.26 3 1`／`226.3,226.111 1 1` |
| ②调用链 | `grep -rn --include=*.go "NewWavInjector(" internal cmd tools` | 全命中 **13** 行／剥定义＝**12** 枚调用点／**非 `_test.go` 且非自身文件＝0 枚**；剩下的命中逐字只有两行注释＝`internal/audio/doc.go:16`＋`internal/audio/gate.go:29` ⇒ 腿那句"响亮报回"我复现 |
| ②根因（夹具编码器） | `sed -n '26p;27p;32p' wavinjector_test.go` | 逐字 `buf.Write(u32le(16))`／`buf.Write(u16le(1)) // PCM`／`buf.Write(u16le(16))` ⇒ 全仓唯一 wav 夹具把标签与位深写死成 **1/16** |
| ②旁路入口 | `sed -n '581p' hotplug_test.go` | 逐字 `parsed, rate, chans, err := parseWav(raw)` ⇒ **同包测试直接可达 `parseWav` 有现成先例**（这一枚是第 2 节裁定唯一需要的可达性凭据） |
| ③硬前提（⛔ 改成引用常量） | 仓外导出树 `git archive HEAD \| tar -x`，`sed` 把 `:192`→`waveFormatExt`、`:218`→`waveFormatFloat` | `rc_linux_mutated=`**1**，逐字两句 `internal\audio\wavinjector.go:192:17: undefined: waveFormatExt` 与 `:218:17: undefined: waveFormatFloat`；同树 `rc_windows_mutated=`**0**；真仓 pristine `GOOS=linux go build ./internal/audio/` `rc=`**0** ⇒ **"平台中立件引不到那三枚常量"＝实测，⛔ 推断**；票面 §边界第 3 条据此成立 |
| ③甲形成本 | `git ls-files \| grep -icE '\.wav$'`／`find internal cmd -type d -name testdata`／`grep -nE "wav" .gitignore` | 被跟踪 `.wav`＝**0**；`testdata` 只有 `internal/agent`、`internal/llm`、`internal/models`、`cmd/wisp` 四枚（audio **没有**）；`.gitignore` 零命中（`rc=1`）⇒ 甲会同时开两处先例，但⛔ 会被静默跳过 |
| ③CI 可见性 | `sed -n '242p;253p' scripts/portable-tests.sh`＋`sed -n '164p;168p;193p;195p'`＋`sed -n '419,420p;487p;533,534p;797p' .github/workflows/ci.yml` | core 档 scope 逐字含 `./internal/audio/...`（`:242`）且 `test-core: runs-on: ubuntu-latest`（`:419-420`，调用点 `:487`）；windows 档逐字含 `./internal/audio/`（`:253`）且 `runs-on: windows-latest`（`:533-534`，调用点 `:797`）；`core_pin` 起 `:164`、audio 在 `:168`，`win_pin` 起 `:193`、audio 在 `:195`，两处逐字都是**包名** ⇒ 新增测试文件⛔ 动 pin（见第 5 节：可见面比票面写的**更宽**） |

### 2. ★裁"要不要升级成接线票"＝**⛔ 不升级**（本票保持仪器票；票面 `AC#0` ② 那句"若根本没有生产调用者就升级由我另裁"，这一裁在本节，⛔ 改原句）

腿现量到 `wavinjector.go` 整枚文件⛔ 生产调用者，并按票面把这一枚交给我裁。我的判据是**盘上契约文本＋已有归口**，⛔ 感觉：

1. **"测试替身没有生产调用者"是 C8 的本意，⛔ 是缺陷**。`docs/PLAN.md:1358` 那行逐字＝`| C8 | `AudioSource` | 真实麦克风 / **wav 注入（测试钩子）** / `[Aec]` 仅接口位 | D16, 测试地基 |`；`PLAN.md:1478` 逐字＝`| **语音链路** | **C8 `AudioSource` 注入 wav 文件**，不用真麦克风 | 这是整个策略的地基 —— 没有它语音链路完全不可自动化测试 |`；`docs/specs/SPEC-04-voice-pipeline.md:124` 逐字＝`- **全部走 C8 WavInjector**，不依赖真麦克风（整个测试策略的地基）`；`internal/audio/doc.go:13-16` 自称 `WavInjector (the test backbone; ALL pipeline tests drive data through it, SPEC-04 §9)`。⇒ 给 injector 找一枚生产调用者＝**把测试钩子接成生产路径**＝撞 C8／D16，要走人工批准，⛔ 落在本票。
2. **接线那半件早就另有家、且已具名**（⛔ 本票该认领）：`docs/evidence/s1/241-audio-level-producer-r1.md:109` 逐字写着生产那一环"具名欠两枚票"＝**票 07 的 GUI 腿**＋**票 228 之后的消费腿**；`docs/reports/pending-and-issues.md:11120`（票 247 `247-a2`）已把采集协程的落点裁成**甲＝挂 `cmd/wisp` 常驻腿**；`docs/evidence/s1/154-host-id-never-closed-r1.md:342` **更早**就量到同一件事（`NewWavInjector` 三者在生产码里"只出现在自身定义与包文档"、"整条麦克风腿今天不接任何宿主"）。⇒ 腿这一枚是**复认既有现量**，⛔ 新洞、⛔ 新票由头。
3. **牙⛔ 依赖生产接线**：`parseWav` 是包内私有函数、同包测试直接可达（现成先例＝上面 `hotplug_test.go:581`），所以 `AC#1`／`AC#2`／新格 `AC#2b` 三枚牙**纯测试侧就能买到**。今天真正"没牙"的是 **R2**（`Drain 0.0%`／`convertPacket 0.0%`，`wasapi_windows.go:378` 那句谓词），票面已划出射程、要动需 owner 另一次批准 ⇒ ⛔ 因它给本票加射程。

⇒ **`AC#0` 我翻**（三张表齐、尺名齐、我逐枚复跑过），本票性质不变＝仪器票。⚠ 具名残余：`306-r1` 的夹具**只能**给"解析出来的值"装牙；"injector 到底该不该被生产调用"这一问自本程起**永久有家**＝票 07／票 228 之后的消费腿／票 247 落点甲，谁引本票说这事＝引第 2 节。

### 3. ★第三枚行内权威值 `:196` 入账（新格 `AC#2b`，落点腿的活，⛔ 我这一程写进票面就完）

- 内容锚逐字＝`fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])`（`internal/audio/wavinjector.go:196`；blob＝工作树已证）。
- 它⛔ 是那一族"值"（`3`／`0xFFFE`），而是**权威偏移 24**＝`WAVEFORMATEXTENSIBLE.SubFormat` 的起点；同一件事在 `wasapi_windows.go:188-193` 的注释里被钉过，而 injector 这一份⛔ 牙——我这发块图逐字 `196.5,196.65 1 0`。
- ⇒ 新格 `AC#2b`＝`body+24`→`body+26` 必须让指名用例红（成对两发：改⇒红／`cmp` 逐字节还原⇒绿）。
- ⚠ **射程写死**：一枚 extensible＋float32 夹具**同时**给 `:192`／`:196`／`:218` 三处装牙（腿的表③ 推演与我的块读数同向），但"一次夹具"⛔ 能当"三格都闭"的证据——三枚各要自己那一发突变红。

### 4. 与腿的读数对不上的**三枚**（具名；⛔ 改腿的件，账记在我这边的派单与复算上）

1. **`writeWav` 枚数**：腿写"调用点 **11** 枚"，它那把尺（`grep -rn "writeWav(" internal | wc -l`）**含定义行**`wavinjector_test.go:18`。我同尺并排＝含定义 **11** 行、剥定义 **10** 枚调用点（按文件 `wavinjector_test.go` 4／`captureopt_test.go` 3／`gate_test.go` 1／`hotplug_test.go` 1／`level_test.go` 1），而它的按文件分解把 `wavinjector_test.go` 写成 5＝**多算了定义那一枚**。结论⛔ 受影响（10 枚全走硬编码 tag=1/bits=16）。
2. **注释跨行**：腿引"doc comment 写着 PCM16…"为 `wavinjector.go:163-164`，盘上"PCM16 (format 1) and IEEE float32 (format 3, incl. WAVE_FORMAT_EXTENSIBLE…"起在 **`:162`** ⇒ 该写 `:162-164`。⚠ 这一枚⛔ 影响任何判语（那条注释是"注释比代码有牙"的例证，语义一样）。
3. **戊类排除枚数**：腿写"排除 **6 枚**"，它自己列出的名册只有 **5** 枚位（`internal/tools/recycle_windows.go:37`／`cmd/wisp/panel_host_windows_test.go:181`／`cmd/wisp/notify_windows.go:27`、`:30`、`:35`）。我这把独立尺（`grep -rnE '0x0003|0x00000001|0xFFFE' internal cmd tools` 剥 `internal/audio/`）现量＝**逐字 5 行**，与名册同。⚠ 它那句"否则会被虚报成 8 枚"用的是**另一把尺**（只 `0x0003` 形状：我复跑＝8 行，其中 6 行在 `wave_format_float_300_windows_test.go` 内＝在族、⛔ 该排除）⇒ 两个数都⛔ 错，错在**一把尺的分母写成了一句**。★定式入派单：**"排除 N 枚"必须＝名册行数，列完自己 `wc -l` 一遍**。

### 5. 腿加出的一枚**更宽射程**（票面 `:25` 只写了 windows 档；原句⛔ 改，以本节为准）

票面"CI 可见性前提"那一枚只写了 `windows` 档（`portable-tests.sh:253`）。腿现量＋我复跑：**core 档也含 audio**（`:242` 逐字 `./internal/audio/...`，`test-core` 跑 `ubuntu-latest`）。⇒ **新增一枚无 `//go:build windows` 的测试文件，CI 上是两档可见（ubuntu＋windows）**，派单里"哪一档会跑到"必写两档。
⚠ 反向那一层我标**〔读文件结构推的，本程⛔ 现量〕**：票 300 那两枚钉带 `//go:build windows` 文件名约束 ⇒ 它们⛔ 进 core 档，其 CI 面仍是 windows 单档；要坐实得取一发 core 档日志（`test-core` 那 job 里 audio 的 `--- PASS` 名册）。⛔ 本程取（要 push 后那一发）。

### 6. 现态与 next

- **现态**＝票 306 **1 勾／5 未勾**（`AC#0` 勾；`AC#1`／`AC#2`／**新 `AC#2b`**／`AC#3`／`AC#4` 未勾），⛔ 改 `-done`。
- next＝派 **`306-r1`**（**乙形**＝在既有 injector 用例里加 40 字节 extensible＋float32 面，⛔ 新二进制夹具、⛔ 加标签；一次装三枚牙＋**三对**突变凭据；`AC#3` 注释格同批；`AC#4` 门禁按票面含"整包两发取交集只到'稳定新增红 0'"那半句新规矩）→ 非实现者 **`306-v1`**；队列 **`305-a1`**（派单必带票 305 现量那两句逐字＋`A826` 第 3 节的射程写法）→ `296-r2` → `295-r1` → `294-r1` → `294-v1` → `298-v1`。
- **记我两笔（⛔ 算成腿的欠账）**：① 我给 `306-a1` 的派单把锚号写成 `d9864baf`（那是 `303-v1` 的终笔），当时真 HEAD＝`f994f956` ⇒ 腿自锚现量并具名报回、没被带歪；定式＝派单锚号写之前现跑 `git rev-parse --short HEAD`。② 票 303 `AC#3` 那枚不可满足判据的**正式登记**（票 141 先例）与票 300 `AC#4` 的 `cmd/wisp` 各 ≥2 发那半，两样在 `A826`／`A825` 都还挂着、本程⛔ 办。
- Progress log：
  - [2026-10-10 22:1x +0800] agent=编排者 did=收 `306-a1`（三笔 `4f210227`→`408c0d06`，写面并集＝只 `probes/306/a1/**`，新 `.go` 0／`.out` 0／票框改动 0，越界 0）＋**我自己同台面复跑**（`porcelain` 起手与终局各 0 行；blob↔工作树 `cmp` rc=0）：覆盖尺 `rc_test=0`、`ok … internal/audio 15.967s coverage: 59.3%`、函数级 `parseWav 70.0%`／`NewWavInjector 100.0%`／`parseWaveFormat 81.8%`／`Drain 0.0%`／`convertPacket 0.0%` 六数逐字对上腿；块级零命中六块逐字复现（`193.5,193.18`／`194.6,195.1`／`196.5,196.65`／`219.3,221.26`／`222.4,223.1`／`224.3,224.43` 全 `1 0` 或 `3 0`）；`NewWavInjector` 剥定义 12 枚调用点、非 `_test.go` 非自身文件＝**0**（只剩 `doc.go:16`＋`gate.go:29` 两行注释，逐字复现）；夹具硬编码逐字 `u32le(16)`／`u16le(1) // PCM`／`u16le(16)`；旁路先例 `hotplug_test.go:581` 逐字；★决定性两发我自己在仓外导出树复跑＝把那两枚行内值换成 `waveFormatExt`/`waveFormatFloat` 后 `GOOS=linux go build ./internal/audio/` **rc=1**（逐字两句 `undefined:`）、同树 `GOOS=windows` **rc=0**、真仓 pristine linux **rc=0**、终局产码面 `cmp` 未动；成本尺复现（被跟踪 `.wav`=0、`testdata` 四枚⛔ 含 audio、`.gitignore` wav 零命中 rc=1）；CI 尺复现（core `:242` 与 windows `:253` 两档都含 audio、`core_pin:168`／`win_pin:195` 逐字都是**包名**、`test-core=ubuntu-latest`／`test-windows=windows-latest`）。⇒ **`AC#0` 翻**。★**裁"要不要升级成接线票"＝⛔ 不升级**：`PLAN.md:1358`／`:1478`＋`SPEC-04:124`＋`doc.go:13-16` 逐字定 C8 的 wav 注入＝**测试钩子/地基本意**，给它加生产调用者＝撞契约要走人工批准；接线那半件早有归口（`241-audio-level-producer-r1.md:109`＝票 07 的 GUI 腿＋票 228 之后的消费腿、`pending-and-issues.md:11120`＝票 247 落点甲、`154-host-id-never-closed-r1.md:342` 更早已量到同形）⇒ 腿这一枚是**复认**⛔ 新洞；牙⛔ 依赖生产接线（`parseWav` 同包可达，先例逐字在手）。★**第三枚 `:196`（权威偏移 `body+24`，hit=0）入账成新格 `AC#2b`**，⚠ 写死"一枚夹具同装三枚牙⛔ 等于三格都闭，三枚各要自己那发突变红"。⚠**与腿对不上的三枚具名**：`writeWav` 腿报 11＝含定义行、剥定义真 **10**（它按文件把 `wavinjector_test.go` 写成 5，真 4）；doc comment 腿引 `:163-164`、盘上起 `:162`；戊类腿写"排除 6 枚"而名册**5** 枚位、我这把独立尺逐字也 5 行（它那句"虚报成 8"用的是只 `0x0003` 的另一把尺＝两数都对、分母⛔ 同），★定式＝"排除 N 枚"必须＝名册行数，列完自己 `wc -l`。★腿加出一枚**更宽射程**＝票面只写 windows 档、实则 core(ubuntu) 也含 audio ⇒ 无标签新文件两档可见（反向那层"带标签的钉⛔ 进 core"我标〔读结构推的〕，坐实要一发 core 日志）。记我两笔＝派单锚号写错成 `d9864baf`（真 HEAD `f994f956`，腿自锚具名报回）＋票 303 `AC#3` 不可满足判据登记与票 300 `AC#4` 的 `cmd/wisp` 那半仍挂着。现态 **1 勾／5 未勾**，⛔ `-done`；next=派 `306-r1`（**乙形**：40 字节 extensible＋float32 面进既有 injector 用例，⛔ 新二进制夹具／⛔ 加标签／⛔ 改成引用常量／⛔ 动 C5-R2／⛔ 删 `waveFormatPCM`）→ `306-v1` → `305-a1` → `296-r2` → `295-r1` → `294-r1` → `294-v1` → `298-v1`；本程零产码、⛔ push 由本笔执行（`A813` fast-forward）
