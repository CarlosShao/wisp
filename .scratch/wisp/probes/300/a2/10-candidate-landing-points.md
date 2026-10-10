# 300-a2 · 10 — `AC#6` 候选落点名册（⛔ 裁形、⛔ 推荐、⛔ 产码）

起手锚＝`4e00357f`（2026-10-10T17:34:47+08:00）；本件全部读自 **blob 层**（`git show HEAD:<path>` / `git grep … HEAD`），⛔ 拿工作树当 HEAD。名册级尺逐条附在下面 F 段。**本件只量"该钉在哪、撞⛔ 撞、买到什么"，⛔ 替编排者挑形状。**

## F. 五枚硬事实（每枚带尺，后面各候选都靠它们）

- **F1 常量只活在 windows-tagged 文件里。** 尺＝`git grep -n -e 'waveFormatFloat' HEAD -- '*.go'` ⇒ 全仓 **4 枚**命中，逐枚＝`internal/audio/wasapi_windows.go:98`（声明 `waveFormatFloat = 3`）、`:378`（`floating := s.format.tag == waveFormatFloat`，在 `Drain()` 体内）、`internal/audio/parse_wave_format_300_windows_test.go:152`＋`:153`（恒等式那两行）。尺＝`git show HEAD:internal/audio/wasapi_windows.go | sed -n '1,4p'` ⇒ 第 1 行逐字 `//go:build windows`。⇒ **硬后果：任何直接引用该常量的判据⛔ 可能不带 windows tag**，linux 档连编译面都进不去。
- **F2 `waveFormatPCM = 1` 声明后零使用。** 尺＝`git grep -n -e 'waveFormatPCM' -e 'waveFormatExt' HEAD -- '*.go'` ⇒ `waveFormatPCM` 只命中声明行（:97），`waveFormatExt` 声明 :99 ＋唯一使用 :210（`if f.tag == waveFormatExt`）。⇒ 顺带买到⛔ 买到，看落点选谁。
- **F3 生产里那枚常量只有⛔ 一个消费者，且没抽出来。** 尺＝同 F1 的第二枚命中（:378 在 `Drain()` 体内）＋`git grep -n -e 'convertPacket' HEAD -- internal cmd` ⇒ 生产调用点**只有 :400 一枚**，测试侧**0 枚**（另 2 枚命中是 :149/:155 的注释与红句文案）。`convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool)` 收的是 **bool** ⇒ 测试要观察"生产怎么用这枚常量"只有两条路：调 `Drain()`（要活流＝真设备，⛔ 本票票面 :31）或在用例里复写同形比较。**这枚事实直接框住候选形状空间。**
- **F4 同包第二枚实现用的是字面量，且那字面量也零尺。** 尺＝`git show HEAD:internal/audio/wavinjector.go`（`case fmtTag == 0xFFFE`、守卫 `if size < 40`、`fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])`、`case fmtTag == 1 && bits == 16`、`case fmtTag == 3 && bits == 32`）；尺＝`git show HEAD:internal/audio/wavinjector_test.go | grep -n -i -e float -e FFFE -e extensib -e tag` ⇒ **除 `sineI16At(... float64 ...)` 那枚生成器 0 命中** ⇒ injector 的 float32/extensible 那两支今天**零覆盖**。
- **F5 枚数口径：两把尺⛔ 同形。** 尺 A＝对 blob 名册逐枚 `git show HEAD:<file> | grep -c '^func Test'` ⇒ 8 枚 `_test.go` 分别 1/4/5/8/10/1/7/6 ＝ **42 枚顶层 `func Test` 行**（其中 windows-tagged＝`capturelevel_windows_test.go` 1 ＋ `hotplug_test.go` 8 ＋ `parse_wave_format_300_windows_test.go` 1 ＝ **10 枚**；tag 判据尺＝逐文件 blob 前 3 行，⚠ 按**文件名**分类会漏 `hotplug_test.go`——它文件名不带 `_windows` 而内有 tag）。尺 B＝`git grep -n -e 'func TestMain' HEAD -- internal/audio` ⇒ **0 命中**（故 42⛔ 由 `TestMain` 解释）。票面/派单引的是"整包 **41** 枚顶层全绿"（`--- PASS` 那把尺）。⇒ **41↔42 是尺差⛔ 是矛盾**，落笔引用时必带尺名（票面 `AC#4` 自己钉过"三数必带尺名"）。

## 候选（逐枚：动哪枚文件／判据形状／换形 3→1 与 3→0 红⛔ 红／买到什么／买不到什么）

