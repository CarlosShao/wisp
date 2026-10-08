# 票 111 `AC#11` — 非实现者对抗验收腿 `111-v1` · 裁决件

- 起手锚：`git rev-parse HEAD` = `a809ab74a65544b8a608cb0f4b60b2d049e8da3b`，分支 `dev`，`date` = `Thu Oct 8 09:11:37 CST 2026`
- 本腿 commit 链：`3fdfa8ca…`（锚）／`e250f0f9…`（Q2）／本件所在笔（Q1+Q3+Q4）
- 被验物：`111-r6` 的六笔（`02505442` 锚／`779714f8` ⓐ／`351e5a5e` 加步／`04b01c6c` ⓑ／`e6981776` 门／`b437d82a` 交付件）
- 票面原文：`.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md:91-94`（本腿 `sed -n '85,100p'` 自取，⛔ 未翻任何框）
- 纪律：**零 push／零开窗（未跑任何 `go test`，更未跑 `-tags winlive` 的 test）／未改产码／未改跟踪测试／未改 `ci.yml`**；
  写面只有 `.scratch/wisp/probes/111/v1/**`。别人的在飞件（` M .gitignore`、`design/**` 16 删＋4 改）一字未动。

**总判**：Q1 部分成立＋一格判不动｜Q2 成立｜Q3 成立｜Q4 成立（有两处非致命缺陷具名）。
**没有任何一问是"不成立"**；也**没有一问是"附条件成立"**——下面逐问给的是拆分后的判，不是笼统的"成立"。

---

## Q1｜那一步在 CI 里到底会不会跑？

### 结论拆分

- **"这一步落在 `test-windows`（`windows-latest`），不是 self-hosted"** ⇒ **成立**。
- **"这个 job 没有任何 `needs:`／作业级 `if:`，所以每次 push 到 main/dev、每个 PR、每晚 cron 都会跑到它"** ⇒ **成立**。
- **"它前面的依赖链是 `Cache third_party`→`build.ps1`→`wisp-cli-tests.sh`，且它紧跟 `:593` 之后"** ⇒ **顺序成立，但"依赖"不成立**（见下，本腿推翻写腿与编排者各一句）。
- **"这一步在某次真实 GitHub run 里出现过 success/failure"** ⇒ **判不动**，缺的是**推送之后的读数**（⛔ 零 push ⇒ 没有任何 run 可查）。
  本腿明确不把"yaml 里有这一步"读成"它跑过"（同台账 `A692`／`A691` §5③）。

### 尺与件

- 作业归属（尺＝`grep -nE '^  [a-z0-9-]+:$|^    runs-on:' .github/workflows/ci.yml`）：
  `lint` 65 ubuntu / `test-core` 395 ubuntu / **`test-windows` 505 → `runs-on: windows-latest`（:506）** /
  `slo-smoke` 778 windows-latest / `slo-full` 836 `[self-hosted, wisp-slo]` / `lint-frontend` 910 ubuntu。
  ⚠ 编排者转述的 739／797／871 是**加步前**的行号，今天一律 **+39**（见 §5）。
- 那一步本体：名 `ci.yml:595`，`shell: bash` `:630`，`if: ${{ !cancelled() }}` **`:631`**，
  `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/` **`:632`**。
  ⇒ 编排者说的"595-632 一带"精确成立（632 就是 run 行）。
- 步骤序列（尺＝`python -c yaml.safe_load` 逐 job 列 steps，件 `logs/q3-gates.txt`）——`test-windows` 共 **10 步**：
  1 `actions/checkout@v4`（**无 `if:`**）｜2 `actions/setup-go@v5`（**无 `if:`**）｜3 winsec ACL｜
  **4 `Cache third_party`（:559）**｜**5 `cgo build smoke`＝`build.ps1`（:566，run 在 :571）**｜
  **6 `cmd/wisp CLI tests`＝`wisp-cli-tests.sh`（:573，run 在 :593）**｜**7 winlive compile gate（:595，本步）**｜
  8 Package coverage census｜9 Portable windows tests｜10 PathResolver junction。
  ⇒ **它确实紧跟在 `:593` 那条 run 之后**，写腿的落点自报为真。
