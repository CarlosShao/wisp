# 136 — AC#14 实现件（settle 报告的自陈门行）**impl**

> **这份文件是什么**：票 136 `AC#14`（`.scratch/wisp/issues/136-...md:251-271` 那一格 ＋
> `:273-292` 那枚编排者 09-24 18:4x 的**具名解冻＋四裁定**）的实现件。实现程
> `worker-ticket136-ac14`（本程），非验收方，**本文件不翻任何勾、不构成任何一格的结案依据**。
> **授权**：解冻只按票面 `:276-277` 那两段（`internal/observe/sampler.go:431-456` ＋ `:462-521`），
> 边界按 `:278-279`（`StateReport` 侧 `:189-190`/`:331`/`:341-346` 冻结、`buildVerdicts` 不许两用）。
> **本文件只引用变量名与文件名，不含任何凭据值**（本程也没有读到过任何凭据）。

---

## §0 锚点、口径、取数时刻

| 项 | 读数 | 取数时刻（本机） |
| --- | --- | --- |
| 本程自量锚点 | `git rev-parse --short HEAD` ⇒ **`ca2c55e`**（分支 `dev`） | `2026-09-24 18:50:04 +0800` |
| 工具链 | `go version go1.27.1 windows/amd64`；`gofumpt -version` ⇒ **`v0.12.0 (go1.27.1)`**（盘上现量，未照抄票面旧值 v0.7.0） | 18:52 |
| 工作树里别人的东西 | `git status --porcelain` 整树只有 owner 的 `design/**` 16 枚未提交删除 ＋ `design/doubao/`、`design/old/` 两枚未跟踪目录 ⇒ **一枚没碰、没还原、没代提交** | 18:50 |
| 临时件落点 | `D:\tmp\wisp136ac14\`（`bin/wisp.exe` ＋ 三枚 DLL · `readings/` · `snap/` 纯净快照 · 两份 `-v` 日志），**仓库目录内没有新建任何东西**；`build/wisp.exe`（mtime 09-21 15:06，是预检程 §1.2 引过的那枚参照件）**未被覆盖** | 18:53-18:56 |

### 0.1 基线四数（改动**之前**，`go test -count=2 -v ./internal/observe/` 跑两遍）

| 口径 | run1 | run2 |
| --- | --- | --- |
| `^=== RUN`（顶层） | 130 | 130 |
| `^--- PASS` | 130 | 130 |
| `^--- FAIL` | 0 | 0 |
| `^--- SKIP` | 0 | 0 |
| 去重后的顶层名数 | **65** | **65** |
| `grep -ci panic`（整文件命中行数） | 8 | 8 |
| `^panic:`（真 panic） | 0 | 0 |

⇒ **四数口径**：本包的 "-count=2" 把顶层名册翻倍，**65 枚 × 2 = 130**，报数一律带这个乘子。
**`grep -ci panic` 的 8 行命中全是用例名里带 `Panic` 的两枚**（`TestPanicInWorkerSurvivesAndCancelsRoot`、
`TestFakeTreeEmptyScriptFailsClosedAndNotPanics`）各 ×2（`=== RUN` ＋ `--- PASS`）×2（`-count=2`）⇒ **真 panic=0**，
按"名册差集＋panic 逐名"这一条口径写死在这里，免得下一位把 8 读成 8 发崩溃。
**名册差集**：`grep '^--- ' \| awk '{print $2,$3}' \| sort` 两遍 ⇒ `comm -3` **输出 0 行**。
凭据文件 `D:\tmp\wisp136ac14\base-v-1.txt`／`base-v-2.txt`／`b1.names`／`b2.names`。

### 0.2 门禁基线（改动之前）

| 门 | 读数 |
| --- | --- |
| `gofmt -l internal/observe cmd/wisp` | 空输出（rc=0） |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/observe` | 空输出 |
| `go vet ./internal/observe/` | rc=0 |
| `sh scripts/d22scan.sh` @ `git archive HEAD` 纯净快照（`D:\tmp\wisp136ac14\snap`，树外） | **rc=0**；八 scope 命中数 `203 / 22 / 40 / 18 / 16 / 40 / 404 / 39`（`bans #1-5 internal/`、`#1-5 cmd/`、`#6 frontend/`、`#7 internal/tools/`、`#8 design/`、`#8 frontend/`、`#8 internal/`、`#8 cmd/`）；正对照 `runtests.sh` 21/21/0/0。日志 `D:\tmp\wisp136ac14\d22-base.txt` |

