# 140 · AC#2(b) 实现腿交件：ticket224 WISP_ENV pin ＋ R-140-1 恒真改码 ＋ ci.yml dev-branch doctor 断言 step 提案 diff（实现方交件）

- **程**：`worker-ticket140-ac2-fixes`（实现腿；按票面规矩，`AC#2` 的编排者终裁与本案的非实现者对抗验收仍待另程，本文件是实现方证据，不是终裁表）
- **日期**：2026-11-06 15:29–16:1x +08
- **授权**：台账 **A1569**（`docs/reports/pending-and-issues.md:31336`，编排者 2026-11-06 15:xx 登记，逐字）：*"授权登记：140 之 ci.yml「dev-branch doctor 断言 step」——编排者按先例规矩允许登记立此条；本次只出 diff 提案落证据＋做两处纯测试码修复（ticket224 WISP_ENV pin、R-140-1 恒真），ci.yml 本体待定是否摘。"*
- **锚点**：派单锚 `290e08d6f`；导出树基 `b471d4b27`（＝`290e08d6f` + wave32b 收口 `0c852c136` + 199 census）。**写目标三径零漂移现量**：`git diff --stat 290e08d6f b471d4b27 -- cmd/wisp .github/workflows/ci.yml internal/proc` ＝ **空**（现量 16:11）。
- **仪器**：Git Bash（`D:\work\soft\Git\bin\bash.exe -lc`）；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`；**全程未设 `CGO_ENABLED=0`**（现量环境无 CGO 变量，cgo 默认开）；gofumpt **v0.12.0**（盘上 `--version` 现量）；go **1.27.1**。
- **规则 21 隔离**：全部测量在共享 clone 内做——`git clone --shared "…/projects plans/Wisp" /d/tmp/wisp/20261106/140ac2fix/repo`（台账 A1296③ 规程补丁同形）。主共享树的三枚在飞未跟踪测试件（`panel_resident_323_windows_test.go`、`panel_resident_assembly_323_test.go`、`internal/tools/ticket23_image_channel_test.go`）**不在 clone 里**，门禁读数因此只含 HEAD＋本腿两文件改动。

---

## §0 结论速览

1. **两处纯测试码修复落地**（write face＝`cmd/wisp/ticket224_ac3_realprocess_windows_test.go` ＋ `cmd/wisp/dataroot_128_test.go`）：ticket224 的真进程子 env 由"裸 ambient"改为"先栽后钉"（`WISP_ENV=dev` 栽、`WISP_ENV=test` 钉）＋新增一枚文件系统牙；`dataroot_128_test.go:209-210` 的恒真断言改为**从文件系统算期望**（票面 §④ 完成判据逐字兑现，含变异一发）。**生产码零改动**。
2. **变异三发全按预期**：M2（早退目标换另一枚真实目录）⇒ 新形 **RED**；M3（`proc.TestDataDir` 自身漂移）⇒ **旧形 GREEN / 新形 RED**（恒真实证）；M1（摘 pin 行）⇒ ticket224 **RED**（子进程落在栽的 dev 根、写了一次性 APPDATA、真 profile 零接触）。
3. **门禁全绿**：`go build ./cmd/... ./internal/... ./tools/...` rc=0；受影响面 scoped `go test`（两文件全部 9 枚顶层＋子例）全 PASS；`./internal/proc` 全包 ok；`gofumpt -l cmd internal tools scripts` 空；`sh scripts/d22scan.sh` **rc=0 clean**。
4. **ci.yml 本体一字未动**（交付 commit 的 tracked diff 仅两枚 `_test.go`）；提案 step 的**逐字 unified diff**见 §5（在导出树里真插入→`git diff` 现量→逐字还原，故 hunk 行号与上下文对基锚 `b471d4b27` 可直接 `git apply`）。
5. **机主真 profile 本程零写入**（实测凭据：`%APPDATA%\wisp-dev\wisp.db` mtime `2026-11-06 15:14:12`，早于本程开工 15:29，收工后复量未变——读数见 §6⑥）。

---

## §1 修复一：ticket224 子进程 hermeticity（票 140 AC#2 支 (b)）

**现状（HEAD `b471d4b27`）**：`ticket224_ac3_realprocess_windows_test.go:44-48` 起真 `wisp run` 子进程，`cmd.Env = append(os.Environ(), "WISP_TEST_DATA_DIR="+h.dir)`——只给注入、不给 `WISP_ENV`。而注入只在子进程解析 env `test` 时生效（`internal/proc/envfork.go:212`；CLI 数据根同一分叉在 `cmd/wisp/doctor.go:255`）。⇒ UNSET 形（dev host 常态）子进程落 dev、无视注入、起在真 `%APPDATA%\wisp-dev`——本日 09:57 已被 READ-MEASURE 腿实测兑现（`140-ac1-realrun-shaping-r1.md` §4.2.1：写入 `wisp.db` 163840 B ＋ `logs\wisp-20261106-001.jsonl`，用例照绿）。

**改后**（`mint` 内）：

```go
cmd.Env = append(os.Environ(),
    "WISP_ENV=dev", // the planted ambient this pin must beat
    "WISP_ENV=test",
    "WISP_TEST_DATA_DIR="+h.dir,
)
```

- **先栽后钉**（票面 AC#2 判据 (b) "钉环境这件事本身要有牙——先栽后钉才算"；形状＝票 128 `c2fa2e9` 的 `pinEnvThatAsksTheOS128` 推广）：`dev` 栽在先、`test` 钉在后（Windows 环境块同名后者胜，`os/exec` 在 Windows 侧做 case-insensitive 去重保留末次）。⇒ 子进程 env 在任何宿主上都由本用例自己决定，不再被 job 级 env 供养；只钉 `test` 不栽 `dev` 的形状会把"摘掉钉"藏在 runner 的 job 级 env 后面（本票要除的正是那个形状）。
- **牙（文件系统侧）**：`id1 := mint()` 之后、父进程自己的 `memory.Open(h.dir)` **之前**，`os.Stat(h.dir/wisp.db)`——`t101boot` 只栽 config.toml＋DPAPI key fixture（`fixtureKeyRef180` 走 `secret.NewStore`，不建 SQLite 文件；两文件注释现量），所以 h.dir 里出现 `wisp.db` 只能是**子进程自己**在注入根上开的 store（`memory/open.go:40` `dbFileName="wisp.db"`）。摘掉钉 ⇒ 子进程落栽的 dev 根 ⇒ 本断言 RED（M1 实测，§4）。
- 名册不变（无新顶层用例，无改名）。

## §2 修复二：R-140-1 恒真断言改码（A188⑦ 自决落地）

**现状**：`dataroot_128_test.go:202-212` 的期望值是 `proc.TestDataDir()`——而 `doctor.go:256` 早退返回的正是同一枚函数。期望由被测函数自己算出 ⇒ 只挡得住"早退被删"（经 `failConfigDir128` 的 err 支），挡不住"守门人把根认错"：`proc.TestDataDir` 自身漂移时两侧同移、永远绿（票面 §④ 判词逐字形状，票 119 `R-119-9` 族）。

**改后**：期望改为 **从文件系统算**（新 helper `wantTestDataDirFromDisk140`）：

- `WISP_TEST_DATA_DIR` 已声明 ⇒ 期望＝该声明值 verbatim（同一性契约，`proc/envfork.go:100-106`；同形＝winsec `dataroot_symlink_119_other_test.go:347-350` 的 declared-literal 钉法，§5.1 点名过的"修过的形状"）；
- 未声明 ⇒ 期望＝`filepath.EvalSymlinks(os.TempDir())` 拼上 `wisp-test-<os.Getpid()>`——OS 事实（temp 根的内核拼写＋本进程 pid），**不问 proc 的任何代码**（`cleanSpelling119` 的先例：故意用 `EvalSymlinks` 而非 `proc.SealableRoot`，避免把被测函数族再请回期望侧）。
- `failConfigDir128` 保留：早退被删 ⇒ OS 读（已栽坏）被咨询 ⇒ err ⇒ Fatalf——原保护完整保留。
- 断言两侧信息：`dir` vs 独立算出的 `want`，红句写明期望的来源与它能抓的变异形状。

## §3 门禁读数（导出树内现量，全部可逐字重跑）

仪器环境每发一致：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`；未设 CGO_ENABLED。

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools scripts` | **输出空**（收尾复跑二次确认） |
| go build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| 受影响面 go test | `go test ./cmd/wisp -count=1 -timeout 420s -v -run '^(TestTicket224ProductionSessionChangesAcrossRealOSProcesses\|TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128\|TestAC2TestDataDirBranchStillResolves128\|TestAC2RefusalMarkersAreNotAShortenableList128\|TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128\|TestAC1TheStartUpDirectoryCheckHasItsOwnTeeth128\|TestAC3TheDoctorLegLeavesTheStartUpDirectoryUntouched128\|TestAC2ARefusalWithATreeInTheStartUpDirectoryIsRed128\|TestAC5RefusalSentenceHasOneWordingPoint128)$'` | **顶层 9/9 全 PASS**（含 `TestTicket224ProductionSessionChangesAcrossRealOSProcesses` 11.27s、`ok …cmd/wisp 11.952s`；128 族子例逐枚 PASS） |
| 第二受影响包 | `go test ./internal/proc -count=1 -timeout 420s` | **ok 2.013s** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；examined 421 production Go files；positive control `runtests.sh -C tools/d22scan ./...` PASS=35 FAIL=0 SKIP=0 |

**关于派单门 "go test the two affected packages -timeout 420s" 的口径说明**：两处修复都在 `./cmd/wisp` 一包（第二包按判据依赖取 `./internal/proc`，已全包跑）。全包 cmd/wisp 本机今日实测 **580.7s–653.2s**（READ-MEASURE 腿 §4.2/§4.2.3，`-timeout=25m` 仪器），420s 上限装不下整包——按 420s 跑整包只会复刻已登记的 deadline panic（机器形状，非本腿引入）。故 420s 形状落在受影响面 scoped run 上（两文件全部顶层用例＋子例，逐枚 PASS），整包基线引 READ-MEASURE 腿今日读数（UNSET 489/443/44/2、TEST 489/442/45/2，红全是 DLL 遮蔽/模型未布/4 枚具名移交红，零枚涉 WISP_ENV/数据根）。**本腿改动零生产码，整包色不变可由 scoped 读数＋构建门＋名册不变推出。** 若编排者仍要一发整包 `25m` 复跑，仪器与基线都在，一条命令可复现。

## §4 变异自证（全部在导出树内做，做完逐字还原；`git diff` 复量非变异文件零残留）

| 变异 | 内容 | 预期 | 实测 |
|---|---|---|---|
| **M2**（票面 §④ 要求的那发） | `doctor.go:256` 早退目标 `proc.TestDataDir()` → `os.TempDir()`（另一枚真实目录） | 新形必须红 | **FAIL 0.103s**：`test data dir = "C:\Users\swq", want "C:\Users\swq\wisp-test-13640"`（红句自带两侧值与期望来源） |
| **M3**（恒真实证） | `internal/proc/envfork.go:125` 回落根拼名漂移为 `wisp-probe140` | 旧形绿、新形红 | **旧形探针 PASS / 新形 FAIL**：新形 `test data dir = "C:\Users\swq\wisp-probe140", want "C:\Users\swq\wisp-test-22400"`。旧形（`r1401_oldform_probe_test.go`，pre-fix 形逐字保存、仅存在于导出树、永不提交）在同一变异下 GREEN ⇒ **"两侧同移、恒真"由测量兑现，不是论证** |
| **M1**（pin 的牙） | 摘掉钉环境里 `WISP_ENV=test` 那行 | ticket224 必须红 | **FAIL 7.73s**：子进程落**栽的 dev 根**（一次性 `APPDATA=D:\tmp\wisp\20261106\140ac2fix\mut-appdata`），在那里新建 `wisp-dev\config.toml`＋`logs\wisp-20261106-001.jsonl`、找不到 provider 配置、无 SESSION-MINT 行 ⇒ 用例红。**机主真 profile 零接触**（APPDATA 全程指一次性根；真 profile `wisp.db` mtime 15:14:12 不变，§6⑥） |
| 复钉复跑 | M1 还原后同一命令、同一重定向 APPDATA | 绿 | **PASS 11.18s**——钉住的子进程对 APPDATA 完全无感（test 支早退，不问 OS），hermetic 正面读数 |

（M3 用的探针件 `cmd/wisp/r1401_oldform_probe_test.go` 按"临时件只建不删"留在导出树，登记于 §6⑦；它不在任何交付 commit 的 pathspec 里。）

## §5 ci.yml「dev-branch doctor 断言 step」提案（PROPOSAL ONLY——ci.yml 本体本程未动）

生成方式：导出树内对 HEAD `b471d4b27` 真插入该 step → `git diff --no-color --unified=3 -- .github/workflows/ci.yml` 现量捕获 → `git show HEAD:.github/workflows/ci.yml` 逐字还原 → `git diff --stat` 复量 ci.yml 零内容差。**hunk 行号与上下文锚 `b471d4b27`（对派单锚 `290e08d6f` 该文件零漂移）**；HEAD 再前进时需重锚。验收形状＝READ-MEASURE 腿 §4.3 探针（`WISP_ENV=dev`＋一次性 APPDATA，断言 data dir 行落 `$APPDATA\wisp-dev`）。

```diff
diff --git a/.github/workflows/ci.yml b/.github/workflows/ci.yml
index 58179b030..3e6c48640 100644
--- a/.github/workflows/ci.yml
+++ b/.github/workflows/ci.yml
@@ -883,6 +883,49 @@ jobs:
         if: ${{ !cancelled() }}
         run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev
 
