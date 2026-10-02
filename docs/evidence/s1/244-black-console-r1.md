# 244-r1 证据件 —— 双击启动黑控制台：`scripts/build.ps1` 的 `$ldflags` 追加 `-H=windowsgui`（甲形）

## §0 起手锚

- 时刻：起手 **2026-10-02 19:29:16 +0800**（`date` 现取）；交件锚 **2026-10-02 19:37:52 +0800**。
- 锚点：进场 `git rev-parse --short HEAD` ＝ **`61c526e0`**，分支 **`dev`**（写面 diff 也对这枚锚负责，diff 基线 `3593e87b`）。本腿全部读数对这枚锚负责。
- 派单：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（⛔ AC 框一枚没碰，归编排者）。
- 修法凭据：普查件 `.scratch/wisp/probes/244/c1/census.md` §3.1 形状甲；台账 **`A546` 裁甲**（build.ps1 唯一权威出口，CI 三处 step 与 slo-check 自建分支自动跟随；乙＝第二真相源被毙、丙＝FreeConsole 被 G4 裁撞 `SPEC-11:50` 字面 flag 要求）。撤销口令：「244 改丙」。
- 规格锚：`docs/specs/SPEC-11-build-deploy-containerization.md:50` 逐字「CLI 与 GUI 同一二进制：无参 = GUI（`-H=windowsgui`）」。
- 纪律自证：写面只动了 `scripts/build.ps1` 一枚（＋证据件＋台件区 `.scratch/probes-244-r1/`）；`docs/BUILD.md`、`.github/workflows/ci.yml`、`docs/PLAN.md`、`docs/specs/**`、SLO/golden、三枚冻结件、`frontend/**`／`design/**` 零读零写；⛔ 没跑 `go test` 整包（同机他腿共用）；未跑 `go build ./cmd/wisp` 单独验编译——全链真跑 build.ps1（其内即 go build ./cmd/wisp）已覆盖且更强。

## §1 改法＋diff 摘要

**改法（甲）**：`scripts/build.ps1:103-110` 的 `$ldflags` 数组在六枚 `-X` 之后追加一枚 `"-H=windowsgui"`，六枚 `-X` 逐字原样未动；数组仍整体 `-join ' '` 喂给 `:123` 的 `go build -trimpath -ldflags $ldflags -o (Join-Path $outDir 'wisp.exe') ./cmd/wisp`。另按派单「读该文件周边注释补一行 why」在数组上方补 5 行注释（why 见 diff）。

**diff 摘要**（`git diff -- scripts/build.ps1` 逐字，7 insertions / 1 deletion）：

```diff
diff --git a/scripts/build.ps1 b/scripts/build.ps1
index 3593e87b..101d39f4 100644
--- a/scripts/build.ps1
+++ b/scripts/build.ps1
@@ -100,13 +100,19 @@ if ($null -ne $git) {
     if ($LASTEXITCODE -eq 0 && $short) { $commit = $short }   ← 摘要行，原文是 `&&`（PowerShell 7 形）此处为摘要排版，逐字 diff 见 git
 }
 $buildDate = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
+# -H=windowsgui (SPEC-11 §2.2: no-arg launch = GUI): the win32 desktop binary must
+# not pop a black console window when double-clicked. CLI legs (`wisp run` etc.)
+# keep their stdout because attachParentConsole (cmd/wisp/console_windows.go:33-43)
+# returns early when a usable stdout already exists (redirected or parent console)
+# and otherwise attaches ATTACH_PARENT_PROCESS and rebinds the std handles.
 $ldflags = (
     "-X $BuildInfoPkg.Version=0.0.0-dev",
     "-X $BuildInfoPkg.Commit=$commit",
     "-X $BuildInfoPkg.BuildDate=$buildDate",
     "-X $BuildInfoPkg.DefaultEnv=$Env",
     "-X $BuildInfoPkg.SherpaOnnxVersion=$sherpaVersion",
-    "-X $BuildInfoPkg.OnnxRuntimeVersion=$ortVersion"
+    "-X $BuildInfoPkg.OnnxRuntimeVersion=$ortVersion",
+    "-H=windowsgui"
 ) -join ' '
```

