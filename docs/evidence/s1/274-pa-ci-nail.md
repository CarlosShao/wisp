# 274 · 「出货 exe 真带着页面包」在 CI 里钉一道直量（写腿 `274-paci`；2026-11-14 14:1x +08）

**腿性质**：写面三处 —— `.github/workflows/ci.yml`（+81 行、零删改）、本件（新）、`.scratch/wisp/probes/274/paci/**`（读数）；票面只追加一节。⛔ 零 push、⛔ 零 AC 框翻动、⛔ `scripts/build.ps1`／产码／测试／规格／阈值／golden 一字节未动。
**起手锚**：`86e5655ba`（branch `dev`）。本笔提交前 HEAD ＝ `2dfd29b98`（他腿同时在上游 commit），但 `git rev-parse HEAD:.github/workflows/ci.yml` 起手＝提交前＝**`631e9bd92`** ⇒ 本腿的 ci.yml diff **只有自己的 +81 行**，没有夹带、也没有覆盖谁的改动。
**一句话**：新步骤是 `slo-smoke` 的第 **5** 步（0-based 索引 4／共 7 步），紧接 `Build wisp.exe` 之后，对出炉的 `build\wisp.exe` 跑 **`panel-assets -check`**，退出码用 `Process.Start` + `WaitForExit` 取（⛔ 不是 `&`）。

---

## 0. 结论表（一行一格，读数出处指到文件）

| 问 | 答 | 读数／出处 |
|---|---|---|
| 加在第几步后面 | `slo-smoke` 第 **5** 步（0-based 4），**紧接 `Build wisp.exe` 之后**、`SLO smoke gate` 之前 | §4（`yaml.safe_load` 抽出的步名册＋`PREVIOUS_STEP`/`NEXT_STEP`） |
| 跑的是什么命令 | `build\wisp.exe panel-assets -check`（对**出货 exe**，不是对 `frontend/dist`、不是 `go test`） | §1.1 |
| 负控①（`frontend/dist` 只剩锚文件） | **rc=1**，逐字 `panel: embedded assets are not built (run npm run build in frontend/)` | §3.3 / `probes/274/paci/neg1-extracted.txt` |
| 负控②（入口在、引用缺，本腿加发） | **rc=1**，逐字点名 `index.html references "./assets/index-QUtCIbHm.js" but the embedded bundle cannot serve it` | §3.4 / `neg2-entry-only.txt` |
| 正控（还原成完整包） | **rc=0**，逐字 `panel assets check: entry=index.html built=true 2 asset refs resolve [./assets/index-QUtCIbHm.js ./assets/index-CZoXrb9N.css]` | §3.5 / `pos-extracted-final.txt` |
| YAML 仍合法 | **rc=0**（PyYAML 6.0.3 `yaml.safe_load`，两次） | §4 / `yaml-and-extract.txt` |
| 我跑的是**提交的那串字节** | 从改后的 ci.yml 抽出 `run:` → sha256 `52109f41…`，负控／正控都用这串跑 | §3.7 |
| 代价 | 检查本身 **83／101／283／319 ms**（四次裸跑）；CI 里没有额外 job／artifact／`needs:` | §5 |
| push 了没有 | **没有**。推不推由编排者核后做 | §7 |

---

## 1. 落点与形状

### 1.1 `git diff .github/workflows/ci.yml` 逐行（全文，无省略）

