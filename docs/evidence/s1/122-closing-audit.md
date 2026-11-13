# 票 122 收案审计（非实现者验收）— 5 个未勾 AC 逐条现跑判定

- 审计腿：122-closing-audit（审计/验收腿，⛔ 非实现者）
- 时刻：2026-11-13 21:1x–22:0x +08
- 开工 HEAD：`10f98f52b0b19050a682373b9bd14b914d0d5a6f`（工作树 tracked `.go` 与 HEAD 一致；仅 `.scratch` 未跟踪件与 1 枚 html 件有他腿在飞改动，均不在 staticcheck/测试分析面）
- 仪器版本（逐字，AC#7 格式要求）：`staticcheck.exe 2026.2.2 (0.8.2)`、`gofumpt v0.12.0 (go1.27.1)`（宿主 `D:\work\base\gopath\bin\`）、`go version go1.27.1 windows/amd64`
- 写面：本表（新）＋票 122 面（勾框＋Progress log＋`git mv` 改名）。⛔ 零代码/测试/脚本/CI/manifest 改动。
- 临时件：`D:\tmp\wisp\20261113\122-close\`（只建不删）

> 判据说明：本表每一格都现跑，⛔ 不转抄任何前腿读数。CI 侧读数由本腿用 `gh` 现取（run/job/step 逐字）。

---

## AC#1 — 当前 HEAD 纯净快照重测条数与分布（两 GOOS），与预做 34/78 并列

**票面要求（复述）**：开工先在当前 HEAD 的纯净快照里重测条数与分布（两 GOOS 都要），把「预做时的 34/78」与今天的数并列贴出；不一致就说明差在哪几枚 commit（不许直接覆盖旧数）。

### 现跑读数

**linux 口径**（`GOOS=linux`，**未设 `CGO_ENABLED=0`**；`-checks=all,-U1000,-SA9009`，逐字照 `ci.yml:530`）：

```
=== MODULE: . ===
-: build constraints exclude all Go files in ...sherpa-onnx-go-linux@v1.13.8 (compile)
--- module . [pass1 GOOS=linux checks=all,-U1000,-SA9009]: exit=1 packages=40 findings=0 toolchain-crash-lines=0 compile-lines=1
=== MODULE: tools/d22scan ===
--- module tools/d22scan [...]: exit=0 packages=1 findings=0 toolchain-crash-lines=0 compile-lines=0
=== MODULE: tools/mockllm ===
--- module tools/mockllm [...]: exit=0 packages=1 findings=0 toolchain-crash-lines=0 compile-lines=0
PASS1 TOTAL: packages=42 findings=0 toolchain-crash-lines=0 compile-lines=1 overall_exit=1
```

- 本机 Windows 主机**没有 linux cgo 工具链**，`go list` 列得出 40 包（含 `cmd/wisp`），但 staticcheck 加载 `cmd/wisp` 时撞 `sherpa-onnx-go-linux` 构建约束 ⇒ 那条 compile 行是**主机形状产物**（A1851 同族），`cmd/wisp` 未被分析。
- **CI 的 linux 读数才是权威**（ubuntu 有 cgo）：本腿 `gh` 现取 run `43431925841` / lint job `131907281319` / step 13 `staticcheck` 的自报行逐字：
  ```
  --- module . [pass 1 GOOS=linux checks=all,-U1000,-SA9009]: exit=0 packages=40 findings=0 toolchain-crash-lines=0 compile-lines=0
  --- module tools/d22scan [...]: exit=0 packages=1 findings=0 ...
  --- module tools/mockllm [...]: exit=0 packages=1 findings=0 ...
  staticcheck self-report: version=staticcheck 2026.2.2 (0.8.2) modules=3 [pass1 ...] packages=42 findings=0 toolchain-crash-lines=0 compile-lines=0 ... (step exit is 0)
  ```

**windows 口径**（`GOOS=windows CGO_ENABLED=1 -tags=winlive`，`-checks=all,-SA9009`，`./build` 按名排除，逐字照 `ci.yml:1067-1085`）：

```
==================== MODULE: . ====================
--- module . [pass2 GOOS=windows tags=winlive checks=all,-SA9009]: exit=0 packages=41 findings=0 toolchain-crash-lines=0 compile-lines=0
==================== MODULE: tools/d22scan ====================
--- module tools/d22scan [...]: exit=0 packages=1 findings=0 toolchain-crash-lines=0 compile-lines=0
==================== MODULE: tools/mockllm ====================
--- module tools/mockllm [...]: exit=0 packages=1 findings=0 toolchain-crash-lines=0 compile-lines=0
staticcheck windows self-report: version=staticcheck.exe 2026.2.2 (0.8.2) modules=3 [pass2 ...] packages=43 findings=0 toolchain-crash-lines=0 compile-lines=0 (step exit is 0)
```

- CI 侧同一步（test-windows job `131907281121` step 8 `staticcheck (windows, cgo, tags=winlive; ticket 122 A1851)`）= **success**。
- **纯净快照复跑**（`git archive HEAD` 导出到 `D:\tmp\wisp\20261113\122-close\snap`，同口径）＝**43 包 0 finding 0 compile exit 0**，与工作树读数一致。
  ⚠ 工作树在本腿跑读期间有**他腿在飞改动** 4 枚 `cmd/wisp/*_test.go`（票 111 波 10：`config_reload_223_test.go`／`panel_host_windows_test.go`／`resident_voice_41a2r1_windows_test.go`／`voice_e2e_demo_16r1_windows_test.go`，+124/−15）——本腿**未触碰、未提交**，且静态check 读数在含它们与不含它们两种形状下**均为 0**；判据以纯净快照 + CI 为准。

### 与预做 34/78 的并列

| 读数 | GOOS=linux | GOOS=windows |
|---|---|---|
| 票面预做（audit-85-preflight，旧 HEAD；标题里的 34/78） | 34 | 78 |
| CI 自报分母（A101②，09-22） | — | 42 |
| 批 1 收口（122sa1，HEAD `6e2eb5845`，11-08） | 20 | 115 |
| 批 2 收口（122sb2，HEAD `31482f271`，11-08） | 14 | 98 |
| 122-r1 收口（HEAD `e31937a90`，11-11） | 57 | 103 |
| **今天 HEAD `10f98f52b` 实测（本腿现跑）** | **0**（CI 权威 42 包 0 finding；本机 39 包 0 finding＋1 主机 compile 行） | **0**（43 包 0 finding） |

**差在哪几枚 commit（不覆盖旧数，逐块归因）**：
- `SA4000` 那 1 条已绝迹：`79041f071`（ci-r1，11-01）删掉重复操作数。
- 第一批 `SA1019`/`SA4006`：`4df9964d1`（batch 1）＋`dc8a2d536`（risk 解冻后修 SA1019）。
- `U1000` 主体：波 2–7（`04ae7d817`→`e505420b2`）按 A1773/A1830 的 GOOS 口径分遍清。
- 口径本身：`61ac182e8` 把 staticcheck 步拆两遍（linux 不含 U1000／windows 含 U1000），分母从「42」变成「pass1 42 包 / pass2 44 包」两段。
- 逐 commit 全量归因不再可行（其间几十枚 commit）；票面点名的三组有一对一证据（见 AC#2），其余块归因见 `docs/evidence/s1/122-r1-staticcheck-backlog.md` §AC#1。

**判定：满足。** 当前 HEAD 两 GOOS 均 0 finding（CI linux 权威读数 42 包 0/0/0；本机 windows 43 包 0/0/0），34/78 与今天的数已并列，差块归因到具体 commit。

---

## AC#2 — 第一批只清三组（SA4000 1＋SA1019 3＋SA4006 1），SA4000 要有行为证据

**票面要求（复述）**：第一批只清三组（`SA4000` 1 条＋`SA1019` 3 条＋`SA4006` 1 条）；`SA4000` 那格要有行为证据——修完之后原来那枚用例仍然红于「它本来该红的形状」（不是变成恒真，而是真的验到了两个维度）。

### 现跑读数（三组剩余条数）

```
=== windows pass, checks=SA4000,SA1019,SA4006 (winlive, build excluded) ===
EXIT=0                       ← 0 条
=== linux pass, checks=SA4000,SA1019,SA4006 (internal+cmd) ===
-: build constraints exclude all Go files in ...sherpa-onnx-go-linux@v1.13.8 (compile)
EXIT=1                       ← 0 条 finding（唯一输出是主机 compile 行）
```

⇒ 三组在两 GOOS 下**均剩 0 条**。

### 那 5 条（票面锚）的处置与 `文件:行`

| # | 形状 | 位置 | 处置 | 落点证据 |
|---|---|---|---|---|
| 1 | SA4000 | `internal/observe/goroutine_test.go:151` | **真修**（删重复操作数） | `79041f071` diff 逐字 `if !errors.As(err, &oe) \|\| !errors.As(err, &oe) {` → `if !errors.As(err, &oe) {` |
| 2 | SA1019 | `cmd/wisp/secret_argv_windows_test.go:174` | **真修**（改读 `os.Getenv("GOROOT")`） | `4df9964d1`；现读 `:166-170` 注释逐字说明弃用与兜底职责 |
| 3 | SA1019 | `internal/speech/cgo_boundary_windows_test.go:39`、`:95` | **真修**（改 `windows.CurrentProcess()`，删恒不触发的 err 分支） | `4df9964d1`；现读 `:36-38` 注释逐字 |
| 4 | SA1019 | `internal/risk/proccpu_windows_test.go:20` | **真修**（改 `windows.CurrentProcess()`） | `dc8a2d536`（A1757 具名解冻后落地）；现读 `:20` |
| 5 | SA4006 | `cmd/wisp/resident_watchdog_windows_test.go:269`、`:305` | **真修**（补 `len(hold) != N` 尺寸断言＝加强） | `4df9964d1`；现读 `:268-272`、`:304-308` |
| 6 | SA4006 | `cmd/wisp/slo_windows.go:803` | **真修**（初值改裸零值 `var last subjectReportRead`） | `4df9964d1`；现读 `:800-803` 注释逐字 |

> 票面写「SA1019 3 条」，实测第一批是 **4 条**（批 1 表 #2–#5，含 risk 那条）；另有第 5 条 `cmd/wisp/resident_approval_live_246_windows_test.go:420` 的 `runtime.GOROOT` 由波 7 落**具名 `//lint:ignore SA1019` 带理由**（`e505420b2`，理由逐字：「这个分支只在 `LookPath("go")` 失败后才跑到，`go env GOROOT` 是循环依赖；本机夹具与构建同机，build root 即正确 root」）。全仓 `git grep -n "lint:ignore SA1019"` 仅此 1 枚。

### SA4000 的行为证据（本腿现跑变异）

用 `go test -overlay`（⛔ 不改树）把生产件 `internal/observe/goroutine.go:317` 的 `err = &Error{...}` 换成普通 `fmt.Errorf`，跑原用例：

```
=== RUN   TestPanicInWorkerSurvivesAndCancelsRoot
    goroutine_test.go:152: panic error not an observe *Error: mutated panic in goroutine "tool-exec-boom" (owner test): boom: injected test panic
--- FAIL: TestPanicInWorkerSurvivesAndCancelsRoot (0.00s)
FAIL
```

⇒ 类型维度（`:151`）**真的会红**；同时类别维度（`:147` `ClassOfOrInternal(err) == ClassInternal`）仍成立（未分类错误按设计落 ClassInternal）。原用例在未变异时 `:147`/`:151` 是两个独立断言，**不是恒真**。变异件只存在于 `D:\tmp\wisp\20261113\122-close\goroutine_mut.go`，主树 `git status` 零改动。

**判定：满足。** 三组两 GOOS 剩 0 条；5＋1 条处置逐条具名（真修 6 条／具名 ignore 1 条，均带理由）；SA4000 行为证据现跑成立。

---

## AC#5 — 门禁：受影响包 `-count=2 -v` 四数＋gofmt/gofumpt/vet/d22scan

**票面要求（复述）**：受影响包 `-count=2 -v` 四数逐条点名（`-count=2` 不缓存；`=== RUN` 行数 == 不同测试名 × 2；非 `-v` 既不印 PASS 也不印 SKIP）；`gofmt -l`＋`gofumpt -l` 真跑；`go vet`；`sh scripts/d22scan.sh` 纯净快照 rc=0、台账各 scope 不降。

### 现跑读数（四数，逐包）

| 包 | rc | `=== RUN` 行 | 不同测试名（全名） | ×2 对上？ | PASS | FAIL | SKIP | 包级结论 |
|---|---|---|---|---|---|---|---|---|
| `internal/observe` | 0 | 238 | 119 | **是**（119×2=238） | 238 | 0 | 0 | `ok` |
| `internal/risk` | 0 | 578 | 289 | **是**（289×2=578） | 576 | 0 | 2 | `ok` |
| `internal/speech` | 0 | 212 | 106 | **是**（106×2=212） | 118 | 0 | 94 | `ok` |
| `cmd/wisp` | 1 | 1594 | 796 | **是**（796×2=1592；差 2 行＝`Test384ResidentWatchdogEndToEnd` 每遍多出一行嵌套 `=== RUN`，非 count=2 失效） | 1505 | 83 | 6 | `FAIL` |

- 四包的 `ok`/`FAIL` 行均无 `(cached)` 字样（observe 10.768s／risk 28.358s／speech 0.502s／cmd/wisp 2040.862s），`-count=2` 真跑未命中缓存。
- `internal/risk` 的 2 条 SKIP 与 `internal/speech` 的 94 条 SKIP 均为既有环境门控（risk 的 `TestSyncRegistryProbeLive` 族；speech 的硬件/夹具门控），非本票引入。
- **`cmd/wisp` 的 83 条 FAIL（42 个不同名）逐条归因，零条属票 122**：
  - 绝大多数是**语音链环境红**——仪器自报逐字 `speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH`（本机 loader 顺序让一枚 1.17.1 遮蔽钉版 1.28），涉及 `Test15e*`／`Test16r1*`／`Test26P*`／`Test26R3*`／`Test381r2*`／`Test382WR2*`／`Test390r1*`／`Test392*`／`Test394R1*`／`Test395r1*`／`Test41a2r1*`／`TestWireR1*`／`Test12r2SettleWalk*`；
  - `TestAC4FocusReturnToPriorWindowGap33r5`＝真窗前台时序（机器耦合，A1848 已登记同形）；
  - `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`＝票 226/223 的既有红（`always_write_no_clobber_226_test.go:119` 逐字「reading a hand edit back is ticket 223」），非本票。
  - **本票改到的面两遍全绿**：`Test384*` PASS=8／FAIL=0、`TestSecret*` PASS=94／FAIL=0、`TestSLO*` PASS=82／FAIL=0（watchdog 尺寸断言、secret GOROOT 兜底、slo 裸零值三处修点的用例逐条 ×2 PASS）。

### 其余门禁（现跑）

| 门禁 | 读数 |
|---|---|
| `gofmt -l internal cmd tools scripts frontend`（**纯净快照** `git archive HEAD`） | **0 行**，rc=0 |
| `gofumpt -l internal cmd tools scripts frontend`（v0.12.0，**纯净快照**） | **0 行**，rc=0 |
| `go vet ./internal/observe/ ./internal/risk/ ./internal/speech/ ./cmd/wisp/` | **rc=0，零输出** |
| `sh scripts/d22scan.sh`（工作树） | **rc=0，clean**；scope：`internal/=403, cmd/=77, frontend/=123, internal/tools/=51, internal/risk/=21, design/=39, ban#8 internal/=1005, ban#8 cmd/=277`（对 122-r1 的 402/74/122/51/21/39/1003/270 **逐格不降**） |

> 工作树 `gofmt -l`/`gofumpt -l` 会点名 `internal/models/acl_sid_121_test.go`——该件 **HEAD blob 为 LF（CR=0）**，工作副本被别的进程动成 CRLF（CR=268）；纯净快照上两把尺均 0 行，CI  checkout 亦是 LF ⇒ 读数取纯净快照。

**判定：满足。** 四包 `-count=2 -v` 四数逐条点名且 `=== RUN` == 不同测试名 ×2（cmd/wisp 差 2 行已归因到嵌套 `=== RUN`）；`cmd/wisp` 的 83 条红逐条归因为零条属本票（语音链环境红＋真窗时序＋票 226 既有红），本票改到的 `Test384`/`TestSecret`/`TestSLO` 两遍全绿；gofmt/gofumpt 纯净快照 0 行、vet rc0、d22scan clean rc0 且八 scope 不降。

---

## AC#6 — 结案判据＝步级 CI 读数

**票面要求（复述）**：结案判据是步级 CI 读数：`lint` job 里 staticcheck 那一步的 run id＋job id＋step 号＋结论；不许用本地 rc=0 冒充 CI 绿。

### 现跑读数（本腿 `gh` 现取，非转抄）

**当前 HEAD `10f98f52b` 的 run（比派单给的 285a5c310 更新，6 枚 commit 之后）**：

| 项 | 值 |
|---|---|
| run id | **43431925841**（headSha `10f98f52b0b19050a682373b9bd14b914d0d5a6f`） |
| job | **lint**（job id **131907281319**），conclusion **success**，16 步全绿 |
| step | **13 `staticcheck`**，conclusion **success** |
| 自报行 | `pass1 GOOS=linux ... packages=42 findings=0 toolchain-crash-lines=0 compile-lines=0 (step exit is 0)` |
| 另侧 | test-windows job `131907281121` step 8 `staticcheck (windows, cgo, tags=winlive; ticket 122 A1851)` = **success** |

派单所给读数（run `43425302222` / job `131885140855` / step 13 / success）在 `285a5c310` 上同样成立，但被 HEAD 这枚**更新更全**的读数覆盖（HEAD 这枚含 6 枚后续 commit，含唯一 .go 改动 `cmd/wisp/config_reload_perm_223_windows_test.go`）。

**判定：满足。** 步级 CI 读数在**当前 HEAD** 上成立：lint job 16 步全绿，staticcheck 步 success、自报 42 包 0 finding 0 compile、exit 0；windows 侧同一步亦 success。

---

## AC#7 — 工具版本未钉的三件后果（读全票面后自行判定）

**票面要求（复述）**：lint 这一族工具版本没钉，跑「CI 同形」门禁会升宿主工具链。三件后果各给判据：(a) 同一纯净快照用两版 gofumpt 各跑一次、差集逐名列出；(b) 门禁读数一律带 `--version`，或别再逐字复制 `go install …@latest`；(c) 报「lint 还剩几枚」必须同时报 staticcheck 版本。修法范围只许动 `ci.yml` 版本字面量与读数格式。

### 现跑读数

**(a) 双版本对齐（同一棵纯净快照 `git archive HEAD`，CI scope `internal cmd tools scripts frontend`）**：

| 版本 | 输出 |
|---|---|
| `v0.12.0`（CI 钉版，宿主 `D:\work\base\gopath\bin\gofumpt.exe`） | **0 行** |
| `v0.7.0`（本腿装到仓外 `D:\tmp\wisp\20261113\122-close\gobin-v070`，`--version` 逐字 `v0.7.0 (go1.27.1)`） | **2 行**：`cmd\wisp\resident_voice_engines_windows.go`、`cmd\wisp\subagent_selfapproval_197_test.go` |

⇒ **差集逐名＝那 2 枚**（v0.7.0 报、v0.12.0 不报）。这坐实了 AC#7 的判断：跨批「0 行」在两版下**不可比**，钉版不是洁癖。

**(b) 读数带版本 / 不复制 `@latest`**：
- `ci.yml:307` 逐字 `go install mvdan.cc/gofumpt@v0.12.0` —— **已钉版，不再是 `@latest`**（AC#7 允许的「两者择一」取后者）。
- `ci.yml:511`／`:1055` 逐字 `go install honnef.co/go/tools/cmd/staticcheck@2026.2.2` —— 已钉版。
- staticcheck 步自报行**自带版本**：`staticcheck self-report: version=staticcheck 2026.2.2 (0.8.2) ... findings=...`（CI 日志逐字，见 AC#6）。

**(c) staticcheck 版本与枚数同报**：自报行把 `version=` 与 `findings=` 打在**同一行**（`version=staticcheck 2026.2.2 (0.8.2) ... packages=42 findings=0 ...`），满足「报枚数必须同时报版本」。旧 42 枚主账在 2026.2.2 口径下已清到 0，本腿 AC#1 的读数同报版本。

**判定：满足。** (a) 双版本差集已逐名落定（2 枚）；(b) 版本已钉、读数自带版本；(c) 版本与枚数同行。修法只动 `ci.yml` 版本字面量与读数格式，未碰 D22 四禁令。

---

## 结案动作

- **5 格全判「满足」** ⇒ 勾框（`- [ ]` → `- [x]`，AC 文字未改）＋ `git mv` 加 `-done` 后缀。
- 新文件名：`122-clear-the-staticcheck-backlog-after-85a-done.md`
- commit：显式 pathspec，⛔ push。

## 越界自证

- 写面：`docs/evidence/s1/122-closing-audit.md`（本表，新）＋ `.scratch/wisp/issues/122-clear-the-staticcheck-backlog-after-85a.md`（勾框＋Progress log＋改名）。
- ⛔ 零 `.go`／测试／`scripts/**`／`internal/**`／`cmd/**`．`.github/**`／`models/**` 改动；⛔ 未加任何 `//lint:ignore`；⛔ 未改任何 checker 开关/阈值/断言/golden。
- 临时件全部落 `D:\tmp\wisp\20261113\122-close\`（只建不删）。