+      # PROPOSAL ONLY (ticket 140, registered ledger A1569) - NOT applied by the
+      # ticket 140 evidence leg; ci.yml itself is untouched by that leg. Inserted
+      # into the export tree solely to produce the exact unified diff recorded in
+      # docs/evidence/s1/140-ac2-fixes.md §5, then reverted.
+      #
+      # WHY THIS STEP (the shape the ticket's read-measure leg measured): the three
+      # build.ps1 doctor steps (this one, slo-smoke, slo-full) all run under the
+      # job-level WISP_ENV=test, so on every runner the userConfigDir read
+      # (cmd/wisp/doctor.go:258) is never consulted: the dev branch of the
+      # data-root fork has had no CI verdict at all, and SPEC-11:182-183 /
+      # SPEC-03:101's "CI asserts WISP_ENV actually takes effect" duty was served
+      # by no case (docs/evidence/s1/140-ac1-realrun-shaping-r1.md §2.3/§5.2①).
+      # This step buys that verdict WITHOUT unpinning anything: it runs the
+      # build\wisp.exe this job just produced with WISP_ENV=dev and APPDATA
+      # pointed at a one-shot directory, then asserts the data-dir line lands in
+      # <one-shot>\wisp-dev - the §4.3 probe is the acceptance shape, and nothing
+      # below pins, deletes or weakens any existing step.
+      #
+      # THE ASSERTION IS THE LINE, NOT DOCTOR'S RC: doctor's overall rc belongs to
+      # the build smoke step above. A data-dir line naming anything other than
+      # <one-shot>\wisp-dev means the OS read was skipped (test-branch early exit)
+      # or answered elsewhere, and the step fails loudly with both values.
+      - name: "Doctor dev-branch assertion (WISP_ENV=dev asks the OS; one-shot APPDATA)"
+        if: ${{ !cancelled() }}
+        shell: pwsh
+        run: |
+          $oneShot = Join-Path $env:RUNNER_TEMP "wisp-devfork-appdata"
+          New-Item -ItemType Directory -Force -Path $oneShot | Out-Null
+          $env:WISP_ENV = 'dev'
+          $env:APPDATA = $oneShot
+          $out = & (Join-Path (Get-Location) 'build\wisp.exe') doctor 2>&1 | Out-String
+          $out | Set-Content -Path (Join-Path $oneShot 'doctor-devbranch.out')
+          $want = Join-Path $oneShot 'wisp-dev'
+          Write-Host $out
+          if ($out -notmatch [regex]::Escape('data dir writable (dev)')) {
+            Write-Error "dev-branch doctor assertion: doctor printed no 'data dir writable (dev)' line; the data-root fork's dev branch produced no verdict"
+            exit 1
+          }
+          if ($out -notlike "*$want*") {
+            Write-Error "dev-branch doctor assertion: the data-dir line did not land in '$want' - the OS read (userConfigDir, cmd/wisp/doctor.go:258) was not consulted or answered elsewhere, which is the job-level WISP_ENV=test dependency ticket 140 registers"
+            exit 1
+          }
+
       - name: "cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"
         # AC#4 asked whether this package can be tested on CI at all, and forbade
         # answering from the laptop. It was measured on both legs - see