```diff
diff --git a/.github/workflows/ci.yml b/.github/workflows/ci.yml
index 631e9bd92..8b8cda531 100644
--- a/.github/workflows/ci.yml
+++ b/.github/workflows/ci.yml
@@ -1407,6 +1407,87 @@ jobs:
         if: ${{ !cancelled() }}
         run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev
 
+      - name: Shipped wisp.exe carries the page bundle (panel-assets -check on the artifact)
+        # Ticket 274: "the exe that ships carries the page" had no direct
+        # measurement in CI. The AC#12 line is an indirect witness - built=true
+        # read off a TEST binary built in the same run from the same
+        # frontend/dist - and the only other sighting of `panel-assets` in this
+        # file was a comment (lint-frontend's L2 render step quotes its output
+        # as fixture input). A comment is not a step and cannot go red. This
+        # step measures the artifact itself: the build\wisp.exe the step above
+        # just produced, asked to serve its own embedded bundle.
+        #
+        # What `-check` proves (cmd/wisp/panel_assets.go:114-124 ->
+        # internal/panel/assets.go:137-160): the entry file resolves AND every
+        # asset reference inside it resolves from that same embedded tree. So
+        # the shape "the embed shipped less than npm run build produced" -
+        # built=true, blank panel, every hashed asset 404 - exits non-zero
+        # here, and the offending reference is named on the line above.
+        #
+        # Its reach, stated so nobody reads more into it: this says the exe
+        # carries a COMPLETE bundle, not that the bundle is THIS checkout's
+        # bytes (that is build.ps1's roster/provenance gate, ticket 274 AC#10,
+        # in the step above). Both sides were measured on 2026-11-14 in an
+        # exported tree outside the repo, readings in
+        # docs/evidence/s1/274-pa-ci-nail.md: frontend/dist reduced to the
+        # anchor-only shape -> exit 1, "embedded assets are not built (run npm
+        # run build in frontend/)"; entry present with its hashed refs absent
+        # -> exit 1, `index.html references "./assets/index-QUtCIbHm.js" but
+        # the embedded bundle cannot serve it`; dist restored to the full
+        # bundle -> exit 0, "entry=index.html built=true 2 asset refs
+        # resolve". Three bare runs of the check on the same bench the same
+        # hour: 319 / 283 / 101 ms, the first two carrying the cold DLL load.
+        #
+        # Ordering: directly below the build, so a pageless exe is red before
+        # the SLO gate below spends its whole sampling run on it.
+        if: ${{ !cancelled() }}
+        shell: pwsh
+        run: |
+          $exe = Join-Path $PWD 'build/wisp.exe'
+          if (-not (Test-Path -LiteralPath $exe)) {
+            Write-Host "the build step above reported success but $exe does not exist - nothing was measured, so this step does not pass"
+            exit 1
+          }
+          # `& $exe panel-assets -check` is NOT the shape used here. build\wisp.exe
+          # links with -H=windowsgui (scripts/build.ps1:244-257), and PowerShell does
+          # not wait for a GUI-subsystem child: measured on this bench 2026-11-14,
+          # `& $exe panel-assets -check *> file` returned after 20 ms with an EMPTY
+          # file and $LASTEXITCODE still reading 99 from the statement before it.
+          # The measured shape is the one build.ps1's own smoke test uses
+          # (scripts/build.ps1:307-361, ticket 244-r2): Process.Start with
+          # UseShellExecute=$false, so the child's stdout still reaches this log,
+          # then WaitForExit and the real ExitCode.
+          $psi = New-Object System.Diagnostics.ProcessStartInfo
+          $psi.FileName = $exe
+          $psi.Arguments = 'panel-assets -check'
+          $psi.UseShellExecute = $false
+          $proc = $null
+          try {
+            $proc = [System.Diagnostics.Process]::Start($psi)
+          } catch {
+            Write-Host "could not start ${exe}: $($_.Exception.Message) - no verdict was produced, so this step does not pass"
+            exit 1
+          }
+          if ($null -eq $proc) {
+            Write-Host "no process object came back for ${exe} - no verdict was produced, so this step does not pass"
+            exit 1
+          }
+          $code = $null
+          try {
+            $proc.WaitForExit()
+            $code = [int]$proc.ExitCode
+          } catch {
+            Write-Host "${exe} started but its exit code could not be read: $($_.Exception.Message) - no verdict was produced, so this step does not pass"
+            exit 1
+          } finally {
+            $proc.Dispose()
+          }
+          Write-Host "wisp.exe panel-assets -check -> exit code $code"
+          if ($code -ne 0) {
+            Write-Host "the shipping executable cannot serve its own embedded page bundle (panel-assets -check exit code $code); the lines above name the missing piece"
+            exit 1
+          }
+
       - name: SLO smoke gate