- 步级 `if:` 枚数（尺三把一起量，件 `logs/q1-rulers-and-q4.txt`）：
  `^ {8}if: `＝**13**｜`^[[:space:]]+if: `＝**13**｜YAML 权威逐 job 求和＝lint 4＋test-core 1＋test-windows 8＋其余 0＝**13**。
  作业级 `if:`/`needs:` ＝ **0**（YAML 尺与文本尺 `^\s*(needs|continue-on-error):` 双确认）。
- `if: ${{ !cancelled() }}` 的语义，具名写清（这是"前一步红会不会把它一起吞掉"那一问的答案）：
  - 前一步 **失败（red）** ⇒ `cancelled()` 为假 ⇒ **本步照跑、照答**。
  - 整个 job 被 **取消**（人工取消、或 PR 侧 `cancel-in-progress` 命中）⇒ **本步不跑**，且它下面 8/9/10 同样不跑。
    可达性：`concurrency`（:50-51）只在 `pull_request` 上开 `cancel-in-progress`，push 走 per-sha 组不取消 ⇒
    **push 到 dev 这条路上"取消"支基本走不到**；能走到的是 PR 强推与手工取消。
  - 反向连带：**本步若红，不会吞掉它下面的 8/9/10**——那三步都带 `!cancelled()`（尺＝YAML 逐步列 `if`，件 `logs/q1-rulers-and-q4.txt` Q1.16），
    所以"新步插在中段把 census／portable／junction 三道门打成 inconclusive"这个风险**不存在**。

### ★ 本腿攻出来的一条：那一步的"依赖"是假的，"理由"因此写错了

写腿注释（`ci.yml:616-619`）逐字：

> `# WHY IT SITS HERE (windows leg, after the third_party + build.ps1 +`
> `# cmd/wisp CLI chain that stages the cgo deps this package compiles`
> `# against; precedent ci.yml:561-571 + scripts/wisp-cli-tests.sh): with`

编排者在派单里把它转述成"⛔ 要编译 cgo＋sherpa 头，紧跟那串之后"。**两处都与盘上不符**：

- `go vet -tags winlive` 的整张依赖图里（尺＝`go list -f '{{.ImportPath}}|{{len .CgoFiles}}|{{len .CFiles}}|{{len .SFiles}}' -deps -tags winlive ./cmd/wisp/ ./internal/ball/`，286 枚包，件 `logs/q1-cgo-check.txt`）
  只有 2 枚含 cgo：`runtime/cgo`（1 枚）与 **`github.com/k2-fsa/sherpa-onnx-go-windows`（2 枚）**；
  后者的 `-L` 指向 **`D:/work/base/gopath/pkg/mod/github.com/k2-fsa/sherpa-onnx-go-windows@v1.13.8/lib/x86_64-pc-windows-gnu`**＝**Go 模块缓存**，不是 `third_party/`。
- 全仓 `#cgo` 指令 **0 处**、`import "C"` 在 `cmd/`＋`internal/` 跟踪文件里 **0 处**；`third_party/sherpa-onnx/` 现量只有 3 枚 `.dll`；
  跟踪 `.go` 里没有任何非测试文件引用 `third_party`。⇒ **`third_party` 是运行期 DLL 投放，不是编译期头/库来源。**
- `build.ps1` 的 cgo 环境是**进程内**的：`:79 $env:PATH = "$ccDir;$env:PATH"`、`:262 $env:CGO_ENABLED = '1'`；
  全 `build.ps1`／全 `ci.yml` **没有一处写 `GITHUB_ENV` 或 `GITHUB_PATH`**（尺＝`grep -nE 'GITHUB_ENV|GITHUB_PATH'`），
  `test-windows` 的作业级 `env:`（:507-508）只有 `WISP_ENV: test`。⇒ **第 5 步设的 gcc/CGO 环境传不到第 7 步。**
- 结论：**第 7 步既不依赖第 4 步，也不被第 5 步"喂"环境**。它真正需要的只有一件：**bash 进程默认 `PATH` 上有 gcc**（cgo 编译期）。
  放在这个位置**不多不少都无害**（`!cancelled()` 保它一定答），但**注释给的理由是错的**——
  一个假的依赖声明会让后续程以为"把它上移就坏"，或以为"cache miss 会红它"。**这是缺陷，不是牙的失效。**

