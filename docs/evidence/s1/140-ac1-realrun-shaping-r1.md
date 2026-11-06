# 140 · AC#1 真跑取证 ＋ AC#2 裁形清单（READ-MEASURE 腿，实现方交件）

- **程**：`worker-ticket140-readmeasure`（实现腿；按票面规矩，**非实现者终裁表仍待另程出**，本文件是测量＋裁形清单，不是 AC#1/AC#2 的终裁）
- **日期**：2026-11-06 09:36–10:2x +08
- **码面锚点**：`84b0a7c4dfcf2b28888fb642df759af0fe01b1b0`（开工 `git rev-parse` 现量，与派单一致；收工未变，见 §9）
- **仪器**：Git Bash（`D:\work\soft\Git\bin\bash.exe -lc`）；Go 构建面 `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`，CGO 默认开（全程未设 `CGO_ENABLED=0`）；gofumpt **v0.12.0**（盘上现量 `--version`，票面 `AC#4` 要求的口径）
- **本程性质**：READ-MEASURE。**零行代码/零行 ci.yml 改动**；唯一交付物＝本文件（＋票面 log 一条）。AC#3 变异、恒真改码（`A188⑦`）均**不在本程**。

---

## §0 结论速览

1. **四包两形真跑（`WISP_ENV` 未设 vs `=test`）名册差集＝0，可归因于 env 的换色用例＝0**（`internal/buildinfo` 2/2/0/0、`internal/secret` 25/25/0/0、`internal/winsec` 76/76/0/0 全两形一致；`cmd/wisp` 完整体 489 枚两形、唯一一枚差名经 `-count=3` 两形专项复跑判为计时偶发，§4.2.3）。**HEAD 84b0a7c4d 上，票面标题描述的"加固腿从不被咨询"形状在测试层已经不存在**——票 128 `c2fa2e9` 的 pin helper 落地后，四包里每一枚碰数据根的用例要么显式 pin、要么显式传 `dataDir`、要么 rebind `userConfigDir`。
2. **"换分支但不换色"那一族今天活在生产面，不在测试面**：`wisp doctor`／`wisp secret`／`wisp run` 的数据根都读 ambient `WISP_ENV`（`doctor.go:246→255`、`secret.go:184`、`run.go:218→2131`），同一枚二进制同一枚命令两形各一发实测**数据根翻面**（§4.3）。而 CI 上跑到这三种生产面的地方（`build.ps1:363-367` 的 doctor smoke，落在 `test-windows:884`／`slo-smoke:1149`／`slo-full:1225` 三枚 step）全部在 job 级 `WISP_ENV=test` 之下 ⇒ **在 runner 上 `userConfigDir` 那一行没有任何 CI 分母咨询它**（M2/M3 现量，§4.3/§4.4）。
3. **测试层仍有一枚"hermeticity 靠 ambient 供养"的隐藏依赖，而且本程实测它咬了真 profile 一口**：`ticket224_ac3_realprocess_windows_test.go:44-48` 起真进程只给 `WISP_TEST_DATA_DIR`、不给 `WISP_ENV`，而注入只在 `env==test` 时生效（`internal/proc/envfork.go:212`）⇒ **UNSET 形（dev host 常态）下该子进程无视注入、起在真 `%APPDATA%\wisp-dev` 上：09:57:22 实测写入 `wisp.db`＋`logs/wisp-20261106-001.jsonl` 两枚真 profile 文件、用例照样 PASS**（归因链见 §4.2.1）。TEST 形下同一用例 hermetic（子进程继承 test）。这是 (a′) 支最硬的一枚实测代价（§7 表）。
4. **裁形推荐（AC#2，待编排者终裁）**：**(b) 为主 ＋ (c) 收窄成三笔具名登记 ＋ (a′) 的价值改用"加一枚显式 dev-分支断言 step"来买，而不是"摘 env"**。凡动 `ci.yml` 一律**先登记**（§8：本程未改、也未获登记，按派单令 STOP 在这一格前面）。
5. **GATES 全绿**：`go build ./cmd/... ./internal/... ./tools/...` rc=0；`gofumpt -l internal cmd tools scripts frontend` 输出**空**（v0.12.0）。零代码改动 ⇒ 与"无代码变更"预期一致。

---

