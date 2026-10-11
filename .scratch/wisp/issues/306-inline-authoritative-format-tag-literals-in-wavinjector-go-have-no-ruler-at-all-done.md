# 票 306 — `wavinjector.go` 里那两枚**行内的权威字面量**＋injector 两支零覆盖：票 300 `AC#6` 那种"常量没人钉"的病，在另一枚文件里原样还在（顺带一格只改注释）

**立票**：2026-10-10 21:1x 编排者（来路＝非实现者验收腿 `300-v4` 的 ⓔ 必答＋**我自己复跑的那四发退码**；它建议"一票两洞"，我认，理由见下面 §为什么要紧 末尾）
**性质**：★**仪器票**（给已经写好的产码装会响的尺），⛔ 缺陷票——今天盘上没有任何证据说 injector 算错了，只有证据说**它错了没人能发现**。
**票号来源**：`306` 是我 2026-10-10 21:1x 现量的下一个空号（`ls .scratch/wisp/issues/ | sed -E 's/^([0-9]+)-.*/\1/' | sort -n | tail -1`＝`305`）。

## 为什么要紧（⛔ 把它读成"补个测试就行"）

票 300 花了一整天才把 `waveFormatFloat` 那枚常量钉住，判据形状是**"改期望侧那枚抄自 `mmreg.h` 的字面量，指名用例必须红"**（`AC#6`）。同一批现量量到：**同源的两枚权威值在 `internal/audio/wavinjector.go` 里是写成行内字面量的，而⛔ 引用那三枚常量**——⇒ 改它们，**全仓没有一把尺会红**。⛔ 让这类事再走一遍"先落地、三个月后发现仪器是瞎的"。

## 现量（每条带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**两枚行内字面量在盘上的形状**（尺＝`git show HEAD:internal/audio/wavinjector.go` 逐字，⛔ 工作树）：
  - `:192` ＝ `			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]`
  - `:218` ＝ `		case fmtTag == 3 && bits == 32:`
  ⚠ 它⛔ 能引用 `waveFormatExt`/`waveFormatFloat`：那两枚常量在 `wasapi_windows.go`，带 `//go:build windows`；`wavinjector.go` 是平台中立件——**这是"为什么它当初被写进行内"的真实原因**，修法⛔ 是"改成引用常量"（那会让非 windows 档编译不过），是**给这枚行内值装一枚钉**。派单必须把这条带上，否则第一腿就会去动 `build` 标签。
- ★**它现在真的零尺（我自己现跑的四发，件＝`probes/300/orch/logs/r3-orch-v4-attacks-rcfixed-20261010-204046.txt`；台面＝`git archive d9864baf | tar -x` 仓外导出树）**：
  - 甲＝`:218` 的 `3`→`2`：**`rc_R1b_jia_pkg=0`**，整包顶层 `--- FAIL` **0**、顶层 `--- PASS` **42**、`ok github.com/CarlosShao/wisp/internal/audio 15.999s`；
  - 乙＝`:192` 的 `0xFFFE`→`0xFFFD`：**`rc_R1b_yi_pkg=0`**，同样 FAIL **0**／PASS **42**、`ok … 15.991s`；
  - 四发各自 `cp`→跑→写回→`cmp`＝**IDENTICAL**，还原后定向正控 `rc=0`、整包正控 `rc=0`，`worktree_internal_audio_dirty_lines=0`。
  ⇒ **把这两枚值换错，包级尺照绿**＝本票唯一需要的"改前必红"凭据已经在盘上（⛔ 落地腿再造一次，它只需复跑并逐字对拉）。
- ★**覆盖尺那一形（R1；腿 `300-v4` 现跑，件 `probes/300/v4/logs/34-cov-run.txt`／`36-cov-injector-tag-branches.txt`／`39-cov-func-named.txt`，⛔ 我复跑，⛔ 引用它的名册只引"我复跑过"的那些）**：`parseWav` 有 **70.0%**，而 injector 里靠这两枚值分支的那两支**逐字零命中**任何名册尺 ⇒ "改行内值不红"⛔ 是巧合，是**那两段代码根本没有执行者**。
- ★**同族第三枚（R2）⛔ 在本票射程**：`wasapi_windows.go:378` 那句 `floating := s.format.tag == waveFormatFloat` 今天 `Drain 0.0%`／`convertPacket 0.0%`（同上覆盖尺）。闭它的**唯一形状**＝把谓词抽成包内可调用（＝票 300 的形 C5＝**动产码语义面**）。⇒ 本票⛔ 碰它，登记在 `A825`，要动得由 owner 一句批准另立票。
- **注释那一处与盘上冲突（本票第二洞的"顺手格"）**：`internal/audio/wave_format_float_300_windows_test.go:38-41` 那句写着 `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` "lives in `ksmedia.h`, not in `mmreg.h`"。现量：`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/mmreg.h:2480-2484` 逐字有 `DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)`，而 `ksmedia.h:851-855` 另有一份**替代定义** ⇒ 两边各有一份、那句"not in mmreg.h"⛔ 成立（其兄弟件 `:101` 引 `mmreg.h:2483` 是对的）。⚠ 这格⛔ 触及任何断言（`AC#6` 钉的是**格式标签** `0x0003`＝`mmreg.h:2110`，⛔ 是 GUID）。
- **CI 可见性前提已量过（省一条派单）**：`scripts/portable-tests.sh:253` 的 `windows` 档逐字含 `./internal/audio/`，`win_pin`（`:195`）是**包名**名册 ⇒ 新增一枚 `internal/audio` 的测试文件⛔ 需要动 pin，而且它在 CI 上真会跑（凭据＝`probes/300/orch/r10-ci-third-run-new-test-visibility.txt`：run `38049322056` 里 `:6109` 顶层 PASS＋`:6110-6115` 六枚子测试）。

## 要建什么（`AC#0` 之前⛔ 任何产码）

- [x] **`AC#0` 只读普查（决定落点，⛔ 预设答案）**：交回三张表——
  ① **行内权威字面量全名册**：仓里还有多少处把 `mmreg.h` 那族值（`3`／`0xFFFE`／`0x0003`／GUID 低字）写成**行内字面量**而⛔ 引常量／引钉过的测试？尺＝`grep -nE '0xFFFE|fmtTag == 3|== 0x0003' internal/**/*.go` ＋逐枚分类"引常量／行内／注释"，名册必带**射程目录＋blob 还是工作树**（本仓老尺）。
  ② **那两支为什么零执行者**：从入口到 `wavinjector.go:192`／`:218` 的调用链逐跳列（谁调 `parseWav`／谁构造 `fmtTag`），并具名答"今天有没有任何夹具喂过 `0xFFFE`＋`bits==32` 这一组合"。⚠ 若"根本没有生产调用者"，本票就升级成**接线票**，由我另裁。
  ③ **夹具落点候选＋代价**：至少两形候选（新增 wav 夹具文件／在既有 injector 用例里加形）各写"要不要动 `third_party`、要不要 windows 标签、CI 哪一档会跑到"。