### 预测（写成"我预计"，⛔ 不是已证）

我预计托管 `windows-latest` 上这一步**最可能的失败形状＝cgo 工具链解析**，不是 tag 拼错：
GOOS 在 windows runner 上已是 `windows`，所以只剩 `winlive` 一维（这一维被我 Q2 的 A/B 对拉直接钉住），
`-tags` 拼写由 Run A 的红句反证有效；而 `CGO_ENABLED` 若因默认 `PATH` 无 gcc 而落到 0，红会长成
`build constraints exclude all Go files`／`runtime/cgo`、`sherpa-onnx-go-windows` 那两枚包的编译错，
**红在依赖包里、不在那 12 枚里**——那将是"设置红"而非"内容红"，读日志的人必须按这个形状分诊。
第二候选我**排除掉**了一格：工具链版本不齐——`go.mod:3 go 1.27`＋`go.mod:5 toolchain go1.27.1`，
而 `setup-go` 用 `go-version-file: go.mod`（`ci.yml:512-514`），本机 `go version` 正是 `go1.27.1 windows/amd64` ⇒ 同版，
vet 分析器行为不会因版本错位。⛔ 以上全是预测，**runner 上的真实颜色＝欠账**。

---

## Q2｜ⓑ 那两发反形是不是真的有牙？ ⇒ **成立**

### 我自造的坏法（⛔ 未复用写腿任何件）

`cmd/wisp/panel_transport_live_35v2_windows_test.go` 末尾追加一行
`var wisp111v1Probe = undefinedIdent111V1` —— 与写腿三处不同：
**另一枚文件**（写腿＝`internal/ball/live_windows_test.go`）、**另一个包**（cmd/wisp 而非 internal/ball）、
**另一个错误阶级**（**类型错** undefined 符号，而非 `expected declaration, found this` 语法错）。
仓外副本 `D:/tmp/wisp111v1/`（`mutant.go`／`pristine.go`／三份 `overlay-*.json`），⛔ 未写任何跟踪文件。

### 读数（四发，件 `logs/q2-runA-with-tag.txt`／`-runB-no-tag.txt`／`-runC-pristine-overlay.txt`／`-runD-backslash-trap.txt`／`-runD2-raw-backslash.txt`，汇总 `logs/q2-summary.txt`）

| 发 | 命令形状 | rc | 含义 |
|---|---|---|---|
| **A** | `go vet -tags winlive -overlay <mutant> ./cmd/wisp/ ./internal/ball/` | **1** | `vet.exe: D:\tmp\wisp111v1\mutant.go:263:22: undefined: undefinedIdent111V1` ⇒ **有牙** |
| **B** | `go vet -overlay <mutant> ./cmd/wisp/ ./internal/ball/`（摘 winlive 维） | **0** | 同一枚错**完全不可见** ⇒ **今天 CI 的形状确实看不见这 12 枚**（写腿 (ii) 由我第二次量到） |
| **C** | `go vet -tags winlive -overlay <pristine>` | **0** | ★**反空操作正控**：红来自 overlay 的**内容**，不是"挂了 overlay"这件事 |
| **D** | overlay 键路径＝`json.dumps` 正确转义的反斜杠 | **1** | 报的是**我的符号** ⇒ 转义过的反斜杠键路径在 go1.27/Windows **工作正常** |
| **D2** | overlay 键路径＝**裸单反斜杠** | **1** | `go: parsing overlay JSON: invalid escape sequence \w`＝逐字复现编排者那句 |

Run A 的另一层信息：红句**只有我这一枚符号**，`cmd/wisp` 包整体（含其 cgo 依赖）在这台机上 type-check 得动 ⇒
ⓐ 的 rc=0 不是"根本没解析到包"的空绿。

### 还原自证（⛔ 不是只比我改过那一枚）

