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

- [ ] **`AC#0` 只读普查（决定落点，⛔ 预设答案）**：交回三张表——
  ① **行内权威字面量全名册**：仓里还有多少处把 `mmreg.h` 那族值（`3`／`0xFFFE`／`0x0003`／GUID 低字）写成**行内字面量**而⛔ 引常量／引钉过的测试？尺＝`grep -nE '0xFFFE|fmtTag == 3|== 0x0003' internal/**/*.go` ＋逐枚分类"引常量／行内／注释"，名册必带**射程目录＋blob 还是工作树**（本仓老尺）。
  ② **那两支为什么零执行者**：从入口到 `wavinjector.go:192`／`:218` 的调用链逐跳列（谁调 `parseWav`／谁构造 `fmtTag`），并具名答"今天有没有任何夹具喂过 `0xFFFE`＋`bits==32` 这一组合"。⚠ 若"根本没有生产调用者"，本票就升级成**接线票**，由我另裁。
  ③ **夹具落点候选＋代价**：至少两形候选（新增 wav 夹具文件／在既有 injector 用例里加形）各写"要不要动 `third_party`、要不要 windows 标签、CI 哪一档会跑到"。
- [ ] **`AC#1` 决定性一发（⛔ 先红再修；本格只产**测试侧**夹具，⛔ 改产码语义）**：让 `0xFFFE`＋`bits==32` 那一支**真的被执行**（extensible 头＋float32 子格式），断言的是**解析出来的那个值**，⛔ 断"被调用过"。★凭据形状⛔ 变：把 `:218` 的 `3` 换成 `2` 之后**必须红**（今天的红＝我 §现量 那两发 `rc=0`＝本票的"改前必红"）。
- [ ] **`AC#2` 同族第二枚的牙**：`:192` 的 `0xFFFE`→`0xFFFD` 也必须让指名用例红。⚠ 两枚各一对（改⇒红／`cmp` 逐字节还原⇒绿），缺一枚⛔ 算本格闭合。
- [ ] **`AC#3` 只改注释那一格（最小面，⛔ 与 `AC#1`/`AC#2` 同批落）**：把 `wave_format_float_300_windows_test.go:38-41` 那句改成实话（两份定义各在何处、`mmreg.h:2110` 才是那枚钉引用的出处）。**硬边界**＝⛔ 动任何一行断言／任何一枚期望值；⛔ 删 `waveFormatPCM`（票 300 `AC#6` 的附带覆盖面正读着它，删它＝那件编译不过）；⛔ 顺手改别的注释。
- [ ] **`AC#4` 门禁＋越界（每把尺自落一行 `rc=N`，0 字节件⛔ 算交）**：`GOFLAGS= go build ./...`／`go vet ./internal/audio/`／`sh scripts/d22scan.sh` 各 rc=0；格式名册**两把并排**（工作树＋仓外 blob，各写射程目录）新增 0 枚；`go test ./internal/audio/ -count=1 -v` 改前改后各 **≥2 发**、同一台面（`git archive` 仓外导出树，⛔ 在共享工作树跑"改前那发"）、红名册 `comm` 双向＝**新增红 0**；三数必带尺名（`^--- PASS` 顶层 vs 含子测试）。写面只落 `internal/audio/**` ＋自家 `probes/306/<腿名>/**`；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；⛔ `frontend/**`／`design/**`／三枚冻结件／`thresholds.go`／`allowlist.txt`／D43 表／`PLAN.md` 零字节；⛔ 零 push、commit 显式 pathspec 写在 `$( … )` 之外。

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
