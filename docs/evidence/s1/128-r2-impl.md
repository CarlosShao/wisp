# 128-r2 实现件（写腿）：AC#5 两半同笔落地 —— 常驻腿自救句改走 `internal/proc` 单点

- 身份＝**写腿（实现者）**。票面 `.scratch/wisp/issues/128-…-memory.md` **一字未改**（未翻框、未改 `Status:`、未加 `-done`）；台账/`PLAN.md`/`specs/**`/`scripts/**`/`internal/observe/**` **零字节**改动；**未 push**。
- 开工锚点 `e29908c…`。落笔期间 `130-r1` 在同一工作树连落三笔（`b19c6612e`／`8f205a6db`／`66e16de00`）：`git diff --name-status e29908c..HEAD` 逐枚只显示 `internal/observe/**`＋`docs/**`＋`docs/reports/pending-and-issues.md`，**与本笔五枚路径零交集**；它的 `internal/observe/logging.go` 我未碰、未入笔。
- 工作树**零突变**：两发反向判据都在**仓外导出树** `D:\tmp\wisp\20261023\128r2\tree-new` 里做，窗口 `12:18:05→12:19:20`（MUT-A）／`12:19:31→12:20:20`（MUT-B），窗口期**未提交任何东西**（两笔提交均在窗口后）。
- 工具：`go1.27.1 windows/amd64`、`gofumpt` 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`＝`v0.12.0 (go1.27.1)`；`unset GOTMPDIR`；`CGO_ENABLED` 未动。
- 读数落盘（只建不删）：`D:\tmp\wisp\20261023\128r2\{apply-check.out,apply.out,apply-head2.out,runs-old\*,runs-new\*,sentences\*,mut\*,full-cmdwisp-{1,2}.out,analysis{1,2}\*,norm\*}`。

## 1. `git apply`（逐字输出＋rc）

| 命令 | rc | 逐字输出 |
|---|---|---|
| `git apply --check /d/tmp/wisp/20261023/128r1/both-halves.patch` | **0** | 空（`apply-check.out`＝0 字节） |
| `git apply /d/tmp/wisp/20261023/128r1/both-halves.patch` | **0** | 空（`apply.out`＝0 字节） |
| `git -c core.autocrlf=false apply` 同一 patch → `git archive HEAD`（`944bdb9e8`）仓外树 | **0** | 空（`apply-head2.out`）；该树五枚与工作树 `cmp` **逐字节相同**、CR=0 |

- 形状自量：`grep -c '^@@' both-halves.patch`＝**11**；`git diff --numstat` 五枚＝`52/0`＋`6/5`＋`13/10`＋`49/1`＋`104/3` ＝ **+224 −19**（与派单逐字一致）。
- ⚠ **一处必要小修（如实登记）**：patch 末行是 `\ No newline at end of file`（`cmd/wisp/dataroot_128_test.go` 结尾无换行）⇒ `gofmt -d` 只报那一句（`-}`＋`\ No newline at end of file` / `+}`），`gofmt -l`／`gofumpt -l` 都点名该文件（`gofumpt` 输出形是 `cmd\wisp\dataroot_128_test.go`，**用 `[\\/]` 判归属**）。我只补 **1 个 `\n`**（32631→32632 字节），**没有重写任何一行改动**；此后 `gofmt -l`／`gofumpt -l` 全空。

## 2. 改前对照（HEAD 版源码；旧句逐字 67 字节）

- `internal/proc/envfork.go:236-239`（HEAD 逐字）：`dir, err := os.UserConfigDir()` / `if err != nil {` / `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)` / `}` ⇒ **只有 `proc:` 前缀＋OS 原话**，即票面 AC#5 的起点。
- `cmd/wisp/doctor.go:279-298`（HEAD）：`errDataDirUnresolved` ＋ `dataDirUnresolved128` 的**本地模板**（三样都写）——五条 CLI 腿当时读的是这一份。
- 真进程（HEAD 导出树编译的 `wisp.exe`，`APPDATA` 真删、`WISP_ENV=dev`、空 CWD，常驻腿）**拒绝行逐字 67 字节**：

  `wisp: boot failed: proc: user config dir: %AppData% is not defined`

  `sha256sum`＝**SHA-256／64 位**＝`5b2604477f82430a4487f3976f4954beec3455b49d4e2679989e0794204a87e5`；整份 stderr 320 字节（3 行横幅＋1 行 INFO＋该行）。同一发里六枚自救 marker 在场 **0/6**（串 `000000`）。

## 3. 改后（490 字节逐字＋六枚 marker 逐枚在场）

- 同形真进程（工作树编译的 `wisp.exe`，同一脚本同一环境）**拒绝行逐字 490 字节**：

  `wisp: boot failed: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。`

  `sha256sum`（SHA-256／64 位）＝`eddd3f3bced0889c32d9d627b4c18714b69ad7da48d95ad9b608a3373a32f52e`；整份 stderr 743 字节（多一行 `winsec: sealing path resolver installed`）。
- 六枚 marker **逐枚在场**（偏移＝在该 490 字节行内的字节位置，`grep -bo` 现量）：`用户配置目录不可得`@43 · `数据根本应是`@115 · `当前工作目录`@196 · `APPDATA`@396 · `可写目录`@416 · `XDG_CONFIG_HOME`@447 ⇒ **6/6**。
- 改后 CLI 单点产句（同一条串去掉 `wisp: boot failed: ` 前缀）＝**471 字节**，`sha256`＝`4410bb0b95e2a4ef179c1c61180ee4b9469caebf3c5a0d6bf0cf10936fe83c55`。

## 4. 四腿 rc ＋ `assertDirEmpty128`

- 真进程四行（`cmd/wisp` 定向 `-count=2 -v`，两轮同值，日志逐字 `dataroot_128_windows_test.go:70`）：**run rc=2／secret-list rc=2／doctor rc=1／resident rc=1**；四条子用例两轮 **8/8 PASS**。
- `assertDirEmpty128`：该 `-v` 流里 `grep -c "wrote into the start-up directory"`＝**0**（两轮 × 四腿＝**8 发全过**）。我另用独立脚本（`run-legs.sh`，四腿+models+providers）现量每个 CWD 事后 `find` 结果＝**空串**（`cwd=[]`），即 **CWD 一个字节都不多**。
- 定向发汇总：`go test -count=2 -v -run '^(TestAC2[A-Z]|TestAC5|TestDefaultLayout|TestUserConfigDir|TestLayout)' ./cmd/wisp/` ⇒ **rc=0**；尺＝行首 `=== RUN`／`--- PASS`：**顶层 RUN 36（18 用例 ×2 轮）＋子 RUN 36＝`=== RUN` 72 行；顶层 `--- PASS` 36＋子 PASS 36；`--- FAIL` 0／`SKIP` 0**。
- `go test -count=2 ./internal/proc/` ⇒ **rc=0**、`ok … 3.386s`。

## 5. 「五条 CLI 腿文本零变化」我自己复现（471 字节／同一 `sha256`）

- 方法：仓外导出树 `tree-head`（＝HEAD `e29908c`）编 `bin-old/wisp.exe`；工作树编 `bin-new/wisp.exe`；同一脚本（`APPDATA` 真删、`WISP_ENV=dev`、每腿一枚新空 CWD）跑 `run`／`models`／`providers`／`secret`／`doctor`；从各自输出里**机械抽取**「从 `数据根无法解析` 到行尾」那段（`sed 's/.*\(数据根无法解析.*\)$/\1/p'`）比字节。
- 结果：五腿 old/new **各 471 字节、`cmp` 全等**，两侧 `sha256` 同一枚＝`4410bb0b…83c55`（SHA-256／64 位）⇒ **proc 单点产句与 cmd 原来的产句逐字节相同**。
- 整份输出（把 `time=…` 归一后）另比：`run`／`models`／`providers`／`secret` **逐字节相同**；`doctor` 差 **2 行**，差的是它自报 exe 所在目录的两行 `[PASS] DLL colocated: …\bin-old\…` ⇄ `…\bin-new\…`（我的两只 exe 放在不同夹），**不是句子的变化**。

## 6. 反向判据（自造两发）＋还原核

- **MUT-A「退回旧句」**（`tree-new`，只改 `DefaultLayout` 那一行 → `fmt.Errorf("proc: user config dir: %w", err)`；grep 证落地 `envfork.go:286`；`go build ./...` **rc=0** 后才读红）：
  - `./internal/proc/` ⇒ **rc=1**，红名 **`TestDefaultLayoutRefusesWithNoUserConfigDir`**，红行逐字：`AC#5 RED: the refusal is not of the class ErrUserConfigDirUnresolved…` ＋ `…it never says 用户配置目录不可得 / 当前工作目录 / APPDATA / XDG_CONFIG_HOME / 可写目录 / 数据根本应是.`（六枚逐名）＋三条语序红＋句尾红。
  - `./cmd/wisp/` 真进程发 ⇒ **rc=1**，红名 **`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident`**（同一发里 `run`／`secret-list`／`doctor` 三行仍 PASS），红行六条＝`…leg "resident" refused but its line never says "<marker>"` 逐枚 ⇒ **红名点到那枚腿**。