```

**只加不减**：`git diff --numstat` ＝ `81 0 .github/workflows/ci.yml`；既有七步一步未挪、未改名、未加 `if:`／`continue-on-error`、未动 `timeout-minutes`、未动 `env:`。

### 1.2 照的是邻居的形状（不是新创）

`slo-smoke` 原有的每步都是 `- name:` ＋（必要时一串 `#` 注释）＋ `if: ${{ !cancelled() }}` ＋ `run:`，**没有** `shell:`／`working-directory:`／`env:`（job 级 `env: WISP_ENV=test` 照旧继承）。新步沿用同一骨架。三处**有意的差别**，逐条给理由：

| 差别 | 为什么 |
|---|---|
| 多了 `shell: pwsh` | 邻居的 `run:` 是 `powershell … -File scripts/xxx.ps1`（把正文放在脚本里）。本腿写面**不含 `scripts/**`**（⛔ 不许新建脚本），正文只能内联；内联时把解释器点名，等于把"跑它的是哪一个 PowerShell"从默认值变成明写的字。仓里 ubuntu 的每一步都显式写 `shell: bash`（`ci.yml:955` 等 8 处），同一个习惯 |
| `run:` 用块标量 `\|` | 内联多行 PowerShell 只能是块标量。文件里现成的块标量先例＝`:306`／`:510`／`:772`／`:1618` |
| `name:` 里带括号说明 | 同文件既有习惯（`(AC#3 render evidence)`、`(tickets 18/20; …)`、`(deps.toml key)`） |

### 1.3 为什么"放得进 `slo-smoke`"＝三条前提逐条现量（所以没走"停手问你"那一支）

派单留了一个分叉：若那步的 dist 是 `build.ps1` 产的、路径不同，就**别硬塞**、写出两个候选落点停手问。本腿逐条量了三条前提，**都成立**：

1. **dist 在同一个 job 里已生成**：`scripts/build.ps1:86-215` 是前端段（`npm ci` 当 `node_modules` 不在时 ＋ `npm run build`），`ci.yml:1408` 那一步就是跑这个脚本；`build.ps1:215` 打印本次产出的 dist 名册。⇒ 不存在"dist 在另一个 job"的问题（`lint-frontend` 那份 dist 确实在另一个 job、且不传给任何人，但**本步不依赖它**）。
2. **exe 在同一个 job 的固定路径**：`build.ps1:270` → `build\wisp.exe`；该 job 没有 `defaults.run.working-directory`，cwd＝仓根。⇒ `Join-Path $PWD 'build/wisp.exe'` 指的就是那枚出炉件。
3. **命令真存在、且正是"对 exe 量页面包"那一支**：`cmd/wisp/panel_assets.go:36-138`，`-check` 分支在 `:114-124`。⛔ 不是猜的：先 `git grep -n "panel-assets" -- .github/ scripts/ cmd/` 看真用法（`.github/` 唯一命中＝`:1587` 的**注释**；`cmd/wisp/main.go:181` 是分派；`scripts/` 零命中）。

---

## 2. 为什么是 `-check`（默认支／`-manifest` 都不行）

`cmd/wisp/panel_assets.go` 的四支能力对照（读码 ＋ 本腿现跑）：

| 支 | 干什么 | 对"出货 exe 带页面"够不够 |
|---|---|---|
| 无参（默认支，`:125-136`） | 只 `assets.Built()`：入口文件在就绿，然后打 `N files … built=true` | **不够**：这正是 `internal/panel/assets.go:128-136` 注释点名的形状——入口在、页面里引的哈希件不在，`built=true` 而面板全白 |
| `-manifest`（`:105-113`） | 逐枚列出**内嵌物**的名字／字节数／sha256 前 8 字节 | 是"带的是哪些字节"的尺，但它**不判**引用能不能解析；本腿用它当旁证（§3.5 末） |
| **`-check`（`:114-124`）** | `assets.Check()`：入口能解析 **且** 入口里每一个 `src=`/`href=` 引用都能从同一棵内嵌树解析 | **正是这一问**：`internal/panel/assets.go:137-160`，解析不到就 `return served, fmt.Errorf("panel: %s references %q but the embedded bundle cannot serve it")` → rc=1 |
| `-render`／`-l2` | 单字节输出／L2 卡 JSON | 与本问无关 |