## §1 锚点漂移（先说清，免得下游按旧行号找）

票面 `:9` 写的四枚 job 级 `WISP_ENV` 在 `:227/:337/:482/:540`。HEAD 84b0a7c4d 现量（`grep -n 'WISP_ENV' .github/workflows/ci.yml`）：

| 今天行号 | 归属 job | runner | job 里有没有 `go test` |
|---|---|---|---|
| `:643` | `test-core`（`:634` 起） | `ubuntu-latest` | 有（`:683` envfork 断言、`:708` core scope、`:764` cer） |
| `:821` | `test-windows`（`:810` 起） | `windows-latest` | 有（`:1042` winsec 层、`:1091` windows scope、`:1116` risk junction；cli 经 `:897` wisp-cli-tests.sh） |
| `:1132` | `slo-smoke`（`:1121` 起） | `windows-latest` | **零**（cache＋build＋slo-check＋upload） |
| `:1205` | `slo-full`（`:1190` 起） | `[self-hosted, wisp-solo]`＝**本机** | **零**（build＋slo-check＋upload） |

票面 `:80` 的三枚"真正依赖 job 级 env 的非测试步骤"（旧 `:400/:497/:560`）现量＝**`:884`／`:1149`／`:1225`**，三处都是 `scripts/build.ps1 -Env dev`，而 `build.ps1:363-367` 收尾跑 `wisp.exe doctor`。
`slo-check.ps1` 自己在 step 级设 `$env:WISP_ENV='test'`（`:319`）＋ `WISP_TEST_DATA_DIR`（`:320`）⇒ **两枚 slo job 的 slo-check 步骤不依赖 job 级 env；job 级 env 在这两枚 job 里唯一到达的生产面就是 build.ps1 的 doctor**。

---

## §2 静态普查：WISP_ENV 的读点与用例名册（全仓现量）

### 2.1 生产读点（非测试码，全部）

| site | 形状 |
|---|---|
| `internal/buildinfo/buildinfo.go:38`（`EnvString`） | `os.Getenv("WISP_ENV")` 非空则赢，否则 build 期注入的 `DefaultEnv`（源码默认 `"dev"`，`:20`；`build.ps1:107 -Env dev` ldflags 注入） |
| `internal/buildinfo/env.go:36`（`ResolveEnv`） | 同序，typed 版 |
| `cmd/wisp/doctor.go:246`（`resolveDataDir`） | `:248-253` portable 分支（本机 `build/portable.txt` 不存在，§4.4）→ **`:255-257` `env=="test"` 早退 `proc.TestDataDir()`（OS 读在它后面，永不执行）** → `:258` 起 `userConfigDir()`（`:276` 变量，绑 `os.UserConfigDir`，128 的 seam） |
| `cmd/wisp/run.go:218` ← `:2131` `buildEnvString()=buildinfo.EnvString()` | `wisp run` 数据根＝ambient |
| `cmd/wisp/secret.go:184` | `wisp secret` 全命令族 `buildinfo.ResolveEnv()`＝ambient |
| `cmd/wisp/slo_windows.go:251-252` | `wisp slo` **自己**：ambient 为空就自设 `test`（自愈，不依赖 job 级 env） |

### 2.2 测试码名册（`grep -rn WISP_ENV --include='*_test.go'`，排除 third_party）

**12 枚文件／43 处**，全部落在四包：`cmd/wisp`（9 枚）、`internal/buildinfo`（1）、`internal/secret`（1）、`internal/winsec`（1）。
（票面 `:50` 那句"10 枚／38 处"是 09-24 的读数；六周开发后枚数涨了，**口径同、枚数漂**，本表为准。）

逐枚的 ambient 取法分三类（这是 AC#1 改写判据的轴 (1)）：