**注释里的 why（逐字，即上面 + 行）**：win32 桌面程序不许带控制台窗（双击 `wisp.exe` 不得连黑框一起起）；run 腿 stdout 重定向由 `attachParentConsole` 判据链保住——`cmd/wisp/console_windows.go:33-43` 逐字：

```go
func attachParentConsole() {
	out, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err == nil && out != 0 && out != windows.InvalidHandle {
		return
	}
	// Fails harmlessly when already attached or when no parent console exists.
	_, _, _ = procAttachConsole.Call(uintptr(attachParentProcess))
	rebindStdHandle(windows.STD_OUTPUT_HANDLE, "CONOUT$")
	rebindStdHandle(windows.STD_ERROR_HANDLE, "CONOUT$")
	rebindStdHandle(windows.STD_INPUT_HANDLE, "CONIN$")
}
```

（`:33-37`：已有可用 stdout 就直接 return ⇒ 重定向/父控制台两形不碰；`:38-43`：否则 `AttachConsole(ATTACH_PARENT_PROCESS)`＋重绑三枚 std handle；explorer 拉起无父控制台时 `:38` 注释自陈 "Fails harmlessly"。）

**继承面确认**：`.github/workflows/ci.yml` 三处 step 与 `scripts/slo-check.ps1:104` 自建分支全部转手 `scripts/build.ps1`（census §1 #2/#4），本改自动跟随，CI 一字未动——这正是选甲的理由。

## §2 三发读数（逐字）

**产生新 exe 的命令与 flag 全文**（AC#1 具名口径：`build/**` 被 gitignore、产物不可自证，故交件必带命令）：
派单给的 `-OutDir` 参数在 build.ps1 里**不存在**（`.USAGE` 只收 `[-Env dev|prod]`，`$outDir` 在 `:112` 硬编码 `Join-Path $RepoRoot 'build'`）⇒ 为不把产物写进 `build/`（现有台上产物地盘），本腿在台件区放了**补丁副本** `.scratch/probes-244-r1/build-probe.ps1`（对 `61c526e0` 时刻的 build.ps1 全文逐字，仅改**三处纯位置行**：`$RepoRoot` 绝对路径、`fetch-deps.ps1` 锚回真 `scripts/`、`$outDir` 指向 `$PSScriptRoot\out`；**旗标零差异**）。真跑：

```
powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/probes-244-r1/build-probe.ps1 -Env dev
```

全链真跑通过：fetch-deps → go build（CGO_ENABLED=1，CC=msys2 gcc 16.2.0，go1.27.1）→ DLL 伴随 → SHA256SUMS → **wisp doctor 冒烟 PASS**（尾行逐字 `wisp doctor: PASS`）。产物 30,973,922 字节，mtime 2026-10-02 19:37。

**① 新 exe subsystem（断言目标）**——`objdump -p .scratch/probes-244-r1/out/wisp.exe | grep -i subsystem` 逐字：

```
MajorSubsystemVersion	10
MinorSubsystemVersion	0
Subsystem		00000002	(Windows GUI)
```

**② 旧 exe 对照（未动）**——`objdump -p build/wisp.exe | grep -i subsystem` 逐字：

```
MajorSubsystemVersion	10
MinorSubsystemVersion	0
Subsystem		00000003	(Windows CUI)
```

对照成立：改后 2（GUI）、改前台上产物 3（CUI）；`build/wisp.exe` 本腿一字未动（mtime 2026-09-30 23:43 保持原样）。

**③ 重定向验证（run 腿 stdout 活着）**：