⇒ 用的是 `-check`，代价是"入口 + 2 枚引用"这一条链全通才绿，比 `built=true` 严一格。

---

## 3. 正控：在**仓外导出树**里做三发（判据＝弄坏必红、还原必绿）

### 3.1 树怎么来的（可逐条重跑）

```
D:\tmp\wisp274paci\tree
  ← git archive HEAD | tar -x -C D:\tmp\wisp274paci\tree
  ← New-Item -ItemType Junction tree\third_party → <repo>\third_party   （third_party 是 .gitignore:12，archive 里没有）
  ← 同置 DLL：third_party\sherpa-onnx\*.dll → tree\build\（照 build.ps1:277-280）
```

- **⛔ 没有碰仓里那份 `frontend/dist`**：它是别人的在飞产物（`git status` 里不脏，但本腿不改它一字节）。三发全部在 `D:\tmp\wisp274paci\` 下动。
- `git archive` 现跑 **rc=1**，逐字原因是 tar 对 **16 行**异常路径名的报错（全部在 `.scratch/` 下、名字带非 ASCII 或引号，`Invalid empty pathname`），末行 `tar.exe: Error exit delayed from previous errors.`；`build` 相关的那棵树照样齐（`tree\frontend\dist\.gitkeep`、`scripts/`、`cmd/`、`internal/`、`models/`、`go.mod` 都在）。逐字 stderr 存 `probes/274/paci/export-tree-anomalies.txt`。
- **为什么不在树里跑 `scripts/build.ps1`**：它的第 2 步（`:86-215`）会**自己跑 npm 把 dist 填满**，且第 2 步的出处闸门（`:196-213`）要求 dist 里出现"本次 npm 产的新名字"——两件事都会把我要的"锚形 dist"抹掉。⇒ 用**同一条 go build 口径**在树里造 exe，照抄 `build.ps1:244-270`：`CGO_ENABLED=1`、`CC=E:\work\base\msys64\mingw64\bin\gcc.exe`（bin 前置进 PATH）、`go build -trimpath -ldflags "… -H=windowsgui" -o build\wisp.exe ./cmd/wisp`。**`-H=windowsgui` 一字不改地保留**，因为那条属性正是"`&` 不等它"的来源，去掉就不算复现。

### 3.2 ⚠ 先记一发失手（不粉饰）：第一次跑出来的是**错的红**

树里第一次跑（dist 已锚形化）得 `rc=-1073741515`＝`0xC0000135`＝`STATUS_DLL_NOT_FOUND`——因为树里 `build\` 只有 exe，没有 `onnxruntime.dll`／`sherpa-onnx-*.dll`（`build.ps1:277-290` 才同置）。**那不是本检查的红**，是缺 DLL 的红。照 `build.ps1` 补同置后立刻变成真读数。⇒ 这枚读数留在 `neg1-anchor-only.txt` 的头一段，任何人复跑前先看它，别把"环境缺件"记成"页面缺件"。

### 3.3 负控①：`frontend/dist` 只剩锚文件（票 274 已有此形）

树里把 `frontend/dist/index.html` 与 `assets/` **移走**（移到 `D:\tmp\wisp274paci\dist-pos-staging\`，⛔ 没删）⇒ dist 只剩 `.gitkeep`。然后 `go list` 与重建 exe：

```
=== go list EmbedFiles with anchor-only dist ===
[dist/.gitkeep]
go_list_rc=0
go_build_rc=0  exe_sha256=a0634c07701ef782717e719f4c76951e3a03e7f87119185ac6030793be041355
```

**用提交的那串字节**（§3.7）跑，cwd＝树根：

```
=== NEG-1: bytes extracted from the committed ci.yml, cwd=tree, dist = anchor-only ===
STEP| time=2026-11-14T14:17:36.138+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
STEP| wisp panel-assets: panel: embedded assets are not built (run npm run build in frontend/)
STEP| wisp.exe panel-assets -check -> exit code 1
STEP| the shipping executable cannot serve its own embedded page bundle (panel-assets -check exit code 1); the lines above name the missing piece
STEP_EXITCODE=1
=== same exe, direct CLI (for the 'what does the tool say' column) ===
DIRECT| wisp panel-assets: panel: embedded assets are not built (run npm run build in frontend/)
direct_rc=1
```

⇒ **rc≠0 成立**，且点名的是缺的那件事（"embedded assets are not built (run npm run build in frontend/)"＝入口 `index.html` 不在，`errNotBuilt` 在 `internal/panel/assets.go:31-34`）。全文 `probes/274/paci/neg1-extracted.txt`。

### 3.4 负控②（本腿加发，为把"点名缺什么"钉到**具体那一枚引用**）

锚形只会说"没构建"，不会说"少了哪个文件"。把 `index.html` 放回、`assets/` 仍不在 ⇒ 重建 exe ⇒ 同一串字节跑：

```
=== NEG-2: entry present, hashed refs absent; committed step body, cwd=tree ===
STEP| wisp panel-assets: panel: index.html references "./assets/index-QUtCIbHm.js" but the embedded bundle cannot serve it
STEP| wisp.exe panel-assets -check -> exit code 1
STEP| the shipping executable cannot serve its own embedded page bundle (panel-assets -check exit code 1); the lines above name the missing piece
STEP_EXITCODE=1
```

⇒ 这一发**逐字点名**少的是 `./assets/index-QUtCIbHm.js`。全文 `neg2-entry-only.txt`。（这一发用的 `step-body.ps1` 与提交串在**注释换行**上差 2 字节，代码行逐字节相同；见 §3.7 的对照。）

### 3.5 正控：`frontend/dist` 还原成完整包 ⇒ 必须 rc=0

把 staging 里的 `index.html` ＋ `assets/` 移回（＝还原），重建 exe：

```
=== tree dist RESTORED ===
.gitkeep 0 sha256[0:16]=e3b0c44298fc1c14
index.html 1044 sha256[0:16]=664520b87528f607
assets\index-CZoXrb9N.css 65264 sha256[0:16]=6d443e15fb042128
assets\index-QUtCIbHm.js 640318 sha256[0:16]=b990516742de811d
[dist/.gitkeep dist/assets/index-CZoXrb9N.css dist/assets/index-QUtCIbHm.js dist/index.html]
go_list_rc=0
go_build_rc=0
exe bytes=51387675 sha256=b5584a149c9df60e6a8cb5ba2f8d173ec84396a241ec4b5214c2cd8246311995
=== POS (final): dist restored, exe rebuilt, committed bytes extracted from ci.yml, cwd=tree ===
STEP| time=2026-11-14T14:18:44.071+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
STEP| panel assets check: entry=index.html built=true 2 asset refs resolve [./assets/index-QUtCIbHm.js ./assets/index-CZoXrb9N.css]
STEP| wisp.exe panel-assets -check -> exit code 0
STEP_EXITCODE=0
=== same exe: panel-assets -manifest (what the artifact carries) ===
MANIFEST| .gitkeep 0 e3b0c44298fc1c14
MANIFEST| assets/index-CZoXrb9N.css 65264 6d443e15fb042128
MANIFEST| assets/index-QUtCIbHm.js 640318 b990516742de811d
MANIFEST| index.html 1044 664520b87528f607
```

⇒ **还原后 rc=0 成立**；右侧 `-manifest` 四行与盘上 dist 四枚的 sha256 前 8 字节**逐枚同值**（这一行只是旁证"这枚 exe 带的就是这棵树"，⛔ 不是本步的判据）。全文 `pos-extracted-final.txt`。

### 3.6 为什么这步不用 `&`（本机两发实测，两种写法**一发等一发不等**）

`build\wisp.exe` 是 `-H=windowsgui`（`build.ps1:244-257`，现量 PE subsystem＝**2**；树里那枚也是 2）。同一枚 exe、同一条命令、只换调用写法（`trap-ampersand.txt`）：

| 写法 | 读数 |
|---|---|
| `& $exe panel-assets -check 2>&1 \| ForEach-Object …` | 等到进程结束（111 ms），`rc_read_back=0`，输出完整 |
| `& $exe panel-assets -check *> file` | **17 ms 就返回**，`rc_read_back=99`（＝上一句留下的陈旧值），重定向文件 **0 字节；800 ms 后仍是 0 字节** |
| `Process.Start` + `WaitForExit()`（**提交的形状**） | 83 ms，`rc=0`，输出完整 |

⇒ 不是"PowerShell 一定不等"，而是**取决于写法与句柄**——这种"有时对、有时悄悄错"最难查。仓里已经为这件事挨过一次（`build.ps1:307-318` 逐字记着 ticket 263／244-r2：`&` 21 ms 返回、`$LASTEXITCODE` 读到 `go build` 留下的陈旧 0，于是 smoke 闸门**结构性地恒绿**）。⇒ 本步直接用那条量过的形状，并在 `run:` 里留注释说明，免得下一个人"化简"回 `&`。

### 3.7 我跑的是**提交的那串字节**（不是"我另写了一份像它的东西"）

1. 改后的 ci.yml 用 `yaml.safe_load` 读出来，按步名取那一步的 `run` 字段，原样落成 `probes/274/paci/extracted-run.ps1`（44 行／1846 字符）。
2. `sha256(extracted-run.ps1) = 52109f410e83fffa22afe085f566bd53166e73cbae40591736a692841991f25f`；**改完注释后又抽了一次，同值**（⇒ 我后面写进 `#` 里的 319/283/101 ms 没有动过 `run:` 一个字节）。
3. §3.3 负控与 §3.5 正控**都用这枚文件跑**，cwd＝树根。⇒ 那两发红／绿，就是推下去之后 CI 里会发生的两发。