### C1 既有 `internal/audio/parse_wave_format_300_windows_test.go` 里加一段
- **动**：只动那枚 tracked 用例文件（表格新增一列或新增一段断言），⛔ 动产码。
- **判据形状**：把浮点那一支的**期望侧**从"随常量挪的来源"换成⛔ 随常量挪的来源（表格里的字面期望列／由夹具 GUID 推出的期望），观察侧仍经该常量。
- **换形**：**3→1 ⇒ 红**（float 面的观察侧判成非浮点，同时 PCM 面判成浮点，两面各红一发）；**3→0 ⇒ 红**（float 面）。依据＝F1/F3：常量只动**一**侧。恒等式那形（:152-153 两侧同枚常量）被这枚形状直接堵死。
- **买到**：改动面最小；⛔ 新增分母枚数（顶层用例名不变，多的是子测试 ⇒ 若要顶层枚数作差，先说清是哪把尺，见 F5）；⛔ 碰 scripts/CI。
- **买不到**：⛔ 钉"该常量⛔ 是 3"这件事本身（它钉的是"浮点判定⇔解析值"的一致性）；⛔ 顺带钉 `waveFormatExt`/`waveFormatPCM`；⛔ 与 `AC#1`/`AC#2` 的凭据隔离——**同一枚文件被两格凭据共用**，改它 ⇒ 那对成对基线（票面 `300-r1` 交的 `432/320`→`437/321` 名册作差）与 `300-v1`/`300-v3` 引用的 `:152-153` 逐字凭据一起漂（docs 面 prose 锚，⛔ 仪器，见 20 件 N1/N5）。
- **撞**：⛔ 放宽既有断言⛔ 成立才算不撞 ⇒ 见 20 件 N1（那枚文件的 `got.tag != tc.wantTag` 是 `AC#1`/`AC#2` 凭据本体，硬约束②射程）。

### C2 新枚同包 windows-tagged 文件（例：`internal/audio/wave_format_float_300_windows_test.go`）
- **动**：新增一枚 `//go:build windows` 的 `_test.go`，⛔ 动产码、⛔ 动既有件。
- **判据形状**：由测试自带的**权威侧**期望（把 `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 的 `Data1` 当从权威文本抄来的字面，或走 `convertPacket` 的**输出字节面**让浮点⇔整型在样本值上不可混淆）× 观察侧经 `parseWaveFormat`＋该常量。
- **换形**：**两形都红**的前提＝期望侧的推导链⛔ 经过该常量（依据 F1＋F3）；**颜色本腿取⛔ 到**（⛔ 编译面授权），那一发成对读数（改⇒红／`cmp` 还原⇒绿）归落地腿。
- **买到**：与 `AC#1`/`AC#2` **物理隔离**（硬约束②零风险）；可同时把 `waveFormatExt`、`waveFormatPCM` 一起钉上（F2 那两枚今天零尺）；windows 档分母 **+1 枚顶层用例**（票 301 `AC#1` 已把 audio 拉进 `windows)` 档，见 30 件）。
- **买不到**：⛔ 外部权威（值到底⛔ 是 3 仍靠测试作者抄的那份表）；⛔ ubuntu 档可见（F1）⇒ CI 上只有**一台 OS** 会求值它。
- **撞**：⛔ 与既有件同名（名册＝F5 那 8 枚）；⛔ 改既有断言；⚠ 新增顶层用例会把该包推进 `test-windows` 那一档的红名册作差面（落地腿要按票面 `AC#4` 改名册尺）。

### C3 钉在权威表那族（读盘上 SDK 头文件 `shared/mmreg.h`，与常量对拉）
- **动**：新枚 windows-tagged 用例（读盘）；⛔ 动产码。
- **判据形状**：开文件→抽 `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 的 `Data1`→与 `waveFormatFloat` 比。**两侧来源⛔ 同源** ⇒ 这是名册里**唯一把"值⇔外部权威"钉住**的形（正是硬约束③要的"匹配语义而⛔ 匹配形状"）。
- **换形**：**3→1 红、3→0 红**，依据＝权威文本⛔ 随 Go 常量变（该形⛔ 依赖任何行为复写）。
- **买到**：语义级钉子；⛔ 设备；⛔ 产码改动。
- **买不到／代价（三枚具名，都⛔ 我裁）**：
  1. **与本仓一句既有声明冲突**：tracked 用例头顶注释逐字写着 the absolute path is machine- and build-version-specific, so it is deliberately not written out here（尺＝`git show HEAD:internal/audio/parse_wave_format_300_windows_test.go | sed -n '17,21p'`）。⇒ 落这枚＝**要改那句话**，属⛔ 由腿自取的形式分歧。
  2. **读盘依赖＝本仓已知红形先例**：`internal/ball/tokens_table_test.go` 那族硬读 `design/assets/tokens.css`（尺＝`git grep -n -e 'tokens.css' HEAD -- '*.go'` ⇒ 命中 `internal/ball/tokens.go` 5 枚＋`tokens_table_test.go` 6 枚），台账口径是"工作树缺它就红⛔ 当归因当程"。
  3. **托管 runner 上有没有那枚头文件＝本腿取⛔ 到的读数**（尺＝`git grep -n -i -e 'Windows Kits' -e 'mmreg.h' HEAD -- scripts tools .github` ⇒ **0 命中** ⇒ 仓内从无 SDK 路径的读取或准备动作）。⚠ 连带一枚**结构性代价**：若该用例在无头文件时走 `t.Skip` ⇒ 未登记 SKIP＝红档（夹具 ledger 的 class 定义逐字 `fixture = the machine or OS cannot supply the subject at all`，尺＝`git show HEAD:scripts/portable-tests.sh | sed -n '586,600p'`），而登记要改 `scripts/portable-tests.sh` ⇒ **撞票面 `AC#4` 的名册（只许 `internal/audio/**`＋`probes/300/**`）**。⇒ 这枚候选要落地**必须同时扩射程**，那一格⛔ 由我拍。

