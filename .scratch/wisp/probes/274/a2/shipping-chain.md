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