- **A 类·自己 pin**（`t.Setenv` 或子进程 `cmd.Env` 显式给值）＝ambient 无关：`dataroot_128_test.go:305-306`（先栽 `test` 再 pin `dev` 的 premise helper）、`dataroot_128_windows_test.go:87-114`（dev 腿 pin `dev`＋APPDATA 指一次性根）、`early_log_nail_130:123`、`leg_sink_nail_131:409/:526`、`resident_sink_nail_127:221`、`resident_task_source_246:397`（test 腿）＆`:457`（dev 腿）、`resident_task_source_live_246:125/:263`、`secret_argv_windows:236/:337/:371`、`buildinfo/env_test.go:23/:28/:33`（三子例全 pin，含 `t.Setenv("WISP_ENV","")` 强制走默认支）。
- **B 类·显式注入绕开 ambient**：secret 测试具 `newProbe`（`secret_test.go:97` 起）直接构造 `secretCmd{env: env, dataDir: ...}`；全部 `runTextTask` 测试调用点（`approval_reply_201:171`、`approval_unanswerable_375:110/:188`、`firstrun_198:32`、`logsink_windows:152/:319`、`present_wiring_30w1:291/:542`、`run_mode101:146` 等）**都显式传 `dataDir:`**，`runTextTask` 只在 `s.dataDir==""` 时才问 ambient（`run.go:214`）。
- **C 类·裸 ambient（无 pin、无注入）＝本票要抓的形状**：现量**测试层只剩 1 枚文件半枚**——
  ① `ticket224_ac3_realprocess_windows_test.go:44-48`：真进程 `wisp run hi` 只给 `WISP_TEST_DATA_DIR`，**不给 `WISP_ENV`** ⇒ 子进程 env 随 ambient：runner 上=test（注入生效、hermetic），dev host 上=dev（**注入被 `envfork.go:212` 无视**，子进程落真 `%APPDATA%\wisp-dev`）。**它的 hermeticity 今天是 ambient/job 级 env 供养的，且 dev-host 一侧的破口本程已实测兑现**（§4.2.1）。
  ② `dataroot_128_test.go:304`（`pinEnvThatAsksTheOS128` 读 ambient）——但它读了就 pin 覆盖，且 premise 双向有牙（`:309`/`:313` 两道 Fatalf），属"消费 ambient 但不依赖 ambient"。
  ③ `config_reload_223:907-910`／`leg_dispatch_gate_133:1690`／`slo_exit_os_156:322-324`／`panel_host_gate:416-418`：子进程继承 ambient，但子进程跑的是测试二进制自再入（内部各自 pin/传显式目录）或 `go build`/`git`，实测两形无读数差（§4.2）。
- `internal/secret/delete_test.go:159`、`internal/winsec/dataroot_symlink_119_other_test.go:133`：**纯注释**。winsec 那枚文件是 `_other`＝**非 Windows** build tag ⇒ 在 Windows 测试二进制里**根本不编译**（本机 winsec 两形跑的都是 windows-tag 的 76 枚顶层）；它的 ubuntu 分母在 core roster（`portable-tests.sh` core_pin 含 `internal/winsec`）。

### 2.3 CI 分母地图（轴 (3)；`ci.yml`＋`portable-tests.sh` 现量）

| 包 | ubuntu（test-core, env=test） | windows（test-windows, env=test） | 两枚 slo job（env=test） |
|---|---|---|---|
| `internal/buildinfo` | **core scope 有**（core_pin `:298`） | **无**（windows scope、cli、winsec 层都不含它）⇒ `A1`/verdict ⑥ 在 HEAD 仍成立 | 无 |
| `internal/secret` | core scope 有 | windows scope 有（`:1091`） | 无 |
| `internal/winsec` | core scope（glob）有 | winsec 层（`winsec-tests.sh`→`--scope=winsec`） | 无 |
| `cmd/wisp` | **无**（core 名册无 cmd/wisp） | **cli 层**（`wisp-cli-tests.sh:304`→`--scope=cli`） | 无 |
| 生产面 `wisp doctor`（build.ps1） | 无 ubuntu 分母 | `:884` build step（env=test） | `:1149`／`:1225`（env=test） |

`ci.yml:668` 那枚"Environment fork assertion (WISP_ENV=test data dir)"步骤：跑的是 `TestLayoutForTestEnv`（`internal/proc/envfork_test.go:75`），它传 `buildinfo.EnvTest` **字面量**＋自己 `t.Setenv(TestDataDirEnv)`，`internal/proc/*.go` 对 `WISP_ENV` 全文只有 `doc.go:13` 一条注释 ⇒ **该步骤不读 job 级变量**。`SPEC-11:182-183`（CI 显式断言 env 生效）与 `SPEC-03:101`（CI 去设它）两条规范义务，今天仍然**零枚用例履行**（现量复认，票面 ③ 的结论在 HEAD 不变）。