- 名册尺＝`grep -rlE '//go:build.*winlive' --include=*.go .`：跑前 12 枚／跑后 12 枚，`diff` 空。
- `git hash-object` 全 12 枚：跑前基线 `D:/tmp/wisp111v1/sha-before.txt` vs 跑后 `sha-after.txt`，`diff` **空（12/12 同）**。
- `git status --porcelain cmd/wisp internal/ball` 现量 **0 行**。
- 全树 `git status --porcelain` 从起手 753 → 754，**多出的那行是本腿自己的探针件**，别人的在飞件一行未动。

### 顺手钉住的一条会被读歪的话

写腿红句是 `internal\ball\live_windows_test.go:691:1`（反斜杠）而命令用正斜杠包名——正常 Go 输出；
我这发的红句路径形是 `D:\tmp\wisp111v1\mutant.go`（**指向替换件**）同样是 `-overlay` 下的正常形态。
⛔ 两处都不是"尺坏了"。

---

## Q3｜会不会踩到既有钉子／把别的门打红？ ⇒ **成立**（找到 1 枚真仪器，它不红）

### 能断言 `ci.yml` 的既有仪器：**有一枚**（⛔ 不是"没有"）

- 尺与射程：
  `grep -rn 'ci\.yml' --include=*.go .` ⇒ **3 处、全是注释**（`cmd/wisp/dataroot_128_test.go:242`、`cmd/wisp/logsink_test.go:27`、
  `cmd/wisp/resident_task_source_live_246_windows_test.go:81`）；
  `grep -rln 'workflows' .github scripts tools` ⇒ 5 枚（`ci.yml` 本体＋`scripts/check-path-length-budget.sh`、
  `portable-tests-selftest.sh`、`slo-check.ps1`、`slo-freshness.sh`）；
  再窄一刀，看**谁真的把 ci.yml 当数据打开**：`grep -rn '\.github/workflows' scripts/*.sh scripts/*.ps1 | grep -v 注释` ⇒
  **只剩 `scripts/slo-freshness.sh:99 ci_file=${SLO_FRESH_CI:-$root/.github/workflows/ci.yml}`**。
  `grep -rn 'steps=[0-9]\|STEP_COUNT\|EXPECTED_STEPS' scripts tools .github` ⇒ **空**＝**没有步数基线常量**
  （唯一数步数的脚本在 `.scratch/wisp/probes/111/r5/logs/yaml-guard-census.py`＝探针，⛔ 不是门）。
- **`scripts/slo-freshness.sh` 的 P1 就是那枚仪器**，且它真在 CI 里当门：`.github/workflows/slo-fresh.yml:70 run: sh scripts/slo-freshness.sh`。
  它的射程＝4 条触发器 grep（`on:`/`true:`、`push:`、`branches:[…dev…]`、`cron:`）＋
  一把**按作业锚定**的禁键尺（`sed -n '/^  slo-full:/,/^  [a-z-]*:$/p'` 里禁 `if:`/`continue-on-error:`）。
  新步在 `test-windows`，**落在 slo-full 锚段之外** ⇒ 不会红。
- **被测物跑法**：我不信推理，直接把它 **`:130-172` 的 P1 代码块原样抠出来对真 `ci.yml` 跑**
  （件 `logs/q3-gates.txt` Q3.8）⇒ **`rc_P1block=0`**；逐条 grep 复量也各 `rc=0`，slo-full 体内禁键命中 **0**（件 `logs/q3-instruments.txt`）。
  整支脚本的 P2/P3 需要 `gh`+`GH_TOKEN`（`gh` 在 PATH、**`GH_TOKEN` 未设**）⇒ 全脚本实跑会落 `exit 2`（"没法看"），
  本腿**没有**把那个 2 读成红，也只裁 P1 那半（**具名限制**）。
- `tools/d22scan/runtests.sh` 那道"任何 `--- SKIP` 判致命"（`:88 skipped=$(count '^--- SKIP')`、`:99`）：
  它只读 `go test` 的输出，而新步**不产 test 输出**；且尺 `grep -rn -- '-tags\|winlive' .github/workflows/` 现量
  `winlive` 只出现在 `ci.yml:595/597-603/632`，**没有任何 `go test` 步骤带这个 tag** ⇒ 不会把哪枚用例带进 SKIP。

### 两发门（各落 rc，件 `logs/q3-gates.txt`）

