# 274 — 写腿 `274-r1` 实现件（落形＝乙-1，编排者已裁）

票面＝`.scratch/wisp/issues/274-no-nail-requires-the-shipped-exe-to-carry-the-page-build-ps-step-2-comment-is-false-in-both-halves.md`
（**已整份读完 98 行**，非只读派单转述）。
唯一写面＝`scripts/build.ps1` 一枚文件。本件与 logs 属取证盘，不属产码面。

---

## 0. 起手锚（第 1 笔 commit，先于一切长跑命令）

原文见 `logs/00-anchor.txt`、`logs/01-toolchain.txt`、`logs/02-dist-baseline.txt`。

| 尺 | 现读数 |
|---|---|
| `date` | `Wed Oct  7 11:46:08 CST 2026` |
| `git rev-parse --short HEAD` | `8976ffab`（⚠ 票面立票锚点是 `ca4ae2af`，起手已在其后） |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal scripts tools .github docs frontend` | **空**（这七枚路径起手干净；工作树里别人的改动在 `design/**`、`.gitignore`、`.scratch/**`，不在本 pathspec 内，见 §7） |
| `git ls-files frontend/dist` 枚数 | **1**＝`frontend/dist/.gitkeep`（锚文件，与票 §2③ 一致） |
| `node --version` | `v24.9.0` |
| `npm --version` | `11.6.0` |
| `Get-Command node,npm`（PowerShell 实际解析） | `node.exe` → `D:\work\server\node14\node.exe`；`npm.ps1` → `D:\work\server\node14\npm.ps1` |
| `Get-Command go.exe,npm.cmd,gcc.exe` | `go.exe` → `D:\work\base\go\bin\go.exe`；`npm.cmd` → `D:\work\server\node14\npm.cmd`；`gcc.exe` → `E:\work\base\msys64\mingw64\bin\gcc.exe` |

★起手一处**必须在落码前定调的仪器事实**：PowerShell 里 `npm` 解析到的是 **`npm.ps1`**（npm 官方 PowerShell shim），
不是 `.cmd`/`.exe`。`build.ps1` 顶部是 `$ErrorActionPreference = 'Stop'` + `Set-StrictMode -Version 2.0`，
经 `.ps1` shim 调 npm 会把 shim 自身的 stderr/异常语义混进构建步骤。⇒ 本腿解析顺序＝**先 `npm.cmd`（批处理垫片，
语义最扁），取不到再退回 `Get-Command npm`**；两者都取不到＝**具名报错＋非零退出**（不许跳过）。

### `sed -n '60,80p' scripts/build.ps1` 逐字（改前）

```
$ccDir = Split-Path -Parent $cc
if (-not ($env:PATH -like "*$ccDir*")) { $env:PATH = "$ccDir;$env:PATH" }
$ccVersionLine = (& $cc --version | Select-Object -First 1)
Write-Host "build.ps1: toolchain: $($goVersionLine); CC=$cc ($ccVersionLine)"

# --- 1. fetch-deps ----------------------------------------------------------
& (Join-Path $PSScriptRoot 'fetch-deps.ps1') -RepoRoot $RepoRoot
if ($LASTEXITCODE -ne 0) { Fail 'fetch-deps.ps1 failed (see output above).' }

# --- 2. frontend ------------------------------------------------------------
Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'

# --- 3. go build ------------------------------------------------------------
function Get-DepsValue([string]$Section, [string]$Key) {
    $inSection = $false
    foreach ($raw in (Get-Content -LiteralPath (Join-Path $RepoRoot 'deps.toml'))) {
        $line = $raw.Trim()
        if ($line.StartsWith('[')) {
            $inSection = ($line.Trim('[', ']').Trim() -eq $Section)
            continue
        }
```

`scripts/build.ps1` 起手 MD5＝`9dee32a3c73317a94df456f5208cf1cc`，总 172 行，**UTF-8 无 BOM**（首三字节 `23 52 65`＝`#Re`）。
非 ASCII 行现量 **5 枚**，全部只是 `§`（U+00A7）：`:4 :14 :74 :103 :136`。本机 Git Bash 与 `powershell` 两侧读出来都**没有乱码**
⇒ 派单 §1.4 预警的「GBK 误读使 `§` 变乱码」现象**今天本机未复现**，如实登记为「未观察到」，⛔ 不因此改整文件编码。
⇒ 新增文案一律 ASCII 英文（派单硬要求），`§` 那 5 枚原有字符保留原样。

### 起手 dist 基线（AC#9 的凭据，`logs/02-dist-baseline.txt`）

已**先备份后动**：`cp -a frontend/dist/. /d/tmp/wisp274r1/dist-stale-backup/`（临时件只建不删）。

```
frontend/dist/.gitkeep                    md5=d41d8cd98f00b204e9800998ecf8427e  sha256[0:8]=e3b0c44298fc1c14  bytes=0
frontend/dist/assets/index-BRKj5OIJ.css   md5=db4db7a27fb34e3de7acf045561d753a  sha256[0:8]=f8ff337b2ce80f3b  bytes=49943
frontend/dist/assets/index-BVKlegVD.js    md5=70128a3dfdbc73e9f603bb421cdf324c  sha256[0:8]=793e606366463f0e  bytes=553469
frontend/dist/index.html                  md5=f98bfc4ba5dd39d74c9b2f85438de70a  sha256[0:8]=2078af7f2343af94  bytes=1044
```

★与票 §4（`274-a2`）第 4 节现量一致：本机这 4 枚＝陈旧件（`index-BRKj5OIJ.css`／`index-BVKlegVD.js`），
HEAD 仓外新建那形＝`index-vdBrT8rM.js`／`index-yy8KMgdf.css` ⇒ **哈希名全不同**。本腿因此把「exe 带页面」的读数
**全部**做成「fresh build 文件名＋sha256 指纹 ↔ exe `-manifest` 指纹逐枚相同」这一形（AC#9），枚数级读数不作证。

`build/wisp.exe` 起手＝`31076405` 字节、md5 见 log、mtime `2026-10-06 16:57`（别人的产物）⇒ 已另存
`/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`，收尾还原（`build/` 是 gitignored，不进 porcelain，但仍按「不是我动的不毁」处置）。

### 现读到的两把可用的尺（决定 §4/AC#9 怎么量）

- `wisp panel-assets -manifest`（`cmd/wisp/panel_assets.go:38` 定义、`:130-136` 走 `assets.Manifest()`）
  ⇒ `internal/panel/assets.go:100-122` 每行打 `"<path> <size> <sha256 前 8 字节的 hex>"`，**这就是 exe 内嵌物的逐枚指纹**，
  可直接与盘上 fresh dist 的 `sha256[0:8]` 逐枚对（`sha256sum | cut -c1-16` 同形）。
- `wisp panel-assets`（不带子命令）＝票 §2② 说的那把现成的牙：`cmd/wisp/panel_assets.go:125-129`
  ⇒ `if !assets.Built() { "wisp panel-assets: assets NOT BUILT"; return 1 }`。本票乙-1 的第 3 小建（构建后校验）
  就是要在**产线侧**接上这根线，让「npm 跑成功但什么也没产出」也红。

### d22scan 射程（派单 §1.4 要我现量并报告的那一条）

`scripts/d22scan.sh` → `tools/d22scan` 的扫描根（现读 `tools/d22scan/main.go`）＝
`internal/`（walkGo）、`cmd/`（walkGo）、`frontend/`（walkText，panel-approval / ban #6）、`internal/tools/`（walkText，ban #7），
外加 allowlist `tools/d22scan/allowlist.txt`。⇒ **`scripts/` 不在 d22scan 射程内**（零命中 `filepath.Join(root, "scripts")`）。
实际后果：本腿在 `scripts/build.ps1` 里写什么形状的 emoji，**仪器今天看不见**；因此「UI 零 emoji」那条规矩在这一枚文件上
实际后果：本腿在 `scripts/build.ps1` 里写什么形状的 emoji，**仪器今天看不见**；因此「UI 零 emoji」那条规矩在这一枚文件上
只能靠自觉，⛔ 不得据此以为写了也没关系。这一处是**射程缺口的事实登记**，不是本腿的请求裁定（改射程属别票）。
★§5 里我补了仪器自己的射程打印（比读码更强）：ban #8 实扫 `design/` 39 枚、`frontend/` 85 枚、`internal/` 514 枚、`cmd/` 104 枚，
**名下列表里没有 `scripts/`** ⇒ 上面那句由「读码推得」升为「仪器自报」。

---

## 1. 未修码读数＝**今天拦不住**那一发（`logs/03`、`logs/04`、`logs/05`）

改前的第 2 步零命令 ⇒ 我不需要从 CI 借这发读数，直接在**本机把树摆成干净检出的那一形**就能演：
`frontend/dist` 只剩跟踪锚文件 `.gitkeep`（⛔ 那 3 枚陈旧件不是删除、是 `mv` 到仓外 staging，§6 有逐枚还原自证），
然后按 CI 逐字同一形调用跑真 `build.ps1`（`ci.yml:571/:757/:820` 三处都是这一形：
`powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`，cwd＝仓根）。

现场（`logs/04-unfixed-buildps1-full.txt`，全文 13 行）：

```
build.ps1: repo root D:\work\workspace\projects plans\Wisp (Env=dev)
build.ps1: toolchain: go version go1.27.1 windows/amd64; CC=E:\work\base\msys64\mingw64\bin\gcc.exe (gcc.exe (Rev3, Built by MSYS2 project) 16.2.0)
fetch-deps: cache hit - third_party/sherpa-onnx matches deps.toml (sherpa-onnx 1.13.8)
build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)
build.ps1: go build ok (cgo linked against sherpa-onnx C API)
build.ps1: DLLs colocated into D:\work\workspace\projects plans\Wisp\build
build.ps1: signed model manifest colocated into D:\work\workspace\projects plans\Wisp\build\models
build.ps1: wrote D:\work\workspace\projects plans\Wisp\build\SHA256SUMS
build.ps1: smoke test - running wisp.exe doctor
build.ps1: done. Artifacts in build\: wisp.exe, onnxruntime.dll, ... SHA256SUMS
BUILD_PS1_PROCESS_EXITCODE=0
```

那枚 exe（`logs/05`，md5 `50846822b62fd3426ced8365408891cb`，`30471221` 字节）自己怎么说的，三个分支都跑了：

```
wisp panel-assets                 -> wisp panel-assets: assets NOT BUILT                        EXITCODE=1
wisp panel-assets -manifest       -> wisp panel-assets: panel: embedded assets are not built (run npm run build in frontend/)   EXITCODE=1
wisp panel-assets -check          -> 同上具名错误                                                EXITCODE=1
```

** ⇒ 这就是票面第 1 节那句「没有任何判据要求它有」的现场，两半都取到了：
**出货通路整条绿（`EXITCODE=0`、走到 smoke test、写了 SHA256SUMS），而交出去的那枚 exe 一个页面字节都不带、
它自己的牙（`rc=1`＋具名 `NOT BUILT`）**明明会喊**，但没有任何脚本或 CI 步骤去咬它。
另附一把改前/改后都成立的尺：`go list -f '{{.EmbedFiles}}' ./frontend` 在锚文件形下回 `[dist/.gitkeep]` 且 `rc=0`
（`logs/03`、`logs/12`）⇒ 「能不能编译」与「里面有没有页面」在这一形下确实毫无关系（票 §2③ 在本机真树上复现，不止仓外台件）。
⚠ 我没有把 `§` 那句谎话改口之前先跑过别的分支：**这一发的红不在 build.ps1 里，在 exe 里**，
所以第 2 步本身零退出＝票面立案的那句谎话在日志里的原样（`logs/04` 第 5 行）。

## 2. 落地内容＝乙-1（唯一写面 `scripts/build.ps1`；逐字 diff 见 `git diff -- scripts/build.ps1`，`logs/` 未单列因 diff 已进 commit）

文件 172 → **265 行**；`Parser::ParseFile` 报 `PARSE_ERROR_COUNT=0`；仍 UTF-8 **无 BOM**（首三字节 `35,82,101`＝`#Re`）；
行尾仍**全 LF**（`grep -c $'\r'` ＝ 0）；非 ASCII 行 5 → **4**（`4,18,196,229`，全是原有那枚 `§`：`:14→:18`、`:103→:196`、`:136→:229` 是我插入行的位移，
原 `:74` 那句谎话随替换**消失**）⇒ **我没有新增任何非 ASCII 字符**，也没动整文件编码。

第 2 步现在做四件事，每一件都只有「成功」和「具名失败」两形，⛔ 无跳过分支：

1. **树对不对先看一眼**：`frontend/package.json` 不存在 ⇒ `Fail`（页面源码不在树上时，后面每一步都是空话）。
2. **解析 node 与 npm**：`Resolve-BuildTool` 先查带扩展名的 `node.exe` / `npm.cmd`，再退无扩展名 `node` / `npm`；
   任一取不到 ⇒ `Fail`，错误文本**指名试过哪两个名字、以及为什么不跳过**（「The page bundle is a build input of wisp.exe,
   so this script fails instead of skipping it」）。为什么优先 `npm.cmd`：起手锚量到 PowerShell 把裸 `npm` 解析成 **`npm.ps1`** 垫片，
   而本文件顶部是 `$ErrorActionPreference='Stop'`＋`Set-StrictMode 2.0` ⇒ 走 `.ps1` 会把垫片自己的异常语义混进构建步，
   我要让这一步唯一的失败来源是 **npm 自己的退出码**。两枚工具的**实际解析路径＋版本**都印进日志（`node = ... (v24.9.0); npm = ... (11.6.0)`），
   这样 CI 日志第一次有了「这台作业机到底有没有 node」的可读证据（票 §7 那枚取不到的数，落地后可由 run 日志直接取）。
3. **真构建**：`node_modules` **存在与否**决定走哪一支，**两支都印出走了哪一支**（派单硬要求）——
   存在 ⇒ 明说 `taking the npm-ci-SKIPPED branch` ＋「这一支不校验已装树与 package-lock.json 是否相符，删掉 node_modules 可强制走 npm ci」；
   缺失 ⇒ 明说 `running npm ci (clean install pinned by package-lock.json)`。两支之后一律 `npm run build`（它自己就是 `tsc -b && vite build`，
   现读 `frontend/package.json:8`）。`npm ci`／`npm run build` 任一步非零 ⇒ `Fail` 带退出码原文＋npm 路径。
4. **构建后校验（本票真正的牙）**：`npm run build` 退 0 **不算证据**。四道各自具名的检查，顺序是
   ①`dist` 目录在不在 → ②**除锚文件外有没有任何一枚**（`$distReal`，锚文件＝`.gitkeep`）→
   ③`index.html` 在不在（`panel.EntryFile` 那一枚）→ ④`index.html` 字节数是否 > 0。
   ★为什么第 ② 道用「除锚文件之外的枚数」而不是「总枚数 > 1」：票 §5 AC#4 要钉的是「只剩锚文件那一形必须红」，
   而 vite 是会清空输出目录的——若某次构建清掉了 `.gitkeep` 却真产出了页面，按总枚数算会把它误判红，
   那是**为错的原因报红**。按「锚文件之外至少一枚」算，锚文件形照旧红（0 枚），真页面不会误红。
   全过时印出 `frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\...` 这一行**名册＋字节**。

⛔ 射程守住了：`.github/workflows/ci.yml` 一个字节没动（零新 job、零 artifact 传递、无 `continue-on-error`）、
`internal/panel/assets.go` 的 `errNotBuilt`／`Built()`／`Resolve` 一字未改、`frontend/.gitignore` 的 `dist/*` 未动、
`frontend/**` 未写一个字节、CI 未加 `-tags winlive`、`docs/PLAN.md`／`docs/specs/**` 未碰（AC#8 属文字面，另案）。

★**一处我主动扩了半格，具名报回，请编排者裁**：票面/派单只点名 `:74` 那一句，但同一文件**第 10 行的 `.DESCRIPTION`**
写的是 `2. frontend: skipped (none yet, lands in S5)`——同一句谎话的第二处、且是我刚实现的那一步的自我描述。
落实第 2 步之后留着它，就是**本腿自己新造的一句假话**（比过期更糟），故一并改成说实话（ASCII，含「No node/npm on this machine means this script FAILS; the step never skips」）。
它仍在唯一写面 `scripts/build.ps1` 内，⛔ 不是别的文件；若编排者认为 AC#3 的边界只容 `:74`，这一处请退回，我另笔改回。

## 3. 两发定向突变（票面 AC#4；两发都真跑，退出码原文与具名红因整行都在）

### (i) PATH 剥掉 node（`logs/09-mutation-i-rerun-after-node-flag-fix.txt`＝终码版；`logs/07`＝首版）

台件＝仓外 `D:\tmp\wisp274r1\mut-i-no-node.ps1`：只在**该进程**的 PATH 里滤掉含 `node14` 的那一枚目录，
⛔ 系统/用户 PATH 一字未写（跑完即随进程消失，无需还原；日志第 6 行自证 `remaining PATH entries matching node14: 0`）。

```
node.exe resolves to:            (空)
node resolves to:                (空)
npm.cmd resolves to:             (空)
npm resolves to:                 (空)
go.exe still resolves to: D:\work\base\go\bin\go.exe        <- 只剥 node，go/gcc 还在
build.ps1: FATAL: frontend step: node could not be resolved on PATH (tried node.exe, then node). The page bundle is a build input of wisp.exe, so this script fails instead of skipping it. Install the toolchain pinned in docs/BUILD.md, or run the build on a machine that has node and npm on PATH.
BUILD_PS1_PROCESS_EXITCODE=1
```

⇒ **非零退出 1＋具名红因**，且红在第 2 步、在 `go build` 之前（日志里没有 `go build ok` 那行）。

### (ii) npm 假装成功而 dist 只剩锚文件（`logs/10-mutation-ii-rerun-after-node-flag-fix.txt`＝终码版；`logs/08`＝首版）

台件＝仓外 `D:\tmp\wisp274r1\fake-npm-shim\npm.cmd`（`exit /b 0`、什么都不产出）＋前置到该进程 PATH；
`frontend/dist` 临时清成只剩 `.gitkeep`（3 枚陈旧件 `mv` 到仓外 staging，非删除）。node 保持**真的**在（这一发的题就是「npm 说成功了怎么办」）。

```
npm.cmd resolves to: D:\tmp\wisp274r1\fake-npm-shim\npm.cmd   <- 台架确实被解析到了
node.exe resolves to: D:\work\server\node14\node.exe          <- node 是真件
build.ps1: frontend: node = D:\work\server\node14\node.exe (v24.9.0); npm = D:\tmp\wisp274r1\fake-npm-shim\npm.cmd (fake-npm 0.0.0 ...)
build.ps1: frontend: node_modules EXISTS -> taking the npm-ci-SKIPPED branch ...
build.ps1: frontend: running npm run build (tsc -b && vite build, output goes to frontend/dist).
fake-npm 0.0.0 (ticket 274 mutation rig - exits 0 without doing any work)
build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds only the anchor file(s) [1 file(s), none beside frontend/dist/.gitkeep]. That is the clean-checkout shape, so the exe would ship with zero page bytes. Named cause: build produced no page artifacts.
BUILD_PS1_PROCESS_EXITCODE=1
--- dist after run ---
D:\work\workspace\projects plans\Wisp\frontend\dist\.gitkeep   <- 桩 npm 什么都没产出，dist 仍是锚文件形
```

⇒ 「npm 跑成功但什么也没产出」这一形**红**，红因具名到「锚文件形／零页面字节」。
这正是票 §2④ 说的 `AC#12` 那把尺两形都放行之缺的补牙：它不看你有没有 built 标志，它看**除锚文件外有没有东西**。

★两发之间的**一个我自己的缺陷**，如实登记（不是隐藏）：首版 `logs/07/08` 里 `node = ... ()` 是空的，日志同时冒出一句
`Error: Cannot find module 'D:\...\Wisp\version'`——我把 node 的版本探测写成了 `& $node.Source version`（那是 `go` 的形状，`node version` 会把 `version` 当脚本路径）。
⇒ 改成 `--version` 后**两发都重跑**（`logs/09`、`logs/10` 即终码版读数，上面贴的都是终码版），
`node = ... (v24.9.0)` 现在有值、那句假错误也随消失。**首版两份 log 我没有删**（临时件只建不删），它们是「这个缺陷怎么被抓到的」的记录。

## 4. fresh build ↔ exe 内嵌物逐枚对照（票面 AC#9，`logs/12`、`logs/13`、`logs/14`）

先把树摆成**干净检出那一形**（`logs/12`：`dist` 只剩 `.gitkeep`、`frontend/node_modules` 整棵 `mv` 到仓外 ⇒ 强制走 `npm ci` 支），
再跑终码 `build.ps1` 一发到底（`logs/13`，`BUILD_PS1_PROCESS_EXITCODE=0`）：

```
build.ps1: frontend: node_modules ABSENT -> running npm ci (clean install pinned by package-lock.json).
build.ps1: frontend: running npm run build (tsc -b && vite build, output goes to frontend/dist).
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\index-B8yINMF1.js=551989 assets\index-yy8KMgdf.css=49540 index.html=1044
build.ps1: go build ok (cgo linked against sherpa-onnx C API)
```

**① fresh build 的名册＋哈希（盘上，同一发产出）** ＝ **② exe 自己报的内嵌名册**（`wisp panel-assets -manifest`，
即 `internal/panel/assets.go:100-122` 从 embed FS 现算的 sha256 前 8 字节）：

| 名 | 盘上 bytes | 盘上 sha256[0:8] | exe 内 bytes | exe 内 sha256[0:8] | 逐枚 |
|---|---|---|---|---|---|
| `.gitkeep` | 0 | `e3b0c44298fc1c14` | 0 | `e3b0c44298fc1c14` | 同 |
| `assets/index-B8yINMF1.js` | 551989 | `8f06145dac449fb6` | 551989 | `8f06145dac449fb6` | 同 |
| `assets/index-yy8KMgdf.css` | 49540 | `d0b198664b250973` | 49540 | `d0b198664b250973` | 同 |
| `index.html` | 1044 | `9b7856b949d63989` | 1044 | `9b7856b949d63989` | 同 |

旁证三把（都出自同一发）：`go list -f '{{.EmbedFiles}}' ./frontend` ＝
`[dist/.gitkeep dist/assets/index-B8yINMF1.js dist/assets/index-yy8KMgdf.css dist/index.html]`；
`wisp panel-assets` ＝ `panel assets embedded: 4 files, entry=index.html built=true`（`EXITCODE=0`）；
`wisp panel-assets -check` ＝ `entry=index.html built=true 2 asset refs resolve [./assets/index-B8yINMF1.js ./assets/index-yy8KMgdf.css]`
⇒ 入口 HTML 里引用的两枚资产名**就是** exe 里那两枚，不是「枚数凑巧对得上」。
这枚 exe ＝ md5 `e6c8e52bb14f15a7983bbc6b093e6058`、`31074357` 字节。

★**陈旧件那坑被正面钉住**：§0 基线里本机陈旧件是 `index-BRKj5OIJ.css`/`index-BVKlegVD.js`，
本次 fresh build 是 `index-B8yINMF1.js`/`index-yy8KMgdf.css`——**名与哈希全不同**，
而 exe 报的是**新**那一组 ⇒ 「4 files, built=true」这一形单看确实能替旧字节报假绿，只有逐枚哈希对拉才能分辨（AC#9 的立格理由，本腿实测复现）。

⚠ **一发未决现读，具名报回（新发现，⛔ 不是我漏做的格）**：票 §5 记的 `274-a2` 仓外 fresh build 是
`index-vdBrT8rM.js`／`index-yy8KMgdf.css`、`index.html` 1,068 字节；我这一发在仓内是 `index-B8yINMF1.js`／`index-yy8KMgdf.css`、`index.html` **1,044** 字节。
⇒ **CSS 名相同、JS 名与 index.html 字节不同**＝同一枚 HEAD 的页面源码，在「仓外 `git archive` 导出树」与「仓内工作树」两处构建出的 bundle **不是同一份字节**。
我做了有界排查，四项都排掉了：`frontend/vite.config.ts` 里 `../`／`process.env`／`loadEnv`／`publicDir`／`envDir`／`define:` **零命中**；
`frontend/.env*` 与仓根 `.env*` **不存在**；`git status --porcelain --ignored -- frontend/src frontend/index.html` **零命中**（无被忽略的本地生成物混进构建）；
`git log -1 -- frontend` ＝ `611ae8b9`（2026-09-26，已确证是 HEAD 祖先）⇒ 两次构建之间**跟踪的页面源码没有变动**。
⇒ 成因**未定**，本腿不再深挖（再挖要读 `frontend/scripts/*.mjs` 与 vite 解析细节，越出本票射程）。
它的意义朝本票有利方向：**「同一份源码」并不自动等于「同一份 bundle 字节」，构建上下文本身也在决定出货字节**——
票 §1 那句「带不带取决于编译机器上碰巧有什么」比立案时更成立。对 AC#9 本身无影响（该格要的是 ①与② 同发逐枚相同，已证）。

## 5. 门禁四数（票面 AC#6／派单 §4，全部本腿现跑，`logs/16`、`logs/17`）

| 尺 | 读数 | rc |
|---|---|---|
| `sh scripts/d22scan.sh` | 正控先绿：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` ⇒ 后面的 clean 不是假绿；扫描 `d22scan: clean - no D22 ban violations` | **0** |
| `gofmt -l`／gofumpt v0.12.0 `-l` | **本腿零 `.go` 改动 ⇒ 无 `.go` 可喂，未喂**（尺＝`git status --porcelain -- '*.go'` 里 ` M ` 形态 0 枚）。⛔ 我没有把 `.md`/`.ps1` 喂给 gofmt（那会出 `U+0023`＋rc=2＝工具误用） | 不适用（如实记） |
| `go vet ./cmd/wisp/ ./internal/panel/` | 输出为空 | **0** |
| 终态 `git status --porcelain -- cmd internal scripts tools .github docs frontend` | ` M scripts/build.ps1`（⛔ 与起手「空」差这一枚，那是我的写面本体；除此之外**零差异**，`frontend` 下零枚＝陈旧 dist 已原样还原） | 0 |

仪器自报的射程（从 `logs/16` 原文，比读码强）：ban #8 实扫 `design/` 39＋`frontend/` 85＋`internal/` 514＋`cmd/` 104，
**名下列表无 `scripts/`** ⇒ 印证 §0 那句「`scripts/build.ps1` 不在 d22scan 射程」。
另记一条仪器行为：本轮扫描时它报 `d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`
⇒ dist 产出被它按 gitignore 跳过，本票的构建产出不会污染扫描面。

## 6. 还原自证（三次动 dist／一次动 node_modules，每一次都有逐枚 md5 对拉）

| 时刻 | 动作 | 还原证据 | 结果 |
|---|---|---|---|
| §1 未修码读数前 | `mv` 3 枚陈旧件到仓外 staging（⛔ 非删除） | `logs/06-restore-proof-run1.txt` | `diff -r rc=0 IDENTICAL`，4 枚 md5 全 MATCH |
| §3(ii) 突变前 | 同一批 `mv` 到同一路径 | `logs/11-restore-proof-after-mutations.txt` | `diff -r rc=0 IDENTICAL`，4 枚全 MATCH |
| §4 e2e 前 | 陈旧件 `mv` 出＋`frontend/node_modules` 整棵 `mv` 到仓外 | `logs/15-restore-proof-after-e2e.txt` | `diff -r rc=0 IDENTICAL`；`node_modules` 原件移回原位，`npm ci` 重建的那棵移存 `/d/tmp/wisp274r1/node-modules-generated-by-npm-ci/` |

终态复核（同一份 `logs/15`）：`go list -f '{{.EmbedFiles}}' ./frontend` 回到陈旧那一形
`[dist/.gitkeep dist/assets/index-BRKj5OIJ.css dist/assets/index-BVKlegVD.js dist/index.html]`
⇒ 本机那份陈旧件**字节未被我改动**（票 AC#9 那句「⛔ 不许为了让读数好看去动本机陈旧件」我守住了：
我全程只**移出再移回**，从未改过其中任何一枚，fresh build 的产出另存仓外）。
fresh build 的 3 枚产物也**没删**（`/d/tmp/wisp274r1/dist-fresh-274r1/`），`build/wisp.exe` 起手那份同样另存在
`/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`（`31076405` 字节／md5 `7032d94d...`）——**我没有把它还原回去**：
它是被本票合法重跑的出货通路产物，而 §4 那枚 `e6c8e52b...` 才是 AC#9 读数的主体；还原它会让「exe 里是哪一份」变成日志与盘上不一致。
这一处置由我具名登记，见 §7 第 4 条，请编排者裁要不要我再移回去。

## 7. 未做完／待定格的格，逐枚具名（⛔ 无空话，⛔ 没测到的不写成通过）

1. **票面 AC#0 的 (a) 半句未做**：`lint-frontend` 那一发真跑成功后 `npm run build` 在 **CI 侧**产出几枚、落在哪——
   本腿取到的是**本机同形调用**的枚数与名册（4 枚／名＋字节＋sha256 指纹，`logs/13/14`），⛔ 那不是 CI 的 run 侧读数。
   CI 侧那半句要推送后由编排者取全量 run 日志（本腿不推送）。
2. **票面 AC#0 的 (b) 半句已在本机做成，但形不同**：票面问的是「同一发 run 里 `build.ps1` 造出的 exe 报 true 还是 false」。
   本腿给的是「同一发**本机** run 里 exe 报 `built=true` 4 枚并与 fresh build 逐枚同哈希」（§4）＋改前那一形报 `NOT BUILT rc=1`（§1）。
   ⇒ **CI run 侧那一枚未做**（同第 1 条，归编排者推送后取数）。
3. **票面 AC#1 未做**（⛔ 不是我漏做，是本腿射程里没有它）：它要求把「只剩锚文件那一形」喂给 `AC#12` 那把尺（`go test -overlay` 或仓外合成 bundle），
   证今天整包全绿，另需一发正向对照。本腿是**写腿**，派单 §4 只给我四数、⛔ 明令我不要跑整包（本票排程 §6 的包级互斥理由），
   且派单改让我走「真 `build.ps1` 出货通路」那一形来证「今天拦不住」（§1）。
   那把尺的两形读数我**只从代码现读**（`cmd/wisp/panel_host_gate_test.go:77` 及 §2④ 那些行号是票面现读，我没复跑），⛔ 没跑过 `go test` 一个用例。
4. **票面 AC#2 的三形代价表不属本腿**（编排者已裁＝乙-1，本腿不重开）；但落地后**新增一条可现量的代价**：
   `build.ps1` 现在把 Node 变成 Windows 产线的硬依赖，本机实测第 2 步耗时＝`npm ci`＋`npm run build`（日志时序见 `logs/13`）；
   CI 时长增量本腿取不到（不推送、不跑 run）。另：`build/wisp.exe` 被本票重跑覆盖（起手那份已另存，见 §6），要不要还原请编排者裁。
5. **票面 AC#3**：`:74` 那句谎话已改成实话（且随批改了文件头 `.DESCRIPTION` 同义那半句，见 §2 末具名条目）。
   ⚠ 派单/票面写的示例文案「这一步在这里跳过，出货判据在 AC#2 所选形里」在乙-1 下**不适用**（这一步现在不跳过），
   我按派单 §1.4 的实际要求写「现在构建前端；没有 node 就整条失败」。**这一处按原文与派单的分工处理，如有越界请退回。**
6. **票面 AC#5 守住**：零新 `C##`、`C17` 白名单零新名、D29/D41/`SPEC-11` 文字零改动。
   `SPEC-11 §2.2` 那句「embed lands in S5」／`--with-frontend` 开关（AC#8）**我一行没碰**——改它＝人工批准，另案。
7. **票面 AC#6 的「终态逐名红名册与起手作差＝0」未做**：本腿**没有起手整包 `go test` 基线**（派单 §4 明令把四数收窄为
   d22scan／gofmt-gofumpt（无 `.go` 可喂）／`go vet ./cmd/wisp/`／终态 porcelain，且票面 §6 的包级互斥理由是针对本波的）。
   ⇒ 我**没有**跑 `go test ./cmd/wisp ./internal/... -count=1`，因此**既无起手名册也无终态名册，作差不成立**；
   这一格**如实记为未做**，⛔ 不写成通过。要补的话请编排者排一枚「写腿跑整包」的时机（与 272/275 波的包级互斥一并裁）。
   已跑的替代读数：`go vet ./cmd/wisp/ ./internal/panel/` rc=0（编译面干净）＋d22scan 正控先绿再 clean。
8. **票面 AC#7 的 `-overlay` 半句未用**：本腿的 dist 突变全走 `mv`（移出／移回，逐枚 md5 对拉，§6），⛔ 全程未删、未改任何一枚陈旧件字节；
   `git status` 终态与起手差 `M scripts/build.ps1` 一枚（我的写面本体）。**这一条我按派单 §3/§4 的形做，未按票面 AC#7 的 `go test -overlay` 做**，具名报回。
9. **票面 AC#9 已做**（§4），含票面点名的反向坑：本机陈旧件我只移出移回，未改一字。
10. **AC#4 的两枚定向突变都已真跑并贴退出码原文与整行红因**（§3）。⛔ 我没有把「exe 侧那根现成的牙」(`cmd/wisp/panel_assets.go:125-129`)
    接进 `build.ps1`：派单 §1 的第 3 小建明确把校验定义在 **dist 侧**（`go build` 在第 3 步、exe 在第 2 步之后才存在），
    而「出货前跑一次 `wisp.exe panel-assets` 逼它咬自己」那一根线**今天仍然没有任何脚本或 CI 步骤去接**——
    ⇒ 票 §2② 末句「现成的牙没人咬」这一半**本票状态未变**，属编排者裁（要么并入丙形另案，要么另开票）。
11. **known cost 已在现场复现形状**：§3(i) 就是 `test-windows`／`slo-smoke` 作业 PATH 若没有 node 时会红的**那一行原文**
    （`build.ps1: FATAL: frontend step: node could not be resolved on PATH ...`＋`EXITCODE=1`）。
    本腿不推送、不判 CI 结果；红名册请编排者取数后按票 111 AC#2 裁例逐条登记（⛔ 不撤步骤、不改产码、不用「跳过」消红）。
12. **别人的未提交改动：我没提交也没还原**（如实登记）：`design/**` 起手即脏（终读 `git status --porcelain -- design` 见
    ` D design/assets/base.css`、` D design/assets/icons.js`、` D design/assets/theme.js`、` D design/assets/tokens.css`、` M design/doubao/README.md`），
    `.gitignore`、`.scratch/**` 若干（含 `.scratch/wisp/probes/*.go` 未跟踪件一大列）同样不是我动的。
    本腿的 `git add`/`git commit` 全程带显式 pathspec，只提 `.scratch/wisp/probes/274/r1/**` 与 `scripts/build.ps1`。
13. **日志里含一枚 U+2713**：`logs/13` 捕到 vite 自己那句 `✓ built in 484ms`。那是**被捕获的工具输出**，⛔ 不是本腿写的界面文案，
    且 d22scan 射程不含 `.scratch/`（§0/§5）。我**没有**为好看而抹它——改日志原文比留着它更糟。

---

## 8. commit 名册与终态自查（`logs/18-final-endstate-and-roster.txt`）

三笔，全程 `git add` 与 `git commit` 都带显式 pathspec 字面量，⛔ 零 `add -A`／零裸 commit／⛔ 零 push：

| 笔 | 时刻 | 内容 | `git show --stat` 的文件面 |
|---|---|---|---|
| `77e278bf` | 11:50 | 起手锚（§0） | `impl.md`＋`logs/00,01,02`＝4 files, 182 insertions |
| `35633445` | 12:02 | 产码（乙-1，§2） | **1 file changed, 95 insertions(+), 2 deletions(-)＝只有 `scripts/build.ps1`** |
| `58f1df8e` | 12:03 | 证据（§1–§7） | `impl.md`＋`logs/03..17`＝16 files, 784 insertions |

★每一笔落地前后我都查了 staged 面：起手那笔之前 `git diff --cached --name-only` ＝ 空；
产码那笔之前同样只有我自己的 `scripts/build.ps1` 进入提交（`git show --stat` 原文已贴在本节与 commit 里）
⇒ 没有把别人已 staged 的文件一起提走（今天 `8d30a862`／`ca4ae2af` 那两发的事故形，本腿零复发）。

终态（全部现读，`logs/18`）：

```
git status --porcelain -- cmd internal scripts tools .github docs frontend   ->  空（＝起手那一条，我的写面已入库）
frontend/dist  4 枚 md5 逐枚 = 起手基线；diff -r frontend/dist /d/tmp/wisp274r1/dist-stale-backup -> rc=0 IDENTICAL
scripts/build.ps1  md5=644c2f6af1fa3fc7f8b3ab38c20753f5  lines=265  非 ASCII 行数=4（全为原有 §，无新增）
frontend/node_modules  present（原树移回原位；npm ci 重建的那棵另存仓外，未删）
```

仓外留存的临时件（⛔ 一枚未删，`issues/README` 规则 8）：
`dist-stale-backup/`（改前 dist 原件）、`dist-stale-staging/`（移出用的空壳）、`dist-fresh-274r1/`（本次 fresh build 的 3 枚产物）、
`node-modules-pre-274r1/`（起手那棵 node_modules 的中转）、`node-modules-generated-by-npm-ci/`（`npm ci` 重建的那棵）、
`wisp.exe.pre-274r1-backup`（起手 exe）、`fake-npm-shim/npm.cmd`（突变 (ii) 台架）、
`mut-i-no-node.ps1`／`mut-ii-stub-npm.ps1`／`run-build.ps1`／`parse-check.ps1`（台件）。