---

## §4 真跑读数（AC#1 判据的"必须真跑"那半）

两形定义：**UNSET 形**＝`env -u WISP_ENV`（本机 shell 现量 ambient 本来就是 UNSET，§0 已核）；**TEST 形**＝`WISP_ENV=test`。四数口径＝`-v` 输出的顶层 `=== RUN`／`--- PASS`／`--- FAIL`／`--- SKIP`（顶层＝名字不含 `/`），panic 另计。

### §4.1 三枚 internal 包（`go test -count=1 -v`，两形各一发，2026-11-06 09:4x）

| 包 | UNSET：RUN/PASS/FAIL/SKIP/panic | TEST：RUN/PASS/FAIL/SKIP/panic | 名册差集 | 换色 |
|---|---|---|---|---|
| `./internal/buildinfo/` | 2/2/0/0/0 | 2/2/0/0/0 | **0** | **0** |
| `./internal/secret/` | 25/25/0/0/0 | 25/25/0/0/0 | **0** | **0** |
| `./internal/winsec/` | 76/76/0/0/0 | 76/76/0/0/0 | **0** | **0** |

三包合计 103 枚顶层用例，两形**逐名同色**（color-aware diff 只剩计时尾数差）。附注：winsec 的 POSIX 那枚（`_other_test.go`）不在 Windows 二进制里；它的 ubuntu 分母读数引用票 111 `AC#9` 的容器实测（RUN=35/PASS=20/FAIL=0/SKIP=0），本程未在容器复跑（本程仪器是宿主两形，不冒充容器）。

### §4.2 `./cmd/wisp/`（两形 × 两仪器；本程最重要的 instrument 教训）

**仪器 A＝裸 `go test`（PATH 带 third_party）**——UNSET 形完整跑完：
**顶层 489/443/44/2，panic=0**。44 枚 FAIL 逐枚归因（逐名 grep 失败句，清单 `.scratch` 外 `/d/tmp/140/cmdwisp.unset.fails`）：

- ~20 枚：`onnxruntime 1.17.1 is older than the 1.28 ... a second onnxruntime.dll is shadowing the one on PATH (the loader searches System32 before PATH)`——**本机 System32 里 Windows 自带的 1.17.1**（`C:\Windows\System32\onnxruntime.dll` 现量 FileVersion `1.17.260311-...`，mtime 05-26）压过 PATH 里的 1.28.2。**机器形状，与 env 无关。**
- ~17 枚：`WISP_SPIKE_MODELS is not set`——语音模型未布。**机器形状。**
- 其余 ~7 枚：同一批语音 E2E 的连带（DLL/模型两种因）；另有 4 枚值得另立账的非语音红：`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`（自述 7 of 11 vs 期望 5）、`TestAC247HandingTheLevelToTheSeamIsNotVisibility`、`TestAC4FocusReturnToPriorWindowGap33r5`（前台焦点竞态）、`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`。**这 44 枚里 0 枚的红句提到 WISP_ENV／数据根／userConfigDir** ⇒ 全部**非 140 账**；那 4 枚非语音红是本机现存红（CI 是否同样红本程未追），按票面规矩**登记不扩权**。
- 2 枚 SKIP＝live 麦克风探针（`TestAC247LiveMicrophoneLevels...`、`Test16r1PremortemMicPermissionDeniedLiveProbe`），两形同 SKIP。
- UNSET 形总时长 580.7s（go test 自报），已贴近 `go test` 默认 10m 上限——这就是 TEST 形第一发撞 `panic: test timed out after 10m0s`（死在 `TestTicket101SessionGrantDoesNotCrossRestart`，仅完成 408/489 顶层）的原因；完整 TEST 形读数用 `-timeout=25m` 补跑（§4.2.3）。