台账口径出处：`docs/reports/pending-and-issues.md:3922`（"台账八 scope … 逐格不降"）、`:4192`（"两处各有出处 ⇒ **一枚不降**，不是'没降是因为没人看'"）、
`:4392`（`ban #8 internal/` 401→402 归因到自己那枚 `_test.go`）。⇒ 引用命中数**必须连归因一起引**，本文照此办（§3）。

---

## §1 `Gate` 那一支：先量的读数，逐次列（票面 `:281-285` 要求的前置）

**裁判问题**：门行做成 `Gate: true` 会把 `wisp slo -settle` 在"部分未测"的真窗口上从 exit 0 变 exit 1，
并经 `scripts/slo-check.ps1:381` → `:396` 打到 CI 的 `slo-smoke`（托管）与 `slo-full`（**这台 6C12T 本机**）。
票面 `:283` 要的实现程第一发＝编队安静窗口内 ≥5 次 `wisp slo -settle`、逐次记 `sample_errors`。

### 1.1 取数前置：编队现量

| 检查 | 读数 |
| --- | --- |
| `scripts/slo-check.ps1:153-155` 名单**盘上现量**（逐字取自该脚本，15 枚，**不是派单简报里的 8 枚**：多 `link.exe`/`g++.exe`/`cc1.exe`/`cc1plus.exe`/`as.exe`/`ld.exe`/`wisp-cli.exe` ⇒ 用更全的那份，只可能更严） | 每发读数**前后各量一次** ⇒ **12 次全空**（`roster_before=` / `roster_after=` 皆无值） |
| `gh run list --limit 6` | 在飞的只有 `35989390166`（`ca2c55e` 那次推送，18:48:47Z 起） |
| 同一发的逐 job | `slo-full`（self-hosted，**就是这台机**）= **completed/success**；`slo-smoke` = success；`lint` = **failure（先于本程，见 1.4）**；`test-core`/`lint-frontend` = success；`test-windows` 在飞但 `runs-on: windows-latest`（`ci.yml:335`，托管，不抢本机） ⇒ **取数窗口内 `slo-full` 不在飞** |

### 1.2 六发读数（≥5，逐次）

二进制：本程 `go build -trimpath -ldflags "-X …buildinfo.Commit=ca2c55e" -o D:\tmp\wisp136ac14\bin\wisp.exe ./cmd/wisp`
（`CGO_ENABLED=1` ＋ `E:/work/base/msys64/mingw64/bin/gcc.exe`，rc=0），**三枚 sherpa/onnx DLL 与 exe 同目录**
（`third_party/sherpa-onnx/{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`，简报 `⚠` 那条加载期依赖照办了）。
命令逐发：`WISP_ENV=test`（＝`ci.yml:540` 那枚 job env）＋ `wisp.exe slo -settle -out readings/settle-N.json`。
取数脚本 `D:\tmp\wisp136ac14\readings.ps1`（脚本内**不调用任何 go 工具**），日志 `D:\tmp\wisp136ac14\readings/readings.txt`。

| # | 起时 | exit | `sample_errors` | `len(samples)` | `pass` | `back_within_cap_ms` | `final_bytes` | 名单前后 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 18:54:59 | 0 | **0** | 40 | true | 267 | 10133504 | 空 / 空 |
| 2 | 18:55:13 | 0 | **0** | 40 | true | 263 | 9252864 | 空 / 空 |
| 3 | 18:55:26 | 0 | **0** | 40 | true | 264 | 9064448 | 空 / 空 |
| 4 | 18:55:39 | 0 | **0** | 40 | true | 265 | 9203712 | 空 / 空 |
| 5 | 18:55:52 | 0 | **0** | 40 | true | 265 | 9248768 | 空 / 空 |
| 6 | 18:56:05 | 0 | **0** | 40 | true | 269 | 10350592 | 空 / 空 |

