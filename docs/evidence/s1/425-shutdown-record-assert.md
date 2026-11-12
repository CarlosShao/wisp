# 票 425 — 断言落后于产码：`shutdown step` 永远挤不出末位

> **尺**：起手 `HEAD=176d267a2`（`dev`），现跑 `2026-11-12 19:42–19:52 +08`。
> **写面**：`cmd/wisp/resident_sink_nail_127_windows_test.go` 一枚 ＋ 本文件一枚。产码零字节（`internal/proc/shutdown.go`、`cmd/wisp/resident_windows.go` 起手＝终态＝HEAD blob，见 §5）。
> **载体**：`D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"`。
> ⚠ 本机 `bash` 在 PATH 上是 WSL 转发器（`C:\Windows\System32\bash.exe`，`execvpe(/bin/bash) failed`，**一枚测试都不跑**）。四发读数**全部**用 Git bash 取。
> **临时件（只建不删）**：`D:\tmp\wisp\20261112\425-shutdown\`（`ac0-red.log` / `ac1-green.log` / `ac2-mut-red.log` / `ac2b-mut2-red.log` / `ac2c-restore-green.log` / `mut/` / `stage-*/`）。

## 0. 一句话

票 43（`5056a365c`）在 `cmd/wisp/resident_windows.go:246` 的 defer 尾部加了**十步回放**，`wisp: D38(e) step 10 (exit)` 必然排在所有 `shutdown step` 之后 ⇒ 原断言"**末位**必须是 `shutdown step`"**结构上不可满足**，与产品好坏无关。
**产品没坏，是断言没跟上。** 本票把断言从"末位"改锚到"**关停记录发生在回放之前**"这个真正的产品承诺（票面乙形），并**实测**它有牙。

## 1. AC#0 复现（改前，逐字红句）

命令（`WISP_CLI_STAGE_DIR=D:\tmp\wisp\20261112\425-shutdown\stage-ac0`）：

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
```

逐字红句（`resident_sink_nail_127_windows_test.go:635`）：

```
    resident_sink_nail_127_windows_test.go:635: the last record in the sink is "wisp: D38(e) step 10 (exit): ok elapsed=0s", want a shutdown-step record: the shutdown must be logged before the sink closes (records: [{time=2026-11-12T19:43:43.1402942+08:00 level=INFO msg="winsec: sealing path resolver installed" ... }])
```

`--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (14.70s)`，`FAIL`，`[exit code: 1]`。

判定通过条件逐字满足：红句含 `the last record in the sink is "wisp: D38(e) step 10`。

**★旁证（这条最硬，同发日志里就有）**：

```
    resident_sink_nail_127_windows_test.go:640: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="shutdown step skipped (module not present)"
```

⇒ **5 条 shutdown 记录确实在盘上**，`:629` 的 `len(trail) != 0` 那关**本来就过了**。红的只是"末位"这一件事。

## 2. AC#1 选形：**乙**（钉"关停记录早于第一条回放记录"），理由

**为什么是乙而不是甲。** 甲形（"`trail` 末位必须是 `recs` 里最后一条 `shutdown step`"）是**同义反复**——`trail` 本来就是从 `recs` 按 `residentShutdownRecord` 过滤出来的，"最后一条 shutdown 记录在 `trail` 末位"恒真，它读的是**自己的过滤器**，不是产品。要让它有内容就得再叠一条外部次序判据，而那正是乙在做的事。乙直接钉住两个 defer 存在的理由本身：**序列在 sink 还开着时逐条记账，回放随后才跑**。

**改后的三个合取项**（`resident_sink_nail_127_windows_test.go:666-721`）：

| # | 判据 | 抓什么退化 |
|---|---|---|
| 1 | `len(trail) != 0`（**保留**，未放宽） | 删掉 `internal/proc/shutdown.go:130` 的 `slog.Info` ⇒ 盘上一条 `shutdown step` 都没有 ⇒ 红。**§3 实测。** |
| 2 | 步序**严格递增** | 序列被重排、或 `run()` 名册被重排 ⇒ 记录全在、次序错 ⇒ 红。原断言对此**完全无感**。 |
| 3 | **最后一条** `shutdown step` 记录的**记录下标** 早于**第一条** `wisp: D38(e) step` 记录的记录下标 | 回放被挪到序列之前、或回放里插了一条记账 ⇒ 红。**§4 实测。** |