**§4.2.1 UNSET 形的真 profile 写入（本程最重要的一枚实测，承认 §4.2 初稿"零写入"那句已被推翻）**：
`TestTicket224ProductionSessionChangesAcrossRealOSProcesses`（PASS，12.44s，RUN 于日志 `:4858`）起的真子进程只带 `WISP_TEST_DATA_DIR`（`:44-48`）→ ambient UNSET ⇒ 子进程 env 落 `dev` ⇒ 注入被 `envfork.go:212` 无视 ⇒ `resolveDataDir(dev)` 答**真 `%APPDATA%\wisp-dev`**。实测痕迹（现量）：`C:\Users\swq\AppData\Roaming\wisp-dev\wisp.db`（163840 B，mtime 09:57）与 `logs\wisp-20261106-001.jsonl`（5833 B，09:57:22.6 起：`persistent log sink installed dir=...wisp-dev\logs`→panel SNAPSHOT→tools C25 audit→config HOT-RELOAD stopped，即一次完整 `wisp run` 生命周期落在机主真 dev profile 上）。**同一用例在 TEST 形下 hermetic**（子进程继承 `WISP_ENV=test`，注入生效）。本程**未删、未改**这两枚文件（无删除政策）；机主 dev profile 里多了一次探针任务记录，如实登记。

TEST 形（仪器 A）：本程第一发被本程**主动取消**（发现 DLL 遮蔽后换仪器重跑；见 §7 纪律），补跑一发完整 TEST 形读数见 §4.2.3。

**仪器 B＝仓库自己的门（`scripts/wisp-cli-tests.sh`，CI test-windows 的原路）**——两形各一发：
- 两形**同一枚用例上硬崩**：跑到 `TestAC13BringUpSurvivesAReusedThreadQuit`（`=== RUN` 已印、结果行永未出现，测试二进制整进程死，`FAIL ... 253.785s`，rc=1）。UNSET 形 10:00:32 起 10:05:12 止，TEST 形 10:05:12 起 10:10:00 止——**两形崩点相同 ⇒ 该崩与 WISP_ENV 无关**（本身就是一条"分支不变"的读数）。同一枚用例在仪器 A 下两形都 **PASS**（0.86s）⇒ 崩是仪器 B 在本机（交互桌面＋LOCALAPPDATA gotmp 根）的形状，不是包的形状。**此崩未追**（非 140 账；登记给 cli 门的主人）。
- 崩点之前的读数（runtests 严格 runner 四数，两形逐位同）：**RUN=277（含子例）／顶层 PASS=209／FAIL=0／SKIP=0**；209 枚顶层名册两形 color-aware diff＝**只有计时尾数差，换色 0**。

**§4.2.3 合并结论（完整体）**：TEST 形补跑（`-timeout=25m`，2026-11-06 10:26:07–10:37:06，653.2s）＝**顶层 489/442/45/2，panic=0**，与 UNSET 形同一 489 枚名册。两形 FAIL 名册差集＝**仅 1 枚**：`TestAC1ResidentLegInstallsItsLogListenerOnDisk`（TEST2 红 8.42s，UNSET 绿 6.90s，第一次 TEST 形绿 8.54s）。**专项消歧**（`-run ^...$ -count=3` 两形各一发）：**6/6 PASS**（UNSET 7.33/7.47/7.07s，TEST 7.15/7.60/5.80s）⇒ 那一枚红是**计时型偶发**（与它常红的邻居 `...BooksItsShutdownBeforeClosingTheSink` 同属 shutdown-record 顺序断言族，后者两形都红、已列移交清单），**不归 WISP_ENV**。⇒ **四包合并终读：可归因于 `WISP_ENV` 分叉的换色用例＝0；名册差集＝0。** ticket224 四枚两形全 PASS；10:25 之后真 profile 零新写入（TEST 形 hermetic 的正面读数）。

### §4.3 生产面两形探针（M2，决定性的一发）

`go build -o /d/tmp/140/wisp.exe ./cmd/wisp`（plain go build ⇒ `DefaultEnv="dev"` 未被 ldflags 覆盖，与 CI `build.ps1 -Env dev` 的注入值同值），APPDATA 指一次性根，同一枚 exe、同一枚命令 `wisp.exe doctor` 两形各一发：

| 形 | doctor 自报 env | data dir 行 | OS 问了吗 |
|---|---|---|---|
| UNSET | `WISP_ENV=dev` | `[PASS] data dir writable (dev)  D:\tmp\140\adu\wisp-dev`（探针文件写了又删，`adu/wisp-dev` 目录被建） | **问了**（`userConfigDir` 执行，落到 APPDATA） |
| TEST | `WISP_ENV=test` | `[PASS] data dir writable (test)  C:\Users\swq\AppData\Local\Temp\wisp-test-6116`（`adt` 根零字节） | **没问**（`:255` 早退） |