### C4 跨实现一致性钉（同一枚 40 字节面喂两枚实现，要求两侧对"浮点"判定一致）
- **动**：新枚用例（可 windows-tag，因为 `parseWaveFormat` 只在 windows 面存在）。
- **判据形状**：injector 侧是**字面量 3**（F4）⇒ 常量挪了它⛔ 挪 ⇒ 观察侧单侧移动 ⇒ 换形必红。
- **换形**：3→1 红／3→0 红（依据＝F4 那枚字面量与常量⛔ 同源）。
- **买到**：第二枚**仓内**独立来源；⛔ 盘上依赖、⛔ 真设备；**顺带买到 injector float/extensible 分支的第一枚覆盖**（F4 现零覆盖）——⚠ 那是**新增覆盖面**、⛔ 属 `AC#6` 的要求，落地时要具名说⛔ 顺手扩射程。
- **买不到**：⛔ 外部权威——两枚仓内实现可以一起错。⚠ **一处要具名报回的形状问题**：票面 `AC#0` 那条反循环硬门原文是"⛔ 拿本仓另一枚实现当裁判"，它当时钉的是**偏移**；这枚候选钉的是**值**。两者算⛔ 同一条禁令＝**交编排者裁，我⛔ 拍**。
- **撞**：`parseWav` 的 `size < 40` 守卫会被拉进判据射程（⛔ 改产码）；⛔ 动 `wavinjector_test.go` 那 6 枚既有断言（见 20 件 N2）。

### C5 生产侧先把 :378 那一行抽成包内谓词，再钉谓词
- **动**：`internal/audio/wasapi_windows.go`（**产码面**，虽在 `AC#4` 允许的 `internal/audio/**` 射程内，但票面只授权过一次产码改动＝`AC#2` 的偏移＋注释那一发）。
- **判据形状**：期望侧走字面／GUID，观察侧调抽出来的谓词。
- **换形**：两形都红（同 C1 依据）。
- **买到**：`Drain()` 那条**唯一消费者**第一次变得可测（F3 说今天只有那一行用常量，而它在 COM 泵里）；这是名册里唯一能覆盖"生产真的这么判"的形。
- **买不到**：⛔ 值⇔权威；⚠ 代价＝**改产码语义面**，且 `Drain()` 那行的形状一旦外提，`300-v2` 判过的"注释逐句可证伪性"凭据也要重跑。
- **撞**：⛔ 碰 `gate.go`/`device.go` 的行数（那两枚被 cmd/wisp 侧按行号引，见 20 件 N5）；⛔ 动 `cmd/wisp`（`303-r1` 在飞＋票面 :32＋`AC#4` 名册）。

## 各候选共有的两枚约束（⛔ 我裁，只是量到的）

1. **文案面**：新判据的红句要"具名指向那枚常量"（`AC#6` 原文），而那句话必然落在**字符串字面量**里；本仓 `d22scan` 的 ban #8 **含 `_test.go`**（尺＝`git grep -n -e '_test.go' HEAD -- tools/d22scan/main.go` ⇒ :44 逐字 `_test.go` INCLUDED (D23). This is the one ban whose…），且**注释豁免、字符串⛔ 豁免** ⇒ 红句里⛔ 能用 `⛔`(U+26D4)、`✓`(U+2713)、`≤`(U+2264) 这几枚（仪器射程＝U+1F000–1FAFF／U+2200–22FF／U+2600–27BF／U+2B00–2BFF／U+FE0F／U+1F1E6–1F1FF）。
2. **CI 面**：F1 决定了⛔ 任何 ubuntu 档的尺都⛔ 可能因这枚常量换形而红 ⇒ "换形必红"与"CI 双 OS 可见"二者⛔ 能同时买到（详见 30 件）。

rc=0