> ⚠ 一处口径：这枚 `.ps1` 在 `.gitattributes` 里是 `*.ps1 text eol=crlf`，**将来 checkout 出来的副本会是 CRLF，哈希必然变**。判据是 ci.yml 里 `run:` 字段的正文，⛔ 不是这枚文件的哈希；哈希只是"我当时跑的就是当时文件里那串"的封条。

4. **提交之后，又对 `HEAD` 的那枚 blob 复验了一次**（判据升级：GitHub 读的是对象，⛔ 不是我工作区那份）：`git show HEAD:.github/workflows/ci.yml` → `yaml.safe_load` **OK**、`steps=7`、新步索引 **4**、`prev='Build wisp.exe'`／`next='SLO smoke gate'`、blob 全文 sha256 **`c9cc3b46…`**；从**blob**里再抽一次 `run:` → sha256 仍是 **`52109f41…`**，与工作区抽出的那枚 **`identical=True`**。全文 `post-commit-blob-check.txt`，被解析的那枚对象留档为 `ci-committed-blob.yml`。

---

## 4. YAML 合法性与"步骤真在我说的地方"

仓里既有手段里**没有** YAML 解析器（`Get-Module -ListAvailable powershell-yaml` 空、`yq` 空；`scripts/**` 里引 ci.yml 的都是 `grep`/`awk` 级的字面尺）。故按派单许可走 **PyYAML**（DSH 自带 python：`…\dependencies\python\python.exe`，`yaml 6.0.3`）：