外加一层**时间戳佐证**（`:709-718`）：`trail` 末位的时间戳不得**晚于**第一条回放记录的时间戳（相等放行——平局不是退化证据）。位置读的是 sink 自己的追加次序（`internal/observe` 的写入契约），时间戳只是旁证。

**为什么这仍然能抓住产品退化、不是把测试废掉**：原断言唯一的判据是"末位"，而末位在票 43 之后**恒假**——一个恒假的判据不测任何东西，它的绿也买不到任何东西。改后三条判据**每一条都有对应的实测红**（§3 的红落在判据 1、§4 的红落在判据 3 且判据 1／2 同时保持绿），而且新增的判据 2／3 覆盖了原断言**从来测不到**的两类退化。

**一处反向保护**：若将来回放被删／改名，`replayAt < 0` 会**红**并逐字说明"这一格要重新锚定到哪里，⛔ 不许删掉这个承诺"（`:693-697`）——不是静默跳过。这是刻意的：本票的整个教训就是"锚点漂了，断言要跟着搬，而不是把断言删掉"。

## 3. AC#2 突变控制 ①：删掉产码那行 `slog.Info`（负向，必响）

**手段**：`internal/proc/shutdown.go` 拷贝到**仓外** `D:\tmp\wisp\20261112\425-shutdown\mut\shutdown_mut.go`，把 `:130` 那一行换成一行注释；`go test -overlay` 用 `mut/overlay-mut.json` 把编译输入指过去。**仓内产码一字未动**（§5 的 hash 读数）。

- 锚点命中数：**1**（`[regex]::Matches` 计数，删前先验，⛔ 不许猜）
- 起手 `shutdown.go` sha256 `E35D66F0…EA550A`，变异副本 sha256 `F0D5128B…FA7035`
- 字节数 6815 → 6814（恰好少掉那一行的净长度）

命令（`GOFLAGS` 注入 overlay，仍走同一枚 carrier）：

```
$env:GOFLAGS="-overlay=D:/tmp/wisp/20261112/425-shutdown/mut/overlay-mut.json"
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
```

**读数（红）**：

```
    resident_sink_nail_127_windows_test.go:601: resident sink ...\logs: 1 file(s), 40 record(s)
    resident_sink_nail_127_windows_test.go:675: the D38(e) trail is not in the file the same process closed: 38 record(s) after the install, all msgs [winsec: sealing path resolver installed wisp: persistent log sink installed ... wisp: D38(e) step 10 (exit): ok elapsed=0s]
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (9.41s)
FAIL
```

逐字要点：`the D38(e) trail is not in the file the same process closed`；记录数 45 → **40**（5 条 `shutdown step` 全部消失，与 `[1 2 5 6 7]` 恰好五条对得上）。

⚠ **本发落在判据 1（trail 非空）上，不是落在新加的次序判据上。** 这是必然的：把每一格的记账都删掉，次序也就无从读起。所以另加了 §4 那一发，专门证明**新判据自己**有牙。

## 4. AC#2 突变控制 ②（补强）：把一条 `shutdown step` 记到**回放之后**

只做 §3 那一发，只能证明"非空守卫还在"（那是原断言就有的），证明不了新锚点有牙。所以补这一发：`cmd/wisp/resident_windows.go` 拷到仓外 `mut/resident_windows_mut2.go`，在回放循环**之后**、退出报告之前插一条 `slog.Info("shutdown step booked after the replay", "step", 99)`。**仓内产码一字未动。**

**读数（红，且红在新判据上）**：