- [x] **`AC#1` 决定性一发（⛔ 先红再修；本格只产**测试侧**夹具，⛔ 改产码语义）**：让 `0xFFFE`＋`bits==32` 那一支**真的被执行**（extensible 头＋float32 子格式），断言的是**解析出来的那个值**，⛔ 断"被调用过"。★凭据形状⛔ 变：把 `:218` 的 `3` 换成 `2` 之后**必须红**（今天的红＝我 §现量 那两发 `rc=0`＝本票的"改前必红"）。
- [x] **`AC#2` 同族第二枚的牙**：`:192` 的 `0xFFFE`→`0xFFFD` 也必须让指名用例红。⚠ 两枚各一对（改⇒红／`cmp` 逐字节还原⇒绿），缺一枚⛔ 算本格闭合。
- [x] **`AC#2b` 第三枚行内权威值的牙（★`306-a1` 表① 现量新增、立票时⛔ 数到；凭据与裁语在下面的收腿节 3）**：`:196` 的 SubFormat 偏移 `body+24`→`body+26` 必须让**同一枚**指名用例红。⚠ 这枚⛔ 是"格式标签"（⛔ 那族 4 枚常量之一），是 `WAVEFORMATEXTENSIBLE.SubFormat` 的**起点偏移**＝另一个权威值，今天块 `196.5,196.65` hit=0（编排者复跑逐字）。⛔ 与 `AC#1`/`AC#2` 共用一发读数＝**三枚各要自己那一发突变红**，"一枚夹具装了牙"⛔ 等于"三格都闭"；成对两发（改⇒红／`cmp` 逐字节还原⇒绿）缺一枚⛔ 算本格闭合。
- [x] **`AC#3` 只改注释那一格（最小面，⛔ 与 `AC#1`/`AC#2` 同批落）**：把 `wave_format_float_300_windows_test.go:38-41` 那句改成实话（两份定义各在何处、`mmreg.h:2110` 才是那枚钉引用的出处）。**硬边界**＝⛔ 动任何一行断言／任何一枚期望值；⛔ 删 `waveFormatPCM`（票 300 `AC#6` 的附带覆盖面正读着它，删它＝那件编译不过）；⛔ 顺手改别的注释。
- [x] **`AC#4` 门禁＋越界（每把尺自落一行 `rc=N`，0 字节件⛔ 算交）**：`GOFLAGS= go build ./...`／`go vet ./internal/audio/`／`sh scripts/d22scan.sh` 各 rc=0；格式名册**两把并排**（工作树＋仓外 blob，各写射程目录）新增 0 枚；`go test ./internal/audio/ -count=1 -v` 改前改后各 **≥2 发**、同一台面（`git archive` 仓外导出树，⛔ 在共享工作树跑"改前那发"）、红名册 `comm` 双向＝**新增红 0**；三数必带尺名（`^--- PASS` 顶层 vs 含子测试）。写面只落 `internal/audio/**` ＋自家 `probes/306/<腿名>/**`；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；⛔ `frontend/**`／`design/**`／三枚冻结件／`thresholds.go`／`allowlist.txt`／D43 表／`PLAN.md` 零字节；⛔ 零 push、commit 显式 pathspec 写在 `$( … )` 之外。

## 边界（⛔ 塞进本票）

- ⛔ **R2**（`Drain`/`convertPacket` 那两枚 0.0%）＝要动产码语义面（票 300 形 C5）＝**另一次批准**，本票⛔ 碰。
- ⛔ **放宽票 300 `AC#6` 那枚钉**（既有 `--- PASS` 顶层 1＋子测试 6 的形状不许为变绿改）；⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`；⛔ 动 `third_party/sherpa-onnx` 那三枚 `.dll`（冻结）。
- ⛔ 把 `wavinjector.go` 改成引用 `waveFormatExt`/`waveFormatFloat`（`build` 标签那面会挂，见 §为什么要紧）——本票买的是**"行内值也有牙"**，⛔ 是"行内值消失"。
- ⛔ 顺手把 injector 的其它分支补全覆盖（那是另一族缺口，先让 `AC#0` ① 那名册说话）。

## 排程约束（派单里逐字带上）

- 车道＝`internal/audio` 的 Go 编译面**当前无在飞写者**（`300-v4` 只写 `.md`、`303-v1` 只写 `.md`）；但 ⛔ 与 `296-r2`／`294-r1` 那几枚抢 `cmd/wisp` 面——本票⛔ 碰 `cmd/wisp`，真冲突时本票让位（本票的写面包级独立）。
- 起手锚必现跑（`git rev-parse HEAD`＋`git status --porcelain -- internal cmd`），与我派单里给的号不一致时**以原文为准并具名报回**。
- 证据件⛔ 叫 `.out`（根 `.gitignore:8` 是全仓 `*.out` ⇒ 会被静默跳过而 commit 仍回显成功）；退码⛔ 用 `命令 | tee 件; echo $?`（那取的是 tee 的退码；纯重定向 `> 件 2>&1; rc=$?`）。★这条是我 2026-10-10 自己踩过的雷，写死在派单里。

## Progress log

