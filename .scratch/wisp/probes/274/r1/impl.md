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
只能靠自觉，⛔ 不得据此以为写了也没关系。这一处是**射程缺口的事实登记**，不是本腿的请求裁定（改射程属别票）。

---