```
YAML_PARSE_OK: yaml.safe_load returned dict
JOBS: lint, test-core, test-windows, slo-smoke, slo-full, lint-frontend
SLO_SMOKE_STEPS:
  0: actions/checkout@v4
  1: actions/setup-go@v5
  2: Cache third_party (deps.toml key)
  3: Build wisp.exe
  4: Shipped wisp.exe carries the page bundle (panel-assets -check on the artifact)
  5: SLO smoke gate
  6: Upload SLO report
STEP_INDEX=4 (0-based) of 7 step(s)
PREVIOUS_STEP='Build wisp.exe'
NEXT_STEP='SLO smoke gate'
STEP_KEYS: ['if', 'name', 'run', 'shell']
STEP_IF: ${{ !cancelled() }}
STEP_SHELL: pwsh
RUN_BLOCK_LINES: 44 chars: 1846
YAML_CHECK_RC=0
```

第二发（同一份脚本再跑一次，用于"改注释后再验"）：

```
YAML_PARSE_OK: yaml.safe_load returned dict
…
RUN_BLOCK_LINES: 44 chars: 1846
YAML_CHECK_RC=0
extracted sha256=52109f410e83fffa22afe085f566bd53166e73cbae40591736a692841991f25f
```

另加一发结构性抽查（同一解析器）：`steps=7`、`runs-on=windows-latest`、`timeout-minutes=20`、`rc-ok`，`rc=0`。全文 `probes/274/paci/yaml-and-extract.txt`；脚本 `yaml_check_and_extract.py`。