- [2026-10-10 21:1x +0800] agent=编排者 did=立票（料源＝`300-v4` 的 ⓑⓒⓔ 三答＋我自己的 `r3` 四发退码＋`r10` 的 CI 可见性；两洞一票＋一格只改注释）；⛔ 派单（先等 `303-v1` 交完，再按 `AC#0` 起一枚只读腿——它⛔ 跑 Go 编译面以外的东西由派单再裁）
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
  - [2026-10-11 07:4x +0800] agent=**`306-r1`（落地腿，⛔ 参与过本票任何前程）** did=按 §6 那句 next 落**测试侧**四格＋交 `AC#4` 全套读数；
    五笔＝`4e73d0a5`(起手锚)→`1b34e475`(`AC#1` 乙形夹具＝新文件 `internal/audio/wavinjector_extensible_float_306_test.go`，
    40 字节 extensible：`tag=0xFFFE`／`cbSize=22`／`fmt` size=40／`SubFormat.Data1` 低字 `3` 在 body+24／`bits=32`，两枚入口＝同包 `parseWav`（先例 `hotplug_test.go:581`）＋`NewWavInjector(t.TempDir())`，
    断言⛔ "分支被调用过"＝`err`／`rate`／`chans`／`len(samples)=8`（这把尺区分"按 float32 解"与"同 32 字节读成 16 枚 PCM16"）／逐枚 int16 **手推**、⛔ 调 `FloatToPCM16`、⛔ 引产码行内值）→
    `c0d5ec24`(`AC#3` 只改注释：现量 `mmreg.h:2480-2484` 逐字有那枚 `DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)`（在 `:2483`）、`ksmedia.h:850-855` 是同一 GUID 的另一份替代定义（`:854`，由 `defined(_INC_MMREG)` 包着）⇒ 原句"not in `mmreg.h`"就地改成实话，并写死那枚钉引用的出处＝**格式标签** `mmreg.h:2110`⛔ 是 GUID；
    自查＝diff 非 `//` 行 grep 零命中、`numstat 16/5`、`gofmt -l` 0 枚、`waveFormatPCM` ⛔ 删（`wasapi_windows.go` 仍 1 行），票 300 两枚钉形状原样（顶层 1＋子测试 6／顶层 1＋子测试 4，各 `rc=0`））→
    `fbd150e3`(**三对**突变凭据，仓外导出树 `C:/Users/swq/tmp/306r1-mut-20261011-0730`，`git archive` 取、⛔ 在共享工作树改产码)→`71249b05`(`AC#4` 门禁＋越界)。
    ★三枚各一对、红句逐字**互异**（⛔ 用"一次夹具三枚都装上"当三格都闭的读数）：`:192` `0xFFFE`→`0xFFFD` ⇒ `parseWav error = unsupported wav format: **tag=65534** bits=32 …want nil`（顶层 FAIL 2 枚，`rc=1`）；
    `:218` `3`→`2` ⇒ 同句形状但 `**tag=3**`（`rc=1`，FAIL 2）；`:196` `body+24`→`body+26` ⇒ `**tag=0**`（`rc=1`，FAIL 2）；
    三对各自 `cmp` 逐字节还原 `rc=0` ＋ 定向复绿 `rc=0`，且每对另有一发"新用例 `mv` 走、产码带突变⇒整包照绿＝`rc=0`／顶层 FAIL 0／**PASS 42**"（＝本票 §现量 那两发 `rc=0` 的同台面复跑与逐字对拉，本腿⛔ 再造该证据）。
    终局产码面三把 `cmp`（导出树↔pristine／↔真仓工作树／↔真仓 HEAD blob）全 `rc=0` ⇒ **`wavinjector.go` 一字未动；⛔ 改成引用常量**（正控＝`GOOS=linux GOFLAGS= go build ./internal/audio/` `rc=0`）。
    `AC#4`：`go build ./...`／`go vet`(mod)／`GOOS=linux go vet` 平台中立那一面／`GOOS=windows go vet` 各 `rc=0`；`sh scripts/d22scan.sh` `rc=0` 且 ban #8 自证行逐字留件（`internal/`＝**527**＝派单快照 526＋本腿那枚新文件，它⛔ 豁免于注释、全 ASCII；`cmd/`＝119 逐字同）；
    `gofmt` 三把并排（工作树 0／`5480434f` blob 21 枚 0／现 HEAD blob 22 枚 0）⇒ 新增未格式 0 枚；
    `go test ./internal/audio/ -count=1 -v` 同一台面改前 2 发＋改后 2 发＋同树 `mv` 正控 1 发＝顶层 PASS `42`→`44`（＋2＝本腿两枚，子测试那把尺 12⛔ 变，两把尺⛔ 相加），
    红名册尺 `^--- (FAIL|SKIP):` `comm -13`／`comm -23` **五对双向 0** ⇒ 成论只写**"稳定新增红 0 枚"**并具名该尺盲区（只红过一次会被排除；五份名册逐字相同＝只有既有 `--- SKIP: TestLiveWasapiSmoke`，⛔ 本腿造的、本腿⛔ 造任何 SKIP）；
    `bash scripts/portable-tests.sh --scope=census` `rc=0`、totals 行逐字未变＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（audio 那一行 `10/0`＝8→10）。
    越界＝逐笔 `git show --name-only --format=%H` 四笔写面外 **0/0/0/0**，并集 51 枚＝`internal/audio/**` **2** 枚＋`.scratch/wisp/probes/306/r1/**` 49 枚；被跟踪 `.wav`＝**0**（⛔ 开甲形那处先例）；`frontend/`＋`design/` 命中 0；
    禁改清单 11 枚逐枚 `git diff --quiet` `rc=0`（含 `wavinjector.go`／`wasapi_windows.go`／`thresholds.go`／`allowlist.txt`／三枚冻结件／`PLAN.md`／`parse_wave_format_300_windows_test.go`）；scoped porcelain 起手与终局各 **0 行**（整棵树 828 行别家脏面⛔ 碰⛔ 暂存）；⛔ push。
    ★**具名报回七处与盘上不符／与派单不符**（详见 `probes/306/r1/90-final.md` §5，本腿一律以盘上为准）：① 起手锚 `a81c2980` 在本腿两次读数之间被编排者的 `5480434f` 顶掉（漂移区间 `-- internal` **0 行**、三枚行号在新旧 blob 逐字相同 ⇒ 行号与判据⛔ 受影响；随后 HEAD 又推进 `c888420e`，同样⛔ 碰 `internal/`），本腿以盘上为准重锚；
    ② `AC#2b` 那句"只改那一处偏移、语句其余字不动"的字面那一发在盘上**是 `panic: index out of range [1] with length 0`**（`data[body+26 : body+26]`＝合法空切片喂 `Uint16`，会炸掉整个测试二进制、另一枚指名用例根本没跑到）⇒ 本格主凭据取"同一枚切片上下界一起移"（`data[body+26 : body+28]`，语句其余字⛔ 动）那一发，字面那一发⛔ 丢、作为**第四发**完整交上（`mut-ding-196-literal-wording-*`）；
    ③ `AC#3` 那格派单写 `ksmedia.h:851-855`，盘上守卫起 `:850`、`DEFINE_GUIDSTRUCT` 在 `:854` ⇒ 注释按盘上写 `ksmedia.h:850-855`（`mmreg.h:2480-2484`/`:2483` 与派单逐字对得上）；
    ④ ban #8 `internal/` 本腿现量 **527** vs 派单快照 526（差的 1 枚＝本腿新文件，⛔ 缺陷）；⑤ §现量 与表② 那句"全仓**唯一**的 wav 夹具编码器"自 `1b34e475` 起⛔ 再唯一（`writeWav` 仍是唯一 PCM16 编码器、10 枚调用点⛔ 动；仓里多了第二枚 in-code 编码器 `extFloatWav306`，只产 extensible＋float32 面）⇒ 那句是**世代读数**，本腿⛔ 改原句；
    ⑥ 乙形落在**新文件**而⛔ 扩写 `wavinjector_test.go`（派单写"新增／扩充"两种都授权；选新增是为了三对突变的红册归因⛔ 脏、且⛔ 动那枚载 10 枚调用点文件的既有断言）；
    ⑦ 改前那一发 **census ⛔ 能在导出树取**（`portable-tests.sh` 要问 git，`tar` 解出的树⛔ 带 `.git`）⇒ 另取一把能跑的尺（顶层 `func Test` 计数 `43`→`45`，差 2＝本腿两枚）；"改前／改后同一台面"本腿取的是两枚仓外导出树（同一台机器）＋一枚**同树 `mv` 正控**，"＋2" 那枚 delta 由同树量出。
    ★**自报三处自己的尺缺陷**（⛔ 藏着，件里逐字保留原状）：突变驱动脚本往它自己要 `grep` 的那枚件里追加汇总 ⇒ GNU grep 拒绝（`grep: input file … is also the output`），名册与红句已**另起一把尺**重导进 `mut-*-SUMMARY.txt`（`$(grep -c)` 那三枚计数走管道故为真读数）；第一版越界件里 `--format=4e73d0a5` ⛔ 是合法 pretty 串（`fatal: invalid --pretty format`）且同块 `rc_show` 取的是管道尾 `tail` 的退码（＝派单明令⛔ 的那枚雷的另一种形态）⇒ 整件改名 `40-outofbounds-attempt1-flawed.txt` 保留、第二版逐枚单独重定向重跑；`echo` 里写了反引号⇒命令替换真跑了 `git archive`，usage 灌进 `40-outofbounds.txt` §8（读数⛔ 受影响，已追加 §9 更正）；另：本腿曾在仓库根 `rm -f` 两枚自己的临时计数件，**撞派单第 3 节第 7 条"临时件只建不删"**，内容⛔ 丢（已逐字在 `40-gofmt.txt`）但动作⛔ 合规，同批另六枚走的是 `mv`。
    ⇒ **本腿⛔ 自勾任何 `AC` 框、⛔ 写任何"本格闭合"判语**（裁决者≠实现者）：现态＝票面 **1 勾／5 未勾**原样、本腿只追加这一条 bullet（⛔ 改原句、⛔ 动 `-done`）；交件＝`probes/306/r1/{00-anchor,10-fixture,20-nails,30-ac3-comment,40-gates,90-final}.md`＋`logs/**`（6＋66 枚，`find -size 0`＝**0**、`.out` 命名件＝**0**）；next＝非实现者 **`306-v1`**（三对突变可逐枚 `bash logs/mutate.sh <树>` 重跑；`AC#3` 那格的"只动注释"自查尺一句贴在 `90-final.md` §6）→ `305-a1` → `296-r2`；⛔ push 由本笔执行；残余⛔ 在本票射程：R2（`wasapi_windows.go:378`／`Drain 0.0%`／`convertPacket 0.0%`＝票 300 形 C5，要 owner 另批）＋`:193-195` 的 `size<40` 错误支与裸 `tag=3` 面（本腿**刻意⛔ 装**，＝"另一族缺口"）＋"injector 该⛔ 该被生产调用"那一问（§2 已永久归口票 07／票 228 消费腿／票 247 落点甲）
  - [2026-10-11 08:2x +0800] agent=**`306-v1`（非实现者验收腿，⛔ 参与过本票任何前程）** did=五格逐格独立取证并交判语（⛔ 翻框、⛔ 动 `-done`、⛔ 抄 `306-r1` 任何一枚读数）；件＝`probes/306/v1/{00-anchor,10-sensitivity,20-scope-and-gates,30-judgments,90-final}.md`＋`logs/**`（三枚驱动脚本，⛔ 用腿的 `logs/mutate.sh`），台面＝仓外导出树 `git archive e4740e35 | tar -x`（⛔ 共享工作树做任何突变）。★三次重锚具名：派单那句 HEAD＝`6284a489` 一进场就腐烂（本腿量到 `e4740e35`→`05b37976`→`17dd4a6f`），而 `git diff --quiet e4740e35 17dd4a6f -- internal/audio`＝**rc=0** ⇒ 漂移⛔ 移动任何一枚读数；所有区间尺右端点改写显式 SHA 重跑（件 `logs/A0-rerun-at-head.txt`）。**判语五句**＝`AC#1` **成立**／`AC#2` **成立**／`AC#2b` **附条件成立（形＝甲）**／`AC#3` **成立**／`AC#4` **附条件成立**。★必答 1（反形敏感度）三形并排本腿自己跑：形一 取坏值⇒红（`:192`→`0xFFFD`＝`rc=1`／FAIL 2／红句 `tag=65534`；`:218`→`2`＝同形 `tag=3`；`:196` 上下界同移⇒同形 `tag=0`；红句三枚互异）；形二 把六枚期望侧逐枚改恒真（`VA_inert_count=6`）⇒三枚取坏值**照红**（`err` 门是第二道独立牙）；形三 **拔掉 `err` 门只留值断言**（`VR`，只跑 `TestParseWavExtensibleFloat32306`）⇒ 仍红且红句逐字是解析值（`rate = 0, want 48000`／`chans = 0, want 2`／`len(samples) = 0, want 8`）⇒ 那枚钉⛔ 是"确实返回了 err 类型"那一形；★另加第四形＝逐枚翻转判据（9 枚 `!=`→`==`）⇒ **9/9 各自红**（含 `sample 0 = 32767, want 32767` 这种"值确实在流"的凭证）⇒ 文件里⛔ 恒真装饰、⛔ 空循环。★第四形那问（"守卫所需的世界到底造出来没"）＝覆盖块图三发并排（set 模式）：夹具 `mv` 走⇒`193.5,193.18`／`196.5,196.65`／`219.3,221.26`／`222.4,223.1`／`224.3,224.43` **五块全 0**（＝票面 §现量 那六块的本腿复跑）；夹具在⇒五块全 1；**只跑那两枚新用例**⇒五块仍全 1 而 `212.3,214.26`（PCM16 支）与 `226.3,226.111`（default 支）**恒 0** ⇒ 走的是 `:192`→`:196`→`:218` **同一支**，⛔ "一支有执行者、其余两支仍零"。⚠ 具名残余（⛔ 扣本格＝票面 §边界⛔ 买那一族）：`194.6,195.1`（`size<40` 错误 return）与裸 `tag=3` 面在三发里恒 0。★必答 3（`AC#2b` 那道选择题）＝**甲**：本腿两发都重跑，字面那一发（`data[body+26 : body+26]`）量到 `rc=1`／FAIL **1 枚**／`panic: runtime error: index out of range [1] with length 0 [recovered, repanicked]`，另补一发整包 `-v` 钉死后果＝顶层 PASS **44→36**、两枚 306 用例里只有第一枚 `=== RUN` 过（第二枚指名用例根本没跑到）；⇒ 判语：**panic 证的是"有人走到 `:196`"（无夹具那一发 `rc=0`、`panic` 0 枚＝可达性正控），⛔ 证"那枚偏移被钉着"**——崩⛔ 产生一句关于读出来的值的判词，且它对"挪到仍越界"之外的取值（界内错偏移）根本不响；炸二进制与断言失败在"有没有牙"这件事上**⛔ 同一回事**。同盘先例支持甲＝姊妹钉 `parse_wave_format_300_windows_test.go:116-119` 把同一形状（`data1At: 26`）钉成一枚带 `wantTag: 0` 的子测试（本腿现量 PASS）。⇒ **请编排者在票面追加一条更正**（append-only、⛔ 删原句、⛔ 动框）：`AC#2b` 字面那发产出⛔ 一句断言红，本格凭据取"上下界同移"那一发，字面那发降为可达性正控。★必答 2（⛔ 顺手放宽）三把尺本腿现跑：`git diff 5480434f..17dd4a6f -- internal/audio`＝只有 `M wave_format_float_300_windows_test.go`（16/5）＋`A` 那枚新测试文件；`git diff --quiet … -- internal/audio/wavinjector.go`＝**rc=0（产码 0 字节）**、`wavinjector_test.go`（载 10 枚 `writeWav` 调用点，本腿逐枚数过＝剥定义 10／含定义 11）＝**rc=0 一字未动**；`AC#3` 那把"只在 `//` 行"的尺逐字重跑（腿那把＋本腿更严的一把，缩进 `//` 也豁免）＝**各零命中**；旁证＝票 300 两枚钉形状本腿自跑（顶层 1＋子测试 6／顶层 1＋子测试 4，`rc=0`）、`waveFormatPCM` 仍 1/7 枚。★必答 5（`AC#3` 三枚 SDK 行号逐枚在本腿机器那份 `10.0.26100.0` 树里核）：`mmreg.h:2418`／`:2110`／`:2376`／`:2480-2484`（`DEFINE_GUIDSTRUCT` 在 `:2483`）／`ksmedia.h:850`（`#if defined(_INC_MMREG)`）／`:854` **全部逐字对上**⇒ 注释写盘上的号＝**对的**；票面与派单那句 `ksmedia.h:851-855` 是守卫起点晚一行（`:851` 是里层 `#if !defined( STATIC_… )`），存在性⛔ 错、号段粒度写偏（本腿再加一枚注脚：外层守卫的 `#endif` 在 `:856`，⛔ 扣分）；票面那句"兄弟件 `:101` 引 `mmreg.h:2483`"本腿逐字核中。★必答 6（门禁与名册）：`GOFLAGS= go vet ./internal/audio/` 工作树**rc=0**、`sh scripts/d22scan.sh` **rc=0**（ban #8 射程 `internal/`＝**527**、`cmd/`＝**119**，＋仪器自证行同数；派单快照 526 已过期）；格式名册两把并排（工作树 `internal/audio` 0 枚／blob `5480434f` 21 枚 `.go` 0 枚／blob HEAD 22 枚 0 枚）⇒ 新增 0 枚；⚠ 全仓那把尺本腿量到 `rc=2`、名册 44 行而 `internal/audio` 命中 **0** 行⇒⛔ 写"全仓 0 枚"；红名册**⛔ 用交集尺**＝单发整包 `-v` **5 发**（每发 `rc=0`／`^--- FAIL`＝0／`^--- PASS`＝44／`^    --- PASS`＝12，两把尺⛔ 相加）逐对 `diff` rc=0、五份名册逐字同一枚 `--- SKIP: TestLiveWasapiSmoke`（本腿钉了它的世代＝`5480434f:internal/audio/hotplug_test.go`，⛔ 本腿所造、⛔ 别家新造）⇒"只红过一次被藏住"那一形在这 5 发里没出现；改⇒红那侧的整包名册 `D6`/`D7`/`B4`＝FAIL 恰好那 2 枚、其余顶层 42 枚逐字绿（⛔ 连带红）；无夹具那侧 5 发（`D1`-`D5`，含 `:196` 两形）＝`rc=0`／PASS 42 ⇒ 票面 §现量 那两发在本票夹具之外确实零尺，本腿把它扩量到第三枚。`AC#4` 附的两枚条件都⛔ 本腿车道（欠编排者）＝① 共享工作树 `GOFLAGS= go build ./...` 那一发（本腿只取到导出树 `rc=0`）② `scripts/portable-tests.sh --scope=census` 那一发（跨包，`cmd/wisp` 缺 sherpa PATH 会退 `0xc0000135` 而零 `--- FAIL` ⇒ 腿那枚读数本腿⛔ 复认、⛔ 采信）。★必答 7（腿 §3 第 14 条那枚 `rm -f`）本腿按盘上裁＝**(a) ⛔ 证据丢失**（那两枚是 gofmt 名册的计数中间件，读数逐字在 `r1/logs/40-gofmt.txt` 那 26 行里，且本腿另起三把尺独立复跑同一问：0/0/0 枚与 21→22 那枚分母全对上）；**(b) ⛔ 因此让任何一枚读数不可信**——但罪名定准：真正的风险面⛔ "删"，是**在共享仓根造未跟踪件**（本腿现量根上别家未跟踪件 47 枚，每一枚都在移动别家"整棵树 porcelain"那一类读数），`rm -f` 只是又动了一次共享台面、⛔ 恢复合规 ⇒ 判**动作违规、后果零枚读数受损**；腿在 §7 那句"⛔ 撤回那枚 `rm` 自报"本腿判**⛔ 该撤**（违规⛔ 因后续改做法而消失，要留的是它自己写下的定式：临时件从一开始就落自家 `logs/`、⛔ 落仓库根）。★**与盘上不符处具名 9 枚**（全文 `90-final.md` §5）：① 派单 HEAD 快照腐烂；② 派单现态"1 勾／5 未勾"＝对；③ 票面/派单 `ksmedia.h:851-855` 守卫起点晚一行；④ **腿件 `r1/00-anchor.md` 那句"工作锚＝`5480434f`"是过期读数**（腿的起手笔 `4e73d0a5` 的**父笔是 `5f891575`＝`305-a1` 的锚**，落锚那一刻 HEAD 已晚一笔；影响面＝0，因 `5480434f..17dd4a6f -- internal/audio` 只有那 2 枚文件）；⑤ 腿 §1"6＋66 枚"过期（盘上 6＋**72**；`find -size 0`＝0、`.out`＝0 两把仍对），另盘上 `r1` 78 枚 vs tracked 77 枚＝腿自己的 `logs/90-commit8.txt`（记录它第 8 笔那枚件）**未提交**；⑥ 腿 §7"六笔并集 75＝2＋72＋1"与本腿"八笔并集 80＝2＋77＋1"＝口径差⛔ 矛盾；⑦ 派单 ban #8 526＝过期（527/119）；⑧ 票面 `AC#2b` 字面判据盘上产出⛔ 一句断言红（要追加更正）；⑨ 票面 `AC#4` 那句"写面只落 `internal/audio/**`＋自家 `probes/306/**`"与"末笔追加票面"两句互⛔ 全相容⇒ 记成口径张力⛔ 记成越界（腿那笔票面 diff＝26 加 0 删＝⛔ 改原句、⛔ 动框，本腿这枚 bullet 同样只追加）。⚠ 台面 caveat（⛔ 归本腿、⛔ 动）：本腿两把共享工作树的尺落在别家的 `go.mod` 脏面上（`go-webview2` indirect→direct，1/1），但 `git diff --quiet 5480434f HEAD -- go.mod go.sum`＝rc=0 ⇒ 腿与本腿用的都是已提交那枚；`.gitignore` ＋9 行、`design/**`、别家 `probes/**`、票 303 文件⛔ 碰；仓根另有一枚 0 字节文件名叫 `-`，`stat` 两次 mtime（`08:11:06`／`08:12:33`）落在本腿两批发之间而本腿列出自家全部重定向目标（件 `logs/80-dashfile.txt`）逐枚⛔ 指向它⇒⛔ 能归因、交编排者裁。★本腿自己一枚同类雷具名自报：驱动里那句 `diff a b | sed … >> 件` 之后 `echo $?` 取的是 `sed` 的退码（派单明令⛔ 的那枚雷的又一种形态）⇒ 那行 `rc_B1_B2_roster_diff=0` 作废、本腿另起一把干净的（`rc_diff_B1_B2=0`），坏件与更正件**并存⛔ 删**。⇒ **本腿⛔ 翻任何框、⛔ 写"本格闭合"这种判语之外的东西**：现态＝票面 **1 勾／5 未勾**原样；本腿四笔＝`ec87076f`(起手锚)→`a458114b`(敏感度＋放宽断言尺)→`ffc41214`(五格判语＋锚件就地补重锚)→本笔（票面追加），逐笔 `git show --name-only` 写面外＝**0/0/0/0**、`.go` 写面枚数＝**0**；⛔ push；残余⛔ 在本票射程：R2（`wasapi_windows.go:378`／`Drain 0.0%`／`convertPacket 0.0%`＝票 300 形 C5，要 owner 另批）＋`:193-195` 错误支与裸 `tag=3` 面（腿刻意⛔ 装，本腿复核它那句自陈与盘上一致）＋"injector 该⛔ 该被生产调用"那一问（§2 已永久归口票 07／票 228 消费腿／票 247 落点甲）

## 编排者收 `306-v1`（2026-10-11 08:2x，五笔 `ec87076f`→`a458114b`→`ffc41214`→`56634d23`→`1eaa91f6`，件 `probes/306/v1/**`＝五份正文＋`logs/**`；⛔ 本节目＝我自己复跑的读数）⇒ **五格全翻**＋★**两枚 append-only 更正落在我账上**＋★**裁决表 1:1 已建**（`docs/evidence/s1/306-inline-literals-v1.md`）

### 0. 台面与越界（⛔ 抄腿自陈）

- 起手锚＝`6284a489`（派单给的号），腿自己**三次重锚** `e4740e35`→`05b37976`→`17dd4a6f`，我复核＝三枚之间 `internal/audio` **0 字节**（`git diff --quiet e4740e35 17dd4a6f -- internal/audio` rc=0）⇒ 漂移⛔ 影响本票任何一行号与判据。**定式再钉一遍：派单里的 HEAD 号只算派发时刻快照。**
- 逐笔越界尺我在 **13 笔**（`306-r1` 八笔＋`306-v1` 五笔）上自己重跑 `git show --name-only --format=%H`：每笔**写面外＝0**、`.out`＝**0**；并集尺＝`5480434f..HEAD` 里碰 `internal/audio` 的**只有那 2 枚**（新测试文件＋`AC#3` 那枚注释文件）。件＝`probes/306/orch/logs/orch-03-roster-306-commits.txt`。
- ⚠ **`306-v1` 第 2 笔 `a458114b` 把两枚 HEAD blob 副本命名 `.go` 落进了 `logs/`**（那一笔 新go＝2），第 5 笔 `1eaa91f6` 走 `git mv` 改成 `.txt`、内容 `cmp rc=0` 两把都过 ⇒ **把 `.go` 的分母还给别家的尺**（`gofumpt`／`d22scan` 的射程按 `.go` 数），⛔ 删、⛔ 追认成缺陷。这一枚我复核过：现 HEAD 上 `git ls-files .scratch/wisp/probes/306/v1 | grep -c '\.go$'`＝**0**。

### 1. 承重读数我自己在盘上现跑（⛔ 抄腿；每条带尺；台面＝共享工作树 HEAD `1eaa91f6`＋一枚仓外导出树 `C:/Users/swq/tmp/wisp-306-orch-20261011-0821`）

| 腿的读数 | 我这发的独立复跑 |
|---|---|
| `AC#1` 覆盖块图三发并排：夹具 `mv` 走⇒五块全 0，夹具在与只跑那两枚新用例⇒五块全 1，PCM16 支与 default 支恒 0 | 我**⛔ 复跑块图**，用的是**等价的那一把**：发 A（pristine 导出树、具名两枚）＝`rc_A=0`／`--- PASS`=**2**／FAIL=0／SKIP=0（这一发同时是**台面自检**——具名用例真发出了颜色，⛔ 拿退码当"跑过了"）；发 D（把 `:218` 那一行的 `bits == 32` 改成 `16`）＝`rc=1`／FAIL **恰好那 2 枚**、红句逐字 `unsupported wav format: tag=3 bits=32 (want PCM16 or float32)`。⚠ **口径差别具名**：票面那枚形是 `fmtTag` 的 **`3`→`2`**（腿 `sens-driver.sh:21` 跑的就是那一形，红句同为 `tag=3 bits=32`），我这发改的是同一行的 `bits` 侧 ⇒ 本格凭据以**腿那一形**为主、我这发作对拉 |
| ★"⛔ 断被调用过"那一问＝拔掉 `err` 门后仍红且红句是解析值；另九枚逐枚翻转全红 | 我另起一把更窄的**等号倒置**尺（⛔ 动产码、只动用例）：**F1**＝`if rate != 48000`→`==` ⇒ `rc=1`／FAIL 1、红句逐字 `rate = 48000, want 48000`；**F2**＝`if len(inj.samples) != …`→`==` ⇒ `rc=1`／FAIL 1、红句 `injector holds 8 samples, want 8` ⇒ 断言**逐枚 Binding、⛔ 恒真装饰**（腿那九枚覆盖的是全族，我这两枚抽验的是"翻转必响"这一形） |
| `AC#2` `:192` `0xFFFE`→`0xFFFD`⇒两枚指名用例当场红、无夹具正控整包照绿 | 发 B 同形＝`rc_B=1`／FAIL **2**、红句逐字（`…306_test.go:177` 与 `:213`）`tag=65534 bits=32`；**发 E**＝把新用例 `mv` 走＋同一坏值＋**整包**＝`rc_E=0`／顶层 PASS **42**／FAIL **0**／SKIP 1 ⇒ 与腿的 `D2` 同形同读数；还原＝`git show HEAD:…` 回写后 `cmp -s`⇒`identical rc=0`，随后发 Z＝`rc=0`／PASS 2 |
| `AC#2b` 甲形（上下界同移）⇒ `tag=0` 的正常判词；乙形（派单字面）⇒ `panic` 且第二枚根本没跑到 | 发 C＝`data[body+24 : body+26]`→`data[body+26 : body+28]`＝`rc_C=1`／FAIL **2**／红句逐字 `unsupported wav format: tag=0 bits=32 …`。**两枚腿（落地腿＋验收腿）各自独立量到"字面那一发⛔ 是断言红而是 `panic`"，我今天⛔ 复跑乙形**（它⛔ 是本格凭据，只算可达性正控，腿件里完整留档）⇒ **裁：本格凭据取甲形**，见下面第 2 节那条追加更正 |
| `AC#3` 只动注释：三把"非注释行"尺全零命中、`numstat` 16/5；SDK 行号逐枚核中 | 我自己跑那把逐字尺＝**零命中**（`rc_grep_expect_zero=1`），`numstat` 逐字 **16 加／5 删**；★**SDK 行号我自己去本机那份头文件里逐枚 `sed`**（路径＝`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/`）：`mmreg.h:2110`＝`#define WAVE_FORMAT_IEEE_FLOAT 0x0003`✓、`:2376`＝`#define WAVE_FORMAT_EXTENSIBLE 0xFFFE`✓、`ksmedia.h:850`＝`#if defined(_INC_MMREG)`✓、`:854`＝`DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`✓、`:856`＝`#endif` ⇒ **腿写进注释的号是对的，票面／派单那句 `851-855` 的守卫起点晚一行** |
| `AC#4` 门禁（腿判**附条件成立**，两枚条件⛔ 在它车道） | **条件①我这发补上了**：共享工作树 `GOFLAGS= go build ./...`＝**rc=0**（件 `g1-build-all-worktree.txt`）；**条件②也补上了**：`bash scripts/portable-tests.sh --scope=census`＝**rc=0**、totals 行逐字 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，与落地腿件 `probes/306/r1/logs/40-census-raw.txt:38` **逐字相等**（⚠ 台面⛔ 同一枚：腿那发在它自己的台面、我这发在共享工作树 HEAD `1eaa91f6` ⇒ 只作对拉、⛔ 相减）。另我自己跑 `GOFLAGS= go vet ./internal/audio/`＝**rc=0**、`sh scripts/d22scan.sh`＝**rc=0** clean（ban #8 `internal/`=**527**／`cmd/`=**119**，仪器自证行 `scan_test.go:1261: verified ban #8 internal/: 527` 同数）、整包 `go test ./internal/audio/ -count=1 -v`＝`rc=0`／顶层 PASS **44**／FAIL **0**／SKIP **1**（唯一那枚 `TestLiveWasapiSmoke`，世代＝`hotplug_test.go`，⛔ 本票所造）⇒ **`AC#4` 判成立** |
| 产码一字未动 | 我这把换成**对象层尺**（比 `cmp` 硬）：`git rev-parse 5480434f:internal/audio/wavinjector.go` ＝ `git rev-parse HEAD:internal/audio/wavinjector.go` ＝ **`871f8eb9b3b27830d5ae7c0062268b0a5a3161e2`**；三枚行内权威值在导出树里逐字仍在原位（`:192` `if fmtTag == 0xFFFE {`／`:196` `data[body+24 : body+26]`／`:218` `case fmtTag == 3 && bits == 32:`）；scoped porcelain `-- internal cmd`＝**0 行** |
| 腿的⛔ 做到的四条＋越格自报那一枚 | 我复核成立。那两枚被 `rm` 的 gofmt 计数中间件＝**动作违规、后果零枚读数受损**（读数逐字在 `r1/logs/40-gofmt.txt`，腿另起三把尺独立复跑同一问）⇒ **⛔ 撤回那枚自报**；罪名定准＝**在共享仓根造未跟踪件**（仓根现量别家未跟踪件 47 枚，⛔ 我名下、⛔ 我删）。⚠ 另：仓根那枚 0 字节、文件名是一个减号 `-` 的件（`17` 个硬链接、mtime `2026-10-03 10:20`）＝某次重定向写歪的产物，腿具名报了两次 mtime 落在它两批发之间但⛔ 能归因 ⇒ **我⛔ 删**（别人的件、共享工作树），只在这里记它存在 |

### 2. ★两枚 append-only 更正（原句⛔ 改、⛔ 删，本节追加）

1. **`AC#2b` 的凭据形**：票面那句"只改那一处偏移、语句其余字不动"在盘上产出⛔ 一句断言红——它给的是 `panic: index out of range [1] with length 0`（合法空切片 `data[body+26 : body+26]`），**炸掉整个测试二进制**：顶层 PASS 从 **44 掉到 36**、两枚 306 用例里只有第一枚 `=== RUN` 过 ⇒ 第二枚**根本没跑到**；无夹具那一发 `rc=0`、`panic` 0 枚 ⇒ 它证的是"**有人走到**"，⛔"有人钉"，而且**对界内错偏移根本不响**。⇒ **本格凭据取甲形**＝同一枚切片的**上下界一起移**（`body+24 : body+26` → `body+26 : body+28`，改的仍只有"SubFormat 起点偏移"这一枚权威值，红句是一句正常的关于值的判词 `tag=0`）；字面那一发**降为可达性正控**、完整留档⛔ 丢。同盘先例＝姊妹钉 `parse_wave_format_300_windows_test.go:116-119` 把同一形状（`data1At: 26`）钉成一枚带 `wantTag: 0` 的**子测试**（腿现量它 PASS）。
2. **`AC#3` 的 SDK 号段**：票面／派单写 `ksmedia.h:851-855`；本机 `10.0.26100.0` 那份盘上真形＝守卫起 **`:850`**、`DEFINE_GUIDSTRUCT` 在 **`:854`**、外层 `#endif` 在 `:856` ⇒ **存在性判断⛔ 错、号段粒度写偏**，落地腿按盘上写 `:850-855` 是**对的**。今后引这一格一律按盘上号。

### 3. ★我自己两处判据措辞被两枚腿否证（记我，⛔ 改原句，指回本节）

1. 派单与票面 `AC#2b` 那句字面判据＝**一枚⛔ 可满足的判据**（我写的时候⛔ 知道它给的是 panic）。**先例已记过两次**（`A826` 票 303 `AC#3` 那半句、本轮 `AC#2b`）⇒ **定式＝派单里任何"只改 X 一处、其余字⛔ 动"的突变判据，落笔前先自己在那枚切片／表达式上算一遍它会产生什么形态的失败**（断言红／panic／空操作）。本轮两枚腿各自量到、结论一致，我⛔ 复跑第二发（它⛔ 本格凭据、⛔ 需要我复跑才能定罪它）。
2. 派单里 `ksmedia.h:851-855` 那一枚号＝我从票面抄来的**过期号**（两枚腿与我自己三把尺都给 `:850`）。⇒ **行号一律当快照**这条我自己也要照做，抄票面行号前先 `sed` 一次。

### 4. 欠着的读数（⛔ 算任何一枚腿的欠账，具名归编排者）

- **core 档（ubuntu）那一步的真 CI 日志**：票面 §5 那句"新文件在 CI 两档可见"现前只有**间接**读数（`portable-tests.sh:242` 的 core scope 逐字含 audio＋新文件第一行逐字 `package audio`＋linux 档 `go vet` rc=0）。坐实＝**push 之后取 CI 日志**＝我车道 ⇒ **随下一次推送一起取**。
- 覆盖块图那一把我⛔ 自己复跑（用等价的突变＋红句尺对拉）；三枚 `.cover` 件在 `306-v1/logs/sens/` 里逐块可查，下一位要复核就点那三枚。

### 5. 现态与 next

- 现态＝**票 306 六格全勾**（`AC#0`…`AC#4` 含 `AC#2b`），裁决表 1:1 已建（`docs/evidence/s1/306-inline-literals-v1.md`，README 规则 6 那道闸过），**⛔ 改 `-done` 由我按住到"新文件在 CI core 档那一发的颜色回来"**（理由见第 4 节第一枚；⛔ 收了再红）。`Status:` 行⛔ 写 `done`。
- 残余两件⛔ 属本票射程：① 块 `194.6,195.1`（`size<40` 错误支）与裸 `tag=3`（非 extensible）那一形**仍零执行者**＝票面 §边界明写⛔ 买，谁要闭另立射程；② 上面那枚 CI 色。
- next＝票 306 收口（改名 `-done`＋更新索引）排在**推送＋core 档色回来**之后；本程在飞＝`307-a1`（票 307 `AC#0` 只读）∥ `303-v2`（补票 303 那张缺两行的 1:1 裁决表）；我自己那六发在 `cmd/wisp` 车道（票 305 判死那一族，件 `probes/305/a2/30-deciding-run-recipe.md`）排在**两枚腿交完、编译面空**之后。
- Progress log：
  - [2026-10-11 08:2x +0800] agent=编排者 did=收 `306-v1`（五笔，13 笔越界逐笔我自己重跑＝写面外 **0**、`.out` **0**；`.go` 分母那一枚它自己 `git mv` 还回去了、现量 `probes/306/v1` 下 tracked `.go`＝0）＋**承重读数我自己现跑**（一枚仓外导出树 `C:/Users/swq/tmp/wisp-306-orch-20261011-0821`，⛔ 入库、⛔ 删）：发 A `rc=0`/PASS **2**（＝台面自检，具名用例真发出颜色）／发 B `rc=1`/FAIL **2** 红句 `tag=65534`／发 C `rc=1`/FAIL **2** 红句 `tag=0`／发 D `rc=1`/FAIL **2** 红句 `tag=3`／发 E（`mv` 走新用例＋坏值＋整包）`rc=0`/顶层 PASS **42**/FAIL **0**／F1·F2 等号倒置各 `rc=1`（`rate = 48000, want 48000`／`injector holds 8 samples, want 8`）／两枚文件对 HEAD blob `cmp` 各 `identical rc=0`／发 Z `rc=0`；★**产码零动用对象层尺**（`5480434f` 与 HEAD 的 `wavinjector.go` blob 同为 `871f8eb9…`）；门禁四把我自己跑＝`go build ./...` **rc=0**（共享工作树，＝腿的条件①）、`--scope=census` **rc=0** 且 totals 行与腿件 `40-census-raw.txt:38` **逐字相等**（条件②；⚠ 两把台面⛔ 同一枚，只作对拉）、`go vet ./internal/audio/` **rc=0**、`d22scan` **rc=0** clean（ban #8 `internal/`=527／`cmd/`=119）、整包 `rc=0`／PASS **44**／FAIL **0**／SKIP **1**；SDK 五枚号我自己去 `10.0.26100.0` 那份里 `sed`＝`mmreg.h:2110`/`:2376`✓、`ksmedia.h:850`/`:854`/`:856` ⇒ **腿对、票面与我派单错**；裁决表 1:1 建好＝`docs/evidence/s1/306-inline-literals-v1.md`（六行、每行带尺＋射程＋出处三态）。⇒ **翻 `AC#1`/`AC#2`/`AC#2b`/`AC#3`/`AC#4` 五格**（判语出自非实现者腿 `306-v1`，⛔ 落地腿自勾）；★**两枚 append-only 更正落我账上**＝`AC#2b` 字面判据⛔ 可满足（盘上给的是 `panic`、第二枚指名用例根本没跑到 ⇒ 本格凭据取**甲形**上下界同移、字面那发降为可达性正控留档）＋`AC#3` 的 `ksmedia.h:851-855` 真形 `:850-856`；★**我自己那枚措辞坑第二次踩同形**（派单里"只改 X 一处"的突变判据落笔前⛔ 自己算一遍它产什么形态的失败——第一次＝`A826` 票 303 `AC#3`）；⛔ 改 `-done`＝**按到"新测试文件在 CI core 档那一发的颜色回来"**（现前只有间接读数：core scope 逐字含 audio＋新文件首行逐字 `package audio`＋linux `go vet` rc=0；⛔ 收了再红）；腿⛔ 做到的四条与越格自报一枚复核成立（**⛔ 撤回那枚自报**，罪名定准＝在共享仓根造未跟踪件⛔"删"，读数⛔ 受损；仓根那枚文件名是 `-` 的 0 字节件 mtime `10-03 10:20`⛔ 能归因⇒**我⛔ 删**、只记它存在）。现态＝票 306 **六格全勾、⛔ `-done`**；本程在飞＝`307-a1`∥`303-v2`；next＝腿交完→我跑票 305 那六发（`cmd/wisp` 车道）→推送＋取 core 档色→票 306 收口；⛔ push（照 `A830` §6 两条理由：推送自启 `slo-full` 抢机主 CPU＋票 300/303/111 那批 CI 色要同批取）