⇒ **六发 `sample_errors` 逐发全 0**（`sample_errors_total=0`；`retries=0` 六次，即没有一发因争用被重取）。
按票面 `:283` 那一支，这枚读数的含义是"**全 0 ⇒ 按 `Gate: true` 落地，这一形今天不会误伤存量合法窗口**"。

### 1.3 为什么本程仍然**没有**把它落成 `Gate: true`（一处前提被盘上推翻，具名停手报回）

票面 `:287` 与派单简报都把"AC#14 要推翻的既有断言"记成**一处**：
`sampler_settle_coverage_136_test.go:208-210` 那一枚 `if !rep.Pass`。预检程 §3.3 末也只点到那一处。
**盘上现量是三处**，同一条"disclosure not verdict"形状，同一枚文件：

| 位置 | 逐字断言 | 它现在要求什么 |
| --- | --- | --- |
| `:143-145`（`TestCheckSettleSingleTrustworthyReadReportsItsLoss`） | `if !rep.Pass {` ⇒ `t.Fatalf("this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v", rep)` | 10 枚读数里只 1 枚可信（丢了 9 枚）那一形**必须仍然 pass** |
| `:208-210`（`TestCheckSettleHalfTheReadsFailedReportsItsLoss`） | `if !rep.Pass {` ⇒ `t.Fatalf("disclosure leg, not a verdict leg: this window still passes, report=%+v", rep)` | 丢一半那一形**必须仍然 pass** ← **唯一被预先授权改的那一处** |
| `:252-254`（`TestCheckSettleZeroFootprintDropsAreCountedToo`） | `if !rep.Pass {` ⇒ `t.Fatalf("report=%+v", rep)` | 1 枚可信 / 2 枚零足迹（丢了 9 枚）那一形**必须仍然 pass** |

**任何"会抓住丢一半那一形"的门行判据必然同时抓住另两形**（那两窗口的覆盖率严格更差：kept=1 / lost≥9，
对照 kept≈lost≈5），除非把判据写成"只红在恰好一半"这种非单调的任意阈值——那是**新造阈值**，
`AGENTS.md §1.1` 与票面 `:279` 都禁。⇒ `Gate: true` ＋折进总 `Pass` 的**最小**代价是**同改三处断言**，
而本程的授权只覆盖其中一处。**扩授权不是本程能做的动作**（`AGENTS.md §0` 第 1 句"未定义即停"、派单"超出即停手报回，不许自己扩"），
故本程落地形状＝票面 `:284` 给的那一档**记录行**（`gate=false`、不折进总 `Pass`），
把"三处 vs 一处"这枚盘上事实连同 §1.2 的全 0 读数一起交回编排者裁：
**批准另外两处之后，`sampler.go` 里要改的是**一枚具名常量的一个布尔字面量**（`settleCoverageRowGates`），
门行的自陈、红句里的 `sample_errors`、折叠处形状**都已就位**，不需要再动别的。
§4 那两发快照变异给出"批准之后会红哪三枚"的**读数**（不是推理）。

本程**没有**把 `AC#14` 判成完成：按票面 `:262-263` 判据① 的原话与预检程 §3.3 末的读法，
记录行那一支会让"不许只靠 `pass=false` 一个总布尔代答"这一维**落空**——本文件如实把这一格标为**未结**。

### 1.4 顺带量到的、不属本格的账（登记，不动手）

- `lint` 在锚点 `ca2c55e` 上就是 **failure**（`gh run view 35989390166 --json jobs`，18:5x 现量）⇒ 本程交件时 `lint`
  的色**不是本程造的**；本程自己的 gofmt/gofumpt/vet/d22scan 逐门读数在 §3。
- 简报里"进程名单 8 枚"与盘上 15 枚不符（§1.1）⇒ 按更全那份量的，只可能更严，不影响"全 0"这一结论。
