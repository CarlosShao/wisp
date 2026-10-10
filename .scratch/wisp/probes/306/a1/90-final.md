# 306-a1 · 90 终局：现态＋我⛔ 做到的＋越界读数

腿＝`306-a1`（票 306 `AC#0`，只读普查）。世代锚＝`f994f956647ebf343d861b588f8d379f7f883b6e`（`dev`）。

## 交件名册（全部在允许写面 `.scratch/wisp/probes/306/a1/**`，⛔ `.go`、⛔ `.out`）

- `00-anchor.md` — 锚＋进程闸门＋台面
- `10-table1-literal-roster.md` — 表①行内权威字面量全名册（自带射程）
- `20-table2-branch-reachability.md` — 表②调用链＋★必答＋响亮报回
- `30-table3-candidate-points-and-cost.md` — 表③候选甲/乙＋代价＋硬前提现量
- `logs/00-tasklist.txt` · `logs/40-cover-run.txt`（`rc_test=0`）· `logs/40-cover.txt`（块图）· `logs/41-cover-func.txt` · `logs/50-linux-build.txt`（`rc_linux_build=0`）· `logs/51-linux-vet.txt`（`rc_linux_vet=0`）· `logs/30-roster-worktree.txt` · `logs/31-roster-wavinjector.txt` · `logs/60-outbound-diff.txt`

## 三张表的结论（各一行，带尺名）

- **①** 同族"行内字面量且零执行者"＝**2 枚**（`internal/audio/wavinjector.go:192` 的 `0xFFFE`、`:218` 的 `3`）＋**1 枚同形未数到的附赠**（`:196` SubFormat 偏移 `body+24`，块 hit=0）；尺＝`grep -rnE … internal cmd tools`（72 行原始读数，`logs/30-roster-worktree.txt`）＋本腿复跑覆盖率块图；同族"引常量但零执行者"＝1 枚（`wasapi_windows.go:378`＝R2，⛔ 本票射程）；戊类同形异族排除 6 枚（具名于表①）。⇒ **洞族没散落到 `cmd/`、`tools/` 或别的 `internal/` 包。**
- **②** 两支零执行者的根＝唯一夹具编码器 `writeWav`（`wavinjector_test.go:27` 逐字 `buf.Write(u16le(1)) // PCM`、`:32` `buf.Write(u16le(16))`，11 枚调用点全走它），手工第二形只走 `:226` default；**★今天⛔ 任何夹具喂过 `fmtTag==0xFFFE`＋`bits==32`**（凭据＝块计数 `193/194/196/219/222/224` 全 0，尺＝本腿自跑 `go test ./internal/audio/ -count=1 -coverprofile=…`，`ok … 15.960s coverage: 59.3%`、`parseWav 70.0%`）；**响亮报回：`wavinjector.go` 整枚文件⛔ 生产调用者**（只剩 `doc.go:16`/`gate.go:29` 两行注释），但牙⛔ 依赖生产接线（`parseWav` 同包可达，先例 `hotplug_test.go:581`）⇒ 是否升级成接线票**归编排者裁**。
- **③** 甲（新增 `.wav` 夹具）与乙（既有 injector 用例里加形）**都能一次买到两枚洞的牙**（extensible+float32 面 ⇒ `:192` 真 ⇒ `:196` 重写 ⇒ `:218` 命中；两枚各自改错都落 `:226` ⇒ 红）⇒ **票面那句"补 float32/extensible 夹具会同时钉住行内 3/0xFFFE"＝验证成立**（推演依据＝逐字分支结构＋块 hit 计数；"改⇒红"两发⛔ 由本腿跑，那要写 `.go`）。两形都⛔ 动 `third_party`、都⛔ 需要 `//go:build windows`、都落进 CI **core(ubuntu)+windows 两档**（尺＝`scripts/portable-tests.sh:242`／`:253`，`ci.yml:420/487`／`:534/797`），pin 是**包名**名册（`core_pin:164` 起、audio 在 `:168`；`win_pin:193` 起、audio 在 `:195`）⇒ 新增测试文件⛔ 动 pin。甲的真代价＝本仓被跟踪 `.wav` 为 **0 枚**（尺＝`git ls-files | grep -icE '\.wav$'`＝0）且 `internal/audio` **无 testdata 目录**（尺＝`find internal cmd -type d -name testdata`）。
- **硬前提已量实**：仓外导出树里把那两枚行内值换成 `waveFormatExt`/`waveFormatFloat` ⇒ `GOOS=linux go build ./internal/audio/` **rc=1**（`undefined: waveFormatExt`／`undefined: waveFormatFloat`），同树 `GOOS=windows` **rc=0** ⇒ "改成引用常量"这一形今天确实让非 windows 档编译不过。

## 我⛔ 复跑 / ⛔ 做到的（越界自查之外的诚实清单）