```
    resident_sink_nail_127_windows_test.go:701: the last shutdown-step record is record 45 of 46, but the first "wisp: D38(e) step" record is record 35: the D38(e) sequence must book its steps BEFORE the replay runs, in the sink the same process closes (records: [...])
    resident_sink_nail_127_windows_test.go:716: the last shutdown-step record is stamped 2026-11-12 19:48:30.2196796 +0800 CST, after the first replay record's 2026-11-12 19:48:30.2120552 +0800 CST: the booking must happen before the replay, not during or after it
    resident_sink_nail_127_windows_test.go:725: RESIDENT LEG RECORDED ON DISK: 6 shutdown record(s), steps [1 2 5 6 7 99], first="shutdown step skipped (module not present)" last="shutdown step booked after the replay"
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (9.37s)
```

**关键读法**：同一发里 `6 shutdown record(s), steps [1 2 5 6 7 99]` ⇒ **判据 1 绿**（trail 非空、还多了），**判据 2 绿**（步序仍严格递增：99 最大、在最后）；**只有判据 3 红**。位置判据与时间戳判据**同时**响。⇒ 新锚点确有牙，且牙长在它自己身上。

## 5. 恢复读数（绿）＋ 产码零改动证明

撤掉 `GOFLAGS`（overlay 只存在于仓外 JSON 与仓外副本，撤掉即"恢复"）：

```
    resident_sink_nail_127_windows_test.go:601: resident sink ...\logs: 1 file(s), 45 record(s)
    resident_sink_nail_127_windows_test.go:719: RESIDENT LEG ORDERING: 5 shutdown record(s) at records [26 27 32 33 34], first replay record 35 of 45, last trail stamp 2026-11-12T19:48:59.2161246+08:00, first replay stamp 2026-11-12T19:48:59.2175758+08:00
    resident_sink_nail_127_windows_test.go:725: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], ...
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (8.39s)
PASS
```

**产码零改动（我自己现跑，不信自己的手）**：

```
$ git status --porcelain -- internal/proc/shutdown.go cmd/wisp/resident_windows.go
（空）
$ git hash-object internal/proc/shutdown.go        = 745d03039b2c5728a4e18add10fe9790784d7d78
$ git rev-parse HEAD:internal/proc/shutdown.go     = 745d03039b2c5728a4e18add10fe9790784d7d78
$ git hash-object cmd/wisp/resident_windows.go     = b7d54320156cac3bc52802b2b704476bd32cfb44
$ git rev-parse HEAD:cmd/wisp/resident_windows.go  = b7d54320156cac3bc52802b2b704476bd32cfb44
```

## 6. AC#3 越界检查（`git diff --name-only` 逐行）

```
$ git diff --name-only
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/resident_sink_nail_127_windows_test.go
frontend/fixtures/composer-states.html
```

**逐行归因**（本票写面之外的三枚**都不是本腿改的**）：

| 路径 | 起手即脏？ | 归因 |
|---|---|---|
| `.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go` | **是**（本腿起手 `git status --porcelain --untracked-files=no` 就已 ` M`） | 票 291 探针腿的残留，⛔ 非本腿 |
| `.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt` | **是**（同上，起手已 ` M`） | 票 346 派单存档，编排者侧，⛔ 非本腿 |
| `frontend/fixtures/composer-states.html` | **是**（同上，起手已 ` M`） | `frontend/**` 归机主外部 agent，⛔ 非本腿 |
| `cmd/wisp/resident_sink_nail_127_windows_test.go` | 否（本腿改） | **本票唯一写面** |

**收窄到本票的三个禁区面**（更硬的一把尺）：

```
$ git diff --name-only -- cmd/ internal/ docs/
cmd/wisp/resident_sink_nail_127_windows_test.go
```

⇒ `internal/`（含 `internal/proc/shutdown.go`、`internal/observe/thresholds.go`）**零行**、`cmd/` 下**只有本测试文件**。
`docs/PLAN.md`／`docs/specs/**`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 逐面零行。
（`docs/evidence/s1/425-shutdown-record-assert.md` 是本票新增件，走 commit 显式 pathspec 入库，`git diff` 看不见未跟踪文件。）

`cmd/wisp/` 的改动量：`1 file changed, 91 insertions(+), 6 deletions(-)`。