（doctor 总 verdict 两形都 FAIL——那是我这枚 build 没有 ldflags 版本 pin 的缘故，与数据根无关；数据根那两行才是本探针的读数。）
⇒ **"换分支但不换色"在生产面是实测事实，不是推理**。把它放到 CI：三枚 doctor step 全在 job 级 `test` 之下 ⇒ **runner 上 `userConfigDir` 那一行（`doctor.go:258`＋`internal/proc` 的同族读）CI 分母＝0**。`wisp secret`（`secret.go:184`）与 `wisp run`（`run.go:218`）同族，CI 上同样只有 test 分母。

### §4.4 M3：`build\portable.txt` 在 wisp-slo 机器上

`ls build/portable.txt` ⇒ **不存在** ⇒ 本机（`[self-hosted, wisp-slo]`）doctor 不走 portable 分支，env 分支就是它的实况——§4.3 的两形探针在本机等价于 runner 实况。

---

## §5 逐枚具名用例三轴表（AC#1 改写判据①②③；只列与 WISP_ENV 分叉有关的，其余用例两形同色已由 §4 盖章）

| 用例（file:line 现量） | (1) ambient 来源 | (2) 纯函数？ | (3) CI 分母 |
|---|---|---|---|
| `buildinfo/env_test.go:20 TestResolveEnv`（3 子例） | 自 pin（`t.Setenv` test/nonsense/**空串**）——空串那枚**在两种 ambient 下都强制走默认支**，pin 有牙 | 纯函数 | **仅 ubuntu core**；Windows 零分母（`A1` 在 HEAD 复认） |
| `cmd/wisp/dataroot_128_test.go:202 TestAC2TestDataDirBranchStillResolves128` | 显式传 `"test"`，不读 ambient | 否（读 `proc.TestDataDir`，不碰 OS config 读） | cli 层（windows） |
| `cmd/wisp/dataroot_128_test.go:319 TestAC2EveryLegRefuses...128`（5 腿表） | `pinEnvThatAsksTheOS128`（`:304`）：读 ambient→栽 `test`→pin `dev`→断言 pin 赢＋seam 被问 | 否（rebind `userConfigDir`） | cli 层 |
| `cmd/wisp/dataroot_128_windows_test.go`（真进程 5 腿） | pin `dev`＋APPDATA 指一次性根（`:87-114`） | 否 | cli 层（windows 才编译） |
| `cmd/wisp/ticket224_ac3_realprocess_windows_test.go:34` | **裸 ambient**（子进程只给 `WISP_TEST_DATA_DIR`） | 否 | cli 层；**hermeticity 由 job 级 env 供养** |
| `cmd/wisp/secret_test.go`（probe 族） | 构造时显式 `env:`，ambient 无关 | 混合（store 落盘但目录显式） | core（ubuntu）＋windows scope；**生产面 `secret.go:184` 的 ambient 分叉 CI 零分母** |
| `internal/winsec/dataroot_symlink_119_other_test.go:133/135` | pin `TMPDIR`／`HOME`（不读 WISP_ENV；`:133` 仅注释提它） | 否（真 symlink 树） | **ubuntu core 有；Windows 二进制不编译** |

### 5.1 恒真断言点名（轴 (4)）

- **活的一枚**：`cmd/wisp/dataroot_128_test.go:209-210`——`resolveDataDir("test")` 的返回值去比 `proc.TestDataDir()`，而后者正是 `:255` 早退**原样返回的那个函数**：挡得住"早退被删"，挡不住"早退去错地方"（`R-140-1` 原判据）。`A188⑦` 自决的改码（撤销口令「撤 140 恒真改码」）**在 HEAD 84b0a7c4d 尚未落地**（该形状仍在）。M5 变异（把早退目标换成真目录⇒必须红）随改码程走，本程不带。
- **声明过的半枚（不计 R-140-1）**：`buildinfo/env_test.go:41-45` 默认支用 `ParseEnv(DefaultEnv)` 算期望——同一包同一常量自比；其注释**自己声明**了为什么只能相对自身（`DefaultEnv` 是 ldflags 注入，测试里只能是源码值）。属"声明过的局限"不是"未登记的恒真"。
- winsec `:349/:375` 两枚 `proc.TestDataDir()` 比的是测试自己声明的字面量（`declared`/`injected`）——**不是**恒真；`cleanSpelling119`（`:125` 注释）明写故意用 `EvalSymlinks` 而非被测函数——修过的形状。

### 5.2 "这条腿测的不是它声称的东西"名册（轴 (5)，每条自带"今天实际测了什么"）

1. **三枚 `wisp.exe doctor` step（`:884/:1149/:1225`）**：声称＝"build 后 doctor 自检"；今天实际测＝**test-env 数据根路线**（早退→TestDataDir→probe 写 temp）。dev/prod 的 `%APPDATA%` 分叉在 CI 的任何 step 都没人测（§4.3/4.4）。
2. **`ci.yml:668` step 名**：名字说 "WISP_ENV=test data dir"；今天实际测＝`LayoutFor(字面量 test)` 的 fork 布局，**不读那枚环境变量**。SPEC-11/SPEC-03 要求的"断言 env 真的生效"今天零履行。
3. **票 128 那 4 枚 CLI 用例（`审批超时（1/300 秒未确认）` 族）**：前置约束要求它们"出现但标同名不同因"。现量：**该名目在 HEAD 的 `cmd/wisp/*_test.go` 与本程两形名册里都是 0 命中**（grep `审批超时｜1/300` 双零）——六周间已被改名/重构，今天名册里没有它们 ⇒ 前置 (iii)"登记清单不许出现票 123"自动满足，(i)(ii) 以"不在分母"的形状满足。它们的 CI 红因（`128-ac4-r1` §9#7 登记）另有主人，本程未追。

---

## §7 裁形清单（AC#2 三支＋实测代价；**编排者终裁**，本程只供读数与推荐）

| 支 | 内容 | 实测代价（本程读数） | 实测收益 | 判 |
|---|---|---|---|---|
| (a) 字面（4 删 3 增：三枚 build step 步骤级补设 env） | 摘 job 级 env、原 step 自己补 | **颜色不变的空转**（`A189④` 判过；本程两形差集 0 复认——补回去等于没摘） | 0 | 否 |
| (a′) 4 删 0 增（真摘，让 runner 问 OS） | 三枚 doctor step 翻 dev 支：`probeWritable` **往 runner 真 `%APPDATA%\wisp-dev` 写**（slo-full 就在本机、复用工作树，`A189④` 的代价被本程 §4.3 探针坐实）；**`ticket224` 的 hermeticity 破**（`:44-48` 只给注入不给 pin，`envfork.go:212` 在 dev 下无视注入 ⇒ 子进程起在 runner 真 profile 上——**这一形状本程已在 dev host 实测兑现一次**，§4.2.1，09:57:22 写 `wisp.db`＋logs）；`resident_task_source_live_246` 等 pin 死的用例不受影响（pin 仍赢） | 让三枚 doctor step 的 OS 读真正被 CI 咨询（买回票面标题主张的那个面） | **不推荐现在摘**——代价全落在"与 140 无关的用例"头上，且 (a′) 一天不登记就一天不能动 `ci.yml` |
| (b) 保 env，受影响用例自己钉（128 helper 形状推广） | 今天四包测试层**已经全是 pin/显式注入**（§2.2 A/B 类）⇒ (b) 的**剩余工作量≈1 枚**：给 `ticket224` 补一行 `WISP_ENV=test` pin（把"靠 job 级 env 供养"变"自己钉"，正是票面判据 (b) 说的"先栽后钉"）；**不修票面标题**（生产 doctor 在 runner 上仍答 test——这格是 (c) 的账） | 零 runner 写风险；与票 135 在 `dataroot_128_test.go` 的撞面只剩 1 行、错开即可 | **主支**（测试层） |
| (c) 现状＋响亮登记 | 登记三笔（**不含票 123 任何名字**，§5.2③）：①三枚 doctor step"测的不是它声称的"（§5.2①）＋`userConfigDir` 的 CI 分母＝0；②`ci.yml:668` 步骤名与所测不符＋SPEC-11/SPEC-03 义务零履行；③`buildinfo` Windows 零分母（`A1`） | 把"装饰"变"挂牌的装饰"，可被后续票领走 | **并做**（登记是文档，不动共享件） |
| (b/c 之后的补刀·可选) | **加一枚显式 step**：doctor 以 `WISP_ENV=dev`＋一次性 `APPDATA` 跑一发，断言 data dir 行落在 `$APPDATA\wisp-dev`（本程 §4.3 探针就是现成判据，搬到 ci.yml 只是一个 step） | 买回 (a′) 想买的"OS 读被 CI 咨询"，而**不摘任何 env、不破任何 hermeticity、不写 runner 真 profile** | **这是买 (a′) 价值的正确形状**；但它是 `ci.yml` 改动 ⇒ **先登记** |

**AC#2 前置约束三问**：(i) 若最终仍选 (a′)：票 123 那 4 枚**不在今天名册**（已改名），"取消 env 后它们仍红"的担忧在字面上已无对象；它们的红因账仍在 `128-ac4-r1 §9#7`。(ii) 选 (b)：它们不在本票分母。(iii) 本清单零票 123 名字。**票面判据 (a) 的"会不会有用例从绿变从没跑过"（98/111 族）**：现量两形名册差集为 0 ⇒ 摘 env 不会让任何**测试**消失分母；会变的只有**生产 step** 的分支（已列入 (a′) 代价行）。**判据 (b) 的"自己给自己设前提的假绿"**：`pinEnvThatAsksTheOS128` 的 premise 双向有牙（先栽后钉，`:305-313`），`t.Setenv` 空串那枚（buildinfo）在两种 ambient 下都强制默认支——现量没有"钉无牙"的用例。

**R-140-1 在终裁表的可查清单**：`dataroot_128_test.go:209-210`（§5.1，改码未落地于 HEAD）。

---

## §8 STOP 格：ci.yml 未登记

`AC#2` 无论选 (a)/(a′) 还是"补一枚 dev 断言 step"，都是 **`.github/workflows/ci.yml`（共享件）改动**；本程**未获登记、未获授权**，按派单令与票面 Rules（"`ci.yml` 属共享件：改前先登记交回编排者"）**停在登记之前**：本程一行 ci.yml 未动（`git status`/`git diff` 双零，§9）。登记之后，实现判据现成：§4.3 探针＝那枚 step 的验收形状；§7 表＝代价清单。

## §9 副作用与纪律自证

- **仓库改动**：本程只有本文件＋票面 log 追加（各一枚显式 pathspec commit）。`git status --porcelain` 的 tracked 集在写档前现量**零改动**；全程未跑 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`push`；**未删任何件**。
- **仓库外临时件**（不删，列出）：`/d/tmp/140/`（两仪器四份 `-v` 日志、doctor 探针 exe＋DLL＋两个一次性 APPDATA 根、名册/diff 中间件）。
- **真 profile**：UNSET 形下 `ticket224` 子进程按 §4.2.1 写入了 `wisp-dev\wisp.db`＋`logs\wisp-20261106-001.jsonl` 两枚（机主 dev profile，非本程起见已有 328 文件的树）；**未删未改任何既有文件**；doctor 探针的 APPDATA 全部指一次性根。
- **取消过一发自己起的跑**：仪器 A 的 TEST 形第一发（裸 `go test`）在发现 System32 DLL 遮蔽后由本程 `job_kill` 主动取消（污染仪器，不污染仓库）；cli 门两形与仪器 A 两形全部**串行**，未同包并行。
- **判为注入的工具输出：0**；真通知回显：后台 job 完成通知 3 条（本程自己起的跑，正常回显）。
- **GATES**（收尾现量，2026-11-06 10:1x）：`go build ./cmd/... ./internal/... ./tools/...` rc=0；`/d/work/base/gopath/bin/gofumpt.exe -l internal cmd tools scripts frontend` 输出空（v0.12.0/go1.27.1）。
- **本机现存红移交清单**（非 140 账、未追）：§4.2 仪器 A 那 4 枚非语音红（256/247/33r5/127 各一枚）＋仪器 B 在本机的 `TestAC13BringUpSurvivesAReusedThreadQuit` 硬崩（两形同点，仪器 A 两形同 PASS）。