```

设计要点：断言**只咬 data dir 行**，不咬 doctor 总 rc（rc 归 build smoke step 管——READ-MEASURE §4.3 注记过无 ldflags 版本 pin 的 build 会让 doctor 总判 FAIL，与数据根无关）；`RUNNER_TEMP` 一次性根保证 slo-full（self-hosted）同形也零真 profile 写入；step 名与所测严格对齐（R-140 家族"名实不符"就是本票 (c) 的账①②）。

## §6 纪律自证

① **仓库改动**：交付面＝两枚 `_test.go`（§1/§2）＋本文件＋票面 log 一条；**`.github/workflows/ci.yml` 本体零改动**（主树 tracked diff 现量仅两枚测试文件，89 insertions / 6 deletions）；**`docs/reports/pending-and-issues.md` 一字未碰**；生产码（doctor.go／envfork.go）的变异只在导出树内发生并已逐字还原（`git diff` 复量零残留）。② commit 全程显式 pathspec、无 `add -A`／amend／reset／rebase／stash／`checkout .`／clean；**不 push**。③ 全程未设 `CGO_ENABLED=0`（每发环境现量无 CGO 变量）。④ cmd/wisp 测试串行单发（"重包同一时刻只跑一个"），未与其它腿并行同包。⑤ 未删任何件（导出树内探针件、一次性 APPDATA、diff 捕获件全部保留）。⑥ **机主真 profile**：`C:\Users\swq\AppData\Roaming\wisp-dev\wisp.db` mtime `2026-11-06 15:14:12`（早于本程开工 15:29），收工复量不变 ⇒ 本程零真 profile 写入；16:02 的 128 族 scoped run 与 16:06 的复钉 run 子进程全落 h.dir/一次性根。⑦ 导出树外临时件（**只建不删**，全部登记）：`/d/tmp/wisp/20261106/140ac2fix/{repo(clone，含拷入的 build\*.dll 与 third_party\sherpa-onnx——两目录仓内不受 git 管，clone 需自备)、ci-yml-proposal.diff、mut-appdata、repo\cmd\wisp\r1401_oldform_probe_test.go}`。⑧ 判为注入的工具输出：0；真通知回显：无后台 job（全部前台串行）。

## §7 交回编排者

- `AC#2` 的编排者终裁仍待（本程按 A1569 登记的"两处纯测试码修复＋diff 提案落证据"执行完毕）；ci.yml step 是否落地由编排者拍（diff 现成可 apply，锚 `b471d4b27`）。
- 非实现者对抗验收可复跑的尺：§3 表逐字命令（导出树或主树均可）；§4 变异三发逐字可复现。
