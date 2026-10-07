# 274-a2 出货通路普查 — 页面字节的落地机在哪台机器上

> 腿：`274-a2`（只读普查，Wisp 仓 `D:/work/workspace/projects plans/Wisp`，分支 `dev`）。
> 交付件本体＝本文件；每条读数指名它出自 `logs/` 哪枚文件。
> 所有读数均为本腿**现跑**；票面/台账旧数只作为「待复量的对象」出现，不当凭据。
> 唯一写面＝`.scratch/wisp/probes/274/a2/**` 与本腿的 commit。

## §0 起手锚

原始输出：`logs/anchor.txt`。

| 尺 | 读数 |
|---|---|
| `date` | `Wed Oct  7 11:02:00 CST 2026` |
| `git rev-parse --short HEAD` | `933a8341` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd internal scripts tools .github docs frontend` | **0 行（干净）**，退出码 0 |
| `git ls-files frontend/src \| wc -l` | **64** |
| `git ls-files frontend \| wc -l` | **85** |

工作树现状说明（本腿不处理）：`design/**`、`.gitignore`、`.scratch/**` 里有别的腿/机主的**未提交改动**（含 `design/**` 的一批 `D` 删除记录与若干 `??` 未跟踪件），**不是本腿的**；本腿既不还原也不提交它们。终态自查一律只按上面那六族路径（`cmd internal scripts tools .github docs frontend`）判定。

复量对照（待复量对象，非凭据）：编排者给的背景与工单 `274` 里「dev 的 64 枚页面源码」这一枚，本腿现跑 `git ls-files frontend/src` **＝ 64，对上**。

---

## §1 产 exe 的通路到底是谁

原始输出：`logs/s1-buildps.txt`（`param(` 尺、main 包普查、param 块、59–89 行、步标记、`wc`）、`logs/s1-frontend-step-raw.txt`（59–89 行**未截断**原文）、`logs/s1-exe-producers.txt`（`grep -rn 'go build' scripts tools .github`）、`logs/s2-steps-bodies.txt`（`git grep 'go build' -- '*.sh' '*.ps1' '*.yml'` 补尺）。

### 1.1 `scripts/build.ps1` 全文结构（现读全文，非节选）

尺＝`wc -l -c scripts/build.ps1` ⇒ **172 行 / 7346 字节**。
尺＝`grep -n '^# --- ' scripts/build.ps1` ⇒ 命中 **8 枚段标记**（步号是脚本自己的编号，含一枚 `4b`）：

| 段标记行 | 段名 | 射程行 | 干什么 |
|---|---|---|---|
| `41` | `--- 0. toolchain discovery` | 41–67 | 找 `go.exe`（42–46）、找 mingw gcc（48–66，候选含 `$env:CC`、`gcc.exe`、`%MINGW64_ROOT%\gcc.exe`、硬编码 `C:\msys64\mingw64\bin\gcc.exe`），找不到即 `Fail` |
| `69` | `--- 1. fetch-deps` | 69–71 | 调 `scripts/fetch-deps.ps1 -RepoRoot $RepoRoot` |
| `73` | `--- 2. frontend` | **73–74** | **整段只两句：一行注释 + 一行 `Write-Host`。没有任何构建动作、没有 `npm`、没有 `frontend`。** |
| `76` | `--- 3. go build` | 76–134 | 读 `deps.toml` 版本（77–94）、取 commit（96–101）、拼 ldflags（108–116，含 `-H=windowsgui`）、建 `build\` 目录（118–119）、设 `CGO_ENABLED=1`/`CC`/`GOOS=windows`/`GOARCH=amd64`（121–125）、**`go build` 唯一次调用在 129** |
| `136` | `--- 4. colocate DLLs` | 136–139 | `Copy-Item third_party\sherpa-onnx\*.dll → build\` |
| `141` | `--- 4b. colocate the signed C29 model manifest` | 141–149 | 拷 `models\manifest.json` + `.minisig` → `build\models\` |
| `151` | `--- 5. SHA256SUMS` | 151–164 | 只对 4 枚产物求哈希：`wisp.exe`、`onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`（`152`） |
| `166` | `--- 6. smoke test` | 166–169 | `build\wisp.exe doctor`，非 0 即 `Fail` |

**① `go build` 的行与入口**：唯一次真正编译在 **`scripts/build.ps1:129`**：

```
& $go.Source build -trimpath -ldflags $ldflags -o (Join-Path $outDir 'wisp.exe') ./cmd/wisp
```

⇒ 入口 **只有 `cmd/wisp` 一枚**，产物 **只有 `build\wisp.exe` 一枚**。`130` 是退出码闸门，`134` 是成功打印。
尺＝`grep -rn '^package main' cmd --include=*.go` ⇒ `cmd/` 下共有 **3 枚 main 包目录**：`balldebug`、`llmrecord`、`wisp`（尺＝`ls cmd`）。`build.ps1` 只编 `cmd/wisp`；另两枚不由这条通路产出（`balldebug` 见 1.3）。

**② 参数开关**：尺＝`grep -n 'param(' scripts/build.ps1` ⇒ **命中 1 枚，在 `22` 行**。param 块整块原文（`logs/s1-buildps.txt`，`18–30` 行）：

```powershell
[CmdletBinding()]
param(
    [ValidateSet('dev', 'prod')]
    [string]$Env = 'dev'
)
```

⇒ **今天只有 `-Env dev|prod` 这一枚开关**（写进 `buildinfo.DefaultEnv`，`112`）。**没有 `--with-frontend`、没有任何形式的 frontend 开关**（这条与本腿 §1.1 表第三行互锁：frontend 段无条件跳过）。三枚调用点传的也都是 `-Env dev`（见 §2）。

**③ `frontend step skipped` 前后各 15 行（59–89 行）逐行**：全文未截断落在 `logs/s1-frontend-step-raw.txt`；这一带最长行只有 102 字节（尺＝`awk 'NR>=59 && NR<=89 {print NR": "length($0)}'` ⇒ 逐行长度 99/1/78/76/74/31/73/59/77/0/78/64/77/0/78/56/23/84/27/36/67/20/9/41/102/5/16/1），**没有需要整行读的超长行**。承重两行整行原文：

```
73:# --- 2. frontend ------------------------------------------------------------
74:Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'
```

**读数（承重的负向）**：第 74 行**是脚本里唯一一句关于 frontend 的话**，它既不调 `npm`、不建 `frontend/dist`、也不拷任何东西进 `build\`。前 15 行（59–72）是 gcc 发现 + `fetch-deps`，后 15 行（75–89）是 `# --- 3. go build` 段头与 `Get-DepsValue` 函数体。

### 1.2 出货 exe 里今天确实没有页面字节（仓内字节层）

尺＝`git grep -n 'frontend[\\/]dist' -- .`（只看跟踪文件，`logs/s2-steps-bodies.txt` 同批）⇒ 代码侧命中：`.gitignore:22` `frontend/dist/*`、`.gitignore:23` `!frontend/dist/.gitkeep`、`cmd/wisp/panel_host_gate_test.go:147`（那句报错文案本身承认「committed under frontend/dist」要超过锚文件才算数）、`.github/workflows/ci.yml:902`（注释）。⇒ 仓里入库的 dist 只有 `.gitkeep` 那一枚锚（与工单 274/248 的说法同向，本腿是自己现跑的）。
尺＝`git grep -n 'go:embed' -- cmd internal` ⇒ `cmd/wisp/panel_host_windows.go:11` 与 `cmd/wisp/panel_assets.go:19` 都指向 `internal/panel` 的 `//go:embed all:dist` bundle（`internal/panel/assets.go:63` 注释：锚文件的作用就是让 go:embed 在纯净树里还过得去）。

### 1.3 仓里除 `build.ps1` 外还有谁调 `go build`

尺＝`grep -rn 'go build' scripts tools .github 2>/dev/null`（`logs/s1-exe-producers.txt`）⇒ 逐枚点名，**能产二进制的只有 2 枚，其中只有 1 枚产 wisp.exe**：

| 文件:行 | 命令 | 算不算「产出货 exe」 |
|---|---|---|
| `scripts/build.ps1:129` | `go build ... -o build\wisp.exe ./cmd/wisp` | **是——唯一通路** |
| `scripts/dev/ball-cycle.ps1:30` | `& $go build -o $exe (Join-Path $repo "cmd\balldebug")` | 否：dev 腿的调试球载体，产 `balldebug`，不入 `build\`、不入 SHA256SUMS |
| `tools/d22scan/scan_test.go:2259`、`tools/d22scan/selftest_test.go:184` | 测试里的 `go build failed: %v` 断言文本/内建编译 | 否：测试自身体内编译，不产出货件 |
| `scripts/spike/run.ps1:35`、`scripts/spike/xy-verdict/main.go:160`、`session_cgo.go:5` | 只出现「cgo builds」字样，无 `go build` 调用 | 否 |
| `.github/workflows/ci.yml:466/520/536/566/809` | 注释或步名里的 `cgo build`，非命令 | 否（真正的命令是 `571/757/820` 三行调 `build.ps1`，见 §2） |

补尺（怕漏 shell 侧）＝`git grep -n 'go build' -- '*.sh' '*.ps1' '*.yml'`（`logs/s2-steps-bodies.txt`）⇒ 除上面这些，其余命中全在 `.scratch/wisp/probes/**` 的历史普查脚本里（`152/e2e-subject-kill*.sh`、`153/*`、`161/r1/bench.sh`、`139/run-mutations.sh`、`111/r4/ci-at-HEAD.yml` 快照），**都不是出货通路**。

---

## §2 三枚产 exe 的 CI job 跑在哪台机器上

原始输出：`logs/s2-ci-jobs.txt`（job 名尺 / `runs-on` 尺 / `setup-node` 尺 / `npm` 尺 / `frontend/dist` 尺）、`logs/s2-ci-detail.txt`（`needs:`、artifact、`uses:`、三段 job 头原文）、`logs/s2-workflows.txt`（workflow 文件清单 + `wisp.exe` 尺）、`logs/s2-steps-bodies.txt`（upload 步正文 + `release` 尺）、`logs/s1s2-extras.txt`（步级结论普查）。

**先量「哪几枚 job 真会跑到 build.ps1」，不按 job 名猜**：尺＝`grep -rn 'build\.ps1' .github` ⇒ 命中 8 行，其中**只有 3 行是 `run:` 命令**（`571`、`757`、`820`），另 5 行是注释/步名（`466`、`566`、`579`、`809`，加 `scripts/build.ps1` 自身）。这 3 行按行号落进 §1.1 之外的 job 区间 ⇒ **产 exe 的 job ＝ `test-windows`(505–738) / `slo-smoke`(739–796) / `slo-full`(797–870)** 这三枚，逐枚点名见下表。

| job（行号） | `runs-on` 原文（行号） | 跑 build.ps1 的那一步（行号） | 步级实际结论 | 有 `setup-node`？ | 有任何一步读 `frontend/dist`？ |
|---|---|---|---|---|---|
| `test-windows`（`505`） | `windows-latest`（`506`） | `"cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)"`，`566` 步名 / **`571` 命令**，带 `if: ${{ !cancelled() }}` | 该步 **success**（整 job failure，红在第 7、9 步） | **无** | **无** |
| `slo-smoke`（`739`） | `windows-latest`（`740`） | `Build wisp.exe`，`756` 步名 / **`757` 命令** | **success**（整 job success） | **无** | **无** |
| `slo-full`（`797`） | `[self-hosted, wisp-slo]`（`798`） | `Build wisp.exe (deps cached on the runner)`，`819` 步名 / **`820` 命令** | **success**（整 job success） | **无** | **无** |

三枚都对 `scripts/build.ps1` 传 `-Env dev`（逐字：`powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`）。
另外 `slo-full` 的 `env:` 块里有一条**承重**的具名指认（`797–825` 原文在 `logs/s2-ci-detail.txt`，`86` 行）：`MINGW64_ROOT: E:\work\base\msys64\mingw64\bin`，注释逐字写着「The runner's job environment carries no msys64 on PATH … (on the interactive user PATH only) … build.ps1's own fallback path is `C:\msys64\...`, which is the wrong drive here」⇒ **同机、但 runner 作业环境 ≠ 交互式用户 PATH**，这条对本腿 §3 的 node 问题是直接可迁移的风险。

### 2.1 每条负向读数的正控（同一把尺在别处的命中数）

| 尺 | 产 exe 三枚 job 内命中 | 正控（全文件命中处） |
|---|---|---|
| `grep -n 'setup-node'` | 0 | **1 枚**：`879` 在 `lint-frontend`（`871`，`ubuntu-latest` `872`） ⇒ 尺能命中，负向成立 |
| `grep -n 'npm'` | 0 | **10 枚**：`884`、`887`、`888`、`891`、`894`、`900`、`903`、`916`、`928`、`933`、`940`（逐名 11 行，含 `cache: npm`），**全部落在 `lint-frontend` 区间 871–941** ⇒ `npm run build` 只在**不产 exe** 的那枚 job 里 |
| `grep -n 'frontend/dist'` | 0 | **1 枚**：`902` 步名注释「build (vite build -> frontend/dist, the bytes go:embed carries)」，也在 `lint-frontend` |
| `grep -rn 'wisp\.exe' .github/workflows` | 命中 4 行但**无一是产物落点**：`541` 注释、`756`/`819` 步名、`912` 注释（指 `wisp.exe panel-assets -l2 fs.delete`） | 尺本身命中 ⇒ 成立：**没有任何一步把 exe 里的页面字节看一眼** |
| `grep -n 'needs:'` | 0 | **全文件 0 枚**——**这把尺没有正控可用** ⇒ 读数只能写成「`needs:` 命中 0 枚」，⛔ 不许由它推断任何顺序/无阻塞结论（本腿只复量了同一枚 0，与 `274-a1` 的 0 同形） |
| `grep -n 'download-artifact'` | 0 | **全文件 0 枚**；同族尺 `upload-artifact` 命中 **2 枚真步**（`763`、`826`）+ 2 枚注释行（`842`、`843`）⇒ upload 尺有正控、download 尺负向成立 |
| `grep -rn 'release' .github/workflows` | 0 | **两个 workflow 文件全 0 命中**，无正控 ⇒ 只能写「命中 0 枚」。job 名可枚举的正控另给：尺＝`grep -n '^  [a-zA-Z0-9_-]*:$' .github/workflows/ci.yml` 命中 12 行（`11/12/36/38` 属 `on:` 块，`65/395/505/739/797/871` 属 jobs）⇒ **ci.yml 里一共就 6 枚 job，没有名为 `release` 的 job**（这与编排者转述的「`release` job 无落点」略有出入：**不是有 job 没落点，而是根本没有 `release` job**，本腿现量具名报回） |

### 2.2 产出的 exe 今天没有被任何通路留下

两枚 `upload-artifact` 步（`slo-smoke:763`、`slo-full:826`）正文逐字＝`name: slo-smoke-report` / `name: slo-full-report`、`path: build/slo/slo-report.json`（`logs/s2-steps-bodies.txt`）⇒ **上传的是 SLO 报告 JSON，不是 `build\wisp.exe`**；`build\` 下那 4 枚进 SHA256SUMS 的产物（§1.1 表 `152` 行）在三枚 job 里**谁都不上传、谁都不下载、谁都不验**。
workflow 文件清单（尺＝`ls -la .github/workflows`）＝`ci.yml`（57948 字节）+ `slo-fresh.yml`（4038 字节，`runs-on: ubuntu-latest` `62`，⛔ 不跑 `build.ps1`：`grep -rl 'build\.ps1' .github/workflows` 只命中 `ci.yml`）⇒ **整个仓里产 exe 的 CI 通路就 §2 那三枚 job，一台 hosted Windows × 2、一台 self-hosted Windows × 1。**

---

## §3 那台机器今天有没有 node/npm（run 侧真读数）

**取数方式**：`gh run list --branch dev --limit 6 --json databaseId,conclusion,headSha,workflowName,createdAt,event`（只读）⇒ 6 发读数在 `logs/` 之外的 `D:/tmp/wisp274a2/gh/run-list.json`（原始 JSON，MB 级日志不进仓也不进上下文，遵守纪律 4）：

| databaseId | workflow | event | conclusion | createdAt | headSha |
|---|---|---|---|---|---|
| `37545246395` | ci | schedule | failure | 2026-10-06T23:12:27Z | `cc31526` |
| `37511647900` | slo-fresh | schedule | success | 2026-10-06T18:29:01Z | `cc31526` |
| `37424357373` | slo-fresh | schedule | success | 2026-10-06T06:33:02Z | `cc31526` |
| `37406757402` | ci | push | failure | 2026-10-06T02:58:39Z | `cc31526` |
| `37406422380` | ci | push | failure | 2026-10-06T02:54:30Z | `b948bcb` |
| `37405698188` | ci | push | failure | 2026-10-06T02:45:36Z | `ec84cff` |

作业级日志＝`gh run view --job <id> --log`，逐枚重定向后才抽段（抽取件在仓内、全量在 `D:/tmp/wisp274a2/gh/logs/`，只建不删）：

| job | jobId | 日志字节 / 行数 | 步结论（`logs/s1s2-extras.txt`） |
|---|---|---|---|
| `test-windows` | `112547577399` | 974277 / 5610 | `cgo build smoke` **success**；第 7 步 `cmd/wisp CLI tests` **failure**、第 9 步 `Portable windows tests` **failure** |
| `slo-smoke` | `112547577433` | 39615 / 372 | `Build wisp.exe` **success**、全 job success |
| `slo-full` | `112547577499` | 31596 / 290 | `Build wisp.exe (deps cached on the runner)` **success**、全 job success |

抽取件：`logs/s3-runside.txt`（三枚 job 的 node/npm 命中行原文 + runner 身份头）、`logs/s3-slofull-detail.txt`（build.ps1 在作业里的逐行输出 + 步边界 + PATH 线索 + 机侧 node.exe 普查）、`logs/s3-machine.txt`（本机身份与交互式 PATH 读数）、`logs/s3-reds.txt`（今天红哪些枚）、`logs/s3-panel-fails.txt`（panel 那几枚的红因原文）。

### 3.1 ① 产 exe 的 job 里有没有任何一步打印过 `node`/`npm` 的版本或路径

**答：没有一步打印过；就「作业环境里 npm 是否可达」这一问，读数＝取不到（零枚步骤试过硬碰硬的答案）。**
三枚 job 的 `grep -ci 'node'` 分别是 `test-windows` 7 枚、`slo-smoke` 14 枚、`slo-full` 8 枚，`grep -ci 'npm'` 是 2/0/0——**逐枚点名后全部不是「一步跑 node/npm」**（`logs/s3-runside.txt` 有原文行号）：

- Actions **自带运行时**的弃用告示：`(node:12240) [DEP0040] DeprecationWarning: The punycode module is deprecated`（`slo-full` 日志 `112/113/257/258/263/272/273` 行、`slo-smoke` `164/165/182/253/254/256/331/332/337/347/348/350/351` 行、`test-windows` `158/159/176/1678/1679/1681` 行）；
- 各 job **末行**的 runner 级告示（`test-windows` `5610`、`slo-smoke` `372`、`slo-full` `290`）逐字：`Node.js 20 is deprecated. The following actions target Node.js 20 but are being forced to run on Node.js 24: actions/checkout@v4, actions/setup-go@v5, actions/upload-artifact@v4` ⇒ **这就是「作业机器上确有 Node 运行时」的 run 侧证据，但它服务的是 actions 本身，不是 `npm run build`**；
- `test-windows` 那 2 枚 `npm` 命中**不是命令**，是我方 Go 测试的错误文案：`panel_host_gate_test.go:109` 与 `panel_resident_windows_test.go:315` 里那句 `panel: embedded assets are not built (run npm run build in frontend/)`（`2908`、`2990` 行）。

⛔ 按纪律 7，本节不写成「那台机器没有 node」——恰恰相反，见 3.3 的机侧具名读数。

### 3.2 ② hosted 还是 self-hosted（run 侧证据具名，不按 `runs-on` 字符串猜）

- `slo-full` 日志头 `2/3/4` 行逐字：`Runner name: 'wisp-selfhosted-01'`、`Runner group name: 'Default'`、`Machine name: 'DESKTOP-LVS7839'`；`41` 行 `Working directory is 'E:\work\base\actions-runner\_work\wisp\wisp'`；`48` 行 `[command]D:\work\soft\Git\cmd\git.exe config --global --add safe.directory E:\work\base\actions-runner\_work\wisp\wisp`；`45` 行 `Copying 'C:\Users\swq\.gitconfig'` ⇒ **self-hosted，身份三件套齐**。
- `test-windows`（`logs/s3-runside.txt` 头 6 行）与 `slo-smoke`（同）逐字出现 `##[group]Runner Image Provisioner` + `Hosted Compute Agent` + `Version: 20260901.588` + `Build Date: 2026-09-01T19:56:44Z` ⇒ **hosted（GitHub 托管 Windows 镜像）**，两枚都没有 `Runner name:` 这一行（尺本身在 `slo-full` 命中过 1 枚 ⇒ 负向有正控）。

### 3.3 ③ 如果 self-hosted：这台机器与跑 `slo-full` 的是不是同一枚

**是，同一枚＝编排者本机。** 双向对上了：

| 作业侧（run 侧）| 本机现量（`logs/s3-machine.txt`） |
|---|---|
| `Machine name: 'DESKTOP-LVS7839'` | `COMPUTERNAME=DESKTOP-LVS7839` |
| `Copying 'C:\Users\swq\.gitconfig'` | `USERNAME=swq` |
| `Working directory is 'E:\work\base\actions-runner\_work\wisp\wisp'` | `E:/work/base/actions-runner` 存在（`_diag`/`_work`/`bin`/`externals`/`config.cmd`/`run.cmd`），且 `E:/work/base/actions-runner/_work/wisp` 里就是 `wisp` |
| `build.ps1: toolchain: go version go1.27.1 windows/amd64; CC=E:\work\base\msys64\mingw64\bin\gcc.exe (gcc.exe (Rev3, Built by MSYS2 project) 16.2.0)`（`slo-full` 日志 `176` 行） | `E:/work/base/msys64/mingw64/bin` 存在；仓内工作目录 `D:/work/workspace/projects plans/Wisp` 与 `D:\work\soft\Git` 同机 |

⇒ **乙形若落在 `slo-full`，真实代价就是编排者已知的那枚（每次 push 自启、在本机跑、抢 CPU），而且它已经在跑 `build.ps1`；若落在 `test-windows`／`slo-smoke`，落点是 GitHub 托管镜像——那两枚镜像上 node 可达性本腿无 run 侧读数（3.1），要买这一条只能现加一步去问。**

**机侧补充读数（具名标注＝机侧、不是作业侧，不可当作业结论用）**：
- 交互式 PATH 上：`node --version` ⇒ **v24.9.0**、`npm --version` ⇒ **11.6.0**；`cmd.exe //c where node` ⇒ `D:\work\server\node14\node.exe`、`where npm` ⇒ `D:\work\server\node14\npm` + `D:\work\server\node14\npm.cmd`。⚠ 目录名叫 `node14` 而装的是 **v24.9.0**——一处名实不符，具名报回。
- runner 自带：`find E:/work/base/actions-runner -maxdepth 4 -iname 'node.exe'` ⇒ **2 枚**：`externals/node20/bin/node.exe`、`externals/node24/bin/node.exe`（与 3.1 末行「forced to run on Node.js 24」同源）。`runner bin` 目录里 maxdepth 2 无 node ⇒ actions 的 node 走 `externals/`。
- ⚠ **不能由机侧推作业侧**：`ci.yml:804-806` 那条注释逐字写着「The runner's job environment carries no msys64 on PATH … (on the interactive user PATH only)」，而 `E:\work\base\msys64\mingw64\bin` 在本机确实存在（上表已量）⇒ **同机、同目录，交互式 PATH 有而作业 PATH 没有，这枚坑已经在仓里被具名过一次**。`npm` 是否在同一条作业 PATH 上，仍是**取不到**。

### 3.4 顺手量到的两条与本票选形直接相关的 run 侧事实

1. **那句撒谎注释今天真的在 CI 里印出来**：`slo-full` 日志 `183` 行逐字 `build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 搂2.2)`（`搂`＝`§` 在作业控制台代码页下的乱码；源码 `scripts/build.ps1:74` 里是 `§`）⇒ 甲形（改注释）会动到这枚真印出来的字符串，**不是死代码注释**。
2. **出货 exe 不带页面字节，run 侧已自证**：`test-windows` `2908` 行 `AC#12 reading (head cc31526): shape=anchor-only built=false entry-bytes=0 entry-ctype="" entry-err=panel: embedded assets are not built (run npm run build in frontend/) refs=0 check-err=… manifest-entries=0 … | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=0`，且这一枚 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 的结论是 **`--- PASS`**（`2910` 行）⇒ **这就是工单 274 要的「今天拦不住」凭据：尺读到了 `built=false`，然后放它绿。**
3. **丙形的落地面今天不干净**：同一枚 `test-windows` 的 `cmd/wisp` 包已有 8 枚具名红（`logs/s3-reds.txt`：`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`、`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`、`TestRunPacketCarriesTheLoadedInstructionFiles`、`TestPanelHostRealWindowHopAndLifecycle`、`TestPanelHostLatencyPercentilesAC2`、`TestAC14GoSideEvalPushReachesThePage`、`FAIL github.com/CarlosShao/wisp/cmd/wisp 371.253s`；`portable-tests.sh` 四数＝`=== RUN=335 PASS=230 FAIL=7 SKIP=1`）。其中 panel 那两枚的红因是**冷拉起 3414.339 ms 超 D32 预算 1500 ms**（`logs/s3-panel-fails.txt` 原文 `panel_host_windows_test.go:660/665`、`:999`），**与 dist 无关** ⇒ 丙形（把 `wisp panel-assets` 的 rc=1 牙接进 CI）不会撞上这两枚既有红，但要新起一步，不能塞进现有 `cmd/wisp CLI tests` 里冒充绿。

---

## §4 dev 那 64 枚页面源码今天能不能构建（仓外构建，`frontend/**` 零写入）

原始输出：`logs/s4-frontend-build.txt`（dist 清点 + 三份 npm 日志 + `package.json` scripts + `frontend/embed.go` 原文 + `go list -f {{.EmbedFiles}}`）；驱动脚本与全量日志在仓外 `D:/tmp/wisp274a2/`（`fe-run.sh`、`fe-driver.log`、`fe-npmci.log`、`fe-typecheck.log`、`fe-build.log`、拷贝树 `fe/frontend/**`）——**只建不删，全在那里**。

**做法逐字**（⛔ 没在仓内 `frontend/` 落过一个字节）：

```
mkdir -p /d/tmp/wisp274a2/fe && git archive HEAD frontend | tar -x -C /d/tmp/wisp274a2/fe
cd /d/tmp/wisp274a2/fe/frontend && (npm ci || npm install) && npm run build ; echo "rc=$?"
```

`archive_rc=0`，导出后 `fe/frontend` 顶层＝`VENDORED.md dist embed.go fixtures index.html package-lock.json package.json scripts src tsconfig.app.json tsconfig.json tsconfig.node.json vite.config.ts`（＝跟踪的 85 枚，其中 `src` 64 枚，与 §0 起手锚一致）。

### 4.1 三档读数（逐档 rc 与耗时，尺＝驱动脚本里的 `date +%s` 差）

| 档 | 命令 | rc | 秒 | 结果 |
|---|---|---|---|---|
| 装依赖 | `npm ci`（走 lockfile，一次成，未落 `npm install` 退路） | **0** | **10** | `node_modules` 顶层 52 枚包目录；`npm audit` 报 **1 high severity vulnerability**（只报不修） |
| 只 `tsc -b` | `npm run typecheck`（`package.json` 里 `typecheck: tsc -b`） | **0** | **7** | **`tsc -b` 那一半不红**，零诊断输出 ⇒ 乙形不欠类型账 |
| 完整构建 | `npm run build`（`tsc -b && vite build`） | **0** | **9** | `vite v8.3.0`，`✓ 2439 modules transformed`、`✓ built in 801ms` |

**总耗时 26 秒**（`ALL_DONE ci=0 tc=0 build=0 total_seconds=26`）。

### 4.2 产物：出得出 `dist/index.html`

尺＝`find …/dist -type f -printf '%p %s bytes\n'` ⇒ **4 枚 / 共 602635 字节**：

| 文件 | 字节 |
|---|---|
| `dist/index.html` | **1068**（vite 自报 `1.06 kB │ gzip 0.64 kB`） |
| `dist/assets/index-vdBrT8rM.js` | **552027**（`552.02 kB │ gzip 170.74 kB`，带 `>500 kB` 分块告警） |
| `dist/assets/index-yy8KMgdf.css` | **49540**（`49.54 kB │ gzip 9.84 kB`） |
| `dist/.gitkeep` | 0（导出的锚文件，构建不清空它） |

⇒ **答：能。dev 这一支的页面源码今天是自洽可构建的，`go build` 之前只差「谁去跑这一句 `npm run build`」。**

### 4.3 一处必须具名的现场差异（本机 dist 里有前人产物，本腿没碰）

尺＝`go list -f '{{.EmbedFiles}}' ./frontend`（在仓内只读跑，不写）⇒ **本机工作树**的嵌入集是 4 枚：`dist/.gitkeep`、`dist/assets/index-BRKj5OIJ.css`、`dist/assets/index-BVKlegVD.js`、`dist/index.html`；而 `git ls-files frontend/dist` 只有 `.gitkeep` 一枚 ⇒ 工作树里躺着**别的会话未跟踪的构建产物**（工单 274 里那句「本机 `frontend/dist` 那 4 枚文件」，本腿现量对上）。
⚠ 关键差异：**哈希名不同**（本机 `index-BRKj5OIJ.css`/`index-BVKlegVD.js` vs 本腿仓外新建 `index-yy8KMgdf.css`/`index-vdBrT8rM.js`）⇒ 本机那批是**旧源码的过期产物**，不是当前 HEAD 的产物。⇒ **「本机 `go build` 出来的 exe 带页面」这一条今天能绿，是这堆未跟踪文件在替它绿**；CI 的纯净检出（`test-windows` 日志 `2908` 行 `shape=anchor-only`）证实 CI 上没有它们。⛔ 本腿没删、没盖、没动这 4 枚，一个字节都没有。

### 4.4 承重原文：`frontend/embed.go` 自己就把乙形写成了唯一通道

`frontend/embed.go`（原文在 `logs/s4-frontend-build.txt`）逐字：「Ticket 77 AC#1: `npm run build` in this directory produces dist/, and the go:embed below is the **ONLY** way those bytes reach wisp.exe. There is no runtime CDN and no local HTTP server (D29 …), so **the binary must be self-sufficient on a machine with no node, no npm and no network**」＋ `//go:embed all:dist` ＋ `const AnchorName = ".gitkeep"`（注释解释：只有锚文件时 `panel.Assets.Built()` 报 false，宿主显示 "assets not built"，永远不会悄悄发占位页）。
⇒ 这条设计口径把 §3 的结论咬死了：**exe 带页面＝出货机上必须有谁跑过 `npm run build`**，而那三枚产 exe 的 job 今天谁都没跑（§2/§3）。