- **MUT-B「摘掉『修法…然后重试。』」**（`tree-new`，删该子句；第一次我的突变把行尾 `+` 留在原处 ⇒ `go build` **rc=1**（`cannot convert … to type interface{Error() string}`），我改的是**突变**、连 `+` 一并去掉；第二次 `go build` **rc=0** 后才读红）：
  - `./internal/proc/` ⇒ **rc=1**，红名同 `TestDefaultLayoutRefusesWithNoUserConfigDir`，红行逐字：`…it never says APPDATA / XDG_CONFIG_HOME / 可写目录.` ＋ `…does not read as a sentence that says what to set…` ＋ `…does not end with the fix and its full stop…` ⇒ **红**。
- **还原核**：还原源＝`mut/pristine-envfork.go`；还原后 `sha256sum internal/proc/envfork.go`＝`a37c4ecc0e867dad1720405feade0fa8fffa4aab080635a4290dca949651bc3a`（**谁算的＝`sha256sum`，算法 SHA-256，64 位**），与还原前 pristine、与工作树同文件**三处逐字相同**；复跑 `./internal/proc` ⇒ **rc=0 绿**。MUT-B 同法还原、同 hash。

## 7. 门禁（逐发 rc，全部本腿现跑）

| 发 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 空输出 |
| `go vet ./internal/proc/ ./cmd/wisp/` | **0** | 空输出 |
| `go test -count=2 ./internal/proc/` | **0** | `ok … 3.386s` |
| 定向 `cmd/wisp`（`-count=2 -v`，见 §4） | **0** | RUN 72／PASS 36×2／FAIL 0／SKIP 0 |
| `sh scripts/d22scan.sh` | **0** | 正控 `PASS=35 FAIL=0 SKIP=0`、`=== RUN=77`；扫描 `clean - no D22 ban violations`（internal/=248、cmd/=42、frontend/=116、internal/tools/=25、design/=39、ban#8 internal/=631、cmd/=148） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；17630 tracked／over-budget 61／roster 61／not-in-roster 0 |