## 7. AC#4 四道门禁读数

| 门禁 | 命令 | 读数 |
|---|---|---|
| `go build ./...` | `go build ./...` | **rc=1**，但红在 **`build/mut_rvw.go`**（`undefined: voiceTtsFiles` 等 10+ 条）——见下方具名说明，⛔ 不是本票 |
| `go vet ./cmd/wisp/` | `go vet ./cmd/wisp/` | **rc=0**，零输出 |
| `gofumpt -l` | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp/resident_sink_nail_127_windows_test.go`（`v0.12.0 (go1.27.1)`） | **零输出**，rc=0（`gofmt -l` 同为零输出） |
| `tools/d22scan` | `sh scripts/d22scan.sh`（经 Git bash） | **rc=0**，逐字 `d22scan: clean - no D22 ban violations`；`examined 480 production Go files`、`ban #8 cmd/ 276 Go files`、自测 `PASS=45 FAIL=0 SKIP=0` |

**`go build ./...` 那发红的具名归因（本票零责任，可逐字复跑）**：

```
$ git ls-files --error-unmatch build/mut_rvw.go
error: pathspec 'build/mut_rvw.go' did not match any file(s) known to git
$ git check-ignore -v build/mut_rvw.go
.gitignore:14:build/	build/mut_rvw.go
$ git cat-file -e HEAD:build/mut_rvw.go
fatal: path 'build/mut_rvw.go' exists on disk, but not in 'HEAD'
$ (git ls-files build/ | Measure-Object).Count   → 0
$ (Get-Item build\mut_rvw.go).LastWriteTime      → 2026-11-04 10:26:45
```

⇒ `build/` 整目录**零枚跟踪件**、`build/mut_rvw.go` 是 2026-11-04 某突变腿留下的**本地未跟踪残留**（被 `.gitignore:14` 吃掉）。**HEAD 导出树上 `go build ./...` rc=0**：

```
$ git archive HEAD -o .../head.tar && tar -x -f .../head.tar -C .../headtree
$ Test-Path .../headtree/build/mut_rvw.go → False
$ (cd headtree && go build ./...) → rc=0
```

工作树侧，排除这一枚残留包后同样 rc=0：

```
$ go build $(go list ./... | ? { $_ -notlike "*wisp/build*" }) → rc=0
$ go build ./cmd/wisp/ ./internal/... → rc=0
```

⇒ **本票的门禁读数是 rc=0**；`go build ./...` 那发 rc=1 由 2026-11-04 的本地残留造成，与票 425 无关，且**不在本票写面内**（`build/` 未被跟踪，删它属别人的动作，⛔ 本腿不动）。

## 8. AC#5 与票 127 的关系

- 本票**只修** `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` **这一枚的断言**。
- ⛔ **不宣称**票 127 的其他格已满足；⛔ **未翻**票 127 与票 425 的任何 AC 复选框（本腿零翻框，翻框归编排者）。
- 票 425 §「本票不解决」那三条照办：`TestAC1ResidentLegInstallsItsLogListenerOnDisk`／`TestAC1ResidentLegOutlivesItsOwnLogFailure` **一字未动**；`resident_windows.go` 的回放顺序 **一字未动**；"把回放挪到 sink 关闭之后" **未做**。

## 9. 复跑清单（逐条可重放）

```powershell
cd "D:\work\workspace\projects plans\Wisp"
# 1) 绿
$env:WISP_CLI_STAGE_DIR="D:\tmp\wisp\20261112\425-shutdown\stage-restore"
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
# 2) 红-①：删掉 shutdown.go:130 的 slog.Info（overlay，仓外副本）
$env:GOFLAGS="-overlay=D:/tmp/wisp/20261112/425-shutdown/mut/overlay-mut.json"
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
# 3) 红-②：回放之后再记一条 shutdown step（overlay，仓外副本）
$env:GOFLAGS="-overlay=D:/tmp/wisp/20261112/425-shutdown/mut/overlay-mut2.json"
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
# 4) 恢复
$env:GOFLAGS=$null
```