**诚实边界**：YAML 合法 ≠ GitHub 会接受它。真正的裁决只能靠一次真跑（⛔ 本腿不 push）。本腿能给的替代证据是：`shell`／`if`／`run` 三个键的形状与同文件其余 **59** 步（`yaml.safe_load` 现数：6 job／60 步）同构、`run:` 正文在本机以同一种 shell（`pwsh`）真跑过三发。

---

## 5. 代价

- **CI 增量**：一枚 exe 起停（无新 job、无新 artifact、无 `needs:`、无 npm 重跑）。检查本身四次裸跑（`Process.Start`→`WaitForExit`）：**319 / 283 / 101 ms**（树里那枚，前两发含冷 DLL 载入）与 **83 ms**（仓里那枚，已热身）。口径与票 274 早前那条"150 ms"同量级。
- 本机 `pos-extracted-final` 那一发的 2034 ms **不是**这一步的代价：那是我为了"跑提交的那串字节"多起了一层 `pwsh -NoProfile -File`（约 1.8 s 是那层壳）。CI 里这一步本来就在 pwsh 进程内，只剩 exe 自己那 0.1–0.3 s。
- 相对于它下游的 `SLO smoke gate`（`-SecondsPerState 4`，多状态采样，分钟级），这 0.1–0.3 s 可以忽略；放在 SLO 之前，是让"exe 都不带页面"这种红**先**红，不白烧一轮采样。

---

## 6. 射程：这一步**不**证明什么；以及两个我考虑过的备选落点

### 6.1 它不证明的两件事（写清，避免被读成"全绿"）

1. **不证明"带的是这份源码的字节"**。`-check` 只判内嵌树**内部完整**。本机现成的反例就在盘上：仓里那枚 `build\wisp.exe`（sha256 `ff941871…`，2026-11-02）跑 `-check` 回 **rc=0**，而它带的引用是 `./assets/index-DpRdJ-7r.js`／`./assets/index-Dxt2y5oW.css`，**与工作树 `frontend/dist` 里的 `index-QUtCIbHm.js` 不是同一套名字** ⇒ 这枚 exe 带的**不是**当前源码的产物，而 `-check` 照样绿。这一格归 `scripts/build.ps1:196-213` 的出处闸门（票 274 AC#10），**已经**在 CI 的同一 job 里跑；新步与它互补，不互相替代。
2. **不证明宿主机上真能开出面板**（WebView2 宿主是票 33/35）。它证明的是"字节在 exe 里、且这条链能自洽解析"。

### 6.2 备选落点（都没选，代价一栏留着给编排者复议）

| 落点 | 代价 | 为什么没选 |
|---|---|---|
| **A（本腿选的）`slo-smoke` 的 `Build wisp.exe` 之后** | exe 起停 0.1–0.3 s／次 push；0 新 job／artifact／`needs:` | ①dist 与 exe 都在同一个 job 里、路径固定；②hosted runner 是每 push 都跑的合入门；③紧挨产出它那一步，"谁先红"最清楚 |
| **B `slo-full`（自托管机）的 `Build wisp.exe (deps cached on the runner)` 之后** | 同样 0.1–0.3 s，但跑在另一台机（另一套 node/gcc/缓存态）⇒ **多一台机器的独立性** | 没选：派单只点名 `slo-smoke`；且 `slo-full` 的红今天已承载"machine-contended＝无结论"这层语义（`ci.yml:1421-1448`），再往它身上叠红面会把"红＝什么"搅浑。**这是一条正当的后续收紧**，不是本腿否掉的 |
| **C `lint-frontend`（注释 `:1587` 所在的那个 job）** | 要新 job／跨 job 传 dist／新 artifact ⇒ 正是票 274 AC#2 的**乙形**，且与 `build.ps1:196-213` 出处闸门**正冲** | 没选：那个 job 是 node-only `ubuntu-latest`，**根本没有 exe**；它能给的只有"dist 里有 3 枚文件"，量不到"出货 exe 带不带页面"——这正是本票要补的那一格 |