- `cd tools/d22scan && go run . -root ../../` ⇒ **`rc_d22scan=0`**，`clean - no D22 ban violations`
  （现量射程：bans#1-5 `internal/`=228、`cmd/`=38；#6/#8 `frontend/`=85；#7 `internal/tools/`=23；#8 `design/`=39、`internal/`=514、`cmd/`=108 枚 `.go`（含注释与 `_test.go`）；
  另报 1 枚 git-ignored 跳过 `frontend/dist/assets/`。⇒ **#8 的射程不含 `.github/`，注释里的 emoji 不会被它管**（本步注释零 emoji，我核过）。
- `python -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml',encoding='utf-8'))"` ⇒ **rc=0**，
  顺带逐 job 数：`lint 13`／`test-core 7`／**`test-windows 10`（if-guarded 8）**／`slo-smoke 6`／`slo-full 5`／`lint-frontend 11`，`TOTAL_STEPS=52`。
  ⇒ 写腿"步数 9→10"自报为真。

### 那枚 step 的注释射程有没有"冒充"？—— **没有冒充执行／真窗载体**，但有一句现在谎话化

`ci.yml:605-615` 逐字（节选，全文见 `logs/q1-cgo-check.txt` 与我 Q1 段）：

> `# WHAT THIS BUYS, STATED AS COMPILE-ONLY SO IT IS NOT READ AS A SLIDE`
> `# (ticket 111 AC#11 c): it reddens the job when any of those 12 files is`
> `# broken at the compile / vet level and nothing else. \`go vet\` type-checks`
> `# and runs analyzers, it does NOT build a test binary and executes ZERO`
> `# cases, so this opens no ball window and needs no native DLL on PATH.`
> `# It is explicitly NOT a CI carrier for ticket 35's real-window readings`
> `# (:52 / :75(c)) - those take the "host-measurable only" route per A690`
> `# and are NOT wired here; do not read this step as "the winlive cases now`
> `# run on CI".`

⇒ **AC#11 ⓒ 要求的"明写只买编译覆盖"是真的落了**，且明确点名票 35 `:52`／`:75(c)` 不走这里、不冒充 AC#7/AC#9。
**但是**：同一注释块 `ci.yml:599-600` 写的是**现在时**——

> `# were never even COMPILED anywhere in this workflow. \`-tags\` hits 0 times`
> `# and \`winlive\` hits 0 times in this file, and the only \`go vet ./...\``

而它所在的文件今天 `winlive` 命中 **7**、`-tags` 命中 **3**（尺＝`grep -c 'winlive'`、`grep -n -- '-tags'`）。
⇒ **这句话在它自己生效的那一刻就变成假话**（票面原话是"改前 0 命中"，写腿把它以现在时抄进了改后的文件）。
这是**注释缺陷**，⛔ 不是牙的失效；建议改成过去时"as of the pre-step HEAD"。同类：`scripts/slo-freshness.sh:156-158` 那段
注释还在说"ci.yml carries 9 indented `if:` lines …（whose key is at :537）"——今天现量 **13** 枚、`slo-full` 键在 **836**；
它是**注释不是断言**（断言体按作业锚定，见 Q3.8 实跑 rc=0），⛔ 不会红，但会误导下一腿。

---

## Q4｜写腿这份交付有没有"死格"或说谎处？ ⇒ **成立**（无死格、无说谎；两处非致命缺陷）

- **件存在性**：`.scratch/wisp/probes/111/r6/logs/` **16 枚**，逐枚 `stat` 非 0 字节（`anchor.txt` 2152／`ciyml-step.diff` 3909／
  `gate-d22scan.txt` 1384／`negform-with-tag-red.txt` 214／`roster-sha.txt` 1139／`vet-base.txt` 293 …）。
  **无一枚 `.out` 结尾**（尺＝`ls | grep -c '\.out$'` = 0）⇒ 没踩根 `.gitignore:8` 的 `*.out` 静默坑（这一条我起手就现量过：`git check-ignore -v` 确实命中）。