- **格式甲乙两尺**：甲（判据）＝`git archive HEAD` 仓外树（`944bdb9e8`＋patch＋1 换行）⇒ `gofmt -l`、`gofumpt -l` **双双零输出**；乙（只作诊断）＝工作树 ⇒ 同样零输出。CRLF 用 `tr -cd "\r" | wc -c` 现量，五枚各 **0**。

## 8. 整包 `cmd/wisp`：两发都终止（我自己复现，⛔ 当本笔红）

两发都是 `go test -count=1 -v ./cmd/wisp/`（DLL 在 PATH、`unset GOTMPDIR`），**都终止、都无 `panic:`、都无 `test timed out`**；红名册与死点**逐发漂**（按名登记，⛔ 用四数当判据）：

| 发 | 起→止 | 包级行 | 顶层 RUN／PASS／FAIL／SKIP | 红名（逐名） | 死点（有 RUN 无判决） |
|---|---|---|---|---|---|
| 1 | `12:20:49→12:26:13` | `FAIL … 318.640s` rc=1 | 169／167／**1**／0 | `TestAC14GoSideEvalPushReachesThePage` | `TestPanelThreadIsSTAAndExitsCleanly` |
| 2 | `12:26:50→12:32:00` | `FAIL … 304.464s` rc=1 | 168／167／**0**／0 | 无 | `TestAC14GoSideEvalPushReachesThePage` |

- 发 1 多出的那一枚是**在册既有红**（台账 `docs/reports/pending-and-issues.md:20471` 逐字「既有红 1 枚 `TestAC14GoSideEvalPushReachesThePage`（与本票无关）」；真窗族、WebView2 Eval push 维度、颜色随会话漂，见 `:14768`／`:17354`／`:17382`），本笔五枚路径与 `panel*`／WebView2 无依赖路径 ⇒ **新增红 0 枚**；两枚死点属票 314 那族（`A1079`／`A1100`）。发 2 的形状与派单描述的「~301–310 s、零 `--- FAIL`、整进程终止」**逐条吻合**；发 1 的 318.6 s 与那枚既有红**照实登记、不并入"零 FAIL"那句话**。

## 9. 提交（两半同笔；逐个 pathspec 字面点名）

- 代码笔：**`4c40a23f2b6f89397dc05cd951fc33df24dc85d4`**，`5 files changed, 224 insertions(+), 19 deletions(-)`，`git show --name-status` 逐枚：
  1. `cmd/wisp/dataroot_128_test.go`
  2. `cmd/wisp/dataroot_128_windows_test.go`
  3. `cmd/wisp/doctor.go`
  4. `internal/proc/envfork.go`
  5. `internal/proc/envfork_test.go`
- 提交前 `git diff --cached --name-status` 只有这五枚 `M`，**零 `D`／零 `A`**；`internal/observe/**`、`.scratch/**`、`docs/reports/**` 未入索引（`git status` 里别人的 `M` 原样留着）。提交后索引空、五枚工作树干净。
- 件笔：本文件自己的提交＝**`811f127ed`**（首笔）；其后若有一处行内更正，以**追加笔**补入（本仓禁 `--amend`），⛔ 不动上面那枚代码笔。