- `./out/wisp.exe run --help > out.txt 2>&1` ⇒ **exit=2、out.txt=1782 字节非空**。⚠ `run --help` 没走 help 分支而是进了 run 腿：因 dev 环境配置未设（llm 未配置）按设计快速失败退 2，**未起任何任务进程、未进常驻**（out.txt 逐字首行 `time=... msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2`，含 `wisp run: 本机没有可交互控制台，本轮没有人能答复卡片…` 的降级声明与三段补救指引）。这一发证的是**重定向管道本身接住了 stdout+stderr**（否则 out.txt 必空）。
- 按派单预案换安全旗标逐个试：`./out/wisp.exe version > out-v.txt 2>&1` ⇒ **exit=0、out-v.txt=267 字节非空**（逐字：`wisp 0.0.0-dev (61c526e0, built 2026-10-02T11:36:57Z)` / `WISP_ENV=dev (data dir rules: SPEC-03 §5)` / `sherpa-onnx runtime version: 1.13.8`，前置一行 winsec INFO）；`./out/wisp.exe -h > out-h.txt 2>&1` ⇒ **exit=0、out-h.txt=2464 字节非空**（逐字首两行：`wisp - personal voice agent for Windows` / `Usage:`）。三发全部证明：GUI 子系统下 stdout 重定向路径活着（`console_windows.go:34-37` 判据链：重定向时 `GetStdHandle` 合法 ⇒ 直接 return，不碰 attach）。
- ⛔ 本腿全程未裸跑 `wisp run`（无参或带任务参数），未进常驻腿。

## §3 我可能写错的条目＋判不动的地方

**可能写错**：

1. **`run --help` 的语义归我猜**：out.txt 里没有 usage 文本而走了 run 腿的未配置失败路；`run` 子命令是否存在 `--help` 旗标我没有读 dispatch 逐字核。若它本该打 usage，则我这发"非空"的成分是错误输出而非帮助文本——**判定不变**（重定向活着），但"用 --help 验"这个动作本身选错了旗标，正式的干净凭据以 `version`/`-h` 两发为准。
2. **补丁副本的"仅三处位置行"声明**：我逐字核对过 diff 基线，但台件副本是手工三处 Edit 产生，若 build.ps1 在我 copy 与 commit 之间被他人推进，副本可能落后于 HEAD 版本——读数对 `61c526e0` 时刻的数组负责。
3. **GUI 子系统下 doctor 冒烟的 stdout 形态**：build-probe.ps1 收尾 `wisp doctor: PASS` 是脚本自打（PowerShell 控制台父进程 ⇒ attach 父控制台通路），这一发不能当作"explorer 拉起三形态"的真机凭据——票面 AC#2 的父控制台跑／`doctor > out.txt`、explorer 双击三形真跑仍欠（本腿只交了重定向形＋安全旗标形）。
4. **`build/wisp.exe` 的归属标签**：沿用 census §2.5 的保留（按 mtime 判它出自 build.ps1 链，未读版本资源段核对）——对"改前 CUI(3)"对照读数无影响。

**判不动（归后续格，本腿不做）**：

1. **票面 AC#2 三条真机读数**（父控制台里跑 run／explorer 双击无黑框且常驻照起）——需要真机交互形，本腿零真跑常驻。
2. **AC#3 两行注释改口**（`console_windows.go:26`／`console_other.go:7`）——写面在 `cmd/wisp`，且与 `228-r1` 串行禁区（票面排程节），本腿只动 build.ps1 一枚。
3. **AC#5 ⓐ′／ⓑ 两截**（GUI 任务源＋停机源前置、会响的尺）——票面明文"未答之前整票不许动构建链"由编排者排程裁，本腿按派单 A546 落甲形，前置账不归我销。
4. **AC#4 文档过期行登记**（`BUILD.md:51/:67-71/:87`）——`BUILD.md` 冻结件一字不动；具名上报：其 `:87` "推迟到票 07" 两行追述在本改落地后即成过期描述，归 G6 裁的并票 225 族一格。