撤销口令：**「274-paci 撤那一步」**＝删掉 ci.yml 里 `Shipped wisp.exe carries the page bundle` 整步（`git diff` 里那 81 行）即可回到 `631e9bd92` 的形状，其余文件无需动。

---

## 7. 纪律自证

- **写面**：`.github/workflows/ci.yml`（+81/−0）· 本件（新）· `.scratch/wisp/probes/274/paci/**`（15 枚，含两枚 **0 字节**文件——它们正是 §3.6 那行"重定向文件 0 字节"的物证，⛔ 不许当垃圾清掉）· 票 274 只追加一节。`scripts/**`／`cmd/**`／`internal/**`／`frontend/**`／`docs/specs/**`／阈值／golden／测试：**本腿没有任何一笔写入**（可核＝本笔 commit 的 `--name-only` 只列本句前三处＋票面）。
- **仓里那份 `frontend/dist` 一字节未动**：三发全部在 `D:\tmp\wisp274paci\` 的导出树里；仓里 `build/wisp.exe` 也一字节未动（它只是被**读**来当反例，§6.1）。
- **git**：只 `commit` 带**显式 pathspec**；⛔ 无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；⛔ **零 push**。索引里他腿暂存的那一枚（`.scratch/wisp/issues/70-ci-actually-green.md`）本腿没碰、也没被卷进本笔 commit（`--name-only` 会给出）。
- **commit**：件里不写死本笔 hash（写死了，下一次改动就变成假话）——用 `git log -- docs/evidence/s1/274-pa-ci-nail.md` 取；票面追加的那一节同法可追。

### 7.1 本腿在盘上留下的四处可核事实

1. `git diff --stat 631e9bd92 -- .github/workflows/ci.yml` ＝ `81 +++`。
2. `git rev-parse HEAD:.github/workflows/ci.yml` 提交后应等于本笔新 blob（不是 `631e9bd92`）。
3. 票面新追加的那一节，`git log -- .scratch/wisp/issues/274-*.md` 可追。
4. 15 枚文件在 `.scratch/wisp/probes/274/paci/`，其中三枚（`neg1-extracted.txt`／`neg2-entry-only.txt`／`pos-extracted-final.txt`）是"提交字节"的红／红／绿。

---

## 8. 复跑清单（谁都能重来一遍）

```powershell
# 1) YAML 合法 + 步骤位置 + 抽出提交字节
& <dsh-python> .scratch\wisp\probes\274\paci\yaml_check_and_extract.py   # rc 应为 0

# 2) 仓外导出树（锚形）
git archive HEAD | tar -x -C D:\tmp\wisp274paci\tree
New-Item -ItemType Junction -Path D:\tmp\wisp274paci\tree\third_party -Target <repo>\third_party
#    ⚠ 一定先同置 DLL，否则得到的是 0xC0000135 的假红（§3.2）
Copy-Item <repo>\third_party\sherpa-onnx\*.dll D:\tmp\wisp274paci\tree\build -Force

# 3) 同一条 go build 口径（照 build.ps1:244-270），dist 只剩 .gitkeep 时：
#    CGO_ENABLED=1、CC=<mingw gcc>、go build -trimpath -ldflags "... -H=windowsgui" -o build\wisp.exe ./cmd/wisp
# 4) 跑提交的字节：cwd = 树根
pwsh -NoProfile -File <repo>\.scratch\wisp\probes\274\paci\extracted-run.ps1   # 期望 exit 1
# 5) 还原 dist（index.html + assets/）→ 重建 → 再跑 ⇒ 期望 exit 0
```