- **两枚 5 字节件**（`vet-winline-base.txt`、`orch-vet-winline.txt`）`od -c` 现量＝`r c = 0 \n` 各一枚——
  即 `go vet` 成功零输出、件里只剩自落的 rc。⚠ 这是本仓"同形连撞两次"的老坑，编排者这次在 `A692` 里自己带了、
  我在 `anchor.md` 里也预先规避；**不算死格**，但它的"命令原文/环境"出处靠 `vet-base.txt`（293 字节，含 go1.27.1 windows/amd64）供。
- **每件事自己落 rc**：`grep -ciE 'rc=|exit code|EXIT='` 逐件命中 ≥1（件 `logs/q1-rulers-and-q4.txt` Q4.1 那张表）。
- **ⓐ 我复跑**：本腿 Run C（`-tags winlive` 对 pristine 内容）＝ **rc=0**，等价于对那两包带 tag 的绿的第二次读数；
  加编排者的 `orch-vet-winline.txt` 与写腿的 `vet-base.txt` ⇒ ⓐ 三度独立复量同色。⛔ 无红名册，且未动那 12 枚任何断言（我的 12/12 hash 对拉就是这条边界的证明）。
- **证据↔结论一致性（首要攻击点）**：全文搜"把编译面说成执行面／把今天绿说成有牙"的滑句 ⇒ `evidence.md` 只在 §5／§6 出现，
  且**方向全部是否定式**（`:128` "不许读成票 35 真窗载体"、`:130` "不冒充 AC#7/AC#9"、`:142` "**⛔ 不把'yaml 里有这一步'读成'它跑过'**，
  因此本腿**不自称 AC#11 完成**"）。`grep -niE 'ci (run|success|green|passed)|run id|workflow run'` ⇒ **0 命中**＝
  **它没有任何一处宣称 CI 跑过**。⛔ 未抓到说谎。
- **自报未跑清单与盘相符**：未跑 `go test`、未跑 `go build ./...`、未造 `scripts/` 包壳——本腿核到 `evidence.md` §4/§6 具名，与令相符。
- **票 111 框数现量**：`- [ ]` = **5**，`- [x]` = **6**（合计 11 枚）⇒ **与编排者给的数一致**；AC#11 仍 `- [ ]`，写腿没翻框（`b437d82a` 的自报为真）。

### 两处非致命缺陷（具名，⛔ 不构成退回"实现"的理由，因为都落在注释/理由层）

1. `ci.yml:616-619` 的**假依赖声明**（third_party/build.ps1 "stages the cgo deps this package compiles against"）——
   由 Q1 那三把尺证伪：cgo 依赖来自模块缓存、`third_party` 只有 DLL、build.ps1 的环境不出进程。
2. `ci.yml:599-600` 的**现在时"0 命中"**在它生效的那一刻失效（同文件今天 7／3 命中）。

---

## 与我转述冲突之处（逐条给可复跑尺）

| # | 编排者的话 | 现量 | 尺 |
|---|---|---|---|
| 1 | "`slo-smoke` 739／`slo-full` 797／`lint-frontend` 871" | **778／836／910**（`lint`65、`test-core`395、`test-windows`505 三枚一致） | `grep -nE '^  [a-z0-9-]+:$|^    runs-on:' .github/workflows/ci.yml` |
| 2 | "`go vet -tags winlive ./cmd/wisp/` 要编译 **cgo＋sherpa 头**，紧跟 cache→build.ps1→cli-tests 那串（＝依赖它）" | **顺序真、依赖假**：cgo 只有 `runtime/cgo` 与 `sherpa-onnx-go-windows` 两包，`-L` 指向 **gopath 模块缓存**；`third_party/` 只有 3 枚 dll；全仓 `#cgo` 0 处、`import "C"` 0 处；`build.ps1:79/:262` 的 `PATH`/`CGO_ENABLED` **不出进程**（无 `GITHUB_ENV`/`GITHUB_PATH`） | `go list -f '{{.ImportPath}}\|{{len .CgoFiles}}' -deps -tags winlive ./cmd/wisp/ ./internal/ball/ \| awk -F'\|' '$2>0'`＋`grep -rn '#cgo' --include=*.go cmd internal`＋`grep -nE 'GITHUB_ENV\|GITHUB_PATH' scripts/build.ps1 .github/workflows/ci.yml` |
| 3 | "严格 8 空格那把 `if:` 尺是 **11**（漏 `:462` 的 `always()`），只有 `^\s+if:` 才得 12→13" | **两把都＝13**；`:462 if: always()` 现量就是 **8 空格**缩进，严格尺抓得到 | `grep -cE '^ {8}if: ' ci.yml` vs `grep -cE '^[[:space:]]+if: ' ci.yml` vs YAML 逐 job 求和 |
| 4 | "JSON 里路径**必须正斜杠**，写反斜杠就会 `invalid escape sequence`" | **过头了**：真陷阱是 **裸单反斜杠**（我 D2 逐字复现该红句）；**正确转义**的反斜杠路径（我 D）照样替换成功、报的是我塞的符号 | 两份 overlay JSON 各跑一发（件 `logs/q2-runD*.txt`） |
| 5 | 票面 `:91-94` 我只转了 ⓐⓑⓒ | 盘上原文**多两枚 ⚠ 交叉核对子弹**（与 AC#7 不同轴／与 AC#9 反向），转述没漏判据但漏了这两条边界；写腿 §6.3 已自报同一差异 | `sed -n '91,94p' .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` |

