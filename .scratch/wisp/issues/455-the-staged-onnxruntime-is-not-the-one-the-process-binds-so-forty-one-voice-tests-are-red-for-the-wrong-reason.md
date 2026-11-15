# 455 — 整包 41 枚语音红：staged 的 17.8 MB onnxruntime 没被绑上，System32 的 1.17.1 抢先

**Status:** ready-for-agent
**Found by:** 编排者（2026-11-15 19:3x，来路＝票 35 G3 写腿 `3a8e776ac` 的整包读数）
**Spec refs:** SPEC-04（语音链路）、SPEC-10（测试与验收）

## What is wrong

票 35 G3 的整包一发：`PASS=618／FAIL=41／SKIP=10`，752.9 秒。那 41 枚红的报错逐字都带同一句：

> `speech: onnxruntime 1.17.1 is older than the 1.28 … a second onnxruntime.dll is shadowing the one on PATH`

**但载具本来是防这个的**：
- `scripts/wisp-cli-staged-tests.sh:230-232` 在测试二进制旁边摆 pinned DLL（读 `deps.toml`）
- `.scratch/wisp/stage-cli/onnxruntime.dll` 现读 **17799168 B**（＝仓根那枚，2026-09-11），**不是** System32 的 10572960 B
- 那载具的头注释 `:12-18` 逐字写着：「a host whose `C:\Windows\System32` carries an OS-provided `onnxruntime.dll`（measured on this machine: **1.17.260311, 10,572,960 bytes**）binds THAT copy…」

⇒ **"摆了正确的 DLL 却仍然绑到 System32 那一枚"** —— 这条连线里有一环没成立。

## Why it matters

1. **41 枚测试的读数今天不可用**：它们红的是环境，不是代码。任何人拿"整包 41 红"当产品结论都会错。
2. 它**污染每一次全量读数**：今天已经有两条腿的整包读数被它搅（其中一发还叠加了票 314 的整态死亡）。
3. 载具里那段注释现在**是假话**（"摆了 pinned DLL 就不会被抢"）—— A1874 家族。

## Acceptance criteria

- [ ] **AC#1 先证"绑的是哪一枚"**：跑一发能打印实际加载的 `onnxruntime.dll` **路径与版本**的读数（⛔ 不许用"看报错说 1.17.1"当证据 —— 那只是被影子那一枚的**版本号**，不是路径）。红句原文抄进证据：至少给出「`stage-cli/onnxruntime.dll` 是 17799168 B 而进程绑到 10572960 B」这一对数字。
- [ ] **AC#2 找到那一环**：说清为什么"应用目录里有 17.8 MB 却绑了 System32 的 10.5 MB"。候选（**逐条现读判，⛔ 猜**）：①`stage-sherpa-dlls.sh` 当发**拒了**（如模型缺失）但载具继续跑 ②sherpa 的 c-api DLL 用 `LoadLibrary` 按**自己的**搜索序找 onnxruntime（不是测试二进制所在目录优先）③`deps.toml` 里 pin 的与 stage 的不是同一枚 ④Windows 的 `SafeDllSearchMode`／`SetDllDirectory` 被某处改过 ⑤测试里有人 `os.Setenv("PATH", …)` 或 `LoadLibrary` 绝对路径写错。
- [ ] **AC#3 修**（按主流 harness 的形状：显式搜索路径 / `LoadLibraryEx` + `LOAD_LIBRARY_SEARCH_APPLICATION_DIR` / 版本前置检查 fail-fast 三者择一，**写明为什么选它**），并给**会红的测试**（把 pinned DLL 挪走 ⇒ 必须红，且红句说得出"绑错版本"而不是"崩了"）。
- [ ] **AC#4 载具那段注释改如实**：`scripts/wisp-cli-staged-tests.sh:12-18` 现在断言"摆了就不会被抢"—— 修好后要么它变成真的（并给读数），要么改成它真正保证的东西。
- [ ] **AC#5 门禁**：`d22scan` rc=0 clean／path-length GREEN／`go vet` rc=0；**⛔ 动 `deps.toml` 里 pin 的版本号**（那不是本票的活）；**⛔ 动 System32**（不改机器，改我们这一侧）。

## Hard constraints

- ⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 表／`frontend/**`／`design/**`
- ⛔ 删 System32 里的任何东西（那是机器，不是仓）
- ⛔ 放宽任何断言来变绿；⛔ 把 41 枚标 `t.Skip` 糊过去
- 两个 `sh` 门用 `D:\work\soft\Git\bin\bash.exe`（裸 `bash` 命中 WSL 会假红，台账 A1884）