1. ⛔ 复跑编排者那四发 mutation（`probes/300/orch/logs/r3-orch-v4-attacks-rcfixed-20261010-204046.txt`）——那要动产码，⛔ 本票 `AC#0` 只读；只做了逐字对拉：件里 `rc_R1b_jia_pkg=0`／`rc_R1b_yi_pkg=0` 与"顶层 PASS 42"我**没有独立复现**，本腿的等价凭据是覆盖率块 hit=0（改值不可能红的**原因**），不是那两发的退码本身。
2. ⛔ 复跑票面 `:25` 的 CI 可见性读数（run `38049322056`，件 `probes/300/orch/r10-ci-third-run-new-test-visibility.txt`）——复跑需 push，⛔ 子代理权限；我只量到 pin/scope 两侧的**结构性**理由（包名已在两档 scope 与两个 pin 名册里）。
3. ⛔ 跑过 `sh scripts/d22scan.sh`／整包 `go build ./...`（`AC#4` 的门禁，属落地腿本格；本腿只跑 `internal/audio` 的 build/vet/test 三面，⛔ 整包）。
4. ⛔ 写任何产码或测试：`internal/**`、`cmd/**`、`tools/**`、`.scratch/wisp/issues/**` 全程只读（自查＝`git status --porcelain -- internal cmd tools` 起手 0、终局仍 0；票面框改动数＝`git status --porcelain -- .scratch/wisp/issues | wc -l`＝**0**）。
5. ⛔ 复跑 `300-v4` 那把覆盖尺（派单点名⛔ 引它的数当凭据）——已独立复跑并逐字对上 `parseWav 70.0%`／`Drain 0.0%`／`convertPacket 0.0%`。

## 越界读数（逐笔尺，⛔ 区间尺）

- 逐笔名册尺＝对本腿全部 commit 逐枚 `git show --name-only --format=<hash> <hash>`（前两笔 `4f210227`、`5940f64a` 的名册已逐字贴在下面；终局笔＝本件所属那笔，用 `git log --oneline -3` 现量其号后同尺重跑即可）——**每一行都在 `.scratch/wisp/probes/306/a1/`** 下：

  ```
  4f210227  00-anchor.md + logs/00-tasklist.txt logs/40-cover-run.txt logs/40-cover.txt logs/41-cover-func.txt logs/50-linux-build.txt
  5940f64a  10-table1-literal-roster.md 20-table2-branch-reachability.md + logs/30-roster-worktree.txt logs/31-roster-wavinjector.txt logs/51-linux-vet.txt
  ```
- 作差件＝`logs/60-outbound-diff.txt`：内容只有自落的那行 `rc_outbound=1 (…) ⇒ grep rc=1 且文件空 = 零越界 = 结论成立`（**空集即成立**，本仓规矩：`grep -v` 型作差件的空＝结论）。
- 禁形状自查：`find .scratch/wisp/probes/306 -name "*.out" | wc -l` ＝ **0**；`-name "*.go" | wc -l` ＝ **0**。
- ⛔ push（本子代理只 commit）；⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；每次 commit 显式 pathspec 写在 `$( … )` 之外；⛔ 在仓库目录内建 worktree/checkout（导出树在 `C:/Users/swq/tmp/306a1-export`，⛔ 仓内）。
- ⛔ 读 `%APPDATA%\wisp-dev\secrets\`；⛔ 读机主 config 值；无凭据进件。⛔ 写 `frontend/src/**`、`design/**`；⛔ 动 `thresholds.go`／`allowlist.txt`／三枚冻结件／SLO 阈值／`PLAN.md`／D43 表。
- 台面：⛔ 真机／⛔ 真窗（`wisp.exe`／`balldebug.exe`／`mockllm.exe` 起手与全程 0/0/0；`msedgewebview2.exe` 18 枚⛔ 归本腿⛔ 归我）。包级 `go test ./internal/audio/ -count=1` 跑在真实工作树（该面 `status` ＝ 0，且派单声明本腿是当时唯一 Go 编译面腿）；退码一律 `> 件 2>&1; rc=$?` 形式取，件里各有一行 `rc=N`。

## 与我读到的既有文本冲突处（具名报回，以盘上为准）

- **派单未给锚号**，只提 `d9864baf`（那是编排者四发攻击的导出树世代）；当前 HEAD＝`f994f95`。两枚行内字面量的行号 `:192`／`:218` 在两世代之间**未漂移**（内容锚逐字抄在 `00-anchor.md`）。
- 票面 `:25` 把 `scripts/portable-tests.sh:253` 称作 "`windows` 档"＝盘上逐字核对通过（`scope=(… ./internal/audio/)` 在 `:250-254`，`windows)` 分支起 `:249`）；但票面**没数到 core 档也含 audio**（`:242` `./internal/audio/...`，跑在 ubuntu-latest）⇒ 新增**无标签**测试文件的可见面比票面写的更宽（一条 CI 事实，⛔ 规矩变更）。