框数 5/6、12 枚名册、`test-windows`/`windows-latest`、`:595`/`:632`、`shell: bash`、`if: !cancelled()`、
两发门 rc=0、ⓐ rc=0——**这些都与转述一致**，我逐条复跑过。

---

## 我自抓的自己的错（⛔ 不藏）

1. **第一版 cgo 尺是坏的**：我曾用 `go list -f '{{.CDependencies}}'` 取 C 依赖，模板直接报错
   `can't evaluate field CDependencies in type *load.PackagePublic`，那一次输出里混进了错误行（件 `logs/q1-dependency-chain.txt` 仍在，我没删）。
   按现量行号换成 `.CgoFiles`/`.CFiles`/`.SFiles` 三字段才对上（件 `logs/q1-cgo-check.txt`）。
2. **我差点把 Run D 误判成"overlay 空操作"**：D（转义反斜杠）与 A 同红、同一句符号名，我第一反应是"键路径没匹配上、红来自别处"；
   是 **Run C（pristine overlay ⇒ rc=0）** 这一发把它掰回来——若没有 C，我会把"红"错记成"台件可疑"。
   这条也正是 35-v4 那程 `{mutant:mutant}` 假绿的镜像教训：**反空操作的那一发必须单独跑，不能靠 A 的红顺便作证。**
3. **起手一次 `git status` 全量输出太大**（753 行）被工具截断，我改用 `awk` 计数＋排除 `.scratch` 的窄尺重取，才拿到别人的在飞名册。
4. 我把"全树 status 从 753→754"当成潜在污染查了一遍，现量多出的那一行是**我自己的探针件**；⛔ 别人的行没动。

---

## 判不动清单（具名，缺什么写什么）

1. **`test-windows` step 7 在真实 GitHub run 里的 success/failure**——缺**推送**。⛔ 零 push ⇒ 无可查 run；本腿不自称 AC#11 完成。
2. **托管 `windows-latest` 的默认 `PATH` 是否有 gcc**（决定那一步首跑是绿还是"设置红"）——缺 runner 侧读数；本腿只给预测（Q1 末段）。
3. `scripts/slo-freshness.sh` **整支**（P2/P3）对新 `ci.yml` 的 rc——缺 `GH_TOKEN`（`gh` 在、token 不在）；本腿只裁 P1 那半并实跑了它的代码块。
4. `tools/d22scan/runtests.sh` 那道 SKIP 门我**没有实跑**（它要跑 `go test`，与"⛔ 不开窗"相邻而本票无需）——按结构裁：新步不产 `go test` 输出、无任何 test 步骤带 `winlive`（尺见 Q3）。

## 交给编排者的两格小修（⛔ 本腿没动，翻框也归你）

- `ci.yml:616-619` 的依赖理由改写为"编译依赖来自模块缓存＋runner 的 gcc；`third_party` 只供运行期 DLL"；`:599-600` 的"0 命中"改过去时。
- `scripts/slo-freshness.sh:156-158` 那段注释的行号/枚数（9→13、`:537`→`:836`）过期，属下一腿顺手项，不影响其 rc=0。
